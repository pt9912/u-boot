#!/bin/bash
# Messstand: Datum, Host, Docker, Umgebung (Namen, nie Werte) samt Gegenprobe der Pruefung, Bestand der Docker-Objekte vorher.
# Aufruf: WORK=<messordner> tools/messung/ollama-cloud/01-messstand.sh      -> $OUT/rohdaten-01.txt, rohdaten-01-vorher.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-01.txt

echo "== Messstand, $(now)"
echo "host: $(uname -sm)"
docker version --format 'docker client={{.Client.Version}} server={{.Server.Version}} os/arch={{.Server.Os}}/{{.Server.Arch}}'
docker info --format 'storage-driver={{.Driver}} cgroup={{.CgroupVersion}}' 2>&1
echo "-- Namen der Umgebungsvariablen dieses Messlaufs (nach env -i; Werte werden nicht ausgegeben)"
env | cut -d= -f1 | sort | tr '\n' ' '; echo
umgebung_pruefen; rc=$?
echo "umgebung-pruefung exit=$rc (0 = kein verbotener Name)"
aussage "umgebung-ohne-verbotene-namen=$([ $rc = 0 ] && echo ja || echo nein)"
echo "-- Gegenprobe der Pruefung: derselbe Aufruf mit gesetzten Namen (Platzhalterwert) muss anschlagen"
p=platzhalter
( ANTHROPIC_API_KEY=$p OLLAMA_API_KEY=$p GH_TOKEN=$p GITLAB_TOKEN=$p GIT_TOKEN=$p umgebung_pruefen; echo "gegenprobe exit=$?" )
echo "-- die Variablen des Hosts, die die Dev-Container-CLI ueber \${localEnv:…} erreichen koennte, sind im Messlauf nicht gesetzt:"
for n in ANTHROPIC_API_KEY GH_TOKEN GITLAB_TOKEN GIT_TOKEN OLLAMA_API_KEY OLLAMA_HOST; do
  printf '%s: %s\n' "$n" "$([ -n "$(printenv "$n")" ] && echo GESETZT || echo ungesetzt)"
done
echo "-- Bestand der Docker-Objekte vorher"
bestand | tee "$OUT/rohdaten-01-vorher.txt"
echo "-- Platz"
docker system df | head -5
free -m | head -2
