package bot

import "gopkg.in/telebot.v3"

func Register(b *telebot.Bot) {
	b.Handle("/start", HandleStart)
	b.Handle("/incident", HandleIncident)
	b.Handle("/timeline", HandleTimeline)
}
