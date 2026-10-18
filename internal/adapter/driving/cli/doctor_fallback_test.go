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

// fallbackStateItems are the five states of one
// `devcontainer.sandbox.fallback.*` item (LH-FA-DEV-007, doctor mapping).
func fallbackStateItems(id string) map[string]domain.Diagnostic {
	return map[string]domain.Diagnostic{
		"ok":             {ID: id, Severity: domain.SeverityOK, Message: "reported ok state."},
		"fallback":       {ID: id, Severity: domain.SeverityWarn, DegradationFallback: true, Message: "reported fallback state.", Hint: "fix the start"},
		"failed":         {ID: id, Severity: domain.SeverityError, Message: "reported failed state.", Hint: "fix the start"},
		"undeterminable": {ID: id, Severity: domain.SeverityInfo, Message: "state not determinable here.", Hint: "rebuild the container"},
		"invalid":        {ID: id, Severity: domain.SeverityWarn, Message: "status file is invalid."},
	}
}

func runDoctorWith(t *testing.T, items []domain.Diagnostic, args ...string) (string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	app := newAppWithDoctor(&fakeInitUseCase{}, &fakeDoctorUseCase{resp: driving.DoctorResponse{Report: domain.DiagnosticReport{Items: items}}},
		cli.WithGetwd(func() (string, error) { return "/tmp/x", nil }))
	err := app.Execute(context.Background(), args, &stdout, &stderr)
	return stdout.String(), err
}

// LH-FA-CLI-005 exception 2: under `--quiet` the fallback item and the
// failed item (an error) stay, everything else of the same ID is hidden:
// ok, not determinable, and the invalid-document warning without the flag.
func TestExecute_Doctor_Quiet_FallbackStates(t *testing.T) {
	for _, id := range []string{"devcontainer.sandbox.fallback.vfs", "devcontainer.sandbox.fallback.daemon", "devcontainer.sandbox.fallback.egress"} {
		t.Run(id, func(t *testing.T) {
			states := fallbackStateItems(id)
			var items []domain.Diagnostic
			for _, s := range []string{"ok", "fallback", "failed", "undeterminable", "invalid"} {
				items = append(items, states[s])
			}
			out, err := runDoctorWith(t, items, "--quiet", "doctor")
			if !errors.Is(err, cli.ErrDoctorFailures) {
				t.Fatalf("err = %v, want ErrDoctorFailures (the failed item is an error)", err)
			}
			for _, want := range []string{"reported fallback state.", "reported failed state.", "fix the start"} {
				if !strings.Contains(out, want) {
					t.Errorf("--quiet hid %q: %q", want, out)
				}
			}
			for _, gone := range []string{"reported ok state", "not determinable here", "status file is invalid", "rebuild the container", "Summary:"} {
				if strings.Contains(out, gone) {
					t.Errorf("--quiet kept %q: %q", gone, out)
				}
			}
		})
	}
}

// The info item of an undeterminable state alone prints nothing under
// `--quiet` and does not change the exit code.
func TestExecute_Doctor_Quiet_UndeterminableOnlyIsSilent(t *testing.T) {
	out, err := runDoctorWith(t, []domain.Diagnostic{fallbackStateItems("devcontainer.sandbox.fallback.vfs")["undeterminable"]}, "--quiet", "doctor")
	if err != nil || out != "" {
		t.Errorf("err = %v, stdout = %q, want silence and exit 0", err, out)
	}
}

// Without `--quiet` the report names every item, the info item of an
// undeterminable state included, with its hint; the summary line counts
// error, warn and ok items as before.
func TestExecute_Doctor_FullReport_InfoItem(t *testing.T) {
	states := fallbackStateItems("devcontainer.sandbox.fallback.vfs")
	ok := domain.Diagnostic{ID: "fs.write-permissions", Severity: domain.SeverityOK, Message: "BaseDir is writable."}
	out, err := runDoctorWith(t, []domain.Diagnostic{ok, states["undeterminable"], states["fallback"]}, "doctor")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	for _, want := range []string{"devcontainer.sandbox.fallback.vfs state not determinable here.", "⚠  devcontainer.sandbox.fallback.vfs reported fallback state.",
		"→ rebuild the container", "→ fix the start", "Summary: 0 error, 1 warn, 1 ok"} {
		if !strings.Contains(out, want) {
			t.Errorf("report lacks %q:\n%s", want, out)
		}
	}
}

// `doctor --strict` treats the flagged warning like any warning (exit 11)
// but not an info item, and names the item ID in the error.
func TestExecute_Doctor_Strict_FallbackWarnVsInfo(t *testing.T) {
	states := fallbackStateItems("devcontainer.sandbox.fallback.daemon")
	_, err := runDoctorWith(t, []domain.Diagnostic{states["fallback"]}, "doctor", "--strict")
	if !errors.Is(err, cli.ErrDoctorFailures) || cli.ExitCode(err) != 11 {
		t.Fatalf("err = %v, want exit 11 for the fallback warning", err)
	}
	if !strings.Contains(err.Error(), "devcontainer.sandbox.fallback.daemon") {
		t.Errorf("error text lacks the item ID: %v", err)
	}
	if _, err := runDoctorWith(t, []domain.Diagnostic{states["undeterminable"], states["ok"]}, "doctor", "--strict"); err != nil {
		t.Errorf("info and ok items must not fail --strict: %v", err)
	}
	if _, err := runDoctorWith(t, []domain.Diagnostic{states["invalid"]}, "doctor", "--strict"); !errors.Is(err, cli.ErrDoctorFailures) {
		t.Errorf("the invalid-document warning fails --strict like any warning: %v", err)
	}
}

// `--json`: the diagnostics carry the codes of the registry; warn and error
// items appear, the info and ok items do not (LH-NFA-USE-004).
func TestDoctorJSON_FallbackItems(t *testing.T) {
	var items []domain.Diagnostic
	for _, id := range []string{"devcontainer.sandbox.fallback.vfs", "devcontainer.sandbox.fallback.daemon", "devcontainer.sandbox.fallback.egress"} {
		states := fallbackStateItems(id)
		items = append(items, states["fallback"], states["undeterminable"], states["ok"])
	}
	items = append(items, fallbackStateItems("devcontainer.sandbox.fallback.vfs")["failed"])
	stdout := &bytes.Buffer{}
	app := newAppWithDoctor(&fakeInitUseCase{}, &fakeDoctorUseCase{resp: driving.DoctorResponse{Report: domain.DiagnosticReport{Items: items}}},
		cli.WithGetwd(func() (string, error) { return "/tmp/proj", nil }))
	err := app.Execute(context.Background(), []string{"--quiet", "--json", "doctor"}, stdout, &bytes.Buffer{})
	if !errors.Is(err, cli.ErrDoctorFailures) {
		t.Fatalf("err = %v, want the failed item to end with exit 11", err)
	}
	raw := bytes.TrimSpace(stdout.Bytes())
	jsontestutil.AssertMinimalEnvelope(t, raw, jsontestutil.WithCommand("doctor"), jsontestutil.WithExitCode(11))
	var env struct {
		Diagnostics []struct{ Level, Code, Message string } `json:"diagnostics"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatal(err)
	}
	var warns, errs int
	for _, d := range env.Diagnostics {
		switch d.Level {
		case "warn":
			warns++
			if d.Message != "reported fallback state." {
				t.Errorf("unexpected warn message %q", d.Message)
			}
		case "error":
			errs++
		default:
			t.Errorf("level %q in %+v", d.Level, d)
		}
	}
	if warns != 3 || errs != 1 {
		t.Errorf("warn=%d error=%d in %+v, want 3 warn and 1 error (--quiet does not filter --json)", warns, errs, env.Diagnostics)
	}
}
