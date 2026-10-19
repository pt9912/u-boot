#!/bin/bash
# Laeuft im Container als vscode, OHNE sudo (ausser dem Lesen der Regeln); $1 = Modus (ohne|ist|skuid), $2 = IPv6-Adresse des Antwort-Containers.
# D3, Zusatz im benutzerdefinierten Docker-Netz (mit IPv6): IPv6 gegen die Tabelle ip6 (Zeilen 67 bis 69 des Templates) und der Docker-eigene
# Resolver 127.0.0.11. Siehe 07-d3-netz.sh.
set -u
mode=${1:?ohne|ist|skuid}
R6=${2:?ipv6-adresse}
orig=/etc/resolv.conf.u-boot-orig; [ -f "$orig" ] || orig=/etc/resolv.conf
echo "resolv.conf-Original: $(grep nameserver "$orig" | tr '\n' ' ')"
echo "IPv6-Adressen (ausser ::1): $(ip -6 addr | grep -c 'inet6 [^:]'; true); eigene Adresse: $(ip -6 addr show scope global | awk '/inet6/ {print $2}' | head -1)"
A() { echo "AUSSAGE $mode d3-netz $1=$2"; }
res() { if [ "$1" = 0 ] && [ -n "$2" ]; then echo offen; else echo gesperrt; fi; }
o=$(dig +time=2 +tries=1 +short -6 @"$R6" example.com A 2>/dev/null | grep -v -e '^;' -e '^$' | head -1); A "ipv6-udp53-R6" "$(res $? "$o")"
o=$(dig +time=2 +tries=1 +short -6 +tcp @"$R6" example.com A 2>/dev/null | grep -v -e '^;' -e '^$' | head -1); A "ipv6-tcp53-R6" "$(res $? "$o")"
printf 'u-boot-messung-v6' | nc -6 -w 2 "$R6" 853 > /dev/null 2>&1; A "ipv6-tcp853-R6" "$([ $? = 0 ] && echo offen || echo gesperrt)"
if ping -6 -c 1 -W 2 "$R6" > /dev/null 2>&1; then A "ipv6-icmp-R6" offen; else A "ipv6-icmp-R6" gesperrt; fi
echo "--- Docker-eigener Resolver 127.0.0.11 (Name ausserhalb der Liste: example.com; Name der Liste: api.github.com)"
for h in example.com api.github.com; do
  out=$(dig +time=3 +tries=1 +short @127.0.0.11 "$h" A 2>/dev/null | grep -v -e '^;' -e '^$' | head -1)
  echo "    @127.0.0.11 $h: ${out:-keine Antwort}"
  A "docker-resolver-127.0.0.11 $h" "$([ -n "$out" ] && echo antwort || echo keine)"
done
out=$(getent hosts example.com | head -1 | cut -d' ' -f1); echo "    getent example.com (resolv.conf): ${out:-keine Adresse}"
A "getent example.com" "$([ -n "$out" ] && echo adresse || echo keine)"
out=$(dig +time=3 +tries=1 +short @127.0.0.1 api.github.com A 2>/dev/null | grep -v -e '^;' -e '^$' | head -1)
echo "    lokaler dnsmasq api.github.com: ${out:-keine Antwort}"
A "lokaler-dnsmasq api.github.com" "$([ -n "$out" ] && echo antwort || echo keine)"
echo "--- Regeln (sudo nur zum Lesen)"
sudo -n nft list tables 2>&1 | sed 's/^/    /'
