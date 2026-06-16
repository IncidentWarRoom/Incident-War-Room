package bot

import (
	"bytes"
	"log"
	"strings"

	"gopkg.in/telebot.v3"

	"github.com/cQu1x/Incident-War-Room/internal/bot/response"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

const incidentUsage = "Usage:\n/incident create <description> — open a new incident\n/incident close — close the active incident\n/incident <message> — add an update to the timeline"

func (h *Handler) HandleIncident(c telebot.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Send(incidentUsage)
	}

	switch args[0] {
	case "create":
		return h.createIncident(c, strings.TrimSpace(strings.Join(args[1:], " ")))
	case "close":
		_, err := h.closeIncident(c)
		return err
	default:
		return h.addUpdate(c, strings.TrimSpace(strings.Join(args, " ")))
	}
}

// createIncident opens a new incident and replies with the interactive card.
func (h *Handler) createIncident(c telebot.Context, description string) error {
	if description == "" {
		return c.Send("Please add a description:\n/incident create <what happened>")
	}

	ctx, cancel := reqContext()
	defer cancel()

	userID, username := sender(c)
	inc, err := h.svc.CreateIncident(ctx, c.Chat().ID, description, "", userID, username)
	if err != nil {
		log.Printf("bot: create incident: %v", err)
		return c.Send(userError(err))
	}

	return c.Send(incidentCard(inc.Title, inc.Severity, inc.Status), incidentMenu)
}

// addUpdate appends a comment to the active incident's timeline.
func (h *Handler) addUpdate(c telebot.Context, message string) error {
	ctx, cancel := reqContext()
	defer cancel()

	userID, username := sender(c)
	if _, err := h.svc.AddTimelineEvent(ctx, c.Chat().ID, userID, username, message); err != nil {
		log.Printf("bot: add timeline event: %v", err)
		return c.Send(userError(err))
	}

	return c.Send("📝 Update added to the timeline.")
}

// closeIncident closes the active incident, generates its report and sends both
// the closing summary and (best effort) the PDF. Report rendering goes first so
// it still sees an active incident; a report failure does not block the close.
//
// On a logical failure (e.g. no active incident) it replies with a friendly
// message and returns (nil, nil); the returned incident is non-nil only when
// the incident was actually closed.
func (h *Handler) closeIncident(c telebot.Context) (*incident.Incident, error) {
	ctx, cancel := reqContext()
	defer cancel()

	chatID := c.Chat().ID
	userID, username := sender(c)

	pdf, reportErr := h.svc.GenerateReport(ctx, chatID)

	inc, err := h.svc.CloseIncident(ctx, chatID, userID, username)
	if err != nil {
		log.Printf("bot: close incident: %v", err)
		return nil, c.Send(userError(err))
	}

	if err := c.Send(response.IncidentClosed(*inc), telebot.ModeHTML); err != nil {
		return inc, err
	}

	if reportErr != nil {
		log.Printf("bot: generate report: %v", reportErr)
		return inc, c.Send("⚠️ The incident was closed, but the report could not be generated right now.")
	}

	return inc, c.Send(&telebot.Document{
		File:     telebot.FromReader(bytes.NewReader(pdf)),
		FileName: "incident_report.pdf",
		MIME:     "application/pdf",
		Caption:  "📄 Incident report",
	})
}

// setSeverity changes the active incident's severity and returns the updated
// incident for re-rendering.
func (h *Handler) setSeverity(c telebot.Context, sev incident.Severity) (*incident.Incident, error) {
	ctx, cancel := reqContext()
	defer cancel()

	return h.svc.SetSeverity(ctx, c.Chat().ID, sev)
}
