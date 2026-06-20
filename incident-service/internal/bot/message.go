package bot

import (
	"errors"
	"log"
	"strings"

	"gopkg.in/telebot.v3"

	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

// mediaNotAllowed is the reply sent when a user posts media in an incident
// topic. Only text is recorded on the timeline.
const mediaNotAllowed = "⚠️ Media messages are not allowed in incident topics. " +
	"Please describe the issue in text — only text messages are recorded on the timeline."

// HandleTopicText records a plain text message posted in an incident topic on
// the timeline. Media messages (screenshots, voice, etc.) are not handled here
// and never reach the timeline.
func (h *Handler) HandleTopicText(c telebot.Context) error {
	m := c.Message()
	if m == nil || m.Text == "" {
		return nil
	}
	return h.captureTopicMessage(c, m.Text)
}

// HandleTopicMedia rejects media messages (photos, video, documents, voice,
// stickers, …) posted in an incident topic. Media is never recorded on the
// timeline, so the sender is told it is not allowed. Topics without an active
// incident, and messages outside a topic, are ignored.
func (h *Handler) HandleTopicMedia(c telebot.Context) error {
	topicID := threadID(c)
	if topicID == 0 {
		return nil
	}

	ctx, cancel := reqContext()
	defer cancel()

	// Only police topics that have an active incident; leave other topics alone.
	if _, _, err := h.svc.GetTimeline(ctx, c.Chat().ID, topicID); err != nil {
		if !errors.Is(err, errs.ErrNoActiveIncident) {
			log.Printf("bot: reject topic media: %v", err)
		}
		return nil
	}

	return c.Send(mediaNotAllowed, &telebot.SendOptions{ThreadID: int(topicID)})
}

// captureTopicMessage silently appends message to the active incident timeline
// for the current topic. Messages outside a topic, or in a topic without an
// active incident, are ignored.
func (h *Handler) captureTopicMessage(c telebot.Context, message string) error {
	message = strings.TrimSpace(message)
	if message == "" {
		return nil
	}

	topicID := threadID(c)
	if topicID == 0 {
		return nil
	}

	ctx, cancel := reqContext()
	defer cancel()

	userID, username := sender(c)
	if _, err := h.svc.AddTimelineEvent(ctx, c.Chat().ID, topicID, userID, username, message); err != nil {
		if errors.Is(err, errs.ErrNoActiveIncident) {
			return nil
		}
		log.Printf("bot: capture topic message: %v", err)
	}

	return nil
}
