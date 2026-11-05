package application

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// The two list paths of the Homebrew block (LH-FA-DEV-013):
// `devcontainer.features.brew.packages` and `devcontainer.brewTaps.allow`.
// Both hang the new entries on the existing list and drop duplicates, like
// the other list paths; the file is marshalled anew (comments are lost). The
// form of every entry and the approval of the taps are asked on this write
// path as on the read path of `generate devcontainer` (one statement of each
// rule, see [validateBrewEntries], [unlistedBrewTaps]).

// brewListPlan applies one list item to the in-memory document and reports
// whether the document changed; it is the common shape of
// [planBrewTapsAllow] and [planBrewPackages].
type brewListPlan func(cfg ubootYAMLConfig, it driving.ConfigSetItem) (ubootYAMLConfig, bool, error)

// brewListPlanFor returns the plan of a list kind and false for any other
// kind. [ConfigService.setNonScalarKind] (config set) and
// [ConfigService.applySetItem] (config set with several pairs) both route by
// it.
func brewListPlanFor(kind domain.ConfigPathKind) (brewListPlan, bool) {
	switch kind {
	case domain.ConfigDevcontainerBrewTapsAllow:
		return planBrewTapsAllow, true
	case domain.ConfigDevcontainerFeaturePackages:
		return planBrewPackages, true
	}
	return nil, false
}

// withWritableDevcontainer returns cfg with a copy of its `devcontainer`
// sub-tree whose feature map can change without touching the document the
// caller still reads (the old value of the response).
func withWritableDevcontainer(cfg ubootYAMLConfig) ubootYAMLConfig {
	var dc ubootYAMLDevcontainer
	if cfg.Devcontainer != nil {
		dc = *cfg.Devcontainer
	}
	features := make(map[string]ubootYAMLDevcontainerFeature, len(dc.Features)+1)
	for k, v := range dc.Features {
		features[k] = v
	}
	dc.Features = features
	cfg.Devcontainer = &dc
	return cfg
}

// mergeBrewList appends the incoming entries to the existing list, each
// incoming entry checked by valid first, and drops duplicates in first-seen
// order (byte-equal). A failing entry is a value error of the call
// ([driving.ErrConfigValueInvalid]); existing entries passed the check when
// the file was read.
func mergeBrewList(path domain.ConfigPath, existing, incoming []string, valid func(string) error) ([]string, error) {
	for _, raw := range incoming {
		if err := valid(raw); err != nil {
			return nil, fmt.Errorf("%w: %s: %w", driving.ErrConfigValueInvalid, path, err)
		}
	}
	seen := make(map[string]struct{}, len(existing)+len(incoming))
	out := make([]string, 0, len(existing)+len(incoming))
	for _, raw := range append(append([]string{}, existing...), incoming...) {
		if _, dup := seen[raw]; dup {
			continue
		}
		seen[raw] = struct{}{}
		out = append(out, raw)
	}
	return out, nil
}

// planBrewTapsAllow merges the item value (comma-separated `<user>/<tap>`
// entries) into `devcontainer.brewTaps.allow`. Allowing a tap cannot break the
// approval of a formula list, so no further rule runs.
func planBrewTapsAllow(cfg ubootYAMLConfig, it driving.ConfigSetItem) (ubootYAMLConfig, bool, error) {
	incoming, err := parseFeatureSourcesArgument(it.Value)
	if err != nil {
		return cfg, false, fmt.Errorf("%w: %s: %v", driving.ErrConfigValueInvalid, it.Path, err)
	}
	var existing []string
	if cfg.Devcontainer != nil && cfg.Devcontainer.BrewTaps != nil {
		existing = cfg.Devcontainer.BrewTaps.Allow
	}
	merged, err := mergeBrewList(it.Path, existing, incoming, func(raw string) error {
		_, err := domain.NewBrewTap(raw)
		return err
	})
	if err != nil {
		return cfg, false, err
	}
	if strings.Join(existing, ",") == strings.Join(merged, ",") {
		return cfg, false, nil
	}
	next := withWritableDevcontainer(cfg)
	next.Devcontainer.BrewTaps = &ubootYAMLBrewTaps{Allow: merged}
	return next, true, nil
}

// planBrewPackages merges the item value (comma-separated entries
// `<formula>` or `<user>/<tap>/<formula>`) into the formula list of the key
// `brew`. The placement rule and the form of the entries are those of the read
// path; after the merge the approval of the taps runs on the result, so a
// formula of a tap that is not in `devcontainer.brewTaps.allow` ends the call
// with exit code 10 and leaves the file unchanged (with several pairs in one
// call, the pair of the tap list stands before the pair of the formulas). The
// approval asks only an active block: while `brew` is not enabled the list is
// checked for its form, and the enabling asks the approval ([brewViolationMessage]).
func planBrewPackages(cfg ubootYAMLConfig, it driving.ConfigSetItem) (ubootYAMLConfig, bool, error) {
	feature := it.Path.Feature.String()
	var entry ubootYAMLDevcontainerFeature
	if cfg.Devcontainer != nil {
		entry = cfg.Devcontainer.Features[feature]
	}
	if err := brewPackagesPlacementError(feature, entry.Source); err != nil {
		return cfg, false, fmt.Errorf("%w: %w", driving.ErrConfigValueInvalid, err)
	}
	incoming, err := parseFeatureSourcesArgument(it.Value)
	if err != nil {
		return cfg, false, fmt.Errorf("%w: %s: %v", driving.ErrConfigValueInvalid, it.Path, err)
	}
	merged, err := mergeBrewList(it.Path, entry.Packages, incoming, func(raw string) error {
		_, err := domain.NewBrewPackage(raw)
		return err
	})
	if err != nil {
		return cfg, false, err
	}
	if strings.Join(entry.Packages, ",") == strings.Join(merged, ",") {
		return cfg, false, nil
	}
	next := withWritableDevcontainer(cfg)
	entry.Packages = merged
	next.Devcontainer.Features[feature] = entry
	if violations := unlistedBrewTaps(next.Devcontainer); len(violations) > 0 {
		return cfg, false, fmt.Errorf("%w: %s: %s", driving.ErrConfigValueInvalid, it.Path, brewTapViolationMessage(violations))
	}
	return next, true, nil
}

// brewListValue is the comma-joined view of one of the two list kinds, the
// form the value of `config get` and the old and new value of `config set`
// take; "" means unset.
func brewListValue(cfg ubootYAMLConfig, path domain.ConfigPath) string {
	if cfg.Devcontainer == nil {
		return ""
	}
	switch path.Kind {
	case domain.ConfigDevcontainerBrewTapsAllow:
		if cfg.Devcontainer.BrewTaps != nil {
			return strings.Join(cfg.Devcontainer.BrewTaps.Allow, ",")
		}
	case domain.ConfigDevcontainerFeaturePackages:
		return strings.Join(cfg.Devcontainer.Features[path.Feature.String()].Packages, ",")
	}
	return ""
}

// setBrewList is the Set route of the two list kinds: plan, marshal-rewrite,
// write through the preview-selected file system. A call that adds nothing new
// is a NoOp and writes nothing.
func (s *ConfigService) setBrewList(
	req driving.ConfigSetRequest, cfg ubootYAMLConfig, plan brewListPlan,
) (driving.ConfigSetResponse, error) {
	logger := s.logger
	if req.SilenceLogger {
		logger = noopLogger{}
	}
	oldValue := brewListValue(cfg, req.Path)
	next, changed, err := plan(cfg, driving.ConfigSetItem{Path: req.Path, Value: req.Value})
	if err != nil {
		return driving.ConfigSetResponse{}, err
	}
	if !changed {
		logger.Debug("config set: no-op", "path", req.Path.String(), "value", oldValue)
		return driving.ConfigSetResponse{Path: req.Path, OldValue: oldValue, NewValue: oldValue}, nil
	}
	rewritten, err := s.yaml.Marshal(next)
	if err != nil {
		return driving.ConfigSetResponse{}, fmt.Errorf("%w: marshal u-boot.yaml: %v", driving.ErrConfigSchemaInvalid, err)
	}
	fs, recorder := s.selectFS(req.PreviewMode)
	path := filepath.Join(req.BaseDir, domain.ConfigFileName)
	if err := fs.WriteFile(path, rewritten, defaultFileMode); err != nil {
		return driving.ConfigSetResponse{}, fmt.Errorf("%w: write %q: %w", driving.ErrConfigFileSystem, path, err)
	}
	newValue := brewListValue(next, req.Path)
	logWrite(logger, req.PreviewMode, "config set: updated", "config set: would update",
		"path", req.Path.String(), "old", oldValue, "new", newValue)
	resp := driving.ConfigSetResponse{Path: req.Path, OldValue: oldValue, NewValue: newValue}
	if recorder != nil {
		resp.PlannedFiles = mapCaptureToPlannedFiles(recorder.Captured(), req.BaseDir)
	}
	return resp, nil
}
