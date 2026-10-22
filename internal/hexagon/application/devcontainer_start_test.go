package application_test

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const devcontainerTestBaseDir = "/tmp/u-boot-devcontainer-test/demo"

// fakeLauncher is the [driven.DevcontainerLauncher] test double: it records
// every call and returns the configured result or error. onUp runs at the
// start of Up, so a test can look at what happened before the call.
type fakeLauncher struct {
	mu       sync.Mutex
	upOpts   []driven.DevcontainerUpOptions
	upResult driven.DevcontainerUpResult
	upErr    error
	onUp     func()
	stopIDs  []string
	stopErr  error
}

func (f *fakeLauncher) Up(_ context.Context, opts driven.DevcontainerUpOptions) (driven.DevcontainerUpResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.upOpts = append(f.upOpts, opts)
	if f.onUp != nil {
		f.onUp()
	}
	return f.upResult, f.upErr
}

func (f *fakeLauncher) Stop(_ context.Context, containerID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.stopIDs = append(f.stopIDs, containerID)
	return f.stopErr
}

// devcontainerSvc builds the service on a fake filesystem that holds
// `.devcontainer/devcontainer.json` (when config is not empty) and u-boot.yaml
// (when ubootYAML is not empty).
func devcontainerSvc(t *testing.T, config, ubootYAML string, l *fakeLauncher, insp *fakeInspector) (*application.DevcontainerStartService, *fakeFS) {
	t.Helper()
	fs := newFakeFS()
	fs.markDirExists(devcontainerTestBaseDir)
	if config != "" {
		if err := fs.WriteFile(filepath.Join(devcontainerTestBaseDir, ".devcontainer", "devcontainer.json"), []byte(config), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if ubootYAML != "" {
		if err := fs.WriteFile(filepath.Join(devcontainerTestBaseDir, "u-boot.yaml"), []byte(ubootYAML), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var launcher driven.DevcontainerLauncher
	if l != nil {
		launcher = l
	}
	var inspector driven.DevcontainerInspector
	if insp != nil {
		inspector = insp
	}
	return application.NewDevcontainerStartService(fs, &fakeYAML{}, launcher, inspector, nil), fs
}

const plainConfig = `{"name": "demo", "image": "debian"}`

func newLauncher() *fakeLauncher {
	return &fakeLauncher{upResult: driven.DevcontainerUpResult{
		ContainerID: fallbackContainerID, RemoteUser: "vscode", RemoteWorkspaceFolder: "/workspaces/demo",
	}}
}

// ADR-0028 decision 4: the names come from the text of devcontainer.json
// as `${localEnv:<NAME>}`, with or without a default.
func TestLocalEnvNames(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		text string
		want []string
	}{
		{"none", `{"image": "debian"}`, nil},
		{"one", `{"remoteEnv": {"GIT_TOKEN": "${localEnv:GIT_TOKEN}"}}`, []string{"GIT_TOKEN"}},
		{"with default", `"${localEnv:HOME_DIR:/root}"`, []string{"HOME_DIR"}},
		{"sorted and distinct", `"${localEnv:B}" "${localEnv:A}" "${localEnv:B:x}"`, []string{"A", "B"}},
		{"two on one line", `"${localEnv:A}/${localEnv:B}"`, []string{"A", "B"}},
		{"other variable forms are not read", `"${env:A}" "${containerEnv:B}" "$localEnv:C" "${localEnvironment:D}"`, nil},
		{"empty name is kept for validation", `"${localEnv:}"`, []string{""}},
		{"a comment counts as text", "// ${localEnv:IN_COMMENT}\n{}", []string{"IN_COMMENT"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := application.LocalEnvNamesForTest(tc.text)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("names = %q, want %q", got, tc.want)
			}
		})
	}
}

// LH-FA-DEV-010: Up hands the project folder, the referenced names, the
// recreate flag and the progress sink to the launcher and returns the result.
func TestDevcontainerUp_PassesRequestToLauncher(t *testing.T) {
	t.Parallel()
	l := newLauncher()
	svc, _ := devcontainerSvc(t, `{"remoteEnv": {"T": "${localEnv:GIT_TOKEN}", "U": "${localEnv:AGENT_KEY}"}}`, "", l, newFakeInspector(true))
	resp, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: devcontainerTestBaseDir, Recreate: true, ProgressSink: io.Discard})
	if err != nil {
		t.Fatalf("Up: %v", err)
	}
	if len(l.upOpts) != 1 {
		t.Fatalf("launcher Up calls = %d, want 1", len(l.upOpts))
	}
	got := l.upOpts[0]
	if got.ProjectFolder != devcontainerTestBaseDir || !got.Recreate || got.ReadOnlyProject || got.Progress != io.Discard ||
		!reflect.DeepEqual(got.EnvNames, []string{"AGENT_KEY", "GIT_TOKEN"}) {
		t.Errorf("launcher options = %+v", got)
	}
	if resp.ContainerID != fallbackContainerID || resp.RemoteUser != "vscode" || resp.RemoteWorkspaceFolder != "/workspaces/demo" {
		t.Errorf("response = %+v", resp)
	}
}

// ADR-0028 decision 4 / ADR-0029 decision 3 / ADR-0032 decision 5: the socket
// disclosure reaches OnWarning before the launcher starts, carries its parts
// (the helper, the generated configuration with a feature that mounts the
// socket, the hand-written configuration, the initializeCommand), and is a
// security relaxation.
func TestDevcontainerUp_SocketDisclosureBeforeStart(t *testing.T) {
	t.Parallel()
	l := newLauncher()
	var warnings []driving.WarningEntry
	l.onUp = func() {
		if len(warnings) != 1 {
			t.Errorf("warnings raised before the launcher started = %d, want 1", len(warnings))
		}
	}
	svc, _ := devcontainerSvc(t, plainConfig, "", l, newFakeInspector(true))
	_, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{
		BaseDir:   devcontainerTestBaseDir,
		OnWarning: func(w driving.WarningEntry) { warnings = append(warnings, w) },
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 1 {
		t.Fatalf("warnings = %+v, want exactly the disclosure", warnings)
	}
	w := warnings[0]
	if w.Code != "LH-FA-DEV-010" || w.Level != "warn" || !w.SecurityRelaxation || w.DegradationFallback {
		t.Errorf("disclosure = %+v, want code LH-FA-DEV-010, level warn, SecurityRelaxation only", w)
	}
	for _, part := range []string{
		"mounts the host container-engine socket into a short-lived helper container", "root-equivalent access",
		"u-boot generates mounts it if a feature that does so is enabled, such as the built-in feature docker-cli (the sandbox profile rejects that feature)",
		"hand-written .devcontainer/devcontainer.json can mount it",
		"initializeCommand declared there runs with this access",
	} {
		if !strings.Contains(strings.Join(strings.Fields(w.Message), " "), part) {
			t.Errorf("disclosure %q lacks %q", w.Message, part)
		}
	}
	if strings.Contains(w.Message, "does not receive it") {
		t.Errorf("disclosure %q keeps the unconditional claim of ADR-0028", w.Message)
	}
	// ADR-0032 decision 5: the generated configuration is not said to
	// leave the socket out, and the text does not say "only" (a feature from
	// an external source can mount the socket without u-boot seeing it).
	if strings.Contains(w.Message, "does not mount it") || strings.Contains(" "+w.Message+" ", " only ") {
		t.Errorf("disclosure %q keeps the old claim or says \"only\"", w.Message)
	}
}

// Up validates before it raises the disclosure and before the launcher runs:
// a request that fails here started no helper container.
func TestDevcontainerUp_ValidationFailuresRaiseNoDisclosureAndStartNothing(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		baseDir string
		config  string
		setup   func(fs *fakeFS)
		want    error
	}{
		{"no configuration", devcontainerTestBaseDir, "", nil, driving.ErrDevcontainerConfigMissing},
		{"comma in the path", "/tmp/a,b", plainConfig, nil, driving.ErrDevcontainerProjectFolderInvalid},
		{"double quote in the path", "/tmp/a\"b", plainConfig, nil, driving.ErrDevcontainerProjectFolderInvalid},
		{"line break in the path", "/tmp/a\nb", plainConfig, nil, driving.ErrDevcontainerProjectFolderInvalid},
		{"relative path", "demo", plainConfig, nil, driving.ErrDevcontainerProjectFolderInvalid},
		{"name that is no plain name", devcontainerTestBaseDir, `"${localEnv:A-B}"`, nil, driving.ErrDevcontainerEnvNameInvalid},
		{"empty name", devcontainerTestBaseDir, `"${localEnv:}"`, nil, driving.ErrDevcontainerEnvNameInvalid},
		{"probe fails", devcontainerTestBaseDir, plainConfig, func(fs *fakeFS) {
			fs.failExistsOn = filepath.Join(devcontainerTestBaseDir, ".devcontainer", "devcontainer.json")
			fs.failExistsErr = errors.New("disk gone")
		}, driving.ErrDevcontainerFileSystem},
		{"read fails", devcontainerTestBaseDir, plainConfig, func(fs *fakeFS) {
			fs.failReadOn = filepath.Join(devcontainerTestBaseDir, ".devcontainer", "devcontainer.json")
			fs.failReadErr = errors.New("io error")
		}, driving.ErrDevcontainerFileSystem},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := newLauncher()
			svc, fs := devcontainerSvc(t, tc.config, "", l, newFakeInspector(true))
			if tc.setup != nil {
				tc.setup(fs)
			}
			raised := 0
			_, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: tc.baseDir, OnWarning: func(driving.WarningEntry) { raised++ }})
			if !errors.Is(err, tc.want) {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
			if len(l.upOpts) != 0 || raised != 0 {
				t.Errorf("launcher calls = %d, warnings = %d, want 0 and 0", len(l.upOpts), raised)
			}
			if strings.Contains(err.Error(), tc.baseDir) && tc.want == driving.ErrDevcontainerProjectFolderInvalid {
				t.Errorf("error %q names the path", err)
			}
		})
	}
}

// ADR-0028 decision 7: the engine sentinels keep their identity through the
// service, and a non-runtime failure reads no status file; the only question
// to the inspector is the query before the launcher (ADR-0030 decision 2).
func TestDevcontainerUp_EngineSentinelsPassThrough(t *testing.T) {
	t.Parallel()
	for _, want := range []error{driven.ErrDockerUnavailable, driven.ErrContainerEngineNotLocal, context.Canceled} {
		t.Run(want.Error(), func(t *testing.T) {
			t.Parallel()
			l := newLauncher()
			l.upErr = want
			insp := newFakeInspector(true)
			insp.setStatus("vfs", "failed")
			svc, _ := devcontainerSvc(t, plainConfig, yamlPodman, l, insp)
			_, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: devcontainerTestBaseDir})
			if !errors.Is(err, want) || errors.Is(err, driving.ErrDevcontainerStartScriptFailed) || errors.Is(err, driven.ErrDevcontainerRuntime) {
				t.Errorf("err = %v, want only %v", err, want)
			}
			if len(insp.readPaths) != 0 || len(insp.findFolders) != 1 {
				t.Errorf("the inspector read %d files and was asked %d times after %v, want 0 reads and the one query before the launcher",
					len(insp.readPaths), len(insp.findFolders), want)
			}
		})
	}
}

// ADR-0029 decision 1 (fitness function): after a runtime failure, a status
// file with `result=failed` makes the exit class 11, any other state keeps 12.
func TestDevcontainerUp_RuntimeFailureClassification(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		ubootYAML  string
		status     map[string]string
		notFound   bool
		wantScript bool
	}{
		{"vfs failed", yamlPodman, map[string]string{"vfs": "failed"}, false, true},
		{"daemon failed", yamlDocker, map[string]string{"daemon": "failed"}, false, true},
		{"egress failed", yamlEgressOnly, map[string]string{"egress": "failed"}, false, true},
		{"one of two failed", yamlDockerEgress, map[string]string{"daemon": "ok", "egress": "failed"}, false, true},
		{"ok", yamlPodman, map[string]string{"vfs": "ok"}, false, false},
		{"fallback is no abort", yamlPodman, map[string]string{"vfs": "fallback"}, false, false},
		{"no status file", yamlPodman, nil, false, false},
		{"no container", yamlPodman, map[string]string{"vfs": "failed"}, true, false},
		{"failed file of an unconfigured capability", yamlPodman, map[string]string{"daemon": "failed"}, false, false},
		{"no u-boot.yaml", "", map[string]string{"vfs": "failed"}, false, false},
		{"profile is not sandbox", sandboxYAMLHead + "  nope: 1\n", map[string]string{"vfs": "failed"}, false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := newLauncher()
			l.upErr = wrapRuntime("devcontainer up failed (exit status 1: boom)")
			insp := newFakeInspector(true)
			insp.found = !tc.notFound
			for file, result := range tc.status {
				insp.setStatus(file, result)
			}
			svc, _ := devcontainerSvc(t, plainConfig, tc.ubootYAML, l, insp)
			_, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: devcontainerTestBaseDir})
			if err == nil {
				t.Fatal("Up succeeded")
			}
			gotScript := errors.Is(err, driving.ErrDevcontainerStartScriptFailed)
			gotRuntime := errors.Is(err, driven.ErrDevcontainerRuntime)
			if gotScript != tc.wantScript || gotRuntime == tc.wantScript {
				t.Errorf("start script failed = %v, runtime = %v, want script %v (classes must not overlap): %v", gotScript, gotRuntime, tc.wantScript, err)
			}
			if !strings.Contains(err.Error(), "boom") {
				t.Errorf("error %q lost the CLI output", err)
			}
		})
	}
}

// wrapRuntime builds an error the way the launcher does: a text with the
// runtime sentinel at its end.
func wrapRuntime(text string) error {
	return &runtimeErr{text: text}
}

type runtimeErr struct{ text string }

func (e *runtimeErr) Error() string { return e.text + ": " + driven.ErrDevcontainerRuntime.Error() }
func (e *runtimeErr) Unwrap() error { return driven.ErrDevcontainerRuntime }

// LH-FA-DEV-007 after the start: a container whose start script reported
// `fallback` yields a degradation-fallback warning in the wording of `up`;
// `ok` and a missing `u-boot.yaml` are silent. `failed` is no warning: it ends
// the call (TestDevcontainerUp_FailedStatusAfterSuccess).
func TestDevcontainerUp_FallbackWarningsAfterStart(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		ubootYAML string
		file      string
		result    string
		code      string
		contains  string
		want      int
	}{
		{"vfs", yamlPodman, "vfs", "fallback", "LH-FA-DEV-007", "vfs storage", 1},
		{"daemon", yamlDocker, "daemon", "fallback", "LH-FA-DEV-007", "no Docker daemon", 1},
		{"egress", yamlEgressOnly, "egress", "fallback", "LH-FA-DEV-008", "egress restriction is not active", 1},
		{"ok is silent", yamlPodman, "vfs", "ok", "", "", 0},
		{"no u-boot.yaml: no disclosure (not an error)", "", "vfs", "fallback", "", "", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			insp := newFakeInspector(true)
			insp.setStatus(tc.file, tc.result)
			svc, _ := devcontainerSvc(t, plainConfig, tc.ubootYAML, newLauncher(), insp)
			resp, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: devcontainerTestBaseDir})
			if err != nil {
				t.Fatal(err)
			}
			if len(resp.Warnings) != tc.want {
				t.Fatalf("warnings = %+v, want %d", resp.Warnings, tc.want)
			}
			if tc.want == 0 {
				return
			}
			w := resp.Warnings[0]
			if w.Code != tc.code || w.Level != "warn" || !w.DegradationFallback || w.SecurityRelaxation || !strings.Contains(w.Message, tc.contains) {
				t.Errorf("warning = %+v", w)
			}
			// The wording of `up`: the container was just started, nothing is
			// regenerated, so neither "existing" nor the rebuild instruction.
			if !strings.Contains(w.Message, "state of the devcontainer container at its last start ("+fallbackAt+")") ||
				strings.Contains(w.Message, "existing") || strings.Contains(w.Message, "rebuilt") {
				t.Errorf("message %q is not the wording of `up`", w.Message)
			}
		})
	}
}

// ADR-0030 decision 1 (fitness function), LH-FA-DEV-010: after a successful
// run of the devcontainer CLI, `result=failed` in the status file of a
// configured capability ends the call with the start-script sentinel, in the
// class of exit 11 and not of exit 12. The message carries the sentence of the
// capability, the time of the report, that the container runs and this call
// did not run the script, and `--recreate`. Every other state stays a success.
func TestDevcontainerUp_FailedStatusAfterSuccess(t *testing.T) {
	t.Parallel()
	const tail = "the container is running and this call did not run the start script; `u-boot devcontainer up --recreate` creates the container anew"
	reported := "(reported at " + fallbackAt + ")"
	vfsPath := fallbackStatusDir + "/vfs"
	cases := []struct {
		name      string
		ubootYAML string
		status    map[string]string
		notFound  bool
		readErrs  map[string]error
		wantErr   bool
		contains  []string
		absent    []string
	}{
		{"vfs failed", yamlPodman, map[string]string{"vfs": "failed"}, false, nil, true,
			[]string{"/dev/fuse was not available in the container (onUnavailable: fail) " + reported + "; " + tail}, nil},
		{"daemon failed", yamlDocker, map[string]string{"daemon": "failed"}, false, nil, true,
			[]string{"the nested Docker daemon did not become ready (onUnavailable: fail) " + reported + "; " + tail}, nil},
		{"egress failed", yamlEgressOnly, map[string]string{"egress": "failed"}, false, nil, true,
			[]string{"the NET_ADMIN capability was not granted (onUnavailable: fail)", reported + "; " + tail}, nil},
		{"one of two failed", yamlDockerEgress, map[string]string{"daemon": "ok", "egress": "failed"}, false, nil, true,
			[]string{"NET_ADMIN"}, []string{"nested Docker daemon"}},
		{"failed wins over a fallback of another capability", yamlDockerEgress, map[string]string{"daemon": "fallback", "egress": "failed"}, false, nil, true,
			[]string{"NET_ADMIN"}, []string{"the container has no Docker daemon"}},
		{"the first failed capability is named", yamlPodmanEgress, map[string]string{"vfs": "failed", "egress": "failed"}, false, nil, true,
			[]string{"/dev/fuse"}, []string{"NET_ADMIN"}},
		{"ok", yamlPodman, map[string]string{"vfs": "ok"}, false, nil, false, nil, nil},
		{"fallback is a warning, no abort", yamlPodman, map[string]string{"vfs": "fallback"}, false, nil, false, nil, nil},
		{"no status file", yamlPodman, nil, false, nil, false, nil, nil},
		{"no container", yamlPodman, map[string]string{"vfs": "failed"}, true, nil, false, nil, nil},
		{"failed file of an unconfigured capability", yamlPodman, map[string]string{"daemon": "failed"}, false, nil, false, nil, nil},
		{"no u-boot.yaml", "", map[string]string{"vfs": "failed"}, false, nil, false, nil, nil},
		{"profile is not sandbox", sandboxYAMLHead + "  nope: 1\n", map[string]string{"vfs": "failed"}, false, nil, false, nil, nil},
		{"invalid status document", yamlPodman, map[string]string{"vfs": "failed"}, false,
			map[string]error{vfsPath: driven.ErrContainerFileInvalid}, false, nil, nil},
		{"status file cannot be read", yamlPodman, map[string]string{"vfs": "failed"}, false,
			map[string]error{vfsPath: errors.New("engine hiccup")}, false, nil, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := newLauncher()
			insp := newFakeInspector(true)
			insp.found = !tc.notFound
			for file, result := range tc.status {
				insp.setStatus(file, result)
			}
			for path, err := range tc.readErrs {
				insp.readErrs[path] = err
			}
			svc, _ := devcontainerSvc(t, plainConfig, tc.ubootYAML, l, insp)
			resp, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: devcontainerTestBaseDir})
			if len(l.upOpts) != 1 {
				t.Fatalf("launcher Up calls = %d, want 1", len(l.upOpts))
			}
			if !tc.wantErr {
				if err != nil || resp.ContainerID != fallbackContainerID {
					t.Fatalf("resp = %+v, err = %v, want a success", resp, err)
				}
				return
			}
			if !errors.Is(err, driving.ErrDevcontainerStartScriptFailed) || errors.Is(err, driven.ErrDevcontainerRuntime) {
				t.Fatalf("err = %v, want the start-script sentinel and not the runtime class", err)
			}
			if resp.ContainerID != "" || len(resp.Warnings) != 0 || resp.AlreadyRunning {
				t.Errorf("resp = %+v, want the zero response", resp)
			}
			for _, want := range tc.contains {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("message %q lacks %q", err, want)
				}
			}
			for _, unwanted := range tc.absent {
				if strings.Contains(err.Error(), unwanted) {
					t.Errorf("message %q holds %q", err, unwanted)
				}
			}
		})
	}
}

// ADR-0030 decision 2 (fitness function), LH-FA-DEV-010: the response says the
// container was already running when the inspector found it running before
// the launcher and the CLI reports the same ID. The query comes before the
// launcher, names the project folder and is skipped with Recreate; a stopped
// or absent container, another ID and a failed query give false, and the failed
// query is no error.
func TestDevcontainerUp_AlreadyRunning(t *testing.T) {
	t.Parallel()
	running := driven.DevcontainerContainer{ID: fallbackContainerID, Running: true}
	cases := []struct {
		name        string
		recreate    bool
		answers     []findAnswer
		want        bool
		wantQueries int
	}{
		{"running before, same ID", false, nil, true, 1},
		{"stopped before", false, []findAnswer{{container: driven.DevcontainerContainer{ID: fallbackContainerID}, found: true}}, false, 1},
		{"no container before", false, []findAnswer{{}}, false, 1},
		{"another ID before", false, []findAnswer{{container: driven.DevcontainerContainer{ID: strings.Repeat("f", 64), Running: true}, found: true}}, false, 1},
		{"the query fails", false, []findAnswer{{container: running, found: true, err: errors.New("engine hiccup")}}, false, 1},
		{"the engine is not local", false, []findAnswer{{err: driven.ErrContainerEngineNotLocal}}, false, 1},
		{"recreate skips the query", true, nil, false, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := newLauncher()
			insp := newFakeInspector(true)
			insp.findAnswers = tc.answers
			queriesBeforeLauncher := -1
			l.onUp = func() { queriesBeforeLauncher = len(insp.findFolders) }
			svc, _ := devcontainerSvc(t, plainConfig, "", l, insp)
			resp, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: devcontainerTestBaseDir, Recreate: tc.recreate})
			if err != nil {
				t.Fatalf("Up: %v", err)
			}
			if resp.AlreadyRunning != tc.want || resp.ContainerID != fallbackContainerID {
				t.Errorf("resp = %+v, want AlreadyRunning %v", resp, tc.want)
			}
			if queriesBeforeLauncher != tc.wantQueries || len(insp.findFolders) != tc.wantQueries {
				t.Errorf("queries before the launcher = %d, in all = %d, want %d and %d (without u-boot.yaml nothing asks afterwards)",
					queriesBeforeLauncher, len(insp.findFolders), tc.wantQueries, tc.wantQueries)
			}
			for _, folder := range insp.findFolders {
				if folder != devcontainerTestBaseDir {
					t.Errorf("the query names %q, want the project folder", folder)
				}
			}
		})
	}
}

// The query before the launcher and the reading of the status files are two
// questions: a failed first answer does not make the status files unreadable,
// and the container found afterwards is the one whose files are read.
func TestDevcontainerUp_QueryBeforeTheLauncherDoesNotSpoilTheStatusReading(t *testing.T) {
	t.Parallel()
	insp := newFakeInspector(true)
	insp.findAnswers = []findAnswer{{err: errors.New("engine hiccup")}}
	insp.setStatus("vfs", "failed")
	svc, _ := devcontainerSvc(t, plainConfig, yamlPodman, newLauncher(), insp)
	_, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: devcontainerTestBaseDir})
	if !errors.Is(err, driving.ErrDevcontainerStartScriptFailed) {
		t.Fatalf("err = %v, want the start-script sentinel", err)
	}
	if len(insp.findFolders) != 2 {
		t.Errorf("queries = %d, want the one before the launcher and the one of the status reading", len(insp.findFolders))
	}
}

// A missing port is an error, never "nothing to do" with exit 0.
func TestDevcontainer_UnwiredPortsAreErrors(t *testing.T) {
	t.Parallel()
	t.Run("up without launcher", func(t *testing.T) {
		t.Parallel()
		svc, _ := devcontainerSvc(t, plainConfig, "", nil, newFakeInspector(true))
		if _, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: devcontainerTestBaseDir}); err == nil || !strings.Contains(err.Error(), "launcher not wired") {
			t.Errorf("err = %v", err)
		}
	})
	t.Run("up without inspector", func(t *testing.T) {
		t.Parallel()
		svc, _ := devcontainerSvc(t, plainConfig, "", newLauncher(), nil)
		if _, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: devcontainerTestBaseDir}); err == nil || !strings.Contains(err.Error(), "inspector not wired") {
			t.Errorf("err = %v", err)
		}
	})
	t.Run("down without inspector", func(t *testing.T) {
		t.Parallel()
		svc, _ := devcontainerSvc(t, "", "", newLauncher(), nil)
		if _, err := svc.Down(context.Background(), driving.DevcontainerDownRequest{BaseDir: devcontainerTestBaseDir}); err == nil || !strings.Contains(err.Error(), "inspector not wired") {
			t.Errorf("err = %v", err)
		}
	})
	t.Run("down without launcher", func(t *testing.T) {
		t.Parallel()
		svc, _ := devcontainerSvc(t, "", "", nil, newFakeInspector(true))
		if _, err := svc.Down(context.Background(), driving.DevcontainerDownRequest{BaseDir: devcontainerTestBaseDir}); err == nil || !strings.Contains(err.Error(), "launcher not wired") {
			t.Errorf("err = %v", err)
		}
	})
	t.Run("empty base directory", func(t *testing.T) {
		t.Parallel()
		svc, _ := devcontainerSvc(t, "", "", newLauncher(), newFakeInspector(true))
		if _, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{}); err == nil {
			t.Error("Up without BaseDir succeeded")
		}
		if _, err := svc.Down(context.Background(), driving.DevcontainerDownRequest{}); err == nil {
			t.Error("Down without BaseDir succeeded")
		}
	})
}

// ADR-0028 decision 3: `down` stops a running container, and "nothing to do"
// (no container, already stopped) is a success without a Stop call.
func TestDevcontainerDown(t *testing.T) {
	t.Parallel()
	t.Run("stops a running container", func(t *testing.T) {
		t.Parallel()
		l := newLauncher()
		svc, _ := devcontainerSvc(t, "", "", l, newFakeInspector(true))
		resp, err := svc.Down(context.Background(), driving.DevcontainerDownRequest{BaseDir: devcontainerTestBaseDir})
		if err != nil || !resp.Stopped || resp.ContainerID != fallbackContainerID {
			t.Fatalf("resp = %+v, err = %v", resp, err)
		}
		if !reflect.DeepEqual(l.stopIDs, []string{fallbackContainerID}) {
			t.Errorf("Stop calls = %v", l.stopIDs)
		}
	})
	t.Run("a stopped container is not stopped again", func(t *testing.T) {
		t.Parallel()
		l := newLauncher()
		svc, _ := devcontainerSvc(t, "", "", l, newFakeInspector(false))
		resp, err := svc.Down(context.Background(), driving.DevcontainerDownRequest{BaseDir: devcontainerTestBaseDir})
		if err != nil || resp.Stopped || resp.ContainerID != fallbackContainerID || len(l.stopIDs) != 0 {
			t.Errorf("resp = %+v, err = %v, Stop calls = %v", resp, err, l.stopIDs)
		}
	})
	t.Run("no container", func(t *testing.T) {
		t.Parallel()
		l := newLauncher()
		insp := newFakeInspector(true)
		insp.found = false
		svc, _ := devcontainerSvc(t, "", "", l, insp)
		resp, err := svc.Down(context.Background(), driving.DevcontainerDownRequest{BaseDir: devcontainerTestBaseDir})
		if err != nil || resp != (driving.DevcontainerDownResponse{}) || len(l.stopIDs) != 0 {
			t.Errorf("resp = %+v, err = %v, Stop calls = %v", resp, err, l.stopIDs)
		}
		if !reflect.DeepEqual(insp.findFolders, []string{devcontainerTestBaseDir}) {
			t.Errorf("FindContainer folders = %v", insp.findFolders)
		}
	})
	t.Run("the container is found by the project folder and needs no configuration", func(t *testing.T) {
		t.Parallel()
		svc, _ := devcontainerSvc(t, "", "", newLauncher(), newFakeInspector(true))
		if _, err := svc.Down(context.Background(), driving.DevcontainerDownRequest{BaseDir: devcontainerTestBaseDir}); err != nil {
			t.Errorf("Down without devcontainer.json: %v", err)
		}
	})
	t.Run("stop failure is the runtime class", func(t *testing.T) {
		t.Parallel()
		l := newLauncher()
		l.stopErr = wrapRuntime("docker stop failed")
		svc, _ := devcontainerSvc(t, "", "", l, newFakeInspector(true))
		if _, err := svc.Down(context.Background(), driving.DevcontainerDownRequest{BaseDir: devcontainerTestBaseDir}); !errors.Is(err, driven.ErrDevcontainerRuntime) {
			t.Errorf("err = %v", err)
		}
	})
	t.Run("stop precondition failure keeps its class", func(t *testing.T) {
		t.Parallel()
		l := newLauncher()
		l.stopErr = driven.ErrDockerUnavailable
		svc, _ := devcontainerSvc(t, "", "", l, newFakeInspector(true))
		if _, err := svc.Down(context.Background(), driving.DevcontainerDownRequest{BaseDir: devcontainerTestBaseDir}); !errors.Is(err, driven.ErrDockerUnavailable) {
			t.Errorf("err = %v", err)
		}
	})
}

// A remote engine keeps its class; any other failure to ask the engine is an
// environment failure (exit 11), never "no container" with exit 0.
func TestDevcontainerDown_EngineFailures(t *testing.T) {
	t.Parallel()
	t.Run("remote engine", func(t *testing.T) {
		t.Parallel()
		insp := newFakeInspector(true)
		insp.findErr = driven.ErrContainerEngineNotLocal
		svc, _ := devcontainerSvc(t, "", "", newLauncher(), insp)
		_, err := svc.Down(context.Background(), driving.DevcontainerDownRequest{BaseDir: devcontainerTestBaseDir})
		if !errors.Is(err, driven.ErrContainerEngineNotLocal) {
			t.Errorf("err = %v", err)
		}
	})
	t.Run("engine cannot be asked", func(t *testing.T) {
		t.Parallel()
		l := newLauncher()
		insp := newFakeInspector(true)
		insp.findErr = errors.New("docker context inspect: exec: \"docker\": executable file not found in $PATH")
		svc, _ := devcontainerSvc(t, "", "", l, insp)
		_, err := svc.Down(context.Background(), driving.DevcontainerDownRequest{BaseDir: devcontainerTestBaseDir})
		if !errors.Is(err, driven.ErrDockerUnavailable) || !strings.Contains(err.Error(), "executable file not found") || len(l.stopIDs) != 0 {
			t.Errorf("err = %v, Stop calls = %v", err, l.stopIDs)
		}
	})
	t.Run("cancelled", func(t *testing.T) {
		t.Parallel()
		insp := newFakeInspector(true)
		insp.findErr = errors.New("signal: killed")
		svc, _ := devcontainerSvc(t, "", "", newLauncher(), insp)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := svc.Down(ctx, driving.DevcontainerDownRequest{BaseDir: devcontainerTestBaseDir})
		if !errors.Is(err, context.Canceled) || errors.Is(err, driven.ErrDockerUnavailable) {
			t.Errorf("err = %v", err)
		}
	})
}
