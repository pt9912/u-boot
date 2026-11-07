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
# golangci-lint: "golangci-lint has version 2.15.0 ..."; der Pin heisst "v2.15.0"
# (repo.mk) — verglichen wird der Versionsteil.
tool golangci-lint "${LINT_PIN#v}" "golangci-lint"

# make: die Mindest-Form der Inventur (der Parse-Guard des Makefiles traegt
# denselben Floor; hier die Laufzeit-Lesart, weil ein make eines anderen
# Kontexts den Floor nicht trägt).
MAKE_FLOOR="$(sed -n 's/^MAKE_MIN_MAJOR := *//p' "$root/Makefile" | head -1).$(sed -n 's/^MAKE_MIN_MINOR := *//p' "$root/Makefile" | head -1)"
[ -n "$MAKE_FLOOR" ] || fail "the MAKE_MIN_VERSION floor of the Makefile is unparseable"
have="$(make --version 2>/dev/null | sed -n '1s/^GNU Make \([0-9.]*\).*$/\1/p')"
[ -n "$have" ] || fail "make reports no GNU make version line"
awk -v have="$have" -v floor="$MAKE_FLOOR" 'BEGIN{ m=(have<floor) ; if (m) exit 1 }' \
  || fail "make reports $have, the floor of the inventory is $MAKE_FLOOR"

# Die Anwesenheits-Form der Inventur (Vollstaendigkeit der Entwicklungsumgebung;
# LX-Feld: je Werkzeug je Pruefform, LH-FA-DEV-011): go/golangci-lint exakt oben,
# make ueber den Floor, und diese ueber die Anwesenheit. Ein fehlendes Werkzeug
# ist ein Befund mit Name; sed/awk/coreutils tragen keine Einzel-Form je
# Basis-Image-Grenze des Inventur-Absatzes.
for t in bash curl gh git docker; do
  command -v "$t" >/dev/null 2>&1 || fail "the tool $t of the development inventory is not installed in the running context"
done
echo "tool-pruefung: the running context carries the pinned toolchain (go $GO_PIN, golangci-lint ${LINT_PIN#v}, make >= $MAKE_FLOOR) and the complete development inventory (bash, curl, gh, git, docker client) — OK"