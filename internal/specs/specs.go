// Package specs embeds the API contracts the provider is built from: a
// copy of every apps/<app>/openapi.json of the monorepo, written by
// `npm run -w terraform-provider-das sync:specs` (part of
// `npm run gen:contracts`). specs_test.go fails when a copy is stale or
// an API is missing, so a provider is never released without a feature its
// APIs have.
package specs

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/ringdas/terraform-provider-das/internal/contract"
	"github.com/ringdas/terraform-provider-das/internal/model"
)

// Suffix of an embedded contract: `<x-das-module>.openapi.json`.
const Suffix = ".openapi.json"

//go:embed *.openapi.json
var files embed.FS

// Files lists the embedded contracts, sorted.
func Files() ([]string, error) {
	entries, err := fs.ReadDir(files, ".")
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), Suffix) {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out, nil
}

// Read returns one embedded contract.
func Read(name string) ([]byte, error) { return files.ReadFile(name) }

// Contracts parses every embedded contract.
func Contracts() ([]*contract.Spec, error) {
	names, err := Files()
	if err != nil {
		return nil, err
	}
	var out []*contract.Spec
	for _, name := range names {
		data, err := files.ReadFile(name)
		if err != nil {
			return nil, err
		}
		s, err := contract.Parse(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if s.Module+Suffix != name {
			return nil, fmt.Errorf("%s declares the module %q", name, s.Module)
		}
		out = append(out, s)
	}
	return out, nil
}

// Model derives the provider from the embedded contracts.
func Model() (*model.Provider, error) {
	specs, err := Contracts()
	if err != nil {
		return nil, err
	}
	return model.Build(specs)
}
