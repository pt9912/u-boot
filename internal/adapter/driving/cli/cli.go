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
	"log/slog"
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

	// generateUseCase implements `u-boot generate <artifact>`
	// (LH-FA-GEN-001..005).
	generateUseCase driving.GenerateUseCase

	// configUseCase implements `u-boot config get|set|show`
	// (LH-FA-CONF-001..005).
	configUseCase driving.ConfigUseCase

	// templateListUseCase implements `u-boot template list`
	// (LH-FA-TPL-004; slice-v1-template-list).
	templateListUseCase driving.TemplateListUseCase

	// removeServiceUseCase implements `u-boot remove <service>`
	// (LH-FA-ADD-007; slice-v1-add-remove).
	removeServiceUseCase driving.RemoveServiceUseCase

	// logsUseCase implements `u-boot logs [service]`
	// (LH-FA-UP-005; slice-v1-logs).
	logsUseCase driving.LogsUseCase

	// getwd is the working-directory probe; defaults to os.Getwd.
	// Tests inject a fake via [WithGetwd] so they do not depend on
	// the host pwd.
	getwd func() (string, error)

	// yes and noInteractive are bound to the root command's
	// PersistentFlags by [buildRootCommand].
	yes           bool
	noInteractive bool

	// quiet, verbose, debug are bound to the LH-FA-CLI-005 root
	// PersistentFlags. --quiet additionally filters SeverityOK
	// items from the doctor render and suppresses the up status
	// table / down success message. --verbose and --debug both
	// raise the logger level to Debug via buildRootCommand's
	// PersistentPreRunE; --quiet lowers it to Warn.
	quiet   bool
	verbose bool
	debug   bool

	// json is bound to the LH-NFA-USE-004 root --json PersistentFlag
	// (slice-v1-cli-json-dry-run-doctor T3). Subcommands that
	// implement the JSON envelope read this state; non-migrated
	// subcommands are rejected by the root PersistentPreRunE with
	// ErrJSONNotImplemented (exit code 2).
	json bool

	// logLevel is the slog level handle the wiring layer also
	// shares with the logger adapter. PersistentPreRunE flips it
	// based on the verbosity flags. Optional — nil-tolerant so
	// tests that do not need level-switching can omit
	// [WithLogLevel].
	logLevel *slog.LevelVar
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

// WithLogLevel hands the App a [*slog.LevelVar] that
// [buildRootCommand]'s PersistentPreRunE mutates from the
// LH-FA-CLI-005 verbosity flags. The wiring layer creates the
// LevelVar at the same point it constructs the logger adapter and
// passes the same instance to both — so the level change applies
// across every Logger.Debug/Info/Warn call. nil is allowed (the
// switching becomes a no-op); tests that do not exercise the level
// path can omit this option.
func WithLogLevel(level *slog.LevelVar) Option {
	return func(a *App) { a.logLevel = level }
}

// New constructs an App. The version string and every use-case
// implementation must be non-nil at call time; the CLI package
// trusts the wiring layer to honor that.
func New(version string, initUC driving.InitProjectUseCase, doctorUC driving.DoctorUseCase, addUC driving.AddServiceUseCase, upUC driving.UpUseCase, downUC driving.DownUseCase, genUC driving.GenerateUseCase, cfgUC driving.ConfigUseCase, tmplUC driving.TemplateListUseCase, removeUC driving.RemoveServiceUseCase, logsUC driving.LogsUseCase, opts ...Option) *App {
	a := &App{
		version:              version,
		initUseCase:          initUC,
		doctorUseCase:        doctorUC,
		addServiceUseCase:    addUC,
		upUseCase:            upUC,
		downUseCase:          downUC,
		generateUseCase:      genUC,
		configUseCase:        cfgUC,
		templateListUseCase:  tmplUC,
		removeServiceUseCase: removeUC,
		logsUseCase:          logsUC,
		getwd:                os.Getwd,
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

// ErrServiceNameMissing is returned by the `u-boot remove`
// custom-Args-validator when the positional `<service>` argument is
// absent (slice-v1-cli-json-dry-run-remove T0-(e) + R11-HIGH-F1 +
// R12-HIGH-F1). Lives in the cli package (not in `driving`)
// because the application service never sees this state — the
// missing-arg condition is a pure CLI form-validation concern
// (analog [ErrConflictingModeFlags]).
//
// Symmetrie-Bruch-Fix vs. `cobra.ExactArgs(1)`-Vorzustand: das
// vorherige `cobra.ExactArgs(1)`-Guard feuerte VOR RunE und ließ
// im --json-Pfad einen Konsumenten ohne Envelope auf stdout
// zurück (Spec §1841 verletzt — eine JSON-Mode-Invocation MUSS
// einen Envelope produzieren). Der neue Custom-`Args`-Validator
// (`validateRemoveArgs`) emittiert den Envelope BEVOR er den
// Sentinel an Cobra returnt; ExitCode-Mapping bleibt 2
// (LH-FA-CLI-006).
var ErrServiceNameMissing = errors.New("service name is required")

// ErrInvalidTimeout is returned by the M6 up subcommand when
// `--timeout` is a negative integer (LH-FA-UP-001 §965). The CLI
// could not delegate that validation to the application service —
// the application takes a `time.Duration` and could not distinguish
// a deliberate negative value from a unit-mistake-mismatch — so the
// rejection happens at the CLI before construction of the request.
// Maps to LH-FA-CLI-006 exit code 2.
var ErrInvalidTimeout = errors.New("--timeout must be >= 0")

// ErrJSONNotImplemented is returned by the root PersistentPreRunE
// when --json is passed to a subcommand form that has not yet been
// migrated to the LH-NFA-USE-004 envelope (slice-v1-cli-json-dry-run
// cluster T0-(g)). Maps to LH-FA-CLI-006 exit code 2. The error
// message itself is built by [jsonRejectError] and includes the
// concrete CommandPath + follow-up slice reference.
//
// Cluster-T_close-Pflicht-Check: every spec-enum subcommand form
// must end up in the allowlist (or the allowlist mechanic is removed
// completely). See docs/user/cli-json-output.md §6.1.
var ErrJSONNotImplemented = errors.New("json output not implemented for this subcommand")

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
//          LH-FA-ADD-007 service-not-registered (ErrServiceUnregistered),
//          LH-FA-ADD-006 add-on dependencies missing (ErrDependenciesRequired),
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
//   - 14 — technischer Persistenz-/Dateisystemfehler: LH-NFA-REL-003
//          backup-suffix exhausted (ErrBackupSuffixExhausted),
//          backup source vanished mid-flight
//          (ErrBackupSourceMissing); LH-FA-TPL-004 catalog adapter
//          failure (ErrTemplateCatalog — filesystem IO / malformed
//          embedded template.yaml); LH-FA-TPL-001 render-loop
//          failure (ErrTemplateRender — text/template parse/exec
//          or IO during the per-file render copy).
//          Footnote (slice-v1-cli-json-dry-run-init T0-(f)): die
//          Backup-Sentinels werden hier auf LH-NFA-REL-003 gezogen,
//          obwohl Spec §595-619 (INIT-005 "Überschreibschutz") sie
//          ursprünglich der INIT-005-Klasse zuordnete — Engineering-
//          Entscheidung im init-Slice, um Envelope-Code und
//          Exit-Code-Klasse (technische Persistenz) zu synchronisieren.
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
		errors.Is(err, driving.ErrComposeFileMissing) ||
		errors.Is(err, driving.ErrConfirmationRequired) ||
		// slice-v1-cli-json-dry-run-remove T0-(e) R2-HIGH-F1:
		// ErrConfirmerUnavailable (Confirmer-I/O-Failure beim --purge-
		// Gate, z.B. stdin EOF / pipe break) teilt die LH-FA-CLI-005A-
		// Klasse mit ErrConfirmationRequired und mappt auf exit 10.
		// Distinct vom User-Refusal-Pfad — beide sind Gate-Failures
		// vom selben Spec-Anker (§254).
		errors.Is(err, driving.ErrConfirmerUnavailable) ||
		errors.Is(err, driving.ErrGenerateManualConflict) ||
		isServiceValidationError(err) ||
		isConfigValidationError(err) ||
		isTemplateInitValidationError(err) ||
		errors.Is(err, domain.ErrInvalidProjectName) ||
		errors.Is(err, domain.ErrInvalidServiceName)
}

// isServiceValidationError bundles the add-on lifecycle sentinels
// (LH-FA-ADD-002/-005/-007) into one helper so [isValidationError]
// stays under the gocyclo threshold. All three share the LH-FA-
// CLI-006 code-10 mapping ("user must fix the add / remove
// invocation"); they are conceptually one cluster.
func isServiceValidationError(err error) bool {
	return errors.Is(err, driving.ErrServiceUnsupported) ||
		errors.Is(err, driving.ErrServiceInconsistent) ||
		errors.Is(err, driving.ErrServiceUnregistered) ||
		errors.Is(err, driving.ErrDependenciesRequired)
}

// isTemplateInitValidationError carves the slice-v1-template-init
// validation sentinels out of [isValidationError] so the parent
// helper stays under the gocyclo threshold. The two sentinels share
// the LH-FA-CLI-006 code-10 mapping but are conceptually one cluster
// ("the user must fix the --template invocation or the template
// content").
func isTemplateInitValidationError(err error) bool {
	return errors.Is(err, driving.ErrTemplateNotFound) ||
		errors.Is(err, driving.ErrInvalidTemplatePath)
}

// isConfigValidationError is the M8-T5 carve-out from
// [isValidationError] so the latter stays under the gocyclo
// threshold. The four config sentinels share the LH-FA-CLI-006
// code-10 mapping but are conceptually one cluster ("user must
// fix the config call"); keeping them in their own helper makes
// the partition obvious.
func isConfigValidationError(err error) bool {
	return errors.Is(err, driving.ErrConfigPathUnknown) ||
		errors.Is(err, driving.ErrConfigValueInvalid) ||
		// slice-v1-cli-json-dry-run-config T0-(m) split two classes
		// out of ErrConfigValueInvalid (T3). Both share the LH-FA-
		// CONF Exit-10 class — they MUST be listed here or ExitCode
		// (a classifier independent of the new T5 mapConfigError…
		// mapper) silently drops them to Exit 1. Independent-Review
		// finding R-IR-1: present regression on the plain CLI path.
		errors.Is(err, driving.ErrConfigWriteRejected) ||
		errors.Is(err, driving.ErrConfigPostPatchSanityFailed) ||
		errors.Is(err, driving.ErrConfigSchemaInvalid) ||
		errors.Is(err, driving.ErrConfigValueNotSet) ||
		// slice-v1-devcontainer-features Review-Followup R1:
		// LH-FA-DEV-003 source-format failures (raised by
		// `init --allow-external-feature-sources` without
		// `--devcontainer`, and by `generate devcontainer`'s
		// pre-write allowlist append) must map to exit-code 10
		// per Spec §720/§1353. The sentinel lives in
		// `internal/hexagon/domain` so this adapter file can
		// reference it without violating the
		// `adapter-no-application` depguard rule.
		errors.Is(err, domain.ErrInvalidFeatureSource)
}

// isFilesystemError returns true for the LH-FA-CLI-006 code-14
// sentinels — technical persistence / filesystem failures the
// application cannot recover from. The user must intervene
// (clean up stale backups, free disk, etc.).
//
// ErrAddFileSystem (slice-v1-cli-json-dry-run-add T0-(j)/T1-C) wraps
// raw os.WriteFile errors from addservice_execute.go so add-mid-write
// failures map to LH-NFA-REL-003 / exit-code 14, not the default
// "1" — without this entry the recorder would surface plannedFiles[]
// to the user but the process would exit with the wrong code class.
// ErrInitFileSystem (slice-v1-cli-json-dry-run-init T2) plays the
// same role for init's WriteFile/MkdirAll/BackupPath sites.
func isFilesystemError(err error) bool {
	return errors.Is(err, driving.ErrBackupSuffixExhausted) ||
		errors.Is(err, driving.ErrBackupSourceMissing) ||
		errors.Is(err, driving.ErrGenerateFileSystem) ||
		errors.Is(err, driving.ErrConfigFileSystem) ||
		errors.Is(err, driving.ErrTemplateCatalog) ||
		errors.Is(err, driving.ErrTemplateRender) ||
		errors.Is(err, driving.ErrAddFileSystem) ||
		errors.Is(err, driving.ErrInitFileSystem) ||
		// slice-v1-cli-json-dry-run-remove T5: ErrRemoveFileSystem
		// wraps raw os.WriteFile / RemoveAll / Exists / ReadFile /
		// Lstat errors in removeservice.go (8 FS-Wrap-Stellen) so
		// remove-mid-write failures map to LH-NFA-REL-003 / exit-code
		// 14, not the default "1" — without this entry the recorder
		// would surface plannedFiles[] to the user but the process
		// would exit with the wrong code class.
		errors.Is(err, driving.ErrRemoveFileSystem) ||
		// slice-v1-cli-json-dry-run-up-down T2/T3: ErrUpFileSystem /
		// ErrDownFileSystem wrap raw os.fs.Exists / ReadFile errors in
		// upservice.go (3 FS-Read-Wrap-Stellen) and downservice.go (2)
		// so up/down FS-read failures map to LH-NFA-REL-003 / exit-
		// code 14, not the default "1". Switch-Order in
		// mapUp/mapDownErrorToDiagnostic checks the FS-Sentinel FIRST
		// (T0-(e) R3-HIGH-1) so synthetic multi-`%w`-wraps with
		// FS+Docker fall to FS+exit-14, not Docker+exit-11.
		errors.Is(err, driving.ErrUpFileSystem) ||
		errors.Is(err, driving.ErrDownFileSystem) ||
		// slice-v1-cli-json-dry-run-logs T2/T3: ErrLogsFileSystem
		// wraps raw os.fs.Exists errors in logsservice.go (2 FS-Read-
		// Wrap-Stellen: checkProjectInitialized + checkComposeFile)
		// so logs FS-read failures map to LH-NFA-REL-003 / exit-code
		// 14, not the default "1". Switch-Order in
		// mapLogsErrorToDiagnostic checks the FS-Sentinel FIRST
		// (T0-(e) R3-HIGH-1) so synthetic multi-`%w`-wraps with
		// FS+Docker fall to FS+exit-14, not Docker+exit-11.
		errors.Is(err, driving.ErrLogsFileSystem)
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
	// (a) u-boot CLI sentinels. ErrArtifactUnknown is in this
	// category by spec mandate (§LH-FA-GEN-001): "Bei unbekanntem
	// Artefakt muss der Befehl mit Exit Code 2 abbrechen" — distinct
	// from `add <unknown-service>` which maps to code 10 via
	// [isValidationError].
	if errors.Is(err, ErrConflictingModeFlags) || errors.Is(err, ErrInvalidTimeout) ||
		errors.Is(err, ErrInvalidLogsTail) ||
		errors.Is(err, ErrFollowJSONNotSupported) ||
		errors.Is(err, ErrJSONNotImplemented) ||
		errors.Is(err, ErrServiceNameMissing) ||
		errors.Is(err, driving.ErrArtifactUnknown) ||
		errors.Is(err, driving.ErrTemplateConflictsWithFlag) {
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
