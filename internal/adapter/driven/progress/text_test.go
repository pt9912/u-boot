package progress_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/progress"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

func TestText_AffectedFiles_RendersHeaderAndRows(t *testing.T) {
	var buf bytes.Buffer
	w := progress.NewText(&buf)

	w.AffectedFiles("/proj", []driven.AffectedFile{
		{Path: "compose.yaml", Action: driven.AffectedReplaceBlock, Backup: false},
		{Path: ".env.example", Action: driven.AffectedOverwriteFull, Backup: true},
	})

	got := buf.String()
	wantLines := []string{
		"Affected files in /proj:",
		"  - compose.yaml — replace managed block",
		"  - .env.example — full overwrite (with backup)",
	}
	for _, line := range wantLines {
		if !strings.Contains(got, line) {
			t.Errorf("output missing %q; got:\n%s", line, got)
		}
	}
}

func TestText_AffectedFiles_BackupMarkerOnReplaceBlock(t *testing.T) {
	// Why: --force + --backup on a managed-block file gets a safety
	// copy; the marker must show even when the action is block-only.
	var buf bytes.Buffer
	w := progress.NewText(&buf)

	w.AffectedFiles("/p", []driven.AffectedFile{
		{Path: "README.md", Action: driven.AffectedReplaceBlock, Backup: true},
	})

	if !strings.Contains(buf.String(), "replace managed block (with backup)") {
		t.Errorf("backup marker missing on replace-block row: %q", buf.String())
	}
}

func TestText_AffectedFiles_UsesEmDash(t *testing.T) {
	// Why: convention pin — the action label is separated from the
	// path by an em-dash (matches the existing M3-T4b user-visible
	// format). ASCII-dash drift would surface here first.
	var buf bytes.Buffer
	w := progress.NewText(&buf)
	w.AffectedFiles("/p", []driven.AffectedFile{
		{Path: "x", Action: driven.AffectedReplaceBlock},
	})
	if !strings.Contains(buf.String(), "—") {
		t.Errorf("expected em-dash separator in output: %q", buf.String())
	}
}

func TestText_AffectedFiles_UnknownActionFallsBackToInt(t *testing.T) {
	// Why: defensive — an enum extension that ships before this
	// adapter is updated should produce a debug-friendly label
	// instead of an empty string.
	var buf bytes.Buffer
	w := progress.NewText(&buf)
	w.AffectedFiles("/p", []driven.AffectedFile{
		{Path: "x", Action: driven.AffectedAction(99)},
	})
	if !strings.Contains(buf.String(), "action(99)") {
		t.Errorf("expected fallback label for unknown action; got %q", buf.String())
	}
}
