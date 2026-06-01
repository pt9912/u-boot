package application_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const (
	yamlBareProject = "schemaVersion: 1\nproject:\n  name: demo\n"

	yamlPostgresRegistered = "schemaVersion: 1\n" +
		"project:\n  name: demo\n" +
		"services:\n  postgres:\n    enabled: true\n"

	yamlPostgresDisabled = "schemaVersion: 1\n" +
		"project:\n  name: demo\n" +
		"services:\n  postgres:\n    enabled: false\n"

	yamlDevcontainerEnabled = "schemaVersion: 1\n" +
		"project:\n  name: demo\n" +
		"devcontainer:\n  enabled: true\n"
)

// Tests reach the unexported [dependenciesFor] catalogue side-table
// via the [application.DependenciesForTest] export_test.go accessor.
// Pattern mirrors [DetectServiceStateForTest] / [PlanAddForTest].

func TestDependenciesFor_PostgresHasNoDeclaredDependencies(t *testing.T) {
	t.Parallel()
	svc, err := domain.NewServiceName("postgres")
	if err != nil {
		t.Fatalf("NewServiceName(postgres): %v", err)
	}
	deps := application.DependenciesForTest(svc)
	if deps != nil {
		t.Errorf("dependenciesFor(postgres) = %v, want nil (postgres has no deps in v0.3.0)", deps)
	}
}

func TestDependenciesFor_UnknownServiceReturnsNil(t *testing.T) {
	t.Parallel()
	// Catalogue defensiveness: an unknown name (rejected earlier
	// by isSupportedService with ErrServiceUnsupported) returns
	// nil instead of panicking on a missing switch case.
	svc, err := domain.NewServiceName("ghost-service")
	if err != nil {
		t.Fatalf("NewServiceName(ghost-service): %v", err)
	}
	if deps := application.DependenciesForTest(svc); deps != nil {
		t.Errorf("dependenciesFor(ghost-service) = %v, want nil", deps)
	}
}

// --- resolveScalarPath (T2 helper) ---------------------------------------

func TestResolveScalarPath_ProjectName(t *testing.T) {
	t.Parallel()
	got := application.ResolveScalarPathForTest(t, []byte(yamlBareProject), "project.name")
	if got != "demo" {
		t.Errorf("project.name = %q, want %q", got, "demo")
	}
}

func TestResolveScalarPath_DevcontainerEnabledTrue(t *testing.T) {
	t.Parallel()
	got := application.ResolveScalarPathForTest(t, []byte(yamlDevcontainerEnabled), "devcontainer.enabled")
	if got != "true" {
		t.Errorf("devcontainer.enabled = %q, want %q", got, "true")
	}
}

func TestResolveScalarPath_DevcontainerEnabledUnset(t *testing.T) {
	t.Parallel()
	got := application.ResolveScalarPathForTest(t, []byte(yamlBareProject), "devcontainer.enabled")
	if got != "" {
		t.Errorf("devcontainer.enabled (unset) = %q, want \"\"", got)
	}
}

func TestResolveScalarPath_ServicesEnabledTrue(t *testing.T) {
	t.Parallel()
	got := application.ResolveScalarPathForTest(t, []byte(yamlPostgresRegistered), "services.postgres.enabled")
	if got != "true" {
		t.Errorf("services.postgres.enabled = %q, want %q", got, "true")
	}
}

func TestResolveScalarPath_ServicesEnabledFalse(t *testing.T) {
	t.Parallel()
	got := application.ResolveScalarPathForTest(t, []byte(yamlPostgresDisabled), "services.postgres.enabled")
	if got != "false" {
		t.Errorf("services.postgres.enabled (disabled) = %q, want %q", got, "false")
	}
}

func TestResolveScalarPath_ServicesEnabledUnregistered(t *testing.T) {
	t.Parallel()
	got := application.ResolveScalarPathForTest(t, []byte(yamlBareProject), "services.postgres.enabled")
	if got != "" {
		t.Errorf("services.postgres.enabled (unregistered) = %q, want \"\"", got)
	}
}

func TestResolveScalarPath_UnknownPathReturnsEmpty(t *testing.T) {
	t.Parallel()
	for _, path := range []string{"services.x.persistence", "schemaVersion", "foo.bar.baz", "project.name.extra"} {
		got := application.ResolveScalarPathForTest(t, []byte(yamlPostgresRegistered), path)
		if got != "" {
			t.Errorf("path %q = %q, want \"\" (unknown path)", path, got)
		}
	}
}

// --- resolveAddDependencies (T2 resolver) --------------------------------

func TestResolveAddDependencies_EmptyDepsReturnsNil(t *testing.T) {
	t.Parallel()
	got := application.ResolveAddDependenciesForTest(t, []byte(yamlPostgresRegistered), nil)
	if got != nil {
		t.Errorf("got = %v, want nil for empty deps", got)
	}
}

func TestResolveAddDependencies_TriggerNotMet_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Dep would require postgres if services.x.enabled == "external-postgres",
	// but in the fixture services.postgres.enabled is "true", not "external-postgres".
	deps := []domain.AddOnDependency{
		{
			Requires:    mustServiceNameForDeps(t, "postgres"),
			WhenPath:    "services.postgres.enabled",
			EqualsValue: "external-postgres",
		},
	}
	got := application.ResolveAddDependenciesForTest(t, []byte(yamlPostgresRegistered), deps)
	if got != nil {
		t.Errorf("got = %v, want nil (trigger value not matched)", got)
	}
}

func TestResolveAddDependencies_TriggerMet_ServicePresent_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Trigger condition matches AND required service IS registered.
	// Spec semantics: registered (even disabled) counts as PRESENT.
	deps := []domain.AddOnDependency{
		{
			Requires:    mustServiceNameForDeps(t, "postgres"),
			WhenPath:    "devcontainer.enabled",
			EqualsValue: "true",
		},
	}
	got := application.ResolveAddDependenciesForTest(t, []byte(yamlDevcontainerEnabled+"services:\n  postgres:\n    enabled: false\n"), deps)
	if got != nil {
		t.Errorf("got = %v, want nil (postgres registered, even though disabled)", got)
	}
}

func TestResolveAddDependencies_TriggerMet_ServiceAbsent_ReturnsMissing(t *testing.T) {
	t.Parallel()
	deps := []domain.AddOnDependency{
		{
			Requires:    mustServiceNameForDeps(t, "postgres"),
			WhenPath:    "devcontainer.enabled",
			EqualsValue: "true",
		},
	}
	got := application.ResolveAddDependenciesForTest(t, []byte(yamlDevcontainerEnabled), deps)
	if len(got) != 1 || got[0].String() != "postgres" {
		t.Errorf("got = %v, want [postgres]", got)
	}
}

func TestResolveAddDependencies_MultipleDeps_OnlyMatchingFires(t *testing.T) {
	t.Parallel()
	deps := []domain.AddOnDependency{
		{
			Requires:    mustServiceNameForDeps(t, "postgres"),
			WhenPath:    "devcontainer.enabled",
			EqualsValue: "true",
		},
		{
			Requires:    mustServiceNameForDeps(t, "keycloak"),
			WhenPath:    "devcontainer.enabled",
			EqualsValue: "external-keycloak", // does NOT match
		},
	}
	got := application.ResolveAddDependenciesForTest(t, []byte(yamlDevcontainerEnabled), deps)
	if len(got) != 1 || got[0].String() != "postgres" {
		t.Errorf("got = %v, want [postgres] only (keycloak trigger did not match)", got)
	}
}

func TestResolveAddDependencies_DuplicateRequires_DeduplicatedAndOrdered(t *testing.T) {
	t.Parallel()
	// Two distinct paths both pointing at the same Requires service —
	// returned once, in order of first encounter.
	deps := []domain.AddOnDependency{
		{
			Requires:    mustServiceNameForDeps(t, "postgres"),
			WhenPath:    "devcontainer.enabled",
			EqualsValue: "true",
		},
		{
			Requires:    mustServiceNameForDeps(t, "postgres"),
			WhenPath:    "project.name",
			EqualsValue: "demo",
		},
	}
	got := application.ResolveAddDependenciesForTest(t, []byte(yamlDevcontainerEnabled), deps)
	if len(got) != 1 || got[0].String() != "postgres" {
		t.Errorf("got = %v, want [postgres] (deduplicated)", got)
	}
}

// --- AddServiceService.checkAddDependencies integration -----------------

func TestCheckAddDependencies_NoMissing_ReturnsNil(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	if err := fs.WriteFile("/proj/u-boot.yaml", []byte(yamlPostgresRegistered), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	svc := application.NewAddServiceService(fs, &fakeYAML{}, nil)

	// Synthetic dep whose trigger condition does NOT match
	// (devcontainer.enabled is unset in the fixture).
	deps := []domain.AddOnDependency{
		{
			Requires:    mustServiceNameForDeps(t, "postgres"),
			WhenPath:    "devcontainer.enabled",
			EqualsValue: "true",
		},
	}
	if err := svc.CheckAddDependenciesForTest("/proj", mustServiceNameForDeps(t, "postgres"), deps); err != nil {
		t.Fatalf("checkAddDependencies: %v", err)
	}
}

func TestCheckAddDependencies_MissingWrapsSentinel(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	if err := fs.WriteFile("/proj/u-boot.yaml", []byte(yamlDevcontainerEnabled), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	svc := application.NewAddServiceService(fs, &fakeYAML{}, nil)

	// Trigger matches AND required service (postgres) is NOT registered
	// → the integration returns ErrDependenciesRequired with the
	// missing names in the message.
	deps := []domain.AddOnDependency{
		{
			Requires:    mustServiceNameForDeps(t, "postgres"),
			WhenPath:    "devcontainer.enabled",
			EqualsValue: "true",
		},
	}
	err := svc.CheckAddDependenciesForTest("/proj", mustServiceNameForDeps(t, "keycloak"), deps)
	if err == nil {
		t.Fatal("checkAddDependencies: want ErrDependenciesRequired, got nil")
	}
	if !errors.Is(err, driving.ErrDependenciesRequired) {
		t.Errorf("err = %v, want wrap of driving.ErrDependenciesRequired", err)
	}
	if !strings.Contains(err.Error(), "postgres") {
		t.Errorf("err = %v, missing service name %q in message", err, "postgres")
	}
	if !strings.Contains(err.Error(), "--with-deps") {
		t.Errorf("err = %v, missing --with-deps hint in message", err)
	}
}

// --- helpers --------------------------------------------------------------

func mustServiceNameForDeps(t *testing.T, raw string) domain.ServiceName {
	t.Helper()
	name, err := domain.NewServiceName(raw)
	if err != nil {
		t.Fatalf("NewServiceName(%q): %v", raw, err)
	}
	return name
}
