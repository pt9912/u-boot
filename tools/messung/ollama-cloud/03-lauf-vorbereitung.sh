#!/bin/bash
# Hilfs-Image der Dev-Container-CLI bauen (docker build -f Dockerfile.cli). Kein Messschritt; Voraussetzung von 04-lauf.sh.
# Aufruf: WORK=<messordner> tools/messung/ollama-cloud/03-lauf-vorbereitung.sh    -> $OUT/rohdaten-03.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-03.txt
echo "== Hilfs-Image $CLI, $(now)"
docker build -q -f "$here/Dockerfile.cli" -t "$CLI" "$here" > "$WORK/cli-build.log" 2>&1; rc=$?
echo "build exit=$rc"; [ $rc = 0 ] || { tail -5 "$WORK/cli-build.log"; exit 1; }
docker run --rm "$CLI" devcontainer --version
echo "-- Namen der Umgebungsvariablen im CLI-Container (Werte nicht ausgegeben)"
docker run --rm "$CLI" env | cut -d= -f1 | sort | tr '\n' ' '; echo
aussage "cli-version=$(docker run --rm "$CLI" devcontainer --version)"
