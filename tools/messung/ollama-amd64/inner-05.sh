#!/bin/bash
# Laeuft IM Container (Image oll) als Benutzer vscode, wie im Devcontainer; aufgerufen von 05-laufzeit.sh.
# Aufruf im Container: inner-05.sh <modus> <modell...>
#   discover    keine Restriktion; protokollierender dnsmasq reicht alle Namen weiter (zeigt, welche Hosts Pull und Lauf brauchen)
#   restricted  das erzeugte Startscript unveraendert (`sudo sh /usr/local/bin/u-boot-egress-init`); danach ersetzt ein protokollierender
#               dnsmasq mit derselben Konfiguration den Resolver des Skripts (Verfahren des Vorgaenger-Protokolls)
set -u
mode=$1; shift
models=("$@")
log=/tmp/dns.log
now() { date -u +%T.%N | cut -c1-12; }
mark() { echo "---- $(now) $*"; }
dauer() { awk -v a="$1" -v b="$2" 'BEGIN{printf "%.1f", b-a}'; }
# Fortschrittsanzeige von `ollama pull` ohne Steuerzeichen und Dubletten, vorsignierte URLs ohne Query
strip() { sed -e 's/\x1b\[[0-9;?]*[a-zA-Z]//g' | tr '\r' '\n' | grep -v -E '^\s*$' | awk '!seen[$0]++' | sed -E 's/\?[^" ]*/?[Query entfernt]/g'; }
sudo touch "$log"; sudo chmod 666 "$log"
up=$(awk '/^nameserver/ && $2 !~ /:/ {print $2}' /etc/resolv.conf | head -1)
echo "resolver des Containers vor dem Skript: $(grep -c '^nameserver' /etc/resolv.conf) Eintrag/Eintraege (erster IPv4-Upstream vorhanden: $([ -n "$up" ] && echo ja || echo nein))"

if [ "$mode" = discover ]; then
  cp /etc/resolv.conf /tmp/resolv.conf.orig
  printf 'listen-address=127.0.0.1\nbind-interfaces\nno-resolv\nno-hosts\nserver=%s\nlog-queries\nlog-facility=%s\n' "$up" "$log" | sudo tee /tmp/dns-discover.conf > /dev/null
  sudo dnsmasq --conf-file=/tmp/dns-discover.conf
  echo "nameserver 127.0.0.1" | sudo tee /etc/resolv.conf > /dev/null
else
  sudo sh /usr/local/bin/u-boot-egress-init; echo "egress-init exit=$?"
  echo "-- Status-Datei des Skripts"; cat /var/lib/u-boot/status/egress
  echo "-- nft-Tabellen"; sudo nft list tables
  grep -E '^(server|nftset)' /etc/u-boot-egress-dnsmasq.conf | sed 's/^/conf: /'
  sudo kill "$(cat /run/u-boot-egress-dnsmasq.pid)"; sleep 1
  { cat /etc/u-boot-egress-dnsmasq.conf; echo "log-queries"; echo "log-facility=$log"; } | sudo tee /tmp/dns-restricted.conf > /dev/null
  sudo dnsmasq --conf-file=/tmp/dns-restricted.conf
fi
sleep 1
hostsline=$(grep '^hosts=' /usr/local/bin/u-boot-egress-init | cut -d'"' -f2)

mark "serve starten (HOME=$HOME)"
env | grep -i '^OLLAMA' || echo "(keine OLLAMA_*-Variable gesetzt: OLLAMA_HOST ungesetzt)"
n0=$(wc -l < $log)
s0=$(date +%s.%N)
nohup ollama serve > /tmp/ollama.log 2>&1 &
spid=$!
for i in $(seq 1 60); do curl -fsS -m 2 http://127.0.0.1:11434/ > /dev/null 2>&1 && break; sleep 0.5; done
s1=$(date +%s.%N)
echo "AUSSAGE serve-bereit=$(curl -sS -m 2 http://127.0.0.1:11434/)"
echo "serve bereit nach $(dauer "$s0" "$s1")s"
curl -sS -m 2 http://127.0.0.1:11434/api/version; echo
echo "-- lauschende Sockets"
ss -ltn | grep -E 'State|11434'
echo "AUSSAGE bindung=$(ss -ltn | awk '$4 ~ /:11434$/ {print $4}' | tr '\n' ' ')"
echo "-- Erreichbarkeit ueber die Container-Adresse (Exit-Code von curl: 7 = abgelehnt, 28 = Zeitueberschreitung/verworfen)"
ip4=$(hostname -I | awk '{print $1}')
curl -sS -m 3 -o /dev/null "http://$ip4:11434/" 2> /tmp/c.err; rc=$?
echo "container-adresse: curl-exit=$rc ($(cut -c1-100 /tmp/c.err))"
echo "AUSSAGE container-adresse-curl-exit=$rc"
if [ "$mode" = restricted ]; then
  echo "-- Gegenstueck: eine verworfene Verbindung (gesperrte Adresse, kein Namens-Aufruf)"
  curl -sS -m 3 -o /dev/null https://1.1.1.1/ 2> /tmp/c.err; rc=$?
  echo "gesperrte Adresse 1.1.1.1: curl-exit=$rc ($(cut -c1-100 /tmp/c.err))"
  echo "AUSSAGE gesperrte-adresse-curl-exit=$rc"
fi
echo "-- Server-Log: Start, Bibliothek, Threads"
grep -i -E 'listening|library=|threads|OLLAMA_HOST|inference compute|no compatible GPU|cpu' /tmp/ollama.log | cut -c1-260 | head -12
echo "AUSSAGE library-cpu=$(grep -c 'library=cpu' /tmp/ollama.log)"
echo "-- ~/.ollama nach dem Start"; ls -la "$HOME/.ollama" | head -6
echo "-- DNS-Anfragen seit Serverstart (vor dem Pull)"
tail -n +$((n0+1)) $log | grep -E 'query\[(A|AAAA)\]' | awk '{print $6}' | sort | uniq -c

for m in "${models[@]}"; do
  mark "pull $m"
  n1=$(wc -l < $log)
  s=$(date +%s.%N)
  ollama pull "$m" > /tmp/pull.txt 2>&1; rc=$?
  e=$(date +%s.%N)
  strip < /tmp/pull.txt | grep -E '100%|success|Error|error|dial|lookup|timeout' | tail -8
  echo "pull $m exit=$rc dauer=$(dauer "$s" "$e")s"
  echo "AUSSAGE pull-$m-exit=$rc"
  echo "-- DNS-Anfragen waehrend des Pulls ($m)"
  tail -n +$((n1+1)) $log | grep -E 'query\[(A|AAAA)\]' | awk '{print $6}' | sort | uniq -c
  if [ "$mode" = discover ]; then
    echo "-- Abgleich der angefragten Namen mit hosts= des erzeugten Skripts (gedeckt: Name oder Oberdomain steht in hosts=)"
    tail -n +$((n1+1)) $log | grep -E 'query\[A\]' | awk '{print $6}' | sort -u | while read -r n; do
      ok=0; for h in $hostsline; do case "$n" in "$h"|*."$h") ok=1;; esac; done
      [ $ok = 1 ] && echo "gedeckt: $n" || echo "NICHT gedeckt: $n"
    done
  else
    echo "-- Antworten fuer Namen ausserhalb der Liste (REFUSED/NXDOMAIN/NODATA/config)"
    tail -n +$((n1+1)) $log | grep -E 'REFUSED|NXDOMAIN|NODATA|config ' | awk '{print $4, $5, $6}' | sort -u | head
  fi
done

if [ "${#models[@]}" -gt 0 ]; then
  mark "Speicher nach den Pulls"
  du -sh "$HOME/.ollama" "$HOME/.ollama/models" 2>&1
  ollama list
  for m in "${models[@]}"; do
    mark "run $m (CPU)"
    n2=$(wc -l < $log)
    s=$(date +%s.%N)
    ollama run "$m" "Reply with the single word OK." > /tmp/run.txt 2>&1; rc=$?
    e=$(date +%s.%N)
    strip < /tmp/run.txt | tail -3
    echo "run $m exit=$rc dauer=$(dauer "$s" "$e")s"
    echo "AUSSAGE run-$m-exit=$rc antwort-nichtleer=$([ -n "$(strip < /tmp/run.txt)" ] && echo ja || echo nein)"
    echo "-- ollama ps"; ollama ps
    echo "-- DNS-Anfragen waehrend run"
    tail -n +$((n2+1)) $log | grep -E 'query\[(A|AAAA)\]' | awk '{print $6}' | sort | uniq -c
  done
  echo "-- cgroup memory.peak (Bytes)"; cat /sys/fs/cgroup/memory.peak 2>&1
  echo "-- Server-Log (Bibliothek/Threads waehrend der Laeufe)"
  grep -i -E 'library=|threads|runner|num_thread' /tmp/ollama.log | cut -c1-260 | sort -u | head -8
fi

if [ "$mode" = restricted ]; then
  mark "gesperrte Hosts (Namensaufloesung ueber den Resolver des Skripts, dann curl)"
  for h in ollama.com example.com huggingface.co; do
    curl -sS -m 6 -o /dev/null "https://$h/" 2> /tmp/c.err; rc=$?
    echo "$h: curl-exit=$rc ($(cut -c1-110 /tmp/c.err))"
    echo "AUSSAGE gesperrt-$h-curl-exit=$rc"
  done
  echo "-- erlaubte Hosts (Gegenprobe)"
  for h in github.com registry.ollama.ai; do
    curl -sS -m 8 -o /dev/null "https://$h/" 2> /tmp/c.err; rc=$?
    echo "$h: curl-exit=$rc ($(cut -c1-110 /tmp/c.err))"
  done
fi
mark "alle DNS-Anfragen der Sitzung (eindeutig, A/AAAA)"
grep -E 'query\[(A|AAAA)\]' $log | awk '{print $6}' | sort | uniq -c
kill $spid 2> /dev/null; sleep 1
