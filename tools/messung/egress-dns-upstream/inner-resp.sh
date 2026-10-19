#!/bin/sh
# Laeuft im Antwort-Container der Messung als root (Image u-boot-messung-dns-img-tool): ein Gegenueber, das DNS beantwortet und Verbindungen annimmt,
# damit D3 und D5 ein Ziel im Wegwerf-Netz haben (kein Zugriff auf fremde Server). Es beantwortet JEDEN Namen mit 203.0.113.7 (TEST-NET-3).
#   dnsmasq Port 53 und Port 5353 (Anfragen im Log /tmp/dns53.log, /tmp/dns5353.log), TCP-Listener 853 und 443 (Eingang in /tmp/tcp853.in,
#   /tmp/tcp443.in).
set -eu
for p in 53 5353; do
  dnsmasq --conf-file=/dev/null --no-resolv --no-hosts --address=/#/203.0.113.7 --log-queries --log-facility=/tmp/dns$p.log \
    --pid-file=/tmp/dnsmasq$p.pid --port=$p --user=root
done
for p in 853 443; do
  nohup sh -c "while :; do nc -l $p >> /tmp/tcp$p.in 2>/dev/null; done" > /dev/null 2>&1 &
done
sleep 1
echo "antwort-container bereit: $(hostname -i)"
