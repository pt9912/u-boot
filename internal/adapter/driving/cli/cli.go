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
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// App holds the driving-port dependencies the CLI needs.
//
// The struct is intentionally small — one field per use-case port,
// plus environment hooks (getwd) that tests substitute via
// functional options.
type App struct {
	// version is the build-time version string, surfaced via
	// `u-boot --version`. The wiring layer passes it in; the CLI
	// package does not own version metadata.
	version string

	// initUseCase implements `u-boot init` (LH-FA-INIT-001..007).
	initUseCase driving.InitProjectUseCase

	// getwd is the working-directory probe; defaults to os.Getwd.
	// Tests inject a fake via [WithGetwd] so they do not depend on
	// the host pwd.
	getwd func() (string, error)
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
func New(version string, initUC driving.InitProjectUseCase, opts ...Option) *App {
	a := &App{
		version:     version,
		initUseCase: initUC,
		getwd:       os.Getwd,
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

// ExitCode classifies a CLI error into the u-boot exit-code scheme
// (LH-FA-CLI-006):
//
//   - 0  — no error
//   - 2  — pure CLI / flag errors (unknown subcommand, unknown flag,
//          missing required arg, too many positional args)
//   - 10 — fachlicher Validierungsfehler: LH-FA-INIT-004 marker
//          collisions (ErrProjectExists), LH-FA-INIT-006 invalid
//          project name (ErrInvalidProjectName), LH-AK-001 missing
//          BaseDir (ErrBaseDirMissing)
//   - 1  — everything else (generic error)
//
// The mapping lives in the driving adapter because exit-code
// semantics are part of the CLI contract (LH-FA-CLI-006), not of
// the application use-cases — the application layer returns
// sentinel errors and lets the adapter translate.
//
// Codes 11/12/13–15 (environment, runtime, technical errors) are
// added by later slices that introduce the corresponding use-case
// sentinels (`u-boot doctor` for 11, `u-boot up`/`down` for 12).
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	if isValidationError(err) {
		return 10
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
		errors.Is(err, driving.ErrBaseDirMissing) ||
		errors.Is(err, domain.ErrInvalidProjectName)
}

// isUsageError detects errors that Cobra raises for malformed CLI
// input. Cobra does not export a sentinel; we look at the message
// prefix because that is what reaches us.
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
