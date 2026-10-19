package application

import (
	"fmt"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// egressEnabled reports whether the LH-FA-DEV-008 egress restriction
// is switched on (`devcontainer.sandbox.egress.enabled: true`).
func egressEnabled(dc *ubootYAMLDevcontainer) bool {
	return dc != nil && dc.Sandbox != nil && dc.Sandbox.Egress != nil &&
		dc.Sandbox.Egress.Enabled != nil && *dc.Sandbox.Egress.Enabled
}

// validateEgressHosts checks every `egress.allow` entry.
func validateEgressHosts(eg *ubootYAMLSandboxEgress) error {
	if eg == nil {
		return nil
	}
	for _, h := range eg.Allow {
		if _, err := domain.NewEgressHost(h); err != nil {
			return err
		}
	}
	return nil
}

// normaliseEgressHosts validates and de-duplicates host names,
// keeping the first-seen order.
func normaliseEgressHosts(raw []string) ([]string, error) {
	seen := map[string]bool{}
	out := make([]string, 0, len(raw))
	for _, r := range raw {
		h, err := domain.NewEgressHost(r)
		if err != nil {
			return nil, err
		}
		if !seen[h] {
			seen[h] = true
			out = append(out, h)
		}
	}
	return out, nil
}

// egressBaseHosts is the common part of the default allowlist
// (documented in spec/spezifikation.md).
func egressBaseHosts() []string {
	return []string{
		"github.com", "api.github.com", "codeload.github.com",
		"raw.githubusercontent.com", "objects.githubusercontent.com",
		"deb.debian.org", "security.debian.org",
	}
}

// egressRegistryHosts are added with `nestedRuntime: podman` or `docker`;
// Docker Hub serves layer blobs from the cloudflare or the cloudfront host.
func egressRegistryHosts() []string {
	return []string{"registry-1.docker.io", "auth.docker.io", "production.cloudflare.docker.com", "production.cloudfront.docker.com", "ghcr.io"}
}

// egressDockerAptHost is the Docker APT repository that stays configured
// in the image with `nestedRuntime: docker`; `apt-get update` in the
// restricted container needs it.
func egressDockerAptHost() string { return "download.docker.com" }

// egressFeatureHosts are added per enabled devcontainer feature. The
// `claude-code` row lists the API host (measured at runtime) and the npm
// registry (the required `node` feature brings npm; its use at run time is
// not measured). The sign-in hosts and the hosts of telemetry, updates, MCP
// and plugins are not listed; the user adds them to `egress.allow`. The
// `ollama` row lists the manifest host and the blob host of the model pull
// (measured for two models at one time and place; `ollama.com` is not
// needed and the install hosts are used by the image build only).
func egressFeatureHosts() map[string][]string {
	return map[string][]string{
		"node":        {"registry.npmjs.org"},
		"go":          {"proxy.golang.org", "sum.golang.org", "storage.googleapis.com"},
		"java":        {"repo.maven.apache.org", "repo1.maven.org"},
		"python":      {"pypi.org", "files.pythonhosted.org"},
		"homebrew":    {"formulae.brew.sh", "ghcr.io", "pkg-containers.githubusercontent.com"},
		"sdkman":      {"api.sdkman.io", "broker.sdkman.io", "repo.maven.apache.org", "repo1.maven.org"},
		"claude-code": {"registry.npmjs.org", "api.anthropic.com"},
		"ollama":      {"registry.ollama.ai", "dd20bb891979d25aebc8bec07b2b3bbc.r2.cloudflarestorage.com"},
	}
}

// cloneSourceHost extracts the host of a clone URL (https://host/…,
// ssh://user@host/…, or scp-like user@host:path); "" if unparsable.
func cloneSourceHost(raw string) string {
	if raw == "" {
		return ""
	}
	rest := raw
	if _, after, ok := strings.Cut(rest, "://"); ok {
		rest, _, _ = strings.Cut(after, "/")
	}
	if _, after, ok := strings.Cut(rest, "@"); ok {
		rest = after
	}
	host, _, _ := strings.Cut(rest, ":")
	return strings.ToLower(host)
}

// egressHostSources computes the effective allowlist with the origin of
// every host: default base + registries (podman, docker) + Docker APT host
// (docker) + enabled features + clone source host + the user's
// `egress.allow`; sorted by host and de-duplicated (ADR-0012). A host
// that several origins contribute lists each of them once, in the order
// base, runtime, docker-apt, features (by key), clone, user. The startup
// script ([egressHosts]) and `u-boot egress list` both take this list.
func egressHostSources(dc *ubootYAMLDevcontainer, cloneURL string) []driving.EgressHostEntry {
	origins := map[string][]string{}
	add := func(source string, hosts ...string) {
		for _, h := range hosts {
			if h != "" && !slices.Contains(origins[h], source) {
				origins[h] = append(origins[h], source)
			}
		}
	}
	add(driving.EgressSourceBase, egressBaseHosts()...)
	if nestedPodman(dc) || nestedDocker(dc) {
		add(driving.EgressSourceRuntime, egressRegistryHosts()...)
	}
	if nestedDocker(dc) {
		add(driving.EgressSourceDockerAPT, egressDockerAptHost())
	}
	byFeature := egressFeatureHosts()
	for _, name := range enabledFeatureNames(dc) {
		add(driving.EgressSourceFeaturePrefix+name, byFeature[name]...)
	}
	add(driving.EgressSourceClone, cloneSourceHost(cloneURL))
	if dc != nil && dc.Sandbox != nil && dc.Sandbox.Egress != nil {
		add(driving.EgressSourceUser, dc.Sandbox.Egress.Allow...)
	}
	entries := make([]driving.EgressHostEntry, 0, len(origins))
	for h, sources := range origins {
		entries = append(entries, driving.EgressHostEntry{Host: h, Sources: sources})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Host < entries[j].Host })
	return entries
}

// enabledFeatureNames returns the keys of the devcontainer features with
// `enabled: true`, sorted; nil for a nil or feature-less `dc`.
func enabledFeatureNames(dc *ubootYAMLDevcontainer) []string {
	if dc == nil {
		return nil
	}
	names := make([]string, 0, len(dc.Features))
	for name, f := range dc.Features {
		if f.Enabled != nil && *f.Enabled {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

// egressHosts is the allowlist baked into the startup script: the host
// names of [egressHostSources], sorted and without duplicates.
func egressHosts(dc *ubootYAMLDevcontainer, cloneURL string) []string {
	entries := egressHostSources(dc, cloneURL)
	hosts := make([]string, 0, len(entries))
	for _, e := range entries {
		hosts = append(hosts, e.Host)
	}
	return hosts
}

// setEgressAllow handles the list path
// `devcontainer.sandbox.egress.allow` (LH-FA-DEV-008): the value is a
// comma-separated host list that is validated, appended to the
// existing list and de-duplicated; u-boot.yaml is marshal-rewritten
// (comments are lost, like for `featureSources.allow`).
func (s *ConfigService) setEgressAllow(
	req driving.ConfigSetRequest, cfg ubootYAMLConfig,
) (driving.ConfigSetResponse, error) {
	logger := s.logger
	if req.SilenceLogger {
		logger = noopLogger{}
	}
	incoming, err := parseFeatureSourcesArgument(req.Value)
	if err != nil {
		return driving.ConfigSetResponse{}, fmt.Errorf("%w: %s: %v", driving.ErrConfigValueInvalid, req.Path, err)
	}
	var existing []string
	if cfg.Devcontainer != nil && cfg.Devcontainer.Sandbox != nil && cfg.Devcontainer.Sandbox.Egress != nil {
		existing = cfg.Devcontainer.Sandbox.Egress.Allow
	}
	merged, err := normaliseEgressHosts(append(append([]string{}, existing...), incoming...))
	if err != nil {
		return driving.ConfigSetResponse{}, fmt.Errorf("%w: %s: %w", driving.ErrConfigValueInvalid, req.Path, err)
	}
	oldValue, newValue := strings.Join(existing, ","), strings.Join(merged, ",")
	if oldValue == newValue {
		logger.Debug("config set: no-op", "path", req.Path.String(), "value", newValue)
		return driving.ConfigSetResponse{Path: req.Path, OldValue: oldValue, NewValue: newValue}, nil
	}
	if cfg.Devcontainer == nil {
		cfg.Devcontainer = &ubootYAMLDevcontainer{}
	}
	if cfg.Devcontainer.Sandbox == nil {
		cfg.Devcontainer.Sandbox = &ubootYAMLDevcontainerSandbox{}
	}
	if cfg.Devcontainer.Sandbox.Egress == nil {
		cfg.Devcontainer.Sandbox.Egress = &ubootYAMLSandboxEgress{}
	}
	cfg.Devcontainer.Sandbox.Egress.Allow = merged
	rewritten, err := s.yaml.Marshal(cfg)
	if err != nil {
		return driving.ConfigSetResponse{}, fmt.Errorf("%w: marshal u-boot.yaml: %v", driving.ErrConfigSchemaInvalid, err)
	}
	fs, recorder := s.selectFS(req.PreviewMode)
	path := filepath.Join(req.BaseDir, "u-boot.yaml")
	if err := fs.WriteFile(path, rewritten, defaultFileMode); err != nil {
		return driving.ConfigSetResponse{}, fmt.Errorf("%w: write %q: %w", driving.ErrConfigFileSystem, path, err)
	}
	logWrite(logger, req.PreviewMode, "config set: updated", "config set: would update", "path", req.Path.String(), "old", oldValue, "new", newValue)
	resp := driving.ConfigSetResponse{Path: req.Path, OldValue: oldValue, NewValue: newValue}
	if recorder != nil {
		resp.PlannedFiles = mapCaptureToPlannedFiles(recorder.Captured(), req.BaseDir)
	}
	return resp, nil
}
