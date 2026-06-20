package bot

import "gopkg.in/telebot.v3"

// Register binds all command and inline-panel handlers to b.
func (h *Handler) Register(b *telebot.Bot) {
	b.Handle("/start", h.HandleStart)
	b.Handle("/incident", h.HandleIncident)
	b.Handle("/timeline", h.HandleTimeline)

	b.Handle(&btnTimeline, h.handleShowTimeline)
	b.Handle(&btnClose, h.handleCloseIncident)
	b.Handle(&btnSeverity, h.handleChangeSeverity)
	b.Handle(&btnSevBack, h.handleSeverityBack)

	b.Handle(&btnSevLow, h.handleSetSeverity)

	b.Handle(telebot.OnText, h.HandleTopicText)

	// Media is not recorded on the timeline; reject it in incident topics.
	for _, ev := range []string{
		telebot.OnPhoto,
		telebot.OnVideo,
		telebot.OnVideoNote,
		telebot.OnDocument,
		telebot.OnVoice,
		telebot.OnAudio,
		telebot.OnAnimation,
		telebot.OnSticker,
	} {
		b.Handle(ev, h.HandleTopicMedia)
	}
}
