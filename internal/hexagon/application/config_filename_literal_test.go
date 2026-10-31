package application_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

const configFileNameLiteral = "u-boot.yaml"

// configFileNameHome is the one non-test file that spells the
// configuration file name as a string literal.
const configFileNameHome = "internal/hexagon/domain/config.go"

// configFilenameRepoRoot walks up from this test file to the
// directory holding go.mod.
func configFilenameRepoRoot(t *testing.T) string {
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

// TestConfigFileNameLiteralOnlyInDomain pins that no non-test Go file
// under internal/ and cmd/ carries a string literal equal to the
// configuration file name, except internal/hexagon/domain/config.go.
// Message texts and comments that merely mention the name are not
// string literals of exactly that value and stay out of scope.
func TestConfigFileNameLiteralOnlyInDomain(t *testing.T) {
	root := configFilenameRepoRoot(t)
	fset := token.NewFileSet()
	scanned := 0
	var offenders []string

	for _, top := range []string{"internal", "cmd"} {
		walkErr := filepath.WalkDir(filepath.Join(root, top), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			rel, relErr := filepath.Rel(root, path)
			if relErr != nil {
				return relErr
			}
			rel = filepath.ToSlash(rel)
			scanned++
			if rel == configFileNameHome {
				return nil
			}
			file, parseErr := parser.ParseFile(fset, path, nil, 0)
			if parseErr != nil {
				return parseErr
			}
			ast.Inspect(file, func(n ast.Node) bool {
				lit, ok := n.(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					return true
				}
				if val, uqErr := strconv.Unquote(lit.Value); uqErr == nil && val == configFileNameLiteral {
					offenders = append(offenders, fset.Position(lit.Pos()).String())
				}
				return true
			})
			return nil
		})
		if walkErr != nil {
			t.Fatalf("walking %s: %v", top, walkErr)
		}
	}

	if scanned == 0 {
		t.Fatal("no non-test Go files scanned — root detection or walk broken")
	}
	for _, pos := range offenders {
		t.Errorf("string literal %q at %s: use domain.ConfigFileName (only %s spells the name)", configFileNameLiteral, pos, configFileNameHome)
	}
}
