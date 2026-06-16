package reportclient

import (
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/report"
)

type request struct {
	Incident     incidentDTO      `json:"incident"`
	Participants []participantDTO `json:"participants"`
	Timeline     []timelineDTO    `json:"timeline"`
}

type incidentDTO struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	CreatedAt time.Time  `json:"createdAt"`
	ClosedAt  *time.Time `json:"closedAt,omitempty"`
	Severity  string     `json:"severity,omitempty"`
	Status    string     `json:"status,omitempty"`
}

type participantDTO struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
}

type timelineDTO struct {
	Timestamp time.Time `json:"timestamp"`
	Username  string    `json:"username"`
	Message   string    `json:"message"`
}

// toRequest maps the domain report aggregate onto the wire contract.
func toRequest(r report.Report) request {
	req := request{
		Incident: incidentDTO{
			ID:        r.Incident.ID.String(),
			Title:     r.Incident.Title,
			CreatedAt: r.Incident.CreatedAt,
			ClosedAt:  r.Incident.ClosedAt,
			Severity:  string(r.Incident.Severity),
			Status:    string(r.Incident.Status),
		},
		Participants: make([]participantDTO, 0, len(r.Participants)),
		Timeline:     make([]timelineDTO, 0, len(r.Timeline)),
	}

	for _, p := range r.Participants {
		req.Participants = append(req.Participants, participantDTO{
			UserID:   p.UserID,
			Username: p.Username,
		})
	}

	for _, e := range r.Timeline {
		req.Timeline = append(req.Timeline, timelineDTO{
			Timestamp: e.CreatedAt,
			Username:  e.Username,
			Message:   e.Message,
		})
	}

	return req
}
