//go:build docker

// Integration tests for the [docker.Launcher] adapter against the Docker
// engine of the runner (`make test-docker`): the guards of ADR-0028
// Entscheidung 9 and ADR-0029 Entscheidung 4 for path (S2), socket,
// environment and lifetime (S6), supply chain (S5) and the Linux engine
// (S7). The projects are hand-written `devcontainer.json` files without
// features on `mcr.microsoft.com/devcontainers/base:debian`; they live
// under the shared `/tmp`, where test container and Docker daemon see the
// same path. Every test removes the containers, volumes and images it
// creates, also on failure; the helper image is kept as the build cache.
// A missing docker binary fails the test, it does not skip it.

package docker_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/adapter/driven/docker"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

const (
	devBaseImage       = "mcr.microsoft.com/devcontainers/base:debian"
	folderLabel        = "devcontainer.local_folder"
	launcherTestBudget = 15 * time.Minute
)

func requireDocker(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH but the test was built with -tags=docker: %v", err)
	}
}

func randomHex(t *testing.T) string {
	t.Helper()
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(b)
}

// tryDocker runs docker and returns its trimmed stdout.
func tryDocker(args ...string) (string, error) {
	var stderr bytes.Buffer
	cmd := exec.Command("docker", args...)
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return strings.TrimSpace(string(out)), errors.New(err.Error() + ": " + strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(string(out)), nil
}

func mustDocker(t *testing.T, args ...string) string {
	t.Helper()
	out, err := tryDocker(args...)
	if err != nil {
		t.Fatalf("docker %s: %v", strings.Join(args, " "), err)
	}
	return out
}

// removeContainers removes the containers that match the `docker ps` filter.
func removeContainers(filter string) {
	ids, err := tryDocker("ps", "-aq", "--filter", filter)
	if err != nil || ids == "" {
		return
	}
	_, _ = tryDocker(append([]string{"rm", "-f"}, strings.Fields(ids)...)...)
}

func removeImages(reference string) {
	ids, err := tryDocker("images", "-q", "--filter", "reference="+reference)
	if err != nil || ids == "" {
		return
	}
	_, _ = tryDocker(append([]string{"rmi", "-f"}, strings.Fields(ids)...)...)
}

// helperPrefix is the helper-container name prefix of one test: the product
// prefix, a hash of the test name and a dash. A test starts, counts and
// removes only containers below its own prefix, so it neither removes nor
// counts the helper container of another run on the same daemon.
func helperPrefix(t *testing.T) string {
	t.Helper()
	sum := sha256.Sum256([]byte(t.Name()))
	return docker.HelperContainerPrefixForTest + "t" + hex.EncodeToString(sum[:4]) + "-"
}

// testLauncher is the production launcher with the helper-container name
// prefix of the test.
func testLauncher(t *testing.T) *docker.Launcher {
	t.Helper()
	return docker.NewLauncherWithNamePrefixForTest(helperPrefix(t))
}

// helperLeftovers lists the containers below the helper name prefix of the test.
func helperLeftovers(t *testing.T) []string {
	t.Helper()
	out := mustDocker(t, "ps", "-a", "--filter", "name="+helperPrefix(t), "--format", "{{.Names}}")
	return strings.Fields(out)
}

// newProject writes .devcontainer/devcontainer.json (and optional extra
// files) into a fresh folder and registers the cleanup of everything the
// devcontainer CLI creates for it.
func newProject(t *testing.T, config map[string]any, extra map[string]string) string {
	t.Helper()
	requireDocker(t)
	folder := filepath.Join(t.TempDir(), "ubl"+randomHex(t))
	if err := os.MkdirAll(filepath.Join(folder, ".devcontainer"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeConfig(t, folder, config)
	for name, content := range extra {
		if err := os.WriteFile(filepath.Join(folder, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		removeContainers("label=" + folderLabel + "=" + folder)
		removeContainers("name=" + helperPrefix(t))
		removeImages("vsc-" + filepath.Base(folder) + "-*")
	})
	return folder
}

func writeConfig(t *testing.T, folder string, config map[string]any) {
	t.Helper()
	raw, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(folder, ".devcontainer", "devcontainer.json"), raw, 0o644); err != nil {
		t.Fatal(err)
	}
}

func imageConfig(extra map[string]any) map[string]any {
	cfg := map[string]any{"name": "ubl-launcher-test", "image": devBaseImage}
	for k, v := range extra {
		cfg[k] = v
	}
	return cfg
}

// up runs the production launcher against the engine of the runner and
// returns the result and the progress text.
func up(t *testing.T, opts driven.DevcontainerUpOptions) (driven.DevcontainerUpResult, string, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), launcherTestBudget)
	defer cancel()
	var progress bytes.Buffer
	opts.Progress = &progress
	res, err := testLauncher(t).Up(ctx, opts)
	return res, progress.String(), err
}

func mustUp(t *testing.T, opts driven.DevcontainerUpOptions) driven.DevcontainerUpResult {
	t.Helper()
	res, progress, err := up(t, opts)
	if err != nil {
		t.Fatalf("Up: %v\nprogress:\n%s", err, progress)
	}
	return res
}

// inspectDoc is the part of `docker inspect` the tests read.
type inspectDoc struct {
	Mounts []struct {
		Source      string
		Destination string
	}
	HostConfig struct {
		Binds  []string
		Mounts []struct{ Source string }
	}
	Config struct {
		Env    []string
		Labels map[string]string
	}
}

func inspect(t *testing.T, target string) (inspectDoc, string) {
	t.Helper()
	raw := mustDocker(t, "inspect", target)
	var docs []inspectDoc
	if err := json.Unmarshal([]byte(raw), &docs); err != nil || len(docs) != 1 {
		t.Fatalf("docker inspect %s: %v (%d documents)", target, err, len(docs))
	}
	return docs[0], raw
}

// hasSocketMount reports whether the container holds a docker.sock mount.
func hasSocketMount(doc inspectDoc) bool {
	for _, m := range doc.Mounts {
		if strings.Contains(m.Source, "docker.sock") || strings.Contains(m.Destination, "docker.sock") {
			return true
		}
	}
	for _, b := range doc.HostConfig.Binds {
		if strings.Contains(b, "docker.sock") {
			return true
		}
	}
	for _, m := range doc.HostConfig.Mounts {
		if strings.Contains(m.Source, "docker.sock") {
			return true
		}
	}
	return false
}

// ADR-0028 Entscheidung 5, S2: the label of the started container is the
// path of the host and FindContainer finds the container.
func TestLauncher_Docker_S2_LabelCarriesHostPath(t *testing.T) {
	folder := newProject(t, imageConfig(nil), nil)
	res := mustUp(t, driven.DevcontainerUpOptions{ProjectFolder: folder})

	doc, _ := inspect(t, res.ContainerID)
	if got := doc.Config.Labels[folderLabel]; got != folder {
		t.Errorf("label %s = %q, want the host path %q", folderLabel, got, folder)
	}
	found, ok, err := docker.NewInspector().FindContainer(context.Background(), folder)
	if err != nil || !ok || found.ID != res.ContainerID || !found.Running {
		t.Errorf("FindContainer(%s) = %+v ok=%v err=%v, want the started container %s", folder, found, ok, err, res.ContainerID)
	}
}

// Guard of S2: the measurement path (the project below `/w/…` in the
// helper container) yields a label that differs from the path of the
// host, and FindContainer does not find the container for the host path.
func TestLauncher_Docker_S2_GuardMeasurementPathDeviates(t *testing.T) {
	requireDocker(t)
	folder := filepath.Join(t.TempDir(), "ubl"+randomHex(t))
	if err := os.MkdirAll(filepath.Join(folder, ".devcontainer"), 0o755); err != nil {
		t.Fatal(err)
	}
	volume := "ubl-s2-" + randomHex(t)
	inHelper := "/w/" + filepath.Base(folder)
	// The workspace is a volume: the daemon needs no `/w/…` on its host.
	writeConfig(t, folder, imageConfig(map[string]any{
		"workspaceMount":  "source=" + volume + ",target=/workspaces/probe,type=volume",
		"workspaceFolder": "/workspaces/probe",
	}))
	t.Cleanup(func() {
		removeContainers("label=" + folderLabel + "=" + inHelper)
		removeImages("vsc-" + filepath.Base(folder) + "-*")
		_, _ = tryDocker("volume", "rm", "-f", volume)
	})
	tag, err := docker.NewLauncher().EnsureImage(context.Background(), nil)
	if err != nil {
		t.Fatalf("EnsureImage: %v", err)
	}
	mustDocker(t, "run", "--rm",
		"--mount", "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock",
		"--mount", "type=bind,source="+folder+",target="+inHelper,
		tag, "devcontainer", "up", "--workspace-folder", inHelper)

	id := mustDocker(t, "ps", "-aq", "--filter", "label="+folderLabel+"="+inHelper)
	if id == "" {
		t.Fatal("the measurement-path container was not created")
	}
	doc, _ := inspect(t, id)
	if got := doc.Config.Labels[folderLabel]; got == folder {
		t.Errorf("label = %q equals the host path: the guard cannot tell the measurement path from the mirror path", got)
	}
	if _, ok, err := docker.NewInspector().FindContainer(context.Background(), folder); ok || err != nil {
		t.Errorf("FindContainer(%s) = ok %v err %v, want not found for the measurement path", folder, ok, err)
	}
}

// ADR-0028 Entscheidung 4, S6: the started container holds no engine socket.
func TestLauncher_Docker_S6_SocketNotInStartedContainer(t *testing.T) {
	folder := newProject(t, imageConfig(nil), nil)
	res := mustUp(t, driven.DevcontainerUpOptions{ProjectFolder: folder})
	doc, raw := inspect(t, res.ContainerID)
	if hasSocketMount(doc) {
		t.Errorf("the started container holds the engine socket:\n%s", raw)
	}
}

// Guard of S6: the same check recognises a hand-written devcontainer.json
// that mounts the socket itself.
func TestLauncher_Docker_S6_GuardHandWrittenSocketMountIsDetected(t *testing.T) {
	folder := newProject(t, imageConfig(map[string]any{
		"mounts": []string{"source=/var/run/docker.sock,target=/var/run/docker.sock,type=bind"},
	}), nil)
	res := mustUp(t, driven.DevcontainerUpOptions{ProjectFolder: folder})
	doc, raw := inspect(t, res.ContainerID)
	if !hasSocketMount(doc) {
		t.Errorf("the socket check missed a mount declared in devcontainer.json:\n%s", raw)
	}
}

// ADR-0029 Entscheidung 4, S6: an initializeCommand of a hand-written
// devcontainer.json runs in the helper container with the engine socket
// (it writes a marker file and gets the answer of the daemon).
func TestLauncher_Docker_S6_InitializeCommandRunsWithSocket(t *testing.T) {
	folder := newProject(t, imageConfig(map[string]any{
		"initializeCommand": "touch .ubl-marker && docker version --format '{{.Server.Version}}' > .ubl-daemon-version",
	}), nil)
	mustUp(t, driven.DevcontainerUpOptions{ProjectFolder: folder})

	if _, err := os.Stat(filepath.Join(folder, ".ubl-marker")); err != nil {
		t.Errorf("marker file of the initializeCommand: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(folder, ".ubl-daemon-version"))
	if err != nil {
		t.Fatalf("daemon answer of the initializeCommand: %v", err)
	}
	want := mustDocker(t, "version", "--format", "{{.Server.Version}}")
	if strings.TrimSpace(string(got)) != want {
		t.Errorf("daemon version seen by the initializeCommand = %q, want %q", got, want)
	}
}

// ADR-0028 Entscheidung 4, S6 readonly: with a read-only project mount the
// CLI itself starts the container for an image-based and a Dockerfile-based
// project and writes nothing into the folder; an initializeCommand that
// writes into the folder fails, and succeeds with a writable mount.
func TestLauncher_Docker_S6_ReadOnlyProject(t *testing.T) {
	t.Run("CLI writes nothing for an image-based project", func(t *testing.T) {
		folder := newProject(t, imageConfig(nil), nil)
		before := listFiles(t, folder)
		mustUp(t, driven.DevcontainerUpOptions{ProjectFolder: folder, ReadOnlyProject: true})
		if after := listFiles(t, folder); after != before {
			t.Errorf("project files changed:\nbefore: %s\nafter:  %s", before, after)
		}
	})
	t.Run("CLI writes nothing for a Dockerfile-based project", func(t *testing.T) {
		folder := newProject(t,
			map[string]any{"name": "ubl-launcher-test", "build": map[string]any{"dockerfile": "Dockerfile"}},
			nil)
		if err := os.WriteFile(filepath.Join(folder, ".devcontainer", "Dockerfile"), []byte("FROM "+devBaseImage+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		before := listFiles(t, folder)
		mustUp(t, driven.DevcontainerUpOptions{ProjectFolder: folder, ReadOnlyProject: true})
		if after := listFiles(t, folder); after != before {
			t.Errorf("project files changed:\nbefore: %s\nafter:  %s", before, after)
		}
	})
	t.Run("initializeCommand that writes", func(t *testing.T) {
		cfg := imageConfig(map[string]any{"initializeCommand": "touch .ubl-written"})
		folder := newProject(t, cfg, nil)
		_, progress, err := up(t, driven.DevcontainerUpOptions{ProjectFolder: folder, ReadOnlyProject: true})
		if !errors.Is(err, driven.ErrDevcontainerRuntime) {
			t.Errorf("read-only: error = %v, want ErrDevcontainerRuntime\n%s", err, progress)
		}
		if _, statErr := os.Stat(filepath.Join(folder, ".ubl-written")); statErr == nil {
			t.Errorf("a read-only project mount let the initializeCommand write")
		}
		mustUp(t, driven.DevcontainerUpOptions{ProjectFolder: folder})
		if _, statErr := os.Stat(filepath.Join(folder, ".ubl-written")); statErr != nil {
			t.Errorf("writable: initializeCommand wrote no file: %v", statErr)
		}
	})
}

func listFiles(t *testing.T, root string) string {
	t.Helper()
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, strings.TrimPrefix(path, root)+":"+info.ModTime().String())
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return strings.Join(files, ",")
}

// ADR-0028 Entscheidung 4, S6 Lebensdauer: no helper container remains
// after Up, after a cancelled Up, and the guard (a run without `--rm`)
// leaves one that the check finds.
func TestLauncher_Docker_S6_Lifetime(t *testing.T) {
	t.Run("after Up", func(t *testing.T) {
		folder := newProject(t, imageConfig(nil), nil)
		mustUp(t, driven.DevcontainerUpOptions{ProjectFolder: folder})
		if left := helperLeftovers(t); len(left) != 0 {
			t.Errorf("helper containers after Up: %v", left)
		}
	})
	t.Run("after a cancelled Up", func(t *testing.T) {
		folder := newProject(t, imageConfig(map[string]any{"initializeCommand": "sleep 120"}), nil)
		ctx, cancel := context.WithTimeout(context.Background(), launcherTestBudget)
		defer cancel()
		launcher := testLauncher(t)
		done := make(chan error, 1)
		go func() {
			_, err := launcher.Up(ctx, driven.DevcontainerUpOptions{ProjectFolder: folder})
			done <- err
		}()
		deadline := time.Now().Add(5 * time.Minute)
		for len(helperLeftovers(t)) == 0 {
			select {
			case err := <-done:
				t.Fatalf("Up ended before the helper container was seen: %v", err)
			default:
			}
			if time.Now().After(deadline) {
				t.Fatal("no helper container appeared")
			}
			time.Sleep(200 * time.Millisecond)
		}
		cancel()
		select {
		case err := <-done:
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Up error = %v, want context.Canceled", err)
			}
		case <-time.After(2 * time.Minute):
			t.Fatal("Up did not return after cancel")
		}
		if left := helperLeftovers(t); len(left) != 0 {
			t.Errorf("helper containers after a cancelled Up: %v", left)
		}
	})
	t.Run("guard: a run without --rm leaves a helper container", func(t *testing.T) {
		requireDocker(t)
		t.Cleanup(func() { removeContainers("name=" + helperPrefix(t)) })
		tag, err := docker.NewLauncher().EnsureImage(context.Background(), nil)
		if err != nil {
			t.Fatalf("EnsureImage: %v", err)
		}
		args := docker.RunArgsForTest(helperPrefix(t)+"guard-"+randomHex(t), tag, nil,
			driven.DevcontainerUpOptions{ProjectFolder: "/tmp"})
		var noRm []string
		for _, a := range args {
			if a == "--rm" {
				continue
			}
			noRm = append(noRm, a)
		}
		// Replace `devcontainer up …` by a command that ends at once.
		for i, a := range noRm {
			if a == tag {
				noRm = append(noRm[:i+1], "devcontainer", "--version")
				break
			}
		}
		mustDocker(t, noRm...)
		if left := helperLeftovers(t); len(left) != 1 {
			t.Errorf("the leftover check found %v, want the one container of the run without --rm", left)
		}
	})
}

// cmdlineHolders returns the command lines of all processes below /proc
// that contain value.
func cmdlineHolders(value string) []string {
	entries, err := filepath.Glob("/proc/[0-9]*/cmdline")
	if err != nil {
		return nil
	}
	var holders []string
	for _, e := range entries {
		raw, err := os.ReadFile(e)
		if err != nil {
			continue
		}
		if cmdline := strings.ReplaceAll(string(raw), "\x00", " "); strings.Contains(cmdline, value) {
			holders = append(holders, cmdline)
		}
	}
	return holders
}

// helperCmdlineHolders returns the lines of `docker top` of the helper
// containers below prefix that contain value: the command lines of the
// processes inside the helper container (the devcontainer CLI and the
// `docker` commands it runs), which the test container's /proc does not show.
func helperCmdlineHolders(prefix, value string) []string {
	ids, err := tryDocker("ps", "-q", "--filter", "name="+prefix)
	if err != nil {
		return nil
	}
	var holders []string
	for _, id := range strings.Fields(ids) {
		out, topErr := tryDocker("top", id)
		if topErr != nil {
			continue
		}
		for _, line := range strings.Split(out, "\n") {
			if strings.Contains(line, value) {
				holders = append(holders, line)
			}
		}
	}
	return holders
}

// scanDuring polls probe until ctx is done and returns the distinct hits.
// The goroutine ends with ctx; the caller cancels it with defer, so a failed
// test does not leave it running.
func scanDuring(ctx context.Context, probe func() []string) <-chan []string {
	out := make(chan []string, 1)
	go func() {
		seen := map[string]bool{}
		var hits []string
		for {
			for _, h := range probe() {
				if !seen[h] {
					seen[h] = true
					hits = append(hits, h)
				}
			}
			select {
			case <-ctx.Done():
				out <- hits
				return
			case <-time.After(15 * time.Millisecond):
			}
		}
	}()
	return out
}

func readInContainer(t *testing.T, id, path string) string {
	t.Helper()
	return mustDocker(t, "exec", id, "sh", "-c", "cat "+path+" 2>/dev/null || true")
}

// LH-FA-DEV-009, ADR-0028 Entscheidung 4, S6 Umgebung: a variable the
// launcher is told to pass reaches the devcontainer through `remoteEnv`; one
// it is not told to pass does not, though devcontainer.json references it.
func TestLauncher_Docker_S6_EnvReachesOnlyPassedNames(t *testing.T) {
	suffix := randomHex(t)
	refName, otherName := "UBL_TEST_REF_"+suffix, "UBL_TEST_OTHER_"+suffix
	refValue, otherValue := "placeholder-ref-"+suffix, "placeholder-other-"+suffix
	t.Setenv(refName, refValue)
	t.Setenv(otherName, otherValue)
	cfg := imageConfig(map[string]any{
		"remoteEnv":         map[string]string{"UBL_REF": "${localEnv:" + refName + "}", "UBL_OTHER": "${localEnv:" + otherName + "}"},
		"postCreateCommand": "printenv UBL_REF > /tmp/ubl-ref; printenv UBL_OTHER > /tmp/ubl-other; true",
	})
	folder := newProject(t, cfg, nil)

	res := mustUp(t, driven.DevcontainerUpOptions{ProjectFolder: folder, EnvNames: []string{refName}})
	if got := readInContainer(t, res.ContainerID, "/tmp/ubl-ref"); got != refValue {
		t.Errorf("passed variable in the devcontainer = %q, want %q", got, refValue)
	}
	if got := readInContainer(t, res.ContainerID, "/tmp/ubl-other"); got != "" {
		t.Errorf("a variable that was not passed reached the devcontainer: %q", got)
	}

	// Guard: the probe sees a variable that does arrive.
	res = mustUp(t, driven.DevcontainerUpOptions{ProjectFolder: folder, EnvNames: []string{refName, otherName}, Recreate: true})
	if got := readInContainer(t, res.ContainerID, "/tmp/ubl-other"); got != otherValue {
		t.Errorf("guard: a passed variable did not arrive (%q); the probe cannot detect a leak", got)
	}
}

// LH-FA-DEV-009, S6: with the `remoteEnv` pattern the value of a passed
// variable is in no image, no label of the started container, no file under
// the project folder and on no command line of a process in the test
// container (`/proc`, sampled every 15 ms plus the duration of the scan; the
// launcher's own `docker run` is such a process). The command lines inside
// the helper container are not covered: see the limit test below.
func TestLauncher_Docker_S6_ValueIsNeverInImageLabelFileOrCommandLine(t *testing.T) {
	suffix := randomHex(t)
	name := "UBL_TEST_SECRET_" + suffix
	value := "placeholder-secret-" + suffix
	t.Setenv(name, value)
	folder := newProject(t, imageConfig(map[string]any{
		"remoteEnv":         map[string]string{"UBL_SECRET": "${localEnv:" + name + "}"},
		"initializeCommand": "sleep 3",
		"postCreateCommand": "printenv UBL_SECRET > /tmp/ubl-secret; true",
	}), nil)

	scanCtx, stopScan := context.WithCancel(context.Background())
	defer stopScan()
	hits := scanDuring(scanCtx, func() []string { return cmdlineHolders(value) })
	res := mustUp(t, driven.DevcontainerUpOptions{ProjectFolder: folder, EnvNames: []string{name}})
	stopScan()
	if got := <-hits; len(got) != 0 {
		t.Errorf("the value was on a command line: %q", got)
	}
	if got := readInContainer(t, res.ContainerID, "/tmp/ubl-secret"); got != value {
		t.Fatalf("the variable did not arrive (%q); the test proves nothing", got)
	}

	tag, err := docker.NewLauncher().EnsureImage(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	for what, target := range map[string]string{"helper image": tag, "started devcontainer": res.ContainerID} {
		if _, raw := inspect(t, target); strings.Contains(raw, value) {
			t.Errorf("the value is in docker inspect of the %s", what)
		}
	}
	err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if raw, readErr := os.ReadFile(path); readErr == nil && bytes.Contains(raw, []byte(value)) {
			t.Errorf("the value is in %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// Limit of LH-FA-DEV-009 (documented in the port and in the slice): the
// Dev Container CLI hands a `remoteEnv` value to a lifecycle command as
// `docker exec -e NAME=<value>`, so the value is on a command line inside
// the helper container for the duration of that command (`docker top`, also
// visible on the host). The test pins this observation; if it turns red, the
// CLI does not do it and the limit may be dropped from the text.
func TestLauncher_Docker_S6_LimitRemoteEnvValueIsOnHelperCommandLine(t *testing.T) {
	suffix := randomHex(t)
	name := "UBL_TEST_LIMIT_" + suffix
	value := "placeholder-limit-" + suffix
	t.Setenv(name, value)
	folder := newProject(t, imageConfig(map[string]any{
		"remoteEnv":         map[string]string{"UBL_SECRET": "${localEnv:" + name + "}"},
		"postCreateCommand": "sleep 4",
	}), nil)
	scanCtx, stopScan := context.WithCancel(context.Background())
	defer stopScan()
	hits := scanDuring(scanCtx, func() []string { return helperCmdlineHolders(helperPrefix(t), value) })
	mustUp(t, driven.DevcontainerUpOptions{ProjectFolder: folder, EnvNames: []string{name}})
	stopScan()
	if len(<-hits) == 0 {
		t.Errorf("no command line in the helper container carried the value: the CLI does not write remoteEnv into `docker exec -e`; drop the limit")
	}
}

// LH-FA-DEV-009, S6: a devcontainer.json that writes `${localEnv:NAME}` into
// `containerEnv` makes the CLI print `docker run … -e X=<value>` on its
// stderr. The launcher replaces the value by a placeholder in the progress it
// forwards and in the error text, on success and on failure. Guard: the same
// value is on a command line inside the helper container (the CLI writes it
// there; the launcher cannot mask a process list), so the helper scan sees
// what it is meant to see, and the stream without masking would carry it.
func TestLauncher_Docker_S6_ContainerEnvValueIsMaskedInProgressAndError(t *testing.T) {
	suffix := randomHex(t)
	name := "UBL_TEST_MASK_" + suffix
	value := "placeholder-mask-" + suffix
	t.Setenv(name, value)
	config := func(extra map[string]any) map[string]any {
		cfg := imageConfig(map[string]any{
			"containerEnv":      map[string]string{"X": "${localEnv:" + name + "}"},
			"remoteEnv":         map[string]string{"UBL_SECRET": "${localEnv:" + name + "}"},
			"initializeCommand": "sleep 3",
		})
		for k, v := range extra {
			cfg[k] = v
		}
		return cfg
	}
	opts := driven.DevcontainerUpOptions{EnvNames: []string{name}}

	t.Run("Up succeeds", func(t *testing.T) {
		opts.ProjectFolder = newProject(t, config(nil), nil)
		scanCtx, stopScan := context.WithCancel(context.Background())
		defer stopScan()
		hits := scanDuring(scanCtx, func() []string { return helperCmdlineHolders(helperPrefix(t), value) })
		_, progress, err := up(t, opts)
		stopScan()
		if err != nil {
			t.Fatalf("Up: %v\nprogress:\n%s", err, progress)
		}
		if len(<-hits) == 0 {
			t.Errorf("guard: the helper scan did not see the value the CLI writes into its own `docker run -e X=…`; the probe cannot tell masking from absence")
		}
		if strings.Contains(progress, value) {
			t.Errorf("the value is in the progress:\n%s", progress)
		}
		if !strings.Contains(progress, "-e X=[masked]") {
			t.Errorf("the progress shows no masked `-e X=` of the CLI (the stream the test must cover):\n%s", progress)
		}
	})
	t.Run("Up fails", func(t *testing.T) {
		opts.ProjectFolder = newProject(t, config(map[string]any{"postCreateCommand": "exit 1"}), nil)
		_, progress, err := up(t, opts)
		if !errors.Is(err, driven.ErrDevcontainerRuntime) {
			t.Fatalf("Up error = %v, want ErrDevcontainerRuntime\nprogress:\n%s", err, progress)
		}
		if strings.Contains(err.Error(), value) || strings.Contains(progress, value) {
			t.Errorf("the value is in the error text or the progress:\n%v\n%s", err, progress)
		}
		if !strings.Contains(progress, "[masked]") {
			t.Errorf("the progress shows no placeholder:\n%s", progress)
		}
	})
}

// Guard of the command-line check: it finds a value that is part of an
// argument (`-e NAME=value`).
func TestLauncher_Docker_S6_GuardCommandLineScanFindsValueInArguments(t *testing.T) {
	requireDocker(t)
	tag, err := docker.NewLauncher().EnsureImage(context.Background(), nil)
	if err != nil {
		t.Fatalf("EnsureImage: %v", err)
	}
	suffix := randomHex(t)
	value := "placeholder-guard-" + suffix
	name := docker.HelperContainerPrefixForTest + "scan-" + suffix
	t.Cleanup(func() { removeContainers("name=" + name) })
	cmd := exec.Command("docker", "run", "--rm", "--name", name, "-e", "GUARD="+value, "--entrypoint", "sleep", tag, "30")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }()
	deadline := time.Now().Add(20 * time.Second)
	for len(cmdlineHolders(value)) == 0 {
		if time.Now().After(deadline) {
			t.Fatal("the command-line scan did not find a value passed as `-e NAME=value`")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// ADR-0028 Entscheidung 2, S5: the image builds from the embedded context
// after its tag was removed, and the CLI inside answers with the pin.
func TestLauncher_Docker_S5_BuildFromEmbeddedContext(t *testing.T) {
	requireDocker(t)
	tag, err := docker.HelperImageTagForTest()
	if err != nil {
		t.Fatal(err)
	}
	pin, err := docker.HelperCLIVersionForTest()
	if err != nil {
		t.Fatal(err)
	}
	_, _ = tryDocker("rmi", tag)
	if _, err := tryDocker("image", "inspect", tag); err == nil {
		t.Fatalf("image %s still exists after rmi", tag)
	}
	var progress bytes.Buffer
	got, err := docker.NewLauncher().EnsureImage(context.Background(), &progress)
	if err != nil || got != tag {
		t.Fatalf("EnsureImage = %q, %v; want %q\n%s", got, err, tag, progress.String())
	}
	if version := mustDocker(t, "run", "--rm", tag, "devcontainer", "--version"); version != pin {
		t.Errorf("devcontainer --version = %q, want the pin %q", version, pin)
	}
}

// Guard of S5: a falsified `integrity` in the lock file breaks the build;
// the unchanged lock file builds.
func TestLauncher_Docker_S5_GuardFalsifiedIntegrityBreaksTheBuild(t *testing.T) {
	requireDocker(t)
	files, err := docker.HelperContextFilesForTest()
	if err != nil {
		t.Fatal(err)
	}
	build := func(t *testing.T, lock []byte) (string, error) {
		t.Helper()
		dir := t.TempDir()
		for name, content := range files {
			if name == "package-lock.json" {
				content = lock
			}
			if err := os.WriteFile(filepath.Join(dir, name), content, 0o644); err != nil {
				t.Fatal(err)
			}
		}
		image := "ubl-s5-guard-" + randomHex(t)
		t.Cleanup(func() { removeImages(image) })
		return tryDocker("build", "--force-rm", "-q", "-t", image, dir)
	}

	falsified := bytes.Replace(files["package-lock.json"], []byte(`"integrity": "sha512-`), []byte(`"integrity": "sha512-AAAA`), 1)
	if bytes.Equal(falsified, files["package-lock.json"]) {
		t.Fatal("the lock file holds no integrity hash to falsify")
	}
	if out, err := build(t, falsified); err == nil {
		t.Errorf("the build accepted a falsified integrity hash: %s", out)
	} else if !strings.Contains(strings.ToLower(err.Error()), "integrity") {
		t.Errorf("the build failed for another reason than the integrity check: %v", err)
	}
	if out, err := build(t, files["package-lock.json"]); err != nil {
		t.Errorf("the unchanged lock file does not build: %v\n%s", err, out)
	}
}

// S7: the adapter runs against the engine of the runner. The environment
// check passes there (a missing container fails at `docker stop`, a runtime
// failure), and an unreachable engine is ErrDockerUnavailable.
func TestLauncher_Docker_S7_EngineOfTheRunner(t *testing.T) {
	requireDocker(t)
	missing := strings.Repeat("0", 64)
	err := docker.NewLauncher().Stop(context.Background(), missing)
	if !errors.Is(err, driven.ErrDevcontainerRuntime) || errors.Is(err, driven.ErrDockerUnavailable) {
		t.Errorf("Stop of a missing container on the runner engine: error = %v, want ErrDevcontainerRuntime", err)
	}

	t.Setenv("DOCKER_HOST", "tcp://127.0.0.1:1")
	if err := docker.NewLauncher().Stop(context.Background(), missing); !errors.Is(err, driven.ErrDockerUnavailable) {
		t.Errorf("Stop with an unreachable engine: error = %v, want ErrDockerUnavailable", err)
	}
}

// ADR-0028 Entscheidung 3: Stop stops the container and does not remove it.
func TestLauncher_Docker_StopStopsButKeepsTheContainer(t *testing.T) {
	folder := newProject(t, imageConfig(nil), nil)
	res := mustUp(t, driven.DevcontainerUpOptions{ProjectFolder: folder})
	if err := docker.NewLauncher().Stop(context.Background(), res.ContainerID); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	found, ok, err := docker.NewInspector().FindContainer(context.Background(), folder)
	if err != nil || !ok || found.ID != res.ContainerID || found.Running {
		t.Errorf("after Stop: %+v ok=%v err=%v, want the stopped container", found, ok, err)
	}
}
