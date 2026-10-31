# Messprotokoll: Messlauf M0 zu `slice-v1-devcontainer-bau-lokal`

Messung vor dem Code des Slice [`slice-v1-devcontainer-bau-lokal`](../../plan/planning/in-progress/slice-v1-devcontainer-bau-lokal.md) (Schritt M0 in dessen §3) für die Entscheidung [`ADR-0037`](../../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md), Entscheidung 7.

## Messstand

| Eigenschaft | Wert |
|---|---|
| Datum | 2026-10-31 (Lauf `messen.sh` 09:44 UTC; Ergänzung E danach) |
| Host | Linux `x86_64`, lokale Docker-Engine 29.9.2 |
| Go (Test-Container) | `golang:1.27.3` (Pin `GO_VERSION` aus `repo.mk`) |
| Lint-Image | `golangci/golangci-lint:v2.14.0` (Pin `GOLANGCI_LINT_VERSION` aus `repo.mk`) |
| Läufe | **ein Lauf je Zelle** (keine Streuung bestimmt) |
| Caches | zwei benannte Volumes (`uboot-m0-gocache`, `uboot-m0-mod` für A bis D; `uboot-m0e-*` für E), vor dem Lauf leer, danach entfernt |

## Reihenfolge und Aufruf

Aus der Repo-Wurzel, Rohdaten in dieses Verzeichnis:

```bash
bash tools/messung/devcontainer-bau-lokal/messen.sh > docs/messungen/slice-v1-devcontainer-bau-lokal/roh.txt
bash tools/messung/devcontainer-bau-lokal/messen-profil.sh >> docs/messungen/slice-v1-devcontainer-bau-lokal/roh.txt
```

Beide Skripte laufen über Docker, der Host braucht keine Go-Toolchain; die Quellen sind schreibgeschützt eingehängt, die Änderung einer Testdatei (B) geschieht in einer Kopie. Rohdaten: [`roh.txt`](roh.txt), Fehlerausgabe: `stderr.log` (leer).

## Ergebnisse

**A. Test-Cache von Go** (Quelle unverändert, 24 Pakete mit Tests):

| Lauf | Wanduhr | aus dem Cache |
|---|---|---|
| `go test ./...` kalt | 11,9 s | 0 von 24 |
| `go test ./...` warm | 0,9 s | 24 von 24 |
| mit `-coverpkg` und `-coverprofile` kalt | 5,5 s | 0 von 24 |
| mit `-coverpkg` und `-coverprofile` warm | 1,2 s | 24 von 24 |
| nur `-coverprofile` kalt / warm | 4,6 s / 0,8 s | 0 von 24 / 24 von 24 |

**B. Eine Testdatei geändert** (`internal/hexagon/domain/addon_dependency_test.go`, eine Kommentarzeile angehängt, in einer Kopie): `go test ./...` 1,4 s, `go test` mit `-coverprofile` 1,8 s; jeweils **1 von 24 Paketen** lief neu, 23 kamen aus dem Cache.

**C. Lint-Image** `golangci/golangci-lint:v2.14.0`: Go `go1.27.1` (nicht `1.27.3`), `golangci-lint` 2.14.0 gebaut mit `go1.27.0`, `bash` vorhanden (`/usr/bin/bash`), `sh` vorhanden.

**D. `make gates` im Docker-Modus:** zwei Läufe, beide grün, Wanduhr 89,9 s und 89,3 s.

**E. Coverage-Profil aus dem Cache** (Ergänzung, `messen-profil.sh`):

| Lauf | aus dem Cache | Profil | Summe `go tool cover -func` |
|---|---|---|---|
| kalt | 0 von 24 | 127585 Zeilen, 10988323 Bytes | 92,1 % |
| warm | 24 von 24 | 127585 Zeilen, 10988323 Bytes | 92,1 % |

Der Hash des Profils ist in beiden Läufen verschieden; Größe und Summe sind gleich. Warum der Hash abweicht (etwa die Reihenfolge der Blöcke), ist nicht untersucht.

**F. Lokaler Modus im Image des Devcontainers** (nach der Umsetzung; Skript `im-container.sh`, Rohdaten [`roh-lokal.txt`](roh-lokal.txt)): Das Image aus `.devcontainer/Dockerfile` (`docker build`, ohne die Features aus `devcontainer.json`), privilegiert gestartet mit einem Volume für `/var/lib/docker`, das Repo als Bind-Mount unter `/workspaces/u-boot`, verschachtelter Docker 29.9.2 mit `overlayfs`. Aufruf:

```bash
docker build -t uboot-devc-test -f .devcontainer/Dockerfile .devcontainer
docker run --rm --privileged -v uboot-test-docker:/var/lib/docker -v "$PWD":/workspaces/u-boot \
  -v "$PWD/tools/messung/devcontainer-bau-lokal":/m:ro -w /workspaces/u-boot uboot-devc-test bash /m/im-container.sh
```

| Messgröße | Ergebnis |
|---|---|
| Werkzeuge im Image | Go 1.27.3, `golangci-lint` 2.14.0; `BUILD_MODE=local`, `GOTOOLCHAIN=local`, `GOCACHE` und `GOMODCACHE` im Workspace |
| `make modus-pruefung` | OK (8 Targets, beide Modi, Teillauf für `gates` abgelehnt) |
| `make test PKG=./internal/hexagon/domain/...` | Teillauf, 5,0 s (erster Aufruf, mit Übersetzen) |
| `make gates PKG=./x` | abgelehnt (`repo.mk`, „PKG and RUN make a partial run, which is no evidence“) |
| `make gates` Lauf 1 (nach dem Teillauf davor, leeres Docker-Volume: das d-check-Image wird im Lauf geladen) | grün, 42,7 s |
| `make gates` Lauf 2 (warm) | grün, 26,4 s |
| `make lint` / `make test` / `make coverage-gate` (warm, einzeln) | 0,7 s / 0,4 s / 1,0 s |
| Coverage | 92,10 % (Schwelle 90 %) |

Zum Vergleich der Docker-Modus auf dem Host nach dem Umbau des `Dockerfile`: `make gates` grün, 91,1 s (vorher 89,9 s und 89,3 s).

**G. `govulncheck`, `build-binaries` und `test-docker` im lokalen Modus** (Skript `im-container-rest.sh`, Rohdaten [`roh-lokal-rest.txt`](roh-lokal-rest.txt); derselbe Aufbau wie F):

| Messgröße | Ergebnis |
|---|---|
| `make govulncheck` | Exit 0, 6,9 s, „No vulnerabilities found.“ |
| `make build-binaries` | Exit 0, 26,8 s, sechs Binaries (`linux`, `darwin`, `windows` je `amd64` und `arm64`) |
| `make test-docker` im verschachtelten Daemon (zwei Läufe) | Exit 2 in beiden: 652,0 s und 646,1 s; im zweiten Lauf 1621 bestanden, **9 fehlgeschlagen**, 1 übersprungen (oberste Ebene) |
| `make test-docker` im Docker-Modus am Host (Vergleich) | Exit 0, 1619 s, 1627 bestanden, **0 fehlgeschlagen**, 4 übersprungen |

Die 9 fehlgeschlagenen Tests des zweiten Laufs (dieselbe Zahl wie im ersten; die Namen des ersten sind nicht gesichert): `TestLauncher_Docker_S6_ValueIsNeverInImageLabelFileOrCommandLine`, `TestLauncher_Docker_S5_GuardFalsifiedIntegrityBreaksTheBuild`, `TestE2E_ForgeCLI_PinMismatchBreaksBuild`, `TestE2E_SandboxDevcontainer_NestedPodman`, `…_NestedPodman_Build`, `…_NestedPodman_BuildEgress`, `…_NestedPodman_PinMismatchBreaksBuild`, `…_Ollama_PinMismatchBreaksBuild`, `TestE2E_SandboxDevcontainer_StatusFiles_VFS`.

Einordnung nach den Meldungen im Log: Bei vier Tests (`NestedPodman`, `NestedPodman_Build`, `NestedPodman_BuildEgress`, `StatusFiles_VFS`) steht „restricted host: init exit 0, want 11“. Der Test erwartet bei `kernel.apparmor_restrict_unprivileged_userns=1` den Exit 11 des Startscripts und überspringt dann; im Docker-Modus am Host sind genau diese vier `SKIP`. Im verschachtelten Daemon endete das Startscript mit Exit 0. Bei vier weiteren (`S5`, `ForgeCLI_PinMismatch`, `NestedPodman_PinMismatch`, `Ollama_PinMismatch`) bricht der Bau im verschachtelten Daemon ab (`npm ci`, `curl`, `apt-get`), und der Test schlägt dennoch fehl; die erwartete Ursache (Prüfsumme) ist aus den gekürzten Zeilen nicht zu erkennen. Bei `S6` steht „the value was on a command line“. **Die Ursache dieser fünf ist nicht geklärt.**

## Lesehilfen

- Lauf 1 von `make gates` in F ist nicht vollständig kalt: Der Teillauf `make test PKG=…` davor hat den Modulcache gefüllt und das Paket `domain` übersetzt; kalt ist nur das Docker-Volume des verschachtelten Daemons.
- Nur der Lauf `go test ./...` kalt ist vollständig kalt. Die beiden `-coverprofile`-Läufe „kalt“ liefen danach im selben Cache-Volume; ihr Build-Cache war schon warm, die Zeiten sind darum nicht mit dem ersten vergleichbar. Aussagekräftig ist die Zahl aus dem Cache (Spalte „aus dem Cache“).
- Die Wanduhr der Docker-Läufe in D enthält Kontext-Übertragung und Image-Layer (Docker-Modus). Die 26,4 s des warmen lokalen Laufs (F) bestehen zum größten Teil aus den Gates, die kein Go brauchen: Die Einzelzeiten von `lint`, `test` und `coverage-gate` liegen zusammen bei etwa 2 s; der Rest ist `ablage-pruefung`, `docs-check` und `baseline-verify` (im Docker-Modus gemessen: 14,4 s, 6,3 s und 0,1 s, in F nicht einzeln gemessen).

## Was daraus folgt

- `go test` cached Ergebnisse auch mit `-coverpkg` und `-coverprofile`; das Profil aus dem gecachten Lauf trägt dieselbe Summe. Test und Coverage müssen im lokalen Modus kein gemeinsamer Lauf werden: [`SPEC-032`](../../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) bleibt unverändert, die Folge-ADR aus dem Architect-Report entfällt.
- Nach der Änderung einer Testdatei läuft ein Paket neu, nicht alle 24.
- Das Lint-Image bringt `bash` mit; seine Go-Version (1.27.1) weicht von `GO_VERSION` (1.27.3) ab. Das gilt für den Docker-Modus unverändert; im lokalen Modus läuft `golangci-lint` mit dem Go des Devcontainers.

- `govulncheck` und `build-binaries` laufen im lokalen Modus (G). `test-docker` im verschachtelten Daemon ist **nicht gleichwertig** zum Docker-Modus am Host (9 Fehlschläge gegen 0); der Nachweis für die Docker-Integrationstests bleibt der Docker-Modus an einem Host mit lokaler Engine und die CI. Die CI auf `main` (Stand `b371e262`) meldet `integration-docker` grün.

## Grenzen

Ein Lauf je Zelle, ein Host (Linux `amd64`). Nicht gemessen: andere Hosts und Architekturen (macOS mit Colima, `arm64`), die Streuung der Zeiten, Verhalten bei einer geänderten Nicht-Testdatei (abhängige Pakete laufen laut Go ebenfalls neu). Zu F: Das Image kam aus `docker build` ohne die Features von `devcontainer.json` (`node`, `claude-code`, `homebrew`); der Workspace war ein Bind-Mount des Arbeitsbaums, nicht das Volume mit dem Clone, das der Devcontainer sonst nutzt; die Dateisystem-Eigenschaften beider können sich unterscheiden. Ein Start über `u-boot devcontainer up` mit dem neuen Dockerfile ist nicht gemessen.
