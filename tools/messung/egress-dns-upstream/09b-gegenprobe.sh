#!/bin/bash
# Gegenprobe der Pruefung in 09-d4.sh und des Befunds "erlaubte Namen loesen weiter auf": die Messvariante mit einer UID, die nicht die des dnsmasq
# ist (DNSUID=1234), sperrt den dnsmasq selbst. Die Pruefung von 09-d4.sh muss den Zustand als auffaellig melden, die lokale Aufloesung
# eines erlaubten Namens muss ausfallen.
# Aufruf: WORK=<messordner> OUT=<ausgabeverzeichnis> 09b-gegenprobe.sh -> $OUT/rohdaten-09b-gegenprobe.txt
# Voraussetzung: 02-projekte.sh.
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-09b-gegenprobe.txt
c=$PFX-gegen
echo "######## Gegenprobe: Variante mit falscher UID (DNSUID=1234), $(now)"
start $c $IMG_TOOL
docker exec -u vscode $c sudo env DNSUID=1234 sh /m/variante.sh skuid; echo "postStart (skuid, DNSUID=1234) exit=$?"
z=$(docker exec -u root $c bash /m/inner-d4.sh skuid gegenprobe-falsche-uid | head -1)
echo "$z"
erw="direkt-example.com=keine"; erws="skuid-regeln=4"
case "$z" in *"dnsmasq-prozesse=1 uid=65534, pid-lebt=ja"*"$erws"*"lokal-api.github.com='status: NOERROR'"*"$erw"*) auf=nein ;; *) auf=ja ;; esac
echo "Pruefung aus 09-d4.sh meldet den Zustand als auffaellig: $auf"
aussage "skuid gegenprobe falsche-uid pruefung-auffaellig=$auf"
aussage "skuid gegenprobe falsche-uid lokal-aufloesung=$(case "$z" in *"status: NOERROR"*) echo noerror ;; *) echo ausgefallen ;; esac)"
docker rm -f -v $c > /dev/null
