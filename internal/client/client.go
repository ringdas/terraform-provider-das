// Package client is the provider's one transport to the DAS APIs, the
// same conventions as the `das` CLI (docs/architecture/cli.md, "Hosts and
// credentials"; docs/engineering/api-conventions.md):
//
//   - the Organisation API answers on `https://<org domain>`, a product API
//     on the Network host `https://<network>.<product>.<org domain>`; the
//     Organization authority may carry a port, every host keeps it;
//   - the credential is a bearer: an API key (`DAS_API_KEY`) or a token
//     (`DAS_TOKEN`); the tenant is always the credential's. An Organization
//     API key (`daso_…`) acts on the Organization layer only: a product
//     call takes the token instead, or fails before any request;
//   - every mutation carries one Idempotency-Key, kept across retries, so
//     a retry never doubles a write;
//   - a network error or a 502/503/504 is retried twice; a redirect is
//     never followed (the DAS API does not redirect: a 3xx is a wrong host);
//   - an error is `application/problem+json`, reported with its code,
//     field errors and request id.
package client

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

// Config is how the provider reaches an Organization.
type Config struct {
	// OrgDomain is the Organization's domain (`das.example.com`, or
	// `lvh.me:6100` locally).
	OrgDomain string
	// Scheme is `https` (default) or `http` (a local stack only).
	Scheme string
	// Network is the slug of the Network a product API is called for;
	// empty, it is the API key's own Network.
	Network string
	// APIKey or Token is the bearer credential.
	APIKey string
	Token  string
	// Endpoints overrides the base URL of a module (`org`, `delivery`, ...).
	Endpoints map[string]string
	UserAgent string
	// HTTPClient defaults to one with a 60 s timeout.
	HTTPClient *http.Client
	// RetryDelays defaults to 200 ms and 800 ms; tests pass zeros.
	RetryDelays []time.Duration
}

// Client calls the DAS APIs of one Organization.
type Client struct {
	cfg     Config
	http    *http.Client
	mu      sync.Mutex
	network string
}

// New validates a configuration.
func New(cfg Config) (*Client, error) {
	if cfg.Scheme == "" {
		cfg.Scheme = "https"
	}
	if cfg.Scheme != "https" && cfg.Scheme != "http" {
		return nil, fmt.Errorf("scheme must be https or http, not %q", cfg.Scheme)
	}
	if cfg.OrgDomain == "" && cfg.Endpoints["org"] == "" {
		return nil, errors.New("no Organization: set org_domain or DAS_ORG_DOMAIN")
	}
	if cfg.APIKey == "" && cfg.Token == "" {
		return nil, errors.New("no credential: set api_key (DAS_API_KEY) or token (DAS_TOKEN)")
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{
			Timeout:       60 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		}
	}
	if cfg.RetryDelays == nil {
		cfg.RetryDelays = []time.Duration{200 * time.Millisecond, 800 * time.Millisecond}
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = "terraform-provider-das/dev"
	}
	return &Client{cfg: cfg, http: cfg.HTTPClient, network: strings.ToLower(cfg.Network)}, nil
}

// Describe names where the provider talks to, for a plan's diagnostics.
func (c *Client) Describe() string {
	if c.cfg.OrgDomain != "" {
		return c.cfg.Scheme + "://" + c.cfg.OrgDomain
	}
	return c.cfg.Endpoints["org"]
}

// Request is one call of an operation.
type Request struct {
	// Module and Product of the contract (`org`/`platform`, `delivery`/`delivery`).
	Module  string
	Product string
	Method  string
	// Path with its parameters substituted.
	Path  string
	Query url.Values
	// Body is marshalled as JSON (nil: no body).
	Body any
	// IdempotencyKey of a mutation; empty, a new one is drawn.
	IdempotencyKey string
	// Layer is the operation's `x-das-layer`: "org" or "network". A
	// Network-layer call never carries an Organization API key; empty, a
	// product module's call is Network-layer and the org module's is not.
	Layer string
}

// Response is a 2xx answer.
type Response struct {
	Status int
	// JSON is the decoded body (numbers as json.Number); nil without one.
	JSON any
}

// Problem is a non-2xx answer of the API.
type Problem struct {
	Status    int
	Code      string
	Title     string
	Detail    string
	RequestID string
	Errors    []FieldError
	Method    string
	URL       string
}

// FieldError is one entry of a problem's `errors`.
type FieldError struct {
	Location string `json:"location"`
	Path     string `json:"path"`
	Message  string `json:"message"`
}

func (p *Problem) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s answered %d", p.Method, p.URL, p.Status)
	if p.Code != "" {
		fmt.Fprintf(&b, " %s", p.Code)
	}
	if p.Title != "" {
		fmt.Fprintf(&b, ": %s", p.Title)
	}
	if p.Detail != "" && p.Detail != p.Title {
		fmt.Fprintf(&b, " - %s", p.Detail)
	}
	for _, e := range p.Errors {
		fmt.Fprintf(&b, "\n  %s %s: %s", e.Location, e.Path, e.Message)
	}
	if p.RequestID != "" {
		fmt.Fprintf(&b, "\n(request id %s)", p.RequestID)
	}
	return b.String()
}

// NotFound reports a 404 answer.
func NotFound(err error) bool {
	var p *Problem
	return errors.As(err, &p) && p.Status == http.StatusNotFound
}

// ErrOrgKeyInProduct: an Organization API key never reaches a Network-layer
// operation (a product API, the Organisation API's /api/network/...).
var ErrOrgKeyInProduct = errors.New("an Organization API key acts on the Organization layer only: a Network-layer operation (a product API, /api/network/...) needs a Network credential - set token (and network for a product API)")

// IsOrgAPIKey reports an Organization API key (`daso_…`).
func IsOrgAPIKey(key string) bool { return strings.HasPrefix(key, "daso_") }

// credentialFor is the bearer a call carries: an Organization API key only
// on the Organization layer, the token on the Network layer.
func (c *Client) credentialFor(module, layer string) (string, error) {
	networkLayer := layer == "network" || (layer == "" && module != "org")
	if networkLayer && IsOrgAPIKey(c.cfg.APIKey) {
		if c.cfg.Token == "" {
			return "", ErrOrgKeyInProduct
		}
		return c.cfg.Token, nil
	}
	return c.credential(), nil
}

// Do calls an operation.
func (c *Client) Do(ctx context.Context, r Request) (*Response, error) {
	bearer, err := c.credentialFor(r.Module, r.Layer)
	if err != nil {
		return nil, err
	}
	origin, err := c.origin(ctx, r.Module, r.Product)
	if err != nil {
		return nil, err
	}
	u := origin + r.Path
	if len(r.Query) > 0 {
		u += "?" + r.Query.Encode()
	}
	var body []byte
	if r.Body != nil {
		if body, err = json.Marshal(r.Body); err != nil {
			return nil, err
		}
	}
	key := r.IdempotencyKey
	if r.Method != http.MethodGet && key == "" {
		key = newUUID()
	}
	for attempt := 0; ; attempt++ {
		req, err := http.NewRequestWithContext(ctx, r.Method, u, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json, application/problem+json")
		req.Header.Set("User-Agent", c.cfg.UserAgent)
		req.Header.Set("X-Request-Id", newUUID())
		req.Header.Set("Authorization", "Bearer "+bearer)
		if r.Body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		if key != "" {
			req.Header.Set("Idempotency-Key", key)
		}
		res, err := c.http.Do(req)
		if err != nil {
			if ctx.Err() == nil && attempt < len(c.cfg.RetryDelays) {
				if !sleep(ctx, c.cfg.RetryDelays[attempt]) {
					return nil, ctx.Err()
				}
				continue
			}
			return nil, fmt.Errorf("cannot reach %s: %w", req.URL.Host, err)
		}
		data, readErr := io.ReadAll(io.LimitReader(res.Body, 32<<20))
		_ = res.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		if (res.StatusCode == 502 || res.StatusCode == 503 || res.StatusCode == 504) && attempt < len(c.cfg.RetryDelays) {
			if !sleep(ctx, c.cfg.RetryDelays[attempt]) {
				return nil, ctx.Err()
			}
			continue
		}
		if res.StatusCode >= 300 && res.StatusCode < 400 {
			return nil, fmt.Errorf("%d redirect to %q: %s does not serve the DAS API here - check org_domain, network and endpoints", res.StatusCode, res.Header.Get("Location"), req.URL.Host)
		}
		if res.StatusCode >= 400 {
			return nil, problemOf(res, data, r.Method, u)
		}
		out := &Response{Status: res.StatusCode}
		if len(bytes.TrimSpace(data)) > 0 {
			dec := json.NewDecoder(bytes.NewReader(data))
			dec.UseNumber()
			if err := dec.Decode(&out.JSON); err != nil {
				return nil, fmt.Errorf("%s %s answered %d with a body that is not JSON: %w", r.Method, u, res.StatusCode, err)
			}
		}
		return out, nil
	}
}

// All follows a page answer (`{items, nextCursor}`) to its end.
func (c *Client) All(ctx context.Context, r Request) ([]any, error) {
	var items []any
	query := url.Values{}
	for k, v := range r.Query {
		query[k] = v
	}
	query.Set("limit", "100")
	for pages := 0; pages < 10000; pages++ {
		r.Query = query
		res, err := c.Do(ctx, r)
		if err != nil {
			return nil, err
		}
		page, ok := res.JSON.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("%s %s answered no page object", r.Method, r.Path)
		}
		batch, _ := page["items"].([]any)
		items = append(items, batch...)
		next, _ := page["nextCursor"].(string)
		if next == "" {
			return items, nil
		}
		query.Set("cursor", next)
	}
	return nil, fmt.Errorf("%s %s: more than 10000 pages", r.Method, r.Path)
}

func (c *Client) credential() string {
	if c.cfg.APIKey != "" {
		return c.cfg.APIKey
	}
	return c.cfg.Token
}

func (c *Client) origin(ctx context.Context, module, product string) (string, error) {
	if base := c.cfg.Endpoints[module]; base != "" {
		return strings.TrimRight(base, "/"), nil
	}
	if c.cfg.OrgDomain == "" {
		return "", fmt.Errorf("no endpoint for the %s API: set org_domain, or endpoints.%s", module, module)
	}
	if module == "org" {
		return c.cfg.Scheme + "://" + c.cfg.OrgDomain, nil
	}
	network, err := c.networkSlug(ctx, module)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s://%s.%s.%s", c.cfg.Scheme, network, product, c.cfg.OrgDomain), nil
}

// networkSlug is the configured Network, or the API key's own
// (`GET <org>/api/network`), read once.
func (c *Client) networkSlug(ctx context.Context, module string) (string, error) {
	c.mu.Lock()
	known := c.network
	c.mu.Unlock()
	if known != "" {
		return known, nil
	}
	// Asked with the product's credential: the Network is the credential's.
	res, err := c.Do(ctx, Request{Module: "org", Layer: "network", Method: http.MethodGet, Path: "/api/network"})
	if err != nil {
		return "", fmt.Errorf("a product API goes to a Network host: set network (DAS_NETWORK), or use the Network's API key; reading the credential's Network failed: %w", err)
	}
	m, _ := res.JSON.(map[string]any)
	slug, _ := m["slug"].(string)
	if slug == "" {
		return "", errors.New("GET /api/network answered no Network slug; set network (DAS_NETWORK)")
	}
	c.mu.Lock()
	c.network = strings.ToLower(slug)
	c.mu.Unlock()
	return c.network, nil
}

func problemOf(res *http.Response, data []byte, method, u string) *Problem {
	p := &Problem{Status: res.StatusCode, Method: method, URL: u, RequestID: res.Header.Get("X-Request-Id")}
	var body struct {
		Code      string       `json:"code"`
		Title     string       `json:"title"`
		Detail    string       `json:"detail"`
		RequestID string       `json:"requestId"`
		Errors    []FieldError `json:"errors"`
	}
	if json.Unmarshal(data, &body) == nil {
		p.Code, p.Title, p.Detail, p.Errors = body.Code, body.Title, body.Detail, body.Errors
		if body.RequestID != "" {
			p.RequestID = body.RequestID
		}
	} else {
		p.Detail = strings.TrimSpace(string(data))
		if len(p.Detail) > 500 {
			p.Detail = p.Detail[:500]
		}
	}
	if p.Title == "" {
		p.Title = http.StatusText(res.StatusCode)
	}
	return p
}

// SubstitutePath fills `{param}` segments, escaping each value.
func SubstitutePath(path string, values map[string]string) (string, error) {
	var missing []string
	out := path
	for {
		start := strings.IndexByte(out, '{')
		if start < 0 {
			break
		}
		end := strings.IndexByte(out[start:], '}')
		if end < 0 {
			return "", fmt.Errorf("malformed path %s", path)
		}
		name := out[start+1 : start+end]
		v, ok := values[name]
		if !ok || v == "" {
			missing = append(missing, name)
			v = ""
		}
		out = out[:start] + url.PathEscape(v) + out[start+end+1:]
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return "", fmt.Errorf("%s: no value for %s", path, strings.Join(missing, ", "))
	}
	return out, nil
}

func sleep(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		return ctx.Err() == nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

func newUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// NewIdempotencyKey draws a key for a mutation the caller retries itself.
func NewIdempotencyKey() string { return newUUID() }
