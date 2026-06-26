package api

import (
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
)

// imageResponse is a single image attached to the incident timeline, with the
// caption it was posted with and its author.
type imageResponse struct {
	EventID   string    `json:"eventId"`
	URL       string    `json:"url"`
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

// newImageResponses maps timeline events to image DTOs. Every event is expected
// to carry a media URL; events without one are skipped defensively.
func newImageResponses(events []event.Event) []imageResponse {
	out := make([]imageResponse, 0, len(events))
	for _, e := range events {
		if e.MediaURL == nil {
			continue
		}
		out = append(out, imageResponse{
			EventID:   e.ID.String(),
			URL:       *e.MediaURL,
			Username:  e.Username,
			Message:   e.Message,
			CreatedAt: e.CreatedAt,
		})
	}
	return out
}
