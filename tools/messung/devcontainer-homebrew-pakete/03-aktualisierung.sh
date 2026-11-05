#!/bin/bash
# M3 Selbst-Aktualisierung: verschiebt `brew install` die Revision des Homebrew-Verzeichnisses? Gemessen an einem Image, das auf den Commit des
# vorherigen Releases (7.0.10) gebaut ist (Build-Argument HOMEBREW_REVISION); die Datei .git/FETCH_HEAD wird auf drei Tage alt gesetzt (die
# Automatik prueft ihr Alter), dann installiert `brew install` eine Formel (a) wie der Block (HOMEBREW_NO_AUTO_UPDATE=1) und (b) ohne den Schalter.
# Aufruf: UBOOT=<binary> WORK=<messordner> ./03-aktualisierung.sh      -> rohdaten-m3.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-m3.txt

echo "== M3 Selbst-Aktualisierung, $(now)"
echo "pin-alt=$PIN_ALT pin-template=$PIN"
projekt alt 0 "devcontainer.features.brew.enabled=true"
bild alt --build-arg HOMEBREW_REVISION="$PIN_ALT" || exit 1
img=$(vimg alt)
script='set -u
repo=/home/linuxbrew/.linuxbrew/Homebrew
echo "vorher=$(git -C "$repo" rev-parse HEAD)"
touch -d "3 days ago" "$repo/.git/FETCH_HEAD"
env $ENVS brew install --formula jq 2>&1 | grep -E "Auto-updating|Updated|error|Error|jq" | head -8
echo "nachher=$(git -C "$repo" rev-parse HEAD)"
echo "branch=$(git -C "$repo" branch --show-current)"'
for run in "mit-schalter:HOMEBREW_NO_AUTO_UPDATE=1 HOMEBREW_NO_ANALYTICS=1" "ohne-schalter:HOMEBREW_NO_ANALYTICS=1"; do
  name=${run%%:*}; envs=${run#*:}
  mark "$name ($envs)"
  out=$(docker run --rm -e ENVS="$envs" "$img" bash -c "$script" 2>&1); echo "$out"
  after=$(echo "$out" | sed -n 's/^nachher=//p')
  if [ "$after" = "$PIN_ALT" ]; then aussage "m3-$name=revision-unveraendert"; else aussage "m3-$name=revision-verschoben-auf-$after"; fi
done
echo "== ende $(now)"
