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

// fakeEgressListUseCase records the request and returns the configured
// response or error.
type fakeEgressListUseCase struct {
	called  int
	lastReq driving.EgressListRequest
	resp    driving.EgressListResponse
	err     error
}

func (f *fakeEgressListUseCase) List(_ context.Context, req driving.EgressListRequest) (driving.EgressListResponse, error) {
	f.called++
	f.lastReq = req
	return f.resp, f.err
}

func newAppWithEgress(uc driving.EgressListUseCase, opts ...cli.Option) *cli.App {
	return newApp(&fakeInitUseCase{}, append([]cli.Option{cli.WithEgressListUseCase(uc), cli.WithGetwd(sandboxGetwd)}, opts...)...)
}

func activeEgressResponse() driving.EgressListResponse {
	return driving.EgressListResponse{
		SandboxProfile: true, EgressEnabled: true, CloneSource: true,
		Hosts: []driving.EgressHostEntry{
			{Host: "api.github.com", Sources: []string{"base"}},
			{Host: "ghcr.io", Sources: []string{"runtime", "feature:homebrew"}},
			{Host: "gitlab.example.org", Sources: []string{"clone"}},
			{Host: "registry.npmjs.org", Sources: []string{"feature:claude-code", "feature:node"}},
		},
	}
}

func runEgress(t *testing.T, app *cli.App, args ...string) (string, string, error) {
	t.Helper()
	return quietRun(t, app, args...)
}

// egressRows maps each host line of the table to its normalised origin text.
func egressRows(out string) map[string]string {
	rows := map[string]string{}
	for _, line := range strings.Split(out, "\n") {
		f := strings.Fields(line)
		if len(f) >= 2 && !strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "  ") &&
			f[0] != "HOST" && f[0] != "Egress" {
			rows[f[0]] = strings.Join(f[1:], " ")
		}
	}
	return rows
}

// LH-FA-DEV-008: the table shows each host with all its origins, the state
// line says the restriction is active, and the use case gets the working
// directory.
func TestExecute_EgressList_ActiveTable(t *testing.T) {
	t.Parallel()
	uc := &fakeEgressListUseCase{resp: activeEgressResponse()}
	out, errOut, err := runEgress(t, newAppWithEgress(uc), "egress", "list")
	if err != nil || errOut != "" {
		t.Fatalf("err=%v stderr=%q", err, errOut)
	}
	if uc.called != 1 || uc.lastReq.BaseDir != "/tmp/x/demo" {
		t.Errorf("use case called %d times with %+v, want once with the working directory", uc.called, uc.lastReq)
	}
	if !strings.HasPrefix(out, "Egress restriction: active.\n") {
		t.Errorf("state line missing:\n%s", out)
	}
	want := map[string]string{
		"api.github.com":     "base",
		"ghcr.io":            "runtime, feature:homebrew",
		"gitlab.example.org": "clone",
		"registry.npmjs.org": "feature:claude-code, feature:node",
	}
	got := egressRows(out)
	for host, origins := range want {
		if got[host] != origins {
			t.Errorf("row %s = %q, want %q\n%s", host, got[host], origins, out)
		}
	}
	if strings.Contains(out, "No clone host") {
		t.Errorf("a clone source exists, the note about a missing one must not appear:\n%s", out)
	}
}

// LH-FA-DEV-008: when the restriction is not active the output says so, names
// the reason and still shows the list "that would apply on activation".
func TestExecute_EgressList_NotActiveStatesTheReason(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		profile, enabled bool
		reason           string
	}{
		"neither":            {false, false, "devcontainer.profile is not `sandbox` and devcontainer.sandbox.egress.enabled is not true"},
		"egress, no sandbox": {false, true, "(devcontainer.profile is not `sandbox`)"},
		"sandbox, no egress": {true, false, "(devcontainer.sandbox.egress.enabled is not true)"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			resp := activeEgressResponse()
			resp.SandboxProfile, resp.EgressEnabled = tc.profile, tc.enabled
			out, _, err := runEgress(t, newAppWithEgress(&fakeEgressListUseCase{resp: resp}), "egress", "list")
			if err != nil {
				t.Fatal(err)
			}
			first, _, _ := strings.Cut(out, "\n")
			if !strings.HasPrefix(first, "Egress restriction: not active (") || !strings.Contains(first, tc.reason) ||
				!strings.Contains(first, "would apply once it is active") {
				t.Errorf("state line = %q", first)
			}
			if egressRows(out)["ghcr.io"] == "" {
				t.Errorf("the list must still be shown:\n%s", out)
			}
		})
	}
}

// LH-FA-DEV-008: the notes of the output — subdomains, a missing clone host and
// the delimitation from a running container (measured statement only).
func TestExecute_EgressList_Notes(t *testing.T) {
	t.Parallel()
	resp := activeEgressResponse()
	resp.CloneSource = false
	out, _, err := runEgress(t, newAppWithEgress(&fakeEgressListUseCase{resp: resp}), "egress", "list")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"A name in the list also allows its subdomains.",
		"No clone host: there is no git remote 'origin' in .git/config and no devcontainer.sandbox.repository.",
		"It is not read from a container: a container keeps the list of the image it was built from",
		"stopping and starting it again kept the old list, `devcontainer up --remove-existing-container` applied the new one;",
		"measured on Colima, arm64; editor actions not measured).",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
}

// LH-FA-CLI-005: the output of the query command stays under `--quiet`, byte
// for byte and with every line (state line and notes included).
func TestExecute_EgressList_QuietKeepsTheCompleteOutput(t *testing.T) {
	t.Parallel()
	for name, resp := range map[string]driving.EgressListResponse{
		"active":     activeEgressResponse(),
		"not active": {Hosts: activeEgressResponse().Hosts},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			plain, _, err := runEgress(t, newAppWithEgress(&fakeEgressListUseCase{resp: resp}), "egress", "list")
			if err != nil {
				t.Fatal(err)
			}
			quiet, qErr, err := runEgress(t, newAppWithEgress(&fakeEgressListUseCase{resp: resp}), "--quiet", "egress", "list")
			if err != nil || qErr != "" {
				t.Fatalf("--quiet: err=%v stderr=%q", err, qErr)
			}
			if quiet != plain || !strings.Contains(quiet, "Egress restriction: ") || !strings.Contains(quiet, "Notes:") {
				t.Errorf("--quiet output differs from the default or lost lines:\n quiet=%q\n plain=%q", quiet, plain)
			}
			jsonPlain, _, _ := runEgress(t, newAppWithEgress(&fakeEgressListUseCase{resp: resp}), "egress", "list", "--json")
			jsonQuiet, _, _ := runEgress(t, newAppWithEgress(&fakeEgressListUseCase{resp: resp}), "--quiet", "egress", "list", "--json")
			if jsonPlain != jsonQuiet {
				t.Errorf("--quiet --json differs from --json:\n%s\n%s", jsonQuiet, jsonPlain)
			}
		})
	}
}

// LH-NFA-USE-004: `--json` emits the envelope with command `egress`, subcommand
// `list` and `data` carrying the state and the hosts with their origins; both
// flag positions work.
func TestExecute_EgressList_JSON(t *testing.T) {
	t.Parallel()
	uc := &fakeEgressListUseCase{resp: activeEgressResponse()}
	out, _, err := runEgress(t, newAppWithEgress(uc), "egress", "list", "--json")
	if err != nil {
		t.Fatal(err)
	}
	jsontestutil.AssertMinimalEnvelope(t, []byte(out),
		jsontestutil.WithCommand("egress"), jsontestutil.WithSubcommand("list"), jsontestutil.WithExitCode(0))
	var env struct {
		Data struct {
			Active         *bool `json:"active"`
			SandboxProfile *bool `json:"sandboxProfile"`
			EgressEnabled  *bool `json:"egressEnabled"`
			CloneSource    *bool `json:"cloneSource"`
			Hosts          []struct {
				Host    string   `json:"host"`
				Sources []string `json:"sources"`
			} `json:"hosts"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	d := env.Data
	if d.Active == nil || !*d.Active || d.SandboxProfile == nil || !*d.SandboxProfile ||
		d.EgressEnabled == nil || !*d.EgressEnabled || d.CloneSource == nil || !*d.CloneSource {
		t.Errorf("state fields wrong: %+v\n%s", d, out)
	}
	if len(d.Hosts) != 4 || d.Hosts[1].Host != "ghcr.io" || strings.Join(d.Hosts[1].Sources, ",") != "runtime,feature:homebrew" {
		t.Errorf("hosts = %+v", d.Hosts)
	}
	flagFirst, _, _ := runEgress(t, newAppWithEgress(&fakeEgressListUseCase{resp: activeEgressResponse()}), "--json", "egress", "list")
	if flagFirst != out {
		t.Errorf("--json before the subcommand differs:\n%s\n%s", flagFirst, out)
	}
	// An inactive state is data, not a diagnostic: no warn, exit 0.
	inactive, _, err := runEgress(t, newAppWithEgress(&fakeEgressListUseCase{resp: driving.EgressListResponse{}}), "egress", "list", "--json")
	if err != nil || !strings.Contains(inactive, `"active":false`) || !strings.Contains(inactive, `"hosts":[]`) || !strings.Contains(inactive, `"status":"ok"`) {
		t.Errorf("inactive/empty: err=%v\n%s", err, inactive)
	}
}

// LH-FA-CLI-006: the error classes of the neighbours — exit code and
// diagnostic code in the envelope.
func TestExecute_EgressList_Errors(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		err  error
		exit int
		code string
	}{
		{"no project", fmt.Errorf("%w: /tmp/x/demo/u-boot.yaml not found", driving.ErrProjectNotInitialized), 10, "LH-FA-INIT-001"},
		{"invalid configuration", fmt.Errorf("%w: parse", driving.ErrConfigSchemaInvalid), 10, "LH-FA-CONF-002"},
		{"unsafe origin", fmt.Errorf("%w: origin embeds credentials", domain.ErrInvalidSandboxSetting), 10, "LH-FA-DEV-006"},
		{"read failure", fmt.Errorf("%w: read /tmp/x/demo/u-boot.yaml", driving.ErrEgressFileSystem), 14, "LH-NFA-REL-003"},
		{"other", errors.New("boom"), 1, "LH-FA-CLI-006"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			out, _, err := runEgress(t, newAppWithEgress(&fakeEgressListUseCase{err: tc.err}), "egress", "list")
			if !errors.Is(err, tc.err) || cli.ExitCode(err) != tc.exit || out != "" {
				t.Fatalf("human: err=%v exit=%d stdout=%q, want %v exit %d and no stdout", err, cli.ExitCode(err), out, tc.err, tc.exit)
			}
			out, _, err = runEgress(t, newAppWithEgress(&fakeEgressListUseCase{err: tc.err}), "egress", "list", "--json")
			if cli.ExitCode(err) != tc.exit {
				t.Fatalf("json: exit=%d, want %d", cli.ExitCode(err), tc.exit)
			}
			jsontestutil.AssertMinimalEnvelope(t, []byte(out),
				jsontestutil.WithCommand("egress"), jsontestutil.WithSubcommand("list"),
				jsontestutil.WithExitCode(tc.exit), jsontestutil.WithExpectedCodes(tc.code))
			if strings.Contains(out, "/tmp/x/demo") {
				t.Errorf("the project path leaks into the diagnostic:\n%s", out)
			}
		})
	}
}

// `egress list` is a read-only query: `--dry-run` and `--diff` are rejected with
// exit 2 (an envelope under `--json`) before the use case runs; extra
// arguments are a usage error.
func TestExecute_EgressList_RejectsPreviewFlagsAndArguments(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{{"egress", "list", "--dry-run"}, {"egress", "list", "--diff"}} {
		uc := &fakeEgressListUseCase{resp: activeEgressResponse()}
		out, _, err := runEgress(t, newAppWithEgress(uc), args...)
		if !errors.Is(err, cli.ErrEgressPreviewNotApplicable) || cli.ExitCode(err) != 2 || uc.called != 0 || out != "" {
			t.Errorf("%v: err=%v exit=%d calls=%d stdout=%q", args, err, cli.ExitCode(err), uc.called, out)
		}
		out, _, err = runEgress(t, newAppWithEgress(uc), append(args, "--json")...)
		if cli.ExitCode(err) != 2 {
			t.Errorf("%v --json: exit=%d, want 2", args, cli.ExitCode(err))
		}
		jsontestutil.AssertMinimalEnvelope(t, []byte(out),
			jsontestutil.WithCommand("egress"), jsontestutil.WithSubcommand("list"),
			jsontestutil.WithExitCode(2), jsontestutil.WithExpectedCodes("LH-FA-CLI-006"))
	}
	uc := &fakeEgressListUseCase{}
	out, _, err := runEgress(t, newAppWithEgress(uc), "egress", "list", "extra", "--json")
	if err == nil || cli.ExitCode(err) != 2 || uc.called != 0 {
		t.Errorf("extra argument: err=%v exit=%d calls=%d", err, cli.ExitCode(err), uc.called)
	}
	jsontestutil.AssertMinimalEnvelope(t, []byte(out),
		jsontestutil.WithCommand("egress"), jsontestutil.WithSubcommand("list"), jsontestutil.WithExitCode(2))
}

// Bare `egress` prints the help (human) and rejects `--json` with exit 2 and
// without an envelope, like bare `template`.
func TestExecute_Egress_Bare(t *testing.T) {
	t.Parallel()
	out, _, err := runEgress(t, newAppWithEgress(&fakeEgressListUseCase{}), "egress")
	if err != nil || !strings.Contains(out, "u-boot egress list") {
		t.Errorf("human: err=%v out=%q, want help naming `u-boot egress list`", err, out)
	}
	out, _, err = runEgress(t, newAppWithEgress(&fakeEgressListUseCase{}), "--json", "egress")
	if !errors.Is(err, cli.ErrEgressSubcommandRequired) || cli.ExitCode(err) != 2 || out != "" {
		t.Errorf("--json: err=%v exit=%d out=%q", err, cli.ExitCode(err), out)
	}
}

// A command that was not wired fails with an error instead of a panic.
func TestExecute_EgressList_NotWired(t *testing.T) {
	t.Parallel()
	_, _, err := runEgress(t, newApp(&fakeInitUseCase{}, cli.WithGetwd(sandboxGetwd)), "egress", "list")
	if err == nil || !strings.Contains(err.Error(), "not wired") {
		t.Errorf("err = %v, want a `not wired` error", err)
	}
}

// A failing working-directory probe is an error, no panic.
func TestExecute_EgressList_GetwdFails(t *testing.T) {
	t.Parallel()
	uc := &fakeEgressListUseCase{}
	app := newAppWithEgress(uc, cli.WithGetwd(func() (string, error) { return "", errors.New("no cwd") }))
	_, _, err := runEgress(t, app, "egress", "list")
	if err == nil || !strings.Contains(err.Error(), "determine working directory") || uc.called != 0 {
		t.Errorf("err=%v calls=%d", err, uc.called)
	}
}

// LH-FA-CLI-002: `u-boot --help` lists the command, `egress list --help` names
// its purpose and carries an example.
func TestExecute_Egress_Help(t *testing.T) {
	t.Parallel()
	var root bytes.Buffer
	if err := newApp(&fakeInitUseCase{}).Execute(context.Background(), []string{"--help"}, &root, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(root.String(), "egress") || !strings.Contains(root.String(), "Inspect the egress allowlist of the sandbox devcontainer") {
		t.Errorf("root help lacks the egress command and its short description:\n%s", root.String())
	}
	group, _, err := runEgress(t, newAppWithEgress(&fakeEgressListUseCase{}), "egress", "--help")
	if err != nil || !strings.Contains(group, "list") || !strings.Contains(group, "Show the egress allowlist with the origin of each host") {
		t.Errorf("egress --help lacks the short description of `list`: err=%v\n%s", err, group)
	}
	out, _, err := runEgress(t, newAppWithEgress(&fakeEgressListUseCase{}), "egress", "list", "--help")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Show the egress allowlist that \"u-boot generate devcontainer\" writes into the",
		"Examples:\n  u-boot egress list\n  u-boot egress list --json",
		"It is not read from a container",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("egress list --help lacks %q:\n%s", want, out)
		}
	}
}
