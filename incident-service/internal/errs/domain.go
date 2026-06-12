package errs

// Sentinel errors for known business-rule failures.
// Compare with errors.Is.
var (
	ErrIncidentNotFound      = New(KindNotFound, "incident", "incident not found")
	ErrNoActiveIncident      = New(KindNotFound, "incident", "no active incident in this chat")
	ErrIncidentAlreadyActive = New(KindConflict, "incident", "an active incident already exists in this chat")
	ErrIncidentAlreadyClosed = New(KindConflict, "incident", "incident is already closed")
)
