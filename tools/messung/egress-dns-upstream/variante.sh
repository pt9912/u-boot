#!/bin/sh
# MESSVARIANTE, kein Produktcode und kein Vorschlag fuer einen Wortlaut. Sie liegt nur im Wegwerf-Container der Messung
# (Mount /m) und aendert weder das Template noch das erzeugte Startscript.
#
# Aufruf (als root im Container): sh /m/variante.sh <ist|skuid> [zaehler]
#   ist     fuehrt das erzeugte Startscript /usr/local/bin/u-boot-egress-init unveraendert aus
#   skuid   fuehrt eine Kopie davon aus, in der die beiden Regeln zum Upstream-DNS (udp und tcp, dport 53) zusaetzlich
#           `meta skuid $DNSUID` verlangen; die Kopie entsteht mit sed aus dem erzeugten Script (DNSUID Vorgabe 65534 = nobody,
#           der Benutzer des dnsmasq laut Messung D1)
#   zaehler haengt Zaehler-Regeln an (nur Messung): vor den Upstream-Regeln zaehlt eine nicht terminierende Regel jeden Versuch einer
#           DNS-Anfrage (dport 53, nicht ueber lo) eines anderen Benutzers als $DNSUID, nach ihnen eine Regel jedes Paket an Port 53,
#           das bis dorthin nicht angenommen wurde ("verworfen")
set -eu
mode=${1:?ist|skuid}
count=${2:-}
DNSUID=${DNSUID:-65534}
src=/usr/local/bin/u-boot-egress-init
case "$mode" in
  ist) script=$src ;;
  skuid)
    script=/tmp/egress-init-skuid.sh
    sed -e "s|output ip daddr \"\$ns\" udp dport 53 accept|output meta skuid $DNSUID ip daddr \"\$ns\" udp dport 53 accept|" \
        -e "s|output ip daddr \"\$ns\" tcp dport 53 accept|output meta skuid $DNSUID ip daddr \"\$ns\" tcp dport 53 accept|" \
        "$src" > "$script"
    if [ "$(grep -c "meta skuid $DNSUID" "$script")" != 2 ]; then echo "variante: sed hat nicht genau zwei Regeln ersetzt" >&2; exit 2; fi
    ;;
  *) echo "variante: unbekannter Modus $mode" >&2; exit 2 ;;
esac
sh "$script"
if [ "$count" = zaehler ]; then
  nft insert rule ip u_boot_egress output oifname != lo meta l4proto { tcp, udp } th dport 53 meta skuid != "$DNSUID" counter
  nft insert rule ip u_boot_egress output oifname != lo meta l4proto { tcp, udp } th dport 53 meta skuid "$DNSUID" counter
  nft add rule ip u_boot_egress output meta l4proto { tcp, udp } th dport 53 counter drop
fi
