package application_test

import (
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"testing"
	"testing/fstest"

	"github.com/pt9912/u-boot/internal/adapter/driven/externaltemplates"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// fakeTemplateFiles implements driven.TemplateFiles for the
// TemplateInitService unit tests. Returns the configured FS for
// the configured name and ErrTemplateNotFound otherwise; openErr
// short-circuits Open if non-nil.
type fakeTemplateFiles struct {
	name    string
	fs      iofs.FS
	openErr error
}

func (f *fakeTemplateFiles) Open(_ context.Context, name string) (iofs.FS, error) {
	if f.openErr != nil {
		return nil, f.openErr
	}
	if name != f.name {
		return nil, fmt.Errorf("%w: %q", driven.ErrTemplateNotFound, name)
	}
	return f.fs, nil
}

func mustProjectNameForTest(t *testing.T, raw string) domain.ProjectName {
	t.Helper()
	name, err := domain.NewProjectName(raw)
	if err != nil {
		t.Fatalf("NewProjectName(%q): %v", raw, err)
	}
	return name
}

func TestTemplateInitService_RendersTmplAndCopiesPlain(t *testing.T) {
	t.Parallel()
	files := &fakeTemplateFiles{
		name: "basic",
		fs: fstest.MapFS{
			"u-boot.yaml.tmpl": &fstest.MapFile{Data: []byte("name: {{ .Name }}\n")},
			".gitignore":       &fstest.MapFile{Data: []byte("*.log\n")},
			"template.yaml":    &fstest.MapFile{Data: []byte("ignored: metadata\n")},
		},
	}
	fs := newFakeFS()
	svc := application.NewTemplateInitService(files, fs, nil)

	resp, err := svc.Init(context.Background(), driving.TemplateInitRequest{
		BaseDir:      "/proj",
		ProjectName:  mustProjectNameForTest(t, "my-app"),
		TemplateName: "basic",
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	// Sorted Created list.
	wantCreated := []string{".gitignore", "u-boot.yaml"}
	if len(resp.Created) != len(wantCreated) {
		t.Fatalf("Created = %v, want %v", resp.Created, wantCreated)
	}
	for i, p := range wantCreated {
		if resp.Created[i] != p {
			t.Errorf("Created[%d] = %q, want %q", i, resp.Created[i], p)
		}
	}

	// Rendered .tmpl substitutes the project name.
	if got := string(fs.files["/proj/u-boot.yaml"]); got != "name: my-app\n" {
		t.Errorf("u-boot.yaml content = %q, want %q", got, "name: my-app\n")
	}
	// Plain file copied byte-identically.
	if got := string(fs.files["/proj/.gitignore"]); got != "*.log\n" {
		t.Errorf(".gitignore content = %q, want %q", got, "*.log\n")
	}
	// template.yaml was skipped.
	if _, ok := fs.files["/proj/template.yaml"]; ok {
		t.Errorf("template.yaml should be skipped, but was written")
	}
}

func TestTemplateInitService_NestedDirectoriesAreCreated(t *testing.T) {
	t.Parallel()
	files := &fakeTemplateFiles{
		name: "basic",
		fs: fstest.MapFS{
			"docker/Dockerfile.tmpl": &fstest.MapFile{Data: []byte("FROM alpine\n# {{ .Name }}\n")},
		},
	}
	fs := newFakeFS()
	svc := application.NewTemplateInitService(files, fs, nil)

	resp, err := svc.Init(context.Background(), driving.TemplateInitRequest{
		BaseDir:      "/proj",
		ProjectName:  mustProjectNameForTest(t, "x"),
		TemplateName: "basic",
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if len(resp.Created) != 1 || resp.Created[0] != "docker/Dockerfile" {
		t.Errorf("Created = %v, want [docker/Dockerfile]", resp.Created)
	}
	if got := string(fs.files["/proj/docker/Dockerfile"]); got != "FROM alpine\n# x\n" {
		t.Errorf("Dockerfile content = %q", got)
	}
	// Parent dir auto-created.
	if !fs.dirs["/proj/docker"] {
		t.Errorf("/proj/docker dir not registered")
	}
}

func TestTemplateInitService_UnknownTemplateWrapsSentinel(t *testing.T) {
	t.Parallel()
	files := &fakeTemplateFiles{name: "basic", fs: fstest.MapFS{}}
	fs := newFakeFS()
	svc := application.NewTemplateInitService(files, fs, nil)

	_, err := svc.Init(context.Background(), driving.TemplateInitRequest{
		BaseDir:      "/proj",
		ProjectName:  mustProjectNameForTest(t, "x"),
		TemplateName: "nonexistent",
	})
	if err == nil {
		t.Fatal("Init: want error, got nil")
	}
	if !errors.Is(err, driving.ErrTemplateNotFound) {
		t.Errorf("err = %v, want wrap of driving.ErrTemplateNotFound", err)
	}
	if !errors.Is(err, driven.ErrTemplateNotFound) {
		t.Errorf("err = %v, want preserve of driven.ErrTemplateNotFound via multi-%%w", err)
	}
}

func TestTemplateInitService_RenderFailureWrapsErrTemplateRender(t *testing.T) {
	t.Parallel()
	files := &fakeTemplateFiles{
		name: "basic",
		fs: fstest.MapFS{
			// Malformed text/template: unclosed action.
			"broken.txt.tmpl": &fstest.MapFile{Data: []byte("{{ .Name")},
		},
	}
	fs := newFakeFS()
	svc := application.NewTemplateInitService(files, fs, nil)

	_, err := svc.Init(context.Background(), driving.TemplateInitRequest{
		BaseDir:      "/proj",
		ProjectName:  mustProjectNameForTest(t, "x"),
		TemplateName: "basic",
	})
	if err == nil {
		t.Fatal("Init: want render error, got nil")
	}
	if !errors.Is(err, driving.ErrTemplateRender) {
		t.Errorf("err = %v, want wrap of driving.ErrTemplateRender", err)
	}
	// Nothing was written when the render failed.
	if len(fs.writes) > 0 {
		t.Errorf("partial writes on render failure: %v", fs.writes)
	}
}

// The path-escape boundary (`..` segments, absolute paths, drive
// letters) is exhaustively pinned in domain/template_path_test.go
// — domain.NewTemplatePath has its own 14-case table. The service
// wrap is a one-line `fmt.Errorf("%w: %w", driving.
// ErrInvalidTemplatePath, err)` transform whose correctness is
// trivially visible (template_init.go renderOne). Reaching it via
// fstest.MapFS would require building a custom iofs.FS that
// bypasses fs.ValidPath's `..` rejection — a significant test-only
// dependency for one line of production code.

func TestTemplateInitService_EmptyBaseDirRejected(t *testing.T) {
	t.Parallel()
	svc := application.NewTemplateInitService(&fakeTemplateFiles{name: "basic", fs: fstest.MapFS{}}, newFakeFS(), nil)
	_, err := svc.Init(context.Background(), driving.TemplateInitRequest{
		BaseDir:      "",
		ProjectName:  mustProjectNameForTest(t, "x"),
		TemplateName: "basic",
	})
	if err == nil {
		t.Fatal("Init: want error, got nil")
	}
}

func TestTemplateInitService_NilLoggerAccepted(t *testing.T) {
	t.Parallel()
	svc := application.NewTemplateInitService(&fakeTemplateFiles{name: "basic", fs: fstest.MapFS{}}, newFakeFS(), nil)
	if svc == nil {
		t.Fatal("New returned nil")
	}
}

// TestTemplateInitService_BasicByteIdenticalToDefaultInit pins the
// slice-v1-template-init T3 byte-identity guarantee: every file the
// `basic` external template renders for `{Name: "demo"}` must equal
// the file the production InitProjectService writes today.
//
// The expected strings were captured from `docker run --rm u-boot
// init demo --no-git` against u-boot:latest at slice-v1-template-init
// T3 commit-time. A diff between the two paths is the regression
// signal that fires either when:
//
//   - the basic/ source templates drift from `internal/hexagon/
//     application/templates/`,
//   - the yaml.v3 default indent / encoding semantics change
//     (currently 4-space indent in `u-boot.yaml`),
//   - or `text/template` evaluation diverges between the two render
//     paths (it shouldn't — both use the same engine).
//
// This is the test that future micronaut/sveltekit slices will use
// as their conformance scaffolding: rendering must match a known
// good string.
func TestTemplateInitService_BasicByteIdenticalToDefaultInit(t *testing.T) {
	t.Parallel()
	files := externaltemplates.New()
	fs := newFakeFS()
	svc := application.NewTemplateInitService(files, fs, nil)

	_, err := svc.Init(context.Background(), driving.TemplateInitRequest{
		BaseDir:      "/proj",
		ProjectName:  mustProjectNameForTest(t, "demo"),
		TemplateName: "basic",
	})
	if err != nil {
		t.Fatalf("Init(basic): %v", err)
	}

	want := map[string]string{
		"/proj/u-boot.yaml": "schemaVersion: 1\n" +
			"project:\n" +
			"    name: demo\n",
		"/proj/compose.yaml": "# BEGIN U-BOOT MANAGED BLOCK: init\n" +
			"# Compose stack for demo. Service add-ons are appended by\n" +
			"# `u-boot add <service>` (LH-FA-ADD-*); manual edits outside the\n" +
			"# u-boot-managed blocks remain untouched (LH-SA-FILE-002).\n" +
			"\n" +
			"name: demo\n" +
			"\n" +
			"networks:\n" +
			"  default:\n" +
			"    name: demo-default\n" +
			"# END U-BOOT MANAGED BLOCK: init\n" +
			"\n" +
			"services: {}\n" +
			"\n" +
			"volumes: {}\n",
		"/proj/README.md": "<!-- BEGIN U-BOOT MANAGED BLOCK: init -->\n" +
			"# demo\n" +
			"\n" +
			"A Docker-based development environment scaffolded by [u-boot](https://github.com/pt9912/u-boot).\n" +
			"\n" +
			"## Quickstart\n" +
			"\n" +
			"```bash\n" +
			"u-boot up         # start the development stack\n" +
			"u-boot doctor     # check local prerequisites\n" +
			"u-boot down       # stop the stack\n" +
			"```\n" +
			"\n" +
			"## Project layout\n" +
			"\n" +
			"- `compose.yaml` — Docker Compose stack.\n" +
			"- `docker/` — Dockerfile(s) and container-related assets.\n" +
			"- `scripts/` — project-local helper scripts.\n" +
			"- `docs/` — project documentation.\n" +
			"- `u-boot.yaml` — u-boot project configuration.\n" +
			"- `.env.example` — copy to `.env` and adapt before `u-boot up`.\n" +
			"\n" +
			"## Prerequisites\n" +
			"\n" +
			"- Docker Engine ≥ 24.0.0\n" +
			"- Docker Compose ≥ 2.20.0\n" +
			"- Git\n" +
			"<!-- END U-BOOT MANAGED BLOCK: init -->\n",
		"/proj/CHANGELOG.md": "<!-- BEGIN U-BOOT MANAGED BLOCK: init -->\n" +
			"# Changelog\n" +
			"\n" +
			"All notable changes to **demo** are documented in this file. The\n" +
			"format is loosely based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);\n" +
			"versions follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).\n" +
			"\n" +
			"## [Unreleased]\n" +
			"\n" +
			"### Added\n" +
			"\n" +
			"- Project scaffolded with `u-boot init`.\n" +
			"<!-- END U-BOOT MANAGED BLOCK: init -->\n",
		"/proj/.env.example": "# BEGIN U-BOOT MANAGED BLOCK: init\n" +
			"# Example environment for demo.\n" +
			"#\n" +
			"# Copy this file to `.env` and adapt the values; `.env` is git-ignored\n" +
			"# and read by `docker compose` automatically.\n" +
			"#\n" +
			"# Service add-ons (`u-boot add postgres` etc.) append their variables\n" +
			"# inside dedicated `U-BOOT MANAGED BLOCK` regions (LH-SA-FILE-002).\n" +
			"# END U-BOOT MANAGED BLOCK: init\n",
		"/proj/.gitignore": "# Generated by `u-boot init` for demo.\n" +
			"#\n" +
			"# Note: this file is treated as fully u-boot-managed (LH-FA-INIT-005).\n" +
			"# Re-running `u-boot init` on an existing project requires `--backup`\n" +
			"# to overwrite this file safely — your customizations are preserved\n" +
			"# as `.gitignore.bak[*]` (LH-SA-FILE-002 §611 list does NOT cover\n" +
			"# .gitignore, so the managed-block escape hatch does not apply).\n" +
			"\n" +
			"# Local environment overrides; never commit.\n" +
			".env\n" +
			".env.local\n" +
			"\n" +
			"# Editor / IDE\n" +
			".idea/\n" +
			".vscode/\n" +
			"*.swp\n" +
			"*.swo\n" +
			"\n" +
			"# OS junk\n" +
			".DS_Store\n" +
			"Thumbs.db\n" +
			"\n" +
			"# Backup files written by u-boot's overwrite-protection\n" +
			"# (LH-FA-INIT-005).\n" +
			"*.bak\n" +
			"*.bak.*\n",
	}

	for path, expected := range want {
		got, ok := fs.files[path]
		if !ok {
			t.Errorf("%s: missing from render output", path)
			continue
		}
		if string(got) != expected {
			t.Errorf("%s: byte-identity mismatch.\n--- expected ---\n%s\n--- got ---\n%s", path, expected, string(got))
		}
	}
}
