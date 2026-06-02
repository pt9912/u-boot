// Package application holds the use-case orchestration of u-boot. It
// imports `internal/hexagon/domain` and `internal/hexagon/port` only;
// `internal/adapter/*` and external I/O libraries are forbidden by
// depguard (LH-FA-ARCH-002, LH-FA-ARCH-003).
package application

import (
	"bytes"
	"embed"
	"fmt"
	iofs "io/fs"
	"sort"
	"strings"
	"text/template"

	"github.com/pt9912/u-boot/internal/hexagon/application/managedblock"
)

//go:embed templates/*.tmpl templates/services/*.tmpl templates/devcontainer/*.tmpl
var templateFS embed.FS

// templateData is the projection of [domain.Project] that the
// embedded templates render against. Keeping it as a tiny dedicated
// struct (instead of passing the domain object directly) makes the
// template surface stable and lets the domain grow without
// inadvertently leaking new fields into the rendered output.
type templateData struct {
	Name string

	// ForwardPorts is the M7-T5 devcontainer-only field. Other
	// templates ignore it (text/template silently drops unreferenced
	// fields); only `templates/devcontainer/devcontainer.json.tmpl`
	// renders it, conditionally — empty/nil omits the JSON
	// `forwardPorts` key entirely per LH-FA-DEV-005 ("darf fehlen").
	ForwardPorts []int

	// Features is the slice-v1-devcontainer-features T3 devcontainer-
	// only field (LH-FA-DEV-003). Pre-sorted alphabetically by Source
	// in [collectDevcontainerFeatures] so the rendered JSONC is
	// byte-deterministic across runs; empty/nil omits the JSON
	// `features` key entirely.
	Features []devcontainerFeatureData
}

// fileTemplate maps an embedded template to its destination path
// (relative to [driving.InitProjectRequest.BaseDir]).
//
// Managed reports whether the template wraps its content in a
// `U-BOOT MANAGED BLOCK: init` marker (LH-SA-FILE-002) of the given
// Style. Managed templates support the LH-FA-INIT-005 §611–§614
// block-only re-init path; whole-file-managed templates
// (Managed=false, e.g. .gitignore) require --backup for re-init
// because the §619 backup-mandatory rule kicks in unconditionally.
type fileTemplate struct {
	Path         string
	TemplateName string
	Managed      bool
	Style        managedblock.Style
}

// fileTemplates returns the project files that [InitProjectService]
// generates from embedded templates. The order is the deterministic
// order in which they are written and listed in
// [driving.InitProjectResponse.Created]. Implemented as a function
// to avoid the gochecknoglobals false-positive on immutable list
// constants.
//
// The Managed flag tracks the LH-FA-INIT-005 §611 list of
// structured configuration files. .gitignore is intentionally left
// off the list (matches the spec verbatim); u-boot.yaml is handled
// outside this slice in [InitProjectService.executeUBootYAML] with
// the same fully-managed semantics.
func fileTemplates() []fileTemplate {
	return []fileTemplate{
		{Path: "README.md", TemplateName: "readme.md.tmpl", Managed: true, Style: managedblock.StyleHTMLComment},
		{Path: "CHANGELOG.md", TemplateName: "changelog.md.tmpl", Managed: true, Style: managedblock.StyleHTMLComment},
		{Path: "compose.yaml", TemplateName: "compose.yaml.tmpl", Managed: true, Style: managedblock.StyleHash},
		{Path: ".env.example", TemplateName: "env.example.tmpl", Managed: true, Style: managedblock.StyleHash},
		{Path: ".gitignore", TemplateName: "gitignore.tmpl"},
	}
}

// renderTemplate executes the named template against the given data
// and returns the rendered bytes.
func renderTemplate(name string, data templateData) ([]byte, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/"+name)
	if err != nil {
		return nil, fmt.Errorf("parse template %s: %w", name, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template %s: %w", name, err)
	}
	return buf.Bytes(), nil
}

// templateNames lists the embedded template filenames in sorted
// order. Walks the embedded tree so subdirectories
// (e.g. templates/services/) are included; returned names are
// relative to the templates/ root (e.g. "compose.yaml.tmpl",
// "services/postgres.compose.tmpl"). Used by the integrity self-test
// in `templates_test.go` and as the source for renderTemplate
// lookups.
func templateNames() ([]string, error) {
	var names []string
	err := iofs.WalkDir(templateFS, "templates", func(path string, d iofs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".tmpl") {
			return nil
		}
		rel := strings.TrimPrefix(path, "templates/")
		names = append(names, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}
