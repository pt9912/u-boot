# Review-Report: slice-cr-devcontainer-forge-cli — 2026-10-23

**Review-Art:** Plan/Design (Lastenheft-CR) — Wortlaut gegen Freigabe, Plan, Decken-Regel.

**Gegenstand:** Diff-Range `dd4d5e1a..HEAD` (Lastenheft 0.22.0, Welle `welle-forge-cli-gh-glab`).

**Skill:** `.harness/skills/reviewer.md` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-23

**Eingangs-Kontext:** Slice `slice-cr-devcontainer-forge-cli` (§1, §3, §6, §7 Freigabe), Entwurf `docs/messungen/slice-v1-forge-cli-quellen-messung/wortlaut-entwurf.md`, `spec/lastenheft.md` (Diff), `AGENTS.md` §3.4, §3.7, §3.10, §5 Regel 6, Welle-Datei, `open/README.md`, `slice-cr-lastenheft-messstaende`.

**Sensor:** `make gates` selbst gelaufen, Exit-Code 0 (Log im Scratchpad). `make test-docker`, `make verify-depguard`: nicht gelaufen, kein Code-, Docker- oder Importregel-Pfad im Diff.

## Satz-für-Satz-Vergleich Lastenheft gegen Freigabe (§7 Freigabe, Absatz) und Entwurf

- Spiegelstriche `GitHub CLI (gh)` / `GitLab CLI (glab)`: wortgleich zum Entwurf.
- S1: Halbsatz „in beiden Profilen (Standard und Sandbox)“ vor „als eingebaute Features … wählbar“. Lastenheft = §7-Absatz wortgleich. Position bewertet: sachgerecht (bindet die Wählbarkeit an beide Profile, nicht nur den Schlüssel). Weitere Abweichungen vom Entwurf siehe F-1.
- S2, S3 (bis auf Ort), S5, S6: Lastenheft = §7-Absatz wortgleich; S6 verdichtet „Nicht zugesagt sind … sowie …“ zu „Nicht zugesagt: …; außerdem …“ ohne Zusageverschiebung (siehe F-1 zu „Aufrufe mit gültigem Token“).
- S4 hostnamenfrei: „stehen die Hosts von `gh` in der Default-Allowlist und kommt für `glab` ein Host hinzu“ ist genau der Entwurf ohne Klammer und ohne `gitlab.com`; trägt Antwort Frage 2, bewertet: richtig, Ollama-Absatz ist das Muster.
- Stelle 2 (drei Passagen): wortgleich zum Entwurf-*Neu*.
- Version 0.22.0, Datum 2026-10-23.

## Findings

### F-1 MEDIUM
- `quelle`: AGENTS.md §3.10 (Aussage über den Plan/Freigabe prüfen), Maintainability
- `pfad`: `docs/plan/planning/in-progress/slice-cr-devcontainer-forge-cli.md` · „S1 mit dem Halbsatz umgestellt, sonst wortgleich“
- `klasse`: freigabe-abweichung-nicht-ausgewiesen
- `befund`: §7 nennt als einzige Abweichung vom Entwurf den Halbsatz (Punkt 1) und die Hostnamen (Punkt 2). Im Lastenheft weichen zusätzlich ab: S1 „`u-boot config set devcontainer.features.gh.enabled true`, analog `glab`“ statt des zweiten ausgeschriebenen Befehls; S2 ohne „im Container“ (Entwurf: „installiert sein Werkzeug im Container aus …“); S6 „Aufrufe mit gültigem Token“ statt „Aufrufe, die ein gültiges Token brauchen“. Die Zeile „sonst wortgleich“ ist damit unzutreffend. Ob der Projektinhaber den verdichteten Wortlaut (der dem §7-Absatz entspricht) gesehen hat, ist aus dem Repo nicht ablesbar.
- `risiko`: Die Freigabe-Spur belegt die Wortgleichheit nicht; der Verifier vergleicht gegen §7 und der Entwurf bleibt unverändert daneben stehen. Inhaltlich klein (S2 Ollama-Absatz hat „im Container“, hier fehlt es), aber eine nicht ausgewiesene Abweichung bei einer Vertragsänderung.
- `verifizierbar`: Review-only (Rückfrage Hauptagent: war der verdichtete Wortlaut vorgelegt? sonst §7 berichtigen, die drei Abweichungen benennen).

### F-2 LOW
- `quelle`: AGENTS.md §3.10 (Beleg-Tabelle / Test-Auftrag je Bedingung)
- `pfad`: Slice §7 · Beleg-Tabelle Zeilen B5, B6, B9
- `klasse`: bedingung-test-auftrag-luecke
- `befund`: Drei Test-Aufträge decken die Bedingung nicht ganz: B5 pinnt nur „fehlende Signatur“, das Lastenheft sagt „weist **das** aus“ (Signatur *und* nicht geprüfte `gh`-Attestierung); B6 „`gh` fügt der Basis-Liste nichts hinzu“ pinnt nicht, dass die `gh`-Hosts in der Default-Allowlist stehen (Lastenheft: „stehen … in der Default-Allowlist“); B9 testet `remoteEnv`-Referenz, die Lastenheft-Zusage „weder im Image noch im Workspace-Volume“ für Forge-Token hat keinen eigenen Auftrag. Alle übrigen Bedingungen (B1–B3, B7, B10, B11) haben Auftrag; B4/B8 begründet ohne Test; die Aufträge stehen auch in Welle §5.
- `risiko`: Der Code-Slice kann grün schließen, während ein Teil der Zusage ungepinnt bleibt (§3.10, „Bedingung ohne Test“, 18 Belege).
- `verifizierbar`: Review-only.

### F-3 LOW
- `quelle`: AGENTS.md §3.10 (Zeilenangaben im Plan nach Änderung nachziehen)
- `pfad`: `docs/plan/planning/open/slice-cr-lastenheft-messstaende.md` · „Zeile 409, `LH-FA-DEV-003` (Absatz `ollama`)“
- `klasse`: plan-aussage-ueber-code
- `befund`: Der Diff zieht die Versionsangaben jenes Slice nach, nicht aber seine Zeilenangaben: Ollama-Absatz steht jetzt in Zeile 411 (Plan: 409), `podman`-Zeilen 467/469 (Plan: 463/465), Historie 0.12.0/0.15.0 Zeilen 1588/1591 (Plan: 1568/1571; 0.7.0 und 0.13.0 analog). Die Verschiebung durch die vier neuen Zeilen vor DEV-007 ist von diesem Diff verursacht, die Historie-Angaben waren schon vorher veraltet.
- `risiko`: Der Implementer des Messstände-CR liest eine Zeile, die nicht mehr stimmt.
- `verifizierbar`: Review-only (`grep -n` im Lastenheft).

### F-4 INFO
- `quelle`: Maintainability
- `pfad`: `spec/lastenheft.md` · „hat die Messung keine Signatur gefunden“
- `klasse`: messung-im-vertragstext
- `befund`: Der freigegebene Absatz nennt „die Messung“ und „Suchmuster“ im Vertragstext (Ollama-Absatz nennt Messstand im selben Abschnitt). Das ist freigegeben und wird vom Folgeslice `slice-cr-lastenheft-messstaende` gegebenenfalls erfasst; kein Befund gegen den Diff.
- `risiko`: keines für diesen Diff; Hinweis für den Messstände-CR.
- `verifizierbar`: Review-only.

## Historie-Zeile 0.22.0 (selbst formuliert) — bewertet
Nennt nur freigegebene Stellen, trägt Frage-1-Antwort („in beiden Profilen“), keine Hostnamen, keine Kennung nach unten, kein Messwert, Anker-Links auf Lastenheft-Abschnitte; Muster der Zeile 0.21.0. Ohne Befund. Die Wendung „ausgewiesen beim Erzeugen“ ist knapp, aber verständlich.

## Geprüft, ohne Befund
- `spec/lastenheft.md` Decken-Regel (AGENTS.md §3.4): keine Verweise auf ADR, Slice, Welle, Roadmap im Diff; Links nur auf Lastenheft-Anker (existieren, `make docs-check` grün).
- Stelle 2 wortgleich; Stelle 3 (keine Änderung an DEV-006/-008/NFA-SEC-004) eingehalten, diese Anforderungen im Diff unberührt.
- Commit-Reihenfolge: Freigabe (`3b2179c7`) vor Lastenheft (`5b55250f`), Moves als reine `git mv`-Commits (`9789f71d`, `e1b98398`), Lifecycle konsistent (Slice nur in `in-progress/`).
- Plan-Tabelle gegen `git diff --stat`: 6 Dateien (roadmap, Slice-Move, `open/README.md`, Messstände-CR, Welle, `spec/lastenheft.md`), jede im Plan, keine ohne Diff; `docs/reviews/` und `observations/` laut Plan erst später.
- Versionsnachzüge Messstände-CR und `open/README.md`: „nach 0.22.0“ / „zuletzt 0.22.0“ korrekt (Zeilenangaben siehe F-3).
- §5 Regel 6: `docs/user/`, `README.md`, `README.de.md`, `CHANGELOG.md` im Diff unberührt.
- Kommentar-/Chronik-Muster (§3.7, §3.9): keine Konjunktiv-/Früher-Sätze in den neuen Texten gefunden; Roadmap-Zelle nennt Zustand mit Anker.

## Nicht geprüft
- Wahrheit der Messaussagen A–F gegen Rohdaten (Verifier/Validator; Messprotokoll ist Gegenstand des Vorgänger-Slice).
- Ob der Projektinhaber den verdichteten Wortlaut gesehen hat (nicht im Repo belegt, siehe F-1).
- Ausgänge der Risiken §6, DoD, Closure-Notiz: Verifier-Aufgabe; Nachlauf-Review nötig, falls eine Fix-Runde Lastenheft-Text ändert (F-1).
