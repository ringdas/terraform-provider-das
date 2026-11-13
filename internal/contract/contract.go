// Package contract reads a DAS API contract (an apps/<app>/openapi.json
// the provider embeds) and its Terraform surface (`x-das-terraform`,
// declared in apps/<app>/src/openapi/terraformSurface.ts and held to the
// contract by the TypeScript gate packages/platform/src/iac/iacCoverage.ts).
//
// Nothing here knows a DAS resource by name: every resource, data source,
// action and ephemeral resource of the provider is derived from what this
// package returns.
package contract

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Spec is one API contract.
type Spec struct {
	// Module is `x-das-module`: `org`, `delivery`, `audience`, `tag-manager`.
	Module string
	// Product is `x-das-product`: the host label of a product API.
	Product string
	Title   string
	Version string
	Surface Surface

	raw map[string]any
}

// Surface is `x-das-terraform`.
type Surface struct {
	Resources          []ResourceDef  `json:"resources"`
	DataSources        []DataDef      `json:"dataSources"`
	Actions            []ActionDef    `json:"actions"`
	EphemeralResources []EphemeralDef `json:"ephemeralResources"`
}

// ResourceDef mirrors TfResourceDef of the TypeScript registry.
type ResourceDef struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Create      string   `json:"create"`
	Read        string   `json:"read"`
	Update      string   `json:"update,omitempty"`
	Delete      string   `json:"delete,omitempty"`
	Destroy     *Destroy `json:"destroy,omitempty"`
	// ChangeInstead is what to do instead of changing an attribute the
	// update does not take (without a delete, it cannot be replaced either).
	ChangeInstead string `json:"changeInstead,omitempty"`
	List          string `json:"list,omitempty"`
	ListName      string `json:"listName,omitempty"`
	// Example is the registry example by wire name (internal/docgen).
	Example map[string]any `json:"example,omitempty"`
	// CreateBody is sent on every create and never an attribute (a minted
	// secret delivered to the secret store, not the answer).
	CreateBody map[string]any `json:"createBody,omitempty"`
	// CreateAnswer is the property of the create's answer holding the entity.
	CreateAnswer string `json:"createAnswer,omitempty"`
}

// Destroy says what destroy does when the API deletes nothing.
type Destroy struct {
	Via    string         `json:"via"`
	Body   map[string]any `json:"body,omitempty"`
	Reason string         `json:"reason"`
}

// DataDef is a read no resource owns.
type DataDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Read        string `json:"read"`
}

// ActionDef is an operation with nothing for Terraform to keep.
type ActionDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Invoke      string         `json:"invoke"`
	Example     map[string]any `json:"example,omitempty"`
}

// EphemeralDef is a short-lived credential opened per run.
type EphemeralDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Open        string `json:"open"`
}

// Operation is one operation of the contract, its references resolved.
type Operation struct {
	Ref         string // "POST /api/groups"
	Method      string // "POST"
	Path        string // "/api/groups"
	OperationID string
	Summary     string
	Description string
	// Layer is `x-das-layer`: "org" or "network".
	Layer       string
	PathParams  []Param
	QueryParams []Param
	// Body is the JSON request body schema (nil without a body).
	Body *Schema
	// Response is the schema of the first 2xx JSON answer (nil without one).
	Response *Schema
	// SuccessStatus is that answer's status, e.g. 201.
	SuccessStatus string
}

// Param is a path or query parameter.
type Param struct {
	Name        string
	In          string
	Required    bool
	Description string
	Schema      *Schema
}

var pathParamRe = regexp.MustCompile(`\{([^}]+)\}`)

// Parse reads one contract document.
func Parse(data []byte) (*Spec, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("contract: %w", err)
	}
	s := &Spec{raw: raw}
	s.Module, _ = raw["x-das-module"].(string)
	s.Product, _ = raw["x-das-product"].(string)
	if info, ok := raw["info"].(map[string]any); ok {
		s.Title, _ = info["title"].(string)
		s.Version, _ = info["version"].(string)
	}
	if s.Module == "" {
		return nil, fmt.Errorf("contract: no x-das-module")
	}
	surface, ok := raw["x-das-terraform"]
	if !ok {
		return nil, fmt.Errorf("contract %s: no x-das-terraform", s.Module)
	}
	buf, _ := json.Marshal(surface)
	dec := json.NewDecoder(strings.NewReader(string(buf)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&s.Surface); err != nil {
		return nil, fmt.Errorf("contract %s: x-das-terraform: %w", s.Module, err)
	}
	return s, nil
}

// Operation resolves "<METHOD> <path>".
func (s *Spec) Operation(ref string) (*Operation, error) {
	method, path, ok := strings.Cut(ref, " ")
	if !ok {
		return nil, fmt.Errorf("%s: %q is not \"<METHOD> <path>\"", s.Module, ref)
	}
	paths, _ := s.raw["paths"].(map[string]any)
	item, _ := paths[path].(map[string]any)
	opRaw, _ := item[strings.ToLower(method)].(map[string]any)
	if opRaw == nil {
		return nil, fmt.Errorf("%s: %s is not an operation of the contract", s.Module, ref)
	}
	op := &Operation{Ref: ref, Method: method, Path: path}
	op.OperationID, _ = opRaw["operationId"].(string)
	op.Summary, _ = opRaw["summary"].(string)
	op.Description, _ = opRaw["description"].(string)
	op.Layer, _ = opRaw["x-das-layer"].(string)
	params, _ := opRaw["parameters"].([]any)
	for _, p := range params {
		pm, err := s.resolveMap(p)
		if err != nil {
			return nil, fmt.Errorf("%s %s: %w", s.Module, ref, err)
		}
		in, _ := pm["in"].(string)
		if in != "path" && in != "query" {
			continue
		}
		param := Param{In: in}
		param.Name, _ = pm["name"].(string)
		param.Required, _ = pm["required"].(bool)
		param.Description, _ = pm["description"].(string)
		if sch, ok := pm["schema"]; ok {
			param.Schema, err = s.Schema(sch)
			if err != nil {
				return nil, fmt.Errorf("%s %s parameter %s: %w", s.Module, ref, param.Name, err)
			}
		}
		if in == "path" {
			op.PathParams = append(op.PathParams, param)
		} else {
			op.QueryParams = append(op.QueryParams, param)
		}
	}
	// Path parameters in path order, whatever order the spec lists them in.
	order := map[string]int{}
	for i, m := range pathParamRe.FindAllStringSubmatch(path, -1) {
		order[m[1]] = i
	}
	sort.SliceStable(op.PathParams, func(i, j int) bool { return order[op.PathParams[i].Name] < order[op.PathParams[j].Name] })
	if len(op.PathParams) != len(order) {
		return nil, fmt.Errorf("%s %s: the path names %d parameters, the operation declares %d", s.Module, ref, len(order), len(op.PathParams))
	}
	if rb, ok := opRaw["requestBody"]; ok {
		rbm, err := s.resolveMap(rb)
		if err != nil {
			return nil, err
		}
		if sch := jsonSchemaOf(rbm); sch != nil {
			op.Body, err = s.Schema(sch)
			if err != nil {
				return nil, fmt.Errorf("%s %s body: %w", s.Module, ref, err)
			}
		}
	}
	responses, _ := opRaw["responses"].(map[string]any)
	statuses := make([]string, 0, len(responses))
	for status := range responses {
		statuses = append(statuses, status)
	}
	sort.Strings(statuses)
	for _, status := range statuses {
		if len(status) != 3 || status[0] != '2' {
			continue
		}
		rm, err := s.resolveMap(responses[status])
		if err != nil {
			return nil, err
		}
		op.SuccessStatus = status
		if sch := jsonSchemaOf(rm); sch != nil {
			op.Response, err = s.Schema(sch)
			if err != nil {
				return nil, fmt.Errorf("%s %s response %s: %w", s.Module, ref, status, err)
			}
		}
		break
	}
	return op, nil
}

func jsonSchemaOf(m map[string]any) any {
	content, _ := m["content"].(map[string]any)
	media, _ := content["application/json"].(map[string]any)
	return media["schema"]
}

func (s *Spec) resolveMap(v any) (map[string]any, error) {
	for depth := 0; depth < 32; depth++ {
		m, ok := v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected an object, got %T", v)
		}
		ref, ok := m["$ref"].(string)
		if !ok {
			return m, nil
		}
		target, err := s.pointer(ref)
		if err != nil {
			return nil, err
		}
		v = target
	}
	return nil, fmt.Errorf("$ref chain too deep")
}

func (s *Spec) pointer(ref string) (any, error) {
	if !strings.HasPrefix(ref, "#/") {
		return nil, fmt.Errorf("only local $refs are supported: %s", ref)
	}
	var cur any = s.raw
	for _, part := range strings.Split(ref[2:], "/") {
		part = strings.ReplaceAll(strings.ReplaceAll(part, "~1", "/"), "~0", "~")
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("unresolved $ref %s", ref)
		}
		cur, ok = m[part]
		if !ok {
			return nil, fmt.Errorf("unresolved $ref %s", ref)
		}
	}
	return cur, nil
}
