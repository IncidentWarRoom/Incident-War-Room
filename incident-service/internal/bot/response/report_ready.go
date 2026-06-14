package response

import (
	"fmt"
	"strings"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

func ReportReady(inc incident.Incident, reportURL string) string {
	var b strings.Builder

	b.WriteString("📄 <b>Report ready</b>\n\n")
	fmt.Fprintf(&b, "<b>Incident:</b> %s\n", escape(inc.Title))
	fmt.Fprintf(&b, "<b>ID:</b> <code>%s</code>", shortID(inc.ID))

	if reportURL != "" {
		fmt.Fprintf(&b, "\n\n<a href=\"%s\">Download report</a>", escape(reportURL))
	}

	return b.String()
}
