#!/bin/bash
# Laeuft IM Container als Benutzer vscode (Image aus `devcontainer up` des erzeugten Projekts); aufgerufen von 04-lauf.sh.
# Aufruf im Container: inner.sh <modus> <cloud-modell>
#   discover    keine Restriktion (Container ohne postStartCommand); ein protokollierender dnsmasq reicht alle Namen weiter
#   restricted  die Restriktion des erzeugten Startscripts lief als postStartCommand; ein protokollierender dnsmasq mit derselben
#               Konfiguration ersetzt den Resolver des Skripts (Verfahren O4 des Vorgaenger-Protokolls)
# Kein Konto, kein Token, keine Anmeldung. Jeder Aufruf laeuft hoechstens einmal je Fall, mit timeout und ohne Eingabe (</dev/null);
# keine Wiederholungsschleife. Dateien werden nur als Pfad und Groesse ausgewiesen, nie ihr Inhalt.
set -u
mode=$1; M=$2
log=/tmp/dns.log
now() { date -u +%T.%N | cut -c1-12; }
mark() { echo "---- $(now) $*"; }
dauer() { awk -v a="$1" -v b="$2" 'BEGIN{printf "%.1f", b-a}'; }
# Ausgabe ohne Steuerzeichen und Dubletten, hoechstens 60 Zeilen
strip() { sed -e 's/\x1b\[[0-9;?]*[a-zA-Z]//g' -e 's/\x1b[()][A-Z0-9]//g' | tr '\r' '\n' | grep -v -E '^\s*$' | awk '!seen[$0]++' | head -60; }
dns_seit() { tail -n +$(($1+1)) $log | grep -E 'query\[(A|AAAA)\]' | awk '{print $6}' | sort | uniq -c; }
# Dateien, die seit der letzten Marke angelegt oder veraendert wurden: Pfad, Groesse, Einhaengepunkt (Volume ja/nein); kein Inhalt
neu() {
  echo "-- Dateien neu/geaendert seit der letzten Marke ($1): Pfad, Groesse, Einhaengepunkt"
  { find "$HOME" /etc /var /usr/local /opt /root /tmp -type f -newer /tmp/m.prev 2> /dev/null; } \
    | grep -v -E '^/tmp/(m\.|dns|ollama\.log|o-)|^/var/(lib/apt|cache|log)|^/var/lib/dpkg|^/var/lib/u-boot' | sort | while read -r f; do
      [ -f "$f" ] || { echo "$f (beim Auflisten schon entfernt)"; continue; }
      printf '%s %s %s\n' "$f" "$(stat -c %s "$f")" "$(df --output=target "$f" 2> /dev/null | tail -1)"
    done | head -40
  touch /tmp/m.prev
}
# <label> <timeout-s> <befehl...>: einmal, ohne Eingabe, mit Dauer, Exit-Code, geschwaerzter Ausgabe und den Namen im Resolver-Log
fall() {
  local label=$1 t=$2; shift 2
  mark "$label: $*"
  local n s e rc; n=$(wc -l < $log); s=$(date +%s.%N)
  timeout -k 5 "$t" "$@" < /dev/null > /tmp/o-$$.txt 2>&1; rc=$?
  e=$(date +%s.%N)
  strip < /tmp/o-$$.txt
  echo "$label exit=$rc dauer=$(dauer "$s" "$e")s (124 = timeout nach ${t}s)"
  echo "AUSSAGE $label-exit=$rc"
  echo "-- Namen im Resolver-Log waehrend $label (Anzahl Anfragen A/AAAA)"
  dns_seit "$n"
  if [ "$mode" = restricted ]; then
    echo "-- abgelehnte/unbekannte Namen waehrend $label"
    tail -n +$((n+1)) $log | grep -E 'REFUSED|NXDOMAIN|NODATA|config ' | awk '{print $4, $5, $6}' | sort | uniq -c | head
  fi
  neu "$label"
}
serve() { # [env-zuweisung...]
  nohup env "$@" ollama serve > /tmp/ollama.log 2>&1 &
  spid=$!
  for i in $(seq 1 60); do curl -fsS -m 2 http://127.0.0.1:11434/ > /dev/null 2>&1 && break; sleep 0.5; done
  echo "serve bereit: $(curl -sS -m 2 http://127.0.0.1:11434/api/version)"
}

sudo touch "$log"; sudo chmod 666 "$log"
up=$(awk '/^nameserver/ && $2 !~ /:/ {print $2}' /etc/resolv.conf | head -1)
if [ "$mode" = discover ]; then
  printf 'listen-address=127.0.0.1\nbind-interfaces\nno-resolv\nno-hosts\nserver=%s\nlog-queries\nlog-facility=%s\n' "$up" "$log" | sudo tee /tmp/dns-discover.conf > /dev/null
  sudo dnsmasq --conf-file=/tmp/dns-discover.conf
  echo "nameserver 127.0.0.1" | sudo tee /etc/resolv.conf > /dev/null
else
  echo "-- Status-Datei des Startscripts (postStartCommand)"; sudo cat /var/lib/u-boot/status/egress
  grep -E '^(server|nftset)' /etc/u-boot-egress-dnsmasq.conf | sed 's/^/conf: /'
  sudo kill "$(cat /run/u-boot-egress-dnsmasq.pid)"; sleep 1
  { cat /etc/u-boot-egress-dnsmasq.conf; echo "log-queries"; echo "log-facility=$log"; } | sudo tee /tmp/dns-restricted.conf > /dev/null
  sudo dnsmasq --conf-file=/tmp/dns-restricted.conf
fi
sleep 1
touch /tmp/m.prev
echo "modus=$mode modell=$M benutzer=$(id -un) HOME=$HOME"
echo "-- Namen der Umgebungsvariablen im Container (Werte nicht ausgegeben)"
env | cut -d= -f1 | sort | tr '\n' ' ' | fold -w 150; echo
echo "AUSSAGE umgebung-container-verbotene-namen=$(env | cut -d= -f1 | grep -c -E 'ANTHROPIC|TOKEN|API_KEY|SECRET|OLLAMA')"
echo "-- Einhaengepunkt von ~/.ollama (Volume?)"
df --output=target,fstype "$HOME/.ollama" | tail -1
df --output=target "$HOME" | tail -1

mark "C0 Version und Unterbefehle"
n=$(wc -l < $log)
ollama --version 2>&1
echo "AUSSAGE ollama-$(ollama --version 2>&1 | grep -o 'version is .*')"
echo "claude im PATH: $(command -v claude || echo nein)"
echo "claude --version: $(timeout 20 claude --version < /dev/null 2>&1 | head -2)"
echo "-- ollama --help (Unterbefehle)"
ollama --help 2>&1 | strip
echo "AUSSAGE c0-launch-in-help=$(ollama --help 2>&1 | grep -c -w launch)"
echo "-- ollama launch --help (ohne Server aufrufbar?)"
ollama launch --help 2>&1 | strip
echo "AUSSAGE c0-launch-claude-in-help=$(ollama launch --help 2>&1 | grep -c -i claude)"
echo "-- Namen der Variablen OLLAMA_* im Binary (nur Namen, gelesen mit grep -a aus dem installierten Binary)"
echo "(Go-Zeichenketten liegen im Binary aneinandergereiht; Namen koennen mit Nachbarn verklebt erscheinen. Ausgewiesen: Namen mit KEY, TOKEN, API, AUTH, CLOUD, SIGN, LAUNCH)"
grep -a -o -E 'OLLAMA_[A-Z0-9_]*(KEY|TOKEN|API|AUTH|CLOUD|SIGN|LAUNCH)[A-Z0-9_]*' "$(command -v ollama)" | sort -u | tr '\n' ' ' | fold -w 150; echo
echo "AUSSAGE c4-ollama-var-mit-key-api-token-im-binary=$(grep -a -o -E 'OLLAMA_[A-Z0-9_]*(KEY|TOKEN|API|AUTH)[A-Z0-9_]*' "$(command -v ollama)" | sort -u | tr '\n' ' ')"
echo "AUSSAGE c5-no-cloud-im-binary=$(grep -a -c 'OLLAMA_NO_CLOUD' "$(command -v ollama)")"
echo "-- Namen im Resolver-Log fuer C0"; dns_seit "$n"
neu C0

mark "serve starten (ohne OLLAMA_NO_CLOUD)"
n=$(wc -l < $log)
serve
echo "-- Namen beim Start des Servers"; dns_seit "$n"
echo "-- ollama help serve (Umgebungsvariablen laut Hilfe)"
ollama help serve 2>&1 | strip
neu serve

fall C1-pull 90 ollama pull "$M"
fall C1-run 90 ollama run "$M" "Reply with the single word OK."
echo "-- ollama list / ps nach den Aufrufen"; ollama list 2>&1 | head -5; ollama ps 2>&1 | head -5
fall C3-launch-claude 40 ollama launch claude --model "$M"
mark "C3 mit Pseudo-Terminal (script), Aufruf bis zum Hindernis"
n=$(wc -l < $log); s=$(date +%s.%N)
script -qec "timeout -k 5 25 ollama launch claude --model $M" /dev/null < /dev/null > /tmp/o-tty.txt 2>&1; rc=$?
echo "-- Prozesse von ollama launch nach dem Aufruf: $(pgrep -c -f 'ollama launch' || true)"; pkill -9 -f 'ollama launch' 2> /dev/null
echo "ausgabe des Aufrufs im Pseudo-Terminal: $(wc -c < /tmp/o-tty.txt) Bytes, $(strip < /tmp/o-tty.txt | wc -l) Zeilen nach Entfernen der Steuerzeichen"
strip < /tmp/o-tty.txt
echo "C3-launch-claude-tty exit=$rc dauer=$(dauer "$s" "$(date +%s.%N)")s"
echo "AUSSAGE C3-launch-claude-tty-exit=$rc"
echo "-- Namen im Resolver-Log (C3 tty)"; dns_seit "$n"
neu C3-tty
if [ "$mode" = restricted ]; then
  # Handvariante: claude durch einen Platzhalter ersetzt, nur fuer diesen einen Aufruf (PATH); zeigt, was `ollama launch claude` uebergibt
  mkdir -p /tmp/stub; cp /tmp/claude-stub.sh /tmp/stub/claude; chmod +x /tmp/stub/claude
  fall C3-launch-claude-stub 40 env PATH="/tmp/stub:$PATH" ollama launch claude --model "$M"
fi
fall C4-signin 40 ollama signin
fall C4-signout 20 ollama signout

mark "serve neu mit OLLAMA_NO_CLOUD=1 (C5)"
kill $spid 2> /dev/null; sleep 2
n=$(wc -l < $log)
serve OLLAMA_NO_CLOUD=1
echo "-- Namen beim Start des Servers (NO_CLOUD)"; dns_seit "$n"
export OLLAMA_NO_CLOUD=1
fall C5-pull 90 ollama pull "$M"
fall C5-run 90 ollama run "$M" "Reply with the single word OK."
fall C5-launch-claude 40 ollama launch claude --model "$M"
fall C5-signin 40 ollama signin
unset OLLAMA_NO_CLOUD

mark "alle DNS-Anfragen der Sitzung (eindeutig, A/AAAA)"
grep -E 'query\[(A|AAAA)\]' $log | awk '{print $6}' | sort | uniq -c
echo "-- Server-Log: Zeilen mit cloud (ohne Zugangsdaten, hoechstens 12)"
grep -i -E 'cloud|sign|auth|unauthor' /tmp/ollama.log | cut -c1-220 | head -12
echo "-- ~/.ollama (Pfad, Groesse; kein Inhalt)"
find "$HOME/.ollama" -maxdepth 2 -printf '%p %s\n' | sort | head -20
kill $spid 2> /dev/null; sleep 1
