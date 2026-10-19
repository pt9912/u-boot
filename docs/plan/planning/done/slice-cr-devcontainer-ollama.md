# Slice CR: Katalog-Feature `ollama` im Sandbox-Devcontainer — Zusage nur für Gemessenes (`LH-FA-DEV-003`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-sandbox-egress-ollama`](welle-sandbox-egress-ollama.md) — sechster Slice der Welle; ihm geht [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md) voraus (Messwerte und Folge-ADR), ihm folgt [`slice-v1-devcontainer-feature-ollama`](../done/slice-v1-devcontainer-feature-ollama.md) (Umsetzung).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Devcontainer-Features; Erweiterung), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil; unverändert: Stelle 2 entfallen, Ollama lauscht nur auf Loopback), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts: unverändert). Verwandt: [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) und die Folge-ADR des Mess-Slice.
Anlass: Projektinhaber, 2026-10-18 („Lastenheft-CR: Katalog erweitert um `ollama`: Zusage nur für Gemessenes“).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features); [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) bleibt unverändert), §16 (Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt der CR nicht; die Spezifikation folgt im letzten Slice der Welle ([`slice-v2-spec-sandbox-egress-ollama`](../done/slice-v2-spec-sandbox-egress-ollama.md)).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Zusage „Ollama läuft als Katalog-Feature `ollama` im Devcontainer, CPU-only, Modelle im Container“ steht im Lastenheft (Änderung auf die **nächste freie Minor-Version**: `0.15.0`, der Anzeige-CR trägt `0.14.0`; ein Neu-Laden-CR entsteht nach R0 nicht; der Enum-CR [`slice-cr-cli-json-command-enum`](../done/slice-cr-cli-json-command-enum.md) folgt getrennt mit `0.16.0`, die Begründung steht dort §1), **mit den Messwerten des Mess-Slice** und in der vom Projektinhaber freigegebenen Fassung. Gemessen sind Installation und Lauf des Archivs, die Bindung an Loopback, das Überleben des Volumes beim Löschen des Containers und die Hosts des Pulls; Verhalten des Generators (Ausweisung beim Erzeugen, Eintrag in die Default-Allowlist, Wirkung im Sandbox-Profil) ist Zielbild, das der Code-Slice einlöst, und nicht gemessen ([`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md), Abschnitt „Nicht gemessen“). Der Code-Slice beginnt erst danach.

**Warum eine Lastenheft-Änderung:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) führt Beispiele auf („Git, Docker CLI, Node.js, Java, SDKMAN, PostgreSQL Client, Kubernetes Tools“), keine abschließende Liste und keine Katalogschlüssel; `python`, `homebrew` und `sdkman` kamen ohne Lastenheft-Änderung (gelesen: [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md) §1, Absatz „Katalog-Schlüssel allein sind keine Zusage“). Für `ollama` verlangt der Projektinhaber ausdrücklich eine Zusage „nur für Gemessenes“; das ist eine Aussage über Produktverhalten (Installation, CPU-only, Modellspeicher) und gehört deshalb ins Lastenheft.

### Gefundener Stand (gelesen am 2026-10-18)

- [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features): Absatz „Beispiele:“ mit sieben Spiegelstrichen; danach „Für optionale externe Feature-Quellen gilt:“ mit den Regeln zur Freigabe ([`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)). Ein eingebauter Katalogeintrag braucht keine Freigabe (`featureCatalogue`, `internal/hexagon/application/devcontainer_features.go`: „Built-in entries are aktivierbar without an Allowlist entry“).
- Die Egress-Hosts der Features stehen nicht im Lastenheft, sondern in [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion); [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) sagt nur, die Default-Allowlist sei „dokumentiert“. Eine Lastenheft-Aussage zu den Ollama-Hosts ist deshalb **nicht** nötig.
- Das Lastenheft nennt heute weder `sshd` noch `claude-code` als Katalogschlüssel (Suche `grep -i "claude-code\|sshd" spec/lastenheft.md` am 2026-10-18: kein Treffer; das Lastenheft spricht vom „Login-Dienst“ und von „Zugangsdaten eines Agenten“, nicht von Schlüsseln des Katalogs); der CR führt `ollama` als **einzigen** namentlich genannten Katalogschlüssel ein, weil seine Zusage an Messbedingungen hängt.

### Wortlaut (alt/neu; vom Projektinhaber freigegeben und im Lastenheft übernommen)

Die Messwerte des Mess-Slice sind eingetragen. Aussagen über das Verhalten des Generators tragen keinen Messwert; sie sind Zielbild (siehe §1 Ziel). Das Lastenheft nennt keine Slice- oder Wellen-Kennung; im Lastenheft stehen die Verweise als Anker `#…` der Datei.

**Stelle 1 — [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Beispiele und neuer Absatz.**

Alt (Spiegelstriche, bleiben):

> - Git
> - Docker CLI
> - Node.js
> - Java
> - SDKMAN
> - PostgreSQL Client
> - Kubernetes Tools

Neu (zusätzlicher Spiegelstrich am Ende der Liste):

> - Ollama (lokale Sprachmodelle, nur CPU)

Neu (zusätzlicher Absatz zwischen der Liste und „Für optionale externe Feature-Quellen gilt:“):

> Ollama ist als eingebautes Feature `ollama` wählbar (`u-boot config set devcontainer.features.ollama.enabled true`) und wirkt im Sandbox-Profil. Es installiert Ollama im Container ohne Grafikbeschleunigung (nur CPU); die Installation stammt aus einem gepinnten Release-Archiv mit geprüfter SHA-256-Prüfsumme. Das Archiv ist nicht signiert; `u-boot` weist das beim Erzeugen aus. Der Server startet nicht von selbst und lauscht nur auf der Loopback-Adresse des Containers. Modelle lädt der Nutzer im Container (`ollama pull`); sie liegen in einem benannten Volume und überstehen das Löschen des Containers. Mit der Egress-Restriktion ([`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) gehören die Hosts des Modell-Pull zur Default-Allowlist. Gemessen ist das auf Docker 29.5.2, Colima 0.10.3, arm64, Ollama 0.40.6; nicht zugesagt sind GPU-Beschleunigung, ein Anschluss an einen Ollama-Dienst des Hosts, andere Architekturen und Container-Engines als gemessen, das Überstehen des Neubaus des Images sowie Auswahl oder Verwaltung von Modellen durch `u-boot`.

**Stelle 2 — [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), Ausweisung eines Netzdienstes: entfallen** (Ollama lauscht nur auf Loopback, kein Netzdienst nach außen; [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 5). Die folgende alt/neu-Fassung ist nicht übernommen.

Alt (letzter Spiegelstrich des Abschnitts „Das erzeugte Ergebnis im Sandbox-Profil muss:“):

> - einen im Container aktivierten Login-Dienst (Passwort-Login) in der Befehlsausgabe ausweisen.

Neu (nicht übernommen):

> - einen im Container aktivierten Login-Dienst (Passwort-Login) und einen im Container aktivierten Dienst, der von außerhalb des Containers erreichbar ist, in der Befehlsausgabe ausweisen.

Stelle 2a ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Ausnahme 1) entfällt mit Stelle 2.

**Keine Änderung** (gelesen, bewusst nicht angefasst): [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) — gewählt ist ein Binär-Archiv, kein Drittanbieter-Feature; die Auslegung (SHA-256-Pin ohne Signatur, Grenze ausgewiesen) trägt [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 3 (Muster [`ADR-0023`](../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md)), die Warnung ist eine eigene Meldung ohne `SecurityRelaxation` und entfällt unter `--quiet`, ohne Lastenheft-Änderung.

**Stelle 3 — Versionsfeld und Historie (§16).**

Alt: `| Version          | 0.14.0 |` (Stand nach dem Anzeige-CR)

Neu: `| Version          | 0.15.0 |` (der Anzeige-CR trägt `0.14.0`). Der Enum-CR [`slice-cr-cli-json-command-enum`](../done/slice-cr-cli-json-command-enum.md) landet danach mit `0.16.0` und wird mit diesem CR **nicht** im selben Lastenheft-Commit gebündelt (Freigabestand, Review und Rücknahme verschieden, Begründung dort §1).

Neue Historie-Zeile:

> | 0.15.0 | 2026-10-18 | [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features): Katalog-Feature `ollama` (Ollama im Container, nur CPU, Modelle im benannten Volume) für die gemessenen Hosts und Versionen; nicht zugesagt: GPU, Anschluss an einen Ollama des Hosts, Modellverwaltung. | Vereinbarung mit dem Projektinhaber |

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Messung, Quellenwahl und Folge-ADR** — [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md); der CR trägt nur, was dort gemessen und entschieden ist.
- **Code, Tests, Spezifikation, Handbuch, CHANGELOG** — [`slice-v1-devcontainer-feature-ollama`](../done/slice-v1-devcontainer-feature-ollama.md) und [`slice-v2-spec-sandbox-egress-ollama`](../done/slice-v2-spec-sandbox-egress-ollama.md); das Lastenheft ist Zielbild, der Code folgt.
- **Eine Zusage zu Hosts, Version oder Modellen** — Hosts und Version gehören in die Spezifikation ([`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)), nicht ins Lastenheft; Modelle wählt der Nutzer.
- **Anschluss an den Ollama des Hosts, GPU, weitere Agenten** — vom Projektinhaber ausgeschlossen bzw. nicht genannt.

## 2. Definition of Done

- [x] **Lastenheft 0.15.0** (Commit `d531479`): Stellen 1 und 3 in der vom Projektinhaber freigegebenen Fassung (Stellen 2/2a entfallen: Loopback), Messwerte eingetragen, Versionsfeld, Historie-Zeile mit Datum (Verweis „Vereinbarung mit dem Projektinhaber“); die Freigabe ist im Gespräch belegt und in §7 vermerkt (Datei-Beleg besteht nicht).
- [x] **Vorab-Suche und Nachzug:** Die Suche gegen ältere Lastenheft-Zusagen (Features, Katalog, Agenten, Sandbox-Regeln, Login-Dienst, `--quiet`) ist in §7 festgehalten; mitzuändernde Stellen: keine über die freigegebenen hinaus.
- [x] `make gates` grün; `make doc-trace` ohne neue Waise. Läufe und Exit-Codes in §7.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8); der Review lief über die **freigegebene** Fassung (Report-Commit `bc3b313`, Verification-Report-Commit `b9ae4e7`).
- [x] Doku-Update: kein Handbuch, kein CHANGELOG (kein Produktverhalten; die Code-Slices ziehen nach).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen von der Welle-Closure geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Stellen 1 und 3 (§3 [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features); Versionsfeld; §16 Historie); [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) unverändert (Stellen 2/2a entfallen); gelesen: Zeilen 394 bis 417, 438 bis 477, 1058 bis 1071, Historie |
| `docs/plan/planning/done/slice-cr-devcontainer-ollama.md` (diese Datei; vorher `open/`) | update | Freigabe, Vorab-Suche, Zustandsmarker auf Ist-Stand; Closure (§7) folgt durch den Planer |
| `docs/plan/planning/in-progress/roadmap.md` | update | Arbeitszeile ersetzt den Ruhe-Marker; Verweis-Nachzug des Moves |
| `docs/plan/planning/done/slice-v1-cli-egress-anzeige.md`, `done/slice-v1-devcontainer-ollama-messung.md`, `done/slice-v1-sandbox-egress-neuladen-entwurf.md`, `open/README.md`, `done/slice-cr-cli-json-command-enum.md`, `done/slice-v1-devcontainer-feature-ollama.md`, `open/slice-v1-devcontainer-start-ohne-cli.md`, `done/slice-v2-spec-sandbox-egress-ollama.md`, `welle-sandbox-egress-ollama.md`, `docs/reviews/2026-10-18-slice-v1-sandbox-egress-neuladen-entwurf.md`, `docs/reviews/2026-10-18-verification-slice-v1-sandbox-egress-neuladen-entwurf.md` | update | reiner Verweis-Nachzug der Moves `open/` → `in-progress/` und `in-progress/` → `done/` (`make slice-mv`), kein Inhalt |
| `docs/reviews/2026-10-18-slice-cr-devcontainer-ollama.md`, `docs/reviews/2026-10-18-verification-slice-cr-devcontainer-ollama.md` | neu | Review und Verifikation nach `harness/review.md` und `harness/verification.md` |
| `docs/plan/planning/done/slice-v1-devcontainer-feature-ollama.md`, `done/slice-v2-spec-sandbox-egress-ollama.md` | update | Closure: Übergabe der Test-Pins V-3 bis V-5 (Feature-Slice §2 und §3) und der Bedingungen aus V-6 (Spec-Slice §3); gelesen vor der Änderung (beide Dateien, §2 und §3) |
| `docs/plan/planning/open/README.md`, `welle-sandbox-egress-ollama.md` | update | Zustand der Zeile dieses Slice auf „geliefert“ und [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) aus Bezug und Tabelle (V-8); Welle §3: Release-Gate zu F-1 und F-2; Verweis-Nachzug des Moves |
| `docs/plan/planning/observations/` (`BEO-WZG/host-werkzeug-ausserhalb-docker`, `BEO-DOK/plan-aussage-ueber-code`: je `state.md` und eine Evidence-Datei; neu `BEO-DOK/zusage-weiter-als-messung`) | update, neu | Register bei der Closure (§7) |
| `CHANGELOG.md` | unverändert | der CR ändert kein Produktverhalten; ein Eintrag entsteht im Code-Slice |

**Übergabe aus dem Entwurf [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (`Accepted` am 2026-10-18, Freigabe-Punkte a bis e wie vorgeschlagen; Bedarf, kein Wortlaut):** Stelle 1 — Quelle und Prüfung: gepinntes Release-Archiv des Upstreams, nur per SHA-256 geprüft, **ohne Signatur** (die Grenze ist auszuweisen); Bindung: nur Loopback (gemessen); Start: nicht von selbst; Modelle im benannten Volume überstehen das **Löschen des Containers**, der Satz „überstehen den Neubau“ entfällt (nicht gemessen); wirksam nur mit Sandbox-Profil; Lauf auf `amd64` nicht verifiziert. Stelle 2 und 2a **entfallen** (Loopback, kein Netzdienst nach außen; Entscheidung 5 der ADR). [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) bleibt unverändert (die Auslegung steht in der ADR); die Hosts stehen nicht im Lastenheft.

Vorab-Suche (Pflicht vor jeder Vertragsänderung, Muster [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md)): `grep -i` über `spec/lastenheft.md` nach `Feature`, `Katalog`, `Login-Dienst` und `quiet`; ausgeführt vor dem Lastenheft-Commit, Treffer in §7.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md) liegt in `done/` (Messwerte und ADR-Status liegen vor); der Projektinhaber hat den Wortlaut freigegeben (Gespräch).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): der Wortlaut verlangt mehr als zwei Anforderungen oder eine Sicherheitsentscheidung, die der Mess-Slice nicht getroffen hat.
- `in-progress` → `open` (blockiert): keine Freigabe des Wortlauts, oder die ADR des Mess-Slice ist noch `Proposed`; das Lastenheft bleibt unverändert, kein Code-Slice startet.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). Der Lastenheft-Commit trägt die Freigabe nur, wenn §7 sie vermerkt.

## 6. Risiken und offene Punkte

- **Ein Satz im Lastenheft stützt sich auf einen Messwert, den der Mess-Slice nicht liefert.** — **Ausgang:** entfallen für die Messaussagen (Installation, Loopback, Volume, Messstand: von Review F-6 und Verifier gegen [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) bestätigt); die Generator-Aussagen sind Zielbild und nicht gemessen (Review F-3, in §1 berichtigt). Die Hosts des Modell-Pull sind das **weiter offene** Risiko der nächsten Zeile.
- **Die Hosts des Modell-Pull im Lastenheft sind allgemeiner als die Messung** (Review F-1, Verification V-1: zwei Modelle, ein Zeitpunkt, ein Netzstandort; Stabilität des Blob-Hosts nicht gemessen; ein Wechsel lässt den Pull scheitern). — **Ausgang:** weiter offen → **Gate vor dem Release** (Entscheidung des Projektinhabers: Bestätigung des Wortlauts oder Änderung; bei Änderung neue Freigabe, neuer Lastenheft-Commit samt Historie, Review des Diffs). Der Wortlaut bleibt unverändert. Bis dahin tragen die Grenze [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 7 und das Risiko „Hosts des Modell-Pull wechseln“ in [`slice-v1-devcontainer-feature-ollama`](../done/slice-v1-devcontainer-feature-ollama.md) §6.
- **„weist das beim Erzeugen aus“ steht neben Ausweisungen, die unter `--quiet` bleiben** (Review F-2, Verification V-2; [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 3 lässt die Warnung unter `--quiet` entfallen). — **Ausgang:** weiter offen → **Gate vor dem Release** (Entscheidung des Projektinhabers: Bestätigung der Lesart der ADR oder Änderung des Wortlauts). Der Wortlaut bleibt unverändert. Der Test unter `--quiet` im Code-Slice richtet sich nach der Entscheidung (Übergabe V-3, §7).
- **Die Zusage „Modelle überstehen den Neubau“ ist ungemessen.** — **Ausgang:** entfallen: der Satz steht nicht im Lastenheft; dort steht „nicht zugesagt: das Überstehen des Neubaus des Images“, Modelle überstehen laut Messung das Löschen des Containers.
- **Die Reihenfolge der Lastenheft-Versionen mit dem Anzeige-CR und dem Enum-CR kollidiert.** — **Ausgang:** eingetreten wie erwartet, ohne Kollision: dieser CR trägt `0.15.0` (Anzeige-CR `0.14.0`); der Enum-CR [`slice-cr-cli-json-command-enum`](../done/slice-cr-cli-json-command-enum.md) nimmt `0.16.0`.
- **Ein Review läuft vor der Freigabe** (siehe DoD). — **Ausgang:** entfallen: der Review lief über die freigegebene Fassung (Inhalts-Commit `d531479`); nach dem Review gab es keine Fix-Runde am Lastenheft (Verification V-7).

## 7. Closure-Notiz

**Freigabe des Wortlauts.** Im Gespräch am 2026-10-18 nach wörtlicher Vorlage der Texte (Stelle 1: Spiegelstrich und Absatz; Stelle 3: Versionsfeld und Historie-Zeile) durch den Projektinhaber freigegeben („ok“). Übernommen unverändert; einzige Anpassung: die Link-Form im Lastenheft (Anker `#lh-fa-dev-008--egress-restriktion` der eigenen Datei). Stellen 2 und 2a entfallen ([`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md)).

**Vorab-Suche (2026-10-18, vor dem Lastenheft-Commit).** `grep -i` über `spec/lastenheft.md` nach `feature`, `katalog`, `login-dienst`, `quiet`, `agent`, `0.14`/`0.15`; dazu gelesen: [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), Glossar. Berührte Stellen, jeweils nicht geändert:

- [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Zeilen 444 bis 456): Ausweisungspflicht nennt Lockerungen und Login-Dienst; Ollama lauscht nur auf Loopback, also kein Anlass; die Regel „Features gelten nach [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)“ trägt das Feature.
- [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Zeilen 473 bis 477): Default-Allowlist „dokumentiert“, Herkunft „Feature“ ist vorhanden; die Hosts stehen in der Spezifikation (Folge-Slice).
- [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Zeilen 1059 bis 1071): erfasst Features und Skripte aus nicht freigegebenen Quellen; ein eingebautes Katalog-Feature ist keine externe Quelle, das nicht signierte Archiv trägt [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md).
- [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Zeilen 184 bis 198): die Warnung zum nicht signierten Archiv ist eine Hinweis-Warnung und entfällt unter `--quiet` (Ausnahme 1 betrifft Lockerungen und Login-Dienst); keine Änderung, Auslegung in der ADR.
- [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) (Befehlsübersicht, Zeile 1108 ff.), Glossar (Zeile 1526, Beispiele „Node.js, Docker CLI“, offene Beispiele), Modulliste (Devcontainer-Generator, Zeilen 132 bis 140) und Nicht-Kernumfang (Zeile 1373, „Kubernetes-Produktionsdeployment“): zählen keine Katalog-Features abschließend auf; die Suche nach den Beispielnamen der Liste (Node.js, Docker CLI, SDKMAN, Kubernetes) fand ihre Aufzählung nur in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und im Glossar.
- [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zeile 482 bis 486, „bei aktivem Agenten-Feature“): betrifft Agenten, nicht Ollama.

**Nicht gelaufen:** `make test-docker` und `make verify-depguard` (kein Code, keine Importregel berührt).

**Belege:** Review-Report [`2026-10-18-slice-cr-devcontainer-ollama`](../../../reviews/2026-10-18-slice-cr-devcontainer-ollama.md) (Commit `bc3b313`; 0 HIGH, 2 MEDIUM, 2 LOW, 4 INFO) und Verification-Report [`2026-10-18-verification-slice-cr-devcontainer-ollama`](../../../reviews/2026-10-18-verification-slice-cr-devcontainer-ollama.md) (Commit `b9ae4e7`; GO mit Auflagen). Der Verifier hat den Lastenheft-Diff (drei Hunks, 5 Einfügungen, 1 Löschung) Satz für Satz mit dem freigegebenen Wortlaut verglichen: wortgleich bis auf die Link-Form. Sensoren des Verifiers und der Closure: unten.

**Entscheidung des Projektinhabers vor dem Release (Gates).** F-1 und F-2 betreffen freigegebenen Text; der Wortlaut bleibt unverändert, und die Entscheidung ist ein **Gate vor dem Release**, nicht vor `done/` (Muster [`slice-cr-sandbox-egress-anzeige`](../done/slice-cr-sandbox-egress-anzeige.md) §7; Anordnung des Projektinhabers vom 2026-10-18, bis zum Release ohne Rückfragen weiterzuarbeiten). Das weicht vom Urteil des Verifiers ab („nicht schließbar bis zur Closure-Notiz“) nur insoweit, als die Closure-Notiz die Gates ausweist; sie entscheidet sie nicht.

- **Gate F-1 (Hosts des Modell-Pull):** Bestätigung des Wortlauts „gehören die Hosts des Modell-Pull zur Default-Allowlist“ oder Änderung. Bei Änderung: neue Freigabe, neuer Lastenheft-Commit samt Historie-Zeile, Review des Diffs (nicht der ganzen Datei). Zustand im Lastenheft: unverändert (0.15.0, Historie-Zeile ohne Bestätigungs-Klausel). **Vermerk 2026-10-19: entschieden** (Delegation durch den Projektinhaber): Wortlaut bleibt unverändert bestätigt; die Aussage ist durch [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 7 begrenzt, ein breiterer Host-Eintrag entsteht nicht ohne Messung (Welle §3, drittes Release-Gate).
- **Gate F-2 („weist aus“ gegen `--quiet`):** Bestätigung der Lesart von [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 3 (die Warnung zum nicht signierten Archiv entfällt unter `--quiet`) oder Änderung des Wortlauts. Zustand im Lastenheft: unverändert. **Vermerk 2026-10-19: entschieden** (Delegation durch den Projektinhaber): Wortlaut bleibt unverändert; die Lesart der ADR gilt (Warnung und Meldung entfallen unter `--quiet`, keine `SecurityRelaxation`), von den Tests des Code-Slice gepinnt (Welle §3, viertes Release-Gate).
- Beide Gates stehen in den Closure-Kriterien der Welle [`welle-sandbox-egress-ollama`](welle-sandbox-egress-ollama.md) §3 (drittes und viertes Release-Gate); der Release-Cut-Slice der Welle nimmt sie in seine Vorbereitung auf. Wer sie bestätigt, trägt den Vermerk mit Datum in diese Zeile und in die Welle-Datei ein.

**Berichtigung der Zustandsmarker (Review F-3, F-4; Verification V-8).** In dieser Datei auf den Ist-Stand gesetzt: Stelle 3 „Alt“ (`0.14.0`), Stelle 2 und 2a als nicht übernommen, Kopf `Bezug` und `Berührte Spec-Stellen` ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) unverändert), Vorab-Suche (ausgeführt, §7), Ziel (`0.15.0`) und die Aussage über Messwerte: Das Lastenheft ist Zielbild; was der Generator tut (Ausweisung, Allowlist, Wirkung im Profil), ist nicht gemessen ([`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md), „Nicht gemessen“). In `open/README.md` und der Tabelle der Welle-Datei steht die Zeile dieses Slice als geliefert und ohne [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) im Bezug.

**Übergabe der Test-Pins (Verification V-3 bis V-5) an [`slice-v1-devcontainer-feature-ollama`](../done/slice-v1-devcontainer-feature-ollama.md)** (verankert in dessen §2 und §3, gelesen vor dem Eintrag):

- **V-3:** Test, dass die Warnung zur Lieferkette und die Meldung „keine Wirkung ohne Sandbox-Profil“ unter `--quiet` entfallen; er folgt der Entscheidung zu Gate F-2 und wird vor dem Test festgeschrieben.
- **V-4:** Golden Case „kein `OLLAMA_HOST` in `devcontainer.json`, Dockerfile und Skript, kein `forwardPorts`“ und Docker-Test „Aufruf über die Container-Adresse scheitert, `127.0.0.1:11434` antwortet“.
- **V-5:** Negativfall „ohne Sandbox-Profil: kein Block im Dockerfile und die ungekennzeichnete Meldung“.

**Übergabe an [`slice-cr-cli-json-command-enum`](../done/slice-cr-cli-json-command-enum.md):** keine Übergabe nötig. Die Verification-Auflage 2 nennt für den Enum-CR keinen Pin; die Abstimmung der Version (`0.16.0`) steht in dessen §1.

**Übergabe an [`slice-v2-spec-sandbox-egress-ollama`](../done/slice-v2-spec-sandbox-egress-ollama.md) (V-6):** Bindung „bei ungesetztem `OLLAMA_HOST`“, `amd64` als nicht verifiziert ([`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 10) und die Grenze des Volumes (Löschen des Containers, nicht Neubau des Images) sind in §1 des Spec-Slice als Bedarf eingetragen.

**Was hat funktioniert:** Der Review las die Stellen unabhängig und fand die zwei Wortlaut-Grenzen (F-1, F-2), die die Vorab-Suche nicht suchen konnte; der Verifier trennte die Frage nach dem Wortlaut (wortgleich) von der Frage nach seiner Reichweite und führte beide MEDIUM-Befunde als Entscheidung des Projektinhabers, statt den freigegebenen Text zu ändern. Es gab keine Fix-Runde am Lastenheft nach dem Review (V-7), also keinen ungeprüften Nachtrag wie beim Anzeige-CR.

**Was ging anders als geplant:** Der Plan sagte „eine Aussage ohne Messwert steht nicht im Text“; der Absatz trägt Generator-Zusagen ohne Messwert (F-3). Die Zustandsmarker blieben nach der Freigabe stehen (F-4, V-8). Die Grenzen der Messung (Modelle, Zeitpunkt, Netzstandort) standen in der ADR und nicht im Wortlaut (F-1).

**Steering-Loop-Eintrag:** gezählt, nicht verkörpert; mit diesem Slice wurde keine Regel und kein Sensor verkörpert. Lerneintrag als **geschärfte Regel für den Wortlaut-Vorschlag** (Kandidat, erstes Auftreten): Vor der Freigabe eines Wortlauts, der „nur Gemessenes“ zusagt, wird die Liste „nicht gemessen“ und die Messbedingungen der Entscheidung (Modelle, Zeitpunkt, Netzstandort, Geltungsbereich eines Begriffs wie „ausweisen“) Satz für Satz gegen den Wortlaut gelesen; jede Grenze steht im Text oder wird ausdrücklich als Entscheidung des Projektinhabers vorgelegt — Register [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), 1×. Der Rest läuft im Register weiter (unten).

**Beobachtungs-Register (`../observations/`):** Neu: [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (1×; F-1, F-2). [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md): `evidence/slice-cr-devcontainer-ollama.md` ergänzt (zwanzigster Vorgang: F-3 und F-4, V-8; Rückfall nach der Verkörperung in `AGENTS.md` §3.10), `state.md` nachgeführt, eine weitere Schärfung nicht entschieden. [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md): `evidence/slice-cr-devcontainer-ollama.md` ergänzt (sechster Vorgang: Review F-8, Verifier V-10 und der Closure-Lauf mit lesenden und einem schreibenden Host-Aufruf, ausgewiesen); die Entscheidung über Verkörperung liegt weiter beim Projektinhaber. [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md): kein neuer Beleg; die Bedingungen des Wortlauts („nur Loopback“, „wirkt im Sandbox-Profil“) sind mit V-4 und V-5 an Tests des Code-Slice gebunden. Kein Eintrag erreicht mit diesem Vorgang erstmals die Schwelle von drei.

**Folge-Slices:** [`slice-v1-devcontainer-feature-ollama`](../done/slice-v1-devcontainer-feature-ollama.md) (Code, Pins V-3 bis V-5), [`slice-cr-cli-json-command-enum`](../done/slice-cr-cli-json-command-enum.md) und [`slice-v2-spec-sandbox-egress-ollama`](../done/slice-v2-spec-sandbox-egress-ollama.md) — Dateien in `open/`; neue Folge-Slices entstehen nicht.

**Risiken aus §6:** jedes trägt in §6 genau einen Ausgang: entfallen (Messwerte der Messaussagen, Neubau, Review vor Freigabe), eingetreten wie erwartet (Version), weiter offen (Hosts des Modell-Pull, „weist aus“ gegen `--quiet`: beide Gate vor dem Release).

**Grenzen:** (1) Die Freigabe des Wortlauts ist nur im Gespräch belegt (Verification-Report Abschnitt 8). (2) Die Wortgleichheit ist von Hand geprüft, nicht mit einem Zeichen-Diff-Werkzeug. (3) Rohdaten des Messprotokolls hat der Review nur zum Teil, der Verifier nicht gelesen; die Messstand-Zahlen stützen sich auf [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md). (4) `make test-docker`, `make verify-depguard` und `make doc-immutable` liefen in der Closure nicht (kein Produktcode, keine ADR im Diff); `make doc-immutable` lief beim Verifier (Exit 0). (5) Der Closure-Lauf nutzte Host-Werkzeuge außerhalb von Docker (siehe Register `BEO-WZG/host-werkzeug-ausserhalb-docker`).

**Drei Paarungen:** Anker: kein Feld `liegt in` in diesem Slice, nichts zu paaren. Folge-Slice: [`slice-v1-devcontainer-feature-ollama`](../done/slice-v1-devcontainer-feature-ollama.md), [`slice-cr-cli-json-command-enum`](../done/slice-cr-cli-json-command-enum.md) und [`slice-v2-spec-sandbox-egress-ollama`](../done/slice-v2-spec-sandbox-egress-ollama.md) liegen in `open/`. Register: die drei genannten Kennungen sind Verzeichnisse im Register und tragen je ein nicht leeres `evidence/`.

**Abgleich der Plan-Tabelle und Sensor-Läufe der Closure:** `git diff --stat 19d8e2a` nach dem letzten Inhalt vor diesem Absatz: 24 Einträge — Lastenheft, dieser Slice (Neuanlage unter `in-progress/` und Löschung unter `open/`), Roadmap, `open/README.md`, vier Slices in `open/` (Feature-Slice und Spec-Slice mit Inhalt, Enum-CR und Start-ohne-CLI als Verweis-Nachzug), Welle-Datei, drei Slices in `done/` (Verweis-Nachzug), Review-Report, Verification-Report, zwei ältere Reports (Verweis-Nachzug) und sieben Dateien des Registers (drei Evidence-Dateien, zwei bestehende `state.md`, ein neues Verzeichnis mit `observation.md` und `state.md`). Jeder Eintrag hat eine Zeile in der Plan-Tabelle §3, jede Zeile mit Datei-Eintrag einen Diff; `CHANGELOG.md` steht nicht im Diff. Läufe (Exit-Code des `make`-Prozesses selbst): `make gates` Exit 0 (Coverage 91,40 % gegen 90 %; `docs-check` 450 Dateien, 0 Befunde) vor und nach den Änderungen der Closure, `make doc-trace` Exit 0 (83 Anforderungen, 0 Waisen). Dieser Absatz ändert nur Text in dieser Datei.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird allein `spec/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); die Schwelle (zwei von drei Achsen) ist laut Modus-Deklaration erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; kein Treffer für Lastenheft-Änderungen unter den offenen Beobachtungen. [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) gelten (verkörpert in `AGENTS.md` §3.10): jede Bedingung im Wortlaut („nur auf der Loopback-Adresse“, „überstehen den Neubau“) braucht im Code-Slice einen Test oder entfällt.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
