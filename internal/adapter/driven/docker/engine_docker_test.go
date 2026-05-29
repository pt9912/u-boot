//go:build docker

// Integration test for the [docker.Engine] adapter — runs only when
// the project is built with `-tags=docker` (the `make test-docker`
// target, extended by `slice-m6-docker-integrationstests`). Without
// the tag the file is excluded from compilation, so the regular
// `make test` does not require a host docker installation.
//
// T2 lays down the build-tag skeleton; the full LH-AK-002 happy-
// path test (init → add postgres → up → healthy port) lives in
// the carveout slice's CI job, not here. This smoke test only
// confirms the build-tag wiring works.

package docker_test

import (
	"context"
	"errors"
	"os/exec"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/docker"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

func TestEngine_RealDocker_PreflightGreen(t *testing.T) {
	// Smoke: with a real docker available, ComposePs against a
	// non-existent compose.yaml must still pass the preflight
	// (LookPath + version + compose version all succeed) and then
	// fail at the actual `compose ps` call — classified as
	// ErrComposeRuntime (code 12), not ErrDockerUnavailable (code 11).
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skipf("docker binary not on PATH: %v", err)
	}
	e := docker.NewEngine()
	_, err := e.ComposePs(context.Background(), "/tmp/u-boot-no-such-project")
	if err == nil {
		t.Fatalf("expected an error against a non-existent compose project, got nil")
	}
	if errors.Is(err, driven.ErrDockerUnavailable) {
		t.Errorf("preflight should be green with real docker; got ErrDockerUnavailable: %v", err)
	}
	if !errors.Is(err, driven.ErrComposeRuntime) {
		t.Errorf("expected ErrComposeRuntime classification, got: %v", err)
	}
}
