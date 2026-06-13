package bot

import (
	"strings"

	"gopkg.in/telebot.v3"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

const incidentUsage = "Usage:\n/incident create <description> — open a new incident\n/incident close — close the active incident\n/incident <message> — add an update to the timeline"

func HandleIncident(c telebot.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Send(incidentUsage)
	}

	switch args[0] {
	case "create":
		description := strings.TrimSpace(strings.Join(args[1:], " "))
		if description == "" {
			return c.Send("Please add a description:\n/incident create <what happened>")
		}
		card := incidentCard(description, incident.SeverityMedium, incident.StatusActive)
		return c.Send(card, incidentMenu)
	case "close":
		return c.Send("[stub] Incident closed. (not implemented yet)")
	default:
		message := strings.Join(args, " ")
		return c.Send("[stub] Update added to timeline: " + message + " (not implemented yet)")
	}
}
