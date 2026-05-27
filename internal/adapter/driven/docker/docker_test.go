package docker_test

import (
	"context"
	"os/exec"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/docker"
)

// dockerAvailable returns true when the host `docker` binary is on
// PATH. Most CI runners do not have it (the project's own gates
// image runs in golang:1.26.3, no docker-in-docker); tests skip
// cleanly in that case. Full daemon-tests live under build tag
// `docker` (slice-m6-docker-integrationstests).
func dockerAvailable(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skipf("docker not available: %v", err)
	}
}

func TestProbe_Version_HappyPath(t *testing.T) {
	dockerAvailable(t)
	got, err := docker.New().Version(context.Background())
	if err != nil {
		t.Fatalf("Version: %v", err)
	}
	if got == "" {
		t.Fatalf("Version returned empty string")
	}
}

func TestProbe_Version_MissingBinaryReturnsError(t *testing.T) {
	p := docker.WithBinary("/does/not/exist/docker-binary")
	_, err := p.Version(context.Background())
	if err == nil {
		t.Fatalf("Version with missing binary: expected error, got nil")
	}
}

func TestProbe_Info_MissingBinaryReturnsError(t *testing.T) {
	p := docker.WithBinary("/does/not/exist/docker-binary")
	if err := p.Info(context.Background()); err == nil {
		t.Fatalf("Info with missing binary: expected error, got nil")
	}
}

func TestProbe_ComposeVersion_MissingBinaryReturnsError(t *testing.T) {
	p := docker.WithBinary("/does/not/exist/docker-binary")
	_, err := p.ComposeVersion(context.Background())
	if err == nil {
		t.Fatalf("ComposeVersion with missing binary: expected error, got nil")
	}
}
