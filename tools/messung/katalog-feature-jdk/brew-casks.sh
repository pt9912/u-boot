#!/usr/bin/env bash
# Messskript M3 des Slices slice-v1-katalog-feature-jdk-gepinnt: die Brew-Form
# des Temurin-Vertriebs (die Auflage des Projektinhabers, 2026-11-02): die
# Casks `temurin` und `temurin@21` — Existenz, Form (macOS `pkg` desselben
# Releases) und Verhältnis zum Katalog-Weg der Archivform.
set -eu

for c in temurin temurin@21; do
  printf '== cask %s\n' "$c"
  curl -fsSL -- "https://formulae.brew.sh/api/cask/${c}.json" \
    | sed -n 's/.*"token": *"\([^"]*\)".*/token: \1/p;s/.*"tap": *"\([^"]*\)".*/tap: \1/p;s/.*"version": *"\([^"]*\)".*/version: \1/p;s/.*"homepage": *"\([^"]*\)".*/homepage: \1/p' \
    | head -8
  curl -fsSL -- "https://formulae.brew.sh/api/cask/${c}.json" \
    | sed 's|.*"url": *"\([^"]*\)".*"url_specs".*|url: \1|' \
    | head -2
  echo
done