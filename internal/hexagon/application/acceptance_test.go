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
	doctorSvc := application.NewDoctorService(fs, y, git, docker, nil)

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
	addSvc := application.NewAddServiceService(fs, y, nil)

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
