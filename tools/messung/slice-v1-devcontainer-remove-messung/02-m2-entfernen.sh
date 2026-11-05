#!/bin/bash
# M2: Wie verhaelt sich `docker rm` bei einem laufenden Container (ohne und mit -f), und bleibt ein Volume, das ein zweiter Container
# haelt, beim `docker volume rm` stehen?
#   a  laufender Container: docker rm, docker rm -f; Volume danach
#   b  Volume, das ein zweiter (Fremd-)Container haelt: docker volume rm
#   c  gestoppter Container (u-boot devcontainer down): docker volume rm [-f], docker rm
#   d  u-boot nach dem Entfernen: down, up; Inhalt des Workspace-Volumes bleibt / ist weg
#   e  anonyme Volumes: docker rm -f ohne und mit -v (Probe-Container und Projekt mit VOLUME im Dockerfile), up --recreate
#   f  u-boot devcontainer up --recreate mit einem benannten Volume (Workspace-Volume, Marker)
# Projekt sbx (Sandbox-Profil, ein Workspace-Volume) und anon (Standardprofil, VOLUME im Dockerfile); Zweitcontainer: alpine:3.
# Aufruf: UBOOT=<binary> WORK=<messordner> ./02-m2-entfernen.sh      -> rohdaten-m2.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-m2.txt
echo "== M2 docker rm und docker volume rm, $(now)"
anfang

existiert() { docker volume inspect "$1" > /dev/null 2>&1 && echo vorhanden || echo entfernt; }
ws_volume() { mount_volumes "$1" | grep -- '-workspace-' | head -n 1; }
marker_lesen() { # <volume> -> Inhalt der Markerdatei oder leer, gelesen mit einem Fremd-Container
  docker run --rm -v "$1":/v:ro "$PROBE_IMG" sh -c 'cat /v/ubmess-marker 2>/dev/null' 2>&1
}

mark "Projekt sbx"
projekt sbx 1 || exit 1

mark "M2a: laufender Container"
hoch sbx ""
c=$(cid sbx); V=$(ws_volume "$c")
echo "container=${c:0:12} workspace-volume=$V"
docker exec "$c" sh -c 'echo ubmess-m2 > /workspaces/'"$(projektname sbx)"'/ubmess-marker' && echo "marker im Workspace des Containers geschrieben"
befehl docker rm "$c"
aussage "m2a-rm-laufend-ohne-f-exit=$RC"
echo "laeuft danach: $(docker inspect --format '{{.State.Running}}' "$c" 2>&1)"
aussage "m2a-container-nach-rm-ohne-f=$(docker inspect --format '{{.State.Status}}' "$c" 2>&1)"
befehl docker rm -f "$c"
aussage "m2a-rm-laufend-mit-f-exit=$RC"
aussage "m2a-container-nach-rm-f=$(docker inspect "$c" > /dev/null 2>&1 && echo vorhanden || echo entfernt)"
aussage "m2a-workspace-volume-nach-rm-f=$(existiert "$V")"
echo "marker im Volume nach rm -f (gelesen mit $PROBE_IMG): $(marker_lesen "$V")"

mark "M2d: u-boot nach docker rm -f (Volume stehen gelassen)"
runter sbx
aussage "m2d-down-ohne-container-exit=$RC"
hoch sbx ""
aussage "m2d-up-exit=$UP_RC"
c2=$(cid sbx)
aussage "m2d-neuer-container=$([ "$c2" != "$c" ] && echo ja || echo nein)"
aussage "m2d-workspace-volume-gleich=$([ "$(ws_volume "$c2")" = "$V" ] && echo ja || echo nein)"
aussage "m2d-marker-im-neuen-container=$(docker exec "$c2" cat /workspaces/"$(projektname sbx)"/ubmess-marker 2>&1)"

mark "M2b: Volume, das ein zweiter Container haelt"
befehl docker run -d --name "$PFX-zweiter" -v "$V":/zweiter "$PROBE_IMG" sleep 3600
echo "zweiter Container sieht den Marker: $(docker exec "$PFX-zweiter" cat /zweiter/ubmess-marker 2>&1)"
echo "-- Entfernen des Devcontainers (docker rm -f -v), das Volume haelt jetzt nur der zweite"
befehl docker rm -f -v "$c2"
befehl docker volume rm "$V"
aussage "m2b-volume-rm-mit-zweitem-halter-exit=$RC"
aussage "m2b-volume-nach-rm=$(existiert "$V")"
echo "marker nach dem fehlgeschlagenen volume rm: $(docker exec "$PFX-zweiter" cat /zweiter/ubmess-marker 2>&1)"
befehl docker volume rm -f "$V"
aussage "m2b-volume-rm-f-mit-zweitem-halter-exit=$RC"
aussage "m2b-volume-nach-rm-f=$(existiert "$V")"
befehl docker rm -f -v "$PFX-zweiter"
befehl docker volume rm "$V"
aussage "m2b-volume-rm-ohne-halter-exit=$RC"
aussage "m2b-volume-nach-letztem-rm=$(existiert "$V")"

mark "M2d (Fortsetzung): Neuanlage nach docker volume rm"
hoch sbx ""
c3=$(cid sbx)
aussage "m2d-up-nach-volume-rm-exit=$UP_RC"
aussage "m2d-workspace-volume-neu-angelegt=$([ "$(ws_volume "$c3")" = "$V" ] && echo gleicher-name || echo anderer-name)"
aussage "m2d-marker-nach-volume-rm=$(docker exec "$c3" sh -c 'cat /workspaces/'"$(projektname sbx)"'/ubmess-marker 2>&1 || true' | head -n 1)"

mark "M2c: gestoppter Container"
V3=$(ws_volume "$c3")
mount_volumes "$c3" > "$TMP/m2-vol-laufend"
runter sbx
aussage "m2c-down-exit=$RC"
aussage "m2c-container-status-nach-down=$(docker inspect --format '{{.State.Status}}' "$c3")"
echo "-- Mounts des gestoppten Containers:"; mounts "$c3"
aussage "m2c-mounts-des-gestoppten-containers-gleich-wie-laufend=$(mount_volumes "$c3" | diff -q "$TMP/m2-vol-laufend" - > /dev/null && echo ja || echo nein)"
befehl docker volume rm "$V3"
aussage "m2c-volume-rm-mit-gestopptem-halter-exit=$RC"
befehl docker volume rm -f "$V3"
aussage "m2c-volume-rm-f-mit-gestopptem-halter-exit=$RC"
aussage "m2c-volume-nach-rm-f=$(existiert "$V3")"
befehl docker rm "$c3"
aussage "m2c-rm-gestoppt-ohne-f-exit=$RC"
aussage "m2c-volume-nach-rm-gestoppt=$(existiert "$V3")"
befehl docker volume rm "$V3"
aussage "m2c-volume-rm-nach-container-rm-exit=$RC"
runter sbx
aussage "m2c-down-nach-rm-exit=$RC"

mark "M2f: u-boot devcontainer up --recreate mit benanntem Volume"
hoch sbx ""
c5=$(cid sbx); V5=$(ws_volume "$c5")
docker exec "$c5" sh -c 'echo ubmess-m2f > /workspaces/'"$(projektname sbx)"'/ubmess-marker'
hoch sbx "" --recreate
aussage "m2f-recreate-exit=$UP_RC"
c6=$(cid sbx)
aussage "m2f-neuer-container=$([ "$c6" != "$c5" ] && echo ja || echo nein) alter-container-noch-da=$(docker inspect "$c5" > /dev/null 2>&1 && echo ja || echo nein)"
aussage "m2f-workspace-volume-gleich=$([ "$(ws_volume "$c6")" = "$V5" ] && echo ja || echo nein) volume-nach-recreate=$(existiert "$V5")"
aussage "m2f-marker-im-neuen-container=$(docker exec "$c6" cat /workspaces/"$(projektname sbx)"/ubmess-marker 2>&1)"

mark "M2e: anonyme Volumes (Probe-Container, nur Docker)"
for variante in ohne-v mit-v; do
  n="$PFX-probe-$variante"
  docker run -d --name "$n" -v /ubmess-anon -v "$PFX-benannt-$variante":/ubmess-benannt "$PROBE_IMG" sleep 3600 > /dev/null
  anon=$(docker inspect --format '{{range .Mounts}}{{if eq .Destination "/ubmess-anon"}}{{.Name}}{{end}}{{end}}' "$n")
  merke_volumes "$n"
  echo "$variante: anonymes Volume ${anon:0:12}... (Namenslaenge ${#anon}), benanntes $PFX-benannt-$variante"
  if [ "$variante" = ohne-v ]; then docker rm -f "$n" > /dev/null; else docker rm -f -v "$n" > /dev/null; fi
  aussage "m2e-probe-$variante-anonym-danach=$(existiert "$anon") benannt-danach=$(existiert "$PFX-benannt-$variante")"
  [ "$variante" = ohne-v ] && echo "  docker volume ls -f dangling=true (Praefix- und Register-Volumes): $(docker volume ls -q -f dangling=true | grep -Fx "$anon")"
done

mark "M2e: Projekt anon (Standardprofil, VOLUME im Dockerfile), docker rm -f ohne -v, danach up --recreate"
projekt anon 0 || exit 1
echo 'VOLUME /ubmess-anon-daten' >> "$(dir_von anon)/.devcontainer/Dockerfile"
hoch anon ""
a=$(cid anon); mounts "$a"
A=$(docker inspect --format '{{range .Mounts}}{{if eq .Destination "/ubmess-anon-daten"}}{{.Name}}{{end}}{{end}}' "$a")
aussage "m2e-anon-projekt-anonymes-volume-in-mounts=$([ -n "$A" ] && echo ja || echo nein)"
hoch anon "" --recreate
a2=$(cid anon)
A2=$(docker inspect --format '{{range .Mounts}}{{if eq .Destination "/ubmess-anon-daten"}}{{.Name}}{{end}}{{end}}' "$a2")
aussage "m2e-up-recreate-neuer-container=$([ "$a2" != "$a" ] && echo ja || echo nein) altes-anonymes-volume-danach=$(existiert "$A") neues-anonymes-volume-gleich=$([ "$A2" = "$A" ] && echo ja || echo nein)"
befehl docker rm -f "$a2"
aussage "m2e-anon-projekt-rm-f-ohne-v: anonymes-volume-danach=$(existiert "$A2")"
hoch anon ""
a3=$(cid anon)
A3=$(docker inspect --format '{{range .Mounts}}{{if eq .Destination "/ubmess-anon-daten"}}{{.Name}}{{end}}{{end}}' "$a3")
befehl docker rm -f -v "$a3"
aussage "m2e-anon-projekt-rm-f-mit-v: anonymes-volume-danach=$(existiert "$A3")"
echo "Volumes in der Engine, die das Register als Rest kennt (vor dem Aufraeumen):"
{ [ -f "$REGISTER" ] && sort -u "$REGISTER"; } | while read -r v; do [ -n "$v" ] && echo "  $v: $(existiert "$v")"; done
