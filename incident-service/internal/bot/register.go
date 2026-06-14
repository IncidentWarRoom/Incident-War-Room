package bot

import "gopkg.in/telebot.v3"

func Register(b *telebot.Bot) {
	b.Handle("/start", HandleStart)
	b.Handle("/incident", HandleIncident)
	b.Handle("/timeline", HandleTimeline)

	// Inline panel callbacks.
	b.Handle(&btnTimeline, handleShowTimeline)
	b.Handle(&btnClose, handleCloseIncident)
	b.Handle(&btnSeverity, handleChangeSeverity)
	b.Handle(&btnSevBack, handleSeverityBack)

	// All severity buttons share one unique, so one registration routes them all.
	b.Handle(&btnSevLow, handleSetSeverity)
}
