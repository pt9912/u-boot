#!/usr/bin/env bash
# Messung M1: loest die Dev-Container-CLI eine Workspace-Pfad-Form in
# `containerEnv` bzw. `remoteEnv` auf?
#
# Zellen (je zwei Laeufe, frisches Projekt je Lauf):
# - m1a: containerEnv mit ${localWorkspaceFolder}, Standard-Form (die CLI
#        bind-mountet den Workspace selbst).
# - m1b: containerEnv mit ${localWorkspaceFolder}, Sandbox-Form
#        (workspaceMount "", workspaceFolder unter /workspaces, named volume).
# - m1c: containerEnv mit ${containerWorkspaceFolder}, Sandbox-Form.
# - m1d: remoteEnv mit ${localWorkspaceFolder}, Standard-Form; Messstelle ist
#   das `docker exec` (Weg des Launchers) und zusaetzlich `devcontainer exec`
#   (Weg der CLI), weil remoteEnv nur auf dem zweiten Weg wirken kann.
#
# Quellen der Formen: die erzeugte devcontainer.json traegt heute nur
# remoteEnv mit ${localEnv:...} (internal/hexagon/application/templates/
# devcontainer/devcontainer.json.tmpl Zeilen 33-38); containerEnv rendert
# kein Pfad der Pipeline. Ziel der Pruefung: GOCACHE/GOMODCACHE je Projekt
# unter dem Workspace-Volume (Slice §1, Punkt ContainerEnv).

set -euo pipefail

# Read the shared helpers of the measurement in place.
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

mess_ensure_cli_image
mkdir -p "$MESS_DOC_DIR"

mess_make_dockerfile() {
  # $1 = project dir, $2 = workspace folder to pre-create (empty = none).
  mkdir -p "$1/.devcontainer"
  if [ -n "$2" ]; then
    printf 'FROM %s\nUSER root\nRUN mkdir -p %s && chown vscode:vscode %s\nUSER vscode\n' \
      "$MESS_BASE_IMAGE" "$2" "$2" > "$1/.devcontainer/Dockerfile"
  else
    printf 'FROM %s\n' "$MESS_BASE_IMAGE" > "$1/.devcontainer/Dockerfile"
  fi
}

mess_write_conf() {
  # $1 = project dir, rest = devcontainer.json body (valid JSON object).
  mkdir -p "$1/.devcontainer"
  printf '%s\n' "$2" > "$1/.devcontainer/devcontainer.json"
}

run_cell() {
  # $1 = cell id (m1a..m1d), $2 = shape (standard|sandbox), $3 = variable form
  # (local|container) — only evaluated for the containerEnv cells, $4 = env key
  # (containerEnv|remoteEnv).
  local cell="$1" shape="$2" var="$3" key="$4" run
  for run in 1 2; do
    local label="${cell}-lauf${run}"
    local dir prefix
    dir="$(mess_new_project "$label")"
    prefix="$MESS_DOC_DIR/rohdaten-${label}"
    {
      echo "Messung M1, Zelle ${cell} (${key}, ${var}-Workspace-Form, ${shape}-Form), Lauf ${run}"
      date -Is
      mess_versions
      echo "Host-Architektur: $(dpkg --print-architecture 2>/dev/null || uname -m)"
    } > "$prefix.txt"
    echo "=== $label ($dir) ==="
    local ws_folder volume_extra volume_name
    ws_folder="/workspaces/$(basename "$dir")"
    volume_extra=""
    if [ "$shape" = "sandbox" ]; then
      mess_make_dockerfile "$dir" "$ws_folder"
      volume_name="$MESS_PREFIX-${label}-ws"
      volume_extra="\"workspaceMount\": \"\",
  \"workspaceFolder\": \"$ws_folder\",
  \"mounts\": [\"source=$volume_name,target=$ws_folder,type=volume\"],
  "
    else
      mess_make_dockerfile "$dir" ""
      volume_name=""
    fi
    local var_value
    if [ "$var" = "local" ]; then
      # Single quotes keep the literal; bash does not re-expand the content
      # when the variable joins a double-quoted string later.
      var_value='${localWorkspaceFolder}'
    else
      var_value='${containerWorkspaceFolder}'
    fi
    if [ "$key" = "containerEnv" ]; then
      mess_write_conf "$dir" "{
  \"name\": \"$label\",
  \"build\": {\"dockerfile\": \"Dockerfile\", \"context\": \".\"},
  $volume_extra\"containerEnv\": {
    \"GOCACHE\": \"$var_value/.cache/go-build\",
    \"GOMODCACHE\": \"$var_value/.cache/gomod\"
  },
  \"remoteUser\": \"vscode\"
}"
    else
      mess_write_conf "$dir" "{
  \"name\": \"$label\",
  \"build\": {\"dockerfile\": \"Dockerfile\", \"context\": \".\"},
  $volume_extra\"remoteEnv\": {
    \"GOCACHE\": \"$var_value/.cache/go-build\",
    \"GOMODCACHE\": \"$var_value/.cache/gomod\"
  },
  \"remoteUser\": \"vscode\"
}"
    fi
    {
      echo "--- devcontainer.json ---"
      cat "$dir/.devcontainer/devcontainer.json"
      echo "--- Dockerfile ---"
      cat "$dir/.devcontainer/Dockerfile"
    } >> "$prefix.txt"
    local id
    id="$(mess_up "$dir" "$prefix")" || id=""
    if [ -z "$id" ]; then
      echo "$label: up fehlgeschlagen (Fehler in $prefix.log)" >> "$prefix.txt"
      mess_teardown "" "$volume_name" "$prefix"
      rm -rf "$dir"
      continue
    fi
    {
      echo "--- docker inspect .Config.Env ---"
      mess_container_env "$id"
      echo "--- docker exec -u vscode env (Weg des Launchers) ---"
      mess_exec_env "$id"
      if [ "$key" = "remoteEnv" ]; then
        echo "--- devcontainer exec env (Weg der CLI; remoteEnv wirkt nur hier) ---"
        mess_cli_exec_env "$dir" "$id" || true
      fi
    } >> "$prefix.txt" 2>&1
    mess_teardown "$id" "$volume_name" "$prefix"
    rm -rf "$dir"
    echo "$label: fertig, Rohdaten unter $prefix.txt"
  done
}

run_cell m1a standard local containerEnv
run_cell m1b sandbox local containerEnv
run_cell m1c sandbox container containerEnv
run_cell m1d standard local remoteEnv
echo "M1 abgeschlossen."