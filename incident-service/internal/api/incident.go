package api

import "net/http"

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
