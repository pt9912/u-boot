#!/bin/bash
# Laeuft IM Container als Benutzer vscode (Image aus `devcontainer up` des erzeugten Projekts); wird von den Treibern per `docker cp` hineinkopiert
# und per `docker exec` aufgerufen. Aufrufe:
#   inner.sh setup <modus> [<dnsmasq-zeile>]   Resolver mit Protokoll einrichten. Modus none: nichts (leeres Log); discover: alle Namen werden weitergereicht;
#                                              restricted: die Konfiguration des Egress-Skripts (laeuft es nicht, wird es wie das
#                                              postStartCommand gestartet) plus Protokoll, optional plus eine zusaetzliche Zeile (K3c)
#   inner.sh fall <label> <timeout-s> <prompt:0|1> <stdin:0|1> <env|-> <befehl...>
#                                              einen Aufruf, einmal, mit timeout -k 5; @PROMPT@ im Befehl steht fuer den Prompt-Satz,
#                                              stdin=1 gibt den Prompt-Satz auf stdin; env = Zuweisungen NAME=wert (durch Leerzeichen
#                                              getrennt) fuer genau diesen Aufruf
#   inner.sh serve [NAME=wert...]              ollama serve im Hintergrund (nur Weg b)
#   inner.sh stopserve                         beendet ihn
#   inner.sh dateien <label>                   Dateien neu/geaendert seit der letzten Marke (Pfad, Groesse, Rechte, Einhaengepunkt; kein Inhalt)
#   inner.sh dns [<ab-zeile>]                  Namen im Resolver-Log (Anzahl A/AAAA), ab Zeile
#   inner.sh nft                               Zustand der nftables-Menge der Egress-Restriktion (nur Elementzahl und ob 172.17.0.1 darin steht)
# Die Ausgabe eines Aufrufs mit prompt=1 traegt KEINEN Antwort-Inhalt: ausgewiesen werden Exit-Code, Dauer, Zeilen- und Byte-Zahl der Ausgabe und
# nur die Zeilen, die nach einer Meldung aussehen (Muster MELDUNG unten). Jeder Aufruf laeuft hoechstens einmal, ohne Wiederholungsschleife.
set -u
log=/tmp/dns.log
now() { date -u +%T.%N | cut -c1-12; }
mark() { echo "---- $(now) $*"; }
dauer() { awk -v a="$1" -v b="$2" 'BEGIN{printf "%.1f", b-a}'; }
PROMPT='Reply with the single word OK.'
MELDUNG='sign|error|warn|fail|unauthori|denied|forbidden|not found|refused|cloud|requires|invalid|expired|limit|quota|rate|timeout|dial|lookup|connect|model'
# Ausgabe ohne Steuerzeichen und Dubletten, hoechstens 60 Zeilen
strip() { sed -e 's/\x1b\[[0-9;?]*[a-zA-Z]//g' -e 's/\x1b[()][A-Z0-9]//g' | tr '\r' '\n' | grep -v -E '^\s*$' | awk '!seen[$0]++' | head -60; }
dns_seit() { tail -n +$(($1+1)) $log | grep -E 'query\[(A|AAAA)\]' | awk '{print $6}' | sort | uniq -c; }
neu() {
  echo "-- Dateien neu/geaendert seit der letzten Marke ($1): Pfad, Groesse, Rechte, Einhaengepunkt"
  { find "$HOME" /etc /var /usr/local /opt /root /tmp -type f -newer /tmp/m.prev 2> /dev/null; } \
    | grep -v -E '^/tmp/(m\.|dns|ollama\.log|o-)|^/var/(lib/apt|cache|log)|^/var/lib/dpkg|^/var/lib/u-boot|^/etc/(resolv|hosts)' | sort | while read -r f; do
      [ -f "$f" ] || { echo "$f (beim Auflisten schon entfernt)"; continue; }
      printf '%s %s %s %s\n' "$f" "$(stat -c %s "$f")" "$(stat -c %a "$f")" "$(df --output=target "$f" 2> /dev/null | tail -1)"
    done | head -40
  touch /tmp/m.prev
}

# Zusatzzeilen fuer die dnsmasq-Konfiguration (HANDVARIANTE): $extra, getrennt durch ';'. Ein Eintrag "host:<name>" steht fuer die zwei Zeilen, die
# das Egress-Skript fuer einen erlaubten Namen erzeugt (server=/<name>/<upstream> und nftset=/<name>/4#ip#u_boot_egress#allowed); jeder andere Eintrag
# ist eine Zeile der Konfiguration.
zusatz() {
  local z; local IFS=';'
  for z in $extra; do
    case "$z" in
      host:*) echo "server=/${z#host:}/$up"; echo "nftset=/${z#host:}/4#ip#u_boot_egress#allowed" ;;
      "") ;;
      *) echo "$z" ;;
    esac
  done
}

cmd=$1; shift
case "$cmd" in
setup)
  mode=$1; extra=${2:-}
  sudo touch "$log"; sudo chmod 666 "$log"; : | sudo tee "$log" > /dev/null
  up=$(awk '/^nameserver/ && $2 !~ /:/ {print $2}' /etc/resolv.conf.u-boot-orig 2> /dev/null | head -1)
  [ -n "$up" ] || up=$(awk '/^nameserver/ && $2 !~ /:/ {print $2}' /etc/resolv.conf | head -1)
  if [ "$mode" = none ]; then
    echo "-- Standardprofil: kein Resolver-Protokoll (das Image traegt kein dnsmasq); das Log bleibt leer"
  elif [ "$mode" = discover ]; then
    printf 'listen-address=127.0.0.1\nbind-interfaces\nno-resolv\nno-hosts\nserver=%s\nlog-queries\nlog-facility=%s\n' "$up" "$log" | sudo tee /tmp/dns-discover.conf > /dev/null
    zusatz | sudo tee -a /tmp/dns-discover.conf > /dev/null
    sudo dnsmasq --conf-file=/tmp/dns-discover.conf
    echo "nameserver 127.0.0.1" | sudo tee /etc/resolv.conf > /dev/null
  else
    if [ ! -f /run/u-boot-egress-dnsmasq.pid ]; then
      echo "-- Egress-Skript lief nicht als postStartCommand; Aufruf wie dieses (HANDVARIANTE): sudo sh /usr/local/bin/u-boot-egress-init"
      sudo sh /usr/local/bin/u-boot-egress-init; echo "egress-init exit=$?"
    fi
    echo "-- Status-Datei des Startscripts"; sudo cat /var/lib/u-boot/status/egress
    grep -E '^(server|nftset)' /etc/u-boot-egress-dnsmasq.conf | sed 's/^/conf: /'
    sudo kill "$(cat /run/u-boot-egress-dnsmasq.pid)"; sleep 1
    { cat /etc/u-boot-egress-dnsmasq.conf; echo "log-queries"; echo "log-facility=$log"; zusatz; } | sudo tee /tmp/dns-restricted.conf > /dev/null
    sudo dnsmasq --conf-file=/tmp/dns-restricted.conf
  fi
  sleep 1
  touch /tmp/m.prev
  echo "modus=$mode benutzer=$(id -un) HOME=$HOME zusatzzeilen=$(zusatz | wc -l)"
  zusatz | sed 's/^/zusatz: /'
  echo "-- Namen der Umgebungsvariablen im Container (Werte nicht ausgegeben)"
  env | cut -d= -f1 | sort | tr '\n' ' ' | fold -w 150; echo
  echo "AUSSAGE umgebung-container-verbotene-namen=$(env | cut -d= -f1 | grep -c -E 'ANTHROPIC|TOKEN|API_KEY|SECRET|OLLAMA')"
  echo "-- /etc/hosts-Zeile fuer host.docker.internal: $(grep host.docker.internal /etc/hosts | tr -s '\t ' ' ' || true)"
  echo "-- Einhaengepunkt von ~/.ollama und von HOME"
  df --output=target,fstype "$HOME/.ollama" 2> /dev/null | tail -1
  df --output=target "$HOME" | tail -1
  ;;
fall)
  label=$1; t=$2; p=$3; si=$4; envs=$5; shift 5
  args=(); for a in "$@"; do args+=("${a//@PROMPT@/$PROMPT}"); done
  names=$(for z in $envs; do [ "$z" = - ] || echo -n "${z%%=*} "; done)
  mark "$label"
  echo "befehl: $* (env-namen: ${names:-keine}; prompt=$p; stdin=$si)"
  n=$(wc -l < $log); s=$(date +%s.%N)
  [ "$envs" = - ] && envs=""
  if [ "$si" = 1 ]; then printf '%s\n' "$PROMPT" > /tmp/o-in.txt; else : > /tmp/o-in.txt; fi
  timeout -k 5 "$t" env $envs "${args[@]}" < /tmp/o-in.txt > /tmp/o-$$.txt 2>&1; rc=$?
  e=$(date +%s.%N)
  if [ "$p" = 1 ]; then
    zl=$(strip < /tmp/o-$$.txt | wc -l); by=$(wc -c < /tmp/o-$$.txt)
    echo "ausgabe: $zl Zeilen nach Entfernen der Steuerzeichen, $by Bytes (Inhalt wird nicht ausgewiesen)"
    echo "-- Zeilen, die nach einer Meldung aussehen (Muster $MELDUNG; hoechstens 10, je 200 Zeichen)"
    strip < /tmp/o-$$.txt | grep -i -E "$MELDUNG" | cut -c1-200 | head -10
    echo "uebrige Zeilen (ohne Meldungsmuster, Inhalt nicht ausgewiesen): $(strip < /tmp/o-$$.txt | grep -v -i -E "$MELDUNG" | wc -l), laengste $(strip < /tmp/o-$$.txt | grep -v -i -E "$MELDUNG" | awk '{ if (length($0) > m) m = length($0) } END { print m + 0 }') Zeichen"
  else
    strip < /tmp/o-$$.txt | cut -c1-300
  fi
  echo "$label exit=$rc dauer=$(dauer "$s" "$e")s (124 = timeout nach ${t}s, 137 = beendet nach -k 5)"
  echo "AUSSAGE $label-exit=$rc"
  echo "-- Namen im Resolver-Log waehrend $label (Anzahl Anfragen A/AAAA)"
  dns_seit "$n"
  echo "-- abgelehnte/unbekannte Namen waehrend $label"
  tail -n +$((n+1)) $log | grep -E 'REFUSED|NXDOMAIN|NODATA|config ' | awk '{print $4, $5, $6}' | sort | uniq -c | head
  rm -f /tmp/o-$$.txt
  neu "$label"
  ;;
serve)
  nohup env "$@" ollama serve > /tmp/ollama.log 2>&1 &
  for i in $(seq 1 60); do curl -fsS -m 2 http://127.0.0.1:11434/ > /dev/null 2>&1 && break; sleep 0.5; done
  echo "serve bereit: $(curl -sS -m 2 http://127.0.0.1:11434/api/version)"
  ;;
stopserve)
  pkill -f 'ollama serve' 2> /dev/null; sleep 2; echo "serve beendet: $(pgrep -c -f 'ollama serve' || true) Prozesse"
  ;;
dateien) neu "$1" ;;
ablage)
  echo "-- ~/.ollama (Pfad, Groesse, Rechte, Eigentuemer; kein Inhalt)"
  find "$HOME/.ollama" -maxdepth 2 -printf '%p %s %m %u\n' 2> /dev/null | sort | head -30
  for f in "$HOME/.ollama/id_ed25519" "$HOME/.ollama/id_ed25519.pub" "$HOME/.ollama/config.json" "$HOME/.claude.json"; do
    [ -e "$f" ] && printf 'einhaengepunkt %s -> %s\n' "$f" "$(df --output=target "$f" | tail -1)"
  done
  echo "-- Einhaengepunkte (df): $(df --output=target | tail -n +2 | tr '\n' ' ')"
  ;;
dns)
  echo "-- Namen im Resolver-Log ab Zeile ${1:-0} (Anzahl Anfragen A/AAAA)"; dns_seit "${1:-0}"
  echo "-- Zeilen im Log: $(wc -l < $log)"
  ;;
nftadd)
  # HANDVARIANTE (K3c): die Adresse von Hand in die nftables-Menge der Egress-Restriktion eintragen; zeigt, ob der Weg offen ist, wenn die Menge die Adresse traegt
  sudo nft add element ip u_boot_egress allowed "{ $1 }"; echo "nft add element exit=$?"
  ;;
nft)
  if sudo nft list set ip u_boot_egress allowed > /tmp/o-nft.txt 2>&1; then
    echo "nft set allowed: Elemente=$(grep -o -E '[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+' /tmp/o-nft.txt | wc -l) 172.17.0.1-enthalten=$(grep -c -F '172.17.0.1' /tmp/o-nft.txt)"
  else
    echo "nft set allowed nicht lesbar: $(head -1 /tmp/o-nft.txt | cut -c1-120)"
  fi
  rm -f /tmp/o-nft.txt
  ;;
*) echo "unbekannter Aufruf: $cmd"; exit 2 ;;
esac
