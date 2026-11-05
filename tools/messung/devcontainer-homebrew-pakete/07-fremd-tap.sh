#!/bin/bash
# M6 Abhaengigkeit einer Tap-Formel aus einem anderen Tap: richtet Homebrew einen Tap, auf den die Abhaengigkeit einer Formel zeigt, im Bau
# selbst ein, oder scheitert der Bau? Die Formel xdebug@8.2 des Taps shivammathur/extensions hat die Abhaengigkeit shivammathur/php/php@8.2
# (anderer Tap). Gemessen: (a) der Bau mit dem erzeugten Block, nur der Tap der Formel ist freigegeben; (b) im Image der Variante alone
# `brew install --dry-run` mit beiden Formeln in beiden Reihenfolgen (welche Taps danach eingerichtet sind). Es wird nichts installiert.
# Voraussetzung: 01-bau.sh (Variante alone).
# Aufruf: UBOOT=<binary> WORK=<messordner> ./07-fremd-tap.sh      -> rohdaten-m6.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-m6.txt

TAP_A=shivammathur/extensions
TAP_B=shivammathur/php
F_A='shivammathur/extensions/xdebug@8.2'
F_B='shivammathur/php/php@8.2'

echo "== M6 Abhaengigkeit aus einem anderen Tap, $(now)"
echo "u-boot: $("$UBOOT" --version 2>&1 | head -1)"
echo "formel=$F_A haengt ab von $F_B"

mark "(a) Bau: nur $TAP_A freigegeben, Liste: $F_A"
projekt fremd 0 "devcontainer.brewTaps.allow=$TAP_A" "devcontainer.features.brew.enabled=true" "devcontainer.features.brew.packages=$F_A"
grep -n 'brewTaps' -A2 "$(vdir fremd)/u-boot.yaml"
if bild fremd --no-cache; then
  aussage "m6-bau=ok"
else
  aussage "m6-bau=FEHLER"
fi
aussage "m6-bau-image-vorhanden=$(docker image inspect "$(vimg fremd)" > /dev/null 2>&1 && echo ja || echo nein)"
# Die Meldung von Homebrew steht nur in der Ausgabe ohne -q.
docker build --progress=plain --no-cache -t "$(vimg fremd)-plain" "$(vdir fremd)/.devcontainer" > "$WORK/bau-fremd-plain.log" 2>&1
echo "-- Ausgabe des Schritts brew install (Bau ohne -q), Auszug"
grep -E "Tapping|Trusted formula|Warning|Error|This command requires|brew tap shivammathur|did not complete" "$WORK/bau-fremd-plain.log" | sed 's/^#[0-9]* [0-9.]* //'
aussage "m6-bau-meldung-tap-noetig=$(grep -c "This command requires the tap $TAP_B" "$WORK/bau-fremd-plain.log")"
aussage "m6-bau-tap-b-eingerichtet=$(grep -c "Tapping $TAP_B" "$WORK/bau-fremd-plain.log")"

alone=$(vimg alone)
docker image inspect "$alone" > /dev/null 2>&1 || { echo "Abbruch: das Image der Variante alone fehlt (zuerst 01-bau.sh)"; exit 2; }
for order in "B A" "A B"; do
  case $order in "B A") args=("$F_B" "$F_A");; *) args=("$F_A" "$F_B");; esac
  mark "(b) --dry-run in einem Aufruf, Reihenfolge: ${args[*]}"
  docker run --rm "$alone" bash -c '
    export HOMEBREW_NO_ANALYTICS=1 HOMEBREW_NO_AUTO_UPDATE=1
    echo "taps vorher: [$(brew tap | tr "\n" " ")]"
    brew install --formula --dry-run "$@" 2>&1 | grep -E "^(==> (Tapping|Trusted|Would)|Error|Warning|This command|If you trust|  brew tap)"
    echo "exit=${PIPESTATUS[0]}"
    echo "taps nachher: [$(brew tap | tr "\n" " ")]"' _ "${args[@]}"
done
echo "== ende $(now)"
