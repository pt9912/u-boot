//go:build docker

// End-to-end check of `devcontainer.resources.memory` and
// `devcontainer.resources.cpus` (LH-FA-CONF-002, LH-FA-DEV-001): the container
// that `u-boot devcontainer up` starts from the generated project carries the
// limits in its `HostConfig`, in the standard profile and in the sandbox
// profile with `nestedRuntime: docker`. The tests need an engine that applies
// limits: a Linux host (`amd64`) with a local engine that is not nested; on an
// engine that cannot apply them (the nested engine of the devcontainer of this
// repository stops every start with `--memory` at the cgroup) they fail at the
// start, and they say nothing about another engine.

package e2e_test

import (
	"context"
	"strings"
	"testing"
)

const (
	dcTestMemory      = "512m"
	dcTestMemoryBytes = "536870912"
	dcTestCPUs        = "1.5"
	dcTestNanoCPUs    = "1500000000"
)

// dcLimits returns `HostConfig.Memory` and `HostConfig.NanoCpus` of a container.
func dcLimits(t *testing.T, id string) (memory, nanoCPUs string) {
	t.Helper()
	out := dcMustDocker(t, "inspect", "--format", "{{.HostConfig.Memory}} {{.HostConfig.NanoCpus}}", id)
	fields := strings.Fields(out)
	if len(fields) != 2 {
		t.Fatalf("HostConfig.Memory and NanoCpus of %s: %q", id, out)
	}
	return fields[0], fields[1]
}

// L1: standard profile. The started container holds both limits; a project
// without the keys holds none (counter-sample).
func TestE2E_Resources_DefaultProfileLimitsReachTheContainer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, false, map[string]string{
		"devcontainer.resources.memory": dcTestMemory,
		"devcontainer.resources.cpus":   dcTestCPUs,
	})
	s := dcMustUp(ctx, t, dcService(), dir, false)
	if !dcInspect(t, s.resp.ContainerID).State.Running {
		t.Fatalf("the container is not running: %+v", s.resp)
	}
	if mem, cpus := dcLimits(t, s.resp.ContainerID); mem != dcTestMemoryBytes || cpus != dcTestNanoCPUs {
		t.Errorf("Memory = %s, NanoCpus = %s, want %s and %s", mem, cpus, dcTestMemoryBytes, dcTestNanoCPUs)
	}

	plain := dcGenerateProject(ctx, t, false, nil)
	ps := dcMustUp(ctx, t, dcService(), plain, false)
	if mem, cpus := dcLimits(t, ps.resp.ContainerID); mem != "0" || cpus != "0" {
		t.Errorf("without the keys Memory = %s, NanoCpus = %s, want 0 and 0", mem, cpus)
	}
}

// L2: sandbox profile with `nestedRuntime: docker` (`--privileged`). The
// started container holds both limits next to the privileged mode. Not
// checked: the effect on the nested Docker and its containers.
func TestE2E_Resources_SandboxNestedDockerLimitsReachTheContainer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, true, map[string]string{
		"devcontainer.sandbox.nestedRuntime": "docker",
		"devcontainer.sandbox.onUnavailable": "warn",
		"devcontainer.resources.memory":      dcTestMemory,
		"devcontainer.resources.cpus":        dcTestCPUs,
	})
	s := dcMustUp(ctx, t, dcService(), dir, false)
	id := s.resp.ContainerID
	if privileged := strings.TrimSpace(dcMustDocker(t, "inspect", "--format", "{{.HostConfig.Privileged}}", id)); privileged != "true" {
		t.Errorf("Privileged = %s, want true", privileged)
	}
	if mem, cpus := dcLimits(t, id); mem != dcTestMemoryBytes || cpus != dcTestNanoCPUs {
		t.Errorf("Memory = %s, NanoCpus = %s, want %s and %s", mem, cpus, dcTestMemoryBytes, dcTestNanoCPUs)
	}
}
