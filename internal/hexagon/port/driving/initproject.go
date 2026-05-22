// Package driving holds the driving-port interfaces of u-boot — the
// use cases that the outside world (CLI, future HTTP daemon) calls
// into. Concrete implementations live in
// `internal/hexagon/application`.
//
// Layer rules (LH-FA-ARCH-002, LH-FA-ARCH-003): driving ports may
// only depend on `internal/hexagon/domain` and the Go standard
// library.
package driving

import (
	"context"
	"errors"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// InitProjectRequest is the input for [InitProjectUseCase.Init]. It is
// the application-layer expression of `u-boot init` flags; the CLI
// adapter (M3-T3) translates Cobra flags into this struct.
type InitProjectRequest struct {
	// Name is the explicit project name; when empty, the service
	// derives it from BaseDir's basename via
	// [domain.NormalizeProjectName] (LH-FA-INIT-002).
	Name string

	// BaseDir is the absolute path of the directory the project is
	// initialized in. Mandatory; the CLI adapter defaults it to the
	// current working directory.
	BaseDir string

	// SkipGit disables the git-init step (LH-FA-INIT-007). Default
	// (false) keeps git init enabled; the CLI adapter sets this from
	// the `--no-git` flag.
	SkipGit bool
}

// InitProjectResponse is the output of [InitProjectUseCase.Init].
type InitProjectResponse struct {
	// Project is the validated, ready-to-persist domain aggregate.
	Project domain.Project

	// Created lists the paths (relative to BaseDir) that were
	// created or written, in deterministic order, for the CLI
	// adapter to report back to the user.
	Created []string
}

// ErrProjectExists signals that BaseDir already looks like an
// initialized u-boot project (`u-boot.yaml`, `compose.yaml`, or
// `.env.example` present). LH-FA-INIT-004 forbids silent overwrite;
// the M3-T4 slice adds `--backup`/`--force` handling.
var ErrProjectExists = errors.New("project already initialized")

// ErrBaseDirMissing signals that req.BaseDir does not exist on the
// filesystem. The acceptance flow LH-AK-001 has the user create the
// directory (`mkdir demo && cd demo`); the use-case refuses to invent
// it because a typoed BaseDir would otherwise quietly initialize an
// unintended path under the typo.
//
// Sentinel lives in the driving port (not in the application package)
// so the CLI adapter can map it to its LH-FA-CLI-006 exit code
// without violating the depguard rule that forbids adapter→application
// imports (LH-FA-ARCH-003).
var ErrBaseDirMissing = errors.New("base directory does not exist")

// ErrBackupSourceMissing signals that the path passed to the backup
// strategy (LH-FA-INIT-005) no longer exists at the moment the
// backup runs. This is a race condition: the caller observed the
// file moments earlier, but the filesystem changed before the
// backup could capture it. Surfaces to the CLI as a technical
// filesystem error.
var ErrBackupSourceMissing = errors.New("backup source does not exist")

// ErrBackupSuffixExhausted signals that the LH-FA-INIT-005 §607
// backup-suffix space (<src>.bak through <src>.bak.999) is fully
// occupied, including after race-retries. A user hitting this has
// accumulated unusually many stale backups and must clean up
// manually. Maps to a technical filesystem exit code.
var ErrBackupSuffixExhausted = errors.New("backup suffix exhausted")

// ErrBackupUnsupportedKind signals that the backup target is neither
// a regular file nor a regular directory (currently only symlinks
// trip this). LH-FA-INIT-005 §608 does not specify symlink
// semantics; rejecting is the safe default until a follow-up slice
// decides between "copy-as-symlink" and "follow-then-copy". Maps to
// a validation exit code because the user gave the tool an input it
// cannot safely act on.
var ErrBackupUnsupportedKind = errors.New("backup source kind unsupported")

// ErrBackupTooLarge signals that a file in the backup scope exceeds
// the MVP size cap. The cap exists because the current FileSystem
// port loads files via ReadFile (full content into memory) and a
// multi-GB asset would OOM the process. Lifted by
// `slice-v1-backup-streaming-copy.md` once a streaming copy
// primitive lands. Maps to a technical filesystem exit code.
var ErrBackupTooLarge = errors.New("backup source exceeds size cap")

// InitProjectUseCase is the driving-port for `u-boot init`. The CLI
// adapter holds a reference and calls [Init] from the Cobra command
// handler.
type InitProjectUseCase interface {
	// Init initializes a new u-boot project in req.BaseDir according
	// to LH-FA-INIT-001..007 and LH-FA-CONF-001..003.
	//
	// Returns wrapped [ErrProjectExists] when BaseDir already
	// contains a project steering file and SkipGit-aware handling is
	// not enough to proceed (LH-FA-INIT-004). Returns wrapped
	// [domain.ErrInvalidProjectName] when the name (explicit or
	// derived) does not pass LH-FA-INIT-006.
	Init(ctx context.Context, req InitProjectRequest) (InitProjectResponse, error)
}
