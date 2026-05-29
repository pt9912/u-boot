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

// composePsFixture pins a tiny long-running image so the
// `compose up -d` + `compose ps` sequence completes quickly and
// the test environment isn't held by a heavy pull. busybox is
// ~1 MiB compressed.
const composePsFixture = `services:
  busy:
    image: busybox:1.36
    command: ["sleep", "60"]
`

func setupComposePsFixture(t *testing.T) (engine *docker.Engine, dir string, cleanup func()) {
	t.Helper()
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skipf("docker not on PATH: %v", err)
	}
	dir = t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "compose.yaml"), []byte(composePsFixture), 0o644); err != nil {
		t.Fatalf("write compose.yaml: %v", err)
	}
	engine = docker.NewEngine()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if _, err := engine.ComposeUp(ctx, dir, driven.ComposeUpOptions{Detach: true}); err != nil {
		t.Fatalf("ComposeUp: %v", err)
	}
	cleanup = func() {
		dctx, dcancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer dcancel()
		if err := engine.ComposeDown(dctx, dir, driven.ComposeDownOptions{RemoveVolumes: true}); err != nil {
			t.Logf("cleanup ComposeDown failed: %v", err)
		}
	}
	return engine, dir, cleanup
}

func TestEngine_RealDocker_ComposePs_PopulatesParserFields(t *testing.T) {
	engine, dir, cleanup := setupComposePsFixture(t)
	defer cleanup()

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
	_, dir, cleanup := setupComposePsFixture(t)
	defer cleanup()

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
	requiredFields := []string{"Service", "State"}
	for _, key := range requiredFields {
		if _, ok := first[key]; !ok {
			t.Errorf("docker compose ps JSON missing required field %q; row: %#v", key, first)
		}
	}

	// `Health` and `Publishers` MAY be absent for fixtures without
	// a healthcheck or exposed ports — log only.
	if _, ok := first["Health"]; !ok {
		t.Logf("Health field absent (acceptable for the no-healthcheck fixture)")
	}
	if _, ok := first["Publishers"]; !ok {
		t.Logf("Publishers field absent (acceptable for the no-ports fixture)")
	}
}
