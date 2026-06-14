// Package reportclient is the infrastructure adapter that renders incident
// reports by calling the Python report-service over HTTP. It implements the
// domain report.Generator port.
//
// The package owns its wire DTOs (the contract of
// POST /api/v1/reports/generate) and maps the domain report.Report onto them,
// so the external schema and the domain model can evolve independently.
// Failures are wrapped into *errs.Error: an unreachable or failing service
// yields errs.KindUnavailable.
package reportclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/report"
	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

const generatePath = "/api/v1/reports/generate"

// Client talks to the report service over HTTP. It implements report.Generator.
type Client struct {
	baseURL string
	http    *http.Client
}

// Option customizes a Client.
type Option func(*Client)

// WithHTTPClient sets a custom *http.Client (e.g. with a tuned transport).
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) { c.http = h }
}

// WithTimeout sets the per-request timeout of the default HTTP client.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.http.Timeout = d }
}

// New creates a report client targeting baseURL (e.g. "http://localhost:8000").
// A trailing slash is trimmed.
func New(baseURL string, opts ...Option) *Client {
	c := &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Generate maps r onto the report-service wire contract, posts it and returns
// the rendered PDF bytes. A network failure or a 5xx response yields
// errs.KindUnavailable; a 4xx response yields errs.KindValidation.
func (c *Client) Generate(ctx context.Context, r report.Report) ([]byte, error) {
	const op = "reportclient.Generate"

	body, err := json.Marshal(toRequest(r))
	if err != nil {
		return nil, errs.Wrapf(errs.KindInternal, op, err, "marshal request")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+generatePath, bytes.NewReader(body))
	if err != nil {
		return nil, errs.Wrapf(errs.KindInternal, op, err, "build request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/pdf")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, errs.Wrapf(errs.KindUnavailable, op, err, "call report service")
	}
	defer resp.Body.Close()

	pdf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.Wrapf(errs.KindUnavailable, op, err, "read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errs.New(kindForStatus(resp.StatusCode), op,
			"report service returned "+resp.Status+": "+snippet(pdf))
	}

	return pdf, nil
}

// kindForStatus maps an HTTP status code to an error Kind.
func kindForStatus(status int) errs.Kind {
	switch {
	case status == http.StatusNotFound:
		return errs.KindNotFound
	case status >= 400 && status < 500:
		return errs.KindValidation
	default:
		return errs.KindUnavailable
	}
}

// snippet returns a short, single-line preview of an error response body.
func snippet(b []byte) string {
	const max = 256
	s := strings.TrimSpace(string(b))
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > max {
		return s[:max] + "…"
	}
	return s
}
