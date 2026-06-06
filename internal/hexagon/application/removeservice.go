package application

import (
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"path/filepath"
	"sync"

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
	fsFactory func(driving.PreviewMode) (driven.FileSystem, driven.RecorderPort)
	yaml      driven.YAMLCodec
	confirmer driven.Confirmer
	logger    driven.Logger
	// removeMu serialises Remove() invocations on a single service
	// instance. The PreviewMode-aware s.fs/s.confirmer-swaps in Remove()
	// mutate shared service fields; concurrent Remove calls would race
	// on the swap/restore pairs and one caller's writes would route
	// through the other caller's recorder (slice-v1-cli-json-dry-run-
	// remove T0-(c) R5-F1; inherited from init T0-(d) / add review
	// #10).
	removeMu sync.Mutex
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
//
// Legacy constructor: the resulting service has no fsFactory wired,
// so PreviewMode is ignored at Remove() time (PlannedFiles stays
// nil, recorder is nil). Production wiring uses
// [NewRemoveServiceServiceWithFactory] instead — this constructor
// remains for test sites that don't exercise preview paths.
func NewRemoveServiceService(fs driven.FileSystem, yaml driven.YAMLCodec, confirmer driven.Confirmer, logger driven.Logger) *RemoveServiceService {
	if confirmer == nil {
		confirmer = noopConfirmer{}
	}
	if logger == nil {
		logger = noopLogger{}
	}
	return &RemoveServiceService{fs: fs, yaml: yaml, confirmer: confirmer, logger: logger}
}

// NewRemoveServiceServiceWithFactory is the slice-v1-cli-json-dry-
// run-remove T3 Composition-Root constructor (analog to
// [NewInitProjectServiceWithFactory] / [NewGenerateServiceWithFactory]):
// instead of a fixed [driven.FileSystem], the service receives a
// factory that picks the FS per [driving.PreviewMode]. Composition-
// Root wires PreviewNone → production FS, PreviewDryRun/
// PreviewAndApply → recording FS.
//
// [Remove] reads the mode from [driving.RemoveServiceRequest.PreviewMode],
// asks the factory for a fresh (fs, recorder) tuple, swaps the
// service-level [fs] field for the request's duration (defer-restored)
// and — if the recorder is non-nil — maps its
// [driven.RecorderPort.Captured] output into
// [driving.RemoveServiceResponse.PlannedFiles] on the way out.
func NewRemoveServiceServiceWithFactory(
	fsFactory func(driving.PreviewMode) (driven.FileSystem, driven.RecorderPort),
	yaml driven.YAMLCodec,
	confirmer driven.Confirmer,
	logger driven.Logger,
) *RemoveServiceService {
	if confirmer == nil {
		confirmer = noopConfirmer{}
	}
	if logger == nil {
		logger = noopLogger{}
	}
	// Bootstrap fs from PreviewNone so methods that read s.fs outside
	// of a Remove() call see a valid adapter (mirrors
	// AddServiceServiceWithFactory / InitProjectServiceWithFactory).
	bootstrapFS, _ := fsFactory(driving.PreviewNone)
	return &RemoveServiceService{
		fs:        bootstrapFS,
		fsFactory: fsFactory,
		yaml:      yaml,
		confirmer: confirmer,
		logger:    logger,
	}
}

// selectFS picks the per-request FS pair (slice-v1-cli-json-dry-run-
// remove T3, inherited from init T0-(d) / generate T3): if the
// factory is wired (Composition-Root path) it returns the
// mode-specific tuple; otherwise the legacy [fs] field with a nil
// recorder (PreviewMode is ignored — the use case keeps writing to
// production).
func (s *RemoveServiceService) selectFS(mode driving.PreviewMode) (driven.FileSystem, driven.RecorderPort) {
	if s.fsFactory == nil {
		return s.fs, nil
	}
	return s.fsFactory(mode)
}

// Remove implements [driving.RemoveServiceUseCase.Remove]. T3 wraps
// the original dispatch in a PreviewMode-aware FS-swap +
// Confirmer-Swap + Mutex guard (slice-v1-cli-json-dry-run-remove
// T0-(c) Control-Flow-Skeleton):
//
//  1. Lock removeMu for the whole request — serialises concurrent
//     calls so the per-request s.fs/s.confirmer swaps stay atomic
//     (R5-F1 race-safety; ALL swaps INSIDE the lock region).
//  2. If req.SilenceConfirmer: swap s.confirmer to noopConfirmer{}
//     with defer-restore (R12-F3 mechanism analog init's
//     s.progress-swap initproject.go:345-349).
//  3. Ask the factory for a mode-specific (fs, recorder) tuple
//     (call-scoped local recorder, R6-F2 — NOT a service field).
//     Swap s.fs for the request's duration (defer-restored).
//  4. Call [runRemove] (the original dispatch body) — detectService
//     State runs INSIDE the swap region so the recorder sees the
//     read captures.
//  5. Drain recorder captures BEFORE unswaps and map to
//     resp.PlannedFiles — also on the error path (R4-Recorder-
//     Realität: failing WriteFile/RemoveAll attempt-content lands
//     in PlannedFiles too, see Mid-Write-Failure-AK).
//
// State-machine flow (delegated to runRemove):
//
//   - Unregistered                       → [driving.ErrServiceUnregistered]
//   - InconsistentYAML                   → [driving.ErrServiceInconsistent]
//   - Deactivated                        → idempotent no-op (Changed=nil)
//   - Active / EnabledUnset / InconsistentBlock → state transition
//
// `--purge` (T0-(h)): the LH-FA-CLI-005A §254 confirmation gate
// fires only when the call WILL transition state (Active /
// EnabledUnset / InconsistentBlock) AND PreviewMode != PreviewDryRun
// (T0-(h)(a) skip-logic: Dry-Run implies null-mutations, no gate).
// Volume removal itself remains deferred — `VolumesPurged` stays
// false even on a passed gate.
func (s *RemoveServiceService) Remove(ctx context.Context, req driving.RemoveServiceRequest) (driving.RemoveServiceResponse, error) {
	// Serialise per-request fs/confirmer-swap (slice-v1-cli-json-
	// dry-run-remove T0-(c) R5-F1, inherited from init T0-(d) /
	// add review #10).
	s.removeMu.Lock()
	defer s.removeMu.Unlock()

	// Confirmer-Swap-Mechanismus (T0-(j) NEW, R12-F3): service-field
	// mutation with defer-restore. Conditional on req.SilenceConfirmer
	// analog init's s.progress-swap (initproject.go:345-349) — R7-F4.
	if req.SilenceConfirmer {
		prevConfirmer := s.confirmer
		s.confirmer = noopConfirmer{}
		defer func() { s.confirmer = prevConfirmer }()
	}

	// PreviewMode-aware FS-swap: route every FS access of this
	// Remove() invocation through the mode-specific adapter
	// (production FS for PreviewNone, RecordingFileSystem for
	// PreviewDryRun/PreviewAndApply). recorder is a CALL-SCOPED
	// local variable (R6-F2 — NOT a service field), so parallel
	// Goroutines can't leak captures via the service.
	fs, recorder := s.selectFS(req.PreviewMode)
	prevFS := s.fs
	s.fs = fs
	defer func() { s.fs = prevFS }()

	resp, removeErr := s.runRemove(ctx, req)

	// Drain recorder captures BEFORE the unswaps (LIFO defer
	// resolves fs/confirmer-unswap after this return). Map to
	// resp.PlannedFiles also on the error path (T0-(i) Mid-Write-
	// Failure-AK: user sees captured calls up to the failure point
	// inclusive of the failing underlying attempt; recordingfs.go:
	// 139 records before delegating).
	if recorder != nil {
		resp.PlannedFiles = mapCaptureToPlannedFiles(recorder.Captured(), req.BaseDir)
	}
	return resp, removeErr
}

// runRemove is the original Remove() body unchanged from the pre-T3
// flow: validation → catalogue check → detect → dispatch on state.
// T3 split it out so [Remove] can wrap it with the PreviewMode-aware
// FS/Confirmer-swap + mutex + recorder-capture mapping without
// duplicating the dispatch logic (analog to init's runInit /
// generate's runGenerate).
func (s *RemoveServiceService) runRemove(ctx context.Context, req driving.RemoveServiceRequest) (driving.RemoveServiceResponse, error) {
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
		// WARN emission BEFORE the gate (T0-(c) R8-F3): WARN is
		// visible even in the ErrConfirmationRequired path (user
		// learns: your --purge would have been deferred anyway).
		// T5 unterdrückt WARN bei Error-Diagnostic (R4-F3-Variante A).
		warnings := s.volumesPurgedWarnings(req)

		// Gate fires here — a state transition (or convergence) IS
		// happening. Spec LH-FA-CLI-005A §254 confirmation lives
		// adjacent to the actual destructive intent.
		//
		// T0-(h)(a) Skip-Logic: in PreviewDryRun the gate is skipped
		// (null-mutations implies no destructive op happens; analog
		// init's initGit-skip T0-(n)).
		if req.PreviewMode != driving.PreviewDryRun {
			if err := s.runPurgeGate(ctx, req); err != nil {
				return driving.RemoveServiceResponse{Warnings: warnings}, err
			}
		}
		resp, execErr := s.executeRemove(req.BaseDir, req.ServiceName, state)
		resp.Warnings = warnings
		return resp, execErr

	default:
		// Defensive: detectServiceState's six wohlgeformte LH-FA-ADD-
		// 005 states are all listed above. An unknown value here is a
		// code bug, not a user error.
		return driving.RemoveServiceResponse{}, fmt.Errorf("remove: unexpected state %s for %q",
			state.String(), req.ServiceName.String())
	}
}

// volumesPurgedWarnings constructs the soft-warning Diagnostics for
// the --purge && !VolumesPurged path (slice-v1-cli-json-dry-run-
// remove T0-(g) R5-F2 + R3-F1 Volume-Presence-Check). Returns a
// single-element slice when:
//
//   - req.Purge is true, AND
//   - the catalogue entry declares a named volume
//     (volumeOptional == false — postgres today; keycloak/otel are
//     volumeOptional=true and produce no WARN, R3-F1).
//
// VolumesPurged is implicitly false in v0.3.0 (deferred auto-removal,
// Out-of-Scope). The WARN-Diagnostic carries code LH-FA-ADD-007 +
// level "warn"; T5 maps it into the JSON envelope's `diagnostics[]`
// array. T5 also unterdrückt WARN on the error path (R4-F3 variante A).
//
// Returns nil for `--purge=false` and for volumeless catalogue
// entries — no semantic-falsche WARN for services without volumes.
func (*RemoveServiceService) volumesPurgedWarnings(req driving.RemoveServiceRequest) []driving.WarningEntry {
	if !req.Purge {
		return nil
	}
	entry, ok := catalogueFor(req.ServiceName)
	if !ok || entry.volumeOptional {
		return nil
	}
	return []driving.WarningEntry{
		{
			Code:    "LH-FA-ADD-007",
			Level:   "warn",
			Message: fmt.Sprintf("--purge requested for service %q but volume removal is deferred (v0.3.0); the named volumes are still on disk and untouched", req.ServiceName.String()),
		},
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
		// Multi-`%w` (T3 R2-F1): two sentinels — ErrConfirmerUnavailable
		// for the I/O-class classification (LH-FA-CLI-005A / Exit 10)
		// and the raw err for errors.Is-matches against Driven-Errors.
		// Switch-Order in mapRemoveErrorToDiagnostic checks
		// ErrConfirmerUnavailable AFTER ErrRemoveFileSystem and BEFORE
		// the fachlich service sentinels (T0-(e) R3-F3 + R5-F2).
		return fmt.Errorf("remove service: confirmer: %w: %w", driving.ErrConfirmerUnavailable, err)
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

	extraDeletes, err := s.planExtraFileDeletes(baseDir, svc)
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
			return driving.RemoveServiceResponse{}, fmt.Errorf("remove: write %s: %w: %w", f.path, driving.ErrRemoveFileSystem, err)
		}
		changed = append(changed, f.relPath)
	}
	for _, f := range extraDeletes {
		if err := s.fs.RemoveAll(f.path); err != nil {
			return driving.RemoveServiceResponse{}, fmt.Errorf("remove: remove-all %s: %w: %w", f.relPath, driving.ErrRemoveFileSystem, err)
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
		return plannedRemoveFile{}, fmt.Errorf("remove: exists %s: %w: %w", filename, driving.ErrRemoveFileSystem, err)
	}
	if !exists {
		return plannedRemoveFile{skip: true}, nil
	}
	body, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			return plannedRemoveFile{skip: true}, nil
		}
		return plannedRemoveFile{}, fmt.Errorf("remove: read %s: %w: %w", filename, driving.ErrRemoveFileSystem, err)
	}
	mode, err := s.fileMode(path, defaultFileMode)
	if err != nil {
		return plannedRemoveFile{}, fmt.Errorf("remove: stat %s: %w: %w", filename, driving.ErrRemoveFileSystem, err)
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
		// R4-HIGH-F1 + R5-MED-F3 Klassifikations-Fix: Scanner-default-
		// branch ist KEIN FS-Wrap (kein s.fs.*-Aufruf); semantisch
		// gleicher Marker/Datenkonsistenz-Defekt wie Z. 304
		// (managedblock-malformed). Wrap mit ErrServiceInconsistent
		// statt ErrRemoveFileSystem — Datenkonsistenz-Klasse (Exit 10
		// User-Action: YAML reparieren), nicht FS-Klasse
		// (Exit 14 retry-safe).
		return plannedRemoveFile{}, fmt.Errorf("%w: scan %s for %q block: %v",
			driving.ErrServiceInconsistent, filename, serviceMarkerName(svc), err)
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
		return plannedRemoveFile{}, fmt.Errorf("remove: read u-boot.yaml: %w: %w", driving.ErrRemoveFileSystem, err)
	}
	mode, err := s.fileMode(path, defaultFileMode)
	if err != nil {
		return plannedRemoveFile{}, fmt.Errorf("remove: stat u-boot.yaml: %w: %w", driving.ErrRemoveFileSystem, err)
	}
	patched, err := s.yaml.PatchScalar(body,
		[]string{"services", svc.String(), "enabled"}, false)
	if err != nil {
		// R4-HIGH-F1 + R5-MED-F3 Klassifikations-Fix: yaml.PatchScalar-
		// Failure ist KEIN FS-I/O (s.yaml.* call, kein s.fs.*).
		// Datenkonsistenz-Klasse (invalides YAML-Schema, Exit 10
		// User-Action: YAML reparieren) — Konsolidierung auf
		// ErrServiceInconsistent statt eigenem ErrYAMLPatchFailed-
		// Sentinel.
		return plannedRemoveFile{}, fmt.Errorf("%w: patch u-boot.yaml: %v",
			driving.ErrServiceInconsistent, err)
	}
	return plannedRemoveFile{
		path:    path,
		relPath: "u-boot.yaml",
		content: patched,
		mode:    mode,
	}, nil
}

// planExtraFileDeletes is the planning step for slice-v1-otel T2:
// catalogue entries with extraFiles get their whole-file artefacts
// deleted by [executeRemove] after the standard managed-block
// removals. Returns only the entries whose files exist on disk —
// non-existent extraFiles are skipped silently (idempotency
// guarantee: removing an already-removed service must not error,
// and a missing extraFile is the same shape as a missing managed
// block in compose.yaml).
func (s *RemoveServiceService) planExtraFileDeletes(baseDir string, svc domain.ServiceName) ([]plannedRemoveFile, error) {
	entry, ok := catalogueFor(svc)
	if !ok || len(entry.extraFiles) == 0 {
		return nil, nil
	}
	var out []plannedRemoveFile
	for _, xf := range entry.extraFiles {
		path := filepath.Join(baseDir, xf.Path)
		exists, err := s.fs.Exists(path)
		if err != nil {
			return nil, fmt.Errorf("remove: exists %s: %w: %w", xf.Path, driving.ErrRemoveFileSystem, err)
		}
		if !exists {
			continue
		}
		out = append(out, plannedRemoveFile{path: path, relPath: xf.Path})
	}
	return out, nil
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
