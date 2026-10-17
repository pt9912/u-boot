# Slice v1: Devcontainer-Feature `sshd` im eingebauten Katalog

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md) — zweiter Slice der Welle; ihm geht [`slice-cr-devcontainer-agent-features-sandbox`](../in-progress/slice-cr-devcontainer-agent-features-sandbox.md) voraus, ihm folgen [`slice-v1-devcontainer-feature-claude-code`](slice-v1-devcontainer-feature-claude-code.md) und [`slice-v2-spec-devcontainer-sshd-claude-code`](slice-v2-spec-devcontainer-sshd-claude-code.md).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Devcontainer-Features), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Hosts je Feature), [`LH-FA-DEV-004`](../../../../spec/lastenheft.md#lh-fa-dev-004--benutzerrechte) (Benutzer, unter dem der Login läuft); [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) (Katalog-Regel: offizielle Einträge brauchen das Präfix `ghcr.io/devcontainers/features/`). Muster: [`slice-v1-devcontainer-feature-python`](../done/slice-v1-devcontainer-feature-python.md).

**Berührte Spec-Stellen:** — (keine). [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) und [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) bleiben bis [`slice-v2-spec-devcontainer-sshd-claude-code`](slice-v2-spec-devcontainer-sshd-claude-code.md) ohne den Schlüssel; das Lastenheft ändert nur der CR.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Schlüssel `sshd` ist ein eingebautes Devcontainer-Feature und rendert das offizielle Feature `ghcr.io/devcontainers/features/sshd:1`; im Sandbox-Profil ist `sshd` erlaubt (Entscheidung des Projektinhabers am 2026-10-17), und die Befehlsausgabe weist den Passwort-Login mit einer Warnung aus; Wortlaut und Sichtbarkeit unter `--quiet` legt der CR fest (offene Teilfrage dort). Ein sshd-Port wird nicht automatisch weitergeleitet.

**Planungsgrundlage (WebFetch 2026-10-17, im Slice gegen `install.sh` und `devcontainer-feature.json` erneut zu prüfen, nicht zu übernehmen):** Das Feature lauscht auf
Port 2222; Optionen `gatewayPorts` (Default `no`) und `version` (unbenutzt); nur Passwort-Login (der Nutzer setzt es mit `sudo passwd $(whoami)`), kein dokumentierter
Schlüssel-Login; es verlangt eine Debian- oder Ubuntu-Basis. Das Basis-Image des erzeugten Dockerfiles ist `mcr.microsoft.com/devcontainers/base:debian`
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
- **Die Spezifikation** ([`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)) — sie folgt in [`slice-v2-spec-devcontainer-sshd-claude-code`](slice-v2-spec-devcontainer-sshd-claude-code.md), nachdem Lastenheft und Code es tragen.

## 2. Definition of Done

- [ ] **Katalog und Generator:** `sshd` im Katalog (offizielle Quelle `ghcr.io/devcontainers/features/sshd`, Standardversion `1`, ohne `publisher`); Egress-Zeile nur, wenn die Messung Laufzeit-Hosts findet; das Sandbox-Verhalten folgt dem CR (erlaubt, mit Warnung; Wortlaut und Sichtbarkeit unter `--quiet` nach der Entscheidung des CR); Golden Cases belegen Render, Sandbox-Verhalten, Hosts und double-run -> NoOp.
- [ ] **Messung:** Die Dev-Container-CLI baut ein Image mit `sshd` auf dem Basis-Image (Verfahren der Messreihe in [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md)), im Profil `default` und im Sandbox-Profil mit und ohne Egress; belegt sind: sshd lauscht auf 2222, die Feature-Metadaten (`mounts`, `capAdd`, `privileged`, `entrypoint`, `installsAfter`, Optionen) und die Laufzeit-Hosts. Der Befund steht in der Evidence; findet er eine Lockerung des Profils, geht der Slice zurück an den CR (§4).
- [ ] **Doku:** Benutzerhandbuch (Katalogtabelle, Passwort-Setzen mit `sudo passwd`, Sicherheitshinweis zum Sandbox-Profil laut CR), `README.md`, `README.de.md`, `CHANGELOG.md` (`[Unreleased]`) und `harness/replay.md` (Golden-Zeile) nennen den Schlüssel.
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
| `internal/hexagon/application/devcontainer_features.go` | update | Eintrag `sshd` in `featureCatalogue()` (Quelle, Standardversion `1`, `shortDesc`) |
| `internal/hexagon/application/devcontainer_egress.go` | update (bedingt) | Zeile in `egressFeatureHosts()` nur, wenn die Messung Hosts findet |
| `internal/hexagon/application/devcontainer_sandbox.go`, `generate.go` | update (bedingt) | Sandbox-Regel laut CR: Warnung in `sandboxRelaxationWarnings` (mit oder ohne `SecurityRelaxation`, je nach Entscheidung des CR); die Stelle bestimmt der Implementer aus dem Code |
| `internal/hexagon/application/devcontainer_features_test.go`, `generate_features_test.go`, `generate_sandbox_test.go` | update | Katalogumfang, Render und Dedupe, Sandbox-Verhalten, Egress-Hosts je Feature — nach [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) |
| `internal/hexagon/application/doctor_features_drift_test.go` | update | `sshd` ist projizierbar (`projectAllFeatureEntries`), kein Drift-Fehlalarm |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | update | öffentlicher Vertrag und Golden-Zeile |

Die Pfade stammen aus dem Quelltext (Lesestand 2026-10-17); der Implementer prüft jede Zeile vor der Änderung gegen den Code und trägt Abweichungen in die Evidence ein.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der CR [`slice-cr-devcontainer-agent-features-sandbox`](../in-progress/slice-cr-devcontainer-agent-features-sandbox.md) liegt in `done/`; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice überschreitet 3 Liefer-Punkte.
- `in-progress` → `open` (blockiert): die Messung findet eine Lockerung des Profils, die der CR nicht trägt (dann zuerst CR erweitern), oder `ghcr.io/devcontainers/features/sshd:1` baut auf dem Basis-Image nicht.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **`sshd` öffnet im Container einen Passwort-Login** (Nutzer-Passwort per `sudo passwd`), im Sandbox-Profil mit Zielsetzung Schadensbegrenzung; entschieden ist: erlaubt mit Warnung. — **Ausgang:** weiter offen bis der CR die Teilfrage (Wortlaut, Sichtbarkeit unter `--quiet`) entschieden hat; danach trägt DoD 1 die gewählte Fassung.
- **Feature-Metadaten unterlaufen das Profil** (`mounts`, `capAdd`, `privileged`, `entrypoint`): u-boot erzeugt nur `devcontainer.json`; die Features bringen eigene Metadaten mit, die der Dev-Container-Lauf zusammenführt. — **Ausgang:** weiter offen: die Messung (DoD 2) entscheidet; bei Eintritt Rückführung an den CR (§4), bei Entfall wird er bei der Closure mit Verweis auf die Evidence gesetzt.
- **Der Port 2222 ist im Sandbox-Profil nur über das Forwarding der Dev-Container-Werkzeuge erreichbar**, `u-boot` veröffentlicht ihn nicht. — **Ausgang:** entfallen: `forwardPorts` bleibt unberührt, und `doctor` warnt nur bei fehlenden Service-Ports (aus `doctor.go` gelesen); die Handbuch-Zeile nennt das.
- **Plan-Aussage über den Code stimmt nicht** (Klasse mit drei Vorgängen an der Schwelle, siehe §8): die Fundstellen in §1 und §3 sind ein Lesestand. — **Ausgang:** weiter offen: der Implementer prüft jede Aussage vor der Änderung im Code; Abweichungen stehen in der Evidence, die Klasse zählt die Closure.

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `internal/hexagon/application` (Kürzel APP) und, für die Doku, `docs/user/`, `harness/` und `spec/` (DOK) — beide Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006) (Abschnitt Modus-Deklaration), beide erfüllen die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. `internal/adapter/driving/cli` (CLI) wird nur berührt, falls der CR die Warnung unter `--quiet` sichtbar hält und das nicht über das Feld `SecurityRelaxation` allein geht; dann gilt `AGENTS.md §3.9` (Kommentar-Selbstcheck).

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer für APP und DOK. Die einzige Beobachtung `BEO-CLI/kommentar-beschreibt-abwesendes` (3×, verkörpert in `AGENTS.md §3.9`) betrifft CLI; sie wird nur bei einer CLI-Berührung relevant. Die Klasse „Plan-Aussage über Code stimmt nicht“ liegt nach dem Auftrag an der Schwelle, steht aber nicht im Register; ihre Belege legt die Slice-Closure an (Register-Regel: eingetragen wird bei der Closure).

Alle berührten Sub-Areas GF.
