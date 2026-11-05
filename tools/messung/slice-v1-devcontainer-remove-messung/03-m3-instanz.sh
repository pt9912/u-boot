#!/bin/bash
# M3: Traegt die Instanz (--instance) eigene Volumes (Namen, Isolation), und loescht das Entfernen einer Instanz nur deren Volumes?
#   Projekt sbx   Sandbox-Profil: Container ohne Instanz, Instanz a, Instanz b
#   Projekt full  Sandbox + nestedRuntime docker + ollama: Instanz a, Instanz b (alle vier Volume-Familien ausser Podman)
#   Projekt std   Standardprofil: Container ohne Instanz und Instanz x (Bind-Mount des Projektordners)
# Entfernen einer Instanz: die Volumes werden ueber `Mounts` DES CONTAINERS bestimmt (vor docker rm), nicht ueber den Namen.
# Aufruf: UBOOT=<binary> WORK=<messordner> ./03-m3-instanz.sh      -> rohdaten-m3.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-m3.txt
echo "== M3 Instanzen, $(now)"
anfang

existiert() { docker volume inspect "$1" > /dev/null 2>&1 && echo vorhanden || echo entfernt; }
ws_volume() { mount_volumes "$1" | grep -- '-workspace-' | head -n 1; }
laeuft() { docker inspect --format '{{.State.Running}}' "$1" 2>&1; }

mark "Projekt sbx: drei Container"
projekt sbx 1 || exit 1
P=$(projektname sbx)
hoch sbx ""; aussage "m3-sbx-up-ohne-instanz-exit=$UP_RC"
hoch sbx a; aussage "m3-sbx-up-a-exit=$UP_RC"
hoch sbx b; aussage "m3-sbx-up-b-exit=$UP_RC"
c0=$(cid sbx); ca=$(cid sbx a); cb=$(cid sbx b)
echo "container: ohne=${c0:0:12} a=${ca:0:12} b=${cb:0:12}"
for x in 0 a b; do eval "c=\$c$x"; echo "-- Mounts ($x):"; mounts "$c"; done
V0=$(ws_volume "$c0"); Va=$(ws_volume "$ca"); Vb=$(ws_volume "$cb")
aussage "m3-sbx-workspace-volumes-verschieden=$([ "$V0" != "$Va" ] && [ "$V0" != "$Vb" ] && [ "$Va" != "$Vb" ] && echo ja || echo nein)"
echo "-- docker volume ls --filter name=$P-workspace- (Namen allein; welcher Name zu welchem Container gehoert, steht nicht darin):"
docker volume ls -q --filter "name=$P-workspace-" | sort | sed 's/^/  /'
echo "-- Namenslaenge: ohne=${#V0} a=${#Va} b=${#Vb}"
echo "-- Labels der drei Container (devcontainer.local_folder / u-boot.instance):"
docker ps -a --no-trunc --filter "label=devcontainer.local_folder=$(dir_von sbx)" --format '  {{.Names}} instance={{.Label "u-boot.instance"}}'
echo "-- Filter auf das Instanz-Label a: $(docker ps -aq --filter "label=u-boot.instance=a" --filter "label=devcontainer.local_folder=$(dir_von sbx)" | wc -l) Container"
echo "-- Marker je Instanz in das Workspace-Volume schreiben und gegenpruefen (Isolation):"
for x in 0 a b; do eval "c=\$c$x"; docker exec "$c" sh -c "echo marker-$x > /workspaces/$P/ubmess-marker"; done
for x in 0 a b; do eval "c=\$c$x"; echo "  container $x liest: $(docker exec "$c" cat /workspaces/$P/ubmess-marker)"; done
aussage "m3-sbx-marker-isoliert=$( [ "$(docker exec "$c0" cat /workspaces/$P/ubmess-marker)" = marker-0 ] && [ "$(docker exec "$ca" cat /workspaces/$P/ubmess-marker)" = marker-a ] && [ "$(docker exec "$cb" cat /workspaces/$P/ubmess-marker)" = marker-b ] && echo ja || echo nein)"

mark "Entfernen der Instanz a (Container und die Volumes AUS DEREN Mounts)"
volsa=$(mount_volumes "$ca")
echo "Volumes aus Mounts von a: $volsa"
befehl docker rm -f "$ca"
for v in $volsa; do docker volume rm "$v" > /dev/null 2>&1; echo "  volume rm $v exit=$?"; done
aussage "m3-sbx-a-container-nach-entfernen=$(docker inspect "$ca" > /dev/null 2>&1 && echo vorhanden || echo entfernt) a-volume=$(existiert "$Va")"
aussage "m3-sbx-ohne-instanz-unberuehrt: container-laeuft=$(laeuft "$c0") volume=$(existiert "$V0") marker=$(docker exec "$c0" cat /workspaces/$P/ubmess-marker)"
aussage "m3-sbx-b-unberuehrt: container-laeuft=$(laeuft "$cb") volume=$(existiert "$Vb") marker=$(docker exec "$cb" cat /workspaces/$P/ubmess-marker)"
echo "-- u-boot sieht nach dem Entfernen von a:"
runter sbx a; aussage "m3-sbx-down-a-nach-entfernen-exit=$RC"
runter sbx b; aussage "m3-sbx-down-b-exit=$RC (stoppt nur b)"
aussage "m3-sbx-nach-down-b: ohne-instanz-laeuft=$(laeuft "$c0") b-laeuft=$(laeuft "$cb")"
echo "-- Entfernen von b und des Containers ohne Instanz bleibt je einzeln: erst b"
vols_b=$(mount_volumes "$cb"); docker rm -f "$cb" > /dev/null; for v in $vols_b; do docker volume rm "$v" > /dev/null 2>&1; done
aussage "m3-sbx-nach-entfernen-von-b: ohne-instanz-laeuft=$(laeuft "$c0") volume-ohne-instanz=$(existiert "$V0") volume-b=$(existiert "$Vb")"
docker rm -f -v "$c0" > /dev/null; docker volume rm "$V0" > /dev/null 2>&1
objekte_entfernen > /dev/null

mark "Projekt full (nestedRuntime docker + ollama): Instanzen a und b"
projekt full 1 devcontainer.sandbox.nestedRuntime=docker devcontainer.features.ollama.enabled=true || exit 1
P=$(projektname full)
hoch full a; aussage "m3-full-up-a-exit=$UP_RC"
hoch full b; aussage "m3-full-up-b-exit=$UP_RC"
ca=$(cid full a); cb=$(cid full b)
echo "-- Mounts (a):"; mounts "$ca"; echo "-- Mounts (b):"; mounts "$cb"
mount_volumes "$ca" > "$TMP/m3-vol-a"; mount_volumes "$cb" > "$TMP/m3-vol-b"
aussage "m3-full-volumes-a=$(wc -l < "$TMP/m3-vol-a") b=$(wc -l < "$TMP/m3-vol-b") gemeinsame=$(comm -12 "$TMP/m3-vol-a" "$TMP/m3-vol-b" | wc -l)"
echo "-- Familien je Volume (Ziel im Container -> Namensinfix):"
for x in a b; do eval "c=\$c$x"; docker inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}{{.Destination}} -> {{.Name}}{{"\n"}}{{end}}{{end}}' "$c" | grep -v "^$" | sed "s/^/  [$x] /"; done
echo "-- Entfernen von a mit den Volumes aus Mounts von a:"
docker rm -f "$ca" > /dev/null; echo "  rm -f exit=$?"
while read -r v; do [ -n "$v" ] && { docker volume rm "$v" > /dev/null 2>&1; echo "  volume rm $v exit=$?"; }; done < "$TMP/m3-vol-a"
aussage "m3-full-a-volumes-nach-entfernen=$(while read -r v; do existiert "$v"; done < "$TMP/m3-vol-a" | sort | uniq -c | tr '\n' ' ')"
aussage "m3-full-b-volumes-nach-entfernen-von-a=$(while read -r v; do existiert "$v"; done < "$TMP/m3-vol-b" | sort | uniq -c | tr '\n' ' ') b-laeuft=$(laeuft "$cb")"
docker rm -f -v "$cb" > /dev/null; while read -r v; do docker volume rm "$v" > /dev/null 2>&1; done < "$TMP/m3-vol-b"
objekte_entfernen > /dev/null

mark "Projekt std (Standardprofil): Container ohne Instanz und Instanz x"
projekt std 0 || exit 1
P=$(projektname std)
hoch std ""; aussage "m3-std-up-ohne-instanz-exit=$UP_RC"
hoch std x;  aussage "m3-std-up-x-exit=$UP_RC"
c0=$(cid std); cx=$(cid std x)
echo "-- Mounts (ohne Instanz):"; mounts "$c0"; echo "-- Mounts (x):"; mounts "$cx"
aussage "m3-std-volumes-ohne=$(mount_volumes "$c0" | wc -l) x=$(mount_volumes "$cx" | wc -l)"
aussage "m3-std-bind-quelle-gleich=$([ "$(docker inspect --format '{{range .Mounts}}{{if eq .Type "bind"}}{{.Source}}{{end}}{{end}}' "$c0")" = "$(docker inspect --format '{{range .Mounts}}{{if eq .Type "bind"}}{{.Source}}{{end}}{{end}}' "$cx")" ] && echo ja || echo nein)"
docker exec "$cx" sh -c "echo std-marker > /workspaces/$P/ubmess-marker"
echo "marker im Projektordner des Hosts: $(cat "$(dir_von std)/ubmess-marker" 2>&1)"
echo "marker im Container ohne Instanz: $(docker exec "$c0" cat /workspaces/$P/ubmess-marker 2>&1)"
befehl docker rm -f -v "$cx"
aussage "m3-std-nach-rm-von-x: marker-im-projektordner=$(cat "$(dir_von std)/ubmess-marker" 2>&1) ohne-instanz-laeuft=$(laeuft "$c0")"
