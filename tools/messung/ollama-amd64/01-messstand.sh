#!/bin/bash
# A0: Messstand erfassen, und die Liste der Docker-Objekte VOR der Messung (A10, Teil 1).
# Aufruf: WORK=<messordner> ./01-messstand.sh        -> rohdaten-01.txt, rohdaten-11-vorher.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-01.txt

echo "== A0 Messstand, $(now)"
echo "-- Docker (Client und Server)"
docker version --format 'client={{.Client.Version}} api={{.Client.APIVersion}} os/arch={{.Client.Os}}/{{.Client.Arch}}'
docker version --format 'server={{.Server.Version}} api={{.Server.APIVersion}} os/arch={{.Server.Os}}/{{.Server.Arch}}'
docker info --format 'storage-driver={{.Driver}} cgroup={{.CgroupVersion}} ncpu={{.NCPU}} memtotal={{.MemTotal}} security={{.SecurityOptions}} runtimes={{range $k, $v := .Runtimes}}{{$k}} {{end}}'
echo "-- Kernel, CPU, RAM, Geraete (aus einem Container: dieselben Werte wie auf dem Host, der Kernel ist geteilt)"
docker run --rm "$BASE" sh -c '
  echo "kernel=$(uname -sr)  arch=$(uname -m)"
  echo "cpu=$(grep -m1 "model name" /proc/cpuinfo | cut -d: -f2 | sed "s/^ //")  anzahl=$(nproc)"
  echo "ram=$(awk "/MemTotal/ {printf \"%.1f GiB\", \$2/1048576}" /proc/meminfo)"'
docker run --rm --privileged "$BASE" sh -c '
  echo "kernel.apparmor_restrict_unprivileged_userns=$(cat /proc/sys/kernel/apparmor_restrict_unprivileged_userns 2>&1)"
  for d in /dev/fuse /dev/net/tun; do [ -e $d ] && echo "$d vorhanden" || echo "$d fehlt"; done'
echo "-- Plattenplatz vor der Messung (Docker-Datenverzeichnis des Daemons, aus einem Container gelesen)"
docker system df
docker run --rm -v /var/lib/docker:/d:ro "$BASE" df -h /d 2>&1 | tail -2 || true
echo "-- Basis-Image"
docker image inspect "$BASE" --format 'id={{.Id}} digest={{index .RepoDigests 0}} created={{.Created}} arch={{.Architecture}} size={{.Size}}'
docker run --rm "$BASE" sh -c '. /etc/os-release; echo "$PRETTY_NAME"; tar --version | head -1; curl --version | head -1'
echo "-- Ollama (Pins der Vorlage, gelesen aus dem Generator-Quelltext)"
grep -E '^ARG OLLAMA_(VERSION|SHA256_AMD64|SHA256_ARM64)=' \
  "$ROOT/internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl"
aussage "docker-server-arch=$(docker version --format '{{.Server.Arch}}')"
aussage "basis-image-digest=$(docker image inspect "$BASE" --format '{{index .RepoDigests 0}}')"
aussage "ollama-version=$OLLAMA_VERSION"
aussage "sha256-amd64=$SHA_AMD64"

echo
echo "== A10 Teil 1: Docker-Objekte VOR der Messung (Anzahl je Art, Objekte mit dem Praefix; auch in rohdaten-11-vorher.txt)"
bestand | tee "$OUT/rohdaten-11-vorher.txt"
