# im-container.sh — M2 zu slice-v1-bau-nur-devcontainer: verschachtelter Daemon im Container.
# Läuft als Benutzer vscode im privilegierten Devcontainer-Container (Aufruf siehe messen.sh Modus dind).
set -u
cd /workspaces/u-boot
echo "== Benutzer: $(whoami) uid=$(id -u)"
sudo sh /usr/local/bin/u-boot-docker-init >/tmp/docker-init.log 2>&1
echo "docker-init-exit=$?"
for i in $(seq 1 60); do docker info >/dev/null 2>&1 && break; sleep 1; done
echo "== docker version im laufenden Container (Client und SERVER = der verschachtelte Daemon)"
docker version
echo "== docker info (Kurzform)"
docker info --format 'server={{.ServerVersion}} storage={{.Driver}} cgroup={{.CgroupDriver}} cgroupv{{.CgroupVersion}}' 2>&1 | head -1
echo
echo "== Test-Container im Container: docker run --rm busybox true (mit Pull)"
t0=$(date +%s); docker pull busybox >/tmp/busybox-pull.log 2>&1; pull_rc=$?; t1=$(date +%s)
echo "pull-exit=$pull_rc wall=$((t1-t0))s"
t0=$(date +%s); docker run --rm busybox true; rc=$?; t1=$(date +%s)
echo "run-exit=$rc wall=$((t1-t0))s"
echo
echo "== docker ps -a (Bestand im verschachtelten Daemon)"
docker ps -a
echo
echo "== make modus-pruefung im Container"
t0=$(date +%s); make modus-pruefung; rc=$?; t1=$(date +%s)
echo "modus-pruefung-exit=$rc wall=$((t1-t0))s"