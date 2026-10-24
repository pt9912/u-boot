# Gemeinsame Definitionen der Skripte (wird mit `. ./lib.sh` eingebunden, nicht aufgerufen).
# Voraussetzungen auf dem Host: docker, make und die Shell mit Standardwerkzeugen (bash, coreutils: date, tee, mkdir, cp, cat, cut,
# tr, head, tail, id, sleep, wc, sort; grep, diff). awk, sed, jq, curl, sha256sum laufen nur im Container.
#
# Umgebung der Skripte:
#   WORK  Messordner AUSSERHALB des Repos (erzeugte Projekte); wird vom Aufrufer gewaehlt
#   OUT   Ausgabeverzeichnis der Rohdaten (Vorgabe: docs/messungen/slice-v1-ollama-cloud-messung)
#   ROOT  Repo-Wurzel (Vorgabe: drei Ebenen ueber diesem Verzeichnis)
#   LAUF  Nummer des Messlaufs fuer 04-lauf.sh (Vorgabe 1); sie steht im Dateinamen der Rohdaten
#
# Kein Skript liest, schreibt oder uebergibt ein Zugangsdatum. Die Messung laeuft ohne Konto: kein Token, keine Anmeldung. Die Skripte
# starten sich mit einer BEREINIGTEN Umgebung neu (`env -i` mit PATH, HOME, USER, DOCKER_HOST, WORK, OUT); danach prueft
# umgebung_pruefen() die Namen (nie die Werte) der verbleibenden Variablen und bricht bei einem Treffer ab. Ein Messcontainer erhaelt
# nie ein Host-Verzeichnis mit Konfigurationsdateien eines Agenten und nie eine Variable des Hosts.

here=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ROOT=${ROOT:-$(cd "$here/../../.." && pwd)}
OUT=${OUT:-$ROOT/docs/messungen/slice-v1-ollama-cloud-messung}
: "${WORK:?WORK (Messordner ausserhalb des Repos) fehlt}"
mkdir -p "$WORK" "$OUT"
WORK=$(cd "$WORK" && pwd)
OUT=$(cd "$OUT" && pwd)

# Neustart mit bereinigter Umgebung (einmal je Skriptaufruf); $0 und die Argumente des Skripts bleiben.
if [ -z "${UB_BEREINIGT:-}" ]; then
  exec env -i UB_BEREINIGT=1 PATH="$PATH" HOME="$HOME" USER="$(id -un)" ${DOCKER_HOST:+DOCKER_HOST="$DOCKER_HOST"} \
    ${LAUF:+LAUF="$LAUF"} WORK="$WORK" OUT="$OUT" ROOT="$ROOT" bash "$0" "$@"
fi

# Verbotene Namen: Zugangs- und Agenten-Variablen. Gibt die Namen aus, die zutreffen (ohne Werte); Rueckgabe 1 bei einem Treffer.
umgebung_pruefen() {
  local n treffer=0
  for n in $(env | cut -d= -f1); do
    case "$n" in
      ANTHROPIC_*|CLAUDE*|GH_TOKEN|GITHUB_TOKEN|GITLAB_TOKEN|GIT_TOKEN|GIT_ASKPASS|OLLAMA*|*TOKEN*|*SECRET*|*PASSWORD*|*API_KEY*|*AUTH*|OPENAI*|GEMINI*)
        echo "VERBOTENER NAME in der Umgebung: $n"; treffer=1 ;;
    esac
  done
  return $treffer
}

# Alle Docker-Objekte der Messung tragen dieses Praefix (Aufraeumen in 07-aufraeumen.sh arbeitet nur darauf).
PFX=u-boot-messung-ocloud
UB=$PFX-ub                       # Runtime-Image von u-boot (make build IMAGE=$UB)
BASE=mcr.microsoft.com/devcontainers/base:debian
CLI=$PFX-devcontainer-cli        # Dev-Container-CLI im Wegwerf-Container (Dockerfile.cli)

# Werkzeug-Container: das Basis-Image selbst (traegt sed, awk, jq, curl). Der Messordner ist unter /w eingehaengt. Keine Umgebung des Hosts.
tool() { docker run --rm -i -u "$(id -u):$(id -g)" -v "$WORK":/w -w /w -e HOME=/tmp "$BASE" "$@"; }

# u-boot im Messordner: u <unterverzeichnis> <argumente...>
u() { local d=$1; shift; docker run --rm -v "$WORK/$d":/w -w /w -u "$(id -u):$(id -g)" -e HOME=/tmp "$UB" "$@"; }

# Schwaerzt vorsignierte URLs, Autorisierungs- und Token-Header (Gross-/Kleinschreibung gleich), Zuweisungen an Namen mit TOKEN, SECRET, PASSWORD
# oder API_KEY, JSON-Felder mit solchen Namen, Zugangsdaten in URLs, bekannte Token-Praefixe, Anmelde-URLs und -Codes von `ollama signin`
# (Query und Fragment einer URL, jede URL mit `connect`/`signin`/`device`/`key=`) und den Home-Pfad, bevor eine Zeile in eine Rohdatei gelangt.
# Grenze: ein Token ohne Praefix und ohne eine dieser Formen (etwa als nacktes Wort in einer Fehlermeldung) wird nicht erkannt; Verweise auf
# eine Variable (Wert beginnt mit $ oder %) und Zeilen AUSSAGE bleiben stehen; redact-probe.sh pinnt die Formen, nicht die Vollstaendigkeit.
# Deshalb werden die Rohdaten zusaetzlich nach Tokenformen durchsucht (06-suche.sh) und von Hand gelesen.
redact() {
  docker run --rm -i "$BASE" sed -u -E \
    -e 's/\?[^" ]*(X-Amz|Signature|Credential|token=|sig=|X-Goog)[^" ]*/?[Query entfernt: vorsignierte URL]/gI' \
    -e 's#(https?://[^ "?]*(connect|signin|sign-in|device|login|oauth)[^ "?]*)[?#][^ "]*#\1?[Query entfernt: Anmelde-URL]#gI' \
    -e 's/(Authorization:? +)([A-Za-z]+ +)?[^ ]+/\1[entfernt]/gI' \
    -e 's/(Bearer|Basic) +[^ ]+/\1 [entfernt]/gI' \
    -e 's/(Authorization|Bearer|X-Amz-[A-Za-z-]*)[:= ]+[^ ]+/\1 [entfernt]/gI' \
    -e 's/((private|job|deploy|auth|access)-token|x-[a-z-]*token): *[^ ]+/\1: [entfernt]/gI' \
    -e '/^AUSSAGE /!s/([A-Za-z0-9_.-]*(token|secret|passw(or)?d|api_?key)[A-Za-z0-9_.-]*)=[^ $%][^ ]*/\1=[entfernt]/gI' \
    -e 's/"([A-Za-z0-9_]*(token|secret|passw(or)?d|api_?key|authorization)[A-Za-z0-9_]*)" *: *"[^"$][^"]*"/"\1":"[entfernt]"/gI' \
    -e 's#(://[^/:@ ]*):[^/@ ]+@#\1:[entfernt]@#g' \
    -e 's#(x-access-token|oauth2):[^/@ ]+@#\1:[entfernt]@#gI' \
    -e 's/(ghp_|gho_|ghu_|ghs_|ghr_|github_pat_|glpat-|glptt-|gldt-|glrt-|sk-ant-|sk-|ollama_)[A-Za-z0-9_-]{8,}/\1[entfernt]/g' \
    -e 's/\b(code|key|nonce|public key|ssh-ed25519|ssh-rsa)([ :=]+)[A-Za-z0-9+\/_=-]{16,}/\1\2[entfernt]/gI' \
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
# Eine stabile Aussage, die ein zweiter Lauf wiederholen muss (siehe README, Abschnitt Lesehilfen).
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

# Variante -> Projektverzeichnis, Projektname. Varianten: std (Standardprofil, C7), dis (Sandbox ohne Egress: Modus discover),
# res (Sandbox mit Egress: Modus restricted), alw (wie res plus Eintrag ollama.com in egress.allow, C2b).
vdir() { echo "$WORK/$PFX-$1"; }
vname() { echo "$PFX-$1"; }
