#!/bin/bash
# O1/O7: Option B bauen (kalt), Umleitungskette der Download-URL, Image-Größe.
# Aufruf: 02-build-b.sh > rohdaten-02.txt 2>&1
set -u
here=$(cd "$(dirname "$0")" && pwd)
echo "== Umleitungskette (Host, nur Header, Range 0-0)"
curl -sIL -r 0-0 https://github.com/ollama/ollama/releases/download/v0.40.6/ollama-linux-arm64.tar.zst \
  | grep -i -E '^(HTTP|location)' | sed 's/?.*//'
echo "== Build Option B (arm64, --no-cache)"
time docker build --no-cache --progress=plain -f "$here/Dockerfile.b" -t u-boot-messung-ollama-b "$here" 2>&1 | grep -v '^#[0-9]* \[internal\]' | tail -80
docker image inspect u-boot-messung-ollama-b --format 'image size={{.Size}} arch={{.Architecture}}'
echo "== Dateien"
docker run --rm u-boot-messung-ollama-b bash -c 'ls -l /usr/local/bin/ollama; du -sh /usr/local/lib/ollama /usr/local/bin/ollama; ls /usr/local/lib/ollama; ls -R /usr/local/lib/ollama | head -40; id vscode; command -v ollama; ollama --version 2>&1'
