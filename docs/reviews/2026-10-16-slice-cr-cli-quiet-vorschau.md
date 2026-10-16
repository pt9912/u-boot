# Review-Report: slice-cr-cli-quiet-vorschau — 2026-10-16

**Review-Art:** Plan/Design — Diff gegen den CR-Kern des Slice, [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), die angrenzenden Lastenheft-Zusagen und Hard Rules (Maintainability; nicht gegen die DoD).

**Gegenstand:** Commit `49f459a` für den Slice [`slice-cr-cli-quiet-vorschau`](../plan/planning/done/slice-cr-cli-quiet-vorschau.md) (Welle `welle-cli-verbosity`): `spec/lastenheft.md` (Version 0.10.0, §3 `LH-FA-CLI-005`, §16-Zeile), `CHANGELOG.md` `[Unreleased]`. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

**Skill:** `.harness/skills/reviewer.md` (aktueller Stand im Arbeitsbaum), HEAD `49f459a` · <!-- d-check:ignore (Skill-Pfad als Zitat) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingangs-Kontext:**

- Slice [`slice-cr-cli-quiet-vorschau`](../plan/planning/done/slice-cr-cli-quiet-vorschau.md) (§1 CR-Kern, §2 DoD, §6 Risiken) und Folge-Slice [`slice-v1-cli-quiet-vorschau`](../plan/planning/done/slice-v1-cli-quiet-vorschau.md)
- [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), [`LH-FA-INIT-005`](../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz)
- Vorherige Reports der Welle: `2026-10-16-slice-cr-cli-quiet-fehler-und-warnungen.md`, `2026-10-16-slice-cr-cli-quiet-weitere-ausnahmen.md` <!-- d-check:ignore (Report-Dateinamen als Zitat) -->
- `AGENTS.md` (Hard Rules §3.4, §3.7)

**Systematische Suche gegen ältere Zusagen** (Lastenheft, Begriffe `--quiet`, `quiet`, `--dry-run`, `--diff`, `Vorschau`, `Zusammenfassung`, `Rückfrage`, `Bestätigung`, `geplant`): Treffer nur in `LH-FA-CLI-005`, `-005A`, `-007`, `-008`, `LH-NFA-USE-004` (JSON-Schema), `--force`-Zusammenfassung (`LH-FA-CLI-005A`, `LH-FA-INIT-005`, `init`-Optionsliste). `LH-FA-GEN-*` und `LH-FA-CONF-*` tragen keine Aussage zu Vorschau oder Verbosity. Keine Kollision über das hinaus, was der Satz auflöst. <!-- d-check:ignore (Kennung als Zitat) -->

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Eintrag nennt die Ausgabe von `--dry-run` als Ergebnis, das unter `--quiet` bleibt, und schiebt im selben Satz „the plan summary of `--dry-run`" auf ein späteres Release. Nach Slice §1 ist die Plan-Zusammenfassung („Would initialize …", „Would add service …") genau die `--dry-run`-Ausgabe; der Satz sagt damit zugleich, sie bleibt, und sie bleibt noch nicht. Im Code schweigt `--dry-run --quiet` heute mit Exit-Code 0. | Doku-Drift (Reviewer-Skill MEDIUM „Doku … driftet gegen Spec oder gegen das implementierte Verhalten"); [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run) | `CHANGELOG.md` · "is the result of the call and stays under `--quiet`; the plan summary of `--dry-run` and the rule that the last verbosity option wins" | nein — Review-only (Release-Cut vor `slice-v1-cli-quiet-vorschau`; `make docs-check` prüft keinen Verhaltensabgleich) | Changelog-Zusage widerspricht sich selbst bzw. dem Ist-Stand | <!-- d-check:ignore (Kennung als Zitat) -->
| F-2 | INFO | Das Lastenheft nennt „geplanter Änderungsplan" und „Diff-Darstellung"; Slice und CHANGELOG sprechen von „plan summary" / „Vorschau-Zusammenfassung". Die Begriffe sind im Vertrag nicht abgegrenzt: Ob die heutige Vergangenheitsform bei `remove`/`generate`/`config set` („Removed service …", „config: a → b.") unter `--dry-run` ein „Änderungsplan" ist, entscheidet der Vertrag nicht. | [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run) („welche Dateien erzeugt, geändert oder gelöscht würden") | `spec/lastenheft.md` · "der geplante Änderungsplan und die Diff-Darstellung" | nein — Review-only; der Slice benennt den „Would …"-Wortlaut als Nicht-Ziel (§1 und Folge-Slice) | Vertragsbegriff „Änderungsplan" nicht von „Zusammenfassung" getrennt |
| F-3 | INFO | Die §16-Zeile weicht vom Entwurf in Slice §1 ab: „die Erfolgsmeldung eines tatsächlichen Schreibens entfällt weiter" wird zu „die Erfolgs- und Zusammenfassungsmeldung eines tatsächlichen Schreibens entfällt auch bei `--diff` ohne `--dry-run`"; „nach LH-FA-CLI-007 und -008" wird durch Verweise im Klammerzusatz ersetzt. Die Lastenheft-Zeile ist wortgleich zum Absatz in §3 und damit korrekter als der Entwurf; der Slice trägt den Entwurf weiter. | Maintainability (Entwurf und Artefakt sollen nicht auseinanderlaufen) | `docs/plan/planning/done/slice-cr-cli-quiet-vorschau.md` · "die Erfolgsmeldung eines tatsächlichen Schreibens entfällt" | nein — Review-only | Slice-Entwurf der Historie-Zeile weicht vom gesetzten Text ab | <!-- d-check:ignore (Kennung als Zitat) -->
| F-4 | INFO | Das Lastenheft sagt zu `--quiet --json` nichts; die Zusage „`--quiet --json` is unchanged" steht in CHANGELOG, Benutzerhandbuch und Spezifikation. Der neue Satz berührt den Fall nicht (die JSON-Ausgabe von `--dry-run --json` bleibt nach [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run) streng maschinenlesbar); es entsteht keine Kombination, die der Vertrag widersprüchlich lässt. | [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) | `spec/lastenheft.md` · "Bei gleichzeitiger Verwendung von `--dry-run` und `--json` muss die Ausgabe streng maschinenlesbar" | nein — Review-only | Vertragsschweigen zu `--quiet --json`, folgenlos |
| F-5 | INFO | Das Benutzerhandbuch (Optionstabelle zu `--quiet`) nennt die Vorschau nicht als Ergebnis und sagt „Erfolgs- und Zusammenfassungsmeldungen … entfallen". Das ist der Code-Stand; der Nachzug steht als Abnahmekriterium 3 im Folge-Slice. Zuständig: Verifier/Folge-Slice. | Slice-Abgrenzung (Doku-Nachzug im Code-Slice) | `docs/user/benutzerhandbuch.md` · "Die Abfragebefehle `config get`, `config list`, `logs` und `template list` behalten ihre Ausgabe." | nein — Review-only | Handbuch folgt dem Code, nicht dem Lastenheft 0.10.0 |

## Prüfpunkte aus dem Auftrag

1. **Wörtlichkeit.** Der neue Satzteil in `LH-FA-CLI-005` entspricht Wort für Wort dem kursiven CR-Kern in Slice §1 (verglichen Satz für Satz, inklusive Gedankenstrichen und beider Links). Unverändert sind: die vier Stufen, der Einleitungsabsatz, die vier Ausnahmen, der Abfragebefehle-Satz („ist das Ergebnis des Befehls und bleibt.") und die Mehrfachangabe-Regel. Ohne Befund. <!-- d-check:ignore (Kennung als Zitat) -->
2. **Eindeutigkeit.**
   - `init --dry-run --quiet` gibt unter dem Satz den Plan aus; dass der Code heute schweigt, ist Gegenstand des Folge-Slice `slice-v1-cli-quiet-vorschau`, nicht dieses CR. Der Vertrag ist insoweit eindeutig; die Abgrenzung der Begriffe siehe F-2. <!-- d-check:ignore (Kennung als Zitat) -->
   - `--diff` ohne `--dry-run`: Diff bleibt (nach [`LH-FA-CLI-008`](../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) „zeigt den geplanten Endzustand als Vorschau"), die Erfolgsmeldung entfällt; der Satz sagt es ausdrücklich. Kein Widerspruch zu CLI-008 (Beispiel „Vorschau mit anschließendem Schreiben").
   - `config set --dry-run --quiet`: `config set` ist in der Liste der ändernden Subkommandos und in CLI-005A als bestätigungspflichtig genannt; der Satz erfasst es über „Vorschau-Optionen". Kein Widerspruch.
   - Wörtlich zu `LH-FA-CLI-007`/`-008`: „Änderungsplan" und „Diff-Darstellung" stammen aus CLI-007 („ausschließlich aus dem geplanten Änderungsplan und der Diff-Darstellung"); die Verweise lösen. Warnungen neben der Vorschau entfallen nach „Alle anderen Hinweise und Warnungen entfallen unter `--quiet`"; CLI-007 „ausschließlich" wird dadurch eher bestätigt. <!-- d-check:ignore (Kennung als Zitat) -->
   - `--quiet --json --dry-run`: siehe F-4; unberührt.
3. **Versionssprung und §16.** 0.9.0 → 0.10.0 folgt dem Muster der Zeilen 0.5.0/0.7.0/0.8.0/0.9.0 (Minor bei Verhaltensänderung mit Code-Folge; Patches 0.3.x für Nachträge ohne geänderte Zusage). Die Zeile trägt Datum, Kennung als Link, Quelle „Vereinbarung mit dem Projektinhaber". Abweichung vom Entwurf: F-3. Das Versionsfeld ist 0.10.0, die Spaltenbreite der Kopftabelle ist erhalten.
4. **Decken-Regel (`AGENTS.md` §3.4).** Der Diff fügt nur Verweise auf `LH-FA-CLI-007` und `-008` (gleiches Stratum, aufwärts bzw. seitwärts) hinzu; kein Verweis auf Slice, Carveout, Roadmap oder ADR. Die Anker `#lh-fa-cli-007--dry-run` und `#lh-fa-cli-008--diff-ausgabe` existieren (Überschriften „LH-FA-CLI-007 – Dry Run", „LH-FA-CLI-008 – Diff-Ausgabe"); `make docs-check` meldet 0 Befunde. <!-- d-check:ignore (Kennung als Zitat) -->
5. **CHANGELOG.** F-1. Konsistent mit Lastenheft 0.10.0 sind die Versionsangabe und der Satz zur Vorschau; im Widerspruch zum Ist-Stand steht die Gleichzeitigkeit von „bleibt" und „folgt später". Failure-Szenario: Ein Release-Cut vor `slice-v1-cli-quiet-vorschau` liefert einen Changelog aus, der `--dry-run` unter `--quiet` als Vertragsteil nennt, während `u-boot init --dry-run --quiet` nichts ausgibt und mit 0 endet; die Verzögerungsklausel nennt „plan summary", das der Leser nicht von der Ausgabe unterscheiden kann. <!-- d-check:ignore (Kennung als Zitat) -->
6. **Sprachvertrag und Auslegungen.** Lastenheft deutsch, CHANGELOG englisch: vertragskonform. `AGENTS.md` §3.7: die neuen Sätze beschreiben den Zustand im Indikativ, keine verworfene Alternative, keine Chronik (die §16-Zeile ist die dafür vorgesehene Historie). Auslegungen des Planners: *Version 0.10.0 statt 0.9.1* — vertragskonform (Muster oben; im Slice als zu bestätigen markiert, §1). *`--diff` ohne `--dry-run`* — vertragskonform; der Slice führt es zutreffend als offenes Risiko mit Projektinhaber-Bestätigung (§6), im Lastenheft sauber getrennt. *„Would …"-Wortlaut bei `remove`/`generate`/`config set`* — außerhalb des Vertragsdiffs, im Slice als Nicht-Ziel begründet; Berührung mit dem Vertrag siehe F-2. DoD-Zeilen (§2): nicht Gegenstand dieses Reviews (Verifier).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `spec/lastenheft.md` Kopftabelle, Wortlaut `LH-FA-CLI-005` gegen CR-Kern | geprüft, ohne Befund | <!-- d-check:ignore (Kennung als Zitat) -->
| Stufen, vier Ausnahmen, Abfragebefehle-Satz, Mehrfachangabe-Regel | geprüft, ohne Befund (unverändert) |
| Verweise und Anker im Lastenheft, Decken-Regel | geprüft, ohne Befund (`make docs-check`: 263 Dateien, 0 Befunde) |
| §16-Historie-Zeile 0.10.0 (Muster, Quelle) | geprüft, ohne Befund (Abweichung vom Entwurf: F-3) |
| Ältere Zusagen `LH-FA-CLI-005A`, `LH-FA-INIT-005` (`--force`-Zusammenfassung), `LH-FA-GEN-*`, `LH-FA-CONF-*` | geprüft, ohne Befund | <!-- d-check:ignore (Kennung als Zitat) -->
| `AGENTS.md` §3.7 in Lastenheft-Text und Slice | geprüft, ohne Befund |
| Sprachvertrag Lastenheft/CHANGELOG | geprüft, ohne Befund |

Nicht geprüft:

- `internal/adapter/driving/cli` (Code) — der Diff berührt keinen Code; Umsetzung liegt in `slice-v1-cli-quiet-vorschau`. <!-- d-check:ignore (Kennung als Zitat) -->
- `make gates`, `make test-docker`, `make verify-depguard` — kein Code-, Docker- oder Importregel-Pfad berührt; nur `make docs-check` gelaufen.
- DoD-Zeilen des Slice — Verifier-Aufgabe.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Changelog-Zusage widerspricht sich selbst bzw. dem Ist-Stand · Vertragsbegriff „Änderungsplan" nicht von „Zusammenfassung" getrennt · Slice-Entwurf der Historie-Zeile weicht vom gesetzten Text ab · Vertragsschweigen zu `--quiet --json`, folgenlos · Handbuch folgt dem Code, nicht dem Lastenheft

## Verdikt

**Merge-blockierend:** nein — kein HIGH; F-1 (MEDIUM) betrifft einen unveröffentlichten Changelog-Eintrag und ist vor dem nächsten Release-Cut zu klären, nicht vor dem Lastenheft-Change.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation.
