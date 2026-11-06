//go:build docker

// End-to-end check of `u-boot devcontainer remove` (LH-FA-DEV-010): the built
// binary removes the containers of generated projects against the engine of
// the runner, and every case reads the state of the engine before and after
// the call: the containers of the project folder and whether they run, the
// volumes by name and by content, the images of the devcontainer CLI and the
// files of the project folder. The question of `--volumes` is answered through
// a pipe (no terminal) and under `script` (a pseudo-terminal). Zusage-Menge: a
// Linux host (`amd64`) with a local engine and Docker; the cases of the
// sandbox profile run with its default settings, those of the feature `ollama`
// with the feature switched on, the standard profile without a catalogue
// feature. A missing docker binary or a missing `script` fails the test, it
// does not skip it.

package e2e_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"
)

// rmAnonymousName is the name of an anonymous volume: 64 hexadecimal characters.
var rmAnonymousName = regexp.MustCompile(`^[0-9a-f]{64}$`)

// rmEnvelope is the envelope of `devcontainer remove --json`.
type rmEnvelope struct {
	Status      string `json:"status"`
	Command     string `json:"command"`
	Subcommand  string `json:"subcommand"`
	ExitCode    int    `json:"exitCode"`
	Diagnostics []struct {
		Level   string `json:"level"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"diagnostics"`
	Data *struct {
		ContainerID        string   `json:"containerId"`
		Stopped            bool     `json:"stopped"`
		Removed            bool     `json:"removed"`
		RemovedVolumeNames []string `json:"removedVolumeNames"`
	} `json:"data"`
}

// rmRun runs the binary in dir with the text on its standard input (a pipe, no
// terminal); env entries replace those of the test process.
func rmRun(ctx context.Context, t *testing.T, bin, dir string, env []string, stdin string, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdin = strings.NewReader(stdin)
	var out, errOut lockedBuffer
	cmd.Stdout, cmd.Stderr = &out, &errOut
	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if !errors.As(err, &ee) {
			t.Fatalf("run %v: %v", args, err)
		}
		code = ee.ExitCode()
	}
	return out.String(), errOut.String(), code
}

// rmRunTerminal runs the binary under `script`, so its standard input is a
// pseudo-terminal, and types answer when the question is on the screen. The
// output of both streams is the terminal's.
func rmRunTerminal(ctx context.Context, t *testing.T, bin, dir, answer string, args ...string) (string, int) {
	t.Helper()
	if _, err := exec.LookPath("script"); err != nil {
		t.Fatalf("`script` not on PATH but the terminal test needs it: %v", err)
	}
	quoted := []string{dcQuote(bin)}
	for _, a := range args {
		quoted = append(quoted, dcQuote(a))
	}
	cmd := exec.CommandContext(ctx, "script", "-q", "-e", "-c", strings.Join(quoted, " "), "/dev/null")
	cmd.Dir = dir
	cmd.Env = os.Environ()
	var out lockedBuffer
	cmd.Stdout, cmd.Stderr = &out, &out
	in, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	deadline := time.After(3 * time.Minute)
	asked := false
	for !asked {
		select {
		case err := <-done:
			done <- err
			asked = true
		case <-deadline:
			_ = cmd.Process.Kill()
			t.Fatalf("no question on the terminal; output:\n%s", out.String())
		case <-time.After(100 * time.Millisecond):
			asked = strings.Contains(out.String(), "[y/N]")
		}
	}
	_, _ = io.WriteString(in, answer+"\n")
	var waitErr error
	select {
	case waitErr = <-done:
	case <-time.After(3 * time.Minute):
		_ = cmd.Process.Kill()
		t.Fatalf("the command did not end; output:\n%s", out.String())
	}
	code := 0
	var ee *exec.ExitError
	if errors.As(waitErr, &ee) {
		code = ee.ExitCode()
	} else if waitErr != nil {
		t.Fatalf("wait: %v", waitErr)
	}
	return out.String(), code
}

// rmVolumesOf returns the volumes the container mounts, sorted (the Lesestand
// of the set `remove` removes).
func rmVolumesOf(t *testing.T, id string) []string {
	t.Helper()
	var names []string
	for _, m := range dcInspect(t, id).Mounts {
		if m.Type == "volume" {
			names = append(names, m.Name)
		}
	}
	sort.Strings(names)
	return names
}

func rmVolumeExists(name string) bool {
	_, err := dcDocker("volume", "inspect", name)
	return err == nil
}

// rmKeepVolumes removes the volumes at the end of the test: the cleanup of the
// project finds only the volumes of containers that still exist.
func rmKeepVolumes(t *testing.T, names ...string) {
	t.Helper()
	t.Cleanup(func() {
		for _, n := range names {
			_, _ = dcDocker("volume", "rm", "-f", n)
		}
	})
}

// rmVolumeFile reads a file of the volume with a container that mounts it
// read-only; ok is false when the file is not there.
func rmVolumeFile(ctx context.Context, t *testing.T, volume, name string) (string, bool) {
	t.Helper()
	out, err := exec.CommandContext(ctx, "docker", "run", "--rm", "-u", "0", "--entrypoint", "sh",
		"-v", volume+":/v:ro", dcBaseImage, "-c", "cat /v/"+name).Output()
	return strings.TrimSpace(string(out)), err == nil
}

// rmTree is the Lesestand of the project folder: every path with a hash of the
// content (or "dir") and the mode.
func rmTree(t *testing.T, dir string) map[string]string {
	t.Helper()
	tree := map[string]string{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		if d.IsDir() {
			tree[rel] = "dir " + info.Mode().String()
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(body)
		tree[rel] = hex.EncodeToString(sum[:]) + " " + info.Mode().String()
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return tree
}

// rmState is the Lesestand of the containers and volumes a case watches.
type rmState struct {
	containers map[string]string // label of the container -> ID
	volumes    map[string]string // label of the volume -> name
}

// require fails unless every container of the state exists and has the
// running state in running (default true), and every volume exists.
func (s rmState) require(t *testing.T, step string, running map[string]bool, skipContainers, skipVolumes map[string]bool) {
	t.Helper()
	for label, id := range s.containers {
		if skipContainers[label] {
			continue
		}
		if !dcExists(id) {
			t.Errorf("%s: container %s (%s) is gone", step, label, id[:12])
			continue
		}
		want := true
		if w, ok := running[label]; ok {
			want = w
		}
		if got := dcRunning(t, id); got != want {
			t.Errorf("%s: container %s running = %v, want %v", step, label, got, want)
		}
	}
	for label, name := range s.volumes {
		if !skipVolumes[label] && !rmVolumeExists(name) {
			t.Errorf("%s: volume %s (%s) is gone", step, label, name)
		}
	}
}

func rmDecode(t *testing.T, out string) rmEnvelope {
	t.Helper()
	var env rmEnvelope
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	if env.Command != "devcontainer" || env.Subcommand != "remove" {
		t.Errorf("envelope names %s %s, want devcontainer remove", env.Command, env.Subcommand)
	}
	return env
}

// LH-FA-DEV-010 (Spiegelstrich zu `remove`, Lastenheft 0.37.1 bis 0.37.3), sandbox
// profile with the default settings, the container without an instance and two
// instances of one project folder. The cases, each with the state of the
// engine before and after the call:
//
//   - invalid instance names (the empty string, `=`, a capital, 33 characters,
//     a hyphen at the start or the end) end with exit 10, a missing docker
//     client with exit 11, an instance without a container with exit 0 (32
//     characters pass); in each case every container still runs and every
//     volume stays;
//   - `--volumes` without `--yes` without a terminal ends with exit 10 after
//     the question named the container and its volume, with `--no-interactive`
//     and with `--json` with exit 10 without a question; a declined question
//     under a terminal ends with exit 10; none of them stops or removes
//     anything;
//   - `remove --instance a` without a container of the instance, next to
//     other containers, ends with exit 0 and leaves them alone;
//   - `remove --instance i1` stops and removes only that container, its
//     volume stays with its content, a following `remove --volumes --yes`
//     finds no container, ends with exit 0 and leaves the volume, and `up`
//     mounts the volume again in a new container;
//   - `remove --volumes --yes --instance i2` removes exactly the volumes the
//     Lesestand names, leaves the others, and `up` creates the volume again,
//     empty;
//   - `remove --volumes --yes` without `--instance` removes only the container
//     without an instance and its volume; with only containers of instances
//     left, running, mixed and stopped, it ends with exit 10 and a hint at
//     `--instance` and removes nothing;
//   - a stopped container is removed with its volume without being stopped.
func TestDevcontainerRemove_Docker_ContainersVolumesAndSelection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	requireDocker(t)
	dir := dcGenerateProject(ctx, t, true, nil)
	svc := dcService()
	bin := ubootBinary(t)
	ws := dcWorkspacePath(dir)

	// The container without an instance, then the instances.
	state := rmState{containers: map[string]string{}, volumes: map[string]string{}}
	for _, name := range []string{"", "i1", "i2"} {
		label := name
		if label == "" {
			label = "base"
		}
		resp := dcUpInstance(ctx, t, svc, dir, name, false)
		state.containers[label] = resp.ContainerID
		state.volumes[label] = dcInspect(t, resp.ContainerID).volumeAt(ws)
		if out, code := dcExec(ctx, t, resp.ContainerID, "root", "echo "+label+" > "+ws+"/marker-"+label); code != 0 {
			t.Fatalf("write the marker of %s: %s", label, out)
		}
	}
	dcRequireDistinct(t, "workspace volume", []string{state.volumes["base"], state.volumes["i1"], state.volumes["i2"]})
	rmKeepVolumes(t, state.volumes["base"], state.volumes["i1"], state.volumes["i2"])
	images := dcProjectImages(t, dir)
	short := func(label string) string { return state.containers[label][:12] }
	state.require(t, "before", nil, nil, nil)

	// Invalid instance names, no docker client, an instance without a container.
	for _, name := range []string{"", "a=b", "Abc", strings.Repeat("a", 33), "-a", "a-"} {
		_, errOut, code := rmRun(ctx, t, bin, dir, nil, "", "devcontainer", "remove", "--instance="+name, "--volumes", "--yes")
		if code != 10 || !strings.Contains(errOut, "instance") {
			t.Errorf("--instance=%q: exit %d, stderr %q, want exit 10 naming the instance", name, code, errOut)
		}
	}
	state.require(t, "invalid instance names", nil, nil, nil)
	if _, errOut, code := rmRun(ctx, t, bin, dir, []string{"PATH=" + t.TempDir()}, "", "devcontainer", "remove", "--volumes", "--yes"); code != 11 {
		t.Errorf("no docker client: exit %d, stderr %q, want 11", code, errOut)
	}
	state.require(t, "no docker client", nil, nil, nil)
	if out, _, code := rmRun(ctx, t, bin, dir, nil, "", "devcontainer", "remove", "--instance", strings.Repeat("a", 32), "--volumes"); code != 0 ||
		!strings.Contains(out, "nothing to remove") {
		t.Errorf("32 characters: exit %d, stdout %q, want 0 and nothing to remove", code, out)
	}
	state.require(t, "32 characters", nil, nil, nil)

	// --volumes without a release.
	out, errOut, code := rmRun(ctx, t, bin, dir, nil, "", "devcontainer", "remove", "--volumes")
	if code != 10 || out != "" {
		t.Errorf("no terminal: exit %d, stdout %q, want 10 and no stdout", code, out)
	}
	for _, part := range []string{"About to remove the devcontainer container " + short("base"), "  - " + state.volumes["base"], "PERMANENTLY DELETED", "[y/N]", "neither stopped nor removed"} {
		if !strings.Contains(errOut, part) {
			t.Errorf("no terminal: stderr lacks %q:\n%s", part, errOut)
		}
	}
	for _, other := range []string{state.volumes["i1"], state.volumes["i2"]} {
		if strings.Contains(errOut, other) {
			t.Errorf("the question names the volume %s of another container:\n%s", other, errOut)
		}
	}
	state.require(t, "no terminal", nil, nil, nil)
	_, errOut, code = rmRun(ctx, t, bin, dir, nil, "", "devcontainer", "remove", "--volumes", "--no-interactive")
	if code != 10 || strings.Contains(errOut, "[y/N]") || !strings.Contains(errOut, "--yes") {
		t.Errorf("--no-interactive: exit %d, stderr %q, want 10, no question and a hint at --yes", code, errOut)
	}
	jsonOut, _, code := rmRun(ctx, t, bin, dir, nil, "", "devcontainer", "remove", "--volumes", "--json")
	env := rmDecode(t, jsonOut)
	if code != 10 || env.ExitCode != 10 || env.Status != "error" || env.Data != nil || len(env.Diagnostics) != 1 ||
		env.Diagnostics[0].Code != "LH-FA-DEV-010" || !strings.Contains(env.Diagnostics[0].Message, "--yes") {
		t.Errorf("--json without --yes: exit %d, envelope %+v", code, env)
	}
	state.require(t, "--no-interactive and --json", nil, nil, nil)
	if _, _, code := rmRun(ctx, t, bin, dir, nil, "", "--yes", "--no-interactive", "devcontainer", "remove"); code != 2 {
		t.Errorf("--yes with --no-interactive: exit %d, want 2", code)
	}
	state.require(t, "--yes with --no-interactive", nil, nil, nil)

	// A declined question under a terminal.
	tty, code := rmRunTerminal(ctx, t, bin, dir, "n", "devcontainer", "remove", "--volumes")
	if code != 10 || !strings.Contains(tty, "About to remove the devcontainer container "+short("base")) || !strings.Contains(tty, state.volumes["base"]) {
		t.Errorf("declined under a terminal: exit %d, output:\n%s", code, tty)
	}
	state.require(t, "declined question", nil, nil, nil)

	// An instance without a container, next to other containers.
	if out, _, code := rmRun(ctx, t, bin, dir, nil, "", "devcontainer", "remove", "--instance", "a", "--volumes", "--yes"); code != 0 ||
		!strings.Contains(out, "No devcontainer container found for this project; nothing to remove") {
		t.Errorf("instance a: exit %d, stdout %q", code, out)
	}
	state.require(t, "instance without a container", nil, nil, nil)

	// remove --instance i1: the container goes, the volume and its content stay.
	oldI1 := state.containers["i1"]
	out, errOut, code = rmRun(ctx, t, bin, dir, nil, "", "devcontainer", "remove", "--instance", "i1")
	if code != 0 || out != "Devcontainer removed: container "+short("i1")+" (stopped first)\n" || errOut != "" {
		t.Errorf("remove i1: exit %d, stdout %q, stderr %q", code, out, errOut)
	}
	if dcExists(oldI1) {
		t.Errorf("remove i1: the container %s still exists", oldI1[:12])
	}
	state.require(t, "remove i1", nil, map[string]bool{"i1": true}, nil)
	if got, ok := rmVolumeFile(ctx, t, state.volumes["i1"], "marker-i1"); !ok || got != "i1" {
		t.Errorf("remove i1: the content of the volume is %q (found %v), want the marker", got, ok)
	}
	if out, _, code := rmRun(ctx, t, bin, dir, nil, "", "devcontainer", "remove", "--instance", "i1", "--volumes", "--yes"); code != 0 ||
		out != "No devcontainer container found for this project; nothing to remove\nVolumes of a container that no longer exists are not removed by u-boot.\n" {
		t.Errorf("remove --volumes of i1 without a container: exit %d, stdout %q", code, out)
	}
	if !rmVolumeExists(state.volumes["i1"]) {
		t.Errorf("remove --volumes without a container removed the volume of i1")
	}
	again := dcUpInstance(ctx, t, svc, dir, "i1", false)
	if again.ContainerID == oldI1 || dcInspect(t, again.ContainerID).volumeAt(ws) != state.volumes["i1"] {
		t.Errorf("up i1 after remove: container %s (old %s), volume %q, want a new container with the old volume", again.ContainerID[:12], oldI1[:12], dcInspect(t, again.ContainerID).volumeAt(ws))
	}
	state.containers["i1"] = again.ContainerID
	if got, code := dcExec(ctx, t, again.ContainerID, "root", "cat "+ws+"/marker-i1"); code != 0 || strings.TrimSpace(got) != "i1" {
		t.Errorf("up i1 after remove: marker %q (exit %d), want the old content", got, code)
	}
	dcRequireImagesKept(t, "remove i1", images, dcProjectImages(t, dir))

	// remove --volumes --yes --instance i2: exactly the volumes of the Lesestand.
	wantVolumes := rmVolumesOf(t, state.containers["i2"])
	if len(wantVolumes) == 0 || wantVolumes[0] != state.volumes["i2"] {
		t.Fatalf("Lesestand of the volumes of i2 = %v, want the workspace volume %s", wantVolumes, state.volumes["i2"])
	}
	oldI2 := state.containers["i2"]
	jsonOut, _, code = rmRun(ctx, t, bin, dir, nil, "", "devcontainer", "remove", "--instance", "i2", "--volumes", "--yes", "--json")
	env = rmDecode(t, jsonOut)
	if code != 0 || env.Status != "ok" || env.Data == nil || env.Data.ContainerID != oldI2 || !env.Data.Stopped || !env.Data.Removed ||
		strings.Join(env.Data.RemovedVolumeNames, ",") != strings.Join(wantVolumes, ",") {
		t.Errorf("remove --volumes i2: exit %d, envelope %+v, want the volumes %v", code, env, wantVolumes)
	}
	if dcExists(oldI2) {
		t.Errorf("remove --volumes i2: the container still exists")
	}
	for _, v := range wantVolumes {
		if rmVolumeExists(v) {
			t.Errorf("remove --volumes i2: the volume %s still exists", v)
		}
	}
	state.require(t, "remove --volumes i2", nil, map[string]bool{"i2": true}, map[string]bool{"i2": true})
	again2 := dcUpInstance(ctx, t, svc, dir, "i2", false)
	state.containers["i2"] = again2.ContainerID
	if _, code := dcExec(ctx, t, again2.ContainerID, "root", "test -e "+ws+"/marker-i2"); code == 0 {
		t.Errorf("up i2 after remove --volumes: the marker of the old volume is there, want a new empty volume")
	}
	if dcInspect(t, again2.ContainerID).volumeAt(ws) != state.volumes["i2"] {
		t.Errorf("up i2 after remove --volumes: the volume is not the one of the same name")
	}

	// remove --volumes --yes without --instance: only the container without an instance.
	baseVolumes := rmVolumesOf(t, state.containers["base"])
	out, errOut, code = rmRun(ctx, t, bin, dir, nil, "", "devcontainer", "remove", "--volumes", "--yes")
	wantOut := "Devcontainer removed: container " + short("base") + " (stopped first)\nVolumes removed: 1\n  - " + baseVolumes[0] + "\n"
	if code != 0 || out != wantOut || errOut != "" {
		t.Errorf("remove --volumes without --instance: exit %d, stdout %q, stderr %q, want %q", code, out, errOut, wantOut)
	}
	if dcExists(state.containers["base"]) || rmVolumeExists(state.volumes["base"]) {
		t.Errorf("remove --volumes without --instance left the container or its volume")
	}
	state.require(t, "remove without --instance", nil, map[string]bool{"base": true}, map[string]bool{"base": true})

	// Only containers of instances: exit 10 and a hint, whatever their state.
	onlyInstances := func(step string, running map[string]bool) {
		t.Helper()
		_, errOut, code := rmRun(ctx, t, bin, dir, nil, "", "devcontainer", "remove", "--volumes", "--yes")
		if code != 10 || !strings.Contains(errOut, "--instance") {
			t.Errorf("%s: exit %d, stderr %q, want 10 and a hint at --instance", step, code, errOut)
		}
		state.require(t, step, running, map[string]bool{"base": true}, map[string]bool{"base": true})
	}
	onlyInstances("only instances, running", nil)
	dcMustDocker(t, "stop", state.containers["i1"])
	onlyInstances("only instances, mixed", map[string]bool{"i1": false})
	dcMustDocker(t, "stop", state.containers["i2"])
	onlyInstances("only instances, stopped", map[string]bool{"i1": false, "i2": false})

	// Stopped containers: removed without a stop, with their volumes.
	for _, name := range []string{"i1", "i2"} {
		volumes := rmVolumesOf(t, state.containers[name])
		out, _, code := rmRun(ctx, t, bin, dir, nil, "", "devcontainer", "remove", "--instance", name, "--volumes", "--yes")
		want := "Devcontainer removed: container " + short(name) + "\nVolumes removed: 1\n  - " + volumes[0] + "\n"
		if code != 0 || out != want {
			t.Errorf("remove stopped %s: exit %d, stdout %q, want %q", name, code, out, want)
		}
		if dcExists(state.containers[name]) || rmVolumeExists(state.volumes[name]) {
			t.Errorf("remove stopped %s left the container or its volume", name)
		}
	}
	if left := dcFolderContainers(t, dir); len(left) != 0 {
		t.Errorf("containers of the project after all: %v, want none", left)
	}
	dcRequireImagesKept(t, "all removals", images, dcProjectImages(t, dir))
}

// LH-FA-DEV-010, Lastenheft 0.37.2 and 0.37.3: without a container to remove
// `remove` ends with exit 0 and asks nothing, also `--volumes` without `--yes`
// without a terminal and with `--no-interactive` (and not with exit 10); with
// `--json` the envelope names no container and an empty array of volumes. A
// missing docker client ends with exit 11, the situation of the released
// container image.
func TestDevcontainerRemove_Docker_NoContainer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	requireDocker(t)
	dir := dcProjectDir(t)
	bin := ubootBinary(t)
	if left := dcFolderContainers(t, dir); len(left) != 0 {
		t.Fatalf("the fresh project folder has containers: %v", left)
	}
	before := rmTree(t, dir)

	for _, args := range [][]string{
		{"devcontainer", "remove"},
		{"devcontainer", "remove", "--volumes"},
		{"devcontainer", "remove", "--volumes", "--no-interactive"},
		{"--no-interactive", "devcontainer", "remove", "--volumes"},
	} {
		out, errOut, code := rmRun(ctx, t, bin, dir, nil, "", args...)
		want := "No devcontainer container found for this project; nothing to remove\n"
		if strings.Contains(strings.Join(args, " "), "--volumes") {
			want += "Volumes of a container that no longer exists are not removed by u-boot.\n"
		}
		if code != 0 || out != want || errOut != "" {
			t.Errorf("%v: exit %d, stdout %q, stderr %q, want exit 0, %q and no stderr", args, code, out, errOut, want)
		}
	}
	jsonOut, _, code := rmRun(ctx, t, bin, dir, nil, "", "devcontainer", "remove", "--volumes", "--json")
	env := rmDecode(t, jsonOut)
	if code != 0 || env.Status != "ok" || env.Data == nil || env.Data.ContainerID != "" || env.Data.Removed || env.Data.Stopped ||
		env.Data.RemovedVolumeNames == nil || len(env.Data.RemovedVolumeNames) != 0 || len(env.Diagnostics) != 0 {
		t.Errorf("--json: exit %d, envelope %+v", code, env)
	}
	if _, errOut, code := rmRun(ctx, t, bin, dir, []string{"PATH=" + t.TempDir()}, "", "devcontainer", "remove"); code != 11 {
		t.Errorf("no docker client: exit %d, stderr %q, want 11", code, errOut)
	}
	if !equalTrees(before, rmTree(t, dir)) {
		t.Errorf("the project folder changed")
	}
}

func equalTrees(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// LH-FA-DEV-010, standard profile without a catalogue feature and without a
// `VOLUME` in the Dockerfile: the container holds no volume, only the bind
// mount of the project folder. The question says so, `remove --volumes`
// removes the container, names no volume, and the project folder, its files
// and the images of the devcontainer CLI are as before. With a `VOLUME` in the
// Dockerfile of the project the container mounts an anonymous volume: `remove`
// without `--volumes` leaves it, `remove --volumes --yes` on the stopped
// container (`docker rm -v` without `-f`) removes it.
func TestDevcontainerRemove_Docker_StandardProfileAndAnonymousVolume(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	requireDocker(t)
	svc := dcService()
	bin := ubootBinary(t)

	// No volume at all.
	dir := dcGenerateProject(ctx, t, false, nil)
	first := dcMustUp(ctx, t, svc, dir, false).resp.ContainerID
	if got := rmVolumesOf(t, first); len(got) != 0 {
		t.Fatalf("the standard profile mounts the volumes %v, want none", got)
	}
	if !dcInspect(t, first).bindsFolder(dir) {
		t.Fatalf("the container does not bind the project folder")
	}
	images, tree := dcProjectImages(t, dir), rmTree(t, dir)
	tty, code := rmRunTerminal(ctx, t, bin, dir, "y", "devcontainer", "remove", "--volumes")
	if code != 0 || !strings.Contains(tty, "About to remove the devcontainer container "+first[:12]+" of this project; it holds no volume.") ||
		!strings.Contains(tty, "No volume removed (the container held no volume)") {
		t.Errorf("confirmed under a terminal: exit %d, output:\n%s", code, tty)
	}
	if dcExists(first) {
		t.Errorf("the container still exists after the confirmed removal")
	}
	if !equalTrees(tree, rmTree(t, dir)) {
		t.Errorf("the project folder changed")
	}
	dcRequireImagesKept(t, "standard profile", images, dcProjectImages(t, dir))
	if next := dcMustUp(ctx, t, svc, dir, false).resp.ContainerID; next == first {
		t.Errorf("up after remove found the old container %s", first[:12])
	}

	// An anonymous volume.
	anon := dcGenerateProject(ctx, t, false, nil)
	dockerfile := filepath.Join(anon, ".devcontainer", "Dockerfile")
	f, err := os.OpenFile(dockerfile, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("\nVOLUME /rm-anon-data\n"); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	id := dcMustUp(ctx, t, svc, anon, false).resp.ContainerID
	volumes := rmVolumesOf(t, id)
	if len(volumes) != 1 || !rmAnonymousName.MatchString(volumes[0]) {
		t.Fatalf("the container mounts %v, want one anonymous volume", volumes)
	}
	rmKeepVolumes(t, volumes[0])
	if out, _, code := rmRun(ctx, t, bin, anon, nil, "", "devcontainer", "remove"); code != 0 || strings.Contains(out, "Volumes removed") {
		t.Errorf("remove without --volumes: exit %d, stdout %q", code, out)
	}
	if dcExists(id) || !rmVolumeExists(volumes[0]) {
		t.Errorf("remove without --volumes: container exists %v, anonymous volume exists %v, want the container gone and the volume kept", dcExists(id), rmVolumeExists(volumes[0]))
	}

	id = dcMustUp(ctx, t, svc, anon, false).resp.ContainerID
	volumes = rmVolumesOf(t, id)
	if len(volumes) != 1 || !rmAnonymousName.MatchString(volumes[0]) {
		t.Fatalf("the second container mounts %v, want one anonymous volume", volumes)
	}
	rmKeepVolumes(t, volumes[0])
	dcMustDocker(t, "stop", id)
	anonTree := rmTree(t, anon)
	jsonOut, _, code := rmRun(ctx, t, bin, anon, nil, "", "devcontainer", "remove", "--volumes", "--yes", "--json")
	env := rmDecode(t, jsonOut)
	if code != 0 || env.Data == nil || env.Data.ContainerID != id || env.Data.Stopped || !env.Data.Removed ||
		strings.Join(env.Data.RemovedVolumeNames, ",") != volumes[0] {
		t.Errorf("remove --volumes of the stopped container: exit %d, envelope %+v, want the volume %s", code, env, volumes[0])
	}
	if dcExists(id) || rmVolumeExists(volumes[0]) {
		t.Errorf("remove --volumes of the stopped container: container exists %v, anonymous volume exists %v, want both gone", dcExists(id), rmVolumeExists(volumes[0]))
	}
	if !equalTrees(anonTree, rmTree(t, anon)) {
		t.Errorf("the project folder changed")
	}
}

// LH-FA-DEV-010, feature `ollama` and a volume that another container holds:
// without `--volumes` the volume of the models stays and a new container mounts
// it again under the same name; with `--volumes --yes` it goes with the
// container. A second container that holds the workspace volume keeps that
// volume: `docker volume rm` fails while a container holds it, `remove` still
// removes the container and the volume of the models, ends with exit 12 and
// names the volume that remains and the one removed.
func TestDevcontainerRemove_Docker_OllamaVolumeAndVolumeHeldByAnotherContainer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	requireDocker(t)
	dir := dcGenerateProject(ctx, t, true, map[string]string{"devcontainer.features.ollama.enabled": "true"})
	svc := dcService()
	bin := ubootBinary(t)
	ws := dcWorkspacePath(dir)

	first := dcUpInstance(ctx, t, svc, dir, "", false).ContainerID
	volumes := rmVolumesOf(t, first)
	workspace := dcInspect(t, first).volumeAt(ws)
	var models string
	for _, v := range volumes {
		if strings.Contains(v, "-ollama-") {
			models = v
		}
	}
	if len(volumes) != 2 || workspace == "" || models == "" {
		t.Fatalf("the container mounts %v (workspace %q), want the workspace volume and the volume of the models", volumes, workspace)
	}
	rmKeepVolumes(t, volumes...)

	// Without --volumes both stay and the next container mounts them again.
	if _, _, code := rmRun(ctx, t, bin, dir, nil, "", "devcontainer", "remove"); code != 0 || dcExists(first) {
		t.Fatalf("remove: exit %d, container exists %v", code, dcExists(first))
	}
	for _, v := range volumes {
		if !rmVolumeExists(v) {
			t.Errorf("remove without --volumes removed the volume %s", v)
		}
	}
	second := dcUpInstance(ctx, t, svc, dir, "", false).ContainerID
	if got := rmVolumesOf(t, second); strings.Join(got, ",") != strings.Join(volumes, ",") {
		t.Errorf("the new container mounts %v, want the old volumes %v", got, volumes)
	}

	// Another container holds the workspace volume.
	holder := "rm-holder-" + dcSuffix()
	t.Cleanup(func() { _, _ = dcDocker("rm", "-f", holder) })
	dcMustDocker(t, "run", "-d", "--name", holder, "-v", workspace+":/held", "--entrypoint", "sleep", dcBaseImage, "600")

	jsonOut, _, code := rmRun(ctx, t, bin, dir, nil, "", "devcontainer", "remove", "--volumes", "--yes", "--json")
	env := rmDecode(t, jsonOut)
	if code != 12 || env.ExitCode != 12 || env.Status != "error" || env.Data != nil || len(env.Diagnostics) != 1 || env.Diagnostics[0].Code != "LH-FA-DEV-010" {
		t.Fatalf("remove --volumes with a held volume: exit %d, envelope %+v, want exit 12", code, env)
	}
	message := env.Diagnostics[0].Message
	for _, part := range []string{second[:12], "volumes removed: " + models, "volumes that remain: " + workspace, "in use"} {
		if !strings.Contains(message, part) {
			t.Errorf("the message lacks %q:\n%s", part, message)
		}
	}
	if dcExists(second) {
		t.Errorf("the container still exists after the failed volume removal")
	}
	if rmVolumeExists(models) {
		t.Errorf("remove --volumes left the volume of the models")
	}
	if !rmVolumeExists(workspace) {
		t.Errorf("the volume held by another container is gone")
	}
	if !dcRunning(t, holder) {
		t.Errorf("the holder container does not run any more")
	}
}
