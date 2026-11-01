#!/usr/bin/env bash
# messen.sh — Messung des Slices slice-v1-devcontainer-container-name (LH-FA-DEV-010).
#
# FRAGEN (§6 des Slices, Messungen 1 bis 5):
#  M1 (Modus cli)        Nimmt die gepinnte Dev-Container-CLI einen Containernamen an — Option an
#                        `devcontainer up` oder Eigenschaft `containerName` der devcontainer.json —
#                        und wie heißt der Container danach?
#  M6 (Modus config)     Trägt `--override-config` eine je Aufruf gerenderte Konfiguration mit
#                        Namen in runArgs an die CLI, und bleibt mit gesetztem Namen die Identität
#                        (devcontainerId, Volumes, Labels) gleich? Ersetzt S1-Option der CLI.
#  M2 (Modus var)        Löst die CLI `${devcontainerId}` oder eine andere Variable in `runArgs`
#                        der devcontainer.json auf?
#  M3 (Modus identitaet) Bleiben devcontainerId, Labels und Volume-Namen mit gesetztem Namen gleich
#                        zum Lauf ohne Namen (beide Profile, mit und ohne Instanz)?
#  M4 (Modus engine)     Namensregeln der Engine: erlaubte Zeichen, Länge, Konfliktverhalten; was
#                        die Engine für --name zurückweist (Fehlerwortlaut für die Kollisionsregel).
#  M5 (Modus bestand)    Verhalten bei Kollision: laufender/fremder Container gleiches Namens,
#                        gestoppter Container, `--remove-existing-container` daneben.
#
# AUFRUF (aus der Repo-Wurzel; Docker-Client + lokale Engine, das u-boot-Binary, Netz):
#   bash tools/messung/devcontainer-container-name/messen.sh <cli|var|identitaet|engine|bestand> \
#     bin/u-boot-linux-amd64 > docs/messungen/slice-v1-devcontainer-container-name/rohdaten-<modus>-lauf<N>.txt 2>&1
#
# Der Aufruf der Dev-Container-CLI spiegelt `runArgs` in
# internal/adapter/driven/docker/devcontainer_launcher.go (`docker run --rm --name <Hilfs-Container>
# --mount Socket --mount Projektordner <Image> devcontainer up --workspace-folder <Ordner> …`); das
# Skript fährt NICHT `u-boot devcontainer up`, sondern die CLI je Frage direkt.
#
# Lockerung (ausgewiesen): der Hilfs-Container bekommt den Docker-Socket des Hosts. Alles, was das
# Skript anlegt, trägt das Präfix ubmess-/vsc-ubmess-; nie ein Glob über fremde Namen. Vorprüfung:
# Abbruch mit Befund, nichts wird am Host geändert. Aufräumen per trap EXIT/INT/TERM (Hilfs-Container,
# Devcontainer nach Projekten und nach Name, ubmess-Container, Volumes Images, Arbeitsordner); danach
# Inventur nachher und Differenz zur Inventur vorher.
set -u

modus=${1:?usage: messen.sh <cli|var|identitaet|engine|bestand|config> <u-boot-binary>}
bin=$(readlink -f "${2:?usage: messen.sh <cli|var|identitaet|engine|bestand|config> <u-boot-binary>}")
here=$(cd "$(dirname "$0")" && pwd)
ctx="$here/../../../internal/adapter/driven/docker/devcontainercli"

HELPER_IMG=ubmess-devcontainer-cli:0.89.0
export GIT_TOKEN="platzhalter-kein-zugangsdatum-0123456789abcdef"

base=$(mktemp -d)
UP_LIMIT=300
STAGE_LIMIT=1200
NAME_BASE="ubmess-name"
HSEQ=0 PSEQ=0
LAST_RC=0 LAST_CID="" PDIR="" PNAME=""
DIRS=() NAMES=()
STATUS="vollstaendig"

section() { printf '\n######## %s\n' "$*"; }
now() { date +%s; }

inventur() {
  docker ps -a --format '{{.Names}}	{{.Image}}' | sort > "$1.ps"
  docker volume ls -q | sort > "$1.vol"
  docker images --format '{{.Repository}}:{{.Tag}}	{{.ID}}' | sort > "$1.img"
}

# ---- Aufräumen ----------------------------------------------------------------------------

teardown_dir() {
  local dir=$1 ids
  ids=$(docker ps -aq --filter "label=devcontainer.local_folder=$dir")
  # shellcheck disable=SC2086
  [ -n "$ids" ] && docker rm -f -v $ids > /dev/null 2>&1
}

cleanup() {
  trap - EXIT INT TERM
  section "Aufraeumen (trap)"
  local helpers ids n
  helpers=$(docker ps -aq --filter "name=ubmess-helper-")
  # shellcheck disable=SC2086
  [ -n "$helpers" ] && docker rm -f $helpers > /dev/null 2>&1
  for d in "${DIRS[@]:-}"; do [ -n "$d" ] && teardown_dir "$d"; done
  for n in "${NAMES[@]:-}"; do
    [ -n "$n" ] && docker rm -f -v "$n" > /dev/null 2>&1
  done
  local vols imgs v
  # Volumes, die die CLI für die Dev-Container anlegte, tragen den Projektnamen ubmess-<seq>- …
  for v in $(docker volume ls -q | grep '^ubmess-' || true); do docker volume rm -f "$v" > /dev/null 2>&1; done
  imgs=$(docker images -q --filter 'reference=vsc-ubmess-*')
  # shellcheck disable=SC2086
  [ -n "$imgs" ] && docker rmi -f $imgs > /dev/null 2>&1
  docker rmi -f "$HELPER_IMG" > /dev/null 2>&1
  rm -rf "$base" 2> /dev/null
  [ -e "$base" ] && echo "BEFUND: Arbeitsordner $base bleibt stehen" || echo "Arbeitsordner entfernt"
  section "Inventur nachher und Differenz"
  if [ -f "$INVDIR/vorher.ps" ]; then
    inventur "$INVDIR/nachher"
    for k in ps vol img; do
      echo "--- Differenz $k (vorher -> nachher; leer = keine)"
      diff "$INVDIR/vorher.$k" "$INVDIR/nachher.$k"
    done
    echo "--- Reste mit Praefix ubmess-/vsc-ubmess-:"
    { docker ps -a --format '{{.Names}}'; docker volume ls -q; docker images --format '{{.Repository}}'; } | grep -E '^(ubmess-|vsc-ubmess-)' || echo "keine"
  fi
  rm -rf "$INVDIR" 2> /dev/null
  section "Ende"
  echo "Stand des Modus: $STATUS"
}
INVDIR=$(mktemp -d)
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

# ---- Helfer -------------------------------------------------------------------------------

# up <ordner> [CLI-Optionen…]: ein `devcontainer up` im Hilfs-Container, gemessen.
up() {
  local dir=$1 hn t0
  shift
  HSEQ=$((HSEQ + 1))
  hn="ubmess-helper-$HSEQ"
  t0=$(now)
  timeout -k 10 "$UP_LIMIT" docker run --rm --name "$hn" \
    --mount "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock" \
    --mount "type=bind,source=$dir,target=$dir" -e GIT_TOKEN \
    "$HELPER_IMG" devcontainer up --workspace-folder "$dir" "$@" > "$base/$hn.out" 2> "$base/$hn.err"
  LAST_RC=$?
  if [ "$LAST_RC" -eq 124 ] || [ "$LAST_RC" -eq 137 ]; then
    echo "BEFUND: Zeitlimit (${UP_LIMIT} s) fuer devcontainer up [$*] erreicht; Hilfs-Container $hn wird entfernt. Letzte Zeilen der CLI:"
    docker rm -f "$hn" > /dev/null 2>&1
    tail -n 4 "$base/$hn.err" | cut -c1-200 | sed 's/^/    /'
    STATUS="Zeitlimit bei up [$*]"
    return 1
  fi
  LAST_CID=$(tail -n 1 "$base/$hn.out" | grep -o '"containerId":"[0-9a-f]*"' | cut -d'"' -f4)
  echo "up [$*] -> exit=$LAST_RC dauer=$(( $(now) - t0 )) s containerId=${LAST_CID:0:12}"
  echo "  Ergebnis: $(tail -n 1 "$base/$hn.out" | cut -c1-260)"
  if [ "$LAST_RC" -ne 0 ] || [ -z "$LAST_CID" ]; then
    echo "  stderr (letzte 10 Zeilen):"
    tail -n 10 "$base/$hn.err" | cut -c1-300 | sed 's/^/    /'
    return 1
  fi
  return 0
}

# cname <cid>: Name des Containers ohne fuehrenden Schraegstrich
cname() { docker inspect --format '{{.Name}}' "$1" | sed 's|^/||'; }

ident_of() { # <cid>: aufgeloester ${devcontainerId} aus dem Workspace-Volume (<projekt>-workspace-…)
  local n
  n=$(docker inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}{{if eq .Destination "/workspaces/'"$PNAME"'"}}{{.Name}}{{end}}{{end}}{{end}}' "$1")
  echo "${n#"$PNAME"-workspace-}"
}

vols_of() { # <cid>: Volume-Namen mit Zielordner
  docker inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}    {{.Name}} -> {{.Destination}}{{"\n"}}{{end}}{{end}}' "$1"
}

labels_of() { # <cid>
  docker inspect --format '  devcontainer.local_folder={{index .Config.Labels "devcontainer.local_folder"}}
  u-boot.instance={{index .Config.Labels "u-boot.instance"}}
  Zahl der Labels insgesamt: {{len .Config.Labels}}
  Image={{.Image}} Config.Image={{.Config.Image}}' "$1"
}

snapshot() { # <kennung> <cid>: dcid, Labels, Volumes — M3
  echo "[$1] containerId=${2:0:12} name=$(cname "$2")"
  echo "[$1] devcontainerId=$(ident_of "$2")"
  echo "[$1] Volumes:"
  vols_of "$2"
  labels_of "$2" | sed "s/^/[$1] /"
}

inject_property() { # <dir> <name>: Eigenschaft containerName vor der "name"-Zeile ergänzt
  local f=$1/.devcontainer/devcontainer.json
  sed -i 's/^  "name": \(.*\)$/  "containerName": "'"$2"'",\n  "name": \1/' "$f"
}

inject_runargs_name() { # <dir> <name>: --name=<name> an die runArgs-Liste angehängt bzw. Liste ergänzt
  local f=$1/.devcontainer/devcontainer.json
  if grep -q '"runArgs"' "$f"; then
    sed -i 's/"runArgs": \[/"runArgs": ["--name='"$2"'", /' "$f"
  else
    sed -i 's/^  "name": \(.*\)$/  "runArgs": ["--name='"$2"'"],\n  "name": \1/' "$f"
  fi
}

new_project() { # <standard|sandbox> -> PDIR, PNAME
  local profile=$1
  PSEQ=$((PSEQ + 1))
  PNAME="ubmess-$PSEQ"
  PDIR="$base/$PNAME"
  mkdir -p "$PDIR"
  DIRS+=("$PDIR")
  if [ "$profile" = sandbox ]; then
    ( cd "$PDIR" && "$bin" init --devcontainer --sandbox --no-git --yes > /dev/null 2>&1 )
    cat > "$PDIR/u-boot.yaml" <<EOF
schemaVersion: 1
project:
    name: $PNAME
devcontainer:
    enabled: true
    profile: sandbox
    sandbox:
        nestedRuntime: docker
EOF
    ( cd "$PDIR" && "$bin" generate devcontainer --yes > /dev/null 2>&1 )
  else
    ( cd "$PDIR" && "$bin" init --devcontainer --no-git --yes > /dev/null 2>&1 )
  fi
  echo "Projekt $PNAME ($profile) in $PDIR; devcontainer.json:"
  head -n 6 "$PDIR/.devcontainer/devcontainer.json" | sed 's/^/  json: /'
}

engine_probe() { # <name> [weitere docker run-Optionen…]: Name mit der Engine direkt (kein Hilfs-Container)
  local n=$1
  shift
  echo "-- Name: $(printf '%q' "$n") (Laenge ${#n})"
  local out rc
  out=$(timeout 30 docker run --rm --name "$n" --entrypoint /bin/true "$HELPER_IMG" "$@" 2>&1)
  rc=$?
  echo "   exit=$rc${out:+ Ausgabe:} $(echo "$out" | cut -c1-260)"
  docker inspect --type container --format '   Nachher am Container: {{.Name}}' "$n" 2>/dev/null || true
  # Nur bei Erfolg aufraeumen; bei einer Kollision gehoert der Name dem anderen Container,
  # dessen Entfernen hier den naechsten Probe-Zustand zerstoeren wuerde.
  if [ "$rc" -eq 0 ]; then
    NAMES+=("$n")
    docker rm -f -v "$n" > /dev/null 2>&1 || true
  fi
}

# ---- Modus cli (M1) -----------------------------------------------------------------------

mode_cli() {
  section "M1: Hilfe der CLI zu devcontainer up"
  docker run --rm "$HELPER_IMG" devcontainer --version 2>&1 | sed 's/^/  version: /'
  docker run --rm "$HELPER_IMG" devcontainer up --help 2>&1 | sed 's/^/  help: /'
  section "M1: Zeilen der Hilfe mit 'name'"
  docker run --rm "$HELPER_IMG" devcontainer up --help 2>&1 | grep -in 'name' | sed 's/^/  treffer: /' || echo "  treffer: keine"
  new_project standard
  local dir=$PDIR got ref
  section "M1: Gegenprobe ohne jede Namensquelle (Zufallsname erwartet)"
  if up "$dir"; then
    ref=$(cname "$LAST_CID")
    echo "  Name am Container: $ref"
    teardown_dir "$dir"
  else
    echo "  BEFUND M1: Referenzlauf endet exit=$LAST_RC; weitere Proben laufen trotzdem"
    teardown_dir "$dir"
  fi
  # Vor jedem weiteren Lauf ist die Eigenschaften-Datei am Bestand; eine Quelle wird gesetzt,
  # gemessen und wieder entfernt — der Name des Containers bleibt je Probe der Messwert.
  section "M1: Probe Eigenschaft containerName der devcontainer.json"
  inject_property "$dir" "$NAME_BASE-prop"
  head -n 4 "$dir/.devcontainer/devcontainer.json" | sed 's/^/    json: /'
  if up "$dir"; then
    got=$(cname "$LAST_CID")
    echo "  Name am Container: $got (gewuenscht: $NAME_BASE-prop)"
    [ "$got" = "$NAME_BASE-prop" ] && echo "  BEFUND M1a: Eigenschaft wird angewendet" || echo "  BEFUND M1a: Eigenschaft wird NICHT angewendet"
    teardown_dir "$dir"
  else
    echo "  BEFUND M1a: Eigenschaft — devcontainer up endet exit=$LAST_RC"
    teardown_dir "$dir"
  fi
  sed -i '/"containerName"/d' "$dir/.devcontainer/devcontainer.json"
  section "M1: Probe runArgs --name=<name> (statisch)"
  inject_runargs_name "$dir" "$NAME_BASE-runargs"
  grep -n 'runArgs' "$dir/.devcontainer/devcontainer.json" | sed 's/^/    json: /'
  if up "$dir"; then
    got=$(cname "$LAST_CID")
    echo "  Name am Container: $got (gewuenscht: $NAME_BASE-runargs)"
    [ "$got" = "$NAME_BASE-runargs" ] && echo "  BEFUND M1c: runArgs --name wird angewendet" || echo "  BEFUND M1c: runArgs --name wird NICHT angewendet"
    teardown_dir "$dir"
  else
    echo "  BEFUND M1c: runArgs --name — devcontainer up endet exit=$LAST_RC"
    teardown_dir "$dir"
  fi
}

# ---- Modus var (M2) -----------------------------------------------------------------------

mode_var() {
  section "M2: \${devcontainerId} in runArgs"
  new_project sandbox
  local dir=$PDIR
  inject_runargs_name "$dir" "${NAME_BASE}-var-\${devcontainerId}"
  echo "  devcontainer.json (Kopf):"; head -n 6 "$dir/.devcontainer/devcontainer.json" | sed 's/^/    json: /'
  if up "$dir"; then
    local got=$(cname "$LAST_CID")
    echo "  Name am Container: $got"
    echo "  devcontainerId des Containers: $(ident_of "$LAST_CID")"
    case "$got" in
      *'${devcontainerId}'*) echo "  BEFUND M2: NICHT aufgeloest (Name enthaelt den Platzhalter wörtlich)" ;;
      *ubmess-*) echo "  BEFUND M2: aufgeloest (Name enthaelt die devcontainerId)" ;;
      *) echo "  BEFUND M2: Name ohne Platzhalter — CLI hat --name nicht angewendet" ;;
    esac
    teardown_dir "$dir"
  else
    echo "  BEFUND M2: devcontainer up endet exit=$LAST_RC (stderr oben)"
    teardown_dir "$dir"
  fi
}

# ---- Modus identitaet (M3) ----------------------------------------------------------------

identitaet_cell() { # <standard|sandbox> <mit|ohne-instanz>
  local profile=$1 inst=$2
  section "M3: Zelle Profil=$profile Instanz=$inst"
  new_project "$profile"
  local dir=$PDIR
  local args=() refc namedc ref_dcid ref_vols ref_inst named_dcid named_vols named_inst named_name
  case $inst in
    mit) args=("--id-label" "devcontainer.local_folder=$dir" "--id-label" "u-boot.instance=i1") ;;
    *)   args=("--id-label" "devcontainer.local_folder=$dir") ;;
  esac
  echo "-- Referenz (ohne Namen)"
  if up "$dir" "${args[@]}"; then
    refc=$LAST_CID
    snapshot "ref" "$refc"
    ref_dcid=$(ident_of "$refc")
    ref_vols=$(vols_of "$refc")
    ref_inst=$(docker inspect --format '{{index .Config.Labels "u-boot.instance"}}' "$refc")
    docker rm -f -v "$refc" > /dev/null 2>&1
  else
    echo "  BEFUND M3: Referenzlauf endet exit=$LAST_RC — Zelle abgebrochen"
    teardown_dir "$dir"; return 1
  fi
  echo "-- Lauf mit Namen"
  inject_runargs_name "$dir" "$NAME_BASE-$profile-$inst"
  grep -n 'runArgs\|"name"' "$dir/.devcontainer/devcontainer.json" | sed 's/^/    json: /'
  if up "$dir" "${args[@]}"; then
    namedc=$LAST_CID
    named_name=$(cname "$namedc")
    snapshot "named" "$namedc"
    named_dcid=$(ident_of "$namedc")
    named_vols=$(vols_of "$namedc")
    named_inst=$(docker inspect --format '{{index .Config.Labels "u-boot.instance"}}' "$namedc")
    teardown_dir "$dir"
  else
    echo "  BEFUND M3: Lauf mit Namen endet exit=$LAST_RC — Zelle abgebrochen"
    teardown_dir "$dir"; return 1
  fi
  echo "  Vergleich devcontainerId: $([ "$ref_dcid" = "$named_dcid" ] && echo gleich || echo 'NICHT gleich') ($ref_dcid <-> $named_dcid)"
  echo "  Vergleich Volume-Namen:   $([ "$ref_vols" = "$named_vols" ] && echo gleich || echo 'NICHT gleich')"
  echo "  Vergleich Label-Wert:     $([ "$ref_inst" = "$named_inst" ] && echo gleich || echo 'NICHT gleich') ($ref_inst <-> $named_inst)"
  echo "  Name am Container (named-Lauf): $named_name (gewuenscht: $NAME_BASE-$profile-$inst)"
}

mode_identitaet() {
  identitaet_cell sandbox mit || true
  identitaet_cell sandbox ohne || true
  identitaet_cell standard mit || true
  identitaet_cell standard ohne || true
}

# ---- Modus bestand (M5) -------------------------------------------------------------------

mode_bestand() {
  new_project standard
  local dir=$PDIR n="$NAME_BASE-konflikt" foreignid stoppedid
  section "M5: Bestandscontainer gleicher Namen, gleicher Ordner (erwartet: uebernommen)"
  inject_runargs_name "$dir" "$n"
  if up "$dir"; then
    echo "  erster Lauf: $(cname "$LAST_CID") (${LAST_CID:0:12})"
    echo "-- zweiter Lauf, gleiches Projekt, gleicher Name"
    if up "$dir"; then
      echo "  zweiter Lauf: exit=0; Container $(cname "$LAST_CID") (${LAST_CID:0:12})"
      teardown_dir "$dir"
    else
      echo "  BEFUND M5: zweiter Lauf exit=$LAST_RC"
      teardown_dir "$dir"
    fi
  else
    echo "  BEFUND M5: erster Lauf exit=$LAST_RC; Zelle Kollision abgebrochen"
    teardown_dir "$dir"
    STATUS="erster Lauf des Modus bestand endet exit=$LAST_RC"
    return 1
  fi
  section "M5: Fremder Container gleicher Namen (andere local_folder, laufend)"
  docker run -d --name "$n" --entrypoint sleep -l "devcontainer.local_folder=$base/fremd" "$HELPER_IMG" 600 > /dev/null 2>&1
  foreignid=$(docker ps -aq --filter "name=^$n$"); NAMES+=("$n")
  echo "  Fremder Container: ${foreignid:0:12} laeuft=$(docker inspect --format '{{.State.Running}}' "$n")"
  echo "-- CLI-Lauf daneben"
  if up "$dir"; then
    echo "  BEFUND M5: CLI startet trotzdem — Container $(cname "$LAST_CID") (${LAST_CID:0:12}); fremder Container noch da: $(docker ps -aq --filter "name=^$n$" | wc -l) Treffer"
    teardown_dir "$dir"
    docker rm -f -v "$n" > /dev/null 2>&1
  else
    echo "  BEFUND M5: CLI endet exit=$LAST_RC; Fremdwortlaut oben; fremder Container noch da: $(docker ps -aq --filter "name=^$n$" | wc -l) Treffer"
  fi
  section "M5: --remove-existing-container neben fremdem Container gleichen Namens"
  if up "$dir" --remove-existing-container; then
    echo "  BEFUND M5: exit=0, Container $(cname "$LAST_CID"); fremder Container noch da: $(docker ps -aq --filter "name=^$n$" | wc -l) Treffer"
    teardown_dir "$dir"
  else
    echo "  BEFUND M5: --remove-existing-container exit=$LAST_RC; fremder Container noch da: $(docker ps -aq --filter "name=^$n$" | wc -l) Treffer"
  fi
  docker rm -f -v "$n" > /dev/null 2>&1
  section "M5: Gestoppter Container gleichen Namens, dann up und up --remove-existing-container"
  docker run -d --name "$n" --entrypoint sleep "$HELPER_IMG" 600 > /dev/null 2>&1
  NAMES+=("$n")
  docker stop "$n" > /dev/null 2>&1
  echo "  Gestoppter Container: ${n} Running=$(docker inspect --format '{{.State.Running}}' "$n")"
  if up "$dir"; then
    echo "  BEFUND M5: up startet trotzdem — Container $(cname "$LAST_CID")"
    teardown_dir "$dir"
  else
    echo "  BEFUND M5: up exit=$LAST_RC"
  fi
  if up "$dir" --remove-existing-container; then
    echo "  BEFUND M5: --remove-existing-container exit=0, Container $(cname "$LAST_CID"); gestoppter Container noch da: $(docker ps -aq --filter "name=^$n$" | wc -l) Treffer"
    teardown_dir "$dir"
  else
    echo "  BEFUND M5: --remove-existing-container exit=$LAST_RC; gestoppter Container noch da: $(docker ps -aq --filter "name=^$n$" | wc -l) Treffer"
  fi
  docker rm -f -v "$n" > /dev/null 2>&1
}

# ---- Modus config (M6: Name je Instanz über --override-config) -----------------------------

mode_config() {
  section "M6: Name je Aufruf über --override-config (Referenz, dann Name je Zelle)"
  # Der Hilfs-Aufruf würde eine je Aufruf gerenderte Konfigurationsdatei an die CLI geben; hier
  # steht die Kopie der erzeugten devcontainer.json, in die nur der Name eingetragen ist.
  new_project sandbox
  local dir=$PDIR
  local ov="$dir/devcontainer-ubmess-override.json"
  local ref_refc ref_dcid="" ref_vols="" ref_inst=""
  echo "-- Referenz (erzeugte Datei, ohne Namen)"
  if up "$dir" --id-label "devcontainer.local_folder=$dir" --id-label "u-boot.instance=i1"; then
    ref_refc=$LAST_CID
    snapshot "cfgref" "$ref_refc"
    ref_dcid=$(ident_of "$ref_refc")
    ref_vols=$(vols_of "$ref_refc")
    ref_inst=$(docker inspect --format '{{index .Config.Labels "u-boot.instance"}}' "$ref_refc")
    docker rm -f -v "$ref_refc" > /dev/null 2>&1
  else
    echo "  BEFUND M6: Referenzlauf endet exit=$LAST_RC"
    teardown_dir "$dir"
  fi
  echo "-- Lauf mit --override-config und Name in runArgs"
  cp "$dir/.devcontainer/devcontainer.json" "$ov"
  # shellcheck disable=SC2016 — sed-Muster mit Einzelanführung, der Name kommt als Einfügung.
  sed -i "s|\"runArgs\": \[|\"runArgs\": [\"--name=$NAME_BASE-cfg-i1\", |" "$ov"
  grep -o '"runArgs": \[[^]]*\]' "$ov" | sed 's/^/  override runArgs: /'
  if up "$dir" --override-config "$ov" --id-label "devcontainer.local_folder=$dir" --id-label "u-boot.instance=i1"; then
    local oc=$LAST_CID oname oname_dcid oname_vols oname_inst
    oname=$(cname "$oc")
    oname_dcid=$(ident_of "$oc")
    oname_vols=$(vols_of "$oc")
    oname_inst=$(docker inspect --format '{{index .Config.Labels "u-boot.instance"}}' "$oc")
    snapshot "cfgnamed" "$oc"
    teardown_dir "$dir"
    echo "  Vergleich devcontainerId: $([ "$ref_dcid" = "$oname_dcid" ] && echo gleich || echo 'NICHT gleich') ($ref_dcid <-> $oname_dcid)"
    echo "  Vergleich Volume-Namen:   $([ "$ref_vols" = "$oname_vols" ] && echo gleich || echo 'NICHT gleich')"
    echo "  Vergleich Label-Wert:     $([ "$ref_inst" = "$oname_inst" ] && echo gleich || echo 'NICHT gleich') ($ref_inst <-> $oname_inst)"
    echo "  Name am Container: $oname (gewuenscht: $NAME_BASE-cfg-i1)"
  else
    echo "  BEFUND M6: Lauf mit --override-config endet exit=$LAST_RC"
    teardown_dir "$dir"
  fi
  # Gegenprobe: ohne --override-config uebernimmt die CLI die erzeugte Datei; der Name im
  # override-File darf der erzeugten Datei auch nach dem Lauf fern bleiben.
  echo "  Erzeugte Datei enthaelt den Namen fern: $(grep -c "$NAME_BASE-cfg-i1" "$dir/.devcontainer/devcontainer.json") (0 erwartet)"
}

# ---- Modus engine (M4) --------------------------------------------------------------------

mode_engine() {
  section "M4: Namensregeln der Engine (docker run --name, kein Hilfs-Container)"
  local fill i long n
  for n in "$NAME_BASE-proj-inst" "9ubmess" "Uubmess" "ubmess.1" "ubmess_1"; do engine_probe "$n"; done
  fill="x"
  for i in 1 63 64 127 128 255 256; do
    long=""
    while [ ${#long} -lt "$i" ]; do long="${long}${fill}${i}"; done
    long=${long:0:i}
    engine_probe "a${long}"
  done
  section "M4: Zurückgewiesene Formen (Eingabewerte)"
  engine_probe "ubmess-led"
  # shellcheck disable=SC2016 — die Namen stehen wörtlich in der Ausgabe, sie sind der Messwert.
  for n in "-ubmess" ".ubmess" "_ubmess" "ub mess" 'ub"mess' "ub/mess" "ub:mess" "ub@mess" "übmess"; do
    engine_probe "$n" || true
  done
  section "M4: Konflikt — laufender Container gleiches Namens"
  docker run -d --name "$NAME_BASE-konf" "$HELPER_IMG" sleep 600 > /dev/null 2>&1
  NAMES+=("$NAME_BASE-konf")
  echo "  Vorhanden: $(docker ps --format '{{.Names}}' --filter "name=^$NAME_BASE-konf$")"
  echo "-- docker run --rm gleicher Name, laufend:"
  engine_probe "$NAME_BASE-konf" || true
  echo "  Vorhandener Container noch da: $(docker ps -aq --filter "name=^$NAME_BASE-konf$" | wc -l) Treffer"
  docker stop "$NAME_BASE-konf" > /dev/null 2>&1
  echo "-- docker run --rm gleicher Name, gestoppt:"
  engine_probe "$NAME_BASE-konf" || true
  echo "  Gestoppter Container noch da: $(docker ps -aq --filter "name=^$NAME_BASE-konf$" | wc -l) Treffer"
  docker rm -f -v "$NAME_BASE-konf" > /dev/null 2>&1
}

# ---- Ablauf -------------------------------------------------------------------------------

section "M0: Messstand"
date -u +"Zeit (UTC): %Y-%m-%dT%H:%M:%SZ"
echo "Modus: $modus"
echo "Host: $(uname -srm); uname -m: $(uname -m)"
echo "u-boot: $("$bin" --version 2>&1 | head -n 1) ($bin)"
docker version --format 'Docker Client {{.Client.Version}} ({{.Client.Os}}/{{.Client.Arch}}), Server {{.Server.Version}} ({{.Server.Os}}/{{.Server.Arch}})' 2>&1
docker info --format 'Engine: {{.OperatingSystem}}; Treiber {{.Driver}}; Cgroup {{.CgroupVersion}} ({{.CgroupDriver}}); CPUs {{.NCPU}}; MemTotal {{.MemTotal}}; Kernel {{.KernelVersion}}; Architektur {{.Architecture}}' 2>&1
free -m | head -n 3
echo "Plattenplatz /: $(df -h / | tail -n 1)"

section "Vorpruefung"
ok=1
docker info > /dev/null 2>&1 && echo "(a) docker info antwortet: ja" || { echo "(a) docker info antwortet: NEIN"; ok=0; }
avail_gb=$(df -BG --output=avail / | tail -n 1 | tr -dc 0-9)
echo "(d) freier Platz der Docker-Wurzel: ${avail_gb} GB (benoetigt 12 GB)"
[ "$avail_gb" -ge 12 ] && echo "(d) Platz: erfuellt" || { echo "(d) Platz: NICHT erfuellt"; ok=0; }
if [ "$ok" != 1 ]; then
  echo "BEFUND: Vorpruefung nicht erfuellt; Modus $modus wird nicht gemessen, nichts wurde angelegt oder geaendert."
  STATUS="nicht gemessen, Vorpruefung nicht erfuellt"
  exit 3
fi

section "Hilfs-Image bauen (ubmess-devcontainer-cli)"
t0=$(now)
docker build -q -t "$HELPER_IMG" "$ctx" > /dev/null 2> "$base/build.err" || { echo "BEFUND: Bau des Hilfs-Images scheitert"; tail -n 10 "$base/build.err"; STATUS="Abbruch: Hilfs-Image nicht baubar"; exit 4; }
echo "Bau: $(( $(now) - t0 )) s; Image $HELPER_IMG: $(docker images -q "$HELPER_IMG")"
echo "Dev-Container-CLI: $(docker run --rm "$HELPER_IMG" devcontainer --version 2>&1)"

# Inventur erst nach dem Bau des Mess-Images — sonst wuerde die Vorpruefung auf dem
# Mess-Image selbst abbrechen, der vorherige Lauf es in seinen Rohdaten veraendert
# und das Aufraeumen ein Objekt entfernt haben, das vor dem Lauf da war.
echo "(b) Inventur vorher (Rohdaten):"
inventur "$INVDIR/vorher"
pre=$( { cat "$INVDIR/vorher.ps" "$INVDIR/vorher.vol"; } | grep -cE '(^|[[:space:]])(vsc-)?ubmess-' || true)
echo "(c) kein Objekt mit Praefix ubmess-/vsc-ubmess- vorhanden: $([ "$pre" = 0 ] && echo ja || echo "NEIN ($pre)")"
[ "$pre" = 0 ] || { echo "BEFUND: Vorpruefung (c) nicht erfuellt (nach dem Bau des Mess-Images); Modus $modus wird nicht gemessen."; STATUS="nicht gemessen, Vorpruefung (c) nicht erfuellt"; exit 3; }

DEADLINE=$(( $(now) + STAGE_LIMIT ))
echo "Zeitlimits: je devcontainer up ${UP_LIMIT} s, je Lauf ${STAGE_LIMIT} s (Werte der Instanz-Messung: ein up dauert bis 7 s warm)"

"mode_$modus"