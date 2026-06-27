package bot

import (
	"testing"

	"gopkg.in/telebot.v3"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/cQu1x/Incident-War-Room/internal/domain/media"
	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

func activeTimeline() func(int64, int64) (*incident.Incident, []event.Event, error) {
	return func(int64, int64) (*incident.Incident, []event.Event, error) {
		return &incident.Incident{}, nil, nil
	}
}

func TestHandleTopicPhotoDisabled(t *testing.T) {
	called := false
	h := New(&fakeService{
		timeline: activeTimeline(),
		addImage: func(int64, int64, *int64, string, string, media.Image) (*event.Event, error) {
			called = true
			return &event.Event{}, nil
		},
	}, newFakeAPI())

	ctx := &mockContext{chatID: 42, message: &telebot.Message{ThreadID: 7, Photo: &telebot.Photo{}}}

	if err := h.HandleTopicPhoto(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "S3 storage is not connected")
	if called {
		t.Error("image should not be uploaded when media is disabled")
	}
}

func TestHandleTopicPhotoAlbumRejected(t *testing.T) {
	called := false
	h := New(&fakeService{
		timeline: activeTimeline(),
		addImage: func(int64, int64, *int64, string, string, media.Image) (*event.Event, error) {
			called = true
			return &event.Event{}, nil
		},
	}, newFakeAPI(), WithMediaEnabled(true))

	ctx := &mockContext{chatID: 42, message: &telebot.Message{ThreadID: 7, AlbumID: "grp1", Photo: &telebot.Photo{}}}

	if err := h.HandleTopicPhoto(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentContains(t, ctx, "Only one photo")
	if called {
		t.Error("album photos should be rejected before upload")
	}
}

func TestHandleTopicPhotoUploadsImage(t *testing.T) {
	var gotCaption string
	var gotImg media.Image
	h := New(&fakeService{
		timeline: activeTimeline(),
		addImage: func(_, _ int64, _ *int64, _, caption string, img media.Image) (*event.Event, error) {
			gotCaption = caption
			gotImg = img
			return &event.Event{}, nil
		},
	}, newFakeAPI(), WithMediaEnabled(true))

	ctx := &mockContext{chatID: 42, message: &telebot.Message{
		ThreadID: 7,
		Caption:  "prod down",
		Photo:    &telebot.Photo{File: telebot.File{FileID: "abc"}},
	}}

	if err := h.HandleTopicPhoto(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotCaption != "prod down" {
		t.Errorf("caption = %q, want %q", gotCaption, "prod down")
	}
	if len(gotImg.Data) == 0 || gotImg.ContentType != "image/jpeg" || gotImg.Ext != "jpg" {
		t.Errorf("unexpected image: %+v", gotImg)
	}
	if len(ctx.sent) != 0 {
		t.Errorf("expected no reply on success, got %v", ctx.sent)
	}
}

func TestHandleTopicPhotoNoActiveIncidentSilent(t *testing.T) {
	h := New(&fakeService{
		timeline: func(int64, int64) (*incident.Incident, []event.Event, error) {
			return nil, nil, errs.ErrNoActiveIncident
		},
	}, newFakeAPI(), WithMediaEnabled(true))

	ctx := &mockContext{chatID: 42, message: &telebot.Message{ThreadID: 7, Photo: &telebot.Photo{}}}

	if err := h.HandleTopicPhoto(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ctx.sent) != 0 {
		t.Errorf("expected no reply without an active incident, got %v", ctx.sent)
	}
}
