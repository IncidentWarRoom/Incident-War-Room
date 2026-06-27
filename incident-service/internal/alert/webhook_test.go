package alert

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

type openedIncident struct {
	title    string
	severity incident.Severity
}

type closedIncident struct {
	chatID  int64
	topicID int64
}

type fakeOpener struct {
	opened    []openedIncident
	closed    []closedIncident
	nextTopic int64
	chatID    int64
}

func (f *fakeOpener) OpenIncidentFromAlert(_ context.Context, title string, severity incident.Severity) (*incident.Incident, error) {
	f.opened = append(f.opened, openedIncident{title: title, severity: severity})
	f.nextTopic++
	return &incident.Incident{Title: title, Severity: severity, ChatID: f.chatID, TopicID: f.nextTopic}, nil
}

func (f *fakeOpener) CloseIncidentFromAlert(_ context.Context, chatID, topicID int64) error {
	f.closed = append(f.closed, closedIncident{chatID: chatID, topicID: topicID})
	return nil
}

func post(t *testing.T, h http.Handler, token, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/webhooks/alertmanager", strings.NewReader(body))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestFiringAlertsOpenIncidents(t *testing.T) {
	opener := &fakeOpener{}
	body := `{"alerts":[
		{"status":"firing","labels":{"alertname":"HighCPU","severity":"critical"},"annotations":{"summary":"CPU is on fire"}},
		{"status":"resolved","labels":{"alertname":"HighCPU","severity":"critical"}},
		{"status":"firing","labels":{"alertname":"DiskFilling","severity":"warning"}}
	]}`

	rec := post(t, NewHandler(opener, ""), "", body)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if len(opener.opened) != 2 {
		t.Fatalf("opened %d incidents, want 2", len(opener.opened))
	}
	if opener.opened[0].title != "CPU is on fire" || opener.opened[0].severity != incident.SeverityHigh {
		t.Errorf("first incident = %+v", opener.opened[0])
	}
	if opener.opened[1].title != "DiskFilling" || opener.opened[1].severity != incident.SeverityMedium {
		t.Errorf("second incident = %+v", opener.opened[1])
	}
}

func TestResolvedAlertClosesIncident(t *testing.T) {
	opener := &fakeOpener{chatID: -100123}
	h := NewHandler(opener, "")

	firing := `{"alerts":[{"status":"firing","fingerprint":"abc","labels":{"alertname":"HighCPU","severity":"critical"}}]}`
	if rec := post(t, h, "", firing); rec.Code != http.StatusOK {
		t.Fatalf("firing status = %d, want 200", rec.Code)
	}

	resolved := `{"alerts":[{"status":"resolved","fingerprint":"abc","labels":{"alertname":"HighCPU","severity":"critical"}}]}`
	if rec := post(t, h, "", resolved); rec.Code != http.StatusOK {
		t.Fatalf("resolved status = %d, want 200", rec.Code)
	}

	if len(opener.closed) != 1 {
		t.Fatalf("closed %d incidents, want 1", len(opener.closed))
	}
	if opener.closed[0] != (closedIncident{chatID: -100123, topicID: 1}) {
		t.Errorf("closed = %+v, want {chatID:-100123 topicID:1}", opener.closed[0])
	}

	if rec := post(t, h, "", resolved); rec.Code != http.StatusOK {
		t.Fatalf("second resolved status = %d, want 200", rec.Code)
	}
	if len(opener.closed) != 1 {
		t.Fatalf("resolved without a known firing alert closed an incident: %+v", opener.closed)
	}
}

func TestTokenRequired(t *testing.T) {
	opener := &fakeOpener{}
	h := NewHandler(opener, "secret")

	if rec := post(t, h, "wrong", `{"alerts":[]}`); rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
	if rec := post(t, h, "secret", `{"alerts":[]}`); rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if len(opener.opened) != 0 {
		t.Fatalf("opened %d incidents, want 0", len(opener.opened))
	}
}

func TestInvalidPayload(t *testing.T) {
	rec := post(t, NewHandler(&fakeOpener{}, ""), "", "not json")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}
