//go:build docker

// End-to-end check of the status files of the devcontainer start scripts
// (LH-FA-DEV-007, LH-FA-DEV-008): generate a sandbox project with the
// production adapters, build the image, start a container that carries the
// label the devcontainer CLI sets (`devcontainer.local_folder`), run the
// generated start scripts the way the CLI does (`docker exec` of the
// lifecycle command as the container user), and read the result from the host
// with the production Inspector adapter and `doctor`. The Dev Containers CLI
// itself is not part of the repository tooling; the label it sets and the
// lifecycle commands it runs are the contract these tests reproduce.
// Hosts with kernel.apparmor_restrict_unprivileged_userns=1 skip the nested
// Podman cases.

package e2e_test

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	dockeradapter "github.com/pt9912/u-boot/internal/adapter/driven/docker"
	fsadapter "github.com/pt9912/u-boot/internal/adapter/driven/fs"
	gitadapter "github.com/pt9912/u-boot/internal/adapter/driven/git"
	loggeradapter "github.com/pt9912/u-boot/internal/adapter/driven/logger"
	yamladapter "github.com/pt9912/u-boot/internal/adapter/driven/yaml"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const statusDirInContainer = "/var/lib/u-boot/status"

func requireDocker(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
}

// buildProjectImage builds the generated Dockerfile of dir and removes the
// image at the end of the test.
func buildProjectImage(ctx context.Context, t *testing.T, dir, tag string) string {
	t.Helper()
	image := tag + ":" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("docker build of the generated Dockerfile failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })
	return image
}

// startLabeled starts a container the way the devcontainer CLI does (the
// label `devcontainer.local_folder` carries the project folder) with the given
// extra `docker run` arguments, and removes it at the end of the test. It
// returns the container ID, read from stdout only. The container carries a
// generated name, and the cleanup is registered before `docker run`, so a
// container whose ID cannot be read is removed as well.
func startLabeled(ctx context.Context, t *testing.T, folder, image string, extra ...string) string {
	t.Helper()
	name := "sbxe2e-lbl-" + strconv.FormatInt(time.Now().UnixNano(), 36)
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rm", "-f", "-v", name) })
	args := append([]string{"run", "-d", "-q", "--name", name, "--label", "devcontainer.local_folder=" + folder}, extra...)
	args = append(args, image, "sleep", "1800")
	var stdout, stderr strings.Builder
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("docker run failed: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	id, err := containerIDFromRunOutput(stdout.String())
	if err != nil {
		t.Fatalf("%v\nstderr:\n%s", err, stderr.String())
	}
	return id
}

func removeContainer(ctx context.Context, t *testing.T, name string) {
	t.Helper()
	if out, code := dockerOutput(ctx, t, "rm", "-f", "-v", name); code != 0 {
		t.Fatalf("docker rm: %s", out)
	}
}

// lifecycle runs a lifecycle command of devcontainer.json in the container as
// the container user, like the devcontainer CLI (`docker exec`).
func lifecycle(ctx context.Context, t *testing.T, container, command string) (string, int) {
	t.Helper()
	return dockerOutput(ctx, t, "exec", "-u", "vscode", container, "sh", "-c", command)
}

// hostStatus reads one status file the way `doctor` does from the host: find
// the container by the project folder, copy the file out, parse it strictly.
func hostStatus(ctx context.Context, t *testing.T, folder, file string) (domain.SandboxStatus, driven.DevcontainerContainer) {
	t.Helper()
	insp := dockeradapter.NewInspector()
	c, found, err := insp.FindContainer(ctx, folder)
	if err != nil || !found {
		t.Fatalf("FindContainer(%s): found=%v err=%v", folder, found, err)
	}
	raw, err := insp.ReadFile(ctx, c.ID, statusDirInContainer+"/"+file, domain.SandboxStatusMaxBytes)
	if err != nil {
		t.Fatalf("ReadFile %s: %v", file, err)
	}
	status, err := domain.ParseSandboxStatus(raw)
	if err != nil {
		t.Fatalf("status file %s is not a valid document: %v\n%s", file, err, raw)
	}
	if age := time.Since(status.At()); age < -time.Minute || age > 10*time.Minute {
		t.Errorf("status file %s carries the time %v, not a recent script run", file, status.At())
	}
	return status, c
}

// doctorItem runs the production `doctor` against dir with the real Docker
// adapters and returns the item with the given ID.
func doctorItem(ctx context.Context, t *testing.T, dir, id string) domain.Diagnostic {
	t.Helper()
	svc := application.NewDoctorService(fsadapter.New(), yamladapter.New(), gitadapter.New(), dockeradapter.New(), nil,
		loggeradapter.New(os.Stderr, loggeradapter.FormatText, nil)).
		WithDevcontainerInspector(dockeradapter.NewInspector())
	resp, err := svc.Check(ctx, driving.DoctorRequest{BaseDir: dir})
	if err != nil {
		t.Fatalf("doctor: %v", err)
	}
	for _, d := range resp.Report.Items {
		if d.ID == id {
			return d
		}
	}
	t.Fatalf("doctor emitted no item %s", id)
	return domain.Diagnostic{}
}

func withoutArg(args []string, drop string) []string {
	var out []string
	for _, a := range args {
		if a != drop {
			out = append(out, a)
		}
	}
	return out
}

// expectItem asserts severity and fallback flag of a doctor item.
func expectItem(t *testing.T, d domain.Diagnostic, severity domain.Severity, fallback bool) {
	t.Helper()
	if d.Severity != severity || d.DegradationFallback != fallback {
		t.Errorf("%s: %v fallback=%v (%q), want %v fallback=%v", d.ID, d.Severity, d.DegradationFallback, d.Message, severity, fallback)
	}
}

// LH-FA-DEV-007: the `vfs` start script writes `ok`, `fallback` and `failed`
// into /var/lib/u-boot/status/vfs; the host reads it and `doctor` maps it.
func TestE2E_SandboxDevcontainer_StatusFiles_VFS(t *testing.T) {
	t.Parallel()
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()
	const startCmd = "sh /usr/local/bin/u-boot-sandbox-init"
	const itemID = "devcontainer.sandbox.fallback.vfs"

	warnDir := generateSandboxProject(ctx, t, "podman", "warn")
	warnImage := buildProjectImage(ctx, t, warnDir, "uboot-e2e-status-vfs")
	runArgs := generatedRunArgs(t, warnDir)
	skipIfNestedUsernsRestricted(ctx, t, warnImage, runArgs)

	// The Dockerfile created the directory for the container user.
	probe := startLabeled(ctx, t, warnDir, warnImage, runArgs...)
	if out, code := lifecycle(ctx, t, probe, "test -d "+statusDirInContainer+" && test -w "+statusDirInContainer+" && ! test -e "+statusDirInContainer+"/vfs"); code != 0 {
		t.Fatalf("status directory missing or not writable for the container user, or a file exists before the script ran (exit %d):\n%s", code, out)
	}

	// Control: with the device the script reports ok.
	if out, code := lifecycle(ctx, t, probe, startCmd); code != 0 {
		t.Fatalf("start script exit %d:\n%s", code, out)
	}
	status, c := hostStatus(ctx, t, warnDir, "vfs")
	if status.Result() != domain.SandboxStatusOK || !c.Running {
		t.Errorf("with /dev/fuse: result=%q running=%v, want ok in a running container", status.Result(), c.Running)
	}
	expectItem(t, doctorItem(ctx, t, warnDir, itemID), domain.SeverityOK, false)
	removeContainer(ctx, t, probe)

	// Without the device and policy warn: fallback to vfs, exit 0.
	noFuse := startLabeled(ctx, t, warnDir, warnImage, withoutArg(runArgs, "--device=/dev/fuse")...)
	if out, code := lifecycle(ctx, t, noFuse, startCmd); code != 0 || !strings.Contains(out, "falling back to vfs storage") {
		t.Fatalf("warn policy without /dev/fuse: exit %d (stderr line kept?):\n%s", code, out)
	}
	if status, _ := hostStatus(ctx, t, warnDir, "vfs"); status.Result() != domain.SandboxStatusFallback {
		t.Errorf("result = %q, want fallback", status.Result())
	}
	expectItem(t, doctorItem(ctx, t, warnDir, itemID), domain.SeverityWarn, true)

	// A stopped container is read; the item says so.
	if out, code := dockerOutput(ctx, t, "stop", "-t", "1", noFuse); code != 0 {
		t.Fatalf("docker stop: %s", out)
	}
	status, c = hostStatus(ctx, t, warnDir, "vfs")
	if status.Result() != domain.SandboxStatusFallback || c.Running {
		t.Errorf("stopped container: result=%q running=%v", status.Result(), c.Running)
	}
	if d := doctorItem(ctx, t, warnDir, itemID); !strings.Contains(d.Message, "container stopped") {
		t.Errorf("stopped container not named: %q", d.Message)
	}
	removeContainer(ctx, t, noFuse)

	// Policy fail: exit 11, status failed, doctor error without the fallback flag.
	failDir := generateSandboxProject(ctx, t, "podman", "fail")
	failImage := buildProjectImage(ctx, t, failDir, "uboot-e2e-status-vfs-fail")
	failBox := startLabeled(ctx, t, failDir, failImage, withoutArg(generatedRunArgs(t, failDir), "--device=/dev/fuse")...)
	if out, code := lifecycle(ctx, t, failBox, startCmd); code != 11 {
		t.Fatalf("fail policy without /dev/fuse: exit %d, want 11:\n%s", code, out)
	}
	if status, _ := hostStatus(ctx, t, failDir, "vfs"); status.Result() != domain.SandboxStatusFailed {
		t.Errorf("result = %q, want failed", status.Result())
	}
	expectItem(t, doctorItem(ctx, t, failDir, itemID), domain.SeverityError, false)
}

// LH-FA-DEV-007: the daemon start script reports `fallback` for a container
// started without `--privileged`, `failed` under policy `fail`, and `ok` when
// the daemon becomes ready.
func TestE2E_SandboxDevcontainer_StatusFiles_Daemon(t *testing.T) {
	t.Parallel()
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
	defer cancel()
	const startCmd = "sudo sh /usr/local/bin/u-boot-docker-init"
	const itemID = "devcontainer.sandbox.fallback.daemon"

	warnDir := generateSandboxProject(ctx, t, "docker", "warn")
	warnImage := buildProjectImage(ctx, t, warnDir, "uboot-e2e-status-daemon")

	// No --privileged: no daemon, policy warn → exit 0 and fallback.
	plain := startLabeled(ctx, t, warnDir, warnImage)
	if out, code := lifecycle(ctx, t, plain, startCmd); code != 0 || !strings.Contains(out, "no daemon available") {
		t.Fatalf("warn policy without --privileged: exit %d:\n%s", code, out)
	}
	if status, _ := hostStatus(ctx, t, warnDir, "daemon"); status.Result() != domain.SandboxStatusFallback {
		t.Errorf("result = %q, want fallback", status.Result())
	}
	expectItem(t, doctorItem(ctx, t, warnDir, itemID), domain.SeverityWarn, true)
	removeContainer(ctx, t, plain)

	// With --privileged and the daemon volume the daemon becomes ready.
	extra := append([]string(nil), generatedRunArgs(t, warnDir)...)
	instance := "stat" + strconv.FormatInt(time.Now().UnixNano(), 10)
	var volumes []string
	for _, m := range generatedMounts(t, warnDir) {
		mount := strings.ReplaceAll(m, "${devcontainerId}", instance)
		extra = append(extra, "--mount", mount)
		volumes = append(volumes, strings.TrimPrefix(strings.Split(mount, ",")[0], "source="))
	}
	// Registered before the container: cleanups run last-in first-out, and a
	// volume can only be removed after the container that uses it.
	t.Cleanup(func() {
		for _, v := range volumes {
			_, _ = dockerOutput(context.Background(), t, "volume", "rm", "-f", v)
		}
	})
	privileged := startLabeled(ctx, t, warnDir, warnImage, extra...)
	if out, code := lifecycle(ctx, t, privileged, startCmd); code != 0 {
		t.Fatalf("daemon start script exit %d:\n%s", code, out)
	}
	if status, _ := hostStatus(ctx, t, warnDir, "daemon"); status.Result() != domain.SandboxStatusOK {
		t.Errorf("with a ready daemon: result = %q, want ok", status.Result())
	}
	expectItem(t, doctorItem(ctx, t, warnDir, itemID), domain.SeverityOK, false)

	// A second run of the script on the running daemon stays ok.
	if out, code := lifecycle(ctx, t, privileged, startCmd); code != 0 {
		t.Fatalf("second start script run exit %d:\n%s", code, out)
	}
	if status, _ := hostStatus(ctx, t, warnDir, "daemon"); status.Result() != domain.SandboxStatusOK {
		t.Errorf("second run: result = %q, want ok", status.Result())
	}

	// Policy fail without --privileged: exit 11, failed.
	failDir := generateSandboxProject(ctx, t, "docker", "fail")
	failImage := buildProjectImage(ctx, t, failDir, "uboot-e2e-status-daemon-fail")
	failBox := startLabeled(ctx, t, failDir, failImage)
	if out, code := lifecycle(ctx, t, failBox, startCmd); code != 11 {
		t.Fatalf("fail policy without --privileged: exit %d, want 11:\n%s", code, out)
	}
	if status, _ := hostStatus(ctx, t, failDir, "daemon"); status.Result() != domain.SandboxStatusFailed {
		t.Errorf("result = %q, want failed", status.Result())
	}
	expectItem(t, doctorItem(ctx, t, failDir, itemID), domain.SeverityError, false)
}

// LH-FA-DEV-008: the egress start script reports `ok` with NET_ADMIN,
// `fallback` without it (policy warn) and `failed` (policy fail).
func TestE2E_SandboxDevcontainer_StatusFiles_Egress(t *testing.T) {
	t.Parallel()
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()
	const startCmd = "sudo sh /usr/local/bin/u-boot-egress-init"
	const itemID = "devcontainer.sandbox.fallback.egress"

	warnDir := generateEgressProject(ctx, t, "warn")
	warnImage := buildProjectImage(ctx, t, warnDir, "uboot-e2e-status-egress")
	runArgs := generatedRunArgs(t, warnDir)

	granted := startLabeled(ctx, t, warnDir, warnImage, runArgs...)
	if out, code := lifecycle(ctx, t, granted, startCmd); code != 0 {
		t.Fatalf("egress start script exit %d:\n%s", code, out)
	}
	if status, _ := hostStatus(ctx, t, warnDir, "egress"); status.Result() != domain.SandboxStatusOK {
		t.Errorf("with NET_ADMIN: result = %q, want ok", status.Result())
	}
	expectItem(t, doctorItem(ctx, t, warnDir, itemID), domain.SeverityOK, false)
	removeContainer(ctx, t, granted)

	denied := startLabeled(ctx, t, warnDir, warnImage)
	if out, code := lifecycle(ctx, t, denied, startCmd); code != 0 || !strings.Contains(out, "egress restriction unavailable") {
		t.Fatalf("warn policy without NET_ADMIN: exit %d:\n%s", code, out)
	}
	if status, _ := hostStatus(ctx, t, warnDir, "egress"); status.Result() != domain.SandboxStatusFallback {
		t.Errorf("without NET_ADMIN: result = %q, want fallback", status.Result())
	}
	expectItem(t, doctorItem(ctx, t, warnDir, itemID), domain.SeverityWarn, true)

	failDir := generateEgressProject(ctx, t, "fail")
	failImage := buildProjectImage(ctx, t, failDir, "uboot-e2e-status-egress-fail")
	failBox := startLabeled(ctx, t, failDir, failImage)
	if out, code := lifecycle(ctx, t, failBox, startCmd); code != 11 {
		t.Fatalf("fail policy without NET_ADMIN: exit %d, want 11:\n%s", code, out)
	}
	if status, _ := hostStatus(ctx, t, failDir, "egress"); status.Result() != domain.SandboxStatusFailed {
		t.Errorf("result = %q, want failed", status.Result())
	}
	expectItem(t, doctorItem(ctx, t, failDir, itemID), domain.SeverityError, false)
}

// ADR-0022 trust boundary against a real engine: the Inspector refuses every
// entry type but a regular file, caps the size, tells a missing file from a
// missing container, reads a stopped container, and picks the newest of
// several containers of the same folder.
func TestE2E_SandboxDevcontainer_StatusFiles_InspectorAgainstEngine(t *testing.T) {
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	folder := t.TempDir()
	insp := dockeradapter.NewInspector()

	if _, found, err := insp.FindContainer(ctx, folder); err != nil || found {
		t.Fatalf("no container yet: found=%v err=%v", found, err)
	}

	if out, code := dockerOutput(ctx, t, "pull", "-q", "docker.io/library/alpine:3"); code != 0 {
		t.Fatalf("docker pull alpine:3: %s", out)
	}
	first := startLabeled(ctx, t, folder, "docker.io/library/alpine:3")
	script := `set -e
mkdir -p ` + statusDirInContainer + `
cd ` + statusDirInContainer + `
printf 'u-boot-status 1\nresult=ok\nat=2026-10-17T18:52:54Z\n' > good
ln -s /etc/shadow link
mkfifo fifo
mkdir dir
head -c 100000 /dev/zero > big
printf 'u-boot-status 1\nresult=fallback\nat=2026-10-17T18:52:54Z\nnote=free text\n' > freetext`
	if out, code := dockerOutput(ctx, t, "exec", first, "sh", "-c", script); code != 0 {
		t.Fatalf("seed files: %s", out)
	}
	c, found, err := insp.FindContainer(ctx, folder)
	if err != nil || !found || !c.Running {
		t.Fatalf("FindContainer: %+v found=%v err=%v", c, found, err)
	}
	read := func(name string) ([]byte, error) {
		return insp.ReadFile(ctx, c.ID, statusDirInContainer+"/"+name, domain.SandboxStatusMaxBytes)
	}
	if raw, err := read("good"); err != nil || !strings.Contains(string(raw), "result=ok") {
		t.Errorf("regular file: %q, %v", raw, err)
	}
	for _, bad := range []string{"link", "fifo", "dir", "big"} {
		if raw, err := read(bad); !errors.Is(err, driven.ErrContainerFileInvalid) || raw != nil {
			t.Errorf("%s: got %q, %v; want ErrContainerFileInvalid and no content", bad, raw, err)
		}
	}
	if _, err := read("absent"); !errors.Is(err, driven.ErrContainerFileNotFound) {
		t.Errorf("absent file: %v, want ErrContainerFileNotFound", err)
	}
	if raw, err := read("freetext"); err != nil {
		t.Fatalf("freetext: %v", err)
	} else if _, perr := domain.ParseSandboxStatus(raw); !errors.Is(perr, domain.ErrInvalidSandboxStatus) {
		t.Errorf("a document with free text must not parse: %v", perr)
	}
	if _, err := insp.ReadFile(ctx, strings.Repeat("0", 64), statusDirInContainer+"/good", domain.SandboxStatusMaxBytes); err == nil || errors.Is(err, driven.ErrContainerFileNotFound) {
		t.Errorf("unknown container: %v, want an engine error", err)
	}

	// Stopped containers stay readable.
	if out, code := dockerOutput(ctx, t, "stop", "-t", "1", first); code != 0 {
		t.Fatalf("docker stop: %s", out)
	}
	c, found, err = insp.FindContainer(ctx, folder)
	if err != nil || !found || c.Running {
		t.Fatalf("stopped container: %+v found=%v err=%v", c, found, err)
	}
	if _, err := read("good"); err != nil {
		t.Errorf("stopped container not readable: %v", err)
	}

	// Several containers of one folder: the most recently created wins, running or not.
	time.Sleep(1500 * time.Millisecond)
	second := startLabeled(ctx, t, folder, "docker.io/library/alpine:3")
	c2, found, err := insp.FindContainer(ctx, folder)
	if err != nil || !found || !c2.Running || c2.ID == c.ID {
		t.Fatalf("newest container: %+v (old %s) found=%v err=%v", c2, c.ID, found, err)
	}
	if out, _ := dockerOutput(ctx, t, "inspect", "--format", "{{.Id}}", second); strings.TrimSpace(out) != c2.ID {
		t.Errorf("FindContainer returned %s, the newest container is %s", c2.ID, strings.TrimSpace(out))
	}
	if _, err := insp.ReadFile(ctx, c2.ID, statusDirInContainer+"/good", domain.SandboxStatusMaxBytes); !errors.Is(err, driven.ErrContainerFileNotFound) {
		t.Errorf("the new container has no file yet: %v, want ErrContainerFileNotFound", err)
	}

	// A remote engine is not queried.
	t.Setenv("DOCKER_HOST", "tcp://127.0.0.1:1")
	if _, found, err := dockeradapter.NewInspector().FindContainer(ctx, folder); !errors.Is(err, driven.ErrContainerEngineNotLocal) || found {
		t.Errorf("remote DOCKER_HOST: found=%v err=%v, want ErrContainerEngineNotLocal", found, err)
	}
}
