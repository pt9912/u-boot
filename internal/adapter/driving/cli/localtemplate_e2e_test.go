package cli_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/externaltemplates"
	"github.com/pt9912/u-boot/internal/adapter/driven/fs"
	"github.com/pt9912/u-boot/internal/adapter/driven/localtemplates"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// This test pins the end-to-end exit-code mapping for
// `u-boot init --template <ref>` across the real wiring chain that
// cmd/uboot assembles (slice-later-local-templates T3): the Composite
// classifies the ref, delegates to the embed.FS catalog or the
// filesystem resolver, TemplateInitService renders, and cli.ExitCode
// maps the resulting sentinel to the LH-FA-CLI-006 class. It exercises
// the production adapters (no fakes) against TempDir fixtures.

const e2eValidYAML = `apiVersion: github.com/pt9912/u-boot/template/v1
name: sample
description: "E2E sample."
version: 1.0.0
`

// newRealTemplateInitService wires the exact chain from cmd/uboot:
// Composite(catalog, local) → TemplateInitService over the real fs.
func newRealTemplateInitService() driving.TemplateInitUseCase {
	composite := localtemplates.NewComposite(externaltemplates.New(), localtemplates.New())
	return application.NewTemplateInitService(composite, fs.New(), nil)
}

func writeLocalTemplate(t *testing.T, yaml string, files map[string]string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "tpl")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if yaml != "" {
		if err := os.WriteFile(filepath.Join(dir, "template.yaml"), []byte(yaml), 0o644); err != nil {
			t.Fatalf("write template.yaml: %v", err)
		}
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	return dir
}

func initFromTemplate(t *testing.T, ref string) (driving.TemplateInitResponse, error) {
	t.Helper()
	name, err := domain.NewProjectName("my-app")
	if err != nil {
		t.Fatalf("NewProjectName: %v", err)
	}
	return newRealTemplateInitService().Init(context.Background(), driving.TemplateInitRequest{
		BaseDir:      t.TempDir(),
		ProjectName:  name,
		TemplateName: ref,
	})
}

func TestInitTemplate_E2E_LocalValid(t *testing.T) {
	t.Parallel()
	dir := writeLocalTemplate(t, e2eValidYAML, map[string]string{
		"u-boot.yaml.tmpl": "name: {{ .Name }}\n",
		".gitignore":       "*.log\n",
	})
	name, err := domain.NewProjectName("my-app")
	if err != nil {
		t.Fatal(err)
	}
	baseDir := t.TempDir()
	_, err = newRealTemplateInitService().Init(context.Background(), driving.TemplateInitRequest{
		BaseDir:      baseDir,
		ProjectName:  name,
		TemplateName: dir,
	})
	if err != nil {
		t.Fatalf("Init(local valid): %v (exit %d)", err, cli.ExitCode(err))
	}
	if cli.ExitCode(err) != 0 {
		t.Errorf("ExitCode = %d, want 0", cli.ExitCode(err))
	}
	// Rendered + copied, template.yaml not leaked.
	got, rerr := os.ReadFile(filepath.Join(baseDir, "u-boot.yaml"))
	if rerr != nil || string(got) != "name: my-app\n" {
		t.Errorf("u-boot.yaml = %q (err %v), want %q", got, rerr, "name: my-app\n")
	}
	if _, serr := os.Stat(filepath.Join(baseDir, "template.yaml")); !os.IsNotExist(serr) {
		t.Errorf("template.yaml must not be rendered into the project")
	}
}

func TestInitTemplate_E2E_CatalogNameStillResolves(t *testing.T) {
	t.Parallel()
	// A bare name routes through the Composite to the embed.FS catalog
	// (the shipped `basic` template) — proving the dispatch did not
	// break the existing catalog path.
	_, err := initFromTemplate(t, "basic")
	if err != nil {
		t.Fatalf("Init(basic): %v (exit %d)", err, cli.ExitCode(err))
	}
}

func TestInitTemplate_E2E_ExitCodeMatrix(t *testing.T) {
	t.Parallel()

	t.Run("missing path -> 10", func(t *testing.T) {
		t.Parallel()
		missing := filepath.Join(t.TempDir(), "nope")
		_, err := initFromTemplate(t, missing)
		if !errors.Is(err, driving.ErrTemplateNotFound) {
			t.Errorf("err = %v, want ErrTemplateNotFound", err)
		}
		if got := cli.ExitCode(err); got != 10 {
			t.Errorf("ExitCode = %d, want 10", got)
		}
	})

	t.Run("malformed metadata -> 10", func(t *testing.T) {
		t.Parallel()
		bad := "apiVersion: github.com/pt9912/u-boot/template/v2\nname: x\ndescription: y\nversion: 1\n"
		dir := writeLocalTemplate(t, bad, nil)
		_, err := initFromTemplate(t, dir)
		if !errors.Is(err, driving.ErrTemplateInvalid) {
			t.Errorf("err = %v, want ErrTemplateInvalid", err)
		}
		if got := cli.ExitCode(err); got != 10 {
			t.Errorf("ExitCode = %d, want 10", got)
		}
	})

	t.Run("symlink in tree -> 10", func(t *testing.T) {
		t.Parallel()
		dir := writeLocalTemplate(t, e2eValidYAML, map[string]string{"ok.tmpl": "{{ .Name }}\n"})
		if err := os.Symlink("/etc/passwd", filepath.Join(dir, "evil")); err != nil {
			t.Skipf("symlink unsupported: %v", err)
		}
		baseDir := t.TempDir()
		name, _ := domain.NewProjectName("my-app")
		_, err := newRealTemplateInitService().Init(context.Background(), driving.TemplateInitRequest{
			BaseDir: baseDir, ProjectName: name, TemplateName: dir,
		})
		if !errors.Is(err, driving.ErrInvalidTemplatePath) {
			t.Errorf("err = %v, want ErrInvalidTemplatePath", err)
		}
		if got := cli.ExitCode(err); got != 10 {
			t.Errorf("ExitCode = %d, want 10", got)
		}
		// No partial output: the symlink reject fires in render phase 1.
		if entries, _ := os.ReadDir(baseDir); len(entries) != 0 {
			t.Errorf("project dir not empty after symlink reject: %v", entries)
		}
	})

	t.Run("render failure -> 14", func(t *testing.T) {
		t.Parallel()
		dir := writeLocalTemplate(t, e2eValidYAML, map[string]string{"broken.tmpl": "{{ .Name"}) // unclosed action
		_, err := initFromTemplate(t, dir)
		if !errors.Is(err, driving.ErrTemplateRender) {
			t.Errorf("err = %v, want ErrTemplateRender", err)
		}
		if got := cli.ExitCode(err); got != 14 {
			t.Errorf("ExitCode = %d, want 14", got)
		}
	})
}
