package application_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// sessionConfig is a configuration in the shape the generator writes for the
// sandbox profile (`remoteEnv` with the entry of the Git token and the entry
// of an agent) plus a reference of another form in `mounts`.
const sessionConfig = `{
  "name": "demo",
  "mounts": ["source=${localEnv:HOME}/.cache,target=/cache,type=bind"],
  "remoteEnv": {
    "GIT_TOKEN": "${localEnv:GIT_TOKEN}",
    "ANTHROPIC_API_KEY": "${localEnv:ANTHROPIC_API_KEY}"
  },
  "remoteUser": "vscode"
}`

// ADR-0031 decision 1 (Umgebung): the session gets the entries
// `"<NAME>": "${localEnv:<NAME>}"` of the configuration and nothing else.
func TestSessionEnvNames(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		text string
		want []string
	}{
		{"none", `{"image": "debian"}`, nil},
		{"the shape the generator writes", sessionConfig, []string{"ANTHROPIC_API_KEY", "GIT_TOKEN"}},
		{"space around the colon", `{"A" :   "${localEnv:A}"}`, []string{"A"}},
		{"sorted and distinct", `{"B": "${localEnv:B}", "A": "${localEnv:A}", "again": {"B": "${localEnv:B}"}}`, []string{"A", "B"}},
		{"another name in the reference", `{"A": "${localEnv:B}"}`, nil},
		{"a default", `{"A": "${localEnv:A:fallback}"}`, nil},
		{"a reference inside a longer value", `{"A": "pre-${localEnv:A}"}`, nil},
		{"a reference in mounts", `{"mounts": ["source=${localEnv:HOME}/x,target=/x,type=bind"]}`, nil},
		{"a key that is no plain name", `{"A-B": "${localEnv:A-B}"}`, nil},
		{"an entry in a line comment", "// \"A\": \"${localEnv:A}\"\n{}", nil},
		{"an entry in a block comment", `{/* "A": "${localEnv:A}" */}`, nil},
		{"an entry after a line comment", "{ // note\n\"A\": \"${localEnv:A}\"}", []string{"A"}},
		{"an entry after a block comment", `{/* note */ "A": "${localEnv:A}"}`, []string{"A"}},
		{"// in a string starts no comment", `{"url": "http://example.test/x", "A": "${localEnv:A}"}`, []string{"A"}},
		{"/* in a string starts no comment", `{"glob": "a/*b", "A": "${localEnv:A}", "z": "*/"}`, []string{"A"}},
		{"an escaped quote does not end the string", `{"s": "a \" // b", "A": "${localEnv:A}"}`, []string{"A"}},
		{"multi-byte text before the entry", `{"s": "ä // ö", "A": "${localEnv:A}"}`, []string{"A"}},
		{"an unterminated block comment hides the rest", `{"A": "${localEnv:A}"} /* "B": "${localEnv:B}"`, []string{"A"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := application.SessionEnvNamesForTest(tc.text); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("names = %q, want %q", got, tc.want)
			}
		})
	}
}

// The comment stripper keeps what is not a comment byte for byte and replaces
// a comment by its line break or one space.
func TestWithoutComments(t *testing.T) {
	t.Parallel()
	cases := []struct{ name, text, want string }{
		{"no comment", `{"a": 1}`, `{"a": 1}`},
		{"line comment keeps its line break", "a // b\nc", "a \nc"},
		{"line comment at the end of the text", "a // b", "a "},
		{"block comment becomes a space", "a/* b */c", "a c"},
		{"block comment over lines", "a/* b\nc */d", "a d"},
		{"unterminated block comment", "a/* b", "a "},
		{"comment marker in a string", `"a//b" "c/*d*/"`, `"a//b" "c/*d*/"`},
		{"escaped quote in a string", `"a\"//b" // c`, `"a\"//b" `},
		{"unterminated string", `"a // b`, `"a // b`},
		{"a lone slash", "a / b", "a / b"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := application.WithoutCommentsForTest(tc.text); got != tc.want {
				t.Errorf("without comments = %q, want %q", got, tc.want)
			}
		})
	}
}

// shellFixture is the service on a project with the given configuration and
// u-boot.yaml, a launcher whose start reports user vscode in /workspaces/demo,
// and an inspector for a running container.
func shellFixture(t *testing.T, config, ubootYAML string) (*application.DevcontainerStartService, *fakeLauncher, *fakeInspector) {
	t.Helper()
	l := newLauncher()
	insp := newFakeInspector(true)
	svc, _ := devcontainerSvc(t, config, ubootYAML, l, insp)
	return svc, l, insp
}

func shellRequest() driving.DevcontainerShellRequest {
	return driving.DevcontainerShellRequest{BaseDir: devcontainerTestBaseDir}
}

// ADR-0031 decision 1 (Sitzungsphase): the session gets the container, the
// user and the workspace of the result of the start, the entries of the
// configuration as its variables, and the streams and the terminal flag of the
// request unchanged. The helper container of the start still gets every
// `${localEnv:…}` reference by name (the path in `mounts` too).
func TestDevcontainerShell_OpensTheSessionWithTheResultOfTheStart(t *testing.T) {
	t.Parallel()
	svc, l, _ := shellFixture(t, sessionConfig, "")
	l.upResult = driven.DevcontainerUpResult{ContainerID: strings.Repeat("ab", 32), RemoteUser: "node", RemoteWorkspaceFolder: "/workspaces/other"}
	stdin, stdout, stderr := strings.NewReader("ls\n"), &bytes.Buffer{}, &bytes.Buffer{}
	req := shellRequest()
	req.Stdin, req.Stdout, req.Stderr, req.Terminal = stdin, stdout, stderr, true

	if err := svc.Shell(context.Background(), req); err != nil {
		t.Fatalf("Shell: %v", err)
	}
	if len(l.shellOpts) != 1 {
		t.Fatalf("Shell calls = %d, want 1", len(l.shellOpts))
	}
	got := l.shellOpts[0]
	if got.ContainerID != strings.Repeat("ab", 32) || got.User != "node" || got.Workspace != "/workspaces/other" {
		t.Errorf("session = container %q user %q workspace %q, want those of the start", got.ContainerID, got.User, got.Workspace)
	}
	if want := []string{"ANTHROPIC_API_KEY", "GIT_TOKEN"}; !reflect.DeepEqual(got.EnvNames, want) {
		t.Errorf("session variables = %q, want %q", got.EnvNames, want)
	}
	if got.Stdin != io.Reader(stdin) || got.Stdout != io.Writer(stdout) || got.Stderr != io.Writer(stderr) || !got.Terminal {
		t.Errorf("session streams or terminal flag differ from the request: %+v", got)
	}
	if len(l.upOpts) != 1 || !reflect.DeepEqual(l.upOpts[0].EnvNames, []string{"ANTHROPIC_API_KEY", "GIT_TOKEN", "HOME"}) {
		t.Errorf("start options = %+v, want the helper to get every reference by name", l.upOpts)
	}
}

// The terminal flag is passed as it is, in both values.
func TestDevcontainerShell_PassesTheTerminalFlagBothWays(t *testing.T) {
	t.Parallel()
	for _, terminal := range []bool{true, false} {
		svc, l, _ := shellFixture(t, plainConfig, "")
		req := shellRequest()
		req.Terminal = terminal
		if err := svc.Shell(context.Background(), req); err != nil {
			t.Fatal(err)
		}
		if len(l.shellOpts) != 1 || l.shellOpts[0].Terminal != terminal {
			t.Errorf("request terminal %v: session options %+v", terminal, l.shellOpts)
		}
	}
}

// ADR-0031 decision 1 (Startphase): the start is the one of `up` without
// `--recreate`: the project folder, the progress sink and no recreate. `shell`
// does not ask the inspector whether the container already ran, because it
// reports nothing about it.
func TestDevcontainerShell_StartIsUpWithoutRecreate(t *testing.T) {
	t.Parallel()
	svc, l, insp := shellFixture(t, plainConfig, "")
	progress := &bytes.Buffer{}
	req := shellRequest()
	req.ProgressSink = progress
	if err := svc.Shell(context.Background(), req); err != nil {
		t.Fatal(err)
	}
	if len(l.upOpts) != 1 {
		t.Fatalf("launcher Up calls = %d, want 1", len(l.upOpts))
	}
	if got := l.upOpts[0]; got.ProjectFolder != devcontainerTestBaseDir || got.Recreate || got.ReadOnlyProject || got.Progress != io.Writer(progress) {
		t.Errorf("start options = %+v", got)
	}
	if len(insp.findFolders) != 0 {
		t.Errorf("the inspector was asked %d times, want none: `shell` reports no \"already running\"", len(insp.findFolders))
	}
}

// ADR-0031 decision 2: the disclosure of the socket names the duration of the
// start, `up` keeps the duration of the call, and the disclosure comes before
// the start of the helper container and carries the flags of a security
// relaxation.
func TestDevcontainerShell_SocketDisclosureNamesTheStart(t *testing.T) {
	t.Parallel()
	svc, l, _ := shellFixture(t, plainConfig, "")
	var shell []driving.WarningEntry
	l.onUp = func() {
		if len(shell) != 1 {
			t.Errorf("warnings before the helper container started = %d, want 1", len(shell))
		}
	}
	req := shellRequest()
	req.OnWarning = func(w driving.WarningEntry) { shell = append(shell, w) }
	if err := svc.Shell(context.Background(), req); err != nil {
		t.Fatal(err)
	}
	if len(shell) != 1 {
		t.Fatalf("warnings = %+v, want exactly the disclosure", shell)
	}
	w := shell[0]
	if w.Code != "LH-FA-DEV-010" || w.Level != "warn" || !w.SecurityRelaxation || w.DegradationFallback {
		t.Errorf("disclosure = %+v", w)
	}
	if !strings.Contains(w.Message, "root-equivalent access to the engine for the duration of the start)") ||
		strings.Contains(w.Message, "duration of the call") {
		t.Errorf("disclosure of shell %q does not name the duration of the start", w.Message)
	}

	// The wording of `up` is unchanged, and the rest of both texts is the same.
	upSvc, _, _ := shellFixture(t, plainConfig, "")
	var up []driving.WarningEntry
	if _, err := upSvc.Up(context.Background(), driving.DevcontainerUpRequest{
		BaseDir: devcontainerTestBaseDir, OnWarning: func(w driving.WarningEntry) { up = append(up, w) },
	}); err != nil || len(up) != 1 {
		t.Fatalf("Up: err=%v warnings=%+v", err, up)
	}
	if !strings.Contains(up[0].Message, "for the duration of the call)") || strings.Contains(up[0].Message, "duration of the start") {
		t.Errorf("disclosure of up %q lost the duration of the call", up[0].Message)
	}
	if strings.Replace(up[0].Message, "the call)", "the start)", 1) != w.Message {
		t.Errorf("the two disclosures differ in more than the duration:\nup:    %s\nshell: %s", up[0].Message, w.Message)
	}
}

// ADR-0031 decision 1 (Fallback-Warnungen): a degradation fallback the
// container reports reaches OnWarning before the session begins, in the wording
// of `up`, after the disclosure; `ok` raises nothing.
func TestDevcontainerShell_FallbackWarningsComeBeforeTheSession(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		result string
		want   int
	}{{"fallback", 1}, {"ok", 0}} {
		t.Run(tc.result, func(t *testing.T) {
			t.Parallel()
			svc, l, insp := shellFixture(t, plainConfig, yamlPodman)
			insp.setStatus("vfs", tc.result)
			var raised []driving.WarningEntry
			atSession := -1
			l.onShell = func() { atSession = len(raised) }
			req := shellRequest()
			req.OnWarning = func(w driving.WarningEntry) { raised = append(raised, w) }
			if err := svc.Shell(context.Background(), req); err != nil {
				t.Fatal(err)
			}
			if len(raised) != 1+tc.want || atSession != len(raised) {
				t.Fatalf("warnings = %+v, raised before the session = %d, want the disclosure and %d fallback(s) before the session", raised, atSession, tc.want)
			}
			if !raised[0].SecurityRelaxation {
				t.Errorf("first warning %+v is not the disclosure", raised[0])
			}
			if tc.want == 1 {
				w := raised[1]
				if w.Code != "LH-FA-DEV-007" || !w.DegradationFallback || !strings.Contains(w.Message, "vfs storage") ||
					!strings.Contains(w.Message, "state of the devcontainer container at its last start ("+fallbackAt+")") {
					t.Errorf("fallback warning = %+v", w)
				}
			}
		})
	}
}

// ADR-0031 decision 1 (Startphase, Regel nach dem Zustand), ADR-0030 decision
// 1: a `result=failed` of a configured capability ends the call with the
// start-script sentinel (class 11, not the runtime class) and opens no
// session, after a successful run of the CLI and after a failed one; the
// message names `shell`.
func TestDevcontainerShell_FailedStatusOpensNoSession(t *testing.T) {
	t.Parallel()
	for _, launchFailed := range []bool{false, true} {
		name := "after a successful run of the CLI"
		if launchFailed {
			name = "after a failed run of the CLI"
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, l, insp := shellFixture(t, plainConfig, yamlPodman)
			insp.setStatus("vfs", "failed")
			if launchFailed {
				l.upErr = wrapRuntime("devcontainer up failed (exit status 1: boom)")
			}
			err := svc.Shell(context.Background(), shellRequest())
			if !errors.Is(err, driving.ErrDevcontainerStartScriptFailed) || errors.Is(err, driven.ErrDevcontainerRuntime) {
				t.Fatalf("err = %v, want the start-script sentinel and not the runtime class", err)
			}
			if len(l.shellOpts) != 0 {
				t.Errorf("a session was opened (%d Shell calls)", len(l.shellOpts))
			}
			if !strings.HasPrefix(err.Error(), "devcontainer shell: ") || strings.Contains(err.Error(), "devcontainer up:") {
				t.Errorf("message %q does not name the command", err)
			}
			// After a failed run the message keeps the output of the CLI; after a
			// successful one it names the way out.
			if !launchFailed && !strings.Contains(err.Error(), "u-boot devcontainer up --recreate") {
				t.Errorf("message %q does not name the way out", err)
			}
			if launchFailed && !strings.Contains(err.Error(), "boom") {
				t.Errorf("message %q lost the output of the CLI", err)
			}
		})
	}
}

// ADR-0030 limit (ii), which holds for `shell` as for `up`: without a
// u-boot.yaml, without the sandbox profile or without the capability no status
// file is read, so the session opens in a container whose start script failed.
func TestDevcontainerShell_FailedStatusNobodyReadsDoesNotStopTheSession(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, ubootYAML string }{
		{"no u-boot.yaml", ""},
		{"capability not configured", yamlDocker},
		{"profile is not sandbox", sandboxYAMLHead + "  nope: 1\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc, l, insp := shellFixture(t, plainConfig, tc.ubootYAML)
			insp.setStatus("vfs", "failed")
			if err := svc.Shell(context.Background(), shellRequest()); err != nil {
				t.Fatalf("Shell: %v", err)
			}
			if len(l.shellOpts) != 1 {
				t.Errorf("Shell calls = %d, want the session to open", len(l.shellOpts))
			}
		})
	}
}

// The engine sentinels and an aborted start keep their identity through
// `shell`, and no session opens after any failure of the start.
func TestDevcontainerShell_StartFailuresOpenNoSession(t *testing.T) {
	t.Parallel()
	for _, want := range []error{driven.ErrDockerUnavailable, driven.ErrContainerEngineNotLocal, context.Canceled, driven.ErrDevcontainerRuntime} {
		t.Run(want.Error(), func(t *testing.T) {
			t.Parallel()
			svc, l, _ := shellFixture(t, plainConfig, "")
			l.upErr = want
			err := svc.Shell(context.Background(), shellRequest())
			if !errors.Is(err, want) || errors.Is(err, driving.ErrDevcontainerStartScriptFailed) {
				t.Errorf("err = %v, want %v", err, want)
			}
			if len(l.shellOpts) != 0 {
				t.Errorf("a session was opened after the failed start")
			}
		})
	}
}

// Validation fails before the disclosure and before any launcher call, with
// the name of the command in the message; the project path stays out of the
// message of the folder error.
func TestDevcontainerShell_ValidationFailuresRaiseNoDisclosureAndStartNothing(t *testing.T) {
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
		{"name that is no plain name", devcontainerTestBaseDir, `"${localEnv:A-B}"`, nil, driving.ErrDevcontainerEnvNameInvalid},
		{"probe fails", devcontainerTestBaseDir, plainConfig, func(fs *fakeFS) {
			fs.failExistsOn = devcontainerTestBaseDir + "/.devcontainer/devcontainer.json"
			fs.failExistsErr = errors.New("disk gone")
		}, driving.ErrDevcontainerFileSystem},
		{"read fails", devcontainerTestBaseDir, plainConfig, func(fs *fakeFS) {
			fs.failReadOn = devcontainerTestBaseDir + "/.devcontainer/devcontainer.json"
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
			err := svc.Shell(context.Background(), driving.DevcontainerShellRequest{BaseDir: tc.baseDir, OnWarning: func(driving.WarningEntry) { raised++ }})
			if !errors.Is(err, tc.want) {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
			if len(l.upOpts) != 0 || len(l.shellOpts) != 0 || raised != 0 {
				t.Errorf("Up calls = %d, Shell calls = %d, warnings = %d, want 0, 0 and 0", len(l.upOpts), len(l.shellOpts), raised)
			}
			if !strings.HasPrefix(err.Error(), "devcontainer shell: ") {
				t.Errorf("message %q does not name the command", err)
			}
		})
	}
}

// The start phase is shared, the name of the command in its texts is not: the
// errors of `up` still begin with "devcontainer up: " (and none with "shell").
func TestDevcontainerUp_ErrorsNameUp(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		baseDir string
		config  string
		yaml    string
		status  bool
	}{
		{"no configuration", devcontainerTestBaseDir, "", "", false},
		{"folder", "/tmp/a,b", plainConfig, "", false},
		{"variable name", devcontainerTestBaseDir, `"${localEnv:A-B}"`, "", false},
		{"failed status after a start", devcontainerTestBaseDir, plainConfig, yamlPodman, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := newLauncher()
			insp := newFakeInspector(true)
			if tc.status {
				insp.setStatus("vfs", "failed")
			}
			svc, _ := devcontainerSvc(t, tc.config, tc.yaml, l, insp)
			_, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: tc.baseDir})
			if err == nil || !strings.HasPrefix(err.Error(), "devcontainer up: ") {
				t.Errorf("err = %v, want a message that begins with \"devcontainer up: \"", err)
			}
		})
	}
}

// The error of the session passes through with its identity: a failed probe
// is the runtime class, which the CLI maps to exit code 12.
func TestDevcontainerShell_SessionErrorPassesThrough(t *testing.T) {
	t.Parallel()
	svc, l, _ := shellFixture(t, plainConfig, "")
	l.shellErr = wrapRuntime("devcontainer shell: no shell runs as user \"vscode\"")
	err := svc.Shell(context.Background(), shellRequest())
	if !errors.Is(err, driven.ErrDevcontainerRuntime) || errors.Is(err, driving.ErrDevcontainerStartScriptFailed) {
		t.Errorf("err = %v, want the runtime class", err)
	}
}

// A missing port is an error, never a session that did not open with exit 0.
func TestDevcontainerShell_UnwiredPortsAreErrors(t *testing.T) {
	t.Parallel()
	t.Run("without launcher", func(t *testing.T) {
		t.Parallel()
		svc, _ := devcontainerSvc(t, plainConfig, "", nil, newFakeInspector(true))
		if err := svc.Shell(context.Background(), shellRequest()); err == nil || !strings.Contains(err.Error(), "devcontainer shell: launcher not wired") {
			t.Errorf("err = %v", err)
		}
	})
	t.Run("without inspector", func(t *testing.T) {
		t.Parallel()
		svc, _ := devcontainerSvc(t, plainConfig, "", newLauncher(), nil)
		if err := svc.Shell(context.Background(), shellRequest()); err == nil || !strings.Contains(err.Error(), "devcontainer shell: inspector not wired") {
			t.Errorf("err = %v", err)
		}
	})
	t.Run("without base directory", func(t *testing.T) {
		t.Parallel()
		svc, _ := devcontainerSvc(t, plainConfig, "", newLauncher(), newFakeInspector(true))
		if err := svc.Shell(context.Background(), driving.DevcontainerShellRequest{}); err == nil || !strings.Contains(err.Error(), "devcontainer shell: BaseDir is required") {
			t.Errorf("err = %v", err)
		}
	})
}
