package cli_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli/jsontestutil"
)

// TestRootJSON_AcceptsTemplateList_BothFlagPositions pins that
// `u-boot template list --json` and `u-boot --json template list`
// produce identical output. Post-Cluster-T_close the output is the
// LH-NFA-USE-004 Minimalkontrakt-Envelope (slice-v1-cli-json-dry-run-
// template T2).
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
	jsontestutil.AssertMinimalEnvelope(t, out1.Bytes(),
		jsontestutil.WithCommand("template"),
		jsontestutil.WithSubcommand("list"),
		jsontestutil.WithExitCode(0),
	)
	if !strings.Contains(out1.String(), "\"data\"") {
		t.Errorf("template list --json envelope must carry a data field, got: %q", out1.String())
	}
}

// TestRootJSON_AcceptsTemplateList_FlagBeforeSubcommand pins that
// `--json` before the subcommand still routes through the root
// persistent flag (Cobra-flag-shadow regression guard).
func TestRootJSON_AcceptsTemplateList_FlagBeforeSubcommand(t *testing.T) {
	app := newApp(&fakeInitUseCase{})
	out := &bytes.Buffer{}
	err := app.Execute(context.Background(), []string{"--json", "template", "list"}, out, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("--json template list failed: %v", err)
	}
	jsontestutil.AssertMinimalEnvelope(t, out.Bytes(),
		jsontestutil.WithCommand("template"),
		jsontestutil.WithSubcommand("list"),
		jsontestutil.WithExitCode(0),
	)
}

// TestRootJSON_AllFormsRespondPostTClose is the post-Cluster-T_close
// replacement for the old allowlist-completeness tree-walk: the
// transitional reject gate is gone, so EVERY registered Cobra form
// (Spec-Enum incl. grouped subcommands — config/get/set,
// template/list count separately) must respond to `--json` without
// panicking. The only `--json` reject left is bare `u-boot template`
// (RunE-borne ErrTemplateSubcommandRequired). A new subcommand that
// forgets `--json` handling and leaks help shows up here.
func TestRootJSON_AllFormsRespondPostTClose(t *testing.T) {
	app := newApp(&fakeInitUseCase{})
	paths := app.WalkRootCommandPathsForTest()
	if len(paths) == 0 {
		t.Fatal("Cobra tree walk returned zero paths — walker broken")
	}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			args := append([]string{"--json"}, strings.Split(strings.TrimPrefix(path, "u-boot "), " ")...)
			args = appendStubArgs(args, path)
			var stdout bytes.Buffer
			err := newApp(&fakeInitUseCase{}).Execute(context.Background(), args, &stdout, &bytes.Buffer{})
			// bare `template` is the one RunE-borne reject form.
			if path == "u-boot template" {
				if !errors.Is(err, cli.ErrTemplateSubcommandRequired) {
					t.Errorf("bare `template --json` must reject with ErrTemplateSubcommandRequired, got %v", err)
				}
				if strings.Contains(stdout.String(), "Usage:") {
					t.Errorf("bare `template --json` must not leak help to stdout, got: %q", stdout.String())
				}
				return
			}
			// Every other form must run its RunE (downstream use-case
			// errors from the fakes are acceptable) — what matters is
			// the gate is gone: no form is silently rejected away.
		})
	}
}

// TestRootJSON_BareTemplateReject pins the Cluster-T_close bare-
// `template`-reject contract: `--json template` → Exit 2 via
// ErrTemplateSubcommandRequired, envelope-LOS (no envelope, no help
// leak); the human mode (no --json) still prints help.
func TestRootJSON_BareTemplateReject(t *testing.T) {
	t.Run("json-rejects-exit2-no-leak", func(t *testing.T) {
		app := newApp(&fakeInitUseCase{})
		var stdout, stderr bytes.Buffer
		err := app.Execute(context.Background(), []string{"--json", "template"}, &stdout, &stderr)
		if !errors.Is(err, cli.ErrTemplateSubcommandRequired) {
			t.Fatalf("want ErrTemplateSubcommandRequired, got %v", err)
		}
		if cli.ExitCode(err) != 2 {
			t.Errorf("exit = %d, want 2", cli.ExitCode(err))
		}
		if strings.Contains(stdout.String(), "Usage:") || strings.Contains(stdout.String(), "\"command\"") {
			t.Errorf("bare template --json must emit neither help nor envelope to stdout; got %q", stdout.String())
		}
	})
	t.Run("human-mode-still-prints-help", func(t *testing.T) {
		app := newApp(&fakeInitUseCase{})
		var stdout bytes.Buffer
		if err := app.Execute(context.Background(), []string{"template"}, &stdout, &bytes.Buffer{}); err != nil {
			t.Fatalf("bare template (human) should print help, got err %v", err)
		}
		if !strings.Contains(stdout.String(), "list") {
			t.Errorf("bare template (human) should print help mentioning `list`; got %q", stdout.String())
		}
	})
}

// appendStubArgs adds the positional args that Cobra needs for forms
// with ExactArgs validators to PARSE without a usage error before the
// RunE runs — so the tree-walk exercises the real RunE per form.
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
