package bot

import (
	"github.com/cQu1x/Incident-War-Room/internal/bot/response"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"gopkg.in/telebot.v3"
)

func HandleTimeline(c telebot.Context) error {
	// Persistence is not wired yet; show the empty timeline for a placeholder
	// incident so the formatted response can be exercised end to end.
	inc := incident.Incident{Title: "Untitled incident"}
	return c.Send(response.Timeline(inc, nil), telebot.ModeHTML)
}
