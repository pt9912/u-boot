# Verification-Report: slice-v1-egress-init-ohne-luecke — 2026-10-19

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice [`slice-v1-egress-init-ohne-luecke`](../plan/planning/in-progress/slice-v1-egress-init-ohne-luecke.md) (`in-progress/`), Commit-Range `12c34da..HEAD` (HEAD `05dc167`; sieben Commits `7deb96d` bis `05dc167`). Entscheidungen: [`ADR-0027`](../plan/adr/0027-egress-init-wiederholungslauf-atomar.md) (`Accepted`), [`ADR-0025`](../plan/adr/0025-egress-liste-im-laufenden-container.md), [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md); [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion).

**Eingang:** Implementer-Zusage (DoD, Plan-Tabelle, Vermerke und Fix-Runden in §7 des Slice), Review [`2026-10-19-slice-v1-egress-init-ohne-luecke`](2026-10-19-slice-v1-egress-init-ohne-luecke.md) (0 HIGH, 2 MEDIUM, 2 LOW, 5 INFO), Fix-Commits `be0b998`, `1ba232d`, `05dc167`.

**Gesamturteil: GO mit Auflagen.** Der Code trägt die Zusage von ADR-0027 Entscheidung 1, 3 und 4; alle vier Review-Findings sind im Code, im Test und in der Doku behoben (selbst gelesen). `make gates` ist grün, der gezielte Docker-Test ist grün. Offen vor der Closure sind zwei Zustandsmarker in Planungsdokumenten, die eine überholte Aussage tragen (V-1, V-2), und die Closure-Handgriffe selbst.

## 1. Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass | eigener Lauf auf `05dc167`, sauberer Arbeitsbaum; `d-check: 533 Datei(en) geprüft, 0 Befund(e)`, `baseline-verify: v6.18.1 OK`, `coverage-gate: OK — coverage 91.40% meets threshold 90%`, `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`; Stempel `.harness/state/gates-passed.head` = HEAD |
| `TestE2E_SandboxDevcontainer_EgressRerun` einzeln (`-tags docker -count=1 -v`, im Image des Ziels `test-docker`, kein Host-Go) | pass, 13,97 s | Unterfälle `repeated_runs` (3,31 s) und `failed_transaction` (0,92 s) PASS, kein SKIP |
| `make test-docker` (voll) | **nicht von mir gefahren** | Die Zusage des Implementers (643 s, voll grün) ist mir nur als Behauptung bekannt; ein Beleg liegt nicht in einer Datei. Der Diff nach dem letzten Lauf (laut Implementer) ist `05dc167` (Doctor-Meldung, Unit-Test, kein Docker-Pfad). Die Lücke ist benannt, nicht geschlossen. |
| `make verify-depguard` | nicht gelaufen | keine Importregel berührt (Diff: Template, Tests, zwei Zeichenfolgen in `devcontainer_fallback.go`, Doku) |
| `git diff --stat 12c34da..HEAD -- docs/plan/adr spec docs/plan/planning/in-progress/carveouts.md` | leer | kein ADR, keine Spec, kein Carveout berührt |

## 2. Review-Findings gegen den Code (nicht dem Report geglaubt)

| Finding | Behoben? | Beleg (selbst gelesen) |
|---|---|---|
| F-1 (MEDIUM) fünfter Abbruchzweig im Handbuch nicht nachgezogen | ja | `docs/user/benutzerhandbuch.md`: Tabelle *Zustand / warn / fail*, `failed`-Zeile der Status-Tabelle, Zustands-Absatz („einem von drei Zuständen“, eigener Spiegelpunkt mit Erstlauf und Wiederholungslauf), Fehlersuche-Satz nennen „die nftables-Regeln lassen sich nicht anwenden“; `CHANGELOG.md` nennt den Zweig und den `ip6`-Fall; `devcontainer_fallback.go` `egressProbe` (`failed`, `hint`) und `TestDoctor_SandboxFallback_EgressFailedNamesEveryCause` pinnen die Zeichenfolgen (`05dc167`); `harness/replay.md` zählt fünf Zweige. |
| F-2 (MEDIUM) `ip6`-Hälfte ungepinnt | ja | Golden Case `TestGenerateDevcontainer_SandboxEgress_RerunIsOneTransaction` fordert beide `ip6`-Zeilen genau einmal; Docker-Probe `egressRerunProbeScript` zählt `ip6drop=1 ip6lo=1` in der Regex jedes der elf geprüften Läufe. Meine Lesung der Regex und der Läufe: der Test ist im eigenen Lauf grün. Die Mutation (beide Zeilen gestrichen) habe ich nicht neu gefahren; der Implementer nennt sie in §7. |
| F-3 (LOW) Bedingung der Baseline-Antwort im Handbuch | ja | Handbuch: „Ein Treffer färbt den Test nur, wenn das Gateway vorher ohne Beschränkung auf die Sonde geantwortet hat …“; stimmt mit dem `switch` im Test überein (`BASELINE answers=0` ergibt nur `t.Logf`). |
| F-4 (LOW) elf/zwölf Läufe | ja | Handbuch: „zwölfmal … Jeder der ersten elf Läufe …, nach dem zwölften Lauf ist die Liste leer“; Probe: Lauf 1, Läufe 2 bis 11 (`report`), zwölfter Lauf ohne `report`, `SET before=1 after=0`. |
| Reihenfolge-Nachtrag `be0b998`: Meldung nennt alten Resolver | ja | Template: `${old_alive:+ (the previous resolver was still running after 5 seconds)}`; Golden Case pinnt die Zeichenfolge in `…UpstreamDNSOnlyForTheResolverUser`. Der Docker-Test erzeugt den Zweig nicht (im Slice §7 ausgewiesen). |

## 3. DoD-Abgleich

- [x] **Entscheidung:** ADR-0027 `Accepted` vor dem Code (Geschichte-Zeile, Freigabe a bis d), Commit `7deb96d` vor `58bc45a`; ADR-0025/-0012/-0026 unverändert (Diff leer); keine Messskripte, kein Eintrag unter `docs/messungen/` im Diff.
- [x] **Template, Golden Case, Docker-Test:** Template: Konfiguration, Transaktion (`add table`, `delete table`, `add table` je Tabelle in einer Datei, ein `nft -f "$rules"`), dann `kill`, begrenztes Warten (50 mal 0,1 s), `dnsmasq`-Start, UID-Prüfung, `resolv.conf` zuletzt, `write_status ok`: so in ADR-0027 Entscheidung 1 und 4 beschrieben. Golden Case pinnt Struktur, Reihenfolge, Abwesenheit von `nft add`/`nft delete` außerhalb der Transaktion, fünf `failed`-Zweige. Docker-Test (`//go:build docker`) führt das **ausgelieferte** Script aus: `generateEgressProject` erzeugt den Baum mit dem Generator, das Image wird aus `.devcontainer/` gebaut, der Aufruf `sudo sh /usr/local/bin/u-boot-egress-init` ist derselbe wie im erzeugten `postStartCommand` (`devcontainer_sandbox.go:124`). Der Fehlerzweig läuft über eine per `sed` verfälschte Kopie des erzeugten Scripts (eine Regelzeile ungültig, mit `MUTATION not applied`-Wächter); getestet wird der echte Zweig `if ! nft -f "$rules"`, die Eingabe ist verfälscht, das Script nicht neu geschrieben.
- [x] **Handbuch und CHANGELOG:** Aussagen über die neue Fassung beschränkt auf Struktur, Test-Inhalt und ausgewiesenes Nicht-Gemessenes; die Zahlen von B-1 sind „frühere Fassung“; „Fenster von null“ und „ununterbrochene Auflösung“ ausdrücklich nicht zugesagt (ADR-0027 Entscheidung 3, `BEO-DOK/zusage-weiter-als-messung`). Keine `LH-`/`ADR-`/`SPEC-`/Slice-Kennung und kein Planungsverweis in den Diff-Zeilen des Handbuchs (gesucht und gelesen); die Kennung im CHANGELOG ist dort zulässig.
- [ ] **Sammelpunkt (`make gates` / `make test-docker` / Review / Closure-Notiz / Register / §6 / Paarungen):** `make gates` grün (eigener Lauf); `make test-docker` voll: siehe Abschnitt 1 (nicht selbst gefahren); Review-Report liegt vor; **Closure-Notiz mit Steering-Loop-Lerneintrag fehlt** (§7 trägt „Wird bei der Closure gefüllt“ und den Vermerk des Implementers); Register: §7 nennt „keine neue Beobachtung“ (Aussage, nicht im Register geprüft, siehe Abschnitt 7); alle vier Risiken in §6 tragen einen Ausgang; die Haken stehen noch nicht, wie vorgesehen.

## 4. Konformität mit den Entscheidungen

| Entscheidung | Befund |
|---|---|
| ADR-0027 Entscheidung 1 (eine Transaktion, Resolver geordnet, `resolv.conf` zuletzt; Exit-Codes, Statusdatei, `skuid`-Regel unverändert) | konform; die `skuid`-Zeilen stehen inhaltlich gleich als `echo`-Zeilen, der Golden Case `…UpstreamDNSOnlyForTheResolverUser` pinnt sie in neuer Form. |
| Entscheidung 2 (kein Neu-Laden vom Host; keine Lastenheft-Änderung) | konform; `spec/`, Lastenheft, `Makefile`, Workflows unverändert. |
| Entscheidung 3 (Struktur zugesagt, Lücken ausgewiesen) | im Handbuch und CHANGELOG konform; nicht konform in zwei Planungsdokumenten (V-1, V-2). |
| Entscheidung 4 (Transaktion schlägt fehl: Exit 11, `failed`, auch `warn`) | konform; Docker-Test belegt Erstlauf (`tables=0 resolv=0`) und Wiederholungslauf (`chain-unchanged=1 same-resolver=1 resolv=1`). Die Atomarität selbst stützt sich auf diesen Test und `nft`; für den `ip6`-Fehlerfall gilt „gelesen, nicht gelaufen“, so ausgewiesen. |
| ADR-0025 (Entscheidungen 1 bis 3), ADR-0026 (`user=`, UID-Prüfung, fünf `failed`-Zweige statt vier) | konform; kein neuer Exit-Code, kein neuer Port, kein Host-Aufruf. |
| Messvariante gegen Auslieferung | konform für Struktur und Fehlerzweig (siehe Abschnitt 3). Nicht abgedeckt, und so ausgewiesen: Die Zeiten der neuen Fassung sind nirgends gemessen; der Docker-Test misst eine Untergrenze auf einem nativen `amd64`-Host und trägt keine Zeitzusage. |

## 5. Plan-Tabelle gegen `git diff --stat 12c34da..HEAD`

Zehn Dateien im Diff: Template, `generate_sandbox_test.go`, `devcontainer_fallback.go`, `devcontainer_fallback_test.go`, `harness/replay.md`, `sandbox_devcontainer_docker_test.go`, `benutzerhandbuch.md`, `CHANGELOG.md`, die Slice-Datei, der Review-Report: jede Zeile hat einen Eintrag in §3 (die Zeile der Doctor-Meldung ist im Plan nachgetragen). Die Zeilen „entfällt“ und „unverändert“ haben keinen Diff. **Abweichung:** §3 führt `open/README.md` und `roadmap.md` „bei Start und Closure“ als Änderung; beide stehen im Diff dieser Range nicht, und beide tragen eine überholte Aussage (V-1, V-2).

## 6. Befunde

| ID | Schwere | Befund | Fundstelle |
|---|---|---|---|
| V-1 | MEDIUM | Der Zustandsmarker der Roadmap sagt: „Messung: Wiederholungslauf von `egress-init.sh` ohne Fenster ohne Restriktion; Umsetzung erst nach Messbefund und Entscheidung“. Die Messung entfällt (Entscheidung des Projektinhabers), die Umsetzung ist geschehen, und „ohne Fenster ohne Restriktion“ ist eine Zusage, die ADR-0027 Entscheidung 3 ausdrücklich nicht gibt. Verstoß gegen `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt) und gegen `BEO-DOK/zusage-weiter-als-messung`. | `docs/plan/planning/in-progress/roadmap.md`, Zeile „In Arbeit: … (Messung: Wiederholungslauf …“ |
| V-2 | MEDIUM | Die Zeile des Slice in `open/README.md` sagt dasselbe: „Wiederholungslauf … ohne Zeitraum ohne Restriktion und ohne Ausfall der Namensauflösung …; Messung, Template, Golden Case, Docker-Test … Vorbehalt: Folge-ADR … möglich, Entscheidung des Projektinhabers vor dem Start“. Messung entfällt, Folge-ADR ist `Accepted`, Zusage ist weiter als der Beleg. Die Zeile steht außerdem für einen Slice, der in `in-progress/` liegt (Link zeigt dorthin); ob die Zeile dort bleiben soll, ist bei der Closure zu entscheiden. | `docs/plan/planning/open/README.md`, Zeile `slice-v1-egress-init-ohne-luecke` |
| V-3 | LOW | §1 „Gefundener Stand“ nennt Handbuch-Zeilen 1355 bis 1366 und Test-Zeilen 180 bis 185 als Fundstelle; beides verschoben (Handbuch-Abschnitt ist umgeschrieben). Der Abschnitt ist als Stand vor der Umsetzung gekennzeichnet; die Zeilenbezüge sind dennoch nicht mehr auffindbar. | Slice §1 |
| V-4 | LOW | Das Handbuch schreibt „Zuvor blieb in diesem Fall die IPv4-Beschränkung bestehen“ (Rückblick auf die frühere Fassung in der Nutzer-Doku, `AGENTS.md` §5 Regel 6: Stand, nicht Chronik). Der Satz ist inhaltlich belegt (gelesen im Script) und als solcher markiert; er gehört nach Ablauf der Version in den CHANGELOG, nicht ins Handbuch. | `docs/user/benutzerhandbuch.md`, Spiegelpunkt „Die nftables-Regeln lassen sich nicht anwenden“ |
| V-5 | INFO | Der Docker-Test färbt das Entfernen des begrenzten Wartens nach `kill` nicht (Mutation M-B aus dem Review: Docker-Test grün, Golden Case rot); der Zweig „alter Resolver lebt nach 5 s“ hat keinen Docker-Beleg. Im Slice §7 und im Review ausgewiesen; die Zusage (Struktur) deckt das ab. | Slice §7 „Lücke“ |

## 7. Auflagen vor Closure

1. V-1: Zeile „In Arbeit“ in `roadmap.md` auf den Ist-Stand setzen (Messung entfällt, ADR-0027, Struktur zugesagt, kein Fenster von null zugesagt) oder bei der Closure entfernen.
2. V-2: Zeile in `open/README.md` auf denselben Stand bringen oder bei der Closure entfernen (Zusage nicht weiter als ADR-0027 Entscheidung 3).
3. §7 des Slice: Closure-Notiz mit Steering-Loop-Lerneintrag schreiben; die Aussage „Beobachtung: keine neue angefallen“ gegen das Register belegen (`../observations/`, insbesondere `BEO-DOK/zusage-weiter-als-messung`: zwei Planungsdokumente mit zu weiter Zusage könnten ein erneuter Fall dieser Beobachtung sein, Schwellenwert siehe dortige `state.md`); den Slice-Satz „Der Verifier-Bericht trägt Abschnitt 5“ (Release-Zuordnung) belassen.
4. DoD-Haken im Sammelpunkt erst nach 1 bis 3 und nach einem vollen `make test-docker` (Beleg als Ausgabe, Dauer und Exit) setzen; Closure-Reihenfolge nach `AGENTS.md` §3.3 Fall 2 (erst Inhalt, dann reiner `git mv`).
5. V-3 und V-4 sind LOW und optional.

## 8. Nicht geprüft

- Voller `make test-docker` (Behauptung des Implementers, 643 s grün, nicht nachgefahren); nur `…EgressRerun` einzeln.
- Mutationen aus Review und §7 (M-A, M-B, M-C, Doctor-Mutation) nicht neu gefahren; ihr Ergebnis gilt als Implementer-/Reviewer-Aussage, gestützt nur durch den Befund, dass die Pins im Test stehen.
- Verhalten auf einem Host, auf dem `nft` die Familie `ip6` nicht anlegen kann; Verhalten, wenn der alte `dnsmasq` nach `kill` länger als 5 s lebt.
- `amd64` (nativ ausgeführt, aber Testlauf, keine Messung), Docker Desktop, CI, Last, `nestedRuntime`: nach ADR-0027 Entscheidung 3 ungemessen und ausgewiesen.
- Beobachtungs-Register (`observations/`) und die drei Paarungen der Closure: nicht gelesen; Closure-Aufgabe.
- `make verify-depguard`, `make ci`/`image-scan`/`govulncheck`: Diff berührt weder Importregeln noch Abhängigkeiten.

## 9. Carveouts

Neu: none. Gelöst: none. Unverändert: Inventar `carveouts.md` (Diff leer).

## 10. Commit / Artefakt

Verifizierter Stand: `05dc167`. Dieser Report: `docs/reviews/2026-10-19-verification-slice-v1-egress-init-ohne-luecke.md`.
