#!/bin/bash
# M5 Plattform: laeuft ein Container der Architektur arm64 auf diesem Host? Gemessen wird nur, ob die Engine einen arm64-Container des
# Basis-Images startet (Emulation); der Bau des Blocks auf arm64 ist Gegenstand der Zusage erst, wenn dieser Schritt gelingt.
# Aufruf: UBOOT=<binary> WORK=<messordner> ./06-arm64.sh      -> rohdaten-m5.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-m5.txt

echo "== M5 Plattform, $(now)"
echo "host: $(uname -m) docker=$(docker version --format '{{.Server.Version}} {{.Server.Os}}/{{.Server.Arch}}')"
echo "-- arm64-Container des Basis-Images"
out=$(docker run --rm --platform linux/arm64 "$BASE" uname -m 2>&1); rc=$?
echo "$out" | tail -3
echo "exit=$rc"
if [ $rc -eq 0 ]; then aussage "m5-arm64=lauffaehig"; else aussage "m5-arm64=nicht-lauffaehig"; fi
echo "== ende $(now)"
