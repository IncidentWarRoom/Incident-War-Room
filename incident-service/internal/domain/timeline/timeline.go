package timeline

import (
	"context"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

type Timeline struct {
	Incident incident.Incident
	Events   []event.Event
}

type Publisher interface {
	Publish(ctx context.Context, t Timeline) ([]string, error)
}
