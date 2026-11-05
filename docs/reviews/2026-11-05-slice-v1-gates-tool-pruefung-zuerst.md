# Review-Report: slice-v1-gates-tool-pruefung-zuerst — 2026-11-05

**Review-Art:** Code-Review gegen Plan + Konventionen (Modul 10).

**Gegenstand:** `git diff d28b3548..36138efc` (5 Commits; 7 Dateien, 284 Einfügungen, 4 Löschungen): `tools/mk/tool-pruefung.mk`, `tools/harness/gate-reihenfolge.sh`, `tools/harness/gate-reihenfolge-test.sh`, Slice [`slice-v1-gates-tool-pruefung-zuerst`](../plan/planning/done/slice-v1-gates-tool-pruefung-zuerst.md), Roadmap, `open/README.md`; [`slice-v1-devcontainer-homebrew-pakete`](../plan/planning/open/slice-v1-devcontainer-homebrew-pakete.md) nur auf Formfehler.

**Skill:** `.harness/skills/reviewer.md` @ `36138efc`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-11-05

**Eingangs-Kontext:**

- Slice-Plan [`slice-v1-gates-tool-pruefung-zuerst`](../plan/planning/done/slice-v1-gates-tool-pruefung-zuerst.md) (§1 Ziel/Nicht-Ziele, §2 DoD, §3 Plan-Tabelle, §6 Risiken)
- [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md)
- `AGENTS.md` §3.6, §3.7, §3.9, §3.10; `MR-003`, `MR-012` (`harness/conventions.md`)
- Roadmap (`docs/plan/planning/in-progress/roadmap.md`), Index `open/README.md`

**Sensors:**

- `make gates`: laut Handoff grün auf `36138efc` (Exit 0); nicht neu gefahren (Review-Rolle).
- Eigene Proben (nur lesend, keine Repo-Änderung): `bash tools/harness/gate-reihenfolge.sh` -> OK; `bash tools/harness/gate-reihenfolge-test.sh` -> 2 Fälle, 0 Fehlschläge (der `+=`-Fall endet wirklich rot: Exit 1 und Muster "is not the first check" werden verlangt, der `sed` greift, sonst wäre der Fall wegen Exit 0 rot); `make -pn gates` liefert genau eine Zeile `record-gates: tool-pruefung baseline-verify docs-check lint test coverage-gate ablage-pruefung`; Läufe mit `MAKEFLAGS=-j4`, `make -pn gates -j4`, `env -i` und `LC_ALL=de_DE.UTF-8` (Locale hier nicht installiert, Ausgabe der Regel-Datenbank enthält keine übersetzten Teile) bleiben grün.
- `make test-docker`, `make verify-depguard`: nicht gefahren — kein Docker-Pfad, keine Importregel berührt.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Zeile des Slice in `open/README.md` sagt "Backlog" und "die Änderung in `tools/mk/tool-pruefung.mk` liegt im Arbeitsbaum"; der Slice liegt in `in-progress/`, und die Änderung ist committet (`6c4edaa9`, `36138efc`). Der Nachzug des Verweises (`7f7b9ba9`) hat nur den Pfad berührt, nicht die Aussage. Ein Leser des Index hält die Änderung für ungepusht/uncommittet. | `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt wird entfernt oder auf den Ist-Stand gesetzt); §3.7 (Zustandsfeld) | `docs/plan/planning/open/README.md` · "die Änderung in `tools/mk/tool-pruefung.mk` liegt im Arbeitsbaum" (Z. 23) | Review-only | Zustandsmarker nach Ablauf stehen geblieben |
| F-2 | LOW | Die Plan-Tabelle (§3) führt nicht alle Dateien des Diffs: `docs/plan/planning/open/README.md` (zwei Zeilen, zwei Commits) und [`docs/plan/planning/open/slice-v1-devcontainer-homebrew-pakete.md`](../plan/planning/open/slice-v1-devcontainer-homebrew-pakete.md) fehlen, ebenso der Slice-Plan selbst; der Commit `36138efc` behauptet "Plan-Tabelle nachgezogen". Gegen `git diff --stat d28b3548..36138efc` stehen 3 Zeilen ohne Plan-Eintrag. | `AGENTS.md` §3.10 (Tabelle listet jede berührte Datei; Abgleich gegen `git diff --stat` nach der letzten Runde) | [`docs/plan/planning/done/slice-v1-gates-tool-pruefung-zuerst.md`](../plan/planning/done/slice-v1-gates-tool-pruefung-zuerst.md) · "## 3. Plan (vor Code)" (Z. 44-53) | Review-only | Plan-Tabelle ohne berührte Dateien |
| F-3 | LOW | Der Kopf von `gate-reihenfolge-test.sh` beschreibt die Mutation als "per += ans Ende der Liste gehaengt, wie vor der Zusage in tools/mk/tool-pruefung.mk" — er erzählt den früheren Zustand des Repos; die geltende Zusage ("Mutation: `+=` statt `:=`, tool-pruefung steht dann am Ende") genügt dem Änderer. | `AGENTS.md` §3.7, §3.9 (a) ("früher/vor der Zusage") | `tools/harness/gate-reihenfolge-test.sh` · "wie vor der Zusage in tools/mk/tool-pruefung.mk" (Z. 5) | Review-only | Kommentar beschreibt früheren Zustand |
| F-4 | LOW | Der Kopf von `gate-reihenfolge.sh` sagt, das Skript "pinnt die Zusage" einschließlich "der Nachweis record-gates selbst haengt an allen Checks"; geprüft wird nur, dass das erste Prerequisite `tool-pruefung` ist. Dass `record-gates` an allen Checks hängt (DoD: "`record-gates` als letzten") steht nirgends in einer Prüfung; die DoD-Zeile fordert, dass ein Test "das" pinnt. | `AGENTS.md` §3.7 (Kommentar beschreibt, was da ist), §3.10 (Bedingung ohne pinnenden Test) | `tools/harness/gate-reihenfolge.sh` · "der Nachweis record-gates selbst haengt an allen Checks" (Z. 3-4) | Review-only (Verifier: DoD-Zeile 1 Teil "record-gates als letzten") | Kommentar-Zusage ohne pinnende Prüfung |
| F-5 | LOW | `tools/harness/` ist in `MR-012` dem Werkzeug `ai-harness-init` zugeschrieben, mit drei benannten repo-geführten Ausnahmen (`fetch-baseline-cache.sh`, `ablage-pruefung.sh`, `ablage-pruefung-test.sh`). Die zwei neuen repo-eigenen Skripte stehen weder dort noch in `MR-012` Punkt 3; die Zeile von `make tool-pruefung` in `harness/README.md` §Sensors nennt den neuen Sensor und seinen Wächter-Test nicht, obwohl das Rezept sie jetzt ausführt. (`tool-pruefung.sh` fehlt in der Ausnahmeliste ebenfalls — Vorbefund.) | `MR-012` (Geltungsbereich, Punkt 3); `AGENTS.md` §3.8 (Gate-Index einmal in `harness/README.md` §Sensors) | `harness/conventions/MR-012-ai-harness-init-werkzeug-des-repos.md` · "außer `fetch-baseline-cache.sh`, `ablage-pruefung.sh` und `ablage-pruefung-test.sh`" | Review-only | Repo-eigene Skripte in Werkzeug-Verzeichnis ohne Ausnahme-Eintrag |
| F-6 | INFO | Die neuen Dateien `gate-reihenfolge.sh` und `gate-reihenfolge-test.sh` haben Modus `100644`; die Geschwister `ablage-pruefung*.sh` sind `100755` (`tool-pruefung.sh` ebenfalls `100644`). Alle Aufrufer nutzen `bash <pfad>`, nichts bricht. | Maintainability | `tools/harness/gate-reihenfolge.sh` · Dateimodus | Review-only | Skript ohne Ausführ-Bit |
| F-7 | INFO | Die Roadmap-Zeile unter "Offene Wellen" steht unmittelbar unter der vorbestehenden Zeile die Zeile mit leerem Linktext zum Slice [`slice-v1-bau-nur-devcontainer`](../plan/planning/done/slice-v1-bau-nur-devcontainer.md), die „liegt in in-progress/“ sagt: leerer Linktext, und der Slice liegt in `done/`. Vorbefund, nicht Teil des Diffs; die neue Zeile selbst folgt dem Muster der Nachbarn ("liegt in `in-progress/`", ohne Welle) und der Rückkehr zum Ruhe-Marker bei der Closure (Plan §3). | `MR-003`; `AGENTS.md` §3.7 | `docs/plan/planning/in-progress/roadmap.md` · "Zeile mit leerem Linktext zum Slice [`slice-v1-bau-nur-devcontainer`](../plan/planning/done/slice-v1-bau-nur-devcontainer.md)" | Review-only | Vorbestehender Zustandsfehler neben dem Diff |
| F-8 | INFO | Zuständigkeit Verifier: DoD-Zeile 1 ist im Plan bereits mit `[x]` abgehakt (vor Review); `docs/maintainer/quality.md` (Z. 83) führt in der `make gates`-Zeile `tool-pruefung` nicht auf (Vorbefund, plan-konform nicht berührt, behauptet keine Reihenfolge). | `harness/verification.md` | [`docs/plan/planning/done/slice-v1-gates-tool-pruefung-zuerst.md`](../plan/planning/done/slice-v1-gates-tool-pruefung-zuerst.md) · "- [x] **Reihenfolge:**" | Review-only | DoD-Abgleich Verifier-Sache |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/mk/tool-pruefung.mk`: `GATE_CHECKS := tool-pruefung $(GATE_CHECKS)` | geprüft, ohne Befund: wirksam (Prerequisite-Liste gelesen); `include tools/mk/*.mk` steht in `repo.mk:256` nach den Anhängestellen, `:=` erfasst `baseline-verify docs-check lint test coverage-gate`; spätere `+=` (`ablage-pruefung`) bleiben dahinter. Kommentar nennt Zusage und `-j`-Grenze im Indikativ (§3.7 erfüllt) |
| Gate-Lockerung (§3.6) | geprüft, ohne Befund: kein Check entfällt, `tool-pruefung` bekommt zwei Aufrufe dazu, Reihenfolge schärfer; keine Schwelle gesenkt, kein `nolint` |
| Rote Probe des Wächter-Tests | geprüft, ohne Befund: Grundzustand grün, `+=`-Mutation rot mit Exit 1 und Meldungsmuster; ein nicht greifendes `sed` ließe den Fall rot werden |
| Robustheit des Sensors (`make -pn gates`, `-j`, `MAKEFLAGS`, Locale, `env -i`, `pipefail` mit `grep -m1` und `|| true`) | geprüft, ohne Befund; Grenze: das Skript liest die Zeile `^record-gates:`, die erste Fundstelle der Datenbank — heute gibt es genau eine |
| Plan-Aussagen über Code (§3.10) für `repo.mk:254/:256`, `README.md:281`, `README.de.md:287`, `harness/README.md:72` | geprüft, ohne Befund: Fundstellen stimmen; die Gate-Zeilen nennen keine Reihenfolge |
| Roadmap: Format und `MR-003` | geprüft, ohne Befund (Datum "Letzte Änderung" nachgezogen, Zeile im Muster der Nachbarn; siehe F-7 für den Vorbefund) |
| [`slice-v1-devcontainer-homebrew-pakete.md`](../plan/planning/open/slice-v1-devcontainer-homebrew-pakete.md) (nur Form) | geprüft, ohne Befund: Abschnitte 1-8 vorhanden, Lifecycle-Absatz, Bezug/Spec-Stellen mit Links, Welle "ohne Welle", Modus-Block; Umsetzung nicht Gegenstand |
| `make test-docker`, `make verify-depguard` | nicht geprüft: kein Docker-Pfad, keine Importregel berührt |
| `make gates` | nicht neu gefahren (Review-Kontext); Evidence des Implementers übernommen als Stand, nicht als Urteil |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 4 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Zustandsmarker nach Ablauf stehen geblieben · Plan-Tabelle ohne berührte Dateien · Kommentar beschreibt früheren Zustand · Kommentar-Zusage ohne pinnende Prüfung · Repo-eigene Skripte in Werkzeug-Verzeichnis ohne Ausnahme-Eintrag

## Verdikt

**Merge-blockierend:** ja, bis F-1 behoben ist (MEDIUM; Textzeile in `open/README.md`). Die LOW-Findings blockieren nicht; die Umsetzung (Reihenfolge, Sensor, rote Probe) ist ohne Befund.

**Übergabe:** Findings gehen an den Implementer; Finding-Klassen in die Slice-Closure §7. Der Report ersetzt keine Verifikation — der DoD-Abgleich bleibt beim Verifier.
