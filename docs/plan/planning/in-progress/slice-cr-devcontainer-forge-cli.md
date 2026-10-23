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

**Verantwortlich:** pt9912 (Rolle Implementer, Start 2026-10-23); Freigaben: Projektinhaber.

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
3. **Version.** `0.22.0`, wenn dieser CR vor [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) nach `done/` geht; sonst die darauf folgende (§6).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktcode, Katalogeinträge, Template, Tests** — der CR schreibt den Vertrag, der Code folgt ihm (`slice-v1-devcontainer-feature-gh-glab`, Welle §4).
- **`spec/spezifikation.md` und `spec/architecture.md`** — nach dem Code in `slice-v2-spec-forge-cli`; [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) und [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) präzisieren das Lastenheft und erweitern es nie ([`AGENTS.md`](../../../../AGENTS.md) §2).
- **Handbuch, READMEs, `CHANGELOG.md`** — es ändert sich kein Verhalten; die Code-Slices bringen die Dokumente mit.
- **Änderungen an [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) und [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)** — der Entwurf begründet je Anforderung „keine Änderung“ (Stelle 3); die Begründung gehört zur Freigabe, sie entscheidet der Projektinhaber.
- **Die Messstände in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Ollama-Absatz, „Gemessen ist das auf Docker 29.5.2 …“) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)** — eigener CR ([`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md)); der neue Absatz dieses CR trägt keinen Messstand (Entwurf, Konvention).
- **Eine Nachmessung mit gültigem Token und eine Zusage über das Gemessene hinaus** — vom Projektinhaber verworfen ([`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) Punkt f); der Erfolgspfad nach der Anmeldung bleibt „nicht zugesagt“ (S5).
- **Die Release-Nummer** — Lastenheft-Version ist nicht Release-Version; `v0.15.0` schneidet der Release-Cut der Welle.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten darunter zählen nicht mit.

- [x] **Freigabe des Projektinhabers vor der Änderung:** Der Wortlaut der Stellen 1 und 2 des Entwurfs ist dem Projektinhaber je Satz (S1 bis S6, die drei Sätze der Stelle 2) mit der Spalte *Nicht gemessen* vorgelegt, die Fragen 1 und 2 aus §1 sind beantwortet, und die Entscheidung je Stelle steht in §7 mit Datum. *Beleg:* §7 *Freigabe*; kein Commit von `spec/lastenheft.md` vor diesem Eintrag (Commit-Reihenfolge in `git log`).
- [x] **Wortlaut im Lastenheft** nach der Freigabe, wortgleich mit ihr; Zielbild ohne Messwert, ohne Messstand, ohne Prozesszustand und ohne Verweis auf ADR, Slice oder Roadmap (Decken-Regel, [`AGENTS.md`](../../../../AGENTS.md) §3.4). Jede Bedingung im Text ist durch eine Messaussage des Protokolls oder durch ein Zielbild des Code-Slice mit benanntem Test belegt (Beleg-Tabelle in §7, [`AGENTS.md`](../../../../AGENTS.md) §3.10). *Beleg:* Vergleich Satz für Satz im Verification-Report, Fundstellen in §7.
- [x] **Versionsfeld und Historie:** die nächste freie Minor-Version im Versionsfeld (Zeile 9, heute `0.21.0`) und eine Historie-Zeile in §16, die nur die freigegebenen Stellen nennt. *Beleg:* Versionsfeld, Historie-Zeile.
- [x] `make gates` grün, `make docs-check` grün. `make doc-trace` ist für [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) ein Beleg (Familie `LH-FA-*`, `.d-check.yml` Abschnitt `trace`). *Beleg:* Läufe und Exit-Codes in §7.
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
| `docs/plan/planning/open/README.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zeile dieses Slice (Anlage durch den Planner; bei der Closure „geliefert (liegt in `done/`)“); Zeile von [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) (gelesen Zeile 33: „nach 0.21.0“) auf die Version nach diesem CR, falls dieser CR zuerst läuft |
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
- **Versionsfolge mit [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md).** Beide CR ändern [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (jener den Ollama-Absatz, dieser Beispiele und einen Absatz dahinter) und beanspruchen „die nächste freie Version nach 0.21.0“. — **Ausgang:** offen bis zur Closure; wer zuerst nach `done/` geht, belegt `0.22.0`, der andere nennt die folgende Version und zieht die Angaben des anderen in der Plan-Tabelle nach; beim zweiten Lauf wird der Diff gegen die Zeilen des ersten gelesen (gleicher Abschnitt).
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

Reihenfolge dieser Sektion: Freigabe je Stelle (Rückfrage beim Projektinhaber, beim Start vermerkt) · Vorab-Suche mit Fundstellen · Beleg-Tabelle (Bedingung → Messaussage oder Test) · Verifikation mit Exit-Codes · Closure.

### Freigabe

**Entscheidungen des Projektinhabers (2026-10-23, Rückfrage im Gespräch; an den Implementer übermittelt durch den Hauptagenten, beim Start vermerkt).** Der Wortlaut ist je Stelle freigegeben; ein geänderter Satz steht hier im freigegebenen Wortlaut. Der Entwurf ([`wortlaut-entwurf.md`](../../../messungen/slice-v1-forge-cli-quellen-messung/wortlaut-entwurf.md)) bleibt unverändert; die Abweichungen tragen die Punkte 1 und 2.

1. **Stelle 1 ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features): zwei Spiegelstriche `gh`/`glab` und ein Absatz) — freigegeben,** mit den Antworten zu Frage 1 und 2 aus §1:
   - **Frage 1 — Ja:** S1 trägt den Halbsatz „in beiden Profilen (Standard und Sandbox)“. Die Bedingung braucht im Code-Slice einen pinnenden Test (Zeile B1 der Beleg-Tabelle unten; Eintrag in der Welle-Datei §5).
   - **Frage 2 — Nein:** Die Hostnamen (`api.github.com`, `github.com`, `gitlab.com`) stehen in der Spezifikation, nicht im Lastenheft; die Hostnamen-Passage von S4 entfällt und S4 lautet hostnamenfrei. Der Implementer hat das so umgesetzt (Umsetzung der Antwort, nicht ein Wortlaut des Projektinhabers; der Reviewer prüft sie).
   - Freigegebener Absatz (Umsetzung im Lastenheft: S1 mit dem Halbsatz umgestellt, sonst wortgleich): „`gh` (GitHub CLI) und `glab` (GitLab CLI) sind in beiden Profilen (Standard und Sandbox) als eingebaute Features `gh` und `glab` wählbar (`u-boot config set devcontainer.features.gh.enabled true`, analog `glab`). Jedes installiert sein Werkzeug aus einem gepinnten Release-Archiv des Upstream-Projekts (`cli/cli`, `gitlab-org/cli`) mit geprüfter SHA-256-Prüfsumme. Für die Archive hat die Messung keine Signatur gefunden (nach dem Suchmuster, nicht abschließend); für `gh` veröffentlicht das Projekt eine GitHub-Attestierung, die `u-boot` nicht prüft. `u-boot` weist das beim Erzeugen aus. Mit der Egress-Restriktion ([`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) stehen die Hosts von `gh` in der Default-Allowlist und kommt für `glab` ein Host hinzu. Anmeldung und Zugangsdaten: siehe [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer); `u-boot` führt weder `gh auth login` noch `glab auth login` aus und speichert kein Token. Nicht zugesagt: Aufrufe mit gültigem Token und deren Hosts, Telemetrie-/Update-Hosts einer interaktiven Sitzung, selbst gehostete Instanzen, andere Versionen als die gepinnten; außerdem `glab duo`, `glab orbit` und das Nachladen weiterer Programme durch `glab`.“
2. **Stelle 2 ([`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer): drei Passagen) — freigegeben** wie im Entwurf (*Neu*-Wortlaut unverändert).
3. **Stelle 3 — keine Änderung** an [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) und [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (kein Widerspruch des Projektinhabers zum Entwurf; die Hostnamen-Antwort bestätigt die Ablage in der Spezifikation).
4. **Stelle 4 — Version `0.22.0` (Frage 3),** wenn dieser CR vor [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) nach `done/` geht; die Versionsangaben jenes Slice und seiner Zeile in `open/README.md` werden mit dem Lastenheft-Commit nachgezogen (§3, §6). Die Historie-Zeile nennt nur die freigegebenen Stellen.

### Vorab-Suche

Gelesen am 2026-10-23 (Implementer) in `spec/lastenheft.md` vor der Änderung, Suchwörter wie in §3 Schritt 4: `gh`/`glab` keine Treffer; `GitHub`/`GitLab`/`Forge`: nur [`LH-QA-003`](../../../../spec/lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions) (Fremdkontext, unverändert); `claude-code`, `GIT_TOKEN`: keine; `Ollama`: Spiegelstrich und Absatz in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Zeilen 407 und 409 vor der Änderung, danach 407 und 411; der neue Absatz steht in Zeile 413 vor „Für optionale externe Feature-Quellen gilt:“) und die Historie-Zeile `0.15.0`; `Katalog`/`Beispiele`: die Liste in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (die anderen „Beispiele:“-Treffer, Zeilen 162, 528, 597, 1006, 1199, 1211 vor der Änderung, gehören zu anderen Anforderungen); `Agenten-Feature`/`Zugangsschlüssel`/`Umgebungsvariable`: [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Einleitungssatz und zwei Spiegelstriche, geändert) und [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zeilen 498 und 500 nach der Änderung, „Zugangsdaten gibt der Befehl nur als Umgebungsvariablen weiter …“ mit Verweis auf [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) — allgemein gefasst, trägt keine Liste von Features, unverändert). Eine Aufzählung aller Katalog-Features außer der Liste in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) gibt es nicht.

### Beleg-Tabelle (Bedingung → Messaussage oder Test)

Messaussagen im [Messprotokoll](../../../messungen/slice-v1-forge-cli-quellen-messung/README.md) (Kennungen wie im Entwurf); „Code-Slice“ ist `slice-v1-devcontainer-feature-gh-glab` (Welle §4; die Datei besteht noch nicht, die Aufträge stehen deshalb auch in der Welle-Datei §5). Die Spalte *Test* nennt den geplanten Test als Auftrag, nicht als vorhandenen.

| # | Bedingung im Lastenheft `0.22.0` | Messaussage | Test bzw. Auftrag an den Code-Slice |
|---|---|---|---|
| B1 | `gh` und `glab` sind **in beiden Profilen** (Standard und Sandbox) als eingebaute Features wählbar | C10 (Blöcke als Handvariante im Sandbox-Dockerfile gebaut) | Golden Cases oder Render-Test: je Profil ein erzeugtes Dockerfile mit dem Block von `gh` und `glab`; im Standardprofil ohne UID-Überschreibung steht `USER root` nicht ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)) |
| B2 | Konfigurationsschlüssel `devcontainer.features.gh.enabled`, analog `glab` | — (Zielbild) | Test des Katalogs (`featureCatalogue`): beide Schlüssel sind gültig und aktivieren den Eintrag |
| B3 | Gepinntes Release-Archiv des Upstream-Projekts mit geprüfter SHA-256-Prüfsumme | C10 (Prüfzeile `…tar.gz: OK`), C11 (verfälschter Pin bricht den Bau ab, kein Image) | Golden Case des Blocks (Pin und Prüfzeile) und Docker-Test des Baus; `arm64` ist nicht gemessen und im Satz nicht zugesagt |
| B4 | Keine Signatur der Archive gefunden; `gh`-Attestierung wird von `u-boot` nicht geprüft | A10 (`gh`: keine Signatur-Assets, zwei Attestierungen über die API), A12 (`glab`: keine Signatur-/Provenance-Assets nach Suchmuster) | Aussage über das Suchergebnis, kein Verhalten des Codes; belegt durch das Messprotokoll |
| B5 | `u-boot` weist das beim Erzeugen aus | — (Zielbild) | Test der Befehlsausgabe von `generate`/`init` für beide Features (Ausweisung der fehlenden Signatur) |
| B6 | Hosts von `gh` stehen in der Default-Allowlist, für `glab` kommt ein Host hinzu | D1 (`gh`: beide Hosts in der Basis-Liste, B6 des Protokolls), D2 und D3 (`glab`: ein Host, nicht in der Basis-Liste; ohne Eintrag scheitern die Aufrufe, mit Eintrag gelingen sie) | Test von `egressFeatureHosts`: `glab` trägt den Host, `gh` fügt der Basis-Liste nichts hinzu; kein Hostname im Lastenheft |
| B7 | `u-boot` führt weder `gh auth login` noch `glab auth login` aus und speichert kein Token | E4 (Platzhalterwert steht in keiner Datei unter `~/.config`, `~/.cache`, `~/.local`, `~/.gitconfig`), E2 (Umgebungsvariable genügt als Quelle) | Test: der erzeugte Block und die Templates enthalten keinen Aufruf von `auth login` und keinen Token-Wert |
| B8 | Nicht zugesagt: Aufrufe mit gültigem Token und deren Hosts, Telemetrie-/Update-Hosts, selbst gehostete Instanzen, andere Versionen, `glab duo`, `glab orbit`, Nachladen weiterer Programme | F, D5, D7 (Grenzen des Messprotokolls; Ausschluss ist die Entscheidung des Projektinhabers) | Keiner nötig: eine Nicht-Zusage hat kein Verhalten, das ein Test pinnt |
| B9 | Das Token von `gh`/`glab` steht weder im Image noch im Workspace-Volume noch in `u-boot.yaml` oder einer erzeugten Datei; die Konfiguration trägt nur die Referenz | E4, B5 des Protokolls (die erzeugte `devcontainer.json` trägt `GIT_TOKEN` als Referenz, keinen Wert) | Test: `remoteEnv` trägt `${localEnv:GH_TOKEN}` bzw. `${localEnv:GITLAB_TOKEN}` und keinen Wert, wenn das Feature aktiv ist; nichts bei inaktivem Feature |
| B10 | Die Dokumentation nennt `GH_TOKEN` für `gh` und `GITLAB_TOKEN` für `glab` | E1 (Hilfe nennt beide an erster Stelle der Rangfolge), E2 (`401` mit gesetzter Variable ohne `auth login`) | Handbuch im Code-Slice; Test der Variablennamen über `agentCredentialEnv` (B9) |
| B11 | `u-boot doctor` prüft die Quelle der Token-Variable bei aktivem `gh`- oder `glab`-Feature: fehlt sie oder liegt sie im Klartext in einer Projektdatei, `warn` (kein Abbruch) | — (Zielbild; Muster `checkSandboxCredentials`, `doctor_sandbox.go`) | Test von `checkSandboxCredentials` je Feature: fehlt die Variable, `warn`; Klartext in einer Projektdatei, `warn`; Feature inaktiv, kein Befund |

Gegenprobe je Wächter (Test des Code-Slice, hier nur der Auftrag): B1 wird rot, wenn der Block im Standardprofil nicht gerendert wird; B3 rot, wenn die Prüfzeile aus dem Block fällt; B6 rot, wenn der Host von `glab` aus der Tabelle fällt; B9 rot, wenn die Referenz in `remoteEnv` fehlt oder ein Wert dort steht; B11 rot, wenn der Eintrag aus `agentCredentialEnv` fehlt. Sie wurden in diesem Slice nicht gesehen (kein Code); der Code-Slice liefert sie.

### Plan-Abgleich (§3.10)

Die Plan-Tabelle in §3 gegen `git diff --stat dd4d5e1a..HEAD` (der Stand nach dem Start, vor der Fix-Runde): `spec/lastenheft.md` (Wortlaut, Versionsfeld, Historie), `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/open/README.md`, `docs/plan/planning/open/slice-cr-lastenheft-messstaende.md`, diese Datei (die zwei Moves durch `make slice-mv` eingeschlossen), `docs/plan/planning/welle-forge-cli-gh-glab.md` (Verweis-Nachzug durch `make slice-mv` und der Auftrag in §5). Keine Zeile ohne Plan-Eintrag und kein Eintrag ohne Diff, mit zwei Ausnahmen, die erst die Reviewer-Rolle schreibt (`docs/reviews/`, `docs/plan/planning/observations/`).

### Verifikation (Implementer, 2026-10-23)

Commit-Reihenfolge (`git log`): `3b2179c7` (Freigabe in §7, `Verantwortlich:`, Roadmap) vor `5b55250f` (Lastenheft im freigegebenen Wortlaut); danach Versionsfeld, Historie-Zeile und die Nachzüge der Versionsangaben.

| Lauf | Exit-Code |
|---|---|
| `make gates` (nach dem Wortlaut-Commit und den Nachzügen) | 0 |
| `make doc-trace` | 0 (84 Anforderungen, 0 Waisen) |
| `make doc-planning` | 0 |
| `make test-docker`, `make verify-depguard` | nicht gelaufen: kein Code, kein Docker-Pfad, keine Importregel berührt |

Ein erster `make gates`-Lauf endete mit Exit-Code 2 (`docs-check`: eine nicht verlinkte Kennung in der Vorab-Suche dieses Abschnitts); die Stelle ist berichtigt, der Lauf oben ist der zweite.

*Offen für Review und Verification:* der Satz-für-Satz-Vergleich Lastenheft gegen Freigabe, die Hostnamen-Umsetzung (Punkt 1), Was hat funktioniert, Was ging anders, Steering-Loop-Eintrag, Beobachtungs-Register, Folge-Slices und die Ausgänge der Risiken aus §6.

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
