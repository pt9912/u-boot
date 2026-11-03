# analogie-ebene2-nop.sh — M5, Ebene 2 ohne --privileged: der Fall, der es zeigen muss.
# Läuft im Devcontainer-Container ohne --privileged. Der Startscript fällt nach
# Policy (onUnavailable: warn) auf fallback zurück und der Server-Anteil kommt nicht auf.
set -u
echo "== Ebene 2 (ohne --privileged): $(id -un) uid=$(id -u)"
sudo sh /usr/local/bin/u-boot-docker-init >/tmp/docker-init.log 2>&1
echo "docker-init-exit=$?"
echo "== Status-Datei (u-boot doctor liest diese): cat /var/lib/u-boot/status/daemon"
cat /var/lib/u-boot/status/daemon 2>/dev/null || echo "(keine Status-Datei)"
echo "== docker info (Erwartung: kein Daemon)"
docker info 2>&1 | grep -i -m1 'cannot connect' || docker info >/dev/null 2>&1 && echo "(info gelang unerwartet)"
echo "== docker run (Erwartung: Fehler)"
docker run --rm busybox true 2>&1 | tail -2