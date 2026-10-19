#!/bin/bash
# Erzeugt das Sandbox-Projekt (Egress aktiv) mit dem u-boot-Binary dieses Repos
# (`make build-binaries`, bin/u-boot-linux-arm64) in <zielverzeichnis>.
# Aufruf: 01-projekt-erzeugen.sh <repo-wurzel> <zielverzeichnis> [egress.allow-Liste]
# Voraussetzung: Docker (Colima); Binary unter <repo-wurzel>/bin/.
set -eu
root=$(cd "${1:?repo-wurzel}" && pwd)
dest=${2:?zielverzeichnis}
allow=${3:-}
rm -rf "$dest"; mkdir -p "$dest"
dest=$(cd "$dest" && pwd)
docker run --rm -v "$root/bin":/b:ro -v "$dest":/w -w /w -u "$(id -u):$(id -g)" -e HOME=/tmp -e ALLOW="$allow" debian:stable-slim bash -c '
set -eu
/b/u-boot-linux-arm64 init egressprobe --devcontainer --sandbox --no-git --yes
/b/u-boot-linux-arm64 config set devcontainer.sandbox.egress.enabled true
[ -z "$ALLOW" ] || /b/u-boot-linux-arm64 config set devcontainer.sandbox.egress.allow "$ALLOW"
/b/u-boot-linux-arm64 generate devcontainer --yes
ls -la .devcontainer
'
