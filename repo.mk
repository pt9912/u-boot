# repo.mk — die Make-Targets dieses Repos.
#
# Diese Datei gehoert dem Repo. Der Bootstrap legt sie an, wo sie fehlt, und
# schreibt sie danach nie wieder; harness/mk/ und das Root-Makefile dagegen
# schreibt jeder Lauf neu. Das Root-Makefile bindet sie nach den Fragmenten
# unter harness/mk/ ein (-include repo.mk); fehlt sie, laeuft make ohne sie.
#
# Ein eigenes Target steht hier wie in jedem Makefile. Ein eigenes Gate haengt
# sich an die Gate-Kette mit
#   GATE_CHECKS += <target>
# und laeuft dann in `make gates`, vor dem Gate-Nachweis. Jedes Target dieser
# Datei, Gate oder nicht, bekommt seine Zeile in harness/README.md, Abschnitt
# Sensors: das Doku-Gate (.d-check.yml, Modul targets) meldet ein Target ohne
# Zeile als gate-undocumented.
#
# Eine weitere Make-Datei, die diese Datei per include einbindet, liest das
# Doku-Gate nur im Wurzelverzeichnis (Glob "*.mk" in makefiles:). Grenze: liegt
# sie in einem Unterverzeichnis, laufen ihre Targets mit make, das Doku-Gate
# sieht sie nicht und meldet nichts.
#
# Eine Vorgabe fuer einen Marker der Fragmente (dort mit ?= belegt) steht hier
# mit einfachem = und gewinnt, weil diese Datei nach den Fragmenten gelesen wird.

# ---- u-boot ----------------------------------------------------------------
#
# Die Targets dieses Repos. Das Makefile schreibt ai-harness-init (MR-012); die
# Gate-Kette bildet GATE_CHECKS, `gates` und `record-gates` stehen dort.

# u-boot — developer environment bootloader for Docker-based projects.
#
# Docker-only workflow (LH-FA-BUILD-007): in the Docker mode (BUILD_MODE=docker,
# the default) build/lint/test/coverage all run via
# `docker build --target <stage>` inside containers. The repo has no
# host-side Go toolchain requirement; only Docker and `make` are expected on
# the host. `make` is a deliberate carveout to LH-NFA-PORT-002 — it is the
# only host dependency beyond Docker. The devcontainer sets BUILD_MODE=local
# (.devcontainer/Dockerfile): there the targets run with the toolchain
# installed in the container (scripts/bau.sh); CI and hosts stay in the
# Docker mode.
#
# Quality gates (LH-FA-BUILD-006):
#   make gates       — inner-loop mandatory gates (lint + test + coverage-gate + docs-check + baseline-verify + ablage-pruefung).
#   make ci          — gates plus govulncheck plus image-scan (mirrors ci.yml).
#   make fullbuild   — ci plus build (runtime image).

IMAGE                   ?= u-boot
# Hebung von GO_VERSION und GOLANGCI_LINT_VERSION betrifft drei Stellen:
# Dockerfile (ARG), diese Datei und .devcontainer/Dockerfile (hand-written part).
GO_VERSION              ?= 1.27.3
GOLANGCI_LINT_VERSION   ?= v2.14.0
GOVULNCHECK_VERSION     ?= v1.8.0
PYTHON_VERSION          ?= 3.13-slim

# VERSION is injected at build time and becomes both `u-boot --version`
# output (via -X main.version) and the org.opencontainers.image.version
# OCI label on the runtime image. Default matches the in-source fallback
# in cmd/uboot/main.go (`var version = "0.16.1-dev"`); the publish.yml
# workflow passes VERSION=<tag-without-v> for tagged releases. CI gates
# (lint/test/coverage/govulncheck/image-scan) and local `make build`
# without override produce a coherent "0.16.1-dev" binary.
VERSION                 ?= 0.16.1-dev
# Trivy pin policy — TWO formats in play, both must be bumped together:
#   - repo.mk (here):  Docker-Hub-Tag-Konvention OHNE `v`-Prefix
#                        → `aquasec/trivy:0.70.0`
#   - ci.yml::image-scan: GitHub-Release-Tag-Konvention MIT `v`-Prefix
#                        → Eingabe `version: 'v0.70.0'` der trivy-action
# Bei jeder Pin-Hebung beide Stellen synchron heben (gleiche Trivy-
# Version, unterschiedliche Schreibweisen). Detector-/DB-Parität
# ist sonst gebrochen.
TRIVY_VERSION           ?= 0.70.0
THRESHOLD               ?= 90

# BUILD_MODE selects where the build targets run: `docker` (default) builds
# the Dockerfile stages, `local` runs the commands with the toolchain of the
# environment (the devcontainer). The mode is explicit; there is no detection
# of a `go` on the PATH.
BUILD_MODE              ?= docker
ifeq ($(filter $(BUILD_MODE),docker local),)
$(error BUILD_MODE must be docker or local, not '$(BUILD_MODE)')
endif

# PKG and RUN narrow `make test` to a partial run (local mode only). Only a
# value given on the command line counts; a PKG or RUN in the environment is
# ignored. A partial run is no evidence: it is refused together with the
# evidence targets.
TEST_PKG := $(if $(filter command%,$(origin PKG)),$(PKG))
TEST_RUN := $(if $(filter command%,$(origin RUN)),$(RUN))
ifneq ($(strip $(TEST_PKG)$(TEST_RUN)),)
ifeq ($(BUILD_MODE),docker)
$(error PKG and RUN exist in the local mode only (BUILD_MODE=local))
endif
ifneq ($(filter gates record-gates ci fullbuild,$(MAKECMDGOALS)),)
$(error PKG and RUN make a partial run, which is no evidence: not with gates, record-gates, ci or fullbuild)
endif
endif

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

.PHONY: deps compile lint test test-docker coverage coverage-gate build build-binaries run clean \
        ci fullbuild govulncheck packages image-scan verify-depguard ablage-pruefung modus-pruefung

# ---- inner-loop ------------------------------------------------------------

ifeq ($(BUILD_MODE),local)

# Local mode: the targets run with the toolchain of the environment (the
# devcontainer, .devcontainer/Dockerfile). The commands come from
# scripts/bau.sh, the script the Dockerfile stages call as well. The build and
# test cache of Go stays between runs, so a repeated run is not a rebuild.

deps: ## Resolve Go module dependencies (Docker mode: deps stage; local mode: go mod download).
	go mod download

compile: ## Fast compile feedback, no tests/lint (Docker mode: compile stage; local mode: go build).
	bash scripts/bau.sh compile

lint: ## golangci-lint with the project profile (Docker mode: lint stage; local mode: golangci-lint run).
	bash scripts/bau.sh lint

test: ## Run the Go tests (Docker mode: test stage; local mode: PKG=… RUN=… makes a partial run, which is no evidence).
	PKG="$(TEST_PKG)" RUN="$(TEST_RUN)" bash scripts/bau.sh test

test-docker: ## Run //go:build docker integration tests against a Docker daemon (Docker mode: the host daemon; local mode: the daemon of the environment).
	go test -v -tags docker -timeout 30m ./...

coverage-gate: ## Coverage threshold gate, bootstrap-aware (LH-FA-BUILD-008; Docker mode: coverage stage; local mode: scripts/bau.sh).
	OUT=coverage COVERAGE_THRESHOLD=$(THRESHOLD) bash scripts/bau.sh coverage

else

deps:
	$(DOCKER_BUILD) --target deps -t $(IMAGE):deps .

compile:
	$(DOCKER_BUILD) --target compile -t $(IMAGE):compile .

lint:
	$(DOCKER_BUILD) $(NO_CACHE_FILTER_LINT) --target lint -t $(IMAGE):lint .

test:
	$(DOCKER_BUILD) $(NO_CACHE_FILTER_TEST) --target test -t $(IMAGE):test .

test-docker:
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
	@#
	@#    `/tmp` wird host-shared zum Test-Container gemountet:
	@#    `t.TempDir()` schreibt unter `/tmp/...`, und Compose-Bind-
	@#    Mount-Pfade (z. B. die `otel-collector-config.yaml` aus
	@#    slice-v1-otel) werden vom Daemon (Host) aufgelöst, nicht
	@#    vom Test-Container. Ohne den Share sieht der Daemon nur
	@#    leere Verzeichnisse statt der vom Test geschriebenen
	@#    Dateien (slice-v1-otel T3 CI-Diagnose: Collector
	@#    `failed to get config: ... is a directory`).
	@#
	@#    `-timeout 30m`: das e2e-Paket braucht zwischen sieben und elf
	@#    Minuten (Colima, arm64), das Standardlimit von `go test` sind
	@#    zehn; der CI-Job (`integration.yml`) hat 30 Minuten.
	@#
	@#    `-v`: Pass und Skip jedes Tests stehen im Log, auch im CI-Lauf; ein
	@#    übersprungener Test ist dort von einem bestandenen unterscheidbar.
	$(DOCKER_BUILD) --target test-docker-tools -t $(IMAGE):test-docker-tools .
	docker run --rm --network=host \
	    -v "$(CURDIR)":/src -w /src \
	    -v /var/run/docker.sock:/var/run/docker.sock \
	    -v /tmp:/tmp \
	    $(IMAGE):test-docker-tools \
	    go test -v -tags docker -timeout 30m ./...

coverage-gate:
	$(DOCKER_BUILD) $(NO_CACHE_FILTER_COVERAGE) \
	    --target coverage \
	    --build-arg COVERAGE_THRESHOLD=$(THRESHOLD) \
	    -t $(IMAGE):coverage .

endif

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
# binary's `--version` matches the surrounding release. The container
# emits a tar stream and the host tar process extracts it into $(BIN_DIR),
# so mounted worktrees do not receive root-owned output files.
BIN_DIR    ?= bin
PLATFORMS  := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

packages: build-binaries ## Build the .deb and .rpm packages (amd64, arm64) with nfpm (slice-v2-distro-pakete).
	./scripts/build-packages.sh $(VERSION) $(BIN_DIR) packages

ifeq ($(BUILD_MODE),local)

build-binaries: ## Cross-compile u-boot binaries for all release platforms (Docker mode: binary-archive stage; local mode: go build; the release pipeline uses the Docker mode).
	@mkdir -p $(BIN_DIR)
	@bash -c 'set -o pipefail; out=$$(mktemp -d); trap "rm -rf $$out" EXIT; UBOOT_VERSION="$(VERSION)" PLATFORMS="$(PLATFORMS)" OUT_DIR="$$out" sh scripts/build-binary-archive.sh | tar --no-same-owner -xf - -C "$(BIN_DIR)"'
	@echo "[build-binaries] $(words $(PLATFORMS)) binaries built in $(BIN_DIR)/ (local mode)"

else

build-binaries:
	@mkdir -p $(BIN_DIR)
	$(DOCKER_BUILD) --target binary-archive -t $(IMAGE):binary-archive .
	@docker run --rm \
	    -e PLATFORMS="$(PLATFORMS)" \
	    $(IMAGE):binary-archive \
	    | tar --no-same-owner -xf - -C $(BIN_DIR)
	@echo "[build-binaries] $(words $(PLATFORMS)) binaries built in $(BIN_DIR)/"

endif

# ---- security gates --------------------------------------------------------

# govulncheck runs inside an ephemeral Go container with the project
# mounted in. Pinned via GOVULNCHECK_VERSION (ADR-0004 pin policy);
# routine upgrade documented in the commit body.
ifeq ($(BUILD_MODE),local)

govulncheck: ## Run govulncheck against the project (LH-FA-BUILD-006; Docker mode: golang container; local mode: the toolchain of the environment).
	GOFLAGS=-buildvcs=false go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)
	GOFLAGS=-buildvcs=false "$$(go env GOPATH)/bin/govulncheck" ./...

else

govulncheck:
	docker run --rm \
	    -v "$(CURDIR)":/src -w /src \
	    -e GOFLAGS=-buildvcs=false \
	    golang:$(GO_VERSION) \
	    sh -c "go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) && govulncheck ./..."

endif

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

# verify-depguard proves each of the nine LH-FA-ARCH-003 layer rules
# fires on a real forbidden import. Manual / on-demand: not part of
# `gates` because each iteration runs `make lint`, so a full pass takes
# several minutes. Re-run whenever the hexagonal layer set or the
# depguard config changes.
verify-depguard: ## Verify all nine depguard layer rules fire (LH-FA-ARCH-003).
	bash scripts/verify-depguard.sh

# modus-pruefung pins the two modes of BUILD_MODE with `make -n`: the targets
# that have a local variant call no `docker build` in the local mode and call
# `docker build` or `docker run` in the Docker mode (the default); a partial
# run (PKG, RUN) is refused together with the evidence targets. On demand, not
# part of `gates` (a dry run costs seconds, but the rule it guards changes
# rarely). A mutation that turns it red: a `docker build` line in a local
# recipe, or a local recipe without its Docker counterpart.
MODUS_TARGETS := deps compile lint test test-docker coverage-gate govulncheck build-binaries
modus-pruefung: ## Check that the local mode calls no docker build and the Docker mode does (make -n).
	@set -e; \
	for t in $(MODUS_TARGETS); do \
	    if $(MAKE) -s -n BUILD_MODE=local $$t 2>&1 | grep -q 'docker build'; then \
	        echo "modus-pruefung: $$t calls docker build in the local mode"; exit 1; fi; \
	    if ! $(MAKE) -s -n BUILD_MODE=docker $$t 2>&1 | grep -q 'docker \(build\|run\)'; then \
	        echo "modus-pruefung: $$t calls no docker in the Docker mode"; exit 1; fi; \
	done; \
	for v in PKG=./x RUN=X; do \
	    if $(MAKE) -s -n BUILD_MODE=local gates $$v >/dev/null 2>&1; then \
	        echo "modus-pruefung: gates accepts $$v (a partial run is no evidence)"; exit 1; fi; \
	done; \
	echo "modus-pruefung: OK ($(words $(MODUS_TARGETS)) targets, both modes; partial run refused for gates)"

# ---- gate chain ------------------------------------------------------------

# docs-check (d-check.mk) und baseline-verify (harness/mk/baseline.mk) haengen sich selbst
# an GATE_CHECKS; hier kommen die Code-Gates dazu.
GATE_CHECKS += lint test coverage-gate

# ablage-pruefung: Positivliste der Ablage im Planning-Baum und unter docs/
# (harness/conventions.md, Abschnitt Ablage von Messprotokollen und Skripten).
# bash + git auf dem Host, kein Docker (Muster tools/harness/history-range-guard.sh);
# das Skript und sein Waechter-Test liegen unter tools/harness/.
ablage-pruefung: ## Ablage im Planning-Baum und unter docs/ gegen die Positivliste pruefen.
	@bash tools/harness/ablage-pruefung.sh
	@bash tools/harness/ablage-pruefung-test.sh

GATE_CHECKS += ablage-pruefung

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

# ---- Devcontainer ----------------------------------------------------------

include tools/mk/devcontainer.mk
