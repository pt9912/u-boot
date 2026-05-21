#!/usr/bin/env bash
# coverage-gate.sh — bootstrap-aware Go coverage gate (LH-FA-BUILD-008).
#
# Usage:
#   coverage-gate.sh <coverage-func.txt> <threshold>
#
# Reads the output of `go tool cover -func=<profile>` and enforces the
# overall total coverage against the given threshold (percent, integer or
# float).
#
# Bootstrap mode:
#   When COVERAGE_BOOTSTRAP=1 is set in the environment, an empty input
#   file is treated as "no production code yet" and the gate passes with
#   threshold 0. This avoids a false-green during the MVP bootstrap phase
#   before ./internal/... contains any production packages.
#
#   COVERAGE_BOOTSTRAP=0 (or unset) treats an empty input as a hard
#   failure (exit 2) so that real `go test` regressions cannot be masked.

set -euo pipefail

if [[ $# -ne 2 ]]; then
  echo "usage: $0 <coverage-func.txt> <threshold>" >&2
  exit 2
fi

func_file="$1"
threshold="$2"

if [[ ! -f "$func_file" ]]; then
  echo "coverage-gate: input file not found: $func_file" >&2
  exit 2
fi

if [[ ! -s "$func_file" ]]; then
  if [[ "${COVERAGE_BOOTSTRAP:-0}" == "1" ]]; then
    echo "coverage-gate: empty coverage input — bootstrap mode (threshold 0)"
    exit 0
  fi
  echo "coverage-gate: empty coverage input and COVERAGE_BOOTSTRAP != 1" >&2
  echo "hint: did `go test` fail before producing coverage data?" >&2
  exit 2
fi

# `go tool cover -func` final line: "total:\t(statements)\tXX.X%"
total_line="$(grep -E '^total:' "$func_file" || true)"
if [[ -z "$total_line" ]]; then
  echo "coverage-gate: no 'total:' line in $func_file" >&2
  echo "hint: did `go test -coverprofile` actually run any tests?" >&2
  exit 2
fi

# Extract trailing percent number (e.g., "85.7%" → "85.7").
total_pct="$(echo "$total_line" | grep -oE '[0-9]+\.[0-9]+%?$' | tr -d '%')"
if [[ -z "$total_pct" ]]; then
  echo "coverage-gate: could not parse coverage percent from: $total_line" >&2
  exit 2
fi

# awk-based comparison handles fractional thresholds.
pass="$(awk -v p="$total_pct" -v t="$threshold" 'BEGIN { print (p+0 >= t+0) ? 1 : 0 }')"
if [[ "$pass" != "1" ]]; then
  printf "coverage-gate: FAIL — coverage %.2f%% below threshold %s%%\n" "$total_pct" "$threshold" >&2
  exit 1
fi

printf "coverage-gate: OK — coverage %.2f%% meets threshold %s%%\n" "$total_pct" "$threshold"
