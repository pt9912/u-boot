//go:build docker

package e2e_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// testNameBudget bounds one test of this file.
const testNameBudget = 10 * time.Minute

// LH-FA-DEV-010: the started devcontainer carries the name
// `u-boot-<project>` beside its ID; the labels and the reuse of the container
// are unchanged, `--recreate` keeps the name on the new container, and the
// response reports the name beside the container ID. The generated
// devcontainer.json is byte-identical across the calls: the name travels in a
// per-call copy, never in the file of the project.
func TestDevcontainerUp_Docker_ContainerName(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testNameBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, false, nil)
	svc := dcService()

	configBefore, err := os.ReadFile(filepath.Join(dir, ".devcontainer", "devcontainer.json"))
	if err != nil {
		t.Fatal(err)
	}
	s := dcMustUp(ctx, t, svc, dir, false)
	want := "u-boot-dcstart"
	if s.resp.ContainerName != want {
		t.Errorf("response name = %q, want %q", s.resp.ContainerName, want)
	}
	doc := dcInspect(t, s.resp.ContainerID)
	if doc.Name != "/"+want {
		t.Errorf("container Name = %q, want %q (the engine holds the name)", doc.Name, "/"+want)
	}
	if got := doc.Config.Labels[dcFolderLbl]; got != dir {
		t.Errorf("label %s = %q, want the project folder (the search keeps its labels)", dcFolderLbl, got)
	}

	// The CLI takes the container of the folder over (the holder of the same
	// folder belongs to this start).
	again := dcMustUp(ctx, t, svc, dir, false)
	if again.resp.ContainerID != s.resp.ContainerID {
		t.Errorf("up again = %s, want the same container (the name is a reuse key of the CLI)", again.resp.ContainerID[:12])
	}

	// --recreate creates the container anew and keeps the name.
	re := dcMustUp(ctx, t, svc, dir, true)
	if re.resp.ContainerID == s.resp.ContainerID {
		t.Errorf("--recreate did not create a new container")
	}
	if re.resp.ContainerName != want {
		t.Errorf("--recreate response name = %q, want %q", re.resp.ContainerName, want)
	}
	if doc2 := dcInspect(t, re.resp.ContainerID); doc2.Name != "/"+want {
		t.Errorf("the new container Name = %q, want %q", doc2.Name, "/"+want)
	}

	configAfter, err := os.ReadFile(filepath.Join(dir, ".devcontainer", "devcontainer.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(configBefore, configAfter) {
		t.Errorf("the generated devcontainer.json changed across the calls (before %d bytes, after %d bytes)",
			len(configBefore), len(configAfter))
	}
}

// LH-FA-DEV-010: a container of another project folder that holds the same
// name (here two projects whose u-boot.yaml carries the same project name)
// ends the start with the name conflict (exit code 10), and the holder is
// neither stopped nor removed.
func TestDevcontainerUp_Docker_ConflictingContainerName(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testNameBudget)
	defer cancel()
	holderDir := dcGenerateProject(ctx, t, false, nil)
	svc := dcService()
	holder := dcMustUp(ctx, t, svc, holderDir, false)

	secondDir := dcGenerateProject(ctx, t, false, nil)
	s2 := dcUp(ctx, svc, secondDir, false)
	if !errors.Is(s2.err, driving.ErrDevcontainerContainerNameConflict) || cli.ExitCode(s2.err) != 10 {
		t.Fatalf("up of the second project: exit %d, err=%v, want the name conflict with exit 10", cli.ExitCode(s2.err), s2.err)
	}
	if doc := dcInspect(t, holder.resp.ContainerID); !doc.State.Running {
		t.Errorf("the holder container does not run: %+v", doc.State)
	}
}
