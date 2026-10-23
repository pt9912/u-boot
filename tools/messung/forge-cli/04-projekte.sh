#!/bin/bash
# u-boot bauen (`make build`, Runtime-Image) und die Messprojekte erzeugen (init, config set, generate devcontainer).
# Der Katalog kennt weder gh noch glab; ein Feature kommt ueber eine externe Quelle (`source:` + Freigabeliste), so wie ein Nutzer es heute
# einrichten wuerde. Varianten (alle ausser ghn: Sandbox-Profil, Egress aktiv, onUnavailable warn):
#   bas   ohne Forge-Werkzeug (Vergleich)        gh    offizielles Feature github-cli
#   ghx   Feature gh-cli (devcontainers-extra)   glabs Feature gitlab-cli von skriptfabrik
#   glabp Feature gitlab-cli von postfinance     glabz Feature gitlab-cli von SrzStephen
#   ghn   wie gh, OHNE Sandbox-Profil          basn  wie bas, OHNE Sandbox-Profil (Vergleich zu ghn)
#   kat   wie bas, dazu das Katalog-Feature git (Form des Schluessels eines Katalogeintrags)
#   ega   wie bas, dazu devcontainer.sandbox.egress.allow gitlab.com (nur sein egress-init.sh wird verwendet, siehe 07-laufzeit.sh)
# Aufruf: WORK=<messordner> ./04-projekte.sh          -> $OUT/rohdaten-04.txt, rohdaten-04-erzeugt-<variante>.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-04.txt

echo "== make build IMAGE=$UB (Repo-Stand: Arbeitsverzeichnis $(basename "$ROOT"), HEAD $(git -C "$ROOT" rev-parse --short HEAD))"
s=$(date +%s.%N)
make -C "$ROOT" build IMAGE="$UB" > "$WORK/make-build.log" 2>&1; rc=$?
e=$(date +%s.%N)
echo "make build exit=$rc dauer=$(sekunden "$s" "$e")s"
tail -3 "$WORK/make-build.log"
[ "$rc" = 0 ] || exit 1
docker image inspect "$UB:latest" --format 'image={{.Id}} size={{.Size}} arch={{.Architecture}}'
docker run --rm "$UB" --version
aussage "u-boot-version=$(docker run --rm "$UB" --version | tr -d '\r')"
echo "-- Arbeitsverzeichnis sauber? (git status --short, Anzahl Zeilen; Skripte der Messung selbst stehen darin)"
git -C "$ROOT" status --short | grep -v -E 'tools/messung/forge-cli|docs/messungen/slice-v1-forge-cli' | wc -l
echo "-- sha256 der drei Vorlagen, aus denen die Projekte entstehen"
for f in Dockerfile.tmpl devcontainer.json.tmpl egress-init.sh.tmpl; do
  docker run --rm -v "$ROOT/internal/hexagon/application/templates/devcontainer":/t:ro "$BASE" sha256sum "/t/$f"
done

projekt() { # <variante> <sandbox:ja|nein> [<oci-variante>] [<egress.allow-host>] [<katalog-feature>]
  local v=$1 sb=$2 f=${3:-} d ref; d=$(vdir "$v")
  rm -rf "$d"; mkdir -p "$d"
  echo; echo "== Projekt $v (sandbox=$sb, feature=${f:-keines})"
  if [ "$sb" = ja ]; then
    u "p-$v" init "$(vname "$v")" --devcontainer --sandbox --no-git --yes > /dev/null 2>&1 || { echo "init fehlgeschlagen"; return 1; }
    u "p-$v" config set devcontainer.sandbox.egress.enabled true > /dev/null 2>&1
    u "p-$v" config set devcontainer.sandbox.onUnavailable warn > /dev/null 2>&1
    [ -n "${4:-}" ] && u "p-$v" config set devcontainer.sandbox.egress.allow "$4" 2>&1 | grep -v "^time="
  else
    u "p-$v" init "$(vname "$v")" --devcontainer --no-git --yes > /dev/null 2>&1 || { echo "init fehlgeschlagen"; return 1; }
  fi
  [ -n "${5:-}" ] && u "p-$v" config set "devcontainer.features.$5.enabled" true 2>&1 | grep -v "^time="
  if [ -n "$f" ]; then
    ref=https://$(oci_ref "$f")
    local key=gh; case "$f" in glab*) key=glab ;; esac
    u "p-$v" config set devcontainer.featureSources.allow "$ref" 2>&1 | grep -v '^time=' | head -3
    u "p-$v" config set "devcontainer.features.$key.source" "$ref" 2>&1 | grep -v '^time=' | head -3
    u "p-$v" config set "devcontainer.features.$key.enabled" true 2>&1 | grep -v '^time=' | head -3
  fi
  u "p-$v" generate devcontainer --yes 2>&1 | grep -v '^time='
  echo "generate-exit=${PIPESTATUS[0]}"
  [ "$sb" = ja ] && u "p-$v" egress list | grep -E '^(Egress|[a-z0-9.-]+ +[a-z:-]+$)' | grep -v '^HOST'
  for x in $(cd "$d/.devcontainer" && ls); do
    echo "=== .devcontainer/$x ==="; cat "$d/.devcontainer/$x"
  done | redact > "$OUT/rohdaten-04-erzeugt-$v.txt"
  echo "-- features und remoteEnv in der erzeugten devcontainer.json"
  grep -n -E '"features"|ghcr.io|remoteEnv|GIT_TOKEN|_TOKEN' "$d/.devcontainer/devcontainer.json"
  echo "-- u-boot.yaml (Abschnitt devcontainer)"
  cat "$d/u-boot.yaml" | grep -n -A12 '^devcontainer:' | head -30
}
projekt bas   ja
projekt basn  nein
projekt gh    ja gh
projekt ghx   ja ghx
projekt glabs ja glabs
projekt glabp ja glabp
projekt glabz ja glabz
projekt ghn   nein gh
projekt ega   ja "" gitlab.com
projekt kat   ja "" "" git
