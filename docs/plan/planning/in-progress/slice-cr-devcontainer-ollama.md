# Slice CR: Katalog-Feature `ollama` im Sandbox-Devcontainer — Zusage nur für Gemessenes (`LH-FA-DEV-003`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) — sechster Slice der Welle; ihm geht [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md) voraus (Messwerte und Folge-ADR), ihm folgt [`slice-v1-devcontainer-feature-ollama`](../open/slice-v1-devcontainer-feature-ollama.md) (Umsetzung).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Devcontainer-Features; Erweiterung), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil; nur bei Netzdienst-Ausweisung), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts: unverändert). Verwandt: [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) und die Folge-ADR des Mess-Slice.
Anlass: Projektinhaber, 2026-10-18 („Lastenheft-CR: Katalog erweitert um `ollama`: Zusage nur für Gemessenes“).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), gegebenenfalls [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)), §16 (Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt der CR nicht; die Spezifikation folgt im letzten Slice der Welle ([`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md)).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Zusage „Ollama läuft als Katalog-Feature `ollama` im Devcontainer, CPU-only, Modelle im Container“ steht im Lastenheft (Änderung auf die **nächste freie Minor-Version**; `0.15.0`, wenn der Anzeige-CR `0.14.0` trägt; ein Neu-Laden-CR entsteht nach R0 nicht; der Enum-CR [`slice-cr-cli-json-command-enum`](../open/slice-cr-cli-json-command-enum.md) folgt getrennt mit `0.16.0`, die Begründung steht dort §1), **mit den Messwerten des Mess-Slice** und in der vom Projektinhaber freigegebenen Fassung. Der Code-Slice beginnt erst danach.

**Warum eine Lastenheft-Änderung:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) führt Beispiele auf („Git, Docker CLI, Node.js, Java, SDKMAN, PostgreSQL Client, Kubernetes Tools“), keine abschließende Liste und keine Katalogschlüssel; `python`, `homebrew` und `sdkman` kamen ohne Lastenheft-Änderung (gelesen: [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md) §1, Absatz „Katalog-Schlüssel allein sind keine Zusage“). Für `ollama` verlangt der Projektinhaber ausdrücklich eine Zusage „nur für Gemessenes“; das ist eine Aussage über Produktverhalten (Installation, CPU-only, Modellspeicher) und gehört deshalb ins Lastenheft.

### Gefundener Stand (gelesen am 2026-10-18)

- [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features): Absatz „Beispiele:“ mit sieben Spiegelstrichen; danach „Für optionale externe Feature-Quellen gilt:“ mit den Regeln zur Freigabe ([`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)). Ein eingebauter Katalogeintrag braucht keine Freigabe (`featureCatalogue`, `internal/hexagon/application/devcontainer_features.go`: „Built-in entries are aktivierbar without an Allowlist entry“).
- Die Egress-Hosts der Features stehen nicht im Lastenheft, sondern in [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion); [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) sagt nur, die Default-Allowlist sei „dokumentiert“. Eine Lastenheft-Aussage zu den Ollama-Hosts ist deshalb **nicht** nötig.
- Das Lastenheft nennt heute weder `sshd` noch `claude-code` als Katalogschlüssel (Suche `grep -i "claude-code\|sshd" spec/lastenheft.md` am 2026-10-18: kein Treffer; das Lastenheft spricht vom „Login-Dienst“ und von „Zugangsdaten eines Agenten“, nicht von Schlüsseln des Katalogs); der CR führt `ollama` als **einzigen** namentlich genannten Katalogschlüssel ein, weil seine Zusage an Messbedingungen hängt.

### Wortlaut (alt/neu; vom Projektinhaber freigegeben und im Lastenheft übernommen)

Die Messwerte des Mess-Slice sind eingetragen; eine Aussage ohne Messwert steht nicht im Text. Das Lastenheft nennt keine Slice- oder Wellen-Kennung; im Lastenheft stehen die Verweise als Anker `#…` der Datei.

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

Neu (Variante, nur bei Bedarf):

> - einen im Container aktivierten Login-Dienst (Passwort-Login) und einen im Container aktivierten Dienst, der von außerhalb des Containers erreichbar ist, in der Befehlsausgabe ausweisen.

Eine Änderung dieser Stelle zieht [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Ausnahme 1 nennt „eines im Container aktivierten Login-Dienstes“) mit; der CR nimmt sie dann als Stelle 2a auf, mit einem eigenen alt/neu.

**Keine Änderung** (gelesen, bewusst nicht angefasst): [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) — gewählt ist ein Binär-Archiv, kein Drittanbieter-Feature; die Auslegung (SHA-256-Pin ohne Signatur, Grenze ausgewiesen) trägt [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 3 (Muster [`ADR-0023`](../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md)), die Warnung ist eine eigene Meldung ohne `SecurityRelaxation` und entfällt unter `--quiet`, ohne Lastenheft-Änderung.

**Stelle 3 — Versionsfeld und Historie (§16).**

Alt: `| Version          | 0.13.0 |` (bzw. der Stand nach dem Landen der Vorgänger-CR in dieser Welle)

Neu: `| Version          | 0.15.0 |` (der Anzeige-CR trägt `0.14.0`). Der Enum-CR [`slice-cr-cli-json-command-enum`](../open/slice-cr-cli-json-command-enum.md) landet danach mit `0.16.0` und wird mit diesem CR **nicht** im selben Lastenheft-Commit gebündelt (Freigabestand, Review und Rücknahme verschieden, Begründung dort §1).

Neue Historie-Zeile:

> | 0.15.0 | 2026-10-18 | [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features): Katalog-Feature `ollama` (Ollama im Container, nur CPU, Modelle im benannten Volume) für die gemessenen Hosts und Versionen; nicht zugesagt: GPU, Anschluss an einen Ollama des Hosts, Modellverwaltung. | Vereinbarung mit dem Projektinhaber |

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Messung, Quellenwahl und Folge-ADR** — [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md); der CR trägt nur, was dort gemessen und entschieden ist.
- **Code, Tests, Spezifikation, Handbuch, CHANGELOG** — [`slice-v1-devcontainer-feature-ollama`](../open/slice-v1-devcontainer-feature-ollama.md) und [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md); das Lastenheft ist Zielbild, der Code folgt.
- **Eine Zusage zu Hosts, Version oder Modellen** — Hosts und Version gehören in die Spezifikation ([`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)), nicht ins Lastenheft; Modelle wählt der Nutzer.
- **Anschluss an den Ollama des Hosts, GPU, weitere Agenten** — vom Projektinhaber ausgeschlossen bzw. nicht genannt.

## 2. Definition of Done

- [x] **Lastenheft 0.15.0:** Stellen 1 und 3 in der vom Projektinhaber freigegebenen Fassung (Stellen 2/2a entfallen: Loopback), Messwerte eingetragen, Versionsfeld, Historie-Zeile mit Datum (Verweis „Vereinbarung mit dem Projektinhaber“); die Freigabe ist im Gespräch belegt und in §7 vermerkt (Datei-Beleg besteht nicht).
- [x] **Vorab-Suche und Nachzug:** Die Suche gegen ältere Lastenheft-Zusagen (Features, Katalog, Agenten, Sandbox-Regeln, Login-Dienst, `--quiet`) ist in §7 festgehalten; mitzuändernde Stellen: keine über die freigegebenen hinaus.
- [ ] `make gates` grün; `make doc-trace` ohne neue Waise.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8); der Review läuft über die **freigegebene** Fassung (Lehre der Vorgänger-Welle).
- [ ] Doku-Update: kein Handbuch, kein CHANGELOG (kein Produktverhalten; die Code-Slices ziehen nach).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen von der Welle-Closure geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Stellen 1 und 3 (§3 [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features); Versionsfeld; §16 Historie); [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) unverändert (Stellen 2/2a entfallen); gelesen: Zeilen 394 bis 417, 438 bis 477, 1058 bis 1071, Historie |
| `docs/plan/planning/in-progress/slice-cr-devcontainer-ollama.md` (diese Datei; vorher `open/`) | update | Freigabe, Vorab-Suche, Zustandsmarker auf Ist-Stand; Closure (§7) folgt durch den Planer |
| `docs/plan/planning/in-progress/roadmap.md` | update | Arbeitszeile ersetzt den Ruhe-Marker; Verweis-Nachzug des Moves |
| `docs/plan/planning/done/slice-v1-cli-egress-anzeige.md`, `done/slice-v1-devcontainer-ollama-messung.md`, `done/slice-v1-sandbox-egress-neuladen-entwurf.md`, `open/README.md`, `open/slice-cr-cli-json-command-enum.md`, `open/slice-v1-devcontainer-feature-ollama.md`, `open/slice-v1-devcontainer-start-ohne-cli.md`, `open/slice-v2-spec-sandbox-egress-ollama.md`, `welle-sandbox-egress-ollama.md`, `docs/reviews/2026-10-18-slice-v1-sandbox-egress-neuladen-entwurf.md`, `docs/reviews/2026-10-18-verification-slice-v1-sandbox-egress-neuladen-entwurf.md` | update | reiner Verweis-Nachzug des Moves `open/` → `in-progress/` (`make slice-mv`), kein Inhalt |
| `docs/reviews/` (Review-Report, Verification Evidence) | neu | Review und Verifikation nach `harness/review.md` und `harness/verification.md` |
| `CHANGELOG.md` | unverändert | der CR ändert kein Produktverhalten; ein Eintrag entsteht im Code-Slice |

**Übergabe aus dem Entwurf [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (`Accepted` am 2026-10-18, Freigabe-Punkte a bis e wie vorgeschlagen; Bedarf, kein Wortlaut):** Stelle 1 — Quelle und Prüfung: gepinntes Release-Archiv des Upstreams, nur per SHA-256 geprüft, **ohne Signatur** (die Grenze ist auszuweisen); Bindung: nur Loopback (gemessen); Start: nicht von selbst; Modelle im benannten Volume überstehen das **Löschen des Containers**, der Satz „überstehen den Neubau“ entfällt (nicht gemessen); wirksam nur mit Sandbox-Profil; Lauf auf `amd64` nicht verifiziert. Stelle 2 und 2a **entfallen** (Loopback, kein Netzdienst nach außen; Entscheidung 5 der ADR). [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) bleibt unverändert (die Auslegung steht in der ADR); die Hosts stehen nicht im Lastenheft.

Vorab-Suche (Pflicht vor jeder Vertragsänderung, Muster [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md)): `grep -i` über `spec/lastenheft.md` nach `Feature`, `Katalog`, `Login-Dienst` und `quiet`; der Implementer führt sie vor dem Lastenheft-Commit aus und trägt die Treffer hier ein.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md) liegt in `done/` (Messwerte und ADR-Status liegen vor); der Projektinhaber hat den Wortlaut freigegeben (Gespräch).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): der Wortlaut verlangt mehr als zwei Anforderungen oder eine Sicherheitsentscheidung, die der Mess-Slice nicht getroffen hat.
- `in-progress` → `open` (blockiert): keine Freigabe des Wortlauts, oder die ADR des Mess-Slice ist noch `Proposed`; das Lastenheft bleibt unverändert, kein Code-Slice startet.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). Der Lastenheft-Commit trägt die Freigabe nur, wenn §7 sie vermerkt.

## 6. Risiken und offene Punkte

- **Ein Satz im Lastenheft stützt sich auf einen Messwert, den der Mess-Slice nicht liefert.** — **Ausgang:** offen bis zum Start; der Satz entfällt (keine Aussage ohne Messwert).
- **Die Zusage „Modelle überstehen den Neubau“ ist ungemessen.** — **Ausgang:** offen bis O5 des Mess-Slice; der Satz bleibt nur mit Messwert.
- **Die Reihenfolge der Lastenheft-Versionen mit dem Anzeige-CR und dem Enum-CR kollidiert.** — **Ausgang:** offen; die Version ist die nächste freie zum Zeitpunkt des Commits, nicht vorab festgeschrieben; erwartet `0.15.0` hier und `0.16.0` im Enum-CR ([`slice-cr-cli-json-command-enum`](../open/slice-cr-cli-json-command-enum.md), getrennt entschieden).
- **Ein Review läuft vor der Freigabe** (siehe DoD). — **Ausgang:** offen; Review über die freigegebene Fassung.

## 7. Closure-Notiz

*(Closure durch den Planer; Steering-Loop-Eintrag, Register, Folge-Slices und Risiken folgen dort.)*

**Freigabe des Wortlauts.** Im Gespräch am 2026-10-18 nach wörtlicher Vorlage der Texte (Stelle 1: Spiegelstrich und Absatz; Stelle 3: Versionsfeld und Historie-Zeile) durch den Projektinhaber freigegeben („ok“). Übernommen unverändert; einzige Anpassung: die Link-Form im Lastenheft (Anker `#lh-fa-dev-008--egress-restriktion` der eigenen Datei). Stellen 2 und 2a entfallen ([`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md)).

**Vorab-Suche (2026-10-18, vor dem Lastenheft-Commit).** `grep -i` über `spec/lastenheft.md` nach `feature`, `katalog`, `login-dienst`, `quiet`, `agent`, `0.14`/`0.15`; dazu gelesen: [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), Glossar. Berührte Stellen, jeweils nicht geändert:

- [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Zeilen 444 bis 456): Ausweisungspflicht nennt Lockerungen und Login-Dienst; Ollama lauscht nur auf Loopback, also kein Anlass; die Regel „Features gelten nach [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)“ trägt das Feature.
- [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Zeilen 473 bis 477): Default-Allowlist „dokumentiert“, Herkunft „Feature“ ist vorhanden; die Hosts stehen in der Spezifikation (Folge-Slice).
- [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Zeilen 1059 bis 1071): erfasst Features und Skripte aus nicht freigegebenen Quellen; ein eingebautes Katalog-Feature ist keine externe Quelle, das nicht signierte Archiv trägt [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md).
- [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Zeilen 184 bis 198): die Warnung zum nicht signierten Archiv ist eine Hinweis-Warnung und entfällt unter `--quiet` (Ausnahme 1 betrifft Lockerungen und Login-Dienst); keine Änderung, Auslegung in der ADR.
- [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) (Befehlsübersicht, Zeile 1108 ff.), Glossar (Zeile 1526, Beispiele „Node.js, Docker CLI“, offene Beispiele), Modulliste (Devcontainer-Generator, Zeilen 132 bis 140) und Nicht-Kernumfang (Zeile 1373, „Kubernetes-Produktionsdeployment“): zählen keine Katalog-Features abschließend auf; die Suche nach den Beispielnamen der Liste (Node.js, Docker CLI, SDKMAN, Kubernetes) fand ihre Aufzählung nur in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und im Glossar.
- [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zeile 482 bis 486, „bei aktivem Agenten-Feature“): betrifft Agenten, nicht Ollama.

**Nicht gelaufen:** `make test-docker` und `make verify-depguard` (kein Code, keine Importregel berührt).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird allein `spec/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); die Schwelle (zwei von drei Achsen) ist laut Modus-Deklaration erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; kein Treffer für Lastenheft-Änderungen unter den offenen Beobachtungen. [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) gelten (verkörpert in `AGENTS.md` §3.10): jede Bedingung im Wortlaut („nur auf der Loopback-Adresse“, „überstehen den Neubau“) braucht im Code-Slice einen Test oder entfällt.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
