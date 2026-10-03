package docgen_test

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/ringdas/terraform-provider-das/internal/docgen"
	"github.com/ringdas/terraform-provider-das/internal/model"
	"github.com/ringdas/terraform-provider-das/internal/specs"
)

// THE DOCUMENTATION GATE: docs/ (what the registries publish) is exactly
// what the embedded contracts generate - every resource, data source,
// action and ephemeral resource has its page, and no page outlives its
// type. Stale? `npm run gen:contracts` (or `npm run -w terraform-provider-das gen:docs`).

func TestDocsAreGenerated(t *testing.T) {
	m, err := specs.Model()
	if err != nil {
		t.Fatal(err)
	}
	want, err := docgen.Generate(m)
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join("..", "..", "docs")
	var have []string
	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			rel, _ := filepath.Rel(dir, p)
			have = append(have, filepath.ToSlash(rel))
		}
		return nil
	})
	sort.Strings(have)
	for _, rel := range have {
		if _, ok := want[rel]; !ok {
			t.Errorf("docs/%s documents nothing the provider serves (run npm run gen:contracts)", rel)
		}
	}
	for rel, content := range want {
		got, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(rel)))
		if err != nil {
			t.Errorf("docs/%s is missing (run npm run gen:contracts)", rel)
			continue
		}
		if string(got) != content {
			t.Errorf("docs/%s is stale (run npm run gen:contracts)", rel)
		}
	}
}

// THE EXAMPLE GATE: an example a surface declares (for what the contract
// cannot say, such as exactly one of email or groupName) names inputs of
// the type, with values of their type, and every required input - so the
// registry never shows a configuration the provider refuses.
func TestExamplesAreValid(t *testing.T) {
	m, err := specs.Model()
	if err != nil {
		t.Fatal(err)
	}
	check := func(typeName string, attrs []*model.Attr, ex map[string]any) {
		if len(ex) == 0 {
			return
		}
		for wire, v := range ex {
			a := docgen.ExampleAttr(attrs, wire)
			if a == nil {
				t.Errorf("%s: the example names %s, no input of the type", typeName, wire)
				continue
			}
			// Numbers as the API's JSON decoding gives them.
			buf, _ := json.Marshal(v)
			dec := json.NewDecoder(strings.NewReader(string(buf)))
			dec.UseNumber()
			var jv any
			_ = dec.Decode(&jv)
			if _, err := model.ToValue(a, jv); err != nil {
				t.Errorf("%s: the example's %s: %v", typeName, wire, err)
			}
			checkPatterns(t, typeName+"."+wire, a, jv)
		}
		for _, a := range attrs {
			if a.Required {
				if _, ok := ex[a.Wire]; !ok {
					t.Errorf("%s: the example leaves out the required %s", typeName, a.Name)
				}
			}
		}
	}
	for _, r := range m.Resources {
		check(r.TypeName, r.Attrs, r.Def.Example)
	}
	for _, a := range m.Actions {
		check(a.TypeName, a.Attrs, a.Def.Example)
	}
}

// checkPatterns holds an example's strings to the contract's patterns (a
// decimal amount, an ISO-8601 instant), nested objects and lists included.
// A pattern Go cannot compile (a lookahead) is left to the API.
func checkPatterns(t *testing.T, at string, a *model.Attr, v any) {
	t.Helper()
	switch x := v.(type) {
	case string:
		if a.Pattern == "" {
			return
		}
		re, err := regexp.Compile(a.Pattern)
		if err == nil && !re.MatchString(x) {
			t.Errorf("%s: %q does not match the contract's pattern %s", at, x, a.Pattern)
		}
	case map[string]any:
		for _, na := range a.Attrs {
			if nv, ok := x[na.Wire]; ok {
				checkPatterns(t, at+"."+na.Wire, na, nv)
			}
		}
	case []any:
		if a.Elem != nil {
			for _, it := range x {
				checkPatterns(t, at+"[]", a.Elem, it)
			}
		}
	}
}
