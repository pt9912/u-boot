# Slice v1: Katalog-Schlüssel `ollama` — Installation, Hosts, Modellspeicher, Start, Tests, Doku

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) — siebter Slice der Welle; ihm gehen [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md) (Quelle, Messwerte, Folge-ADR) und [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) (Lastenheft) voraus; ihm folgt [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Devcontainer-Features; um `ollama` erweitert durch den CR), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Lieferkette), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts des Modell-Pull), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Warnungen unter `--quiet`); [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) und die Folge-ADR des Mess-Slice (Quelle, Pin, Prüfung).

**Berührte Spec-Stellen:** [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (Katalog) und [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Hosts) bleiben bis [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md) ohne den Schlüssel; das Lastenheft ändert nur der CR.

**Verantwortlich:** pt9912

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Schlüssel `ollama` ist ein eingebautes Katalog-Feature: Mit `devcontainer.features.ollama.enabled: true` installiert das erzeugte Devcontainer-Ergebnis Ollama (CPU-only) aus der Quelle, die die Folge-ADR des Mess-Slice wählt, mit der dort entschiedenen Prüfung; die Default-Allowlist trägt die **gemessenen** Laufzeit-Hosts des Modell-Pull; die Modelle liegen in einem benannten Volume (soweit gemessen); der Server startet so, wie die ADR es festlegt; ein Golden Case belegt das Ergebnis, ein Docker-Test den Lauf. Zugesagt ist nur, was der CR und die Messung tragen.

### Gefundener Stand im Code (gelesen am 2026-10-18; die Zeilen für Variante A oder B/E gelten erst nach der ADR)

- **Katalog:** `featureCatalogue()` (`internal/hexagon/application/devcontainer_features.go`) enthält heute `git`, `docker-cli`, `node`, `java`, `go`, `cpp`, `kubectl-helm`, `postgres-client`, `python`, `homebrew`, `sdkman`, `sshd` und `claude-code`; kein `ollama`. Das Hinzufügen eines OCI-Features (Variante A) ist laut Kommentar zu `featureCatalogueEntry` „drop a new entry into [featureCatalogue], done — no other code change“ (für die Render-Seite; Hosts, Warnungen und Tests kommen hinzu). Ein Eintrag mit `publisher` löst `thirdPartyFeatureWarnings` (gleiche Datei) aus.
- **Dockerfile-Block statt Feature (Varianten B, C, E):** `Dockerfile.tmpl` führt Blöcke je Fähigkeit hinter `{{- if .NestedPodman}}`, `.NestedDocker`, `.Egress`; die Schalter sind Felder von `templateData` (`internal/hexagon/application/templates.go`), gesetzt in `devcontainerTemplateData` (`internal/hexagon/application/devcontainer_sandbox.go`) aus der Konfiguration. Ein Katalogschlüssel, der kein OCI-Feature ist, ist im Katalog-Typ heute nicht vorgesehen; `collectDevcontainerFeatures` würde einen Eintrag ohne OCI-Quelle nicht rendern (entschieden in [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 2: Marker statt Quelle).
- **Egress-Hosts je Feature:** `egressFeatureHosts()` (`devcontainer_egress.go`) ist eine Karte Schlüssel → Hosts; `egressHosts` liest sie für jedes Feature mit `Enabled` wahr nach **Schlüsselname**; die Hosts kommen in die Allowlist des Startscripts. Der Eintrag `ollama` fehlt.
- **Volumes und Start:** `devcontainer.json.tmpl` kennt im Sandbox-Zweig drei Fälle für `mounts` (`NestedPodman` mit Workspace- und Containers-Volume, `NestedDocker` mit Workspace- und Docker-Volume, sonst nur das Workspace-Volume); ein zusätzliches Modell-Volume verlangt, diese Verzweigung so umzubauen, dass Volumes sich zusammensetzen lassen (Entwurf, byte-identisch für bestehende Konfigurationen, Golden Cases). `postStartCommand` setzt `devcontainerTemplateData` aus `postStart`: `sudo sh /usr/local/bin/u-boot-egress-init` vor dem Docker-Start-Script; ein Auto-Start von `ollama serve` (Variante O6-ii der Messung) käme hinter dieses Script.
- **Sandbox-Regeln für Features:** Beide bisherigen Feature-Schlüssel mit Sandbox-Regel erzeugen nur den Schlüssel unter `features`; kein `runArgs`, `mounts` oder Startbefehl ([`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten), Absatz *Sandbox-Regel*). Das Modell-Volume ist das erste Feature-bedingte `mounts` ([`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 8); einen Startbefehl ergänzt `ollama` nicht (Entscheidung 4, kein Auto-Start). Der Satz der *Sandbox-Regel* gilt weiter für `sshd` und `claude-code` und darf `ollama` nicht mitmeinen; die Spezifikation führt das nach.
- **Gegen den Entwurf [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) gelesen (2026-10-18):** `featureFor` entscheidet über den Orphan (`maybeWarnOrphanFeatureActivation` in `config.go`, `classifyFeatureEntries` in `doctor.go`), `ollama` gehört deshalb in den Katalog; `devcontainerTemplateData` setzt die Schalter nur im Sandbox-Zweig (Rückkehr bei `!sandbox`), `USER root` steht in `Dockerfile.tmpl` nur bei Sandbox-Profil oder UID-Überschreibung; `egressHosts` liest `egressFeatureHosts()` nach Schlüsselnamen, die Zeile `ollama` braucht keinen Katalogbezug; `devcontainer.json.tmpl` hat im Sandbox-Zweig drei feste Fälle für `mounts`.
- **Tests der Vorgänger als Muster:** `TestGenerateDevcontainer_SshdAndClaudeCodeTogether` (`internal/hexagon/application/generate_features_test.go`) pinnt das Zusammenspiel zweier Features (Render-Schlüssel genau einmal, zweiter Lauf NoOp, Egress-Vereinigung); `TestE2E_SandboxDevcontainer_Egress` und `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` (`internal/e2e/sandbox_devcontainer_docker_test.go`) sind die Docker-Tests für Egress und den verschachtelten Build; der Katalog-Test prüft, dass jeder Eintrag offiziell ist oder seinen Herausgeber als Namensraum nennt (`devcontainer_features_test.go`, Fitness Function von [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md)).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Messung und Quellenwahl** — [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md); der Slice setzt die ADR um, er wiederholt die Messung nicht (er prüft den Lauf, den die Messung belegt hat, in einem Docker-Test).
- **Anschluss an den Ollama des Hosts, GPU, Modellverwaltung, Feature-Optionen** (`pull`) — vom Projektinhaber ausgeschlossen bzw. Re-Evaluierungs-Trigger von [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md); der Nutzer holt Modelle im Container.
- **Weitere Agenten** — kein Agent außer Ollama genannt.
- **Portweiterleitung des Servers nach außen** (`forwardPorts`) — [`LH-FA-DEV-005`](../../../../spec/lastenheft.md#lh-fa-dev-005--ports) lässt `forwardPorts` aus den Services entstehen; eine Änderung wäre eine Plan-Änderung (Muster `sshd`).
- **Die Spezifikation** ([`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)) — sie folgt in [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md).

## 2. Definition of Done

- [ ] **Katalog und Generator:** `ollama` ist wählbarer Schlüssel; Quelle, Version und Prüfung der Lieferkette sind genau die der Folge-ADR (gepinnt, SHA-256-Prüfung vor dem Entpacken; die Warnung ist eine eigene Meldung nach [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 3, keine Drittanbieter-Warnung nach [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md)); `egressFeatureHosts()` trägt die zwei gemessenen Laufzeit-Hosts; Modell-Volume nach Entscheidung 8, kein Auto-Start (Entscheidung 4); bestehende Konfigurationen rendern byte-identisch; **ein Golden Case** aktiviert `ollama` mit Sandbox-Profil und Egress (auch zusammen mit `claude-code`): `devcontainer.json`, Dockerfile und `egress-init.sh` stabil, zweiter Lauf NoOp, Egress-Liste die sortierte, duplikatfreie Vereinigung, Hinweise genau wie die ADR vorschreibt.
- [ ] **Lauf belegt:** Ein Docker-Test (`//go:build docker`) baut das erzeugte Ergebnis, startet den Container hinter der Egress-Restriktion und belegt das, was der CR zusagt (zum Beispiel `ollama --version`, Server erreichbar auf der gemessenen Adresse, ein Pull-Host erreichbar und ein fremder Host gesperrt); jede Bedingung in Handbuch und CHANGELOG hat einen Test oder steht dort nicht (`AGENTS.md` §3.10). Läuft der Test in der Umgebung nicht, nennt die Closure den Grund (kein stilles Rot).
- [ ] **Pins aus der Verification des CR** ([`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) §7, Verification-Report V-3 bis V-5): (a) ein Test, dass die Warnung zur Lieferkette und die Meldung „keine Wirkung ohne Sandbox-Profil“ unter `--quiet` entfallen — oder das Gegenteil, falls der Projektinhaber das Gate zu „weist aus“ gegen `--quiet` anders entscheidet; die Entscheidung steht vor dem Test fest; (b) ein Golden Case „kein `OLLAMA_HOST` in `devcontainer.json`, Dockerfile und Skript, kein `forwardPorts`“ und ein Docker-Test „Aufruf über die Container-Adresse scheitert, `127.0.0.1:11434` antwortet“ (Bedingung „nur Loopback“ des Lastenhefts); (c) ein Negativfall „ohne Sandbox-Profil: kein Block im Dockerfile und die ungekennzeichnete Meldung“ (Bedingung „wirkt im Sandbox-Profil“). Ohne Test steht die Bedingung nicht in Handbuch und CHANGELOG (`AGENTS.md` §3.10).
- [ ] **Doku:** Handbuch (Katalogtabelle, Egress-Tabelle, ein Abschnitt *Ollama (optional)* mit Messstand, Grenzen und den Schritten „Modell holen, Server nutzen“), `README.md`, `README.de.md`, `CHANGELOG.md` (`[Unreleased]`) und `harness/replay.md` (Golden-Zeile) nennen den Schlüssel; Aussagen zu den Hosts stehen in Handbuch, CHANGELOG und (nach dem Spezifikations-Slice) Spezifikation gleich.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update, falls ein öffentlicher Vertrag berührt ist (hier ja).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen von der Welle-Closure geprüft.

## 3. Plan (vor Code)

Die Pfade stammen aus dem Quelltext (Lesestand 2026-10-18); der Implementer prüft jede Zeile vor der Änderung gegen den Code und trägt Abweichungen in die Evidence ein. **Stand gegen den Entwurf [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (`Accepted` am 2026-10-18, Freigabe-Punkte a bis e wie vorgeschlagen):** Variante B — Dockerfile-Block, Katalogeintrag ohne OCI-Quelle, nur mit Sandbox-Profil, kein Auto-Start; die Zeilen „A“ entfallen, die Zeilen sind darauf berichtigt und am Symbol gelesen (2026-10-18), gegen den Accepted-Stand erneut gelesen (die ADR ist seit dem Entwurf inhaltlich unverändert, nur Status, Index und Geschichte). Jede weitere Änderung der Entscheidung ist eine neue ADR (`AGENTS.md` §3.5), danach ist die Tabelle erneut zu berichtigen.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/application/devcontainer_features.go` | update | Eintrag `ollama` in `featureCatalogue()` mit einem Marker (neues Feld in `featureCatalogueEntry`) statt `source`, `defaultVersion`, `publisher`; `projectFeatureSourceVersion` liefert für ihn ohne `source:` `ok=false` (damit rendern `collectDevcontainerFeatures` und `projectFeatureEntry` nichts, `projectAllFeatureEntries` bildet keinen Schlüssel `:<version>`, `checkDevcontainerFeaturesDrift` bleibt still); `validateDevcontainerFeatures`: `version:` ohne `source:` bei `ollama` → Exit-Code `10`; `thirdPartyFeatureWarnings` bleibt unberührt (hängt an `publisher`) |
| `internal/hexagon/application/devcontainer_egress.go` | update | Zeile `ollama` in `egressFeatureHosts()` aus der Messung |
| `internal/hexagon/application/templates.go` | update | Felder in `templateData` für den Schalter und den Volume-Namen (der Start entfällt: kein Auto-Start) |
| `internal/hexagon/application/devcontainer_sandbox.go` | update | `devcontainerTemplateData`: Schalter und Volume-Name im Sandbox-Zweig (kein Anteil an `postStart`); `sandboxRelaxationWarnings`: die ungekennzeichnete Meldung „keine Wirkung ohne Sandbox-Profil“ (Muster `nestedRuntime`); `catalogueFeatureActive` vergleicht die aufgelöste **Quelle** und trägt den Typ nicht, es wird nicht verwendet |
| `internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl` | update | Block hinter `USER root` im Sandbox-Zweig: Pin (`ARG`), `sha256sum -c` vor `tar`, `--exclude`-Muster, `zstd` nur im Bau, `mkdir`/`chown` von `/home/vscode/.ollama` (Muster Block `NestedPodman`) |
| `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` | update | zusammensetzbare `mounts` für das Modell-Volume, byte-identisch für bestehende Konfigurationen |
| `internal/hexagon/application/generate.go` | update | die Warnung zur Lieferkette (eigene Meldung mit Code [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), ohne `SecurityRelaxation`, entfällt unter `--quiet`) neben dem Aufruf von `thirdPartyFeatureWarnings` und `agentNodeWarnings` |
| `internal/hexagon/application/devcontainer_features_test.go`, `export_test.go`, `generate_features_test.go`, `generate_sandbox_test.go` | update | Pins V-3 bis V-5 (DoD), Katalogumfang (`TestFeatureCatalogue_EntriesAreValid` verlangt heute `Source` und `DefaultVersion` je Eintrag und die Namensraum-Regel: für den Typ Dockerfile-Block anzupassen; `FeatureCatalogueEntryForTest` in `export_test.go` trägt den Marker), Render, Hosts, Volume, Start, **Golden Case** mit `ollama` (auch zusammen mit `claude-code`), double-run → NoOp — nach [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) |
| `internal/e2e/sandbox_devcontainer_docker_test.go` | update | Docker-Test: Lauf hinter Egress (Build-Tag `docker`), Bedingungen aus dem CR |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | update | öffentlicher Vertrag und Golden-Zeile |
| `docs/plan/planning/observations/BEO-APP/…` | entfällt | der zweite Beleg zur Pin-Alterung liegt aus dem Mess-Slice vor; hier nur bei einem dritten Auftreten |
| `docs/reviews/` (Review-Report, Verification Evidence) | neu | Review und Verifikation |
| diese Slice-Datei | update | Umsetzungsstand, Closure (§7) |

Abgleich der Tabelle gegen `git diff --stat` nach der **letzten** Fix-Runde (`AGENTS.md` §3.10); jede Zeile mit „bedingt“ wird als eingetreten oder entfallen markiert.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md) und [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) liegen in `done/` (die ADR ist `Accepted`, Lastenheft mit freigegebenem Wortlaut); kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice überschreitet 3 Liefer-Punkte, etwa weil der Umbau der `mounts`-Verzweigung oder ein zweiter Katalog-Eintragstyp (B/E) einen eigenen Slice verlangt.
- `in-progress` → `open` (blockiert): die Umsetzung verlangt eine Lockerung des Profils (Capability, Device, Mount), die ADR und CR nicht tragen (dann zuerst Entscheidung, wie bei [`slice-v1-devcontainer-podman-build`](../done/slice-v1-devcontainer-podman-build.md)); die gepinnte Quelle ist nicht mehr abrufbar.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). Verifier: Stichprobe der Zusagen aus dem CR gegen den Docker-Test; Validator-Hinweis: der Nutzerpfad (`config set … enabled true`, `generate devcontainer`, Container starten, Modell holen) ist nachzuspielen.

## 6. Risiken und offene Punkte

- **Der Katalog-Typ trägt einen Eintrag ohne OCI-Quelle nicht** (B/E). — **Ausgang:** offen bis zur Umsetzung; [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 2 legt den Marker-Typ fest (kein Render-Schlüssel, Drift-Prüfung still, kein Orphan); trägt der Code ihn nicht, Rückführung nach `next/` mit eigenem Slice für den Katalog-Typ.
- **Der Umbau der `mounts`-Verzweigung ändert die Ausgabe bestehender Konfigurationen.** — **Ausgang:** offen bis zur Closure; erwartet: entfallen (Golden Cases der bisherigen Fälle bleiben byte-identisch).
- **Der Docker-Test ist langsam, groß (Image, Download) und netzabhängig**; auf einem Runner ohne Cache rot ([`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md), [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md)). — **Ausgang:** offen; der Test trägt das Build-Tag, die Closure nennt Lauf oder Grund; kein Workflow-Eingriff (Gate-Änderung, `AGENTS.md` §3.6).
- **Der Pin altert** (Ollama veröffentlicht häufig). — **Ausgang:** weiter offen: Trigger der ADR und Register-Eintrag (Muster `podman-static`).
- **Eine Bedingung in Handbuch oder CHANGELOG hat keine Entsprechung im Code** („Server nur auf Loopback“, „Modelle überstehen den Neubau“). — **Ausgang:** offen bis zur Closure; jede Bedingung hat einen Test oder entfällt ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), `AGENTS.md` §3.10).
- **Hosts des Modell-Pull wechseln** (Blob-Speicher unter Subdomains). — **Ausgang:** offen; Rückgabe an [`slice-v1-egress-wildcard-pruefung`](../done/slice-v1-egress-wildcard-pruefung.md) laut Welle §5, nicht still in die Default-Liste.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: umgesetzte Variante, Fundstellen je Aussage, Lauf des Docker-Tests, Steering-Loop-Eintrag, Register, Folge-Slices, Risiken.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/application` (APP), `internal/e2e` (E2E), `docs/`, `harness/` (DOK) — alle Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); jede erfüllt die Schwelle (zwei von drei Achsen) laut Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-APP/podman-static-pin-alterung`](../observations/BEO-APP/podman-static-pin-alterung/observation.md) (offen, zwei Belege; der Ollama-Pin ist der zweite gepinnte Download), [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) und [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) (offen, je ein Beleg; neuer Docker-Test), [`BEO-WZG/kopplung-nicht-nachgezogen`](../observations/BEO-WZG/kopplung-nicht-nachgezogen/observation.md) (offen, zwei Belege; Katalogtabelle, Egress-Tabelle, README, Golden-Zeile sind zusammen nachzuziehen), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (verkörpert in `AGENTS.md` §3.10). Evidenz-/Diskrepanz-Risiko: mittel.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
