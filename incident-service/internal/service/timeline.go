package service

import (
	"context"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/cQu1x/Incident-War-Room/internal/domain/timeline"
)

// GetActiveIncident returns the chat's active incident, or
// errs.ErrNoActiveIncident if there is none.
func (s *Service) GetActiveIncident(ctx context.Context, chatID int64) (*incident.Incident, error) {
	return s.incidents.GetActiveByChatID(ctx, chatID)
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

	return s.timelines.Publish(ctx, timeline.Timeline{Incident: *inc, Events: events})
}
