package main

import (
	"bytes"
	"context"
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
