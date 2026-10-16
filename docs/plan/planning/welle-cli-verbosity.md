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
die Spezifikation sagen dasselbe über die Verbosity: Die zuletzt angegebene Option gewinnt, und `--quiet` hat die Wirkung, die der
Projektinhaber im CR festlegt. Ein Test pro Aussage belegt es; die Spezifikation beschreibt das Verhalten erst, wenn Lastenheft und
Code es tragen.

## 2. Trigger (Welle startet)

- Die Welle [`welle-harness-wartung`](done/welle-harness-wartung.md) liegt in `done/` (geschlossen 2026-10-16).
- Der Projektinhaber hat entschieden, ob `--quiet` Ausgaben der Sandbox-Hinweise unterdrücken darf; die Entscheidung steht im CR
  [`slice-cr-cli-quiet-fehler-und-warnungen`](open/slice-cr-cli-quiet-fehler-und-warnungen.md) (Prüfauftrag §1 abgeschlossen).

## 3. Closure-Trigger (Welle schließt)

Das *Mehr* gegenüber den beiden Slice-DoDs: Die drei Stellen, an denen der Vertrag steht, sagen dasselbe. Der CR ändert das
Lastenheft, der Code-Slice den Adapter, und die Spezifikation beschreibt erst danach das Verhalten.

- Beide Slices der Welle liegen in `done/`.
- Die Spezifikation trägt einen Eintrag zur Pegelabbildung und zur Rangfolge, der Lastenheft (0.8.0) und Code deckungsgleich
  beschreibt; er entstand nach dem Abschluss beider Slices (Commit-Reihenfolge belegt es).
- `make doc-trace` meldet keine Waise unter [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging).
- `make gates` ist grün.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-cr-cli-quiet-fehler-und-warnungen`](open/slice-cr-cli-quiet-fehler-und-warnungen.md) | CR: Wirkung von `--quiet` (zuerst) | [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |
| [`slice-v1-cli-verbosity-letzte-option-gewinnt`](open/slice-v1-cli-verbosity-letzte-option-gewinnt.md) | Code: die zuletzt angegebene Option gewinnt, danach der Spezifikations-Eintrag | [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |

## 5. Abhängigkeiten

- Blockiert: keine Welle.
- Wird blockiert von: keiner Welle, aber von einer Entscheidung des Projektinhabers (Wirkung von `--quiet`, siehe §2). Innerhalb der Welle
  geht der CR dem Code-Slice voraus, weil der Code-Slice die Pegelabbildung des CR übernimmt.

## 6. Out-of-Scope für diese Welle

- Eine neue Pegelstufe zwischen Info und Debug: `--verbose` und `--debug` bleiben beide Debug; der Kommentar in `root.go` nennt das als vorläufig.
- Konfigurations- oder Umgebungsquellen für den Pegel: es gibt nur die drei Flags.
- Die Sandbox-Hinweise bei `init` und `generate` auf den Logger umzustellen: ein anderer Vorgang; der CR hält den Befund fest, entscheidet aber nur über `--quiet`.
- Die Homebrew-Smoke-Härtung ([`slice-v2-homebrew-smoke-haertung`](open/slice-v2-homebrew-smoke-haertung.md)) und die Docker-Desktop-Messung: ohne Bezug zur Verbosity.

## 7. Closure-Notiz

Offen.
