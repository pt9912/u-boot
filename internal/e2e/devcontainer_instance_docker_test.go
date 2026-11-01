//go:build docker

// End-to-end check of `--instance` of `u-boot devcontainer up`, `down` and
// `shell` (LH-FA-DEV-010, ADR-0038): several containers of one project folder
// run side by side, started with the production start service and, for the
// exit codes and the shell, with the built binary, against the engine of the
// runner. The sandbox test starts two instances and the container without an
// instance with the default settings of the sandbox profile: no nested
// runtime, no feature `ollama`, no resource limit, no check of the host
// resources. It proves the isolation of the workspace volume; the Docker-storage
// and model volumes and the limit of one CPU and 3 GiB are not run by the suite.

package e2e_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const (
	// dcInstanceLabel is the id label of the container of an instance.
	dcInstanceLabel = "u-boot.instance"
)

// dcUpInstance starts the container of the instance (empty: without an
// instance) with the start service.
func dcUpInstance(ctx context.Context, t *testing.T, svc driving.DevcontainerUseCase, dir, instance string, recreate bool) driving.DevcontainerUpResponse {
	t.Helper()
	resp, err := svc.Up(ctx, driving.DevcontainerUpRequest{BaseDir: dir, Instance: instance, Recreate: recreate})
	if err != nil {
		t.Fatalf("Up(instance %q, recreate %v): %v", instance, recreate, err)
	}
	return resp
}

// dcRequireDistinct fails when a value is empty or two values are equal.
func dcRequireDistinct(t *testing.T, what string, values []string) {
	t.Helper()
	seen := map[string]int{}
	for i, v := range values {
		if v == "" {
			t.Errorf("%s of container %d is empty", what, i)
			continue
		}
		if j, dup := seen[v]; dup {
			t.Errorf("%s of container %d equals that of container %d: %s", what, i, j, v)
		}
		seen[v] = i
	}
}

// dcFolderContainers returns the IDs of the containers (running or not) with
// the folder label of dir.
func dcFolderContainers(t *testing.T, dir string) map[string]bool {
	t.Helper()
	ids := map[string]bool{}
	for _, id := range strings.Fields(dcMustDocker(t, "ps", "-aq", "--no-trunc", "--filter", "label="+dcFolderLbl+"="+dir)) {
		ids[id] = true
	}
	return ids
}

// dcRunning reports whether the container runs.
func dcRunning(t *testing.T, id string) bool {
	t.Helper()
	return dcMustDocker(t, "inspect", "--format", "{{.State.Running}}", id) == "true"
}

// dcMarkerPath is the file a marker of an instance lands in, in the workspace
// volume.
func dcMarkerPath(dir, instance string) string {
	return dcWorkspacePath(dir) + "/marker-" + instance
}

// dcWatchNoHelper runs the function and fails when a helper container of the
// binary appeared meanwhile: the call must end before the CLI is started.
func dcWatchNoHelper(t *testing.T, what string, run func()) {
	t.Helper()
	stop := dcWatchHelpers(dcHelperNames())
	run()
	if started := stop(); len(started) != 0 {
		t.Errorf("%s started the helper container %v, want none", what, started)
	}
}

// LH-FA-DEV-010 (ADR-0038), sandbox profile with the default settings, two
// instances of one project folder and the container without an instance:
//
//   - The containers of the instances carry the folder label and the instance
//     label and no `devcontainer.config_file`; the container without an
//     instance carries the folder label and the config label and no instance
//     label. Container IDs and workspace volumes are pairwise different, also
//     against the container without an instance.
//   - A file in the workspace volume of a container is in no other container
//     (every ordered pair).
//   - `up`, `down` and `shell` without `--instance` meet the container without
//     an instance although the instances are newer; a second `up --instance`
//     without `--recreate` meets the same container; `down --instance` stops
//     only that instance; `up --instance --recreate` replaces only that
//     container, the volume and the file stay.
//   - `down --instance <unknown>` stops nothing and ends with 0.
//   - Without a container without an instance, `up`, `down` and `shell`
//     without `--instance` end with exit 10 and a hint at `--instance`, with
//     no helper container and no change to the containers.
func TestDevcontainerInstance_Docker_TwoInstancesAreIsolated(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	requireDocker(t)
	dir := dcGenerateProject(ctx, t, true, nil)
	svc := dcService()
	bin := ubootBinary(t)

	// The container without an instance, then the instances, each newer.
	existing := dcUpInstance(ctx, t, svc, dir, "", false)
	names := []string{"i1", "i2"}
	ids := []string{existing.ContainerID}
	for _, name := range names {
		resp := dcUpInstance(ctx, t, svc, dir, name, false)
		if resp.AlreadyRunning {
			t.Errorf("the first start of instance %s reports already running", name)
		}
		ids = append(ids, resp.ContainerID)
	}
	t.Logf("containers: without instance %s, instances %v", existing.ContainerID[:12], ids[1:])

	// Identity: IDs, labels, workspace volumes.
	dcRequireDistinct(t, "container ID", ids)
	var workspaces []string
	for n, id := range ids {
		doc := dcInspect(t, id)
		if !doc.State.Running {
			t.Errorf("container %s does not run", id[:12])
		}
		if got := doc.Config.Labels[dcFolderLbl]; got != dir {
			t.Errorf("container %s: label %s = %q, want %q", id[:12], dcFolderLbl, got, dir)
		}
		_, hasInstance := doc.Config.Labels[dcInstanceLabel]
		_, hasConfig := doc.Config.Labels[dcConfigLbl]
		if n == 0 {
			if hasInstance || !hasConfig {
				t.Errorf("container without an instance: instance label %v, config label %v, want no instance label and a config label", hasInstance, hasConfig)
			}
		} else {
			if got := doc.Config.Labels[dcInstanceLabel]; got != names[n-1] {
				t.Errorf("container %s: label %s = %q, want %q", id[:12], dcInstanceLabel, got, names[n-1])
			}
			if hasConfig {
				t.Errorf("instance %s carries the label %s, the measurement found none", names[n-1], dcConfigLbl)
			}
		}
		workspaces = append(workspaces, doc.volumeAt(dcWorkspacePath(dir)))
	}
	dcRequireDistinct(t, "workspace volume", workspaces)

	// Content probe over all ordered pairs.
	markers := []string{"base", names[0], names[1]}
	for n, id := range ids {
		if out, code := dcExec(ctx, t, id, "root", "echo "+markers[n]+" > "+dcMarkerPath(dir, markers[n])); code != 0 {
			t.Fatalf("write the marker of %s: %s", markers[n], out)
		}
	}
	for i := range ids {
		for j, idJ := range ids {
			_, code := dcExec(ctx, t, idJ, "root", "test -e "+dcMarkerPath(dir, markers[i]))
			if (code == 0) != (i == j) {
				t.Errorf("file of %s: visible in container %d = %v, want %v", markers[i], j, code == 0, i == j)
			}
		}
	}

	// Without --instance the container without an instance is met although
	// the instances are newer.
	again := dcUpInstance(ctx, t, svc, dir, "", false)
	if again.ContainerID != existing.ContainerID || !again.AlreadyRunning {
		t.Errorf("up without --instance: container %s already running %v, want the container without an instance %s", again.ContainerID, again.AlreadyRunning, existing.ContainerID)
	}
	if out, _, code := runBinary(ctx, t, bin, dir, "hostname\nexit\n", "devcontainer", "shell"); code != 0 || strings.TrimSpace(out) != existing.ContainerID[:12] {
		t.Errorf("shell without --instance: exit %d, hostname %q, want exit 0 and %s", code, strings.TrimSpace(out), existing.ContainerID[:12])
	}
	if out, _, code := runBinary(ctx, t, bin, dir, "hostname\nexit\n", "devcontainer", "shell", "--instance", names[1]); code != 0 || strings.TrimSpace(out) != ids[2][:12] {
		t.Errorf("shell --instance %s: exit %d, hostname %q, want exit 0 and %s", names[1], code, strings.TrimSpace(out), ids[2][:12])
	}
	// A second start of an instance without --recreate meets the same container.
	repeat := dcUpInstance(ctx, t, svc, dir, names[1], false)
	if repeat.ContainerID != ids[2] || !repeat.AlreadyRunning {
		t.Errorf("second up --instance %s: container %s already running %v, want %s already running", names[1], repeat.ContainerID, repeat.AlreadyRunning, ids[2])
	}
	if got := dcFolderContainers(t, dir); len(got) != len(ids) {
		t.Errorf("containers of the project = %d, want %d", len(got), len(ids))
	}

	// down without --instance stops only the container without an instance.
	down, err := svc.Down(ctx, driving.DevcontainerDownRequest{BaseDir: dir})
	if err != nil || !down.Stopped || down.ContainerID != existing.ContainerID {
		t.Fatalf("Down without --instance = %+v err=%v, want the container without an instance stopped", down, err)
	}
	for n, id := range ids {
		if want := n != 0; dcRunning(t, id) != want {
			t.Errorf("after down without --instance: container %d running = %v, want %v", n, !want, want)
		}
	}
	// down for an instance that does not exist stops nothing.
	none, err := svc.Down(ctx, driving.DevcontainerDownRequest{BaseDir: dir, Instance: "nope"})
	if err != nil || none.Stopped || none.ContainerID != "" {
		t.Errorf("Down --instance nope = %+v err=%v, want nothing to stop and no error", none, err)
	}
	// down --instance stops only that instance; up starts the same container again.
	stopped, err := svc.Down(ctx, driving.DevcontainerDownRequest{BaseDir: dir, Instance: names[1]})
	if err != nil || !stopped.Stopped || stopped.ContainerID != ids[2] {
		t.Fatalf("Down --instance %s = %+v err=%v", names[1], stopped, err)
	}
	for n, id := range ids[1:] {
		if want := n != 1; dcRunning(t, id) != want {
			t.Errorf("after down --instance %s: instance %s running = %v, want %v", names[1], names[n], !want, want)
		}
	}
	restarted := dcUpInstance(ctx, t, svc, dir, names[1], false)
	if restarted.ContainerID != ids[2] {
		t.Errorf("up --instance %s after down: container %s, want the same %s", names[1], restarted.ContainerID, ids[2])
	}
	// --recreate replaces only the container of the instance; volume and file stay.
	recreated := dcUpInstance(ctx, t, svc, dir, names[0], true)
	if recreated.ContainerID == ids[1] || dcExists(ids[1]) {
		t.Errorf("recreate of %s: new ID %s, old container still there = %v", names[0], recreated.ContainerID, dcExists(ids[1]))
	}
	doc := dcInspect(t, recreated.ContainerID)
	if doc.volumeAt(dcWorkspacePath(dir)) != workspaces[1] {
		t.Errorf("recreate of %s changed the workspace volume: %+v", names[0], doc.Mounts)
	}
	if out, code := dcExec(ctx, t, recreated.ContainerID, "root", "test -e "+dcMarkerPath(dir, names[0])); code != 0 {
		t.Errorf("recreate of %s: the file is gone: %s", names[0], out)
	}
	if !dcExists(ids[2]) || !dcRunning(t, ids[2]) {
		t.Errorf("recreate of %s touched instance %s (exists %v)", names[0], names[1], dcExists(ids[2]))
	}

	// Without a container without an instance: exit 10, no helper container,
	// no change to the containers.
	baseVolume := dcInspect(t, existing.ContainerID).volumeAt(dcWorkspacePath(dir))
	dcMustDocker(t, "rm", "-f", existing.ContainerID)
	t.Cleanup(func() {
		if baseVolume != "" {
			_, _ = dcDocker("volume", "rm", "-f", baseVolume)
		}
	})
	before := dcFolderContainers(t, dir)
	if len(before) != len(names) {
		t.Fatalf("containers of the project = %d, want only the %d instances", len(before), len(names))
	}
	for _, args := range [][]string{{"devcontainer", "up"}, {"devcontainer", "up", "--recreate"}, {"devcontainer", "down"}, {"devcontainer", "shell"}} {
		dcWatchNoHelper(t, strings.Join(args, " "), func() {
			stdout, stderr, code := runBinary(ctx, t, bin, dir, "hostname\nexit\n", args...)
			if code != 10 || !strings.Contains(stderr+stdout, "--instance") {
				t.Errorf("%s with instances only: exit %d, want 10 with a hint at --instance\nstdout:\n%s\nstderr:\n%s", strings.Join(args, " "), code, stdout, stderr)
			}
		})
	}
	after := dcFolderContainers(t, dir)
	if len(after) != len(before) {
		t.Errorf("containers of the project changed: %d before, %d after", len(before), len(after))
	}
	for id := range before {
		if !after[id] || !dcRunning(t, id) {
			t.Errorf("instance container %s was changed by a call without --instance", id[:12])
		}
	}
	if _, err := svc.Up(ctx, driving.DevcontainerUpRequest{BaseDir: dir}); !errors.Is(err, driving.ErrDevcontainerOnlyInstances) {
		t.Errorf("Up without --instance and with instances only: err = %v, want ErrDevcontainerOnlyInstances", err)
	}
}

// LH-FA-DEV-010 (ADR-0038 decision 6): an instance name the rule rejects ends
// the call with exit 10 before the engine is asked and before the Dev
// Container CLI runs, for `up`, `down` and `shell`; the name `a=b` keeps
// `devcontainer up` from returning (measurement M8), and the empty flag value
// is an invalid name, not "no instance". A valid name passes the check; `up`
// and `shell` then fail for another reason, the missing configuration, which
// is exit 10 as well and which the message tells apart.
func TestDevcontainerInstance_Docker_InvalidNamesEndWithExit10(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	dir := dcProjectDir(t)
	bin := ubootBinary(t)
	for _, sub := range []string{"up", "down", "shell"} {
		for _, name := range []string{"a=b", "Abc", "-a", "a-", "a_b", strings.Repeat("a", 33), ""} {
			dcWatchNoHelper(t, sub+" --instance "+name, func() {
				// A timeout turns a hang of the CLI into a failure of the test.
				callCtx, callCancel := context.WithTimeout(ctx, 90*time.Second)
				defer callCancel()
				stdout, stderr, code := runBinary(callCtx, t, bin, dir, "", "devcontainer", sub, "--instance", name)
				if code != 10 || !strings.Contains(stderr+stdout, "instance") {
					t.Errorf("%s --instance %q: exit %d, want 10 with a message about the instance\nstdout:\n%s\nstderr:\n%s", sub, name, code, stdout, stderr)
				}
			})
		}
		_, stderr, code := runBinary(ctx, t, bin, dir, "", "devcontainer", sub, "--instance", "a-1")
		if sub != "down" && (code != 10 || !strings.Contains(stderr, "devcontainer.json")) {
			t.Errorf("%s --instance a-1 in a folder without configuration: exit %d, want 10 about the missing configuration\n%s", sub, code, stderr)
		}
	}
}

// LH-FA-DEV-010 (ADR-0038), standard profile: instances of one project folder
// share the project folder (bind mount), no volume isolates them. Three
// containers (without an instance and two instances) have different IDs and
// see the same file in the project folder; `up` without `--instance` meets the
// container without an instance although the instances are newer.
func TestDevcontainerInstance_Docker_StandardProfileSharesTheProjectFolder(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, false, nil)
	svc := dcService()

	existing := dcUpInstance(ctx, t, svc, dir, "", false)
	a := dcUpInstance(ctx, t, svc, dir, "a", false)
	b := dcUpInstance(ctx, t, svc, dir, "b", false)
	ids := []string{existing.ContainerID, a.ContainerID, b.ContainerID}
	dcRequireDistinct(t, "container ID", ids)
	for n, id := range ids {
		if doc := dcInspect(t, id); !doc.State.Running || !doc.bindsFolder(dir) {
			t.Errorf("container %d: running %v, binds the project folder %v, want both", n, doc.State.Running, doc.bindsFolder(dir))
		}
	}
	if out, code := dcExec(ctx, t, ids[1], "vscode", "echo shared > "+existing.RemoteWorkspaceFolder+"/shared-marker"); code != 0 {
		t.Fatalf("write the marker in instance a: %s", out)
	}
	for n, id := range ids {
		if out, code := dcExec(ctx, t, id, "vscode", "cat "+existing.RemoteWorkspaceFolder+"/shared-marker"); code != 0 || strings.TrimSpace(out) != "shared" {
			t.Errorf("container %d does not see the file of instance a in the project folder: %q (exit %d)", n, out, code)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "shared-marker")); err != nil {
		t.Errorf("the file of instance a is not in the project folder of the host: %v", err)
	}
	again := dcUpInstance(ctx, t, svc, dir, "", false)
	if again.ContainerID != existing.ContainerID || !again.AlreadyRunning {
		t.Errorf("up without --instance: container %s already running %v, want the container without an instance %s", again.ContainerID, again.AlreadyRunning, existing.ContainerID)
	}
}
