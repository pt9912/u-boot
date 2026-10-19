#!/bin/bash
# A10: Aufraeumen. Entfernt nur Objekte, die die Messung angelegt hat: Namen mit dem Praefix u-boot-messung-amd64 (Container, Images,
# Volumes, Netze), die Images der Dev-Container-CLI vsc-u-boot-messung-amd64-*, unbenannte Images, die erst waehrend der Messung
# entstanden sind (nicht in der Liste von 01-messstand.sh), und Images, die der Bau der Hilfs-Images gezogen hat (EXTRA_IMAGES,
# nur wenn sie VOR der Messung fehlten). Build-Cache und fremde Objekte bleiben. Liste vorher: rohdaten-11-vorher.txt (01-messstand.sh).
# Aufruf: WORK=<messordner> tools/messung/ollama-amd64/11-aufraeumen.sh      -> $OUT/rohdaten-11.txt, rohdaten-11-nachher.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-11.txt

EXTRA_IMAGES="docker:29-cli"
vorher="$OUT/rohdaten-11-vorher.txt"
[ -s "$vorher" ] || { echo "Abbruch: die Liste vorher ($vorher) fehlt oder ist leer; ohne sie wird nichts entfernt (zuerst 01-messstand.sh)"; exit 2; }
echo "== A10 Aufraeumen, $(now)"
echo "-- Objekte der Messung vor dem Aufraeumen"
docker ps -a --format '{{.Names}}' | grep -E "^$PFX" || echo "  keine Container"
docker images --format '{{.Repository}}:{{.Tag}}' | grep -E "^(vsc-)?$PFX" || echo "  keine Images mit Praefix"
docker volume ls -q | grep -E "^$PFX" || echo "  keine Volumes"
docker network ls --format '{{.Name}}' | grep -E "^$PFX" || echo "  keine Netze"

docker ps -a --format '{{.Names}}' | grep -E "^$PFX" | while read -r c; do docker rm -f -v "$c" > /dev/null; done
docker ps -a -q --filter "label=devcontainer.local_folder=/w/$PFX-oll" | while read -r c; do docker rm -f -v "$c" > /dev/null; done
# Zwischencontainer eines fehlgeschlagenen Baus (A8, Negativprobe): kein eigener Name, erkennbar am RUN-Befehl des Ollama-Blocks
docker ps -a -q --filter status=exited | while read -r c; do
  if docker inspect "$c" --format '{{.Path}} {{.Args}}' | grep -q 'OLLAMA_SHA256_AMD64'; then echo "entferne Zwischencontainer $c des fehlgeschlagenen Baus"; docker rm -v "$c" > /dev/null; fi
done
docker volume ls -q | grep -E "^$PFX" | while read -r v; do docker volume rm "$v" > /dev/null; done
docker network ls --format '{{.Name}}' | grep -E "^$PFX" | while read -r n; do docker network rm "$n" > /dev/null; done
docker images --format '{{.Repository}}:{{.Tag}}' | grep -E "^(vsc-)?$PFX" | while read -r i; do docker rmi "$i" > /dev/null; done
for i in $EXTRA_IMAGES; do
  if ! grep -q -F "$i" "$vorher" && docker image inspect "$i" > /dev/null 2>&1; then echo "entferne durch die Messung gezogenes Image $i"; docker rmi "$i" > /dev/null; fi
done
# unbenannte Images (<none>), die in der Liste vorher nicht standen: Reste der Bauten der Messung
docker images --format '{{.Repository}}:{{.Tag}}\t{{.ID}}' | grep -E '^<none>:<none>' | while IFS=$'\t' read -r _ id; do
  if ! grep -q -F "$id" "$vorher"; then echo "entferne unbenanntes Image $id (nicht in der Liste vorher)"; docker rmi "$id" > /dev/null 2>&1; fi
done

{
  echo "## docker ps -a"; docker ps -a --format '{{.Names}}\t{{.Image}}' | sort
  echo "## docker images"; docker images --format '{{.Repository}}:{{.Tag}}\t{{.ID}}' | sort
  echo "## docker volume ls"; docker volume ls -q | sort
  echo "## docker network ls"; docker network ls --format '{{.Name}}' | sort
} > "$OUT/rohdaten-11-nachher.txt"
echo
echo "-- Objekte mit dem Praefix nachher"
docker ps -a --format '{{.Names}}' | grep -E "^$PFX" || echo "  keine Container"
docker images --format '{{.Repository}}:{{.Tag}}' | grep -E "^(vsc-)?$PFX" || echo "  keine Images"
docker volume ls -q | grep -E "^$PFX" || echo "  keine Volumes"
docker network ls --format '{{.Name}}' | grep -E "^$PFX" || echo "  keine Netze"
aussage "objekte-mit-praefix-nachher=$( (docker ps -a --format '{{.Names}}'; docker images --format '{{.Repository}}'; docker volume ls -q; docker network ls --format '{{.Name}}') | grep -c -E "^(vsc-)?$PFX")"
echo
echo "-- Unterschied der Listen vorher -> nachher ('<' nur vorher, '>' nur nachher; Fremdes, das waehrend der Messung auf dem Host kam oder ging, steht hier ebenfalls)"
diff "$vorher" "$OUT/rohdaten-11-nachher.txt" && echo "(kein Unterschied)"
echo
echo "-- Build-Cache und Platz (der Build-Cache bleibt; er wurde nicht geleert)"
docker system df
docker run --rm -v /var/lib/docker:/d:ro "$BASE" df -h /d | tail -1
