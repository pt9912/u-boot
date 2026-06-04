package jsontestutil_test

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli/jsontestutil"
	"github.com/pt9912/u-boot/internal/hexagon/application"
)

// repoRoot returns the repository root by walking up from this test
// file's directory until it finds go.mod.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, here, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(here)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found while walking up")
		}
		dir = parent
	}
}

// TestDriftGate1_MapVsDoctorCheckIDs is the T0-(h) Gate 1: every
// check ID returned by application.DoctorCheckIDs() must have an
// entry in DefaultAllowedCodes.
//
// Previously this test source-parsed application/doctor.go with
// go/parser (depguard-loophole). The application package now ships
// an explicit DoctorCheckIDs() public helper — the gate is a clean
// import (Review M3-Findings adressiert).
func TestDriftGate1_MapVsDoctorCheckIDs(t *testing.T) {
	checkIDs := application.DoctorCheckIDs()
	if len(checkIDs) == 0 {
		t.Fatal("application.DoctorCheckIDs() returned empty — helper broken or doctor.go moved")
	}

	registry := jsontestutil.DefaultAllowedCodes()
	for _, code := range checkIDs {
		if _, ok := registry[code]; !ok {
			t.Errorf("check ID %q from application.DoctorCheckIDs() missing in DefaultAllowedCodes (T0-(h) Gate 1)", code)
		}
	}
}

// TestDriftGate2_MapVsMarkdownDoc is the T0-(h) Gate 2: every Map
// entry must have a Markdown table row in docs/user/cli-json-output.md
// §5 Code-Registry, and vice versa. Both drift directions are checked.
//
// Sektion-Begrenzung via HTML-Markers `<!-- code-registry:start -->`
// und `<!-- code-registry:end -->` in der Doku — robust gegen
// Doku-Erweiterungen um weitere Tabellen in anderen Sektionen
// (Review M1-Findings adressiert).
func TestDriftGate2_MapVsMarkdownDoc(t *testing.T) {
	docPath := filepath.Join(repoRoot(t), "docs", "user", "cli-json-output.md")
	bytes, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("read cli-json-output.md: %v", err)
	}

	section := extractRegistrySection(string(bytes))
	if section == "" {
		t.Fatal("code-registry markers not found in cli-json-output.md — Doku-Anker entfernt?")
	}

	docCodes := extractMarkdownCodes(section)
	if len(docCodes) == 0 {
		t.Fatal("found zero Markdown table rows in code-registry section — extractor broken")
	}

	registry := jsontestutil.DefaultAllowedCodes()

	for code := range registry {
		if _, ok := docCodes[code]; !ok {
			t.Errorf("registry code %q missing from Markdown doc code-registry section (T0-(h) Gate 2, map → doc)", code)
		}
	}

	for code := range docCodes {
		if _, ok := registry[code]; !ok {
			t.Errorf("Markdown doc lists code %q absent from registry (T0-(h) Gate 2, doc → map)", code)
		}
	}
}

// extractRegistrySection isolates the registry block between the
// `<!-- code-registry:start -->` and `<!-- code-registry:end -->`
// markers. Returns "" if either marker is missing.
func extractRegistrySection(content string) string {
	const startMarker = "<!-- code-registry:start -->"
	const endMarker = "<!-- code-registry:end -->"
	start := strings.Index(content, startMarker)
	if start < 0 {
		return ""
	}
	end := strings.Index(content, endMarker)
	if end < 0 || end < start {
		return ""
	}
	return content[start+len(startMarker) : end]
}

// extractMarkdownCodes scans the Markdown for table rows shaped like
// `| `code` | description |` and returns a set of the codes. The
// regex is liberal on surrounding whitespace.
var markdownRowRe = regexp.MustCompile(`^\s*\|\s*` + "`" + `([^` + "`" + `]+)` + "`" + `\s*\|\s*[^|]+\|\s*$`)

func extractMarkdownCodes(content string) map[string]bool {
	codes := make(map[string]bool)
	for _, line := range strings.Split(content, "\n") {
		m := markdownRowRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		code := strings.TrimSpace(m[1])
		if !isCodeLike(code) {
			continue
		}
		codes[code] = true
	}
	return codes
}

// isCodeLike accepts entries that look like dotted-identifier codes
// (e.g. `docker.installed`). LH-IDs are NOT accepted here because
// the code-registry section documents Tool-internal codes only;
// LH-IDs are handled at runtime by the helper, not in the doc.
func isCodeLike(s string) bool {
	if s == "" {
		return false
	}
	if !strings.Contains(s, ".") {
		return false
	}
	for _, r := range s {
		if r != '.' && r != '-' && r != '_' &&
			(r < 'a' || r > 'z') &&
			(r < 'A' || r > 'Z') &&
			(r < '0' || r > '9') {
			return false
		}
	}
	return true
}
