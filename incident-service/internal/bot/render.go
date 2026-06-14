package bot

import (
	"fmt"
	"strings"

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

// parseCard recovers the description and severity from a rendered card.
//
// There is no persistence yet, so callbacks re-read the incident state from
// the message they are attached to in order to keep the description visible
// after an action. Missing fields fall back to sensible defaults.
func parseCard(text string) (description string, sev incident.Severity) {
	sev = incident.SeverityMedium
	for _, line := range strings.Split(text, "\n") {
		switch {
		case strings.HasPrefix(line, "📝 "):
			description = strings.TrimSpace(strings.TrimPrefix(line, "📝 "))
		case strings.Contains(line, "Severity: "):
			_, val, _ := strings.Cut(line, "Severity: ")
			sev = incident.Severity(strings.TrimSpace(val))
		}
	}
	return description, sev
}
