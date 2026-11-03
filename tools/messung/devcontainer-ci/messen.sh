#!/usr/bin/env bash
# messen.sh — Fahrer der Messung M1 zu slice-v1-bau-nur-devcontainer (Plan-Risiko 1,
# Messung vor dem Vertrags-Wortlaut). Misst: Image-Bau (M1), verschachtelter Daemon (M2),
# make gates im Container (M3), make test-docker im Container (M4), CI-Runner-Analogie (M5).
#
# AUFRUF (aus der Repo-Wurzel; Rohdaten je Modus nach
# docs/messungen/slice-v1-bau-nur-devcontainer/ umgeleitet):
#   bash tools/messung/devcontainer-ci/messen.sh bau        > docs/messungen/slice-v1-bau-nur-devcontainer/roh-bau.txt 2>&1
#   bash tools/messung/devcontainer-ci/messen.sh dind       > docs/messungen/slice-v1-bau-nur-devcontainer/roh-dind.txt 2>&1
#   bash tools/messung/devcontainer-ci/messen.sh gates      > docs/messungen/slice-v1-bau-nur-devcontainer/roh-gates.txt 2>&1
#   bash tools/messung/devcontainer-ci/messen.sh test-docker > docs/messungen/slice-v1-bau-nur-devcontainer/roh-test-docker.txt 2>&1
#   bash tools/messung/devcontainer-ci/messen.sh analogie   > docs/messungen/slice-v1-bau-nur-devcontainer/roh-analogie.txt 2>&1
set -uo pipefail

[[ -f .devcontainer/Dockerfile && -f repo.mk ]] || { echo "messen.sh: aus der Repo-Wurzel aufrufen" >&2; exit 2; }
mode="${1:?Modus: bau|dind|gates|test-docker|analogie}"
root="$(pwd)"
tools_rel="tools/messung/devcontainer-ci"
tag=uboot-devc-m1

# Die containerEnv aus .devcontainer/devcontainer.json als Aufrufoptionen nachgebaut
# (die bare `docker run`-Form uebernimmt sie nicht; die CLI tut es beim Start).
env_args=(-e BUILD_MODE=local -e GOTOOLCHAIN=local
  -e GOFLAGS="-mod=readonly -buildvcs=false"
  -e GOCACHE=/workspaces/u-boot/.cache/go-build
  -e GOMODCACHE=/workspaces/u-boot/.cache/go-mod)

case "$mode" in
  bau)
    log="$(mktemp)"
    echo "== M1a Kontext-Probelauf: docker build -f .devcontainer/Dockerfile ."
    t0=$(date +%s)
    docker build -t "$tag-ctx-probe" -f .devcontainer/Dockerfile . >"$log" 2>&1
    rc=$?
    t1=$(date +%s)
    echo "exit=$rc wall=$((t1-t0))s"
    tail -8 "$log"
    if [ "$rc" -eq 0 ]; then docker rmi "$tag-ctx-probe" >/dev/null 2>&1; fi
    echo
    echo "== M1b Image-Bau: docker build -t $tag -f .devcontainer/Dockerfile .devcontainer"
    t0=$(date +%s)
    docker build -t "$tag" -f .devcontainer/Dockerfile .devcontainer >"$log" 2>&1
    rc=$?
    t1=$(date +%s)
    echo "exit=$rc wall=$((t1-t0))s"
    tail -12 "$log"
    rm -f "$log"
    [ "$rc" -eq 0 ] || exit "$rc"
    echo
    echo "== Image-Steckbrief"
    docker image inspect "$tag" --format 'id={{.Id}}'
    docker image inspect "$tag" --format 'size_bytes={{.Size}} created={{.Created}} arch={{.Architecture}}'
    echo "== Werkzeuge im Image"
    docker run --rm "$tag" bash -c 'id -un; echo HOME=$HOME; make --version 2>&1 | head -1; bash --version | head -1; go version; golangci-lint version; docker --version; git --version; ollama --version 2>/dev/null; gh --version | head -1; node --version'
    echo "== Pfade der Werkzeuge"
    docker run --rm "$tag" bash -c 'command -v make bash go golangci-lint docker dockerd ollama gh || true'
    echo "== docker info ohne gestarteten Daemon (Erwartung: Fehler)"
    docker run --rm "$tag" docker info 2>&1 | grep -i 'cannot connect' | head -1
    ;;
  dind)
    docker run --rm --privileged \
      -v uboot-m1-docker:/var/lib/docker \
      -v "$root":/workspaces/u-boot -w /workspaces/u-boot \
      "${env_args[@]}" "$tag" bash "$tools_rel/im-container.sh"
    ;;
  gates)
    # Die Go-Caches sitzen (wie im echten Devcontainer je containerEnv) unter
    # dem Workspace — bei dem Bind-Mount-Form im Baum ./.cache (uid 1000, im
    # .gitignore). Eine kalte Lage herzustellen, ist Sache des Aufrufs
    # (rm -rf .cache vor diesem Modus). Frische NAMED Volumes an diesen
    # Unterpfaden sind root-owned (Mountpoint vom Daemon angelegt, ohne
    # Vorbild-Verzeichnis im Image) und brechen den Go-Cache; die gemessene
    # Form ist diese.
    for lauf in 1 2; do
      echo "== M3 Lauf $lauf"
      docker run --rm --privileged \
        -v uboot-m1-docker:/var/lib/docker \
        -v "$root":/workspaces/u-boot -w /workspaces/u-boot \
        "${env_args[@]}" "$tag" bash "$tools_rel/im-gates.sh" "$lauf" || echo "container-exit=$?"
    done
    ;;
  test-docker)
    docker run --rm --privileged \
      -v uboot-m1-docker:/var/lib/docker \
      -v "$root":/workspaces/u-boot -w /workspaces/u-boot \
      "${env_args[@]}" "$tag" bash "$tools_rel/im-test-docker.sh" || echo "container-exit=$?"
    ;;
  analogie)
    echo "== M5: Container A (unprivilegiert, nur den Host-Docker-Socket eingehaengt; die Analogie zum CI-Job-Container)"
    docker run --rm \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -v "$root/$tools_rel":/m:ro \
      "$tag" bash /m/analogie-ebene-a.sh
    ;;
  *)
    echo "unbekannter Modus: $mode" >&2; exit 2 ;;
esac