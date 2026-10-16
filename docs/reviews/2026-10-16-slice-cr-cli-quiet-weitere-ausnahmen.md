# Review-Report: slice-cr-cli-quiet-weitere-ausnahmen — 2026-10-16

**Review-Art:** Plan/Design — Diff gegen den CR-Kern des Slice, [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), die angrenzenden Lastenheft-Zusagen und Hard Rules (Maintainability; nicht gegen die DoD).

**Gegenstand:** Commit `0052e1c` für den Slice [`slice-cr-cli-quiet-weitere-ausnahmen`](../plan/planning/done/slice-cr-cli-quiet-weitere-ausnahmen.md) (Welle `welle-cli-verbosity`): `spec/lastenheft.md` (Version 0.9.0, §3 `LH-FA-CLI-005`, §16-Zeile), `CHANGELOG.md` `[Unreleased]`, Streichung von „Nichts in Arbeit." in der Roadmap. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

**Skill:** `.harness/skills/reviewer.md` @ `3ec4535` (letzter Skill-Commit), HEAD `0052e1c` · <!-- d-check:ignore (Skill-Pfad als Zitat) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingangs-Kontext:**

- Slice [`slice-cr-cli-quiet-weitere-ausnahmen`](../plan/planning/done/slice-cr-cli-quiet-weitere-ausnahmen.md) (§1 CR-Kern, §2 DoD, §6 Risiken)
- [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-NFA-PERF-002`](../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker), [`LH-FA-INIT-005`](../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz)
- `AGENTS.md` (§3.4 Decken-Regel, §3.7), `harness/review.md`
- Sensors: `make docs-check` (256 Dateien, 0 Befunde). Nicht gelaufen: `make gates`, `make test`, `make lint` (reiner Dokumentations-Diff, kein Code berührt); `make test-docker`, `make verify-depguard` (kein Docker-Pfad, keine Importregel berührt).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Satz zu Ergebnis-Ausgaben steht im Lastenheft als normativ („bleiben“), ist aber keine der „Drei Ausnahmen“ und nennt kein Abgrenzungskriterium: „Zweck eines Befehls“ trägt auch `config set` (Zusammenfassungszeile), `add` („Changed:“), `init`/`generate` (Zusammenfassung) und `doctor` (der Bericht ist der Zweck). Szenario: Ein Implementer lässt die Erfolgszeilen von `config set`/`add`/`init` als „Ergebnis“ stehen oder lässt `doctor --quiet` Warn-Einträge zeigen; beide Lesarten sind vom Wortlaut gedeckt, und Abnahmekriterium 1 des Code-Slice („Erfolgs- und Zusammenfassungsmeldungen unterdrückt“) hat kein Vertragsmaß, gegen das es prüft. Der Slice delegiert die Grenze („trifft der Code-Slice im Prüfpunkt dort“) an Code-Ebene, obwohl sie am Vertrags-Stratum liegt. Zusätzlich wechselt die Begrifflichkeit zwischen „Erfolgs- und Zusammenfassungsmeldungen“, „Hinweise“ und „Status- oder Erfolgsmeldungen“. | [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-LESE-004`](../../spec/lastenheft.md#lh-lese-004--prüfbarkeit) | spec/lastenheft.md · „sind keine Status- oder Erfolgsmeldungen und bleiben.“ | nein — Review-only (nach Umsetzung: `make test` mit Pin je Subkommando) | Normsatz ohne Abgrenzungskriterium (Schlupfloch) |
| F-2 | MEDIUM | `LH-FA-CLI-005A` und `LH-FA-INIT-005` verlangen für `--force` eine Zusammenfassung der betroffenen Pfade („gibt immer eine Zusammenfassung der betroffenen Pfade aus“; „vor dem Schreiben muss eine Zusammenfassung der betroffenen Pfade ausgegeben werden“); CLI-005 lässt „Zusammenfassungsmeldungen der Subkommandos“ unter `--quiet` entfallen und benennt für diese Zusammenfassung keine Ausnahme. Szenario: Der Implementer unterdrückt sie unter `--force --quiet` (Überschreiben ohne Spur, verletzt 005A/INIT-005) oder behält sie (verletzt CLI-005); der Vertrag entscheidet den Vorrang nicht. Gleiches Muster, schwächer, bei „vor dem Entfernen warnen“ (Abhängigkeiten, `remove`) gegen „Alle anderen Hinweise entfallen“. | [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) | spec/lastenheft.md · „Unter `--quiet` entfallen die Erfolgs- und Zusammenfassungsmeldungen der Subkommandos.“ | nein — Review-only | Ausgabepflicht-Zusage ohne Vorrang gegen `--quiet` | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| F-3 | MEDIUM | Der Implementer-Eintrag im CHANGELOG behauptet den Ist-Stand des Codes, der gegen Lastenheft 0.9.0 verstößt: „`u-boot doctor --quiet` lists only error items … so warn items disappear too“ und „drop every warning under `--quiet` except the security-relaxation notices“ widersprechen Ausnahme 2 (Fallbacks in Befehlsausgabe und `doctor`); „the "no git remote", "has no effect", recreate and third-party notices no longer appear“ lässt init/generate/add/remove-Meldungen unerwähnt. Derselbe Abschnitt trägt direkt darüber den CR-Eintrag („Lastenheft 0.9.0“, „the implementation follows in a later release“); beide Einträge zusammen sagen zugleich „Fallbacks stay visible“ und „warn items disappear“. Szenario: Release-Cut vor dem Code-Slice `slice-v1-cli-quiet-nur-fehler` veröffentlicht beide Sätze; ein Nutzer liest „stays visible“, der Code unterdrückt Fallback-Warnungen. Dieselbe Aussage steht im Benutzerhandbuch (`doctor --quiet` zeigt nur Fehler-Einträge); das Handbuch ist nicht Teil des Diffs und gehört nach Abnahmekriterium 3 dem Code-Slice. | [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), harness/review.md MEDIUM „Doku/README/CHANGELOG driftet gegen Spec“ | CHANGELOG.md · „so warn items disappear too“ | nein — Review-only (Release-Cut-Prüfung) | CHANGELOG-Ist-Stand gegen Spec-Soll | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| F-4 | LOW | Die §16-Zeile 0.9.0 nennt die drei Ausnahmen und den Wegfall der Erfolgs- und Zusammenfassungsmeldungen, aber nicht den neuen Satz zu Ergebnis-Ausgaben (normativer Inhalt, der in 0.8.0 nicht stand). §16 definiert die Historie als „Fußabdruck angenommener Vertragsänderungen“; die Zeile bildet die Änderung unvollständig ab. | Lastenheft §16 (Historie-Regel, Zeilen 0.3.x–0.8.0 nennen jede inhaltliche Änderung) | spec/lastenheft.md · „Erfolgs- und Zusammenfassungsmeldungen aller Subkommandos entfallen unter `--quiet`; die Stufen und die Mehrfachangabe-Regel bleiben unverändert.“ | nein — Review-only | Historie-Zeile unvollständig gegen Anforderungstext |
| F-5 | INFO | Wortlaut-Abgleich Lastenheft gegen CR-Kern (§1): Der Anforderungstext ist Zeichen für Zeichen der CR-Kern, einschließlich des Ergebnis-Satzes (bewusst übernommen, siehe F-1). Die §16-Zeile weicht vom Entwurf ab, korrekt: „Compose-Fortschritt von `up`“ wird „Fortschritt von `u-boot up`“ (an Ausnahme 3 angeglichen), „bleiben“ wird „bleiben unverändert“ (0.8.0-Muster), der Entwurfs-Satz „Die geltende Zeile steht in spec/lastenheft.md §16“ entfällt zu Recht (er wäre ein Abwärts-Verweis). Der Slice-§1 trägt weiter „Entwurf“/„Vorschlag“ für einen angenommenen Text; Pflege bei der Closure (§7), Zuständigkeit Implementer. | — | docs/plan/planning/done/slice-cr-cli-quiet-weitere-ausnahmen.md · „CR-Kern (vorgeschlagene Fassung von“ | Review-only | Slice-Zustand hinkt dem angenommenen Text nach | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| F-6 | INFO | Stufenliste gegen Ausnahmen: „`--quiet` – nur Fehler“ bleibt als Stufe; die Ausnahmen stehen in einem eigenen Absatz mit „Drei Ausnahmen bleiben … sichtbar“ und „Alle anderen Hinweise entfallen“, also als Ausnahmen, nicht als vierte Stufe. Die Formulierung ist tragfähig, mit Ausnahme des Ergebnis-Satzes (F-1), der als vierte, nicht gezählte Ausnahme wirkt. Verträglichkeit mit [`LH-NFA-PERF-002`](../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer): ihre unbedingten Sätze („muss … sichtbar sein“, „jeder Fallback wird … ausgewiesen“) sind jetzt konsistent mit CLI-005; sie nennen `--quiet` nicht und brauchen es nicht. Ausnahme 2 nennt „Befehlsausgabe“ ohne Befehl; die Fallbacks entstehen teils beim Containerstart außerhalb des u-boot-Prozesses (Slice §6 trägt das als „entfallen“). | [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | spec/lastenheft.md · „Drei Ausnahmen bleiben unter `--quiet` sichtbar:“ | Review-only | — |
| F-7 | INFO | Versionsbezeichnung: Der CHANGELOG-Eintrag spricht von „Lastenheft 0.9.0“, der Abschnitt darunter ist das Produkt-Release „[0.9.0] - 2026-10-15“; zwei unabhängige Zählungen mit gleicher Nummer. Das Lastenheft 0.8.0 und der Eintrag davor trugen dasselbe Muster; kein Vertragsbruch. | Maintainability | CHANGELOG.md · „(`LH-FA-CLI-005`](spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Lastenheft 0.9.0)“ | Review-only | Doppelte Versionsnummer Lastenheft/Produkt | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Decken-Regel (`AGENTS.md` §3.4): Lastenheft-Diff enthält nur Verweise auf Lastenheft-Anker (`#lh-fa-dev-006--sandbox-profil`, `#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer`, `#lh-fa-dev-008--egress-restriktion`, `#lh-nfa-perf-002--startzeit-abhängig-von-docker`, `#lh-fa-cli-005--verbosity-und-logging`); alle Überschriften existieren; kein Verweis auf Slice, Carveout, Roadmap | geprüft, ohne Befund (`make docs-check`: 0 Befunde) |
| Widerspruch zu [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (Interaktivität) für Rückfragen: Bestätigungsfragen sind keine Erfolgs- oder Zusammenfassungsmeldung und bleiben unter `--quiet` unberührt; die Kollision liegt nur bei der `--force`-Zusammenfassung (F-2) | geprüft, ohne Befund (Rückfragen) |
| Suche `quiet`/„sichtbar“/„ausgewiesen“ im Lastenheft: weitere Zusagen unter `--quiet` stehen nur in `LH-FA-DEV-006` (Zeile „jede Sicherheitslockerung … ausweisen“), `LH-FA-DEV-007`, `LH-NFA-PERF-002`; alle sind durch die drei Ausnahmen gedeckt | geprüft, ohne Befund | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| Versionssprung 0.8.0 → 0.9.0: Versionsfeld, Historie-Zeile, Verweis „Vereinbarung mit dem Projektinhaber“, Datum 2026-10-16, Muster der Zeilen 0.3.x–0.8.0 | geprüft, ohne Befund (Inhaltslücke siehe F-4) |
| Sprachvertrag: Lastenheft deutsch, CHANGELOG englisch | geprüft, ohne Befund |
| `AGENTS.md` §3.7 im Lastenheft-Text und CHANGELOG: Indikativ über den Zustand, keine Chronik, keine verworfene Alternative | geprüft, ohne Befund |
| Commit-Nachricht: nennt `LH-FA-CLI-005`, `LH-NFA-PERF-002` und die Slice-Kennung | geprüft, ohne Befund | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| `docs/plan/planning/in-progress/roadmap.md`: Streichung „Nichts in Arbeit.“ bei gelisteter Welle | geprüft, ohne Befund |

Nicht geprüft:

- `docs/user/benutzerhandbuch.md` — nicht Teil des Diffs; die Drift ist in F-3 vermerkt, die Nachführung liegt nach Slice-Abnahmekriterium 3 beim Code-Slice.
- Code und Tests — der Commit ändert keinen; Umsetzungsprüfung gehört dem Review von `slice-v1-cli-quiet-nur-fehler`. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- DoD-Closure des CR-Slice — Aufgabe des Verifiers.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Normsatz ohne Abgrenzungskriterium (Schlupfloch) · Ausgabepflicht-Zusage ohne Vorrang gegen `--quiet` · CHANGELOG-Ist-Stand gegen Spec-Soll · Historie-Zeile unvollständig gegen Anforderungstext

## Verdikt

**Merge-blockierend:** ja — die drei MEDIUM-Findings sind vor Closure des CR-Slice zu klären (F-1 und F-2 am Vertrags-Stratum, F-3 vor jedem Release-Cut). Entscheidungen zum Ergebnis-Satz und zum Vorrang der `--force`-Zusammenfassung liegen beim Projektinhaber über den Architect (Lastenheft-Änderung ist eine Vereinbarung nach §16).

**Übergabe:** Findings gehen an den Implementer (F-1/F-2 zusätzlich an den Architect, Rückkante Review → Plan bei Plan-Defekt); die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den Zähler. Der Report ersetzt keine Verifikation.
