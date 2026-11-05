#!/bin/bash
# M1: Welche Volumes haelt der Container (benannte und anonyme), trägt `docker inspect` sie vollstaendig in `Mounts`, und folgt daraus
# eine Erkennung ueber den Container statt ueber die Aufloesung von ${devcontainerId}?
# Varianten (je ein Projekt, `u-boot devcontainer up` ohne Instanz, danach Container und Volumes weg):
#   std   Standardprofil               sbx   Sandbox-Profil
#   dock  Sandbox + nestedRuntime docker   pod   Sandbox + nestedRuntime podman
#   oll   Sandbox + Feature ollama         full  Sandbox + nestedRuntime docker + ollama
#   anon  Standardprofil, im Dockerfile des Projekts von Hand `VOLUME /ubmess-anon-daten` (anonymes Volume des Images)
# Vollstaendigkeit: die Menge der Volumes, die der Container im Lauf von `up` in der Engine angelegt hat (Differenz von `docker volume ls`
# vor und nach `up`), gegen die Volumes in `Mounts`.
# Aufruf: UBOOT=<binary> WORK=<messordner> ./01-m1-volumes.sh      -> rohdaten-m1.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-m1.txt
echo "== M1 Volumes eines Devcontainers, $(now)"
anfang

variante() { # <variante> <sandbox:0|1> [pfad=wert]...
  local v=$1; shift
  mark "Variante $v"
  projekt "$v" "$@" || { aussage "m1-$v-projekt=fehlt"; return; }
  if [ "$v" = anon ]; then
    echo 'VOLUME /ubmess-anon-daten' >> "$(dir_von "$v")/.devcontainer/Dockerfile"
    echo "  Dockerfile: Zeile VOLUME /ubmess-anon-daten angehaengt"
  fi
  vol_alle > "$TMP/m1-vol-vorher"
  hoch "$v" ""
  aussage "m1-$v-up-exit=$UP_RC"
  local c; c=$(cid "$v")
  if [ -z "$c" ]; then aussage "m1-$v-container=fehlt"; return; fi
  echo "container: ${c:0:12} name=$(docker inspect --format '{{.Name}}' "$c") image=$(docker inspect --format '{{.Config.Image}}' "$c" | cut -c1-40)..."
  echo "-- docker inspect, Mounts:"; mounts "$c"
  echo "-- docker inspect, Config.Volumes (Volumes, die das Image deklariert): $(docker inspect --format '{{json .Config.Volumes}}' "$c")"
  mount_volumes "$c" > "$TMP/m1-mountvol"
  vol_alle > "$TMP/m1-vol-nachher"
  comm -13 "$TMP/m1-vol-vorher" "$TMP/m1-vol-nachher" > "$TMP/m1-delta"
  echo "-- Volumes, die up in der Engine angelegt hat (docker volume ls vorher/nachher):"; sed 's/^/  /' "$TMP/m1-delta"
  local n_mv n_delta n_anon n_bind
  n_mv=$(wc -l < "$TMP/m1-mountvol"); n_delta=$(wc -l < "$TMP/m1-delta")
  n_anon=$(grep -cE '^[0-9a-f]{64}$' "$TMP/m1-mountvol")
  n_bind=$(docker inspect --format '{{range .Mounts}}{{if eq .Type "bind"}}b{{end}}{{end}}' "$c" | wc -c)
  aussage "m1-$v-volumes-in-mounts=$n_mv benannt=$((n_mv - n_anon)) anonym=$n_anon"
  aussage "m1-$v-bind-mounts=$((n_bind - 1))"
  if diff -q "$TMP/m1-mountvol" "$TMP/m1-delta" > /dev/null; then aussage "m1-$v-mounts-vollstaendig=ja (Mounts == Differenz der Engine, $n_delta Volumes)"
  else aussage "m1-$v-mounts-vollstaendig=NEIN"; diff "$TMP/m1-mountvol" "$TMP/m1-delta" | sed 's/^/  diff: /'; fi
  echo "-- docker volume inspect (Treiber, Labels, Optionen):"
  while read -r vol; do
    [ -n "$vol" ] && docker volume inspect --format '  {{.Name}} driver={{.Driver}} labels={{json .Labels}} options={{json .Options}}' "$vol"
  done < "$TMP/m1-mountvol"
  local nlab; nlab=$(while read -r vol; do [ -n "$vol" ] && docker volume inspect --format '{{len .Labels}}' "$vol"; done < "$TMP/m1-mountvol" | awk '{s+=$1} END{print s+0}')
  aussage "m1-$v-labels-an-volumes-insgesamt=$nlab"
  # das Namensende der benannten Volumes (die aufgeloeste ${devcontainerId}) in Labels und Umgebung des Containers?
  local suffix; suffix=$(grep -vE '^[0-9a-f]{64}$' "$TMP/m1-mountvol" | head -n 1 | sed -E 's/.*-([0-9a-z]{52})$/\1/')
  if [ -n "$suffix" ] && [ "${#suffix}" -eq 52 ]; then
    echo "-- Namensende des ersten benannten Volumes (${#suffix} Zeichen): $suffix"
    aussage "m1-$v-id-in-labels=$(docker inspect --format '{{json .Config.Labels}}' "$c" | grep -c "$suffix" | sed 's/^0$/nein/;s/^[1-9].*$/ja/') id-in-umgebung=$(docker inspect --format '{{json .Config.Env}}' "$c" | grep -c "$suffix" | sed 's/^0$/nein/;s/^[1-9].*$/ja/')"
    local k; k=$(grep -vE '^[0-9a-f]{64}$' "$TMP/m1-mountvol" | sed -E 's/.*-([0-9a-z]{52})$/\1/' | sort -u | wc -l)
    aussage "m1-$v-verschiedene-ids-ueber-alle-benannten-volumes=$k"
  fi
  echo "-- Labels des Containers (Namen): $(docker inspect --format '{{range $k,$v := .Config.Labels}}{{$k}} {{end}}' "$c")"
  # Laufzeit der verschachtelten Runtime, soweit die Variante sie anfordert
  case $v in
    dock|full|pod)
      echo "-- Statusdateien der Startscripts:"
      docker exec "$c" sh -c 'for f in /var/lib/u-boot/status/*; do [ -f "$f" ] && echo "  $f: $(cat "$f")"; done; true' 2>&1 | cut -c1-300
      if [ "$v" = pod ]; then echo "-- podman: $(docker exec "$c" sh -c 'podman --version 2>&1; podman info >/dev/null 2>&1; echo "podman info exit=$?"' 2>&1 | tr '\n' ' ')"
      else echo "-- docker im Container: $(docker exec "$c" sh -c 'timeout 20 docker info --format "server={{.ServerVersion}} driver={{.Driver}}" 2>&1 | head -n 2' 2>&1 | tr '\n' ' ')"; fi ;;
  esac
  # Entfernen des Containers mit -v: Wirkung auf die Volumes (Vorgriff auf M2)
  echo "-- docker rm -f -v: Volumes danach (Name: vorhanden?)"
  docker rm -f -v "$c" > /dev/null 2>&1; echo "  rm -f -v exit=$?"
  while read -r vol; do [ -n "$vol" ] && echo "  $vol: $(docker volume inspect "$vol" > /dev/null 2>&1 && echo vorhanden || echo entfernt)"; done < "$TMP/m1-mountvol"
  objekte_entfernen > /dev/null
}

variante std 0
variante sbx 1
variante dock 1 devcontainer.sandbox.nestedRuntime=docker
variante pod 1 devcontainer.sandbox.nestedRuntime=podman
variante oll 1 devcontainer.features.ollama.enabled=true
variante full 1 devcontainer.sandbox.nestedRuntime=docker devcontainer.features.ollama.enabled=true
variante anon 0
echo "base-image Config.Volumes: $(docker image inspect --format '{{json .Config.Volumes}}' mcr.microsoft.com/devcontainers/base:debian)"
