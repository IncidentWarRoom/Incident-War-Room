package service

import (
	"context"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/report"
)

// GenerateReport renders a PDF report for the chat's active incident. It loads
// the incident together with its timeline, derives the participants from the
// timeline authors and delegates rendering to the report.Generator port.
//
// Returns errs.ErrNoActiveIncident if the chat has no active incident, or an
// errs.KindUnavailable error if the report service is unreachable.
func (s *Service) GenerateReport(ctx context.Context, chatID int64) ([]byte, error) {
	inc, events, err := s.GetTimeline(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return s.reports.Generate(ctx, report.Report{
		Incident:     *inc,
		Participants: participantsFromEvents(events),
		Timeline:     events,
	})
}

// participantsFromEvents returns the distinct authors of the events, preserving
// first-seen order. Events without an author (system events) are skipped.
func participantsFromEvents(events []event.Event) []report.Participant {
	seen := make(map[int64]struct{})
	var participants []report.Participant
	for _, e := range events {
		if e.AuthorID == nil {
			continue
		}
		if _, ok := seen[*e.AuthorID]; ok {
			continue
		}
		seen[*e.AuthorID] = struct{}{}
		participants = append(participants, report.Participant{
			UserID:   *e.AuthorID,
			Username: e.Username,
		})
	}
	return participants
}
