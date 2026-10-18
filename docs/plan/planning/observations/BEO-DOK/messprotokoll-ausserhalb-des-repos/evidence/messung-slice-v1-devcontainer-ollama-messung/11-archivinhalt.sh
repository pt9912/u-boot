#!/bin/bash
# O1 (Groesse) fuer beide Architekturen: entpackte Groesse je Verzeichnis, aus dem Strom gelesen (nichts wird entpackt oder gespeichert).
# Aufruf: 11-archivinhalt.sh > rohdaten-11.txt 2>&1   (Image u-boot-messung-ollama-b aus 02-build-b.sh liefert zstd)
set -u
for a in arm64 amd64; do
  echo "== ollama-linux-$a.tar.zst: Bytes je Verzeichnis (zwei Ebenen unter lib/ollama), entpackt"
  docker run --rm --platform linux/arm64 --entrypoint bash u-boot-messung-ollama-b -c "
    curl -fsSL https://github.com/ollama/ollama/releases/download/v0.40.6/ollama-linux-$a.tar.zst | zstd -dc | tar -tv --quote-chars= 2>/dev/null |
    awk '{ n=\$6; s=\$3; split(n, p, \"/\"); k=p[1]\"/\"p[2]\"/\"((p[3]!=\"\" && p[2]==\"ollama\" && p[3] ~ /^(cuda|rocm|vulkan|mlx)/) ? p[3] : \"(sonst)\"); sum[k]+=s } END { for (k in sum) printf \"%12d %s\\n\", sum[k], k }' | sort -rn"
done
