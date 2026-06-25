package service

import (
	"context"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

// ListIncidents returns all incidents ordered from newest to oldest.
func (s *Service) ListIncidents(ctx context.Context) ([]incident.Incident, error) {
	return s.incidents.List(ctx)
}
