package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli/jsontestutil"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// slice-v1-cli-json-envelope-consolidation T2: add/init/generate now
// route their positional-args validation through the shared
// jsonArgsValidator (T1), so a wrong-arg invocation under --json
// emits the spec envelope on stdout (§1841) instead of a bare Cobra
// stderr message — and --dry-run/--diff selects the Voll-Schema
// (§1842). These pins lock the contract per command.
//
// Matrix (User-Review): add/generate carry ExactArgs(1) → both NoArg
// AND TooMany error; init carries MaximumNArgs(1) → only TooMany
// errors (0 args is a valid success, separately pinned).

func assertArgsErrorEnvelope(t *testing.T, app *cli.App, args []string, command string, full bool) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), args, &stdout, &stderr)
	if err == nil {
		t.Fatalf("args %v: want error, got nil (stdout=%s)", args, stdout.String())
	}
	if code := cli.ExitCode(err); code != 2 {
		t.Errorf("args %v: ExitCode = %d, want 2", args, code)
	}
	if stdout.Len() == 0 {
		t.Fatalf("args %v: --json must emit the envelope on stdout, got empty (Spec §1841); stderr=%s", args, stderr.String())
	}
	opts := []jsontestutil.AssertOption{
		jsontestutil.WithCommand(command),
		jsontestutil.WithExitCode(2),
	}
	if full {
		jsontestutil.AssertFullEnvelope(t, stdout.Bytes(), opts...)
	} else {
		jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(), opts...)
	}
}

func TestConsolidation_AddArgsError_Envelope(t *testing.T) {
	cases := []struct {
		name string
		args []string
		full bool
	}{
		// --dry-run/--diff are LOCAL subcommand flags → they must
		// follow the subcommand (a root-position local flag is an
		// "unknown flag" before arg validation even runs).
		{"noarg-minimal", []string{"--json", "add"}, false},
		{"noarg-full", []string{"--json", "add", "--dry-run"}, true},
		{"toomany-minimal", []string{"--json", "add", "postgres", "extra"}, false},
		{"toomany-full", []string{"--json", "add", "postgres", "extra", "--diff"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertArgsErrorEnvelope(t, newAppWithAddStub(&addUseCaseStub{}), tc.args, "add", tc.full)
		})
	}
}

func TestConsolidation_GenerateArgsError_Envelope(t *testing.T) {
	cases := []struct {
		name string
		args []string
		full bool
	}{
		{"noarg-minimal", []string{"--json", "generate"}, false},
		{"noarg-full", []string{"--json", "generate", "--dry-run"}, true},
		{"toomany-minimal", []string{"--json", "generate", "readme", "extra"}, false},
		{"toomany-full", []string{"--json", "generate", "readme", "extra", "--diff"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertArgsErrorEnvelope(t, newAppWithGenerateStub(&fakeGenerateUseCase{}), tc.args, "generate", tc.full)
		})
	}
}

func TestConsolidation_InitArgsError_Envelope(t *testing.T) {
	// init uses MaximumNArgs(1) (SD-C): only len>1 errors.
	cases := []struct {
		name string
		args []string
		full bool
	}{
		{"toomany-minimal", []string{"--json", "init", "proj", "extra"}, false},
		{"toomany-full", []string{"--json", "init", "proj", "extra", "--dry-run"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertArgsErrorEnvelope(t, newAppWithInitStub(&initUseCaseStub{}), tc.args, "init", tc.full)
		})
	}
}

// TestConsolidation_InitNoArg_StaysSuccess pins SD-C: init with zero
// positional args remains a valid success (default project name),
// NOT an args error — the MaximumNArgs(1) base only rejects len>1.
func TestConsolidation_InitNoArg_StaysSuccess(t *testing.T) {
	stub := &initUseCaseStub{
		resp: driving.InitProjectResponse{Project: domain.NewProject(mustProjectName(t, "myproj"))},
	}
	app := newAppWithInitStub(stub)

	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "init"}, &stdout, &stderr); err != nil {
		t.Fatalf("init --json (0 args) must succeed, got %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("init"),
		jsontestutil.WithExitCode(0),
	)
}

// ----------------------------------------------------------------------
// T3 — greedy sanitizeBaseDir on the UC-error path (Path-Leak-Defense)
// ----------------------------------------------------------------------

// assertSanitizedNoLeak pins that the single diagnostic message has
// the absolute baseDir stripped (→ `.`) but keeps the project-relative
// remainder. Mirrors remove's TestRemove_FSErrorWithAbsolutePath.
func assertSanitizedNoLeak(t *testing.T, raw []byte, abs, rel string) {
	t.Helper()
	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v\nraw=%s", err, raw)
	}
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 1 {
		t.Fatalf("want 1 diagnostic, got %d (raw=%s)", len(diags), raw)
	}
	diag, _ := diags[0].(map[string]any)
	msg, _ := diag["message"].(string)
	if strings.Contains(msg, abs) {
		t.Errorf("path-leak: absolute baseDir %q must not appear in diagnostic.message; got %q", abs, msg)
	}
	if !strings.Contains(msg, rel) {
		t.Errorf("sanitized message should keep project-relative %q; got %q", rel, msg)
	}
}

func TestConsolidation_AddFSError_SanitizesPath(t *testing.T) {
	const cwd = "/tmp/u-boot-add-test/demo"
	stub := &addUseCaseStub{err: fmt.Errorf("add write %s/u-boot.yaml: %w: %w",
		cwd, driving.ErrAddFileSystem, errors.New("disk full"))}
	app := newAppWithAddStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "add", "postgres"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrAddFileSystem) {
		t.Fatalf("errors.Is(ErrAddFileSystem) broken by sanitizer; got %v", err)
	}
	if code := cli.ExitCode(err); code != 14 {
		t.Errorf("exit: want 14, got %d", code)
	}
	assertSanitizedNoLeak(t, stdout.Bytes(), cwd, "u-boot.yaml")
}

func TestConsolidation_GenerateFSError_SanitizesPath(t *testing.T) {
	const cwd = "/tmp/u-boot-generate-test"
	stub := &fakeGenerateUseCase{err: fmt.Errorf("generate write %s/README.md: %w: %w",
		cwd, driving.ErrGenerateFileSystem, errors.New("disk full"))}
	app := newAppWithGenerateStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "generate", "readme"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrGenerateFileSystem) {
		t.Fatalf("errors.Is(ErrGenerateFileSystem) broken by sanitizer; got %v", err)
	}
	if code := cli.ExitCode(err); code != 14 {
		t.Errorf("exit: want 14, got %d", code)
	}
	assertSanitizedNoLeak(t, stdout.Bytes(), cwd, "README.md")
}

func TestConsolidation_InitFSError_SanitizesPath(t *testing.T) {
	const cwd = "/tmp/u-boot-init-test"
	stub := &initUseCaseStub{err: fmt.Errorf("init write %s/u-boot.yaml: %w: %w",
		cwd, driving.ErrInitFileSystem, errors.New("disk full"))}
	app := newAppWithInitStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "init", "myproj"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrInitFileSystem) {
		t.Fatalf("errors.Is(ErrInitFileSystem) broken by sanitizer; got %v", err)
	}
	if code := cli.ExitCode(err); code != 14 {
		t.Errorf("exit: want 14, got %d", code)
	}
	assertSanitizedNoLeak(t, stdout.Bytes(), cwd, "u-boot.yaml")
}
