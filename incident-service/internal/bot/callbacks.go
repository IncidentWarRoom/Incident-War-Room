package bot

import (
	"gopkg.in/telebot.v3"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

// Callback handlers for the inline incident panel.
//
// Each handler first answers the callback (c.Respond) to stop Telegram's
// loading spinner, then updates the message in place (c.Edit) so the chat
// stays tidy. Until persistence lands, the incident state is recovered from
// the card text (parseCard) so the description survives every action.

func handleShowTimeline(c telebot.Context) error {
	if err := c.Respond(); err != nil {
		return err
	}
	return c.Send("[stub] Incident timeline is empty. (not implemented yet)")
}

func handleCloseIncident(c telebot.Context) error {
	if err := c.Respond(&telebot.CallbackResponse{Text: "Incident closed ✅"}); err != nil {
		return err
	}
	description, sev := parseCard(c.Message().Text)
	card := incidentCard(description, sev, incident.StatusClosed)
	// Empty markup removes the buttons — a closed incident has no actions.
	return c.Edit(card, &telebot.ReplyMarkup{})
}

func handleChangeSeverity(c telebot.Context) error {
	if err := c.Respond(); err != nil {
		return err
	}
	// Swap only the keyboard so the card text (and its description) stays put.
	return c.Edit(severityMenu)
}

func handleSetSeverity(c telebot.Context) error {
	sev := incident.Severity(c.Data())
	if err := c.Respond(&telebot.CallbackResponse{Text: "Severity set to " + string(sev)}); err != nil {
		return err
	}
	description, _ := parseCard(c.Message().Text)
	card := incidentCard(description, sev, incident.StatusActive)
	return c.Edit(card, incidentMenu)
}

func handleSeverityBack(c telebot.Context) error {
	if err := c.Respond(); err != nil {
		return err
	}
	// Restore the action panel; card text is untouched.
	return c.Edit(incidentMenu)
}
