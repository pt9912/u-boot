# Slice M3: Build-Polish (govulncheck-Pin + PROGRESS_FLAG)

> **Status:** Done
> **DoD:** Commit `987c164` + Review-Fix-Commit

## Auslöser

Vergleich des `u-boot`-Makefiles mit `../c-hsm-doc/Makefile` (gleiche
Build-Familie, gleiche Docker-only-Philosophie) hat zwei kleine
Lücken offengelegt, die u-boot übernehmen sollte:

1. **`govulncheck @latest`** in `make govulncheck` ist ein
   nicht-dokumentierter Carveout — Reproduzierbarkeit kippt, sobald
   eine neue Release rauskommt, und die Pin-Politik aus ADR-0004
   („alle Tool-Versionen pin-bar") gilt sonst überall (Go, golangci-lint,
   distroless-Image). `c-hsm-doc` setzt `GOVULNCHECK_VERSION ?= v1.1.4`
   und installiert `@$(GOVULNCHECK_VERSION)`.
2. **CI-Logs** waren in `make`-Aufrufen abhängig vom Docker-Default-
   Progress (`auto`), der im non-TTY-CI-Modus auf `tty` zurückfällt und
   knapp wird. `c-hsm-doc` setzt `PROGRESS_FLAG := --progress=plain`
   wenn `CI=1` — vollständige Logs für die GitHub-Actions, lokal
   bleibt `auto`.

Beides ist reine Build-Infra-Übernahme, kein neuer Carveout. Eigener
Mini-Slice, weil M3-T5 (depguard) bewusst auf den eigentlichen
Carveout-Scope beschränkt war (User-Wahl: „M3-T5 erst, dann
Polish-Slice").

## Geliefert

- `Makefile`:
  - `GOVULNCHECK_VERSION ?= v1.1.4` neu (Match mit c-hsm-doc).
  - `make govulncheck` installiert `@$(GOVULNCHECK_VERSION)` statt
    `@latest`. Routine-Bump dokumentiert sich künftig im Commit-Body
    (wie `GO_VERSION` / `GOLANGCI_LINT_VERSION`).
  - `PROGRESS_FLAG := --progress=plain` für `ifeq ($(CI),1)`; in
    `DOCKER_BUILD` vorangestellt. Lokaler Default bleibt `auto`.

## Out of Scope

- **Digest-Pinning per `*_BASE_IMAGE`-ARGs** aus dem c-hsm-doc-Dockerfile:
  passt thematisch in [`slice-v1-release-pipeline`](../open/slice-v1-release-pipeline.md)
  zusammen mit GHCR-Publish + Trivy-Scan (alle Release-Bauformen
  brauchen den Pinning-Pfad gleichzeitig, einzeln eingeführt wäre er
  totes Gewicht).
- **Markdown-Link-Validator (`docs-check`-Target + `tools/check_refs.py`)**:
  Aufhebungsplan steht in [`slice-v1-markdown-link-validator`](slice-v1-markdown-link-validator.md);
  c-hsm-doc's Python-Tool kann dort als Vorlage dienen. (Mittlerweile
  in derselben Sitzung 2026-05-27 vorgezogen + abgeschlossen.)
- **`HEALTHCHECK`** im runtime-Image (aus c-hsm-doc): nicht anwendbar
  — `u-boot` ist ein One-Shot-CLI ohne langlebigen Prozess.

## Bezug

- Auslösende Quelle: Vergleich `c-hsm-doc/Makefile` vs.
  `u-boot/Makefile` (Sitzung 2026-05-27).
- ADR-0004 (Pin-Politik) — `govulncheck`-Pin schließt die Lücke.
