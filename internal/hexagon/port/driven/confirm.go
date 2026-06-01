package driven

import "context"

// Confirmer is the application's side-channel for asking the user a
// single yes/no question during a use-case run. Introduced in the M4
// soft-existing-detection slice (LH-FA-INIT-004) to keep `os.Stdin`
// out of the application layer.
//
// The single method [Confirmer.ConfirmTreatAsExisting] is narrowly
// scoped to the soft-detection scenario: `u-boot init` found ≥3
// LH-FA-INIT-003 structure elements in BaseDir without a hard marker
// (`u-boot.yaml`/`compose.yaml`/`.env.example`) and must decide
// whether to treat the directory as an already-initialized project.
// Subsequent confirmation needs (e.g. M5 `add postgres` overwriting
// existing service blocks) should add their own narrowly-scoped
// methods on this interface rather than a generic `Confirm(prompt
// string)` — the explicit names keep call sites readable and let the
// adapter render context-specific prompts.
//
// Policy ownership: the application decides *whether* to call the
// Confirmer (see [LH-FA-INIT-004] non-interactive carve-outs);
// `--no-interactive` makes the service skip the call entirely. The
// adapter only owns the *how* of the prompt (text formatting, source
// of user input, `--yes` auto-confirmation).
type Confirmer interface {
	// ConfirmTreatAsExisting asks the user whether BaseDir should be
	// treated as an existing u-boot project. indicators is the
	// deterministic list of LH-FA-INIT-003 structure elements that
	// were found and triggered the soft-detection threshold (≥3); the
	// adapter is free to render the list in the prompt for context.
	//
	// Returns (true, nil) when the user confirms (or the adapter is
	// configured for auto-yes via `--yes`); (false, nil) when the
	// user declines (or refuses to answer / sends EOF — adapter's
	// choice).
	//
	// Returns a non-nil error only for I/O failures on the input
	// channel that the adapter cannot interpret as "user said no"; in
	// that case the application aborts the use case and surfaces the
	// error to the CLI.
	ConfirmTreatAsExisting(ctx context.Context, baseDir string, indicators []string) (bool, error)

	// ConfirmRemoveVolumes asks the user whether `u-boot down
	// --volumes` should proceed and drop named Compose volumes
	// alongside the containers. Surfaces the M6 LH-FA-CLI-005A
	// §254 destructive-confirmation prompt.
	//
	// Returns (true, nil) when the user confirms (the adapter
	// MUST default-show `N` capitalized, mirroring the
	// soft-detection prompt's safer-default convention); (false,
	// nil) when the user declines or sends EOF without input.
	//
	// Returns a non-nil error only for I/O failures on the input
	// channel that the adapter cannot interpret as "user said no";
	// in that case the application aborts the use case and
	// surfaces the error to the CLI.
	//
	// Caller responsibility (slice plan §T5 truth table): the
	// application service decides *whether* to invoke this method.
	// LH-FA-CLI-005A §254 requires that `--no-interactive` without
	// `--yes` skip the call entirely and return
	// `driving.ErrConfirmationRequired` directly — the adapter is
	// NOT expected to surface that case.
	ConfirmRemoveVolumes(ctx context.Context, baseDir string) (bool, error)

	// ConfirmAddDependency asks the user whether `u-boot add <svc>`
	// should auto-install missing add-on dependencies declared via
	// [domain.AddOnDependency] (LH-FA-ADD-006). The prompt is the
	// interactive default path; `--with-deps` / `--yes` /
	// `--no-interactive` short-circuit it at the application layer.
	//
	// svc is the add-on the user invoked; missing is the ordered list
	// of service names that must be registered first. Both are passed
	// as strings to keep the port domain-free (same convention as
	// [ConfirmTreatAsExisting]).
	//
	// Returns (true, nil) when the user confirms (the adapter MUST
	// default-show `N` capitalized — installing an extra service is
	// non-trivial, so the safer-default is "no, abort"); (false, nil)
	// when the user declines or sends EOF without input.
	//
	// Returns a non-nil error only for I/O failures on the input
	// channel that the adapter cannot interpret as "user said no".
	//
	// Caller responsibility (slice-v1-addons-deps T3 four-mode
	// dispatch): the application service decides *whether* to
	// invoke this method. `--no-interactive` without `--yes` /
	// `--with-deps` skips the call entirely and returns
	// [driving.ErrDependenciesRequired] directly.
	ConfirmAddDependency(ctx context.Context, svc string, missing []string) (bool, error)
}
