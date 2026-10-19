#!/bin/bash
# A1: Bau des vom Generator erzeugten Dockerfiles (Variante oll) kalt (--no-cache) und aus dem Cache; Gegenprobe ohne Ollama (bas).
# Aufruf: WORK=<messordner> tools/messung/ollama-amd64/03-bau.sh
#   -> $OUT/rohdaten-03.txt, rohdaten-03-kalt.log, rohdaten-03-cache.log, rohdaten-03-bas.log
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-03.txt

bau() { # <variante> <log> [docker-build-argumente...]
  local v=$1 log=$2; shift 2
  local s e rc
  s=$(date +%s.%N)
  docker build --progress=plain "$@" -t "$(vimg "$v")" "$(vdir "$v")/.devcontainer" > "$WORK/$log" 2>&1; rc=$?
  e=$(date +%s.%N)
  echo "build $v $* exit=$rc dauer=$(sekunden "$s" "$e")s"
  redact < "$WORK/$log" > "$OUT/rohdaten-03-$log.log"
  return $rc
}

echo "== A1 Bau, $(now)"
echo "-- Groesse des Downloads (Metrik: Content-Length des Archivs, nicht Image-Groesse)"
tool sh -c "curl -sIL '$ARCHIVE_URL' | grep -i '^content-length' | tail -1" | tr -d '\r'
aussage "archiv-name=$ARCHIVE"

mark "Bau kalt (--no-cache), Variante oll"
bau oll kalt --no-cache; rc=$?
aussage "bau-kalt-exit=$rc"
[ "$rc" = 0 ] || { tail -20 "$WORK/kalt"; exit 1; }
echo "-- Zeilen des Builds zur Pruefsumme und zum Archiv"
grep -E "$ARCHIVE|OK$|FAILED" "$WORK/kalt" | grep -v 'RUN ' | head
aussage "sha256sum-c=$(grep -c "$ARCHIVE: OK" "$WORK/kalt")"
mark "Bau aus dem Cache, Variante oll"
bau oll cache; rc=$?
aussage "bau-cache-exit=$rc"
echo "-- CACHED-Schritte im zweiten Bau: $(grep -c ' CACHED' "$WORK/cache") von $(grep -c -E '^#[0-9]+ \[' "$WORK/cache") Schritten"

echo "-- Image (Metrik .Size, Bytes und GB)"
docker image inspect "$(vimg oll)" --format '{{.Id}} size={{.Size}} arch={{.Architecture}}'
docker image inspect "$(vimg oll)" --format '{{.Size}}' | awk '{printf "image-size-gb(10^9)=%.2f\n", $1/1e9}'
echo "-- Inhalt des Images (als Benutzer des Images)"
docker run --rm "$(vimg oll)" bash -c '
  id
  echo "-- ollama --version"; ollama --version 2>&1
  echo "-- Eigentuemer"; stat -c "%U:%G %n" /usr/local/bin/ollama /home/vscode/.ollama
  echo "-- lib/ollama"; ls -1 /usr/local/lib/ollama
  echo "-- zstd im Image"; command -v zstd || echo "kein zstd"
  echo "-- Groesse /usr/local"; du -sh /usr/local/lib/ollama /usr/local/bin/ollama
  echo "-- OLLAMA_*-Umgebung"; env | grep -i "^OLLAMA" || echo "keine OLLAMA_*-Variable"'
aussage "ollama-version-im-image=$(docker run --rm "$(vimg oll)" ollama --version 2>&1 | grep -o 'version is .*')"
aussage "eigentuemer-ollama-dir=$(docker run --rm "$(vimg oll)" stat -c '%U:%G' /home/vscode/.ollama)"
aussage "zstd-im-image=$(docker run --rm "$(vimg oll)" bash -c 'command -v zstd >/dev/null && echo ja || echo nein')"

mark "Gegenprobe: Variante bas (ohne Ollama-Block), Paketliste gegen oll"
bau bas bas || exit 1
docker run --rm "$(vimg bas)" dpkg-query -W -f '${Package} ${Version}\n' | sort > "$WORK/dpkg-bas.txt"
docker run --rm "$(vimg oll)" dpkg-query -W -f '${Package} ${Version}\n' | sort > "$WORK/dpkg-oll.txt"
echo "Pakete bas=$(grep -c . "$WORK/dpkg-bas.txt") oll=$(grep -c . "$WORK/dpkg-oll.txt")"
echo "-- diff dpkg bas -> oll (leer: der Ollama-Block hinterlaesst keine Paketaenderung)"
diff "$WORK/dpkg-bas.txt" "$WORK/dpkg-oll.txt" && echo "(kein Unterschied)"
aussage "paketdiff-bas-oll=$(diff "$WORK/dpkg-bas.txt" "$WORK/dpkg-oll.txt" | grep -c '^[<>]')"
docker image inspect "$(vimg bas)" --format 'bas size={{.Size}}'
echo "-- Plattenplatz nach dem Bau"
docker system df | head -3
