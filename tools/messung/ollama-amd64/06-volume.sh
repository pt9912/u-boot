#!/bin/bash
# A5: Modell-Volume. Container mit den `mounts` der erzeugten devcontainer.json (${devcontainerId} durch eine feste Kennung ersetzt),
# Pull, Container loeschen, neuer Container mit demselben Volume, Neubau des Images von Hand mit demselben Volume-Namen.
# Voraussetzung: 03-bau.sh. Aufruf: WORK=<messordner> tools/messung/ollama-amd64/06-volume.sh      -> $OUT/rohdaten-06.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-06.txt

ID=m1
C1=$PFX-c1; C2=$PFX-c2; C3=$PFX-c3
IMG=$(vimg oll)
aufraeumen() { docker rm -f -v "$C1" "$C2" "$C3" > /dev/null 2>&1; }
volumes() { docker volume ls -q | grep "^$PFX-oll-" | sort; }
aufraeumen
for v in $(volumes); do docker volume rm "$v" > /dev/null; done

echo "== A5 Volume, $(now)"
echo "-- mounts laut devcontainer.json (Reihenfolge wie dort)"
mapfile -t MOUNTS < <(tool sh -c "grep -v '^[[:space:]]*//' p-oll/.devcontainer/devcontainer.json | jq -r '.mounts[]'" | tr -d '\r')
for m in "${MOUNTS[@]}"; do echo "  $m"; done
ARGS=(); for m in "${MOUNTS[@]}"; do ARGS+=(--mount "${m//\$\{devcontainerId\}/$ID}"); done
aussage "mounts-anzahl=${#MOUNTS[@]}"
aussage "mount-2-ziel=$(echo "${MOUNTS[1]}" | grep -o 'target=[^,]*')"

start() { docker run -d --name "$1" --cap-add=NET_ADMIN "${ARGS[@]}" "$2" sleep 3600 > /dev/null; }
serve='sudo sh /usr/local/bin/u-boot-egress-init && (nohup ollama serve > /tmp/ollama.log 2>&1 &) ; for i in $(seq 1 60); do curl -fsS -m 2 http://127.0.0.1:11434/ >/dev/null 2>&1 && break; sleep 0.5; done'

mark "Container 1: neues Volume, Eigentuemer, Pull hinter der Restriktion"
start "$C1" "$IMG"
docker exec "$C1" bash -c 'stat -c "%U:%G %a %n" /home/vscode/.ollama /workspaces/*; ls -la /home/vscode/.ollama | tail -n +2'
aussage "volume-eigentuemer-vor-dem-start=$(docker exec "$C1" stat -c '%U:%G' /home/vscode/.ollama)"
docker exec "$C1" bash -c "$serve; ollama pull smollm2:135m 2>&1 | sed -e 's/\x1b\[[0-9;?]*[a-zA-Z]//g' | tr '\r' '\n' | grep -E 'success|Error' | tail -1; ollama list"
echo "-- Schluesselpaar (nur Name, Rechte, Eigentuemer; Pruefsumme des PUBLIC-Teils)"
docker exec "$C1" bash -c 'stat -c "%U:%G %a %n" /home/vscode/.ollama/id_ed25519 /home/vscode/.ollama/id_ed25519.pub; sha256sum /home/vscode/.ollama/id_ed25519.pub | cut -c1-16 | sed "s/^/pub-sha256-anfang=/"'
K1=$(docker exec "$C1" sha256sum /home/vscode/.ollama/id_ed25519.pub | cut -d' ' -f1)
echo "-- Volumes nach dem Start"; volumes
docker volume inspect "$PFX-oll-ollama-$ID" --format 'volume {{.Name}} driver={{.Driver}} scope={{.Scope}}'

mark "Container 1 loeschen (docker rm -f), Volume bleibt"
docker rm -f "$C1" > /dev/null
volumes
mark "Container 2: neuer Container, dasselbe Volume"
start "$C2" "$IMG"
docker exec "$C2" bash -c "$serve; ollama list; echo -- Lauf ohne neuen Pull; ollama run smollm2:135m 'Reply with the single word OK.' 2>&1 | sed -e 's/\x1b\[[0-9;?]*[a-zA-Z]//g' | tr '\r' '\n' | grep -v '^ *\$' | tail -1"
K2=$(docker exec "$C2" sha256sum /home/vscode/.ollama/id_ed25519.pub | cut -d' ' -f1)
aussage "modell-nach-rm-gelistet=$(docker exec "$C2" bash -c 'ollama list | grep -c "^smollm2:135m"')"
aussage "schluessel-nach-rm-gleich=$([ "$K1" = "$K2" ] && echo ja || echo nein)"
docker rm -f "$C2" > /dev/null

mark "Neubau des Images von Hand (angehaengte Schicht), dasselbe Volume"
rm -rf "$WORK/p-oll-neubau"; mkdir -p "$WORK/p-oll-neubau"
cp "$(vdir oll)/.devcontainer/"* "$WORK/p-oll-neubau/"
printf '\nUSER root\nRUN echo neubau > /etc/u-boot-messung-neubau\nUSER vscode\n' >> "$WORK/p-oll-neubau/Dockerfile"
docker build -q -t "$IMG-neubau" "$WORK/p-oll-neubau" > /dev/null; echo "build-neubau exit=$?"
echo "image alt: $(docker image inspect "$IMG" --format '{{.Id}}')"
echo "image neu: $(docker image inspect "$IMG-neubau" --format '{{.Id}}')"
aussage "neubau-image-id-verschieden=$([ "$(docker image inspect "$IMG" --format '{{.Id}}')" != "$(docker image inspect "$IMG-neubau" --format '{{.Id}}')" ] && echo ja || echo nein)"
start "$C3" "$IMG-neubau"
docker exec "$C3" bash -c "$serve; ollama list; cat /etc/u-boot-messung-neubau"
K3=$(docker exec "$C3" sha256sum /home/vscode/.ollama/id_ed25519.pub | cut -d' ' -f1)
aussage "modell-nach-neubau-gelistet=$(docker exec "$C3" bash -c 'ollama list | grep -c "^smollm2:135m"')"
aussage "schluessel-nach-neubau-gleich=$([ "$K1" = "$K3" ] && echo ja || echo nein)"
echo "-- Nicht gemessen: ob \${devcontainerId} beim Neuaufbau durch die Dev-Container-CLI stabil bleibt (A9)."
aufraeumen
docker rmi "$IMG-neubau" > /dev/null
for v in $(volumes); do docker volume rm "$v" > /dev/null; done
echo "-- Volumes danach"; volumes || true
