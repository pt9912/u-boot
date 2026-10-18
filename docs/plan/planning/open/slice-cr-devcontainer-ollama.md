# Slice CR: Katalog-Feature `ollama` im Sandbox-Devcontainer — Zusage nur für Gemessenes (`LH-FA-DEV-003`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) — sechster Slice der Welle; ihm geht [`slice-v1-devcontainer-ollama-messung`](../in-progress/slice-v1-devcontainer-ollama-messung.md) voraus (Messwerte und Folge-ADR), ihm folgt [`slice-v1-devcontainer-feature-ollama`](slice-v1-devcontainer-feature-ollama.md) (Umsetzung).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Devcontainer-Features; Erweiterung), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil; nur bei Netzdienst-Ausweisung), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts: unverändert). Verwandt: [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) und die Folge-ADR des Mess-Slice.
Anlass: Projektinhaber, 2026-10-18 („Lastenheft-CR: Katalog erweitert um `ollama`: Zusage nur für Gemessenes“).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), gegebenenfalls [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)), §16 (Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt der CR nicht; die Spezifikation folgt im letzten Slice der Welle ([`slice-v2-spec-sandbox-egress-ollama`](slice-v2-spec-sandbox-egress-ollama.md)).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Zusage „Ollama läuft als Katalog-Feature `ollama` im Devcontainer, CPU-only, Modelle im Container“ steht im Lastenheft (Änderung auf die **nächste freie Minor-Version**; `0.15.0` oder höher, je nach Reihenfolge des Landens mit dem Neu-Laden-CR, der nur bei Variante V2 entsteht), **mit den Messwerten des Mess-Slice** und in der vom Projektinhaber freigegebenen Fassung. Der Code-Slice beginnt erst danach.

**Warum eine Lastenheft-Änderung:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) führt Beispiele auf („Git, Docker CLI, Node.js, Java, SDKMAN, PostgreSQL Client, Kubernetes Tools“), keine abschließende Liste und keine Katalogschlüssel; `python`, `homebrew` und `sdkman` kamen ohne Lastenheft-Änderung (gelesen: [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md) §1, Absatz „Katalog-Schlüssel allein sind keine Zusage“). Für `ollama` verlangt der Projektinhaber ausdrücklich eine Zusage „nur für Gemessenes“; das ist eine Aussage über Produktverhalten (Installation, CPU-only, Modellspeicher) und gehört deshalb ins Lastenheft.

### Gefundener Stand (gelesen am 2026-10-18)

- [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features): Absatz „Beispiele:“ mit sieben Spiegelstrichen; danach „Für optionale externe Feature-Quellen gilt:“ mit den Regeln zur Freigabe ([`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)). Ein eingebauter Katalogeintrag braucht keine Freigabe (`featureCatalogue`, `internal/hexagon/application/devcontainer_features.go`: „Built-in entries are aktivierbar without an Allowlist entry“).
- Die Egress-Hosts der Features stehen nicht im Lastenheft, sondern in [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion); [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) sagt nur, die Default-Allowlist sei „dokumentiert“. Eine Lastenheft-Aussage zu den Ollama-Hosts ist deshalb **nicht** nötig.
- Das Lastenheft nennt heute weder `sshd` noch `claude-code` als Katalogschlüssel (Suche `grep -i "claude-code\|sshd" spec/lastenheft.md` am 2026-10-18: kein Treffer; das Lastenheft spricht vom „Login-Dienst“ und von „Zugangsdaten eines Agenten“, nicht von Schlüsseln des Katalogs); der CR führt `ollama` als **einzigen** namentlich genannten Katalogschlüssel ein, weil seine Zusage an Messbedingungen hängt.

### Wortlaut-Vorschlag (alt/neu; zur Freigabe durch den Projektinhaber vor dem Lastenheft-Commit)

Die in eckigen Klammern gesetzten **Messwerte** (`[…]`) sind Leerstellen, die der Mess-Slice füllt; ohne Messwert gilt die Aussage nicht, und der Satz entfällt. Das Lastenheft nennt keine Slice- oder Wellen-Kennung; im Lastenheft stehen die Verweise als Anker `#…` der Datei.

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

> Ollama ist als eingebautes Feature `ollama` wählbar (`u-boot config set devcontainer.features.ollama.enabled true`). Es installiert Ollama im Container ohne Grafikbeschleunigung (nur CPU); die Installation stammt aus einer gepinnten Quelle mit geprüfter Prüfsumme [`[Quelle und Prüfung aus der ADR des Mess-Slice]`]. Der Server läuft im Container [`[Bindung und Start aus der Messung: z. B. „nur auf der Loopback-Adresse“, „startet nicht von selbst“]`], Modelle lädt der Nutzer im Container (`ollama pull`), sie liegen in einem benannten Volume und überstehen den Neubau des Containers [`[nur, wenn gemessen]`]. Mit der Egress-Restriktion ([`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) gehören die Hosts des Modell-Pull zur Default-Allowlist. Gemessen ist das auf [`[Messstand: Docker-Version, Colima-Version, Architektur, Ollama-Version]`]; nicht zugesagt sind GPU-Beschleunigung, ein Anschluss an einen Ollama-Dienst des Hosts, andere Architekturen und Container-Engines als gemessen sowie Auswahl oder Verwaltung von Modellen durch `u-boot`.

**Stelle 2 — [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), Ausweisung eines Netzdienstes (nur wenn der Mess-Slice zeigt, dass der Server auf einer anderen als der Loopback-Adresse lauscht oder über `forwardPorts`/Engine von außen erreichbar wird; sonst entfällt die Stelle).**

Alt (letzter Spiegelstrich des Abschnitts „Das erzeugte Ergebnis im Sandbox-Profil muss:“):

> - einen im Container aktivierten Login-Dienst (Passwort-Login) in der Befehlsausgabe ausweisen.

Neu (Variante, nur bei Bedarf):

> - einen im Container aktivierten Login-Dienst (Passwort-Login) und einen im Container aktivierten Dienst, der von außerhalb des Containers erreichbar ist, in der Befehlsausgabe ausweisen.

Eine Änderung dieser Stelle zieht [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Ausnahme 1 nennt „eines im Container aktivierten Login-Dienstes“) mit; der CR nimmt sie dann als Stelle 2a auf, mit einem eigenen alt/neu.

**Keine Änderung** (gelesen, bewusst nicht angefasst): [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) — wird ein Drittanbieter-Feature gewählt, gilt die Regel von [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) (Warnung, entfällt unter `--quiet`) ohne Lastenheft-Änderung; wird ein Binär-Archiv gewählt, trägt die Auslegung die ADR des Mess-Slice (Muster [`ADR-0023`](../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md)).

**Stelle 3 — Versionsfeld und Historie (§16).**

Alt: `| Version          | 0.13.0 |` (bzw. der Stand nach dem Landen der Vorgänger-CR in dieser Welle)

Neu: die nächste freie Minor-Version (Lastenheft-Commit-Reihenfolge entscheidet; nach dem Anzeige-CR `0.15.0`, nach dem Neu-Laden-CR `0.16.0`).

Neue Historie-Zeile (Datum: das des Lastenheft-Commits):

> | `<Version>` | `<Datum des Lastenheft-Commits>` | [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features): Katalog-Feature `ollama` (Ollama im Container, nur CPU, Modelle im benannten Volume) für die gemessenen Hosts und Versionen; nicht zugesagt: GPU, Anschluss an einen Ollama des Hosts, Modellverwaltung. | Vereinbarung mit dem Projektinhaber |

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Messung, Quellenwahl und Folge-ADR** — [`slice-v1-devcontainer-ollama-messung`](../in-progress/slice-v1-devcontainer-ollama-messung.md); der CR trägt nur, was dort gemessen und entschieden ist.
- **Code, Tests, Spezifikation, Handbuch, CHANGELOG** — [`slice-v1-devcontainer-feature-ollama`](slice-v1-devcontainer-feature-ollama.md) und [`slice-v2-spec-sandbox-egress-ollama`](slice-v2-spec-sandbox-egress-ollama.md); das Lastenheft ist Zielbild, der Code folgt.
- **Eine Zusage zu Hosts, Version oder Modellen** — Hosts und Version gehören in die Spezifikation ([`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)), nicht ins Lastenheft; Modelle wählt der Nutzer.
- **Anschluss an den Ollama des Hosts, GPU, weitere Agenten** — vom Projektinhaber ausgeschlossen bzw. nicht genannt.

## 2. Definition of Done

- [ ] **Lastenheft `<Version>`:** Stellen 1 und 3 (und 2/2a, falls der Mess-Slice einen Netzdienst zeigt) in der vom Projektinhaber freigegebenen Fassung, **jede Klammer mit gemessenem Wert gefüllt oder der Satz gestrichen**, Versionsfeld, Historie-Zeile mit Datum (Verweis „Vereinbarung mit dem Projektinhaber“); die Freigabe ist im Gespräch belegt und in §7 vermerkt (Datei-Beleg besteht nicht).
- [ ] **Vorab-Suche und Nachzug:** Die Suche gegen ältere Lastenheft-Zusagen (Features, Sandbox-Regeln, Login-Dienst, `--quiet`) ist festgehalten; mitzuändernde Stellen sind genannt und geändert.
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
| `spec/lastenheft.md` | update | Stellen 1 bis 3 (§3 [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), gegebenenfalls [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging); Versionsfeld; §16 Historie); erst nach Freigabe des Wortlauts und nach dem Mess-Slice |
| `docs/plan/planning/open/slice-cr-devcontainer-ollama.md` (diese Datei) | update | Freigabe, Vorab-Suche, Closure (§7) |
| `docs/reviews/` (Review-Report, Verification Evidence) | neu | Review und Verifikation nach `harness/review.md` und `harness/verification.md` |
| `CHANGELOG.md` | unverändert | der CR ändert kein Produktverhalten; ein Eintrag entsteht im Code-Slice |

**Übergabe aus dem Entwurf [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (`Proposed`; gilt mit `Accepted`; Bedarf, kein Wortlaut):** Stelle 1 — Quelle und Prüfung: gepinntes Release-Archiv des Upstreams, nur per SHA-256 geprüft, **ohne Signatur** (die Grenze ist auszuweisen); Bindung: nur Loopback (gemessen); Start: nicht von selbst; Modelle im benannten Volume überstehen das **Löschen des Containers**, der Satz „überstehen den Neubau“ entfällt (nicht gemessen); wirksam nur mit Sandbox-Profil; Lauf auf `amd64` nicht verifiziert. Stelle 2 und 2a **entfallen** (Loopback, kein Netzdienst nach außen; Entscheidung 5 der ADR). [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) bleibt unverändert (die Auslegung steht in der ADR); die Hosts stehen nicht im Lastenheft.

Vorab-Suche (Pflicht vor jeder Vertragsänderung, Muster [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md)): `grep -i` über `spec/lastenheft.md` nach `Feature`, `Katalog`, `Login-Dienst` und `quiet`; der Implementer führt sie vor dem Lastenheft-Commit aus und trägt die Treffer hier ein.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; [`slice-v1-devcontainer-ollama-messung`](../in-progress/slice-v1-devcontainer-ollama-messung.md) liegt in `done/` (Messwerte und ADR-Status liegen vor); der Projektinhaber hat den Wortlaut freigegeben (Gespräch).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): der Wortlaut verlangt mehr als zwei Anforderungen oder eine Sicherheitsentscheidung, die der Mess-Slice nicht getroffen hat.
- `in-progress` → `open` (blockiert): keine Freigabe des Wortlauts, oder die ADR des Mess-Slice ist noch `Proposed`; das Lastenheft bleibt unverändert, kein Code-Slice startet.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). Der Lastenheft-Commit trägt die Freigabe nur, wenn §7 sie vermerkt.

## 6. Risiken und offene Punkte

- **Ein Satz im Lastenheft stützt sich auf einen Messwert, den der Mess-Slice nicht liefert.** — **Ausgang:** offen bis zum Start; der Satz entfällt (keine Aussage ohne Messwert).
- **Die Zusage „Modelle überstehen den Neubau“ ist ungemessen.** — **Ausgang:** offen bis O5 des Mess-Slice; der Satz bleibt nur mit Messwert.
- **Die Reihenfolge der Lastenheft-Versionen mit dem Anzeige-CR (und dem Neu-Laden-CR) kollidiert.** — **Ausgang:** offen; die Version ist die nächste freie zum Zeitpunkt des Commits, nicht vorab festgeschrieben.
- **Ein Review läuft vor der Freigabe** (siehe DoD). — **Ausgang:** offen; Review über die freigegebene Fassung.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: Freigabe des Wortlauts (Gespräch), Vorab-Suche, Steering-Loop-Eintrag, Register, Folge-Slices, Risiken.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird allein `spec/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); die Schwelle (zwei von drei Achsen) ist laut Modus-Deklaration erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; kein Treffer für Lastenheft-Änderungen unter den offenen Beobachtungen. [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) gelten (verkörpert in `AGENTS.md` §3.10): jede Bedingung im Wortlaut („nur auf der Loopback-Adresse“, „überstehen den Neubau“) braucht im Code-Slice einen Test oder entfällt.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
