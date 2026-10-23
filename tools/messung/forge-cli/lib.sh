# Gemeinsame Definitionen der Skripte (wird mit `. ./lib.sh` eingebunden, nicht aufgerufen).
# Voraussetzungen auf dem Host: docker, make und die Shell mit Standardwerkzeugen (bash, coreutils: date, tee, mkdir, cp, cat, cut,
# tr, head, tail, id, sleep, wc, sort; grep, diff). awk, sed, jq, curl, sha256sum, tar laufen nur im Container.
#
# Umgebung der Skripte:
#   WORK  Messordner AUSSERHALB des Repos (erzeugte Projekte, Downloads); wird vom Aufrufer gewaehlt
#   OUT   Ausgabeverzeichnis der Rohdaten (Vorgabe: docs/messungen/slice-v1-forge-cli-quellen-messung)
#   ROOT  Repo-Wurzel (Vorgabe: drei Ebenen ueber diesem Verzeichnis)
#
# Kein Skript liest, schreibt oder uebergibt ein echtes Zugangsdatum. Wo ein Werkzeug eine Umgebungsvariable mit Token erwartet, steht
# der Platzhalter $DUMMY (kein gueltiges Token); die Antwort der Forge darauf ist ein Messwert. Hostseitige Anmeldungen (gh auth)
# werden nicht verwendet; die Abfragen der Forge-APIs laufen anonym aus einem Container.

here=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ROOT=${ROOT:-$(cd "$here/../../.." && pwd)}
OUT=${OUT:-$ROOT/docs/messungen/slice-v1-forge-cli-quellen-messung}
: "${WORK:?WORK (Messordner ausserhalb des Repos) fehlt}"
mkdir -p "$WORK" "$OUT"
WORK=$(cd "$WORK" && pwd)
OUT=$(cd "$OUT" && pwd)

# Alle Docker-Objekte der Messung tragen dieses Praefix (Aufraeumen in 08-aufraeumen.sh arbeitet nur darauf).
PFX=u-boot-messung-forge
UB=$PFX-ub                       # Runtime-Image von u-boot (make build IMAGE=$UB)
BASE=mcr.microsoft.com/devcontainers/base:debian
DUMMY=MESSUNG-KEIN-ECHTES-TOKEN  # Platzhalter fuer eine Token-Variable; kein gueltiges Zugangsdatum

# Werkzeug-Container: das Basis-Image selbst (traegt jq, curl, sed, awk, gpg). Der Messordner ist unter /w eingehaengt, das Verzeichnis
# dieser Skripte unter /m (nur lesend).
tool() { docker run --rm -i -u "$(id -u):$(id -g)" -v "$WORK":/w -v "$here":/m:ro -w /w -e HOME=/tmp "$BASE" "$@"; }

# u-boot im Messordner: u <unterverzeichnis> <argumente...>
u() { local d=$1; shift; docker run --rm -v "$WORK/$d":/w -w /w -u "$(id -u):$(id -g)" -e HOME=/tmp "$UB" "$@"; }

# Schwaerzt vorsignierte URLs, Autorisierungs-Header, bekannte Token-Praefixe und den Home-Pfad, bevor eine Zeile in eine Rohdatei gelangt.
redact() {
  docker run --rm -i "$BASE" sed -u -E \
    -e 's/\?[^" ]*(X-Amz|Signature|Credential|token=|sig=|X-Goog)[^" ]*/?[Query entfernt: vorsignierte URL]/g' \
    -e 's/(Authorization|Bearer|X-Amz-[A-Za-z-]*)[:= ]+[^ ]+/\1 [entfernt]/g' \
    -e 's/(ghp_|gho_|ghu_|ghs_|ghr_|github_pat_|glpat-)[A-Za-z0-9_-]{8,}/\1[entfernt]/g' \
    -e "s#/home/$(id -un)#/home/<benutzer>#g"
}

# Leitet stdout/stderr des Skripts in $OUT/<datei> (geschwaerzt) und auf die Konsole.
rohdaten() {
  exec > >(redact | tee "$OUT/$1") 2>&1
  trap 'exec >&- 2>&-; sleep 3' EXIT   # laesst den Filter die letzten Zeilen schreiben
}

now() { date -u +%Y-%m-%dT%H:%M:%S.%3NZ; }
mark() { echo "---- $(now) $*"; }
# Differenz zweier `date +%s.%N`-Werte in Sekunden mit einer Nachkommastelle (Ganzzahl-Arithmetik der Shell).
sekunden() { local a=${1/./} b=${2/./}; local d=$(( (b - a) / 100000000 )); echo "$((d / 10)).$((d % 10))"; }
# Eine stabile Aussage, die ein zweiter Lauf wiederholen muss (siehe README, Abschnitt Wiederholung).
aussage() { echo "AUSSAGE $*"; }

# Bestand der Docker-Objekte auf stdout: je Art nur die Anzahl und die Objekte, deren Name das Praefix der Messung traegt (auch die
# Images der Dev-Container-CLI, vsc-<Praefix>-…). Namen anderer Objekte des Hosts werden nirgends ausgegeben.
bestand() {
  local pat="^(vsc-)?$PFX" n
  echo "Namen fremder Objekte nicht ausgewiesen, nur Anzahl"
  echo "container=$(docker ps -a -q | wc -l) images=$(docker images -q | wc -l) volumes=$(docker volume ls -q | wc -l) netze=$(docker network ls -q | wc -l)"
  echo "-- Objekte mit dem Praefix $PFX (Name beginnt mit $PFX oder vsc-$PFX)"
  docker ps -a --format '{{.Names}}' | while read -r n; do [[ $n =~ $pat ]] && echo "container $n"; done
  docker images --format '{{.Repository}}:{{.Tag}}' | while read -r n; do [[ $n =~ $pat ]] && echo "image $n"; done
  docker volume ls -q | while read -r n; do [[ $n =~ $pat ]] && echo "volume $n"; done
  docker network ls --format '{{.Name}}' | while read -r n; do [[ $n =~ $pat ]] && echo "netz $n"; done
  return 0
}
# Entfernt die Objekte, die bestand() mit dem Praefix ausweist. Loescht nach dem Namen, nie nach Argumenten, ungetaggten Images oder Zeitfenstern.
praefix_entfernen() {
  local pat="^(vsc-)?$PFX" n
  docker ps -a --format '{{.Names}}' | while read -r n; do [[ $n =~ $pat ]] && { echo "entferne container $n"; docker rm -f -v "$n" > /dev/null; }; done
  docker volume ls -q | while read -r n; do [[ $n =~ $pat ]] && { echo "entferne volume $n"; docker volume rm "$n" > /dev/null; }; done
  docker network ls --format '{{.Name}}' | while read -r n; do [[ $n =~ $pat ]] && { echo "entferne netz $n"; docker network rm "$n" > /dev/null; }; done
  docker images --format '{{.Repository}}:{{.Tag}}' | while read -r n; do [[ $n =~ $pat ]] && { echo "entferne image $n"; docker rmi "$n" > /dev/null; }; done
  return 0
}

# Variante -> Projektverzeichnis, Projektname und Image-Tag
vdir() { echo "$WORK/p-$1"; }
vname() { echo "$PFX-$1"; }
vimg() { echo "$PFX-img-$1"; }

# Quellen der Messung (OCI-Features): Variante -> Referenz ohne Schema
oci_ref() {
  case "$1" in
    gh)    echo ghcr.io/devcontainers/features/github-cli ;;
    ghx)   echo ghcr.io/devcontainers-extra/features/gh-cli ;;
    glabs) echo ghcr.io/skriptfabrik/devcontainer-features/gitlab-cli ;;
    glabp) echo ghcr.io/postfinance/devcontainer-features/gitlab-cli ;;
    glabz) echo ghcr.io/srzstephen/devcontainer-features/gitlab-cli ;;
  esac
}
