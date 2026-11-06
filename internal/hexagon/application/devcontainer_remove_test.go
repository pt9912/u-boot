package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// removeEngine plays the engine for `devcontainer remove`: it implements the
// launcher and the inspector over one state, so a test reads the order of the
// engine calls in log and the state after the call in containers and volumes.
// The container list is newest first, like the port.
type removeEngine struct {
	mu         sync.Mutex
	log        []string
	containers []*removeContainer
	// volumes is the set of volumes that exist in the engine.
	volumes map[string]bool
	// anonymous names the volumes that `docker rm -v` removes together with the
	// container.
	anonymous map[string]bool
	// held names the volumes another container holds: `docker volume rm` fails.
	held map[string]bool

	findErr, mountsErr, stopErr, rmErr error
	volumeRmErr                        map[string]error
	// lsFailAt makes the n-th call of ExistingVolumes fail (1-based).
	lsFailAt int
	lsCalls  int
}

type removeContainer struct {
	id       string
	instance string
	running  bool
	mounts   []string
}

const (
	removeIDRunning  = "aaaaaaaaaaaa1111111111111111111111111111111111111111111111111111"
	removeIDStopped  = "bbbbbbbbbbbb2222222222222222222222222222222222222222222222222222"
	removeIDInstance = "cccccccccccc3333333333333333333333333333333333333333333333333333"
	removeIDOther    = "dddddddddddd4444444444444444444444444444444444444444444444444444"
)

func newRemoveEngine(containers ...*removeContainer) *removeEngine {
	e := &removeEngine{
		containers: containers, volumes: map[string]bool{}, anonymous: map[string]bool{},
		held: map[string]bool{}, volumeRmErr: map[string]error{},
	}
	for _, c := range containers {
		for _, v := range c.mounts {
			e.volumes[v] = true
		}
	}
	return e
}

func (e *removeEngine) record(event string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.log = append(e.log, event)
}

func (e *removeEngine) events() []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]string(nil), e.log...)
}

func (e *removeEngine) containerIDs() []string {
	var ids []string
	for _, c := range e.containers {
		ids = append(ids, c.id)
	}
	return ids
}

func (e *removeEngine) volumeNames() []string {
	var names []string
	for n := range e.volumes {
		names = append(names, n)
	}
	return sortedStrings(names)
}

func sortedStrings(in []string) []string {
	out := append([]string{}, in...)
	sort.Strings(out)
	return out
}

func (e *removeEngine) FindContainer(_ context.Context, _, instance string) (driven.DevcontainerContainer, bool, error) {
	e.record("find")
	if e.findErr != nil {
		return driven.DevcontainerContainer{}, false, e.findErr
	}
	for _, c := range e.containers {
		if c.instance == instance {
			return driven.DevcontainerContainer{ID: c.id, Running: c.running}, true, nil
		}
	}
	if instance == "" && len(e.containers) > 0 {
		return driven.DevcontainerContainer{}, false, driven.ErrOnlyInstanceContainers
	}
	return driven.DevcontainerContainer{}, false, nil
}

func (e *removeEngine) ContainerVolumes(_ context.Context, id string) ([]string, error) {
	e.record("mounts")
	if e.mountsErr != nil {
		return nil, e.mountsErr
	}
	for _, c := range e.containers {
		if c.id == id {
			return append([]string(nil), c.mounts...), nil
		}
	}
	return nil, errors.New("no such container")
}

func (e *removeEngine) ExistingVolumes(_ context.Context, names []string) ([]string, error) {
	e.record("ls")
	e.lsCalls++
	if e.lsFailAt == e.lsCalls {
		return nil, errors.New("ls-boom")
	}
	out := []string{}
	for _, n := range names {
		if e.volumes[n] {
			out = append(out, n)
		}
	}
	return out, nil
}

func (e *removeEngine) Stop(_ context.Context, id string) error {
	e.record("stop " + short(id))
	if e.stopErr != nil {
		return e.stopErr
	}
	for _, c := range e.containers {
		if c.id == id {
			c.running = false
		}
	}
	return nil
}

func (e *removeEngine) RemoveContainer(_ context.Context, id string, withVolumes bool) error {
	if withVolumes {
		e.record("rm -v " + short(id))
	} else {
		e.record("rm " + short(id))
	}
	if e.rmErr != nil {
		return e.rmErr
	}
	for i, c := range e.containers {
		if c.id != id {
			continue
		}
		if c.running {
			return errors.New("fake engine: cannot remove a running container")
		}
		e.containers = append(e.containers[:i], e.containers[i+1:]...)
		if withVolumes {
			for _, v := range c.mounts {
				if e.anonymous[v] {
					delete(e.volumes, v)
				}
			}
		}
		return nil
	}
	return errors.New("fake engine: no such container")
}

func (e *removeEngine) RemoveVolume(_ context.Context, name string) error {
	e.record("volume rm " + name)
	if err := e.volumeRmErr[name]; err != nil {
		return err
	}
	if e.held[name] {
		return errors.New("volume is in use")
	}
	delete(e.volumes, name)
	return nil
}

func (e *removeEngine) Up(context.Context, driven.DevcontainerUpOptions) (driven.DevcontainerUpResult, error) {
	e.record("up")
	return driven.DevcontainerUpResult{}, errors.New("the helper container must not start")
}

func (e *removeEngine) Shell(context.Context, driven.DevcontainerShellOptions) error {
	e.record("shell")
	return errors.New("no session")
}

func (e *removeEngine) ReadFile(context.Context, string, string, int) ([]byte, error) {
	e.record("readfile")
	return nil, driven.ErrContainerFileNotFound
}

func short(id string) string { return id[:12] }

// removeFixture builds the service on the engine double, a confirmer that logs
// its place in the sequence and a filesystem that holds the two project files,
// so a read of either shows.
type removeFixture struct {
	svc  *application.DevcontainerStartService
	eng  *removeEngine
	conf *fakeConfirmer
	fs   *fakeFS
}

func newRemoveFixture(t *testing.T, eng *removeEngine) *removeFixture {
	t.Helper()
	fs := newFakeFS()
	fs.markDirExists(devcontainerTestBaseDir)
	for rel, body := range map[string]string{".devcontainer/devcontainer.json": plainConfig, "u-boot.yaml": "schemaVersion: 1\n"} {
		if err := fs.WriteFile(filepath.Join(devcontainerTestBaseDir, rel), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	conf := &fakeConfirmer{onRemoveDevcontainer: func() { eng.record("confirm") }}
	svc := application.NewDevcontainerStartService(fs, &fakeYAML{}, eng, eng, nil).WithConfirmer(conf)
	return &removeFixture{svc: svc, eng: eng, conf: conf, fs: fs}
}

func (f *removeFixture) remove(req driving.DevcontainerRemoveRequest) (driving.DevcontainerRemoveResponse, error) {
	req.BaseDir = devcontainerTestBaseDir
	return f.svc.Remove(context.Background(), req)
}

func runningWithVolumes() *removeContainer {
	return &removeContainer{id: removeIDRunning, running: true, mounts: []string{"demo-workspace-xyz", "demo-docker-xyz"}}
}

// SPEC-058 steps 1 to 9: a running container with two named volumes and a
// confirmed `--volumes` ends in this order of engine calls; the volumes are
// named to the confirmer before anything is stopped, and the sorted set is
// what the confirmer, the removal and the response carry.
func TestDevcontainerRemove_VolumesRunInTheOrderOfTheSpec(t *testing.T) {
	t.Parallel()
	f := newRemoveFixture(t, newRemoveEngine(runningWithVolumes()))
	f.conf.removeDevcontainerAnswer = true

	resp, err := f.remove(driving.DevcontainerRemoveRequest{Volumes: true})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"find", "mounts", "confirm", "stop aaaaaaaaaaaa", "rm -v aaaaaaaaaaaa",
		"ls", "volume rm demo-docker-xyz", "volume rm demo-workspace-xyz", "ls",
	}
	if got := f.eng.events(); !reflect.DeepEqual(got, want) {
		t.Errorf("engine calls:\n got  %v\n want %v", got, want)
	}
	sorted := []string{"demo-docker-xyz", "demo-workspace-xyz"}
	if len(f.conf.removeDevcontainerCalls) != 1 || f.conf.removeDevcontainerCalls[0].ContainerID != removeIDRunning ||
		!reflect.DeepEqual(f.conf.removeDevcontainerCalls[0].Volumes, sorted) {
		t.Errorf("confirmer calls = %+v, want the container and the sorted volumes", f.conf.removeDevcontainerCalls)
	}
	wantResp := driving.DevcontainerRemoveResponse{ContainerID: removeIDRunning, Stopped: true, Removed: true, RemovedVolumeNames: sorted}
	if !reflect.DeepEqual(resp, wantResp) {
		t.Errorf("response = %+v, want %+v", resp, wantResp)
	}
	if len(f.eng.containerIDs()) != 0 || len(f.eng.volumeNames()) != 0 {
		t.Errorf("state after: containers %v, volumes %v, want none", f.eng.containerIDs(), f.eng.volumeNames())
	}
}

// SPEC-058 step 7 and 8: without `--volumes` the command asks nothing, reads
// no volume, passes no `-v` and removes no volume; a stopped container is not
// stopped again.
func TestDevcontainerRemove_WithoutVolumesTouchesNoVolume(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		running bool
		want    []string
		stopped bool
	}{
		{"running", true, []string{"find", "stop aaaaaaaaaaaa", "rm aaaaaaaaaaaa"}, true},
		{"stopped", false, []string{"find", "rm aaaaaaaaaaaa"}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := runningWithVolumes()
			c.running = tc.running
			f := newRemoveFixture(t, newRemoveEngine(c))
			resp, err := f.remove(driving.DevcontainerRemoveRequest{})
			if err != nil {
				t.Fatal(err)
			}
			if got := f.eng.events(); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("engine calls:\n got  %v\n want %v", got, tc.want)
			}
			if resp.Stopped != tc.stopped || !resp.Removed || resp.ContainerID != removeIDRunning || resp.RemovedVolumeNames == nil || len(resp.RemovedVolumeNames) != 0 {
				t.Errorf("response = %+v", resp)
			}
			if len(f.conf.removeDevcontainerCalls) != 0 {
				t.Errorf("the confirmer was asked without --volumes: %+v", f.conf.removeDevcontainerCalls)
			}
			if got := f.eng.volumeNames(); !reflect.DeepEqual(got, []string{"demo-docker-xyz", "demo-workspace-xyz"}) {
				t.Errorf("volumes after = %v, want both to stay", got)
			}
		})
	}
}

// SPEC-058, table of the release: the six rows of the table. A row without
// release changes nothing, and a container that ran keeps running.
func TestDevcontainerRemove_ReleaseTable(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		req        driving.DevcontainerRemoveRequest
		noMounts   bool
		answer     bool
		answerErr  error
		wantErr    error
		wantAsked  bool
		wantAction bool
	}{
		{name: "no volumes flag", req: driving.DevcontainerRemoveRequest{NonInteractive: true}, wantAction: true},
		{name: "yes", req: driving.DevcontainerRemoveRequest{Volumes: true, AssumeYes: true}, wantAction: true},
		{name: "no-interactive", req: driving.DevcontainerRemoveRequest{Volumes: true, NonInteractive: true}, wantErr: driving.ErrConfirmationRequired},
		{name: "json", req: driving.DevcontainerRemoveRequest{Volumes: true, SilenceConfirmer: true}, wantErr: driving.ErrConfirmationRequired},
		{name: "json with yes", req: driving.DevcontainerRemoveRequest{Volumes: true, SilenceConfirmer: true, AssumeYes: true}, wantAction: true},
		{name: "asked, agreed", req: driving.DevcontainerRemoveRequest{Volumes: true}, answer: true, wantAsked: true, wantAction: true},
		{name: "asked, declined", req: driving.DevcontainerRemoveRequest{Volumes: true}, wantErr: driving.ErrConfirmationRequired, wantAsked: true},
		{name: "asked, unreadable", req: driving.DevcontainerRemoveRequest{Volumes: true}, answerErr: errors.New("read boom"), wantErr: driving.ErrConfirmerUnavailable, wantAsked: true},
		{name: "asked although no volume", req: driving.DevcontainerRemoveRequest{Volumes: true}, noMounts: true, answer: true, wantAsked: true, wantAction: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := runningWithVolumes()
			if tc.noMounts {
				c.mounts = nil
			}
			f := newRemoveFixture(t, newRemoveEngine(c))
			f.conf.removeDevcontainerAnswer = tc.answer
			f.conf.removeDevcontainerErr = tc.answerErr
			resp, err := f.remove(tc.req)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
				if errors.Is(err, driven.ErrDevcontainerRuntime) {
					t.Errorf("a missing release is no runtime failure: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if asked := len(f.conf.removeDevcontainerCalls) > 0; asked != tc.wantAsked {
				t.Errorf("asked = %t, want %t", asked, tc.wantAsked)
			}
			for _, e := range f.eng.events() {
				if !tc.wantAction && (strings.HasPrefix(e, "stop") || strings.HasPrefix(e, "rm") || strings.HasPrefix(e, "volume rm")) {
					t.Errorf("%q ran without a release: %v", e, f.eng.events())
				}
			}
			if !tc.wantAction {
				if resp.Removed || resp.Stopped {
					t.Errorf("response = %+v, want nothing done", resp)
				}
				if !c.running || len(f.eng.containerIDs()) != 1 || len(f.eng.volumeNames()) != 2 {
					t.Errorf("the container must keep running and keep its volumes: running=%t volumes=%v", c.running, f.eng.volumeNames())
				}
			}
			if tc.wantAsked && tc.noMounts && (len(f.conf.removeDevcontainerCalls) != 1 || len(f.conf.removeDevcontainerCalls[0].Volumes) != 0) {
				t.Errorf("confirmer calls = %+v, want one call with no volume", f.conf.removeDevcontainerCalls)
			}
		})
	}
}

// SPEC-058, messages of the exit code 10 without release: they name the
// container, the volumes and that nothing was stopped or removed; `--yes` is
// named where no question can be asked, and only there.
func TestDevcontainerRemove_RefusalMessages(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		req     driving.DevcontainerRemoveRequest
		wantYes bool
	}{
		{"no-interactive", driving.DevcontainerRemoveRequest{Volumes: true, NonInteractive: true}, true},
		{"json", driving.DevcontainerRemoveRequest{Volumes: true, SilenceConfirmer: true}, true},
		{"declined", driving.DevcontainerRemoveRequest{Volumes: true}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := newRemoveFixture(t, newRemoveEngine(runningWithVolumes()))
			_, err := f.remove(tc.req)
			if !errors.Is(err, driving.ErrConfirmationRequired) {
				t.Fatalf("error = %v", err)
			}
			msg := err.Error()
			for _, part := range []string{"aaaaaaaaaaaa", "demo-docker-xyz, demo-workspace-xyz", "neither stopped nor removed"} {
				if !strings.Contains(msg, part) {
					t.Errorf("message %q lacks %q", msg, part)
				}
			}
			if strings.Contains(msg, "--yes") != tc.wantYes {
				t.Errorf("message %q: --yes named = %t, want %t", msg, strings.Contains(msg, "--yes"), tc.wantYes)
			}
		})
	}
	f := newRemoveFixture(t, newRemoveEngine(runningWithVolumes()))
	f.conf.removeDevcontainerErr = errors.New("read boom")
	_, err := f.remove(driving.DevcontainerRemoveRequest{Volumes: true})
	if !errors.Is(err, driving.ErrConfirmerUnavailable) || !strings.Contains(err.Error(), "read boom") || !strings.Contains(err.Error(), "neither stopped nor removed") {
		t.Errorf("unreadable answer: %v", err)
	}
}

// A service without a wired confirmer does not release `--volumes`.
func TestDevcontainerRemove_WithoutConfirmerDeclines(t *testing.T) {
	t.Parallel()
	eng := newRemoveEngine(runningWithVolumes())
	fs := newFakeFS()
	svc := application.NewDevcontainerStartService(fs, &fakeYAML{}, eng, eng, nil)
	_, err := svc.Remove(context.Background(), driving.DevcontainerRemoveRequest{BaseDir: devcontainerTestBaseDir, Volumes: true})
	if !errors.Is(err, driving.ErrConfirmationRequired) {
		t.Errorf("error = %v, want ErrConfirmationRequired", err)
	}
	if got := eng.events(); !reflect.DeepEqual(got, []string{"find", "mounts"}) {
		t.Errorf("engine calls = %v", got)
	}
}

// SPEC-058, selection of the container: the instance, the container without an
// instance next to instances, and the project with only instances.
func TestDevcontainerRemove_SelectionOfTheContainer(t *testing.T) {
	t.Parallel()
	newWorld := func() *removeEngine {
		return newRemoveEngine(
			&removeContainer{id: removeIDInstance, instance: "agent", running: true, mounts: []string{"inst-workspace"}},
			&removeContainer{id: removeIDRunning, running: true, mounts: []string{"plain-workspace"}},
		)
	}
	t.Run("instance", func(t *testing.T) {
		t.Parallel()
		f := newRemoveFixture(t, newWorld())
		resp, err := f.remove(driving.DevcontainerRemoveRequest{Instance: "agent", Volumes: true, AssumeYes: true})
		if err != nil || resp.ContainerID != removeIDInstance || !reflect.DeepEqual(resp.RemovedVolumeNames, []string{"inst-workspace"}) {
			t.Fatalf("resp = %+v, err = %v", resp, err)
		}
		if got := f.eng.containerIDs(); !reflect.DeepEqual(got, []string{removeIDRunning}) {
			t.Errorf("containers after = %v, want the one without instance", got)
		}
		if got := f.eng.volumeNames(); !reflect.DeepEqual(got, []string{"plain-workspace"}) {
			t.Errorf("volumes after = %v, want the volume of the other container", got)
		}
	})
	t.Run("without instance next to an instance", func(t *testing.T) {
		t.Parallel()
		f := newRemoveFixture(t, newWorld())
		resp, err := f.remove(driving.DevcontainerRemoveRequest{Volumes: true, AssumeYes: true})
		if err != nil || resp.ContainerID != removeIDRunning {
			t.Fatalf("resp = %+v, err = %v", resp, err)
		}
		if got := f.eng.containerIDs(); !reflect.DeepEqual(got, []string{removeIDInstance}) {
			t.Errorf("containers after = %v, want the one of the instance", got)
		}
		if got := f.eng.volumeNames(); !reflect.DeepEqual(got, []string{"inst-workspace"}) {
			t.Errorf("volumes after = %v, want the volume of the instance", got)
		}
	})
	t.Run("only instances", func(t *testing.T) {
		t.Parallel()
		for _, running := range []bool{true, false} {
			eng := newRemoveEngine(&removeContainer{id: removeIDInstance, instance: "agent", running: running, mounts: []string{"inst-workspace"}})
			f := newRemoveFixture(t, eng)
			resp, err := f.remove(driving.DevcontainerRemoveRequest{Volumes: true, AssumeYes: true})
			if !errors.Is(err, driving.ErrDevcontainerOnlyInstances) || !strings.Contains(err.Error(), "--instance") || resp.Removed {
				t.Fatalf("running=%t: resp = %+v, err = %v", running, resp, err)
			}
			if got := eng.events(); !reflect.DeepEqual(got, []string{"find"}) {
				t.Errorf("running=%t: engine calls = %v, want the search only", running, got)
			}
		}
	})
	t.Run("instance without a container, others present", func(t *testing.T) {
		t.Parallel()
		f := newRemoveFixture(t, newWorld())
		resp, err := f.remove(driving.DevcontainerRemoveRequest{Instance: "none", Volumes: true})
		if err != nil || resp.Removed || resp.ContainerID != "" || resp.RemovedVolumeNames == nil {
			t.Fatalf("resp = %+v, err = %v", resp, err)
		}
		if got := f.eng.events(); !reflect.DeepEqual(got, []string{"find"}) {
			t.Errorf("engine calls = %v, want the search only", got)
		}
		if len(f.conf.removeDevcontainerCalls) != 0 || len(f.eng.containerIDs()) != 2 || len(f.eng.volumeNames()) != 2 {
			t.Errorf("the other containers and volumes must stay and nothing is asked")
		}
	})
	t.Run("no container", func(t *testing.T) {
		t.Parallel()
		f := newRemoveFixture(t, newRemoveEngine())
		resp, err := f.remove(driving.DevcontainerRemoveRequest{Volumes: true, NonInteractive: true})
		if err != nil || resp.Removed || resp.ContainerID != "" || resp.RemovedVolumeNames == nil || len(resp.RemovedVolumeNames) != 0 {
			t.Fatalf("resp = %+v, err = %v, want exit 0 although --no-interactive and no --yes", resp, err)
		}
		if len(f.conf.removeDevcontainerCalls) != 0 {
			t.Errorf("asked although no container")
		}
	})
}

// SPEC-058: of several containers with the same labels the port returns the
// newest, and the command acts on that one; the others stay.
func TestDevcontainerRemove_ActsOnTheContainerTheSearchReturns(t *testing.T) {
	t.Parallel()
	eng := newRemoveEngine(
		&removeContainer{id: removeIDStopped, mounts: []string{"newest-vol"}},
		&removeContainer{id: removeIDOther, mounts: []string{"older-vol"}},
	)
	f := newRemoveFixture(t, eng)
	if _, err := f.remove(driving.DevcontainerRemoveRequest{Volumes: true, AssumeYes: true}); err != nil {
		t.Fatal(err)
	}
	if got := eng.containerIDs(); !reflect.DeepEqual(got, []string{removeIDOther}) {
		t.Errorf("containers after = %v, want the older one", got)
	}
	if got := eng.volumeNames(); !reflect.DeepEqual(got, []string{"older-vol"}) {
		t.Errorf("volumes after = %v, want the volume of the older one", got)
	}
}

// SPEC-058, set of the volumes: sorted, without repetition, and the set of the
// container without a volume is empty (a JSON array, never nil).
func TestDevcontainerRemove_SetOfTheVolumes(t *testing.T) {
	t.Parallel()
	t.Run("sorted and distinct", func(t *testing.T) {
		t.Parallel()
		c := &removeContainer{id: removeIDStopped, mounts: []string{"b", "a", "b", "c"}}
		f := newRemoveFixture(t, newRemoveEngine(c))
		f.conf.removeDevcontainerAnswer = true
		resp, err := f.remove(driving.DevcontainerRemoveRequest{Volumes: true})
		if err != nil || !reflect.DeepEqual(resp.RemovedVolumeNames, []string{"a", "b", "c"}) {
			t.Fatalf("resp = %+v, err = %v", resp, err)
		}
		if !reflect.DeepEqual(f.conf.removeDevcontainerCalls[0].Volumes, []string{"a", "b", "c"}) {
			t.Errorf("confirmer volumes = %v", f.conf.removeDevcontainerCalls[0].Volumes)
		}
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		f := newRemoveFixture(t, newRemoveEngine(&removeContainer{id: removeIDStopped}))
		resp, err := f.remove(driving.DevcontainerRemoveRequest{Volumes: true, AssumeYes: true})
		if err != nil || !resp.Removed || resp.RemovedVolumeNames == nil || len(resp.RemovedVolumeNames) != 0 {
			t.Fatalf("resp = %+v, err = %v", resp, err)
		}
		if got := f.eng.events(); !reflect.DeepEqual(got, []string{"find", "mounts", "rm -v bbbbbbbbbbbb"}) {
			t.Errorf("engine calls = %v, want no listing of volumes", got)
		}
	})
}

// SPEC-058, anonymous volumes: `docker rm -v` removes them with the container;
// the step on the volumes finds them gone, does not remove them again and
// counts them as removed.
func TestDevcontainerRemove_AnonymousVolumeGoneWithTheContainer(t *testing.T) {
	t.Parallel()
	anon := strings.Repeat("0123456789abcdef", 4)
	c := &removeContainer{id: removeIDStopped, mounts: []string{"named-vol", anon}}
	eng := newRemoveEngine(c)
	eng.anonymous[anon] = true
	f := newRemoveFixture(t, eng)
	resp, err := f.remove(driving.DevcontainerRemoveRequest{Volumes: true, AssumeYes: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range eng.events() {
		if e == "volume rm "+anon {
			t.Errorf("the anonymous volume was removed twice: %v", eng.events())
		}
	}
	want := sortedStrings([]string{"named-vol", anon})
	if !reflect.DeepEqual(resp.RemovedVolumeNames, want) || len(eng.volumeNames()) != 0 {
		t.Errorf("removed = %v, volumes after = %v, want both removed", resp.RemovedVolumeNames, eng.volumeNames())
	}
}

// SPEC-058, errors on removal (exit 12), one row of the table each: the state
// after, the sentinel identity and the content of the message.
func TestDevcontainerRemove_FailuresKeepTheStateOfTheTable(t *testing.T) {
	t.Parallel()
	req := driving.DevcontainerRemoveRequest{Volumes: true, AssumeYes: true}
	t.Run("docker stop", func(t *testing.T) {
		t.Parallel()
		eng := newRemoveEngine(runningWithVolumes())
		eng.stopErr = errors.Join(errors.New("stop-boom"), driven.ErrDevcontainerRuntime)
		f := newRemoveFixture(t, eng)
		_, err := f.remove(req)
		if !errors.Is(err, driven.ErrDevcontainerRuntime) || !strings.Contains(err.Error(), "docker stop") || !strings.Contains(err.Error(), "stop-boom") ||
			!strings.Contains(err.Error(), "may still run") {
			t.Fatalf("error = %v", err)
		}
		if got := eng.events(); got[len(got)-1] != "stop aaaaaaaaaaaa" {
			t.Errorf("engine calls = %v, want nothing after the failed stop", got)
		}
		if len(eng.containerIDs()) != 1 || len(eng.volumeNames()) != 2 {
			t.Errorf("state after: containers %v, volumes %v", eng.containerIDs(), eng.volumeNames())
		}
	})
	t.Run("docker rm", func(t *testing.T) {
		t.Parallel()
		eng := newRemoveEngine(runningWithVolumes())
		eng.rmErr = errors.Join(errors.New("rm-boom"), driven.ErrDevcontainerRuntime)
		f := newRemoveFixture(t, eng)
		_, err := f.remove(req)
		if !errors.Is(err, driven.ErrDevcontainerRuntime) || !strings.Contains(err.Error(), "docker rm") || !strings.Contains(err.Error(), "rm-boom") ||
			!strings.Contains(err.Error(), "(stopped)") {
			t.Fatalf("error = %v", err)
		}
		if got := eng.events(); got[len(got)-1] != "rm -v aaaaaaaaaaaa" {
			t.Errorf("engine calls = %v, want no volume step after the failed rm", got)
		}
		if len(eng.containerIDs()) != 1 || eng.containers[0].running || len(eng.volumeNames()) != 2 {
			t.Errorf("state after: containers %v, volumes %v", eng.containerIDs(), eng.volumeNames())
		}
	})
	t.Run("docker volume rm of a volume another container holds", func(t *testing.T) {
		t.Parallel()
		eng := newRemoveEngine(runningWithVolumes())
		eng.held["demo-docker-xyz"] = true
		f := newRemoveFixture(t, eng)
		_, err := f.remove(req)
		if !errors.Is(err, driven.ErrDevcontainerRuntime) {
			t.Fatalf("error = %v, want the runtime class", err)
		}
		msg := err.Error()
		for _, part := range []string{"aaaaaaaaaaaa", "volumes removed: demo-workspace-xyz", "volumes that remain: demo-docker-xyz", "volume is in use"} {
			if !strings.Contains(msg, part) {
				t.Errorf("message %q lacks %q", msg, part)
			}
		}
		// every volume was tried although the first failed
		wantTail := []string{"volume rm demo-docker-xyz", "volume rm demo-workspace-xyz", "ls"}
		got := eng.events()
		if !reflect.DeepEqual(got[len(got)-3:], wantTail) {
			t.Errorf("engine calls = %v, want %v at the end", got, wantTail)
		}
		if len(eng.containerIDs()) != 0 || !reflect.DeepEqual(eng.volumeNames(), []string{"demo-docker-xyz"}) {
			t.Errorf("state after: containers %v, volumes %v", eng.containerIDs(), eng.volumeNames())
		}
	})
	t.Run("a volume of the set remains after the step", func(t *testing.T) {
		t.Parallel()
		// the engine answers the removal with success and keeps the volume
		eng := newRemoveEngine(&removeContainer{id: removeIDStopped, mounts: []string{"ghost"}})
		keep := &ghostEngine{removeEngine: eng}
		svc := application.NewDevcontainerStartService(newFakeFS(), &fakeYAML{}, keep, keep, nil)
		_, err := svc.Remove(context.Background(), driving.DevcontainerRemoveRequest{BaseDir: devcontainerTestBaseDir, Volumes: true, AssumeYes: true})
		if !errors.Is(err, driven.ErrDevcontainerRuntime) || !strings.Contains(err.Error(), "volumes that remain: ghost") {
			t.Fatalf("error = %v", err)
		}
	})
	t.Run("listing before the removal fails", func(t *testing.T) {
		t.Parallel()
		eng := newRemoveEngine(runningWithVolumes())
		eng.lsFailAt = 1
		f := newRemoveFixture(t, eng)
		_, err := f.remove(req)
		if !errors.Is(err, driven.ErrDevcontainerRuntime) || !strings.Contains(err.Error(), "ls-boom") || !strings.Contains(err.Error(), "is removed") {
			t.Fatalf("error = %v", err)
		}
		for _, e := range eng.events() {
			if strings.HasPrefix(e, "volume rm") {
				t.Errorf("a volume was removed without knowing its state: %v", eng.events())
			}
		}
	})
	t.Run("listing after the removal fails", func(t *testing.T) {
		t.Parallel()
		eng := newRemoveEngine(runningWithVolumes())
		eng.lsFailAt = 2
		f := newRemoveFixture(t, eng)
		_, err := f.remove(req)
		if !errors.Is(err, driven.ErrDevcontainerRuntime) || !strings.Contains(err.Error(), "ls-boom") || !strings.Contains(err.Error(), "not established") {
			t.Fatalf("error = %v", err)
		}
	})
}

// ghostEngine reports every volume as existing after `docker volume rm`, as
// an engine that answers the removal with success and keeps the volume.
type ghostEngine struct{ *removeEngine }

func (g *ghostEngine) RemoveVolume(_ context.Context, name string) error {
	g.record("volume rm " + name)
	return nil
}

// SPEC-058, exit code 11: the engine cannot be asked in the search or while
// reading the volumes; nothing is asked and nothing is changed.
func TestDevcontainerRemove_EngineNotAskable(t *testing.T) {
	t.Parallel()
	t.Run("search fails", func(t *testing.T) {
		t.Parallel()
		eng := newRemoveEngine(runningWithVolumes())
		eng.findErr = errors.New("docker context inspect: not found")
		f := newRemoveFixture(t, eng)
		_, err := f.remove(driving.DevcontainerRemoveRequest{Volumes: true, AssumeYes: true})
		if !errors.Is(err, driven.ErrDockerUnavailable) || !reflect.DeepEqual(eng.events(), []string{"find"}) {
			t.Errorf("error = %v, calls = %v", err, eng.events())
		}
	})
	t.Run("remote engine", func(t *testing.T) {
		t.Parallel()
		eng := newRemoveEngine(runningWithVolumes())
		eng.findErr = driven.ErrContainerEngineNotLocal
		f := newRemoveFixture(t, eng)
		if _, err := f.remove(driving.DevcontainerRemoveRequest{}); !errors.Is(err, driven.ErrContainerEngineNotLocal) {
			t.Errorf("error = %v", err)
		}
	})
	t.Run("reading the volumes fails", func(t *testing.T) {
		t.Parallel()
		eng := newRemoveEngine(runningWithVolumes())
		eng.mountsErr = errors.New("inspect-boom")
		f := newRemoveFixture(t, eng)
		_, err := f.remove(driving.DevcontainerRemoveRequest{Volumes: true, AssumeYes: true})
		if !errors.Is(err, driven.ErrDockerUnavailable) || !strings.Contains(err.Error(), "inspect-boom") {
			t.Fatalf("error = %v", err)
		}
		if got := eng.events(); !reflect.DeepEqual(got, []string{"find", "mounts"}) || len(f.conf.removeDevcontainerCalls) != 0 {
			t.Errorf("calls = %v: nothing may be asked or changed", got)
		}
	})
}

// SPEC-058 step 2: an instance name that is not admissible ends with
// ErrDevcontainerInstanceInvalid before the engine is asked (limit: 32
// characters pass, 33 do not).
func TestDevcontainerRemove_InstanceNameIsCheckedBeforeTheEngine(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"a=b", "Abc", "-a", "a-", "a_b", "a b", strings.Repeat("a", 33)} {
		eng := newRemoveEngine(runningWithVolumes())
		f := newRemoveFixture(t, eng)
		_, err := f.remove(driving.DevcontainerRemoveRequest{Instance: name, Volumes: true, AssumeYes: true})
		if !errors.Is(err, driving.ErrDevcontainerInstanceInvalid) || len(eng.events()) != 0 {
			t.Errorf("instance %q: error = %v, calls = %v", name, err, eng.events())
		}
	}
	eng := newRemoveEngine()
	f := newRemoveFixture(t, eng)
	if _, err := f.remove(driving.DevcontainerRemoveRequest{Instance: strings.Repeat("a", 32)}); err != nil {
		t.Errorf("32 characters: %v", err)
	}
	if got := eng.events(); !reflect.DeepEqual(got, []string{"find"}) {
		t.Errorf("32 characters: calls = %v, want the search", got)
	}
}

// SPEC-058: the command reads no project file, starts no helper container and
// opens no session; every project file stays unread and unwritten, in each of
// the paths.
func TestDevcontainerRemove_ReadsNoProjectFileAndStartsNoHelper(t *testing.T) {
	t.Parallel()
	paths := []driving.DevcontainerRemoveRequest{
		{},
		{Volumes: true, AssumeYes: true},
		{Volumes: true, NonInteractive: true},
		{Instance: "none"},
	}
	for _, req := range paths {
		f := newRemoveFixture(t, newRemoveEngine(runningWithVolumes()))
		_, _ = f.remove(req)
		for _, e := range f.eng.events() {
			if e == "up" || e == "readfile" || e == "shell" {
				t.Errorf("%+v: engine call %q", req, e)
			}
		}
		for _, rel := range []string{".devcontainer/devcontainer.json", "u-boot.yaml"} {
			if n := f.fs.readFileCallCount(filepath.Join(devcontainerTestBaseDir, rel)); n != 0 {
				t.Errorf("%+v: %s was read %d times", req, rel, n)
			}
		}
		if len(f.fs.writtenPaths()) != 2 {
			t.Errorf("%+v: the filesystem was written after the fixture: %v", req, f.fs.writtenPaths())
		}
	}
}

// A missing BaseDir or a missing port is a wiring error, not "nothing to do".
func TestDevcontainerRemove_NotWired(t *testing.T) {
	t.Parallel()
	svc := application.NewDevcontainerStartService(newFakeFS(), &fakeYAML{}, nil, nil, nil)
	if _, err := svc.Remove(context.Background(), driving.DevcontainerRemoveRequest{BaseDir: devcontainerTestBaseDir}); err == nil || !strings.Contains(err.Error(), "not wired") {
		t.Errorf("error = %v", err)
	}
	eng := newRemoveEngine()
	svc = application.NewDevcontainerStartService(newFakeFS(), &fakeYAML{}, eng, eng, nil)
	if _, err := svc.Remove(context.Background(), driving.DevcontainerRemoveRequest{}); err == nil || !strings.Contains(err.Error(), "BaseDir") {
		t.Errorf("error = %v", err)
	}
}
