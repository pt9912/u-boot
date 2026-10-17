#!/usr/bin/env bash
# tap-formula-check.sh — prueft, dass die Formel zum Tag gehoert (tap-nachzug.sh, ADR-0016).
# Die Formel traegt keine `version`-Zeile (Homebrew liest die Version aus den URLs); die
# Pruefung liest daher die URLs: genau vier `url`-Zeilen (vier Plattformen), jede mit
# `/releases/download/<tag>/`. Ein falscher Tag in einer URL, eine fehlende oder eine
# zusaetzliche URL ist ein Fehler (Exit 1).
#
# AUFRUF: tap-formula-check.sh <tag> <formel-datei>
set -euo pipefail

tag="${1:-}"
file="${2:-}"
[ -n "$tag" ] && [ -f "$file" ] || { echo "tap-formula-check: usage: tap-formula-check.sh <tag> <formula-file>" >&2; exit 2; }

urls="$(grep -E '^[[:space:]]*url[[:space:]]' "$file" || true)"
total="$(printf '%s' "$urls" | grep -c . || true)"
ok="$(printf '%s' "$urls" | grep -cF "/releases/download/${tag}/" || true)"
if [ "$total" -ne 4 ] || [ "$ok" -ne 4 ]; then
	echo "tap-nachzug: the u-boot.rb asset of ${tag} does not carry four download URLs of ${tag} (found ${total} url lines, ${ok} of ${tag})." >&2
	exit 1
fi
