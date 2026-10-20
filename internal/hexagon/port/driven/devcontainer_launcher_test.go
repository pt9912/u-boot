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
