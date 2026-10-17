# Welle welle-cli-verbosity: Verbosity-Vertrag (Lastenheft, Code und Spezifikation)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/`. Der Zustand ist die Verzeichnis-Position — kein Status-Feld.
**Geplante Wellen bekommen noch keine Datei:** Sie stehen in der Roadmap unter
*Nächste Wellen* und nirgends sonst — zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-10-16.

---

## 1. Welle-Ziel

[`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), der Code in `internal/adapter/driving/cli/` und
die Spezifikation sagen dasselbe über die Verbosity: `--quiet` zeigt nur Fehler (die Erfolgs- und Zusammenfassungsmeldung aller
Subkommandos entfällt), mit vier Ausnahmen — die Ausweisung von Sicherheitslockerungen im Sandbox-Profil
([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)), die Ausweisung der Degradations-Fallbacks
([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)), der Compose-Fortschritt von `up`
([`LH-NFA-PERF-002`](../../../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker)) und die Pfad-Zusammenfassung von `--force`
([`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)) bleiben sichtbar —, die Ausgabe der Vorschau-Optionen `--dry-run` und `--diff` ist das Ergebnis des Aufrufs und bleibt (wie die der Abfragebefehle; [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe)), und die zuletzt
angegebene Option gewinnt. Ein Test pro Aussage belegt es; die Spezifikation beschreibt das Verhalten erst, wenn Lastenheft und
Code es tragen.

## 2. Trigger (Welle startet)

- Die Welle [`welle-harness-wartung`](welle-harness-wartung.md) liegt in `done/` (geschlossen 2026-10-16).
- Die Entscheidungen des Projektinhabers liegen vor (2026-10-16): `--quiet` zeigt nur Fehler; die Ausnahmen sind die
  Sandbox-Lockerungs-Hinweise (CR [`slice-cr-cli-quiet-fehler-und-warnungen`](slice-cr-cli-quiet-fehler-und-warnungen.md))
  sowie Degradations-Fallbacks, Compose-Fortschritt von `up` und `--force`-Zusammenfassung (CR
  [`slice-cr-cli-quiet-weitere-ausnahmen`](slice-cr-cli-quiet-weitere-ausnahmen.md); angenommene Fassung von
  [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)); der Trigger ist damit eingetreten.

## 3. Closure-Trigger (Welle schließt)

Das *Mehr* gegenüber den sechs Slice-DoDs: Die drei Stellen, an denen der Vertrag steht, sagen dasselbe. Die drei CRs ändern das
Lastenheft, die drei Code-Slices den Adapter, und die Spezifikation beschreibt erst danach das Verhalten.

- Alle sechs Slices der Welle liegen in `done/` (die zwei geschlossenen CRs, der CR
  [`slice-cr-cli-quiet-vorschau`](slice-cr-cli-quiet-vorschau.md), [`slice-v1-cli-quiet-nur-fehler`](slice-v1-cli-quiet-nur-fehler.md),
  [`slice-v1-cli-quiet-vorschau`](slice-v1-cli-quiet-vorschau.md) und [`slice-v1-cli-verbosity-letzte-option-gewinnt`](slice-v1-cli-verbosity-letzte-option-gewinnt.md)); das Lastenheft steht auf 0.10.0.
- Die Spezifikation trägt einen Eintrag zur Pegelabbildung und zur Rangfolge, der Lastenheft (0.10.0) und Code deckungsgleich
  beschreibt; er entstand nach dem Abschluss aller sechs Slices (Commit-Reihenfolge belegt es).
- `make doc-trace` meldet keine Waise unter [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging).
- `make gates` ist grün.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-cr-cli-quiet-fehler-und-warnungen`](slice-cr-cli-quiet-fehler-und-warnungen.md) | CR: Ausnahme der Sicherheitslockerungen unter `--quiet` (zuerst, geschlossen) | [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) |
| [`slice-cr-cli-quiet-weitere-ausnahmen`](slice-cr-cli-quiet-weitere-ausnahmen.md) | CR: Ausnahmen Degradations-Fallbacks, Compose-Fortschritt und `--force`-Zusammenfassung unter `--quiet` (Lastenheft 0.9.0) | [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-NFA-PERF-002`](../../../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker), [`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) |
| [`slice-v1-cli-quiet-nur-fehler`](slice-v1-cli-quiet-nur-fehler.md) | Code: `--quiet` zeigt nur Fehler in allen Subkommandos, vier Ausnahmen bleiben (Lastenheft 0.9.0; geschlossen) | [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) |
| [`slice-cr-cli-quiet-vorschau`](slice-cr-cli-quiet-vorschau.md) | CR: die Vorschau von `--dry-run` und `--diff` bleibt unter `--quiet` (Lastenheft 0.10.0) | [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) |
| [`slice-v1-cli-quiet-vorschau`](slice-v1-cli-quiet-vorschau.md) | Code: `--dry-run --quiet` und `--diff --quiet` geben die Vorschau aus (nach dem CR, Lastenheft 0.10.0) | [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) |
| [`slice-v1-cli-verbosity-letzte-option-gewinnt`](slice-v1-cli-verbosity-letzte-option-gewinnt.md) | Code: die zuletzt angegebene Option gewinnt, danach der Spezifikations-Eintrag (letzter Slice) | [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |

## 5. Abhängigkeiten

- Blockiert: keine Welle.
- Wird blockiert von: nichts; die Entscheidungen des Projektinhabers (Wirkung von `--quiet`, siehe §2) liegen vor. Innerhalb der
  Welle geht der CR [`slice-cr-cli-quiet-fehler-und-warnungen`](slice-cr-cli-quiet-fehler-und-warnungen.md) dem CR
  [`slice-cr-cli-quiet-weitere-ausnahmen`](slice-cr-cli-quiet-weitere-ausnahmen.md) voraus; dieser geht dem Slice
  [`slice-v1-cli-quiet-nur-fehler`](slice-v1-cli-quiet-nur-fehler.md) voraus, weil die Ausnahmen erst nach Lastenheft 0.9.0
  im Vertrag stehen. Der CR [`slice-cr-cli-quiet-vorschau`](slice-cr-cli-quiet-vorschau.md) folgt auf diesen Code-Slice und geht dem Code-Slice
  [`slice-v1-cli-quiet-vorschau`](slice-v1-cli-quiet-vorschau.md) voraus, weil der Vorschau-Satz erst nach Lastenheft 0.10.0 im Vertrag steht; der Code-Slice
  beginnt erst nach der Closure von [`slice-v1-cli-quiet-nur-fehler`](slice-v1-cli-quiet-nur-fehler.md), weil beide dieselben Quiet-Zweige berühren.
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](slice-v1-cli-verbosity-letzte-option-gewinnt.md) kommt zuletzt (er trägt den Spezifikations-Eintrag nach allen anderen); logisch ist er von den Vorschau-Slices unabhängig,
  berührt aber ebenfalls `root.go`.

## 6. Out-of-Scope für diese Welle

- Eine neue Pegelstufe zwischen Info und Debug: `--verbose` und `--debug` bleiben beide Debug; der Kommentar in `root.go` nennt das als vorläufig.
- Konfigurations- oder Umgebungsquellen für den Pegel: es gibt nur die drei Flags.
- Die Sandbox-Hinweise bei `init` und `generate` auf den Logger umzustellen: ein anderer Vorgang; der CR hält den Befund fest, entscheidet aber nur über `--quiet`.
- Die Homebrew-Smoke-Härtung ([`slice-v2-homebrew-smoke-haertung`](../done/slice-v2-homebrew-smoke-haertung.md)) und die Docker-Desktop-Messung: ohne Bezug zur Verbosity.

## 7. Closure-Notiz

Detailquellen: die Closure-Notizen der sechs Slices
([`slice-cr-cli-quiet-fehler-und-warnungen`](slice-cr-cli-quiet-fehler-und-warnungen.md) §7,
[`slice-cr-cli-quiet-weitere-ausnahmen`](slice-cr-cli-quiet-weitere-ausnahmen.md) §7,
[`slice-v1-cli-quiet-nur-fehler`](slice-v1-cli-quiet-nur-fehler.md) §7,
[`slice-cr-cli-quiet-vorschau`](slice-cr-cli-quiet-vorschau.md) §7,
[`slice-v1-cli-quiet-vorschau`](slice-v1-cli-quiet-vorschau.md) §7,
[`slice-v1-cli-verbosity-letzte-option-gewinnt`](slice-v1-cli-verbosity-letzte-option-gewinnt.md) §7); die Welle ist kein Release,
einen Release-Cut-Slice gibt es nicht (der CHANGELOG ist vollständig, ein Cut ist nicht mehr blockiert).

- **Geliefert:** Das Lastenheft ging von 0.7.0 auf 0.10.0; 0.7.0 war der Stand der vorigen Welle
  ([`welle-macos-verifikation`](welle-macos-verifikation.md), Zusage für Docker Desktop zurückgenommen), in dieser Welle entstanden 0.8.0
  (CR 1: Ausnahme der Sicherheitslockerungen), 0.9.0 (CR 2: vier Ausnahmen) und 0.10.0 (CR 3: Vorschau bleibt); Historie in
  [`spec/lastenheft.md`](../../../../spec/lastenheft.md) §16. Im Code zeigt `--quiet` in allen ändernden Subkommandos nur Fehler, vier
  Ausnahmen bleiben sichtbar (Sicherheitslockerungen, Degradations-Fallbacks, Compose-Fortschritt von `up`, `--force`-Pfad-Zusammenfassung),
  `doctor` listet nur Fehler-Einträge, `down` puffert seinen Strom, die Ausgabe von `--dry-run` und `--diff` bleibt, und die zuletzt
  angegebene Option gewinnt. [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)
  beschreibt Pegelabbildung, Rangfolge und Wirkung deckungsgleich mit Lastenheft 0.10.0 und Code.
- **Closure-Trigger (§3):** Alle sechs Slices liegen in `done/`. [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) entstand nach dem Abschluss der fünf anderen Slices
  (`git log`: erst die Moves der fünf, dann der Spezifikations-Commit, dann der Move des sechsten). `make doc-trace`: 83 Anforderungen, 0 Waisen,
  keine unter [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging). `make gates` grün.
- **Was hat funktioniert:** Die Vorab-Suche gegen ältere Zusagen im Lastenheft vor einer Änderung (Lehre aus dem zweiten CR) fand beim dritten CR
  nur noch [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) und
  [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe); der Review fand danach keine Vertragskollision mehr. Der Review
  gegen den Code lief vor der Spezifikation: [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) kam als letzter Slice und stimmt Aussage für Aussage mit Code und Lastenheft überein.
- **Was ging anders als geplant:** Die Welle wuchs von zwei auf sechs Slices: Zwei CRs wurden umgekehrt oder erweitert (erst „Fehler und
  Warnungen“, dann „nur Fehler“ mit einer, dann mit vier Ausnahmen), die Entscheidung zu `--quiet` fiel in drei Schritten (Sicherheitslockerungen,
  weitere Ausnahmen, Vorschau), und der Code-Slice `slice-v1-cli-quiet-nur-fehler` brauchte drei Review-Runden (zwei Reports und der Verifier); die Vorschau und
  die Rangfolge schnitten sich als eigene Slices ab. Docker- und Mac-Ressourcen waren nicht beteiligt.
- **Trigger-Audit:** Carveout: keiner berührt, [`carveouts.md`](../in-progress/carveouts.md) unverändert; 0 offen. Bootstrap-aware Gate:
  `coverage-gate` 91,10 % gegen 90 %; 0 offen. ADR: keine ADR trägt einen eingetretenen Re-Evaluierungs-Trigger (eine Suche nach `--quiet` und
  Verbosity in [`docs/plan/adr/`](../../adr/README.md) trifft keine ADR; [`ADR-0003`](../../adr/0003-solid-nahes-lint-profil.md) nennt `log/slog` nur als Linter-Regel); 0 offen. Hard Rule: keine hat
  einen eingetretenen Auflösungs-Trigger; die neue Regel `AGENTS.md` §3.9 ist permanent.
- **Steering-Loop-Einträge:**
  - **`AGENTS.md`** ergänzt: Pre-Commit-Selbstcheck des Implementers auf Kommentar-Muster (abwesender Text, Verweis auf Geändertes, Helfer zwischen
    Doc-Kommentar und Funktion) — liegt in `AGENTS.md §3.9`. Auslöser:
    [`BEO-CLI/kommentar-beschreibt-abwesendes`](../observations/BEO-CLI/kommentar-beschreibt-abwesendes/observation.md)
    (`slice-v1-cli-quiet-nur-fehler`, `slice-v1-cli-quiet-vorschau`, `slice-v1-cli-verbosity-letzte-option-gewinnt` — 3×; vier Review-Befunde in drei Vorgängen).
    Herkunft: `seit welle-cli-verbosity`. Entscheidung des Projektinhabers 2026-10-16: Die Regel stand schon in `AGENTS.md` §3.7 und in der HIGH-Liste des
    Reviewers; die Implementer verletzten sie trotzdem, deshalb steht jetzt der Selbstcheck vor dem Commit. Prosa-Form: Erreicht dieselbe Klasse ein
    **viertes** Mal die Schwelle, benennt der nächste Eintrag einen mechanischen Sensor oder begründet, warum keiner möglich ist.
  - Gezählt, nicht verkörpert (aus den Slice-§7, ohne `liegt in`): ein CR, dessen Begründung auf einer Annahme über den Code ruht, prüft die Annahme im
    Code vor der Änderung (einmal); eine Ausnahme-Regel im Lastenheft wird vor der Änderung gegen alle älteren Zusagen geprüft (zweimal, mit
    gegenteiligem Ergebnis); Plan-Aussagen über Code (Quellen, Dateien) werden vor dem Schreiben im Code geprüft (`slice-v1-cli-quiet-nur-fehler`:
    `doctor`-Items, die es nicht gibt; `slice-v1-cli-quiet-vorschau` §3: falsche Testdatei; `slice-v1-cli-verbosity-letzte-option-gewinnt` §3: unvollständige Dateiliste).
    Für die dritte Klasse ist keine Beobachtung angelegt; die nächste Eröffnung sichtet sie.
- **Beobachtungs-Register:** [`BEO-CLI/kommentar-beschreibt-abwesendes`](../observations/BEO-CLI/kommentar-beschreibt-abwesendes/observation.md) angelegt,
  drei Belege (`evidence/`), Ausgang `verkörpert` ([`state.md`](../observations/BEO-CLI/kommentar-beschreibt-abwesendes/state.md)). Zeiger auf das Register:
  [`../observations/`](../observations/README.md).
- **Folge-Slices:** [`slice-v1-sandbox-fallback-ausweisung`](../done/slice-v1-sandbox-fallback-ausweisung.md) (Fallback-Ausweisung, Entscheidung Code oder CR beim Projektinhaber),
  [`slice-v1-cli-dry-run-would-formulierung`](../done/slice-v1-cli-dry-run-would-formulierung.md) (Zeitform der Dry-Run-Zusammenfassung),
  [`slice-v2-homebrew-smoke-haertung`](../done/slice-v2-homebrew-smoke-haertung.md) (aus der Vorwelle, unverändert offen),
  [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md) (ohne Welle, wartet auf einen Docker-Desktop-Host).
  Alle vier liegen in `open/`.
- **Verifikation:** `make gates` Exit 0 (Coverage 91,10 %, `docs-check` 0 Befunde); Verification Evidence je Slice unter `docs/reviews/`:
  [`slice-cr-cli-quiet-fehler-und-warnungen`](../../../reviews/2026-10-16-verification-slice-cr-cli-quiet-fehler-und-warnungen.md),
  [`slice-cr-cli-quiet-weitere-ausnahmen`](../../../reviews/2026-10-16-verification-slice-cr-cli-quiet-weitere-ausnahmen.md),
  [`slice-v1-cli-quiet-nur-fehler`](../../../reviews/2026-10-16-verification-slice-v1-cli-quiet-nur-fehler.md),
  [`slice-cr-cli-quiet-vorschau`](../../../reviews/2026-10-16-verification-slice-cr-cli-quiet-vorschau.md),
  [`slice-v1-cli-quiet-vorschau`](../../../reviews/2026-10-16-verification-slice-v1-cli-quiet-vorschau.md),
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../../../reviews/2026-10-16-verification-slice-v1-cli-verbosity-letzte-option-gewinnt.md).
