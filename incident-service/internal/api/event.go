package api

import (
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
)

type eventResponse struct {
	ID         string    `json:"id"`
	IncidentID string    `json:"incidentId"`
	Type       string    `json:"type"`
	UserID     *int64    `json:"userId"`
	Username   string    `json:"username"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"createdAt"`
}

func newEventResponses(events []event.Event) []eventResponse {
	out := make([]eventResponse, 0, len(events))
	for _, e := range events {
		out = append(out, eventResponse{
			ID:         e.ID.String(),
			IncidentID: e.IncidentID.String(),
			Type:       string(e.Type),
			UserID:     e.UserID,
			Username:   e.Username,
			Message:    e.Message,
			CreatedAt:  e.CreatedAt,
		})
	}
	return out
}
