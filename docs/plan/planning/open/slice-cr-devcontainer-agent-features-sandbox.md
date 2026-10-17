# Slice CR: Sandbox-Regel für `sshd` und Claude Code und Zugangsdaten des Agenten (`LH-FA-DEV-006`, `LH-FA-DEV-009`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md) — erster Slice der Welle; ihm folgen [`slice-v1-devcontainer-feature-sshd`](slice-v1-devcontainer-feature-sshd.md), [`slice-v1-devcontainer-feature-claude-code`](slice-v1-devcontainer-feature-claude-code.md) und [`slice-v2-spec-devcontainer-sshd-claude-code`](slice-v2-spec-devcontainer-sshd-claude-code.md).

**Bezug:** [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil; Änderung), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten; Erweiterung auf Agenten-Zugangsdaten), [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Features, Zustimmung), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts des Agenten), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (nur, falls die Teilfrage zu Entscheidung 1 „ja“ lautet). Verwandt: [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), und nur bei Teilfrage „ja“ [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)), §16 (Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt der CR nicht; die Spezifikation folgt im letzten Slice der Welle.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Zusagen aus den Entscheidungen des Projektinhabers vom 2026-10-17 zu `sshd` und Claude Code im Sandbox-Profil stehen im Lastenheft (Änderung
**0.10.0 → 0.11.0**): die Zugangsdaten des Agenten in [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) und die Ausweisung des Passwort-Logins von `sshd` in [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil).
Die Code-Slices der Welle beginnen erst danach.

**Gefundener Stand im Lastenheft und im Code (Ausgangslage, aus dem Quelltext gelesen):**

- [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) zählt die Sicherheitslockerungen auf (Capabilities, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices; mit [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) und [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)
  `--privileged` und `NET_ADMIN`) und sagt über Features nur: [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) gelten unverändert, das Profil gibt keine externe Quelle implizit frei.
  Ein Login-Dienst im Container (`sshd`) ist weder Lockerung im Sinn der Aufzählung noch ausgeschlossen.
- [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) gilt für **Git**-Zugangsdaten. Im Code passt dazu: `templates/devcontainer/devcontainer.json.tmpl` reicht im Sandbox-Zweig nur `GIT_TOKEN` als
  `remoteEnv` durch (ohne Sandbox gibt es kein `remoteEnv`), `credentialPatterns()` in `doctor_sandbox.go` erkennt Git-Token-Muster, `hasTokenSource` prüft nur
  `GIT_TOKEN` und einen Secret-Mount. Ein Anthropic-Schlüssel ist nirgends erfasst.
- Ob eine Warnung zur Befehlsausgabe unter `--quiet` sichtbar bleibt, entscheidet das Feld `SecurityRelaxation` des `WarningEntry`
  (`sandboxRelaxationWarnings` in `devcontainer_sandbox.go`; Ausnahme 1 von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) nennt Capabilities, `--privileged`, Seccomp-/AppArmor-Anpassungen und
  zusätzliche Devices). Die Drittanbieter-Warnung (`thirdPartyFeatureWarnings`) setzt das Feld nicht.
- Katalog-Schlüssel allein sind keine Zusage: `python`, `sdkman` und `homebrew` kamen ohne Lastenheft-Änderung
  ([`slice-v1-devcontainer-features-homebrew-sdkman`](../done/slice-v1-devcontainer-features-homebrew-sdkman.md),
  [`slice-v1-devcontainer-feature-python`](../done/slice-v1-devcontainer-feature-python.md)); [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) führt Beispiele auf, keine abschließende Liste.

**Entscheidungen des Projektinhabers am 2026-10-17** (gewählte Option, Begründung, bestätigt):

| # | Frage | Gewählt | Begründung | Folge für den Vertrag |
|---|---|---|---|---|
| 1 | `sshd` im Sandbox-Profil | erlaubt, die Befehlsausgabe weist den Passwort-Login aus (Warnung); die Teilfrage unten ist offen | Das Profil ist Schadensbegrenzung, kein Verbot von Diensten; der Nutzer entscheidet bewusst und wird informiert. | neue Ausweisungs-Zusage in [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil); bei Teilfrage „ja“ zusätzlich Ausnahme 1 in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |
| 2 | Claude Code im Sandbox-Profil und seine Egress-Hosts | erlaubt, Hosts gemessen (erwartet: npm-Registry beim Build, `api.anthropic.com` zur Laufzeit; Anmeldung, Telemetrie, Update und Anmeldeweg gemessen, nicht geraten) | Der Agent ist der Anlass des Sandbox-Profils; seine Hosts gehören in die Default-Allowlist je Feature. | keine Lastenheft-Änderung (die Hosts stehen in [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), nicht im Lastenheft) |
| 3 | Einstufung des Features von Anthropic | Drittanbieter-Feature nach [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) (Warnung bei Aktivierung, Herausgeber `anthropics`) | [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) trägt die Regel bereits; eine zweite Regel wäre eine Parallelwelt. | keine; keine Folge-ADR |
| 4 | Zugangsdaten des Agenten | [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) wird auf Zugangsdaten des Agenten erweitert: nie im Image, Workspace-Volume, in `u-boot.yaml` oder einer erzeugten Datei; Übergabe zur Laufzeit per Umgebungsvariable; `doctor` prüft die Quelle wie bei Git-Tokens | Das Sandbox-Profil führt mit dem Agenten einen zweiten Schlüsseltyp ein; eine Zusage gehört in das Lastenheft, nicht in die Spezifikation. | neue Zusage in [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) |

**Offene Teilfrage zu Entscheidung 1 (der Projektinhaber entscheidet beim Review des CR):** Ein Passwort-Login im Container ist eine Ausweisung, die der Slice vorschlagen darf.
Offen sind der Wortlaut der Warnung und ob sie zur Ausnahme 1 von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Sicherheitslockerung; im Code das Kennzeichen `SecurityRelaxation` am `WarningEntry`) zählt und damit unter `--quiet` sichtbar bleibt:

- Variante „ja“: die Warnung ist eine Sicherheitslockerung; Ausnahme 1 nennt den Login-Dienst, die Warnung trägt `SecurityRelaxation`, und sie bleibt unter `--quiet`.
- Variante „nein“: gewöhnliche Warnung; sie entfällt unter `--quiet` (die vier Ausnahmen bleiben unverändert).

Der Vorschlag für den Inhalt der Warnung: Der Login-Dienst `sshd` ist aktiv (Port 2222, nur Passwort-Login nach `sudo passwd`); das Profil begrenzt Schaden und ist keine Isolationsgrenze.
Ein sshd-Port wird nicht automatisch weitergeleitet: `forwardPorts` stammt nur aus den Services ([`LH-FA-DEV-005`](../../../../spec/lastenheft.md#lh-fa-dev-005--ports)).

**Kern des CR — Vorschlag für den erweiterten Text von [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)** (Stil und Aufbau des heutigen Textes; im Lastenheft stehen die Verweise als Anker `#…` der Datei; die Überschrift bleibt, damit die Anker im Repo stabil sind; das Lastenheft nennt keine Slice- oder Wellen-Kennung):

> *Priorität: V1.* Im Sandbox-Profil ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)) dürfen Git-Zugangsdaten und Zugangsdaten eines Agenten (etwa der Zugangsschlüssel eines Coding-Agenten, den ein aktiviertes Devcontainer-Feature nutzt) weder im Image, im Workspace-Volume noch in `u-boot.yaml` oder einer anderen erzeugten Datei stehen.
>
> - Übergabe zur Laufzeit per Umgebungsvariable; Git-Zugangsdaten alternativ als schreibgeschützter Secret-Mount. Die erzeugte Konfiguration trägt nur die Referenz auf die Umgebungsvariable des Hosts, nie den Wert.
> - Die Dokumentation beschreibt kurzlebige, auf das Repository begrenzte Tokens (privater Schlüssel nie im Container) und die Forderung nach Branch-Protection auf dem Remote; für den Agenten nennt sie die Umgebungsvariable, aus der der Zugangsschlüssel stammt.
> - `u-boot doctor` prüft die Quelle der Zugangsdaten (Git-Token-Quelle, bei aktivem Agenten-Feature die Umgebungsvariable des Agenten): fehlt sie oder liegt sie im Klartext in einer Projektdatei, `warn` (kein Abbruch).

Vorschlag für den Zusatz in [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), Liste „muss“ (der Wortlaut hängt an der Teilfrage): *„- einen im Container aktivierten Login-Dienst (Passwort-Login) in der Befehlsausgabe ausweisen.“*
Vorschlag für die Historie (§16): Zeile 0.11.0, 2026-10-17, Quelle „Vereinbarung mit dem Projektinhaber“: [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer): Zusage auf Zugangsdaten eines Agenten erweitert (nie in Image, Workspace-Volume, `u-boot.yaml` oder erzeugter Datei; Übergabe per Umgebungsvariable; `doctor` prüft die Quelle); [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): Ausweisung eines Login-Dienstes im Sandbox-Profil.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code, Katalog und Egress-Hosts** — die Code-Slices [`slice-v1-devcontainer-feature-sshd`](slice-v1-devcontainer-feature-sshd.md) und
  [`slice-v1-devcontainer-feature-claude-code`](slice-v1-devcontainer-feature-claude-code.md) setzen um, was dieser CR festlegt; der CR ändert nur den Vertrag.
- **`spezifikation.md` und `architecture.md`** — die Spezifikation beschreibt das Verhalten erst, wenn Lastenheft und Code es tragen
  ([`slice-v2-spec-devcontainer-sshd-claude-code`](slice-v2-spec-devcontainer-sshd-claude-code.md)); die Sicht nennt keine Features.
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
| [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | Zugangsdaten des Agenten | angepasst: Erweiterung (siehe Kern des CR); die Git-Zusagen bleiben wörtlich |
| [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (vier Ausnahmen unter `--quiet`) | Sichtbarkeit der sshd-Warnung | offen: hängt an der Teilfrage zu Entscheidung 1; Variante „nein“ lässt den Text unverändert |
| [`LH-NFA-SEC-001`](../../../../spec/lastenheft.md#lh-nfa-sec-001--keine-secrets-einchecken) | Referenz auf die Umgebungsvariable in `devcontainer.json` | verträglich: die Datei trägt die Referenz, nie den Wert |
| [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | Anthropic-Installationsskript beim Image-Build | verträglich: das Feature ist externer Code aus einer nicht freigegebenen Quelle und braucht die Zustimmung nach [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features); die Drittanbieter-Warnung weist es aus |

Die Zeilen sind gegen den Normtext gelesen; der Reviewer des CR bestätigt sie.

## 2. Definition of Done

- [ ] Die Entscheidungen vom 2026-10-17 stehen in §1 (Entscheidungsdokumentation); die Prüfungstabelle ist bestätigt, die Teilfrage zu Entscheidung 1 ist vom Projektinhaber entschieden und in §1 festgehalten.
- [ ] Lastenheft **0.11.0**: die Zusagen stehen in [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Agenten-Zugangsdaten) und [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Ausweisung des Login-Dienstes), bei Teilfrage „ja“ zusätzlich in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging); Versionsfeld und Historie (§16, Quelle „Vereinbarung mit dem Projektinhaber“) tragen den Eintrag; die Suche in §1 belegt, dass keine ältere Zusage widerspricht.
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
| `spec/lastenheft.md` (§3 [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), bei Teilfrage „ja“ [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Versionsfeld, §16 Historie) | update | die entschiedenen Zusagen; Historie nach dem Muster der CRs zu [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |
| dieser Slice (§1 Entscheidungsdokumentation, Prüfungstabelle, §7) | update | Entscheidungsdokumentation |

Kein Code, keine Spezifikation, keine Tests; `make docs-check` und `make gates` laufen für die Dokumente.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; die Entscheidungen 1 bis 4 des Projektinhabers (2026-10-17) liegen vor, die Teilfrage zu Entscheidung 1 ist beim Review des CR entschieden; kein anderer Slice in `in-progress/`
(WIP-Limit 1; solange [`slice-v2-homebrew-smoke-haertung`](../in-progress/slice-v2-homebrew-smoke-haertung.md) dort liegt, beginnt kein Slice der Welle).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Text verlangt eine neue Anforderung statt der Änderung von zwei bestehenden — dann ein CR je Anforderung.
- `in-progress` → `open` (blockiert): Die Teilfrage zu Entscheidung 1 ist nicht entschieden, eine Entscheidung wechselt, oder die Messung in einem Code-Slice findet eine Lockerung, die der CR noch nicht trägt (der CR wird dann erweitert, bevor der Code-Slice weiterläuft).

## 5. Closure-Trigger

DoD vollständig (Lastenheft 0.11.0 mit Historie), Review-Report liegt vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Die Teilfrage zu Entscheidung 1 bleibt offen** (Wortlaut der Warnung, Sichtbarkeit unter `--quiet`): der Code-Slice [`slice-v1-devcontainer-feature-sshd`](slice-v1-devcontainer-feature-sshd.md) kann die Warnung nicht festlegen. — **Ausgang:** weiter offen bis zum Review des CR; die DoD 1 verlangt die Entscheidung vor dem Move nach `done/`.
- **Der Anmeldeweg des Agenten ist nicht die Umgebungsvariable** (Anmeldung im Container statt Schlüssel vom Host): die `doctor`-Warnung „Quelle fehlt“ wäre dann ein Fehlalarm. — **Ausgang:** weiter offen: die Messung in [`slice-v1-devcontainer-feature-claude-code`](slice-v1-devcontainer-feature-claude-code.md) belegt den Weg; weicht er ab, wird der Text von [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) im CR vor dem Code-Slice nachgezogen (Rückführung `in-progress` → `open`).
- **Die Features bringen Metadaten mit** (`mounts`, `capAdd`, `privileged`, `entrypoint` in ihrer `devcontainer-feature.json`) und unterlaufen damit „kein Host-Bind-Mount, kein `--privileged`“ des Profils; weder Plan noch Lastenheft sagen es voraus. — **Ausgang:** weiter offen: die Messung der Code-Slices entscheidet; findet sie eine Lockerung, wird dieser CR vor dem Code-Slice um die Ausweisung erweitert (Rückführung `in-progress` → `open`).
- **Widerspruch zu einer älteren Zusage** bleibt unentdeckt, weil die Suche zu eng war. — **Ausgang:** entfallen: die Prüfungstabelle in §1 trägt die Suche (Stand 2026-10-17); der Review bestätigt sie.

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `spec/` (Kürzel DOK) — Modus Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006) (Abschnitt Modus-Deklaration), erfüllt die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. Code-Sub-Areas berührt der CR nicht.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer für diese Sub-Area (die einzige Beobachtung `BEO-CLI/kommentar-beschreibt-abwesendes` betrifft `internal/adapter/driving/cli` und ist verkörpert in `AGENTS.md §3.9`).

Alle berührten Sub-Areas GF.
