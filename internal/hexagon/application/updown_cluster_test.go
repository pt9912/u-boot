package application_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// slice-v1-down-volumes-named-list: before/after snapshot of the
// project's declared volumes → names of the removed ones.
func TestDown_RemovedVolumeNames(t *testing.T) {
	t.Parallel()
	f := newDownFixture(t)
	f.engine.projectInfo = driven.ComposeProjectInfo{Name: "demo", Volumes: []driven.ComposeVolume{
		{Key: "postgres-data", Name: "demo_postgres-data"}, {Key: "other", Name: "demo_other"}}}
	// before: both exist; after: only "other" remains (still in use).
	f.engine.volumeSnapshots = [][]string{{"demo_postgres-data", "demo_other", "unrelated"}, {"demo_other", "unrelated"}}

	resp, err := f.svc.Down(context.Background(), driving.DownRequest{BaseDir: "/proj", RemoveVolumes: true, AssumeYes: true})
	if err != nil {
		t.Fatalf("Down: %v", err)
	}
	if len(resp.RemovedVolumeNames) != 1 || resp.RemovedVolumeNames[0] != "demo_postgres-data" {
		t.Errorf("RemovedVolumeNames = %v, want [demo_postgres-data]", resp.RemovedVolumeNames)
	}
	if !resp.RemovedVolumes {
		t.Errorf("RemovedVolumes bool must stay true")
	}
}

// Without --volumes nothing is listed; a failing lookup never fails down.
func TestDown_RemovedVolumeNames_NoneAndBestEffort(t *testing.T) {
	t.Parallel()
	f := newDownFixture(t)
	resp, err := f.svc.Down(context.Background(), driving.DownRequest{BaseDir: "/proj"})
	if err != nil || resp.RemovedVolumeNames != nil {
		t.Errorf("without --volumes: err=%v names=%v, want nil", err, resp.RemovedVolumeNames)
	}
	f2 := newDownFixture(t)
	f2.engine.projectErr = errors.New("compose config failed")
	resp, err = f2.svc.Down(context.Background(), driving.DownRequest{BaseDir: "/proj", RemoveVolumes: true, AssumeYes: true})
	if err != nil || resp.RemovedVolumeNames != nil || !resp.RemovedVolumes {
		t.Errorf("failing lookup: err=%v names=%v removed=%v, want best-effort nil/true", err, resp.RemovedVolumeNames, resp.RemovedVolumes)
	}
}

func purgeEngine() *fakeDockerEngine {
	e := newFakeDockerEngine()
	e.projectInfo = driven.ComposeProjectInfo{Name: "proj", Volumes: []driven.ComposeVolume{{Key: "postgres-data", Name: "proj_postgres-data"}}}
	e.volumeSnapshots = [][]string{{"proj_postgres-data", "unrelated"}}
	return e
}

// slice-v1-volume-auto-removal: --purge really removes the service's
// volume (name resolved before the compose block disappears).
func TestRemove_Purge_RemovesVolume(t *testing.T) {
	t.Parallel()
	fs := seedProjectWithActiveService(t)
	engine := purgeEngine()
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, &fakeConfirmer{}, nil).WithDockerEngine(engine)
	resp, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir: "/proj", ServiceName: mustServiceName(t, "postgres"), Purge: true, Yes: true})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if !resp.VolumesPurged || !resp.PurgeAttempted || len(resp.PurgedVolumes) != 1 || resp.PurgedVolumes[0] != "proj_postgres-data" {
		t.Errorf("resp = purged:%v attempted:%v volumes:%v", resp.VolumesPurged, resp.PurgeAttempted, resp.PurgedVolumes)
	}
	if len(engine.removedVolumes) != 1 || engine.removedVolumes[0] != "proj_postgres-data" {
		t.Errorf("engine removed %v, want only proj_postgres-data (never unrelated volumes)", engine.removedVolumes)
	}
	for _, w := range resp.Warnings {
		t.Errorf("unexpected warning: %+v", w)
	}
}

// A volume that cannot be removed (in use) is a per-volume warning, not
// an error; the service removal itself still succeeded.
func TestRemove_Purge_InUseVolumeWarns(t *testing.T) {
	t.Parallel()
	fs := seedProjectWithActiveService(t)
	engine := purgeEngine()
	engine.removeVolumeErr = map[string]error{"proj_postgres-data": errors.New("volume is in use")}
	svc := application.NewRemoveServiceService(fs, &fakeYAML{}, &fakeConfirmer{}, nil).WithDockerEngine(engine)
	resp, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir: "/proj", ServiceName: mustServiceName(t, "postgres"), Purge: true, Yes: true})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if resp.VolumesPurged || !resp.PurgeAttempted || len(resp.Warnings) != 1 || !strings.Contains(resp.Warnings[0].Message, "in use") {
		t.Errorf("resp = purged:%v attempted:%v warnings:%+v", resp.VolumesPurged, resp.PurgeAttempted, resp.Warnings)
	}
	if resp.State != domain.ServiceStateDeactivated {
		t.Errorf("State = %v, want Deactivated", resp.State)
	}
}

// Dry-run and the declined gate remove nothing; without an engine the
// legacy "deferred" warning stays.
func TestRemove_Purge_NoRemovalWhenDryRunDeclinedOrNoEngine(t *testing.T) {
	t.Parallel()
	req := func(mode driving.PreviewMode, yes bool) driving.RemoveServiceRequest {
		return driving.RemoveServiceRequest{BaseDir: "/proj", ServiceName: mustServiceName(t, "postgres"), Purge: true, Yes: yes, PreviewMode: mode}
	}
	engine := purgeEngine()
	svc := application.NewRemoveServiceService(seedProjectWithActiveService(t), &fakeYAML{}, &fakeConfirmer{}, nil).WithDockerEngine(engine)
	if _, err := svc.Remove(context.Background(), req(driving.PreviewDryRun, false)); err != nil {
		t.Fatalf("dry-run: %v", err)
	}
	if len(engine.removedVolumes) != 0 {
		t.Errorf("dry-run removed volumes: %v", engine.removedVolumes)
	}
	declined := &fakeConfirmer{}
	declined.removeVolumesAnswer = false
	svc2 := application.NewRemoveServiceService(seedProjectWithActiveService(t), &fakeYAML{}, declined, nil).WithDockerEngine(engine)
	if _, err := svc2.Remove(context.Background(), req(driving.PreviewNone, false)); !errors.Is(err, driving.ErrConfirmationRequired) {
		t.Errorf("declined gate: err = %v, want ErrConfirmationRequired", err)
	}
	if len(engine.removedVolumes) != 0 {
		t.Errorf("declined gate removed volumes: %v", engine.removedVolumes)
	}
	svc3 := application.NewRemoveServiceService(seedProjectWithActiveService(t), &fakeYAML{}, &fakeConfirmer{}, nil)
	resp, err := svc3.Remove(context.Background(), req(driving.PreviewNone, true))
	if err != nil || resp.PurgeAttempted || len(resp.Warnings) != 1 || !strings.Contains(resp.Warnings[0].Message, "deferred") {
		t.Errorf("no engine: err=%v attempted=%v warnings=%+v, want the legacy deferred warning", err, resp.PurgeAttempted, resp.Warnings)
	}
}

// slice-v1-recreate-detection: one warning per planned Recreate.
func TestUp_RecreateWarnings(t *testing.T) {
	t.Parallel()
	f := newUpFixture(t, composePostgres)
	f.engine.scriptUp(driven.ComposeUpResult{}, nil)
	f.engine.plan = []driven.ComposePlanAction{{Container: "proj-postgres-1", Action: "Recreate"}, {Container: "proj-web-1", Action: "Start"}}
	f.engine.scriptPsReply([]driven.ComposeService{{Name: "postgres", State: "running", Health: "healthy", Ports: []string{"5432:5432"}}}, nil)
	resp, err := f.svc.Up(context.Background(), driving.UpRequest{BaseDir: "/proj", Timeout: 60 * time.Second})
	if err != nil {
		t.Fatalf("Up: %v", err)
	}
	if len(resp.Warnings) != 1 || resp.Warnings[0].Code != "LH-FA-UP-003" || resp.Warnings[0].Subject != "proj-postgres-1" ||
		!strings.Contains(resp.Warnings[0].Message, "recreated") {
		t.Errorf("warnings = %+v, want one LH-FA-UP-003 recreate warning", resp.Warnings)
	}
	// A failing plan never fails up.
	f2 := newUpFixture(t, composePostgres)
	f2.engine.scriptUp(driven.ComposeUpResult{}, nil)
	f2.engine.planErr = errors.New("dry-run unsupported")
	f2.engine.scriptPsReply([]driven.ComposeService{{Name: "postgres", State: "running", Health: "healthy", Ports: []string{"5432:5432"}}}, nil)
	if resp, err := f2.svc.Up(context.Background(), driving.UpRequest{BaseDir: "/proj", Timeout: 60 * time.Second}); err != nil || len(resp.Warnings) != 0 {
		t.Errorf("failing plan: err=%v warnings=%+v, want none", err, resp.Warnings)
	}
}

// slice-v1-up-partial-snapshot-on-failure: terminal state and timeout
// carry the services seen so far.
func TestUp_PartialSnapshotOnFailure(t *testing.T) {
	t.Parallel()
	f := newUpFixture(t, composePostgres)
	f.engine.scriptUp(driven.ComposeUpResult{}, nil)
	f.engine.scriptPsReply([]driven.ComposeService{{Name: "postgres", State: "exited"}}, nil)
	resp, err := f.svc.Up(context.Background(), driving.UpRequest{BaseDir: "/proj", Timeout: 60 * time.Second})
	if !errors.Is(err, driven.ErrComposeRuntime) {
		t.Fatalf("err = %v, want ErrComposeRuntime", err)
	}
	if len(resp.PartialServices) != 1 || resp.PartialServices[0].Name != "postgres" || resp.PartialServices[0].ContainerStatus != domain.StateDead {
		t.Errorf("PartialServices = %+v", resp.PartialServices)
	}

	f2 := newUpFixture(t, composePostgres)
	f2.engine.scriptUp(driven.ComposeUpResult{}, nil)
	for i := 0; i < 200; i++ {
		f2.engine.scriptPsReply([]driven.ComposeService{{Name: "postgres", State: "starting"}}, nil)
	}
	resp, err = f2.svc.Up(context.Background(), driving.UpRequest{BaseDir: "/proj", Timeout: 2 * time.Second})
	if !errors.Is(err, driving.ErrStabilizationTimeout) || len(resp.PartialServices) != 1 {
		t.Errorf("timeout: err=%v partial=%+v, want ErrStabilizationTimeout with the snapshot", err, resp.PartialServices)
	}
}
