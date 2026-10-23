# HAND-Variante (kein Katalogeintrag, nicht von u-boot erzeugt), Form des Ollama-Blocks: gepinntes Release-Archiv von gitlab-org/cli.
# Version und Pruefsummen (aus checksums.txt des Release, gelesen am 2026-10-23) stehen hier; der Bau bricht vor dem Entpacken ab, wenn die
# SHA-256 nicht passt. Das Projekt veroeffentlicht keine Signatur und keine Attestierung fuer die Archive (rohdaten-03.txt).
ARG GLAB_VERSION=1.124.0
ARG GLAB_SHA256_AMD64=63df684842b0697e418649dd64374d37ac1da0fbcdcb422d6a5557222f8bb40f
ARG GLAB_SHA256_ARM64=0716395a0efeb991fc844b33fe55310645afcbd9b9746e50c725e9d1cd1b8d16
RUN arch="$(dpkg --print-architecture)" \
 && case "$arch" in amd64) sha256="${GLAB_SHA256_AMD64}" ;; arm64) sha256="${GLAB_SHA256_ARM64}" ;; *) echo "glab has no archive for $arch" >&2; exit 1 ;; esac \
 && work="$(mktemp -d)" \
 && cd "$work" \
 && curl -fsSLO "https://gitlab.com/gitlab-org/cli/-/releases/v${GLAB_VERSION}/downloads/glab_${GLAB_VERSION}_linux_${arch}.tar.gz" \
 && echo "${sha256}  glab_${GLAB_VERSION}_linux_${arch}.tar.gz" | sha256sum -c - \
 && tar -xzf "glab_${GLAB_VERSION}_linux_${arch}.tar.gz" -C /usr/local --no-same-owner bin/glab \
 && cd / \
 && rm -rf "$work"
