package service

import (
	"context"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

// SetSeverity changes the severity of the chat's active incident and returns
// the updated incident.
//
// Returns errs.ErrNoActiveIncident if the chat has no active incident, or an
// errs.KindValidation error if the severity is invalid.
func (s *Service) SetSeverity(ctx context.Context, chatID int64, severity incident.Severity) (*incident.Incident, error) {
	const op = "service.SetSeverity"

	if !validSeverity(severity) {
		return nil, errs.New(errs.KindValidation, op, "invalid severity")
	}

	active, err := s.incidents.GetActiveByChatID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if err := s.incidents.UpdateSeverity(ctx, active.ID, severity); err != nil {
		return nil, err
	}

	active.Severity = severity
	return active, nil
}
