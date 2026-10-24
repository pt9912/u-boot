#!/bin/bash
# Aufraeumen. Entfernt nur Objekte, deren NAME das Praefix u-boot-messung-ocloud-konto oder vsc-u-boot-messung-ocloud-konto traegt (Container, Volumes,
# Netze, Images; praefix_entfernen in lib.sh), und die Container der Dev-Container-CLI ueber das Label devcontainer.local_folder=/w/<Praefix>-<variante>.
# Nicht entfernt werden Tags ohne Praefix (Basis-Images der Bauten, alpine/socat), unbenannte Images, der Build-Cache und alle fremden Objekte.
#   weg-a   Container, Volumes, Netze; die Images bleiben (Weg b braucht das Image von res)
#   alles   zusaetzlich die Images mit dem Praefix
# Voraussetzung: 00-vorlauf.sh, rohdaten-01-vorher.txt (01-messstand.sh).
# Aufruf: WORK=<messordner> tools/messung/ollama-cloud-konto/08-aufraeumen.sh weg-a|alles  -> $ROH/rohdaten-08-<stufe>.txt, rohdaten-08-<stufe>-nachher.txt
set -u
. "$(dirname "$0")/lib.sh"
vorlauf_pruefen || exit 2
stufe=${1:-}
case "$stufe" in weg-a) KEEP_IMAGES=1 ;; alles) KEEP_IMAGES=0 ;; *) echo "Aufruf: 08-aufraeumen.sh weg-a|alles"; exit 2 ;; esac
rohdaten "rohdaten-08-$stufe.txt"
vorher="$ROH/rohdaten-01-vorher.txt"
[ -s "$vorher" ] || { echo "Abbruch: die Liste vorher fehlt oder ist leer (zuerst 01-messstand.sh)"; exit 2; }
echo "== Aufraeumen ($stufe), $(now)"
echo "-- vorher (rohdaten-01-vorher.txt)"; cat "$vorher"
echo "-- vor dem Aufraeumen"; bestand
echo "-- Entfernen"
for v in std res alw; do
  docker ps -a -q --filter "label=devcontainer.local_folder=$(wdir "$v")" | while read -r c; do echo "entferne container der Dev-Container-CLI $c ($v)"; docker rm -f -v "$c" > /dev/null; done
done
praefix_entfernen
echo "-- nachher (rohdaten-08-$stufe-nachher.txt)"
bestand | tee "$ROH/rohdaten-08-$stufe-nachher.txt"
aussage "objekte-mit-praefix-nachher-$stufe=$(bestand | grep -c -E '^(container|image|volume|netz) ')"
aussage "relay-container-nachher=$(docker ps -a --format '{{.Names}}' | grep -c -x "$RELAY")"
echo "vorher:  $(grep '^container=' "$vorher")"
echo "nachher: $(grep '^container=' "$ROH/rohdaten-08-$stufe-nachher.txt")"
echo "(Aenderungen fremder Arbeiten auf dem Host waehrend der Messung stehen in diesem Vergleich ebenfalls)"
prompt_stand
