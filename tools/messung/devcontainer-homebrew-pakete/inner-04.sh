#!/bin/bash
# Laeuft IM Container (Image einer Variante) als Benutzer vscode, wie im Devcontainer; aufgerufen von 04-egress.sh.
# Aufruf im Container: inner-04.sh <modus>
#   modus   discover  keine Restriktion; ein protokollierender dnsmasq reicht alle Namen weiter (zeigt, welche Hosts `brew install` braucht)
#           restricted das erzeugte Startscript des Images (`sudo sh /usr/local/bin/u-boot-egress-init`); danach ersetzt ein protokollierender
#                      dnsmasq mit derselben Konfiguration den Resolver des Skripts (Verfahren von tools/messung/forge-cli)
#           restricted-ohne  wie restricted, aber die drei Hosts der Zeile des Features aus der Zeile hosts= des Startscripts entfernt (die Liste ohne
#                      die Zeile des Features; das Startscript ist sonst unveraendert)
#           restricted-formelhost  wie restricted, aber die Zeile hosts= traegt zusaetzlich release-assets.githubusercontent.com (den Host, von dem die
#                      gemessene Tap-Formel laedt; der Nutzer traegt ihn in egress.allow ein)
# Die Installationen laufen mit der Umgebung eines Nutzers (keine Schalter): die Automatik der Aktualisierung und die Analytik von Homebrew
# sind an. Kein Zugangsdatum.
set -u
mode=$1
log=/tmp/dns.log
now() { date -u +%T.%N | cut -c1-12; }
mark() { echo "---- $(now) $*"; }
sudo touch "$log"; sudo chmod 666 "$log"
up=$(awk '/^nameserver/ && $2 !~ /:/ {print $2}' /etc/resolv.conf | head -1)

if [ "$mode" = discover ]; then
  printf 'listen-address=127.0.0.1\nbind-interfaces\nno-resolv\nno-hosts\nserver=%s\nlog-queries\nlog-facility=%s\n' "$up" "$log" | sudo tee /tmp/dns-discover.conf > /dev/null
  sudo dnsmasq --conf-file=/tmp/dns-discover.conf
  echo "nameserver 127.0.0.1" | sudo tee /etc/resolv.conf > /dev/null
else
  if [ "$mode" = restricted-ohne ]; then
    hosts=$(grep "^hosts=" /usr/local/bin/u-boot-egress-init | cut -d'"' -f2)
    rest=$(echo "$hosts" | tr " " "\n" | grep -v -x -e formulae.brew.sh -e ghcr.io -e pkg-containers.githubusercontent.com | tr "\n" " ")
    sudo sed -i "s|^hosts=.*|hosts=\"$rest\"|" /usr/local/bin/u-boot-egress-init
  fi
  if [ "$mode" = restricted-formelhost ]; then
    hosts=$(grep "^hosts=" /usr/local/bin/u-boot-egress-init | cut -d'"' -f2)
    sudo sed -i "s|^hosts=.*|hosts=\"$hosts release-assets.githubusercontent.com\"|" /usr/local/bin/u-boot-egress-init
  fi
  sudo sh /usr/local/bin/u-boot-egress-init; echo "egress-init exit=$?"
  echo "-- Status-Datei des Skripts"; cat /var/lib/u-boot/status/egress
  sudo kill "$(cat /run/u-boot-egress-dnsmasq.pid)"; sleep 1
  { cat /etc/u-boot-egress-dnsmasq.conf; echo "log-queries"; echo "log-facility=$log"; } | sudo tee /tmp/dns-restricted.conf > /dev/null
  sudo dnsmasq --conf-file=/tmp/dns-restricted.conf
fi
sleep 1
hostsline=$(grep '^hosts=' /usr/local/bin/u-boot-egress-init 2>/dev/null | cut -d'"' -f2)
echo "modus=$mode user=$(id -un)"
[ "$mode" != discover ] && echo "Hosts der Allowlist im Startscript: $hostsline"

# probe <etikett> <befehl...>: fuehrt aus (Zeitlimit 150 s), zeigt das Ende der Ausgabe, den Exit-Code und die DNS-Namen waehrend des Aufrufs.
probe() {
  local label=$1; shift
  local n0 out rc
  n0=$(wc -l < "$log"); out=$(mktemp); sleep 0.2
  mark "$label"
  echo "\$ $*"
  timeout 150 "$@" > "$out" 2>&1 < /dev/null; rc=$?
  sed -e 's/\x1b\[[0-9;?]*[a-zA-Z]//g' "$out" | tr '\r' '\n' | grep -v -E '^\s*$|^#|%\s*$' | tail -8 | cut -c1-180
  rm -f "$out"
  echo "exit=$rc"
  echo "AUSSAGE $mode-$label-exit=$rc"
  echo "-- DNS-Anfragen (A/AAAA, eindeutig)"
  tail -n +$((n0 + 1)) "$log" | grep -E 'query\[(A|AAAA)\]' | awk '{print $6}' | sort | uniq -c
  if [ "$mode" != discover ]; then
    echo "-- Antworten ausserhalb der Liste (REFUSED/NXDOMAIN/NODATA/config)"
    tail -n +$((n0 + 1)) "$log" | grep -E 'REFUSED|NXDOMAIN|NODATA|config ' | awk '{print $4, $5, $6}' | sort -u | head -8
  fi
}

probe version brew --version
probe core-formel brew install --formula jq
probe tap-formel brew install --formula pt9912/u-boot/u-boot
probe aufruf-tap-formel u-boot --version
echo
mark "alle DNS-Anfragen der Sitzung (eindeutig, A/AAAA)"
grep -E 'query\[(A|AAAA)\]' "$log" | awk '{print $6}' | sort | uniq -c
