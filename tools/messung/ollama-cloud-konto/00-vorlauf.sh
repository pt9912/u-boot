#!/bin/bash
# Sicherung vor dem ersten Lauf (DoD 1, Verification-Report V-3): Gegenprobe der Schwaerzung (redact-probe.sh), Pruefung der Umgebung nach Namen samt
# Gegenprobe, Hash des Stands von lib.sh und redact-probe.sh. Das Ergebnis ist ein Artefakt mit Zeitstempel (Stufe 1: $ROH/rohdaten-00-vorlauf.txt);
# jedes weitere Skript ruft vorlauf_pruefen (lib.sh) und bricht ab, wenn das Artefakt fehlt, nicht bestanden ist oder zu einem anderen Stand gehoert.
# Aufruf: WORK=<messordner> tools/messung/ollama-cloud-konto/00-vorlauf.sh
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-00-vorlauf.txt

echo "== Vorlauf (Sicherung vor dem ersten Lauf), $(now)"
echo "-- Namen der Umgebungsvariablen dieses Skripts (nach der Bereinigung im selben Prozess; Werte werden nicht ausgegeben)"
env | cut -d= -f1 | sort | tr '\n' ' '; echo
umgebung_pruefen; rc_env=$?
echo "umgebung-pruefung exit=$rc_env (0 = kein verbotener Name)"
echo "-- Gegenprobe der Pruefung: derselbe Aufruf mit gesetzten Namen (Platzhalterwert) muss anschlagen"
p=platzhalter
( ANTHROPIC_API_KEY=$p OLLAMA_API_KEY=$p GH_TOKEN=$p GITLAB_TOKEN=$p GIT_TOKEN=$p ANTHROPIC_AUTH_TOKEN=$p OLLAMA_HOST=$p umgebung_pruefen; echo "gegenprobe exit=$?" )
echo "-- die Variablen des Hosts, die einen Zugang tragen koennten, sind im Messlauf nicht gesetzt:"
for n in ANTHROPIC_API_KEY ANTHROPIC_AUTH_TOKEN GH_TOKEN GITLAB_TOKEN GIT_TOKEN OLLAMA_API_KEY OLLAMA_HOST; do
  printf '%s: %s\n' "$n" "$([ -n "$(printenv "$n")" ] && echo GESETZT || echo ungesetzt)"
done
echo "-- redact-probe.sh"
"$here/redact-probe.sh"; rc_probe=$?
echo "redact-probe exit=$rc_probe"
h=$(stand_hash)
echo "-- Stand der Sicherung"
aussage "sicherung-stand=$h"
[ "$rc_env" = 0 ] && aussage "umgebung-pruefung=ok" || aussage "umgebung-pruefung=FEHLER"
[ "$rc_probe" = 0 ] && aussage "redact-probe=ok" || aussage "redact-probe=FEHLER"
[ "$rc_env" = 0 ] && [ "$rc_probe" = 0 ]
