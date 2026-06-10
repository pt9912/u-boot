# Slice M4: Soft-Existing-Detection für `u-boot init`

> **Status:** Done
> **DoD:** Commit `5415d9f`

## Auslöser

M3-T2 hatte nur die **Hard**-Marker-Erkennung aus [`LH-FA-INIT-004`](../../../../spec/lastenheft.md#lh-fa-init-004--bestehendes-projekt-erkennen)
implementiert (`u-boot.yaml`/`compose.yaml`/`.env.example` präsent
→ `ErrProjectExists`). Die **Soft**-Erkennung fehlte komplett — Spec
verlangt:

> *„Liegt keine Projektsteuerdatei vor, gilt das Verzeichnis nur als
> wahrscheinliches bestehendes Projekt, wenn mindestens drei Elemente
> aus dem Mindestumfang der Projektstruktur bereits vorhanden sind.
> In diesem Fall muss `u-boot init` im interaktiven Modus explizit
> nachfragen ..."*

Bewusste M3-Lücke (Carveout, [`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005--carveout-disziplin)); M3-T4c hatte das
`--assume-existing`-Flag bereits durch den CLI bis in die Request
durchgereicht („load-bearing when [slice-m4-soft-existing-detection](slice-m4-soft-existing-detection.md)
lands"). Dieser Slice macht es load-bearing.

## Aufhebung

Soft-Detection ist im `InitProjectService` aktiv. Sechs Indikatoren
(deterministische Reihenfolge):

1. `README.md`
2. `CHANGELOG.md`
3. `docs/`
4. `scripts/`
5. `docker/`
6. `.devcontainer/devcontainer.json`

Schwelle: **≥ 3** Indikatoren → „wahrscheinliches bestehendes
Projekt". Entscheidungs-Baum vor `planTemplatedFiles`:

- `--force` oder `--backup` gesetzt → Detection skippen (planFile
  übernimmt).
- `< 3` Indikatoren → skippen.
- `--assume-existing` → abort mit `ErrProjectExists` (Trigger:
  „--assume-existing").
- `--no-interactive` (ohne `--assume-existing`) → skippen (spec
  §247-Carveout); per-File-Kollision in planFile bleibt aktiv.
- sonst → `driven.Confirmer.ConfirmTreatAsExisting`. „Yes" → abort
  mit `ErrProjectExists` (Trigger: „user confirmation"); „No" →
  skippen.

## Geliefert

- **Driving-Port-Erweiterung**: `driving.InitProjectRequest` neue
  Bool-Feld `NoInteractive` (mit Doc-Comment zum spec-§247-Carveout).
  `AssumeExisting` ist jetzt load-bearing (Doc-Update entfernt das
  „M3 no-op"-Hinweis).
- **Driven-Port** `port/driven.Confirmer` mit `ConfirmTreatAsExisting(
  ctx, baseDir, indicators) (bool, error)`. Narrow-scoped (eine
  Methode pro Confirm-Kontext, statt generisches `Confirm(prompt)`)
  — zukünftige Bestätigungsfälle bekommen eigene Methoden.
- **Application-Service**:
  - Konstruktor um `confirmer driven.Confirmer` erweitert; `nil`
    routed auf internen `noopConfirmer` (refuses, modelliert
    deterministisches Non-Interactive).
  - `checkSoftExisting(ctx, req)` als erster Schritt nach `baseExists`-
    Check.
  - Helpers `softIndicators()`, `detectSoftIndicators(baseDir)`,
    `softExistingThreshold = 3`, `softExistingAbort(indicators, trigger)`.
- **Driven-Adapter** `internal/adapter/driven/confirm/Confirmer`:
  `bufio.Scanner` über `io.Reader`, Prompt auf `io.Writer`. Default
  `[y/N]` (sicherer Default: no = proceed). EOF → no.
- **Driving-Adapter (CLI)**:
  - `--no-interactive` propagiert in `req.NoInteractive`.
  - `--assume-existing`-Flag-Doc auf „assert existing project; aborts
    unless --backup/--force" aktualisiert; M3-stderr-Note entfernt.
  - `Long`-Doc um die Soft-Detection-Sektion erweitert.
- **cmd/uboot-Wiring**: `confirm.New(os.Stdin, stderr)` konstruiert
  und in `application.NewInitProjectService` injiziert.
- **Tests**:
  - Service: 7 neue Tests (under-3, AssumeExisting, NoInteractive,
    Confirmer-Yes, Confirmer-No, Force/Backup-Skip, Confirmer-Error)
    plus `fakeConfirmer` in `fakes_test.go`. `seedSoftIndicators`-
    Helper mit collision-safer Reihenfolge (Dirs zuerst, Template-
    Files zuletzt).
  - CLI: `TestExecute_InitAssumeExisting_NoLongerEmitsM3Note`
    (Regression-Guard gegen Re-Add der Stderr-Note),
    `TestExecute_NoInteractive_PassThrough` (Flag-Propagation),
    bestehendes `TestExecute_NoAssumeExisting_NoStderrNote`
    weiterhin aktiv.
  - Adapter: 4 Tests für `Confirmer` (yes-variants, no-variants
    inkl. EOF, Indicator-Display, Read-Error).

## Out of Scope

- **`--yes`-Verhalten** für die Soft-Detection-Prompt: heute
  unspezifiziert. `--yes` ([`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)) bedeutet „answer the
  default" — für `[y/N]` wäre das `N` (proceed). Eine saubere
  Implementierung bräuchte eine Konstruktions-Zeit-Aware Confirmer-
  Variante; das Wiring ist heute Service-vor-Flag-Parsing
  (Chicken-and-Egg). Behoben in einem späteren Slice, sobald weitere
  Confirm-Prompts hinzukommen und das Pattern den Aufwand
  rechtfertigt. Heutiger Workaround: User passt `--no-interactive`
  oder `--assume-existing` zusätzlich, um die Prompt zu umgehen.
- **Devcontainer-Erkennung als Soft-Marker via semantischen Parse**:
  nur Pfad-Existenz wird geprüft (`.devcontainer/devcontainer.json`).
- **Per-Indicator-Gewichtung** (z. B. `docker/` zählt mehr als
  `README.md`): out of scope; einfache Anzahl ≥ 3 reicht für MVP.

## Bezug

- Auslösende Spec: [`LH-FA-INIT-004`](../../../../spec/lastenheft.md#lh-fa-init-004--bestehendes-projekt-erkennen), [`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung).
- Aufhebung dokumentiert in: [`carveouts.md`](../in-progress/carveouts.md)
  (Zeile entfernt) und [`roadmap.md`](../in-progress/roadmap.md)
  (Carveout-Auflösungs-Slice-Tabelle).
- Hängt von: M3-T3 (CLI-Adapter), M3-T4c (`--assume-existing`-
  Durchreichung).
