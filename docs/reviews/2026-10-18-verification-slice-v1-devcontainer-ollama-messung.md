# Verification-Report: slice-v1-devcontainer-ollama-messung — 2026-10-18 <!-- d-check:ignore (Slice-Kennung im Titel, Zitat) -->

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice `slice-v1-devcontainer-ollama-messung` (`docs/plan/planning/in-progress/`), Commit-Range `6acb1c0..HEAD` (8 Commits, HEAD `2289928`, 49 Dateien), [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) (`Proposed`). <!-- d-check:ignore (Kennungen als Zitat) -->

**Eingang:** Review `docs/reviews/2026-10-18-slice-v1-devcontainer-ollama-messung.md` (F-1 MEDIUM, F-2..F-9 LOW, F-10..F-12 INFO), Fix-Runde `0430356`, Folge-ADR `2289928`. <!-- d-check:ignore (Kennungen als Zitat) -->

**Gesamturteil: GO mit Auflage.** Die Messung und ihre Dokumentation tragen; kein Review-Finding ist offen (F-9 bewusst). Der Slice ist **nicht schliessbar**, bis der Projektinhaber [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) auf `Accepted` setzt (Haken Folge-ADR und Uebergabe). Auflagen siehe Abschnitt 6. <!-- d-check:ignore (Kennung als Zitat) -->

## 1. Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass, Exit 0 (Exit-Code von `make` selbst erfasst) | Lauf im Verifier-Kontext |
| `make doc-trace` | pass, Exit 0 | „83 Anforderung(en), 0 Waise(n).“ |
| `make doc-immutable RANGE="6acb1c0..HEAD"` | pass, Exit 0 | Range aufgeloest, 8 Commits; „408 Datei(en) geprueft, 0 Befund(e)“; im ADR-Diff keine geloeschte Zeile, keine Accepted-ADR veraendert |
| `make test-docker`, `make verify-depguard` | nicht gelaufen | `git diff --stat 6acb1c0..HEAD -- internal spec docs/user README.md CHANGELOG.md Makefile` ist leer: kein Produktcode, keine Importregeln, kein Docker-Pfad des Produkts |
| Neuer Lauf der Messskripte | nicht gelaufen | braucht Colima, Netz und `gh` mit Anmeldung auf dem Host; Rohdaten gelesen (Docker-only, `AGENTS.md` 3.1) |

## 2. Plan-vs-Diff (Plan-Tabelle §3 gegen `git diff 6acb1c0..HEAD --stat`)

| Diff-Gruppe | In Tabelle | Befund |
|---|---|---|
| ADR-0024 und ADR-Index | Zeile 1 | ok; Index-Zeile 52 vorhanden | <!-- d-check:ignore (Kennung als Zitat) -->
| `BEO-APP/podman-static-pin-alterung` (Evidence neu, `state.md`) | Zeile 2 | ok |
| `open/slice-cr-…`, `open/slice-v1-devcontainer-feature-ollama`, `open/slice-v2-spec-…` | Zeile 3 und Zeile 8 (Verweis-Nachzug) | ok | <!-- d-check:ignore (Kennung als Zitat) -->
| `docs/reviews/` (nur Review-Report; dieser Verification-Report folgt mit diesem Commit) | Zeile 4 | ok; „Verification Evidence folgt“ wird mit diesem Commit eingeloest |
| Slice-Datei (Move `open` nach `in-progress`) | Zeile 5 | ok |
| Messprotokoll-Verzeichnis, 32 Dateien (13 Skripte, 3 Dockerfiles, 15 Rohdaten, README) | Zeile 6 | Zaehlung stimmt (`ls | wc -l` = 32) |
| `BEO-WZG/host-werkzeug-ausserhalb-docker` (Evidence neu, `state.md`) | Zeile 7 | ok |
| `roadmap.md` | Zeile 8 | ok (Arbeitszeile; Rueckstellung auf „Nichts in Arbeit.“ bei Closure ist Planner-Aufgabe) |
| `welle-sandbox-egress-ollama.md`, `open/README.md`, `done/slice-v1-egress-wildcard-pruefung.md`, `BEO-DOK/messprotokoll-…/state.md` | Verweis-Nachzuege | ok | <!-- d-check:ignore (Kennung als Zitat) -->
| `internal/**`, `spec/**`, `docs/user/`, Lastenheft | „unveraendert“ | durch leeren Diff bestaetigt |

Befund V-3 (LOW, Zeilentext): Plan-Tabelle Zeile 1 nennt „Wortlaut-Entwurf“; die ADR traegt die Bedarfsliste (im Umsetzungsstand ausgewiesen als Abweichung). Zeilentext bei der Closure angleichen.

## 3. DoD-Abgleich

| DoD-Zeile | Stand | Beleg |
|---|---|---|
| Messung O1..O7 „gemessen, soweit ausgewiesen“ | **erfuellt, der Haken traegt** | Wortlaut ist die engere Fassung (Review F-5): Option C nur Manifest, Option E nicht, O6 (ii) mit `docker exec` nachgestellt, O6 (iii) nicht gesondert, `amd64` ohne Lauf, jeweils in §3 und der Liste „Nicht gemessen“ benannt. Stichproben gegen Rohdaten: Pins (Version 0.40.6, beide SHA-256) in `rohdaten-08.txt`, `Dockerfile.b`, `Dockerfile.probe`, `rohdaten-02/-12` identisch mit ADR-0024 Entscheidung 1; Groessen (1 438 364 593 amd64, 1 560 414 226 arm64, vulkan 40 275 720, Image 2 182 601 467, Probe-Image 402 069 870) stimmen zu Slice und ADR; Messstand-Digest `6ee2b409…` in README und ADR gleich; Ausnahme-Messprotokoll im Repo (Wiederholbarkeit bis auf Netzquellen). | <!-- d-check:ignore (Kennung als Zitat) -->
| Folge-ADR | **offen, zu Recht** | ADR-0024 `Proposed`; fuenf Freigabe-Punkte a bis e stehen beim Projektinhaber | <!-- d-check:ignore (Kennung als Zitat) -->
| Uebergabe | **offen, zu Recht** | Folge-Slices gegen den Entwurf berichtigt (Diff vorhanden), Haken folgt mit `Accepted`; Beobachtung `BEO-APP/podman-static-pin-alterung` hat zweiten Beleg |
| `make gates` gruen | erfuellt | Abschnitt 1 |
| Review, Report liegt vor | erfuellt, Haken noch nicht gesetzt | Report vorhanden, alle Findings behandelt (Abschnitt 4); Haken ist Closure-Schritt |
| Doku-Update ADR-Index | erfuellt im Inhalt | Zeile 52; Haken folgt |
| Closure-Notiz §7, Steering-Loop-Eintrag | offen | §7 „Wird bei der Closure gefuellt.“ |
| Beobachtungs-Register | teils | zwei Evidence-Dateien und zwei `state.md` fortgeschrieben (Zaehlung 2 von 3 bei Pin-Alterung und Host-Werkzeug); keine Schwelle erreicht; Verknuepfung in §7 steht aus |
| Risiken §6 mit Ausgang | im Wortlaut „Stand“ gefuellt, „Ausgang: setzt die Closure“ | mit Closure auf eingetreten/entfallen/offen setzen (Risiko Lieferkette haengt an der Entscheidung zu ADR-0024) | <!-- d-check:ignore (Kennung als Zitat) -->
| Drei Paarungen | offen | Welle-Closure |

## 4. Review-Findings F-1..F-12

| Finding | Stand | Verifier-Beleg |
|---|---|---|
| F-1 MEDIUM vorsignierte URL | **behoben** | `grep -rniE 'X-Amz|Authorization|bearer|token|PRIVATE KEY|password|secret|/Users/|ghp_|gho_'` ueber das Messprotokoll: kein `X-Amz`, kein Schluessel, keine Home-Pfade; Treffer nur Skripte (`ghcr.io/token` wird zur Laufzeit geholt, nicht gespeichert), Meldung „Generating new private key“ und `docker-credential-secret get` im Build-Log (ohne Inhalt). `inner-05.sh` Zeile 56 schwaerzt die Query kuenftig (`?[Query entfernt: vorsignierte URL]`); README benennt die Schwaerzung. Die URL bleibt im Git-Verlauf (Commit `47575bd`); nach Review: oeffentliches Modell, Ablauf 2026-10-19 16:58 UTC, kein Konto-Token. |
| F-2 Groessenmetrik | behoben | Metrik `.Size` und Nicht-Einzelfaktor in Slice §3 und ADR |
| F-3 Zaehlfehler Releases | behoben | „sieben Releases … nicht nachweislich der Anfang“ |
| F-4 Zustandsmarker | behoben | Tabelle Zeile 4 und §6 berichtigt; §6-Ausgaenge „setzt die Closure“ |
| F-5 Haken-Wortlaut | behoben | Haken „soweit ausgewiesen“, siehe DoD |
| F-6 Deutungen | behoben | Codesigning als „nicht belegt/Deutung“, Praefix als „Deutung“, Volume-Lauf mit Hinweis auf `rohdaten-10.txt` |
| F-7 Probe als Erzeugnis | behoben | `Dockerfile.probe` Z. 3 und README („von Hand nachgeschrieben“) |
| F-8 Host-Python | behoben | im Umsetzungsstand ausgewiesen, zweiter Beleg `BEO-WZG/host-werkzeug-ausserhalb-docker` |
| F-9 `Verantwortlich:` | **bewusst offen, begruendet** | Slice §4 „Ist-Stand des Starts“ (open direkt nach in-progress auf Beauftragung, Verantwortlicher setzt der Projektinhaber). Akzeptiert; setzt der Projektinhaber das Feld nicht, bleibt `—` und ist in §7 festzuhalten |
| F-10 Vorschlag breiter als Messung | behoben | Quellenwahl-Vorschlag klammert GPU-Auslassung aus (Entwurfsfrage 5) |
| F-11 Hygiene Rest | ok | vier ephemere oeffentliche `ssh-ed25519`-Schluessel (nur oeffentlich, je Wegwerf-Container), kein privater Schluessel |
| F-12 Wiederholbarkeit | behoben, Rest LOW (V-4) | README nennt `Dockerfile.cli` im Voraussetzungsteil (Z. 18) und `06-`/`07-`/`09-…sh` nennen es im Kopf; die Reihenfolgetabelle fuehrt es nur bei Schritt 6 |

## 5. ADR-0024 (Proposed)

- **Status:** `Proposed`, Index-Zeile `Proposed`; keine Accepted-ADR veraendert (`doc-immutable` gruen). ADR-0020 und -0023 bleiben unveraendert. <!-- d-check:ignore (Kennung als Zitat) -->
- **Pin-Werte gegen Rohdaten:** Version, beide SHA-256, Download-Groessen, Archivinhalte, Digest des Basis-Images stimmen (Abschnitt 3).
- **Aussagen am Code:** gelesen und gefunden: `featureCatalogue`/`featureFor` (`devcontainer_features.go:251`), `projectFeatureSourceVersion` (Z. 413), `projectAllFeatureEntries`, `classifyFeatureEntries` (`doctor.go`), `maybeWarnOrphanFeatureActivation` (`config.go:319`), `catalogueFeatureActive` (`devcontainer_sandbox.go:339`), `TestFeatureCatalogue_EntriesAreValid`; ADR-0015, -0022 und `SPEC-045` existieren. Folgerungen zu `projectFeatureSourceVersion` (Marker -> `ok=false`) sind als Entwurf der Aenderung formuliert, nicht als Ist-Verhalten. <!-- d-check:ignore (Kennungen als Zitat) -->
- **Nicht-Verifiziertes ausgewiesen:** eigener Absatz „Nicht verifiziert“ und Liste „Nicht gemessen“ (Generator-Dockerfile, `amd64`-Lauf und `--exclude`-Muster, `vulkan`, Auto-Start mit Dev-Container-CLI, Stabilitaet des Blob-Hosts, Neubau-Volume, `OLLAMA_NO_CLOUD`, Signierschritte ausserhalb des Suchmusters). Grenze „kein Signaturschluessel, nur Pin gegen spaeteren Austausch“ steht in Entscheidung 3 und in den Konsequenzen.
- **Entscheidungsform:** jede Entscheidung nennt Gegenoption; fuenf Freigabe-Punkte a bis e liegen beim Projektinhaber (Empfehlung, keine Entscheidung). Konform mit der Rollentrennung: kein Lastenheft-, Spec- oder Produkttext.
- **Hinweis V-1 (INFO, Architect-Entscheidung):** Entscheidung 6 (CUDA-Ausschluss) ist als Zusage formuliert, obwohl nur `cuda_*` auf `arm64` gemessen ist; das ist in der Entscheidung selbst ausgewiesen, der Umsetzungs-Slice muss `amd64` beim Bau messen (Entscheidung 10).
- **Hinweis V-2 (INFO):** Entscheidung 4 (kein Auto-Start) haelt eine Folge-Entscheidung offen; `Accepted` schliesst sie nicht.

## 6. Auflagen fuer die Closure

1. **Closure erst nach `Accepted` von ADR-0024** durch den Projektinhaber (Haken Folge-ADR und Uebergabe; der Slice ist bis dahin nicht schliessbar). Die Entscheidung zur Lieferkette (Freigabe-Punkt a) ist laut §5 im Gespraech zu belegen und in §7 zu vermerken. Wird eine Option der Gegenseite gewaehlt (Auto-Start, Wildcard, vollstaendig entpacken), ist ADR und Folge-Slice-Plan neu zu pruefen. <!-- d-check:ignore (Kennung als Zitat) -->
2. Nach `Accepted`: Plan-Tabellen der drei Folge-Slices gegen den endgueltigen ADR-Stand pruefen (Hakensatz „Stand 2026-10-18“ in DoD ersetzen), Slice-Plan Zeile 1 („Wortlaut-Entwurf“ -> Bedarfsliste, V-3).
3. §7 fuellen (Messergebnis, Quelle, ADR-Nummer und -Status, Entscheidung, Steering-Loop-Eintrag, Register, Folge-Slices); alle §6-Risiken auf eingetreten/entfallen/offen setzen; Abgleich Plan-Tabelle gegen `git diff --stat` nach der **letzten** Aenderung wiederholen (`AGENTS.md` 3.10; dieser Report ist die letzte Datei der Range).
4. Roadmap auf „Nichts in Arbeit.“ zuruecksetzen; `git mv` nach `done/` als reiner Move nach der Closure-Notiz (`AGENTS.md` 3.3).
5. Haken Review und Doku-Update (ADR-Index) setzen; Register-Zaehlung (Pin-Alterung, Host-Werkzeug) bleibt 2 von 3, keine Verkoerperung.
6. Optional (V-4): Reihenfolgetabelle des Messprotokoll-README um `Dockerfile.cli` fuer Schritte 7 und 9 ergaenzen.

## 7. Nicht verifiziert

Ob der Upstream ausserhalb der gelesenen Orte signiert; Stabilitaet des Blob-Hosts; `amd64`-Lauf; Wiederholung der Messung (Netz und `gh` noetig); Python-Einmalaufruf auf dem Host (Aussage des Implementers, nicht aus dem Repo ableitbar).
