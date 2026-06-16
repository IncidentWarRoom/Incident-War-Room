package bot

import "gopkg.in/telebot.v3"

// Inline buttons and the menus that arrange them.
//
// Buttons are package-level so the callback router (Register) and the menu
// builders can reference the same Btn values — telebot routes a callback by
// the button's Unique string.
//
// The menus, in contrast, are built fresh on every send by incidentMenu and
// severityMenu. telebot's processButtons mutates a markup's callback_data in
// place before each send, prepending "\f<unique>|"; reusing a single shared
// *ReplyMarkup makes that prefix accumulate ("\fset_severity|\fset_severity|…")
// and corrupts the callback payload. Handing Send/Edit a new markup each time
// keeps the payload clean.
var (
	btnTimeline = telebot.Btn{Unique: "show_timeline", Text: "📜 Show Timeline"}
	btnClose    = telebot.Btn{Unique: "close_incident", Text: "✅ Close Incident"}
	btnSeverity = telebot.Btn{Unique: "change_severity", Text: "⚠️ Change Severity"}

	// All three severity buttons share the "set_severity" unique, so a single
	// handler serves them; the chosen level travels as the callback payload
	// (read via telebot.Context.Data).
	btnSevLow    = telebot.Btn{Unique: "set_severity", Text: "🟢 Low", Data: "LOW"}
	btnSevMedium = telebot.Btn{Unique: "set_severity", Text: "🟡 Medium", Data: "MEDIUM"}
	btnSevHigh   = telebot.Btn{Unique: "set_severity", Text: "🔴 High", Data: "HIGH"}
	btnSevBack   = telebot.Btn{Unique: "severity_back", Text: "⬅️ Back"}
)

// incidentMenu is the main action panel attached to an incident message.
func incidentMenu() *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{}
	m.Inline(
		m.Row(btnTimeline),
		m.Row(btnClose, btnSeverity),
	)
	return m
}

// severityMenu is the pop-up shown after "Change Severity" is tapped.
func severityMenu() *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{}
	m.Inline(
		m.Row(btnSevLow, btnSevMedium, btnSevHigh),
		m.Row(btnSevBack),
	)
	return m
}
