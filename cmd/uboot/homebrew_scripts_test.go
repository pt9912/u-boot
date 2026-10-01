package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// slice-v2-homebrew-formula: the formula fill script substitutes tag,
// version and the four platform digests, and fails closed (no output
// file) on a missing or malformed digest.
const sumsFull = `1111111111111111111111111111111111111111111111111111111111111111  u-boot-darwin-amd64
2222222222222222222222222222222222222222222222222222222222222222  u-boot-darwin-arm64
3333333333333333333333333333333333333333333333333333333333333333  u-boot-linux-amd64
4444444444444444444444444444444444444444444444444444444444444444  u-boot-linux-arm64
5555555555555555555555555555555555555555555555555555555555555555  u-boot-windows-amd64.exe
`

func runFill(t *testing.T, tag, sums string) (string, string, error) {
	t.Helper()
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("bash not available")
	}
	dir := t.TempDir()
	sumsPath := filepath.Join(dir, "SHA256SUMS")
	if err := os.WriteFile(sumsPath, []byte(sums), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "u-boot.rb")
	cmd := exec.Command("bash", "../../scripts/homebrew-formula-fill.sh", tag, sumsPath, "../../scripts/homebrew-formula.rb.tmpl", out)
	msg, err := cmd.CombinedOutput()
	body, _ := os.ReadFile(out)
	return string(body), string(msg), err
}

func TestHomebrewFormulaFill_Fills(t *testing.T) {
	body, msg, err := runFill(t, "v1.2.3", sumsFull)
	if err != nil {
		t.Fatalf("fill: %v\n%s", err, msg)
	}
	for _, want := range []string{
		`version "1.2.3"`,
		"releases/download/v1.2.3/u-boot-darwin-arm64",
		`sha256 "1111111111111111111111111111111111111111111111111111111111111111"`, // darwin-amd64
		`sha256 "2222222222222222222222222222222222222222222222222222222222222222"`, // darwin-arm64
		`sha256 "3333333333333333333333333333333333333333333333333333333333333333"`, // linux-amd64
		`sha256 "4444444444444444444444444444444444444444444444444444444444444444"`, // linux-arm64
		"class UBoot < Formula",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("formula lacks %q", want)
		}
	}
	if strings.Contains(body, "__") {
		t.Errorf("unfilled placeholder left in the formula:\n%s", body)
	}
	if strings.Contains(body, "windows") {
		t.Errorf("Homebrew carries no Windows platform")
	}
}

func TestHomebrewFormulaFill_FailsClosed(t *testing.T) {
	missing := strings.Replace(sumsFull, "u-boot-linux-arm64", "u-boot-linux-armv7", 1)
	malformed := strings.Replace(sumsFull, "4444444444444444444444444444444444444444444444444444444444444444", "not-a-digest", 1)
	for name, sums := range map[string]string{"missing platform": missing, "malformed digest": malformed} {
		body, _, err := runFill(t, "v1.2.3", sums)
		if err == nil || body != "" {
			t.Errorf("%s: err=%v body-written=%v, want a failure before writing", name, err, body != "")
		}
	}
	if _, _, err := runFill(t, "", sumsFull); err == nil {
		t.Errorf("empty tag must fail")
	}
}
