#!/bin/bash
# Laeuft IM Werkzeug-Container; aufgerufen von 02-oci-quellen.sh. Liest den oeffentlichen Feature-Index.
set -u
curl -sS -m 60 -o /w/index.json -w 'http=%{http_code} bytes=%{size_download}\n' https://containers.dev/static/devcontainer-index.json
echo "-- Eintraege, deren id gitlab-cli, glab, github-cli, gh-cli oder gh enthaelt (Sammlung, id, Version)"
jq -r '.collections[] | .sourceInformation.name as $c | .features[]? | select((.id//"")|test("gitlab-cli|glab|github-cli|gh-cli|^gh$")) | [$c, .id, .version] | @tsv' /w/index.json
echo "-- Eintraege mit glab oder GitLab CLI in id, Name oder Beschreibung (alle Sammlungen)"
jq -r '.collections[] | .features[]? | select(((.id//"")+(.name//"")+(.description//""))|test("glab|gitlab cli";"i")) | .id' /w/index.json | sort -u
