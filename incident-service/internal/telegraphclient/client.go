// Package telegraphclient is the infrastructure adapter that publishes an
// incident timeline to Telegraph (https://telegra.ph) over its HTTP API. It
// implements the domain timeline.Publisher port.
//
// A long timeline is split across several Telegraph pages; Publish returns
// their URLs in reading order. Failures are wrapped into *errs.Error: an
// unreachable or failing service yields errs.KindUnavailable.
package telegraphclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/cQu1x/Incident-War-Room/internal/domain/timeline"
	"github.com/cQu1x/Incident-War-Room/internal/errs"
)

const (
	defaultBaseURL   = "https://api.telegra.ph"
	defaultShortName = "Incident War Room"
)

// Client talks to the Telegraph API over HTTP. It implements
// timeline.Publisher.
//
// An access token can be supplied up front (WithAccessToken); otherwise the
// client creates a throwaway Telegraph account on first use and caches the
// token for the rest of its lifetime.
type Client struct {
	baseURL    string
	authorName string
	http       *http.Client

	mu    sync.Mutex
	token string
}

// Option customizes a Client.
type Option func(*Client)

// WithHTTPClient sets a custom *http.Client (e.g. with a tuned transport).
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) { c.http = h }
}

// WithAccessToken pins a Telegraph access token, skipping account creation.
func WithAccessToken(token string) Option {
	return func(c *Client) { c.token = strings.TrimSpace(token) }
}

// WithAuthorName sets the author shown on published pages.
func WithAuthorName(name string) Option {
	return func(c *Client) { c.authorName = name }
}

// New creates a Telegraph client. With no options it targets the public
// Telegraph API and creates an anonymous account lazily.
func New(opts ...Option) *Client {
	c := &Client{
		baseURL:    defaultBaseURL,
		authorName: defaultShortName,
		http:       &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Publish renders t into one or more Telegraph pages and returns their URLs in
// reading order. A network failure or a Telegraph error yields
// errs.KindUnavailable.
func (c *Client) Publish(ctx context.Context, t timeline.Timeline) ([]string, error) {
	const op = "telegraphclient.Publish"

	token, err := c.ensureToken(ctx)
	if err != nil {
		return nil, err
	}

	pages := buildPages(t.Incident, t.Events)
	urls := make([]string, 0, len(pages))
	for _, p := range pages {
		u, err := c.createPage(ctx, token, p)
		if err != nil {
			return nil, errs.Wrapf(errs.KindUnavailable, op, err, "create telegraph page")
		}
		urls = append(urls, u)
	}

	return urls, nil
}

// ensureToken returns the cached access token, creating a Telegraph account on
// first use.
func (c *Client) ensureToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != "" {
		return c.token, nil
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := c.call(ctx, "createAccount", url.Values{
		"short_name":  {defaultShortName},
		"author_name": {c.authorName},
	}, &result); err != nil {
		return "", errs.Wrapf(errs.KindUnavailable, "telegraphclient.ensureToken", err, "create telegraph account")
	}

	c.token = result.AccessToken
	return c.token, nil
}

// createPage creates a single Telegraph page and returns its public URL.
func (c *Client) createPage(ctx context.Context, token string, p page) (string, error) {
	content, err := json.Marshal(p.content)
	if err != nil {
		return "", errs.Wrapf(errs.KindInternal, "telegraphclient.createPage", err, "marshal content")
	}

	var result struct {
		URL string `json:"url"`
	}
	if err := c.call(ctx, "createPage", url.Values{
		"access_token": {token},
		"title":        {p.title},
		"author_name":  {c.authorName},
		"content":      {string(content)},
	}, &result); err != nil {
		return "", err
	}

	return result.URL, nil
}

// call posts form values to a Telegraph API method and decodes result into out.
// A Telegraph-level failure (ok=false) is returned as an *errs.Error.
func (c *Client) call(ctx context.Context, method string, form url.Values, out any) error {
	const op = "telegraphclient.call"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/"+method, strings.NewReader(form.Encode()))
	if err != nil {
		return errs.Wrapf(errs.KindInternal, op, err, "build request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return errs.Wrapf(errs.KindUnavailable, op, err, "call telegraph")
	}
	defer resp.Body.Close()

	var envelope struct {
		OK     bool            `json:"ok"`
		Error  string          `json:"error"`
		Result json.RawMessage `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return errs.Wrapf(errs.KindUnavailable, op, err, "decode telegraph response")
	}
	if !envelope.OK {
		return errs.New(errs.KindUnavailable, op, "telegraph error: "+envelope.Error)
	}

	if out == nil {
		return nil
	}
	if err := json.Unmarshal(envelope.Result, out); err != nil {
		return errs.Wrapf(errs.KindInternal, op, err, "decode telegraph result")
	}
	return nil
}
