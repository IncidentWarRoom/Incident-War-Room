package bot

import (
	"log"

	"gopkg.in/telebot.v3"

	"github.com/cQu1x/Incident-War-Room/internal/bot/response"
)

func (h *Handler) HandleTimeline(c telebot.Context) error {
	ctx, cancel := reqContext()
	defer cancel()

	chatID, topicID := c.Chat().ID, threadID(c)

	inc, events, err := h.svc.GetTimeline(ctx, chatID, topicID)
	if err != nil {
		log.Printf("bot: get timeline: %v", err)
		return c.Send(userError(err))
	}

	msg := response.Timeline(*inc, events)
	if len(events) > 0 {
		urls, err := h.svc.PublishTimeline(ctx, chatID, topicID)
		if err != nil {
			log.Printf("bot: publish timeline: %v", err)
			msg += response.TimelineUnavailable()
		} else {
			msg += response.TimelineLink(urls)
		}
	}

	return c.Send(msg, &telebot.SendOptions{
		ThreadID:  int(topicID),
		ParseMode: telebot.ModeHTML,
	})
}
