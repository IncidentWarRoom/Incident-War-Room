package telegraphclient

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
	"github.com/cQu1x/Incident-War-Room/internal/domain/timeline"
)

type fakeTelegraph struct {
	mu          sync.Mutex
	createCalls int
	editCalls   int
	editBodies  []string
}

func (f *fakeTelegraph) handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/createAccount", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `{"ok":true,"result":{"access_token":"tok"}}`)
	})

	mux.HandleFunc("/createPage", func(w http.ResponseWriter, _ *http.Request) {
		f.mu.Lock()
		f.createCalls++
		n := f.createCalls
		f.mu.Unlock()
		fmt.Fprintf(w, `{"ok":true,"result":{"url":"https://telegra.ph/p%d","path":"p%d"}}`, n, n)
	})

	mux.HandleFunc("/editPage/", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		f.mu.Lock()
		f.editCalls++
		f.editBodies = append(f.editBodies, r.FormValue("content"))
		f.mu.Unlock()
		fmt.Fprint(w, `{"ok":true,"result":{"url":"https://telegra.ph/edited"}}`)
	})

	return mux
}

func newClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()
	return New(WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
}

func bigTimeline(n int) timeline.Timeline {
	events := make([]event.Event, n)
	base := time.Date(2026, 6, 20, 10, 0, 0, 0, time.UTC)
	for i := range events {
		events[i] = event.Event{
			Username:  "alice",
			Message:   strings.Repeat("z", 500),
			CreatedAt: base.Add(time.Duration(i) * time.Minute),
		}
	}
	return timeline.Timeline{
		Incident: incident.Incident{Title: "outage", Severity: incident.SeverityHigh, Status: incident.StatusActive},
		Events:   events,
	}
}

func TestPublishSinglePageDoesNotEdit(t *testing.T) {
	fake := &fakeTelegraph{}
	srv := httptest.NewServer(fake.handler())
	defer srv.Close()

	urls, err := newClient(t, srv).Publish(context.Background(), bigTimeline(2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(urls) != 1 {
		t.Fatalf("expected 1 url, got %v", urls)
	}
	if fake.editCalls != 0 {
		t.Fatalf("single page should not be edited, got %d edits", fake.editCalls)
	}
}

func TestPublishPaginatesAcrossPages(t *testing.T) {
	fake := &fakeTelegraph{}
	srv := httptest.NewServer(fake.handler())
	defer srv.Close()

	urls, err := newClient(t, srv).Publish(context.Background(), bigTimeline(400))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(urls) < 2 {
		t.Fatalf("expected the timeline to span multiple pages, got %v", urls)
	}
	if fake.editCalls != len(urls) {
		t.Fatalf("expected one edit per page for nav, got %d edits for %d pages", fake.editCalls, len(urls))
	}

	first := fake.editBodies[0]
	if !strings.Contains(first, urls[len(urls)-1]) {
		t.Errorf("first page nav should link to the last page %q: %s", urls[len(urls)-1], first)
	}
}
