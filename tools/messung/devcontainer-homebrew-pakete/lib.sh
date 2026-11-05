# Gemeinsame Definitionen der Skripte (wird mit `. ./lib.sh` eingebunden, nicht aufgerufen).
# Voraussetzungen auf dem Host: docker, die Shell mit Standardwerkzeugen (bash, coreutils: date, tee, mkdir, cat, cut, tr, head, wc,
# sort, id; grep, sed, awk). Alles andere laeuft in Containern.
#
# Umgebung der Skripte:
#   UBOOT  Pfad des u-boot-Binaries, das den Stand des Slice traegt (erzeugt die Projekte; `make compile BAU_COMPILE_OUT=<pfad>`)
#   WORK   Messordner AUSSERHALB des Repos (erzeugte Projekte); wird vom Aufrufer gewaehlt
#   OUT    Ausgabeverzeichnis der Rohdaten (Vorgabe: docs/messungen/slice-v1-devcontainer-homebrew-pakete)
#   ROOT   Repo-Wurzel (Vorgabe: drei Ebenen ueber diesem Verzeichnis)
#
# Die Messung liest keine Zugangsdaten und setzt keine: kein Token, keine Anmeldung. Alle Docker-Objekte tragen das Praefix PFX.

here=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ROOT=${ROOT:-$(cd "$here/../../.." && pwd)}
OUT=${OUT:-$ROOT/docs/messungen/slice-v1-devcontainer-homebrew-pakete}
: "${UBOOT:?UBOOT (Pfad des u-boot-Binaries) fehlt}"
: "${WORK:?WORK (Messordner ausserhalb des Repos) fehlt}"
mkdir -p "$WORK" "$OUT"
WORK=$(cd "$WORK" && pwd)
OUT=$(cd "$OUT" && pwd)
UBOOT=$(cd "$(dirname "$UBOOT")" && pwd)/$(basename "$UBOOT")

PFX=u-boot-messung-brew
BASE=mcr.microsoft.com/devcontainers/base:debian

# Die Pin-Revision des Templates (Release 7.0.12) und die des vorherigen Releases (7.0.10), gelesen mit
# `git ls-remote --tags https://github.com/Homebrew/brew` (Tags 7.0.12 und 7.0.10; die Tags sind leichtgewichtig, der Hash ist der Commit).
PIN=bf464f57275cae0be106b43d190c98564ac59b0e
PIN_ALT=8f10429b56be28a25032b87463fffc95e556b5bc

# Leitet stdout/stderr des Skripts in $OUT/<datei> und auf die Konsole.
rohdaten() {
  exec > >(tee "$OUT/$1") 2>&1
  trap 'exec >&- 2>&-; sleep 2' EXIT   # laesst den Filter die letzten Zeilen schreiben
}

now() { date -u +%Y-%m-%dT%H:%M:%SZ; }
mark() { echo "---- $(now) $*"; }
# Differenz zweier `date +%s` in Sekunden.
dauer() { echo $(( $2 - $1 )); }
# Eine stabile Aussage, die ein zweiter Lauf wiederholen muss (siehe README, Lesehilfen).
aussage() { echo "AUSSAGE $*"; }

vdir() { echo "$WORK/p-$1"; }
vimg() { echo "$PFX-img-$1"; }

# projekt <variante> <sandbox:0|1> <pfad=wert>...  Erzeugt das Projekt mit dem Binary: init --devcontainer, die Einstellungen in der
# angegebenen Reihenfolge per `config set`, dann `generate devcontainer`. Ein Pfad mit Komma im Wert ist eine Liste.
projekt() {
  local v=$1 sbx=$2; shift 2
  local d; d=$(vdir "$v")
  rm -rf "$d"; mkdir -p "$d"
  ( cd "$d" || exit 1
    local flag=(); [ "$sbx" = 1 ] && flag=(--sandbox)
    "$UBOOT" init --devcontainer "${flag[@]}" --no-git --no-interactive "$PFX-$v" > /dev/null || exit 1
    local kv
    for kv in "$@"; do
      "$UBOOT" config set "${kv%%=*}" "${kv#*=}" > /dev/null || { echo "config set ${kv%%=*} ist gescheitert"; exit 1; }
    done
    "$UBOOT" generate devcontainer 2>&1 | sed 's/^/generate: /'
  )
}

# bild <variante> [docker-build-argumente...]  Baut das Dockerfile des Projekts; gibt Dauer und Groesse aus, Rueckgabe = Exit-Code des Baus.
bild() {
  local v=$1; shift
  local s e rc
  s=$(date +%s)
  docker build -q -t "$(vimg "$v")" "$@" "$(vdir "$v")/.devcontainer" > "$WORK/bau-$v.log" 2>&1; rc=$?
  e=$(date +%s)
  echo "bau $v: exit=$rc dauer=$(dauer "$s" "$e")s"
  [ $rc -ne 0 ] && tail -15 "$WORK/bau-$v.log" | sed 's/^/  bau-log: /'
  [ $rc -eq 0 ] && echo "image $v: bytes=$(docker image inspect --format '{{.Size}}' "$(vimg "$v")")"
  return $rc
}

# Entfernt die Objekte, deren Name das Praefix traegt (Container, Images). Loescht nach dem Namen, nie nach Zeitfenstern.
praefix_entfernen() {
  local pat="^$PFX" n
  docker ps -a --format '{{.Names}}' | while read -r n; do [[ $n =~ $pat ]] && { echo "entferne container $n"; docker rm -f -v "$n" > /dev/null; }; done
  docker images --format '{{.Repository}}:{{.Tag}}' | while read -r n; do [[ $n =~ $pat ]] && { echo "entferne image $n"; docker rmi "$n" > /dev/null; }; done
  return 0
}
