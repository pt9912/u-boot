#!/bin/sh
# MESSVARIANTE, kein Produktcode: ersetzt die Allowlist im laufenden Container mit einem atomaren Wechsel der
# nftables-Tabellen (eine Transaktion) und wartet auf das Ende des alten dnsmasq, bevor der neue startet.
# Aufruf (als root im Container): reload-atomar.sh "<host> <host> ..."
# Gleiche Regeln wie egress-init.sh (Vorlage egress-init.sh.tmpl); nur die Reihenfolge der Schritte unterscheidet sich.
hosts="$1"
pidfile=/run/u-boot-egress-dnsmasq.pid
conf=/etc/u-boot-egress-dnsmasq.conf
upstreams=$(awk '/^nameserver/ && $2 !~ /:/ {print $2}' /etc/resolv.conf.u-boot-orig)

{
  echo "listen-address=127.0.0.1"
  echo "bind-interfaces"
  echo "no-resolv"
  echo "no-hosts"
  echo "pid-file=$pidfile"
  for h in $hosts; do
    for ns in $upstreams; do echo "server=/$h/$ns"; done
    echo "nftset=/$h/4#ip#u_boot_egress#allowed"
  done
} > "$conf.new"

rules=/tmp/u-boot-egress-atomar.nft
{
  echo "add table ip u_boot_egress"
  echo "delete table ip u_boot_egress"
  echo "add table ip u_boot_egress"
  echo "add set ip u_boot_egress allowed { type ipv4_addr; flags timeout; timeout 1h; }"
  echo "add chain ip u_boot_egress output { type filter hook output priority 0; policy drop; }"
  echo "add rule ip u_boot_egress output oifname lo accept"
  echo "add rule ip u_boot_egress output ct state established,related accept"
  echo "add rule ip u_boot_egress output ip daddr @allowed accept"
  for ns in $upstreams; do
    echo "add rule ip u_boot_egress output ip daddr $ns udp dport 53 accept"
    echo "add rule ip u_boot_egress output ip daddr $ns tcp dport 53 accept"
  done
  echo "add table ip6 u_boot_egress6"
  echo "delete table ip6 u_boot_egress6"
  echo "add table ip6 u_boot_egress6"
  echo "add chain ip6 u_boot_egress6 output { type filter hook output priority 0; policy drop; }"
  echo "add rule ip6 u_boot_egress6 output oifname lo accept"
} > "$rules"

if [ -f "$pidfile" ]; then
  old=$(cat "$pidfile")
  kill "$old" 2>/dev/null || true
  # "beendet" heisst: Prozess weg oder Zombie (Z); `kill -0` meldet auch einen Zombie als lebend
  i=0
  while [ "$i" -lt 500 ]; do
    st=$(ps -o stat= -p "$old" 2>/dev/null | cut -c1)
    [ -z "$st" ] || [ "$st" = Z ] && break
    sleep 0.01; i=$((i + 1))
  done
fi
nft -f "$rules"
mv -f "$conf.new" "$conf"
dnsmasq --conf-file="$conf"
echo "nameserver 127.0.0.1" > /etc/resolv.conf
