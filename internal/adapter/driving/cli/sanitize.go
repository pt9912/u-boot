package cli

import (
	"path/filepath"
	"strings"
)

// baseDirSanitizedError wraps an error such that Error() returns the
// inner-error's message with occurrences of `baseDir` rewritten as
// project-relative paths (slice-v1-cli-json-dry-run-remove T7 R14-
// MED-1 fix für Pfad-Leak in `diagnostic.message`; slice-v1-cli-json-
// dry-run-up-down T5 R2-MED-5 Helper-Extraktion für cluster-weite
// Wiederverwendung).
//
// Use-Case-Wraps der Form
// `fmt.Errorf("<cmd>: <action> %s: %w: %w", absPath, sentinel, raw)`
// tunneln den absoluten Pfad in `err.Error()`. Ohne Sanitisierung
// würde die User-facing Diagnostic den Filesystem-Layout des Users
// preisgeben (im JSON-Mode auch maschinen-lesbar abgreifbar).
//
// Sanitisierung-Regeln:
//   - `<baseDir>/foo` → `foo` (path-Separator project-relative)
//   - bare `<baseDir>` → `.` (project-root reference)
//   - leerer baseDir → unverändert (defensive identity)
//
// errors.Is/As bleiben intakt via Unwrap — der Wrapper ersetzt nur
// die Error()-Text-Form, nicht die Identity der Sentinels in der
// Chain.
type baseDirSanitizedError struct {
	inner   error
	baseDir string
}

func (e *baseDirSanitizedError) Error() string {
	msg := e.inner.Error()
	if e.baseDir == "" {
		return msg
	}
	sep := string(filepath.Separator)
	msg = strings.ReplaceAll(msg, e.baseDir+sep, "")
	return replaceBareBaseDir(msg, e.baseDir)
}

func (e *baseDirSanitizedError) Unwrap() error { return e.inner }

// replaceBareBaseDir replaces standalone occurrences of baseDir with
// `.`. A standalone occurrence is one followed by end-of-string or a
// byte that cannot continue a path-component name (i.e. not a letter/
// digit/`-`/`_`/`.`). Defense against Substring-Kollisionen wie
// `<baseDir>-cache/...` (R15-LOW-1): naive `strings.ReplaceAll(msg,
// baseDir, ".")` würde `proj-cache` zu `.-cache` mangeln.
//
// Use-Case-Layer löst dasselbe Pfad-Boundary-Problem für einzelne Pfade
// via `filepath.Rel` ([application.relativizePath]); für Error-Messages
// mit eingebettetem Pfad innerhalb von Prosa ist ein Boundary-Check
// die direkte Form (Rel braucht einen Pfad-String, hier haben wir eine
// Mischung aus Prosa und Pfad).
func replaceBareBaseDir(msg, baseDir string) string {
	if baseDir == "" || !strings.Contains(msg, baseDir) {
		return msg
	}
	var b strings.Builder
	b.Grow(len(msg))
	for i := 0; i < len(msg); {
		if strings.HasPrefix(msg[i:], baseDir) {
			end := i + len(baseDir)
			if end == len(msg) || !isPathComponentByte(msg[end]) {
				b.WriteByte('.')
				i = end
				continue
			}
		}
		b.WriteByte(msg[i])
		i++
	}
	return b.String()
}

// isPathComponentByte reports whether c can appear inside a single
// path-component name (no separator). Used by [replaceBareBaseDir] to
// distinguish a bare-baseDir match (`/foo/proj: error`) from a
// Substring-Kollision (`/foo/proj-cache/lock`). Non-ASCII bytes
// (UTF-8-Multibyte ≥ 0x80) gelten als Nicht-Continuation — ein
// non-ASCII-Pfad-Suffix direkt nach baseDir ist heute kein realer
// Case in unseren FS-Wraps; falsche Positives (Sanitisierung zu
// aggressiv) sind besser als falsche Negatives (Path-Leak).
func isPathComponentByte(c byte) bool {
	return c == '-' || c == '_' || c == '.' ||
		(c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9')
}

// sanitizeBaseDir wraps err with a baseDirSanitizedError, or returns
// err unchanged when err is nil. Convenience-Wrapper damit der
// Call-Site keinen Nil-Check braucht.
func sanitizeBaseDir(err error, baseDir string) error {
	if err == nil {
		return nil
	}
	return &baseDirSanitizedError{inner: err, baseDir: baseDir}
}
