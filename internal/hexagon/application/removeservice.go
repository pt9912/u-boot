package application

import (
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"path/filepath"

	"github.com/pt9912/u-boot/internal/hexagon/application/managedblock"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// RemoveServiceService implements [driving.RemoveServiceUseCase] for
// the `u-boot remove <service>` flow (LH-FA-ADD-007). It is the
// mirror of [AddServiceService]: same dependencies (FileSystem,
// YAMLCodec, Confirmer, Logger), same detect → execute layout,
// opposite direction.
//
// Detect-phase delegates to the shared package-level
// [detectServiceState] function so both add and remove classify
// services against the same LH-FA-ADD-005 state table. The execute
// phase is independent because the per-state actions differ
// (remove sets enabled=false + cuts blocks; add sets enabled=true
// + inserts blocks).
type RemoveServiceService struct {
	fs        driven.FileSystem
	yaml      driven.YAMLCodec
	confirmer driven.Confirmer
	logger    driven.Logger
}

// Static check: RemoveServiceService satisfies the driving port.
var _ driving.RemoveServiceUseCase = (*RemoveServiceService)(nil)

// NewRemoveServiceService constructs the service with the driven
// adapters injected by the wiring layer. fs and yaml are mandatory
// (the constructor trusts the wiring layer to provide non-nil
// instances, matching the existing [NewAddServiceService] /
// [NewConfigService] pattern); confirmer and logger are nil-
// tolerant and fall back to the package-internal no-op
// implementations so callers (tests, scripts that do not care
// about prompts) need not wire a stub.
func NewRemoveServiceService(fs driven.FileSystem, yaml driven.YAMLCodec, confirmer driven.Confirmer, logger driven.Logger) *RemoveServiceService {
	if confirmer == nil {
		confirmer = noopConfirmer{}
	}
	if logger == nil {
		logger = noopLogger{}
	}
	return &RemoveServiceService{fs: fs, yaml: yaml, confirmer: confirmer, logger: logger}
}

// Remove implements [driving.RemoveServiceUseCase.Remove].
//
// State-machine flow (T2):
//
//   - Unregistered                       → [driving.ErrServiceUnregistered]
//   - InconsistentYAML / InconsistentBlock → [driving.ErrServiceInconsistent]
//   - Deactivated                        → idempotent no-op (Changed=nil)
//   - Active                             → execute 3 actions (compose-block-
//     remove, env-block-remove, yaml-patch enabled=false)
//   - EnabledUnset                       → execute the same 3 actions
//     (normalises a service that lacks the explicit enabled key)
//
// `--purge` is T3-scope; T2 ignores the flag and never touches
// volumes.
func (s *RemoveServiceService) Remove(_ context.Context, req driving.RemoveServiceRequest) (driving.RemoveServiceResponse, error) {
	if req.BaseDir == "" {
		return driving.RemoveServiceResponse{}, errors.New("BaseDir is required")
	}

	if !isSupportedService(req.ServiceName) {
		return driving.RemoveServiceResponse{}, fmt.Errorf("%w: %q is not in the built-in catalogue %v",
			driving.ErrServiceUnsupported, req.ServiceName.String(), supportedServices())
	}

	state, err := detectServiceState(s.fs, s.yaml, req.BaseDir, req.ServiceName)
	if err != nil {
		return driving.RemoveServiceResponse{}, err
	}

	switch state {
	case domain.ServiceStateUnregistered:
		return driving.RemoveServiceResponse{}, fmt.Errorf("%w: %q was never added; nothing to remove",
			driving.ErrServiceUnregistered, req.ServiceName.String())

	case domain.ServiceStateInconsistentYAML, domain.ServiceStateInconsistentBlock:
		return driving.RemoveServiceResponse{}, fmt.Errorf("%w: %q has a mismatched u-boot.yaml / compose.yaml state; clean up manually before removing",
			driving.ErrServiceInconsistent, req.ServiceName.String())

	case domain.ServiceStateDeactivated:
		s.logger.Debug("remove: idempotent no-op", "service", req.ServiceName.String())
		return driving.RemoveServiceResponse{
			ServiceName: req.ServiceName,
			PriorState:  state,
			State:       state,
		}, nil

	case domain.ServiceStateActive, domain.ServiceStateEnabledUnset:
		return s.executeRemove(req.BaseDir, req.ServiceName, state)

	default:
		// Defensive: detectServiceState covers the six LH-FA-ADD-005
		// states; an unknown value here is a code bug, not a user
		// error.
		return driving.RemoveServiceResponse{}, fmt.Errorf("remove: unexpected state %s for %q",
			state.String(), req.ServiceName.String())
	}
}

// executeRemove performs the three filesystem mutations for an
// Active / EnabledUnset → Deactivated transition:
//
//  1. Strip the `service.<name>` managed block from compose.yaml.
//  2. Strip the same-named block from .env.example.
//  3. Patch `services.<name>.enabled: false` in u-boot.yaml.
//
// Order is u-boot.yaml-last so a mid-flight failure on compose or
// env leaves u-boot.yaml's enabled flag unchanged — the project
// stays self-consistent for a retry. Files that don't exist or
// don't contain the block are silently skipped (idempotent
// per-file cleanup; the spec says compose-block-remove and env-
// block-remove are best-effort).
//
// Volumes are NOT touched here — T3 will layer the `--purge` flag
// + confirmation gate on top of this method.
func (s *RemoveServiceService) executeRemove(baseDir string, svc domain.ServiceName, priorState domain.ServiceState) (driving.RemoveServiceResponse, error) {
	var changed []string

	composeChanged, err := s.removeBlock(baseDir, "compose.yaml", svc)
	if err != nil {
		return driving.RemoveServiceResponse{}, err
	}
	if composeChanged {
		changed = append(changed, "compose.yaml")
	}

	envChanged, err := s.removeBlock(baseDir, ".env.example", svc)
	if err != nil {
		return driving.RemoveServiceResponse{}, err
	}
	if envChanged {
		changed = append(changed, ".env.example")
	}

	yamlPath := filepath.Join(baseDir, "u-boot.yaml")
	yamlBody, err := s.fs.ReadFile(yamlPath)
	if err != nil {
		return driving.RemoveServiceResponse{}, fmt.Errorf("read u-boot.yaml: %w", err)
	}
	patched, err := s.yaml.PatchScalar(yamlBody,
		[]string{"services", svc.String(), "enabled"}, false)
	if err != nil {
		return driving.RemoveServiceResponse{}, fmt.Errorf("patch u-boot.yaml: %w", err)
	}
	if err := s.fs.WriteFile(yamlPath, patched, defaultFileMode); err != nil {
		return driving.RemoveServiceResponse{}, fmt.Errorf("write u-boot.yaml: %w", err)
	}
	changed = append(changed, "u-boot.yaml")

	s.logger.Debug("remove: state transition",
		"service", svc.String(),
		"prior", priorState.String(),
		"changed", len(changed))

	return driving.RemoveServiceResponse{
		ServiceName: svc,
		PriorState:  priorState,
		State:       domain.ServiceStateDeactivated,
		Changed:     changed,
	}, nil
}

// removeBlock strips the `service.<name>` managed block from a
// host file. Returns true when the host file existed AND contained
// the block (and was therefore rewritten); false when either the
// file was absent or did not contain the block — both are no-op
// outcomes the caller skips silently.
//
// Block-malformed surfaces as [driving.ErrServiceInconsistent] for
// consistency with [detectServiceState]'s pre-classification (a
// malformed block stops the operation rather than auto-repairing).
func (s *RemoveServiceService) removeBlock(baseDir, filename string, svc domain.ServiceName) (bool, error) {
	path := filepath.Join(baseDir, filename)
	exists, err := s.fs.Exists(path)
	if err != nil {
		return false, fmt.Errorf("check %s: %w", filename, err)
	}
	if !exists {
		return false, nil
	}
	body, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("read %s: %w", filename, err)
	}
	marker := managedblock.Marker{
		Style: managedblock.StyleHash,
		Name:  serviceMarkerName(svc),
	}
	out, err := managedblock.Replace(body, marker, nil)
	switch {
	case err == nil:
		// Block was present and has been cut out (empty replacement).
	case errors.Is(err, managedblock.ErrBlockNotFound):
		// File doesn't carry the block — no-op, no rewrite needed.
		return false, nil
	case errors.Is(err, managedblock.ErrBlockMalformed):
		return false, fmt.Errorf("%w: malformed managed block for %q in %s: %v",
			driving.ErrServiceInconsistent, svc.String(), filename, err)
	default:
		return false, fmt.Errorf("scan %s for %q block: %w",
			filename, serviceMarkerName(svc), err)
	}
	if err := s.fs.WriteFile(path, out, defaultFileMode); err != nil {
		return false, fmt.Errorf("write %s: %w", filename, err)
	}
	return true, nil
}
