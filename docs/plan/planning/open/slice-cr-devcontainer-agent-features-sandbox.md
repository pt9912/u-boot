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

**Bezug:** [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil; Änderung), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten; Änderung, falls Frage 4 (b) gilt), [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Features, Zustimmung), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts des Agenten), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (nur bei Teilfrage „ja“ zu Frage 1). Verwandt: [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), ggf. [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)), §16 (Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt der CR nicht; die Spezifikation folgt im letzten Slice der Welle.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Entscheidungen des Projektinhabers zu `sshd` und Claude Code im Sandbox-Profil sind festgehalten, und soweit sie eine neue Zusage tragen, steht sie
im Lastenheft (Änderung **0.10.0 → 0.11.0**). Verlangt keine Entscheidung eine neue Zusage, entfällt der Slice (§7 `Gegenstand:`), und die Code-Slices der Welle
beginnen ohne Lastenheft-Änderung.

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

**Entscheidungen des Projektinhabers (offen, nicht beschlossen):** Der Slice liegt in `open/`, bis sie fallen; der Wortlaut kommt mit Datum hierher.

| # | Frage | Optionen | Folge für den Vertrag |
|---|---|---|---|
| 1 | `sshd` im Sandbox-Profil | (a) erlaubt; (b) erlaubt, die Befehlsausgabe weist den Passwort-Login aus (Teilfrage: zählt die Warnung zur Ausnahme 1 unter `--quiet`?); (c) im Sandbox-Profil ausgeschlossen (fachlicher Fehler, Exit-Code `10`) | (a) keine; (b) neue Ausweisungs-Zusage in [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), bei der Teilfrage „ja“ zusätzlich Ausnahme 1 in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging); (c) neue Zusage in [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) |
| 2 | Claude Code im Sandbox-Profil und seine Egress-Hosts | (a) erlaubt, Hosts in der Default-Allowlist je Feature; (b) erlaubt mit Ausweisung; (c) ausgeschlossen | (a) keine (die Hosts stehen in [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), nicht im Lastenheft); (b), (c) wie bei Frage 1 |
| 3 | Einstufung des Features von Anthropic | (a) Drittanbieter-Feature nach [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) (Warnung bei Aktivierung, Herausgeber `anthropics`); (b) eine andere Regel | (a) keine; (b) Folge-ADR, die [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) Entscheidung 1 und 2 schärft |
| 4 | Zugangsdaten des Agenten | (a) nur Handbuch und `doctor`, kein Lastenheft; (b) [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) wird auf Zugangsdaten des Agenten erweitert: nie im Image, Workspace-Volume, in `u-boot.yaml` oder einer erzeugten Datei, Übergabe zur Laufzeit per Umgebungsvariable | (a) keine; (b) neue Zusage in [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) |

**Empfehlung der Planung (kein Beschluss):** Die Frage 4 (b) ist wahrscheinlich: Sie überträgt die Zusage von [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) auf den zweiten Schlüsseltyp, den das Sandbox-Profil
neu einführt, und eine Zusage gehört in das Lastenheft, nicht in die Spezifikation. Bei Frage 1 und 2 genügt (a) ohne Lastenheft-Änderung nur dann, wenn die Messung in den
Code-Slices keine Lockerung findet (Feature-Metadaten mit `mounts`, `capAdd`, `privileged` oder `entrypoint`; siehe §6).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code, Katalog und Egress-Hosts** — die Code-Slices [`slice-v1-devcontainer-feature-sshd`](slice-v1-devcontainer-feature-sshd.md) und
  [`slice-v1-devcontainer-feature-claude-code`](slice-v1-devcontainer-feature-claude-code.md) setzen um, was dieser CR festlegt; der CR ändert nur den Vertrag.
- **`spezifikation.md` und `architecture.md`** — die Spezifikation beschreibt das Verhalten erst, wenn Lastenheft und Code es tragen
  ([`slice-v2-spec-devcontainer-sshd-claude-code`](slice-v2-spec-devcontainer-sshd-claude-code.md)); die Sicht nennt keine Features.
- **Eine Aufnahme von `sshd` oder Claude Code in die Beispielliste von [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)** — Beispiele sind keine Zusage; der Katalog führt sie ohne Lastenheft-Änderung (Muster `python`).
- **Schlüssel-Login, Passwort-Setzen und Agent-Konfiguration durch `u-boot`** — [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) nennt den Agenten als Nicht-Ziel; ein Passwort setzt der Nutzer im Container.
- **Eine ADR** — Frage 3 (b) löst, falls sie fällt, eine Folge-ADR im Slice [`slice-v1-devcontainer-feature-claude-code`](slice-v1-devcontainer-feature-claude-code.md) aus; hier entscheidet sie nichts.

### Prüfung gegen die älteren Zusagen (Steering-Loop-Lehre der vorigen CRs: vor der Änderung per Suche im Lastenheft)

Die Tabelle wird beim Schreiben des CR ausgefüllt; jede Zeile nennt Fundstelle und Ergebnis (verträglich / angepasst / Widerspruch).

| Stelle | Berührung | Ergebnis |
|---|---|---|
| [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (externe Quellen, Zustimmung) | Anthropic-Feature im Katalog | offen |
| [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Lockerungen ausweisen, Features unverändert) | `sshd`, Claude Code im Profil | offen |
| [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) und [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) | Degradations-Fallbacks, Hosts des Agenten | offen |
| [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | Zugangsdaten des Agenten | offen |
| [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (vier Ausnahmen unter `--quiet`) | Sichtbarkeit einer sshd-Warnung | offen |
| [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | Anthropic-Installationsskript beim Image-Build | offen |

## 2. Definition of Done

- [ ] Die Entscheidungen zu Frage 1 bis 4 stehen mit Datum und Wortlaut in §1; die Prüfungstabelle ist ausgefüllt.
- [ ] Lastenheft **0.11.0**: die Zusagen aus den Entscheidungen stehen in [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und, falls Frage 4 (b) gilt, in [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (bei „Teilfrage ja“ in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)); Versionsfeld und Historie (§16) tragen den Eintrag; die Suche in §1 belegt, dass keine ältere Zusage widerspricht. Verlangt keine Entscheidung eine neue Zusage: §7 trägt `Gegenstand: entfallen` mit den Entscheidungen, die Liefer-Punkte bleiben leer.
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
| `spec/lastenheft.md` (§3 [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), ggf. [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) und [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Versionsfeld, §16 Historie) | update | die entschiedenen Zusagen; Historie nach dem Muster der CRs zu [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |
| dieser Slice (§1 Entscheidungen, Prüfungstabelle, §7) | update | Entscheidungs-Protokoll |

Kein Code, keine Spezifikation, keine Tests; `make docs-check` und `make gates` laufen für die Dokumente.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; die Entscheidungen zu Frage 1 bis 4 des Projektinhabers liegen vor; kein anderer Slice in `in-progress/`
(WIP-Limit 1; solange [`slice-v2-homebrew-smoke-haertung`](../in-progress/slice-v2-homebrew-smoke-haertung.md) dort liegt, beginnt kein Slice der Welle).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Entscheidung verlangt eine neue Anforderung statt der Änderung von zwei bestehenden — dann ein CR je Anforderung.
- `in-progress` → `open` (blockiert): Eine Entscheidung fehlt oder wechselt, oder die Messung in einem Code-Slice findet eine Lockerung, die der CR noch nicht trägt (der CR wird dann erweitert, bevor der Code-Slice weiterläuft).
- `open` → `done` (ohne Arbeit): Der Projektinhaber verlangt keine neue Zusage; §7 nennt `Gegenstand: entfallen`.

## 5. Closure-Trigger

DoD vollständig (Lastenheft 0.11.0 mit Historie, oder `Gegenstand: entfallen`), Review-Report liegt vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Der CR schwächt oder verschärft das Sandbox-Profil** (etwa `sshd` im Profil ausgeschlossen): ein Fehlerpfad mit Exit-Code `10` im Code-Slice. — **Ausgang:** weiter offen bis zur Entscheidung des Projektinhabers; bei Eintritt trägt der Code-Slice [`slice-v1-devcontainer-feature-sshd`](slice-v1-devcontainer-feature-sshd.md) die Folge.
- **Die Features bringen Metadaten mit** (`mounts`, `capAdd`, `privileged`, `entrypoint` in ihrer `devcontainer-feature.json`) und unterlaufen damit „kein Host-Bind-Mount, kein `--privileged`“ des Profils; weder Plan noch Lastenheft sagen es voraus. — **Ausgang:** weiter offen: die Messung der Code-Slices entscheidet; findet sie eine Lockerung, wird dieser CR vor dem Code-Slice um die Ausweisung erweitert (Rückführung `in-progress` → `open`).
- **Widerspruch zu einer älteren Zusage** bleibt unentdeckt, weil die Suche zu eng war. — **Ausgang:** entfallen: die Prüfungstabelle in §1 trägt die Suche; sie wird bei der Closure bestätigt.

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `spec/` (Kürzel DOK) — Modus Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006) (Abschnitt Modus-Deklaration), erfüllt die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. Code-Sub-Areas berührt der CR nicht.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer für diese Sub-Area (die einzige Beobachtung `BEO-CLI/kommentar-beschreibt-abwesendes` betrifft `internal/adapter/driving/cli` und ist verkörpert in `AGENTS.md §3.9`).

Alle berührten Sub-Areas GF.
