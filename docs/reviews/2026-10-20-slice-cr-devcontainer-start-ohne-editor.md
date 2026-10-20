# Review-Report: [`slice-cr-devcontainer-start-ohne-editor`](../plan/planning/in-progress/slice-cr-devcontainer-start-ohne-editor.md) — 2026-10-20

**Review-Art:** Lastenheft-CR (Vertragstext, kein Produktcode). Maßstab ist der freigegebene, unverändert zu übernehmende Wortlaut in [`slice-v1-devcontainer-start-ohne-cli`](../plan/planning/done/slice-v1-devcontainer-start-ohne-cli.md) §7 (Stellen 1 bis 6), [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) und [`ADR-0029`](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md) (beide `Accepted`).

**Gegenstand:** `git diff 9e6933b..69280df -- spec/ docs/plan/planning/` und `git diff --stat 9e6933b..HEAD` (4 Commits, `HEAD` = `69280df`).

**Skill:** `.harness/skills/reviewer.md` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-20

**Selbst gelesen:** der gesamte Diff von `spec/lastenheft.md`; der Wortlaut Stellen 1 bis 6 im Eltern-Slice §7 (Stelle 1 mit `diff` gegen den Lastenheft-Text verglichen, Linkformen normalisiert); der Slice (voll); Lastenheft-Stellen [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), -006, -007, [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), -010, [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), §2 Betriebsumgebung, [`LH-NFA-PORT-004`](../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image), LH-MVP, §16; der Befehl-Slice [`slice-v1-cli-devcontainer-up-down`](../plan/planning/open/slice-v1-cli-devcontainer-up-down.md) (§2, §3 Tests S1, S3, S4, S8, S9) und der Fundament-Slice [`slice-v1-devcontainer-start-fundament`](../plan/planning/open/slice-v1-devcontainer-start-fundament.md) (Tests S5, S6); die Diffs der übrigen Planungs-Dateien.

**Sensor-Läufe des Reviewers:** `make gates` Exit 0; `make doc-trace` (`84 Anforderung(en), 0 Waise(n)`, [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) Zeile `ok`); `make doc-immutable RANGE="9e6933b..HEAD"` (4 Commits, `d-check: 555 Datei(en) geprüft, 0 Befund(e)`). Nicht gelaufen: `make test-docker` und `make verify-depguard` (kein Docker-Pfad, keine Importregel berührt).

---

## Findings

Keine HIGH-, keine MEDIUM-Findings.

```text
LOW docs/plan/planning/in-progress/slice-cr-devcontainer-start-ohne-editor.md · "Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt" - Plan-Tabelle nennt nicht alle berührten Dateien
Klasse: plan-tabelle-ohne-diff-abgleich
Quelle: AGENTS.md §3.10
Befund: git diff --stat 9e6933b..HEAD berührt zwölf Dateien; die Tabelle führt vier Zeilen (lastenheft.md, roadmap.md, open/README.md, docs/reviews/). Ohne Eintrag: [`open/slice-cr-lastenheft-messstaende.md`](../plan/planning/open/slice-cr-lastenheft-messstaende.md), welle-devcontainer-start-ohne-editor.md, open/slice-v1-cli-devcontainer-up-down.md, open/slice-v1-devcontainer-start-fundament.md, open/slice-v2-spec-devcontainer-start.md, done/slice-v1-devcontainer-start-ohne-cli.md, der Slice selbst (Verantwortlich) sowie die beiden bestehenden Reports unter docs/reviews/ (Linknachzug durch slice-mv; die Zeile nennt "neu (Reviewer, Verifier)"). Die Zeile roadmap.md sagt "bei der Closure", der Diff berührt sie schon beim Start.
Risiko: Der Abgleich gilt nach der letzten Fix-Runde; solange er aussteht, ist es kein Vertragsbruch, aber die Tabelle bildet den Diff heute nicht ab, und der Verifier liest sie als Maßstab.
Verifizierbar: Review-only (git diff --stat gegen §3)
```

```text
LOW docs/plan/planning/open/slice-v1-cli-devcontainer-up-down.md · "`down` entfernt den Container nicht (`docker inspect` kennt ihn, `Running` falsch)" - "entfernt weder den Container noch seine Volumes" ohne Test für die Volumes nach down
Klasse: spec-bedingung-ohne-test
Quelle: AGENTS.md §3.10 (Bedingung im Wortlaut braucht Test oder Streichung)
Befund: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) sagt "`u-boot devcontainer down` stoppt den Container; er entfernt weder den Container noch seine Volumes." Test S4 pinnt den Container nach down; das benannte Volume mit Marker wird nur über `--recreate` geprüft, nicht über `down`. Auch `--recreate` bleibt "benannte Volumes bleiben erhalten" nur im Sandbox-Profil gepinnt (das Standardprofil hat keins).
Risiko: Eine Änderung, die bei down Volumes mitentfernt, bleibt grün; der Lastenheft-Satz hat dann keinen Beleg. Schwach, weil `docker stop` Volumes nicht berührt.
Verifizierbar: Review-only (Prüfauftrag an den Befehl-Slice, S4 um einen Marker-Fall nach `down` ergänzen oder den Halbsatz im Folge-CR streichen)
```

```text
INFO spec/lastenheft.md · "`egress`, `devcontainer` oder `template`" - Wert `devcontainer` im Enum, Verhalten bei --dry-run/--diff steht nicht im Lastenheft
Klasse: enum-wert-ohne-vorschau-aussage
Quelle: [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run)
Befund: [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run) verlangt Dry-Run für "dateiverändernde Befehle"; die neuen Befehle ändern keine Dateien, der Befehl-Slice weist `--dry-run`/`--diff` mit Exit 2 ab. Das Lastenheft sagt dazu nichts (wie bei `egress`, gleiches Muster, gleicher Wortlaut-Typ).
Risiko: keines gegen den freigegebenen Wortlaut; Hinweis für den Spezifikations-Slice (SPEC-001).
Verifizierbar: Review-only
```

```text
INFO spec/lastenheft.md · "Fehlerklassen ([`LH-FA-CLI-006`](#lh-fa-cli-006--exit-codes)): Exit-Code `10`" - Exit 14 nicht genannt
Klasse: wortlaut-fehlerklassen-teilmenge
Quelle: [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
Befund: Der Befehl-Slice nennt 2/10/11/12/14; [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) nennt 10, 11, 12. 14 (technischer Dateisystemfehler) und 2 (Usage) gelten über [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) allgemein.
Risiko: keines; der Wortlaut ist freigegeben und weicht nicht ab.
Verifizierbar: Review-only
```

```text
INFO doc-trace · "[`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | – Devcontainer ohne Editor starten | —" - ADR-Spalte leer
Klasse: rtm-adr-spalte-leer
Quelle: [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) (Accepted, immutable, AGENTS.md §3.5)
Befund: Die Kennung bestand beim Accept von [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) nicht; die Spalte bleibt leer. Der Slice-Beleg trägt die Zeile (`ok`).
Risiko: keines; die ADR wird nicht umgeschrieben. Hinweis, falls eine Folge-ADR das Schärft-Feld setzt.
Verifizierbar: make doc-trace
```

```text
INFO docs/plan/planning/in-progress/slice-cr-devcontainer-start-ohne-editor.md · "Freigabe des Projektinhabers** im Gespräch belegt und in §7 vermerkt" - Nachlauf-Review und Freigabe-Vermerk stehen aus
Klasse: vertragsaenderung-nach-review
Quelle: BEO-DOK/vertragsaenderung-nach-review
Befund: §7 ist leer (Closure steht aus). Jede Änderung von spec/lastenheft.md nach diesem Report braucht einen Nachlauf-Review über den Diff. Zuständig für DoD-Closure: Verifier.
Risiko: siehe Slice §6, dritter Punkt.
Verifizierbar: Review-only
```

---

## Prüfung nach Auftrag

1. **Wortlaut-Treue.** Stelle 1: der Text von [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) ist wortgleich zum Vorschlag (`diff` nach Normalisierung der Linkform: einzige Abweichung ist der Anker `#lh-fa-up-001--umgebung-starten` statt Pfad plus Anker, erlaubt). Stelle 2: beide Zeilen, Priorität V2, wortgleich. Stelle 3: (a) und (b) wortgleich, [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) in der Verweisliste von Ausnahme 1. Stelle 4: `devcontainer` zwischen `egress` und `template`, wortgleich. Stelle 5: Zusatz wortgleich. Stelle 6: Versionsfeld 0.18.0 und Historie-Zeile wortgleich zum Vorschlag (Datum 2026-10-20, Spalte "Vereinbarung mit dem Projektinhaber"). Keine inhaltliche Abweichung.
2. **Zielbild-Regel.** Die hinzugefügten Zeilen in `spec/lastenheft.md` enthalten keinen Verweis auf ADR, Slice, Carveout, Roadmap oder Welle, keinen Messstand und keinen Prozesszustand ("Linux-Host `amd64`" ist eine Zusage-Menge). Gesucht über den gesamten Diff-Plus-Teil.
3. **Konsistenz im Lastenheft.** Gelesen, kein Widerspruch: [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) ("keinen Container-Runtime-Socket des Hosts einbinden") trägt jetzt die Abgrenzung des Hilfs-Containers; `up`/`down` (Compose) und `devcontainer up`/`down` sind in der `--quiet`-Liste, in [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) und in [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) ("startet nicht den Compose-Stack") getrennt benannt; die Zeile `u-boot up` der MVP-Liste und [`LH-FA-UP-001`](../../spec/lastenheft.md#lh-fa-up-001--umgebung-starten)/-004 bleiben unberührt und eindeutig; "weist aus" (Socket) steht in Ausnahme 1, die Fallback-Ausweisung in Ausnahme 2, beide bleiben unter `--quiet`; Verb "soll" bei Priorität V2 entspricht [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (V2, "soll"); [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Klassen 10, 11, 12 bestehen); §2 "sekundäre Betriebsumgebungen" (macOS, Windows) und "Docker Engine oder kompatible Docker-Laufzeit" stehen nicht im Widerspruch, DEV-010 nennt seine engere Menge ausdrücklich; [`LH-NFA-PORT-004`](../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image) (Image ohne Shell) stützt "im ausgelieferten Container-Image läuft er nicht"; Sprachvertrag [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache) und Traceability (§13) nicht berührt; [`LH-ABG-002`](../../spec/lastenheft.md#lh-abg-002--keine-ide-abhängigkeit) gestützt.
4. **Version und Verweise.** 0.17.0 → 0.18.0, die Nummer war frei (letzte Historie-Zeile 0.17.0). Die Verweise in [`slice-cr-lastenheft-messstaende.md`](../plan/planning/open/slice-cr-lastenheft-messstaende.md) (Titel, Frage 2, DoD "zuletzt `0.18.0`"), `open/README.md` und `welle-devcontainer-start-ohne-editor.md` lesen "nächste freie Version nach 0.18.0" beziehungsweise nennen 0.18.0 als vergeben. Suche nach `0.18.0` im Planungsbaum: keine Stelle behauptet mehr, 0.18.0 sei für den Messstände-CR vorgesehen. Stimmt.
5. **Plan-Tabelle gegen `git diff --stat`.** Lücken, siehe das erste LOW-Finding.
6. **Bedingungen im Wortlaut gegen Tests.** Zugeordnet: "Fallback-Ausweisung setzt `u-boot.yaml` voraus" (Befehl-Slice S1 Fall (a)); "ohne diese Datei Exit 12 statt 11" (S1 Fall (b)); "Exit 11 bei Startscript-Abbruch, sonst 12" (S8: `onUnavailable: fail` → 11, `postStartCommand: exit 3` → 12; Unit-Fall `result=failed` → 11, sonst 12); "Docker nicht erreichbar oder Engine nicht lokal → 11" (S8: Binary nicht auf dem `PATH`, `DOCKER_HOST=tcp://…`); "im ausgelieferten Container-Image läuft er nicht" (S1 (c) über den S8-Fall); "kein devcontainer.json → 10" (S8); "nie in Datei oder Image" (Fundament S6 Umgebung); "gestarteter Devcontainer erhält den Socket nicht" (Fundament S6); "erster Aufruf braucht Netzzugriff" und "gepinnte Quellen" (Fundament S5); "`--recreate` neuer Container, Volumes bleiben" (S4, nur Sandbox-Volume). Nicht zugeordnet: "down entfernt keine Volumes" (siehe zweites LOW-Finding).

## Geprüft ohne Befund

- `spec/lastenheft.md` Stellen 1 bis 6 gegen den freigegebenen Vorschlag (Wortlaut-Treue)
- Zielbild-Regel (AGENTS.md §3.4) im Lastenheft-Diff
- Konsistenz [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), -006, -007, [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), -007, -008, [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), [`LH-NFA-PORT-004`](../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image), [`LH-ABG-002`](../../spec/lastenheft.md#lh-abg-002--keine-ide-abhängigkeit)
- Versionsfeld, Historie-Zeile und Nachzug der Versionsverweise (messstaende-Slice, `open/README.md`, Welle)
- Referenzmodell: keine Verweise von Rang 1 abwärts; die Links in den Planungs-Dateien zeigen auf `in-progress/` (nach dem Move), `make docs-check` im Gate grün
- ADR-Immutabilität: `make doc-immutable` ohne Befund; [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md)/-0029 nicht berührt
- Dual-Classifier: nicht anwendbar (kein Code)

## Nicht geprüft

- `make test-docker`, `make verify-depguard`: kein Code-, Docker- oder Importregel-Pfad im Diff.
- Die Tragfähigkeit der Zusage (Linux `amd64`, Standardprofil mit UID ≠ 1000) gegen reale Läufe: Aufgabe der Docker-Tests des Befehl-Slice, nicht dieses CR.
- Spezifikation und Architektur-Sicht: im Diff nicht berührt, Nachführung liegt beim Spezifikations-Slice.
- DoD-Closure des Slice (Verifier-Aufgabe): Freigabe-Vermerk in §7, Closure-Notiz, Steering-Loop-Eintrag.
