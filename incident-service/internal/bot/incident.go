package bot

import (
	"strings"
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/bot/response"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/google/uuid"
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
		// Persistence is not wired yet; build a placeholder incident so the
		// formatted response can already be exercised end to end.
		inc := incident.Incident{
			ID:        uuid.New(),
			Title:     "Untitled incident",
			Severity:  incident.SeverityMedium,
			Status:    incident.StatusActive,
			CreatedAt: time.Now(),
		}
		return c.Send(response.IncidentCreated(inc), telebot.ModeHTML)
	case "close":
		// Persistence is not wired yet; build a placeholder closed incident so
		// the formatted response can already be exercised end to end.
		now := time.Now()
		inc := incident.Incident{
			ID:        uuid.New(),
			Title:     "Untitled incident",
			Severity:  incident.SeverityMedium,
			Status:    incident.StatusClosed,
			CreatedAt: now,
			ClosedAt:  &now,
		}
		return c.Send(response.IncidentClosed(inc), telebot.ModeHTML)
	default:
		message := strings.Join(args, " ")
		return c.Send("[stub] Update added to timeline: " + message + " (not implemented yet)")
	}
}
