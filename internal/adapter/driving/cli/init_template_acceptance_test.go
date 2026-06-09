package cli_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/externaltemplates"
	"github.com/pt9912/u-boot/internal/adapter/driven/fs"
	"github.com/pt9912/u-boot/internal/adapter/driven/git"
	"github.com/pt9912/u-boot/internal/adapter/driven/localtemplates"
	"github.com/pt9912/u-boot/internal/adapter/driven/yaml"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/application"
)

// realInitApp wires the full `u-boot init` command path exactly as
// cmd/uboot does — real fs/yaml/git adapters and a real
// InitProjectService whose template delegate is the production
// Composite(catalog, local) → TemplateInitService chain. getwd is
// pinned to baseDir so the command writes into a TempDir without
// chdir (slice-later-local-templates T4). git init is exercised only
// when the caller omits --no-git; the tests pass --no-git to stay
// hermetic.
func realInitApp(baseDir string) *cli.App {
	composite := localtemplates.NewComposite(externaltemplates.New(), localtemplates.New())
	tmplInit := application.NewTemplateInitService(composite, fs.New(), nil)
	initSvc := application.NewInitProjectService(fs.New(), yaml.New(), git.New(), nil, nil, nil, application.WithTemplateInit(tmplInit))
	return newApp(initSvc, cli.WithGetwd(func() (string, error) { return baseDir, nil }))
}

// TestInitCLI_LocalTemplate_RendersProject drives the real `init`
// command with a filesystem `--template ./path` and asserts the
// project is rendered: the template's `.tmpl` is rendered with the
// project name, the LH-FA-INIT-003 structure dirs are created, and
// template.yaml is not leaked into the project.
func TestInitCLI_LocalTemplate_RendersProject(t *testing.T) {
	t.Parallel()
	tplDir := writeLocalTemplate(t, e2eValidYAML, map[string]string{
		"u-boot.yaml.tmpl": "name: {{ .Name }}\n",
		".gitignore":       "*.log\n",
	})
	baseDir := t.TempDir()

	var stdout, stderr bytes.Buffer
	err := realInitApp(baseDir).Execute(
		context.Background(),
		[]string{"init", "myproj", "--template", tplDir, "--no-git"},
		&stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("init --template: %v (stderr=%s)", err, stderr.String())
	}

	if got, rerr := os.ReadFile(filepath.Join(baseDir, "u-boot.yaml")); rerr != nil || string(got) != "name: myproj\n" {
		t.Errorf("u-boot.yaml = %q (err %v), want %q", got, rerr, "name: myproj\n")
	}
	if got, rerr := os.ReadFile(filepath.Join(baseDir, ".gitignore")); rerr != nil || string(got) != "*.log\n" {
		t.Errorf(".gitignore = %q (err %v), want %q", got, rerr, "*.log\n")
	}
	if _, serr := os.Stat(filepath.Join(baseDir, "template.yaml")); !os.IsNotExist(serr) {
		t.Errorf("template.yaml must not be rendered into the project")
	}
	// LH-FA-INIT-003 structure dirs are an init-flow concern, created
	// alongside the template render.
	for _, dir := range []string{"docker", "scripts", "docs"} {
		if info, serr := os.Stat(filepath.Join(baseDir, dir)); serr != nil || !info.IsDir() {
			t.Errorf("structure dir %q missing (err %v)", dir, serr)
		}
	}
}

// TestInitCLI_LocalTemplate_MissingPathExits10 drives the command with
// a non-existent local path and asserts the CLI surfaces the
// validation error (exit-code 10 mapping pinned in the T3 matrix); the
// command must not partially populate the project dir.
func TestInitCLI_LocalTemplate_MissingPath(t *testing.T) {
	t.Parallel()
	baseDir := t.TempDir()
	missing := filepath.Join(t.TempDir(), "nope")

	var stdout, stderr bytes.Buffer
	err := realInitApp(baseDir).Execute(
		context.Background(),
		[]string{"init", "myproj", "--template", missing, "--no-git"},
		&stdout, &stderr,
	)
	if err == nil {
		t.Fatal("init --template <missing>: want error, got nil")
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10 (err=%v)", got, err)
	}
	if _, serr := os.Stat(filepath.Join(baseDir, "u-boot.yaml")); !os.IsNotExist(serr) {
		t.Errorf("u-boot.yaml must not be written when the template path is missing")
	}
}
