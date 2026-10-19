#!/bin/bash
# Messstand (Host, Docker, Kernel, Datum) und Bestand der Docker-Objekte VOR der Messung.
# Aufruf: WORK=<messordner> OUT=<ausgabeverzeichnis> tools/messung/egress-dns-upstream/01-messstand.sh
#   -> $OUT/rohdaten-01-messstand.txt, rohdaten-11-vorher.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-01-messstand.txt
echo "== Messstand, $(now)"
echo "host: $(uname -srm)"
echo "docker-client: $(docker version --format '{{.Client.Version}}')"
echo "docker-server: $(docker version --format '{{.Server.Version}} ({{.Server.Os}}/{{.Server.Arch}})')"
docker info --format 'engine-os: {{.OperatingSystem}}; cgroup {{.CgroupVersion}}; cpus {{.NCPU}}; storage {{.Driver}}'
echo "kernel.apparmor_restrict_unprivileged_userns=$(cat /proc/sys/kernel/apparmor_restrict_unprivileged_userns 2>&1)"
echo "-- Resolver des Hosts (Anzahl der nameserver-Zeilen; Adressen siehe Upstream-Liste im Container)"
echo "nameserver-zeilen=$(grep -c '^nameserver' /etc/resolv.conf)"
echo "-- Bestand vorher"
bestand | tee "$OUT/rohdaten-11-vorher.txt"
