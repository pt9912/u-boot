package application_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// messageUserPatterns match identifiers (requirements, specification,
// decisions, slices, tranches), German words, announcements and history
// remarks. A message of the binary describes what is, in words the
// user reads without a lookup.
func messageUserPatterns() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`\bLH-[A-Z]+(-[A-Z0-9]+)*-?[0-9]+[A-Z]?\b`),
		regexp.MustCompile(`\bLH-AK-[0-9]+\b`),
		regexp.MustCompile(`\b(SPEC|ADR|ARC)-[0-9]+\b`),
		regexp.MustCompile(`\bslice-[a-z0-9][a-z0-9.-]*`),
		regexp.MustCompile(`§`),
		regexp.MustCompile(`\bT[0-9]+-\(`),
		regexp.MustCompile(`(?i)\b(voll-schema|fachlich\w*)\b`),
		regexp.MustCompile(`\bComing\b`),
		regexp.MustCompile(`(?i)\bcoming soon\b`),
		regexp.MustCompile(`will follow`),
		regexp.MustCompile(`\bjoin in\b`),
		regexp.MustCompile(`(?i)\bnot yet\b`),
		regexp.MustCompile(`\bdeferred\b`),
		regexp.MustCompile(`\bplanned for\b`),
		regexp.MustCompile(`\bTODO\b`),
		regexp.MustCompile(`\bV[0-9]+\b`),
		regexp.MustCompile(`\bv[0-9]+\.[0-9]+(\.[0-9]+)?\b`),
		regexp.MustCompile(`\btoday\b`),
		regexp.MustCompile(`\bMVP\b`),
		regexp.MustCompile(`\bpreviously\b`),
		regexp.MustCompile(`\bformerly\b`),
		regexp.MustCompile(`\bused to\b`),
		regexp.MustCompile(`\bwas (removed|renamed|replaced)\b`),
		regexp.MustCompile(`deprecated since`),
		regexp.MustCompile(`\bsince v?[0-9]`),
		regexp.MustCompile(`\bchanged in\b`),
		regexp.MustCompile(`\bin v[0-9]`),
		regexp.MustCompile(`\b(früher|bisher|vorher|seit)\b`),
	}
}

func messageFindings(text string) []string {
	var found []string
	for _, re := range messageUserPatterns() {
		if m := re.FindString(text); m != "" {
			found = append(found, m)
		}
	}
	return found
}

func TestMessageFindings_Patterns(t *testing.T) {
	for _, s := range []string{
		`%q does not match %s (LH-FA-INIT-006)`,
		`the LH-FA-ADD-005 state machine`,
		`see slice-v0.1.1-doctor-container-awareness`,
		`volume removal is deferred (v0.3.0)`,
		`LH-SA-FILE-002 §611`,
	} {
		if len(messageFindings(s)) == 0 {
			t.Errorf("no pattern matches %q", s)
		}
	}
	for _, s := range []string{
		`remove them after confirming the data is no longer needed`,
		`the profile no longer isolates the container from its host`,
		`a later pair may depend on an earlier one`,
	} {
		if f := messageFindings(s); len(f) != 0 {
			t.Errorf("%q must stay green, matched %v", s, f)
		}
	}
}

// messageCallees are the calls whose string arguments become text the
// user reads (error values, formatted messages, written output).
func messageCallees() map[string]bool {
	return map[string]bool{
		"errors.New": true, "fmt.Errorf": true, "fmt.Sprintf": true,
		"fmt.Fprintf": true, "fmt.Fprintln": true, "fmt.Fprint": true,
	}
}

// messageFields are the struct fields whose string value is a message.
func messageFields() map[string]bool {
	return map[string]bool{"Message": true, "Hint": true, "MissingHint": true, "Detail": true}
}

// messageNameWords mark a constant or variable that holds message text.
var messageNameWords = regexp.MustCompile(`(?i)(hint|message|msg|prose|reason)`)

// calleeName renders "pkg.Func" for a selector call, or the plain
// function name for a local call.
func calleeName(call *ast.CallExpr) string {
	switch fn := call.Fun.(type) {
	case *ast.SelectorExpr:
		if id, ok := fn.X.(*ast.Ident); ok {
			return id.Name + "." + fn.Sel.Name
		}
	case *ast.Ident:
		return fn.Name
	}
	return ""
}

// stringLiterals collects the string literals below n.
func stringLiterals(n ast.Node, fset *token.FileSet, report func(pos token.Pos, text string)) {
	ast.Inspect(n, func(c ast.Node) bool {
		lit, ok := c.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}
		if val, err := strconv.Unquote(lit.Value); err == nil {
			report(lit.Pos(), val)
		}
		return true
	})
}

// TestMessageLiteralsCarryNoIdentifiersOrHistory (seit slice-v1-hilfetexte-ohne-kennungen)
// walks the non-test Go
// sources of the driving CLI adapter, the application, the ports, the
// domain and the driven adapters and reports every string literal that
// forms message text — an argument of errors.New/fmt.Errorf/fmt.Sprintf/
// fmt.Fprint*, the second argument of the sandbox helpers warn/relax,
// the value of a Message/Hint/MissingHint/Detail field, or the value of
// a constant or variable named like a message — and carries a pattern.
// The Code fields and the first argument of warn/relax are codes, not
// text, and stay out.
func TestMessageLiteralsCarryNoIdentifiersOrHistory(t *testing.T) {
	root := configFilenameRepoRoot(t)
	fset := token.NewFileSet()
	scanned := 0
	var offenders []string
	callees, fields := messageCallees(), messageFields()

	roots := []string{
		"internal/adapter/driving/cli", "internal/adapter/driven",
		"internal/hexagon/application", "internal/hexagon/port",
		"internal/hexagon/domain",
	}
	for _, top := range roots {
		walkErr := filepath.WalkDir(filepath.Join(root, top), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == "jsontestutil" || d.Name() == "testdata" {
					return fs.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			file, parseErr := parser.ParseFile(fset, path, nil, 0)
			if parseErr != nil {
				return parseErr
			}
			scanned++
			report := func(pos token.Pos, text string) {
				for _, f := range messageFindings(text) {
					offenders = append(offenders, fset.Position(pos).String()+": "+f)
				}
			}
			ast.Inspect(file, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.CallExpr:
					name := calleeName(x)
					switch {
					case callees[name]:
						for _, a := range x.Args {
							stringLiterals(a, fset, report)
						}
					case (name == "warn" || name == "relax") && len(x.Args) == 2:
						stringLiterals(x.Args[1], fset, report)
					}
				case *ast.KeyValueExpr:
					if id, ok := x.Key.(*ast.Ident); ok && fields[id.Name] {
						stringLiterals(x.Value, fset, report)
					}
				case *ast.ValueSpec:
					for i, id := range x.Names {
						if messageNameWords.MatchString(id.Name) && i < len(x.Values) {
							stringLiterals(x.Values[i], fset, report)
						}
					}
				}
				return true
			})
			return nil
		})
		if walkErr != nil {
			t.Fatalf("walking %s: %v", top, walkErr)
		}
	}

	if scanned < 50 {
		t.Fatalf("only %d non-test Go files scanned — root detection or walk broken", scanned)
	}
	seen := map[string]bool{}
	for _, o := range offenders {
		if !seen[o] {
			seen[o] = true
			t.Errorf("message text with identifier, announcement or history: %s", o)
		}
	}
}
