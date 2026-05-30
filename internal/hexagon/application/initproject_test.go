package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/application/managedblock"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
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
	return application.NewInitProjectService(fs, y, g, nil, nil, nil), fs, y, g
}

// newServiceWithProgress is newService plus a fakeProgress that
// records every AffectedFiles call. Tests that assert on the
// LH-FA-INIT-005 §609 affected-paths events use this constructor.
func newServiceWithProgress(t *testing.T) (*application.InitProjectService, *fakeFS, *fakeYAML, *fakeGit, *fakeProgress) {
	t.Helper()
	fs := newFakeFS()
	fs.markDirExists(testBaseDir)
	y := &fakeYAML{}
	g := &fakeGit{}
	progress := &fakeProgress{}
	return application.NewInitProjectService(fs, y, g, progress, nil, nil), fs, y, g, progress
}

// newServiceWithConfirmer is newService plus a fakeConfirmer for the
// LH-FA-INIT-004 soft-existing-detection prompt. Tests that exercise
// the soft-detection paths use this constructor.
func newServiceWithConfirmer(t *testing.T, c driven.Confirmer) (*application.InitProjectService, *fakeFS, *fakeYAML, *fakeGit) {
	t.Helper()
	fs := newFakeFS()
	fs.markDirExists(testBaseDir)
	y := &fakeYAML{}
	g := &fakeGit{}
	return application.NewInitProjectService(fs, y, g, nil, c, nil), fs, y, g
}

// newServiceWithLogger is newService plus a fakeLogger that records
// every Debug/Info/Warn/Error call. Tests that assert on the
// LH-QA-004 logger port (slice-m4-logging-port) use this constructor.
func newServiceWithLogger(t *testing.T, l driven.Logger) (*application.InitProjectService, *fakeFS, *fakeYAML, *fakeGit) {
	t.Helper()
	fs := newFakeFS()
	fs.markDirExists(testBaseDir)
	y := &fakeYAML{}
	g := &fakeGit{}
	return application.NewInitProjectService(fs, y, g, nil, nil, l), fs, y, g
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
	if !errors.Is(err, driving.ErrBaseDirMissing) {
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
	// Why: ubootYAMLConfig.Services has `omitempty`; fresh init must
	// not emit an empty services-block. Regression-guard for the
	// M5-T1 schema extension.
	if strings.Contains(got, "services") {
		t.Errorf("u-boot.yaml unexpectedly contains `services` on fresh init; got:\n%s", got)
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
	// Why: Pins the error-propagation path through executeTemplatedFiles
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
		"devcontainer/Dockerfile.tmpl",
		"devcontainer/devcontainer.json.tmpl",
		"env.example.tmpl",
		"gitignore.tmpl",
		"readme.md.tmpl",
		"services/postgres.compose.tmpl",
		"services/postgres.env.tmpl",
		"services/postgres.volume.tmpl",
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

// --- T4b: Re-Init with --force / --backup (LH-FA-INIT-005 §611–§619) ---

// seedManagedBlockFile writes a synthetic file that already contains a
// canonical `U-BOOT MANAGED BLOCK: init` (hash style), plus user
// content outside the block. Used by the --force re-init tests to
// verify that non-managed content survives.
func seedManagedBlockFile(t *testing.T, fs *fakeFS, path, marker, userContent string) {
	t.Helper()
	body := marker + "\n# managed: old\n# END U-BOOT MANAGED BLOCK: init\n" + userContent
	if err := fs.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("seed %s: %v", path, err)
	}
}

func TestInit_Force_ManagedBlock_ReplacesOnlyBlock(t *testing.T) {
	// Why: LH-FA-INIT-005 §613–§614 — non-managed content must survive
	// --force when a marker block is present.
	svc, fs, _, _ := newService(t)
	composePath := filepath.Join(testBaseDir, "compose.yaml")
	userContent := "\n# user-added service stub below\nservices:\n  app: { build: . }\n"
	seedManagedBlockFile(t, fs, composePath, "# BEGIN U-BOOT MANAGED BLOCK: init", userContent)

	resp, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Force:   true,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	got, err := fs.ReadFile(composePath)
	if err != nil {
		t.Fatalf("ReadFile compose.yaml: %v", err)
	}
	gotStr := string(got)
	if strings.Contains(gotStr, "# managed: old") {
		t.Errorf("old block content not replaced: %q", gotStr)
	}
	if !strings.Contains(gotStr, "# user-added service stub below") {
		t.Errorf("user content outside block was clobbered: %q", gotStr)
	}
	if !strings.Contains(gotStr, "name: demo") {
		t.Errorf("new block content missing: %q", gotStr)
	}
	if len(resp.Backups) != 0 {
		t.Errorf("--force alone must not back up, got backups=%v", resp.Backups)
	}
}

func TestInit_ForceWithoutBlock_RequiresBackup(t *testing.T) {
	// Why: LH-FA-INIT-005 §619 — when a managed block is absent, full
	// overwrite is mandatory and the spec requires --backup.
	svc, fs, _, _ := newService(t)
	composePath := filepath.Join(testBaseDir, "compose.yaml")
	if err := fs.WriteFile(composePath, []byte("services: {}\n"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Force:   true,
		SkipGit: true,
	})
	if err == nil {
		t.Fatalf("Init: expected error, got nil")
	}
	if !errors.Is(err, driving.ErrForceRequiresBackup) {
		t.Errorf("Init: error %v does not wrap ErrForceRequiresBackup", err)
	}
}

func TestInit_ForceWithBackup_NoBlock_OverwritesAndBacksUp(t *testing.T) {
	// Why: --force + --backup on a no-block file is the path §619
	// allows for full overwrite — backup happens first.
	svc, fs, _, _ := newService(t)
	composePath := filepath.Join(testBaseDir, "compose.yaml")
	original := []byte("services: { custom: { image: foo } }\n")
	if err := fs.WriteFile(composePath, original, 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	resp, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Force:   true,
		Backup:  true,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	// Backup recorded.
	var composeBackup *driving.BackupAction
	for i := range resp.Backups {
		if resp.Backups[i].Original == "compose.yaml" {
			composeBackup = &resp.Backups[i]
		}
	}
	if composeBackup == nil {
		t.Fatalf("Backups missing compose.yaml entry: %v", resp.Backups)
	}
	backupContent, err := fs.ReadFile(composeBackup.Backup)
	if err != nil {
		t.Fatalf("ReadFile %s: %v", composeBackup.Backup, err)
	}
	if !reflect.DeepEqual(backupContent, original) {
		t.Errorf("backup content mismatch")
	}

	// Live file overwritten with template content.
	live, _ := fs.ReadFile(composePath)
	if !strings.Contains(string(live), "BEGIN U-BOOT MANAGED BLOCK: init") {
		t.Errorf("live file missing rendered template markers: %q", live)
	}
}

func TestInit_BackupOnly_FullOverwrite(t *testing.T) {
	// Why: --backup alone (literal-spec) means backup + full overwrite,
	// regardless of whether a managed block exists.
	svc, fs, _, _ := newService(t)
	composePath := filepath.Join(testBaseDir, "compose.yaml")
	seedManagedBlockFile(t, fs, composePath, "# BEGIN U-BOOT MANAGED BLOCK: init", "\n# user-added stuff\n")

	resp, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Backup:  true,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	if len(resp.Backups) == 0 {
		t.Fatalf("expected at least one backup, got none")
	}
	live, _ := fs.ReadFile(composePath)
	if strings.Contains(string(live), "# user-added stuff") {
		t.Errorf("--backup alone should full-overwrite (lose user content): %q", live)
	}
}

func TestInit_ForceAndBackup_ManagedBlock_BackupsAndReplacesBlock(t *testing.T) {
	// Why: --force + --backup + managed block → block-only edit with
	// safety backup. Backup taken, user content outside block kept.
	svc, fs, _, _ := newService(t)
	composePath := filepath.Join(testBaseDir, "compose.yaml")
	userContent := "\nservices: { app: { image: foo } }\n"
	seedManagedBlockFile(t, fs, composePath, "# BEGIN U-BOOT MANAGED BLOCK: init", userContent)

	resp, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Force:   true,
		Backup:  true,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	var composeBackup *driving.BackupAction
	for i := range resp.Backups {
		if resp.Backups[i].Original == "compose.yaml" {
			composeBackup = &resp.Backups[i]
		}
	}
	if composeBackup == nil {
		t.Errorf("--backup with --force should still create a backup")
	}
	live, _ := fs.ReadFile(composePath)
	if !strings.Contains(string(live), "services: { app: { image: foo } }") {
		t.Errorf("user content lost: %q", live)
	}
}

func TestInit_NonManagedFile_ForceWithoutBackup_RequiresBackup(t *testing.T) {
	// Why: .gitignore is intentionally not in the §611 managed-block
	// list — re-init treats it as fully managed (--backup required).
	svc, fs, _, _ := newService(t)
	if err := fs.WriteFile(filepath.Join(testBaseDir, ".gitignore"), []byte("*.tmp\n"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Force:   true,
		SkipGit: true,
	})
	if err == nil {
		t.Fatalf("Init: expected error, got nil")
	}
	if !errors.Is(err, driving.ErrForceRequiresBackup) {
		t.Errorf("Init: error %v does not wrap ErrForceRequiresBackup", err)
	}
}

func TestInit_UBootYAML_ForceWithoutBackup_RequiresBackup(t *testing.T) {
	// Why: u-boot.yaml is fully u-boot-managed (no inline block) per
	// LH-SA-FILE-002 §615 strict-JSON/steering-file fallback. Re-init
	// requires --backup.
	svc, fs, _, _ := newService(t)
	if err := fs.WriteFile(filepath.Join(testBaseDir, "u-boot.yaml"), []byte("schemaVersion: 0\n"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Force:   true,
		SkipGit: true,
	})
	if err == nil {
		t.Fatalf("Init: expected error, got nil")
	}
	if !errors.Is(err, driving.ErrForceRequiresBackup) {
		t.Errorf("Init: error %v does not wrap ErrForceRequiresBackup", err)
	}
}

func TestInit_Summary_EmittedOnReInit(t *testing.T) {
	// Why: LH-FA-INIT-005 §609 / LH-FA-CLI-005A §262 — affected
	// paths must be reported BEFORE the write. With T4c-review the
	// reporting goes through a structured port (not a text writer);
	// assert on the recorded event shape.
	svc, fs, _, _, progress := newServiceWithProgress(t)
	composePath := filepath.Join(testBaseDir, "compose.yaml")
	seedManagedBlockFile(t, fs, composePath, "# BEGIN U-BOOT MANAGED BLOCK: init", "")

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Force:   true,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	if len(progress.calls) != 1 {
		t.Fatalf("AffectedFiles calls = %d, want 1", len(progress.calls))
	}
	call := progress.calls[0]
	if call.BaseDir != testBaseDir {
		t.Errorf("baseDir = %q, want %q", call.BaseDir, testBaseDir)
	}
	found := false
	for _, row := range call.Rows {
		if row.Path == "compose.yaml" && row.Action == driven.AffectedReplaceBlock {
			found = true
		}
	}
	if !found {
		t.Errorf("compose.yaml/ReplaceBlock event missing: %v", call.Rows)
	}
}

func TestInit_Summary_QuietOnFreshInit(t *testing.T) {
	// Why: defensive — fresh init must not emit a summary; nothing
	// is being overwritten. Port not called at all (so a no-op
	// adapter is unobservable).
	svc, _, _, _, progress := newServiceWithProgress(t)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if len(progress.calls) != 0 {
		t.Errorf("expected 0 AffectedFiles calls on fresh init, got %d", len(progress.calls))
	}
}

func TestInit_Summary_WithBackupMarker(t *testing.T) {
	// Why: pin that the event for a --backup overwrite carries
	// Backup=true so the adapter can render "(with backup)".
	svc, fs, _, _, progress := newServiceWithProgress(t)
	if err := fs.WriteFile(filepath.Join(testBaseDir, "compose.yaml"), []byte("services: {}\n"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Backup:  true,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if len(progress.calls) != 1 {
		t.Fatalf("AffectedFiles calls = %d, want 1", len(progress.calls))
	}
	hasBackupRow := false
	for _, row := range progress.calls[0].Rows {
		if row.Backup {
			hasBackupRow = true
		}
	}
	if !hasBackupRow {
		t.Errorf("expected at least one row with Backup=true, got %v", progress.calls[0].Rows)
	}
}

func TestInit_NilProgress_TolerantToNoop(t *testing.T) {
	// Why: constructor must accept nil progress without panicking;
	// it falls back to an internal no-op ProgressPort.
	fs := newFakeFS()
	fs.markDirExists(testBaseDir)
	svc := application.NewInitProjectService(fs, &fakeYAML{}, &fakeGit{}, nil, nil, nil)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
}

func TestInit_Backup_RecordedInResponse(t *testing.T) {
	// Why: response.Backups must enumerate the actions so the CLI
	// (T4c) can render "backed up X to Y" lines to the user.
	svc, fs, _, _ := newService(t)
	if err := fs.WriteFile(filepath.Join(testBaseDir, ".env.example"), []byte("OLD=1\n"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	resp, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Force:   true,
		Backup:  true,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if len(resp.Backups) == 0 {
		t.Fatalf("Backups empty, want at least 1 entry")
	}
	// Match either .env.example.bak or .env.example.bak.<n> so a
	// future setup that seeds a stale backup doesn't flip the test
	// (review finding #7).
	suffixRE := regexp.MustCompile(`\.env\.example\.bak(\.\d+)?$`)
	found := false
	for _, b := range resp.Backups {
		if b.Original == ".env.example" && suffixRE.MatchString(b.Backup) {
			found = true
		}
	}
	if !found {
		t.Errorf(".env.example backup entry not found in %v", resp.Backups)
	}
}

func TestInit_PlanErrorAbortsBeforeAnyWrite(t *testing.T) {
	// Why: spec requires no partial side effects when a plan-phase
	// error fires. Setup: .env.example exists without flags →
	// ErrProjectExists; assert that README.md / CHANGELOG.md /
	// compose.yaml are NOT written.
	svc, fs, _, _ := newService(t)
	if err := fs.WriteFile(filepath.Join(testBaseDir, ".env.example"), []byte("X=1\n"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	})
	if !errors.Is(err, driving.ErrProjectExists) {
		t.Fatalf("Init: want ErrProjectExists, got %v", err)
	}
	for _, name := range []string{"README.md", "CHANGELOG.md", "compose.yaml"} {
		if exists, _ := fs.Exists(filepath.Join(testBaseDir, name)); exists {
			t.Errorf("plan-error did not prevent write of %s", name)
		}
	}
}

func TestInit_RenderedTemplate_ContainsManagedBlockMarkers(t *testing.T) {
	// Why: pin that the embedded templates produce LH-SA-FILE-002-
	// compliant markers so the re-init replace path can find them.
	svc, fs, _, _ := newService(t)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	cases := map[string][2]string{
		"compose.yaml":  {"# BEGIN U-BOOT MANAGED BLOCK: init", "# END U-BOOT MANAGED BLOCK: init"},
		".env.example":  {"# BEGIN U-BOOT MANAGED BLOCK: init", "# END U-BOOT MANAGED BLOCK: init"},
		"README.md":     {"<!-- BEGIN U-BOOT MANAGED BLOCK: init -->", "<!-- END U-BOOT MANAGED BLOCK: init -->"},
		"CHANGELOG.md":  {"<!-- BEGIN U-BOOT MANAGED BLOCK: init -->", "<!-- END U-BOOT MANAGED BLOCK: init -->"},
	}
	for path, markers := range cases {
		body, err := fs.ReadFile(filepath.Join(testBaseDir, path))
		if err != nil {
			t.Errorf("ReadFile %s: %v", path, err)
			continue
		}
		s := string(body)
		if !strings.Contains(s, markers[0]) {
			t.Errorf("%s missing BEGIN marker %q in: %q", path, markers[0], s)
		}
		if !strings.Contains(s, markers[1]) {
			t.Errorf("%s missing END marker %q in: %q", path, markers[1], s)
		}
	}
	// .gitignore intentionally has no markers (whole-file managed per
	// §611 list exclusion).
	gitignore, err := fs.ReadFile(filepath.Join(testBaseDir, ".gitignore"))
	if err != nil {
		t.Fatalf("ReadFile .gitignore: %v", err)
	}
	if strings.Contains(string(gitignore), "BEGIN U-BOOT MANAGED BLOCK") {
		t.Errorf(".gitignore must not have managed-block markers (spec §611 list excludes it): %q", gitignore)
	}
}

// --- T4b-Review-Fixes: Lstat / Mode-Preservation / Sentinel-Split ---

func TestInit_Symlink_AtTemplatePath_Rejected(t *testing.T) {
	// Why: review finding #2 — silently following a `.env.example ->
	// /etc/passwd` symlink would let the re-init read and overwrite
	// the link target. Reject with ErrBackupUnsupportedKind (same
	// sentinel BackupPath uses in T4a for the same class of bug).
	svc, fs, _, _ := newService(t)
	fs.markSymlink(filepath.Join(testBaseDir, ".env.example"))

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	})
	if err == nil {
		t.Fatalf("Init: expected error, got nil")
	}
	if !errors.Is(err, driving.ErrBackupUnsupportedKind) {
		t.Errorf("Init: error %v does not wrap ErrBackupUnsupportedKind", err)
	}
}

func TestInit_Mode_PreservedAcrossReplaceBlock(t *testing.T) {
	// Why: review finding #1 — T4a-review enforced mode preservation
	// for backups; T4b's write paths must do the same. Setup: existing
	// compose.yaml with mode 0o600 (e.g. user chmod'd it). Re-init
	// with --force should refresh the managed block but keep 0o600.
	svc, fs, _, _ := newService(t)
	composePath := filepath.Join(testBaseDir, "compose.yaml")
	body := "# BEGIN U-BOOT MANAGED BLOCK: init\nold\n# END U-BOOT MANAGED BLOCK: init\n"
	if err := fs.WriteFile(composePath, []byte(body), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Force:   true,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	info, err := fs.Lstat(composePath)
	if err != nil {
		t.Fatalf("Lstat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("mode = %o, want 0o600 (preserved across ReplaceBlock)", info.Mode().Perm())
	}
}

func TestInit_Mode_PreservedAcrossOverwriteFull(t *testing.T) {
	// Why: same as the ReplaceBlock test, but for the OverwriteFull
	// path (no managed block in source, --backup forces full overwrite).
	svc, fs, _, _ := newService(t)
	envPath := filepath.Join(testBaseDir, ".env.example")
	if err := fs.WriteFile(envPath, []byte("FOO=1\n"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Backup:  true,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	info, err := fs.Lstat(envPath)
	if err != nil {
		t.Fatalf("Lstat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("mode = %o, want 0o600 (preserved across OverwriteFull)", info.Mode().Perm())
	}
}

func TestInit_FreshFile_GetsDefaultMode(t *testing.T) {
	// Why: pin the fallback — actionWrite (file didn't exist) uses
	// 0o644 because there's no source mode to preserve.
	svc, fs, _, _ := newService(t)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	info, err := fs.Lstat(filepath.Join(testBaseDir, "README.md"))
	if err != nil {
		t.Fatalf("Lstat: %v", err)
	}
	if info.Mode().Perm() != 0o644 {
		t.Errorf("fresh README.md mode = %o, want 0o644 (default)", info.Mode().Perm())
	}
}

func TestInit_NonMarkerFile_Collision_ReturnsErrFileExists(t *testing.T) {
	// Why: review finding #5 — a stray README.md is not proof of an
	// existing u-boot project. Sentinel split keeps the message
	// honest: ErrFileExists for non-markers, ErrProjectExists for
	// real markers (separately tested in TestInit_ExistingProjectRejected).
	svc, fs, _, _ := newService(t)
	if err := fs.WriteFile(filepath.Join(testBaseDir, "README.md"), []byte("My personal README\n"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	})
	if err == nil {
		t.Fatalf("Init: expected error, got nil")
	}
	if !errors.Is(err, driving.ErrFileExists) {
		t.Errorf("Init: error %v does not wrap ErrFileExists", err)
	}
	if errors.Is(err, driving.ErrProjectExists) {
		t.Errorf("Init: stray README.md should NOT trip ErrProjectExists (markers only)")
	}
}

func TestInit_MarkerFile_Collision_StillReturnsErrProjectExists(t *testing.T) {
	// Why: cross-check of the split — u-boot.yaml IS a marker, so
	// ErrProjectExists fires (not ErrFileExists). ErrProjectExists is
	// also distinct enough that callers branching on it keep working.
	svc, fs, _, _ := newService(t)
	if err := fs.WriteFile(filepath.Join(testBaseDir, "u-boot.yaml"), []byte("x\n"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	})
	if !errors.Is(err, driving.ErrProjectExists) {
		t.Fatalf("Init: want ErrProjectExists for u-boot.yaml marker, got %v", err)
	}
	if errors.Is(err, driving.ErrFileExists) {
		t.Errorf("Init: u-boot.yaml is a marker — must not double-trip ErrFileExists")
	}
}

func TestInit_PlanCachesBody_NoDoubleRead(t *testing.T) {
	// Why: review finding #8 — plan caches the body in filePlan.Body
	// so execute does not re-read (closes one TOCTOU window and one
	// I/O syscall). Asserted via the ReadFile call counter on fakeFS:
	// after a --force re-init with managed-block-replace, the file
	// must have been read exactly once.
	svc, fs, _, _ := newService(t)
	composePath := filepath.Join(testBaseDir, "compose.yaml")
	originalBody := "# BEGIN U-BOOT MANAGED BLOCK: init\nold\n# END U-BOOT MANAGED BLOCK: init\nbelow-block-user-content\n"
	if err := fs.WriteFile(composePath, []byte(originalBody), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Force:   true,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	if got := fs.readFileCallCount(composePath); got != 1 {
		t.Errorf("ReadFile(compose.yaml) called %d times, want 1 (plan must cache body for execute)", got)
	}
	// And the user content outside the block survived — proves the
	// cached body really was used for the splice.
	final, _ := fs.ReadFile(composePath)
	if !strings.Contains(string(final), "below-block-user-content") {
		t.Errorf("user content not preserved: %q", final)
	}
}

// ----------------------------------------------------------------------------
// LH-FA-INIT-004 — soft-existing-detection paths
// ----------------------------------------------------------------------------

// seedSoftIndicators primes the fake FS with N of the LH-FA-INIT-004
// soft indicators so a test can choose how many are present.
//
// Collision-safe ordering: the directories `docs`, `scripts`, `docker`
// are seeded first (init's MkdirAll is idempotent, so seeding does
// not trip planFile's per-file collision) followed by
// `.devcontainer/devcontainer.json` (not in planFile's template list,
// also collision-free). `README.md` and `CHANGELOG.md` come last
// because they ARE in the template list — seeding them additionally
// causes planFile to abort with ErrFileExists, useful for tests that
// want to observe the post-detection collision path.
//
// Pass n=2 to stay below the threshold; n=3..4 to cross the threshold
// without per-file collisions; n=5..6 to cross the threshold AND
// trigger per-file collisions.
func seedSoftIndicators(t *testing.T, fs *fakeFS, n int) {
	t.Helper()
	candidates := []string{
		"docs",                            // dir, collision-safe
		"scripts",                         // dir, collision-safe
		"docker",                          // dir, collision-safe
		".devcontainer/devcontainer.json", // file, not in template list
		"README.md",                       // file, IN template list (collides)
		"CHANGELOG.md",                    // file, IN template list (collides)
	}
	if n > len(candidates) {
		t.Fatalf("seedSoftIndicators: n=%d > %d candidates", n, len(candidates))
	}
	for _, rel := range candidates[:n] {
		full := filepath.Join(testBaseDir, rel)
		if rel == "docs" || rel == "scripts" || rel == "docker" {
			fs.markDirExists(full)
			continue
		}
		if err := fs.WriteFile(full, []byte("seed\n"), 0o644); err != nil {
			t.Fatalf("seedSoftIndicators: write %s: %v", rel, err)
		}
	}
}

func TestInit_SoftDetect_Under3Indicators_Proceeds(t *testing.T) {
	// Why: <3 soft indicators must NOT trigger the soft-detection
	// abort — the init proceeds as fresh.
	confirmer := &fakeConfirmer{}
	svc, fs, _, _ := newServiceWithConfirmer(t, confirmer)
	seedSoftIndicators(t, fs, 2)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v (expected fresh-init success)", err)
	}
	if len(confirmer.calls) != 0 {
		t.Errorf("Confirmer called %d times, want 0 (below threshold)", len(confirmer.calls))
	}
}

func TestInit_SoftDetect_AssumeExisting_AbortsWithoutPrompt(t *testing.T) {
	// Why: --assume-existing turns soft-detection into a deterministic
	// abort. Confirmer must not be called (no prompt in non-
	// interactive runs that asserted existence).
	confirmer := &fakeConfirmer{}
	svc, fs, _, _ := newServiceWithConfirmer(t, confirmer)
	seedSoftIndicators(t, fs, 4)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:           "demo",
		BaseDir:        testBaseDir,
		SkipGit:        true,
		AssumeExisting: true,
	})
	if !errors.Is(err, driving.ErrProjectExists) {
		t.Fatalf("err = %v, want wrapped ErrProjectExists", err)
	}
	if !strings.Contains(err.Error(), "--assume-existing") {
		t.Errorf("err message should name --assume-existing trigger: %v", err)
	}
	if len(confirmer.calls) != 0 {
		t.Errorf("Confirmer called %d times, want 0 (AssumeExisting short-circuits)", len(confirmer.calls))
	}
}

func TestInit_SoftDetect_NoInteractive_SkipsDetectionAndConfirmer(t *testing.T) {
	// Why: LH-FA-INIT-004 §247 — in non-interactive mode without
	// --assume-existing, the soft-detection does not fire. The
	// service must proceed; the per-file collision logic in planFile
	// will still trip on README.md → ErrFileExists, but that is the
	// pre-existing behaviour, not the soft-detection abort.
	// n=5 includes README.md so the post-detection collision path is
	// observable in the same test.
	confirmer := &fakeConfirmer{answer: true} // would abort if asked
	svc, fs, _, _ := newServiceWithConfirmer(t, confirmer)
	seedSoftIndicators(t, fs, 5)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:          "demo",
		BaseDir:       testBaseDir,
		SkipGit:       true,
		NoInteractive: true,
	})
	if len(confirmer.calls) != 0 {
		t.Errorf("Confirmer called %d times, want 0 (NoInteractive skips prompt)", len(confirmer.calls))
	}
	// We do hit the planFile collision on README.md (ErrFileExists).
	if !errors.Is(err, driving.ErrFileExists) {
		t.Errorf("err = %v, want wrapped ErrFileExists (per-file collision)", err)
	}
}

func TestInit_SoftDetect_Interactive_ConfirmerYesAborts(t *testing.T) {
	// Why: interactive mode + Confirmer says yes → soft-detection
	// abort fires with ErrProjectExists naming "user confirmation".
	confirmer := &fakeConfirmer{answer: true}
	svc, fs, _, _ := newServiceWithConfirmer(t, confirmer)
	seedSoftIndicators(t, fs, 4)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	})
	if !errors.Is(err, driving.ErrProjectExists) {
		t.Fatalf("err = %v, want wrapped ErrProjectExists", err)
	}
	if !strings.Contains(err.Error(), "user confirmation") {
		t.Errorf("err message should name user-confirmation trigger: %v", err)
	}
	if len(confirmer.calls) != 1 {
		t.Fatalf("Confirmer called %d times, want 1", len(confirmer.calls))
	}
	if confirmer.calls[0].BaseDir != testBaseDir {
		t.Errorf("BaseDir passed to Confirmer = %q, want %q", confirmer.calls[0].BaseDir, testBaseDir)
	}
	if len(confirmer.calls[0].Indicators) != 4 {
		t.Errorf("Indicators passed = %v (want 4)", confirmer.calls[0].Indicators)
	}
}

func TestInit_SoftDetect_Interactive_ConfirmerNoProceeds(t *testing.T) {
	// Why: interactive mode + Confirmer says no → soft-detection does
	// not abort; the service proceeds. planFile still trips on the
	// per-file collision (README.md → ErrFileExists), which is the
	// deterministic pre-existing behaviour. n=5 includes README.md so
	// the collision is observable.
	confirmer := &fakeConfirmer{answer: false}
	svc, fs, _, _ := newServiceWithConfirmer(t, confirmer)
	seedSoftIndicators(t, fs, 5)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	})
	if len(confirmer.calls) != 1 {
		t.Errorf("Confirmer called %d times, want 1", len(confirmer.calls))
	}
	if !errors.Is(err, driving.ErrFileExists) {
		t.Errorf("err = %v, want wrapped ErrFileExists (Confirmer said no, planFile collides)", err)
	}
	if errors.Is(err, driving.ErrProjectExists) {
		t.Errorf("err must not be ErrProjectExists when Confirmer said no: %v", err)
	}
}

func TestInit_SoftDetect_ForceBackup_SkipsDetection(t *testing.T) {
	// Why: --force / --backup already opt into re-init explicitly.
	// Soft-detection must NOT call the Confirmer in that path, and
	// must NOT abort — planFile owns the per-file action choice.
	confirmer := &fakeConfirmer{answer: true} // would abort if asked
	svc, fs, _, _ := newServiceWithConfirmer(t, confirmer)
	seedSoftIndicators(t, fs, 4)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
		Backup:  true, // explicit opt-in
	})
	if err != nil {
		t.Fatalf("Init with --backup: %v", err)
	}
	if len(confirmer.calls) != 0 {
		t.Errorf("Confirmer called %d times, want 0 (--backup skips detection)", len(confirmer.calls))
	}
}

func TestInit_SoftDetect_Logger_DebugFiresAboveThreshold(t *testing.T) {
	// Why: LH-QA-004 driven.Logger port — verifies the soft-detection
	// emits a Debug entry once the ≥3-indicator threshold is crossed.
	// Below the threshold the logger must stay silent (test below).
	logger := &fakeLogger{}
	svc, fs, _, _ := newServiceWithLogger(t, logger)
	seedSoftIndicators(t, fs, 4) // no per-file collision at n=4

	_, _ = svc.Init(context.Background(), driving.InitProjectRequest{
		Name:           "demo",
		BaseDir:        testBaseDir,
		SkipGit:        true,
		AssumeExisting: true, // aborts cleanly so test stays focused on the log call
	})

	var found bool
	for _, e := range logger.entries {
		if e.Level == "DEBUG" && strings.Contains(e.Msg, "soft-existing-detection") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected DEBUG log entry for soft-existing-detection, got %+v", logger.entries)
	}
}

func TestInit_SoftDetect_Logger_SilentBelowThreshold(t *testing.T) {
	// Why: the soft-detection Debug must NOT fire when the indicator
	// count is below the threshold — keeps the log volume meaningful.
	logger := &fakeLogger{}
	svc, fs, _, _ := newServiceWithLogger(t, logger)
	seedSoftIndicators(t, fs, 2)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v (expected fresh-init success)", err)
	}
	for _, e := range logger.entries {
		if strings.Contains(e.Msg, "soft-existing-detection") {
			t.Errorf("soft-detection logged below threshold: %+v", e)
		}
	}
}

func TestInit_SoftDetect_ConfirmerError_Propagates(t *testing.T) {
	// Why: Confirmer I/O failures must surface to the CLI as an
	// abort, not be silently coerced to "no". The init does not
	// proceed in that case.
	confirmer := &fakeConfirmer{err: errors.New("simulated stdin failure")}
	svc, fs, _, _ := newServiceWithConfirmer(t, confirmer)
	seedSoftIndicators(t, fs, 4)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	})
	if err == nil {
		t.Fatal("expected propagated Confirmer error, got nil")
	}
	if !strings.Contains(err.Error(), "simulated stdin failure") {
		t.Errorf("error did not wrap Confirmer error: %v", err)
	}
}

// --- M5-T4a: compose scaffold restructuring + renderManagedBlockOnly + Ensure-Scaffold ----

// TestInit_T4a_FreshInit_ComposeSplitBlockForm pins the new T4a
// scaffold shape: the init-block contains name + networks but NOT
// services: {}, and services:/volumes: live as empty top-level maps
// OUTSIDE the init block. Add-ons placed under those maps by `u-boot
// add` survive a later `u-boot init --force`.
func TestInit_T4a_FreshInit_ComposeSplitBlockForm(t *testing.T) {
	svc, fs, _, _ := newService(t)
	if _, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	}); err != nil {
		t.Fatalf("Init: %v", err)
	}
	body, err := fs.ReadFile(filepath.Join(testBaseDir, "compose.yaml"))
	if err != nil {
		t.Fatalf("ReadFile compose.yaml: %v", err)
	}
	got := string(body)

	// init block contains the name + networks block, not services.
	beginIdx := strings.Index(got, "# BEGIN U-BOOT MANAGED BLOCK: init")
	endIdx := strings.Index(got, "# END U-BOOT MANAGED BLOCK: init")
	if beginIdx < 0 || endIdx < 0 || endIdx < beginIdx {
		t.Fatalf("init block markers missing or reversed; got:\n%s", got)
	}
	initBlock := got[beginIdx:endIdx]
	if !strings.Contains(initBlock, "name: demo") {
		t.Errorf("init block missing name: demo; got:\n%s", initBlock)
	}
	if !strings.Contains(initBlock, "networks:") {
		t.Errorf("init block missing networks; got:\n%s", initBlock)
	}
	if strings.Contains(initBlock, "services:") {
		t.Errorf("init block must NOT contain services anymore (T4a); got:\n%s", initBlock)
	}
	if strings.Contains(initBlock, "volumes:") {
		t.Errorf("init block must NOT contain volumes; got:\n%s", initBlock)
	}

	// services: and volumes: live OUTSIDE the init block.
	afterBlock := got[endIdx:]
	if !strings.Contains(afterBlock, "\nservices: {}") {
		t.Errorf("top-level services: {} missing after init block; got:\n%s", afterBlock)
	}
	if !strings.Contains(afterBlock, "\nvolumes: {}") {
		t.Errorf("top-level volumes: {} missing after init block; got:\n%s", afterBlock)
	}
}

// TestInit_T4a_Force_PreservesAddonBlocks is the central T4a
// regression: a re-init with --force on a compose.yaml that already
// holds add-on managed blocks (service.postgres / volume.postgres)
// MUST keep those blocks. T3's executeReplaceBlock used to splice
// the entire rendered template into the init region; T4a splices
// only the init-block byte range so the surrounding services:/
// volumes: maps (with their add-on sub-markers) survive untouched.
func TestInit_T4a_Force_PreservesAddonBlocks(t *testing.T) {
	svc, fs, _, _ := newService(t)
	composePath := filepath.Join(testBaseDir, "compose.yaml")
	seeded := "# BEGIN U-BOOT MANAGED BLOCK: init\n" +
		"# managed: old\n" +
		"# END U-BOOT MANAGED BLOCK: init\n" +
		"\n" +
		"services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  postgres:\n" +
		"    image: postgres:16-alpine\n" +
		"  # END U-BOOT MANAGED BLOCK: service.postgres\n" +
		"\n" +
		"volumes:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: volume.postgres\n" +
		"  postgres-data: {}\n" +
		"  # END U-BOOT MANAGED BLOCK: volume.postgres\n"
	if err := fs.WriteFile(composePath, []byte(seeded), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if _, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Force:   true,
		SkipGit: true,
	}); err != nil {
		t.Fatalf("Init: %v", err)
	}

	body, err := fs.ReadFile(composePath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	got := string(body)

	// Old init body gone, new init body in.
	if strings.Contains(got, "# managed: old") {
		t.Errorf("old init block content not replaced; got:\n%s", got)
	}
	if !strings.Contains(got, "name: demo") {
		t.Errorf("new init body missing name: demo; got:\n%s", got)
	}

	// Add-on markers survived.
	for _, want := range []string{
		"# BEGIN U-BOOT MANAGED BLOCK: service.postgres",
		"# END U-BOOT MANAGED BLOCK: service.postgres",
		"# BEGIN U-BOOT MANAGED BLOCK: volume.postgres",
		"# END U-BOOT MANAGED BLOCK: volume.postgres",
		"image: postgres:16-alpine",
		"postgres-data: {}",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("add-on artefact %q missing after --force; got:\n%s", want, got)
		}
	}

	// Exactly one top-level services: / volumes: (no duplication).
	if c := strings.Count(got, "\nservices:"); c != 1 {
		t.Errorf("expected exactly 1 top-level services:, got %d; full:\n%s", c, got)
	}
	if c := strings.Count(got, "\nvolumes:"); c != 1 {
		t.Errorf("expected exactly 1 top-level volumes:, got %d; full:\n%s", c, got)
	}

	// Exactly one init block.
	if c := strings.Count(got, "# BEGIN U-BOOT MANAGED BLOCK: init"); c != 1 {
		t.Errorf("expected exactly 1 init BEGIN marker, got %d", c)
	}
}

// TestInit_T4a_AltM3_Migration covers the most likely upgrade path:
// an existing M3-form compose.yaml that holds `services: {}` and
// `networks: default:` INSIDE the init block and no `volumes:` at
// all. `u-boot init --force` must produce a clean split-block form
// (services:/volumes: as top-level maps outside the block) without
// dropping content the user might rely on.
func TestInit_T4a_AltM3_Migration(t *testing.T) {
	svc, fs, _, _ := newService(t)
	composePath := filepath.Join(testBaseDir, "compose.yaml")
	altM3 := "# BEGIN U-BOOT MANAGED BLOCK: init\n" +
		"# Compose stack for demo.\n" +
		"\n" +
		"name: demo\n" +
		"\n" +
		"services: {}\n" +
		"\n" +
		"networks:\n" +
		"  default:\n" +
		"    name: demo-default\n" +
		"# END U-BOOT MANAGED BLOCK: init\n"
	if err := fs.WriteFile(composePath, []byte(altM3), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if _, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Force:   true,
		SkipGit: true,
	}); err != nil {
		t.Fatalf("Init: %v", err)
	}

	body, err := fs.ReadFile(composePath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	got := string(body)

	// init block stays, name + networks still in it, services no longer in it.
	beginIdx := strings.Index(got, "# BEGIN U-BOOT MANAGED BLOCK: init")
	endIdx := strings.Index(got, "# END U-BOOT MANAGED BLOCK: init")
	if beginIdx < 0 || endIdx < beginIdx {
		t.Fatalf("init markers reversed/missing; got:\n%s", got)
	}
	initBlock := got[beginIdx:endIdx]
	if strings.Contains(initBlock, "services:") {
		t.Errorf("init block must NOT contain services after migration; got:\n%s", initBlock)
	}
	if !strings.Contains(initBlock, "networks:") {
		t.Errorf("init block must retain networks after migration; got:\n%s", initBlock)
	}

	// New top-level services: and volumes: outside the block, both empty.
	afterBlock := got[endIdx:]
	if !strings.Contains(afterBlock, "\nservices: {}") {
		t.Errorf("top-level services: {} missing after migration; got:\n%s", afterBlock)
	}
	if !strings.Contains(afterBlock, "\nvolumes: {}") {
		t.Errorf("top-level volumes: {} missing after migration (Alt-M3 had no volumes); got:\n%s", afterBlock)
	}
}

// TestInit_T4a_EnsureScaffold_UserCustomKey is the Ensure-Scaffold
// safety check: a user-added top-level key (e.g. `x-user-config`)
// outside any managed block must not be touched, duplicated, or
// re-ordered by the ensure-pass.
func TestInit_T4a_EnsureScaffold_UserCustomKey(t *testing.T) {
	svc, fs, _, _ := newService(t)
	composePath := filepath.Join(testBaseDir, "compose.yaml")
	seeded := "# BEGIN U-BOOT MANAGED BLOCK: init\n" +
		"name: demo\n" +
		"# END U-BOOT MANAGED BLOCK: init\n" +
		"\n" +
		"x-user-config:\n" +
		"  region: eu-central-1\n"
	if err := fs.WriteFile(composePath, []byte(seeded), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if _, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Force:   true,
		SkipGit: true,
	}); err != nil {
		t.Fatalf("Init: %v", err)
	}

	body, _ := fs.ReadFile(composePath)
	got := string(body)

	// User block intact, exactly once.
	if c := strings.Count(got, "x-user-config:"); c != 1 {
		t.Errorf("x-user-config: expected exactly once, got %d; full:\n%s", c, got)
	}
	if !strings.Contains(got, "region: eu-central-1") {
		t.Errorf("user value lost; got:\n%s", got)
	}
	// Ensure-Scaffold appended services and volumes too.
	if !strings.Contains(got, "\nservices: {}") {
		t.Errorf("services: {} missing; got:\n%s", got)
	}
	if !strings.Contains(got, "\nvolumes: {}") {
		t.Errorf("volumes: {} missing; got:\n%s", got)
	}
}

// TestInit_T4a_BackupOnly_FullOverwrite_DropsAddonsToBak pins the
// documented-destructive behavior of `--backup` alone: live file is
// the rendered template (add-on blocks gone from live), backup file
// contains them. Users who want the add-on blocks preserved on
// re-init must use `--force --backup`.
func TestInit_T4a_BackupOnly_FullOverwrite_DropsAddonsToBak(t *testing.T) {
	svc, fs, _, _ := newService(t)
	composePath := filepath.Join(testBaseDir, "compose.yaml")
	seeded := "# BEGIN U-BOOT MANAGED BLOCK: init\n" +
		"name: demo\n" +
		"# END U-BOOT MANAGED BLOCK: init\n" +
		"\n" +
		"services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  postgres:\n" +
		"    image: postgres:16-alpine\n" +
		"  # END U-BOOT MANAGED BLOCK: service.postgres\n"
	if err := fs.WriteFile(composePath, []byte(seeded), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	resp, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Backup:  true,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	live, _ := fs.ReadFile(composePath)
	liveStr := string(live)
	if strings.Contains(liveStr, "service.postgres") {
		t.Errorf("live file still contains add-on block after --backup-only (expected destructive); got:\n%s", liveStr)
	}
	if !strings.Contains(liveStr, "services: {}") || !strings.Contains(liveStr, "volumes: {}") {
		t.Errorf("live file missing fresh top-level scaffold; got:\n%s", liveStr)
	}

	// Backup retains the add-on block.
	var bak *driving.BackupAction
	for i := range resp.Backups {
		if resp.Backups[i].Original == "compose.yaml" {
			bak = &resp.Backups[i]
		}
	}
	if bak == nil {
		t.Fatalf("no backup recorded; resp.Backups=%v", resp.Backups)
	}
	bakBody, _ := fs.ReadFile(bak.Backup)
	if !strings.Contains(string(bakBody), "service.postgres") {
		t.Errorf("backup file missing add-on block; got:\n%s", bakBody)
	}
}

// TestInit_T4a_ForceWithBackup_ManagedBlock_PreservesAddons mirrors
// the --force-alone test but with --backup also set: the splice
// mechanic is the same, plus a backup of the pre-state.
func TestInit_T4a_ForceWithBackup_ManagedBlock_PreservesAddons(t *testing.T) {
	svc, fs, _, _ := newService(t)
	composePath := filepath.Join(testBaseDir, "compose.yaml")
	seeded := "# BEGIN U-BOOT MANAGED BLOCK: init\n" +
		"name: stale\n" +
		"# END U-BOOT MANAGED BLOCK: init\n" +
		"\n" +
		"services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  postgres:\n" +
		"    image: postgres:16-alpine\n" +
		"  # END U-BOOT MANAGED BLOCK: service.postgres\n" +
		"\n" +
		"volumes:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: volume.postgres\n" +
		"  postgres-data: {}\n" +
		"  # END U-BOOT MANAGED BLOCK: volume.postgres\n"
	if err := fs.WriteFile(composePath, []byte(seeded), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	resp, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		Force:   true,
		Backup:  true,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	live, _ := fs.ReadFile(composePath)
	if !strings.Contains(string(live), "service.postgres") {
		t.Errorf("--force --backup must preserve add-on blocks; got:\n%s", live)
	}

	// Backup recorded too.
	hasBackup := false
	for _, b := range resp.Backups {
		if b.Original == "compose.yaml" {
			hasBackup = true
		}
	}
	if !hasBackup {
		t.Errorf("--force --backup must record a backup; got Backups=%v", resp.Backups)
	}
}

// TestRenderManagedBlockOnly_BlockNotFound covers the helper-direct
// programmer-error path: a template body without the requested
// marker.
func TestRenderManagedBlockOnly_BlockNotFound(t *testing.T) {
	_, err := application.RenderManagedBlockOnlyForTest(
		[]byte("name: demo\n"),
		"init",
	)
	if err == nil {
		t.Fatalf("expected ErrBlockNotFound, got nil")
	}
	if !errors.Is(err, managedblock.ErrBlockNotFound) {
		t.Errorf("error %v does not wrap ErrBlockNotFound", err)
	}
}

// TestRenderManagedBlockOnly_BlockMalformed covers the
// duplicate-BEGIN malformed path.
func TestRenderManagedBlockOnly_BlockMalformed(t *testing.T) {
	body := []byte(
		"# BEGIN U-BOOT MANAGED BLOCK: init\n" +
			"# BEGIN U-BOOT MANAGED BLOCK: init\n" +
			"# END U-BOOT MANAGED BLOCK: init\n",
	)
	_, err := application.RenderManagedBlockOnlyForTest(body, "init")
	if err == nil {
		t.Fatalf("expected ErrBlockMalformed, got nil")
	}
	if !errors.Is(err, managedblock.ErrBlockMalformed) {
		t.Errorf("error %v does not wrap ErrBlockMalformed", err)
	}
}

// TestRenderManagedBlockOnly_HappyPath: extracts the BEGIN..END
// region including marker lines, leaving surrounding template
// content (which is the whole point of T4a) untouched.
func TestRenderManagedBlockOnly_HappyPath(t *testing.T) {
	body := []byte(
		"# BEGIN U-BOOT MANAGED BLOCK: init\n" +
			"name: demo\n" +
			"# END U-BOOT MANAGED BLOCK: init\n" +
			"\n" +
			"services: {}\n",
	)
	got, err := application.RenderManagedBlockOnlyForTest(body, "init")
	if err != nil {
		t.Fatalf("RenderManagedBlockOnly: %v", err)
	}
	want := "# BEGIN U-BOOT MANAGED BLOCK: init\n" +
		"name: demo\n" +
		"# END U-BOOT MANAGED BLOCK: init\n"
	if string(got) != want {
		t.Errorf("got\n%q\nwant\n%q", got, want)
	}
}

// TestEnsureComposeScaffold_AddsBothMissing pins the bare ensure-
// pass: empty input gets both add-on host maps appended.
func TestEnsureComposeScaffold_AddsBothMissing(t *testing.T) {
	out := application.EnsureComposeScaffoldForTest([]byte("name: demo\n"))
	got := string(out)
	if !strings.Contains(got, "\nservices: {}") {
		t.Errorf("missing services; got:\n%s", got)
	}
	if !strings.Contains(got, "\nvolumes: {}") {
		t.Errorf("missing volumes; got:\n%s", got)
	}
}

// TestEnsureComposeScaffold_NoOpWhenBothPresent: when both keys are
// already top-level and outside any managed block, the body is
// returned byte-identical.
func TestEnsureComposeScaffold_NoOpWhenBothPresent(t *testing.T) {
	in := []byte("name: demo\nservices: {}\nvolumes: {}\n")
	out := application.EnsureComposeScaffoldForTest(in)
	if !reflect.DeepEqual(out, in) {
		t.Errorf("expected byte-identical no-op; got:\n%q\nwant:\n%q", out, in)
	}
}

// TestEnsureComposeScaffold_IgnoresKeyInsideManagedBlock pins the
// Alt-M3 motivation: a `services: {}` *inside* the init block does
// NOT count as "top-level present" — ensure-pass still appends it
// outside.
func TestEnsureComposeScaffold_IgnoresKeyInsideManagedBlock(t *testing.T) {
	in := []byte("# BEGIN U-BOOT MANAGED BLOCK: init\n" +
		"services: {}\n" +
		"# END U-BOOT MANAGED BLOCK: init\n")
	out := application.EnsureComposeScaffoldForTest(in)
	got := string(out)
	// Original block left untouched (still contains `services: {}` inside)
	if !strings.Contains(got, "# BEGIN U-BOOT MANAGED BLOCK: init") {
		t.Errorf("init block missing; got:\n%s", got)
	}
	// New top-level keys appended after.
	endIdx := strings.Index(got, "# END U-BOOT MANAGED BLOCK: init")
	after := got[endIdx:]
	if !strings.Contains(after, "\nservices: {}") {
		t.Errorf("top-level services: missing outside block; got:\n%s", got)
	}
	if !strings.Contains(after, "\nvolumes: {}") {
		t.Errorf("top-level volumes: missing; got:\n%s", got)
	}
}
