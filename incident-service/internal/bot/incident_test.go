package bot

import (
	"strings"
	"testing"
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

func TestHandleIncidentNoArgsShowsUsage(t *testing.T) {
	h := New(&fakeService{})
	ctx := &mockContext{}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "Usage:")
}

func TestHandleIncidentUsageListsAllSubcommands(t *testing.T) {
	h := New(&fakeService{})
	ctx := &mockContext{}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	usage := lastSent(t, ctx)
	for _, part := range []string{"/incident create", "/incident close", "/incident <message>"} {
		if !strings.Contains(usage, part) {
			t.Errorf("usage %q does not mention %q", usage, part)
		}
	}
}

func TestHandleIncidentCreateWithoutDescriptionAsksForOne(t *testing.T) {
	h := New(&fakeService{})
	ctx := &mockContext{args: []string{"create"}}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "Please add a description")
}

func TestHandleIncidentCreatePersistsAndShowsCard(t *testing.T) {
	var gotTitle string
	h := New(&fakeService{
		create: func(chatID int64, title string, sev incident.Severity, _ *int64, _ string) (*incident.Incident, error) {
			gotTitle = title
			return &incident.Incident{Title: title, Severity: incident.SeverityMedium, Status: incident.StatusActive}, nil
		},
	})
	ctx := &mockContext{args: []string{"create", "db", "is", "down"}}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotTitle != "db is down" {
		t.Errorf("service got title %q, want %q", gotTitle, "db is down")
	}
	sentContains(t, ctx, "db is down")
}

func TestHandleIncidentCreateReportsConflict(t *testing.T) {
	h := New(&fakeService{
		create: func(int64, string, incident.Severity, *int64, string) (*incident.Incident, error) {
			return nil, errs.ErrIncidentAlreadyActive
		},
	})
	ctx := &mockContext{args: []string{"create", "again"}}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "already active")
}

func TestHandleIncidentMessageAddsTimelineUpdate(t *testing.T) {
	var gotMsg string
	h := New(&fakeService{
		addEvent: func(_ int64, _ *int64, _, message string) (*event.Event, error) {
			gotMsg = message
			return &event.Event{}, nil
		},
	})
	ctx := &mockContext{args: []string{"db", "is", "down"}}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMsg != "db is down" {
		t.Errorf("service got message %q, want %q", gotMsg, "db is down")
	}
	sentContains(t, ctx, "Update added to the timeline")
}

func TestHandleIncidentCloseSendsSummaryAndReport(t *testing.T) {
	now := time.Now()
	h := New(&fakeService{
		report: func(int64) ([]byte, error) { return []byte("%PDF-1.4 fake"), nil },
		closeInc: func(int64, *int64, string) (*incident.Incident, error) {
			return &incident.Incident{Title: "outage", CreatedAt: now, ClosedAt: &now}, nil
		},
	})
	ctx := &mockContext{args: []string{"close"}}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "Incident closed")
	sentContains(t, ctx, "telebot.Document")
}

func TestHandleIncidentCloseStillClosesWhenReportFails(t *testing.T) {
	now := time.Now()
	h := New(&fakeService{
		report: func(int64) ([]byte, error) { return nil, errs.New(errs.KindUnavailable, "report", "down") },
		closeInc: func(int64, *int64, string) (*incident.Incident, error) {
			return &incident.Incident{Title: "outage", CreatedAt: now, ClosedAt: &now}, nil
		},
	})
	ctx := &mockContext{args: []string{"close"}}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "Incident closed")
	sentContains(t, ctx, "report could not be generated")
}

func TestHandleIncidentCloseReportsNoActive(t *testing.T) {
	h := New(&fakeService{
		report:   func(int64) ([]byte, error) { return nil, errs.ErrNoActiveIncident },
		closeInc: func(int64, *int64, string) (*incident.Incident, error) { return nil, errs.ErrNoActiveIncident },
	})
	ctx := &mockContext{args: []string{"close"}}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "no active incident")
}
