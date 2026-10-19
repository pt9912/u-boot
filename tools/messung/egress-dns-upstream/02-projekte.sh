#!/bin/bash
# u-boot bauen (`make build`, Runtime-Image), zwei Messprojekte erzeugen (init, config set, generate devcontainer) und die Images bauen.
# Projekte (beide: Sandbox-Profil, Egress aktiv, onUnavailable warn): sbx ohne nestedRuntime, pod mit nestedRuntime podman.
# Aufruf: WORK=<messordner> OUT=<ausgabeverzeichnis> 02-projekte.sh -> $OUT/rohdaten-02.txt, rohdaten-02-erzeugt-egress-init.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-02.txt

echo "== make build IMAGE=$UB (Repo-Stand: Arbeitsverzeichnis $(basename "$ROOT"), Commit $(git -C "$ROOT" rev-parse --short HEAD))"
make -C "$ROOT" build IMAGE="$UB" > "$WORK/make-build.log" 2>&1; rc=$?
echo "make build exit=$rc"; tail -2 "$WORK/make-build.log" | cut -c1-200
[ "$rc" = 0 ] || exit 1
docker run --rm "$UB" --version
echo "-- sha256 der Vorlage, aus der das Startscript entsteht"
docker run --rm -v "$ROOT/internal/hexagon/application/templates/devcontainer":/t:ro "$BASE" sha256sum /t/egress-init.sh.tmpl
aussage "u-boot-version=$(docker run --rm "$UB" --version | tr -d '\r')"

projekt() { # <variante> <nestedRuntime:none|podman>
  local v=$1 rt=$2 d; d=$(vdir "$v")
  rm -rf "$d"; mkdir -p "$d"
  echo; echo "== Projekt $v (nestedRuntime=$rt)"
  u "p-$v" init "$PFX-$v" --devcontainer --sandbox --no-git --yes > /dev/null 2>&1 || { echo "init fehlgeschlagen"; return 1; }
  u "p-$v" config set devcontainer.sandbox.egress.enabled true > /dev/null 2>&1
  u "p-$v" config set devcontainer.sandbox.onUnavailable warn > /dev/null 2>&1
  [ "$rt" != none ] && u "p-$v" config set devcontainer.sandbox.nestedRuntime "$rt" > /dev/null 2>&1
  u "p-$v" generate devcontainer --yes 2>&1 | grep -v '^time='
  u "p-$v" egress list | grep -v '^HOST' | head -40
  echo "-- runArgs laut devcontainer.json ($v)"
  grep -F -h '"runArgs"' "$d/.devcontainer/devcontainer.json"
  echo "-- postStartCommand, user, mounts ($v)"
  grep -E '"(postStartCommand|remoteUser|containerUser|mounts)"' "$d/.devcontainer/devcontainer.json" | cut -c1-300
  aussage "projekt-$v-runargs=$(grep -F -h '"runArgs"' "$d/.devcontainer/devcontainer.json" | tr -d ' ')"
}
projekt sbx none   || exit 1
projekt pod podman || exit 1
cp "$(vdir sbx)/.devcontainer/egress-init.sh" "$WORK/egress-init.sh"
{ echo "=== .devcontainer/egress-init.sh (Projekt sbx; die Datei des Projekts pod ist gleich: $(cmp -s "$(vdir sbx)/.devcontainer/egress-init.sh" "$(vdir pod)/.devcontainer/egress-init.sh" && echo ja || echo nein)) ==="; cat "$WORK/egress-init.sh"; } | redact > "$OUT/rohdaten-02-erzeugt-egress-init.txt"

echo; echo "== docker build der Images"
for v in sbx pod; do
  s=$(date +%s)
  docker build -q -t "$PFX-img-$v" "$(vdir $v)/.devcontainer" > /dev/null 2> "$WORK/build-$v.err"; echo "build $v exit=$? dauer=$(( $(date +%s) - s ))s"
done
docker build -q -f "$here/Dockerfile.messung" -t "$IMG_TOOL" "$here" > /dev/null 2> "$WORK/build-tool.err"; echo "build tool exit=$?"
echo "-- Versionen im Werkzeug-Image"
docker run --rm --entrypoint sh "$IMG_TOOL" -c 'dnsmasq --version | head -1; nft --version; sudo --version | head -1; dig -v 2>&1; head -2 /etc/os-release; id vscode; getent passwd nobody dnsmasq'
echo "-- podman im Image pod"
docker run --rm --entrypoint sh "$IMG_POD" -c '/usr/local/bin/podman --version; id vscode'
