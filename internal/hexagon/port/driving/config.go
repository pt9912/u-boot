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

	// AllowExternalFeatureSources carries additional LH-FA-DEV-003
	// source URLs from the `--allow-external-feature-sources` flag
	// (Spec §714). Only meaningful when
	// Path.Kind == domain.ConfigDevcontainerFeatureSourcesAllow; the
	// use case merges these entries with the comma-separated
	// positional [Value] before validation + dedupe. Slice-v1-
	// devcontainer-features T4.
	AllowExternalFeatureSources []string

	// PreviewMode selects the FS-mutation regime for `config set`
	// per slice-v1-cli-json-dry-run-config T0-(g)/T2 — the same
	// three-state contract every modifying subcommand shares
	// (add/init/generate/remove). The CLI maps the --dry-run/--diff
	// flag combination to one of [PreviewNone]/[PreviewDryRun]/
	// [PreviewAndApply]; the Composition-Root's fsFactory routes the
	// `s.fs.WriteFile(u-boot.yaml, …)` call through the recording
	// adapter accordingly. Default zero value [PreviewNone] keeps
	// today's production-write behaviour, so the non-JSON path is
	// unchanged. T3 reads it in `ConfigService.Set`; T5 sets it
	// from the parsed flags.
	PreviewMode PreviewMode

	// SilenceLogger switches the application-layer logger to a no-op
	// sink for the duration of this request so the five
	// `s.logger.*` sites in the Set path (three Info, two Debug)
	// do not pollute stdout-bound machine output with stderr lines
	// (slice-v1-cli-json-dry-run-config T0-(n)). The CLI sets this
	// to `true` when `--json` is active. Pattern is symmetric to
	// [UpRequest.SilenceProgress] / [RemoveServiceRequest.
	// SilenceConfirmer]: a boolean request flag that lets the use
	// case branch internally without the CLI knowing about the
	// no-op sink. T3 wires the branch; T5 sets the field.
	SilenceLogger bool
}

// ConfigSetResponse is the output of [ConfigUseCase.Set]. The CLI
// surfaces the OldValue → NewValue transition in its summary line
// so the user sees what they replaced. Both are stringified in
// the same canonical form [ConfigGetResponse.Value] uses.
type ConfigSetResponse struct {
	Path     domain.ConfigPath
	OldValue string
	NewValue string

	// Warnings carries soft-warning diagnostics the use case wants
	// the CLI adapter to surface via the JSON envelope's
	// `diagnostics[]` array (slice-v1-cli-json-dry-run-config T0-(n)
	// Orphan-Feature-WARN-Migration; [WarningEntry] type inherited
	// from remove/up T2). Empty / nil on the happy path; populated
	// when `maybeWarnOrphanFeatureActivation` fires for a
	// `services.<svc>.enabled` write whose service is not
	// registered (LH-FA-DEV-003 / `level: "warn"`). T3 appends the
	// entries; T5 maps them to `diagnostics[]`.
	Warnings []WarningEntry

	// PlannedFiles carries the recorder-captured `u-boot.yaml`
	// mutation when the request ran under a non-PreviewNone
	// [PreviewMode] (slice-v1-cli-json-dry-run-config T4-Review
	// R-T4-1; Pattern-Erbe [AddServiceResponse.PlannedFiles]). It is
	// the data source the CLI's shared diff renderer
	// (`mapPlannedFilesToWire` / `writeDiff`) needs for `config set
	// --diff`: the [PlannedFile.NewContent] (patched bytes) and
	// [PlannedFile.OldContent] (current bytes) only exist in the
	// recorder, not in the scalar OldValue/NewValue fields.
	//
	// nil on three paths: PreviewNone (production write, no recorder —
	// the CLI doesn't render a plan), the NoOp short-circuit (no write
	// captured → `plannedFiles: []` per T0-(d)), and the legacy
	// [NewConfigService] constructor (nil factory). Always exactly
	// one entry (`u-boot.yaml`) when populated. T5 maps it to the
	// envelope's `plannedFiles[]` + `hunks[]`.
	PlannedFiles []PlannedFile
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
// `application` (LH-FA-ARCH-003 depguard rule). They map to
// LH-FA-CLI-006 exit codes:
//
//   - ErrConfigPathUnknown            → 10 (validation; unknown dotted path)
//   - ErrConfigValueInvalid           → 10 (validation; bad value coercion)
//   - ErrConfigWriteRejected          → 10 (validation; non-writable path,
//                                        `u-boot add <svc>` hint)
//   - ErrConfigPostPatchSanityFailed  → 10 (validation; post-patch
//                                        roundtrip mismatch)
//   - ErrConfigSchemaInvalid          → 10 (validation; post-patch schema
//                                        roundtrip failed)
//   - ErrConfigValueNotSet            → 10 (validation; optional path unset)
//   - ErrConfigFileSystem             → 14 (technical; IO/permission)
//
// ErrConfigWriteRejected and ErrConfigPostPatchSanityFailed are new
// in slice-v1-cli-json-dry-run-config T2: they split the three
// semantic classes that [ErrConfigValueInvalid] conflated today
// (T0-(m)) so JSON consumers can disambiguate by `code` instead of
// a message substring. T3 redirects the corresponding wrap-sites in
// `application/config.go`; until then nothing returns them.

// ErrConfigPathUnknown signals that the CLI received a `<path>`
// argument that does not match the [domain.ConfigPath] whitelist.
// The CLI wraps the [domain.ErrInvalidConfigPath] cause with this
// sentinel so [ExitCode] sees code 10 (LH-FA-CLI-006 validation)
// regardless of the Go error chain depth.
var ErrConfigPathUnknown = errors.New("config: unknown path")

// ErrConfigValueInvalid signals a value-coercion failure: the raw
// value cannot be coerced to the path's expected scalar type (e.g.
// `set devcontainer.enabled vielleicht`). Code 10. The CLI emits
// the wrapped detail verbatim so the user sees the specific reason.
//
// Historically this sentinel also carried the write-rejected-path
// and post-patch-sanity classes; slice-v1-cli-json-dry-run-config
// T0-(m) split those into [ErrConfigWriteRejected] and
// [ErrConfigPostPatchSanityFailed] so consumers can disambiguate by
// `code`. T3 performs the redirect in `application/config.go`.
var ErrConfigValueInvalid = errors.New("config: invalid value")

// ErrConfigWriteRejected signals that the requested path's
// [domain.ConfigPath.WriteAllowed] flag is false, so `config set`
// refuses the write (slice-m8-config.md §D1: only
// `services.<svc>.enabled` falls under this rejection today; the
// user is pointed at `u-boot add <svc>`). Code 10.
//
// Split out of [ErrConfigValueInvalid] in
// slice-v1-cli-json-dry-run-config T2 (T0-(m)). The message keeps
// the `u-boot add <svc>` hint so the human path is unchanged; the
// new sentinel only adds a distinct `code` for JSON consumers.
// T3 redirects the WriteAllowed-reject sites
// (`application/config.go` ~Z. 251-256) onto this sentinel.
var ErrConfigWriteRejected = errors.New("config: write rejected for non-writable path")

// ErrConfigPostPatchSanityFailed signals that the value round-trip
// after [PatchScalar] did not reproduce the expected value — a rare
// schema-drift indicator distinct from the structural / domain
// re-validation covered by [ErrConfigSchemaInvalid]. Code 10. The
// use case MUST NOT write the file in this case.
//
// Split out of [ErrConfigValueInvalid] in
// slice-v1-cli-json-dry-run-config T2 (T0-(m)) so consumers can
// tell post-patch sanity failures apart from value-coercion
// failures. T3 redirects the post-patch-sanity sites
// (`application/config.go` ~Z. 376, 388) onto this sentinel.
var ErrConfigPostPatchSanityFailed = errors.New("config: post-patch sanity check failed")

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
