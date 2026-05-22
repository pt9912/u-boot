package application_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"path/filepath"
	"reflect"
	"regexp"
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
	return application.NewInitProjectService(fs, y, g, io.Discard), fs, y, g
}

// newServiceWithProgress is newService plus a buffer for the
// LH-FA-INIT-005 §609 affected-paths summary. Tests that assert on
// summary text use this constructor.
func newServiceWithProgress(t *testing.T) (*application.InitProjectService, *fakeFS, *fakeYAML, *fakeGit, *bytes.Buffer) {
	t.Helper()
	fs := newFakeFS()
	fs.markDirExists(testBaseDir)
	y := &fakeYAML{}
	g := &fakeGit{}
	progress := &bytes.Buffer{}
	return application.NewInitProjectService(fs, y, g, progress), fs, y, g, progress
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
	// Why: LH-FA-INIT-005 §609 / LH-FA-CLI-005A §262 — affected paths
	// must be reported BEFORE the write. We assert that the progress
	// writer received the summary lines.
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

	out := progress.String()
	if !strings.Contains(out, "Affected files in") {
		t.Errorf("summary header missing: %q", out)
	}
	if !strings.Contains(out, "compose.yaml") {
		t.Errorf("compose.yaml not mentioned: %q", out)
	}
	if !strings.Contains(out, "replace managed block") {
		t.Errorf("action label missing: %q", out)
	}
}

func TestInit_Summary_QuietOnFreshInit(t *testing.T) {
	// Why: defensive — fresh init must not emit a summary; nothing
	// is being overwritten.
	svc, _, _, _, progress := newServiceWithProgress(t)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: testBaseDir,
		SkipGit: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if progress.Len() != 0 {
		t.Errorf("expected empty progress on fresh init, got %q", progress.String())
	}
}

func TestInit_Summary_WithBackupMarker(t *testing.T) {
	// Why: pin that the summary line for a --backup overwrite carries
	// the "(with backup)" marker so the user knows their original is
	// safe.
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
	if !strings.Contains(progress.String(), "(with backup)") {
		t.Errorf("backup marker missing in summary: %q", progress.String())
	}
}

func TestInit_NilProgress_TolerantToDiscard(t *testing.T) {
	// Why: constructor must accept nil progress without panicking;
	// it falls back to io.Discard.
	fs := newFakeFS()
	fs.markDirExists(testBaseDir)
	svc := application.NewInitProjectService(fs, &fakeYAML{}, &fakeGit{}, nil)

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
