# Slice CR: Forge-CLIs `gh` und `glab` als Katalog-Features — Lastenheft-CR ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) Katalog und Zusage-Grenzen, [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) Token-Variablen und `doctor`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-forge-cli-gh-glab`](../welle-forge-cli-gh-glab.md) (Release `v0.15.0`). Zweiter von fünf Slices (Reihenfolge Messung → CR → Code → Spezifikation → Release-Cut, Welle §4; Muster [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) und [`slice-cr-devcontainer-shell-und-zweites-up`](../done/slice-cr-devcontainer-shell-und-zweites-up.md)).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog, Beispiele, Zusage-Grenzen), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten eines Forge-Werkzeugs, Token-Variablen, `doctor`); [`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) (`Accepted`, Entscheidungen des Projektinhabers a bis d, f und h vom 2026-10-23). Unverändert und nur gelesen: [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Begründung §1, *Stelle 3*).

**Berührte Spec-Stellen:** `lastenheft.md` §4.3 ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)), §16 (Historie), Versionsfeld. `SPEC-*` und `ARC-*` berührt der CR nicht (das tut `slice-v2-spec-forge-cli` nach dem Code). Der Verweis zeigt aufwärts: Das Lastenheft nennt diesen Slice nie.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-23.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das Lastenheft sagt zu, was [`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) entschieden und das [Messprotokoll](../../../messungen/slice-v1-forge-cli-quellen-messung/README.md) getragen hat, und nicht mehr: `gh` und `glab` sind wählbare Katalog-Features mit gepinntem, geprüftem Release-Archiv, ausgewiesener Grenze der Lieferkette, ausgewiesenen Hosts und einem Token, das nur als Umgebungsvariable des Hosts in den Container kommt; `doctor` prüft die Quelle der Variable. Es ist das **Zielbild**: kein Messwert, kein Messstand, kein Prozesszustand, keine Slice-, Wellen- oder ADR-Kennung. Ob der Code es trägt, entscheiden die Tests des Code-Slice, nicht dieser CR.

**Wortlaut und Freigabe.** Der Wortlaut liegt als Entwurf vor ([`wortlaut-entwurf.md`](../../../messungen/slice-v1-forge-cli-quellen-messung/wortlaut-entwurf.md): Stelle 1 = [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) mit den Sätzen S1 bis S6, Stelle 2 = [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) mit drei Sätzen, Stelle 3 = keine Änderung, Stelle 4 = Version). Er ist bis zur Freigabe kein Vertragstext. Dieser Slice führt die Freigabe als **eigenen Schritt vor der Lastenheft-Änderung** (§3 Ablauf, Schritte 1 und 2): Der Implementer legt dem Projektinhaber den Wortlaut je Satz mit der Spalte *Nicht gemessen* vor, vermerkt die Entscheidung je Stelle in §7 und ändert das Lastenheft erst danach, wortgleich mit der Freigabe. Der Slice enthält keinen Lastenheft-Text; was freigegeben wird, trägt §7 (Fix-Runden wie im Muster als Ergänzung, nicht als stilles Umschreiben).

**Fragen für die Freigabe** (Planner, aus dem Lesen des Entwurfs gegen [`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md); keine Änderung des Entwurfs):

1. **Wirkung in beiden Profilen.** Entscheidung d und Entscheidung 3 der ADR sagen: der Block wirkt mit und ohne Sandbox-Profil. Der Ollama-Absatz des Lastenhefts nennt die Wirkung („wirkt im Sandbox-Profil“, gelesen Zeile 409); die Sätze S1 bis S6 nennen sie nicht. Soll ein Halbsatz zu S1 sie zusagen („in beiden Profilen“), oder genügt das Schweigen, bis die Spezifikation sie festlegt? Ohne Halbsatz schweigt der Vertrag zu einer Entscheidung des Projektinhabers (Muster [`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](../observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md)); mit ihm braucht die Bedingung einen Test (§6).
2. **Hosts im Lastenheft.** S4 nennt `api.github.com`, `github.com`, `gitlab.com`; der Ollama-Absatz nennt keinen Host (gelesen Zeile 409), die Hosts je Feature stehen in der Spezifikation (Entwurf, Stelle 3). Soll S4 die Namen tragen, oder entfällt der zweite Teil von S4?
3. **Version.** `0.22.0`, wenn dieser CR vor [`slice-cr-lastenheft-messstaende`](slice-cr-lastenheft-messstaende.md) nach `done/` geht; sonst die darauf folgende (§6).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktcode, Katalogeinträge, Template, Tests** — der CR schreibt den Vertrag, der Code folgt ihm (`slice-v1-devcontainer-feature-gh-glab`, Welle §4).
- **`spec/spezifikation.md` und `spec/architecture.md`** — nach dem Code in `slice-v2-spec-forge-cli`; [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) und [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) präzisieren das Lastenheft und erweitern es nie ([`AGENTS.md`](../../../../AGENTS.md) §2).
- **Handbuch, READMEs, `CHANGELOG.md`** — es ändert sich kein Verhalten; die Code-Slices bringen die Dokumente mit.
- **Änderungen an [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) und [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)** — der Entwurf begründet je Anforderung „keine Änderung“ (Stelle 3); die Begründung gehört zur Freigabe, sie entscheidet der Projektinhaber.
- **Die Messstände in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Ollama-Absatz, „Gemessen ist das auf Docker 29.5.2 …“) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)** — eigener CR ([`slice-cr-lastenheft-messstaende`](slice-cr-lastenheft-messstaende.md)); der neue Absatz dieses CR trägt keinen Messstand (Entwurf, Konvention).
- **Eine Nachmessung mit gültigem Token und eine Zusage über das Gemessene hinaus** — vom Projektinhaber verworfen ([`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) Punkt f); der Erfolgspfad nach der Anmeldung bleibt „nicht zugesagt“ (S5).
- **Die Release-Nummer** — Lastenheft-Version ist nicht Release-Version; `v0.15.0` schneidet der Release-Cut der Welle.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten darunter zählen nicht mit.

- [ ] **Freigabe des Projektinhabers vor der Änderung:** Der Wortlaut der Stellen 1 und 2 des Entwurfs ist dem Projektinhaber je Satz (S1 bis S6, die drei Sätze der Stelle 2) mit der Spalte *Nicht gemessen* vorgelegt, die Fragen 1 und 2 aus §1 sind beantwortet, und die Entscheidung je Stelle steht in §7 mit Datum. *Beleg:* §7 *Freigabe*; kein Commit von `spec/lastenheft.md` vor diesem Eintrag (Commit-Reihenfolge in `git log`).
- [ ] **Wortlaut im Lastenheft** nach der Freigabe, wortgleich mit ihr; Zielbild ohne Messwert, ohne Messstand, ohne Prozesszustand und ohne Verweis auf ADR, Slice oder Roadmap (Decken-Regel, [`AGENTS.md`](../../../../AGENTS.md) §3.4). Jede Bedingung im Text ist durch eine Messaussage des Protokolls oder durch ein Zielbild des Code-Slice mit benanntem Test belegt (Beleg-Tabelle in §7, [`AGENTS.md`](../../../../AGENTS.md) §3.10). *Beleg:* Vergleich Satz für Satz im Verification-Report, Fundstellen in §7.
- [ ] **Versionsfeld und Historie:** die nächste freie Minor-Version im Versionsfeld (Zeile 9, heute `0.21.0`) und eine Historie-Zeile in §16, die nur die freigegebenen Stellen nennt. *Beleg:* Versionsfeld, Historie-Zeile.
- [ ] `make gates` grün, `make docs-check` grün. `make doc-trace` ist für [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) ein Beleg (Familie `LH-FA-*`, `.d-check.yml` Abschnitt `trace`). *Beleg:* Läufe und Exit-Codes in §7.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8). Ändert eine Fix-Runde Lastenheft-Text, liest ein Nachlauf-Review genau diesen Diff, bevor der Verifier startet ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)).
- [ ] Doku-Update: keines — der CR ändert kein Verhalten; Roadmap und `open/README.md` zeigen den Zustand.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — u-boot führt Wellen; die Welle-Closure prüft sie.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

**Ablauf** (die Reihenfolge ist Teil des Plans):

1. **Vorlage.** Der Implementer fasst den Wortlaut der Stellen 1 und 2 des Entwurfs mit der Spalte *Nicht gemessen* und den Fragen aus §1 für den Projektinhaber zusammen (Rückfrage im Gespräch; kein Lastenheft-Commit).
2. **Freigabe.** Die Antwort je Stelle (Ja / Ja mit Änderung / Nein) steht in §7 *Freigabe*; ein geänderter Satz steht dort im freigegebenen Wortlaut. Der Entwurf wird nicht umgeschrieben (er ist Messprotokoll-Anlage des fertigen Slice `slice-v1-forge-cli-quellen-messung`); die Abweichung trägt §7.
3. **Lastenheft.** `spec/lastenheft.md` ändert sich genau im freigegebenen Wortlaut (ein Commit), danach Versionsfeld und Historie-Zeile.
4. **Vorab-Suche** vor Schritt 3 mit Fundstellen in §7 (Muster [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md)): `gh`, `glab`, `GitHub`, `GitLab`, `Forge`, `Ollama`, `claude-code`, `Katalog`, `Beispiele`, `GIT_TOKEN`, `Agenten-Feature`, `Zugangsschlüssel`, `Umgebungsvariable`. Gelesen am 2026-10-23 (Planner): `gh`/`glab`-Treffer keine; `Ollama` in Zeilen 407, 409 und 1587 (Historie `0.15.0`); `GitHub` nur in [`LH-QA-003`](../../../../spec/lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions) (Zeile 1253, Fremdkontext, unverändert); `claude-code`, `ANTHROPIC`, `GIT_TOKEN` und `Forge` kommen im Lastenheft nicht vor. Eine Aufzählung aller Katalog-Features außer den Spiegelstrichen in Zeilen 396 bis 407 fand der Planner nicht; die Suche des Implementers wiederholt sie.
5. Review, gegebenenfalls Fix-Runde mit Nachlauf-Review, Verification, Closure.

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; sie wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Eine Zeile „prüfen“ nennt eine Datei, die der Planner nicht gelesen hat.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Stellen 1 und 2 und Stelle 4 des Entwurfs im freigegebenen Wortlaut. Gelesen am 2026-10-23: Versionsfeld `0.21.0` (Zeile 9); [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) Spiegelstriche Zeilen 396 bis 407 (letzter: „Ollama (lokale Sprachmodelle, nur CPU)“), Ollama-Absatz Zeile 409, „Für optionale externe Feature-Quellen gilt:“ Zeile 411 (der neue Absatz steht dazwischen); [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) Einleitungssatz Zeile 482, Spiegelstriche zur Dokumentation und zu `doctor` Zeilen 485 und 486; §16 beginnt Zeile 1552, letzte Historie-Zeile `0.21.0` Zeile 1593 (Datei 1600 Zeilen) |
| `docs/plan/planning/in-progress/roadmap.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zeile „Nichts in Arbeit.“ (gelesen) bei Start auf diesen Slice, bei der Closure zurück (`.d-check.yml`, Abschnitt `planning`; `make doc-planning`); die Zeile „Offen:“ nennt „die Kennungen ab dem zweiten vorläufig“ und zählt ab dem dritten Slice (Anlage dieses Plans, Planner-Commit) |
| `docs/plan/planning/open/README.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zeile dieses Slice (Anlage durch den Planner; bei der Closure „geliefert (liegt in `done/`)“); Zeile von [`slice-cr-lastenheft-messstaende`](slice-cr-lastenheft-messstaende.md) (gelesen Zeile 33: „nach 0.21.0“) auf die Version nach diesem CR, falls dieser CR zuerst läuft |
| `docs/plan/planning/open/slice-cr-lastenheft-messstaende.md` <!-- d-check:ignore (Pfad-Zitat) --> | update (bedingt) | Versionsangaben (gelesen: Zeile 52 „nach 0.21.0“, Zeile 65 „zuletzt `0.21.0`“) auf die Version dieses CR, falls dieser CR zuerst nach `done/` geht (Muster [`slice-cr-devcontainer-shell-und-zweites-up`](../done/slice-cr-devcontainer-shell-und-zweites-up.md) §6); entfällt, wenn jener vorher läuft — dann ändert dieser CR seine eigene Versionsangabe |
| diese Datei | update | `Verantwortlich:`, Freigabe in §7, Plan-Abgleich, Closure; Moves per `make slice-mv` (eigene Commits, [`AGENTS.md`](../../../../AGENTS.md) §3.3) |
| `docs/plan/planning/welle-forge-cli-gh-glab.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | §4 Tabelle: Zeile dieses Slice wird Link statt „(Datei entsteht nach der Messung)“ (gelesen Zeile 67), der Satz „Die Kennungen der Slices ab Zeile 2 sind vorläufig“ (gelesen Zeile 72) zählt ab Zeile 3 (Anlage, Planner-Commit); Verweis-Nachzug durch `make slice-mv` bei den Moves; die vergebene Lastenheft-Version bei der Closure |
| `docs/reviews/` <!-- d-check:ignore (Pfad-Zitat) --> | neu | Review-Report, Nachlauf-Review(s) und Verification-Report, von Reviewer bzw. Verifier geschrieben |
| `docs/plan/planning/observations/` <!-- d-check:ignore (Pfad-Zitat) --> | update (Closure) | Evidence-Dateien für angefallene Beobachtungen (§7) |
| `spec/spezifikation.md`, `spec/architecture.md`, `docs/user/`, `README.md`, `README.de.md`, `CHANGELOG.md`, `docs/maintainer/`, `harness/replay.md`, `internal/`, `cmd/`, `tools/`, `docs/messungen/` <!-- d-check:ignore (Pfad-Zitate) --> | **unverändert** | kein Code, kein Verhalten. Gelesen als Stand, den die Folge-Slices ändern oder der Beleg bleibt: der Katalog kennt weder `gh` noch `glab` (`featureCatalogue`, `internal/hexagon/application/devcontainer_features.go`; ADR-0035 *Stand im Code*, am Typ `featureCatalogueEntry` Feld `dockerfileBlock` gelesen); `remoteEnv` mit `${localEnv:…}` steht in `devcontainer.json.tmpl` Zeilen 33 bis 36, die Namen kommen aus `agentCredentialEnv` (`devcontainer_sandbox.go`, heute nur `ANTHROPIC_API_KEY` für `claude-code`); `checkSandboxCredentials` in `doctor_sandbox.go` Zeile 133; `Dockerfile.tmpl` führt `USER root` unter `{{- if or .UID .Sandbox}}` (Zeilen 6 und 7). Der Entwurf `docs/messungen/slice-v1-forge-cli-quellen-messung/wortlaut-entwurf.md` bleibt unverändert |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; [`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) ist `Accepted` (eingetreten: `4655bcff`) und [`slice-v1-forge-cli-quellen-messung`](../done/slice-v1-forge-cli-quellen-messung.md) liegt in `done/` (eingetreten); kein anderer Slice in `in-progress/` (WIP-Limit 1; die Slices ohne Welle in `open/` teilen sich das Limit). Die Freigabe des Wortlauts ist **kein** Start-Trigger, sondern Schritt 2 des Ablaufs (§3) und Bedingung für den Lastenheft-Commit.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Projektinhaber verlangt eine Zusage, die eine neue Entscheidung oder Messung braucht (etwa ein Erfolgspfad mit gültigem Token oder `glab duo`); dann zuerst eine Folge-ADR.
- `in-progress` → `open` (blockiert): die Freigabe des Wortlauts liegt nicht vor oder der Projektinhaber verlangt einen anderen Wortlaut als den Entwurf, der Rückfragen nach sich zieht; §7 nennt die offene Stelle.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; der Wortlaut im Lastenheft ist mit der Freigabe Satz für Satz verglichen (Verification-Report), der Review-Report liegt vor, und die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Der Wortlaut sagt mehr zu, als die Messung trägt.** Die Spalte *Nicht gemessen* des Entwurfs nennt je Satz die Grenze (kein Erfolgspfad mit gültigem Token, `arm64`, `gh attestation verify`, Hosts nach Anmeldung); S5 und S6 sagen die Grenzen als „nicht zugesagt“ ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), 17 Belege). — **Ausgang:** offen bis zur Closure; Prüfung im Review: jeder Satz gegen Messaussage und Spalte *Nicht gemessen*.
- **Eine Bedingung des Wortlauts ohne Test** („wählbar“, „gepinntes Archiv mit geprüfter Prüfsumme“, „Ausweisung beim Erzeugen“, „Hosts in der Default-Allowlist“, „führt weder `gh auth login` noch `glab auth login` aus“, „`doctor` prüft die Variable“, gegebenenfalls „in beiden Profilen“ nach Frage 1; [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), [`AGENTS.md`](../../../../AGENTS.md) §3.10). Der Entwurf führt sie als Zielbild des Code-Slice. — **Ausgang:** offen bis zur Closure; die Beleg-Tabelle in §7 nennt je Bedingung den geplanten Test im Code-Slice `slice-v1-devcontainer-feature-gh-glab` (Eigentümer in dessen DoD, Welle §4) oder die Streichung.
- **Der Wortlaut ändert sich nach dem Review.** — **Ausgang:** offen bis zur Closure; jede Fix-Runde mit Lastenheft-Diff bekommt ein Nachlauf-Review über genau diesen Diff ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md), 14 Belege).
- **Versionsfolge mit [`slice-cr-lastenheft-messstaende`](slice-cr-lastenheft-messstaende.md).** Beide CR ändern [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (jener den Ollama-Absatz, dieser Beispiele und einen Absatz dahinter) und beanspruchen „die nächste freie Version nach 0.21.0“. — **Ausgang:** offen bis zur Closure; wer zuerst nach `done/` geht, belegt `0.22.0`, der andere nennt die folgende Version und zieht die Angaben des anderen in der Plan-Tabelle nach; beim zweiten Lauf wird der Diff gegen die Zeilen des ersten gelesen (gleicher Abschnitt).
- **Eine Aussage über den Code oder das Lastenheft im Plan stimmt nicht** (Zeilennummern, Symbole, Zustand der Welle). — **Ausgang:** offen bis zur Closure; Abgleich der Plan-Tabelle gegen `git diff --stat` nach der letzten Fix-Runde ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), 44 Belege, verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10).
- **Die Freigabe fällt anders aus als der Entwurf** (etwa S4 ohne Hostnamen, S3 ohne `gh`-Attestierung, anderer Halbsatz zu den Profilen). — **Ausgang:** offen bis zur Closure; die Freigabe in §7 trägt den geänderten Wortlaut, die Beleg-Tabelle folgt ihm; trägt der geänderte Wortlaut eine Zusage über die Messung hinaus, ist das ein Rückführungsgrund (§4).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

*(Wird bei der Umsetzung und Closure gefüllt, in dieser Reihenfolge: **Freigabe** je Stelle mit Datum und freigegebenem Wortlaut · Vorab-Suche mit Fundstellen · Beleg-Tabelle (Bedingung → Messaussage oder Test) · Verifikation mit Exit-Codes · Was hat funktioniert · Was ging anders · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird `spec/` (Lastenheft) und `docs/plan/` (Planung), beide in der Sub-Area mit dem Kürzel DOK laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) (Greenfield; auditiert, drei Achsen erfüllt). Keine Code-Sub-Area.

**Vorgelagert — offene Beobachtungen sichten** (Register gelesen am 2026-10-23; Zähler = Dateien in `evidence/`): [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (17) — trifft diesen CR unmittelbar, Risiko in §6; [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (14) — Risiko in §6; [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) (4) — Vorab-Suche in §3; [`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](../observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md) (3) — Frage 1 in §1; [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) (4) — S1 und S2 stützen sich auf Handvarianten (`ghtgz`, `glabtgz`), die Spalte *Nicht gemessen* des Entwurfs nennt das; [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (44, verkörpert) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (18, verkörpert) — Prüfpflicht, Risiken in §6. Evidenz-/Diskrepanz-Risiko: mittel (Zusage über Werkzeuge Dritter, Erfolgspfad ungemessen).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
