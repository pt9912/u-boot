#!/bin/bash
# Laeuft im Container als vscode, OHNE sudo (ausser dem Abschnitt "Zaehler"); $1 = Modus (ist|skuid), nur fuer die Beschriftung.
# D2: Wirkung der Regel auf den direkten Upstream-Pfad und auf erlaubte Namen. Siehe 04-d2.sh.
set -u
mode=${1:?ist|skuid}
UPS=$(awk '/^nameserver/ && $2 !~ /:/ {print $2}' /etc/resolv.conf.u-boot-orig)
hosts=$(sed -n 's/^hosts="\(.*\)"$/\1/p' /usr/local/bin/u-boot-egress-init)
echo "id: $(id)"
echo "Upstream-Liste (resolv.conf.u-boot-orig): $UPS"
echo "resolv.conf: $(grep -v '^#' /etc/resolv.conf | tr '\n' ' ')"
echo "Liste der erlaubten Namen: $hosts"

# q <beschriftung> <dig-argumente...> : Ergebnis der Antwort (Adresse) oder "keine Antwort"
q() {
  local label=$1; shift
  local out rc
  out=$(dig +time=3 +tries=1 +short "$@" 2>&1); rc=$?
  out=$(printf '%s\n' "$out" | grep -v -e '^;' -e '^$' | head -2 | tr '\n' ' ')
  if [ "$rc" = 0 ] && [ -n "$out" ]; then echo "antwort"; echo "    $label: antwort $out" >&2
  elif [ "$rc" = 9 ]; then echo "keine"; echo "    $label: keine Antwort (dig rc=9, Zeitueberschreitung)" >&2
  else echo "leer"; echo "    $label: dig rc=$rc, keine Adresse ($out)" >&2; fi
}
code() { curl -sS -m 5 -o /dev/null -w '%{http_code}' "$@" 2>/dev/null; true; }

echo "--- 1. direkt an jeden Upstream (UDP und TCP), Name ausserhalb der Liste (example.com)"
for ns in $UPS; do
  for p in udp tcp; do
    r=$(q "@$ns $p example.com" $([ $p = tcp ] && echo +tcp) @"$ns" example.com A 2>&2)
    echo "AUSSAGE $mode upstream-direkt $ns $p example.com=$r"
  done
done
echo "--- 1b. direkt an den ersten Upstream, Name der Liste (api.github.com)"
ns1=$(echo $UPS | cut -d' ' -f1)
r=$(q "@$ns1 udp api.github.com" @"$ns1" api.github.com A 2>&2); echo "AUSSAGE $mode upstream-direkt $ns1 udp api.github.com=$r"

echo "--- 2. ueber den lokalen Resolver (127.0.0.1, dnsmasq)"
r=$(q "@127.0.0.1 example.com" @127.0.0.1 example.com A 2>&2); echo "AUSSAGE $mode lokal example.com=$r"
echo "    Status @127.0.0.1 example.com: $(dig +time=3 +tries=1 @127.0.0.1 example.com A 2>&1 | grep -o 'status: [A-Z]*' | head -1)"
for h in $hosts; do
  r=$(q "@127.0.0.1 $h" @127.0.0.1 "$h" A 2>&2); echo "AUSSAGE $mode lokal $h=$r"
done
echo "--- 3. ueber getent (resolv.conf) wie ein Programm"
for h in api.github.com deb.debian.org example.com; do
  s=$(date +%s%N); a=$(getent hosts "$h" | head -1 | cut -d' ' -f1); e=$(date +%s%N)
  echo "    getent $h: ${a:-keine Adresse} ($(( (e - s) / 1000000 )) ms)"
  echo "AUSSAGE $mode getent $h=$([ -n "$a" ] && echo adresse || echo keine)"
done
echo "--- 4. Verbindung"
echo "AUSSAGE $mode https api.github.com=$(code https://api.github.com)"
echo "AUSSAGE $mode https deb.debian.org=$(code https://deb.debian.org)"
echo "AUSSAGE $mode https example.com (Name)=$(code https://example.com)"
ip=$(dig +time=3 +tries=1 +short @"$ns1" example.com A 2>/dev/null | grep -E '^[0-9.]+$' | head -1)
if [ -n "$ip" ]; then
  echo "AUSSAGE $mode https zur direkt aufgeloesten Adresse von example.com=$(code --resolve "example.com:443:$ip" https://example.com)"
else
  echo "    (keine direkt aufgeloeste Adresse: Verbindungsprobe entfaellt)"
  echo "AUSSAGE $mode https zur direkt aufgeloesten Adresse von example.com=entfaellt"
fi
echo "--- 5. Zaehler der Messvariante (sudo nur zum Lesen der Regeln)"
sudo -n nft list table ip u_boot_egress 2>&1 | grep -E 'counter|skuid|dport 53' | sed 's/^[[:space:]]*/    /'
