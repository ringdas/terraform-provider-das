package contract

import (
	"fmt"
	"sort"
	"strings"
)

// Kind is the shape of a value, as Terraform can hold it.
type Kind int

const (
	KindString Kind = iota
	KindInteger
	KindNumber
	KindBool
	KindArray
	KindObject
	// KindMap is an object with arbitrary keys and one value schema.
	KindMap
	// KindAny is a value of no fixed shape (a setting's value: a string or
	// an object): a Terraform dynamic attribute.
	KindAny
)

func (k Kind) String() string {
	return [...]string{"string", "integer", "number", "bool", "array", "object", "map", "any"}[k]
}

// Schema is a JSON Schema normalized to what the provider maps.
type Schema struct {
	Kind        Kind
	Nullable    bool
	Description string
	Enum        []string
	Format      string
	// Pattern is the regular expression a string matches (examples are checked against it).
	Pattern   string
	WriteOnly bool
	// Items is the element of an array.
	Items *Schema
	// Properties of an object, sorted by name.
	Properties []Property
	// Values is the value schema of a map.
	Values *Schema
	// Component names the `#/components/schemas/<name>` it came from.
	Component string
}

// Property is one property of an object.
type Property struct {
	Name     string
	Required bool
	Schema   *Schema
}

// Property returns the named property, or nil.
func (s *Schema) Property(name string) *Property {
	if s == nil {
		return nil
	}
	for i := range s.Properties {
		if s.Properties[i].Name == name {
			return &s.Properties[i]
		}
	}
	return nil
}

// Keywords that only validate a value: the API validates them, the
// provider passes the value through. Any other keyword is one the provider
// does not know how to map, and Schema refuses it - so a new construct in
// a contract fails the provider's tests (provider_test.go), never a user's
// plan.
var validationOnly = map[string]bool{
	"description": true, "title": true, "pattern": true, "format": true, "default": true, "examples": true, "example": true,
	"minLength": true, "maxLength": true, "minimum": true, "maximum": true, "exclusiveMinimum": true, "exclusiveMaximum": true,
	"minItems": true, "maxItems": true, "uniqueItems": true, "minProperties": true, "maxProperties": true, "not": true,
	"readOnly": true, "deprecated": true, "multipleOf": true,
}

var structural = map[string]bool{
	"type": true, "enum": true, "$ref": true, "items": true, "properties": true, "required": true,
	"additionalProperties": true, "anyOf": true, "oneOf": true, "writeOnly": true, "const": true,
}

// Schema normalizes a raw JSON Schema of this contract.
func (s *Spec) Schema(raw any) (*Schema, error) {
	return s.schema(raw, 0)
}

func (s *Spec) schema(raw any, depth int) (*Schema, error) {
	if depth > 24 {
		return nil, fmt.Errorf("schema nested too deep (a recursive $ref?)")
	}
	m, ok := raw.(map[string]any)
	if !ok {
		if b, isBool := raw.(bool); isBool && b {
			return &Schema{Kind: KindAny}, nil
		}
		return nil, fmt.Errorf("a schema must be an object, got %T", raw)
	}
	if ref, ok := m["$ref"].(string); ok {
		target, err := s.pointer(ref)
		if err != nil {
			return nil, err
		}
		out, err := s.schema(target, depth+1)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", ref, err)
		}
		cp := *out
		if d, ok := m["description"].(string); ok && d != "" {
			cp.Description = d
		}
		if strings.HasPrefix(ref, "#/components/schemas/") {
			cp.Component = strings.TrimPrefix(ref, "#/components/schemas/")
		}
		return &cp, nil
	}
	for k := range m {
		if !validationOnly[k] && !structural[k] && !strings.HasPrefix(k, "x-") {
			return nil, fmt.Errorf("the provider does not map the JSON Schema keyword %q (teach internal/contract/schema.go, or keep the contract to the mapped constructs)", k)
		}
	}
	out := &Schema{}
	out.Description, _ = m["description"].(string)
	out.Format, _ = m["format"].(string)
	out.Pattern, _ = m["pattern"].(string)
	out.WriteOnly, _ = m["writeOnly"].(bool)

	for _, key := range []string{"anyOf", "oneOf"} {
		alts, ok := m[key].([]any)
		if !ok {
			continue
		}
		var kept []*Schema
		nullable := false
		for _, alt := range alts {
			if am, ok := alt.(map[string]any); ok && am["type"] == "null" && len(am) == 1 {
				nullable = true
				continue
			}
			a, err := s.schema(alt, depth+1)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", key, err)
			}
			kept = append(kept, a)
		}
		var picked *Schema
		if len(kept) == 1 {
			cp := *kept[0]
			picked = &cp
		} else {
			// Alternatives of different shapes: a value of no fixed shape.
			picked = &Schema{Kind: KindAny}
		}
		picked.Nullable = picked.Nullable || nullable
		if out.Description != "" {
			picked.Description = out.Description
		}
		picked.WriteOnly = picked.WriteOnly || out.WriteOnly
		return picked, nil
	}

	var types []string
	switch t := m["type"].(type) {
	case string:
		types = []string{t}
	case []any:
		for _, x := range t {
			if xs, ok := x.(string); ok {
				types = append(types, xs)
			}
		}
	case nil:
	default:
		return nil, fmt.Errorf("type must be a string or an array, got %T", t)
	}
	var nonNull []string
	for _, t := range types {
		if t == "null" {
			out.Nullable = true
		} else {
			nonNull = append(nonNull, t)
		}
	}
	if enum, ok := m["enum"].([]any); ok {
		for _, e := range enum {
			if es, ok := e.(string); ok {
				out.Enum = append(out.Enum, es)
			}
		}
		if len(nonNull) == 0 {
			nonNull = []string{"string"}
		}
	}
	if len(nonNull) == 0 {
		switch {
		case m["properties"] != nil:
			nonNull = []string{"object"}
		case m["items"] != nil:
			nonNull = []string{"array"}
		case m["const"] != nil:
			nonNull = []string{"string"}
		}
	}
	if len(nonNull) != 1 {
		out.Kind = KindAny
		return out, nil
	}
	switch nonNull[0] {
	case "string":
		out.Kind = KindString
	case "integer":
		out.Kind = KindInteger
	case "number":
		out.Kind = KindNumber
	case "boolean":
		out.Kind = KindBool
	case "array":
		out.Kind = KindArray
		items, ok := m["items"]
		if !ok {
			items = true
		}
		it, err := s.schema(items, depth+1)
		if err != nil {
			return nil, fmt.Errorf("items: %w", err)
		}
		out.Items = it
	case "object":
		props, _ := m["properties"].(map[string]any)
		required := map[string]bool{}
		if req, ok := m["required"].([]any); ok {
			for _, r := range req {
				if rs, ok := r.(string); ok {
					required[rs] = true
				}
			}
		}
		if len(props) == 0 {
			switch ap := m["additionalProperties"].(type) {
			case map[string]any:
				v, err := s.schema(ap, depth+1)
				if err != nil {
					return nil, fmt.Errorf("additionalProperties: %w", err)
				}
				out.Kind = KindMap
				out.Values = v
			default:
				// A free-form object (additionalProperties true or absent).
				out.Kind = KindAny
			}
			return out, nil
		}
		out.Kind = KindObject
		names := make([]string, 0, len(props))
		for n := range props {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, n := range names {
			ps, err := s.schema(props[n], depth+1)
			if err != nil {
				return nil, fmt.Errorf("property %s: %w", n, err)
			}
			out.Properties = append(out.Properties, Property{Name: n, Required: required[n], Schema: ps})
		}
	default:
		return nil, fmt.Errorf("unknown type %q", nonNull[0])
	}
	return out, nil
}
