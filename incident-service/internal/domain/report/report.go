// Package report defines the incident report aggregate and the port for
// rendering it into a document. The concrete renderer (the report-service
// HTTP client) lives in the infrastructure layer and implements Generator,
// the same way repositories implement the domain repository ports.
package report

import (
	"context"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

// Participant is a person who took part in an incident.
type Participant struct {
	UserID   int64
	Username string
}

// Report is everything needed to render an incident report: the incident
// itself, who took part, and the full timeline of events in chronological
// order.
type Report struct {
	Incident     incident.Incident
	Participants []Participant
	Timeline     []event.Event
}

type Generator interface {
	Generate(ctx context.Context, r Report) (string, error)
}
