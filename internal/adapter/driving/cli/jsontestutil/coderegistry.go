// Package jsontestutil bietet schema-konforme Assertion-Helper für
// `u-boot --json`-Ausgaben (slice-v1-cli-json-dry-run-doctor T2).
// Spec-Anker: docs/user/cli-json-output.md zitiert das Lastenheft
// §1823-1842 (Minimalkontrakt) und §322-417 (Voll-Schema) verbatim;
// dieses Package prüft Schema-Konformität als Go-Code, kein embedded
// JSON-Schema, kein zusätzlicher Dep.
//
// Public-API: [AssertMinimalEnvelope] und [AssertFullEnvelope].
// Code-Registry: [DefaultAllowedCodes].
package jsontestutil

// DefaultAllowedCodes liefert die Code-Registry für
// `diagnostics[].code`. Spec §1835 / §445 erlaubt zwei Quellen:
// LH-Kennungen (`LH-FA-DEV-003`, …) oder tool-interne Codes, falls
// ihre Bedeutung dokumentiert ist. u-boot verwendet tool-interne
// Codes mit Dotted-Notation; diese Funktion ist die Source-of-Truth
// für den Helper.
//
// Source-of-Truth-Disziplin: die zurückgegebene Map ist der
// kanonische Code-Satz. Markdown-Doku-Form lebt in
// docs/user/cli-json-output.md §5; der drift_test.go im
// jsontestutil-Package erzwingt symmetrische Synchronisation
// (Gate 2 aus dem T0-(h)-Outcome).
//
// Folge-Slice-Pflicht: jeder Slice, der neue Diagnostic-Codes
// einführt, ergänzt sie in dieser Map UND in der Markdown-Doku
// im selben Slice-Closure-Commit — sonst bricht Gate 2.
//
// gochecknoglobals-konform via Funktions-Wrapper (Repo-Konvention).
func DefaultAllowedCodes() map[string]string {
	return map[string]string{
		// Doctor-Checks (LH-FA-DIAG-002, application/doctor.go:74-114)
		"fs.write-permissions":                   "doctor: Schreib-Permission im Working Directory",
		"git.installed":                          "doctor: Git-Binary verfügbar",
		"docker.installed":                       "doctor: Docker-Binary verfügbar",
		"docker.reachable":                       "doctor: Docker-Daemon erreichbar",
		"docker.compose.installed":               "doctor: Compose-Plugin verfügbar",
		"uboot.yaml.valid":                       "doctor: u-boot.yaml syntaktisch valide",
		"compose.yaml.valid":                     "doctor: compose.yaml syntaktisch valide",
		"devcontainer.json.valid":                "doctor: .devcontainer/devcontainer.json syntaktisch valide",
		"devcontainer.dockerfile.valid":          "doctor: .devcontainer/Dockerfile parsebar",
		"services.enabled-key":                   "doctor: u-boot.yaml services-Block konsistent",
		"devcontainer.forwardPorts.consistency":  "doctor: devcontainer.json forwardPorts konsistent",
		"devcontainer.features.allowlist":        "doctor: devcontainer features auf Allowlist",
		"devcontainer.features.drift":            "doctor: devcontainer features ohne Drift",
	}
}
