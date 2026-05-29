//go:build docker

// LH-FA-DIAG-002 Compose-ps JSON-schema pin (M6-docker-int Sub-T1).
//
// Two complementary tests:
//
//  1. Engine-level: ComposePs against a real `up`'d service must
//     populate the [driven.ComposeService] fields the application
//     reads (Name, State; Health and Ports may be empty depending
//     on the fixture).
//
//  2. Raw-JSON: bypass the T2 parser entirely and assert the
//     expected JSON field NAMES (`Service`, `State`, `Health`,
//     `Publishers`) appear in `docker compose ps --format json`
//     output. If a Compose release renamed any of these the T2
//     parser would silently drop data — this pin catches the
//     schema drift directly.

package docker_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/adapter/driven/docker"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// composePsFixture pins a tiny long-running image with both a port
// mapping and a healthcheck so the LH-FA-DIAG-002 schema pin can
// hard-assert that Publishers AND Health appear in the JSON output
// — the two fields the T2 parser depends on for port-probe gating
// and stabilization classification. A pure busybox+sleep fixture
// would leave both as soft-log assertions and let a schema rename
// of either field slip past CI.
//
// `127.0.0.1:0:1` binds container port 1 to a kernel-assigned host
// port on the loopback interface (no collision risk between
// parallel test runs). The healthcheck runs `true` — trivially
// passes after the first 1 s interval so Compose can report
// `healthy` quickly.
const composePsFixture = `services:
  busy:
    image: busybox:1.36
    command: ["sleep", "60"]
    ports:
      - "127.0.0.1:0:1"
    healthcheck:
      test: ["CMD", "true"]
      interval: 1s
      timeout: 1s
      retries: 1
`

func setupComposePsFixture(t *testing.T) (engine *docker.Engine, dir string) {
	t.Helper()
	// Hard-fail (not skip) on missing docker — see the equivalent
	// note in engine_progressstream_docker_test.go.
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH but the test was built with -tags=docker; install docker (e.g. via Sub-T4 Makefile wiring): %v", err)
	}
	dir = t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "compose.yaml"), []byte(composePsFixture), 0o644); err != nil {
		t.Fatalf("write compose.yaml: %v", err)
	}
	engine = docker.NewEngine()

	// Register cleanup BEFORE ComposeUp so even a partial pull /
	// half-started network is torn down. ComposeDown on a not-yet-
	// up'd project is mostly a no-op; the t.Cleanup hook fires on
	// any test exit (success or t.Fatalf later in the test).
	t.Cleanup(func() {
		dctx, dcancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer dcancel()
		if err := engine.ComposeDown(dctx, dir, driven.ComposeDownOptions{RemoveVolumes: true}); err != nil {
			t.Logf("cleanup ComposeDown failed (may be expected if Up never completed): %v", err)
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if _, err := engine.ComposeUp(ctx, dir, driven.ComposeUpOptions{Detach: true}); err != nil {
		t.Fatalf("ComposeUp: %v", err)
	}
	return engine, dir
}

func TestEngine_RealDocker_ComposePs_PopulatesParserFields(t *testing.T) {
	engine, dir := setupComposePsFixture(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	services, err := engine.ComposePs(ctx, dir)
	if err != nil {
		t.Fatalf("ComposePs: %v", err)
	}
	if len(services) == 0 {
		t.Fatal("ComposePs returned no services; expected at least the `busy` fixture")
	}

	var busy *driven.ComposeService
	for i := range services {
		if services[i].Name == "busy" {
			busy = &services[i]
			break
		}
	}
	if busy == nil {
		t.Fatalf("expected service `busy` in ps output, got: %+v", services)
	}
	if busy.State == "" {
		t.Errorf("Service %q has empty State; expected a Compose state string", busy.Name)
	}
	// Health and Ports may legitimately be empty for the
	// healthcheck-less / no-ports fixture — that's a
	// behaviour assertion, not a schema gap.
}

func TestEngine_RealDocker_ComposePs_RawJSONHasExpectedFieldNames(t *testing.T) {
	_, dir := setupComposePsFixture(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, "docker", "compose",
		"-f", filepath.Join(dir, "compose.yaml"),
		"ps", "--format", "json").Output()
	if err != nil {
		t.Fatalf("docker compose ps: %v", err)
	}
	raw := bytes.TrimSpace(out)
	if len(raw) == 0 {
		t.Fatal("empty output from docker compose ps")
	}

	// Compose v2.20: NDJSON, one object per line.
	// Compose v2.21+: JSON array.
	// Parse either by inspecting the first byte.
	var first map[string]any
	if raw[0] == '[' {
		var arr []map[string]any
		if err := json.Unmarshal(raw, &arr); err != nil {
			t.Fatalf("parse JSON array: %v", err)
		}
		if len(arr) == 0 {
			t.Fatal("docker compose ps returned an empty JSON array")
		}
		first = arr[0]
	} else {
		idx := bytes.IndexByte(raw, '\n')
		line := raw
		if idx >= 0 {
			line = raw[:idx]
		}
		if err := json.Unmarshal(line, &first); err != nil {
			t.Fatalf("parse NDJSON line: %v", err)
		}
	}

	// Hard-pin every field name the T2 parser reads. A drift in
	// Compose's JSON schema (rename, removal) makes this test fail
	// fast — bringing the parser's drift-detection forward to CI
	// rather than to a user-reported missing field at runtime.
	//
	// Publishers and Health are pinned hard (not soft-logged) since
	// the fixture above declares both a port and a healthcheck —
	// the most behavior-relevant fields for the polling-loop gating
	// must drift-fail this test, not silently shape-shift.
	requiredFields := []string{"Service", "State", "Health", "Publishers"}
	for _, key := range requiredFields {
		if _, ok := first[key]; !ok {
			t.Errorf("docker compose ps JSON missing required field %q; row: %#v", key, first)
		}
	}
}
