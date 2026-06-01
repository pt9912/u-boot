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

// serviceMarkerNamePrefix is the prefix every per-service managed
// block uses inside compose.yaml / .env.example. The full marker
// name is `service.<servicename>`; the [domain.ServiceName] regex
// guarantees the suffix is safe for both compose service keys and
// managed-block comment lines.
const serviceMarkerNamePrefix = "service."

// supportedServices returns the MVP catalogue of services
// [AddServiceService] knows how to add (LH-FA-ADD-002). M5 ships
// only `postgres`; LH-FA-ADD-003 / LH-FA-ADD-004 (keycloak, otel)
// are V1.
//
// Function instead of package var to avoid the gochecknoglobals
// false-positive on immutable list constants (same pattern as
// [projectStructureDirs] in initproject.go).
func supportedServices() []string {
	return []string{"postgres"}
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
	fs     driven.FileSystem
	yaml   driven.YAMLCodec
	logger driven.Logger
}

// Static check: AddServiceService satisfies the driving port.
var _ driving.AddServiceUseCase = (*AddServiceService)(nil)

// NewAddServiceService constructs the service with the driven
// adapters injected by the wiring layer. logger accepts nil and is
// routed to the package-local [noopLogger] so tests / scripts that
// do not wire a logger need no stub. fs and yaml are mandatory —
// the service does not invent fallbacks for missing infrastructure
// ports.
func NewAddServiceService(fs driven.FileSystem, yaml driven.YAMLCodec, logger driven.Logger) *AddServiceService {
	if logger == nil {
		logger = noopLogger{}
	}
	return &AddServiceService{fs: fs, yaml: yaml, logger: logger}
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

	if !isSupportedService(req.ServiceName) {
		return driving.AddServiceResponse{}, fmt.Errorf("%w: %q is not in the built-in catalogue %v",
			driving.ErrServiceUnsupported, req.ServiceName.String(), supportedServices())
	}

	state, err := detectServiceState(s.fs, s.yaml, req.BaseDir, req.ServiceName)
	if err != nil {
		return driving.AddServiceResponse{}, err
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
