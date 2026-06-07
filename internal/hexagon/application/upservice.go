package application

import (
	"context"
	"errors"
	"fmt"
	"io"
	iofs "io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// UpService implements [driving.UpUseCase]. It orchestrates
// `u-boot up`: validate the project layout, hand off to the
// [driven.DockerEngine] for the actual `compose up`, then drive a
// polling loop that classifies each declared service against
// LH-FA-UP-001 (`healthy` / `running` + TCP-probe) until every
// service stabilizes or the request's timeout fires.
//
// LH-FA-UP-001 §970 (--timeout=0): on Timeout==0 the service
// short-circuits to fire-and-forget — no `ComposePs` roundtrip, no
// probes, and the result carries a single `up.fire-and-forget`
// [domain.SeverityInfo] diagnostic.
type UpService struct {
	fs     driven.FileSystem
	yaml   driven.YAMLCodec
	engine driven.DockerEngine
	net    driven.NetProbe
	clock  driven.Clock
	logger driven.Logger
}

// Static check: UpService satisfies the driving port.
var _ driving.UpUseCase = (*UpService)(nil)

// NewUpService constructs the service with the driven adapters the
// M6 polling loop needs. A nil logger is routed to the package-level
// noopLogger so tests and dry-runs do not need a stub.
func NewUpService(fs driven.FileSystem, yaml driven.YAMLCodec, engine driven.DockerEngine, net driven.NetProbe, clock driven.Clock, logger driven.Logger) *UpService {
	if logger == nil {
		logger = noopLogger{}
	}
	return &UpService{fs: fs, yaml: yaml, engine: engine, net: net, clock: clock, logger: logger}
}

// pollInterval is the wall-clock gap between `ComposePs` polling
// iterations. Domain-level constant so a future per-project override
// would land here.
const pollInterval = 500 * time.Millisecond

// dialTimeout bounds each individual TCP probe; small relative to
// pollInterval so multiple probes per iteration fit comfortably.
const dialTimeout = 300 * time.Millisecond

// Up implements [driving.UpUseCase.Up].
func (s *UpService) Up(ctx context.Context, req driving.UpRequest) (driving.UpResponse, error) {
	if req.BaseDir == "" {
		return driving.UpResponse{}, errors.New("up service: BaseDir is empty")
	}
	if req.Timeout < 0 {
		return driving.UpResponse{}, fmt.Errorf("up service: Timeout must be >= 0, got %v", req.Timeout)
	}
	if err := s.checkProjectInitialized(req.BaseDir); err != nil {
		return driving.UpResponse{}, err
	}
	compose, err := s.readComposeFile(req.BaseDir)
	if err != nil {
		return driving.UpResponse{}, err
	}

	// ProgressSink-Silencing (slice-v1-cli-json-dry-run-up-down
	// T0-(c) form (d) + T3 wiring): CLI sets req.SilenceProgress =
	// flags.JSON; the use case swaps the effective sink to
	// io.Discard so the Compose stderr stream does not pollute
	// machine-consumable output. nil-Default stays in the adapter
	// (`progressSinkOrDiscard`, R5-MED-1 DRY-Prinzip).
	effective := req.ProgressSink
	if req.SilenceProgress {
		effective = io.Discard
	}
	if _, err := s.engine.ComposeUp(ctx, req.BaseDir, driven.ComposeUpOptions{
		Detach:       true,
		ProgressSink: effective,
	}); err != nil {
		return driving.UpResponse{}, fmt.Errorf("up service: ComposeUp on %q: %w", req.BaseDir, err)
	}

	if req.Timeout == 0 {
		return driving.UpResponse{Result: domain.UpResult{
			Stabilized: false,
			Diagnostics: []domain.Diagnostic{{
				ID:       "up.fire-and-forget",
				Severity: domain.SeverityInfo,
				Message:  "Started with --timeout=0; status check skipped.",
				Hint:     "Run `u-boot doctor` or `docker compose ps` to inspect status.",
			}},
		}}, nil
	}

	return s.pollUntilStabilized(ctx, req.BaseDir, req.Timeout, compose)
}

// checkProjectInitialized verifies that `<BaseDir>/u-boot.yaml`
// exists. Permission/IO errors are wrapped without the project-not-
// initialized sentinel — they are technical, not fachlich.
func (s *UpService) checkProjectInitialized(baseDir string) error {
	path := filepath.Join(baseDir, "u-boot.yaml")
	exists, err := s.fs.Exists(path)
	if err != nil {
		return fmt.Errorf("up service: Exists(%q): %w: %w", path, driving.ErrUpFileSystem, err)
	}
	if !exists {
		return fmt.Errorf("up service: %q absent: %w", path, driving.ErrProjectNotInitialized)
	}
	return nil
}

// composeFileDecode is the local YAML-decoding shape for the
// `compose.yaml` parts the polling loop cares about. Only `Ports`
// and `Healthcheck` are read; everything else (image, environment,
// volumes, …) flows through Compose untouched.
type composeFileDecode struct {
	Services map[string]composeServiceDecode `yaml:"services"`
}

type composeServiceDecode struct {
	Ports       []any                 `yaml:"ports"`
	Healthcheck *composeHealthcheckYAML `yaml:"healthcheck"`
}

type composeHealthcheckYAML struct {
	Disable bool `yaml:"disable"`
}

// readComposeFile loads + parses `<BaseDir>/compose.yaml`. Missing
// file → [driving.ErrComposeFileMissing]; permission/IO errors →
// technical wrap (no fachlich sentinel).
func (s *UpService) readComposeFile(baseDir string) (composeFileDecode, error) {
	var out composeFileDecode
	path := filepath.Join(baseDir, "compose.yaml")
	exists, err := s.fs.Exists(path)
	if err != nil {
		return out, fmt.Errorf("up service: Exists(%q): %w: %w", path, driving.ErrUpFileSystem, err)
	}
	if !exists {
		return out, fmt.Errorf("up service: %q absent: %w", path, driving.ErrComposeFileMissing)
	}
	data, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			return out, fmt.Errorf("up service: %q vanished after Exists: %w", path, driving.ErrComposeFileMissing)
		}
		return out, fmt.Errorf("up service: ReadFile(%q): %w: %w", path, driving.ErrUpFileSystem, err)
	}
	if err := s.yaml.Unmarshal(data, &out); err != nil {
		return out, fmt.Errorf("up service: parse compose.yaml: %w", err)
	}
	return out, nil
}

// servicePollState carries per-service iteration-to-iteration state
// for the polling loop. The map keys are Compose service names.
type servicePollState struct {
	// restartObservations counts consecutive stateRestarting
	// observations. Reset to 0 the moment the service shows a
	// different state.
	restartObservations int
	// probeTargets is parsed once at loop start; reused every
	// iteration. nil means no probable port.
	probeTargets []portProbeTarget
	// healthcheckRequired is `true` iff the compose.yaml service
	// has a healthcheck mapping with `disable: false` (or omitted).
	healthcheckRequired bool
	// unknownStateReported guarantees the
	// `up.state.<service>.unknown` diagnostic is emitted only once
	// per service per Up() call, even if the loop runs many
	// iterations with the same unknown state.
	unknownStateReported bool
	// portUnreachableReported guarantees the
	// `up.port.<service>.unreachable` warn-diagnostic is emitted
	// only once per service per Up() call, even if the healthcheck-
	// dominated probe fails over many iterations after the service
	// already reached `healthy`. LH-FA-UP-001 §968 + slice plan §141.
	portUnreachableReported bool
}

// pollUntilStabilized drives the polling loop until every service
// reaches [domain.OutcomeStabilized], a service Fails, the request
// timeout elapses, or the context is cancelled.
func (s *UpService) pollUntilStabilized(ctx context.Context, baseDir string, timeout time.Duration, compose composeFileDecode) (driving.UpResponse, error) {
	startTime := s.clock.Now()
	pollStates := s.initPollStates(compose)
	diagnostics := s.parseAllPortDiagnostics(compose, pollStates)

	for {
		if err := ctx.Err(); err != nil {
			return driving.UpResponse{}, fmt.Errorf("up service: poll cancelled at t=%v: %w", s.clock.Now().Sub(startTime), err)
		}

		services, err := s.engine.ComposePs(ctx, baseDir)
		if err != nil {
			return driving.UpResponse{}, fmt.Errorf("up service: ComposePs at t=%v: %w", s.clock.Now().Sub(startTime), err)
		}

		stabilized, failedName, failedState := s.classifyAllServices(ctx, services, compose, pollStates, &diagnostics)
		if failedName != "" {
			return driving.UpResponse{}, fmt.Errorf("up service: %q reached terminal state %q: %w", failedName, failedState, driven.ErrComposeRuntime)
		}
		if stabilized {
			return driving.UpResponse{Result: buildResult(services, true, diagnostics)}, nil
		}

		if s.clock.Now().Sub(startTime) >= timeout {
			pending := pendingServiceNames(services, compose, pollStates)
			return driving.UpResponse{}, fmt.Errorf("up service: stabilization timeout after %v, pending: %s: %w", timeout, strings.Join(pending, ", "), driving.ErrStabilizationTimeout)
		}

		s.clock.Sleep(pollInterval)
	}
}

// initPollStates builds an entry for every service declared in
// compose.yaml so iterations can update counters without checking
// existence. Compose may return extra services (manually added);
// those are tolerated and ignored by the classifier.
func (*UpService) initPollStates(compose composeFileDecode) map[string]*servicePollState {
	states := make(map[string]*servicePollState, len(compose.Services))
	for name, def := range compose.Services {
		probeTargets, _ := parseServicePortsForState(def.Ports)
		states[name] = &servicePollState{
			probeTargets:        probeTargets,
			healthcheckRequired: healthcheckRequired(def),
		}
	}
	return states
}

// parseAllPortDiagnostics walks every service's ports[] at loop
// start and emits a Severity-warn diagnostic for each non-probable
// element (one diagnostic per service per port-index — IDs of the
// form `up.port.<service>.<index>`). LH-FA-UP-001 §969.
func (*UpService) parseAllPortDiagnostics(compose composeFileDecode, _ map[string]*servicePollState) []domain.Diagnostic {
	names := sortedServiceNames(compose.Services)
	diagnostics := make([]domain.Diagnostic, 0)
	for _, name := range names {
		def := compose.Services[name]
		for idx, raw := range def.Ports {
			if _, probable := parseComposePort(raw); !probable {
				diagnostics = append(diagnostics, domain.Diagnostic{
					ID:       fmt.Sprintf("up.port.%s.%d", name, idx),
					Severity: domain.SeverityWarn,
					Message:  fmt.Sprintf("Service %q port entry #%d (%v) is not TCP-probable; stabilization will rely on healthcheck or running-only.", name, idx, raw),
					Hint:     "Use a single host:container TCP mapping if a port probe is required.",
				})
			}
		}
	}
	return diagnostics
}

// classifyAllServices classifies each *compose.yaml-declared*
// service against the latest `ComposePs` snapshot. Returns
// stabilized=true when every declared service reaches
// [domain.OutcomeStabilized]; returns failedName != "" on the
// first service that lands in [domain.OutcomeFailed] (terminal,
// abort the loop immediately).
func (s *UpService) classifyAllServices(ctx context.Context, ps []driven.ComposeService, compose composeFileDecode, states map[string]*servicePollState, diagnostics *[]domain.Diagnostic) (stabilized bool, failedName, failedState string) {
	psByName := make(map[string]driven.ComposeService, len(ps))
	for _, svc := range ps {
		psByName[svc.Name] = svc
	}

	allStabilized := true
	// Iterate in deterministic order so a multi-service failure
	// in the same iteration always reports the alphabetically
	// first name — keeps CLI error strings stable across runs and
	// makes grep-based scripts reproducible.
	for _, name := range sortedServiceNames(compose.Services) {
		state := states[name]
		live, ok := psByName[name]
		if !ok {
			// Service declared but not yet visible to Compose Ps;
			// treat as RunningOnly so the loop keeps waiting.
			allStabilized = false
			continue
		}
		cs := domain.ParseContainerState(live.State)

		// Update restart-observation counter.
		// MVP-Limitation (slice plan §176-178): the counter resets
		// on ANY non-restarting observation, including a brief
		// `running` between restart ticks. A pathological restart-
		// loop that flashes through `running` (e.g. `restarting →
		// running → exited → restarting → ...`) can mask itself in
		// the no-healthcheck/no-port stabilization path. V1
		// follow-up tracks multi-tick history; M6 accepts the
		// trade-off for simpler counter semantics.
		if cs == domain.StateRestarting {
			state.restartObservations++
		} else {
			state.restartObservations = 0
		}

		// Emit unknown-state diagnostic once per service per call.
		if cs == domain.StateUnknown && !state.unknownStateReported {
			*diagnostics = append(*diagnostics, domain.Diagnostic{
				ID:       fmt.Sprintf("up.state.%s.unknown", name),
				Severity: domain.SeverityWarn,
				Message:  fmt.Sprintf("Service %q reports container state %q which u-boot does not classify; treated as running-only.", name, live.State),
				Hint:     "If the state persists past --timeout, the call ends with a stabilization timeout.",
			})
			state.unknownStateReported = true
		}

		outcome := s.classifyOne(ctx, cs, live, state, name, diagnostics)
		switch outcome {
		case domain.OutcomeFailed:
			return false, name, live.State
		case domain.OutcomeStabilized:
			// keep allStabilized true so far
		default:
			allStabilized = false
		}
	}
	return allStabilized, "", ""
}

// classifyOne computes the per-service outcome from a normalized
// container state plus the live healthcheck/port observations.
// May append warn diagnostics to `diagnostics` via classifyRunning
// when the healthcheck-dominated path discovers an unreachable
// declared TCP port (LH-FA-UP-001 §968).
func (s *UpService) classifyOne(ctx context.Context, cs domain.ContainerState, live driven.ComposeService, state *servicePollState, name string, diagnostics *[]domain.Diagnostic) domain.StabilizationOutcome {
	switch cs {
	case domain.StateDead:
		return domain.OutcomeFailed
	case domain.StateRestarting:
		if state.restartObservations >= domain.RestartLoopThreshold {
			return domain.OutcomeFailed
		}
		return domain.OutcomeRunningOnly
	case domain.StateStarting, domain.StateUnknown:
		return domain.OutcomeRunningOnly
	case domain.StateRunning:
		return s.classifyRunning(ctx, live, state, name, diagnostics)
	default:
		return domain.OutcomeRunningOnly
	}
}

// classifyRunning handles the LH-FA-UP-001 §966–§969 matrix:
//
//   - Healthcheck required + `healthy` → Stabilized. Declared TCP
//     ports are still probed (LH-FA-UP-001 §968) but a probe
//     failure does NOT veto stabilization — it emits a one-shot
//     `up.port.<service>.unreachable` Warn diagnostic instead
//     (slice plan §141: "Healthcheck dominiert die Klassifikation").
//   - Healthcheck required + not yet `healthy` → RunningOnly,
//     no probe (Compose owns the health gate).
//   - No healthcheck → port probe gates stabilization: any failure
//     drops to RunningOnly, all probes pass means Stabilized
//     (§967 + §968).
func (s *UpService) classifyRunning(ctx context.Context, live driven.ComposeService, state *servicePollState, name string, diagnostics *[]domain.Diagnostic) domain.StabilizationOutcome {
	if state.healthcheckRequired {
		if !strings.EqualFold(live.Health, "healthy") {
			return domain.OutcomeRunningOnly
		}
		s.probePortsForWarn(ctx, name, state, diagnostics)
		return domain.OutcomeStabilized
	}
	for _, target := range state.probeTargets {
		if err := s.net.DialTCP(ctx, target.Host, target.Port, dialTimeout); err != nil {
			return domain.OutcomeRunningOnly
		}
	}
	return domain.OutcomeStabilized
}

// probePortsForWarn implements the slice plan §141 "Healthcheck
// dominiert" branch: when the service has reached `healthy`, the
// declared TCP ports are STILL probed per LH-FA-UP-001 §968, but
// a probe failure emits a one-shot Warn diagnostic instead of
// vetoing stabilization. The first unreachable port short-circuits
// the loop and arms portUnreachableReported so subsequent
// iterations do not duplicate the diagnostic.
func (s *UpService) probePortsForWarn(ctx context.Context, name string, state *servicePollState, diagnostics *[]domain.Diagnostic) {
	if state.portUnreachableReported {
		return
	}
	for _, target := range state.probeTargets {
		if err := s.net.DialTCP(ctx, target.Host, target.Port, dialTimeout); err != nil {
			*diagnostics = append(*diagnostics, domain.Diagnostic{
				ID:       fmt.Sprintf("up.port.%s.unreachable", name),
				Severity: domain.SeverityWarn,
				Message:  fmt.Sprintf("Service %q reached `healthy` but the declared port %s:%d is not reachable on the host (%v).", name, target.Host, target.Port, err),
				Hint:     "Verify the host port mapping in compose.yaml and any host firewall rules.",
			})
			state.portUnreachableReported = true
			return
		}
	}
}

// parseServicePortsForState walks ports[] and returns the probable
// targets only — non-probable elements are skipped (the caller
// emits the warn diagnostics separately to keep state-build pure).
func parseServicePortsForState(ports []any) ([]portProbeTarget, int) {
	targets := make([]portProbeTarget, 0, len(ports))
	skipped := 0
	for _, raw := range ports {
		t, probable := parseComposePort(raw)
		if probable {
			targets = append(targets, t)
		} else {
			skipped++
		}
	}
	return targets, skipped
}

// healthcheckRequired returns true iff the compose service declares
// a healthcheck mapping that is not disabled. The Compose convention
// `healthcheck: { disable: true }` opts out of healthcheck-based
// stabilization, falling back to running-only.
func healthcheckRequired(def composeServiceDecode) bool {
	if def.Healthcheck == nil {
		return false
	}
	return !def.Healthcheck.Disable
}

// buildResult assembles the LH-FA-UP-003 status snapshot from the
// latest ComposePs reply plus the polling-loop's diagnostics.
// Services are sorted alphabetically for deterministic CLI output.
func buildResult(ps []driven.ComposeService, stabilized bool, diagnostics []domain.Diagnostic) domain.UpResult {
	statuses := make([]domain.ServiceStatus, 0, len(ps))
	for _, svc := range ps {
		statuses = append(statuses, domain.ServiceStatus{
			Name:            svc.Name,
			ContainerStatus: domain.ParseContainerState(svc.State),
			Port:            strings.Join(svc.Ports, ", "),
			Healthcheck:     svc.Health,
		})
	}
	sort.Slice(statuses, func(i, j int) bool { return statuses[i].Name < statuses[j].Name })
	return domain.UpResult{
		Services:    statuses,
		Stabilized:  stabilized,
		Diagnostics: diagnostics,
	}
}

// pendingServiceNames returns the alphabetical list of services
// that did not reach [domain.OutcomeStabilized] in the final
// iteration — surfaces in the timeout error so the CLI can render
// "still waiting for X, Y".
//
// The function mirrors [UpService.classifyOne]'s logic per service
// instead of the cruder "not running+healthy" heuristic that an
// earlier T4-svc revision used: a running-only service without a
// healthcheck (LH-FA-UP-001 §967 stabilization path) must NOT be
// reported as pending. Re-running NetProbe.DialTCP here is
// avoided — the network is async-unsafe from a no-context callsite,
// and over-reporting "running + no healthcheck + has probable
// ports" as pending is the conservative choice (matches the user-
// facing wording: "still waiting for the port to answer").
func pendingServiceNames(ps []driven.ComposeService, compose composeFileDecode, states map[string]*servicePollState) []string {
	psByName := make(map[string]driven.ComposeService, len(ps))
	for _, svc := range ps {
		psByName[svc.Name] = svc
	}
	pending := make([]string, 0)
	for _, name := range sortedServiceNames(compose.Services) {
		live, ok := psByName[name]
		if !ok {
			pending = append(pending, name)
			continue
		}
		cs := domain.ParseContainerState(live.State)
		if cs != domain.StateRunning {
			// Not running → not yet stabilized (StateStarting,
			// StateRestarting under threshold, StateUnknown).
			pending = append(pending, name)
			continue
		}
		state := states[name]
		if state == nil {
			// Shouldn't happen — every declared service has a
			// poll state — but treat missing state as pending to
			// avoid false-positive stabilization.
			pending = append(pending, name)
			continue
		}
		if state.healthcheckRequired {
			// Running with healthcheck required: only stabilized
			// when health == healthy.
			if !strings.EqualFold(live.Health, "healthy") {
				pending = append(pending, name)
			}
			continue
		}
		if len(state.probeTargets) > 0 {
			// Running, no healthcheck, but has probable ports:
			// we can't re-probe from here without a ctx/NetProbe;
			// conservatively list as pending (better to over-
			// report than to hide a real port-down).
			pending = append(pending, name)
			continue
		}
		// Running, no healthcheck, no probable ports → stabilized
		// per LH-FA-UP-001 §967, NOT pending.
	}
	return pending
}

// sortedServiceNames returns the alphabetical list of service
// names from compose.yaml — used so diagnostics emit in
// deterministic order regardless of map iteration.
func sortedServiceNames(services map[string]composeServiceDecode) []string {
	names := make([]string, 0, len(services))
	for name := range services {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
