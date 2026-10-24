#!/bin/bash
# C1 bis C6: Lauf im Container, der aus dem erzeugten Projekt entsteht (Dev-Container-CLI `devcontainer up`, wie ausgeliefert).
#   discover    Variante res: die CLI baut das Image und startet den Container (Egress-Skript als postStartCommand); danach ein ZWEITER
#               Container aus demselben Image per `docker run --cap-add=NET_ADMIN` OHNE postStartCommand (Handvariante des Aufrufs:
#               runArgs wie im Projekt, kein Workspace-Volume, kein Modell-Volume ausser dem unten benannten) -> Modus discover
#   restricted  der Container der CLI aus Variante res (Liste des Generators)
#   allow       Container der CLI aus Variante alw (Liste plus ollama.com, erzeugt ueber `config set`)
# Der Aufruf der Messung ist inner.sh (wird per `docker cp` in den Container kopiert). Voraussetzung: 02-projekte.sh, 03-lauf-vorbereitung.sh.
# Aufruf: WORK=<messordner> tools/messung/ollama-cloud/04-lauf.sh [modell]
#   LAUF=<n> setzt die Nummer des Messlaufs (Vorgabe 1) -> $OUT/rohdaten-04-<n>.txt (Treiber), rohdaten-04-<modus>-<n>.txt,
#   rohdaten-04-up-<variante>-<n>.txt
set -u
. "$(dirname "$0")/lib.sh"
LAUF=${LAUF:-1}
rohdaten "rohdaten-04-$LAUF.txt"
M=${1:-gpt-oss:120b-cloud}   # Name eines Cloud-Modells aus der Bibliothek von Ollama; ob es existiert, ist nicht Gegenstand
umgebung_pruefen || { echo "Abbruch: verbotene Namen in der Umgebung"; exit 2; }

cli() { # <variante> <argumente...>
  local v=$1 P; P=$(vname "$1"); shift
  docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "$(vdir "$v")":/w/"$P" -w /w/"$P" "$CLI" devcontainer "$@" --workspace-folder /w/"$P"
}
entfernen() { # <variante>
  docker ps -a -q --filter "label=devcontainer.local_folder=/w/$(vname "$1")" | while read -r c; do docker rm -f -v "$c" > /dev/null; done
}
hoch() { # <variante>: Container der CLI starten, Laufzeit und Exit ausweisen
  local v=$1 s; s=$(date +%s.%N)
  entfernen "$v"
  cli "$v" up --log-level info --remove-existing-container 2>&1 | grep -v -E '^\[[0-9]+ ms\]|^\s*$' | tail -12 | redact > "$OUT/rohdaten-04-up-$v-$LAUF.txt"
  local rc=${PIPESTATUS[0]}
  echo "devcontainer up $v exit=$rc dauer=$(sekunden "$s" "$(date +%s.%N)")s"
  aussage "up-$v-exit=$rc"
  docker ps -q --filter "label=devcontainer.local_folder=/w/$(vname "$v")" | head -1
}
lauf() { # <container> <modus> <datei>
  docker cp "$here/inner.sh" "$1":/tmp/inner.sh
  docker cp "$here/claude-stub.sh" "$1":/tmp/claude-stub.sh
  ( rohdaten "$3"; echo "== Modus $2, Modell $M, $(now)"; docker exec -u vscode "$1" bash /tmp/inner.sh "$2" "$M"; echo "docker exec exit=$?" )
}

echo "== C1 bis C6, $(now)"
echo "-- Variante res: devcontainer up"
c=$(hoch res | tee /dev/stderr | tail -1)
[ -n "$c" ] || { echo "kein Container"; exit 1; }
img=$(docker inspect -f '{{.Config.Image}}' "$c")
echo "image des Containers: $img"
echo "-- Egress-Status im Container der CLI"; docker exec "$c" cat /var/lib/u-boot/status/egress
echo "-- Einhaengepunkte des Containers (Quelle -> Ziel, Art)"
docker inspect -f '{{range .Mounts}}{{.Type}} {{.Destination}}{{"\n"}}{{end}}' "$c"

echo "-- Modus discover: zweiter Container aus demselben Image (Handvariante, siehe Kopf)"
d=$PFX-discover
docker rm -f -v "$d" > /dev/null 2>&1
docker run -d --name "$d" --cap-add=NET_ADMIN -v "$PFX-modelle:/home/vscode/.ollama" --entrypoint sleep "$img" 1800 > /dev/null
lauf "$d" discover "rohdaten-04-discover-$LAUF.txt"
docker rm -f -v "$d" > /dev/null; docker volume rm "$PFX-modelle" > /dev/null

echo "-- Modus restricted: Container der CLI (Variante res)"
lauf "$c" restricted "rohdaten-04-restricted-$LAUF.txt"
entfernen res

echo "-- Modus allow: Variante alw (egress.allow ollama.com)"
c=$(hoch alw | tee /dev/stderr | tail -1)
[ -n "$c" ] || { echo "kein Container"; exit 1; }
lauf "$c" restricted "rohdaten-04-allow-$LAUF.txt"
entfernen alw
echo "== Ende, $(now)"
