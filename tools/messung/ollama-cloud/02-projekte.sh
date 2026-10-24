#!/bin/bash
# u-boot bauen (`make build`, Runtime-Image) und die drei Messprojekte erzeugen (init, config set, generate devcontainer).
# Varianten (Features: ollama, node, claude-code; Ergebnis von `generate devcontainer` mit dem gebauten Binary, nichts von Hand):
#   std  Standardprofil (ohne Sandbox): C7
#   res  Sandbox-Profil, Egress an (Modus restricted: erzeugte Liste)
#   alw  wie res plus egress.allow=ollama.com (C2b, ueber `config set`)
# Aufruf: WORK=<messordner> tools/messung/ollama-cloud/02-projekte.sh     -> $OUT/rohdaten-02.txt, rohdaten-02-erzeugt-<variante>.txt, rohdaten-02-c7.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-02.txt

echo "== make build IMAGE=$UB (Repo-Stand: Arbeitsverzeichnis $(basename "$ROOT"))"
s=$(date +%s.%N)
make -C "$ROOT" build IMAGE="$UB" > "$WORK/make-build.log" 2>&1; rc=$?
echo "make build exit=$rc dauer=$(sekunden "$s" "$(date +%s.%N)")s"
tail -3 "$WORK/make-build.log"
[ "$rc" = 0 ] || exit 1
docker image inspect "$UB:latest" --format 'image={{.Id}} size={{.Size}} arch={{.Architecture}}'
aussage "u-boot-version=$(docker run --rm "$UB" --version | tr -d '\r')"

projekt() { # <variante> <sandbox:true|false> <egress:true|false> <allow-eintrag|->
  local v=$1 sb=$2 eg=$3 al=$4 d f; d=$(vdir "$v")
  rm -rf "$d"; mkdir -p "$d"
  echo; echo "== Projekt $v (sandbox=$sb, egress=$eg, allow=$al)"
  if [ "$sb" = true ]; then
    u "$(basename "$d")" init "$(vname "$v")" --devcontainer --sandbox --no-git --yes > /dev/null 2>&1 || { echo "init fehlgeschlagen"; return 1; }
  else
    u "$(basename "$d")" init "$(vname "$v")" --devcontainer --no-git --yes > /dev/null 2>&1 || { echo "init fehlgeschlagen"; return 1; }
  fi
  for k in ollama node claude-code; do
    u "$(basename "$d")" config set "devcontainer.features.$k.enabled" true 2>&1 | grep -v '^time='
  done
  [ "$sb" = true ] && u "$(basename "$d")" config set devcontainer.sandbox.egress.enabled "$eg" > /dev/null 2>&1
  [ "$sb" = true ] && u "$(basename "$d")" config set devcontainer.sandbox.onUnavailable warn > /dev/null 2>&1
  [ "$al" != - ] && u "$(basename "$d")" config set devcontainer.sandbox.egress.allow "$al" > /dev/null 2>&1
  echo "-- generate devcontainer --yes (Ausgabe des Werkzeugs, Exit-Code folgt)"
  u "$(basename "$d")" generate devcontainer --yes 2>&1 | grep -v '^time='; echo "generate exit=${PIPESTATUS[0]}"
  echo "-- erzeugte Dateien (Pfad, Groesse)"
  ( cd "$d" && find . -type f -not -path './.git/*' -printf '%p %s\n' | sort )
  for f in $(cd "$d/.devcontainer" 2>/dev/null && ls); do
    echo "=== .devcontainer/$f ==="; cat "$d/.devcontainer/$f"
  done | redact > "$OUT/rohdaten-02-erzeugt-$v.txt"
  aussage "projekt-$v-ollama-block=$(grep -c 'ARG OLLAMA_VERSION' "$d/.devcontainer/Dockerfile" 2>/dev/null)"
}
projekt std false false -
projekt res true  true  -
projekt alw true  true  ollama.com
echo
echo "-- Ollama-Pin im erzeugten Dockerfile (Variante res)"
grep -n 'OLLAMA_VERSION' "$(vdir res)/.devcontainer/Dockerfile" | head -3
aussage "ollama-pin=$(grep -m1 -o 'OLLAMA_VERSION=[0-9.]*' "$(vdir res)/.devcontainer/Dockerfile")"
echo "-- Egress-Liste je Variante (u-boot egress list)"
for v in res alw; do echo "# $v"; u "$(basename "$(vdir $v)")" egress list 2>&1 | grep -v '^time='; done
echo "-- C7: Standardprofil, Meldungen und Dateien (Ausgabe von generate in rohdaten-02-c7.txt)"
d=$(vdir std)
{ u "$(basename "$d")" generate devcontainer --yes 2>&1 | grep -v '^time=' ; echo "zweiter Lauf (idempotent) exit=${PIPESTATUS[0]}"; } > "$OUT/rohdaten-02-c7.txt"
cat "$OUT/rohdaten-02-c7.txt"
grep -n -i 'ollama' "$d/.devcontainer/Dockerfile" | head -3 || true
aussage "c7-std-ollama-im-dockerfile=$(grep -c -i 'ollama' "$d/.devcontainer/Dockerfile")"
aussage "c7-std-claude-code-feature-in-json=$(grep -c 'claude-code' "$d/.devcontainer/devcontainer.json")"
