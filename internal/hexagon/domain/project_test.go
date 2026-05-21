package domain_test

import (
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

func TestNewProject_SetsCurrentSchemaVersion(t *testing.T) {
	name, err := domain.NewProjectName("demo")
	if err != nil {
		t.Fatalf("NewProjectName: %v", err)
	}

	p := domain.NewProject(name)

	if p.Name != name {
		t.Fatalf("Project.Name = %q, want %q", p.Name, name)
	}
	if p.SchemaVersion != domain.SchemaVersionCurrent {
		t.Fatalf("Project.SchemaVersion = %d, want %d", p.SchemaVersion, domain.SchemaVersionCurrent)
	}
}

func TestSchemaVersionCurrent_IsOne(t *testing.T) {
	// Why: anchors LH-DA-003 (initial schemaVersion is 1). If this
	// changes, LH-DA-004 (migration) becomes relevant and the test
	// reminds us.
	if domain.SchemaVersionCurrent != 1 {
		t.Fatalf("SchemaVersionCurrent = %d, want 1 (LH-DA-003)", domain.SchemaVersionCurrent)
	}
}
