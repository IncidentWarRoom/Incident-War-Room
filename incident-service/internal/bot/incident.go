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

// topicNameLimit is Telegram's maximum forum topic name length (in characters).
const topicNameLimit = 128

// topicForumRequired is shown when a forum topic cannot be created, which almost
// always means the chat is not a forum supergroup or the bot lacks the Manage
// Topics admin right.
const topicForumRequired = "Couldn't open a topic for this incident. Use /incident create in a forum supergroup where the bot is an admin with the \"Manage Topics\" right."

// topicName trims the incident title to Telegram's topic-name limit.
func topicName(title string) string {
	r := []rune(title)
	if len(r) > topicNameLimit {
		return string(r[:topicNameLimit])
	}
	return title
}

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

func (h *Handler) createIncident(c telebot.Context, description string) error {
	if description == "" {
		return c.Send("Please add a description:\n/incident create <what happened>")
	}

	ctx, cancel := reqContext()
	defer cancel()

	chat := c.Chat()

	topic, err := h.api.CreateTopic(chat, &telebot.Topic{Name: topicName(description)})
	if err != nil {
		log.Printf("bot: create topic: %v", err)
		return c.Send(topicForumRequired)
	}

	userID, username := sender(c)
	inc, err := h.svc.CreateIncident(ctx, chat.ID, int64(topic.ThreadID), description, "", userID, username)
	if err != nil {
		log.Printf("bot: create incident: %v", err)
		if delErr := h.api.DeleteTopic(chat, topic); delErr != nil {
			log.Printf("bot: delete orphan topic %d: %v", topic.ThreadID, delErr)
		}
		return c.Send(userError(err))
	}

	_, err = h.api.Send(
		chat,
		incidentCard(inc.Title, inc.Severity, inc.Status),
		incidentMenu(),
		&telebot.SendOptions{ThreadID: topic.ThreadID},
	)
	return err
}

func (h *Handler) addUpdate(c telebot.Context, message string) error {
	ctx, cancel := reqContext()
	defer cancel()

	userID, username := sender(c)
	if _, err := h.svc.AddTimelineEvent(ctx, c.Chat().ID, threadID(c), userID, username, message); err != nil {
		log.Printf("bot: add timeline event: %v", err)
		return c.Send(userError(err))
	}

	return c.Send("📝 Update added to the timeline.")
}

// closeIncident closes the incident bound to the current topic. The closing
// summary and the PDF report are posted to the chat's General thread (the topic
// is about to be removed), then the topic is deleted.
func (h *Handler) closeIncident(c telebot.Context) (*incident.Incident, error) {
	ctx, cancel := reqContext()
	defer cancel()

	chat := c.Chat()
	topicID := threadID(c)
	userID, username := sender(c)

	pdf, reportErr := h.svc.GenerateReport(ctx, chat.ID, topicID)

	inc, err := h.svc.CloseIncident(ctx, chat.ID, topicID, userID, username)
	if err != nil {
		log.Printf("bot: close incident: %v", err)
		return nil, c.Send(userError(err))
	}

	if _, err := h.api.Send(chat, response.IncidentClosed(*inc), telebot.ModeHTML); err != nil {
		return inc, err
	}

	if reportErr != nil {
		log.Printf("bot: generate report: %v", reportErr)
		if _, err := h.api.Send(chat, "⚠️ The incident was closed, but the report could not be generated right now."); err != nil {
			return inc, err
		}
	} else if _, err := h.api.Send(chat, &telebot.Document{
		File:     telebot.FromReader(bytes.NewReader(pdf)),
		FileName: "incident_report.pdf",
		MIME:     "application/pdf",
		Caption:  "📄 Incident report",
	}); err != nil {
		return inc, err
	}

	if topicID != 0 {
		if err := h.api.DeleteTopic(chat, &telebot.Topic{ThreadID: int(topicID)}); err != nil {
			log.Printf("bot: delete topic %d: %v", topicID, err)
		}
	}

	return inc, nil
}

func (h *Handler) setSeverity(c telebot.Context, sev incident.Severity) (*incident.Incident, error) {
	ctx, cancel := reqContext()
	defer cancel()

	return h.svc.SetSeverity(ctx, c.Chat().ID, threadID(c), sev)
}
