package response

import (
	"fmt"
	"strings"

	"github.com/cQu1x/Incident-War-Room/internal/domain/event"
	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

func Timeline(inc incident.Incident, events []event.Event) string {
	var b strings.Builder

	fmt.Fprintf(&b, "📋 <b>Timeline</b> — %s\n", escape(inc.Title))

	if len(events) == 0 {
		b.WriteString("\nThe incident timeline is empty.")
		return b.String()
	}

	for _, e := range events {
		fmt.Fprintf(&b, "\n<b>%s</b> — %s: %s",
			formatTime(e.CreatedAt), escape(e.Username), escape(e.Message))
	}

	return b.String()
}
