// Package alert receives Alertmanager webhook notifications from an existing
// Prometheus deployment and opens incidents in the war room for firing alerts.
package alert

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

// Opener opens an incident from an external monitoring alert. It is implemented
// by the bot handler.
type Opener interface {
	OpenIncidentFromAlert(ctx context.Context, title string, severity incident.Severity) (*incident.Incident, error)
}

// Payload is the subset of the Alertmanager webhook body (schema version 4)
// that we consume.
type Payload struct {
	Alerts []Alert `json:"alerts"`
}

type Alert struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

// Handler exposes an HTTP endpoint that Alertmanager posts to. When token is
// non-empty, requests must carry a matching Bearer Authorization header.
type Handler struct {
	opener Opener
	token  string
}

func NewHandler(opener Opener, token string) *Handler {
	return &Handler{opener: opener, token: token}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.authorized(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	for _, a := range payload.Alerts {
		if !strings.EqualFold(a.Status, "firing") {
			continue
		}
		if _, err := h.opener.OpenIncidentFromAlert(r.Context(), title(a), severity(a)); err != nil {
			log.Printf("alert: open incident for %q: %v", title(a), err)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) authorized(r *http.Request) bool {
	if h.token == "" {
		return true
	}
	return r.Header.Get("Authorization") == "Bearer "+h.token
}

func title(a Alert) string {
	for _, key := range []string{"summary", "title", "description"} {
		if v := strings.TrimSpace(a.Annotations[key]); v != "" {
			return v
		}
	}
	if v := strings.TrimSpace(a.Labels["alertname"]); v != "" {
		return v
	}
	return "Prometheus alert"
}

func severity(a Alert) incident.Severity {
	switch strings.ToLower(a.Labels["severity"]) {
	case "critical", "page", "emergency", "high":
		return incident.SeverityHigh
	case "info", "information", "low":
		return incident.SeverityLow
	default:
		return incident.SeverityMedium
	}
}
