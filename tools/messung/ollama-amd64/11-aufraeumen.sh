#!/bin/bash
# A10: Aufraeumen. Entfernt nur Objekte, deren NAME das Praefix u-boot-messung-amd64 oder vsc-u-boot-messung-amd64 traegt (Container,
# Images, Volumes, Netze; praefix_entfernen in lib.sh), und den Container der Dev-Container-CLI ueber das Label
# devcontainer.local_folder=/w/<Praefix>-oll, das 10-devcontainer-cli.sh setzt. Nicht entfernt werden Tags ohne Praefix (auch Images, die
# der Bau der Hilfs-Images gezogen hat), unbenannte Images, Zwischencontainer fehlgeschlagener Bauten (09-tests.sh weist sie aus), der
# Build-Cache und alle fremden Objekte. In die Rohdaten gelangen nur Anzahlen und Namen mit dem Praefix.
# Voraussetzung: rohdaten-11-vorher.txt aus 01-messstand.sh.
# Aufruf: WORK=<messordner> tools/messung/ollama-amd64/11-aufraeumen.sh      -> $OUT/rohdaten-11.txt, rohdaten-11-nachher.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-11.txt

vorher="$OUT/rohdaten-11-vorher.txt"
[ -s "$vorher" ] || { echo "Abbruch: die Liste vorher fehlt oder ist leer (zuerst 01-messstand.sh)"; exit 2; }
echo "== A10 Aufraeumen, $(now)"
echo "-- vorher (rohdaten-11-vorher.txt)"
cat "$vorher"
echo "-- vor dem Aufraeumen"
bestand

echo "-- Entfernen"
docker ps -a -q --filter "label=devcontainer.local_folder=/w/$PFX-oll" | while read -r c; do echo "entferne container der Dev-Container-CLI $c"; docker rm -f -v "$c" > /dev/null; done
praefix_entfernen

echo "-- nachher (rohdaten-11-nachher.txt)"
bestand | tee "$OUT/rohdaten-11-nachher.txt"
aussage "objekte-mit-praefix-nachher=$(bestand | grep -c -E '^(container|image|volume|netz) ')"
echo
echo "-- Anzahl je Art vorher gegen nachher (Aenderungen fremder Arbeiten auf dem Host waehrend der Messung stehen hier ebenfalls)"
echo "vorher:  $(grep '^container=' "$vorher")"
echo "nachher: $(grep '^container=' "$OUT/rohdaten-11-nachher.txt")"
echo
echo "-- Build-Cache und Platz (der Build-Cache bleibt; er wurde nicht geleert)"
docker system df
docker run --rm -v /var/lib/docker:/d:ro "$BASE" df -h /d | tail -1
