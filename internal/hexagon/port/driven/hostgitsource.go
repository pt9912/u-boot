package driven

import "context"

// TokenSourceEnvironment and TokenSourceGh are the source names of the
// token resolution ([HostGitSource.GitToken]): the variable of the host
// environment first, then the gh session. They are display values the
// caller may print, never secrets.
const (
	// TokenSourceEnvironment is the source name of the environment variable
	// GIT_TOKEN (LH-FA-DEV-009), the host variable itself.
	TokenSourceEnvironment = "GIT_TOKEN"
	// TokenSourceGh is the source name of the gh session of the host
	// (`gh auth token`).
	TokenSourceGh = "gh auth token"
)

// HostGitSource reads the git token and the git identity of the host for the
// devcontainer commands (LH-FA-DEV-009, LH-FA-DEV-010). It shells out to the
// `gh` and `git` binaries of the host; it opens no configuration file itself
// and reads no other key than the ones the methods name.
//
// Layer rules (LH-FA-ARCH-002, LH-FA-ARCH-003): driven port; the production
// implementation lives in `internal/adapter/driven/git`.
//
// The port answers three questions, and only `up`/`shell` ask them: which
// token the host carries ([HostGitSource.GitToken], LH-FA-DEV-009), whether
// the git configuration of the project folder on the host carries an identity
// ([HostGitSource.ConfigValue], LH-FA-DEV-010) and no other. `doctor` reads
// project files only and stays there (LH-FA-DEV-009).
type HostGitSource interface {
	// GitToken resolves the git token of the host in the fixed order of the
	// sources: the variable `GIT_TOKEN` of the environment of the process
	// (the explicit choice of the user) first, then the token of the gh
	// session (`gh auth token`). It returns the value and the name of the
	// source that provided it ("GIT_TOKEN" or "gh auth token"). With no
	// source it returns empty values without an error; a failure of the
	// resolution of one source (gh missing, gh not logged in) skips that
	// source the same way. The error carries only a failed environment
	// read; the adapters read the environment of their own process.
	GitToken(ctx context.Context) (value string, source string, err error)

	// ConfigValue returns the effective value of the git configuration key
	// (for example <code>user.name</code>) in the given project folder, as
	// `git -C <folder> config <key>` reports it (local before global before
	// system, with the conditional includes applied). A missing key reads as
	// an empty value without an error; an empty value is an empty value.
	ConfigValue(ctx context.Context, dir, key string) (string, error)
}
