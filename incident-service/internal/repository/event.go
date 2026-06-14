package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

// EventRepository stores timeline events in the "incident_events" table.
type EventRepository struct {
	db Querier
}

func NewEventRepository(db Querier) *EventRepository {
	return &EventRepository{db: db}
}

// Create inserts a new timeline event. ID and CreatedAt are generated
// by the database and written back into e.
func (r *EventRepository) Create(ctx context.Context, e *event.Event) error {
	const query = `
		INSERT INTO incident_events (incident_id, type, author_id, username, message)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	err := r.db.
		QueryRow(ctx, query, e.IncidentID, e.Type, e.AuthorID, e.Username, e.Message).
		Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		return errs.Wrapf(errs.KindInternal, "repository.Event.Create", err, "insert event")
	}

	return nil
}

// ListByIncidentID returns all events of an incident in chronological
// order — this is the incident timeline.
func (r *EventRepository) ListByIncidentID(ctx context.Context, incidentID uuid.UUID) ([]event.Event, error) {
	const query = `
		SELECT id, incident_id, type, author_id, username, message, created_at
		FROM incident_events
		WHERE incident_id = $1
		ORDER BY created_at ASC`

	rows, err := r.db.Query(ctx, query, incidentID)
	if err != nil {
		return nil, errs.Wrapf(errs.KindInternal, "repository.Event.ListByIncidentID", err, "select events")
	}

	events, err := pgx.CollectRows(rows, scanEvent)
	if err != nil {
		return nil, errs.Wrapf(errs.KindInternal, "repository.Event.ListByIncidentID", err, "scan events")
	}

	return events, nil
}

// ListParticipants returns distinct Telegram user IDs of everyone who
// produced at least one event in the incident. Events without an author
// (system events) are skipped.
func (r *EventRepository) ListParticipants(ctx context.Context, incidentID uuid.UUID) ([]int64, error) {
	const query = `
		SELECT DISTINCT author_id
		FROM incident_events
		WHERE incident_id = $1 AND author_id IS NOT NULL`

	rows, err := r.db.Query(ctx, query, incidentID)
	if err != nil {
		return nil, errs.Wrapf(errs.KindInternal, "repository.Event.ListParticipants", err, "select participants")
	}

	participants, err := pgx.CollectRows(rows, pgx.RowTo[int64])
	if err != nil {
		return nil, errs.Wrapf(errs.KindInternal, "repository.Event.ListParticipants", err, "scan participants")
	}

	return participants, nil
}

func scanEvent(row pgx.CollectableRow) (event.Event, error) {
	var e event.Event
	err := row.Scan(
		&e.ID,
		&e.IncidentID,
		&e.Type,
		&e.AuthorID,
		&e.Username,
		&e.Message,
		&e.CreatedAt,
	)
	return e, err
}
