#!/usr/bin/env bash
# gate-reihenfolge — pinnt die Zusage aus tools/mk/tool-pruefung.mk: in der
# Prerequisite-Liste von record-gates steht tool-pruefung an erster Stelle, der
# Nachweis record-gates selbst haengt an allen Checks. Gelesen wird die Liste aus
# der Regel-Datenbank von make (make -pn), nicht der Text einer Rezeptzeile.
#
# Aufruf: tools/harness/gate-reihenfolge.sh   (Exit 0 gruen, 1 bei Abweichung)
# GATE_REIHENFOLGE_WURZEL liest eine andere Make-Wurzel als dieses Repo (der
# Waechter-Test gate-reihenfolge-test.sh nutzt das fuer die rote Probe).
set -euo pipefail
root="${GATE_REIHENFOLGE_WURZEL:-$(cd "$(dirname "$0")/../.." && pwd)}"
cd "$root"

zeile="$(make -pn gates 2>/dev/null | grep -m1 '^record-gates:' || true)"
[ -n "$zeile" ] || { echo "gate-reihenfolge: no record-gates rule found in the make database" >&2; exit 1; }

# Die Prerequisites stehen hinter dem Doppelpunkt; das erste Wort ist der erste Check.
voraus="${zeile#record-gates:}"
set -- $voraus
if [ "${1:-}" != "tool-pruefung" ]; then
  echo "gate-reihenfolge: tool-pruefung is not the first check of record-gates (first: '${1:-<none>}')" >&2
  exit 1
fi
echo "gate-reihenfolge: tool-pruefung is the first check of record-gates — OK"
