package bot

import (
	"strings"

	"gopkg.in/telebot.v3"
)

const incidentUsage = "Usage:\n/incident create — open a new incident\n/incident close — close the active incident\n/incident <message> — add an update to the timeline"

func HandleIncident(c telebot.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Send(incidentUsage)
	}

	switch args[0] {
	case "create":
		return c.Send("[stub] Incident created. (not implemented yet)")
	case "close":
		return c.Send("[stub] Incident closed. (not implemented yet)")
	default:
		message := strings.Join(args, " ")
		return c.Send("[stub] Update added to timeline: " + message + " (not implemented yet)")
	}
}
