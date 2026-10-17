#!/usr/bin/env bash
# package-asset-outputs.sh — gibt die Job-Outputs des Package-Smoke (publish.yml) auf
# stdout aus: `deb_asset=<name>` und `rpm_asset=<name>`; Asset-Namen auf stdin.
# Beide Namen werden zuerst aufgeloest und geprueft; bei einem Fehler (kein oder mehr als
# ein Treffer, leerer Name) endet das Skript mit Exit != 0 und schreibt nichts.
#
# AUFRUF (Workflow): scripts/package-asset-outputs.sh <<<"$NAMES" >> "$GITHUB_OUTPUT"
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
names="$(cat)"

deb="$("$here/package-asset-name.sh" deb <<<"$names")"
rpm="$("$here/package-asset-name.sh" rpm <<<"$names")"
[ -n "$deb" ] && [ -n "$rpm" ] || { echo "package-asset-outputs: empty asset name" >&2; exit 1; }

printf 'deb_asset=%s\nrpm_asset=%s\n' "$deb" "$rpm"
