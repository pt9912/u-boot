#!/bin/bash
# M4: Welche Volumes ueberleben das Entfernen des Containers absichtlich (das Volume der Ollama-Modelle, laut Lastenheft), und wie sind
# sie zu unterscheiden?
# Projekt full (Sandbox + nestedRuntime docker + ollama): ein Container ohne Instanz mit Workspace-, Docker-Speicher- und Ollama-Volume.
#   - Marker in jedes Volume; im Ollama-Verzeichnis zusaetzlich das Schluesselpaar, das `ollama serve` beim ersten Start anlegt
#     (ein Modell wird nicht geladen: der Marker steht stellvertretend, das Laden braucht Netz und Platz)
#   - docker rm -f -v: welche Volumes bleiben?  Was sagt die Engine ueber ein verwaistes Volume (volume inspect, ls)?
#   - up (neuer Container): haengen dieselben Volumes ein, ist der Inhalt da?
#   - Auswahl nach dem Entfernen allein ueber den Namen: Infix `-ollama-` und gemeinsames Namensende der ${devcontainerId}
# Aufruf: UBOOT=<binary> WORK=<messordner> ./04-m4-ueberleben.sh      -> rohdaten-m4.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-m4.txt
echo "== M4 Volumes, die das Entfernen des Containers ueberstehen, $(now)"
anfang

existiert() { docker volume inspect "$1" > /dev/null 2>&1 && echo vorhanden || echo entfernt; }

mark "Projekt full, up"
projekt full 1 devcontainer.sandbox.nestedRuntime=docker devcontainer.features.ollama.enabled=true || exit 1
P=$(projektname full)
hoch full ""; aussage "m4-up-exit=$UP_RC"
c=$(cid full)
echo "-- Mounts vor dem Entfernen:"; mounts "$c"
mount_volumes "$c" > "$TMP/m4-vol"
docker inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}{{.Destination}} {{.Name}}{{"\n"}}{{end}}{{end}}' "$c" | sort > "$TMP/m4-ziel-name"
Vws=$(grep -- '/workspaces/' "$TMP/m4-ziel-name" | cut -d' ' -f2)
Vdk=$(grep -- '^/var/lib/docker ' "$TMP/m4-ziel-name" | cut -d' ' -f2)
Vol=$(grep -- '^/home/vscode/.ollama ' "$TMP/m4-ziel-name" | cut -d' ' -f2)
echo "workspace=$Vws"; echo "docker=$Vdk"; echo "ollama=$Vol"

mark "Inhalt anlegen"
docker exec "$c" sh -c "echo m4-ws > /workspaces/$P/ubmess-marker"
docker exec "$c" sh -c 'echo m4-docker | sudo tee /var/lib/docker/ubmess-marker > /dev/null'
docker exec "$c" sh -c 'echo m4-ollama > /home/vscode/.ollama/ubmess-marker'
docker exec "$c" sh -c 'ollama --version 2>&1 | head -n 2; (ollama serve > /tmp/ollama-serve.log 2>&1 &) ; sleep 6; ls -la /home/vscode/.ollama; tail -n 3 /tmp/ollama-serve.log | cut -c1-200'
KEY_VOR=$(docker exec "$c" sh -c 'sha256sum /home/vscode/.ollama/id_ed25519.pub 2>/dev/null | cut -d" " -f1')
aussage "m4-schluesselpaar-angelegt=$([ -n "$KEY_VOR" ] && echo ja || echo nein)"
docker exec "$c" sh -c 'pkill ollama; true' > /dev/null 2>&1

mark "docker rm -f -v des Containers"
befehl docker rm -f -v "$c"
aussage "m4-rm-f-v-exit=$RC"
for v in "$Vws" "$Vdk" "$Vol"; do echo "  $v: $(existiert "$v")"; done
aussage "m4-volumes-nach-rm-f-v: workspace=$(existiert "$Vws") docker=$(existiert "$Vdk") ollama=$(existiert "$Vol")"

mark "Was die Engine ueber die verwaisten Volumes sagt"
for v in "$Vws" "$Vdk" "$Vol"; do docker volume inspect --format '  {{.Name}} labels={{json .Labels}} options={{json .Options}} created={{.CreatedAt}}' "$v" | cut -c1-260; done
echo "-- docker volume ls --filter dangling=true (nur Volumes dieses Projekts):"
docker volume ls -q --filter dangling=true | grep -F "$P-" | sed 's/^/  /'
echo "-- docker system df -v, Volumes dieses Projekts (Platz):"
docker system df -v 2>/dev/null | grep -F "$P-" | grep -v '^vsc-' | awk '{print "  " $1 " links=" $(NF-1) " groesse=" $NF}' | cut -c1-200
echo "-- Auswahl nach dem Namen allein:"
echo "  Infix -ollama-: $(docker volume ls -q | grep -F -- "$P-ollama-" | wc -l) Volume(s)"
echo "  gemeinsames Namensende (52 Zeichen) der drei: $(printf '%s\n' "$Vws" "$Vdk" "$Vol" | sed -E 's/.*-([0-9a-z]{52})$/\1/' | sort -u | wc -l) verschiedene"
echo "  docker volume ls -q --filter name=$P- : $(docker volume ls -q --filter "name=$P-" | wc -l) Volume(s)"

mark "up nach dem Entfernen: neuer Container, dieselben Volumes?"
hoch full ""; aussage "m4-up-nach-rm-exit=$UP_RC"
c2=$(cid full)
aussage "m4-neuer-container=$([ "$c2" != "$c" ] && echo ja || echo nein)"
mount_volumes "$c2" > "$TMP/m4-vol2"
aussage "m4-volume-namen-gleich=$(diff -q "$TMP/m4-vol" "$TMP/m4-vol2" > /dev/null && echo ja || echo nein)"
aussage "m4-marker-workspace=$(docker exec "$c2" cat /workspaces/$P/ubmess-marker 2>&1)"
aussage "m4-marker-docker=$(docker exec "$c2" sudo cat /var/lib/docker/ubmess-marker 2>&1)"
aussage "m4-marker-ollama=$(docker exec "$c2" cat /home/vscode/.ollama/ubmess-marker 2>&1)"
KEY_NACH=$(docker exec "$c2" sh -c 'sha256sum /home/vscode/.ollama/id_ed25519.pub 2>/dev/null | cut -d" " -f1')
aussage "m4-oeffentlicher-schluessel-gleich=$([ -n "$KEY_VOR" ] && [ "$KEY_VOR" = "$KEY_NACH" ] && echo ja || echo nein)"

mark "Auswahl nach dem Entfernen allein ueber den Namen: alles ausser dem Ollama-Volume entfernen"
docker rm -f -v "$c2" > /dev/null
for v in $(docker volume ls -q --filter "name=$P-" | grep -v -- "$P-ollama-"); do docker volume rm "$v" > /dev/null 2>&1; echo "  volume rm ${v%-*} exit=$?"; done
aussage "m4-nach-auswahl: workspace=$(existiert "$Vws") docker=$(existiert "$Vdk") ollama=$(existiert "$Vol")"
mark "Neuanlage mit dem verbliebenen Ollama-Volume"
hoch full ""; aussage "m4-up-mit-nur-ollama-volume-exit=$UP_RC"
c3=$(cid full)
aussage "m4-marker-ollama-im-dritten-container=$(docker exec "$c3" cat /home/vscode/.ollama/ubmess-marker 2>&1)"
aussage "m4-marker-workspace-im-dritten-container=$(docker exec "$c3" cat /workspaces/$P/ubmess-marker 2>&1 | head -n 1)"
