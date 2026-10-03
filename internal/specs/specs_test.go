package specs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// THE CONTRACT FRESHNESS GATE: the provider embeds every API contract of
// the monorepo, byte for byte as released (apps/<app>/openapi.json). A
// changed contract without `npm run gen:contracts`, or a new API the
// provider does not embed, fails here - the provider cannot fall behind
// the CLI. Outside the monorepo (the public mirror) there is nothing to
// compare with, and the test says so.

func TestEmbeddedContractsEqualTheReleasedOnes(t *testing.T) {
	apps := filepath.Join("..", "..", "..")
	if _, err := os.Stat(filepath.Join(apps, "org-api", "openapi.json")); err != nil {
		t.Skip("not in the DAS monorepo: no released contracts to compare with")
	}
	entries, err := os.ReadDir(apps)
	if err != nil {
		t.Fatal(err)
	}
	var want []string
	for _, e := range entries {
		file := filepath.Join(apps, e.Name(), "openapi.json")
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		var head struct {
			Module string `json:"x-das-module"`
		}
		if err := json.Unmarshal(data, &head); err != nil || head.Module == "" {
			t.Fatalf("%s declares no x-das-module", file)
		}
		name := head.Module + Suffix
		want = append(want, name)
		embedded, err := Read(name)
		if err != nil {
			t.Errorf("%s is not embedded (run npm run gen:contracts)", name)
			continue
		}
		if string(embedded) != string(data) {
			t.Errorf("internal/specs/%s differs from apps/%s/openapi.json (run npm run gen:contracts)", name, e.Name())
		}
	}
	got, err := Files()
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(want)
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("embedded %v, the monorepo releases %v", got, want)
	}
}
