package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_VersionFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run(context.Background(), []string{"--version"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("--version exit code: got %d, want 0", code)
	}
	if got := strings.TrimSpace(stdout.String()); got == "" {
		t.Fatalf("--version stdout: empty; expected non-empty version string")
	}
}

func TestRun_HelpFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run(context.Background(), []string{"--help"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("--help exit code: got %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "u-boot") {
		t.Fatalf("--help stdout: missing 'u-boot'; got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "init") {
		t.Fatalf("--help stdout: missing subcommand 'init'; got %q", stdout.String())
	}
}

func TestRun_UnknownCommandExits2(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run(context.Background(), []string{"frobnicate"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("unknown command exit code: got %d, want 2 (LH-FA-CLI-006)", code)
	}
	if !strings.Contains(stderr.String(), "u-boot:") {
		t.Fatalf("unknown command stderr: missing 'u-boot:' prefix; got %q", stderr.String())
	}
}

func TestRun_UnknownFlagExits2(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run(context.Background(), []string{"init", "--no-such-flag"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("unknown flag exit code: got %d, want 2 (LH-FA-CLI-006)", code)
	}
}

func TestVersionConstantNonEmpty(t *testing.T) {
	if version == "" {
		t.Fatal("package-level version must not be empty")
	}
}

func TestRun_InitHappyPath_WiresAllAdapters(t *testing.T) {
	// Integration-style test against the real wiring: real fs, yaml,
	// git adapters all exercised once. Without --no-git this would
	// also exercise the git adapter (host git binary required);
	// staying with --no-git keeps the test self-contained on lean
	// CI runners.
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	dir := t.TempDir()
	t.Chdir(dir) // Go 1.24+; PWD is the resolveProjectName fallback source.

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"init", "demo", "--no-git"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("init demo --no-git: exit %d, stderr=%q", code, stderr.String())
	}

	// Spot-check the artefacts produced (LH-FA-INIT-003 mandatory
	// set). Full per-file content is asserted in the application
	// layer's unit tests.
	for _, rel := range []string{
		"u-boot.yaml", "compose.yaml", "README.md", "CHANGELOG.md",
		".env.example", ".gitignore",
		"docker", "scripts", "docs",
	} {
		if _, err := os.Stat(filepath.Join(dir, rel)); err != nil {
			t.Errorf("init did not create %s: %v", rel, err)
		}
	}

	if !strings.Contains(stdout.String(), "demo") {
		t.Errorf("stdout does not mention project name; got %q", stdout.String())
	}
}
