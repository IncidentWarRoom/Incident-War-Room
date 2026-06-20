package telegraphclient

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

// Telegraph limits a page title to 256 characters and the rendered content to
// roughly 64 KB. We chunk the timeline so a single page stays comfortably
// under the content limit.
const (
	titleLimit      = 256
	maxContentBytes = 60 * 1024
	timeLayout      = "2006-01-02 15:04 MST"
)

// node is a Telegraph DOM node. A node with an empty Tag and a single string
// child is encoded as a plain text node by buildChildren; everything else is a
// NodeElement. Children may hold strings or further nodes.
type node struct {
	Tag      string `json:"tag,omitempty"`
	Attrs    *attrs `json:"attrs,omitempty"`
	Children []any  `json:"children,omitempty"`
}

type attrs struct {
	Href string `json:"href,omitempty"`
}

// page is a single Telegraph page ready to be created: a title and its content
// (a list of top-level nodes).
type page struct {
	title   string
	content []any
}

// buildPages renders a timeline into one or more Telegraph pages. The header
// (incident metadata) is repeated on every page so each page is readable on
// its own; events are split across pages to respect the content size limit.
func buildPages(inc incident.Incident, events []event.Event) []page {
	header := headerNodes(inc)

	var pages []page
	content := append([]any{}, header...)
	size := nodesSize(content)

	flush := func() {
		pages = append(pages, page{content: content})
		content = append([]any{}, header...)
		size = nodesSize(header)
	}

	for _, e := range events {
		n := eventNode(e)
		ns := nodeSize(n)
		if len(content) > len(header) && size+ns > maxContentBytes {
			flush()
		}
		content = append(content, n)
		size += ns
	}
	pages = append(pages, page{content: content})

	titlePages(pages, inc.Title)
	return pages
}

// titlePages sets each page's title, appending a "(i/n)" suffix when the
// timeline spans more than one page.
func titlePages(pages []page, title string) {
	total := len(pages)
	for i := range pages {
		t := title
		if total > 1 {
			t = fmt.Sprintf("%s (%d/%d)", title, i+1, total)
		}
		pages[i].title = truncate(t, titleLimit)
	}
}

func headerNodes(inc incident.Incident) []any {
	meta := fmt.Sprintf("Severity: %s · Status: %s · Opened: %s",
		inc.Severity, inc.Status, formatTime(inc.CreatedAt))
	if inc.ClosedAt != nil {
		meta += " · Closed: " + formatTime(*inc.ClosedAt)
	}

	return []any{
		element("h3", text("Timeline")),
		element("p", element("i", text(meta))),
	}
}

func eventNode(e event.Event) node {
	who := e.Username
	if who == "" {
		who = "system"
	}

	children := []any{
		element("b", text(formatTime(e.CreatedAt))),
		text(fmt.Sprintf(" — %s", who)),
	}
	if e.Message != "" {
		children = append(children, text(": "+e.Message))
	}

	return node{Tag: "p", Children: children}
}

// element builds a NodeElement with the given tag and children.
func element(tag string, children ...any) node {
	return node{Tag: tag, Children: children}
}

// text builds a plain text node (Telegraph encodes a bare string as a text
// node).
func text(s string) string { return s }

func formatTime(t time.Time) string { return t.Format(timeLayout) }

func truncate(s string, limit int) string {
	r := []rune(s)
	if len(r) <= limit {
		return s
	}
	return string(r[:limit])
}

// nodeSize estimates the serialized size of a node; errors are treated as zero
// since they only affect chunk boundaries, not correctness.
func nodeSize(n node) int {
	b, err := json.Marshal(n)
	if err != nil {
		return 0
	}
	return len(b)
}

func nodesSize(nodes []any) int {
	b, err := json.Marshal(nodes)
	if err != nil {
		return 0
	}
	return len(b)
}
