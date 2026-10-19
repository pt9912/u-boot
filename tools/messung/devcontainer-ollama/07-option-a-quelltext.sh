#!/bin/bash
# O1 (Option A): den veroeffentlichten Feature-Quelltext (OCI-Artefakt 1.0.1) lesen: install.sh und Metadaten.
# Aufruf: 07-option-a-quelltext.sh > rohdaten-07.txt 2>&1   (Image u-boot-messung-devcontainer-cli aus Dockerfile.cli)
set -u
docker run --rm -i u-boot-messung-devcontainer-cli bash -s <<'SH'
set -eu
repo=prulloac/devcontainer-features/ollama
tok=$(curl -fsS "https://ghcr.io/token?scope=repository:$repo:pull" | sed -e 's/.*"token":"\([^"]*\)".*/\1/')
dig=sha256:daee60f59e843b210b3d89f420afff0d9713f92f8c6ebe40508b798eb85be5a3
mkdir -p /x && cd /x
curl -fsSL -H "Authorization: Bearer $tok" "https://ghcr.io/v2/$repo/blobs/$dig" -o f.tar
echo "blob sha256: $(sha256sum f.tar | cut -d' ' -f1)  (erwartet ${dig#sha256:})"
tar -xf f.tar
ls -la
echo "== devcontainer-feature.json"; cat devcontainer-feature.json
echo "== install.sh"; cat install.sh
SH
