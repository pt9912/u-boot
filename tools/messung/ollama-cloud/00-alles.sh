#!/bin/bash
# Fuehrt die Schritte der Reihe nach aus (Reihenfolge und Voraussetzungen: README.md des Protokolls); 04-lauf.sh zweimal (LAUF=1, LAUF=2).
# Aufruf: WORK=<messordner> [OUT=<ausgabeverzeichnis>] tools/messung/ollama-cloud/00-alles.sh [cloud-modell]
# redact-probe.sh laeuft getrennt VOR dem ersten Lauf; 06-suche.sh getrennt NACH dem Lauf, vor dem Commit.
set -u
d=$(cd "$(dirname "$0")" && pwd)
for s in 01-messstand 02-projekte 03-lauf-vorbereitung; do
  echo "######## $s $(date -u +%H:%M:%S)"
  "$d/$s.sh" > /dev/null 2>&1; echo "$s exit=$?"
done
for n in 1 2; do
  echo "######## 04-lauf LAUF=$n $(date -u +%H:%M:%S)"
  LAUF=$n "$d/04-lauf.sh" "$@" > /dev/null 2>&1; echo "04-lauf LAUF=$n exit=$?"
done
echo "######## 05-aufraeumen $(date -u +%H:%M:%S)"
"$d/05-aufraeumen.sh" > /dev/null 2>&1; echo "05-aufraeumen exit=$?"
