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
	h := New(&fakeService{}, newFakeAPI())
	ctx := &mockContext{}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "Usage:")
}

func TestHandleIncidentUsageListsAllSubcommands(t *testing.T) {
	h := New(&fakeService{}, newFakeAPI())
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
	h := New(&fakeService{}, newFakeAPI())
	ctx := &mockContext{args: []string{"create"}}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "Please add a description")
}

func TestHandleIncidentCreateOpensTopicAndShowsCard(t *testing.T) {
	var gotTitle string
	var gotTopicID int64
	api := newFakeAPI()
	api.createdTopic.ThreadID = 777
	h := New(&fakeService{
		create: func(_ int64, topicID int64, title string, sev incident.Severity, _ *int64, _ string) (*incident.Incident, error) {
			gotTitle, gotTopicID = title, topicID
			return &incident.Incident{Title: title, TopicID: topicID, Severity: incident.SeverityMedium, Status: incident.StatusActive}, nil
		},
	}, api)
	ctx := &mockContext{args: []string{"create", "db", "is", "down"}}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotTitle != "db is down" {
		t.Errorf("service got title %q, want %q", gotTitle, "db is down")
	}
	if gotTopicID != 777 {
		t.Errorf("service got topic id %d, want 777", gotTopicID)
	}
	if api.createdTopic.Name != "db is down" {
		t.Errorf("topic name %q, want %q", api.createdTopic.Name, "db is down")
	}
	if len(api.sent) != 1 || api.sent[0].threadID != 777 {
		t.Fatalf("expected card sent to topic 777, got %v", api.sent)
	}
	if !strings.Contains(api.sent[0].what, "db is down") {
		t.Errorf("card %q does not contain the title", api.sent[0].what)
	}
}

func TestHandleIncidentCreateReportsForumRequiredOnTopicError(t *testing.T) {
	api := newFakeAPI()
	api.createErr = errs.New(errs.KindUnavailable, "tg", "not a forum")
	h := New(&fakeService{
		create: func(int64, int64, string, incident.Severity, *int64, string) (*incident.Incident, error) {
			t.Fatal("service should not be called when the topic cannot be created")
			return nil, nil
		},
	}, api)
	ctx := &mockContext{args: []string{"create", "outage"}}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "forum supergroup")
}

func TestHandleIncidentCreateDeletesTopicOnConflict(t *testing.T) {
	api := newFakeAPI()
	api.createdTopic.ThreadID = 42
	h := New(&fakeService{
		create: func(int64, int64, string, incident.Severity, *int64, string) (*incident.Incident, error) {
			return nil, errs.ErrIncidentAlreadyActive
		},
	}, api)
	ctx := &mockContext{args: []string{"create", "again"}}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "already active")
	if len(api.deleted) != 1 || api.deleted[0] != 42 {
		t.Fatalf("expected orphan topic 42 to be deleted, got %v", api.deleted)
	}
}

func TestHandleIncidentMessageAddsTimelineUpdate(t *testing.T) {
	var gotMsg string
	var gotTopicID int64
	h := New(&fakeService{
		addEvent: func(_ int64, topicID int64, _ *int64, _, message string) (*event.Event, error) {
			gotMsg, gotTopicID = message, topicID
			return &event.Event{}, nil
		},
	}, newFakeAPI())
	ctx := &mockContext{args: []string{"db", "is", "down"}, threadID: 9}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMsg != "db is down" {
		t.Errorf("service got message %q, want %q", gotMsg, "db is down")
	}
	if gotTopicID != 9 {
		t.Errorf("service got topic id %d, want 9", gotTopicID)
	}
	sentContains(t, ctx, "Update added to the timeline")
}

func TestHandleIncidentCloseSendsSummaryAndReportToGeneralAndDeletesTopic(t *testing.T) {
	now := time.Now()
	api := newFakeAPI()
	h := New(&fakeService{
		report: func(int64, int64) ([]byte, error) { return []byte("%PDF-1.4 fake"), nil },
		closeInc: func(int64, int64, *int64, string) (*incident.Incident, error) {
			return &incident.Incident{Title: "outage", CreatedAt: now, ClosedAt: &now}, nil
		},
	}, api)
	ctx := &mockContext{args: []string{"close"}, threadID: 333}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	apiSentContains(t, api, "Incident closed")
	apiSentContains(t, api, "telebot.Document")
	for _, s := range api.sent {
		if s.threadID != 0 {
			t.Errorf("close output sent to thread %d, want General (0)", s.threadID)
		}
	}
	if len(api.deleted) != 1 || api.deleted[0] != 333 {
		t.Fatalf("expected topic 333 to be deleted, got %v", api.deleted)
	}
}

func TestHandleIncidentCloseStillClosesWhenReportFails(t *testing.T) {
	now := time.Now()
	api := newFakeAPI()
	h := New(&fakeService{
		report: func(int64, int64) ([]byte, error) { return nil, errs.New(errs.KindUnavailable, "report", "down") },
		closeInc: func(int64, int64, *int64, string) (*incident.Incident, error) {
			return &incident.Incident{Title: "outage", CreatedAt: now, ClosedAt: &now}, nil
		},
	}, api)
	ctx := &mockContext{args: []string{"close"}, threadID: 1}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	apiSentContains(t, api, "Incident closed")
	apiSentContains(t, api, "report could not be generated")
	if len(api.deleted) != 1 {
		t.Fatalf("expected the topic to be deleted, got %v", api.deleted)
	}
}

func TestHandleIncidentCloseReportsNoActive(t *testing.T) {
	api := newFakeAPI()
	h := New(&fakeService{
		report:   func(int64, int64) ([]byte, error) { return nil, errs.ErrNoActiveIncident },
		closeInc: func(int64, int64, *int64, string) (*incident.Incident, error) { return nil, errs.ErrNoActiveIncident },
	}, api)
	ctx := &mockContext{args: []string{"close"}, threadID: 5}

	if err := h.HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "no active incident")
	if len(api.deleted) != 0 {
		t.Fatalf("expected no topic deletion when there is no active incident, got %v", api.deleted)
	}
}
