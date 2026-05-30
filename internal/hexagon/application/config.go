package application

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

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

// Get implements [driving.ConfigUseCase.Get]. T2 ships only the
// project-state gate + the stub handler; T3 fills the body.
func (s *ConfigService) Get(_ context.Context, req driving.ConfigGetRequest) (driving.ConfigGetResponse, error) {
	if req.BaseDir == "" {
		return driving.ConfigGetResponse{}, errors.New("BaseDir is required")
	}
	if err := s.checkProjectInitialized(req.BaseDir); err != nil {
		return driving.ConfigGetResponse{}, err
	}
	return driving.ConfigGetResponse{}, fmt.Errorf("config get %s: %w", req.Path, errStubConfigHandler)
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

// Show implements [driving.ConfigUseCase.Show]. T2 ships only
// the project-state gate + the stub handler; T3 fills the body
// (read u-boot.yaml byte-identically into the response).
func (s *ConfigService) Show(_ context.Context, req driving.ConfigShowRequest) (driving.ConfigShowResponse, error) {
	if req.BaseDir == "" {
		return driving.ConfigShowResponse{}, errors.New("BaseDir is required")
	}
	if err := s.checkProjectInitialized(req.BaseDir); err != nil {
		return driving.ConfigShowResponse{}, err
	}
	return driving.ConfigShowResponse{}, fmt.Errorf("config show: %w", errStubConfigHandler)
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
