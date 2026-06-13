package response

import (
	"fmt"
	"strings"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

// IncidentCreated renders the confirmation shown when a new incident is opened.
func IncidentCreated(inc incident.Incident) string {
	var b strings.Builder

	b.WriteString("🚨 <b>Incident created</b>\n\n")
	fmt.Fprintf(&b, "<b>Title:</b> %s\n", escape(inc.Title))
	fmt.Fprintf(&b, "<b>Severity:</b> %s\n", severityIcon(inc.Severity))
	fmt.Fprintf(&b, "<b>Status:</b> %s\n", escape(string(inc.Status)))
	fmt.Fprintf(&b, "<b>ID:</b> <code>%s</code>\n", shortID(inc.ID))
	fmt.Fprintf(&b, "<b>Created:</b> %s", formatTime(inc.CreatedAt))

	return b.String()
}
