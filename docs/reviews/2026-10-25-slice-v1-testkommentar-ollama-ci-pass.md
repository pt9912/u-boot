# Review-Report: slice-v1-testkommentar-ollama-ci-pass — 2026-10-25

**Review-Art:** Code (Diff gegen Plan, Hard Rules und Konventionen; kein DoD-Urteil)

**Gegenstand:** Commit-Range `09ddacdd..18c6a68c` (Inhalts-Commit `18c6a68c`)

**Skill:** `.harness/skills/reviewer.md` @ 3ec45359 · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-25

**Eingangs-Kontext:**

- [`slice-v1-testkommentar-ollama-ci-pass`](../plan/planning/done/slice-v1-testkommentar-ollama-ci-pass.md): Befund, DoD, Plan-/Beleg-Tabelle, Risiken, §7
- [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 10 (gelesen, unverändert)
- [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)
- `AGENTS.md` §3.7, §3.9, §3.10
- `internal/e2e/sandbox_devcontainer_docker_test.go` Zeilen 1225 bis 1335, `repo.mk` Ziel `test-docker`, `.github/workflows/integration.yml`

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der neue Schluss „The architecture is not asserted in the test; its pass shows in the CI output only while …“ hat im Kommentar keinen Bezugspunkt für „the architecture“: der Satz über den „architecture claim“ ist gestrichen, die Zusage nennt keine Architektur. „its pass“ bindet grammatisch an „the architecture“ statt an den Test. | `AGENTS.md` §3.7 (Kommentar schreibt an den Änderer; Zusage/Abgrenzung muss aus dem Kommentar selbst lesbar sein) | `internal/e2e/sandbox_devcontainer_docker_test.go` · „The architecture is not asserted in the test; its pass shows in the CI output only while“ | Review-only | Kommentar-Verweiswort ohne Bezug im Text |
| F-2 | LOW | Die Plan-Tabelle führt die Testdatei weiter als „update (nur Kommentarzeilen 1246 bis 1253)“; der Kommentar belegt nach dem Diff die Zeilen 1246 bis 1254 (eine Zeile mehr). | `AGENTS.md` §3.10 (Abgleich der Plan-Tabelle nach der letzten Fix-Runde) | Slice-Plan, Plan-Tabelle · „update (nur Kommentarzeilen 1246 bis 1253)“ | Review-only (`git diff --stat` gegen Tabelle) | Plan-Tabelle zeigt Vor-Diff-Zeilenbereich |
| F-3 | LOW | Die Zeile des Slice in `open/README.md` verlinkt jetzt auf `in-progress/`, trägt aber weiter „Plan liegt in `open/`“. Der Zustandsmarker stimmt nicht mit dem Link in derselben Zeile überein. | `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt), §3.7 (Zustandsfelder) | `docs/plan/planning/open/README.md` · „Plan liegt in `open/`“ | Review-only (`make docs-check` ist grün, prüft das nicht) | Zustandsmarker nach Lifecycle-Move nicht nachgezogen |
| F-4 | INFO | Zuständigkeit Verifier: Die Wahrheit der Kopplung an `-v` ist hier gegen `repo.mk` (`go test -v -tags docker -timeout 30m ./...` im Ziel `test-docker`) und `integration.yml` (`run: make test-docker`) gelesen und stimmt. Ob ein CI-Log mit `-v` den Pass zeigt, ist kein Gegenstand dieses Reviews. | Maintainability | `repo.mk` · `go test -v -tags docker` | Review-only | – |

## Prüfpunkte im Einzelnen

1. **Diff nur Kommentar:** `git diff 09ddacdd..HEAD -- internal/` zeigt 3 Einfügungen und 2 Löschungen, alle Zeilen mit `//`; keine Testlogik, `repo.mk` unverändert. `git diff --stat` nennt vier Dateien: Testdatei, Slice, `open/README.md`, `roadmap.md`.
2. **Kommentar gegen §3.7/§3.9:**
   - Zusage (Image, Loopback, CPU, Pull über Allowlist, Fremdhost gesperrt, Modell überlebt `rm`) steht unverändert; gegen `ollamaImageScript`, `ollamaServerScript`, `ollamaRestartScript` und die `want`-Listen gelesen, stimmt.
   - Grenze (Netz, rund 1,5 GB) bleibt.
   - „The architecture is not asserted in the test“ stimmt gegen den Rumpf: nur `requireDocker(t)`, keine Architektur-Bedingung, kein `runtime.GOARCH`.
   - Kopplung stimmt gegen `repo.mk` und `integration.yml`.
   - (a) kein „former/früher/until/no longer/ohne X wäre“, keine Run-ID im Code-Kommentar.
   - (b) die genannten Symbole `test-docker`, `repo.mk`, `go test -v` bestehen unverändert; der Diff benennt keinen um.
   - (c) zwischen Doc-Kommentar und `TestE2E_SandboxDevcontainer_Ollama` steht nichts Eingeschobenes.
   - „its pass shows“ ist mit dem Bezug aus F-1 mehrdeutig, sonst verständlich.
3. **Alter Wortlaut im Baum:** `grep -rn "visible pass"` und `"not a verification"` selbst wiederholt. Treffer nur in `CHANGELOG.md` Zeile 190 (ausgelieferte Version, anderer Wortlaut), im Slice selbst, in [`slice-cr-lastenheft-messstaende`](../plan/planning/done/slice-cr-lastenheft-messstaende.md), [`slice-v1-devcontainer-ollama-amd64-messung`](../plan/planning/done/slice-v1-devcontainer-ollama-amd64-messung.md) und im Verifier-Report [`2026-10-25-slice-cr-lastenheft-messstaende-verify`](2026-10-25-slice-cr-lastenheft-messstaende-verify.md). Kein Treffer in Code, Spec, Handbuch, README, `harness/`. Die ADR zu Ollama im Sandbox-Devcontainer ist unverändert (`git diff` zeigt keine Änderung unter `docs/plan/adr/`).
4. **Build (Docker-only, Image `u-boot:test-docker-tools`):** `gofmt -l internal/e2e/` ohne Ausgabe; `go vet -tags docker ./internal/e2e/` ohne Befund; `go test -tags docker -run XXX -count=1 ./internal/e2e/` ergibt `ok … [no tests to run]`. Kein Docker-Testlauf, kein volles `make test-docker` (Kommentar-Diff; Begründung im Slice-DoD, hier geteilt).
5. **Plan-Tabelle gegen `git diff --stat`:** alle vier Diff-Dateien stehen in der Tabelle; Roadmap-Zeile ist als „update, bei Closure zurück auf ‚Nichts in Arbeit.‘“ geführt und nennt den Grund (`planning-drift`). Roadmap-Zwischenzeile „Ohne Welle in Arbeit: …“ trägt Zustand und Anker, keine Chronik; `make docs-check` im Gate-Lauf grün. Abweichungen: F-2, F-3.
6. **`make gates`:** Exit 0 (lint, test, coverage-gate, docs-check, baseline-verify, ablage-pruefung inklusive 26 Fälle ohne Fehlschlag).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Diff `internal/e2e/sandbox_devcontainer_docker_test.go` (nur `//`, keine Logik) | geprüft, ohne Befund |
| Kommentar gegen §3.9 (a), (b), (c) | geprüft, ohne Befund |
| Zusage- und Grenz-Teil des Kommentars gegen Testcode | geprüft, ohne Befund |
| `repo.mk`, `integration.yml` (Kopplung an `-v`) | geprüft, ohne Befund |
| ADR zu Ollama im Sandbox-Devcontainer, Entscheidung 10 (unverändert) | geprüft, ohne Befund |
| Alter Wortlaut im Baum | geprüft, ohne Befund |
| `roadmap.md` Zwischenzeile | geprüft, ohne Befund |
| Build mit Tag `docker` (gofmt, vet, Kompilat) | geprüft, ohne Befund |
| `make gates` | geprüft, grün |
| `make test-docker` (voll) | nicht gelaufen: Kommentar-Diff, rund 15 Minuten und 1,5 GB Download, pinnt nichts, was der Diff ändert |
| `make verify-depguard` | nicht gelaufen: keine Importregeln berührt |
| Messprotokoll unter `docs/messungen/` zum Lastenheft-Messstand | nur als Kontext gelesen; nicht Teil des Diffs |
| DoD-Konformität des Slice | nicht geprüft: Verifier-Aufgabe |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Kommentar-Verweiswort ohne Bezug im Text · Plan-Tabelle zeigt Vor-Diff-Zeilenbereich · Zustandsmarker nach Lifecycle-Move nicht nachgezogen

## Verdikt

**Merge-blockierend:** nein. Kein HIGH, kein MEDIUM; die drei LOW-Findings sind Wartbarkeit/Doku-Präzision mit Anker in `AGENTS.md` §3.7 bzw. §3.10.

**Übergabe:** F-1 bis F-3 an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Der Report ersetzt keine Verifikation.
