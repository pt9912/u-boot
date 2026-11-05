#!/bin/bash
# Fuehrt die Messungen M1 bis M4 der Reihe nach aus (01 bis 04); 05-aufraeumen.sh laeuft getrennt.
# Aufruf: UBOOT=<binary> WORK=<messordner> [OUT=<ausgabeverzeichnis>] tools/messung/slice-v1-devcontainer-remove-messung/00-alles.sh
set -u
d=$(cd "$(dirname "$0")" && pwd)
for s in 01-m1-volumes 02-m2-entfernen 03-m3-instanz 04-m4-ueberleben; do
  echo "######## $s $(date -u +%H:%M:%S)"
  "$d/$s.sh" > /dev/null 2>&1; echo "$s exit=$?"
done
