package provider_test

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/ringdas/terraform-provider-das/internal/model"
	"github.com/ringdas/terraform-provider-das/internal/provider"
	"github.com/ringdas/terraform-provider-das/internal/specs"
)

// THE PROVIDER GATE (docs/architecture/terraform-provider.md): the provider
// built from the embedded contracts serves a valid schema for every entry
// of every API's Terraform surface. A contract construct the generic
// engine cannot map, a reserved or clashing attribute name, or a surface
// entry that yields no Terraform type fails here - in CI, never in a
// user's plan. No resource is written by hand, so nothing else has to be
// remembered when an API gains a feature.

func TestEverySurfaceEntryIsServed(t *testing.T) {
	m, err := specs.Model()
	if err != nil {
		t.Fatalf("the provider cannot be derived from the embedded contracts: %v", err)
	}
	want := map[string][]string{}
	for _, s := range m.Specs {
		prefix := model.Prefix(s.Module)
		for _, r := range s.Surface.Resources {
			want["resource"] = append(want["resource"], prefix+r.Name)
			want["data source"] = append(want["data source"], prefix+r.Name)
			if r.List != "" {
				name := r.ListName
				if name == "" {
					name = r.Name + "s"
				}
				want["data source"] = append(want["data source"], prefix+name)
			}
		}
		for _, d := range s.Surface.DataSources {
			want["data source"] = append(want["data source"], prefix+d.Name)
		}
		for _, a := range s.Surface.Actions {
			want["action"] = append(want["action"], prefix+a.Name)
		}
		for _, e := range s.Surface.EphemeralResources {
			want["ephemeral resource"] = append(want["ephemeral resource"], prefix+e.Name)
		}
	}

	server, err := providerserver.NewProtocol6WithError(provider.New("test")())()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := server.GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range resp.Diagnostics {
		t.Errorf("%s: %s: %s", d.Severity, d.Summary, d.Detail)
	}
	got := map[string][]string{}
	for name := range resp.ResourceSchemas {
		got["resource"] = append(got["resource"], name)
	}
	for name := range resp.DataSourceSchemas {
		got["data source"] = append(got["data source"], name)
	}
	for name := range resp.ActionSchemas {
		got["action"] = append(got["action"], name)
	}
	for name := range resp.EphemeralResourceSchemas {
		got["ephemeral resource"] = append(got["ephemeral resource"], name)
	}
	for _, kind := range []string{"resource", "data source", "action", "ephemeral resource"} {
		sort.Strings(want[kind])
		sort.Strings(got[kind])
		if strings.Join(want[kind], ",") != strings.Join(got[kind], ",") {
			t.Errorf("%ss served %v, the surfaces declare %v", kind, got[kind], want[kind])
		}
	}
	if len(got["resource"]) == 0 {
		t.Fatal("no resources served")
	}
}

func TestNoAttributeLeaksASecretIntoState(t *testing.T) {
	m, err := specs.Model()
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range m.Resources {
		for _, a := range r.Attrs {
			if a.WriteOnly && (a.Computed || !a.Sensitive) {
				t.Errorf("%s.%s: a write-only secret is sensitive and never computed", r.TypeName, a.Name)
			}
		}
		for _, op := range []string{r.Def.Create, r.Def.Update} {
			if op == "" {
				continue
			}
			o, _ := r.Spec.Operation(op)
			if o.Body == nil {
				continue
			}
			for _, p := range o.Body.Properties {
				if !p.Schema.WriteOnly {
					continue
				}
				a := model.Find(r.Attrs, model.SnakeCase(p.Name)+"_wo")
				if a == nil || !a.WriteOnly {
					t.Errorf("%s: the secret %s is not a write-only attribute", r.TypeName, p.Name)
				}
			}
		}
	}
}

// An optional input of a resource the API takes as null (in its update,
// or in its create when the update does not take it) is null when left
// out (P-110): optional, sent as null, never computed - so a value removed
// from the configuration plans `-> null` instead of staying. An optional input the API never takes as null stays
// optional and computed: left out, it keeps the API's value.
func TestALeftOutNullableInputIsNull(t *testing.T) {
	m, err := specs.Model()
	if err != nil {
		t.Fatal(err)
	}
	// The update decides when it takes the input: a null it refuses could
	// never be sent after a removal. Otherwise the create decides.
	takesNull := func(r *model.Resource, wire string) bool {
		for _, op := range []string{r.Def.Update, r.Def.Create} {
			if op == "" {
				continue
			}
			o, _ := r.Spec.Operation(op)
			if p := o.Body.Property(wire); p != nil {
				return p.Schema.Nullable
			}
		}
		return false
	}
	for _, r := range m.Resources {
		for _, a := range r.Attrs {
			if a.Source != model.SourceBody || a.Required || a.UpdateOnly || a.WriteOnly || a.VersionOf != "" {
				continue
			}
			nullable := takesNull(r, a.Wire)
			if nullable && (!a.Optional || !a.SendNull || a.Computed) {
				t.Errorf("%s.%s: the API takes it as null, so it is optional, sent as null and never computed", r.TypeName, a.Name)
			}
			if !nullable && a.SendNull {
				t.Errorf("%s.%s: the API never takes it as null, yet it is sent as null", r.TypeName, a.Name)
			}
		}
	}
	for _, c := range []struct {
		resource, attr     string
		sendNull, computed bool
	}{
		{"das_tag_manager_policy", "end_at", true, false},
		{"das_tag_manager_policy", "description", true, false},
		// Nullable only in the update: the create leaves it out.
		{"das_tag_manager_policy_type", "description", true, false},
		{"das_identity_provider", "oidc_issuer", true, false},
		// Optional, never null: the API's default when left out.
		{"das_tag_manager_policy", "status", false, true},
		// Its create schema has no type: not nullable for the provider (P-110).
		{"das_tag_manager_policy_mapping", "cell", false, true},
	} {
		var r *model.Resource
		for _, x := range m.Resources {
			if x.TypeName == c.resource {
				r = x
			}
		}
		a := model.Find(r.Attrs, c.attr)
		if a == nil || !a.Optional || a.SendNull != c.sendNull || a.Computed != c.computed {
			t.Errorf("%s.%s: want optional, SendNull %v, Computed %v; got %+v", c.resource, c.attr, c.sendNull, c.computed, a)
		}
	}
}
