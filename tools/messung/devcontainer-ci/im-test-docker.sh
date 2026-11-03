# im-test-docker.sh — M4 zu slice-v1-bau-nur-devcontainer: make test-docker im Container
# (der e2e-Zug mit dem verschachtelten Daemon). Aufruf siehe messen.sh Modus test-docker.
set -u
cd /workspaces/u-boot
echo "== Benutzer: $(whoami) uid=$(id -u)"
sudo sh /usr/local/bin/u-boot-docker-init >/tmp/docker-init.log 2>&1
echo "docker-init-exit=$?"
for i in $(seq 1 60); do docker info >/dev/null 2>&1 && break; sleep 1; done
echo "== Docker: $(docker info --format 'server {{.ServerVersion}}, storage {{.Driver}}' 2>&1 | head -1)"
git config --global --add safe.directory /workspaces/u-boot
echo "BUILD_MODE=${BUILD_MODE:-<ungesetzt>} GOCACHE=${GOCACHE:-<ungesetzt>}"
echo "== make test-docker (bis zu 30 Minuten)"
t0=$(date +%s); timeout 1800 make test-docker >/tmp/test-docker.log 2>&1; rc=$?; t1=$(date +%s)
echo "test-docker-exit=$rc wall=$(awk -v a="$t1" -v b="$t0" 'BEGIN{printf "%.1f", a-b}')s"
echo "== Zählung (Testfunktionen der obersten Ebene)"
printf 'PASS=%s FAIL=%s SKIP=%s\n' \
  "$(grep -c -E '^--- PASS' /tmp/test-docker.log)" \
  "$(grep -c -E '^--- FAIL' /tmp/test-docker.log)" \
  "$(grep -c -E '^--- SKIP' /tmp/test-docker.log)"
echo "== fehlgeschlagene Tests (Namen)"
grep -E '^--- FAIL' /tmp/test-docker.log | cut -c1-160
echo "== übersprungene Tests (Namen)"
grep -E '^--- SKIP' /tmp/test-docker.log | cut -c1-160
echo "== Paketergebnisse (alles außer ok)"
grep -E '^(FAIL|---|\?)\s' /tmp/test-docker.log | grep -v -E '^--- (PASS|FAIL|SKIP)' | head -12 | cut -c1-160
echo "== Kontext je Fehlschlag (je 14 Zeilen nach der FAIL-Marker, gekürzt)"
grep -A 14 -E '^--- FAIL' /tmp/test-docker.log | cut -c1-200 | head -220
echo "== letzte 15 Zeilen"
tail -15 /tmp/test-docker.log | cut -c1-200