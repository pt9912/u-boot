#!/bin/bash
# A3 und A4: Lauf im Container (Image oll) ohne und hinter der Egress-Restriktion; Modelle smollm2:135m und qwen2.5:0.5b.
# Voraussetzung: 03-bau.sh. Der Container traegt --cap-add=NET_ADMIN wie die erzeugten runArgs und kein Volume (A5: 06-volume.sh).
# Aufruf: WORK=<messordner> tools/messung/ollama-amd64/05-laufzeit.sh
#   -> $OUT/rohdaten-05-discover.txt, rohdaten-05-restricted.txt
set -u
. "$(dirname "$0")/lib.sh"

# Aus devcontainer.json: runArgs der Variante oll (zur Kontrolle gegen den Aufruf unten)
grep -F -h '"runArgs"' "$(vdir oll)/.devcontainer/devcontainer.json" > "$WORK/runargs-oll.txt"

for mode in discover restricted; do
  (
    rohdaten "rohdaten-05-$mode.txt"
    echo "== A3/A4 Modus $mode, $(now)"
    echo "runArgs laut devcontainer.json: $(cat "$WORK/runargs-oll.txt")"
    s=$(date +%s.%N)
    docker run --rm --cap-add=NET_ADMIN -v "$here":/m:ro "$(vimg oll)" bash /m/inner-05.sh "$mode" smollm2:135m qwen2.5:0.5b
    echo "container-exit=$?  gesamtdauer=$(sekunden "$s" "$(date +%s.%N)")s"
  )
done
