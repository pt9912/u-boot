//go:build docker

// End-to-end check of `u-boot devcontainer up` and `down` (LH-FA-DEV-010,
// ADR-0028 decision 9, ADR-0029): projects are generated with the production
// adapters, started with the production start service (launcher, inspector,
// filesystem, YAML codec) against the engine of the runner, and inspected with
// `docker inspect`. The S8 table runs the built binary, so the exit code of the
// process and the `--json` envelope are the observed result.
//
// The tests need the engine of the runner, the network (base images, the
// registry of the Dev Container features) and a project folder that engine and
// test share (`make test-docker` mounts `/tmp` under the same path).

package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	confirmadapter "github.com/pt9912/u-boot/internal/adapter/driven/confirm"
	dockeradapter "github.com/pt9912/u-boot/internal/adapter/driven/docker"
	fsadapter "github.com/pt9912/u-boot/internal/adapter/driven/fs"
	gitadapter "github.com/pt9912/u-boot/internal/adapter/driven/git"
	loggeradapter "github.com/pt9912/u-boot/internal/adapter/driven/logger"
	progressadapter "github.com/pt9912/u-boot/internal/adapter/driven/progress"
	yamladapter "github.com/pt9912/u-boot/internal/adapter/driven/yaml"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const (
	dcBaseImage  = "mcr.microsoft.com/devcontainers/base:debian"
	dcFolderLbl  = "devcontainer.local_folder"
	dcConfigLbl  = "devcontainer.config_file"
	dcTestBudget = 20 * time.Minute
)

// --- helpers -----------------------------------------------------------------

func dcDocker(args ...string) (string, error) {
	var stderr bytes.Buffer
	cmd := exec.Command("docker", args...)
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return strings.TrimSpace(string(out)), errors.New(err.Error() + ": " + strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(string(out)), nil
}

func dcMustDocker(t *testing.T, args ...string) string {
	t.Helper()
	out, err := dcDocker(args...)
	if err != nil {
		t.Fatalf("docker %s: %v", strings.Join(args, " "), err)
	}
	return out
}

// dcInspectDoc is the part of `docker inspect` of a container the tests read.
type dcInspectDoc struct {
	ID     string
	State  struct{ Running bool }
	Mounts []struct {
		Type, Name, Source, Destination string
	}
	HostConfig struct {
		Binds  []string
		Mounts []struct{ Source string }
	}
	Config struct {
		Image  string
		Env    []string
		Labels map[string]string
	}
}

func dcInspect(t *testing.T, target string) dcInspectDoc {
	t.Helper()
	raw := dcMustDocker(t, "inspect", target)
	var docs []dcInspectDoc
	if err := json.Unmarshal([]byte(raw), &docs); err != nil || len(docs) != 1 {
		t.Fatalf("docker inspect %s: %v (%d documents)", target, err, len(docs))
	}
	return docs[0]
}

func dcExists(id string) bool {
	_, err := dcDocker("inspect", "--type", "container", id)
	return err == nil
}

func (d dcInspectDoc) hasSocket() bool {
	for _, m := range d.Mounts {
		if strings.Contains(m.Source, "docker.sock") || strings.Contains(m.Destination, "docker.sock") {
			return true
		}
	}
	for _, b := range d.HostConfig.Binds {
		if strings.Contains(b, "docker.sock") {
			return true
		}
	}
	for _, m := range d.HostConfig.Mounts {
		if strings.Contains(m.Source, "docker.sock") {
			return true
		}
	}
	return false
}

// bindsFolder reports whether the container has the folder (or a path below
// it) of the host bind-mounted.
func (d dcInspectDoc) bindsFolder(folder string) bool {
	for _, m := range d.Mounts {
		if m.Type == "bind" && (m.Source == folder || strings.HasPrefix(m.Source, folder+"/")) {
			return true
		}
	}
	return false
}

// volumeAt returns the named volume mounted at destination.
func (d dcInspectDoc) volumeAt(destination string) string {
	for _, m := range d.Mounts {
		if m.Type == "volume" && m.Destination == destination {
			return m.Name
		}
	}
	return ""
}

func dcSuffix() string { return strconv.FormatInt(time.Now().UnixNano(), 36) }

// dcProjectDir returns a fresh project folder below the test's temp dir; its
// base name is unique, so the images `vsc-<name>-…` of the devcontainer CLI
// belong to this test, and registers the cleanup of everything the start
// creates for it: containers (found by the label), their volumes and the
// images.
func dcProjectDir(t *testing.T) string {
	t.Helper()
	requireDocker(t)
	dir := filepath.Join(t.TempDir(), "dcs"+dcSuffix())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	dcRegisterCleanup(t, dir)
	return dir
}

func dcRegisterCleanup(t *testing.T, dir string) {
	t.Helper()
	t.Cleanup(func() {
		ids, _ := dcDocker("ps", "-aq", "--filter", "label="+dcFolderLbl+"="+dir)
		var volumes []string
		for _, id := range strings.Fields(ids) {
			raw, err := dcDocker("inspect", "--format", `{{range .Mounts}}{{if eq .Type "volume"}}{{.Name}} {{end}}{{end}}`, id)
			if err == nil {
				volumes = append(volumes, strings.Fields(raw)...)
			}
		}
		if len(strings.Fields(ids)) > 0 {
			_, _ = dcDocker(append([]string{"rm", "-f", "-v"}, strings.Fields(ids)...)...)
		}
		if len(volumes) > 0 {
			_, _ = dcDocker(append([]string{"volume", "rm", "-f"}, volumes...)...)
		}
		if imgs, err := dcDocker("images", "-q", "--filter", "reference=vsc-"+filepath.Base(dir)+"-*"); err == nil && imgs != "" {
			_, _ = dcDocker(append([]string{"rmi", "-f"}, strings.Fields(imgs)...)...)
		}
	})
}

// dcGenerateProject produces a project (standard or sandbox profile) with the
// production adapters: `init --devcontainer`, the settings through `config
// set`, then `generate devcontainer`.
func dcGenerateProject(ctx context.Context, t *testing.T, sandbox bool, settings map[string]string) string {
	t.Helper()
	dir := dcProjectDir(t)
	fsys := fsadapter.New()
	yaml := yamladapter.New()
	logger := loggeradapter.New(os.Stderr, loggeradapter.FormatText, nil)
	initSvc := application.NewInitProjectService(fsys, yaml, gitadapter.New(),
		progressadapter.NewText(os.Stderr), confirmadapter.New(strings.NewReader(""), os.Stderr), logger)
	if _, err := initSvc.Init(ctx, driving.InitProjectRequest{
		BaseDir: dir, Name: "dcstart", SkipGit: true, Devcontainer: true, Sandbox: sandbox, SilenceProgress: true,
	}); err != nil {
		t.Fatalf("init --devcontainer (sandbox=%v): %v", sandbox, err)
	}
	cfgSvc := application.NewConfigService(fsys, yaml, logger)
	for path, value := range settings {
		p, err := domain.NewConfigPath(path)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := cfgSvc.Set(ctx, driving.ConfigSetRequest{BaseDir: dir, Path: p, Value: value}); err != nil {
			t.Fatalf("config set %s: %v", path, err)
		}
	}
	genSvc := application.NewGenerateService(fsys, yaml, logger)
	if _, err := genSvc.Generate(ctx, driving.GenerateRequest{BaseDir: dir, Artifact: domain.ArtifactDevcontainer}); err != nil {
		t.Fatalf("generate devcontainer: %v", err)
	}
	// The project folder is bind-mounted in the standard profile: the
	// container user must be able to write into it.
	if err := os.Chmod(dir, 0o777); err != nil {
		t.Fatal(err)
	}
	return dir
}

// dcWithoutNetAdmin replaces the `NET_ADMIN` capability in the generated
// runArgs by a capability every container has, so the egress start script
// meets a container without it.
func dcWithoutNetAdmin(t *testing.T, dir string) {
	t.Helper()
	path := filepath.Join(dir, ".devcontainer", "devcontainer.json")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(body, []byte(`"--cap-add=NET_ADMIN"`)) {
		t.Fatalf("devcontainer.json has no NET_ADMIN runArg:\n%s", body)
	}
	body = bytes.Replace(body, []byte(`"--cap-add=NET_ADMIN"`), []byte(`"--cap-add=NET_BIND_SERVICE"`), 1)
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatal(err)
	}
}

// dcHandWritten writes a hand-made project: devcontainer.json (and further
// files below .devcontainer) with the given content.
func dcHandWritten(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := dcProjectDir(t)
	for rel, content := range files {
		path := filepath.Join(dir, ".devcontainer", rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func dcService() *application.DevcontainerStartService {
	return application.NewDevcontainerStartService(fsadapter.New(), yamladapter.New(),
		dockeradapter.NewLauncher(), dockeradapter.NewInspector(),
		loggeradapter.New(os.Stderr, loggeradapter.FormatText, nil))
}

// dcStart is the outcome of one `Up` call.
type dcStart struct {
	resp      driving.DevcontainerUpResponse
	disclosed []driving.WarningEntry
	progress  string
	err       error
}

func dcUp(ctx context.Context, svc *application.DevcontainerStartService, dir string, recreate bool) dcStart {
	var out dcStart
	var progress bytes.Buffer
	out.resp, out.err = svc.Up(ctx, driving.DevcontainerUpRequest{
		BaseDir: dir, Recreate: recreate, ProgressSink: &progress,
		OnWarning: func(w driving.WarningEntry) { out.disclosed = append(out.disclosed, w) },
	})
	out.progress = progress.String()
	return out
}

func dcMustUp(ctx context.Context, t *testing.T, svc *application.DevcontainerStartService, dir string, recreate bool) dcStart {
	t.Helper()
	s := dcUp(ctx, svc, dir, recreate)
	if s.err != nil {
		t.Fatalf("Up: %v\nprogress:\n%s", s.err, s.progress)
	}
	return s
}

func dcExec(ctx context.Context, t *testing.T, id, user, script string) (string, int) {
	t.Helper()
	return dockerOutput(ctx, t, "exec", "-u", user, id, "sh", "-c", script)
}

func dcOwnerUID(t *testing.T, path string) uint32 {
	t.Helper()
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		t.Fatalf("no Stat_t for %s", path)
	}
	return st.Uid
}

// dcProjectImages returns the names (`repository:tag`) of the images
// `vsc-<project>-…` the devcontainer CLI built for the project folder. The
// name, not the ID: `--recreate` builds the image again.
func dcProjectImages(t *testing.T, dir string) map[string]bool {
	t.Helper()
	out := dcMustDocker(t, "images", "--format", "{{.Repository}}:{{.Tag}}", "--filter", "reference=vsc-"+filepath.Base(dir)+"-*")
	ids := map[string]bool{}
	for _, id := range strings.Fields(out) {
		ids[id] = true
	}
	return ids
}

// dcRequireImagesKept fails when an image name of before is gone.
func dcRequireImagesKept(t *testing.T, step string, before, after map[string]bool) {
	t.Helper()
	for id := range before {
		if !after[id] {
			t.Errorf("%s: the image %s of the devcontainer CLI is gone", step, id)
		}
	}
}

// dcVolumeHolds reports whether a file in the volume contains needle.
func dcVolumeHolds(ctx context.Context, t *testing.T, image, volume, needle string) bool {
	t.Helper()
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-i", "-u", "0", "--entrypoint", "sh",
		"-v", volume+":/scan:ro", image, "-c", "grep -rlF -f - /scan > /dev/null")
	cmd.Stdin = strings.NewReader(needle + "\n")
	err := cmd.Run()
	if err == nil {
		return true
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) && ee.ExitCode() == 1 {
		return false
	}
	t.Fatalf("scan of volume %s: %v", volume, err)
	return false
}

// --- S1: start of a generated project ---------------------------------------

// LH-FA-DEV-010, S1 / S3 / S9, standard profile: the generated project starts,
// the container user is the generated one, a catalogue feature is installed,
// the labels are the ones of the devcontainer CLI, the started container holds
// no engine socket (the project has no feature that mounts it) and no label of
// u-boot's own. With the default UID the file a
// container user creates in the bind-mounted project folder does not belong to
// UID 4321 (guard of the S3 test).
func TestDevcontainerUp_Docker_S1_S9_StandardProfile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, false, map[string]string{"devcontainer.features.node.enabled": "true"})
	svc := dcService()

	s := dcMustUp(ctx, t, svc, dir, false)
	doc := dcInspect(t, s.resp.ContainerID)
	if !doc.State.Running || s.resp.RemoteUser != "vscode" || s.resp.RemoteWorkspaceFolder != "/workspaces/"+filepath.Base(dir) {
		t.Errorf("running=%v response=%+v", doc.State.Running, s.resp)
	}
	if len(s.resp.Warnings) != 0 {
		t.Errorf("a standard project has no start script, got warnings %+v", s.resp.Warnings)
	}
	if len(s.disclosed) != 1 || s.disclosed[0].Code != "LH-FA-DEV-010" || !s.disclosed[0].SecurityRelaxation {
		t.Errorf("disclosures = %+v, want exactly the socket relaxation", s.disclosed)
	}
	if out, code := dcExec(ctx, t, s.resp.ContainerID, "vscode", "id -un && node --version"); code != 0 || !strings.HasPrefix(out, "vscode\n") {
		t.Errorf("remote user or feature `node` not usable (exit %d): %s", code, out)
	}

	// S9: the labels of the devcontainer CLI, nothing of u-boot.
	if got := doc.Config.Labels[dcFolderLbl]; got != dir {
		t.Errorf("label %s = %q, want the project folder %q", dcFolderLbl, got, dir)
	}
	if got := doc.Config.Labels[dcConfigLbl]; got != filepath.Join(dir, ".devcontainer", "devcontainer.json") {
		t.Errorf("label %s = %q", dcConfigLbl, got)
	}
	for key := range doc.Config.Labels {
		if lower := strings.ToLower(key); strings.Contains(lower, "u-boot") || strings.Contains(lower, "uboot") {
			t.Errorf("label %q is a label of u-boot's own", key)
		}
	}
	// The standard profile mounts the project folder (guard of the sandbox check).
	if !doc.bindsFolder(dir) {
		t.Errorf("the standard profile does not mount the project folder: %+v", doc.Mounts)
	}
	// ADR-0032 decision 5: a project without a feature that mounts the socket
	// mounts none (the counter-sample of
	// TestDevcontainerUp_Docker_StandardProfileDockerCLIMountsTheSocket).
	if doc.hasSocket() {
		t.Errorf("the started container holds the engine socket: %+v", doc.Mounts)
	}
	// Guard: FindContainer finds the container by the label value only.
	insp := dockeradapter.NewInspector()
	if c, ok, err := insp.FindContainer(ctx, dir); err != nil || !ok || c.ID != s.resp.ContainerID {
		t.Errorf("FindContainer(project) = %+v ok=%v err=%v", c, ok, err)
	}
	if _, ok, err := insp.FindContainer(ctx, "/w/"+filepath.Base(dir)); err != nil || ok {
		t.Errorf("FindContainer finds the container for another label value (ok=%v err=%v)", ok, err)
	}

	// S3 guard: the default UID is 1000, not 4321.
	if out, code := dcExec(ctx, t, s.resp.ContainerID, "vscode", "touch "+s.resp.RemoteWorkspaceFolder+"/probe-default"); code != 0 {
		t.Fatalf("create a file in the project folder: %s", out)
	}
	if uid := dcOwnerUID(t, filepath.Join(dir, "probe-default")); uid != 1000 {
		t.Errorf("owner of the file with the default UID = %d, want 1000", uid)
	}
}

// dcFeatureOptions sets the options of one rendered feature key in the
// generated devcontainer.json. The generator renders every key with an empty
// option object and u-boot has no setting for options, so the file is the
// only place to set one.
func dcFeatureOptions(t *testing.T, dir, key, options string) {
	t.Helper()
	path := filepath.Join(dir, ".devcontainer", "devcontainer.json")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	rendered := []byte(`"` + key + `": {}`)
	if !bytes.Contains(body, rendered) {
		t.Fatalf("devcontainer.json has no feature key %q with empty options:\n%s", key, body)
	}
	body = bytes.Replace(body, rendered, []byte(`"`+key+`": `+options), 1)
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatal(err)
	}
}

// LH-FA-DEV-006, ADR-0032 decision 5, standard profile with the catalogue
// feature `docker-cli`: the start mounts the engine socket of the host into
// the container. It is the measurement behind the statement of `up` that the
// generated configuration mounts the socket if a feature that does so is
// enabled; the project of the test above, without that feature, is the
// counter-sample. The project is generated and started with the production
// adapters, the container is read with `docker inspect`.
//
// Limit: the project differs from the generated one in one option. With the
// default option `moby: true` the install script of the feature stopped on
// the base image `mcr.microsoft.com/devcontainers/base:debian` (Debian trixie
// at the time; one run on 2026-10-22, "The 'moby' option is not supported on
// debian 'trixie'"), so the rendered key carries `{"moby": false}`. That the
// mount does not depend on the option is read from the static `mounts` of the
// feature metadata, not measured with `moby: true`.
func TestDevcontainerUp_Docker_StandardProfileDockerCLIMountsTheSocket(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, false, map[string]string{"devcontainer.features.docker-cli.enabled": "true"})
	dcFeatureOptions(t, dir, "ghcr.io/devcontainers/features/docker-outside-of-docker:1", `{"moby": false}`)

	s := dcMustUp(ctx, t, dcService(), dir, false)
	doc := dcInspect(t, s.resp.ContainerID)
	if !doc.State.Running {
		t.Fatalf("the container is not running: %+v", s.resp)
	}
	if !doc.hasSocket() {
		t.Errorf("the container started with the feature docker-cli holds no engine socket: %+v", doc.Mounts)
	}
	for _, m := range doc.Mounts {
		if strings.Contains(m.Source, "docker.sock") || strings.Contains(m.Destination, "docker.sock") {
			t.Logf("socket mount: type=%s source=%s destination=%s", m.Type, m.Source, m.Destination)
		}
	}
}

// LH-FA-DEV-010, S1 / S6 (volume) / S4 (volume), sandbox profile with the
// egress restriction: the start script reports `ok`, no fallback is disclosed,
// the credential placeholder of the host is in no volume of the container, a
// marker in the workspace volume survives `down`, a following `up` and `up
// --recreate`; `down` stops and keeps the container; `down` without a container
// and a second `down` end without an error.
func TestDevcontainerUp_Docker_S1_S4_SandboxProfileVolumes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	t.Setenv("GIT_TOKEN", sandboxCanaryToken)
	dir := dcGenerateProject(ctx, t, true, map[string]string{
		"devcontainer.sandbox.egress.enabled": "true",
		"devcontainer.sandbox.onUnavailable":  "warn",
	})
	svc := dcService()
	workspace := "/workspaces/dcstart"

	first := dcMustUp(ctx, t, svc, dir, false)
	if first.resp.RemoteUser != "vscode" || first.resp.RemoteWorkspaceFolder != workspace || len(first.resp.Warnings) != 0 {
		t.Fatalf("response = %+v, want user vscode, workspace %s and no fallback", first.resp, workspace)
	}
	status, _ := hostStatus(ctx, t, dir, "egress")
	if status.Result() != domain.SandboxStatusOK {
		t.Errorf("egress status = %v, want ok", status.Result())
	}
	doc := dcInspect(t, first.resp.ContainerID)
	if doc.hasSocket() {
		t.Errorf("the started sandbox container holds the engine socket: %+v", doc.Mounts)
	}
	if doc.bindsFolder(dir) {
		t.Errorf("the sandbox container has the project folder of the host mounted: %+v", doc.Mounts)
	}
	volume := doc.volumeAt(workspace)
	if volume == "" {
		t.Fatalf("no named volume at %s: %+v", workspace, doc.Mounts)
	}
	// The images `vsc-<project>-…` are there after `down` and after
	// `--recreate` (the manual states it); the guard is that the start created
	// some.
	images := dcProjectImages(t, dir)
	if len(images) == 0 {
		t.Fatalf("the start left no image vsc-%s-…, the check of their survival would be empty", filepath.Base(dir))
	}

	// S6, word "volume": the placeholder value of GIT_TOKEN, which the
	// start passed by name, is in no volume of the container; the guard puts it
	// into the workspace volume and the scan finds it.
	var volumes []string
	for _, m := range doc.Mounts {
		if m.Type == "volume" {
			volumes = append(volumes, m.Name)
		}
	}
	for _, v := range volumes {
		if dcVolumeHolds(ctx, t, doc.Config.Image, v, sandboxCanaryToken) {
			t.Errorf("volume %s holds the credential placeholder", v)
		}
	}
	id := first.resp.ContainerID
	if out, code := dockerOutput(ctx, t, "exec", "-u", "root", id, "sh", "-c", `printf '%s\n' "$1" > "$2"`, "sh", sandboxCanaryToken, workspace+"/leak.txt"); code != 0 {
		t.Fatalf("write the guard value: %s", out)
	}
	if !dcVolumeHolds(ctx, t, doc.Config.Image, volume, sandboxCanaryToken) {
		t.Errorf("guard: the scan does not find a value written into the workspace volume")
	}
	if out, code := dcExec(ctx, t, id, "root", "rm -f "+workspace+"/leak.txt && echo marker-1 > "+workspace+"/u-boot-marker"); code != 0 {
		t.Fatalf("write the marker: %s", out)
	}

	// S4: down stops and keeps; a second down is a no-op; up starts the same container.
	down, err := svc.Down(ctx, driving.DevcontainerDownRequest{BaseDir: dir})
	if err != nil || !down.Stopped || down.ContainerID != id {
		t.Fatalf("Down = %+v err=%v", down, err)
	}
	if after := dcInspect(t, id); after.State.Running {
		t.Errorf("the container still runs after down")
	}
	dcRequireImagesKept(t, "down", images, dcProjectImages(t, dir))
	again, err := svc.Down(ctx, driving.DevcontainerDownRequest{BaseDir: dir})
	if err != nil || again.Stopped || again.ContainerID != id {
		t.Errorf("second Down = %+v err=%v, want the same container and nothing stopped", again, err)
	}
	restarted := dcMustUp(ctx, t, svc, dir, false)
	if restarted.resp.ContainerID != id {
		t.Errorf("up after down started container %s, want the same %s", restarted.resp.ContainerID, id)
	}
	if out, code := dcExec(ctx, t, id, "root", "cat "+workspace+"/u-boot-marker"); code != 0 || strings.TrimSpace(out) != "marker-1" {
		t.Errorf("marker after down and up: %q (exit %d)", out, code)
	}

	// S4: --recreate gives a new container, the volume and the marker stay.
	recreated := dcMustUp(ctx, t, svc, dir, true)
	newID := recreated.resp.ContainerID
	if newID == id || dcExists(id) {
		t.Errorf("recreate: new ID %s, old container still there = %v", newID, dcExists(id))
	}
	dcRequireImagesKept(t, "recreate", images, dcProjectImages(t, dir))
	if got := dcInspect(t, newID).volumeAt(workspace); got != volume {
		t.Errorf("workspace volume after recreate = %q, want %q", got, volume)
	}
	if out, code := dcExec(ctx, t, newID, "root", "cat "+workspace+"/u-boot-marker"); code != 0 || strings.TrimSpace(out) != "marker-1" {
		t.Errorf("marker after recreate: %q (exit %d)", out, code)
	}

	// Guard of the marker test: without the volume the marker is gone.
	dcMustDocker(t, "rm", "-f", "-v", newID)
	dcMustDocker(t, "volume", "rm", "-f", volume)
	empty, err := svc.Down(ctx, driving.DevcontainerDownRequest{BaseDir: dir})
	if err != nil || empty.Stopped || empty.ContainerID != "" {
		t.Errorf("Down without a container = %+v err=%v, want nothing and no error", empty, err)
	}
	wiped := dcMustUp(ctx, t, svc, dir, false)
	if out, code := dcExec(ctx, t, wiped.resp.ContainerID, "root", "test -e "+workspace+"/u-boot-marker"); code == 0 {
		t.Errorf("guard: the marker survived the removal of its volume: %s", out)
	}
}

// LH-FA-DEV-010, S1 conditions of the wording: the fallback disclosure after
// `up` needs a u-boot.yaml in the project folder. With one, a container that
// lost the capability yields the disclosure in the wording of `up`; without
// one, `up` ends with exit 0 and discloses nothing.
func TestDevcontainerUp_Docker_S1_FallbackDisclosureNeedsUbootYAML(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, true, map[string]string{
		"devcontainer.sandbox.egress.enabled": "true",
		"devcontainer.sandbox.onUnavailable":  "warn",
	})
	dcWithoutNetAdmin(t, dir)
	svc := dcService()

	with := dcMustUp(ctx, t, svc, dir, false)
	if len(with.resp.Warnings) != 1 {
		t.Fatalf("warnings with u-boot.yaml = %+v, want the egress fallback", with.resp.Warnings)
	}
	w := with.resp.Warnings[0]
	if w.Code != "LH-FA-DEV-008" || !w.DegradationFallback || !strings.Contains(w.Message, "egress restriction is not active") ||
		strings.Contains(w.Message, "rebuilt") || strings.Contains(w.Message, "existing") {
		t.Errorf("warning = %+v", w)
	}

	if err := os.Remove(filepath.Join(dir, "u-boot.yaml")); err != nil {
		t.Fatal(err)
	}
	without := dcUp(ctx, svc, dir, false)
	if without.err != nil || len(without.resp.Warnings) != 0 {
		t.Errorf("without u-boot.yaml: err=%v warnings=%+v, want exit 0 and no disclosure", without.err, without.resp.Warnings)
	}
}

// --- S3: UID -----------------------------------------------------------------

// LH-FA-DEV-004, S3: with `devcontainer.user.uid: 4321` the container user has
// UID 4321 and a file it creates in the bind-mounted project folder belongs to
// 4321 on the host. The guard is the default-UID case in the standard-profile
// test (owner 1000 ≠ 4321).
func TestDevcontainerUp_Docker_S3_ConfiguredUID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, false, map[string]string{"devcontainer.user.uid": "4321"})
	s := dcMustUp(ctx, t, dcService(), dir, false)
	id := s.resp.ContainerID
	if out, code := dcExec(ctx, t, id, "vscode", "id -u"); code != 0 || strings.TrimSpace(out) != "4321" {
		t.Errorf("UID of the container user = %q (exit %d), want 4321", out, code)
	}
	if out, code := dcExec(ctx, t, id, "vscode", "touch "+s.resp.RemoteWorkspaceFolder+"/probe-4321"); code != 0 {
		t.Fatalf("create a file in the project folder: %s", out)
	}
	if uid := dcOwnerUID(t, filepath.Join(dir, "probe-4321")); uid != 4321 {
		t.Errorf("owner of the file = %d, want 4321", uid)
	}
}

// LH-FA-DEV-004, S3, observation: the result of `updateRemoteUserUID` when the
// project folder belongs to another UID than the generated one. Observed: the
// UID of the container user is not adapted to the owner of the project folder,
// so the user keeps UID 1000 and the file it creates belongs to 1000 on the
// host. The helper container runs as root; no test separates that from the
// cause, so neither this test nor the manual states one. The UID of the host
// comes from `devcontainer.user.uid`.
func TestDevcontainerUp_Docker_S3_ObserveUpdateRemoteUserUID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	if os.Geteuid() != 0 {
		t.Skip("giving the project folder to UID 2000 needs root; make test-docker runs the tests as root")
	}
	dir := dcGenerateProject(ctx, t, false, nil)
	if err := os.Chown(dir, 2000, 2000); err != nil {
		t.Fatalf("chown of the project folder to UID 2000: %v", err)
	}
	s := dcMustUp(ctx, t, dcService(), dir, false)
	id := s.resp.ContainerID
	out, code := dcExec(ctx, t, id, "vscode", "id -u")
	t.Logf("UID of the container user for a project folder owned by 2000: %q (exit %d)", strings.TrimSpace(out), code)
	if code != 0 || strings.TrimSpace(out) != "1000" {
		t.Errorf("UID of the container user = %q (exit %d), want 1000: updateRemoteUserUID does not adapt it to the folder owner", out, code)
	}
	if out, code := dcExec(ctx, t, id, "vscode", "touch "+s.resp.RemoteWorkspaceFolder+"/probe-observed"); code != 0 {
		t.Fatalf("create a file in the project folder: %s", out)
	}
	if uid := dcOwnerUID(t, filepath.Join(dir, "probe-observed")); uid != 1000 {
		t.Errorf("owner of the file = %d, want 1000", uid)
	}
}

// --- S8: failure picture ------------------------------------------------------

// ubootBinary builds the command into the temp dir of the test. The Go build
// cache makes every build after the first of a test run fast.
func ubootBinary(t *testing.T) string {
	t.Helper()
	_, file, _, _ := runtime.Caller(0)
	binary := filepath.Join(t.TempDir(), "u-boot")
	cmd := exec.Command("go", "build", "-o", binary, "./cmd/uboot")
	cmd.Dir = filepath.Join(filepath.Dir(file), "..", "..")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build ./cmd/uboot: %v\n%s", err, out)
	}
	return binary
}

// dcRun runs the binary in dir and returns stdout, stderr and the process exit
// code. env entries replace those of the test process.
func dcRun(ctx context.Context, t *testing.T, dir string, env map[string]string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.CommandContext(ctx, ubootBinary(t), args...)
	cmd.Dir = dir
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	code := 0
	if err != nil {
		var ee *exec.ExitError
		if !errors.As(err, &ee) {
			t.Fatalf("run %v: %v", args, err)
		}
		code = ee.ExitCode()
	}
	return stdout.String(), stderr.String(), code
}

// LH-FA-CLI-006, LH-FA-DEV-010, ADR-0028 decision 3, ADR-0029 decision 1, S8:
// the table "case → expected exit" of the real process. Under `--json` the
// `exitCode` of the envelope equals the exit code of the process.
func TestDevcontainerUp_Docker_S8_FailurePicture(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	const hand = `{"name": "s8", "image": "` + dcBaseImage + `"}`

	type testCase struct {
		name    string
		project func(t *testing.T) string
		env     func(t *testing.T) map[string]string
		want    int
		message string // a text the first error diagnostic must contain
	}
	handProject := func(files map[string]string) func(*testing.T) string {
		return func(t *testing.T) string { return dcHandWritten(t, files) }
	}
	noEnv := func(*testing.T) map[string]string { return nil }
	egressFail := func(removeYAML bool) func(*testing.T) string {
		return func(t *testing.T) string {
			dir := dcGenerateProject(ctx, t, true, map[string]string{
				"devcontainer.sandbox.egress.enabled": "true",
				"devcontainer.sandbox.onUnavailable":  "fail",
			})
			dcWithoutNetAdmin(t, dir)
			if removeYAML {
				if err := os.Remove(filepath.Join(dir, "u-boot.yaml")); err != nil {
					t.Fatal(err)
				}
			}
			return dir
		}
	}
	cases := []testCase{
		{"no docker binary on the PATH", handProject(map[string]string{"devcontainer.json": hand}),
			func(t *testing.T) map[string]string { return map[string]string{"PATH": t.TempDir()} }, 11, "docker binary not on PATH"},
		{"DOCKER_HOST is a remote engine", handProject(map[string]string{"devcontainer.json": hand}),
			func(*testing.T) map[string]string { return map[string]string{"DOCKER_HOST": "tcp://127.0.0.1:1"} }, 11, ""},
		{"no devcontainer.json", handProject(map[string]string{"README": "x"}), noEnv, 10, "devcontainer.json"},
		{"Dockerfile with RUN false", handProject(map[string]string{
			"devcontainer.json": `{"name": "s8", "build": {"dockerfile": "Dockerfile"}}`,
			"Dockerfile":        "FROM " + dcBaseImage + "\nRUN false\n",
		}), noEnv, 12, "false"},
		{"start script aborts with onUnavailable fail and u-boot.yaml", egressFail(false), noEnv, 11, "start script"},
		{"start script aborts with onUnavailable fail and no u-boot.yaml", egressFail(true), noEnv, 12, ""},
		{"postStartCommand exit 3", handProject(map[string]string{
			"devcontainer.json": `{"name": "s8", "image": "` + dcBaseImage + `", "postStartCommand": "exit 3"}`,
		}), noEnv, 12, "devcontainer up failed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := tc.project(t)
			stdout, stderr, code := dcRun(ctx, t, dir, tc.env(t), "devcontainer", "up", "--json")
			var env struct {
				Status      string                                  `json:"status"`
				Command     string                                  `json:"command"`
				Subcommand  string                                  `json:"subcommand"`
				ExitCode    int                                     `json:"exitCode"`
				Diagnostics []struct{ Level, Code, Message string } `json:"diagnostics"`
			}
			if err := json.Unmarshal([]byte(stdout), &env); err != nil {
				t.Fatalf("stdout is not one JSON envelope: %v\nstdout:\n%s\nstderr:\n%s", err, stdout, stderr)
			}
			if code != tc.want || env.ExitCode != code || env.Status != "error" || env.Command != "devcontainer" || env.Subcommand != "up" {
				t.Errorf("process exit %d, envelope %+v, want exit %d with status error", code, env, tc.want)
			}
			last := env.Diagnostics[len(env.Diagnostics)-1]
			t.Logf("last diagnostic (exit %d): %.600s", code, last.Message)
			if last.Level != "error" || !strings.Contains(last.Message, tc.message) {
				t.Errorf("last diagnostic = %+v, want an error containing %q", last, tc.message)
			}
		})
	}
}

// ADR-0028 decision 3, S8: a path with a character a `--mount` value cannot
// carry ends with exit 10 before any docker call.
func TestDevcontainerUp_Docker_S8_UnusablePathExits10(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := filepath.Join(t.TempDir(), "with,comma")
	if err := os.MkdirAll(filepath.Join(dir, ".devcontainer"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".devcontainer", "devcontainer.json"), []byte(`{"image": "`+dcBaseImage+`"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	stdout, stderr, code := dcRun(ctx, t, dir, nil, "devcontainer", "up", "--json")
	if code != 10 || !strings.Contains(stdout, `"exitCode":10`) || strings.Contains(stdout, "with,comma") {
		t.Errorf("exit %d, stdout=%q stderr=%q, want exit 10 without the path in the envelope", code, stdout, stderr)
	}
}

// ADR-0028 decision 4, S8: a project folder the engine cannot see ends with
// exit 12, because `--mount type=bind` aborts where `-v` creates the missing
// source on the daemon side (observed here, not assumed).
func TestDevcontainerUp_Docker_S8_FolderTheEngineCannotSeeExits12(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	requireDocker(t)
	missing := filepath.Join(os.TempDir(), "uboot-s8-invisible-"+dcSuffix())
	// The guard below makes the daemon create `missing` (owned by the daemon's
	// user); a container removes it, so the test leaves nothing behind whoever
	// runs it.
	t.Cleanup(func() {
		if os.RemoveAll(missing) == nil {
			return
		}
		_, _ = dcDocker("run", "--rm", "-v", filepath.Dir(missing)+":/parent", "--entrypoint", "rm",
			dcBaseImage, "-rf", "/parent/"+filepath.Base(missing))
	})

	_, err := dockeradapter.NewLauncher().Up(ctx, driven.DevcontainerUpOptions{ProjectFolder: missing})
	if !errors.Is(err, driven.ErrDevcontainerRuntime) || cli.ExitCode(err) != 12 {
		t.Fatalf("Up for a folder the engine cannot see: exit %d, err=%v, want exit 12", cli.ExitCode(err), err)
	}
	if _, statErr := os.Stat(missing); statErr == nil {
		t.Errorf("--mount created the missing source %s", missing)
	}

	// The guard: the same source with `-v` is created by the daemon.
	if out, err := dcDocker("run", "--rm", "-v", missing+":/probe", "--entrypoint", "true", dcBaseImage); err != nil {
		t.Fatalf("docker run -v with a missing source: %v (%s)", err, out)
	}
	if _, statErr := os.Stat(missing); statErr != nil {
		t.Errorf("guard: `-v` did not create the missing source: %v", statErr)
	}
}

// ADR-0028 decision 3, S8, observation: which output of the CLI is on stdout
// and which on stderr. The result line is the last non-empty line of stdout; it
// does not appear on stderr. The launcher passes stderr on as progress.
func TestDevcontainerUp_Docker_S8_ResultLineIsOnStdout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcHandWritten(t, map[string]string{"devcontainer.json": `{"name": "s8", "image": "` + dcBaseImage + `"}`})
	launcher := dockeradapter.NewLauncher()
	var progress bytes.Buffer
	image, err := launcher.EnsureImage(ctx, &progress)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"--mount", "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock",
		"--mount", "type=bind,source="+dir+",target="+dir,
		image, "devcontainer", "up", "--workspace-folder", dir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("devcontainer up: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	lastLine := strings.TrimSpace(lines[len(lines)-1])
	var result struct{ Outcome, ContainerID string }
	if err := json.Unmarshal([]byte(lastLine), &result); err != nil || result.Outcome != "success" || result.ContainerID == "" {
		t.Fatalf("the last stdout line is not the result line: %q (%v)", lastLine, err)
	}
	t.Logf("stdout lines: %d, stderr bytes: %d", len(lines), stderr.Len())
	if strings.Contains(stderr.String(), `"outcome":"success"`) {
		t.Errorf("the result line is on stderr as well:\n%s", stderr.String())
	}
	if stderr.Len() == 0 {
		t.Errorf("the CLI wrote nothing to stderr; the launcher's progress would stay empty")
	}
	t.Cleanup(func() { _, _ = dcDocker("rm", "-f", "-v", result.ContainerID) })
}

// Abort at the shipped command (ADR-0028 decision 4): SIGINT while the CLI runs
// a lifecycle command ends the process within the grace period of the launcher
// plus the removal of the helper container. The measured duration is logged;
// the manual states it as a single measurement.
func TestDevcontainerUp_Docker_AbortDuration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcHandWritten(t, map[string]string{
		"devcontainer.json": `{"name": "s8", "image": "` + dcBaseImage + `", "postStartCommand": "sleep 600"}`,
	})
	cmd := exec.CommandContext(ctx, ubootBinary(t), "devcontainer", "up", "--json")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	// The dev container appears once the CLI created it; its postStartCommand
	// (sleep) then keeps the helper container busy.
	deadline := time.Now().Add(5 * time.Minute)
	for {
		if ids, _ := dcDocker("ps", "-q", "--filter", "label="+dcFolderLbl+"="+dir); ids != "" {
			break
		}
		select {
		case err := <-done:
			t.Fatalf("the command ended before the dev container appeared: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
		default:
		}
		if time.Now().After(deadline) {
			_ = cmd.Process.Kill()
			t.Fatal("the dev container did not appear")
		}
		time.Sleep(500 * time.Millisecond)
	}
	time.Sleep(2 * time.Second) // the lifecycle command is running
	started := time.Now()
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}
	err := <-done
	took := time.Since(started)
	t.Logf("abort took %s", took.Round(100*time.Millisecond))
	var ee *exec.ExitError
	if !errors.As(err, &ee) || ee.ExitCode() == 0 {
		t.Errorf("the aborted command ended with %v, want a non-zero exit", err)
	}
	if took > 30*time.Second {
		t.Errorf("the abort took %s, want at most the 10 s grace plus the removal", took)
	}
	if ids, _ := dcDocker("ps", "-aq", "--filter", "label="+dcFolderLbl+"="+dir); ids == "" {
		t.Errorf("the dev container the CLI had created is gone after the abort")
	}
}
