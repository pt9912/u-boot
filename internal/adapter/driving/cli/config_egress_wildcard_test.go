package cli_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/fs"
	"github.com/pt9912/u-boot/internal/adapter/driven/yaml"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/application"
)

// LH-FA-DEV-008, LH-FA-CLI-006: `config set` with a wildcard egress
// entry fails with exit code 10, names the host to enter in the error
// text, and leaves u-boot.yaml untouched. Single pair (setEgressAllow)
// and several pairs (planEgressAllow) run the production chain:
// ConfigService over the real filesystem and YAML codec.
func TestConfigSet_EgressAllowWildcard_Exit10WithHint(t *testing.T) {
	t.Parallel()
	const seed = "schemaVersion: 1\nproject:\n  name: t-uboot-config\ndevcontainer:\n  enabled: true\n"
	const key = "devcontainer.sandbox.egress.allow"
	cases := []struct {
		name string
		args []string
		want []string
	}{
		{"single pair", []string{"config", "set", key, "*.example.org"},
			[]string{`entry "*.example.org" is a wildcard`, "an allowed name already covers its subdomains", `enter "example.org" instead`}},
		{"single pair, wildcard after a valid entry", []string{"config", "set", key, "api.example.org, *.example.org"},
			[]string{`entry "*.example.org" is a wildcard`, "an allowed name already covers its subdomains", `enter "example.org" instead`}},
		{"single pair, no valid name behind the star", []string{"config", "set", key, "*.*.example.org"},
			[]string{`entry "*.*.example.org" is a wildcard`, "an allowed name already covers its subdomains", "without `*`"}},
		{"several pairs", []string{"config", "set", "devcontainer.profile", "sandbox", key, "*.example.org"},
			[]string{`entry "*.example.org" is a wildcard`, "an allowed name already covers its subdomains", `enter "example.org" instead`}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := filepath.Join(dir, "u-boot.yaml")
			if err := os.WriteFile(path, []byte(seed), 0o644); err != nil {
				t.Fatalf("seed: %v", err)
			}
			svc := application.NewConfigService(fs.New(), yaml.New(), nil)
			app := newAppWithConfig(svc, cli.WithGetwd(func() (string, error) { return dir, nil }))
			var stdout, stderr bytes.Buffer
			err := app.Execute(context.Background(), tc.args, &stdout, &stderr)
			if err == nil {
				t.Fatal("want an error for a wildcard entry")
			}
			if got := cli.ExitCode(err); got != 10 {
				t.Errorf("exit = %d, want 10 (err=%v)", got, err)
			}
			for _, w := range tc.want {
				if !strings.Contains(err.Error(), w) {
					t.Errorf("error %q lacks %q", err.Error(), w)
				}
			}
			body, _ := os.ReadFile(path)
			if string(body) != seed {
				t.Errorf("u-boot.yaml changed on rejected set:\n%s", body)
			}
		})
	}
}
