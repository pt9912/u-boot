#!/bin/bash
# A6: Kombination des Features ollama mit nestedRuntime podman und mit nestedRuntime docker.
# Je Laufzeit: Bau mit und ohne Ollama-Block, Paketliste davor/danach, Versionen, `mounts`-Reihenfolge, Lauf eines kleinen Images
# im Container (ohne Egress, dann mit Egress), gleichzeitig `ollama serve` und ein Pull. Die Container tragen die `runArgs` der
# erzeugten devcontainer.json.
# Voraussetzung: 02-projekte.sh. Aufruf: WORK=<messordner> tools/messung/ollama-amd64/07-kombination.sh [podman|docker]
#   -> $OUT/rohdaten-07-podman.txt, rohdaten-07-docker.txt (+ Build-Logs rohdaten-07-<variante>.log)
set -u
. "$(dirname "$0")/lib.sh"

kombination() { # <laufzeit: podman|docker> <variante mit ollama> <variante ohne ollama>
  local rt=$1 v=$2 vn=$3 id=k1 rc s e
  local C=$PFX-k-$rt
  rohdaten "rohdaten-07-$rt.txt"
  ex() { docker exec -u vscode "$C" bash -c "$1"; }
  echo "== A6 nestedRuntime $rt + ollama, $(now)"
  docker rm -f -v "$C" > /dev/null 2>&1
  for x in $(docker volume ls -q | grep "^$(vname "$v")-"); do docker volume rm "$x" > /dev/null; done

  for x in "$v" "$vn"; do
    s=$(date +%s.%N)
    docker build --progress=plain -t "$(vimg "$x")" "$(vdir "$x")/.devcontainer" > "$WORK/build-$x.log" 2>&1; rc=$?
    e=$(date +%s.%N)
    echo "build $x exit=$rc dauer=$(sekunden "$s" "$e")s size=$(docker image inspect "$(vimg "$x")" --format '{{.Size}}' 2>/dev/null)"
    redact < "$WORK/build-$x.log" > "$OUT/rohdaten-07-$x.log"
    aussage "bau-$x-exit=$rc"
    [ "$rc" = 0 ] || { tail -15 "$WORK/build-$x.log"; return 1; }
  done
  echo "-- Paketliste ohne Ollama-Block ($vn) gegen mit Ollama-Block ($v) (dpkg-query; leer: der Block entfernt und hinterlaesst nichts)"
  docker run --rm "$(vimg "$vn")" dpkg-query -W -f '${Package} ${Version}\n' | sort > "$WORK/dpkg-$vn.txt"
  docker run --rm "$(vimg "$v")" dpkg-query -W -f '${Package} ${Version}\n' | sort > "$WORK/dpkg-$v.txt"
  echo "Pakete $vn=$(grep -c . "$WORK/dpkg-$vn.txt") $v=$(grep -c . "$WORK/dpkg-$v.txt")"
  diff "$WORK/dpkg-$vn.txt" "$WORK/dpkg-$v.txt" && echo "(kein Unterschied)"
  aussage "paketdiff-$rt=$(diff "$WORK/dpkg-$vn.txt" "$WORK/dpkg-$v.txt" | grep -c '^[<>]')"
  echo "-- Pakete der Runtime im Image (Auszug)"
  grep -E '^(uidmap|ca-certificates|docker-ce|docker-ce-cli|containerd.io|iptables|libzstd1|zstd|fuse|dnsmasq-base|nftables) ' "$WORK/dpkg-$v.txt"
  echo "-- Versionen im Image"
  if [ "$rt" = podman ]; then
    docker run --rm "$(vimg "$v")" bash -c 'podman --version; crun --version | head -1; ls -l /usr/local/bin/docker | cut -c1-80'
  else
    docker run --rm "$(vimg "$v")" bash -c 'docker --version; dockerd --version; containerd --version | cut -c1-80'
  fi
  aussage "ollama-version-$rt=$(docker run --rm "$(vimg "$v")" ollama --version 2>&1 | grep -o 'version is .*')"

  echo "-- mounts (Reihenfolge) und runArgs der erzeugten devcontainer.json"
  mapfile -t MOUNTS < <(tool sh -c "grep -v '^[[:space:]]*//' p-$v/.devcontainer/devcontainer.json | jq -r '.mounts[]'" | tr -d '\r')
  mapfile -t RA < <(tool sh -c "grep -v '^[[:space:]]*//' p-$v/.devcontainer/devcontainer.json | jq -r '.runArgs[]'" | tr -d '\r')
  local n=0; for m in "${MOUNTS[@]}"; do n=$((n+1)); echo "  mount $n: $(echo "$m" | grep -o 'target=[^,]*')"; done
  echo "  runArgs: ${RA[*]}"
  aussage "mount-ziele-$rt=$(for m in "${MOUNTS[@]}"; do echo "$m" | grep -o 'target=[^,]*' | tr '\n' ' '; done)"
  ARGS=(); for m in "${MOUNTS[@]}"; do ARGS+=(--mount "${m//\$\{devcontainerId\}/$id}"); done
  docker run -d --name "$C" "${RA[@]}" "${ARGS[@]}" "$(vimg "$v")" sleep 3600 > /dev/null; rc=$?
  echo "start exit=$rc"

  if [ "$rt" = podman ]; then
    mark "postCreateCommand (sandbox-init), ohne Egress"
    ex 'sh /usr/local/bin/u-boot-sandbox-init; echo exit=$?; cat /var/lib/u-boot/status/sandbox 2>/dev/null'
    mark "ohne Egress: podman run eines kleinen Images (Pull von Docker Hub)"
    ex 'timeout 180 podman run --rm docker.io/library/alpine:3 sh -c "echo NESTED_RUN_OK; id -u"; echo exit=$?'
    aussage "nested-run-ohne-egress-$rt=$(ex 'timeout 120 podman run --rm --pull=never docker.io/library/alpine:3 echo ok >/dev/null 2>&1; echo $?')"
  else
    mark "ohne Egress: Daemon im Container starten (docker-init), docker run eines kleinen Images"
    ex 'sudo sh /usr/local/bin/u-boot-docker-init; echo docker-init-exit=$?; cat /var/lib/u-boot/status/docker 2>/dev/null'
    ex 'for i in $(seq 1 60); do docker info >/dev/null 2>&1 && break; sleep 1; done; docker version --format "client={{.Client.Version}} server={{.Server.Version}} os/arch={{.Server.Os}}/{{.Server.Arch}}"'
    ex 'timeout 180 docker run --rm alpine:3 sh -c "echo NESTED_RUN_OK; id -u"; echo exit=$?'
    aussage "nested-run-ohne-egress-$rt=$(ex 'timeout 120 docker run --rm --pull=never alpine:3 echo ok >/dev/null 2>&1; echo $?')"
  fi

  mark "postStartCommand: Egress-Restriktion"
  ex 'sudo sh /usr/local/bin/u-boot-egress-init; echo egress-init-exit=$?; cat /var/lib/u-boot/status/egress'
  if [ "$rt" = podman ]; then run='podman'; img=docker.io/library/alpine:3; img2=docker.io/library/busybox:1.36; else run='docker'; img=alpine:3; img2=busybox:1.36; fi
  echo "-- mit Egress: Lauf des bereits geladenen Images (--pull=never)"
  ex "timeout 120 $run run --rm --pull=never $img sh -c 'echo NESTED_RUN_OK'; echo exit=\$?"
  aussage "nested-run-mit-egress-$rt=$(ex "timeout 120 $run run --rm --pull=never $img echo ok >/dev/null 2>&1; echo \$?")"
  echo "-- mit Egress: Pull eines weiteren Images ($img2; Hosts der Laufzeit stehen in der Liste)"
  ex "timeout 180 $run pull -q $img2; echo exit=\$?"
  aussage "nested-pull-mit-egress-$rt=$(ex "timeout 180 $run image inspect $img2 >/dev/null 2>&1; echo \$?")"
  if ! ex "$run image inspect $img2 > /dev/null 2>&1"; then
    echo "-- Pull ist gescheitert; erneuter Versuch nach 15 s (ein Prozess, der /etc/resolv.conf vor dem Start der Restriktion gelesen hat, kann den alten Resolver weiter nutzen; hier nur gemessen, nicht diagnostiziert)"
    sleep 15
    ex "timeout 180 $run pull -q $img2; echo exit=\$?"
    aussage "nested-pull-mit-egress-nach-15s-$rt=$(ex "timeout 180 $run image inspect $img2 >/dev/null 2>&1; echo \$?")"
  fi
  echo "-- mit Egress: Zugriff aus dem verschachtelten Container auf einen gesperrten Host (example.com; Exit 0 = erreicht)"
  ex "timeout 60 $run run --rm --pull=never $img wget -T 5 -qO /dev/null http://example.com; echo wget-exit=\$?"
  aussage "nested-gesperrt-example.com-$rt=$(ex "timeout 60 $run run --rm --pull=never $img wget -T 5 -qO /dev/null http://example.com >/dev/null 2>&1; echo \$?")"
  echo "-- mit Egress: Zugriff aus dem verschachtelten Container auf eine gesperrte Adresse ohne Namen (1.1.1.1, Port 80; Exit 0 = erreicht)"
  ex "timeout 60 $run run --rm --pull=never $img wget -T 5 -qO /dev/null http://1.1.1.1/; echo wget-exit=\$?"
  aussage "nested-gesperrt-1.1.1.1-$rt=$(ex "timeout 60 $run run --rm --pull=never $img wget -T 5 -qO /dev/null http://1.1.1.1/ >/dev/null 2>&1; echo \$?")"
  echo "-- Gegenstueck: derselbe Zugriff aus dem Container selbst (Exit 6 = Name nicht aufgeloest)"
  ex 'curl -sS -m 6 -o /dev/null http://example.com; echo curl-exit=$?' 2>&1

  mark "gleichzeitig: ollama serve + Pull und ein Lauf in der verschachtelten Laufzeit"
  ex "(nohup ollama serve > /tmp/ollama.log 2>&1 &); for i in \$(seq 1 60); do curl -fsS -m 2 http://127.0.0.1:11434/ >/dev/null 2>&1 && break; sleep 0.5; done
    ( ollama pull smollm2:135m > /tmp/pull.txt 2>&1; echo \$? > /tmp/pull.rc ) &
    ( timeout 120 $run run --rm --pull=never $img sh -c 'sleep 3; echo CONCURRENT_NESTED_OK' > /tmp/nested.txt 2>&1; echo \$? > /tmp/nested.rc ) &
    wait
    echo pull-exit=\$(cat /tmp/pull.rc) nested-exit=\$(cat /tmp/nested.rc); cat /tmp/nested.txt; ollama list"
  aussage "gleichzeitig-pull-exit-$rt=$(ex 'cat /tmp/pull.rc')"
  aussage "gleichzeitig-nested-exit-$rt=$(ex 'cat /tmp/nested.rc')"
  aussage "ollama-library-cpu-$rt=$(ex 'grep -c library=cpu /tmp/ollama.log')"
  echo "-- Volumes der Laufzeit"
  docker volume ls -q | grep "^$(vname "$v")-"
  docker rm -f -v "$C" > /dev/null
  for x in $(docker volume ls -q | grep "^$(vname "$v")-"); do docker volume rm "$x" > /dev/null; done
}

case "${1:-alle}" in
  podman) (kombination podman pod podn) ;;
  docker) (kombination docker dok dokn) ;;
  *) (kombination podman pod podn); (kombination docker dok dokn) ;;
esac
