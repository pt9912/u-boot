#!/bin/bash
# Laeuft IM Container als root (Basis-Image); aufgerufen von 03-lieferkette.sh. Liest nur Paketindizes; installiert nichts.
set -u
. /etc/os-release; echo "$PRETTY_NAME, arch=$(dpkg --print-architecture)"
apt-get update -qq > /dev/null 2>&1; echo "apt-get update exit=$?"
echo "-- Paketquellen des Basis-Images"
grep -h -E '^(URIs|Suites|Components|Signed-By):' /etc/apt/sources.list.d/*.sources 2>/dev/null
echo "-- Kandidaten gh und glab"
apt-cache policy gh glab | grep -E '^[a-z]|Candidate|^     [0-9]'
echo "-- Abhaengigkeiten der Pakete"
apt-cache depends gh glab | head -20
for p in gh glab; do
  echo "-- Debian-Paket $p: Herkunft der Quelle (apt-cache show)"
  apt-cache show "$p" | grep -E '^(Package|Version|Maintainer|Installed-Size|Size|Homepage):'
done
