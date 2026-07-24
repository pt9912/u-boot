#!/usr/bin/env bash
# fetch-baseline-cache — materialisiert die lokale, netzlose Lese-Form des
# adoptierten Betriebsregelwerks (AI-Harness-Kurs) als committet-vendored
# Baseline (harness/conventions.md MR-004 „committet vendored", MR-007 Ortswahl
# `.harness/`):
#
#   regelwerk → .harness/baseline/<tag>/regelwerk/   (COMMITTET, vendored)
#   templates → .harness/baseline/<tag>/templates/   (COMMITTET, vendored)
#               + .harness/baseline/<tag>/SHA256SUMS  (Integritäts-/Provenienz-
#                 Manifest über BEIDE Bäume)
#
# u-boot vendored BEIDE Bäume (Upstream-Default, MR-004): das Regelwerk verweist
# mit ../templates/… auf die Templates als „Ziel-Form" — netzlos nur auflösbar,
# wenn templates/ parallel zu regelwerk/ vendored ist. Die Templates tragen zwei
# Rollen: Referenz-Form (Verweis-Ziel) und Kopiervorlage (neue Artefakte werden
# aus templates/ kopiert und ausgefüllt, nicht frei formuliert).
#
# Modi:
#   (default)  re-vendor: zieht das Release-Asset lab-regelwerk.zip, entpackt es
#              in den committeten Vendor-Pfad, (re)generiert SHA256SUMS und
#              verifiziert. Netz nötig (Release-Download) — Anlass: Baseline-Bump.
#   --verify   nur Integritätsprüfung des committeten Bestands gegen SHA256SUMS.
#              Offline, kein Netz — für CI/Audit/frischen Checkout.
#
# Tag-Quelle: ohne Argument die §Baseline-`**Stand:**`-Zeile in
# harness/conventions.md (Skript-Eingabe; der Pin ist nicht vollumfänglicher
# SSoT — ein Bump zieht zusätzlich Vendor-Pfad, AGENTS.md und harness/README.md
# nach, MR-004 Bump-Prozedur). Mit Argument ein expliziter Tag (z. B. `v3.5.1`);
# der erste Vendor-Lauf während der Erst-Adoption nutzt das explizite Argument,
# weil conventions.md im selben Slice erst entsteht (Bootstrap-Barriere).
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
root="$(pwd -P)"

repo="pt9912/ai-harness-course"
conventions="harness/conventions.md"

mode="vendor"
if [ "${1:-}" = "--verify" ]; then mode="verify"; shift; fi

tag="${1:-}"
if [ -z "$tag" ]; then
  tag="$(grep -m1 '\*\*Stand:\*\*' "$conventions" 2>/dev/null \
    | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | head -1 || true)"
fi
if ! [[ "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "fetch-baseline-cache: ungültiger/leerer Tag '${tag}' — Argument vMAJOR.MINOR.PATCH angeben oder §Baseline in ${conventions} prüfen" >&2
  exit 1
fi

baseline=".harness/baseline/${tag}"
sums="${baseline}/SHA256SUMS"

verify() {
  # Integritätsprüfung des committeten Bestands gegen SHA256SUMS (offline).
  for c in sha256sum find; do
    command -v "$c" >/dev/null 2>&1 \
      || { echo "fetch-baseline-cache: '$c' nicht gefunden (Host-Werkzeug)" >&2; exit 1; }
  done
  [ -f "$sums" ] \
    || { echo "fetch-baseline-cache: ${sums} fehlt — erst re-vendor (ohne --verify) laufen" >&2; exit 1; }
  echo "fetch-baseline-cache: verify ${baseline} gegen SHA256SUMS"
  ( cd "$baseline" && sha256sum -c SHA256SUMS )
  # Manifest-Deckung: die Datei-Anzahl auf Platte muss der Manifest-Zeilenzahl
  # gleichen. Das faengt Post-Vendor-Drift ab (nachträglich gelöschte/zusätzliche
  # Datei, unvollständiger Checkout, manuell editiertes Manifest) - `sha256sum -c`
  # allein würde eine untermengige, aber in sich konsistente SHA256SUMS grün
  # passieren. Der Leer-Fall scheitert ohnehin laut (exit 1). Die Under-Copy-
  # Barriere (Quelle vs. vendored) sitzt im re-vendor-Pfad, nicht hier.
  local on_disk manifest
  on_disk="$(find "${baseline}/regelwerk" "${baseline}/templates" -type f 2>/dev/null | wc -l)"
  manifest="$(grep -c . "$sums" || true)"
  [ "$on_disk" -gt 0 ] \
    || { echo "fetch-baseline-cache: 0 Dateien — leeres/kaputtes Vendoring" >&2; exit 1; }
  [ "$on_disk" = "$manifest" ] \
    || { echo "fetch-baseline-cache: Manifest (${manifest} Zeilen) != Dateien auf Platte (${on_disk}) — unvollständig" >&2; exit 1; }
  echo "fetch-baseline-cache: verify ok (${manifest} Dateien, vollständig)"
}

if [ "$mode" = "verify" ]; then
  verify
  exit 0
fi

# --- re-vendor (Netz) ---
for cmd in curl unzip sha256sum find; do
  command -v "$cmd" >/dev/null 2>&1 \
    || { echo "fetch-baseline-cache: '${cmd}' nicht gefunden (Host-Werkzeug)" >&2; exit 1; }
done

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

url="https://github.com/${repo}/releases/download/${tag}/lab-regelwerk.zip"
echo "fetch-baseline-cache: ${tag}/lab-regelwerk.zip -> ${baseline}/{regelwerk,templates}/"
curl -fsSL -o "${tmp}/lab-regelwerk.zip" "$url"
stage="${tmp}/stage"
mkdir -p "$stage"
unzip -oq "${tmp}/lab-regelwerk.zip" -d "$stage"

# Bundle-Layout tolerant auflösen: das Regelwerk-Verzeichnis ist das, das die
# Module trägt (modul-00-*.md), egal ob flach oder unter verschachteltem
# regelwerk/. Die Templates liegen als GESCHWISTER daneben (archive-root/templates).
modul0="$(find "$stage" -type f -name 'modul-00-*.md' | head -1)"
[ -n "$modul0" ] \
  || { echo "fetch-baseline-cache: kein Regelwerk (modul-00-*.md) im ZIP gefunden" >&2; exit 1; }
src_regelwerk="$(dirname "$modul0")"
archive_root="$(dirname "$src_regelwerk")"
src_templates="${archive_root}/templates"
[ -d "$src_templates" ] \
  || { echo "fetch-baseline-cache: kein templates/-Geschwisterbaum neben regelwerk/ gefunden (v3.5.1 bündelt beide; Schritt-0-Annahme verletzt)" >&2; exit 1; }

rm -rf "${baseline}/regelwerk" "${baseline}/templates"
mkdir -p "${baseline}/regelwerk" "${baseline}/templates"
# Regelwerk-Markdown (README + grundlagen-* + modul-*) REKURSIV, Struktur erhalten
# (rekursiv statt maxdepth-1, damit ein künftiger Bundle-Reorg mit
# Unterverzeichnissen nicht still unter-kopiert). Ziel absolut über $root, weil
# die Subshell nach $src_regelwerk wechselt.
( cd "$src_regelwerk" && find . -type f -name '*.md' -print0 | while IFS= read -r -d '' f; do
    dest="${root}/${baseline}/regelwerk/${f#./}"
    mkdir -p "$(dirname "$dest")"
    cp "$f" "$dest"
  done )
# Templates vollständig (rekursiv, Struktur erhalten) — Referenz-Form + Vorlage.
( cd "$src_templates" && find . -type f -print0 | while IFS= read -r -d '' f; do
    dest="${root}/${baseline}/templates/${f#./}"
    mkdir -p "$(dirname "$dest")"
    cp "$f" "$dest"
  done )

# Under-Copy-Barriere (die eigentliche Vollständigkeitsprüfung): der vendorte
# Bestand muss exakt so viele Dateien tragen wie die Quelle im ZIP. Ein reines
# Post-Copy-vs-Manifest-Vergleich (verify) könnte einen übergangenen Quell-Zweig
# nicht sehen, weil Manifest und Platte beide post-copy sind.
src_n="$(( $(find "$src_regelwerk" -type f -name '*.md' | wc -l) + $(find "$src_templates" -type f | wc -l) ))"
dst_n="$(find "${baseline}/regelwerk" "${baseline}/templates" -type f | wc -l)"
[ "$src_n" = "$dst_n" ] \
  || { echo "fetch-baseline-cache: Quelle (${src_n}) != vendored (${dst_n}) — Kopierschritt hat Dateien übergangen" >&2; exit 1; }

# Manifest über den TATSÄCHLICHEN Dateibestand BEIDER Bäume (find, nicht
# Top-Level-Glob), damit reorganisierte/verschachtelte Dateien nicht still
# übergangen werden.
( cd "$baseline" && find regelwerk templates -type f | LC_ALL=C sort | xargs sha256sum > SHA256SUMS )
verify

echo "fetch-baseline-cache: fertig — vendored ${baseline}/{regelwerk,templates} (+SHA256SUMS)"
