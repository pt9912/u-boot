//go:build docker

// End-to-end check of the OAuth token of Claude Code in the sandbox
// devcontainer (LH-FA-DEV-009, LH-FA-DEV-003): `CLAUDE_CODE_OAUTH_TOKEN`
// reaches the session of `u-boot devcontainer shell` through the generated
// reference, an unset host variable stays unset, an empty one arrives empty,
// and `claude -p` with a placeholder gets an answer of the service. The tests
// use an invented value, never a valid token; the build installs `node` and
// `claude-code` from the network.

package e2e_test

import (
	"context"
	"os"
	"strings"
	"testing"
)

const (
	oauthName        = "CLAUDE_CODE_OAUTH_TOKEN"
	oauthPlaceholder = "sk-ant-oat01-placeholderPLACEHOLDERplaceholder0123456789"
)

// oauthSession reads the variable in one session of `devcontainer shell` and
// then runs `claude -p` with a short timeout; the lines are "VALUE:<value>"
// or "VALUE-UNSET", "KEY-UNSET" for the API key, and "CLAUDE-EXIT:<n>" with
// the output of the call after it.
const oauthSession = `if [ "${CLAUDE_CODE_OAUTH_TOKEN+set}" = set ]; then echo "VALUE:$CLAUDE_CODE_OAUTH_TOKEN"; else echo VALUE-UNSET; fi
if [ "${ANTHROPIC_API_KEY+set}" = set ]; then echo KEY-SET; else echo KEY-UNSET; fi
exit
`

const oauthClaudeCall = `timeout 120 claude -p "Reply with the single word OK." 2>&1 | head -c 600; echo "CLAUDE-EXIT:${PIPESTATUS[0]}"
exit
`

// LH-FA-DEV-009 (agent credential): the reference generated for
// `claude-code` carries the variable of the host into the session. Cells:
// set to a placeholder, unset, empty. Guard: the placeholder stands in no
// generated file. `claude -p` with the placeholder ends with an answer of the
// service (an authentication error), not a hang.
func TestE2E_ClaudeOAuthToken_ReachesTheSession(t *testing.T) {
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	t.Setenv("GIT_TOKEN", sandboxCanaryToken)
	t.Setenv("ANTHROPIC_API_KEY", "")
	_ = os.Unsetenv("ANTHROPIC_API_KEY")
	t.Setenv(oauthName, oauthPlaceholder)
	dir := dcGenerateProject(ctx, t, true, map[string]string{
		"devcontainer.features.node.enabled":        "true",
		"devcontainer.features.claude-code.enabled": "true",
	})
	raw, err := os.ReadFile(dir + "/.devcontainer/devcontainer.json")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"`+oauthName+`": "${localEnv:`+oauthName+`}"`) {
		t.Fatalf("devcontainer.json lacks the reference:\n%s", raw)
	}
	if strings.Contains(string(raw), oauthPlaceholder) {
		t.Errorf("the placeholder stands in devcontainer.json")
	}
	bin := ubootBinary(t)

	for run := 1; run <= 2; run++ {
		stdout, stderr, code := runDCShell(ctx, t, bin, dir, oauthSession)
		if code != 0 {
			t.Fatalf("set, run %d: exit %d\nstdout:\n%s\nstderr:\n%.1500s", run, code, stdout, stderr)
		}
		if lines := strings.Split(strings.TrimSpace(stdout), "\n"); len(lines) != 2 || lines[0] != "VALUE:"+oauthPlaceholder || lines[1] != "KEY-UNSET" {
			t.Errorf("set, run %d: session lines %q, want the placeholder and an unset API key", run, stdout)
		}
		if strings.Contains(stderr, oauthPlaceholder) {
			t.Errorf("set, run %d: standard error carries the value", run)
		}
	}

	stdout, stderr, code := runDCShell(ctx, t, bin, dir, oauthClaudeCall)
	t.Logf("claude -p with the placeholder (exit %d):\n%s", code, stdout)
	if code != 0 || !strings.Contains(stdout, "CLAUDE-EXIT:") || strings.Contains(stdout, "CLAUDE-EXIT:124") {
		t.Errorf("claude -p: exit %d, stdout %q, stderr %.500q; want an end with an answer, not the timeout", code, stdout, stderr)
	}
	if body := strings.TrimSpace(strings.SplitN(stdout, "CLAUDE-EXIT:", 2)[0]); body == "" {
		t.Errorf("claude -p wrote no answer of the service: %q", stdout)
	}

	_ = os.Unsetenv(oauthName)
	for run := 1; run <= 2; run++ {
		stdout, stderr, code := runDCShell(ctx, t, bin, dir, oauthSession)
		if code != 0 {
			t.Fatalf("unset, run %d: exit %d\nstdout:\n%s\nstderr:\n%.1500s", run, code, stdout, stderr)
		}
		if lines := strings.Split(strings.TrimSpace(stdout), "\n"); len(lines) != 2 || lines[0] != "VALUE-UNSET" {
			t.Errorf("unset, run %d: session lines %q, want VALUE-UNSET", run, stdout)
		}
	}

	t.Setenv(oauthName, "")
	for run := 1; run <= 2; run++ {
		stdout, stderr, code := runDCShell(ctx, t, bin, dir, oauthSession)
		if code != 0 {
			t.Fatalf("empty, run %d: exit %d\nstdout:\n%s\nstderr:\n%.1500s", run, code, stdout, stderr)
		}
		if lines := strings.Split(strings.TrimSpace(stdout), "\n"); len(lines) != 2 || lines[0] != "VALUE:" {
			t.Errorf("empty, run %d: session lines %q, want an empty variable that is set", run, stdout)
		}
	}
}
