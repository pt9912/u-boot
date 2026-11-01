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

// ADR-0038 decision 6 (LH-FA-DEV-010): the rule for an instance name has one
// definition; each edge of the rule has a case.
func TestValidateDevcontainerInstance(t *testing.T) {
	for _, ok := range []string{"a", "0", "a-b", "agent-1", "a--b", "1-2", strings.Repeat("a", 32), "a" + strings.Repeat("-", 30) + "a", "0" + strings.Repeat("9", 31)} {
		if err := driven.ValidateDevcontainerInstance(ok); err != nil {
			t.Errorf("instance %q rejected: %v", ok, err)
		}
	}
	bad := map[string]string{
		"empty":             "",
		"upper case":        "Abc",
		"upper case inside": "aBc",
		"leading hyphen":    "-a",
		"trailing hyphen":   "a-",
		"only a hyphen":     "-",
		"equals sign":       "a=b",
		"only equals sign":  "=",
		"underscore":        "a_b",
		"dot":               "a.b",
		"space":             "a b",
		"non-ASCII":         "ä",
		"slash":             "a/b",
		"line break":        "a\nb",
		"NUL":               "a\x00b",
		"33 characters":     strings.Repeat("a", 33),
		"128 characters":    strings.Repeat("a", 128),
	}
	for name, value := range bad {
		if err := driven.ValidateDevcontainerInstance(value); err == nil {
			t.Errorf("%s: instance %q accepted", name, value)
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

// LH-FA-DEV-010: the character set the engine accepts for a container name
// has one definition, measured with Docker Engine 29.9.2 (measurement 4 of the
// container-name measurement). Each edge has a case; a length limit was not
// measured down to 257 characters, so the rule carries none.
func TestValidateDockerContainerName(t *testing.T) {
	for _, ok := range []string{"u", "u-boot-demo", "u-boot-demo-agent-1", "9ubmess", "Uubmess", "a.b_c-d", "u-boot-9"} {
		if err := driven.ValidateDockerContainerName(ok); err != nil {
			t.Errorf("container name %q rejected: %v", ok, err)
		}
	}
	for _, bad := range []string{"", "-x", ".x", "_x", "a b", `a"b`, "a/b", "a:b", "a@b", "ä", "a\nb", "a\x00b"} {
		if err := driven.ValidateDockerContainerName(bad); err == nil {
			t.Errorf("container name %q accepted", bad)
		}
	}
}
