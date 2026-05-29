package driving

import (
	"context"
	"errors"
	"io"
)

// DownRequest is the input for [DownUseCase.Down]. It is the
// application-layer expression of `u-boot down` per LH-FA-UP-004;
// the CLI adapter (M6-T7) translates the `--volumes`, `--yes` and
// the persistent `--no-interactive` root flag into the three
// boolean fields below.
//
// LH-FA-CLI-005A truth table for the destructive `--volumes`
// confirmation path (M6 slice §T5):
//
//	RemoveVolumes | AssumeYes | NonInteractive | behaviour
//	false         | *         | *              | proceed, no confirmer call
//	true          | true      | *              | proceed, no confirmer call
//	true          | false     | true           | fail-fast with ErrConfirmationRequired,
//	              |           |                | no confirmer call, no engine call
//	true          | false     | false          | call confirmer; (true, nil) → proceed,
//	              |           |                | (false, nil) → ErrConfirmationRequired,
//	              |           |                | error → wrap and return
//
// `NonInteractive` is a *separate* field from `AssumeYes`, not
// derivable from it: the spec distinguishes "automatically agree"
// (`--yes`) from "refuse to ask" (`--no-interactive`). The CLI
// adapter sets both fields verbatim from their respective flags.
// The orthogonal §235 rule (`--yes` AND `--no-interactive` set
// together → exit code 2) is enforced earlier in the root resolver
// and never reaches this use case.
type DownRequest struct {
	// BaseDir is the absolute path of the initialized u-boot project.
	// Mandatory; the CLI adapter defaults it to the current working
	// directory.
	BaseDir string

	// RemoveVolumes mirrors the `--volumes` CLI flag. When true,
	// `ComposeDown` is invoked with the equivalent of
	// `docker compose down -v`, which deletes named volumes
	// alongside containers. LH-FA-UP-004 §1015 isolates this from
	// the regular stop path so a non-destructive `down` cannot
	// accidentally drop persisted data.
	RemoveVolumes bool

	// AssumeYes mirrors the `--yes` CLI flag. When true (combined
	// with RemoveVolumes), the destructive confirmation step is
	// skipped — LH-FA-CLI-005A explicit auto-approve path.
	AssumeYes bool

	// NonInteractive mirrors the `--no-interactive` persistent
	// root flag. When true (combined with RemoveVolumes and
	// !AssumeYes), the use case returns [ErrConfirmationRequired]
	// without calling the confirmer at all — LH-FA-CLI-005A §254
	// "im nicht-interaktiven Modus ohne `--yes` ist der Befehl mit
	// Exit-Code 10 abzubrechen". Modeled as an explicit request
	// field so the application service does not need to know the
	// confirmer adapter's interactivity mode.
	NonInteractive bool

	// ProgressSink is the writer the application passes to the
	// `DockerEngine.ComposeDown` adapter for the Compose stderr
	// stream (per LH-NFA-PERF-002). The CLI adapter wires this to
	// `os.Stderr`; `nil` is treated as `io.Discard`.
	ProgressSink io.Writer
}

// DownResponse is the output of [DownUseCase.Down]. The CLI adapter
// renders a one-line success message keyed off [RemovedVolumes].
//
// No stop / removed counters — `docker compose down` emits a
// human-readable progress stream rather than a structured count,
// and inventing an "unknown" sentinel value would force every
// caller to special-case it. If a future slice needs precise
// counts (e.g. for `--json` output, LH-NFA-USE-004 V1), it would
// add a `ComposePs` diff before/after the call rather than parse
// the stderr stream.
type DownResponse struct {
	// RemovedVolumes echoes [DownRequest.RemoveVolumes] on success.
	// The CLI uses it to choose between
	// "environment stopped" and "environment stopped, volumes removed".
	RemovedVolumes bool
}

// ErrConfirmationRequired signals the destructive-confirmation
// abort path from LH-FA-CLI-005A §254 — `u-boot down --volumes` in
// non-interactive mode without `--yes`, or with an interactive
// confirmer that returned `(false, nil)`. Maps to LH-FA-CLI-006
// exit code 10 (fachliche Validierung).
//
// Distinct from the §235 root-level exclusivity error (`--yes` AND
// `--no-interactive` set simultaneously → exit code 2): the
// exclusivity error is a CLI-validation failure that never reaches
// the use case, while this sentinel signals a use-case-level
// refusal triggered by a real flag combination. Keeping the two
// failure modes on separate exit codes lets CI dashboards
// distinguish "user typoed flags" from "user tried to destroy data
// without approval".
var ErrConfirmationRequired = errors.New("confirmation required")

// DownUseCase is the driving-port for `u-boot down`. The CLI
// adapter holds a reference and calls [Down] from the Cobra command
// handler.
//
// Contract:
//
//   - On success the response carries [RemovedVolumes] = req.RemoveVolumes
//     and the error is nil.
//   - On a use-case failure the response is the zero value and the
//     error wraps one of: [ErrConfirmationRequired],
//     [ErrComposeFileMissing] (from `up.go`),
//     [ErrProjectNotInitialized] (from `addservice.go`), or a
//     `driven.*` sentinel forwarded from the engine.
//   - The returned error carries any Engine sentinel identity
//     intact via `fmt.Errorf("down service: ...: %w", err)`, so
//     `errors.Is(err, driven.ErrDockerUnavailable)` and
//     `errors.Is(err, driven.ErrComposeRuntime)` continue to hold
//     at the CLI level.
type DownUseCase interface {
	Down(ctx context.Context, req DownRequest) (DownResponse, error)
}
