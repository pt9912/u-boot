#!/bin/bash
# Vergleicht die stabilen Aussagen (Zeilen "AUSSAGE ...") zweier Laeufe der Skripte.
# Aufruf: tools/messung/ollama-amd64/vergleich.sh <ausgabeverzeichnis-1> <ausgabeverzeichnis-2>
# Ausgenommen sind Zeilen, die je Lauf andere Zahlen tragen duerfen: die Zaehlungen der Wiederholung A7 (A6/A7 unterscheiden sich je
# Lauf in der Zahl der Treffer) und der Pull des verschachtelten Docker hinter der Restriktion (A6; zeitabhaengig, README).
# Exit 0, wenn die Mengen uebereinstimmen; sonst 1 und der Unterschied.
set -u
a=${1:?verzeichnis 1}; b=${2:?verzeichnis 2}
holen() { grep -h '^AUSSAGE ' "$1"/rohdaten-*.txt | grep -v -E 'laeufe-mit-|sonden-ausserhalb|nested-pull-mit-egress' | sort; }
holen "$a" > "$a/.aussagen.tmp"; holen "$b" > "$b/.aussagen.tmp"
if diff "$a/.aussagen.tmp" "$b/.aussagen.tmp"; then echo "gleich: $(grep -c . "$a/.aussagen.tmp") Aussagen"; rc=0; else rc=1; fi
rm -f "$a/.aussagen.tmp" "$b/.aussagen.tmp"
exit $rc
