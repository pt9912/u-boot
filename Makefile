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
#   make ci          — gates plus govulncheck.
#   make fullbuild   — ci plus build (runtime image).

IMAGE                   ?= u-boot
GO_VERSION              ?= 1.26.3
GOLANGCI_LINT_VERSION   ?= v2.12.2
GOVULNCHECK_VERSION     ?= v1.1.4
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
    --build-arg GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION)

.DEFAULT_GOAL := help

.PHONY: help deps compile lint test coverage coverage-gate build run clean \
        gates ci fullbuild govulncheck verify-depguard

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

# verify-depguard proves each of the eight LH-FA-ARCH-003 layer rules
# fires on a real forbidden import. Manual / on-demand: not part of
# `gates` because each iteration runs `make lint`, so a full pass takes
# several minutes. Re-run whenever the hexagonal layer set or the
# depguard config changes.
verify-depguard: ## Verify all eight depguard layer rules fire (LH-FA-ARCH-003).
	bash scripts/verify-depguard.sh

# ---- aggregators -----------------------------------------------------------

gates: lint test coverage-gate ## Inner-loop mandatory gates.
	@echo "[gates] lint + test + coverage-gate green"

ci: gates govulncheck ## Gates plus govulncheck.
	@echo "[ci] gates + govulncheck green"

fullbuild: ci build ## CI plus runtime image (full closure).
	@echo "[fullbuild] ci + runtime image green"

# ---- maintenance -----------------------------------------------------------

clean: ## Remove local build artefacts and built images.
	@rm -rf out coverage *.out *.test
	@-docker image rm \
	    $(IMAGE):latest $(IMAGE):deps $(IMAGE):compile \
	    $(IMAGE):lint $(IMAGE):test $(IMAGE):coverage 2>/dev/null || true
	@echo "[clean] artefacts and images removed"
