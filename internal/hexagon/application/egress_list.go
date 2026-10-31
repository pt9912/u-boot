package application

import (
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"path/filepath"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// EgressListService implements [driving.EgressListUseCase] for
// `u-boot egress list` (LH-FA-DEV-008). It reads `u-boot.yaml` and
// `.git/config` through the [driven.FileSystem] port and uses no other
// port: no Docker, no network, no write.
type EgressListService struct {
	fs     driven.FileSystem
	yaml   driven.YAMLCodec
	logger driven.Logger
}

// Static check: EgressListService satisfies the driving port.
var _ driving.EgressListUseCase = (*EgressListService)(nil)

// NewEgressListService constructs the service; logger accepts nil
// (routed to [noopLogger]).
func NewEgressListService(fs driven.FileSystem, yaml driven.YAMLCodec, logger driven.Logger) *EgressListService {
	if logger == nil {
		logger = noopLogger{}
	}
	return &EgressListService{fs: fs, yaml: yaml, logger: logger}
}

// List computes the allowlist of `generate devcontainer` with the origin
// of each host. The clone source is resolved as for `generate devcontainer
// --sandbox` ([resolveSandboxClone]): the list is the one the sandbox
// profile bakes into the start script, also while the profile or the
// restriction is not active. The "no remote" warning of that resolution is
// not an error here; the response says that no clone host contributes.
func (s *EgressListService) List(_ context.Context, req driving.EgressListRequest) (driving.EgressListResponse, error) {
	if req.BaseDir == "" {
		return driving.EgressListResponse{}, errors.New("BaseDir is required")
	}
	cfg, err := s.readConfig(req.BaseDir)
	if err != nil {
		return driving.EgressListResponse{}, err
	}
	cloneURL, _, err := resolveSandboxClone(s.fs, req.BaseDir, true, cfg.Devcontainer)
	if err != nil {
		return driving.EgressListResponse{}, err
	}
	resp := driving.EgressListResponse{
		SandboxProfile: profileIsSandbox(cfg.Devcontainer),
		EgressEnabled:  egressEnabled(cfg.Devcontainer),
		CloneSource:    cloneURL != "",
		Hosts:          egressHostSources(cfg.Devcontainer, cloneURL),
	}
	s.logger.Debug("egress list", "hosts", len(resp.Hosts), "active", resp.Active())
	return resp, nil
}

// readConfig reads and validates `<baseDir>/u-boot.yaml`: a missing file
// wraps [driving.ErrProjectNotInitialized], a read failure
// [driving.ErrEgressFileSystem], an unparsable or schema-invalid document
// [driving.ErrConfigSchemaInvalid] (as `config list` classifies them).
func (s *EgressListService) readConfig(baseDir string) (ubootYAMLConfig, error) {
	path := filepath.Join(baseDir, domain.ConfigFileName)
	exists, err := s.fs.Exists(path)
	if err != nil {
		return ubootYAMLConfig{}, fmt.Errorf("%w: probe %q: %w", driving.ErrEgressFileSystem, path, err)
	}
	if !exists {
		return ubootYAMLConfig{}, fmt.Errorf("%w: %s not found", driving.ErrProjectNotInitialized, path)
	}
	body, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			return ubootYAMLConfig{}, fmt.Errorf("%w: %s vanished between Exists and ReadFile", driving.ErrProjectNotInitialized, path)
		}
		return ubootYAMLConfig{}, fmt.Errorf("%w: read %q: %w", driving.ErrEgressFileSystem, path, err)
	}
	var cfg ubootYAMLConfig
	if err := s.yaml.Unmarshal(body, &cfg); err != nil {
		return ubootYAMLConfig{}, fmt.Errorf("%w: parse u-boot.yaml: %v", driving.ErrConfigSchemaInvalid, err)
	}
	if err := validateDevcontainer(cfg.Devcontainer); err != nil {
		return ubootYAMLConfig{}, fmt.Errorf("%w: u-boot.yaml devcontainer schema invalid: %v", driving.ErrConfigSchemaInvalid, err)
	}
	return cfg, nil
}
