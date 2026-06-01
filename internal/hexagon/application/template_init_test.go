package application_test

import (
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"testing"
	"testing/fstest"

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
