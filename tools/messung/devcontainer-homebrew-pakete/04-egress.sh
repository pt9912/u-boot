#!/bin/bash
# M4 Laufzeit hinter der Egress-Restriktion: welche Hosts braucht `brew install` im laufenden Container? Image der Variante egr (Sandbox-Profil,
# Restriktion an, Homebrew allein) in vier Modi (inner-04.sh): discover (ohne Restriktion, Namen protokolliert), restricted (erzeugtes Startscript
# mit der Zeile des Features), restricted-ohne (die drei Hosts der Zeile des Features aus der Liste entfernt) und restricted-formelhost (die
# Liste plus der Host, von dem die gemessene Tap-Formel laedt). Der Container traegt --cap-add=NET_ADMIN wie die erzeugten runArgs und kein Volume.
# Voraussetzung: 01-bau.sh (Variante egr).
# Aufruf: UBOOT=<binary> WORK=<messordner> ./04-egress.sh [variante:modus...]   (Vorgabe: egr:discover egr:restricted egr:restricted-ohne egr:restricted-formelhost)
#   -> rohdaten-m4-<variante>-<modus>.txt
set -u
. "$(dirname "$0")/lib.sh"

pairs=("$@"); [ ${#pairs[@]} -gt 0 ] || pairs=(egr:discover egr:restricted egr:restricted-ohne egr:restricted-formelhost)
for p in "${pairs[@]}"; do
  v=${p%%:*}; mode=${p##*:}
  docker image inspect "$(vimg "$v")" > /dev/null 2>&1 || { echo "Abbruch: das Image der Variante $v fehlt (zuerst 01-bau.sh)"; exit 2; }
  (
    rohdaten "rohdaten-m4-$v-$mode.txt"
    echo "== M4 Laufzeit $v, Modus $mode, $(now)"
    echo "runArgs laut devcontainer.json: $(grep -F -h '"runArgs"' "$(vdir "$v")/.devcontainer/devcontainer.json" | tr -d ' ')"
    s=$(date +%s)
    docker run --rm --cap-add=NET_ADMIN -v "$here":/m:ro "$(vimg "$v")" bash /m/inner-04.sh "$mode"
    echo "container-exit=$?  gesamtdauer=$(dauer "$s" "$(date +%s)")s"
  )
done
