package jsontestutil_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli/jsontestutil"
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
// checkID* constant declared in internal/hexagon/application/doctor.go
// must have an entry in DefaultAllowedCodes. If a future doctor
// check lands without a registry entry, this test breaks.
//
// Parses doctor.go via go/parser (stdlib, no new dep). We look for
// any const declaration whose name starts with "checkID" and reads
// the assigned string literal.
func TestDriftGate1_MapVsDoctorCheckIDs(t *testing.T) {
	doctorPath := filepath.Join(repoRoot(t), "internal", "hexagon", "application", "doctor.go")
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, doctorPath, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse doctor.go: %v", err)
	}

	checkIDs := extractCheckIDs(file)
	if len(checkIDs) == 0 {
		t.Fatal("found zero checkID* constants in doctor.go — extractor broken or file moved")
	}

	registry := jsontestutil.DefaultAllowedCodes()
	for _, code := range checkIDs {
		if _, ok := registry[code]; !ok {
			t.Errorf("checkID %q from doctor.go missing in DefaultAllowedCodes (T0-(h) Gate 1)", code)
		}
	}
}

// extractCheckIDs walks the AST for `const checkID... = "literal"`
// declarations and returns the literal values.
func extractCheckIDs(file *ast.File) []string {
	var ids []string
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, name := range vs.Names {
				if !strings.HasPrefix(name.Name, "checkID") {
					continue
				}
				if i >= len(vs.Values) {
					continue
				}
				lit, ok := vs.Values[i].(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					continue
				}
				unquoted := strings.Trim(lit.Value, `"`)
				ids = append(ids, unquoted)
			}
		}
	}
	return ids
}

// TestDriftGate2_MapVsMarkdownDoc is the T0-(h) Gate 2: every Map
// entry must have a Markdown table row in docs/user/cli-json-output.md,
// and vice versa. Bricht in both drift directions.
//
// Markdown row format expected (from docs/user/cli-json-output.md §5.1):
//   | `code.name` | description text |
//
// Code-Spalte ist immer Backtick-Quoted, Description folgt nach `|`.
func TestDriftGate2_MapVsMarkdownDoc(t *testing.T) {
	docPath := filepath.Join(repoRoot(t), "docs", "user", "cli-json-output.md")
	bytes, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("read cli-json-output.md: %v", err)
	}

	docCodes := extractMarkdownCodes(string(bytes))
	if len(docCodes) == 0 {
		t.Fatal("found zero Markdown table rows — extractor broken or doc moved")
	}

	registry := jsontestutil.DefaultAllowedCodes()

	for code := range registry {
		if _, ok := docCodes[code]; !ok {
			t.Errorf("registry code %q missing from Markdown doc §5 (T0-(h) Gate 2, map → doc)", code)
		}
	}

	for code := range docCodes {
		if _, ok := registry[code]; !ok {
			t.Errorf("Markdown doc lists code %q absent from registry (T0-(h) Gate 2, doc → map)", code)
		}
	}
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
		// Skip non-code rows: schema/spec markers, table separators,
		// description-headers. Real codes follow dotted-identifier or
		// LH-prefix shape; everything else gets filtered.
		if !isCodeLike(code) {
			continue
		}
		codes[code] = true
	}
	return codes
}

// isCodeLike accepts entries that look like dotted-identifier
// codes (e.g. `docker.installed`) or LH-IDs. Rejects code-blocks,
// inline-spans of regular prose, paths, file-names.
func isCodeLike(s string) bool {
	if s == "" {
		return false
	}
	if strings.HasPrefix(s, "LH-") {
		return true
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
