#!/usr/bin/env bash
# tap-check.sh — prueft, ob das Tap-Token (Secret HOMEBREW_TAP_GITHUB_TOKEN) auf das
# Homebrew-Tap SCHREIBEN darf (slice-v2-homebrew-formula, ADR-0016). Ein Rerun des
# Jobs `tap` prueft das nicht, wenn die Formel schon im Tap liegt (nichts zu pushen).
#
# Umgebung: TAP_TOKEN (Pflicht), TAP_REPO (optional, Default pt9912/homebrew-u-boot).
# Exit 0 = Push erlaubt; 1 = Token gueltig, aber ohne Schreibrecht; 2 = Token ungueltig/leer.
set -euo pipefail

: "${TAP_TOKEN:?tap-check: TAP_TOKEN is required}"
tap="${TAP_REPO:-pt9912/homebrew-u-boot}"

body="$(mktemp)"
trap 'rm -f "$body"' EXIT
code="$(curl -sS -o "$body" -w '%{http_code}' \
	-H "Authorization: Bearer ${TAP_TOKEN}" -H "Accept: application/vnd.github+json" \
	"https://api.github.com/repos/${tap}")"

case "$code" in
200) ;;
401) echo "tap-check: token rejected (401) — invalid, expired or pasted with extra characters." >&2; exit 2 ;;
404) echo "tap-check: ${tap} not visible to this token (404) — wrong repository or token without access." >&2; exit 1 ;;
*) echo "tap-check: unexpected HTTP ${code} from the GitHub API." >&2; exit 2 ;;
esac

# Die API-Sicht allein reicht NICHT: bei klassischen PATs zeigt `permissions.push` die Rechte
# des BENUTZERS, nicht die Scopes des Tokens — ein Token ohne `repo`-Scope meldet dort
# `push: true` und scheitert beim echten Push mit 403. Deshalb der echte Test: Klon und
# `git push --dry-run` (verhandelt die Schreib-Verbindung, schreibt aber nichts).
work="$(mktemp -d)"
trap 'rm -f "$body"; rm -rf "$work"' EXIT
git clone --quiet --depth 1 "https://x-access-token:${TAP_TOKEN}@github.com/${tap}.git" "$work/tap"
branch="$(git -C "$work/tap" rev-parse --abbrev-ref HEAD)"
if out="$(git -C "$work/tap" push --dry-run origin "HEAD:refs/heads/${branch}" 2>&1)"; then
	echo "tap-check: token may push to ${tap} (git push --dry-run succeeded)."
else
	echo "tap-check: token is valid but git push to ${tap} is DENIED:" >&2
	printf '  %s\n' "${out//$'\n'/$'\n'  }" >&2
	echo "tap-check: classic PAT -> needs the 'repo' scope; fine-grained PAT -> add the repository and set Contents to Read and write." >&2
	exit 1
fi
