#!/bin/bash
# Negativprobe zu Option B: ein falscher SHA-256-Pin muss den Build vor dem Entpacken abbrechen.
# Aufruf: 12-negativprobe-pin.sh > rohdaten-12.txt 2>&1
set -u
here=$(cd "$(dirname "$0")" && pwd)
docker build --no-cache --progress=plain -f "$here/Dockerfile.b" \
  --build-arg OLLAMA_SHA256_ARM64=0000000000000000000000000000000000000000000000000000000000000000 \
  -t u-boot-messung-ollama-b-negativ "$here" 2>&1 | grep -E 'download-ende|FAILED|WARNING|ERROR|exit code|sha256sum|did not match' | cut -c1-220
echo "build-exit=${PIPESTATUS[0]}"
docker image inspect u-boot-messung-ollama-b-negativ >/dev/null 2>&1 && echo "Image entstanden (unerwartet)" || echo "kein Image entstanden"
