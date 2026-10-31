set -u
# im-container-rest.sh — Ergänzung G zum Messlauf des Slice slice-v1-devcontainer-bau-lokal:
# die drei lokalen Rezepte, die im-container.sh nicht ausführt (govulncheck, build-binaries,
# test-docker). Aufruf wie im-container.sh (siehe docs/messungen/…/README.md, Abschnitt F).
# STEPS wählt die Schritte (Standard: alle drei), etwa `-e STEPS=test-docker`.
# Die vollständigen Logs liegen danach unter .cache/messung/ im gemounteten Repo (gitignored).
cd /workspaces/u-boot
steps="${STEPS:-govulncheck build-binaries test-docker}"
mkdir -p .cache/messung
echo "== Benutzer: $(whoami) uid=$(id -u)"
sudo sh /usr/local/bin/u-boot-docker-init >/tmp/docker-init.log 2>&1 || echo "docker-init Exit $?"
for i in $(seq 1 30); do docker info >/dev/null 2>&1 && break; sleep 1; done
echo "== Docker: $(docker info --format 'server {{.ServerVersion}}, storage {{.Driver}}' 2>&1 | head -1)"
git config --global --add safe.directory /workspaces/u-boot
echo "BUILD_MODE=${BUILD_MODE:-<ungesetzt>} STEPS=$steps"

timed() { # timed <name> <befehl…>
  local name="$1"; shift
  local t0 t1 rc
  t0=$(date +%s.%N); "$@" >".cache/messung/$name.log" 2>&1; rc=$?; t1=$(date +%s.%N)
  printf '%-16s exit=%s wall=%ss\n' "$name" "$rc" "$(awk -v a="$t1" -v b="$t0" 'BEGIN{printf "%.1f", a-b}')"
  tail -n "${TAIL:-3}" ".cache/messung/$name.log" | cut -c1-200
}

has() { case " $steps " in *" $1 "*) return 0 ;; esac; return 1; }

if has govulncheck; then
  echo "== make govulncheck"
  timed govulncheck make govulncheck
fi
if has build-binaries; then
  echo "== make build-binaries"
  timed build-binaries make build-binaries
  ls -1 bin 2>/dev/null | head -8
fi
if has test-docker; then
  echo "== make test-docker (bis zu 30 Minuten)"
  TAIL=6 timed test-docker timeout 1800 make test-docker
  log=.cache/messung/test-docker.log
  echo "== test-docker: Zählung (oberste Ebene)"
  printf 'PASS=%s FAIL=%s SKIP=%s\n' "$(grep -c -E '^--- PASS' "$log")" "$(grep -c -E '^--- FAIL' "$log")" "$(grep -c -E '^--- SKIP' "$log")"
  echo "== fehlgeschlagene Tests"
  grep -E '^--- FAIL' "$log" | cut -c1-160
  echo "== Paketergebnisse"
  grep -E '^(ok|FAIL|\?)\s' "$log" | grep -v '^ok' | cut -c1-160
fi
