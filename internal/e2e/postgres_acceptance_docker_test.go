//go:build docker

// LH-AK-002 PostgreSQL-Acceptance-Flow Pin (M6-docker-int Sub-T3).
//
// Spec §LH-AK-002: the sequence
//
//   u-boot init
//   u-boot add postgres
//   u-boot up
//
// must succeed end-to-end. Acceptance criteria:
//
//   - PostgreSQL-Service exists in `compose.yaml`;
//   - `.env.example` lists `POSTGRES_USER`, `POSTGRES_PASSWORD`,
//     `POSTGRES_DB`;
//   - the container reaches healthcheck status `healthy` within 60 s;
//   - port 5432 is reachable on `localhost`.
//
// The test wires every production driven adapter (no fakes) and
// drives the three application services in sequence.
//
// Test prerequisites (slice §Strukturelle Bedingungen): docker CLI
// + compose plugin available; test process and docker daemon must
// share a network namespace. `make test-docker` satisfies both via
// the `test-docker-tools` Dockerfile stage and `--network=host`.

package e2e_test

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	clockadapter "github.com/pt9912/u-boot/internal/adapter/driven/clock"
	confirmadapter "github.com/pt9912/u-boot/internal/adapter/driven/confirm"
	dockeradapter "github.com/pt9912/u-boot/internal/adapter/driven/docker"
	fsadapter "github.com/pt9912/u-boot/internal/adapter/driven/fs"
	gitadapter "github.com/pt9912/u-boot/internal/adapter/driven/git"
	loggeradapter "github.com/pt9912/u-boot/internal/adapter/driven/logger"
	netprobeadapter "github.com/pt9912/u-boot/internal/adapter/driven/netprobe"
	progressadapter "github.com/pt9912/u-boot/internal/adapter/driven/progress"
	yamladapter "github.com/pt9912/u-boot/internal/adapter/driven/yaml"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

func TestE2E_LHAK002_PostgresAcceptanceFlow(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH but the test was built with -tags=docker; install docker (e.g. via Sub-T4 Makefile wiring): %v", err)
	}

	dir := t.TempDir()

	fs := fsadapter.New()
	yaml := yamladapter.New()
	git := gitadapter.New()
	prog := progressadapter.NewText(os.Stderr)
	conf := confirmadapter.New(strings.NewReader(""), os.Stderr)
	level := new(slog.LevelVar)
	level.Set(slog.LevelWarn)
	logger := loggeradapter.New(os.Stderr, loggeradapter.FormatText, level)
	engine := dockeradapter.NewEngine()
	netprobe := netprobeadapter.New()
	clock := clockadapter.New()

	// Cleanup registered BEFORE any potentially-failing step so a
	// half-started Compose project is still torn down. ComposeDown
	// is a no-op when nothing is up.
	t.Cleanup(func() {
		dctx, dcancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer dcancel()
		if err := engine.ComposeDown(dctx, dir, driven.ComposeDownOptions{RemoveVolumes: true}); err != nil {
			t.Logf("cleanup ComposeDown failed: %v", err)
		}
	})

	initSvc := application.NewInitProjectService(fs, yaml, git, prog, conf, logger)
	addSvc := application.NewAddServiceService(fs, yaml, conf, logger)
	upSvc := application.NewUpService(fs, yaml, engine, netprobe, clock, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Step 1: u-boot init (SkipGit so the test env stays clean).
	// Name is set explicitly because t.TempDir() returns a numeric
	// counter as the leaf segment (`001`/`002`/…), which fails the
	// LH-FA-INIT-006 regex `^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$` when
	// the service falls back to filepath.Base(dir).
	if _, err := initSvc.Init(ctx, driving.InitProjectRequest{
		BaseDir: dir,
		Name:    "t-uboot-e2e-acc",
		SkipGit: true,
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	// Step 2: u-boot add postgres.
	serviceName, err := domain.NewServiceName("postgres")
	if err != nil {
		t.Fatalf("NewServiceName(postgres): %v", err)
	}
	if _, err := addSvc.Add(ctx, driving.AddServiceRequest{
		BaseDir:     dir,
		ServiceName: serviceName,
	}); err != nil {
		t.Fatalf("add postgres: %v", err)
	}

	// LH-AK-002 (1): compose.yaml contains the postgres service.
	composeBytes, err := os.ReadFile(filepath.Join(dir, "compose.yaml"))
	if err != nil {
		t.Fatalf("read compose.yaml: %v", err)
	}
	if !strings.Contains(string(composeBytes), "postgres") {
		t.Errorf("compose.yaml does not contain postgres service:\n%s", composeBytes)
	}

	// LH-AK-002 (2): .env.example lists the three POSTGRES_* variables.
	envBytes, err := os.ReadFile(filepath.Join(dir, ".env.example"))
	if err != nil {
		t.Fatalf("read .env.example: %v", err)
	}
	envStr := string(envBytes)
	for _, key := range []string{"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB"} {
		if !strings.Contains(envStr, key) {
			t.Errorf(".env.example missing %s:\n%s", key, envStr)
		}
	}

	// Step 3: u-boot up.
	upResp, err := upSvc.Up(ctx, driving.UpRequest{
		BaseDir: dir,
		Timeout: 90 * time.Second,
	})
	if err != nil {
		t.Fatalf("up: %v", err)
	}

	// LH-AK-002 (3): container reaches `healthy`.
	if !upResp.Result.Stabilized {
		t.Errorf("Stabilized = false, want true")
	}
	foundHealthy := false
	for _, s := range upResp.Result.Services {
		if s.Name == "postgres" && s.Healthcheck == "healthy" {
			foundHealthy = true
		}
	}
	if !foundHealthy {
		t.Errorf("postgres did not reach `healthy`; services snapshot: %+v", upResp.Result.Services)
	}

	// LH-AK-002 (4): port 5432 reachable on localhost. With a
	// healthcheck UpService treats port-probe failure as warn-only,
	// so the assertion above does not by itself prove the port is
	// reachable — we dial directly here to lock that part of the
	// contract.
	var d net.Dialer
	dialCtx, dialCancel := context.WithTimeout(ctx, 5*time.Second)
	defer dialCancel()
	conn, err := d.DialContext(dialCtx, "tcp", "localhost:5432")
	if err != nil {
		t.Errorf("port 5432 not reachable on localhost (LH-AK-002 §2321): %v", err)
	} else {
		_ = conn.Close()
	}
}
