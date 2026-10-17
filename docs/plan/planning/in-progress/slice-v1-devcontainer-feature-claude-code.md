# Slice v1: Devcontainer-Feature `claude-code` im eingebauten Katalog

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md) — vierter Slice der Welle; ihm gehen [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md), [`slice-v1-cli-quiet-sicherheitswarnungen`](../done/slice-v1-cli-quiet-sicherheitswarnungen.md) und [`slice-v1-devcontainer-feature-sshd`](../done/slice-v1-devcontainer-feature-sshd.md) voraus, ihm folgt [`slice-v2-spec-devcontainer-sshd-claude-code`](../open/slice-v2-spec-devcontainer-sshd-claude-code.md).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Devcontainer-Features), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Installationsskript eines Drittanbieters), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Hosts), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten des Agenten, erweitert durch den CR); [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) (Drittanbieter-Features im Katalog). Muster: [`slice-v1-devcontainer-features-homebrew-sdkman`](../done/slice-v1-devcontainer-features-homebrew-sdkman.md).

**Berührte Spec-Stellen:** — (keine). [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) und [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) bleiben bis [`slice-v2-spec-devcontainer-sshd-claude-code`](../open/slice-v2-spec-devcontainer-sshd-claude-code.md) ohne den Schlüssel; das Lastenheft ändert nur der CR.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Schlüssel `claude-code` ist ein eingebautes Devcontainer-Feature und rendert das Feature von Anthropic (`ghcr.io/anthropics/devcontainer-features/claude-code`, als
Drittanbieter-Feature nach [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) ausgewiesen, Entscheidung des Projektinhabers am 2026-10-17); im Sandbox-Profil ist es erlaubt; die Egress-Allowlist trägt die gemessenen Hosts, wenn es aktiv ist, und ein Zugangsschlüssel des Agenten steht
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

**Anforderung an den Slice — Zugangsdaten des Agenten ([`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), vom CR auf Agenten-Zugangsdaten erweitert):** Ein Anthropic-Zugangsschlüssel steht **weder im Image noch im Workspace-Volume noch in `u-boot.yaml` oder
`devcontainer.json` als Wert**; zur Laufzeit gelangt er per Umgebungsvariable aus der Host-Umgebung in den Container (Referenz `${localEnv:<NAME>}`, nie der Wert; den Namen der Variable belegt die
Dokumentation des Features). Die Zusage steht nach dem CR im Lastenheft.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Feature-Optionen** — der Katalog kennt keine Optionen ([`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md), Re-Evaluierungs-Trigger); der Nutzer überschreibt den Tag über `devcontainer.features.claude-code.version`.
- **Konfiguration, Start oder Berechtigungsmodus des Agenten, `~/.claude`-Mounts** — [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) nennt es als Nicht-Ziel; ein Mount von Host-Dateien widerspricht dem Profil (kein Bind-Mount, keine Host-Geheimnisse).
- **Ein neuer `doctor`-Check-Code** — er wäre an Registry und [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) gebunden und liefe der Regel „Spezifikation nach dem Code“ zuwider; die vorhandene Prüfung (`devcontainer.sandbox.credentials`) wird erweitert, ohne neue Kennung.
- **Weitere Agenten-Features** (Codex, Gemini, andere) — jede Quelle braucht eine eigene Entscheidung.
- **Die Spezifikation** ([`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)) — sie folgt in [`slice-v2-spec-devcontainer-sshd-claude-code`](../open/slice-v2-spec-devcontainer-sshd-claude-code.md).

## 2. Definition of Done

- [ ] **Katalog und Generator:** `claude-code` im Katalog (Quelle `ghcr.io/anthropics/devcontainer-features/claude-code`, Herausgeber `anthropics`, Standardversion nach der Messung); die Drittanbieter-Warnung folgt [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md); Egress-Hosts nur bei aktivem Feature und aus der Messung; im Sandbox-Profil ist das Feature erlaubt; Golden Cases belegen Render, Warnung, Hosts und double-run -> NoOp, **ein Golden Case aktiviert `claude-code` zusammen mit `sshd`** (Welle-Closure).
- [ ] **Zugangsdaten zur Laufzeit:** `devcontainer.json` enthält bei aktivem Feature die Referenz der Umgebungsvariable (kein Wert); ohne das Feature bleibt die Ausgabe byte-identisch; die vorhandene `doctor`-Prüfung erkennt einen Klartext-Schlüssel in den gescannten Dateien (Muster aus der Dokumentation des Schlüsselformats belegt, keine neue Check-Kennung); jede neue `doctor`-Warnung zur Agenten-Variable (Quelle fehlt, Klartext) trägt das Kennzeichen `SecurityRelaxation` (Ausnahme 1 von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Kennzeichen aus [`slice-v1-cli-quiet-sicherheitswarnungen`](../done/slice-v1-cli-quiet-sicherheitswarnungen.md)) und bleibt unter `doctor --quiet` sichtbar; Tests je Fall.
- [ ] **Messung und Doku:** Die Dev-Container-CLI baut ein Image mit `claude-code` auf dem Basis-Image samt Node-Abhängigkeit (Verfahren der Messreihe in [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md)); belegt sind: der Tag (`1` oder `1.0`), die Feature-Metadaten (`mounts`, `capAdd`, `privileged`, `entrypoint`, `installsAfter`, Optionen), die Laufzeit-Hosts (DNS-Abfragen) und der Anmeldeweg; eine Folge-ADR (nächste freie Nummer) nur, wenn die Messung etwas findet, das [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) nicht trägt (etwa ein Optionsmodell). Benutzerhandbuch (Katalogtabelle, Egress-Tabelle, Abschnitt Zugangsdaten des Agenten; der Satz „Die Hosts des Agenten selbst kennt `u-boot` nicht“ ist angepasst), `README.md`, `README.de.md`, `CHANGELOG.md` (`[Unreleased]`) und `harness/replay.md` nennen den Schlüssel; der Satz „follow in a later release“ im CHANGELOG-Eintrag des CR (Lastenheft 0.11.0) ist für den Teil dieses Slice angepasst (der letzte landende Code-Slice der Welle entfernt ihn).
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
| `internal/hexagon/application/generate.go`, `devcontainer_sandbox.go` | update | Drittanbieter-Warnung (`thirdPartyFeatureWarnings`, `publisher: anthropics`) und Sandbox-Verhalten (erlaubt) |
| `internal/hexagon/application/doctor_sandbox.go` | update | `credentialPatterns()` um das Schlüsselformat erweitert, Meldungstext nennt nicht nur Git; Warn-Einträge zur Agenten-Variable mit `SecurityRelaxation`; **keine** neue Check-Kennung |
| `internal/hexagon/application/devcontainer_features_test.go`, `generate_features_test.go`, `generate_sandbox_test.go`, `doctor_sandbox_test.go`, `doctor_features_drift_test.go` | update | Katalogumfang (Herausgeber als Namensraum), Render, Warnung, Hosts, Zugangsdaten-Referenz, Klartext-Erkennung, beide Features zusammen — nach [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | update | öffentlicher Vertrag und Golden-Zeile |
| `docs/plan/adr/` (Folge-ADR, Index) | neu / update (bedingt) | nur wenn die Messung etwas findet, das [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) nicht trägt (DoD 3) |

Die Pfade stammen aus dem Quelltext (Lesestand 2026-10-17); der Implementer prüft jede Zeile vor der Änderung gegen den Code und trägt Abweichungen in die Evidence ein.

**Umsetzungsstand und Abweichungen von der Tabelle (Implementer, 2026-10-17):**

- Jede Zeile der Tabelle wurde vor der Änderung im Code gelesen; sie stimmt. `generate.go` trägt eine Zeile (`generateDevcontainer` hängt `agentNodeWarnings` an die Warnungen); `devcontainer_sandbox.go` trägt die Logik (`catalogueFeatureActive` verallgemeinert die Erkennung von `loginServiceActive`; `agentCredentialEnv` liefert `ANTHROPIC_API_KEY`; `devcontainerTemplateData` setzt `AgentEnv`). Die Drittanbieter-Warnung (`thirdPartyFeatureWarnings`) blieb unverändert: `publisher: anthropics` im Katalog genügt, sie entsteht ohne `SecurityRelaxation`. Die Folge-ADR entfällt (siehe Messung, Folgerungen).
- **Neu gegenüber dem Plan:** `generate devcontainer` weist mit einer Meldung (Code [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), ohne `SecurityRelaxation`) darauf hin, dass `claude-code` ohne das Feature `node` nicht baut (`agentNodeWarnings`); die Messung zeigte den Build-Abbruch (unten). Das Feature `node` wird nicht automatisch aktiviert.
- `doctor`: `credentialPatterns` (`doctor_sandbox.go`) trägt zwei Ausdrücke dazu (`sk-ant-`-Schlüssel, Zuweisung an `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`); `checkSandboxCredentials` meldet eine fehlende Referenz `${localEnv:ANTHROPIC_API_KEY}` bei aktivem Feature (neuer Warn-Eintrag mit `SecurityRelaxation`, kein neuer Check-Code; `coderegistry.go` und [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) blieben unberührt). Die Prüfung der Klartext-Schlüssel läuft im Sandbox-Profil unabhängig davon, ob das Feature aktiv ist.
- **Das Schlüsselformat ist belegt, aus dem Fremdcode gelesen:** Die installierte CLI (Claude Code 2.1.302, Binary `claude.exe` des npm-Pakets) enthält den eigenen Erkennungsausdruck `(?:^|[^A-Za-z0-9_-])sk-ant-[A-Za-z0-9_-]{40,}` (Anzeigename `Anthropic Key`), die Unterscheidung `sk-ant-api…` (API-Schlüssel, Header `x-api-key`) und `sk-ant-oat…` (OAuth-Token) sowie den Admin-Schlüssel `sk-ant-admin01-…`; `doctor` übernimmt den Ausdruck wörtlich. Die drei Variablennamen stehen in der Dokumentation (Umgebungsvariablen, Stand 2026-10-17) und im Binary.
- Der Name der Variable (`ANTHROPIC_API_KEY`): Quelle Dokumentation (Stand 2026-10-17), im Binary vorhanden; `u-boot` reicht nur diese eine durch.

**Messung (Dev-Container-CLI 0.89.0 im Wegwerf-Container mit Docker-Socket, Colima, arm64, Basis `mcr.microsoft.com/devcontainers/base:debian` = Debian trixie; Messordner unter dem Home-Verzeichnis, danach gelöscht; Projekt `generate devcontainer` von `u-boot` mit Sandbox-Profil, Egress, `sshd`, `claude-code`, `node`):**

| Prüfung | Ergebnis |
|---|---|
| Tags (`devcontainer features info tags`) | `1`, `1.0`, `1.0.0` bis `1.0.5`, `latest`; `1` löst auf 1.0.5 auf: Standardversion `1` |
| Feature-Metadaten (OCI-Manifest, `devcontainer-feature.json`) | `options: {}`, `installsAfter: ghcr.io/devcontainers/features/node`, `customizations.vscode.extensions: anthropic.claude-code`; **kein** `mounts`, `capAdd`, `privileged`, `entrypoint`, `containerEnv`, `dependsOn` |
| Build auf dem Basis-Image **ohne** `node` | **bricht ab** („Failed to install Node.js and npm“): der Auto-Install holt Debians `nodejs` 20.19 (der Nodesource-Eintrag für 18.x unterliegt), ohne `npm` |
| Build mit `node:1` und `claude-code:1` | baut; `claude --version` → `2.1.302 (Claude Code)`; Reihenfolge `node` vor `claude-code` (`installsAfter`) |
| Build-Hosts (Build-Protokoll) | `deb.debian.org` (Firewall-Pakete des Features; ohne `node` auch Debians `nodejs`) und, ohne `node`, `deb.nodesource.com`; der Schritt `npm install -g @anthropic-ai/claude-code` nennt seine Registry im Protokoll nicht (Standard-Registry `registry.npmjs.org`: **nicht gemessen**); die Quellen des Features `node` sind nicht aufgeschlüsselt |
| Container (`docker inspect`, Sandbox mit Egress) | `Privileged=false`, `CapAdd=[NET_ADMIN]` (von `u-boot`), keine Devices, `SecurityOpt` leer, einziges Volume `/workspaces/<projekt>`; keine Umgebungsvariable und kein Image-Metadatum mit dem Schlüssel |
| `~/.claude` | nach dem Build **nicht vorhanden**; nach dem ersten Aufruf der CLI entstehen `~/.claude/` (`backups`, `projects`, `sessions`) und `~/.claude.json` im Home von `vscode`; mit dem Schlüssel aus der Umgebung enthielt keine dieser Dateien den Schlüssel |
| Zusätzlich vom Feature installiert | `/usr/local/bin/init-firewall.sh` (nicht ausgeführt; `install.sh`: „Always set up the firewall script, but don't execute it“); `u-boot` ruft es nicht auf |
| Laufzeit-Hosts (DNS-Mitschnitt mit protokollierendem Resolver) | Aufruf ohne Schlüssel und mit leerem `ANTHROPIC_API_KEY`: `api.anthropic.com` („Not logged in“); mit einem absichtlich falschen Schlüssel: `api.anthropic.com` und `http-intake.logs.us5.datadoghq.com` („Invalid API key“) |
| Referenz im erzeugten `remoteEnv` | `devcontainer exec` löste `${localEnv:ANTHROPIC_API_KEY}` aus der Umgebung des Aufrufers auf; ein nicht gesetztes `GIT_TOKEN` blieb leer |
| Egress-Lauf mit der erzeugten Allowlist | `api.anthropic.com` erreichbar (HTTP-Antwort), `example.com` und `http-intake.logs.us5.datadoghq.com` gesperrt; `claude -p` mit falschem Schlüssel: „Invalid API key“ (Antwort des Servers), ohne Schlüssel „Not logged in“; Port 2222 lauscht (`sshd`) |

Folgerungen: Die Egress-Zeile trägt `registry.npmjs.org`, `api.anthropic.com` (gemessen) sowie `claude.ai`, `platform.claude.com`, `claude.com` (Anmeldung; **laut Dokumentation, nicht gemessen**). Das Feature unterläuft das Profil nicht (Risiko „Feature-Metadaten“ entfällt für Version 1.0.5, der Pin `1` schwimmt). Eine Folge-ADR ist nicht nötig: [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) trägt Herausgeber als Namensraum, Hinweis und Egress-Hosts; ein Optionsmodell wird nicht gebraucht. Die Voraussetzung `node` ist eine Folge des Features, kein Mangel des Modells.

**Nicht aufgenommen, optional (Entscheidung des Projektinhabers; laut Netzwerk-Dokumentation von Anthropic, 2026-10-17):** `mcp-proxy.anthropic.com` (MCP-Connectors; abschaltbar mit `ENABLE_CLAUDEAI_MCP_SERVERS=false`), `downloads.claude.ai` (nativer Installer, Updater, Plugin-Downloads), `storage.googleapis.com`, `github.com` und `raw.githubusercontent.com` (Plugins, Changelog; die beiden GitHub-Hosts stehen ohnehin in der Basisliste), `http-intake.logs.us5.datadoghq.com` und `browser-intake-us5-datadoghq.com` (Telemetrie und Fehlerberichte; abschaltbar mit `DISABLE_TELEMETRY`, `DISABLE_ERROR_REPORTING` oder gemeinsam mit `CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC`). Begründung: Die Default-Allowlist trägt nur, was Installation und Anmeldung oder API-Nutzung brauchen; die CLI lief im Test mit gesperrtem Telemetrie-Host. Der Nutzer ergänzt Hosts über `devcontainer.sandbox.egress.allow`, das Handbuch nennt es.

**Grenzen der Messung:** (1) Ein echter Login und eine Anfrage mit gültigem Schlüssel wurden nicht ausgeführt (kein Konto, kein Schlüssel); die Hosts `claude.ai`, `platform.claude.com`, `claude.com` und der Speicherort der Zugangsdaten bei einem Login im Container sind „laut Dokumentation, nicht gemessen“. (2) Die Hosts des Builds stammen aus dem Build-Protokoll, nicht aus einem DNS-Mitschnitt; die Laufzeit-Hosts aus einem protokollierenden Resolver für die Fälle ohne, mit leerem und mit falschem Schlüssel. (3) Nur arm64 und Colima; Feature `1.0.5`, Claude Code 2.1.302; der Pin `1` schwimmt, und `npm install -g` holt die jeweils aktuelle CLI, spätere Versionen können abweichen. (4) Der Feature-Quelltext kam über das OCI-Artefakt von `ghcr.io/anthropics/devcontainer-features/claude-code`, nicht über das Repository. (5) Der `sk-ant-`-Ausdruck stammt aus dem Binary der installierten CLI (Version 2.1.302), nicht aus einer Dokumentationsaussage über das Schlüsselformat; Anthropic kann das Format ändern.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der CR [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md) und die Slices
[`slice-v1-cli-quiet-sicherheitswarnungen`](../done/slice-v1-cli-quiet-sicherheitswarnungen.md) (liefert das Kennzeichen am `doctor`-Eintrag) und [`slice-v1-devcontainer-feature-sshd`](../done/slice-v1-devcontainer-feature-sshd.md) liegen in `done/` (die beiden Feature-Slices berühren Katalog, Egress und Tests; WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice überschreitet 3 Liefer-Punkte, etwa weil die Zugangsdaten-Übergabe einen eigenen Slice verlangt.
- `in-progress` → `open` (blockiert): die Messung findet eine Lockerung des Profils, die der CR nicht trägt (dann zuerst CR erweitern); das Feature von Anthropic baut nicht oder ist nicht mehr unter diesem Namen verfügbar.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Das Feature von Anthropic ändert sich oder verschwindet** (Drittanbieter, Tag-Pin): ein Build-Fehler fällt beim Image-Build auf. — **Ausgang:** weiter offen: Re-Evaluierungs-Trigger von [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) (Änderung oder Wegfall des Features des Herausgebers).
- **Die Egress-Hosts reichen nicht** (Anmeldung, Telemetrie, Update-Prüfung des Agenten sind nicht belegt): der Agent läuft im Sandbox-Profil mit Egress-Restriktion nicht. — **Ausgang:** entfallen, sobald die Messung (DoD 3) die Hosts belegt; weitere Hosts trägt der Nutzer über `devcontainer.sandbox.egress.allow` ein, das Handbuch nennt es (bei Closure bestätigen).
- **Das Feature bringt Metadaten mit** (`mounts`, `capAdd`, `privileged`, `entrypoint`) und unterläuft damit das Profil (kein Host-Bind-Mount, keine Host-Geheimnisse). — **Ausgang:** weiter offen: die Messung entscheidet; bei Eintritt Rückführung an den CR (§4).
- **Der Anmeldeweg des Agenten legt Zugangsdaten im Container ab** (Anmeldung im Container statt Umgebungsvariable): sie lägen im Home des Container-Benutzers, nicht im Image oder Volume des Profils. — **Ausgang:** weiter offen: die Messung (DoD 3) belegt den Weg; liegt der Zugangsschlüssel im Container statt in der Umgebungsvariable, wird der Text von [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) im CR nachgezogen (Rückführung, §4); das Handbuch nennt den gemessenen Weg.
- **Schlüsselformat und Variablenname sind geraten**, wenn die Dokumentation sie nicht nennt. — **Ausgang:** entfallen: DoD 2 verlangt den Beleg aus der Dokumentation; ohne Beleg wird das Muster gestrichen und der Slice meldet die Lücke.
- **Plan-Aussage über den Code stimmt nicht** (siehe §8): die Fundstellen in §1 und §3 sind ein Lesestand. — **Ausgang:** weiter offen: → [BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) im Register (verkörpert in `AGENTS.md §3.10`); der Implementer prüft jede Aussage vor der Änderung im Code, Abweichungen stehen in der Evidence.

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `internal/hexagon/application` (Kürzel APP) und, für die Doku, `docs/user/`, `harness/`, `docs/plan/` und `spec/` (DOK) — beide Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006) (Abschnitt Modus-Deklaration), beide erfüllen die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. `internal/adapter/driving/cli` (CLI) wird nur berührt, falls ein neuer Warn-Code in `coderegistry.go` stehen muss; dann gilt `AGENTS.md §3.9` (Kommentar-Selbstcheck); `AGENTS.md §3.10` gilt für jede Plan-Aussage über Code.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer für APP und DOK. Die einzige Beobachtung `BEO-CLI/kommentar-beschreibt-abwesendes` (3×, verkörpert in `AGENTS.md §3.9`) betrifft CLI; sie wird nur bei einer CLI-Berührung relevant. Die Klasse „Plan-Aussage über Code stimmt nicht“ steht als [BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) im Register (vier Belege, verkörpert in `AGENTS.md §3.10`): jede Aussage über den Code wird vor der Änderung an der Stelle gelesen; Abweichungen stehen in der Evidence.

Alle berührten Sub-Areas GF.
