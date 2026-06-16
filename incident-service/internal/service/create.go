package service

import (
	"context"
	"strings"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

// CreateIncident opens a new incident in the chat and records an
// INCIDENT_CREATED event on its timeline. Both writes happen in one
// transaction.
//
// Returns errs.ErrIncidentAlreadyActive if the chat already has an active
// incident, or an errs.KindValidation error if the input is invalid. An empty
// severity defaults to incident.SeverityMedium.
func (s *Service) CreateIncident(
	ctx context.Context,
	chatID int64,
	title string,
	severity incident.Severity,
	userID *int64,
	username string,
) (*incident.Incident, error) {
	const op = "service.CreateIncident"

	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errs.New(errs.KindValidation, op, "incident title is required")
	}

	if severity == "" {
		severity = incident.SeverityMedium
	}
	if !validSeverity(severity) {
		return nil, errs.New(errs.KindValidation, op, "invalid severity")
	}

	inc := &incident.Incident{
		Title:     title,
		Severity:  severity,
		Status:    incident.StatusActive,
		ChatID:    chatID,
		CreatedBy: userID,
	}

	err := s.tx.WithTx(ctx, func(incidents incident.Repository, events event.Repository) error {
		if err := incidents.Create(ctx, inc); err != nil {
			return err
		}
		return events.Create(ctx, &event.Event{
			IncidentID: inc.ID,
			Type:       event.TypeIncidentCreated,
			UserID:     userID,
			Username:   username,
			Message:    title,
		})
	})
	if err != nil {
		return nil, err
	}

	return inc, nil
}

// AddTimelineEvent appends a comment to the chat's active incident timeline.
//
// Returns errs.ErrNoActiveIncident if the chat has no active incident, or an
// errs.KindValidation error if the message is empty.
func (s *Service) AddTimelineEvent(
	ctx context.Context,
	chatID int64,
	userID *int64,
	username string,
	message string,
) (*event.Event, error) {
	const op = "service.AddTimelineEvent"

	message = strings.TrimSpace(message)
	if message == "" {
		return nil, errs.New(errs.KindValidation, op, "message is required")
	}

	active, err := s.incidents.GetActiveByChatID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	e := &event.Event{
		IncidentID: active.ID,
		Type:       event.TypeCommentAdded,
		UserID:     userID,
		Username:   username,
		Message:    message,
	}
	if err := s.events.Create(ctx, e); err != nil {
		return nil, err
	}

	return e, nil
}

func validSeverity(s incident.Severity) bool {
	switch s {
	case incident.SeverityLow, incident.SeverityMedium, incident.SeverityHigh:
		return true
	default:
		return false
	}
}
