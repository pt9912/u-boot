//go:build docker

// End-to-end check of the Lastenheft 0.3.x sandbox devcontainer
// (LH-FA-DEV-004/-006/-007/-009): generate the files with the
// production adapters, build the image with the real Docker engine,
// start it with the generated runArgs and run the generated startup
// script plus a nested `podman run` and `podman build` (ADR 0014,
// ADR 0021) or an own Docker daemon (ADR 0018). Podman as the host
// engine is not covered; Colima (Linux VM, cgroup v2) was measured for
// ADR 0018. Hosts with kernel.apparmor_restrict_unprivileged_userns=1
// skip the nested-Podman tests.

package e2e_test

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	confirmadapter "github.com/pt9912/u-boot/internal/adapter/driven/confirm"
	fsadapter "github.com/pt9912/u-boot/internal/adapter/driven/fs"
	gitadapter "github.com/pt9912/u-boot/internal/adapter/driven/git"
	loggeradapter "github.com/pt9912/u-boot/internal/adapter/driven/logger"
	progressadapter "github.com/pt9912/u-boot/internal/adapter/driven/progress"
	yamladapter "github.com/pt9912/u-boot/internal/adapter/driven/yaml"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const sandboxCanaryToken = "ghp_canaryCANARYcanary0123456789abcdef"

func dockerOutput(ctx context.Context, t *testing.T, args ...string) (string, int) {
	t.Helper()
	out, err := exec.CommandContext(ctx, "docker", args...).CombinedOutput()
	code := 0
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("docker %v: %v", args, err)
		}
		code = ee.ExitCode()
	}
	return string(out), code
}

// generateSandboxProject produces a sandbox project with the given
// nested runtime (`podman` | `docker`) in a temp dir and returns the dir.
func generateSandboxProject(ctx context.Context, t *testing.T, runtime, onUnavailable string) string {
	t.Helper()
	return generateSandboxProjectWith(ctx, t, map[string]string{
		"devcontainer.sandbox.nestedRuntime": runtime,
		"devcontainer.sandbox.onUnavailable": onUnavailable,
	})
}

// generateSandboxProjectWith produces a sandbox project with the given
// `devcontainer.sandbox.*` settings and returns the dir.
func generateSandboxProjectWith(ctx context.Context, t *testing.T, settings map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	fsys := fsadapter.New()
	yaml := yamladapter.New()
	logger := loggeradapter.New(os.Stderr, loggeradapter.FormatText, nil)
	initSvc := application.NewInitProjectService(fsys, yaml, gitadapter.New(),
		progressadapter.NewText(os.Stderr), confirmadapter.New(strings.NewReader(""), os.Stderr), logger)
	if _, err := initSvc.Init(ctx, driving.InitProjectRequest{
		BaseDir: dir, Name: "sbxe2e", SkipGit: true, Devcontainer: true, Sandbox: true, SilenceProgress: true,
	}); err != nil {
		t.Fatalf("init --devcontainer --sandbox: %v", err)
	}
	cfgSvc := application.NewConfigService(fsys, yaml, logger)
	for path, value := range settings {
		p, err := domain.NewConfigPath(path)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := cfgSvc.Set(ctx, driving.ConfigSetRequest{BaseDir: dir, Path: p, Value: value}); err != nil {
			t.Fatalf("config set %s: %v", path, err)
		}
	}
	genSvc := application.NewGenerateService(fsys, yaml, logger)
	if _, err := genSvc.Generate(ctx, driving.GenerateRequest{BaseDir: dir, Artifact: domain.ArtifactDevcontainer}); err != nil {
		t.Fatalf("generate devcontainer: %v", err)
	}
	return dir
}

// generateEgressProject renders a sandbox project with the egress
// restriction enabled (no nested Podman) and returns the dir.
func generateEgressProject(ctx context.Context, t *testing.T, onUnavailable string) string {
	t.Helper()
	dir := t.TempDir()
	fsys := fsadapter.New()
	yaml := yamladapter.New()
	logger := loggeradapter.New(os.Stderr, loggeradapter.FormatText, nil)
	initSvc := application.NewInitProjectService(fsys, yaml, gitadapter.New(),
		progressadapter.NewText(os.Stderr), confirmadapter.New(strings.NewReader(""), os.Stderr), logger)
	if _, err := initSvc.Init(ctx, driving.InitProjectRequest{
		BaseDir: dir, Name: "egrese2e", SkipGit: true, Devcontainer: true, Sandbox: true, SilenceProgress: true,
	}); err != nil {
		t.Fatalf("init --devcontainer --sandbox: %v", err)
	}
	cfgSvc := application.NewConfigService(fsys, yaml, logger)
	for path, value := range map[string]string{
		"devcontainer.sandbox.egress.enabled": "true",
		"devcontainer.sandbox.onUnavailable":  onUnavailable,
	} {
		p, err := domain.NewConfigPath(path)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := cfgSvc.Set(ctx, driving.ConfigSetRequest{BaseDir: dir, Path: p, Value: value}); err != nil {
			t.Fatalf("config set %s: %v", path, err)
		}
	}
	genSvc := application.NewGenerateService(fsys, yaml, logger)
	if _, err := genSvc.Generate(ctx, driving.GenerateRequest{BaseDir: dir, Artifact: domain.ArtifactDevcontainer}); err != nil {
		t.Fatalf("generate devcontainer: %v", err)
	}
	return dir
}

// LH-FA-DEV-008: the generated egress script blocks everything but the
// allowlist; without NET_ADMIN it degrades per onUnavailable.
func TestE2E_SandboxDevcontainer_Egress(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	build := func(dir, tag string) string {
		image := tag + ":" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
		if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
			t.Fatalf("docker build failed:\n%s", out)
		}
		t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })
		return image
	}
	warnDir := generateEgressProject(ctx, t, "warn")
	image := build(warnDir, "uboot-e2e-egress")
	runArgs := generatedRunArgs(t, warnDir)
	if len(runArgs) != 1 || runArgs[0] != "--cap-add=NET_ADMIN" {
		t.Fatalf("runArgs = %v, want only --cap-add=NET_ADMIN", runArgs)
	}

	probe := "sudo sh /usr/local/bin/u-boot-egress-init && " +
		"curl -sS -m 10 --retry 3 --retry-delay 2 --retry-all-errors -o /dev/null -w 'allowed=%{http_code}\\n' https://github.com; " +
		"curl -sS -m 5 -o /dev/null https://example.com && echo example-REACHABLE || echo example-blocked"
	args := append(append([]string{"run", "--rm"}, runArgs...), image, "bash", "-c", probe)
	out, code := dockerOutput(ctx, t, args...)
	if code != 0 || !strings.Contains(out, "allowed=200") || !strings.Contains(out, "example-blocked") {
		t.Errorf("egress restriction: exit %d, want github reachable and example.com blocked:\n%s", code, out)
	}

	// Without NET_ADMIN and policy warn: no restriction, exit 0, warning.
	out, code = dockerOutput(ctx, t, "run", "--rm", image, "sudo", "sh", "/usr/local/bin/u-boot-egress-init")
	if code != 0 || !strings.Contains(out, "egress restriction unavailable") {
		t.Errorf("warn policy without NET_ADMIN: exit %d:\n%s", code, out)
	}

	// Policy fail: exit 11.
	failDir := generateEgressProject(ctx, t, "fail")
	failImage := build(failDir, "uboot-e2e-egress-fail")
	if out, code := dockerOutput(ctx, t, "run", "--rm", failImage, "sudo", "sh", "/usr/local/bin/u-boot-egress-init"); code != 11 {
		t.Errorf("fail policy without NET_ADMIN: exit %d, want 11:\n%s", code, out)
	}
}

// egressSubdomainProbeScript rewrites the `hosts=` line of the generated
// egress script (the allowlist is the only thing that differs) and probes
// names under it. `reach` prints the HTTP status, `000` when the name does
// not resolve or the address is dropped.
const egressSubdomainProbeScript = `reach() { curl -sS -m 8 -o /dev/null -w '%{http_code}' "https://$1" 2>/dev/null || true; }
run_with() {
  sudo sed -i "s|^hosts=.*|hosts=\"$1\"|" /usr/local/bin/u-boot-egress-init
  sudo sh /usr/local/bin/u-boot-egress-init
}
run_with "example.org amazonaws.com"
echo "A example.org=$(reach example.org)"
echo "A www.example.org=$(reach www.example.org)"
echo "A s3.us-east-1.amazonaws.com=$(reach s3.us-east-1.amazonaws.com)"
echo "A example.com=$(reach example.com)"
run_with "www.example.org"
echo "B www.example.org=$(reach www.example.org)"
echo "B example.org=$(reach example.org)"
`

// LH-FA-DEV-008, ADR-0012: the entry example.org makes example.org and
// www.example.org reachable and example.com not; the entry amazonaws.com
// makes s3.us-east-1.amazonaws.com (two levels below it) reachable; the
// entry www.example.org makes www.example.org reachable and example.org not
// by name.
func TestE2E_SandboxDevcontainer_EgressSubdomains(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	dir := generateEgressProject(ctx, t, "warn")
	image := "uboot-e2e-egress-sub:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("docker build failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })
	args := append(append([]string{"run", "--rm"}, generatedRunArgs(t, dir)...), image, "bash", "-c", egressSubdomainProbeScript)
	out, code := dockerOutput(ctx, t, args...)
	if code != 0 {
		t.Fatalf("subdomain probe exit %d:\n%s", code, out)
	}
	got := map[string]string{}
	for _, line := range strings.Split(out, "\n") {
		if tag, rest, ok := strings.Cut(strings.TrimSpace(line), " "); ok && (tag == "A" || tag == "B") {
			if name, status, ok := strings.Cut(rest, "="); ok {
				got[tag+" "+name] = status
			}
		}
	}
	reachable := func(key string) bool { s, ok := got[key]; return ok && s != "000" && s != "" }
	for _, key := range []string{"A example.org", "A www.example.org", "A s3.us-east-1.amazonaws.com", "B www.example.org"} {
		if !reachable(key) {
			t.Errorf("%s: want reachable (entry covers the name and its subdomains), got %q\n%s", key, got[key], out)
		}
	}
	for _, key := range []string{"A example.com", "B example.org"} {
		if s, ok := got[key]; !ok || s != "000" {
			t.Errorf("%s: want blocked (000), got %q (present %v)\n%s", key, s, ok, out)
		}
	}
}

// egressUpstreamDNSProbeScript runs in the container with the generated
// egress script. `dnsq` asks one server directly and prints `answer` when a
// reply with a successful code and at least one answer record arrives within
// three seconds, `none` otherwise.
const egressUpstreamDNSProbeScript = `set -u
cat > /tmp/dnsq.pl <<'PERL'
use strict; use warnings;
use IO::Socket::INET; use IO::Select;
my ($proto, $srv, $name) = @ARGV;
my $id = int(rand 65536);
my $q = pack('nnnnnn', $id, 0x0100, 1, 0, 0, 0) . join('', map { pack('C', length $_) . $_ } split(/\./, $name)) . "\0" . pack('nn', 1, 1);
my $buf = '';
if ($proto eq 'udp') {
  my $s = IO::Socket::INET->new(PeerAddr => $srv, PeerPort => 53, Proto => 'udp') or do { print "none\n"; exit };
  $s->send($q);
  $s->recv($buf, 512) if IO::Select->new($s)->can_read(3);
} else {
  my $s = IO::Socket::INET->new(PeerAddr => $srv, PeerPort => 53, Proto => 'tcp', Timeout => 3) or do { print "none\n"; exit };
  $s->print(pack('n', length $q) . $q);
  if (IO::Select->new($s)->can_read(3)) { my $l = ''; $s->read($l, 2); $s->read($buf, unpack('n', $l)) if length $l == 2; }
}
my ($rid, $fl, $qd, $an) = length $buf >= 12 ? unpack('nnnn', $buf) : (-1, 0xf, 0, 0);
print(($rid == $id && ($fl & 0xf) == 0 && $an > 0) ? "answer\n" : "none\n");
PERL
dnsq() { perl /tmp/dnsq.pl "$@"; }
sudo sh /usr/local/bin/u-boot-egress-init; echo "INIT_EXIT=$?"
UPS=$(awk '/^nameserver/ && $2 !~ /:/ {print $2}' /etc/resolv.conf.u-boot-orig)
echo "UPSTREAMS=$(echo $UPS | wc -w)"
for ns in $UPS; do
  for p in udp tcp; do echo "DIRECT $p example.com=$(dnsq $p $ns example.com)"; done
done
ns1=$(echo $UPS | cut -d' ' -f1)
echo "LIMIT sudo-as-resolver-user=$(sudo setpriv --reuid="$(id -u nobody)" --regid="$(id -g nobody)" --clear-groups perl /tmp/dnsq.pl udp "$ns1" example.com)"
echo "LOCAL api.github.com=$(getent hosts api.github.com | head -1 | cut -d' ' -f1 | grep -c .)"
echo "LOCAL example.com=$(getent hosts example.com | head -1 | cut -d' ' -f1 | grep -c .)"
echo "HTTPS api.github.com=$(curl -sS -m 10 --retry 3 --retry-delay 2 --retry-all-errors -o /dev/null -w '%{http_code}' https://api.github.com 2>/dev/null)"
echo "UID dnsmasq=$(awk '/^Uid:/ {print $2}' /proc/$(cat /run/u-boot-egress-dnsmasq.pid)/status) rule=$(id -u nobody)"
echo "RULES skuid=$(sudo nft list chain ip u_boot_egress output | grep -c "meta skuid $(id -u nobody) ip daddr .* dport 53 accept") upstreams=$(echo $UPS | wc -w)"
echo "STATUS $(grep '^result=' /var/lib/u-boot/status/egress)"
`

// egressAbortProbeScript starts mutated copies of the generated script:
// the uid of the upstream rules differs from the uid of the resolver, and
// the resolver user does not exist. Both end with exit 11 and status
// `failed` although the policy is `warn`.
const egressAbortProbeScript = `set -u
src=/usr/local/bin/u-boot-egress-init
sed 's/^dns_uid=.*/dns_uid=1234/' $src > /tmp/uid-differs.sh
sed 's/^dns_user=.*/dns_user=u_boot_no_such_user/' $src > /tmp/user-missing.sh
for m in uid-differs user-missing; do
  sudo sh /tmp/$m.sh 2>/tmp/$m.err; echo "ABORT $m exit=$? $(grep '^result=' /var/lib/u-boot/status/egress)"
  sed "s/^/ERR $m: /" /tmp/$m.err
done
`

// LH-FA-DEV-008, ADR-0026: a process without sudo gets no answer for a name
// outside the allowlist from the upstream resolvers (UDP and TCP port 53);
// allowed names resolve through the local resolver; the uid of the resolver
// matches the uid in the rules; the limit holds: as the resolver user (via
// sudo) the upstream answers. A uid mismatch or a missing resolver user
// stops the start script with exit 11 and status `failed`, also under
// onUnavailable warn.
func TestE2E_SandboxDevcontainer_EgressUpstreamDNS(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	dir := generateEgressProject(ctx, t, "warn")
	image := "uboot-e2e-egress-dns:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("docker build failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })
	runArgs := generatedRunArgs(t, dir)

	t.Run("direct path closed", func(t *testing.T) {
		args := append(append([]string{"run", "--rm"}, runArgs...), image, "bash", "-c", egressUpstreamDNSProbeScript)
		out, code := dockerOutput(ctx, t, args...)
		if code != 0 {
			t.Fatalf("probe exit %d:\n%s", code, out)
		}
		lines := map[string]bool{}
		for _, l := range strings.Split(out, "\n") {
			lines[strings.TrimSpace(l)] = true
		}
		if !lines["INIT_EXIT=0"] {
			t.Fatalf("start script did not exit 0:\n%s", out)
		}
		if lines["UPSTREAMS=0"] {
			t.Skipf("no IPv4 upstream resolver in /etc/resolv.conf of the container:\n%s", out)
		}
		for _, want := range []string{
			"DIRECT udp example.com=none", "DIRECT tcp example.com=none",
			"LIMIT sudo-as-resolver-user=answer",
			"LOCAL api.github.com=1", "LOCAL example.com=0", "HTTPS api.github.com=200",
			"STATUS result=ok",
		} {
			if !lines[want] {
				t.Errorf("output lacks %q:\n%s", want, out)
			}
		}
		for _, l := range strings.Split(out, "\n") {
			if strings.HasPrefix(l, "DIRECT ") && !strings.HasSuffix(l, "=none") {
				t.Errorf("direct upstream query got an answer: %q\n%s", l, out)
			}
			if rest, ok := strings.CutPrefix(l, "UID dnsmasq="); ok {
				got, want, _ := strings.Cut(rest, " rule=")
				if got == "" || got != want {
					t.Errorf("resolver uid %q differs from the uid in the rules %q", got, want)
				}
			}
			if rest, ok := strings.CutPrefix(l, "RULES skuid="); ok {
				n, ups, _ := strings.Cut(rest, " upstreams=")
				if upN, err := strconv.Atoi(ups); err != nil || n != strconv.Itoa(2*upN) {
					t.Errorf("rules with meta skuid = %s, want 2 per upstream (%s upstreams)", n, ups)
				}
			}
		}
	})

	t.Run("uid mismatch and missing user abort", func(t *testing.T) {
		args := append(append([]string{"run", "--rm"}, runArgs...), image, "bash", "-c", egressAbortProbeScript)
		out, code := dockerOutput(ctx, t, args...)
		if code != 0 {
			t.Fatalf("probe exit %d:\n%s", code, out)
		}
		for _, want := range []string{
			"ABORT uid-differs exit=11 result=failed",
			"ERR uid-differs: u-boot sandbox: egress restriction failed: the local resolver is not running as 'nobody'",
			"ABORT user-missing exit=11 result=failed",
			"ERR user-missing: u-boot sandbox: egress restriction not set up: the user 'u_boot_no_such_user' for the local resolver does not exist",
		} {
			if !strings.Contains(out, want) {
				t.Errorf("output lacks %q:\n%s", want, out)
			}
		}
	})
}

func generatedMounts(t *testing.T, dir string) []string {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(dir, ".devcontainer", "devcontainer.json"))
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	for _, l := range strings.Split(string(body), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(l), "//") {
			lines = append(lines, l)
		}
	}
	var dc struct {
		WorkspaceMount string   `json:"workspaceMount"`
		Mounts         []string `json:"mounts"`
	}
	if err := json.Unmarshal([]byte(strings.Join(lines, "\n")), &dc); err != nil {
		t.Fatalf("devcontainer.json: %v", err)
	}
	if dc.WorkspaceMount != "" {
		t.Fatalf("workspaceMount = %q, want empty (no host bind mount)", dc.WorkspaceMount)
	}
	return dc.Mounts
}

// Two instances of the same project (different ${devcontainerId},
// as the Dev Containers tooling resolves it per folder) must not
// share the workspace or the Podman storage volumes.
func TestE2E_SandboxDevcontainer_TwoInstancesIsolated(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	dir := generateSandboxProject(ctx, t, "podman", "warn")
	image := "uboot-e2e-sandbox-multi:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("docker build failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })

	mounts := generatedMounts(t, dir)
	if len(mounts) != 2 {
		t.Fatalf("mounts = %v, want workspace + containers volume", mounts)
	}
	runArgs := generatedRunArgs(t, dir)
	start := func(id string) string {
		name := "uboot-e2e-inst-" + id
		args := []string{"run", "-d", "--name", name}
		for _, m := range mounts {
			if strings.Contains(m, "type=bind") {
				t.Fatalf("bind mount in sandbox mounts: %s", m)
			}
			args = append(args, "--mount", strings.ReplaceAll(m, "${devcontainerId}", id))
		}
		args = append(append(args, runArgs...), image, "sleep", "600")
		if out, code := dockerOutput(ctx, t, args...); code != 0 {
			t.Fatalf("start instance %s: %s", id, out)
		}
		t.Cleanup(func() {
			_, _ = dockerOutput(context.Background(), t, "rm", "-f", name)
			for _, m := range mounts {
				src := strings.TrimPrefix(strings.Split(strings.ReplaceAll(m, "${devcontainerId}", id), ",")[0], "source=")
				_, _ = dockerOutput(context.Background(), t, "volume", "rm", "-f", src)
			}
		})
		return name
	}
	a, b := start("idaaaa"), start("idbbbb")
	ws := "/workspaces/sbxe2e"
	if out, code := dockerOutput(ctx, t, "exec", a, "sh", "-c", "echo instance-a > "+ws+"/marker"); code != 0 {
		t.Fatalf("write in instance a: %s", out)
	}
	if out, code := dockerOutput(ctx, t, "exec", b, "sh", "-c", "test ! -e "+ws+"/marker && echo instance-b > "+ws+"/marker"); code != 0 {
		t.Errorf("instance b sees instance a's workspace: %s", out)
	}
	if out, _ := dockerOutput(ctx, t, "exec", a, "cat", ws+"/marker"); strings.TrimSpace(out) != "instance-a" {
		t.Errorf("instance a workspace content = %q", out)
	}
}

func generatedRunArgs(t *testing.T, dir string) []string {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(dir, ".devcontainer", "devcontainer.json"))
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	for _, l := range strings.Split(string(body), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(l), "//") {
			lines = append(lines, l)
		}
	}
	var dc struct {
		RunArgs []string `json:"runArgs"`
	}
	if err := json.Unmarshal([]byte(strings.Join(lines, "\n")), &dc); err != nil {
		t.Fatalf("devcontainer.json: %v", err)
	}
	if len(dc.RunArgs) == 0 {
		t.Fatalf("generated devcontainer.json has no runArgs")
	}
	return dc.RunArgs
}

// skipIfNestedUsernsRestricted skips the test on hosts with
// kernel.apparmor_restrict_unprivileged_userns=1 (Ubuntu 24.04+, Colima's
// Ubuntu VM): they block nested user namespaces even with
// apparmor=unconfined, so the startup script must degrade with exit 11
// (LH-FA-DEV-007) and nested Podman cannot run. The degradation is
// verified before the skip.
func skipIfNestedUsernsRestricted(ctx context.Context, t *testing.T, image string, runArgs []string) {
	t.Helper()
	sysctl, err := os.ReadFile("/proc/sys/kernel/apparmor_restrict_unprivileged_userns")
	if err != nil || strings.TrimSpace(string(sysctl)) != "1" {
		return
	}
	out, code := dockerOutput(ctx, t, append(append([]string{"run", "--rm"}, runArgs...), image, "sh", "/usr/local/bin/u-boot-sandbox-init")...)
	if code != 11 || !strings.Contains(out, "nested user namespaces are blocked") {
		t.Errorf("restricted host: init exit %d, want 11 with the blocked-namespaces message:\n%s", code, out)
	}
	t.Skip("host restricts unprivileged user namespaces (kernel.apparmor_restrict_unprivileged_userns=1); nested Podman cannot run here, degradation verified")
}

// nestedPodmanBuildScript runs as `vscode` in the sandbox container and
// builds a Dockerfile with FROM, RUN and COPY: cold, from the layer cache,
// with --no-cache, and through the `docker` alias.
const nestedPodmanBuildScript = `set -e
sh /usr/local/bin/u-boot-sandbox-init
mkdir -p "$HOME/ctx" && cd "$HOME/ctx"
echo copied > datei.txt
printf 'FROM docker.io/library/alpine:3\nRUN echo ran > /run.txt\nCOPY datei.txt /datei\n' > Dockerfile
podman build -t probe . && echo BUILD_COLD_OK
test "$(podman run --rm probe cat /datei /run.txt | tr '\n' ' ')" = "copied ran " && echo RUN_OF_BUILT_IMAGE_OK
podman build -t probe . > "$HOME/cached.log" 2>&1
grep -q 'Using cache' "$HOME/cached.log" && echo BUILD_CACHE_OK
podman build --no-cache -t probe . > "$HOME/nocache.log" 2>&1
! grep -q 'Using cache' "$HOME/nocache.log" && echo BUILD_NOCACHE_OK
docker build -t probe-alias . && echo DOCKER_ALIAS_OK`

// LH-FA-DEV-007 (nestedRuntime podman): `podman build` and, through the
// `docker` symlink, `docker build` build a Dockerfile inside the
// container without a host socket. Runs only on hosts that allow nested
// user namespaces; elsewhere it verifies the degradation and skips.
func TestE2E_SandboxDevcontainer_NestedPodman_Build(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
	defer cancel()

	dir := generateSandboxProject(ctx, t, "podman", "warn")
	image := "uboot-e2e-sandbox-build:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("docker build of the generated Dockerfile failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })

	runArgs := generatedRunArgs(t, dir)
	skipIfNestedUsernsRestricted(ctx, t, image, runArgs)

	args := append(append([]string{"run", "--rm"}, runArgs...), image, "bash", "-c", nestedPodmanBuildScript)
	out, code := dockerOutput(ctx, t, args...)
	if code != 0 {
		t.Fatalf("nested podman build script exit %d:\n%s", code, out)
	}
	for _, want := range []string{"BUILD_COLD_OK", "RUN_OF_BUILT_IMAGE_OK", "BUILD_CACHE_OK", "BUILD_NOCACHE_OK", "DOCKER_ALIAS_OK"} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %s:\n%s", want, out)
		}
	}
}

// nestedPodmanEgressBuildScript runs as `vscode` behind the egress
// restriction: the image pull needs the registry hosts of the default
// allowlist, a `RUN` reaches an allowed host and not an unlisted one.
const nestedPodmanEgressBuildScript = `set -e
sudo sh /usr/local/bin/u-boot-egress-init
sh /usr/local/bin/u-boot-sandbox-init
mkdir -p "$HOME/ctx" && cd "$HOME/ctx"
printf 'FROM docker.io/library/debian:stable-slim\nRUN getent hosts deb.debian.org > /dev/null && echo ALLOWED_RESOLVES\nRUN ! getent hosts example.com && echo UNLISTED_BLOCKED\n' > Dockerfile
podman build -t probe-egress . && echo EGRESS_BUILD_OK`

// LH-FA-DEV-007 with LH-FA-DEV-008: `podman build` behind the egress
// restriction pulls from Docker Hub (registry and blob hosts of the default
// allowlist) and resolves only listed names in `RUN`.
func TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
	defer cancel()

	dir := generateSandboxProjectWith(ctx, t, map[string]string{
		"devcontainer.sandbox.nestedRuntime":  "podman",
		"devcontainer.sandbox.onUnavailable":  "warn",
		"devcontainer.sandbox.egress.enabled": "true",
	})
	image := "uboot-e2e-sandbox-build-egress:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("docker build of the generated Dockerfile failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })

	runArgs := generatedRunArgs(t, dir)
	skipIfNestedUsernsRestricted(ctx, t, image, runArgs)

	args := append(append([]string{"run", "--rm"}, runArgs...), image, "bash", "-c", nestedPodmanEgressBuildScript)
	out, code := dockerOutput(ctx, t, args...)
	if code != 0 {
		t.Fatalf("nested podman build behind egress exit %d:\n%s", code, out)
	}
	for _, want := range []string{"ALLOWED_RESOLVES", "UNLISTED_BLOCKED", "EGRESS_BUILD_OK"} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %s:\n%s", want, out)
		}
	}
}

func TestE2E_SandboxDevcontainer_NestedPodman(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	t.Setenv("GIT_TOKEN", sandboxCanaryToken)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	dir := generateSandboxProject(ctx, t, "podman", "warn")

	// LH-FA-DEV-009 / AK 4: no credential value in any generated file.
	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if b, rerr := os.ReadFile(p); rerr == nil && strings.Contains(string(b), sandboxCanaryToken) {
				t.Errorf("credential value found in generated file %s", p)
			}
		}
		return nil
	})

	image := "uboot-e2e-sandbox:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("docker build of the generated Dockerfile failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })

	runArgs := generatedRunArgs(t, dir)
	skipIfNestedUsernsRestricted(ctx, t, image, runArgs)

	// ADR 0023: Podman of the pinned major line from /usr/local (no Debian
	// package), the registries file with Docker Hub only, the `docker` alias.
	args := append(append([]string{"run", "--rm"}, runArgs...), image, "bash", "-c", nestedPodmanVersionScript)
	out, code := dockerOutput(ctx, t, args...)
	if code != 0 {
		t.Fatalf("podman source/version checks failed (exit %d):\n%s", code, out)
	}
	for _, want := range []string{"PODMAN_MAJOR_6", "PODMAN_FROM_USR_LOCAL", "NO_DEBIAN_PODMAN", "NO_ARCHIVE_README", "REGISTRIES_DOCKER_IO", "DOCKER_IS_PODMAN"} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %s:\n%s", want, out)
		}
	}

	// AK 2: docker run inside the container without a socket.
	nested := "sh /usr/local/bin/u-boot-sandbox-init && docker run --rm --network=host docker.io/library/alpine:3 echo NESTED_OK"
	args = append(append([]string{"run", "--rm"}, runArgs...), image, "bash", "-c", nested)
	if out, code := dockerOutput(ctx, t, args...); code != 0 || !strings.Contains(out, "NESTED_OK") {
		t.Fatalf("nested podman run failed (exit %d):\n%s", code, out)
	}

	// Degradation: without the relaxations the start script aborts
	// with exit 11 (blocked user namespaces).
	if out, code := dockerOutput(ctx, t, "run", "--rm", image, "sh", "/usr/local/bin/u-boot-sandbox-init"); code != 11 {
		t.Errorf("without relaxations: exit %d, want 11:\n%s", code, out)
	}

	// Degradation: container started without the /dev/fuse device and
	// policy warn → vfs fallback.
	// The fallback must work: the effective driver is vfs, `podman run`,
	// `podman build` and `docker build` run. The storage volume is mounted
	// like the generated devcontainer.json mounts it (`image_copy_tmp_dir`
	// lies there).
	volume := "sbxe2e-vfs-" + strconv.FormatInt(time.Now().UnixNano(), 36)
	if out, code := dockerOutput(ctx, t, "volume", "create", volume); code != 0 {
		t.Fatalf("docker volume create: %s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "volume", "rm", "-f", volume) })
	noFuse := append(withoutArg(runArgs, "--device=/dev/fuse"), "-v", volume+":"+containersStorageTarget)
	args = append(append([]string{"run", "--rm"}, noFuse...), image, "bash", "-c", nestedPodmanVFSScript)
	out, code = dockerOutput(ctx, t, args...)
	if code != 0 {
		t.Errorf("warn policy without /dev/fuse: exit %d:\n%s", code, out)
	}
	for _, want := range []string{"VFS_DRIVER", "VFS_RUN_OK", "VFS_BUILD_OK", "VFS_DOCKER_BUILD_OK"} {
		if !strings.Contains(out, want) {
			t.Errorf("vfs fallback output lacks %s:\n%s", want, out)
		}
	}
}

// containersStorageTarget is the target of the storage volume in the
// generated devcontainer.json (pinned by TestGenerateDevcontainer_SandboxPodman).
const containersStorageTarget = "/home/vscode/.local/share/containers"

// nestedPodmanVersionScript checks, as `vscode` after the start script, where
// Podman comes from (ADR 0023): major line 6, the static archive under
// /usr/local, no Debian package, Docker Hub as the only short-name registry,
// `docker` as the Podman alias.
const nestedPodmanVersionScript = `set -e
sh /usr/local/bin/u-boot-sandbox-init
podman --version | grep -q '^podman version 6\.' && echo PODMAN_MAJOR_6
test "$(command -v podman)" = /usr/local/bin/podman && echo PODMAN_FROM_USR_LOCAL
! dpkg -l podman fuse-overlayfs passt 2> /dev/null | grep -q '^ii' && echo NO_DEBIAN_PODMAN
test ! -e /README.md && echo NO_ARCHIVE_README
grep -Fx "unqualified-search-registries = ['docker.io']" /etc/containers/registries.conf > /dev/null && echo REGISTRIES_DOCKER_IO
test "$(readlink -f "$(command -v docker)")" = "$(readlink -f /usr/local/bin/podman)" && echo DOCKER_IS_PODMAN`

// nestedPodmanVFSScript runs without /dev/fuse and with the storage volume:
// the start script falls back to vfs; the effective storage driver, a run and
// builds through `podman` and the `docker` alias must work.
const nestedPodmanVFSScript = `set -e
sh /usr/local/bin/u-boot-sandbox-init
test "$(podman info --format '{{.Store.GraphDriverName}}')" = vfs && echo VFS_DRIVER
podman run --rm --network=host docker.io/library/alpine:3 echo VFS_RUN_OK
mkdir -p "$HOME/ctx" && cd "$HOME/ctx"
printf 'FROM docker.io/library/alpine:3\nRUN echo ran > /run.txt\n' > Dockerfile
podman build -t probe . && echo VFS_BUILD_OK
docker build -t probe-alias . && echo VFS_DOCKER_BUILD_OK`

// ADR 0023 negative probe: the build stops before the installation when the
// pinned checksum or the pinned signing key does not match. The pins are
// `ARG`s of the generated Dockerfile; `--build-arg` replaces one value per
// probe, the control build keeps the generated pins and must succeed (the
// other nested-Podman tests build with them too). The checksum probe replaces
// both architecture values, so it holds on `amd64` and `arm64` hosts.
func TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild(t *testing.T) {
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	dir := generateSandboxProject(ctx, t, "podman", "warn")
	contextDir := filepath.Join(dir, ".devcontainer")
	zeros := strings.Repeat("0", 64)
	for name, tc := range map[string]struct {
		args []string
		want string
	}{
		"wrong checksum":                 {[]string{"PODMAN_SHA256_AMD64=" + zeros, "PODMAN_SHA256_ARM64=" + zeros}, "FAILED"},
		"signing subkey as the pin":      {[]string{"PODMAN_SIGNING_KEY=6E8EF2F5B1F26FF78431E3D1364FA5A62B410BA4"}, ".tar.gz: OK"},
		"unknown fingerprint as the pin": {[]string{"PODMAN_SIGNING_KEY=" + strings.Repeat("0", 40)}, ".tar.gz: OK"},
	} {
		t.Run(name, func(t *testing.T) {
			args := []string{"build", "-q", "-t", "uboot-e2e-sandbox-pin-negative:probe"}
			for _, a := range tc.args {
				args = append(args, "--build-arg", a)
			}
			out, code := dockerOutput(ctx, t, append(args, contextDir)...)
			if code == 0 {
				t.Cleanup(func() {
					_, _ = dockerOutput(context.Background(), t, "rmi", "-f", "uboot-e2e-sandbox-pin-negative:probe")
				})
				t.Fatalf("build with %v succeeded, want a failure before the installation:\n%s", tc.args, out)
			}
			if tc.want != "" && !strings.Contains(out, tc.want) {
				t.Errorf("build output lacks %q:\n%s", tc.want, out)
			}
		})
	}
}

// LH-FA-DEV-007: a `/dev/fuse` that is missing on the container host is no
// degradation step; the engine refuses to start the container, so the start
// script never runs. The host device is replaced by a path that does not
// exist; the same runArgs with the real device start the container.
func TestE2E_SandboxDevcontainer_MissingHostFuseDeviceRefusesStart(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	dir := generateSandboxProject(ctx, t, "podman", "warn")
	image := "uboot-e2e-sandbox-nofusehost:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("docker build failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })

	runArgs := generatedRunArgs(t, dir)
	if out, code := dockerOutput(ctx, t, append(append([]string{"run", "--rm"}, runArgs...), image, "true")...); code != 0 {
		t.Fatalf("control: container with the generated runArgs does not start (exit %d):\n%s", code, out)
	}

	var missing []string
	replaced := false
	for _, a := range runArgs {
		if a == "--device=/dev/fuse" {
			a, replaced = "--device=/dev/fuse-nonexistent", true
		}
		missing = append(missing, a)
	}
	if !replaced {
		t.Fatalf("generated runArgs carry no --device=/dev/fuse: %v", runArgs)
	}
	out, code := dockerOutput(ctx, t, append(append([]string{"run", "--rm"}, missing...), image, "true")...)
	if code == 0 || !strings.Contains(out, "error gathering device information") {
		t.Errorf("missing host device: exit %d, want an engine refusal (error gathering device information):\n%s", code, out)
	}
}

// LH-FA-DEV-007: a container started without the `/dev/fuse` device (changed
// runArgs) ends the start script with exit 11 under the `fail` policy.
func TestE2E_SandboxDevcontainer_FailPolicyWithoutFuse(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	dir := generateSandboxProject(ctx, t, "podman", "fail")
	image := "uboot-e2e-sandbox-fail:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("docker build failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })

	var noFuse []string
	for _, a := range generatedRunArgs(t, dir) {
		if a != "--device=/dev/fuse" {
			noFuse = append(noFuse, a)
		}
	}
	args := append(append([]string{"run", "--rm"}, noFuse...), image, "sh", "/usr/local/bin/u-boot-sandbox-init")
	if out, code := dockerOutput(ctx, t, args...); code != 11 {
		t.Errorf("fail policy without /dev/fuse: exit %d, want 11:\n%s", code, out)
	}
}

// overlayDriver reports whether `docker info --format '{{.Driver}}'` names
// an overlay storage backend: `overlay2` (classic graph driver) or
// `overlayfs` (containerd image store, default of fresh Docker 29 installs).
func overlayDriver(info string) bool {
	return strings.HasPrefix(info, "overlay2") || strings.HasPrefix(info, "overlayfs")
}

// dockerExec runs `bash -c script` in the container as user `vscode`
// (no sudo, group `docker` from /etc/group).
func dockerExec(ctx context.Context, t *testing.T, container, script string) (string, int) {
	t.Helper()
	return dockerOutput(ctx, t, "exec", "-u", "vscode", container, "bash", "-c", script)
}

const nestedDockerBuildScript = `set -e
mkdir -p /tmp/ctx && cd /tmp/ctx
printf 'FROM alpine:3 AS base\nRUN echo base > /base.txt\nFROM base AS final\nRUN --network=none echo nonet > /n.txt\n' > Dockerfile
printf 'services:\n  hello:\n    image: alpine:3\n    command: ["echo","COMPOSE_OK"]\n' > compose.yaml
docker buildx build --target base -o type=cacheonly . && echo CACHEONLY_OK
docker build -t same:latest --target final . && echo NETWORK_NONE_OK
docker compose up 2>&1 | grep -q COMPOSE_OK && echo COMPOSE_OK
docker rm -f same-name >/dev/null 2>&1 || true
docker run -d --name same-name same:latest sleep 600 >/dev/null && echo SAME_NAMES_OK
docker run --rm --privileged docker:dind docker --version && echo THIRD_LEVEL_OK`

// LH-FA-DEV-007 (nestedRuntime docker), acceptance criteria 2 and 3:
// two instances of the same project each run an own daemon with
// BuildKit, Buildx and Compose, as `vscode` without sudo, and
// identically named images and containers do not collide.
func TestE2E_SandboxDevcontainer_NestedDocker(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
	defer cancel()

	dir := generateSandboxProject(ctx, t, "docker", "warn")
	image := "uboot-e2e-sandbox-docker:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("docker build of the generated Dockerfile failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })

	runArgs := generatedRunArgs(t, dir)
	if len(runArgs) != 1 || runArgs[0] != "--privileged" {
		t.Fatalf("runArgs = %v, want [--privileged]", runArgs)
	}
	mounts := generatedMounts(t, dir)
	if len(mounts) != 2 || !strings.Contains(mounts[1], "-docker-${devcontainerId}") {
		t.Fatalf("mounts = %v, want workspace + docker volume", mounts)
	}
	start := func(id string) string {
		name := "uboot-e2e-dockerinst-" + id
		args := []string{"run", "-d", "--name", name}
		for _, m := range mounts {
			args = append(args, "--mount", strings.ReplaceAll(m, "${devcontainerId}", id))
		}
		args = append(append(args, runArgs...), image, "sleep", "1800")
		if out, code := dockerOutput(ctx, t, args...); code != 0 {
			t.Fatalf("start instance %s: %s", id, out)
		}
		t.Cleanup(func() {
			_, _ = dockerOutput(context.Background(), t, "rm", "-f", name)
			for _, m := range mounts {
				src := strings.TrimPrefix(strings.Split(strings.ReplaceAll(m, "${devcontainerId}", id), ",")[0], "source=")
				_, _ = dockerOutput(context.Background(), t, "volume", "rm", "-f", src)
			}
		})
		if out, code := dockerOutput(ctx, t, "exec", "-u", "vscode", name, "sudo", "sh", "/usr/local/bin/u-boot-docker-init"); code != 0 {
			t.Fatalf("daemon start script (instance %s) exit %d:\n%s", id, code, out)
		}
		return name
	}
	a, b := start("idaaaa"), start("idbbbb")
	for _, c := range []string{a, b} {
		out, code := dockerExec(ctx, t, c, "docker info --format '{{.Driver}} {{.ID}}'")
		if code != 0 || !overlayDriver(out) {
			t.Fatalf("%s: docker info exit %d: %s", c, code, out)
		}
		if out, code := dockerExec(ctx, t, c, nestedDockerBuildScript); code != 0 {
			t.Fatalf("%s: nested docker script exit %d:\n%s", c, code, out)
		} else {
			for _, want := range []string{"CACHEONLY_OK", "NETWORK_NONE_OK", "COMPOSE_OK", "SAME_NAMES_OK", "THIRD_LEVEL_OK"} {
				if !strings.Contains(out, want) {
					t.Errorf("%s: output lacks %s:\n%s", c, want, out)
				}
			}
		}
	}
	// Origin and age: Docker's own packages (docker-ce), not the Debian
	// archive (docker.io, 26.1); a current major version.
	for _, c := range []string{a, b} {
		if out, code := dockerExec(ctx, t, c, "dpkg-query -W -f='${Package}\\n' docker-ce containerd.io docker-buildx-plugin docker-compose-plugin"); code != 0 {
			t.Errorf("%s: packages from the Docker APT repository missing (exit %d): %s", c, code, out)
		}
		out, code := dockerExec(ctx, t, c, "docker version --format '{{.Server.Version}}'")
		major, convErr := strconv.Atoi(strings.SplitN(strings.TrimSpace(out), ".", 2)[0])
		if code != 0 || convErr != nil || major < 28 {
			t.Errorf("%s: Docker server version %q (exit %d) is not a major version >= 28", c, strings.TrimSpace(out), code)
		}
	}
	idA, _ := dockerExec(ctx, t, a, "docker info --format '{{.ID}}'")
	idB, _ := dockerExec(ctx, t, b, "docker info --format '{{.ID}}'")
	if strings.TrimSpace(idA) == "" || idA == idB {
		t.Errorf("instances share one daemon: %q vs %q", idA, idB)
	}
}

// LH-FA-DEV-007 degradation, acceptance criterion 5: a start without
// `--privileged` ends without a daemon at `warn` (exit 0, warning) and
// with exit code 11 at `fail`; podman never replaces the daemon.
func TestE2E_SandboxDevcontainer_NestedDocker_Degradation(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	for policy, wantCode := range map[string]int{"warn": 0, "fail": 11} {
		dir := generateSandboxProject(ctx, t, "docker", policy)
		image := "uboot-e2e-sandbox-docker-" + policy + ":" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
		if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
			t.Fatalf("docker build failed:\n%s", out)
		}
		t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })
		out, code := dockerOutput(ctx, t, "run", "--rm", image, "sudo", "sh", "/usr/local/bin/u-boot-docker-init")
		if code != wantCode || !strings.Contains(out, "no daemon available") {
			t.Errorf("policy %s without --privileged: exit %d, want %d with a warning:\n%s", policy, code, wantCode, out)
		}
		if out, code := dockerOutput(ctx, t, "run", "--rm", image, "sh", "-c", "command -v podman"); code == 0 {
			t.Errorf("policy %s: podman present in the image: %s", policy, out)
		}
	}
}

// LH-FA-DEV-007: a container restart kills the daemon without cleanup;
// the start script brings a new daemon up on the same storage volume.
func TestE2E_SandboxDevcontainer_NestedDocker_Restart(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	dir := generateSandboxProject(ctx, t, "docker", "fail")
	image := "uboot-e2e-sandbox-docker-restart:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("docker build failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })

	name, volume := "uboot-e2e-dockerrestart", "uboot-e2e-dockerrestart-vol"
	args := []string{"run", "-d", "--name", name, "--mount", "source=" + volume + ",target=/var/lib/docker,type=volume"}
	args = append(append(args, generatedRunArgs(t, dir)...), image, "sleep", "1800")
	if out, code := dockerOutput(ctx, t, args...); code != 0 {
		t.Fatalf("start: %s", out)
	}
	t.Cleanup(func() {
		_, _ = dockerOutput(context.Background(), t, "rm", "-f", name)
		_, _ = dockerOutput(context.Background(), t, "volume", "rm", "-f", volume)
	})
	initDaemon := func(step string) {
		if out, code := dockerOutput(ctx, t, "exec", "-u", "vscode", name, "sudo", "sh", "/usr/local/bin/u-boot-docker-init"); code != 0 {
			t.Fatalf("%s: daemon start script exit %d:\n%s", step, code, out)
		}
		if out, code := dockerExec(ctx, t, name, "docker info --format '{{.Driver}}'"); code != 0 || !overlayDriver(out) {
			t.Fatalf("%s: docker info exit %d: %s", step, code, out)
		}
	}
	initDaemon("first start")
	if out, code := dockerExec(ctx, t, name, "docker pull -q alpine:3 >/dev/null && docker image inspect alpine:3 >/dev/null && echo PULLED"); code != 0 || !strings.Contains(out, "PULLED") {
		t.Fatalf("pull into the daemon storage: exit %d: %s", code, out)
	}
	if out, code := dockerOutput(ctx, t, "restart", "-t", "0", name); code != 0 {
		t.Fatalf("restart: %s", out)
	}
	// The killed daemon leaves its PID file; a live process (PID 1) under
	// that number makes dockerd refuse to start unless the script cleans up.
	if out, code := dockerOutput(ctx, t, "exec", "-u", "root", name, "sh", "-c", "echo 1 > /var/run/docker.pid"); code != 0 {
		t.Fatalf("seed stale PID file: %s", out)
	}
	initDaemon("after restart")
	if out, code := dockerExec(ctx, t, name, "docker image inspect alpine:3 >/dev/null && echo STORAGE_KEPT"); code != 0 || !strings.Contains(out, "STORAGE_KEPT") {
		t.Errorf("daemon storage volume not reused after restart: exit %d: %s", code, out)
	}
}

// LH-FA-DEV-008 with nestedRuntime docker: the Docker APT repository stays
// configured in the image, so `apt-get update` must work behind the egress
// restriction (download.docker.com is in the default allowlist) while an
// unlisted host stays blocked.
func TestE2E_SandboxDevcontainer_NestedDocker_EgressAptUpdate(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	dir := generateSandboxProjectWith(ctx, t, map[string]string{
		"devcontainer.sandbox.nestedRuntime":  "docker",
		"devcontainer.sandbox.onUnavailable":  "fail",
		"devcontainer.sandbox.egress.enabled": "true",
	})
	image := "uboot-e2e-sandbox-docker-egress:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("docker build failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })
	args := append([]string{"run", "--rm"}, generatedRunArgs(t, dir)...)
	args = append(args, image, "bash", "-c",
		"sudo sh /usr/local/bin/u-boot-egress-init && sudo apt-get update --error-on=any && echo APT_UPDATE_OK && "+
			"! curl -sS --max-time 8 -o /dev/null https://example.com && echo BLOCKED_OK")
	out, code := dockerOutput(ctx, t, args...)
	if code != 0 || !strings.Contains(out, "APT_UPDATE_OK") || !strings.Contains(out, "BLOCKED_OK") {
		t.Errorf("egress + docker: exit %d, want apt-get update to succeed and example.com to stay blocked:\n%s", code, out)
	}
}

// ollamaImageScript runs as `vscode` in the image without volume and without
// egress: the pinned version, no GPU libraries, no `zstd`, a model store the
// container user owns, and no environment variable that sets a listen
// address.
const ollamaImageScript = `ollama --version 2>&1 | grep -F '0.40.6' > /dev/null && echo OLLAMA_VERSION_PINNED
test -d /usr/local/lib/ollama && echo OLLAMA_LIB_PRESENT
test -z "$(ls /usr/local/lib/ollama | grep -E '^(cuda_|vulkan)')" && echo NO_GPU_LIBRARIES
! command -v zstd > /dev/null && echo NO_ZSTD
test "$(stat -c %U:%G /home/vscode/.ollama)" = vscode:vscode && echo STORE_OWNED_BY_VSCODE
test "$(env | grep -c '^OLLAMA' || true)" = 0 && echo NO_OLLAMA_ENV`

// ollamaServerScript runs as `vscode` in the container behind the egress
// restriction: it starts the server like the user does, then reports the
// listening sockets, the container-address probe, the compute library, the
// pull of a small model through the two allowlisted hosts, a run on the CPU, the
// key file in the model store, and an unlisted host.
const ollamaServerScript = `sudo sh /usr/local/bin/u-boot-egress-init
nohup ollama serve > /tmp/ollama.log 2>&1 &
for i in $(seq 1 30); do curl -fsS -m 2 http://127.0.0.1:11434/ > /dev/null 2>&1 && break; sleep 1; done
curl -fsS -m 2 http://127.0.0.1:11434/ | grep -F 'Ollama is running' > /dev/null && echo LOOPBACK_ANSWERS
echo "LISTEN:$(ss -ltn | awk '$4 ~ /:11434$/ {print $4}' | tr '\n' ' ')"
ip4="$(hostname -I | awk '{print $1}')"
if curl -sS -m 3 -o /dev/null "http://$ip4:11434/" 2> /dev/null; then echo CONTAINER_ADDRESS_REACHABLE; else echo CONTAINER_ADDRESS_REFUSED; fi
grep -F 'library=cpu' /tmp/ollama.log > /dev/null && echo LIBRARY_CPU
if ollama pull smollm2:135m > /tmp/pull.log 2>&1; then echo PULL_OK; else echo PULL_FAILED; tail -5 /tmp/pull.log; fi
ollama list | grep -E '^smollm2:135m' > /dev/null && echo MODEL_LISTED
test -f /home/vscode/.ollama/id_ed25519 && echo KEY_IN_MODEL_STORE
run_out="$(ollama run smollm2:135m 'Reply with the single word OK.' 2>&1)" && test -n "$run_out" && echo RUN_OK
if curl -sS -m 6 -o /dev/null https://example.com 2> /dev/null; then echo FOREIGN_REACHABLE; else echo FOREIGN_BLOCKED; fi`

// ollamaRestartScript starts the server in a fresh container that mounts the
// volumes of the removed one: the model survives the removal.
const ollamaRestartScript = `nohup ollama serve > /tmp/ollama.log 2>&1 &
for i in $(seq 1 30); do curl -fsS -m 2 http://127.0.0.1:11434/ > /dev/null 2>&1 && break; sleep 1; done
ollama list | grep -E '^smollm2:135m' > /dev/null && echo MODEL_SURVIVES_CONTAINER_REMOVAL`

// LH-FA-DEV-003 (feature `ollama`), LH-FA-DEV-008: the generated Dockerfile
// builds; the image carries the pinned Ollama without GPU libraries and
// `zstd`; behind the egress restriction the server the user starts listens on
// loopback only (the container address refuses), runs on the CPU, pulls a
// model through the two default-allowlist hosts while an unlisted host stays
// blocked; the model survives the removal of the container. The test needs
// network access and about 1.5 GB of download; the amd64 run is not a
// verification of the architecture claim until a CI run with a visible pass.
func TestE2E_SandboxDevcontainer_Ollama(t *testing.T) {
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
	defer cancel()
	dir := generateSandboxProjectWith(ctx, t, map[string]string{
		"devcontainer.features.ollama.enabled": "true",
		"devcontainer.sandbox.egress.enabled":  "true",
		"devcontainer.sandbox.onUnavailable":   "warn",
	})
	image := buildProjectImage(ctx, t, dir, "uboot-e2e-ollama")

	out, code := dockerOutput(ctx, t, "run", "--rm", image, "bash", "-c", ollamaImageScript)
	if code != 0 {
		t.Fatalf("image checks exit %d:\n%s", code, out)
	}
	for _, want := range []string{"OLLAMA_VERSION_PINNED", "OLLAMA_LIB_PRESENT", "NO_GPU_LIBRARIES", "NO_ZSTD", "STORE_OWNED_BY_VSCODE", "NO_OLLAMA_ENV"} {
		if !strings.Contains(out, want) {
			t.Errorf("image checks lack %s:\n%s", want, out)
		}
	}

	mounts := generatedMounts(t, dir)
	if len(mounts) != 2 || !strings.Contains(mounts[1], "target=/home/vscode/.ollama,type=volume") {
		t.Fatalf("mounts = %v, want workspace volume and model volume", mounts)
	}
	id := "ollama" + strconv.FormatInt(time.Now().UnixNano(), 36)
	t.Cleanup(func() {
		for _, m := range mounts {
			src := strings.TrimPrefix(strings.Split(strings.ReplaceAll(m, "${devcontainerId}", id), ",")[0], "source=")
			_, _ = dockerOutput(context.Background(), t, "volume", "rm", "-f", src)
		}
	})
	start := func(name string, withRunArgs bool) {
		args := []string{"run", "-d", "--name", name}
		for _, m := range mounts {
			args = append(args, "--mount", strings.ReplaceAll(m, "${devcontainerId}", id))
		}
		if withRunArgs {
			args = append(args, generatedRunArgs(t, dir)...)
		}
		args = append(args, image, "sleep", "1800")
		t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rm", "-f", name) })
		if out, code := dockerOutput(ctx, t, args...); code != 0 {
			t.Fatalf("start %s: %s", name, out)
		}
	}

	first := "uboot-e2e-ollama-" + id
	start(first, true)
	out, code = dockerExec(ctx, t, first, ollamaServerScript)
	if code != 0 {
		t.Errorf("server script exit %d:\n%s", code, out)
	}
	for _, want := range []string{"LOOPBACK_ANSWERS", "LISTEN:127.0.0.1:11434 \n", "CONTAINER_ADDRESS_REFUSED", "LIBRARY_CPU", "PULL_OK", "MODEL_LISTED", "KEY_IN_MODEL_STORE", "RUN_OK", "FOREIGN_BLOCKED"} {
		if !strings.Contains(out, want) {
			t.Errorf("server script output lacks %q:\n%s", want, out)
		}
	}
	for _, bad := range []string{"CONTAINER_ADDRESS_REACHABLE", "FOREIGN_REACHABLE", "PULL_FAILED"} {
		if strings.Contains(out, bad) {
			t.Errorf("server script output has %q:\n%s", bad, out)
		}
	}

	if out, code := dockerOutput(ctx, t, "rm", "-f", first); code != 0 {
		t.Fatalf("remove %s: %s", first, out)
	}
	second := "uboot-e2e-ollama-again-" + id
	start(second, false)
	out, code = dockerExec(ctx, t, second, ollamaRestartScript)
	if code != 0 || !strings.Contains(out, "MODEL_SURVIVES_CONTAINER_REMOVAL") {
		t.Errorf("model volume after removing the container: exit %d:\n%s", code, out)
	}
}

// ADR-0024 negative probe: the build stops before the extraction when the
// pinned checksum does not match; `--build-arg` replaces both architecture
// values, so the probe holds on `amd64` and `arm64` hosts.
func TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild(t *testing.T) {
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
	defer cancel()
	dir := generateSandboxProjectWith(ctx, t, map[string]string{"devcontainer.features.ollama.enabled": "true"})
	zeros := strings.Repeat("0", 64)
	const tag = "uboot-e2e-ollama-pin-negative:probe"
	pin := "OLLAMA_SHA256_AMD64=" + zeros
	before := containerIDs(ctx, t)
	// The classic builder (the `test-docker-tools` image has no buildx) keeps
	// the container of the failed step unless `--force-rm` is set.
	out, code := dockerOutput(ctx, t, "build", "-q", "--force-rm", "-t", tag,
		"--build-arg", pin, "--build-arg", "OLLAMA_SHA256_ARM64="+zeros,
		filepath.Join(dir, ".devcontainer"))
	t.Cleanup(func() {
		left := newContainersWithCommand(context.Background(), t, before, pin)
		for _, id := range left {
			if rmOut, rmCode := dockerOutput(context.Background(), t, "rm", "-f", id); rmCode != 0 {
				t.Errorf("the failed build left container %.12s with the zero pin; removing it failed: %s", id, rmOut)
				continue
			}
			t.Errorf("the failed build left container %.12s with the zero pin; cleanup removed it", id)
		}
	})
	if code == 0 {
		t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", tag) })
		t.Fatalf("build with a wrong checksum succeeded, want a failure before the extraction:\n%s", out)
	}
	if !strings.Contains(out, "FAILED") {
		t.Errorf("build output lacks the checksum failure:\n%s", out)
	}
}

// containerIDs lists the IDs of all containers (`docker ps -a`) at call time.
func containerIDs(ctx context.Context, t *testing.T) map[string]bool {
	t.Helper()
	out, code := dockerOutput(ctx, t, "ps", "-a", "-q", "--no-trunc")
	if code != 0 {
		t.Fatalf("docker ps: %s", out)
	}
	ids := map[string]bool{}
	for _, id := range strings.Fields(out) {
		ids[id] = true
	}
	return ids
}

// newContainersWithCommand returns the containers that are not in `before`
// and whose Path, Args or Env contain needle; every other container is left out.
// Two runs against the same daemon at the same time are not told apart: both
// would see the other's container as new.
func newContainersWithCommand(ctx context.Context, t *testing.T, before map[string]bool, needle string) []string {
	t.Helper()
	var found []string
	for id := range containerIDs(ctx, t) {
		if before[id] {
			continue
		}
		out, code := dockerOutput(ctx, t, "inspect", "--format", "{{.Path}} {{.Args}} {{.Config.Env}}", id)
		if code == 0 && strings.Contains(out, needle) {
			found = append(found, id)
		}
	}
	return found
}
