package git_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/git"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// hostSourceTests run with script binaries (the pattern of
// [`TestGit_IsRepository_MissingBinaryReturnsError`]): the tests set the
// binary paths of the adapter to skripts of the test, so the answers come
// from controlled exits and outputs, never from a `gh` of the host.
func writeScript(t *testing.T, dir, name, body string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return path
}

// scriptBinary returns a script that prints its lines and exits with the
// given status; a body `echo value` prints one line.
func TestHostSource_GitToken_HostVariableWins(t *testing.T) {
	dir := t.TempDir()
	gh := writeScript(t, dir, "gh", "echo gh-token-placeholder-0123\n")
	t.Setenv("GIT_TOKEN", "env-token-placeholder-0123")
	hs := git.NewHostSource().WithHostGhBinary(gh)
	value, source, err := hs.GitToken(context.Background())
	if err != nil {
		t.Fatalf("GitToken: %v", err)
	}
	if source != driven.TokenSourceEnvironment {
		t.Fatalf("source = %q, want the environment source", source)
	}
	if value == "" {
		t.Fatalf("value empty, want the environment value")
	}
}

func TestHostSource_GitToken_FallsBackToGh(t *testing.T) {
	dir := t.TempDir()
	gh := writeScript(t, dir, "gh", "echo gh-token-placeholder-0123")
	hs := git.NewHostSource().WithHostGhBinary(gh)
	value, source, err := hs.GitToken(context.Background())
	if err != nil {
		t.Fatalf("GitToken: %v", err)
	}
	if source != driven.TokenSourceGh {
		t.Fatalf("source = %q, want %q", source, driven.TokenSourceGh)
	}
	if value != "gh-token-placeholder-0123" {
		t.Fatalf("value = %q, want the script output", value)
	}
}

func TestHostSource_GitToken_GhMissingReadsAsNoSource(t *testing.T) {
	// A missing binary: the adapter reports no source without an error.
	dir := t.TempDir()
	gh := filepath.Join(dir, "gh-absent")
	hs := git.NewHostSource().WithHostGhBinary(gh)
	value, source, err := hs.GitToken(context.Background())
	if err != nil {
		t.Fatalf("GitToken: %v (a missing gh is a skipped source, no error)", err)
	}
	if value != "" || source != "" {
		t.Fatalf("value=%q source=%q, want empty both", value, source)
	}
}

func TestHostSource_GitToken_GhFailingReadsAsNoSource(t *testing.T) {
	dir := t.TempDir()
	// exit 1 without output: the "not logged in" of gh (M3).
	gh := writeScript(t, dir, "gh", "exit 1")
	hs := git.NewHostSource().WithHostGhBinary(gh)
	value, source, err := hs.GitToken(context.Background())
	if err != nil {
		t.Fatalf("GitToken: %v", err)
	}
	if value != "" || source != "" {
		t.Fatalf("value=%q source=%q, want empty both", value, source)
	}
}

func TestHostSource_GitToken_GhEmptyOutputReadsAsNoSource(t *testing.T) {
	dir := t.TempDir()
	gh := writeScript(t, dir, "gh", "")
	hs := git.NewHostSource().WithHostGhBinary(gh)
	value, source, err := hs.GitToken(context.Background())
	if err != nil {
		t.Fatalf("GitToken: %v", err)
	}
	if value != "" || source != "" {
		t.Fatalf("value=%q source=%q, want empty both", value, source)
	}
}

func TestHostSource_ConfigValue_ReadsTheEffectiveValue(t *testing.T) {
	dir := t.TempDir()
	if out := run(t, dir, "git", "init", "-q", "."); out != "" {
		t.Fatalf("git init: %s", out)
	}
	g := git.NewHostSource()
	value, err := g.ConfigValue(context.Background(), dir, "user.name")
	if err != nil || value != "" {
		t.Fatalf("before set: value=%q err=%v", value, err)
	}
	run(t, dir, "git", "config", "user.name", "h-name")
	value, err = g.ConfigValue(context.Background(), dir, "user.name")
	if err != nil {
		t.Fatalf("ConfigValue: %v", err)
	}
	if value != "h-name" {
		t.Fatalf("value = %q, want the set value", value)
	}
}

func TestHostSource_ConfigValue_MissingKeyIsEmptyWithoutError(t *testing.T) {
	dir := t.TempDir()
	run(t, dir, "git", "init", "-q", ".")
	value, err := git.NewHostSource().ConfigValue(context.Background(), dir, "user.email")
	if err != nil {
		t.Fatalf("ConfigValue: %v", err)
	}
	if value != "" {
		t.Fatalf("value = %q, want an empty one", value)
	}
}

func TestHostSource_ConfigValue_OtherFailureIsAnError(t *testing.T) {
	// exit code 2 of a script as git binary: a failure beyond the missing
	// key is an error (M5: fatal failures carry 128 and beyond).
	dir := t.TempDir()
	script := writeScript(t, dir, "git", "echo broken 1>&2; exit 128")
	_, err := git.NewHostSource().WithHostGitBinary(script).ConfigValue(
		context.Background(), t.TempDir(), "user.name")
	if err == nil || !strings.Contains(err.Error(), "failed") {
		t.Fatalf("err = %v, want a failure of the call", err)
	}
}

// run runs the given binary in dir; a failing call fails the test.
func run(t *testing.T, dir, binary string, args ...string) string {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("%s %v: %v (%s)", binary, args, err, out)
	}
	return string(out)
}