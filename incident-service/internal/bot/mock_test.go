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
	args   []string
	chatID int64
	user   *telebot.User
	sent   []string
}

func (m *mockContext) Args() []string { return m.args }

func (m *mockContext) Chat() *telebot.Chat { return &telebot.Chat{ID: m.chatID} }

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
	create   func(chatID int64, title string, sev incident.Severity, authorID *int64, username string) (*incident.Incident, error)
	addEvent func(chatID int64, authorID *int64, username, message string) (*event.Event, error)
	closeInc func(chatID int64, authorID *int64, username string) (*incident.Incident, error)
	setSev   func(chatID int64, sev incident.Severity) (*incident.Incident, error)
	timeline func(chatID int64) (*incident.Incident, []event.Event, error)
	report   func(chatID int64) ([]byte, error)
}

func (f *fakeService) CreateIncident(_ context.Context, chatID int64, title string, sev incident.Severity, authorID *int64, username string) (*incident.Incident, error) {
	return f.create(chatID, title, sev, authorID, username)
}

func (f *fakeService) AddTimelineEvent(_ context.Context, chatID int64, authorID *int64, username, message string) (*event.Event, error) {
	return f.addEvent(chatID, authorID, username, message)
}

func (f *fakeService) CloseIncident(_ context.Context, chatID int64, authorID *int64, username string) (*incident.Incident, error) {
	return f.closeInc(chatID, authorID, username)
}

func (f *fakeService) SetSeverity(_ context.Context, chatID int64, sev incident.Severity) (*incident.Incident, error) {
	return f.setSev(chatID, sev)
}

func (f *fakeService) GetTimeline(_ context.Context, chatID int64) (*incident.Incident, []event.Event, error) {
	return f.timeline(chatID)
}

func (f *fakeService) GenerateReport(_ context.Context, chatID int64) ([]byte, error) {
	return f.report(chatID)
}
