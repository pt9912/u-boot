package cli

import (
	"github.com/pt9912/u-boot/internal/adapter/driving/cli/diff"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// mapPlannedFilesToWire converts a recorder-captured
// []driving.PlannedFile slice into the CLI wire-types
// (plannedFile + changeEntry). The Hunks-field is populated only
// when withHunks is true (--diff / preview-and-apply); `changes[].
// count` always follows the T0-(g) semantics regardless of flag
// state — for modify-actions that means we compute hunks even
// without --diff just to sum their additions, since the alternative
// (CountLines on the whole new file) overstated the count by orders
// of magnitude for any add-on-into-existing-file case (add review
// finding #1 / Spec §477).
//
// Originally lived in add.go as mapResponseToWire(resp, withHunks)
// taking the concrete AddServiceResponse; slice-v1-cli-json-dry-run-
// init T0-(e) decomposed-Helper-Signatur takes the
// []driving.PlannedFile directly so init/remove/generate/config-set
// can reuse it without depending on AddServiceResponse-shape (their
// own InitProjectResponse/etc. have separate PlannedFiles fields).
func mapPlannedFilesToWire(planned []driving.PlannedFile, withHunks bool) ([]plannedFile, []changeEntry) {
	pfs := make([]plannedFile, 0, len(planned))
	chs := make([]changeEntry, 0, len(planned))
	for _, pf := range planned {
		wirePF := plannedFile{Path: pf.Path, Action: pf.Action}
		count, hunks := computeChangeCountAndHunks(pf)
		if withHunks && len(hunks) > 0 {
			wirePF.Hunks = toCLIHunks(hunks)
		}
		pfs = append(pfs, wirePF)
		chs = append(chs, changeEntry{Path: pf.Path, Count: count})
	}
	return pfs, chs
}

// computeChangeCountAndHunks applies the T0-(g) `changes[].count`
// semantics AND returns the hunks (or nil for binary/no-change paths)
// so the caller can re-use them for the wire-Hunks field when --diff
// is set. The double-return keeps the diff invocation single per
// PlannedFile regardless of flag combination.
//
// Action-rules:
//   - "create": count = CountLines(NewContent), hunks computed for
//     full-file insertion shape.
//   - "modify": count = CountAdditions(hunks) — only the `+` lines
//     (Spec §477 pins `count: 6` for a 6-line postgres-block append,
//     NOT 6 + context lines; add review-round-7 finding B).
//   - "delete" (slice-v1-cli-json-dry-run-remove T0-(p)): count = 0
//     (delete contributes zero added lines per Spec §477 semantics);
//     hunks rendert den Old-Inhalt als full-file-Remove-Block, damit
//     `--diff --json`-Konsumenten sehen WAS gelöscht wird. Add review
//     #8 binary-delete-Trap bleibt geschützt: ein Pre-Compute
//     IsBinary-Check unterdrückt Hunks (und `CountBytesDiff` — wir
//     wollen `0` und nicht den Byte-Delta) für binary content.
//   - binary content (non-delete): count = CountBytesDiff, hunks=nil
//     so wirePF.Hunks remains omitted (T0-(l) Spec-konformes
//     Fallback).
func computeChangeCountAndHunks(pf driving.PlannedFile) (int, []driving.Hunk) {
	if pf.Action == "delete" {
		if diff.IsBinary(pf.OldContent, pf.NewContent) {
			return 0, nil
		}
		return 0, diff.Compute(pf.OldContent, pf.NewContent)
	}
	if diff.IsBinary(pf.OldContent, pf.NewContent) {
		return diff.CountBytesDiff(pf.OldContent, pf.NewContent), nil
	}
	hunks := diff.Compute(pf.OldContent, pf.NewContent)
	switch pf.Action {
	case "create":
		return diff.CountLines(pf.NewContent), hunks
	case "modify":
		return diff.CountAdditions(hunks), hunks
	default:
		// Unknown action — keep parity with the create branch as the
		// safe fallback; the spec restricts action to {create, modify,
		// delete} (Spec §354) so this branch is unreachable today.
		return diff.CountLines(pf.NewContent), hunks
	}
}

// toCLIHunks copies driving.Hunk values into the CLI hunk wire-type.
// The field-level JSON tags are identical (T0-(l)); the copy is a
// schicht-separation guarantee, not a re-shape.
func toCLIHunks(src []driving.Hunk) []hunk {
	if len(src) == 0 {
		return nil
	}
	out := make([]hunk, len(src))
	for i, h := range src {
		out[i] = hunk{
			OldStart: h.OldStart,
			OldLines: h.OldLines,
			NewStart: h.NewStart,
			NewLines: h.NewLines,
			Content:  h.Content,
		}
	}
	return out
}
