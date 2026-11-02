#!/usr/bin/env bash
# Messung M3-Zusatz: uebersteht der Cache unter dem Workspace-Volume ein
# `up --recreate`, wenn die ENV-Form die `m1c`-Form ist - containerEnv mit
# ${containerWorkspaceFolder} (Entscheidung des Projektinhabers 2026-11-02,
# Slice §4 Frage 3)?
#
# Zellen (je zwei Laeufe, frisches Projekt je Lauf, Sandbox-Form wie M3:
# workspaceMount "", workspaceFolder /workspaces/m, named Volume; das
# Dockerfile traegt keinen ENV - die ENV kommt allein aus containerEnv):
# - m3z-gocache: containerEnv {"GOCACHE": "${containerWorkspaceFolder}/.cache/go-build"}
# - m3z-gomod:   containerEnv {"GOMODCACHE": "${containerWorkspaceFolder}/.cache/gomod"}
#
# Ablauf je Lauf: up -> Marker + Pad-Datei in die Variable des Containers ->
# up --remove-existing-container (die recreate-Form des Adapters,
# devcontainer_launcher.go runArgs Zeilen 502-504) -> Variable, Marker, Pads
# und Eigentumsrechte im neuen Container pruefen.

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

mess_ensure_cli_image
mkdir -p "$MESS_DOC_DIR"

run_cell() {
  # $1 = cell id, $2 = env key written to containerEnv.
  local cell="$1" key="$2" run
  for run in 1 2; do
    local label="${cell}-lauf${run}"
    local dir prefix volume
    dir="$(mess_new_project "$label")"
    prefix="$MESS_DOC_DIR/rohdaten-${label}"
    volume="$MESS_PREFIX-${label}-ws"
    mkdir -p "$dir/.devcontainer"
    {
      echo "Messung M3-Zusatz, Zelle ${cell} (containerEnv ${key} = \${containerWorkspaceFolder}/.cache, Sandbox-Form, named volume ${volume})"
      date -Is
      mess_versions
    } > "$prefix.txt"
    echo "=== $label ($dir) ==="
    printf 'FROM %s\nUSER root\nRUN mkdir -p /workspaces/m && chown vscode:vscode /workspaces/m\nUSER vscode\n# Kein ENV im Block: die ENV-Form dieser Reihe ist containerEnv allein\n' \
      "$MESS_BASE_IMAGE" > "$dir/.devcontainer/Dockerfile"
    printf '{\n  "name": "%s",\n  "build": {"dockerfile": "Dockerfile", "context": "."},\n  "workspaceMount": "",\n  "workspaceFolder": "/workspaces/m",\n  "mounts": ["source=%s,target=/workspaces/m,type=volume"],\n  "containerEnv": {\n    "%s": "${containerWorkspaceFolder}/.cache"\n  },\n  "remoteUser": "vscode"\n}\n' \
      "$label" "$volume" "$key" > "$dir/.devcontainer/devcontainer.json"
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
      continue
    fi
    {
      echo "--- Container nach erstem up: $id ---"
      echo "--- ${key} im Container (docker exec env) ---"
      docker exec -u vscode "$id" env | grep "^${key}=" || echo "${key}: nicht gesetzt"
      echo "--- Marker schreiben ---"
      mess_write_marker "$id" "\$${key}" "${cell}-lauf${run}"
    } >> "$prefix.txt" 2>&1
    echo "$label: Marker geschrieben, jetzt recreate"

    local id2
    id2="$(mess_up "$dir" "$prefix" --remove-existing-container)" || id2=""
    if [ -z "$id2" ]; then
      echo "$label: recreate fehlgeschlagen (Fehler in $prefix.log)" >> "$prefix.txt"
      mess_teardown "$id" "$volume" "$prefix"
      rm -rf "$dir"
      continue
    fi
    {
      echo "--- Container nach recreate: $id2 (vorher: $id) ---"
      test "$id" != "$id2" && echo "Container-ID gewechselt: ja" \
        || echo "Container-ID gewechselt: nein"
      echo "--- ${key} nach recreate (docker exec env) ---"
      docker exec -u vscode "$id2" env | grep "^${key}=" || echo "${key}: nicht gesetzt"
      echo "--- Cache-Verzeichnis nach recreate ---"
      docker exec -u vscode "$id2" sh -c 'ls -la "$'"${key}"'" 2>&1'
      echo "--- Marker nach recreate ---"
      docker exec -u vscode "$id2" sh -c 'cat "$'"${key}"'/marker.txt" 2>&1' || true
      echo "--- Volume nach recreate: $(docker volume inspect "$volume" --format '{{.Name}} created {{.CreatedAt}}') ---"
    } >> "$prefix.txt" 2>&1
    mess_teardown "$id2" "$volume" "$prefix"
    rm -rf "$dir"
    echo "$label: fertig, Rohdaten unter $prefix.txt"
  done
}

run_cell m3z-gocache GOCACHE
run_cell m3z-gomod GOMODCACHE
echo "M3-Zusatz abgeschlossen."