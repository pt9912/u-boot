package driven

import "errors"

// ErrYAMLPathInvalid signals that a [YAMLCodec.PatchScalar] call
// targeted a path that cannot exist in a YAML mapping tree — e.g. an
// intermediate node is a scalar or sequence (cannot become a mapping
// without destructive replacement), the path is empty, or value is
// not a supported scalar kind.
var ErrYAMLPathInvalid = errors.New("yaml path invalid")

// YAMLCodec abstracts YAML serialization. The application layer uses
// it to marshal the u-boot.yaml structure (LH-FA-CONF-001..003), to
// read it back for idempotency checks, and — via [PatchScalar] — to
// set a single scalar value inside an existing document without
// destroying comments and unknown fields (the LH-SA-FILE-002 "managed
// edits" surface the M5 add-flow needs for u-boot.yaml).
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
}
