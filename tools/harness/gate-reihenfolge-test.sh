#!/usr/bin/env bash
# gate-reihenfolge-test — Waechter-Test fuer gate-reihenfolge.sh.
#
# Der Grundzustand (Kopie der Make-Dateien dieses Repos) muss gruen sein; die
# Mutation (tool-pruefung per += ans Ende der Liste gehaengt, die Form, die die Zusage
# in tools/mk/tool-pruefung.mk ausschliesst) muss das Skript rot faerben und tool-pruefung
# nennen. Ein Fall, der nicht greift, bricht den Test.
#
# Kein Docker, kein Netz: bash + make.
set -euo pipefail

root="$(cd "$(dirname "$0")/../.." && pwd)"
skript="$root/tools/harness/gate-reihenfolge.sh"
arbeit="$(mktemp -d)"
trap 'rm -rf "$arbeit"' EXIT

versagt=0
faelle=0

kopie() {
  rm -rf "${arbeit:?}/r"
  mkdir -p "$arbeit/r/tools"
  (cd "$root" && cp -r Makefile repo.mk d-check.mk harness "$arbeit/r/")
  cp -r "$root/tools/mk" "$arbeit/r/tools/mk"
}

erwarte() {
  local name="$1" soll="$2" muster="$3" ist=0 aus
  faelle=$((faelle + 1))
  aus="$(GATE_REIHENFOLGE_WURZEL="$arbeit/r" bash "$skript" 2>&1)" || ist=$?
  if [ "$ist" -ne "$soll" ] || ! grep -qF -- "$muster" <<<"$aus"; then
    echo "FAIL  $name: Exit $ist (erwartet $soll), Muster '$muster'" >&2
    echo "$aus" | sed 's/^/      /' >&2
    versagt=$((versagt + 1))
    return
  fi
  echo "ok    $name"
}

kopie
erwarte "tool-pruefung steht vorn" 0 "is the first check"

kopie
sed -i 's/^GATE_CHECKS := tool-pruefung \$(GATE_CHECKS)/GATE_CHECKS += tool-pruefung/' "$arbeit/r/tools/mk/tool-pruefung.mk"
erwarte "tool-pruefung per += am Ende" 1 "tool-pruefung is not the first check"

echo "gate-reihenfolge-test: $faelle Faelle, $versagt Fehlschlaege"
[ "$versagt" -eq 0 ]
