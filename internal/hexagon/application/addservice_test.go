package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
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

// composeBlock returns a minimal managed compose-block body — just
// the structural marker plus a stub image. Used by the T3 state-
// detection tests where only marker presence matters; T4c
// content-checks would flag this body as stale (no environment /
// ports / healthcheck), so completeness-sensitive tests use
// composeBlockComplete instead.
func composeBlock(svc string) string {
	return "services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service." + svc + "\n" +
		"  " + svc + ":\n" +
		"    image: example\n" +
		"  # END U-BOOT MANAGED BLOCK: service." + svc + "\n"
}

// composeBlockComplete returns a full LH-FA-ADD-002 / LH-AK-002
// compose body for the given service: services.<svc> with all
// required fields (image, environment with the three POSTGRES_*
// keys, volumes referencing <svc>-data, ports, healthcheck) plus
// the matching volumes.<svc>-data block. Used by every T4c test
// that expects the Active state to be a true no-op.
func composeBlockComplete(svc string) string {
	return "services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service." + svc + "\n" +
		"  " + svc + ":\n" +
		"    image: " + svc + ":16-alpine\n" +
		"    environment:\n" +
		"      POSTGRES_USER: " + svc + "\n" +
		"      POSTGRES_PASSWORD: changeme\n" +
		"      POSTGRES_DB: " + svc + "\n" +
		"    volumes:\n" +
		"      - " + svc + "-data:/var/lib/postgresql/data\n" +
		"    ports:\n" +
		"      - \"5432:5432\"\n" +
		"    healthcheck:\n" +
		"      test: [\"CMD\", \"pg_isready\"]\n" +
		"  # END U-BOOT MANAGED BLOCK: service." + svc + "\n" +
		"\n" +
		"volumes:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: volume." + svc + "\n" +
		"  " + svc + "-data: {}\n" +
		"  # END U-BOOT MANAGED BLOCK: volume." + svc + "\n"
}

// envBlockComplete returns a .env.example body that contains a
// well-formed managed env block with all three POSTGRES_* keys
// present.
func envBlockComplete(svc string) string {
	return "# BEGIN U-BOOT MANAGED BLOCK: service." + svc + "\n" +
		"POSTGRES_USER=" + svc + "\n" +
		"POSTGRES_PASSWORD=changeme\n" +
		"POSTGRES_DB=" + svc + "\n" +
		"# END U-BOOT MANAGED BLOCK: service." + svc + "\n"
}

// seedEnv writes a .env.example under addTestBaseDir.
func seedEnv(t *testing.T, fs *fakeFS, body string) {
	t.Helper()
	if err := fs.WriteFile(filepath.Join(addTestBaseDir, ".env.example"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed .env.example: %v", err)
	}
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

// TestAdd_ActiveWithAllArtifactsIsNoOp pins the T4c contract: an
// Active state with all LH-FA-ADD-002 / LH-AK-002 artefacts present
// (service block + volume block + .env.example block) returns nil-
// error with Changed=nil and writes no files.
func TestAdd_ActiveWithAllArtifactsIsNoOp(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	seedCompose(t, fs, composeBlockComplete("postgres"))
	seedEnv(t, fs, envBlockComplete("postgres"))

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

// TestAdd_MutatingStates_FullExecute pins the T4c happy-path for
// every mutating LH-FA-ADD-005 state: the dispatch reaches
// executeAdd, the slots are filled per the action's plan rule, and
// the deterministic order u-boot.yaml → compose.yaml → .env.example
// is observed in Changed.
func TestAdd_MutatingStates_FullExecute(t *testing.T) {
	cases := []struct {
		name         string
		ubootYAML    string
		wantUBoot    bool   // is u-boot.yaml written?
		wantCompose  bool   // compose.yaml written?
		wantEnv      bool   // .env.example written?
		wantPrior    domain.ServiceState
	}{
		{
			name:        "unregistered",
			ubootYAML:   "schemaVersion: 1\nproject:\n  name: demo\n",
			wantUBoot:   true,
			wantCompose: true,
			wantEnv:     true,
			wantPrior:   domain.ServiceStateUnregistered,
		},
		{
			name: "deactivated",
			ubootYAML: "schemaVersion: 1\nproject:\n  name: demo\n" +
				"services:\n  postgres:\n    enabled: false\n",
			wantUBoot:   true,
			wantCompose: true,
			wantEnv:     true,
			wantPrior:   domain.ServiceStateDeactivated,
		},
		{
			name: "enabled-unset",
			ubootYAML: "schemaVersion: 1\nproject:\n  name: demo\n" +
				"services:\n  postgres: {}\n",
			wantUBoot:   true,
			wantCompose: true,
			wantEnv:     true,
			wantPrior:   domain.ServiceStateEnabledUnset,
		},
		{
			name: "inconsistent-block",
			ubootYAML: "schemaVersion: 1\nproject:\n  name: demo\n" +
				"services:\n  postgres:\n    enabled: true\n",
			wantUBoot:   false, // enabled already true
			wantCompose: true,
			wantEnv:     true,
			wantPrior:   domain.ServiceStateInconsistentBlock,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, fs, _ := newAddService(t)
			seedUBootYAML(t, fs, tc.ubootYAML)
			resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
				BaseDir:     addTestBaseDir,
				ServiceName: postgresName(t),
			})
			if err != nil {
				t.Fatalf("Add: %v", err)
			}
			if resp.PriorState != tc.wantPrior {
				t.Errorf("PriorState = %s, want %s",
					resp.PriorState.String(), tc.wantPrior.String())
			}
			if resp.State != domain.ServiceStateActive {
				t.Errorf("State = %s, want active", resp.State.String())
			}
			// Verify Changed order matches actual writes.
			wantChanged := []string{}
			if tc.wantUBoot {
				wantChanged = append(wantChanged, "u-boot.yaml")
			}
			if tc.wantCompose {
				wantChanged = append(wantChanged, "compose.yaml")
			}
			if tc.wantEnv {
				wantChanged = append(wantChanged, ".env.example")
			}
			if !reflect.DeepEqual(resp.Changed, wantChanged) {
				t.Errorf("Changed = %v, want %v", resp.Changed, wantChanged)
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

