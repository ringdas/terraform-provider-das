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
