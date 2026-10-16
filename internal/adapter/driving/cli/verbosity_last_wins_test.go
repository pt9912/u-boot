package cli_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// LH-FA-CLI-005: of several verbosity options the one given last on
// the command line decides — also for the output filters, not only
// for the logger level. `--quiet` is in force exactly when it is the
// last verbosity option; a trailing --verbose or --debug lifts it.

// Init summary: --quiet drops it, a trailing --verbose/--debug keeps it.
func TestExecute_Verbosity_LastWins_InitSummary(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name        string
		args        []string
		wantSummary bool
	}{
		{"quiet-verbose", []string{"--quiet", "--verbose", "init"}, true},
		{"quiet-debug", []string{"--quiet", "--debug", "init"}, true},
		{"verbose-quiet", []string{"--verbose", "--quiet", "init"}, false},
		{"debug-quiet", []string{"--debug", "--quiet", "init"}, false},
		{"quiet-verbose-quiet", []string{"--quiet", "--verbose", "--quiet", "init"}, false},
		{"verbose-quiet-verbose", []string{"--verbose", "--quiet", "--verbose", "init"}, true},
		{"quiet-alone", []string{"--quiet", "init"}, false},
		{"verbose-alone", []string{"--verbose", "init"}, true},
		{"none", []string{"init"}, true},
		{"quiet-then-false", []string{"--quiet", "--quiet=false", "init"}, true},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			out, errOut, err := quietRun(t, newApp(initFixture(t), cli.WithGetwd(sandboxGetwd)), tc.args...)
			if err != nil {
				t.Fatalf("Execute %v: %v (a repeated or conflicting option is not rejected)", tc.args, err)
			}
			if got := strings.Contains(out, "Initialized u-boot project"); got != tc.wantSummary {
				t.Errorf("summary visible = %v, want %v:\nstdout=%q stderr=%q", got, tc.wantSummary, out, errOut)
			}
		})
	}
}

// Warning filter: the last option decides which warnings are printed.
func TestExecute_Verbosity_LastWins_WarningFilter(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		args     []string
		wantAll  bool
		wantSome bool
	}{
		{"quiet-verbose", []string{"--quiet", "--verbose", "init", "--devcontainer", "--sandbox"}, true, true},
		{"verbose-quiet", []string{"--verbose", "--quiet", "init", "--devcontainer", "--sandbox"}, false, true},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			uc := &fakeInitUseCase{resp: driving.InitProjectResponse{Warnings: quietAllWarnings()}}
			out, _, err := quietRun(t, newApp(uc, cli.WithGetwd(sandboxGetwd)), tc.args...)
			if err != nil {
				t.Fatalf("Execute: %v", err)
			}
			if got := strings.Contains(out, "has no effect"); got != tc.wantAll {
				t.Errorf("non-relaxation warning visible = %v, want %v:\n%s", got, tc.wantAll, out)
			}
			if got := strings.Contains(out, "security relaxation"); got != tc.wantSome {
				t.Errorf("security relaxation visible = %v, want %v:\n%s", got, tc.wantSome, out)
			}
		})
	}
}

// Doctor: the report is full unless --quiet is the last option; the
// error item stays in both cases.
func TestExecute_Verbosity_LastWins_DoctorReport(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		args     []string
		wantFull bool
	}{
		{"quiet-verbose", []string{"--quiet", "--verbose", "doctor"}, true},
		{"verbose-quiet", []string{"--verbose", "--quiet", "doctor"}, false},
		{"debug-quiet", []string{"--debug", "--quiet", "doctor"}, false},
		{"quiet-debug", []string{"--quiet", "--debug", "doctor"}, true},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			doctorUC := &fakeDoctorUseCase{resp: driving.DoctorResponse{
				Report: domain.DiagnosticReport{Items: []domain.Diagnostic{
					{ID: "fs.write-permissions", Severity: domain.SeverityOK, Message: "BaseDir is writable."},
					{ID: "docker.daemon", Severity: domain.SeverityError, Message: "daemon unreachable."},
				}},
			}}
			var stdout, stderr bytes.Buffer
			err := newAppWithDoctor(&fakeInitUseCase{}, doctorUC, cli.WithGetwd(sandboxGetwd)).Execute(
				context.Background(), tc.args, &stdout, &stderr)
			if !errors.Is(err, cli.ErrDoctorFailures) {
				t.Fatalf("err = %v, want ErrDoctorFailures", err)
			}
			out := stdout.String()
			if !strings.Contains(out, "daemon unreachable.") {
				t.Errorf("the error item is missing: %q", out)
			}
			if got := strings.Contains(out, "BaseDir is writable"); got != tc.wantFull {
				t.Errorf("OK item visible = %v, want %v: %q", got, tc.wantFull, out)
			}
		})
	}
}

// Down: the Compose stream is held back only while --quiet is the
// last option (buffer decision reads the same value).
func TestExecute_Verbosity_LastWins_DownComposeStream(t *testing.T) {
	t.Parallel()
	const compose = " Container demo-db-1  Stopping\n"
	for _, tc := range []struct {
		name       string
		args       []string
		wantStream bool
	}{
		{"quiet-verbose", []string{"--quiet", "--verbose", "down"}, true},
		{"verbose-quiet", []string{"--verbose", "--quiet", "down"}, false},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			uc := &fakeDownUseCase{compose: compose}
			var stdout, stderr bytes.Buffer
			if err := newAppWithDown(uc, cli.WithGetwd(sandboxGetwd)).Execute(
				context.Background(), tc.args, &stdout, &stderr); err != nil {
				t.Fatalf("Execute: %v", err)
			}
			if got := strings.Contains(stderr.String(), "Container demo-db-1"); got != tc.wantStream {
				t.Errorf("compose stream on stderr = %v, want %v:\n%s", got, tc.wantStream, stderr.String())
			}
		})
	}
}

// Negative: the order of the options does not change `--json`; the
// envelope is the same for every combination.
func TestExecute_Verbosity_LastWins_JSONUnchanged(t *testing.T) {
	t.Parallel()
	run := func(args ...string) string {
		out, _, err := quietRun(t, newApp(initFixture(t), cli.WithGetwd(sandboxGetwd)), args...)
		if err != nil {
			t.Fatalf("Execute %v: %v", args, err)
		}
		return out
	}
	want := run("--json", "init")
	for _, args := range [][]string{
		{"--json", "--quiet", "--verbose", "init"},
		{"--json", "--verbose", "--quiet", "init"},
		{"--quiet", "--json", "init"},
	} {
		if got := run(args...); got != want {
			t.Errorf("%v changed the --json output:\n%s\nwant:\n%s", args, got, want)
		}
	}
}

// Boundary: the verbosity state does not leak between two Execute
// calls on one App.
func TestExecute_Verbosity_StateIsPerExecute(t *testing.T) {
	t.Parallel()
	app := newApp(initFixture(t), cli.WithGetwd(sandboxGetwd))
	if out, _, err := quietRun(t, app, "--quiet", "init"); err != nil || out != "" {
		t.Fatalf("first run: out=%q err=%v, want silent success", out, err)
	}
	out, _, err := quietRun(t, app, "init")
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if !strings.Contains(out, "Initialized u-boot project") {
		t.Errorf("the second run inherited --quiet from the first: %q", out)
	}
}

// Negative: a value that is no boolean is a flag-parse error (exit
// code 2 class) — only a multiple or conflicting option is accepted.
func TestExecute_Verbosity_InvalidValueIsRejected(t *testing.T) {
	t.Parallel()
	_, _, err := quietRun(t, newApp(initFixture(t), cli.WithGetwd(sandboxGetwd)), "--quiet=maybe", "init")
	if err == nil {
		t.Fatal("--quiet=maybe was accepted")
	}
}
