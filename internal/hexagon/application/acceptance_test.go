// Package application_test acceptance pins for the MVP
// LH-AK-* criteria from spec/lastenheft.md §9.
//
// Each test is named `TestLHAK00X_…` to make the spec-anchor
// obvious to grep. LH-AK-002 (postgres flow) and LH-AK-007
// (changelog generator) live in their own home — the Docker-tagged
// e2e package (`internal/e2e/`) and `generate_test.go`
// respectively — because they need real Docker (LH-AK-002) or
// already had a focused e2e test from M7-T4 (LH-AK-007). This
// file collects the remaining MVP-pin tests that run in the plain
// application-layer test harness (no build tags, no real Docker).
//
// MVP-Closure-T2 (slice-mvp-closure.md).
package application_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TestLHAK001_InitFlow_DoctorClean pins `LH-AK-001` (spec §2281):
//
//	mkdir demo && cd demo && u-boot init && u-boot doctor
//
// must produce a project structure plus a doctor report with **no**
// `error`-severity diagnostics. The application-layer harness
// substitutes the mkdir+cd step (the fakeFS pre-registers the
// directory) and stubs the Docker probe at LH-FA-DIAG-002-OK
// versions so the doctor's docker checks pass without a real
// daemon. The spec's "vorhandene Dateien wurden nicht ungewollt
// überschrieben" clause is satisfied vacuously (fresh tempDir);
// the negative-overwrite path is independently pinned by the
// existing `TestInit_ExistingProjectRejected`.
func TestLHAK001_InitFlow_DoctorClean(t *testing.T) {
	fs := newFakeFS()
	fs.markDirExists(testBaseDir)
	y := &fakeYAML{}
	git := &fakeGit{}
	docker := &fakeDockerProbe{
		version:        "24.0.0",
		composeVersion: "v2.20.0",
	}

	initSvc := application.NewInitProjectService(fs, y, git, nil, nil, nil)
	doctorSvc := application.NewDoctorService(fs, y, git, docker, nil, nil)

	if _, err := initSvc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir: testBaseDir,
		Name:    "lhak001demo",
		SkipGit: true,
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	resp, err := doctorSvc.Check(context.Background(), driving.DoctorRequest{
		BaseDir: testBaseDir,
	})
	if err != nil {
		t.Fatalf("doctor: %v", err)
	}

	if resp.Report.HasErrors() {
		var errorIDs []string
		for _, d := range resp.Report.Items {
			if d.Severity == domain.SeverityError {
				errorIDs = append(errorIDs, d.ID+": "+d.Message)
			}
		}
		t.Errorf("doctor reported errors after fresh init (LH-AK-001 spec: no error):\n  %s",
			strings.Join(errorIDs, "\n  "))
	}
}

// TestLHAK006_DoubleAddPostgres_NoDuplicate pins `LH-AK-006`
// (spec §2387):
//
//	u-boot add postgres
//	u-boot add postgres
//
// Erwartetes Ergebnis: PostgreSQL ist genau einmal in der
// Konfiguration vorhanden + verständliche Meldung.
//
// In application-layer terms: the second Add call returns
// `(PriorState=Active, State=Active, Changed=nil)` and produces
// zero new WriteFile invocations. The "verständliche Meldung" is
// the CLI's `printAddSummary("Service \"postgres\" is already
// active; no changes.")` — covered separately by
// `cli_test.go::TestExecute_Add_AlreadyActiveNoOp`.
//
// Mirrors the LH-AK-007 e2e shape from M7-T4
// (`TestGenerateChangelog_LHAK007_FlowEndToEnd`): direct service
// calls, no Docker, no CLI layer.
func TestLHAK006_DoubleAddPostgres_NoDuplicate(t *testing.T) {
	fs := newFakeFS()
	fs.markDirExists(testBaseDir)
	y := &fakeYAML{}
	git := &fakeGit{}

	initSvc := application.NewInitProjectService(fs, y, git, nil, nil, nil)
	addSvc := application.NewAddServiceService(fs, y, nil, nil)

	if _, err := initSvc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir: testBaseDir,
		Name:    "lhak006demo",
		SkipGit: true,
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	postgres, err := domain.NewServiceName("postgres")
	if err != nil {
		t.Fatalf("NewServiceName: %v", err)
	}

	// First add: registers postgres.
	resp1, err := addSvc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     testBaseDir,
		ServiceName: postgres,
	})
	if err != nil {
		t.Fatalf("first add: %v", err)
	}
	if resp1.State != domain.ServiceStateActive {
		t.Errorf("first add State = %v, want Active", resp1.State)
	}

	writesAfterFirst := len(fs.writtenPaths())

	// Second add: idempotent no-op.
	resp2, err := addSvc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     testBaseDir,
		ServiceName: postgres,
	})
	if err != nil {
		t.Fatalf("second add: %v", err)
	}
	if resp2.PriorState != domain.ServiceStateActive || resp2.State != domain.ServiceStateActive {
		t.Errorf("second add states = (%v, %v), want (Active, Active)",
			resp2.PriorState, resp2.State)
	}
	if resp2.Changed != nil {
		t.Errorf("second add Changed = %v, want nil (LH-AK-006: no duplicate)", resp2.Changed)
	}
	if delta := len(fs.writtenPaths()) - writesAfterFirst; delta != 0 {
		t.Errorf("second add produced %d WriteFile call(s), want 0; writes = %v",
			delta, fs.writtenPaths())
	}

	// "Genau einmal in der Konfiguration": u-boot.yaml has the
	// postgres key once. Read the rendered file and grep — the
	// fakeYAML round-trip preserves the source verbatim, so a
	// substring count is deterministic.
	body, err := fs.ReadFile(testBaseDir + "/u-boot.yaml")
	if err != nil {
		t.Fatalf("read u-boot.yaml: %v", err)
	}
	if got := strings.Count(string(body), "postgres:"); got != 1 {
		t.Errorf("u-boot.yaml contains %d `postgres:` entries, want 1; body:\n%s",
			got, body)
	}
}

// TestLHFADEV003_CatalogueActivation pins `LH-FA-DEV-003` happy
// path (spec/lastenheft.md:692-721 + slice-v1-devcontainer-features
// AK „Spec-Pin"):
//
//	u-boot init --devcontainer
//	u-boot config set devcontainer.features.node.enabled true
//	u-boot generate devcontainer
//
// must produce a `.devcontainer/devcontainer.json` whose
// `features:` block contains the key
// `ghcr.io/devcontainers/features/node:1` (catalogue lookup + T3
// renderer projection). No Allowlist needed because `node` is a
// built-in catalogue feature (Spec §711).
func TestLHFADEV003_CatalogueActivation(t *testing.T) {
	fs := newFakeFS()
	fs.markDirExists(testBaseDir)
	y := &fakeYAML{}
	git := &fakeGit{}

	initSvc := application.NewInitProjectService(fs, y, git, nil, nil, nil)
	configSvc := application.NewConfigService(fs, y, nil)
	generateSvc := application.NewGenerateService(fs, y, nil)

	// init --devcontainer
	if _, err := initSvc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir:      testBaseDir,
		Name:         "demo",
		SkipGit:      true,
		Devcontainer: true,
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	// config set devcontainer.features.node.enabled true
	path, err := domain.NewConfigPath("devcontainer.features.node.enabled")
	if err != nil {
		t.Fatalf("NewConfigPath: %v", err)
	}
	if _, err := configSvc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: testBaseDir,
		Path:    path,
		Value:   "true",
	}); err != nil {
		t.Fatalf("config set: %v", err)
	}

	// generate devcontainer
	if _, err := generateSvc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:  testBaseDir,
		Artifact: domain.ArtifactDevcontainer,
	}); err != nil {
		t.Fatalf("generate: %v", err)
	}

	// devcontainer.json carries the canonical OCI-ref key.
	body, err := fs.ReadFile(testBaseDir + "/.devcontainer/devcontainer.json")
	if err != nil {
		t.Fatalf("read devcontainer.json: %v", err)
	}
	wantKey := `"ghcr.io/devcontainers/features/node:1": {}`
	if !strings.Contains(string(body), wantKey) {
		t.Errorf("devcontainer.json missing feature key %q; body:\n%s",
			wantKey, body)
	}
}

// TestLHFADEV003_AllowlistEnforcement pins `LH-FA-DEV-003` negative
// path (spec/lastenheft.md:720, LH-NFA-SEC-004): an attempt to
// register a `features.<name>.source` URL that is not in
// `featureSources.allow` fails with [driving.ErrConfigValueInvalid]
// (LH-FA-CLI-006 exit-code 10). The seed via
// `--allow-external-feature-sources` at init time then unblocks the
// same call. `--yes` is not exercised (it doesn't apply to the
// non-interactive Set path) but the sentinel-chain contract is the
// LH-NFA-SEC-004 surface.
func TestLHFADEV003_AllowlistEnforcement(t *testing.T) {
	fs := newFakeFS()
	fs.markDirExists(testBaseDir)
	y := &fakeYAML{}
	git := &fakeGit{}

	initSvc := application.NewInitProjectService(fs, y, git, nil, nil, nil)
	configSvc := application.NewConfigService(fs, y, nil)

	// init --devcontainer (no Allowlist seed yet).
	if _, err := initSvc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir:      testBaseDir,
		Name:         "demo",
		SkipGit:      true,
		Devcontainer: true,
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	// Negative path: source override without Allowlist entry → error.
	sourcePath, err := domain.NewConfigPath("devcontainer.features.custom.source")
	if err != nil {
		t.Fatalf("NewConfigPath: %v", err)
	}
	externalURL := "https://ghcr.io/orgX/features/custom-rust"
	_, err = configSvc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: testBaseDir,
		Path:    sourcePath,
		Value:   externalURL,
	})
	if err == nil {
		t.Fatalf("Set: expected error for source not in allowlist, got nil")
	}
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Errorf("err = %v, want wrap of ErrConfigValueInvalid (LH-FA-CLI-006 exit-10)", err)
	}
	if !strings.Contains(err.Error(), "LH-FA-DEV-003") {
		t.Errorf("err message %q does not name LH-FA-DEV-003", err.Error())
	}
	if !strings.Contains(err.Error(), "LH-NFA-SEC-004") {
		t.Errorf("err message %q does not name LH-NFA-SEC-004", err.Error())
	}

	// Positive path after seeding the Allowlist via the Spec §717-
	// `config set devcontainer.featureSources.allow` route — same
	// shape the CLI flag uses internally.
	allowPath, err := domain.NewConfigPath("devcontainer.featureSources.allow")
	if err != nil {
		t.Fatalf("NewConfigPath(allow): %v", err)
	}
	if _, err := configSvc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: testBaseDir,
		Path:    allowPath,
		Value:   externalURL,
	}); err != nil {
		t.Fatalf("seed allowlist: %v", err)
	}
	// Now the same source override succeeds.
	if _, err := configSvc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: testBaseDir,
		Path:    sourcePath,
		Value:   externalURL,
	}); err != nil {
		t.Errorf("source override after Allowlist seed: %v", err)
	}
}

