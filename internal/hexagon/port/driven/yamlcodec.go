package driven

// YAMLCodec abstracts YAML serialization. The application layer uses
// it to marshal the u-boot.yaml structure (LH-FA-CONF-001..003) and
// to read it back for idempotency checks.
//
// Managed-block-aware edits (LH-SA-FILE-002) are a future concern and
// will extend this interface; the M3 surface is intentionally minimal.
type YAMLCodec interface {
	// Marshal serializes v into a YAML byte slice with the
	// project-wide indent and key-ordering conventions implemented by
	// the adapter.
	Marshal(v any) ([]byte, error)

	// Unmarshal parses YAML data into v.
	Unmarshal(data []byte, v any) error
}
