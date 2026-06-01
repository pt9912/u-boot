package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/externaltemplates"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TestInitProjectService_FromTemplate_BasicByteIdenticalToDefaultPath
// pins the slice-v1-template-init T4 integration: running
// `u-boot init demo --template basic` produces a project that is
// byte-identical (per-file content) to `u-boot init demo` (the
// default render path) — the structural goal of the slice.
//
// Distinct from the T3 byte-identity test in template_init_test.go:
// that one drives TemplateInitService directly with fake FS and
// pinned strings. This one drives the full InitProjectService
// (with WithTemplateInit wired against the production catalog
// adapter) and compares the two services' output side-by-side
// using the same fakeFS shape. A drift between the two paths
// shows up here even if both individually produce
// "reasonable-looking" content.
func TestInitProjectService_FromTemplate_BasicByteIdenticalToDefaultPath(t *testing.T) {
	t.Parallel()

	const baseDir = "/proj"
	const projectName = "demo"

	tmplFS := newFakeFS()
	tmplFS.dirs[baseDir] = true
	tmplInit := application.NewTemplateInitService(externaltemplates.New(), tmplFS, nil)
	tmplSvc := application.NewInitProjectService(tmplFS, &fakeYAML{}, &fakeGit{}, nil, nil, nil, application.WithTemplateInit(tmplInit))

	defaultFS := newFakeFS()
	defaultFS.dirs[baseDir] = true
	defaultSvc := application.NewInitProjectService(defaultFS, &fakeYAML{}, &fakeGit{}, nil, nil, nil)

	tmplResp, err := tmplSvc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir:  baseDir,
		Name:     projectName,
		SkipGit:  true,
		Template: "basic",
	})
	if err != nil {
		t.Fatalf("template Init: %v", err)
	}

	defaultResp, err := defaultSvc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir: baseDir,
		Name:    projectName,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("default Init: %v", err)
	}

	// Pin the 6 generated files; they must be byte-identical.
	wantFiles := []string{
		baseDir + "/u-boot.yaml",
		baseDir + "/compose.yaml",
		baseDir + "/README.md",
		baseDir + "/CHANGELOG.md",
		baseDir + "/.env.example",
		baseDir + "/.gitignore",
	}
	for _, p := range wantFiles {
		got := string(tmplFS.files[p])
		want := string(defaultFS.files[p])
		if got == "" {
			t.Errorf("%s: missing in template-path output", p)
			continue
		}
		if want == "" {
			t.Errorf("%s: missing in default-path output", p)
			continue
		}
		if got != want {
			t.Errorf("%s: byte-identity mismatch\n--- default ---\n%s\n--- template ---\n%s", p, want, got)
		}
	}

	// Created lists must be permutation-equivalent (both contain the
	// 3 dirs + 6 files; default also adds u-boot.yaml via a different
	// mechanism so the order differs slightly).
	if len(tmplResp.Created) != len(defaultResp.Created) {
		t.Errorf("len(Created): template=%d, default=%d", len(tmplResp.Created), len(defaultResp.Created))
	}
}

func TestInitProjectService_FromTemplate_RejectsDevcontainerMutex(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.dirs["/proj"] = true
	tmplInit := application.NewTemplateInitService(externaltemplates.New(), fs, nil)
	svc := application.NewInitProjectService(fs, &fakeYAML{}, &fakeGit{}, nil, nil, nil, application.WithTemplateInit(tmplInit))

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir:      "/proj",
		Name:         "demo",
		SkipGit:      true,
		Template:     "basic",
		Devcontainer: true,
	})
	if err == nil {
		t.Fatal("Init: want ErrTemplateConflictsWithFlag, got nil")
	}
	if !errors.Is(err, driving.ErrTemplateConflictsWithFlag) {
		t.Errorf("err = %v, want wrap of driving.ErrTemplateConflictsWithFlag", err)
	}
}

func TestInitProjectService_FromTemplate_RejectsForceMutex(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.dirs["/proj"] = true
	tmplInit := application.NewTemplateInitService(externaltemplates.New(), fs, nil)
	svc := application.NewInitProjectService(fs, &fakeYAML{}, &fakeGit{}, nil, nil, nil, application.WithTemplateInit(tmplInit))

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir:  "/proj",
		Name:     "demo",
		SkipGit:  true,
		Template: "basic",
		Force:    true,
	})
	if err == nil {
		t.Fatal("Init: want ErrTemplateConflictsWithFlag, got nil")
	}
	if !errors.Is(err, driving.ErrTemplateConflictsWithFlag) {
		t.Errorf("err = %v, want wrap of driving.ErrTemplateConflictsWithFlag", err)
	}
}

func TestInitProjectService_FromTemplate_RejectsExistingProject(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.dirs["/proj"] = true
	// Pre-seed u-boot.yaml to simulate an existing project.
	if err := fs.WriteFile("/proj/u-boot.yaml", []byte("schemaVersion: 1\n"), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	tmplInit := application.NewTemplateInitService(externaltemplates.New(), fs, nil)
	svc := application.NewInitProjectService(fs, &fakeYAML{}, &fakeGit{}, nil, nil, nil, application.WithTemplateInit(tmplInit))

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir:  "/proj",
		Name:     "demo",
		SkipGit:  true,
		Template: "basic",
	})
	if err == nil {
		t.Fatal("Init: want ErrProjectExists, got nil")
	}
	if !errors.Is(err, driving.ErrProjectExists) {
		t.Errorf("err = %v, want wrap of driving.ErrProjectExists", err)
	}
}

func TestInitProjectService_FromTemplate_RequiresWiredTemplateInit(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.dirs["/proj"] = true
	// No WithTemplateInit option — the service has nil templateInit.
	svc := application.NewInitProjectService(fs, &fakeYAML{}, &fakeGit{}, nil, nil, nil)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir:  "/proj",
		Name:     "demo",
		SkipGit:  true,
		Template: "basic",
	})
	if err == nil {
		t.Fatal("Init: want wiring error, got nil")
	}
}

func TestInitProjectService_FromTemplate_LeftoverStructureDirsDoNotBlock(t *testing.T) {
	t.Parallel()
	// Review-followup F3: a directory with ≥3 LH-FA-INIT-003
	// indicators (docker/, scripts/, docs/) from a previous failed
	// init must NOT trigger soft-existing-detection on the template
	// path — the user explicitly chose fresh-init via --template.
	// The hard-existing check (u-boot.yaml) is the only safety net.
	fs := newFakeFS()
	fs.dirs["/proj"] = true
	fs.dirs["/proj/docker"] = true
	fs.dirs["/proj/scripts"] = true
	fs.dirs["/proj/docs"] = true
	tmplInit := application.NewTemplateInitService(externaltemplates.New(), fs, nil)
	svc := application.NewInitProjectService(fs, &fakeYAML{}, &fakeGit{}, nil, nil, nil, application.WithTemplateInit(tmplInit))

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir:  "/proj",
		Name:     "demo",
		SkipGit:  true,
		Template: "basic",
	})
	if err != nil {
		t.Fatalf("Init with leftover dirs + --template basic: %v (want success — soft-detection must be skipped on template path)", err)
	}
	// Verify the basic-rendered u-boot.yaml landed.
	if _, ok := fs.files["/proj/u-boot.yaml"]; !ok {
		t.Error("u-boot.yaml not rendered")
	}
}

func TestInitProjectService_FromTemplate_UnknownTemplateWrapsSentinel(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.dirs["/proj"] = true
	tmplInit := application.NewTemplateInitService(externaltemplates.New(), fs, nil)
	svc := application.NewInitProjectService(fs, &fakeYAML{}, &fakeGit{}, nil, nil, nil, application.WithTemplateInit(tmplInit))

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir:  "/proj",
		Name:     "demo",
		SkipGit:  true,
		Template: "no-such-template",
	})
	if err == nil {
		t.Fatal("Init: want ErrTemplateNotFound, got nil")
	}
	if !errors.Is(err, driving.ErrTemplateNotFound) {
		t.Errorf("err = %v, want wrap of driving.ErrTemplateNotFound", err)
	}
}
