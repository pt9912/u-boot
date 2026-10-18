#!/bin/bash
# M2: erneuter Lauf des Startscripts im laufenden Container (mit anderer Liste) und die atomare Variante.
# Aufruf: 03-m2-neuladen.sh > rohdaten-03-m2.txt 2>&1
# Voraussetzung: Docker (Colima); Image u-boot-messung-egress-neuladen-dig (README.md, Schritt 0).
# Jeder Abschnitt laeuft in einem frischen Container (sleep infinity, --cap-add=NET_ADMIN, postStartCommand als vscode,
# danach die Messung als root: der postStartCommand ruft `sudo sh` auf).
set -u
here=$(cd "$(dirname "$0")" && pwd)
IMG=u-boot-messung-egress-neuladen-dig
title() { case "$1" in
  A) echo "A: erneuter Lauf des gebackenen Scripts mit anderer Liste (Fenster ohne Filter, DNS)";;
  E) echo "E: atomare Variante (reload-atomar.sh): ein nft-Wechsel, dnsmasq erst nach dem Ende des alten";;
  B) echo "B: 50 erneute Laeufe des gebackenen Scripts hintereinander (Startrennen)";;
  F) echo "F: 50 Laeufe der atomaren Variante hintereinander";;
  G) echo "G: Zeit vom SIGTERM bis zum Ende des alten dnsmasq (20 Messungen, sortiert)";;
  C) echo "C: laufende Verbindung ueber einen Neu-Lauf, dessen Liste den Host nicht mehr enthaelt";;
  D) echo "D: Set nach dem Neu-Lauf; Verbindung zu einer zuvor aufgeloesten Adresse ohne neue Aufloesung";;
esac; }
for s in ${*:-A E B F G C D}; do
  echo "######## $(title "$s")"
  c=egress-m2
  docker rm -f $c >/dev/null 2>&1
  docker run -d --name $c --cap-add=NET_ADMIN -v "$here":/m:ro $IMG sleep infinity >/dev/null
  docker exec -u vscode $c sh -c 'sudo sh /usr/local/bin/u-boot-egress-init'; echo "postStart exit=$?"
  docker exec -u root $c bash /m/inner-m2.sh "$s"
  docker rm -f $c >/dev/null 2>&1
  echo
done
