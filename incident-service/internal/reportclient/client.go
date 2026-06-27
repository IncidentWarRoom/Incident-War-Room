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

const (
	generateURLPath    = "/api/v1/reports/generate-url"
	generateInlinePath = "/api/v1/reports/generate-inline"
)

// Client talks to the report service over HTTP. It implements report.Generator.
type Client struct {
	baseURL   string
	http      *http.Client
	s3Enabled bool
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

// WithS3Enabled selects the generation mode. When enabled the client calls
// generate-url and expects a public report URL; otherwise it calls
// generate-inline and reads the PDF bytes directly.
func WithS3Enabled(enabled bool) Option {
	return func(c *Client) { c.s3Enabled = enabled }
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

func (c *Client) Generate(ctx context.Context, r report.Report) (report.Document, error) {
	const op = "reportclient.Generate"

	body, err := json.Marshal(toRequest(r))
	if err != nil {
		return report.Document{}, errs.Wrapf(errs.KindInternal, op, err, "marshal request")
	}

	if c.s3Enabled {
		return c.generateURL(ctx, body)
	}
	return c.generateInline(ctx, body)
}

func (c *Client) generateURL(ctx context.Context, body []byte) (report.Document, error) {
	const op = "reportclient.generateURL"

	data, err := c.post(ctx, op, generateURLPath, "application/json", body)
	if err != nil {
		return report.Document{}, err
	}

	var res response
	if err := json.Unmarshal(data, &res); err != nil {
		return report.Document{}, errs.Wrapf(errs.KindUnavailable, op, err, "decode response body")
	}
	if res.ReportURL == "" {
		return report.Document{}, errs.New(errs.KindUnavailable, op, "report service returned an empty report URL")
	}

	return report.Document{URL: res.ReportURL}, nil
}

func (c *Client) generateInline(ctx context.Context, body []byte) (report.Document, error) {
	const op = "reportclient.generateInline"

	data, err := c.post(ctx, op, generateInlinePath, "application/pdf", body)
	if err != nil {
		return report.Document{}, err
	}
	if len(data) == 0 {
		return report.Document{}, errs.New(errs.KindUnavailable, op, "report service returned an empty report")
	}

	return report.Document{PDF: data}, nil
}

func (c *Client) post(ctx context.Context, op, path, accept string, body []byte) ([]byte, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, errs.Wrapf(errs.KindInternal, op, err, "build request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", accept)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, errs.Wrapf(errs.KindUnavailable, op, err, "call report service")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.Wrapf(errs.KindUnavailable, op, err, "read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errs.New(kindForStatus(resp.StatusCode), op,
			"report service returned "+resp.Status+": "+snippet(data))
	}

	return data, nil
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
