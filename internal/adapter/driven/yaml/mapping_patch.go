package yaml

import (
	"bytes"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// Marker-line literals, redundant to the
// internal/hexagon/application/managedblock convention. The adapter
// must not import application packages (depguard rule
// `adapter-no-application`), so the hash-style markers live here too.
// Format: `# BEGIN U-BOOT MANAGED BLOCK: <name>` (and END).
const (
	beginMarkerPrefix = "# BEGIN U-BOOT MANAGED BLOCK: "
	endMarkerPrefix   = "# END U-BOOT MANAGED BLOCK: "
)

// errBlockNotFound is the adapter-local equivalent of
// managedblock.ErrBlockNotFound. Kept private; PatchMappingEntryYAML
// translates it into the absent-block branch and LocateMarkedEntry
// reports it through MarkerSomewhereElse=false / MarkerInEntry=false.
var errBlockNotFound = errors.New("managed block not found")

// errBlockMalformed is the adapter-local equivalent of
// managedblock.ErrBlockMalformed. Surfaces to the caller as a wrapped
// error from both PatchMappingEntryYAML and LocateMarkedEntry.
var errBlockMalformed = errors.New("managed block malformed")

// markerSpan is the byte range of a managed block inside content.
type markerSpan struct {
	beginLineStart int // start of BEGIN line (column 0 even when the marker itself is indented)
	beginLineEnd   int // past BEGIN line's terminating newline
	endLineStart   int // start of END line
	endLineEnd     int // past END line's terminating newline
}

// parentSpan describes a top-level YAML mapping key inside content.
type parentSpan struct {
	headerLineStart  int // start of the `parent:` header line (column 0)
	headerLineEnd    int // past header line's terminating newline
	bodyEnd          int // first byte past the last line that is part of parent's block-style body
	isInlineEmpty    bool // `parent:` / `parent: {}` / `parent: null`
	isInlineNonEmpty bool // `parent: { sub: ... }` or `parent: scalar` (anything non-empty after the colon on the same line)
}

// findManagedBlock locates the first managed block with the given
// name in content. Returns errBlockNotFound when no BEGIN line
// matches, and errBlockMalformed for unmatched-BEGIN or duplicate-
// BEGIN-before-END.
//
// Markers may be indented (block lives nested under a mapping key in
// compose.yaml); the scanner anchors to the literal prefix and
// captures the full enclosing line via newline boundaries.
func findManagedBlock(content []byte, name string) (markerSpan, error) {
	begin := []byte(beginMarkerPrefix + name)
	end := []byte(endMarkerPrefix + name)

	beginIdx := bytes.Index(content, begin)
	if beginIdx == -1 {
		return markerSpan{}, errBlockNotFound
	}

	// BEGIN line: rewind to previous newline (or content start) and
	// advance to next newline (inclusive).
	beginLineStart := beginIdx
	for beginLineStart > 0 && content[beginLineStart-1] != '\n' {
		beginLineStart--
	}
	beginLineEnd := beginIdx + len(begin)
	for beginLineEnd < len(content) && content[beginLineEnd] != '\n' {
		beginLineEnd++
	}
	if beginLineEnd < len(content) {
		beginLineEnd++ // consume \n
	}

	// END line: search forward from beginLineEnd, reject if a second
	// BEGIN line appears between BEGIN and END.
	rest := content[beginLineEnd:]
	endRel := bytes.Index(rest, end)
	if endRel == -1 {
		return markerSpan{}, fmt.Errorf("%w: %q without matching %q",
			errBlockMalformed, beginMarkerPrefix+name, endMarkerPrefix+name)
	}
	if dup := bytes.Index(rest[:endRel], begin); dup != -1 {
		return markerSpan{}, fmt.Errorf("%w: duplicate %q before %q",
			errBlockMalformed, beginMarkerPrefix+name, endMarkerPrefix+name)
	}
	endAbs := beginLineEnd + endRel
	endLineStart := endAbs
	for endLineStart > 0 && content[endLineStart-1] != '\n' {
		endLineStart--
	}
	endLineEnd := endAbs + len(end)
	for endLineEnd < len(content) && content[endLineEnd] != '\n' {
		endLineEnd++
	}
	if endLineEnd < len(content) {
		endLineEnd++ // consume \n
	}

	return markerSpan{
		beginLineStart: beginLineStart,
		beginLineEnd:   beginLineEnd,
		endLineStart:   endLineStart,
		endLineEnd:     endLineEnd,
	}, nil
}

// findTopLevelKey locates a top-level mapping key by literal `<key>:`
// at column 0. Returns (parentSpan, true) on match, (zero, false)
// when the key is absent. Lines inside managed blocks at column 0
// still count as top-level for this scanner — the disambiguation is
// the caller's concern (see PatchMappingEntryYAML).
func findTopLevelKey(content []byte, key string) (parentSpan, bool) {
	keyColon := []byte(key + ":")
	pos := 0
	for pos < len(content) {
		nl := bytes.IndexByte(content[pos:], '\n')
		var line []byte
		var nextPos int
		if nl == -1 {
			line = content[pos:]
			nextPos = len(content)
		} else {
			line = content[pos : pos+nl]
			nextPos = pos + nl + 1
		}
		if bytes.HasPrefix(line, keyColon) {
			rest := line[len(keyColon):]
			if len(rest) == 0 || rest[0] == ' ' || rest[0] == '\t' {
				return classifyParent(content, pos, nextPos, rest), true
			}
		}
		pos = nextPos
	}
	return parentSpan{}, false
}

// classifyParent inspects the rest of the parent header line and the
// following indented body to fill the parentSpan. Reads no further
// than the next column-0 (non-blank) line.
func classifyParent(content []byte, headerStart, headerEnd int, rest []byte) parentSpan {
	ps := parentSpan{
		headerLineStart: headerStart,
		headerLineEnd:   headerEnd,
	}
	trimmed := bytes.TrimSpace(rest)
	switch {
	case len(trimmed) == 0:
		// bare `parent:` — block-style body may follow
		ps.isInlineEmpty = true
	case bytes.Equal(trimmed, []byte("{}")):
		ps.isInlineEmpty = true
	case bytes.Equal(trimmed, []byte("null")) || bytes.Equal(trimmed, []byte("~")):
		ps.isInlineEmpty = true
	default:
		// Anything else on the same line — flow mapping with content,
		// scalar, sequence start, etc. The adapter rejects this case
		// with ErrYAMLFragmentInvalid further up the stack.
		ps.isInlineNonEmpty = true
	}
	ps.bodyEnd = scanIndentedBodyEnd(content, headerEnd)
	return ps
}

// scanIndentedBodyEnd advances past every consecutive line that is
// blank or starts with whitespace, returning the position of the
// first non-blank line that begins at column 0 (or the end of
// content).
func scanIndentedBodyEnd(content []byte, pos int) int {
	for pos < len(content) {
		nl := bytes.IndexByte(content[pos:], '\n')
		var line []byte
		var nextPos int
		if nl == -1 {
			line = content[pos:]
			nextPos = len(content)
		} else {
			line = content[pos : pos+nl]
			nextPos = pos + nl + 1
		}
		if len(bytes.TrimSpace(line)) == 0 {
			pos = nextPos
			continue
		}
		// Non-blank: indented stays in body, column-0 ends body
		if line[0] != ' ' && line[0] != '\t' {
			break
		}
		pos = nextPos
	}
	return pos
}

// LocateMarkedEntry implements [driven.YAMLCodec.LocateMarkedEntry].
func (Codec) LocateMarkedEntry(content []byte, parentKey, entryKey, markerName string) (driven.LocateResult, error) {
	if len(bytes.TrimSpace(content)) > 0 {
		var probe yaml.Node
		if err := yaml.Unmarshal(content, &probe); err != nil {
			return driven.LocateResult{}, fmt.Errorf("locate: parse yaml: %w", err)
		}
	}

	parent, parentFound := findTopLevelKey(content, parentKey)
	mark, markerExists, err := lookupManagedBlock(content, markerName)
	if err != nil {
		return driven.LocateResult{}, err
	}

	res := driven.LocateResult{
		EntryExists: parentEntryPresent(content, parent, parentFound, entryKey),
	}
	if markerExists {
		populateMarkerFlags(&res, content, parent, parentFound, entryKey, mark)
	}
	return res, nil
}

// lookupManagedBlock wraps findManagedBlock with the not-found-vs-
// malformed distinction so the caller only sees (mark, exists, err)
// — keeps LocateMarkedEntry within the cyclomatic-complexity budget.
func lookupManagedBlock(content []byte, markerName string) (markerSpan, bool, error) {
	mark, err := findManagedBlock(content, markerName)
	switch {
	case err == nil:
		return mark, true, nil
	case errors.Is(err, errBlockNotFound):
		return markerSpan{}, false, nil
	default:
		return markerSpan{}, false, err
	}
}

// parentEntryPresent reports whether parentKey.entryKey is present.
// Flow-style non-empty parents (`parent: { sub: … }`) are reported
// as "absent" — they trip ErrYAMLFragmentInvalid at the patch path
// before any decision matters, and the application layer prefers
// the explicit reject over a soft entry-exists branch.
func parentEntryPresent(content []byte, parent parentSpan, parentFound bool, entryKey string) bool {
	if !parentFound || parent.isInlineNonEmpty {
		return false
	}
	return parentContainsEntry(content, parent, entryKey)
}

// populateMarkerFlags fills MarkerInEntry / MarkerSomewhereElse /
// BlockBody on res when a managed block was located. Separated so
// LocateMarkedEntry stays within the nestif budget.
func populateMarkerFlags(res *driven.LocateResult, content []byte, parent parentSpan, parentFound bool, entryKey string, mark markerSpan) {
	insideParent := parentFound &&
		mark.beginLineStart >= parent.headerLineEnd &&
		mark.endLineEnd <= parent.bodyEnd
	if !insideParent {
		res.MarkerSomewhereElse = true
		return
	}
	if res.EntryExists && entryLineInsideMarker(content, parent, entryKey, mark) {
		res.MarkerInEntry = true
		res.BlockBody = sliceBlockBody(content, mark)
		return
	}
	res.MarkerSomewhereElse = true
}

// parentContainsEntry scans the parent's block-style body for a
// mapping entry `entryKey:` whose key sits at the parent's child
// indent. Detection is line-based; markers themselves are skipped
// because their lines never carry a YAML key.
func parentContainsEntry(content []byte, parent parentSpan, entryKey string) bool {
	body := content[parent.headerLineEnd:parent.bodyEnd]
	keyColon := []byte(entryKey + ":")
	for _, line := range splitLines(body) {
		trimmed := bytes.TrimLeft(line, " \t")
		if isMarkerLine(trimmed) {
			continue
		}
		if lineDefinesKey(trimmed, keyColon) {
			return true
		}
	}
	return false
}

// splitLines yields the slice of byte-lines in body (without the
// terminating newline characters). Used by the line-based scanners
// so each one can stay readable without juggling positions.
func splitLines(body []byte) [][]byte {
	var lines [][]byte
	pos := 0
	for pos < len(body) {
		nl := bytes.IndexByte(body[pos:], '\n')
		if nl == -1 {
			lines = append(lines, body[pos:])
			break
		}
		lines = append(lines, body[pos:pos+nl])
		pos = pos + nl + 1
	}
	return lines
}

// isMarkerLine reports whether the trimmed line begins with one of
// the managed-block marker prefixes.
func isMarkerLine(trimmed []byte) bool {
	return bytes.HasPrefix(trimmed, []byte(beginMarkerPrefix)) ||
		bytes.HasPrefix(trimmed, []byte(endMarkerPrefix))
}

// lineDefinesKey reports whether the trimmed line starts with the
// `<key>:` byte sequence and the next byte (if any) is whitespace —
// the YAML cue for a mapping entry header.
func lineDefinesKey(trimmed, keyColon []byte) bool {
	if !bytes.HasPrefix(trimmed, keyColon) {
		return false
	}
	rest := trimmed[len(keyColon):]
	return len(rest) == 0 || rest[0] == ' ' || rest[0] == '\t'
}

// entryLineInsideMarker reports whether the `entryKey:` line in the
// parent body lies between the BEGIN and END marker lines.
func entryLineInsideMarker(content []byte, parent parentSpan, entryKey string, mark markerSpan) bool {
	body := content[parent.headerLineEnd:parent.bodyEnd]
	keyColon := []byte(entryKey + ":")
	pos := 0
	for pos < len(body) {
		nl := bytes.IndexByte(body[pos:], '\n')
		var line []byte
		var nextPos int
		if nl == -1 {
			line = body[pos:]
			nextPos = len(body)
		} else {
			line = body[pos : pos+nl]
			nextPos = pos + nl + 1
		}
		trimmed := bytes.TrimLeft(line, " \t")
		if bytes.HasPrefix(trimmed, keyColon) {
			rest := trimmed[len(keyColon):]
			if len(rest) == 0 || rest[0] == ' ' || rest[0] == '\t' {
				abs := parent.headerLineEnd + pos
				if abs >= mark.beginLineEnd && abs < mark.endLineStart {
					return true
				}
			}
		}
		pos = nextPos
	}
	return false
}

// sliceBlockBody returns a copy of the bytes between the BEGIN and
// END marker lines (exclusive). Used to populate
// LocateResult.BlockBody so the caller (T4c) can run a content
// presence check without a second pass through the adapter.
func sliceBlockBody(content []byte, mark markerSpan) []byte {
	body := content[mark.beginLineEnd:mark.endLineStart]
	out := make([]byte, len(body))
	copy(out, body)
	return out
}

// PatchMappingEntryYAML implements [driven.YAMLCodec.PatchMappingEntryYAML].
func (Codec) PatchMappingEntryYAML(content []byte, parentKey, entryKey string, valueYAML []byte, markerName string) ([]byte, error) {
	if err := validateFragmentMappingRoot(valueYAML); err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(content)) > 0 {
		if err := assertNoTopLevelDuplicate(content); err != nil {
			return nil, err
		}
	}

	// Symmetric anchor check via the shared scanner.
	mark, mErr := findManagedBlock(content, markerName)
	if mErr != nil && !errors.Is(mErr, errBlockNotFound) {
		return nil, mErr
	}
	markerExists := mErr == nil

	parent, parentFound := findTopLevelKey(content, parentKey)

	if markerExists {
		// Marker must be inside parent body AND surround the entryKey
		// mapping entry; otherwise reject as anchor-mismatch.
		if !parentFound || mark.beginLineStart < parent.headerLineEnd || mark.endLineEnd > parent.bodyEnd {
			return nil, fmt.Errorf("%w: marker %q is outside %s mapping",
				driven.ErrYAMLAnchorMismatch, markerName, parentKey)
		}
		if parent.isInlineNonEmpty {
			return nil, fmt.Errorf("%w: parent %s is a non-empty flow mapping",
				driven.ErrYAMLFragmentInvalid, parentKey)
		}
		if !entryLineInsideMarker(content, parent, entryKey, mark) {
			return nil, fmt.Errorf("%w: marker %q does not enclose %s.%s",
				driven.ErrYAMLAnchorMismatch, markerName, parentKey, entryKey)
		}
		return replaceExistingBlock(content, mark, entryKey, valueYAML, markerName), nil
	}

	// No marker yet. Check entry presence under parent: if it
	// already exists as a user-manual entry, refuse.
	if parentFound {
		if parent.isInlineNonEmpty {
			return nil, fmt.Errorf("%w: parent %s is a non-empty flow mapping",
				driven.ErrYAMLFragmentInvalid, parentKey)
		}
		if parentContainsEntry(content, parent, entryKey) {
			return nil, fmt.Errorf("%w: %s.%s exists without a managed marker",
				driven.ErrYAMLAnchorMismatch, parentKey, entryKey)
		}
		return insertUnderParent(content, parent, parentKey, entryKey, valueYAML, markerName), nil
	}
	return appendNewParent(content, parentKey, entryKey, valueYAML, markerName), nil
}

// validateFragmentMappingRoot ensures valueYAML parses as a YAML
// mapping root. Scalars and sequences are rejected with
// ErrYAMLFragmentInvalid because their insertion would not produce
// a valid managed entry shape (`<entryKey>: <mapping body>`).
func validateFragmentMappingRoot(valueYAML []byte) error {
	if len(bytes.TrimSpace(valueYAML)) == 0 {
		// Empty mapping fragment is fine (renders as `{}`).
		return nil
	}
	var node yaml.Node
	if err := yaml.Unmarshal(valueYAML, &node); err != nil {
		return fmt.Errorf("%w: valueYAML parse: %v", driven.ErrYAMLFragmentInvalid, err)
	}
	root := &node
	if root.Kind == yaml.DocumentNode {
		if len(root.Content) == 0 {
			return nil
		}
		root = root.Content[0]
	}
	switch root.Kind {
	case yaml.MappingNode:
		return nil
	case yaml.ScalarNode:
		return fmt.Errorf("%w: valueYAML is a scalar, want mapping",
			driven.ErrYAMLFragmentInvalid)
	case yaml.SequenceNode:
		return fmt.Errorf("%w: valueYAML is a sequence, want mapping",
			driven.ErrYAMLFragmentInvalid)
	default:
		return fmt.Errorf("%w: valueYAML kind %d unsupported",
			driven.ErrYAMLFragmentInvalid, root.Kind)
	}
}

// assertNoTopLevelDuplicate parses content with yaml.v3 in
// duplicate-key-strict mode and reports duplicate top-level keys as
// ErrYAMLFragmentInvalid. Sub-tree duplicates are not the patcher's
// concern.
func assertNoTopLevelDuplicate(content []byte) error {
	var node yaml.Node
	if err := yaml.Unmarshal(content, &node); err != nil {
		return fmt.Errorf("patch: parse content: %w", err)
	}
	root := &node
	if root.Kind == yaml.DocumentNode {
		if len(root.Content) == 0 {
			return nil
		}
		root = root.Content[0]
	}
	if root.Kind != yaml.MappingNode {
		return nil
	}
	seen := make(map[string]bool, len(root.Content)/2)
	for i := 0; i+1 < len(root.Content); i += 2 {
		k := root.Content[i]
		if k.Kind != yaml.ScalarNode {
			continue
		}
		if seen[k.Value] {
			return fmt.Errorf("%w: duplicate top-level key %q",
				driven.ErrYAMLFragmentInvalid, k.Value)
		}
		seen[k.Value] = true
	}
	return nil
}

// renderManagedEntry returns the bytes of the managed entry, ready
// to splice under the parent's body indent. The output starts at
// column 0 and includes BEGIN line, entryKey + indented valueYAML,
// END line — each terminated by `\n`. The caller indents the whole
// chunk by the desired parent-body indent.
func renderManagedEntry(entryKey string, valueYAML []byte, markerName string) []byte {
	var b bytes.Buffer
	b.WriteString(beginMarkerPrefix)
	b.WriteString(markerName)
	b.WriteByte('\n')
	b.WriteString(entryKey)
	b.WriteString(":")
	trimmedValue := bytes.TrimRight(valueYAML, "\n\r ")
	if len(trimmedValue) == 0 {
		// empty mapping body — render as `entryKey: {}`
		b.WriteString(" {}\n")
	} else {
		b.WriteByte('\n')
		// indent valueYAML by 2 spaces so it sits under entryKey
		for _, line := range bytes.Split(trimmedValue, []byte("\n")) {
			b.WriteString("  ")
			b.Write(line)
			b.WriteByte('\n')
		}
	}
	b.WriteString(endMarkerPrefix)
	b.WriteString(markerName)
	b.WriteByte('\n')
	return b.Bytes()
}

// indentLines prefixes every non-empty line with the given indent
// string. Used to push the rendered managed entry under the parent's
// child-indent level.
func indentLines(in []byte, indent string) []byte {
	var b bytes.Buffer
	for i, line := range bytes.Split(in, []byte("\n")) {
		if i > 0 {
			b.WriteByte('\n')
		}
		if len(line) > 0 {
			b.WriteString(indent)
			b.Write(line)
		}
	}
	return b.Bytes()
}

// replaceExistingBlock substitutes a freshly rendered managed entry
// in place of the existing [mark.beginLineStart, mark.endLineEnd]
// byte range, preserving everything outside.
func replaceExistingBlock(content []byte, mark markerSpan, entryKey string, valueYAML []byte, markerName string) []byte {
	// Determine the marker's leading indent (number of spaces before
	// "# BEGIN") so the new block sits at the same indent.
	indent := ""
	for i := mark.beginLineStart; i < mark.beginLineEnd && (content[i] == ' ' || content[i] == '\t'); i++ {
		indent += string(content[i])
	}
	rendered := renderManagedEntry(entryKey, valueYAML, markerName)
	indented := indentLines(rendered, indent)
	out := make([]byte, 0, len(content)-(mark.endLineEnd-mark.beginLineStart)+len(indented))
	out = append(out, content[:mark.beginLineStart]...)
	out = append(out, indented...)
	out = append(out, content[mark.endLineEnd:]...)
	return out
}

// insertUnderParent appends a fresh managed entry at the end of the
// parent's block-style body (or converts the parent's inline-empty
// form into block-style first). Existing siblings of the new entry
// stay byte-identical.
func insertUnderParent(content []byte, parent parentSpan, parentKey, entryKey string, valueYAML []byte, markerName string) []byte {
	const childIndent = "  "

	// Convert `parent:` / `parent: {}` / `parent: null` to block
	// style by replacing the header line with `parent:\n` and
	// inserting the body straight afterwards.
	header := content[parent.headerLineStart:parent.headerLineEnd]
	rewriteHeader := false
	if parent.isInlineEmpty {
		// detect if header has anything after the colon (e.g. `{}` or `null`)
		colonIdx := bytes.IndexByte(header, ':')
		if colonIdx != -1 {
			tail := bytes.TrimSpace(header[colonIdx+1:])
			if len(tail) > 0 {
				rewriteHeader = true
			}
		}
	}

	rendered := renderManagedEntry(entryKey, valueYAML, markerName)
	indented := indentLines(rendered, childIndent)

	var out bytes.Buffer
	out.Grow(len(content) + len(indented) + 4)
	out.Write(content[:parent.headerLineStart])
	if rewriteHeader {
		out.WriteString(parentKey)
		out.WriteString(":\n")
	} else {
		out.Write(header)
	}
	// existing body (lines between header and bodyEnd)
	out.Write(content[parent.headerLineEnd:parent.bodyEnd])
	// ensure the appended entry sits on its own line — if the body
	// ends without a trailing newline (e.g. EOF without newline),
	// add one.
	if parent.bodyEnd > parent.headerLineEnd && content[parent.bodyEnd-1] != '\n' {
		out.WriteByte('\n')
	}
	out.Write(indented)
	out.Write(content[parent.bodyEnd:])
	return out.Bytes()
}

// appendNewParent appends a new top-level parent block at the end of
// content. The new block is separated from the previous content by
// exactly one blank line.
func appendNewParent(content []byte, parentKey, entryKey string, valueYAML []byte, markerName string) []byte {
	const childIndent = "  "
	rendered := renderManagedEntry(entryKey, valueYAML, markerName)
	indented := indentLines(rendered, childIndent)
	var out bytes.Buffer
	out.Grow(len(content) + len(indented) + len(parentKey) + 8)
	out.Write(content)
	if len(content) > 0 && content[len(content)-1] != '\n' {
		out.WriteByte('\n')
	}
	if len(content) > 0 {
		out.WriteByte('\n') // separator blank line
	}
	out.WriteString(parentKey)
	out.WriteString(":\n")
	out.Write(indented)
	return out.Bytes()
}

