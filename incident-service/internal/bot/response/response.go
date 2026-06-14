package response

import (
	"html"
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/google/uuid"
)

const timeLayout = "2006-01-02 15:04 MST"

func escape(s string) string {
	return html.EscapeString(s)
}

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

func shortID(id uuid.UUID) string {
	return id.String()[:8]
}

func formatTime(t time.Time) string {
	return t.Format(timeLayout)
}
