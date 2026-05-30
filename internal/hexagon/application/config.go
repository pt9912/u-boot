package application

import (
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"path/filepath"
	"strconv"

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

// errStubConfigHandler is the package-internal pin that proves
// `u-boot config <subcommand>` is reachable but the per-method
// handler has not been wired yet. Unexported; exported via
// [ErrStubConfigHandlerForTest] in export_test.go so external
// `_test` packages can pin without leaking the sentinel into the
// driving-port API surface. T3 removes the Get/Show stubs, T4
// removes the Set stub and this sentinel along with it.
var errStubConfigHandler = errors.New("config: handler not implemented")

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

// Set implements [driving.ConfigUseCase.Set]. T2 ships only the
// project-state gate + the stub handler; T4 fills the body with
// the two-stage schema validation per slice-m8-config.md §D3.
func (s *ConfigService) Set(_ context.Context, req driving.ConfigSetRequest) (driving.ConfigSetResponse, error) {
	if req.BaseDir == "" {
		return driving.ConfigSetResponse{}, errors.New("BaseDir is required")
	}
	if err := s.checkProjectInitialized(req.BaseDir); err != nil {
		return driving.ConfigSetResponse{}, err
	}
	return driving.ConfigSetResponse{}, fmt.Errorf("config set %s: %w", req.Path, errStubConfigHandler)
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
	}
	// Unreachable: domain.NewConfigPath only constructs the three
	// kinds above. Defensive branch surfaces the int so a future
	// enum addition without dispatch-switch update is loud.
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
