package driven

import "errors"

// ErrYAMLPathInvalid signals that a [YAMLCodec.PatchScalar] call
// targeted a path that cannot exist in a YAML mapping tree — e.g. an
// intermediate node is a scalar or sequence (cannot become a mapping
// without destructive replacement), the path is empty, or value is
// not a supported scalar kind.
var ErrYAMLPathInvalid = errors.New("yaml path invalid")

// ErrYAMLFragmentInvalid signals that a
// [YAMLCodec.PatchMappingEntryYAML] call received a fragment or
// parent-shape that the adapter refuses to splice — non-mapping
// valueYAML (scalar/sequence root), a top-level key duplicate in
// content, a parent value that is neither a mapping nor an empty
// inline `{}` / `null`, or a non-empty flow-style parent map.
//
// The adapter rejects rather than rewriting because a destructive
// re-layout would violate the byte-preservation contract of
// PatchMappingEntryYAML.
var ErrYAMLFragmentInvalid = errors.New("yaml fragment invalid")

// ErrYAMLParse signals that YAML *content* failed to parse. It is
// returned by every [YAMLCodec] method that parses content bytes:
// [YAMLCodec.Unmarshal], [YAMLCodec.PatchScalar],
// [YAMLCodec.PatchMappingEntryYAML] (content side), and
// [YAMLCodec.LocateMarkedEntry]. Application callers branch on it
// via [errors.Is] to surface a "user must fix the YAML" repair hint
// — typically mapped to LH-FA-CLI-006 exit code 10 — without ever
// importing a YAML library themselves (depguard
// `application-no-yaml`).
//
// Scope and what is NOT covered:
//
//   - The sentinel covers parse failures of *content* the caller
//     hands the codec. It does NOT cover parse failures of
//     [YAMLCodec.PatchMappingEntryYAML]'s `valueYAML` argument —
//     that path is structural fragment validation and stays on
//     [ErrYAMLFragmentInvalid] because the caller controls the
//     fragment, not the user.
//   - The sentinel only signals "parse failed"; finer-grained
//     classification (TypeError vs SyntaxError vs unknown-key) is
//     intentionally out of scope; a future slice can layer it on
//     top.
//
// Introduced by [`slice-v1-yaml-parse-error-sentinel.md`] to close
// the M7-T5-N2 classification gap (`u-boot generate devcontainer`
// on a corrupt `compose.yaml` would otherwise surface as
// `ErrGenerateFileSystem` → exit 14 instead of the spec-mandated
// exit 10).
var ErrYAMLParse = errors.New("yaml parse error")

// ErrYAMLAnchorMismatch signals that
// [YAMLCodec.PatchMappingEntryYAML] would have inserted or replaced
// a managed block at an anchor that disagrees with the surrounding
// YAML mapping shape. Two symmetric cases:
//
//  1. A block with the requested markerName exists in content, but
//     outside the parentKey.entryKey byte range (marker is anchored
//     elsewhere — for example the postgres marker landed under
//     volumes:).
//  2. A mapping entry parentKey.entryKey exists in content but
//     contains no managed block with the requested markerName
//     (a manual user-edited entry under the requested key).
//
// PatchMappingEntryYAML is the only method that returns this sentinel
// — [LocateMarkedEntry] expresses both cases through the
// [LocateResult] flags instead, so application callers can branch
// without catching a driven-port error.
var ErrYAMLAnchorMismatch = errors.New("yaml anchor mismatch")

// LocateResult is the read-only outcome of
// [YAMLCodec.LocateMarkedEntry]. The four orthogonal flags let the
// application layer branch on the symmetric anchor states without
// importing a YAML library or catching a driven-port sentinel for
// the well-formed wrong-anchor cases.
type LocateResult struct {
	// EntryExists reports whether parentKey.entryKey exists as a
	// mapping entry in content (regardless of whether the entry
	// holds a managed block or user-manual content).
	EntryExists bool

	// MarkerInEntry reports whether a managed block with the
	// requested markerName exists AND encloses the parentKey.entryKey
	// mapping entry (the regular u-boot-managed state).
	MarkerInEntry bool

	// MarkerSomewhereElse reports whether a managed block with the
	// requested markerName exists but does NOT enclose the
	// parentKey.entryKey entry — either anchored under a different
	// parent or sibling-to but not surrounding the entry.
	MarkerSomewhereElse bool

	// BlockBody is the byte range between the BEGIN and END marker
	// lines (exclusive of the marker lines themselves) when
	// MarkerInEntry is true. Nil otherwise. Suitable for adapter-
	// free content inspection (e.g. the LH-FA-ADD-002 stale-content
	// check in M5-T4c).
	BlockBody []byte
}

// YAMLCodec abstracts YAML serialization and structural patching.
// The application layer uses it to marshal the u-boot.yaml structure
// (LH-FA-CONF-001..003), to read it back for idempotency checks, to
// set a single scalar value via [PatchScalar], and — via
// [PatchMappingEntryYAML] / [LocateMarkedEntry] — to manage
// structured Compose service/volume blocks under the LH-SA-FILE-002
// managed-block convention without re-marshaling the whole file.
type YAMLCodec interface {
	// Marshal serializes v into a YAML byte slice with the
	// project-wide indent and key-ordering conventions implemented by
	// the adapter.
	Marshal(v any) ([]byte, error)

	// Unmarshal parses YAML data into v.
	Unmarshal(data []byte, v any) error

	// PatchScalar sets a scalar value at the given mapping-path
	// inside content and returns the rewritten document.
	//
	// Path is a list of mapping keys interpreted top-down; missing
	// intermediate mappings are created. value must be a scalar
	// type the adapter can serialize (bool, string, int / int64,
	// float64) — passing a slice, map or struct returns
	// [ErrYAMLPathInvalid].
	//
	// Comments and sibling keys outside the path are preserved as
	// far as the underlying yaml.v3 node API preserves them
	// (head/foot/line comments stay on their nodes; key ordering of
	// the unmodified siblings is retained).
	//
	// Errors:
	//   - empty path ⇒ [ErrYAMLPathInvalid].
	//   - an intermediate node already exists as a non-mapping
	//     (scalar / sequence) ⇒ [ErrYAMLPathInvalid] wrapped with
	//     the offending key.
	//   - value is not a supported scalar kind ⇒
	//     [ErrYAMLPathInvalid] wrapped with the Go type name.
	//   - parse errors of content ⇒ wrapped, non-sentinel error.
	PatchScalar(content []byte, path []string, value any) ([]byte, error)

	// PatchMappingEntryYAML inserts or replaces a u-boot-managed
	// block (`# BEGIN/END U-BOOT MANAGED BLOCK: markerName`) under
	// parentKey.entryKey inside content. valueYAML is the YAML body
	// (a mapping root) that becomes the new value of the entryKey
	// mapping entry; the adapter wraps it in the BEGIN/END marker
	// lines and splices into content as a byte-range edit.
	//
	// parentKey is a top-level mapping key (e.g. "services",
	// "volumes"); nested paths are deliberately out of scope (see
	// the M5-T4b slice rationale).
	//
	// Byte-preservation contract: only three byte ranges may change:
	//   1. the existing managed block byte range, when present;
	//   2. an insertion range at the entryKey mapping position when
	//      the block is newly created;
	//   3. the parent header line (`parent: {}` → `parent:` +
	//      body) when the parent had an inline empty form, or a
	//      newly appended parent block at end-of-file when parent
	//      was missing entirely.
	// All other content (comments, sibling entries, other top-level
	// blocks) MUST remain byte-identical.
	//
	// Errors:
	//   - valueYAML is not a mapping root, content has a top-level
	//     key duplicate, or the parent is a scalar/sequence/non-empty
	//     flow mapping ⇒ [ErrYAMLFragmentInvalid].
	//   - managed-block anchor mismatch (marker elsewhere, or
	//     entry-without-marker) ⇒ [ErrYAMLAnchorMismatch].
	//   - malformed managed block (BEGIN without END / duplicate
	//     BEGIN) ⇒ wrapped malformed-block error.
	//   - parse errors of content or valueYAML ⇒ wrapped,
	//     non-sentinel error.
	PatchMappingEntryYAML(content []byte, parentKey, entryKey string, valueYAML []byte, markerName string) ([]byte, error)

	// LocateMarkedEntry inspects content without mutating it and
	// reports the symmetric anchor state for parentKey.entryKey
	// versus a managed block with markerName. The result is used
	// by the application layer for pre-patch anchor checks and
	// LH-FA-ADD-002 stale-content inspection.
	//
	// All well-formed inputs return (result, nil) — including
	// wrong-anchor states (MarkerSomewhereElse=true, or EntryExists
	// without MarkerInEntry). Sentinel errors only fire for true
	// adapter failures: content parse errors and malformed managed
	// blocks (BEGIN without END / duplicate BEGIN).
	LocateMarkedEntry(content []byte, parentKey, entryKey, markerName string) (LocateResult, error)
}
