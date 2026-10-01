package application

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// GetMany implements [driving.ConfigUseCase.GetMany]
// (slice-v1-config-multi-path-get): one read of u-boot.yaml, values in
// request order, all-or-nothing (the first failing path aborts with
// its error, which already names the path).
func (s *ConfigService) GetMany(_ context.Context, req driving.ConfigGetManyRequest) (driving.ConfigGetManyResponse, error) {
	if req.BaseDir == "" {
		return driving.ConfigGetManyResponse{}, errors.New("BaseDir is required")
	}
	if len(req.Paths) == 0 {
		return driving.ConfigGetManyResponse{}, errors.New("at least one path is required")
	}
	cfg, err := s.readUbootYAML(req.BaseDir)
	if err != nil {
		return driving.ConfigGetManyResponse{}, err
	}
	entries := make([]driving.ConfigPathValue, 0, len(req.Paths))
	for _, p := range req.Paths {
		v, err := extractConfigValue(cfg, p)
		if err != nil {
			return driving.ConfigGetManyResponse{}, err
		}
		entries = append(entries, driving.ConfigPathValue{Path: p, Value: v})
	}
	return driving.ConfigGetManyResponse{Entries: entries}, nil
}

// List implements [driving.ConfigUseCase.List]
// (slice-v1-config-list-subcommand): every whitelisted path that
// currently has a value, sorted by its dotted form. Per-service and
// per-feature paths are derived from what u-boot.yaml contains.
func (s *ConfigService) List(_ context.Context, req driving.ConfigListRequest) (driving.ConfigListResponse, error) {
	if req.BaseDir == "" {
		return driving.ConfigListResponse{}, errors.New("BaseDir is required")
	}
	cfg, err := s.readUbootYAML(req.BaseDir)
	if err != nil {
		return driving.ConfigListResponse{}, err
	}
	entries := []driving.ConfigPathValue{}
	for _, p := range listCandidatePaths(cfg) {
		if v, err := extractConfigValueRaw(cfg, p); err == nil && v != "" {
			entries = append(entries, driving.ConfigPathValue{Path: p, Value: v})
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path.String() < entries[j].Path.String() })
	return driving.ConfigListResponse{Entries: entries}, nil
}

// listCandidatePaths builds every path that could have a value in cfg:
// the fixed scalar/list paths plus `services.<svc>.enabled` and the
// three `devcontainer.features.<name>.*` leaves for every key present.
func listCandidatePaths(cfg ubootYAMLConfig) []domain.ConfigPath {
	var out []domain.ConfigPath
	add := func(raw string) {
		if p, err := domain.NewConfigPath(raw); err == nil {
			out = append(out, p)
		}
	}
	for raw := range scalarConfigPathNames() {
		add(raw)
	}
	for name := range cfg.Services {
		add("services." + name + ".enabled")
	}
	if cfg.Devcontainer != nil {
		for name := range cfg.Devcontainer.Features {
			for _, leaf := range []string{"enabled", "source", "version"} {
				add("devcontainer.features." + name + "." + leaf)
			}
		}
	}
	return out
}

// scalarConfigPathNames returns the fixed dotted paths (also the list
// paths) the domain whitelists.
func scalarConfigPathNames() map[string]struct{} {
	out := map[string]struct{}{}
	for _, raw := range domain.FixedConfigPaths() {
		out[raw] = struct{}{}
	}
	return out
}

// SetMany implements [driving.ConfigUseCase.SetMany]
// (slice-v1-config-multi-path-set). All items are applied in order on
// one in-memory document — coercion, schema re-validation and the
// LH-FA-DEV-003 allowlist enforcement per item — and only then written
// once; any error aborts before the write, so u-boot.yaml stays
// byte-identical (atomic).
func (s *ConfigService) SetMany(_ context.Context, req driving.ConfigSetManyRequest) (driving.ConfigSetManyResponse, error) {
	if req.BaseDir == "" {
		return driving.ConfigSetManyResponse{}, errors.New("BaseDir is required")
	}
	if len(req.Items) == 0 {
		return driving.ConfigSetManyResponse{}, errors.New("at least one path/value pair is required")
	}
	logger := s.logger
	if req.SilenceLogger {
		logger = noopLogger{}
	}
	for _, it := range req.Items {
		if !it.Path.WriteAllowed {
			return driving.ConfigSetManyResponse{}, writeRejectedError(it.Path)
		}
	}
	body, cfg, err := s.readUbootYAMLBody(req.BaseDir)
	if err != nil {
		return driving.ConfigSetManyResponse{}, err
	}
	cur, curCfg := body, cfg
	resp := driving.ConfigSetManyResponse{Entries: make([]driving.ConfigSetEntry, 0, len(req.Items))}
	for _, it := range req.Items {
		old := extractConfigValueLenient(curCfg, it.Path)
		next, nextCfg, err := s.applySetItem(cur, curCfg, it, req.AllowExternalFeatureSources)
		if err != nil {
			return driving.ConfigSetManyResponse{}, err
		}
		resp.Entries = append(resp.Entries, driving.ConfigSetEntry{
			Path: it.Path, OldValue: old, NewValue: extractConfigValueLenient(nextCfg, it.Path)})
		resp.Warnings = append(resp.Warnings, maybeWarnOrphanFeatureActivation(logger, nextCfg, it.Path)...)
		cur, curCfg = next, nextCfg
	}
	if bytes.Equal(cur, body) {
		logger.Debug("config set (many): no-op", "items", len(req.Items))
		return resp, nil
	}
	fs, recorder := s.selectFS(req.PreviewMode)
	path := filepath.Join(req.BaseDir, "u-boot.yaml")
	if err := fs.WriteFile(path, cur, defaultFileMode); err != nil {
		return driving.ConfigSetManyResponse{}, fmt.Errorf("%w: write %q: %w", driving.ErrConfigFileSystem, path, err)
	}
	logger.Info("config set (many): updated", "items", len(req.Items))
	if recorder != nil {
		resp.PlannedFiles = mapCaptureToPlannedFiles(recorder.Captured(), req.BaseDir)
	}
	return resp, nil
}

// applySetItem applies one path/value pair to the in-memory document
// and returns the new bytes and parsed config. An unchanged value
// returns the inputs untouched.
func (s *ConfigService) applySetItem(
	body []byte, cfg ubootYAMLConfig, it driving.ConfigSetItem, extraSources []string,
) ([]byte, ubootYAMLConfig, error) {
	switch it.Path.Kind {
	case domain.ConfigDevcontainerFeatureSourcesAllow:
		next, changed, err := planFeatureSourcesAllow(cfg, it, extraSources)
		return s.marshalIfChanged(body, cfg, next, changed, err)
	case domain.ConfigDevcontainerSandboxEgressAllow:
		next, changed, err := planEgressAllow(cfg, it)
		return s.marshalIfChanged(body, cfg, next, changed, err)
	}
	coerced, formatted, err := coerceConfigValue(it.Path, it.Value)
	if err != nil {
		return nil, cfg, err
	}
	if extractConfigValueLenient(cfg, it.Path) == formatted {
		return body, cfg, nil
	}
	patched, err := s.yaml.PatchScalar(body, configPathToYAMLPath(it.Path), coerced)
	if err != nil {
		return nil, cfg, fmt.Errorf("%w: PatchScalar(%s): %v", driving.ErrConfigSchemaInvalid, it.Path, err)
	}
	var patchedCfg ubootYAMLConfig
	if err := s.yaml.Unmarshal(patched, &patchedCfg); err != nil {
		return nil, cfg, fmt.Errorf("%w: post-patch re-unmarshal failed: %v", driving.ErrConfigSchemaInvalid, err)
	}
	if err := revalidateConfigDomain(patchedCfg, it.Path); err != nil {
		return nil, cfg, err
	}
	return patched, patchedCfg, nil
}

// marshalIfChanged is the list-path tail of [applySetItem]: marshal-
// rewrite (comments are lost, like the single-path list route) only
// when the list actually changed.
func (s *ConfigService) marshalIfChanged(
	body []byte, cfg, next ubootYAMLConfig, changed bool, err error,
) ([]byte, ubootYAMLConfig, error) {
	if err != nil {
		return nil, cfg, err
	}
	if !changed {
		return body, cfg, nil
	}
	out, err := s.yaml.Marshal(next)
	if err != nil {
		return nil, cfg, fmt.Errorf("%w: marshal u-boot.yaml: %v", driving.ErrConfigSchemaInvalid, err)
	}
	return out, next, nil
}

// planFeatureSourcesAllow merges the item value (comma-separated) and
// the `--allow-external-feature-sources` entries into the existing
// `devcontainer.featureSources.allow` list (validated, de-duplicated).
func planFeatureSourcesAllow(cfg ubootYAMLConfig, it driving.ConfigSetItem, extra []string) (ubootYAMLConfig, bool, error) {
	incoming, err := parseFeatureSourcesArgument(it.Value)
	if err != nil {
		return cfg, false, fmt.Errorf("%w: %s: %v", driving.ErrConfigValueInvalid, it.Path, err)
	}
	var existing []string
	if cfg.Devcontainer != nil && cfg.Devcontainer.FeatureSources != nil {
		existing = cfg.Devcontainer.FeatureSources.Allow
	}
	merged, err := normaliseFeatureSources(append(append(append([]string{}, existing...), incoming...), extra...))
	if err != nil {
		return cfg, false, fmt.Errorf("%w: %s: %v", driving.ErrConfigValueInvalid, it.Path, err)
	}
	if strings.Join(existing, ",") == strings.Join(merged, ",") {
		return cfg, false, nil
	}
	if cfg.Devcontainer == nil {
		cfg.Devcontainer = &ubootYAMLDevcontainer{}
	}
	if cfg.Devcontainer.FeatureSources == nil {
		cfg.Devcontainer.FeatureSources = &ubootYAMLFeatureSources{}
	}
	cfg.Devcontainer.FeatureSources.Allow = merged
	return cfg, true, nil
}

// planEgressAllow merges the item value (comma-separated host names)
// into the existing `devcontainer.sandbox.egress.allow` list.
func planEgressAllow(cfg ubootYAMLConfig, it driving.ConfigSetItem) (ubootYAMLConfig, bool, error) {
	incoming, err := parseFeatureSourcesArgument(it.Value)
	if err != nil {
		return cfg, false, fmt.Errorf("%w: %s: %v", driving.ErrConfigValueInvalid, it.Path, err)
	}
	var existing []string
	if cfg.Devcontainer != nil && cfg.Devcontainer.Sandbox != nil && cfg.Devcontainer.Sandbox.Egress != nil {
		existing = cfg.Devcontainer.Sandbox.Egress.Allow
	}
	merged, err := normaliseEgressHosts(append(append([]string{}, existing...), incoming...))
	if err != nil {
		return cfg, false, fmt.Errorf("%w: %s: %w", driving.ErrConfigValueInvalid, it.Path, err)
	}
	if strings.Join(existing, ",") == strings.Join(merged, ",") {
		return cfg, false, nil
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
	return cfg, true, nil
}
