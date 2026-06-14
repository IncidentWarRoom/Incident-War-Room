// Package service implements the core business logic (use cases) of the
// incident war room on top of the domain repositories. It works purely with
// domain models and abstractions and has no knowledge of Postgres, telebot or
// any other infrastructure.
package service

import (
	"context"
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/cQu1x/Incident-War-Room/internal/domain/report"
)

// TxManager runs a unit of work inside a single transaction, exposing the
// incident and event repositories bound to it. It is implemented by the
// repository layer; the service depends only on this abstraction.
type TxManager interface {
	WithTx(ctx context.Context, fn func(incident.Repository, event.Repository) error) error
}

// Service holds the use cases of the incident war room.
//
// Read-only use cases call the repositories directly; use cases that write to
// more than one table (CreateIncident, CloseIncident) run through the
// TxManager so they commit atomically. GenerateReport delegates rendering to
// the report.Generator port.
type Service struct {
	incidents incident.Repository
	events    event.Repository
	tx        TxManager
	reports   report.Generator
	now       func() time.Time
}

func New(incidents incident.Repository, events event.Repository, tx TxManager, reports report.Generator) *Service {
	return &Service{
		incidents: incidents,
		events:    events,
		tx:        tx,
		reports:   reports,
		now:       func() time.Time { return time.Now().UTC() },
	}
}
