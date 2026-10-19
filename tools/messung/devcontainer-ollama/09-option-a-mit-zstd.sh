#!/bin/bash
# O1/O6 (Option A, Zusatz): Das Community-Feature auf einem Basis-Image mit zstd (ohne zstd bricht der Build ab, siehe rohdaten-06).
# Aufruf: 09-option-a-mit-zstd.sh <arbeitsverzeichnis> > rohdaten-09.txt 2>&1   (Image u-boot-messung-devcontainer-cli noetig)
set -u
w=${1:?arbeitsverzeichnis}
rm -rf "$w"; mkdir -p "$w/proj/.devcontainer"; w=$(cd "$w" && pwd)
cat > "$w/proj/.devcontainer/devcontainer.json" <<'JSON'
{
  "name": "ollama-a-zstd",
  "image": "u-boot-messung-base-zstd",
  "features": { "ghcr.io/prulloac/devcontainer-features/ollama:1": {} },
  "remoteUser": "vscode"
}
JSON
docker build -q -t u-boot-messung-base-zstd - <<'DF'
FROM mcr.microsoft.com/devcontainers/base@sha256:6ee2b4097c75d33e9a94d941577a3286e5d629434fba384688b06daeeaaeac99
RUN apt-get update -qq && apt-get install -y -qq --no-install-recommends zstd && rm -rf /var/lib/apt/lists/*
DF
docker run --rm -i -v /var/run/docker.sock:/var/run/docker.sock -v "$w":"$w" u-boot-messung-devcontainer-cli bash -s "$w" <<'SH'
set -u
w=$1
cd "$w"
date -u +'build-start %T'
devcontainer build --workspace-folder "$w/proj" --image-name u-boot-messung-ollama-a --log-level debug > build-a.log 2>&1
echo "build exit=$?"; date -u +'build-ende %T'
grep -n -i -E 'error|ERROR|nvidia|cuda|systemd|systemctl|Install complete|Linux|Downloading|Adding|WARNING|useradd|ollama user|ollama.service|GPU' build-a.log | cut -c1-200 | head -50
tail -3 build-a.log | cut -c1-200
echo "== devcontainer up (Feature-Entrypoint: ollama serve)"
devcontainer up --workspace-folder "$w/proj" --id-label u-boot-messung=a --log-level info 2>&1 | tail -3 | cut -c1-300
cid=$(docker ps -q --filter label=u-boot-messung=a)
echo "container=$cid"
sleep 4
docker exec "$cid" bash -c 'echo "-- Prozesse"; ps -eo pid,user,args | grep -E "ollama|sleep|init" | grep -v grep | cut -c1-160; echo "-- Sockets"; ss -ltnp 2>/dev/null | grep -E "State|11434"; echo "-- Umgebung"; env | grep -i ollama; echo "-- Dateien"; ls -l /usr/local/bin/ollama; ls -ld /usr/share/ollama /home/vscode/.ollama /usr/local/lib/ollama 2>&1; id ollama 2>&1; du -sh /usr/local/lib/ollama; ls /usr/local/lib/ollama | head -20; echo "-- Wer startet serve"; ps -o pid,user,ppid,args -C ollama'
echo "-- Image-Metadaten (devcontainer.metadata)"
docker image inspect u-boot-messung-ollama-a --format '{{index .Config.Labels "devcontainer.metadata"}}' | head -c 2500; echo
docker image inspect u-boot-messung-ollama-a --format 'size={{.Size}}'
docker rm -f "$cid" >/dev/null
SH
