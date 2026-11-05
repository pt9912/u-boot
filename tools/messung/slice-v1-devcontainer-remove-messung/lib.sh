# Gemeinsame Definitionen der Skripte (wird mit `. ./lib.sh` eingebunden, nicht aufgerufen).
# Voraussetzungen: docker (Client und lokale Engine), Bash und Standardwerkzeuge (date, tee, mkdir, sort, comm, grep, sed, awk, timeout).
# Alles andere laeuft in Containern; die Messung fasst keine Host-Toolchain an.
#
# Umgebung der Skripte:
#   UBOOT  Pfad des u-boot-Binaries, das den Stand des Slice traegt (`make compile BAU_COMPILE_OUT=<pfad>`)
#   WORK   Messordner AUSSERHALB des Repos (erzeugte Projekte, Logs, Register); wird vom Aufrufer gewaehlt
#   OUT    Ausgabeverzeichnis der Rohdaten (Vorgabe: docs/messungen/slice-v1-devcontainer-remove-messung)
#
# Alle Docker-Objekte tragen das Praefix PFX: Projekt- und Volume-Namen beginnen mit ubmess-rm, die Container mit u-boot-ubmess-rm
# (so benennt `u-boot devcontainer up` den Container: u-boot-<projekt>[-<instanz>]), die Images mit vsc-ubmess-rm.
# Die Messung liest keine Zugangsdaten und setzt keine: GIT_TOKEN ist ein Platzhalter.

here=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ROOT=${ROOT:-$(cd "$here/../../.." && pwd)}
OUT=${OUT:-$ROOT/docs/messungen/slice-v1-devcontainer-remove-messung}
: "${UBOOT:?UBOOT (Pfad des u-boot-Binaries) fehlt}"
: "${WORK:?WORK (Messordner ausserhalb des Repos) fehlt}"
mkdir -p "$WORK" "$OUT"
WORK=$(cd "$WORK" && pwd)
OUT=$(cd "$OUT" && pwd)
UBOOT=$(cd "$(dirname "$UBOOT")" && pwd)/$(basename "$UBOOT")

PFX=ubmess-rm
PROBE_IMG=alpine:3
TMP=$WORK/$PFX-tmp   # Logs, Register und Zwischenergebnisse; 05-aufraeumen.sh entfernt den Ordner
mkdir -p "$TMP"
REGISTER=$TMP/volumes.txt
export GIT_TOKEN="platzhalter-kein-zugangsdatum-0123456789abcdef"

# Leitet stdout/stderr des Skripts in $OUT/<datei> und auf die Konsole.
rohdaten() {
  exec > >(tee "$OUT/$1") 2>&1
  trap 'exec >&- 2>&-; sleep 2' EXIT   # laesst den Filter die letzten Zeilen schreiben
}

now() { date -u +%Y-%m-%dT%H:%M:%SZ; }
mark() { echo; echo "---- $(now) $*"; }
# Eine stabile Aussage, die ein zweiter Lauf wiederholen muss (siehe README, Lesehilfen).
aussage() { echo "AUSSAGE $*"; }

# befehl <argumente...>  Gibt den Befehl, seine Ausgabe (eingerueckt) und den Exit-Code aus; setzt RC.
befehl() {
  echo "\$ $*"
  "$@" 2>&1 | sed 's/^/  /'
  RC=${PIPESTATUS[0]}
  echo "  -> exit=$RC"
}

dir_von() { echo "$WORK/$PFX-$1"; }
projektname() { echo "$PFX-$1"; }
# cname <variante> [instanz]: der Containername, den `u-boot devcontainer up` vergibt.
cname() { echo "u-boot-$PFX-$1${2:+-$2}"; }
# cid <variante> [instanz]: volle Container-ID oder leer.
cid() { docker ps -aq --no-trunc --filter "name=^/$(cname "$1" "${2:-}")\$"; }

# projekt <variante> <sandbox:0|1> [pfad=wert]...  Erzeugt das Projekt mit dem Binary: init --devcontainer, die Einstellungen per
# `config set`, dann `generate devcontainer`.
projekt() {
  local v=$1 sbx=$2; shift 2
  local d; d=$(dir_von "$v")
  rm -rf "$d"; mkdir -p "$d"
  ( cd "$d" || exit 1
    local flag=(); [ "$sbx" = 1 ] && flag=(--sandbox)
    "$UBOOT" init --devcontainer "${flag[@]}" --no-git --no-interactive "$(projektname "$v")" > /dev/null 2>&1 || { echo "init ist gescheitert"; exit 1; }
    local kv
    for kv in "$@"; do
      "$UBOOT" config set "${kv%%=*}" "${kv#*=}" > /dev/null 2>&1 || { echo "config set ${kv%%=*} ist gescheitert"; exit 1; }
    done
    [ $# -gt 0 ] && { "$UBOOT" generate devcontainer > /dev/null 2>&1 || { echo "generate devcontainer ist gescheitert"; exit 1; }; }
    exit 0
  ) || return 1
  echo "projekt $v: sandbox=$sbx $*"
  if grep -q '"mounts"' "$d/.devcontainer/devcontainer.json"; then
    grep -n '"mounts"' -A6 "$d/.devcontainer/devcontainer.json" | grep -E 'mounts|source=' | sed 's/^/  json: /'
  else
    echo "  json: kein Schluessel mounts"
  fi
  return 0
}

# hoch <variante> <instanz|""> [argumente...]  `u-boot devcontainer up`, gemessen. Setzt UP_RC. Die Zeile der Ausgabe nennt den Container.
hoch() {
  local v=$1 inst=$2; shift 2
  local args=(devcontainer up); [ -n "$inst" ] && args+=(--instance "$inst"); args+=("$@")
  local log="$TMP/up-$v-${inst:-ohne}.log" s e
  s=$(date +%s)
  ( cd "$(dir_von "$v")" && timeout 900 "$UBOOT" "${args[@]}" > "$log" 2>&1 ); UP_RC=$?
  e=$(date +%s)
  echo "u-boot ${args[*]} (Projekt $v): exit=$UP_RC dauer=$((e - s))s"
  if [ $UP_RC -eq 0 ]; then
    grep -E '^Devcontainer started|already running' "$log" | tail -n 1 | sed 's/^/  /'
  else
    tail -n 6 "$log" | cut -c1-300 | sed 's/^/  log: /'
  fi
  merke_volumes "$(cid "$v" "$inst")"
}

# runter <variante> [instanz]  `u-boot devcontainer down`.
runter() {
  local v=$1 inst=${2:-}
  local args=(devcontainer down); [ -n "$inst" ] && args+=(--instance "$inst")
  ( cd "$(dir_von "$v")" && "$UBOOT" "${args[@]}" 2>&1 | sed 's/^/  /'; exit "${PIPESTATUS[0]}" ); RC=$?
  echo "u-boot ${args[*]} (Projekt $v): exit=$RC"
}

# mounts <container-id>  Die Einhaengungen des Containers: Typ, Name, Ziel, schreibbar.
mounts() {
  docker inspect --format '{{range .Mounts}}  mount type={{.Type}} name={{.Name}} ziel={{.Destination}} rw={{.RW}}{{"\n"}}{{end}}' "$1"
}
# mount_volumes <container-id>  Die Namen der Volumes in Mounts, sortiert, eine je Zeile.
mount_volumes() {
  docker inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}{{.Name}}{{"\n"}}{{end}}{{end}}' "$1" | grep -v '^$' | sort
}
# merke_volumes <container-id>  Haelt die Volume-Namen des Containers im Register fest, damit das Aufraeumen auch anonyme findet.
merke_volumes() {
  [ -n "${1:-}" ] || return 0
  mount_volumes "$1" >> "$REGISTER"
}
vol_alle() { docker volume ls -q | sort; }

# objekte_entfernen [mit-images]  Entfernt Container, Volumes (Register und Praefix) und, mit Argument, Images, deren Name das Praefix traegt.
# Loescht nach dem Namen, nie nach Zeitfenstern, nie mit einem Glob ueber fremde Namen.
objekte_entfernen() {
  local n v pat_c="^(u-boot-)?$PFX" pat_v="^$PFX" pat_i="^vsc-$PFX"
  docker ps -a --format '{{.Names}}' | while read -r n; do
    [[ $n =~ $pat_c ]] && { merke_volumes "$n"; echo "entferne container $n"; docker rm -f -v "$n" > /dev/null; }
  done
  { [ -f "$REGISTER" ] && cat "$REGISTER"; docker volume ls -q | grep -E "$pat_v"; } | sort -u | while read -r v; do
    [ -n "$v" ] && docker volume inspect "$v" > /dev/null 2>&1 && { echo "entferne volume $v"; docker volume rm "$v" > /dev/null; }
  done
  if [ "${1:-}" = mit-images ]; then
    docker images --format '{{.Repository}}:{{.Tag}}' | while read -r n; do
      [[ $n =~ $pat_i ]] && { echo "entferne image $n"; docker rmi "$n" > /dev/null; }
    done
  fi
  return 0
}

praefix_objekte() {
  { docker ps -a --format 'container {{.Names}}'; docker volume ls -q | sed 's/^/volume /'; docker images --format 'image {{.Repository}}:{{.Tag}}'; } \
    | grep -E "^(container (u-boot-)?$PFX|volume $PFX|image vsc-$PFX)" || true
}

# inventur <praefix-datei>  Schreibt Container, Volumes und Images der Engine sortiert in <praefix>.ps/.vol/.img.
inventur() {
  docker ps -a --format '{{.Names}}' | sort > "$1.ps"
  docker volume ls -q | sort > "$1.vol"
  docker images --format '{{.Repository}}:{{.Tag}}' | sort > "$1.img"
}

# anfang / ende  Vorpruefung (keine Objekte mit dem Praefix), Inventur vorher; am Ende Aufraeumen der Container und Volumes und Differenz.
anfang() {
  echo "host: $(uname -m) docker=$(docker version --format '{{.Server.Version}} {{.Server.Os}}/{{.Server.Arch}}') u-boot: $("$UBOOT" --version 2>&1 | head -1)"
  echo "engine: $(docker info --format 'storage-driver={{.Driver}} cgroup={{.CgroupDriver}}/v{{.CgroupVersion}} security={{range .SecurityOptions}}{{.}};{{end}} os={{.OperatingSystem}} kernel={{.KernelVersion}}' 2>&1 | cut -c1-300)"
  echo "hilfs-image von u-boot devcontainer up: $(docker images --format '{{.Repository}}:{{.Tag}}' | grep '^u-boot-devcontainer-cli:' | tr '\n' ' ')"
  local rest; rest=$(praefix_objekte | grep -v '^image ' || true)
  if [ -n "$rest" ]; then echo "BEFUND: Objekte mit dem Praefix $PFX sind vorhanden (Aufraeumen davor):"; echo "$rest"; objekte_entfernen; fi
  inventur "$TMP/vorher-$$"
  trap 'ende' EXIT
}
ende() {
  trap - EXIT
  mark "Aufraeumen (Container und Volumes mit dem Praefix; Images bleiben bis 05-aufraeumen.sh)"
  objekte_entfernen
  inventur "$TMP/nachher-$$"
  for k in ps vol; do
    echo "--- Differenz $k (vorher -> nachher; leer = keine)"
    diff "$TMP/vorher-$$.$k" "$TMP/nachher-$$.$k"
  done
  echo "--- Reste mit dem Praefix (Container, Volumes): $(praefix_objekte | grep -vc '^image ' || true)"
  rm -f "$TMP/vorher-$$".* "$TMP/nachher-$$".*
  echo "== ende $(now)"
  exec >&- 2>&-; sleep 2
}
