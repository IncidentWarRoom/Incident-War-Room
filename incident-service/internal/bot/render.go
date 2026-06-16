package bot

import (
	"fmt"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

// severityEmoji maps a severity level to a coloured dot for the card/menus.
func severityEmoji(s incident.Severity) string {
	switch s {
	case incident.SeverityLow:
		return "🟢"
	case incident.SeverityMedium:
		return "🟡"
	case incident.SeverityHigh:
		return "🔴"
	default:
		return "⚪"
	}
}

// incidentCard renders the human-readable summary shown under the inline panel.
func incidentCard(title string, sev incident.Severity, status incident.Status) string {
	return fmt.Sprintf(
		"🚨 Incident\n\n"+
			"📝 %s\n"+
			"%s Severity: %s\n"+
			"🔵 Status: %s",
		title, severityEmoji(sev), sev, status,
	)
}
