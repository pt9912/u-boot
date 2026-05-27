package application

import (
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"path/filepath"
	"strconv"
	"strings"

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
	yaml   driven.YAMLCodec
	git    driven.Git
	docker driven.DockerProbe
	logger driven.Logger
}

// Static check: DoctorService satisfies the driving port.
var _ driving.DoctorUseCase = (*DoctorService)(nil)

// NewDoctorService constructs the service with the driven adapters
// the M4 checks need. logger accepts nil (routed to noopLogger) so
// tests and dry-runs do not need a stub. Future tranches may add
// more ports (devcontainer probe); the constructor signature grows
// accordingly.
func NewDoctorService(fs driven.FileSystem, yaml driven.YAMLCodec, git driven.Git, docker driven.DockerProbe, logger driven.Logger) *DoctorService {
	if logger == nil {
		logger = noopLogger{}
	}
	return &DoctorService{fs: fs, yaml: yaml, git: git, docker: docker, logger: logger}
}

// doctorCheckID enumerates the stable machine-readable IDs the
// service emits. Pin them as constants so future tranches can extend
// the set without typos that would silently break CI-side log-
// scraping. The naming convention is `<area>.<probe>` (e.g.
// `fs.write-permissions`).
const (
	checkIDWritePermissions = "fs.write-permissions"
	checkIDGitInstalled     = "git.installed"
	checkIDDockerInstalled  = "docker.installed"
	checkIDDockerReachable  = "docker.reachable"
	checkIDComposeInstalled = "docker.compose.installed"
	checkIDUbootYaml        = "uboot.yaml.valid"
	checkIDComposeYaml      = "compose.yaml.valid"
)

// Minimum versions per LH-FA-DIAG-002. The thresholds are MAJOR.MINOR
// pairs; PATCH is informational only.
const (
	minDockerMajor  = 24
	minDockerMinor  = 0
	minComposeMajor = 2
	minComposeMinor = 20
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
		// Use the shared ErrBaseDirMissing sentinel — same semantics
		// as [InitProjectService] for the LH-FA-CLI-006 §10
		// validation mapping. Doctor never invents a default BaseDir
		// (would silently diagnose an unintended path).
		return driving.DoctorResponse{}, fmt.Errorf("%w: doctor.BaseDir is empty", driving.ErrBaseDirMissing)
	}
	s.logger.Debug("doctor: starting checks", "baseDir", req.BaseDir)
	items := []domain.Diagnostic{
		s.checkWritePermissions(ctx, req.BaseDir),
		s.checkGitInstalled(ctx),
		s.checkDockerInstalled(ctx),
		s.checkDockerReachable(ctx),
		s.checkComposeInstalled(ctx),
		s.checkUbootYaml(ctx, req.BaseDir),
		s.checkComposeYaml(ctx, req.BaseDir),
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

// checkGitInstalled probes the git binary availability. Any error
// from [driven.Git.Version] classifies as Error — the M3 init flow
// relies on `git init`, so a missing git binary blocks the typical
// LH-AK-001 use case.
func (s *DoctorService) checkGitInstalled(ctx context.Context) domain.Diagnostic {
	version, err := s.git.Version(ctx)
	if err != nil {
		return domain.Diagnostic{
			ID:       checkIDGitInstalled,
			Severity: domain.SeverityError,
			Message:  "git binary not available: " + err.Error() + ".",
			Hint:     "Install git (e.g. `apt install git`, `brew install git`).",
		}
	}
	return domain.Diagnostic{
		ID:       checkIDGitInstalled,
		Severity: domain.SeverityOK,
		Message:  "git " + version + " available.",
	}
}

// checkDockerInstalled probes the docker client binary + version.
// Missing binary → Error; present but below LH-FA-DIAG-002 minimum
// (24.0) → Error; parseable but unrecognized semver → Warn (we
// observed the binary but can't validate the version, so the user
// should look). At-or-above the minimum → OK.
func (s *DoctorService) checkDockerInstalled(ctx context.Context) domain.Diagnostic {
	version, err := s.docker.Version(ctx)
	if err != nil {
		return domain.Diagnostic{
			ID:       checkIDDockerInstalled,
			Severity: domain.SeverityError,
			Message:  "docker binary not available: " + err.Error() + ".",
			Hint:     "Install Docker Desktop or Docker Engine (https://docs.docker.com/engine/install/).",
		}
	}
	return classifyVersionAtLeast(checkIDDockerInstalled, "docker", version, minDockerMajor, minDockerMinor)
}

// checkDockerReachable probes the docker daemon socket. Reachability
// failures classify as Error — every meaningful u-boot subcommand
// (init, add, up, down, doctor itself for compose-validation) needs
// the daemon eventually.
func (s *DoctorService) checkDockerReachable(ctx context.Context) domain.Diagnostic {
	if err := s.docker.Info(ctx); err != nil {
		return domain.Diagnostic{
			ID:       checkIDDockerReachable,
			Severity: domain.SeverityError,
			Message:  "docker daemon is not reachable: " + err.Error() + ".",
			Hint:     "Start Docker (or check /var/run/docker.sock permissions for the current user).",
		}
	}
	return domain.Diagnostic{
		ID:       checkIDDockerReachable,
		Severity: domain.SeverityOK,
		Message:  "docker daemon is reachable.",
	}
}

// checkComposeInstalled probes the docker compose plugin + version.
// Same classification logic as checkDockerInstalled, scoped to the
// compose plugin (`docker compose version --short`).
func (s *DoctorService) checkComposeInstalled(ctx context.Context) domain.Diagnostic {
	version, err := s.docker.ComposeVersion(ctx)
	if err != nil {
		return domain.Diagnostic{
			ID:       checkIDComposeInstalled,
			Severity: domain.SeverityError,
			Message:  "docker compose plugin not available: " + err.Error() + ".",
			Hint:     "Install the Docker Compose v2 plugin (https://docs.docker.com/compose/install/linux/).",
		}
	}
	return classifyVersionAtLeast(checkIDComposeInstalled, "docker compose", version, minComposeMajor, minComposeMinor)
}

// classifyVersionAtLeast builds the Diagnostic for a version-vs-
// minimum comparison. Shared between `docker` and `docker compose`
// (and reusable for any future semver-min check).
//
// Outcomes:
//   - parse OK + at-or-above minimum → SeverityOK
//   - parse OK + below minimum       → SeverityError with concrete versions in the message
//   - parse fail                     → SeverityWarn (we saw the tool but can't validate)
func classifyVersionAtLeast(id, label, version string, minMajor, minMinor int) domain.Diagnostic {
	major, minor, ok := parseSemverMajorMinor(version)
	if !ok {
		return domain.Diagnostic{
			ID:       id,
			Severity: domain.SeverityWarn,
			Message: fmt.Sprintf("%s reports unrecognized version %q (expected `<major>.<minor>.<patch>`).",
				label, version),
			Hint: fmt.Sprintf("Cannot validate against the LH-FA-DIAG-002 minimum %d.%d — verify manually.",
				minMajor, minMinor),
		}
	}
	if major < minMajor || (major == minMajor && minor < minMinor) {
		return domain.Diagnostic{
			ID:       id,
			Severity: domain.SeverityError,
			Message: fmt.Sprintf("%s %s is below the LH-FA-DIAG-002 minimum %d.%d.",
				label, version, minMajor, minMinor),
			Hint: fmt.Sprintf("Upgrade %s to %d.%d or newer.", label, minMajor, minMinor),
		}
	}
	return domain.Diagnostic{
		ID:       id,
		Severity: domain.SeverityOK,
		Message:  fmt.Sprintf("%s %s available (≥ %d.%d).", label, version, minMajor, minMinor),
	}
}

// parseSemverMajorMinor extracts MAJOR + MINOR from a version string
// like `"24.0.7"` or `"2.20.0-rc1"`. Returns ok=false when the input
// has fewer than two dot-separated leading numeric components.
//
// MAJOR.MINOR is enough for the LH-FA-DIAG-002 thresholds (24.0 /
// 2.20); PATCH and pre-release suffixes are informational only.
//
// The parser is stdlib-only (no semver library) — u-boot's needs are
// narrow and pulling in `golang.org/x/mod/semver` (gomodguard-blocked
// today) is not justified.
func parseSemverMajorMinor(version string) (int, int, bool) {
	parts := strings.SplitN(version, ".", 3)
	if len(parts) < 2 {
		return 0, 0, false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, false
	}
	// MINOR can be followed by `-rcN`, `+build`, etc. Strip
	// non-digit suffix at the first non-digit position.
	minorRaw := trimNonDigitSuffix(parts[1])
	minor, err := strconv.Atoi(minorRaw)
	if err != nil {
		return 0, 0, false
	}
	return major, minor, true
}

// trimNonDigitSuffix returns the prefix of s up to the first
// non-digit character. `"20"` → `"20"`, `"20-rc1"` → `"20"`,
// `"abc"` → `""`.
func trimNonDigitSuffix(s string) string {
	for i, r := range s {
		if r < '0' || r > '9' {
			return s[:i]
		}
	}
	return s
}

// projectFileProbe carries the per-call configuration for
// [DoctorService.loadProjectYAML]. Bundling the args in a struct
// (instead of 5 positional parameters) keeps the call sites readable
// and the helper signature stable as more probes are added (e.g.
// devcontainer.json in T6 will reuse the same scaffold).
type projectFileProbe struct {
	// ID is the diagnostic check ID for every Diagnostic the helper
	// emits.
	ID string
	// BaseDir / RelPath compose the on-disk path (joined via filepath.Join).
	BaseDir, RelPath string
	// Label is the human-readable file name used in
	// probe/read/parse-error messages (e.g. `"u-boot.yaml"`).
	Label string
	// MissingMsg / MissingHint are emitted when the file does not
	// exist (Warn diagnostic). Other LH-FA-DIAG-002 scenarios
	// (probe I/O error, read error, parse error) get standard wording
	// to keep the operator-facing language consistent across files.
	MissingMsg, MissingHint string
}

// loadProjectYAML is the shared scaffold for the file-based YAML
// validations (checkUbootYaml, checkComposeYaml; T6 devcontainer
// will follow the same pattern). Probes existence, reads the body,
// unmarshals into the caller-supplied destination — emitting a
// ready-to-return Diagnostic for the missing/probe-error/read-error/
// parse-error branches. On success it returns (nil, true); the caller
// then runs its file-specific post-parse validation.
//
// The helper centralizes the standard hint wording for the four
// non-success branches so all file probes use identical language;
// branch-specific texts (the success message + validation hints) stay
// in the calling check.
func (s *DoctorService) loadProjectYAML(p projectFileProbe, dst any) (*domain.Diagnostic, bool) {
	path := filepath.Join(p.BaseDir, p.RelPath)
	exists, err := s.fs.Exists(path)
	if err != nil {
		return &domain.Diagnostic{
			ID:       p.ID,
			Severity: domain.SeverityError,
			Message:  "Cannot probe " + p.Label + ": " + err.Error() + ".",
			Hint:     "Check filesystem permissions on " + path + ".",
		}, false
	}
	if !exists {
		return &domain.Diagnostic{
			ID:       p.ID,
			Severity: domain.SeverityWarn,
			Message:  p.MissingMsg,
			Hint:     p.MissingHint,
		}, false
	}
	body, err := s.fs.ReadFile(path)
	if err != nil {
		return &domain.Diagnostic{
			ID:       p.ID,
			Severity: domain.SeverityError,
			Message:  "Cannot read " + p.Label + ": " + err.Error() + ".",
			Hint:     "Check filesystem permissions on " + path + ".",
		}, false
	}
	if err := s.yaml.Unmarshal(body, dst); err != nil {
		return &domain.Diagnostic{
			ID:       p.ID,
			Severity: domain.SeverityError,
			Message:  p.Label + " is not valid YAML: " + err.Error() + ".",
			Hint:     "Fix YAML syntax (indentation, missing colons, mismatched quotes).",
		}, false
	}
	return nil, true
}

// checkUbootYaml validates the `u-boot.yaml` steering file against
// LH-FA-CONF-001..003 / LH-FA-INIT-006:
//
//   - missing file       → Warn (directory is not a u-boot project;
//                          might be intentional, e.g. running doctor
//                          before init).
//   - I/O error on probe → Error.
//   - invalid YAML       → Error with parser message.
//   - schemaVersion ≠ 1  → Error.
//   - missing project.name → Error.
//   - invalid project.name (per LH-FA-INIT-006 regex) → Error.
//   - all checks pass    → OK with project name + schemaVersion in
//                          the message.
//
// The check shares the `ubootYAMLConfig` struct with
// [InitProjectService] (same package, unexported) and uses
// [domain.NewProjectName] for the regex enforcement, so the two
// use-cases stay in lock-step on what "valid u-boot.yaml" means.
func (s *DoctorService) checkUbootYaml(_ context.Context, baseDir string) domain.Diagnostic {
	var cfg ubootYAMLConfig
	diag, ok := s.loadProjectYAML(projectFileProbe{
		ID:          checkIDUbootYaml,
		BaseDir:     baseDir,
		RelPath:     "u-boot.yaml",
		Label:       "u-boot.yaml",
		MissingMsg:  "u-boot.yaml not present — directory is not a u-boot project.",
		MissingHint: "Run `u-boot init` to create one (LH-FA-INIT-001).",
	}, &cfg)
	if !ok {
		return *diag
	}
	if cfg.SchemaVersion != domain.SchemaVersionCurrent {
		return domain.Diagnostic{
			ID:       checkIDUbootYaml,
			Severity: domain.SeverityError,
			Message: fmt.Sprintf("u-boot.yaml schemaVersion is %d (expected %d).",
				cfg.SchemaVersion, domain.SchemaVersionCurrent),
			Hint: fmt.Sprintf("Set `schemaVersion: %d` at the top of the file.", domain.SchemaVersionCurrent),
		}
	}
	if cfg.Project.Name == "" {
		return domain.Diagnostic{
			ID:       checkIDUbootYaml,
			Severity: domain.SeverityError,
			Message:  "u-boot.yaml is missing required `project.name`.",
			Hint:     "Add `project: { name: <valid-name> }` per LH-FA-INIT-006.",
		}
	}
	if _, err := domain.NewProjectName(cfg.Project.Name); err != nil {
		return domain.Diagnostic{
			ID:       checkIDUbootYaml,
			Severity: domain.SeverityError,
			Message:  fmt.Sprintf("u-boot.yaml project.name %q is invalid: %s.", cfg.Project.Name, err.Error()),
			Hint:     "Use a lowercase name like `my-service` (LH-FA-INIT-006 regex).",
		}
	}
	return domain.Diagnostic{
		ID:       checkIDUbootYaml,
		Severity: domain.SeverityOK,
		Message: fmt.Sprintf("u-boot.yaml is valid (project %q, schemaVersion %d).",
			cfg.Project.Name, cfg.SchemaVersion),
	}
}

// composeYAMLShape captures the minimum top-level Compose shape the
// LH-FA-DIAG-002-compose-validation cares about: just the `services:`
// key as a free-form map. Per spec ("minimal Top-Level-Shape"), no
// deeper schema validation happens at this layer — a deeper
// validator (service-level `image`/`build`, port format, network
// references) would be a follow-up slice.
type composeYAMLShape struct {
	Services map[string]any `yaml:"services"`
}

// checkComposeYaml validates the `compose.yaml` Docker Compose file
// per LH-FA-DIAG-002 / spec/lastenheft.md §4.7:
//
//   - missing file        → Warn (LH-FA-INIT-003 names compose.yaml
//                          as part of the mandatory project layout,
//                          but doctor running before init or in a
//                          partial directory is a soft signal).
//   - I/O error on probe  → Error.
//   - invalid YAML        → Error with parser message.
//   - parsed but no `services:` → Error (a Compose file without
//                          services is not a meaningful one).
//   - parsed with services → OK with service count in message.
//
// The exists/read/parse scaffold is shared with [checkUbootYaml] via
// [DoctorService.loadProjectYAML]; the Compose-specific validation
// (`services:`-key presence) lives below.
func (s *DoctorService) checkComposeYaml(_ context.Context, baseDir string) domain.Diagnostic {
	var shape composeYAMLShape
	diag, ok := s.loadProjectYAML(projectFileProbe{
		ID:          checkIDComposeYaml,
		BaseDir:     baseDir,
		RelPath:     "compose.yaml",
		Label:       "compose.yaml",
		MissingMsg:  "compose.yaml not present — directory has no Docker Compose configuration.",
		MissingHint: "Run `u-boot init` (LH-FA-INIT-003 ships a compose.yaml).",
	}, &shape)
	if !ok {
		return *diag
	}
	if len(shape.Services) == 0 {
		return domain.Diagnostic{
			ID:       checkIDComposeYaml,
			Severity: domain.SeverityError,
			Message:  "compose.yaml has no `services:` entries.",
			Hint:     "Add at least one service block (e.g. `services: { app: { image: ... } }`).",
		}
	}
	return domain.Diagnostic{
		ID:       checkIDComposeYaml,
		Severity: domain.SeverityOK,
		Message:  fmt.Sprintf("compose.yaml is valid (%d service(s)).", len(shape.Services)),
	}
}
