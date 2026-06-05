package diff_test

import (
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli/diff"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TestCountLines_TrailingNewlineEdgeCases pins the four cases the
// slice plan (T0-(g)) explicitly calls out — `bytes.Count("\n")` plus
// `HasSuffix` correction.
func TestCountLines_TrailingNewlineEdgeCases(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    int
	}{
		{name: "empty", content: "", want: 0},
		{name: "single line without newline", content: "a", want: 1},
		{name: "single line with newline", content: "a\n", want: 1},
		{name: "two lines no trailing", content: "a\nb", want: 2},
		{name: "two lines with trailing", content: "a\nb\n", want: 2},
		{name: "twelve lines fresh-postgres-fixture", content: strings.Repeat("x\n", 12), want: 12},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := diff.CountLines([]byte(tc.content))
			if got != tc.want {
				t.Fatalf("CountLines(%q) = %d, want %d", tc.content, got, tc.want)
			}
		})
	}
}

func TestIsBinary(t *testing.T) {
	cases := []struct {
		name       string
		oldC, newC []byte
		wantBinary bool
	}{
		{name: "both empty", oldC: nil, newC: nil, wantBinary: false},
		{name: "text both sides", oldC: []byte("a\nb\n"), newC: []byte("a\nB\n"), wantBinary: false},
		{name: "invalid utf-8 in old", oldC: []byte{0xff, 0xfe}, newC: []byte("x\n"), wantBinary: true},
		{name: "invalid utf-8 in new", oldC: []byte("x\n"), newC: []byte{0xff, 0xfe}, wantBinary: true},
		{name: "valid multi-byte utf-8", oldC: []byte("Ümlaut\n"), newC: []byte("Ümlaut!\n"), wantBinary: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := diff.IsBinary(tc.oldC, tc.newC); got != tc.wantBinary {
				t.Fatalf("IsBinary = %v, want %v", got, tc.wantBinary)
			}
		})
	}
}

func TestCountBytesDiff(t *testing.T) {
	cases := []struct {
		name       string
		oldC, newC []byte
		want       int
	}{
		{name: "both empty", oldC: nil, newC: nil, want: 0},
		{name: "create from empty", oldC: nil, newC: []byte("abc"), want: 3},
		{name: "delete to empty", oldC: []byte("abc"), newC: nil, want: 3},
		{name: "growth", oldC: []byte("abc"), newC: []byte("abcdef"), want: 3},
		{name: "shrink", oldC: []byte("abcdef"), newC: []byte("abc"), want: 3},
		{name: "same size different content", oldC: []byte("abc"), newC: []byte("xyz"), want: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := diff.CountBytesDiff(tc.oldC, tc.newC); got != tc.want {
				t.Fatalf("CountBytesDiff = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestCountFromHunks(t *testing.T) {
	hunks := []driving.Hunk{
		{NewLines: 3},
		{NewLines: 0},
		{NewLines: 5},
	}
	if got := diff.CountFromHunks(hunks); got != 8 {
		t.Fatalf("CountFromHunks = %d, want 8", got)
	}
	if got := diff.CountFromHunks(nil); got != 0 {
		t.Fatalf("CountFromHunks(nil) = %d, want 0", got)
	}
}

func TestCompute_IdenticalIsNil(t *testing.T) {
	hunks := diff.Compute([]byte("a\nb\nc\n"), []byte("a\nb\nc\n"))
	if hunks != nil {
		t.Fatalf("identical inputs: expected nil hunks, got %#v", hunks)
	}
}

func TestCompute_BinaryIsNil(t *testing.T) {
	hunks := diff.Compute([]byte{0xff}, []byte("text\n"))
	if hunks != nil {
		t.Fatalf("binary input: expected nil hunks, got %#v", hunks)
	}
}

func TestCompute_BothEmpty(t *testing.T) {
	hunks := diff.Compute(nil, nil)
	if hunks != nil {
		t.Fatalf("both empty: expected nil hunks, got %#v", hunks)
	}
}

func TestCompute_PureAddition(t *testing.T) {
	hunks := diff.Compute(nil, []byte("a\nb\n"))
	if len(hunks) != 1 {
		t.Fatalf("pure addition: want 1 hunk, got %d (%#v)", len(hunks), hunks)
	}
	h := hunks[0]
	if h.OldStart != 0 || h.OldLines != 0 {
		t.Errorf("pure addition: OldStart=%d OldLines=%d, want 0/0", h.OldStart, h.OldLines)
	}
	if h.NewStart != 1 || h.NewLines != 2 {
		t.Errorf("pure addition: NewStart=%d NewLines=%d, want 1/2", h.NewStart, h.NewLines)
	}
	if h.Content != "+a\n+b\n" {
		t.Errorf("pure addition: Content = %q, want %q", h.Content, "+a\n+b\n")
	}
}

func TestCompute_PureDeletion(t *testing.T) {
	hunks := diff.Compute([]byte("a\nb\n"), nil)
	if len(hunks) != 1 {
		t.Fatalf("pure deletion: want 1 hunk, got %d (%#v)", len(hunks), hunks)
	}
	h := hunks[0]
	if h.OldStart != 1 || h.OldLines != 2 {
		t.Errorf("pure deletion: OldStart=%d OldLines=%d, want 1/2", h.OldStart, h.OldLines)
	}
	if h.NewStart != 0 || h.NewLines != 0 {
		t.Errorf("pure deletion: NewStart=%d NewLines=%d, want 0/0", h.NewStart, h.NewLines)
	}
	if h.Content != "-a\n-b\n" {
		t.Errorf("pure deletion: Content = %q, want %q", h.Content, "-a\n-b\n")
	}
}

func TestCompute_MiddleModify(t *testing.T) {
	hunks := diff.Compute([]byte("a\nb\nc\n"), []byte("a\nB\nc\n"))
	if len(hunks) != 1 {
		t.Fatalf("middle modify: want 1 hunk, got %d (%#v)", len(hunks), hunks)
	}
	h := hunks[0]
	if h.OldStart != 1 || h.OldLines != 3 {
		t.Errorf("middle modify: OldStart=%d OldLines=%d, want 1/3", h.OldStart, h.OldLines)
	}
	if h.NewStart != 1 || h.NewLines != 3 {
		t.Errorf("middle modify: NewStart=%d NewLines=%d, want 1/3", h.NewStart, h.NewLines)
	}
	const want = " a\n-b\n+B\n c\n"
	if h.Content != want {
		t.Errorf("middle modify: Content = %q, want %q", h.Content, want)
	}
}

func TestCompute_TrailingAppend(t *testing.T) {
	// Adding lines at the end — single hunk, no leading context
	// available beyond the entire existing file (3 lines ≤ default
	// context of 3, so all three precede the change).
	hunks := diff.Compute([]byte("a\nb\nc\n"), []byte("a\nb\nc\nd\n"))
	if len(hunks) != 1 {
		t.Fatalf("trailing append: want 1 hunk, got %d (%#v)", len(hunks), hunks)
	}
	h := hunks[0]
	if h.NewLines != 4 || h.OldLines != 3 {
		t.Errorf("trailing append: NewLines=%d OldLines=%d, want 4/3", h.NewLines, h.OldLines)
	}
	if !strings.HasSuffix(h.Content, "+d\n") {
		t.Errorf("trailing append: Content should end with +d, got %q", h.Content)
	}
}

func TestCompute_TwoSeparateHunks(t *testing.T) {
	// Changes at positions 1 and 11 — with 10 equal lines between
	// and default context 3, the gap (10) > 2*context (6), so they
	// stay as two separate hunks.
	oldLines := []string{"X", "a", "a", "a", "a", "a", "a", "a", "a", "a", "a", "Y"}
	newLines := []string{"x", "a", "a", "a", "a", "a", "a", "a", "a", "a", "a", "y"}
	oldC := []byte(strings.Join(oldLines, "\n") + "\n")
	newC := []byte(strings.Join(newLines, "\n") + "\n")
	hunks := diff.Compute(oldC, newC)
	if len(hunks) != 2 {
		t.Fatalf("two changes far apart: want 2 hunks, got %d (%#v)", len(hunks), hunks)
	}
}

func TestCompute_MergedHunksWithinContext(t *testing.T) {
	// Changes at positions 1 and 5 — only 3 equal lines between.
	// 3 ≤ 2*context (6), so they merge into one hunk.
	oldC := []byte("X\na\nb\nc\nY\nd\n")
	newC := []byte("x\na\nb\nc\ny\nd\n")
	hunks := diff.Compute(oldC, newC)
	if len(hunks) != 1 {
		t.Fatalf("two changes close together: want 1 merged hunk, got %d (%#v)", len(hunks), hunks)
	}
}

func TestRender_EmptyHunksIsEmpty(t *testing.T) {
	if got := diff.Render(nil); got != "" {
		t.Errorf("Render(nil) = %q, want empty", got)
	}
}

func TestRender_HeaderFormat(t *testing.T) {
	hunks := []driving.Hunk{
		{OldStart: 5, OldLines: 3, NewStart: 5, NewLines: 4, Content: " a\n-b\n+B\n c\n+d\n"},
	}
	got := diff.Render(hunks)
	const wantHeader = "@@ -5,3 +5,4 @@\n"
	if !strings.HasPrefix(got, wantHeader) {
		t.Errorf("Render: missing header %q in output %q", wantHeader, got)
	}
	if !strings.Contains(got, " a\n-b\n+B\n c\n+d\n") {
		t.Errorf("Render: missing hunk body in output %q", got)
	}
}

func TestRender_MultipleHunksConcatenated(t *testing.T) {
	hunks := []driving.Hunk{
		{OldStart: 1, OldLines: 1, NewStart: 1, NewLines: 1, Content: "-a\n+A\n"},
		{OldStart: 10, OldLines: 1, NewStart: 10, NewLines: 1, Content: "-z\n+Z\n"},
	}
	got := diff.Render(hunks)
	want := "@@ -1,1 +1,1 @@\n-a\n+A\n@@ -10,1 +10,1 @@\n-z\n+Z\n"
	if got != want {
		t.Errorf("Render multi-hunk:\n got %q\nwant %q", got, want)
	}
}

// TestCompute_PostgresComposeFresh pins the create-from-scratch
// scenario from Slice §Aufhebungsbedingung Variante A: a fresh
// compose.yaml gets a 12-line postgres block. Result must yield
// exactly one hunk with NewLines == 12 so the consumer-side
// CountFromHunks reports `count: 12` (Spec §430).
func TestCompute_PostgresComposeFresh(t *testing.T) {
	const block = `services:
  postgres:
    image: postgres:16
    restart: unless-stopped
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
`
	hunks := diff.Compute(nil, []byte(block))
	if len(hunks) != 1 {
		t.Fatalf("postgres fresh: want 1 hunk, got %d", len(hunks))
	}
	if got := hunks[0].NewLines; got != 12 {
		t.Errorf("postgres fresh: NewLines = %d, want 12 (Spec §430)", got)
	}
	if got := diff.CountFromHunks(hunks); got != 12 {
		t.Errorf("postgres fresh: CountFromHunks = %d, want 12 (Spec §430)", got)
	}
}

// TestCompute_PostgresComposeExisting pins the modify-existing
// scenario from Slice §Aufhebungsbedingung Variante B: an existing
// compose.yaml with another service gets 6 lines appended. The
// formal T0-(g) definition is `count = sum(hunk.NewLines)` — context
// lines count too — so NewLines > 6 is correct here; the floor
// invariant is that the 6 actual additions show up as `+` lines.
func TestCompute_PostgresComposeExisting(t *testing.T) {
	const existing = `services:
  redis:
    image: redis:7
    restart: unless-stopped
`
	const updated = `services:
  redis:
    image: redis:7
    restart: unless-stopped
  postgres:
    image: postgres:16
    restart: unless-stopped
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
    ports:
      - "5432:5432"
`
	hunks := diff.Compute([]byte(existing), []byte(updated))
	if len(hunks) != 1 {
		t.Fatalf("postgres existing: want 1 hunk, got %d", len(hunks))
	}
	if got := diff.CountFromHunks(hunks); got <= 6 {
		t.Errorf("postgres existing: CountFromHunks = %d, want > 6 (NewLines = inserts + context)", got)
	}
	additions := strings.Count(hunks[0].Content, "\n+")
	if additions < 6 {
		t.Errorf("postgres existing: counted %d '+' lines, want ≥ 6", additions)
	}
}
