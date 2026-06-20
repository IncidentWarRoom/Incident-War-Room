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

func TestHandleTimelineRepliesInTheTopic(t *testing.T) {
	h := New(&fakeService{
		timeline: func(_, topicID int64) (*incident.Incident, []event.Event, error) {
			return &incident.Incident{Title: "outage", TopicID: topicID}, nil, nil
		},
	}, newFakeAPI())
	ctx := &mockContext{threadID: 444}

	if err := h.HandleTimeline(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ctx.sentThread) != 1 || ctx.sentThread[0] != 444 {
		t.Fatalf("expected timeline sent to topic 444, got threads %v", ctx.sentThread)
	}
}

func TestHandleTimelineLinksTelegraph(t *testing.T) {
	h := New(&fakeService{
		timeline: func(int64, int64) (*incident.Incident, []event.Event, error) {
			return &incident.Incident{Title: "outage"}, []event.Event{
				{Username: "alice", Message: "looking into it"},
			}, nil
		},
		publish: func(int64, int64) ([]string, error) {
			return []string{"https://telegra.ph/timeline-1"}, nil
		},
	}, newFakeAPI())
	ctx := &mockContext{}

	if err := h.HandleTimeline(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "https://telegra.ph/timeline-1")
}

func TestHandleTimelineShowsTelegraphFailure(t *testing.T) {
	h := New(&fakeService{
		timeline: func(int64, int64) (*incident.Incident, []event.Event, error) {
			return &incident.Incident{Title: "outage"}, []event.Event{
				{Username: "alice", Message: "looking into it"},
			}, nil
		},
		publish: func(int64, int64) ([]string, error) {
			return nil, errs.New(errs.KindUnavailable, "telegraph", "down")
		},
	}, newFakeAPI())
	ctx := &mockContext{}

	if err := h.HandleTimeline(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "Telegraph is unavailable")
}

func TestHandleTimelineEmptyDoesNotPublish(t *testing.T) {
	published := false
	h := New(&fakeService{
		timeline: func(int64, int64) (*incident.Incident, []event.Event, error) {
			return &incident.Incident{Title: "outage"}, nil, nil
		},
		publish: func(int64, int64) ([]string, error) {
			published = true
			return nil, nil
		},
	}, newFakeAPI())

	if err := h.HandleTimeline(&mockContext{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if published {
		t.Fatal("empty timeline should not be published to Telegraph")
	}
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
