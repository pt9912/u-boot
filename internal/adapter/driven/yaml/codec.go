// Package yaml is the gopkg.in/yaml.v3-backed implementation of the
// `port/driven.YAMLCodec` interface (LH-FA-ARCH-002).
//
// The package is intentionally a thin wrapper around yaml.v3 so the
// application layer can stay free of YAML-library imports
// (LH-FA-ARCH-003, depguard-enforced).
package yaml

import (
	"bytes"
	"fmt"
	"strconv"

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

// Unmarshal delegates to yaml.v3 and wraps any parse failure with
// [driven.ErrYAMLParse] so application callers can branch on
// `errors.Is(err, driven.ErrYAMLParse)` without importing yaml.v3.
// See [wrapYAMLParse] for the wrap convention.
func (Codec) Unmarshal(data []byte, v any) error {
	if err := yaml.Unmarshal(data, v); err != nil {
		return wrapYAMLParse("unmarshal", err)
	}
	return nil
}

// wrapYAMLParse converts a raw yaml.v3 error into a
// [driven.ErrYAMLParse]-wrapped error with a `<context>: <yaml-msg>`
// detail. The yaml.v3 error message often carries a leading
// `yaml: ` prefix (e.g. `yaml: line 3: did not find expected key`);
// [stripYAMLPrefix] removes it so a downstream `%v` does not
// surface a doubled prefix in user-facing messages.
//
// This is the only place the production adapter classifies parse
// failures; the four call sites (`Unmarshal`, `PatchScalar`,
// `LocateMarkedEntry` content parse, `assertNoTopLevelDuplicate`
// content parse) all route through here.
func wrapYAMLParse(context string, err error) error {
	return fmt.Errorf("%w: %s: %s", driven.ErrYAMLParse, context, stripYAMLPrefix(err))
}

// stripYAMLPrefix removes a leading `yaml: ` from a yaml.v3 error
// message. yaml.v3 emits parse errors as plain `error` values with
// the prefix baked into the message; the prefix is redundant after
// the [driven.ErrYAMLParse] wrap so we strip it for readability.
//
// Pinned by `TestCodec_StripYAMLPrefix` in codec_test.go (M1
// review-followup of the V1 sentinel slice).
func stripYAMLPrefix(err error) string {
	msg := err.Error()
	const prefix = "yaml: "
	if len(msg) >= len(prefix) && msg[:len(prefix)] == prefix {
		return msg[len(prefix):]
	}
	return msg
}

// PatchScalar implements the [driven.YAMLCodec.PatchScalar] contract:
// it parses content into a yaml.v3 Node tree, walks path through
// nested mappings (creating missing intermediate mappings), and sets
// a scalar at the leaf. The Node-API preserves comments and sibling
// keys; only the touched scalar / created mappings change.
//
// Empty content is treated as an empty mapping document — the patch
// then creates the full path. The output is re-marshaled with a
// 2-space indent (yaml.v3 default) so the result is deterministic
// across runs.
func (Codec) PatchScalar(content []byte, path []string, value any) ([]byte, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("%w: empty path", driven.ErrYAMLPathInvalid)
	}
	scalar, err := scalarNode(value)
	if err != nil {
		return nil, err
	}

	var doc yaml.Node
	if len(bytes.TrimSpace(content)) > 0 {
		if err := yaml.Unmarshal(content, &doc); err != nil {
			return nil, wrapYAMLParse("patch-scalar", err)
		}
	}

	root := documentRoot(&doc)
	if root.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%w: document root is %s, want mapping",
			driven.ErrYAMLPathInvalid, kindName(root.Kind))
	}

	if err := setMappingPath(root, path, scalar); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&doc); err != nil {
		return nil, fmt.Errorf("encode yaml: %w", err)
	}
	if err := enc.Close(); err != nil {
		return nil, fmt.Errorf("close yaml encoder: %w", err)
	}
	return buf.Bytes(), nil
}

// documentRoot returns the inner top-level node of doc. yaml.Unmarshal
// on an empty document leaves doc with Kind=0 and no children; we
// promote that to an empty MappingNode so the patch can populate it.
// For a parsed non-empty document yaml.v3 wraps the top level in a
// DocumentNode whose single child is the actual root mapping.
func documentRoot(doc *yaml.Node) *yaml.Node {
	if doc.Kind == 0 {
		doc.Kind = yaml.DocumentNode
		doc.Content = []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}}
		return doc.Content[0]
	}
	if doc.Kind == yaml.DocumentNode {
		if len(doc.Content) == 0 {
			doc.Content = []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}}
		}
		return doc.Content[0]
	}
	// Defensive: a bare non-document node — return as-is so the
	// caller's MappingNode check can reject non-mapping roots.
	return doc
}

// setMappingPath descends path through nested mappings inside root,
// creating missing intermediates, and assigns scalar to the leaf key.
// Returns [driven.ErrYAMLPathInvalid] when an intermediate key is
// present but its value is not a mapping (cannot become one without
// destructive replacement).
func setMappingPath(root *yaml.Node, path []string, scalar *yaml.Node) error {
	cur := root
	for i, key := range path {
		isLast := i == len(path)-1
		valueNode, found := mappingChild(cur, key)
		switch {
		case isLast && found:
			// In-place stamp of the new scalar onto the existing
			// value node — assign only Kind/Tag/Value so the
			// value-side metadata (LineComment, HeadComment,
			// FootComment, Anchor, Style) survives the patch.
			// yaml.v3 attaches inline trailing comments
			// (`enabled: false  # CHANGEME`) to the VALUE node, so a
			// full node copy (`*valueNode = *scalar`) would silently
			// drop them; the key node carries head comments and is
			// already untouched here.
			valueNode.Kind = scalar.Kind
			valueNode.Tag = scalar.Tag
			valueNode.Value = scalar.Value
			return nil
		case isLast && !found:
			appendMappingChild(cur, key, scalar)
			return nil
		case !isLast && found:
			if valueNode.Kind == yaml.AliasNode {
				return fmt.Errorf("%w: intermediate key %q is a YAML alias; "+
					"PatchScalar cannot traverse aliases — inline the anchored target first",
					driven.ErrYAMLPathInvalid, key)
			}
			if valueNode.Kind != yaml.MappingNode {
				return fmt.Errorf("%w: intermediate key %q is %s, want mapping",
					driven.ErrYAMLPathInvalid, key, kindName(valueNode.Kind))
			}
			cur = valueNode
		case !isLast && !found:
			child := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
			appendMappingChild(cur, key, child)
			cur = child
		}
	}
	return nil
}

// mappingChild returns the value node paired with key in mapping and
// reports whether the pair was found. mapping must be a MappingNode;
// its Content is laid out as [k0, v0, k1, v1, ...].
func mappingChild(mapping *yaml.Node, key string) (*yaml.Node, bool) {
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		k := mapping.Content[i]
		if k.Kind == yaml.ScalarNode && k.Value == key {
			return mapping.Content[i+1], true
		}
	}
	return nil, false
}

// appendMappingChild appends a fresh (keyScalar, value) pair to a
// MappingNode. Used both for missing intermediate mappings and for
// the leaf when the leaf key did not exist before.
func appendMappingChild(mapping *yaml.Node, key string, value *yaml.Node) {
	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: key,
	}
	mapping.Content = append(mapping.Content, keyNode, value)
}

// scalarNode builds a yaml.Node from a Go scalar. The supported
// kinds match the [driven.YAMLCodec.PatchScalar] contract; other
// kinds return [driven.ErrYAMLPathInvalid].
func scalarNode(value any) (*yaml.Node, error) {
	switch v := value.(type) {
	case bool:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: strconv.FormatBool(v)}, nil
	case string:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: v}, nil
	case int:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.FormatInt(int64(v), 10)}, nil
	case int64:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.FormatInt(v, 10)}, nil
	case float64:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: strconv.FormatFloat(v, 'f', -1, 64)}, nil
	default:
		return nil, fmt.Errorf("%w: unsupported scalar type %T", driven.ErrYAMLPathInvalid, value)
	}
}

// kindName returns a human-readable label for a yaml.Kind, used in
// error messages so the caller sees "scalar" / "sequence" instead of
// the numeric constant.
func kindName(k yaml.Kind) string {
	switch k {
	case yaml.DocumentNode:
		return "document"
	case yaml.MappingNode:
		return "mapping"
	case yaml.SequenceNode:
		return "sequence"
	case yaml.ScalarNode:
		return "scalar"
	case yaml.AliasNode:
		return "alias"
	default:
		return "unknown"
	}
}
