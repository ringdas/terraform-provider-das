// Package provider is the DAS Terraform provider (`ringdas/das`): one
// generic resource, data source, action and ephemeral resource,
// instantiated once per entry of the API contracts' Terraform surfaces
// (internal/model). Nothing here is written for a particular DAS entity,
// so a new feature of an API reaches Terraform with the API's contract.
package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ringdas/terraform-provider-das/internal/client"
	"github.com/ringdas/terraform-provider-das/internal/model"
	"github.com/ringdas/terraform-provider-das/internal/specs"
)

// Address is the provider's registry address.
const Address = "registry.terraform.io/ringdas/das"

type dasProvider struct {
	version string
	model   *model.Provider
}

var (
	_ provider.ProviderWithActions            = (*dasProvider)(nil)
	_ provider.ProviderWithEphemeralResources = (*dasProvider)(nil)
)

// New returns the provider built from the embedded contracts.
func New(version string) func() provider.Provider {
	m, err := specs.Model()
	if err != nil {
		// The contracts are embedded at build time and the provider's tests
		// build this model: a failure here never reaches a release.
		panic(err)
	}
	return func() provider.Provider { return &dasProvider{version: version, model: m} }
}

// NewWithModel serves a model built elsewhere (tests).
func NewWithModel(version string, m *model.Provider) func() provider.Provider {
	return func() provider.Provider { return &dasProvider{version: version, model: m} }
}

type providerConfig struct {
	OrgDomain types.String `tfsdk:"org_domain"`
	Scheme    types.String `tfsdk:"scheme"`
	Network   types.String `tfsdk:"network"`
	APIKey    types.String `tfsdk:"api_key"`
	Token     types.String `tfsdk:"token"`
	Endpoints types.Map    `tfsdk:"endpoints"`
}

// ProviderSchemaDescription documents the provider block.
const ProviderSchemaDescription = "The DAS provider configures a DAS Organization as code: everything the `das` CLI can do. " +
	"Its resources, data sources, actions and ephemeral resources are derived from the API contracts of the DAS version it was released with."

func (p *dasProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "das"
	resp.Version = p.version
}

// ProviderAttributes are the provider block's attributes, also documented by internal/docgen.
var ProviderAttributes = []struct {
	Name, Env, Description string
	Sensitive              bool
}{
	{"org_domain", "DAS_ORG_DOMAIN", "The Organization's domain, e.g. `das.example.com`. The Organisation API answers on it, a product API on `<network>.<product>.<org_domain>`.", false},
	{"scheme", "DAS_ORG_SCHEME", "`https` (default), or `http` for a local stack.", false},
	{"network", "DAS_NETWORK", "The slug of the Network a product API (Delivery, Audience, Tag Manager) is called for. Omitted, it is the API key's own Network.", false},
	{"api_key", "DAS_API_KEY", "An API key: a Network's, `das_<networkId>_<secret>`, or the Organization's, `daso_<secret>` (`das_org_api_key`), which administers the Organization layer only - a product API then takes `token`. The tenant of a call is always the credential's.", true},
	{"token", "DAS_TOKEN", "A bearer token instead of an API key: an Organization session token, or a Network token.", true},
}

func (p *dasProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	attrs := map[string]schema.Attribute{}
	for _, a := range ProviderAttributes {
		d := a.Description + " Environment variable `" + a.Env + "`."
		attrs[a.Name] = schema.StringAttribute{Optional: true, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}
	}
	attrs["endpoints"] = schema.MapAttribute{
		ElementType: types.StringType, Optional: true,
		Description:         "Base URLs by API module (`org`, `delivery`, `audience`, `tag-manager`) instead of the hosts derived from `org_domain`: local development and tests.",
		MarkdownDescription: "Base URLs by API module (`org`, `delivery`, `audience`, `tag-manager`) instead of the hosts derived from `org_domain`: local development and tests.",
	}
	resp.Schema = schema.Schema{Description: ProviderSchemaDescription, MarkdownDescription: ProviderSchemaDescription, Attributes: attrs}
}

func pick(v types.String, env string) string {
	if !v.IsNull() && !v.IsUnknown() && v.ValueString() != "" {
		return v.ValueString()
	}
	return os.Getenv(env)
}

func (p *dasProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg providerConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	endpoints := map[string]string{}
	if !cfg.Endpoints.IsNull() && !cfg.Endpoints.IsUnknown() {
		resp.Diagnostics.Append(cfg.Endpoints.ElementsAs(ctx, &endpoints, false)...)
	}
	c, err := client.New(client.Config{
		OrgDomain: pick(cfg.OrgDomain, "DAS_ORG_DOMAIN"),
		Scheme:    pick(cfg.Scheme, "DAS_ORG_SCHEME"),
		Network:   pick(cfg.Network, "DAS_NETWORK"),
		APIKey:    pick(cfg.APIKey, "DAS_API_KEY"),
		Token:     pick(cfg.Token, "DAS_TOKEN"),
		Endpoints: endpoints,
		UserAgent: "terraform-provider-das/" + p.version,
	})
	if err != nil {
		resp.Diagnostics.AddError("Configuring the DAS provider", err.Error())
		return
	}
	resp.DataSourceData = c
	resp.ResourceData = c
	resp.ActionData = c
	resp.EphemeralResourceData = c
}

func (p *dasProvider) Resources(context.Context) []func() resource.Resource {
	out := make([]func() resource.Resource, 0, len(p.model.Resources))
	for _, m := range p.model.Resources {
		out = append(out, func() resource.Resource { return &genericResource{m: m} })
	}
	return out
}

func (p *dasProvider) DataSources(context.Context) []func() datasource.DataSource {
	out := make([]func() datasource.DataSource, 0, len(p.model.DataSources))
	for _, m := range p.model.DataSources {
		out = append(out, func() datasource.DataSource { return &genericData{m: m} })
	}
	return out
}

func (p *dasProvider) Actions(context.Context) []func() action.Action {
	out := make([]func() action.Action, 0, len(p.model.Actions))
	for _, m := range p.model.Actions {
		out = append(out, func() action.Action { return &genericAction{m: m} })
	}
	return out
}

func (p *dasProvider) EphemeralResources(context.Context) []func() ephemeral.EphemeralResource {
	out := make([]func() ephemeral.EphemeralResource, 0, len(p.model.Ephemerals))
	for _, m := range p.model.Ephemerals {
		out = append(out, func() ephemeral.EphemeralResource { return &genericEphemeral{m: m} })
	}
	return out
}
