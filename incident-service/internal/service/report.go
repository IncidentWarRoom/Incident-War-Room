package service

import (
	"context"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/report"
)

func (s *Service) GenerateReport(ctx context.Context, chatID, topicID int64) (report.Document, error) {
	inc, events, err := s.GetTimeline(ctx, chatID, topicID)
	if err != nil {
		return report.Document{}, err
	}

	doc, err := s.reports.Generate(ctx, report.Report{
		Incident:     *inc,
		Participants: participantsFromEvents(events),
		Timeline:     events,
	})
	if err != nil {
		return report.Document{}, err
	}

	if doc.URL != "" {
		if err := s.incidents.UpdateReportURL(ctx, inc.ID, doc.URL); err != nil {
			return report.Document{}, err
		}
		inc.ReportURL = &doc.URL
	}

	return doc, nil
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
