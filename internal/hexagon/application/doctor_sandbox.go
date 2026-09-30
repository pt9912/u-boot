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

const (
	// checkIDSandboxRuntime is the LH-FA-DEV-007 check: state of the
	// nested-runtime prerequisites that can be determined from the
	// host (`/dev/fuse`) plus the profile/runtime consistency.
	checkIDSandboxRuntime = "devcontainer.sandbox.runtime"

	// checkIDSandboxEgress is the LH-FA-DEV-008 check: the egress
	// restriction is only meaningful together with the sandbox
	// profile; the NET_ADMIN capability itself cannot be determined
	// from the host and is verified by the container start script.
	checkIDSandboxEgress = "devcontainer.sandbox.egress"

	// checkIDSandboxCredentials is the LH-FA-DEV-009 check: git
	// credentials must not sit in plain text in a project file, and
	// an https clone needs a runtime token source.
	checkIDSandboxCredentials = "devcontainer.sandbox.credentials"
)

// credentialPatterns detect plaintext git credentials in project
// files (LH-FA-DEV-009): userinfo in an http(s) URL and the common
// token prefixes / assignments. References such as
// `${localEnv:GIT_TOKEN}` do not match.
func credentialPatterns() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`https?://[^/\s:@"']+(:[^/\s@"']*)?@`),
		regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9]{20,}`),
		regexp.MustCompile(`\bgithub_pat_[A-Za-z0-9_]{20,}`),
		regexp.MustCompile(`\bglpat-[A-Za-z0-9_-]{20,}`),
		regexp.MustCompile(`(?i)\b(?:git|gh|github|gitlab)_?token\s*[=:]\s*["']?[A-Za-z0-9_\-]{12,}`),
	}
}

// credentialScanFiles are the project files scanned for plaintext
// credentials (relative to the project root).
func credentialScanFiles() []string {
	return []string{
		"u-boot.yaml", "compose.yaml", ".env.example",
		".devcontainer/devcontainer.json", ".devcontainer/Dockerfile", ".devcontainer/sandbox-init.sh",
	}
}

// checkSandboxRuntime implements LH-FA-DEV-007 from the host's point
// of view. Skips (OK) without the sandbox profile or without a
// nested runtime; warns when `nestedRuntime: podman` has no effect;
// checks `/dev/fuse` on Linux hosts only (on macOS/Colima the engine
// runs in a VM, so the host cannot tell — the container start
// script checks there). `onUnavailable: fail` turns the finding into
// an error (exit 11).
func (s *DoctorService) checkSandboxRuntime(_ context.Context, baseDir string) domain.Diagnostic {
	cfg, err := s.loadUbootYAML(baseDir)
	if err != nil || cfg.Devcontainer == nil || !nestedPodman(cfg.Devcontainer) && !profileIsSandbox(cfg.Devcontainer) {
		return domain.Diagnostic{ID: checkIDSandboxRuntime, Severity: domain.SeverityOK,
			Message: "No sandbox nested runtime configured; check skipped."}
	}
	dc := cfg.Devcontainer
	if nestedPodman(dc) && !profileIsSandbox(dc) {
		return domain.Diagnostic{ID: checkIDSandboxRuntime, Severity: domain.SeverityWarn,
			Message: "devcontainer.sandbox.nestedRuntime is `podman` but devcontainer.profile is not `sandbox`; the setting has no effect (LH-FA-DEV-007).",
			Hint:    "Run `u-boot config set devcontainer.profile sandbox` and `u-boot generate devcontainer`, or set nestedRuntime to `none`."}
	}
	if !nestedPodman(dc) {
		return domain.Diagnostic{ID: checkIDSandboxRuntime, Severity: domain.SeverityOK,
			Message: "Sandbox profile without nested runtime; nothing to check."}
	}
	if linux, _ := s.fs.Exists("/proc/sys/kernel/osrelease"); !linux {
		return domain.Diagnostic{ID: checkIDSandboxRuntime, Severity: domain.SeverityOK,
			Message: "Nested Podman prerequisites cannot be determined from this host (container engine runs in a VM); the container start script checks /dev/fuse and user namespaces."}
	}
	if fuse, _ := s.fs.Exists("/dev/fuse"); fuse {
		return domain.Diagnostic{ID: checkIDSandboxRuntime, Severity: domain.SeverityOK,
			Message: "/dev/fuse is available for nested Podman."}
	}
	sev := domain.SeverityWarn
	msg := "/dev/fuse is not available on this host; nested Podman will fall back to vfs storage (slower) (LH-FA-DEV-007)."
	if onUnavailablePolicy(dc) == string(domain.OnUnavailableFail) {
		sev = domain.SeverityError
		msg = "/dev/fuse is not available on this host and devcontainer.sandbox.onUnavailable is `fail` (LH-FA-DEV-007)."
	}
	return domain.Diagnostic{ID: checkIDSandboxRuntime, Severity: sev, Message: msg,
		Hint: "Load the fuse module (`modprobe fuse`) or set `devcontainer.sandbox.onUnavailable` to `warn`."}
}

// checkSandboxCredentials implements LH-FA-DEV-009 (warn only). It is
// skipped without the sandbox profile.
func (s *DoctorService) checkSandboxCredentials(_ context.Context, baseDir string) domain.Diagnostic {
	cfg, err := s.loadUbootYAML(baseDir)
	if err != nil || !profileIsSandbox(cfg.Devcontainer) {
		return domain.Diagnostic{ID: checkIDSandboxCredentials, Severity: domain.SeverityOK,
			Message: "Sandbox profile not active; git credential check skipped."}
	}
	if found := s.plaintextCredentialFiles(baseDir); len(found) > 0 {
		return domain.Diagnostic{ID: checkIDSandboxCredentials, Severity: domain.SeverityWarn,
			Message: fmt.Sprintf("possible plaintext git credential in: %s (LH-FA-DEV-009).", strings.Join(found, ", ")),
			Hint:    "Remove it from the file, revoke the token, and pass a short-lived repository-scoped token at runtime via $GIT_TOKEN."}
	}
	origin := effectiveCloneSource(s.fs, baseDir, cfg.Devcontainer)
	if !strings.HasPrefix(origin, "https://") && !strings.HasPrefix(origin, "http://") {
		return domain.Diagnostic{ID: checkIDSandboxCredentials, Severity: domain.SeverityOK,
			Message: "No plaintext git credentials found; no https clone source needs a token."}
	}
	if !s.hasTokenSource(baseDir) {
		return domain.Diagnostic{ID: checkIDSandboxCredentials, Severity: domain.SeverityWarn,
			Message: "sandbox profile clones over https but devcontainer.json has no runtime token source (remoteEnv GIT_TOKEN or a read-only secret mount) (LH-FA-DEV-009).",
			Hint:    "Run `u-boot generate devcontainer` to restore the GIT_TOKEN pass-through; never store tokens in project files."}
	}
	return domain.Diagnostic{ID: checkIDSandboxCredentials, Severity: domain.SeverityOK,
		Message: "No plaintext git credentials found; a runtime token source is configured."}
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

// checkSandboxEgress implements the static part of LH-FA-DEV-008.
func (s *DoctorService) checkSandboxEgress(_ context.Context, baseDir string) domain.Diagnostic {
	cfg, err := s.loadUbootYAML(baseDir)
	if err != nil || !egressEnabled(cfg.Devcontainer) {
		return domain.Diagnostic{ID: checkIDSandboxEgress, Severity: domain.SeverityOK,
			Message: "Egress restriction not enabled; check skipped."}
	}
	if !profileIsSandbox(cfg.Devcontainer) {
		return domain.Diagnostic{ID: checkIDSandboxEgress, Severity: domain.SeverityWarn,
			Message: "devcontainer.sandbox.egress.enabled is set but devcontainer.profile is not `sandbox`; the setting has no effect (LH-FA-DEV-008).",
			Hint:    "Run `u-boot config set devcontainer.profile sandbox` and `u-boot generate devcontainer`, or disable the egress restriction."}
	}
	return domain.Diagnostic{ID: checkIDSandboxEgress, Severity: domain.SeverityOK,
		Message: "Egress restriction configured; the NET_ADMIN capability is verified when the container starts (guardrail, not a sandbox boundary)."}
}
