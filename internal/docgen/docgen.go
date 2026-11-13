// Package docgen writes the provider's registry documentation (docs/,
// the layout registry.terraform.io and the OpenTofu registry read) from the
// same model the provider serves. Nothing is written by hand per resource:
// docs_test.go fails while docs/ differs from what Generate returns, so the
// documentation of a new feature lands with the feature's contract
// (`npm run gen:contracts`).
package docgen

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/ringdas/terraform-provider-das/internal/contract"
	"github.com/ringdas/terraform-provider-das/internal/model"
	"github.com/ringdas/terraform-provider-das/internal/provider"
)

// Generate returns every documentation file, by path under docs/.
func Generate(m *model.Provider) (map[string]string, error) {
	out := map[string]string{"index.md": index(m)}
	for _, r := range m.Resources {
		out["resources/"+short(r.TypeName)+".md"] = resourcePage(r)
	}
	for _, d := range m.DataSources {
		out["data-sources/"+short(d.TypeName)+".md"] = dataPage(d)
	}
	for _, a := range m.Actions {
		out["actions/"+short(a.TypeName)+".md"] = actionPage(a)
	}
	for _, e := range m.Ephemerals {
		out["ephemeral-resources/"+short(e.TypeName)+".md"] = ephemeralPage(e)
	}
	return out, nil
}

func short(typeName string) string { return strings.TrimPrefix(typeName, "das_") }

func subcategory(s *contract.Spec) string {
	switch s.Module {
	case "org":
		return "Organization"
	case "delivery":
		return "DAS Delivery"
	case "audience":
		return "DAS Audience"
	case "tag-manager":
		return "DAS Tag Manager"
	}
	return s.Title
}

func frontmatter(b *strings.Builder, title, sub, description string) {
	fmt.Fprintf(b, "---\npage_title: %q\nsubcategory: %q\ndescription: |-\n  %s\n---\n\n", title, sub, oneLine(description))
}

func oneLine(s string) string { return strings.Join(strings.Fields(s), " ") }

func index(m *model.Provider) string {
	var b strings.Builder
	frontmatter(&b, "Provider: DAS", "", provider.ProviderSchemaDescription)
	b.WriteString("# DAS Provider\n\n")
	b.WriteString(provider.ProviderSchemaDescription + "\n\n")
	b.WriteString("Every resource, data source, action and ephemeral resource below is derived from the API contracts of the DAS release this provider version ships with; ")
	b.WriteString("a newer DAS adds them to a newer provider version. Requires Terraform 1.11 or later (write-only attributes); actions need Terraform 1.14.\n\n")
	b.WriteString("Secrets never enter Terraform state: a client secret is a write-only attribute, a minted credential is an ephemeral resource.\n\n")
	b.WriteString("## Example Usage\n\n```terraform\n")
	b.WriteString("terraform {\n  required_providers {\n    das = {\n      source = \"ringdas/das\"\n    }\n  }\n}\n\n")
	b.WriteString("# The credential comes from DAS_API_KEY or DAS_TOKEN.\nprovider \"das\" {\n  org_domain = \"das.example.com\"\n}\n```\n\n")
	b.WriteString("## Contracts\n\n| API | Module | Version |\n| --- | --- | --- |\n")
	for _, s := range m.Specs {
		fmt.Fprintf(&b, "| %s | `%s` | %s |\n", s.Title, s.Module, s.Version)
	}
	b.WriteString("\n## Schema\n\n### Optional\n\n")
	for _, a := range provider.ProviderAttributes {
		sens := ""
		if a.Sensitive {
			sens = ", Sensitive"
		}
		fmt.Fprintf(&b, "- `%s` (String%s) %s Environment variable `%s`.\n", a.Name, sens, a.Description, a.Env)
	}
	b.WriteString("- `endpoints` (Map of String) Base URLs by API module (`org`, `delivery`, `audience`, `tag-manager`) instead of the hosts derived from `org_domain`: local development and tests.\n")
	return b.String()
}

func resourcePage(r *model.Resource) string {
	var b strings.Builder
	frontmatter(&b, short(r.TypeName)+" Resource - das", subcategory(r.Spec), r.Description)
	fmt.Fprintf(&b, "# %s (Resource)\n\n%s\n\n", r.TypeName, r.Description)
	fmt.Fprintf(&b, "API: `%s` creates, `%s` reads", r.Create.Ref, r.Read.Ref)
	if r.Update != nil {
		fmt.Fprintf(&b, ", `%s` updates", r.Update.Ref)
	}
	if r.Delete != nil {
		fmt.Fprintf(&b, ", `%s` deletes", r.Delete.Ref)
	}
	b.WriteString(".")
	if r.Def.Destroy != nil {
		fmt.Fprintf(&b, " Destroy: %s", r.Def.Destroy.Reason)
	}
	if r.Def.ChangeInstead != "" {
		fmt.Fprintf(&b, " %s", r.Def.ChangeInstead)
	}
	b.WriteString("\n\n## Example Usage\n\n```terraform\n")
	b.WriteString(exampleBlock("resource", r.TypeName, r.Attrs, r.Def.Example))
	b.WriteString("```\n\n")
	schemaSection(&b, r.Attrs)
	b.WriteString("## Import\n\n")
	id := "<" + strings.Join(r.Identity, ">/<") + ">"
	if len(r.Identity) == 0 {
		id = r.Def.Name
	}
	fmt.Fprintf(&b, "```terraform\nimport {\n  to = %s.example\n  id = %q\n}\n```\n", r.TypeName, id)
	return b.String()
}

func dataPage(d *model.DataSource) string {
	var b strings.Builder
	frontmatter(&b, short(d.TypeName)+" Data Source - das", subcategory(d.Spec), d.Description)
	fmt.Fprintf(&b, "# %s (Data Source)\n\n%s\n\nAPI: `%s`.\n\n## Example Usage\n\n```terraform\n%s```\n\n", d.TypeName, d.Description, d.Op.Ref, example("data", d.TypeName, exampleInputs(d)))
	schemaSection(&b, d.Attrs)
	return b.String()
}

// exampleInputs are a data source's inputs as its example sets them: a query
// parameter the API requires stays optional in Terraform (the model makes
// every query parameter optional), yet a read without it is a 400, so the
// example sets it, to the first alternative of its pattern when it has them.
func exampleInputs(d *model.DataSource) []*model.Attr {
	required := map[string]*contract.Schema{}
	for _, p := range d.Op.QueryParams {
		if p.Required {
			required[p.Name] = p.Schema
		}
	}
	if len(required) == 0 {
		return d.Attrs
	}
	out := make([]*model.Attr, len(d.Attrs))
	for i, a := range d.Attrs {
		s, ok := required[a.Wire]
		if a.Source != model.SourceQuery || !ok {
			out[i] = a
			continue
		}
		c := *a
		c.Required = true
		if len(c.Enum) == 0 && s != nil {
			c.Enum = patternAlternatives(s.Pattern)
		}
		out[i] = &c
	}
	return out
}

var alternativesRe = regexp.MustCompile(`^\^\(([a-z0-9_]+(?:\|[a-z0-9_]+)*)\)\$$`)

// patternAlternatives reads a pattern of plain words, `^(string|map)$`, as
// its words; nil for any other pattern.
func patternAlternatives(pattern string) []string {
	m := alternativesRe.FindStringSubmatch(pattern)
	if m == nil {
		return nil
	}
	return strings.Split(m[1], "|")
}

func actionPage(a *model.Action) string {
	var b strings.Builder
	frontmatter(&b, short(a.TypeName)+" Action - das", subcategory(a.Spec), a.Description)
	fmt.Fprintf(&b, "# %s (Action)\n\n%s Requires Terraform 1.14 or later.\n\nAPI: `%s`.\n\n## Example Usage\n\n```terraform\n", a.TypeName, a.Description, a.Op.Ref)
	b.WriteString(strings.Replace(exampleBlock("action", a.TypeName, a.Attrs, a.Def.Example), "{\n", "{\n  config {\n", 1))
	b.WriteString("```\n\nRun it with `terraform apply -invoke=action." + a.TypeName + ".example`, or from a resource's `lifecycle { action_trigger { ... } }`.\n\n")
	schemaSection(&b, a.Attrs)
	return b.String()
}

func ephemeralPage(e *model.Ephemeral) string {
	var b strings.Builder
	frontmatter(&b, short(e.TypeName)+" Ephemeral Resource - das", subcategory(e.Spec), e.Description)
	fmt.Fprintf(&b, "# %s (Ephemeral Resource)\n\n%s Requires Terraform 1.10 or later; its values never enter state or plan files.\n\nAPI: `%s`.\n\n## Example Usage\n\n```terraform\n%s```\n\n", e.TypeName, e.Description, e.Op.Ref, example("ephemeral", e.TypeName, e.Attrs))
	schemaSection(&b, e.Attrs)
	return b.String()
}

// ExampleAttr resolves a wire name of a surface example to its input
// attribute (a write-only secret's is `<name>_wo`), or nil.
func ExampleAttr(attrs []*model.Attr, wire string) *model.Attr {
	for _, a := range attrs {
		if a.Wire == wire && a.VersionOf == "" && (a.Required || a.Optional) {
			return a
		}
	}
	return nil
}

// exampleBlock writes the surface's example, or every required attribute.
func exampleBlock(kind, typeName string, attrs []*model.Attr, ex map[string]any) string {
	if len(ex) == 0 {
		return example(kind, typeName, attrs)
	}
	type line struct{ name, value string }
	var lines []line
	for wire, v := range ex {
		a := ExampleAttr(attrs, wire)
		if a == nil {
			continue // ExamplesAreValid (docs_test.go) refuses it
		}
		lines = append(lines, line{a.Name, hcl(v, "  ")})
		if a.WriteOnly {
			lines = append(lines, line{a.Name + "_version", "1"})
		}
	}
	sort.Slice(lines, func(i, j int) bool { return lines[i].name < lines[j].name })
	width := 0
	for _, l := range lines {
		if len(l.name) > width {
			width = len(l.name)
		}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s %q \"example\" {\n", kind, typeName)
	for _, l := range lines {
		fmt.Fprintf(&b, "  %-*s = %s\n", width, l.name, l.value)
	}
	if kind == "action" {
		b.WriteString("  }\n")
	}
	b.WriteString("}\n")
	return b.String()
}

// hcl renders a JSON value as an HCL expression.
func hcl(v any, indent string) string {
	switch x := v.(type) {
	case nil:
		return "null"
	case string:
		return fmt.Sprintf("%q", x)
	case []any:
		parts := make([]string, 0, len(x))
		for _, it := range x {
			parts = append(parts, hcl(it, indent))
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case map[string]any:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var b strings.Builder
		b.WriteString("{\n")
		for _, k := range keys {
			fmt.Fprintf(&b, "%s  %s = %s\n", indent, model.SnakeCase(k), hcl(x[k], indent+"  "))
		}
		b.WriteString(indent + "}")
		return b.String()
	default:
		return fmt.Sprint(x)
	}
}

// example writes a block with every required attribute.
func example(kind, typeName string, attrs []*model.Attr) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %q \"example\" {\n", kind, typeName)
	width := 0
	var req []*model.Attr
	for _, a := range attrs {
		if a.Required {
			req = append(req, a)
			if len(a.Name) > width {
				width = len(a.Name)
			}
		}
	}
	for _, a := range req {
		fmt.Fprintf(&b, "  %-*s = %s\n", width, a.Name, placeholder(a, "  "))
	}
	if kind == "action" {
		b.WriteString("  }\n")
	}
	b.WriteString("}\n")
	return b.String()
}

func placeholder(a *model.Attr, indent string) string {
	if len(a.Enum) > 0 {
		return fmt.Sprintf("%q", a.Enum[0])
	}
	switch a.Kind {
	case contract.KindString:
		switch {
		case strings.HasSuffix(a.Wire, "Id") || a.Wire == "id":
			return fmt.Sprintf("%q", "<"+a.Wire+">")
		case strings.Contains(strings.ToLower(a.Wire), "email") || a.Wire == "to":
			return `"someone@example.com"`
		case strings.Contains(strings.ToLower(a.Wire), "timezone"):
			return `"Europe/Warsaw"`
		case strings.Contains(strings.ToLower(a.Wire), "currency"):
			return `"PLN"`
		}
		return fmt.Sprintf("%q", "example-"+strings.ReplaceAll(a.Name, "_", "-"))
	case contract.KindInteger, contract.KindNumber:
		return "1"
	case contract.KindBool:
		return "true"
	case contract.KindArray:
		return "[" + placeholder(a.Elem, indent) + "]"
	case contract.KindMap:
		return "{ key = " + placeholder(a.Elem, indent) + " }"
	case contract.KindObject:
		var b strings.Builder
		b.WriteString("{\n")
		for _, na := range a.Attrs {
			if na.Required {
				fmt.Fprintf(&b, "%s  %s = %s\n", indent, na.Name, placeholder(na, indent+"  "))
			}
		}
		b.WriteString(indent + "}")
		return b.String()
	default:
		return `"example"`
	}
}

func schemaSection(b *strings.Builder, attrs []*model.Attr) {
	b.WriteString("## Schema\n\n")
	groups := []struct {
		title string
		pick  func(*model.Attr) bool
	}{
		{"Required", func(a *model.Attr) bool { return a.Required }},
		{"Optional", func(a *model.Attr) bool { return a.Optional }},
		{"Read-Only", func(a *model.Attr) bool { return a.Computed && !a.Optional && !a.Required }},
	}
	var nested []*model.Attr
	for _, g := range groups {
		var list []*model.Attr
		for _, a := range attrs {
			if g.pick(a) {
				list = append(list, a)
			}
		}
		if len(list) == 0 {
			continue
		}
		sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })
		fmt.Fprintf(b, "### %s\n\n", g.title)
		for _, a := range list {
			fmt.Fprintf(b, "- `%s` (%s) %s\n", a.Name, typeLabel(a), provider.Description(a))
			if hasNested(a) {
				nested = append(nested, a)
			}
		}
		b.WriteString("\n")
	}
	for _, a := range nested {
		nestedSection(b, []string{a.Name}, a)
	}
}

func hasNested(a *model.Attr) bool {
	return a.Kind == contract.KindObject || (a.Elem != nil && a.Elem.Kind == contract.KindObject)
}

// nestedSection documents an object's attributes, then each nested object
// inside it at every depth, anchored as tfplugindocs does
// (`nestedatt--rows--cells`).
func nestedSection(b *strings.Builder, path []string, a *model.Attr) {
	inner := a.Attrs
	if a.Elem != nil {
		inner = a.Elem.Attrs
	}
	fmt.Fprintf(b, "<a id=\"nestedatt--%s\"></a>\n### Nested Schema for `%s`\n\n", strings.Join(path, "--"), strings.Join(path, "."))
	var deeper []*model.Attr
	for _, na := range inner {
		fmt.Fprintf(b, "- `%s` (%s) %s\n", na.Name, typeLabel(na), provider.Description(na))
		if hasNested(na) {
			deeper = append(deeper, na)
		}
	}
	b.WriteString("\n")
	for _, na := range deeper {
		nestedSection(b, append(append([]string{}, path...), na.Name), na)
	}
}

func typeLabel(a *model.Attr) string {
	var t string
	switch a.Kind {
	case contract.KindString:
		t = "String"
	case contract.KindInteger:
		t = "Number, integer"
	case contract.KindNumber:
		t = "Number"
	case contract.KindBool:
		t = "Boolean"
	case contract.KindAny:
		t = "Dynamic"
	case contract.KindObject:
		t = "Attributes"
	case contract.KindArray:
		collection := "List"
		if a.Set {
			collection = "Set"
		}
		if a.Elem.Kind == contract.KindObject {
			t = "Attributes " + collection
		} else {
			t = collection + " of " + typeLabel(a.Elem)
		}
	case contract.KindMap:
		if a.Elem.Kind == contract.KindObject {
			t = "Attributes Map"
		} else {
			t = "Map of " + typeLabel(a.Elem)
		}
	}
	if a.Sensitive {
		t += ", Sensitive"
	}
	if a.WriteOnly {
		t += ", Write-only"
	}
	return t
}
