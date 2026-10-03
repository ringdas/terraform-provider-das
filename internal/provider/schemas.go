package provider

import (
	"context"
	"fmt"
	"strings"

	aschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	eschema "github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ringdas/terraform-provider-das/internal/contract"
	"github.com/ringdas/terraform-provider-das/internal/model"
)

// Description is an attribute's documentation: the contract's words, the
// allowed values and what a change does.
func Description(a *model.Attr) string {
	d := strings.TrimSpace(a.Description)
	if d != "" && !strings.HasSuffix(d, ".") {
		d += "."
	}
	if len(a.Enum) > 0 {
		d += " One of `" + strings.Join(a.Enum, "`, `") + "`."
	}
	switch {
	case a.ForceNew:
		d += " Changing it replaces the resource."
	case a.Immutable:
		d += " It cannot change: the API offers no update and never deletes it."
	case a.UpdateOnly:
		d += " Applied after create, by an update."
	}
	return strings.TrimSpace(d)
}

// elemType is the framework type of a collection's element.
func elemType(a *model.Attr) (attr.Type, error) {
	switch a.Kind {
	case contract.KindString:
		return types.StringType, nil
	case contract.KindInteger:
		return types.Int64Type, nil
	case contract.KindNumber:
		return types.NumberType, nil
	case contract.KindBool:
		return types.BoolType, nil
	case contract.KindArray:
		e, err := elemType(a.Elem)
		return types.ListType{ElemType: e}, err
	case contract.KindMap:
		e, err := elemType(a.Elem)
		return types.MapType{ElemType: e}, err
	case contract.KindObject:
		attrs := map[string]attr.Type{}
		for _, na := range a.Attrs {
			t, err := elemType(na)
			if err != nil {
				return nil, err
			}
			attrs[na.Name] = t
		}
		return types.ObjectType{AttrTypes: attrs}, nil
	default:
		return nil, fmt.Errorf("a value of no fixed shape inside a list or map is not supported by Terraform")
	}
}

// oneOf validates an enum.
type oneOf []string

func (o oneOf) Description(context.Context) string { return "one of " + strings.Join(o, ", ") }
func (o oneOf) MarkdownDescription(ctx context.Context) string {
	return o.Description(ctx)
}
func (o oneOf) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	v := req.ConfigValue.ValueString()
	for _, allowed := range o {
		if v == allowed {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(req.Path, "Invalid value", fmt.Sprintf("%q is not one of %s", v, strings.Join(o, ", ")))
}

func stringValidators(a *model.Attr) []validator.String {
	if len(a.Enum) == 0 || a.Computed && !a.Optional && !a.Required {
		return nil
	}
	return []validator.String{oneOf(a.Enum)}
}

func isNestedObject(a *model.Attr) bool {
	return (a.Kind == contract.KindArray || a.Kind == contract.KindMap) && a.Elem.Kind == contract.KindObject
}

// ---------------------------------------------------------------- resource

func resourceSchema(r *model.Resource) (rschema.Schema, error) {
	attrs, err := resourceAttrs(r.Attrs)
	if err != nil {
		return rschema.Schema{}, err
	}
	if id, ok := attrs["id"].(rschema.StringAttribute); ok {
		id.PlanModifiers = []planmodifier.String{stringplanmodifier.UseStateForUnknown()}
		attrs["id"] = id
	}
	return rschema.Schema{Description: r.Description, MarkdownDescription: r.Description, Attributes: attrs}, nil
}

func resourceAttrs(list []*model.Attr) (map[string]rschema.Attribute, error) {
	out := map[string]rschema.Attribute{}
	for _, a := range list {
		at, err := resourceAttr(a)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", a.Name, err)
		}
		out[a.Name] = at
	}
	return out, nil
}

func resourceAttr(a *model.Attr) (rschema.Attribute, error) {
	d := Description(a)
	switch a.Kind {
	case contract.KindString:
		return rschema.StringAttribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, WriteOnly: a.WriteOnly, Description: d, MarkdownDescription: d, Validators: stringValidators(a)}, nil
	case contract.KindInteger:
		return rschema.Int64Attribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, WriteOnly: a.WriteOnly, Description: d, MarkdownDescription: d}, nil
	case contract.KindNumber:
		return rschema.NumberAttribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, WriteOnly: a.WriteOnly, Description: d, MarkdownDescription: d}, nil
	case contract.KindBool:
		return rschema.BoolAttribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, WriteOnly: a.WriteOnly, Description: d, MarkdownDescription: d}, nil
	case contract.KindAny:
		return rschema.DynamicAttribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, WriteOnly: a.WriteOnly, Description: d, MarkdownDescription: d}, nil
	case contract.KindObject:
		nested, err := resourceAttrs(a.Attrs)
		if err != nil {
			return nil, err
		}
		return rschema.SingleNestedAttribute{Attributes: nested, Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, WriteOnly: a.WriteOnly, Description: d, MarkdownDescription: d}, nil
	case contract.KindArray, contract.KindMap:
		if isNestedObject(a) {
			nested, err := resourceAttrs(a.Elem.Attrs)
			if err != nil {
				return nil, err
			}
			obj := rschema.NestedAttributeObject{Attributes: nested}
			if a.Kind == contract.KindArray {
				return rschema.ListNestedAttribute{NestedObject: obj, Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, WriteOnly: a.WriteOnly, Description: d, MarkdownDescription: d}, nil
			}
			return rschema.MapNestedAttribute{NestedObject: obj, Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, WriteOnly: a.WriteOnly, Description: d, MarkdownDescription: d}, nil
		}
		et, err := elemType(a.Elem)
		if err != nil {
			return nil, err
		}
		if a.Kind == contract.KindArray {
			return rschema.ListAttribute{ElementType: et, Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, WriteOnly: a.WriteOnly, Description: d, MarkdownDescription: d}, nil
		}
		return rschema.MapAttribute{ElementType: et, Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, WriteOnly: a.WriteOnly, Description: d, MarkdownDescription: d}, nil
	}
	return nil, fmt.Errorf("unsupported kind %s", a.Kind)
}

// ------------------------------------------------------------- data source

func dataSchema(ds *model.DataSource) (dschema.Schema, error) {
	attrs, err := dataAttrs(ds.Attrs)
	return dschema.Schema{Description: ds.Description, MarkdownDescription: ds.Description, Attributes: attrs}, err
}

func dataAttrs(list []*model.Attr) (map[string]dschema.Attribute, error) {
	out := map[string]dschema.Attribute{}
	for _, a := range list {
		at, err := dataAttr(a)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", a.Name, err)
		}
		out[a.Name] = at
	}
	return out, nil
}

func dataAttr(a *model.Attr) (dschema.Attribute, error) {
	d := Description(a)
	switch a.Kind {
	case contract.KindString:
		return dschema.StringAttribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d, Validators: stringValidators(a)}, nil
	case contract.KindInteger:
		return dschema.Int64Attribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
	case contract.KindNumber:
		return dschema.NumberAttribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
	case contract.KindBool:
		return dschema.BoolAttribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
	case contract.KindAny:
		return dschema.DynamicAttribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
	case contract.KindObject:
		nested, err := dataAttrs(a.Attrs)
		if err != nil {
			return nil, err
		}
		return dschema.SingleNestedAttribute{Attributes: nested, Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
	case contract.KindArray, contract.KindMap:
		if isNestedObject(a) {
			nested, err := dataAttrs(a.Elem.Attrs)
			if err != nil {
				return nil, err
			}
			obj := dschema.NestedAttributeObject{Attributes: nested}
			if a.Kind == contract.KindArray {
				return dschema.ListNestedAttribute{NestedObject: obj, Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
			}
			return dschema.MapNestedAttribute{NestedObject: obj, Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
		}
		et, err := elemType(a.Elem)
		if err != nil {
			return nil, err
		}
		if a.Kind == contract.KindArray {
			return dschema.ListAttribute{ElementType: et, Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
		}
		return dschema.MapAttribute{ElementType: et, Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
	}
	return nil, fmt.Errorf("unsupported kind %s", a.Kind)
}

// -------------------------------------------------------------- ephemeral

func ephemeralSchema(e *model.Ephemeral) (eschema.Schema, error) {
	attrs, err := ephemeralAttrs(e.Attrs)
	return eschema.Schema{Description: e.Description, MarkdownDescription: e.Description, Attributes: attrs}, err
}

func ephemeralAttrs(list []*model.Attr) (map[string]eschema.Attribute, error) {
	out := map[string]eschema.Attribute{}
	for _, a := range list {
		at, err := ephemeralAttr(a)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", a.Name, err)
		}
		out[a.Name] = at
	}
	return out, nil
}

func ephemeralAttr(a *model.Attr) (eschema.Attribute, error) {
	d := Description(a)
	switch a.Kind {
	case contract.KindString:
		return eschema.StringAttribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d, Validators: stringValidators(a)}, nil
	case contract.KindInteger:
		return eschema.Int64Attribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
	case contract.KindNumber:
		return eschema.NumberAttribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
	case contract.KindBool:
		return eschema.BoolAttribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
	case contract.KindAny:
		return eschema.DynamicAttribute{Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
	case contract.KindObject:
		nested, err := ephemeralAttrs(a.Attrs)
		if err != nil {
			return nil, err
		}
		return eschema.SingleNestedAttribute{Attributes: nested, Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
	case contract.KindArray, contract.KindMap:
		if isNestedObject(a) {
			nested, err := ephemeralAttrs(a.Elem.Attrs)
			if err != nil {
				return nil, err
			}
			obj := eschema.NestedAttributeObject{Attributes: nested}
			if a.Kind == contract.KindArray {
				return eschema.ListNestedAttribute{NestedObject: obj, Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
			}
			return eschema.MapNestedAttribute{NestedObject: obj, Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
		}
		et, err := elemType(a.Elem)
		if err != nil {
			return nil, err
		}
		if a.Kind == contract.KindArray {
			return eschema.ListAttribute{ElementType: et, Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
		}
		return eschema.MapAttribute{ElementType: et, Required: a.Required, Optional: a.Optional, Computed: a.Computed, Sensitive: a.Sensitive, Description: d, MarkdownDescription: d}, nil
	}
	return nil, fmt.Errorf("unsupported kind %s", a.Kind)
}

// ----------------------------------------------------------------- action

func actionSchema(a *model.Action) (aschema.Schema, error) {
	attrs, err := actionAttrs(a.Attrs)
	return aschema.Schema{Description: a.Description, MarkdownDescription: a.Description, Attributes: attrs}, err
}

func actionAttrs(list []*model.Attr) (map[string]aschema.Attribute, error) {
	out := map[string]aschema.Attribute{}
	for _, a := range list {
		at, err := actionAttr(a)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", a.Name, err)
		}
		out[a.Name] = at
	}
	return out, nil
}

func actionAttr(a *model.Attr) (aschema.Attribute, error) {
	d := Description(a)
	switch a.Kind {
	case contract.KindString:
		return aschema.StringAttribute{Required: a.Required, Optional: a.Optional, Description: d, MarkdownDescription: d, Validators: stringValidators(a)}, nil
	case contract.KindInteger:
		return aschema.Int64Attribute{Required: a.Required, Optional: a.Optional, Description: d, MarkdownDescription: d}, nil
	case contract.KindNumber:
		return aschema.NumberAttribute{Required: a.Required, Optional: a.Optional, Description: d, MarkdownDescription: d}, nil
	case contract.KindBool:
		return aschema.BoolAttribute{Required: a.Required, Optional: a.Optional, Description: d, MarkdownDescription: d}, nil
	case contract.KindAny:
		return aschema.DynamicAttribute{Required: a.Required, Optional: a.Optional, Description: d, MarkdownDescription: d}, nil
	case contract.KindObject:
		nested, err := actionAttrs(a.Attrs)
		if err != nil {
			return nil, err
		}
		return aschema.SingleNestedAttribute{Attributes: nested, Required: a.Required, Optional: a.Optional, Description: d, MarkdownDescription: d}, nil
	case contract.KindArray, contract.KindMap:
		if isNestedObject(a) {
			nested, err := actionAttrs(a.Elem.Attrs)
			if err != nil {
				return nil, err
			}
			obj := aschema.NestedAttributeObject{Attributes: nested}
			if a.Kind == contract.KindArray {
				return aschema.ListNestedAttribute{NestedObject: obj, Required: a.Required, Optional: a.Optional, Description: d, MarkdownDescription: d}, nil
			}
			return aschema.MapNestedAttribute{NestedObject: obj, Required: a.Required, Optional: a.Optional, Description: d, MarkdownDescription: d}, nil
		}
		et, err := elemType(a.Elem)
		if err != nil {
			return nil, err
		}
		if a.Kind == contract.KindArray {
			return aschema.ListAttribute{ElementType: et, Required: a.Required, Optional: a.Optional, Description: d, MarkdownDescription: d}, nil
		}
		return aschema.MapAttribute{ElementType: et, Required: a.Required, Optional: a.Optional, Description: d, MarkdownDescription: d}, nil
	}
	return nil, fmt.Errorf("unsupported kind %s", a.Kind)
}
