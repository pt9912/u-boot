// Package application holds the use-case orchestration of u-boot. It
// imports `internal/hexagon/domain` and `internal/hexagon/port` only;
// `internal/adapter/*` and external I/O libraries are forbidden by
// depguard (LH-FA-ARCH-002, LH-FA-ARCH-003).
package application

import (
	"bytes"
	"embed"
	"fmt"
	"sort"
	"text/template"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// templateData is the projection of [domain.Project] that the
// embedded templates render against. Keeping it as a tiny dedicated
// struct (instead of passing the domain object directly) makes the
// template surface stable and lets the domain grow without
// inadvertently leaking new fields into the rendered output.
type templateData struct {
	Name string
}

// fileTemplate maps an embedded template to its destination path
// (relative to [driving.InitProjectRequest.BaseDir]).
type fileTemplate struct {
	Path         string
	TemplateName string
}

// fileTemplates returns the project files that [InitProjectService]
// generates from embedded templates. The order is the deterministic
// order in which they are written and listed in
// [driving.InitProjectResponse.Created]. Implemented as a function
// to avoid the gochecknoglobals false-positive on immutable list
// constants.
func fileTemplates() []fileTemplate {
	return []fileTemplate{
		{Path: "README.md", TemplateName: "readme.md.tmpl"},
		{Path: "CHANGELOG.md", TemplateName: "changelog.md.tmpl"},
		{Path: "compose.yaml", TemplateName: "compose.yaml.tmpl"},
		{Path: ".env.example", TemplateName: "env.example.tmpl"},
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
// order. Used by the integrity self-test in `templates_test.go`.
func templateNames() ([]string, error) {
	entries, err := templateFS.ReadDir("templates")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	sort.Strings(names)
	return names, nil
}
