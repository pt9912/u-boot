package cli_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli/jsontestutil"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// instanceOf returns the instance the use case received in the request of the
// subcommand, and how often it was called.
func instanceOf(uc *fakeDevcontainerUseCase, sub string) (instance string, calls int) {
	switch sub {
	case "up":
		return uc.upReq.Instance, uc.upCalls
	case "down":
		return uc.downReq.Instance, uc.downCalls
	}
	return uc.shellReq.Instance, uc.shellCalls
}

// LH-FA-DEV-010 (ADR-0038 decision 1): `--instance <name>` of `up`, `down` and
// `shell` reaches the request of the use case; without the flag the request
// carries no instance.
func TestExecute_Devcontainer_InstanceFlagReachesTheRequest(t *testing.T) {
	t.Parallel()
	for _, sub := range []string{"up", "down", "shell"} {
		for _, args := range [][]string{{"--instance", "agent-1"}, {"--instance=agent-1"}} {
			t.Run(sub+" "+strings.Join(args, " "), func(t *testing.T) {
				t.Parallel()
				uc := sessionUseCase()
				if _, _, err := shellRun(t, uc, append([]string{"devcontainer", sub}, args...)...); err != nil {
					t.Fatal(err)
				}
				if got, calls := instanceOf(uc, sub); got != "agent-1" || calls != 1 {
					t.Errorf("instance = %q after %d calls, want agent-1 after 1", got, calls)
				}
			})
		}
		t.Run(sub+" without the flag", func(t *testing.T) {
			t.Parallel()
			uc := sessionUseCase()
			if _, _, err := shellRun(t, uc, "devcontainer", sub); err != nil {
				t.Fatal(err)
			}
			if got, calls := instanceOf(uc, sub); got != "" || calls != 1 {
				t.Errorf("instance = %q after %d calls, want none after 1", got, calls)
			}
		})
	}
}

// ADR-0038 decision 6 (LH-FA-DEV-010): a `--instance` given empty is an
// invalid name (exit 10), not "no instance"; the use case is not called. With
// `--json`, `up` and `down` write the envelope of the error with the code of
// the requirement (`shell` rejects `--json` with exit 2 before).
func TestExecute_Devcontainer_EmptyInstanceIsInvalid(t *testing.T) {
	t.Parallel()
	for _, sub := range []string{"up", "down", "shell"} {
		for _, args := range [][]string{{"--instance", ""}, {"--instance="}} {
			t.Run(sub+" "+strings.Join(args, "|"), func(t *testing.T) {
				t.Parallel()
				uc := sessionUseCase()
				out, _, err := shellRun(t, uc, append([]string{"devcontainer", sub}, args...)...)
				if !errors.Is(err, driving.ErrDevcontainerInstanceInvalid) || cli.ExitCode(err) != 10 || out != "" {
					t.Errorf("err = %v, exit = %d, stdout = %q, want ErrDevcontainerInstanceInvalid, exit 10, no stdout", err, cli.ExitCode(err), out)
				}
				if _, calls := instanceOf(uc, sub); calls != 0 {
					t.Errorf("the use case was called %d times", calls)
				}
				if sub == "shell" {
					return
				}
				out, _, err = shellRun(t, sessionUseCase(), append([]string{"devcontainer", sub, "--json"}, args...)...)
				if cli.ExitCode(err) != 10 {
					t.Fatalf("json: exit = %d, want 10", cli.ExitCode(err))
				}
				jsontestutil.AssertMinimalEnvelope(t, []byte(out),
					jsontestutil.WithCommand("devcontainer"), jsontestutil.WithSubcommand(sub),
					jsontestutil.WithExitCode(10), jsontestutil.WithExpectedCodes("LH-FA-DEV-010"))
			})
		}
	}
}

// ADR-0038 decision 1: the helps name `--instance`, the rule of the name, the
// behaviour without the flag (only the container without an instance, exit 10
// when the project has instances only) and the exit-10 reasons; the group
// help points at the flag. The behaviour behind each sentence has its test
// (`TestDevcontainerInstance_*` of the application, the Docker test of the
// instances). No help carries a requirement, ADR or specification identifier.
func TestExecute_Devcontainer_HelpNamesInstance(t *testing.T) {
	t.Parallel()
	text := func(args ...string) string {
		out, errOut, err := quietRun(t, newAppWithDevcontainer(&fakeDevcontainerUseCase{}), args...)
		if err != nil {
			t.Fatalf("%v: %v", args, err)
		}
		return strings.Join(strings.Fields(out+errOut), " ")
	}
	rule := "1 to 32 characters of a-z, 0-9 and '-'"
	only := "if the project has only containers of instances, it ends with exit 10 and "
	cases := []struct {
		args  []string
		parts []string
	}{
		{[]string{"devcontainer", "--help"}, []string{"Each subcommand takes --instance NAME to act on one of several devcontainers of the same project folder"}},
		{[]string{"devcontainer", "up", "--help"}, []string{
			"--instance NAME starts a further devcontainer of the same project folder, next to the one without an instance;",
			rule, "Without --instance the command acts on the container without an instance only; " + only + "starts nothing.",
			"an instance name that is not admissible; a container-name part the engine would refuse for --name, or a container of another project folder that already holds the name; no --instance although the project has only containers of instances",
			"u-boot devcontainer up --instance agent-1",
		}},
		{[]string{"devcontainer", "down", "--help"}, []string{
			"--instance NAME stops the container of that instance", rule, "with no container of that instance the command ends with exit 0.",
			"Without --instance the command acts on the container without an instance only; " + only + "stops nothing.",
			"10 an instance name that is not admissible; no --instance although the project has only containers of instances",
			"u-boot devcontainer down --instance agent-1",
		}},
		{[]string{"devcontainer", "shell", "--help"}, []string{
			"--instance NAME opens the shell in the container of that instance", rule,
			"Without --instance the command acts on the container without an instance only; " + only + "starts nothing.",
			"Before the session begins, the command names the container of the session beside its ID on the standard error stream.",
			"an instance name that is not admissible; a container-name part the engine would refuse for --name, or a container of another project folder that already holds the name; no --instance although the project has only containers of instances",
			"u-boot devcontainer shell --instance agent-1",
		}},
	}
	for _, tc := range cases {
		got := text(tc.args...)
		for _, id := range []string{"LH-FA-DEV", "ADR-", "SPEC-", "ARC-"} {
			if strings.Contains(got, id) {
				t.Errorf("%v: the help carries %q:\n%s", tc.args, id, got)
			}
		}
		for _, part := range tc.parts {
			if !strings.Contains(got, part) {
				t.Errorf("%v: the help lacks %q:\n%s", tc.args, part, got)
			}
		}
	}
}
