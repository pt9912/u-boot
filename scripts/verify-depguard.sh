#!/usr/bin/env bash
# verify-depguard.sh — verify all eight depguard rules from .golangci.yml
# fire on a real forbidden import (LH-FA-ARCH-003, spec/architecture.md §4).
#
# Background:
#   depguard has been active since M2b but matched nothing while
#   ./internal/... was empty. With the M3 init-flow, every layer has at
#   least one production package; this script proves each rule is alive
#   by injecting one deliberate forbidden import per rule, running
#   `make lint`, asserting the expected desc appears in the output, and
#   reverting the injection. Resolves the depguard-leer-Match carveout
#   (slice-m3-depguard-aktivierung-verifizieren).
#
# Usage:
#   scripts/verify-depguard.sh
#
# Exit codes:
#   0 — all eight rules fired with the expected desc
#   1 — at least one rule did not fire as expected (see log files)
#   2 — preconditions not met (dirty working tree, missing tools)
#
# Notes:
#   - Refuses to run on a dirty tree so injection / cleanup cannot be
#     misattributed to user changes.
#   - Forbidden imports are chosen to avoid Go import cycles; see the
#     CASES table for the rationale per rule.
#   - One `make lint` invocation per rule; expect ~3–5 minutes total
#     wall-clock with a warm Docker cache.

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

LAYER_DIRS=(
  internal/hexagon/domain
  internal/hexagon/application
  internal/hexagon/port
  internal/adapter
)
dirty="$(git status --porcelain -- "${LAYER_DIRS[@]}")"
if [ -n "$dirty" ]; then
  echo "[verify-depguard] FATAL: working tree dirty in a layer dir; commit or stash first" >&2
  echo "$dirty" >&2
  exit 2
fi

VIOLATION_FILE=""
cleanup() {
  if [ -n "$VIOLATION_FILE" ] && [ -f "$VIOLATION_FILE" ]; then
    rm -f "$VIOLATION_FILE"
  fi
}
trap cleanup EXIT INT TERM

# Each row: rule|target-dir|forbidden-import|expected-desc-substring
#
# The forbidden import is picked so it cannot form a Go import cycle
# with the target package — depguard only runs after typecheck, so a
# cycle would mask the rule. The expected desc must match a deny entry
# in .golangci.yml for the named rule.
CASES=(
  "domain-isoliert|internal/hexagon/domain|github.com/pt9912/u-boot/internal/hexagon/port/driven|domain must not depend on port"
  "application-no-adapter|internal/hexagon/application|github.com/pt9912/u-boot/internal/adapter/driven/clock|application must depend on ports, not on adapter implementations"
  "port-no-application|internal/hexagon/port/driving|github.com/pt9912/u-boot/internal/adapter/driven/clock|port must not depend on adapter"
  "port-driving-no-driven|internal/hexagon/port/driving|github.com/pt9912/u-boot/internal/hexagon/port/driven|driving port must not depend on driven port"
  "port-driven-no-driving|internal/hexagon/port/driven|github.com/pt9912/u-boot/internal/hexagon/port/driving|driven port must not depend on driving port"
  "adapter-no-application|internal/adapter/driven/fs|github.com/pt9912/u-boot/internal/hexagon/application|adapter must implement ports, not consume application"
  "adapter-driving-no-driven|internal/adapter/driving/cli|github.com/pt9912/u-boot/internal/adapter/driven/fs|driving adapter must not depend on driven adapter"
  "adapter-driven-no-driving|internal/adapter/driven/clock|github.com/pt9912/u-boot/internal/adapter/driving/cli|driven adapter must not depend on driving adapter"
)

LOG_DIR="$(mktemp -d -t verify-depguard-XXXXXX)"
echo "[verify-depguard] logs in $LOG_DIR"

failures=0
for case in "${CASES[@]}"; do
  IFS='|' read -r rule dir bad_import desc <<< "$case"

  pkg_name="$(basename "$dir")"
  VIOLATION_FILE="$dir/verify_depguard_violation.go"
  log="$LOG_DIR/$rule.log"

  # No leading comment (would trip revive's package-comments per-file
  # rule against an extra-doc form), and the blank import sits inside an
  # import block with an explanatory comment (revive's blank-imports
  # requires a justification adjacent to the import).
  cat > "$VIOLATION_FILE" <<EOF
package $pkg_name

import (
	// Forbidden import to verify the '$rule' depguard rule;
	// removed automatically by scripts/verify-depguard.sh.
	_ "$bad_import"
)
EOF

  echo "[verify-depguard] $rule: injecting $bad_import into $dir/"

  set +e
  make lint > "$log" 2>&1
  rc=$?
  set -e

  rm -f "$VIOLATION_FILE"
  VIOLATION_FILE=""

  if [ $rc -eq 0 ]; then
    echo "[verify-depguard] FAIL $rule: make lint succeeded but should have failed"
    echo "  log: $log"
    failures=$((failures + 1))
    continue
  fi

  if ! grep -F -q "$desc" "$log"; then
    echo "[verify-depguard] FAIL $rule: expected desc not found in lint output"
    echo "  expected: $desc"
    echo "  log:      $log"
    failures=$((failures + 1))
    continue
  fi

  echo "[verify-depguard] PASS $rule (desc: \"$desc\")"
done

if [ $failures -ne 0 ]; then
  echo "[verify-depguard] $failures rule(s) failed verification"
  exit 1
fi

echo "[verify-depguard] all 8 depguard rules verified"
