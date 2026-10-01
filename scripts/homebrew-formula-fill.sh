#!/usr/bin/env bash
# homebrew-formula-fill.sh — befuellt das Homebrew-Formel-Skeleton mit den Werten
# eines Release-Tags (slice-v2-homebrew-formula, ADR-0016).
#
# FAIL-CLOSED VOR DEM SCHREIBEN: fehlt einer der vier brew-relevanten Plattform-
# Digests (darwin/linux x amd64/arm64 — Homebrew traegt kein Windows) in der
# SHA256SUMS, bricht der Lauf, statt eine Formel mit leerem Digest zu schreiben.
#
# AUFRUF: homebrew-formula-fill.sh <tag> <sums-datei> <skeleton> <ziel>
#   tag         Release-Tag (z. B. v0.5.1) -> __TAG__; ohne fuehrendes "v" -> __VERSION__
#   sums-datei  SHA256SUMS des Tags (Format: sha256sum, "<hash>  <name>")
#   skeleton    scripts/homebrew-formula.rb.tmpl
#   ziel        Pfad der befuellten Formel
set -euo pipefail

tag="${1:-}"
sums="${2:-}"
skeleton="${3:-}"
ziel="${4:-}"
[ -n "$tag" ] && [ -n "$sums" ] && [ -n "$skeleton" ] && [ -n "$ziel" ] || {
	echo "homebrew-formula-fill: usage: homebrew-formula-fill.sh <tag> <sums-datei> <skeleton> <ziel>" >&2
	exit 2
}
[ -f "$sums" ] || { echo "homebrew-formula-fill: $sums does not exist — no digests without SHA256SUMS." >&2; exit 2; }
[ -f "$skeleton" ] || { echo "homebrew-formula-fill: $skeleton does not exist." >&2; exit 2; }

version="${tag#v}"
digest() { awk -v n="u-boot-$1" '$2 == n {print $1}' "$sums"; }

for plat in darwin-amd64 darwin-arm64 linux-amd64 linux-arm64; do
	d="$(digest "$plat")"
	[[ "$d" =~ ^[0-9a-f]{64}$ ]] || {
		echo "homebrew-formula-fill: no valid SHA256SUMS entry for u-boot-$plat — aborting before writing." >&2
		exit 1
	}
done

sed \
	-e "s/__VERSION__/${version}/g" \
	-e "s/__TAG__/${tag}/g" \
	-e "s/__SHA256_DARWIN_AMD64__/$(digest darwin-amd64)/g" \
	-e "s/__SHA256_DARWIN_ARM64__/$(digest darwin-arm64)/g" \
	-e "s/__SHA256_LINUX_AMD64__/$(digest linux-amd64)/g" \
	-e "s/__SHA256_LINUX_ARM64__/$(digest linux-arm64)/g" \
	"$skeleton" >"$ziel"
echo "homebrew-formula-fill: $ziel written (tag $tag, four platform digests filled)."
