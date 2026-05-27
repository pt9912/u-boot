package driven

import "context"

// Git abstracts the git operations u-boot uses. M3 needs only repo
// existence check and init (LH-FA-INIT-007). Add operations (commit,
// status, etc.) are added by later slices as needed.
//
// Methods take a [context.Context] because the underlying adapter
// shells out to the `git` binary — a process that can block on
// network (e.g. submodule fetch in future operations) or hang on a
// stale filesystem; the application layer must be able to cancel.
// [Clock], [FileSystem], and [YAMLCodec] do not take Context because
// their implementations are non-blocking syscalls / library calls.
type Git interface {
	// IsRepository reports whether the given directory is already
	// inside a git repository (i.e. `git rev-parse --is-inside-
	// work-tree` would succeed).
	IsRepository(ctx context.Context, dir string) (bool, error)

	// Init runs `git init` in the given directory. It must be a no-op
	// (or return a clear error) when the directory is already a repo
	// — the caller is responsible for the IsRepository pre-check
	// (LH-FA-INIT-007).
	Init(ctx context.Context, dir string) error

	// Version returns the installed git version string (e.g.
	// `"2.43.0"`). A non-nil error signals that the git binary is not
	// available — used by the M4 doctor's `git.installed` check to
	// distinguish "git missing" from "git present but with a quirk".
	// The CLI-level format of the underlying `git --version` output
	// (typically `"git version 2.43.0"`) is parsed by the adapter so
	// the application layer sees only the bare version number.
	Version(ctx context.Context) (string, error)
}
