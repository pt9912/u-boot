#!/bin/bash
# u-boot bauen (`make build`, Runtime-Image) und die drei Messprojekte erzeugen (init, config set, generate devcontainer).
# Varianten (Ergebnis von `generate devcontainer` mit dem gebauten Binary; danach genau EINE Handänderung, ausgewiesen):
#   std  Standardprofil (ohne Sandbox), Features node und claude-code, kein Ollama im Container: Weg (a), K2
#   res  Sandbox-Profil, Egress an, Features ollama, node, claude-code: Weg (a) K3a und Weg (b)
#   alw  wie res plus egress.allow=host.docker.internal (ueber `config set`): K3b
# HANDVARIANTE: in `runArgs` jeder Variante steht zusaetzlich --add-host=host.docker.internal:host-gateway (die Auslieferung erzeugt
# keinen solchen Eintrag); das Skript setzt sie nach `generate` mit sed im Werkzeug-Container.
# Voraussetzung: 00-vorlauf.sh. Aufruf: WORK=<messordner> tools/messung/ollama-cloud-konto/02-projekte.sh -> $ROH/rohdaten-02.txt, rohdaten-02-erzeugt-<variante>.txt
set -u
. "$(dirname "$0")/lib.sh"
vorlauf_pruefen || exit 2
rohdaten rohdaten-02.txt

echo "== make build IMAGE=$UB (Repo-Stand: $(git -C "$ROOT" rev-parse --short HEAD))"
s=$(date +%s.%N)
make -C "$ROOT" build IMAGE="$UB" > "$WORK/make-build.log" 2>&1; rc=$?
echo "make build exit=$rc dauer=$(sekunden "$s" "$(date +%s.%N)")s"
tail -3 "$WORK/make-build.log"
[ "$rc" = 0 ] || exit 1
docker image inspect "$UB:latest" --format 'image={{.Id}} size={{.Size}} arch={{.Architecture}}'
aussage "u-boot-version=$(docker run --rm "$UB" --version | tr -d '\r')"

projekt() { # <variante> <sandbox:true|false> <allow-eintrag|-> <features...>
  local v=$1 sb=$2 al=$3 d f k; shift 3; d=$(vdir "$v")
  rm -rf "$d"; mkdir -p "$d"
  echo; echo "== Projekt $v (sandbox=$sb, allow=$al, features=$*)"
  if [ "$sb" = true ]; then
    u "$(basename "$d")" init "$(vname "$v")" --devcontainer --sandbox --no-git --yes > /dev/null 2>&1 || { echo "init fehlgeschlagen"; return 1; }
  else
    u "$(basename "$d")" init "$(vname "$v")" --devcontainer --no-git --yes > /dev/null 2>&1 || { echo "init fehlgeschlagen"; return 1; }
  fi
  for k in "$@"; do
    u "$(basename "$d")" config set "devcontainer.features.$k.enabled" true 2>&1 | grep -v '^time='
  done
  if [ "$sb" = true ]; then
    u "$(basename "$d")" config set devcontainer.sandbox.egress.enabled true > /dev/null 2>&1
    u "$(basename "$d")" config set devcontainer.sandbox.onUnavailable warn > /dev/null 2>&1
  fi
  [ "$al" != - ] && u "$(basename "$d")" config set devcontainer.sandbox.egress.allow "$al" > /dev/null 2>&1
  echo "-- generate devcontainer --yes (Ausgabe des Werkzeugs, Exit-Code folgt)"
  u "$(basename "$d")" generate devcontainer --yes 2>&1 | grep -v '^time='; echo "generate exit=${PIPESTATUS[0]}"
  # HANDVARIANTE: --add-host in runArgs
  local j=/w/$(basename "$d")/.devcontainer/devcontainer.json
  if grep -q '"runArgs"' "$d/.devcontainer/devcontainer.json"; then
    tool sed -i 's#"runArgs": \["--cap-add=NET_ADMIN"\],#"runArgs": ["--cap-add=NET_ADMIN", "--add-host=host.docker.internal:host-gateway"],#' "$j"
  else
    tool sed -i 's#^  "features": {#  "runArgs": ["--add-host=host.docker.internal:host-gateway"],\n  "features": {#' "$j"
  fi
  aussage "projekt-$v-handvariante-add-host=$(grep -c -- '--add-host=host.docker.internal:host-gateway' "$d/.devcontainer/devcontainer.json")"
  echo "-- erzeugte Dateien (Pfad, Groesse)"
  ( cd "$d" && find . -type f -not -path './.git/*' -printf '%p %s\n' | sort )
  for f in $(cd "$d/.devcontainer" 2>/dev/null && ls); do
    echo "=== .devcontainer/$f ==="; cat "$d/.devcontainer/$f"
  done | redact > "$ROH/rohdaten-02-erzeugt-$v.txt"
  aussage "projekt-$v-ollama-block=$(grep -c 'ARG OLLAMA_VERSION' "$d/.devcontainer/Dockerfile" 2>/dev/null)"
}
projekt std false - node claude-code
projekt res true  - ollama node claude-code
projekt alw true  host.docker.internal ollama node claude-code
echo
echo "-- Ollama-Pin im erzeugten Dockerfile (Variante res)"
aussage "ollama-pin=$(grep -m1 -o 'OLLAMA_VERSION=[0-9.]*' "$(vdir res)/.devcontainer/Dockerfile")"
echo "-- Egress-Liste je Variante (u-boot egress list)"
for v in res alw; do echo "# $v"; u "$(basename "$(vdir $v)")" egress list 2>&1 | grep -v '^time='; done
echo "-- Zugangs-Namen in den erzeugten devcontainer.json (remoteEnv)"
for v in std res alw; do echo "# $v: $(grep -o '\${localEnv:[A-Z_]*}' "$(vdir $v)/.devcontainer/devcontainer.json" | tr '\n' ' ')"; done
