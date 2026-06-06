package cli

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// Test-only bridges for unexported helpers in statusview.go. The
// helpers stay package-private in production; tests in
// `cli_test` package access them through these wrappers. Pattern
// borrowed from `internal/hexagon/application/export_test.go`.

// RenderUpStatusForTest exposes [renderUpStatus] (M6-T6).
func RenderUpStatusForTest(out io.Writer, services []domain.ServiceStatus) error {
	return renderUpStatus(out, services)
}

// RenderUpDiagnosticsForTest exposes [renderUpDiagnostics] (M6-T6).
func RenderUpDiagnosticsForTest(out io.Writer, diagnostics []domain.Diagnostic, quiet bool) {
	renderUpDiagnostics(out, diagnostics, quiet)
}

// RenderDownSuccessForTest exposes [renderDownSuccess] (M6-T6).
func RenderDownSuccessForTest(out io.Writer, removedVolumes, quiet bool) {
	renderDownSuccess(out, removedVolumes, quiet)
}

// MinimalEnvelopeForTest exposes newMinimalEnvelope as a JSON-bytes
// helper for slice-v1-cli-json-dry-run-doctor T2 Marshal-Pins.
// Returns the marshalled bytes plus the error from json.Marshal,
// so tests can pin both the body and the absence of errors.
func MinimalEnvelopeForTest(command, subcommand string, diagnostics []DiagnosticItemForTest, exitCode int) ([]byte, error) {
	items := make([]diagnosticItem, len(diagnostics))
	for i, d := range diagnostics {
		items[i] = diagnosticItem(d)
	}
	return marshalEnvelopeForTest(newMinimalEnvelope(command, subcommand, items, exitCode))
}

// DataEnvelopeForTest exposes newDataEnvelope as a JSON-bytes
// helper for slice-v1-cli-json-dry-run-generate T0-(p) Marshal-
// Pins. Minimalkontrakt-Form mit zusätzlichem `data`-Träger;
// `data == nil` lässt das Feld via omitempty aus dem JSON
// fallen — anti-drift Pin gegen versehentliches `"data":null`.
func DataEnvelopeForTest(command, subcommand string, data any, diagnostics []DiagnosticItemForTest, exitCode int) ([]byte, error) {
	items := make([]diagnosticItem, len(diagnostics))
	for i, d := range diagnostics {
		items[i] = diagnosticItem(d)
	}
	return marshalEnvelopeForTest(newDataEnvelope(command, subcommand, data, items, exitCode))
}

// FullEnvelopeForTest exposes newFullEnvelope, same shape as
// MinimalEnvelopeForTest. The dryRun=false/diff=false pin (M1
// anti-drift) lives in jsonenvelope_test.go and exercises the
// *bool-vs-bool boundary.
func FullEnvelopeForTest(
	command, subcommand string,
	dryRun, diff bool,
	planned []PlannedFileForTest,
	changes []ChangeEntryForTest,
	diagnostics []DiagnosticItemForTest,
	exitCode int,
) ([]byte, error) {
	pfs := make([]plannedFile, len(planned))
	for i, p := range planned {
		// Field-by-field: PlannedFileForTest predates the Hunks field
		// (slice-v1-cli-json-dry-run-add T4 / T0-(l)); test code that
		// doesn't set Hunks gets the omitempty path.
		pfs[i] = plannedFile{Path: p.Path, Action: p.Action}
	}
	chs := make([]changeEntry, len(changes))
	for i, c := range changes {
		chs[i] = changeEntry(c)
	}
	items := make([]diagnosticItem, len(diagnostics))
	for i, d := range diagnostics {
		items[i] = diagnosticItem(d)
	}
	return marshalEnvelopeForTest(newFullEnvelope(command, subcommand, dryRun, diff, pfs, chs, nil, items, exitCode))
}

// DiagnosticItemForTest mirrors the unexported diagnosticItem so
// the external _test package can construct instances without
// reaching into the cli package directly.
type DiagnosticItemForTest struct {
	Level   string `json:"level"`
	Code    string `json:"code"`
	Message string `json:"message"`
	File    string `json:"file,omitempty"`
}

// PlannedFileForTest mirrors plannedFile.
type PlannedFileForTest struct {
	Path   string `json:"path"`
	Action string `json:"action"`
}

// ChangeEntryForTest mirrors changeEntry.
type ChangeEntryForTest struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

// marshalEnvelopeForTest centralises the marshal call so both
// constructors share the same encoding (no SetIndent, no
// SetEscapeHTML tweaks — Spec §1809 wire-level JSON, byte order
// of fields per struct definition).
func marshalEnvelopeForTest(env cliJSONEnvelope) ([]byte, error) {
	return json.Marshal(env)
}

// JSONSliceSuffixForTest exposes [jsonSliceSuffix] (slice-v1-cli-
// json-dry-run-doctor T3). External _test packages call this to
// pin the path → slice-suffix mapping used in error messages.
func JSONSliceSuffixForTest(cmdPath string) string {
	return jsonSliceSuffix(cmdPath)
}

// JSONAllowlistPathsForTest returns the keys of the runtime allowlist
// (slice-v1-cli-json-dry-run-doctor T3 §M2-Drift). Used by the Cobra-
// Tree-Walk-Test to compare the registered Cobra paths against the
// migrated set.
func JSONAllowlistPathsForTest() []string {
	m := jsonAllowlist()
	out := make([]string, 0, len(m))
	for path := range m {
		out = append(out, path)
	}
	return out
}

// WalkRootCommandPathsForTest enumerates every command path under the
// constructed root: all runnable leaves PLUS every help-parent
// (cobra.NoArgs without RunE → prints help; cobra.NoArgs with
// RunE=cmd.Help() also counts as a user-facing form, e.g. bare
// `template`). The walk powers the Cobra-Tree-Walk-Anti-Drift-Pin
// in jsonallowlist_test.go.
func (a *App) WalkRootCommandPathsForTest() []string {
	root := buildRootCommand(a)
	var paths []string
	walkCommandPaths(root, &paths)
	return paths
}

func walkCommandPaths(cmd *cobra.Command, out *[]string) {
	for _, child := range cmd.Commands() {
		if child.Hidden || child.Name() == "help" || child.Name() == "completion" {
			continue
		}
		*out = append(*out, child.CommandPath())
		walkCommandPaths(child, out)
	}
}
