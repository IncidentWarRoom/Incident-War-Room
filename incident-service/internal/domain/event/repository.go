package event

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// Create inserts a new event and fills in the generated ID and CreatedAt.
	Create(ctx context.Context, e *Event) error

	// ListByIncidentID returns all events of an incident in chronological order.
	ListByIncidentID(ctx context.Context, incidentID uuid.UUID) ([]Event, error)

	// ListParticipants returns distinct author IDs of an incident's events.
	ListParticipants(ctx context.Context, incidentID uuid.UUID) ([]int64, error)
}
