package bot

import (
	"strings"
	"testing"
)

func TestHandleTimeline(t *testing.T) {
	ctx := &mockContext{}

	if err := HandleTimeline(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := lastSent(t, ctx); !strings.Contains(got, "timeline is empty") {
		t.Errorf("unexpected reply: %q", got)
	}
}
