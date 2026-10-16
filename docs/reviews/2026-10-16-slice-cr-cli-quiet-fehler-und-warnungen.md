# Review-Report: slice-cr-cli-quiet-fehler-und-warnungen — 2026-10-16

**Review-Art:** Plan — Diff gegen den CR-Slice (§1 CR-Kern, §2, §3, §6), `LH-FA-DEV-006`/`-007`/`-008`, das Muster in §16 der Lastenheft-Historie und die Hard Rules (Maintainability; nicht gegen die DoD). <!-- d-check:ignore (Kennungen als Zitat) -->

**Gegenstand:** Commit `b5fa4b4` (Lastenheft 0.8.0, `CHANGELOG.md` `[Unreleased]`, Wegfall der Zeile „Nichts in Arbeit." in der Roadmap) für `slice-cr-cli-quiet-fehler-und-warnungen` (Welle `welle-cli-verbosity`). <!-- d-check:ignore (Kennungen als Zitat im Lauf-Beleg) -->

**Skill:** `.harness/skills/reviewer.md` @ HEAD `b5fa4b4` · <!-- d-check:ignore (Skill-Pfad als Zitat) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingangs-Kontext:**

- `slice-cr-cli-quiet-fehler-und-warnungen` (§1 CR-Kern, §2, §3, §6), `slice-v1-cli-quiet-nur-fehler` (`open/`) und `welle-cli-verbosity` (nur gelesen) <!-- d-check:ignore (Kennungen als Zitat im Lauf-Beleg) -->
- [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), Lastenheft §16
- Code nur gelesen: `internal/hexagon/application/devcontainer_sandbox.go` (`sandboxRelaxationWarnings`) <!-- d-check:ignore (Pfad als Zitat) -->
- `AGENTS.md` (Hard Rules §3.4, §3.7), `harness/conventions.md`
- Sensors: `make docs-check` (nach Erstellung dieses Reports: 0 Befunde). Nicht gelaufen: `make gates`, `make test`, `make test-docker`, `make verify-depguard` (Diff ohne Go-, Docker- oder Importregel-Änderung).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Ausnahme-Satz grenzt „Sicherheitslockerung" nicht ab; die Auslegung des CR §1 (nur die „security relaxation"-Hinweise der Codes `LH-FA-DEV-007`/`-008`, nicht „no git remote", „has no effect", Recreate, Drittanbieter) steht nicht im Lastenheft. Der Satz zitiert nur `-006` und `-007`, die Hinweise tragen aber auch `-008`; `-006` nennt „Capabilities" nur in der Klammer der Ausweispflicht. Dieselbe Kennung `LH-FA-DEV-007` trägt in `sandboxRelaxationWarnings` sowohl den Lockerungs-Hinweis als auch den „has no effect"-Hinweis, die Code-Kennung trennt beides nicht. Die Degradations-Fallbacks aus `-007` („Wegfall der Restriktion", Container ohne Daemon; „in der Befehlsausgabe ausgewiesen") sind weder ein- noch ausgeschlossen. | [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) | spec/lastenheft.md · "Die Ausweisung jeder Sicherheitslockerung in der Befehlsausgabe" | nein — Review-only (Plan-/Spec-Review) | Ausnahme im Vertrag ohne Abgrenzung des Geltungsbereichs <!-- d-check:ignore (Kennungen als Zitat) --> |
| F-2 | INFO | Der Slice-§6 schließt das Risiko „Reichweite der Ausnahme" mit „entfallen" und legt die Code-Grenze in den Umsetzungs-Slice. Das Lastenheft ist dann der Vertrag, gegen den die Abnahmetests geschrieben werden, trägt die Grenze aber nicht (F-1). Zuständig für die Entscheidung: Architect bzw. Projektinhaber (Vertragstext). | Maintainability | docs/plan/planning/done/slice-cr-cli-quiet-fehler-und-warnungen.md · "der Umsetzungs-Slice [`slice-v1-cli-quiet-nur-fehler`] legt die Code-Grenze fest" <!-- d-check:ignore (Zitat) --> | nein — Review-only | Risiko-Ausgang verschiebt Vertragsgrenze in den Code |
| F-3 | INFO | Die Versionsfortschreibung 0.7.0 → 0.8.0 ist mit §16 vereinbar, aber nicht erzwungen. Patch trugen Präzisierungen bei gleicher Zusage (0.3.1, 0.3.3) und Statuswechsel (0.3.4, 0.3.5); Minor trugen einen neuen Wert (0.5.0), die Einengung einer Zusage ohne neuen Wert (0.7.0) und die Strukturänderung (0.6.0). Eine Ausnahme von einer Stufenzusage liegt näher an 0.7.0 als an 0.3.1; §16 enthält keine ausgesprochene Regel, 0.7.1 wäre nicht falsch. | Maintainability | spec/lastenheft.md · "| 0.8.0 | 2026-10-16 |" | nein — Review-only | Versionsschritt ohne dokumentierte Schnittregel |
| F-4 | INFO | Die §16-Zeile weicht im Wortlaut von der Historie-Zeile des CR §1 ab (CR: „`--quiet` bleibt „nur Fehler"; …"; Lastenheft: „Ausnahme … ergänzt: …, ein Skript mit `--quiet` schaltet sie nicht stumm. Die Stufenbeschreibungen … bleiben unverändert."). Die DoD fordert nur das Vorhandensein samt Verweis; der CR-Kern (Anforderungstext) stimmt wörtlich überein. | Maintainability | spec/lastenheft.md · "ein Skript mit `--quiet` schaltet sie nicht stumm" | nein — Review-only | Historie-Zeile weicht vom Plan-Entwurf ab |
| F-5 | INFO | Der CHANGELOG-Eintrag sagt im Präsens, „`--quiet` shows only errors; the last verbosity option wins", beides entspricht erst nach den Folge-Slices dem Code (Pegel heute `Warn`, Mehrfachangabe offen). Der Satz „the implementation follows in a later release" deckt das ab und stimmt mit Plan und Welle überein (`slice-v1-cli-quiet-nur-fehler` in `open/`, `slice-v1-cli-verbosity-letzte-option-gewinnt` in `open/`). Hinweis für den Verifier: Ein Release vor den Folge-Slices trägt einen Eintrag, dessen Präsens dem Binary widerspricht. | AGENTS.md §5 Regel 5 | CHANGELOG.md · "the implementation follows in a later release" | nein — Review-only | CHANGELOG-Präsens vor Umsetzung <!-- d-check:ignore (Kennungen als Zitat) --> |

## Konsistenzprüfung (Auftrag 1–5)

- **CR-Kern ↔ Lastenheft:** Der Absatz „Eine Ausnahme gilt für Sicherheitslockerungen: …" in `LH-FA-CLI-005` ist wörtlich identisch mit dem CR-Kern. Der Diff fügt dort nur Zeilen hinzu; die vier Stufenbeschreibungen und der Absatz zur Mehrfachangabe sind unverändert. Der Abschnitt steht vor der Mehrfachangabe-Regel, wie im CR. <!-- d-check:ignore (Kennungen als Zitat) -->
- **Eindeutigkeit der Ausnahme:** siehe F-1. Ableitbar aus dem Lastenheft ist: „has no effect"-Hinweise sind keine Lockerung (nichts wird erzeugt), „no git remote" ist keine Lockerung, Recreate-Warnung (`LH-FA-UP-003`) und Drittanbieter-Hinweise nicht. Nicht ableitbar: ob `-008` (Egress, `NET_ADMIN`) trotz fehlender Zitation trägt (nur über das Wort „Capabilities" in `-006`) und ob Degradations-Fallbacks zählen. Failure-Szenario: Ein Implementer filtert nach Kennung `LH-FA-DEV-007`, zeigt dadurch „has no effect" unter `--quiet` (zu viel) und lässt `LH-FA-DEV-008` weg (zu wenig), oder behandelt Degradations-Hinweise abweichend vom Projektinhaber; Test und Vertrag stimmen formal, der Wortlaut trägt beide Lesarten. <!-- d-check:ignore (Kennungen als Zitat) -->
- **Versionssprung:** siehe F-3.
- **Decken-Regel (`AGENTS.md` §3.4):** Das Lastenheft verweist im Diff nur auf Lastenheft-Anker (`#lh-fa-dev-006--sandbox-profil`, `#lh-fa-dev-007--…`, `#lh-fa-cli-005--…` existieren); kein Verweis auf Slice, Carveout, Roadmap oder ADR. Die Spalte „Verweis" nennt den externen Vorgang („Vereinbarung mit dem Projektinhaber"). Der CHANGELOG verweist aufwärts auf `LH-FA-CLI-005`; seine Aussage „implementation follows in a later release" stimmt mit Slice-Plan und Welle überein. <!-- d-check:ignore (Kennungen als Zitat) -->
- **Traceability:** Commit-Nachricht nennt `LH-FA-CLI-005`, `LH-FA-DEV-006` und den Slice; Versionsfeld, Anforderung und §16-Zeile (die drei Spuren laut §16-Kopf) sind alle gesetzt; Datum 2026-10-16 trägt das Annahmedatum. <!-- d-check:ignore (Kennungen als Zitat) -->
- **§3.7 / Sprachvertrag:** Im Lastenheft- und Slice-Diff keine Kommentare oder Zustandsfelder mit Konjunktiv über verworfene Alternativen; Die Roadmap-Änderung entfernt „Nichts in Arbeit." und führt die Welle als Zustand, keine Chronik. Lastenheft deutsch, CHANGELOG englisch: eingehalten.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `spec/lastenheft.md` Versionsfeld, §16-Zeile (Format, Reihenfolge, Verweis-Spalte) | geprüft, ohne Befund (F-3/F-4 nur Hinweis) |
| `spec/lastenheft.md` Link-Anker der Ausnahme | geprüft, ohne Befund |
| `CHANGELOG.md` `[Unreleased]` (Sprache, Abschnitt `Changed`, Link) | geprüft, ohne Befund (F-5 nur Hinweis) |
| `docs/plan/planning/in-progress/roadmap.md` (Entfall „Nichts in Arbeit.") | geprüft, ohne Befund |
| Referenzmodell Rang 1–4 (kein Abwärtsverweis) | geprüft, ohne Befund |
| Lifecycle des Slice (liegt in `in-progress/`, kein Doppelort) | geprüft, ohne Befund |
| Hexagon-/depguard-/Dual-Classifier-Linsen, Exit-Codes | nicht berührt, kein Go-Code im Diff |
| `docs/user/benutzerhandbuch.md` | nicht berührt; die Nachführung ist Abnahmekriterium 2 des Umsetzungs-Slice |
| `make gates`, `make test`, `make test-docker`, `make verify-depguard` | nicht ausgeführt: kein Go-/Docker-Pfad und keine Importregel im Diff |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Ausnahme im Vertrag ohne Abgrenzung des Geltungsbereichs · Risiko-Ausgang verschiebt Vertragsgrenze in den Code · Versionsschritt ohne dokumentierte Schnittregel · Historie-Zeile weicht vom Plan-Entwurf ab · CHANGELOG-Präsens vor Umsetzung

## Verdikt

**Merge-blockierend:** nein für Datei-Konsistenz; F-1 (MEDIUM) ist vor der Closure des CR zu klären, weil der Umsetzungs-Slice seine Abnahmetests gegen diesen Vertragstext schreibt. Ob der Satz präzisiert oder die Auslegung bewusst dem Umsetzungs-Slice überlassen bleibt, entscheidet der Architect mit dem Projektinhaber (Lastenheft-Änderung = Vereinbarung); der Reviewer entscheidet das nicht.

**Übergabe:** F-1 und F-2 an den Implementer bzw. Architect (Rückkante Review → Plan); die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den Zähler. F-3 bis F-5 ohne Änderungspflicht. Der Report ersetzt keine Verifikation.
