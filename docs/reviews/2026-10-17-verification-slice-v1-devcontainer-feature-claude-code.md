# Verification Evidence: slice-v1-devcontainer-feature-claude-code — 2026-10-17 <!-- d-check:ignore (Slice-Kennung im Titel) -->

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** Slice `slice-v1-devcontainer-feature-claude-code` (Welle `welle-devcontainer-sshd-claude-code`, vierter Slice; liegt in `in-progress/`), Commits `18cdfbc` (Code, Tests), `b7f888e` (Handbuch, READMEs, CHANGELOG, Replay-Zeile), `501cabb` (Slice §3: Umsetzungsstand, Messung, Grenzen), `6320e4e` (Review-Fixes im Code) und `153a600` (Review-Fixes in Doku und Slice). Maßstab: §1, §2 (DoD), §3 (Plan, Umsetzungsstand, Messung), §6 (Risiken), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Lastenheft 0.11.0), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md), `harness/replay.md`. Review-Report: [`2026-10-17-slice-v1-devcontainer-feature-claude-code`](2026-10-17-slice-v1-devcontainer-feature-claude-code.md) (0 HIGH, 1 MEDIUM, 2 LOW, 3 INFO; lief gegen `501cabb`, die Fixes stehen in `6320e4e` und `153a600`). Entscheidung des Projektinhabers vom 2026-10-17: Default-Allowlist von `claude-code` nur `registry.npmjs.org` und `api.anthropic.com`; Anmelde-Hosts trägt der Nutzer ein; OAuth-Weg nicht zugesagt. <!-- d-check:ignore (Kennungen und Pfade als Zitat) -->

## Verification Evidence

Scope:
- Slice: `slice-v1-devcontainer-feature-claude-code` <!-- d-check:ignore (Slice-Kennung als Zitat) -->
- IDs: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md)
- Artefakte: `devcontainer_features.go`, `devcontainer_egress.go`, `devcontainer_sandbox.go`, `doctor_sandbox.go`, `generate.go`, `templates.go`, `devcontainer.json.tmpl`, `coderegistry.go`, die Tests (`generate_features_test.go`, `generate_sandbox_test.go`, `doctor_sandbox_test.go`, `doctor_features_drift_test.go`, `devcontainer_features_test.go`), `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` <!-- d-check:ignore (Pfade als Zitat) -->

DoD-Abgleich (§2 des Slice; die Häkchen im Slice stehen noch aus und sind Closure-Arbeit):
- [x] **Katalog und Generator** — Eintrag `claude-code` (Quelle `ghcr.io/anthropics/devcontainer-features/claude-code`, Version `1`, `publisher: anthropics`); Drittanbieter-Warnung über das vorhandene `thirdPartyFeatureWarnings` ohne `SecurityRelaxation` (`TestGenerateDevcontainer_ClaudeCodeRender`); Sandbox: erlaubt, Erkennung per aufgelöster Quelle (`catalogueFeatureActive`, Enabled-Filter); Egress-Zeile `{"registry.npmjs.org", "api.anthropic.com"}` nur bei aktivem Feature (`…SandboxEgress_ClaudeCodeHostsOnlyWhenActive`, `…_FeatureHosts/claude-code` mit Negativliste der Anmelde-, Telemetrie-, MCP- und Update-Hosts); Golden Case mit `sshd` zusammen (`TestGenerateDevcontainer_SshdAndClaudeCodeTogether`: vier Profile, jeder Schlüssel einmal, double-run → NoOp, Egress als sortierte duplikatfreie Vereinigung). **Plan-Diskrepanz D-1:** der Wortlaut „Egress-Hosts … aus der Messung“ trifft für `registry.npmjs.org` nicht zu (siehe unten).
- [x] **Zugangsdaten zur Laufzeit** — `AgentEnv` aus `agentCredentialEnv` (nur `ANTHROPIC_API_KEY`), Template rendert `"{{.}}": "${localEnv:{{.}}}"`; `TestGenerateDevcontainer_ClaudeCodeCredentialReference` (Referenzform, kein `sk-ant`, ohne Feature byte-identischer `remoteEnv`-Block); `credentialPatterns` um `sk-ant-`-Ausdruck und Zuweisung der drei Variablen erweitert, Herkunft im Kommentar und im Slice benannt (Binary der CLI 2.1.302); keine neue Check-Kennung (`coderegistry.go` ändert nur den Label-Text, Schlüssel unverändert); jede neue Warnung (Klartext, fehlende Referenz) trägt `SecurityRelaxation`, kein `DegradationFallback` (`TestDoctor_SandboxCredentials_Agent`, 17 Fälle; der Test prüft je Fall das Kennzeichen und dass der Fund den Schlüssel nicht wiederholt). Sichtbarkeit unter `doctor --quiet` folgt aus dem Kennzeichen (Mechanismus aus dem Vorgänger-Slice, im Review geprüft).
- [~] **Messung und Doku** — Doku erfüllt, Messung mit zwei benannten Abweichungen vom DoD-Wortlaut: **D-2** „Anmeldeweg belegt“ ist nicht erfüllt (ein echter Login ist ohne Konto nicht messbar; belegt ist nur der Weg über `ANTHROPIC_API_KEY`, ohne und mit falschem Schlüssel); **D-1** (Hosts). Tag (`1` → 1.0.5), Metadaten (`options: {}`, `installsAfter`, kein `mounts`/`capAdd`/`privileged`/`entrypoint`), Laufzeit-Host `api.anthropic.com` (DNS-Mitschnitt) sind belegt; die Folge-ADR entfällt zu Recht ([`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md) trägt Herausgeber als Namensraum). Handbuch (Katalogtabelle, Egress-Tabelle, Abschnitt Claude Code, „Zugangsdaten (Git und Agent)“), `README.md`, `README.de.md`, `CHANGELOG.md` (`### Added`, Satz im Lastenheft-0.11.0-Absatz), `harness/replay.md` nennen den Schlüssel; der Satz „follow in a later release“ ist im Lastenheft-0.11.0-Absatz entfernt (der Restsatz „follow in V1“ in Zeile 1138 betrifft OpenTelemetry, nicht diesen Slice).
- [x] `make gates` grün — selbst gefahren, Exit 0.
- [~] Review durchgeführt, Report liegt vor — Report liegt vor; Bewertung der Findings unten (F-1 bis F-4 und F-6 behoben, F-5 Grenze). Die Fixes `6320e4e`/`153a600` sind nicht erneut gereviewt; ich habe sie selbst am Code gelesen.
- [ ] Doku-Update bei öffentlichem Vertrag — inhaltlich erledigt (siehe DoD „Messung und Doku“); Häkchen ist Closure.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen (§7 „Wird bei der Closure gefüllt“).
- [ ] Beobachtungs-Register fortgeschrieben, oder „keine Beobachtung“ in §7 — offen, Closure; Beleg-Kandidaten unter (j).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — Ausgang steht an allen sechs; zwei Ausgänge sind veraltet (V-2, siehe (k)).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) — Closure.

### (a) Plan-Diskrepanzen mit Vorschlag

| ID | Stelle im Slice | Befund | Vorschlag |
|---|---|---|---|
| D-1 | §1 Ziel: „die Egress-Allowlist trägt die gemessenen Hosts“; DoD 1: „Egress-Hosts nur bei aktivem Feature und aus der Messung“ | Die Default-Liste hat zwei Hosts: `api.anthropic.com` (gemessen, DNS-Mitschnitt) und `registry.npmjs.org` (nicht gemessen; Begründung: das nötige Feature `node` bringt ihn ohnehin, ein Laufzeit-Bedarf der CLI ist nicht gemessen). Der Wortlaut „gemessenen Hosts“ ist damit falsch, die Entscheidung des Projektinhabers vom 2026-10-17 steht nur in §3 („Folgerungen“). | §1 und DoD 1 umschreiben: „die Egress-Allowlist trägt für das Feature `api.anthropic.com` (gemessen) und `registry.npmjs.org` (wegen des nötigen Features `node`; Laufzeit-Bedarf nicht gemessen); Anmelde-, Telemetrie-, MCP- und Update-Hosts trägt der Nutzer über `devcontainer.sandbox.egress.allow` (Entscheidung des Projektinhabers 2026-10-17)“. |
| D-2 | DoD 3: „belegt sind: … der Anmeldeweg“ | Nur der Weg über die Umgebungsvariable ist belegt; Login im Container, Abo-Login und `CLAUDE_CODE_OAUTH_TOKEN` sind nicht gemessen und (Entscheidung) nicht zugesagt; Grenze 1 in §3 sagt es, DoD 3 nicht. | DoD 3: „der Anmeldeweg über `ANTHROPIC_API_KEY` (Login und OAuth-Token sind nicht zugesagt, nicht gemessen)“. |
| D-3 | §6 Risiko 4 (Anmeldeweg legt Zugangsdaten im Container ab): „die Messung (DoD 3) belegt den Weg“ | Die Messung belegt den Login-Weg nicht (siehe D-2); der Ausgang verspricht einen Beleg, den es nicht gibt. | Ausgang auf „weiter offen: Login im Container nicht gemessen, nicht zugesagt (Entscheidung 2026-10-17); `~/.claude` entsteht erst nach dem ersten Aufruf und enthielt den Schlüssel aus der Umgebung nicht (gemessen); Rückführung nur, wenn ein Login-Weg zugesagt werden soll“. |
| D-4 | §1 „Planungsgrundlage“, „erwartet wird die npm-Registry beim Build und `api.anthropic.com` zur Laufzeit“ | Erwartung, durch die Messung nur halb bestätigt (Build-Host der npm-Registry nicht im Protokoll, Laufzeit-Host bestätigt). Kein Mangel, weil der Absatz als „erneut zu prüfen“ markiert ist. | keine Änderung nötig; bei der Closure im Lerneintrag erwähnen (Erwartung vs. Befund). |

### (b) Katalog, Generator, Render

- `featureCatalogue()["claude-code"]`: Quelle, Version `1`, `publisher: anthropics`, `shortDesc` mit Hinweis auf `node`; `TestFeatureCatalogue_KeysCoverSpecExamples` führt den Schlüssel (13 Einträge); `TestDoctor_FeaturesDrift_ClaudeCode` (in sync / fehlender Render-Schlüssel); das Zusammenlegen der sshd-Drift-Tests (`checkCatalogueFeatureDrift`) lässt beide Fälle für beide Schlüssel laufen, nichts gelockert.
- `agentNodeWarnings` (Code [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), `warn`, ohne `SecurityRelaxation`): Bedingung „`claude-code` aktiv, `node` nicht aktiv“; die Messung belegt den Build-Abbruch ohne `node`. Der Hinweis ist eine Neuerung gegenüber dem Plan (im Umsetzungsstand ausgewiesen); `…ClaudeCodeNodeNotice` deckt mit/ohne `node`, Feature aus, eigene Quelle. Kein bestehender Golden Case ändert sich (Hinweis nur bei aktivem Feature).
- Egress: `egressFeatureHosts()["claude-code"] = {"registry.npmjs.org", "api.anthropic.com"}`; Kommentar nennt Messung und Nicht-Messung zutreffend (F-2 behoben). Anmelde-Hosts (`claude.ai`, `platform.claude.com`, `claude.com`) stehen in keinem Code-Pfad außer als Negativ-Erwartung in Tests (Suche im Baum).

### (c) Messung (Slice §3) — belegt vs. „laut Dokumentation“

Die Messung (Dev-Container-CLI 0.89.0, Colima, arm64, Basis `mcr.microsoft.com/devcontainers/base:debian`) habe ich **nicht wiederholt**: im Repo nicht ausführbar (Docker-Socket-Image ad hoc, wie bei den Vorgängern; Grenze 4 des Slice). Bewertung nach Formulierung:

| Aussage | Belegart |
|---|---|
| Tag `1` → 1.0.5, `options: {}`, `installsAfter`, kein `mounts`/`capAdd`/`privileged`/`entrypoint` | gemessen (OCI-Manifest, `devcontainer.json` des Features) |
| Build ohne `node` bricht ab | gemessen |
| `api.anthropic.com` zur Laufzeit; `http-intake.logs.us5.datadoghq.com` bei falschem Schlüssel | gemessen (protokollierender Resolver) |
| CLI läuft mit gesperrtem Telemetrie-Host | gemessen (Egress-Lauf) |
| `~/.claude` erst nach erstem Aufruf, Schlüssel nicht in den Dateien | gemessen |
| `${localEnv:ANTHROPIC_API_KEY}` wird aufgelöst, nicht gesetztes `GIT_TOKEN` bleibt leer | gemessen |
| `registry.npmjs.org` zur Laufzeit | **nicht gemessen** (im Slice, Handbuch und Code-Kommentar so benannt) |
| `claude.ai`, `platform.claude.com`, `claude.com`; Speicherort bei Login; Abo-Login | **laut Dokumentation, nicht gemessen** (Grenze 1; Handbuch ebenso) |
| Variablenname `ANTHROPIC_API_KEY`, Vorrang vor Abo-Login | Dokumentation (Stand 2026-10-17), Variablenname auch im Binary |
| `sk-ant-`-Ausdruck | aus dem Binary der CLI 2.1.302 gelesen, nicht aus der Dokumentation (Grenze 5) |

Vorsichtig formuliert: ja. Zwei Stellen gehen nicht über den Beleg hinaus, sind aber in der Closure mitzunehmen: „Risiko ‚Feature-Metadaten‘ entfällt“ gilt für 1.0.5 (Pin `1` schwimmt; im Text benannt), und die Nicht-Wiederholbarkeit im Repo gehört als Satz in §7 (CLI 0.89.0, Feature 1.0.5, Basis-Image `…/base:debian`, Claude Code 2.1.302, Profil Sandbox mit Egress, `sshd`, `claude-code`, `node`).

### (d) Review-Findings am Code geprüft

| Finding | Bewertung | Beleg am Code |
|---|---|---|
| F-1 MEDIUM (nicht gemessene Hosts) | **behoben** | `devcontainer_egress.go`: nur `registry.npmjs.org`, `api.anthropic.com`; `…ClaudeCodeHostsOnlyWhenActive` verlangt: Anmelde-Hosts nie in der Liste; `…_FeatureHosts/claude-code` verlangt Negativliste; CHANGELOG nennt zwei Hosts und verweist auf `egress.allow`; Handbuch ebenso mit „Entscheidung des Projektinhabers“; Suche nach `claude.ai`/`platform.claude`/`claude.com` im Baum: nur Handbuch (als „nicht in der Liste“), Slice, Tests (Negativ). Entscheidung steht in §3 „Folgerungen“, nicht im DoD-Wortlaut (D-1). |
| F-2 LOW (`registry.npmjs.org` mit „Installation“ begründet) | **behoben** | Kommentar in `devcontainer_egress.go` und Handbuch: Host bringt das nötige Feature `node` mit, Laufzeit-Bedarf der CLI nicht gemessen, Beschränkung wirkt erst beim Containerstart. |
| F-3 LOW (Label und OK-Meldungen nur Git) | **behoben** | `coderegistry.go`: „Sandbox git and agent credentials“; OK-Meldungen „No plaintext git or agent credentials found; …“, „credential check skipped“; Test prüft Teiltexte (`reference only`, `no https clone`). |
| F-4 INFO (Format aus Binary) | **behoben im DoD-Wortlaut** | DoD 2 sagt „Muster aus dem Binary der installierten CLI belegt, die Dokumentation nennt das Schlüsselformat nicht“; Kommentar von `credentialPatterns`, Handbuch und Test-Kommentar nennen die Herkunft. |
| F-5 INFO ([`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md) nennt nur `homebrew`/`sdkman`) | **Grenze, nicht blockierend** | Katalogeintrag stützt sich auf die generische Regel (Herausgeber als Namensraum); ADR bleibt unverändert (`AGENTS.md` §3.5). Fürs Register als Beobachtung möglich, kein Closure-Blocker. |
| F-6 INFO (Platzhalter löst Warnung aus) | **behoben als dokumentierte Grenze** | Test `placeholder value is reported (known limit)`, Handbuch „Bekannte Grenze“, Slice §3 „Bekannte Grenze (Review F-6)“. Das Verhalten ist unverändert und gewollt (RE2 ohne Vorausschau). |

### (e) Plan-vs-Code-Diff (§3 gegen `git diff --stat d795199..HEAD`, 19 Dateien)

| Datei im Diff | Plan-Zeile §3 | Befund |
|---|---|---|
| `devcontainer_features.go`, `devcontainer_egress.go`, `devcontainer_sandbox.go`, `doctor_sandbox.go`, `generate.go`, `templates.go`, `devcontainer.json.tmpl` | ja | konform |
| `devcontainer_features_test.go`, `generate_features_test.go`, `generate_sandbox_test.go`, `doctor_sandbox_test.go`, `doctor_features_drift_test.go` | ja | konform |
| `coderegistry.go` (nur Label-Text) | in §3 nicht als Zeile; §3 nennt es als „nur bei neuem Code“ und §1 „Neue `doctor`-Check-Codes … gebunden“ | **Abweichung klein:** Label geändert, kein neuer Code (Review F-3). Im Plan-Umsetzungsstand steht „`coderegistry.go` … blieben unberührt“ — das ist seit `6320e4e` nicht mehr richtig (D-5, siehe unten). |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | ja | konform |
| Slice-Datei | nicht in §3 | Prozess-Artefakt |
| `docs/plan/adr/` | bedingt | nicht berührt, Bedingung nicht eingetreten (bestätigt) |

**D-5 (Plan-Aussage über Code, `AGENTS.md` §3.10):** Umsetzungsstand sagt, `coderegistry.go` und [`SPEC-017`](../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) blieben unberührt; `coderegistry.go` ist im Diff (Label-Text von `devcontainer.sandbox.credentials`, Schlüssel unverändert). Vorschlag: Satz berichtigen („`coderegistry.go` ändert nur den Label-Text, keinen Schlüssel; SPEC-017 unberührt“) und die Datei in die Plan-Tabelle aufnehmen. Gleichzeitig steht `spec/spezifikation.md` (SPEC-017-Zeile „keine Klartext-Git-Zugangsdaten“) hinter dem Code; das trägt der Folge-Slice `slice-v2-spec-devcontainer-sshd-claude-code` (laut Welle gewollt). <!-- d-check:ignore (Slice-Kennung als Zitat) -->

Nicht-Ziele §1 eingehalten: keine Optionen, kein `~/.claude`-Mount, kein neuer Check-Code, keine weiteren Agenten-Features, Spezifikation und Lastenheft unberührt (`spec/` nicht im Diff), keine Gate-Lockerung (`.golangci.yml` unberührt, kein `//nolint`).

### (f) Replay-Evidence (`harness/replay.md`)

Die Zeile 55 nennt `claude-code` allein und zusammen mit `sshd`, die Egress-Hosts (nur die zwei, ohne Anmelde-/Telemetrie-/MCP-/Update-Hosts), `doctor`-Fälle und die Testnamen; die genannten Funktionen stehen im Diff (`TestGenerateDevcontainer_ClaudeCode*`, `…SshdAndClaudeCodeTogether`, `…ClaudeCodeHostsOnlyWhenActive`, `TestDoctor_SandboxCredentials_Agent`). Der Test vergleicht Egress gegen eine ausgeschriebene Liste, nicht gegen `egressHosts()`.

### (g) CHANGELOG und Handbuch gegen den Code

- CHANGELOG `### Added`: `claude-code:1`, Drittanbieter-Hinweis, Pflicht zu `node`, `remoteEnv`-Referenz neben `GIT_TOKEN`, nur zwei Hosts, Verweis auf `egress.allow`: stimmt mit (b). Absatz Lastenheft 0.11.0: letzter Satz nennt nur Geliefertes (Check `devcontainer.sandbox.credentials` ohne neue Kennung, `sk-ant-`, drei Variablen, Referenz bei aktivem Feature). Release-Cut-Szenario: der Eintrag ist im Indikativ und vollständig lieferbar.
- Handbuch: Katalogtabelle, Egress, Abschnitt Claude Code (Quellen und Grenzen je Aussage benannt), „Zugangsdaten (Git und Agent)“ inkl. Platzhalter-Grenze: stimmt mit Code und Messung. Der Satz „Die Hosts des Agenten selbst kennt `u-boot` nicht“ ist angepasst. Keine Behauptung, die OAuth-Weg oder Anmelde-Hosts zusagt.
- README und README.de: 13 Features, stimmt mit `featureCatalogue()` überein (Review geprüft; Test trägt die Liste).

### (h) Eigene Sensor-Kontrolle

Mutationen habe ich nicht wiederholt (das Review fuhr sieben; die Fixes `6320e4e` ändern die Host-Liste, die Texte und das Label: die Host-Mutation des Reviews, „Host `claude.com` entfernt“, ist durch die neue Fassung gegenstandslos, ihr Gegenstück ist `…ClaudeCodeHostsOnlyWhenActive` mit Negativliste und `…_FeatureHosts/claude-code`, die jede Zusatz- oder fehlende Host-Zeile als rot erwarten).

### (i) Risiken §6 — Form der Ausgänge und Bewertung

| Risiko | Ausgang im Slice | Befund |
|---|---|---|
| 1 Feature ändert sich oder verschwindet | `weiter offen:` Re-Evaluierungs-Trigger von [`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md) | kanonisch, vertretbar (Drittanbieter, Pin `1` schwimmt) |
| 2 Egress-Hosts reichen nicht | `weiter offen:` API-Host gemessen, Anmelde-Hosts nicht gemessen und nicht in der Default-Liste, Nutzer trägt ein; „bei Closure bestätigen“ | **vertretbar, siehe unten** |
| 3 Feature bringt Metadaten mit | `weiter offen: die Messung entscheidet` | **veraltet (V-2):** die Messung hat entschieden → `entfallen: gemessen (§3): kein mounts/capAdd/privileged/entrypoint, nur installsAfter; Feature 1.0.5, Pin 1 schwimmt` |
| 4 Anmeldeweg legt Zugangsdaten im Container ab | `weiter offen: die Messung (DoD 3) belegt den Weg` | **veraltet (D-3)** |
| 5 Format/Variablenname geraten | `entfallen:` Variablen dokumentiert, Format aus Binary belegt | kanonisch, stimmt (Grenze 5 mitgenommen) |
| 6 Plan-Aussage über Code | `weiter offen: → BEO-DOK/plan-aussage-ueber-code` | kanonisch; Eintritt belegt (D-5 und F-3: Plan sagt „unberührt“, Code berührt; Plan sagte „gemessene Hosts“, Code trägt einen nicht gemessenen) |

**Risiko 2 „weiter offen“ — für die Closure vertretbar.** Gründe: (1) [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) verlangt eine dokumentierte Default-Allowlist (Basis plus Ergänzung je Feature) und die Nutzer-Liste `egress.allow`, aber nicht, dass jeder Agenten-Weg ohne Zutun läuft; beides ist geliefert und getestet. (2) Der gemessene Weg (Schlüssel aus der Umgebung) läuft im generierten Allowlist-Container (API erreichbar, `example.com` gesperrt, Telemetrie-Host gesperrt ohne Folgen). (3) Der Rest ist durch die Entscheidung des Projektinhabers eingegrenzt und im Handbuch, CHANGELOG und Slice als Nutzerpflicht benannt (Handbuch nennt Hosts, Optionen und `egress.allow`). (4) Eine Schließung wäre ohne Konto nicht messbar, ein Raten widerspräche [`ADR-0012`](../plan/adr/0012-devcontainer-egress-firewall.md) („zu weit ist wertlos“). Es ist ein offenes Risiko mit Ausgang, kein verschwiegenes. **Vorschlag für §7 (Wortlaut):** „weiter offen: Die Default-Allowlist von `claude-code` trägt `api.anthropic.com` (gemessen) und `registry.npmjs.org` (wegen `node`, Laufzeit-Bedarf nicht gemessen), Entscheidung des Projektinhabers 2026-10-17. Die Hosts der Anmeldung (`claude.ai`, `platform.claude.com`, `claude.com`; laut Dokumentation, nicht gemessen) und die Hosts von Telemetrie, MCP, Update und Plugins trägt der Nutzer über `devcontainer.sandbox.egress.allow` ein; das Handbuch nennt sie. OAuth- und Login-Weg sind nicht zugesagt. Trigger: ein Nutzer-Bericht oder eine Messung mit Konto; dann Re-Evaluierung der Liste in `slice-v2-spec-devcontainer-sshd-claude-code` (SPEC-019) oder einem Folge-Slice.“ <!-- d-check:ignore (Kennungen im Vorschlagstext) -->

Die Spezifikation ([`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)) muss im Folge-Slice dieselben zwei Hosts tragen; das ist dort vorzumerken (Welle-Übergabe, V-3).

### (j) Beobachtungs-Zähler (Belege für die Closure)

- **Plan-Aussage über Code ([`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md), `AGENTS.md` §3.10): ja, trat auf**, zweimal: (1) Plan §1/DoD 1 „gemessene Hosts“ vs. Code mit nicht gemessenem Host (D-1, Review F-1); (2) Umsetzungsstand „`coderegistry.go` unberührt“ vs. Diff (D-5). Die Regel ist verkörpert; die Evidence-Datei im Register ist Closure-Arbeit.
- **Kommentar beschreibt Abwesendes ([`BEO-CLI/kommentar-beschreibt-abwesendes`](../plan/planning/observations/BEO-CLI/kommentar-beschreibt-abwesendes/observation.md), `AGENTS.md` §3.9): nein.** Der Diff der Fixes habe ich auf (a) bis (c) gelesen; Kommentare in `devcontainer_egress.go`, `doctor_sandbox.go`, `devcontainer_sandbox.go` sind im Indikativ; ein Treffer der Suche nach alten Namen (`git credentials`, `Git-Zugangsdaten`) steht nur in `Dockerfile.tmpl` (Git-Zeile, zutreffend), im Lastenheft/der Spezifikation (Spec-Scope) und im Handbuch-Abschnitt über Git. Kein Zähler-Eintrag.
- Review-Finding-Klassen (aus dem Report, in §7 und den Zähler): Egress-Host ohne Messbeleg in der Default-Allowlist · Egress-Begründung nennt Build-Zeit-Bedarf für eine Laufzeit-Allowlist · Beschreibung verweist auf zu engen Prüfumfang · Beleg-Quelle weicht vom DoD-Wortlaut ab · Katalogeintrag ohne eigene ADR-Entscheidung · Platzhalter-Wert löst Klartext-Warnung aus.

## Sensors

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0 (`baseline-verify: v6.18.1 OK`, `coverage-gate: OK — coverage 91.20% meets threshold 90%`, `d-check: 299 Datei(en) geprüft, 0 Befund(e)`; `lint` und `test` im selben Lauf grün) |
| `make doc-trace` | pass | selbst gefahren: `83 Anforderung(en), 0 Waise(n)` |
| `make docs-check` | pass | Teil von `make gates`; nach Anlage dieser Datei erneut gefahren (siehe Übergabe) |
| `make test` mit Mutationen | not run (eigene) | Review fuhr sieben Mutationen; die Fixes ändern Host-Liste, Texte, Label: durch Tests mit Negativliste gedeckt (siehe (h)) |
| `make test-docker` | not run | kein Docker-Laufzeitpfad berührt: Generator, Doctor-Check, Doku; keine Compose-/Engine-Aufrufe |
| `make verify-depguard` | not run | keine Importregel berührt, kein neuer Import in Nicht-Testcode |
| Messung mit der Dev-Container-CLI | nicht wiederholt | im Repo nicht ausführbar; Bewertung nach Formulierung, siehe (c) |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | Katalogeintrag, `TestFeatureCatalogue_KeysCoverSpecExamples`, `TestGenerateDevcontainer_ClaudeCodeRender`, `…ClaudeCodeNodeNotice`, `TestDoctor_FeaturesDrift_ClaudeCode` |
| [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | Drittanbieter-Warnung über `publisher: anthropics` (`…ClaudeCodeRender`) |
| [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) | Sandbox erlaubt, `catalogueFeatureActive`, `…SshdAndClaudeCodeTogether` (Sandbox-Fälle) |
| [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) | `…SandboxEgress_ClaudeCodeHostsOnlyWhenActive`, `…_FeatureHosts/claude-code` |
| [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | `…ClaudeCodeCredentialReference`, `TestDoctor_SandboxCredentials_Agent` (17 Fälle) |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | `SecurityRelaxation` an jeder neuen `doctor`-Warnung (`TestDoctor_SandboxCredentials_Agent`) |
| [`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md) | Herausgeber als Namensraum, keine Optionen, Messgrenze wie dort |

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: none berührt

Nicht ausgeführt: siehe Sensors (`make test-docker`, `make verify-depguard`, eigene Mutationen, Wiederholung der Messung).

Commit / Artefakt: `18cdfbc`, `b7f888e`, `501cabb`, `6320e4e`, `153a600` (Range `d795199..HEAD`).

## Verdikt

**Verifiziert: ja, ohne blockierende Abweichung.** Die Liefer-Punkte (Katalog und Generator, Zugangsdaten zur Laufzeit, Doku) sind am Code, an den Tests und an Handbuch, CHANGELOG und Replay-Zeile getragen; die Review-Findings F-1 bis F-4 und F-6 sind am Code behoben (F-6 als dokumentierte Grenze), F-5 ist Grenze. Die Messung ist belegt und vorsichtig formuliert; der Nicht-Messbereich (Anmeldung, `registry.npmjs.org` zur Laufzeit) ist in Slice, Handbuch und CHANGELOG ausgewiesen.

**Auflagen für die Closure (keine blockiert den Übergang nach `done/`, alle gehören vorher in den Slice):**
1. Plan-Wortlaut an die Entscheidung des Projektinhabers anpassen: §1 und DoD 1 (D-1), DoD 3 (D-2), §6 Risiko 4 (D-3), Umsetzungsstand und Plan-Tabelle zu `coderegistry.go` (D-5).
2. §6 Risiko 3 auf `entfallen` mit gemessener Begründung; Risiko 2 mit dem Vorschlagstext aus (i) stehen lassen (`weiter offen`).
3. §7 nennt: Messlauf im Repo nicht wiederholbar (CLI 0.89.0, Feature 1.0.5, Claude Code 2.1.302, Basis `…/base:debian`, Colima/arm64), Pin `1` schwimmt; `registry.npmjs.org` nicht gemessen; Platzhalter-Grenze von `doctor`.
4. Beobachtungs-Register: `BEO-DOK/plan-aussage-ueber-code` bekommt die Evidence-Datei dieses Slice mit den zwei Befunden aus (j); für `BEO-CLI/kommentar-beschreibt-abwesendes` „keine Beobachtung“ in §7 notieren.
5. Planner, Welle-Übergabe: `slice-v2-spec-devcontainer-sshd-claude-code` trägt in `SPEC-019` genau die zwei Hosts, in `SPEC-017` die erweiterte Check-Beschreibung und in `SPEC-023` den Schlüssel; die Welle ist erst fertig, wenn der Folge-Slice geschlossen ist (der CHANGELOG-Restsatz „follow in a later release“ ist bereits entfernt). <!-- d-check:ignore (Slice-Kennung im Auflagentext) -->
6. Häkchen der offenen DoD-Punkte erst nach Closure-Notiz, Register und Risiko-Ausgängen setzen.

Restpunkte: V-2 (Risiko 3 und 4 veraltet), V-3 (Spec-Folge-Slice trägt die Hostliste).
