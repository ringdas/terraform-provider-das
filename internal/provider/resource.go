package provider

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/ringdas/terraform-provider-das/internal/client"
	"github.com/ringdas/terraform-provider-das/internal/contract"
	"github.com/ringdas/terraform-provider-das/internal/model"
)

// genericResource is every `das_*` resource: its behaviour is the model's
// (internal/model), its calls the contract's operations.
type genericResource struct {
	m *model.Resource
	c *client.Client
}

var (
	_ resource.ResourceWithConfigure   = (*genericResource)(nil)
	_ resource.ResourceWithImportState = (*genericResource)(nil)
	_ resource.ResourceWithModifyPlan  = (*genericResource)(nil)
)

func (r *genericResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.m.TypeName
}

func (r *genericResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	s, err := resourceSchema(r.m)
	if err != nil {
		resp.Diagnostics.AddError("Schema of "+r.m.TypeName, err.Error())
		return
	}
	resp.Schema = s
}

func (r *genericResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if c, ok := req.ProviderData.(*client.Client); ok {
		r.c = c
	}
}

// attrNamed is the attribute with this Terraform name, or nil.
func attrNamed(attrs []*model.Attr, name string) *model.Attr {
	for _, a := range attrs {
		if a.Name == name {
			return a
		}
	}
	return nil
}

func (r *genericResource) identity(vals values) func(string) string {
	return func(param string) string { return vals.stringOf(r.m.IdentityAttr(param)) }
}

// stateFrom builds the state from an answer of the API. With preferBase
// (after a write), a configured input keeps the planned value - Terraform
// requires the applied value to equal the plan; a read (preferBase false)
// shows the API's value, so drift appears in the next plan. A property the
// answer lacks keeps the base value.
func (r *genericResource) stateFrom(answer any, base values, preferBase bool) (tftypes.Value, error) {
	obj, _ := answer.(map[string]any)
	if obj == nil {
		return tftypes.Value{}, fmt.Errorf("the API answered no object")
	}
	vals := map[string]tftypes.Value{}
	for _, a := range r.m.Attrs {
		baseVal, hasBase := base[a.Name]
		known := hasBase && baseVal.IsKnown() && !baseVal.IsNull()
		raw, present := obj[a.Wire]
		switch {
		case a.WriteOnly:
			vals[a.Name] = tftypes.NewValue(a.Type(), nil)
			continue
		case a.VersionOf != "":
			if hasBase && baseVal.IsKnown() {
				vals[a.Name] = baseVal
			} else {
				vals[a.Name] = tftypes.NewValue(a.Type(), nil)
			}
			continue
		case a.Name == "id":
			continue
		case a.Source == model.SourcePath && known,
			preferBase && (a.Required || a.Optional) && known,
			!present && known:
			vals[a.Name] = baseVal
			continue
		}
		v, err := model.ToValue(a, raw)
		if err != nil {
			return tftypes.Value{}, err
		}
		vals[a.Name] = v
	}
	id := r.m.Def.Name
	switch {
	case r.m.IDParam != "":
		s, _ := obj["id"].(string)
		if s == "" {
			s = base.stringOf("id")
		}
		id = s
	case len(r.m.Identity) > 0:
		parts := make([]string, 0, len(r.m.Identity))
		for _, p := range r.m.Identity {
			parts = append(parts, textOf(vals[r.m.IdentityAttr(p)]))
		}
		id = strings.Join(parts, "/")
	}
	vals["id"] = tftypes.NewValue(tftypes.String, id)
	return tftypes.NewValue(model.ObjectType(r.m.Attrs), vals), nil
}

func (r *genericResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan, err := attrsOf(req.Plan.Raw)
	config, err2 := attrsOf(req.Config.Raw)
	if err != nil || err2 != nil {
		resp.Diagnostics.AddError("Reading the plan", fmt.Sprint(err, err2))
		return
	}
	call, err := request(r.m.Spec, r.m.Create, r.m.Attrs, plan, config, pathFrom(r.m.Attrs, plan), func(a *model.Attr) bool { return !a.UpdateOnly })
	if err != nil {
		resp.Diagnostics.AddError("Creating "+r.m.TypeName, err.Error())
		return
	}
	if body, ok := call.Body.(map[string]any); ok {
		for k, v := range r.m.Def.CreateBody {
			body[k] = v
		}
	}
	res, err := r.c.Do(ctx, call)
	if err != nil {
		resp.Diagnostics.AddError("Creating "+r.m.TypeName, err.Error())
		return
	}
	if r.m.Def.CreateAnswer != "" {
		// The entity inside the answer; the rest (a secret) is dropped here.
		wrapped, _ := res.JSON.(map[string]any)
		res.JSON = wrapped[r.m.Def.CreateAnswer]
	}
	state, err := r.stateFrom(res.JSON, plan, true)
	if err != nil {
		resp.Diagnostics.AddError("Creating "+r.m.TypeName, err.Error())
		return
	}
	// Saved before anything else can fail: the entity exists now.
	resp.State.Raw = state
	created, _ := attrsOf(state)

	// What only the update takes (a person's status) is applied next.
	if r.m.Update != nil {
		answer, _ := res.JSON.(map[string]any)
		pending := map[string]bool{}
		for _, a := range r.m.Attrs {
			if !a.UpdateOnly || a.WriteOnly || a.VersionOf != "" {
				continue
			}
			pv := plan[a.Name]
			if !pv.IsKnown() || pv.IsNull() {
				continue
			}
			current, err := model.ToValue(a, answer[a.Wire])
			if err == nil && !current.Equal(pv) {
				pending[a.Name] = true
			}
		}
		if len(pending) > 0 {
			call, err := request(r.m.Spec, r.m.Update, r.m.Attrs, plan, config, r.identity(created), func(a *model.Attr) bool { return pending[a.Name] })
			if err == nil {
				res, err = r.c.Do(ctx, call)
			}
			if err != nil {
				resp.Diagnostics.AddError("Creating "+r.m.TypeName+": applying "+strings.Join(keys(pending), ", "), err.Error())
				return
			}
			if state, err = r.stateFrom(res.JSON, plan, true); err != nil {
				resp.Diagnostics.AddError("Creating "+r.m.TypeName, err.Error())
				return
			}
			resp.State.Raw = state
		}
	}
}

func (r *genericResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	prior, err := attrsOf(req.State.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Reading the state", err.Error())
		return
	}
	call, err := request(r.m.Spec, r.m.Read, nil, prior, nil, r.identity(prior), nil)
	if err != nil {
		resp.Diagnostics.AddError("Reading "+r.m.TypeName, err.Error())
		return
	}
	res, err := r.c.Do(ctx, call)
	if client.NotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Reading "+r.m.TypeName, err.Error())
		return
	}
	state, err := r.stateFrom(res.JSON, prior, false)
	if err != nil {
		resp.Diagnostics.AddError("Reading "+r.m.TypeName, err.Error())
		return
	}
	resp.State.Raw = state
}

func (r *genericResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan, err := attrsOf(req.Plan.Raw)
	prior, err2 := attrsOf(req.State.Raw)
	config, err3 := attrsOf(req.Config.Raw)
	if err != nil || err2 != nil || err3 != nil {
		resp.Diagnostics.AddError("Reading the plan", fmt.Sprint(err, err2, err3))
		return
	}
	if r.m.Update == nil {
		// Only attributes Terraform keeps alone changed (a secret's version
		// on a resource without an update forces a replacement instead).
		resp.State.Raw = req.Plan.Raw
		return
	}
	rotate := map[string]bool{}
	for _, a := range r.m.Attrs {
		if a.VersionOf != "" && !plan[a.Name].Equal(prior[a.Name]) {
			rotate[a.VersionOf] = true
		}
	}
	patch := r.m.Update.Method == http.MethodPatch
	sent := 0
	include := func(a *model.Attr) bool {
		ok := false
		switch {
		case !a.InUpdate:
		case a.WriteOnly:
			ok = rotate[a.Name]
		case patch:
			ok = plan[a.Name].IsKnown() && !plan[a.Name].Equal(prior[a.Name])
		default:
			ok = true
		}
		if ok {
			sent++
		}
		return ok
	}
	call, err := request(r.m.Spec, r.m.Update, r.m.Attrs, plan, config, r.identity(prior), include)
	if err != nil {
		resp.Diagnostics.AddError("Updating "+r.m.TypeName, err.Error())
		return
	}
	var answer any
	if sent > 0 || !patch {
		res, err := r.c.Do(ctx, call)
		if err != nil {
			resp.Diagnostics.AddError("Updating "+r.m.TypeName, err.Error())
			return
		}
		answer = res.JSON
	}
	if answer == nil {
		read, err := request(r.m.Spec, r.m.Read, nil, prior, nil, r.identity(prior), nil)
		if err == nil {
			var res *client.Response
			if res, err = r.c.Do(ctx, read); err == nil {
				answer = res.JSON
			}
		}
		if err != nil {
			resp.Diagnostics.AddError("Updating "+r.m.TypeName, err.Error())
			return
		}
	}
	state, err := r.stateFrom(answer, plan, true)
	if err != nil {
		resp.Diagnostics.AddError("Updating "+r.m.TypeName, err.Error())
		return
	}
	resp.State.Raw = state
}

func (r *genericResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	prior, err := attrsOf(req.State.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Reading the state", err.Error())
		return
	}
	var call client.Request
	switch {
	case r.m.Delete != nil:
		call, err = request(r.m.Spec, r.m.Delete, nil, prior, nil, r.identity(prior), nil)
	case r.m.Def.Destroy != nil && r.m.Def.Destroy.Via == "update":
		call, err = request(r.m.Spec, r.m.Update, nil, prior, nil, r.identity(prior), nil)
		call.Body = r.m.Def.Destroy.Body
	default:
		// Forgotten, not deleted: ModifyPlan warned at plan.
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Destroying "+r.m.TypeName, err.Error())
		return
	}
	if _, err := r.c.Do(ctx, call); err != nil && !client.NotFound(err) {
		resp.Diagnostics.AddError("Destroying "+r.m.TypeName, err.Error())
	}
}

func (r *genericResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if len(r.m.Identity) == 0 {
		// A singleton: any id imports the one there is.
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(r.m.Def.Name))...)
		return
	}
	parts := strings.Split(req.ID, "/")
	if len(parts) != len(r.m.Identity) {
		resp.Diagnostics.AddError("Import id of "+r.m.TypeName, fmt.Sprintf("expected %s, got %q", strings.Join(r.m.Identity, "/"), req.ID))
		return
	}
	for i, p := range r.m.Identity {
		if parts[i] == "" {
			resp.Diagnostics.AddError("Import id of "+r.m.TypeName, fmt.Sprintf("empty %s in %q", p, req.ID))
			return
		}
		name := r.m.IdentityAttr(p)
		var value any = types.StringValue(parts[i])
		if a := attrNamed(r.m.Attrs, name); a != nil && a.Kind == contract.KindInteger {
			// A number in the path (a version): the attribute holds it as a number.
			n, err := strconv.ParseInt(parts[i], 10, 64)
			if err != nil {
				resp.Diagnostics.AddError("Import id of "+r.m.TypeName, fmt.Sprintf("%s is not an integer in %q", p, req.ID))
				return
			}
			value = types.Int64Value(n)
		}
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(name), value)...)
	}
	if r.m.IDParam == "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
	}
}

// ModifyPlan: an input the update does not take replaces the resource (or,
// on one the API never deletes, is refused); an input left out of the
// configuration keeps its value instead of turning unknown; a destroy that
// only forgets says so.
func (r *genericResource) ModifyPlan(_ context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		if r.m.Delete == nil && r.m.Def.Destroy != nil && r.m.Def.Destroy.Via == "forget" {
			resp.Diagnostics.AddWarning(r.m.TypeName+" is only forgotten", r.m.Def.Destroy.Reason)
		}
		return
	}
	if req.State.Raw.IsNull() {
		return
	}
	plan, err := attrsOf(req.Plan.Raw)
	prior, err2 := attrsOf(req.State.Raw)
	config, err3 := attrsOf(req.Config.Raw)
	if err != nil || err2 != nil || err3 != nil {
		resp.Diagnostics.AddError("Reading the plan", fmt.Sprint(err, err2, err3))
		return
	}
	changedPlan := false
	for _, a := range r.m.Attrs {
		pv, sv := plan[a.Name], prior[a.Name]
		cv, configured := config[a.Name]
		unconfigured := !configured || cv.IsNull()
		// Data conventions (docs/engineering/data-conventions.md): who created
		// an entity and when never changes.
		if a.Computed && !a.Optional && !pv.IsKnown() && sv.IsKnown() && !sv.IsNull() && immutableComputed[a.Wire] {
			plan[a.Name] = sv
			changedPlan = true
			continue
		}
		if a.Computed && (a.Optional || a.Required) && unconfigured && !pv.IsKnown() {
			plan[a.Name] = sv
			changedPlan = true
			continue
		}
		if !a.ForceNew && !a.Immutable {
			continue
		}
		if pv.IsKnown() && pv.Equal(sv) {
			continue
		}
		if a.Immutable {
			if pv.IsKnown() {
				resp.Diagnostics.AddAttributeError(path.Root(a.Name), "Cannot change "+a.Name,
					fmt.Sprintf("%s.%s cannot change: the API offers no update of it and never deletes a %s, so it cannot be replaced either.", r.m.TypeName, a.Name, strings.TrimPrefix(r.m.TypeName, "das_")))
			}
			continue
		}
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root(a.Name))
	}
	if changedPlan {
		resp.Plan.Raw = tftypes.NewValue(model.ObjectType(r.m.Attrs), map[string]tftypes.Value(plan))
	}
}

// immutableComputed are answer fields no update changes.
var immutableComputed = map[string]bool{"createdAt": true, "createdBy": true}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
