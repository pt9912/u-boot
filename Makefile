# u-boot — developer environment bootloader for Docker-based projects.
#
# Docker-only workflow (LH-FA-BUILD-007): build/lint/test/coverage all
# run via `docker build --target <stage>` inside containers. The repo
# has no host-side Go toolchain requirement; only Docker and `make` are
# expected on the host. `make` is a deliberate carveout to
# LH-NFA-PORT-002 — it is the only host dependency beyond Docker.
#
# Quality gates (LH-FA-BUILD-006):
#   make gates       — inner-loop mandatory gates (lint + test + coverage-gate).
#   make ci          — gates plus govulncheck plus image-scan (mirrors ci.yml).
#   make fullbuild   — ci plus build (runtime image).

IMAGE                   ?= u-boot
GO_VERSION              ?= 1.26.3
GOLANGCI_LINT_VERSION   ?= v2.12.2
GOVULNCHECK_VERSION     ?= v1.1.4
PYTHON_VERSION          ?= 3.13-slim

# VERSION is injected at build time and becomes both `u-boot --version`
# output (via -X main.version) and the org.opencontainers.image.version
# OCI label on the runtime image. Default matches the in-source fallback
# in cmd/uboot/main.go (`var version = "0.1.0-dev"`); the publish.yml
# workflow passes VERSION=<tag-without-v> for tagged releases. CI gates
# (lint/test/coverage/govulncheck/image-scan) and local `make build`
# without override produce a coherent "0.1.0-dev" binary.
VERSION                 ?= 0.1.0-dev
# Trivy pin policy — TWO formats in play, both must be bumped together:
#   - Makefile (here):  Docker-Hub-Tag-Konvention OHNE `v`-Prefix
#                        → `aquasec/trivy:0.70.0`
#   - ci.yml::image-scan: GitHub-Release-Tag-Konvention MIT `v`-Prefix
#                        → `trivy-version: 'v0.70.0'` für trivy-action
# Bei jeder Pin-Hebung beide Stellen synchron heben (gleiche Trivy-
# Version, unterschiedliche Schreibweisen). Detector-/DB-Parität
# ist sonst gebrochen.
TRIVY_VERSION           ?= 0.70.0
THRESHOLD               ?= 90

# `--progress=plain` gives full, line-by-line BuildKit logs that survive
# CI log truncation. Locally the default (`auto`) keeps the compact TUI.
PROGRESS_FLAG :=
ifeq ($(CI),1)
PROGRESS_FLAG := --progress=plain
endif

# `--no-cache-filter <stage>` forces BuildKit to re-evaluate the given
# stage without invalidating the `deps` cache layer. Without this, a
# stale layer hash could mask test/lint/coverage failures.
NO_CACHE_FILTER_TEST     := --no-cache-filter test
NO_CACHE_FILTER_LINT     := --no-cache-filter lint
NO_CACHE_FILTER_COVERAGE := --no-cache-filter coverage

DOCKER_BUILD := docker build $(PROGRESS_FLAG) \
    --build-arg GO_VERSION=$(GO_VERSION) \
    --build-arg GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION) \
    --build-arg UBOOT_VERSION=$(VERSION)

.DEFAULT_GOAL := help

.PHONY: help deps compile lint test test-docker coverage coverage-gate build build-binaries run clean \
        gates ci fullbuild govulncheck image-scan verify-depguard docs-check

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# ---- inner-loop ------------------------------------------------------------

deps: ## Resolve Go module dependencies (deps-cache layer).
	$(DOCKER_BUILD) --target deps -t $(IMAGE):deps .

compile: ## Fast compile feedback (no tests/lint).
	$(DOCKER_BUILD) --target compile -t $(IMAGE):compile .

lint: ## golangci-lint with the project profile.
	$(DOCKER_BUILD) $(NO_CACHE_FILTER_LINT) --target lint -t $(IMAGE):lint .

test: ## Run `go test ./...` inside Docker.
	$(DOCKER_BUILD) $(NO_CACHE_FILTER_TEST) --target test -t $(IMAGE):test .

test-docker: ## Run //go:build docker integration tests against the host docker daemon.
	@# Two-step target (slice-m6-docker-integrationstests Sub-T4):
	@# 1) Build the `test-docker-tools` Dockerfile stage — golang +
	@#    docker-ce-cli + docker-compose-plugin (from the official
	@#    Docker repo). Cached after first build.
	@# 2) Run the test binary in that image with `--network=host` so
	@#    the test process and the Docker daemon share a network
	@#    namespace (required for `NetProbe.DialTCP("localhost", ...)`
	@#    to reach Compose-published ports — slice §Strukturelle
	@#    Bedingungen, Netzwerk-Namespace-Voraussetzung). The
	@#    Docker socket is also mounted so the test's Compose calls
	@#    reach the host daemon.
	$(DOCKER_BUILD) --target test-docker-tools -t $(IMAGE):test-docker-tools .
	docker run --rm --network=host \
	    -v "$(CURDIR)":/src -w /src \
	    -v /var/run/docker.sock:/var/run/docker.sock \
	    $(IMAGE):test-docker-tools \
	    go test -tags docker ./...

coverage-gate: ## Coverage threshold gate (bootstrap-aware, LH-FA-BUILD-008).
	$(DOCKER_BUILD) $(NO_CACHE_FILTER_COVERAGE) \
	    --target coverage \
	    --build-arg COVERAGE_THRESHOLD=$(THRESHOLD) \
	    -t $(IMAGE):coverage .

# Alias for ergonomics; the gate is the same target.
coverage: coverage-gate

build: ## Build the runtime image (distroless static, nonroot).
	$(DOCKER_BUILD) --target runtime -t $(IMAGE):latest .

run: build ## Smoke test: run `u-boot --help` from the built image.
	docker run --rm $(IMAGE):latest --help

# Cross-compiled binaries for the slice-v2-binary-distribution path.
# Output naming convention: bin/u-boot-<os>-<arch> (no ".tar.gz" — the
# release workflow handles packaging). Static linking + CGO=0 so the
# binaries run on minimal hosts. -ldflags `-X main.version=$(VERSION)`
# mirrors the runtime-image VERSION-Pin from the Dockerfile so the
# binary's `--version` matches the surrounding release. Docker-only
# build (LH-FA-BUILD-007): everything runs inside the pinned
# golang:$(GO_VERSION) container, no host Go toolchain required.
BIN_DIR    ?= bin
PLATFORMS  := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

build-binaries: ## Cross-compile u-boot binaries for all release platforms.
	@mkdir -p $(BIN_DIR)
	@for p in $(PLATFORMS); do \
	    os=$${p%/*}; arch=$${p#*/}; \
	    ext=""; [ "$$os" = "windows" ] && ext=".exe"; \
	    out=$(BIN_DIR)/u-boot-$$os-$$arch$$ext; \
	    echo "==> $$out (UBOOT_VERSION=$(VERSION))"; \
	    docker run --rm \
	        -v "$(CURDIR)":/src -w /src \
	        -e GOOS=$$os -e GOARCH=$$arch -e CGO_ENABLED=0 \
	        -e GOFLAGS=-buildvcs=false \
	        golang:$(GO_VERSION) \
	        go build -ldflags="-s -w -X main.version=$(VERSION)" -o $$out ./cmd/uboot \
	        || exit 1; \
	done
	@echo "[build-binaries] $(words $(PLATFORMS)) binaries built in $(BIN_DIR)/"

# ---- security gates --------------------------------------------------------

# govulncheck runs inside an ephemeral Go container with the project
# mounted in. Pinned via GOVULNCHECK_VERSION (ADR-0004 pin policy);
# routine upgrade documented in the commit body.
govulncheck: ## Run govulncheck against the project (LH-FA-BUILD-006).
	docker run --rm \
	    -v "$(CURDIR)":/src -w /src \
	    -e GOFLAGS=-buildvcs=false \
	    golang:$(GO_VERSION) \
	    sh -c "go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) && govulncheck ./..."

# image-scan reproduces the third PR-blocking ci.yml job locally.
# Builds the runtime image, then runs Trivy against it inside an
# ephemeral container; HIGH/CRITICAL findings fail the build
# (slice-v1-release-pipeline T3 / ADR-0007). Mounts the host Docker
# socket so Trivy can resolve the just-built local tag.
image-scan: build ## Trivy scan of the runtime image (LH-QA-003 third gate).
	docker run --rm \
	    -v /var/run/docker.sock:/var/run/docker.sock \
	    aquasec/trivy:$(TRIVY_VERSION) image \
	    --severity HIGH,CRITICAL --exit-code 1 \
	    --no-progress \
	    $(IMAGE):latest

# verify-depguard proves each of the eight LH-FA-ARCH-003 layer rules
# fires on a real forbidden import. Manual / on-demand: not part of
# `gates` because each iteration runs `make lint`, so a full pass takes
# several minutes. Re-run whenever the hexagonal layer set or the
# depguard config changes.
verify-depguard: ## Verify all eight depguard layer rules fire (LH-FA-ARCH-003).
	bash scripts/verify-depguard.sh

# ---- docs gates ------------------------------------------------------------

# docs-check validates relative markdown links across docs/, spec/, and
# root *.md. Docker-encapsulated Python (stdlib-only) — no host Python
# requirement. Adapted from c-hsm-doc.
docs-check: ## Validate relative markdown links in docs/, spec/, README*.
	docker run --rm \
	    -v "$(CURDIR)":/src -w /src \
	    python:$(PYTHON_VERSION) \
	    python tools/check_refs.py

# ---- aggregators -----------------------------------------------------------

gates: lint test coverage-gate docs-check ## Inner-loop mandatory gates.
	@echo "[gates] lint + test + coverage-gate + docs-check green"

ci: gates govulncheck image-scan ## Gates plus govulncheck plus image-scan (mirrors ci.yml).
	@echo "[ci] gates + govulncheck + image-scan green"

fullbuild: ci build ## CI plus runtime image (full closure).
	@echo "[fullbuild] ci + runtime image green"

# ---- maintenance -----------------------------------------------------------

clean: ## Remove local build artefacts and built images.
	@rm -rf out coverage $(BIN_DIR) *.out *.test
	@-docker image rm \
	    $(IMAGE):latest $(IMAGE):deps $(IMAGE):compile \
	    $(IMAGE):lint $(IMAGE):test $(IMAGE):coverage 2>/dev/null || true
	@echo "[clean] artefacts and images removed"
