#!/bin/bash
# Laeuft im Container als root (docker exec -u root); $1 = Modus (ist|skuid), $2 = Beschriftung, $3 = optional "kurz" (ohne Ausgabe der Regeln).
# D4: Zustand nach einem Lauf des Startscripts. Die direkte Abfrage laeuft als vscode (setpriv, ohne sudo). Siehe 09-d4.sh.
set -u
mode=${1:?ist|skuid}
label=${2:?beschriftung}
dm=$(pgrep -x dnsmasq | wc -l)
uid=$(for p in $(pgrep -x dnsmasq); do awk '/^Uid:/ {print $2}' /proc/$p/status; done | sort -u | tr '\n' ',')
pidfile=$(cat /run/u-boot-egress-dnsmasq.pid 2>/dev/null || echo 0)
alive=$([ -d "/proc/$pidfile" ] && echo ja || echo nein)
accept=$(nft list table ip u_boot_egress 2>/dev/null | grep -c 'accept')
skuid=$(nft list table ip u_boot_egress 2>/dev/null | grep -c 'meta skuid 65534 ip daddr')
tabellen=$(nft list tables 2>/dev/null | tr '\n' ' ' | sed 's/ $//')
lokal=$(dig +time=2 +tries=1 @127.0.0.1 api.github.com A 2>&1 | grep -o 'status: [A-Z]*' | head -1)
direkt=$(setpriv --reuid=1000 --regid=1000 --clear-groups dig +time=1 +tries=1 +short @"$(awk '/^nameserver/ && $2 !~ /:/ {print $2; exit}' /etc/resolv.conf.u-boot-orig)" example.com A 2>/dev/null | grep -c -E '^[0-9a-f.:]+$')
resolv=$(grep -c '^nameserver 127.0.0.1$' /etc/resolv.conf)
status=$(grep '^result=' /var/lib/u-boot/status/egress 2>/dev/null)
echo "[$label] dnsmasq-prozesse=$dm uid=$uid pid-lebt=$alive accept-regeln=$accept skuid-regeln=$skuid tabellen='$tabellen' lokal-api.github.com='$lokal' direkt-example.com=$([ "$direkt" -gt 0 ] && echo antwort || echo keine) resolv.conf=127.0.0.1:$resolv status=$status"
echo "AUSSAGE $mode d4 $label dnsmasq-prozesse=$dm uid=$uid accept-regeln=$accept skuid-regeln=$skuid tabellen='$tabellen' lokal='$lokal' direkt=$([ "$direkt" -gt 0 ] && echo antwort || echo keine) resolv=$resolv status=$status"
