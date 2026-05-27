package application

import (
	"context"
	"errors"
	iofs "io/fs"
	"path/filepath"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// DoctorService implements [driving.DoctorUseCase]. It runs a fixed
// set of checks against BaseDir and the local environment, returning
// a [domain.DiagnosticReport] that aggregates every check's outcome.
//
// LH-FA-DIAG-001..004: every Check method appends one Diagnostic per
// invocation; severity comes from the check's own success / warn /
// fail logic, not from the service. The service is severity-agnostic
// and just collects.
type DoctorService struct {
	fs     driven.FileSystem
	logger driven.Logger
}

// Static check: DoctorService satisfies the driving port.
var _ driving.DoctorUseCase = (*DoctorService)(nil)

// NewDoctorService constructs the service with the driven adapters
// the M4-T2 checks need. logger accepts nil (routed to noopLogger)
// so tests and dry-runs do not need a stub. Future tranches will add
// more ports (Git, DockerProbe, YAMLCodec, Confirmer-free since
// doctor is non-interactive); the constructor signature grows
// accordingly.
func NewDoctorService(fs driven.FileSystem, logger driven.Logger) *DoctorService {
	if logger == nil {
		logger = noopLogger{}
	}
	return &DoctorService{fs: fs, logger: logger}
}

// doctorCheckID enumerates the stable machine-readable IDs the
// service emits. Pin them as constants so future tranches can extend
// the set without typos that would silently break CI-side log-
// scraping. The naming convention is `<area>.<probe>` (e.g.
// `fs.write-permissions`).
const (
	checkIDWritePermissions = "fs.write-permissions"
)

// Check runs every M4-T2 check against req.BaseDir and assembles
// the diagnostic report. Checks run in a deterministic order; the
// service does not parallelize (filesystem and external-binary
// checks are I/O-bound but cheap enough sequentially for the MVP).
//
// The use-case-level error is reserved for *fatal* problems that
// prevent any check from running (e.g. an invalid request). Per-
// check failures become [domain.SeverityError] Diagnostics in the
// report, not Go errors.
func (s *DoctorService) Check(ctx context.Context, req driving.DoctorRequest) (driving.DoctorResponse, error) {
	if req.BaseDir == "" {
		return driving.DoctorResponse{}, errors.New("BaseDir is required")
	}
	s.logger.Debug("doctor: starting checks", "baseDir", req.BaseDir)
	items := []domain.Diagnostic{
		s.checkWritePermissions(ctx, req.BaseDir),
	}
	report := domain.DiagnosticReport{Items: items}
	s.logger.Info("doctor: checks complete",
		"baseDir", req.BaseDir,
		"items", len(items),
		"maxSeverity", report.MaxSeverity().String(),
	)
	return driving.DoctorResponse{Report: report}, nil
}

// checkWritePermissions verifies the service can create+remove a
// file in BaseDir. The actual filesystem probe is a
// [driven.FileSystem.WriteFileExclusive] + [driven.FileSystem.RemoveAll]
// pair on a sentinel path. The choice of WriteFileExclusive (instead
// of WriteFile) means: when the sentinel already exists for some
// reason, the check classifies as Error rather than silently
// over-writing — that's the user's own footprint we don't want to
// disturb.
//
// Classifications:
//   - success → SeverityOK, no hint.
//   - any write error → SeverityError with a `chmod` hint and the
//     underlying error in the message.
func (s *DoctorService) checkWritePermissions(_ context.Context, baseDir string) domain.Diagnostic {
	sentinel := filepath.Join(baseDir, ".u-boot-doctor-probe")
	err := s.fs.WriteFileExclusive(sentinel, []byte("probe\n"), 0o600)
	if err != nil {
		// Distinguish "sentinel already exists" (likely user-side
		// junk, not a permission problem) from a real permission
		// problem so the hint is honest.
		if errors.Is(err, iofs.ErrExist) {
			return domain.Diagnostic{
				ID:       checkIDWritePermissions,
				Severity: domain.SeverityError,
				Message:  "Cannot probe write permissions: sentinel file already exists at " + sentinel + ".",
				Hint:     "Remove " + sentinel + " and re-run doctor.",
			}
		}
		return domain.Diagnostic{
			ID:       checkIDWritePermissions,
			Severity: domain.SeverityError,
			Message:  "BaseDir is not writable: " + err.Error() + ".",
			Hint:     "Check directory ownership and permissions (e.g. `chmod u+w " + baseDir + "`).",
		}
	}
	// Cleanup; ignore the error: if the sentinel cannot be removed
	// the next doctor run will hit the ErrExist branch above and the
	// user gets a focused message. Logging it here at Warn would
	// double-emit.
	_ = s.fs.RemoveAll(sentinel)
	return domain.Diagnostic{
		ID:       checkIDWritePermissions,
		Severity: domain.SeverityOK,
		Message:  "BaseDir is writable.",
	}
}
