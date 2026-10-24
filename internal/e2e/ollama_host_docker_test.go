//go:build docker

// End-to-end check of the catalogue feature `ollama-host` (LH-FA-DEV-003,
// LH-FA-DEV-006): the container that `u-boot devcontainer up` starts from the
// generated project carries the entry `host.docker.internal:host-gateway` in
// its `ExtraHosts` and resolves the name. The tests make no call to
// `ollama.com`, use no account and no token, and start no service in the
// container: they check the resolution of one name, not the path to a service
// of the host. They run on the engine of the runner (Linux, local engine is
// the measured case) and say nothing about another engine.

package e2e_test

import (
	"context"
	"encoding/json"
	"net"
	"runtime"
	"strings"
	"testing"
)

// dcExtraHosts returns the `HostConfig.ExtraHosts` of a container.
func dcExtraHosts(t *testing.T, id string) []string {
	t.Helper()
	raw := dcMustDocker(t, "inspect", "--format", "{{json .HostConfig.ExtraHosts}}", id)
	var hosts []string
	if err := json.Unmarshal([]byte(raw), &hosts); err != nil {
		t.Fatalf("ExtraHosts of %s: %v (%q)", id, err, raw)
	}
	return hosts
}

// dcResolvesHostGateway reports whether `getent hosts host.docker.internal`
// in the container ends with exit 0 and prints an address first.
func dcResolvesHostGateway(ctx context.Context, t *testing.T, id string) (addr string, exit int) {
	t.Helper()
	out, code := dcExec(ctx, t, id, "vscode", "getent hosts host.docker.internal")
	if code != 0 {
		return "", code
	}
	fields := strings.Fields(out)
	if len(fields) == 0 || net.ParseIP(fields[0]) == nil {
		t.Fatalf("getent hosts host.docker.internal printed %q, want an address first", out)
	}
	return fields[0], 0
}

// D1: default profile. The container started by `u-boot devcontainer up` from
// the generated project holds the entry in `ExtraHosts` and resolves the name
// to an address; the container of the same project without the feature does
// not resolve it (exit 2, the case of the measurement on Linux).
func TestE2E_OllamaHost_DefaultProfileResolvesHostDockerInternal(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, false, map[string]string{"devcontainer.features.ollama-host.enabled": "true"})
	s := dcMustUp(ctx, t, dcService(), dir, false)
	if !dcInspect(t, s.resp.ContainerID).State.Running {
		t.Fatalf("the container is not running: %+v", s.resp)
	}
	if got := dcExtraHosts(t, s.resp.ContainerID); len(got) != 1 || got[0] != "host.docker.internal:host-gateway" {
		t.Errorf("ExtraHosts = %v, want exactly [host.docker.internal:host-gateway]", got)
	}
	addr, code := dcResolvesHostGateway(ctx, t, s.resp.ContainerID)
	if code != 0 {
		t.Fatalf("getent hosts host.docker.internal: exit %d, want 0", code)
	}
	t.Logf("host.docker.internal -> %s", addr)

	// Counter-sample: the same profile without the feature.
	if runtime.GOOS != "linux" {
		t.Logf("counter-sample skipped: measured on Linux only (GOOS=%s)", runtime.GOOS)
		return
	}
	plain := dcGenerateProject(ctx, t, false, nil)
	ps := dcMustUp(ctx, t, dcService(), plain, false)
	if got := dcExtraHosts(t, ps.resp.ContainerID); len(got) != 0 {
		t.Errorf("ExtraHosts without the feature = %v, want none", got)
	}
	if _, code := dcResolvesHostGateway(ctx, t, ps.resp.ContainerID); code != 2 {
		t.Errorf("getent hosts host.docker.internal without the feature: exit %d, want 2 (not found)", code)
	}
}

// D2: sandbox profile with the egress restriction. The container that `u-boot
// devcontainer up` starts (it carries NET_ADMIN and runs the egress script)
// holds the entry and resolves the name through `/etc/hosts`; the allowlist
// baked into the script it runs names neither `ollama.com`, nor
// `host.docker.internal`, nor the address. Not checked: that a service of the
// host is reachable.
func TestE2E_OllamaHost_SandboxWithEgressResolvesAndAddsNoHost(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, true, map[string]string{
		"devcontainer.features.ollama-host.enabled": "true",
		"devcontainer.sandbox.egress.enabled":       "true",
		"devcontainer.sandbox.onUnavailable":        "warn",
	})
	s := dcMustUp(ctx, t, dcService(), dir, false)
	id := s.resp.ContainerID
	if got := dcExtraHosts(t, id); len(got) != 1 || got[0] != "host.docker.internal:host-gateway" {
		t.Errorf("ExtraHosts = %v, want exactly [host.docker.internal:host-gateway]", got)
	}
	addr, code := dcResolvesHostGateway(ctx, t, id)
	if code != 0 {
		t.Fatalf("getent hosts host.docker.internal: exit %d, want 0", code)
	}
	hosts, code := dcExec(ctx, t, id, "root", `grep '^hosts=' /usr/local/bin/u-boot-egress-init`)
	if code != 0 || !strings.HasPrefix(hosts, "hosts=") {
		t.Fatalf("hosts= line of the egress script: exit %d, %q", code, hosts)
	}
	for _, bad := range []string{"ollama.com", "host.docker.internal", addr} {
		if strings.Contains(hosts, bad) {
			t.Errorf("the allowlist of the script names %q: %s", bad, hosts)
		}
	}
	t.Logf("host.docker.internal -> %s; allowlist: %s", addr, strings.TrimSpace(hosts))
}
