#!/usr/bin/env bash
# Shared helpers of the measurement "katalog-gepinnte-toolchain" (M1-M3).
# Sourced by the three measurement scripts; runs nothing on its own.
#
# The measurement starts the pinned Dev Container CLI the same way
# `u-boot devcontainer up` starts it (internal/adapter/driven/docker/
# devcontainer_launcher.go, runArgs): a short-lived helper container from
# the CLI image, the engine socket mounted, the project folder mounted at
# the identical path, `devcontainer up --workspace-folder <folder>`
# (recreate: plus --remove-existing-container).

set -euo pipefail

readonly MESS_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly REPO_ROOT="$(cd "$MESS_ROOT/../../.." && pwd)"
readonly MESS_DOC_DIR="$REPO_ROOT/docs/messungen/slice-v1-katalog-gepinnte-toolchain-entwurf"

# Namespace of the measurement: every helper container name starts with this
# prefix, every volume with it too, so a leftover is recognisable by name
# (Observation BEO-E2E/docker-test-hinterlaesst-zwischencontainer).
readonly MESS_PREFIX="ub-messung"

# CLI helper image: built from the vendored context of the docker driven
# adapter, so the measurement runs the exact pin the launcher runs
# (@devcontainers/cli 0.89.0, package.json + package-lock.json copied
# verbatim); the tag mirrors helperImageRepo/helperImageTag without the
# context hash, because the context is byte-identical to the embed.
readonly MESS_CLI_IMAGE="$MESS_PREFIX-devcontainer-cli:0.89.0"
readonly MESS_HELPER_CONTEXT="$REPO_ROOT/internal/adapter/driven/docker/devcontainercli"
readonly MESS_BASE_IMAGE="mcr.microsoft.com/devcontainers/base:debian"

# mess_ensure_cli_image builds the CLI helper image once, from the vendored
# context. Needs the npm registry (npm ci --ignore-scripts inside the build)
# and docker.io for the two digest-pinned base images.
mess_ensure_cli_image() {
  if docker image inspect "$MESS_CLI_IMAGE" >/dev/null 2>&1; then
    return 0
  fi
  docker build -t "$MESS_CLI_IMAGE" "$MESS_HELPER_CONTEXT"
}

# mess_engine_version prints the engine and CLI versions for a raw-data header.
# $1 = target file.
mess_versions() {
  docker version --format 'Engine: {{.Server.Version}} ({{.Server.Os}}/{{.Server.Arch}})'
  docker run --rm "$MESS_CLI_IMAGE" devcontainer --version
}

# mess_up runs the CLI `up` like the launcher does and echoes the container ID
# of the result line. $1 = project folder, $2 = output prefix (raw-data file
# without extension), rest = extra CLI arguments.
# Print the container id on stdout; the full CLI output is appended to
# "<prefix>.log". On a failed run the file grows with the error and the call
# returns non-zero with no output — the caller records the failure and
# tears down.
mess_up() {
  local folder="$1" prefix="$2"
  shift 2
  local helper_name status out last_line id
  helper_name="$MESS_PREFIX-cli-$(date +%s%N | tail -c 8)"
  out="$(docker run --rm --name "$helper_name" \
    --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock \
    --mount type=bind,source="$folder",target="$folder" \
    "$MESS_CLI_IMAGE" \
    devcontainer up --workspace-folder "$folder" "$@" 2>&1 | tee "$prefix.log")" \
    && status=0 || status=$?
  last_line="$(printf '%s\n' "$out" | sed '/^[[:space:]]*$/d' | tail -n 1 || true)"
  id="$(printf '%s' "$last_line" | sed -n 's/.*"containerId":"\([^"]*\)".*/\1/p' || true)"
  if [ -n "$id" ]; then
    printf '%s\n' "$id"
    return 0
  fi
  echo "mess_up: no containerId in result line (status $status): $last_line" >&2
  return 1
}

# mess_container_env writes the configured env of a container to stdout.
mess_container_env() {
  docker inspect "$1" --format '{{range .Config.Env}}{{println .}}{{end}}'
}

# mess_exec_env prints the effective environment the launcher path sees
# (docker exec, as the remote user).
mess_exec_env() {
  docker exec -u vscode "$1" env
}

# mess_cli_exec_env applies remoteEnv the way the CLI itself does: a second
# helper run with `devcontainer exec`, which is the path `docker exec` (the
# launcher path) does not take. $1 = project folder, $2 = container id.
mess_cli_exec_env() {
  local folder="$1" id="$2"
  docker run --rm \
    --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock \
    --mount type=bind,source="$folder",target="$folder" \
    "$MESS_CLI_IMAGE" \
    devcontainer exec --workspace-folder "$folder" --container-id "$id" env
}

# mess_cache_list shows the cache directory contents as the remote user.
# $1 = container id, $2 = path.
mess_cache_list() {
  docker exec -u vscode "$1" sh -c "ls -la \"$2\" 2>&1 || echo 'ls failed'"
}

# mess_write_marker writes a marker plus padding files into the cache
# directory as the remote user. $1 = container id, $2 = path, $3 = label.
mess_write_marker() {
  docker exec -u vscode "$1" sh -c "mkdir -p \"$2\" && printf '%s %s\n' \"$3\" \"\$(date -Is)\" > \"$2/marker.txt\" && dd if=/dev/zero of=\"$2/pad.bin\" bs=1024 count=64 2>/dev/null && ls -la \"$2\""
}

# mess_teardown removes a container and every volume of the current cell run;
# both errors are reported but the cleanup continues.
# $1 = container id (or empty), $2 = volume name (or empty), $3 = raw prefix.
mess_teardown() {
  local id="${1:-}" volume="${2:-}" prefix="$3"
  {
    if [ -n "$id" ]; then
      docker rm -f "$id" 2>&1 || true
    fi
    if [ -n "$volume" ]; then
      docker volume rm "$volume" 2>&1 || true
    fi
    echo "--- leftovers after teardown ---"
    docker ps -a --filter "name=devcontainer-" --format '{{.Names}} {{.Status}}' 2>&1
    docker volume ls --filter "name=$MESS_PREFIX" --format '{{.Name}}' 2>&1
  } >> "$prefix.log" 2>&1
}

# mess_new_project creates a fresh project directory (fresh folder means a
# fresh devcontainer.local_folder label, so two runs of one cell start two
# containers independently). $1 = label, e.g. "m1a-lauf1"; echoes the path.
mess_new_project() {
  local label="$1" dir
  dir="$(mktemp -d "/tmp/$MESS_PREFIX-${label}-XXXXXX")"
  printf '%s\n' "$dir"
}