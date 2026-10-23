#!/bin/bash
# Fuehrt die Schritte 01 bis 07 der Reihe nach aus (Reihenfolge und Voraussetzungen: README.md des Protokolls); 08-aufraeumen.sh laeuft getrennt.
# Aufruf: WORK=<messordner> [OUT=<ausgabeverzeichnis>] tools/messung/forge-cli/00-alles.sh
set -u
d=$(cd "$(dirname "$0")" && pwd)
for s in 01-messstand 02-oci-quellen 03-lieferkette 04-projekte 05-bau 06-dockerfile-bloecke 07-laufzeit; do
  echo "######## $s $(date -u +%H:%M:%S)"
  if [ "$s" = 07-laufzeit ]; then
    "$d/$s.sh" ghk:gh ghxk:gh glabsk:glab glabpk:glab glabzk:glab ghdeb:gh glabdeb:glab ghtgz:gh glabtgz:glab > /dev/null 2>&1; echo "$s exit=$?"
  else
    "$d/$s.sh" > /dev/null 2>&1; echo "$s exit=$?"
  fi
done
