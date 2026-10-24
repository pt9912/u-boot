#!/bin/bash
# Weg (a): Claude Code und ollama-Client im Container, angemeldeter Ollama-Server des HOSTS ueber das Netz (host.docker.internal -> Relay -> 127.0.0.1:11434).
# Der Host-Server wird nur aus dem Container angesprochen; kein Host-Werkzeug, keine Datei des Hosts. Ein Container zur Zeit.
#
# Aufruf: WORK=<messordner> [LAUF=<n>] tools/messung/ollama-cloud-konto/05-lauf-host.sh <teil> [<fall>...]
#   teil     std         Standardprofil (Projekt std, Container der Dev-Container-CLI, kein Resolver-Protokoll)
#            discover    Sandbox-Projekt res: zweiter Container aus dem Image der CLI per `docker run --cap-add=NET_ADMIN`, ohne postStartCommand (Handvariante)
#            restricted  Container der CLI aus Projekt res (Liste des Generators; K3a)
#            allow       Container der CLI aus Projekt alw (Liste plus host.docker.internal; K3b)
#            c           wie allow (Projekt alw), dazu die Zusatzzeile address=/host.docker.internal/<Gateway> im Resolver (K3c, Handvariante)
#            bauen       nur `devcontainer up` der drei Projekte (kein Container bleibt), kein Prompt
#   fall     k1          Erreichbarkeit ohne und mit Relay, `ollama --version` mit OLLAMA_HOST (kein Prompt)
#            nft         Zustand der nftables-Menge
#            nftadd      HANDVARIANTE: traegt das Gateway von Hand in die nftables-Menge ein (K3c)
#            curl        api/version des Host-Servers mit Relay (kein Prompt)
#            claude:1|2  claude -p mit Modell 1 (glm-5.3-flash:cloud[1m]) oder 2 (gpt-oss:120b-cloud)       (1 Prompt)
#            run:1|2     ollama run <modell> "<prompt>" mit OLLAMA_HOST des Clients auf den Host-Server       (1 Prompt)
#            launch:1|2  ollama launch claude --model <modell> --yes, Prompt ueber stdin                     (1 Prompt)
#   Ohne <fall>: nur k1.
# Voraussetzung: 00-vorlauf.sh, 02-projekte.sh, 03-lauf-vorbereitung.sh. Rohdaten (Stufe 1): $ROH/rohdaten-05-<teil>-<lauf>-<HHMMSS>.txt, rohdaten-05-up-<variante>-<lauf>-<HHMMSS>.txt
# Der Prompt-Zaehler steht im Messordner ($WORK/prompts.txt); die Obergrenze (40) bricht den Lauf mit Exit 99 ab.
set -u
. "$(dirname "$0")/lib.sh"
vorlauf_pruefen || exit 2
LAUF=${LAUF:-1}
zeit=$(date -u +%H%M%S)   # Teil des Dateinamens: jeder Aufruf schreibt eine eigene Rohdatei
teil=${1:-}; shift
[ -n "$teil" ] || { echo "Aufruf: 05-lauf-host.sh <teil> [<fall>...]"; exit 2; }
faelle=("$@"); [ ${#faelle[@]} -gt 0 ] || faelle=(k1)
HOSTURL="http://host.docker.internal:$PORT"
gateway=172.17.0.1   # Gateway des Netzes bridge (gelesen in 04-relay.sh: host-gateway im Container)

cli() { # <variante> <argumente...>
  local v=$1 P; P=$(vname "$1"); shift
  docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "$(vdir "$v")":"$(wdir "$v")" -w "$(wdir "$v")" "$CLI" devcontainer "$@" --workspace-folder "$(wdir "$v")"
}
entfernen() { # <variante>: Container der CLI und Volumes des Projekts (frisches Volume je Lauf)
  local v=$1 pat="^$(vname "$1")-" n
  docker ps -a -q --filter "label=devcontainer.local_folder=$(wdir "$v")" | while read -r c; do docker rm -f -v "$c" > /dev/null; done
  docker volume ls -q | while read -r n; do [[ $n =~ $pat ]] && docker volume rm "$n" > /dev/null; done
}
hoch() { # <variante>: Container der CLI starten; gibt die Container-ID auf stdout aus, Meldungen auf stderr
  local v=$1 s rc; s=$(date +%s.%N)
  entfernen "$v"
  cli "$v" up --log-level info --remove-existing-container 2>&1 | grep -v -E '^\[[0-9]+ ms\]|^\s*$' | tail -12 | redact > "$ROH/rohdaten-05-up-$v-$LAUF-$zeit.txt"
  rc=${PIPESTATUS[0]}
  echo "devcontainer up $v exit=$rc dauer=$(sekunden "$s" "$(date +%s.%N)")s" >&2
  docker ps -q --filter "label=devcontainer.local_folder=$(wdir "$v")" | head -1
}

if [ "$teil" = bauen ]; then
  rohdaten "rohdaten-05-bauen-$LAUF.txt"
  echo "== bauen (devcontainer up der drei Projekte, kein Prompt), $(now)"
  liste=("${faelle[@]}"); [ "${liste[0]}" = k1 ] && liste=(res alw std)   # ohne Angabe alle drei; sonst die genannten Varianten
  for v in "${liste[@]}"; do c=$(hoch "$v"); echo "container $v: ${c:-keiner}"; aussage "bauen-$v-container=$([ -n "$c" ] && echo ja || echo nein)"; entfernen "$v"; done
  echo "image res: $(bild)"
  exit 0
fi

rohdaten "rohdaten-05-$teil-$LAUF-$zeit.txt"
echo "== Weg (a), Teil $teil, Lauf $LAUF, Faelle: ${faelle[*]}, $(now)"
prompt_stand
extra=""
case "$teil" in
  std)        v=std; mode=none; c=$(hoch std) ;;
  restricted) v=res; mode=restricted; c=$(hoch res) ;;
  allow)      v=alw; mode=restricted; c=$(hoch alw) ;;
  c)          v=alw; mode=restricted; extra="address=/host.docker.internal/$gateway"; c=$(hoch alw) ;;
  discover)
    v=res; mode=discover
    c0=$(hoch res); img=$(docker inspect -f '{{.Config.Image}}' "$c0"); entfernen res
    c=$PFX-discover
    docker rm -f -v "$c" > /dev/null 2>&1
    docker run -d --name "$c" --cap-add=NET_ADMIN --add-host=host.docker.internal:host-gateway --entrypoint sleep "$img" 3000 > /dev/null
    ;;
  *) echo "unbekannter Teil: $teil"; exit 2 ;;
esac
[ -n "${c:-}" ] || { echo "kein Container"; exit 1; }
echo "image des Containers: $(docker inspect -f '{{.Config.Image}}' "$c")"
echo "-- Einhaengepunkte des Containers (Art, Ziel)"
docker inspect -f '{{range .Mounts}}{{.Type}} {{.Destination}}{{"\n"}}{{end}}' "$c"
echo "-- Egress-Status im Container: $(docker exec "$c" cat /var/lib/u-boot/status/egress 2> /dev/null | tr '\n' ' ')"
docker cp "$here/inner.sh" "$c":/tmp/inner.sh
ex() { docker exec -u vscode -w /home/vscode "$c" bash /tmp/inner.sh "$@"; }

relay_an=0
relay() { # start|stop
  local i
  if [ "$1" = start ]; then
    # ein fremder Container kann kurz im Netz haengen (etwa ein Gate-Lauf im Nachbar-Repo): bis zu vier Versuche, dann Abbruch
    for i in 1 2 3 4; do "$here/04-relay.sh" start > /dev/null 2>&1 && break; [ "$i" = 4 ] && { echo "Abbruch: Relay nicht gestartet"; exit 3; }; sleep 15; done
    relay_an=1; echo "relay: gestartet"
  else "$here/04-relay.sh" stop > /dev/null 2>&1; relay_an=0; echo "relay: entfernt (docker ps -a: $(docker ps -a --format '{{.Names}}' | grep -c -x "$RELAY") Relay-Container)"; fi
}
teil_aufraeumen() { if [ "$teil" = discover ]; then docker rm -f -v "$c" > /dev/null 2>&1; else entfernen "$v"; fi; }
trap '[ "$relay_an" = 1 ] && "$here/04-relay.sh" stop > /dev/null 2>&1; teil_aufraeumen; exec >&- 2>&-; sleep 3' EXIT

ex setup "$mode" "$extra"
CLAUDE_ENV="ANTHROPIC_BASE_URL=$HOSTURL ANTHROPIC_AUTH_TOKEN=ollama"
modell() { case "$1" in 1) echo "$M1" ;; 2) echo "$M2" ;; o) echo "${M1%\[1m\]}" ;; esac; }   # o = Modell 1 ohne den Zusatz [1m]
# ein unbekannter Modellschluessel bricht ab, bevor ein Prompt gezaehlt wird
modell_pruefen() { [ -n "$(modell "$1")" ] || { echo "Abbruch: unbekannter Modellschluessel '$1'"; exit 2; }; }
for f in "${faelle[@]}"; do
  case "$f" in
  k1)
    # K1: ohne Hilfe (Relay aus), dann mit Relay (K1b)
    ex fall "K1-getent" 10 0 0 - getent hosts host.docker.internal
    ex fall "K1-api-ohne-relay" 15 0 0 - curl -sS -m 5 "$HOSTURL/api/version"
    relay start
    ex fall "K1b-api-mit-relay" 15 0 0 - curl -sS -m 5 "$HOSTURL/api/version"
    ex fall "K1b-ollama-version" 25 0 0 "OLLAMA_HOST=$HOSTURL" ollama --version
    ;;
  nft) ex nft ;;
  versionen)
    ex fall "V-claude" 30 0 0 - claude --version
    ex fall "V-node" 20 0 0 - node --version
    ex fall "V-ollama-client" 20 0 0 - ollama --version
    ;;
  nftadd) ex nftadd "$gateway"; ex nft ;;
  curl)
    [ "$relay_an" = 1 ] || relay start
    ex fall "K-curl-api" 15 0 0 - curl -sS -m 5 "$HOSTURL/api/version"
    ;;
  claude:*)
    modell_pruefen "${f#*:}"; m=$(modell "${f#*:}"); [ "$relay_an" = 1 ] || relay start
    prompt_zaehlen "weg-a $teil lauf$LAUF claude:${f#*:}"
    ex fall "K-claude-m${f#*:}" 120 1 0 "$CLAUDE_ENV" claude --model "$m" -p "@PROMPT@"
    ;;
  run:*)
    modell_pruefen "${f#*:}"; m=$(modell "${f#*:}"); [ "$relay_an" = 1 ] || relay start
    prompt_zaehlen "weg-a $teil lauf$LAUF run:${f#*:}"
    ex fall "K-run-m${f#*:}" 90 1 0 "OLLAMA_HOST=$HOSTURL" ollama run "$m" "@PROMPT@"
    ;;
  launch:*)
    modell_pruefen "${f#*:}"; m=$(modell "${f#*:}"); [ "$relay_an" = 1 ] || relay start
    prompt_zaehlen "weg-a $teil lauf$LAUF launch:${f#*:}"
    ex fall "K-launch-m${f#*:}" 120 1 1 "OLLAMA_HOST=$HOSTURL" ollama launch claude --model "$m" --yes
    ;;
  *) echo "unbekannter Fall: $f" ;;
  esac
done
[ "$relay_an" = 0 ] || relay stop
echo "-- Aufraeumen des Containers und der Volumes dieses Teils"
teil_aufraeumen
prompt_stand
echo "== Ende, $(now)"
