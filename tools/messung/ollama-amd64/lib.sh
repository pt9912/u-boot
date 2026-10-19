# Gemeinsame Definitionen der Skripte (wird mit `. ./lib.sh` eingebunden, nicht aufgerufen).
# Voraussetzungen auf dem Host: docker, make und die Shell (bash, date, tee, mkdir, cp, rm, ls, wc, grep, diff).
# Jede Auswertung (jq, sed, awk, sha256sum, tar, curl, ss, ip) laeuft in einem Container.
#
# Umgebung der Skripte:
#   WORK  Messordner AUSSERHALB des Repos (Downloads, erzeugte Projekte); wird vom Aufrufer gewaehlt
#   OUT   Ausgabeverzeichnis der Rohdaten und der erzeugten Dateien (Vorgabe: das Evidence-Verzeichnis dieser Messung im Repo)
#   ROOT  Repo-Wurzel (Vorgabe: drei Ebenen ueber diesem Verzeichnis)

here=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ROOT=${ROOT:-$(cd "$here/../../.." && pwd)}
OUT=${OUT:-$ROOT/docs/messungen/slice-v1-devcontainer-ollama-amd64-messung}
: "${WORK:?WORK (Messordner ausserhalb des Repos) fehlt}"
mkdir -p "$WORK" "$OUT"
WORK=$(cd "$WORK" && pwd)
OUT=$(cd "$OUT" && pwd)

# Alle Docker-Objekte der Messung tragen dieses Praefix (Aufraeumen in 11-aufraeumen.sh arbeitet nur darauf).
PFX=u-boot-messung-amd64
UB=$PFX-ub                       # Runtime-Image von u-boot (make build IMAGE=$UB)
BASE=mcr.microsoft.com/devcontainers/base:debian
OLLAMA_VERSION=0.40.6
SHA_AMD64=c164ab8b1b28e2a247c29719b8e404c423c5b83ae64c3ab4e59d50c7ff7fc5e0
ARCHIVE=ollama-linux-amd64.tar.zst
ARCHIVE_URL=https://github.com/ollama/ollama/releases/download/v$OLLAMA_VERSION/$ARCHIVE

# Werkzeug-Container: das Basis-Image selbst (traegt jq, curl, sed, awk, ss, ip). Der Messordner ist unter /w eingehaengt.
tool() { docker run --rm -i -u "$(id -u):$(id -g)" -v "$WORK":/w -w /w -e HOME=/tmp "$BASE" "$@"; }

# u-boot im Messordner: u <unterverzeichnis> <argumente...>
u() { local d=$1; shift; docker run --rm -v "$WORK/$d":/w -w /w -u "$(id -u):$(id -g)" -e HOME=/tmp "$UB" "$@"; }

# Schwaerzt vorsignierte URLs und Zugangsdaten, bevor eine Zeile in eine Rohdatei gelangt.
redact() {
  docker run --rm -i "$BASE" sed -u -E \
    -e 's/\?[^" ]*(X-Amz|Signature|Credential|token=|sig=)[^" ]*/?[Query entfernt: vorsignierte URL]/g' \
    -e 's/(Authorization|Bearer|X-Amz-[A-Za-z-]*)[:= ]+[^ ]+/\1 [entfernt]/g' \
    -e "s#/home/$(id -un)#/home/<benutzer>#g"
}

# Leitet stdout/stderr des Skripts in $OUT/<datei> (geschwaerzt) und auf die Konsole.
rohdaten() {
  exec > >(redact | tee "$OUT/$1") 2>&1
  trap 'exec >&- 2>&-; sleep 3' EXIT   # lässt den Filter die letzten Zeilen schreiben
}

now() { date -u +%Y-%m-%dT%H:%M:%S.%3NZ; }
mark() { echo "---- $(now) $*"; }
sekunden() { awk -v a="$1" -v b="$2" 'BEGIN{printf "%.1f", b-a}'; }
# Eine stabile Aussage, die ein zweiter Lauf wiederholen muss (siehe README, Abschnitt Wiederholung).
aussage() { echo "AUSSAGE $*"; }

# Anzahl der Zeilen mit einem Muster (fuer Aussagen).
anzahl() { grep -c -E "$1" || true; }

# Variante -> Projektverzeichnis, Projektname und Image-Tag
vdir() { echo "$WORK/p-$1"; }
vname() { echo "$PFX-$1"; }
vimg() { echo "$PFX-img-$1"; }
