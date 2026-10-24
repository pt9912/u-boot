# Gemeinsame Definitionen der Skripte (wird mit `. ./lib.sh` eingebunden, nicht aufgerufen).
# Voraussetzungen auf dem Host: docker, make und die Shell mit Standardwerkzeugen (bash, coreutils, grep, diff). sed, awk, jq, curl,
# sha256sum und socat laufen nur in Containern.
#
# Umgebung der Skripte:
#   WORK  Messordner AUSSERHALB des Repos (erzeugte Projekte, Zaehler, Stufe 1 der Rohdaten); wird vom Aufrufer gewaehlt
#   ROH   Stufe 1 der Rohdaten (Vorgabe: $WORK/stufe1); die Skripte schreiben nur dorthin. Von dort kommt eine Datei erst nach der
#         Handlesung ins Repo (10-uebernehmen.sh = Stufe 2)
#   ROOT  Repo-Wurzel (Vorgabe: drei Ebenen ueber diesem Verzeichnis)
#   LAUF  Nummer des Messlaufs (Vorgabe 1); sie steht im Dateinamen der Rohdaten
#
# Zugang: Kein Skript liest, schreibt oder uebergibt ein Zugangsdatum und keine Datei des Hosts (kein ~/.ollama, kein ~/.claude). Der
# angemeldete Host-Server wird nur ueber das Netz aus einem Container angesprochen. Der Platzhalter `ollama` fuer ANTHROPIC_AUTH_TOKEN und
# `platzhalter` fuer OLLAMA_API_KEY stehen nur im Aufruf IM Container (inner.sh), nie in der Umgebung des Hosts. Die Umgebung des Skripts wird
# im selben Prozess bereinigt (Positivliste, `unset` aller anderen Namen); danach prueft umgebung_pruefen() die Namen (nie die Werte).

here=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ROOT=${ROOT:-$(cd "$here/../../.." && pwd)}
: "${WORK:?WORK (Messordner ausserhalb des Repos) fehlt}"
mkdir -p "$WORK"
WORK=$(cd "$WORK" && pwd)
ROH=${ROH:-$WORK/stufe1}
mkdir -p "$ROH"
ROH=$(cd "$ROH" && pwd)
BASEDIR=${BASEDIR:-$ROOT/docs/messungen/slice-v1-ollama-cloud-messung-mit-konto}

# Bereinigung im selben Prozess: Positivliste der exportierten Namen, alles andere wird entfernt. Von der aufrufenden Shell ererbte
# Funktionen (Eintraege BASH_FUNC_… in der Umgebung) werden vor den eigenen Definitionen entfernt.
unset -f $(compgen -A function) 2> /dev/null
bereinigen() {
  local n
  for n in $(compgen -e); do
    case "$n" in
      PATH|HOME|USER|DOCKER_HOST|WORK|ROH|ROOT|LAUF|BASEDIR|PWD|SHLVL|OLDPWD|_|PROMPT_MAX|KEEP_IMAGES) ;;
      *) unset "$n" 2> /dev/null ;;
    esac
  done
}
bereinigen
export PATH HOME

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

# Alle Docker-Objekte der Messung tragen dieses Praefix (Aufraeumen arbeitet nur darauf). Es beginnt NICHT mit dem Praefix der Messung ohne
# Konto, ist aber dessen Erweiterung; die Muster unten sind deshalb auf das Praefix mit "-konto" verankert.
PFX=u-boot-messung-ocloud-konto
UB=$PFX-ub                       # Runtime-Image von u-boot (make build IMAGE=$UB)
BASE=mcr.microsoft.com/devcontainers/base:debian
CLI=$PFX-devcontainer-cli        # Dev-Container-CLI im Wegwerf-Container (Dockerfile.cli)
RELAY=$PFX-relay                 # Hilfs-Container des Relays (K1b)
PORT=11434
PROMPT_MAX=${PROMPT_MAX:-40}     # Obergrenze der Prompts (Entscheidung des Projektinhabers vom 2026-10-24)
ZAEHLER=$WORK/prompts.txt        # eine Zeile je gestelltem Prompt: <Nr> <Zeit> <Fall>
M1='glm-5.3-flash:cloud[1m]'
M2='gpt-oss:120b-cloud'
PROMPT_TEXT='Reply with the single word OK.'

# Werkzeug-Container: das Basis-Image selbst (traegt sed, awk, jq, curl). Der Messordner ist unter /w eingehaengt. Keine Umgebung des Hosts.
tool() { docker run --rm -i -u "$(id -u):$(id -g)" -v "$WORK":/w -w /w -e HOME=/tmp "$BASE" "$@"; }

# u-boot im Messordner: u <unterverzeichnis> <argumente...>
u() { local d=$1; shift; docker run --rm -v "$WORK/$d":/w -w /w -u "$(id -u):$(id -g)" -e HOME=/tmp "$UB" "$@"; }

# Schwaerzt vorsignierte URLs, Autorisierungs- und Token-Header, Zuweisungen an Namen mit TOKEN, SECRET, PASSWORD oder API_KEY, JSON-Felder mit
# solchen Namen, Zugangsdaten in URLs, bekannte Token-Praefixe, Anmelde-URLs und -Codes von `ollama signin`, Schluesselbloecke und lange
# Base64-Zeilen (Geraeteschluessel), die Zeile "signed in as <Name>", E-Mail-Adressen, JSON-Felder einer Modell-Antwort (response, content,
# result, text, thinking), Zeilen "ANTWORT ..." und den Home-Pfad, bevor eine Zeile in eine Rohdatei gelangt.
# Grenze: ein Token ohne Praefix und ohne eine dieser Formen (etwa als nacktes Wort in einer Fehlermeldung) wird nicht erkannt; Verweise auf
# eine Variable (Wert beginnt mit $ oder %) und Zeilen AUSSAGE bleiben stehen; redact-probe.sh pinnt die Formen, nicht die Vollstaendigkeit.
# Deshalb nehmen die Skripte keinen Antwort-Inhalt auf (inner.sh), durchsucht 09-suche.sh die Rohdaten und werden sie von Hand gelesen.
redact() {
  docker run --rm -i "$BASE" sed -u -E \
    -e '/^-----BEGIN [A-Z ]*PRIVATE KEY-----/,/^-----END [A-Z ]*PRIVATE KEY-----/c\[Schluesselblock entfernt]' \
    -e 's/^[A-Za-z0-9+\/]{40,}={0,2}$/[lange Base64-Zeile entfernt]/' \
    -e 's/^ANTWORT.*$/ANTWORT [entfernt]/' \
    -e 's/"(response|content|result|text|thinking)" *: *"([^"\\]|\\.)*"/"\1":"[Antwort entfernt]"/gI' \
    -e 's/(signed in as) +[^ ]+/\1 [entfernt]/gI' \
    -e 's/[A-Za-z0-9._%+-]+@[A-Za-z0-9-]+(\.[A-Za-z0-9-]+)*\.[A-Za-z]{2,}/[E-Mail entfernt]/g' \
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

# Leitet stdout/stderr des Skripts in $ROH/<datei> (geschwaerzt, Stufe 1) und auf die Konsole.
rohdaten() {
  exec > >(redact | tee "$ROH/$1") 2>&1
  trap 'exec >&- 2>&-; sleep 3' EXIT   # laesst den Filter die letzten Zeilen schreiben
}

now() { date -u +%Y-%m-%dT%H:%M:%S.%3NZ; }
mark() { echo "---- $(now) $*"; }
# Differenz zweier `date +%s.%N`-Werte in Sekunden mit einer Nachkommastelle (Ganzzahl-Arithmetik der Shell).
sekunden() { local a=${1/./} b=${2/./}; local d=$(( (b - a) / 100000000 )); echo "$((d / 10)).$((d % 10))"; }
# Eine stabile Aussage, die ein zweiter Lauf wiederholen muss (siehe README, Abschnitt Lesehilfen).
aussage() { echo "AUSSAGE $*"; }

# Sicherung vor dem ersten Lauf (DoD 1): der Lauf bricht ab, wenn das Artefakt von 00-vorlauf.sh fehlt oder zu einem anderen Stand von
# lib.sh / redact-probe.sh gehoert (Hash im Artefakt). Artefakt: $ROH/rohdaten-00-vorlauf.txt, Zeile "AUSSAGE sicherung-stand=<sha256>".
stand_hash() { docker run --rm -v "$here":/s:ro "$BASE" sh -c 'cat /s/lib.sh /s/redact-probe.sh | sha256sum | cut -d" " -f1'; }
vorlauf_pruefen() {
  local art=$ROH/rohdaten-00-vorlauf.txt h
  [ -s "$art" ] || { echo "Abbruch: Artefakt der Gegenprobe fehlt ($art); zuerst 00-vorlauf.sh"; return 2; }
  h=$(stand_hash)
  grep -q "^AUSSAGE sicherung-stand=$h\$" "$art" || { echo "Abbruch: Artefakt gehoert zu einem anderen Stand von lib.sh / redact-probe.sh; 00-vorlauf.sh erneut"; return 2; }
  grep -q '^AUSSAGE redact-probe=ok$' "$art" || { echo "Abbruch: Artefakt meldet keine bestandene Gegenprobe"; return 2; }
  grep -q '^AUSSAGE umgebung-pruefung=ok$' "$art" || { echo "Abbruch: Artefakt meldet keine saubere Umgebung"; return 2; }
  umgebung_pruefen || { echo "Abbruch: verbotene Namen in der Umgebung"; return 2; }
  return 0
}

# Prompt-Zaehler: jeder Aufruf, der einen Prompt an ein Cloud-Modell senden KANN (auch wenn er vorher an einer Meldung scheitert), zaehlt.
# Bricht mit Exit 99 ab, wenn die Obergrenze erreicht ist; schreibt die Zeile ins Zaehler-Verzeichnis ausserhalb des Repos.
prompt_zaehlen() { # <fall>
  local n; n=$( [ -f "$ZAEHLER" ] && wc -l < "$ZAEHLER" || echo 0 )
  if [ "$n" -ge "$PROMPT_MAX" ]; then echo "Abbruch: Obergrenze von $PROMPT_MAX Prompts erreicht ($n)"; exit 99; fi
  echo "$((n + 1)) $(now) $1" >> "$ZAEHLER"
  echo "PROMPT $((n + 1))/$PROMPT_MAX $1"
}
prompt_stand() { echo "AUSSAGE prompts-gestellt=$( [ -f "$ZAEHLER" ] && wc -l < "$ZAEHLER" || echo 0 ) obergrenze=$PROMPT_MAX"; }

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
# Mit KEEP_IMAGES=1 bleiben die Images (der Folgeschritt braucht das Image des Projekts).
praefix_entfernen() {
  local pat="^(vsc-)?$PFX" n
  docker ps -a --format '{{.Names}}' | while read -r n; do [[ $n =~ $pat ]] && { echo "entferne container $n"; docker rm -f -v "$n" > /dev/null; }; done
  docker volume ls -q | while read -r n; do [[ $n =~ $pat ]] && { echo "entferne volume $n"; docker volume rm "$n" > /dev/null; }; done
  docker network ls --format '{{.Name}}' | while read -r n; do [[ $n =~ $pat ]] && { echo "entferne netz $n"; docker network rm "$n" > /dev/null; }; done
  [ "${KEEP_IMAGES:-0}" = 1 ] || docker images --format '{{.Repository}}:{{.Tag}}' | while read -r n; do [[ $n =~ $pat ]] && { echo "entferne image $n"; docker rmi "$n" > /dev/null; }; done
  return 0
}

# Image des Projekts res (Name, den die Dev-Container-CLI vergibt; der Container der CLI nutzt das Image mit der Endung -uid).
bild() { docker images --format '{{.Repository}}:{{.Tag}}' | grep -E "^vsc-$PFX-res-.*-uid:" | head -1; }

# Variante -> Projektverzeichnis, Projektname. Varianten: std (Standardprofil, ohne Ollama im Container), res (Sandbox mit Egress), alw (wie res
# plus Eintrag host.docker.internal in egress.allow). Alle drei tragen die Handvariante --add-host=host.docker.internal:host-gateway in runArgs.
vdir() { echo "$WORK/$PFX-$1"; }
vname() { echo "$PFX-$1"; }
