package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

const orgKey = "daso_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

// An Organization API key (#119) administers the Organization; a product
// call never sends it - not even to look up a Network - and fails before
// any request, unless a token is configured for the product.
func TestOrgAPIKeyActsOnTheOrganizationLayerOnly(t *testing.T) {
	var calls atomic.Int32
	var seen atomic.Value
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		seen.Store(r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"slug":"north"}`))
	}))
	defer srv.Close()
	ctx := context.Background()

	c, err := New(Config{APIKey: orgKey, Endpoints: map[string]string{"org": srv.URL}, RetryDelays: []time.Duration{}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Do(ctx, Request{Module: "org", Method: http.MethodGet, Path: "/api/groups"}); err != nil {
		t.Fatal(err)
	}
	if got := seen.Load(); got != "Bearer "+orgKey {
		t.Errorf("the Organisation API got %v", got)
	}

	before := calls.Load()
	_, err = c.Do(ctx, Request{Module: "delivery", Product: "delivery", Method: http.MethodGet, Path: "/api/lineitems"})
	if !errors.Is(err, ErrOrgKeyInProduct) {
		t.Fatalf("a product call with only an Organization API key: %v", err)
	}
	// A Network-layer operation of the Organisation API (its own Network's keys) is refused the same way.
	_, err = c.Do(ctx, Request{Module: "org", Layer: "network", Method: http.MethodGet, Path: "/api/network/api-keys"})
	if !errors.Is(err, ErrOrgKeyInProduct) {
		t.Fatalf("a Network-layer call with only an Organization API key: %v", err)
	}
	if calls.Load() != before {
		t.Error("a Network-layer call with an Organization API key sent a request")
	}

	// With a token as well, the product (and its Network lookup) takes the token.
	product := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen.Store(r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items":[],"nextCursor":null}`))
	}))
	defer product.Close()
	c2, err := New(Config{APIKey: orgKey, Token: "network-token", Endpoints: map[string]string{"org": srv.URL, "delivery": product.URL}, RetryDelays: []time.Duration{}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c2.Do(ctx, Request{Module: "delivery", Product: "delivery", Method: http.MethodGet, Path: "/api/lineitems"}); err != nil {
		t.Fatal(err)
	}
	if got := seen.Load(); got != "Bearer network-token" {
		t.Errorf("the product got %v", got)
	}
	// And so does a Network-layer operation of the Organisation API; an Organization-layer one keeps the key.
	if _, err := c2.Do(ctx, Request{Module: "org", Layer: "network", Method: http.MethodGet, Path: "/api/network/api-keys"}); err != nil {
		t.Fatal(err)
	}
	if got := seen.Load(); got != "Bearer network-token" {
		t.Errorf("a Network-layer call of the Organisation API got %v", got)
	}
	if _, err := c2.Do(ctx, Request{Module: "org", Layer: "org", Method: http.MethodGet, Path: "/api/groups"}); err != nil {
		t.Fatal(err)
	}
	if got := seen.Load(); got != "Bearer "+orgKey {
		t.Errorf("an Organization-layer call got %v", got)
	}
}
