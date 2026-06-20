package service

import (
	"context"
	"testing"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

func TestSetSeverity(t *testing.T) {
	ctx := context.Background()

	t.Run("updates the active incident's severity", func(t *testing.T) {
		svc, incidents, _ := newTestService()
		if _, err := svc.CreateIncident(ctx, 200, 200, "DB is down", incident.SeverityLow, nil, "alice"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		inc, err := svc.SetSeverity(ctx, 200, 200, incident.SeverityHigh)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if inc.Severity != incident.SeverityHigh {
			t.Fatalf("returned severity = %q, want HIGH", inc.Severity)
		}

		got, _ := incidents.GetActiveByTopicID(ctx, 200, 200)
		if got.Severity != incident.SeverityHigh {
			t.Fatalf("stored severity = %q, want HIGH", got.Severity)
		}
	})

	t.Run("invalid severity is rejected", func(t *testing.T) {
		svc, _, _ := newTestService()
		if _, err := svc.CreateIncident(ctx, 201, 201, "outage", incident.SeverityLow, nil, "alice"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err := svc.SetSeverity(ctx, 201, 201, incident.Severity("CRITICAL"))
		if errs.KindOf(err) != errs.KindValidation {
			t.Fatalf("expected validation error, got %v", err)
		}
	})

	t.Run("no active incident", func(t *testing.T) {
		svc, _, _ := newTestService()

		_, err := svc.SetSeverity(ctx, 999, 999, incident.SeverityHigh)
		if errs.KindOf(err) != errs.KindNotFound {
			t.Fatalf("expected not-found, got %v", err)
		}
	})
}
