#!/bin/bash
# A9: Dev-Container-CLI im Wegwerf-Container (Docker-Socket nur fuer die Messung): `devcontainer up` mit dem erzeugten Projekt (oll),
# Modell laden, `devcontainer up --remove-existing-container --build-no-cache`, Verhalten des Volumes <projekt>-ollama-${devcontainerId}.
# Die CLI vergibt Container- und Image-Namen selbst (Image: vsc-<ordnername>-<hash>, Container: zufaellig); der Ordnername traegt das
# Praefix der Messung, die Container werden ueber das Label devcontainer.local_folder gefunden und von diesem Skript entfernt.
# Entfaellt mit Vermerk, wenn das Hilfs-Image nicht baut. Voraussetzung: 02-projekte.sh.
# Aufruf: WORK=<messordner> tools/messung/ollama-amd64/10-devcontainer-cli.sh      -> $OUT/rohdaten-10.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-10.txt

CLI=$PFX-devcontainer-cli
P=$(vname oll)
FOLDER=/w/$P
label="devcontainer.local_folder=$FOLDER"
cli() { docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "$(vdir oll)":"$FOLDER" -w "$FOLDER" "$CLI" "$@"; }
entfernen() {
  docker ps -a -q --filter "label=$label" | while read -r c; do docker rm -f -v "$c" > /dev/null; done
  docker images --format '{{.Repository}}:{{.Tag}}' | grep "^vsc-$P" | while read -r i; do docker rmi "$i" > /dev/null; done
  for x in $(docker volume ls -q | grep "^$P-"); do docker volume rm "$x" > /dev/null; done
}

echo "== A9 Dev-Container-CLI, $(now)"
if ! docker build -q -f "$here/Dockerfile.cli" -t "$CLI" "$here" > "$WORK/cli-build.log" 2>&1; then
  echo "ENTFAELLT: das Hilfs-Image baut nicht (Exit ungleich 0); letzte Zeilen:"; tail -5 "$WORK/cli-build.log"
  aussage "a9=entfaellt-hilfsimage-baut-nicht"; exit 0
fi
echo "devcontainer-cli $(cli devcontainer --version)"
aussage "a9-cli-version=$(cli devcontainer --version)"
entfernen

mark "devcontainer up (erster Aufbau)"
s=$(date +%s.%N)
cli devcontainer up --workspace-folder "$FOLDER" --log-level info 2>&1 | grep -v -E '^\[[0-9]+ ms\]|^\s*$' | tail -25
echo "up-exit=${PIPESTATUS[0]} dauer=$(sekunden "$s" "$(date +%s.%N)")s"
C1=$(docker ps -q --filter "label=$label")
echo "container laeuft: $([ -n "$C1" ] && echo ja || echo nein)"
aussage "a9-up-container-laeuft=$([ -n "$C1" ] && echo ja || echo nein)"
[ -n "$C1" ] || { entfernen; exit 0; }
echo "-- Volumes (Name mit devcontainerId)"
docker volume ls -q | grep "^$P-"
V1=$(docker volume ls -q | grep "^$P-ollama-")
echo "-- Egress-Skript lief als postStartCommand? (Status-Datei)"
docker exec "$C1" cat /var/lib/u-boot/status/egress
docker exec -u vscode "$C1" bash -c '(nohup ollama serve > /tmp/ollama.log 2>&1 &); for i in $(seq 1 60); do curl -fsS -m 2 http://127.0.0.1:11434/ >/dev/null 2>&1 && break; sleep 0.5; done; ollama pull smollm2:135m 2>&1 | sed -e "s/\x1b\[[0-9;?]*[a-zA-Z]//g" | tr "\r" "\n" | grep -E "success|Error" | tail -1; ollama list'

mark "devcontainer up --remove-existing-container --build-no-cache (Neuaufbau)"
s=$(date +%s.%N)
cli devcontainer up --workspace-folder "$FOLDER" --remove-existing-container --build-no-cache --log-level info 2>&1 | grep -v -E '^\[[0-9]+ ms\]|^\s*$' | tail -8
echo "up-exit=${PIPESTATUS[0]} dauer=$(sekunden "$s" "$(date +%s.%N)")s"
C2=$(docker ps -q --filter "label=$label")
echo "neuer container: $([ -n "$C2" ] && [ "$C2" != "$C1" ] && echo ja || echo nein)"
echo "-- Volumes nach dem Neuaufbau"
docker volume ls -q | grep "^$P-"
V2=$(docker volume ls -q | grep "^$P-ollama-")
aussage "a9-volume-name-gleich-nach-neuaufbau=$([ "$V1" = "$V2" ] && echo ja || echo nein)"
aussage "a9-anzahl-ollama-volumes=$(echo "$V2" | grep -c .)"
docker exec -u vscode "$C2" bash -c '(nohup ollama serve > /tmp/ollama.log 2>&1 &); for i in $(seq 1 60); do curl -fsS -m 2 http://127.0.0.1:11434/ >/dev/null 2>&1 && break; sleep 0.5; done; ollama list'
aussage "a9-modell-nach-neuaufbau=$(docker exec -u vscode "$C2" bash -c 'ollama list | grep -c "^smollm2:135m"')"
entfernen
