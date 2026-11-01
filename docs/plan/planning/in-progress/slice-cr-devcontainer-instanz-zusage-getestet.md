# Slice CR: Zusage zu gleichzeitigen Instanzen auf das eingrenzen, was die Suite belegt — Lastenheft-CR zu `LH-FA-DEV-010`, Wortlaut zur Freigabe des Projektinhabers

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. In Arbeit (Variante 1 freigegeben, §7). Zweiter von zwei Slices: ihm geht [`slice-v1-devcontainer-instanz-test-leicht`](../done/slice-v1-devcontainer-instanz-test-leicht.md) voraus (Wege entschieden, Entscheidung Projektinhaber 2026-11-01: zwei Instanzen, Workspace-Isolation, kein Messskript). Vorbild des Ablaufs: [`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md).

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

- [ ] **Wortlaut ohne offene Entscheidung:** Der Wortlaut ist dem Projektinhaber als Alt/Neu mit Alternativen vorgelegt und freigegeben; die Freigabe (Datum, gewählte Variante) ist in §7 vermerkt, bevor das Lastenheft geschrieben wird.
- [ ] **Lastenheft** in der nächsten freien Version (`0.26.0`, falls frei) geschrieben: Absatz zu den Instanzen in der freigegebenen Fassung, Versionsfeld (Zeile 9), Historie-Zeile mit Datum (Verweis „Vereinbarung mit dem Projektinhaber“).
- [ ] **Folgen nachgezogen:** `docs/maintainer/zusagen-und-messstaende.md` (Zusage-Teil), Handbuch und Spec nur, wo die Lesung (Tabelle „Folgen“) eine Aussage findet, die der Wortlaut nicht mehr trägt; jede „unverändert“-Zeile mit Fundstelle der gelesenen Stelle.
- [ ] **Vorab-Suche und Beleg-Tabelle:** `grep -i` über `spec/lastenheft.md` nach `Instanz`, `vier`, `3 GiB`, `ollama`, `nestedRuntime`, `Swap` ist in §7 festgehalten; jede Bedingung im neuen Wortlaut hat einen Test oder eine Messung als Beleg (Tabelle Bedingung → Beleg, `AGENTS.md` §3.10, letzter Punkt).
- [ ] `make gates` grün; `make doc-trace` ohne neue Waise. Läufe und Exit-Codes in §7.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8); der Review läuft über die **freigegebene** Fassung.
- [ ] Doku-Update: `CHANGELOG.md` nur, wenn der Wortlaut Zeile 16 betrifft; kein Produktverhalten.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

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
| `spec/spezifikation.md` | unverändert | [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) Zeilen 837 bis 844, 892 und 900 gelesen: keine Zahl, keine Grenze |
| `docs/plan/planning/in-progress/roadmap.md` | update | Zeiger auf diesen Slice unter *Offene Wellen* |
| `docs/plan/planning/done/slice-v1-devcontainer-instanz-test-leicht.md`, `docs/reviews/2026-11-01-slice-v1-devcontainer-instanz-test-leicht.md`, `docs/reviews/2026-11-01-slice-v1-devcontainer-instanz-test-leicht-verify.md` | update (Werkzeug) | Verweis-Nachzug der `slice-mv`-Commits (`git show --stat` von `a8f0ece1`, `c077a6a0`, `8296197c`, `44d14987` gelesen); kein Inhalt des CR |
| `docs/reviews/2026-11-01-slice-cr-devcontainer-instanz-zusage-getestet.md` | neu | Review-Report dieses Slice |
| `docs/plan/planning/in-progress/slice-cr-devcontainer-instanz-zusage-getestet.md` (diese Datei) | update | Wortlaut-Entwurf, Freigabe und Vorab-Suche in §7; Zustandsmarker auf Ist-Stand |
| `docs/plan/planning/open/README.md` | update | Zeile dieses Slice; beim Move Verweis-Nachzug (`make slice-mv`) |
| `docs/reviews/` | neu | Review und Verification nach `harness/review.md` und `harness/verification.md` |
| `docs/plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/` | update | Closure (§7): Evidence-Datei dieses Slice und `state.md` |
| `internal/**`, `spec/architecture.md`, `README.md`, `README.de.md`, `CHANGELOG.md` | unverändert, außer `CHANGELOG.md` nach DoD-Punkt 7 | der CR ändert kein Produktverhalten; `CHANGELOG.md` Zeile 16 nennt den Test (gelesen), dessen Nachzug gehört dem Test-Slice |

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

- **Der neue Wortlaut sagt mehr zu, als Messung und Suite belegen** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)): etwa „zwei Instanzen“, die nur die Messung belegt, oder „Isolation“ als Oberbegriff, wo nur Workspace-Volumes getestet sind. — **Ausgang:** offen bis zur Freigabe; vor ihr wird jedes Wort des Vorschlags gegen Messprotokoll und Test gelesen, die Beleg-Tabelle (DoD 4) ist dieser Abgleich.
- **Die Suite belegt „zwei“ erst nach einem grünen Lauf** von `TestDevcontainerInstance_Docker_TwoInstancesAreIsolated` auf dem Runner (der Test ist nicht gelaufen). — **Ausgang:** offen bis zum Lauf `integration-docker` auf `main`; der Wortlaut wird vorher nicht geschrieben.
- **Ein Satz im Handbuch behauptet Isolation von Docker-Speicher und Modellen** (`docs/user/benutzerhandbuch.md` Zeilen 903 bis 906) und ist von der Suite nicht pinnbar. — **Ausgang:** offen bis zur Lesung; gepinnt (Text-Pin der Mounts) oder gestrichen.
- **Die Version `0.26.0` ist vergeben**, wenn der Slice startet: andere Entwürfe in `open/` mit Lastenheft-CR können sie früher nehmen. — **Ausgang:** offen bis zum Start; Version wird dann aus Zeile 9 und der Historie gelesen.
- **Test-Lücken der Eingrenzung von `shell` auf die Engine 29.9.2** (Review F-7): (a) der Fall „die Sitzung schreibt die Zeile selbst“ ist von keinem Test gepinnt (Lesung des Codes); (b) die Form der Engine 28.0.4 pinnen nur Einheitstests mit Fake, kein Lauf gegen eine solche Engine; (c) ob der GitHub-Runner die Engine 29.9.2 stellt, ist nicht geprüft. — **Ausgang:** weiter offen (Anker: Review-Report F-7 und `docs/maintainer/zusagen-und-messstaende.md` Abschnitt 3.15).
- **Das Lastenheft ist Vertrag und wird enger, nicht weiter:** eine Rücknahme einer Zusage ist eine Vertragsänderung, auch ohne neues Verhalten. — **Ausgang:** entfallen: der Slice ist als CR mit Freigabe des Projektinhabers geschnitten.
- **Eine Accepted-ADR trägt den Satz:** [`ADR-0038`](../../adr/0038-devcontainer-instanz-zweites-id-label.md) wurde nicht auf Zahl und Grenze der Zusage gelesen. — **Ausgang:** offen; der Implementer liest ihre Konsequenzen; ein Satz, der überholt wird, löst eine Folge-ADR aus (`AGENTS.md` §3.5), die ADR wird nicht umgeschrieben.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren) · `grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln.

*Wird bei der Closure gefüllt.* Hier steht vor dem Lastenheft-Commit der Vermerk der Freigabe (Datum, Variante, Wortlaut) und der Entscheidungen zu den Fragen 1 bis 4 aus §1.

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

- **Was hat funktioniert:**
- **Was ging anders als geplant:**
- **Gegenstand:** *(nur falls Variante 2 gewählt wird: entfallen: Der Wortlaut bleibt, die Messung ist der Beleg; sonst Zeile löschen)*
- **Steering-Loop-Eintrag:**
- **Beobachtungs-Register (`../observations/`):**
- **Folge-Slices:**
- **Risiken aus §6:**
- **Drei Paarungen:**

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
