# im-gates.sh — M3 zu slice-v1-bau-nur-devcontainer: make gates im Container (BUILD_MODE=local).
# Argument: Lauf-Nummer (nur fuer die Beschriftung). Aufruf siehe messen.sh Modus gates.
set -u
lauf="${1:-?}"
cd /workspaces/u-boot
echo "== Benutzer: $(whoami) uid=$(id -u)"
sudo sh /usr/local/bin/u-boot-docker-init >/tmp/docker-init.log 2>&1
echo "docker-init-exit=$?"
for i in $(seq 1 120); do docker info >/dev/null 2>&1 && break; sleep 1; done
docker info --format 'daemon server={{.ServerVersion}} storage={{.Driver}}' | head -1
echo "BUILD_MODE=${BUILD_MODE:-<ungesetzt>} GOTOOLCHAIN=${GOTOOLCHAIN:-<ungesetzt>} GOFLAGS=${GOFLAGS:-<ungesetzt>}"
echo "GOCACHE=${GOCACHE:-<ungesetzt>} GOMODCACHE=${GOMODCACHE:-<ungesetzt>}"
git config --global --add safe.directory /workspaces/u-boot
echo "== make gates, Lauf $lauf"
t0=$(date +%s); make gates >"/tmp/gates-$lauf.log" 2>&1; rc=$?; t1=$(date +%s)
echo "gates-exit=$rc wall=$(awk -v a="$t1" -v b="$t0" 'BEGIN{printf "%.1f", a-b}')s"
echo "== Ergebniszeilen des Laufes $lauf"
grep -E 'Befund|OK|FAIL|Error' "/tmp/gates-$lauf.log" | tail -10
echo "== Ende des Protokolls des Laufes $lauf"
tail -6 "/tmp/gates-$lauf.log"