# Verification-Report: slice-v1-gates-tool-pruefung-zuerst — 2026-11-05

**Rolle:** Verifier (frischer Kontext; Baseline-Regelwerk `modul-11`)
**Gegenstand:** [`slice-v1-gates-tool-pruefung-zuerst`](../plan/planning/done/slice-v1-gates-tool-pruefung-zuerst.md) — Commit-Range `d28b3548..4a4cbd2a`. Code-Commits `6c4edaa9` (Reihenfolge), `36138efc` (Sensor, Wächter-Test), Review-Fix-Runde `3ef9d47b`; Abgleichstand **nach der letzten Fix-Runde = `4a4cbd2a`** (= Repo-Tip, Arbeitsbaum sauber). Die Move-Commits `7fc612f9`, `7f7b9ba9` sind Lifecycle-Bewegungen. Der Anlege-Commit `9182d4bd` und `4a4cbd2a` tragen zusätzlich den anderen Slice [`slice-v1-devcontainer-homebrew-pakete`](../plan/planning/next/slice-v1-devcontainer-homebrew-pakete.md); er ist nicht Gegenstand.

**Eingangs-Kontext:**

- Slice-Plan §1 (Ziel/Nicht-Ziele), §2 (DoD), §3 (Plan-Tabelle), §6 (zwei Risiken), §7 (Closure leer)
- [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md)
- Review-Report [`2026-11-05-slice-v1-gates-tool-pruefung-zuerst.md`](2026-11-05-slice-v1-gates-tool-pruefung-zuerst.md) (0 HIGH / 1 MEDIUM / 4 LOW / 3 INFO; als Eingang, nicht als Evidence)
- [`harness/verification.md`](../../harness/verification.md), [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.6, §3.7, §3.9, §3.10), [`harness/roles.md`](../../harness/roles.md)

---

## Verification Evidence

Scope:
- Slice: [`slice-v1-gates-tool-pruefung-zuerst`](../plan/planning/done/slice-v1-gates-tool-pruefung-zuerst.md)
- IDs: [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md)
- Artefakte: `tools/mk/tool-pruefung.mk`, `tools/harness/gate-reihenfolge.sh`, `tools/harness/gate-reihenfolge-test.sh`; Doku `harness/README.md`, `harness/conventions.md`, [`MR-012`](../../harness/conventions/MR-012-ai-harness-init-werkzeug-des-repos.md); Planung `roadmap.md`, `open/README.md`

DoD-Abgleich (Plan §2):

| # | DoD-Punkt | Zustand | Evidence |
|---|---|---|---|
| 1 | **Reihenfolge** (Plan Z. 35, `[x]`) | **erbracht** | `tools/mk/tool-pruefung.mk:18` `GATE_CHECKS := tool-pruefung $(GATE_CHECKS)`. Eigener Lauf: `make -n gates` beginnt mit `bash tools/harness/tool-pruefung.sh`, dann die zwei Sensor-Aufrufe, dann `baseline-verify.sh`. `make -pn gates` liefert genau eine Zeile `record-gates: tool-pruefung baseline-verify docs-check lint test coverage-gate ablage-pruefung`. Pin: `gate-reihenfolge.sh:14-23` (erstes Wort der Prerequisite-Liste), Wächter-Test `gate-reihenfolge-test.sh:40-45`. Der Haken stand vor Review und Verifikation; die Verifikation bestätigt ihn. |
| 2 | `make gates` grün (Plan Z. 39, Closure-Pflicht, offen) | **erbracht, Haken steht aus** | Eigener Lauf, Exit 0 (siehe Sensors). Der Haken gehört zur Closure. |
| 3 | Review, Report unter `docs/reviews/` (Z. 40) | **erbracht, Haken steht aus** | Report vorhanden; die Behebung ist unten geprüft. Der Haken gehört zur Closure. |
| 4 | Doku-Zeiger nachgesehen (Z. 41) | **gelesen, bestätigt; Haken steht aus** | `README.md:281`, `README.de.md:287`, `harness/README.md:72` führen die Gates als Menge mit `+` und nennen keine Reihenfolge. Zusatz: `docs/maintainer/quality.md:83` listet in der `make gates`-Zeile `tool-pruefung` nicht (Vorbefund, behauptet keine Reihenfolge; nicht im Zeigerkreis des Plans). |
| 5 | Closure-Notiz mit Steering-Loop-Lerneintrag, Ausgang je Risiko (Z. 42, §7 Z. 79-81) | **offen** | §7 trägt nur den Platzhalter. Ausgänge: siehe Risiken. |

Prüfung der Review-Behebung (`3ef9d47b`) gegen den Ist-Stand, nicht gegen den Report:

| Finding | Behebung trifft? | Beleg |
|---|---|---|
| F-1 (MEDIUM, Zustandsmarker in `open/README.md`) | ja | `open/README.md:23` nennt „liegt in `in-progress/`, umgesetzt, wartet auf Verifikation“; kein „Backlog“, kein „im Arbeitsbaum“ mehr. Der Marker „wartet auf Verifikation“ hat einen Ablaufpunkt (dieser Report) und ist bei der Closure zu entfernen oder auf den Ist-Stand zu setzen (§3.10). |
| F-2 (Plan-Tabelle unvollständig) | teilweise | Die drei vom Report genannten Zeilen stehen jetzt in §3 (Z. 53-55). Neu fehlt der Review-Report selbst, siehe Abweichung A-1. |
| F-3 (Kommentar mit früherem Zustand) | ja | `gate-reihenfolge-test.sh:4-7` nennt die Mutation im Indikativ („die Form, die die Zusage … ausschliesst“); kein „früher/vor der Zusage“. |
| F-4 (Kommentar-Zusage ohne Prüfung) | ja | `gate-reihenfolge.sh:2-5` sagt „nur diese Stelle wird geprueft“; der Satz über „record-gates haengt an allen Checks“ ist entfernt. Die DoD (Z. 35) fordert die Zusage über `record-gates` nicht mehr. |
| F-5 (MR-012, Sensors-Zeile) | ja | `harness/README.md:73` (Zeile `make tool-pruefung`) nennt `gate-reihenfolge.sh` und `gate-reihenfolge-test.sh`; `harness/conventions.md` (Zeile 012) und [`MR-012`](../../harness/conventions/MR-012-ai-harness-init-werkzeug-des-repos.md) führen die zwei Skripte als repo-geführte Ausnahme. `tool-pruefung.sh` fehlt weiterhin in der Ausnahmeliste (Vorbefund, siehe Hinweis H-2). |
| F-6 (Modus) | ja | `git ls-files -s`: beide Skripte `100755`. |
| F-7, F-8 (INFO) | unverändert | Vorbefunde, nicht Teil des Diffs; F-8: der DoD-Abgleich ist hiermit der Verifier-Teil. |

Plan-vs-Code (§3 gegen `git diff --stat d28b3548..4a4cbd2a`, Stand nach der letzten Fix-Runde):

| Diff-Datei | Plan-Zeile (§3) | Abgleich |
|---|---|---|
| `tools/mk/tool-pruefung.mk` (+6/−4) | Z. 48 | deckt; Kommentar (Z. 14-17) nennt Zusage und `-j`-Grenze im Indikativ |
| `tools/harness/gate-reihenfolge.sh` (neu, 24) | Z. 49 | deckt |
| `tools/harness/gate-reihenfolge-test.sh` (neu, 48) | Z. 50 | deckt |
| `docs/plan/planning/in-progress/roadmap.md` (±4) | Z. 51 | deckt |
| `docs/plan/planning/open/README.md` (+2) | Z. 53 | deckt |
| [`next/slice-v1-devcontainer-homebrew-pakete.md`](../plan/planning/next/slice-v1-devcontainer-homebrew-pakete.md) (neu) | Z. 54 | deckt (anderer Slice) |
| [`done/slice-v1-gates-tool-pruefung-zuerst.md`](../plan/planning/done/slice-v1-gates-tool-pruefung-zuerst.md) (neu) | Z. 55 | deckt |
| `harness/README.md`, `harness/conventions.md`, `MR-012` (je ±1) | Z. 56 | deckt |
| [`docs/reviews/2026-11-05-slice-v1-gates-tool-pruefung-zuerst.md`](2026-11-05-slice-v1-gates-tool-pruefung-zuerst.md) (neu, 67) | **keine Zeile** | **Abweichung A-1** |
| `repo.mk`, `Makefile`, `README.md`, `README.de.md` | Z. 52, 57 „unverändert“ | `git diff --stat` ohne Eintrag: deckt |

Gelesene Fundstellen des Plans, nachgeprüft: `repo.mk:254` (`GATE_CHECKS += lint test coverage-gate`) und `repo.mk:256` (`include tools/mk/*.mk`) stimmen; `Makefile:4` hält `.NOTPARALLEL` bewusst offen; die Suche nach `-j` in `scripts/bau.sh` und `.github/workflows/` ergibt 0 Treffer (Nicht-Ziel Z. 29 belegt). Der Commit `6c4edaa9` aus Z. 48 existiert.

Sensors:

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `bash tools/harness/gate-reihenfolge.sh` | pass | Exit 0, „tool-pruefung is the first check of record-gates — OK“ |
| `bash tools/harness/gate-reihenfolge-test.sh` | pass | 2 Fälle, 0 Fehlschläge |
| Rote Probe, selbst nachvollzogen | pass | Kopie der Make-Dateien unter dem Scratchpad (kein Eingriff ins Repo): `+=`-Mutation von `tool-pruefung.mk:18` → Exit 1, „first: 'baseline-verify'“; zweite Mutation (`lint` vorangestellt) → Exit 1, „first: 'lint'“. Der Sensor greift auf Position, nicht auf Text; ein nicht greifendes `sed` ließe im Test den Fall rot werden (Exit-0-Prüfung, `gate-reihenfolge-test.sh:31`). |
| `make -pn gates \| grep -m1 '^record-gates:'` | pass | `record-gates: tool-pruefung baseline-verify docs-check lint test coverage-gate ablage-pruefung`; unter `-j4` dieselbe Liste (die Liste, nicht die Laufreihenfolge) |
| `make tool-pruefung` | pass | Gate-Kern OK, Sensor OK, Wächter-Test OK |
| `make docs-check` | pass | d-check: 1130 Dateien, 0 Befunde; `pinsync` 3 Pins |
| `make gates` | **pass (eigener Lauf)** | Im Devcontainer-Container (`/.dockerenv`, Docker 29.9.3, `go` und `golangci-lint` vorhanden), Exit 0 auf `4a4cbd2a`. Log-Beginn: `tool-pruefung` OK, dann `gate-reihenfolge` OK und Wächter-Test 2/0 als erste Ausgabe, danach die übrigen Gates. Stempel `.harness/state/gates-passed.head` = `4a4cbd2a…` = HEAD, `gates-passed.diffsha` = `working-tree-hash.sh` (`df2bd213…`); der Stempel liegt in der ignorierten Ablage und ist kein Repo-Beleg. |

Traceability:

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) | Das Lastenheft nennt `tool-pruefung` als Gate der Pin-Inventur und sagt keine Reihenfolge zu (Plan Z. 13 gelesen: bestätigt). Der Slice verschärft die Wirkung ohne Spec-Änderung; Belege: `gate-reihenfolge.sh`, Wächter-Test, `make gates` grün. |
| [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md) | Die Werkzeuge des laufenden Kontexts tragen: unverändert; die Prüfung läuft früher, nicht anders. `make gates` im Devcontainer-Container (§3.1). |

Konformität und Gate-Lockerung (§3.6): kein Check entfällt (`tool-pruefung` bekommt zwei Aufrufe dazu); keine Schwelle, kein Lint-Profil, kein `.golangci.yml` berührt; kein `nolint` in den neuen Skripten. Keine Lockerung. Kommentare (§3.7, §3.9): `tool-pruefung.mk:14-17`, `gate-reihenfolge.sh:2-9`, `gate-reihenfolge-test.sh:2-9` beschreiben Zusage, Kopplung und Grenze im Indikativ; kein früherer Zustand.

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: [`CO-001`](../plan/carveouts/CO-001-crlf-logs.md) (Diff berührt `docs/plan/carveouts/` nicht)

Offene Risiken (Plan §6) und was die Closure braucht:

| Risiko | Zustand | Für die Closure |
|---|---|---|
| 1 — Sensor hängt am Text von `make -n` (Z. 76) | **Prämisse überholt, Ausgang steht noch auf „offen; der Implementer wählt …“** | Gewählt ist die Position in der Regel-Datenbank (`gate-reihenfolge.sh:14`, `make -pn`), nicht der Rezepttext. Der Eintrag ist zu schließen mit dem Ausgang „Position statt Text; Restgrenze: erste Zeile `^record-gates:` der Datenbank (heute genau eine), Format der `make -p`-Ausgabe“. |
| 2 — späterer `-j`-Lauf (Z. 77) | Ausgang benannt („weiter offen“), Grenze im Kommentar `tool-pruefung.mk:16` belegt | Closure-Notiz hält fest, dass die Rückführung auf `open` aus §4 greift, sobald `scripts/bau.sh` `-j` einführt; der Eintrag ist ein akzeptiertes, benanntes Restrisiko (die Prerequisite-Liste bleibt unter `-j` gepinnt, die Laufreihenfolge nicht). |

Nicht ausgeführt:
- `make test-docker` — kein Docker-Pfad berührt (nur Make-Fragment, Shell-Skripte, Doku).
- `make verify-depguard` — keine Importregel berührt.
- `make ci`, `make fullbuild` — kein Security-/Release-Pfad.
- Unter `make -j` tatsächlich gefahrener `make gates` — bewusst nicht zugesagt (Plan Z. 29).

Commit / Artefakt:
- `4a4cbd2a` (Tip); Code-Commits `6c4edaa9`, `36138efc`, Fix-Runde `3ef9d47b`

## Abweichungen und Hinweise

| ID | Gewicht | Befund | Fundstelle |
|---|---|---|---|
| A-1 | LOW | Die Plan-Tabelle führt den Review-Report [`docs/reviews/2026-11-05-slice-v1-gates-tool-pruefung-zuerst.md`](2026-11-05-slice-v1-gates-tool-pruefung-zuerst.md) nicht, der im Diff steht (`AGENTS.md` §3.10: jede berührte Datei, Abgleich nach der letzten Fix-Runde). Mit diesem Verification-Report kommt eine zweite Datei hinzu. Zeile nachziehen und den Abgleich bei der Closure wiederholen. | [`slice-v1-gates-tool-pruefung-zuerst.md`](../plan/planning/done/slice-v1-gates-tool-pruefung-zuerst.md) Z. 46-57` |
| A-2 | LOW | Risiko 1 trägt einen überholten Ausgang („offen; der Implementer wählt …“), obwohl die Wahl getroffen ist (Zustandsmarker mit Ablaufpunkt, §3.10). Bei der Closure schließen (siehe Tabelle oben). | [Slice-Plan Z. 76](../plan/planning/done/slice-v1-gates-tool-pruefung-zuerst.md) |
| A-3 | INFO | Der Marker „wartet auf Verifikation“ in `open/README.md` läuft mit diesem Report ab; bei der Closure auf den Ist-Stand setzen oder die Zeile gemäß Lifecycle behandeln. | `docs/plan/planning/open/README.md:23` |
| H-1 | INFO | DoD-Wortlaut „`make -n gates` zeigt `tool-pruefung` als ersten Check“ ist per Hand bestätigt; der Sensor liest bewusst die Datenbank (`make -pn`), nicht die `-n`-Ausgabe. Kein Mangel, die Wahl ist durch Risiko 1 gedeckt. | [Slice-Plan Z. 35](../plan/planning/done/slice-v1-gates-tool-pruefung-zuerst.md) |
| H-2 | INFO (Vorbefund) | `tools/harness/tool-pruefung.sh` ist `100644` und fehlt in der Ausnahmeliste von [`MR-012`](../../harness/conventions/MR-012-ai-harness-init-werkzeug-des-repos.md); `docs/maintainer/quality.md:83` listet `tool-pruefung` nicht. Nicht Gegenstand dieses Slice; als Beobachtung zu führen, falls die Closure sie aufnimmt. | `tools/harness/tool-pruefung.sh`, `docs/maintainer/quality.md:83` |

## Urteil

**Verifiziert mit Auflagen.** Die einzige inhaltliche DoD-Zeile (Reihenfolge samt Pin) ist erfüllt und durch eigenen Sensor-Lauf, rote Probe und grünes `make gates` im Devcontainer belegt; die Review-Behebung trifft F-1 bis F-6 (F-2 bis auf A-1). Gate nicht gelockert, keine Spec- oder ADR-Abweichung, Carveouts unverändert. Vor dem Verschieben nach `done/` fehlen: A-1 und A-2 beheben, A-3 beachten, Closure-Pflichten abhaken (`make gates`, Review, Doku-Zeiger, Closure-Notiz mit Lerneintrag und Ausgängen). Kein Blocker für die Closure-Arbeit.
