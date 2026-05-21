// Package git is the os/exec-backed implementation of the
// `port/driven.Git` interface (LH-FA-ARCH-002, LH-FA-INIT-007).
//
// The adapter shells out to the host `git` binary. A future
// optimization could swap this for go-git, but for u-boot's narrow
// surface (IsRepository, Init) os/exec is the simplest stable bet.
package git

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

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
// `git -C <dir> rev-parse --is-inside-work-tree` and treats a non-zero
// exit as "not a repository" (no error). A real I/O / `git` invocation
// problem is returned as an error.
func (g Git) IsRepository(dir string) (bool, error) {
	cmd := exec.CommandContext(context.Background(), g.binary, "-C", dir, "rev-parse", "--is-inside-work-tree")
	out, err := cmd.CombinedOutput()
	if err == nil {
		return true, nil
	}
	if _, ok := err.(*exec.ExitError); ok {
		// Plain exit-code failure → not a repo.
		return false, nil
	}
	return false, fmt.Errorf("git rev-parse failed: %w (output: %s)", err, string(out))
}

// Init runs `git init` in dir. The caller is expected to have checked
// [IsRepository] first (LH-FA-INIT-007 forbids re-initializing).
func (g Git) Init(dir string) error {
	cmd := exec.CommandContext(context.Background(), g.binary, "-C", dir, "init")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git init failed: %w (output: %s)", err, string(out))
	}
	return nil
}
