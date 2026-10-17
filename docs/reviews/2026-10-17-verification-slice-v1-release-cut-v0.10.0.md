# Verification Evidence: slice-v1-release-cut-v0.10.0 — 2026-10-17

**Rolle:** Verifier (frischer Kontext, kein Self-Review, repariert nichts). Frage: Wurde der Cut richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach [`harness/verification.md`](../../harness/verification.md).

**Gegenstand:** Slice [`slice-v1-release-cut-v0.10.0`](../plan/planning/done/slice-v1-release-cut-v0.10.0.md), Commit-Range `630f427..HEAD` (`6536a66`, sechs Commits, 14 Dateien einschließlich Review-Report). Maßstab: §2 (DoD), §3 (Plan-Tabelle, Maintainer-Aktionen M0–M9), §6 (Risiken), [`releasing.md`](../maintainer/releasing.md) §2–§7, Vorbild [`slice-v1-release-cut-v0.9.0`](../plan/planning/done/slice-v1-release-cut-v0.9.0.md). Review-Report: [`2026-10-17-slice-v1-release-cut-v0.10.0.md`](2026-10-17-slice-v1-release-cut-v0.10.0.md) (0 HIGH, 0 MEDIUM, 3 LOW, 6 INFO). Der Slice bleibt bis nach dem Tag in `in-progress/`; geprüft ist der Teil der DoD, der **vor dem Tag** erfüllbar ist.

## Verification Evidence

Scope:
- Slice: [`slice-v1-release-cut-v0.10.0`](../plan/planning/done/slice-v1-release-cut-v0.10.0.md)
- IDs: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache), [ADR-0007](../plan/adr/0007-distributionswege-ghcr.md), [ADR-0016](../plan/adr/0016-homebrew-distribution-per-tap.md), [ADR-0017](../plan/adr/0017-linux-pakete-nfpm.md), [ADR-0021](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md)
- Artefakte: `CHANGELOG.md`, `cmd/uboot/main.go`, `repo.mk`, `Dockerfile`, `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, `docs/maintainer/releasing.md`, `AGENTS.md`, Roadmap und Planungs-READMEs

DoD-Abgleich (§2 des Slice; Stand vor M1):
- [ ] `v0.10.0` getaggt, `publish`-Lauf grün, rc.1 geprobt und aufgeräumt — **offen, erwartet** (M4–M7 nach Freigabe). Voraussetzungen belegt: Secret `HOMEBREW_TAP_GITHUB_TOKEN` existiert (nur Existenz, Eintrag vom 2026-10-01; Gültigkeit nicht prüfbar); `publish.yml` leitet `is_stable` aus dem Tag ab (`*-*` → `false`, Log-Zeile `:85`), `tap` läuft nur bei `is_stable == 'true'` (`:297`), Formel (`u-boot.rb`), `:latest`-Push und Tap-Nachzug nur im `true`-Zweig (`:157`–`:163`, `:235`–`:240`), `--prerelease` bei `false` (`:182`): die Aussage des Slice „rc.1 beweist nur den `false`-Zweig“ stimmt. Die Formel-Füllung mit Beispiel-`SHA256SUMS` lokal geprüft (Scratchpad, ohne externe Wirkung): `scripts/homebrew-formula-fill.sh v0.10.0 …` schreibt `u-boot.rb` mit Tag und vier Digests, ohne `version`-Zeile. Tap-Repo (lesend): `smoke.yml` im Tap ist identisch zur Vorlage `packaging/homebrew-tap/.github/workflows/smoke.yml`; letzter Smoke-Lauf auf `5129ebc` (Formel 0.9.0 ohne `version`-Zeile) `success`.
- [x] `make gates`, `make ci`, `make test-docker`, `make fullbuild` vor dem Tag grün — selbst gefahren (Sensors unten), Exit 0 je Target. Der Teil „letzter CI-Lauf auf dem Tag-Commit grün“ ist erst nach M1/M2 prüfbar: **offen**.
- [x] Doku-Update (vor dem Tag erfüllbarer Teil) — CHANGELOG, READMEs, Handbuch tragen `0.10.0`; Messstand genannt; `--quiet`-Doppelung im Handbuch durch Verweis ersetzt (Einzelnachweise unten). **Offen bis nach dem Tag:** Roadmap-Zeile unter „Abgeschlossene Wellen“ (Tag-Hash, Datum) und `roadmap.md` „Letzte Änderung“ — bewusst erst nach M7 (§3, M9); bis dahin steht ein Zeiger unter „Offene Wellen“.
- [x] Review durchgeführt, Report liegt vor — [`2026-10-17-slice-v1-release-cut-v0.10.0.md`](2026-10-17-slice-v1-release-cut-v0.10.0.md); behoben am Text: F-1, F-2, F-3, F-5, F-9 (Nachweis unten); F-4, F-6, F-7, F-8 bewusst ohne Änderung.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen, §7 steht auf Platzhalter (Closure nach dem Tag, Planner).
- [ ] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen“ — offen (Closure). Hinweis: [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../plan/planning/observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) bleibt offen; siehe Befund V-3.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — offen (Closure); §6 nennt für jedes Risiko den vorgesehenen Ausgang oder „weiter offen“ mit Anker.
- [ ] Drei Paarungen (Anker · Folge-Slice · Register) — offen (Closure).

Einzelnachweise zur Doku-DoD (selbst gelesen und gegrept):
- CHANGELOG: `## [Unreleased]` leer, `## [0.10.0] - 2026-10-17`, Lead „Tenth release.“ nennt alle vier Inhalte (Features, `podman build` mit sechster Lockerung und Wirkung auf bestehende Projekte, `--quiet`-Vertrag, Docker-Desktop-Rücknahme, Lastenheft 0.12.0); Compare-Links `[Unreleased]: …v0.10.0...HEAD`, `[0.10.0]: …v0.9.0...v0.10.0`. Messstand: „Measured with Docker 29.5.2 on Colima 0.10.3 (arm64, Ubuntu 24.04 VM) … **Not measured:** Docker Desktop, amd64 and CI“ — und jetzt zusätzlich eigener Nachlauf (Sensors, Zeile `NestedPodman_Build`). Lastenheft-Version `0.12.0` (`spec/lastenheft.md:9`) stimmt mit dem CHANGELOG. Glättung der `doctor`-Warn-Einträge („credential source and egress restriction“) steht; Sprachvertrag-Zeile (Link in der Traceability-Tabelle) (englische Beschreibung, keine `version`-Zeile) steht und stimmt mit dem Formel-Template (`desc` englisch, keine `version`-Zeile).
- Versionsstrings: `git grep 0.10.0-dev` findet `cmd/uboot/main.go:43`, `repo.mk:51,54,55`, `Dockerfile:35,37,39` (drei Stellen, sieben Zeilen); `git grep 0.9.0-dev` außerhalb von Plänen und Reports: kein Treffer. Alte Werte `0.9.0` im Baum außerhalb der Historie: nur Releases-Zeilen `v0.9.0`, Roadmap-Zeile, [`MR-003`](../../harness/conventions/MR-003-roadmap-wellen.md), `carveouts.md` (Formel-Smoke 0.9.0, Befund der Vergangenheit), Testkommentar „Lastenheft 0.9.0“ in `quiet_modifying_test.go` (historische Spec-Version, kein Versionsstand des Produkts) — alles zulässig.
- Image-Version: `make fullbuild` baut das Runtime-Image mit dem Default; `docker run --rm u-boot:latest --version` → `u-boot version 0.10.0-dev`. `make build VERSION=0.10.0` (lokal, Docker, nichts gepusht): `--version` → `u-boot version 0.10.0`, OCI-Label `org.opencontainers.image.version=0.10.0`. Das ist dieselbe Kette wie im Tag-Build (`publish.yml`: `make build VERSION="$VERSION"`). Danach lokal wieder `make build` (Image `0.10.0-dev`).
- READMEs: `grep` nach `LH-`, `ADR-`, `SPEC-`, `ARC-`, `lastenheft.md#`, `spezifikation.md`, `docs/plan`, `slice-`, `welle-`, `roadmap` in `README.md`, `README.de.md`, `docs/user/`: kein Treffer. Status-Block nennt den Stand (aktuelle Version, Zusagen, Grenzen), kein Datum, kein „released“. Feature-Zahl 13 gegen `featureCatalogue()` (`devcontainer_features.go:175`; Schlüssel `git`, `docker-cli`, `node`, `java`, `go`, …, `sshd`, `claude-code`, vom Review gezählt) — Fließtext und Status-Block beider READMEs sagen 13. EN/DE-Spiegelung: Status-Block, Tabellenzeile `v0.10.0` (unmittelbar nach `v0.9.0`, Link `releases/tag/v0.10.0`) und `VERSION=0.10.0` inhaltsgleich. Links: `make docs-check` 0 Befunde (318 Dateien).
- Handbuch: `Software-Version: v0.10.0`, `Stand: 2026-10-17`, Beispielausgabe `u-boot version 0.10.0`, `download/v0.10.0`, `VERSION=0.10.0`, Gültigkeitsbereich `v0.10.0`; die `--quiet`-Zeile in „Ausführlichkeit der Ausgabe“ ist ein Verweis auf [Allgemeine Optionen](../user/benutzerhandbuch.md#allgemeine-optionen), die Langfassung dort ist in vier nummerierte Ausnahmen gegliedert; „Gemessen“-Absatz trennt `debian:stable-slim` (`RUN` mit Netz unter Egress) von `alpine:3` (Pull und Build unter Egress).
- `AGENTS.md` §5 Zeile 6: Wortlaut trägt Kennungsverbot, Verweis-Ausschluss (Roadmap, Slices, Wellen, ADRs, Spezifikation, Architektur-Sicht), die Ausnahme „Link ohne Anker auf `spec/lastenheft.md`“ und „Status = Stand“; deckt sich mit `releasing.md` §3 Nr. 3 und dem Befund der READMEs (nur Link auf `spec/lastenheft.md` ohne Anker, `CHANGELOG.md`).
- `releasing.md`: §3 (Verweis auf [`slice-v1-release-cut-v0.9.0`](../plan/planning/done/slice-v1-release-cut-v0.9.0.md)), §4 (`baseline-verify` im `make gates`-Kommentar, deckt sich mit [`harness/README.md`](../../harness/README.md) §Sensors), §3 Nr. 4 („Offene Wellen“), §5 (Datumsprüfung vor jedem Tag) und §7 (Datumsprüfung danach) zueinander konsistent.
- Roadmap: „Nichts in Arbeit.“ ersetzt durch Zeiger auf den Slice unter „Offene Wellen“ mit „(Tag-Push steht aus)“; Zeile in `in-progress/README.md` gesetzt, Zeile in `open/README.md` entfernt (Ruhe-Marker-Logik trägt, `docs-check` grün).
- Review-Findings am Text: F-1 (kein „Three tagged releases“, kein „10 architecture decision records“ mehr in beiden READMEs), F-2 (kein „Phase table“, „M1..M8“, „Nächste Schritte“, kein Planungs-Link in den READMEs; Verlauf verweist auf Tabelle und CHANGELOG), F-3 (Handbuch-Absatz getrennt, „in der VM“ einmal), F-5 (`releasing.md` §5 trägt die Vor-Tag-Datumsprüfung), F-9 (DE-Zeile `v0.10.0` heißt `„sshd, Claude Code und `podman build`“` ohne „Milestone“; ältere Zeilen behalten es, wie vom Review als INFO belassen).
- Plan-Tabelle §3 gegen `git diff --stat 630f427..HEAD`: alle 14 Dateien stehen in der Tabelle (13 Zeilen des Reviews plus Review-Report, der als „durch den Reviewer“ geführt ist; `docs/plan/planning/open/README.md`, `in-progress/README.md`, Roadmap, Slice selbst); kein „update“ ohne Diff; `packaging/`, `scripts/homebrew-formula*` und `observations/` ohne Diff, wie geführt.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0; Lint `0 issues`; `coverage-gate: OK — coverage 91.20% meets threshold 90%`; Stempel `.harness/state/gates-passed.head` = `6536a66` |
| `make ci` | pass | selbst gefahren, Exit 0; `baseline-verify: v6.18.1 OK — 54 Dateien`; `d-check: 318 Datei(en) geprüft, 0 Befund(e)`; `govulncheck`: `No vulnerabilities found.`; Trivy `u-boot:latest` Debian 12.15 und `usr/local/bin/u-boot`: 0 Funde; `[ci] gates + govulncheck + image-scan green` |
| `make doc-trace` | pass | selbst gefahren: `83 Anforderung(en), 0 Waise(n)` |
| `make docs-check` | pass | `d-check: 318 Datei(en) geprüft, 0 Befunde` (mit dieser Datei) |
| `make fullbuild` | pass | selbst gefahren, Exit 0; `[fullbuild] ci + runtime image green`; `--version` des gebauten Images siehe oben |
| `make test-docker` | pass | selbst gefahren (Colima, Schalter `kernel.apparmor_restrict_unprivileged_userns = 0`, per `colima ssh` gelesen), Exit 0; `internal/e2e` ok in 442 s, alle Pakete `ok` |
| Nachlauf `…NestedPodman_Build`, `…_BuildEgress` mit `-v -count=1` | pass, **nicht übersprungen** | `docker run … u-boot:test-docker-tools go test -tags docker ./internal/e2e/ -run 'TestE2E_SandboxDevcontainer_NestedPodman_Build' -v`: `--- PASS: …_Build (25.98s)`, `--- PASS: …_BuildEgress (69.21s)`, kein `SKIP`; Beleg gilt für Colima, arm64, nicht für amd64 oder CI |
| `make verify-depguard` | not run | keine Importregel berührt (Diff enthält eine Versionszeile in Go-Code) |
| `gh run list --branch main --limit 3` | lesend | `ci` (`39030308799`) und `integration-docker` (`39030308808`) auf `be7e1bb` (= `origin/main`) `success`; `git log origin/main..HEAD` = 36 Commits, also ist der Cut-Stand noch nie in CI gelaufen |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [ADR-0021](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) | CHANGELOG-Eintrag und Handbuch „Gemessen/Nicht gemessen“; Nachlauf der beiden Build-Tests (pass, nicht übersprungen) |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | Handbuch „Allgemeine Optionen“ (vier Ausnahmen) gegen Lastenheft 0.12.0; Review prüfte gegen Code (`generate.go`); Abweichung F-7 benannt |
| [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | 13 Features in READMEs gegen `featureCatalogue()` |
| [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache) | englische Formel-/Paketbeschreibung im CHANGELOG, Formel-Template gelesen |
| [ADR-0016](../plan/adr/0016-homebrew-distribution-per-tap.md), [ADR-0017](../plan/adr/0017-linux-pakete-nfpm.md), [ADR-0007](../plan/adr/0007-distributionswege-ghcr.md) | `publish.yml` gelesen; Formel-Füllung lokal; `make build VERSION=0.10.0` → `--version` und OCI-Label `0.10.0` |

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: alle (kein Gate, keine Schwelle, kein Workflow im Diff); `carveouts.md` nicht berührt.

Nicht ausgeführt:
- `git push`, `git tag`, `gh release`, `gh workflow run`, Schreiben in Fremd-Repos — Auftragsgrenze; M1–M8 liegen beim Projektinhaber.
- `make verify-depguard` — keine Importregel berührt.
- Die `true`-Zweige von `publish.yml` (Formel-Upload, `:latest`, Job `tap`) im echten Lauf — erst mit `v0.10.0` ausführbar (M7).
- amd64 und CI für die Build-Tests — nicht messbar ohne CI-Log mit `-v` (siehe V-3).
- Gültigkeit des Tap-Tokens — nur Existenz sichtbar.

Commit / Artefakt:
- Geprüft auf `6536a66`; lokale Images `u-boot:latest` (`0.10.0-dev`), kein Release-Artefakt.

## Befunde und Diskrepanzen

| ID | Klasse | Befund | Wirkung |
|---|---|---|---|
| V-1 | INFO | Der Plan nennt den Report `docs/reviews/<Datum>-release-cut-v0.10.0.md`; er liegt als [`2026-10-17-slice-v1-release-cut-v0.10.0.md`](2026-10-17-slice-v1-release-cut-v0.10.0.md). | keine; Platzhalter, kein Gegenstand |
| V-2 | INFO | DoD-Zeile 3 nennt „Roadmap trägt `v0.10.0`“: vor dem Tag steht nur der Zeiger unter „Offene Wellen“; die Zeile in „Abgeschlossene Wellen“ (Tag-Hash, Datum) kommt per Plan erst nach M7. | keine; Closure-Aufgabe, nicht vor dem Tag erfüllbar |
| V-3 | INFO | Der Aufruf `go test -tags docker ./...` (`repo.mk:132`) läuft ohne `-v`; im CI-Log ist ein Skip der Build-Tests nicht von einem Pass trennbar. Der lokale Nachlauf zeigt Pass auf Colima/arm64, belegt aber weder amd64 noch CI. Die Aussage „amd64 und CI nicht gemessen“ bleibt zutreffend; die Beobachtung [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../plan/planning/observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) bleibt offen (Entscheidung des Projektinhabers, keine Workflow-Änderung). | keine Änderung der Release-Aussage |
| V-4 | INFO | CHANGELOG und READMEs sagen „with and without the egress restriction“ für den Messstand ohne Aufteilung nach Image, das Handbuch trennt (`alpine:3` unter Egress nur Pull und Build). Beides ist durch Messprotokoll und Validator-Zusatzbeleg gedeckt (Review F-4); der Handbuch-Wortlaut ist enger als der CHANGELOG, nicht weiter. | keine |
| V-5 | INFO | README-Zeile „16 diagnostic checks“ (Subkommando-Referenz, `doctor`) ist nicht Teil des Diffs und wurde in diesem Lauf nicht gegen den Code gezählt. | keine; nicht verifiziert |
| V-6 | INFO | Vom Review offen gelassen und hier bestätigt offen: F-4, F-6 (Status-Block und Link `releases/tag/v0.10.0` behaupten bis M6 einen Stand ohne Tag), F-7 (Lastenheft nennt noch „Sicherheitslockerungen des Docker-Daemons“ als `doctor`-Warn-Eintrag), F-8 (Wortlaut `AGENTS.md` §5 Zeile 6), Hinweise L-4 und L-5 des Validators. | Entscheidungen beim Projektinhaber, keine Cut-Aufgabe |

Keine DoD-Verletzung vor dem Tag gefunden; keine Behauptung ohne Beleg in den geprüften Doku-Aussagen.

## Go/No-Go

**M1 (Push von `main`): GO** — vorbehaltlich der ausdrücklichen Freigabe des Projektinhabers im Gespräch. Gates, CI-Target-Lauf, Docker-Tests und Runtime-Image-Build sind grün; `make docs-check` 0 Befunde; Review ohne HIGH/MEDIUM, behobene Findings am Text bestätigt; Plan-Tabelle deckt den Diff.

**Auflagen vor M1:**
1. Diese Verification-Datei committen (sonst liegt der geprüfte Stand nicht im Push); danach `make docs-check` und `make gates` erneut (der Stempel gilt für `6536a66`).
2. M0: `date +%F` gegen das Datum in CHANGELOG, beiden Tabellenzeilen und Handbuch-`Stand` (heute `2026-10-17`, stimmt; bei Push an einem anderen Tag ein weiterer Cut-Commit).
3. Der Projektinhaber bestätigt in der Freigabe von M1, dass der Verbosity-Vertrag in dieser Version gewünscht ist (§6 des Slice).

**Auflagen vor rc.1 (M4):** M2 grün auf dem Push-Tip (`ci`, `integration-docker`; bei rotem `image-scan` Go-Pin heben, [`releasing.md`](../maintainer/releasing.md) §4); M3 CI-Log der Build-Tests sichten (ohne `-v` vermutlich nicht entscheidbar, dann nur festhalten); M0 erneut; rc.1-Erwartungen: `Resolved version=0.10.0-rc.1, is_stable=false`, `tap` übersprungen, Prerelease, kein `:latest`, keine `u-boot.rb`, `docker run … :0.10.0-rc.1 --version` = `0.10.0-rc.1`.

**Auflagen vor dem stabilen Tag (M6):** rc.1 aufgeräumt (M5, GHCR-Löschrecht des Tokens ungewiss, Risiko im Slice §6); der Tag auf dem Commit mit grünem CI; M0 erneut; keine weiteren Commits zwischen M2 und M6 ohne neuen CI-Lauf; danach M7 (zwölf Assets, `:latest` gleicher Digest, erster echter Lauf des `true`-Zweigs mit Job `tap`; Token nur auf Existenz geprüft) und M8; Roadmap-Zeile, Closure, Steering-Loop-Eintrag, Register und die drei Paarungen erst nach M7.
