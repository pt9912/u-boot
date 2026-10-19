#!/bin/bash
# u-boot bauen (`make build`, Runtime-Image) und die sechs Messprojekte erzeugen (init, config set, generate devcontainer).
# Varianten (alle: Sandbox-Profil, Egress aktiv, onUnavailable warn):
#   oll   ollama                      pod   ollama + nestedRuntime podman      dok   ollama + nestedRuntime docker
#   bas   ohne ollama                 podn  nestedRuntime podman ohne ollama    dokn  nestedRuntime docker ohne ollama
# Aufruf: WORK=<messordner> ./02-projekte.sh          -> $OUT/rohdaten-02.txt, rohdaten-02-erzeugt-<variante>.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-02.txt

echo "== make build IMAGE=$UB (Repo-Stand: Arbeitsverzeichnis $(basename "$ROOT"))"
s=$(date +%s.%N)
make -C "$ROOT" build IMAGE="$UB" > "$WORK/make-build.log" 2>&1; rc=$?
e=$(date +%s.%N)
echo "make build exit=$rc dauer=$(sekunden "$s" "$e")s"
tail -3 "$WORK/make-build.log"
[ "$rc" = 0 ] || exit 1
docker image inspect "$UB:latest" --format 'image={{.Id}} size={{.Size}} arch={{.Architecture}}'
docker run --rm "$UB" --version
echo "-- sha256 der beiden Vorlagen, aus denen die Projekte entstehen"
for f in Dockerfile.tmpl egress-init.sh.tmpl; do
  docker run --rm -v "$ROOT/internal/hexagon/application/templates/devcontainer":/t:ro "$BASE" sha256sum "/t/$f"
done
aussage "u-boot-version=$(docker run --rm "$UB" --version | tr -d '\r')"

projekt() { # <variante> <ollama:true|false> <nestedRuntime:none|podman|docker>
  local v=$1 oll=$2 rt=$3 d; d=$(vdir "$v")
  rm -rf "$d"; mkdir -p "$d"
  echo; echo "== Projekt $v (ollama=$oll, nestedRuntime=$rt)"
  u "p-$v" init "$(vname "$v")" --devcontainer --sandbox --no-git --yes > /dev/null 2>&1 || { echo "init fehlgeschlagen"; return 1; }
  u "p-$v" config set devcontainer.sandbox.egress.enabled true > /dev/null 2>&1
  u "p-$v" config set devcontainer.sandbox.onUnavailable warn > /dev/null 2>&1
  [ "$oll" = true ] && u "p-$v" config set devcontainer.features.ollama.enabled true > /dev/null 2>&1
  [ "$rt" != none ] && u "p-$v" config set devcontainer.sandbox.nestedRuntime "$rt" > /dev/null 2>&1
  u "p-$v" generate devcontainer --yes 2>&1 | grep -v '^time='
  u "p-$v" egress list | grep -E '^(Egress|[a-z0-9.-]+ +[a-z:-]+$)' | grep -v '^HOST'
  # Die erzeugten Dateien stehen als Text in rohdaten-02-erzeugt-<variante>.txt (Abschnitt '=== <Datei> ===')
  for f in $(cd "$d/.devcontainer" && ls); do
    echo "=== .devcontainer/$f ==="; cat "$d/.devcontainer/$f"
  done | redact > "$OUT/rohdaten-02-erzeugt-$v.txt"
  aussage "projekt-$v-ollama-block=$(grep -c 'ARG OLLAMA_VERSION' "$d/.devcontainer/Dockerfile")"
}
# Das u-boot-Image hat den Eintrittspunkt `u-boot`; die Projekte entstehen im Messordner.
projekt oll  true  none   || exit 1
projekt bas  false none   || exit 1
projekt pod  true  podman || exit 1
projekt podn false podman || exit 1
projekt dok  true  docker || exit 1
projekt dokn false docker || exit 1
echo
echo "-- Block 'Ollama' im erzeugten Dockerfile (Variante oll)"
grep -n -A28 'ARG OLLAMA_VERSION' "$(vdir oll)/.devcontainer/Dockerfile" | head -34
