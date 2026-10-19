#!/bin/bash
# O3 bis O5: Laufzeit im Probe-Container (Image aus 04-build-probe.sh).
# Aufruf: 05-laufzeit.sh discover "" <modell...>        > rohdaten-05-discover.txt 2>&1
#         05-laufzeit.sh restricted "<hosts>" <modell...> > rohdaten-05-restricted-<n>.txt 2>&1
set -u
here=$(cd "$(dirname "$0")" && pwd)
mode=$1; hosts=$2; shift 2
echo "== Modus $mode, Hosts: [$hosts], Modelle: $*"
time docker run --rm --cap-add=NET_ADMIN -v "$here":/m:ro u-boot-messung-ollama-probe bash /m/inner-05.sh "$mode" "$hosts" "$@"
