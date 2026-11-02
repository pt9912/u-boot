#!/usr/bin/env bash
# Messung M3: uebersteht der Cache unter dem Workspace-Volume ein
# `up --recreate`?
#
# Form, die gemessen wird: der Dockerfile-Block traegt die ENV (die
# M2-Zellen zeigen, welche Rang-Frage containerEnv hat); das Projekt folgt
# der Sandbox-Form der erzeugten Konfiguration (workspaceMount "",
# workspaceFolder /workspaces/<name>, named volume - templates/devcontainer/
# devcontainer.json.tmpl Zeilen 12-19, Dockerfile.tmpl Zeilen 151-157).
# Je Lauf ein frisches Projekt mit frischem Volume:
#   up -> Marker + Pad-Dateien in GOCACHE/GOMODCACHE -> up
#   --remove-existing-container (die recreate-Form des Launchers,
#   devcontainer_launcher.go runArgs Zeilen 502-504) -> Marker pruefen.
#
# Zwei Laeufe.

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

mess_ensure_cli_image
mkdir -p "$MESS_DOC_DIR"

readonly CACHE_GOCACHE="/workspaces/m/.cache/go-build"
readonly CACHE_GOMOD="/workspaces/m/.cache/gomod"

run() {
  local run_no="$1"
  local label="m3-lauf${run_no}"
  local dir prefix
  dir="$(mess_new_project "$label")"
  prefix="$MESS_DOC_DIR/rohdaten-${label}"
  local volume="$MESS_PREFIX-${label}-ws"
  mkdir -p "$dir/.devcontainer"
  {
    echo "Messung M3, Lauf ${run_no} (Dockerfile-ENV-Form, Sandbox-Form, named volume ${volume})"
    date -Is
    mess_versions
  } > "$prefix.txt"
  echo "=== $label ($dir) ==="
  printf 'FROM %s\nUSER root\nRUN mkdir -p /workspaces/m && chown vscode:vscode /workspaces/m\nUSER vscode\n# ENV im Block, Pfade unter dem Workspace-Volume\nENV GOCACHE=%s \\\n    GOMODCACHE=%s\n' \
    "$MESS_BASE_IMAGE" "$CACHE_GOCACHE" "$CACHE_GOMOD" > "$dir/.devcontainer/Dockerfile"
  printf '{\n  "name": "%s",\n  "build": {"dockerfile": "Dockerfile", "context": "."},\n  "workspaceMount": "",\n  "workspaceFolder": "/workspaces/m",\n  "mounts": ["source=%s,target=/workspaces/m,type=volume"],\n  "remoteUser": "vscode"\n}\n' \
    "$label" "$volume" > "$dir/.devcontainer/devcontainer.json"
  {
    echo "--- devcontainer.json ---"
    cat "$dir/.devcontainer/devcontainer.json"
    echo "--- Dockerfile ---"
    cat "$dir/.devcontainer/Dockerfile"
  } >> "$prefix.txt"

  local id
  id="$(mess_up "$dir" "$prefix")" || id=""
  if [ -z "$id" ]; then
    echo "$label: erster up fehlgeschlagen (Fehler in $prefix.log)" >> "$prefix.txt"
    mess_teardown "" "$volume" "$prefix"
    rm -rf "$dir"
    return 0
  fi
  {
    echo "--- Container nach erstem up: $id ---"
    echo "--- GOCACHE: Marker schreiben ---"
    mess_write_marker "$id" "$CACHE_GOCACHE" "go-build-lauf${run_no}"
    echo "--- GOMODCACHE: Marker schreiben ---"
    mess_write_marker "$id" "$CACHE_GOMOD" "gomod-lauf${run_no}"
  } >> "$prefix.txt" 2>&1
  echo "$label: Marker geschrieben, jetzt recreate"

  local id2
  id2="$(mess_up "$dir" "$prefix" --remove-existing-container)" || id2=""
  if [ -z "$id2" ]; then
    echo "$label: recreate fehlgeschlagen (Fehler in $prefix.log)" >> "$prefix.txt"
    mess_teardown "$id" "$volume" "$prefix"
    rm -rf "$dir"
    return 0
  fi
  {
    echo "--- Container nach recreate: $id2 (vorher: $id) ---"
    test "$id" != "$id2" && echo "Container-ID gewechselt: ja" \
      || echo "Container-ID gewechselt: nein"
    echo "--- GOCACHE nach recreate ---"
    mess_cache_list "$id2" "$CACHE_GOCACHE"
    echo "--- Marker GOCACHE nach recreate ---"
    docker exec -u vscode "$id2" cat "$CACHE_GOCACHE/marker.txt" 2>&1 || true
    echo "--- GOMODCACHE nach recreate ---"
    mess_cache_list "$id2" "$CACHE_GOMOD"
    echo "--- Marker GOMODCACHE nach recreate ---"
    docker exec -u vscode "$id2" cat "$CACHE_GOMOD/marker.txt" 2>&1 || true
    echo "--- Volume nach recreate: $(docker volume inspect "$volume" --format '{{.Name}} created {{.CreatedAt}}') ---"
  } >> "$prefix.txt" 2>&1
  mess_teardown "$id2" "$volume" "$prefix"
  rm -rf "$dir"
  echo "$label: fertig, Rohdaten unter $prefix.txt"
}

run 1
run 2
echo "M3 abgeschlossen."