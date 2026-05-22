package managedblock_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application/managedblock"
)

func TestMarker_BeginEnd_PerStyle(t *testing.T) {
	cases := []struct {
		name      string
		marker    managedblock.Marker
		wantBegin string
		wantEnd   string
	}{
		{
			name:      "hash style",
			marker:    managedblock.Marker{Style: managedblock.StyleHash, Name: "init"},
			wantBegin: "# BEGIN U-BOOT MANAGED BLOCK: init",
			wantEnd:   "# END U-BOOT MANAGED BLOCK: init",
		},
		{
			name:      "html-comment style",
			marker:    managedblock.Marker{Style: managedblock.StyleHTMLComment, Name: "init"},
			wantBegin: "<!-- BEGIN U-BOOT MANAGED BLOCK: init -->",
			wantEnd:   "<!-- END U-BOOT MANAGED BLOCK: init -->",
		},
		{
			name:      "double-slash style",
			marker:    managedblock.Marker{Style: managedblock.StyleDoubleSlash, Name: "postgres"},
			wantBegin: "// BEGIN U-BOOT MANAGED BLOCK: postgres",
			wantEnd:   "// END U-BOOT MANAGED BLOCK: postgres",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.marker.Begin(); got != tc.wantBegin {
				t.Errorf("Begin() = %q, want %q", got, tc.wantBegin)
			}
			if got := tc.marker.End(); got != tc.wantEnd {
				t.Errorf("End() = %q, want %q", got, tc.wantEnd)
			}
		})
	}
}

func TestStyle_String(t *testing.T) {
	cases := map[managedblock.Style]string{
		managedblock.StyleHash:        "hash",
		managedblock.StyleHTMLComment: "html-comment",
		managedblock.StyleDoubleSlash: "double-slash",
		managedblock.Style(99):        "Style(99)",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("Style(%d).String() = %q, want %q", int(s), got, want)
		}
	}
}

func TestFind_HappyPath_HashStyle(t *testing.T) {
	content := []byte("before\n# BEGIN U-BOOT MANAGED BLOCK: init\nmanaged\n# END U-BOOT MANAGED BLOCK: init\nafter\n")
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "init"}

	start, end, err := managedblock.Find(content, m)
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	region := string(content[start:end])
	want := "# BEGIN U-BOOT MANAGED BLOCK: init\nmanaged\n# END U-BOOT MANAGED BLOCK: init\n"
	if region != want {
		t.Errorf("region =\n%q\nwant:\n%q", region, want)
	}
}

func TestFind_HappyPath_HTMLCommentStyle(t *testing.T) {
	content := []byte("# title\n<!-- BEGIN U-BOOT MANAGED BLOCK: init -->\nmanaged\n<!-- END U-BOOT MANAGED BLOCK: init -->\nfooter\n")
	m := managedblock.Marker{Style: managedblock.StyleHTMLComment, Name: "init"}

	start, end, err := managedblock.Find(content, m)
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	region := string(content[start:end])
	if !strings.HasPrefix(region, "<!-- BEGIN") {
		t.Errorf("region does not start with BEGIN marker: %q", region)
	}
	if !strings.HasSuffix(region, "<!-- END U-BOOT MANAGED BLOCK: init -->\n") {
		t.Errorf("region does not end with END marker + newline: %q", region)
	}
}

func TestFind_IndentedMarker_Detected(t *testing.T) {
	// Why: spec example shows indented markers (nested under
	// `services:`). The matcher must accept leading whitespace.
	content := []byte("services:\n  # BEGIN U-BOOT MANAGED BLOCK: postgres\n  postgres:\n    image: postgres:16\n  # END U-BOOT MANAGED BLOCK: postgres\n")
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "postgres"}

	start, end, err := managedblock.Find(content, m)
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	region := string(content[start:end])
	if !strings.Contains(region, "postgres:") {
		t.Errorf("indented block missed content: %q", region)
	}
	// start should land at column 0 of the BEGIN line (consuming the
	// leading whitespace), so a splice does not leave hanging spaces.
	if content[start] != ' ' {
		t.Errorf("region must begin at column 0 (leading whitespace), got byte %q", content[start])
	}
}

func TestFind_CRLFLineEndings_Detected(t *testing.T) {
	// Why: defensive — Windows-edited templates would otherwise miss.
	content := []byte("a\r\n# BEGIN U-BOOT MANAGED BLOCK: init\r\nmanaged\r\n# END U-BOOT MANAGED BLOCK: init\r\nb\r\n")
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "init"}

	start, end, err := managedblock.Find(content, m)
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	region := string(content[start:end])
	if !strings.Contains(region, "# END U-BOOT MANAGED BLOCK: init") {
		t.Errorf("CRLF region missed END marker: %q", region)
	}
	// The trailing-newline consumer must eat both \r and \n.
	if end != len(content)-len("b\r\n") {
		t.Errorf("region end = %d, want %d (just before 'b' line)", end, len(content)-len("b\r\n"))
	}
}

func TestFind_BlockMissing_ReturnsErrBlockNotFound(t *testing.T) {
	content := []byte("just some content\nno markers here\n")
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "init"}

	_, _, err := managedblock.Find(content, m)
	if !errors.Is(err, managedblock.ErrBlockNotFound) {
		t.Errorf("Find: want ErrBlockNotFound, got %v", err)
	}
}

func TestFind_BeginWithoutEnd_ReturnsErrBlockMalformed(t *testing.T) {
	content := []byte("# BEGIN U-BOOT MANAGED BLOCK: init\nstuff\nbut no end\n")
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "init"}

	_, _, err := managedblock.Find(content, m)
	if !errors.Is(err, managedblock.ErrBlockMalformed) {
		t.Errorf("Find: want ErrBlockMalformed, got %v", err)
	}
}

func TestFind_EndBeforeBegin_TreatedAsMissing(t *testing.T) {
	// Why: an END-only file (no BEGIN) is missing as far as the
	// algorithm cares — the END after BEGIN is what counts. The
	// reverse-ordered pair is treated as "no block".
	content := []byte("# END U-BOOT MANAGED BLOCK: init\nstuff\n")
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "init"}

	_, _, err := managedblock.Find(content, m)
	if !errors.Is(err, managedblock.ErrBlockNotFound) {
		t.Errorf("Find: want ErrBlockNotFound (no BEGIN), got %v", err)
	}
}

func TestFind_MultipleNamedBlocks_TargetsCorrectOne(t *testing.T) {
	// Why: a real compose.yaml will accumulate blocks per service.
	// Find must isolate the requested name and skip others.
	content := []byte(`services:
  # BEGIN U-BOOT MANAGED BLOCK: postgres
  postgres: { image: postgres:16 }
  # END U-BOOT MANAGED BLOCK: postgres
  # BEGIN U-BOOT MANAGED BLOCK: redis
  redis: { image: redis:7 }
  # END U-BOOT MANAGED BLOCK: redis
`)
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "redis"}

	start, end, err := managedblock.Find(content, m)
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	region := string(content[start:end])
	if !strings.Contains(region, "redis:") {
		t.Errorf("region missed redis content: %q", region)
	}
	if strings.Contains(region, "postgres") {
		t.Errorf("region leaked postgres content: %q", region)
	}
}

func TestFind_BlockNameWithRegexMetachars_TreatedAsLiteral(t *testing.T) {
	// Why: defensive — a name like "foo.bar" must not be interpreted
	// as a regex (`.` is any-char). regexp.QuoteMeta handles this.
	content := []byte("# BEGIN U-BOOT MANAGED BLOCK: foo.bar\nx\n# END U-BOOT MANAGED BLOCK: foo.bar\n")
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "foo.bar"}

	if _, _, err := managedblock.Find(content, m); err != nil {
		t.Errorf("Find(foo.bar): %v", err)
	}

	// A different name that would match the regex `.` for the `.`
	// must NOT match — proves we escape.
	mWrong := managedblock.Marker{Style: managedblock.StyleHash, Name: "fooxbar"}
	if _, _, err := managedblock.Find(content, mWrong); !errors.Is(err, managedblock.ErrBlockNotFound) {
		t.Errorf("Find(fooxbar): want ErrBlockNotFound, got %v", err)
	}
}

func TestFind_LastLineEndMarker_NoTrailingNewline(t *testing.T) {
	// Why: a file that ends with the END marker (no trailing \n)
	// must still yield a clean region. end will equal len(content).
	content := []byte("# BEGIN U-BOOT MANAGED BLOCK: init\nbody\n# END U-BOOT MANAGED BLOCK: init")
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "init"}

	start, end, err := managedblock.Find(content, m)
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	if end != len(content) {
		t.Errorf("end = %d, want %d (EOF)", end, len(content))
	}
	if string(content[start:end]) != string(content) {
		t.Errorf("region should equal full content when block spans whole file")
	}
}

func TestHas_PresentAndMissing(t *testing.T) {
	present := []byte("# BEGIN U-BOOT MANAGED BLOCK: init\nx\n# END U-BOOT MANAGED BLOCK: init\n")
	missing := []byte("# Some unrelated comment\n")
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "init"}

	if !managedblock.Has(present, m) {
		t.Error("Has(present) = false, want true")
	}
	if managedblock.Has(missing, m) {
		t.Error("Has(missing) = true, want false")
	}
}

func TestReplace_PreservesContentAroundBlock(t *testing.T) {
	content := []byte("before line\n# BEGIN U-BOOT MANAGED BLOCK: init\nold body\n# END U-BOOT MANAGED BLOCK: init\nafter line\n")
	replacement := []byte("# BEGIN U-BOOT MANAGED BLOCK: init\nnew body\n# END U-BOOT MANAGED BLOCK: init\n")
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "init"}

	got, err := managedblock.Replace(content, m, replacement)
	if err != nil {
		t.Fatalf("Replace: %v", err)
	}
	want := "before line\n# BEGIN U-BOOT MANAGED BLOCK: init\nnew body\n# END U-BOOT MANAGED BLOCK: init\nafter line\n"
	if string(got) != want {
		t.Errorf("Replace =\n%q\nwant:\n%q", got, want)
	}
}

func TestReplace_BlockMissing_PropagatesErr(t *testing.T) {
	content := []byte("no block here\n")
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "init"}

	_, err := managedblock.Replace(content, m, []byte("anything"))
	if !errors.Is(err, managedblock.ErrBlockNotFound) {
		t.Errorf("Replace: want ErrBlockNotFound, got %v", err)
	}
}

func TestReplace_DoesNotMutateInputContent(t *testing.T) {
	content := []byte("# BEGIN U-BOOT MANAGED BLOCK: init\nold\n# END U-BOOT MANAGED BLOCK: init\n")
	original := string(content)
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "init"}

	if _, err := managedblock.Replace(content, m, []byte("xyz")); err != nil {
		t.Fatalf("Replace: %v", err)
	}
	if string(content) != original {
		t.Errorf("Replace mutated input: got %q, want %q", content, original)
	}
}

func TestReplace_BlockMalformed_PropagatesErr(t *testing.T) {
	content := []byte("# BEGIN U-BOOT MANAGED BLOCK: init\nopen but never closed\n")
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "init"}

	_, err := managedblock.Replace(content, m, []byte("x"))
	if !errors.Is(err, managedblock.ErrBlockMalformed) {
		t.Errorf("Replace: want ErrBlockMalformed, got %v", err)
	}
}

func TestFind_BeginNotAtEndOfLine_NotFound(t *testing.T) {
	// Why: review finding #3 — spec §2099 shows markers on separate
	// lines, and the BEGIN regex anchors to `$`. A single-line
	// `BEGIN…--><…END…-->` has the END text appended after BEGIN, so
	// the BEGIN regex never matches → ErrBlockNotFound (the stricter
	// "no BEGIN here" answer rather than "malformed body").
	content := []byte("<!-- BEGIN U-BOOT MANAGED BLOCK: init --><!-- END U-BOOT MANAGED BLOCK: init -->\n")
	m := managedblock.Marker{Style: managedblock.StyleHTMLComment, Name: "init"}

	_, _, err := managedblock.Find(content, m)
	if !errors.Is(err, managedblock.ErrBlockNotFound) {
		t.Errorf("Find: want ErrBlockNotFound for invalid single-line markers, got %v", err)
	}
}

func TestFind_EndAppendedMidLine_Malformed(t *testing.T) {
	// Why: complements the single-line case — BEGIN cleanly matches
	// on its own line, but END is concatenated after body text on a
	// later line. The END regex requires `$` so it never matches →
	// BEGIN-without-END → ErrBlockMalformed (covers the
	// skipLineEnding path explicitly).
	content := []byte("# BEGIN U-BOOT MANAGED BLOCK: init\nbody # END U-BOOT MANAGED BLOCK: init trailing\n")
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "init"}

	_, _, err := managedblock.Find(content, m)
	if !errors.Is(err, managedblock.ErrBlockMalformed) {
		t.Errorf("Find: want ErrBlockMalformed for END mid-line, got %v", err)
	}
}

func TestFind_DuplicatedBeginMarker_RejectedAsMalformed(t *testing.T) {
	// Why: review finding #4 — a botched manual edit can leave two
	// BEGIN markers before the END. Silent auto-repair would let
	// Replace absorb both into the "managed body"; pin the rejection.
	content := []byte("# BEGIN U-BOOT MANAGED BLOCK: init\nfirst body\n# BEGIN U-BOOT MANAGED BLOCK: init\nsecond body\n# END U-BOOT MANAGED BLOCK: init\n")
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "init"}

	_, _, err := managedblock.Find(content, m)
	if !errors.Is(err, managedblock.ErrBlockMalformed) {
		t.Errorf("Find: want ErrBlockMalformed for duplicated BEGIN, got %v", err)
	}
}

func TestHas_MalformedBlock_ReturnsFalse(t *testing.T) {
	// Why: pin the safe-direction behaviour — a half-edited block
	// (BEGIN-only) makes Has return false, so re-init falls into the
	// "no block" branch (ErrProjectExists / ErrForceRequiresBackup)
	// instead of attempting a Replace that would itself error.
	content := []byte("# BEGIN U-BOOT MANAGED BLOCK: init\nopen forever\n")
	m := managedblock.Marker{Style: managedblock.StyleHash, Name: "init"}

	if managedblock.Has(content, m) {
		t.Errorf("Has(malformed) = true, want false")
	}
}

func TestInitName_MatchesTemplateConvention(t *testing.T) {
	// Why: lock the package-level constant to the literal "init"
	// that every M3 template embeds in its BEGIN/END markers. A
	// rename here without a template update would silently break
	// re-init detection.
	if managedblock.InitName != "init" {
		t.Errorf("InitName = %q, want %q", managedblock.InitName, "init")
	}
}
