# repo.mk — die Make-Targets dieses Repos.
#
# Diese Datei gehoert dem Repo. Der Bootstrap legt sie an, wo sie fehlt, und
# schreibt sie danach nie wieder; harness/mk/ und das Root-Makefile dagegen
# schreibt jeder Lauf neu. Das Root-Makefile bindet sie nach den Fragmenten
# ein (-include repo.mk); fehlt sie, laeuft make ohne sie.
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
# Docker-only workflow (LH-FA-BUILD-007), single build path (ADR-0043): the
# devcontainer of the repo is the build environment, and the recipes of
# scripts/bau.sh call the toolchain installed in the container (the
# containerEnv of .devcontainer/devcontainer.json carries the Go caches under
# the workspace). `make gates` runs inside the devcontainer container; its
# container work (docs-check with the d-check image) is carried by the nested
# daemon of the container. `make test-docker` binds to the Docker daemon
# outside the devcontainer (socket form): the integration tests are not
# promised in the nested daemon (SPEC-031, Messbasis M4). The build of the
# runtime image, the release binaries, packages and govulncheck run against a
# Docker daemon as before. The host carries Docker and GNU make, no host-side
# Go toolchain requirement (`make` is the deliberate carveout to
# LH-NFA-PORT-002 — it is the only host dependency beyond Docker).
#
# Quality gates (LH-FA-BUILD-006):
#   make gates       — inner-loop mandatory gates (lint + test + coverage-gate + docs-check + baseline-verify + ablage-pruefung).
#   make ci          — gates plus govulncheck plus image-scan (mirrors ci.yml).
#   make fullbuild   — ci plus build (runtime image).

IMAGE                   ?= u-boot
# Hebung von GO_VERSION und GOLANGCI_LINT_VERSION (LH-FA-DEV-011): die
# Quellen-Stelle des Pins sind die Build-Argumente des Katalog-Templates der
# zwei Toolchain-Schluessel (`go-toolchain`, `golangci-lint`); gleich sind
# vier Stellen: dieses Dockerfile (ARG), jenes Template, diese Datei (?=) und
# die erzeugte .devcontainer/Dockerfile des eigenen Projekts (ARG). Der
# Pin-Stellen-Test (generate_toolchain_test.go) liest alle vier. Die Routine
# der Hebung ohne eigenen ADR traegt ADR-0041, Abschnitt Konsequenzen
# (docs/plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md).
GO_VERSION              ?= 1.27.3
GOLANGCI_LINT_VERSION   ?= v2.15.0
GOVULNCHECK_VERSION     ?= v1.8.0
PYTHON_VERSION          ?= 3.13-slim

# VERSION is injected at build time and becomes both `u-boot --version`
# output (via -X main.version) and the org.opencontainers.image.version
# OCI label on the runtime image. Default matches the in-source fallback
# in cmd/uboot/main.go (`var version = "0.24.0-dev"`); the publish.yml
# workflow passes VERSION=<tag-without-v> for tagged releases. CI gates
# (lint/test/coverage/govulncheck/image-scan) and local `make build`
# without override produce a coherent "0.24.0-dev" binary.
VERSION                 ?= 0.24.0-dev
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

# The build-mode decision is retired (ADR-0043): the devcontainer is the build
# environment, and there is no second mode. The value `local`, which the
# containerEnv of the devcontainer sets, is accepted; every other value,
# notably the retired `docker`, ends with this make error.
BUILD_MODE              ?= local
ifneq ($(filter-out local,$(BUILD_MODE)),)
$(error BUILD_MODE must be local or unset, not '$(BUILD_MODE)'; the repo carries a single build path, the devcontainer (docs/plan/adr/0043-nur-devcontainer-bau-pfad.md))
endif

# PKG and RUN narrow `make test` to a partial run. Only a value given on the
# command line counts; a PKG or RUN in the environment is ignored. A partial
# run is no evidence: it is refused together with the evidence targets.
TEST_PKG := $(if $(filter command%,$(origin PKG)),$(PKG))
TEST_RUN := $(if $(filter command%,$(origin RUN)),$(RUN))
ifneq ($(strip $(TEST_PKG)$(TEST_RUN)),)
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

DOCKER_BUILD := docker build $(PROGRESS_FLAG) \
    --build-arg GO_VERSION=$(GO_VERSION) \
    --build-arg GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION) \
    --build-arg UBOOT_VERSION=$(VERSION)

.DEFAULT_GOAL := help

.PHONY: deps compile lint test test-docker coverage coverage-gate build build-binaries run clean \
        ci fullbuild govulncheck packages image-scan verify-depguard ablage-pruefung

# ---- inner-loop ------------------------------------------------------------

# The targets that start no container call the recipes of scripts/bau.sh — the
# same commands the Dockerfile stages call. PKG and RUN reach the script
# through the environment (BAU_PKG, BAU_RUN), not through the command line of
# the recipe: a quote in the value stays text. They apply to `test` only.
export BAU_PKG := $(TEST_PKG)
export BAU_RUN := $(TEST_RUN)

deps: ## Resolve Go module dependencies (the toolchain of the build environment; the containerEnv of the devcontainer sets the Go caches under the workspace).
	go mod download

compile: ## Fast compile feedback, no tests/lint (the recipe of scripts/bau.sh, the same command as the compile stage).
	bash scripts/bau.sh compile

lint: ## golangci-lint with the project profile (the recipe of scripts/bau.sh, the same command as the lint stage).
	bash scripts/bau.sh lint

test: ## Run the Go tests (the recipe of scripts/bau.sh; PKG=… RUN=… makes a partial run, which is no evidence; write a $ in RUN as $$).
	bash scripts/bau.sh test

# Two-step target (slice-m6-docker-integrationstests Sub-T4), socket form
# (SPEC-031): the tests bind to the Docker daemon outside the devcontainer
# container — the daemon that creates it.
test-docker: ## Run //go:build docker integration tests against the Docker daemon outside the devcontainer (socket form; SPEC-031, the tests are not promised in the nested daemon).
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

coverage-gate: ## Coverage threshold gate, bootstrap-aware (LH-FA-BUILD-008; the recipes of scripts/bau.sh, the same commands as the coverage stage).
	OUT=coverage COVERAGE_THRESHOLD=$(THRESHOLD) bash scripts/bau.sh coverage

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

build-binaries: ## Cross-compile u-boot binaries for all release platforms (the binary-archive stage; the release pipeline uses this path).
	@mkdir -p $(BIN_DIR)
	$(DOCKER_BUILD) --target binary-archive -t $(IMAGE):binary-archive .
	@docker run --rm \
	    -e PLATFORMS="$(PLATFORMS)" \
	    $(IMAGE):binary-archive \
	    | tar --no-same-owner -xf - -C $(BIN_DIR)
	@echo "[build-binaries] $(words $(PLATFORMS)) binaries built in $(BIN_DIR)/"

# packages builds its binaries over the binary-archive stage: a package is a
# release artefact (SPEC-031).
packages: ## Build the .deb and .rpm packages (amd64, arm64) with nfpm (slice-v2-distro-pakete); binaries from the binary-archive stage.
	$(MAKE) build-binaries
	./scripts/build-packages.sh $(VERSION) $(BIN_DIR) packages

# ---- security gates --------------------------------------------------------

# govulncheck runs inside an ephemeral Go container with the project
# mounted in. Pinned via GOVULNCHECK_VERSION (ADR-0004 pin policy);
# routine upgrade documented in the commit body.
govulncheck: ## Run govulncheck against the project (LH-FA-BUILD-006; ephemeral Go container).
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

# verify-depguard proves each of the nine LH-FA-ARCH-003 layer rules
# fires on a real forbidden import. Manual / on-demand: not part of
# `gates` because each iteration runs `make lint`, so a full pass takes
# several minutes. Re-run whenever the hexagonal layer set or the
# depguard config changes.
verify-depguard: ## Verify all nine depguard layer rules fire (LH-FA-ARCH-003).
	bash scripts/verify-depguard.sh

# ---- gate chain ------------------------------------------------------------

# docs-check (d-check.mk) und baseline-verify (harness/mk/baseline.mk) haengen sich selbst
# an GATE_CHECKS; hier kommen die Code-Gates dazu.
GATE_CHECKS += lint test coverage-gate

include tools/mk/*.mk  # the repo-owned fragments (ablage-pruefung, tool-pruefung, devcontainer)

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
# (das Fragment tools/mk/devcontainer.mk traegt der Glob oben ein)