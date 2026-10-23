#!/bin/bash
# Bau der Dockerfile-Block-Varianten (HAND-Varianten): das von u-boot erzeugte Sandbox-Dockerfile der Variante bas, dazu ein Block vor der
# letzten Zeile `USER vscode`. Der Katalog hat fuer gh und glab keinen solchen Eintrag; der Block hat die Form des Ollama-Blocks
# (Dockerfile.tmpl). Gebaut mit `docker build` (nicht mit der Dev-Container-CLI; es gibt keine Features), kalt (--no-cache).
# Varianten: basd (Referenz, unveraendert), ghdeb, glabdeb (Debian-Paket), ghtgz, glabtgz (gepinntes Release-Archiv, SHA-256 geprueft).
# Voraussetzung: 04-projekte.sh. Aufruf: WORK=<messordner> ./06-dockerfile-bloecke.sh      -> rohdaten-06.txt, rohdaten-06-<variante>.log
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-06.txt

echo "== Dockerfile-Bloecke, $(now)"
for v in basd ghdeb glabdeb ghtgz glabtgz; do
  rm -rf "$(vdir "$v")"; cp -a "$(vdir bas)" "$(vdir "$v")"
  f="$(vdir "$v")/.devcontainer/Dockerfile"
  if [ "$v" != basd ]; then
    cp "$f" "$WORK/Dockerfile.erzeugt"
    { head -n -2 "$WORK/Dockerfile.erzeugt"; cat "$here/block-$v.dockerfile"; tail -n 2 "$WORK/Dockerfile.erzeugt"; } > "$f"
  fi
  echo; mark "docker build $v (kalt)"
  s=$(date +%s.%N)
  docker build --progress=plain --no-cache -t "$(vimg "$v")" "$(vdir "$v")/.devcontainer" > "$WORK/bau-$v.log" 2>&1; rc=$?
  e=$(date +%s.%N)
  redact < "$WORK/bau-$v.log" > "$OUT/rohdaten-06-$v.log"
  echo "build $v exit=$rc dauer=$(sekunden "$s" "$e")s"
  aussage "bau-$v-exit=$rc"
  [ "$rc" = 0 ] || { tail -12 "$WORK/bau-$v.log" | cut -c1-240; continue; }
  docker image inspect "$(vimg "$v")" --format 'image={{.Id}} size={{.Size}} arch={{.Architecture}}'
  docker run --rm "$(vimg "$v")" bash -c '
    for t in gh glab; do
      p=$(command -v $t) && { echo "$t: $p"; $t --version 2>&1 | head -1; stat -c "%U:%G %a %s" "$p"; } || echo "$t: nicht im Image"
    done'
  case "$v" in *tgz) echo "-- Zeilen des Builds zur Pruefsumme"; grep -E 'OK$|FAILED' "$WORK/bau-$v.log" | head -3 ;; esac
done

echo
echo "== Imagegroesse gegen basd (Bytes, Differenz)"
b=$(docker image inspect "$(vimg basd)" --format '{{.Size}}')
for v in ghdeb glabdeb ghtgz glabtgz; do
  s=$(docker image inspect "$(vimg "$v")" --format '{{.Size}}' 2>/dev/null) || continue
  echo "$v size=$s differenz=$((s - b))"
  aussage "groesse-differenz-$v=$((s - b))"
done

echo
echo "== Negativprobe: falsche Pruefsumme bricht den Bau ab (ghtgz, SHA-256 mit geaenderter erster Stelle)"
v=ghtgz; d="$WORK/p-negativ"; rm -rf "$d"; cp -a "$(vdir "$v")" "$d"
sed -i 's/ARG GH_SHA256_AMD64=f/ARG GH_SHA256_AMD64=e/' "$d/.devcontainer/Dockerfile"
docker build --progress=plain --no-cache -t "$PFX-img-negativ" "$d/.devcontainer" > "$WORK/bau-negativ.log" 2>&1; rc=$?
echo "build negativ exit=$rc"
grep -E 'FAILED|did NOT match|WARNING: 1 computed' "$WORK/bau-negativ.log" | head -3
aussage "negativprobe-pin-exit=$rc"
docker image inspect "$PFX-img-negativ" > /dev/null 2>&1 && echo "Image entstanden: ja" || echo "Image entstanden: nein"
