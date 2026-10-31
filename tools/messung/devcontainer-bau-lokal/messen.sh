#!/usr/bin/env bash
# messen.sh — Messlauf M0 des Slice slice-v1-devcontainer-bau-lokal.
#
# WAS ES MISST (alles in Docker, der Host braucht keine Go-Toolchain):
#   A. Test-Cache von Go: Laeuft ein zweiter Lauf von `go test ./...` aus dem
#      Cache, auch mit -coverprofile und -coverpkg (so ruft das Dockerfile die
#      Stage `coverage`)? Zaehlung der Zeilen `(cached)` je Lauf.
#   B. Eine Aenderung an EINER Testdatei: Wie viele Pakete laufen danach neu?
#   C. Das Lint-Image golangci/golangci-lint:$GOLANGCI_LINT_VERSION: Go-Version
#      und ob `bash` vorhanden ist.
#   D. Wanduhr von `make gates` im Docker-Modus (zwei Laeufe).
#
# AUFRUF (aus der Repo-Wurzel):
#   bash tools/messung/devcontainer-bau-lokal/messen.sh > docs/messungen/slice-v1-devcontainer-bau-lokal/roh.txt
#
# SCHREIBT NICHT in das Repo, außer nach stdout; die Quellen werden schreibgeschuetzt
# eingehaengt, die Aenderung an einer Testdatei (B) passiert in einer Kopie unter
# $TMPDIR. Zwei benannte Volumes (uboot-m0-gocache, uboot-m0-mod) tragen die Caches
# und werden am Ende entfernt.
set -euo pipefail

root="$(pwd)"
[[ -f repo.mk && -f Dockerfile ]] || { echo "messen.sh: aus der Repo-Wurzel aufrufen" >&2; exit 2; }

go_version="$(sed -n 's/^GO_VERSION[[:space:]]*?=[[:space:]]*//p' repo.mk | head -1)"
lint_version="$(sed -n 's/^GOLANGCI_LINT_VERSION[[:space:]]*?=[[:space:]]*//p' repo.mk | head -1)"
[[ -n "$go_version" && -n "$lint_version" ]] || { echo "messen.sh: Pins in repo.mk nicht gefunden" >&2; exit 2; }

vol_cache=uboot-m0-gocache
vol_mod=uboot-m0-mod
tmp="$(mktemp -d)"
cleanup() { docker volume rm -f "$vol_cache" "$vol_mod" >/dev/null 2>&1 || true; rm -rf "$tmp"; }
trap cleanup EXIT

echo "# Messlauf M0 slice-v1-devcontainer-bau-lokal"
echo "# Datum: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
echo "# Host: $(uname -sm), Docker $(docker version --format '{{.Server.Version}}')"
echo "# GO_VERSION=$go_version GOLANGCI_LINT_VERSION=$lint_version"
echo

docker volume rm -f "$vol_cache" "$vol_mod" >/dev/null 2>&1 || true

# go_run <quelle> <skript>: fuehrt <skript> (sh) im Go-Container aus, Quelle schreibgeschuetzt.
go_run() {
  local src="$1" script="$2"
  docker run --rm \
    -v "$src":/src:ro -v "$vol_cache":/root/.cache/go-build -v "$vol_mod":/go/pkg/mod \
    -w /src -e GOFLAGS="-mod=readonly -buildvcs=false" -e CGO_ENABLED=0 \
    "golang:$go_version" sh -c "$script"
}

# lauf <name> <quelle> <skript>: misst die Wanduhr und zaehlt die Ergebniszeilen.
lauf() {
  local name="$1" src="$2" script="$3" t0 t1 out
  t0="$(date +%s.%N)"
  out="$(go_run "$src" "$script" 2>&1 || true)"
  t1="$(date +%s.%N)"
  printf 'LAUF %-26s wall=%ss  ok_cached=%s ok_neu=%s ohne_tests=%s FAIL=%s\n' "$name" \
    "$(awk -v a="$t1" -v b="$t0" 'BEGIN{printf "%6.1f", a-b}')" \
    "$(grep -c -E '^ok .*\(cached\)' <<<"$out" || true)" \
    "$(grep -E '^ok ' <<<"$out" | grep -vc '(cached)' || true)" \
    "$(grep -c -E 'no test files' <<<"$out" || true)" \
    "$(grep -c -E '^(FAIL|--- FAIL)' <<<"$out" || true)"
}

coverpkg='$(go list ./internal/... | tr "\n" "," | sed "s/,$//")'

echo "## A. Test-Cache (Quelle unveraendert, Caches in benannten Volumes)"
lauf "test_kalt"            "$root" 'go test ./... 2>&1'
lauf "test_warm"            "$root" 'go test ./... 2>&1'
lauf "cover_kalt"           "$root" "go test -coverpkg=\"$coverpkg\" -coverprofile=/tmp/c.out -covermode=atomic ./... 2>&1"
lauf "cover_warm"           "$root" "go test -coverpkg=\"$coverpkg\" -coverprofile=/tmp/c.out -covermode=atomic ./... 2>&1"
lauf "cover_ohne_pkgliste"  "$root" 'go test -coverprofile=/tmp/c.out ./... 2>&1'
lauf "cover_ohne_pkgliste2" "$root" 'go test -coverprofile=/tmp/c.out ./... 2>&1'
echo

echo "## B. Eine Testdatei geaendert (Kopie der Quellen, eine Zeile angehaengt)"
cp -a "$root/." "$tmp/src" 2>/dev/null || true
rm -rf "$tmp/src/.git"
testfile="$(cd "$tmp/src" && find internal/hexagon/domain -name '*_test.go' | sort | head -1)"
printf '\n// messung M0: Kommentarzeile\n' >> "$tmp/src/$testfile"
echo "# geaenderte Datei: $testfile"
lauf "test_nach_aenderung"  "$tmp/src" 'go test ./... 2>&1'
lauf "cover_nach_aenderung" "$tmp/src" "go test -coverpkg=\"$coverpkg\" -coverprofile=/tmp/c.out -covermode=atomic ./... 2>&1"
echo

echo "## C. Lint-Image golangci/golangci-lint:$lint_version"
docker run --rm --entrypoint sh "golangci/golangci-lint:$lint_version" -c \
  'go version; printf "bash: "; command -v bash || echo FEHLT; printf "sh: "; command -v sh || echo FEHLT; golangci-lint --version | head -1' 2>&1
echo

echo "## D. make gates im Docker-Modus (Wanduhr)"
for i in 1 2; do
  t0="$(date +%s.%N)"
  if make gates >/dev/null 2>&1; then st=ok; else st=ROT; fi
  t1="$(date +%s.%N)"
  printf 'GATES lauf=%s status=%s wall=%ss\n' "$i" "$st" "$(awk -v a="$t1" -v b="$t0" 'BEGIN{printf "%6.1f", a-b}')"
done
