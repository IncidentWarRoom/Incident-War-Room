package incident

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	// Create inserts a new incident and fills in the generated ID and CreatedAt.
	// Returns errs.ErrIncidentAlreadyActive if the chat already has an active incident.
	Create(ctx context.Context, inc *Incident) error

	// GetByID returns errs.ErrIncidentNotFound if no incident exists with the given ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Incident, error)

	// GetActiveByChatID returns the chat's active incident,
	// or errs.ErrNoActiveIncident if there is none.
	GetActiveByChatID(ctx context.Context, chatID int64) (*Incident, error)

	GetActiveByTopicID(ctx context.Context, chatID, topicID int64) (*Incident, error)

	// UpdateSeverity returns errs.ErrIncidentNotFound if no incident exists with the given ID.
	UpdateSeverity(ctx context.Context, id uuid.UUID, severity Severity) error

	UpdateTopicID(ctx context.Context, id uuid.UUID, topicID int64) error
	UpdateReport(ctx context.Context, id uuid.UUID, telegraphURLs []string, reportURL string) error
	UpdateReportURL(ctx context.Context, id uuid.UUID, reportURL string) error

	// Close marks an active incident as closed.
	// Returns errs.ErrIncidentNotFound if the incident does not exist,
	// or errs.ErrIncidentAlreadyClosed if it is already closed.
	Close(ctx context.Context, id uuid.UUID, closedAt time.Time) error
}
