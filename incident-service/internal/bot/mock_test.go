package bot

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"gopkg.in/telebot.v3"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

// mockContext is a minimal telebot.Context for handler tests. It records the
// messages a handler sends and lets a test set the chat and sender.
type mockContext struct {
	telebot.Context
	args     []string
	chatID   int64
	threadID int64
	user     *telebot.User
	sent     []string
}

func (m *mockContext) Args() []string { return m.args }

func (m *mockContext) Chat() *telebot.Chat { return &telebot.Chat{ID: m.chatID} }

func (m *mockContext) Message() *telebot.Message {
	return &telebot.Message{ThreadID: int(m.threadID)}
}

func (m *mockContext) Sender() *telebot.User { return m.user }

// Send records string payloads verbatim and other payloads (e.g. a PDF
// document) by their type, so tests can assert on either.
func (m *mockContext) Send(what interface{}, opts ...interface{}) error {
	if s, ok := what.(string); ok {
		m.sent = append(m.sent, s)
	} else {
		m.sent = append(m.sent, fmt.Sprintf("<%T>", what))
	}
	return nil
}

func lastSent(t *testing.T, m *mockContext) string {
	t.Helper()
	if len(m.sent) != 1 {
		t.Fatalf("expected exactly 1 message sent, got %d: %v", len(m.sent), m.sent)
	}
	return m.sent[0]
}

func sentContains(t *testing.T, m *mockContext, substr string) {
	t.Helper()
	for _, s := range m.sent {
		if strings.Contains(s, substr) {
			return
		}
	}
	t.Fatalf("no sent message contains %q; sent: %v", substr, m.sent)
}

// fakeService is a configurable IncidentService for tests. Unset hooks return
// zero values.
type fakeService struct {
	create   func(chatID, topicID int64, title string, sev incident.Severity, userID *int64, username string) (*incident.Incident, error)
	addEvent func(chatID, topicID int64, userID *int64, username, message string) (*event.Event, error)
	closeInc func(chatID, topicID int64, userID *int64, username string) (*incident.Incident, error)
	setSev   func(chatID, topicID int64, sev incident.Severity) (*incident.Incident, error)
	timeline func(chatID, topicID int64) (*incident.Incident, []event.Event, error)
	report   func(chatID, topicID int64) ([]byte, error)
}

func (f *fakeService) CreateIncident(_ context.Context, chatID, topicID int64, title string, sev incident.Severity, userID *int64, username string) (*incident.Incident, error) {
	return f.create(chatID, topicID, title, sev, userID, username)
}

func (f *fakeService) AddTimelineEvent(_ context.Context, chatID, topicID int64, userID *int64, username, message string) (*event.Event, error) {
	return f.addEvent(chatID, topicID, userID, username, message)
}

func (f *fakeService) CloseIncident(_ context.Context, chatID, topicID int64, userID *int64, username string) (*incident.Incident, error) {
	return f.closeInc(chatID, topicID, userID, username)
}

func (f *fakeService) SetSeverity(_ context.Context, chatID, topicID int64, sev incident.Severity) (*incident.Incident, error) {
	return f.setSev(chatID, topicID, sev)
}

func (f *fakeService) GetTimeline(_ context.Context, chatID, topicID int64) (*incident.Incident, []event.Event, error) {
	return f.timeline(chatID, topicID)
}

func (f *fakeService) GenerateReport(_ context.Context, chatID, topicID int64) ([]byte, error) {
	return f.report(chatID, topicID)
}

// fakeAPI is a configurable TelegramAPI for tests. It records topic lifecycle
// calls and the payloads sent, so tests can assert where messages land.
type fakeAPI struct {
	createdTopic *telebot.Topic
	createErr    error
	deleted      []int // ThreadIDs passed to DeleteTopic
	sent         []sentMessage
}

// sentMessage records one api.Send call: its thread (0 == General) and a
// string/type rendering of the payload, mirroring mockContext.Send.
type sentMessage struct {
	threadID int
	what     string
}

func newFakeAPI() *fakeAPI {
	return &fakeAPI{createdTopic: &telebot.Topic{Name: "topic", ThreadID: 555}}
}

func (a *fakeAPI) Send(_ telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error) {
	var thread int
	for _, o := range opts {
		if so, ok := o.(*telebot.SendOptions); ok {
			thread = so.ThreadID
		}
	}

	msg := sentMessage{threadID: thread}
	if s, ok := what.(string); ok {
		msg.what = s
	} else {
		msg.what = fmt.Sprintf("<%T>", what)
	}
	a.sent = append(a.sent, msg)
	return &telebot.Message{}, nil
}

func (a *fakeAPI) CreateTopic(_ *telebot.Chat, topic *telebot.Topic) (*telebot.Topic, error) {
	if a.createErr != nil {
		return nil, a.createErr
	}
	a.createdTopic.Name = topic.Name
	return a.createdTopic, nil
}

func (a *fakeAPI) DeleteTopic(_ *telebot.Chat, topic *telebot.Topic) error {
	a.deleted = append(a.deleted, topic.ThreadID)
	return nil
}

// apiSentContains fails the test unless some api.Send payload contains substr.
func apiSentContains(t *testing.T, a *fakeAPI, substr string) {
	t.Helper()
	for _, s := range a.sent {
		if strings.Contains(s.what, substr) {
			return
		}
	}
	t.Fatalf("no api-sent message contains %q; sent: %v", substr, a.sent)
}
