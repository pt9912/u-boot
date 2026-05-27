# syntax=docker/dockerfile:1.7

# ---------------------------------------------------------------------------
# u-boot — developer environment bootloader for Docker-based projects.
#
# Docker-only workflow per LH-FA-BUILD-007: the repository carries no
# host-side Go toolchain requirement. Build / lint / test / coverage all
# run via `docker build --target <stage> -t u-boot:<stage> .` and are
# wrapped in the Makefile.
#
# Stages (LH-FA-BUILD-001):
#   deps      — Go module resolution (cache layer).
#   compile   — Fast compile feedback without tests/lint.
#   lint      — golangci-lint with the project profile.
#   test      — `go test ./...`.
#   coverage  — `go test -coverprofile` + coverage-gate.sh
#               (bootstrap-aware, LH-FA-BUILD-008).
#   build     — Statically linked binary (CGO=0, -ldflags="-s -w").
#   runtime   — distroless/static:nonroot (LH-FA-BUILD-002).
#
# Pin policy (LH-FA-BUILD-003): GO_VERSION and GOLANGCI_LINT_VERSION are
# routine pins, hebung without separate ADR; rationale goes into the
# commit body.
# ---------------------------------------------------------------------------

# Global build args. Both must be declared before the first FROM so they
# are usable in every stage's FROM-line.
ARG GO_VERSION=1.26.3
ARG GOLANGCI_LINT_VERSION=v2.12.2

# ---- deps ------------------------------------------------------------------
FROM golang:${GO_VERSION} AS deps

WORKDIR /src
ENV GOFLAGS="-mod=readonly -buildvcs=false" \
    GOMODCACHE=/go/pkg/mod \
    GOCACHE=/root/.cache/go-build

COPY go.mod ./
# Same go.sum trick as k-deskflight: the [m] character class matches
# go.sum if present and silently matches nothing if it does not exist
# yet (pre-`go mod tidy` bootstrap state). Single COPY line covers both
# cases.
COPY go.su[m] ./

RUN mkdir -p "$GOMODCACHE" && go mod download

# ---- compile ---------------------------------------------------------------
FROM deps AS compile

COPY . .
RUN CGO_ENABLED=0 go build -o /tmp/u-boot ./cmd/uboot

# ---- lint ------------------------------------------------------------------
FROM golangci/golangci-lint:${GOLANGCI_LINT_VERSION}-alpine AS lint

WORKDIR /src
COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY . .
RUN golangci-lint run ./...

# ---- test ------------------------------------------------------------------
FROM deps AS test

COPY . .
RUN CGO_ENABLED=0 go test ./...

# ---- coverage --------------------------------------------------------------
# Bootstrap-aware (LH-FA-BUILD-008): when ./internal/... is empty, the
# stage sets COVERAGE_BOOTSTRAP=1 and coverage-gate.sh accepts an empty
# input as bootstrap-OK. Once production packages land in ./internal/...,
# the stage measures coverage against them; THRESHOLD is overridable via
# `make coverage-gate THRESHOLD=…`.
#
# `pipefail` is set explicitly via SHELL so that `go test … | tee …`
# propagates the `go test` exit code instead of being masked by tee.
FROM deps AS coverage

SHELL ["/bin/bash", "-eo", "pipefail", "-c"]

ARG COVERAGE_THRESHOLD=90
ENV COVERAGE_THRESHOLD=${COVERAGE_THRESHOLD}

COPY . .
RUN mkdir -p /out && \
    COVERPKG=$(go list ./internal/... 2>/dev/null | tr '\n' ',' | sed 's/,$//') && \
    if [ -z "$COVERPKG" ]; then \
        echo "coverage: no production packages in ./internal/... yet — bootstrap mode"; \
        : > /out/coverage.out; \
        : > /out/coverage-func.txt; \
        export COVERAGE_BOOTSTRAP=1; \
    else \
        CGO_ENABLED=0 go test \
            -coverpkg="$COVERPKG" \
            -coverprofile=/out/coverage.out \
            -covermode=atomic \
            ./... && \
        go tool cover -func=/out/coverage.out | tee /out/coverage-func.txt; \
        export COVERAGE_BOOTSTRAP=0; \
    fi && \
    bash scripts/coverage-gate.sh /out/coverage-func.txt "$COVERAGE_THRESHOLD"

# ---- build -----------------------------------------------------------------
FROM deps AS build

COPY . .
RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -o /out/u-boot \
    ./cmd/uboot

# ---- runtime ---------------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

LABEL org.opencontainers.image.source="https://github.com/pt9912/u-boot" \
      org.opencontainers.image.description="u-boot — a developer environment bootloader for Docker-based projects." \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.title="u-boot" \
      org.opencontainers.image.vendor="pt9912"

COPY --from=build /out/u-boot /usr/local/bin/u-boot

USER 65532:65532
ENTRYPOINT ["/usr/local/bin/u-boot"]
