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

type fakeOpener struct {
	opened []openedIncident
}

func (f *fakeOpener) OpenIncidentFromAlert(_ context.Context, title string, severity incident.Severity) (*incident.Incident, error) {
	f.opened = append(f.opened, openedIncident{title: title, severity: severity})
	return &incident.Incident{Title: title, Severity: severity}, nil
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
