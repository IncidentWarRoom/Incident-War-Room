package api

import (
	"strings"
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

// timelineResponse is the incident timeline as the frontend needs it: the
// incident header (title, status, severity), its time span and duration, the
// list of responders, and the chronological events.
type timelineResponse struct {
	IncidentID      string          `json:"incidentId"`
	Title           string          `json:"title"`
	Status          string          `json:"status"`
	Severity        string          `json:"severity"`
	StartedAt       time.Time       `json:"startedAt"`
	EndedAt         *time.Time      `json:"endedAt"`
	DurationSeconds *int64          `json:"durationSeconds"`
	Responders      []string        `json:"responders"`
	Events          []eventResponse `json:"events"`
}

func newTimelineResponse(inc incident.Incident, events []event.Event) timelineResponse {
	resp := timelineResponse{
		IncidentID: inc.ID.String(),
		Title:      inc.Title,
		Status:     string(inc.Status),
		Severity:   string(inc.Severity),
		StartedAt:  inc.CreatedAt,
		EndedAt:    inc.ClosedAt,
		Responders: responders(events),
		Events:     newEventResponses(events),
	}

	if inc.ClosedAt != nil {
		seconds := int64(inc.ClosedAt.Sub(inc.CreatedAt).Seconds())
		resp.DurationSeconds = &seconds
	}

	return resp
}

// responders returns the distinct usernames that appear on the timeline, in
// first-seen order. Events without a username (system events) are skipped.
func responders(events []event.Event) []string {
	seen := make(map[string]struct{})
	out := make([]string, 0)
	for _, e := range events {
		name := strings.TrimSpace(e.Username)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	return out
}
