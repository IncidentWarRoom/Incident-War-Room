package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

// GetIncident returns the incident with the given ID, or
// errs.ErrIncidentNotFound if it does not exist.
func (s *Service) GetIncident(ctx context.Context, id uuid.UUID) (*incident.Incident, error) {
	return s.incidents.GetByID(ctx, id)
}
