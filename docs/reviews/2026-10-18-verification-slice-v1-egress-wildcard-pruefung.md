# Verification-Report: slice-v1-egress-wildcard-pruefung — 2026-10-18 <!-- d-check:ignore (Slice-Kennung im Titel, Zitat) -->

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD, ADR und Entscheidung des Projektinhabers). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice `slice-v1-egress-wildcard-pruefung` (`docs/plan/planning/in-progress/`), Commit-Range `1cec55f..HEAD` (6 Commits, HEAD `dba70e2`), [`ADR-0012`](../plan/adr/0012-devcontainer-egress-firewall.md) (`Accepted`, unverändert), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes). Entscheidung des Projektinhabers: kein Wildcard-Alias; erklärender Fehlertext und Tests. <!-- d-check:ignore (Slice-Kennung als Zitat) -->

**Eingang:** Review `docs/reviews/2026-10-18-slice-v1-egress-wildcard-pruefung.md` (F-1, F-2 MEDIUM; F-3..F-7 LOW; F-8..F-10 INFO) und die Zusage des Implementers (Commit `dba70e2`), alle Findings seien behoben oder begründet. Geprüft wurden Belege am Code und Text; die Sensoren wurden selbst gefahren (Colima, Docker, kein Host-go). <!-- d-check:ignore (Kennungen als Zitat) -->

**Gesamturteil: GO** für die Closure. Alle sachlichen DoD-Zeilen sind erfüllt, alle Review-Findings sind behoben oder begründet, die Entscheidung des Projektinhabers ist im Code umgesetzt (kein Wildcard-Parsing, anderer Text zu demselben Exit-Code `10`). Offen sind nur die Closure-Zeilen, die dieser Rolle nachgelagert sind (Auflagen unten). Keine blockierende Diskrepanz.

## 1. Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass, Exit 0 (Exit-Code von `make` selbst erfasst, nicht von einer Pipe) | `coverage-gate`: 91,20 % bei Schwelle 90 %; `docs-check`: 398 Dateien, 0 Befunde; `baseline-verify` v6.18.1 OK (54 Dateien); Baum danach sauber |
| `make doc-trace` | pass, Exit 0 | „83 Anforderung(en), 0 Waise(n).“ |
| `make test-docker` (Colima) | pass, Exit 0 | alle Pakete `ok`, e2e-Paket 762,7 s, gesamt 12 min 49 s; kein `FAIL`, `panic`, `SKIP`; danach keine Container `sbxe2e-*` und keine Images `uboot-e2e-*` auf dem Host. `TestE2E_SandboxDevcontainer_EgressSubdomains` ist Teil des Pakets (ohne `-v` nicht einzeln ausgewiesen, siehe N-1). |
| `make verify-depguard` | nicht gelaufen | der Diff ändert keine Importregeln (`domain` importiert nur Standardbibliothek; `strings` im Test ist Test-Import); zudem ist `domain-isoliert` laut Vorgänger-Report Bestand-rot |
| `make doc-immutable` | nicht gelaufen; ersatzweise gelesen | `git diff 1cec55f --stat -- docs/plan/adr spec` ist leer: ADR-0012 und Spec unverändert | <!-- d-check:ignore (Kennung als Zitat) -->

## 2. DoD-Abgleich (jede Zeile)

| # | DoD-Zeile | Stand | Evidence (am Code und Text geprüft) |
|---|---|---|---|
| 1 | Messung a–d mit Messstand in Evidence und §3, Ergebnis E1/E2/E3 | erfüllt (Haken gesetzt) | §3 Tabelle a–d, Messstand (Docker 29.5.2, Colima 0.10.3, arm64, `dnsmasq` 2.91, in den Rohdaten als `Dnsmasq version 2.91` bestätigt); Evidence-Datei trägt Auswertung, Reproduktion und Grenzen; Ergebnis E1 |
| 2 | Belegte Aussage: Handbuch im Wortlaut der Messung, Test pinnt, Fehlertext mit Test | erfüllt (Haken gesetzt) | Handbuch-Punkte *Subdomains*, *Namen hinter einem CNAME*, *Messstand*; e2e `TestE2E_SandboxDevcontainer_EgressSubdomains` (Build-Tag `docker`); Fehlertext in `domain.NewEgressHost` mit `TestNewEgressHost`, `TestNewEgressHost_WildcardLengthAndRawQuote` und CLI-Test `TestConfigSet_EgressAllowWildcard_Exit10WithHint` |
| 3 | Ergebnis für die Welle: §7 `Gegenstand:`, Welle nachgeführt | offen | §7 steht auf „Wird bei der Closure gefüllt.“; E1 zieht keinen Slice in die Welle, die Welle verweist auf den Slice im Ort `in-progress/` — Auflage A-1 |
| 4 | `make gates` grün | erfüllt | §1 |
| 5 | Review durchgeführt, Report liegt vor | erfüllt (abzuhaken) | `docs/reviews/2026-10-18-slice-v1-egress-wildcard-pruefung.md`; Haken im Slice noch leer | <!-- d-check:ignore (Kennung als Zitat) -->
| 6 | Doku-Update Handbuch; CHANGELOG bei geändertem Fehlertext | erfüllt (Haken gesetzt) | Handbuch-Diff gelesen; CHANGELOG `[Unreleased]` → Changed nennt Text, Exit-Code `10` und Hostname-Vorschlag, stimmt mit Code und Tests überein |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | offen | Auflage A-1 |
| 8 | Beobachtungs-Register fortgeschrieben | teils erfüllt | Evidence `evidence/slice-v1-egress-wildcard-pruefung.md` unter `BEO-DOK/spec-bedingung-ohne-code-entsprechung` liegt vor (Fund: Aussage stand ohne Beleg, war wahr); die Zählung im Register (`observation.md`) und die Verknüpfung in §7 stehen aus — Auflage A-3 | <!-- d-check:ignore (Kennung als Zitat) -->
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | erfüllt im Wortlaut | §6: zwei „weiter offen“ mit Begründung, eine „eingetreten“, eine „entfallen“, `offen` bei Ollama (Rückweg benannt). Abzuhaken mit der Closure |
| 10 | Die drei Paarungen (Anker, Folge-Slice, Register) | offen | Closure-Schritt (Welle-Closure prüft) — Auflage A-1/A-3 |

## 3. Entscheidung und ADR-Konformität

| Prüfpunkt | Ergebnis |
|---|---|
| Kein Wildcard-Alias, keine Wildcard-Syntax im Code | konform: `NewEgressHost` gibt bei `*` immer einen Fehler zurück; die Annahme-Zweige (`h` gültig) lassen `*` nicht durch (Muster ohne `*`); kein Parsing von `*.`-Einträgen. Nur `strings.TrimPrefix(h, "*.")` für den Vorschlagstext |
| Erklärender Fehlertext | konform: `entry "*.example.org" is a wildcard; subdomains of an allowed name are already included, enter "example.org" instead`; ohne gültigen Namen hinter dem Stern: „enter the host name without `*`“; sonstige Fehler behalten den alten Text (Test `wantNot`) |
| Exit-Code unverändert | konform: gleicher Sentinel `ErrInvalidSandboxSetting` → `ErrConfigValueInvalid`, Exit `10` über Produktionskette gepinnt (`config_egress_wildcard_test.go`: ein Paar, mehrere Paare, Wildcard nach gültigem Eintrag, Datei unverändert) |
| ADR-0012 unverändert `Accepted`, keine Spec-/Lastenheft-Änderung | konform: Diff über `docs/plan/adr` und `spec` leer; Handbuch-Zusagen stehen bis zum Folgeslice ohne Spec-Anker (Review F-10, in §6 als „weiter offen“ geführt) | <!-- d-check:ignore (Kennung als Zitat) -->
| Aussage „Subdomains mitgemeint“ stimmt mit Entscheidung 2 der ADR | konform; Messung a/a2/b bestätigen sie |

## 4. Review-Findings am Code und Text

| Finding | Stand | Prüfung |
|---|---|---|
| F-1 MEDIUM CNAME-Aussage allgemein | behoben | Handbuch begrenzt auf „in den drei gemessenen Fällen“, nennt gleichen CDN-Anbieter, einen Lauf und „es gibt keinen Test dafür“; §6 begründet den fehlenden Test (kein stabiles Ziel, kein zweiter Resolver) |
| F-2 MEDIUM Messprotokoll außerhalb des Repos | behoben | `evidence/messung-slice-v1-egress-wildcard-pruefung/` mit `run.sh`, `inner.sh`, `inner-cd.sh`, `rohdaten-1..3.txt`; Evidence-Datei erklärt das „Address already in use“-Artefakt (Zweig `d`, `pkill` + `sleep 1`) und stützt `.example.org` auf Lauf d3; „aus Antworten geschlossen, nicht getrennt belegt“ ist ehrlich benannt | <!-- d-check:ignore (Kennung als Zitat) -->
| F-3 „gesperrt“ nur Namensauflösung | behoben | Handbuch: „öffnet den übergeordneten Namen nicht für die Namensauflösung“; teilen beide Namen eine Adresse, ist sie auch per Adresse offen. Rohdaten-1 belegen identische Adressen für `example.org` und `www.example.org`; die Folge (Adresse im Set) ist Mechanismus-Aussage ohne eigenen Test (N-2, INFO) |
| F-4 Plan-Tabelle gegen Diff | behoben | §5 unten: jede der 24 Dateien hat eine Zeile; Satz im Perfekt mit Stand |
| F-5 Zustandsmarker | behoben | Suche nach „bedingt“ und „offen bis zur Messung“ findet nur den Satz in §3, der das Fehlen feststellt; §8 nennt den Fehlertext ohne „bedingt“ |
| F-6 `slice-v2-spec-…` Aussage veraltet | behoben | Aussage gegen das Handbuch neu gelesen und berichtigt (Punkte *Subdomains*, *Namen hinter einem CNAME*, *Messstand*) | <!-- d-check:ignore (Kennung als Zitat) -->
| F-7 Fehlertext nur auf Domänenebene | behoben | CLI-Test über die Produktionskette (Exit 10, Text, Datei unverändert); Grenze 253 und rohes Zitat im Test `TestNewEgressHost_WildcardLengthAndRawQuote`. Mutationen nicht selbst gefahren (Produktionscode darf der Verifier nicht ändern); aus dem Test gelesen: 253/254 und Trim-Unterschied werden getrennt geprüft |
| F-8 INFO Wächter-Grenze | begründet | §6 „weiter offen“ mit Register `BEO-E2E/test-gruen-nur-mit-lokalem-cache`; falsch-grün ausgeschlossen (Positivfälle im selben Lauf) |
| F-9 INFO Startrennen | begründet | §6 „weiter offen, bewusst nicht behoben“ mit Begründung (Generator-Output, Replay-Fälle) und Zuordnung zu `slice-v1-sandbox-egress-neuladen-entwurf`. Hinweis: Dieser Lauf (`make test-docker`, grün, 762,7 s) ist der zweite Lauf ohne Rennen; Flakefreiheit bleibt unbelegt | <!-- d-check:ignore (Kennung als Zitat) -->
| F-10 INFO Handbuch ohne Spec-Anker | begründet | §6 „weiter offen“; Folgeslice `slice-v2-spec-sandbox-egress-ollama` | <!-- d-check:ignore (Kennung als Zitat) -->

## 5. Plan-Tabelle §3 gegen `git diff 1cec55f --stat` (§3.10)

24 berührte Dateien: Alle stehen in der Tabelle. Zugeordnet: Code (`sandbox.go`, `sandbox_test.go`, `config_egress_wildcard_test.go`, e2e-Test) · Handbuch · CHANGELOG · Slice-Datei · Roadmap · Review-Report · Evidence-Datei und Messprotokoll (6 Dateien, als Gruppe in einer Zeile) · Verweis-Nachzüge (`welle-sandbox-egress-ollama.md`, `open/README.md`, fünf `open/slice-…`-Dateien, `open/slice-v2-spec-…`). Zeilen ohne Diff sind ausschließlich „unverändert“ (`egress-init.sh.tmpl`, mit gelesener Fundstelle `kill "$(cat "$pidfile")"`) und „entfallen“ (Folge-ADR). Kein Eintrag ohne Diff, keine Diff-Zeile ohne Eintrag. Eine Nebenbeobachtung: die Verifikations-Datei dieses Reports kommt mit diesem Commit hinzu und gehört in die Closure-Tabelle (A-1). <!-- d-check:ignore (Kennungen als Zitat) -->

## 6. Handbuch gegen Messung und Test (§3.10, §5 Zeile 6)

| Aussage im Handbuch | Beleg | Stand |
|---|---|---|
| `example.org` erlaubt `www.example.org` | Messung a, e2e `A www.example.org` | gemessen und gepinnt |
| `amazonaws.com` erlaubt `s3.us-east-1.amazonaws.com` (zwei Ebenen) | Messung a2, e2e `A s3.us-east-1.amazonaws.com` | gemessen und gepinnt |
| Eintrag `www.example.org` lässt `example.org` nicht auflösen | Messung b (analog `api.github.com`), e2e `B example.org` = `000` | gemessen und gepinnt (`000` unterscheidet Auflösung nicht von Verwurf, Grenze in §6 des Slice) |
| gleiche IP-Adresse: bei Verbindung per Adresse offen | Rohdaten-1 (identische Adressen), Mechanismus des `nftset` | nicht getestet, als Folge benannt (N-2) |
| `*.example.org` wird mit Hinweis abgewiesen | Unit- und CLI-Tests | gepinnt |
| CNAME: in drei Fällen erreichbar, Zielname löst nicht auf, kein Test | Messung c; Begrenzung ausgewiesen | auf Messung begrenzt |
| Messstand: Docker 29.5.2, Colima 0.10.3, arm64, `dnsmasq` 2.91 | Rohdaten | gedeckt |

Kennungen im Handbuch (Ausschnitt Egress und Diff): keine `LH-*`, `ADR-*`, `SPEC-*`, Slice- oder Planungs-Verweise (Suche über den Egress-Abschnitt ohne Treffer; `docs-check` 0 Befunde). Regel §5 Zeile 6 erfüllt.

## 7. Diskrepanzen und Annahmen

- **N-1 (INFO):** Dieser `make test-docker`-Lauf zeigt den neuen Test nicht einzeln (kein `-v`); er ist im e2e-Paket (`ok`, 762,7 s) enthalten, kein Skip. Ein Einzellauf mit `-run EgressSubdomains -v` wurde nicht gefahren, weil der Paketlauf ihn bereits enthält.
- **N-2 (INFO):** Die Aussage „Adresse ist per Adresse offen, wenn beide Namen sie teilen“ ist Mechanismus-Wissen (nftables-Set je Adresse) und nur über die identischen Adressen der Rohdaten gestützt, nicht durch einen Verbindungstest über die IP. Sie ist als Bedingung („teilen sich … eine Adresse“) formuliert und trägt keine Zusage über Namen; für die Closure in §7 als Grenze führen, oder mit einem Test pinnen, wenn der Spec-Slice sie übernimmt.
- **N-3 (LOW, Zählbasis):** Die Evidence zu `BEO-DOK/spec-bedingung-ohne-code-entsprechung` ist ein weiterer Vorgang (Beleg-Lücke, Aussage wahr). Das Register trägt „verkörpert“ (`AGENTS.md` §3.10); der Retirement-Check ist nicht berührt, die Zählung ist in `observation.md` nachzuführen oder §7 sagt „keine neue Verkörperung“ (A-3).
- **Annahme ohne Bestätigung:** keine; die Entscheidung zum Wildcard-Alias ist vom Projektinhaber getroffen und im Code umgesetzt.

## 8. Auflagen für die Closure

- **A-1 Closure-Notiz (§7)** mit `Gegenstand:` „entfallen“ (Wildcard-Syntax; Begründung aus den Messungen a–d und der Entscheidung des Projektinhabers), Messergebnis, Steering-Loop-Eintrag (Klassen: CNAME-Aussage zu allgemein aus einem CDN; Messprotokoll außerhalb des Repos; „gesperrt“ meint nur DNS; Plan-Tabelle ohne Link-Nachzüge; Zustandsmarker nach Ablauf), Grenzen (Colima/arm64, Linux, Docker Desktop, amd64, IPv6; Fall c ungetestet; Startrennen und Wächter `000` offen), und den Haken in DoD 3, 5, 7–10 erst nach dem Schreiben. Plan-Tabelle um diese Verification-Datei ergänzen, Abgleich gegen `git diff --stat` nach der letzten Änderung wiederholen.
- **A-2 Ruheort und Roadmap:** nach dem Inhalt der reine `git mv` nach `done/` (`AGENTS.md` §3.3); Marker „Nichts in Arbeit.“ zurück in `roadmap.md`; `next/README.md`/`done/README.md` nachziehen; Pfade vom Ruheort aus schreiben. Die Welle bleibt (weitere Slices offen).
- **A-3 Register:** `observation.md` nachführen (Vorgang, Zählbasis) oder in §7 „keine neue Verkörperung“ nennen; Evidence liegt vor.
- **A-4 Folgeslices:** `slice-cr-sandbox-egress-anzeige` trägt die Bedingung „gilt nur, wenn … E1 bestätigt“; mit der Closure ist die Bedingung erfüllt und beim Start dieses Slice zu streichen. <!-- d-check:ignore (Kennungen als Zitat) -->

## 9. Verification Evidence (Block nach `harness/verification.md`)

Scope:
- Slice: `slice-v1-egress-wildcard-pruefung` <!-- d-check:ignore (Slice-Kennung als Zitat) -->
- IDs: [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`ADR-0012`](../plan/adr/0012-devcontainer-egress-firewall.md)
- Artefakte: `internal/hexagon/domain/sandbox.go`, `sandbox_test.go`, `internal/adapter/driving/cli/config_egress_wildcard_test.go`, `internal/e2e/sandbox_devcontainer_docker_test.go`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, Messprotokoll unter `docs/plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/evidence/`

Sensors:

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` | pass | Exit 0, 91,20 % Coverage, `docs-check` 398 Dateien 0 Befunde |
| `make test-docker` | pass | Exit 0, e2e 762,7 s, kein Skip, keine Reste auf dem Host |
| `make doc-trace` | pass | 0 Waisen |

Traceability:

| ID | Beleg |
|---|---|
| `LH-FA-DEV-008` | `TestE2E_SandboxDevcontainer_EgressSubdomains`, `TestNewEgressHost`, `TestNewEgressHost_WildcardLengthAndRawQuote`, Handbuch-Abschnitt Egress mit Messstand, Messprotokoll | <!-- d-check:ignore (Kennung als Zitat) -->
| `LH-FA-CLI-006` | `TestConfigSet_EgressAllowWildcard_Exit10WithHint` (Exit `10`) | <!-- d-check:ignore (Kennung als Zitat) -->
| `ADR-0012` | unverändert `Accepted`; Subdomain-Aussage (Entscheidung 2) durch Messung und e2e-Test belegt | <!-- d-check:ignore (Kennung als Zitat) -->

Carveouts:
- Neu: none · Gelöst: none · Unverändert: bestehende Einträge

Nicht ausgeführt:
- `make verify-depguard`, `make doc-immutable` — keine Importregel- und keine ADR-Änderung im Diff (ADR/Spec per `git diff --stat` leer)
- `make ci` (`govulncheck`, `image-scan`) — keine Abhängigkeits- oder Image-Änderung
- Mutationen an `NewEgressHost` und am e2e-Test — Produktionscode darf der Verifier nicht ändern
- Linux ohne VM, Docker Desktop, amd64, IPv6 — nicht verfügbar; im Handbuch und in §6 als nicht gemessen geführt

Commit / Artefakt:
- HEAD `dba70e2` (Range `1cec55f..HEAD`)
