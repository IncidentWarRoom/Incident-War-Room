package reportclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/cQu1x/Incident-War-Room/internal/domain/report"
	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

func sampleReport() report.Report {
	closedAt := time.Date(2026, 6, 1, 14, 18, 0, 0, time.UTC)
	user := int64(1)
	return report.Report{
		Incident: incident.Incident{
			ID:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Title:     "Payment Service Down",
			Severity:  incident.SeverityHigh,
			Status:    incident.StatusClosed,
			CreatedAt: time.Date(2026, 6, 1, 14, 3, 0, 0, time.UTC),
			ClosedAt:  &closedAt,
		},
		Participants: []report.Participant{{UserID: 1, Username: "rolan"}},
		Timeline: []event.Event{{
			UserID:    &user,
			Username:  "rolan",
			Message:   "Started investigating database issues",
			CreatedAt: time.Date(2026, 6, 1, 14, 5, 0, 0, time.UTC),
		}},
	}
}

func TestGenerateSendsContractAndReturnsPDF(t *testing.T) {
	var gotPath, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Content-Type", "application/pdf")
		_, _ = w.Write([]byte("%PDF-1.4 fake"))
	}))
	defer srv.Close()

	pdf, err := New(srv.URL).Generate(context.Background(), sampleReport())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(pdf) != "%PDF-1.4 fake" {
		t.Fatalf("unexpected pdf bytes: %q", pdf)
	}
	if gotPath != generatePath {
		t.Fatalf("expected path %q, got %q", generatePath, gotPath)
	}

	// The body must match the report-service wire contract (camelCase keys,
	// RFC3339 timestamps).
	var decoded map[string]any
	if err := json.Unmarshal([]byte(gotBody), &decoded); err != nil {
		t.Fatalf("body is not valid JSON: %v", err)
	}
	inc := decoded["incident"].(map[string]any)
	if inc["id"] != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("unexpected incident id: %v", inc["id"])
	}
	if inc["createdAt"] != "2026-06-01T14:03:00Z" {
		t.Fatalf("unexpected createdAt: %v", inc["createdAt"])
	}
	if inc["closedAt"] != "2026-06-01T14:18:00Z" {
		t.Fatalf("unexpected closedAt: %v", inc["closedAt"])
	}
	parts := decoded["participants"].([]any)
	if len(parts) != 1 || parts[0].(map[string]any)["userId"].(float64) != 1 {
		t.Fatalf("unexpected participants: %v", parts)
	}
}

func TestGenerateMapsStatusToErrorKind(t *testing.T) {
	cases := []struct {
		status int
		want   errs.Kind
	}{
		{http.StatusUnprocessableEntity, errs.KindValidation},
		{http.StatusNotFound, errs.KindNotFound},
		{http.StatusInternalServerError, errs.KindUnavailable},
	}
	for _, tc := range cases {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(tc.status)
			_, _ = w.Write([]byte("boom"))
		}))

		_, err := New(srv.URL).Generate(context.Background(), sampleReport())
		if errs.KindOf(err) != tc.want {
			t.Errorf("status %d: expected kind %s, got %v", tc.status, tc.want, err)
		}
		srv.Close()
	}
}

func TestGenerateUnreachableServiceIsUnavailable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	url := srv.URL
	srv.Close()

	_, err := New(url).Generate(context.Background(), sampleReport())
	if errs.KindOf(err) != errs.KindUnavailable {
		t.Fatalf("expected unavailable, got %v", err)
	}
}
