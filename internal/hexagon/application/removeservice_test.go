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

// --- T1-skeleton-pin -------------------------------------------------------

func TestRemoveServiceService_New(t *testing.T) {
	t.Parallel()
	svc := application.NewRemoveServiceService(newFakeFS(), &fakeYAML{}, nil, nil)
	if svc == nil {
		t.Fatal("NewRemoveServiceService returned nil")
	}
}

func TestRemoveServiceService_Remove_EmptyBaseDirRejected(t *testing.T) {
	t.Parallel()
	svc := application.NewRemoveServiceService(newFakeFS(), &fakeYAML{}, nil, nil)
	_, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "",
		ServiceName: mustServiceName(t, "postgres"),
	})
	if err == nil {
		t.Fatal("Remove: want error for empty BaseDir, got nil")
	}
}

// --- T2: state-machine paths ----------------------------------------------

func TestRemoveServiceService_Remove_ProjectNotInitialized(t *testing.T) {
	t.Parallel()
	// u-boot.yaml missing → detect returns ErrProjectNotInitialized.
	fs := newFakeFS()
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, nil, nil)
	_, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustServiceName(t, "postgres"),
	})
	if err == nil {
		t.Fatal("Remove: want ErrProjectNotInitialized, got nil")
	}
	if !errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Errorf("err = %v, want wrap of driving.ErrProjectNotInitialized", err)
	}
}

func TestRemoveServiceService_Remove_UnsupportedService(t *testing.T) {
	t.Parallel()
	fs := seedProjectWithoutService(t)
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, nil, nil)
	_, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustServiceName(t, "keycloak"), // not in catalogue yet
	})
	if err == nil {
		t.Fatal("Remove: want ErrServiceUnsupported, got nil")
	}
	if !errors.Is(err, driving.ErrServiceUnsupported) {
		t.Errorf("err = %v, want wrap of driving.ErrServiceUnsupported", err)
	}
}

func TestRemoveServiceService_Remove_UnregisteredService(t *testing.T) {
	t.Parallel()
	// u-boot.yaml present, no services.postgres entry, no compose block.
	fs := seedProjectWithoutService(t)
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, nil, nil)
	_, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustServiceName(t, "postgres"),
	})
	if err == nil {
		t.Fatal("Remove: want ErrServiceUnregistered, got nil")
	}
	if !errors.Is(err, driving.ErrServiceUnregistered) {
		t.Errorf("err = %v, want wrap of driving.ErrServiceUnregistered", err)
	}
}

func TestRemoveServiceService_Remove_DeactivatedIsIdempotentNoOp(t *testing.T) {
	t.Parallel()
	// services.postgres.enabled: false in u-boot.yaml, no compose block.
	fs := seedProjectWithDeactivatedService(t)
	preWrites := len(fs.writes) // snapshot seed-step writes
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, nil, nil)
	resp, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustServiceName(t, "postgres"),
	})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if resp.PriorState != domain.ServiceStateDeactivated {
		t.Errorf("PriorState = %v, want Deactivated", resp.PriorState)
	}
	if resp.State != domain.ServiceStateDeactivated {
		t.Errorf("State = %v, want Deactivated", resp.State)
	}
	if len(resp.Changed) != 0 {
		t.Errorf("Changed = %v, want nil (idempotent no-op must not write)", resp.Changed)
	}
	if newWrites := len(fs.writes) - preWrites; newWrites != 0 {
		t.Errorf("Remove emitted %d writes; want 0 (no-op path):\n%v", newWrites, fs.writes[preWrites:])
	}
}

func TestRemoveServiceService_Remove_ActiveTransitionsToDeactivated(t *testing.T) {
	t.Parallel()
	fs := seedProjectWithActiveService(t)
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, nil, nil)
	resp, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustServiceName(t, "postgres"),
	})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if resp.PriorState != domain.ServiceStateActive {
		t.Errorf("PriorState = %v, want Active", resp.PriorState)
	}
	if resp.State != domain.ServiceStateDeactivated {
		t.Errorf("State = %v, want Deactivated", resp.State)
	}
	wantChanged := []string{"compose.yaml", ".env.example", "u-boot.yaml"}
	if !equalStringsTest(resp.Changed, wantChanged) {
		t.Errorf("Changed = %v, want %v", resp.Changed, wantChanged)
	}
	// Compose-block should be removed (no BEGIN/END for service.postgres).
	body := string(fs.files["/proj/compose.yaml"])
	if strings.Contains(body, "service.postgres") {
		t.Errorf("compose.yaml still contains service.postgres marker:\n%s", body)
	}
	// u-boot.yaml should now carry enabled: false.
	yamlBody := string(fs.files["/proj/u-boot.yaml"])
	if !strings.Contains(yamlBody, "enabled: false") {
		t.Errorf("u-boot.yaml does not carry enabled: false:\n%s", yamlBody)
	}
}

func TestRemoveServiceService_Remove_ActiveWithoutEnvBlockSkipsEnv(t *testing.T) {
	t.Parallel()
	// Compose block present + u-boot.yaml entry enabled=true, but no
	// .env.example file. Remove still succeeds — env-block-remove is
	// idempotent per-file.
	fs := seedProjectWithActiveService(t)
	delete(fs.files, "/proj/.env.example")
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, nil, nil)
	resp, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustServiceName(t, "postgres"),
	})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	// Changed should be just compose.yaml + u-boot.yaml (env was absent).
	wantChanged := []string{"compose.yaml", "u-boot.yaml"}
	if !equalStringsTest(resp.Changed, wantChanged) {
		t.Errorf("Changed = %v, want %v (env should be skipped silently)", resp.Changed, wantChanged)
	}
}

func TestRemoveServiceService_Remove_ComposeBlockMalformedSurfacesInconsistent(t *testing.T) {
	t.Parallel()
	// compose.yaml has a BEGIN marker but no matching END.
	// removeBlock should surface ErrServiceInconsistent (review-
	// followup path symmetric to detectServiceState's malformed-
	// block branch).
	fs := newFakeFS()
	if err := fs.WriteFile("/proj/u-boot.yaml",
		[]byte("schemaVersion: 1\nproject:\n  name: demo\nservices:\n  postgres:\n    enabled: true\n"), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	// BEGIN-only block — Find returns ErrBlockMalformed, but the
	// state classifier in detectServiceState surfaces it first.
	composeBody := "name: demo\n" +
		"services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  postgres:\n" +
		"    image: postgres:16\n"
	if err := fs.WriteFile("/proj/compose.yaml", []byte(composeBody), 0o644); err != nil {
		t.Fatalf("seed compose.yaml: %v", err)
	}
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, nil, nil)
	_, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustServiceName(t, "postgres"),
	})
	if err == nil {
		t.Fatal("Remove: want ErrServiceInconsistent, got nil")
	}
	if !errors.Is(err, driving.ErrServiceInconsistent) {
		t.Errorf("err = %v, want wrap of driving.ErrServiceInconsistent", err)
	}
}

func TestRemoveServiceService_Remove_FileSystemErrorsPropagate(t *testing.T) {
	t.Parallel()
	// fs.Exists fails on .env.example during the per-file remove
	// loop — surfaces as a wrapped "check .env.example" error,
	// not a fachlicher sentinel. Pin via substring so the path
	// remains debuggable in CI logs.
	fs := seedProjectWithActiveService(t)
	fs.failExistsOn = "/proj/.env.example"
	fs.failExistsErr = errors.New("stat: permission denied")

	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, nil, nil)
	_, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustServiceName(t, "postgres"),
	})
	if err == nil {
		t.Fatal("Remove: want fs.Exists error, got nil")
	}
	if !strings.Contains(err.Error(), ".env.example") {
		t.Errorf("err = %v, want filename in wrap", err)
	}
}

func TestRemoveServiceService_Remove_WriteFailurePropagates(t *testing.T) {
	t.Parallel()
	// Simulate a disk-full / permission error on the final
	// u-boot.yaml write. The compose- and env-block-remove already
	// landed; the YAML write fails and surfaces as a plain error
	// (not a fachlicher sentinel) so the CLI maps it to exit 1.
	fs := seedProjectWithActiveService(t)
	fs.failOn = "/proj/u-boot.yaml"
	fs.failErr = errors.New("write u-boot.yaml: disk full")

	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, nil, nil)
	_, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustServiceName(t, "postgres"),
	})
	if err == nil {
		t.Fatal("Remove: want write-failure error, got nil")
	}
	if !strings.Contains(err.Error(), "write u-boot.yaml") {
		t.Errorf("err = %v, want write-failure wrap", err)
	}
}

func TestRemoveServiceService_Remove_InconsistentBlockConverges(t *testing.T) {
	t.Parallel()
	// Review-followup F1: InconsistentBlock state (services.postgres.
	// enabled: true + no compose block) is now allowed to converge
	// forwards via Remove. The dispatch sends it into executeRemove
	// instead of rejecting with ErrServiceInconsistent. Use case:
	// a previous Remove crashed mid-flight (compose-write OK, env-
	// write failed) — re-running Remove should complete the unfinished
	// work, not block on a manual-cleanup error.
	fs := newFakeFS()
	if err := fs.WriteFile("/proj/u-boot.yaml",
		[]byte("schemaVersion: 1\nproject:\n  name: demo\nservices:\n  postgres:\n    enabled: true\n"), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	if err := fs.WriteFile("/proj/compose.yaml",
		[]byte("name: demo\nservices: {}\n"), 0o644); err != nil {
		t.Fatalf("seed compose.yaml: %v", err)
	}
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, nil, nil)
	resp, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustServiceName(t, "postgres"),
	})
	if err != nil {
		t.Fatalf("Remove: %v (want convergence, not error)", err)
	}
	if resp.State != domain.ServiceStateDeactivated {
		t.Errorf("State = %v, want Deactivated", resp.State)
	}
	// u-boot.yaml should now have enabled: false; compose.yaml's
	// no-block case is a skip (planBlockRemoval returns skip=true);
	// .env.example didn't exist → skip. So Changed should be just
	// u-boot.yaml.
	wantChanged := []string{"u-boot.yaml"}
	if !equalStringsTest(resp.Changed, wantChanged) {
		t.Errorf("Changed = %v, want %v", resp.Changed, wantChanged)
	}
	if !strings.Contains(string(fs.files["/proj/u-boot.yaml"]), "enabled: false") {
		t.Errorf("u-boot.yaml not patched to enabled: false:\n%s", string(fs.files["/proj/u-boot.yaml"]))
	}
}

func TestRemoveServiceService_Remove_EnabledUnsetIsTreatedLikeActive(t *testing.T) {
	t.Parallel()
	// services.postgres entry exists but enabled-key is missing.
	// Compose block present. EnabledUnset → Deactivated transition.
	fs := seedProjectWithEnabledUnset(t)
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, nil, nil)
	resp, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustServiceName(t, "postgres"),
	})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if resp.PriorState != domain.ServiceStateEnabledUnset {
		t.Errorf("PriorState = %v, want EnabledUnset", resp.PriorState)
	}
	if resp.State != domain.ServiceStateDeactivated {
		t.Errorf("State = %v, want Deactivated", resp.State)
	}
}

// --- T3: --purge confirmation gate ----------------------------------------

func TestRemoveServiceService_Remove_PurgeYesAutoApproves(t *testing.T) {
	t.Parallel()
	// --purge + --yes: gate auto-passes, no confirmer call, normal
	// executeRemove path runs. VolumesPurged stays false (T3 defers
	// the actual volume removal; T4 CLI surfaces the gap).
	fs := seedProjectWithActiveService(t)
	conf := &fakeConfirmer{}
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, conf, nil)

	resp, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustServiceName(t, "postgres"),
		Purge:       true,
		Yes:         true,
	})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if resp.State != domain.ServiceStateDeactivated {
		t.Errorf("State = %v, want Deactivated", resp.State)
	}
	if resp.VolumesPurged {
		t.Errorf("VolumesPurged = true; T3 defers volume removal — must stay false")
	}
	if len(conf.removeVolumesCalls) != 0 {
		t.Errorf("Confirmer was called %d times; --yes must short-circuit the prompt", len(conf.removeVolumesCalls))
	}
}

func TestRemoveServiceService_Remove_PurgeNonInteractiveRefuses(t *testing.T) {
	t.Parallel()
	// --purge + --no-interactive without --yes: ErrConfirmationRequired,
	// no FS writes.
	fs := seedProjectWithActiveService(t)
	preWrites := len(fs.writes)
	conf := &fakeConfirmer{}
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, conf, nil)

	_, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:       "/proj",
		ServiceName:   mustServiceName(t, "postgres"),
		Purge:         true,
		NoInteractive: true,
	})
	if err == nil {
		t.Fatal("Remove: want ErrConfirmationRequired, got nil")
	}
	if !errors.Is(err, driving.ErrConfirmationRequired) {
		t.Errorf("err = %v, want wrap of driving.ErrConfirmationRequired", err)
	}
	if newWrites := len(fs.writes) - preWrites; newWrites != 0 {
		t.Errorf("Remove emitted %d writes; refuse path must not touch FS", newWrites)
	}
	if len(conf.removeVolumesCalls) != 0 {
		t.Errorf("Confirmer was called %d times; --no-interactive must short-circuit before prompt", len(conf.removeVolumesCalls))
	}
}

func TestRemoveServiceService_Remove_PurgeInteractiveAccepted(t *testing.T) {
	t.Parallel()
	// --purge + interactive, confirmer says yes → proceed.
	fs := seedProjectWithActiveService(t)
	conf := &fakeConfirmer{removeVolumesAnswer: true}
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, conf, nil)

	resp, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustServiceName(t, "postgres"),
		Purge:       true,
	})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if resp.State != domain.ServiceStateDeactivated {
		t.Errorf("State = %v, want Deactivated", resp.State)
	}
	if len(conf.removeVolumesCalls) != 1 {
		t.Fatalf("Confirmer was called %d times; want 1", len(conf.removeVolumesCalls))
	}
	if conf.removeVolumesCalls[0].BaseDir != "/proj" {
		t.Errorf("Confirmer.BaseDir = %q, want /proj", conf.removeVolumesCalls[0].BaseDir)
	}
}

func TestRemoveServiceService_Remove_PurgeInteractiveDeclined(t *testing.T) {
	t.Parallel()
	// --purge + interactive, confirmer says no → ErrConfirmationRequired.
	fs := seedProjectWithActiveService(t)
	preWrites := len(fs.writes)
	conf := &fakeConfirmer{removeVolumesAnswer: false}
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, conf, nil)

	_, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustServiceName(t, "postgres"),
		Purge:       true,
	})
	if err == nil {
		t.Fatal("Remove: want ErrConfirmationRequired, got nil")
	}
	if !errors.Is(err, driving.ErrConfirmationRequired) {
		t.Errorf("err = %v, want wrap of driving.ErrConfirmationRequired", err)
	}
	if newWrites := len(fs.writes) - preWrites; newWrites != 0 {
		t.Errorf("Remove emitted %d writes; declined path must not touch FS", newWrites)
	}
}

func TestRemoveServiceService_Remove_PurgeOnDeactivatedSkipsGate(t *testing.T) {
	t.Parallel()
	// Review-followup F6: --purge on already-deactivated service.
	// Prior behaviour fired the confirmation gate even though no
	// destructive op would happen — the user's YES to the prompt
	// was theatre. Fixed: gate only fires when state will actually
	// transition, so Deactivated returns the idempotent no-op
	// directly without bothering the confirmer.
	fs := seedProjectWithDeactivatedService(t)
	preWrites := len(fs.writes)
	conf := &fakeConfirmer{removeVolumesAnswer: false} // would refuse if asked
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, conf, nil)

	resp, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustServiceName(t, "postgres"),
		Purge:       true,
	})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if resp.State != domain.ServiceStateDeactivated {
		t.Errorf("State = %v, want Deactivated", resp.State)
	}
	if len(resp.Changed) != 0 {
		t.Errorf("Changed = %v, want nil (Deactivated path is idempotent)", resp.Changed)
	}
	if len(conf.removeVolumesCalls) != 0 {
		t.Errorf("Confirmer was called %d times; want 0 (gate must NOT fire for Deactivated)", len(conf.removeVolumesCalls))
	}
	if newWrites := len(fs.writes) - preWrites; newWrites != 0 {
		t.Errorf("Remove emitted %d writes; Deactivated+Purge must not write", newWrites)
	}
}

func TestRemoveServiceService_Remove_PurgeOnUnregisteredSkipsGate(t *testing.T) {
	t.Parallel()
	// --purge on a service that was never added: ErrServiceUnregistered
	// fires BEFORE the gate so the user gets the most informative
	// error rather than being prompted to confirm cleanup of nothing.
	fs := seedProjectWithoutService(t)
	conf := &fakeConfirmer{}
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, conf, nil)

	_, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:       "/proj",
		ServiceName:   mustServiceName(t, "postgres"),
		Purge:         true,
		NoInteractive: true,
	})
	if err == nil {
		t.Fatal("Remove: want ErrServiceUnregistered, got nil")
	}
	if !errors.Is(err, driving.ErrServiceUnregistered) {
		t.Errorf("err = %v, want ErrServiceUnregistered (not ErrConfirmationRequired)", err)
	}
	if len(conf.removeVolumesCalls) != 0 {
		t.Errorf("Confirmer was called %d times; gate must not fire for Unregistered", len(conf.removeVolumesCalls))
	}
}

// --- fixture helpers ------------------------------------------------------

// seedProjectWithoutService builds a fakeFS containing an
// initialised u-boot project (u-boot.yaml + compose.yaml) without
// any registered service. Used to exercise the "Unregistered"
// state-machine branch.
func seedProjectWithoutService(t *testing.T) *fakeFS {
	t.Helper()
	fs := newFakeFS()
	if err := fs.WriteFile("/proj/u-boot.yaml",
		[]byte("schemaVersion: 1\nproject:\n  name: demo\n"), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	if err := fs.WriteFile("/proj/compose.yaml",
		[]byte("name: demo\nservices: {}\n"), 0o644); err != nil {
		t.Fatalf("seed compose.yaml: %v", err)
	}
	return fs
}

// seedProjectWithDeactivatedService builds a fakeFS with
// services.postgres.enabled: false in u-boot.yaml and no compose
// block. State machine should classify as Deactivated.
func seedProjectWithDeactivatedService(t *testing.T) *fakeFS {
	t.Helper()
	fs := newFakeFS()
	if err := fs.WriteFile("/proj/u-boot.yaml",
		[]byte("schemaVersion: 1\nproject:\n  name: demo\nservices:\n  postgres:\n    enabled: false\n"), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	if err := fs.WriteFile("/proj/compose.yaml",
		[]byte("name: demo\nservices: {}\n"), 0o644); err != nil {
		t.Fatalf("seed compose.yaml: %v", err)
	}
	return fs
}

// seedProjectWithActiveService builds a fakeFS with
// services.postgres.enabled: true, the matching managed block in
// compose.yaml, and a .env.example managed block. Active state.
func seedProjectWithActiveService(t *testing.T) *fakeFS {
	t.Helper()
	fs := newFakeFS()
	if err := fs.WriteFile("/proj/u-boot.yaml",
		[]byte("schemaVersion: 1\nproject:\n  name: demo\nservices:\n  postgres:\n    enabled: true\n"), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	composeBody := "name: demo\n" +
		"services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  postgres:\n" +
		"    image: postgres:16\n" +
		"  # END U-BOOT MANAGED BLOCK: service.postgres\n"
	if err := fs.WriteFile("/proj/compose.yaml", []byte(composeBody), 0o644); err != nil {
		t.Fatalf("seed compose.yaml: %v", err)
	}
	envBody := "# fixed user content\n" +
		"# BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"POSTGRES_USER=demo\n" +
		"# END U-BOOT MANAGED BLOCK: service.postgres\n"
	if err := fs.WriteFile("/proj/.env.example", []byte(envBody), 0o644); err != nil {
		t.Fatalf("seed .env.example: %v", err)
	}
	return fs
}

// seedProjectWithEnabledUnset builds a fakeFS with services.postgres
// present but without the enabled key, and a matching compose block.
// State machine should classify as EnabledUnset.
func seedProjectWithEnabledUnset(t *testing.T) *fakeFS {
	t.Helper()
	fs := newFakeFS()
	// `services.postgres: {}` parses with Enabled-pointer = nil.
	if err := fs.WriteFile("/proj/u-boot.yaml",
		[]byte("schemaVersion: 1\nproject:\n  name: demo\nservices:\n  postgres: {}\n"), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	composeBody := "name: demo\n" +
		"services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  postgres:\n" +
		"    image: postgres:16\n" +
		"  # END U-BOOT MANAGED BLOCK: service.postgres\n"
	if err := fs.WriteFile("/proj/compose.yaml", []byte(composeBody), 0o644); err != nil {
		t.Fatalf("seed compose.yaml: %v", err)
	}
	return fs
}

// equalStringsTest is a small helper to compare slices ignoring slice-
// equality nuances (nil vs empty).
func equalStringsTest(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func mustServiceName(t *testing.T, raw string) domain.ServiceName {
	t.Helper()
	name, err := domain.NewServiceName(raw)
	if err != nil {
		t.Fatalf("NewServiceName(%q): %v", raw, err)
	}
	return name
}
