package bot

import (
	"testing"

	"gopkg.in/telebot.v3"
)

type mockContext struct {
	telebot.Context
	args []string
	sent []string
}

func (m *mockContext) Args() []string {
	return m.args
}

func (m *mockContext) Send(what interface{}, opts ...interface{}) error {
	m.sent = append(m.sent, what.(string))
	return nil
}

func lastSent(t *testing.T, m *mockContext) string {
	t.Helper()
	if len(m.sent) != 1 {
		t.Fatalf("expected exactly 1 message sent, got %d: %v", len(m.sent), m.sent)
	}
	return m.sent[0]
}
