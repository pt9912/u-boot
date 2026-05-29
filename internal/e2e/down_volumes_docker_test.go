//go:build docker

// LH-FA-UP-004 §1015 down --volumes-removes-named-volumes pin
// (M6-docker-int Sub-T3).
//
// Spec §1015: "Das Produkt muss zwischen einem regulären Stopp
// (Container stoppen) und einem vollständigen Aufräumen (Container
// und Volumes entfernen) unterscheiden". This test pins the
// `--volumes` half of that contract: after `down --volumes`, every
// named volume created by `compose up` must be gone at the Docker
// engine level.
//
// Approach: snapshot the host's `docker volume ls` BEFORE compose
// up, after up, and after down --volumes. Volumes present after up
// but absent after down must include every named volume the
// fixture declares — proves the Down-service actually delegated
// `-v` rather than dropping it.

package e2e_test

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	confirmadapter "github.com/pt9912/u-boot/internal/adapter/driven/confirm"
	dockeradapter "github.com/pt9912/u-boot/internal/adapter/driven/docker"
	fsadapter "github.com/pt9912/u-boot/internal/adapter/driven/fs"
	loggeradapter "github.com/pt9912/u-boot/internal/adapter/driven/logger"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// downVolumesFixture pins a fixture with two named volumes so the
// pin proves both were dropped (not just one). busybox+sleep keeps
// the test small and fast.
const downVolumesFixture = `services:
  data:
    image: busybox:1.36
    command: ["sleep", "60"]
    volumes:
      - aaa:/aaa
      - bbb:/bbb
volumes:
  aaa: {}
  bbb: {}
`

const downVolumesUbootYAML = `schemaVersion: 1
project:
  name: t-uboot-int
`

// listHostVolumes snapshots `docker volume ls -q` into a set.
func listHostVolumes(t *testing.T) map[string]struct{} {
	t.Helper()
	out, err := exec.Command("docker", "volume", "ls", "-q").Output()
	if err != nil {
		t.Fatalf("docker volume ls: %v", err)
	}
	set := make(map[string]struct{})
	for _, name := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if name != "" {
			set[name] = struct{}{}
		}
	}
	return set
}

func TestE2E_LHFAUP004_DownVolumesRemovesNamedVolumes(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH but the test was built with -tags=docker; install docker (e.g. via Sub-T4 Makefile wiring): %v", err)
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "u-boot.yaml"), []byte(downVolumesUbootYAML), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "compose.yaml"), []byte(downVolumesFixture), 0o644); err != nil {
		t.Fatalf("seed compose.yaml: %v", err)
	}

	fs := fsadapter.New()
	engine := dockeradapter.NewEngine()
	conf := confirmadapter.New(strings.NewReader(""), os.Stderr)
	level := new(slog.LevelVar)
	level.Set(slog.LevelWarn)
	logger := loggeradapter.New(os.Stderr, loggeradapter.FormatText, level)

	// Defense-in-depth cleanup: even after the test's own Down,
	// trigger one more Down --volumes to flush any leak from a
	// half-completed run.
	t.Cleanup(func() {
		dctx, dcancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer dcancel()
		_ = engine.ComposeDown(dctx, dir, driven.ComposeDownOptions{RemoveVolumes: true})
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	before := listHostVolumes(t)

	if _, err := engine.ComposeUp(ctx, dir, driven.ComposeUpOptions{Detach: true}); err != nil {
		t.Fatalf("ComposeUp: %v", err)
	}

	afterUp := listHostVolumes(t)
	created := diffVolumeSets(afterUp, before)
	if len(created) < 2 {
		t.Fatalf("expected at least 2 new named volumes after up, got %d: %v", len(created), created)
	}

	downSvc := application.NewDownService(fs, engine, conf, logger)
	if _, err := downSvc.Down(ctx, driving.DownRequest{
		BaseDir:       dir,
		RemoveVolumes: true,
		AssumeYes:     true,
	}); err != nil {
		t.Fatalf("Down --volumes --yes: %v", err)
	}

	afterDown := listHostVolumes(t)
	for vol := range created {
		if _, present := afterDown[vol]; present {
			t.Errorf("volume %q still present after down --volumes (LH-FA-UP-004 §1015 violation)", vol)
		}
	}
}

// diffVolumeSets returns the set of volumes in `after` that are not
// in `before`.
func diffVolumeSets(after, before map[string]struct{}) map[string]struct{} {
	out := make(map[string]struct{})
	for v := range after {
		if _, ok := before[v]; !ok {
			out[v] = struct{}{}
		}
	}
	return out
}
