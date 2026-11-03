# analogie-ebene2.sh — M5, Ebene 2 mit --privileged: startet der verschachtelte Daemon?
# Läuft im Devcontainer-Container (per socket-Container A auf Ebene 2 erzeugt — der
# layer-1-Daemon honorert --privileged, nicht die Privilegien von A). Aufgerufen aus
# analogie-ebene-a.sh (a).
set -u
echo "== Ebene 2: $(id -un) uid=$(id -u)"
sudo sh /usr/local/bin/u-boot-docker-init >/tmp/docker-init.log 2>&1
echo "docker-init-exit=$?"
for i in $(seq 1 60); do docker info >/dev/null 2>&1 && break; sleep 1; done
echo "== Ebene-2-Daemon: $(docker info --format 'server={{.ServerVersion}} storage={{.Driver}}' 2>&1 | head -1)"
echo "== Test-Container auf Ebene 3: docker run --rm busybox true (mit Pull)"
t0=$(date +%s); docker pull busybox >/tmp/pull.log 2>&1; pull_rc=$?; t1=$(date +%s)
echo "pull-exit=$pull_rc wall=$((t1-t0))s"
t0=$(date +%s); docker run --rm busybox true; rc=$?; t1=$(date +%s)
echo "ebene3-run-exit=$rc wall=$((t1-t0))s"