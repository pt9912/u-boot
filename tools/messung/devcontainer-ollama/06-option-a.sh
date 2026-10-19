#!/bin/bash
# O1 (Option A): Community-Feature ghcr.io/prulloac/devcontainer-features/ollama mit der Dev-Container-CLI messen.
# Aufruf: 06-option-a.sh <arbeitsverzeichnis> > rohdaten-06.txt 2>&1
# OUT ist das Ausgabeverzeichnis des Build-Logs rohdaten-06-build-a.log (Vorgabe: das Arbeitsverzeichnis des Aufrufs,
# im Protokoll docs/messungen/slice-v1-devcontainer-ollama-messung/).
# Voraussetzung: docker build -f Dockerfile.cli -t u-boot-messung-devcontainer-cli .   (Dev-Container-CLI + Docker-CLI)
# Die CLI laeuft in einem Wegwerf-Container mit dem Docker-Socket der Colima-VM.
set -u
w=${1:?arbeitsverzeichnis}
rm -rf "$w"; mkdir -p "$w/proj/.devcontainer"; w=$(cd "$w" && pwd)
cat > "$w/proj/.devcontainer/devcontainer.json" <<'JSON'
{
  "name": "ollama-a",
  "image": "mcr.microsoft.com/devcontainers/base:debian",
  "features": { "ghcr.io/prulloac/devcontainer-features/ollama:1": {} },
  "remoteUser": "vscode"
}
JSON
cat > "$w/inner-a.sh" <<'SH'
set -u
echo "devcontainer-cli $(devcontainer --version)"
echo "== Tags"
devcontainer features info tags ghcr.io/prulloac/devcontainer-features/ollama --output-format json 2>&1 | head -c 600; echo
echo "== Manifest / Metadaten"
devcontainer features info manifest ghcr.io/prulloac/devcontainer-features/ollama:1 --output-format json 2>&1 | head -c 3500; echo
echo "== Build (Debian 13, arm64)"
date -u +%T
time devcontainer build --workspace-folder /w/proj --image-name u-boot-messung-ollama-a --log-level debug 2>&1 | tail -150 > /w/build-a.log
date -u +%T
grep -n -i -E 'error|fail|curl|ollama|systemctl|nvidia|cuda|zstd' /w/build-a.log | cut -c1-200 | head -60
echo "== Ende Build; letzte Zeilen"
tail -5 /w/build-a.log | cut -c1-200
SH
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "$w":/w u-boot-messung-devcontainer-cli bash /w/inner-a.sh
echo "== Ergebnis-Image"
docker image inspect u-boot-messung-ollama-a --format 'size={{.Size}} arch={{.Architecture}}' 2>&1
docker image inspect u-boot-messung-ollama-a --format '{{index .Config.Labels "devcontainer.metadata"}}' 2>&1 | head -c 3000; echo
docker run --rm --entrypoint '' u-boot-messung-ollama-a bash -c 'command -v ollama && ls -l $(command -v ollama); ollama --version 2>&1; du -sh /usr/local/lib/ollama 2>&1; ls /usr/local/lib/ollama 2>&1 | head; ls -ld /usr/share/ollama /home/vscode/.ollama 2>&1; id ollama 2>&1; ls /etc/systemd/system/ollama.service 2>&1'
cp "$w/build-a.log" "${OUT:-.}/rohdaten-06-build-a.log" 2>/dev/null || true
