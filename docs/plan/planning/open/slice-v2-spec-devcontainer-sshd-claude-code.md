# Slice v2: Spezifikation zu `sshd` und `claude-code` (Katalog, Egress, Sandbox-Verhalten)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md) — letzter Slice der Welle; ihm gehen [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md), [`slice-v1-cli-quiet-sicherheitswarnungen`](../done/slice-v1-cli-quiet-sicherheitswarnungen.md), [`slice-v1-devcontainer-feature-sshd`](../done/slice-v1-devcontainer-feature-sshd.md) und [`slice-v1-devcontainer-feature-claude-code`](../in-progress/slice-v1-devcontainer-feature-claude-code.md) voraus.

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Hosts), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten); [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md). Muster: [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../done/slice-v1-cli-verbosity-letzte-option-gewinnt.md) (Spezifikations-Eintrag nach allen anderen Slices).

**Berührte Spec-Stellen:** `spezifikation.md`: [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (eingebaute Features), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Egress-Hosts je Feature), [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) (Beschreibung der Zeile `devcontainer.sandbox.credentials`), [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Wirkung von `--quiet`: Ausnahme 1 und Zeile zu `doctor`), neuer Eintrag ab `SPEC-044` und §7 (Historie). `lastenheft.md` und `architecture.md` berührt der Slice nicht. <!-- d-check:ignore (SPEC-044 ist noch nicht vergeben) -->

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Spezifikation beschreibt Katalog, Egress-Hosts und Sandbox-Verhalten der Features `sshd` und `claude-code` deckungsgleich mit Lastenheft (nach dem CR) und Code: die
Katalogzeilen in [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), die Host-Zeilen in [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) und ein neuer Eintrag ab `SPEC-044` für die Sandbox-Regel beider Features, die Übergabe der Zugangsdaten des Agenten und die Warnungen; außerdem zieht der Slice [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) nach: Ausnahme 1 nennt den Login-Dienst und die Warn-Einträge von `u-boot doctor`, die eine Sicherheitseigenschaft des Sandbox-Profils betreffen (Fassung der Ausnahme 1 in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)). <!-- d-check:ignore (SPEC-044 ist noch nicht vergeben) -->

**Einstufung der `doctor`-Einträge unter `--quiet` ([`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)):** Die Spezifikation trägt die Einstufung des Quiet-Slice ([`slice-v1-cli-quiet-sicherheitswarnungen`](../done/slice-v1-cli-quiet-sicherheitswarnungen.md) §1) gegen den Code gelesen, darunter zwei Aussagen, die der Wortlaut der Ausnahme 1 nicht selbst trägt: (i) `nestedRuntime: docker` hat im `doctor` die Stufe OK (kein Warn-Eintrag; `--privileged` weist die Befehlsausgabe aus), die Aufzählung „Sicherheitslockerungen des Docker-Daemons“ nennt also keinen `doctor`-Eintrag; (ii) die Warn-Einträge der Feature-Prüfungen `devcontainer.features.allowlist` und `devcontainer.features.drift` sind keine Sicherheitseigenschaft des Sandbox-Profils im Sinn der Ausnahme 1 und entfallen unter `--quiet`. Dazu die Entscheidung des Projektinhabers vom 2026-10-17: Egress „has no effect“ bleibt unter `--quiet` sichtbar, `nestedRuntime` „has no effect“ entfällt.

**Übergabe aus `slice-v1-devcontainer-feature-claude-code` (Verification, 2026-10-17):** (i) [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) trägt für `claude-code` genau zwei Hosts, `registry.npmjs.org` und `api.anthropic.com` (Fundstelle `egressFeatureHosts` in `devcontainer_egress.go`; Entscheidung des Projektinhabers 2026-10-17, Anmelde-Hosts trägt der Nutzer über `egress.allow`); (ii) die Zeile `devcontainer.sandbox.credentials` in [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) nennt heute nur Klartext-Git-Zugangsdaten und Token-Quelle, der Check prüft auch Agenten-Schlüssel und die Referenz `${localEnv:ANTHROPIC_API_KEY}` (Fundstelle `checkSandboxCredentials` in `doctor_sandbox.go`; kein neuer Check-Code, daher keine neue Zeile, nur die Beschreibung); (iii) [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) trägt den Schlüssel `claude-code` samt Herausgeber `anthropics`.

**Lehre der vorigen Welle (`welle-cli-verbosity`):** Der Entwurf wird gegen den **Code** gelesen, nicht gegen die Plan-Prosa der Slices. Jede Aussage trägt in der Evidence eine Fundstelle (Datei und Funktion
oder Test); eine Aussage ohne Fundstelle wird gestrichen. Zahlenwerte und Hosts stammen aus dem Quelltext; weicht der Code vom Lastenheft ab, wird die Abweichung vermerkt und nicht in der Spezifikation geglättet.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Änderungen an Code, Tests, Handbuch oder CHANGELOG** — sie liegen in den Code-Slices ([`slice-v1-devcontainer-feature-sshd`](../done/slice-v1-devcontainer-feature-sshd.md),
  [`slice-v1-devcontainer-feature-claude-code`](../in-progress/slice-v1-devcontainer-feature-claude-code.md)); ein abweichender Wert wird vermerkt, nicht geändert (Folge-Slice).
- **Änderungen am Lastenheft** — sie gehören dem CR [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md); die Spezifikation verfeinert und erweitert das Lastenheft nie.
- **`architecture.md`** — die Sicht nennt keine einzelnen Features und, nach `AGENTS.md` §3.4, weder Slices noch ADRs.
- **Neue `doctor`-Check-Codes in [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry)** — kein Code-Slice der Welle führt einen ein; entsteht dennoch einer, gehört die Zeile in diesen Slice. Die Beschreibung der vorhandenen Zeile `devcontainer.sandbox.credentials` zieht der Slice nach (siehe Übergabe oben).

## 2. Definition of Done

- [ ] **Spezifikation:** [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) nennt die Schlüssel `sshd` (offizielle Quelle, Standardversion) und `claude-code` (Quelle, Herausgeber, Standardversion, Drittanbieter-Einstufung); [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) trägt die Host-Zeilen beider Features, soweit der Code welche führt (`claude-code`: genau `registry.npmjs.org` und `api.anthropic.com`); die Zeile `devcontainer.sandbox.credentials` in [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) nennt auch Agenten-Schlüssel und die Referenz bei aktivem Feature; der neue Eintrag ab `SPEC-044` beschreibt Sandbox-Regel, Zugangsdaten-Übergabe und Warnungen; [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) trägt in Ausnahme 1 den Login-Dienst und die `doctor`-Sicherheits-Einträge, und die Zeile zu `doctor` unter `--quiet` nennt neben Fehler-Einträgen und Degradations-Fallbacks auch sie; `SPEC-043` hält fest, dass `nestedRuntime: docker` die Stufe OK hat (kein Warn-Eintrag), dass die Feature-Allowlist- und Feature-Drift-Warnungen unter `--quiet` entfallen und dass Egress „has no effect“ sichtbar bleibt, `nestedRuntime` „has no effect“ dagegen entfällt (Fundstellen im Code: `writeDoctorReport` in `doctor.go`, Kennzeichen am `domain.Diagnostic`, `sandboxRelaxationWarnings` in `devcontainer_sandbox.go`); §7 (Historie) trägt die Zeile. <!-- d-check:ignore (SPEC-044 ist noch nicht vergeben) -->
- [ ] **Gegen den Code gelesen:** Die Evidence führt je Aussage die Fundstelle; die Spezifikation entstand nach dem Move der vier anderen Slices nach `done/` (Commit-Reihenfolge belegt es).
- [ ] **Traceability:** `make doc-trace` meldet keine Waise; `make docs-check` meldet 0 Befunde (jede neue `SPEC-`Kennung ist ein auflösbarer Anker).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update, falls ein öffentlicher Vertrag berührt ist.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft, im Repo mit Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` ([`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) (Beschreibung), neuer Eintrag ab `SPEC-044`, §7 Historie) | update | Katalog, Hosts, Sandbox-Regel, Wirkung von `--quiet`; Platzierung des neuen Eintrags im Kapitel zu Devcontainer, die nächste freie Nummer ist beim Schreiben zu prüfen (Stand 2026-10-17: höchste vergebene Nummer [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)); [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) selbst ändert sich in Ausnahme 1 und der `doctor`-Zeile <!-- d-check:ignore (SPEC-044 ist noch nicht vergeben) --> |
| `internal/hexagon/application/devcontainer_features.go`, `devcontainer_egress.go`, `devcontainer_sandbox.go`, `generate.go`, `doctor_sandbox.go`, `internal/adapter/driving/cli/doctor.go`, `templates/devcontainer/devcontainer.json.tmpl` | lesen (Quelle der Fundstellen) | die Spezifikation wird gegen diese Dateien gelesen; keine Änderung |

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; die vier Slices [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md),
[`slice-v1-cli-quiet-sicherheitswarnungen`](../done/slice-v1-cli-quiet-sicherheitswarnungen.md), [`slice-v1-devcontainer-feature-sshd`](../done/slice-v1-devcontainer-feature-sshd.md) und [`slice-v1-devcontainer-feature-claude-code`](../in-progress/slice-v1-devcontainer-feature-claude-code.md) liegen in `done/`;
kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice überschreitet 3 Liefer-Punkte.
- `in-progress` → `open` (blockiert): Lastenheft und Code widersprechen sich in einem Punkt; der Widerspruch geht als Folge-Slice (CR oder Code) vor, bevor die Spezifikation ihn beschreibt.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). Mit diesem Slice ist die Welle
[`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md) schließbar.

## 6. Risiken und offene Punkte

- **Der Entwurf folgt der Plan-Prosa statt dem Code** (Klasse „Plan-Aussage über Code stimmt nicht“, Beobachtung [BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), verkörpert in `AGENTS.md §3.10`): der Reviewer fand in `welle-harness-wartung` so einen Fehler in einem Eintrag, den die Inventur nicht fand. — **Ausgang:** entfallen: DoD 2 verlangt je Aussage eine Fundstelle im Code; bei Abweichung gilt der Code, die Abweichung steht in der Evidence (bei Closure bestätigen).
- **Lastenheft und Code widersprechen sich** (etwa eine Warnung, die der CR anders beschreibt als der Code sie setzt). — **Ausgang:** weiter offen: gemeldet als Folge-Slice (CR oder Code), nicht in der Spezifikation geglättet (§4, Rückführung).
- **Die Abfolge „Spezifikation zuletzt“ lässt [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) und [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) zwischen den Code-Slices und diesem Slice hinter dem Code zurück.** — **Ausgang:** entfallen: Kein Gate prüft die Tabellen gegen den Code (nur die `doctor`-Code-Registry ist durch einen Drift-Test gebunden); die Lücke ist in der Welle benannt (Closure-Trigger) und endet mit diesem Slice.

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `spec/` (Kürzel DOK) — Modus Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006) (Abschnitt Modus-Deklaration), erfüllt die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. Der Code wird nur gelesen.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer für diese Sub-Area (die einzige Beobachtung `BEO-CLI/kommentar-beschreibt-abwesendes` betrifft CLI und ist verkörpert in `AGENTS.md §3.9`). Die Klasse „Plan-Aussage über Code stimmt nicht“ steht als [BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) im Register (verkörpert in `AGENTS.md §3.10`); dieser Slice ist ihr Gegenmittel (DoD 2).

Alle berührten Sub-Areas GF.
