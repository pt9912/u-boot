#!/bin/bash
# D2 (ohne nestedRuntime): Wirkung der Messvariante "Upstream-DNS nur fuer den dnsmasq (nft meta skuid)" auf den direkten Pfad
# und auf erlaubte Namen, gegen den Ist-Zustand (erzeugtes Startscript unveraendert).
# Aufruf: WORK=<messordner> OUT=<ausgabeverzeichnis> 04-d2.sh -> $OUT/rohdaten-04-d2.txt
# Voraussetzung: 02-projekte.sh (Image u-boot-messung-dns-img-tool).
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-04-d2.txt
c=$PFX-d2
for mode in ist skuid; do
  echo; echo "######## D2, Modus $mode, $(now)"
  start $c $IMG_TOOL
  egress $c $mode zaehler
  docker exec -u vscode $c bash /m/inner-d2.sh $mode
  docker rm -f -v $c > /dev/null
done
