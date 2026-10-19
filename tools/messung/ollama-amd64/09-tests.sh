#!/bin/bash
# A8: die beiden Docker-Tests des Ollama-Features mit sichtbarem Pass oder Skip (go test -v), im Image test-docker-tools des Repos.
# Derselbe Aufruf wie `make test-docker` (Bau-Zeile aus `make -n`, gleiche Mounts und Optionen), nur mit -v und -run statt ./... .
# Das ist KEIN CI-Beleg im Sinn von ADR-0024 Entscheidung 10: der CI-Aufruf selbst bleibt unveraendert.
# Aufruf: WORK=<messordner> tools/messung/ollama-amd64/09-tests.sh      -> $OUT/rohdaten-09.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-09.txt

echo "== A8 Docker-Tests mit -v, $(now)"
buildline=$(make -s -C "$ROOT" -n test-docker IMAGE="$UB" | grep '^docker build')
echo "bau: ${buildline//$ROOT/<repo>}"
(cd "$ROOT" && eval "$buildline") > "$WORK/test-tools-build.log" 2>&1; echo "bau-exit=$?"
s=$(date +%s.%N)
docker run --rm --network=host \
  -v "$ROOT":/src -w /src \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /tmp:/tmp \
  "$UB:test-docker-tools" \
  go test -tags docker -v -timeout 30m -count=1 -run '^TestE2E_SandboxDevcontainer_Ollama(_PinMismatchBreaksBuild)?$' ./internal/e2e/ 2>&1 \
  | grep -E '^(=== RUN|--- |\s+--- |PASS|FAIL|ok |panic|SKIP)|_test.go:[0-9]+:' | cut -c1-240
rc=${PIPESTATUS[0]}
e=$(date +%s.%N)
echo "go-test-exit=$rc gesamtdauer=$(sekunden "$s" "$e")s"
aussage "go-test-exit=$rc"
