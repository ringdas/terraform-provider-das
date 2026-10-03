// Command docgen writes docs/ (the registry documentation) from the
// embedded contracts: `npm run -w terraform-provider-das gen:docs`, part of
// `npm run gen:contracts`. A file no longer generated is removed.
package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/ringdas/terraform-provider-das/internal/docgen"
	"github.com/ringdas/terraform-provider-das/internal/specs"
)

func main() {
	m, err := specs.Model()
	if err != nil {
		log.Fatal(err)
	}
	files, err := docgen.Generate(m)
	if err != nil {
		log.Fatal(err)
	}
	const dir = "docs"
	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(dir, p)
		if _, keep := files[filepath.ToSlash(rel)]; !keep {
			_ = os.Remove(p)
			fmt.Println("removed", p)
		}
		return nil
	})
	for rel, content := range files {
		p := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			log.Fatal(err)
		}
		if old, err := os.ReadFile(p); err == nil && string(old) == content {
			continue
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			log.Fatal(err)
		}
		fmt.Println("wrote", p)
	}
}
