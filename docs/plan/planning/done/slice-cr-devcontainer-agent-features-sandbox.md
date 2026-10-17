# Slice CR: Sandbox-Regel für `sshd` und Claude Code und Zugangsdaten des Agenten (`LH-FA-DEV-006`, `LH-FA-DEV-009`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md) — erster Slice der Welle; ihm folgen [`slice-v1-cli-quiet-sicherheitswarnungen`](../done/slice-v1-cli-quiet-sicherheitswarnungen.md), [`slice-v1-devcontainer-feature-sshd`](../done/slice-v1-devcontainer-feature-sshd.md), [`slice-v1-devcontainer-feature-claude-code`](../in-progress/slice-v1-devcontainer-feature-claude-code.md) und [`slice-v2-spec-devcontainer-sshd-claude-code`](../open/slice-v2-spec-devcontainer-sshd-claude-code.md).

**Bezug:** [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil; Änderung), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten; Erweiterung auf Agenten-Zugangsdaten), [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Features, Zustimmung), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts des Agenten), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Verbosity; Änderung der Ausnahme 1). Verwandt: [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), und [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 1), §16 (Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt der CR nicht; die Spezifikation folgt im letzten Slice der Welle.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Zusagen aus den Entscheidungen des Projektinhabers vom 2026-10-17 zu `sshd` und Claude Code im Sandbox-Profil stehen im Lastenheft (Änderung
**0.10.0 → 0.11.0**): die Zugangsdaten des Agenten in [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) die Ausweisung des Passwort-Logins von `sshd` in [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und die Sichtbarkeit der Sicherheitswarnungen unter `--quiet` in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Ausnahme 1).
Die Code-Slices der Welle beginnen erst danach.

**Gefundener Stand im Lastenheft und im Code (Ausgangslage, aus dem Quelltext gelesen):**

- [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) zählt die Sicherheitslockerungen auf (Capabilities, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices; mit [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) und [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)
  `--privileged` und `NET_ADMIN`) und sagt über Features nur: [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) gelten unverändert, das Profil gibt keine externe Quelle implizit frei.
  Ein Login-Dienst im Container (`sshd`) ist weder Lockerung im Sinn der Aufzählung noch ausgeschlossen.
- [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) gilt für **Git**-Zugangsdaten. Im Code passt dazu: `templates/devcontainer/devcontainer.json.tmpl` reicht im Sandbox-Zweig nur `GIT_TOKEN` als
  `remoteEnv` durch (ohne Sandbox gibt es kein `remoteEnv`), `credentialPatterns()` in `doctor_sandbox.go` erkennt Git-Token-Muster, `hasTokenSource` prüft nur
  `GIT_TOKEN` und einen Secret-Mount. Ein Anthropic-Schlüssel ist nirgends erfasst.
- Ob eine Warnung zur Befehlsausgabe unter `--quiet` sichtbar bleibt, entscheidet das Feld `SecurityRelaxation` des `WarningEntry`
  (`sandboxRelaxationWarnings` in `devcontainer_sandbox.go`; Ausnahme 1 von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) nennt heute Capabilities, `--privileged`, Seccomp-/AppArmor-Anpassungen und
  zusätzliche Devices in der Befehlsausgabe). Die Drittanbieter-Warnung (`thirdPartyFeatureWarnings`) setzt das Feld nicht.
- `u-boot doctor --quiet` zeigt heute nur Fehler-Einträge und Einträge mit dem Kennzeichen `DegradationFallback` (`writeDoctorReport` in `doctor.go`, `domain.Diagnostic`);
  die Warn-Einträge der Prüfungen `devcontainer.sandbox.credentials` und `devcontainer.sandbox.egress` (`doctor_sandbox.go`) entfallen unter `--quiet`.
- Katalog-Schlüssel allein sind keine Zusage: `python`, `sdkman` und `homebrew` kamen ohne Lastenheft-Änderung
  ([`slice-v1-devcontainer-features-homebrew-sdkman`](../done/slice-v1-devcontainer-features-homebrew-sdkman.md),
  [`slice-v1-devcontainer-feature-python`](../done/slice-v1-devcontainer-feature-python.md)); [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) führt Beispiele auf, keine abschließende Liste.

**Entscheidungen des Projektinhabers am 2026-10-17** (gewählte Option, Begründung, bestätigt):

| # | Frage | Gewählt | Begründung | Folge für den Vertrag |
|---|---|---|---|---|
| 1 | `sshd` im Sandbox-Profil | erlaubt, die Befehlsausgabe weist den Passwort-Login aus (Warnung); die Warnung ist eine Sicherheitslockerung im Sinn der Ausnahme 1, trägt `SecurityRelaxation` und bleibt unter `--quiet` | Das Profil ist Schadensbegrenzung, kein Verbot von Diensten; der Nutzer entscheidet bewusst und wird informiert, auch in einem stillen Lauf. | neue Ausweisungs-Zusage in [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil); Ausnahme 1 in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) nennt den Login-Dienst |
| 2 | Claude Code im Sandbox-Profil und seine Egress-Hosts | erlaubt, Hosts gemessen (erwartet: npm-Registry beim Build, `api.anthropic.com` zur Laufzeit; Anmeldung, Telemetrie, Update und Anmeldeweg gemessen, nicht geraten) | Der Agent ist der Anlass des Sandbox-Profils; seine Hosts gehören in die Default-Allowlist je Feature. | keine Lastenheft-Änderung (die Hosts stehen in [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), nicht im Lastenheft) |
| 3 | Einstufung des Features von Anthropic | Drittanbieter-Feature nach [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) (Warnung bei Aktivierung, Herausgeber `anthropics`) | [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) trägt die Regel bereits; eine zweite Regel wäre eine Parallelwelt. | keine; keine Folge-ADR |
| 4 | Zugangsdaten des Agenten | [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) wird auf Zugangsdaten des Agenten erweitert: nie im Image, Workspace-Volume, in `u-boot.yaml` oder einer erzeugten Datei; Übergabe zur Laufzeit per Umgebungsvariable; `doctor` prüft die Quelle wie bei Git-Tokens | Das Sandbox-Profil führt mit dem Agenten einen zweiten Schlüsseltyp ein; eine Zusage gehört in das Lastenheft, nicht in die Spezifikation. | neue Zusage in [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) |
| 5 | Sicherheitswarnungen unter `--quiet` | gewählt: sie bleiben sichtbar, und zwar zusätzlich die `u-boot doctor`-Warn-Einträge, die eine Sicherheitseigenschaft des Sandbox-Profils betreffen (Quelle der Zugangsdaten, Egress-Restriktion, Sicherheitslockerungen des Docker-Daemons); nicht gewählt: die Drittanbieter-Code-Warnung beim Image-Build ([`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md)) | Ein Skript soll die Sicherheitslage des Profils nicht stumm schalten können; die Drittanbieter-Warnung betrifft die Herkunft eines Features, nicht eine Eigenschaft des laufenden Profils. | Ausnahme 1 in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) wird erweitert, die Zahl der Ausnahmen bleibt vier; die Drittanbieter-Warnung entfällt weiter unter `--quiet`; Code folgt in [`slice-v1-cli-quiet-sicherheitswarnungen`](../done/slice-v1-cli-quiet-sicherheitswarnungen.md) |

**Entscheidung zu `sshd` (Projektinhaber, 2026-10-17):** Ein Passwort-Login im Container ist eine Sicherheitslockerung im Sinn der Ausnahme 1 von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (im Code das Kennzeichen `SecurityRelaxation` am `WarningEntry`): Ausnahme 1 nennt den Login-Dienst, die Warnung trägt `SecurityRelaxation`, und sie bleibt unter `--quiet` bei `init`, `generate` und `up`.

Inhalt der Warnung: Der Login-Dienst `sshd` ist aktiv (Port 2222, nur Passwort-Login nach `sudo passwd`); das Profil begrenzt Schaden und ist keine Isolationsgrenze.
Ein sshd-Port wird nicht automatisch weitergeleitet: `forwardPorts` stammt nur aus den Services ([`LH-FA-DEV-005`](../../../../spec/lastenheft.md#lh-fa-dev-005--ports)).

**Kern des CR — Vorschlag für den erweiterten Text von [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)** (Stil und Aufbau des heutigen Textes; im Lastenheft stehen die Verweise als Anker `#…` der Datei; die Überschrift bleibt, damit die Anker im Repo stabil sind; das Lastenheft nennt keine Slice- oder Wellen-Kennung):

> *Priorität: V1.* Im Sandbox-Profil ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)) dürfen Git-Zugangsdaten und Zugangsdaten eines Agenten (etwa der Zugangsschlüssel eines Coding-Agenten, den ein aktiviertes Devcontainer-Feature nutzt) weder im Image, im Workspace-Volume noch in `u-boot.yaml` oder einer anderen erzeugten Datei stehen.
>
> - Übergabe zur Laufzeit per Umgebungsvariable; Git-Zugangsdaten alternativ als schreibgeschützter Secret-Mount. Die erzeugte Konfiguration trägt nur die Referenz auf die Umgebungsvariable des Hosts, nie den Wert.
> - Die Dokumentation beschreibt kurzlebige, auf das Repository begrenzte Tokens (privater Schlüssel nie im Container) und die Forderung nach Branch-Protection auf dem Remote; für den Agenten nennt sie die Umgebungsvariable, aus der der Zugangsschlüssel stammt.
> - `u-boot doctor` prüft die Quelle der Zugangsdaten (Git-Token-Quelle, bei aktivem Agenten-Feature die Umgebungsvariable des Agenten): fehlt sie oder liegt sie im Klartext in einer Projektdatei, `warn` (kein Abbruch).

**Kern des CR — Wortlaut der Ausnahme 1 von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)** (Entscheidung des Projektinhabers; die Zahl der Ausnahmen bleibt vier; im Lastenheft stehen die Verweise als Anker `#…` der Datei, hier mit Pfad der Datei):

> 1. die Ausweisung von Sicherheitslockerungen (zusätzliche Capabilities, `--privileged`, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices) und eines im Container aktivierten Login-Dienstes (Passwort-Login) in der Befehlsausgabe ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) sowie die Warn-Einträge von `u-boot doctor`, die eine Sicherheitseigenschaft des Sandbox-Profils betreffen (Quelle der Zugangsdaten nach [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), Egress-Restriktion nach [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), Sicherheitslockerungen des Docker-Daemons nach [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)); dazu gehört die Meldung, dass die Egress-Restriktion ohne Sandbox-Profil nicht greift; die Meldung, dass die Container-Runtime ohne Sandbox-Profil keine Wirkung hat, ist eine Konfigurationsmeldung und gehört nicht dazu;

Ausnahme 2 bis 4 und der übrige Text von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) bleiben wörtlich.

Vorschlag für den Zusatz in [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), Liste „muss“: *„- einen im Container aktivierten Login-Dienst (Passwort-Login) in der Befehlsausgabe ausweisen.“*
Vorschlag für die Historie (§16): Zeile 0.11.0, 2026-10-17, Quelle „Vereinbarung mit dem Projektinhaber“: [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer): Zusage auf Zugangsdaten eines Agenten erweitert (nie in Image, Workspace-Volume, `u-boot.yaml` oder erzeugter Datei; Übergabe per Umgebungsvariable; `doctor` prüft die Quelle); [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): Ausweisung eines im Container aktivierten Login-Dienstes (Passwort-Login) im Sandbox-Profil; [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging): Ausnahme 1 unter `--quiet` nennt auch diesen Login-Dienst und die Warn-Einträge von `u-boot doctor` zu Sicherheitseigenschaften des Sandbox-Profils, darunter die Meldung, dass die Egress-Restriktion ohne Sandbox-Profil nicht greift.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code, Katalog und Egress-Hosts** — die Code-Slices [`slice-v1-cli-quiet-sicherheitswarnungen`](../done/slice-v1-cli-quiet-sicherheitswarnungen.md) (`doctor --quiet`), [`slice-v1-devcontainer-feature-sshd`](../done/slice-v1-devcontainer-feature-sshd.md) und
  [`slice-v1-devcontainer-feature-claude-code`](../in-progress/slice-v1-devcontainer-feature-claude-code.md) setzen um, was dieser CR festlegt; der CR ändert nur den Vertrag.
- **Die Drittanbieter-Code-Warnung beim Image-Build unter `--quiet`** — der Projektinhaber hat sie nicht gewählt (Entscheidung 5); sie bleibt außerhalb der Ausnahme 1 und entfällt unter `--quiet`.
- **`spezifikation.md` und `architecture.md`** — die Spezifikation beschreibt das Verhalten erst, wenn Lastenheft und Code es tragen
  ([`slice-v2-spec-devcontainer-sshd-claude-code`](../open/slice-v2-spec-devcontainer-sshd-claude-code.md)); die Sicht nennt keine Features.
- **Eine Aufnahme von `sshd` oder Claude Code in die Beispielliste von [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)** — Beispiele sind keine Zusage; der Katalog führt sie ohne Lastenheft-Änderung (Muster `python`).
- **Schlüssel-Login, Passwort-Setzen und Agent-Konfiguration durch `u-boot`** — [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) nennt den Agenten als Nicht-Ziel; ein Passwort setzt der Nutzer im Container.
- **Eine ADR** — Entscheidung 3 folgt [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md); keine Folge-ADR.

### Prüfung gegen die älteren Zusagen (Steering-Loop-Lehre der vorigen CRs: zuerst per Suche im Lastenheft)

**Suche (Stand 2026-10-17, Lastenheft 0.10.0):** `rg -n -i 'Zugangsdaten|Token|Secret|Credential|Passwort|Geheimnis' spec/lastenheft.md` <!-- d-check:ignore (Suchbefehl, kein Link-Pfad) --> findet im Normtext genau diese Stellen:
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Clone-URL mit Zugangsdaten ist Fehler; keine Host-Dateien mit Geheimnissen),
[`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Git-Zugangsdaten, Token-Quelle in `doctor`),
[`LH-NFA-SEC-001`](../../../../spec/lastenheft.md#lh-nfa-sec-001--keine-secrets-einchecken) (keine echten Secrets in erzeugte Dateien),
die Admin-Beispielwerte der Services (`.env.example` mit Platzhaltern; [`LH-NFA-SEC-002`](../../../../spec/lastenheft.md#lh-nfa-sec-002--beispielwerte-markieren)) und die Nicht-Ziele (Secret-Management für Produktion).
Die Suche nach `Login`, `SSH` und `ssh` findet nur `~/.ssh` in [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil).

| Stelle | Berührung | Ergebnis |
|---|---|---|
| [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (externe Quellen, Zustimmung) | Anthropic-Feature im Katalog | verträglich: es ist ein Drittanbieter-Feature nach [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md); Zustimmung und `featureSources.allow` gelten wie bei `homebrew` |
| [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Lockerungen ausweisen, keine Host-Geheimnisse, Features unverändert) | `sshd`, Claude Code im Profil | angepasst: neue Ausweisung des Login-Dienstes; „keine Host-Dateien mit Geheimnissen“ bleibt verträglich, weil der Schlüssel als Umgebungsvariable und nicht als Datei übergeben wird; Nicht-Ziel „kein Agent-Start, kein Berechtigungsmodus“ bleibt |
| [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) und [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) | Degradations-Fallbacks, Hosts des Agenten | verträglich: die Default-Allowlist ist „dokumentiert“ ohne Hostliste im Lastenheft; die Hosts je Feature gehören in die Spezifikation |
| [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | Zugangsdaten des Agenten; `doctor`-Warnung zur Quelle | angepasst: Erweiterung (siehe Kern des CR); die Git-Zusagen bleiben wörtlich; die `doctor`-Warnung zur Quelle der Zugangsdaten bleibt unter `--quiet` sichtbar (Ausnahme 1) |
| [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (vier Ausnahmen unter `--quiet`; `doctor` zeigt unter `--quiet` nur Fehler-Einträge) | Sichtbarkeit der sshd-Warnung und der `doctor`-Sicherheitswarnungen | angepasst (entschieden): Ausnahme 1 nennt den Login-Dienst und die Warn-Einträge von `u-boot doctor`, die eine Sicherheitseigenschaft des Sandbox-Profils betreffen; die Zahl der Ausnahmen bleibt vier, Ausnahme 2 bis 4 bleiben wörtlich; die Drittanbieter-Warnung bleibt außerhalb |
| [`LH-NFA-SEC-001`](../../../../spec/lastenheft.md#lh-nfa-sec-001--keine-secrets-einchecken) | Referenz auf die Umgebungsvariable in `devcontainer.json` | verträglich: die Datei trägt die Referenz, nie den Wert |
| [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | Anthropic-Installationsskript beim Image-Build | verträglich: das Feature ist externer Code aus einer nicht freigegebenen Quelle und braucht die Zustimmung nach [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features); die Drittanbieter-Warnung weist es aus |

Die Zeilen sind gegen den Normtext gelesen; der Reviewer des CR bestätigt sie.

## 2. Definition of Done

- [x] Die Entscheidungen vom 2026-10-17 stehen in §1 (Entscheidungsdokumentation); die Prüfungstabelle ist bestätigt.
- [x] Lastenheft **0.11.0**: die Zusagen stehen in [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Agenten-Zugangsdaten) und [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Ausweisung des Login-Dienstes), außerdem in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Ausnahme 1 mit dem Wortlaut aus §1); Versionsfeld und Historie (§16, Quelle „Vereinbarung mit dem Projektinhaber“; die Zeile nennt: „[`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging): Ausnahme 1 unter `--quiet` nennt auch diesen Login-Dienst und die Warn-Einträge von `u-boot doctor` zu Sicherheitseigenschaften des Sandbox-Profils, darunter die Meldung, dass die Egress-Restriktion ohne Sandbox-Profil nicht greift.“) tragen den Eintrag; die Suche in §1 belegt, dass keine ältere Zusage widerspricht.
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
| `spec/lastenheft.md` (§3 [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 1, Versionsfeld, §16 Historie) | update | die entschiedenen Zusagen; Historie nach dem Muster der CRs zu [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |
| dieser Slice (§1 Entscheidungsdokumentation, Prüfungstabelle, §7) | update | Entscheidungsdokumentation |

Kein Code, keine Spezifikation, keine Tests; `make docs-check` und `make gates` laufen für die Dokumente.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; die Entscheidungen 1 bis 5 des Projektinhabers (2026-10-17) liegen vor; kein anderer Slice in `in-progress/`
(WIP-Limit 1; solange [`slice-v2-homebrew-smoke-haertung`](../done/slice-v2-homebrew-smoke-haertung.md) dort liegt, beginnt kein Slice der Welle).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Text verlangt eine neue Anforderung statt der Änderung von zwei bestehenden — dann ein CR je Anforderung.
- `in-progress` → `open` (blockiert): eine Entscheidung wechselt, oder die Messung in einem Code-Slice findet eine Lockerung, die der CR noch nicht trägt (der CR wird dann erweitert, bevor der Code-Slice weiterläuft).

## 5. Closure-Trigger

DoD vollständig (Lastenheft 0.11.0 mit Historie), Review-Report liegt vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Der Wortlaut der Ausnahme 1 bildet eine Stelle nicht ab, die der Code setzt** (etwa ein `doctor`-Warn-Eintrag, den die Einstufung „Sicherheitseigenschaft des Sandbox-Profils“ anders trifft als die Code-Slices). — **Ausgang:** entfallen: Wortlaut und Einstufung sind gegen den Code gelesen; die Menge „Sicherheitslockerungen des Docker-Daemons“ als `doctor`-Warn-Eintrag ist leer, weil der Eintrag zu `nestedRuntime: docker` die Stufe OK hat (Verification-Evidence [`2026-10-17-verification-slice-cr-devcontainer-agent-features-sandbox`](../../../reviews/2026-10-17-verification-slice-cr-devcontainer-agent-features-sandbox.md), Abschnitt (d) und (e)); die Einstufung je Eintrag trägt [`slice-v1-cli-quiet-sicherheitswarnungen`](../done/slice-v1-cli-quiet-sicherheitswarnungen.md) §1.
- **Der Anmeldeweg des Agenten ist nicht die Umgebungsvariable** (Anmeldung im Container statt Schlüssel vom Host): die `doctor`-Warnung „Quelle fehlt“ wäre dann ein Fehlalarm. — **Ausgang:** eingetreten: [`slice-v1-devcontainer-feature-claude-code`](../in-progress/slice-v1-devcontainer-feature-claude-code.md) trägt es (Messung des Anmeldewegs in der DoD, Risiko dort); weicht der Weg ab, wird der Text von [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) im CR vor dem Code-Slice nachgezogen (Rückführung `in-progress` → `open`).
- **Die Features bringen Metadaten mit** (`mounts`, `capAdd`, `privileged`, `entrypoint` in ihrer `devcontainer-feature.json`) und unterlaufen damit „kein Host-Bind-Mount, kein `--privileged`“ des Profils; weder Plan noch Lastenheft sagen es voraus. — **Ausgang:** eingetreten: [`slice-v1-devcontainer-feature-sshd`](../done/slice-v1-devcontainer-feature-sshd.md) und [`slice-v1-devcontainer-feature-claude-code`](../in-progress/slice-v1-devcontainer-feature-claude-code.md) tragen es (je ein Risiko mit Messung); findet die Messung eine Lockerung, wird dieser CR vor dem Code-Slice um die Ausweisung erweitert (Rückführung `in-progress` → `open`).
- **Widerspruch zu einer älteren Zusage** bleibt unentdeckt, weil die Suche zu eng war. — **Ausgang:** entfallen: die Prüfungstabelle in §1 trägt die Suche (Stand 2026-10-17); der Review bestätigt sie.

## 7. Closure-Notiz

**Belege:** zwei Review-Reports, beide ohne HIGH: [`2026-10-17-slice-cr-devcontainer-agent-features-sandbox`](../../../reviews/2026-10-17-slice-cr-devcontainer-agent-features-sandbox.md)
(0 HIGH, 1 MEDIUM, 4 INFO; lief gegen `5de388c`; F-1 CHANGELOG im Indikativ behoben in `4df2b54`, F-2 bis F-4 ohne Vertragsfolge, F-5 nachgezogen) und der Nachlauf-Review
[`2026-10-17-slice-cr-devcontainer-agent-features-sandbox-2`](../../../reviews/2026-10-17-slice-cr-devcontainer-agent-features-sandbox-2.md)
(0 HIGH, 1 MEDIUM, 4 INFO; deckt die Erweiterung der Ausnahme 1 ab; F-1 „Ausnahme 1 pinnt die Unterscheidung Egress/Runtime nicht“ behoben in `abd906f`, F-2 bis F-5 INFO ohne Vertragsfolge);
Verification Evidence [`2026-10-17-verification-slice-cr-devcontainer-agent-features-sandbox`](../../../reviews/2026-10-17-verification-slice-cr-devcontainer-agent-features-sandbox.md)
(erfüllt mit Auflagen A-1 bis A-6; Anhang „Nachverifikation abd906f“: bestanden, 0 offene Befunde); `make gates` grün, `make doc-trace`: 83 Anforderungen, 0 Waisen.

- **Was hat funktioniert:** Die Vorab-Suche gegen ältere Zusagen im Lastenheft fand die Stellen, die die Erweiterung berührt; Review, Nachlauf-Review und Verifier wiederholten sie ohne weiteren Treffer. Beide Reviews fanden je ein MEDIUM, beide zur Wortwahl, keines zum Vertragsinhalt. Der Verifier belegte die leere Menge „Sicherheitslockerungen des Docker-Daemons“ (V-5) statt sie zu behaupten. Die Entscheidungen 1 bis 5 des Projektinhabers sind in den Vertrag eingegangen, und die Unterscheidung Egress/Runtime steht im Normtext von Ausnahme 1.
- **Was ging anders als geplant:** Ausnahme 1 wurde mitten im Slice um die `doctor`-Sicherheits-Einträge erweitert (Entscheidung 5); der erste Review lief davor und deckte die Erweiterung nicht ab, daher waren ein Nachlauf-Review und eine Nachverifikation nötig. Die Teilfrage zu `sshd` und die Entscheidung zu Egress/Runtime kamen vom Projektinhaber, nicht aus dem Plan. Die Historie-Zeile des Slice wich vom geltenden Lastenheft ab (V-1), ein Link-Text trug einen Tippfehler (V-2), und die Welle nannte einen veralteten Stand (V-7).
- **Steering-Loop-Eintrag:** Ein Review, der vor einer Erweiterung des Vertrags lief, deckt die Erweiterung nicht ab. Verfahren: Ändert sich der Vertrag nach dem Review, wird der Diff der Änderung erneut reviewt (Nachlauf-Review), bevor die Closure beginnt. Gezählt, nicht verkörpert; erstes Auftreten.
- **Beobachtungs-Register (`../observations/`):** [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) angelegt (Entscheidung des Projektinhabers 2026-10-17), vier Belege (`slice-v1-cli-quiet-nur-fehler`, `slice-v1-cli-quiet-vorschau`, `slice-v1-cli-verbosity-letzte-option-gewinnt`, `slice-v2-homebrew-smoke-haertung`); verkörpert in `AGENTS.md §3.10`, Herkunfts-Anker `seit welle-devcontainer-sshd-claude-code`. <!-- d-check:ignore (Pfad-Zitat im Code-Span) -->
- **Folge-Slices:** [`slice-v1-cli-quiet-sicherheitswarnungen`](../done/slice-v1-cli-quiet-sicherheitswarnungen.md), [`slice-v1-devcontainer-feature-sshd`](../done/slice-v1-devcontainer-feature-sshd.md), [`slice-v1-devcontainer-feature-claude-code`](../in-progress/slice-v1-devcontainer-feature-claude-code.md), [`slice-v2-spec-devcontainer-sshd-claude-code`](../open/slice-v2-spec-devcontainer-sshd-claude-code.md) — Dateien in `open/`.
- **Risiken aus §6:** alle vier tragen ihren Ausgang in kanonischer Form (entfallen: Wortlaut gegen Code, Suche zu eng; eingetreten: Anmeldeweg und Feature-Metadaten, getragen von den Code-Slices).
- **Release-Hinweis:** Kein Release-Cut-Hindernis aus diesem CR: der CHANGELOG-Eintrag nennt nur den Vertrag. Der Satz „follow in a later release“ ist anzupassen, sobald die Code-Slices landen: er steht als Punkt in der DoD von [`slice-v1-cli-quiet-sicherheitswarnungen`](../done/slice-v1-cli-quiet-sicherheitswarnungen.md), [`slice-v1-devcontainer-feature-sshd`](../done/slice-v1-devcontainer-feature-sshd.md) und [`slice-v1-devcontainer-feature-claude-code`](../in-progress/slice-v1-devcontainer-feature-claude-code.md); der letzte Code-Slice entfernt ihn.
- **Grenzen:** (1) Der Slice ändert nur den Vertrag; Zusagen zu `--quiet`, Login-Dienst und Agenten-Zugangsdaten sind bis zu den Code-Slices ohne Code-Beleg. (2) Handbuch und [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) nennen die neue Ausnahme 1 erst nach Code- und Spec-Slice; bis dahin gilt Lastenheft vor Spezifikation. (3) „Login-Dienst (Passwort-Login)“: die Klammer erklärt den Dienst; ein `sshd` nur mit Schlüssel-Login wäre nicht eindeutig geregelt (Review F-3); der sshd-Slice legt den Fall fest. (4) „Sicherheitslockerungen des Docker-Daemons“ ist als `doctor`-Warn-Eintrag heute eine leere Menge (V-5). (5) Die Überschrift von [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) „Git-Zugangsdaten …“ passt nicht mehr ganz zum Text; sie bleibt bewusst unverändert, weil die Anker im Repo stabil bleiben (Review F-4). (6) Die Drittanbieter-Code-Warnung beim Image-Build entfällt weiter unter `--quiet`; der Projektinhaber hat sie nicht gewählt (Entscheidung 5).
- **Drei Paarungen:** (nach dem `git mv` nach `done/` zu prüfen) <!-- d-check:ignore (Pfad-Zitat im Code-Span) -->

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `spec/` (Kürzel DOK) — Modus Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006) (Abschnitt Modus-Deklaration), erfüllt die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. Code-Sub-Areas berührt der CR nicht.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer für diese Sub-Area (die einzige Beobachtung `BEO-CLI/kommentar-beschreibt-abwesendes` betrifft `internal/adapter/driving/cli` und ist verkörpert in `AGENTS.md §3.9`).

Alle berührten Sub-Areas GF.
