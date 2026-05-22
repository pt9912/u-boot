// Package managedblock parses and rewrites U-BOOT MANAGED BLOCK
// regions in u-boot-generated files (LH-SA-FILE-002). A managed
// block is a contiguous line range delimited by BEGIN/END comment
// markers in the host file's comment syntax — see [Style] for the
// supported syntaxes:
//
//	# BEGIN U-BOOT MANAGED BLOCK: <name>         (YAML, .env, …)
//	<!-- BEGIN U-BOOT MANAGED BLOCK: <name> -->  (Markdown)
//	// BEGIN U-BOOT MANAGED BLOCK: <name>        (JSONC, JS/TS, …)
//
// The package is a pure-text parser — it does not know about Go
// templates, YAML, or any embedding library. Callers render the
// replacement block themselves and ask [Replace] to splice it in.
// The LH-FA-INIT-005 re-init flow (M3-T4b) consumes this package to
// support the §611-§614 "only the block is changed" behaviour for
// structured configuration files (compose.yaml, .env.example,
// README.md, CHANGELOG.md, .devcontainer/devcontainer.json).
package managedblock

import (
	"errors"
	"fmt"
	"regexp"
)

// Style names the comment syntax of the host file. The chosen style
// determines the literal BEGIN/END marker text [Find] looks for.
type Style int

const (
	// StyleHash is the `#`-comment style — used by YAML, .env,
	// Dockerfile, shell scripts.
	StyleHash Style = iota
	// StyleDoubleSlash is the `//`-comment style — used by JSONC,
	// JavaScript, TypeScript, Go.
	StyleDoubleSlash
	// StyleHTMLComment is the `<!-- … -->` HTML-comment style —
	// used by Markdown.
	StyleHTMLComment
)

// String returns the human-readable name of the style.
func (s Style) String() string {
	switch s {
	case StyleHash:
		return "hash"
	case StyleDoubleSlash:
		return "double-slash"
	case StyleHTMLComment:
		return "html-comment"
	}
	return fmt.Sprintf("Style(%d)", int(s))
}

// Marker identifies a single managed block: comment style + block
// name. The name is the user-visible label after the colon in the
// marker line (e.g. `init` for the u-boot init-scaffolding block,
// `postgres` for the spec example).
type Marker struct {
	Style Style
	Name  string
}

// Begin returns the BEGIN-marker line text (without leading
// whitespace or trailing newline).
func (m Marker) Begin() string {
	return wrap(m.Style, "BEGIN U-BOOT MANAGED BLOCK: "+m.Name)
}

// End returns the END-marker line text.
func (m Marker) End() string {
	return wrap(m.Style, "END U-BOOT MANAGED BLOCK: "+m.Name)
}

// wrap formats inner into the comment syntax for style.
func wrap(style Style, inner string) string {
	switch style {
	case StyleHash:
		return "# " + inner
	case StyleDoubleSlash:
		return "// " + inner
	case StyleHTMLComment:
		return "<!-- " + inner + " -->"
	}
	return inner
}

// ErrBlockNotFound is returned by [Find] / [Replace] when no BEGIN
// marker for the requested [Marker] is present in the content.
// Callers branch on this sentinel to decide between "patch the
// existing block" and "fall back to a full overwrite or abort".
var ErrBlockNotFound = errors.New("managed block not found")

// ErrBlockMalformed is returned when a BEGIN marker exists but no
// matching END marker is found after it (or vice-versa). Indicates
// the file was hand-edited into an invalid state; callers must not
// silently auto-repair.
var ErrBlockMalformed = errors.New("managed block malformed")

// Find returns the byte offsets [start, end) of the managed-block
// region named by m in content. The region spans from the start of
// the BEGIN line (column 0, including any leading whitespace) to
// just past the END line's terminating newline (or end of content
// if the END line is the last line). The returned offsets are
// suitable for direct splice into content[:start] + … + content[end:].
//
// Returns [ErrBlockNotFound] when no BEGIN marker is present;
// [ErrBlockMalformed] when BEGIN is present but END is not (or END
// appears before BEGIN). The regex matches markers regardless of
// leading whitespace, so indented blocks (e.g. nested under
// `services:` in compose.yaml) are detected.
func Find(content []byte, m Marker) (int, int, error) {
	beginRE, endRE, err := compileMarkerRegexps(m)
	if err != nil {
		return 0, 0, err
	}
	beginLoc := beginRE.FindIndex(content)
	if beginLoc == nil {
		return 0, 0, ErrBlockNotFound
	}
	endLoc := endRE.FindIndex(content[beginLoc[1]:])
	if endLoc == nil {
		return 0, 0, fmt.Errorf("%w: %s present without %s",
			ErrBlockMalformed, m.Begin(), m.End())
	}
	start := beginLoc[0]
	end := beginLoc[1] + endLoc[1]
	// Consume the END line's trailing newline (\r\n or \n) so the
	// returned region is line-aligned and a splice does not leave
	// a blank line behind.
	if end < len(content) && content[end] == '\r' {
		end++
	}
	if end < len(content) && content[end] == '\n' {
		end++
	}
	return start, end, nil
}

// Has reports whether content contains the managed block named by m.
// Thin wrapper around [Find]; treat as "is the block here at all?",
// not "is the block well-formed?" — see [Find] for the malformed
// case.
func Has(content []byte, m Marker) bool {
	_, _, err := Find(content, m)
	return err == nil
}

// Replace returns content with the managed-block region for m
// replaced by replacement. The replacement must include BEGIN/END
// marker lines for m — Replace does not synthesize them; it just
// splices the byte range from [Find].
//
// Returns [ErrBlockNotFound] or [ErrBlockMalformed] (both wrapped)
// when [Find] cannot locate the block. The returned slice is a
// freshly-allocated copy; content is not modified.
func Replace(content []byte, m Marker, replacement []byte) ([]byte, error) {
	start, end, err := Find(content, m)
	if err != nil {
		return nil, err
	}
	out := make([]byte, 0, len(content)-(end-start)+len(replacement))
	out = append(out, content[:start]...)
	out = append(out, replacement...)
	out = append(out, content[end:]...)
	return out, nil
}

// compileMarkerRegexps builds the BEGIN/END regexps for marker m.
// The patterns use `(?m)` multiline mode so `^` and `$` match line
// boundaries and `[\t ]*` (not `\s`) so the matcher does not jump
// across line breaks. `\r?$` tolerates CRLF line endings.
func compileMarkerRegexps(m Marker) (*regexp.Regexp, *regexp.Regexp, error) {
	beginPattern := `(?m)^[\t ]*` + regexp.QuoteMeta(m.Begin()) + `[\t ]*\r?$`
	endPattern := `(?m)^[\t ]*` + regexp.QuoteMeta(m.End()) + `[\t ]*\r?$`
	beginRE, err := regexp.Compile(beginPattern)
	if err != nil {
		return nil, nil, fmt.Errorf("compile begin pattern: %w", err)
	}
	endRE, err := regexp.Compile(endPattern)
	if err != nil {
		return nil, nil, fmt.Errorf("compile end pattern: %w", err)
	}
	return beginRE, endRE, nil
}
