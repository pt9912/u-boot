#!/usr/bin/env bash
# bau.sh — die Befehle der Stages compile, lint, test und coverage.
#
# ZUSAGE. Das Dockerfile (Stages compile, lint, test, coverage) und der lokale
# Modus von repo.mk (BUILD_MODE=local) rufen dieses Skript; beide führen damit
# dieselben Befehle aus. Das Skript setzt CGO_ENABLED=0 bei compile, test und
# coverage und sonst nichts. GOFLAGS kommt aus der Umgebung (Docker: ENV der
# Stage deps; lokal: ENV im .devcontainer/Dockerfile). lint setzt nichts: Die
# Stage baut auf dem Image golangci/golangci-lint und kennt weder GOFLAGS noch
# CGO_ENABLED.
#
# AUFRUF.   bau.sh compile | lint | test | coverage
# UMGEBUNG. compile:  BAU_COMPILE_OUT  Ausgabepfad des Binary (Standard /tmp/u-boot)
#           test:     BAU_PKG  Paketliste, durch Leerzeichen getrennt (Standard ./...),
#                     BAU_RUN  Filter für -run (ein Teillauf; er ist kein Nachweis,
#                     siehe repo.mk). Die Werte kommen über die Umgebung, nicht über
#                     die Kommandozeile: ein Anführungszeichen im Filter bleibt Text.
#           coverage: OUT  Verzeichnis für Profil und Funktionsliste (Standard /out),
#                     COVERAGE_THRESHOLD  Schwelle in Prozent (Pflicht)
#
# KOPPLUNG. Der Aufruf für coverage steht im Dockerfile (Stage coverage) und in
# repo.mk (coverage-gate); das Bootstrap-Verhalten bei leerem ./internal/...
# beschreibt scripts/coverage-gate.sh.
set -eo pipefail

stage="${1:-}"

case "$stage" in
  compile)
    CGO_ENABLED=0 go build -o "${BAU_COMPILE_OUT:-/tmp/u-boot}" ./cmd/uboot
    ;;
  lint)
    golangci-lint run ./...
    ;;
  test)
    args=()
    if [[ -n "${BAU_RUN:-}" ]]; then
      args+=(-run "$BAU_RUN")
    fi
    # Die Paketliste wird an Leerzeichen getrennt, ohne dass die Shell Muster
    # (etwa ein Stern) gegen Dateien auflöst.
    read -r -a pkgs <<<"${BAU_PKG:-./...}"
    CGO_ENABLED=0 go test ${args[@]+"${args[@]}"} "${pkgs[@]}"
    ;;
  coverage)
    : "${COVERAGE_THRESHOLD:?COVERAGE_THRESHOLD is required}"
    out="${OUT:-/out}"
    mkdir -p "$out"
    COVERPKG=$(go list ./internal/... 2>/dev/null | tr '\n' ',' | sed 's/,$//')
    if [ -z "$COVERPKG" ]; then
      echo "coverage: no production packages in ./internal/... yet — bootstrap mode"
      : > "$out/coverage.out"
      : > "$out/coverage-func.txt"
      export COVERAGE_BOOTSTRAP=1
    else
      CGO_ENABLED=0 go test \
        -coverpkg="$COVERPKG" \
        -coverprofile="$out/coverage.out" \
        -covermode=atomic \
        ./...
      go tool cover -func="$out/coverage.out" | tee "$out/coverage-func.txt"
      export COVERAGE_BOOTSTRAP=0
    fi
    bash scripts/coverage-gate.sh "$out/coverage-func.txt" "$COVERAGE_THRESHOLD"
    ;;
  *)
    echo "usage: $0 compile|lint|test|coverage" >&2
    exit 2
    ;;
esac
