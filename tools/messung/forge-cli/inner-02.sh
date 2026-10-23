#!/bin/bash
# Laeuft IM Werkzeug-Container (Basis-Image); aufgerufen von 02-oci-quellen.sh.
# Aufruf im Container: inner-02.sh <referenz-ohne-schema> <verzeichnisname>
# Liest anonym aus ghcr.io: Tags, Manifest (Digest, Annotationen), Referrers, das Feature-Archiv (Metadaten, install.sh).
set -u
ref=$1; name=$2
path=${ref#ghcr.io/}
d=/w/oci-$name; rm -rf "$d"; mkdir -p "$d"
tok=$(curl -sS -m 20 "https://ghcr.io/token?scope=repository:$path:pull" | jq -r .token)
get() { curl -sS -m 30 -H "Authorization: Bearer $tok" "$@"; }
acc='application/vnd.oci.image.manifest.v1+json'
echo "== $ref"
echo "-- Tags"
get "https://ghcr.io/v2/$path/tags/list" | jq -c '.tags|sort_by(.)'
echo "-- Manifest :latest (Digest laut Antwort-Header, Layer, Annotationen)"
get -D "$d/h.txt" -H "Accept: $acc" "https://ghcr.io/v2/$path/manifests/latest" -o "$d/manifest.json"
grep -i '^docker-content-digest' "$d/h.txt" | tr -d '\r'
jq -c '{layers: [.layers[]|{mediaType,digest,size}], annotationKeys: (.annotations|keys)}' "$d/manifest.json"
digest=$(grep -i '^docker-content-digest' "$d/h.txt" | tr -d '\r' | awk '{print $2}')
echo "-- Gegenprobe: Digest des Manifests = sha256 der Antwort"
echo "sha256:$(sha256sum "$d/manifest.json" | cut -d' ' -f1)"
echo "-- Referrers (Signaturen, Attestierungen) zum Digest"
code=$(get -o "$d/ref.json" -w '%{http_code}' -H 'Accept: application/vnd.oci.image.index.v1+json' "https://ghcr.io/v2/$path/referrers/$digest")
echo "referrers http=$code"
if [ "$code" = 200 ]; then jq -c '{manifests: (.manifests|length)}' "$d/ref.json"; else head -c 200 "$d/ref.json"; echo; fi
echo "-- Tags der Form sha256-*, *.sig, *.att (Cosign-Schema)"
get "https://ghcr.io/v2/$path/tags/list" | jq -r '.tags[]' | grep -E '^sha256-|\.sig$|\.att$' || echo "(keine)"
echo "-- Metadaten (Annotation dev.containers.metadata): Sicherheitsfelder und Abhaengigkeiten"
jq -r '.annotations["dev.containers.metadata"]' "$d/manifest.json" > "$d/meta.json"
jq -c '{id, version, name, mounts, privileged, init, capAdd, securityOpt, entrypoint, dependsOn, installsAfter, containerEnv, remoteEnv, onCreateCommand, postCreateCommand, postStartCommand, postAttachCommand, customizations, options: (.options|keys)}' "$d/meta.json"
echo "-- Layer laden, Pruefsumme gegen den Digest, Inhalt"
ld=$(jq -r '.layers[0].digest' "$d/manifest.json")
get -L "https://ghcr.io/v2/$path/blobs/$ld" -o "$d/layer.tar"
echo "layer-sha256 gleich Digest: $([ "sha256:$(sha256sum "$d/layer.tar" | cut -d' ' -f1)" = "$ld" ] && echo ja || echo NEIN) ($(wc -c < "$d/layer.tar") Bytes)"
mkdir -p "$d/x"; tar -xf "$d/layer.tar" -C "$d/x"
(cd "$d/x" && find . -type f | sort | while read -r f; do echo "$(sha256sum "$f" | cut -c1-16) $(wc -c < "$f") $f"; done)
echo "-- install.sh: Zeilen zu Herkunft, Pruefung und Fremdskripten (grep -n)"
inst="$d/x/install.sh"
if [ -f "$inst" ]; then
  echo "zeilen=$(wc -l < "$inst")"
  grep -n -i -E 'curl|wget|https?://|apt-get|gpg|keyring|sha256|checksum|cosign|sigstore|verify|git clone|dpkg|latest|releases' "$inst" | cut -c1-200
else
  echo "(kein install.sh)"; ls "$d/x"
fi
echo "-- weitere Dateien des Features (ausser install.sh und devcontainer-feature.json)"
(cd "$d/x" && find . -type f ! -name install.sh ! -name devcontainer-feature.json | sort) | sed 's/^/  /'
