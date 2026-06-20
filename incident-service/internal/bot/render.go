package bot

import (
	"fmt"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

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

func incidentCard(title string, sev incident.Severity, status incident.Status) string {
	return fmt.Sprintf(
		"🚨 Incident Investigation Started\n\n"+
			"📝 %s\n"+
			"%s Severity: %s\n"+
			"🔵 Status: %s\n\n"+
			"This topic is dedicated to the investigation of the current incident.\n\n"+
			"Rules:\n"+
			"• Every message sent in this topic will be recorded in the incident timeline.\n"+
			"• Only incident-related information should be posted here.\n"+
			"• Media messages are not supported in this version.\n"+
			"• This topic will be permanently deleted after the incident is closed.\n\n"+
			"Available commands:\n"+
			"• /timeline - see the timeline of the incident events in Telegraph pages format\n"+
			"• /incident close - close the incident",
		title, severityEmoji(sev), sev, status,
	)
}
