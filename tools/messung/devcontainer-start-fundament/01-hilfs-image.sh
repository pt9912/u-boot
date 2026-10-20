#!/bin/bash
# Messung 1: Hilfs-Image aus dem eingebetteten Kontext — Plattformen der Basis-Images, Bauzeit kalt und warm, Groesse,
# Trivy HIGH/CRITICAL.
# Aufruf: 01-hilfs-image.sh <repo-wurzel> > rohdaten-01-hilfs-image.txt 2>&1
# Voraussetzung: nur Docker und tar/date auf dem Host. Der Kontext ist der Tar-Strom, den der Adapter an `docker build -`
# uebergibt (dieselben drei Dateien, Reihenfolge Dockerfile, package.json, package-lock.json). "Kalt" heisst: ohne
# Layer-Cache (`--no-cache`), die beiden Basis-Images sind lokal vorhanden; ein Pull der Basis-Images ist nicht Teil.
# Das Image traegt einen Mess-Tag und wird am Ende entfernt.
set -u
root=$(cd "${1:?repo-wurzel}" && pwd)
ctx="$root/internal/adapter/driven/docker/devcontainercli"
tag=u-boot-messung-hilfs-image:messung
trivy=aquasec/trivy:0.70.0   # Version wie TRIVY_VERSION in repo.mk

now() { date +%s.%N; }
dauer() { awk -v a="$1" -v b="$2" 'BEGIN { printf "%.1f s", b - a }'; }

echo "######## 0: Messstand"
echo "datum_utc=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
echo "host_arch=$(uname -m)"
docker version --format 'client={{.Client.Version}} server={{.Server.Version}} server_os_arch={{.Server.Os}}/{{.Server.Arch}}'
tar_ctx=$(mktemp)
tar -C "$ctx" --mtime='@0' --owner=0 --group=0 --numeric-owner -cf "$tar_ctx" Dockerfile package.json package-lock.json
echo "kontext_dateien=$(tar -tf "$tar_ctx" | tr '\n' ' ')"
echo "kontext_bytes=$(stat -c %s "$tar_ctx")"
grep -E '^FROM ' "$ctx/Dockerfile"

echo "######## 1: Plattformen der Basis-Images (Index-Digest gilt fuer alle gelisteten Plattformen)"
for ref in $(sed -n 's/^FROM \([^ ]*\).*/\1/p' "$ctx/Dockerfile"); do
  echo "-- $ref"
  docker buildx imagetools inspect "$ref" 2>&1 | grep -E '^(Name|Digest|MediaType):|Platform:' | grep -v 'unknown/unknown' | cut -c1-200
done

echo "######## 2: Bauzeit kalt (--no-cache) und warm (Layer-Cache), Bau aus Tar-Strom auf stdin"
docker rmi "$tag" >/dev/null 2>&1
t0=$(now); docker build --no-cache -q -t "$tag" - < "$tar_ctx" >/dev/null 2>&1; rc=$?; t1=$(now)
echo "kalt: exit=$rc dauer=$(dauer "$t0" "$t1")"
docker rmi "$tag" >/dev/null 2>&1
t0=$(now); docker build -q -t "$tag" - < "$tar_ctx" >/dev/null 2>&1; rc=$?; t1=$(now)
echo "warm, Tag fehlt (Layer im Cache): exit=$rc dauer=$(dauer "$t0" "$t1")"
t0=$(now); docker build -q -t "$tag" - < "$tar_ctx" >/dev/null 2>&1; rc=$?; t1=$(now)
echo "warm, Tag vorhanden: exit=$rc dauer=$(dauer "$t0" "$t1")"
rm -f "$tar_ctx"

echo "######## 3: Groesse"
docker image inspect "$tag" --format 'groesse_bytes={{.Size}} architektur={{.Architecture}}'
docker image inspect "$tag" --format '{{.Size}}' | awk '{ printf "groesse_mib=%.0f\n", $1 / 1048576 }'
docker run --rm "$tag" devcontainer --version | sed 's/^/devcontainer --version: /'

echo "######## 4: Trivy (Aufruf wie make image-scan, ohne --exit-code, nur Schwachstellen; die Zusammenfassung nennt je Ziel die Anzahl, hier nur Ziele mit Anzahl > 0)"
docker pull -q "$trivy" >/dev/null
for sev in HIGH,CRITICAL CRITICAL; do
  echo "-- severity=$sev"
  docker run --rm -v /var/run/docker.sock:/var/run/docker.sock "$trivy" image \
    --scanners vuln --severity "$sev" --no-progress --report summary "$tag" 2>&1 \
    | grep -E "Detected OS|Target|│ +[1-9][0-9]* +│" | cut -c1-200
done

echo "######## Aufraeumen"
docker rmi "$tag" >/dev/null 2>&1 && echo "Mess-Image entfernt"
echo "fertig"
