//go:build docker

// End-to-end check of the host git resolution (LH-FA-DEV-009,
// LH-FA-DEV-010): the built binary runs against the engine of the runner,
// and the fake sources are scripts on the PATH of the binary — a gh that
// answers like the not-logged-in case or with a placeholder token, a fake
// git that answers the two identity keys. The runner carries no gh session
// and no host identity, so the tests prove the way with the fakes; a real
// gh of the runner would also be exercised (the switch and the order are
// the contract, not the token of a particular account).
//
// The value never appears in a command line, in a file of the project, in
// the metadata of the container or in an output of u-boot: the tests search
// every output and the inspect documents for the placeholder.

package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	// hostGitPlaceholders are the values of the fake sources; their length
	// (33 + the newline of printenv = 34) is the expected count in the session.
	hostGitTokenPlaceholder  = "platzh-hostgit-token-0123456789ab"
	hostGitNamePlaceholder   = "platzh-hostgit-name"
	hostGitEmailPlaceholder  = "platzh-hostgit-email"
	hostGitSessionLen        = "34"
	hostGitIdentityNameKey   = "user.name"
	hostGitIdentityEmailKey  = "user.email"
)

// hostGitEnvelope is the envelope of `devcontainer up --json` as far as the
// tests read it.
type hostGitEnvelope struct {
	Data struct {
		ContainerID           string `json:"containerId"`
		GitTokenSource        string `json:"gitTokenSource"`
		RemoteUser            string `json:"remoteUser"`
		RemoteWorkspaceFolder string `json:"remoteWorkspaceFolder"`
	} `json:"data"`
	Diagnostics []struct {
		Level   string `json:"level"`
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"diagnostics"`
}

// hostGitRun runs the binary with a PATH that leads to dir first (the fake
// sources) and stdin; gitToken, when set, is set for the process.
func hostGitRun(ctx context.Context, t *testing.T, bin, dir, stdin string, envPath string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "PATH="+envPath)
	cmd.Env = append(cmd.Env, "UBOOT_NO_COLOR=1")
	cmd.Stdin = strings.NewReader(stdin)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	code := 0
	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if !errors.As(err, &ee) {
			t.Fatalf("run %v: %v", args, err)
		}
		code = ee.ExitCode()
	}
	return stdout.String(), stderr.String(), code
}

// hostGitFakeWrites writes one executable script into its own directory and
// returns the directory.
func hostGitFakeWrites(t *testing.T, name, body string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// hostGitPrepare generates the project with the settings.
func hostGitPrepare(ctx context.Context, t *testing.T, sandbox bool, hostGitToken bool) string {
	t.Helper()
	settings := map[string]string{}
	if hostGitToken {
		settings["devcontainer.sandbox.hostGitToken"] = "true"
	}
	return dcGenerateProject(ctx, t, sandbox, settings)
}

// hostGitUnsetToken removes GIT_TOKEN from the test process for the test and
// restores it afterwards.
func hostGitUnsetToken(t *testing.T) {
	t.Helper()
	old, had := os.LookupEnv("GIT_TOKEN")
	if had {
		if err := os.Unsetenv("GIT_TOKEN"); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		if had {
			_ = os.Setenv("GIT_TOKEN", old)
		}
	})
}

// hostGitEnvelopes parses the up --json output.
func hostGitEnvelopes(t *testing.T, out string) hostGitEnvelope {
	t.Helper()
	var doc hostGitEnvelope
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("envelope: %v (%.600s)", err, out)
	}
	return doc
}

// hostGitHasWarning reports whether one diagnostic carries the text.
func hostGitHasWarning(env hostGitEnvelope, fragment string) bool {
	for _, d := range env.Diagnostics {
		if strings.Contains(d.Message, fragment) {
			return true
		}
	}
	return false
}

// TestE2E_HostGitToken_OptInReachesTheSession proves the opt-in way of the
// token with the fake gh: the resolution names the source, the container
// metadata carries no value, and the session and its credential helper carry
// the placeholder.
func TestE2E_HostGitToken_OptInReachesTheSession(t *testing.T) {
	hostGitUnsetToken(t)
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := hostGitPrepare(ctx, t, true, true)
	bin := ubootBinary(t)

	// A fake gh: answers the token of the session with the placeholder.
	ghDir := hostGitFakeWrites(t, "gh", "echo "+hostGitTokenPlaceholder)
	envPath := ghDir + ":" + os.Getenv("PATH")
	stdout, stderr, code := hostGitRun(ctx, t, bin, dir, "", envPath, "devcontainer", "up", "--json")
	if code != 0 {
		t.Fatalf("up: exit %d\nstdout %.600s\nstderr %.600s", code, stdout, stderr)
	}
	env := hostGitEnvelopes(t, stdout)
	if env.Data.GitTokenSource != "gh auth token" {
		t.Errorf("gitTokenSource = %q, want the gh source (%.600s)", env.Data.GitTokenSource, stdout)
	}
	if hostGitHasWarning(env, hostGitTokenPlaceholder) {
		t.Errorf("the envelope carries the value")
	}
	id := env.Data.ContainerID
	t.Cleanup(func() { _, _ = dcDocker("rm", "-f", "-v", id) })

	// The metadata of the devcontainer carries no value.
	doc := dcInspect(t, id)
	for _, v := range doc.Config.Env {
		if strings.Contains(v, hostGitTokenPlaceholder) {
			t.Errorf("Config.Env holds the token: %q", v)
		}
	}
	if strings.Contains(doc.Config.Labels["devcontainer.metadata"], hostGitTokenPlaceholder) {
		t.Errorf("devcontainer.metadata holds the token")
	}

	// The session carries the value, and the credential helper of the image
	// returns it: the helper works in the session.
	probe := "printenv GIT_TOKEN | wc -c\n" +
		`printf "protocol=https\nhost=host.example\n\n" | git credential fill | grep -F "password=platzh" && echo fill-ok` + "\n" +
		"exit\n"
	// The shell of the session resolves its own start phase with the fake
	// source: the shell run carries the same PATH as the up run.
	out, errOut, code := hostGitRun(ctx, t, bin, dir, probe, envPath, "devcontainer", "shell")
	if code != 0 {
		t.Fatalf("shell: exit %d\nstdout %.600s\nstderr %.600s", code, out, errOut)
	}
	if !strings.Contains(out, hostGitSessionLen) {
		t.Errorf("the session carries no length line %s; stdout %.600q stderr %.600q", hostGitSessionLen, out, errOut)
	}
	if !strings.Contains(out, "fill-ok") || !strings.Contains(out, hostGitTokenPlaceholder) {
		t.Errorf("the session or the credential helper carries no value; stdout %.600q stderr %.600q", out, errOut)
	}
}

// TestE2E_HostGitToken_NoSourceWarnsAndCarriesNoValue proves the missing
// source: the warning stands, the command does not abort, and the session
// carries no token.
func TestE2E_HostGitToken_NoSourceWarnsAndCarriesNoValue(t *testing.T) {
	hostGitUnsetToken(t)
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := hostGitPrepare(ctx, t, true, true)
	bin := ubootBinary(t)

	// A gh that is not logged in: exit 1 without a token (the form of M3).
	ghDir := hostGitFakeWrites(t, "gh", "exit 1")
	envPath := ghDir + ":" + os.Getenv("PATH")
	stdout, stderr, code := hostGitRun(ctx, t, bin, dir, "", envPath, "devcontainer", "up", "--json")
	if code != 0 {
		t.Fatalf("up (a missing source is no abort): exit %d\nstdout %.600s\nstderr %.600s", code, stdout, stderr)
	}
	env := hostGitEnvelopes(t, stdout)
	if env.Data.GitTokenSource != "" {
		t.Errorf("gitTokenSource = %q, want empty", env.Data.GitTokenSource)
	}
	if !hostGitHasWarning(env, "no host source carries a git token") {
		t.Errorf("warnings lack the missing-source text: %.600s", stdout)
	}
	if hostGitHasWarning(env, hostGitTokenPlaceholder) {
		t.Errorf("the envelope carries a value")
	}
}

// TestE2E_HostGitToken_SwitchOffDoesNotResolve proves the opt-in gate: with
// the switch absent the fake gh is not asked at all.
func TestE2E_HostGitToken_SwitchOffDoesNotResolve(t *testing.T) {
	hostGitUnsetToken(t)
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := hostGitPrepare(ctx, t, true, false)
	bin := ubootBinary(t)

	asked := t.TempDir() + "/asked"
	ghDir := hostGitFakeWrites(t, "gh", "echo placeholder > "+hostGitQuote(asked)+"; echo "+hostGitTokenPlaceholder)
	envPath := ghDir + ":" + os.Getenv("PATH")
	stdout, _, code := hostGitRun(ctx, t, bin, dir, "", envPath, "devcontainer", "up", "--json")
	if code != 0 {
		t.Fatalf("up: exit %d\nstdout %.600s", code, stdout)
	}
	env := hostGitEnvelopes(t, stdout)
	if env.Data.GitTokenSource != "" {
		t.Errorf("gitTokenSource = %q, want empty (the switch is off)", env.Data.GitTokenSource)
	}
	if _, err := os.Stat(asked); !os.IsNotExist(err) {
		t.Errorf("the fake gh ran (the switch is off); asked file: %v", err)
	}
}

// TestE2E_HostGitIdentity_KeepsTheContainerValue proves D2 over the binary:
// the start sets the identity of the host git config, a value the container
// already carries stays, and the identity is commit-ready afterwards.
func TestE2E_HostGitIdentity_KeepsTheContainerValue(t *testing.T) {
	hostGitUnsetToken(t)
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := hostGitPrepare(ctx, t, false, false)
	bin := ubootBinary(t)

	// A fake git: answers the two identity keys with placeholders, first the
	// value "1", after the marker "second" the value "2" — the container
	// must keep "1" over the second start.
	state := t.TempDir() + "/second"
	// The fake git answers the two identity keys; after the marker file it
		// answers the value suffix "1", so a second start must keep the value
		// the container carries.
		gitBody := "case \" $* \" in\n" +
			"  *' config " + hostGitIdentityNameKey + "'*) " + hostGitKeepEcho(state, hostGitNamePlaceholder) + " ;;\n" +
			"  *' config " + hostGitIdentityEmailKey + "'*) " + hostGitKeepEcho(state, hostGitEmailPlaceholder) + " ;;\n" +
			"  *) exit 1 ;;\nesac\n"
	gitDir := hostGitFakeWrites(t, "git", gitBody)
	envPath := gitDir + ":" + os.Getenv("PATH")

	stdout, stderr, code := hostGitRun(ctx, t, bin, dir, "", envPath, "devcontainer", "up", "--json")
	if code != 0 {
		t.Fatalf("up: exit %d\nstdout %.600s\nstderr %.600s", code, stdout, stderr)
	}
	env := hostGitEnvelopes(t, stdout)
	id := env.Data.ContainerID
	t.Cleanup(func() { _, _ = dcDocker("rm", "-f", "-v", id) })

	user := env.Data.RemoteUser
	if out, code := dcExec(ctx, t, id, user, "git config --global --get "+hostGitIdentityNameKey); code != 0 || strings.TrimSpace(out) != hostGitNamePlaceholder {
		t.Fatalf("user.name in the container = %q (exit %d), want %s", out, code, hostGitNamePlaceholder)
	}
	if out, code := dcExec(ctx, t, id, user, "git var "+hostGitIdentityGitVar); code != 0 || !strings.Contains(out, hostGitNamePlaceholder) {
		t.Fatalf("git var of the container = %q (exit %d), want an identity with the placeholder", out, code)
	}

	// The host changes its value; the container keeps the value it carries.
	if err := os.WriteFile(state, []byte("second"), 0o644); err != nil {
		t.Fatal(err)
	}
	stdout, stdErr2, code := hostGitRun(ctx, t, bin, dir, "", envPath, "devcontainer", "up", "--json")
	if code != 0 {
		t.Fatalf("second up: exit %d\nstdout %.600s\nstderr %.600s", code, stdout, stdErr2)
	}
	env2 := hostGitEnvelopes(t, stdout)
	if id2 := env2.Data.ContainerID; id2 != id {
		t.Fatalf("second up = %s, want the same container", id2)
	}
	if out, code := dcExec(ctx, t, id, user, "git config --global --get "+hostGitIdentityNameKey); code != 0 || strings.TrimSpace(out) != hostGitNamePlaceholder {
		t.Fatalf("user.name after the host change = %q (exit %d), want %s (the container keeps its value)", out, code, hostGitNamePlaceholder)
	}
}

// hostGitKeepEcho is the echo of one fake git answer: the value plus "1"
// after the marker file exists.
func hostGitKeepEcho(state, value string) string {
	return fmt.Sprintf("if test -f '%s'; then echo '%s1'; else echo '%s'; fi", state, value, value)
}

// hostGitIdentityGitVar is the git variable whose output proves a
// commit-ready identity (LH-FA-DEV-010).
const hostGitIdentityGitVar = "GIT_COMMITTER_IDENT"

// hostGitQuote quotes a value for a fake script; the placeholders carry no
// quote, so the check of the round-trip is the test itself.
func hostGitQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

// hostGitUnused keeps time in the import list; the budget of the tests uses
// it through dcTestBudget only.
var _ = time.Second