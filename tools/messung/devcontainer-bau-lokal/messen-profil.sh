#!/usr/bin/env bash
# messen-profil.sh — Ergänzung E zum Messlauf M0 des Slice slice-v1-devcontainer-bau-lokal.
#
# FRAGE. Ein gecachter Lauf von `go test -coverprofile` meldet `(cached)` je Paket
# (siehe messen.sh, Teil A). Liefert er trotzdem das vollständige Profil, und ist
# die Summe von `go tool cover -func` dieselbe wie bei einem kalten Lauf? Davon
# hängt ab, ob das Coverage-Gate im lokalen Modus aus dem Cache laufen darf.
#
# AUFRUF (aus der Repo-Wurzel):
#   bash tools/messung/devcontainer-bau-lokal/messen-profil.sh >> docs/messungen/slice-v1-devcontainer-bau-lokal/roh.txt
set -euo pipefail

[[ -f repo.mk && -f Dockerfile ]] || { echo "messen-profil.sh: aus der Repo-Wurzel aufrufen" >&2; exit 2; }
go_version="$(sed -n 's/^GO_VERSION[[:space:]]*?=[[:space:]]*//p' repo.mk | head -1)"
root="$(pwd)"
vol_cache=uboot-m0e-gocache
vol_mod=uboot-m0e-mod
cleanup() { docker volume rm -f "$vol_cache" "$vol_mod" >/dev/null 2>&1 || true; }
trap cleanup EXIT
docker volume rm -f "$vol_cache" "$vol_mod" >/dev/null 2>&1 || true

inner='
COVERPKG=$(go list ./internal/... | tr "\n" "," | sed "s/,$//")
go test -coverpkg="$COVERPKG" -coverprofile=/tmp/c.out -covermode=atomic ./... >/tmp/test.log 2>&1 || echo "go test: Exit ungleich 0"
printf "gecachte Pakete: %s von %s\n" "$(grep -c "(cached)" /tmp/test.log || true)" "$(grep -c "^ok " /tmp/test.log || true)"
printf "Profil: %s Zeilen, %s Bytes\n" "$(wc -l < /tmp/c.out)" "$(wc -c < /tmp/c.out)"
printf "Profil sha256: %s\n" "$(sha256sum /tmp/c.out | cut -d" " -f1)"
go tool cover -func=/tmp/c.out | tail -1
'

echo
echo "## E. Coverage-Profil aus dem Cache (Go $go_version, Quelle unveraendert)"
for lauf in kalt warm; do
  echo "--- Lauf $lauf"
  docker run --rm \
    -v "$root":/src:ro -v "$vol_cache":/root/.cache/go-build -v "$vol_mod":/go/pkg/mod \
    -w /src -e GOFLAGS="-mod=readonly -buildvcs=false" -e CGO_ENABLED=0 \
    "golang:$go_version" sh -c "$inner" 2>&1
done
