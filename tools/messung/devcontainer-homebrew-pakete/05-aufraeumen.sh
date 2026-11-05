#!/bin/bash
# Aufraeumen: entfernt nur Container und Images, deren NAME das Praefix u-boot-messung-brew traegt (praefix_entfernen in lib.sh). Nicht entfernt
# werden das Basis-Image, der Build-Cache und fremde Objekte. Die Projekte unter WORK bleiben (der Messordner liegt ausserhalb des Repos).
# Aufruf: UBOOT=<binary> WORK=<messordner> ./05-aufraeumen.sh
set -u
. "$(dirname "$0")/lib.sh"
praefix_entfernen
echo "verbleibende Objekte mit dem Praefix: $(docker images --format '{{.Repository}}' | grep -c "^$PFX")"
