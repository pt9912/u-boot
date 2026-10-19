# Gemeinsame Definitionen der Skripte (wird mit `. ./lib.sh` eingebunden, nicht aufgerufen).
# Messung des Slice slice-v1-egress-dns-upstream-schliessen (Fragen D1 bis D5). Ablage: docs/messungen/README.md.
# Voraussetzungen auf dem Host: docker, make und die Shell mit Standardwerkzeugen (bash, coreutils, grep, diff). awk, sed, jq, curl,
# dig laufen nur im Container.
#
# Umgebung der Skripte:
#   WORK  Messordner AUSSERHALB des Repos (erzeugte Projekte); wird vom Aufrufer gewaehlt
#   OUT   Ausgabeverzeichnis der Rohdaten (Pflicht: ${OUT:?}); das Protokoll des Slice liegt unter
#         docs/messungen/slice-v1-egress-dns-upstream-schliessen
#   ROOT  Repo-Wurzel (Vorgabe: drei Ebenen ueber diesem Verzeichnis)
set -u
here=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ROOT=${ROOT:-$(cd "$here/../../.." && pwd)}
: "${WORK:?WORK (Messordner ausserhalb des Repos) fehlt}"
: "${OUT:?OUT (Ausgabeverzeichnis der Rohdaten) fehlt}"
mkdir -p "$WORK" "$OUT"
WORK=$(cd "$WORK" && pwd)
OUT=$(cd "$OUT" && pwd)

# Alle Docker-Objekte der Messung tragen dieses Praefix (Aufraeumen arbeitet nur darauf).
PFX=u-boot-messung-dns
UB=$PFX-ub                       # Runtime-Image von u-boot (make build IMAGE=$UB)
BASE=debian:stable               # Werkzeug-Container (sed fuer die Schwaerzung, Projekt-Erzeugung)

# u-boot im Messordner: u <unterverzeichnis> <argumente...>
u() { local d=$1; shift; docker run --rm -v "$WORK/$d":/w -w /w -u "$(id -u):$(id -g)" -e HOME=/tmp "$UB" "$@"; }

# Schwaerzt vorsignierte URLs, Zugangsdaten und Home-Pfade, bevor eine Zeile in eine Rohdatei gelangt.
redact() {
  docker run --rm -i "$BASE" sed -u -E \
    -e 's/\?[^" ]*(X-Amz|Signature|Credential|token=|sig=)[^" ]*/?[Query entfernt: vorsignierte URL]/g' \
    -e 's/(Authorization|Bearer|X-Amz-[A-Za-z-]*)[:= ]+[^ ]+/\1 [entfernt]/g' \
    -e "s#/home/$(id -un)#/home/<benutzer>#g"
}

# Leitet stdout/stderr des Skripts in $OUT/<datei> (geschwaerzt) und auf die Konsole.
rohdaten() {
  exec > >(redact | tee "$OUT/$1") 2>&1
  trap 'exec >&- 2>&-; sleep 3' EXIT   # laesst den Filter die letzten Zeilen schreiben
}

now() { date -u +%Y-%m-%dT%H:%M:%S.%3NZ; }
# Eine stabile Aussage, die ein zweiter Lauf wiederholen muss (vergleich.sh).
aussage() { echo "AUSSAGE $*"; }

# Bestand der Docker-Objekte auf stdout: je Art nur die Anzahl und die Objekte, deren Name das Praefix der Messung traegt.
# Namen anderer Objekte des Hosts werden nirgends ausgegeben.
bestand() {
  local pat="^$PFX" n
  echo "Namen fremder Objekte nicht ausgewiesen, nur Anzahl"
  echo "container=$(docker ps -a -q | wc -l) images=$(docker images -q | wc -l) volumes=$(docker volume ls -q | wc -l) netze=$(docker network ls -q | wc -l)"
  echo "-- Objekte mit dem Praefix $PFX"
  docker ps -a --format '{{.Names}}' | while read -r n; do [[ $n =~ $pat ]] && echo "container $n"; done
  docker images --format '{{.Repository}}:{{.Tag}}' | while read -r n; do [[ $n =~ $pat ]] && echo "image $n"; done
  docker volume ls -q | while read -r n; do [[ $n =~ $pat ]] && echo "volume $n"; done
  docker network ls --format '{{.Name}}' | while read -r n; do [[ $n =~ $pat ]] && echo "netz $n"; done
  return 0
}
# Entfernt die Objekte, die bestand() mit dem Praefix ausweist. Loescht nach dem Namen, nie nach Argumenten, ungetaggten Images
# oder Zeitfenstern.
praefix_entfernen() {
  local pat="^$PFX" n
  docker ps -a --format '{{.Names}}' | while read -r n; do [[ $n =~ $pat ]] && { echo "entferne container $n"; docker rm -f -v "$n" > /dev/null; }; done
  docker volume ls -q | while read -r n; do [[ $n =~ $pat ]] && { echo "entferne volume $n"; docker volume rm "$n" > /dev/null; }; done
  docker network ls --format '{{.Name}}' | while read -r n; do [[ $n =~ $pat ]] && { echo "entferne netz $n"; docker network rm "$n" > /dev/null; }; done
  docker images --format '{{.Repository}}:{{.Tag}}' | while read -r n; do [[ $n =~ $pat ]] && { echo "entferne image $n"; docker rmi "$n" > /dev/null; }; done
  return 0
}

# Images der Messung
IMG_SBX=$PFX-img-sbx            # erzeugtes Sandbox-Image mit Egress (kein nestedRuntime)
IMG_POD=$PFX-img-pod            # erzeugtes Sandbox-Image mit Egress und nestedRuntime podman
IMG_TOOL=$PFX-img-tool          # IMG_SBX plus dig, nc, ping (Dockerfile.messung)
vdir() { echo "$WORK/p-$1"; }

# Container der Messung: start <name> <image> [docker-run-argumente...]. Wie die erzeugten runArgs: --cap-add=NET_ADMIN; Prozess 1 ist
# `sleep infinity` unter --init (Waisen werden eingesammelt). Das Verzeichnis der Skripte liegt unter /m (nur lesbar).
start() {
  local n=$1 i=$2; shift 2
  docker rm -f -v "$n" > /dev/null 2>&1
  docker run -d --init --name "$n" --cap-add=NET_ADMIN -v "$here":/m:ro "$@" "$i" sleep infinity > /dev/null
}
# postStartCommand als vscode (wie devcontainer.json): `sudo sh /usr/local/bin/u-boot-egress-init`; Modus skuid und Zaehler laufen ueber
# variante.sh (Messvariante). egress <name> <ist|skuid> [zaehler]
egress() {
  local n=$1 m=$2 z=${3:-}
  if [ "$m" = ist ] && [ -z "$z" ]; then docker exec -u vscode "$n" sudo sh /usr/local/bin/u-boot-egress-init
  else docker exec -u vscode "$n" sudo sh /m/variante.sh "$m" $z; fi
  echo "postStart ($m${z:+, $z}) exit=$?"
}
# Upstream-Liste des Containers (IPv4-nameserver der gesicherten Original-Datei)
upstreams() { docker exec "$1" awk '/^nameserver/ && $2 !~ /:/ {print $2}' /etc/resolv.conf.u-boot-orig | tr '\n' ' '; }
