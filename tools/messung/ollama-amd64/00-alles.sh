#!/bin/bash
# Fuehrt die Schritte 01 bis 11 der Reihe nach aus (Reihenfolge und Voraussetzungen: README.md des Protokolls).
# Aufruf: WORK=<messordner> [OUT=<ausgabeverzeichnis>] tools/messung/ollama-amd64/00-alles.sh
# Ein zweiter Lauf mit anderem OUT erzeugt neue Rohdaten; `AUSSAGE`-Zeilen beider Laeufe lassen sich vergleichen (vergleich.sh).
set -u
d=$(cd "$(dirname "$0")" && pwd)
for s in 01-messstand 02-projekte 03-bau 04-archiv 05-laufzeit 06-volume 07-kombination 08-wiederholungslauf 09-tests 10-devcontainer-cli 11-aufraeumen; do
  echo "######## $s $(date -u +%H:%M:%S)"
  "$d/$s.sh" > /dev/null 2>&1; echo "$s exit=$?"
done
