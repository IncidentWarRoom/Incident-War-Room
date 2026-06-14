package bot

import (
	"strings"
	"testing"
)

func TestHandleIncident(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantPart string
	}{
		{
			name:     "no args shows usage",
			args:     nil,
			wantPart: "Usage:",
		},
		{
			name:     "create",
			args:     []string{"create"},
			wantPart: "<b>Incident created</b>",
		},
		{
			name:     "close",
			args:     []string{"close"},
			wantPart: "<b>Incident closed</b>",
		},
		{
			name:     "message adds timeline update",
			args:     []string{"db", "is", "down"},
			wantPart: "Update added to timeline: db is down",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &mockContext{args: tt.args}

			if err := HandleIncident(ctx); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := lastSent(t, ctx); !strings.Contains(got, tt.wantPart) {
				t.Errorf("reply %q does not contain %q", got, tt.wantPart)
			}
		})
	}
}

func TestHandleIncidentUsageListsAllSubcommands(t *testing.T) {
	ctx := &mockContext{}

	if err := HandleIncident(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	usage := lastSent(t, ctx)
	for _, part := range []string{"/incident create", "/incident close", "/incident <message>"} {
		if !strings.Contains(usage, part) {
			t.Errorf("usage %q does not mention %q", usage, part)
		}
	}
}
