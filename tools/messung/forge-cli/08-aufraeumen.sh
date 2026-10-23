#!/bin/bash
# Aufraeumen. Entfernt nur Objekte, deren NAME das Praefix u-boot-messung-forge oder vsc-u-boot-messung-forge traegt (Container, Images,
# Volumes, Netze; praefix_entfernen in lib.sh). Nicht entfernt werden Tags ohne Praefix (auch Images, die der Bau der Hilfs-Images oder die
# Dev-Container-CLI gezogen hat), unbenannte Images, der Build-Cache und alle fremden Objekte. In die Rohdaten gelangen nur Anzahlen und Namen
# mit dem Praefix. Voraussetzung: rohdaten-08-vorher.txt aus 01-messstand.sh.
# Aufruf: WORK=<messordner> ./08-aufraeumen.sh      -> rohdaten-08.txt, rohdaten-08-nachher.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-08.txt

vorher="$OUT/rohdaten-08-vorher.txt"
[ -s "$vorher" ] || { echo "Abbruch: die Liste vorher fehlt oder ist leer (zuerst 01-messstand.sh)"; exit 2; }
echo "== Aufraeumen, $(now)"
echo "-- vorher (rohdaten-08-vorher.txt)"
cat "$vorher"
echo "-- vor dem Aufraeumen"
bestand
echo "-- Entfernen"
praefix_entfernen
echo "-- nachher (rohdaten-08-nachher.txt)"
bestand | tee "$OUT/rohdaten-08-nachher.txt"
aussage "objekte-mit-praefix-nachher=$(bestand | grep -c -E '^(container|image|volume|netz) ')"
echo
echo "-- Anzahl je Art vorher gegen nachher (Aenderungen fremder Arbeiten auf dem Host waehrend der Messung stehen hier ebenfalls)"
echo "vorher:  $(grep '^container=' "$vorher")"
echo "nachher: $(grep '^container=' "$OUT/rohdaten-08-nachher.txt")"
echo
echo "-- Build-Cache und Platz (der Build-Cache bleibt; er wurde nicht geleert)"
docker system df
