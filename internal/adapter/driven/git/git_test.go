package git_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/git"
)

// gitAvailable returns true when the host `git` binary is reachable.
// Tests skip with a clear message when it isn't — that keeps the
// default `make test` path green on `golang:1.26.3` (which has git
// pre-installed) and on lean CI runners that don't.
func gitAvailable(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not available: %v", err)
	}
}

func TestGit_IsRepository_OnFreshDirReturnsFalse(t *testing.T) {
	gitAvailable(t)
	dir := t.TempDir()

	got, err := git.New().IsRepository(context.Background(), dir)
	if err != nil {
		t.Fatalf("IsRepository: %v", err)
	}
	if got {
		t.Fatalf("IsRepository(empty dir) = true, want false")
	}
}

func TestGit_Init_ThenIsRepositoryReturnsTrue(t *testing.T) {
	gitAvailable(t)
	dir := t.TempDir()
	ctx := context.Background()

	if err := git.New().Init(ctx, dir); err != nil {
		t.Fatalf("Init: %v", err)
	}

	got, err := git.New().IsRepository(ctx, dir)
	if err != nil {
		t.Fatalf("IsRepository post-init: %v", err)
	}
	if !got {
		t.Fatalf("IsRepository post-init = false, want true")
	}

	// Sanity: .git directory exists.
	if _, err := os.Stat(filepath.Join(dir, ".git")); err != nil {
		t.Fatalf("post-init .git missing: %v", err)
	}
}

func TestGit_IsRepository_MissingBinaryReturnsError(t *testing.T) {
	// Why: Hardens the error path of the os/exec branch — when the
	// binary itself is missing, the adapter must return an error
	// rather than silently reporting "not a repo".
	g := git.WithBinary("/does/not/exist/git-binary")
	_, err := g.IsRepository(context.Background(), t.TempDir())
	if err == nil {
		t.Fatalf("IsRepository with missing binary: expected error, got nil")
	}
}

