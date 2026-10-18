package application

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// apparmorUsernsSysctl is the Ubuntu 24.04+ switch that blocks
// unprivileged user namespaces for unconfined processes.
const apparmorUsernsSysctl = "/proc/sys/kernel/apparmor_restrict_unprivileged_userns"

const (
	// checkIDSandboxRuntime is the LH-FA-DEV-007 check: state of the
	// nested-runtime prerequisites that can be determined from a Linux
	// host outside a container (`/dev/fuse`) plus the profile/runtime
	// consistency.
	checkIDSandboxRuntime = "devcontainer.sandbox.runtime"

	// checkIDSandboxEgress is the LH-FA-DEV-008 check: the egress
	// restriction is only meaningful together with the sandbox
	// profile; the state of the NET_ADMIN capability at the container
	// start is reported by [checkIDSandboxFallbackEgress].
	checkIDSandboxEgress = "devcontainer.sandbox.egress"

	// checkIDSandboxCredentials is the LH-FA-DEV-009 check: git and
	// agent credentials must not sit in plain text in a project file,
	// an https clone needs a runtime token source, and an active agent
	// feature needs the reference of its host variable.
	checkIDSandboxCredentials = "devcontainer.sandbox.credentials"
)

// credentialPatterns detect plaintext credentials in project files
// (LH-FA-DEV-009): userinfo in an http(s) URL, the common git token
// prefixes / assignments, and Anthropic keys. The `sk-ant-` expression is
// the one the Claude Code CLI 2.1.302 itself uses to find a key in text
// (read from the installed binary); the assignment pattern covers the
// three variables of the Anthropic environment documentation.
// References such as `${localEnv:GIT_TOKEN}` do not match.
func credentialPatterns() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`https?://[^/\s:@"']+(:[^/\s@"']*)?@`),
		regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9]{20,}`),
		regexp.MustCompile(`\bgithub_pat_[A-Za-z0-9_]{20,}`),
		regexp.MustCompile(`\bglpat-[A-Za-z0-9_-]{20,}`),
		regexp.MustCompile(`(?i)\b(?:git|gh|github|gitlab)_?token\s*[=:]\s*["']?[A-Za-z0-9_\-]{12,}`),
		regexp.MustCompile(`(?:^|[^A-Za-z0-9_-])sk-ant-[A-Za-z0-9_-]{40,}`),
		regexp.MustCompile(`\b(?:ANTHROPIC_API_KEY|ANTHROPIC_AUTH_TOKEN|CLAUDE_CODE_OAUTH_TOKEN)\b["']?\s*[=:]\s*["']?[A-Za-z0-9_\-]{12,}`),
	}
}

// credentialScanFiles are the project files scanned for plaintext
// credentials (relative to the project root).
func credentialScanFiles() []string {
	return []string{
		"u-boot.yaml", "compose.yaml", ".env.example",
		".devcontainer/devcontainer.json", ".devcontainer/Dockerfile", ".devcontainer/sandbox-init.sh", ".devcontainer/docker-init.sh",
	}
}

// nestedRuntimeSet reports whether a nested runtime (`podman` or
// `docker`) is configured.
func nestedRuntimeSet(dc *ubootYAMLDevcontainer) bool {
	return nestedPodman(dc) || nestedDocker(dc)
}

// checkSandboxRuntime implements LH-FA-DEV-007. Skips (OK) without the
// sandbox profile or without a nested runtime; warns when
// `nestedRuntime: podman` has no effect; checks `/dev/fuse` on Linux
// hosts outside a container only (on macOS/Colima the engine runs in a
// VM, so the host cannot tell; inside a container the device is the
// container's own, so the item is info and the state of a started
// sandbox container is reported by `devcontainer.sandbox.fallback.vfs`).
// A missing `/dev/fuse` on the host is an error (exit 11) under both
// `onUnavailable` values and carries no fallback flag; `/dev/net/tun` is
// not checked.
func (s *DoctorService) checkSandboxRuntime(_ context.Context, baseDir string) domain.Diagnostic {
	cfg, err := s.loadUbootYAML(baseDir)
	if err != nil || cfg.Devcontainer == nil || !nestedRuntimeSet(cfg.Devcontainer) && !profileIsSandbox(cfg.Devcontainer) {
		return domain.Diagnostic{ID: checkIDSandboxRuntime, Severity: domain.SeverityOK,
			Message: "No sandbox nested runtime configured; check skipped."}
	}
	dc := cfg.Devcontainer
	if nestedRuntimeSet(dc) && !profileIsSandbox(dc) {
		return domain.Diagnostic{ID: checkIDSandboxRuntime, Severity: domain.SeverityWarn,
			Message: fmt.Sprintf("devcontainer.sandbox.nestedRuntime is `%s` but devcontainer.profile is not `sandbox`; the setting has no effect (LH-FA-DEV-007).", dc.Sandbox.NestedRuntime),
			Hint:    "Run `u-boot config set devcontainer.profile sandbox` and `u-boot generate devcontainer`, or set nestedRuntime to `none`."}
	}
	if nestedDocker(dc) {
		return domain.Diagnostic{ID: checkIDSandboxRuntime, Severity: domain.SeverityOK,
			Message: "Nested Docker daemon configured; it needs `--privileged`, which removes the isolation from the container host; whether the daemon became ready at the container start is reported by devcontainer.sandbox.fallback.daemon (LH-FA-DEV-007)."}
	}
	if !nestedPodman(dc) {
		return domain.Diagnostic{ID: checkIDSandboxRuntime, Severity: domain.SeverityOK,
			Message: "Sandbox profile without nested runtime; nothing to check."}
	}
	if linux, _ := s.fs.Exists("/proc/sys/kernel/osrelease"); !linux {
		return domain.Diagnostic{ID: checkIDSandboxRuntime, Severity: domain.SeverityOK,
			Message: "Nested Podman prerequisites cannot be determined from this host (container engine runs in a VM); the container engine refuses to start without /dev/fuse and /dev/net/tun on the container host, and the container start script checks user namespaces."}
	}
	if restricted, _ := s.fs.ReadFile(apparmorUsernsSysctl); strings.TrimSpace(string(restricted)) == "1" {
		return domain.Diagnostic{ID: checkIDSandboxRuntime, Severity: domain.SeverityError,
			Message: "kernel.apparmor_restrict_unprivileged_userns=1 on this host blocks nested user namespaces; rootless Podman in the sandbox cannot start (LH-FA-DEV-007, exit 11 at container start).",
			Hint:    "Set `sysctl kernel.apparmor_restrict_unprivileged_userns=0` on the container host, or set `devcontainer.sandbox.nestedRuntime` to `none`."}
	}
	if s.inContainer() {
		return domain.Diagnostic{ID: checkIDSandboxRuntime, Severity: domain.SeverityInfo,
			Message: "/dev/fuse not checked: u-boot is running inside a container, and its /dev/fuse is the device of this container, not of the container host; a sandbox container started without the device is reported by devcontainer.sandbox.fallback.vfs (LH-FA-DEV-007).",
			Hint:    "Run `u-boot doctor` on the container host to check /dev/fuse."}
	}
	if fuse, _ := s.fs.Exists("/dev/fuse"); fuse {
		return domain.Diagnostic{ID: checkIDSandboxRuntime, Severity: domain.SeverityOK,
			Message: "/dev/fuse is available for nested Podman."}
	}
	// A missing device on the container host is not a degradation
	// step: the engine refuses to start the container, whatever
	// `onUnavailable` says. The item is an error and no fallback.
	return domain.Diagnostic{ID: checkIDSandboxRuntime, Severity: domain.SeverityError,
		Message: "/dev/fuse is not available on this host; the container engine refuses to start the sandbox container with `--device=/dev/fuse` (LH-FA-DEV-007, exit 11).",
		Hint:    "Load the fuse module (`modprobe fuse`) on the container host, or set `devcontainer.sandbox.nestedRuntime` to `none`."}
}

// checkSandboxCredentials implements LH-FA-DEV-009 (warn only). It is
// skipped without the sandbox profile. All warn items (plaintext
// credential, https clone without token source, active agent feature
// without the reference of its host variable) carry
// [domain.Diagnostic.SecurityRelaxation].
func (s *DoctorService) checkSandboxCredentials(_ context.Context, baseDir string) domain.Diagnostic {
	cfg, err := s.loadUbootYAML(baseDir)
	if err != nil || !profileIsSandbox(cfg.Devcontainer) {
		return domain.Diagnostic{ID: checkIDSandboxCredentials, Severity: domain.SeverityOK,
			Message: "Sandbox profile not active; credential check skipped."}
	}
	if found := s.plaintextCredentialFiles(baseDir); len(found) > 0 {
		return domain.Diagnostic{ID: checkIDSandboxCredentials, Severity: domain.SeverityWarn,
			Message: fmt.Sprintf("possible plaintext credential in: %s (LH-FA-DEV-009).", strings.Join(found, ", ")),
			Hint: "Remove it from the file, revoke the credential, and pass it at runtime: a short-lived repository-scoped token via $GIT_TOKEN, " +
				"an agent key via its host environment variable (for example $ANTHROPIC_API_KEY).",
			SecurityRelaxation: true}
	}
	agentVars := agentCredentialEnv(cfg.Devcontainer)
	origin := effectiveCloneSource(s.fs, baseDir, cfg.Devcontainer)
	httpsClone := strings.HasPrefix(origin, "https://") || strings.HasPrefix(origin, "http://")
	if httpsClone && !s.hasTokenSource(baseDir) {
		return domain.Diagnostic{ID: checkIDSandboxCredentials, Severity: domain.SeverityWarn,
			Message:            "sandbox profile clones over https but devcontainer.json has no runtime token source (remoteEnv GIT_TOKEN or a read-only secret mount) (LH-FA-DEV-009).",
			Hint:               "Run `u-boot generate devcontainer` to restore the GIT_TOKEN pass-through; never store tokens in project files.",
			SecurityRelaxation: true}
	}
	if missing := s.missingAgentSources(baseDir, agentVars); len(missing) > 0 {
		return domain.Diagnostic{ID: checkIDSandboxCredentials, Severity: domain.SeverityWarn,
			Message: fmt.Sprintf("an active agent feature needs its credential at runtime but devcontainer.json has no reference for: %s (LH-FA-DEV-009).",
				strings.Join(missing, ", ")),
			Hint:               "Run `u-boot generate devcontainer` to restore the `${localEnv:<NAME>}` pass-through and export the variable on the host; never store the key in project files.",
			SecurityRelaxation: true}
	}
	msg := "No plaintext git or agent credentials found; no https clone source needs a token."
	if httpsClone {
		msg = "No plaintext git or agent credentials found; a runtime token source is configured."
	}
	if len(agentVars) > 0 {
		msg += " The agent credential is passed through " + strings.Join(agentVars, ", ") + " (reference only)."
	}
	return domain.Diagnostic{ID: checkIDSandboxCredentials, Severity: domain.SeverityOK, Message: msg}
}

// missingAgentSources returns the agent variables whose
// `${localEnv:<NAME>}` reference is absent from devcontainer.json (the
// file counts as absent when it cannot be read).
func (s *DoctorService) missingAgentSources(baseDir string, vars []string) []string {
	if len(vars) == 0 {
		return nil
	}
	body, err := s.fs.ReadFile(filepath.Join(baseDir, ".devcontainer", "devcontainer.json"))
	text := string(body)
	var missing []string
	for _, name := range vars {
		if err != nil || !strings.Contains(text, "${localEnv:"+name+"}") {
			missing = append(missing, name)
		}
	}
	return missing
}

// plaintextCredentialFiles returns the sorted project files that match
// a credential pattern.
func (s *DoctorService) plaintextCredentialFiles(baseDir string) []string {
	patterns := credentialPatterns()
	var found []string
	for _, rel := range credentialScanFiles() {
		body, err := s.fs.ReadFile(filepath.Join(baseDir, rel))
		if err != nil {
			continue
		}
		for _, re := range patterns {
			if re.Match(body) {
				found = append(found, rel)
				break
			}
		}
	}
	sort.Strings(found)
	return found
}

// hasTokenSource reports whether devcontainer.json passes a token at
// runtime: a `GIT_TOKEN` environment reference or a read-only secret
// mount below /run/secrets.
func (s *DoctorService) hasTokenSource(baseDir string) bool {
	body, err := s.fs.ReadFile(filepath.Join(baseDir, ".devcontainer", "devcontainer.json"))
	if err != nil {
		return false
	}
	text := string(body)
	return strings.Contains(text, "${localEnv:GIT_TOKEN}") ||
		(strings.Contains(text, "/run/secrets") && strings.Contains(text, "readonly"))
}

// checkSandboxEgress implements the static part of LH-FA-DEV-008. Both
// warn items (egress without the sandbox profile, egress together with
// `nestedRuntime: docker`) carry
// [domain.Diagnostic.SecurityRelaxation].
func (s *DoctorService) checkSandboxEgress(_ context.Context, baseDir string) domain.Diagnostic {
	cfg, err := s.loadUbootYAML(baseDir)
	if err != nil || !egressEnabled(cfg.Devcontainer) {
		return domain.Diagnostic{ID: checkIDSandboxEgress, Severity: domain.SeverityOK,
			Message: "Egress restriction not enabled; check skipped."}
	}
	if !profileIsSandbox(cfg.Devcontainer) {
		return domain.Diagnostic{ID: checkIDSandboxEgress, Severity: domain.SeverityWarn,
			Message:            "devcontainer.sandbox.egress.enabled is set but devcontainer.profile is not `sandbox`; the setting has no effect (LH-FA-DEV-008).",
			Hint:               "Run `u-boot config set devcontainer.profile sandbox` and `u-boot generate devcontainer`, or disable the egress restriction.",
			SecurityRelaxation: true}
	}
	if nestedDocker(cfg.Devcontainer) {
		return domain.Diagnostic{ID: checkIDSandboxEgress, Severity: domain.SeverityWarn,
			Message:            "Egress restriction together with nestedRuntime `docker`: the traffic of nested containers and the iptables rules of the nested daemon bypass the restriction, and `--privileged` lets any process lift it (LH-FA-DEV-008).",
			Hint:               "Treat the restriction as a guardrail only, or set `devcontainer.sandbox.nestedRuntime` to `none` or `podman`.",
			SecurityRelaxation: true}
	}
	return domain.Diagnostic{ID: checkIDSandboxEgress, Severity: domain.SeverityOK,
		Message: "Egress restriction configured; whether the NET_ADMIN capability was granted at the container start is reported by devcontainer.sandbox.fallback.egress (guardrail, not a sandbox boundary)."}
}

// checkSandboxFallbacks reports, per configured capability, the state its
// container start script wrote into the status file (LH-FA-DEV-007,
// LH-FA-DEV-008): one item per `vfs` storage, daemon and egress
// restriction, only for the capabilities the configuration asks for. A
// state that cannot be read is an info item, never a fallback. The
// status file is a hint from the container, not proof: a missing report
// does not mean that no fallback occurred.
func (s *DoctorService) checkSandboxFallbacks(ctx context.Context, baseDir string) []domain.Diagnostic {
	cfg, err := s.loadUbootYAML(baseDir)
	if err != nil {
		return nil
	}
	probes := configuredFallbackProbes(cfg.Devcontainer, profileIsSandbox(cfg.Devcontainer))
	if len(probes) == 0 {
		return nil
	}
	reader := &fallbackReader{inspector: s.inspector, fs: s.fs, inContainer: s.inContainer(), baseDir: baseDir}
	items := make([]domain.Diagnostic, 0, len(probes))
	for _, p := range probes {
		items = append(items, fallbackDiagnostic(p, reader.read(ctx, p)))
	}
	return items
}
