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
// State-machine flow:
//
//   - Unregistered                       → [driving.ErrServiceUnregistered]
//   - InconsistentYAML                   → [driving.ErrServiceInconsistent]
//     (orphan compose block without YAML entry — requires manual
//     cleanup; remove cannot infer the intended YAML state)
//   - Deactivated                        → idempotent no-op (Changed=nil);
//     NO `--purge` confirmation gate fires (no destructive op happens)
//   - Active                             → 3-action transition via
//     two-phase plan-then-write (review-followup F1+F2)
//   - EnabledUnset                       → same 3 actions (normalises
//     a service that lacks the explicit enabled key)
//   - InconsistentBlock                  → forwards-only convergence
//     (review-followup F1): YAML says enabled=true but no compose
//     block. Remove sets enabled=false + idempotent block-removes —
//     `removeBlock` is no-op-on-absent, so a partial-write retry
//     completes the unfinished work. Asymmetric to
//     [AddServiceService] which rejects InconsistentBlock — add
//     cannot auto-converge to "active" without knowing the original
//     state; remove can converge to "disabled" unambiguously.
//
// `--purge` (T3 + review-followup F6): the LH-FA-CLI-005A §254
// confirmation gate fires only when the call WILL transition state
// (Active / EnabledUnset / InconsistentBlock). On Deactivated the
// gate is a no-op for a no-op — we skip it. Volume removal itself
// remains deferred — `VolumesPurged` stays false even on a passed
// gate; the CLI summary (T4) flags the deferred work for the user.
func (s *RemoveServiceService) Remove(ctx context.Context, req driving.RemoveServiceRequest) (driving.RemoveServiceResponse, error) {
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

	// Reject the truly-unrecoverable states BEFORE any gate or write.
	// InconsistentBlock is INTENTIONALLY not in this list (F1):
	// remove can converge it forwards.
	switch state {
	case domain.ServiceStateUnregistered:
		return driving.RemoveServiceResponse{}, fmt.Errorf("%w: %q was never added; nothing to remove",
			driving.ErrServiceUnregistered, req.ServiceName.String())
	case domain.ServiceStateInconsistentYAML:
		return driving.RemoveServiceResponse{}, fmt.Errorf("%w: %q has an orphan compose block without a YAML entry; clean up manually before removing",
			driving.ErrServiceInconsistent, req.ServiceName.String())
	}

	switch state {
	case domain.ServiceStateDeactivated:
		// No-op + no destructive op → no gate (F6).
		s.logger.Debug("remove: idempotent no-op", "service", req.ServiceName.String(), "purge", req.Purge)
		return driving.RemoveServiceResponse{
			ServiceName: req.ServiceName,
			PriorState:  state,
			State:       state,
		}, nil

	case domain.ServiceStateActive,
		domain.ServiceStateEnabledUnset,
		domain.ServiceStateInconsistentBlock:
		// Gate fires here — a state transition (or convergence) IS
		// happening. Spec LH-FA-CLI-005A §254 confirmation lives
		// adjacent to the actual destructive intent.
		if err := s.runPurgeGate(ctx, req); err != nil {
			return driving.RemoveServiceResponse{}, err
		}
		return s.executeRemove(req.BaseDir, req.ServiceName, state)

	default:
		// Defensive: detectServiceState's six wohlgeformte LH-FA-ADD-
		// 005 states are all listed above. An unknown value here is a
		// code bug, not a user error.
		return driving.RemoveServiceResponse{}, fmt.Errorf("remove: unexpected state %s for %q",
			state.String(), req.ServiceName.String())
	}
}

// runPurgeGate implements the LH-FA-CLI-005A §254 confirmation truth
// table for `u-boot remove --purge`. Mirrors
// [DownService.runConfirmationGate]; the two flows share the
// destructive-op confirmation semantics and reuse the same
// [driven.Confirmer.ConfirmRemoveVolumes] adapter method.
//
// Returns nil when the caller should proceed; returns a wrapped
// [driving.ErrConfirmationRequired] (CLI code 10) when the
// destructive op is refused or skipped per the spec table:
//
//	--purge | --yes | --no-interactive | result
//	-------+-------+------------------+----------------------
//	  no    |  any  |       any        | proceed (no gate)
//	  yes   |  yes  |       any        | proceed (auto-yes)
//	  yes   |  no   |       yes        | refuse  (ErrConfirmationRequired)
//	  yes   |  no   |       no         | ask     (Confirmer.ConfirmRemoveVolumes)
func (s *RemoveServiceService) runPurgeGate(ctx context.Context, req driving.RemoveServiceRequest) error {
	if !req.Purge {
		return nil
	}
	if req.Yes {
		return nil
	}
	if req.NoInteractive {
		return fmt.Errorf("remove service: --purge refused in --no-interactive without --yes: %w",
			driving.ErrConfirmationRequired)
	}
	confirmed, err := s.confirmer.ConfirmRemoveVolumes(ctx, req.BaseDir)
	if err != nil {
		return fmt.Errorf("remove service: confirmer error: %w", err)
	}
	if !confirmed {
		return fmt.Errorf("remove service: --purge declined by user: %w",
			driving.ErrConfirmationRequired)
	}
	return nil
}

// plannedRemoveFile is one entry in the in-memory execute plan: the
// project-relative path, the new content, the file mode to preserve,
// and a skip flag for "no change needed" (file absent or block not
// present). The executeRemove apply-loop iterates over the plan
// after every per-file render succeeds.
type plannedRemoveFile struct {
	path    string
	relPath string
	content []byte
	mode    iofs.FileMode
	skip    bool
}

// executeRemove drives the Active / EnabledUnset / InconsistentBlock
// → Deactivated transition via TWO PHASES (review-followup F1):
//
//  1. Plan: read every input file, compute its new bytes, capture
//     its existing mode. No disk writes. A render error
//     (managedblock-malformed, yaml-patch-failure) aborts here
//     before ANY file is touched.
//  2. Apply: write the planned files in order. Per-file write
//     failures still leave partial state on disk, but a retry now
//     converges because [Remove] dispatches InconsistentBlock back
//     into this method instead of rejecting it (F1).
//
// File mode (review-followup F3): each plan entry captures the
// original file's mode via [driven.FileSystem.Lstat], so a chmod-
// hardened compose.yaml (e.g. 0o600) stays at 0o600 after the
// rewrite. Symmetric to addservice's loadForPatch helper.
func (s *RemoveServiceService) executeRemove(baseDir string, svc domain.ServiceName, priorState domain.ServiceState) (driving.RemoveServiceResponse, error) {
	composePlan, err := s.planBlockRemoval(baseDir, "compose.yaml", svc)
	if err != nil {
		return driving.RemoveServiceResponse{}, err
	}
	envPlan, err := s.planBlockRemoval(baseDir, ".env.example", svc)
	if err != nil {
		return driving.RemoveServiceResponse{}, err
	}
	yamlPlan, err := s.planYAMLDisable(baseDir, svc)
	if err != nil {
		return driving.RemoveServiceResponse{}, err
	}

	plan := []plannedRemoveFile{composePlan, envPlan, yamlPlan}
	var changed []string
	for _, f := range plan {
		if f.skip {
			continue
		}
		if err := s.fs.WriteFile(f.path, f.content, f.mode); err != nil {
			return driving.RemoveServiceResponse{}, fmt.Errorf("write %s: %w", f.path, err)
		}
		changed = append(changed, f.relPath)
	}

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

// planBlockRemoval is the planning step for one managed-block host
// file (compose.yaml or .env.example). Reads the file, captures its
// mode, computes the new bytes after `managedblock.Replace(... nil)`
// removes the service block. A file that doesn't exist or doesn't
// carry the block produces skip=true — the apply-loop omits it.
//
// Block-malformed surfaces as [driving.ErrServiceInconsistent]
// (symmetric to [detectServiceState]'s pre-classification: a
// malformed block stops the operation rather than auto-repairing).
func (s *RemoveServiceService) planBlockRemoval(baseDir, filename string, svc domain.ServiceName) (plannedRemoveFile, error) {
	path := filepath.Join(baseDir, filename)
	exists, err := s.fs.Exists(path)
	if err != nil {
		return plannedRemoveFile{}, fmt.Errorf("check %s: %w", filename, err)
	}
	if !exists {
		return plannedRemoveFile{skip: true}, nil
	}
	body, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			return plannedRemoveFile{skip: true}, nil
		}
		return plannedRemoveFile{}, fmt.Errorf("read %s: %w", filename, err)
	}
	mode, err := s.fileMode(path, defaultFileMode)
	if err != nil {
		return plannedRemoveFile{}, fmt.Errorf("stat %s: %w", filename, err)
	}
	marker := managedblock.Marker{
		Style: managedblock.StyleHash,
		Name:  serviceMarkerName(svc),
	}
	out, err := managedblock.Replace(body, marker, nil)
	switch {
	case err == nil:
		return plannedRemoveFile{
			path:    path,
			relPath: filename,
			content: out,
			mode:    mode,
		}, nil
	case errors.Is(err, managedblock.ErrBlockNotFound):
		return plannedRemoveFile{skip: true}, nil
	case errors.Is(err, managedblock.ErrBlockMalformed):
		return plannedRemoveFile{}, fmt.Errorf("%w: malformed managed block for %q in %s: %v",
			driving.ErrServiceInconsistent, svc.String(), filename, err)
	default:
		return plannedRemoveFile{}, fmt.Errorf("scan %s for %q block: %w",
			filename, serviceMarkerName(svc), err)
	}
}

// planYAMLDisable is the planning step for u-boot.yaml. Reads the
// file, captures its mode, and patches the enabled flag to false.
// The file is mandatory at this point — [detectServiceState] only
// classifies states above [domain.ServiceStateUnregistered] when
// u-boot.yaml is present and parseable.
func (s *RemoveServiceService) planYAMLDisable(baseDir string, svc domain.ServiceName) (plannedRemoveFile, error) {
	path := filepath.Join(baseDir, "u-boot.yaml")
	body, err := s.fs.ReadFile(path)
	if err != nil {
		return plannedRemoveFile{}, fmt.Errorf("read u-boot.yaml: %w", err)
	}
	mode, err := s.fileMode(path, defaultFileMode)
	if err != nil {
		return plannedRemoveFile{}, fmt.Errorf("stat u-boot.yaml: %w", err)
	}
	patched, err := s.yaml.PatchScalar(body,
		[]string{"services", svc.String(), "enabled"}, false)
	if err != nil {
		return plannedRemoveFile{}, fmt.Errorf("patch u-boot.yaml: %w", err)
	}
	return plannedRemoveFile{
		path:    path,
		relPath: "u-boot.yaml",
		content: patched,
		mode:    mode,
	}, nil
}

// fileMode captures the existing file's permission bits via Lstat
// so a rewrite preserves them (review-followup F3 — addservice
// already preserves via loadForPatch; remove was asymmetrically
// downgrading to 0o644). Falls back to fallbackMode when Lstat
// fails on a non-existence path (TOCTOU race; the apply-loop
// would error on the write anyway, but a sane mode keeps the
// error sentinel-classifiable).
func (s *RemoveServiceService) fileMode(path string, fallbackMode iofs.FileMode) (iofs.FileMode, error) {
	info, err := s.fs.Lstat(path)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			return fallbackMode, nil
		}
		return 0, err
	}
	return info.Mode().Perm(), nil
}
