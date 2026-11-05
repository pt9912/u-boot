#!/bin/bash
# M5: Welchen Container meint `u-boot devcontainer down` ohne --instance, und was geschieht, wenn es im Projekt keinen Container ohne
# Instanz gibt (nur Container von Instanzen) oder gar keinen? Derselbe Suchweg (FindContainer) steht hinter `up` und `shell`.
#   A  kein Container im Projekt: `down` und `down --instance a`
#   B  Container ohne Instanz plus Instanzen a und b: `down` ohne --instance (welcher Container wird gestoppt?)
#   C  nur Instanz-Container (der Container ohne Instanz ist mit `docker rm -f -v` entfernt): `down` und `up` ohne --instance
#   D  Instanz-Container, einer gestoppt: `down` ohne --instance
#   E  alle Instanz-Container gestoppt: `down` ohne --instance
# Projekt sbx (Sandbox-Profil). `shell` ist interaktiv und nicht Teil der Messung.
# Aufruf: UBOOT=<binary> WORK=<messordner> ./04b-m5-down-ohne-instanz.sh      -> rohdaten-m5.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-m5.txt
echo "== M5 down ohne --instance, $(now)"
anfang

laeuft() { docker inspect --format '{{.State.Running}}' "$1" 2>&1; }
# zustand: die Container des Projektordners mit Zustand und Wert des Labels u-boot.instance
zustand() {
  echo "-- Container des Projektordners (Name, Zustand, Label u-boot.instance):"
  docker ps -a --no-trunc --filter "label=devcontainer.local_folder=$(dir_von sbx)" --format '  {{.Names}} state={{.State}} instance={{.Label "u-boot.instance"}}' | sort
}

mark "A: Projekt sbx ohne Container"
projekt sbx 1 || exit 1
zustand
runter sbx; aussage "m5-a-down-ohne-container-exit=$RC"
runter sbx a; aussage "m5-a-down-instanz-a-ohne-container-exit=$RC"

mark "B: Container ohne Instanz, Instanz a, Instanz b"
hoch sbx "";  aussage "m5-b-up-ohne-instanz-exit=$UP_RC"
hoch sbx a;   aussage "m5-b-up-a-exit=$UP_RC"
hoch sbx b;   aussage "m5-b-up-b-exit=$UP_RC"
c0=$(cid sbx); ca=$(cid sbx a); cb=$(cid sbx b)
echo "container: ohne=${c0:0:12} a=${ca:0:12} b=${cb:0:12}"
zustand
runter sbx; aussage "m5-b-down-ohne-instanz-exit=$RC"
aussage "m5-b-nach-down-ohne-instanz: ohne-instanz-laeuft=$(laeuft "$c0") a-laeuft=$(laeuft "$ca") b-laeuft=$(laeuft "$cb")"
zustand

mark "C: nur Instanz-Container (Container ohne Instanz mit docker rm -f -v entfernt)"
befehl docker rm -f -v "$c0"
zustand
runter sbx; aussage "m5-c-down-ohne-instanz-nur-instanzen-exit=$RC"
aussage "m5-c-nach-down: a-laeuft=$(laeuft "$ca") b-laeuft=$(laeuft "$cb")"
hoch sbx ""; aussage "m5-c-up-ohne-instanz-nur-instanzen-exit=$UP_RC"
aussage "m5-c-container-ohne-instanz-angelegt=$([ -n "$(cid sbx)" ] && echo ja || echo nein)"
zustand

mark "D: Instanz a gestoppt (down --instance a), b laeuft"
runter sbx a; aussage "m5-d-down-a-exit=$RC"
aussage "m5-d-a-laeuft=$(laeuft "$ca") b-laeuft=$(laeuft "$cb")"
runter sbx; aussage "m5-d-down-ohne-instanz-exit=$RC"
aussage "m5-d-nach-down: a-laeuft=$(laeuft "$ca") b-laeuft=$(laeuft "$cb")"

mark "E: alle Instanz-Container gestoppt"
runter sbx b; aussage "m5-e-down-b-exit=$RC"
aussage "m5-e-a-laeuft=$(laeuft "$ca") b-laeuft=$(laeuft "$cb")"
zustand
runter sbx; aussage "m5-e-down-ohne-instanz-nur-gestoppte-instanzen-exit=$RC"
