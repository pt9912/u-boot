#!/bin/bash
# Messstand: Datum, Host, Docker, Namen der Umgebung (nie Werte), Bestand der Docker-Objekte vorher, Speicher.
# Voraussetzung: 00-vorlauf.sh (Artefakt). Aufruf: WORK=<messordner> tools/messung/ollama-cloud-konto/01-messstand.sh
#   -> $ROH/rohdaten-01.txt, rohdaten-01-vorher.txt
set -u
. "$(dirname "$0")/lib.sh"
vorlauf_pruefen || exit 2
rohdaten rohdaten-01.txt

echo "== Messstand, $(now)"
echo "host: $(uname -sm)"
docker version --format 'docker client={{.Client.Version}} server={{.Server.Version}} os/arch={{.Server.Os}}/{{.Server.Arch}}'
docker info --format 'storage-driver={{.Driver}} cgroup={{.CgroupVersion}} os={{.OperatingSystem}}' 2>&1
echo "-- Namen der Umgebungsvariablen dieses Messlaufs (Werte werden nicht ausgegeben)"
env | cut -d= -f1 | sort | tr '\n' ' '; echo
umgebung_pruefen; rc=$?
echo "umgebung-pruefung exit=$rc (0 = kein verbotener Name)"
aussage "umgebung-ohne-verbotene-namen=$([ $rc = 0 ] && echo ja || echo nein)"
echo "-- Bestand der Docker-Objekte vorher"
bestand | tee "$ROH/rohdaten-01-vorher.txt"
echo "-- Container im Netz bridge (Anzahl; Namen fremder Container werden nicht ausgegeben)"
echo "bridge: gateway=$(docker network inspect bridge --format '{{range .IPAM.Config}}{{.Gateway}}{{end}}') container=$(docker network inspect bridge --format '{{len .Containers}}')"
echo "-- Platz"
docker system df | head -5
free -m | head -2
prompt_stand
