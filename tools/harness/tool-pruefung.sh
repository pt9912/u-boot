#!/usr/bin/env bash
# tool-pruefung.sh — das Laufzeit-Gate der Pin-Inventur (LH-FA-DEV-011).
# pinsync prueft die Werte an den Stellen; dieses Gate prueft die INSTALLIERTEN
# Werkzeuge des laufenden Kontexts gegen dieselben Pin-Werte — ein Lauf, dessen
# Werkzeugkette nicht die Pins traegt, endet rot (der Lauf-Ort ist egal, die
# Werkzeuge sind die Trageschaft: ADR-0043, LH-FA-DEV-011 Pin-Inventur).
#
# Gelesen wird der Soll-Wert aus der Pin-Stelle selbst (repo.mk), nicht aus
# einer Kopie — aendert jemand die Stelle, zwingt pinsync zur Konsistenz und
# dieses Gate zur Ausfuehrung.
#
# Aufruf: tools/harness/tool-pruefung.sh   (Exit 0 gruen, 1 bei Abweichung)
set -euo pipefail
root="$(cd "$(dirname "$0")/../.." && pwd)"

fail() { echo "tool-pruefung: $1" >&2; exit 1; }

tool() {
  local cmd="$1" want="$2" name="$3"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    fail "the tool $name of the pin $(echo "$want" | head -1) is not installed in the running context"
  fi
  local have
  have="$("$cmd" --version 2>/dev/null | head -1 || "$cmd" version 2>/dev/null | head -1 || true)"
  case "$have" in
    *"$want"*) : ;;
    *) fail "$name reports \"$have\" (first line), the pin is \"$want\" — the running toolchain does not carry the pinned version" ;;
  esac
}

GO_PIN=$(sed -n 's/^GO_VERSION[[:space:]]*?=[[:space:]]*//p' "$root/repo.mk" | head -1)
LINT_PIN=$(sed -n 's/^GOLANGCI_LINT_VERSION[[:space:]]*?=[[:space:]]*//p' "$root/repo.mk" | head -1)
[ -n "$GO_PIN" ] || fail "the GO_VERSION pin of repo.mk is unparseable"
[ -n "$LINT_PIN" ] || fail "the GOLANGCI_LINT_VERSION pin of repo.mk is unparseable"

# go: die Ausgabe heisst "go version go1.27.3 linux/amd64"; der Pin heisst
# "1.27.3" (repo.mk) — der Versionsteil der Ausgabe trägt ihn.
tool go "$GO_PIN" "go"
# golangci-lint: "golangci-lint has version 2.14.0 ..."; der Pin heisst "2.14.0"
# (repo.mk trägt v2.14.0? Nein — die Pin-Stelle der lint-Version heisst "v2.14.0"
# am Template und "2.14.0" je repo.mk-Lesung; verglichen wird der Teil nach dem v).
tool golangci-lint "${LINT_PIN#v}" "golangci-lint"

echo "tool-pruefung: the running context carries the pinned toolchain (go $GO_PIN, golangci-lint ${LINT_PIN#v}) — OK"