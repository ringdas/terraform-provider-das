package provider_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

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
	// keptAttachments counts attachments a PATCH of a policy kept by configuration (their ids stay).
	keptAttachments    int
	patchedAttachments int
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
	case r.Method == http.MethodPost && p == "/api/sales-channels":
		if refusesNull(w, body, "description") {
			return
		}
		f.seq++
		id := fmt.Sprintf("sch_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq)
		obj := map[string]any{"id": id, "key": body["key"], "name": body["name"], "description": body["description"], "status": "DISABLED", "adProductIds": orEmpty(body["adProductIds"]), "apiKeyIds": orEmpty(body["apiKeyIds"]), "createdAt": now, "updatedAt": now, "archivedAt": nil}
		f.objects["/api/sales-channels/"+id] = obj
		reply(w, 201, obj)
	case strings.HasPrefix(p, "/api/sales-channels/"):
		// A sales channel (GET, PATCH): archivable, never deleted; a list given replaces the stored one.
		obj, ok := f.objects[p]
		if !ok || (r.Method != http.MethodGet && r.Method != http.MethodPatch) {
			problem(w, 404, "NOT_FOUND")
			return
		}
		if r.Method == http.MethodPatch {
			if _, ok := body["key"]; ok {
				problem(w, 400, "VALIDATION_FAILED")
				return
			}
			for k, v := range body {
				obj[k] = v
			}
			obj["updatedAt"] = now
		}
		reply(w, 200, obj)
	case r.Method == http.MethodPost && p == "/api/policy-types":
		if refusesNull(w, body, "description") {
			return
		}
		f.seq++
		id := fmt.Sprintf("poltype_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq)
		obj := map[string]any{"id": id, "key": body["key"], "label": body["label"], "description": body["description"], "additive": orFalse(body["additive"]), "contexts": body["contexts"], "status": orString(body["status"], "ACTIVE"), "platforms": body["platforms"], "dimensions": orEmpty(body["dimensions"]), "createdAt": now, "createdBy": "service:test", "updatedAt": now, "updatedBy": "service:test", "archivedAt": nil}
		obj["configurationTypes"] = f.configurationTypes(body["configurationTypes"])
		f.objects["/api/policy-types/"+id] = obj
		reply(w, 201, obj)
	case strings.HasPrefix(p, "/api/policy-types/"):
		// A policy type (GET, PATCH): archivable, never deleted; its key never changes.
		obj, ok := f.objects[p]
		if !ok || (r.Method != http.MethodGet && r.Method != http.MethodPatch) {
			problem(w, 404, "NOT_FOUND")
			return
		}
		if r.Method == http.MethodPatch {
			if _, ok := body["key"]; ok {
				problem(w, 400, "VALIDATION_FAILED")
				return
			}
			if _, ok := body["expectedUpdatedAt"]; ok {
				problem(w, 400, "VALIDATION_FAILED") // the provider never sends it (x-das-terraform-skip)
				return
			}
			for k, v := range body {
				obj[k] = v
			}
			if cts, ok := body["configurationTypes"]; ok {
				obj["configurationTypes"] = f.configurationTypes(cts)
			}
			obj["archivedAt"] = nil
			if obj["status"] == "ARCHIVED" {
				obj["archivedAt"] = now
			}
			obj["updatedAt"] = now
		}
		reply(w, 200, obj)
	case r.Method == http.MethodPost && p == "/api/policy-target-dimensions":
		f.seq++
		id := fmt.Sprintf("dim_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq)
		weight := body["weight"]
		if weight == nil {
			weight = 5
		}
		obj := map[string]any{"id": id, "key": body["key"], "label": body["label"], "input": body["input"], "weight": weight}
		f.objects["/api/policy-target-dimensions/"+id] = obj
		reply(w, 201, obj)
	case r.Method == http.MethodPost && p == "/api/policy-site-platforms":
		f.seq++
		id := fmt.Sprintf("spmap_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq)
		obj := map[string]any{"id": id, "position": body["position"], "pattern": body["pattern"], "platform": body["platform"]}
		f.objects["/api/policy-site-platforms/"+id] = obj
		reply(w, 201, obj)
	case r.Method == http.MethodPost && p == "/api/policy-labels":
		f.seq++
		id := fmt.Sprintf("label_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq)
		obj := map[string]any{"id": id, "name": body["name"]}
		f.objects["/api/policy-labels/"+id] = obj
		reply(w, 201, obj)
	case r.Method == http.MethodPost && p == "/api/policy-configurations":
		if refusesNull(w, body, "description") {
			return
		}
		f.seq++
		id := fmt.Sprintf("cfg_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq)
		obj := map[string]any{
			"id": id, "policyTypeId": body["policyTypeId"], "configurationTypeId": body["configurationTypeId"], "name": body["name"], "status": orString(body["status"], "ACTIVE"),
			"description": body["description"], "body": body["body"], "usage": map[string]any{"policies": 0, "activePolicies": 0},
			"createdAt": now, "createdBy": "service:test", "updatedAt": now, "updatedBy": "service:test", "archivedAt": nil,
		}
		f.objects["/api/policy-configurations/"+id] = obj
		reply(w, 201, obj)
	case strings.HasPrefix(p, "/api/policy-configurations/"):
		// A configuration (GET, PATCH): archivable, never deleted; its policy type and configuration type never change.
		obj, ok := f.objects[p]
		if !ok || (r.Method != http.MethodGet && r.Method != http.MethodPatch) {
			problem(w, 404, "NOT_FOUND")
			return
		}
		if r.Method == http.MethodPatch {
			_, typed := body["policyTypeId"]
			_, configurationTyped := body["configurationTypeId"]
			if typed || configurationTyped {
				problem(w, 400, "VALIDATION_FAILED")
				return
			}
			delete(body, "expectedUpdatedAt")
			for k, v := range body {
				obj[k] = v
			}
			obj["archivedAt"] = nil
			if obj["status"] == "ARCHIVED" {
				obj["archivedAt"] = now
			}
			obj["updatedAt"] = now
		}
		reply(w, 200, obj)
	case r.Method == http.MethodPost && p == "/api/policies":
		f.seq++
		id := fmt.Sprintf("pol_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq)
		obj := map[string]any{
			"id": id, "policyTypeId": body["policyTypeId"], "name": body["name"], "status": orString(body["status"], "ACTIVE"), "context": orString(body["context"], "PAGE"),
			"startAt": body["startAt"], "endAt": body["endAt"], "description": body["description"],
			"labels": labelRows(body["labels"]), "inventory": inventoryRows(body["inventory"]), "target": targetRows(body["target"]), "attachments": []any{},
			"createdAt": now, "createdBy": "service:test", "updatedAt": now, "updatedBy": "service:test", "archivedAt": nil,
		}
		obj["attachments"] = f.attachments(nil, body["attachments"])
		f.objects["/api/policies/"+id] = obj
		reply(w, 201, obj)
	case strings.HasPrefix(p, "/api/policies/"):
		// A policy (GET, PATCH): archivable, never deleted; its type never changes; a section sent replaces that part.
		obj, ok := f.objects[p]
		if !ok || (r.Method != http.MethodGet && r.Method != http.MethodPatch) {
			problem(w, 404, "NOT_FOUND")
			return
		}
		if r.Method == http.MethodPatch {
			if _, ok := body["policyTypeId"]; ok {
				problem(w, 400, "VALIDATION_FAILED")
				return
			}
			if _, ok := body["expectedUpdatedAt"]; ok {
				problem(w, 400, "VALIDATION_FAILED") // the provider never sends it (x-das-terraform-skip)
				return
			}
			for k, v := range body {
				switch k {
				case "labels":
					obj[k] = labelRows(v)
				case "inventory":
					obj[k] = inventoryRows(v)
				case "target":
					obj[k] = targetRows(v)
				case "attachments":
					obj[k] = f.attachments(obj["attachments"], v)
					f.patchedAttachments++
				default:
					obj[k] = v
				}
			}
			obj["archivedAt"] = nil
			if obj["status"] == "ARCHIVED" {
				obj["archivedAt"] = now
			}
			obj["updatedAt"] = now
		}
		reply(w, 200, obj)
	case strings.HasPrefix(p, "/api/policy-target-dimensions/"), strings.HasPrefix(p, "/api/policy-site-platforms/"), strings.HasPrefix(p, "/api/policy-labels/"):
		// A target dimension, a site-platform rule or a label (GET, PATCH, DELETE); a dimension's key never changes.
		obj, ok := f.objects[p]
		if !ok {
			problem(w, 404, "NOT_FOUND")
			return
		}
		switch r.Method {
		case http.MethodDelete:
			delete(f.objects, p)
			w.WriteHeader(204)
			return
		case http.MethodPatch:
			if _, ok := body["key"]; ok {
				problem(w, 400, "VALIDATION_FAILED")
				return
			}
			for k, v := range body {
				obj[k] = v
			}
		}
		reply(w, 200, obj)
	case r.Method == http.MethodPost && p == "/api/policy-catalogues":
		// A catalogue: its key never changes; a value's label left out is null.
		f.seq++
		id := fmt.Sprintf("cat_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq)
		obj := map[string]any{"id": id, "key": body["key"], "label": body["label"], "values": catalogueValues(body["values"]), "createdAt": now, "updatedAt": now}
		f.objects["/api/policy-catalogues/"+id] = obj
		reply(w, 201, obj)
	case r.Method == http.MethodPost && p == "/api/policy-mappings":
		// A mapping names its catalogues by key; a cell without a value has no `value`.
		for _, side := range []string{"from", "to"} {
			if f.catalogueByKey(body[side]) == nil {
				problem(w, 400, "VALIDATION_FAILED")
				return
			}
		}
		for _, o := range f.objects {
			if o["key"] == body["key"] && strings.HasPrefix(fmt.Sprint(o["id"]), "map_") {
				problem(w, 409, "ALREADY_EXISTS")
				return
			}
		}
		f.seq++
		id := fmt.Sprintf("map_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq)
		obj := map[string]any{"id": id, "key": body["key"], "label": body["label"], "from": body["from"], "to": body["to"], "cell": body["cell"], "rows": orEmpty(body["rows"]), "createdAt": now, "updatedAt": now}
		f.objects["/api/policy-mappings/"+id] = obj
		reply(w, 201, obj)
	case r.Method == http.MethodGet && p == "/api/policy-option-sources":
		// The sources of one kind of field, `kind` required: `map` the mappings (by key); `string` the
		// catalogues (by key). One source a page, so the data source must follow `nextCursor` and
		// send `kind` again with every cursor.
		kind := r.URL.Query().Get("kind")
		if kind != "map" && kind != "string" {
			problem(w, 400, "VALIDATION_FAILED")
			return
		}
		prefix, source, entries, entryKey := "/api/policy-mappings/", "mapping.", "rows", "key"
		if kind == "string" {
			prefix, source, entries, entryKey = "/api/policy-catalogues/", "catalogue.", "values", "value"
		}
		items := []any{}
		for k, o := range f.objects {
			if !strings.HasPrefix(k, prefix) {
				continue
			}
			list := orEmpty(o[entries]).([]any)
			sample := []any{}
			for _, e := range list[:min(3, len(list))] {
				sample = append(sample, e.(map[string]any)[entryKey])
			}
			items = append(items, map[string]any{"key": source + fmt.Sprint(o["key"]), "label": o["label"], "count": len(list), "sample": sample})
		}
		sort.Slice(items, func(i, j int) bool {
			return fmt.Sprint(items[i].(map[string]any)["key"]) < fmt.Sprint(items[j].(map[string]any)["key"])
		})
		at := 0
		if c := r.URL.Query().Get("cursor"); c != "" {
			n, err := strconv.Atoi(c)
			if err != nil || n < 1 || n >= len(items) {
				problem(w, 400, "INVALID_CURSOR")
				return
			}
			at = n
		}
		var next any
		if at+1 < len(items) {
			next = strconv.Itoa(at + 1)
		}
		reply(w, 200, map[string]any{"items": items[at:min(at+1, len(items))], "nextCursor": next})
	case strings.HasPrefix(p, "/api/policy-catalogues/"), strings.HasPrefix(p, "/api/policy-mappings/"):
		// A catalogue or a mapping (GET, PATCH, DELETE); a key, a mapping's catalogues and its cell never change.
		obj, ok := f.objects[p]
		if !ok {
			problem(w, 404, "NOT_FOUND")
			return
		}
		switch r.Method {
		case http.MethodDelete:
			if strings.HasPrefix(p, "/api/policy-catalogues/") {
				for _, o := range f.objects {
					if strings.HasPrefix(fmt.Sprint(o["id"]), "map_") && (o["from"] == obj["key"] || o["to"] == obj["key"]) {
						problem(w, 409, "IN_USE") // a mapping still names it: destroy removes the mapping first
						return
					}
				}
			}
			delete(f.objects, p)
			w.WriteHeader(204)
			return
		case http.MethodPatch:
			for _, fixed := range []string{"key", "from", "to", "cell", "expectedUpdatedAt"} {
				if _, ok := body[fixed]; ok {
					problem(w, 400, "VALIDATION_FAILED")
					return
				}
			}
			for k, v := range body {
				if k == "values" {
					v = catalogueValues(v)
				}
				obj[k] = v
			}
			obj["updatedAt"] = now
		}
		reply(w, 200, obj)
	case r.Method == http.MethodPost && p == "/api/ad-products":
		if refusesNull(w, body, "description") {
			return
		}
		f.seq++
		id := fmt.Sprintf("adp_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq)
		obj := map[string]any{"id": id, "key": body["key"], "name": body["name"], "description": body["description"], "status": "DISABLED", "currentVersion": nil, "createdAt": now, "updatedAt": now, "archivedAt": nil}
		f.objects["/api/ad-products/"+id] = obj
		reply(w, 201, obj)
	case r.Method == http.MethodPost && strings.HasPrefix(p, "/api/ad-products/") && strings.HasSuffix(p, "/versions"):
		// The next number, which becomes the ad product's current version.
		productPath := strings.TrimSuffix(p, "/versions")
		product, ok := f.objects[productPath]
		if !ok {
			problem(w, 404, "NOT_FOUND")
			return
		}
		n := 1
		if current, ok := product["currentVersion"].(int); ok {
			n = current + 1
		}
		f.seq++
		version := map[string]any{"id": fmt.Sprintf("adpv_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq), "adProductId": strings.TrimPrefix(productPath, "/api/ad-products/"), "version": n, "definition": body["definition"], "createdAt": now}
		product["currentVersion"] = n
		f.objects[fmt.Sprintf("%s/versions/%d", productPath, n)] = version
		reply(w, 201, version)
	case strings.HasPrefix(p, "/api/ad-products/"):
		// An ad product (GET, PATCH) or one of its versions (GET only: immutable).
		obj, ok := f.objects[p]
		if !ok || (r.Method != http.MethodGet && strings.Contains(p, "/versions/")) {
			problem(w, 404, "NOT_FOUND")
			return
		}
		if r.Method == http.MethodPatch {
			// As the API: an ad product is enabled only once it has a version.
			if body["status"] == "ENABLED" && obj["currentVersion"] == nil {
				problem(w, 400, "VALIDATION_FAILED")
				return
			}
			for k, v := range body {
				obj[k] = v
			}
			obj["updatedAt"] = now
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

// configurationTypes are a policy type's configuration types as the API
// stores them: each with an id, its label null when left out.
// catalogueValues answers a catalogue's values as the API does: a label left out is null.
func catalogueValues(v any) []any {
	out := []any{}
	for _, it := range orEmpty(v).([]any) {
		m := it.(map[string]any)
		out = append(out, map[string]any{"value": m["value"], "label": m["label"]})
	}
	return out
}

func (f *fakeAPI) catalogueByKey(key any) map[string]any {
	for k, o := range f.objects {
		if strings.HasPrefix(k, "/api/policy-catalogues/") && o["key"] == key {
			return o
		}
	}
	return nil
}

func (f *fakeAPI) configurationTypes(v any) []any {
	in, _ := v.([]any)
	out := make([]any, 0, len(in))
	for _, c := range in {
		m, _ := c.(map[string]any)
		f.seq++
		out = append(out, map[string]any{"id": fmt.Sprintf("cfgtype_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq), "key": m["key"], "label": m["label"], "fields": m["fields"]})
	}
	return out
}

// The parts of a policy as the API answers them: defaults filled in and in
// the API's canonical order, whatever order the request gave.

func labelRows(v any) []any {
	out := append([]any{}, orEmpty(v).([]any)...)
	sort.SliceStable(out, func(i, j int) bool {
		return fmt.Sprint(out[i].(map[string]any)["name"]) < fmt.Sprint(out[j].(map[string]any)["name"])
	})
	return out
}

func inventoryRows(v any) []any {
	in, _ := v.([]any)
	out := make([]any, 0, len(in))
	for _, x := range in {
		m, _ := x.(map[string]any)
		out = append(out, map[string]any{"path": m["path"], "exclude": orFalse(m["exclude"])})
	}
	key := func(m map[string]any) string {
		return fmt.Sprint(m["exclude"], m["path"])
	}
	sort.SliceStable(out, func(i, j int) bool { return key(out[i].(map[string]any)) < key(out[j].(map[string]any)) })
	return out
}

func targetRows(v any) []any {
	in, _ := v.([]any)
	out := make([]any, 0, len(in))
	for _, x := range in {
		m, _ := x.(map[string]any)
		out = append(out, map[string]any{"targetGroup": m["targetGroup"], "dimension": m["dimension"], "op": m["op"], "values": m["values"]})
	}
	key := func(m map[string]any) string {
		values, _ := json.Marshal(m["values"])
		return fmt.Sprint(m["targetGroup"], "|", m["dimension"], "|", m["op"], "|", string(values))
	}
	sort.SliceStable(out, func(i, j int) bool { return key(out[i].(map[string]any)) < key(out[j].(map[string]any)) })
	return out
}

// attachments are a policy's attachments in the order given, each keeping
// the id it had for the same configuration.
func (f *fakeAPI) attachments(stored any, v any) []any {
	ids := map[any]any{}
	if old, ok := stored.([]any); ok {
		for _, x := range old {
			m := x.(map[string]any)
			ids[m["configurationId"]] = m["id"]
		}
	}
	in, _ := v.([]any)
	out := make([]any, 0, len(in))
	for i, x := range in {
		m, _ := x.(map[string]any)
		id, kept := ids[m["configurationId"]]
		if kept {
			f.keptAttachments++
		} else {
			f.seq++
			id = fmt.Sprintf("polcfg_01M40KW4GB68ZT4D4HAE7J3X%02d", f.seq)
		}
		enabled := m["enabled"]
		if enabled == nil {
			enabled = true
		}
		out = append(out, map[string]any{"id": id, "configurationId": m["configurationId"], "position": i + 1, "enabled": enabled, "target": targetRows(m["target"]), "inventory": inventoryRows(m["inventory"])})
	}
	return out
}

// refusesNull answers 400 to an explicit null the create takes only as left
// out, as the API (its description is nullable in the update alone).
func refusesNull(w http.ResponseWriter, body map[string]any, key string) bool {
	if v, ok := body[key]; ok && v == nil {
		problem(w, 400, "VALIDATION_FAILED")
		return true
	}
	return false
}

func orFalse(v any) any {
	if v == nil {
		return false
	}
	return v
}

func orString(v any, def string) any {
	if v == nil {
		return def
	}
	return v
}

// orEmpty is a list the request left out, as the API stores it: empty.
func orEmpty(v any) any {
	if v == nil {
		return []any{}
	}
	return v
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

// config: the ad product's status is its own, as it can be ENABLED only in
// an apply after the one that added its first version.
func config(url, groupName string, sessionDays int, enabled bool, adProductStatus string) string {
	// The policy's attachments, in an order that follows sessionDays: a reorder between applies.
	first := `{ configuration_id = das_tag_manager_policy_configuration.cfg.id }`
	second := `{ configuration_id = das_tag_manager_policy_configuration.cfg2.id, enabled = false, target = [{ target_group = 1, dimension = "seg_ab", op = "IN", values = [jsonencode(7)] }], inventory = [{ path = "NEWS_APP_IOS/*/right" }] }`
	if sessionDays%2 == 1 {
		first, second = second, first
	}
	attachments := "[\n    " + first + ",\n    " + second + ",\n  ]"
	return fmt.Sprintf(`
provider "das" {
  token     = "test"
  endpoints = { org = %[1]q, delivery = %[1]q, "tag-manager" = %[1]q }
}

# Archivable: destroy switches it off. Renamed with the group.
resource "das_delivery_ad_product" "a" {
  key    = "homepage-takeover"
  name   = "Homepage takeover for %[2]s"
  status = %[5]q
}

# Immutable, numbered by the API: identity <adProductId>/<version>.
resource "das_delivery_ad_product_version" "v" {
  ad_product_id = das_delivery_ad_product.a.id
  definition = {
    goals           = ["AWARENESS"]
    billing_models  = ["CPM", "FLAT_FEE"]
    creative_source = "SUPPLIED"
    steps = [
      { key = "brief", name = "Take the brief", kind = "WORK", role = "OPERATOR" },
      { key = "approve", name = "Approve the plan", kind = "APPROVAL", role = "APPROVER" },
    ]
  }
}

# Archivable: destroy switches it off and keeps its links. The offer and
# the keys are lists replaced whole.
resource "das_delivery_sales_channel" "c" {
  key            = "onetads"
  name           = "OnetAds for %[2]s"
  status         = %[5]q
  ad_product_ids = [das_delivery_ad_product.a.id]
  api_key_ids    = ["key_01M40KW4GB68ZT4D4HAE7J3X99"]
}

# Archivable: destroy archives it. Configuration types are a list of
# objects replaced whole; each field definition is one jsonencode().
# Its dimensions are a set of keys, the dimension's own key referenced.
resource "das_tag_manager_policy_type" "pt" {
  key        = "prebid_timeout_app"
  label      = "Prebid timeout for %[2]s"
  contexts   = ["PAGE"]
  platforms  = ["APP"]
  dimensions = [das_tag_manager_policy_target_dimension.d.key]
  configuration_types = [{
    key    = "prebid_timeout_app"
    fields = [jsonencode({ id = "timeoutMs", kind = "int", min = 100 })]
  }]
}

# Deleted on destroy; its input is a value of no fixed shape.
resource "das_tag_manager_policy_target_dimension" "d" {
  key   = "seg_ab"
  label = "A/B for %[2]s"
  input = { kind = "int", min = 0, max = 99 }
}

resource "das_tag_manager_policy_site_platform" "s" {
  position = %[3]d
  pattern  = "*_APP_*"
  platform = "APP"
}

# Deleted on destroy. Values are a list kept in the order written, a label
# left out (null in the answer) included; renamed in place.
resource "das_tag_manager_policy_catalogue" "sizes" {
  key    = "sizes_app"
  label  = "Sizes for %[2]s"
  values = [{ value = "320x50" }, { value = "native", label = "Native" }, { value = "300x250" }]
}

resource "das_tag_manager_policy_catalogue" "consumers" {
  key    = "size_consumers"
  label  = "Size consumers"
  values = [{ value = "PREBID", label = "Prebid" }, { value = "GAM" }]
}

# Membership: no cell type, cells without value_json. Rows and cells keep
# the order written; renamed in place.
resource "das_tag_manager_policy_mapping" "map" {
  key   = "size_map_app"
  label = "Size map for %[2]s"
  from  = das_tag_manager_policy_catalogue.sizes.key
  to    = das_tag_manager_policy_catalogue.consumers.key
  rows = [
    { key = "native", cells = [{ to = "PREBID" }] },
    { key = "320x50", cells = [{ to = "PREBID" }, { to = "GAM" }] },
  ]
}

# A cell type (decimal): each cell's value is one jsonencode(). The cell type
# follows the session days, so the second apply replaces the mapping.
resource "das_tag_manager_policy_mapping" "floors" {
  key   = "size_floor_app"
  label = "Size floors"
  from  = das_tag_manager_policy_catalogue.sizes.key
  to    = das_tag_manager_policy_catalogue.consumers.key
  cell  = { kind = "decimal", min = 0, max = %[3]d, unit = "EUR" }
  rows = [
    { key = "300x250", cells = [{ to = "GAM", value_json = jsonencode(0.4) }, { to = "PREBID", value_json = jsonencode(0.35) }] },
  ]
}

# Read only: the sources a map field of a policy type may read.
data "das_tag_manager_policy_option_sources" "maps" {
  kind       = "map"
  depends_on = [das_tag_manager_policy_mapping.map, das_tag_manager_policy_mapping.floors]
}

# Read only: the sources a string field may read: the catalogues.
data "das_tag_manager_policy_option_sources" "strings" {
  kind       = "string"
  depends_on = [das_tag_manager_policy_mapping.map, das_tag_manager_policy_mapping.floors]
}

# Deleted on destroy, and with it removed from its policies.
resource "das_tag_manager_policy_label" "l" {
  name = "floorian-%[2]s"
}

# Archivable: destroy archives it. Its body is one jsonencode(); its policy
# type and configuration type never change.
resource "das_tag_manager_policy_configuration" "cfg" {
  policy_type_id        = das_tag_manager_policy_type.pt.id
  configuration_type_id = "cfgtype_01M40KW4GB68ZT4D4HAE7J3X99"
  name                  = "Standard for %[2]s"
  body_json             = jsonencode({ timeoutMs = 1500 + %[3]d })
}

resource "das_tag_manager_policy_configuration" "cfg2" {
  policy_type_id        = das_tag_manager_policy_type.pt.id
  configuration_type_id = "cfgtype_01M40KW4GB68ZT4D4HAE7J3X99"
  name                  = "Slow"
  body_json             = jsonencode({ timeoutMs = 2500 })
}

resource "das_tag_manager_policy_label" "l2" {
  name = "aardvark-%[2]s"
}

# Archivable: destroy archives it. Labels, inventory and targets are sets:
# written here out of the API's canonical order and without their defaults
# (exclude, enabled, an attachment's empty target and inventory), and still
# an empty plan; a row's values keep the order written (5, 10, 1). The
# attachments are an ordered list, swapped between the first and the second
# apply.
resource "das_tag_manager_policy" "p" {
  policy_type_id = das_tag_manager_policy_type.pt.id
  name           = "Run of network for %[2]s"
  labels = [
    { id = das_tag_manager_policy_label.l.id, name = das_tag_manager_policy_label.l.name },
    { id = das_tag_manager_policy_label.l2.id, name = das_tag_manager_policy_label.l2.name },
  ]
  inventory = [{ path = "NEWS_APP_IOS/home", exclude = true }, { path = "*" }]
  target = [
    { target_group = 2, dimension = "app_version", op = "IN", values = [jsonencode("4.0.%[3]d")] },
    { target_group = 1, dimension = "seg_ab", op = "BETWEEN", values = [jsonencode(0), jsonencode(49)] },
    { target_group = 1, dimension = "seg_ab", op = "IN", values = [jsonencode(5), jsonencode(10), jsonencode(1)] },
  ]
  attachments = %[6]s
}

resource "das_group" "g" {
  name = %[2]q
}

resource "das_sign_in_policy" "p" {
  local_sign_in = "IDP_ONLY"
  mfa           = "REQUIRED"
  session_days  = %[3]d
}

resource "das_network_feature" "f" {
  network_id = "net_01M40KW4H978FZ2GVQYS00DSQH"
  feature_id = "lineitem-summary"
  enabled    = %[4]t
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
`, url, groupName, sessionDays, enabled, adProductStatus, attachments)
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
				Config: config(server.URL, "traffickers", 30, true, "DISABLED"),
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
					resource.TestMatchResourceAttr("das_delivery_ad_product.a", "id", regexp.MustCompile(`^adp_`)),
					resource.TestCheckResourceAttr("das_delivery_ad_product.a", "status", "DISABLED"),
					resource.TestCheckResourceAttr("das_delivery_ad_product_version.v", "version", "1"),
					resource.TestMatchResourceAttr("das_delivery_ad_product_version.v", "id", regexp.MustCompile(`^adp_[0-9A-Z]+/1$`)),
					resource.TestCheckResourceAttr("das_delivery_ad_product_version.v", "definition.steps.1.role", "APPROVER"),
					resource.TestMatchResourceAttr("das_delivery_sales_channel.c", "id", regexp.MustCompile(`^sch_`)),
					resource.TestCheckResourceAttr("das_delivery_sales_channel.c", "status", "DISABLED"),
					resource.TestMatchResourceAttr("das_delivery_sales_channel.c", "ad_product_ids.0", regexp.MustCompile(`^adp_`)),
					resource.TestCheckResourceAttr("das_delivery_sales_channel.c", "api_key_ids.0", "key_01M40KW4GB68ZT4D4HAE7J3X99"),
					resource.TestMatchResourceAttr("das_tag_manager_policy_type.pt", "id", regexp.MustCompile(`^poltype_`)),
					resource.TestCheckResourceAttr("das_tag_manager_policy_type.pt", "status", "ACTIVE"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_type.pt", "configuration_types.0.fields.0", `{"id":"timeoutMs","kind":"int","min":100}`),
					resource.TestCheckTypeSetElemAttr("das_tag_manager_policy_type.pt", "dimensions.*", "seg_ab"),
					resource.TestMatchResourceAttr("das_tag_manager_policy_target_dimension.d", "id", regexp.MustCompile(`^dim_`)),
					resource.TestCheckResourceAttr("das_tag_manager_policy_target_dimension.d", "weight", "5"),
					resource.TestMatchResourceAttr("das_tag_manager_policy_site_platform.s", "id", regexp.MustCompile(`^spmap_`)),
					resource.TestMatchResourceAttr("das_tag_manager_policy_catalogue.sizes", "id", regexp.MustCompile(`^cat_`)),
					resource.TestCheckResourceAttr("das_tag_manager_policy_catalogue.sizes", "values.1.value", "native"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_catalogue.sizes", "values.1.label", "Native"),
					resource.TestCheckNoResourceAttr("das_tag_manager_policy_catalogue.sizes", "values.0.label"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_catalogue.sizes", "values.2.value", "300x250"),
					resource.TestMatchResourceAttr("das_tag_manager_policy_mapping.map", "id", regexp.MustCompile(`^map_`)),
					resource.TestCheckResourceAttr("das_tag_manager_policy_mapping.map", "from", "sizes_app"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_mapping.map", "rows.0.key", "native"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_mapping.map", "rows.1.cells.0.to", "PREBID"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_mapping.map", "rows.1.cells.1.to", "GAM"),
					resource.TestCheckNoResourceAttr("das_tag_manager_policy_mapping.map", "rows.1.cells.0.value_json"),
					resource.TestCheckNoResourceAttr("das_tag_manager_policy_mapping.map", "cell"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_mapping.floors", "rows.0.cells.0.value_json", "0.4"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_mapping.floors", "rows.0.cells.1.value_json", "0.35"),
					resource.TestCheckResourceAttr("data.das_tag_manager_policy_option_sources.maps", "items.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("data.das_tag_manager_policy_option_sources.maps", "items.*", map[string]string{"key": "mapping.size_map_app", "count": "2", "sample.0": "native", "sample.1": "320x50"}),
					resource.TestCheckResourceAttr("data.das_tag_manager_policy_option_sources.strings", "items.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("data.das_tag_manager_policy_option_sources.strings", "items.*", map[string]string{"key": "catalogue.sizes_app", "count": "3", "sample.0": "320x50", "sample.2": "300x250"}),
					resource.TestMatchResourceAttr("das_tag_manager_policy_label.l", "id", regexp.MustCompile(`^label_`)),
					resource.TestMatchResourceAttr("das_tag_manager_policy_configuration.cfg", "id", regexp.MustCompile(`^cfg_`)),
					resource.TestCheckResourceAttr("das_tag_manager_policy_configuration.cfg", "status", "ACTIVE"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_configuration.cfg", "body_json", `{"timeoutMs":1530}`),
					resource.TestCheckResourceAttr("das_tag_manager_policy_configuration.cfg", "usage.policies", "0"),
					resource.TestCheckResourceAttrPair("das_tag_manager_policy_configuration.cfg", "policy_type_id", "das_tag_manager_policy_type.pt", "id"),
					resource.TestMatchResourceAttr("das_tag_manager_policy.p", "id", regexp.MustCompile(`^pol_`)),
					resource.TestCheckResourceAttr("das_tag_manager_policy.p", "status", "ACTIVE"),
					resource.TestCheckResourceAttr("das_tag_manager_policy.p", "context", "PAGE"),
					resource.TestCheckTypeSetElemNestedAttrs("das_tag_manager_policy.p", "target.*", map[string]string{"target_group": "1", "op": "BETWEEN", "values.0": "0", "values.1": "49"}),
					resource.TestCheckTypeSetElemNestedAttrs("das_tag_manager_policy.p", "target.*", map[string]string{"target_group": "1", "op": "IN", "values.0": "5", "values.1": "10", "values.2": "1"}),
					resource.TestCheckTypeSetElemNestedAttrs("das_tag_manager_policy.p", "inventory.*", map[string]string{"path": "NEWS_APP_IOS/home", "exclude": "true"}),
					resource.TestCheckTypeSetElemNestedAttrs("das_tag_manager_policy.p", "inventory.*", map[string]string{"path": "*", "exclude": "false"}),
					resource.TestCheckResourceAttrPair("das_tag_manager_policy.p", "attachments.0.configuration_id", "das_tag_manager_policy_configuration.cfg", "id"),
					resource.TestCheckResourceAttr("das_tag_manager_policy.p", "attachments.0.enabled", "true"),
					resource.TestCheckResourceAttr("das_tag_manager_policy.p", "attachments.0.target.#", "0"),
					resource.TestCheckTypeSetElemNestedAttrs("das_tag_manager_policy.p", "attachments.1.inventory.*", map[string]string{"path": "NEWS_APP_IOS/*/right", "exclude": "false"}),
					resource.TestCheckTypeSetElemAttrPair("das_tag_manager_policy.p", "labels.*.id", "das_tag_manager_policy_label.l", "id"),
					resource.TestCheckTypeSetElemAttrPair("das_tag_manager_policy.p", "labels.*.id", "das_tag_manager_policy_label.l2", "id"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
			},
			{
				// In place: PATCH the group, PUT the policy and the feature.
				// The ad product has its version now: enabled in place.
				Config: config(server.URL, "traffickers-renamed", 31, false, "ENABLED"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("das_group.g", "name", "traffickers-renamed"),
					resource.TestCheckResourceAttr("das_sign_in_policy.p", "session_days", "31"),
					resource.TestCheckResourceAttr("das_network_feature.f", "enabled", "false"),
					resource.TestCheckResourceAttr("das_delivery_ad_product.a", "name", "Homepage takeover for traffickers-renamed"),
					resource.TestCheckResourceAttr("das_delivery_ad_product.a", "status", "ENABLED"),
					resource.TestCheckResourceAttr("das_delivery_ad_product.a", "current_version", "1"),
					resource.TestCheckResourceAttr("das_delivery_sales_channel.c", "name", "OnetAds for traffickers-renamed"),
					resource.TestCheckResourceAttr("das_delivery_sales_channel.c", "status", "ENABLED"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_type.pt", "label", "Prebid timeout for traffickers-renamed"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_target_dimension.d", "label", "A/B for traffickers-renamed"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_site_platform.s", "position", "31"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_label.l", "name", "floorian-traffickers-renamed"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_catalogue.sizes", "label", "Sizes for traffickers-renamed"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_mapping.map", "label", "Size map for traffickers-renamed"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_mapping.floors", "rows.0.cells.1.value_json", "0.35"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_configuration.cfg", "name", "Standard for traffickers-renamed"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_configuration.cfg", "body_json", `{"timeoutMs":1531}`),
					// Reordered in place: the same policy, cfg2 first now.
					resource.TestCheckResourceAttr("das_tag_manager_policy.p", "name", "Run of network for traffickers-renamed"),
					resource.TestCheckResourceAttrPair("das_tag_manager_policy.p", "attachments.0.configuration_id", "das_tag_manager_policy_configuration.cfg2", "id"),
					resource.TestCheckTypeSetElemNestedAttrs("das_tag_manager_policy.p", "target.*", map[string]string{"target_group": "2", "op": "IN", "values.0": `"4.0.31"`}),
					resource.TestCheckTypeSetElemNestedAttrs("das_tag_manager_policy.p", "labels.*", map[string]string{"name": "floorian-traffickers-renamed"}),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					// A new cell type is a new mapping; a label is changed in place.
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("das_tag_manager_policy_mapping.floors", plancheck.ResourceActionReplace),
						plancheck.ExpectResourceAction("das_tag_manager_policy_mapping.map", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("das_tag_manager_policy_catalogue.sizes", plancheck.ResourceActionUpdate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				// On again, so destroy has something to switch off.
				Config: config(server.URL, "traffickers-renamed", 31, true, "ENABLED"),
				Check:  resource.TestCheckResourceAttr("das_network_feature.f", "enabled", "true"),
			},
			{ResourceName: "das_group.g", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_network_feature.f", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_sign_in_policy.p", ImportState: true, ImportStateId: "sign_in_policy", ImportStateVerify: true},
			{ResourceName: "das_scim_token.t", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_org_api_key.ci", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_delivery_ad_product.a", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_delivery_sales_channel.c", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_tag_manager_policy_type.pt", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_tag_manager_policy_target_dimension.d", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_tag_manager_policy_site_platform.s", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_tag_manager_policy_catalogue.sizes", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_tag_manager_policy_mapping.map", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_tag_manager_policy_mapping.floors", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_tag_manager_policy_label.l", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_tag_manager_policy_configuration.cfg", ImportState: true, ImportStateVerify: true},
			{ResourceName: "das_tag_manager_policy.p", ImportState: true, ImportStateVerify: true},
			// A number in the import id: `<adProductId>/<version>`.
			{ResourceName: "das_delivery_ad_product_version.v", ImportState: true, ImportStateVerify: true},
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
			products, versions := 0, 0
			for p, o := range api.objects {
				switch {
				case strings.Contains(p, "/versions/"):
					versions++
				case strings.HasPrefix(p, "/api/ad-products/"):
					products++
					if o["status"] != "DISABLED" {
						return fmt.Errorf("destroy did not switch the ad product off: %v", o)
					}
				}
			}
			for p, o := range api.objects {
				if strings.HasPrefix(p, "/api/sales-channels/") && (o["status"] != "DISABLED" || len(o["apiKeyIds"].([]any)) != 1) {
					return fmt.Errorf("destroy must switch the sales channel off and keep its keys linked (fail closed): %v", o)
				}
			}
			types, configurations, policies := 0, 0, 0
			for p, o := range api.objects {
				switch {
				case strings.HasPrefix(p, "/api/policies/"):
					policies++
					if o["status"] != "ARCHIVED" || o["archivedAt"] == nil {
						return fmt.Errorf("destroy must archive the policy, never delete it: %v", o)
					}
				case strings.HasPrefix(p, "/api/policy-configurations/"):
					configurations++
					if o["status"] != "ARCHIVED" || o["archivedAt"] == nil {
						return fmt.Errorf("destroy must archive the configuration, never delete it: %v", o)
					}
				case strings.HasPrefix(p, "/api/policy-types/"):
					types++
					if o["status"] != "ARCHIVED" || o["archivedAt"] == nil {
						return fmt.Errorf("destroy must archive the policy type, never delete it: %v", o)
					}
				case strings.HasPrefix(p, "/api/policy-target-dimensions/"), strings.HasPrefix(p, "/api/policy-site-platforms/"),
					strings.HasPrefix(p, "/api/policy-labels/"),
					strings.HasPrefix(p, "/api/policy-catalogues/"), strings.HasPrefix(p, "/api/policy-mappings/"):
					return fmt.Errorf("destroy did not delete %s", p)
				}
			}
			if types != 1 || configurations != 2 || policies != 1 {
				return fmt.Errorf("destroy must keep the archived policy type, configurations and policy: %d, %d, %d", types, configurations, policies)
			}
			if products != 1 || versions != 1 {
				return fmt.Errorf("destroy must keep the ad product and its version (archived, immutable): %d, %d", products, versions)
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
	policiesCreated := 0
	for _, r := range api.requests {
		if r == "POST /api/policies" {
			policiesCreated++
		}
	}
	if policiesCreated != 1 || api.patchedAttachments == 0 || api.keptAttachments < 2 {
		t.Errorf("the reorder must be an update in place keeping the attachments: %d creates, %d attachment patches, %d kept", policiesCreated, api.patchedAttachments, api.keptAttachments)
	}
	mappingsCreated, mappingsPatched := 0, 0
	for _, r := range api.requests {
		if r == "POST /api/policy-mappings" {
			mappingsCreated++
		}
		if strings.HasPrefix(r, "PATCH /api/policy-mappings/") {
			mappingsPatched++
		}
	}
	if mappingsCreated != 3 || mappingsPatched != 1 {
		t.Errorf("a new cell type must replace the mapping (3 creates: map, floors twice) and a label change patch it (1): %d creates, %d patches", mappingsCreated, mappingsPatched)
	}
}

// nullConfig: a policy and a policy type with their nullable inputs set, or
// left out; a second policy type never sets its description.
func nullConfig(url string, set bool) string {
	optional := ""
	typeDescription := ""
	if set {
		optional = "end_at      = \"2027-01-01T00:00:00Z\"\n  description = \"Until the new year\""
		typeDescription = "description = \"Timeouts\""
	}
	return fmt.Sprintf(`
provider "das" {
  token     = "test"
  endpoints = { org = %[1]q, delivery = %[1]q, "tag-manager" = %[1]q }
}

resource "das_tag_manager_policy_type" "pt" {
  key                 = "prebid_timeout_app"
  label               = "Prebid timeout"
  contexts            = ["PAGE"]
  platforms           = ["APP"]
  configuration_types = [{ key = "prebid_timeout_app", fields = [jsonencode({ id = "timeoutMs", kind = "int" })] }]
  %[3]s
}

# Its description is nullable only in the update: the create leaves it out.
resource "das_tag_manager_policy_type" "bare" {
  key                 = "prebid_bare_app"
  label               = "Prebid bare"
  contexts            = ["PAGE"]
  platforms           = ["APP"]
  configuration_types = [{ key = "prebid_bare_app", fields = [jsonencode({ id = "timeoutMs", kind = "int" })] }]
}

resource "das_tag_manager_policy" "p" {
  policy_type_id = das_tag_manager_policy_type.pt.id
  name           = "Run of network"
  inventory      = [{ path = "*" }]
  %[2]s
}
`, url, optional, typeDescription)
}

// A nullable input removed from the configuration plans `-> null`, is
// written as null and leaves an empty plan (P-110); left out from the
// start, it never drifts, also where only the update takes null.
func TestALeftOutInputBecomesNull(t *testing.T) {
	if os.Getenv("TFENV_TERRAFORM_VERSION") == "" {
		t.Setenv("TFENV_TERRAFORM_VERSION", terraformVersion)
	}
	api := newFakeAPI()
	server := httptest.NewServer(api)
	defer server.Close()
	stored := func(prefix string, want map[string]any) resource.TestCheckFunc {
		return func(*terraform.State) error {
			api.mu.Lock()
			defer api.mu.Unlock()
			for p, o := range api.objects {
				if !strings.HasPrefix(p, prefix) || (prefix == "/api/policy-types/" && o["key"] != "prebid_timeout_app") {
					continue
				}
				for k, v := range want {
					if got, ok := o[k]; !ok || got != v {
						return fmt.Errorf("%s: %s is %v (present %v), want %v", p, k, got, ok, v)
					}
				}
				return nil
			}
			return fmt.Errorf("nothing stored under %s", prefix)
		}
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"das": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: nullConfig(server.URL, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("das_tag_manager_policy.p", "end_at", "2027-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("das_tag_manager_policy.p", "description", "Until the new year"),
					resource.TestCheckNoResourceAttr("das_tag_manager_policy.p", "start_at"),
					resource.TestCheckResourceAttr("das_tag_manager_policy_type.pt", "description", "Timeouts"),
					resource.TestCheckNoResourceAttr("das_tag_manager_policy_type.bare", "description"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
			},
			{
				Config: nullConfig(server.URL, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("das_tag_manager_policy.p", "end_at"),
					resource.TestCheckNoResourceAttr("das_tag_manager_policy.p", "description"),
					resource.TestCheckNoResourceAttr("das_tag_manager_policy_type.pt", "description"),
					stored("/api/policies/", map[string]any{"endAt": nil, "description": nil}),
					stored("/api/policy-types/", map[string]any{"description": nil}),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("das_tag_manager_policy.p", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("das_tag_manager_policy.p", tfjsonpath.New("end_at"), knownvalue.Null()),
						plancheck.ExpectKnownValue("das_tag_manager_policy.p", tfjsonpath.New("description"), knownvalue.Null()),
						plancheck.ExpectResourceAction("das_tag_manager_policy_type.pt", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("das_tag_manager_policy_type.pt", tfjsonpath.New("description"), knownvalue.Null()),
						plancheck.ExpectResourceAction("das_tag_manager_policy_type.bare", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}
