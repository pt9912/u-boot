#!/usr/bin/env bash
# messen-schwer.sh — schwere Messungen des Slice slice-v1-devcontainer-host-git-uebernehmen
# (LH-FA-DEV-009, LH-FA-DEV-010): remoteEnv mit ${localEnv:…} in Lebenszyklus-Befehlen der
# Dev-Container-CLI (Messung 2), Transport der Host-Variable in Container und Sitzung über
# `u-boot devcontainer up` und `shell` in beiden Profilen und mit Instanz (Messung 6) und
# die Stellen, an denen der Wert nach dem Start liegt (Messung 8).
#
# AUFRUF (aus der Repo-Wurzel; Docker-Client + lokale Engine, u-boot-Binary, beim ersten Lauf Netz):
#   bash tools/messung/devcontainer-host-git-uebernehmen/messen-schwer.sh <m2|m6|m8> <u-boot-binary> \
#     > docs/messungen/slice-v1-devcontainer-host-git-uebernehmen/rohdaten-<gruppe>-lauf<L>.txt 2>&1
#
# SICHERHEIT: Alle Werte der Messung sind Platzhalter (platzh-m8-…), nie echte Zugangsdaten.
# Lockerung (ausgewiesen): Der Hilfs-Container von `u-boot devcontainer up` bekommt den
# Docker-Socket des Hosts (wie runArgs); der M2-Aufruf setzt denselben Weg in einer reinen
# CLI-Probe. Alles, was das Skript anlegt, trägt die Präfixe mubmess-/ubmess-/vsc-ubmess-;
# aufgeräumt wird per trap EXIT/INT/TERM über Container-Labels und Volume-/Image-Namen dieser
# Präfixe, nie über ein Glob auf fremde Namen.
# Zeitlimits: je Dev-Container-Start 600 s (timeout -k 10); der erste Lauf eines Kontexts
# baut Images (Netz nötig, Minutenskala); die Wiederholungen laufen auf warmem Build-Cache.
set -u

GRUPPE="${1:?usage: messen-schwer.sh <m2|m6|m8> <u-boot-binary>}"
UB=$(readlink -f "${2:?usage: messen-schwer.sh <m2|m6|m8> <u-boot-binary>}")
WIEDERHOLUNGEN="${WIEDERHOLUNGEN:-2}"
LIMIT_UP=600
TMP="$(mktemp -d)"
export TMP
HELPER_IMG=ubmess-devcontainer-cli:0.89.0
TOKEN_PLACEH="platzh-m8-git-token-0123456789abcdef"
export HELPER_IMG TOKEN_PLACEH UB LIMIT_UP
DIRS=()

kopf() {
  echo "# Messung schwer ($GRUPPE) — $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "# Host: $(uname -srm)"
  echo "# Docker: $(docker version --format '{{.Client.Version}} / Server {{.Server.Version}} ({{.Server.Os}}/{{.Server.Arch}})')"
  echo "# u-boot: $(timeout 10 "$UB" --version)"
  echo "# Hilfs-Image: $HELPER_IMG"
  echo "# Token-Platzhalter (Länge): $(echo -n "$TOKEN_PLACEH" | wc -c)"
  echo "# Wiederholungen je Zelle: $WIEDERHOLUNGEN"
  echo
}

aufraeumen() {
  local dir ids vols imgs
  for dir in "${DIRS[@]}"; do
    ids="$(docker ps -aq --filter "label=devcontainer.local_folder=$dir" 2>/dev/null)"
    # shellcheck disable=SC2086
    [ -n "$ids" ] && docker rm -f -v $ids >/dev/null 2>&1
  done
  ids="$(docker ps -aq --filter name=u-boot-devcontainer-up- 2>/dev/null)"
  # shellcheck disable=SC2086
  [ -n "$ids" ] && docker rm -f $ids >/dev/null 2>&1
  vols="$(docker volume ls -q --filter name=mubmess 2>/dev/null)$(docker volume ls -q --filter name=ubmess 2>/dev/null)"
  # shellcheck disable=SC2086
  [ -n "$vols" ] && docker volume rm $vols >/dev/null 2>&1
  imgs="$(docker images --format '{{.Repository}}' | grep -E '^(vsc-)?mubmess|^vsc-ubmess|^ubmess-dev' 2>/dev/null)"
  # shellcheck disable=SC2086
  [ -n "$imgs" ] && docker rmi $imgs >/dev/null 2>&1
  rm -rf "$TMP"
}
trap aufraeumen EXIT INT TERM

# prepare_project <profil> -> echo dir : Wegwerf-Projektordner mit init/Config/generate
prepare_project() {
  local profil=$1 dir
  dir="$TMP/proj-$(echo "$profil" | tr / -)-$RANDOM"
  mkdir -p "$dir"
  ( cd "$dir" && "$UB" init --no-git "mubmess$RANDOM" >/dev/null 2>&1 \
    && "$UB" config set devcontainer.profile "$profil" >/dev/null 2>&1 \
    && "$UB" generate devcontainer --yes >/dev/null 2>&1 )
  DIRS+=("$dir")
  echo "$dir"
}

# cont_of <dir> -> letzter Container-Id des Projektordners
cont_of() { docker ps -aq --filter "label=devcontainer.local_folder=$1" | tail -1; }

# zelle <Kennung> <Beschreibung> <Befehl-String>
zelle() {
  local id="$1" beschr="$2" befehl="$3" n
  for n in $(seq 1 "$WIEDERHOLUNGEN"); do
    echo "### Zelle $id, Wiederholung $n: $beschr"
    echo "befehl: $befehl"
    timeout -k 10 "$LIMIT_UP" bash -c "$befehl" </dev/null >"$TMP/out" 2>"$TMP/err"
    echo "exit-status: $?"
    echo "--- stdout"
    cat "$TMP/out"
    echo "--- stderr"
    cat "$TMP/err"
    echo
  done
}

export -f prepare_project cont_of

case "$GRUPPE" in
m2)
  kopf
  if ! docker image inspect "$HELPER_IMG" >/dev/null 2>&1; then
    REPO="$(cd "$(dirname "$0")/../../.." && pwd)"
    echo "## Vorab: Bau des Hilfs-Images aus dem gepinnten Kontext $REPO/internal/adapter/driven/docker/devcontainercli"
    docker build -t "$HELPER_IMG" "$REPO/internal/adapter/driven/docker/devcontainercli" 2>&1 | tail -5
  fi
  echo "## M2: remoteEnv mit \${localEnv:…} in den Lebenszyklus-Befehlen der Dev-Container-CLI"
  echo "## (reine CLI-Probe, nicht über u-boot; minimaler Devcontainer mit postCreateCommand,"
  echo "##  Alpine-Basis, Workspace als Bind-Mount, Ergebnis-Datei im Wegwerf-Ordner)"
  # prepare_m2 <leer| gesetzt wird beim Aufruf über -e entschieden ->
prepare_m2() {
  local dir="$TMP/m2ws-$RANDOM"
  mkdir -p "$dir/.devcontainer"
  printf 'FROM alpine:3\n' > "$dir/.devcontainer/Dockerfile"
  printf '{\n  "name": "ubmess-m2",\n  "build": {"dockerfile": "Dockerfile", "context": "."},\n  "workspaceMount": "type=bind,source=%s,target=/workspaces/x",\n  "workspaceFolder": "/workspaces/x",\n  "remoteEnv": {"GIT_TOKEN": "${localEnv:GIT_TOKEN}"},\n  "containerUser": "root",\n  "postCreateCommand": "echo gt-laenge=$(printenv GIT_TOKEN 2>/dev/null | wc -c) > /workspaces/x/m2-ergebnis.txt"\n}\n' "$dir" > "$dir/.devcontainer/devcontainer.json"
  DIRS+=("$dir")
  echo "$dir"
}
export -f prepare_m2
  zelle M2-A "remoteEnv-Referenz, GIT_TOKEN im Host-Umfeld gesetzt (Platzhalter)" \
    "dir=\$(prepare_m2); docker run --rm -i --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock --mount type=bind,source=\"\$dir\",target=\"\$dir\" -e GIT_TOKEN=$TOKEN_PLACEH $HELPER_IMG devcontainer up --workspace-folder \"\$dir\" >\"\$dir/cli-out.txt\" 2>&1; echo \"cli-exit=\$?\"; cat \"\$dir/m2-ergebnis.txt\" 2>/dev/null || echo 'm2-ergebnis.txt fehlt'; tail -c 1500 \"\$dir/cli-out.txt\""
  zelle M2-B "remoteEnv-Referenz, GIT_TOKEN im Host-Umfeld ungesetzt" \
    "dir=\$(prepare_m2); docker run --rm -i --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock --mount type=bind,source=\"\$dir\",target=\"\$dir\" $HELPER_IMG devcontainer up --workspace-folder \"\$dir\" >\"\$dir/cli-out.txt\" 2>&1; echo \"cli-exit=\$?\"; cat \"\$dir/m2-ergebnis.txt\" 2>/dev/null || echo 'm2-ergebnis.txt fehlt'; tail -c 1500 \"\$dir/cli-out.txt\""
  ;;

m6)
  kopf
  echo "## M6: Transport der Host-Variable über u-boot up/shell (Platzhalter), beide Profile, mit Instanz"
  echo "## (Sitzungsumgebung über Eingaben der Sitzung gelesen; Container-Umgebung über docker exec)"
  sess_cmds='printenv GIT_TOKEN | wc -c; printf "protocol=https\nhost=host.example\n\n" | GIT_TERMINAL_PROMPT=0 git credential fill; echo fill-exit=$?; exit'
  zelle M6-A "Sandbox, GIT_TOKEN gesetzt, ohne Instanz, up; exec-Umgebung und fill-Probe im Container" \
    "dir=\$(prepare_project sandbox); export dir; GIT_TOKEN=$TOKEN_PLACEH timeout -k 10 $LIMIT_UP bash -c 'cd \"\$dir\" && '$UB' devcontainer up' >/dev/null 2>&1; echo \"up-exit=\$?\"; c=\$(cont_of \"\$dir\"); echo \"container=\$c\"; docker exec \"\$c\" sh -c 'printenv GIT_TOKEN | wc -c' 2>&1; docker exec \"\$c\" sh -c 'printf \"protocol=https\nhost=host.example\n\n\" | GIT_TERMINAL_PROMPT=0 git credential fill; echo exit=\$?' 2>&1"
  zelle M6-B "Sandbox, GIT_TOKEN gesetzt, ohne Instanz, shell mit Eingaben; Sitzungsumgebung und fill-Probe" \
    "dir=$(prepare_project sandbox); export dir; GIT_TOKEN=$TOKEN_PLACEH timeout -k 10 $LIMIT_UP bash -c 'cd \"\$dir\" && '$UB' devcontainer shell' <<'EOF'
$sess_cmds
EOF
"
  zelle M6-C "Sandbox, GIT_TOKEN ungesetzt, shell mit Eingaben (Sitzungsumgebung ohne Variable)" \
    "dir=$(prepare_project sandbox); export dir; env -u GIT_TOKEN timeout -k 10 $LIMIT_UP bash -c 'cd \"\$dir\" && '$UB' devcontainer shell' <<'EOF'
$sess_cmds
EOF
"
  zelle M6-D "Standardprofil, GIT_TOKEN gesetzt, shell mit Eingaben (kein remoteEnv im Standardprofil)" \
    "dir=$(prepare_project default); export dir; GIT_TOKEN=$TOKEN_PLACEH timeout -k 10 $LIMIT_UP bash -c 'cd \"\$dir\" && '$UB' devcontainer shell' <<'EOF'
$sess_cmds
EOF
"
  zelle M6-E "Sandbox, GIT_TOKEN gesetzt, --instance ubmess-i1, up und shell mit Eingaben" \
    "dir=\$(prepare_project sandbox); export dir; GIT_TOKEN=$TOKEN_PLACEH timeout -k 10 $LIMIT_UP bash -c 'cd \"\$dir\" && '$UB' devcontainer up --instance ubmess-i1' >/dev/null 2>&1; echo \"up-exit=\$?\"; c=\$(cont_of \"\$dir\"); echo \"container=\$c\"; docker inspect --format '{{.Config.Labels}}' \"\$c\" 2>&1 | head -c 400; echo; GIT_TOKEN=$TOKEN_PLACEH timeout -k 10 $LIMIT_UP bash -c 'cd \"\$dir\" && '$UB' devcontainer shell --instance ubmess-i1' <<'EOF'
$sess_cmds
EOF
"
  ;;

m8)
  kopf
  echo "## M8: Wo steht der Wert nach dem Start (Platzhalter): Prozesslisten während up, Container-Metadaten"
  echo "## (ein up-Lauf im Sandbox-Profil; Suche nach dem Platzhalter in Prozesslisten, Config.Env und"
  echo "##  Labels; Zeilen ohne Fund sind der Beleg des Ortes 'nicht dort')"
m8_run() {
  local dir c uppid i pid hc
  dir="$(prepare_project sandbox)"
  echo "Projektordner: $dir"
  GIT_TOKEN="$TOKEN_PLACEH" dir="$dir" bash -c 'cd "$dir" && "$UB" devcontainer up' &
  uppid=$!
  echo "--- Prozesslisten während des Starts (10 Proben à 4 s)"
  for i in $(seq 1 10); do
    sleep 4
    echo "== Probe $i (Host, Prozess-Argumentlisten mit Wort GIT_TOKEN):"
    ps -eo args | grep -F GIT_TOKEN | grep -v grep || echo "(keine)"
    echo "== Probe $i (Host, Platzhalter in Prozess-Argumentlisten):"
    ps -eo args | grep -F "$TOKEN_PLACEH" | grep -v grep || echo "(keine)"
    echo "== Probe $i (Host, Platzhalter in /proc/<pid>/environ der Hilfs-Client-Prozesse):"
    for pid in $(pgrep -f 'docker run.*(u-boot-devcontainer-up-|devcontainer)' 2>/dev/null); do
      { tr '\0' '\n' < "/proc/$pid/environ" 2>/dev/null | grep -F "$TOKEN_PLACEH" | head -2; } || true
    done
    echo "(Ende der environ-Suche)"
    echo "== Probe $i (Hilfs-Container, Platzhalter in Prozesslisten der Container):"
    for hc in $(docker ps --filter name=u-boot-devcontainer-up- --format '{{.Names}}' 2>/dev/null); do
      docker exec "$hc" sh -c "ps -ef 2>/dev/null | grep -F GIT_TOKEN | grep -v grep; for p in /proc/[0-9]*/environ; do grep -F '$TOKEN_PLACEH' \$p 2>/dev/null; done; true" || true
    done
    [ -z "$hc" ] && echo "(kein Hilfs-Container sichtbar)"
  done
  wait "$uppid"; echo "up-exit=$?"
  echo "--- nach dem Start: docker inspect des Devcontainers"
  c="$(cont_of "$dir")"; echo "container=$c"
  echo "== Config.Env:"
  docker inspect --format '{{json .Config.Env}}' "$c"
  echo "== Platzhalter in Config.Env: $(docker inspect --format '{{json .Config.Env}}' "$c" | grep -F "$TOKEN_PLACEH" | head -1 || echo nein)"
  echo "== Label devcontainer.metadata (erste 4000 Zeichen):"
  docker inspect --format '{{index .Labels "devcontainer.metadata"}}' "$c" | head -c 4000
  echo
  echo "== Platzhalter in devcontainer.metadata: $(docker inspect --format '{{index .Labels "devcontainer.metadata"}}' "$c" | grep -F "$TOKEN_PLACEH" | head -1 || echo nein)"
  echo "== localEnv in devcontainer.metadata: $(docker inspect --format '{{index .Labels "devcontainer.metadata"}}' "$c" | grep -F localEnv | head -c 400 || echo nein)"
  echo "== alle Labels:"
  docker inspect --format '{{json .Config.Labels}}' "$c"
}
export -f m8_run
  zelle M8-A "Probe während und nach up (Sandbox, Platzhalter)" "m8_run"
  ;;
*)
  echo "unbekannte Gruppe: $GRUPPE (erlaubt: m2, m6, m8)" >&2
  exit 2
  ;;
esac

echo
echo "# Ende $(date -u +%Y-%m-%dT%H:%M:%SZ)"