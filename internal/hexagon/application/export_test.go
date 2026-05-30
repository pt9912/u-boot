package application

import (
	"github.com/pt9912/u-boot/internal/hexagon/application/managedblock"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// PortProbeTargetForTest is the test-package-visible view of the
// internal portProbeTarget produced by [parseComposePort] (M6-T4-
// fund). Same shape — only renamed because the internal type is
// unexported and the bridge convention prefers a distinct exported
// name over an alias.
type PortProbeTargetForTest struct {
	Host string
	Port int
}

// ParseComposePortForTest exposes the package-internal
// parseComposePort helper (M6-T4-fund) to external _test packages
// so the eight Compose port-syntax cases can be table-tested
// without going through a full Up() run.
func ParseComposePortForTest(raw any) (PortProbeTargetForTest, bool) {
	t, probable := parseComposePort(raw)
	return PortProbeTargetForTest(t), probable
}

// RenderManagedBlockOnlyForTest exposes the package-internal
// renderManagedBlockOnly helper (M5-T4a) to external _test packages
// so the programmer-error paths (template-missing-marker, malformed)
// and the happy-path byte-extract can be tested without going through
// a full Init() run.
func RenderManagedBlockOnlyForTest(rendered []byte, markerName string) ([]byte, error) {
	return renderManagedBlockOnly(rendered, managedblock.Marker{
		Style: managedblock.StyleHash,
		Name:  markerName,
	})
}

// EnsureComposeScaffoldForTest exposes the package-internal
// ensureComposeScaffold helper (M5-T4a) to external _test packages.
func EnsureComposeScaffoldForTest(content []byte) []byte {
	return ensureComposeScaffold(content)
}

// RenderEnvManagedBlockForTest exposes the package-internal
// renderEnvManagedBlock helper (M5-T4c) to external _test packages
// so the wrap-contract can be pinned without going through Add().
func RenderEnvManagedBlockForTest(svcName string, varsBody []byte) []byte {
	svc, err := domain.NewServiceName(svcName)
	if err != nil {
		panic(err) // test bridge only — invalid name = test bug
	}
	return renderEnvManagedBlock(svc, varsBody)
}

// TemplateNamesForTest exposes the package-internal templateNames
// helper to external _test packages. The `_test.go` suffix means
// the symbol only exists in the test binary; production callers
// cannot reach it.
func TemplateNamesForTest() ([]string, error) {
	return templateNames()
}

// RenderTemplateForTest exposes the package-internal renderTemplate
// helper to external _test packages so the error path
// (template-not-found) is reachable.
func RenderTemplateForTest(name, projectName string) ([]byte, error) {
	return renderTemplate(name, templateData{Name: projectName})
}

// ErrStubHandlerForTest exposes the package-internal
// [errStubHandler] sentinel so external `_test` packages can pin
// `errors.Is(err, errStubHandler)` for the four M7-T1 stub handlers
// without leaking the sentinel into the driving-port API surface.
// T5 removes the test that uses this once the last stub is
// replaced (slice-m7-generate.md T5 DoD).
var ErrStubHandlerForTest = errStubHandler

// AddServicePlanForTest is the test-only projection of the unexported
// [servicePlan] returned by [AddServiceService.planAdd]. T3 tests use
// it to assert plan shape for each mutating state without exposing
// the production type to non-test callers.
type AddServicePlanForTest struct {
	Service    domain.ServiceName
	PriorState domain.ServiceState
	Action     string
}

// DetectServiceStateForTest exposes the unexported
// [AddServiceService.detectServiceState] helper so T3 fixtures can
// assert state classification directly, without going through
// Add()'s dispatch.
func (s *AddServiceService) DetectServiceStateForTest(baseDir string, svc domain.ServiceName) (domain.ServiceState, error) {
	return s.detectServiceState(baseDir, svc)
}

// PlanAddForTest exposes the unexported [AddServiceService.planAdd]
// helper. The returned struct is the test-only projection so the
// production [servicePlan] stays unexported.
func (s *AddServiceService) PlanAddForTest(svc domain.ServiceName, state domain.ServiceState) (AddServicePlanForTest, error) {
	plan, err := s.planAdd(svc, state)
	if err != nil {
		return AddServicePlanForTest{}, err
	}
	return AddServicePlanForTest{
		Service:    plan.Service,
		PriorState: plan.PriorState,
		Action:     plan.Action.String(),
	}, nil
}
