package bot

import (
	"log"

	"gopkg.in/telebot.v3"

	"github.com/cQu1x/Incident-War-Room/internal/bot/response"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

func (h *Handler) handleShowTimeline(c telebot.Context) error {
	if err := c.Respond(); err != nil {
		return err
	}

	ctx, cancel := reqContext()
	defer cancel()

	inc, events, err := h.svc.GetTimeline(ctx, c.Chat().ID)
	if err != nil {
		log.Printf("bot: show timeline: %v", err)
		return c.Send(userError(err))
	}

	return c.Send(response.Timeline(*inc, events), telebot.ModeHTML)
}

func (h *Handler) handleCloseIncident(c telebot.Context) error {
	if err := c.Respond(&telebot.CallbackResponse{Text: "Closing incident…"}); err != nil {
		return err
	}

	inc, err := h.closeIncident(c)
	if err != nil {
		return err
	}
	if inc == nil {
		return nil
	}

	return c.Edit(incidentCard(inc.Title, inc.Severity, incident.StatusClosed), &telebot.ReplyMarkup{})
}

func (h *Handler) handleChangeSeverity(c telebot.Context) error {
	if err := c.Respond(); err != nil {
		return err
	}

	return c.Edit(severityMenu())
}

func (h *Handler) handleSetSeverity(c telebot.Context) error {
	sev := incident.Severity(c.Data())
	if err := c.Respond(&telebot.CallbackResponse{Text: "Severity set to " + string(sev)}); err != nil {
		return err
	}

	inc, err := h.setSeverity(c, sev)
	if err != nil {
		log.Printf("bot: set severity: %v", err)
		return c.Send(userError(err))
	}

	return c.Edit(incidentCard(inc.Title, inc.Severity, inc.Status), incidentMenu())
}

func (h *Handler) handleSeverityBack(c telebot.Context) error {
	if err := c.Respond(); err != nil {
		return err
	}

	return c.Edit(incidentMenu())
}
