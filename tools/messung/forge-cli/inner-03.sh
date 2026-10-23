#!/bin/bash
# Laeuft IM Werkzeug-Container (Basis-Image, anonym, ohne Token); aufgerufen von 03-lieferkette.sh.
# Misst die Lieferkette der Release-Artefakte von gh (cli/cli, GitHub) und glab (gitlab-org/cli, GitLab).
set -u
t=/w/lieferkette; rm -rf "$t"; mkdir -p "$t"; cd "$t" || exit 1
export GNUPGHOME="$t/gnupg"; mkdir -m 700 "$GNUPGHOME"

echo "################ gh (github.com/cli/cli)"
curl -sS https://api.github.com/repos/cli/cli/releases/latest > gh-release.json
echo "-- neuestes Release (Tag, Veroeffentlichung, immutable, Prerelease, Zahl der Assets)"
jq -c '{tag_name, published_at, immutable, prerelease, assets: (.assets|length)}' gh-release.json
tag=$(jq -r .tag_name gh-release.json); ver=${tag#v}
echo "-- Assets fuer Linux (Name, Groesse, GitHub-Digest)"
jq -r '.assets[]|select(.name|test("linux_(amd64|arm64)\\.(tar\\.gz|deb)$|checksums"))|"\(.name) \(.size) \(.digest)"' gh-release.json
echo "-- Assets, deren Name nach Signatur, Provenance oder SBOM aussieht (Anzahl)"
jq -r '[.assets[].name|select(test("\\.(asc|sig|pem|sigstore|intoto\\.jsonl|sbom|spdx|cdx)|attest|cosign|provenance|sbom|bundle";"i"))]|length' gh-release.json
a="gh_${ver}_linux_amd64.tar.gz"
curl -fsSL "https://github.com/cli/cli/releases/download/$tag/$a" -o "$a"
curl -fsSL "https://github.com/cli/cli/releases/download/$tag/gh_${ver}_checksums.txt" -o gh-checksums.txt
echo "-- Archiv $a: sha256 der Datei, Eintrag in checksums.txt, GitHub-Digest des Assets"
d=$(sha256sum "$a" | cut -d' ' -f1)
echo "datei:    $d"
echo "checksums: $(awk -v f="$a" '$2 == f {print $1}' gh-checksums.txt)"
echo "api:      $(jq -r --arg n "$a" '.assets[]|select(.name==$n)|.digest' gh-release.json)"
aussage() { echo "AUSSAGE $*"; }
aussage "gh-tag=$tag"
aussage "gh-amd64-sha256=$d"
aussage "gh-checksums-gleich-datei=$([ "$(awk -v f="$a" '$2 == f {print $1}' gh-checksums.txt)" = "$d" ] && echo ja || echo nein)"
echo "-- GitHub-Attestierung (API, anonym) fuer das Archiv"
code=$(curl -sS -o gh-att.json -w '%{http_code}' "https://api.github.com/repos/cli/cli/attestations/sha256:$d")
echo "http=$code"
jq -c '{attestations: (.attestations|length), inline_bundle: [.attestations[]|(.bundle != null)], initiator: [.attestations[].initiator]}' gh-att.json 2>/dev/null || head -c 200 gh-att.json
if [ "$code" = 200 ]; then
  echo "-- Attestierungen: Bundle liegt hinter einer signierten Adresse der Antwort (Adresse nicht ausgegeben); Format der Bundle-Datei:"
  n=$(jq '.attestations|length' gh-att.json)
  for i in $(seq 0 $((n - 1))); do
    curl -fsSL "$(jq -r ".attestations[$i].bundle_url" gh-att.json)" -o "gh-bundle-$i.bin"
    echo "bundle $i: $(wc -c < "gh-bundle-$i.bin") Bytes, JSON lesbar: $(jq -e . "gh-bundle-$i.bin" > /dev/null 2>&1 && echo ja || echo nein)"
  done
  echo "(Inhalt nicht dekodiert: die Prueflogik von gh attestation verify ist nicht gemessen)"
fi
echo "-- Release-Workflow des Projekts: Zeilen zu Signatur, Attestierung, Pruefsumme"
curl -sS https://raw.githubusercontent.com/cli/cli/trunk/.github/workflows/deployment.yml -o gh-deploy.yml
echo "http-datei-groesse=$(wc -c < gh-deploy.yml)"
grep -n -i -E 'attest|cosign|sigstore|gpg|id-token|provenance|sha256sum|checksum|immutable' gh-deploy.yml | cut -c1-180 | head -30
echo "-- apt-Repository des Projekts: Schluessel, Fingerabdruck, Signatur von InRelease"
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg -o gh-keyring.gpg
gpg --show-keys --with-colons gh-keyring.gpg | awk -F: '/^(pub|fpr|uid)/ {print $1, $2, $10}' | cut -c1-160
echo "(Fingerabdruck, den das offizielle Feature fest eintraegt: 7F38BBB59D064DBCB3D84D725612B36462313325)"
curl -fsSL https://cli.github.com/packages/dists/stable/InRelease -o gh-InRelease
gpg --no-default-keyring --keyring "$t/gh-keyring.gpg" --status-fd 1 --verify gh-InRelease 2>&1 | grep -E 'GOODSIG|VALIDSIG|BADSIG|NO_PUBKEY|ERRSIG' | cut -c1-140
echo "-- Version in diesem apt-Repository (Paket gh, amd64)"
curl -fsSL https://cli.github.com/packages/dists/stable/main/binary-amd64/Packages -o gh-Packages
awk '/^Package: gh$/ {p=1} p && /^Version:/ {print; exit}' gh-Packages
aussage "gh-apt-fingerprint-gleich-feature=$(gpg --show-keys --with-colons gh-keyring.gpg | awk -F: '/^fpr/ {print $10}' | grep -c 7F38BBB59D064DBCB3D84D725612B36462313325)"

echo
echo "################ glab (gitlab.com/gitlab-org/cli)"
curl -sS 'https://gitlab.com/api/v4/projects/gitlab-org%2Fcli/releases?per_page=1' > glab-release.json
echo "-- neuestes Release (Tag, Zeitpunkt, Zahl der Asset-Links, Quellen)"
jq -c '.[0]|{tag_name, released_at, links: (.assets.links|length), sources: (.assets.sources|length)}' glab-release.json
tag=$(jq -r '.[0].tag_name' glab-release.json); ver=${tag#v}
echo "-- Asset-Links fuer Linux amd64/arm64 und Pruefsummen (Name, Typ)"
jq -r '.[0].assets.links[]|select(.name|test("linux_(amd64|arm64)\\.(tar\\.gz|deb)$|checksums"))|"\(.name) \(.link_type)"' glab-release.json
echo "-- Assets, deren Name nach Signatur, Provenance oder SBOM aussieht (Anzahl)"
jq -r '[.[0].assets.links[].name|select(test("\\.(asc|sig|pem|sigstore|intoto\\.jsonl|sbom|spdx|cdx)|attest|cosign|provenance|sbom|bundle";"i"))]|length' glab-release.json
a="glab_${ver}_linux_amd64.tar.gz"
url=$(jq -r --arg n "$a" '.[0].assets.links[]|select(.name==$n)|.direct_asset_url' glab-release.json)
curl -fsSL "$url" -o "$a"
curl -fsSL "$(jq -r '.[0].assets.links[]|select(.name=="checksums.txt")|.direct_asset_url' glab-release.json)" -o glab-checksums.txt
echo "-- Archiv $a: sha256 der Datei gegen checksums.txt"
d=$(sha256sum "$a" | cut -d' ' -f1)
echo "datei:    $d"
echo "checksums: $(awk -v f="$a" '$2 == f {print $1}' glab-checksums.txt)"
aussage "glab-tag=$tag"
aussage "glab-amd64-sha256=$d"
aussage "glab-checksums-gleich-datei=$([ "$(awk -v f="$a" '$2 == f {print $1}' glab-checksums.txt)" = "$d" ] && echo ja || echo nein)"
echo "-- checksums.txt aus dem Paketregister (Adresse, die das Feature von SrzStephen nutzt) gegen die Datei des Release"
curl -fsSL "https://gitlab.com/api/v4/projects/gitlab-org%2Fcli/packages/generic/glab/${ver}/checksums.txt" -o glab-checksums-pkg.txt && { cmp glab-checksums.txt glab-checksums-pkg.txt && echo "gleich"; } || echo "abweichend oder nicht abrufbar"
echo "-- Release-Konfiguration des Projekts: Zeilen zu Signatur, Attestierung, Pruefsumme"
for f in .goreleaser.yml .gitlab-ci.yml; do
  curl -sS "https://gitlab.com/gitlab-org/cli/-/raw/main/$f" -o "glab-$f"
  echo "$f groesse=$(wc -c < "glab-$f")"
  grep -n -i -E 'cosign|sigstore|gpg|sign|provenance|slsa|attest|checksum|sbom' "glab-$f" | cut -c1-180 | head -20
done
echo "-- Include-Dateien der CI-Konfiguration (Namen) mit Signatur-Bezug"
grep -n -E 'include|local:|project:' glab-.gitlab-ci.yml | cut -c1-160 | head -20

echo
echo "################ Rhythmus der Releases (die acht neuesten, fuer die Pflege der Pins)"
echo "-- gh"
curl -sS 'https://api.github.com/repos/cli/cli/releases?per_page=8' | jq -r '.[]|"\(.tag_name) \(.published_at) prerelease=\(.prerelease) immutable=\(.immutable)"'
echo "-- glab"
curl -sS 'https://gitlab.com/api/v4/projects/gitlab-org%2Fcli/releases?per_page=8' | jq -r '.[]|"\(.tag_name) \(.released_at)"'
