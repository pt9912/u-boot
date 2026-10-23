# HAND-Variante (kein Katalogeintrag, nicht von u-boot erzeugt), Form des Ollama-Blocks: gepinntes Release-Archiv von cli/cli.
# Version und Pruefsummen (aus gh_<version>_checksums.txt des Release, gelesen am 2026-10-23) stehen hier; der Bau bricht vor dem
# Entpacken ab, wenn die SHA-256 nicht passt. Das Projekt signiert die Archive nicht; es veroeffentlicht eine GitHub-Attestierung.
ARG GH_VERSION=2.103.0
ARG GH_SHA256_AMD64=fdd768eb70f6409f80fcfa9f3a1c0adf0fc98d827ab3ce0bed2ac22762c5d786
ARG GH_SHA256_ARM64=7f322df34e624d14f5ec0c4acb121fd6dfa14392fdbbd1066aa8e0a05ac3721e
RUN arch="$(dpkg --print-architecture)" \
 && case "$arch" in amd64) sha256="${GH_SHA256_AMD64}" ;; arm64) sha256="${GH_SHA256_ARM64}" ;; *) echo "gh has no archive for $arch" >&2; exit 1 ;; esac \
 && work="$(mktemp -d)" \
 && cd "$work" \
 && curl -fsSLO "https://github.com/cli/cli/releases/download/v${GH_VERSION}/gh_${GH_VERSION}_linux_${arch}.tar.gz" \
 && echo "${sha256}  gh_${GH_VERSION}_linux_${arch}.tar.gz" | sha256sum -c - \
 && tar -xzf "gh_${GH_VERSION}_linux_${arch}.tar.gz" -C /usr/local --strip-components=1 --no-same-owner "gh_${GH_VERSION}_linux_${arch}/bin/gh" \
 && cd / \
 && rm -rf "$work"
