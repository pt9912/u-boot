# Review-Report: slice-cr-lastenheft-messstaende — 2026-10-25

**Review-Art:** Plan-/Spec-Review (Lastenheft 0.23.2 auf 0.24.0: Zusage-Menge statt Messstand in LH-FA-DEV-003 und LH-FA-DEV-007) mit Abgleich nach unten (Handbuch, READMEs, CHANGELOG, Messprotokoll). Kein Produktcode. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

**Gegenstand:** `git diff 0d50bfd7..HEAD` (10 Dateien, 120 Einfügungen, 25 Löschungen): Commits `9ffbfda4` (Start und Freigabe-Eintrag), `e2ea8398` (nur `spec/lastenheft.md`), `939c3739` (Abgleich nach unten). Move-Commits `e06b7d81`/`0d50bfd7` gelesen, rein.

**Skill:** `.harness/skills/reviewer.md` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-25

**Eingangs-Kontext:** Slice (§1 Wortlaut-Entwurf, Messlage M1 bis M5, Herkunftstabelle, Fragen; §3 Plan; §6; §7 Freigabe, Abgleich, Vorab-Suche), `spec/lastenheft.md`, Messprotokolle der Vorläufer (per Slice-Tabelle), Vorreport `2026-10-24-slice-cr-lastenheft-ollama-cloud-modelle.md`, `AGENTS.md` §3.4, §3.7, §3.10, §5 Regel 6. Das CI-Log des Belegs M5 wurde selbst gelesen (`gh api repos/pt9912/u-boot/actions/jobs/121082471598/logs`, 11608 Zeilen). <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

**Sensoren:** `make gates` grün (Exit 0); `make doc-trace`: 84 Anforderungen, 0 Waisen; `make doc-commits RANGE="origin/main..HEAD"`: 7 Commits, 0 Befunde. `make test-docker` und `make verify-depguard` nicht gelaufen: Diff ohne Code, Docker-Pfad und Importregeln nicht berührt. Der Wortvergleich lief per `grep -F` mit Shell-Ersetzung auf dem Host (reiner Text, kein Gate-Nachweis).

**Ergebnis:** kein HIGH. Ein MEDIUM, zwei LOW, vier INFO. Der Wortlaut im Lastenheft ist mit dem freigegebenen Wortlaut Zeichen für Zeichen identisch.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der CHANGELOG-Eintrag sagt „This resolves the "Limit of the Lastenheft" of `0.16.0` (two sets for `ollama`)“. Das Limit von `0.16.0` lautet: Das Lastenheft führt zwei Mengen (Installation; Cloud-Modelle und Host-Anschluss), „unifying them is a later change“. Das Lastenheft `0.24.0` führt weiter zwei Mengen: Installation `amd64` und Colima (`arm64`), Cloud-Weg und Host-Anschluss nur `amd64`; Zeile 414 weist die Trennung aus („von der Zusage zur Installation … getrennt“). Es vereinheitlicht nichts. | Doku driftet gegen Spec (Skill MEDIUM); AGENTS.md §3.10 (Aussage über Zustand ungelesen); BEO-DOK/zusage-weiter-als-messung | `CHANGELOG.md` · „This resolves the "Limit of the Lastenheft" of `0.16.0` (two sets for `ollama`)“; Gegenstelle `CHANGELOG.md` · „unifying them is a later change of the Lastenheft“ | Review-only | aussage-ueber-zustand-weiter-als-text |
| F-2 | LOW | Das Messprotokoll listet unter „Belegt“ Loopback-Bindung, `library=cpu`, Pull hinter der Egress-Restriktion und das Überleben des Modells als vom Pass des Tests belegt. Im Log steht für `TestE2E_SandboxDevcontainer_Ollama` nur die Zeile `--- PASS` (42,97 s) ohne die Marker `LOOPBACK_ANSWERS`, `LIBRARY_CPU`, `PULL_OK`; die Aussagen stammen aus den Assertions im Testcode (`sandbox_devcontainer_docker_test.go`, Prüfliste um Zeile 1307 und Zeile 1324). Das Protokoll sagt nicht, dass es aus dem Testcode ableitet. | AGENTS.md §3.10 (Fundstelle nennen); `harness/verification.md` Evidence-Form | `docs/messungen/slice-cr-lastenheft-messstaende/README.md` · „**Belegt** ist der Pass von `TestE2E_SandboxDevcontainer_Ollama` (Bau des erzeugten Dockerfiles, Server auf der Loopback-Adresse“ | Review-Only (Log und Testcode gelesen) | beleg-ableitung-nicht-ausgewiesen | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| F-3 | LOW | Das Protokoll nennt „Die drei `SKIP`-Zeilen“ der Podman-Tests. Das Log hat vier `--- SKIP`: zusätzlich `TestE2E_SandboxDevcontainer_StatusFiles_VFS` (`sandbox_status_docker_test.go:173`, gleiche Meldung zu `…=1`). Der Auszug ist als „betroffene Tests“ gefiltert, das Protokoll sagt nicht, dass ein vierter Skip existiert, während `README.md` und `README.de.md` („the four tests of nested Podman are skipped there“) vier nennen. | Maintainability; Doku-Präzision | `docs/messungen/slice-cr-lastenheft-messstaende/README.md` · „Die drei `SKIP`-Zeilen tragen im Log dieselbe Meldung“ | Review-only | zaehlung-nicht-gegen-log | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| F-4 | INFO | Wortlaut-Vergleich ohne Abweichung: Stellen 1, 2, 3, 4 (beide Halbsätze, Zeilen 414 und 418) und die Historie-Zeile `0.24.0` stehen im Lastenheft Zeichen für Zeichen wie in Slice §1 *Neu* (Datei-Links durch Anker, `JJJJ-MM-TT` durch `2026-10-25` ersetzt); Versionsfeld `0.24.0`; Klassen B, C, D und die Historie `0.23.x` unberührt (Diff zeigt nur 7 Einfügungen, 6 Löschungen an Zeilen 9, 412, 414, 418, 472, 474, 1607). Die Freigabe (`9ffbfda4`, 09:06:52) liegt vor der Lastenheft-Änderung (`e2ea8398`, 09:07:33), `9ffbfda4` enthält keinen Lastenheft-Diff. Ob die vorgelegte Vorlage den Entwurf ungekürzt zeigte, ist aus dem Repo nicht prüfbar (Gespräch). | BEO-DOK/vertragsaenderung-nach-review, freigabe-vorlage-verdichtet-den-entwurf | `spec/lastenheft.md` · „Zugesagt ist das auf einem Linux-Host mit lokaler Docker-Engine (`amd64`) und unter Colima (`arm64`)“ | `git log`, `grep -F` | — |
| F-5 | INFO | Die neue Zusage „`podman` startet … auf einem Linux-Host mit lokaler Engine (`amd64`)“ und `podman build` auf `amd64` stützen sich auf zwei Messungen mit Wert `…=0` (M1 A6, M4), nicht auf die CI: Der CI-Runner hat `…=1`, alle drei Podman-Tests skippen (im Log gelesen). Handbuch, README und CHANGELOG weisen das aus; die Zusage ist ohne CI-Gate. Zuständigkeit: Verifier (Evidence), Projektinhaber (Entscheidung ist am 2026-10-25 gefallen). | BEO-DOK/zusage-weiter-als-messung | `spec/lastenheft.md` · „Bei `podman` startet das Ergebnis unter Docker auf einem Linux-Host mit lokaler Engine (`amd64`) und unter Colima (`arm64`)“ | Review-only | zusage-ohne-ci-beleg |
| F-6 | INFO | Die Zusage für Colima (`arm64`) bei `ollama` steht auf einem von Hand nachgeschriebenen Probe-Image (M2), nicht auf dem Erzeugnis. Der Slice führt das als bekannte Einschränkung (§6, Herkunftssatz 1.6), das Handbuch nennt es („Unter Colima (`arm64`) ist ein von Hand nachgeschriebenes Probe-Image gemessen, der Bau des erzeugten Dockerfiles nicht“); der Lastenheft-Wortlaut sagt dazu nichts Neues. Keine Ausweitung, nichts zu ändern. | BEO-DOK/messvariante-belegt-nicht-die-auslieferung | `docs/user/benutzerhandbuch.md` · „nachgeschriebenes Probe-Image gemessen, der Bau des erzeugten Dockerfiles nicht“ | Review-only | — |
| F-7 | INFO | Offener Fund aus Vorab-Suche §7, nicht in diesem Slice: `LH-FA-DEV-003` Absatz `gh`/`glab` enthält „hat die Messung keine Signatur gefunden“ (Messbefund im Anforderungstext, Frage 8) und der Kommentar vor `TestE2E_SandboxDevcontainer_Ollama` nennt einen durch den CI-Pass erledigten Zustand. Beides ist im Slice benannt und für Folgeschritte gemeldet. | Maintainability | `spec/lastenheft.md` · „hat die Messung keine Signatur gefunden“ | Review-only | — | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

---

## Prüfung nach Schwerpunkten

1. **Wortlaut gleich freigegeben:** gleich (F-4).
2. **Kennungen und Messstände im Lastenheft-Text:** keine ADR-, Slice-, Welle-, SPEC- oder ARC-Kennung in den geänderten Zeilen; keine Engine-Version, kein Hostname, kein Messbild; Verweise sind Anker. „Docker Hub (`docker.io`)“ und „Colima“ sind Zusage-Umgebungen, keine Messstände.
3. **Deckung der Zusagen:** amd64-Installation: M1 (zwei Läufe) plus CI-Pass (im Log gelesen: `PASS TestE2E_SandboxDevcontainer_Ollama` und `…_PinMismatchBreaksBuild`, Run 40164437737, Job 121082471598, Commit `f0d85e5b`, Runner-Image `ubuntu-24.04` Version 20261012.334.1; Angaben des Protokolls stimmen). Colima/arm64 bleibt (F-6). `podman build` nur `docker.io` (M3, M4 beide `docker.io/library/…`). Podman-Start: beide Umgebungen mit Erzeugnis gemessen, nicht in CI (F-5). Widersprüche im Lastenheft: keiner gefunden; Absatz `ollama` (Colima zugesagt) und Zeilen 414/418 (Colima nicht zugesagt) sind durch den Satz „sie ist von der Zusage zur Installation … getrennt“ ausgewiesen; Zeile 474 `docker`-Satz und Zeile 507 unverändert; die Suche über Akzeptanzkriterien und Befehlsübersicht (Zeilen mit `ollama`, `podman`, `nestedRuntime`) zeigt keinen Messstand. Die Spannung der zwei Mengen bleibt bestehen (F-1).
4. **Abgleich nach unten:** Handbuch (Zeilen um 1382, 2158 bis 2175, 2461 bis 2473), beide READMEs und CHANGELOG tragen keine Kennung (README, `docs/user`: keine `LH-`/`ADR-`/`SPEC-`-Treffer im Diff); „amd64 nicht gemessen“ und „not in CI“ sind am Log berichtigt: „CI bestand den Docker-Test der Installation“ und „CI übersprang die Podman-Tests wegen `…=1`“ stimmen mit dem Log überein. Die Gegenstelle zum CHANGELOG-Limit ist F-1.
5. **Messprotokoll:** Ablage nach Ablage-Regel (`docs/messungen/<Kennung>/README.md`, kein Skript, keine Rohdatei außer Auszug), Index `docs/messungen/README.md` ergänzt, keine Geheimnisse, Run/Job/Commit/Runner stimmen mit dem Log; Auszug zeilentreu (F-2, F-3).
6. **Plan-Tabelle gegen `git diff --stat`:** alle zehn Dateien im Plan oder in §7 geführt (`roadmap.md` über den Startvermerk, `open/README.md` über „Zeile 33 und Zeile 73“, `docs/messungen/README.md` als nachträglicher Eintrag genannt). `carveouts.md` unverändert, Auflösung laut Frage 1d bei der Closure (Zustandsmarker in `open/README.md`: „Review, Verifier und Closure stehen aus“, im Ist-Stand korrekt).
7. **Gates:** siehe oben.

## Negativbefunde („geprüft, ohne Befund“)

- `spec/lastenheft.md` Diff (Wortlaut, Kennungsfreiheit, Anker, Historie, Version): ohne Befund.
- Freigabe-Reihenfolge per `git log`: ohne Befund.
- `README.md` und `README.de.md` (Aussagen gegen Log; Kennungsfreiheit): ohne Befund.
- `docs/user/benutzerhandbuch.md` (neue Sätze gegen Messung/Log): ohne Befund.
- `docs/plan/planning/open/README.md`, `roadmap.md` (Zustandsmarker, keine Chronik): ohne Befund.

## Nicht geprüft

- Gespräch mit dem Projektinhaber (Vorlage ungekürzt?): nicht aus dem Repo prüfbar.
- Englisch/Deutsch-Parität von README.de.md im Detail: nur die geänderten Zeilen auf Aussagegleichheit gelesen.
- `make test-docker`, `make verify-depguard`: kein Code, kein Importbezug.
