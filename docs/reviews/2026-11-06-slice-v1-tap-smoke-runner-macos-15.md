# Review-Report: slice-v1-tap-smoke-runner-macos-15 — 2026-11-06

**Review-Art:** Code-Review gegen Plan und Konventionen: Diff gegen Slice-Plan, ADR, Lastenheft und Hard Rules. Der Diff enthält eine CI-Workflow-Vorlage für das Tap-Repo, keinen Go-Code.

**Gegenstand:** Commit `737dae18` (HEAD, Basis `026fefb5`), eine Datei, 4 Einfügungen, 1 Löschung: `packaging/homebrew-tap/.github/workflows/smoke.yml`.

**Skill:** `.harness/skills/reviewer.md` @ `3ec45359`. **Modell:** claude-sonnet-5-5. **Datum:** 2026-11-06.

## Eingangs-Kontext und Sensors

**Gelesen:**

- Slice [`slice-v1-tap-smoke-runner-macos-15`](../plan/planning/in-progress/slice-v1-tap-smoke-runner-macos-15.md): Anlass, §1 (Planentscheidung `macos-15` statt `macos-latest`, Nicht-Ziele), §2 (DoD), §3 (Plan-Tabelle), §4 bis §6, §7, §8.
- Lastenheft: [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung) (Paketierung: Homebrew, eigener Tap).
- [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) (Entscheidung 5: Smoke im Tap; nennt kein Runner-Label).
- Die Datei `packaging/homebrew-tap/.github/workflows/smoke.yml` im Endstand (vollständig gelesen) und `git show HEAD`.
- `AGENTS.md` §3.1, §3.2, §3.6, §3.7, §3.9, §3.10, §5; `harness/review.md`; `harness/conventions.md`; die letzten Reports in `docs/reviews/`.

**Sensors:**

- `git diff 026fefb5 HEAD --stat`: genau die eine Datei, 4 Einfügungen, 1 Löschung. Der Diff besteht aus der Matrix-Zeile und einem dreizeiligen Kommentar davor (DoD „Quelle“).
- Suche nach `macos-14`, `macos-15` und `macos-latest` über den Baum (ohne `.git/` und vendorte Baseline): Treffer außerhalb der Slice-Datei sind `CHANGELOG.md` Zeile 454 (Abschnitt `[0.10.0]`, Stand dieses Release) und Abschlussnotizen unter `docs/plan/planning/done/` (Stand des jeweiligen Slice). Keine Treffer in `docs/user/`, `docs/maintainer/`, `README.md`, `README.de.md`, `packaging/` außer der Datei selbst, `scripts/`, `.github/workflows/`.
- Nicht-ASCII-Zeichen in `smoke.yml`: nur Zeile 6 (Gedankenstrich im bestehenden Kommentar); der neue Kommentar ist ASCII.
- `make docs-check`: nach Anlage dieses Reports gefahren, Ergebnis siehe Abschnitt *Sensor-Lauf* am Ende.
- `make gates`, `make test-docker`, `make verify-depguard`: nicht gefahren. Der Diff enthält keinen Go-Code, keine Importregel und keine Datei, die ein Gate oder ein Test liest (Plan-Tabelle, Zeile „Tests“; die Suche oben bestätigt, dass keine Datei außer der Vorlage `smoke.yml` nennt). `make gates` im Devcontainer-Container ist Closure-Pflicht des Implementers und Nachweis des Verifiers.
- Lauf `smoke` im Tap und Vergleich Quelle gegen Tap-Datei: nicht gefahren (kein Gegenstand des Commits; der Push in den Tap folgt nach Review und Verification, DoD „Tap“ und „Smoke“).

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der Plan führt Zeilennummern von `smoke.yml` im Stand vor dem Diff: „Zeile 24“ für die Matrix, „Zeilen 58 bis 66“ für das Audit-Gate. Im Endstand steht die Matrix-Zeile auf 27 und der Schritt „Formel pruefen (gatend)“ ab Zeile 61; die drei Kommentarzeilen verschieben alles dahinter. | `AGENTS.md` §3.10 (Abgleich nach der letzten Änderung); [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) | [`slice-v1-tap-smoke-runner-macos-15`](../plan/planning/in-progress/slice-v1-tap-smoke-runner-macos-15.md) §1 · „`smoke.yml` Zeilen 58 bis 66“ und §Anlass · „Zeile 24 (`runner: [macos-14, macos-15-intel]`)“ | nein — Review-only (Zeilen in `smoke.yml` gegen die Plan-Aussagen gelesen) | Zeilenverweis im Plan nach dem Diff veraltet |
| F-2 | LOW | Die Plan-Tabelle begründet den Kommentar-Wortlaut mit „ASCII wie die übrigen Kommentare der Datei (Zeilen 3 bis 8)“. Zeile 6 enthält einen Gedankenstrich (`—`, nicht ASCII). Der neue Kommentar selbst ist ASCII. | `AGENTS.md` §3.10 (Aussage über Code gelesen, nicht angenommen) | [`slice-v1-tap-smoke-runner-macos-15`](../plan/planning/in-progress/slice-v1-tap-smoke-runner-macos-15.md) §3 · „Wortlaut ASCII wie die übrigen Kommentare der Datei (Zeilen 3 bis 8)“ | nein — Review-only (Suche nach Nicht-ASCII-Zeichen in `smoke.yml`) | Plan-Aussage über Code ungenau |
| F-3 | LOW | §7 des Plans sagt „der Slice liegt in `open/` (Entwurf)“; die Datei liegt in `in-progress/`, und der Diff setzt den Plan um. Der Zustandsmarker ist abgelaufen. | `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt wird entfernt oder auf den Ist-Stand gesetzt) | [`slice-v1-tap-smoke-runner-macos-15`](../plan/planning/in-progress/slice-v1-tap-smoke-runner-macos-15.md) §7 · „Noch nicht geschlossen; der Slice liegt in `open/` (Entwurf).“ | nein — Review-only | Abgelaufener Zustandsmarker im Plan |
| F-4 | INFO | Dass `macos-15` von GitHub als arm64-Runner bedient wird und dass die Stilllegung von `macos-14` die Ursache des roten Laufs ist, folgt aus den im Plan genannten Quellen (README von `actions/runner-images`, Ankündigung, Annotation), nicht aus einem Lauf. Dieser Review las die Quellen nicht nach; der Plan benennt den Lauf mit dem neuen Label als Beleg (§6, DoD „Smoke“). | `AGENTS.md` §3.10; [`BEO-DOK/plan-annahme-ueber-fremdzustand-ohne-messstelle`](../plan/planning/observations/BEO-DOK/plan-annahme-ueber-fremdzustand-ohne-messstelle/observation.md) | `smoke.yml` · „runner: [macos-15, macos-15-intel]“ | ja — Lauf `smoke` im Tap mit Job-Metadaten (`labels`, `runner_name`) und „Set up job“; Zuständigkeit Verifier-Nachlauf | Annahme über Fremdzustand bis zur Messung |
| F-5 | INFO | Der Kommentar nennt die Kopplung „Aenderung dieser Zeile in Vorlage und Tap“. Kein Skript und kein Gate hält Quelle und Tap-Datei gleich; der Plan führt das als offenes Risiko (§6, Datei-Drift) mit Frage an den Projektinhaber. | [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) Entscheidung 5 (Smoke im Tap); Slice §6 | `smoke.yml` · „Zeile in Vorlage und Tap, kein wanderndes macos-latest.“ | nein — Review-only; `cmp` Quelle gegen Tap-Datei ist DoD „Tap“ (Verifier) | Kopplung zweier Orte ohne Sensor |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `packaging/homebrew-tap/.github/workflows/smoke.yml`, Matrix-Zeile | geprüft, ohne Befund: `runner: [macos-15, macos-15-intel]` entspricht der Planentscheidung (arm64 ausdrücklich `macos-15`, Intel unverändert, kein `macos-latest`); das YAML bleibt gültig (Kommentarzeilen mit gleicher Einrückung wie `runner:`, innerhalb von `matrix:`) |
| `smoke.yml`, übrige Datei | geprüft, ohne Befund: Trigger (`push` auf `Formula/u-boot.rb`, `workflow_dispatch`), `permissions: contents: read`, alle Schritte und das gatende `brew audit --strict --online` sind unverändert; keine Lockerung nach `AGENTS.md` §3.6, kein `|| true`, kein `continue-on-error` |
| Kommentar im Diff (`AGENTS.md` §3.7, §3.9) | geprüft, ohne Befund: Indikativ; trägt die Klassen Abgrenzung (`macos-latest` wandert, die beiden Labels nennen die Version) und Kopplung (Vorlage und Tap); keine Chronik, kein „früher“, kein „ohne X wäre“; kein Verweis auf ein Symbol, das der Diff umbenennt; Suche nach dem alten Wert `macos-14` im Baum fand nur Stände älterer Releases und Slices (siehe Sensors) |
| Plan-Tabelle gegen `git diff --stat` (`AGENTS.md` §3.10) | geprüft, bis auf F-1 bis F-3 ohne Befund: der Commit berührt nur die eine Datei, die der Plan als „update“ führt; die Zeilen „unverändert“ (`.github/workflows/`, `scripts/tap-nachzug.sh`, `scripts/tap-check.sh`, `packaging/homebrew-tap/README.md`, `docs/maintainer/releasing.md`, `CHANGELOG.md`, [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md)) sind im Diff nicht enthalten; Plan-Zeilen ohne Diff (Slice-Datei, Open-README, Reports, Beobachtungs-Register, Tap-Repo) sind spätere oder frühere Commits und extern |
| Umfang gegen DoD „Quelle“ | geprüft, ohne Befund: der Diff besteht aus der Matrix-Zeile und einem Kommentar, sonst nichts |
| Commit-Nachricht (Traceability) | geprüft, ohne Befund: Slice-Kennung genannt, Attribution-Zeile vorhanden; die Anforderungsbezüge [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung) und [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) stehen im Slice, den die Nachricht kennzeichnet |
| Referenzmodell, ADR-Immutabilität, Lifecycle | geprüft, ohne Befund: keine ADR, Spec oder Architektur-Sicht berührt; kein Verweis nach unten aus Rang 1 bis 3; der Slice liegt nur in `in-progress/`, die Moves sind reine `git mv`-Commits (`d87425da`, `4b767878`); die Datei im Hauptrepo ist keine Workflow-Datei dieses Repos (liegt unter `packaging/`), löst hier keinen Lauf aus |
| Öffentliche Verträge, Nutzer-Doku (`AGENTS.md` §5) | geprüft, ohne Befund: kein Runner-Label in README, `docs/user/`, `docs/maintainer/`; `CHANGELOG.md` bleibt (CI des Tap ist nicht das ausgelieferte Produkt, Plan §1) |
| Docker-only, Suppression, Gates | geprüft, ohne Befund: kein `//nolint`, keine Änderung an `Makefile`, `.golangci.yml`, `.d-check.yml`, `.github/workflows/` |
| Sprachvertrag, Exit-Codes, Dual-Classifier, Dateisicherheit, Architekturgrenzen | geprüft, ohne Befund: kein Go-Code im Diff; die CI-Datei des Tap ist kein CLI-Output oder generiertes Artefakt von `u-boot` im Sinne von [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache) |

**Nicht geprüft:**

- Ob GitHub `macos-15` als arm64-Runner bedient und ob der Lauf mit dem neuen Label grün ist — Fremdzustand, Lauf im Tap nach dem Push; DoD „Smoke“, Sache des Verifier-Nachlaufs (F-4).
- Byte-Gleichheit Quelle gegen Tap-Datei — der Push in den Tap steht nach Review und Verification (DoD „Tap“, Verifier-Nachlauf).
- Die Quellen des Plans zur Stilllegung von `macos-14` (Ankündigung, README von `actions/runner-images`, Lauf `42191356009`) — kein Netzzugriff in diesem Lauf vorgesehen; der Plan gibt sie als gelesen an.
- Die Zeile des Slice in `docs/plan/planning/open/README.md` und die DoD-Häkchen — nicht im Diff dieses Commits, Verifier-Aufgabe.
- `make gates` — nicht Gegenstand dieses Reviews (siehe Sensors).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zeilenverweis im Plan nach dem Diff veraltet · Plan-Aussage über Code ungenau · Abgelaufener Zustandsmarker im Plan · Annahme über Fremdzustand bis zur Messung · Kopplung zweier Orte ohne Sensor

## Verdikt

**Merge-blockierend:** nein. Kein HIGH und kein MEDIUM. Der Diff ist die geplante Ein-Zeilen-Änderung mit Kommentar, berührt keinen Schritt, kein Audit-Gate und keine weitere Datei. F-1 bis F-3 sind Berichtigungen des Plans (Zeilenverweise, ASCII-Aussage, Zustandsmarker in §7), die der Abgleich nach der letzten Fix-Runde bei der Closure ohnehin verlangt (`AGENTS.md` §3.10); F-4 und F-5 sind Übergaben an den Verifier und den Projektinhaber.

**Übergabe:** F-1 bis F-3 an den Implementer (Plan-Texte vor der Closure berichtigen und die Plan-Tabelle gegen `git diff --stat` des Endstands abgleichen). F-4 an den Verifier (Nachlauf: Lauf `smoke` im Tap, Job-Metadaten, „Set up job“, Formel-Version). F-5 an den Hauptagenten (Frage an den Projektinhaber, Plan §6). Die Finding-Klassen gehen in die Closure-Notiz (§7) des Slice.

## Sensor-Lauf

`make docs-check` nach Anlage dieses Reports: `1202 Datei(en) geprüft, 0 Befund(e)` (der erste Lauf mit dem Report fand eine unverlinkte ADR-Kennung in einer Negativbefund-Zeile; berichtigt).
