package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

func (s *Server) listIncidents(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.context(r)
	defer cancel()

	incidents, err := s.svc.ListIncidents(ctx)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newIncidentResponses(incidents))
}

func (s *Server) getIncident(w http.ResponseWriter, r *http.Request) {
	id, err := incidentID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	ctx, cancel := s.context(r)
	defer cancel()

	inc, err := s.svc.GetIncident(ctx, id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newIncidentResponse(*inc))
}

func incidentID(r *http.Request) (uuid.UUID, error) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return uuid.Nil, errs.New(errs.KindValidation, "api.incidentID", "invalid incident id")
	}
	return id, nil
}
