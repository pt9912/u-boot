package application_test

import "testing"

// envFileOption is the first `runArgs` entry of every generated
// devcontainer.json (LH-FA-DEV-012); the literal is pinned here and not
// taken from the production constant.
const envFileOption = "--env-file=${localWorkspaceFolder}/.devcontainer/devcontainer.env"

// runArgsAfterEnvFile returns the `runArgs` entries of a parsed
// devcontainer.json behind the environment file line and fails when that
// line is not the first entry.
func runArgsAfterEnvFile(t *testing.T, m map[string]any) []any {
	t.Helper()
	all, ok := m["runArgs"].([]any)
	if !ok || len(all) == 0 || all[0] != envFileOption {
		t.Fatalf("runArgs = %v, want %s as first entry", m["runArgs"], envFileOption)
	}
	return all[1:]
}
