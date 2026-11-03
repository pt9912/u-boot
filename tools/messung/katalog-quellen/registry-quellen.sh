#!/usr/bin/env bash
# Messskript M1 des Vorgangs "Katalogquellen cpp/postgres-client extern": die
# Registry-Lage der zwei Katalog-Quellen ghcr.io/devcontainers/features/{cpp,
# postgresql-client} samt der Steuerprobe typescript (die Quelle, die der
# Katalog-Kontext als existierende Referenz liest) — Tags-Listing (anonymous
# Bearer) und die Upstream-Quell-Bestandslesung (devcontainers/features/src).
# Aufruf: registry-quellen.sh > rohdaten-m1.txt
set -eu

tok() {
  curl -fsSL -- "https://ghcr.io/token?scope=repository:devcontainers/features/$1:pull" \
    | sed -n 's/.*"token": *"\([^"]*\)".*/\1/p'
}

tags() {
  T="$(tok "$1")"
  curl -sS -o /tmp/tags.json -w '%{http_code}' -H "Authorization: Bearer $T" \
    "https://ghcr.io/v2/devcontainers/features/$1/tags/list"
  cat /tmp/tags.json
}

upstream() {
  curl -sS -o /tmp/up.json -w '%{http_code}' \
    "https://raw.githubusercontent.com/devcontainers/features/main/src/$1/devcontainer-feature.json"
  echo
}

for x in cpp postgresql-client typescript; do
  echo "== ghcr.io/devcontainers/features/$x/tags/list (HTTP-Code):"
  tags "$x"
  echo
  echo "== Upstream devcontainers/features/src/$x/devcontainer-feature.json (HTTP-Code):"
  upstream "$x"
  echo
done

echo "== Bestandslesung src/features (GitHub API, Seiten 1-3, Treffer cpp/postgres):"
curl -fsSL 'https://api.github.com/repos/devcontainers/features/contents/src/features?per_page=100' \
  | grep -o '"name": *"[^"]*"' | grep -i -E 'cpp|postgres' || echo "(kein Treffer auf Seite 1)"
curl -fsSL 'https://api.github.com/repos/devcontainers/features/contents/src/features?per_page=100&page=2' \
  | grep -o '"name": *"[^"]*"' | grep -i -E 'cpp|postgres' || echo "(kein Treffer auf Seite 2)"
curl -fsSL 'https://api.github.com/repos/devcontainers/features/contents/src/features?per_page=100&page=3' \
  | grep -o '"name": *"[^"]*"' | grep -i -E 'cpp|postgres' || echo "(kein Treffer auf Seite 3)"