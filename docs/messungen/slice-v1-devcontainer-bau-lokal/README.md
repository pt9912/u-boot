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
| `make gates` Lauf 1 (kalt) | grün, 42,7 s |
| `make gates` Lauf 2 (warm) | grün, 26,4 s |
| `make lint` / `make test` / `make coverage-gate` (warm, einzeln) | 0,7 s / 0,4 s / 1,0 s |
| Coverage | 92,10 % (Schwelle 90 %) |

Zum Vergleich der Docker-Modus auf dem Host nach dem Umbau des `Dockerfile`: `make gates` grün, 91,1 s (vorher 89,9 s und 89,3 s).

## Lesehilfen

- Nur der Lauf `go test ./...` kalt ist vollständig kalt. Die beiden `-coverprofile`-Läufe „kalt“ liefen danach im selben Cache-Volume; ihr Build-Cache war schon warm, die Zeiten sind darum nicht mit dem ersten vergleichbar. Aussagekräftig ist die Zahl aus dem Cache (Spalte „aus dem Cache“).
- Die Wanduhr der Docker-Läufe in D enthält Kontext-Übertragung und Image-Layer (Docker-Modus). Die 26,4 s des warmen lokalen Laufs (F) bestehen zum größten Teil aus den Gates, die kein Go brauchen: Die Einzelzeiten von `lint`, `test` und `coverage-gate` liegen zusammen bei etwa 2 s; der Rest ist `ablage-pruefung`, `docs-check` und `baseline-verify` (im Docker-Modus gemessen: 14,4 s, 6,3 s und 0,1 s, in F nicht einzeln gemessen).

## Was daraus folgt

- `go test` cached Ergebnisse auch mit `-coverpkg` und `-coverprofile`; das Profil aus dem gecachten Lauf trägt dieselbe Summe. Test und Coverage müssen im lokalen Modus kein gemeinsamer Lauf werden: [`SPEC-032`](../../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) bleibt unverändert, die Folge-ADR aus dem Architect-Report entfällt.
- Nach der Änderung einer Testdatei läuft ein Paket neu, nicht alle 24.
- Das Lint-Image bringt `bash` mit; seine Go-Version (1.27.1) weicht von `GO_VERSION` (1.27.3) ab. Das gilt für den Docker-Modus unverändert; im lokalen Modus läuft `golangci-lint` mit dem Go des Devcontainers.

## Grenzen

Ein Lauf je Zelle, ein Host (Linux `amd64`). Nicht gemessen: andere Hosts und Architekturen (macOS mit Colima, `arm64`), die Streuung der Zeiten, Verhalten bei einer geänderten Nicht-Testdatei (abhängige Pakete laufen laut Go ebenfalls neu). Zu F: Das Image kam aus `docker build` ohne die Features von `devcontainer.json` (`node`, `claude-code`, `homebrew`); der Workspace war ein Bind-Mount des Arbeitsbaums, nicht das Volume mit dem Clone, das der Devcontainer sonst nutzt; die Dateisystem-Eigenschaften beider können sich unterscheiden. Ein Start über `u-boot devcontainer up` mit dem neuen Dockerfile ist nicht gemessen.
