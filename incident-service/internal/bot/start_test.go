package bot

import "testing"

func TestHandleStart(t *testing.T) {
	h := New(&fakeService{}, newFakeAPI())
	ctx := &mockContext{}

	if err := h.HandleStart(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := lastSent(t, ctx); got != "Incident War Room is running." {
		t.Errorf("unexpected reply: %q", got)
	}
}
