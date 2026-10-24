#!/bin/bash
# Weg (b): Ollama-Server IM Container (Sandbox-Profil, Image des Projekts res), angemeldet vom Projektinhaber mit 07-anmeldung.sh an dessen Terminal.
# Ein Container zur Zeit. Das Modell-Volume ($PFX-b-modelle auf /home/vscode/.ollama) ueberlebt die Neuanlage des Containers (K5b).
#
# Aufruf: WORK=<messordner> [LAUF=<n>] tools/messung/ollama-cloud-konto/06-lauf-container.sh <stufe> ...
#   vorbereiten                    Container $PFX-b aus dem Image von res (docker run, ohne postStartCommand: Modus discover), Resolver mit Protokoll,
#                                  ollama serve, Marke fuer K4. Kein Prompt, keine Anmeldung. (07-anmeldung.sh ruft diese Stufe selbst auf.)
#   messen discover <fall>...      Faelle im Container $PFX-b (Modus discover)
#   messen restricted <extras> <fall>...
#                                  Container $PFX-b wird entfernt und als $PFX-b2 NEU angelegt (gleiches Volume, K5b); das Egress-Skript laeuft wie
#                                  das postStartCommand (Handvariante), der Resolver mit Protokoll uebernimmt; <extras> = "-" oder "host:<name>;host:<name>"
#                                  (je Name die zwei Zeilen, die das Egress-Skript fuer einen erlaubten Namen erzeugt; Handvariante)
#   ende                           ollama signout im Container, Server beenden, Container und Volume entfernen (Rueckabwicklung; den Schluessel in den
#                                  Kontoeinstellungen entfernt der Projektinhaber)
#   Faelle: daten                  K4: Namen im Resolver-Log seit der Marke vor der Anmeldung, Dateien neu/geaendert, ~/.ollama (Pfad/Groesse/Rechte), Einhaengepunkte
#           run:1|2|o              ollama run (1 = glm-5.3-flash:cloud[1m], 2 = gpt-oss:120b-cloud, o = glm-5.3-flash:cloud ohne [1m])  (1 Prompt)
#           claude:1|2|o           claude -p direkt gegen http://127.0.0.1:11434                                                       (1 Prompt)
#           launch:1|2             ollama launch claude --model <m> --yes, Prompt ueber stdin                                         (1 Prompt)
#           nocloud                Server und Client mit OLLAMA_NO_CLOUD=1: run:1 und launch:1                                        (2 Prompts)
#           apikey                 Server und Client mit OLLAMA_API_KEY=platzhalter (Platzhalter, kein gueltiger Schluessel): run:1, launch:1 (2 Prompts)
#           version                ollama --version und api/version im Container
# Voraussetzung: 00-vorlauf.sh, 02-projekte.sh, 03-lauf-vorbereitung.sh und ein gebautes Image von res (05-lauf-host.sh bauen).
# Rohdaten (Stufe 1): $ROH/rohdaten-06-<stufe>-<lauf>.txt. Zaehler wie 05-lauf-host.sh.
set -u
. "$(dirname "$0")/lib.sh"
vorlauf_pruefen || exit 2
LAUF=${LAUF:-1}
stufe=${1:-}; shift
B=$PFX-b; B2=$PFX-b2; VOL=$PFX-b-modelle
CLAUDE_ENV="ANTHROPIC_BASE_URL=http://127.0.0.1:$PORT ANTHROPIC_AUTH_TOKEN=ollama"
modell() { case "$1" in 1) echo "$M1" ;; 2) echo "$M2" ;; o) echo "${M1%\[1m\]}" ;; esac; }
ex() { docker exec -u vscode -w /home/vscode "$c" bash /tmp/inner.sh "$@"; }
neu_container() { # <name>: aus dem Image von res, mit dem Modell-Volume
  local n=$1 img; img=$(bild); [ -n "$img" ] || { echo "kein Image von res (05-lauf-host.sh bauen)"; exit 1; }
  docker volume create "$VOL" > /dev/null
  docker rm -f -v "$n" > /dev/null 2>&1
  docker run -d --name "$n" --cap-add=NET_ADMIN -v "$VOL":/home/vscode/.ollama --entrypoint sleep "$img" 86400 > /dev/null || exit 1
  docker cp "$here/inner.sh" "$n":/tmp/inner.sh
  echo "container $n aus $img, Volume $VOL auf /home/vscode/.ollama"
}

case "$stufe" in
vorbereiten)
  rohdaten "rohdaten-06-vorbereiten-$LAUF.txt"
  echo "== Weg (b) vorbereiten, $(now)"
  c=$B; neu_container "$B"
  ex setup discover
  ex serve
  echo "-- Version des Servers und des Clients im Container"
  ex fall "B-version" 25 0 0 - ollama --version
  docker exec "$c" sh -c 'wc -l < /tmp/dns.log > /tmp/sign.mark'
  echo "-- Zustand vor der Anmeldung: Dateien (der erste Start des Servers legt den Geraeteschluessel an)"
  ex dateien vor-anmeldung
  ex ablage
  echo "Marke fuer K4 gesetzt (Zeile im Resolver-Log: $(docker exec "$c" cat /tmp/sign.mark))"
  prompt_stand
  ;;
messen)
  modus=${1:-}; shift
  case "$modus" in
  discover) c=$B ;;
  restricted)
    extras=${1:--}; shift
    [ "$extras" = - ] && extras=""
    ;;
  *) echo "Aufruf: messen discover|restricted ..."; exit 2 ;;
  esac
  rohdaten "rohdaten-06-$modus-$LAUF.txt"
  echo "== Weg (b), Modus $modus, Lauf $LAUF, Faelle: $*, $(now)"
  prompt_stand
  if [ "$modus" = restricted ]; then
    echo "-- Neuanlage des Containers mit erhaltenem Volume (K5b)"
    docker rm -f -v "$B" > /dev/null 2>&1
    c=$B2; neu_container "$B2"
    ex setup restricted "$extras"
    ex serve
  else
    docker exec "$c" true || { echo "Container $c laeuft nicht (zuerst vorbereiten und anmelden)"; exit 1; }
  fi
  ex nft
  for f in "$@"; do
    k=${f#*:}; m=$(modell "$k"); lbl="B-$modus-${f%%:*}-m$k"
    case "$f" in
    daten)
      echo "-- K4: Resolver-Namen seit der Marke vor der Anmeldung (Anmeldeablauf, Serverstart)"
      ex dns "$(docker exec "$c" cat /tmp/sign.mark)"
      ex dateien nach-anmeldung
      ex ablage
      ;;
    version) ex fall "B-$modus-version" 25 0 0 - ollama --version; ex fall "B-$modus-api-version" 15 0 0 - curl -sS -m 5 "http://127.0.0.1:$PORT/api/version" ;;
    run:*)    prompt_zaehlen "weg-b $modus $f"; ex fall "$lbl" 90 1 0 - ollama run "$m" "@PROMPT@" ;;
    claude:*) prompt_zaehlen "weg-b $modus $f"; ex fall "$lbl" 120 1 0 "$CLAUDE_ENV" claude --model "$m" -p "@PROMPT@" ;;
    launch:*) prompt_zaehlen "weg-b $modus $f"; ex fall "$lbl" 120 1 1 - ollama launch claude --model "$m" --yes ;;
    nocloud)
      echo "-- K7: Server und Client mit OLLAMA_NO_CLOUD=1"
      ex stopserve; ex serve OLLAMA_NO_CLOUD=1
      prompt_zaehlen "weg-b $modus nocloud run:1";    ex fall "B-$modus-nocloud-run" 90 1 0 OLLAMA_NO_CLOUD=1 ollama run "$M1" "@PROMPT@"
      prompt_zaehlen "weg-b $modus nocloud launch:1"; ex fall "B-$modus-nocloud-launch" 120 1 1 OLLAMA_NO_CLOUD=1 ollama launch claude --model "$M1" --yes
      ex stopserve; ex serve
      ;;
    apikey)
      echo "-- K7: Server und Client mit OLLAMA_API_KEY=platzhalter (Platzhalter, kein gueltiger Schluessel)"
      ex stopserve; ex serve OLLAMA_API_KEY=platzhalter
      prompt_zaehlen "weg-b $modus apikey run:1";    ex fall "B-$modus-apikey-run" 90 1 0 OLLAMA_API_KEY=platzhalter ollama run "$M1" "@PROMPT@"
      prompt_zaehlen "weg-b $modus apikey launch:1"; ex fall "B-$modus-apikey-launch" 120 1 1 OLLAMA_API_KEY=platzhalter ollama launch claude --model "$M1" --yes
      ex stopserve; ex serve
      ;;
    *) echo "unbekannter Fall: $f" ;;
    esac
  done
  prompt_stand
  echo "== Ende, $(now)"
  ;;
ende)
  rohdaten "rohdaten-06-ende-$LAUF.txt"
  echo "== Weg (b) Rueckabwicklung, $(now)"
  c=$(docker ps --format '{{.Names}}' | grep -x -E "$B2|$B" | head -1)
  if [ -n "$c" ]; then
    ex fall "B-signout" 30 0 0 - ollama signout
    ex stopserve
    ex ablage
  else echo "kein Container des Weges (b) laeuft"; fi
  docker rm -f -v "$B" "$B2" > /dev/null 2>&1
  docker volume rm "$VOL" > /dev/null 2>&1
  echo "container/volume nach dem Entfernen: $(bestand | grep -c -E "^(container|volume) $PFX-b")"
  echo "SCHRITT FUER DEN PROJEKTINHABER: den Geraeteschluessel dieses Wegwerf-Containers in den Kontoeinstellungen von ollama.com entfernen."
  prompt_stand
  ;;
*) echo "Aufruf: 06-lauf-container.sh vorbereiten | messen discover|restricted ... | ende"; exit 2 ;;
esac
