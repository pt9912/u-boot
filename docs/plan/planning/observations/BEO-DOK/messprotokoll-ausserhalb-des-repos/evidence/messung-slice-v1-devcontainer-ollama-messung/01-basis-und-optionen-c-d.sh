#!/bin/bash
# O1 (Teil): Messstand, Basis-Image, Option D (Debian-Paket), Option C (Docker-Hub-Image, nur Metadaten).
# Aufruf: 01-basis-und-optionen-c-d.sh > rohdaten-01.txt 2>&1
set -u
BASE=mcr.microsoft.com/devcontainers/base:debian
echo "== Messstand"; date -u +%FT%TZ
docker version --format 'docker client {{.Client.Version}} server {{.Server.Version}} {{.Server.Os}}/{{.Server.Arch}}'
colima version | head -1
docker info --format 'VM: CPUs={{.NCPU}} MemTotal={{.MemTotal}} KernelVersion={{.KernelVersion}} OS={{.OperatingSystem}}'
docker pull -q "$BASE" >/dev/null
docker image inspect "$BASE" --format 'basis {{index .RepoDigests 0}} arch={{.Architecture}} size={{.Size}}'
echo "== Option D: Debian-Paket"
docker run --rm "$BASE" bash -c '
cat /etc/os-release | head -2
apt-get update -qq >/dev/null 2>&1
echo "-- policy"; apt-cache policy ollama
echo "-- search"; apt-cache search --names-only "^ollama"
echo "-- zstd"; apt-cache policy zstd | head -3
echo "-- im Basis-Image vorhanden:"
for t in curl zstd ca-certificates xz-utils; do dpkg -s $t >/dev/null 2>&1 && echo "$t ja" || echo "$t nein"; done'
echo "== Option C: ollama/ollama (nur Manifest, kein Pull der Schichten)"
docker buildx imagetools inspect ollama/ollama:0.40.6 2>&1 | head -30
