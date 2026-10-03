// Package model derives the provider's resources, data sources, actions
// and ephemeral resources from the API contracts (internal/contract):
// one attribute tree per type, used to build the Terraform schema, to
// convert between JSON and Terraform values, and to write the registry
// documentation. The derivation is deterministic - the same contracts
// always give the same provider.
package model

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/ringdas/terraform-provider-das/internal/contract"
)

// Source says where an input attribute goes in a request.
type Source int

const (
	SourceNone Source = iota // computed only
	SourcePath
	SourceQuery
	SourceBody
)

// Attr is one attribute (or, nameless, the element of a list or map).
type Attr struct {
	// Name is the Terraform name (snake_case); Wire the JSON name.
	Name        string
	Wire        string
	Description string
	Kind        contract.Kind
	Enum        []string
	Pattern     string
	Elem        *Attr
	Attrs       []*Attr

	Required  bool
	Optional  bool
	Computed  bool
	Sensitive bool
	WriteOnly bool

	Source Source
	// ForceNew: a change replaces the resource (an input the update does
	// not take, on a resource with a delete).
	ForceNew bool
	// Immutable: a change is refused at plan (an input the update does not
	// take, on a resource the API never deletes).
	Immutable bool
	// InUpdate: the update operation takes it.
	InUpdate bool
	// UpdateOnly: only the update operation takes it (set after create).
	UpdateOnly bool
	// VersionOf names the write-only attribute this version attribute
	// triggers: a change sends the secret again.
	VersionOf string
	// JSON: a value of no fixed shape inside a list or map (Terraform has
	// no dynamic type there), held as a JSON-encoded string `<name>_json`.
	JSON bool
	// SendNull: a body property the API requires but accepts as null (an
	// email sender's Reply-To): optional in Terraform, and left out it is
	// sent as an explicit null - never omitted, never filled by the API.
	SendNull bool
}

// jsonInCollection turns a value of no fixed shape inside a collection
// into a JSON-encoded string attribute.
func jsonInCollection(a *Attr, inCollection bool) {
	if a.Kind != contract.KindAny || !inCollection {
		return
	}
	a.Kind = contract.KindString
	a.JSON = true
	if a.Name != "" {
		a.Name += "_json"
	}
	a.Description = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(a.Description), ".") + ". JSON-encoded (`jsondecode()` it).")
}

// Find returns the top-level attribute with the Terraform name, or nil.
func Find(attrs []*Attr, name string) *Attr {
	for _, a := range attrs {
		if a.Name == name {
			return a
		}
	}
	return nil
}

func findWire(attrs []*Attr, wire string) *Attr {
	for _, a := range attrs {
		if a.Wire == wire && a.VersionOf == "" {
			return a
		}
	}
	return nil
}

// Resource is one `das_<name>` resource.
type Resource struct {
	TypeName    string
	Spec        *contract.Spec
	Def         contract.ResourceDef
	Create      *contract.Operation
	Read        *contract.Operation
	Update      *contract.Operation
	Delete      *contract.Operation
	List        *contract.Operation
	Attrs       []*Attr
	Description string
	// Identity is the read path's parameters, in order: the import id is
	// their values joined by "/".
	Identity []string
	// IDParam is the identity parameter the create's answer gives as `id`
	// ("" when every identity parameter is an input).
	IDParam string
}

// CreatedEntity is the schema of the entity the create answers: the answer,
// or its CreateAnswer property.
func (r *Resource) CreatedEntity() *contract.Schema {
	if r.Def.CreateAnswer == "" {
		return r.Create.Response
	}
	if p := r.Create.Response.Property(r.Def.CreateAnswer); p != nil {
		return p.Schema
	}
	return nil
}

// IdentityAttr is the attribute holding an identity parameter's value.
func (r *Resource) IdentityAttr(param string) string {
	if param == r.IDParam {
		return "id"
	}
	return SnakeCase(param)
}

// DataSource is one `data "das_<name>"`.
type DataSource struct {
	TypeName    string
	Spec        *contract.Spec
	Op          *contract.Operation
	Attrs       []*Attr
	Description string
	// Page: the answer is `{items, nextCursor}`, followed to its end into `items`.
	Page bool
	// IDParam maps the `id` input to this path parameter (a resource's read).
	IDParam string
}

// Action is one `action "das_<name>"`.
type Action struct {
	Def         contract.ActionDef
	TypeName    string
	Spec        *contract.Spec
	Op          *contract.Operation
	Attrs       []*Attr
	Description string
}

// Ephemeral is one `ephemeral "das_<name>"`.
type Ephemeral struct {
	TypeName    string
	Spec        *contract.Spec
	Op          *contract.Operation
	Attrs       []*Attr
	Description string
}

// Provider is everything derived from the contracts.
type Provider struct {
	Specs       []*contract.Spec
	Resources   []*Resource
	DataSources []*DataSource
	Actions     []*Action
	Ephemerals  []*Ephemeral
}

// Prefix is the type-name prefix of a module: `das_` for the
// Organisation API, `das_<module>_` for a product.
func Prefix(module string) string {
	if module == "org" {
		return "das_"
	}
	return "das_" + strings.ReplaceAll(module, "-", "_") + "_"
}

// Build derives the provider from the contracts.
func Build(specs []*contract.Spec) (*Provider, error) {
	p := &Provider{Specs: specs}
	seen := map[string]string{}
	taken := func(kind, name, by string) error {
		key := kind + " " + name
		if prev, ok := seen[key]; ok {
			return fmt.Errorf("%s %s is derived twice (%s and %s)", kind, name, prev, by)
		}
		seen[key] = by
		return nil
	}
	for _, spec := range specs {
		prefix := Prefix(spec.Module)
		for _, def := range spec.Surface.Resources {
			r, err := buildResource(spec, prefix, def)
			if err != nil {
				return nil, fmt.Errorf("resource %s%s: %w", prefix, def.Name, err)
			}
			if err := taken("resource", r.TypeName, spec.Module); err != nil {
				return nil, err
			}
			p.Resources = append(p.Resources, r)
			ds, err := dataFromResource(r)
			if err != nil {
				return nil, fmt.Errorf("data source %s: %w", r.TypeName, err)
			}
			if err := taken("data source", ds.TypeName, spec.Module); err != nil {
				return nil, err
			}
			p.DataSources = append(p.DataSources, ds)
			if r.List != nil {
				name := def.ListName
				if name == "" {
					name = def.Name + "s"
				}
				lds, err := buildData(spec, prefix+name, "Every "+strings.TrimPrefix(r.TypeName, "das_")+", all pages: "+lowerFirst(FirstSentence(def.Description)), r.List, "")
				if err != nil {
					return nil, fmt.Errorf("data source %s%s: %w", prefix, name, err)
				}
				if err := taken("data source", lds.TypeName, spec.Module); err != nil {
					return nil, err
				}
				p.DataSources = append(p.DataSources, lds)
			}
		}
		for _, def := range spec.Surface.DataSources {
			op, err := spec.Operation(def.Read)
			if err != nil {
				return nil, err
			}
			ds, err := buildData(spec, prefix+def.Name, def.Description, op, "")
			if err != nil {
				return nil, fmt.Errorf("data source %s%s: %w", prefix, def.Name, err)
			}
			if err := taken("data source", ds.TypeName, spec.Module); err != nil {
				return nil, err
			}
			p.DataSources = append(p.DataSources, ds)
		}
		for _, def := range spec.Surface.Actions {
			op, err := spec.Operation(def.Invoke)
			if err != nil {
				return nil, err
			}
			a := &Action{Def: def, TypeName: prefix + def.Name, Spec: spec, Op: op, Description: def.Description}
			if a.Attrs, err = requestInputs(op); err != nil {
				return nil, fmt.Errorf("action %s: %w", a.TypeName, err)
			}
			if err := taken("action", a.TypeName, spec.Module); err != nil {
				return nil, err
			}
			p.Actions = append(p.Actions, a)
		}
		for _, def := range spec.Surface.EphemeralResources {
			op, err := spec.Operation(def.Open)
			if err != nil {
				return nil, err
			}
			e := &Ephemeral{TypeName: prefix + def.Name, Spec: spec, Op: op, Description: def.Description}
			if e.Attrs, err = requestInputs(op); err != nil {
				return nil, fmt.Errorf("ephemeral resource %s: %w", e.TypeName, err)
			}
			if op.Response != nil {
				if err := mergeResponse(&e.Attrs, op.Response); err != nil {
					return nil, fmt.Errorf("ephemeral resource %s: %w", e.TypeName, err)
				}
			}
			markSensitive(e.Attrs)
			if err := taken("ephemeral resource", e.TypeName, spec.Module); err != nil {
				return nil, err
			}
			p.Ephemerals = append(p.Ephemerals, e)
		}
	}
	sort.Slice(p.Resources, func(i, j int) bool { return p.Resources[i].TypeName < p.Resources[j].TypeName })
	sort.Slice(p.DataSources, func(i, j int) bool { return p.DataSources[i].TypeName < p.DataSources[j].TypeName })
	sort.Slice(p.Actions, func(i, j int) bool { return p.Actions[i].TypeName < p.Actions[j].TypeName })
	sort.Slice(p.Ephemerals, func(i, j int) bool { return p.Ephemerals[i].TypeName < p.Ephemerals[j].TypeName })
	return p, nil
}

func buildResource(spec *contract.Spec, prefix string, def contract.ResourceDef) (*Resource, error) {
	r := &Resource{TypeName: prefix + def.Name, Spec: spec, Def: def, Description: def.Description}
	var err error
	if r.Create, err = spec.Operation(def.Create); err != nil {
		return nil, err
	}
	if r.Read, err = spec.Operation(def.Read); err != nil {
		return nil, err
	}
	for _, x := range []struct {
		ref string
		op  **contract.Operation
	}{{def.Update, &r.Update}, {def.Delete, &r.Delete}, {def.List, &r.List}} {
		if x.ref == "" {
			continue
		}
		if *x.op, err = spec.Operation(x.ref); err != nil {
			return nil, err
		}
	}
	if r.Delete == nil && def.Destroy == nil {
		return nil, fmt.Errorf("no delete and no destroy override")
	}
	if r.Read.Response == nil || r.Read.Response.Kind != contract.KindObject {
		return nil, fmt.Errorf("the read %s answers no object", r.Read.Ref)
	}
	createParams := map[string]bool{}
	for _, p := range r.Create.PathParams {
		createParams[p.Name] = true
	}
	for _, p := range r.Read.PathParams {
		r.Identity = append(r.Identity, p.Name)
		if createParams[p.Name] || r.Create.Body.Property(p.Name) != nil {
			continue
		}
		if r.IDParam != "" {
			return nil, fmt.Errorf("%s and %s cannot both come from the create's answer", r.IDParam, p.Name)
		}
		if r.CreatedEntity().Property("id") == nil {
			return nil, fmt.Errorf("%s must come from the create's answer, which has no id", p.Name)
		}
		r.IDParam = p.Name
	}

	updateBody := map[string]bool{}
	if r.Update != nil && r.Update.Body != nil {
		for _, p := range r.Update.Body.Properties {
			updateBody[p.Name] = true
		}
	}
	changeable := func(a *Attr, inUpdate bool) {
		a.InUpdate = inUpdate
		if inUpdate {
			return
		}
		if r.Delete != nil {
			a.ForceNew = true
		} else {
			a.Immutable = true
		}
	}

	r.Attrs = append(r.Attrs, &Attr{Name: "id", Wire: "id", Kind: contract.KindString, Computed: true, Description: idDescription(r)})
	for _, p := range r.Create.PathParams {
		a := paramAttr(p)
		a.Required = true
		changeable(a, false)
		// An identity parameter can never change in place.
		r.Attrs = append(r.Attrs, a)
	}
	if r.Create.Body != nil {
		if r.Create.Body.Kind != contract.KindObject {
			return nil, fmt.Errorf("the create body of %s is not an object", r.Create.Ref)
		}
		for _, prop := range r.Create.Body.Properties {
			if _, fixed := def.CreateBody[prop.Name]; fixed {
				continue // always sent, never an attribute
			}
			attrs, err := bodyAttrs(prop, true, true)
			if err != nil {
				return nil, err
			}
			for _, a := range attrs {
				if a.VersionOf != "" {
					// The version of a write-only secret: changing it sends
					// the secret again, in place when the update takes it.
					changeable(a, updateBody[prop.Name])
					a.InUpdate = false
				} else {
					changeable(a, updateBody[prop.Name])
				}
				r.Attrs = append(r.Attrs, a)
			}
		}
	}
	if r.Update != nil && r.Update.Body != nil {
		for _, prop := range r.Update.Body.Properties {
			if r.Create.Body.Property(prop.Name) != nil || findWire(r.Attrs, prop.Name) != nil {
				continue
			}
			attrs, err := bodyAttrs(prop, false, true)
			if err != nil {
				return nil, err
			}
			for _, a := range attrs {
				a.InUpdate = true
				a.UpdateOnly = true
				a.Required = false
				if !a.WriteOnly && a.VersionOf == "" {
					a.Optional, a.Computed = true, true
				}
				r.Attrs = append(r.Attrs, a)
			}
		}
	}
	if err := mergeResponse(&r.Attrs, r.Read.Response); err != nil {
		return nil, err
	}
	if r.IDParam != "" {
		// The id the create answers is the identity parameter: no second attribute.
		r.Attrs = without(r.Attrs, SnakeCase(r.IDParam))
	}
	return r, checkNames(r.Attrs)
}

func idDescription(r *Resource) string {
	if r.IDParam != "" {
		return "The id the API assigned."
	}
	if len(r.Identity) == 0 {
		return "A constant: the Organization has one " + strings.TrimPrefix(r.TypeName, "das_") + "."
	}
	return "The identity, " + strings.Join(r.Identity, "/") + " joined by `/` (the import id)."
}

func without(attrs []*Attr, name string) []*Attr {
	out := attrs[:0]
	for _, a := range attrs {
		if a.Name != name || a.Source == SourcePath {
			out = append(out, a)
		}
	}
	return out
}

// dataFromResource is the data source reading one entity by its identity.
func dataFromResource(r *Resource) (*DataSource, error) {
	// What the entity is (the first sentence), never what the resource's
	// destroy does: a data source only reads.
	description := FirstSentence(r.Description)
	if len(r.Identity) > 0 {
		description += " Read by its identity."
	}
	return buildData(r.Spec, r.TypeName, description, r.Read, r.IDParam)
}

func buildData(spec *contract.Spec, typeName, description string, op *contract.Operation, idParam string) (*DataSource, error) {
	ds := &DataSource{TypeName: typeName, Spec: spec, Op: op, Description: description, IDParam: idParam}
	for _, p := range op.PathParams {
		a := paramAttr(p)
		a.Required = true
		if p.Name == idParam {
			a.Name = "id"
		}
		ds.Attrs = append(ds.Attrs, a)
	}
	for _, p := range op.QueryParams {
		if p.Name == "limit" || p.Name == "cursor" {
			continue
		}
		a := paramAttr(p)
		a.Source = SourceQuery
		a.Optional = true
		ds.Attrs = append(ds.Attrs, a)
	}
	if op.Response == nil {
		return nil, fmt.Errorf("%s answers no JSON", op.Ref)
	}
	resp := op.Response
	if items := resp.Property("items"); items != nil && resp.Property("nextCursor") != nil && items.Schema.Kind == contract.KindArray {
		ds.Page = true
		elem, err := computedAttrIn(contract.Property{Name: "items", Schema: items.Schema.Items}, true)
		if err != nil {
			return nil, err
		}
		elem.Name, elem.Wire = "", ""
		ds.Attrs = append(ds.Attrs, &Attr{Name: "items", Wire: "items", Kind: contract.KindArray, Elem: elem, Computed: true, Description: "Every item, all pages followed."})
	} else if err := mergeResponse(&ds.Attrs, resp); err != nil {
		return nil, err
	}
	if idParam != "" {
		// The response's own `id` is the input.
		ds.Attrs = dedupeID(ds.Attrs)
	}
	return ds, checkNames(ds.Attrs)
}

func dedupeID(attrs []*Attr) []*Attr {
	out := attrs[:0]
	seen := false
	for _, a := range attrs {
		if a.Name == "id" {
			if seen {
				continue
			}
			seen = true
		}
		out = append(out, a)
	}
	return out
}

// requestInputs are an operation's path parameters and body properties.
func requestInputs(op *contract.Operation) ([]*Attr, error) {
	var out []*Attr
	for _, p := range op.PathParams {
		a := paramAttr(p)
		a.Required = true
		out = append(out, a)
	}
	for _, p := range op.QueryParams {
		a := paramAttr(p)
		a.Source = SourceQuery
		a.Required = p.Required
		a.Optional = !p.Required
		out = append(out, a)
	}
	if op.Body != nil {
		if op.Body.Kind != contract.KindObject {
			return nil, fmt.Errorf("the body of %s is not an object", op.Ref)
		}
		for _, prop := range op.Body.Properties {
			attrs, err := bodyAttrs(prop, true, false)
			if err != nil {
				return nil, err
			}
			out = append(out, attrs...)
		}
	}
	return out, checkNames(out)
}

func paramAttr(p contract.Param) *Attr {
	a := &Attr{Name: SnakeCase(p.Name), Wire: p.Name, Kind: contract.KindString, Description: p.Description, Source: SourcePath}
	if p.Schema != nil {
		a.Kind = p.Schema.Kind
		a.Enum = p.Schema.Enum
		if a.Description == "" {
			a.Description = p.Schema.Description
		}
	}
	return a
}

// bodyAttrs maps one body property: one attribute, or for a write-only
// secret the write-only attribute and its version.
func bodyAttrs(prop contract.Property, requiredHere, resource bool) ([]*Attr, error) {
	a, err := inputAttr(prop)
	if err != nil {
		return nil, err
	}
	a.Source = SourceBody
	switch {
	case prop.Required && requiredHere && prop.Schema.Nullable:
		a.Optional = true
		a.SendNull = true
		a.Description = strings.TrimSuffix(strings.TrimSpace(a.Description), ".") + ". Leave it out for null."
	case prop.Required && requiredHere:
		a.Required = true
	default:
		a.Optional = true
	}
	if !prop.Schema.WriteOnly {
		return []*Attr{a}, nil
	}
	if !resource {
		// An action's or ephemeral resource's configuration is never stored.
		a.Sensitive = true
		return []*Attr{a}, nil
	}
	a.Name += "_wo"
	a.WriteOnly = true
	a.Sensitive = true
	a.Description = strings.TrimSuffix(strings.TrimSpace(a.Description), ".") + ". Write-only: never stored in Terraform state (Terraform >= 1.11); change `" + a.Name + "_version` to send a new value."
	version := &Attr{
		Name:        a.Name + "_version",
		Wire:        prop.Name,
		Kind:        contract.KindInteger,
		Optional:    true,
		Source:      SourceNone,
		VersionOf:   a.Name,
		Description: "Any number; change it to send `" + a.Name + "` again (a rotation).",
	}
	return []*Attr{a, version}, nil
}

func inputAttr(prop contract.Property) (*Attr, error) {
	return inputAttrIn(prop, false)
}

func inputAttrIn(prop contract.Property, inCollection bool) (*Attr, error) {
	s := prop.Schema
	a := &Attr{Name: SnakeCase(prop.Name), Wire: prop.Name, Kind: s.Kind, Enum: s.Enum, Pattern: s.Pattern, Description: s.Description}
	jsonInCollection(a, inCollection)
	switch s.Kind {
	case contract.KindArray:
		elem, err := inputAttrIn(contract.Property{Name: "", Schema: s.Items, Required: true}, true)
		if err != nil {
			return nil, err
		}
		elem.Name, elem.Wire = "", ""
		a.Elem = elem
	case contract.KindMap:
		elem, err := inputAttrIn(contract.Property{Name: "", Schema: s.Values, Required: true}, true)
		if err != nil {
			return nil, err
		}
		elem.Name, elem.Wire = "", ""
		a.Elem = elem
	case contract.KindObject:
		for _, p := range s.Properties {
			na, err := inputAttrIn(p, inCollection)
			if err != nil {
				return nil, err
			}
			na.Source = SourceBody
			if p.Required {
				na.Required = true
			} else {
				na.Optional = true
			}
			a.Attrs = append(a.Attrs, na)
		}
	}
	return a, nil
}

func computedAttr(prop contract.Property) (*Attr, error) {
	return computedAttrIn(prop, false)
}

func computedAttrIn(prop contract.Property, inCollection bool) (*Attr, error) {
	s := prop.Schema
	a := &Attr{Name: SnakeCase(prop.Name), Wire: prop.Name, Kind: s.Kind, Enum: s.Enum, Description: s.Description, Computed: true}
	jsonInCollection(a, inCollection)
	switch s.Kind {
	case contract.KindArray, contract.KindMap:
		inner := s.Items
		if s.Kind == contract.KindMap {
			inner = s.Values
		}
		elem, err := computedAttrIn(contract.Property{Schema: inner}, true)
		if err != nil {
			return nil, err
		}
		elem.Name, elem.Wire, elem.Computed = "", "", false
		a.Elem = elem
	case contract.KindObject:
		for _, p := range s.Properties {
			na, err := computedAttrIn(p, inCollection)
			if err != nil {
				return nil, err
			}
			a.Attrs = append(a.Attrs, na)
		}
	}
	return a, nil
}

// mergeResponse adds the answer's properties: an input it echoes becomes
// optional and computed (the API may fill it), the rest computed.
func mergeResponse(attrs *[]*Attr, resp *contract.Schema) error {
	if resp.Kind != contract.KindObject {
		return fmt.Errorf("the answer is not an object")
	}
	for _, prop := range resp.Properties {
		if prop.Schema.Format == "password" {
			continue // secret material is never an attribute (P12)
		}
		if existing := findWire(*attrs, prop.Name); existing != nil {
			if existing.WriteOnly {
				continue
			}
			if existing.Optional && !existing.SendNull {
				existing.Computed = true
			}
			if existing.Kind == contract.KindObject && prop.Schema.Kind == contract.KindObject {
				if err := mergeResponse(&existing.Attrs, prop.Schema); err != nil {
					return err
				}
			}
			continue
		}
		a, err := computedAttr(prop)
		if err != nil {
			return err
		}
		*attrs = append(*attrs, a)
	}
	return nil
}

var sensitiveWords = []string{"token", "secret", "password"}

func markSensitive(attrs []*Attr) {
	for _, a := range attrs {
		for _, w := range sensitiveWords {
			if strings.Contains(a.Name, w) && a.Kind == contract.KindString {
				a.Sensitive = true
			}
		}
	}
}

// Names Terraform reserves at the top level of a block.
var reserved = map[string]bool{"count": true, "depends_on": true, "for_each": true, "lifecycle": true, "provider": true, "providers": true, "connection": true, "provisioner": true}

func checkNames(attrs []*Attr) error {
	seen := map[string]bool{}
	for _, a := range attrs {
		if reserved[a.Name] {
			return fmt.Errorf("the attribute name %q is reserved by Terraform (rename the wire field %q)", a.Name, a.Wire)
		}
		if seen[a.Name] {
			return fmt.Errorf("two attributes named %q", a.Name)
		}
		seen[a.Name] = true
	}
	return nil
}

// SnakeCase maps a camelCase wire name to a Terraform attribute name.
func SnakeCase(s string) string {
	var b strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 && (unicode.IsLower(runes[i-1]) || unicode.IsDigit(runes[i-1]) || (i+1 < len(runes) && unicode.IsLower(runes[i+1]))) {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		if r == '-' || r == '.' || r == ' ' {
			b.WriteByte('_')
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// FirstSentence is a description's first sentence (up to the first ". ").
func FirstSentence(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.Index(s, ". "); i >= 0 {
		return s[:i+1]
	}
	return s
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
