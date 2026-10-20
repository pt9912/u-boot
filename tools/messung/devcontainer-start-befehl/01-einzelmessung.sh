#!/usr/bin/env bash
# Einzelmessung des ausgelieferten Befehls `u-boot devcontainer up|down`
# (Messfrage: trägt der Ablauf von S1 und S4 auf einem anderen Host als dem
# CI-Runner, insbesondere unter macOS mit Colima?).
#
# Aufruf, im Verzeichnis dieses Skripts:
#   bash 01-einzelmessung.sh <u-boot-Binary> <Ordner im geteilten Bereich> <Ordner außerhalb> \
#       > ../../../docs/messungen/slice-v1-cli-devcontainer-up-down/rohdaten-01-<host>.txt 2>&1
#
# <u-boot-Binary>: das ausgelieferte Binary des Hosts (kein selbst gebautes Testprogramm).
# <Ordner im geteilten Bereich>: ein vorhandener Ordner, den die VM sieht (unter Colima
#   unterhalb von $HOME). <Ordner außerhalb>: ein vorhandener Ordner, den die VM nicht
#   sieht (Probe der Grenze: der Befehl soll dort mit Exit-Code 12 enden, nicht mit einem
#   leeren Bind-Mount weiterlaufen).
#
# Voraussetzung: Bash, Docker (Client und Engine), das u-boot-Binary, Netzzugriff auf die
# Registries (npm, Basis-Images). Das Skript braucht weder Go noch Node noch die
# devcontainer-CLI. Es setzt GIT_TOKEN auf einen Platzhalter, nie auf einen Zugangsdaten-Wert,
# und gibt keine Umgebung aus. Je Projekt legt es ein Verzeichnis unter dem übergebenen Ordner an
# und entfernt Container, Volumes, Images und Verzeichnis am Ende wieder. Zeitangaben haben die
# Auflösung einer Sekunde (portabel zwischen Linux und macOS).

set -u

bin=${1:?usage: 01-einzelmessung.sh <u-boot-binary> <shared-dir> <outside-dir>}
shared=${2:?usage: 01-einzelmessung.sh <u-boot-binary> <shared-dir> <outside-dir>}
outside=${3:?usage: 01-einzelmessung.sh <u-boot-binary> <shared-dir> <outside-dir>}

export GIT_TOKEN="platzhalter-kein-zugangsdatum-0123456789abcdef"

section() { printf '\n######## %s\n' "$*"; }
now() { date +%s; }

# timed runs the command, prints its exit code and its duration in seconds, and
# returns the exit code.
timed() {
  local t0 rc
  t0=$(now)
  "$@"
  rc=$?
  echo "-> exit=$rc dauer=$(( $(now) - t0 )) s"
  return $rc
}

container_of() {
  docker ps -aq --no-trunc --filter "label=devcontainer.local_folder=$1" | head -n 1
}

cleanup_project() {
  local dir=$1 ids vols
  ids=$(docker ps -aq --filter "label=devcontainer.local_folder=$dir")
  vols=""
  for id in $ids; do
    vols="$vols $(docker inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}{{.Name}} {{end}}{{end}}' "$id")"
  done
  # shellcheck disable=SC2086
  [ -n "$ids" ] && docker rm -f -v $ids > /dev/null 2>&1
  # shellcheck disable=SC2086
  [ -n "${vols// /}" ] && docker volume rm -f $vols > /dev/null 2>&1
  local imgs
  imgs=$(docker images -q --filter "reference=vsc-$(basename "$dir" | tr '[:upper:]' '[:lower:]')-*")
  # shellcheck disable=SC2086
  [ -n "$imgs" ] && docker rmi -f $imgs > /dev/null 2>&1
  rm -rf "$dir"
}

# project_flow runs the flow of S1 and S4 in the current directory (a subshell of
# `project`), so a return ends only the flow.
project_flow() {
  local label=$1 dir=$2
  shift 2
  echo "Projektpfad: $dir"
  timed "$bin" init probe --no-git "$@"

  section "$label: up"
  local t0 rc
  t0=$(now)
  "$bin" devcontainer up --json > up1.json 2> up1.err
  rc=$?
  echo "-> exit=$rc dauer=$(( $(now) - t0 )) s"
  echo "Umschlag (Auszug):"; grep -o '"status":"[a-z]*"\|"exitCode":[0-9]*\|"containerId":"[0-9a-f]*"' up1.json | tr '\n' ' '; echo
  echo "stderr (letzte 5 Zeilen):"; tail -n 5 up1.err
  if [ "$rc" -ne 0 ]; then echo "up endete mit $rc; weitere Schritte entfallen"; return 0; fi

  local id
  id=$(container_of "$dir")
  echo "Container: ${id:0:12}"
  echo "Label devcontainer.local_folder: $(docker inspect --format '{{index .Config.Labels "devcontainer.local_folder"}}' "$id")"
  echo "Label devcontainer.config_file: $(docker inspect --format '{{index .Config.Labels "devcontainer.config_file"}}' "$id")"
  echo "Einhaengungen mit docker.sock: $(docker inspect --format '{{range .Mounts}}{{.Source}} {{end}}' "$id" | grep -c docker.sock)"
  echo "Benutzer im Container: $(docker exec "$id" id -un) (uid $(docker exec "$id" id -u))"

  # Marker im Workspace-Volume (nur das Sandbox-Profil hat eines; im Standardprofil ist der
  # Workspace der Bind-Mount des Projektordners, dort entfällt der Marker).
  local ws
  ws=$(docker inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}{{.Destination}} {{end}}{{end}}' "$id" | tr ' ' '\n' | grep '^/workspaces/' | head -n 1)
  if [ -n "$ws" ]; then
    docker exec -u root "$id" sh -c "echo marker > '$ws/u-boot-marker'" && echo "Marker geschrieben in $ws"
  else
    echo "kein Workspace-Volume: kein Marker"
  fi

  section "$label: down"
  timed "$bin" devcontainer down
  echo "laeuft nach down: $(docker inspect --format '{{.State.Running}}' "$id" 2>&1)"
  timed "$bin" devcontainer down

  section "$label: up (nach down)"
  timed "$bin" devcontainer up > /dev/null 2>&1
  local id2
  id2=$(container_of "$dir")
  [ "$id2" = "$id" ] && echo "gleicher Container: ja" || echo "gleicher Container: nein ($id2)"

  section "$label: up --recreate"
  timed "$bin" devcontainer up --recreate > /dev/null 2>&1
  local id3
  id3=$(container_of "$dir")
  [ "$id3" != "$id" ] && echo "neuer Container: ja" || echo "neuer Container: nein"
  echo "alter Container noch da: $(docker inspect --type container "$id" > /dev/null 2>&1 && echo ja || echo nein)"
  if [ -n "$ws" ]; then
    echo "Marker nach --recreate: $(docker exec "$id3" cat "$ws/u-boot-marker" 2>&1)"
  fi

  section "$label: down (Aufräumen folgt)"
  timed "$bin" devcontainer down
}

# project <label> <base> <init flags...>: the flow for one project below <base>;
# the directory and everything the flow created are removed afterwards.
project() {
  local label=$1 base=$2
  shift 2
  section "Projekt: $label (Ordner unter $base)"
  local dir
  dir=$(mktemp -d "$base/ubl-XXXXXX") || { echo "mktemp in $base failed"; return 1; }
  ( cd "$dir" && project_flow "$label" "$dir" "$@" )
  cleanup_project "$dir"
}

section "Messstand"
date -u +"Zeit (UTC): %Y-%m-%dT%H:%M:%SZ"
echo "Host: $(uname -srm)"
echo "u-boot: $("$bin" --version 2>&1 | head -n 1)"
echo "docker context: $(docker context show 2>&1)"
docker version --format 'Docker Client {{.Client.Version}} ({{.Client.Os}}/{{.Client.Arch}}), Server {{.Server.Version}} ({{.Server.Os}}/{{.Server.Arch}})' 2>&1
docker info --format 'Engine: {{.OperatingSystem}}, Architektur {{.Architecture}}, Cgroup {{.CgroupVersion}}' 2>&1
if command -v colima > /dev/null 2>&1; then
  echo "colima: $(colima version 2>&1 | head -n 2 | tr '\n' ' ')"
else
  echo "colima: nicht gefunden"
fi
echo "Ordner im geteilten Bereich: $shared"
echo "Ordner außerhalb: $outside"

project "Standardprofil, geteilter Bereich" "$shared" --devcontainer
project "Sandbox-Profil, geteilter Bereich" "$shared" --devcontainer --sandbox

section "Probe der Grenze: Ordner außerhalb des geteilten Bereichs"
out_dir=$(mktemp -d "$outside/ubl-XXXXXX") || { echo "mktemp in $outside failed"; exit 1; }
(
  cd "$out_dir" || exit 1
  echo "Projektpfad: $out_dir"
  "$bin" init probe --no-git --devcontainer > /dev/null 2>&1
  t0=$(now)
  "$bin" devcontainer up --json > up.json 2> up.err
  rc=$?
  echo "-> exit=$rc dauer=$(( $(now) - t0 )) s"
  echo "Umschlag (Auszug):"; grep -o '"status":"[a-z]*"\|"exitCode":[0-9]*' up.json | tr '\n' ' '; echo
  echo "stderr (letzte 5 Zeilen):"; tail -n 5 up.err
)
cleanup_project "$out_dir"

section "Ende"
echo "fertig"
