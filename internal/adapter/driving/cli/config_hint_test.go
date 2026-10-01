package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// slice-v1-config-structured-hint: the error envelope of `config get`
// / `config set` carries `data.hint`; exit code stays 10; the human
// message is unchanged.
func TestConfig_ErrorEnvelope_StructuredHint(t *testing.T) {
	cases := []struct {
		name string
		args []string
		uc   *fakeConfigUseCase
		want map[string]string
	}{
		{"get not set (add)", []string{"--json", "config", "get", "services.postgres.enabled"},
			&fakeConfigUseCase{getErr: driving.NewConfigHintError(fmt.Errorf("%w: services.postgres.enabled — run `u-boot add postgres`", driving.ErrConfigValueNotSet),
				driving.ConfigHint{Command: "u-boot add postgres", Action: "add", Argument: "postgres"})},
			map[string]string{"command": "u-boot add postgres", "action": "add", "argument": "postgres"}},
		{"set write rejected", []string{"--json", "config", "set", "services.postgres.enabled", "true"},
			&fakeConfigUseCase{setErr: driving.NewConfigHintError(fmt.Errorf("%w: not writable", driving.ErrConfigWriteRejected),
				driving.ConfigHint{Command: "u-boot add postgres", Action: "add", Argument: "postgres"})},
			map[string]string{"command": "u-boot add postgres", "action": "add", "argument": "postgres"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			err := newAppWithConfigStub(tc.uc).Execute(context.Background(), tc.args, &stdout, &stderr)
			if got := cli.ExitCode(err); got != 10 {
				t.Errorf("ExitCode = %d, want 10", got)
			}
			var env struct {
				Data struct {
					Hint map[string]string `json:"hint"`
				} `json:"data"`
			}
			if jerr := json.Unmarshal(stdout.Bytes(), &env); jerr != nil {
				t.Fatalf("envelope: %v\n%s", jerr, stdout.String())
			}
			for k, v := range tc.want {
				if env.Data.Hint[k] != v {
					t.Errorf("data.hint[%s] = %q, want %q\n%s", k, env.Data.Hint[k], v, stdout.String())
				}
			}
		})
	}
}

// Errors without a hint keep the envelope without data.hint.
func TestConfig_ErrorEnvelope_NoHint(t *testing.T) {
	uc := &fakeConfigUseCase{getErr: fmt.Errorf("%w: boom", driving.ErrConfigSchemaInvalid)}
	var stdout, stderr bytes.Buffer
	_ = newAppWithConfigStub(uc).Execute(context.Background(), []string{"--json", "config", "get", "project.name"}, &stdout, &stderr)
	if bytes.Contains(stdout.Bytes(), []byte(`"hint"`)) {
		t.Errorf("unexpected hint in envelope:\n%s", stdout.String())
	}
}
