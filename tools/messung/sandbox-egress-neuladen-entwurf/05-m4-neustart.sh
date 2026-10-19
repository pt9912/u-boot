#!/bin/bash
# M4: Liste aendern, `u-boot generate devcontainer`, Container neu starten ohne Neubau: gilt danach die alte oder die neue Liste?
# Aufruf: 05-m4-neustart.sh <repo-wurzel> <arbeitsverzeichnis> > rohdaten-05-m4.txt 2>&1
# Voraussetzung: Docker (Colima); Binary <repo-wurzel>/bin/u-boot-linux-arm64 (`make build-binaries`); Image
# u-boot-messung-devcontainer-cli (Dockerfile.cli: Dev-Container-CLI 0.89.0 mit Docker-CLI, spricht mit dem Docker-Socket).
# Das Projekt entsteht wie in 01-projekt-erzeugen.sh; der Container wird von der Dev-Container-CLI angelegt, nicht von Hand:
# der postStartCommand laeuft damit so, wie ihn ein Entwickler-Werkzeug ausfuehrt.
set -u
root=$(cd "${1:?repo-wurzel}" && pwd)
w=${2:?arbeitsverzeichnis}
rm -rf "$w"; mkdir -p "$w"; w=$(cd "$w" && pwd)
folder=/w/proj
CLI=u-boot-messung-devcontainer-cli

ub() { docker run --rm -v "$root/bin":/b:ro -v "$w":/w -w $folder -u "$(id -u):$(id -g)" -e HOME=/tmp debian:stable-slim /b/u-boot-linux-arm64 "$@" 2>&1 | grep -v -E '^(time=|Warning)'; }
dc() { docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "$w":/w $CLI devcontainer "$@" --workspace-folder $folder 2>&1 | tail -n 3 | cut -c1-300; }
cid() { docker ps -a -q --filter "label=devcontainer.local_folder=$folder" | head -1; }
state() {
  c=$(cid)
  echo "    Container $c, Running=$(docker inspect -f '{{.State.Running}}' "$c"), Image $(docker inspect -f '{{.Image}}' "$c" | cut -c1-19), erzeugt $(docker inspect -f '{{.Created}}' "$c")"
  docker exec -u root "$c" sh -c '
    echo "    gebackene Liste (hosts= in /usr/local/bin/u-boot-egress-init): $(sed -n "s/^hosts=\"\(.*\)\"\$/\1/p" /usr/local/bin/u-boot-egress-init)"
    echo "    nftset-Namen der laufenden dnsmasq-Konfiguration: $(sed -n "s|^nftset=/\([^/]*\)/.*|\1|p" /etc/u-boot-egress-dnsmasq.conf | tr "\n" " ")"
    echo "    Statusdatei: $(tr "\n" " " < /var/lib/u-boot/status/egress)"
    echo "    https://example.org (neuer Eintrag): $(curl -sS -m 6 -o /dev/null -w "http=%{http_code}" https://example.org 2>&1 | tail -1 | cut -c1-70)"
    echo "    https://deb.debian.org (Standardeintrag): $(curl -sS -m 6 -o /dev/null -w "http=%{http_code}" https://deb.debian.org 2>&1 | tail -1 | cut -c1-70)"
    echo "    Zombie-Prozesse: $(ps -eo stat= | grep -c "^Z")"
    echo "    nft-Tabellen: [$(nft list tables | tr "\n" " ")]"' 2>&1
}

echo "######## 0: Projekt erzeugen (Standardliste) und Container mit der Dev-Container-CLI starten"
ub init egressprobe --devcontainer --sandbox --no-git --yes >/dev/null
mkdir -p "$w/proj" 2>/dev/null
ub config set devcontainer.sandbox.egress.enabled true | tail -1
ub generate devcontainer --yes | tail -1
dc up
echo "-- Zustand nach dem ersten Start"; state

echo "######## 1: Liste aendern (example.org dazu) und neu erzeugen"
ub config set devcontainer.sandbox.egress.allow example.org | tail -1
ub generate devcontainer --yes | tail -1
docker run --rm -v "$w":/w debian:stable-slim sh -c 'echo "    Host-Datei .devcontainer/egress-init.sh: $(sed -n "s/^hosts=\"\(.*\)\"\$/\1/p" /w/proj/.devcontainer/egress-init.sh)"'
echo "-- Zustand des Containers (unveraendert, er wurde nicht angefasst)"; state

echo "######## 2: Container stoppen und mit der Dev-Container-CLI wieder starten (kein Neubau, der Container bleibt)"
docker stop -t 2 "$(cid)" >/dev/null
dc up
echo "-- Zustand nach Stop und up"; state

echo "######## 3: devcontainer up gegen den laufenden Container (Dockerfile-Kontext hat sich geaendert)"
dc up
echo "-- Zustand nach up"; state

echo "######## 4: devcontainer up --remove-existing-container (Container neu angelegt; baut die CLI das Image neu?)"
dc up --remove-existing-container
echo "-- Zustand nach up --remove-existing-container"; state

echo "######## Aufraeumen"
docker rm -f "$(cid)" >/dev/null 2>&1
echo "fertig"
