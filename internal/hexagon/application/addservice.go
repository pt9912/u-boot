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

// serviceMarkerNamePrefix is the prefix every per-service managed
// block uses inside compose.yaml / .env.example. The full marker
// name is `service.<servicename>`; the [domain.ServiceName] regex
// guarantees the suffix is safe for both compose service keys and
// managed-block comment lines.
const serviceMarkerNamePrefix = "service."

// supportedServices returns the catalogue of services
// [AddServiceService] knows how to add. M5 shipped only `postgres`
// (LH-FA-ADD-002); slice-v1-keycloak T2 added `keycloak`
// (LH-FA-ADD-003); slice-v1-otel T2 added `otel` (LH-FA-ADD-004)
// — last of the v0.3.0 add-on catalogue.
//
// Function instead of package var to avoid the gochecknoglobals
// false-positive on immutable list constants (same pattern as
// [projectStructureDirs] in initproject.go).
func supportedServices() []string {
	return []string{"postgres", "keycloak", "otel"}
}

// isSupportedService reports whether name is in [supportedServices].
func isSupportedService(name domain.ServiceName) bool {
	got := name.String()
	for _, s := range supportedServices() {
		if s == got {
			return true
		}
	}
	return false
}

// dependenciesFor returns the static [domain.AddOnDependency]
// declarations the catalogue knows about for svc — the source-of-
// truth side-table for the LH-FA-ADD-006 resolver
// (`resolveAddDependencies`, slice-v1-addons-deps T2).
//
// Today's catalogue (postgres only) declares no dependencies, so
// this function returns nil for every supported name; nil for
// unknown names too (the catalogue check fires earlier with
// [driving.ErrServiceUnsupported]). slice-v1-keycloak will add the
// keycloak case with the persistence:external-postgres → postgres
// declaration; slice-v1-otel may add OTel-specific entries.
//
// Function instead of package var to avoid the gochecknoglobals
// false-positive on immutable list constants (same pattern as
// [supportedServices] above).
//
//nolint:unparam // result is always nil today (postgres-only catalogue);
// slice-v1-keycloak adds the first non-nil row (Keycloak → Postgres),
// at which point the lint exemption can be removed.
func dependenciesFor(svc domain.ServiceName) []domain.AddOnDependency {
	switch svc.String() {
	case "postgres":
		return nil
	}
	return nil
}

// serviceMarkerName returns the canonical managed-block name for svc.
// Lives in addservice.go (not in [managedblock]) because managedblock
// must not import domain (architecture-layer hygiene).
func serviceMarkerName(svc domain.ServiceName) string {
	return serviceMarkerNamePrefix + svc.String()
}

// serviceAction classifies what the add use-case should do once it
// has detected the [domain.ServiceState] of the target service. T3
// plans Register / Reactivate / RebuildBlock for the four mutating
// states; [actionRepairArtifacts] is reserved for the T4 active-but-
// incomplete-artefacts path and is not yet planned in T3.
type serviceAction int

const (
	// actionRegister handles Unregistered → Active: a fresh service
	// entry in u-boot.yaml plus a new managed compose-block plus a
	// new .env.example block.
	actionRegister serviceAction = iota

	// actionReactivate handles Deactivated / EnabledUnset → Active:
	// flip `enabled` to true (or set it for the first time) and
	// re-emit the compose-block deterministically.
	actionReactivate

	// actionRebuildBlock handles InconsistentBlock → Active: the
	// YAML anchor is correct but the compose-block is missing;
	// re-emit only the compose-block.
	actionRebuildBlock

	// actionRepairArtifacts handles Active-but-incomplete: the core
	// LH-FA-ADD-005 state machine reports Active, but a service-
	// specific artefact (volume managed block, .env.example block)
	// is missing or stale. T3 reserves the enum value; T4 wires the
	// per-service repair logic.
	actionRepairArtifacts
)

// String returns the canonical lowercase name of the action. Used by
// the logger and by error messages so users see "register" /
// "reactivate" instead of an integer.
func (a serviceAction) String() string {
	switch a {
	case actionRegister:
		return "register"
	case actionReactivate:
		return "reactivate"
	case actionRebuildBlock:
		return "rebuild-block"
	case actionRepairArtifacts:
		return "repair-artifacts"
	default:
		return "unknown"
	}
}

// servicePlan is the planned mutation for a single `u-boot add`
// invocation. T3 fills Service / PriorState / Action; T4c adds
// RepairFlags so the executeAdd plan-phase knows which artefact
// the classifier flagged when state==Active.
type servicePlan struct {
	Service     domain.ServiceName
	PriorState  domain.ServiceState
	Action      serviceAction
	RepairFlags activeArtifactsStatus
}

// activeArtifactsStatus is the pure-classifier output of
// [detectActiveArtifacts]. The three orthogonal flags say which
// LH-FA-ADD-002 artefacts need repair when the LH-FA-ADD-005 core
// state is Active. Abort-class conditions (malformed marker, wrong
// anchor, user-manual entry) are not modelled here — the classifier
// surfaces them as ErrServiceInconsistent before any flag is set.
//
// The MVP volume template is an empty mapping (`{}`) and User-
// customisation of the volume body (driver:, driver_opts:, etc.) is
// explicitly allowed; there is therefore no `VolumeStale` flag. If a
// later add-on needs a content-check on the volume body, it joins as
// its own slice with the corresponding flag added back here.
type activeArtifactsStatus struct {
	// ServiceStale is true when the service.postgres managed
	// compose-block exists but its content lacks a required field
	// (image, environment.POSTGRES_USER/PASSWORD/DB, volumes
	// reference to postgres-data, ports, healthcheck).
	ServiceStale bool
	// VolumeMissing is true when the volume.postgres managed
	// compose-block does not exist.
	VolumeMissing bool
	// EnvMissingOrStale is true when .env.example is absent, has no
	// service.postgres block, or the block lacks a POSTGRES_USER /
	// POSTGRES_PASSWORD / POSTGRES_DB line.
	EnvMissingOrStale bool
}

// needsRepair reports whether any artefact flag is set.
func (s activeArtifactsStatus) needsRepair() bool {
	return s.ServiceStale || s.VolumeMissing || s.EnvMissingOrStale
}

// AddServiceService implements [driving.AddServiceUseCase]. It
// orchestrates the FileSystem, YAMLCodec and Logger driven ports to
// realize the LH-FA-ADD-001 / -002 / -005 add-service flow.
//
// Reachable from the CLI subcommand wired up in M5-T6
// (`u-boot add <service>`); see [cli.newAddCommand]. The full plan-
// and-execute pipeline (template rendering, anchor-checks, file
// writes, Active-Repair) lives in addservice_execute.go and
// addservice_detect.go.
type AddServiceService struct {
	fs        driven.FileSystem
	fsFactory func(driving.AddPreviewMode) (driven.FileSystem, driven.RecorderPort)
	yaml      driven.YAMLCodec
	confirmer driven.Confirmer
	logger    driven.Logger
	// addMu serialises Add() invocations on a single service
	// instance. The PreviewMode-aware fs-swap (line 320-322) mutates
	// the shared s.fs field; concurrent Add calls would race on the
	// swap/restore pair and one caller's writes would route through
	// the other caller's recorder (slice-v1-cli-json-dry-run-add
	// review #10). Holding the mutex around the whole Add body keeps
	// the per-request FS scoping load-bearing.
	addMu sync.Mutex
}

// Static check: AddServiceService satisfies the driving port.
var _ driving.AddServiceUseCase = (*AddServiceService)(nil)

// NewAddServiceService constructs the service with the driven
// adapters injected by the wiring layer. logger and confirmer
// accept nil — both are routed to the package-local noopLogger /
// noopConfirmer so tests / scripts that do not wire them need no
// stub. fs and yaml are mandatory — the service does not invent
// fallbacks for missing infrastructure ports.
//
// confirmer is consulted only on the LH-FA-ADD-006 interactive path
// (missing dependencies, neither --with-deps nor --yes nor
// --no-interactive passed). All other Add flows skip the prompt.
func NewAddServiceService(fs driven.FileSystem, yaml driven.YAMLCodec, confirmer driven.Confirmer, logger driven.Logger) *AddServiceService {
	if logger == nil {
		logger = noopLogger{}
	}
	if confirmer == nil {
		confirmer = noopConfirmer{}
	}
	return &AddServiceService{fs: fs, yaml: yaml, confirmer: confirmer, logger: logger}
}

// NewAddServiceServiceWithFactory is the slice-v1-cli-json-dry-run-add
// T0-(e) Option 4 constructor: instead of a fixed [driven.FileSystem],
// the service receives a factory that picks the FS per
// [driving.AddPreviewMode]. The Composition-Root in `cmd/uboot/main.go`
// wires PreviewNone → production FS, PreviewDryRun/PreviewAndApply →
// recording FS (with the matching passthrough switch).
//
// [Add] reads the mode from [driving.AddServiceRequest.PreviewMode],
// asks the factory for a fresh (fs, recorder) tuple, swaps the
// service-level [fs] field for the request's duration (defer-restored
// so legacy code paths still find the bootstrap FS), and — if the
// recorder is non-nil — maps its [driven.RecorderPort.Captured]
// output into [driving.AddServiceResponse.PlannedFiles] on the way
// out.
//
// Legacy callers (`NewAddServiceService(fs, ...)`) keep working
// unchanged: the factory stays nil and [Add] falls back to the
// stored [fs] field, ignoring PreviewMode (the recorder is nil so
// PlannedFiles stays empty as well).
func NewAddServiceServiceWithFactory(
	fsFactory func(driving.AddPreviewMode) (driven.FileSystem, driven.RecorderPort),
	yaml driven.YAMLCodec,
	confirmer driven.Confirmer,
	logger driven.Logger,
) *AddServiceService {
	if logger == nil {
		logger = noopLogger{}
	}
	if confirmer == nil {
		confirmer = noopConfirmer{}
	}
	// Bootstrap fs from the PreviewNone branch so methods that read
	// s.fs outside of an Add() call (today: none, but slice-followup
	// scope might add some) see a valid adapter even before any
	// request lands.
	bootstrapFS, _ := fsFactory(driving.PreviewNone)
	return &AddServiceService{
		fs:        bootstrapFS,
		fsFactory: fsFactory,
		yaml:      yaml,
		confirmer: confirmer,
		logger:    logger,
	}
}

// selectFS picks the per-request FS pair: if the factory is wired
// (Composition-Root path) it returns the mode-specific tuple;
// otherwise the legacy [fs] field with a nil recorder (PreviewMode
// is ignored — the use case keeps writing to production).
func (s *AddServiceService) selectFS(mode driving.AddPreviewMode) (driven.FileSystem, driven.RecorderPort) {
	if s.fsFactory == nil {
		return s.fs, nil
	}
	return s.fsFactory(mode)
}

// mapCaptureToPlannedFiles converts the driven recorder's mutation
// log into the driving-port wire-shape consumed by the CLI adapter.
// Returns nil for an empty capture so the JSON envelope renders
// `plannedFiles` omitted (slice T0-(d) — Voll-Schema needs the field
// present, but the CLI adapter explicitly populates it; an empty
// capture maps to nil here and the adapter chooses its own form).
//
// Hunks stay nil — the CLI-side diff renderer (T2) fills them from
// NewContent/OldContent.
func mapCaptureToPlannedFiles(records []driven.FileMutationRecord) []driving.PlannedFile {
	if len(records) == 0 {
		return nil
	}
	out := make([]driving.PlannedFile, len(records))
	for i, r := range records {
		out[i] = driving.PlannedFile{
			Path:       r.Path,
			Action:     r.Action,
			NewContent: r.NewContent,
			OldContent: r.OldContent,
		}
	}
	return out
}

// Add implements [driving.AddServiceUseCase.Add]. The dispatch order
// is documented in slice-m5-add-postgres.md §T3:
//
//  1. validate BaseDir is non-empty (non-sentinel error).
//  2. catalogue-check the service name — projects don't influence
//     which services exist, so this check runs before any disk I/O.
//  3. detect the LH-FA-ADD-005 state from u-boot.yaml + compose.yaml.
//  4. branch on the state:
//     - Active           → nil-error no-op (core-state idempotent).
//     - InconsistentYAML → ErrServiceInconsistent + repair hint.
//     - all other states → planAdd → executeAdd (T3: stub).
//
// ctx is threaded to executeAdd so the T4 implementation can honour
// cancellation across its (multi-file) write phase without changing
// the call site here.
func (s *AddServiceService) Add(ctx context.Context, req driving.AddServiceRequest) (driving.AddServiceResponse, error) {
	if req.BaseDir == "" {
		return driving.AddServiceResponse{}, errors.New("BaseDir is required")
	}

	// Serialise the per-request fs-swap (slice-v1-cli-json-dry-run-add
	// review #10): concurrent Add() calls on the same instance would
	// otherwise race on the mutable s.fs field. handleMissingDependencies'
	// recursive dep-install bypasses Add() (calls runAdd directly) so
	// the mutex doesn't self-deadlock; the inner runAdd reuses the
	// outer call's fs/recorder swap, which is exactly the right
	// semantics for dep installs sharing the parent's PreviewMode.
	s.addMu.Lock()
	defer s.addMu.Unlock()

	// slice-v1-cli-json-dry-run-add T1-C: route every FS access of
	// this Add() invocation through the mode-specific adapter. The
	// defer restores the bootstrap [fs] field so tests or scripts
	// that reuse the service instance across requests start each call
	// from the same baseline.
	fs, recorder := s.selectFS(req.PreviewMode)
	prevFS := s.fs
	s.fs = fs
	defer func() { s.fs = prevFS }()

	resp, addErr := s.runAdd(ctx, req)

	// Ensure resp.ServiceName is populated on the error path so the
	// CLI envelope's plannedFiles[] view doesn't ship with an unset
	// ServiceName (review #7): runExecutePlan returns zero-value
	// Response{} on mid-write FS-failure, and the recorder mapping
	// below populates PlannedFiles independently — without this
	// fallback the response carries ServiceName="" while PlannedFiles
	// is non-empty, violating the port contract that
	// `ServiceName echoes the name that was processed`.
	if addErr != nil && resp.ServiceName.String() == "" {
		resp.ServiceName = req.ServiceName
	}

	// The recorder is non-nil only on preview/dry-run paths. Map its
	// log even on the error path — T0-(b) Mid-Write-Failure scenario
	// shows the user the captured calls up to the failure point.
	if recorder != nil {
		resp.PlannedFiles = mapCaptureToPlannedFiles(recorder.Captured())
	}
	return resp, addErr
}

// runAdd is the original Add() body unchanged from M5+: the dispatch
// chain (BaseDir guard → catalogue check → state detection →
// dependency check → state-machine branching). T1-C pulled it out so
// [Add] can wrap it with the PreviewMode-aware FS swap and the
// recorder-capture mapping without duplicating the dispatch logic.
func (s *AddServiceService) runAdd(ctx context.Context, req driving.AddServiceRequest) (driving.AddServiceResponse, error) {
	if !isSupportedService(req.ServiceName) {
		return driving.AddServiceResponse{}, fmt.Errorf("%w: %q is not in the built-in catalogue %v",
			driving.ErrServiceUnsupported, req.ServiceName.String(), supportedServices())
	}

	state, err := detectServiceState(s.fs, s.yaml, req.BaseDir, req.ServiceName)
	if err != nil {
		return driving.AddServiceResponse{}, err
	}

	// slice-v1-addons-deps T2/T3: dependency-check between detection
	// and state-machine dispatch. Postgres has no declared deps
	// today, so the `len(deps) > 0` guard short-circuits the cfg-
	// load for the MVP catalogue. Future Keycloak/OTel slices
	// populate `dependenciesFor` and reach this branch, where the
	// four-mode dispatch (--with-deps / --yes / --no-interactive /
	// interactive prompt) decides whether to auto-install,
	// fail-fast, or recursively Add each missing service.
	if deps := dependenciesFor(req.ServiceName); len(deps) > 0 {
		if depErr := s.checkAddDependencies(ctx, req, deps); depErr != nil {
			return driving.AddServiceResponse{}, depErr
		}
	}

	s.logger.Debug("detected service state",
		"baseDir", req.BaseDir,
		"service", req.ServiceName.String(),
		"state", state.String(),
	)

	switch state {
	case domain.ServiceStateActive:
		// LH-FA-ADD-005 core state is Active, but LH-FA-ADD-002
		// requires the full postgres artefact set. detectActiveArtifacts
		// runs the per-artefact content check and either returns flags
		// for repair or surfaces ErrServiceInconsistent for malformed
		// / wrong-anchor / user-manual-entry states.
		status, err := s.detectActiveArtifacts(req.BaseDir, req.ServiceName)
		if err != nil {
			return driving.AddServiceResponse{}, err
		}
		if !status.needsRepair() {
			return driving.AddServiceResponse{
				ServiceName: req.ServiceName,
				PriorState:  domain.ServiceStateActive,
				State:       domain.ServiceStateActive,
				Changed:     nil,
			}, nil
		}
		plan := servicePlan{
			Service:     req.ServiceName,
			PriorState:  domain.ServiceStateActive,
			Action:      actionRepairArtifacts,
			RepairFlags: status,
		}
		return s.executeAdd(ctx, req.BaseDir, plan)
	case domain.ServiceStateInconsistentYAML:
		return driving.AddServiceResponse{}, fmt.Errorf(
			"%w: managed compose-block for service %q has no matching u-boot.yaml anchor; "+
				"remove the block manually or restore the anchor",
			driving.ErrServiceInconsistent, req.ServiceName.String())
	case domain.ServiceStateUnregistered,
		domain.ServiceStateDeactivated,
		domain.ServiceStateEnabledUnset,
		domain.ServiceStateInconsistentBlock:
		plan, err := s.planAdd(req.ServiceName, state)
		if err != nil {
			return driving.AddServiceResponse{}, err
		}
		return s.executeAdd(ctx, req.BaseDir, plan)
	default:
		// Defensive: detectServiceState only returns the six values
		// above. A new ServiceState added without a switch case
		// surfaces here instead of silently no-oping.
		return driving.AddServiceResponse{}, fmt.Errorf(
			"unhandled service state %s for %q",
			state.String(), req.ServiceName.String())
	}
}

// checkAddDependencies is the slice-v1-addons-deps T3 orchestrator
// for LH-FA-ADD-006. It loads u-boot.yaml, resolves which declared
// deps are missing, and — if any — hands off to
// [handleMissingDependencies] for the four-mode dispatch.
//
// Postgres has no declared deps so this path is only reached once
// a future add-on (Keycloak, OTel) populates [dependenciesFor].
// Tests drive it directly via the [CheckAddDependenciesForTest]
// export seam with synthetic [domain.AddOnDependency] inputs.
//
// The redundant cfg-load (detectServiceState already parsed
// u-boot.yaml) is the price for keeping detectServiceState's
// signature unchanged and shared with [RemoveServiceService].
func (s *AddServiceService) checkAddDependencies(ctx context.Context, req driving.AddServiceRequest, deps []domain.AddOnDependency) error {
	missing, err := s.findMissingDependencies(req.BaseDir, deps)
	if err != nil {
		return err
	}
	if len(missing) == 0 {
		return nil
	}
	return s.handleMissingDependencies(ctx, req, missing)
}

// findMissingDependencies loads + parses u-boot.yaml and runs the
// pure [resolveAddDependencies] resolver against deps. Returned
// list is in deterministic order of first encounter (see resolver
// godoc).
func (s *AddServiceService) findMissingDependencies(baseDir string, deps []domain.AddOnDependency) ([]domain.ServiceName, error) {
	_, _, cfg, err := s.loadAndParseUBootYAML(baseDir)
	if err != nil {
		return nil, fmt.Errorf("dependency check: %w", err)
	}
	return resolveAddDependencies(cfg, deps), nil
}

// handleMissingDependencies applies the LH-FA-ADD-006 four-mode
// dispatch when at least one required service is missing:
//
//   - --with-deps OR --yes: auto-install. Recursively Add each
//     missing service, propagating the same flags so transitive
//     deps inherit the auto-confirm decision.
//   - --no-interactive (without --with-deps / --yes): fail-fast
//     with [driving.ErrDependenciesRequired] (exit code 10).
//   - default (no flags): prompt via [driven.Confirmer.
//     ConfirmAddDependency]. "yes" promotes to auto-install; "no"
//     or I/O error → [driving.ErrDependenciesRequired].
//
// The recursive Add carries the parent's BaseDir so transitive
// deps land in the same project, and inherits the flag set so a
// single top-level `--with-deps` installs the whole chain
// non-interactively. The recursion terminates because each Add
// either registers a dep (shrinking the missing set on subsequent
// runs) or fails-fast.
func (s *AddServiceService) handleMissingDependencies(ctx context.Context, req driving.AddServiceRequest, missing []domain.ServiceName) error {
	missingStrings := missingServiceNamesAsStrings(missing)
	if !req.WithDeps && !req.Yes {
		if req.NoInteractive {
			return fmt.Errorf("%w: %q requires %v which is/are not registered — add them first or rerun with --with-deps",
				driving.ErrDependenciesRequired, req.ServiceName.String(), missingStrings)
		}
		confirmed, err := s.confirmer.ConfirmAddDependency(ctx, req.ServiceName.String(), missingStrings)
		if err != nil {
			return fmt.Errorf("confirm add dependencies for %q: %w", req.ServiceName.String(), err)
		}
		if !confirmed {
			return fmt.Errorf("%w: %q requires %v which is/are not registered — add them first or rerun with --with-deps",
				driving.ErrDependenciesRequired, req.ServiceName.String(), missingStrings)
		}
	}
	for _, dep := range missing {
		subReq := driving.AddServiceRequest{
			BaseDir:       req.BaseDir,
			ServiceName:   dep,
			WithDeps:      req.WithDeps,
			Yes:           req.Yes,
			NoInteractive: req.NoInteractive,
			// Inherit PreviewMode from the outer request (review #9):
			// without this, `u-boot add keycloak --with-deps --dry-run`
			// would write the dep install (postgres) to the production
			// FS — the outer recorder wouldn't see it and the user's
			// dry-run promise would be silently violated. Setting it on
			// subReq is defensive (runAdd doesn't consume PreviewMode
			// today, only Add() does), but pins the contract for any
			// future helper that reaches for req.PreviewMode.
			PreviewMode: req.PreviewMode,
		}
		// Bypass s.Add() so the outer call's mutex/fs-swap stay active
		// across the dep install — runAdd reuses the parent's swapped
		// s.fs (the recorder if PreviewMode != PreviewNone), so dep
		// writes land in the same recorder log as the parent's writes.
		if _, err := s.runAdd(ctx, subReq); err != nil {
			return fmt.Errorf("install dependency %q: %w", dep.String(), err)
		}
	}
	return nil
}

// detectServiceState classifies the LH-FA-ADD-005 state of svc inside
// the project rooted at baseDir. The classification reads at most two
// files (u-boot.yaml + compose.yaml) and never writes.
//
// Package-level free function (extracted from the original
// [AddServiceService].detectServiceState method in slice-v1-add-
// remove T2) so both add and remove flows share the same state-
// detection without inheriting one service's struct.
//
// Error model: technical I/O / permission failures bubble up as
// non-sentinel wrapped errors — they must not be mistaken for a
// fachlicher "project not initialized" or "service inconsistent"
// state. Only two failure modes carry add-sentinels:
//
//   - missing or unparsable u-boot.yaml ⇒ ErrProjectNotInitialized.
//   - malformed `service.<svc>` managed compose-block (any YAML side)
//     ⇒ ErrServiceInconsistent. This pre-classification abort is the
//     Spec §895 repair-hint path applied to the malformed case;
//     LH-FA-ADD-005 only describes six wohlgeformte states, so
//     malformed lives outside the state machine.
func detectServiceState(fs driven.FileSystem, yaml driven.YAMLCodec, baseDir string, svc domain.ServiceName) (domain.ServiceState, error) {
	yamlPath := filepath.Join(baseDir, "u-boot.yaml")
	yamlExists, err := fs.Exists(yamlPath)
	if err != nil {
		return 0, fmt.Errorf("check u-boot.yaml: %w", err)
	}
	if !yamlExists {
		return 0, fmt.Errorf("%w: %s missing", driving.ErrProjectNotInitialized, yamlPath)
	}

	yamlBody, err := fs.ReadFile(yamlPath)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			// TOCTOU between Exists and ReadFile — surface the same
			// sentinel as the missing-file path so the result stays
			// stable regardless of the race winner. This is
			// deliberately asymmetric to the compose.yaml branch
			// further down (which maps a vanished file to
			// block-absent): u-boot.yaml's role is to *be* the
			// project-initialization marker, while compose.yaml only
			// hosts the per-service managed block.
			return 0, fmt.Errorf("%w: %s vanished between Exists and ReadFile",
				driving.ErrProjectNotInitialized, yamlPath)
		}
		return 0, fmt.Errorf("read u-boot.yaml: %w", err)
	}

	var cfg ubootYAMLConfig
	if err := yaml.Unmarshal(yamlBody, &cfg); err != nil {
		return 0, fmt.Errorf("%w: parse u-boot.yaml: %v", driving.ErrProjectNotInitialized, err)
	}

	entry, entryFound := cfg.Services[svc.String()]

	composePath := filepath.Join(baseDir, "compose.yaml")
	composeExists, err := fs.Exists(composePath)
	if err != nil {
		return 0, fmt.Errorf("check compose.yaml: %w", err)
	}

	blockPresent := false
	if composeExists {
		composeBody, err := fs.ReadFile(composePath)
		switch {
		case err == nil:
			marker := managedblock.Marker{
				Style: managedblock.StyleHash,
				Name:  serviceMarkerName(svc),
			}
			_, _, findErr := managedblock.Find(composeBody, marker)
			switch {
			case findErr == nil:
				blockPresent = true
			case errors.Is(findErr, managedblock.ErrBlockNotFound):
				blockPresent = false
			case errors.Is(findErr, managedblock.ErrBlockMalformed):
				return 0, fmt.Errorf("%w: malformed managed compose-block for service %q: %v",
					driving.ErrServiceInconsistent, svc.String(), findErr)
			default:
				return 0, fmt.Errorf("scan compose.yaml for %q block: %w",
					serviceMarkerName(svc), findErr)
			}
		case errors.Is(err, iofs.ErrNotExist):
			// TOCTOU: compose.yaml disappeared between Exists and
			// ReadFile. Treat as block-absent — the file isn't
			// there to host a block. This is deliberately
			// asymmetric to the u-boot.yaml branch above (which
			// maps a vanished file to ErrProjectNotInitialized),
			// because compose.yaml is not the project-existence
			// marker — only the per-service block container.
			blockPresent = false
		default:
			return 0, fmt.Errorf("read compose.yaml: %w", err)
		}
	}

	return classifyServiceState(entryFound, entry, blockPresent), nil
}

// classifyServiceState applies the LH-FA-ADD-005 combination table to
// the three orthogonal observations: YAML-entry presence, the
// Enabled-pointer's three-valued state (nil / *false / *true), and
// the compose-block presence.
//
// Pure function (no I/O) so the table is unit-testable in isolation
// and the I/O orchestration in [detectServiceState] stays linear.
func classifyServiceState(entryFound bool, entry ubootYAMLService, blockPresent bool) domain.ServiceState {
	switch {
	case !entryFound && !blockPresent:
		return domain.ServiceStateUnregistered
	case !entryFound && blockPresent:
		return domain.ServiceStateInconsistentYAML
	case entry.Enabled == nil:
		return domain.ServiceStateEnabledUnset
	case !*entry.Enabled:
		return domain.ServiceStateDeactivated
	case blockPresent:
		return domain.ServiceStateActive
	default:
		return domain.ServiceStateInconsistentBlock
	}
}

// planAdd derives the [servicePlan] for a mutating state. T3 maps the
// four mutating LH-FA-ADD-005 states to three actions
// (EnabledUnset is treated identically to Deactivated per Spec §893,
// so both map to actionReactivate). Active / InconsistentYAML are
// handled in [Add] before planAdd is reached. Stays a method (with
// unused receiver) so T4 can grow it into a port-touching planner
// without re-threading every call site.
func (*AddServiceService) planAdd(svc domain.ServiceName, state domain.ServiceState) (servicePlan, error) {
	plan := servicePlan{Service: svc, PriorState: state}
	switch state {
	case domain.ServiceStateUnregistered:
		plan.Action = actionRegister
	case domain.ServiceStateDeactivated, domain.ServiceStateEnabledUnset:
		plan.Action = actionReactivate
	case domain.ServiceStateInconsistentBlock:
		plan.Action = actionRebuildBlock
	default:
		return servicePlan{}, fmt.Errorf("planAdd: non-mutating state %s for %q",
			state.String(), svc.String())
	}
	return plan, nil
}

// executeAdd is implemented in addservice_execute.go (M5-T4c).
// detectActiveArtifacts lives in addservice_detect.go (M5-T4c).
