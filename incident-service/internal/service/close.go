package service

import (
	"context"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

// CloseIncident closes the chat's active incident and records an
// INCIDENT_CLOSED event on its timeline. Resolving the active incident,
// closing it and writing the event all happen in one transaction.
//
// Returns errs.ErrNoActiveIncident if the chat has no active incident.
func (s *Service) CloseIncident(
	ctx context.Context,
	chatID int64,
	userID *int64,
	username string,
) (*incident.Incident, error) {
	closedAt := s.now()

	var active *incident.Incident
	err := s.tx.WithTx(ctx, func(incidents incident.Repository, events event.Repository) error {
		inc, err := incidents.GetActiveByChatID(ctx, chatID)
		if err != nil {
			return err
		}

		if err := incidents.Close(ctx, inc.ID, closedAt); err != nil {
			return err
		}

		if err := events.Create(ctx, &event.Event{
			IncidentID: inc.ID,
			Type:       event.TypeIncidentClosed,
			UserID:     userID,
			Username:   username,
		}); err != nil {
			return err
		}

		inc.Status = incident.StatusClosed
		inc.ClosedAt = &closedAt
		active = inc
		return nil
	})
	if err != nil {
		return nil, err
	}

	return active, nil
}
