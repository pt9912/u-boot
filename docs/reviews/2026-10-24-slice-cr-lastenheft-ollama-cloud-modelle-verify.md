# Verification-Report: slice-cr-lastenheft-ollama-cloud-modelle — 2026-10-24

**Rolle:** Verifier (Modul 8) · **Frage:** Wurde der Slice richtig gebaut, gemessen an Plan, DoD, ADR und ausgeführten Sensors? (nicht: Review, nicht: Validierung) · **Modell:** claude-sonnet-5-5 · frischer Kontext.

**Gegenstand:** Slice [`slice-cr-lastenheft-ollama-cloud-modelle`](../plan/planning/done/slice-cr-lastenheft-ollama-cloud-modelle.md) (liegt in `in-progress/`, DoD-Häkchen offen außer Punkt 1), Stand HEAD `7e32651b` (11 Commits `origin/main..HEAD`), Lastenheft `0.23.2` ([`spec/lastenheft.md`](../../spec/lastenheft.md), Commits `3b135e0e`, `9c8c16eb`, `29cb728b`). Geprüft gegen Slice §1 (freigegebener Wortlaut: Fassung 1 Stellen 1, 2, 5, 6, 8; Fassung 2 Stelle 4; Fassung 3 Stelle 3; Stelle 7 Historie), §2 (DoD), §3 (Plan-Tabelle), §6, §7 (Freigaben je Stelle), die Reports [`…-ollama-cloud-modelle.md`](2026-10-24-slice-cr-lastenheft-ollama-cloud-modelle.md) (F-1 bis F-12), [`…-nachlauf.md`](2026-10-24-slice-cr-lastenheft-ollama-cloud-modelle-nachlauf.md) (N-1 bis N-7) und [`…-nachlauf-2.md`](2026-10-24-slice-cr-lastenheft-ollama-cloud-modelle-nachlauf-2.md) (L-1, I-1 bis I-4), [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) (Accepted; Z1 bis Z13; *Folgepflicht*) und die Messprotokolle (nur Gemessenes: [`slice-v1-ollama-cloud-messung`](../messungen/slice-v1-ollama-cloud-messung/README.md), [`…-mit-konto`](../messungen/slice-v1-ollama-cloud-messung-mit-konto/README.md)).

**Nicht Gegenstand (Auftrag):** Closure-Notiz und Welle-Nachzug (Planner kommt danach); das gesprochene Gespräch mit dem Projektinhaber.

---

## Verification Evidence

Scope:
- Slice: [`slice-cr-lastenheft-ollama-cloud-modelle`](../plan/planning/done/slice-cr-lastenheft-ollama-cloud-modelle.md)
- IDs: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (geändert); [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (unverändert); [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md), [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) (unverändert)
- Artefakte (Diff `origin/main..HEAD`, 10 Dateien): `spec/lastenheft.md` (+12/−4, fünf Hunks: Zeile 9, 191, 407 bis 418, 460, 1604 bis 1606), Slice-Plan (neu in `in-progress/`), `roadmap.md`, `carveouts.md` (+2), `open/README.md`, `open/slice-cr-lastenheft-messstaende.md`, `welle-ollama-cloud-modelle.md`, drei Review-Reports. Kein Go-Code, kein Docker-Pfad, keine Importregel. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

### Sensors (selbst gefahren)

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | Exit 0 (`gates=0` im Log); `baseline-verify: v6.18.1 OK — 54 Dateien`; `d-check: 859 Datei(en) geprüft, 0 Befund(e)`; Coverage-Gate und Tests im Docker-Stage grün; `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege` |
| `make doc-trace` | pass | `84 Anforderung(en), 0 Waise(n)` |
| `make doc-commits RANGE="origin/main..HEAD"` | pass | 11 Commits aufgelöst, 0 Befunde |
| `make doc-immutable RANGE="origin/main..HEAD"` | pass | 0 Befunde |
| `git diff origin/main..HEAD` über [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) und [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) | pass | 0 Byte; ebenso über `docs/plan/adr/` insgesamt |
| Wortvergleich Lastenheft gegen Slice §1 (`awk`/Bash-Parameterersetzung, `grep -c -F -x -f -`) | pass | siehe Prüfung 1 |
| `make test-docker`, `make verify-depguard`, `make ci` (`govulncheck`, `image-scan`) | not run | nur Markdown berührt, kein Docker-Pfad, keine Importregel; `verification.md` §Sensor-Auswahl: Mindest-Sensor `make docs-check`, hier durch `make gates` überdeckt |

Der Wortvergleich lief auf dem Host als reiner Textvergleich ohne Go-Toolchain (kein Nachweis im Sinn von AGENTS.md §3.1); Skript und Zwischendateien liegen im Scratchpad, nicht im Repo.

---

## Prüfung 1: Zeichen-für-Zeichen-Treue Lastenheft ↔ freigegebener Wortlaut

Methode: je Stelle die Blockzitat-Zeile aus Slice §1 (ohne `> `), Pfadpräfix `../../../../spec/lastenheft.md#` durch `#` ersetzt, `JJJJ-MM-TT` durch `2026-10-24`, bei Stelle 6 das Auslassungszeichen am Ende entfernt; dann ganze Zeile (`-x`) bzw. Teilzeichenkette gegen `spec/lastenheft.md` (HEAD) und, für die Zwischenstände, gegen `git show 3b135e0e:` und `git show 9c8c16eb:`.

| Stelle | Fassung (Freigabe) | Treffer im Lastenheft | Ergebnis |
|---|---|---|---|
| 1 Spiegelstrich | 1 | `0.23.2`: 1 (ganze Zeile 410) | identisch |
| 2 letzter Satz Absatz `ollama` | 1 | `0.23.2`: 1 (Teilzeichenkette der Zeile 412); Diff: einziger Unterschied ist der entfallene Halbsatz „ein Anschluss an einen Ollama-Dienst des Hosts,“ | identisch |
| 3 Absatz Cloud-Modelle, *Neu* (3553 Zeichen) | 3 | `0.23.2`: 1 (ganze Zeile 414); *Alt* (Fassung 2, 3351 Zeichen): 0 Treffer | identisch |
| 4 Absatz Anschluss an den Host, *Neu* (3017 Zeichen) | 2 | `0.23.2`: 1 (ganze Zeile 418); *Alt* (Fassung 1): 0 Treffer; unverändert seit `9c8c16eb` | identisch |
| 5 `LH-FA-DEV-006` Aufzählung | 1 | `0.23.2`: 1 (ganze Zeile 460); *Alt*: 0 Treffer | identisch | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| 6 `LH-FA-CLI-005` Ausnahme 1 | 1 | `0.23.2`: 1 (Teilzeichenkette der Zeile 191); „Vier Ausnahmen“ unverändert | identisch | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| 7 Historie `0.23.0` | 1 | Zwischenstand `3b135e0e`: 1; HEAD Zeile 1604 | identisch, bleibt als Chronik |
| 7 Historie `0.23.1` | 2 | Zwischenstand `9c8c16eb`: 1; HEAD Zeile 1605 | identisch, bleibt als Chronik |
| 7 Historie `0.23.2` + Versionsfeld | 3 | HEAD Zeile 1606: 1; Versionsfeld Zeile 9 `0.23.2` (Tabellenbreite erhalten) | identisch |
| 3 und 4 in den Zwischenständen | 1, 2 | `3b135e0e`: Fassung-1-Zeilen 72 und 76: je 1; `9c8c16eb`: Fassung-2-Zeilen 110 und 118: je 1 | identisch |
| 8 (`LH-FA-DEV-008`, `-009`, `-010`) | 1 | Lastenheft-Diff berührt nur Kopf, `LH-FA-CLI-005`, `LH-FA-DEV-003`, `LH-FA-DEV-006`, §16 | unverändert | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

Die fünf Hunks des Gesamtdiffs (`origin/main..HEAD`) sind genau: Zeile 9, 191, 407 bis 418, 460, 1604 bis 1606. Jede dieser Änderungen liegt in einer freigegebenen Stelle; es gibt keine Änderung ohne Freigabe-Eintrag und keinen Freigabe-Eintrag ohne Änderung.

Decken-Regel ([`AGENTS.md`](../../AGENTS.md) §3.4): Suche über die hinzugefügten Lastenheft-Zeilen auf `ADR-0`, `slice-`, `welle-`, `SPEC-0`, `ARC-0`, `roadmap`, `host.docker.internal`, Versionsnummern, Modellnamen: kein Treffer in den neuen Zeilen (der einzige Treffer liegt im bestehenden Messstand-Halbsatz des Absatzes `ollama`, Zeile 412, den Frage F3 stehen lässt und der Carveout 2 trägt). Hostname im Neutext nur `ollama.com`; Links ausschließlich Anker innerhalb der Datei.

### Reihenfolge Freigabe → Änderung (Prüfung 1, zweiter Teil)

| Fassung | Wortlaut-Entwurf im Plan (Planner) | Freigabe-Eintrag in §7 | Lastenheft-Änderung |
|---|---|---|---|
| 1 | `33fc62f3` | Commit `3b135e0e` (2026-10-24 19:19) | Commit `3b135e0e` |
| 2 | `a4c9d8ae` | Commit `9c8c16eb` (20:01) | Commit `9c8c16eb` |
| 3 | `9b3870f5` (20:10) | Commit `29cb728b` (20:16) | Commit `29cb728b` |

`git log origin/main..HEAD -- spec/lastenheft.md` nennt genau diese drei Commits; kein Lastenheft-Commit liegt vor einem Wortlaut-Entwurf. Der Entwurf stand also je Fassung vor der Änderung im Repo, und `spec/lastenheft.md` blieb zwischen Entwurf und Freigabe unverändert (die Planner-Commits `a4c9d8ae`, `9b3870f5` tragen „spec/lastenheft.md unverändert“ und berühren die Datei nicht).

**Grenze, ausgewiesen statt geraten:** Der Freigabe-Eintrag und die Lastenheft-Änderung liegen je Fassung **im selben Commit**. Die Reihenfolge „Eintrag vor Änderung“ (DoD Punkt 1: „kein Commit von `spec/lastenheft.md` vor diesem Eintrag“) ist aus `git log` nicht ablesbar; sie ist nur durch die Aussage des Implementers in §7 belegt. Ob die Vorlage an den Projektinhaber den Entwurf ungekürzt mit *Alt* und *Neu* zeigte ([`BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf`](../plan/planning/observations/BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf/observation.md)), ist aus dem Repo nicht prüfbar. Belegt ist: §7 trägt je Fassung „Freigeben wie vorgelegt“ und „Abweichung gegen den Entwurf in §1: keine“, und das Lastenheft stimmt mit §1 Zeichen für Zeichen überein. Der Freigabe-Nachweis ist damit ein Eintrag des Implementers, kein Artefakt des Projektinhabers. Das ist dieselbe Lage wie in `N-7` und `F-12` der Reviews; sie ändert sich durch den Verifier nicht.

---

## DoD-Abgleich (Slice §2, Punkt für Punkt)

| # | DoD-Punkt | Urteil | Evidence |
|---|---|---|---|
| 1 | Freigabe des Projektinhabers vor der Änderung, je Stelle in §7 mit Datum | erfüllt, mit der Grenze oben | §7 *Freigabe je Stelle* (Stellen 1 bis 8, Fassung 1), *Freigabe der Fassung 2* (3, 4, 7), *Freigabe der Fassung 3* (3, 7); Fragen F1 bis F6, G1 bis G8, H1 bis H4 alle Standard; Abweichung: keine. Die Freigabe der Fassung 1 trägt Stellen 3, 4 und 7 nicht mehr (am Review abgelehnt, §7 *Review und Fassung 2*), die Stellen 1, 2, 5, 6, 8 bleiben freigegeben (§7 Tabelle Fassung 3, letzte Zeile) |
| 2 | Wortlaut im Lastenheft wortgleich zur Freigabe; Zielbild ohne Messwert, Hostnamen, Prozesszustand, Kennung; Ausschluss „Anschluss an einen Ollama-Dienst des Hosts“ im Absatz `ollama` entfällt, Zeile `0.15.0` bleibt; jede Bedingung belegt (Beleg-Tabelle in §7) | **teilweise** | Wortgleich und frei von Kennungen/Hostnamen: erfüllt (Prüfung 1). Ausschluss entfällt (Zeile 412), `0.15.0` bleibt (Zeile 1596; Treffer „Anschluss an einen Ollama“ nur 410, 1596, 1604). Offen: (a) **die Beleg-Tabelle in §7 fehlt** (§7 schließt mit „Beleg-Tabelle … werden von Reviewer, Verifier und bei der Closure gefüllt“); ihre Funktion tragen die Tabellen *Herkunft jedes Satzes, Fassung 2* und *Fassung 3* in §1, die ich Satz für Satz gegen ADR und Protokoll gelesen habe (Abschnitt unten). (b) Der Aufbau-Rest L-1 (siehe Prüfung 2) ist nicht entschieden. (c) Z13 hat bis zum Code-Slice keinen Beleg (ausgewiesen, siehe unten) |
| 3 | Versionsfeld und Historie: `0.23.0` / `0.23.1` / `0.23.2`, Zeile nur mit freigegebenen Stellen, ohne Messangabe | erfüllt | Zeile 9 `0.23.2`; Zeilen 1604 bis 1606; keine Engine-Version, kein Host, kein Modell in den drei Zeilen |
| 4 | `make gates`, `make docs-check` grün; `make doc-trace` Beleg für `LH-FA-DEV-003`, `-006`, `LH-FA-CLI-005` | erfüllt | Sensors-Tabelle: `gates` Exit 0, `d-check` 859 Dateien 0 Befunde, `doc-trace` 84 Anforderungen 0 Waisen | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| 5 | Review durchgeführt; Nachlauf-Review über jeden Lastenheft-Diff nach einer Fix-Runde, als eigene Datei, vor dem Verifier | erfüllt | drei Reports (Review über `3b135e0e`, Nachlauf über `9c8c16eb`, Nachlauf 2 über `29cb728b`), je ein anderer Kontext als der Implementer; Befund-Status siehe Prüfung „Review-Befunde“ |
| 6 | Verifier: Satz-für-Satz-Vergleich gegen §7 und Zusage-Menge der ADR; Plan-Tabelle gegen `git diff --stat` nach der letzten Fix-Runde | erfüllt durch diesen Report | Prüfung 1, Satz-gegen-Z-Abschnitt, Prüfung 4 |
| 7 | Doku-Update: keines | erfüllt | `docs/user/`, `README*`, `CHANGELOG.md`, `spec/spezifikation.md`, `spec/architecture.md` nicht im Diff |
| 8 | Closure-Notiz, Steering-Loop-Eintrag | offen | Planner nach diesem Report (nicht Gegenstand) |
| 9 | Beobachtungs-Register fortgeschrieben | offen | Closure (nicht Gegenstand); Diff berührt `observations/` nicht |
| 10 | Jedes Risiko aus §6 trägt einen Ausgang | offen, mit Berichtigung | siehe V-2: Risiken 1, 4 und 5 tragen „weiter offen bis zum zweiten Nachlauf-Review“ bzw. „… vor dem Verifier“, ein Zustand, der durch Nachlauf 2 und diesen Report abgelaufen ist |
| 11 | Drei Paarungen (Anker · Folge-Slice · Register) | offen | Welle-Closure |

---

## Satz-gegen-Zusage-Menge (Abgleich der Herkunftstabellen mit ADR-0036 und Protokoll)

Die Zusage-Menge ist [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) Z1 bis Z13; „Nicht gemessen, nicht zugesagt“ ist der Gegenabschnitt. Gelesen wurden die Zeilen Z1 bis Z13 und der Gegenabschnitt der ADR sowie Abschnitt F und die Tabelle der Handvarianten des Protokolls ohne Konto.

| Satz im Lastenheft | Zusage-Zeile / Test | Urteil |
|---|---|---|
| Stelle 3: Cloud-Modelle aufrufbar im benannten Aufbau; Ablauf (Server selbst starten, `ollama signin`, `ollama run`/`ollama launch claude`/`claude -p`, kein `ollama pull`) | Z1, Z2, Z3 (Pfad: Handvariante `discover`-Container; die Pfad-Spalte entspricht dem Aufbau-Satz „Container aus dem Image …, vom Nutzer gestartet, `NET_ADMIN`, ohne Egress-Skript, ohne Workspace-Volume, mit Modell-Volume, ohne Terminal“) | gedeckt |
| `claude -p` mit Basis-URL, Platzhalter-Token, kein `ANTHROPIC_API_KEY` | Z1; Entscheidung 6 und Punkt k | gedeckt |
| `u-boot` führt `ollama signin` nicht aus, speichert nichts, reicht keine Variable durch | Test der Fitness Function im Code-Slice (`agentCredentialEnv`, keine Namen in erzeugten Dateien) | **Beleg folgt mit dem Code**, in der ADR-Fitness-Function als Zeile vorhanden |
| Geräteschlüssel im Modell-Volume, nicht im Workspace-Volume; im Aufbau des Ablaufs und im CLI-Container | Z6 (Protokoll ohne Konto Abschnitt F „in allen Modi“; bestätigt `rohdaten-06-vorbereiten-1-114212.txt`) | gedeckt, **bis auf die Lage CLI-Container ohne Egress-Restriktion (L-1)** |
| Lookup-Fehler und Anmeldeaufforderung nur für `ollama run`, im CLI-Container, ohne Anmeldung | Z5 (`ollama run`, Exit 1, `lookup ollama.com`; mit `egress.allow=ollama.com` per `config set` die Anmeldeaufforderung; ohne Konto) | gedeckt |
| Default-Allowlist ohne `ollama.com` | Entscheidung 3; Test: Zeile `ollama` in `egressFeatureHosts` hat genau zwei Namen | Beleg folgt mit dem Code |
| Zusage nur Linux-Host `amd64`, lokale Docker-Engine; Colima, Docker Desktop, rootless, Podman, entfernt, `arm64`, Windows nicht zugesagt | Bedingung jeder Zeile (Messlage); Form wie [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | gedeckt |
| Neuanlage mit demselben Volume: **nicht zugesagt**, in „Nicht zugesagt“ geführt | Z4 ist Handvariante (`restricted`, Eintrag von Hand, 1 Lauf); der Weg über `config set` mit Konto steht in der ADR unter „nicht gemessen“ | Vertrag enger als ADR (V-4) |
| Stelle 3, „Nicht zugesagt“ | Gegenabschnitt der ADR: Hosts des Servers zum Cloud-Dienst; `restricted` mit Konto und Eintrag über `config set`; interaktiver Ablauf; Kontingent; bestimmtes Modell; `OLLAMA_API_KEY`; `ollama signout`; abgelaufene Anmeldung; zweiter Lauf; das von `u-boot` erzeugte Sandbox-Profil ohne Egress-Restriktion | vollständig gegen die Auswahl der ADR (Stelle-3-Liste nennt jeden Punkt) |
| Stelle 4: `ollama-host` in beiden Profilen, genau ein Eintrag, kein Relay, keine Variable, kein Egress-Eintrag | Z13 (Produktversprechen) | **Beleg folgt mit dem Code**; die ADR führt Z13 ausdrücklich als „nicht belegt“ |
| Stelle 4: Loopback-Bindung: ohne Relay nicht erreichbar; andere Bindung keine Zusage | Z11; Gegenabschnitt „Host-Server mit anderer Bindung als `127.0.0.1`“ | gedeckt |
| Stelle 4: mit Relay antwortet `claude -p`, Bedingungen Basis-URL, Platzhalter, kein API-Key | Z9, Z10; Entscheidung 6, Punkt k | gedeckt |
| Stelle 4: Ausweisung beim Erzeugen, Sicherheitslockerung im Sandbox-Profil; Warnung ohne Abbruch hinter der Egress-Restriktion, keine Adresse in der Allowlist | Entscheidung 8, O2, O3; Z12 (trägt nicht); Test der Fitness Function | Teil Messung (Z12) gedeckt; Ausweisung und Warnung ohne Z-Zeile, Beleg folgt mit dem Code (Frage F6, freigegeben) |
| Stelle 4: Aufbau (Features `claude-code`, `node`; Standardprofil oder Container des Sandbox-Profils des Nutzers mit `NET_ADMIN`, ohne Egress-Skript; Eintrag von Hand; Relay des Nutzers) | Z9 (Standardprofil), Z10 (`discover`), Pfad Handvariante | gedeckt |
| Stelle 4, „Nicht zugesagt“ (Wirkung des erzeugten gegenüber dem von Hand gesetzten Eintrag, andere Container, `ollama run`/`ollama launch` hinter der Egress-Restriktion, Hosts des Dienstes, erzeugtes Sandbox-Profil ohne Egress-Restriktion, Modell, Kontingent, interaktiv, `ANTHROPIC_API_KEY`, zweiter Lauf, Versionen) | Gegenabschnitt der ADR; Entscheidung 4 („erzeugte Zeile statt von Hand eingesetzter: nicht gemessen“) | gedeckt |
| Stellen 5, 6: Ausweisung als Sicherheitslockerung; unter `--quiet` im Sandbox-Profil sichtbar | Entscheidung 8, O3; Test der Fitness Function | Beleg folgt mit dem Code |
| Stelle 7: Historie-Zeilen | Chronik, keine Zusage-Zeile nötig | — |

**Z13 und die Sätze ohne Messzeile.** Z13, die Ausweisung, die Warnung und die Sichtbarkeit unter `--quiet` stehen als Produktversprechen im Lastenheft, ihr Beleg ist der Test des Code-Slice. Das ist freigegeben (Frage F6, Standard) und in ADR und Slice §6 Risiko 4 ausgewiesen. Der Auftrag an den Code-Slice steht in der Welle-Zeile `slice-v1-devcontainer-ollama-cloud` („je Profil ein pinnender Test“, Welle §4), nicht in „Welle §5“, wie Slice §6 Risiko 4 sagt (V-3). Das Lastenheft sagt hier ein Zielbild zu, das Code und Spezifikation noch nicht tragen; das trägt der Carveout 1 (Prüfung 3). <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

---

## Prüfung 2: Nachlauf-2-Befund L-1 — Annahme als bekannte Grenze oder weitere Wortlaut-Freigabe?

**Befund (L-1, LOW):** Der Schlüssel-Satz gilt „in einem Container, den die Dev-Container-CLI aus der erzeugten Konfiguration des Sandbox-Profils startet (mit Workspace-Volume und Startbefehl; im Folgenden der CLI-Container)“. Die Definition nennt die Egress-Restriktion nicht. Z6 ist im CLI-Container nur mit Egress-Skript als `postStartCommand` gemessen (Modi `restricted`, `allow`); die ADR führt „das von `u-boot` erzeugte Sandbox-Profil ohne Egress-Restriktion (ohne `NET_ADMIN`, ohne protokollierenden Resolver, mit Workspace-Volume und `postStartCommand`)“ als nicht gemessen. Die Egress-Sätze desselben Absatzes sind durch „Hinter der Egress-Restriktion“ begrenzt, der Schlüssel-Satz nicht.

**Bestätigt am Protokoll.** Abschnitt F nennt die Datei-Anlage „in allen Modi (`discover`, `restricted`, `allow`), vor jedem Anmeldeversuch“; der Modus `discover` ist nach der Tabelle der Handvarianten ein zweiter Container per `docker run`, **ohne** Egress-Skript, **ohne** Workspace-Volume. `restricted` und `allow` sind Container der Dev-Container-CLI **mit** Egress-Skript. Ein CLI-Container ohne Egress-Skript ist in keinem Modus gefahren. Der Befund trifft die Lage zu.

**Bewertung für diese Verifikation.**

- *Gegen eine Annahme ohne Text-Änderung:* Die ADR führt genau diese Lage unter „nicht gemessen“; die Klasse „Zusage weiter als Messung“ hat im Register 24 Belege ([`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md)) und ist Slice §6 Risiko 1; [`AGENTS.md`](../../AGENTS.md) §3.10 verlangt für eine Bedingung in der Spec einen pinnenden Test oder die Streichung. Die Zusage ist im Wortlaut formal weiter als die Messung.
- *Für eine Annahme:* Der Satz sagt nur, **wo** `ollama serve` den Schlüssel ablegt (Modell-Volume statt Workspace-Volume). Die Anlage ist in drei Modi gemessen, darunter einem ohne Egress-Skript (`discover`); das Egress-Skript läuft als `postStartCommand` und legt nichts im Modell-Volume an. Eine Abhängigkeit des Ablageorts von der Egress-Restriktion ist durch die Messung nicht angezeigt (nicht gemessen heißt nicht widerlegt). Der Absatz nimmt den Ablauf mit Anmeldung im Container von `u-boot devcontainer up` ohne Egress-Restriktion ausdrücklich aus „Nicht zugesagt“ aus, sodass die Lage, an der L-1 hängt, für den Anmelde-Teil nicht zugesagt ist; offen bleibt nur der Schlüssel-Satz, der keine Anmeldung braucht. Das Risiko ist gering und der Reviewer hat es so bewertet.
- *Kosten der Alternative:* Eine Wortlaut-Änderung (etwa „im CLI-Container mit Egress-Restriktion“) ändert Lastenheft-Text und braucht eine Freigabe des Projektinhabers (Fassung 4, Version `0.23.3`), ein drittes Nachlauf-Review als eigene Datei vor dem Verifier ([`BEO-DOK/vertragsaenderung-nach-review`](../plan/planning/observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)) und erneute Läufe aller Sensoren; es wäre die vierte Fassung in einem Tag.

**Urteil des Verifiers.** Eine Annahme als bekannte Grenze ist **vertretbar und ausreichend**, **aber nur unter drei Bedingungen, die heute nicht erfüllt sind** (§7 hat keinen Eintrag zu L-1):

1. **Entscheidung des Projektinhabers**, nicht des Implementers oder Planners (Nachlauf 2 Verdikt: ein Widerspruch des Implementers stuft den Befund nicht herab; der Wortlaut ist sein Vertrag). Der Planner stellt die Frage mit beiden Wegen und dem Standard „Annahme, Text unverändert“.
2. **Eintrag in §7** mit Begründung (Z6 in allen drei Modi gemessen, darunter ohne Egress-Skript; Ablageort unabhängig vom Egress-Skript, aber für den CLI-Container ohne Egress-Restriktion nicht gemessen; Wortlaut bleibt) und **Ausgang in §6 Risiko 1** (derzeit „weiter offen bis zum zweiten Nachlauf-Review“, V-2).
3. **Folgepflicht an die Folge-Slices** (Welle-Zeilen `slice-v1-devcontainer-ollama-cloud` und `slice-v2-spec-ollama-cloud`): Handbuch und Spezifikation dehnen den Schlüssel-Satz nicht auf den Container der CLI ohne Egress-Restriktion aus, und ein Test des Code-Slice oder die Validator-Prüfung darf sich nicht auf die Egress-lose Lage stützen. Ohne diese Zeile wiederholt sich die Klasse ein Stratum tiefer. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

Ändert der Projektinhaber den Text, gilt die Alternative: Fassung 4 mit Freigabe, Nachlauf 3, anschließend erneute Verifikation der Stelle 3 (Prüfung 1 und der Z6-Zeile). **Vor `done/` ist eine weitere Wortlaut-Freigabe also nicht zwingend, wohl aber die Entscheidung über L-1.** Bis sie vorliegt, bleibt DoD Punkt 2 (b) offen.

---

## Review-Befunde (alle behoben oder begründet angenommen?)

| Befund | Urteil | Beleg |
|---|---|---|
| F-1 bis F-4 (MEDIUM, Reichweite, Aufrufbedingungen, Bedingungen am zweiten Absatz, Nicht-zugesagt-Listen) | behoben (Fassung 2), Rest F-4 in N-1 | Nachlauf 1 Tabelle „Prüfung der Vorreport-Befunde“; von mir gegen das Lastenheft `0.23.2` gelesen (Zeilen 414, 418) |
| F-5 (MEDIUM, zwei Zusage-Mengen) | **begründet angenommen**, nicht behoben | Wortlaut weist die Trennung aus („Diese Zusage gilt für den Weg zu Cloud-Modellen; sie ist von der Zusage zur Installation … getrennt“); Carveout 2 mit Plan-Anker [`slice-cr-lastenheft-messstaende`](../plan/planning/open/slice-cr-lastenheft-messstaende.md) und Trigger; Frage G6 Standard „Carveout, Welle-Closure prüft vor dem Release-Cut“. Abweichung von [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) *Folgepflicht* („der Messstand im Absatz entfällt“) ist als Planänderung ausgewiesen und freigegeben (Frage F3); der Messstand steht weiter in Zeile 412 |
| F-6 (LOW, Bindung) | behoben | Fassung 2, Zeile 418: „Lauscht der Dienst nur auf der Loopback-Adresse …, ist er … ohne einen Relay des Nutzers nicht erreichbar; für einen Dienst mit anderer Bindung gibt es keine Zusage“ |
| F-7 (LOW, Carveouts) | behoben | `carveouts.md` +2 Zeilen (Prüfung 3) |
| F-8 (LOW, Zeilenangaben im Messstände-CR) | behoben | Diff `slice-cr-lastenheft-messstaende.md`: Zeilen 1101, 433, 1509 und 1584, 412, 472, 474, 1588 bis 1596 auf den Stand `0.23.2`; am Bestand gegengelesen (`gemessen ohne Docker-Kommunikation` in Zeile 1101). Rest-Hinweis V-5 | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| F-9 (LOW, Zustandsmarker „erwartet 0.23.0“) | behoben | `grep -F "erwartet \`0.23.0\`"` über den Slice: kein Treffer |
| F-10 (INFO, Zusatzwarnung unter `--quiet`) | **begründet angenommen** | Frage G8 Standard; Stelle 6 unverändert; die Warnung gehört nach Entscheidung 8 zur Meldung der Ausweisung; der Code-Slice pinnt es |
| F-11 (INFO, Z13 ohne Beleg) | bleibt offen, ausgewiesen | ADR führt Z13 als „nicht belegt“; Beleg ist der Test des Code-Slice (Welle §4 Zeile des Code-Slice) |
| F-12 (INFO, Freigabe-Vorlage) | nicht prüfbar | siehe Prüfung 1, Grenze |
| N-1 (MEDIUM, Neuanlage-Zusage) | behoben | Satz aus der Zusage entfernt (Nachlauf 2; von mir: Zeile 414 enthält „neu angelegt“ nur im Block „Nicht zugesagt“); Folge für die ADR: V-4 |
| N-2 (LOW, „Aufruf“) | behoben | Lookup-Fehler und Anmeldeaufforderung an `ollama run` gebunden; `ollama launch claude`, `claude -p` in „Nicht zugesagt“ |
| N-3 (LOW, „heute 0.23.0“ in der Welle) | behoben | `grep -F "heute 0.23.0" docs/plan/planning`: kein Treffer; Welle §4 und Frage 3 nennen `0.23.2` als gelesenen Stand |
| N-4 (LOW, Aufbau von Z5/Z6) | behoben bis auf L-1 | CLI-Container im Wortlaut benannt; Rest siehe Prüfung 2 |
| N-5 (INFO, `--quiet`) | angenommen (G8) | wie F-10 |
| N-6 (INFO, „geprüft ist“ im Zielbild) | **begründet angenommen** | der Reviewer stuft es als Prüfzustand ohne Messwert ein; freigegeben mit Fassung 2 und 3; Hinweis V-6 |
| N-7 (INFO, Freigabe im selben Commit) | angenommen, Grenze ausgewiesen | Prüfung 1 |
| I-1 (INFO, zwei Container-Begriffe) | angenommen | freigegeben mit G3 und H3; der CLI-Container ist nur über `u-boot devcontainer up` erreichbar, ein Test der Egress-Sätze im Code-Slice hängt am Befehl von `u-boot` (Hinweis an den Code-Slice) |
| I-2 (INFO, Negation im Aufbau ohne Workspace-Volume) | angenommen | der Satz bleibt wahr |
| I-3 (INFO, Z4) | siehe V-4 | |
| I-4 (INFO) | wie N-5 bis N-7 | |
| L-1 (LOW) | **offen, Entscheidung beim Projektinhaber** | Prüfung 2 |

---

## Prüfung 3: Spannung Lastenheft ↔ Spezifikation, Handbuch, README.de.md, Messstände-CR

Das Lastenheft `0.23.2` sagt zu, was die tieferen Straten und der Code nicht tragen. Die Spannungen im Bestand, gelesen am HEAD:

| Stelle | Aussage (gelesen) | Carveout |
|---|---|---|
| [`SPEC-047`](../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen), `spec/spezifikation.md` Zeile 739 | „**Nicht zugesagt.** GPU-Beschleunigung, ein Anschluss an einen Ollama-Dienst des Hosts, Auswahl oder Verwaltung von Modellen …“ | Carveout 1 |
| `docs/user/benutzerhandbuch.md` Zeile 2023 | „**Cloud-Modelle sind nicht enthalten:** …“ (Zeile 2043 führt „ein Anschluss an einen Ollama-Dienst des Hosts“ weiter unter den Ausschlüssen, am HEAD gelesen; im Carveout benannt) | Carveout 1 |
| `README.de.md` Zeile 14 | „Cloud-Modelle von Ollama sind nicht enthalten“ | Carveout 1 |
| Katalog im Code (`featureCatalogue`) | kennt `ollama-host` nicht (Beleg: ADR *Stand im Code und im Vertrag*) | Carveout 1 |
| Zwei Zusage-Mengen im Feature `ollama` (Zeile 412 Colima/`arm64`, Zeilen 414 und 418 Linux `amd64`) | im Lastenheft selbst | Carveout 2 |

Das Inventar [`carveouts.md`](../plan/planning/in-progress/carveouts.md) trägt beide als **temporär** (Diff +2 Zeilen):

- **Carveout 1:** Plan-Anker [`welle-ollama-cloud-modelle`](../plan/planning/done/welle-ollama-cloud-modelle.md) §4 (Folge-Slices `slice-v1-devcontainer-ollama-cloud`, `slice-v2-spec-ollama-cloud`, Kennungen vorläufig) und die *Folgepflicht* der ADR; **Trigger:** beide Folge-Slices liegen in `done/`; der Zwischenzustand erreicht keinen Release. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- **Carveout 2:** Plan-Anker [`slice-cr-lastenheft-messstaende`](../plan/planning/open/slice-cr-lastenheft-messstaende.md) (Datei besteht, Frage 1); **Trigger:** der Slice liegt in `done/`; die Bindung vor dem Release-Cut entscheidet der Projektinhaber (Frage G6).

Urteil: **erfüllt**. Beide Spannungen sind ausgewiesen, haben Plan-Anker und Trigger; die Reihenfolge CR → Code → Spezifikation ist die der ADR. Zwei Anmerkungen ohne Blockade:

- Der Anker von Carveout 1 ist eine **Welle-Zeile** mit „Kennung vorläufig, noch keine Datei“. Der Trigger „beide Folge-Slices liegen in `done/`“ ist erst mechanisch prüfbar, wenn die Dateien existieren. Beim Anlegen der Folge-Slices ist der Anker auf die Dateien zu ziehen (Hinweis an den Planner; Memory-Regel „Carveouts brauchen Slice-Pläne“).
- Carveout 2 hängt an einem Slice ohne Welle, der keinen Release aufhält; die Bindung vor `v0.16.0` ist an die Welle-Closure delegiert (G6). Das ist freigegeben (Standard), aber ein Release mit zwei Mengen im Lastenheft bleibt möglich, wenn der Projektinhaber ihn so freigibt.

Die ADR-*Folgepflicht* nennt außerdem `CHANGELOG.md` und die Golden-Zeile in `harness/replay.md`; beide sind Aufgaben der Folge-Slices (Welle-Zeile `slice-v2-spec-ollama-cloud`) und kein Widerspruch zum Lastenheft; sie brauchen keinen Carveout. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

---

## Prüfung 4: Plan-Tabelle gegen `git diff --stat origin/main..HEAD`, Zustandsmarker (§3.10)

| Datei im Diff | Zeile in Plan-Tabelle §3 |
|---|---|
| `spec/lastenheft.md` | ja (Stellen 1 bis 7, Zeilen am Stand vor dem CR gelesen) |
| `docs/plan/planning/in-progress/roadmap.md` | ja |
| `docs/plan/planning/open/README.md` | ja |
| `docs/plan/planning/open/slice-cr-lastenheft-messstaende.md` | ja | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| `docs/plan/planning/welle-ollama-cloud-modelle.md` | ja |
| `docs/plan/planning/in-progress/carveouts.md` | ja |
| `docs/plan/planning/done/slice-cr-lastenheft-ollama-cloud-modelle.md` („diese Datei“) | ja | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| `docs/reviews/` (drei Reports) | ja (Zeile „neu“; der Verification-Report kommt mit diesem Commit hinzu) |
| `docs/plan/planning/observations/` | Plan-Zeile „update (Closure)“; **kein** Diff, weil die Closure noch aussteht |

Jede Datei im Diff hat eine Zeile, und jede Zeile ohne Diff ist die Closure-Zeile `observations/`. Die Zeile „unverändert“ (`spec/spezifikation.md`, Handbuch, READMEs, `CHANGELOG.md`, `harness/replay.md`, ADR, `internal/`, `cmd/`, `tools/`) wird durch den Diff bestätigt: keine dieser Dateien ist berührt. Die gelesenen Fundstellen der Zeile (Spezifikation Zeile 739, Handbuch Zeile 2023, `README.de.md` Zeile 14) stimmen am HEAD.

Die Aussagen der Plan-Tabelle zu den Zeilen des Messstände-CR (412, 472, 474, 1588 bis 1596, 1101, 433, 1509, 1584) und zu den Lastenheft-Zeilen (9, 191, 410, 412, 414, 418, 460, 1604 bis 1606) habe ich am HEAD gegengelesen; sie stimmen.

**Zustandsmarker (§3.10), Befunde:**

- **V-2 (LOW)** `slice-cr-lastenheft-ollama-cloud-modelle.md` §6: Risiko 1 („weiter offen bis zum zweiten Nachlauf-Review“), Risiko 4 („weiter offen: Nachlauf-Review über genau den Diff der Fassung 3 … vor dem Verifier“), Risiko 5/7 sind Zustandsmarker mit Ablaufpunkt; das zweite Nachlauf-Review liegt vor. Ebenso §1 Kopf der Fassung 3: „Es folgen das Nachlauf-Review über genau diesen Diff und der Verifier.“ Aufgabe: auf den Ist-Stand setzen (Planner, mit der Closure; jedes Risiko braucht dort einen Ausgang). <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- **V-3 (LOW)** §7 *Beleg-Tabelle* fehlt (DoD Punkt 2); §6 Risiko 4 verweist für die Aufträge an den Code-Slice auf „Welle-Datei §5“, die Aufträge stehen in Welle §4, Zeile `slice-v1-devcontainer-ollama-cloud` (§5 enthält Abhängigkeiten und Fragen). Die Herkunftstabellen in §1 tragen die Beleg-Funktion; entweder §7 verweist auf sie, oder die Tabelle wird aus dem Abschnitt „Satz-gegen-Zusage-Menge“ dieses Reports übernommen. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- **V-5 (INFO)** `open/slice-cr-lastenheft-messstaende.md` Zeile 29: „… der vor diesem Slice läuft und Lastenheft `0.23.0` erzeugt“; der CR endet bei `0.23.2`. Der Satz beschreibt die Vergabe beim Start; mit Zustandswort „erzeugt“ ist er für einen Leser der Closure missverständlich. Vorschlag: „… ab `0.23.0`, am Stand `0.23.2` gelesen“. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

---

## Prüfung 5: Z4-Hinweis — als Folgepflicht für den Spec-Slice ausgewiesen?

**Befund V-4 (LOW):** [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) führt Z4 („Neuer Container mit demselben Volume … Antwort ohne neue Anmeldung“) in der Tabelle **Zugesagt** (Pfad: Handvariante) und die Zusage-Menge als Z1 bis Z13 (Entscheidung 10, Folgepflicht). Das Lastenheft `0.23.2` sagt die Neuanlage **nicht mehr** zu und führt sie in „Nicht zugesagt“ (Fassung 3, Standard H1; der Weg über `config set` mit Konto steht in der ADR unter „nicht gemessen“). Die Accepted-ADR bleibt unverändert (`git diff` leer, [`AGENTS.md`](../../AGENTS.md) §3.5). Der Vertrag ist enger als die ADR, was zulässig ist („kein Satz ohne Zeile“, nicht „jede Zeile ein Satz“; Entscheidung 10).

**Wo der Hinweis heute steht:** nur in Slice §1 (Herkunftstabelle Fassung 3, Zeile „3.8 entfällt“: „in der Zusage-Menge als Handvariante geführt, ohne Zusage“; Frage H1) und als INFO I-3 im Nachlauf 2 („Hinweis für Architect und Spezifikations-Slice“). **Nicht** steht er in einer Zeile, die der Spec-Slice liest: weder in der Welle-Zeile `slice-v2-spec-ollama-cloud` (Welle §4, Zeile 81) noch im Carveout 1 noch in der Slice-Closure (die ist nicht geschrieben). `grep -n "Z4"` über Welle, `carveouts.md` und `open/README.md`: kein Treffer. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

**Risiko:** Die Spezifikation und das Handbuch lesen die ADR-*Folgepflicht* („Zusage-Menge Z1 bis Z13“) und führen die Neuanlage als zugesagt, obwohl das Lastenheft sie ausnimmt; die Decken-Regel (Spezifikation präzisiert das Lastenheft, erweitert es nie, [`AGENTS.md`](../../AGENTS.md) §2) würde verletzt. Das ist die Klasse „Bedingung ohne Code-Entsprechung“ eine Schicht tiefer ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)).

**Auftrag (Planner, Closure oder Anlage der Folge-Slices):** in die Welle-Zeile `slice-v2-spec-ollama-cloud` und in die Folge-Slice-Pläne die Zeile „Z4 (Neuanlage mit demselben Volume) ist im Lastenheft `0.23.2` nicht zugesagt; Spezifikation und Handbuch führen sie als Handvariante ohne Zusage; die ADR bleibt (Accepted), die Abweichung ist hier benannt“ aufnehmen. Dasselbe für den Aufbau-Rest nach L-1 (Prüfung 2, Bedingung 3). Eine Änderung der ADR ist nicht nötig und nicht erlaubt; ein Folge-ADR nur, wenn der Projektinhaber die Neuanlage wieder zusagen will. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

---

## Weitere Beobachtungen

- **V-6 (INFO) Prüfzustand im Zielbild.** „Zugesagt ist der Ablauf nur in dem Aufbau, in dem er geprüft ist“ nennt einen Prüfzustand; die Vorgabe vom 2026-10-19 schließt Messstand, Messwert, Prozesszustand aus dem Lastenheft aus. Der Reviewer hat es als INFO geführt (N-6); der Projektinhaber hat den Wortlaut zweimal freigegeben. Hinweis für die Spec-Schicht: die Aufbau-Bedingungen stehen im Wortlaut ausformuliert, sodass der Halbsatz „in dem er geprüft ist“ streichbar wäre, ohne dass die Zusage sich ändert. Keine Aktion in diesem Slice.
- **ADR-Abweichung F3 ausgewiesen.** ADR-*Folgepflicht* („der Messstand im Absatz entfällt“) und Slice (Messstand bleibt in Zeile 412) weichen ab; die Abweichung ist als Planänderung in §1 (Frage F3), in der Welle (Frage 3) und im Carveout 2 benannt und freigegeben. Kein Befund.
- **Sprache und Kommentar-Selbstcheck.** Lastenheft-Text deutsch; die neuen Zeilen enthalten kein „früher“/„no longer“; die Historie-Zeilen `0.23.1` und `0.23.2` beschreiben Änderungen einer Version (Chronik). Kein Befund zu [`AGENTS.md`](../../AGENTS.md) §3.7/§3.9.

---

## Carveouts (verification.md-Feld)

- **Neu:** keine zusätzlichen durch den Verifier. Der Slice hat zwei temporäre Carveouts im Inventar angelegt (Zeilen in `carveouts.md`, Plan-Anker Welle §4 und `slice-cr-lastenheft-messstaende`). <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- **Gelöst:** keine.
- **Unverändert:** der Eintrag zu `u-boot logs --json` (CRLF).

## Nicht ausgeführt

- `make test-docker`, `make verify-depguard`, `make ci` — Markdown-only, kein Docker-Pfad, keine Importregel.
- Rohdaten einzeln: nicht alle geöffnet; gelesen wurden Protokoll ohne Konto Abschnitt F und die Tabelle der Handvarianten sowie die Z-Zeilen der ADR. `rohdaten-06-restricted-1-114624.txt` (Z4) nicht geöffnet, weil das Lastenheft die Neuanlage nicht mehr zusagt.
- Nutzer-Doku (`README.md`, `docs/user/`): nicht geändert, daher nicht geprüft außer den Fundstellen der Spannung.

## Findings

```text
LOW  slice §7 / §6 · L-1 (Nachlauf 2) · Entscheidung über den Aufbau-Rest des Schlüssel-Satzes fehlt
Befund: Annahme als bekannte Grenze vertretbar (Prüfung 2), aber nur mit Entscheidung des Projektinhabers, §7-Eintrag, §6-Ausgang und Folgepflicht an die Folge-Slices. Ohne diese drei Dinge bleibt DoD Punkt 2 offen.
Zuständig: Planner (Frage an den Projektinhaber), Implementer (§7-Eintrag nach Entscheidung)

LOW  slice §6, §1 Kopf Fassung 3 · V-2 · Zustandsmarker mit abgelaufenem Ablaufpunkt (Risiken 1, 4, 5; „Es folgen das Nachlauf-Review … und der Verifier“)
Zuständig: Planner (Closure)

LOW  slice §7 / §6 Risiko 4 · V-3 · Beleg-Tabelle in §7 fehlt; Verweis „Welle §5“ trifft §4
Zuständig: Planner (Closure)

LOW  Welle §4 Zeile slice-v2-spec-ollama-cloud, carveouts.md · V-4 · Z4 (ADR Zugesagt) im Lastenheft nicht mehr zugesagt, nicht als Folgepflicht der Spec-Schicht ausgewiesen
Zuständig: Planner (Welle-Zeile, Folge-Slice-Pläne), Architect (Hinweis)

INFO slice-cr-lastenheft-messstaende.md Zeile 29 · V-5 · „Lastenheft 0.23.0 erzeugt“ gegenüber Stand 0.23.2
INFO Lastenheft Zeile 414/418 · V-6 · „in dem er geprüft ist“ als Prüfzustand (freigegeben; N-6)
INFO Freigabe · V-7 · Freigabe-Eintrag und Lastenheft-Commit im selben Commit; Vorlage an den Projektinhaber aus dem Repo nicht prüfbar
INFO Carveout 1 · V-8 · Plan-Anker ist eine Welle-Zeile ohne Datei der Folge-Slices; beim Anlegen der Dateien auf diese ziehen
```

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 4 (V-1 = L-1, V-2, V-3, V-4) |
| INFO | 4 (V-5 bis V-8) |

## Verdikt

**Verifikation: bestanden mit Auflagen vor `done/`.** Das Lastenheft `0.23.2` ist Zeichen für Zeichen mit dem freigegebenen Wortlaut identisch (alle acht Stellen, drei Fassungen), trägt keine Kennung, keinen Messstand und keinen fremden Hostnamen, und kein Satz geht über Z1 bis Z13 und das Protokoll hinaus, außer dem Aufbau-Rest L-1. `make gates`, `make doc-trace`, `make doc-commits`, `make doc-immutable` sind grün, ADR-0036 und ADR-0024 sind unverändert. Alle Review-Befunde sind behoben oder mit Grund angenommen, bis auf L-1. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

**Vor der Verschiebung nach `done/` nötig:**
1. L-1 entscheiden (Projektinhaber): Annahme mit §7-Eintrag und §6-Ausgang (Standard-Empfehlung des Verifiers) oder Fassung 4 mit Freigabe und Nachlauf 3.
2. §7 *Beleg-Tabelle* ergänzen oder auf die Herkunftstabellen und den Abschnitt „Satz-gegen-Zusage-Menge“ dieses Reports verweisen (V-3); Zustandsmarker in §6 und §1 auf den Ist-Stand setzen, jedes Risiko mit Ausgang (V-2).
3. Z4-Hinweis und die Folgepflicht aus L-1 in die Welle-Zeile `slice-v2-spec-ollama-cloud` und `slice-v1-devcontainer-ollama-cloud` (V-4). <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

**Hinweis zur Freigabe:** Dass der Projektinhaber die Vorlage ungekürzt sah und die Freigabe vor der Änderung gab, ist aus dem Repo nicht prüfbar; der Verifier bestätigt die Übereinstimmung von Lastenheft und §1/§7, nicht das Gespräch.

**Übergabe:** an den Planner (Closure, Frage zu L-1, V-2 bis V-5, V-8). Der Code-Slice erbt: Z13, Ausweisung, Warnung, `--quiet`-Sichtbarkeit, Nicht-Hängen der Egress-Sätze-Tests an der Egress-losen Lage.
