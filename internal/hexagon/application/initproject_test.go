package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const testBaseDir = "/tmp/u-boot-test/demo"

func newService(t *testing.T) (*application.InitProjectService, *fakeFS, *fakeYAML, *fakeGit) {
	t.Helper()
	fs := newFakeFS()
	// The service refuses to initialize a non-existent BaseDir
	// (LH-AK-001 has the user mkdir first); the fake treats a
	// pre-registered directory as "exists".
	fs.markDirExists(testBaseDir)
	y := &fakeYAML{}
	g := &fakeGit{}
	return application.NewInitProjectService(fs, y, g), fs, y, g
}

func TestInit_HappyPath_CreatesStructureAndConfig(t *testing.T) {
	svc, fs, _, git := newService(t)

	resp, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	// Project aggregate.
	if resp.Project.Name.String() != "demo" {
		t.Errorf("Project.Name = %q, want %q", resp.Project.Name, "demo")
	}
	if resp.Project.SchemaVersion != domain.SchemaVersionCurrent {
		t.Errorf("Project.SchemaVersion = %d, want %d", resp.Project.SchemaVersion, domain.SchemaVersionCurrent)
	}

	// Directories per LH-FA-INIT-003, in the deterministic call order
	// asserted by projectStructureDirs().
	wantDirs := []string{
		filepath.Join(testBaseDir, "docker"),
		filepath.Join(testBaseDir, "scripts"),
		filepath.Join(testBaseDir, "docs"),
	}
	if got := fs.mkdirPaths(); !reflect.DeepEqual(got, wantDirs) {
		t.Errorf("MkdirAll paths = %v, want %v", got, wantDirs)
	}

	// Files per LH-FA-INIT-003 + LH-FA-CONF-002.
	wantFiles := []string{
		filepath.Join(testBaseDir, "README.md"),
		filepath.Join(testBaseDir, "CHANGELOG.md"),
		filepath.Join(testBaseDir, "compose.yaml"),
		filepath.Join(testBaseDir, ".env.example"),
		filepath.Join(testBaseDir, ".gitignore"),
		filepath.Join(testBaseDir, "u-boot.yaml"),
	}
	if got := fs.writtenPaths(); !reflect.DeepEqual(got, wantFiles) {
		t.Errorf("WriteFile paths =\n  %v\nwant:\n  %v", got, wantFiles)
	}

	// Created list mirrors the write order.
	wantCreated := []string{
		"docker/", "scripts/", "docs/",
		"README.md", "CHANGELOG.md", "compose.yaml", ".env.example", ".gitignore",
		"u-boot.yaml",
	}
	if !reflect.DeepEqual(resp.Created, wantCreated) {
		t.Errorf("Created =\n  %v\nwant:\n  %v", resp.Created, wantCreated)
	}

	// Git default-on (LH-FA-INIT-007).
	if len(git.isRepoCalls) != 1 || git.isRepoCalls[0] != testBaseDir {
		t.Errorf("git.IsRepository calls = %v, want exactly 1 for %q", git.isRepoCalls, testBaseDir)
	}
	if len(git.initCalls) != 1 {
		t.Errorf("git.Init calls = %v, want exactly 1", git.initCalls)
	}
}

func TestInit_NameDerivedFromBaseDir(t *testing.T) {
	svc, fs, _, _ := newService(t)
	customBase := "/tmp/u-boot-test/My_Project Name"
	fs.markDirExists(customBase)

	resp, err := svc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir: customBase,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	// "My_Project Name" → "my-project-name" via LH-FA-INIT-002.
	if got := resp.Project.Name.String(); got != "my-project-name" {
		t.Errorf("Project.Name = %q, want %q", got, "my-project-name")
	}
}

func TestInit_MissingBaseDirRejects(t *testing.T) {
	// Why: pins the LH-AK-001-driven contract — the user must mkdir
	// before `u-boot init`. A typoed BaseDir must not silently
	// initialize a fresh tree under the typo.
	svc, _, _, _ := newService(t)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: "/tmp/u-boot-test/this-path-does-not-exist",
	})
	if err == nil {
		t.Fatalf("Init(missing BaseDir): expected error, got nil")
	}
	if !errors.Is(err, application.ErrBaseDirMissing) {
		t.Errorf("Init(missing BaseDir): error %v does not wrap ErrBaseDirMissing", err)
	}
}

func TestInit_InvalidNameRejects(t *testing.T) {
	svc, _, _, _ := newService(t)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "Invalid Name!",
		BaseDir: testBaseDir,
	})
	if err == nil {
		t.Fatalf("Init(invalid name): expected error, got nil")
	}
	if !errors.Is(err, domain.ErrInvalidProjectName) {
		t.Errorf("Init(invalid name): error %v does not wrap ErrInvalidProjectName", err)
	}
}

func TestInit_EmptyBaseDirRejects(t *testing.T) {
	svc, _, _, _ := newService(t)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name: "demo",
	})
	if err == nil {
		t.Fatalf("Init(empty BaseDir): expected error, got nil")
	}
}

func TestInit_ExistingProjectRejected(t *testing.T) {
	for _, marker := range []string{"u-boot.yaml", "compose.yaml", ".env.example"} {
		t.Run(marker, func(t *testing.T) {
			svc, fs, _, _ := newService(t)
			if err := fs.WriteFile(filepath.Join(testBaseDir, marker), []byte("preexisting"), 0o644); err != nil {
				t.Fatalf("setup: %v", err)
			}

			_, err := svc.Init(context.Background(), driving.InitProjectRequest{
				Name:    "demo",
				BaseDir: testBaseDir,
			})
			if err == nil {
				t.Fatalf("Init: expected error, got nil")
			}
			if !errors.Is(err, driving.ErrProjectExists) {
				t.Errorf("Init: error %v does not wrap ErrProjectExists", err)
			}
		})
	}
}

func TestInit_SkipGitDoesNotCallAdapter(t *testing.T) {
	svc, _, _, git := newService(t)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	if len(git.isRepoCalls) != 0 || len(git.initCalls) != 0 {
		t.Errorf("SkipGit=true but git was called: isRepoCalls=%v initCalls=%v", git.isRepoCalls, git.initCalls)
	}
}

func TestInit_GitAlreadyRepository_NoInitCall(t *testing.T) {
	// Why: LH-FA-INIT-007 forbids re-initializing an existing repo.
	svc, _, _, git := newService(t)
	git.isRepo = true

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	if len(git.initCalls) != 0 {
		t.Errorf("git.Init called %v despite existing repo", git.initCalls)
	}
}

func TestInit_UBootYAMLContainsSchemaAndName(t *testing.T) {
	svc, fs, _, _ := newService(t)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	body, err := fs.ReadFile(filepath.Join(testBaseDir, "u-boot.yaml"))
	if err != nil {
		t.Fatalf("ReadFile u-boot.yaml: %v", err)
	}
	got := string(body)
	for _, want := range []string{"schemaVersion: 1", "name: demo"} {
		if !strings.Contains(got, want) {
			t.Errorf("u-boot.yaml missing %q; got:\n%s", want, got)
		}
	}
}

func TestInit_RenderedTemplatesContainName(t *testing.T) {
	svc, fs, _, _ := newService(t)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "my-service",
		BaseDir: testBaseDir,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	for _, file := range []string{"README.md", "CHANGELOG.md", "compose.yaml", ".env.example", ".gitignore"} {
		body, err := fs.ReadFile(filepath.Join(testBaseDir, file))
		if err != nil {
			t.Fatalf("ReadFile %s: %v", file, err)
		}
		if !strings.Contains(string(body), "my-service") {
			t.Errorf("%s does not contain project name", file)
		}
	}
}

func TestInit_YAMLMarshalErrorPropagates(t *testing.T) {
	svc, _, yamlCodec, _ := newService(t)
	yamlCodec.failMarshal = true

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
	})
	if err == nil {
		t.Fatalf("Init: expected error from yaml.Marshal, got nil")
	}
	if !strings.Contains(err.Error(), "marshal") {
		t.Errorf("Init error %v does not mention marshal", err)
	}
}

func TestInit_GitInitErrorPropagates(t *testing.T) {
	svc, _, _, git := newService(t)
	git.initErr = errors.New("git failed")

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
	})
	if err == nil {
		t.Fatalf("Init: expected git init error, got nil")
	}
}

func TestInit_FSWriteFailurePropagates(t *testing.T) {
	// Why: Pins the error-propagation path through writeTemplatedFiles
	// — fakeFS's failOn lets us force a failure on a specific path.
	svc, fs, _, _ := newService(t)
	fs.failOn = filepath.Join(testBaseDir, "compose.yaml")
	fs.failErr = errors.New("disk full")

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
	})
	if err == nil {
		t.Fatalf("Init: expected error, got nil")
	}
	if !strings.Contains(err.Error(), "compose.yaml") {
		t.Errorf("Init error %v does not mention failing file", err)
	}
}

func TestInit_GitIsRepositoryErrorPropagates(t *testing.T) {
	svc, _, _, git := newService(t)
	git.isRepoErr = errors.New("git not installed")

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
	})
	if err == nil {
		t.Fatalf("Init: expected git is-repository error, got nil")
	}
}

func TestRenderTemplate_UnknownNameReturnsError(t *testing.T) {
	// Why: ParseFS-error path is the only renderTemplate failure
	// reachable from outside; pins it so a future template-loader
	// refactor cannot accidentally swallow it.
	_, err := application.RenderTemplateForTest("does-not-exist.tmpl", "demo")
	if err == nil {
		t.Fatalf("RenderTemplate(unknown): expected error, got nil")
	}
	if !strings.Contains(err.Error(), "parse template") {
		t.Errorf("RenderTemplate(unknown): error %v does not mention parse", err)
	}
}

func TestTemplateNames_AreSorted(t *testing.T) {
	// Why: pins that the embed.FS-glob picks up exactly the
	// templates we expect; a missing template would break Init
	// silently.
	want := []string{
		"changelog.md.tmpl",
		"compose.yaml.tmpl",
		"env.example.tmpl",
		"gitignore.tmpl",
		"readme.md.tmpl",
	}
	got, err := application.TemplateNamesForTest()
	if err != nil {
		t.Fatalf("TemplateNamesForTest: %v", err)
	}
	// templateNames() already sorts; assert the sorted contract.
	if !reflect.DeepEqual(got, want) {
		t.Errorf("template names =\n  %v\nwant:\n  %v", got, want)
	}
}
