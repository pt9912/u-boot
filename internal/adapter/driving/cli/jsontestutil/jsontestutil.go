package jsontestutil

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

// assertConfig sammelt Options für [AssertMinimalEnvelope] und
// [AssertFullEnvelope]. Per Functional-Options gesetzt — niemand
// muss alle Felder kennen, wenn er nur Subset-Pins braucht.
type assertConfig struct {
	expectedCommand    string
	expectedSubcommand string
	expectedCodes      []string
	expectedExitCode   *int
	allowedCodes       map[string]string
}

// AssertOption ist die Functional-Options-Form für die Helper.
// Die Helper-API ist absichtlich klein: vier Options decken alle
// Pin-Wünsche der Folge-Slices ab. Insbesondere fehlt eine
// `WithAllowedCodes`-Option BEWUSST — globale Code-Erweiterungen
// passieren nur in [DefaultAllowedCodes] plus Markdown-Doku
// (Single-Source-of-Truth-Disziplin, T0-(h)).
type AssertOption func(*assertConfig)

// WithCommand pinnt den erwarteten `command`-Wert. Bei Mismatch
// fail-grade.
func WithCommand(cmd string) AssertOption {
	return func(c *assertConfig) { c.expectedCommand = cmd }
}

// WithSubcommand pinnt den `subcommand`-Wert. Spec §1838 fordert
// `subcommand`-Pflicht bei `command ∈ {template, config}` —
// wird vom Helper bei den beiden Commands auch ohne Option geprüft.
func WithSubcommand(sub string) AssertOption {
	return func(c *assertConfig) { c.expectedSubcommand = sub }
}

// WithExpectedCodes pinnt einen Subset-Erwartung an
// `diagnostics[].code`: der konkrete Test-Output MUSS genau
// diese Codes enthalten (als Set, Reihenfolge egal). Nicht
// für globale Allowlist-Erweiterung — siehe AssertOption-Doc.
func WithExpectedCodes(codes ...string) AssertOption {
	return func(c *assertConfig) { c.expectedCodes = codes }
}

// WithExitCode pinnt den erwarteten `exitCode`-Wert. Häufiger
// Pin-Wunsch (`LH-FA-CLI-006`-Klasse).
func WithExitCode(code int) AssertOption {
	return func(c *assertConfig) { c.expectedExitCode = &code }
}

// AssertMinimalEnvelope prüft den Minimalkontrakt aus LH-NFA-USE-
// 004 §1841 gegen raw. Failures via t.Errorf — Helper failt nicht
// fatal, damit Tests mehrere Befunde gleichzeitig sehen.
//
// Geprüft wird:
//   - Pflicht-Set: status, command, diagnostics, exitCode
//     (subcommand Pflicht bei template/config, sonst optional)
//   - status ∈ {ok, warn, error}
//   - diagnostics[i].level ∈ {warn, error} (NICHT ok, NICHT info)
//   - diagnostics[i].code in DefaultAllowedCodes oder LH-konform
//   - status-Kopplung an höchsten level
//   - exitCode ≥ 0
//   - Voll-Schema-Felder dryRun/diff/plannedFiles/changes FEHLEN
//     (Minimalkontrakt rejected sie, Spec §1841)
func AssertMinimalEnvelope(t testing.TB, raw []byte, opts ...AssertOption) {
	t.Helper()
	cfg := buildConfig(opts...)

	env, ok := parseEnvelope(t, raw)
	if !ok {
		return
	}

	checkRequiredMinimal(t, env)
	checkStatus(t, env)
	checkCommand(t, env, cfg)
	checkSubcommand(t, env, cfg)
	checkExitCode(t, env, cfg)
	checkDiagnostics(t, env, cfg)
	checkStatusCoupling(t, env)
	checkNoFullFields(t, env)
}

// AssertFullEnvelope prüft das Voll-Schema aus LH-FA-CLI-007 §326
// gegen raw. Zusätzlich zu den Minimal-Checks werden alle vier
// Voll-Felder als Pflicht geprüft (`dryRun`/`diff`/`plannedFiles`/
// `changes` müssen ALLE im JSON erscheinen — auch wenn dryRun/diff
// `false` und plannedFiles/changes `[]` sind).
//
// In diesem Slice (Doctor) wird der Voll-Helper als Stub angelegt,
// aber NICHT verwendet — Erstnutzung im Folge-Slice
// slice-v1-cli-json-dry-run-add. Die Tests in
// jsontestutil_test.go decken trotzdem positive/negative Cases ab.
func AssertFullEnvelope(t testing.TB, raw []byte, opts ...AssertOption) {
	t.Helper()
	cfg := buildConfig(opts...)

	env, ok := parseEnvelope(t, raw)
	if !ok {
		return
	}

	checkRequiredFull(t, env)
	checkStatus(t, env)
	checkCommand(t, env, cfg)
	checkSubcommand(t, env, cfg)
	checkExitCode(t, env, cfg)
	checkDiagnostics(t, env, cfg)
	checkStatusCoupling(t, env)
	checkPlannedFiles(t, env)
	checkChanges(t, env)
}

func buildConfig(opts ...AssertOption) assertConfig {
	cfg := assertConfig{
		allowedCodes: DefaultAllowedCodes(),
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func parseEnvelope(t testing.TB, raw []byte) (map[string]any, bool) {
	t.Helper()
	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Errorf("envelope JSON unparsable: %v\nraw=%s", err, raw)
		return nil, false
	}
	return env, true
}

func checkRequiredMinimal(t testing.TB, env map[string]any) {
	t.Helper()
	for _, k := range []string{"status", "command", "diagnostics", "exitCode"} {
		if _, present := env[k]; !present {
			t.Errorf("missing required field %q (Spec §1823-1839)", k)
		}
	}
}

func checkRequiredFull(t testing.TB, env map[string]any) {
	t.Helper()
	required := []string{"status", "command", "dryRun", "diff", "plannedFiles", "changes", "diagnostics", "exitCode"}
	for _, k := range required {
		if _, present := env[k]; !present {
			t.Errorf("missing required field %q (Spec §326)", k)
		}
	}
}

func checkStatus(t testing.TB, env map[string]any) {
	t.Helper()
	s, ok := env["status"].(string)
	if !ok {
		return
	}
	if !slices.Contains([]string{"ok", "warn", "error"}, s) {
		t.Errorf("status %q not in {ok, warn, error}", s)
	}
}

func checkCommand(t testing.TB, env map[string]any, cfg assertConfig) {
	t.Helper()
	cmd, _ := env["command"].(string)
	if cfg.expectedCommand != "" && cmd != cfg.expectedCommand {
		t.Errorf("command: want %q, got %q", cfg.expectedCommand, cmd)
	}
}

func checkSubcommand(t testing.TB, env map[string]any, cfg assertConfig) {
	t.Helper()
	cmd, _ := env["command"].(string)
	sub, hasSub := env["subcommand"].(string)

	// Spec §1838: subcommand Pflicht bei template/config.
	if cmd == "template" || cmd == "config" {
		if !hasSub || sub == "" {
			t.Errorf("subcommand required when command == %q (Spec §1838)", cmd)
		}
	}

	if cfg.expectedSubcommand != "" && sub != cfg.expectedSubcommand {
		t.Errorf("subcommand: want %q, got %q", cfg.expectedSubcommand, sub)
	}
}

func checkExitCode(t testing.TB, env map[string]any, cfg assertConfig) {
	t.Helper()
	code, ok := env["exitCode"].(float64)
	if !ok {
		return
	}
	if code < 0 {
		t.Errorf("exitCode %v must be ≥ 0 (Spec §388)", code)
	}
	if cfg.expectedExitCode != nil && int(code) != *cfg.expectedExitCode {
		t.Errorf("exitCode: want %d, got %d", *cfg.expectedExitCode, int(code))
	}
}

func checkDiagnostics(t testing.TB, env map[string]any, cfg assertConfig) {
	t.Helper()
	diags, ok := env["diagnostics"].([]any)
	if !ok {
		t.Errorf("diagnostics must be an array (Spec §373)")
		return
	}

	foundCodes := make(map[string]bool)
	for i, raw := range diags {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Errorf("diagnostics[%d] must be an object", i)
			continue
		}
		checkDiagnosticItem(t, i, item, cfg)
		if code, _ := item["code"].(string); code != "" {
			foundCodes[code] = true
		}
	}

	if cfg.expectedCodes != nil {
		for _, want := range cfg.expectedCodes {
			if !foundCodes[want] {
				t.Errorf("expected code %q missing from diagnostics", want)
			}
		}
	}
}

func checkDiagnosticItem(t testing.TB, i int, item map[string]any, cfg assertConfig) {
	t.Helper()
	for _, k := range []string{"level", "code", "message"} {
		if _, present := item[k]; !present {
			t.Errorf("diagnostics[%d] missing required key %q (Spec §377)", i, k)
		}
	}

	level, _ := item["level"].(string)
	if level != "warn" && level != "error" {
		t.Errorf("diagnostics[%d].level %q must be warn or error (Spec §1834)", i, level)
	}

	code, _ := item["code"].(string)
	if !codeAllowed(code, cfg.allowedCodes) {
		t.Errorf("diagnostics[%d].code %q not in DefaultAllowedCodes and not LH-conform (Spec §1835 / §445)", i, code)
	}
}

// codeAllowed validates `code` against the two allowed forms of
// Spec §445: LH-IDs (prefix "LH-") OR documented tool-internal
// codes (must appear in the registry map).
func codeAllowed(code string, allowed map[string]string) bool {
	if code == "" {
		return false
	}
	if strings.HasPrefix(code, "LH-") {
		return true
	}
	_, ok := allowed[code]
	return ok
}

func checkStatusCoupling(t testing.TB, env map[string]any) {
	t.Helper()
	status, _ := env["status"].(string)
	diags, _ := env["diagnostics"].([]any)

	highestLevel := "ok"
	for _, raw := range diags {
		item, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		level, _ := item["level"].(string)
		if level == "error" {
			highestLevel = "error"
			break
		}
		if level == "warn" && highestLevel != "error" {
			highestLevel = "warn"
		}
	}

	if status != highestLevel {
		t.Errorf("status %q decoupled from highest diagnostics level %q (Spec §447 / §1837)", status, highestLevel)
	}
}

func checkNoFullFields(t testing.TB, env map[string]any) {
	t.Helper()
	for _, k := range []string{"dryRun", "diff", "plannedFiles", "changes"} {
		if _, present := env[k]; present {
			t.Errorf("minimal envelope must not contain %q (Spec §1841 / §1842)", k)
		}
	}
}

func checkPlannedFiles(t testing.TB, env map[string]any) {
	t.Helper()
	arr, ok := env["plannedFiles"].([]any)
	if !ok {
		t.Errorf("plannedFiles must be an array (Spec §346)")
		return
	}
	for i, raw := range arr {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Errorf("plannedFiles[%d] must be an object", i)
			continue
		}
		for _, k := range []string{"path", "action"} {
			if _, present := item[k]; !present {
				t.Errorf("plannedFiles[%d] missing required key %q (Spec §350)", i, k)
			}
		}
		action, _ := item["action"].(string)
		if !slices.Contains([]string{"create", "modify", "delete"}, action) {
			t.Errorf("plannedFiles[%d].action %q not in {create, modify, delete} (Spec §354)", i, action)
		}
	}
}

func checkChanges(t testing.TB, env map[string]any) {
	t.Helper()
	arr, ok := env["changes"].([]any)
	if !ok {
		t.Errorf("changes must be an array (Spec §361)")
		return
	}
	for i, raw := range arr {
		item, ok := raw.(map[string]any)
		if !ok {
			t.Errorf("changes[%d] must be an object", i)
			continue
		}
		for _, k := range []string{"path", "count"} {
			if _, present := item[k]; !present {
				t.Errorf("changes[%d] missing required key %q (Spec §365)", i, k)
			}
		}
		count, ok := item["count"].(float64)
		if ok && count < 0 {
			t.Errorf("changes[%d].count %v must be ≥ 0 (Spec §368)", i, count)
		}
	}
}

