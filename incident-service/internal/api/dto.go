package api

import (
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

type incidentResponse struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Severity      string     `json:"severity"`
	Status        string     `json:"status"`
	ChatID        int64      `json:"chatId"`
	TopicID       int64      `json:"topicId"`
	CreatedBy     *int64     `json:"createdBy"`
	CreatedAt     time.Time  `json:"createdAt"`
	ClosedAt      *time.Time `json:"closedAt"`
	TelegraphURLs []string   `json:"telegraphUrls"`
	ReportURL     *string    `json:"reportUrl"`
}

func newIncidentResponse(inc incident.Incident) incidentResponse {
	urls := inc.TelegraphURLs
	if urls == nil {
		urls = []string{}
	}
	return incidentResponse{
		ID:            inc.ID.String(),
		Title:         inc.Title,
		Severity:      string(inc.Severity),
		Status:        string(inc.Status),
		ChatID:        inc.ChatID,
		TopicID:       inc.TopicID,
		CreatedBy:     inc.CreatedBy,
		CreatedAt:     inc.CreatedAt,
		ClosedAt:      inc.ClosedAt,
		TelegraphURLs: urls,
		ReportURL:     inc.ReportURL,
	}
}

func newIncidentResponses(incidents []incident.Incident) []incidentResponse {
	out := make([]incidentResponse, 0, len(incidents))
	for _, inc := range incidents {
		out = append(out, newIncidentResponse(inc))
	}
	return out
}
