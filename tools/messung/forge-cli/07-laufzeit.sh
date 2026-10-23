#!/bin/bash
# Laufzeit im Container (Image einer Variante) ohne und hinter der Egress-Restriktion: Hosts, Anmeldung per Umgebungsvariable, Dateien, Git.
# Voraussetzung: 05-bau.sh (und 04-projekte.sh fuer die Variante ega). Der Container traegt --cap-add=NET_ADMIN wie die erzeugten runArgs und
# kein Volume (Zustand je Lauf frisch). Kein echtes Zugangsdatum: die Skripte setzen nur den Platzhalter MESSUNG-KEIN-ECHTES-TOKEN.
# Modi: discover (ohne Restriktion, Namen protokolliert), restricted (erzeugtes Startscript der Variante), restricted-allow (nur glab: das von
#   `u-boot generate devcontainer` mit `egress.allow gitlab.com` erzeugte Startscript der Variante ega ersetzt das des Images, lesend eingehaengt).
# Aufruf: WORK=<messordner> ./07-laufzeit.sh [variante:werkzeug...]    (Vorgabe: ghk:gh ghxk:gh glabsk:glab glabpk:glab glabzk:glab)
#   -> $OUT/rohdaten-07-<variante>-<modus>.txt
set -u
. "$(dirname "$0")/lib.sh"

pairs=("$@"); [ ${#pairs[@]} -gt 0 ] || pairs=(ghk:gh ghxk:gh glabsk:glab glabpk:glab glabzk:glab)
for p in "${pairs[@]}"; do
  v=${p%%:*}; t=${p##*:}
  # Aus devcontainer.json: runArgs der Variante (zur Kontrolle gegen den Aufruf unten)
  runargs=$(grep -F -h '"runArgs"' "$(vdir "$v")/.devcontainer/devcontainer.json" | tr -d ' ')
  modes=(discover restricted); [ "$t" = glab ] && modes+=(restricted-allow)
  for mode in "${modes[@]}"; do
    (
      rohdaten "rohdaten-07-$v-$mode.txt"
      echo "== Laufzeit $v ($t), Modus $mode, $(now)"
      echo "runArgs laut devcontainer.json: $runargs"
      mnt=(); [ "$mode" = restricted-allow ] && mnt=(-v "$(vdir ega)/.devcontainer/egress-init.sh":/usr/local/bin/u-boot-egress-init:ro)
      s=$(date +%s.%N)
      docker run --rm --cap-add=NET_ADMIN -e DUMMY="$DUMMY" -v "$here":/m:ro "${mnt[@]}" "$(vimg "$v")" bash /m/inner-07.sh "$mode" "$t"
      echo "container-exit=$?  gesamtdauer=$(sekunden "$s" "$(date +%s.%N)")s"
    )
  done
done
