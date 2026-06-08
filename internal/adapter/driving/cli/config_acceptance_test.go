package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

// TestConfigJSON_QuietJSONIdenticalToJSON pins Cluster-T0-(a)
// doctor-Pattern for all three forms: `--quiet --json` is
// semantically identical to `--json` (quiet is a no-op in JSON
// mode). Same exit 0, same envelope shape.
func TestConfigJSON_QuietJSONIdenticalToJSON(t *testing.T) {
	cases := []struct {
		name string
		args []string
		uc   *fakeConfigUseCase
		sub  string
	}{
		{"show", []string{"config"}, &fakeConfigUseCase{showResp: driving.ConfigShowResponse{Body: []byte("schemaVersion: 1\n")}}, "show"},
		{"get", []string{"config", "get", "project.name"}, &fakeConfigUseCase{getResp: driving.ConfigGetResponse{Path: mustCfgPath(t, "project.name"), Value: "demo"}}, "get"},
		{"set", []string{"config", "set", "project.name", "x"}, &fakeConfigUseCase{setResp: driving.ConfigSetResponse{Path: mustCfgPath(t, "project.name"), OldValue: "a", NewValue: "x"}}, "set"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := newAppWithConfigStub(tc.uc)
			var out bytes.Buffer
			args := append([]string{"--quiet", "--json"}, tc.args...)
			if err := app.Execute(context.Background(), args, &out, &bytes.Buffer{}); err != nil {
				t.Fatalf("execute: %v", err)
			}
			jsontestutil.AssertMinimalEnvelope(t, out.Bytes(),
				jsontestutil.WithCommand("config"),
				jsontestutil.WithSubcommand(tc.sub),
				jsontestutil.WithExitCode(0),
			)
		})
	}
}

// TestConfigJSON_UnknownSubcommandEmitsEnvelope pins R2-HIGH-3/
// R3-MED-1: `u-boot config foo` (foo is no registered sub-verb)
// dispatches to the bare command's configArgsValidator (cobra.NoArgs
// with args=["foo"]) which emits an Envelope-konformen Exit-2 reject
// instead of Cobra's raw stderr "unknown command". Guards against
// Cobra-mechanic drift on version upgrades.
func TestConfigJSON_UnknownSubcommandEmitsEnvelope(t *testing.T) {
	app := newAppWithConfigStub(&fakeConfigUseCase{})
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "config", "foo"}, &stdout, &stderr)
	if cli.ExitCode(err) != 2 {
		t.Fatalf("exit = %d, want 2 (err=%v)", cli.ExitCode(err), err)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("config"),
		jsontestutil.WithSubcommand("show"),
		jsontestutil.WithExitCode(2),
	)
}

// TestConfigJSON_HelpEdgeCaseNoEnvelope pins R1-MED-6: `config
// --help --json` runs through the Help-Escape-Hatch (no RunE, no
// envelope) — Cobra renders help on stdout. The subcommand-pflicht
// applies only to RunE-emitted envelopes.
func TestConfigJSON_HelpEdgeCaseNoEnvelope(t *testing.T) {
	app := newAppWithConfigStub(&fakeConfigUseCase{})
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "config", "--help"}, &stdout, &stderr); err != nil {
		t.Fatalf("help should not error: %v", err)
	}
	if strings.Contains(stdout.String(), "\"exitCode\"") {
		t.Errorf("--help --json must NOT emit an envelope; got %s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Errorf("--help should render usage text; got %s", stdout.String())
	}
}

// TestConfigJSON_CONF005Disambiguation pins R3-MED-2: the three
// LH-FA-CONF-005 sentinels (Path-Unknown / Write-Rejected /
// Value-Not-Set) all carry code LH-FA-CONF-005 + exit 10, but the
// diagnostic messages differ so consumers can disambiguate by
// message-prefix (code alone is insufficient).
func TestConfigJSON_CONF005Disambiguation(t *testing.T) {
	messages := map[string]string{}
	type row struct {
		name string
		args []string
		uc   *fakeConfigUseCase
	}
	rows := []row{
		{"path-unknown", []string{"--json", "config", "get", "bogus.path"}, &fakeConfigUseCase{}},
		{"write-rejected", []string{"--json", "config", "set", "services.postgres.enabled", "true"}, &fakeConfigUseCase{setErr: driving.ErrConfigWriteRejected}},
		{"value-not-set", []string{"--json", "config", "get", "devcontainer.enabled"}, &fakeConfigUseCase{getErr: driving.ErrConfigValueNotSet}},
	}
	for _, r := range rows {
		t.Run(r.name, func(t *testing.T) {
			app := newAppWithConfigStub(r.uc)
			var out bytes.Buffer
			err := app.Execute(context.Background(), r.args, &out, &bytes.Buffer{})
			if cli.ExitCode(err) != 10 {
				t.Fatalf("exit = %d, want 10 (err=%v)", cli.ExitCode(err), err)
			}
			var env map[string]any
			if e := json.Unmarshal(out.Bytes(), &env); e != nil {
				t.Fatalf("unmarshal: %v", e)
			}
			diags, _ := env["diagnostics"].([]any)
			if len(diags) != 1 {
				t.Fatalf("want 1 diagnostic, got %v", diags)
			}
			d, _ := diags[0].(map[string]any)
			if d["code"] != "LH-FA-CONF-005" {
				t.Errorf("code = %v, want LH-FA-CONF-005", d["code"])
			}
			messages[r.name] = d["message"].(string)
		})
	}
	if len(messages) == 3 {
		if messages["path-unknown"] == messages["write-rejected"] ||
			messages["write-rejected"] == messages["value-not-set"] ||
			messages["path-unknown"] == messages["value-not-set"] {
			t.Errorf("LH-FA-CONF-005 messages must be pairwise distinct for disambiguation: %+v", messages)
		}
	}
}

// TestConfigJSON_SanitizerStripsAbsolutePath pins T0-(p): an FS
// error carrying the absolute baseDir path is sanitized to a
// project-relative path before it reaches diagnostic.message (no
// filesystem-layout leak in the machine-readable output).
func TestConfigJSON_SanitizerStripsAbsolutePath(t *testing.T) {
	uc := &fakeConfigUseCase{setErr: fmt.Errorf(
		"%w: write %q: permission denied",
		driving.ErrConfigFileSystem, "/tmp/u-boot-config-test/demo/u-boot.yaml")}
	app := newAppWithConfigStub(uc)
	var out bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "config", "set", "project.name", "x"}, &out, &bytes.Buffer{})
	if cli.ExitCode(err) != 14 {
		t.Fatalf("exit = %d, want 14 (FS)", cli.ExitCode(err))
	}
	if strings.Contains(out.String(), "/tmp/u-boot-config-test") {
		t.Errorf("absolute baseDir leaked into envelope: %s", out.String())
	}
	if !strings.Contains(out.String(), "u-boot.yaml") {
		t.Errorf("expected project-relative u-boot.yaml in message; got %s", out.String())
	}
}

// TestConfigJSON_SubcommandAlwaysSet pins T0-(h)/§322: every
// RunE-emitted config envelope carries a non-empty subcommand,
// across success AND error paths and all three forms.
func TestConfigJSON_SubcommandAlwaysSet(t *testing.T) {
	cases := []struct {
		name string
		args []string
		uc   *fakeConfigUseCase
		want string
	}{
		{"show-success", []string{"--json", "config"}, &fakeConfigUseCase{}, "show"},
		{"get-success", []string{"--json", "config", "get", "project.name"}, &fakeConfigUseCase{getResp: driving.ConfigGetResponse{Path: mustCfgPath(t, "project.name"), Value: "x"}}, "get"},
		{"set-success", []string{"--json", "config", "set", "project.name", "x"}, &fakeConfigUseCase{}, "set"},
		{"set-error", []string{"--json", "config", "set", "project.name", "x"}, &fakeConfigUseCase{setErr: driving.ErrConfigValueInvalid}, "set"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := newAppWithConfigStub(tc.uc)
			var out bytes.Buffer
			_ = app.Execute(context.Background(), tc.args, &out, &bytes.Buffer{})
			var env map[string]any
			if err := json.Unmarshal(out.Bytes(), &env); err != nil {
				t.Fatalf("unmarshal: %v (out=%s)", err, out.String())
			}
			if env["subcommand"] != tc.want {
				t.Errorf("subcommand = %v, want %q (§322 subcommand-pflicht)", env["subcommand"], tc.want)
			}
		})
	}
}

// TestConfigJSON_SetErrorShapesByPreviewFlag pins the Mid-Stage
// envelope-shape contract (T0-(o)): a set error without preview
// flags emits the minimal shape (no dryRun field); with --dry-run
// it emits the voll-schema shape (dryRun:true present).
func TestConfigJSON_SetErrorShapesByPreviewFlag(t *testing.T) {
	t.Run("plain-minimal", func(t *testing.T) {
		app := newAppWithConfigStub(&fakeConfigUseCase{setErr: driving.ErrConfigValueInvalid})
		var out bytes.Buffer
		_ = app.Execute(context.Background(), []string{"--json", "config", "set", "project.name", "x"}, &out, &bytes.Buffer{})
		var env map[string]any
		_ = json.Unmarshal(out.Bytes(), &env)
		if _, present := env["dryRun"]; present {
			t.Errorf("plain set error must be minimal (no dryRun field); got %s", out.String())
		}
	})
	t.Run("dryrun-full", func(t *testing.T) {
		app := newAppWithConfigStub(&fakeConfigUseCase{setErr: driving.ErrConfigValueInvalid})
		var out bytes.Buffer
		_ = app.Execute(context.Background(), []string{"--json", "config", "set", "project.name", "x", "--dry-run"}, &out, &bytes.Buffer{})
		var env map[string]any
		_ = json.Unmarshal(out.Bytes(), &env)
		if env["dryRun"] != true {
			t.Errorf("dry-run set error must be voll-schema (dryRun:true); got %s", out.String())
		}
	})
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
