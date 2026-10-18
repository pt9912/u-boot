# Verification-Report: slice-v1-cli-dry-run-log-would — 2026-10-18 <!-- d-check:ignore (Kennung im Titel) -->

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan und DoD; kein Review, keine Validierung).

## Scope

- Slice: [`slice-v1-cli-dry-run-log-would`](../plan/planning/done/slice-v1-cli-dry-run-log-would.md) (in-progress), Commits der Range `fa019d0..HEAD` ohne `044e0c9` und ohne die `open/`-Änderungen des Vorgangs [`slice-v1-cli-help-handbuch-link`](../plan/planning/open/slice-v1-cli-help-handbuch-link.md): `1f42dfd`, `a79bf59`, `64e853a`, `878e453`, `b35ad74`, `aacee08`, `bf2d592`, `c94eeb6`.
- IDs: [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging).
- Eingang: DoD-Bestätigung des Implementers; Review [`2026-10-18-slice-v1-cli-dry-run-log-would.md`](2026-10-18-slice-v1-cli-dry-run-log-would.md) (0 HIGH/MEDIUM, F-1/F-2 LOW, F-3 bis F-5 INFO).
- Host: Colima, kein Host-Go.

## Sensors (selbst gefahren)

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` | pass (Exit 0) | `coverage-gate: OK — coverage 91.10% meets threshold 90%`; lint, test, docs-check, baseline-verify im Lauf grün |
| `make doc-trace` | pass | `83 Anforderung(en), 0 Waise(n).` |
| `make doc-immutable RANGE="fa019d0..HEAD"` | pass | `9 Commit(s) — OK`, `0 Befund(e)` |
| `make docs-check` | pass | nach Ablage dieser Datei: 0 Befunde (siehe Übergabe) |
| `make test-docker` | nicht gelaufen | `git diff --stat fa019d0..HEAD`: kein Template, kein e2e, kein Port, kein Adapter im Diff (nur `internal/hexagon/application`, Doku) |
| `make verify-depguard` | nicht gelaufen | keine Importregel berührt (`dryrun_log.go` importiert nur `port/driven`, `port/driving`) |
| `make govulncheck`, `make image-scan` | nicht gelaufen | kein Security-/Release-Pfad |

## Praxislauf (Runtime-Images aus `git archive fa019d0` und `HEAD`, `docker build --target runtime`)

Eigene Tags `verify-dryrunlog:old|new`, Verzeichnis `~/uboot-verify-dryrunlog`; danach `docker rmi` beider Tags und Verzeichnis gelöscht, kein prune.

| Aufruf | vorher (`fa019d0`) | nachher (HEAD) |
|---|---|---|
| `config set project.name renamed --dry-run` | `config set: updated` path old new | `config set: would update`, Felder `path`, `old=demo`, `new=renamed` |
| `config set devcontainer.sandbox.egress.allow a.example.com --dry-run` | `config set: updated` | `config set: would update`, Felder `path`, `old=""`, `new` |
| `generate devcontainer --dry-run` | `generate devcontainer: created` | `generate devcontainer: would create`, Felder `project`, `forwardPorts`, `features`, `changed` gleich |
| `generate changelog --dry-run` | `generate: created` | `generate: would create`, Felder `artifact`, `path`, `project` |
| `generate changelog --diff` allein | `generate: created` | `generate: created` (unverändert) |
| `config set project.name r1 --diff` allein | `config set: updated` | `config set: updated` (unverändert) |
| ohne `--dry-run` (`config set`, `generate devcontainer`) | `updated` / `created` | `updated` / `created` (unverändert) |
| `generate changelog --dry-run --quiet` | keine Zeile | keine Zeile |
| `u-boot.yaml` nach den Dry-Runs | Prüfsumme gleich | Prüfsumme gleich; kein weiteres Artefakt im Verzeichnis |

## DoD-Abgleich (§2, Zeile für Zeile)

- [x] **Dry-Run ohne Info-Zeile in Vergangenheitsform (Code)** — `dryrun_log.go` `logWrite` (Kriterium genau `mode == driving.PreviewDryRun`), `devcontainerWriteMessage`. Alle neun Fundstellen der Tabelle umgestellt: `config.go` 2x (`Set`, `setFeatureSourcesAllow`), `devcontainer_egress.go` `setEgressAllow`, `config_many.go` `SetMany`, `generate.go` 4x `logWrite` (Zeilen 432, 487, 647, 706), `generateDevcontainer` (823), `applyAllowExternalFeatureSources` (1220, neuer Parameter `mode`, am Aufrufer in `writeDevcontainerPhase` übergeben). Restsuche `.Info(`/`.Warn(` in der Anwendungsschicht (ohne Tests): übrig nur Zustands-/Prüfzeilen (`maybeWarnOrphanFeatureActivation`, `generate_rollback.go`, `doctor.go`), wie im Plan begründet. Tests: `TestDryRunLog_WriteLines` (11 Fälle: Happy Dry-Run mit `would`, keine Vergangenheitswörter auf `INFO`; Boundary `PreviewNone` und `PreviewAndApply` Vergangenheitsform und kein `would`; Felder beider Zeilen als Schlüssel-/Wertfolge gepinnt, darunter `old`, `changed`, `items`, `added`/`total`), `…GenerateDevcontainerUpdatedBlock`, `…DryRunWritesNothingAndNoOpStaysQuiet` (Negative: Datei unverändert, No-op und Fehler ohne `INFO`). Die `--json`-Ausgabe liegt im CLI-Adapter, der im Diff nicht berührt ist (kein eigener Test in diesem Slice; Praxislauf unverändert). Erfüllt.
- [x] **Doku-Update** — Handbuch `### Vorschau statt Überraschung`: Absatz mit Beispielzeilen, `--diff` allein, `--quiet`; keine Kennungen; stimmt mit dem Praxislauf überein. CHANGELOG `[Unreleased]` · `### Changed` nennt alle acht Zeilen, Skript-Hinweis, `--json` als Maschinenvertrag; die acht Zeilen stimmen mit den acht `would`-Texten im Code überein. Erfüllt.
- [x] **`make gates` grün; `make doc-trace` ohne neue Waise** — siehe Sensors. Erfüllt.
- [x] **Review durchgeführt** — Report liegt vor. Findings am Code geprüft: F-1 behoben (`TestDryRunLog_WriteLines` vergleicht `fmt.Sprint(e.Args)` je Fall, auch `old` im Egress-Fall mit vorhandenem Wert, `changed` und `items`); F-2 behoben (Doc-Kommentar nennt die geltende Zusage: „sonst `<Aktion>`“); F-3 behoben (Plan nennt `TestConfigSet_SilenceLoggerFalse_EmitsInfo` ohne Zeilenhinweis; Symbol steht in `config_t3_test.go`); F-4 behoben (Doc-Kommentar von `applyAllowExternalFeatureSources` nennt `mode`); F-5 behoben (Zweig `would change` fehlt im Code, Suche im Baum ohne Treffer; Doc-Kommentar nennt die zwei Aktionen). Erfüllt.
- [ ] **Closure-Notiz (§7)** — offen: §7 steht auf „Wird bei der Closure gefüllt.“ (Closure-Schritt).
- [ ] **Register-Ausgang** — offen: `observations/BEO-CLI/dry-run-log-vergangenheitsform/state.md` steht auf „offen — ein Beleg“ (Closure-Schritt).
- [ ] **Risiken §6 mit Ausgang** — offen (Closure-Schritt); Vorschlag unten.
- [ ] **Drei Paarungen** — offen bis zur Closure.

Die DoD-Punkte 1 bis 4 sind erfüllt; 5 bis 8 sind Closure-Arbeit und nach dem Plan vor dem Move zu leisten.

## Plan-Tabelle §3 gegen `git diff --stat` der Slice-Commits (AGENTS.md §3.10)

Diff-Dateien: `CHANGELOG.md`, `roadmap.md`, die Slice-Datei, `docs/reviews/2026-10-18-slice-v1-cli-dry-run-log-would.md`, `benutzerhandbuch.md`, `config.go`, `config_many.go`, `devcontainer_egress.go`, `dryrun_log.go`, `dryrun_log_test.go`, `generate.go` (elf). Jede hat eine Plan-Zeile; jede Plan-Zeile (außer „unverändert“ und `state.md` als Closure-Eintrag) hat einen Diff. „Unverändert“-Aussagen (`generate_rollback.go`, `doctor.go`, `config_t3_test.go`, `logger.go`): keine Änderung im Diff, Aussage am Code gelesen (Test läuft ohne Vorschau-Modus und erwartet die Vergangenheitsform). Nicht im Diff und nicht im Plan: nichts Fehlendes. <!-- d-check:ignore (Kennungen als Zitat im Report) -->

## Diskrepanzen

- INFO: Zeilenhinweise in Plan §3 sind teils verschoben: `generate.go` Allowlist-Zeile steht bei `:1220` (Plan `:1218`), die Signatur bei `:1172` (Plan `:1170`); der Plan sagt selbst, Zeilennummern nennten den Fundort beim Schreiben. Symbole stimmen. Kein Handlungsbedarf, eine Zeilenangabe in der Closure-Notiz sollte entfallen oder aktualisiert werden.
- INFO: Plan §3 listet `docs/plan/planning/open/README.md` (Zeile für den Vorgang `slice-v1-cli-help-handbuch-link`) nicht; sie gehört dem fremden Vorgang `044e0c9` und ist nicht Gegenstand. <!-- d-check:ignore (Kennungen als Zitat im Report) -->
- Kein Fund: Hilfs-Image `verify-build-ctx:latest` im Docker-Bestand ist vier Wochen alt und nicht von diesem Lauf.

## Carveouts

Neu: none. Gelöst: none. Unverändert: alle.

## Auflagen für die Closure

1. **§7 füllen** mit Steering-Loop-Lerneintrag (Finding-Klassen aus dem Review: Test pinnt Zusage „Felder unverändert“ nicht · Kommentar beschreibt Vorzustand · Zeilenhinweis im Plan ungenau · Doc-Kommentar kennt neuen Parameter nicht · toter Rückfallzweig ohne Test; Zähler und Schwelle nach den Register-Regeln, keine Verkörperung ohne `liegt in`).
2. **Register `BEO-CLI/dry-run-log-vergangenheitsform`:** Vorschlag Ausgang `gestrichen` (nicht `verkörpert`: Zähler 1, unter der Schwelle von drei, und es gibt keinen Zielort). Begründung für `state.md`: die Ursache fiel weg, der Dry-Run meldet auf `INFO` nur noch den Konjunktiv (acht Zeilen über `logWrite` und `devcontainerWriteMessage`), Test `TestDryRunLog_WriteLines` pinnt ihn; Beleg `evidence/slice-v1-cli-dry-run-would-formulierung.md` bleibt, ein Beleg des eigenen Slice entfällt (der Slice beantwortet die Beobachtung, er zählt sie nicht). Das Verzeichnis bleibt liegen. <!-- d-check:ignore (Kennungen als Zitat im Report) -->
3. **Risiken §6:** (a) „Weitere Fundstellen“ — entfallen (Zeilenlektüre der Logger-Aufrufe, Review-Vollständigkeitsprüfung, Restsuche des Verifiers); (b) „Wortlaut als Vertrag“ — eingetreten als sichtbare Änderung, behandelt (CHANGELOG nennt sie, `--json` als Vertrag); (c) „`--diff` allein und Dry-Run verwechselt“ — entfallen (Boundary-Test, Praxislauf).
4. **Roadmap-Marker nach dem Move:** den Satz „In Arbeit: …“ in `in-progress/roadmap.md` zurück auf „Nichts in Arbeit.“ setzen (d-check `planning`); Move als reiner `git mv`-Commit nach Inhaltsänderung (AGENTS.md §3.3, Fall 2), `make slice-mv` zieht Verweise nach. <!-- d-check:ignore (Kennungen als Zitat im Report) -->
5. **`open/README.md`:** die Zeile von `slice-v1-cli-dry-run-log-would` zeigt auf `../in-progress/…`; nach dem Move auf `../done/…` und „geliefert (liegt in `done/`)“ ergänzen, wie bei den Nachbarzeilen. Die Zeile von `slice-v1-cli-help-handbuch-link` gehört dem anderen Vorgang und bleibt unberührt. <!-- d-check:ignore (Kennungen als Zitat im Report) -->
6. **DoD-Häkchen** in §2 setzen, wenn die Closure-Punkte stehen; danach `make gates` und `make docs-check` erneut.

## Übergabe

Verifier-Ergebnis: Bau entspricht Plan und DoD (Code, Doku, Gates, Praxislauf); offen sind nur die Closure-Punkte. Hinzugefügt wurde nur diese Datei.
