//go:build docker

// Shared end-to-end helpers for the LH-AK-* Acceptance-Flow tests
// (slice-v1-keycloak T3 — extracted at the second Acceptance-Docker-
// Test landing point to avoid copy-paste of the init + add + up
// pipeline). The Postgres-spezifische Endpoint-Probe and the
// Keycloak-spezifische HTTP-Probe stay in their respective
// _test.go files; this file only owns the Compose-orchestration
// scaffolding the two flows share.

package e2e_test

import (
	"context"
	"log/slog"
	"net"
	"net/http"
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

// acceptanceFlow describes one end-to-end Acceptance-Flow:
// `u-boot init` + `u-boot add <ServiceName>` + `u-boot up`.
//
// EnvKeys are the `.env.example` keys the test asserts after add
// — the Postgres flow checks `POSTGRES_*`, Keycloak checks
// `KEYCLOAK_*`. Empty list = no env-block assertion.
//
// UpTimeout is passed verbatim into UpRequest.Timeout. Two cases
// in practice:
//   - Postgres ~5 s boot + healthcheck → 90 s is comfortable.
//   - Keycloak 30–90 s JVM-Boot + Realm-Init → 240 s; in CI
//     the cold image pull can push past 180 s.
//
// CtxTimeout bounds the whole-test context (must exceed UpTimeout
// by enough for the Compose-Up + assertions).
type acceptanceFlow struct {
	projectName string
	serviceName string
	envKeys     []string
	upTimeout   time.Duration
	ctxTimeout  time.Duration
}

// acceptanceResult is the harness output the per-flow assertions
// branch on (Healthy-Probe, Port-Dial, HTTP-Probe).
type acceptanceResult struct {
	dir    string
	upResp driving.UpResponse
}

// runAcceptanceFlow wires every production driven adapter (no
// fakes) and drives the three application services in sequence —
// init, add, up. Compose-Down is registered as t.Cleanup BEFORE
// any potentially-failing step so a half-started Compose project
// is still torn down.
//
// Asserts the LH-AK pre-up invariants (compose contains service,
// env file contains every EnvKey). The post-up assertions
// (Stabilized/Healthcheck/Probe) stay in the caller so each flow
// can phrase them in its own LH-AK wording.
func runAcceptanceFlow(t *testing.T, flow acceptanceFlow) acceptanceResult {
	t.Helper()
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH but the test was built with -tags=docker; install docker (e.g. via M6 Sub-T4 Makefile wiring): %v", err)
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

	ctx, cancel := context.WithTimeout(context.Background(), flow.ctxTimeout)
	defer cancel()

	if _, err := initSvc.Init(ctx, driving.InitProjectRequest{
		BaseDir: dir,
		Name:    flow.projectName,
		SkipGit: true,
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	svcName, err := domain.NewServiceName(flow.serviceName)
	if err != nil {
		t.Fatalf("NewServiceName(%s): %v", flow.serviceName, err)
	}
	if _, err := addSvc.Add(ctx, driving.AddServiceRequest{
		BaseDir:     dir,
		ServiceName: svcName,
	}); err != nil {
		t.Fatalf("add %s: %v", flow.serviceName, err)
	}

	// LH-AK pre-up: compose.yaml mentions the service, env.example
	// lists the required keys.
	composeBytes, err := os.ReadFile(filepath.Join(dir, "compose.yaml"))
	if err != nil {
		t.Fatalf("read compose.yaml: %v", err)
	}
	if !strings.Contains(string(composeBytes), flow.serviceName) {
		t.Errorf("compose.yaml does not contain %s service:\n%s", flow.serviceName, composeBytes)
	}
	if len(flow.envKeys) > 0 {
		envBytes, err := os.ReadFile(filepath.Join(dir, ".env.example"))
		if err != nil {
			t.Fatalf("read .env.example: %v", err)
		}
		envStr := string(envBytes)
		for _, key := range flow.envKeys {
			if !strings.Contains(envStr, key) {
				t.Errorf(".env.example missing %s:\n%s", key, envStr)
			}
		}
	}

	upResp, err := upSvc.Up(ctx, driving.UpRequest{
		BaseDir: dir,
		Timeout: flow.upTimeout,
	})
	if err != nil {
		t.Fatalf("up: %v", err)
	}

	return acceptanceResult{dir: dir, upResp: upResp}
}

// dialTCP verifies that addr is reachable; logs the error as
// t.Errorf so the test marks the contract as broken without
// aborting the whole flow (the Postgres / Keycloak ports are
// already healthcheck-pinned by Stabilized).
func dialTCP(ctx context.Context, t *testing.T, addr string) {
	t.Helper()
	var d net.Dialer
	dctx, dcancel := context.WithTimeout(ctx, 5*time.Second)
	defer dcancel()
	conn, err := d.DialContext(dctx, "tcp", addr)
	if err != nil {
		t.Errorf("address %q not reachable: %v", addr, err)
		return
	}
	_ = conn.Close()
}

// probeHTTPEndpoint issues a GET against url and asserts that the
// response status code is in the want list. Used for HTTP-based
// services (Keycloak admin console, future OTel collector front-
// end). Failures are t.Errorf, not t.Fatalf, so cleanup still
// runs.
func probeHTTPEndpoint(ctx context.Context, t *testing.T, url string, want ...int) {
	t.Helper()
	pctx, pcancel := context.WithTimeout(ctx, 10*time.Second)
	defer pcancel()
	req, err := http.NewRequestWithContext(pctx, http.MethodGet, url, nil)
	if err != nil {
		t.Errorf("build HTTP probe request for %s: %v", url, err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("HTTP probe to %s failed: %v", url, err)
		return
	}
	defer func() { _ = resp.Body.Close() }()
	for _, code := range want {
		if resp.StatusCode == code {
			return
		}
	}
	t.Errorf("HTTP probe %s returned status %d; want one of %v", url, resp.StatusCode, want)
}

// stabilizationCheck pins the LH-AK post-up invariants: Stabilized
// flag set, the named service reaches healthcheck `healthy`. Used
// by every Acceptance-Flow because the upservice contract is
// service-independent.
func stabilizationCheck(t *testing.T, res acceptanceResult, serviceName string) {
	t.Helper()
	if !res.upResp.Result.Stabilized {
		t.Errorf("Stabilized = false, want true")
	}
	for _, s := range res.upResp.Result.Services {
		if s.Name == serviceName && s.Healthcheck == "healthy" {
			return
		}
	}
	t.Errorf("%s did not reach `healthy`; services snapshot: %+v", serviceName, res.upResp.Result.Services)
}
