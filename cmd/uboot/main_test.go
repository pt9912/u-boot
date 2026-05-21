package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_VersionFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"--version"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("--version exit code: got %d, want 0", code)
	}
	if got := strings.TrimSpace(stdout.String()); got == "" {
		t.Fatalf("--version stdout: empty; expected non-empty version string")
	}
	if stderr.Len() != 0 {
		t.Fatalf("--version stderr: got %q, want empty", stderr.String())
	}
}

func TestRun_HelpFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"--help"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("--help exit code: got %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("--help stdout: missing 'Usage:'; got %q", stdout.String())
	}
}

func TestRun_NoArgsPrintsHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run(nil, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("no-args exit code: got %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("no-args stdout: missing 'Usage:'; got %q", stdout.String())
	}
}

func TestRun_UnknownCommandExits2(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"frobnicate"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("unknown command exit code: got %d, want 2 (LH-FA-CLI-006)", code)
	}
	if !strings.Contains(stderr.String(), "not implemented") {
		t.Fatalf("unknown command stderr: missing 'not implemented'; got %q", stderr.String())
	}
}

func TestRun_UnknownFlagExits2(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"--no-such-flag"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("unknown flag exit code: got %d, want 2 (LH-FA-CLI-006)", code)
	}
}

func TestVersionConstantNonEmpty(t *testing.T) {
	if version == "" {
		t.Fatal("package-level version must not be empty")
	}
}
