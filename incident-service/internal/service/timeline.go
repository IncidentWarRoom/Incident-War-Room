package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/cQu1x/Incident-War-Room/internal/domain/timeline"
)

// IncidentTimeline returns the incident with the given ID together with its
// events in chronological order. Returns errs.ErrIncidentNotFound if the
// incident does not exist.
func (s *Service) IncidentTimeline(ctx context.Context, id uuid.UUID) (*incident.Incident, []event.Event, error) {
	inc, err := s.incidents.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	events, err := s.events.ListByIncidentID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return inc, events, nil
}

// GetTimeline returns the chat's active incident together with its events in
// chronological order. Returns errs.ErrNoActiveIncident if the chat has no
// active incident.
func (s *Service) GetTimeline(ctx context.Context, chatID, topicID int64) (*incident.Incident, []event.Event, error) {
	active, err := s.incidents.GetActiveByTopicID(ctx, chatID, topicID)
	if err != nil {
		return nil, nil, err
	}

	events, err := s.events.ListByIncidentID(ctx, active.ID)
	if err != nil {
		return nil, nil, err
	}

	return active, events, nil
}

func (s *Service) PublishTimeline(ctx context.Context, chatID, topicID int64) ([]string, error) {
	inc, events, err := s.GetTimeline(ctx, chatID, topicID)
	if err != nil {
		return nil, err
	}

	urls, err := s.timelines.Publish(ctx, timeline.Timeline{Incident: *inc, Events: events})
	if err != nil {
		return nil, err
	}

	if err := s.incidents.UpdateTelegraphURLs(ctx, inc.ID, urls); err != nil {
		return nil, err
	}
	inc.TelegraphURLs = urls

	return urls, nil
}
