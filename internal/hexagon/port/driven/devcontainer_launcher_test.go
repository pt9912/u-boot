package driven_test

import (
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// ADR-0028 Entscheidung 4: the admissible character set of the project folder
// has one definition.
func TestValidateDevcontainerProjectFolder(t *testing.T) {
	for _, ok := range []string{"/", "/tmp/project", "/tmp/with space/p", "/tmp/ümlaut"} {
		if err := driven.ValidateDevcontainerProjectFolder(ok); err != nil {
			t.Errorf("folder %q rejected: %v", ok, err)
		}
	}
	for _, bad := range []string{"", "relative/dir", "./x", "/tmp/a,b", `/tmp/a"b`, "/tmp/a\nb", "/tmp/a\rb", "/tmp/a\x00b"} {
		if err := driven.ValidateDevcontainerProjectFolder(bad); err == nil {
			t.Errorf("folder %q accepted", bad)
		}
	}
}

// LH-FA-DEV-009: a variable name never carries a value into `docker run -e`.
func TestValidateDevcontainerEnvName(t *testing.T) {
	for _, ok := range []string{"A", "_", "git_token", "GIT_TOKEN_2", "a1"} {
		if err := driven.ValidateDevcontainerEnvName(ok); err != nil {
			t.Errorf("name %q rejected: %v", ok, err)
		}
	}
	for _, bad := range []string{"", "1A", "A=x", "A B", "A-B", "Ä", "A\n"} {
		if err := driven.ValidateDevcontainerEnvName(bad); err == nil {
			t.Errorf("name %q accepted", bad)
		} else if !strings.Contains(err.Error(), "plain name") {
			t.Errorf("name %q: error %v", bad, err)
		}
	}
}

// ADR-0031 Entscheidung 1: the user of the session is one argument of
// `docker exec -u`.
func TestValidateDevcontainerShellUser(t *testing.T) {
	for _, ok := range []string{"vscode", "root", "1000", "vscode:vscode", "1000:1000", "a b"} {
		if err := driven.ValidateDevcontainerShellUser(ok); err != nil {
			t.Errorf("user %q rejected: %v", ok, err)
		}
	}
	for _, bad := range []string{"", "-u", "--privileged", "a\nb", "a\rb", "a\x00b"} {
		if err := driven.ValidateDevcontainerShellUser(bad); err == nil {
			t.Errorf("user %q accepted", bad)
		}
	}
}

// ADR-0031 Entscheidung 1: the workspace of the session is one argument of
// `docker exec -w`.
func TestValidateDevcontainerWorkspace(t *testing.T) {
	for _, ok := range []string{"/", "/workspaces/project", "/workspaces/with space", "/w/a,b", `/w/a"b`} {
		if err := driven.ValidateDevcontainerWorkspace(ok); err != nil {
			t.Errorf("workspace %q rejected: %v", ok, err)
		}
	}
	for _, bad := range []string{"", "relative/dir", "-w", "/w/a\nb", "/w/a\rb", "/w/a\x00b"} {
		if err := driven.ValidateDevcontainerWorkspace(bad); err == nil {
			t.Errorf("workspace %q accepted", bad)
		}
	}
}
