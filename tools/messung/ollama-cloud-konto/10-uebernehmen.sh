#!/bin/bash
# Stufe 2 der Rohdaten (DoD 1, iv): kopiert Dateien aus Stufe 1 ($ROH, ausserhalb des Repos) in das Verzeichnis des Messprotokolls, NACHDEM sie von Hand
# gelesen wurden (Zeile fuer Zeile). Das Skript prueft das Lesen nicht; es verlangt die ausdrueckliche Option --gelesen, kopiert nur Dateinamen
# rohdaten-*.txt und fuehrt danach 09-suche.sh ueber das Zielverzeichnis aus; bei einem Treffer werden die eben kopierten Dateien wieder entfernt.
# Aufruf: WORK=<messordner> tools/messung/ollama-cloud-konto/10-uebernehmen.sh --gelesen <datei>...    (Dateinamen ohne Pfad, aus $ROH)
set -u
. "$(dirname "$0")/lib.sh"
[ "${1:-}" = --gelesen ] || { echo "Abbruch: Stufe 2 nur mit --gelesen (die Dateien sind von Hand gelesen)"; exit 2; }
shift
[ $# -gt 0 ] || { echo "Abbruch: keine Dateien genannt"; exit 2; }
mkdir -p "$BASEDIR"
kopiert=()
for f in "$@"; do
  case "$f" in rohdaten-*.txt) ;; *) echo "Abbruch: nur rohdaten-*.txt ($f)"; exit 2 ;; esac
  case "$f" in */*) echo "Abbruch: Dateiname ohne Pfad ($f)"; exit 2 ;; esac
  [ -f "$ROH/$f" ] || { echo "Abbruch: $ROH/$f fehlt"; exit 2; }
done
for f in "$@"; do cp "$ROH/$f" "$BASEDIR/$f"; kopiert+=("$f"); done
if "$here/09-suche.sh" "$BASEDIR"; then
  echo "uebernommen: ${kopiert[*]}"
else
  for f in "${kopiert[@]}"; do rm -f "$BASEDIR/$f"; done
  echo "Abbruch: die Suche hat angeschlagen; die Dateien wurden wieder entfernt"; exit 1
fi
