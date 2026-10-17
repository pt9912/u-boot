# Slice v1: Devcontainer-Feature `claude-code` im eingebauten Katalog

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md) — dritter Slice der Welle; ihm gehen [`slice-cr-devcontainer-agent-features-sandbox`](slice-cr-devcontainer-agent-features-sandbox.md) und [`slice-v1-devcontainer-feature-sshd`](slice-v1-devcontainer-feature-sshd.md) voraus, ihm folgt [`slice-v2-spec-devcontainer-sshd-claude-code`](slice-v2-spec-devcontainer-sshd-claude-code.md).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Devcontainer-Features), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Installationsskript eines Drittanbieters), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Hosts), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten, sinngemäß); [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) (Drittanbieter-Features im Katalog). Muster: [`slice-v1-devcontainer-features-homebrew-sdkman`](../done/slice-v1-devcontainer-features-homebrew-sdkman.md).

**Berührte Spec-Stellen:** — (keine). [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) und [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) bleiben bis [`slice-v2-spec-devcontainer-sshd-claude-code`](slice-v2-spec-devcontainer-sshd-claude-code.md) ohne den Schlüssel; das Lastenheft ändert nur der CR.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Schlüssel `claude-code` ist ein eingebautes Devcontainer-Feature und rendert das Feature von Anthropic (`ghcr.io/anthropics/devcontainer-features/claude-code`, als
Drittanbieter-Feature ausgewiesen, soweit der CR es so einstuft); die Egress-Allowlist trägt die gemessenen Hosts, wenn es aktiv ist, und ein Zugangsschlüssel des Agenten steht
in keiner erzeugten Datei, sondern wird zur Laufzeit übergeben.

**Planungsgrundlage (WebFetch 2026-10-17, im Slice gegen `install.sh` und `devcontainer-feature.json` erneut zu prüfen, nicht zu übernehmen):** Das Feature liegt **nicht**
unter `ghcr.io/devcontainers/features`, sondern unter dem Namensraum von Anthropic; die README nennt den Tag `:1.0` (ob `:1` existiert, ist unbestätigt, die Standardversion folgt der
Messung); es hängt von Node.js ab und installiert es bei Bedarf. Optionen und Netz-Hosts sind **nicht belegt**: erwartet wird die npm-Registry beim Build und `api.anthropic.com` zur
Laufzeit; weitere Hosts (Anmeldung, Telemetrie, Update-Prüfung) und der Authentifizierungsweg sind zu messen, nicht zu raten.

**Gefundener Stand im Code (aus dem Quelltext gelesen):**

- Drittanbieter-Muster: `featureCatalogueEntry.publisher` (`devcontainer_features.go`) kennzeichnet `homebrew` (`meaningful-ooo`); `thirdPartyFeatureWarnings` erzeugt je aktiviertem Eintrag mit
  `publisher` eine Warnung mit Code [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Herausgeber, Quelle, Version, Hinweis auf das Installationsskript beim Image-Build; ohne `SecurityRelaxation`). Ein Eintrag mit
  `source:`-Überschreibung löst sie nicht aus. [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) verlangt, dass ein Katalogeintrag außerhalb des offiziellen Namensraums seinen Herausgeber als Registry-Namensraum trägt (`anthropics` erfüllt das).
- Der Schlüsselname `claude-code` erfüllt das Muster von `domain.NewFeatureName` (`^[a-z]([a-z0-9-]{0,30}[a-z0-9])?$`).
- Egress: `egressFeatureHosts()` (`devcontainer_egress.go`); das Handbuch sagt heute, `u-boot` kenne die Hosts des Agenten nicht (Abschnitt Egress-Restriktion, Punkt „Wichtig“).
- Zugangsdaten: `devcontainer.json.tmpl` reicht im Sandbox-Zweig `GIT_TOKEN` per `remoteEnv` als `${localEnv:GIT_TOKEN}` durch; `credentialPatterns()` und `hasTokenSource` in `doctor_sandbox.go` kennen nur Git.
  Neue `doctor`-Check-Codes sind an `coderegistry.go` (`internal/adapter/driving/cli/jsontestutil`) und [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) gebunden (Drift-Test `drift_test.go`).

**Anforderung an den Slice — Zugangsdaten des Agenten ([`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) sinngemäß):** Ein Anthropic-Zugangsschlüssel steht **weder im Image noch im Workspace-Volume noch in `u-boot.yaml` oder
`devcontainer.json` als Wert**; zur Laufzeit gelangt er per Umgebungsvariable aus der Host-Umgebung in den Container (Referenz `${localEnv:<NAME>}`, nie der Wert; den Namen der Variable belegt die
Dokumentation des Features). Ob die Regel eine Zusage im Lastenheft wird, entscheidet der CR (Frage 4); die Umsetzung gilt unabhängig davon.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Feature-Optionen** — der Katalog kennt keine Optionen ([`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md), Re-Evaluierungs-Trigger); der Nutzer überschreibt den Tag über `devcontainer.features.claude-code.version`.
- **Konfiguration, Start oder Berechtigungsmodus des Agenten, `~/.claude`-Mounts** — [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) nennt es als Nicht-Ziel; ein Mount von Host-Dateien widerspricht dem Profil (kein Bind-Mount, keine Host-Geheimnisse).
- **Ein neuer `doctor`-Check-Code** — er wäre an Registry und [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) gebunden und liefe der Regel „Spezifikation nach dem Code“ zuwider; die vorhandene Prüfung (`devcontainer.sandbox.credentials`) wird erweitert, ohne neue Kennung.
- **Weitere Agenten-Features** (Codex, Gemini, andere) — jede Quelle braucht eine eigene Entscheidung.
- **Die Spezifikation** ([`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)) — sie folgt in [`slice-v2-spec-devcontainer-sshd-claude-code`](slice-v2-spec-devcontainer-sshd-claude-code.md).

## 2. Definition of Done

- [ ] **Katalog und Generator:** `claude-code` im Katalog (Quelle `ghcr.io/anthropics/devcontainer-features/claude-code`, Herausgeber `anthropics`, Standardversion nach der Messung); die Drittanbieter-Warnung folgt der Entscheidung des CR (Frage 3); Egress-Hosts nur bei aktivem Feature und aus der Messung; das Sandbox-Verhalten folgt dem CR (Frage 2); Golden Cases belegen Render, Warnung, Hosts und double-run -> NoOp, **ein Golden Case aktiviert `claude-code` zusammen mit `sshd`** (Welle-Closure).
- [ ] **Zugangsdaten zur Laufzeit:** `devcontainer.json` enthält bei aktivem Feature die Referenz der Umgebungsvariable (kein Wert); ohne das Feature bleibt die Ausgabe byte-identisch; die vorhandene `doctor`-Prüfung erkennt einen Klartext-Schlüssel in den gescannten Dateien (Muster aus der Dokumentation des Schlüsselformats belegt, keine neue Check-Kennung); Tests je Fall.
- [ ] **Messung und Doku:** Die Dev-Container-CLI baut ein Image mit `claude-code` auf dem Basis-Image samt Node-Abhängigkeit (Verfahren der Messreihe in [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md)); belegt sind: der Tag (`1` oder `1.0`), die Feature-Metadaten (`mounts`, `capAdd`, `privileged`, `entrypoint`, `installsAfter`, Optionen), die Laufzeit-Hosts (DNS-Abfragen) und der Anmeldeweg; eine Folge-ADR (nächste freie Nummer) nur, wenn Messung oder CR etwas entscheiden, das [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) nicht trägt (andere Einstufung, Optionsmodell). Benutzerhandbuch (Katalogtabelle, Egress-Tabelle, Abschnitt Zugangsdaten des Agenten; der Satz „Die Hosts des Agenten selbst kennt `u-boot` nicht“ ist angepasst), `README.md`, `README.de.md`, `CHANGELOG.md` (`[Unreleased]`) und `harness/replay.md` nennen den Schlüssel.
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
| `internal/hexagon/application/devcontainer_features.go` | update | Eintrag `claude-code` in `featureCatalogue()` (Quelle, Version, `publisher: anthropics`, `shortDesc`) |
| `internal/hexagon/application/devcontainer_egress.go` | update | Zeile in `egressFeatureHosts()` aus der Messung |
| `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl`, `internal/hexagon/application/templates.go` | update | `remoteEnv` trägt die Referenz des Agenten-Schlüssels bei aktivem Feature; `GIT_TOKEN` im Sandbox-Zweig bleibt byte-identisch |
| `internal/hexagon/application/generate.go`, `devcontainer_sandbox.go` | update (bedingt) | Drittanbieter-Warnung (`thirdPartyFeatureWarnings`) und Sandbox-Regel laut CR |
| `internal/hexagon/application/doctor_sandbox.go` | update | `credentialPatterns()` um das Schlüsselformat erweitert, Meldungstext nennt nicht nur Git; **keine** neue Check-Kennung |
| `internal/hexagon/application/devcontainer_features_test.go`, `generate_features_test.go`, `generate_sandbox_test.go`, `doctor_sandbox_test.go`, `doctor_features_drift_test.go` | update | Katalogumfang (Herausgeber als Namensraum), Render, Warnung, Hosts, Zugangsdaten-Referenz, Klartext-Erkennung, beide Features zusammen — nach [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | update | öffentlicher Vertrag und Golden-Zeile |
| `docs/plan/adr/` (Folge-ADR, Index) | neu / update (bedingt) | nur wenn die Messung oder der CR es verlangt (DoD 3) |

Die Pfade stammen aus dem Quelltext (Lesestand 2026-10-17); der Implementer prüft jede Zeile vor der Änderung gegen den Code und trägt Abweichungen in die Evidence ein.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der CR [`slice-cr-devcontainer-agent-features-sandbox`](slice-cr-devcontainer-agent-features-sandbox.md) und der Slice
[`slice-v1-devcontainer-feature-sshd`](slice-v1-devcontainer-feature-sshd.md) liegen in `done/` (beide berühren Katalog, Egress und Tests; WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice überschreitet 3 Liefer-Punkte, etwa weil die Zugangsdaten-Übergabe einen eigenen Slice verlangt.
- `in-progress` → `open` (blockiert): die Messung findet eine Lockerung des Profils, die der CR nicht trägt (dann zuerst CR erweitern); das Feature von Anthropic baut nicht oder ist nicht mehr unter diesem Namen verfügbar.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Das Feature von Anthropic ändert sich oder verschwindet** (Drittanbieter, Tag-Pin): ein Build-Fehler fällt beim Image-Build auf. — **Ausgang:** weiter offen: Re-Evaluierungs-Trigger von [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) (Änderung oder Wegfall des Features des Herausgebers).
- **Die Egress-Hosts reichen nicht** (Anmeldung, Telemetrie, Update-Prüfung des Agenten sind nicht belegt): der Agent läuft im Sandbox-Profil mit Egress-Restriktion nicht. — **Ausgang:** entfallen, sobald die Messung (DoD 3) die Hosts belegt; weitere Hosts trägt der Nutzer über `devcontainer.sandbox.egress.allow` ein, das Handbuch nennt es (bei Closure bestätigen).
- **Das Feature bringt Metadaten mit** (`mounts`, `capAdd`, `privileged`, `entrypoint`) und unterläuft damit das Profil (kein Host-Bind-Mount, keine Host-Geheimnisse). — **Ausgang:** weiter offen: die Messung entscheidet; bei Eintritt Rückführung an den CR (§4).
- **Der Anmeldeweg des Agenten legt Zugangsdaten im Container ab** (Anmeldung im Container statt Umgebungsvariable): sie lägen im Home des Container-Benutzers, nicht im Image oder Volume des Profils. — **Ausgang:** weiter offen bis zur Entscheidung des Projektinhabers (Frage 4 im CR); das Handbuch nennt den gemessenen Weg.
- **Schlüsselformat und Variablenname sind geraten**, wenn die Dokumentation sie nicht nennt. — **Ausgang:** entfallen: DoD 2 verlangt den Beleg aus der Dokumentation; ohne Beleg wird das Muster gestrichen und der Slice meldet die Lücke.
- **Plan-Aussage über den Code stimmt nicht** (siehe §8): die Fundstellen in §1 und §3 sind ein Lesestand. — **Ausgang:** weiter offen: der Implementer prüft jede Aussage vor der Änderung im Code; Abweichungen stehen in der Evidence, die Klasse zählt die Closure.

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `internal/hexagon/application` (Kürzel APP) und, für die Doku, `docs/user/`, `harness/`, `docs/plan/` und `spec/` (DOK) — beide Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006) (Abschnitt Modus-Deklaration), beide erfüllen die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. `internal/adapter/driving/cli` (CLI) wird nur berührt, falls ein neuer Warn-Code in `coderegistry.go` stehen muss; dann gilt `AGENTS.md §3.9` (Kommentar-Selbstcheck).

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer für APP und DOK. Die einzige Beobachtung `BEO-CLI/kommentar-beschreibt-abwesendes` (3×, verkörpert in `AGENTS.md §3.9`) betrifft CLI; sie wird nur bei einer CLI-Berührung relevant. Die Klasse „Plan-Aussage über Code stimmt nicht“ steht nicht im Register (siehe [`slice-v1-devcontainer-feature-sshd`](slice-v1-devcontainer-feature-sshd.md) §8); ihre Belege legt die Slice-Closure an.

Alle berührten Sub-Areas GF.
