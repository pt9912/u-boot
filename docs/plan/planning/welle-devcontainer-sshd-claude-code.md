# Welle welle-devcontainer-sshd-claude-code: Devcontainer-Features `sshd` und Claude Code im eingebauten Katalog

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/`. Der Zustand ist die Verzeichnis-Position — kein Status-Feld.
**Geplante Wellen bekommen noch keine Datei:** Sie stehen in der Roadmap unter
*Nächste Wellen* und nirgends sonst — zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug. Die Zuordnung zu einem Release-Cut (nach v0.9.0) liegt beim Projektinhaber und ist offen; die Welle ist keine Release-Version.

**Verantwortlich:** pt9912. **Datum:** 2026-10-17.

---

## 1. Welle-Ziel

Der eingebaute Devcontainer-Feature-Katalog bietet `sshd` und `claude-code` als wählbare Schlüssel (wie `python`, `sdkman` und `homebrew`), jeweils mit
Egress-Hosts für das Sandbox-Profil. Im Sandbox-Profil
([`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)) sind beide Features erlaubt (Entscheidung des Projektinhabers am 2026-10-17): `sshd` mit einer Warnung zum Passwort-Login, Claude Code mit gemessenen Egress-Hosts. Sicherheitswarnungen bleiben unter `--quiet` sichtbar (Entscheidung des Projektinhabers am 2026-10-17): die Warnung zum Login-Dienst und zusätzlich die `u-boot doctor`-Warn-Einträge, die eine Sicherheitseigenschaft des Sandbox-Profils betreffen (Ausnahme 1 von [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)). Ein Zugangsschlüssel des Agenten steht in keiner erzeugten Datei, sondern wird zur Laufzeit übergeben
([`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) gilt für Git-Zugangsdaten und Zugangsdaten eines Agenten; Lastenheft 0.11.0). Beide Features
zusammen in einem Devcontainer erzeugen ein stabiles `devcontainer.json` und eine gültige Egress-Allowlist.

Anlass (Projektinhaber, 2026-10-17): Entscheidung für Devcontainer-Features statt Add-on-Services, als eigene Welle.

## 2. Trigger (Welle startet)

- Die Entscheidungen des Projektinhabers vom 2026-10-17 liegen vor (eingetreten): `sshd` und Claude Code kommen als eingebaute Katalog-Features (Devcontainer-Features,
  nicht Add-on-Services), als eigene Welle; `sshd` im Sandbox-Profil erlaubt mit Warnung, Claude Code im Sandbox-Profil erlaubt mit gemessenen Egress-Hosts, das Anthropic-Feature
  ein Drittanbieter-Feature nach [`ADR-0020`](../adr/0020-drittanbieter-features-im-katalog.md), [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) wird auf Agenten-Zugangsdaten erweitert, und die Sicherheitswarnungen (Login-Dienst, `doctor`-Einträge) bleiben unter `--quiet` sichtbar, die Drittanbieter-Warnung nicht. Die Einzelheiten stehen im CR.
- Der Slice [`slice-v2-homebrew-smoke-haertung`](done/slice-v2-homebrew-smoke-haertung.md) liegt in `done/` (WIP-Limit 1: solange er in
  `in-progress/` liegt, beginnt kein Slice dieser Welle; beobachtbar an der Verzeichnis-Position des Slice). Erst danach wird der CR nach `in-progress/` bewegt.

## 3. Closure-Trigger (Welle schließt)

Das *Mehr* gegenüber den fünf Slice-DoDs: Die beiden Features tragen **zusammen**, und die Stellen, an denen der Vertrag steht, sagen dasselbe.

- Alle fünf Slices der Welle liegen in `done/`; der CR [`slice-cr-devcontainer-agent-features-sandbox`](done/slice-cr-devcontainer-agent-features-sandbox.md) liegt mit Lastenheft 0.11.0 dort.
- **Beide Features zusammen:** Ein Golden Case aktiviert `sshd` und `claude-code` gemeinsam (einmal ohne, einmal mit Sandbox-Profil und Egress): `devcontainer.json`
  trägt jeden Render-Schlüssel genau einmal in stabiler Reihenfolge, der zweite Lauf ist ein NoOp, die Egress-Allowlist ist die sortierte, duplikatfreie Vereinigung
  der Hosts beider Features, und die Befehlsausgabe trägt genau die Hinweise, die der CR vorschreibt.
- **Sicherheitswarnungen unter `--quiet`:** `--quiet` zeigt bei `init`, `generate` und `up` die Warnung zum Login-Dienst und bei `doctor` die Warn-Einträge zur Quelle der Zugangsdaten (Git und Agent) und zur Egress-Restriktion mit `docker`; die Egress-Meldung „has no effect“ bleibt sichtbar; die Drittanbieter-Warnung, die `nestedRuntime`-Meldung „has no effect“ und die Feature-Allowlist- und Drift-Warnungen entfallen. Lastenheft (Ausnahme 1), Code, Handbuch und CHANGELOG sagen dasselbe; die Spezifikation trägt es in [`SPEC-043`](../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet).
- **Lauffähigkeit belegt, nicht angenommen:** Eine Messung mit der Dev-Container-CLI (Verfahren der Messreihe in
  [`ADR-0020`](../adr/0020-drittanbieter-features-im-katalog.md)) baut ein Image mit beiden Features und startet es; das Ergebnis steht in der Evidence der beiden
  Feature-Slices. Läuft die Messung in der Umgebung nicht, nennt die Closure den Grund (kein stilles Rot).
- **Eine Aussage über die Hosts des Agenten:** Handbuch, README, CHANGELOG und Spezifikation sagen dasselbe; der Satz im Handbuch, `u-boot` kenne die Hosts des
  Agenten nicht (Abschnitt Egress), ist an den neuen Stand angepasst.
- Die Spezifikation trägt die Katalogzeilen (`sshd`, `claude-code`), die Egress-Zeilen, die Nachführung von `SPEC-043` (Ausnahme 1) und den neuen Eintrag ab `SPEC-044`; sie entstand **nach** dem Abschluss der <!-- d-check:ignore (SPEC-044 ist noch nicht vergeben) -->
  vier anderen Slices und ist gegen den Code gelesen (Commit-Reihenfolge belegt es).
- `make doc-trace` meldet keine Waise; `make docs-check` meldet 0 Befunde; `make gates` ist grün.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-cr-devcontainer-agent-features-sandbox`](done/slice-cr-devcontainer-agent-features-sandbox.md) | CR: Sandbox-Regel für `sshd` und Claude Code, Zugangsdaten des Agenten (zuerst; Lastenheft 0.11.0) | [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |
| [`slice-v1-cli-quiet-sicherheitswarnungen`](in-progress/slice-v1-cli-quiet-sicherheitswarnungen.md) | Code: `doctor --quiet` zeigt die Sicherheits-Warnungen des Sandbox-Profils (Kennzeichen am Eintrag; nach dem CR, vor `claude-code` zwingend, vor `sshd` aus Wahl) | [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) |
| [`slice-v1-devcontainer-feature-sshd`](open/slice-v1-devcontainer-feature-sshd.md) | Code: Katalog-Schlüssel `sshd` (offizielles Feature), Egress, Sandbox-Regel, Tests, Doku | [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |
| [`slice-v1-devcontainer-feature-claude-code`](open/slice-v1-devcontainer-feature-claude-code.md) | Code: Katalog-Schlüssel `claude-code` (Feature von Anthropic), Egress, Zugangsdaten zur Laufzeit, Tests, Doku | [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) |
| [`slice-v2-spec-devcontainer-sshd-claude-code`](open/slice-v2-spec-devcontainer-sshd-claude-code.md) | Spezifikation: Katalog, Egress-Hosts, Sandbox-Verhalten beider Features und Wirkung von `--quiet` ([`SPEC-043`](../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet); letzter Slice) | [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) |

## 5. Abhängigkeiten

- Blockiert: keine Welle.
- Wird blockiert von: [`slice-v2-homebrew-smoke-haertung`](done/slice-v2-homebrew-smoke-haertung.md) (WIP-Limit 1, siehe §2). Die fünf **Entscheidungen des
  Projektinhabers** vom 2026-10-17 liegen vor (Entscheidungsdokumentation im CR).
- Release-Zuordnung: Die Zuordnung der Welle zu einem Release-Cut liegt beim Projektinhaber und ist offen; kein Slice der Welle hängt an ihr.
- Reihenfolge in der Welle: der CR geht den Code-Slices voraus, weil sie die Sandbox-Regel erst nach dem Lastenheft umsetzen; [`slice-v1-cli-quiet-sicherheitswarnungen`](in-progress/slice-v1-cli-quiet-sicherheitswarnungen.md) geht `claude-code` voraus, weil dessen `doctor`-Warnung zur Agenten-Variable das dort eingeführte Kennzeichen trägt (zwingend), und `sshd` aus Wahl (die sshd-Warnung nutzt das bestehende `SecurityRelaxation` der Befehlsausgabe, keine Code-Abhängigkeit); `sshd` geht `claude-code` voraus
  (beide berühren Katalog, Egress und Tests in denselben Dateien, WIP-Limit 1); der Spezifikations-Slice kommt zuletzt, weil die Spezifikation das Verhalten erst
  beschreibt, wenn Lastenheft und Code es tragen.

## 6. Out-of-Scope für diese Welle

- Add-on-Services (Compose) für `sshd` oder Claude Code: entschieden gegen (Anlass 2026-10-17).
- Feature-Optionen (`gatewayPorts` und `version` von `sshd`, Optionen von Claude Code): der Katalog kennt keine Optionen; ein Optionsmodell ist ein
  Re-Evaluierungs-Trigger von [`ADR-0020`](../adr/0020-drittanbieter-features-im-katalog.md) und ein eigener Vorgang.
- Schlüssel-Login und Schlüsselverwaltung für `sshd` (`authorized_keys`, Host-Schlüssel): das Feature kennt nur Passwort-Login; `u-boot` verwaltet keine Schlüssel
  und setzt kein Passwort (neue Zusage, nicht Teil dieser Welle).
- Automatisches Eintragen von Port 2222 in `forwardPorts`: `forwardPorts` stammt aus den Services
  ([`LH-FA-DEV-005`](../../../spec/lastenheft.md#lh-fa-dev-005--ports)); ändert der Projektinhaber das, ist es eine Plan-Änderung.
- Konfiguration oder Start des Agenten (Berechtigungsmodus, Einstellungen, Hooks): [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)
  nennt es als Nicht-Ziel.
- Weitere Agenten- oder Remote-Zugangs-Features (Codex, Gemini, andere SSH-Varianten): jede Quelle braucht eine eigene Entscheidung.
- Der Release-Cut und die Zuordnung der Welle zu ihm: Sache des Projektinhabers, eigener Slice nach dieser Welle, wenn er ihn wünscht.
- Die Docker-Desktop-Messung: andere Gegenstände, ohne Bezug zu den Features.

## 7. Closure-Notiz

Wird bei der Closure geschrieben. Form wie in [`welle-cli-verbosity`](done/welle-cli-verbosity.md) §7: Detailquellen, Geliefert, Closure-Trigger (§3), Was hat
funktioniert, Was ging anders als geplant, Trigger-Audit, Steering-Loop-Einträge, Beobachtungs-Register, Folge-Slices, Verifikation. Die Pfade sind so zu
schreiben, wie sie vom Ruheort `done/` auflösen.
