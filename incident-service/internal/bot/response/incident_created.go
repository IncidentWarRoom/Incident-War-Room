package response

import (
	"fmt"
	"strings"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

// IncidentCreated renders the announcement posted in the main chat when an
// incident is opened. topicURL links to the dedicated incident topic; when
// empty the link line is omitted.
func IncidentCreated(inc incident.Incident, topicURL string) string {
	var b strings.Builder

	b.WriteString("🚨 <b>Incident created</b>\n\n")
	fmt.Fprintf(&b, "<b>Title:</b> %s\n", escape(inc.Title))
	fmt.Fprintf(&b, "<b>Severity:</b> %s\n", severityIcon(inc.Severity))
	fmt.Fprintf(&b, "<b>Status:</b> %s\n", escape(string(inc.Status)))
	fmt.Fprintf(&b, "<b>ID:</b> <code>%s</code>\n", shortID(inc.ID))
	fmt.Fprintf(&b, "<b>Created:</b> %s", formatTime(inc.CreatedAt))

	if topicURL != "" {
		fmt.Fprintf(&b, "\n\n📌 <a href=\"%s\">Open incident topic</a>", escape(topicURL))
	}

	return b.String()
}
