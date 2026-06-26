package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
)

// IncidentImages returns the timeline events of the incident that carry an
// attached image, in chronological order. Returns errs.ErrIncidentNotFound if
// the incident does not exist.
func (s *Service) IncidentImages(ctx context.Context, id uuid.UUID) ([]event.Event, error) {
	if _, err := s.incidents.GetByID(ctx, id); err != nil {
		return nil, err
	}

	events, err := s.events.ListByIncidentID(ctx, id)
	if err != nil {
		return nil, err
	}

	images := make([]event.Event, 0)
	for _, e := range events {
		if e.MediaURL != nil && *e.MediaURL != "" {
			images = append(images, e)
		}
	}

	return images, nil
}
