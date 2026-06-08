package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli/jsontestutil"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// newAppWithConfigStub wires the config use-case fake plus a
// deterministic getwd so tests do not depend on the runner's CWD.
func newAppWithConfigStub(uc driving.ConfigUseCase) *cli.App {
	return newAppWithConfig(uc, cli.WithGetwd(func() (string, error) { return "/tmp/u-boot-config-test/demo", nil }))
}

func mustCfgPath(t *testing.T, raw string) domain.ConfigPath {
	t.Helper()
	p, err := domain.NewConfigPath(raw)
	if err != nil {
		t.Fatalf("NewConfigPath(%q): %v", raw, err)
	}
	return p
}

// TestConfigJSON_ShowMinimalDataEnvelope pins bare `config --json`:
// minimal+data envelope, subcommand "show" (T0-(b)), data.body the
// full u-boot.yaml as a string (T0-(c)).
func TestConfigJSON_ShowMinimalDataEnvelope(t *testing.T) {
	uc := &fakeConfigUseCase{showResp: driving.ConfigShowResponse{Body: []byte("schemaVersion: 1\nproject:\n  name: demo\n")}}
	app := newAppWithConfigStub(uc)

	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "config"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("config"),
		jsontestutil.WithSubcommand("show"),
		jsontestutil.WithExitCode(0),
		jsontestutil.WithDataKeyPresent("body", "schemaVersion: 1\nproject:\n  name: demo\n"),
	)
}

// TestConfigJSON_GetMinimalDataEnvelope pins `config get --json`:
// subcommand "get", data.path + data.value (T0-(c)).
func TestConfigJSON_GetMinimalDataEnvelope(t *testing.T) {
	uc := &fakeConfigUseCase{getResp: driving.ConfigGetResponse{Path: mustCfgPath(t, "project.name"), Value: "demo"}}
	app := newAppWithConfigStub(uc)

	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "config", "get", "project.name"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("config"),
		jsontestutil.WithSubcommand("get"),
		jsontestutil.WithExitCode(0),
		jsontestutil.WithDataKeyPresent("path", "project.name"),
		jsontestutil.WithDataKeyPresent("value", "demo"),
	)
}

// TestConfigJSON_SetPlainDataEnvelope pins `config set --json` w/o
// preview flags: minimal+data with the configSetData carrier
// (T0-(c)/(d)), noOp=false on a real change.
func TestConfigJSON_SetPlainDataEnvelope(t *testing.T) {
	uc := &fakeConfigUseCase{setResp: driving.ConfigSetResponse{
		Path: mustCfgPath(t, "project.name"), OldValue: "old", NewValue: "new",
	}}
	app := newAppWithConfigStub(uc)

	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "config", "set", "project.name", "new"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("config"),
		jsontestutil.WithSubcommand("set"),
		jsontestutil.WithExitCode(0),
		jsontestutil.WithDataKeyPresent("oldValue", "old"),
		jsontestutil.WithDataKeyPresent("newValue", "new"),
		jsontestutil.WithDataKeyPresent("noOp", false),
	)
}

// TestConfigJSON_SetNoOp pins the idempotent path: OldValue ==
// NewValue → data.noOp=true, no plannedFiles in the plain envelope.
func TestConfigJSON_SetNoOp(t *testing.T) {
	uc := &fakeConfigUseCase{setResp: driving.ConfigSetResponse{
		Path: mustCfgPath(t, "project.name"), OldValue: "demo", NewValue: "demo",
	}}
	app := newAppWithConfigStub(uc)

	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "config", "set", "project.name", "demo"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v", err)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("config"),
		jsontestutil.WithSubcommand("set"),
		jsontestutil.WithExitCode(0),
		jsontestutil.WithDataKeyPresent("noOp", true),
	)
}

// TestConfigJSON_SetDryRunFullEnvelope pins `config set --dry-run
// --json`: voll-schema, dryRun=true, plannedFiles from the
// recorder-surfaced resp.PlannedFiles (R-T4-1).
func TestConfigJSON_SetDryRunFullEnvelope(t *testing.T) {
	uc := &fakeConfigUseCase{setResp: driving.ConfigSetResponse{
		Path: mustCfgPath(t, "project.name"), OldValue: "old", NewValue: "new",
		PlannedFiles: []driving.PlannedFile{
			{Path: "u-boot.yaml", Action: "modify", OldContent: []byte("name: old\n"), NewContent: []byte("name: new\n")},
		},
	}}
	app := newAppWithConfigStub(uc)

	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"config", "set", "project.name", "new", "--dry-run", "--json"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("config"),
		jsontestutil.WithSubcommand("set"),
		jsontestutil.WithExitCode(0),
	)
	var env map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env["dryRun"] != true {
		t.Errorf("dryRun = %v, want true", env["dryRun"])
	}
}

// TestConfigJSON_SetDiffHunks pins `config set --diff --json`:
// voll-schema, diff=true, plannedFiles[].hunks rendered from the
// surfaced byte content (R-T4-1 makes the bytes available).
func TestConfigJSON_SetDiffHunks(t *testing.T) {
	uc := &fakeConfigUseCase{setResp: driving.ConfigSetResponse{
		Path: mustCfgPath(t, "project.name"), OldValue: "old", NewValue: "new",
		PlannedFiles: []driving.PlannedFile{
			{Path: "u-boot.yaml", Action: "modify", OldContent: []byte("name: old\n"), NewContent: []byte("name: new\n")},
		},
	}}
	app := newAppWithConfigStub(uc)

	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"config", "set", "project.name", "new", "--diff", "--json"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("config"),
		jsontestutil.WithSubcommand("set"),
		jsontestutil.WithExitCode(0),
	)
	if !strings.Contains(stdout.String(), "hunks") {
		t.Errorf("--diff --json envelope should carry hunks; got %s", stdout.String())
	}
}

// TestConfigJSON_DryRunRejectedOnBare pins T0-(g): --dry-run on the
// read-only bare form is rejected Envelope-konform with
// ErrDryRunNotApplicable / Exit 2.
func TestConfigJSON_DryRunRejectedOnBare(t *testing.T) {
	app := newAppWithConfigStub(&fakeConfigUseCase{})

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "config", "--dry-run"}, &stdout, &stderr)
	if !strings.Contains(err.Error(), "only valid for") {
		t.Fatalf("want ErrDryRunNotApplicable, got %v", err)
	}
	if cli.ExitCode(err) != 2 {
		t.Errorf("exit = %d, want 2", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("config"),
		jsontestutil.WithSubcommand("show"),
		jsontestutil.WithExitCode(2),
	)
}

// TestConfigJSON_DiffRejectedOnGet pins the symmetric reject on
// `config get --diff --json` (Exit 2, subcommand "get").
func TestConfigJSON_DiffRejectedOnGet(t *testing.T) {
	app := newAppWithConfigStub(&fakeConfigUseCase{})

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "config", "get", "project.name", "--diff"}, &stdout, &stderr)
	if cli.ExitCode(err) != 2 {
		t.Fatalf("exit = %d, want 2 (err=%v)", cli.ExitCode(err), err)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("config"),
		jsontestutil.WithSubcommand("get"),
		jsontestutil.WithExitCode(2),
	)
}

// TestConfigJSON_SetArgsMismatch pins the custom args validator
// (T0-(l)): `config set <path>` (1 arg) emits the Envelope-konformen
// Exit-2 reject on stdout instead of Cobra's raw stderr error.
func TestConfigJSON_SetArgsMismatch(t *testing.T) {
	app := newAppWithConfigStub(&fakeConfigUseCase{})

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "config", "set", "project.name"}, &stdout, &stderr)
	if cli.ExitCode(err) != 2 {
		t.Fatalf("exit = %d, want 2 (err=%v)", cli.ExitCode(err), err)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("config"),
		jsontestutil.WithSubcommand("set"),
		jsontestutil.WithExitCode(2),
	)
}

// TestConfigJSON_SetWriteRejectedMapsExit10 pins the mapper +
// ExitCode for the write-rejected class (T0-(f)/(m)): exit 10, code
// LH-FA-CONF-005, subcommand "set".
func TestConfigJSON_SetWriteRejectedMapsExit10(t *testing.T) {
	uc := &fakeConfigUseCase{setErr: driving.ErrConfigWriteRejected}
	app := newAppWithConfigStub(uc)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "config", "set", "services.postgres.enabled", "true"}, &stdout, &stderr)
	if cli.ExitCode(err) != 10 {
		t.Fatalf("exit = %d, want 10 (err=%v)", cli.ExitCode(err), err)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("config"),
		jsontestutil.WithSubcommand("set"),
		jsontestutil.WithExitCode(10),
		jsontestutil.WithExpectedCodes("LH-FA-CONF-005"),
	)
}

// TestConfigJSON_SetWarningsMapToDiagnostics pins the Orphan-Feature-
// WARN migration (T0-(n)): resp.Warnings surface as diagnostics[]
// with level "warn"; warn-only keeps exit 0, status "warn".
func TestConfigJSON_SetWarningsMapToDiagnostics(t *testing.T) {
	uc := &fakeConfigUseCase{setResp: driving.ConfigSetResponse{
		Path: mustCfgPath(t, "devcontainer.features.unknown-thing.enabled"), OldValue: "", NewValue: "true",
		Warnings: []driving.WarningEntry{{Code: "LH-FA-DEV-003", Level: "warn", Message: "orphan feature activation", Subject: "unknown-thing"}},
	}}
	app := newAppWithConfigStub(uc)

	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "config", "set", "devcontainer.features.unknown-thing.enabled", "true"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v", err)
	}
	var env map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env["status"] != "warn" {
		t.Errorf("status = %v, want warn (warn-coupling Spec §447)", env["status"])
	}
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 1 {
		t.Fatalf("want 1 diagnostic, got %d: %v", len(diags), diags)
	}
}

// TestConfig_HumanMode_ShowAndSet pins the non-JSON paths stay
// intact: bare config writes the raw body; set prints the summary.
func TestConfig_HumanMode_ShowAndSet(t *testing.T) {
	show := newAppWithConfigStub(&fakeConfigUseCase{showResp: driving.ConfigShowResponse{Body: []byte("schemaVersion: 1\n")}})
	var out bytes.Buffer
	if err := show.Execute(context.Background(), []string{"config"}, &out, &bytes.Buffer{}); err != nil {
		t.Fatalf("show: %v", err)
	}
	if out.String() != "schemaVersion: 1\n" {
		t.Errorf("human show body = %q, want raw yaml", out.String())
	}

	set := newAppWithConfigStub(&fakeConfigUseCase{setResp: driving.ConfigSetResponse{
		Path: mustCfgPath(t, "project.name"), OldValue: "old", NewValue: "new",
	}})
	var setOut bytes.Buffer
	if err := set.Execute(context.Background(), []string{"config", "set", "project.name", "new"}, &setOut, &bytes.Buffer{}); err != nil {
		t.Fatalf("set: %v", err)
	}
	if !strings.Contains(setOut.String(), "→") {
		t.Errorf("human set summary = %q, want transition arrow", setOut.String())
	}
}
