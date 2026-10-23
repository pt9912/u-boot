#!/bin/bash
# Lieferkette der Release-Artefakte von gh und glab (Signatur, Attestierung, Pruefsumme, apt-Quellen, Debian-Pakete); alles anonym.
# Aufruf: WORK=<messordner> ./03-lieferkette.sh      -> rohdaten-03.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-03.txt

echo "== Lieferkette, $(now)"
tool bash /m/inner-03.sh
echo
echo "################ Debian-Pakete im Basis-Image (Debian 13)"
docker run --rm -u root -v "$here":/m:ro "$BASE" bash /m/inner-03-debian.sh
