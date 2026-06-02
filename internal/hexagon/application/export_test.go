package application

import (
	"context"
	"testing"

	yamladapter "github.com/pt9912/u-boot/internal/adapter/driven/yaml"
	"github.com/pt9912/u-boot/internal/hexagon/application/managedblock"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
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

// ServiceCatalogueEntryForTest is the test-package projection of the
// unexported [serviceCatalogueEntry]. Same fields — only renamed
// because the internal type is unexported. slice-v1-otel T1 added
// `ExtraFiles` so tests can pin the per-service whole-file artefact
// declarations (`otel-collector-config.yaml` for OTel).
type ServiceCatalogueEntryForTest struct {
	ComposeTmpl string
	EnvTmpl     string
	VolumeTmpl  string
	ExtraFiles  []ExtraFileEntryForTest
}

// ExtraFileEntryForTest mirrors the unexported [extraFileEntry] so
// slice-v1-otel T1 tests can assert the (path, tmpl) tuples per
// service without depending on package internals.
type ExtraFileEntryForTest struct {
	Path string
	Tmpl string
}

// ServiceCatalogueForTest exposes the unexported [serviceCatalogue]
// lookup so slice-v1-keycloak / slice-v1-otel tests can pin the
// per-service template paths.
func ServiceCatalogueForTest() map[string]ServiceCatalogueEntryForTest {
	out := map[string]ServiceCatalogueEntryForTest{}
	for k, v := range serviceCatalogue() {
		entry := ServiceCatalogueEntryForTest{
			ComposeTmpl: v.composeTmpl,
			EnvTmpl:     v.envTmpl,
			VolumeTmpl:  v.volumeTmpl,
		}
		for _, xf := range v.extraFiles {
			entry.ExtraFiles = append(entry.ExtraFiles, ExtraFileEntryForTest(xf))
		}
		out[k] = entry
	}
	return out
}

// RenderedExtraFileForTest is the test-package view of a single
// rendered ExtraFiles entry from [renderServiceTemplates] (slice-v1-
// otel T1).
type RenderedExtraFileForTest struct {
	Path    string
	Content []byte
}

// RenderServiceTemplatesForTest exposes the per-service render
// pipeline so T1 tests can pin Postgres Byte-Identity, Keycloak's
// nil-VolumeFragment for volume-less services, and slice-v1-otel
// T1's ExtraFiles for whole-file artefacts.
func RenderServiceTemplatesForTest(svc domain.ServiceName) (composeFrag, volumeFrag, envVars []byte, extraFiles []RenderedExtraFileForTest, err error) {
	s := &AddServiceService{}
	tmpls, err := s.renderServiceTemplates(svc)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	for _, xf := range tmpls.ExtraFiles {
		extraFiles = append(extraFiles, RenderedExtraFileForTest(xf))
	}
	return tmpls.ServiceFragment, tmpls.VolumeFragment, tmpls.EnvVariables, extraFiles, nil
}

// IsSupportedServiceForTest exposes the unexported
// [isSupportedService] catalogue check so slice-v1-keycloak tests
// can pin the T1-vs-T2 Catalogue-Erweiterung.
func IsSupportedServiceForTest(svc domain.ServiceName) bool {
	return isSupportedService(svc)
}

// HasRequiredEnvKeysForTest exposes the unexported
// [hasRequiredEnvKeysFor] env-block completeness check so slice-v1-
// keycloak T2 tests can pin the per-service required-keys lookup.
func HasRequiredEnvKeysForTest(svc domain.ServiceName, blockBody []byte) bool {
	return hasRequiredEnvKeysFor(svc, blockBody)
}

// HasRequiredServiceFieldsForTest exposes the unexported
// [hasRequiredServiceFieldsFor] service-block completeness check so
// slice-v1-keycloak T2 tests can pin the per-service scan rules
// (env-keys + volume-ref + volumeOptional skip).
func HasRequiredServiceFieldsForTest(svc domain.ServiceName, blockBody []byte) bool {
	return hasRequiredServiceFieldsFor(svc, blockBody)
}

// CollectActiveServicePortsForTest exposes the package-internal
// `collectActiveServicePorts` helper so the T5 anti-drift test can
// pin that the generator and the doctor `devcontainer.forwardPorts.
// consistency` check share the same forwardPorts source.
func CollectActiveServicePortsForTest(
	fs driven.FileSystem, yamlCodec driven.YAMLCodec, baseDir string, services []string,
) ([]int, error) {
	return collectActiveServicePorts(fs, yamlCodec, baseDir, services)
}

// StripJSONCForTest exposes the package-internal `stripJSONC` helper
// so the T5 devcontainer.json validity tests can pre-process the
// rendered JSONC into plain JSON before passing it to
// `encoding/json.Valid` / `Unmarshal`.
func StripJSONCForTest(src []byte) []byte {
	return stripJSONC(src)
}

// AddServicePlanForTest is the test-only projection of the unexported
// [servicePlan] returned by [AddServiceService.planAdd]. T3 tests use
// it to assert plan shape for each mutating state without exposing
// the production type to non-test callers.
type AddServicePlanForTest struct {
	Service    domain.ServiceName
	PriorState domain.ServiceState
	Action     string
}

// DependenciesForTest exposes the unexported [dependenciesFor]
// catalogue side-table so the slice-v1-addons-deps tests can pin
// the per-service dependency declarations without depending on
// the package internals. Returns nil for every service today
// (postgres has no deps; keycloak / otel land in their own slices).
func DependenciesForTest(svc domain.ServiceName) []domain.AddOnDependency {
	return dependenciesFor(svc)
}

// ResolveAddDependenciesForTest exposes the unexported
// [resolveAddDependencies] resolver so slice-v1-addons-deps T2
// tests can drive it with synthetic [domain.AddOnDependency]
// inputs against a fixture u-boot.yaml without going through the
// AddServiceUseCase wiring. Returns the list of services that
// must be registered before the add request can proceed.
func ResolveAddDependenciesForTest(t *testing.T, yamlBody []byte, deps []domain.AddOnDependency) []domain.ServiceName {
	t.Helper()
	cfg := mustParseUBootYAML(t, yamlBody)
	return resolveAddDependencies(cfg, deps)
}

// ResolveScalarPathForTest exposes the unexported
// [resolveScalarPath] helper so slice-v1-addons-deps T2 tests can
// pin the path → canonical-string mapping per known path.
func ResolveScalarPathForTest(t *testing.T, yamlBody []byte, path string) string {
	t.Helper()
	cfg := mustParseUBootYAML(t, yamlBody)
	return resolveScalarPath(cfg, path)
}

// CheckAddDependenciesForTest exposes the unexported
// [AddServiceService.checkAddDependencies] orchestrator so slice-
// v1-addons-deps tests can drive the full integration path (load
// + resolve + four-mode dispatch) with synthetic dependency
// declarations. The wrapped req carries the four-mode flags so the
// caller selects the dispatch arm under test.
func (s *AddServiceService) CheckAddDependenciesForTest(baseDir string, svc domain.ServiceName, deps []domain.AddOnDependency) error {
	return s.checkAddDependencies(context.Background(), driving.AddServiceRequest{BaseDir: baseDir, ServiceName: svc}, deps)
}

// HandleMissingDependenciesForTest exposes the unexported four-
// mode dispatch directly so slice-v1-addons-deps T3 tests can pin
// each arm (--with-deps, --yes, --no-interactive, interactive
// prompt) without seeding a fixture u-boot.yaml that triggers the
// missing-deps condition. The full Add() recursion is exercised
// because the dispatch calls back into [AddServiceService.Add].
func (s *AddServiceService) HandleMissingDependenciesForTest(ctx context.Context, req driving.AddServiceRequest, missing []domain.ServiceName) error {
	return s.handleMissingDependencies(ctx, req, missing)
}

// FindMissingDependenciesForTest exposes the unexported load +
// resolve helper so tests can pin the resolver wiring against a
// fixture u-boot.yaml without going through the dispatch.
func (s *AddServiceService) FindMissingDependenciesForTest(baseDir string, deps []domain.AddOnDependency) ([]domain.ServiceName, error) {
	return s.findMissingDependencies(baseDir, deps)
}

// mustParseUBootYAML deserialises a u-boot.yaml fixture into the
// in-memory config struct via the production yaml.v3 adapter,
// failing the test if the fixture is malformed.
func mustParseUBootYAML(t *testing.T, body []byte) ubootYAMLConfig {
	t.Helper()
	var cfg ubootYAMLConfig
	codec := yamladapter.New()
	if err := codec.Unmarshal(body, &cfg); err != nil {
		t.Fatalf("parse u-boot.yaml fixture: %v", err)
	}
	return cfg
}

// DetectServiceStateForTest exposes the unexported
// [AddServiceService.detectServiceState] helper so T3 fixtures can
// assert state classification directly, without going through
// Add()'s dispatch.
func (s *AddServiceService) DetectServiceStateForTest(baseDir string, svc domain.ServiceName) (domain.ServiceState, error) {
	return detectServiceState(s.fs, s.yaml, baseDir, svc)
}

// ValidateFeatureSourceForTest exposes the slice-v1-devcontainer-
// features T1 validator for a single allowlist entry. External tests
// drive the failure-table (empty / no-scheme / bad-scheme / no-host)
// without depending on the unexported function.
func ValidateFeatureSourceForTest(raw string) error {
	return validateFeatureSource(raw)
}

// DedupeFeatureSourcesForTest exposes the T1 silent-dedupe helper so
// tests can pin `spec/lastenheft.md:1352` (silent dedupe, first-
// occurrence-order preserved).
func DedupeFeatureSourcesForTest(in []string) []string {
	return dedupeFeatureSources(in)
}

// NormaliseFeatureSourcesForTest exposes the T1 validate-then-dedupe
// pipeline so external tests can pin both the validation-error path
// and the canonical-shape success path.
func NormaliseFeatureSourcesForTest(in []string) ([]string, error) {
	return normaliseFeatureSources(in)
}

// ValidateDevcontainerFeaturesForTest parses a u-boot.yaml fixture
// and runs validateDevcontainerFeatures on the resulting
// devcontainer subtree. The bridge lets external tests pin the full
// load → validate flow for the T1 schema additions without exposing
// the unexported config struct.
func ValidateDevcontainerFeaturesForTest(t *testing.T, yamlBody []byte) error {
	t.Helper()
	cfg := mustParseUBootYAML(t, yamlBody)
	return validateDevcontainerFeatures(cfg.Devcontainer)
}

// FeatureCatalogueEntryForTest is the test-package projection of the
// unexported [featureCatalogueEntry]. Same fields — only renamed
// because the internal type is unexported. Slice-v1-devcontainer-
// features T2 added the catalogue.
type FeatureCatalogueEntryForTest struct {
	Source         string
	DefaultVersion string
	ShortDesc      string
}

// FeatureCatalogueForTest exposes the unexported [featureCatalogue]
// lookup so T2 tests can pin the per-feature source/version triples
// without depending on package internals.
func FeatureCatalogueForTest() map[string]FeatureCatalogueEntryForTest {
	out := map[string]FeatureCatalogueEntryForTest{}
	for k, v := range featureCatalogue() {
		out[k] = FeatureCatalogueEntryForTest{
			Source:         v.source,
			DefaultVersion: v.defaultVersion,
			ShortDesc:      v.shortDesc,
		}
	}
	return out
}

// FeatureForTest exposes the unexported [featureFor] catalogue
// lookup so T2 tests can pin the (name → entry, ok) contract.
func FeatureForTest(name domain.FeatureName) (FeatureCatalogueEntryForTest, bool) {
	entry, ok := featureFor(name)
	if !ok {
		return FeatureCatalogueEntryForTest{}, false
	}
	return FeatureCatalogueEntryForTest{
		Source:         entry.source,
		DefaultVersion: entry.defaultVersion,
		ShortDesc:      entry.shortDesc,
	}, true
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
