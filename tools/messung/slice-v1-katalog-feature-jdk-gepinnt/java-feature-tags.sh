#!/usr/bin/env bash
# Messskript M2 des Slices slice-v1-katalog-feature-jdk-gepinnt: die OCI-Tags
# und die Metadaten des offiziellen `java`-Features (die Handbuch-Berichtigung
# der Plan-Frage "welche Tags trägt das java-OCI-Feature"; der Tag-Suffix des
# Render-Schlüssels ist das Feature-Release, nicht die Java-Version).
set -eu

T="$(curl -fsSL -- 'https://ghcr.io/token?scope=repository:devcontainers/features/java:pull' \
  | sed -n 's/.*"token": *"\([^"]*\)".*/\1/p')"
curl -fsSL -H "Authorization: Bearer $T" \
  'https://ghcr.io/v2/devcontainers/features/java/tags/list'

echo
echo "== Metadaten (main):"
curl -fsSL 'https://raw.githubusercontent.com/devcontainers/features/main/src/java/devcontainer-feature.json' \
  | sed -n 's/.*"id": *"\([^"]*\)".*/id: \1/p;s/.*"version": *"\([^"]*\)".*/feature version: \1/p;s/.*"name": *"\([^"]*\)".*/name: \1/p;s/.*"description": *"\([^"]*\)".*/description: \1/p' | head -6