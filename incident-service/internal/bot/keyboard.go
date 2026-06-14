package bot

import "gopkg.in/telebot.v3"

// Inline keyboards and their buttons.
//
// incidentMenu is the main action panel attached to an incident message.
// severityMenu is the pop-up shown after "Change Severity" is tapped.
//
// Buttons are package-level so that both the builders (Inline) and the
// callback router (Register) can reference the same Btn values — telebot
// routes a callback by the button's Unique string.
var (
	incidentMenu = &telebot.ReplyMarkup{}

	btnTimeline = incidentMenu.Data("📜 Show Timeline", "show_timeline")
	btnClose    = incidentMenu.Data("✅ Close Incident", "close_incident")
	btnSeverity = incidentMenu.Data("⚠️ Change Severity", "change_severity")

	severityMenu = &telebot.ReplyMarkup{}

	// All three severity buttons share the "set_severity" unique, so a single
	// handler serves them; the chosen level travels as the callback payload
	// (read via telebot.Context.Data).
	btnSevLow    = severityMenu.Data("🟢 Low", "set_severity", "LOW")
	btnSevMedium = severityMenu.Data("🟡 Medium", "set_severity", "MEDIUM")
	btnSevHigh   = severityMenu.Data("🔴 High", "set_severity", "HIGH")
	btnSevBack   = severityMenu.Data("⬅️ Back", "severity_back")
)

func init() {
	incidentMenu.Inline(
		incidentMenu.Row(btnTimeline),
		incidentMenu.Row(btnClose, btnSeverity),
	)

	severityMenu.Inline(
		severityMenu.Row(btnSevLow, btnSevMedium, btnSevHigh),
		severityMenu.Row(btnSevBack),
	)
}
