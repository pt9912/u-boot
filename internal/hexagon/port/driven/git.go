package driven

// Git abstracts the git operations u-boot uses. M3 needs only repo
// existence check and init (LH-FA-INIT-007). Add operations (commit,
// status, etc.) are added by later slices as needed.
type Git interface {
	// IsRepository reports whether the given directory is already
	// inside a git repository (i.e. `git rev-parse --is-inside-
	// work-tree` would succeed).
	IsRepository(dir string) (bool, error)

	// Init runs `git init` in the given directory. It must be a no-op
	// (or return a clear error) when the directory is already a repo
	// — the caller is responsible for the IsRepository pre-check
	// (LH-FA-INIT-007).
	Init(dir string) error
}
