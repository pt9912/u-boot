package driving

import (
	"context"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// DoctorRequest is the input for [DoctorUseCase.Check]. M4 doctor
// runs a fixed set of checks against the local environment + the
// project at BaseDir; future requests may carry filters (e.g.
// `--only docker.*`) when the check matrix grows.
type DoctorRequest struct {
	// BaseDir is the absolute path of the project directory the
	// doctor inspects. Write-permission, u-boot.yaml, compose.yaml
	// and devcontainer checks are scoped to this path.
	// Environment checks (docker, compose, git availability) are
	// independent of BaseDir.
	BaseDir string
}

// DoctorResponse is the output of [DoctorUseCase.Check]. It carries
// the [domain.DiagnosticReport] verbatim; the CLI adapter decides
// presentation and exit-code mapping. `--strict` (LH-FA-DIAG-003)
// is a CLI-level concern and lives in the adapter, not in the
// response.
type DoctorResponse struct {
	// Report is the aggregate of every check the service ran. Items
	// appear in the service's run order; the CLI may re-sort via
	// [domain.DiagnosticReport.SortedByIssuesFirst].
	Report domain.DiagnosticReport
}

// DoctorUseCase is the driving-port for `u-boot doctor`. The CLI
// adapter holds a reference and calls [Check] from the Cobra command
// handler.
//
// Contract: Check returns a populated DoctorResponse on success. It
// does not return an error for "checks failed" — failures appear as
// SeverityError [domain.Diagnostic] entries in the report. A non-nil
// error from Check signals a *use-case* failure (e.g. the service
// could not run at all — invalid request, fatal I/O error). This
// keeps the report's contract clean: every check's outcome is
// observable in the report, not in the Go error.
type DoctorUseCase interface {
	Check(ctx context.Context, req DoctorRequest) (DoctorResponse, error)
}
