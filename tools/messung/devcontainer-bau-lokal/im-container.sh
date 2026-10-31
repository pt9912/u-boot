set -u
cd /workspaces/u-boot
echo "== Benutzer: $(whoami) uid=$(id -u)"
sudo sh /usr/local/bin/u-boot-docker-init >/tmp/docker-init.log 2>&1 || echo "docker-init Exit $?"
for i in $(seq 1 30); do docker info >/dev/null 2>&1 && break; sleep 1; done
echo "== Docker: $(docker info --format 'server {{.ServerVersion}}, storage {{.Driver}}' 2>&1 | head -1)"
echo "== Werkzeuge"
go version
golangci-lint --version | head -1
echo "BUILD_MODE=${BUILD_MODE:-<ungesetzt>} GOTOOLCHAIN=${GOTOOLCHAIN:-<ungesetzt>}"
echo "GOCACHE=${GOCACHE:-<ungesetzt>} GOMODCACHE=${GOMODCACHE:-<ungesetzt>} GOFLAGS=${GOFLAGS:-<ungesetzt>}"
git config --global --add safe.directory /workspaces/u-boot
echo "== make modus-pruefung"
make modus-pruefung 2>&1 | tail -2
echo "== make -n test (lokal)"
make -n test | head -2
echo "== Teillauf: make test PKG=./internal/hexagon/domain/..."
t0=$(date +%s.%N); make test PKG=./internal/hexagon/domain/... 2>&1 | tail -2; t1=$(date +%s.%N)
echo "wall=$(awk -v a="$t1" -v b="$t0" 'BEGIN{printf "%.1f", a-b}')s"
echo "== make gates PKG=x muss abgelehnt werden"
make gates PKG=./x 2>&1 | tail -1
echo "== make gates (lokaler Modus), Lauf 1 (kalt)"
t0=$(date +%s.%N); make gates > /tmp/gates1.log 2>&1; rc=$?; t1=$(date +%s.%N)
echo "gates exit=$rc wall=$(awk -v a="$t1" -v b="$t0" 'BEGIN{printf "%.1f", a-b}')s"
grep -E 'Befund\(e\)|coverage-gate: OK|FAIL|Error' /tmp/gates1.log | tail -4
echo "== make gates, Lauf 2 (warm)"
t0=$(date +%s.%N); make gates > /tmp/gates2.log 2>&1; rc=$?; t1=$(date +%s.%N)
echo "gates exit=$rc wall=$(awk -v a="$t1" -v b="$t0" 'BEGIN{printf "%.1f", a-b}')s"
grep -E 'Befund\(e\)|coverage-gate: OK|FAIL|Error' /tmp/gates2.log | tail -4
echo "== Einzelzeiten, Lauf 3 (warm)"
for t in lint test coverage-gate; do
  t0=$(date +%s.%N); make $t >/dev/null 2>&1; rc=$?; t1=$(date +%s.%N)
  printf '%-14s exit=%s wall=%ss\n' "$t" "$rc" "$(awk -v a="$t1" -v b="$t0" 'BEGIN{printf "%.1f", a-b}')"
done
