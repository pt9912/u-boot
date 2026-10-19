#!/bin/bash
# Fuehrt die Schritte 01 bis 10 (mit 09b) der Reihe nach aus (Reihenfolge und Voraussetzungen: README.md des Protokolls).
# Aufruf: WORK=<messordner> OUT=<ausgabeverzeichnis> tools/messung/egress-dns-upstream/00-alles.sh
# Ein zweiter Lauf mit anderem OUT erzeugt neue Rohdaten; die AUSSAGE-Zeilen beider Laeufe lassen sich vergleichen (vergleich.sh).
set -u
: "${OUT:?OUT (Ausgabeverzeichnis der Rohdaten) fehlt}"
d=$(cd "$(dirname "$0")" && pwd)
for s in 01-messstand 02-projekte 03-d1 04-d2 05-d2-nested 06-d3 07-d3-netz 08-d5 09-d4 09b-gegenprobe 10-aufraeumen; do
  echo "######## $s $(date -u +%H:%M:%S)"
  "$d/$s.sh" > /dev/null 2>&1; echo "$s exit=$?"
done
