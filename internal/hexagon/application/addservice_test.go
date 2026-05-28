package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const addTestBaseDir = "/tmp/u-boot-add-test/demo"

// newAddService constructs a service against an in-memory FS plus a
// real yaml.v3-backed fake codec. addTestBaseDir is pre-registered
// as an existing directory because the production service refuses an
// uninitialized BaseDir at the Exists check on u-boot.yaml.
func newAddService(t *testing.T) (*application.AddServiceService, *fakeFS, *fakeYAML) {
	t.Helper()
	fs := newFakeFS()
	fs.markDirExists(addTestBaseDir)
	y := &fakeYAML{}
	svc := application.NewAddServiceService(fs, y, nil)
	return svc, fs, y
}

// seedUBootYAML writes a u-boot.yaml under addTestBaseDir.
func seedUBootYAML(t *testing.T, fs *fakeFS, body string) {
	t.Helper()
	if err := fs.WriteFile(filepath.Join(addTestBaseDir, "u-boot.yaml"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
}

// seedCompose writes a compose.yaml under addTestBaseDir.
func seedCompose(t *testing.T, fs *fakeFS, body string) {
	t.Helper()
	if err := fs.WriteFile(filepath.Join(addTestBaseDir, "compose.yaml"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed compose.yaml: %v", err)
	}
}

// postgresName is shorthand for the validated postgres ServiceName
// used in every fixture.
func postgresName(t *testing.T) domain.ServiceName {
	t.Helper()
	n, err := domain.NewServiceName("postgres")
	if err != nil {
		t.Fatalf("NewServiceName(postgres): %v", err)
	}
	return n
}

// composeBlock returns a minimal managed compose-block body that
// matches the production marker layout (`# BEGIN/END U-BOOT MANAGED
// BLOCK: service.<svc>`).
func composeBlock(svc string) string {
	return "services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service." + svc + "\n" +
		"  " + svc + ":\n" +
		"    image: example\n" +
		"  # END U-BOOT MANAGED BLOCK: service." + svc + "\n"
}

// TestDetectServiceState_AllSixStates pins the LH-FA-ADD-005
// classification table by exercising each combination of YAML-entry
// presence, Enabled-pointer state and compose-block presence.
func TestDetectServiceState_AllSixStates(t *testing.T) {
	cases := []struct {
		name        string
		ubootYAML   string
		composeBody string
		wantState   domain.ServiceState
	}{
		{
			name:        "unregistered: no entry, no block",
			ubootYAML:   "schemaVersion: 1\nproject:\n  name: demo\n",
			composeBody: "",
			wantState:   domain.ServiceStateUnregistered,
		},
		{
			name:        "inconsistent-yaml: no entry, block present",
			ubootYAML:   "schemaVersion: 1\nproject:\n  name: demo\n",
			composeBody: composeBlock("postgres"),
			wantState:   domain.ServiceStateInconsistentYAML,
		},
		{
			name: "active: enabled=true, block present",
			ubootYAML: "schemaVersion: 1\nproject:\n  name: demo\n" +
				"services:\n  postgres:\n    enabled: true\n",
			composeBody: composeBlock("postgres"),
			wantState:   domain.ServiceStateActive,
		},
		{
			name: "inconsistent-block: enabled=true, block missing",
			ubootYAML: "schemaVersion: 1\nproject:\n  name: demo\n" +
				"services:\n  postgres:\n    enabled: true\n",
			composeBody: "",
			wantState:   domain.ServiceStateInconsistentBlock,
		},
		{
			name: "deactivated: enabled=false (block irrelevant)",
			ubootYAML: "schemaVersion: 1\nproject:\n  name: demo\n" +
				"services:\n  postgres:\n    enabled: false\n",
			composeBody: composeBlock("postgres"),
			wantState:   domain.ServiceStateDeactivated,
		},
		{
			name: "enabled-unset: services.postgres present but no enabled key",
			ubootYAML: "schemaVersion: 1\nproject:\n  name: demo\n" +
				"services:\n  postgres: {}\n",
			composeBody: "",
			wantState:   domain.ServiceStateEnabledUnset,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, fs, _ := newAddService(t)
			seedUBootYAML(t, fs, tc.ubootYAML)
			if tc.composeBody != "" {
				seedCompose(t, fs, tc.composeBody)
			}
			got, err := svc.DetectServiceStateForTest(addTestBaseDir, postgresName(t))
			if err != nil {
				t.Fatalf("detectServiceState: unexpected error %v", err)
			}
			if got != tc.wantState {
				t.Errorf("state = %s, want %s", got.String(), tc.wantState.String())
			}
		})
	}
}

// TestDetectServiceState_MalformedComposeBlock_Aborts pins the T3
// pre-classification abort: a malformed `service.postgres` block
// (BEGIN without END) returns ErrServiceInconsistent regardless of
// the YAML side. The wrapped error must satisfy errors.Is so the CLI
// adapter (T6) can map to exit code 10.
func TestDetectServiceState_MalformedComposeBlock_Aborts(t *testing.T) {
	cases := []struct {
		name      string
		ubootYAML string
	}{
		{
			name: "malformed block + enabled=true",
			ubootYAML: "schemaVersion: 1\nproject:\n  name: demo\n" +
				"services:\n  postgres:\n    enabled: true\n",
		},
		{
			name: "malformed block + enabled=false (must not pass as deactivated)",
			ubootYAML: "schemaVersion: 1\nproject:\n  name: demo\n" +
				"services:\n  postgres:\n    enabled: false\n",
		},
		{
			name: "malformed block + no services entry (must not pass as inconsistent-yaml)",
			ubootYAML: "schemaVersion: 1\nproject:\n  name: demo\n",
		},
	}

	malformed := "services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  postgres:\n" +
		"    image: example\n" +
		"# (missing END marker)\n"

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, fs, _ := newAddService(t)
			seedUBootYAML(t, fs, tc.ubootYAML)
			seedCompose(t, fs, malformed)
			_, err := svc.DetectServiceStateForTest(addTestBaseDir, postgresName(t))
			if !errors.Is(err, driving.ErrServiceInconsistent) {
				t.Fatalf("err = %v, want ErrServiceInconsistent", err)
			}
		})
	}
}

// TestAdd_ActiveIsNilErrorNoOp pins the LH-FA-ADD-005 idempotent
// core-state behaviour. T3 does not yet apply the T4 artefact
// repair check, so a present enabled:true + compose-block pair is
// the no-op end-state.
func TestAdd_ActiveIsNilErrorNoOp(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	seedCompose(t, fs, composeBlock("postgres"))

	// Reset the writes log so any side-effect from this call is
	// visible against an empty baseline.
	writesBefore := len(fs.writtenPaths())

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add: unexpected error %v", err)
	}
	if resp.PriorState != domain.ServiceStateActive || resp.State != domain.ServiceStateActive {
		t.Errorf("states = (%s, %s), want (active, active)",
			resp.PriorState.String(), resp.State.String())
	}
	if resp.Changed != nil {
		t.Errorf("Changed = %v, want nil (no-op)", resp.Changed)
	}
	if got := len(fs.writtenPaths()); got != writesBefore {
		t.Errorf("WriteFile count = %d, want %d (no-op must not write)",
			got, writesBefore)
	}
}

// TestAdd_InconsistentYAML_ReturnsSentinel pins Spec §895: a managed
// compose-block without a YAML anchor aborts with a repair hint
// instead of silently re-creating the anchor.
func TestAdd_InconsistentYAML_ReturnsSentinel(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs, "schemaVersion: 1\nproject:\n  name: demo\n")
	seedCompose(t, fs, composeBlock("postgres"))

	_, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if !errors.Is(err, driving.ErrServiceInconsistent) {
		t.Fatalf("err = %v, want ErrServiceInconsistent", err)
	}
}

// TestAdd_UnsupportedService_ChecksCatalogueBeforeDisk pins that the
// catalogue check runs before any FS read — a typo in the service
// name fails fast without depending on project state.
func TestAdd_UnsupportedService_ChecksCatalogueBeforeDisk(t *testing.T) {
	svc, fs, _ := newAddService(t)
	// No seeded files at all — if the catalogue check were after
	// the FS read, the test would hit ErrProjectNotInitialized
	// instead.
	redis, err := domain.NewServiceName("redis")
	if err != nil {
		t.Fatalf("NewServiceName(redis): %v", err)
	}
	_, addErr := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: redis,
	})
	if !errors.Is(addErr, driving.ErrServiceUnsupported) {
		t.Fatalf("err = %v, want ErrServiceUnsupported", addErr)
	}
	if got := fs.readFileCallCount(filepath.Join(addTestBaseDir, "u-boot.yaml")); got != 0 {
		t.Errorf("ReadFile(u-boot.yaml) called %d times, want 0 (catalogue must short-circuit)", got)
	}
}

// TestAdd_ProjectNotInitialized_ReturnsSentinel pins LH-FA-ADD-001:
// `u-boot add` only works in an initialized project.
func TestAdd_ProjectNotInitialized_ReturnsSentinel(t *testing.T) {
	svc, _, _ := newAddService(t)
	// No u-boot.yaml seeded.
	_, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if !errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Fatalf("err = %v, want ErrProjectNotInitialized", err)
	}
}

// TestAdd_UnparsableUBootYAML_MapsToNotInitialized pins the T3
// decision that an unparsable u-boot.yaml is treated as "project
// not initialized" (the schema we expect is broken, so we cannot
// trust the project state). The error wraps the parse detail for
// diagnostics.
func TestAdd_UnparsableUBootYAML_MapsToNotInitialized(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs, ":\n  not yaml")
	_, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if !errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Fatalf("err = %v, want ErrProjectNotInitialized", err)
	}
}

// TestAdd_EmptyBaseDirRejects mirrors the InitProjectService policy
// for a missing required input — a non-sentinel error so CLI
// surfaces it as a clear, generic validation failure (LH-FA-CLI-006).
func TestAdd_EmptyBaseDirRejects(t *testing.T) {
	svc, _, _ := newAddService(t)
	_, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     "",
		ServiceName: postgresName(t),
	})
	if err == nil {
		t.Fatalf("err = nil, want validation error")
	}
	// Must NOT be a fachlicher Sentinel (no false project state).
	for _, sentinel := range []error{
		driving.ErrProjectNotInitialized,
		driving.ErrServiceInconsistent,
		driving.ErrServiceUnsupported,
	} {
		if errors.Is(err, sentinel) {
			t.Errorf("err = %v, must not wrap %v", err, sentinel)
		}
	}
}

// TestAdd_UBootYAMLReadFailure_IsTechnical pins the T3 contract that
// a Permission/I/O failure on u-boot.yaml ReadFile is a technical
// error — it must not surface as ErrProjectNotInitialized because
// the file exists; we just cannot read it.
func TestAdd_UBootYAMLReadFailure_IsTechnical(t *testing.T) {
	svc, fs, _ := newAddService(t)
	yamlPath := filepath.Join(addTestBaseDir, "u-boot.yaml")
	seedUBootYAML(t, fs, "schemaVersion: 1\nproject:\n  name: demo\n")
	fs.failReadOn = yamlPath
	fs.failReadErr = errors.New("permission denied")

	_, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err == nil {
		t.Fatalf("err = nil, want technical I/O error")
	}
	if errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Errorf("err = %v, must not wrap ErrProjectNotInitialized "+
			"(permission failure is technical)", err)
	}
}

// TestAdd_ComposeYAMLReadFailure_IsTechnical pins the symmetric
// contract for compose.yaml: an I/O failure must not silently fall
// through as block-absent (which would mis-classify Active services
// as InconsistentBlock) and must not be reported as
// ErrServiceInconsistent (the block could well be present; we just
// cannot tell).
func TestAdd_ComposeYAMLReadFailure_IsTechnical(t *testing.T) {
	svc, fs, _ := newAddService(t)
	composePath := filepath.Join(addTestBaseDir, "compose.yaml")
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	seedCompose(t, fs, composeBlock("postgres"))
	fs.failReadOn = composePath
	fs.failReadErr = errors.New("permission denied")

	_, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err == nil {
		t.Fatalf("err = nil, want technical I/O error")
	}
	if errors.Is(err, driving.ErrServiceInconsistent) {
		t.Errorf("err = %v, must not wrap ErrServiceInconsistent", err)
	}
}

// TestAdd_ExistsFailure_IsTechnical pins the same contract for an
// Exists check: a Permission failure on the Exists probe of either
// file must surface as a non-sentinel error.
func TestAdd_ExistsFailure_IsTechnical(t *testing.T) {
	svc, fs, _ := newAddService(t)
	yamlPath := filepath.Join(addTestBaseDir, "u-boot.yaml")
	seedUBootYAML(t, fs, "schemaVersion: 1\nproject:\n  name: demo\n")
	fs.failExistsOn = yamlPath
	fs.failExistsErr = errors.New("permission denied")

	_, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err == nil {
		t.Fatalf("err = nil, want technical I/O error")
	}
	for _, sentinel := range []error{
		driving.ErrProjectNotInitialized,
		driving.ErrServiceInconsistent,
	} {
		if errors.Is(err, sentinel) {
			t.Errorf("err = %v, must not wrap %v (Exists failure is technical)",
				err, sentinel)
		}
	}
}

// TestAdd_MutatingStates_HitExecuteStub pins the T3↔T4 contract: for
// each of the four mutating LH-FA-ADD-005 states the dispatch
// reaches planAdd → executeAdd, and the stub returns its
// "not yet implemented (M5-T4)" message. No FS write happens before
// the stub fires.
func TestAdd_MutatingStates_HitExecuteStub(t *testing.T) {
	cases := []struct {
		name        string
		ubootYAML   string
		composeBody string
	}{
		{
			name:        "unregistered",
			ubootYAML:   "schemaVersion: 1\nproject:\n  name: demo\n",
			composeBody: "",
		},
		{
			name: "deactivated",
			ubootYAML: "schemaVersion: 1\nproject:\n  name: demo\n" +
				"services:\n  postgres:\n    enabled: false\n",
			composeBody: "",
		},
		{
			name: "enabled-unset",
			ubootYAML: "schemaVersion: 1\nproject:\n  name: demo\n" +
				"services:\n  postgres: {}\n",
			composeBody: "",
		},
		{
			name: "inconsistent-block",
			ubootYAML: "schemaVersion: 1\nproject:\n  name: demo\n" +
				"services:\n  postgres:\n    enabled: true\n",
			composeBody: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, fs, _ := newAddService(t)
			seedUBootYAML(t, fs, tc.ubootYAML)
			if tc.composeBody != "" {
				seedCompose(t, fs, tc.composeBody)
			}
			writesBefore := len(fs.writtenPaths())
			_, err := svc.Add(context.Background(), driving.AddServiceRequest{
				BaseDir:     addTestBaseDir,
				ServiceName: postgresName(t),
			})
			if err == nil {
				t.Fatalf("err = nil, want execute-stub error")
			}
			// Stub message is the contract; tests intentionally
			// pin the M5-T4 marker so accidentally implementing
			// execute in T3 surfaces here.
			if !strings.Contains(err.Error(), "not yet implemented (M5-T4)") {
				t.Errorf("err = %v, want execute-stub marker", err)
			}
			// Sentinels reserved for fachliche Fehler must not
			// leak through the stub path.
			for _, sentinel := range []error{
				driving.ErrProjectNotInitialized,
				driving.ErrServiceInconsistent,
				driving.ErrServiceUnsupported,
			} {
				if errors.Is(err, sentinel) {
					t.Errorf("err = %v, must not wrap %v (stub is generic)",
						err, sentinel)
				}
			}
			// The stub returns before any write; pre-write
			// validation hasn't side-effects either.
			if got := len(fs.writtenPaths()) - writesBefore; got != 0 {
				t.Errorf("WriteFile count delta = %d, want 0 (no write before T4)", got)
			}
		})
	}
}

// TestPlanAdd_MutatingStates_AssignsExpectedAction pins the
// state→action mapping at the planAdd boundary so T4 can rely on
// the action enum without re-deriving the rules.
func TestPlanAdd_MutatingStates_AssignsExpectedAction(t *testing.T) {
	cases := []struct {
		state      domain.ServiceState
		wantAction string
	}{
		{state: domain.ServiceStateUnregistered, wantAction: "register"},
		{state: domain.ServiceStateDeactivated, wantAction: "reactivate"},
		{state: domain.ServiceStateEnabledUnset, wantAction: "reactivate"},
		{state: domain.ServiceStateInconsistentBlock, wantAction: "rebuild-block"},
	}
	svc, _, _ := newAddService(t)
	name := postgresName(t)
	for _, tc := range cases {
		t.Run(tc.state.String(), func(t *testing.T) {
			plan, err := svc.PlanAddForTest(name, tc.state)
			if err != nil {
				t.Fatalf("planAdd: %v", err)
			}
			if plan.Service != name {
				t.Errorf("plan.Service = %q, want %q", plan.Service, name)
			}
			if plan.PriorState != tc.state {
				t.Errorf("plan.PriorState = %s, want %s",
					plan.PriorState.String(), tc.state.String())
			}
			if plan.Action != tc.wantAction {
				t.Errorf("plan.Action = %s, want %s", plan.Action, tc.wantAction)
			}
		})
	}
}

// TestPlanAdd_NonMutatingStates_ReturnError pins that planAdd
// refuses to plan for Active or InconsistentYAML — Add()'s dispatch
// must handle those before reaching planAdd.
func TestPlanAdd_NonMutatingStates_ReturnError(t *testing.T) {
	svc, _, _ := newAddService(t)
	name := postgresName(t)
	for _, state := range []domain.ServiceState{
		domain.ServiceStateActive,
		domain.ServiceStateInconsistentYAML,
	} {
		t.Run(state.String(), func(t *testing.T) {
			_, err := svc.PlanAddForTest(name, state)
			if err == nil {
				t.Fatalf("planAdd(%s): want error, got nil", state.String())
			}
		})
	}
}

