# Option B: Release-Archiv mit gepinnter Version und SHA-256 je Architektur (Muster ADR-0023 ohne Signatur).
FROM mcr.microsoft.com/devcontainers/base@sha256:6ee2b4097c75d33e9a94d941577a3286e5d629434fba384688b06daeeaaeac99
ARG OLLAMA_VERSION=0.40.6
ARG OLLAMA_SHA256_AMD64=c164ab8b1b28e2a247c29719b8e404c423c5b83ae64c3ab4e59d50c7ff7fc5e0
ARG OLLAMA_SHA256_ARM64=0a279d9428ad17b06fdb7dddf3b059b2c148433585fa3b30676a102fc459bb9f
USER root
RUN set -eu; \
    apt-get update -qq; \
    apt-get install -y -qq --no-install-recommends zstd; \
    rm -rf /var/lib/apt/lists/*
RUN set -eu; \
    case "$(dpkg --print-architecture)" in \
      amd64) a=amd64; sha="$OLLAMA_SHA256_AMD64" ;; \
      arm64) a=arm64; sha="$OLLAMA_SHA256_ARM64" ;; \
      *) echo "architektur nicht unterstuetzt" >&2; exit 1 ;; \
    esac; \
    date -u +"download-start %T"; \
    curl -fsSL -o /tmp/ollama.tar.zst "https://github.com/ollama/ollama/releases/download/v${OLLAMA_VERSION}/ollama-linux-${a}.tar.zst"; \
    date -u +"download-ende %T"; \
    ls -l /tmp/ollama.tar.zst; \
    echo "${sha}  /tmp/ollama.tar.zst" | sha256sum -c -; \
    zstd -dc /tmp/ollama.tar.zst | tar -t | awk -F/ '{print $1"/"$2}' | sort | uniq -c | sort -rn | head -40; \
    tar --zstd -xf /tmp/ollama.tar.zst -C /usr/local --no-same-owner; \
    rm /tmp/ollama.tar.zst; \
    ollama --version
