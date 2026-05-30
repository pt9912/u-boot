package driving

import (
	"context"
	"errors"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// ConfigGetRequest is the input for [ConfigUseCase.Get]. The CLI
// adapter parses the positional `<path>` argument through
// [domain.NewConfigPath] before constructing the request, so by
// the time the use case sees Path it is already whitelist-checked.
type ConfigGetRequest struct {
	// BaseDir is the absolute path of the initialized u-boot
	// project. Mandatory; the CLI adapter defaults it to the
	// current working directory.
	BaseDir string

	// Path is the validated config path to read.
	Path domain.ConfigPath
}

// ConfigGetResponse is the output of [ConfigUseCase.Get]. The
// stringified Value is the canonical scalar form (e.g. `true` /
// `false` for bool, the raw string for project.name); the CLI
// prints it bare with a single trailing newline.
type ConfigGetResponse struct {
	Path  domain.ConfigPath
	Value string
}

// ConfigSetRequest is the input for [ConfigUseCase.Set]. Path is
// whitelist-validated by the CLI before reaching the use case;
// the use case still consults [domain.ConfigPath.WriteAllowed]
// before writing, because the same Path values are reachable
// both ways and a write-rejected path produces
// [ErrConfigValueInvalid] with a `u-boot add`-flavoured hint
// (slice-m8-config.md §D1 review-followup M1).
type ConfigSetRequest struct {
	BaseDir string
	Path    domain.ConfigPath

	// Value is the raw string the user supplied. The use case
	// coerces it (e.g. via [strconv.ParseBool] for bool paths)
	// before patching the YAML; coercion failure produces
	// [ErrConfigValueInvalid].
	Value string
}

// ConfigSetResponse is the output of [ConfigUseCase.Set]. The CLI
// surfaces the OldValue → NewValue transition in its summary line
// so the user sees what they replaced. Both are stringified in
// the same canonical form [ConfigGetResponse.Value] uses.
type ConfigSetResponse struct {
	Path     domain.ConfigPath
	OldValue string
	NewValue string
}

// ConfigShowRequest is the input for [ConfigUseCase.Show]. No
// path argument — Show always returns the full `u-boot.yaml`
// body.
type ConfigShowRequest struct {
	BaseDir string
}

// ConfigShowResponse is the output of [ConfigUseCase.Show]. Body
// is the raw u-boot.yaml content byte-identical to disk
// (slice-m8-config.md §D5): no re-parse, comments preserved.
type ConfigShowResponse struct {
	Body []byte
}

// All Config sentinels below live in the `driving` package so the
// CLI adapter can branch via [errors.Is] without importing
// `application` (LH-FA-ARCH-003 depguard rule). All four map to
// LH-FA-CLI-006 exit codes:
//
//   - ErrConfigPathUnknown    → 10 (validation; unknown dotted path)
//   - ErrConfigValueInvalid   → 10 (validation; bad value coercion
//                                or write-rejected path)
//   - ErrConfigSchemaInvalid  → 10 (validation; post-patch schema
//                                roundtrip failed)
//   - ErrConfigFileSystem     → 14 (technical; IO/permission)

// ErrConfigPathUnknown signals that the CLI received a `<path>`
// argument that does not match the [domain.ConfigPath] whitelist.
// The CLI wraps the [domain.ErrInvalidConfigPath] cause with this
// sentinel so [ExitCode] sees code 10 (LH-FA-CLI-006 validation)
// regardless of the Go error chain depth.
var ErrConfigPathUnknown = errors.New("config: unknown path")

// ErrConfigValueInvalid signals one of two value-rejection
// shapes:
//
//  1. The raw value cannot be coerced to the path's expected
//     scalar type (e.g. `set devcontainer.enabled vielleicht`).
//  2. The path's [domain.ConfigPath.WriteAllowed] flag is false
//     and the use case rejects the Set call (M8 §D1: only
//     `services.<svc>.enabled` falls under this rejection today;
//     the user is pointed at `u-boot add <svc>`).
//
// Code 10. The CLI emits the wrapped detail verbatim so the user
// sees the specific reason.
var ErrConfigValueInvalid = errors.New("config: invalid value")

// ErrConfigSchemaInvalid signals that the post-patch YAML body
// fails the two-stage schema validation
// (slice-m8-config.md §D3):
//
//  1. yaml.v3 Unmarshal into ubootYAMLConfig succeeds (catches
//     structural damage from PatchScalar).
//  2. Per-path domain re-validators on the resulting struct
//     pass (e.g. domain.NewProjectName on cfg.Project.Name)
//     so a lenient yaml.v3 accepting "@@@" into a string field
//     is still caught.
//
// Either stage failing returns this sentinel. The use case
// MUST NOT write the file in either case
// (writesBefore == writesAfter). Code 10.
var ErrConfigSchemaInvalid = errors.New("config: schema validation failed")

// ErrConfigFileSystem wraps an unexpected IO/permissions error
// from [driven.FileSystem] (Read/Write/Stat). Code 14
// (technical persistence failure) via `isFilesystemError` in
// cli.go.
//
// Wrap (not direct errors.Is against a driven sentinel) because
// the driven layer exports no ErrFileSystem* sentinel today —
// same architectural decision as M7's ErrGenerateFileSystem.
var ErrConfigFileSystem = errors.New("config: filesystem error")

// ErrConfigValueNotSet signals that a Get call addressed an
// optional path whose backing field has never been written. Two
// concrete shapes today:
//
//   - `devcontainer.enabled` when the u-boot.yaml `devcontainer:`
//     block is missing or has no `enabled:` key.
//   - `services.<svc>.enabled` when the service is not registered
//     in `services:` (or registered without an `enabled:` key).
//
// The CLI surfaces it with a per-path hint pointing at the
// canonical write path (`u-boot init --devcontainer` /
// `u-boot config set ...` for devcontainer; `u-boot add <svc>`
// for services). Code 10 (LH-FA-CLI-006 validation —
// "user must do something").
//
// `project.name` is required by LH-FA-CONF-002 §1308, so a
// missing name surfaces as [ErrConfigSchemaInvalid] instead,
// not this sentinel: the schema is corrupt, not just unset.
var ErrConfigValueNotSet = errors.New("config: value not set")

// ConfigUseCase is the driving-port for `u-boot config get/set/
// show`. Three methods instead of one action-enum dispatcher
// because the request shapes differ structurally (Show has no
// path) and Cobra's natural one-Cmd-per-subcommand layout maps
// 1:1 to per-method handlers.
//
// Contract:
//
//   - Each method validates the project state first
//     (`<BaseDir>/u-boot.yaml` exists) and returns
//     [ErrProjectNotInitialized] if not. The check is shared
//     so all three methods produce identical sentinel-mapping
//     behaviour at the CLI.
//   - On failure the response is the zero value and the error
//     wraps one of the sentinels above.
//
// Idempotence: Get and Show never mutate the filesystem and are
// trivially idempotent. Set is idempotent only when the new
// value equals the old value (OldValue == NewValue in the
// response, no WriteFile call).
type ConfigUseCase interface {
	Get(ctx context.Context, req ConfigGetRequest) (ConfigGetResponse, error)
	Set(ctx context.Context, req ConfigSetRequest) (ConfigSetResponse, error)
	Show(ctx context.Context, req ConfigShowRequest) (ConfigShowResponse, error)
}
