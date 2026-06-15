package event

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID         uuid.UUID
	IncidentID uuid.UUID
	Type       EventType
	AuthorID   *int64
	Username   string
	Message    string
	CreatedAt  time.Time
}

type EventType string

const (
	TypeIncidentCreated EventType = "INCIDENT_CREATED"
	TypeCommentAdded    EventType = "COMMENT_ADDED"
	TypeIncidentClosed  EventType = "INCIDENT_CLOSED"
)
