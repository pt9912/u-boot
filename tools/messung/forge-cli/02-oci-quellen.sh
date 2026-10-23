#!/bin/bash
# Quellen als OCI-Feature: gh (offiziell, Dritt-Sammlung) und glab (drei Drittanbieter-Features aus dem Feature-Index).
# Aufruf: WORK=<messordner> ./02-oci-quellen.sh      -> rohdaten-02.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-02.txt

echo "== Quellen als OCI-Feature, $(now)"
echo "-- Feature-Index https://containers.dev/static/devcontainer-index.json: Eintraege zu gh und glab"
tool bash /m/inner-02-index.sh
aussage "index-abgefragt=$(date -u +%F)"

for v in gh ghx glabs glabp glabz; do
  echo; tool bash /m/inner-02.sh "$(oci_ref $v)" "$v"
done

echo
echo "-- Sondierung weiterer Namen (tags/list; 404 = unter diesem Namen kein Paket, KEINE Aussage ueber andere Namen)"
tool bash /m/inner-02-sonde.sh
