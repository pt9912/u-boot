#!/bin/bash
# Laeuft im Container als vscode; $1 = ohne | mit  (sudo). Siehe 02-m1-rechte.sh.
set -u
phase=${1:?ohne|mit}
http() { # http <label> <url> [curl-args...]
  label=$1; url=$2; shift 2
  out=$(curl -sS -m 5 -o /dev/null -w '%{http_code}' "$@" "$url" 2>&1 | tr '\n' ' ' | cut -c1-120)
  echo "$label -> $out"
}
upstream=$(awk '/^nameserver/ && $2 !~ /:/ {print $2; exit}' /etc/resolv.conf.u-boot-orig)

if [ "$phase" = ohne ]; then
  echo "id: $(id)"
  grep -E '^Cap(Prm|Eff|Bnd)' /proc/self/status
  sudo -n true; echo "sudo -n true rc=$? (0 = passwortloses sudo)"
  sudo -n -l 2>&1 | tail -2
  echo "--- nft ohne sudo"
  nft list tables 2>&1 | head -2
  echo "--- Schreibversuch (Anhaengen ohne Inhalt) auf die Dateien der Beschraenkung"
  for f in /usr/local/bin/u-boot-egress-init /etc/u-boot-egress-dnsmasq.conf /etc/resolv.conf /etc/resolv.conf.u-boot-orig /run/u-boot-egress-dnsmasq.pid; do
    ls -l "$f" | awk '{print $1, $3, $4, $9}'
    ( : >> "$f" ) 2>&1 | sed 's/^/    /'
    if ( : >> "$f" ) 2>/dev/null; then echo "    $f: schreibbar"; else echo "    $f: NICHT schreibbar"; fi
  done
  echo "--- dnsmasq: Prozess und Signal"
  ps -o user,pid,args -C dnsmasq
  pid=$(cat /run/u-boot-egress-dnsmasq.pid); echo "    PID laut Pid-Datei: $pid"
  kill -TERM "$pid" 2>&1 | sed 's/^/    /'; sleep 1
  if ps -p "$pid" >/dev/null; then echo "    dnsmasq $pid lebt nach SIGTERM"; else echo "    dnsmasq $pid beendet"; fi
  ps -o user,pid,args -C dnsmasq | tail -n +2
  echo "--- Statusdatei"
  ls -ld /var/lib/u-boot/status | awk '{print $1, $3, $4, $9}'
  ls -l /var/lib/u-boot/status | tail -n +2 | awk '{print $1, $3, $4, $9}'
  if ( : >> /var/lib/u-boot/status/egress ) 2>/dev/null; then echo "    status/egress: schreibbar durch vscode"; else echo "    status/egress: NICHT schreibbar"; fi
  echo "--- Statusdatei ersetzen (das Verzeichnis gehoert vscode: loeschen und neu anlegen)"
  rm -f /var/lib/u-boot/status/egress && printf 'u-boot-status 1\nresult=ok\nat=2000-01-01T00:00:00Z\n' > /var/lib/u-boot/status/egress; echo "ersetzen rc=$?"
  ls -l /var/lib/u-boot/status | tail -n +2 | awk '{print $1, $3, $4, $9}'
  cat /var/lib/u-boot/status/egress
  echo "--- Verkehr unter der Beschraenkung (ohne sudo)"
  http "erlaubt api.github.com" https://api.github.com
  http "gesperrt example.com (Name)" https://example.com
  http "gesperrt 1.1.1.1 (IP-Literal)" https://1.1.1.1
  echo "--- DNS direkt an den Upstream-Resolver $upstream (ohne sudo)"
  dig +time=3 +tries=1 @"$upstream" example.com A +short | head -3 | sed 's/^/    A /'
  ip=$(dig +time=3 +tries=1 +short @"$upstream" example.com A | head -1)
  http "Verbindung zur aufgeloesten Adresse $ip" https://example.com --resolve "example.com:443:$ip"
  echo "--- nft nach den Versuchen (ohne sudo nicht lesbar; mit sudo):"
  sudo nft list tables
else
  sudo -n true; echo "sudo -n true rc=$?"
  echo "--- nft mit sudo"
  sudo nft list tables
  echo "--- Tabellen entfernen"
  sudo nft delete table ip u_boot_egress; echo "delete ip rc=$?"
  sudo nft delete table ip6 u_boot_egress6; echo "delete ip6 rc=$?"
  sudo nft list tables | sed 's/^/    verbleibend: /'
  http "gesperrt 1.1.1.1 (IP-Literal), Tabellen weg" https://1.1.1.1
  http "gesperrt example.com (Name), Tabellen weg, resolv.conf unveraendert" https://example.com
  echo "--- resolv.conf auf die Original-Datei zuruecksetzen"
  sudo cp /etc/resolv.conf.u-boot-orig /etc/resolv.conf; echo "cp rc=$?"
  http "gesperrt example.com (Name), Tabellen weg, resolv.conf zurueck" https://example.com
fi
