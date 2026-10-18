#!/bin/bash
# O1: Probe-Image (Sandbox mit Egress + Ollama ohne CUDA) bauen; Größenvergleich zu Dockerfile.b (mit CUDA).
# Aufruf: 04-build-probe.sh <projekt> > rohdaten-04.txt 2>&1   (projekt = Ziel von 03-projekt-erzeugen.sh)
set -u
here=$(cd "$(dirname "$0")" && pwd)
proj=${1:?projekt}
echo "== Build Probe (arm64, --no-cache)"
time docker build --no-cache -q -f "$here/Dockerfile.probe" -t u-boot-messung-ollama-probe "$proj/.devcontainer"
for i in u-boot-messung-ollama-b u-boot-messung-ollama-probe; do
  docker image inspect $i --format "$i size={{.Size}}"
done
docker run --rm u-boot-messung-ollama-probe bash -c 'du -sh /usr/local/lib/ollama /usr/local/bin/ollama; ls /usr/local/lib/ollama | head -50; ls -l /usr/local/bin/ollama; dpkg -s zstd >/dev/null 2>&1 && echo "zstd noch da" || echo "zstd entfernt"'
