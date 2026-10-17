# Slice v1: Devcontainer-Feature `sshd` im eingebauten Katalog

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md) — dritter Slice der Welle; ihm gehen [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md) und [`slice-v1-cli-quiet-sicherheitswarnungen`](../done/slice-v1-cli-quiet-sicherheitswarnungen.md) voraus, ihm folgen [`slice-v1-devcontainer-feature-claude-code`](../open/slice-v1-devcontainer-feature-claude-code.md) und [`slice-v2-spec-devcontainer-sshd-claude-code`](../open/slice-v2-spec-devcontainer-sshd-claude-code.md).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Devcontainer-Features), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Hosts je Feature), [`LH-FA-DEV-004`](../../../../spec/lastenheft.md#lh-fa-dev-004--benutzerrechte) (Benutzer, unter dem der Login läuft); [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) (Katalog-Regel: offizielle Einträge brauchen das Präfix `ghcr.io/devcontainers/features/`). Muster: [`slice-v1-devcontainer-feature-python`](../done/slice-v1-devcontainer-feature-python.md).

**Berührte Spec-Stellen:** — (keine). [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) und [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) bleiben bis [`slice-v2-spec-devcontainer-sshd-claude-code`](../open/slice-v2-spec-devcontainer-sshd-claude-code.md) ohne den Schlüssel; das Lastenheft ändert nur der CR.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Schlüssel `sshd` ist ein eingebautes Devcontainer-Feature und rendert das offizielle Feature `ghcr.io/devcontainers/features/sshd:1`; im Sandbox-Profil ist `sshd` erlaubt (Entscheidung des Projektinhabers am 2026-10-17), und die Befehlsausgabe weist den Login-Dienst (Passwort- und Schlüssel-Login, Root-Login erlaubt) mit einer Warnung aus; die Warnung ist eine Sicherheitslockerung im Sinn der Ausnahme 1 von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Entscheidung des Projektinhabers am 2026-10-17, festgehalten im CR): sie trägt `SecurityRelaxation` und bleibt unter `--quiet` sichtbar (erzeugt wird sie von `generate devcontainer`; der CLI-Filter hält sie bei `init`, `generate` und `up`); den Inhalt legt der CR fest, die englische Formulierung der Meldung der Slice. Ein sshd-Port wird nicht automatisch weitergeleitet.

**Planungsgrundlage (WebFetch 2026-10-17, im Slice gegen `install.sh` und `devcontainer-feature.json` erneut zu prüfen, nicht zu übernehmen):** Das Feature lauscht auf
Port 2222; Optionen `gatewayPorts` (Default `no`) und `version` (unbenutzt); die README beschreibt nur Passwort-Login (der Nutzer setzt es mit `sudo passwd $(whoami)`), die Messung in
§3 fand zusätzlich Schlüssel-Login und erlaubten Root-Login; es verlangt eine Debian- oder Ubuntu-Basis. Das Basis-Image des erzeugten Dockerfiles ist `mcr.microsoft.com/devcontainers/base:debian`
(`templates/devcontainer/Dockerfile.tmpl`, [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)); die Voraussetzung ist damit erfüllt, solange das Basis-Image nicht wechselt.

**Gefundener Stand im Code (aus dem Quelltext gelesen):**

- Katalog: `featureCatalogue()` in `devcontainer_features.go` (Felder `source`, `defaultVersion`, `shortDesc`, `publisher`); ein Eintrag mit leerem `publisher` ist
  offiziell (Namensraum `ghcr.io/devcontainers/features`) und löst keine Warnung aus; ein neuer Schlüssel braucht sonst keine Codeänderung ([`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md)).
- Egress: `egressFeatureHosts()` in `devcontainer_egress.go` ist eine Karte Schlüssel → Hosts; ein Schlüssel ohne Zeile trägt keinen Host bei (`egressHosts`).
- `forwardPorts` stammt aus den aktiven Services (`collectDevcontainerForwardPorts` in `generate.go`, [`LH-FA-DEV-005`](../../../../spec/lastenheft.md#lh-fa-dev-005--ports)); die Prüfung `devcontainer.forwardPorts.consistency`
  warnt nur bei **fehlenden** Service-Ports (`portsNotForwarded` in `doctor.go`), ein zusätzlicher Port stört sie nicht.
- Warnungen im Sandbox-Zweig entstehen in `sandboxRelaxationWarnings` (`devcontainer_sandbox.go`, Feld `SecurityRelaxation` bestimmt die Sichtbarkeit unter `--quiet`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Optionen `gatewayPorts` und `version`** — der Katalog kennt keine Optionen; ein Optionsmodell ist ein Re-Evaluierungs-Trigger von [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md). Der Nutzer überschreibt den Tag über `devcontainer.features.sshd.version` (der Wert ist der OCI-Tag, nicht die SSH-Version).
- **Schlüssel-Login, `authorized_keys`, Passwort-Setzen durch `u-boot`** — das Feature kennt keinen Schlüssel-Login, und `u-boot` setzt kein Passwort; eine solche Funktion wäre eine neue Zusage ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): `u-boot` setzt keinen Agenten-Modus und verwaltet keine Zugänge).
- **Port 2222 in `forwardPorts` eintragen** — `forwardPorts` hat nur Service-Ports als Quelle ([`LH-FA-DEV-005`](../../../../spec/lastenheft.md#lh-fa-dev-005--ports)); ein Eintrag wäre eine Vertragsänderung und gehört in den CR.
- **Die Sandbox-Zusage selbst** — sie legt der CR fest; dieser Slice setzt um, was dort steht.
- **Die Spezifikation** ([`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)) — sie folgt in [`slice-v2-spec-devcontainer-sshd-claude-code`](../open/slice-v2-spec-devcontainer-sshd-claude-code.md), nachdem Lastenheft und Code es tragen.

## 2. Definition of Done

- [x] **Katalog und Generator:** `sshd` im Katalog (offizielle Quelle `ghcr.io/devcontainers/features/sshd`, Standardversion `1`, ohne `publisher`); Egress-Zeile nur, wenn die Messung Laufzeit-Hosts findet; das Sandbox-Verhalten folgt dem CR (erlaubt, mit Warnung, die `SecurityRelaxation` trägt); Golden Cases belegen Render, Sandbox-Verhalten, Hosts und double-run -> NoOp; ein Test belegt, dass die Warnung unter `--quiet` sichtbar bleibt (Happy; der Erzeuger ist `generate devcontainer`, der Test prüft den Filter bei `init`, `generate` und `up` mit eingespeistem Eintrag), ohne `sshd` keine Warnung erscheint (Negative) und die Drittanbieter-Warnung eines anderen Features unter `--quiet` entfällt (Boundary); die Mutation „ohne `SecurityRelaxation`“ fällt rot.
- [x] **Messung:** Die Dev-Container-CLI baut ein Image mit `sshd` auf dem Basis-Image (Verfahren der Messreihe in [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md)), im Profil `default` und im Sandbox-Profil mit und ohne Egress; belegt sind: sshd lauscht auf 2222, die Feature-Metadaten (`mounts`, `capAdd`, `privileged`, `entrypoint`, `installsAfter`, Optionen) und die Laufzeit-Hosts. Der Befund steht in der Evidence; findet er eine Lockerung des Profils, geht der Slice zurück an den CR (§4).
- [x] **Doku:** Benutzerhandbuch (Katalogtabelle, Passwort-Setzen mit `sudo passwd`, Sicherheitshinweis zum Sandbox-Profil laut CR), `README.md`, `README.de.md`, `CHANGELOG.md` (`[Unreleased]`) und `harness/replay.md` (Golden-Zeile) nennen den Schlüssel; der Satz „follow in a later release“ im CHANGELOG-Eintrag des CR (Lastenheft 0.11.0) ist für den Teil dieses Slice angepasst (der letzte landende Code-Slice der Welle entfernt ihn).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update, falls ein öffentlicher Vertrag berührt ist.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft, im Repo mit Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/application/devcontainer_features.go` | update | Eintrag `sshd` in `featureCatalogue()` (Quelle, Standardversion `1`, `shortDesc`) |
| `internal/hexagon/application/devcontainer_egress.go` | update (bedingt) | Zeile in `egressFeatureHosts()` nur, wenn die Messung Hosts findet |
| `internal/hexagon/application/devcontainer_sandbox.go`, `generate.go` | update (bedingt) | Sandbox-Regel laut CR: Warnung in `sandboxRelaxationWarnings` mit `SecurityRelaxation` (über `relax`); die Stelle bestimmt der Implementer aus dem Code |
| `internal/hexagon/application/devcontainer_features_test.go`, `generate_features_test.go`, `generate_sandbox_test.go`, `internal/adapter/driving/cli/quiet_test.go` | update | Katalogumfang, Render und Dedupe, Sandbox-Verhalten, Sichtbarkeit der Warnung unter `--quiet`, Egress-Hosts je Feature — nach [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) |
| `internal/hexagon/application/doctor_features_drift_test.go` | update | `sshd` ist projizierbar (`projectAllFeatureEntries`), kein Drift-Fehlalarm |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | update | öffentlicher Vertrag und Golden-Zeile |

Die Pfade stammen aus dem Quelltext (Lesestand 2026-10-17); der Implementer prüft jede Zeile vor der Änderung gegen den Code und trägt Abweichungen in die Evidence ein.

**Umsetzungsstand und Abweichungen von der Tabelle (Implementer, 2026-10-17):**

- `devcontainer_egress.go` und `generate.go` blieben unverändert: die Messung (unten) findet keinen Host, und `sandboxRelaxationWarnings` wird in `generate.go` bereits aufgerufen (`generateDevcontainer`). Neu ist `loginServiceActive` in `devcontainer_sandbox.go`; es prüft die aufgelöste Quelle der aktiven Features gegen den Katalogeintrag `sshd`.
- `harness/replay.md` trägt eine neue Golden-Zeile für `sshd`.
- Die Warnung entsteht nur in `generate devcontainer`: `init` rendert ohne Feature-Konfiguration (`initTemplateData` übergibt `nil`) und `up` erzeugt nur Neuanlege-Warnungen (`upservice.go`). Der Test zu `--quiet` bei `init`, `generate` und `up` belegt deshalb die Sichtbarkeit des Kennzeichens im CLI-Filter (`visibleWarnings` in `cli/warnings.go`) mit einem Eintrag im Wortlaut der Warnung; dass `init` oder `up` den Eintrag selbst erzeugen, belegt er nicht.
- Warnungstext nach dem Review: Er nennt Passwort- und Schlüssel-Login, den erlaubten Root-Login und `sudo passwd <user>` (gemessen: `sshd -T`, siehe Messung); `generate_features_test.go` pinnt „key login“ und „root login“ und das Fehlen von „only“.
- Grenze der Erkennung (`loginServiceActive`): Der Zweig `f.Source` ist ungepinnt, weil kein erreichbarer Fall existiert: ein `source:`-Override braucht ein URL-Schema (`validateFeatureSource`), die Katalogquelle `ghcr.io/devcontainers/features/sshd` hat keines und kann daher nie gleich einem Override sein; ein anderes SSH-Feature über `source` bleibt unausgewiesen.

**Messung (Dev-Container-CLI 0.89.0 im Docker-Container, Colima, Docker 29.x, arm64, Basis `mcr.microsoft.com/devcontainers/base:debian` = Debian trixie; Feature `sshd:1` löste auf Version 1.1.0 auf; Projekt mit dem erzeugten Sandbox-Profil von `u-boot`):**

| Prüfung | Ergebnis |
|---|---|
| Feature-Metadaten (`devcontainer-feature.json`, OCI-Manifest) | `entrypoint: /usr/local/share/ssh-init.sh`, `installsAfter: common-utils`; **kein** `mounts`, `capAdd`, `privileged`, `securityOpt`; Optionen `version` (ungenutzt) und `gatewayPorts` |
| Container (`docker inspect`) in `default`, `sandbox`, `sandbox` mit Egress | `Privileged=false`; `CapAdd` leer, mit Egress nur `NET_ADMIN` (von `u-boot`); keine Devices, keine Port-Bindungen; einziger Mount das Workspace-Volume (Sandbox); im Lauf `default` ein Mount (der Workspace, Typ nicht geprüft) |
| Port | `ss -ltn`: `0.0.0.0:2222` und `[::]:2222` in allen drei Läufen; mit Egress zusätzlich der lokale Resolver auf `127.0.0.1:53` |
| sshd-Konfiguration (`sshd -T`) | `passwordauthentication yes`, `pubkeyauthentication yes`, `permitrootlogin yes`, `permitemptypasswords no`, `gatewayports no` |
| Konten | `passwd -S`: `root` und `vscode` gesperrt (`L`); `~/.ssh` fehlt |
| Netz beim Build | `apt-get` über `deb.debian.org` (Basis-Host der Allowlist); das Feature selbst holt die Dev-Container-CLI von `ghcr.io` auf dem Host, nicht im Container |
| Netz zur Laufzeit | `ssh-init.sh` startet nur `/etc/init.d/ssh`; kein Netzaufruf; Lauf mit Egress-Beschränkung erfolgreich, Port 2222 lauscht |

Folgerungen: keine Egress-Zeile; das Feature unterläuft das Profil nicht (Risiko „Feature-Metadaten“ entfällt); die CR-Formulierung „nur Passwort-Login“ ist enger als der Befund: die Schlüssel-Anmeldung ist in der sshd-Konfiguration eingeschaltet und wirkt, sobald der Nutzer selbst einen Schlüssel einträgt, und Root-Login ist erlaubt. `u-boot` legt keine Schlüssel an und setzt kein Passwort; das Handbuch nennt beides.

**Grenzen der Messung:** (1) Ein tatsächlicher SSH-Login mit gesetztem Passwort wurde nicht ausgeführt; belegt sind Port, Konfiguration und gesperrte Konten. (2) Die Netz-Hosts des Builds stammen aus dem Build-Protokoll der CLI, nicht aus einem DNS-Mitschnitt; zur Laufzeit folgt „kein Netzaufruf“ aus dem gelesenen `ssh-init.sh` und dem erfolgreichen Egress-Lauf. (3) Nur arm64 und Colima; das Feature `sshd:1` ist ein Major-Pin, spätere Versionen können abweichen. (4) Der Feature-Quelltext kam über das OCI-Artefakt von `ghcr.io/devcontainers/features/sshd`, nicht über das Repository.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der CR [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md) liegt in `done/`; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice überschreitet 3 Liefer-Punkte.
- `in-progress` → `open` (blockiert): die Messung findet eine Lockerung des Profils, die der CR nicht trägt (dann zuerst CR erweitern), oder `ghcr.io/devcontainers/features/sshd:1` baut auf dem Basis-Image nicht.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **`sshd` öffnet im Container einen Login-Dienst** (gemessen: Passwort- und Schlüssel-Login, Root-Login erlaubt; Konten gesperrt, bis der Nutzer ein Passwort setzt oder einen Schlüssel einträgt), im Sandbox-Profil mit Zielsetzung Schadensbegrenzung; entschieden ist: erlaubt mit Warnung, die unter `--quiet` sichtbar bleibt. — **Ausgang:** entfallen: die Entscheidung des CR ist umgesetzt (DoD 1) und bei der Closure bestätigt; die Warnung nennt den gemessenen Umfang (Passwort- und Schlüssel-Login, Root-Login), Text und Doku folgen dem Befund (Review F-1, F-2, behoben in `0574231`).
- **Feature-Metadaten unterlaufen das Profil** (`mounts`, `capAdd`, `privileged`, `entrypoint`): u-boot erzeugt nur `devcontainer.json`; die Features bringen eigene Metadaten mit, die der Dev-Container-Lauf zusammenführt. — **Ausgang:** entfallen: gemessen (Slice §3): kein `mounts`, `capAdd`, `privileged`, `securityOpt`, nur `entrypoint` und `installsAfter`; gemessen an `sshd` 1.1.0, der Pin `1` schwimmt, spätere Versionen können abweichen.
- **Der Port 2222 ist im Sandbox-Profil nur über das Forwarding der Dev-Container-Werkzeuge erreichbar**, `u-boot` veröffentlicht ihn nicht. — **Ausgang:** entfallen: `forwardPorts` bleibt unberührt, und `doctor` warnt nur bei fehlenden Service-Ports (aus `doctor.go` gelesen); die Handbuch-Zeile nennt das.
- **Plan-Aussage über den Code stimmt nicht** (Klasse „Plan-Aussage über Code stimmt nicht“, siehe §8): die Fundstellen in §1 und §3 sind ein Lesestand. — **Ausgang:** weiter offen: → [BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) im Register (verkörpert in `AGENTS.md §3.10`); der Implementer prüft jede Aussage vor der Änderung im Code, Abweichungen stehen in der Evidence.

## 7. Closure-Notiz

**Belege:** Review-Report [`2026-10-17-slice-v1-devcontainer-feature-sshd`](../../../reviews/2026-10-17-slice-v1-devcontainer-feature-sshd.md) (0 HIGH, 2 MEDIUM, 4 INFO; F-1 Warnungstext zu eng und F-2 Doku zu eng behoben in `0574231`; F-3 Befehls-Formulierung berichtigt; F-4 Grenze (Erkennungszweig `f.Source` ungepinnt); F-5 und F-6 Grenze und Folge-Arbeit). Verification Evidence [`2026-10-17-verification-slice-v1-devcontainer-feature-sshd`](../../../reviews/2026-10-17-verification-slice-v1-devcontainer-feature-sshd.md): verifiziert, die Auflagen 1 bis 5 sind hier und in der Welle erfüllt. `make gates` grün, Coverage 91,10 % (Schwelle 90 %); `make doc-trace`: 83 Anforderungen, 0 Waisen. Mutationen: fünf des Implementers, vier des Reviewers, drei des Verifiers, die Wächter-Mutationen rot (Ausnahme: die Mutation „Quelle ignorieren“ blieb grün, siehe Review F-4 und Grenzen).

- **Was hat funktioniert:** Die Messung in der Dev-Container-CLI belegte, dass das Feature das Profil nicht unterläuft (`Privileged=false`, keine Capabilities, Devices oder Port-Bindungen von `sshd`); der Rückführungs-Trigger aus §4 trat nicht ein. Die Warnung ist eine [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)-Ausweisung mit `SecurityRelaxation` und bleibt unter `--quiet` sichtbar ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Ausnahme 1), die Drittanbieter-Warnung daneben entfällt weiter; der Boundary-Test pinnt beides. Die Regel `AGENTS.md` §3.10 fand die falsche Plan-Aussage zum Erzeuger der Warnung: der Implementer trug sie vor der Änderung in den Umsetzungsstand ein.
- **Was ging anders als geplant:** Der Planungsgrundlage und dem CR („nur Passwort-Login“) widersprach die Messung (`sshd -T`: `passwordauthentication yes`, `pubkeyauthentication yes`, `permitrootlogin yes`); Warntext und Doku waren zunächst zu eng (Review F-1, F-2, Klassen „Warnungstext enger als gemessenes Verhalten“ und „Doku widerspricht Messbefund und eigenem Abschnitt“) und stehen jetzt auf dem Befund. Die Plan-Aussage „`init`, `generate` und `up` erzeugen die Warnung“ traf nur auf `generate devcontainer` zu (Review F-3, Klasse „Warnungs-Erzeuger enger als Befehlsliste in der Doku“); `init` rendert ohne Feature-Konfiguration, `up` erzeugt nur Neuanlege-Warnungen. Die übrigen Klassen des Reviews: „Erkennungszweig ohne Wächter“ (F-4), „Folgerung weiter als Messbereich“ (F-5), „Spec-Katalogliste hinter dem Code“ (F-6).
- **Steering-Loop-Eintrag:** **Gezählt, nicht verkörpert.** Die Klasse „Plan-Aussage über Code stimmt nicht“ trat erneut auf: zwei Befunde in einem Vorgang (Passwort-Login; Erzeuger der Warnung), Beleg für [BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md); die Regel `AGENTS.md` §3.10 existiert (Stand `verkörpert`), ein neuer Zielort entsteht nicht. **Lehre:** Die Planungsgrundlage stammte aus einem WebFetch, den die Messung widerlegte; eine Aussage über Fremdcode ist erst nach der Messung eine Aussage, davor ein Lesestand. Das Mittel war die Messung, die der Slice ohnehin vorsah (DoD 2): sie fand den Widerspruch, bevor die Doku ihn trug, und die Review fand, was die Messung zwar notierte, aber nicht in Text und Doku zog. Eine Entscheidung zu `state.md` trifft der Projektinhaber; sie bleibt unverändert.
- **Grenzen:** (1) **Die Messung ist im Repo nicht wiederholbar:** kein Make-Target, kein Skript; sie lief mit der Dev-Container-CLI 0.89.0 in einem Wegwerf-Image mit Docker-Socket, mit dem Feature `sshd` 1.1.0 (der Pin `1` schwimmt), dem Basis-Image `mcr.microsoft.com/devcontainers/base:debian` (Debian trixie) und den drei Profilen (default, sandbox, sandbox mit Egress), auf arm64 unter Colima; der Feature-Quelltext kam als OCI-Artefakt, nicht aus der Repository-Quelle. Ein Nachlauf braucht genau diese Angaben. [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) deckt diese Art Messung als Grenze („im Repo nicht ausgeführt“); der Test `SandboxEgress_SshdAddsNoHost` hält die Folge (kein Host) fest, eine abweichende spätere Messung würde als Test-Änderung sichtbar. (2) Kein echter SSH-Login: der Schlüssel-Login folgt aus `sshd -T`, nicht aus einem Login. (3) Der schwimmende Pin `1`: spätere Feature-Versionen können Metadaten und Verhalten ändern; „Feature-Metadaten unterlaufen das Profil“ ist für 1.1.0 entfallen. (4) Der Zweig `f.Source` in `loginServiceActive` ist ungepinnt und nicht erreichbar (ein Override braucht ein URL-Schema, die Katalogquelle keines); ein anderes SSH-Feature über `source` bleibt unausgewiesen. (5) Das Lastenheft nennt „(Passwort-Login)“ und untertreibt den gemessenen Zustand; der breitere Warntext erfüllt die Zusage, eine Präzisierung bei der nächsten Berührung des Lastenhefts ist Sache des Projektinhabers. (6) Port 2222 lauscht im Container auf allen Schnittstellen (`0.0.0.0` und `[::]`), wird aber nicht weitergeleitet und nicht veröffentlicht. (7) Der CLI-Test zu `--quiet` belegt für `init`, `generate` und `up` nur den Filter mit eingespeistem Eintrag; Erzeuger ist allein `generate devcontainer`. (8) [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) und [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) führen `sshd` bis zum Spec-Slice nicht; bis dahin gilt Lastenheft vor Spezifikation. (9) `make test-docker` und `make verify-depguard` liefen nicht: kein Docker-Laufzeitpfad und keine Importregel berührt.
- **Beobachtungs-Register (`../observations/`):** Evidence-Datei `BEO-DOK/plan-aussage-ueber-code/evidence/slice-v1-devcontainer-feature-sshd.md` angelegt (zwei Befunde: „nur Passwort-Login“; Erzeuger der Warnung); der Zähler steht damit bei **fünf** Vorgängen. `state.md` unverändert (die Regel ist verkörpert in `AGENTS.md` §3.10). `BEO-CLI/kommentar-beschreibt-abwesendes`: keine Beobachtung für diesen Slice (`git diff` auf die Muster (a) bis (c) gelesen, kein Treffer). <!-- d-check:ignore (Pfad-Zitat im Code-Span) -->
- **Folge-Slices:** [`slice-v1-devcontainer-feature-claude-code`](../open/slice-v1-devcontainer-feature-claude-code.md) (entfernt als letzter landender Code-Slice den CHANGELOG-Restsatz „follow in a later release“) und [`slice-v2-spec-devcontainer-sshd-claude-code`](../open/slice-v2-spec-devcontainer-sshd-claude-code.md) (trägt `sshd` in [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) und [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), Review F-6) — Dateien in `open/`.
- **Risiken aus §6:** alle vier tragen einen kanonischen Ausgang: entfallen (Login-Dienst: Entscheidung des CR umgesetzt und mit gemessenem Umfang ausgewiesen; Feature-Metadaten: gemessen, `sshd` 1.1.0, Pin schwimmt; Port 2222) und weiter offen (Plan-Aussage über Code, Register; im Slice eingetreten, siehe Steering-Loop-Eintrag).
- **Drei Paarungen:** Anker (kein `liegt in`, nichts verkörpert) · Folge-Slice (zwei Dateien in `open/`) · Register (Evidence-Datei angelegt, Zähler fünf) — nach dem `git mv` nach `done/` zu prüfen.


## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `internal/hexagon/application` (Kürzel APP) und, für die Doku, `docs/user/`, `harness/` und `spec/` (DOK) — beide Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006) (Abschnitt Modus-Deklaration), beide erfüllen die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. `internal/adapter/driving/cli` (CLI) wird berührt, soweit der Test unter `--quiet` dort liegt (die Sichtbarkeit selbst folgt dem Feld `SecurityRelaxation`, Filter in [`slice-v1-cli-quiet-sicherheitswarnungen`](../done/slice-v1-cli-quiet-sicherheitswarnungen.md)); dann gilt `AGENTS.md §3.9` (Kommentar-Selbstcheck); `AGENTS.md §3.10` gilt für jede Plan-Aussage über Code.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer für APP und DOK. Die einzige Beobachtung `BEO-CLI/kommentar-beschreibt-abwesendes` (3×, verkörpert in `AGENTS.md §3.9`) betrifft CLI; sie wird nur bei einer CLI-Berührung relevant. Die Klasse „Plan-Aussage über Code stimmt nicht“ steht als [BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) im Register (vier Belege, verkörpert in `AGENTS.md §3.10`): jede Aussage über den Code wird vor der Änderung an der Stelle gelesen; Abweichungen stehen in der Evidence.

Alle berührten Sub-Areas GF.
