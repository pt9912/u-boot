package application

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// validateDevcontainer runs every schema check on the `devcontainer:`
// sub-tree: the LH-FA-DEV-003 feature checks plus the Lastenheft
// 0.3.0 checks of [validateDevcontainerSandbox]. Load-time callers
// (config, generate, doctor) use this single entry point so a new
// key family cannot be forgotten at one of the call sites.
func validateDevcontainer(dc *ubootYAMLDevcontainer) error {
	if err := validateDevcontainerFeatures(dc); err != nil {
		return err
	}
	return validateDevcontainerSandbox(dc)
}

// validateDevcontainerSandbox checks the closed value sets of
// `devcontainer.user.uid` (LH-FA-DEV-004), `devcontainer.profile`
// (LH-FA-DEV-006) and `devcontainer.sandbox.*` (LH-FA-DEV-007).
// Empty / absent values are valid (defaults apply). Errors wrap
// [domain.ErrInvalidSandboxSetting]; callers classify them as
// LH-FA-CLI-006 exit 10.
func validateDevcontainerSandbox(dc *ubootYAMLDevcontainer) error {
	if dc == nil {
		return nil
	}
	if dc.User != nil && dc.User.UID != nil {
		if _, err := domain.NewContainerUID(*dc.User.UID); err != nil {
			return err
		}
	}
	if dc.Profile != "" {
		if _, err := domain.NewDevcontainerProfile(dc.Profile); err != nil {
			return err
		}
	}
	if dc.Sandbox == nil {
		return nil
	}
	if dc.Sandbox.NestedRuntime != "" {
		if _, err := domain.NewNestedRuntime(dc.Sandbox.NestedRuntime); err != nil {
			return err
		}
	}
	if dc.Sandbox.OnUnavailable != "" {
		if _, err := domain.NewOnUnavailable(dc.Sandbox.OnUnavailable); err != nil {
			return err
		}
	}
	if dc.Sandbox.Repository != "" {
		if err := validateCloneURL(dc.Sandbox.Repository); err != nil {
			return err
		}
	}
	return validateEgressHosts(dc.Sandbox.Egress)
}

// sandboxWorkspaceRoot is the parent directory of the sandbox
// workspace mount inside the container (LH-FA-DEV-006).
const sandboxWorkspaceRoot = "/workspaces"

// cloneURLSafeRE is the character whitelist for the `origin` URL that
// lands verbatim in a generated shell command inside a JSON string:
// no quotes, spaces, `$`, backticks, `;`, `&`, `|` or backslashes.
var cloneURLSafeRE = regexp.MustCompile(`^[A-Za-z0-9._~:/@%+=-]+$`)

// profileIsSandbox reports whether the persisted profile is `sandbox`.
func profileIsSandbox(dc *ubootYAMLDevcontainer) bool {
	return dc != nil && dc.Profile == string(domain.ProfileSandbox)
}

// devcontainerTemplateData builds the template data for the two
// devcontainer templates from the project name, the parsed
// `devcontainer:` sub-tree (may be nil) and the resolved sandbox
// state. UID stays 0 (= emit nothing) for the default 1000 so the
// pre-0.3.0 output remains byte-identical (LH-FA-DEV-004).
func devcontainerTemplateData(name string, dc *ubootYAMLDevcontainer, sandbox bool, cloneURL string) templateData {
	data := templateData{Name: name}
	if dc != nil && dc.User != nil && dc.User.UID != nil && *dc.User.UID != domain.DefaultContainerUID {
		data.UID = *dc.User.UID
	}
	if !sandbox {
		return data
	}
	data.Sandbox = true
	data.WorkspaceVolume = name + "-workspace"
	data.WorkspaceFolder = sandboxWorkspaceRoot + "/" + name
	data.CloneURL = cloneURL
	var postCreate []string
	if nestedPodman(dc) {
		data.NestedPodman = true
		data.OnUnavailable = onUnavailablePolicy(dc)
		data.ContainersVolume = name + "-containers"
		data.RunArgs = append(data.RunArgs, podmanRelaxations()...)
		postCreate = append(postCreate, "sh /usr/local/bin/u-boot-sandbox-init")
	}
	if egressEnabled(dc) {
		data.Egress = true
		data.OnUnavailable = onUnavailablePolicy(dc)
		data.EgressHosts = egressHosts(dc, cloneURL)
		data.RunArgs = append(data.RunArgs, egressRelaxation())
		data.PostStart = "sudo sh /usr/local/bin/u-boot-egress-init"
	}
	data.PostCreate = composePostCreate(postCreate, cloneURL)
	return data
}

// composePostCreate joins the optional init steps and the clone step
// into one `postCreateCommand`: init steps first; a failing step
// (exit 11) stops the chain before the clone.
func composePostCreate(steps []string, cloneURL string) string {
	clone := ""
	if cloneURL != "" {
		clone = "[ -d .git ] || git clone -- " + cloneURL + " ."
	}
	switch {
	case len(steps) == 0:
		return clone
	case clone == "":
		return strings.Join(steps, " && ")
	}
	return strings.Join(steps, " && ") + " && { " + clone + "; }"
}

// egressRelaxation is the capability the egress restriction needs.
func egressRelaxation() string { return "--cap-add=NET_ADMIN" }

// validateCloneURL rejects `origin` URLs that must not be written
// into a generated file: unsafe characters (shell/JSON injection) or
// embedded credentials (LH-FA-DEV-009 — none in any generated file).
// scp-like (`git@host:path`) and `ssh://user@host/…` URLs carry only
// a user name and pass.
func validateCloneURL(raw string) error {
	if !cloneURLSafeRE.MatchString(raw) || strings.HasPrefix(raw, "-") {
		return fmt.Errorf("%w: git remote 'origin' URL contains characters that cannot be written into the devcontainer safely",
			domain.ErrInvalidSandboxSetting)
	}
	scheme, rest, hasScheme := strings.Cut(raw, "://")
	if !hasScheme {
		return nil
	}
	authority, _, _ := strings.Cut(rest, "/")
	userinfo, _, hasUser := strings.Cut(authority, "@")
	if !hasUser {
		return nil
	}
	if scheme == "http" || scheme == "https" || strings.Contains(userinfo, ":") {
		return fmt.Errorf("%w: git remote 'origin' URL embeds credentials; remove them from the remote (LH-FA-DEV-009) and pass tokens at runtime",
			domain.ErrInvalidSandboxSetting)
	}
	return nil
}

// readGitOriginURL returns the URL of remote `origin` from
// `<baseDir>/.git/config`, or "" when there is no repository, no
// `.git/config` (worktrees / submodules have a `.git` file), or no
// such remote. Best-effort discovery through the FileSystem port.
func readGitOriginURL(fs driven.FileSystem, baseDir string) string {
	path := filepath.Join(baseDir, ".git", "config")
	exists, err := fs.Exists(path)
	if err != nil || !exists {
		return ""
	}
	body, err := fs.ReadFile(path)
	if err != nil {
		return ""
	}
	inOrigin := false
	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") {
			inOrigin = line == `[remote "origin"]`
			continue
		}
		if key, value, ok := strings.Cut(line, "="); inOrigin && ok && strings.TrimSpace(key) == "url" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

// effectiveCloneSource returns the clone URL of the sandbox profile:
// `devcontainer.sandbox.repository` when set, else the `origin` URL.
func effectiveCloneSource(fs driven.FileSystem, baseDir string, dc *ubootYAMLDevcontainer) string {
	if dc != nil && dc.Sandbox != nil && dc.Sandbox.Repository != "" {
		return dc.Sandbox.Repository
	}
	return readGitOriginURL(fs, baseDir)
}

// resolveSandboxClone returns the validated clone URL for the sandbox
// profile plus the LH-FA-DEV-006 warning when there is no remote
// (no clone step is generated then). For a non-sandbox render it
// returns ("", nil, nil). A URL with credentials or unsafe
// characters is a domain error (exit 10).
func resolveSandboxClone(fs driven.FileSystem, baseDir string, sandbox bool, dc *ubootYAMLDevcontainer) (string, []driving.WarningEntry, error) {
	if !sandbox {
		return "", nil, nil
	}
	url := effectiveCloneSource(fs, baseDir, dc)
	if url == "" {
		return "", []driving.WarningEntry{{
			Code:  "LH-FA-DEV-006",
			Level: "warn",
			Message: "sandbox profile: no git remote 'origin' found; no clone step generated. " +
				"Add a remote (or set devcontainer.sandbox.repository) and run `u-boot generate devcontainer` to add it",
			Subject: ".devcontainer/devcontainer.json",
		}}, nil
	}
	if err := validateCloneURL(url); err != nil {
		return "", nil, err
	}
	return url, nil, nil
}

// nestedPodman reports whether `devcontainer.sandbox.nestedRuntime`
// is `podman` (LH-FA-DEV-007).
func nestedPodman(dc *ubootYAMLDevcontainer) bool {
	return dc != nil && dc.Sandbox != nil && dc.Sandbox.NestedRuntime == string(domain.NestedRuntimePodman)
}

// onUnavailablePolicy returns the effective degradation policy
// (`warn` default, LH-FA-DEV-007).
func onUnavailablePolicy(dc *ubootYAMLDevcontainer) string {
	if dc != nil && dc.Sandbox != nil && dc.Sandbox.OnUnavailable != "" {
		return dc.Sandbox.OnUnavailable
	}
	return string(domain.OnUnavailableWarn)
}

// podmanRelaxations returns the container-level security relaxations the
// nested rootless Podman setup needs under Docker (measured, see the
// nested-Podman ADR). Each one is reported to the user individually
// (LH-FA-DEV-006 / -007).
func podmanRelaxations() []string {
	return []string{
		"--cap-add=SYS_ADMIN",
		"--security-opt=seccomp=unconfined",
		"--security-opt=apparmor=unconfined",
		"--security-opt=systempaths=unconfined",
		"--device=/dev/fuse",
	}
}

// sandboxRelaxationWarnings returns one warning per security
// relaxation the rendered sandbox needs (nested Podman: five options,
// LH-FA-DEV-007; egress restriction: NET_ADMIN, LH-FA-DEV-008), plus
// a "no effect" notice per feature that is configured without the
// sandbox profile.
func sandboxRelaxationWarnings(dc *ubootYAMLDevcontainer, sandbox bool) []driving.WarningEntry {
	var out []driving.WarningEntry
	warn := func(code, msg string) {
		out = append(out, driving.WarningEntry{Code: code, Level: "warn", Message: msg, Subject: ".devcontainer/devcontainer.json"})
	}
	if nestedPodman(dc) {
		if !sandbox {
			warn("LH-FA-DEV-007", "devcontainer.sandbox.nestedRuntime=podman has no effect without devcontainer.profile=sandbox; no nested runtime was generated")
		} else {
			for _, r := range podmanRelaxations() {
				warn("LH-FA-DEV-007", "security relaxation for nested Podman: "+r+" (the sandbox is much weaker than without nestedRuntime)")
			}
		}
	}
	if egressEnabled(dc) {
		if !sandbox {
			warn("LH-FA-DEV-008", "devcontainer.sandbox.egress.enabled has no effect without devcontainer.profile=sandbox; no egress restriction was generated")
		} else {
			warn("LH-FA-DEV-008", "security relaxation for the egress restriction: "+egressRelaxation()+" (a guardrail, not a sandbox boundary)")
		}
	}
	return out
}
