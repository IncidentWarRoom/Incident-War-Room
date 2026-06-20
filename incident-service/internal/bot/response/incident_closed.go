package response

import (
	"fmt"
	"strings"
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

func IncidentClosed(inc incident.Incident, timelineURLs []string, reportURL string) string {
	var b strings.Builder

	b.WriteString("✅ <b>Incident closed</b>\n\n")
	fmt.Fprintf(&b, "<b>Title:</b> %s\n", escape(inc.Title))
	fmt.Fprintf(&b, "<b>ID:</b> <code>%s</code>\n", shortID(inc.ID))

	if inc.ClosedAt != nil {
		fmt.Fprintf(&b, "<b>Closed:</b> %s\n", formatTime(*inc.ClosedAt))
		fmt.Fprintf(&b, "<b>Duration:</b> %s", escape(formatDuration(inc.ClosedAt.Sub(inc.CreatedAt))))
	} else {
		b.WriteString("<b>Closed:</b> just now")
	}

	b.WriteString("\n\n📋 <b>Timeline</b>\n")
	if len(timelineURLs) == 0 {
		b.WriteString("<i>Telegraph timeline pages will be linked here.</i>")
	} else {
		for i, url := range timelineURLs {
			if i > 0 {
				b.WriteByte('\n')
			}
			b.WriteString(escape(url))
		}
	}

	b.WriteString("\n\n📄 <b>Report</b>\n")
	if reportURL == "" {
		b.WriteString("<i>The report could not be generated right now.</i>")
	} else {
		fmt.Fprintf(&b, "<a href=\"%s\">Download report</a>", escape(reportURL))
	}

	return b.String()
}

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
