#!/usr/bin/env bash
# Messung der Eigenschaften von `docker run --env-file`, auf die die Zusage zu
# .devcontainer/devcontainer.env sich stuetzt (slice-v1-devcontainer-env-datei).
# Nur Platzhalterwerte, kein Zugangsdatum. Aufruf: bash messen.sh > ../../../docs/messungen/slice-v1-devcontainer-env-datei/roh.txt
set -u
D="$(mktemp -d)"
trap 'rm -rf "$D"' EXIT
cd "$D" || exit 1
img="${IMG:-busybox:1.36}"
echo "datum: $(date -u +%Y-%m-%dT%H:%MZ)"
echo "engine: $(docker version --format '{{.Server.Version}} {{.Server.Os}}/{{.Server.Arch}}' 2>&1 | head -1)"
echo "image: $img"
echo
echo "== E1 fehlende Datei"
docker run --rm --env-file "$D/fehlt.env" "$img" true 2>&1 | head -1 | sed "s#$D#<tmp>#g"
echo "exit=${PIPESTATUS[0]}"
echo "== E2 Datei nur mit Kommentar"
printf '# nur ein Kommentar\n' > leer.env
docker run --rm --env-file "$D/leer.env" "$img" sh -c 'echo ok' 2>&1 | head -1
echo "exit=${PIPESTATUS[0]}"
echo "== E3 Werte: Leerzeichen, Anfuehrungszeichen, leerer Wert, Name ohne Wert"
printf '# Kommentar\nDEMO_A=eins\nDEMO_B=zwei woerter\nDEMO_Q="mit anfuehrung"\nDEMO_EMPTY=\nDEMO_NOVAL\n' > wert.env
docker run --rm --env-file "$D/wert.env" "$img" sh -c 'echo "A=[$DEMO_A] B=[$DEMO_B] Q=[$DEMO_Q] EMPTY=[$DEMO_EMPTY]"; echo "NOVAL gesetzt: $(env | grep -c "^DEMO_NOVAL=")"' 2>&1 | head -2
echo "== E4 Praefix export"
printf 'export DEMO_X=1\n' > export.env
docker run --rm --env-file "$D/export.env" "$img" true 2>&1 | head -1 | sed "s#$D#<tmp>#g"
echo "exit=${PIPESTATUS[0]}"
echo "== E5 relativer Pfad, Arbeitsverzeichnis ohne die Datei"
mkdir -p rel && printf 'DEMO_R=1\n' > rel/a.env
( cd / && docker run --rm --env-file "rel/a.env" "$img" true 2>&1 | head -1 )
echo "== E6 Vorrang von -e gegen die Datei"
docker run --rm --env-file "$D/wert.env" -e DEMO_A=cli "$img" sh -c 'echo "A=[$DEMO_A]"' 2>&1 | head -1
