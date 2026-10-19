#!/bin/bash
# D1: Benutzer von dnsmasq und vscode; unterscheidet `meta skuid` eine DNS-Anfrage des dnsmasq von der eines Prozesses des Benutzers?
# Aufruf: WORK=<messordner> OUT=<ausgabeverzeichnis> 03-d1.sh -> $OUT/rohdaten-03-d1.txt
# Voraussetzung: 02-projekte.sh. Die Zaehler-Regeln sind Messinstrument im Wegwerf-Container (nicht terminierend, vor den Upstream-Regeln).
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-03-d1.txt
c=$PFX-d1
echo "######## D1, $(now)"
start $c $IMG_TOOL
egress $c ist
UPS=$(upstreams $c); echo "Upstream-Liste: $UPS"
echo "-- Benutzer im Image"
docker exec -u vscode $c sh -c 'id; getent passwd vscode nobody dnsmasq'
echo "-- Prozess dnsmasq nach dem Startscript"
docker exec -u root $c sh -c 'ps -eo user:12,uid,pid,ppid,args | grep -E "^(USER|.*dnsmasq --conf)" | grep -v grep'
pid=$(docker exec -u root $c cat /run/u-boot-egress-dnsmasq.pid)
echo "pid-datei=$pid; /proc/$pid/status: $(docker exec -u root $c grep -E '^Uid:' /proc/$pid/status | tr '\t' ' ')"
echo "-- dnsmasq-Konfiguration: user=? (Zeilen mit user oder group)"
docker exec -u root $c sh -c 'grep -c -E "^(user|group)=" /etc/u-boot-egress-dnsmasq.conf; true'
echo "-- Zaehler je Benutzer-ID (Messinstrument): vor den Upstream-Regeln zaehlen zwei nicht terminierende Regeln"
docker exec -u root $c nft insert rule ip u_boot_egress output oifname != lo meta l4proto '{ tcp, udp }' th dport 53 meta skuid 1000 counter
docker exec -u root $c nft insert rule ip u_boot_egress output oifname != lo meta l4proto '{ tcp, udp }' th dport 53 meta skuid 65534 counter
echo "-- Anfrage 1: vscode fragt den lokalen Resolver nach einem erlaubten Namen (der dnsmasq fragt den Upstream)"
docker exec -u vscode $c dig +time=3 +tries=1 +short @127.0.0.1 api.github.com A | head -2
echo "-- Anfrage 2: vscode fragt den ersten Upstream direkt"
ns1=$(echo $UPS | cut -d' ' -f1)
docker exec -u vscode $c dig +time=3 +tries=1 +short @"$ns1" example.com A | head -2
echo "-- Zaehler"
docker exec -u root $c sh -c 'nft list chain ip u_boot_egress output | grep counter | sed "s/^[[:space:]]*/    /"'
docker exec -u root $c sh -c 'nft list chain ip u_boot_egress output' > "$WORK/d1-chain.txt"
c1000=$(grep 'skuid 1000 counter' "$WORK/d1-chain.txt" | sed -n 's/.*packets \([0-9]*\).*/\1/p')
c65534=$(grep 'skuid 65534 counter' "$WORK/d1-chain.txt" | sed -n 's/.*packets \([0-9]*\).*/\1/p')
echo "pakete skuid 1000 (vscode, direkte Anfrage): $c1000; pakete skuid 65534 (nobody, dnsmasq): $c65534"
aussage "d1 dnsmasq-benutzer=$(docker exec -u root $c sh -c "ps -o user= -p $pid")"
aussage "d1 dnsmasq-uid=$(docker exec -u root $c sh -c "awk '/^Uid:/ {print \$2}' /proc/$pid/status")"
aussage "d1 vscode-uid=$(docker exec -u vscode $c id -u)"
aussage "d1 zaehler-vscode-groesser-null=$([ "${c1000:-0}" -gt 0 ] && echo ja || echo nein)"
aussage "d1 zaehler-dnsmasq-groesser-null=$([ "${c65534:-0}" -gt 0 ] && echo ja || echo nein)"
docker rm -f -v $c > /dev/null
