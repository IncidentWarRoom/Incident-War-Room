// Package response builds the typed, user-facing messages the bot sends to
// Telegram. Builders take domain models and return strings formatted with
// Telegram HTML markup, so callers must send them with telebot.ModeHTML.
package response

import (
	"html"
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/google/uuid"
)

// timeLayout is the human-readable timestamp format used across all responses.
const timeLayout = "2006-01-02 15:04 MST"

// escape makes user-supplied text safe to embed in an HTML-formatted message.
func escape(s string) string {
	return html.EscapeString(s)
}

// severityIcon returns a colored indicator for an incident severity. Unknown
// values fall back to the raw string without an icon.
func severityIcon(s incident.Severity) string {
	switch s {
	case incident.SeverityLow:
		return "🟢 LOW"
	case incident.SeverityMedium:
		return "🟠 MEDIUM"
	case incident.SeverityHigh:
		return "🔴 HIGH"
	default:
		return escape(string(s))
	}
}

// shortID renders the first segment of a UUID, enough to identify an incident
// in chat without overwhelming the message.
func shortID(id uuid.UUID) string {
	return id.String()[:8]
}

// formatTime renders a timestamp in the shared layout.
func formatTime(t time.Time) string {
	return t.Format(timeLayout)
}
