#!/bin/bash
# Aufraeumen. Entfernt nur Objekte, deren NAME das Praefix u-boot-messung-ocloud oder vsc-u-boot-messung-ocloud traegt (Container, Images,
# Volumes, Netze; praefix_entfernen in lib.sh), und die Container der Dev-Container-CLI ueber das Label devcontainer.local_folder=/w/<Praefix>-<variante>.
# Nicht entfernt werden Tags ohne Praefix (die Basis-Images der Bauten), unbenannte Images, der Build-Cache und alle fremden Objekte.
# Voraussetzung: rohdaten-01-vorher.txt aus 01-messstand.sh.
# Aufruf: WORK=<messordner> tools/messung/ollama-cloud/05-aufraeumen.sh      -> $OUT/rohdaten-05.txt, rohdaten-05-nachher.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-05.txt

vorher="$OUT/rohdaten-01-vorher.txt"
[ -s "$vorher" ] || { echo "Abbruch: die Liste vorher fehlt oder ist leer (zuerst 01-messstand.sh)"; exit 2; }
echo "== Aufraeumen, $(now)"
echo "-- vorher (rohdaten-01-vorher.txt)"; cat "$vorher"
echo "-- vor dem Aufraeumen"; bestand
echo "-- Entfernen"
for v in std res alw; do
  docker ps -a -q --filter "label=devcontainer.local_folder=/w/$(vname "$v")" | while read -r c; do echo "entferne container der Dev-Container-CLI $c ($v)"; docker rm -f -v "$c" > /dev/null; done
done
praefix_entfernen
echo "-- nachher (rohdaten-05-nachher.txt)"
bestand | tee "$OUT/rohdaten-05-nachher.txt"
aussage "objekte-mit-praefix-nachher=$(bestand | grep -c -E '^(container|image|volume|netz) ')"
echo "vorher:  $(grep '^container=' "$vorher")"
echo "nachher: $(grep '^container=' "$OUT/rohdaten-05-nachher.txt")"
echo "(Aenderungen fremder Arbeiten auf dem Host waehrend der Messung stehen in diesem Vergleich ebenfalls)"
