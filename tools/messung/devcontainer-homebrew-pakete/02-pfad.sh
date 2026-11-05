#!/bin/bash
# M2 PATH und Sitzungsarten: steht `brew` im PATH des Image-Benutzers in Login- und Nicht-Login-Shells? Gemessen im Image der Variante def
# (01-bau.sh) fuer bash, sh und zsh, jeweils nicht interaktiv ohne und mit -l und interaktiv ohne und mit -l, und fuer die Sitzung, die
# `docker exec` startet. Die Shell von `u-boot devcontainer shell` ist keine Login-Shell (Beobachtung der Spezifikation des Befehls).
# Aufruf: UBOOT=<binary> WORK=<messordner> ./02-pfad.sh      -> rohdaten-m2.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-m2.txt

echo "== M2 PATH, $(now)"
img=$(vimg def)
docker image inspect "$img" > /dev/null 2>&1 || { echo "Abbruch: das Image der Variante def fehlt (zuerst 01-bau.sh)"; exit 2; }
echo "-- Umgebung des Images (PATH und HOMEBREW*)"
docker run --rm "$img" bash -c 'env | grep -E "^(PATH|HOMEBREW)"'
echo "-- Shells im Image"
docker run --rm "$img" bash -c 'for s in bash sh zsh; do printf "%s=%s\n" "$s" "$(command -v $s || echo fehlt)"; done'
for mode in "bash -c" "bash -lc" "bash -ic" "bash -lic" "sh -c" "sh -lc" "zsh -c" "zsh -lc" "zsh -ic" "zsh -lic"; do
  r=$(docker run --rm "$img" $mode 'command -v brew || echo MISSING' 2>&1 | tail -1)
  printf '%-10s %s\n' "$mode" "$r"
  aussage "m2-$(echo "$mode" | tr ' ' '_')=$r"
done
echo "-- docker exec in einen laufenden Container (Sitzung ohne Login-Shell)"
c="$PFX-pfad"
docker rm -f "$c" > /dev/null 2>&1
docker run -d --name "$c" "$img" sleep 600 > /dev/null
docker exec "$c" bash -c 'command -v brew || echo MISSING'
docker exec -u vscode "$c" sh -c 'command -v brew || echo MISSING'
echo "-- Login-Shell: Quelle des Verlusts (Debian setzt PATH in /etc/profile neu)"
docker exec "$c" bash -c 'grep -n "PATH" /etc/profile; ls /etc/profile.d/'
docker rm -f "$c" > /dev/null
echo "== ende $(now)"
