package model

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/ringdas/terraform-provider-das/internal/contract"
)

// Type is the Terraform type of an attribute.
func (a *Attr) Type() tftypes.Type {
	switch a.Kind {
	case contract.KindString:
		return tftypes.String
	case contract.KindInteger, contract.KindNumber:
		return tftypes.Number
	case contract.KindBool:
		return tftypes.Bool
	case contract.KindArray:
		if a.Set {
			return tftypes.Set{ElementType: a.Elem.Type()}
		}
		return tftypes.List{ElementType: a.Elem.Type()}
	case contract.KindMap:
		return tftypes.Map{ElementType: a.Elem.Type()}
	case contract.KindObject:
		return ObjectType(a.Attrs)
	default:
		return tftypes.DynamicPseudoType
	}
}

// ObjectType is the Terraform object type of a block of attributes.
func ObjectType(attrs []*Attr) tftypes.Object {
	types := map[string]tftypes.Type{}
	for _, a := range attrs {
		types[a.Name] = a.Type()
	}
	return tftypes.Object{AttributeTypes: types}
}

// ToValue converts a JSON value (as encoding/json decodes it, numbers as
// json.Number) to the attribute's Terraform value; nil is null.
func ToValue(a *Attr, v any) (tftypes.Value, error) {
	t := a.Type()
	if v == nil {
		return tftypes.NewValue(t, nil), nil
	}
	if a.JSON {
		b, err := json.Marshal(v)
		if err != nil {
			return tftypes.Value{}, fmt.Errorf("%s: %w", a.Wire, err)
		}
		return tftypes.NewValue(t, string(b)), nil
	}
	switch a.Kind {
	case contract.KindString:
		s, ok := v.(string)
		if !ok {
			return tftypes.Value{}, fmt.Errorf("%s: expected a string, got %T", a.Wire, v)
		}
		return tftypes.NewValue(t, s), nil
	case contract.KindInteger, contract.KindNumber:
		n, err := toBigFloat(v)
		if err != nil {
			return tftypes.Value{}, fmt.Errorf("%s: %w", a.Wire, err)
		}
		return tftypes.NewValue(t, n), nil
	case contract.KindBool:
		b, ok := v.(bool)
		if !ok {
			return tftypes.Value{}, fmt.Errorf("%s: expected a boolean, got %T", a.Wire, v)
		}
		return tftypes.NewValue(t, b), nil
	case contract.KindArray:
		items, ok := v.([]any)
		if !ok {
			return tftypes.Value{}, fmt.Errorf("%s: expected an array, got %T", a.Wire, v)
		}
		out := make([]tftypes.Value, 0, len(items))
		for _, it := range items {
			ev, err := ToValue(a.Elem, it)
			if err != nil {
				return tftypes.Value{}, fmt.Errorf("%s[]: %w", a.Wire, err)
			}
			out = append(out, ev)
		}
		return tftypes.NewValue(t, out), nil
	case contract.KindMap:
		m, ok := v.(map[string]any)
		if !ok {
			return tftypes.Value{}, fmt.Errorf("%s: expected an object, got %T", a.Wire, v)
		}
		out := map[string]tftypes.Value{}
		for k, it := range m {
			ev, err := ToValue(a.Elem, it)
			if err != nil {
				return tftypes.Value{}, fmt.Errorf("%s.%s: %w", a.Wire, k, err)
			}
			out[k] = ev
		}
		return tftypes.NewValue(t, out), nil
	case contract.KindObject:
		m, ok := v.(map[string]any)
		if !ok {
			return tftypes.Value{}, fmt.Errorf("%s: expected an object, got %T", a.Wire, v)
		}
		return ObjectValue(a.Attrs, m)
	default:
		return dynamicValue(v)
	}
}

// ObjectValue converts a JSON object to the Terraform object of a block;
// a property the block does not name is ignored, an attribute the object
// lacks is null.
func ObjectValue(attrs []*Attr, m map[string]any) (tftypes.Value, error) {
	vals := map[string]tftypes.Value{}
	for _, a := range attrs {
		raw, ok := m[a.Wire]
		if !ok || a.WriteOnly || a.VersionOf != "" {
			raw = nil
		}
		v, err := ToValue(a, raw)
		if err != nil {
			return tftypes.Value{}, err
		}
		vals[a.Name] = v
	}
	return tftypes.NewValue(ObjectType(attrs), vals), nil
}

// dynamicValue infers a Terraform value for JSON of no fixed shape.
func dynamicValue(v any) (tftypes.Value, error) {
	switch x := v.(type) {
	case nil:
		return tftypes.NewValue(tftypes.DynamicPseudoType, nil), nil
	case string:
		return tftypes.NewValue(tftypes.String, x), nil
	case bool:
		return tftypes.NewValue(tftypes.Bool, x), nil
	case json.Number, float64, int, int64:
		n, err := toBigFloat(x)
		if err != nil {
			return tftypes.Value{}, err
		}
		return tftypes.NewValue(tftypes.Number, n), nil
	case []any:
		vals := make([]tftypes.Value, 0, len(x))
		types := make([]tftypes.Type, 0, len(x))
		for _, it := range x {
			ev, err := dynamicValue(it)
			if err != nil {
				return tftypes.Value{}, err
			}
			vals = append(vals, ev)
			types = append(types, ev.Type())
		}
		return tftypes.NewValue(tftypes.Tuple{ElementTypes: types}, vals), nil
	case map[string]any:
		vals := map[string]tftypes.Value{}
		types := map[string]tftypes.Type{}
		for k, it := range x {
			ev, err := dynamicValue(it)
			if err != nil {
				return tftypes.Value{}, err
			}
			vals[k] = ev
			types[k] = ev.Type()
		}
		return tftypes.NewValue(tftypes.Object{AttributeTypes: types}, vals), nil
	default:
		return tftypes.Value{}, fmt.Errorf("unsupported JSON value %T", v)
	}
}

func toBigFloat(v any) (*big.Float, error) {
	switch x := v.(type) {
	case json.Number:
		f, _, err := big.ParseFloat(string(x), 10, 512, big.ToNearestEven)
		return f, err
	case float64:
		return big.NewFloat(x), nil
	case int:
		return new(big.Float).SetInt64(int64(x)), nil
	case int64:
		return new(big.Float).SetInt64(x), nil
	default:
		return nil, fmt.Errorf("expected a number, got %T", v)
	}
}

// ErrUnknown reports a value not known yet (it cannot be sent).
var ErrUnknown = fmt.Errorf("value not known yet")

// FromValue converts a Terraform value to JSON (nil for null).
func FromValue(a *Attr, v tftypes.Value) (any, error) {
	if !v.IsKnown() {
		return nil, ErrUnknown
	}
	if v.IsNull() {
		return nil, nil
	}
	switch a.Kind {
	case contract.KindString:
		var s string
		if err := v.As(&s); err != nil {
			return nil, err
		}
		if a.JSON {
			var j any
			dec := json.NewDecoder(strings.NewReader(s))
			dec.UseNumber()
			if err := dec.Decode(&j); err != nil {
				return nil, fmt.Errorf("%s: not JSON: %w", a.Name, err)
			}
			return j, nil
		}
		return s, nil
	case contract.KindInteger, contract.KindNumber:
		var f big.Float
		if err := v.As(&f); err != nil {
			return nil, err
		}
		return json.Number(numberString(&f, a.Kind == contract.KindInteger)), nil
	case contract.KindBool:
		var b bool
		err := v.As(&b)
		return b, err
	case contract.KindArray:
		var items []tftypes.Value
		if err := v.As(&items); err != nil {
			return nil, err
		}
		out := make([]any, 0, len(items))
		for _, it := range items {
			j, err := FromValue(a.Elem, it)
			if err != nil {
				return nil, err
			}
			out = append(out, j)
		}
		return out, nil
	case contract.KindMap:
		var m map[string]tftypes.Value
		if err := v.As(&m); err != nil {
			return nil, err
		}
		out := map[string]any{}
		for k, it := range m {
			j, err := FromValue(a.Elem, it)
			if err != nil {
				return nil, err
			}
			out[k] = j
		}
		return out, nil
	case contract.KindObject:
		var m map[string]tftypes.Value
		if err := v.As(&m); err != nil {
			return nil, err
		}
		out := map[string]any{}
		for _, na := range a.Attrs {
			if na.Source == SourceNone && na.Computed && !na.Optional && !na.Required {
				continue // computed only: never sent
			}
			if na.Computed && !m[na.Name].IsKnown() {
				continue // left to the API (a catalogue value's id): known after apply
			}
			j, err := FromValue(na, m[na.Name])
			if err != nil {
				return nil, err
			}
			if j != nil {
				out[na.Wire] = j
			}
		}
		return out, nil
	default:
		return dynamicJSON(v)
	}
}

func dynamicJSON(v tftypes.Value) (any, error) {
	if !v.IsKnown() {
		return nil, ErrUnknown
	}
	if v.IsNull() {
		return nil, nil
	}
	t := v.Type()
	switch {
	case t.Is(tftypes.String):
		var s string
		err := v.As(&s)
		return s, err
	case t.Is(tftypes.Bool):
		var b bool
		err := v.As(&b)
		return b, err
	case t.Is(tftypes.Number):
		var f big.Float
		if err := v.As(&f); err != nil {
			return nil, err
		}
		return json.Number(numberString(&f, false)), nil
	case t.Is(tftypes.List{}), t.Is(tftypes.Set{}), t.Is(tftypes.Tuple{}):
		var items []tftypes.Value
		if err := v.As(&items); err != nil {
			return nil, err
		}
		out := make([]any, 0, len(items))
		for _, it := range items {
			j, err := dynamicJSON(it)
			if err != nil {
				return nil, err
			}
			out = append(out, j)
		}
		return out, nil
	case t.Is(tftypes.Map{}), t.Is(tftypes.Object{}):
		var m map[string]tftypes.Value
		if err := v.As(&m); err != nil {
			return nil, err
		}
		out := map[string]any{}
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			j, err := dynamicJSON(m[k])
			if err != nil {
				return nil, err
			}
			out[k] = j
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported Terraform type %s", t)
	}
}

func numberString(f *big.Float, integer bool) string {
	if integer || f.IsInt() {
		i, _ := f.Int(nil)
		return i.String()
	}
	return f.Text('g', -1)
}

// Attributes splits a block's object value into its attributes.
func Attributes(v tftypes.Value) (map[string]tftypes.Value, error) {
	m := map[string]tftypes.Value{}
	if v.IsNull() || !v.IsKnown() {
		return m, nil
	}
	err := v.As(&m)
	return m, err
}
