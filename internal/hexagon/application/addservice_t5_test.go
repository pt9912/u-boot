package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/application/managedblock"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// --- M5-T5: LH-FA-ADD-005 state-machine end-to-end -----------------

// newInitAndAddServices builds an InitProjectService and an
// AddServiceService against the same in-memory FS so a test can run
// the full `u-boot init && u-boot add postgres && u-boot init --force`
// pipeline that the T4a slice promised.
func newInitAndAddServices(t *testing.T) (*application.InitProjectService, *application.AddServiceService, *fakeFS, *fakeYAML) {
	t.Helper()
	fs := newFakeFS()
	fs.markDirExists(addTestBaseDir)
	y := &fakeYAML{}
	g := &fakeGit{}
	initSvc := application.NewInitProjectService(fs, y, g, nil, nil, nil)
	addSvc := application.NewAddServiceService(fs, y, nil)
	return initSvc, addSvc, fs, y
}

// TestStateMachine_LH_FA_ADD_005_AllOutcomes is the central spec-
// pinned table for the LH-FA-ADD-005 state machine. Each case sets up
// one combination of u-boot.yaml / compose.yaml / .env.example and
// asserts the documented outcome (error sentinel or final state +
// Changed shape). Companion to the per-tranche unit tests in
// addservice_test.go and addservice_t4c_test.go — those pin the
// individual primitives; this one pins the user-visible contract.
func TestStateMachine_LH_FA_ADD_005_AllOutcomes(t *testing.T) {
	const (
		yamlUnregistered = "schemaVersion: 1\nproject:\n  name: demo\n"
		yamlActive       = "schemaVersion: 1\nproject:\n  name: demo\n" +
			"services:\n  postgres:\n    enabled: true\n"
		yamlDeactivated = "schemaVersion: 1\nproject:\n  name: demo\n" +
			"services:\n  postgres:\n    enabled: false\n"
		yamlEnabledUnset = "schemaVersion: 1\nproject:\n  name: demo\n" +
			"services:\n  postgres: {}\n"
	)

	cases := []struct {
		name        string
		ubootYAML   string
		composeBody string
		envBody     string
		wantErr     error // non-nil ⇒ Add must error with this sentinel
		wantPrior   domain.ServiceState
		wantChanged []string
	}{
		{
			name:        "unregistered → active (fresh register)",
			ubootYAML:   yamlUnregistered,
			wantPrior:   domain.ServiceStateUnregistered,
			wantChanged: []string{"u-boot.yaml", "compose.yaml", ".env.example"},
		},
		{
			name:        "active + complete artefacts → nil-error no-op",
			ubootYAML:   yamlActive,
			composeBody: composeBlockComplete("postgres"),
			envBody:     envBlockComplete("postgres"),
			wantPrior:   domain.ServiceStateActive,
			wantChanged: nil,
		},
		{
			name: "deactivated → active (reactivate)",
			ubootYAML:   yamlDeactivated,
			composeBody: composeBlockComplete("postgres"),
			envBody:     envBlockComplete("postgres"),
			wantPrior:   domain.ServiceStateDeactivated,
			// reactivate flips enabled + re-emits compose blocks;
			// env stays untouched (user values complete).
			wantChanged: []string{"u-boot.yaml", "compose.yaml"},
		},
		{
			name: "inconsistent-yaml → ErrServiceInconsistent",
			ubootYAML: yamlUnregistered,
			composeBody: "services:\n" +
				"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
				"  postgres: {}\n" +
				"  # END U-BOOT MANAGED BLOCK: service.postgres\n",
			wantErr: driving.ErrServiceInconsistent,
		},
		{
			name:      "malformed managed compose-block → ErrServiceInconsistent",
			ubootYAML: yamlActive,
			composeBody: "services:\n" +
				"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
				"  postgres:\n",
			wantErr: driving.ErrServiceInconsistent,
		},
		{
			name:        "inconsistent-block → active (compose rebuild only)",
			ubootYAML:   yamlActive,
			composeBody: "",
			envBody:     envBlockComplete("postgres"),
			wantPrior:   domain.ServiceStateInconsistentBlock,
			wantChanged: []string{"compose.yaml"},
		},
		{
			name:        "enabled-key missing → treated as deactivated → reactivate",
			ubootYAML:   yamlEnabledUnset,
			composeBody: composeBlockComplete("postgres"),
			envBody:     envBlockComplete("postgres"),
			wantPrior:   domain.ServiceStateEnabledUnset,
			wantChanged: []string{"u-boot.yaml", "compose.yaml"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, fs, _ := newAddService(t)
			seedUBootYAML(t, fs, tc.ubootYAML)
			if tc.composeBody != "" {
				seedCompose(t, fs, tc.composeBody)
			}
			if tc.envBody != "" {
				seedEnv(t, fs, tc.envBody)
			}
			writesBefore := len(fs.writtenPaths())

			resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
				BaseDir:     addTestBaseDir,
				ServiceName: postgresName(t),
			})

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("err = %v, want errors.Is(%v)", err, tc.wantErr)
				}
				if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
					t.Errorf("wrote %d files after error path; expected 0", delta)
				}
				return
			}

			if err != nil {
				t.Fatalf("Add: unexpected error %v", err)
			}
			if resp.PriorState != tc.wantPrior {
				t.Errorf("PriorState = %s, want %s",
					resp.PriorState.String(), tc.wantPrior.String())
			}
			if resp.State != domain.ServiceStateActive {
				t.Errorf("State = %s, want active", resp.State.String())
			}
			if !reflect.DeepEqual(resp.Changed, tc.wantChanged) {
				t.Errorf("Changed = %v, want %v", resp.Changed, tc.wantChanged)
			}
		})
	}
}

// TestStateMachine_ActiveWithMalformedVolume_Aborts pins the
// LH-FA-ADD-002 / -005 corner: even in Active state, a malformed
// volume.postgres block must abort with ErrServiceInconsistent. The
// Active branch goes through detectActiveArtifacts → LocateMarkedEntry,
// which returns the malformed-marker error.
func TestStateMachine_ActiveWithMalformedVolume_Aborts(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	malformed := strings.Replace(composeBlockComplete("postgres"),
		"  # END U-BOOT MANAGED BLOCK: volume.postgres\n", "", 1)
	seedCompose(t, fs, malformed)
	seedEnv(t, fs, envBlockComplete("postgres"))

	_, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if !errors.Is(err, driving.ErrServiceInconsistent) {
		t.Fatalf("err = %v, want ErrServiceInconsistent", err)
	}
}

// TestStateMachine_ActiveWithMalformedEnvBlock_Aborts pins the same
// LH-FA-ADD-002 promise on the env-block side.
func TestStateMachine_ActiveWithMalformedEnvBlock_Aborts(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	seedCompose(t, fs, composeBlockComplete("postgres"))
	// BEGIN without END.
	seedEnv(t, fs,
		"# BEGIN U-BOOT MANAGED BLOCK: service.postgres\n"+
			"POSTGRES_USER=a\n")

	_, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if !errors.Is(err, driving.ErrServiceInconsistent) {
		t.Fatalf("err = %v, want ErrServiceInconsistent", err)
	}
}

// TestStateMachine_DeactivatedWithMalformedServiceMarker_Aborts pins
// the spec promise that detectServiceState's malformed-marker
// pre-classification fires before the reactivate plan reaches
// PatchMappingEntryYAML.
func TestStateMachine_DeactivatedWithMalformedServiceMarker_Aborts(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: false\n")
	seedCompose(t, fs,
		"services:\n"+
			"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n"+
			"  postgres:\n")

	writesBefore := len(fs.writtenPaths())
	_, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if !errors.Is(err, driving.ErrServiceInconsistent) {
		t.Fatalf("err = %v, want ErrServiceInconsistent", err)
	}
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("wrote %d files despite malformed-marker abort", delta)
	}
}

// --- Round-trip scenarios -------------------------------------------

// TestRoundTrip_Idempotency_AddTwice pins the spec contract that a
// second `u-boot add postgres` after a successful first run is a
// true no-op (Changed=nil, no FS writes).
func TestRoundTrip_Idempotency_AddTwice(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs, "schemaVersion: 1\nproject:\n  name: demo\n")

	if _, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	}); err != nil {
		t.Fatalf("Add 1: %v", err)
	}
	writesBetween := len(fs.writtenPaths())

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add 2: %v", err)
	}
	if resp.Changed != nil {
		t.Errorf("second Add expected no-op, got Changed=%v", resp.Changed)
	}
	if resp.PriorState != domain.ServiceStateActive || resp.State != domain.ServiceStateActive {
		t.Errorf("second Add states = (%s, %s), want (active, active)",
			resp.PriorState.String(), resp.State.String())
	}
	if delta := len(fs.writtenPaths()) - writesBetween; delta != 0 {
		t.Errorf("second Add wrote %d files; expected 0", delta)
	}
}

// TestRoundTrip_InitAddInitForce_PreservesAddons pins the T4a
// scaffold decision: `u-boot init && u-boot add postgres && u-boot
// init --force` must keep both service.postgres and volume.postgres
// markers in the live compose.yaml. The whole point of the split-block
// scaffold was to make this case non-destructive.
func TestRoundTrip_InitAddInitForce_PreservesAddons(t *testing.T) {
	initSvc, addSvc, fs, _ := newInitAndAddServices(t)
	ctx := context.Background()

	if _, err := initSvc.Init(ctx, driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: addTestBaseDir,
		SkipGit: true,
	}); err != nil {
		t.Fatalf("Init 1: %v", err)
	}
	if _, err := addSvc.Add(ctx, driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	}); err != nil {
		t.Fatalf("Add: %v", err)
	}
	// Realistic re-init: --force on its own would reject .gitignore
	// (whole-file-managed, no marker block, LH-FA-INIT-005 §619);
	// users who actually re-init pair --force with --backup.
	if _, err := initSvc.Init(ctx, driving.InitProjectRequest{
		Name:    "demo",
		BaseDir: addTestBaseDir,
		Force:   true,
		Backup:  true,
		SkipGit: true,
	}); err != nil {
		t.Fatalf("Init 2 (--force --backup): %v", err)
	}

	body, err := fs.ReadFile(filepath.Join(addTestBaseDir, "compose.yaml"))
	if err != nil {
		t.Fatalf("ReadFile compose.yaml: %v", err)
	}
	for _, want := range []string{
		"BEGIN U-BOOT MANAGED BLOCK: service.postgres",
		"BEGIN U-BOOT MANAGED BLOCK: volume.postgres",
	} {
		if !strings.Contains(string(body), want) {
			t.Errorf("after init --force, marker %q missing; got:\n%s", want, body)
		}
	}
	// Exactly one init block, exactly one of each add-on marker.
	for _, marker := range []string{
		"BEGIN U-BOOT MANAGED BLOCK: init",
		"BEGIN U-BOOT MANAGED BLOCK: service.postgres",
		"BEGIN U-BOOT MANAGED BLOCK: volume.postgres",
	} {
		if c := strings.Count(string(body), marker); c != 1 {
			t.Errorf("expected exactly one %q, got %d", marker, c)
		}
	}
	// Final compose must still be parsable by managedblock.Find for
	// both add-on markers (T4c adapter-E2E rule re-applied here).
	for _, m := range []managedblock.Marker{
		{Style: managedblock.StyleHash, Name: "service.postgres"},
		{Style: managedblock.StyleHash, Name: "volume.postgres"},
	} {
		if _, _, err := managedblock.Find(body, m); err != nil {
			t.Errorf("managedblock.Find(%s) failed after round-trip: %v", m.Name, err)
		}
	}
}

// TestRoundTrip_MissingComposeRecovery covers the
// `enabled: true` + missing compose.yaml path: state-detection
// reports InconsistentBlock, executeAdd bootstraps a fresh compose
// from the project name in u-boot.yaml, then patches the service and
// volume blocks into it.
func TestRoundTrip_MissingComposeRecovery(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	// No compose.yaml at all.

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if resp.State != domain.ServiceStateActive {
		t.Errorf("State = %s, want active", resp.State.String())
	}
	body, err := fs.ReadFile(filepath.Join(addTestBaseDir, "compose.yaml"))
	if err != nil {
		t.Fatalf("compose.yaml not written: %v", err)
	}
	for _, want := range []string{
		"name: demo",
		"BEGIN U-BOOT MANAGED BLOCK: service.postgres",
		"BEGIN U-BOOT MANAGED BLOCK: volume.postgres",
	} {
		if !strings.Contains(string(body), want) {
			t.Errorf("missing %q in recovered compose; got:\n%s", want, body)
		}
	}
}

// TestRoundTrip_EnvIdempotency_KeepsUserLines pins that two `u-boot
// add postgres` runs leave the user-managed sections of
// .env.example untouched and produce exactly one
// service.postgres-managed block.
func TestRoundTrip_EnvIdempotency_KeepsUserLines(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs, "schemaVersion: 1\nproject:\n  name: demo\n")
	// Seed an env with user lines.
	seedEnv(t, fs, "# user-owned section below\nMY_API_KEY=secret\n")

	ctx := context.Background()
	if _, err := svc.Add(ctx, driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	}); err != nil {
		t.Fatalf("Add 1: %v", err)
	}
	body1, _ := fs.ReadFile(filepath.Join(addTestBaseDir, ".env.example"))

	if _, err := svc.Add(ctx, driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	}); err != nil {
		t.Fatalf("Add 2: %v", err)
	}
	body2, _ := fs.ReadFile(filepath.Join(addTestBaseDir, ".env.example"))

	if string(body1) != string(body2) {
		t.Errorf("second add changed env file:\nbody1:\n%s\nbody2:\n%s", body1, body2)
	}
	if !strings.Contains(string(body1), "# user-owned section below") {
		t.Errorf("user comment lost; got:\n%s", body1)
	}
	if !strings.Contains(string(body1), "MY_API_KEY=secret") {
		t.Errorf("user key lost; got:\n%s", body1)
	}
	if c := strings.Count(string(body1), "BEGIN U-BOOT MANAGED BLOCK: service.postgres"); c != 1 {
		t.Errorf("expected exactly one managed env block, got %d; body:\n%s", c, body1)
	}
}
