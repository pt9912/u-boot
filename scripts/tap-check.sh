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

if grep -Eq '"push"[[:space:]]*:[[:space:]]*true' "$body"; then
	echo "tap-check: token may push to ${tap}."
else
	echo "tap-check: token is valid but has NO write permission on ${tap} (fine-grained PAT: add the repository and set Contents to Read and write)." >&2
	exit 1
fi
