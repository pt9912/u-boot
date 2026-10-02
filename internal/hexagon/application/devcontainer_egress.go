package application

import (
	"fmt"
	"path/filepath"
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

// egressPodmanHosts are added with `nestedRuntime: podman`.
func egressPodmanHosts() []string {
	return []string{"registry-1.docker.io", "auth.docker.io", "production.cloudflare.docker.com", "ghcr.io"}
}

// egressFeatureHosts are added per enabled devcontainer feature.
func egressFeatureHosts() map[string][]string {
	return map[string][]string{
		"node": {"registry.npmjs.org"},
		"go":   {"proxy.golang.org", "sum.golang.org", "storage.googleapis.com"},
		"java": {"repo.maven.apache.org", "repo1.maven.org"},
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

// egressHosts computes the effective allowlist baked into the startup
// script: default base + podman + enabled features + clone source host
// + the user's `egress.allow`; sorted and de-duplicated (ADR-0012).
func egressHosts(dc *ubootYAMLDevcontainer, cloneURL string) []string {
	set := map[string]bool{}
	add := func(hosts ...string) {
		for _, h := range hosts {
			if h != "" {
				set[h] = true
			}
		}
	}
	add(egressBaseHosts()...)
	if nestedPodman(dc) {
		add(egressPodmanHosts()...)
	}
	if dc != nil {
		byFeature := egressFeatureHosts()
		for name, f := range dc.Features {
			if f.Enabled != nil && *f.Enabled {
				add(byFeature[name]...)
			}
		}
	}
	add(cloneSourceHost(cloneURL))
	if dc != nil && dc.Sandbox != nil && dc.Sandbox.Egress != nil {
		add(dc.Sandbox.Egress.Allow...)
	}
	hosts := make([]string, 0, len(set))
	for h := range set {
		hosts = append(hosts, h)
	}
	sort.Strings(hosts)
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
	logger.Info("config set: updated", "path", req.Path.String(), "old", oldValue, "new", newValue)
	resp := driving.ConfigSetResponse{Path: req.Path, OldValue: oldValue, NewValue: newValue}
	if recorder != nil {
		resp.PlannedFiles = mapCaptureToPlannedFiles(recorder.Captured(), req.BaseDir)
	}
	return resp, nil
}
