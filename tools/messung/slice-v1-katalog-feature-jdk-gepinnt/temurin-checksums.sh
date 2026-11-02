#!/usr/bin/env bash
# Messskript M1 des Slices slice-v1-katalog-feature-jdk-gepinnt: die
# Vertriebs-Lage des gepinnten Temurin-Releases. Der Aufruf
#   temurin-checksums.sh <tag> <download|verifyx64>
# <download> liest die zwei Prüfsummen-Dateien (je Architektur) und die
# Asset-Liste; <verifyx64> prüft den Archiv-Stream gegen die gelesene
# Prüfsumme (Same-Origin-Gegenprobe, ~200 MB).
set -eu

tag="${1:?Aufruf: temurin-checksums.sh <tag> <download|verifyx64>}"
mode="${2:?Aufruf: temurin-checksums.sh <tag> <download|verifyx64>}"
base="https://github.com/adoptium/temurin21-binaries/releases/download/${tag}"

# Die Asset-Slug-Konvention des Schlüssels: 21.0.13+9 -> 21.0.13_9.
slug="$(printf '%s' "${tag#jdk-}" | tr '+' '_')"

case "$mode" in
download)
  for s in x64 aarch64; do
    printf '== %s\n' "$s"
    curl -fsSL -- "$base/OpenJDK21U-jdk_${s}_linux_hotspot_${slug}.tar.gz.sha256.txt"
    echo
    curl -s -o /dev/null -w 'http=%{http_code} redirects=%{num_redirects} final-host=%{url_effective} \n' \
      -L -- "$base/OpenJDK21U-jdk_${s}_linux_hotspot_${slug}.tar.gz.sha256.txt" \
      | sed 's|https://[^/]*|&|;s|?.*||'
  done
  curl -s -- "https://api.github.com/repos/adoptium/temurin21-binaries/releases/tags/${tag}" \
    | sed -n 's/.*"tag_name": *"\([^"]*\)".*/tag_name: \1/p;s/.*"published_at": *"\([^"]*\)".*/published_at: \1/p'
  ;;
verifyx64)
  want="$(curl -fsSL -- "$base/OpenJDK21U-jdk_x64_linux_hotspot_${slug}.tar.gz.sha256.txt" | awk '{print $1}')"
  got="$(curl -fsSL -- "$base/OpenJDK21U-jdk_x64_linux_hotspot_${slug}.tar.gz" | sha256sum | awk '{print $1}')"
  printf 'want=%s\ngot=%s\n' "$want" "$got"
  test "$want" = "$got"
  echo "verifyx64: OK"
  ;;
*)
  echo "unbekannter Modus: $mode (download|verifyx64)" >&2
  exit 2
  ;;
esac