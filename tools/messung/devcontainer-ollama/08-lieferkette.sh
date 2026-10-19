#!/bin/bash
# O2 und O7: Lieferkette des Upstream-Releases (Signatur, Provenance, Pruefsummen, Veraenderlichkeit, Pin-Alter).
# Aufruf: 08-lieferkette.sh > rohdaten-08.txt 2>&1   (braucht gh mit Anmeldung, jq, curl auf dem Host)
set -u
tag=v0.40.6
t=$(mktemp -d)
echo "== Release $tag"
gh api repos/ollama/ollama/releases/tags/$tag --jq '{tag_name, published_at, immutable, prerelease, target_commitish, assets: (.assets|length)}'
echo "== Assets (Name, Groesse, GitHub-Digest) mit Signatur-/Provenance-Muster"
gh api repos/ollama/ollama/releases/tags/$tag --jq '.assets[]|"\(.name) \(.size) \(.digest)"' | sort
echo "-- Assets, deren Name nach Signatur, Provenance oder SBOM aussieht:"
gh api repos/ollama/ollama/releases/tags/$tag --jq '[.assets[].name|select(test("\\.(asc|sig|pem|sigstore|intoto\\.jsonl|sbom|spdx|cdx)|attest|cosign|provenance|sbom";"i"))]|length'
echo "== sha256sum.txt gegen GitHub-Digest (alle Assets)"
curl -fsSL https://github.com/ollama/ollama/releases/download/$tag/sha256sum.txt -o /dev/stdout | sed 's|  \./|  |' | sort -k2 > "$t/sums.txt"
gh api repos/ollama/ollama/releases/tags/$tag --jq '.assets[]|select(.name!="sha256sum.txt")|"\(.digest|sub("sha256:";""))  \(.name)"' | sort -k2 > "$t/digests.txt"
if diff "$t/sums.txt" "$t/digests.txt"; then echo "gleich ($(wc -l < "$t/sums.txt") Zeilen)"; else echo "ABWEICHUNG"; fi
rm -rf "$t"
echo "== amd64-Archiv: Pruefsumme im Strom (ohne Datei)"
curl -fsSL https://github.com/ollama/ollama/releases/download/$tag/ollama-linux-amd64.tar.zst | sha256sum
echo "(erwartet c164ab8b1b28e2a247c29719b8e404c423c5b83ae64c3ab4e59d50c7ff7fc5e0)"
echo "== GitHub-Attestation (API) fuer die beiden Linux-Archive"
for d in c164ab8b1b28e2a247c29719b8e404c423c5b83ae64c3ab4e59d50c7ff7fc5e0 0a279d9428ad17b06fdb7dddf3b059b2c148433585fa3b30676a102fc459bb9f; do
  printf '%s: ' "$d"; gh api repos/ollama/ollama/attestations/sha256:$d 2>&1 | head -c 120; echo
done
echo "== Tag-Commit"
sha=$(gh api repos/ollama/ollama/git/ref/tags/$tag --jq .object.sha)
echo "tag-objekt: $(gh api repos/ollama/ollama/git/ref/tags/$tag --jq .object.type) $sha"
gh api repos/ollama/ollama/git/commits/$sha --jq '{verified: .verification.verified, reason: .verification.reason, author_date: .author.date, committer_date: .committer.date}'
echo "== Release-Workflow: Signatur-, Attestierungs- und Pruefsummen-Schritte (Zeilen)"
gh api repos/ollama/ollama/contents/.github/workflows/release.yaml --jq .content | base64 -d | grep -n -i -E 'attest|cosign|sigstore|gpg|id-token|provenance|sha256sum|clobber'
echo "== neuere Releases (Stand der Messung)"
gh api 'repos/ollama/ollama/releases?per_page=8' --jq '.[]|"\(.tag_name) \(.published_at) prerelease=\(.prerelease) immutable=\(.immutable)"'
echo "== Docker-Hub-Image: Attestierungen im Index (Anzahl Manifeste, Plattformen)"
docker buildx imagetools inspect --raw ollama/ollama:0.40.6 | jq -c '[.manifests[]|.platform]'
