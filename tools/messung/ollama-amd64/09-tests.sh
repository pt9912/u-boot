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
docker ps -a -q --no-trunc > "$WORK/container-vor.txt"
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

# Container, die waehrend des Testlaufs entstanden und danach noch da sind. Auf dem geteilten Host koennen fremde Container darunter
# sein; ausgewiesen werden deshalb nur ihre Anzahl und jene, deren Befehl den Ollama-Block traegt (Zwischencontainer des
# fehlgeschlagenen Baus der Negativprobe). Sie werden nicht entfernt (kein Name, kein Label, das sie der Messung zuordnet).
docker ps -a -q --no-trunc > "$WORK/container-nach.txt"
neu=$(grep -v -x -F -f "$WORK/container-vor.txt" "$WORK/container-nach.txt")
echo "-- Container nach dem Testlauf, die vorher nicht da waren: $(echo -n "$neu" | grep -c .)"
zw=0
for c in $neu; do
  if docker inspect "$c" --format '{{.Path}} {{.Args}}' | grep -q 'OLLAMA_SHA256_AMD64'; then
    zw=$((zw + 1))
    docker inspect "$c" --format 'zwischencontainer id={{printf "%.12s" .Id}} status={{.State.Status}} exit={{.State.ExitCode}} name={{.Name}} erstellt={{.Created}}'
  fi
done
echo "zwischencontainer-des-fehlgeschlagenen-baus=$zw"
aussage "zwischencontainer-nach-negativprobe=$zw"
