#!/bin/bash
# D3: Wege um die Regel (UDP/TCP, anderer Resolver, andere Portnummer, 853/443, ICMP, IPv6, erlaubte Adresse, Docker-eigener Resolver, sudo).
# Drei Faelle je Sonde: ohne = Kontrolle ohne Startscript (die Sonde erkennt "offen"), ist = erzeugtes Startscript, skuid = Messvariante.
# Aufruf: WORK=<messordner> OUT=<ausgabeverzeichnis> 06-d3.sh -> $OUT/rohdaten-06-d3.txt
# Voraussetzung: 02-projekte.sh. Das Gegenueber ist ein Container der Messung (inner-resp.sh) im Standardnetz, kein fremder Server.
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-06-d3.txt
r=$PFX-resp
c=$PFX-d3
docker rm -f -v $r > /dev/null 2>&1
docker run -d --init --name $r -u root --entrypoint sleep $IMG_TOOL infinity > /dev/null
docker cp "$here/inner-resp.sh" $r:/tmp/inner-resp.sh
docker exec -u root $r sh /tmp/inner-resp.sh
R=$(docker inspect -f '{{.NetworkSettings.Networks.bridge.IPAddress}}' $r)
echo "Antwort-Container R=$R"
for mode in ohne ist skuid; do
  echo; echo "######## D3, Fall $mode, $(now)"
  start $c $IMG_TOOL
  case "$mode" in
    ohne) echo "(kein Startscript: Kontrolle, die Sonden erkennen einen offenen Pfad)" ;;
    *)    egress $c $mode zaehler ;;
  esac
  docker exec -u vscode $c bash /m/inner-d3.sh "$mode" "$R"
  docker rm -f -v $c > /dev/null
done
echo; echo "######## Eingang am Antwort-Container (Anfragen und Nutzdaten, die ankamen)"
for p in 53 5353; do echo "dns$p.log: Anfragen (query[...]) gesamt=$(docker exec $r grep -c 'query\[' /tmp/dns$p.log)"; done
for p in 853 443; do echo "tcp$p.in: $(docker exec $r sh -c "cat /tmp/tcp$p.in 2>/dev/null | wc -c") Bytes, davon Zeichenketten u-boot-messung-: $(docker exec $r sh -c "grep -o 'u-boot-messung-[a-z0-9 ]*' /tmp/tcp$p.in 2>/dev/null | sort | uniq -c | tr '\n' ';'")"; done
docker rm -f -v $r > /dev/null
