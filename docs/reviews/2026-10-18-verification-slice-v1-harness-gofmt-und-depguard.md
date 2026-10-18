# Verification-Report: slice-v1-harness-gofmt-und-depguard — 2026-10-18 <!-- d-check:ignore (Slice-Kennung im Titel, Zitat) -->

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD, ADR und Spec). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice [`slice-v1-harness-gofmt-und-depguard`](../plan/planning/in-progress/slice-v1-harness-gofmt-und-depguard.md) (`in-progress/`), Commit-Range `38f17a5..HEAD` (9 Commits, HEAD `de76298`), [`LH-QA-004`](../../spec/lastenheft.md#lh-qa-004--linting-solid-nahes-lint-profil), [`ADR-0003`](../plan/adr/0003-solid-nahes-lint-profil.md), [`SPEC-034`](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln).

**Eingang:** Review `docs/reviews/2026-10-18-slice-v1-harness-gofmt-und-depguard.md` (0 HIGH/MEDIUM, F-1 LOW, F-2..F-4 INFO) und die Zusage des Implementers, F-1, F-3, F-4 seien behoben, F-2 sei als Hinweis geführt. Geprüft wurden die Belege am Code und Text; die Sensoren wurden selbst gefahren (Host: Colima, Docker, kein Host-go). <!-- d-check:ignore (Kennungen als Zitat) -->

**Gesamturteil:** Alle sachlichen DoD-Zeilen sind erfüllt. Offen sind nur die Closure-Zeilen, die dieser Rolle vorgelagert sind (Review-Häkchen, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen) — Auflagen unten. Keine Diskrepanz zwischen Plan und Code.

## 1. Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass (Exit 0, 35,7 s, warmer Cache) | `lint` „0 issues.“, `coverage-gate` 91,10 % (Schwelle 90 %), `baseline-verify` v6.18.1 OK (54 Dateien), `docs-check` im Lauf enthalten |
| `make doc-trace` | pass | „83 Anforderung(en), 0 Waise(n).“; `LH-QA-004` ok, Slice in der Spalte | <!-- d-check:ignore (Kennung als Zitat) -->
| `make verify-depguard` | pass, 9 von 9 PASS, Exit 0 (zweimal gefahren; Dauer 57 s) | Log `domain-isoliert`: genau ein Befund (`1 issues`), `is not allowed from list 'domain-isoliert': domain must not depend on port … (depguard)`, kein `import cycle`, kein `revive`; danach `git status` leer, `internal/hexagon/port` enthält nur `driven`, `driving` |
| `make doc-immutable RANGE="38f17a5..HEAD"` | pass | 9 Commits aufgelöst, 379 Dateien, 0 Befunde |
| `make test-docker` | pass, Exit 0, 15 min 8 s (`internal/e2e` 896,6 s, `application` 21,0 s) | alle Pakete `ok`. Grenze: ohne `-v`, die Ausgabe nennt keine Einzeltests; der Beleg ist „Paket ok“, nicht je Test sichtbar |
| Rot-Probe Formatter | pass (rot wie erwartet) | absichtlich falsch eingerückte, sonst saubere Datei `internal/hexagon/domain/verifier_probe.go`: `make lint` meldet `verifier_probe.go:5:1: File is not properly formatted (gofmt)`, `1 issues`; Datei gelöscht, `git status` leer, `make lint` danach Exit 0; nicht committet |
| Gegenprobe `verify-depguard` | pass (FAIL wie erwartet) | Klon aus `git archive HEAD` unter `~/uboot-verifier-clone` (eigenes `git init`), Deny-Eintrag `…/internal/hexagon/port` der Regel `domain-isoliert` entfernt: `FAIL domain-isoliert: make lint succeeded but should have failed`, übrige 8 PASS, Exit 1; Klon und Logs gelöscht |

Die Gegenprobe zeigt gegenüber dem Review einen Unterschied, der aus dem Fix F-3 folgt: Weil der `revive`-Begleitbefund entfällt, schlägt bei entfernter Regel schon die Exit-Code-Prüfung an (`make lint` grün); im Review (vor F-3) war es die Beschreibungsprüfung. Die Prüfung ist jetzt trennscharf.

## 2. Plan-vs-Code-Diff

- **Format-Commit `d94da25`:** `git diff -w d94da25~1 d94da25` ganz gelesen: 14 Dateien, 6 Einfügungen, 13 Entfernungen; ohne Leerzeilen sind es ein `//`-Trenner in `erroremission.go`, das umbrochene `t.Cleanup`-Literal in `internal/e2e/sandbox_devcontainer_docker_test.go` (gleiche Anweisung, gleiche Argumente) und die durch `gofmt` hinter die `//`-Zeile versetzte Direktive an `dependenciesFor` in `internal/hexagon/application/addservice.go`. Verhalten unverändert. Gesamt-Stat 29 Dateien, 170/177 wie im Plan. Die 15 weiteren Dateien des Commits ändern nur Whitespace.
- **F-1 (Nachzug `23d841f`):** `addservice.go:65` trägt Begründung und `//nolint:unparam` wieder in einer Zeile; `make lint` grün (Wirkung unverändert).
- **Plan-Tabelle §3 gegen `git diff --stat 38f17a5..HEAD`:** 38 Dateien = 29 formatierte (Anhang A, Namen verglichen) + 9: `.git-blame-ignore-revs`, `.golangci.yml`, `docs/maintainer/quality.md`, `docs/plan/planning/in-progress/roadmap.md`, der Slice, der Review-Report, `harness/README.md`, `harness/sensors/verify-depguard.md`, `scripts/verify-depguard.sh`. Jede Datei hat eine Zeile, jede Zeile einen Diff. „Unverändert“ (`repo.mk`, `Dockerfile`): kein Diff in der Range; die Aussage zum `lint`-Target ist durch die grüne Rot-Probe gedeckt (der Formatter greift ohne weitere Option). Plan und Diff stimmen überein; ein Wortlaut-Detail: die Plan-Zeile `quality.md` nennt die Nutzerangabe nicht, §3 „Nachzug“ nennt sie (F-4), `quality.md` enthält sie.
- **`.golangci.yml`:** `git diff 38f17a5..HEAD` zeigt nur Kopfkommentar und `formatters.enable: [gofmt]`; keine Regel, kein Linter, keine Schwelle berührt, nichts gelockert (`AGENTS.md` §3.6: Verschärfung).
- **`.git-blame-ignore-revs`:** enthält den vollen Hash `d94da2561469a7bdbec2d0f8b05fdee3315d37af` (Format-Commit), Kommentarform korrekt.
- **`scripts/verify-depguard.sh`:** Wegwerf-Paket `internal/hexagon/port/verifydepguardprobe` (fünftes Feld in `CASES`) mit Paketkommentar, ohne Importe; `trap cleanup EXIT INT TERM` räumt Verletzungsdatei und Probe-Verzeichnis; vorhandenes Probe-Verzeichnis führt zu Exit 2 vor dem Setzen von `PROBE_DIR` (kein `rm -rf` auf Fremdes); Exit-Codes 0/1/2 wie im Kopfkommentar und in `harness/sensors/verify-depguard.md`. Lauf zeigt keine liegengebliebene Datei.
- **Doku:** `quality.md` §1 nennt Formatter, Meldung, Docker-Aufruf, Nutzerangabe und `simplify` (F-4); „29 Linter“ bleibt wahr; `harness/README.md` Zeile `make lint` nennt `gofmt`, kein Target hinzugekommen; `verify-depguard.md` stimmt mit dem Skript überein. `.golangci.yml`-Kopfkommentar trägt den Formatter.
- **Bestands-`nolint` (F-2):** `grep -rn nolint internal cmd` liefert 9 Treffer, exakt die in §3 „Hinweis“ benannten: 6 wirksame (`addservice.go:65` `unparam`, `upservice_portprobe_docker_test.go:73` `gosec`, `template_list_test.go:21` `containedctx`, je `dupl` in `up_/down_/logs_acceptance_test.go`) und 3 wirkungslose `// nolint:dupl` (`up.go:265`, `down.go:216`, `logs.go:341`). Fundstellen stimmen; Hinweis ist nicht Gegenstand des Slice.

## 3. DoD-Abgleich (Zeile für Zeile)

| DoD-Zeile | Stand | Beleg |
|---|---|---|
| Teil A: Bestand `gofmt`-konform, Formatter in `.golangci.yml`, Rot-Probe, Format-Commit getrennt | erfüllt | `make lint` grün (Bestand formatiert); `formatters.enable: [gofmt]`; Rot-Probe oben; `d94da25` allein, Gate-Änderung `2265909` getrennt (`AGENTS.md` §3.3) |
| Teil B: `domain-isoliert` ohne Zyklus, 9 von 9, Exit 0, Vertrag stimmt | erfüllt | Lauf und Log oben, Gegenprobe oben |
| Doku-Update (`quality.md` §1, `.golangci.yml`-Kopfkommentar) | erfüllt | Diff oben; §1.2 bewusst unverändert (Formatter ist kein Linter, Begründung steht im Plan) |
| `make gates`, `make verify-depguard` grün, `make doc-trace` ohne neue Waise | erfüllt | Sensor-Tabelle |
| Review durchgeführt, Report liegt vor | erfüllt (Häkchen offen) | `docs/reviews/2026-10-18-slice-v1-harness-gofmt-und-depguard.md`; im Slice-§2 ungehakt: beim Closure zu setzen | <!-- d-check:ignore (Kennung als Zitat) -->
| Closure-Notiz mit Steering-Loop-Lerneintrag | offen | §7 „Wird bei der Closure gefüllt.“; Closure-Aufgabe |
| Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ | offen | Closure-Aufgabe, siehe Auflagen |
| Risiken aus §6 mit Ausgang | offen | fünf Risiken ohne Ausgang, Closure-Aufgabe (Vorschläge unten) |
| Drei Paarungen (Anker · Folge-Slice · Register) | offen | Closure-Aufgabe |

Die Häkchen in §2 des Slice sind sämtlich noch leer; der Verifier hakt nicht.

## 4. Traceability

| ID | Beleg |
|---|---|
| [`LH-QA-004`](../../spec/lastenheft.md#lh-qa-004--linting-solid-nahes-lint-profil) | Formatter-Gate (`make lint`, Rot-Probe), `make doc-trace` ok (Slice in der Spalte), Commits tragen die ID |
| [`ADR-0003`](../plan/adr/0003-solid-nahes-lint-profil.md) | `.golangci.yml` ohne gelockerte Regel, `make verify-depguard` 9 von 9 |
| [`SPEC-034`](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln) | Beweis der neun `depguard`-Regeln einschließlich `domain-isoliert` |

## 5. Carveouts

Neu: keine. Gelöst: keine. Unverändert: alle (keiner berührt; `make doc-immutable` und `docs-check` ohne Befund).

## 6. Nicht ausgeführt

- `make ci` (`govulncheck`, `image-scan`): kein Security-/Release-Pfad berührt; Abhängigkeiten, Dockerfile und Workflows unverändert.
- `make test-docker` ohne `-v`: Einzeltest-Ausgabe nicht gesehen (siehe Sensor-Tabelle).
- Unter Linux-Hosts nicht gemessen: Datei-Eigentümer nach `golangci-lint fmt` im Container (in `quality.md` als nicht gemessen benannt).

## 7. Diskrepanzen

Keine zwischen Plan, DoD und Code. Beobachtungen am Rand:

- Der e2e-Lauf hinterlässt, wie schon vor dem Slice, beendete Container (aus `docker run` ohne `--rm`; drei je Lauf). Die drei Container dieses Laufs wurden entfernt; ältere Altbestände (nicht von dieser Verifikation) blieben unberührt. Kein Befund dieses Slice.
- `make verify-depguard` hinterlässt Logs unter `$TMPDIR/verify-depguard-XXXXXX.*` (Verhalten des Skripts, Exit 0 wie Exit 1 „Logs erhalten“); die Verzeichnisse dieses Laufs wurden gelöscht. Kein Repo-Zustand.

## 8. Auflagen für die Closure (Planner)

1. **§7 Closure-Notiz** mit Steering-Loop-Lerneintrag. Finding-Klassen des Reviews eintragen (Kommentar-Satz durch mechanische Umordnung getrennt · Inline-`nolint` gegen Suppression-Verbot · Begleitbefund maskiert Exit-Code-Prüfung · Formatter-Aufruf ohne Nutzer-/`simplify`-Hinweis). Wird eine Regel verkörpert, trägt sie `liegt in` mit dem Anker `seit slice-v1-harness-gofmt-und-depguard`; sonst „keine Verkörperung“ benennen. <!-- d-check:ignore (Kennung als Zitat) -->
2. **Register:** F-2 (Bestands-`nolint`, 9 Treffer, `nolintlint` nicht aktiviert) als neue Beobachtung führen. Sub-Area ist Code-weit; Vorschlag `BEO-WZG` ist zweifelhaft, weil die Treffer in `internal/**` (CLI, APP) liegen und die Lücke ein fehlender Sensor ist: Kürzel nach Tabelle in `harness/conventions.md` wählen (`BEO-CLI`/`BEO-APP` je Ort, oder `BEO-WZG` für den Sensor), Entscheidung beim Planner. Zusätzlich: `BEO-WZG/kopplung-nicht-nachgezogen` (Zähler 1) — der Slice ist ein zweiter Fall (das Beweis-Skript folgte der Import-Lage nicht), der Planer entscheidet, ob er den Zähler auf 2 hebt; Schwelle 3 nicht erreicht. Alternativ in §7 notieren, dass nichts angefallen ist.
3. **Risiken §6** je ein Ausgang. Vorschläge nach Befund: (1) Formatierungs-Diff verdeckt Logik: entfallen (`git diff -w` gelesen, nur Kommentarzeile versetzt, F-1 behoben). (2) Formatter bricht Test oder Sensor: entfallen (`make gates`, `make test-docker` grün). (3) Gate bricht Zweig-Arbeit: entfallen (WIP-Limit 1, kein anderer Zweig). (4) Wegwerf-Import bildet Zyklus oder trifft Regel nicht: entfallen (9 von 9, Gegenprobe). (5) Skript hinterlässt Paket bei Abbruch: entfallen (`trap`; Abbruchprobe des Implementers per `SIGTERM` im Plan §3 belegt, hier nicht wiederholt).
4. **DoD-Häkchen** in §2 setzen (die Review-Zeile und alle sachlichen Zeilen sind belegt); `make doc-immutable`/`docs-check` nach der Closure-Notiz erneut.
5. **Move nach `done/`** gemäß `AGENTS.md` §3.3 Fall 2: erst Inhalt (Häkchen, Closure-Notiz), dann reiner `git mv` (`make slice-mv`). <!-- d-check:ignore (Kennung als Zitat) -->
6. **Roadmap-Marker nach dem Move:** `roadmap.md` Zeile 20 („In Arbeit: …“) auf „Nichts in Arbeit.“ zurück; Zeile 26 (Tabelle „Nach v0.11.0: Hygiene“) trägt den Link `../in-progress/slice-v1-harness-gofmt-und-depguard.md`, der auf `../done/…` zu ziehen ist; der Slice-Move zieht Verweise nach, `make slice-mv` prüft nichts (Wächter: `make docs-check`). <!-- d-check:ignore (Kennung als Zitat) -->
7. **`docs/plan/planning/open/README.md` Zeile 21:** nennt den Slice mit Link auf `../in-progress/…` in der Open-Liste; nach dem Move den Link auf `../done/…` ziehen und prüfen, ob die Zeile in dieser Liste weiter stehen soll (Lifecycle: der Slice liegt nicht mehr in `open/`).
8. **Optional (nicht blockierend):** `git config blame.ignoreRevsFile .git-blame-ignore-revs` im Klon dokumentieren (steht als Kommentar in der Datei).

## Verdikt

**Closure-fähig nach den Auflagen 1–7:** DoD sachlich erfüllt, keine Diskrepanz zwischen Plan und Code, alle Sensoren grün, Rot-Probe und Gegenprobe wie erwartet.
