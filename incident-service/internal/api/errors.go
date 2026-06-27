package api

import (
	"log"
	"net/http"

	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

type errorResponse struct {
	Error string `json:"error"`
}

// writeError maps a service error to its HTTP status and a JSON body. Internal
// failures are logged and hidden behind a generic message; client-facing kinds
// carry their own description.
func writeError(w http.ResponseWriter, err error) {
	status := statusForKind(errs.KindOf(err))
	message := err.Error()

	if status == http.StatusInternalServerError {
		log.Printf("api: %v", err)
		message = "internal error"
	}

	writeJSON(w, status, errorResponse{Error: message})
}

func statusForKind(kind errs.Kind) int {
	switch kind {
	case errs.KindNotFound:
		return http.StatusNotFound
	case errs.KindValidation:
		return http.StatusBadRequest
	case errs.KindConflict:
		return http.StatusConflict
	case errs.KindUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
