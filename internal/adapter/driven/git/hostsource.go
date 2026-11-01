// Package git — see the package comment of the [driven.Git] adapter
// (git.go). hostsource.go is the [driven.HostGitSource] adapter: it reads
// the git token and the git identity of the host for the devcontainer
// commands (LH-FA-DEV-009, LH-FA-DEV-010) over `gh` and `git` of the host.
package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// gitConfigMissingExitCode is the exit code of `git config <key>` when no
// line carries the key; an empty value is exit 0 with empty output, so the
// adapter reports an empty value only for exit code 1.
const gitConfigMissingExitCode = 1

// tokenEnvName is the variable of source 1, LH-FA-DEV-009.
const tokenEnvName = "GIT_TOKEN"

// HostSource is the production [driven.HostGitSource] adapter. It calls
// `git` and `gh` of the host by binary path; construct with [NewHostSource].
// Tests substitute script paths through the setters ([WithHostGitBinary],
// [WithHostGhBinary]).
type HostSource struct {
	// gitBinary is the git binary that [HostSource.ConfigValue] uses.
	gitBinary string
	// ghBinary is the gh binary that [HostSource.GitToken] source 2 uses.
	ghBinary string
}

// Static check: HostSource satisfies the HostGitSource port.
var _ driven.HostGitSource = (*HostSource)(nil)

// NewHostSource returns a HostSource that calls `git` and `gh` from the
// PATH.
func NewHostSource() *HostSource {
	return &HostSource{gitBinary: "git", ghBinary: "gh"}
}

// WithHostGitBinary overrides the git binary path; intended for tests.
func (h *HostSource) WithHostGitBinary(path string) *HostSource {
	h.gitBinary = path
	return h
}

// WithHostGhBinary overrides the gh binary path; intended for tests.
func (h *HostSource) WithHostGhBinary(path string) *HostSource {
	h.ghBinary = path
	return h
}

// GitToken implements [driven.HostGitSource.GitToken]. Source 1 is the
// GIT_TOKEN variable of the environment of the process; the adapter reads
// the environment of its own process (the launcher of the same run reads
// the same process environment, so both see one state of the host). Source
// 2 is `gh auth token`: any exit code other than 0 — gh missing, gh not
// logged in — skips the source, so a host without a gh session stays silent
// instead of failing the call.
func (h *HostSource) GitToken(ctx context.Context) (string, string, error) {
	if value := os.Getenv(tokenEnvName); value != "" {
		return value, driven.TokenSourceEnvironment, nil
	}
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, h.ghBinary, "auth", "token")
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil || strings.TrimSpace(out.String()) == "" {
		return "", "", nil
	}
	return strings.TrimSpace(out.String()), driven.TokenSourceGh, nil
}

// ConfigValue implements [driven.HostGitSource.ConfigValue]: it runs
// `git -C <dir> config <key>` and treats exit code 1 as a missing key.
func (h *HostSource) ConfigValue(ctx context.Context, dir, key string) (string, error) {
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, h.gitBinary, "-C", dir, "config", key)
	cmd.Stdout = &out
	err := cmd.Run()
	if err == nil {
		return strings.TrimSpace(out.String()), nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == gitConfigMissingExitCode {
		return "", nil
	}
	return "", fmt.Errorf("git config %s in %s failed: %w", key, dir, err)
}