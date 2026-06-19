package bot

import (
	"testing"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

func TestHandleTimelineEmpty(t *testing.T) {
	h := New(&fakeService{
		timeline: func(int64, int64) (*incident.Incident, []event.Event, error) {
			return &incident.Incident{Title: "outage"}, nil, nil
		},
	}, newFakeAPI())
	ctx := &mockContext{}

	if err := h.HandleTimeline(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "timeline is empty")
}

func TestHandleTimelineNoActiveIncident(t *testing.T) {
	h := New(&fakeService{
		timeline: func(int64, int64) (*incident.Incident, []event.Event, error) {
			return nil, nil, errs.ErrNoActiveIncident
		},
	}, newFakeAPI())
	ctx := &mockContext{}

	if err := h.HandleTimeline(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "no active incident")
}
