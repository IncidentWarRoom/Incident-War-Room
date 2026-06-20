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

const handlerTimeout = 30 * time.Second

type IncidentService interface {
	CreateIncident(ctx context.Context, chatID, topicID int64, title string, severity incident.Severity, userID *int64, username string) (*incident.Incident, error)
	AddTimelineEvent(ctx context.Context, chatID, topicID int64, userID *int64, username, message string) (*event.Event, error)
	CloseIncident(ctx context.Context, chatID, topicID int64, userID *int64, username string) (*incident.Incident, error)
	SetSeverity(ctx context.Context, chatID, topicID int64, severity incident.Severity) (*incident.Incident, error)
	GetTimeline(ctx context.Context, chatID, topicID int64) (*incident.Incident, []event.Event, error)
	PublishTimeline(ctx context.Context, chatID, topicID int64) ([]string, error)
	GenerateReport(ctx context.Context, chatID, topicID int64) ([]byte, error)
}

type TelegramAPI interface {
	Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error)
	Edit(msg telebot.Editable, what interface{}, opts ...interface{}) (*telebot.Message, error)
	CreateTopic(chat *telebot.Chat, topic *telebot.Topic) (*telebot.Topic, error)
	DeleteTopic(chat *telebot.Chat, topic *telebot.Topic) error
}

type announceKey struct {
	chatID  int64
	topicID int64
}

type Handler struct {
	svc IncidentService
	api TelegramAPI

	mu            sync.Mutex
	announcements map[announceKey]telebot.Editable
}

func New(svc IncidentService, api TelegramAPI) *Handler {
	return &Handler{
		svc:           svc,
		api:           api,
		announcements: make(map[announceKey]telebot.Editable),
	}
}

func (h *Handler) rememberAnnouncement(chatID, topicID int64, msg telebot.Editable) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.announcements[announceKey{chatID, topicID}] = msg
}

func (h *Handler) announcement(chatID, topicID int64) (telebot.Editable, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	msg, ok := h.announcements[announceKey{chatID, topicID}]
	return msg, ok
}

func (h *Handler) forgetAnnouncement(chatID, topicID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.announcements, announceKey{chatID, topicID})
}

func reqContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), handlerTimeout)
}

func threadID(c telebot.Context) int64 {
	if m := c.Message(); m != nil {
		return int64(m.ThreadID)
	}
	return 0
}

func sender(c telebot.Context) (*int64, string) {
	u := c.Sender()
	if u == nil {
		return nil, ""
	}
	id := u.ID
	return &id, u.Username
}

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
