#!/bin/bash
# Messung 2: Dauer von `devcontainer up` im Hilfs-Container mit der Argumentliste des Adapters — erster Start mit kaltem
# Basis-Image des Projekts, Wiederholung gegen den laufenden Container, Neustart nach `docker stop`, `--remove-existing-container`.
# Aufruf: 02-erster-up.sh <repo-wurzel> <arbeitsverzeichnis> > rohdaten-02-erster-up.txt 2>&1
# Voraussetzung: Linux-Host mit Docker, Bash und GNU-Werkzeugen (tar mit --mtime/--owner/--group, date +%N, awk, mktemp, head, cut);
# auf macOS nicht belegt (siehe 01-hilfs-image.sh). Das Hilfs-Image wird aus dem eingebetteten Kontext gebaut (Tag
# u-boot-messung-hilfs-image:messung). Das Projekt ist handgeschrieben: ein Image ohne Features und ohne Startscripte.
# Das Basis-Image `base:ubuntu` darf lokal nicht vorhanden sein (sonst bricht das Skript ab): "kalt" heisst, der erste `up`
# zieht es. Nach dem Lauf werden Container, Image und Volume dieser Messung entfernt, die Basis-Images des Hilfs-Images bleiben.
set -u
root=$(cd "${1:?repo-wurzel}" && pwd)
w=${2:?arbeitsverzeichnis}
rm -rf "$w"; mkdir -p "$w"; w=$(cd "$w" && pwd)
ctx="$root/internal/adapter/driven/docker/devcontainercli"
tag=u-boot-messung-hilfs-image:messung
base=mcr.microsoft.com/devcontainers/base:ubuntu
folder="$w/projekt"

if docker image inspect "$base" >/dev/null 2>&1; then
  echo "abgebrochen: $base ist lokal vorhanden (kein kalter Start messbar)"; exit 1
fi

now() { date +%s.%N; }
dauer() { awk -v a="$1" -v b="$2" 'BEGIN { printf "%.1f s", b - a }'; }
cid() { docker ps -a -q --filter "label=devcontainer.local_folder=$folder" | head -1; }
# Argumentliste wie runArgs im Adapter (ohne -e, Projektordner schreibend).
up() {
  docker run --rm --name "u-boot-devcontainer-up-messung" \
    --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock \
    --mount type=bind,source="$folder",target="$folder" \
    "$tag" devcontainer up --workspace-folder "$folder" "$@" 2>/dev/null | tail -n 1 | cut -c1-120
}
messen() {
  local was=$1; shift
  local t0 t1; t0=$(now); up "$@"; t1=$(now)
  echo "$was: dauer=$(dauer "$t0" "$t1") container=$(cid | cut -c1-12)"
}
aufraeumen() {
  docker rm -f "$(cid)" >/dev/null 2>&1
  docker rmi -f "$(docker images -q --filter "reference=vsc-projekt-*")" >/dev/null 2>&1
  docker rmi "$base" "$tag" >/dev/null 2>&1
}
trap aufraeumen EXIT

echo "######## 0: Messstand"
echo "datum_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
echo "host_arch=$(uname -m)"
docker version --format 'client={{.Client.Version}} server={{.Server.Version}} server_os_arch={{.Server.Os}}/{{.Server.Arch}}'

mkdir -p "$folder/.devcontainer"
printf '{ "name": "messung", "image": "%s" }\n' "$base" > "$folder/.devcontainer/devcontainer.json"
tar_ctx=$(mktemp)
tar -C "$ctx" --mtime='@0' --owner=0 --group=0 --numeric-owner -cf "$tar_ctx" Dockerfile package.json package-lock.json
docker build -q -t "$tag" - < "$tar_ctx" >/dev/null 2>&1 || { echo "Bau des Hilfs-Images fehlgeschlagen"; exit 1; }
rm -f "$tar_ctx"

echo "######## 1: erster up (Basis-Image des Projekts wird gezogen, Container wird angelegt)"
messen "erster up (kalt)"
echo "######## 2: up gegen den laufenden Container"
messen "up, laeuft"
echo "######## 3: Container stoppen, dann up"
docker stop -t 2 "$(cid)" >/dev/null
messen "up nach stop"
echo "######## 4: up --remove-existing-container"
messen "up --remove-existing-container" --remove-existing-container
echo "######## Aufraeumen"
echo "fertig"
