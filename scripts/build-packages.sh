#!/usr/bin/env bash
# build-packages.sh — baut die Linux-Pakete (.deb und .rpm, amd64 und arm64) aus den
# Release-Binaries (slice-v2-distro-pakete, ADR-0017). Docker-only: nfpm laeuft im
# gepinnten Image, auf dem Host genuegt Docker.
#
# AUFRUF: build-packages.sh <version> [bin-dir] [out-dir]
#   version  z. B. 0.6.0 (ohne "v"); Vorabversionen (0.7.0-rc.1) wandelt nfpm in die
#            Paket-Konvention (0.7.0~rc.1)
#   bin-dir  Verzeichnis mit u-boot-linux-amd64 / u-boot-linux-arm64 (Default: bin)
#   out-dir  Zielverzeichnis der Pakete (Default: packages)
set -euo pipefail

version="${1:-}"
bin_dir="${2:-bin}"
out_dir="${3:-packages}"
[ -n "$version" ] || { echo "build-packages: usage: build-packages.sh <version> [bin-dir] [out-dir]" >&2; exit 2; }
version="${version#v}"

# NFPM_IMAGE ist gepinnt (Supply-Chain, wie die uebrigen Tool-Images).
nfpm_image="${NFPM_IMAGE:-goreleaser/nfpm:v2.47.0}"

rm -rf "$out_dir"
mkdir -p "$out_dir"

for arch in amd64 arm64; do
	binary="${bin_dir}/u-boot-linux-${arch}"
	[ -f "$binary" ] || { echo "build-packages: $binary does not exist — run 'make build-binaries' first." >&2; exit 1; }
	for format in deb rpm; do
		docker run --rm --user "$(id -u):$(id -g)" -v "$PWD":/work -w /work \
			-e NFPM_ARCH="$arch" -e NFPM_VERSION="$version" -e NFPM_BINARY="$binary" \
			"$nfpm_image" package -f packaging/nfpm.yaml -p "$format" -t "$out_dir/"
	done
done
echo "build-packages: $(find "$out_dir" -type f | wc -l) packages in $out_dir/"
ls -l "$out_dir"
