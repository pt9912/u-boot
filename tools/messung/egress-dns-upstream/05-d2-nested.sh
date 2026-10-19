#!/bin/bash
# D2, verschachtelter Podman (roter Test NestedPodman_BuildEgress): fragt der verschachtelte Container als vscode, und was bewirkt die
# Messvariante (nft meta skuid) auf Aufloesung und Bau?
# Aufruf: WORK=<messordner> OUT=<ausgabeverzeichnis> 05-d2-nested.sh -> $OUT/rohdaten-05-d2-nested.txt
# Voraussetzung: 02-projekte.sh (Image u-boot-messung-dns-img-pod, das Produkt-Image mit nestedRuntime podman). Der Container
# bekommt die runArgs der erzeugten devcontainer.json (Messung 02: rohdaten-02.txt).
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-05-d2-nested.txt
c=$PFX-nested
args=$(sed -n 's/^ *"runArgs": \[\(.*\)\],$/\1/p' "$(vdir pod)/.devcontainer/devcontainer.json" | tr -d '",')
echo "runArgs der erzeugten devcontainer.json (pod): $args"
for mode in ist skuid; do
  echo; echo "######## D2 nested, Modus $mode, $(now)"
  docker rm -f -v $c > /dev/null 2>&1
  # shellcheck disable=SC2086
  docker run -d --init --name $c $args -v "$here":/m:ro $IMG_POD sleep infinity > /dev/null
  docker exec -u vscode $c bash /m/inner-nested.sh $mode
  docker rm -f -v $c > /dev/null
done
