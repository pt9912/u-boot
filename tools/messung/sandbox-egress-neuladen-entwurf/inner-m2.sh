#!/bin/bash
# Laeuft im Container als root (siehe 03-m2-neuladen.sh); $1 = Abschnitt A..E.
# Voraussetzung: der Container hat den postStartCommand einmal ausgefuehrt (Baseline-Liste = hosts= des erzeugten Scripts);
# /m enthaelt die Dateien des Messverzeichnisses (probe.pl, reload-atomar.sh).
set -u
sect=${1:?A..E}
L2="api.github.com example.org"
mk() { sed -e "s|^hosts=.*|hosts=\"$1\"|" /usr/local/bin/u-boot-egress-init > "$2"; }
mk "$(sed -n 's/^hosts="\(.*\)"$/\1/p' /usr/local/bin/u-boot-egress-init)" /tmp/l1.sh
mk "$L2" /tmp/l2.sh
echo "Liste L1 (gebacken): $(sed -n 's/^hosts="\(.*\)"$/\1/p' /tmp/l1.sh)"
echo "Liste L2: $L2"
alive() { ps -p "$(cat /run/u-boot-egress-dnsmasq.pid 2>/dev/null || echo 0)" >/dev/null 2>&1 && echo ja || echo nein; }

case "$sect" in
A|E) # A: Fenster beim erneuten Lauf des Scripts; E: dasselbe mit der atomaren Variante
  if [ "$sect" = A ]; then run() { sh /tmp/l2.sh; }; else run() { sh /m/reload-atomar.sh "$L2"; }; fi
  perl /m/probe.pl 6 api.github.com egress > /tmp/probe-egress.txt &
  perl /m/probe.pl 6 api.github.com dns > /tmp/probe-dns.txt &
  sleep 2
  s=$(date +%s.%N); run > /tmp/run.out 2>&1; rc=$?; e=$(date +%s.%N)
  wait
  echo "Lauf: exit=$rc start=$s ende=$e dauer_s=$(echo "$e - $s" | bc 2>/dev/null || awk "BEGIN{printf \"%.3f\", $e - $s}")"
  echo "Ausgabe des Laufs:"; sed 's/^/    /' /tmp/run.out
  echo "Sonde Egress (Zustandswechsel; offen = kein Filter-Regelwerk):"; cat /tmp/probe-egress.txt
  echo "Sonde DNS (Zustandswechsel; lokaler dnsmasq):"; cat /tmp/probe-dns.txt
  echo "Zombies im Container danach:"; ps -eo stat,user,pid,args | awk '$1 ~ /^Z/' | sed 's/^/    /'
  ;;
B|F) # B: 50 erneute Laeufe hintereinander; F: dasselbe mit der atomaren Variante (Rennen um den Port von dnsmasq)
  n=50; bad=0; badbind=0
  for i in $(seq 1 $n); do
    if [ $((i % 2)) = 1 ]; then list="$L2"; f=/tmp/l2.sh; else list="x.example.net"; mk "$list" /tmp/lx.sh; f=/tmp/lx.sh; fi
    if [ "$sect" = B ]; then sh "$f" > /tmp/o.txt 2>&1; else sh /m/reload-atomar.sh "$list" > /tmp/o.txt 2>&1; fi
    if grep -q -i -E 'address already in use|failed to create|error' /tmp/o.txt; then badbind=$((badbind+1)); echo "lauf $i: Meldung: $(grep -i -E 'address already in use|failed to create|error' /tmp/o.txt | head -1)"; fi
    host=$(echo "$list" | awk '{print $NF}')
    r=$(dig +time=2 +tries=1 @127.0.0.1 "$host" A 2>&1 | grep -o 'status: [A-Z]*' | head -1)
    nd=$(ps -u nobody -o args= | grep -c '^dnsmasq')
    if [ -z "$r" ] || [ "$nd" != 1 ] || ! grep -q "^server=/$host/" /etc/u-boot-egress-dnsmasq.conf || [ "$(alive)" != ja ]; then
      bad=$((bad+1)); echo "lauf $i: AUFFAELLIG liste=$list antwort='$r' dnsmasq-prozesse=$nd pid-lebt=$(alive)"
    fi
  done
  echo "Zusammenfassung: $n Laeufe, $bad mit auffaelligem Zustand danach, $badbind mit Fehlermeldung in der Ausgabe"
  ;;
G) # Zeit vom SIGTERM bis zum Ende des alten dnsmasq (Spielraum des `kill` ohne Warten im gebackenen Script)
  for i in $(seq 1 20); do
    sh /tmp/l1.sh > /dev/null 2>&1
    pid=$(cat /run/u-boot-egress-dnsmasq.pid)
    s=$(date +%s%N); kill "$pid"
    while :; do st=$(ps -o stat= -p "$pid" 2>/dev/null | cut -c1); [ -z "$st" ] || [ "$st" = Z ] && break; done
    e=$(date +%s%N)
    echo "lauf $i: dnsmasq $pid beendet nach $(( (e - s) / 1000 )) us"
  done | sort -k7 -n | awk '{print} END{print NR " Laeufe"}'
  ;;
C) # laufende Verbindung ueber einen Neu-Lauf, dessen Liste den Host nicht mehr enthaelt
  sh /tmp/l1.sh
  url=https://deb.debian.org/debian/dists/stable/main/Contents-arm64.gz
  curl -sS -o /tmp/dl.bin --limit-rate 300k "$url" 2> /tmp/dl.err &
  cpid=$!
  sleep 4; echo "vor dem Neu-Lauf:   Bytes=$(stat -c %s /tmp/dl.bin) curl-lebt=$(ps -p $cpid >/dev/null && echo ja || echo nein)"
  sh /tmp/l2.sh; echo "Neu-Lauf mit L2 (ohne deb.debian.org) exit=$?"
  for t in 1 3 6 9; do sleep $([ $t = 1 ] && echo 1 || echo 2); echo "nach dem Neu-Lauf +${t}s: Bytes=$(stat -c %s /tmp/dl.bin) curl-lebt=$(ps -p $cpid >/dev/null && echo ja || echo nein)"; done
  echo "neue Verbindung zum selben Host:"
  curl -sS -m 5 -o /dev/null -w '    http=%{http_code}\n' https://deb.debian.org/debian/ 2>&1 | head -2
  kill $cpid 2>/dev/null; cat /tmp/dl.err | head -3
  ;;
D) # Set nach dem Neu-Lauf leer: Verbindung zu einer zuvor aufgeloesten Adresse ohne neue Namensaufloesung
  sh /tmp/l1.sh
  ip=$(dig +short +time=3 +tries=1 @127.0.0.1 api.github.com A | head -1)
  echo "aufgeloest: api.github.com -> $ip"
  echo "Set vor dem Neu-Lauf: $(nft list set ip u_boot_egress allowed | tr -d '\n\t' | cut -c1-160)"
  curl -sS -m 5 --resolve "api.github.com:443:$ip" -o /dev/null -w 'vor Neu-Lauf, ohne DNS: http=%{http_code}\n' https://api.github.com 2>&1 | tail -1
  sh /tmp/l2.sh
  echo "Set nach dem Neu-Lauf: $(nft list set ip u_boot_egress allowed | tr -d '\n\t' | cut -c1-160)"
  curl -sS -m 5 --resolve "api.github.com:443:$ip" -o /dev/null -w 'nach Neu-Lauf, ohne DNS: http=%{http_code}\n' https://api.github.com 2>&1 | tail -1
  curl -sS -m 5 -o /dev/null -w 'nach Neu-Lauf, mit DNS: http=%{http_code}\n' https://api.github.com 2>&1 | tail -1
  curl -sS -m 5 --resolve "api.github.com:443:$ip" -o /dev/null -w 'danach erneut ohne DNS (Adresse nun im Set): http=%{http_code}\n' https://api.github.com 2>&1 | tail -1
  ;;
esac
