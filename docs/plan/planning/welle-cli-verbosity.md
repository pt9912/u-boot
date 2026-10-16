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

[`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), der Code in `internal/adapter/driving/cli/` und
die Spezifikation sagen dasselbe über die Verbosity: `--quiet` zeigt nur Fehler (die Erfolgs- und Zusammenfassungsmeldung aller
Subkommandos entfällt), mit vier Ausnahmen — die Ausweisung von Sicherheitslockerungen im Sandbox-Profil
([`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)), die Ausweisung der Degradations-Fallbacks
([`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)), der Compose-Fortschritt von `up`
([`LH-NFA-PERF-002`](../../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker)) und die Pfad-Zusammenfassung von `--force`
([`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)) bleiben sichtbar —, die Ausgabe der Vorschau-Optionen `--dry-run` und `--diff` ist das Ergebnis des Aufrufs und bleibt (wie die der Abfragebefehle; [`LH-FA-CLI-007`](../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe)), und die zuletzt
angegebene Option gewinnt. Ein Test pro Aussage belegt es; die Spezifikation beschreibt das Verhalten erst, wenn Lastenheft und
Code es tragen.

## 2. Trigger (Welle startet)

- Die Welle [`welle-harness-wartung`](done/welle-harness-wartung.md) liegt in `done/` (geschlossen 2026-10-16).
- Die Entscheidungen des Projektinhabers liegen vor (2026-10-16): `--quiet` zeigt nur Fehler; die Ausnahmen sind die
  Sandbox-Lockerungs-Hinweise (CR [`slice-cr-cli-quiet-fehler-und-warnungen`](done/slice-cr-cli-quiet-fehler-und-warnungen.md))
  sowie Degradations-Fallbacks, Compose-Fortschritt von `up` und `--force`-Zusammenfassung (CR
  [`slice-cr-cli-quiet-weitere-ausnahmen`](done/slice-cr-cli-quiet-weitere-ausnahmen.md); angenommene Fassung von
  [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)); der Trigger ist damit eingetreten.

## 3. Closure-Trigger (Welle schließt)

Das *Mehr* gegenüber den sechs Slice-DoDs: Die drei Stellen, an denen der Vertrag steht, sagen dasselbe. Die drei CRs ändern das
Lastenheft, die drei Code-Slices den Adapter, und die Spezifikation beschreibt erst danach das Verhalten.

- Alle sechs Slices der Welle liegen in `done/` (die zwei geschlossenen CRs, der CR
  [`slice-cr-cli-quiet-vorschau`](in-progress/slice-cr-cli-quiet-vorschau.md), [`slice-v1-cli-quiet-nur-fehler`](done/slice-v1-cli-quiet-nur-fehler.md),
  [`slice-v1-cli-quiet-vorschau`](open/slice-v1-cli-quiet-vorschau.md) und [`slice-v1-cli-verbosity-letzte-option-gewinnt`](open/slice-v1-cli-verbosity-letzte-option-gewinnt.md)); das Lastenheft steht auf 0.10.0.
- Die Spezifikation trägt einen Eintrag zur Pegelabbildung und zur Rangfolge, der Lastenheft (0.10.0) und Code deckungsgleich
  beschreibt; er entstand nach dem Abschluss aller sechs Slices (Commit-Reihenfolge belegt es).
- `make doc-trace` meldet keine Waise unter [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging).
- `make gates` ist grün.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-cr-cli-quiet-fehler-und-warnungen`](done/slice-cr-cli-quiet-fehler-und-warnungen.md) | CR: Ausnahme der Sicherheitslockerungen unter `--quiet` (zuerst, geschlossen) | [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) |
| [`slice-cr-cli-quiet-weitere-ausnahmen`](done/slice-cr-cli-quiet-weitere-ausnahmen.md) | CR: Ausnahmen Degradations-Fallbacks, Compose-Fortschritt und `--force`-Zusammenfassung unter `--quiet` (Lastenheft 0.9.0) | [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-NFA-PERF-002`](../../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker), [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) |
| [`slice-v1-cli-quiet-nur-fehler`](done/slice-v1-cli-quiet-nur-fehler.md) | Code: `--quiet` zeigt nur Fehler in allen Subkommandos, vier Ausnahmen bleiben (Lastenheft 0.9.0; geschlossen) | [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) |
| [`slice-cr-cli-quiet-vorschau`](in-progress/slice-cr-cli-quiet-vorschau.md) | CR: die Vorschau von `--dry-run` und `--diff` bleibt unter `--quiet` (Lastenheft 0.10.0) | [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-007`](../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) |
| [`slice-v1-cli-quiet-vorschau`](open/slice-v1-cli-quiet-vorschau.md) | Code: `--dry-run --quiet` und `--diff --quiet` geben die Vorschau aus (nach dem CR, Lastenheft 0.10.0) | [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-007`](../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) |
| [`slice-v1-cli-verbosity-letzte-option-gewinnt`](open/slice-v1-cli-verbosity-letzte-option-gewinnt.md) | Code: die zuletzt angegebene Option gewinnt, danach der Spezifikations-Eintrag (letzter Slice) | [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |

## 5. Abhängigkeiten

- Blockiert: keine Welle.
- Wird blockiert von: nichts; die Entscheidungen des Projektinhabers (Wirkung von `--quiet`, siehe §2) liegen vor. Innerhalb der
  Welle geht der CR [`slice-cr-cli-quiet-fehler-und-warnungen`](done/slice-cr-cli-quiet-fehler-und-warnungen.md) dem CR
  [`slice-cr-cli-quiet-weitere-ausnahmen`](done/slice-cr-cli-quiet-weitere-ausnahmen.md) voraus; dieser geht dem Slice
  [`slice-v1-cli-quiet-nur-fehler`](done/slice-v1-cli-quiet-nur-fehler.md) voraus, weil die Ausnahmen erst nach Lastenheft 0.9.0
  im Vertrag stehen. Der CR [`slice-cr-cli-quiet-vorschau`](in-progress/slice-cr-cli-quiet-vorschau.md) folgt auf diesen Code-Slice und geht dem Code-Slice
  [`slice-v1-cli-quiet-vorschau`](open/slice-v1-cli-quiet-vorschau.md) voraus, weil der Vorschau-Satz erst nach Lastenheft 0.10.0 im Vertrag steht; der Code-Slice
  beginnt erst nach der Closure von [`slice-v1-cli-quiet-nur-fehler`](done/slice-v1-cli-quiet-nur-fehler.md), weil beide dieselben Quiet-Zweige berühren.
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](open/slice-v1-cli-verbosity-letzte-option-gewinnt.md) kommt zuletzt (er trägt den Spezifikations-Eintrag nach allen anderen); logisch ist er von den Vorschau-Slices unabhängig,
  berührt aber ebenfalls `root.go`.

## 6. Out-of-Scope für diese Welle

- Eine neue Pegelstufe zwischen Info und Debug: `--verbose` und `--debug` bleiben beide Debug; der Kommentar in `root.go` nennt das als vorläufig.
- Konfigurations- oder Umgebungsquellen für den Pegel: es gibt nur die drei Flags.
- Die Sandbox-Hinweise bei `init` und `generate` auf den Logger umzustellen: ein anderer Vorgang; der CR hält den Befund fest, entscheidet aber nur über `--quiet`.
- Die Homebrew-Smoke-Härtung ([`slice-v2-homebrew-smoke-haertung`](open/slice-v2-homebrew-smoke-haertung.md)) und die Docker-Desktop-Messung: ohne Bezug zur Verbosity.

## 7. Closure-Notiz

Offen.
