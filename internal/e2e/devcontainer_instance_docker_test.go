//go:build docker

// End-to-end check of `--instance` of `u-boot devcontainer up`, `down` and
// `shell` (LH-FA-DEV-010, ADR-0038): several containers of one project folder
// run side by side, started with the production start service and, for the
// exit codes and the shell, with the built binary, against the engine of the
// runner. The values the tests pin are those of the measurement of the
// instances (docs/messungen/slice-v1-devcontainer-instanz-messung): a deviation
// is a finding for the measurement, not a reason to adapt the test.
//
// The sandbox test needs the host the commitment names: Linux amd64, a local
// engine with memory and swap limits, and room for four instances of up to
// one CPU and 3 GiB each. It fails with a finding when the host does not carry
// the stage; it does not leave the stage out.

package e2e_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const (
	// dcInstanceCount is the number of instances of the sandbox test: the
	// highest number the commitment names.
	dcInstanceCount = 4
	// dcInstanceLabel is the id label of the container of an instance.
	dcInstanceLabel = "u-boot.instance"
	// The memory the host must offer for four instances (MemAvailable, bytes):
	// the stage the measurement ran under.
	dcInstanceMinMemAvailable = 9 << 30
	dcInstanceMinSwapFree     = 8 << 30
	// The limit of each instance: the condition the commitment names.
	dcInstanceMemory  = 3 << 30
	dcInstanceNanoCPU = 1_000_000_000
	// dcInstanceMinFreeBytes is the free space of the Docker root the stage needs.
	dcInstanceMinFreeBytes = 40 << 30

	// The three volume families of a sandbox container with nestedRuntime
	// docker and the feature ollama, by their mount point.
	dcWorkspaceDest = "/workspaces/dcstart"
	dcDockerDest    = "/var/lib/docker"
	dcOllamaDest    = "/home/vscode/.ollama"
)

// dcMeminfoBytes reads a field of /proc/meminfo in bytes.
func dcMeminfoBytes(t *testing.T, field string) int64 {
	t.Helper()
	body, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		t.Fatalf("FINDING: the memory of the host is not readable: %v", err)
	}
	for _, line := range strings.Split(string(body), "\n") {
		if rest, ok := strings.CutPrefix(line, field+":"); ok {
			kb, err := strconv.ParseInt(strings.Fields(rest)[0], 10, 64)
			if err != nil {
				t.Fatalf("FINDING: /proc/meminfo %s: %v", field, err)
			}
			return kb << 10
		}
	}
	t.Fatalf("FINDING: /proc/meminfo has no field %s", field)
	return 0
}

// dcRequireInstanceHost fails with a finding when the host does not carry
// four instances under the limit of the commitment: the engine must enforce
// memory and swap limits and offer at least as many CPUs as instances, the
// memory must be available, and the Docker root must have room when the test
// can see it.
func dcRequireInstanceHost(t *testing.T) {
	t.Helper()
	requireDocker(t)
	info := dcMustDocker(t, "info", "--format", "{{.MemoryLimit}} {{.SwapLimit}} {{.NCPU}} {{.DockerRootDir}}")
	fields := strings.Fields(info)
	if len(fields) != 4 || fields[0] != "true" || fields[1] != "true" {
		t.Fatalf("FINDING: the engine does not enforce memory and swap limits (docker info: %q); the stage of %d instances is not carried", info, dcInstanceCount)
	}
	if cpus, err := strconv.Atoi(fields[2]); err != nil || cpus < dcInstanceCount {
		t.Fatalf("FINDING: the engine offers %q CPUs, the stage needs %d (one per instance)", fields[2], dcInstanceCount)
	}
	if avail := dcMeminfoBytes(t, "MemAvailable"); avail < dcInstanceMinMemAvailable {
		t.Fatalf("FINDING: MemAvailable %d MiB is below the %d MiB the stage of %d instances ran under", avail>>20, int64(dcInstanceMinMemAvailable)>>20, dcInstanceCount)
	}
	if free := dcMeminfoBytes(t, "SwapFree"); free < dcInstanceMinSwapFree {
		t.Fatalf("FINDING: SwapFree %d MiB is below the %d MiB the stage of %d instances ran under", free>>20, int64(dcInstanceMinSwapFree)>>20, dcInstanceCount)
	}
	var st syscall.Statfs_t
	if err := syscall.Statfs(fields[3], &st); err != nil {
		t.Logf("the free space of the Docker root %s is not readable from the test (%v); not checked", fields[3], err)
	} else if free := int64(st.Bavail) * int64(st.Bsize); free < dcInstanceMinFreeBytes {
		t.Fatalf("FINDING: %d MiB free in %s, the stage needs %d MiB", free>>20, fields[3], int64(dcInstanceMinFreeBytes)>>20)
	}
}

// dcLimitRunArgs adds the limit of the commitment (one CPU, 3 GiB with swap)
// to the runArgs of the generated devcontainer.json, so each container of the
// project runs under it.
func dcLimitRunArgs(t *testing.T, dir string) {
	t.Helper()
	path := filepath.Join(dir, ".devcontainer", "devcontainer.json")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	const head = `"runArgs": [`
	if !strings.Contains(string(body), head) {
		t.Fatalf("devcontainer.json has no runArgs:\n%s", body)
	}
	limited := strings.Replace(string(body), head, head+`"--memory=3g", "--memory-swap=3g", "--cpus=1", `, 1)
	if err := os.WriteFile(path, []byte(limited), 0o644); err != nil {
		t.Fatal(err)
	}
}

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

// dcMarkerPaths are the three files a marker of an instance lands in, one per
// volume family.
func dcMarkerPaths(instance string) []string {
	return []string{
		dcWorkspaceDest + "/marker-" + instance,
		dcDockerDest + "/marker-" + instance,
		dcOllamaDest + "/marker-" + instance,
	}
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

// LH-FA-DEV-010 (ADR-0038), sandbox profile with nestedRuntime docker and the
// feature ollama, four instances of one project folder, each under one CPU and
// 3 GiB: the values of the measurement.
//
//   - The containers of the instances carry the folder label and the instance
//     label and no `devcontainer.config_file`; the container without an
//     instance carries the folder label and the config label and no instance
//     label. Container IDs and the volumes of the three families (workspace,
//     Docker storage, models) are pairwise different, also against the
//     container without an instance. The limit is read back.
//   - A file in each of the three volumes of an instance is in no other
//     container (every ordered pair).
//   - `up`, `down` and `shell` without `--instance` meet the container without
//     an instance although the instances are newer; a second `up --instance`
//     without `--recreate` meets the same container; `down --instance` stops
//     only that instance; `up --instance --recreate` replaces only that
//     container, the volumes and the file stay.
//   - `down --instance <unknown>` stops nothing and ends with 0.
//   - Without a container without an instance, `up`, `down` and `shell`
//     without `--instance` end with exit 10 and a hint at `--instance`, with
//     no helper container and no change to the containers.
func TestDevcontainerInstance_Docker_FourInstancesAreIsolated(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*dcTestBudget)
	defer cancel()
	dcRequireInstanceHost(t)
	dir := dcGenerateProject(ctx, t, true, map[string]string{
		"devcontainer.sandbox.nestedRuntime":   "docker",
		"devcontainer.sandbox.onUnavailable":   "warn",
		"devcontainer.features.ollama.enabled": "true",
	})
	dcLimitRunArgs(t, dir)
	svc := dcService()
	bin := ubootBinary(t)

	// The container without an instance, then the instances, each newer.
	existing := dcUpInstance(ctx, t, svc, dir, "", false)
	names := make([]string, dcInstanceCount)
	ids := make([]string, 0, dcInstanceCount+1)
	ids = append(ids, existing.ContainerID)
	for i := range names {
		names[i] = fmt.Sprintf("i%d", i+1)
		resp := dcUpInstance(ctx, t, svc, dir, names[i], false)
		if resp.AlreadyRunning {
			t.Errorf("the first start of instance %s reports already running", names[i])
		}
		ids = append(ids, resp.ContainerID)
	}
	t.Logf("containers: without instance %s, instances %v", existing.ContainerID[:12], ids[1:])

	// Identity: IDs, labels, volumes, limit.
	dcRequireDistinct(t, "container ID", ids)
	var workspaces, dockers, ollamas []string
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
		if doc.HostConfig.Memory != dcInstanceMemory || doc.HostConfig.MemorySwap != dcInstanceMemory || doc.HostConfig.NanoCpus != dcInstanceNanoCPU {
			t.Errorf("container %s: limit memory %d swap %d nanoCpus %d, want 3 GiB with swap and one CPU",
				id[:12], doc.HostConfig.Memory, doc.HostConfig.MemorySwap, doc.HostConfig.NanoCpus)
		}
		workspaces = append(workspaces, doc.volumeAt(dcWorkspaceDest))
		dockers = append(dockers, doc.volumeAt(dcDockerDest))
		ollamas = append(ollamas, doc.volumeAt(dcOllamaDest))
	}
	dcRequireDistinct(t, "workspace volume", workspaces)
	dcRequireDistinct(t, "Docker-storage volume", dockers)
	dcRequireDistinct(t, "model volume", ollamas)

	// Content probe over all ordered pairs.
	for n, id := range ids {
		instance := "base"
		if n > 0 {
			instance = names[n-1]
		}
		script := ""
		for _, p := range dcMarkerPaths(instance) {
			script += "echo " + instance + " > " + p + " && "
		}
		if out, code := dcExec(ctx, t, id, "root", script+"true"); code != 0 {
			t.Fatalf("write the markers of %s: %s", instance, out)
		}
	}
	for i := range ids {
		instI := "base"
		if i > 0 {
			instI = names[i-1]
		}
		for j, idJ := range ids {
			for _, p := range dcMarkerPaths(instI) {
				_, code := dcExec(ctx, t, idJ, "root", "test -e "+p)
				if (code == 0) != (i == j) {
					t.Errorf("file %s of %s: visible in container %d = %v, want %v", p, instI, j, code == 0, i == j)
				}
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
	if out, _, code := runBinary(ctx, t, bin, dir, "hostname\nexit\n", "devcontainer", "shell", "--instance", names[2]); code != 0 || strings.TrimSpace(out) != ids[3][:12] {
		t.Errorf("shell --instance %s: exit %d, hostname %q, want exit 0 and %s", names[2], code, strings.TrimSpace(out), ids[3][:12])
	}
	// A second start of an instance without --recreate meets the same container.
	repeat := dcUpInstance(ctx, t, svc, dir, names[1], false)
	if repeat.ContainerID != ids[2] || !repeat.AlreadyRunning {
		t.Errorf("second up --instance %s: container %s already running %v, want %s already running", names[1], repeat.ContainerID, repeat.AlreadyRunning, ids[2])
	}
	if got := dcFolderContainers(t, dir); len(got) != dcInstanceCount+1 {
		t.Errorf("containers of the project = %d, want %d", len(got), dcInstanceCount+1)
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
	// --recreate replaces only the container of the instance; volumes and file stay.
	recreated := dcUpInstance(ctx, t, svc, dir, names[0], true)
	if recreated.ContainerID == ids[1] || dcExists(ids[1]) {
		t.Errorf("recreate of %s: new ID %s, old container still there = %v", names[0], recreated.ContainerID, dcExists(ids[1]))
	}
	doc := dcInspect(t, recreated.ContainerID)
	if doc.volumeAt(dcWorkspaceDest) != workspaces[1] || doc.volumeAt(dcDockerDest) != dockers[1] || doc.volumeAt(dcOllamaDest) != ollamas[1] {
		t.Errorf("recreate of %s changed the volumes: %+v", names[0], doc.Mounts)
	}
	for _, p := range dcMarkerPaths(names[0]) {
		if out, code := dcExec(ctx, t, recreated.ContainerID, "root", "test -e "+p); code != 0 {
			t.Errorf("recreate of %s: file %s is gone: %s", names[0], p, out)
		}
	}
	for n, id := range ids[2:] {
		if !dcExists(id) || !dcRunning(t, id) {
			t.Errorf("recreate of %s touched instance %s (exists %v)", names[0], names[n+1], dcExists(id))
		}
	}
	ids[1] = recreated.ContainerID

	// Without a container without an instance: exit 10, no helper container,
	// no change to the containers.
	baseDoc := dcInspect(t, existing.ContainerID)
	baseVolumes := []string{baseDoc.volumeAt(dcWorkspaceDest), baseDoc.volumeAt(dcDockerDest), baseDoc.volumeAt(dcOllamaDest)}
	dcMustDocker(t, "rm", "-f", existing.ContainerID)
	t.Cleanup(func() {
		for _, v := range baseVolumes {
			if v != "" {
				_, _ = dcDocker("volume", "rm", "-f", v)
			}
		}
	})
	before := dcFolderContainers(t, dir)
	if len(before) != dcInstanceCount {
		t.Fatalf("containers of the project = %d, want only the %d instances", len(before), dcInstanceCount)
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
