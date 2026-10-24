#!/bin/bash
# Relay fuer Weg (a), K1b (HANDVARIANTE, kein Teil der Auslieferung): ein Hilfs-Container mit --network host reicht TCP-Verbindungen auf
# <Gateway-Adresse>:11434 an 127.0.0.1:11434 des Hosts weiter. Die Gateway-Adresse ist die, die host-gateway in einem Container ergibt
# (getent im Container); gebunden wird nur diese Adresse, nie 0.0.0.0. Der Relay laeuft nur waehrend eines Falls.
#   start   prueft, dass kein fremder Container im Netz bridge haengt, startet den Relay (Image per Digest), prueft die Erreichbarkeit
#   stop    entfernt den Relay und weist per docker ps nach, dass er weg ist
# Solange der Relay laeuft, kann jeder Container im Netz bridge den angemeldeten Server des Hosts nutzen; deshalb die Pruefung beim Start.
# Voraussetzung: 00-vorlauf.sh. Aufruf: WORK=<messordner> tools/messung/ollama-cloud-konto/04-relay.sh start|stop
set -u
. "$(dirname "$0")/lib.sh"
vorlauf_pruefen || exit 2
RELAY_IMAGE='alpine/socat@sha256:603174b77eebaba25ba8f8ebb15e1234d6600409554e95d6fb203fd106150287'   # gelesen beim ersten Pull von alpine/socat:latest am 2026-10-24
modus=${1:-}
case "$modus" in
start)
  rohdaten "rohdaten-04-relay-start-$(date -u +%H%M%S).txt"
  echo "== Relay start, $(now)"
  docker rm -f "$RELAY" > /dev/null 2>&1
  fremde=0; pat="^(vsc-)?$PFX"; lpat="/$PFX-(std|res|alw)\$"
  # zur Messung gehoeren: Namen mit dem Praefix (auch die Filter-Container der Schwaerzung) und Container der Dev-Container-CLI mit dem Label
  # devcontainer.local_folder auf ein Projekt der Messung
  for id in $(docker network inspect bridge --format '{{range .Containers}}{{.Name}} {{end}}'); do
    lf=$(docker inspect -f '{{index .Config.Labels "devcontainer.local_folder"}}' "$id" 2> /dev/null)
    [[ $id =~ $pat ]] || [[ $lf =~ $lpat ]] || fremde=$((fremde + 1))
  done
  echo "container im Netz bridge, die nicht zur Messung gehoeren: $fremde (Namen nicht ausgewiesen)"
  [ "$fremde" = 0 ] || { echo "Abbruch: fremde Container im Messnetz"; exit 3; }
  read -r adr _ < <(docker run --rm --add-host=host.docker.internal:host-gateway "$BASE" getent hosts host.docker.internal)
  echo "host-gateway im Container ergibt: $adr"
  [ -n "$adr" ] || { echo "Abbruch: keine Adresse"; exit 3; }
  echo "image: $RELAY_IMAGE"
  docker image inspect "$RELAY_IMAGE" > /dev/null 2>&1 || docker pull -q "$RELAY_IMAGE" > /dev/null
  docker run -d --name "$RELAY" --network host --restart no "$RELAY_IMAGE" "TCP-LISTEN:$PORT,bind=$adr,fork,reuseaddr" "TCP:127.0.0.1:$PORT" > /dev/null || { echo "Relay nicht gestartet"; exit 1; }
  sleep 1
  echo "-- Erreichbarkeit aus einem Container mit add-host (api/version, kein Prompt)"
  docker run --rm --add-host=host.docker.internal:host-gateway "$BASE" curl -sS -m 5 "http://host.docker.internal:$PORT/api/version"; echo " rc=$?"
  echo "-- Relay lauscht auf (ss im Host-Netz-Container, nur Zeilen mit Port $PORT)"
  docker run --rm --network host "$BASE" sh -c "ss -ltn | grep ':$PORT '"
  aussage "relay-gateway-adresse=$adr relay-gebunden-nur-an-diese-adresse=$(docker run --rm --network host "$BASE" sh -c "ss -ltn | grep -c '0.0.0.0:$PORT'" | grep -q '^0$' && echo ja || echo nein)"
  ;;
stop)
  rohdaten "rohdaten-04-relay-stop-$(date -u +%H%M%S).txt"
  echo "== Relay stop, $(now)"
  docker rm -f "$RELAY" > /dev/null 2>&1
  echo "docker ps -a nach dem Entfernen: Relay-Container vorhanden: $(docker ps -a --format '{{.Names}}' | grep -c -x "$RELAY")"
  echo "lauscht noch etwas auf $PORT ausser 127.0.0.1: $(docker run --rm --network host "$BASE" sh -c "ss -ltn | grep ':$PORT ' | grep -v '127.0.0.1:$PORT' | wc -l")"
  aussage "relay-entfernt=$(docker ps -a --format '{{.Names}}' | grep -c -x "$RELAY" | grep -q '^0$' && echo ja || echo nein)"
  ;;
*) echo "Aufruf: 04-relay.sh start|stop"; exit 2 ;;
esac
