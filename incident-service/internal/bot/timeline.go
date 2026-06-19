package bot

import (
	"log"

	"gopkg.in/telebot.v3"

	"github.com/cQu1x/Incident-War-Room/internal/bot/response"
)

func (h *Handler) HandleTimeline(c telebot.Context) error {
	ctx, cancel := reqContext()
	defer cancel()

	inc, events, err := h.svc.GetTimeline(ctx, c.Chat().ID, threadID(c))
	if err != nil {
		log.Printf("bot: get timeline: %v", err)
		return c.Send(userError(err))
	}

	return c.Send(response.Timeline(*inc, events), telebot.ModeHTML)
}
