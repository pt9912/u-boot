#!/usr/bin/env bash
# messen.sh — Messung des Slice slice-v1-shell-erkennung-fehlender-benutzer-eng (LH-FA-DEV-010).
#
# FRAGE: Wo entsteht die Zeile `unable to find user <Benutzer>: no matching entries in passwd file`
# bei `docker exec -u <Benutzer> [-i] [-t] <Container> <Shell>`, auf welchem Strom (Standardausgabe
# oder Fehlerausgabe), mit welchem Exit-Status, mit und ohne Terminal, bei einem laufenden und bei
# einem gestoppten Container? Und: wie sieht dieselbe Zeile aus, wenn der Benutzer existiert, die
# Sitzung die Zeile aber selbst erzeugt (`docker exec -u <fehlender Benutzer> <anderer Container>`
# innerhalb der Sitzung)?
#
# AUFRUF (aus der Repo-Wurzel; nur Docker-Client und lokale Engine, kein Netz nötig, wenn die
# Images vorhanden sind):
#   bash tools/messung/shell-benutzer-fehlt/messen.sh \
#     > docs/messungen/slice-v1-shell-erkennung-fehlender-benutzer-eng/rohdaten-lauf<L>.txt 2>&1
#
# Umgebungsvariablen: WIEDERHOLUNGEN (Vorgabe 2: jede Zelle läuft so oft in einem Skriptlauf),
# IMG_SESSION (Vorgabe mcr.microsoft.com/devcontainers/base:debian; muss die Binary-Bibliotheken
# des Hosts-`docker` tragen, hier glibc), IMG_OTHER (Vorgabe alpine:3).
#
# Lockerung (ausgewiesen): der Sitzungs-Container `ubmess-shell-sitzung` bekommt den Docker-Socket
# des Hosts und das `docker`-Binary des Hosts (nur lesend), damit die Sitzung wirklich
# `docker exec` gegen den anderen Container ausführen kann. Das Skript ändert sonst nichts am Host.
# Alles, was es anlegt, trägt das Präfix ubmess-shell- und das Label ubmess-shell=1; aufgeräumt
# wird per trap EXIT/INT/TERM nach Label, nie über ein Glob auf fremde Namen.
#
# Zeitlimits: jeder Docker-Aufruf läuft unter `timeout -k 5 <Sekunden>`; ein hängender Aufruf wird
# als Exit-Status 124 (Zeitlimit) oder 137 (kill) in der Rohdatei sichtbar, das Skript bricht nicht ab.
set -u

WIEDERHOLUNGEN="${WIEDERHOLUNGEN:-2}"
IMG_SESSION="${IMG_SESSION:-mcr.microsoft.com/devcontainers/base:debian}"
IMG_OTHER="${IMG_OTHER:-alpine:3}"
LIMIT=30
C_SESSION=ubmess-shell-sitzung  # läuft; Benutzer root vorhanden, Benutzer nouser fehlt; Docker-Zugriff
C_OTHER=ubmess-shell-anderer    # läuft; Ziel des `docker exec` innerhalb der Sitzung
C_STOPPED=ubmess-shell-gestoppt # gestoppt
USERNAME=nouser                 # existiert in keinem der Container
TMP="$(mktemp -d)"

aufraeumen() {
  local ids
  ids="$(docker ps -aq --filter label=ubmess-shell=1 2>/dev/null)"
  if [ -n "$ids" ]; then
    # shellcheck disable=SC2086
    docker rm -f $ids >/dev/null 2>&1
  fi
  rm -rf "$TMP"
}
trap aufraeumen EXIT INT TERM

D() { timeout -k 5 "$LIMIT" docker "$@"; }

kopf() {
  echo "# Messung shell-benutzer-fehlt — $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "# Host: $(uname -srm)"
  echo "# Client: $(docker version --format '{{.Client.Version}} (API {{.Client.APIVersion}})' 2>&1)"
  echo "# Engine: $(docker version --format '{{.Server.Version}} (API {{.Server.APIVersion}}, {{.Server.Os}}/{{.Server.Arch}})' 2>&1)"
  echo "# Images: $IMG_SESSION ; $IMG_OTHER"
  echo "# Wiederholungen je Zelle: $WIEDERHOLUNGEN"
  echo
}

# zelle <Kennung> <Beschreibung> <Befehl als eine Zeichenkette für bash -c>
# Standardeingabe der Zelle: $STDIN_MODE = null (/dev/null) | pipe (Text aus $STDIN_TEXT)
# Terminal-Zellen laufen unter `script` (ein Pseudo-Terminal; Standard- und Fehlerausgabe sind dort
# ein Strom und nicht trennbar).
zelle() {
  local id="$1" beschr="$2" befehl="$3" n rc
  for n in $(seq 1 "$WIEDERHOLUNGEN"); do
    echo "### Zelle $id, Wiederholung $n: $beschr"
    echo "befehl: $befehl"
    if [ "${STDIN_MODE:-null}" = pipe ]; then
      printf '%b' "${STDIN_TEXT:-}" | timeout -k 5 "$LIMIT" bash -c "$befehl" >"$TMP/out" 2>"$TMP/err"
    else
      timeout -k 5 "$LIMIT" bash -c "$befehl" </dev/null >"$TMP/out" 2>"$TMP/err"
    fi
    rc=$?
    echo "exit-status: $rc"
    echo "stdout-bytes: $(wc -c <"$TMP/out")  stderr-bytes: $(wc -c <"$TMP/err")"
    echo "--- stdout"
    cat "$TMP/out"
    echo "--- stderr"
    cat "$TMP/err"
    echo
  done
}

# Vorbereitung: Container anlegen (laufend: Sitzung, anderer; gestoppt: gestoppt).
D rm -f "$C_SESSION" "$C_OTHER" "$C_STOPPED" >/dev/null 2>&1
kopf
echo "## Vorbereitung"
D run -d --name "$C_OTHER" --label ubmess-shell=1 "$IMG_OTHER" sleep 600 >/dev/null || echo "Vorbereitung: anderer Container fehlgeschlagen"
D run -d --name "$C_SESSION" --label ubmess-shell=1 \
  -v /var/run/docker.sock:/var/run/docker.sock -v "$(command -v docker)":/usr/local/bin/docker:ro \
  "$IMG_SESSION" sleep 600 >/dev/null || echo "Vorbereitung: Sitzungs-Container fehlgeschlagen"
D run -d --name "$C_STOPPED" --label ubmess-shell=1 "$IMG_SESSION" sleep 600 >/dev/null || echo "Vorbereitung: Container fehlgeschlagen"
D stop -t 1 "$C_STOPPED" >/dev/null
echo "Container: $(docker ps -a --filter label=ubmess-shell=1 --format '{{.Names}}={{.Status}}' | sort | tr '\n' ' ')"
# Der Benutzer vscode (nur im Sitzungs-Container, nicht im anderen) bekommt Zugriff auf den Socket,
# damit E5 und E6 mit gleichem Benutzernamen in Sitzung und innerem Aufruf laufen können.
D exec -u root "$C_SESSION" bash -c '
  g=$(stat -c %g /var/run/docker.sock)
  id -G vscode | tr " " "\n" | grep -qx "$g" && exit 0
  getent group "$g" >/dev/null || groupadd -g "$g" ubsock
  usermod -aG "$(getent group "$g" | cut -d: -f1)" vscode' >/dev/null 2>&1
echo "Kontrolle: Benutzer vscode: im Sitzungs-Container $(D exec "$C_SESSION" getent passwd vscode >/dev/null 2>&1 && echo vorhanden || echo fehlt), im anderen Container $(D exec "$C_OTHER" getent passwd vscode >/dev/null 2>&1 && echo vorhanden || echo fehlt); docker-Zugriff von vscode in der Sitzung: $(D exec -u vscode "$C_SESSION" docker version --format '{{.Client.Version}}/{{.Server.Version}}' 2>&1)"
echo "Kontrolle: Benutzer $USERNAME im Sitzungs-Container:$(D exec "$C_SESSION" getent passwd "$USERNAME" >/dev/null 2>&1 && echo vorhanden || echo fehlt)"
echo "Kontrolle: docker-Client in der Sitzung: $(D exec "$C_SESSION" docker version --format '{{.Client.Version}}' 2>&1)"
echo

echo "## A. Der Client, Benutzer fehlt, Container läuft (Probe-Form und Sitzungs-Form; kein Terminal)"
zelle A1 "Probe-Form: -u, -w, ohne -i und -t, bash -c true; Standardeingabe /dev/null" \
  "docker exec -u $USERNAME -w / $C_SESSION bash -c true"
zelle A2 "Sitzungs-Form: -i, ohne -t, bash; Standardeingabe /dev/null" \
  "docker exec -i -u $USERNAME -w / $C_SESSION bash"
STDIN_MODE=pipe STDIN_TEXT='true\nexit\n' \
  zelle A3 "Sitzungs-Form: -i, ohne -t, bash; Standardeingabe eine Pipe" \
  "docker exec -i -u $USERNAME -w / $C_SESSION bash"
zelle A4 "ohne -i und -t, bash; Standardeingabe /dev/null" \
  "docker exec -u $USERNAME -w / $C_SESSION bash"

echo "## B. Der Client, Benutzer fehlt, mit Terminal (Pseudo-Terminal über script; ein Strom)"
zelle B1 "-i -t, bash, Terminal" \
  "script -qefc 'docker exec -i -t -u $USERNAME -w / $C_SESSION bash' /dev/null"
zelle B2 "-t ohne -i, bash, Terminal" \
  "script -qefc 'docker exec -t -u $USERNAME -w / $C_SESSION bash' /dev/null"

zelle B3 "-i -t, bash, Terminal; Fehlerausgabe des Clients getrennt (Dateideskriptor 3 = stderr der Zelle)" \
  "script -qefc 'docker exec -i -t -u $USERNAME -w / $C_SESSION bash 2>&3' /dev/null 3>&2"

echo "## C. Der Client, Benutzer fehlt, Container läuft nicht"
zelle C1 "Probe-Form, gestoppter Container, Benutzer fehlt" \
  "docker exec -u $USERNAME -w / $C_STOPPED bash -c true"
zelle C2 "Sitzungs-Form, gestoppter Container, Benutzer fehlt" \
  "docker exec -i -u $USERNAME -w / $C_STOPPED bash"
zelle C3 "Probe-Form, gestoppter Container, Benutzer root (vorhanden)" \
  "docker exec -u root -w / $C_STOPPED bash -c true"

echo "## D. Kontrolle: Benutzer vorhanden, laufender Container"
zelle D1 "Probe-Form, Benutzer root" \
  "docker exec -u root -w / $C_SESSION bash -c true"
STDIN_MODE=pipe STDIN_TEXT='true\nexit 3\n' \
  zelle D2 "Sitzungs-Form, Benutzer root, Shell endet mit Status 3" \
  "docker exec -i -u root -w / $C_SESSION bash"

echo "## E. Die Sitzung erzeugt die Zeile selbst (Benutzer der Sitzung vorhanden; im Container ein"
echo "##    echter docker exec -u $USERNAME gegen den anderen Container, über den Socket des Hosts)"
STDIN_MODE=pipe STDIN_TEXT="docker exec -u $USERNAME $C_OTHER sh -c true\n" \
  zelle E1 "Sitzung ohne Terminal: Skript auf der Standardeingabe der Shell, innerer Aufruf ohne -i" \
  "docker exec -i -u root -w / $C_SESSION bash"
STDIN_MODE=pipe STDIN_TEXT="docker exec -i -u $USERNAME $C_OTHER sh\n" \
  zelle E2 "Sitzung ohne Terminal: innerer Aufruf mit -i" \
  "docker exec -i -u root -w / $C_SESSION bash"
STDIN_MODE=pipe STDIN_TEXT="docker exec -u $USERNAME $C_OTHER sh -c true\nexit\n" \
  zelle E3 "Sitzung mit Terminal (script): innerer Aufruf ohne -i" \
  "script -qefc 'docker exec -i -t -u root -w / $C_SESSION bash' /dev/null"
STDIN_MODE=pipe STDIN_TEXT="docker exec -u root $C_OTHER sh -c true\n" \
  zelle E4 "Kontrolle: Sitzung ohne Terminal, innerer Aufruf mit vorhandenem Benutzer root" \
  "docker exec -i -u root -w / $C_SESSION bash"

STDIN_MODE=pipe STDIN_TEXT="docker exec -u $USERNAME $C_OTHER sh -c true\nexit\n" \
  zelle E3b "Sitzung mit Terminal (script): wie E3, Fehlerausgabe des äußeren Clients getrennt" \
  "script -qefc 'docker exec -i -t -u root -w / $C_SESSION bash 2>&3' /dev/null 3>&2"
echo "## E5/E6: derselbe Benutzername in Sitzung und innerem Aufruf (vscode: im Sitzungs-Container"
echo "##    vorhanden, im anderen Container fehlend) — die Sitzung läuft als der Benutzer der Sitzung"
STDIN_MODE=pipe STDIN_TEXT="docker exec -u vscode $C_OTHER sh -c true\n" \
  zelle E5 "Sitzung ohne Terminal als vscode, innerer Aufruf -u vscode gegen den anderen Container" \
  "docker exec -i -u vscode -w / $C_SESSION bash"
STDIN_MODE=pipe STDIN_TEXT="docker exec -u vscode $C_OTHER sh -c true\nexit 0\n" \
  zelle E6 "wie E5, die Shell endet danach ausdrücklich mit exit 0" \
  "docker exec -i -u vscode -w / $C_SESSION bash"
STDIN_MODE=pipe STDIN_TEXT="docker exec -u vscode $C_OTHER sh -c true 2>&1\n" \
  zelle E7 "wie E5, der innere Aufruf mit 2>&1 (Zeile auf der Standardausgabe der Sitzung)" \
  "docker exec -i -u vscode -w / $C_SESSION bash"

echo "## F. Nachgestellte Zeile der Sitzung (synthetisch: die Sitzung schreibt den Text von Hand; das ist"
echo "##    keine Beobachtung des Clients, nur das Muster für die Erkennung)"
zelle F1 "Sitzung schreibt die Zeile auf die Standardausgabe und endet mit 126" \
  "docker exec -i -u root $C_SESSION bash -c 'echo \"unable to find user $USERNAME: no matching entries in passwd file\"; exit 126'"

echo "## G. Die erneute Probe nach der Sitzung: dieselbe Probe-Form wie A1/D1, einzeln"
zelle G1 "Probe des Benutzers der Sitzung (root, vorhanden) nach einer Sitzung" \
  "docker exec -u root -w / $C_SESSION bash -c true 2>&1; echo \"gemeinsamer Strom oben\""
zelle G2 "Probe des fehlenden Benutzers, beide Ströme zusammen (CombinedOutput)" \
  "docker exec -u $USERNAME -w / $C_SESSION bash -c true 2>&1"
echo "## Ende $(date -u +%Y-%m-%dT%H:%M:%SZ)"
