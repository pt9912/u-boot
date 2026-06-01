package application_test

import (
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
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
