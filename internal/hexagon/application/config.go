package application

import (
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// ConfigService implements [driving.ConfigUseCase] for the M8
// `u-boot config get/set/show` flow (LH-FA-CONF-001..005). T2
// ships the skeleton: the three method dispatchers + the shared
// `<BaseDir>/u-boot.yaml`-exists gate; T3 fills Get + Show, T4
// fills Set with the two-stage schema validation
// (slice-m8-config.md §D3).
type ConfigService struct {
	fs     driven.FileSystem
	yaml   driven.YAMLCodec
	logger driven.Logger
}

// Static check: ConfigService satisfies the driving port.
var _ driving.ConfigUseCase = (*ConfigService)(nil)

// NewConfigService constructs the service with the driven
// adapters injected by the wiring layer. logger accepts nil and
// is routed to the package-local [noopLogger]; fs and yaml are
// mandatory.
func NewConfigService(fs driven.FileSystem, yaml driven.YAMLCodec, logger driven.Logger) *ConfigService {
	if logger == nil {
		logger = noopLogger{}
	}
	return &ConfigService{fs: fs, yaml: yaml, logger: logger}
}


// Get implements [driving.ConfigUseCase.Get]. Reads u-boot.yaml,
// unmarshals into [ubootYAMLConfig], and extracts the value at
// req.Path. Per-path semantics (slice-m8-config.md §T3 +
// §D1):
//
//   - ConfigProjectName: project.name is required by
//     LH-FA-CONF-002 §1308; an empty / missing name surfaces as
//     [ErrConfigSchemaInvalid] (corrupt config, not just unset).
//   - ConfigDevcontainerEnabled: the `devcontainer:` block is
//     optional. Missing block OR missing `enabled:` key both
//     surface as [ErrConfigValueNotSet] with a hint pointing
//     at `u-boot init --devcontainer` / `u-boot config set
//     devcontainer.enabled <bool>`.
//   - ConfigServiceEnabled: the service entry is optional.
//     Missing service OR missing `enabled:` key both surface
//     as [ErrConfigValueNotSet] with a hint pointing at
//     `u-boot add <svc>`.
//
// Bool values render as canonical `true` / `false` strings; the
// CLI prints the bare scalar with a trailing newline.
func (s *ConfigService) Get(_ context.Context, req driving.ConfigGetRequest) (driving.ConfigGetResponse, error) {
	if req.BaseDir == "" {
		return driving.ConfigGetResponse{}, errors.New("BaseDir is required")
	}
	cfg, err := s.readUbootYAML(req.BaseDir)
	if err != nil {
		return driving.ConfigGetResponse{}, err
	}
	value, err := extractConfigValue(cfg, req.Path)
	if err != nil {
		return driving.ConfigGetResponse{}, err
	}
	return driving.ConfigGetResponse{Path: req.Path, Value: value}, nil
}

// Set implements [driving.ConfigUseCase.Set]. The pipeline
// follows slice-m8-config.md §D3 strictly: every check below
// runs **before** the disk-mutating WriteFile. If any check
// fails, the file stays byte-identical
// (`writesBefore == writesAfter`):
//
//  1. WriteAllowed gate (M1 review-fix): paths whose
//     [domain.ConfigPath.WriteAllowed] is false (today:
//     `services.<svc>.enabled`) reject the Set with
//     [driving.ErrConfigValueInvalid] and a hint pointing at
//     the canonical write path (`u-boot add <svc>`).
//  2. Value coercion: per-path string→Go-scalar parse
//     (`domain.NewProjectName` for project.name,
//     `strconv.ParseBool` for *.enabled). Failure ⇒
//     [driving.ErrConfigValueInvalid] with the raw value in
//     the message.
//  3. NoOp short-circuit: if the new canonical-stringified
//     value equals the existing one, skip the patch + the
//     WriteFile entirely. The response carries
//     `OldValue == NewValue` so the CLI summary can render
//     a "no change" line. Avoids touching the file when
//     there is nothing to change (idempotency).
//  4. PatchScalar in memory.
//  5. Stage-2 schema roundtrip: re-unmarshal the patched
//     bytes into [ubootYAMLConfig]. yaml.v3 parse failures
//     surface via the V1-yaml-parse sentinel chain and route
//     to [driving.ErrConfigSchemaInvalid].
//  6. Stage-3 domain re-validation: per-path domain
//     validators on the patched config (yaml.v3 Unmarshal is
//     lenient — it accepts a garbage string into a string
//     field without applying domain rules). Failure ⇒
//     [driving.ErrConfigValueInvalid].
//  7. WriteFile (only now).
func (s *ConfigService) Set(_ context.Context, req driving.ConfigSetRequest) (driving.ConfigSetResponse, error) {
	if req.BaseDir == "" {
		return driving.ConfigSetResponse{}, errors.New("BaseDir is required")
	}

	// Stage 0: WriteAllowed gate (M1 review-fix). Runs **before**
	// any FS read so the user sees the kind-specific "use
	// `u-boot add <svc>`" hint even in an uninitialized directory
	// or against a corrupt u-boot.yaml. The pre-fix order
	// (read first, then gate) made the rejection conditional on
	// a successful read, which contradicted the slice plan's
	// claim that Stage 0 is the first stage. M8-T5-Review S1.
	if !req.Path.WriteAllowed {
		return driving.ConfigSetResponse{}, writeRejectedError(req.Path)
	}

	body, cfg, err := s.readUbootYAMLBody(req.BaseDir)
	if err != nil {
		return driving.ConfigSetResponse{}, err
	}

	// LH-FA-DEV-003 list-path branch: `featureSources.allow` is a
	// LIST, not a scalar — PatchScalar cannot represent it. The
	// list-append/dedupe + marshal-rewrite code-route lives in
	// [setFeatureSourcesAllow]. Trade-off: marshal-rewrite loses
	// comments in u-boot.yaml; the Spec (§711-721) does not require
	// comment preservation for this list. The scalar feature paths
	// (.enabled, .source, .version) take the standard PatchScalar
	// path below.
	if req.Path.Kind == domain.ConfigDevcontainerFeatureSourcesAllow {
		return s.setFeatureSourcesAllow(req, cfg)
	}

	// Stage 1: value coercion (catches LH-FA-INIT-006 / bool-
	// parse errors before any patch).
	coerced, formatted, err := coerceConfigValue(req.Path, req.Value)
	if err != nil {
		return driving.ConfigSetResponse{}, err
	}

	// Stage 2: NoOp short-circuit. Compute OldValue from the
	// current cfg; if it equals the canonical form of the new
	// value, skip the write entirely.
	oldValue := extractConfigValueLenient(cfg, req.Path)
	if oldValue == formatted {
		s.logger.Debug("config set: no-op",
			"path", req.Path.String(), "value", formatted)
		return driving.ConfigSetResponse{
			Path: req.Path, OldValue: oldValue, NewValue: formatted,
		}, nil
	}

	// Stage 3: PatchScalar in memory.
	yamlPath := configPathToYAMLPath(req.Path)
	patched, err := s.yaml.PatchScalar(body, yamlPath, coerced)
	if err != nil {
		if errors.Is(err, driven.ErrYAMLParse) {
			return driving.ConfigSetResponse{}, fmt.Errorf(
				"%w: PatchScalar parse failure on %s: %v",
				driving.ErrConfigSchemaInvalid, req.Path, err)
		}
		return driving.ConfigSetResponse{}, fmt.Errorf(
			"%w: PatchScalar(%s): %v",
			driving.ErrConfigSchemaInvalid, req.Path, err)
	}

	// Stage 4: re-unmarshal into ubootYAMLConfig.
	var patchedCfg ubootYAMLConfig
	if err := s.yaml.Unmarshal(patched, &patchedCfg); err != nil {
		return driving.ConfigSetResponse{}, fmt.Errorf(
			"%w: post-patch re-unmarshal failed: %v",
			driving.ErrConfigSchemaInvalid, err)
	}

	// Stage 5: per-path domain re-validation on the patched
	// config (yaml.v3 Unmarshal is lenient).
	if err := revalidateConfigDomain(patchedCfg, req.Path); err != nil {
		return driving.ConfigSetResponse{}, err
	}

	// Stage 6: WriteFile (only now).
	path := filepath.Join(req.BaseDir, "u-boot.yaml")
	if err := s.fs.WriteFile(path, patched, defaultFileMode); err != nil {
		return driving.ConfigSetResponse{}, fmt.Errorf("%w: write %q: %v",
			driving.ErrConfigFileSystem, path, err)
	}

	newValue := extractConfigValueLenient(patchedCfg, req.Path)
	s.logger.Info("config set: updated",
		"path", req.Path.String(), "old", oldValue, "new", newValue)
	return driving.ConfigSetResponse{
		Path: req.Path, OldValue: oldValue, NewValue: newValue,
	}, nil
}

// writeRejectedError builds the ErrConfigValueInvalid response for
// a path whose [domain.ConfigPath.WriteAllowed] is false. The
// message is kind-conditional so a future write-disallowed
// non-service kind would not leak the "register the service"
// phrasing through a wrong-shape interpolation (M8-T5-Review N1).
func writeRejectedError(path domain.ConfigPath) error {
	switch path.Kind {
	case domain.ConfigServiceEnabled:
		return fmt.Errorf(
			"%w: %s is not writable via `u-boot config set` because the LH-FA-ADD-005 state machine owns the lifecycle; use `u-boot add %s` to register the service",
			driving.ErrConfigValueInvalid, path, path.Service.String())
	}
	return fmt.Errorf(
		"%w: %s is not writable via `u-boot config set`",
		driving.ErrConfigValueInvalid, path)
}

// coerceConfigValue parses raw into the path's expected Go
// scalar form and returns:
//   - coerced: the typed value to hand to YAMLCodec.PatchScalar
//     (e.g. string for project.name, bool for *.enabled).
//   - formatted: the canonical stringified form used for the
//     NoOp short-circuit comparison and for response.NewValue.
//   - err: [driving.ErrConfigValueInvalid] on coerce failure.
//
// ConfigDevcontainerFeatureSourcesAllow takes a different route in
// [ConfigService.Set] (list-append marshal-rewrite, not PatchScalar)
// and therefore short-circuits before this function — it is not
// listed in the switch and a defensive panic-safe fallthrough surfaces
// as the unknown-kind error if the caller bypasses Set().
func coerceConfigValue(path domain.ConfigPath, raw string) (coerced any, formatted string, err error) {
	switch path.Kind {
	case domain.ConfigProjectName:
		name, err := domain.NewProjectName(raw)
		if err != nil {
			return nil, "", fmt.Errorf("%w: project.name %q: %v",
				driving.ErrConfigValueInvalid, raw, err)
		}
		return name.String(), name.String(), nil
	case domain.ConfigDevcontainerEnabled,
		domain.ConfigServiceEnabled,
		domain.ConfigDevcontainerFeatureEnabled:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, "", fmt.Errorf("%w: %s expects a bool (true/false/1/0/T/F/…), got %q",
				driving.ErrConfigValueInvalid, path, raw)
		}
		return b, strconv.FormatBool(b), nil
	case domain.ConfigDevcontainerFeatureSource:
		// LH-FA-DEV-003 Stage-1 format check (empty / scheme /
		// host). Allowlist-membership runs in
		// revalidateConfigDomain after the patch is in place.
		if err := validateFeatureSource(raw); err != nil {
			return nil, "", fmt.Errorf("%w: %s: %v",
				driving.ErrConfigValueInvalid, path, err)
		}
		trimmed := strings.TrimSpace(raw)
		return trimmed, trimmed, nil
	case domain.ConfigDevcontainerFeatureVersion:
		// Version pins are free-form strings ("1", "1.2.0", "21"
		// for Java, …); reject empty so a deliberate omission
		// stays expressed as a missing key, not a sentinel-empty
		// override.
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			return nil, "", fmt.Errorf("%w: %s expects a non-empty version pin",
				driving.ErrConfigValueInvalid, path)
		}
		return trimmed, trimmed, nil
	}
	return nil, "", fmt.Errorf("%w: unknown ConfigPathKind %d",
		driving.ErrConfigValueInvalid, int(path.Kind))
}

// configPathToYAMLPath translates a domain.ConfigPath into the
// []string nested-key list YAMLCodec.PatchScalar expects.
// Kept package-local because the YAML traversal convention is a
// codec-protocol detail, not domain semantics.
//
// Only the two write-allowed ConfigPath kinds are handled —
// ConfigServiceEnabled is intentionally absent because Set
// rejects it at the WriteAllowed gate before reaching this
// helper. A future `config force-set`-style relaxation would
// add the case and the corresponding revalidator branch in one
// PR; until then, dead-branch coverage churn is avoided.
func configPathToYAMLPath(path domain.ConfigPath) []string {
	switch path.Kind {
	case domain.ConfigProjectName:
		return []string{"project", "name"}
	case domain.ConfigDevcontainerEnabled:
		return []string{"devcontainer", "enabled"}
	case domain.ConfigDevcontainerFeatureEnabled:
		return []string{"devcontainer", "features", path.Feature.String(), "enabled"}
	case domain.ConfigDevcontainerFeatureSource:
		return []string{"devcontainer", "features", path.Feature.String(), "source"}
	case domain.ConfigDevcontainerFeatureVersion:
		return []string{"devcontainer", "features", path.Feature.String(), "version"}
	}
	// Unreachable for scalar kinds: Set only reaches this helper
	// after the WriteAllowed gate, and the list-path
	// (ConfigDevcontainerFeatureSourcesAllow) bypasses this helper
	// in favour of a marshal-rewrite branch. Returning a non-empty
	// sentinel so PatchScalar would surface a loud
	// ErrYAMLPathInvalid on a programmer-error future enum addition
	// without dispatch-switch update.
	return []string{"__unknown_config_path_kind__"}
}

// revalidateConfigDomain runs the per-path domain validators on
// the patched config (Stage 5 of the Set pipeline). yaml.v3
// Unmarshal accepts strings without applying domain rules; this
// closes that gap. Returns [driving.ErrConfigValueInvalid] on
// failure.
//
// Symmetric to [configPathToYAMLPath]: only the two write-allowed
// kinds are handled because Set rejects ConfigServiceEnabled at
// the WriteAllowed gate before reaching here. A future relaxation
// would extend both helpers in lock-step.
func revalidateConfigDomain(cfg ubootYAMLConfig, path domain.ConfigPath) error {
	switch path.Kind {
	case domain.ConfigProjectName:
		// Stage 1 (coerceConfigValue) already runs NewProjectName
		// on the raw user value, and PatchScalar writes the
		// string verbatim with the !!str tag; this re-validation
		// is therefore redundant for project.name in production
		// today. Kept for symmetry with future kinds where
		// Stage 1 cannot pre-validate (e.g. nested objects where
		// the canonical form is only knowable after the patch
		// re-marshal) — and as defense-in-depth against a
		// YAML-codec bug that would silently corrupt the leaf
		// value (M8-T5-Review N2 informational).
		if _, err := domain.NewProjectName(cfg.Project.Name); err != nil {
			return fmt.Errorf(
				"%w: post-patch project.name failed domain re-validation: %v",
				driving.ErrConfigValueInvalid, err)
		}
		return nil
	case domain.ConfigDevcontainerEnabled:
		// Pointer-check + shape: the value must exist post-patch
		// and have parsed as bool (PatchScalar writes the !!bool
		// tag). The yaml-v3 round-trip already gave us a *bool;
		// nil here means PatchScalar wrote something the unmarshal
		// could not bind, which is structurally invalid.
		if cfg.Devcontainer == nil || cfg.Devcontainer.Enabled == nil {
			return fmt.Errorf(
				"%w: post-patch devcontainer.enabled is absent or unbound",
				driving.ErrConfigValueInvalid)
		}
		return nil
	case domain.ConfigDevcontainerFeatureEnabled,
		domain.ConfigDevcontainerFeatureSource,
		domain.ConfigDevcontainerFeatureVersion:
		return revalidateFeatureEntry(cfg, path)
	}
	return fmt.Errorf("%w: unknown ConfigPathKind %d",
		driving.ErrConfigValueInvalid, int(path.Kind))
}

// revalidateFeatureEntry is the post-patch sanity check for the
// three scalar `devcontainer.features.<name>.*` kinds. Extracted
// from [revalidateConfigDomain] to keep the latter under the
// gocyclo threshold.
//
// For .source it additionally enforces the LH-FA-DEV-003 /
// LH-NFA-SEC-004 allowlist rule: the newly written URL must appear
// verbatim in `devcontainer.featureSources.allow`. `--yes` does not
// suffice — the user must have populated the allowlist first via
// `config set devcontainer.featureSources.allow <url>` or
// `--allow-external-feature-sources <url>`.
func revalidateFeatureEntry(cfg ubootYAMLConfig, path domain.ConfigPath) error {
	entry, ok := lookupFeatureEntry(cfg, path.Feature)
	if !ok {
		return fmt.Errorf("%w: post-patch %s is absent",
			driving.ErrConfigValueInvalid, path)
	}
	switch path.Kind {
	case domain.ConfigDevcontainerFeatureEnabled:
		if entry.Enabled == nil {
			return fmt.Errorf("%w: post-patch %s is unbound",
				driving.ErrConfigValueInvalid, path)
		}
		return nil
	case domain.ConfigDevcontainerFeatureSource:
		if entry.Source == "" {
			return fmt.Errorf("%w: post-patch %s is empty",
				driving.ErrConfigValueInvalid, path)
		}
		if !featureSourceInAllow(cfg, entry.Source) {
			return fmt.Errorf(
				"%w: %s: external source %q is not in devcontainer.featureSources.allow; add it via `u-boot config set devcontainer.featureSources.allow %s` (LH-FA-DEV-003 / LH-NFA-SEC-004 — --yes is not sufficient)",
				driving.ErrConfigValueInvalid, path, entry.Source, entry.Source)
		}
		return nil
	case domain.ConfigDevcontainerFeatureVersion:
		if entry.Version == "" {
			return fmt.Errorf("%w: post-patch %s is empty",
				driving.ErrConfigValueInvalid, path)
		}
		return nil
	}
	return fmt.Errorf("%w: unknown ConfigPathKind %d",
		driving.ErrConfigValueInvalid, int(path.Kind))
}

// lookupFeatureEntry returns the ubootYAMLDevcontainerFeature entry
// for the given name, or zero+false on miss. Used by
// [revalidateConfigDomain] and [extractConfigValue] for the
// devcontainer.features.<feature>.* kinds — shared helper avoids
// repeating the nil-walk on every branch.
func lookupFeatureEntry(cfg ubootYAMLConfig, name domain.FeatureName) (ubootYAMLDevcontainerFeature, bool) {
	if cfg.Devcontainer == nil {
		return ubootYAMLDevcontainerFeature{}, false
	}
	entry, ok := cfg.Devcontainer.Features[name.String()]
	return entry, ok
}

// featureSourceInAllow reports whether src appears verbatim in
// `cfg.Devcontainer.FeatureSources.Allow`. Comparison is
// case-sensitive, byte-equal — the user's source override must
// match the allowlist entry literally. The trimmed-whitespace
// convention from [normaliseFeatureSources] is the canonical form
// of both sides, so leading/trailing whitespace is not a
// confounding factor in practice.
func featureSourceInAllow(cfg ubootYAMLConfig, src string) bool {
	if cfg.Devcontainer == nil || cfg.Devcontainer.FeatureSources == nil {
		return false
	}
	for _, allowed := range cfg.Devcontainer.FeatureSources.Allow {
		if allowed == src {
			return true
		}
	}
	return false
}

// extractConfigValueLenient is the Get-helper without the
// NotSet/SchemaInvalid error mapping — it returns the empty
// string for unset / corrupt values so Set can compute
// OldValue/NewValue without breaking the pipeline on edge cases
// (e.g. setting devcontainer.enabled for the first time, where
// the OldValue is legitimately absent).
func extractConfigValueLenient(cfg ubootYAMLConfig, path domain.ConfigPath) string {
	v, err := extractConfigValue(cfg, path)
	if err != nil {
		return ""
	}
	return v
}

// readUbootYAMLBody is the Set-helper that returns both the raw
// body (for PatchScalar) and the parsed cfg (for OldValue
// extraction). Same sentinel mapping as [ConfigService.readUbootYAML].
func (s *ConfigService) readUbootYAMLBody(baseDir string) ([]byte, ubootYAMLConfig, error) {
	if err := s.checkProjectInitialized(baseDir); err != nil {
		return nil, ubootYAMLConfig{}, err
	}
	path := filepath.Join(baseDir, "u-boot.yaml")
	body, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			return nil, ubootYAMLConfig{}, fmt.Errorf(
				"%w: %s vanished between Exists and ReadFile",
				driving.ErrProjectNotInitialized, path)
		}
		return nil, ubootYAMLConfig{}, fmt.Errorf("%w: read %q: %v",
			driving.ErrConfigFileSystem, path, err)
	}
	var cfg ubootYAMLConfig
	if err := s.yaml.Unmarshal(body, &cfg); err != nil {
		return nil, ubootYAMLConfig{}, fmt.Errorf("%w: parse u-boot.yaml: %v",
			driving.ErrConfigSchemaInvalid, err)
	}
	return body, cfg, nil
}

// Show implements [driving.ConfigUseCase.Show]. Reads
// u-boot.yaml byte-identically into the response (no re-parse,
// no re-marshal). Comments and formatting are preserved
// (slice-m8-config.md §D5).
func (s *ConfigService) Show(_ context.Context, req driving.ConfigShowRequest) (driving.ConfigShowResponse, error) {
	if req.BaseDir == "" {
		return driving.ConfigShowResponse{}, errors.New("BaseDir is required")
	}
	if err := s.checkProjectInitialized(req.BaseDir); err != nil {
		return driving.ConfigShowResponse{}, err
	}
	path := filepath.Join(req.BaseDir, "u-boot.yaml")
	body, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			// TOCTOU between Exists and ReadFile — surface as the
			// same sentinel as the absent path so the CLI message
			// is consistent.
			return driving.ConfigShowResponse{}, fmt.Errorf(
				"%w: %s vanished between Exists and ReadFile",
				driving.ErrProjectNotInitialized, path)
		}
		return driving.ConfigShowResponse{}, fmt.Errorf("%w: read %q: %v",
			driving.ErrConfigFileSystem, path, err)
	}
	return driving.ConfigShowResponse{Body: body}, nil
}

// readUbootYAML reads and parses `<baseDir>/u-boot.yaml`. Shared
// between Get (T3) and Set (T4). Sentinel mapping mirrors the
// slice-m8-config.md §D6 table:
//
//   - missing file ⇒ [driving.ErrProjectNotInitialized] (the
//     gate is shared with the M5/M7 helpers).
//   - read failure ⇒ [driving.ErrConfigFileSystem].
//   - parse failure ⇒ [driving.ErrConfigSchemaInvalid] via the
//     V1-yaml-parse sentinel chain (driven.ErrYAMLParse).
func (s *ConfigService) readUbootYAML(baseDir string) (ubootYAMLConfig, error) {
	if err := s.checkProjectInitialized(baseDir); err != nil {
		return ubootYAMLConfig{}, err
	}
	path := filepath.Join(baseDir, "u-boot.yaml")
	body, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			return ubootYAMLConfig{}, fmt.Errorf(
				"%w: %s vanished between Exists and ReadFile",
				driving.ErrProjectNotInitialized, path)
		}
		return ubootYAMLConfig{}, fmt.Errorf("%w: read %q: %v",
			driving.ErrConfigFileSystem, path, err)
	}
	var cfg ubootYAMLConfig
	if err := s.yaml.Unmarshal(body, &cfg); err != nil {
		return ubootYAMLConfig{}, fmt.Errorf("%w: parse u-boot.yaml: %v",
			driving.ErrConfigSchemaInvalid, err)
	}
	return cfg, nil
}

// extractConfigValue returns the stringified value at path inside
// cfg, or the per-path NotSet / SchemaInvalid sentinel when the
// field is absent. Pure function over cfg + path; no I/O. Used
// by Get (T3) and by Set (T4) to compute OldValue.
func extractConfigValue(cfg ubootYAMLConfig, path domain.ConfigPath) (string, error) {
	switch path.Kind {
	case domain.ConfigProjectName:
		if cfg.Project.Name == "" {
			return "", fmt.Errorf(
				"%w: u-boot.yaml has no `project.name` value; this is a corrupt config (LH-FA-CONF-002 §1308 requires it)",
				driving.ErrConfigSchemaInvalid)
		}
		return cfg.Project.Name, nil

	case domain.ConfigDevcontainerEnabled:
		if cfg.Devcontainer == nil || cfg.Devcontainer.Enabled == nil {
			return "", fmt.Errorf(
				"%w: %s — run `u-boot init --devcontainer` or `u-boot config set devcontainer.enabled <true|false>` to initialize",
				driving.ErrConfigValueNotSet, path)
		}
		return strconv.FormatBool(*cfg.Devcontainer.Enabled), nil

	case domain.ConfigServiceEnabled:
		entry, ok := cfg.Services[path.Service.String()]
		if !ok || entry.Enabled == nil {
			return "", fmt.Errorf(
				"%w: %s — run `u-boot add %s` to register the service",
				driving.ErrConfigValueNotSet, path, path.Service.String())
		}
		return strconv.FormatBool(*entry.Enabled), nil

	case domain.ConfigDevcontainerFeatureSourcesAllow,
		domain.ConfigDevcontainerFeatureEnabled,
		domain.ConfigDevcontainerFeatureSource,
		domain.ConfigDevcontainerFeatureVersion:
		return extractDevcontainerFeatureValue(cfg, path)
	}
	// Unreachable: domain.NewConfigPath only constructs the kinds
	// above. Defensive branch surfaces the int so a future enum
	// addition without dispatch-switch update is loud.
	return "", fmt.Errorf("%w: unknown ConfigPathKind %d",
		driving.ErrConfigPathUnknown, int(path.Kind))
}

// extractDevcontainerFeatureValue handles the four LH-FA-DEV-003
// devcontainer.* kinds. Extracted from [extractConfigValue] to keep
// the latter under the gocognit threshold.
func extractDevcontainerFeatureValue(cfg ubootYAMLConfig, path domain.ConfigPath) (string, error) {
	if path.Kind == domain.ConfigDevcontainerFeatureSourcesAllow {
		if cfg.Devcontainer == nil || cfg.Devcontainer.FeatureSources == nil ||
			len(cfg.Devcontainer.FeatureSources.Allow) == 0 {
			return "", fmt.Errorf(
				"%w: %s — populate via `u-boot config set %s <url>` or `--allow-external-feature-sources <url>`",
				driving.ErrConfigValueNotSet, path, path)
		}
		// Comma-joined view mirrors the input format the user
		// passes on the way in. Get/Set keeps the input/output
		// formats symmetric so a `get | set` round-trip is a
		// NoOp.
		return strings.Join(cfg.Devcontainer.FeatureSources.Allow, ","), nil
	}
	entry, ok := lookupFeatureEntry(cfg, path.Feature)
	if !ok {
		return "", fmt.Errorf("%w: %s — set it via `u-boot config set %s <value>`",
			driving.ErrConfigValueNotSet, path, path)
	}
	switch path.Kind {
	case domain.ConfigDevcontainerFeatureEnabled:
		if entry.Enabled == nil {
			return "", fmt.Errorf("%w: %s — set it via `u-boot config set %s true`",
				driving.ErrConfigValueNotSet, path, path)
		}
		return strconv.FormatBool(*entry.Enabled), nil
	case domain.ConfigDevcontainerFeatureSource:
		if entry.Source == "" {
			return "", fmt.Errorf("%w: %s — set it via `u-boot config set %s <url>` after adding the URL to devcontainer.featureSources.allow",
				driving.ErrConfigValueNotSet, path, path)
		}
		return entry.Source, nil
	case domain.ConfigDevcontainerFeatureVersion:
		if entry.Version == "" {
			return "", fmt.Errorf("%w: %s — set it via `u-boot config set %s <version>`",
				driving.ErrConfigValueNotSet, path, path)
		}
		return entry.Version, nil
	}
	return "", fmt.Errorf("%w: unknown ConfigPathKind %d",
		driving.ErrConfigPathUnknown, int(path.Kind))
}

// checkProjectInitialized mirrors the M5 / M6 / M7 helper of the
// same shape so the three Config methods produce identical
// sentinel-mapping behaviour at the CLI. Shared between all
// three methods (slice-m8-config.md §T2).
func (s *ConfigService) checkProjectInitialized(baseDir string) error {
	path := filepath.Join(baseDir, "u-boot.yaml")
	exists, err := s.fs.Exists(path)
	if err != nil {
		return fmt.Errorf("config service: Exists(%q): %w", path, err)
	}
	if !exists {
		return fmt.Errorf("config service: %q absent: %w", path, driving.ErrProjectNotInitialized)
	}
	return nil
}

// setFeatureSourcesAllow handles the LH-FA-DEV-003 list-path:
// `devcontainer.featureSources.allow`. The value is parsed as a
// comma-separated list of URLs, each entry runs through
// [validateFeatureSource], the resulting set is appended to the
// existing list with silent-dedupe per spec/lastenheft.md:1352, and
// the whole config is marshal-rewritten back to u-boot.yaml.
//
// Unlike the scalar PatchScalar path, this branch loses comment
// preservation for the rewritten file — documented in the Set
// dispatch comment.
//
// Idempotency: NoOp short-circuit when the comma-joined
// before/after lists are byte-equal (i.e., the user added entries
// that were already present).
func (s *ConfigService) setFeatureSourcesAllow(
	req driving.ConfigSetRequest, cfg ubootYAMLConfig,
) (driving.ConfigSetResponse, error) {
	// Stage 1: parse the user-provided inputs (positional value +
	// `--allow-external-feature-sources` flag cumulation per Spec
	// §714-718) and validate them as user input. Failures here
	// map to ErrConfigValueInvalid (Code 10) because the user can
	// fix them by passing a different URL.
	incoming, err := parseFeatureSourcesArgument(req.Value)
	if err != nil {
		return driving.ConfigSetResponse{}, fmt.Errorf("%w: %s: %v",
			driving.ErrConfigValueInvalid, req.Path, err)
	}
	incoming = append(incoming, req.AllowExternalFeatureSources...)
	if _, err := normaliseFeatureSources(incoming); err != nil {
		return driving.ConfigSetResponse{}, fmt.Errorf("%w: %s: %v",
			driving.ErrConfigValueInvalid, req.Path, err)
	}

	// Stage 2: read existing list, merge, validate the merged
	// result. The second normaliseFeatureSources call runs against
	// (existing + incoming) — if it fails here while the incoming-
	// only check passed, the existing list is corrupt
	// (impossible in practice because T1 validates on every read,
	// but defensive). Map that case to ErrConfigSchemaInvalid.
	var existing []string
	if cfg.Devcontainer != nil && cfg.Devcontainer.FeatureSources != nil {
		existing = cfg.Devcontainer.FeatureSources.Allow
	}
	oldValue := strings.Join(existing, ",")
	mergedRaw := make([]string, 0, len(existing)+len(incoming))
	mergedRaw = append(mergedRaw, existing...)
	mergedRaw = append(mergedRaw, incoming...)
	merged, err := normaliseFeatureSources(mergedRaw)
	if err != nil {
		return driving.ConfigSetResponse{}, fmt.Errorf("%w: %s: existing u-boot.yaml allow list contains invalid entry: %v",
			driving.ErrConfigSchemaInvalid, req.Path, err)
	}
	newValue := strings.Join(merged, ",")

	if oldValue == newValue {
		s.logger.Debug("config set: no-op",
			"path", req.Path.String(), "value", newValue)
		return driving.ConfigSetResponse{
			Path: req.Path, OldValue: oldValue, NewValue: newValue,
		}, nil
	}

	// Stage 3: in-memory mutation (no PatchScalar — list path).
	if cfg.Devcontainer == nil {
		cfg.Devcontainer = &ubootYAMLDevcontainer{}
	}
	if cfg.Devcontainer.FeatureSources == nil {
		cfg.Devcontainer.FeatureSources = &ubootYAMLFeatureSources{}
	}
	cfg.Devcontainer.FeatureSources.Allow = merged

	// Stage 4: marshal-rewrite (loses comments; acceptable per
	// Spec §711-721 which has no comment-preservation requirement
	// for this list).
	rewritten, err := s.yaml.Marshal(cfg)
	if err != nil {
		return driving.ConfigSetResponse{}, fmt.Errorf("%w: marshal u-boot.yaml: %v",
			driving.ErrConfigSchemaInvalid, err)
	}

	// Stage 5: WriteFile.
	path := filepath.Join(req.BaseDir, "u-boot.yaml")
	if err := s.fs.WriteFile(path, rewritten, defaultFileMode); err != nil {
		return driving.ConfigSetResponse{}, fmt.Errorf("%w: write %q: %v",
			driving.ErrConfigFileSystem, path, err)
	}

	s.logger.Info("config set: updated",
		"path", req.Path.String(), "old", oldValue, "new", newValue)
	return driving.ConfigSetResponse{
		Path: req.Path, OldValue: oldValue, NewValue: newValue,
	}, nil
}

// parseFeatureSourcesArgument splits a comma-separated user
// argument (positional value of `config set devcontainer.
// featureSources.allow <value>` or the value of
// `--allow-external-feature-sources`) into individual entries with
// per-element whitespace trim and empty-entry rejection. The
// resulting slice is the raw user input — validation + dedupe
// happen at the next pipeline stage.
//
// Spec §718: comma-separation is mandatory. An empty argument
// (no entries after split) returns nil so callers can NoOp-check
// against an empty list.
func parseFeatureSourcesArgument(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			return nil, fmt.Errorf("comma-separated list contains an empty entry near %q", raw)
		}
		out = append(out, trimmed)
	}
	return out, nil
}
