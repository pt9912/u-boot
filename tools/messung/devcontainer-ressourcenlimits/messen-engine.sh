#!/usr/bin/env bash
# messen-engine.sh — Messung 1, 2, 4 und 6 des Slice slice-v1-devcontainer-ressourcenlimits
# an der Engine (docker run mit --memory und --cpus), ohne u-boot.
#
# WAS ES MISST:
#   6.  Werte, die die Engine für --memory und --cpus annimmt oder abweist.
#   1.  Greifen die Grenzen im Container (cgroup-Dateien), ohne und mit --privileged.
#   2.  Greifen sie im verschachtelten Docker (docker:dind, --privileged)?
#   4.  Ein Prozess über dem Limit: Exit-Code des Containers, OOMKilled.
# Die Messung belegt das Verhalten der Engine; die Auslieferung durch
# `u-boot devcontainer up` belegt der Docker-Test des Slice.
#
# AUFRUF: bash tools/messung/devcontainer-ressourcenlimits/messen-engine.sh > roh.txt 2>&1
set -u
IMG=alpine:3.20
DIND=docker:27-dind
echo "# Datum: $(date -u +%Y-%m-%dT%H:%M:%SZ)"; echo "# Host: $(uname -sm)"; docker version --format '# Engine: {{.Server.Version}}'
docker pull -q "$IMG" >/dev/null; docker pull -q "$DIND" >/dev/null

echo "=== 6 Grammatik --memory"
for v in 3g 512m 512M 4G 64k 1024b 1048576 0 -1g 1.5g 4t 3gb "3 g" 6m 5m 4m 1g1; do
  out=$(docker run --rm --memory="$v" "$IMG" true 2>&1); echo "memory=[$v] exit=$? ${out:+msg: $out}"
done
echo "=== 6 Grammatik --cpus"
for v in 2 1.5 0.5 0.01 0.001 0 -1 .5 2. 1e1 2x " 2" 1,5 9999; do
  out=$(docker run --rm --cpus="$v" "$IMG" true 2>&1); echo "cpus=[$v] exit=$? ${out:+msg: $out}"
done

for run in 1 2; do
 for priv in "" "--privileged"; do
  echo "=== 1 Lauf $run priv=[$priv]"
  docker run --rm $priv --memory=512m --cpus=1.5 "$IMG" sh -c 'echo memory.max=$(cat /sys/fs/cgroup/memory.max); echo cpu.max=$(cat /sys/fs/cgroup/cpu.max)'
 done
done

for run in 1 2; do
  echo "=== 2 Lauf $run verschachtelter Docker (dind, --privileged, --memory=1g --cpus=1)"
  c=$(docker run -d --privileged --memory=1g --cpus=1 -e DOCKER_TLS_CERTDIR= "$DIND")
  for i in $(seq 1 30); do docker exec "$c" docker info >/dev/null 2>&1 && break; sleep 1; done
  docker exec "$c" sh -c 'echo aussen memory.max=$(cat /sys/fs/cgroup/memory.max) cpu.max=$(cat /sys/fs/cgroup/cpu.max)'
  docker exec "$c" docker pull -q "$IMG" >/dev/null 2>&1 || echo "pull im inneren fehlgeschlagen"
  echo "innen ohne eigene Grenze:"
  docker exec "$c" docker run --rm "$IMG" sh -c 'echo memory.max=$(cat /sys/fs/cgroup/memory.max); echo cpu.max=$(cat /sys/fs/cgroup/cpu.max)'
  echo "innen mit --memory=2g (größer als aussen):"
  docker exec "$c" docker run --rm --memory=2g "$IMG" sh -c 'echo memory.max=$(cat /sys/fs/cgroup/memory.max)'
  docker rm -f "$c" >/dev/null
done

for run in 1 2; do
  echo "=== 4 Lauf $run Prozess über dem Limit (--memory=32m)"
  c=$(docker run -d --memory=32m --memory-swap=32m "$IMG" sh -c 'head -c 200000000 /dev/zero | tail -c 1 >/dev/null; sleep 1')
  docker wait "$c"; docker inspect -f 'exit={{.State.ExitCode}} oomkilled={{.State.OOMKilled}}' "$c"; docker rm -f "$c" >/dev/null
done
