#!/bin/bash
# Aufraeumen: entfernt nur Objekte, deren NAME das Praefix u-boot-messung-dns traegt (Container, Images, Volumes, Netze; praefix_entfernen in
# lib.sh). Nicht entfernt werden Tags ohne Praefix (auch Images, die der Bau gezogen hat), unbenannte Images, der Build-Cache und alle fremden
# Objekte. In die Rohdaten gelangen nur Anzahlen und Namen mit dem Praefix.
# Voraussetzung: rohdaten-11-vorher.txt aus 01-messstand.sh.
# Aufruf: WORK=<messordner> OUT=<ausgabeverzeichnis> 10-aufraeumen.sh -> $OUT/rohdaten-11.txt, rohdaten-11-nachher.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-11.txt
vorher="$OUT/rohdaten-11-vorher.txt"
[ -s "$vorher" ] || { echo "Abbruch: die Liste vorher fehlt oder ist leer (zuerst 01-messstand.sh)"; exit 2; }
echo "== Aufraeumen, $(now)"
echo "-- vorher (rohdaten-11-vorher.txt)"
cat "$vorher"
echo "-- vor dem Aufraeumen"
bestand
echo "-- Entfernen"
praefix_entfernen
echo "-- nachher (rohdaten-11-nachher.txt)"
bestand | tee "$OUT/rohdaten-11-nachher.txt"
aussage "objekte-mit-praefix-nachher=$(bestand | grep -c -E '^(container|image|volume|netz) ')"
echo
echo "vorher:  $(grep '^container=' "$vorher")"
echo "nachher: $(grep '^container=' "$OUT/rohdaten-11-nachher.txt")"
