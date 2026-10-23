# Review-Report (Nachlauf 2): slice-v1-forge-cli-quellen-messung, Closure `8d475164..HEAD` — 2026-10-23

**Rolle:** Reviewer (Modul 10, [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md), [`harness/review.md`](../../harness/review.md)), frischer Kontext. Prüfung gegen Plan, Entscheidungen und Hard Rules, nicht gegen die DoD (Verifier).

**Gegenstand:** Commit-Range `8d475164..HEAD` (22 Dateien): Closure von [`slice-v1-forge-cli-quellen-messung`](../plan/planning/done/slice-v1-forge-cli-quellen-messung.md); [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) von `Proposed` auf `Accepted`; Wortlaut-Entwurf; Welle-Datei §4/§5; Roadmap; Register. Vorangehende Reports: [Erst-Review](2026-10-23-slice-v1-forge-cli-quellen-messung.md), [Nachlauf](2026-10-23-slice-v1-forge-cli-quellen-messung-nachlauf.md), [Verification](2026-10-23-slice-v1-forge-cli-quellen-messung-verify.md).

**Eingang gelesen:** Diff der ADR (vollständig), des Wortlaut-Entwurfs, der Welle-Datei, der Roadmap, des `open/README.md`, der Register-Dateien (`state.md`, `evidence/`); Slice §3, §6, §7 am Ruheort; Rohdaten-Zuordnung der `arm64`-Pins per `grep`; `AGENTS.md` §3.5, §3.7, §3.10.

**Sensor:** `make gates` selbst ausgeführt (Docker): **Exit-Code 0**. `make test-docker`, `make verify-depguard`: nicht gelaufen (kein Go-Code, keine Importregel im Diff).

## Ergebnis

Kein HIGH, kein MEDIUM. Zwei LOW, zwei INFO. Die ADR gibt die Entscheidungen vom 2026-10-23 richtig wieder; die Änderung verengt die Zusage, eine neue unbelegte Aussage über Code kam nicht hinzu. Befunde an der Accepted-ADR wären nach `AGENTS.md` §3.5 eine Folge-ADR; F-2 ist ein solcher Fall und ist INFO (siehe dort).

## Prüfung nach Fragen

1. **ADR gegen die Entscheidungen.** a und b (Release-Archiv, SHA-256, `gh` und `glab`): Entscheidung 1 und Tabelle, Herkunft „Projektinhaber“. c (`GH_TOKEN`/`GITLAB_TOKEN` getrennt): Entscheidung 5, Herkunft „Projektinhaber“. d (beide Profile): Entscheidung 3, „Projektinhaber“. f (nur Gemessenes, keine Nachmessung): Tabelle und Konsequenzen, „Projektinhaber“; der Trigger „Nachmessung (Punkt f)“ ist zu „spätere Messung“ umformuliert, konsistent. h (`glab duo`, `glab orbit`, Nachladen nicht zugesagt): Entscheidung 4, Tabelle, Konsequenzen, Trigger. e und g: „Standard aus der Empfehlung“, Definition der Spalte *Herkunft* steht vor der Tabelle; Kontextabsatz und Geschichte-Zeile nennen dieselbe Aufteilung (a bis d, f, h gefragt; e, g Standard). Welle: `v0.15.0` bestätigt, „Forge vor Ollama“ bestätigt (Welle §5 Fragen 1 und 3). Dass e und g nicht gefragt wurden und dass die Antworten am 2026-10-23 so fielen, ist Aussage des Auftrags; der Reviewer hat sie nicht an einer Quelle außerhalb von Auftrag und Dokumenten prüfen können (INFO I-1). Ohne Befund bei der Wiedergabe.
2. **Verengung oder neue Aussagen (§3.10).** Neu und über die Entscheidungen hinaus ist nur der `arm64`-Satz in Entscheidung 1: „stehen in keiner Rohdatei als gelesene Ausgabe (nur im Dockerfile-Block und im Bau-Log von `glabtgz`; beim `gh` zusätzlich in `rohdaten-03.txt`)“. Gegengeprüft: `7f322df3…` (gh arm64) steht in `block-ghtgz.dockerfile`, `rohdaten-03.txt` und `rohdaten-06-ghtgz.log`; `0716395a…` (glab arm64) steht in `block-glabtgz.dockerfile` und `rohdaten-06-glabtgz.log`. Die Aussage trifft zu (genauer wäre „und im Bau-Log von `ghtgz`“, kein Befund). „Live nachgelesen am 2026-10-23 (V-4)“ verweist auf den Verification-Report und ist dort belegt. Sonst entfallen oder verengen die Änderungen: „Vorschlag“, „Empfehlung“ und „offen“ sind getilgt (Restsuche in der ADR: nur die Definition der Spalte und Alternative G „N2/N3 offen“, die eine Tatsache zum Weg `source:` ist). Die Zusage zu `LH-NFA-SEC-004` wird enger, nicht weiter.
3. **Wortlaut-Entwurf.** Ausdrücklich „nicht freigegeben“ in der Voraussetzung und im Schlussabschnitt; der CR-Slice legt ihn vor. S6, Auslegung zu `LH-NFA-SEC-004` (Stelle 3) und S4-Zeile nennen dieselbe Menge wie ADR Punkt h: `glab duo`, `glab orbit`, Nachladen weiterer Programme durch `glab`. Eine Ungenauigkeit in der ADR selbst: F-2.
4. **Closure-Notiz, Register, Paarungen.** §6: genau ein Ausgang je Risiko (1 bis 5). §7: Risiken-Zeile, Drei Paarungen (kein `liegt in`, kein Folge-Slice-Link, sechs Register-Verzeichnisse mit nicht leerem `evidence/`) stimmen mit dem Baum überein. Register: je Eintrag eine Evidence-Datei mit **einem** Vorgang (`slice-v1-forge-cli-quellen-messung`), `state.md` nennt jeweils „Ausgang wird nicht zugewiesen“ beziehungsweise „unverändert“, keine Zählerzeile gesetzt; bei `messprotokoll-rohdaten-mit-zugangsdaten` nur eine Notiz im `state.md` ohne neue Evidence-Datei, die §7 und die Notiz selbst ausdrücklich so führen. Roadmap: „Nichts in Arbeit.“ und die Offen-Zeile ohne Slice-Zeiger („fünf Slices“) stimmen mit der Welle-Tabelle (fünf Zeilen) überein; `open/README.md`: reiner Verweis-Nachzug.
5. **Nutzer-Doku und Spec-Straten.** `git diff --stat`: weder `spec/`, `docs/user/`, `README*.md`, `CHANGELOG.md` noch `internal/` berührt. Kennungen in den Rängen 1 bis 3 und in der Nutzer-Doku: nicht berührt. Die ADR verweist aufwärts (Lastenheft, Spezifikation), nicht abwärts auf Slice oder Roadmap (der Verweis auf Messprotokoll und Reports ist Beleg-Kante, kein Slice-Verweis).

## Findings

### F-1 — LOW — Zustandsmarker `Proposed` im Messprotokoll
- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt wird beim Ablauf auf den Ist-Stand gesetzt)
- `pfad`: `docs/messungen/slice-v1-forge-cli-quellen-messung/README.md` · „Die Entscheidungsvorlage zur Quelle ist [`ADR-0035`]… (`Proposed`)“
- `klasse`: zustandsmarker-nach-ablauf
- `befund`: Der Link im Kopf wurde auf `done/` nachgezogen, der Status der ADR steht weiter als `Proposed` und „Entscheidungsvorlage“; die ADR ist `Accepted`. Der Wortlaut-Entwurf im selben Ordner wurde berichtigt.
- `risiko`: Ein Leser des Protokolls hält die Quellenwahl für offen. Kein Gate sieht es.
- `verifizierbar`: Review-only (`grep -n Proposed`).

### F-2 — INFO — Reichweite des Nachladens in ADR Entscheidung 6
- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.5 (Accepted-ADR immutable), ADR-0035 Punkt h
- `pfad`: `docs/plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md` · „das Nachladen weiterer Binaries durch `glab duo` und `glab orbit` ist von der Zusage ausgenommen“
- `klasse`: formulierung-weicht-von-entscheidung-ab
- `befund`: Entscheidung 4 und Punkt h nehmen „`glab duo`, `glab orbit` und das Nachladen weiterer Programme durch `glab`“ aus; Entscheidung 6 bindet das Nachladen an die beiden Unterbefehle. Der Wortlaut-Entwurf folgt Entscheidung 4.
- `risiko`: Gering, weil Entscheidung 4 die engere Menge der Zusage ist und der CR-Slice dem Entwurf folgt. Eine Berichtigung wäre eine Folge-ADR, keine Umschreibung; sie lohnt allein nicht und kann mit der nächsten Folge-ADR zu ADR-0035 mitlaufen.
- `verifizierbar`: Review-only.

### F-3 — LOW — Plan-Tabelle beschreibt den Roadmap-Zustand vor der Closure
- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.10 (Abgleich der Plan-Tabelle nach der letzten Runde)
- `pfad`: `docs/plan/planning/done/slice-v1-forge-cli-quellen-messung.md` · „Roadmap: Zeiger unter *Offene Wellen* statt „Nichts in Arbeit.““
- `klasse`: plan-tabelle-stand-vor-closure
- `befund`: Der Diff setzt die Roadmap auf „Nichts in Arbeit.“ und entfernt den Slice-Zeiger; die Zeile der Plan-Tabelle nennt die Gegenrichtung.
- `risiko`: Ein Leser der Closure-Notiz liest einen Roadmap-Zustand, den der Baum nicht hat.
- `verifizierbar`: Review-only (Abgleich mit `git diff 8d475164..HEAD -- docs/plan/planning/in-progress/roadmap.md`).

### F-4 — INFO — Host-Werkzeug-Nutzung des Planners
- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.1; Register `BEO-WZG/host-werkzeug-ausserhalb-docker`
- `pfad`: `docs/plan/planning/observations/BEO-WZG/host-werkzeug-ausserhalb-docker/evidence/slice-v1-forge-cli-quellen-messung.md` · „Selbstauskunft“
- `klasse`: host-werkzeug-ausserhalb-docker
- `befund`: Selbstauskunft des Planners: `python3 --version` auf dem Host und ein Heredoc im Scratchpad mit `sed -i`. Kein Nachweis des Handoffs stützt sich darauf; der Gate-Nachweis ist `make gates`, hier erneut mit Exit 0 gelaufen. Das ist keine Gate-Lockerung (kein HIGH nach dem Skill); die Versionsabfrage ist ohne Wirkung, `sed -i` und Heredoc im Scratchpad berühren keine Repo-Datei. Das Muster (Heredoc statt Write-Datei bei Commits) steht im Register und ist dort mit dieser Closure als weiterer Beleg geführt. Die Einordnung folgt dem Register: eine Gelegenheit im Vorgang, Ausgang liegt beim Projektinhaber.
- `risiko`: Ein Wiederholungsmuster, kein Schaden in diesem Diff.
- `verifizierbar`: Review-only; eine Aussage des Laufs, an keiner Datei nachprüfbar.

### I-1 — INFO — Herkunft der Entscheidungen
Dass die Punkte a bis d, f, h am 2026-10-23 so entschieden wurden und e, g nicht gefragt wurden, stützt sich auf den Auftrag und die Dokumente, nicht auf eine unabhängige Quelle; die Spalte *Herkunft* macht die Unterscheidung für Leser sichtbar. Zuständig für eine Bestätigung ist der Projektinhaber.

## Negativbefunde (geprüft, ohne Befund)
- `docs/plan/adr/0035-…md` und `docs/plan/adr/README.md`: Wiedergabe der Entscheidungen, Status, Index-Zeile, Geschichte-Zeile; keine Abwärtsverweise auf Slice oder Roadmap.
- `docs/messungen/…/wortlaut-entwurf.md`: Nicht-Freigabe, S4/S6, Stelle 3 konsistent mit Punkt h.
- `docs/plan/planning/welle-forge-cli-gh-glab.md`, `in-progress/roadmap.md`, `open/README.md`: Tabelle (ein Code-Slice), Fragen 1 bis 4, Ruhe-Marker.
- Slice §6/§7 und Register: Ausgänge, Paarungen, ein Vorgang je Evidence-Datei, keine Zähler, keine zugewiesenen Ausgänge.
- Kommentare und Zustandsfelder (§3.7/§3.9) im Diff: keine „früher“-/„ohne X wäre“-Kommentare; Zustandsfelder nennen Zustand und Beleg.
- Rang 1 bis 3 und Nutzer-Doku: nicht berührt.

## Nicht geprüft
- `make test-docker`, `make verify-depguard`: kein Go-Code, keine Importregel im Diff.
- Live-Gegenprüfung der `arm64`-Pins gegen die Prüfsummen-Dateien im Netz: nicht wiederholt, liegt beim Verification-Report V-4.
- Messwerte A bis F des Protokolls: nicht neu gerechnet (unverändert seit dem Verification-Report).
