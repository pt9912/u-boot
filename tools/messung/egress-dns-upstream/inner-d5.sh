#!/bin/bash
# Laeuft im Container als vscode, OHNE sudo; $1 = Modus (ist|skuid), $2 = Adresse R des Antwort-Containers, der hier der Upstream des Containers ist.
# D5: Datenkanal ueber erlaubte Namen. Siehe 08-d5.sh.
set -u
mode=${1:?ist|skuid}
R=${2:?adresse}
A() { echo "AUSSAGE $mode d5 $1=$2"; }
echo "Upstream-Liste: $(awk '/^nameserver/ {print $2}' /etc/resolv.conf.u-boot-orig | tr '\n' ' ')"
echo "--- 1. direkt an den Upstream R, Name ausserhalb der Liste"
o=$(dig +time=2 +tries=1 +short @"$R" direkt.example.com A 2>/dev/null | grep -v -e '^;' -e '^$' | head -1); A "direkt-an-upstream" "$([ -n "$o" ] && echo antwort || echo keine)"
echo "--- 2. ueber den lokalen Resolver: Untername eines erlaubten Namens (kanal1.api.github.com) und Name ausserhalb der Liste (kanal1.example.com)"
o=$(dig +time=3 +tries=1 +short @127.0.0.1 kanal1.api.github.com A 2>/dev/null | grep -v -e '^;' -e '^$' | head -1); echo "    kanal1.api.github.com: ${o:-keine Antwort}"
A "untername-erlaubter-name" "$([ -n "$o" ] && echo antwort || echo keine)"
o=$(dig +time=3 +tries=1 +short @127.0.0.1 kanal1.example.com A 2>/dev/null | grep -v -e '^;' -e '^$' | head -1); echo "    kanal1.example.com: ${o:-keine Antwort}"
A "untername-unlisteter-name" "$([ -n "$o" ] && echo antwort || echo keine)"
echo "--- 3. Durchsatz: 5 s nacheinander Anfragen mit drei Labels zu je 63 Hex-Zeichen vor api.github.com an den lokalen Resolver (MESSWERT, zeitabhaengig)"
echo "MESSWERT $mode lokal-api.github.com: $(perl /m/dnsq.pl 127.0.0.1 api.github.com 5 3 63)"
echo "--- 4. dasselbe, aber ein Name ausserhalb der Liste (example.com) an den lokalen Resolver"
echo "MESSWERT $mode lokal-example.com: $(perl /m/dnsq.pl 127.0.0.1 example.com 3 3 63)"
echo "--- 5. Vergleich: direkt an den Upstream R (nur im Fall ist offen)"
echo "MESSWERT $mode direkt-R: $(perl /m/dnsq.pl "$R" example.com 5 3 63)"
