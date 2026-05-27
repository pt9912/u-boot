// Package git is the os/exec-backed implementation of the
// `port/driven.Git` interface (LH-FA-ARCH-002, LH-FA-INIT-007).
//
// The adapter shells out to the host `git` binary. A future
// optimization could swap this for go-git, but for u-boot's narrow
// surface (IsRepository, Init) os/exec is the simplest stable bet.
package git

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// gitNotARepoExitCode is the exit code git uses for
// "not a git repository" and similar fatal-but-categorical
// failures from rev-parse and friends.
const gitNotARepoExitCode = 128

// Git is the production git adapter. Construct with [New].
type Git struct {
	// binary lets tests substitute a stub via [WithBinary]; production
	// code uses the default "git".
	binary string
}

// Static check: Git satisfies the Git port.
var _ driven.Git = (*Git)(nil)

// New returns a Git adapter that shells out to the `git` binary on
// `$PATH`.
func New() *Git { return &Git{binary: "git"} }

// WithBinary overrides the git binary path; intended for tests.
func WithBinary(path string) *Git { return &Git{binary: path} }

// IsRepository reports whether dir is inside a git work tree. It runs
// `git -C <dir> rev-parse --is-inside-work-tree` and treats only
// exit code 128 ("not a git repository") as "no repo, no error".
// Any other exit code, an I/O problem, or a missing binary is
// returned as an error so subtle environmental issues do not
// silently masquerade as "no repo".
func (g Git) IsRepository(ctx context.Context, dir string) (bool, error) {
	cmd := exec.CommandContext(ctx, g.binary, "-C", dir, "rev-parse", "--is-inside-work-tree")
	out, err := cmd.CombinedOutput()
	if err == nil {
		return true, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == gitNotARepoExitCode {
		return false, nil
	}
	return false, fmt.Errorf("git rev-parse failed: %w (output: %s)", err, string(out))
}

// Init runs `git init` in dir. The caller is expected to have checked
// [IsRepository] first (LH-FA-INIT-007 forbids re-initializing).
func (g Git) Init(ctx context.Context, dir string) error {
	cmd := exec.CommandContext(ctx, g.binary, "-C", dir, "init")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git init failed: %w (output: %s)", err, string(out))
	}
	return nil
}

// Version runs `git --version` and returns the bare version string
// (`"2.43.0"`). The raw `git --version` output is typically
// `"git version 2.43.0\n"`; the adapter strips the `"git version "`
// prefix and trims whitespace so the application layer can pass the
// string straight into a semver comparator.
//
// A non-nil error signals that the git binary is unavailable or
// failed to run — the application maps that to `git.installed: error`.
func (g Git) Version(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, g.binary, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git --version failed: %w (output: %s)", err, string(out))
	}
	raw := strings.TrimSpace(string(out))
	// Typical output: "git version 2.43.0" — strip the prefix.
	const prefix = "git version "
	if strings.HasPrefix(raw, prefix) {
		return strings.TrimSpace(raw[len(prefix):]), nil
	}
	// Unrecognized format: return the trimmed output verbatim so the
	// service can still log it; the semver parser will simply mark it
	// invalid downstream.
	return raw, nil
}
