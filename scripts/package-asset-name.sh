#!/usr/bin/env bash
# package-asset-name.sh — waehlt aus den Namen der Release-Assets (stdin, eine Zeile je
# Name) das Linux-Paket fuer das Package-Smoke (publish.yml, ADR-0017).
#
# AUFRUF: gh release view <tag> --json assets -q '.assets[].name' | package-asset-name.sh <deb|rpm>
#
# Der Name wird gewaehlt, nicht aus der Version abgeleitet: nfpm schreibt eine
# Vorabversion als 0.10.0~rc.1 (Paket-Konvention), GitHub speichert das Asset als
# 0.10.0.rc.1; fuer stabile Tags sind beide Schreibweisen gleich. Genau ein Treffer
# fuer amd64, sonst Exit 1.
set -euo pipefail

format="${1:-}"
case "$format" in
	deb) pattern='^u-boot_.*_amd64\.deb$' ;;
	rpm) pattern='^u-boot-.*-1\.x86_64\.rpm$' ;;
	*) echo "package-asset-name: usage: package-asset-name.sh <deb|rpm> (asset names on stdin)" >&2; exit 2 ;;
esac

matches="$(grep -E "$pattern" || true)"
count="$(printf '%s' "$matches" | grep -c . || true)"
if [ "$count" -ne 1 ]; then
	echo "package-asset-name: expected exactly one $format asset for amd64, found $count" >&2
	exit 1
fi
printf '%s\n' "$matches"
