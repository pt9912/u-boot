package domain

// Project is the aggregate root for a u-boot project. It is constructed
// by the application layer from a validated [ProjectName] and the
// initial schema version (LH-FA-CONF-002, LH-DA-003).
//
// Project is intentionally minimal at M3 — it grows as the
// application layer adds services, devcontainer configuration, and
// template metadata (`LH-FA-CONF-002` optional V1 fields).
type Project struct {
	Name          ProjectName
	SchemaVersion int
}

// SchemaVersionCurrent is the schema version u-boot writes into newly
// created `u-boot.yaml` files (LH-FA-CONF-002, LH-DA-003).
const SchemaVersionCurrent = 1

// NewProject builds a Project with the current schema version.
// Validation of the name happens in [NewProjectName]; NewProject
// trusts an already-validated input and just bundles it.
func NewProject(name ProjectName) Project {
	return Project{
		Name:          name,
		SchemaVersion: SchemaVersionCurrent,
	}
}
