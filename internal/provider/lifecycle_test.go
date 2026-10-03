package provider_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ringdas/terraform-provider-das/internal/provider"
)

// THE ENGINE GATE: the generic resource runs a whole Terraform lifecycle
// (create, an empty second plan, update in place, import, destroy) against
// an in-memory API that answers the Organisation API's contract for three
// shapes of resource - an entity with an id (POST/PATCH/DELETE), a
// singleton set by PUT whose destroy only forgets, and a value keyed by the
// path whose destroy is an update. Every resource of the provider is one of
// these shapes, so a regression of the engine fails here, not in a user's
// apply.

type fakeAPI struct {
	mu       sync.Mutex
	objects  map[string]map[string]any
	seq      int
	requests []string
}

func newFakeAPI() *fakeAPI { return &fakeAPI{objects: map[string]map[string]any{}} }

func (f *fakeAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.requests = append(f.requests, r.Method+" "+r.URL.Path)
	if r.Method != http.MethodGet && r.Header.Get("Idempotency-Key") == "" {
		problem(w, 400, "IDEMPOTENCY_KEY_REQUIRED")
		return
	}
	var body map[string]any
	_ = json.NewDecoder(r.Body).Decode(&body)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	p := r.URL.Path
	switch {
	case r.Method == http.MethodPost && p == "/api/groups":
		f.seq++
		id := fmt.Sprintf("grp_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq)
		obj := map[string]any{"id": id, "name": body["name"], "externalId": body["externalId"], "source": "DAS", "memberCount": 0, "createdBy": "service:test", "createdAt": now, "updatedAt": now}
		f.objects["/api/groups/"+id] = obj
		reply(w, 201, obj)
	case strings.HasPrefix(p, "/api/groups/"):
		obj, ok := f.objects[p]
		if !ok {
			problem(w, 404, "NOT_FOUND")
			return
		}
		switch r.Method {
		case http.MethodGet:
			reply(w, 200, obj)
		case http.MethodPatch:
			for k, v := range body {
				obj[k] = v
			}
			obj["updatedAt"] = now
			reply(w, 200, obj)
		case http.MethodDelete:
			delete(f.objects, p)
			w.WriteHeader(204)
		}
	case p == "/api/identity/sign-in-policy":
		if r.Method == http.MethodPut {
			obj := map[string]any{"localSignIn": body["localSignIn"], "mfa": body["mfa"], "sessionDays": body["sessionDays"], "state": "APPLIED", "appliedAt": now, "error": nil}
			f.objects[p] = obj
			reply(w, 200, obj)
			return
		}
		if obj, ok := f.objects[p]; ok {
			reply(w, 200, obj)
			return
		}
		reply(w, 200, map[string]any{"localSignIn": "ALLOWED", "mfa": "OFF", "sessionDays": 30, "state": "APPLIED", "appliedAt": nil, "error": nil})
	case r.Method == http.MethodPost && p == "/api/scim-tokens":
		// The provider always asks for the secret store (createBody); the
		// answer wraps the entity (createAnswer) next to the secret.
		if body["secretDelivery"] != "SECRET_STORE" {
			problem(w, 400, "VALIDATION_FAILED")
			return
		}
		f.seq++
		id := fmt.Sprintf("scimtok_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq)
		tok := map[string]any{"id": id, "hint": "…abcd", "createdBy": "service:test", "createdAt": now, "expiresAt": nil, "lastUsedAt": nil, "lastRequest": nil, "requestCount": 0, "secretRef": "arn:aws:secretsmanager:eu-central-1:111111111111:secret:das/scim-tokens/" + id + "#v1"}
		f.objects["/api/scim-tokens/"+id] = tok
		reply(w, 201, map[string]any{"token": tok, "secret": nil})
	case r.Method == http.MethodPost && p == "/api/org-api-keys":
		// Always minted into the secret store (createBody); the answer is the key itself, token null.
		if body["secretDelivery"] != "SECRET_STORE" {
			problem(w, 400, "VALIDATION_FAILED")
			return
		}
		f.seq++
		id := fmt.Sprintf("okey_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq)
		key := map[string]any{"id": id, "name": body["name"], "role": body["role"], "createdBy": "human:owner@example.com", "createdAt": now, "expiresAt": nil, "revokedAt": nil, "secretRef": "arn:aws:secretsmanager:eu-central-1:111111111111:secret:das/org-api-keys/" + id + "#v1"}
		f.objects["/api/org-api-keys/"+id] = key
		answer := map[string]any{"token": nil}
		for k, v := range key {
			answer[k] = v
		}
		reply(w, 201, answer)
	case strings.HasPrefix(p, "/api/org-api-keys/"):
		obj, ok := f.objects[p]
		if !ok {
			problem(w, 404, "NOT_FOUND")
			return
		}
		if r.Method == http.MethodDelete {
			// Revoked, never deleted: the row stays, its stored token goes.
			obj["revokedAt"] = now
			obj["secretRef"] = nil
		}
		reply(w, 200, obj)
	case strings.HasPrefix(p, "/api/scim-tokens/"):
		obj, ok := f.objects[p]
		if !ok {
			problem(w, 404, "NOT_FOUND")
			return
		}
		if r.Method == http.MethodDelete {
			delete(f.objects, p)
			w.WriteHeader(204)
			return
		}
		reply(w, 200, obj)
	case strings.HasPrefix(p, "/api/email/senders/"):
		// replyTo is required and nullable: an omitted key is refused, as by the API.
		if r.Method == http.MethodPut {
			if _, ok := body["replyTo"]; !ok {
				problem(w, 400, "VALIDATION_FAILED")
				return
			}
			obj := map[string]any{"purpose": strings.TrimPrefix(p, "/api/email/senders/"), "localPart": body["localPart"], "displayName": body["displayName"], "replyTo": body["replyTo"], "isDefault": false}
			f.objects[p] = obj
			reply(w, 200, obj)
			return
		}
		if r.Method == http.MethodDelete {
			delete(f.objects, p)
		}
		obj, ok := f.objects[p]
		if !ok {
			obj = map[string]any{"purpose": strings.TrimPrefix(p, "/api/email/senders/"), "localPart": "no-reply", "displayName": "DAS", "replyTo": nil, "isDefault": true}
		}
		reply(w, 200, obj)
	case strings.HasPrefix(p, "/api/networks/") && strings.Contains(p, "/features/"):
		parts := strings.Split(p, "/")
		obj, ok := f.objects[p]
		if !ok {
			obj = map[string]any{"networkId": parts[3], "featureId": parts[5], "description": "A feature", "requires": []any{}, "enabled": false, "effective": false, "blockedBy": []any{}, "updatedBy": nil, "updatedAt": nil}
		}
		if r.Method == http.MethodPut {
			obj["enabled"] = body["enabled"]
			obj["effective"] = body["enabled"]
			obj["updatedBy"] = "service:test"
			obj["updatedAt"] = now
			f.objects[p] = obj
		}
		reply(w, 200, obj)
	default:
		problem(w, 404, "NOT_FOUND")
	}
}

func reply(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func problem(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"status": status, "code": code, "title": code})
}

func config(url, groupName string, sessionDays int, enabled bool) string {
	return fmt.Sprintf(`
provider "das" {
  token     = "test"
  endpoints = { org = %q }
}

resource "das_group" "g" {
  name = %q
}

resource "das_sign_in_policy" "p" {
  local_sign_in = "IDP_ONLY"
  mfa           = "REQUIRED"
  session_days  = %d
}

resource "das_network_feature" "f" {
  network_id = "net_01M40KW4H978FZ2GVQYS00DSQH"
  feature_id = "lineitem-summary"
  enabled    = %t
}

# Minted into the secret store: only its reference reaches state.
resource "das_scim_token" "t" {}

resource "das_org_api_key" "ci" {
  name = "ci"
  role = "ORG_ADMIN"
}

# reply_to left out: an explicit null on the wire.
resource "das_email_sender" "s" {
  purpose      = "platform.notifications"
  local_part   = "notifications"
  display_name = "DAS Notifications"
}
`, url, groupName, sessionDays, enabled)
}

// terraformVersion is the Terraform the gate runs on, as in CI
// (.github/workflows/ci.yml, hashicorp/setup-terraform): a version manager
// (tenv) on a developer's machine picks it too, whatever its default.
const terraformVersion = "1.15.3"

func TestLifecycleAgainstTheContract(t *testing.T) {
	if os.Getenv("TFENV_TERRAFORM_VERSION") == "" {
		t.Setenv("TFENV_TERRAFORM_VERSION", terraformVersion)
	}
	api := newFakeAPI()
	server := httptest.NewServer(api)
	defer server.Close()
	factories := map[string]func() (tfprotov6.ProviderServer, error){
		"das": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{
				Config: config(server.URL, "traffickers", 30, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("das_group.g", "id"),
					resource.TestCheckResourceAttr("das_group.g", "source", "DAS"),
					resource.TestCheckResourceAttr("das_sign_in_policy.p", "id", "sign_in_policy"),
					resource.TestCheckResourceAttr("das_sign_in_policy.p", "state", "APPLIED"),
					resource.TestCheckResourceAttr("das_network_feature.f", "id", "net_01M40KW4H978FZ2GVQYS00DSQH/lineitem-summary"),
					resource.TestCheckResourceAttr("das_network_feature.f", "effective", "true"),
					resource.TestCheckNoResourceAttr("das_email_sender.s", "reply_to"),
					resource.TestMatchResourceAttr("das_scim_token.t", "id", regexp.MustCompile(`^scimtok_`)),
					resource.TestMatchResourceAttr("das_scim_token.t", "secret_ref", regexp.MustCompile(`^arn:aws:secretsmanager:`)),
					resource.TestCheckNoResourceAttr("das_scim_token.t", "secret"),
					resource.TestCheckNoResourceAttr("das_scim_token.t", "secret_delivery"),
					resource.TestMatchResourceAttr("das_org_api_key.ci", "id", regexp.MustCompile(`^okey_`)),
					resource.TestMatchResourceAttr("das_org_api_key.ci", "secret_ref", regexp.MustCompile(`^arn:aws:secretsmanager:`)),
					resource.TestCheckNoResourceAttr("das_org_api_key.ci", "token"),
					resource.TestCheckNoResourceAttr("das_org_api_key.ci", "secret_delivery"),
				),
			},
			{
				// In place: PATCH the group, PUT the policy and the feature.
				Config: config(server.URL, "traffickers-renamed", 31, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("das_group.g", "name", "traffickers-renamed"),
					resource.TestCheckResourceAttr("das_sign_in_policy.p", "session_days", "31"),
					resource.TestCheckResourceAttr("das_network_feature.f", "enabled", "false"),
				),
			},
			{
				// On again, so destroy has something to switch off.
				Config: config(server.URL, "traffickers-renamed", 31, true),
				Check:  resource.TestCheckResourceAttr("das_network_feature.f", "enabled", "true"),
			},
			{ResourceName: "das_group.g", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_network_feature.f", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_sign_in_policy.p", ImportState: true, ImportStateId: "sign_in_policy", ImportStateVerify: true},
			{ResourceName: "das_scim_token.t", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_org_api_key.ci", ImportState: true, ImportStateVerify: true},
		},
		CheckDestroy: func(*terraform.State) error {
			api.mu.Lock()
			defer api.mu.Unlock()
			for p := range api.objects {
				if strings.HasPrefix(p, "/api/groups/") {
					return fmt.Errorf("the group was not deleted: %s", p)
				}
			}
			if f := api.objects["/api/networks/net_01M40KW4H978FZ2GVQYS00DSQH/features/lineitem-summary"]; f["enabled"] != false {
				return fmt.Errorf("destroy did not switch the feature off: %v", f)
			}
			for p, o := range api.objects {
				if strings.HasPrefix(p, "/api/org-api-keys/") && o["revokedAt"] == nil {
					return fmt.Errorf("destroy did not revoke the Organization API key %s", p)
				}
			}
			if _, ok := api.objects["/api/identity/sign-in-policy"]; !ok {
				return fmt.Errorf("destroy must leave the sign-in policy as it is")
			}
			return nil
		},
	})
	patched := false
	for _, r := range api.requests {
		if strings.HasPrefix(r, "PATCH /api/groups/") {
			patched = true
		}
	}
	if !patched {
		t.Errorf("the rename was not an update in place: %v", api.requests)
	}
}
