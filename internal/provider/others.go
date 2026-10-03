package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/ringdas/terraform-provider-das/internal/client"
	"github.com/ringdas/terraform-provider-das/internal/model"
)

// resultFrom is the object of a data source or ephemeral resource: the
// inputs as configured, the rest from the answer.
func resultFrom(attrs []*model.Attr, answer any, config values) (tftypes.Value, error) {
	obj, _ := answer.(map[string]any)
	vals := map[string]tftypes.Value{}
	for _, a := range attrs {
		if (a.Required || a.Optional) && config[a.Name].IsKnown() && !config[a.Name].IsNull() {
			vals[a.Name] = config[a.Name]
			continue
		}
		var raw any
		if obj != nil {
			raw = obj[a.Wire]
		}
		v, err := model.ToValue(a, raw)
		if err != nil {
			return tftypes.Value{}, err
		}
		vals[a.Name] = v
	}
	return tftypes.NewValue(model.ObjectType(attrs), vals), nil
}

// -------------------------------------------------------------- data source

type genericData struct {
	m *model.DataSource
	c *client.Client
}

var _ datasource.DataSourceWithConfigure = (*genericData)(nil)

func (d *genericData) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = d.m.TypeName
}

func (d *genericData) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	s, err := dataSchema(d.m)
	if err != nil {
		resp.Diagnostics.AddError("Schema of "+d.m.TypeName, err.Error())
		return
	}
	resp.Schema = s
}

func (d *genericData) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if c, ok := req.ProviderData.(*client.Client); ok {
		d.c = c
	}
}

func (d *genericData) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	config, err := attrsOf(req.Config.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Reading the configuration", err.Error())
		return
	}
	pathValue := func(param string) string {
		if param == d.m.IDParam {
			return config.stringOf("id")
		}
		return pathFrom(d.m.Attrs, config)(param)
	}
	call, err := request(d.m.Spec, d.m.Op, d.m.Attrs, config, config, pathValue, nil)
	if err != nil {
		resp.Diagnostics.AddError("Reading "+d.m.TypeName, err.Error())
		return
	}
	var answer any
	if d.m.Page {
		items, err := d.c.All(ctx, call)
		if err != nil {
			resp.Diagnostics.AddError("Reading "+d.m.TypeName, err.Error())
			return
		}
		if items == nil {
			items = []any{}
		}
		answer = map[string]any{"items": items}
	} else {
		res, err := d.c.Do(ctx, call)
		if err != nil {
			resp.Diagnostics.AddError("Reading "+d.m.TypeName, err.Error())
			return
		}
		answer = res.JSON
	}
	state, err := resultFrom(d.m.Attrs, answer, config)
	if err != nil {
		resp.Diagnostics.AddError("Reading "+d.m.TypeName, err.Error())
		return
	}
	resp.State.Raw = state
}

// ------------------------------------------------------------------ action

type genericAction struct {
	m *model.Action
	c *client.Client
}

var _ action.ActionWithConfigure = (*genericAction)(nil)

func (a *genericAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = a.m.TypeName
}

func (a *genericAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	s, err := actionSchema(a.m)
	if err != nil {
		resp.Diagnostics.AddError("Schema of "+a.m.TypeName, err.Error())
		return
	}
	resp.Schema = s
}

func (a *genericAction) Configure(_ context.Context, req action.ConfigureRequest, _ *action.ConfigureResponse) {
	if c, ok := req.ProviderData.(*client.Client); ok {
		a.c = c
	}
}

func (a *genericAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	config, err := attrsOf(req.Config.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Reading the configuration", err.Error())
		return
	}
	call, err := request(a.m.Spec, a.m.Op, a.m.Attrs, config, config, pathFrom(a.m.Attrs, config), nil)
	if err != nil {
		resp.Diagnostics.AddError("Invoking "+a.m.TypeName, err.Error())
		return
	}
	res, err := a.c.Do(ctx, call)
	if err != nil {
		resp.Diagnostics.AddError("Invoking "+a.m.TypeName, err.Error())
		return
	}
	if resp.SendProgress != nil {
		resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("%s %s", a.m.Op.Ref, summary(res))})
	}
}

// --------------------------------------------------------------- ephemeral

type genericEphemeral struct {
	m *model.Ephemeral
	c *client.Client
}

var _ ephemeral.EphemeralResourceWithConfigure = (*genericEphemeral)(nil)

func (e *genericEphemeral) Metadata(_ context.Context, _ ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = e.m.TypeName
}

func (e *genericEphemeral) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	s, err := ephemeralSchema(e.m)
	if err != nil {
		resp.Diagnostics.AddError("Schema of "+e.m.TypeName, err.Error())
		return
	}
	resp.Schema = s
}

func (e *genericEphemeral) Configure(_ context.Context, req ephemeral.ConfigureRequest, _ *ephemeral.ConfigureResponse) {
	if c, ok := req.ProviderData.(*client.Client); ok {
		e.c = c
	}
}

func (e *genericEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	config, err := attrsOf(req.Config.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Reading the configuration", err.Error())
		return
	}
	call, err := request(e.m.Spec, e.m.Op, e.m.Attrs, config, config, pathFrom(e.m.Attrs, config), nil)
	if err != nil {
		resp.Diagnostics.AddError("Opening "+e.m.TypeName, err.Error())
		return
	}
	res, err := e.c.Do(ctx, call)
	if err != nil {
		resp.Diagnostics.AddError("Opening "+e.m.TypeName, err.Error())
		return
	}
	result, err := resultFrom(e.m.Attrs, res.JSON, config)
	if err != nil {
		resp.Diagnostics.AddError("Opening "+e.m.TypeName, err.Error())
		return
	}
	resp.Result.Raw = result
}
