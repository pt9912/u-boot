//go:build docker

// LH-FA-UP-001 §966 healthcheck-domination pin
// (M6-docker-int Sub-T2).
//
// Spec §966: "Für Dienste mit Healthcheck ist `healthy` als
// Zielzustand erforderlich." UpService MUST keep polling while a
// service is running-but-not-yet-healthy; it MUST NOT stabilize on
// `running` alone.
//
// Operational pin: a fixture whose healthcheck deliberately takes
// ~3 s to flip to `healthy`. UpService.Up is given a 30 s budget;
// the call must return Stabilized=true, and the elapsed wall-clock
// MUST be ≥ 2 s (slack for the actual transition + Compose's own
// healthcheck schedule). A regression that mis-stabilizes on
// `running` alone would return in <500 ms (the very first polling
// iteration).
//
// Fixture: busybox with a script that touches `/tmp/ready` after
// 3 s; the healthcheck `test -f /tmp/ready` flips healthy at that
// point.

package application_test

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	clockadapter "github.com/pt9912/u-boot/internal/adapter/driven/clock"
	dockeradapter "github.com/pt9912/u-boot/internal/adapter/driven/docker"
	fsadapter "github.com/pt9912/u-boot/internal/adapter/driven/fs"
	loggeradapter "github.com/pt9912/u-boot/internal/adapter/driven/logger"
	netprobeadapter "github.com/pt9912/u-boot/internal/adapter/driven/netprobe"
	yamladapter "github.com/pt9912/u-boot/internal/adapter/driven/yaml"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// healthcheckFixture pins the §966 contract. The service:
//   - touches `/tmp/ready` after 3 s, then sleeps for a minute
//     (so it stays `running` after the healthcheck flip);
//   - healthcheck polls every 1 s with `test -f /tmp/ready` —
//     transitions running→healthy at t≈3-4 s.
//
// `start_period: 1s` keeps Compose from marking the service
// failed during the 3 s ramp-up.
const healthcheckFixture = `services:
  slow:
    image: busybox:1.36
    command: ["sh", "-c", "sleep 3 && touch /tmp/ready && sleep 60"]
    healthcheck:
      test: ["CMD", "test", "-f", "/tmp/ready"]
      interval: 1s
      timeout: 1s
      retries: 10
      start_period: 1s
`

// minUbootYAML is the minimum project scaffold UpService.Up needs to
// pass its checkProjectInitialized gate.
const minUbootYAML = `schemaVersion: 1
project:
  name: t-uboot-int
`

func TestUpService_RealDocker_StabilizesOnHealthyNotOnRunning(t *testing.T) {
	// Hard-fail (not skip) on missing docker — see Sub-T1 review
	// rationale (engine_progressstream_docker_test.go).
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH but the test was built with -tags=docker; install docker (e.g. via Sub-T4 Makefile wiring): %v", err)
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "u-boot.yaml"), []byte(minUbootYAML), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "compose.yaml"), []byte(healthcheckFixture), 0o644); err != nil {
		t.Fatalf("seed compose.yaml: %v", err)
	}

	engine := dockeradapter.NewEngine()
	// Register cleanup BEFORE Up so a half-pulled / half-started
	// service is still torn down.
	t.Cleanup(func() {
		dctx, dcancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer dcancel()
		if err := engine.ComposeDown(dctx, dir, driven.ComposeDownOptions{RemoveVolumes: true}); err != nil {
			t.Logf("cleanup ComposeDown failed: %v", err)
		}
	})

	level := new(slog.LevelVar)
	level.Set(slog.LevelWarn) // keep test stderr quiet
	logger := loggeradapter.New(os.Stderr, loggeradapter.FormatText, level)

	svc := application.NewUpService(
		fsadapter.New(),
		yamladapter.New(),
		engine,
		netprobeadapter.New(),
		clockadapter.New(),
		logger,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	start := time.Now()
	resp, err := svc.Up(ctx, driving.UpRequest{
		BaseDir: dir,
		Timeout: 30 * time.Second,
	})
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("Up: %v", err)
	}

	if !resp.Result.Stabilized {
		t.Errorf("Stabilized = false, want true (healthcheck must flip eventually)")
	}

	// Load-bearing pin: a regression that mis-stabilizes on
	// `running` alone would return in <500 ms. Our fixture flips
	// healthy at t≈3-4 s; ≥ 2 s is the conservative threshold
	// that catches the regression while tolerating CI-runner
	// timing slack.
	const minimumWaitForHealthy = 2 * time.Second
	if elapsed < minimumWaitForHealthy {
		t.Errorf(
			"Up returned after %v; want ≥ %v (would mean we stabilized on `running` not `healthy` — LH-FA-UP-001 §966 violation)",
			elapsed, minimumWaitForHealthy,
		)
	}

	// Sanity: the status snapshot should report healthy.
	foundHealthy := false
	for _, s := range resp.Result.Services {
		if s.Name == "slow" && s.Healthcheck == "healthy" {
			foundHealthy = true
		}
	}
	if !foundHealthy {
		t.Errorf("expected service `slow` with Healthcheck=healthy in result; got %+v", resp.Result.Services)
	}
}
