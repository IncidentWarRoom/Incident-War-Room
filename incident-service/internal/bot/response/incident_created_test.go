package response

import (
	"strings"
	"testing"
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/google/uuid"
)

func TestIncidentCreated(t *testing.T) {
	inc := incident.Incident{
		ID:        uuid.MustParse("a1b2c3d4-0000-0000-0000-000000000000"),
		Title:     "DB is down",
		Severity:  incident.SeverityHigh,
		Status:    incident.StatusActive,
		CreatedAt: time.Date(2026, 6, 13, 10, 30, 0, 0, time.UTC),
	}

	got := IncidentCreated(inc, "")

	for _, want := range []string{
		"<b>Incident created</b>",
		"DB is down",
		"🔴 HIGH",
		"ACTIVE",
		"a1b2c3d4",
		"2026-06-13 10:30 UTC",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("IncidentCreated() = %q, missing %q", got, want)
		}
	}
}

func TestIncidentCreatedIncludesTopicLink(t *testing.T) {
	inc := incident.Incident{ID: uuid.New(), Title: "outage", Status: incident.StatusActive}

	got := IncidentCreated(inc, "https://t.me/c/123/45")

	if !strings.Contains(got, `href="https://t.me/c/123/45"`) {
		t.Errorf("IncidentCreated() = %q, expected a link to the topic", got)
	}
}

func TestIncidentCreatedOmitsTopicLinkWhenEmpty(t *testing.T) {
	inc := incident.Incident{ID: uuid.New(), Title: "outage", Status: incident.StatusActive}

	got := IncidentCreated(inc, "")

	if strings.Contains(got, "Open incident topic") {
		t.Errorf("IncidentCreated() = %q, expected no topic link line", got)
	}
}

func TestIncidentCreatedEscapesTitle(t *testing.T) {
	inc := incident.Incident{
		ID:       uuid.New(),
		Title:    "<script>alert(1)</script>",
		Severity: incident.SeverityLow,
		Status:   incident.StatusActive,
	}

	got := IncidentCreated(inc, "")

	if strings.Contains(got, "<script>") {
		t.Errorf("IncidentCreated() did not escape title: %q", got)
	}
	if !strings.Contains(got, "&lt;script&gt;") {
		t.Errorf("IncidentCreated() = %q, expected escaped title", got)
	}
}
