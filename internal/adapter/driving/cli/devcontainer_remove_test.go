package cli_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli/jsontestutil"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const (
	removeVolA = "demo-workspace-abc"
	removeVolB = "demo-docker-abc"
)

// removeCases are the results `remove` renders: the text on the standard
// output, and the fields of `data`.
func removeCases() []struct {
	name    string
	resp    driving.DevcontainerRemoveResponse
	volumes bool
	human   string
	data    map[string]any
} {
	return []struct {
		name    string
		resp    driving.DevcontainerRemoveResponse
		volumes bool
		human   string
		data    map[string]any
	}{
		{
			"running container, no volumes flag",
			driving.DevcontainerRemoveResponse{ContainerID: testContainerID, Stopped: true, Removed: true},
			false, "Devcontainer removed: container 0123456789ab (stopped first)\n",
			map[string]any{"containerId": testContainerID, "stopped": true, "removed": true, "removedVolumeNames": []any{}},
		},
		{
			"stopped container, no volumes flag",
			driving.DevcontainerRemoveResponse{ContainerID: testContainerID, Removed: true},
			false, "Devcontainer removed: container 0123456789ab\n",
			map[string]any{"containerId": testContainerID, "stopped": false, "removed": true, "removedVolumeNames": []any{}},
		},
		{
			"volumes removed",
			driving.DevcontainerRemoveResponse{ContainerID: testContainerID, Stopped: true, Removed: true, RemovedVolumeNames: []string{removeVolB, removeVolA}},
			true, "Devcontainer removed: container 0123456789ab (stopped first)\nVolumes removed: 2\n  - " + removeVolB + "\n  - " + removeVolA + "\n",
			map[string]any{"containerId": testContainerID, "stopped": true, "removed": true, "removedVolumeNames": []any{removeVolB, removeVolA}},
		},
		{
			"volumes flag, container without a volume",
			driving.DevcontainerRemoveResponse{ContainerID: testContainerID, Removed: true},
			true, "Devcontainer removed: container 0123456789ab\nNo volume removed (the container held no volume)\n",
			map[string]any{"containerId": testContainerID, "stopped": false, "removed": true, "removedVolumeNames": []any{}},
		},
		{
			"no container",
			driving.DevcontainerRemoveResponse{},
			false, "No devcontainer container found for this project; nothing to remove\n",
			map[string]any{"containerId": "", "stopped": false, "removed": false, "removedVolumeNames": []any{}},
		},
		{
			"no container, volumes flag",
			driving.DevcontainerRemoveResponse{},
			true, "No devcontainer container found for this project; nothing to remove\nVolumes of a container that no longer exists are not removed by u-boot.\n",
			map[string]any{"containerId": "", "stopped": false, "removed": false, "removedVolumeNames": []any{}},
		},
	}
}

// SPEC-058, output: the text, `--quiet` (no result line), and `--json` (the
// four fields of `data`, no field optional, `removedVolumeNames` always an
// array, status ok and no diagnostic, `--quiet` changes nothing).
func TestExecute_DevcontainerRemove_Output(t *testing.T) {
	t.Parallel()
	for _, tc := range removeCases() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			args := []string{"devcontainer", "remove"}
			if tc.volumes {
				args = append(args, "--volumes", "--yes")
			}
			newUC := func() *fakeDevcontainerUseCase { return &fakeDevcontainerUseCase{removeResp: tc.resp} }
			uc := newUC()
			out, errOut, err := quietRun(t, newAppWithDevcontainer(uc), args...)
			if err != nil || out != tc.human || errOut != "" || uc.removeReq.BaseDir != "/tmp/x/demo" {
				t.Errorf("human: err=%v stdout=%q stderr=%q req=%+v", err, out, errOut, uc.removeReq)
			}
			out, errOut, err = quietRun(t, newAppWithDevcontainer(newUC()), append([]string{"--quiet"}, args...)...)
			if err != nil || out != "" || errOut != "" {
				t.Errorf("--quiet: err=%v stdout=%q stderr=%q, want nothing", err, out, errOut)
			}
			jsonOut, _, err := quietRun(t, newAppWithDevcontainer(newUC()), append([]string{"--json"}, args...)...)
			if err != nil {
				t.Fatal(err)
			}
			opts := []jsontestutil.AssertOption{jsontestutil.WithCommand("devcontainer"), jsontestutil.WithSubcommand("remove"), jsontestutil.WithExitCode(0)}
			for key, value := range tc.data {
				opts = append(opts, jsontestutil.WithDataKeyPresent(key, value))
			}
			jsontestutil.AssertMinimalEnvelope(t, []byte(jsonOut), opts...)
			var env struct {
				Status      string          `json:"status"`
				Diagnostics []any           `json:"diagnostics"`
				Data        json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal([]byte(jsonOut), &env); err != nil {
				t.Fatal(err)
			}
			var keys map[string]json.RawMessage
			if err := json.Unmarshal(env.Data, &keys); err != nil || len(keys) != 4 {
				t.Errorf("data has %d keys (%v), want the four fields:\n%s", len(keys), err, jsonOut)
			}
			if env.Status != "ok" || len(env.Diagnostics) != 0 || !strings.Contains(jsonOut, `"diagnostics":[]`) {
				t.Errorf("remove raises no diagnostic:\n%s", jsonOut)
			}
			quietJSON, _, _ := quietRun(t, newAppWithDevcontainer(newUC()), append([]string{"--quiet", "--json"}, args...)...)
			if quietJSON != jsonOut {
				t.Errorf("--quiet --json differs from --json:\n%s\n%s", quietJSON, jsonOut)
			}
		})
	}
}

// SPEC-058: the flags reach the request: `--volumes`, `--yes` and
// `--no-interactive` (root options, before or after the subcommand), `--json`
// as the silence of the confirmer.
func TestExecute_DevcontainerRemove_FlagsReachTheRequest(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		args []string
		want driving.DevcontainerRemoveRequest
	}{
		{"none", []string{"devcontainer", "remove"}, driving.DevcontainerRemoveRequest{BaseDir: "/tmp/x/demo"}},
		{"volumes", []string{"devcontainer", "remove", "--volumes"}, driving.DevcontainerRemoveRequest{BaseDir: "/tmp/x/demo", Volumes: true}},
		{"yes after", []string{"devcontainer", "remove", "--volumes", "--yes"}, driving.DevcontainerRemoveRequest{BaseDir: "/tmp/x/demo", Volumes: true, AssumeYes: true}},
		{"yes before", []string{"--yes", "devcontainer", "remove", "--volumes"}, driving.DevcontainerRemoveRequest{BaseDir: "/tmp/x/demo", Volumes: true, AssumeYes: true}},
		{"no-interactive", []string{"--no-interactive", "devcontainer", "remove", "--volumes"}, driving.DevcontainerRemoveRequest{BaseDir: "/tmp/x/demo", Volumes: true, NonInteractive: true}},
		{"json", []string{"--json", "devcontainer", "remove", "--volumes"}, driving.DevcontainerRemoveRequest{BaseDir: "/tmp/x/demo", Volumes: true, SilenceConfirmer: true}},
		{"yes without volumes", []string{"--yes", "devcontainer", "remove"}, driving.DevcontainerRemoveRequest{BaseDir: "/tmp/x/demo", AssumeYes: true}},
		{"instance", []string{"devcontainer", "remove", "--instance", "a1"}, driving.DevcontainerRemoveRequest{BaseDir: "/tmp/x/demo", Instance: "a1"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			uc := &fakeDevcontainerUseCase{}
			if _, _, err := quietRun(t, newAppWithDevcontainer(uc), tc.args...); err != nil {
				t.Fatal(err)
			}
			if uc.removeCalls != 1 || uc.removeReq != tc.want {
				t.Errorf("request = %+v after %d calls, want %+v after 1", uc.removeReq, uc.removeCalls, tc.want)
			}
		})
	}
}

// SPEC-058, usage (exit 2), before the use case runs: arguments, `--dry-run`,
// `--diff` and `--yes` together with `--no-interactive`, with and without
// `--volumes`; with `--json` each is an envelope with the code of the
// requirement that rejects it. `--json` itself is not rejected.
func TestExecute_DevcontainerRemove_UsageErrors(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		args []string
		code string
		is   error
	}{
		{"dry-run", []string{"devcontainer", "remove", "--dry-run"}, "LH-FA-CLI-006", cli.ErrDevcontainerPreviewNotApplicable},
		{"diff", []string{"devcontainer", "remove", "--diff"}, "LH-FA-CLI-006", cli.ErrDevcontainerPreviewNotApplicable},
		{"extra argument", []string{"devcontainer", "remove", "extra"}, "LH-FA-CLI-006", nil},
		{"yes and no-interactive", []string{"--yes", "--no-interactive", "devcontainer", "remove"}, "LH-FA-CLI-005A", cli.ErrConflictingModeFlags},
		{"yes and no-interactive with volumes", []string{"--yes", "--no-interactive", "devcontainer", "remove", "--volumes"}, "LH-FA-CLI-005A", cli.ErrConflictingModeFlags},
		{"unknown flag", []string{"devcontainer", "remove", "--purge"}, "", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			uc := &fakeDevcontainerUseCase{}
			out, errOut, err := quietRun(t, newAppWithDevcontainer(uc), tc.args...)
			if err == nil || cli.ExitCode(err) != 2 || uc.removeCalls != 0 || out != "" || errOut != "" {
				t.Fatalf("err=%v exit=%d calls=%d stdout=%q stderr=%q, want exit 2 and no call", err, cli.ExitCode(err), uc.removeCalls, out, errOut)
			}
			if tc.is != nil && !errors.Is(err, tc.is) {
				t.Errorf("err = %v, want %v", err, tc.is)
			}
			if tc.code == "" {
				return
			}
			jsonOut, _, err := quietRun(t, newAppWithDevcontainer(&fakeDevcontainerUseCase{}), append([]string{"--json"}, tc.args...)...)
			if cli.ExitCode(err) != 2 {
				t.Fatalf("--json: exit=%d, want 2", cli.ExitCode(err))
			}
			jsontestutil.AssertMinimalEnvelope(t, []byte(jsonOut), jsontestutil.WithCommand("devcontainer"), jsontestutil.WithSubcommand("remove"),
				jsontestutil.WithExitCode(2), jsontestutil.WithExpectedCodes(tc.code))
		})
	}
	uc := &fakeDevcontainerUseCase{removeResp: driving.DevcontainerRemoveResponse{}}
	if _, _, err := quietRun(t, newAppWithDevcontainer(uc), "--json", "devcontainer", "remove"); err != nil || uc.removeCalls != 1 {
		t.Errorf("--json was rejected: err=%v calls=%d", err, uc.removeCalls)
	}
}

// SPEC-058, exit codes of the table in both classifiers: `ExitCode` and the
// diagnostic mapping of `remove`; the code of the diagnostic is the one of the
// requirement for 10, 11 and 12, also for the missing release, and the path of
// the project stays out of the message. The filesystem class (exit 14) does
// not occur: an error of that class is not mapped to its code by `remove`.
func TestExecute_DevcontainerRemove_ExitCodesInBothClassifiers(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		err  error
		exit int
		code string
	}{
		{"declined", fmt.Errorf("devcontainer remove: %w: the confirmation was declined; /tmp/x/demo", driving.ErrConfirmationRequired), 10, "LH-FA-DEV-010"},
		{"answer unreadable", fmt.Errorf("devcontainer remove: %w: boom", driving.ErrConfirmerUnavailable), 10, "LH-FA-DEV-010"},
		{"instance name", fmt.Errorf("%w: Abc", driving.ErrDevcontainerInstanceInvalid), 10, "LH-FA-DEV-010"},
		{"only instances", fmt.Errorf("%w; pass --instance <name>", driving.ErrDevcontainerOnlyInstances), 10, "LH-FA-DEV-010"},
		{"docker unavailable", fmt.Errorf("docker binary not on PATH: %w", driven.ErrDockerUnavailable), 11, "LH-FA-DEV-010"},
		{"engine not local", driven.ErrContainerEngineNotLocal, 11, "LH-FA-DEV-010"},
		{"runtime", fmt.Errorf("docker volume rm failed: %w", driven.ErrDevcontainerRuntime), 12, "LH-FA-DEV-010"},
		{"other", errors.New("boom"), 1, "LH-FA-CLI-006"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := cli.ExitCode(tc.err); got != tc.exit {
				t.Errorf("ExitCode(%v) = %d, want %d", tc.err, got, tc.exit)
			}
			out, errOut, err := quietRun(t, newAppWithDevcontainer(&fakeDevcontainerUseCase{removeErr: tc.err}), "devcontainer", "remove", "--volumes")
			if !errors.Is(err, tc.err) || cli.ExitCode(err) != tc.exit || out != "" || errOut != "" {
				t.Fatalf("human: err=%v exit=%d stdout=%q stderr=%q", err, cli.ExitCode(err), out, errOut)
			}
			out, _, err = quietRun(t, newAppWithDevcontainer(&fakeDevcontainerUseCase{removeErr: tc.err}), "--json", "devcontainer", "remove", "--volumes")
			if cli.ExitCode(err) != tc.exit {
				t.Fatalf("json: exit=%d, want %d", cli.ExitCode(err), tc.exit)
			}
			jsontestutil.AssertMinimalEnvelope(t, []byte(out),
				jsontestutil.WithCommand("devcontainer"), jsontestutil.WithSubcommand("remove"),
				jsontestutil.WithExitCode(tc.exit), jsontestutil.WithExpectedCodes(tc.code), jsontestutil.WithDataKeyAbsent("containerId"))
			if strings.Contains(out, "/tmp/x/demo") {
				t.Errorf("the project path leaks into the diagnostic:\n%s", out)
			}
		})
	}
	// exit 14 does not occur: the diagnostic of `remove` has no code of the
	// filesystem class, and nothing of the project is read.
	fsErr := fmt.Errorf("read: %w", driving.ErrDevcontainerFileSystem)
	out, _, err := quietRun(t, newAppWithDevcontainer(&fakeDevcontainerUseCase{removeErr: fsErr}), "--json", "devcontainer", "remove")
	if err == nil || strings.Contains(out, "LH-NFA-REL-003") {
		t.Errorf("remove maps the filesystem class: err=%v\n%s", err, out)
	}
}

// SPEC-058: the command reads no project file; a project folder with an
// invalid `u-boot.yaml` and an invalid `.devcontainer/devcontainer.json` does
// not end with exit 14, and the adapter reads neither.
func TestExecute_DevcontainerRemove_ReadsNoProjectFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	for rel, body := range map[string]string{"u-boot.yaml": "{ invalid: [", ".devcontainer/devcontainer.json": "{ not json"} {
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	uc := &fakeDevcontainerUseCase{removeResp: driving.DevcontainerRemoveResponse{ContainerID: testContainerID, Removed: true}}
	app := newApp(&fakeInitUseCase{}, cli.WithDevcontainerUseCase(uc), cli.WithGetwd(func() (string, error) { return dir, nil }))
	_, _, err := quietRun(t, app, "devcontainer", "remove")
	if err != nil || cli.ExitCode(err) == 14 || uc.removeReq.BaseDir != dir {
		t.Errorf("err=%v exit=%d req=%+v, want success on the folder", err, cli.ExitCode(err), uc.removeReq)
	}
}

// The command is wired: without its use case it fails with a wiring error.
func TestExecute_DevcontainerRemove_NotWired(t *testing.T) {
	t.Parallel()
	_, _, err := quietRun(t, newApp(&fakeInitUseCase{}, cli.WithGetwd(sandboxGetwd)), "devcontainer", "remove")
	if err == nil || !strings.Contains(err.Error(), "not wired") || cli.ExitCode(err) != 1 {
		t.Errorf("err=%v exit=%d", err, cli.ExitCode(err))
	}
}

// SPEC-058, help text: it separates `remove` from `u-boot remove` and `down`,
// warns that `--volumes` deletes data, names the question, and carries no
// requirement, ADR or specification identifier; the group help lists it, and
// the help of `down` points at it.
func TestExecute_DevcontainerRemove_Help(t *testing.T) {
	t.Parallel()
	text := func(args ...string) string {
		out, errOut, err := quietRun(t, newAppWithDevcontainer(&fakeDevcontainerUseCase{}), args...)
		if err != nil {
			t.Fatalf("%v: %v", args, err)
		}
		return strings.Join(strings.Fields(out+errOut), " ")
	}
	cases := []struct {
		args  []string
		parts []string
	}{
		{[]string{"devcontainer", "remove", "--help"}, []string{
			`It is not "u-boot remove"`, `and not "u-boot devcontainer down"`,
			"THE DATA IN THESE VOLUMES IS PERMANENTLY DELETED",
			"asks for confirmation on the standard error stream",
			"--no-interactive exit 10, nothing is stopped or removed",
			"with no container the command ends with exit 0 and asks nothing",
			"A volume that another container holds cannot be removed",
			"--quiet holds back the result lines; the question and every error stay.",
			"12 docker stop, docker rm or docker volume rm failed, or a volume remains",
		}},
		{[]string{"devcontainer", "--help"}, []string{"u-boot devcontainer remove [--volumes] remove the container, with --volumes its volumes"}},
		{[]string{"devcontainer", "down", "--help"}, []string{`"u-boot devcontainer remove" removes the container`}},
	}
	for _, tc := range cases {
		got := text(tc.args...)
		for _, id := range []string{"LH-FA", "ADR-", "SPEC-", "ARC-"} {
			if strings.Contains(got, id) {
				t.Errorf("%v: the help carries %q", tc.args, id)
			}
		}
		for _, part := range tc.parts {
			if !strings.Contains(got, part) {
				t.Errorf("%v: the help lacks %q:\n%s", tc.args, part, got)
			}
		}
	}
}
