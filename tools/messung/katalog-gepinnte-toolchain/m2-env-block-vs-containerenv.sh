#!/usr/bin/env bash
# Messung M2: greift ein `ENV` aus einem Dockerfile-Block im Container, und
# wie verhalten sich Block-`ENV` und `containerEnv` zu einander?
#
# Zellen (je zwei Laeufe, frisches Projekt je Lauf, Standard-Form; die Werte
# sind feste Pfade ohne Workspace-Bezug, damit die Zelle nur die Rang-Frage
# traegt, nicht die Pfad-Frage von M1):
# - m2a: ENV nur im Dockerfile-Block (containerEnv fehlt).
# - m2b: ENV in beiden, gleicher Wert.
# - m2c: ENV in beiden, verschiedener Wert.
#
# Messstellen je Lauf: `docker inspect .Config.Env` (die ENV-Liste des
# Containers) und `docker exec -u vscode env` (die wirksame Umgebung, der Weg
# des Launchers). Der Block-Slice im Dockerfile folgt dem Muster der
# Katalog-Bloecke (internal/hexagon/application/templates/devcontainer/
# Dockerfile.tmpl, ARG/RUN-Bloecke; der ENV-Block des Hand-written-Teils der
# eigenen .devcontainer/Dockerfile ab Zeile 107).

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

mess_ensure_cli_image
mkdir -p "$MESS_DOC_DIR"

readonly BLOCK_VALUE="/ub-messung-block/go-build"
readonly BLOCK_VALUE_MOD="/ub-messung-block/gomod"
readonly CONF_VALUE="/ub-messung-config/go-build"
readonly CONF_VALUE_MOD="/ub-messung-config/gomod"

mess_make_dockerfile() {
  # ENV only in the Dockerfile block; containerEnv is a separate file.
  mkdir -p "$1/.devcontainer"
  printf 'FROM %s\n# ENV-Namen im Block: GOCACHE/GOMODCACHE\nENV GOCACHE=%s \\\n    GOMODCACHE=%s\n' \
    "$MESS_BASE_IMAGE" "$BLOCK_VALUE" "$BLOCK_VALUE_MOD" > "$1/.devcontainer/Dockerfile"
}

run_cell() {
  local cell="$1" mode="$2" run
  for run in 1 2; do
    local label="${cell}-lauf${run}"
    local dir prefix
    dir="$(mess_new_project "$label")"
    prefix="$MESS_DOC_DIR/rohdaten-${label}"
    {
      echo "Messung M2, Zelle ${cell} (${mode}), Lauf ${run}"
      date -Is
      mess_versions
    } > "$prefix.txt"
    echo "=== $label ($dir) ==="
    mess_make_dockerfile "$dir"
    local conf_env=""
    if [ "$mode" = "gleich" ]; then
      conf_env="\"containerEnv\": {
    \"GOCACHE\": \"$BLOCK_VALUE\",
    \"GOMODCACHE\": \"$BLOCK_VALUE_MOD\"
  }"
    elif [ "$mode" = "verschieden" ]; then
      conf_env="\"containerEnv\": {
    \"GOCACHE\": \"$CONF_VALUE\",
    \"GOMODCACHE\": \"$CONF_VALUE_MOD\"
  }"
    fi
    if [ -n "$conf_env" ]; then
      printf '{\n  "name": "%s",\n  "build": {"dockerfile": "Dockerfile", "context": "."},\n  %s,\n  "remoteUser": "vscode"\n}\n' \
        "$label" "$conf_env" > "$dir/.devcontainer/devcontainer.json"
    else
      printf '{\n  "name": "%s",\n  "build": {"dockerfile": "Dockerfile", "context": "."},\n  "remoteUser": "vscode"\n}\n' \
        "$label" > "$dir/.devcontainer/devcontainer.json"
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
      mess_teardown "" "" "$prefix"
      rm -rf "$dir"
      continue
    fi
    {
      echo "--- docker inspect .Config.Env ---"
      mess_container_env "$id"
      echo "--- docker exec -u vscode env (Weg des Launchers) ---"
      mess_exec_env "$id"
    } >> "$prefix.txt" 2>&1
    mess_teardown "$id" "" "$prefix"
    rm -rf "$dir"
    echo "$label: fertig, Rohdaten unter $prefix.txt"
  done
}

run_cell m2a "nur-block"
run_cell m2b gleich
run_cell m2c verschieden
echo "M2 abgeschlossen."