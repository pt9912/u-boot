package localtemplates_test

import (
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/localtemplates"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

const validYAML = `apiVersion: github.com/pt9912/u-boot/template/v1
name: sample
description: "A local sample template."
version: 1.0.0
`

// writeTemplateDir creates dir with a template.yaml (content yaml) plus
// any extra files, and returns dir.
func writeTemplateDir(t *testing.T, dir, yaml string, extra map[string]string) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	if yaml != "" {
		if err := os.WriteFile(filepath.Join(dir, "template.yaml"), []byte(yaml), 0o644); err != nil {
			t.Fatalf("write template.yaml: %v", err)
		}
	}
	for name, content := range extra {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	return dir
}

func TestResolver_Open_Valid(t *testing.T) {
	t.Parallel()
	dir := writeTemplateDir(t, filepath.Join(t.TempDir(), "tpl"), validYAML, map[string]string{
		"u-boot.yaml.tmpl": "name: {{ .Name }}\n",
	})

	root, err := localtemplates.New().Open(context.Background(), dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	// The returned FS is rooted at the template dir and includes both
	// template.yaml and the renderable file.
	if _, err := iofs.ReadFile(root, "template.yaml"); err != nil {
		t.Errorf("ReadFile(template.yaml): %v", err)
	}
	if _, err := iofs.ReadFile(root, "u-boot.yaml.tmpl"); err != nil {
		t.Errorf("ReadFile(u-boot.yaml.tmpl): %v", err)
	}
}

func TestResolver_Open_Errors(t *testing.T) {
	t.Parallel()

	t.Run("empty path", func(t *testing.T) {
		t.Parallel()
		_, err := localtemplates.New().Open(context.Background(), "")
		if !errors.Is(err, driven.ErrTemplateNotFound) {
			t.Errorf("err = %v, want ErrTemplateNotFound", err)
		}
	})

	t.Run("nonexistent path", func(t *testing.T) {
		t.Parallel()
		missing := filepath.Join(t.TempDir(), "nope")
		_, err := localtemplates.New().Open(context.Background(), missing)
		if !errors.Is(err, driven.ErrTemplateNotFound) {
			t.Errorf("err = %v, want ErrTemplateNotFound", err)
		}
	})

	t.Run("not a directory", func(t *testing.T) {
		t.Parallel()
		file := filepath.Join(t.TempDir(), "afile")
		if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := localtemplates.New().Open(context.Background(), file)
		if !errors.Is(err, driven.ErrTemplateNotFound) {
			t.Errorf("err = %v, want ErrTemplateNotFound", err)
		}
	})

	t.Run("missing template.yaml", func(t *testing.T) {
		t.Parallel()
		dir := writeTemplateDir(t, filepath.Join(t.TempDir(), "tpl"), "", nil)
		_, err := localtemplates.New().Open(context.Background(), dir)
		if !errors.Is(err, driven.ErrTemplateNotFound) {
			t.Errorf("err = %v, want ErrTemplateNotFound (missing template.yaml)", err)
		}
	})

	t.Run("malformed template.yaml", func(t *testing.T) {
		t.Parallel()
		bad := "apiVersion: github.com/pt9912/u-boot/template/v2\nname: x\ndescription: y\nversion: 1\n"
		dir := writeTemplateDir(t, filepath.Join(t.TempDir(), "tpl"), bad, nil)
		_, err := localtemplates.New().Open(context.Background(), dir)
		if !errors.Is(err, driven.ErrTemplateInvalid) {
			t.Errorf("err = %v, want ErrTemplateInvalid (unsupported apiVersion)", err)
		}
		// Exit-class pin: invalid metadata is NOT a not-found.
		if errors.Is(err, driven.ErrTemplateNotFound) {
			t.Errorf("err = %v, must not classify as ErrTemplateNotFound", err)
		}
	})
}

// Home expansion uses os.UserHomeDir, which reads $HOME on unix. These
// subtests mutate the environment, so they cannot run in parallel.
func TestResolver_Open_HomeExpansion(t *testing.T) {
	t.Run("tilde-slash", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		writeTemplateDir(t, filepath.Join(home, "tpl"), validYAML, nil)

		root, err := localtemplates.New().Open(context.Background(), "~/tpl")
		if err != nil {
			t.Fatalf("Open(~/tpl): %v", err)
		}
		if _, err := iofs.ReadFile(root, "template.yaml"); err != nil {
			t.Errorf("ReadFile(template.yaml): %v", err)
		}
	})

	t.Run("bare tilde", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		writeTemplateDir(t, home, validYAML, nil) // home dir itself is the template

		if _, err := localtemplates.New().Open(context.Background(), "~"); err != nil {
			t.Fatalf("Open(~): %v", err)
		}
	})

	t.Run("home unresolvable", func(t *testing.T) {
		t.Setenv("HOME", "")
		_, err := localtemplates.New().Open(context.Background(), "~/tpl")
		if err == nil {
			t.Fatal("Open(~/tpl) with empty HOME: want error, got nil")
		}
		// Environment failure is technical, not a template sentinel.
		if errors.Is(err, driven.ErrTemplateNotFound) || errors.Is(err, driven.ErrTemplateInvalid) {
			t.Errorf("err = %v, want non-sentinel technical error", err)
		}
	})
}

func TestResolver_Open_CtxCancelled(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := localtemplates.New().Open(ctx, "/whatever"); !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

// recordingFiles is a fake TemplateFiles that records the ref it was
// asked for and returns a tagged marker error so the Composite test can
// assert which resolver handled the dispatch.
type recordingFiles struct {
	tag      string
	called   bool
	lastName string
}

func (r *recordingFiles) Open(_ context.Context, name string) (iofs.FS, error) {
	r.called = true
	r.lastName = name
	return nil, fmt.Errorf("%s:%s", r.tag, name)
}

func TestComposite_DispatchesByKind(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		ref       string
		wantLocal bool
	}{
		{"catalog name", "basic", false},
		{"kebab name", "micronaut-sveltekit", false},
		{"bare dot is path", ".", true},
		{"bare dotdot is path", "..", true},
		{"dot-slash path", "./my-template", true},
		{"absolute path", "/abs/tpl", true},
		{"home path", "~/tpl", true},
		{"nested slash", "vendor/tpl", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cat := &recordingFiles{tag: "cat"}
			loc := &recordingFiles{tag: "loc"}
			_, _ = localtemplates.NewComposite(cat, loc).Open(context.Background(), tc.ref)

			if tc.wantLocal {
				if !loc.called || cat.called {
					t.Errorf("ref %q: local.called=%v cat.called=%v, want local only", tc.ref, loc.called, cat.called)
				}
				if loc.lastName != tc.ref {
					t.Errorf("local got %q, want raw ref %q", loc.lastName, tc.ref)
				}
			} else {
				if !cat.called || loc.called {
					t.Errorf("ref %q: cat.called=%v loc.called=%v, want catalog only", tc.ref, cat.called, loc.called)
				}
				if cat.lastName != tc.ref {
					t.Errorf("catalog got %q, want raw ref %q", cat.lastName, tc.ref)
				}
			}
		})
	}
}
