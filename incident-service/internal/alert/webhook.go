// Package alert receives Alertmanager webhook notifications from an existing
// Prometheus deployment, opening incidents in the war room for firing alerts
// and closing them again once the alerts resolve.
package alert

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/cQu1x/Incident-War-Room/internal/domain/incident"
)

// Incidents opens and closes incidents in response to monitoring alerts. It is
// implemented by the bot handler.
type Incidents interface {
	OpenIncidentFromAlert(ctx context.Context, title string, severity incident.Severity) (*incident.Incident, error)
	CloseIncidentFromAlert(ctx context.Context, chatID, topicID int64) error
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
	Fingerprint string            `json:"fingerprint"`
}

type incidentRef struct {
	chatID  int64
	topicID int64
}

// Handler exposes an HTTP endpoint that Alertmanager posts to. When token is
// non-empty, requests must carry a matching Bearer Authorization header.
//
// It tracks which incident each firing alert opened so the matching resolved
// notification can close it again. The mapping is kept in memory and is
// rebuilt as new alerts arrive.
type Handler struct {
	incidents Incidents
	token     string

	mu     sync.Mutex
	opened map[string]incidentRef
}

func NewHandler(incidents Incidents, token string) *Handler {
	return &Handler{
		incidents: incidents,
		token:     token,
		opened:    make(map[string]incidentRef),
	}
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
		if strings.EqualFold(a.Status, "resolved") {
			h.resolve(r.Context(), a)
			continue
		}
		h.fire(r.Context(), a)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) fire(ctx context.Context, a Alert) {
	inc, err := h.incidents.OpenIncidentFromAlert(ctx, title(a), severity(a))
	if err != nil {
		log.Printf("alert: open incident for %q: %v", title(a), err)
		return
	}
	h.remember(alertKey(a), incidentRef{chatID: inc.ChatID, topicID: inc.TopicID})
}

func (h *Handler) resolve(ctx context.Context, a Alert) {
	ref, ok := h.forget(alertKey(a))
	if !ok {
		return
	}
	if err := h.incidents.CloseIncidentFromAlert(ctx, ref.chatID, ref.topicID); err != nil {
		log.Printf("alert: close incident for %q: %v", title(a), err)
	}
}

func (h *Handler) remember(key string, ref incidentRef) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.opened[key] = ref
}

func (h *Handler) forget(key string) (incidentRef, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	ref, ok := h.opened[key]
	delete(h.opened, key)
	return ref, ok
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

// alertKey identifies an alert across its firing and resolved notifications. It
// prefers the Alertmanager fingerprint and falls back to the label set.
func alertKey(a Alert) string {
	if a.Fingerprint != "" {
		return a.Fingerprint
	}
	keys := make([]string, 0, len(a.Labels))
	for k := range a.Labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(a.Labels[k])
		b.WriteByte(',')
	}
	return b.String()
}
