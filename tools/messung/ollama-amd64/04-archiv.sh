#!/bin/bash
# A2: --exclude-Muster im amd64-Archiv: Archivliste, Entpacken mit und ohne Muster, Muster ohne Treffer, Gegenprobe mit falschem Muster.
# Voraussetzung: 03-bau.sh (Image oll).
# Aufruf: WORK=<messordner> tools/messung/ollama-amd64/04-archiv.sh      -> $OUT/rohdaten-04.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-04.txt

TARIMG=$PFX-tar
echo "== A2 Archiv und --exclude-Muster, $(now)"
mkdir -p "$WORK/archiv"
if [ ! -s "$WORK/archiv/$ARCHIVE" ]; then
  s=$(date +%s.%N)
  tool curl -fsSL -o "archiv/$ARCHIVE" "$ARCHIVE_URL"; rc=$?
  e=$(date +%s.%N)
  echo "download exit=$rc dauer=$(sekunden "$s" "$e")s"
fi
tool sh -c "ls -l archiv/$ARCHIVE | awk '{print \"archivgroesse-bytes=\" \$5}'; echo \"${SHA_AMD64}  archiv/$ARCHIVE\" | sha256sum -c -"
aussage "archiv-sha256=$(tool sha256sum "archiv/$ARCHIVE" | cut -d' ' -f1)"

docker build -q -f "$here/Dockerfile.tar" -t "$TARIMG" "$here" > /dev/null || exit 1
T() { docker run --rm -u root -v "$WORK/archiv":/a:ro "$TARIMG" bash -c "$1"; }

echo "-- tar-Version"; T 'tar --version | head -1'
echo "-- Eintraege des Archivs je Wurzelverzeichnis"
T "tar --zstd -tf /a/$ARCHIVE | awk -F/ '{print \$1\"/\"\$2}' | sort | uniq -c"
echo "-- Eintraege unter lib/ollama/ je Unterverzeichnis (Dateien direkt unter lib/ollama/ als '.')"
T "tar --zstd -tf /a/$ARCHIVE | awk -F/ '\$1==\"lib\" && \$2==\"ollama\" && NF>=3 {print (NF==3 ? \".\" : \$3)}' | sort | uniq -c"
aussage "archiv-hat-cuda-verzeichnisse=$(T "tar --zstd -tf /a/$ARCHIVE | grep -c '^lib/ollama/cuda_'")"
aussage "archiv-hat-vulkan-verzeichnis=$(T "tar --zstd -tf /a/$ARCHIVE | grep -c '^lib/ollama/vulkan'")"

echo; echo "-- Entpacken mit den Mustern des Generators (Pfad: -C /x1)"
T "mkdir /x1 && tar --zstd -xf /a/$ARCHIVE -C /x1 --no-same-owner --exclude='lib/ollama/cuda_*' --exclude='lib/ollama/vulkan'; echo tar-exit=\$?; du -sm /x1 | cut -f1 | sed 's/^/entpackt-mib=/'; ls /x1/lib/ollama | grep -c -E '^(cuda_|vulkan)' | sed 's/^/ausgeschlossene-treffer=/'; ls -d /x1/lib/ollama/*/ 2>/dev/null || echo 'keine Unterverzeichnisse unter lib/ollama'"
echo "-- Entpacken ohne Muster (Pfad: -C /x2)"
T "mkdir /x2 && tar --zstd -xf /a/$ARCHIVE -C /x2 --no-same-owner; echo tar-exit=\$?; du -sm /x2 | cut -f1 | sed 's/^/entpackt-mib=/'; ls -d /x2/lib/ollama/*/ 2>/dev/null || echo 'keine Unterverzeichnisse unter lib/ollama'"

echo; echo "-- Vergleich: Inhalt von /usr/local/lib/ollama im gebauten Image (03) gegen Entpacken mit Mustern"
docker run --rm "$(vimg oll)" bash -c 'cd /usr/local && find lib/ollama bin/ollama -type f | sort' > "$WORK/image-dateien.txt"
T "mkdir /x3 && tar --zstd -xf /a/$ARCHIVE -C /x3 --no-same-owner --exclude='lib/ollama/cuda_*' --exclude='lib/ollama/vulkan' && cd /x3 && find lib/ollama bin/ollama -type f | sort" > "$WORK/archiv-dateien.txt"
if diff "$WORK/image-dateien.txt" "$WORK/archiv-dateien.txt" > /dev/null; then echo "Dateiliste Image == Dateiliste Archiv mit Mustern ($(grep -c . "$WORK/image-dateien.txt") Dateien)"; aussage "image-gleich-archiv-mit-mustern=ja"; else echo "UNTERSCHIED"; diff "$WORK/image-dateien.txt" "$WORK/archiv-dateien.txt" | head; aussage "image-gleich-archiv-mit-mustern=nein"; fi

echo; echo "-- Muster OHNE Treffer (Kette wie im Dockerfile: tar && naechster Schritt)"
T "mkdir /x4 && (tar --zstd -xf /a/$ARCHIVE -C /x4 --no-same-owner --exclude='lib/ollama/gibt_es_nicht_*' --exclude='lib/ollama/auch_nicht' && echo KETTE_LAEUFT_WEITER); echo kette-exit=\$?" 2>&1
T "mkdir /x5; tar --zstd -xf /a/$ARCHIVE -C /x5 --no-same-owner --exclude='lib/ollama/gibt_es_nicht_*' > /tmp/o.txt 2>&1; echo tar-exit=\$?; echo stdout-stderr-bytes=\$(wc -c < /tmp/o.txt)"
aussage "muster-ohne-treffer-tar-exit=$(T "mkdir /x6 && tar --zstd -xf /a/$ARCHIVE -C /x6 --no-same-owner --exclude='lib/ollama/gibt_es_nicht_*' >/dev/null 2>&1; echo \$?")"

echo; echo "-- Gegenprobe: falsch geschriebenes Muster schliesst NICHTS aus (zeigt, dass die Muster wirken)"
T "mkdir /x7 && tar --zstd -xf /a/$ARCHIVE -C /x7 --no-same-owner --exclude='lib/ollama/cudaa_*' --exclude='lib/ollama/vulkann'; echo tar-exit=\$?; du -sm /x7 | cut -f1 | sed 's/^/entpackt-mib=/'; ls -d /x7/lib/ollama/*/ 2>/dev/null || echo 'keine Unterverzeichnisse'"
echo "-- Gegenprobe: das Muster 'lib/ollama/cuda_*' ohne Anker (kein Praefix 'lib/ollama/') trifft ebenfalls (GNU tar: Muster passen auf Namensteile)"
T "mkdir /x8 && tar --zstd -xf /a/$ARCHIVE -C /x8 --no-same-owner --exclude='cuda_*' --exclude='vulkan'; ls -d /x8/lib/ollama/*/ 2>/dev/null || echo 'keine Unterverzeichnisse'"
