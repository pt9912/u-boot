#!/bin/bash
# Aufraeumen: entfernt Container, Volumes und Images, deren NAME das Praefix ubmess-rm traegt (objekte_entfernen in lib.sh), und die
# Projekte unter WORK. Nicht entfernt werden das Basis-Image, das Hilfs-Image von `u-boot devcontainer up`, der Build-Cache und fremde Objekte.
# Vorher werden Groesse und Zahl der Images mit dem Praefix und der freie Platz der Docker-Wurzel gelesen (der Platzbedarf der Messung).
# Aufruf: UBOOT=<binary> WORK=<messordner> ./05-aufraeumen.sh      -> rohdaten-aufraeumen.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-aufraeumen.txt
echo "== Aufraeumen, $(now)"
frei() { df -B1 --output=avail /var/lib/docker 2> /dev/null | tail -n 1 | tr -d ' '; }
F0=$(frei)
echo "vorher, Objekte mit dem Praefix:"; praefix_objekte | sed 's/^/  /'
echo "vorher, Images mit dem Praefix (Repository Tag Groesse; die Groesse zaehlt gemeinsame Schichten je Image):"
docker images --format '{{.Repository}} {{.Tag}} {{.Size}}' | grep -E "^vsc-$PFX" | sed 's/^/  /'
echo "vorher, freier Platz der Docker-Wurzel: $F0 Byte"
objekte_entfernen mit-images
F1=$(frei)
echo "nachher, freier Platz der Docker-Wurzel: $F1 Byte (Differenz $((F1 - F0)) Byte; die Zahl streut mit anderer Last auf dem Host)"
rm -rf "$WORK/$PFX"-*   # die Projekte und TMP (alles mit dem Praefix im Namen)
echo "nachher, Objekte mit dem Praefix: $(praefix_objekte | wc -l)"
praefix_objekte | sed 's/^/  /'
echo "Projekte unter WORK: $(ls -d "$WORK/$PFX"-* 2> /dev/null | wc -l)"
echo "== ende $(now)"
