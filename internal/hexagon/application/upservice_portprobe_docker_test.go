//go:build docker

// LH-FA-UP-001 §968 TCP-port-probe-lands pin
// (M6-docker-int Sub-T2).
//
// Spec §968: "Bei definierten Ports wird auf Erreichbarkeit auf
// `localhost` geprüft, sofern es sich um TCP-basierten Zugriff
// handelt." For a service WITHOUT a healthcheck, the port probe
// is load-bearing for stabilization classification — slice plan
// §141 "no healthcheck + TCP port → port probe gates".
//
// Pin shape: wrap the production NetProbe in a small spy that
// counts DialTCP calls plus forwards them to the real adapter.
// After UpService.Up stabilizes, assert:
//   - the spy recorded at least one call;
//   - the call was against `localhost` + the host-mapped port
//     declared in compose.yaml.
//
// A regression that mis-classifies (e.g. stabilizes on `running`
// without probing) would leave spy.calls() empty.
//
// Fixture: nginx-alpine (small, fast pull) exposing port 80 with an
// explicit host port 18080 on the loopback interface. nginx
// listens immediately on container start, so the probe returns
// nil-error during the first or second polling iteration.

package application_test

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
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

const portProbeFixture = `services:
  web:
    image: nginx:alpine
    ports:
      - "127.0.0.1:18080:80"
`

// spyingNetProbe records every DialTCP call and forwards to a real
// netprobe adapter. Used to assert that UpService actually probed
// the declared TCP port rather than mis-classifying on running-
// state alone.
type spyingNetProbe struct {
	delegate driven.NetProbe
	mu       sync.Mutex
	calls    []spyDialCall
}

type spyDialCall struct {
	Host string
	Port int
}

func (s *spyingNetProbe) DialTCP(ctx context.Context, host string, port int, timeout time.Duration) error {
	s.mu.Lock()
	s.calls = append(s.calls, spyDialCall{Host: host, Port: port})
	s.mu.Unlock()
	return s.delegate.DialTCP(ctx, host, port, timeout)
}

func (s *spyingNetProbe) snapshot() []spyDialCall {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]spyDialCall, len(s.calls))
	copy(out, s.calls)
	return out
}

func TestUpService_RealDocker_PortProbeRunsForNoHealthcheckService(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH but the test was built with -tags=docker; install docker (e.g. via Sub-T4 Makefile wiring): %v", err)
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "u-boot.yaml"), []byte(minUbootYAML), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "compose.yaml"), []byte(portProbeFixture), 0o644); err != nil {
		t.Fatalf("seed compose.yaml: %v", err)
	}

	engine := dockeradapter.NewEngine()
	t.Cleanup(func() {
		dctx, dcancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer dcancel()
		if err := engine.ComposeDown(dctx, dir, driven.ComposeDownOptions{RemoveVolumes: true}); err != nil {
			t.Logf("cleanup ComposeDown failed: %v", err)
		}
	})

	level := new(slog.LevelVar)
	level.Set(slog.LevelWarn)
	logger := loggeradapter.New(os.Stderr, loggeradapter.FormatText, level)

	spy := &spyingNetProbe{delegate: netprobeadapter.New()}

	svc := application.NewUpService(
		fsadapter.New(),
		yamladapter.New(),
		engine,
		spy,
		clockadapter.New(),
		logger,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := svc.Up(ctx, driving.UpRequest{
		BaseDir: dir,
		Timeout: 30 * time.Second,
	})
	if err != nil {
		t.Fatalf("Up: %v", err)
	}

	if !resp.Result.Stabilized {
		t.Errorf("Stabilized = false, want true (nginx listens immediately on container port 80)")
	}

	// Load-bearing pin: the probe MUST have been called. A
	// regression that mis-classifies (stabilizing on `running` for
	// a no-healthcheck service without probing) leaves spy.calls
	// empty.
	calls := spy.snapshot()
	if len(calls) == 0 {
		t.Fatal("NetProbe.DialTCP was never called — LH-FA-UP-001 §968 violation: TCP-port-probe must run for services with declared TCP ports")
	}

	// Sanity: at least one call hit the declared host port on
	// loopback. The probe target is normalized by parseComposePort
	// to `localhost:18080`.
	foundExpected := false
	for _, c := range calls {
		if c.Port == 18080 && (c.Host == "localhost" || c.Host == "127.0.0.1") {
			foundExpected = true
			break
		}
	}
	if !foundExpected {
		t.Errorf("expected NetProbe call against localhost:18080 (declared compose port); got %+v", calls)
	}
}
