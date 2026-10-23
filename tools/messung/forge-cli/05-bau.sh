#!/bin/bash
# Bau der erzeugten Projekte mit der Dev-Container-CLI (Docker-Socket nur fuer die Messung), kalt (--no-cache), und Inhalt der Images.
# Voraussetzung: 04-projekte.sh. Das Hilfs-Image mit der CLI baut Dockerfile.cli.
# Aufruf: WORK=<messordner> ./05-bau.sh [variante...]    (Vorgabe: alle, siehe unten)
#   -> $OUT/rohdaten-05.txt, rohdaten-05-<variante>.log
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-05.txt

CLI=$PFX-devcontainer-cli
variants=("$@"); [ ${#variants[@]} -gt 0 ] || variants=(bas basn gh ghx glabs glabp glabz ghn ghk ghxk glabsk glabpk glabzk ghnk)
# Varianten mit Endung k sind HAND-Varianten: Kopie des erzeugten Projekts, in der devcontainer.json steht der Feature-Schluessel ohne das
# Schema (`ghcr.io/…` statt `https://ghcr.io/…`). Das ist NICHT die Ausgabe von `u-boot generate devcontainer` (Befund in rohdaten-05.txt).
for v in "${variants[@]}"; do
  case "$v" in
    *k) src=${v%k}; rm -rf "$(vdir "$v")"; cp -a "$(vdir "$src")" "$(vdir "$v")"
        tool sed -i 's#"https://ghcr.io/#"ghcr.io/#' "/w/p-$v/.devcontainer/devcontainer.json"
        echo "Hand-Variante $v, Schluessel: $(grep -h "ghcr.io" "$(vdir "$v")/.devcontainer/devcontainer.json" | tr -d " ")" ;;
  esac
done
cli() { local v=$1; shift; docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "$(vdir "$v")":"/w/$(vname "$v")" -w "/w/$(vname "$v")" "$CLI" "$@"; }

echo "== Bau, $(now)"
docker build -q -f "$here/Dockerfile.cli" -t "$CLI" "$here" > "$WORK/cli-build.log" 2>&1 || { echo "Hilfs-Image baut nicht"; tail -5 "$WORK/cli-build.log"; exit 1; }
echo "devcontainer-cli $(cli bas devcontainer --version)"
aussage "cli-version=$(cli bas devcontainer --version)"

for v in "${variants[@]}"; do
  echo; mark "devcontainer build $v (kalt)"
  s=$(date +%s.%N)
  cli "$v" devcontainer build --workspace-folder "/w/$(vname "$v")" --image-name "$(vimg "$v")" --no-cache --log-level info > "$WORK/bau-$v.log" 2>&1; rc=$?
  e=$(date +%s.%N)
  redact < "$WORK/bau-$v.log" > "$OUT/rohdaten-05-$v.log"
  echo "build $v exit=$rc dauer=$(sekunden "$s" "$e")s"
  aussage "bau-$v-exit=$rc"
  if [ "$rc" != 0 ]; then echo "-- letzte Zeilen"; tail -12 "$WORK/bau-$v.log" | cut -c1-240; continue; fi
  docker image inspect "$(vimg "$v")" --format 'image={{.Id}} size={{.Size}} arch={{.Architecture}}'
  echo "-- Werkzeuge im Image (als Benutzer des Images)"
  docker run --rm "$(vimg "$v")" bash -c '
    id -un
    for t in gh glab; do
      p=$(command -v $t) && { echo "$t: $p"; $t --version 2>&1 | head -2; stat -c "%U:%G %a %s" "$p"; } || echo "$t: nicht im Image"
    done'
  docker run --rm "$(vimg "$v")" dpkg-query -W -f '${Package} ${Version}\n' | sort > "$WORK/dpkg-$v.txt"
  echo "pakete=$(grep -c . "$WORK/dpkg-$v.txt")"
done

echo
echo "== Paketunterschied gegen bas (mit Sandbox-Profil) bzw. basn (ohne); nur die Varianten, die gebaut haben"
for v in "${variants[@]}"; do
  case "$v" in bas|basn) continue ;; esac
  ref=bas; case "$v" in ghn|ghnk) ref=basn ;; esac
  [ -s "$WORK/dpkg-$v.txt" ] && [ -s "$WORK/dpkg-$ref.txt" ] || continue
  echo "-- $v gegen $ref"
  diff "$WORK/dpkg-$ref.txt" "$WORK/dpkg-$v.txt" | grep '^[<>]' | head -30
  aussage "paketdiff-$v=$(diff "$WORK/dpkg-$ref.txt" "$WORK/dpkg-$v.txt" | grep -c '^[<>]')"
done
echo "-- Plattenplatz"
docker system df | head -3
