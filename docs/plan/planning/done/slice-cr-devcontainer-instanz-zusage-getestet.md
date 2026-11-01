# Slice CR: Zusage zu gleichzeitigen Instanzen auf das eingrenzen, was die Suite belegt — Lastenheft-CR zu `LH-FA-DEV-010`, Wortlaut zur Freigabe des Projektinhabers

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Variante 1 vom Projektinhaber freigegeben (§7). Zweiter von zwei Slices: ihm geht [`slice-v1-devcontainer-instanz-test-leicht`](../done/slice-v1-devcontainer-instanz-test-leicht.md) voraus (Wege entschieden, Entscheidung Projektinhaber 2026-11-01: zwei Instanzen, Workspace-Isolation, kein Messskript). Vorbild des Ablaufs: [`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Änderung des Absatzes zu den Instanzen, Zeile 504).
Anlass: Entscheidung des Projektinhabers, 2026-11-01: Der Docker-Test mit vier Instanzen ist nicht praxistauglich und entfällt ersatzlos (kein Messskript); „Workspace-Isolation bei zwei Instanzen reicht“. Die Suite belegt danach weniger, als der Absatz zusagt.

**Berührte Spec-Stellen:** `lastenheft.md` Absatz zu den Instanzen in [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), §16 (Historie), Versionsfeld (Zeile 9). [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes): gelesen Zeile 837 (Absatz „Instanzen“); Suche nach `3 GiB`, `vier gleichzeitig`, `bis zu vier` in `spec/spezifikation.md` ohne Treffer, die Spec trägt Zahl und Grenze nicht. Weitere Zeilen von [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (838 bis 844, 892, 900) vom Implementer zu lesen, bevor sie als „unverändert“ geführt werden (`AGENTS.md` §3.10).

**Verantwortlich:** pt9912

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-11-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Absatz zu den Instanzen in [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) sagt nur zu, was die Suite (`integration-docker`) belegt; was nur die Messung vom 2026-10-31 belegt, steht nicht mehr als Zusage im Lastenheft. Der Wortlaut ist vom Projektinhaber freigegeben, bevor er geschrieben wird.

### Gelesener Stand (2026-11-01)

Der Absatz steht in `spec/lastenheft.md` Zeile 504 (Lastenheft `0.25.0`, Zeile 9; Historie-Zeile `0.25.0` Zeile 1609, die letzte Zeile mit Version `0.2x`). Zwei Sätze tragen die Zusage:

1. Struktur: „Bis zu vier gleichzeitig laufende Instanzen desselben Projekts teilen im Sandbox-Profil kein Volume: Jede Instanz hat ein eigenes Workspace-Volume und, soweit aktiviert, ein eigenes Volume für den Docker-Speicher (`nestedRuntime: docker`) und für die Modelle des Features `ollama`; eine Datei im Workspace-Volume einer Instanz ist in keiner anderen Instanz sichtbar.“
2. Zusage: „Zugesagt ist das für zwei, drei und vier gleichzeitig laufende Instanzen, jede mit höchstens einer CPU und 3 GiB Arbeitsspeicher (Grenze samt Swap), im Sandbox-Profil mit `nestedRuntime: docker` und dem Feature `ollama`, im Zustand direkt nach dem Start, auf einem Linux-Host (`amd64`) mit lokaler Docker-Engine. Nicht zugesagt sind mehr als vier Instanzen, der Speicher von Podman (`nestedRuntime: podman`), Instanzen mit mehr oder ohne Begrenzung von CPU und Arbeitsspeicher, das Verhalten unter Last (…), mit kaltem Cache oder bei kaltem Bau des Images, Dauer und Platzbedarf des Starts, andere Hosts und Engines und die Isolation im Standardprofil: Dort teilen alle Instanzen den Projektordner (Bind-Mount).“

Belegt war die Zusage bisher durch die Messung vom 2026-10-31 (`docs/messungen/slice-v1-devcontainer-instanz-messung/`: zwei, drei und vier Instanzen, Stufen 2 bis 4) und den Docker-Test mit vier Instanzen, der ersatzlos entfernt ist (Entscheidung Projektinhaber 2026-11-01).

Seit [`slice-v1-devcontainer-instanz-test-leicht`](../done/slice-v1-devcontainer-instanz-test-leicht.md) belegt die Suite (Weg A, gelesen in `internal/e2e/devcontainer_instance_docker_test.go`, `TestDevcontainerInstance_Docker_TwoInstancesAreIsolated`): Sandbox-Profil mit Standard-Einstellungen, Container ohne Instanz plus **zwei** Instanzen, ein Workspace-Volume je Container, Marker-Isolation; der Lauf `integration-docker` auf `main` ist erfolgreich (Stand 2026-11-01, §7). **Nicht** belegt sind Docker-Speicher- und Modell-Volumes (nur Text-Pins der Mounts), die Grenze von einer CPU und 3 GiB sowie drei und vier Instanzen. Der Absatz geht also in vier Merkmalen (`nestedRuntime: docker` und `ollama`, Grenze, drei und vier Instanzen) über die Suite hinaus.

### Wortlaut-Vorschläge (zur Freigabe; NICHT im Lastenheft geschrieben)

Die Wege des Test-Slice sind entschieden: Weg A mit zwei Instanzen neben dem Container ohne Instanz.

**Variante 1 (Eingrenzen auf die Suite; Vorschlag).** Satz 1 bleibt (er beschreibt die Bauart, die Mounts sind als Text gepinnt; die Zahl „bis zu vier“ im Satz 1 ist vom Projektinhaber mit zu entscheiden, siehe Frage 2). Satz 2 lautet:

> Zugesagt ist die Isolation der Workspace-Volumes für zwei gleichzeitig laufende Instanzen im Sandbox-Profil, im Zustand direkt nach dem Start, auf einem Linux-Host (`amd64`) mit lokaler Docker-Engine. Nicht zugesagt sind drei und mehr Instanzen, die Isolation des Docker-Speichers (`nestedRuntime: docker`) und der Modelle (`ollama`), der Speicher von Podman (`nestedRuntime: podman`), Instanzen mit oder ohne Begrenzung von CPU und Arbeitsspeicher, das Verhalten unter Last (…), mit kaltem Cache oder bei kaltem Bau des Images, Dauer und Platzbedarf des Starts, andere Hosts und Engines und die Isolation im Standardprofil: Dort teilen alle Instanzen den Projektordner (Bind-Mount).

Gegen die Messung gelesen: „zwei“ ist gemessen (Stufe 2, zwei Läufe) und von der Suite belegt (zwei Instanzen plus Container ohne Instanz; „zwei“ zählt die Instanzen, nicht den Container ohne Instanz). Das Wort „Isolation der Workspace-Volumes“ ist Satz 1 entnommen; die Messung belegt für Workspace-Volumes paarweise verschiedene Namen und die Inhaltsprobe (`BEO-DOK/zusage-weiter-als-messung`: jedes Wort an der Messung gegengelesen, vom Implementer in §7 zu belegen).

**Variante 2 (Zusage wie bisher, Beleg ausdrücklich außerhalb der Suite).** Der Wortlaut bleibt; es entsteht kein CR, der Slice entfällt (Gegenstand: entfallen, Grund: die Messung vom 2026-10-31 ist der Beleg). Das Lastenheft sagt dann weiter zu, was kein Gate der CI belegt (es gibt kein Messskript); `docs/maintainer/zusagen-und-messstaende.md` nennt die Messung als Beleg. Risiko: eine Zusage ohne Gate widerspricht dem Muster der übrigen Zusagen mit Test.

**Variante 3 (Zahl und Grenze aus der Zusage, Bauart bleibt).** Satz 2 nennt weder Zahl noch Grenze: „Zugesagt ist die Isolation für gleichzeitig laufende Instanzen im Sandbox-Profil … Zahl und Grenze je Instanz sind nicht zugesagt.“ Der Wortlaut ist knapper, trägt aber weniger als die Messung hergibt und nimmt dem Nutzer die Auskunft „vier mit 1 CPU und 3 GiB sind gemessen“; diese Auskunft stünde nur noch in `docs/maintainer/zusagen-und-messstaende.md`.

### Fragen für den Projektinhaber

1. **Welche Variante** (1, 2, 3) oder ein eigener Wortlaut? Variante 1 ist ohne Wortlaut-Freigabe nicht zu schreiben.
2. **Welche Zahl von Instanzen** nennt der Satz: „zwei“ (Variante 1, wie die Suite sie belegt), „bis zu zwei“, oder bleibt „bis zu vier“ in Satz 1 stehen, wenn Satz 2 nur zwei zusagt? Als Alternativen für den Wortlaut bleiben ferner: zwei Instanzen mit der Messung als genannter Beleg (Satz „gemessen sind zwei, drei und vier“ als Auskunft, nicht als Zusage), und der Wegfall von Zahl und Grenze (Variante 3).
3. **Steht die Auskunft über Docker-Speicher/Modelle und die Grenze weiter im Lastenheft** (als „nicht zugesagt“, wie Variante 1) **oder gar nicht?**
4. **Version:** `0.26.0` (nächste Minor nach dem Muster `0.25.0`), falls nicht inzwischen von einem anderen Entwurf mit Lastenheft-CR vergeben; die Version ist beim Start aus Zeile 9 und der Historie zu lesen und nicht vorab zu reservieren.

### Folgen (gelesen, soweit angegeben; sonst vom Implementer zu lesen)

| Stelle | Aussage | Folge |
|---|---|---|
| `docs/maintainer/zusagen-und-messstaende.md` Zeilen 168 bis 181 | wiederholt die Zusage mit „vier gleichzeitig … 3 GiB … nestedRuntime: docker und dem Feature ollama“ und nennt den schweren Test als Fundstelle (Zeile 178) | folgt dem freigegebenen Wortlaut; die Fundstelle (Test) hat `slice-v1-devcontainer-instanz-test-leicht` bereits auf den leichten Test gesetzt, der CR zieht den Zusage-Teil nach |
| `spec/spezifikation.md` [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), Zeile 837 und folgende | Instanzen, Namensregel, Suche, Exit-Code `10`; keine Zahl, keine Grenze (Suche ohne Treffer) | voraussichtlich unverändert; der Implementer liest 837 bis 844, 892 und 900 und trägt die Fundstelle ein |
| `docs/user/benutzerhandbuch.md` Zeilen 888 bis 922 | Zeilen 903 bis 906: „Im Sandbox-Profil hat sie außerdem eigene Volumes: den Workspace und, soweit aktiviert, den Docker-Speicher (`nestedRuntime: docker`) und die Modelle des Features `ollama`. Eine Datei im Workspace einer Instanz sehen die anderen Instanzen nicht.“ | Bedingung über das Verhalten des Codes („soweit aktiviert“): braucht einen Test, der sie pinnt, oder wird gestrichen (`AGENTS.md` §3.10, letzter Punkt); Pin: Text-Pins der drei Mounts (`generate_ollama_test.go` Zeilen 27 bis 30) und der leichte Test für den Workspace; ob das reicht, entscheidet der Implementer mit Fundstelle. Das Handbuch trägt keine Kennung (`AGENTS.md` §5 Regel 6) |
| `CHANGELOG.md` Zeile 16 (`[Unreleased]`) | „Pinned by unit tests and by Docker tests with two instances and the container without an instance in the sandbox profile …“ (von `slice-v1-devcontainer-instanz-test-leicht` gesetzt) | der CR berührt sie nur, wenn der Wortlaut sie ändert |
| `spec/lastenheft.md` §16 | Historie-Zeile `0.25.0` nennt die Zusage („Isolation der Volumes für bis zu vier gleichzeitige Instanzen … jede mit höchstens einer CPU und 3 GiB …“) | wird nicht umgeschrieben; eine neue Zeile `0.26.0` nennt die Eingrenzung (Muster: Historie-Zeilen sind additiv) |

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Umbau der Tests** — [`slice-v1-devcontainer-instanz-test-leicht`](../done/slice-v1-devcontainer-instanz-test-leicht.md); der CR beschreibt das, was die Suite danach belegt.
- **Eine neue Messung** — der CR sagt nur zu, was gemessen oder von der Suite belegt ist; neue Messungen sind eigene Vorgänge.
- **Änderung von Produktcode oder von [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)** — der CR ändert keine Verhaltenszusage, er nimmt eine zurück; kein Verhalten ändert sich.
- **Eine Änderung der Bauart-Aussage (Satz 1)**, solange Variante 1 gilt — sie bleibt, weil die Mounts als Text gepinnt sind.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Wortlaut ohne offene Entscheidung:** Der Wortlaut ist dem Projektinhaber als Alt/Neu mit Alternativen vorgelegt und freigegeben; die Freigabe (Datum, gewählte Variante) ist in §7 vermerkt, bevor das Lastenheft geschrieben wird.
- [x] **Lastenheft** in der nächsten freien Version (`0.26.0`, falls frei) geschrieben: Absatz zu den Instanzen in der freigegebenen Fassung, Versionsfeld (Zeile 9), Historie-Zeile mit Datum (Verweis „Vereinbarung mit dem Projektinhaber“).
- [x] **Folgen nachgezogen:** `docs/maintainer/zusagen-und-messstaende.md` (Zusage-Teil), Handbuch und Spec nur, wo die Lesung (Tabelle „Folgen“) eine Aussage findet, die der Wortlaut nicht mehr trägt; jede „unverändert“-Zeile mit Fundstelle der gelesenen Stelle.
- [x] **Vorab-Suche und Beleg-Tabelle:** `grep -i` über `spec/lastenheft.md` nach `Instanz`, `vier`, `3 GiB`, `ollama`, `nestedRuntime`, `Swap` ist in §7 festgehalten; jede Bedingung im neuen Wortlaut hat einen Test oder eine Messung als Beleg (Tabelle Bedingung → Beleg, `AGENTS.md` §3.10, letzter Punkt).
- [x] `make gates` grün; `make doc-trace` ohne neue Waise. Läufe und Exit-Codes in §7.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8); der Review läuft über die **freigegebene** Fassung.
- [x] Doku-Update: `CHANGELOG.md` nur, wenn der Wortlaut Zeile 16 betrifft; kein Produktverhalten.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Absatz Zeile 504 in der freigegebenen Fassung; Versionsfeld Zeile 9; neue Historie-Zeile hinter Zeile 1609; gelesen vor der Änderung: Zeilen 497 bis 508, 9, 1609 |
| `docs/maintainer/zusagen-und-messstaende.md` | update | Zusage-Teil des Abschnitts „Mehrere Devcontainer je Projekt“ (Zeilen 168 bis 181) |
| `docs/user/benutzerhandbuch.md` | update, falls die Lesung es verlangt | Zeilen 903 bis 906 (Bedingung „soweit aktiviert“), nur mit Fundstelle der Lesung |
| `spec/spezifikation.md` | unverändert für den CR-Anteil | [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) Zeilen 837 bis 844, 892 und 900 gelesen: keine Zahl, keine Grenze. `HEAD` trägt in der Datei Änderungen aus Änderung A (`9fec8752`, Unterabschnitt `shell` auf die Engine 29.9.2) und aus der Fix-Runde `1bf451c4`; sie gehören zu [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md), nicht zum Wortlaut des CR (Verification V-2) |
| `docs/plan/planning/in-progress/roadmap.md` | update | Zeiger auf diesen Slice unter *Offene Wellen* |
| `docs/plan/planning/done/slice-v1-devcontainer-instanz-test-leicht.md`, `docs/reviews/2026-11-01-slice-v1-devcontainer-instanz-test-leicht.md`, `docs/reviews/2026-11-01-slice-v1-devcontainer-instanz-test-leicht-verify.md` | update (Werkzeug) | Verweis-Nachzug der `slice-mv`-Commits (`git show --stat` von `a8f0ece1`, `c077a6a0`, `8296197c`, `44d14987` gelesen); kein Inhalt des CR |
| `docs/reviews/2026-11-01-slice-cr-devcontainer-instanz-zusage-getestet.md` | neu | Review-Report dieses Slice |
| `docs/plan/planning/done/slice-cr-devcontainer-instanz-zusage-getestet.md` (diese Datei) | update | Wortlaut-Entwurf, Freigabe und Vorab-Suche in §7; Zustandsmarker auf Ist-Stand |
| `docs/plan/planning/open/README.md` | update | Zeile dieses Slice; beim Move Verweis-Nachzug (`make slice-mv`) |
| `docs/reviews/` | neu | Review und Verification nach `harness/review.md` und `harness/verification.md` |
| `docs/plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/` | update | Closure (§7): Evidence-Datei dieses Slice und `state.md` |
| `internal/**`, `spec/architecture.md`, `README.md`, `README.de.md`, `CHANGELOG.md` | unverändert für den CR-Anteil | der CR ändert kein Produktverhalten; `CHANGELOG.md` Zeile 16 nennt den Test (gelesen), dessen Nachzug gehört dem Test-Slice. In `internal/` berühren `9fec8752` und `1bf451c4` den Doc-Kommentar von `internal/hexagon/port/driven/devcontainer_launcher.go` und `internal/adapter/driven/docker/devcontainer_shell.go` (Änderung A, Review F-1 und F-3); sie gehören zu [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md) (Verification V-2) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; [`slice-v1-devcontainer-instanz-test-leicht`](../done/slice-v1-devcontainer-instanz-test-leicht.md) liegt in `done/` (die Wege 1 bis 3 stehen fest und der Lauf der CI belegt, was die Suite belegt); der Projektinhaber hat Variante und Wortlaut freigegeben (Vermerk in §7); WIP-Limit 1.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Der Wortlaut verlangt mehr als eine Anforderung oder eine neue Entscheidung des Projektinhabers.
- `in-progress` → `open` (blockiert): keine Freigabe des Wortlauts; das Lastenheft bleibt unverändert.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). Der Lastenheft-Commit trägt die Freigabe nur, wenn §7 sie vermerkt.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Der neue Wortlaut sagt mehr zu, als Messung und Suite belegen** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)): etwa „zwei Instanzen“, die nur die Messung belegt, oder „Isolation“ als Oberbegriff, wo nur Workspace-Volumes getestet sind. — **Ausgang:** entfallen: Die Beleg-Tabelle in §7 ist der Abgleich jedes Worts von Satz 2 gegen Test und Zusage-Menge; Review (Prüfpunkt 1) und Verifier (maschineller Vergleich gegen die Freigabe) fanden keine Aussage weiter als die Belege. Der Wortlaut bleibt als Beleg in §7 stehen (Verification-Report, DoD-Punkt 2).
- **Die Suite belegt „zwei“ erst nach einem grünen Lauf** von `TestDevcontainerInstance_Docker_TwoInstancesAreIsolated` auf dem Runner (der Test ist nicht gelaufen). — **Ausgang:** entfallen: Lauf 41285105546 war erfolgreich, das Job-Log zeigt `PASS` des Tests, `internal/e2e` ist seither unverändert (§7, Verification-Report).
- **Ein Satz im Handbuch behauptet Isolation von Docker-Speicher und Modellen** (`docs/user/benutzerhandbuch.md` Zeilen 903 bis 906) und ist von der Suite nicht pinnbar. — **Ausgang:** entfallen: Der Satz ist ersetzt („wenn aktiviert“, ohne Isolationsaussage) und durch die Text-Pins `dockerMount` und `ollamaMount` gepinnt (§7, Folgen).
- **Die Version `0.26.0` ist vergeben**, wenn der Slice startet: andere Entwürfe in `open/` mit Lastenheft-CR können sie früher nehmen. — **Ausgang:** entfallen: `0.26.0` war bei Start frei (Zeile 9 und Historie gelesen) und ist geschrieben.
- **Test-Lücken der Eingrenzung von `shell` auf die Engine 29.9.2** (Review F-7): (a) der Fall „die Sitzung schreibt die Zeile selbst“ ist von keinem Test gepinnt (Lesung des Codes); (b) die Form der Engine 28.0.4 pinnen nur Einheitstests mit Fake, kein Lauf gegen eine solche Engine; (c) ob der GitHub-Runner die Engine 29.9.2 stellt, ist nicht geprüft. — **Ausgang:** weiter offen (Anker: Review-Report F-7 und `docs/maintainer/zusagen-und-messstaende.md` Abschnitt 3.15; Verification V-3). Die Lücken gehören zu Änderung A, nicht zum Wortlaut des CR.
- **Das Lastenheft ist Vertrag und wird enger, nicht weiter:** eine Rücknahme einer Zusage ist eine Vertragsänderung, auch ohne neues Verhalten. — **Ausgang:** entfallen: der Slice ist als CR mit Freigabe des Projektinhabers geschnitten.
- **Eine Accepted-ADR trägt den Satz:** [`ADR-0038`](../../adr/0038-devcontainer-instanz-zweites-id-label.md) wurde nicht auf Zahl und Grenze der Zusage gelesen. — **Ausgang:** entfallen: Zeilen 73 und 97 gelesen (§7, Folgen); die Aussage „die Zusage nennt eine Menge“ bleibt wahr, die ADR ist nicht im Diff (Verification-Report).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren) · `grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln.

Reports: [Review](../../../reviews/2026-11-01-slice-cr-devcontainer-instanz-zusage-getestet.md) (0 HIGH, 0 MEDIUM, 3 LOW, 5 INFO; F-1 bis F-4 und F-7 in der Fix-Runde `1bf451c4` berichtigt oder benannt), [Verification](../../../reviews/2026-11-01-slice-cr-devcontainer-instanz-zusage-getestet-verify.md) (0 HIGH, 0 MEDIUM, 0 LOW, 4 INFO). Gegenstand: Lastenheft `0.26.0` (`d38c0908`), Freigabe `3ec5cda5`.

Vor dem Lastenheft-Commit steht hier der Vermerk der Freigabe (Datum, Variante, Wortlaut) und der Entscheidungen zu den Fragen 1 bis 4 aus §1.

**Freigabe und Entscheidung des Projektinhabers (2026-11-01): Variante 1**, Satz 2 im Wortlaut von „Variante 1“ in §1: „Zugesagt ist die Isolation der Workspace-Volumes für zwei gleichzeitig laufende Instanzen im Sandbox-Profil, im Zustand direkt nach dem Start, auf einem Linux-Host (`amd64`) mit lokaler Docker-Engine. Nicht zugesagt sind drei und mehr Instanzen, die Isolation des Docker-Speichers (`nestedRuntime: docker`) und der Modelle (`ollama`), der Speicher von Podman (`nestedRuntime: podman`), Instanzen mit oder ohne Begrenzung von CPU und Arbeitsspeicher, das Verhalten unter Last (…), mit kaltem Cache oder bei kaltem Bau des Images, Dauer und Platzbedarf des Starts, andere Hosts und Engines und die Isolation im Standardprofil: Dort teilen alle Instanzen den Projektordner (Bind-Mount).“ Die Fragen aus §1 sind entschieden:

1. Variante 1.
2. „zwei“ auch in Satz 1 (statt „bis zu vier“).
3. Die Auskunft über Docker-Speicher, Modelle und die Grenze bleibt als „nicht zugesagt“ im Lastenheft.
4. Version `0.26.0` (Zeile 9 und Historie bei Start gelesen: `0.25.0`, keine höhere vergeben).

Beleg für „zwei“ in der Suite: der Lauf `integration-docker` auf `main` (Lauf 41285105546, Kopf `a6264402`) war erfolgreich; sein Kopf enthält den Commit `73b42e21` mit `TestDevcontainerInstance_Docker_TwoInstancesAreIsolated` (`git merge-base --is-ancestor`, gelesen). Das Job-Log ist über `gh api repos/pt9912/u-boot/actions/jobs/124784712075/logs` lesbar und zeigt `--- PASS: TestDevcontainerInstance_Docker_TwoInstancesAreIsolated (16.32s)` (Lauf 41285105546, Kopf `a6264402`; vom Reviewer gelesen, `docs/reviews/2026-11-01-slice-cr-devcontainer-instanz-zusage-getestet.md`).

**Vorab-Suche** (`grep -i` in `spec/lastenheft.md` nach `bis zu vier`, `vier`, `1 CPU`, `einer CPU`, `3 GiB`, `Swap`, `Instanz`, am 2026-11-01): Treffer mit Zahl oder Grenze nur in Zeile 504 (Absatz, geändert) und Zeile 1609 (Historie `0.25.0`, unberührt, additiv ergänzt durch Zeile `0.26.0`). Die übrigen `Instanz`-Treffer (501 bis 503, 507, 509) tragen weder Zahl noch Grenze. [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung) und Zeile 416 tragen keine Zusage zu Instanzen und sind nicht berührt.

**Folgen, gelesen:** `spec/spezifikation.md` [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) Zeilen 837 bis 844, 892 und 900 gelesen: keine Zahl, keine Grenze, unverändert. `CHANGELOG.md` Zeile 16 (`[Unreleased]`) gelesen: nennt Docker-Tests mit zwei Instanzen und keine Zusage über vier; unverändert. [`ADR-0038`](../../adr/0038-devcontainer-instanz-zweites-id-label.md) Zeile 73 („die Zusage im Lastenheft nennt die gemessene Menge“) und Zeile 97 („mehr als die gemessene Zahl von Instanzen“) sind Accepted und nicht geändert; die Menge ist jetzt „zwei“, die Aussage bleibt wahr in der Form „die Zusage nennt eine Menge“. Handbuch Abschnitt 4.9: der Satz „soweit aktiviert“ ist ersetzt durch „wenn aktiviert“ für Docker-Speicher und Modelle, gepinnt durch die Text-Pins `dockerMount` und `ollamaMount` in `internal/hexagon/application/generate_ollama_test.go` (Zeilen 363 und 365: Mount nur bei aktivierter Option); ohne Aussage zur Isolation dieser Volumes.

**Beleg-Tabelle (Bedingung → Beleg):**

| Bedingung im neuen Wortlaut | Beleg |
|---|---|
| zwei gleichzeitig laufende Instanzen, Workspace-Volumes getrennt, Datei nicht in anderer Instanz sichtbar | `TestDevcontainerInstance_Docker_TwoInstancesAreIsolated` (`internal/e2e/devcontainer_instance_docker_test.go`); `integration-docker` auf `main` erfolgreich, Kopf enthält den Test (siehe oben) |
| Sandbox-Profil | derselbe Test (Sandbox-Profil mit Standard-Einstellungen); im Standardprofil `TestDevcontainerInstance_Docker_StandardProfileSharesTheProjectFolder` |
| im Zustand direkt nach dem Start | der Test prüft nach dem Start, ohne Last |
| Linux `amd64`, lokale Docker-Engine | Zusage-Menge von [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) und Runner der CI; kein Test pinnt die Architektur |
| Satz 1: eigene Volumes für Docker-Speicher und Modelle, soweit aktiviert | Text-Pins der Mounts, `generate_ollama_test.go` (Zeilen 363 und 365); keine Zusage der Isolation (Satz 2: nicht zugesagt) |

**Läufe** (Verifier, HEAD `1bf451c4`, Docker-Modus, selbst gefahren; Verification-Report, Kopf): `make gates` Exit 0, `coverage-gate: OK — coverage 92.20% meets threshold 90%`, `baseline-verify: v6.18.1 OK — 54 Dateien`, `ablage-pruefung`: 26 Fälle, 0 Fehlschläge, `docs-check` Teil von `gates`. `make doc-trace` Exit 0: 84 Anforderungen, 0 Waisen; [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) führt diesen Slice in der Slice-Spalte, Status ok. Nicht ausgeführt: `make test-docker` (Entscheidung des Projektinhabers; Ersatz ist der gelesene CI-Lauf 41285105546), `make verify-depguard` (keine Importregel berührt). Der Lauf von `make gates` nach der Closure steht im Commit der Closure.

- **Was hat funktioniert:**
  - **Die Beleg-Tabelle (Bedingung → Beleg) trennte Belegtes von Übernommenem:** Jedes Merkmal von Satz 2 steht im Test (`TwoInstancesAreIsolated`) oder in der Zusage-Menge der Anforderung; Review und Verifier verglichen den Wortlaut maschinell mit der Freigabe in §7 (Teilstring wahr) und fanden keine Aussage weiter als die Belege.
  - **Die Freigabe lag vor dem Lastenheft-Commit:** `3ec5cda5` (Freigabe, Variante 1, Fragen 1 bis 4) vor `d38c0908` (Lastenheft `0.26.0`); der Review lief über den freigegebenen Wortlaut, nicht danach (anders als bei [`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md)).
  - **Der CI-Beleg wurde im Job-Log gelesen**, nicht aus dem Status geschlossen (Review F-2 berichtigte §7).
- **Was ging anders als geplant:**
  - **Der Review umfasste zwei Änderungen zugleich:** den CR und Änderung A (`9fec8752`, [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) `shell` auf die Engine 29.9.2). Die Plan-Zeilen „unverändert“ galten nur für den CR-Anteil, `HEAD` trug Änderungen aus beiden (Verification V-2); §3 sagt es jetzt so.
  - **Die Plan-Tabelle führte berührte Dateien nicht** (Review F-3: `devcontainer_launcher.go`; F-4: Werkzeug-Nachzüge von `slice-mv`) und §7 nannte den Beleg schwächer als vorhanden (F-2). Alle in der Fix-Runde `1bf451c4` berichtigt.
  - **Die Tür zu Änderung A blieb offen:** Die Zusage von `shell` auf der Engine 29.9.2 hat keinen Lauf auf einer solchen Engine im Beleg (Review F-7, Verification V-3); sie ist als „weiter offen“ geführt (§6).
  - **Prozess:** Der Verifier fand `.harness/skills/verifier.md` nicht (Verification-Report, Kopf); die Rolle folgte `harness/roles.md` und `harness/verification.md`.
- **Steering-Loop-Eintrag:** Lehre (zweiteilig): (a) Entfällt ein Test, der eine Zusage belegt, ist die Zusage im selben Vorgang zu prüfen; der Test-Slice [`slice-v1-devcontainer-instanz-test-leicht`](../done/slice-v1-devcontainer-instanz-test-leicht.md) ließ das Lastenheft mit der Zusage für vier Instanzen stehen, bis dieser CR sie auf die Suite zurücknahm. (b) Hängen zwei Änderungen im selben Review, gilt jede „unverändert“-Zeile des Plans nur für den eigenen Anteil, und die Zeile sagt es (V-2). Gezählt, nicht verkörpert: Es wurde keine Regel, kein Sensor und keine Spec-Lücke verkörpert. Teil (a) zählt in [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (die Zusage von `0.25.0` reichte nach dem Wegfall des Tests weiter als die Suite; der neue Wortlaut tut es nach Review und Verifier nicht); Teil (b) und die Plan-Tabellen-Befunde F-2 bis F-4 gehören in die Klasse [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (bereits verkörpert in `AGENTS.md` §3.10); eine Evidence-Datei dort ist für diesen Slice nicht angelegt, der Zähler dieser Beobachtung ist nicht neu gelesen.
- **Beobachtungs-Register (`../observations/`):** keine neue Beobachtung angelegt; zwei Evidence-Dateien `evidence/slice-cr-devcontainer-instanz-zusage-getestet.md` ergänzt, `state.md` beider nachgeführt, kein Zähler und kein Ausgang gesetzt:
  - [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md): Rücknahme der Zusage von `0.25.0` auf das, was die Suite belegt.
  - [`BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf`](../observations/BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf/observation.md): Die Kombination steht im Lastenheft nur noch als „nicht zugesagt“ (Zeile 504, gelesen: „Nicht zugesagt sind … die Isolation des Docker-Speichers (`nestedRuntime: docker`) und der Modelle (`ollama`)“); die Beobachtung bleibt offen, der Lauf der Kombination fehlt in der Suite weiter.
- **Folge-Slices:** keiner. Die offenen Punkte aus §6 (Test-Lücken von `shell`, Verification V-3 und V-4) hängen am Anker Review F-7 und `docs/maintainer/zusagen-und-messstaende.md` Abschnitt 3.15, nicht an einem Slice.
- **Risiken aus §6:** je Risiko genau ein Ausgang, in der Reihenfolge von §6:
  - Wortlaut sagt mehr zu als Messung und Suite belegen: **entfallen** (Beleg-Tabelle, Review und Verifier ohne Befund dieser Klasse für den Wortlaut des CR).
  - Suite belegt „zwei“ erst nach einem grünen Lauf: **entfallen** (Lauf 41285105546, `PASS` im Job-Log, `internal/e2e` unverändert).
  - Handbuch-Satz zu Docker-Speicher und Modellen: **entfallen** (ersetzt, durch Text-Pins gepinnt).
  - Version `0.26.0`: **entfallen** (bei Start frei, geschrieben).
  - Test-Lücken der Eingrenzung von `shell` (Review F-7): **weiter offen** (Anker: Review-Report F-7 und `docs/maintainer/zusagen-und-messstaende.md` Abschnitt 3.15). Dazu Verification V-3 (Engine des Runners im Job-Log nicht geprüft; das Log nennt nur den Client `docker-ce-cli 29.9.2`) und V-4 (die Bedingung „Linux `amd64`“ steht in der Beleg-Tabelle ohne pinnenden Test, seit `0.25.0` im Lastenheft): **weiter offen** unter demselben Anker; die Architektur ist ungepinnt.
  - Lastenheft wird enger, nicht weiter: **entfallen** (CR mit Freigabe des Projektinhabers).
  - Accepted-ADR trägt den Satz: **entfallen** ([`ADR-0038`](../../adr/0038-devcontainer-instanz-zweites-id-label.md) Zeilen 73 und 97 gelesen, die ADR ist nicht im Diff).
- **Drei Paarungen:** Anker: kein Eintrag mit dem Feld `liegt in` in diesem Slice (nichts verkörpert), nichts zu paaren. Folge-Slice: keiner; die Risiken mit Ausgang „weiter offen“ zeigen auf einen Anker in der Wartungsdoku und im Review-Report, nicht auf einen Slice. Register: die zwei Evidence-Dateien `evidence/slice-cr-devcontainer-instanz-zusage-getestet.md` existieren. Das Paarungen-Häkchen der DoD setzt der Planner nach dem Move.
- **Carveouts:** keine neu, keine gelöst, keine berührt.
- **Abgleich der Plan-Tabelle gegen `git diff --stat`** (Stand nach der letzten Fix-Runde `1bf451c4`, Verification-Report und Closure): Der Verifier glich `3ec5cda5..HEAD` (12 Dateien) gegen §3 ab; die Closure fügt diesen Slice, die zwei Evidence-Dateien, die zwei `state.md`, `open/README.md`, `roadmap.md` und den Verification-Report hinzu, alle in §3 geführt (Beobachtungsregister: Zeile „update“; Verification: Zeile `docs/reviews/`). Der Abgleich gegen den Diff mit der Closure steht im Commit der Closure.
- **Grenzen der Closure:** (1) Die Läufe von `make gates` und `make doc-trace` stammen vom Verifier; der Planner liest sie im Verification-Report. (2) `make test-docker` lief nicht; der Beleg für „zwei Instanzen“ ist der gelesene CI-Lauf. (3) Die Engine des Runners (29.9.2 oder 28.0.4) ist nicht geprüft. (4) Die Architektur `amd64` pinnt kein Test.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/` (Lastenheft), `docs/maintainer/` und `docs/user/`, `docs/plan/` (alle DOK; laut Modus-Deklaration in `harness/conventions.md` eine Sub-Area, die Schwelle ist erfüllt); keine zu grobe.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-DOK` und `BEO-E2E` durchgegangen:
- `BEO-DOK/zusage-weiter-als-messung` — Treffer, Kern dieses Slice: Zähler-Stand aus den Dateien unter `evidence/` (Verzeichnisliste gelesen, u. a. `slice-cr-devcontainer-ollama.md`, `slice-v1-devcontainer-instanz-messung.md`); der Implementer zählt und ergänzt eine Evidence-Datei.
- `BEO-DOK/messvariante-belegt-nicht-die-auslieferung` — Treffer: die Messung nutzt `messen.sh`, das `u-boot devcontainer up` ausdrücklich nicht fährt (Kopfkommentar Zeilen 18 bis 21); die Zusage „vier Instanzen“ stützt sich auf diese Messvariante und auf den schweren Test, der die Auslieferung (`UpService` des Starts) fuhr. Zähler-Stand: vom Implementer zu lesen.
- `BEO-DOK/spec-bedingung-ohne-code-entsprechung` und `BEO-DOK/plan-aussage-ueber-code` — Treffer: das Handbuch behauptet „soweit aktiviert“ (siehe Tabelle „Folgen“); Zähler-Stand vom Implementer zu lesen.
- `BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf` — Treffer: die Kombination hat nach dem Test-Slice keinen Lauf in der Suite; der CR nimmt sie aus der Zusage (Variante 1) oder lässt sie ohne Gate stehen (Variante 2).

### Sub-Area: `spec/`, `docs/plan/` (DOK)

- **Modus:** GF
- **Konventionen-Dichte:** Lastenheft-CR mit Historie-Zeile und Freigabe (Muster `slice-cr-devcontainer-instanz`, `slice-cr-devcontainer-ollama`); ID-Schema [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) fest.
- **Phase-Reife:** Phase 4 (Zusage gemessen, Test entfällt, Wortlaut folgt).
- **Evidenz-/Diskrepanz-Risiko:** mittel: `BEO-DOK/zusage-weiter-als-messung` hat bereits Belege; Gegenmaßnahme: Beleg-Tabelle.
- **Reconciliation-Aufwand:** gering; kein Folge-Slice vorgesehen.
