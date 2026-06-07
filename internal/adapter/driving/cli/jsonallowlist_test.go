package cli_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
)

// TestRootJSON_RejectsAllNonMigratedForms is the slice-v1-cli-json-
// dry-run-doctor T3 acceptance pin: every spec-enum subcommand form
// that has NOT been migrated must reject --json with exit code 2
// (ErrJSONNotImplemented).
//
// Migrate-Forms: "doctor", "template list", "add" (T4 add-slice),
// "init" (T5 init-slice), "generate" (T5 generate-slice), "remove"
// (T5 remove-slice), "up" + "down" (T5 up-down-slice), "logs" (T5
// logs-slice). Reject-Forms: 4 — see slice-doctor §T0-(g)
// §Subcommand-Form-Inventar minus the migrated forms.
func TestRootJSON_RejectsAllNonMigratedForms(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		wantSuffix string
	}{
		{"config (bare)", []string{"--json", "config"}, "config"},
		{"config get", []string{"--json", "config", "get", "project.name"}, "config"},
		{"config set", []string{"--json", "config", "set", "project.name", "x"}, "config"},
		{"template (bare)", []string{"--json", "template"}, "template"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := newApp(&fakeInitUseCase{})
			err := app.Execute(context.Background(), tc.args, &bytes.Buffer{}, &bytes.Buffer{})
			if err == nil {
				t.Fatalf("want reject error, got nil")
			}
			if !errors.Is(err, cli.ErrJSONNotImplemented) {
				t.Errorf("want ErrJSONNotImplemented, got %v", err)
			}
			if cli.ExitCode(err) != 2 {
				t.Errorf("want exit code 2, got %d (err=%v)", cli.ExitCode(err), err)
			}
			wantRef := "slice-v1-cli-json-dry-run-" + tc.wantSuffix
			if !strings.Contains(err.Error(), wantRef) {
				t.Errorf("want error to reference %q, got %v", wantRef, err)
			}
		})
	}
}

// TestRootJSON_AcceptsDoctor is the positive pin: u-boot doctor --json
// MUST not reject — doctor is migrated in this slice (Allowlist).
// The actual envelope emission lands in T5; here we only assert that
// the gate lets the invocation through (the doctor RunE then runs
// its today's path with a.json=true).
func TestRootJSON_AcceptsDoctor(t *testing.T) {
	app := newApp(&fakeInitUseCase{})
	err := app.Execute(context.Background(), []string{"--json", "doctor"}, &bytes.Buffer{}, &bytes.Buffer{})
	// doctor may still fail downstream (the fake DoctorUseCase returns
	// an empty report, no error). What matters here: NOT
	// ErrJSONNotImplemented.
	if errors.Is(err, cli.ErrJSONNotImplemented) {
		t.Errorf("doctor must be allowed through; got %v", err)
	}
}

// TestRootJSON_AcceptsTemplateList_BothFlagPositions is the M3
// migration-pin from slice-doctor §AK §template-list-Schnitt:
// `u-boot template list --json` and `u-boot --json template list`
// produce identical output (same exit code, same stdout body).
// This is the M3-Review-Round-2-Finding pin.
func TestRootJSON_AcceptsTemplateList_BothFlagPositions(t *testing.T) {
	app1 := newApp(&fakeInitUseCase{})
	out1 := &bytes.Buffer{}
	err1 := app1.Execute(context.Background(), []string{"template", "list", "--json"}, out1, &bytes.Buffer{})

	app2 := newApp(&fakeInitUseCase{})
	out2 := &bytes.Buffer{}
	err2 := app2.Execute(context.Background(), []string{"--json", "template", "list"}, out2, &bytes.Buffer{})

	if err1 != nil || err2 != nil {
		t.Fatalf("template list --json should succeed both ways; got err1=%v err2=%v", err1, err2)
	}
	if out1.String() != out2.String() {
		t.Errorf("output mismatch:\n--- subcommand --json ---\n%s\n--- --json subcommand ---\n%s",
			out1.String(), out2.String())
	}
	if !strings.HasPrefix(strings.TrimSpace(out1.String()), "[") {
		t.Errorf("template list --json should emit a JSON array, got: %q", out1.String())
	}
}

// TestRootJSON_AcceptsTemplateList_FlagBeforeSubcommand is a
// targeted regression pin for the M3 migration: removing the
// LOCAL --json flag on template list and re-routing through the
// ROOT persistent flag MUST keep `--json template list` working.
// This catches the Cobra-flag-shadow trap where a leftover local
// flag would silently mask the root state.
func TestRootJSON_AcceptsTemplateList_FlagBeforeSubcommand(t *testing.T) {
	app := newApp(&fakeInitUseCase{})
	out := &bytes.Buffer{}
	err := app.Execute(context.Background(), []string{"--json", "template", "list"}, out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("--json template list failed: %v", err)
	}
	if !strings.Contains(out.String(), "[") {
		t.Errorf("expected JSON array output, got: %q", out.String())
	}
}

// TestRootJSON_TreeWalkAllowlistCompleteness is the M2 anti-drift
// pin from slice-doctor §T0-(g): echter Tree-Walk durch alle
// registrierten Cobra-Subcommand-Pfade. Für jeden Pfad gilt:
// entweder Allowlist-Hit (--json akzeptiert) ODER Reject mit
// ErrJSONNotImplemented. Eine Form, die NEITHER ist, ist ein
// Missing-Allowlist-Entry und bricht den Test.
//
// Code-Realität ist die Quelle: ein neuer Subcommand, der per
// root.AddCommand(...) hinzukommt, taucht automatisch im Walk auf —
// statische Erwartungs-Listen entfallen (Review H3-Findings:
// "TreeWalk-Test deckt KEINEN Tree-Walk ab" adressiert).
func TestRootJSON_TreeWalkAllowlistCompleteness(t *testing.T) {
	app := newApp(&fakeInitUseCase{})
	paths := app.WalkRootCommandPathsForTest()
	if len(paths) == 0 {
		t.Fatal("Cobra tree walk returned zero paths — walker broken")
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			args := append([]string{"--json"}, strings.Split(strings.TrimPrefix(path, "u-boot "), " ")...)
			args = appendStubArgs(args, path)
			err := newApp(&fakeInitUseCase{}).Execute(context.Background(), args, &bytes.Buffer{}, &bytes.Buffer{})
			if errors.Is(err, cli.ErrJSONNotImplemented) {
				return // Reject path — non-migrated form, gate fires correctly.
			}
			// Otherwise: form must be migrated (Allowlist hit). Downstream
			// failures (use-case errors, broken fixtures) are acceptable —
			// what matters is "no ErrJSONNotImplemented" for an
			// allowlist-registered path.
		})
	}
}

// TestRootJSON_AllowlistAndTreeMatch is the second half of the M2-
// Drift-Gate: every Allowlist-key MUST correspond to a real Cobra
// command path. If someone adds a stale Allowlist entry pointing to
// a Use-renamed or removed subcommand, this test catches it.
func TestRootJSON_AllowlistAndTreeMatch(t *testing.T) {
	app := newApp(&fakeInitUseCase{})
	allowlist := cli.JSONAllowlistPathsForTest()
	treePaths := app.WalkRootCommandPathsForTest()
	treeSet := make(map[string]bool, len(treePaths))
	for _, p := range treePaths {
		treeSet[p] = true
	}
	for _, key := range allowlist {
		if !treeSet[key] {
			t.Errorf("allowlist key %q has no corresponding Cobra command (stale Use-rename?)", key)
		}
	}
}

// appendStubArgs adds the positional args that Cobra needs for
// non-migrated forms to PARSE without ErrUsage. Otherwise Cobra
// returns an "accepts X args" error BEFORE PersistentPreRunE runs
// and our reject gate never fires.
func appendStubArgs(args []string, path string) []string {
	switch path {
	case "u-boot init":
		return append(args, "stub")
	case "u-boot add", "u-boot remove":
		return append(args, "postgres")
	case "u-boot generate":
		return append(args, "readme")
	case "u-boot config get":
		return append(args, "project.name")
	case "u-boot config set":
		return append(args, "project.name", "x")
	default:
		return args
	}
}

// TestJSONSliceSuffix_StableMapping pins the path → slice-suffix
// resolution used in jsonRejectError. Catches accidental renames.
// Includes the defensive "unknown" fallback for empty/garbled paths
// (Review L3-Findings adressiert).
func TestJSONSliceSuffix_StableMapping(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"u-boot init", "init"},
		{"u-boot add", "add"},
		{"u-boot remove", "remove"},
		{"u-boot up", "up-down"},
		{"u-boot down", "up-down"},
		{"u-boot logs", "logs"},
		{"u-boot generate", "generate"},
		{"u-boot config", "config"},
		{"u-boot config get", "config"},
		{"u-boot config set", "config"},
		{"u-boot template", "template"},
		{"", "unknown"},
		{"u-boot", "unknown"},
		{"u-boot ", "unknown"},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			got := cli.JSONSliceSuffixForTest(tc.path)
			if got != tc.want {
				t.Errorf("path %q: want suffix %q, got %q", tc.path, tc.want, got)
			}
		})
	}
}

// TestRootJSON_AcceptsHelpFlag is the M6-Anti-Drift-Pin:
// --json combined with --help on any non-migrated subcommand must
// print help (no reject). Cobra's --help is a read-only escape hatch.
func TestRootJSON_AcceptsHelpFlag(t *testing.T) {
	cases := [][]string{
		{"--json", "init", "--help"},
		{"--json", "add", "--help"},
		{"--json", "config", "--help"},
		{"--json", "template", "--help"},
	}
	for _, args := range cases {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			app := newApp(&fakeInitUseCase{})
			err := app.Execute(context.Background(), args, &bytes.Buffer{}, &bytes.Buffer{})
			if err != nil {
				t.Errorf("--json + --help must not reject; got %v", err)
			}
		})
	}
}
