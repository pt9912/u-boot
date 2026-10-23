#!/bin/bash
# Messstand erfassen, und die Liste der Docker-Objekte VOR der Messung.
# Aufruf: WORK=<messordner> ./01-messstand.sh        -> rohdaten-01.txt, rohdaten-08-vorher.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-01.txt

echo "== Messstand, $(now)"
echo "-- Docker (Client und Server)"
docker version --format 'client={{.Client.Version}} api={{.Client.APIVersion}} os/arch={{.Client.Os}}/{{.Client.Arch}}'
docker version --format 'server={{.Server.Version}} api={{.Server.APIVersion}} os/arch={{.Server.Os}}/{{.Server.Arch}}'
docker info --format 'storage-driver={{.Driver}} cgroup={{.CgroupVersion}} ncpu={{.NCPU}} memtotal={{.MemTotal}} security={{.SecurityOptions}}'
echo "-- Kernel, CPU, RAM (aus einem Container: der Kernel ist geteilt)"
docker run --rm "$BASE" sh -c '
  echo "kernel=$(uname -sr)  arch=$(uname -m)"
  echo "cpu=$(grep -m1 "model name" /proc/cpuinfo | cut -d: -f2 | sed "s/^ //")  anzahl=$(nproc)"
  echo "ram=$(awk "/MemTotal/ {printf \"%.1f GiB\", \$2/1048576}" /proc/meminfo)"'
echo "-- Basis-Image (Vorlage Dockerfile.tmpl: FROM $BASE)"
docker pull -q "$BASE" > /dev/null
docker image inspect "$BASE" --format 'id={{.Id}} digest={{index .RepoDigests 0}} created={{.Created}} arch={{.Architecture}} size={{.Size}}'
docker run --rm "$BASE" sh -c '. /etc/os-release; echo "$PRETTY_NAME"; curl --version | head -1; git --version; jq --version'
echo "-- Zeile FROM der Vorlage und Zeile des Credential-Helpers (gelesen)"
grep -n -E '^FROM |credential.helper' "$ROOT/internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl"
echo "-- Netz: Erreichbarkeit der Forge-Hosts (HTTP-Status der Startseite, ohne Anmeldung)"
for h in api.github.com github.com gitlab.com ghcr.io cli.github.com; do
  printf '%s: ' "$h"; tool curl -sS -m 15 -o /dev/null -w '%{http_code}\n' "https://$h/"
done
aussage "docker-server-arch=$(docker version --format '{{.Server.Arch}}')"
aussage "basis-image-digest=$(docker image inspect "$BASE" --format '{{index .RepoDigests 0}}')"
echo
echo "== Docker-Objekte VOR der Messung (auch in rohdaten-08-vorher.txt)"
bestand | tee "$OUT/rohdaten-08-vorher.txt"
