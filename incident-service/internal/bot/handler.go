// Package bot is the Telegram-facing layer of the incident war room. It
// translates telebot updates into service use-case calls and renders the
// results back into chat messages. It owns no business rules: every command
// and inline-panel action delegates to the IncidentService.
package bot

import (
	"context"
	"errors"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

// handlerTimeout bounds every use-case call triggered by a Telegram update so a
// slow database or report service cannot wedge a handler forever.
const handlerTimeout = 30 * time.Second

// IncidentService is the set of use cases the bot drives. It is implemented by
// *service.Service; the bot depends only on this interface so it can be tested
// with a fake.
type IncidentService interface {
	CreateIncident(ctx context.Context, chatID, topicID int64, title string, severity incident.Severity, userID *int64, username string) (*incident.Incident, error)
	AddTimelineEvent(ctx context.Context, chatID, topicID int64, userID *int64, username, message string) (*event.Event, error)
	CloseIncident(ctx context.Context, chatID, topicID int64, userID *int64, username string) (*incident.Incident, error)
	SetSeverity(ctx context.Context, chatID, topicID int64, severity incident.Severity) (*incident.Incident, error)
	GetTimeline(ctx context.Context, chatID, topicID int64) (*incident.Incident, []event.Event, error)
	GenerateReport(ctx context.Context, chatID, topicID int64) ([]byte, error)
}

// TelegramAPI is the slice of telebot.Bot the handler needs to manage forum
// topics and post messages to a specific chat/thread. *telebot.Bot satisfies
// it; tests provide a fake. Using it keeps topic-aware handlers off the
// concrete c.Bot().
type TelegramAPI interface {
	Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error)
	CreateTopic(chat *telebot.Chat, topic *telebot.Topic) (*telebot.Topic, error)
	DeleteTopic(chat *telebot.Chat, topic *telebot.Topic) error
}

// Handler wires Telegram updates to the incident use cases.
type Handler struct {
	svc IncidentService
	api TelegramAPI
}

// New returns a Handler backed by svc and the Telegram API.
func New(svc IncidentService, api TelegramAPI) *Handler {
	return &Handler{svc: svc, api: api}
}

// reqContext derives a bounded context for a single update.
func reqContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), handlerTimeout)
}

// threadID returns the forum topic (message_thread_id) the update arrived in.
// For messages outside any topic it is 0 (the chat's General thread).
func threadID(c telebot.Context) int64 {
	if m := c.Message(); m != nil {
		return int64(m.ThreadID)
	}
	return 0
}

// sender extracts the Telegram user of an update as (id, username). The id is
// returned as a pointer so it can be stored as a nullable column; it is nil for
// updates without a sender.
func sender(c telebot.Context) (*int64, string) {
	u := c.Sender()
	if u == nil {
		return nil, ""
	}
	id := u.ID
	return &id, u.Username
}

// userError turns a service error into a chat-friendly message. The internal
// op/cause detail is intentionally dropped — it is logged elsewhere, not shown
// to users.
func userError(err error) string {
	switch {
	case errors.Is(err, errs.ErrNoActiveIncident):
		return "There is no active incident in this chat. Open one with /incident create <description>."
	case errors.Is(err, errs.ErrIncidentAlreadyActive):
		return "An incident is already active in this chat. Close it before opening a new one."
	case errs.Is(err, errs.KindValidation):
		return "Sorry, that input is not valid. " + incidentUsage
	case errs.Is(err, errs.KindUnavailable):
		return "The service is temporarily unavailable. Please try again in a moment."
	default:
		return "Something went wrong. Please try again."
	}
}
