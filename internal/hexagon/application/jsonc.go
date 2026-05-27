package application

// stripJSONC returns src with JSONC-extensions removed so the
// result is plain JSON consumable by [encoding/json]. JSONC, used
// by VS Code Dev Containers, extends JSON with:
//
//   - Line comments `// ...` to end-of-line.
//   - Block comments `/* ... */`.
//   - Trailing commas before `}` and `]`.
//
// The stripper is string-aware: comment characters inside a JSON
// string literal are preserved verbatim. Escape sequences inside
// strings (`\\"`, `\\\\`, etc.) are honored so a quote inside a
// string does not prematurely close it.
//
// Known limitations:
//   - The stripper does not validate the JSONC structure beyond
//     comment / trailing-comma removal. Malformed input is passed
//     through unchanged from the comment/comma standpoint;
//     [encoding/json] then reports the parse error downstream.
//   - Unterminated block comments swallow the rest of the file.
//     That is consistent with how a strict JSONC parser would react
//     (the input is malformed); the downstream JSON parser will see
//     a truncated document and report an unexpected-EOF error.
//
// Stdlib-only. A more rigorous alternative would be a real JSONC
// library (e.g. tailscale/hujson); avoided here to keep `go.mod`
// minimal and the LH-QA-004-gomodguard-blocklist untouched.
func stripJSONC(src []byte) []byte {
	out := make([]byte, 0, len(src))
	for i := 0; i < len(src); {
		switch {
		case src[i] == '"':
			end := scanString(src, i)
			out = append(out, src[i:end]...)
			i = end
		case isCommentStart(src, i, '/'):
			i = skipLineComment(src, i)
		case isCommentStart(src, i, '*'):
			i = skipBlockComment(src, i)
		case src[i] == ',' && isTrailingComma(src, i):
			i = nextNonSpaceIdx(src, i+1)
		default:
			out = append(out, src[i])
			i++
		}
	}
	return out
}

// isCommentStart reports whether src[i:] begins a `//` (when c is
// `/`) or `/*` (when c is `*`) JSONC comment marker.
func isCommentStart(src []byte, i int, second byte) bool {
	return src[i] == '/' && i+1 < len(src) && src[i+1] == second
}

// skipLineComment returns the index past the `\n` that terminates
// the line comment starting at src[i] (which must be `//`). The
// newline is preserved by the caller so JSON parser-error line
// numbers stay accurate.
func skipLineComment(src []byte, i int) int {
	i += 2
	for i < len(src) && src[i] != '\n' {
		i++
	}
	return i
}

// skipBlockComment returns the index past the `*/` that terminates
// the block comment starting at src[i] (which must be `/*`). An
// unterminated block comment swallows the rest of the input.
func skipBlockComment(src []byte, i int) int {
	i += 2
	for i+1 < len(src) && (src[i] != '*' || src[i+1] != '/') {
		i++
	}
	if i+1 < len(src) {
		return i + 2
	}
	return len(src)
}

// isTrailingComma reports whether the comma at src[i] is a JSONC
// trailing comma (i.e. followed only by whitespace before `}` or
// `]`).
func isTrailingComma(src []byte, i int) bool {
	j := nextNonSpaceIdx(src, i+1)
	return j < len(src) && (src[j] == '}' || src[j] == ']')
}

// nextNonSpaceIdx returns the index of the first non-whitespace
// byte at-or-after start, or len(src) if none.
func nextNonSpaceIdx(src []byte, start int) int {
	for start < len(src) && isJSONCSpace(src[start]) {
		start++
	}
	return start
}

// scanString returns the index past the closing `"` of a JSON
// string literal that starts at src[start] (which must be `"`).
// Escape sequences (`\"`, `\\`, etc.) are honored.
func scanString(src []byte, start int) int {
	i := start + 1
	for i < len(src) {
		switch src[i] {
		case '\\':
			// Skip the escape pair; out-of-range is safe because the
			// loop bound stops us before going past the buffer.
			i += 2
		case '"':
			return i + 1
		default:
			i++
		}
	}
	return i
}

// isJSONCSpace reports whether b is a JSON whitespace character.
func isJSONCSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}
