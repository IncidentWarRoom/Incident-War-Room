package response

import (
	"strings"
	"testing"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/google/uuid"
)

func TestReportReadyWithURL(t *testing.T) {
	inc := incident.Incident{
		ID:    uuid.MustParse("a1b2c3d4-0000-0000-0000-000000000000"),
		Title: "DB is down",
	}

	got := ReportReady(inc, "https://example.com/report.pdf")

	for _, want := range []string{
		"<b>Report ready</b>",
		"DB is down",
		"a1b2c3d4",
		"<a href=\"https://example.com/report.pdf\">Download report</a>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("ReportReady() = %q, missing %q", got, want)
		}
	}
}

func TestReportReadyWithoutURL(t *testing.T) {
	inc := incident.Incident{ID: uuid.New(), Title: "DB is down"}

	got := ReportReady(inc, "")

	if strings.Contains(got, "Download report") {
		t.Errorf("ReportReady() should omit link without URL: %q", got)
	}
}

func TestReportReadyEscapesTitle(t *testing.T) {
	inc := incident.Incident{ID: uuid.New(), Title: "<b>oops</b>"}

	got := ReportReady(inc, "")

	if strings.Contains(got, "<b>oops</b>") {
		t.Errorf("ReportReady() did not escape title: %q", got)
	}
}
