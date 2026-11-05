#!/bin/bash
# Fuehrt die Schritte 01 bis 04, 06 und 07 der Reihe nach aus (Reihenfolge und Voraussetzungen: README.md des Protokolls); 05-aufraeumen.sh laeuft getrennt.
# Aufruf: UBOOT=<binary> WORK=<messordner> [OUT=<ausgabeverzeichnis>] tools/messung/devcontainer-homebrew-pakete/00-alles.sh
set -u
d=$(cd "$(dirname "$0")" && pwd)
for s in 01-bau 02-pfad 03-aktualisierung 04-egress 06-arm64 07-fremd-tap; do
  echo "######## $s $(date -u +%H:%M:%S)"
  "$d/$s.sh" > /dev/null 2>&1; echo "$s exit=$?"
done
