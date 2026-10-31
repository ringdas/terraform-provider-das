package provider

import (
	"errors"
	"fmt"
	"math/big"
	"net/url"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/ringdas/terraform-provider-das/internal/client"
	"github.com/ringdas/terraform-provider-das/internal/contract"
	"github.com/ringdas/terraform-provider-das/internal/model"
)

// values is a block's attribute values by Terraform name.
type values map[string]tftypes.Value

func attrsOf(v tftypes.Value) (values, error) {
	m, err := model.Attributes(v)
	return values(m), err
}

// stringOf is the string value of an attribute ("" when null or unknown).
func (v values) stringOf(name string) string {
	tv, ok := v[name]
	if !ok {
		return ""
	}
	return textOf(tv)
}

// textOf is a primitive value as it goes into a path or an import id: a
// string as it is, a number in decimal ("" when null, unknown or neither).
func textOf(tv tftypes.Value) string {
	if !tv.IsKnown() || tv.IsNull() {
		return ""
	}
	if tv.Type().Is(tftypes.Number) {
		var n big.Float
		if err := tv.As(&n); err != nil {
			return ""
		}
		return n.Text('f', -1)
	}
	var s string
	if err := tv.As(&s); err != nil {
		return ""
	}
	return s
}

// request builds the call of an operation from attribute values: path
// parameters, query parameters and the body properties `include` admits
// (writeOnly values come from `config`, the only place Terraform keeps them).
func request(spec *contract.Spec, op *contract.Operation, attrs []*model.Attr, vals, config values, pathValue func(param string) string, include func(*model.Attr) bool) (client.Request, error) {
	req := client.Request{Module: spec.Module, Product: spec.Product, Method: op.Method, Layer: op.Layer}
	params := map[string]string{}
	for _, p := range op.PathParams {
		params[p.Name] = pathValue(p.Name)
	}
	path, err := client.SubstitutePath(op.Path, params)
	if err != nil {
		return req, err
	}
	req.Path = path
	query := url.Values{}
	var body map[string]any
	if op.Body != nil {
		body = map[string]any{}
	}
	for _, a := range attrs {
		switch a.Source {
		case model.SourceQuery:
			j, err := model.FromValue(a, vals[a.Name])
			if err != nil {
				return req, fmt.Errorf("%s: %w", a.Name, err)
			}
			if j != nil {
				query.Set(a.Wire, fmt.Sprint(j))
			}
		case model.SourceBody:
			if body == nil || op.Body.Property(a.Wire) == nil || (include != nil && !include(a)) {
				continue
			}
			src := vals
			if a.WriteOnly {
				src = config
			}
			tv, ok := src[a.Name]
			if !ok {
				continue
			}
			j, err := model.FromValue(a, tv)
			if errors.Is(err, model.ErrUnknown) {
				continue // an optional value the API fills
			}
			if err != nil {
				return req, fmt.Errorf("%s: %w", a.Name, err)
			}
			if j != nil || a.SendNull {
				body[a.Wire] = j
			}
		}
	}
	if len(query) > 0 {
		req.Query = query
	}
	if body != nil {
		req.Body = body
	}
	return req, nil
}

// pathFrom resolves path parameters from attribute values by wire name.
func pathFrom(attrs []*model.Attr, vals values) func(string) string {
	return func(param string) string {
		for _, a := range attrs {
			if a.Wire == param && a.Source == model.SourcePath {
				return vals.stringOf(a.Name)
			}
		}
		return vals.stringOf(model.SnakeCase(param))
	}
}

// summary is a short rendering of an answer for an action's progress.
func summary(res *client.Response) string {
	if res == nil || res.JSON == nil {
		return fmt.Sprintf("answered %d", res.Status)
	}
	s := fmt.Sprintf("answered %d: %v", res.Status, res.JSON)
	if len(s) > 400 {
		s = s[:400] + "..."
	}
	return s
}
