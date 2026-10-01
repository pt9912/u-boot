#!/usr/bin/env bash
# tap-nachzug.sh — zieht die Formel eines Release-Tags ins Homebrew-Tap nach
# (slice-v2-homebrew-formula, ADR-0016). Laedt das Release-Asset `u-boot.rb` des
# Tags, klont das Tap, schreibt Formula/u-boot.rb und pusht — nur bei einer
# Aenderung (idempotent).
#
# Umgebung (nicht als Argumente, damit der Workflow-Text keines davon nennt):
#   TAP_TOKEN  PAT mit Schreibrecht auf das Tap (Secret HOMEBREW_TAP_GITHUB_TOKEN)
#   TAG        Release-Tag (z. B. v0.5.1)
#   GH_REPO    optional, Quell-Repo des Releases (Default pt9912/u-boot)
#   TAP_REPO   optional, Tap-Repo (Default pt9912/homebrew-u-boot)
set -euo pipefail

: "${TAP_TOKEN:?tap-nachzug: TAP_TOKEN is required}"
: "${TAG:?tap-nachzug: TAG is required}"
src="${GH_REPO:-pt9912/u-boot}"
tap="${TAP_REPO:-pt9912/homebrew-u-boot}"

work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT

curl -fsSL "https://github.com/${src}/releases/download/${TAG}/u-boot.rb" -o "$work/u-boot.rb"
grep -q "version \"${TAG#v}\"" "$work/u-boot.rb" || {
	echo "tap-nachzug: the u-boot.rb asset of ${TAG} does not carry version ${TAG#v}." >&2
	exit 1
}

git clone --depth 1 "https://x-access-token:${TAP_TOKEN}@github.com/${tap}.git" "$work/tap"
mkdir -p "$work/tap/Formula"
cp "$work/u-boot.rb" "$work/tap/Formula/u-boot.rb"
cd "$work/tap"
if git diff --quiet -- Formula/u-boot.rb && [ -z "$(git status --porcelain Formula/u-boot.rb)" ]; then
	echo "tap-nachzug: tap already at ${TAG}; nothing to do."
	exit 0
fi
git -c user.name="github-actions[bot]" -c user.email="41898282+github-actions[bot]@users.noreply.github.com" \
	add Formula/u-boot.rb
git -c user.name="github-actions[bot]" -c user.email="41898282+github-actions[bot]@users.noreply.github.com" \
	commit -m "u-boot ${TAG#v}"
git push origin HEAD
echo "tap-nachzug: ${tap} updated to ${TAG}."
