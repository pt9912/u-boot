package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli/jsontestutil"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// allOKReport seeds the fake doctor with two OK items so the All-OK
// envelope case ships diagnostics: [] (Spec §1834 forbids
// level: "ok", §1846-1852 example pins the empty array).
func allOKReport() driving.DoctorResponse {
	return driving.DoctorResponse{Report: domain.DiagnosticReport{
		Items: []domain.Diagnostic{
			{ID: "docker.installed", Severity: domain.SeverityOK, Message: "Docker is installed."},
			{ID: "git.installed", Severity: domain.SeverityOK, Message: "Git is installed."},
		},
	}}
}

func warnReport() driving.DoctorResponse {
	return driving.DoctorResponse{Report: domain.DiagnosticReport{
		Items: []domain.Diagnostic{
			{ID: "docker.installed", Severity: domain.SeverityOK, Message: "Docker is installed."},
			{ID: "devcontainer.json.valid", Severity: domain.SeverityWarn, Message: "devcontainer.json missing 'features' key."},
		},
	}}
}

func errorReport() driving.DoctorResponse {
	return driving.DoctorResponse{Report: domain.DiagnosticReport{
		Items: []domain.Diagnostic{
			{ID: "docker.reachable", Severity: domain.SeverityError, Message: "Docker daemon not reachable."},
			{ID: "devcontainer.json.valid", Severity: domain.SeverityWarn, Message: "Minor issue."},
		},
	}}
}

// TestDoctorJSON_AllOK pins the canonical empty-diagnostics case
// (Lastenheft §1846-1852 example).
func TestDoctorJSON_AllOK(t *testing.T) {
	stdout := &bytes.Buffer{}
	app := newAppWithDoctor(&fakeInitUseCase{}, &fakeDoctorUseCase{resp: allOKReport()},
		cli.WithGetwd(func() (string, error) { return "/tmp/proj", nil }))
	err := app.Execute(context.Background(), []string{"--json", "doctor"}, stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw := bytes.TrimSpace(stdout.Bytes())
	jsontestutil.AssertMinimalEnvelope(t, raw,
		jsontestutil.WithCommand("doctor"),
		jsontestutil.WithExitCode(0),
	)

	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env["status"] != "ok" {
		t.Errorf("status: want %q, got %v", "ok", env["status"])
	}
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 0 {
		t.Errorf("All-OK case must ship diagnostics: [] (Spec §1834 forbids level:ok), got %v", diags)
	}
}

// TestDoctorJSON_Warn pins the Warn-Fall: status="warn",
// diagnostics array hat den Warn-Eintrag, exitCode=0 (non-strict).
func TestDoctorJSON_Warn(t *testing.T) {
	stdout := &bytes.Buffer{}
	app := newAppWithDoctor(&fakeInitUseCase{}, &fakeDoctorUseCase{resp: warnReport()},
		cli.WithGetwd(func() (string, error) { return "/tmp/proj", nil }))
	err := app.Execute(context.Background(), []string{"--json", "doctor"}, stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw := bytes.TrimSpace(stdout.Bytes())
	jsontestutil.AssertMinimalEnvelope(t, raw,
		jsontestutil.WithCommand("doctor"),
		jsontestutil.WithExitCode(0),
		jsontestutil.WithExpectedCodes("devcontainer.json.valid"),
	)

	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env["status"] != "warn" {
		t.Errorf("status: want %q, got %v", "warn", env["status"])
	}
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 1 {
		t.Errorf("Warn-Fall must ship one diagnostics entry (OK filtered out), got %d", len(diags))
	}
}

// TestDoctorJSON_Error pins the Error-Fall: status="error", exitCode=11.
func TestDoctorJSON_Error(t *testing.T) {
	stdout := &bytes.Buffer{}
	app := newAppWithDoctor(&fakeInitUseCase{}, &fakeDoctorUseCase{resp: errorReport()},
		cli.WithGetwd(func() (string, error) { return "/tmp/proj", nil }))
	err := app.Execute(context.Background(), []string{"--json", "doctor"}, stdout, &bytes.Buffer{})
	if !errors.Is(err, cli.ErrDoctorFailures) {
		t.Fatalf("want ErrDoctorFailures, got %v", err)
	}
	if cli.ExitCode(err) != 11 {
		t.Errorf("want exit code 11, got %d", cli.ExitCode(err))
	}

	raw := bytes.TrimSpace(stdout.Bytes())
	jsontestutil.AssertMinimalEnvelope(t, raw,
		jsontestutil.WithCommand("doctor"),
		jsontestutil.WithExitCode(11),
	)

	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env["status"] != "error" {
		t.Errorf("status: want %q, got %v", "error", env["status"])
	}
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 2 {
		t.Errorf("Error-Fall must ship 2 diagnostics entries (Error + Warn, OK filtered), got %d", len(diags))
	}
}

// TestDoctorJSON_QuietIsSemanticNoOp pins the slice-doctor T0-(e)
// quiet/json-interaction decision (Review-Finding M4: semantisch
// identisch, nicht byte-identisch — robust gegen zukünftige
// zeitabhängige Messages).
func TestDoctorJSON_QuietIsSemanticNoOp(t *testing.T) {
	stdout1, stdout2 := &bytes.Buffer{}, &bytes.Buffer{}
	app1 := newAppWithDoctor(&fakeInitUseCase{}, &fakeDoctorUseCase{resp: warnReport()},
		cli.WithGetwd(func() (string, error) { return "/tmp/proj", nil }))
	app2 := newAppWithDoctor(&fakeInitUseCase{}, &fakeDoctorUseCase{resp: warnReport()},
		cli.WithGetwd(func() (string, error) { return "/tmp/proj", nil }))

	if err := app1.Execute(context.Background(), []string{"--json", "doctor"}, stdout1, &bytes.Buffer{}); err != nil {
		t.Fatalf("--json doctor: %v", err)
	}
	if err := app2.Execute(context.Background(), []string{"--json", "--quiet", "doctor"}, stdout2, &bytes.Buffer{}); err != nil {
		t.Fatalf("--json --quiet doctor: %v", err)
	}

	env1, env2 := parseEnv(t, stdout1.Bytes()), parseEnv(t, stdout2.Bytes())
	if env1["status"] != env2["status"] {
		t.Errorf("status mismatch: --json=%v --quiet --json=%v", env1["status"], env2["status"])
	}
	if env1["exitCode"] != env2["exitCode"] {
		t.Errorf("exitCode mismatch: --json=%v --quiet --json=%v", env1["exitCode"], env2["exitCode"])
	}
	d1 := codesAndLevels(env1["diagnostics"])
	d2 := codesAndLevels(env2["diagnostics"])
	if strings.Join(d1, "|") != strings.Join(d2, "|") {
		t.Errorf("diagnostics sequence mismatch:\n  --json:        %v\n  --quiet --json: %v", d1, d2)
	}
}

// TestDoctorJSON_StrictWarnExits11 pins that --strict --json with
// a Warn report still upgrades exitCode to 11; status remains
// "warn" because Spec §1837 couples status to the highest level,
// not to --strict.
func TestDoctorJSON_StrictWarnExits11(t *testing.T) {
	stdout := &bytes.Buffer{}
	app := newAppWithDoctor(&fakeInitUseCase{}, &fakeDoctorUseCase{resp: warnReport()},
		cli.WithGetwd(func() (string, error) { return "/tmp/proj", nil }))
	err := app.Execute(context.Background(), []string{"--json", "doctor", "--strict"}, stdout, &bytes.Buffer{})
	if !errors.Is(err, cli.ErrDoctorFailures) {
		t.Fatalf("want ErrDoctorFailures, got %v", err)
	}
	if cli.ExitCode(err) != 11 {
		t.Errorf("want exit code 11, got %d", cli.ExitCode(err))
	}

	env := parseEnv(t, stdout.Bytes())
	if env["status"] != "warn" {
		t.Errorf("status: want %q (Spec §1837 couples to highest level), got %v", "warn", env["status"])
	}
	if env["exitCode"] != float64(11) {
		t.Errorf("exitCode: want 11, got %v", env["exitCode"])
	}
}

// TestDoctorJSON_NoIndent pins that the wire format is single-line
// JSON (no SetIndent) — important for line-based machine consumers.
func TestDoctorJSON_NoIndent(t *testing.T) {
	stdout := &bytes.Buffer{}
	app := newAppWithDoctor(&fakeInitUseCase{}, &fakeDoctorUseCase{resp: allOKReport()},
		cli.WithGetwd(func() (string, error) { return "/tmp/proj", nil }))
	if err := app.Execute(context.Background(), []string{"--json", "doctor"}, stdout, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	body := strings.TrimSpace(stdout.String())
	if strings.Count(body, "\n") != 0 {
		t.Errorf("expected single-line JSON, got %d newlines: %q", strings.Count(body, "\n"), body)
	}
}

func parseEnv(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var env map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(raw), &env); err != nil {
		t.Fatalf("unmarshal: %v\nraw=%s", err, raw)
	}
	return env
}

func codesAndLevels(diags any) []string {
	arr, _ := diags.([]any)
	out := make([]string, 0, len(arr))
	for _, raw := range arr {
		item, _ := raw.(map[string]any)
		level, _ := item["level"].(string)
		code, _ := item["code"].(string)
		out = append(out, level+":"+code)
	}
	return out
}
