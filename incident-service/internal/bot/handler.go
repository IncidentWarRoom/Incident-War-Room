// Package bot is the Telegram-facing layer of the incident war room. It
// translates telebot updates into service use-case calls and renders the
// results back into chat messages. It owns no business rules: every command
// and inline-panel action delegates to the IncidentService.
package bot

import (
	"context"
	"errors"
	"sync"
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

type TelegramAPI interface {
	Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error)
	Edit(msg telebot.Editable, what interface{}, opts ...interface{}) (*telebot.Message, error)
	CreateTopic(chat *telebot.Chat, topic *telebot.Topic) (*telebot.Topic, error)
	DeleteTopic(chat *telebot.Chat, topic *telebot.Topic) error
}

// announceKey identifies the main-chat announcement of a single incident.
type announceKey struct {
	chatID  int64
	topicID int64
}

// Handler wires Telegram updates to the incident use cases.
type Handler struct {
	svc IncidentService
	api TelegramAPI

	// announcements remembers the main-chat message announcing each incident so
	// it can be edited in place when the incident metadata (e.g. severity)
	// changes from inside the topic. It is best-effort, in-memory state: a bot
	// restart simply leaves older announcements static.
	mu            sync.Mutex
	announcements map[announceKey]telebot.Editable
}

// New returns a Handler backed by svc and the Telegram API.
func New(svc IncidentService, api TelegramAPI) *Handler {
	return &Handler{
		svc:           svc,
		api:           api,
		announcements: make(map[announceKey]telebot.Editable),
	}
}

// rememberAnnouncement stores the main-chat announcement message for an incident.
func (h *Handler) rememberAnnouncement(chatID, topicID int64, msg telebot.Editable) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.announcements[announceKey{chatID, topicID}] = msg
}

// announcement returns the stored main-chat announcement message for an incident.
func (h *Handler) announcement(chatID, topicID int64) (telebot.Editable, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	msg, ok := h.announcements[announceKey{chatID, topicID}]
	return msg, ok
}

// forgetAnnouncement drops the stored announcement for a closed incident.
func (h *Handler) forgetAnnouncement(chatID, topicID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.announcements, announceKey{chatID, topicID})
}

// reqContext derives a bounded context for a single update.
func reqContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), handlerTimeout)
}

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
