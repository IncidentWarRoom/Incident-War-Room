package service

import (
	"context"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/report"
)

// GenerateReport renders a PDF report for the chat's active incident. It loads
// the incident together with its timeline, derives the participants from the
// timeline users and delegates rendering to the report.Generator port.
//
// Returns errs.ErrNoActiveIncident if the chat has no active incident, or an
// errs.KindUnavailable error if the report service is unreachable.
func (s *Service) GenerateReport(ctx context.Context, chatID, topicID int64) ([]byte, error) {
	inc, events, err := s.GetTimeline(ctx, chatID, topicID)
	if err != nil {
		return nil, err
	}

	return s.reports.Generate(ctx, report.Report{
		Incident:     *inc,
		Participants: participantsFromEvents(events),
		Timeline:     events,
	})
}

// participantsFromEvents returns the distinct users of the events, preserving
// first-seen order. Events without a user (system events) are skipped.
func participantsFromEvents(events []event.Event) []report.Participant {
	seen := make(map[int64]struct{})
	var participants []report.Participant
	for _, e := range events {
		if e.UserID == nil {
			continue
		}
		if _, ok := seen[*e.UserID]; ok {
			continue
		}
		seen[*e.UserID] = struct{}{}
		participants = append(participants, report.Participant{
			UserID:   *e.UserID,
			Username: e.Username,
		})
	}
	return participants
}
