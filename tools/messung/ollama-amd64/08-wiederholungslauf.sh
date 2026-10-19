#!/bin/bash
# A7: Wiederholungslauf von egress-init.sh: wiederholte Laeufe des erzeugten Startscripts in einem Container mit aktiver Restriktion,
# waehrenddessen Sonden in kurzem Takt (inner-08.sh). Ein Befund wird nicht behoben.
# Voraussetzung: 03-bau.sh. Aufruf: WORK=<messordner> tools/messung/ollama-amd64/08-wiederholungslauf.sh [laeufe, Vorgabe 30]
#   -> $OUT/rohdaten-08.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-08.txt

runs=${1:-30}
C=$PFX-c8
docker rm -f -v "$C" > /dev/null 2>&1
echo "== A7 Wiederholungslauf von egress-init.sh, $runs Laeufe, $(now)"
docker run --rm --name "$C" --cap-add=NET_ADMIN -v "$here":/m:ro "$(vimg oll)" bash /m/inner-08.sh "$runs"
echo "container-exit=$?"
