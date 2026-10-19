#!/bin/bash
# Laeuft im Container als vscode, OHNE sudo (ausser Abschnitt 9 und dem Eintrag der Adresse in das Set); $1 = Modus (ist|skuid), $2 = Adresse des
# Antwort-Containers R (inner-resp.sh). D3: Wege um die Regel. Siehe 06-d3.sh.
# Jede Probe nennt das Ergebnis "offen" (Antwort oder Verbindung kam zustande) oder "gesperrt" (Zeitueberschreitung/Fehler).
set -u
mode=${1:?ist|skuid}
R=${2:?adresse des antwort-containers}
orig=/etc/resolv.conf.u-boot-orig; [ -f "$orig" ] || orig=/etc/resolv.conf   # Fall ohne: kein Startscript, die Datei des Containers
UPS=$(awk '/^nameserver/ && $2 !~ /:/ {print $2}' "$orig")
ns1=$(echo $UPS | cut -d' ' -f1)
echo "Upstream-Liste: $(echo $UPS); Antwort-Container R=$R (nicht in der Upstream-Liste, nicht im Set)"
ans() { # <beschriftung> <dig-argumente...>
  local label=$1; shift
  local out rc
  out=$(dig +time=2 +tries=1 +short "$@" 2>&1); rc=$?
  out=$(printf '%s\n' "$out" | grep -v -e '^;' -e '^$' | head -1)
  if [ "$rc" = 0 ] && [ -n "$out" ]; then echo "offen"; echo "    $label: offen ($out)" >&2
  else echo "gesperrt"; echo "    $label: gesperrt (dig rc=$rc)" >&2; fi
}
tcp() { # <beschriftung> <adresse> <port>: Verbindung aufbauen und ein Byte-Paket senden
  local label=$1 a=$2 p=$3 rc
  printf 'u-boot-messung-%s' "$label" | nc -w 2 "$a" "$p" > /dev/null 2>&1; rc=$?
  if [ "$rc" = 0 ]; then echo "offen"; echo "    $label: offen (nc rc=0)" >&2; else echo "gesperrt"; echo "    $label: gesperrt (nc rc=$rc)" >&2; fi
}
pg() { # <beschriftung> <adresse>
  if ping -c 1 -W 2 "$2" > /dev/null 2>&1; then echo "offen"; echo "    $1: offen (Echo-Antwort)" >&2; else echo "gesperrt"; echo "    $1: gesperrt (keine Echo-Antwort)" >&2; fi
}
A() { echo "AUSSAGE $mode d3 $1=$2"; }

echo "--- 1. Grundfall: UDP und TCP an den ersten Upstream $ns1 (Name example.com)"
A "udp53-upstream" "$(ans "udp53 @$ns1" @"$ns1" example.com A 2>&2)"
A "tcp53-upstream" "$(ans "tcp53 @$ns1" +tcp @"$ns1" example.com A 2>&2)"
echo "--- 2. anderer Upstream: R (Port 53, UDP und TCP), nicht in der Upstream-Liste"
A "udp53-anderer-resolver-R" "$(ans "udp53 @R" @"$R" example.com A 2>&2)"
A "tcp53-anderer-resolver-R" "$(ans "tcp53 @R" +tcp @"$R" example.com A 2>&2)"
echo "--- 3. andere Portnummer und Port 853 / 443 am Upstream $ns1 (TCP; UDP 5353)"
A "tcp853-upstream" "$(tcp "tcp853 $ns1" "$ns1" 853 2>&2)"
A "tcp443-upstream" "$(tcp "tcp443 $ns1" "$ns1" 443 2>&2)"
A "udp5353-upstream" "$(ans "udp5353 @$ns1" -p 5353 @"$ns1" example.com A 2>&2)"
echo "--- 4. ICMP"
A "icmp-upstream" "$(pg "ping $ns1" "$ns1" 2>&2)"
A "icmp-R" "$(pg "ping R" "$R" 2>&2)"
echo "--- 5. IPv6"
echo "    IPv6-Adressen im Container (ausser ::1): $(ip -6 addr | grep -c 'inet6 [^:]' || true); IPv6-Routen: $(ip -6 route | wc -l)"
A "ipv6-extern-udp53" "$(ans "udp53 @2001:4860:4860::8888" @2001:4860:4860::8888 example.com A 2>&2)"
echo "--- 6. R als erlaubte Adresse (Messvorbereitung: die Adresse wird per sudo in das Set eingetragen, wie der dnsmasq es fuer einen erlaubten Namen taete)"
sudo -n nft add element ip u_boot_egress allowed "{ $R }"; echo "    eintrag rc=$?"
A "udp53-erlaubte-adresse-R" "$(ans "udp53 @R (im Set)" @"$R" example.com A 2>&2)"
A "tcp53-erlaubte-adresse-R" "$(ans "tcp53 @R (im Set)" +tcp @"$R" example.com A 2>&2)"
A "udp5353-erlaubte-adresse-R" "$(ans "udp5353 @R (im Set)" -p 5353 @"$R" example.com A 2>&2)"
A "tcp853-erlaubte-adresse-R" "$(tcp "tcp853 R" "$R" 853 2>&2)"
A "tcp443-erlaubte-adresse-R" "$(tcp "tcp443 R" "$R" 443 2>&2)"
A "icmp-erlaubte-adresse-R" "$(pg "ping R (im Set)" "$R" 2>&2)"
echo "--- 7. Zaehler und Regeln am Ende"
sudo -n nft list chain ip u_boot_egress output | sed 's/^[[:space:]]*/    /'
echo "--- 8. Docker-eigener Resolver 127.0.0.11 (bei benutzerdefiniertem Netz); hier Standard-Bruecke:"
echo "    resolv.conf-Original: $(grep -v '^#' /etc/resolv.conf.u-boot-orig | grep nameserver | tr '\n' ' ')"
A "udp53-127.0.0.11" "$(ans "udp53 @127.0.0.11" @127.0.0.11 example.com A 2>&2)"
echo "--- 9. mit sudo (Grenze: Guardrail; die sudoers-Regel des Images erlaubt nur (root) NOPASSWD: ALL)"
echo "    sudoers: $(sudo -n -l 2>&1 | tail -1 | sed 's/^[[:space:]]*//')"
a=$(sudo -n dig +time=2 +tries=1 +short @"$ns1" example.com A 2>/dev/null | grep -c -E '^[0-9a-f.:]+$')
A "sudo-root-udp53-upstream" "$([ "$a" -gt 0 ] && echo offen || echo gesperrt)"
b=$(sudo -n setpriv --reuid=65534 --regid=65534 --clear-groups dig +time=2 +tries=1 +short @"$ns1" example.com A 2>/dev/null | grep -c -E '^[0-9a-f.:]+$')
A "sudo-setpriv-nobody-udp53-upstream" "$([ "$b" -gt 0 ] && echo offen || echo gesperrt)"
