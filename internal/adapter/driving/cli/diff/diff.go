// Package diff computes unified-diff hunks between two file snapshots
// captured by the recordingfs driven adapter. It is the CLI-adapter-
// internal renderer for LH-FA-CLI-008 `--diff` output, both as a
// human-readable unified string (see [Render]) and as the structured
// `plannedFiles[].hunks` array of the LH-FA-CLI-007 §326 voll-schema
// JSON envelope (see [Compute]).
//
// Algorithmic choice (slice-v1-cli-json-dry-run-add T0-(d)): pure-Go
// LCS via dynamic programming, no new go.mod dependency. The DP table
// is O(m*n) in memory which is acceptable for the file sizes the add
// surface generates (compose.yaml additions ≤ 50 lines, .env.example
// updates ≤ 20 lines).
//
// Binary content (either side fails [utf8.Valid]) cannot meaningfully
// participate in line-oriented diffing. T0-(l) prescribes a
// Spec-konformes Fallback: callers detect via [IsBinary], render no
// hunks, and use [CountBytesDiff] for `changes[].count` instead of
// the line-based [CountFromHunks]/[CountLines] semantics.
package diff

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// defaultContext mirrors the unified-diff convention of three context
// lines around each change cluster. Hard-coded for V1 — exposing it
// as a knob would proliferate options without a real use case.
const defaultContext = 3

// IsBinary reports whether either content fails UTF-8 validation. A
// nil/empty byte slice is considered text (utf8.Valid returns true for
// empty input), so create- and delete-actions with only one non-empty
// side still flow through the line-based path.
func IsBinary(oldContent, newContent []byte) bool {
	return !utf8.Valid(oldContent) || !utf8.Valid(newContent)
}

// CountLines returns a trailing-newline-robust line count of content
// per slice-v1-cli-json-dry-run-add T0-(g). Idioms pinned in the slice
// plan:
//
//	""        → 0
//	"a"       → 1   (unterminated last line counts)
//	"a\n"     → 1
//	"a\nb"    → 2
//	"a\nb\n"  → 2
//
// Generated YAML and .env templates conventionally end with a trailing
// newline; this form keeps Spec §430 (`count: 12` for a 12-line block)
// stable regardless of which convention the template author followed.
func CountLines(content []byte) int {
	n := bytes.Count(content, []byte("\n"))
	if len(content) > 0 && !bytes.HasSuffix(content, []byte("\n")) {
		n++
	}
	return n
}

// CountFromHunks sums the NewLines field across all hunks — the
// `changes[].count` value for action "modify" per T0-(g).
func CountFromHunks(hunks []driving.Hunk) int {
	total := 0
	for _, h := range hunks {
		total += h.NewLines
	}
	return total
}

// CountBytesDiff is the binary-content fallback for `changes[].count`
// (T0-(l)). For text content, prefer [CountLines] (create) or
// [CountFromHunks] (modify) — they carry richer line-oriented
// semantics. Returns the absolute byte-length difference.
func CountBytesDiff(oldContent, newContent []byte) int {
	diff := len(newContent) - len(oldContent)
	if diff < 0 {
		return -diff
	}
	return diff
}

// Compute returns the unified-diff hunks between oldContent and
// newContent. Identical inputs produce nil. Binary content (per
// [IsBinary]) also produces nil — callers must check IsBinary first
// and fall back to byte-count semantics.
//
// Each Hunk.Content carries the line block with leading-character
// markers (space for context, `+` for additions, `-` for deletions),
// one line per text-line, separated by `\n` and terminated with `\n`.
// The unified header (`@@ -...,... +...,... @@`) is NOT included in
// Hunk.Content — it belongs to the human-mode wrapping rendered by
// [Render]; JSON consumers read the four integer fields directly.
func Compute(oldContent, newContent []byte) []driving.Hunk {
	if IsBinary(oldContent, newContent) {
		return nil
	}
	oldLines := splitLines(oldContent)
	newLines := splitLines(newContent)
	ops := lcsOps(oldLines, newLines)
	annotated := annotateOps(ops)
	return groupHunks(annotated, defaultContext)
}

// Render produces the human-mode unified-diff string for `--diff`
// without `--json` (LH-FA-CLI-008). Hunks are concatenated with their
// standard `@@ -oldStart,oldLines +newStart,newLines @@` headers; an
// empty hunk slice produces an empty string. Output is suitable for
// printing to stdout as-is.
func Render(hunks []driving.Hunk) string {
	if len(hunks) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, h := range hunks {
		fmt.Fprintf(&sb, "@@ -%d,%d +%d,%d @@\n", h.OldStart, h.OldLines, h.NewStart, h.NewLines)
		sb.WriteString(h.Content)
	}
	return sb.String()
}

// splitLines turns content into a line slice with stable round-trip
// semantics for trailing newlines: a final `\n` is treated as a line
// terminator (not a separator), so `"a\nb\n"` and `"a\nb"` both yield
// `["a", "b"]`. Empty input yields nil (not a one-element slice with
// the empty string). The trailing-newline difference is invisible to
// LCS; for V1 add templates (UNIX-style trailing newline by
// convention) that lossy form is acceptable, and slice-v1-cli-json-
// dry-run-add T0-(g) anchors `changes[].count` against the line-
// counted form so the loss doesn't bleed into observable output.
func splitLines(b []byte) []string {
	if len(b) == 0 {
		return nil
	}
	s := strings.TrimSuffix(string(b), "\n")
	return strings.Split(s, "\n")
}

type opKind int

const (
	opEqual opKind = iota
	opDelete
	opInsert
)

type lineOp struct {
	kind opKind
	text string
}

type annotatedOp struct {
	lineOp
	oldLine int // 1-based; 0 for inserts (no line in old)
	newLine int // 1-based; 0 for deletes (no line in new)
}

// lcsOps returns the edit operations transforming a into b, computed
// via the classic Wagner–Fischer LCS dynamic-programming table. Time
// and space are O(len(a)*len(b)); fine for the file sizes the add
// surface produces. For larger files (≥ ~10k lines) Myers' O(ND)
// would be the obvious replacement, but V1 doesn't need it.
func lcsOps(a, b []string) []lineOp {
	if len(a) == 0 && len(b) == 0 {
		return nil
	}
	dp := lcsTable(a, b)
	return backtrackOps(a, b, dp)
}

func lcsTable(a, b []string) [][]int {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				continue
			}
			if dp[i-1][j] >= dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}
	return dp
}

func backtrackOps(a, b []string, dp [][]int) []lineOp {
	ops := make([]lineOp, 0, len(a)+len(b))
	i, j := len(a), len(b)
	for i > 0 || j > 0 {
		switch {
		case i > 0 && j > 0 && a[i-1] == b[j-1]:
			ops = append(ops, lineOp{kind: opEqual, text: a[i-1]})
			i--
			j--
		case j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]):
			ops = append(ops, lineOp{kind: opInsert, text: b[j-1]})
			j--
		default:
			ops = append(ops, lineOp{kind: opDelete, text: a[i-1]})
			i--
		}
	}
	// Reverse in-place.
	for l, r := 0, len(ops)-1; l < r; l, r = l+1, r-1 {
		ops[l], ops[r] = ops[r], ops[l]
	}
	return ops
}

// annotateOps tags each op with its 1-based line numbers in the old
// and new files (0 where the op type doesn't appear on that side).
func annotateOps(ops []lineOp) []annotatedOp {
	out := make([]annotatedOp, len(ops))
	oldNo, newNo := 1, 1
	for i, op := range ops {
		out[i].lineOp = op
		switch op.kind {
		case opEqual:
			out[i].oldLine = oldNo
			out[i].newLine = newNo
			oldNo++
			newNo++
		case opDelete:
			out[i].oldLine = oldNo
			oldNo++
		case opInsert:
			out[i].newLine = newNo
			newNo++
		}
	}
	return out
}

// groupHunks clusters non-equal ops into hunks, padded by `context`
// equal lines on each side. Two clusters merge if the gap between
// them is ≤ 2*context (they would share context lines anyway).
func groupHunks(ops []annotatedOp, context int) []driving.Hunk {
	var changes []int
	for i, op := range ops {
		if op.kind != opEqual {
			changes = append(changes, i)
		}
	}
	if len(changes) == 0 {
		return nil
	}
	type cluster struct{ first, last int }
	clusters := []cluster{{changes[0], changes[0]}}
	for i := 1; i < len(changes); i++ {
		gap := changes[i] - clusters[len(clusters)-1].last - 1
		if gap <= 2*context {
			clusters[len(clusters)-1].last = changes[i]
		} else {
			clusters = append(clusters, cluster{changes[i], changes[i]})
		}
	}
	hunks := make([]driving.Hunk, 0, len(clusters))
	for _, c := range clusters {
		start := c.first - context
		if start < 0 {
			start = 0
		}
		end := c.last + context
		if end >= len(ops) {
			end = len(ops) - 1
		}
		hunks = append(hunks, buildHunk(ops, start, end))
	}
	return hunks
}

// buildHunk renders ops[start..end] (inclusive) into a single Hunk.
// Coordinates derive from the first op whose kind contributes to that
// side: pure additions (no equal/delete in range) leave OldStart at 0
// per the unified-diff convention "@@ -0,0 +N,M @@" for new-file
// hunks. Symmetrically for pure deletions and NewStart.
func buildHunk(ops []annotatedOp, start, end int) driving.Hunk {
	var oldStart, newStart, oldLines, newLines int
	var content strings.Builder
	for i := start; i <= end; i++ {
		op := ops[i]
		switch op.kind {
		case opEqual:
			if oldStart == 0 {
				oldStart = op.oldLine
			}
			if newStart == 0 {
				newStart = op.newLine
			}
			oldLines++
			newLines++
			content.WriteByte(' ')
		case opDelete:
			if oldStart == 0 {
				oldStart = op.oldLine
			}
			oldLines++
			content.WriteByte('-')
		case opInsert:
			if newStart == 0 {
				newStart = op.newLine
			}
			newLines++
			content.WriteByte('+')
		}
		content.WriteString(op.text)
		content.WriteByte('\n')
	}
	return driving.Hunk{
		OldStart: oldStart,
		OldLines: oldLines,
		NewStart: newStart,
		NewLines: newLines,
		Content:  content.String(),
	}
}
