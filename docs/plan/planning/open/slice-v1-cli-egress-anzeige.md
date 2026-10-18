# Slice v1: Befehl für die wirksame Egress-Liste mit Herkunft je Host (`LH-FA-DEV-008`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) — dritter Slice der Welle; ihm gehen [`slice-v1-egress-wildcard-pruefung`](../done/slice-v1-egress-wildcard-pruefung.md) und [`slice-cr-sandbox-egress-anzeige`](slice-cr-sandbox-egress-anzeige.md) voraus; die Spezifikation folgt in [`slice-v2-spec-sandbox-egress-ollama`](slice-v2-spec-sandbox-egress-ollama.md).

**Bezug:** [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Restriktion; um die Abfrage erweitert durch den CR, Lastenheft 0.14.0), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Abfragebefehl behält unter `--quiet` seine Ausgabe), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil als Bedingung). Verwandt: [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md), [`ADR-0013`](../../adr/0013-dokumentationsreferenzmodell.md).

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Quellen der Hosts), [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Aufzählung der Abfragebefehle), `architecture.md` §2 (Driving-Port, neuer Use Case) — der Slice **liest** sie; die Nachführung nimmt [`slice-v2-spec-sandbox-egress-ollama`](slice-v2-spec-sandbox-egress-ollama.md) nach dem Code vor („Spezifikation nach dem Code“). Das Lastenheft ändert nur der CR.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Befehl (Arbeitsname `u-boot egress list`, endgültig nach CR-Freigabe) zeigt die wirksame Egress-Allowlist und je Host seine Herkunft (Basis, Container-Runtime, Docker-APT-Host, Feature mit Namen, Clone-Quelle, Nutzerliste) — **offline, ohne Docker, ohne laufenden Container, ohne Schreibzugriff**, mit `--json`. Die angezeigte Liste ist **dieselbe Berechnung**, die das Startscript bekommt (eine Funktion, kein zweiter Algorithmus), und ein Test beweist die Gleichheit.

### Gefundener Stand im Code (gelesen am 2026-10-18)

- **Berechnung:** `egressHosts(dc, cloneURL)` in `internal/hexagon/application/devcontainer_egress.go` setzt die Vereinigung aus Basis (`egressBaseHosts`), Registries (`egressRegistryHosts`, bei `nestedPodman(dc) || nestedDocker(dc)`), `egressDockerAptHost` (bei `nestedDocker`), `egressFeatureHosts()` je Feature mit `Enabled` wahr, `cloneSourceHost(cloneURL)` und `dc.Sandbox.Egress.Allow`; sie liefert `[]string` sortiert und ohne Duplikate und **verwirft die Herkunft**. Der einzige Aufrufer außerhalb der Tests ist `devcontainerTemplateData` (`internal/hexagon/application/devcontainer_sandbox.go`, `data.EgressHosts = egressHosts(dc, cloneURL)`); die Clone-URL kommt aus `resolveSandboxClone` bzw. `effectiveCloneSource` (dort `dc.Sandbox.Repository`, sonst `readGitOriginURL`, das `.git/config` über den `FileSystem`-Port liest).
- **Aktivierung:** `egressEnabled(dc)` (gleiche Datei) ist wahr bei `devcontainer.sandbox.egress.enabled: true`; `profileIsSandbox(dc)` liest das Profil aus `u-boot.yaml` (der Parameter `sandbox` von `devcontainerTemplateData` kommt von `generate.go` bzw. `initproject.go`; sein Zustandekommen, auch über `--sandbox`, liest der Architect dort). `generate devcontainer` rendert die Liste nur, wenn `sandbox` und `egressEnabled(dc)` gelten (`devcontainerTemplateData` kehrt bei `!sandbox` vor dem Egress-Zweig zurück); `u-boot doctor` weist den Fall „Egress ohne Sandbox-Profil“ als `warn` aus (`checkSandboxEgress`, `internal/hexagon/application/doctor_sandbox.go`).
- **Vorbild für einen offline lesenden Use Case:** `DoctorService.loadUbootYAML` (`internal/hexagon/application/doctor.go`) liest `u-boot.yaml` über `FileSystem` und `YAMLCodec`; `ConfigUseCase.List` (`internal/hexagon/port/driving/config.go`) und `runConfigList` (`internal/adapter/driving/cli/config.go`) zeigen die Form eines Abfragebefehls mit `--json` (Hülle `newDataEnvelope`, `writeEnvelope`) und die Ablehnung von `--dry-run`/`--diff` (`ErrDryRunNotApplicable`, `registerConfigPreviewRejectFlags`).
- **Verdrahtung:** `cmd/uboot/main.go` baut die Services und übergibt sie `cli.New(version, initSvc, doctorSvc, …, logsSvc, cli.WithLogLevel(logLevel))`; `internal/adapter/driving/cli/root.go` registriert jeden Befehl mit einer `AddCommand`-Zeile (`newConfigCommand(a)` usw.). Ein neuer Use Case braucht deshalb einen Parameter oder eine Option an `cli.New`, eine Zeile in `main.go` und eine `AddCommand`-Zeile.

### Entwurfsfragen (der Architect entscheidet und begründet; Einordnung des Planners, nicht entscheidend)

1. **Name und Ort des Befehls:** (i) neues Top-Level-Kommando `egress` mit Unterbefehl `list` (ein `reload` entsteht nach der Entscheidung R0 nicht, siehe [`slice-v1-sandbox-egress-neuladen-entwurf`](../done/slice-v1-sandbox-egress-neuladen-entwurf.md)); (ii) `config egress` (neben `get`, `list`, `set`; aber die Liste ist abgeleitet, kein Konfigurationspfad); (iii) eine Option von `doctor` (aber `doctor` prüft, es listet nicht). Maßgebend ist der Wortlaut des CR; der Planner ordnet (i) als den kleinsten Eingriff in die Bedeutung der bestehenden Befehle ein.
2. **Zuschnitt der Berechnung:** `egressHosts` wird zu einer Funktion, die je Host die Liste der Quellen liefert (Typ im Paket `application`), und die bisherige `[]string`-Form wird aus ihr abgeleitet — nur so kann die Anzeige nicht von der Skript-Liste abweichen. Ob der Typ in `domain` gehört (reiner Datentyp, wie `domain.NewEgressHost`) oder im Paket bleibt, entscheidet der Architect nach `spec/architecture.md`.
3. **Herkunfts-Namen:** feste, kleine Menge (`base`, `runtime`, `docker-apt`, `feature:<name>`, `clone`, `user`) mit stabiler Schreibweise für `--json`; ein Host mit mehreren Quellen nennt alle (belegt: `ghcr.io`, `registry.npmjs.org`).
4. **Ohne aktive Restriktion:** zeigt der Befehl die Liste, die bei Aktivierung gälte, und sagt das in einer Zeile (so der Wortlaut des CR), oder er endet mit einem Hinweis ohne Liste; der Wortlaut des CR legt das fest.
5. **Clone-Quelle ohne Remote und ungültige URL:** `resolveSandboxClone` macht daraus beim Erzeugen eine Warnung bzw. einen Fachfehler (Exit-Code `10`); der Befehl darf weder abbrechen noch raten — er lässt die Quelle `clone` weg und nennt den Grund, oder gibt den Fachfehler wie `generate` zurück. Der Architect wählt und begründet.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Zustand des laufenden Containers** (griff die Restriktion beim Start?) — das leistet `doctor` (`checkSandboxFallbacks`); der Befehl bleibt offline, sonst bräuchte er den `DevcontainerInspector` und verlöre die Eigenschaft „ohne Docker“.
- **Neu-Laden der Liste im laufenden Container** — [`slice-v1-sandbox-egress-neuladen-entwurf`](../done/slice-v1-sandbox-egress-neuladen-entwurf.md), eigene Zusage und Sicherheitsentscheidung.
- **Ändern der Liste** (`egress add` o. ä.) — `u-boot config set devcontainer.sandbox.egress.allow …` (`setEgressAllow`) bleibt der Weg; ein zweiter Schreibweg wäre eine neue Zusage.
- **Wildcards in der Eingabe** — [`slice-v1-egress-wildcard-pruefung`](../done/slice-v1-egress-wildcard-pruefung.md).
- **Die Hosts von Ollama** — [`slice-v1-devcontainer-feature-ollama`](slice-v1-devcontainer-feature-ollama.md) trägt sie in `egressFeatureHosts()` ein; die Anzeige zeigt sie dann ohne weitere Änderung (belegt durch den Golden Case der Welle).

## 2. Definition of Done

- [ ] **Befehl und Berechnung:** Der Befehl zeigt je Host die Quellen (menschenlesbar und `--json`), offline, ohne Schreibzugriff, ohne `--dry-run`/`--diff` (abgelehnt wie bei `config list`); die Berechnung der Herkunft ist **dieselbe Funktion**, aus der `devcontainerTemplateData` die Liste des Startscripts ableitet; Tests belegen: die Liste des Befehls ist gleich der `hosts`-Zeile von `egress-init.sh` für dieselbe Konfiguration (Fälle: Standard, `podman`, `docker`, Features mit zwei Quellen für denselben Host, Clone-Quelle, Nutzerliste), Aktivierungszustand (aktiv; Egress an ohne Sandbox; Egress aus), ohne Remote, und `--quiet` behält die Ausgabe.
- [ ] **Vertrag nachgezogen:** `docs/user/benutzerhandbuch.md` (Abschnitt Egress und die Befehlsübersicht), `README.md` und `README.de.md` (falls sie Befehle aufzählen), `CHANGELOG.md` (`[Unreleased]`), `harness/replay.md` (Golden-Zeile) nennen den Befehl; `make doc-trace` meldet für die neue Aussage in [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) einen Beleg.
- [ ] **Hilfe-Ausgabe:** `u-boot --help` und die Hilfe des Befehls nennen Zweck und ein Beispiel ([`LH-FA-CLI-002`](../../../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe)); Test pinnt die Kurzbeschreibung.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update, falls ein öffentlicher Vertrag berührt ist (hier ja: neuer Befehl).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen von der Welle-Closure geprüft.

## 3. Plan (vor Code)

Die Pfade stammen aus dem Quelltext (Lesestand 2026-10-18); der Implementer prüft jede Zeile vor der Änderung gegen den Code und trägt Abweichungen in die Evidence ein. Die Zeilen für Befehlsort und Typ-Ort hängen an der Entscheidung des Architect (§1, Fragen 1 und 2).

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/application/devcontainer_egress.go` | update | `egressHosts` liefert je Host die Quellen; die `[]string`-Form wird daraus abgeleitet |
| `internal/hexagon/application/devcontainer_sandbox.go` | update (bedingt) | `devcontainerTemplateData` nutzt die abgeleitete Form; bleibt der Aufruf gleich, entfällt die Zeile |
| `internal/hexagon/application/egress_list.go` | neu | Use-Case-Implementierung: liest `u-boot.yaml` (`FileSystem`, `YAMLCodec`), Clone-Quelle über `effectiveCloneSource`; kein Docker-Port |
| `internal/hexagon/port/driving/egress.go` | neu | Driving-Port: Request/Response des Use Case |
| `internal/adapter/driving/cli/egress.go` | neu | Cobra-Befehl, menschenlesbare Ausgabe und `--json` (Hülle wie `config list`), Ablehnung von `--dry-run`/`--diff` |
| `internal/adapter/driving/cli/root.go`, `internal/adapter/driving/cli/cli.go` | update | `AddCommand`-Zeile; Use Case an `cli.New` (Parameter oder Option) |
| `cmd/uboot/main.go` | update | Service bauen und an `cli.New` übergeben |
| `internal/hexagon/application/devcontainer_egress_test.go` (neu) und `generate_sandbox_test.go` | neu / update | Herkunft je Host, Gleichheit mit `hosts=` des Startscripts, Mehrfach-Herkunft, Aktivierungszustände — nach [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) |
| `internal/hexagon/application/egress_list_test.go` | neu | Use Case: kein Remote, ungültige Clone-URL, keine `u-boot.yaml`, Egress aus |
| `internal/adapter/driving/cli/egress_test.go`, `quiet_test.go` | neu / update | Ausgabeform, `--json`-Hülle, `--quiet` behält die Ausgabe ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)), Ablehnung von `--dry-run`/`--diff`, Hilfetext |
| `internal/adapter/driving/cli/jsontestutil/coderegistry.go` | update (bedingt) | nur wenn der Befehl einen neuen Diagnose-Code einführt (`drift_test.go` verlangt dann [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry)); im Plan: kein neuer Code |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | update | öffentlicher Vertrag und Golden-Zeile |
| `docs/reviews/` (Review-Report, Verification Evidence) | neu | Review und Verifikation |
| diese Slice-Datei, `docs/plan/planning/observations/…/evidence/` | update | Umsetzungsstand, Closure (§7), Register |

Abgleich der Tabelle gegen `git diff --stat` nach der **letzten** Fix-Runde (`AGENTS.md` §3.10); jede Zeile mit „bedingt“ wird als eingetreten oder entfallen markiert.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; [`slice-cr-sandbox-egress-anzeige`](slice-cr-sandbox-egress-anzeige.md) liegt in `done/` (Lastenheft 0.14.0 mit freigegebenem Wortlaut und Befehlsnamen); [`slice-v1-egress-wildcard-pruefung`](../done/slice-v1-egress-wildcard-pruefung.md) liegt in `done/`; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice überschreitet 3 Liefer-Punkte, etwa weil die Umstellung der Berechnung auf Quellen je Host die Golden Cases von `generate devcontainer` ändert (ein geänderter Golden Case ist eine bewusste Produktentscheidung, `harness/replay.md`).
- `in-progress` → `open` (blockiert): der Wortlaut des CR verlangt eine Ausgabe, die nur mit Docker oder einem laufenden Container ermittelbar ist.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). Verifier und Validator-Hinweis: Die Gleichheit mit dem Startscript ist am Nutzerpfad nachzuspielen (`generate devcontainer` mit Egress, dann Befehl, Vergleich der Liste).

## 6. Risiken und offene Punkte

- **Die Anzeige weicht von der Skript-Liste ab**, weil zwei Berechnungen entstehen. — **Ausgang:** offen bis zur Closure; erwartet: entfallen, weil eine Funktion beide trägt und ein Test die Gleichheit prüft.
- **Die Umstellung von `egressHosts` ändert Golden Cases oder die Reihenfolge der Hosts im Startscript.** — **Ausgang:** offen bis zur Closure; erwartet: entfallen (sortierte, duplikatfreie Liste bleibt byte-identisch; Test).
- **Die Clone-Quelle ist offline nicht immer bestimmbar** (kein Remote, ungültige URL; Fragen 4 und 5). — **Ausgang:** offen; entschieden im Entwurf, belegt durch Test.
- **Ein Befehl mehr in der Oberfläche** erhöht die Pflege (Handbuch, Hilfe, `--json`-Hülle). — **Ausgang:** weiter offen als bewusste Zusage; kein Register-Eintrag.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: Entscheidung zu Befehlsort, Typ-Ort und Herkunftsnamen, Fundstellen je Aussage, Steering-Loop-Eintrag, Register, Folge-Slices, Risiken.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/application` (APP), `internal/hexagon/port` (PORT), `internal/adapter/driving/cli` (CLI), `cmd/uboot` (keine Sub-Area, Aspirantin; eine Verdrahtungszeile, keine Konvention), `docs/` (DOK) und `harness/` (DOK) — alle Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); jede erfüllt die Schwelle (zwei von drei Achsen) laut Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (beide `verkörpert` in `AGENTS.md` §3.10: Plan-Tabelle gegen `git diff --stat` nach der letzten Fix-Runde; jede Bedingung in Handbuch oder Spec hat einen Test), [`BEO-WZG/kopplung-nicht-nachgezogen`](../observations/BEO-WZG/kopplung-nicht-nachgezogen/observation.md) (offen, zwei Belege: hier berührt, weil Befehlsliste, Handbuch, README und Golden-Zeile zusammen nachzuziehen sind). Evidenz-/Diskrepanz-Risiko: niedrig bis mittel.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
