#!/bin/bash
# Laeuft im Probe-Container als Benutzer vscode (wie im Devcontainer).
# Aufruf: inner-05.sh <modus> "<hosts>" <modell...>
#   modus discover : keine Restriktion, protokollierender dnsmasq-Resolver (alle Namen werden weitergereicht)
#   modus restricted: erzeugtes Startscript mit ersetzter hosts-Zeile, danach protokollierender Resolver mit derselben Konfiguration
set -u
mode=$1; hosts=$2; shift 2
models=("$@")
log=/tmp/dns.log
now() { date -u +%T.%N | cut -c1-12; }
mark() { echo "---- $(now) $*"; }
strip() { sed -e 's/\x1b\[[0-9;?]*[a-zA-Z]//g' | tr '\r' '\n' | grep -v -E '^\s*$' | awk '!seen[$0]++'; }
sudo touch "$log"; sudo chmod 666 "$log"
up=$(awk '/^nameserver/ && $2 !~ /:/ {print $2}' /etc/resolv.conf | head -1)

if [ "$mode" = discover ]; then
  sudo cp /etc/resolv.conf /etc/resolv.conf.orig
  printf 'listen-address=127.0.0.1\nbind-interfaces\nno-resolv\nno-hosts\nserver=%s\nlog-queries\nlog-facility=%s\n' "$up" "$log" | sudo tee /tmp/dns-discover.conf >/dev/null
  sudo dnsmasq --conf-file=/tmp/dns-discover.conf
  echo "nameserver 127.0.0.1" | sudo tee /etc/resolv.conf >/dev/null
else
  sed -e "s|^hosts=.*|hosts=\"$hosts\"|" /usr/local/bin/u-boot-egress-init > /tmp/egress-init.sh
  sudo sh /tmp/egress-init.sh; echo "egress-init exit=$?"
  grep -E '^(server|nftset)' /etc/u-boot-egress-dnsmasq.conf | sed 's/^/conf: /'
  sudo kill "$(cat /run/u-boot-egress-dnsmasq.pid)"; sleep 1
  { cat /etc/u-boot-egress-dnsmasq.conf; echo "log-queries"; echo "log-facility=$log"; } | sudo tee /tmp/dns-restricted.conf >/dev/null
  sudo dnsmasq --conf-file=/tmp/dns-restricted.conf
fi
sleep 1

mark "serve starten (Umgebung: OLLAMA_HOST=${OLLAMA_HOST:-<ungesetzt>}, HOME=$HOME)"
env | grep -i '^OLLAMA' || echo "(keine OLLAMA_*-Variable)"
ls -la "$HOME/.ollama" 2>&1 | head -3
n0=$(wc -l < $log)
nohup ollama serve > /tmp/ollama.log 2>&1 &
spid=$!
for i in $(seq 1 30); do curl -fsS -m 2 http://127.0.0.1:11434/ >/dev/null 2>&1 && break; sleep 1; done
mark "serve bereit nach ${i}s"
curl -sS -m 2 http://127.0.0.1:11434/; echo
curl -sS -m 2 http://127.0.0.1:11434/api/version; echo
echo "-- lauschende Sockets"
ss -ltnp 2>/dev/null | grep -E 'State|11434|ollama'
echo "-- Erreichbarkeit ueber die Container-Adresse"
ip4=$(hostname -I | awk '{print $1}')
echo "container-ip=$ip4"
curl -sS -m 3 -o /dev/null -w "http://$ip4:11434 -> %{http_code}\n" "http://$ip4:11434/" 2>&1 | tail -1
echo "-- Server-Log (Start)"; head -40 /tmp/ollama.log | cut -c1-220
echo "-- ~/.ollama nach dem Start"; ls -laR "$HOME/.ollama" 2>&1 | head -20
echo "-- DNS-Anfragen seit Serverstart (vor dem Pull)"
tail -n +$((n0+1)) $log | grep -E 'query\[(A|AAAA)\]' | awk '{print $6}' | sort | uniq -c

for m in "${models[@]}"; do
  mark "pull $m"
  n1=$(wc -l < $log)
  s=$(date +%s.%N)
  ollama pull "$m" 2>&1 | strip | sed -E 's/\?[^" ]*/?[Query entfernt: vorsignierte URL]/g' | grep -E '100%|success|Error|error|dial|lookup|timeout' | tail -8
  rc=${PIPESTATUS[0]}
  e=$(date +%s.%N)
  echo "pull exit=$rc dauer=$(echo "$e - $s" | bc 2>/dev/null || awk "BEGIN{print $e-$s}")s"
  echo "-- DNS-Anfragen waehrend des Pulls ($m)"
  tail -n +$((n1+1)) $log | grep -E 'query\[(A|AAAA)\]' | awk '{print $6}' | sort | uniq -c
  echo "-- DNS-Antworten (CNAME-Kette) waehrend des Pulls"
  tail -n +$((n1+1)) $log | grep -E ' is <CNAME>|config .* is NODATA|SERVFAIL|REFUSED' | awk '{print $4, $5, $6}' | sort -u | head -20
done

if [ "${#models[@]}" -gt 0 ] && [ "$rc" = 0 ]; then
  m=${models[0]}
  mark "Speicher nach dem Pull"
  du -sh "$HOME/.ollama" "$HOME/.ollama/models" 2>&1
  ls -la "$HOME/.ollama" "$HOME/.ollama/models" "$HOME/.ollama/models/blobs" 2>&1 | head -30
  df -h "$HOME" | tail -1
  mark "run $m (CPU-only)"
  n2=$(wc -l < $log)
  s=$(date +%s.%N)
  ollama run "$m" "Reply with the single word OK." 2>&1 | strip | tail -3
  e=$(date +%s.%N)
  awk "BEGIN{printf \"run dauer=%.1fs\n\", $e-$s}"
  echo "-- ollama ps"; ollama ps
  echo "-- Prozesse (RSS in kB)"; ps -eo pid,rss,vsz,args --sort=-rss | grep -E 'ollama|llama' | grep -v grep | cut -c1-160 | head -5
  echo "-- cgroup memory.peak"; cat /sys/fs/cgroup/memory.peak 2>&1
  echo "-- DNS-Anfragen waehrend run"
  tail -n +$((n2+1)) $log | grep -E 'query\[(A|AAAA)\]' | awk '{print $6}' | sort | uniq -c
  echo "-- Server-Log (GPU/CPU-Zeilen)"
  grep -i -E 'gpu|cpu|compute|library|inference|total_vram|memory' /tmp/ollama.log | cut -c1-240 | head -12
fi
mark "alle DNS-Anfragen (eindeutig)"
grep -E 'query\[(A|AAAA)\]' $log | awk '{print $6}' | sort | uniq -c
kill $spid 2>/dev/null; sleep 1
