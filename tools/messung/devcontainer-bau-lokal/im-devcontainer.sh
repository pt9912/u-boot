# im-devcontainer.sh — Ergänzung H zum Messlauf des Slice slice-v1-devcontainer-bau-lokal:
# Prüfung im ECHTEN Devcontainer (gestartet mit `make dev-up`, Features aus devcontainer.json).
# Vorbereitung auf dem Host (der Arbeitsbaum liegt dann unter /tmp/ws im Container; das
# Workspace-Volume mit dem Clone bleibt unangetastet):
#   tar -C "$PWD" -c --exclude=./.cache --exclude=./bin . | \
#     docker exec -i -u vscode <container> sh -c 'rm -rf /tmp/ws && mkdir -p /tmp/ws && tar -x -C /tmp/ws'
# Aufruf: u-boot --quiet devcontainer shell < tools/messung/devcontainer-bau-lokal/im-devcontainer.sh
timed() { # timed <name> <befehl…>
  local name="$1"; shift
  local t0 t1 rc
  t0=$(date +%s.%N); "$@" >"/tmp/$name.log" 2>&1; rc=$?; t1=$(date +%s.%N)
  printf '%-22s exit=%s wall=%ss\n' "$name" "$rc" "$(awk -v a="$t1" -v b="$t0" 'BEGIN{printf "%.1f", a-b}')"
  grep -E 'Befund\(e\)|coverage-gate: OK|FAIL|Error' "/tmp/$name.log" | tail -3 | cut -c1-160
}

echo "== Benutzer und Umgebung"
echo "whoami=$(whoami) uid=$(id -u) shell=$SHELL"
echo "BUILD_MODE=${BUILD_MODE:-<ungesetzt>} GOTOOLCHAIN=${GOTOOLCHAIN:-<ungesetzt>}"
echo "GOFLAGS=${GOFLAGS:-<ungesetzt>}"
echo "GOCACHE=${GOCACHE:-<ungesetzt>}"
echo "GOMODCACHE=${GOMODCACHE:-<ungesetzt>}"
echo "== Werkzeuge"
for c in go golangci-lint node claude ollama gh docker make git brew; do
  printf '%-14s' "$c"
  if command -v "$c" >/dev/null 2>&1; then
    case "$c" in
      go) go version 2>&1 | head -1 ;;
      ollama) ollama --version 2>&1 | tail -1 ;;
      golangci-lint) golangci-lint --version 2>&1 | head -1 | cut -c1-70 ;;
      *) "$c" --version 2>&1 | head -1 | cut -c1-70 ;;
    esac
  else
    echo FEHLT
  fi
done
echo "go liegt in: $(command -v go)"
echo "== Workspace-Volume"
ls -ld /workspaces/u-boot 2>&1
ls -ld /workspaces/u-boot/.cache 2>&1 | head -1
if mkdir -p "$GOCACHE" "$GOMODCACHE" && touch "$GOCACHE/.probe" "$GOMODCACHE/.probe"; then
  echo "GOCACHE und GOMODCACHE schreibbar"
  ls -ld "$GOCACHE" "$GOMODCACHE" | cut -c1-80
else
  echo "GOCACHE oder GOMODCACHE NICHT schreibbar"
fi
echo "== Docker-Daemon im Container"
docker info --format 'server {{.ServerVersion}}, storage {{.Driver}}' 2>&1 | head -1
echo "== make im lokalen Modus (Arbeitsbaum unter /tmp/ws)"
cd /tmp/ws || { echo "/tmp/ws fehlt"; exit 1; }
git config --global --add safe.directory /tmp/ws
timed modus-pruefung make modus-pruefung
timed gates-lauf-1 make gates
timed gates-lauf-2 make gates
for t in lint test coverage-gate; do timed "einzeln-$t" make "$t"; done
exit
