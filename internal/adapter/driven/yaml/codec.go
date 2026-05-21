// Package yaml is the gopkg.in/yaml.v3-backed implementation of the
// `port/driven.YAMLCodec` interface (LH-FA-ARCH-002).
//
// The package is intentionally a thin wrapper around yaml.v3 so the
// application layer can stay free of YAML-library imports
// (LH-FA-ARCH-003, depguard-enforced).
package yaml

import (
	"gopkg.in/yaml.v3"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// Codec is the production YAML codec. The zero value is usable; use
// [New] for readability and future extension (indent, line-width).
type Codec struct{}

// Static check: Codec satisfies the YAMLCodec port.
var _ driven.YAMLCodec = (*Codec)(nil)

// New returns a ready-to-use Codec.
func New() *Codec { return &Codec{} }

// Marshal delegates to yaml.v3 with the library defaults (block style,
// 4-space indent, key ordering as defined in the source struct tags).
func (Codec) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

// Unmarshal delegates to yaml.v3.
func (Codec) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}
