// Package timeline defines the port for publishing an incident's full timeline
// to an external, shareable destination (a set of Telegraph pages). The
// concrete publisher (the Telegraph HTTP client) lives in the infrastructure
// layer and implements Publisher, the same way repositories implement the
// domain repository ports.
package timeline

import (
	"context"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

// Timeline is everything needed to publish an incident's timeline: the
// incident itself and its events in chronological order.
type Timeline struct {
	Incident incident.Incident
	Events   []event.Event
}

// Publisher renders a Timeline into one or more shareable pages and returns
// their URLs in reading order. A long timeline may be split across several
// pages, hence the slice.
//
// It is implemented by the infrastructure layer; the service depends only on
// this abstraction. A publisher that is unreachable or fails yields an
// errs.KindUnavailable error.
type Publisher interface {
	Publish(ctx context.Context, t Timeline) ([]string, error)
}
