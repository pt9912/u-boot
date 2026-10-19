#!/bin/bash
# Vergleicht die stabilen Aussagen (Zeilen "AUSSAGE ...") zweier Laeufe der Skripte.
# Aufruf: tools/messung/ollama-amd64/vergleich.sh <lauf-1> <lauf-2>
#   <lauf-N> ist ein Ausgabeverzeichnis (OUT) mit rohdaten-*.txt oder eine Datei rohdaten-lauf2-aussagen.txt (Zeilen "<rohdatei> AUSSAGE ...").
#   Dateien rohdaten-lauf2-* in einem Verzeichnis zaehlen nicht mit.
# Ausgenommen sind Zeilen, die je Lauf andere Zahlen tragen duerfen: die Zaehlungen der Wiederholung A7 und der Pull des verschachtelten
# Docker hinter der Restriktion (A6; zeitabhaengig, README).
# Exit 0, wenn die Mengen uebereinstimmen; sonst 1 und der Unterschied.
set -u
holen() { # <verzeichnis|datei>
  if [ -d "$1" ]; then
    for f in "$1"/rohdaten-*.txt; do
      case "$(basename "$f")" in rohdaten-lauf2-*) continue ;; esac
      grep -h '^AUSSAGE ' "$f" | while read -r l; do echo "$(basename "$f") $l"; done
    done
  else
    cat "$1"
  fi | grep -v -E 'laeufe-mit-|sonden-ausserhalb|nested-pull-mit-egress' | sort
}
a=$(mktemp); b=$(mktemp)
holen "${1:?lauf 1}" > "$a"; holen "${2:?lauf 2}" > "$b"
if diff "$a" "$b"; then echo "gleich: $(grep -c . "$a") Aussagen"; rc=0; else rc=1; fi
rm -f "$a" "$b"
exit $rc
