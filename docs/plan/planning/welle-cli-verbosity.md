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
die Spezifikation sagen dasselbe über die Verbosity: `--quiet` zeigt nur Fehler, mit einer Ausnahme — die Ausweisung von
Sicherheitslockerungen im Sandbox-Profil ([`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)) bleibt sichtbar —, und die zuletzt
angegebene Option gewinnt. Ein Test pro Aussage belegt es; die Spezifikation beschreibt das Verhalten erst, wenn Lastenheft und
Code es tragen.

## 2. Trigger (Welle startet)

- Die Welle [`welle-harness-wartung`](done/welle-harness-wartung.md) liegt in `done/` (geschlossen 2026-10-16).
- Die Entscheidung des Projektinhabers liegt vor (2026-10-16): `--quiet` zeigt nur Fehler; die einzige Ausnahme sind die
  Sandbox-Lockerungs-Hinweise. Sie steht im CR
  [`slice-cr-cli-quiet-fehler-und-warnungen`](in-progress/slice-cr-cli-quiet-fehler-und-warnungen.md) (vorgeschlagene Fassung von
  [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)); der Trigger ist damit eingetreten.

## 3. Closure-Trigger (Welle schließt)

Das *Mehr* gegenüber den drei Slice-DoDs: Die drei Stellen, an denen der Vertrag steht, sagen dasselbe. Der CR ändert das
Lastenheft, die beiden Code-Slices den Adapter, und die Spezifikation beschreibt erst danach das Verhalten.

- Alle drei Slices der Welle liegen in `done/`.
- Die Spezifikation trägt einen Eintrag zur Pegelabbildung und zur Rangfolge, der Lastenheft (0.8.0) und Code deckungsgleich
  beschreibt; er entstand nach dem Abschluss aller drei Slices (Commit-Reihenfolge belegt es).
- `make doc-trace` meldet keine Waise unter [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging).
- `make gates` ist grün.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-cr-cli-quiet-fehler-und-warnungen`](in-progress/slice-cr-cli-quiet-fehler-und-warnungen.md) | CR: Ausnahme der Sicherheitslockerungen unter `--quiet` (zuerst) | [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) |
| [`slice-v1-cli-quiet-nur-fehler`](open/slice-v1-cli-quiet-nur-fehler.md) | Code: `--quiet` zeigt nur Fehler, Sandbox-Lockerungs-Hinweise bleiben | [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) |
| [`slice-v1-cli-verbosity-letzte-option-gewinnt`](open/slice-v1-cli-verbosity-letzte-option-gewinnt.md) | Code: die zuletzt angegebene Option gewinnt, danach der Spezifikations-Eintrag | [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |

## 5. Abhängigkeiten

- Blockiert: keine Welle.
- Wird blockiert von: nichts; die Entscheidung des Projektinhabers (Wirkung von `--quiet`, siehe §2) liegt vor. Innerhalb der
  Welle geht der CR dem Slice [`slice-v1-cli-quiet-nur-fehler`](open/slice-v1-cli-quiet-nur-fehler.md) voraus, weil die Ausnahme
  erst nach Lastenheft 0.8.0 im Vertrag steht; dieser Slice geht
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](open/slice-v1-cli-verbosity-letzte-option-gewinnt.md) voraus, weil beide
  `root.go` berühren (logisch sind sie unabhängig).

## 6. Out-of-Scope für diese Welle

- Eine neue Pegelstufe zwischen Info und Debug: `--verbose` und `--debug` bleiben beide Debug; der Kommentar in `root.go` nennt das als vorläufig.
- Konfigurations- oder Umgebungsquellen für den Pegel: es gibt nur die drei Flags.
- Die Sandbox-Hinweise bei `init` und `generate` auf den Logger umzustellen: ein anderer Vorgang; der CR hält den Befund fest, entscheidet aber nur über `--quiet`.
- Die Homebrew-Smoke-Härtung ([`slice-v2-homebrew-smoke-haertung`](open/slice-v2-homebrew-smoke-haertung.md)) und die Docker-Desktop-Messung: ohne Bezug zur Verbosity.

## 7. Closure-Notiz

Offen.
