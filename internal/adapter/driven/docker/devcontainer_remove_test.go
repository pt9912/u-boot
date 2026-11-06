package docker_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/docker"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// removeCalls returns the argv lists of the calls that change or ask for the
// state of containers and volumes, without the checks of the environment.
func removeCalls(t *testing.T, f *fakeDocker) [][]string {
	t.Helper()
	var out [][]string
	for _, c := range f.calls(t) {
		switch c[0] {
		case "rm", "volume", "inspect", "stop":
			out = append(out, c)
		}
	}
	return out
}

// SPEC-058 steps 8 and 9: the argument lists. Without volumes `docker rm
// <id>`, with them `docker rm -v <id>`, never `-f`; a volume goes out with
// `docker volume rm <name>`, never `-f`.
func TestLauncher_RemoveContainerAndVolumeArguments(t *testing.T) {
	id := strings.Repeat("ab", 32)
	cases := []struct {
		name string
		run  func(l *docker.Launcher) error
		want []string
	}{
		{"container without volumes", func(l *docker.Launcher) error { return l.RemoveContainer(context.Background(), id, false) }, []string{"rm", id}},
		{"container with volumes", func(l *docker.Launcher) error { return l.RemoveContainer(context.Background(), id, true) }, []string{"rm", "-v", id}},
		{"volume", func(l *docker.Launcher) error { return l.RemoveVolume(context.Background(), "demo-workspace_1.x") }, []string{"volume", "rm", "demo-workspace_1.x"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeDocker(t)
			if err := tc.run(f.launcher(nil)); err != nil {
				t.Fatal(err)
			}
			if got := removeCalls(t, f); !reflect.DeepEqual(got, [][]string{tc.want}) {
				t.Errorf("docker calls = %v, want %v", got, tc.want)
			}
		})
	}
}

// SPEC-058: an argument that is no container ID or volume name is rejected as
// a runtime failure without any docker call; so is a name that could pass as
// an option.
func TestLauncher_RemoveRejectsInvalidArguments(t *testing.T) {
	f := newFakeDocker(t)
	l := f.launcher(nil)
	for _, bad := range []string{"", "xyz", "--all", "-f", "abc def", strings.Repeat("a", 11), "ABCDEF123456"} {
		if err := l.RemoveContainer(context.Background(), bad, true); !errors.Is(err, driven.ErrDevcontainerRuntime) {
			t.Errorf("RemoveContainer(%q) error = %v, want ErrDevcontainerRuntime", bad, err)
		}
	}
	for _, bad := range []string{"", "-f", "--all", "a b", "_a", ".a", "a/b", "a;b", "a\nb"} {
		if err := l.RemoveVolume(context.Background(), bad); !errors.Is(err, driven.ErrDevcontainerRuntime) {
			t.Errorf("RemoveVolume(%q) error = %v, want ErrDevcontainerRuntime", bad, err)
		}
	}
	if len(f.calls(t)) != 0 {
		t.Errorf("docker ran for rejected input: %v", f.flatCalls(t))
	}
}

// SPEC-058, errors on removal: the failure of the engine is a runtime failure
// (exit 12) and carries the message of the engine, so the caller can name it.
func TestLauncher_RemoveFailures(t *testing.T) {
	id := strings.Repeat("ab", 32)
	f := newFakeDocker(t)
	f.flag(t, "rm-fail")
	err := f.launcher(nil).RemoveContainer(context.Background(), id, true)
	if !errors.Is(err, driven.ErrDevcontainerRuntime) || !strings.Contains(err.Error(), "rm-boom") {
		t.Errorf("failing docker rm: error = %v", err)
	}
	f.flag(t, "volume-rm-fail")
	err = f.launcher(nil).RemoveVolume(context.Background(), "demo-volume")
	if !errors.Is(err, driven.ErrDevcontainerRuntime) || !strings.Contains(err.Error(), "volume is in use") {
		t.Errorf("failing docker volume rm: error = %v", err)
	}
}

// SPEC-058: the calls of `remove` start no helper container and mount no
// socket: across the launcher and the inspector, no call is `run`, none
// carries `--mount` or `-f`, and none is `exec`.
func TestRemoveCalls_StartNoHelperAndMountNothing(t *testing.T) {
	id := strings.Repeat("ab", 32)
	f := newFakeDocker(t)
	f.write(t, "inspect-output", `[{"Type":"volume","Name":"demo-vol"}]`)
	f.write(t, "volume-ls-output", "demo-vol\n")
	l := f.launcher(nil)
	insp := docker.WithInspectorBinary(f.bin)
	ctx := context.Background()
	if _, err := insp.ContainerVolumes(ctx, id); err != nil {
		t.Fatal(err)
	}
	if _, err := insp.ExistingVolumes(ctx, []string{"demo-vol"}); err != nil {
		t.Fatal(err)
	}
	for _, err := range []error{l.Stop(ctx, id), l.RemoveContainer(ctx, id, true), l.RemoveVolume(ctx, "demo-vol")} {
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, c := range f.calls(t) {
		if c[0] == "run" || c[0] == "exec" || c[0] == "build" || c[0] == "cp" {
			t.Errorf("a call of remove is docker %s: %v", c[0], c)
		}
		for _, a := range c {
			if a == "--mount" || a == "-f" || a == "--force" {
				t.Errorf("a call of remove carries %s: %v", a, c)
			}
		}
	}
	if len(removeCalls(t, f)) != 5 {
		t.Errorf("docker calls = %v, want inspect, volume ls, stop, rm, volume rm", removeCalls(t, f))
	}
}

// SPEC-058, set of the volumes: the mounts of the type `volume` of `docker
// inspect`, named and anonymous, in the order of the engine; a bind mount and
// a `tmpfs` mount are not in it; a container without a mount has an empty set.
func TestInspector_ContainerVolumes(t *testing.T) {
	t.Parallel()
	anon := strings.Repeat("0123456789abcdef", 4)
	cases := []struct {
		name    string
		out     string
		want    []string
		wantErr bool
	}{
		{"named and anonymous", `[{"Type":"volume","Name":"demo-workspace-x","Destination":"/w"},{"Type":"volume","Name":"` + anon + `"}]`, []string{"demo-workspace-x", anon}, false},
		{"bind and tmpfs are not volumes", `[{"Type":"bind","Source":"/home/x/proj","Destination":"/workspaces/proj"},{"Type":"tmpfs","Destination":"/tmp"},{"Type":"volume","Name":"v1"}]`, []string{"v1"}, false},
		{"only a bind mount", `[{"Type":"bind","Source":"/home/x/proj"}]`, []string{}, false},
		{"no mounts", `[]`, []string{}, false},
		{"null", "null\n", []string{}, false},
		{"not JSON", "oops", nil, true},
		{"name that is no volume name", `[{"Type":"volume","Name":"-rf"}]`, nil, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := &fakeRunner{byCommand: map[string]func([]string) (docker.Process, error){"inspect": textProcess(tc.out)}}
			got, err := newTestInspector(r, "").ContainerVolumes(context.Background(), testContainerID)
			if (err != nil) != tc.wantErr || !reflect.DeepEqual(got, tc.want) {
				t.Errorf("volumes = %#v, err = %v, want %#v (error %t)", got, err, tc.want, tc.wantErr)
			}
			wantCall := []string{"inspect", "--format", "{{json .Mounts}}", testContainerID}
			if len(r.calls) != 1 || !reflect.DeepEqual(r.calls[0], wantCall) {
				t.Errorf("docker calls = %v, want %v", r.calls, wantCall)
			}
		})
	}
	t.Run("engine fails", func(t *testing.T) {
		t.Parallel()
		r := &fakeRunner{byCommand: map[string]func([]string) (docker.Process, error){"inspect": failingProcess("No such container")}}
		if _, err := newTestInspector(r, "").ContainerVolumes(context.Background(), testContainerID); err == nil || !strings.Contains(err.Error(), "No such container") {
			t.Errorf("error = %v", err)
		}
	})
	t.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		r := &fakeRunner{}
		for _, bad := range []string{"", "--all", "xyz"} {
			if _, err := newTestInspector(r, "").ContainerVolumes(context.Background(), bad); err == nil {
				t.Errorf("ContainerVolumes(%q): no error", bad)
			}
		}
		if len(r.calls) != 0 {
			t.Errorf("docker ran for rejected input: %v", r.calls)
		}
	})
}

// SPEC-058 step 9: the volumes that exist are read with `docker volume ls`
// and compared by exact name, in the order asked for; a volume whose name is
// part of another name does not count as existing.
func TestInspector_ExistingVolumes(t *testing.T) {
	t.Parallel()
	r := &fakeRunner{byCommand: map[string]func([]string) (docker.Process, error){"volume": textProcess("other\ndemo-vol-2\r\ndemo-vol-10\n\n")}}
	got, err := newTestInspector(r, "").ExistingVolumes(context.Background(), []string{"demo-vol-1", "demo-vol-2", "gone", "demo-vol-10"})
	if err != nil || !reflect.DeepEqual(got, []string{"demo-vol-2", "demo-vol-10"}) {
		t.Errorf("existing = %#v, err = %v", got, err)
	}
	wantCall := []string{"volume", "ls", "--format", "{{.Name}}"}
	if len(r.calls) != 1 || !reflect.DeepEqual(r.calls[0], wantCall) {
		t.Errorf("docker calls = %v, want %v", r.calls, wantCall)
	}
	none, err := newTestInspector(&fakeRunner{byCommand: map[string]func([]string) (docker.Process, error){"volume": textProcess("")}}, "").
		ExistingVolumes(context.Background(), []string{"a"})
	if err != nil || none == nil || len(none) != 0 {
		t.Errorf("no volume listed: %#v, err = %v, want an empty list", none, err)
	}
	failing := &fakeRunner{byCommand: map[string]func([]string) (docker.Process, error){"volume": failingProcess("daemon gone")}}
	if _, err := newTestInspector(failing, "").ExistingVolumes(context.Background(), []string{"a"}); err == nil || !strings.Contains(err.Error(), "daemon gone") {
		t.Errorf("failing docker volume ls: error = %v", err)
	}
	huge := &fakeRunner{byCommand: map[string]func([]string) (docker.Process, error){"volume": textProcess(strings.Repeat("v\n", 1<<20))}}
	if _, err := newTestInspector(huge, "").ExistingVolumes(context.Background(), []string{"a"}); err == nil {
		t.Errorf("a list cut at the limit must not count as complete")
	}
}
