package response

import (
	"fmt"
	"strings"
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

// IncidentClosed renders the confirmation shown when an incident is closed,
// including how long it stayed open when the closing time is known.
func IncidentClosed(inc incident.Incident) string {
	var b strings.Builder

	b.WriteString("✅ <b>Incident closed</b>\n\n")
	fmt.Fprintf(&b, "<b>Title:</b> %s\n", escape(inc.Title))
	fmt.Fprintf(&b, "<b>ID:</b> <code>%s</code>\n", shortID(inc.ID))

	if inc.ClosedAt != nil {
		fmt.Fprintf(&b, "<b>Closed:</b> %s\n", formatTime(*inc.ClosedAt))
		fmt.Fprintf(&b, "<b>Duration:</b> %s", formatDuration(inc.ClosedAt.Sub(inc.CreatedAt)))
	} else {
		b.WriteString("<b>Closed:</b> just now")
	}

	return b.String()
}

// formatDuration renders a duration as a compact "1h 5m" string, dropping
// zero-valued leading units. Sub-minute durations render as "<1m".
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "<1m"
	}

	h := int(d / time.Hour)
	m := int((d % time.Hour) / time.Minute)

	switch {
	case h == 0:
		return fmt.Sprintf("%dm", m)
	case m == 0:
		return fmt.Sprintf("%dh", h)
	default:
		return fmt.Sprintf("%dh %dm", h, m)
	}
}
