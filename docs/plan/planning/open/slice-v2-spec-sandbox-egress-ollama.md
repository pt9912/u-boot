# Slice v2: Spezifikation — Egress-Anzeige, Subdomain-Aussage, Katalog `ollama` und Hosts nachziehen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) — letzter Slice der Welle. Er startet **erst**, wenn jeder Code- und CR-Slice der Welle in `done/` liegt (Reihenfolge-Regel „Spezifikation nach dem Code“, Muster [`slice-v2-spec-devcontainer-sshd-claude-code`](../done/slice-v2-spec-devcontainer-sshd-claude-code.md)): [`slice-v1-egress-wildcard-pruefung`](../done/slice-v1-egress-wildcard-pruefung.md), [`slice-cr-sandbox-egress-anzeige`](slice-cr-sandbox-egress-anzeige.md), [`slice-v1-cli-egress-anzeige`](slice-v1-cli-egress-anzeige.md), [`slice-v1-sandbox-egress-neuladen-entwurf`](../done/slice-v1-sandbox-egress-neuladen-entwurf.md) (Ergebnis R0: keine Folge-Slices), [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md), [`slice-cr-devcontainer-ollama`](slice-cr-devcontainer-ollama.md), [`slice-v1-devcontainer-feature-ollama`](slice-v1-devcontainer-feature-ollama.md).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (die Anforderungen, die die Spezifikation verfeinert); [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) und die Folge-ADRs der Welle.

**Berührte Spec-Stellen:** `spezifikation.md`: [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Hosts, Zeile *Nutzer*, Zeile `ollama`), [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (Katalog), [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Aufzählung der Abfragebefehle; dazu die Aufzählung der ungekennzeichneten Meldungen), [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) (Warnungstabelle und Satz der *Sandbox-Regel*, nur bei gebauter Ollama-Warnung), [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (durch R0 nicht berührt: kein Neu-Laden vom Host; die offene Frage zum Absatz *Grenze* steht beim Projektinhaber, siehe den Entwurfs-Slice §6), ein neuer Eintrag zur Egress-Anzeige (die auf die höchste vergebene Kennung folgende; geprüft 2026-10-18: die höchste ist [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)), gegebenenfalls ein Eintrag zu `ollama` (nächste freie Nummer danach), [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) (nur bei neuem Diagnose-Code), §7 Historie; `architecture.md` §2.3 (Driving-Port, neuer Use Case) und §2.5 (Adapter). Das Lastenheft ändert dieser Slice **nicht**.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Spezifikation beschreibt das Verhalten, das die Code-Slices der Welle geliefert haben — gelesen gegen den Code, je Aussage mit Fundstelle (Datei und Symbol) —, und die Stellen, an denen der Vertrag steht, sagen dasselbe (Lastenheft, Spezifikation, Handbuch, CHANGELOG).

### Gefundener Stand (gelesen am 2026-10-18; nach den Code-Slices neu zu lesen)

- [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion): Quellen-Tabelle ohne `ollama`; Zeile *Nutzer*: „kleingeschriebene Hostnamen, ohne Schema, Port und Wildcard“; **kein Satz über Subdomains**. Die Subdomain-Aussage steht im Handbuch (`docs/user/benutzerhandbuch.md`, Abschnitt *Internet-Beschränkung (Egress)*, Punkte *Subdomains*, *Namen hinter einem CNAME* und *Messstand*, mit Messung; gelesen am 2026-10-18) und in [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) (Entscheidung 2); die Messung stützt sie, der Fall hinter einem `CNAME` ist dort als nicht allgemein gemessen ausgewiesen.
- [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen): Katalog mit Schlüsseln bis `claude-code`, kein `ollama`.
- [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet): Aufzählung „Die Abfragebefehle `config get`, `config list`, `logs` und `template list` behalten ihre Ausgabe.“
- [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung): „Ermittlung durch den Host“ schließt ein Programm im Container aus; das bleibt unverändert, weil nach R0 kein Neu-Laden vom Host entsteht.
- Höchste vergebene Kennung in `spezifikation.md`: [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code, Tests, Handbuch, CHANGELOG** — die Code-Slices haben sie geliefert; dieser Slice ändert sie nicht (stellt er eine Abweichung fest, ist das ein Befund und ein neuer Slice, kein stilles Umschreiben).
- **Das Lastenheft** — Zielbild; die CR-Slices haben es geändert.
- **Eine Spezifikation von Gemessenem ohne Beleg** — jede Zahl und jeder Host stammt aus Code oder Messprotokoll (Fundstellen-Tabelle); eine Aussage, die weder Code noch Messung trägt, steht nicht im Text.
- **Ein Eintrag zu Varianten, die nicht gebaut wurden** (Neu-Laden bei R0) — die Spezifikation beschreibt den Ist-Stand.

## 2. Definition of Done

- [ ] **Spezifikation nach dem Code:** Die Einträge zu Egress-Anzeige (Befehl, Ausgabeformen, Herkunftsnamen, Aktivierungszustände), Subdomain-Aussage ([`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), im Wortlaut der Messung), Katalog und Hosts von `ollama` ([`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), gegebenenfalls ein eigener Eintrag), die Nachführung von [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Abfragebefehl) stehen in der Spezifikation; jede Aussage trägt in der Closure eine **Fundstellen-Tabelle** (Datei und Symbol, gelesen).
- [ ] **Sicht und Historie:** `architecture.md` (§2.3, §2.5, ohne Wellen-, Slice- und ADR-Bezug, `AGENTS.md` §3.4) zeigt den neuen Use Case und Adapter; §7 Historie der Spezifikation trägt je Änderung eine Zeile; `make doc-trace` meldet keine Waise; `make docs-check` 0 Befunde.
- [ ] **Ein Vertrag, vier Stellen:** Stichprobe des Verifiers: Lastenheft, Spezifikation, Handbuch und CHANGELOG sagen über Anzeige-Befehl, Subdomain-Semantik, den Weg „Liste ändern“ (Neuaufbau) und `ollama` (Hosts, Volume, Start, Grenzen) dasselbe wie der Code.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update, falls ein öffentlicher Vertrag berührt ist (Handbuch nur bei festgestellter Abweichung).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen von der Welle-Closure geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | update | [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) (Ollama-Warnung und „keine Wirkung ohne Sandbox-Profil“; der Satz „keines von beiden ergänzt `runArgs`, `mounts` oder einen Startbefehl“ bleibt auf `sshd` und `claude-code` begrenzt), bedingt [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), neue Einträge (nächste freie Nummern), §7 Historie; je Aussage gegen den Code gelesen |
| `spec/architecture.md` | update | §2.3 und §2.5: Use Case und Adapter der Egress-Anzeige (Pfade, keine Slice-/Wellen-/ADR-Bezüge) |
| `internal/adapter/driving/cli/jsontestutil/coderegistry.go`, [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) | unverändert (bedingt update) | nur bei einem neuen Diagnose-Code; der Drift-Test (`drift_test.go`) gleicht beide ab |
| `docs/reviews/` (Review-Report, Verification Evidence) | neu | Review und Verifikation |
| diese Slice-Datei | update | Fundstellen-Tabelle je Aussage (§3), Closure (§7) |

**Übergabe aus dem Entwurf [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (`Accepted` am 2026-10-18, Freigabe-Punkte a bis e wie vorgeschlagen; Bedarf, kein Wortlaut):** [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen): Katalogeintrag `ollama` vom Typ Dockerfile-Block (keine OCI-Quelle, kein Render-Schlüssel), Pin und Prüfung, nur mit Sandbox-Profil, `version:` ohne `source:` ist Exit-Code `10`, `source:` macht ihn extern; [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion): Zeile `ollama` mit zwei Hosts und dem Hinweis auf die ungemessene Stabilität des Blob-Hosts. Der zweite Eintragstyp ist damit gebaut und gehört in die Beschreibung.

Die Fundstellen-Tabelle wird beim Schreiben der Spezifikation hier eingetragen (Aussage → Datei und Symbol, gelesen am Tag des Schreibens). Abgleich der Tabelle oben gegen `git diff --stat` nach der **letzten** Fix-Runde (`AGENTS.md` §3.10).

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; alle Code- und CR-Slices der Welle (siehe Feld *Welle*) liegen in `done/`; kein anderer Slice in `in-progress/` (WIP-Limit 1). Die Reihenfolge belegt der `git log` (Move-Commits der Code-Slices vor dem ersten Commit dieses Slice).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): die Spezifikation verlangt mehr als drei Einträge neu (Zerlegung nach Anzeige, Egress-Liste, `ollama`).
- `in-progress` → `open` (blockiert): ein Code-Slice der Welle ist nicht geschlossen oder stellt sich beim Lesen als abweichend vom Lastenheft heraus (Befund, neuer Slice).

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7); die Fundstellen-Tabelle trägt jede neue Aussage.

## 6. Risiken und offene Punkte

- **Die Spezifikation beschreibt, was nur das Handbuch sagt, ohne Code-Entsprechung** (etwa die Subdomain-Aussage ohne Test). — **Ausgang:** offen bis zur Closure; erwartet: entfallen, weil der Wildcard-Slice Test oder Streichung geliefert hat.
- **Die auf die höchste vergebene Kennung folgenden Nummern sind zum Start vergeben.** — **Ausgang:** offen; die Nummer ist die nächste freie zum Zeitpunkt des Schreibens (Kennung nach dem Schema der `harness/conventions.md`), nicht vorab festgeschrieben.
- **Ein Code-Slice der Welle ändert sich nach dem Start** (Fix-Runde, Rückführung). — **Ausgang:** offen; die Fundstellen-Tabelle wird nach der letzten Änderung erneut gegen den Code gelesen (`AGENTS.md` §3.10, Abgleich nach der letzten Fix-Runde).
- **Einträge zur Ollama-Warnung hängen an gebautem Code.** Der zweite Eintragstyp und die Ollama-Warnung sind mit [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (`Accepted`) entschieden und gehören in die Spezifikation, sobald der Code-Slice sie gebaut hat. — **Ausgang:** offen bis zum Start; die DoD zählt nur, was gebaut ist.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: Fundstellen-Tabelle, Stichprobe „ein Vertrag, vier Stellen“, Steering-Loop-Eintrag, Register, Folge-Slices, Risiken.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `spec/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); die Schwelle (zwei von drei Achsen) ist laut Modus-Deklaration erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (beide `verkörpert` in `AGENTS.md` §3.10; die Fundstellen-Tabelle ist die Antwort), [`BEO-WZG/kopplung-nicht-nachgezogen`](../observations/BEO-WZG/kopplung-nicht-nachgezogen/observation.md) (offen, zwei Belege; Spezifikation, Sicht und Historie sind zusammen nachzuziehen). Evidenz-/Diskrepanz-Risiko: mittel.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
