// Package cli is the Cobra-based driving adapter for u-boot. It
// translates command-line invocations into driving-port use-case
// calls (LH-FA-ARCH-002, LH-FA-CLI-001..006).
//
// Layer rules (LH-FA-ARCH-003, depguard-enforced): this package may
// import `hexagon/domain`, `hexagon/port/driving`, and external
// libraries (Cobra). It may NOT import `hexagon/application` or
// `adapter/driven` — the wiring layer (`cmd/uboot`) constructs the
// application services and the driven adapters and injects fully-
// constructed driving-port implementations into [New].
package cli

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// App holds the driving-port dependencies the CLI needs.
//
// The struct is intentionally small — one field per use-case port,
// plus environment hooks (getwd) that tests substitute via
// functional options. The LH-FA-CLI-005 persistent verbosity flags
// (--quiet / --verbose / --debug) and the LH-FA-CLI-005A interaction
// flags (--yes / --no-interactive) live here too so subcommands can
// read the parsed values without grovelling through
// cmd.Root().PersistentFlags().
type App struct {
	// version is the build-time version string, surfaced via
	// `u-boot --version`. The wiring layer passes it in; the CLI
	// package does not own version metadata.
	version string

	// initUseCase implements `u-boot init` (LH-FA-INIT-001..007).
	initUseCase driving.InitProjectUseCase

	// doctorUseCase implements `u-boot doctor` (LH-FA-DIAG-001..004).
	doctorUseCase driving.DoctorUseCase

	// addServiceUseCase implements `u-boot add <service>`
	// (LH-FA-ADD-001..002, LH-FA-ADD-005).
	addServiceUseCase driving.AddServiceUseCase

	// upUseCase implements `u-boot up` (LH-FA-UP-001..003).
	upUseCase driving.UpUseCase

	// downUseCase implements `u-boot down` (LH-FA-UP-004).
	downUseCase driving.DownUseCase

	// getwd is the working-directory probe; defaults to os.Getwd.
	// Tests inject a fake via [WithGetwd] so they do not depend on
	// the host pwd.
	getwd func() (string, error)

	// yes and noInteractive are bound to the root command's
	// PersistentFlags by [buildRootCommand].
	yes           bool
	noInteractive bool

	// quiet, verbose, debug are bound to the LH-FA-CLI-005 root
	// PersistentFlags. The doctor subcommand reads --quiet to filter
	// SeverityOK items from the rendered report. --verbose / --debug
	// are accepted per spec but currently do not change the doctor
	// output (logger-level wiring is a follow-up).
	quiet   bool
	verbose bool
	debug   bool
}

// Option mutates an [App] during [New]; the Go-idiomatic functional-
// options pattern keeps the constructor signature stable while
// optional behaviour (test seams, future timeouts) is added.
type Option func(*App)

// WithGetwd overrides the working-directory probe. Intended for
// tests; production callers use [New] without options.
func WithGetwd(fn func() (string, error)) Option {
	return func(a *App) { a.getwd = fn }
}

// New constructs an App. The version string and every use-case
// implementation must be non-nil at call time; the CLI package
// trusts the wiring layer to honor that.
func New(version string, initUC driving.InitProjectUseCase, doctorUC driving.DoctorUseCase, addUC driving.AddServiceUseCase, upUC driving.UpUseCase, downUC driving.DownUseCase, opts ...Option) *App {
	a := &App{
		version:           version,
		initUseCase:       initUC,
		doctorUseCase:     doctorUC,
		addServiceUseCase: addUC,
		upUseCase:         upUC,
		downUseCase:       downUC,
		getwd:             os.Getwd,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Execute parses args and dispatches to the matching subcommand. It
// reads stdin / writes stdout/stderr through the provided streams so
// the wiring layer (and tests) can substitute buffers. Returns the
// CLI-level error (non-nil on bad flag, unknown command, use-case
// failure); the wiring layer maps it to an exit code.
func (a *App) Execute(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	cmd := buildRootCommand(a)
	cmd.SetArgs(args)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	return cmd.ExecuteContext(ctx)
}

// ErrConflictingModeFlags is returned by the init subcommand when
// `--yes` and `--no-interactive` are both set — LH-FA-CLI-005A §235
// declares them mutually exclusive. Lives in the cli package (not
// in `driving`) because the application layer never sees these
// flags; they are pure CLI-level mode switches.
var ErrConflictingModeFlags = errors.New("--yes and --no-interactive are mutually exclusive")

// ErrInvalidTimeout is returned by the M6 up subcommand when
// `--timeout` is a negative integer (LH-FA-UP-001 §965). The CLI
// could not delegate that validation to the application service —
// the application takes a `time.Duration` and could not distinguish
// a deliberate negative value from a unit-mistake-mismatch — so the
// rejection happens at the CLI before construction of the request.
// Maps to LH-FA-CLI-006 exit code 2.
var ErrInvalidTimeout = errors.New("--timeout must be >= 0")

// ErrDoctorFailures signals that `u-boot doctor` ran successfully
// (use-case returned no error) but the diagnostic report contained
// at least one SeverityError item — or at least one SeverityWarn
// when `--strict` was set (LH-FA-DIAG-003). Maps to exit code 11.
//
// Lives in the cli package because the LH-FA-CLI-006 exit-code
// mapping is a CLI concern; the application's DoctorUseCase
// faithfully returns a report and lets the adapter decide.
var ErrDoctorFailures = errors.New("doctor report contains failures")

// ExitCode classifies a CLI error into the u-boot exit-code scheme
// (LH-FA-CLI-006):
//
//   - 0  — no error
//   - 2  — pure CLI / flag errors (unknown subcommand, unknown flag,
//          missing required arg, too many positional args,
//          ErrConflictingModeFlags, M6 `--timeout=-1`)
//   - 10 — fachlicher Validierungsfehler: LH-FA-INIT-004 marker
//          collisions (ErrProjectExists), non-marker file collision
//          (ErrFileExists), LH-FA-INIT-006 invalid project name
//          (ErrInvalidProjectName) or service name
//          (ErrInvalidServiceName), LH-AK-001 missing BaseDir
//          (ErrBaseDirMissing), LH-FA-INIT-005 unsupported
//          backup-source kind (ErrBackupUnsupportedKind), LH-FA-INIT-005
//          §619 force-without-backup (ErrForceRequiresBackup),
//          LH-FA-ADD-001 missing u-boot.yaml
//          (ErrProjectNotInitialized), LH-FA-ADD-002 unknown
//          service (ErrServiceUnsupported), LH-FA-ADD-005
//          inconsistent service state (ErrServiceInconsistent),
//          M6 missing compose.yaml (ErrComposeFileMissing) and
//          destructive confirmation refused (ErrConfirmationRequired)
//   - 11 — fachlicher Umgebungsfehler: `u-boot doctor` reported at
//          least one SeverityError, or at least one SeverityWarn
//          with `--strict` (ErrDoctorFailures, LH-FA-DIAG-003);
//          M6 `u-boot up`/`down` saw a Docker-environment failure
//          before the actual Compose call (driven.ErrDockerUnavailable)
//   - 12 — fachlicher Ausführungsfehler (M6): Compose-runtime error
//          after passing preflight (driven.ErrComposeRuntime) or
//          stabilization timeout (driving.ErrStabilizationTimeout)
//   - 14 — technischer Persistenz-/Dateisystemfehler: LH-FA-INIT-005
//          backup-suffix exhausted (ErrBackupSuffixExhausted),
//          backup source vanished mid-flight
//          (ErrBackupSourceMissing)
//   - 1  — everything else (generic error)
//
// The mapping lives in the driving adapter because exit-code
// semantics are part of the CLI contract (LH-FA-CLI-006), not of
// the application use-cases — the application layer returns
// sentinel errors and lets the adapter translate.
//
// Sentinel-Reihenfolge in der Klassifikation (slice plan §T7
// §Sentinel-Schichtung): Driven-Sentinels werden ZUERST geprüft
// (driven.ErrDockerUnavailable / driven.ErrComposeRuntime), erst
// danach die Driving-/Application-Sentinels. Sentinels überschneiden
// sich nicht; die Reihenfolge folgt der Schicht-Hierarchie.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	// Driven-port sentinels first (M6 slice §Sentinel-Schichtung).
	if errors.Is(err, driven.ErrDockerUnavailable) {
		return 11
	}
	if errors.Is(err, driven.ErrComposeRuntime) {
		return 12
	}
	if errors.Is(err, driving.ErrStabilizationTimeout) {
		return 12
	}
	// Driving / application sentinels.
	if isValidationError(err) {
		return 10
	}
	if errors.Is(err, ErrDoctorFailures) {
		return 11
	}
	if isFilesystemError(err) {
		return 14
	}
	if isUsageError(err) {
		return 2
	}
	return 1
}

// isValidationError returns true for the LH-FA-CLI-006 code-10
// sentinels currently known to u-boot. Add new sentinels here as
// later slices introduce them; the [ExitCode] doc-comment is the
// authoritative list.
func isValidationError(err error) bool {
	return errors.Is(err, driving.ErrProjectExists) ||
		errors.Is(err, driving.ErrFileExists) ||
		errors.Is(err, driving.ErrBaseDirMissing) ||
		errors.Is(err, driving.ErrBackupUnsupportedKind) ||
		errors.Is(err, driving.ErrForceRequiresBackup) ||
		errors.Is(err, driving.ErrProjectNotInitialized) ||
		errors.Is(err, driving.ErrServiceUnsupported) ||
		errors.Is(err, driving.ErrServiceInconsistent) ||
		errors.Is(err, driving.ErrComposeFileMissing) ||
		errors.Is(err, driving.ErrConfirmationRequired) ||
		errors.Is(err, domain.ErrInvalidProjectName) ||
		errors.Is(err, domain.ErrInvalidServiceName)
}

// isFilesystemError returns true for the LH-FA-CLI-006 code-14
// sentinels — technical persistence / filesystem failures the
// application cannot recover from. The user must intervene
// (clean up stale backups, free disk, etc.).
func isFilesystemError(err error) bool {
	return errors.Is(err, driving.ErrBackupSuffixExhausted) ||
		errors.Is(err, driving.ErrBackupSourceMissing)
}

// isUsageError detects two distinct classes of usage-level errors:
//
//   (a) u-boot-defined CLI sentinels — currently
//       [ErrConflictingModeFlags]. New sentinels in this class
//       belong in the errors.Is block at the top.
//   (b) Cobra-raised errors for malformed CLI input. Cobra does
//       not export sentinels for these; we string-match the
//       message prefix because that is the only stable handle we
//       have.
//
// The two classes coexist on purpose — splitting into two helpers
// would obscure the shared "return code 2" intent. Add to the
// right block based on whether the error has a Go sentinel or
// only a message prefix.
//
// Pinned against github.com/spf13/cobra v1.10.2 (see go.mod). A
// major Cobra upgrade must verify these prefixes still match the
// strings Cobra emits — the integration tests
// TestExecute_UnknownCommand / TestExecute_UnknownFlag /
// TestExecute_InitTooManyArgs exercise the real Cobra path and
// will fail loudly if the wording changes.
func isUsageError(err error) bool {
	if err == nil {
		return false
	}
	// (a) u-boot CLI sentinels.
	if errors.Is(err, ErrConflictingModeFlags) || errors.Is(err, ErrInvalidTimeout) {
		return true
	}
	// (b) Cobra usage-error string prefixes.
	msg := err.Error()
	prefixes := []string{
		"unknown command",
		"unknown flag",
		"flag needs an argument",
		"invalid argument",
		"requires at",
		"accepts ",
	}
	for _, p := range prefixes {
		if strings.HasPrefix(msg, p) {
			return true
		}
	}
	return false
}
