# Slice v1: `doctor --quiet` zeigt die Sicherheits-Warnungen des Sandbox-Profils

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md) — zweiter Slice der Welle; ihm geht [`slice-cr-devcontainer-agent-features-sandbox`](../in-progress/slice-cr-devcontainer-agent-features-sandbox.md) voraus (Lastenheft 0.11.0, Ausnahme 1 von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)), ihm folgen [`slice-v1-devcontainer-feature-sshd`](slice-v1-devcontainer-feature-sshd.md), [`slice-v1-devcontainer-feature-claude-code`](slice-v1-devcontainer-feature-claude-code.md) und [`slice-v2-spec-devcontainer-sshd-claude-code`](slice-v2-spec-devcontainer-sshd-claude-code.md).

**Bezug:** [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Scope; Ausnahme 1), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Quelle der Zugangsdaten), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Restriktion), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Docker-Daemon), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil), [`LH-FA-DIAG-003`](../../../../spec/lastenheft.md#lh-fa-diag-003--fehlerklassifikation) (Klassifikation der `doctor`-Einträge). Herkunft: Entscheidung 5 des Projektinhabers vom 2026-10-17 im CR [`slice-cr-devcontainer-agent-features-sandbox`](../in-progress/slice-cr-devcontainer-agent-features-sandbox.md). Muster: [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) (Kennzeichen `DegradationFallback` am Eintrag).

**Berührte Spec-Stellen:** — (keine). [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Wirkung von `--quiet`, Ausnahme 1 und Zeile zu `doctor`) und [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) bleiben bis [`slice-v2-spec-devcontainer-sshd-claude-code`](slice-v2-spec-devcontainer-sshd-claude-code.md) unverändert; das Lastenheft ändert nur der CR.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** `u-boot doctor --quiet` zeigt neben den Fehler-Einträgen und den Einträgen eines Degradations-Fallbacks auch die Warn-Einträge, die eine Sicherheitseigenschaft des Sandbox-Profils betreffen (Ausnahme 1 von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Fassung des CR); die Kennzeichnung sitzt am Eintrag, nicht an der Check-Kennung, wie bei `SecurityRelaxation` und `DegradationFallback`.

**Gefundener Stand im Code (aus dem Quelltext gelesen):**

- `writeDoctorReport` (`internal/adapter/driving/cli/doctor.go`) <!-- d-check:ignore (Pfad-Zitat im Code-Span) --> schreibt unter `--quiet` nur Einträge mit `Severity == SeverityError` oder `DegradationFallback`; die Hilfe des Befehls sagt „errors and degradation fallbacks only“.
- `domain.Diagnostic` (`internal/hexagon/domain/diagnostic.go`) <!-- d-check:ignore (Pfad-Zitat im Code-Span) --> trägt nur `DegradationFallback`; die Struktur hat keine JSON-Tags, die Wire-Form entsteht in `mapDiagnosticsToJSON` aus `ID`, `Severity` und `Message`. Das Gegenstück bei Befehlsausgaben ist `driving.WarningEntry.SecurityRelaxation` (`json:"-"`), gesetzt in `sandboxRelaxationWarnings` (`devcontainer_sandbox.go`) und ausgewertet in `warnings.go` (CLI).
- Die Warn-Einträge der drei Sandbox-Prüfungen in `doctor_sandbox.go` und ihre Einstufung (Sicherheitseigenschaft des Sandbox-Profils: ja/nein):

| Check-Kennung | Warn-Eintrag | Sicherheitseigenschaft | Kennzeichen nach dem Slice |
|---|---|---|---|
| `devcontainer.sandbox.credentials` | möglicher Klartext-Zugangsschlüssel in einer Projektdatei | ja ([`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)) | `SecurityRelaxation` |
| `devcontainer.sandbox.credentials` | https-Klonquelle ohne Quelle des Tokens zur Laufzeit | ja ([`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)) | `SecurityRelaxation` |
| `devcontainer.sandbox.egress` | Egress-Restriktion zusammen mit `nestedRuntime: docker` (Umgehung, `--privileged` hebt sie auf) | ja ([`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion); zugleich die einzige Warnung zur Docker-Daemon-Lockerung) | `SecurityRelaxation` |
| `devcontainer.sandbox.egress` | `egress.enabled` ohne Sandbox-Profil („has no effect“) | ja (Entscheidung des Projektinhabers, 2026-10-17): die Restriktion greift nur im Sandbox-Profil; die Meldung unter `--quiet` zu verbergen, ließe ein falsches Sicherheitsgefühl stehen ([`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) | `SecurityRelaxation` |
| `devcontainer.sandbox.runtime` | `nestedRuntime` ohne Sandbox-Profil („has no effect“) | nein (Entscheidung des Projektinhabers, 2026-10-17): Konfigurationsmeldung | keines (entfällt unter `--quiet`) |
| `devcontainer.sandbox.runtime` | `/dev/fuse` fehlt, Politik `warn` (`vfs`-Storage) | nein: Degradations-Fallback, Ausnahme 2 | `DegradationFallback` (unverändert) |
| `devcontainer.sandbox.runtime` | `nestedRuntime: docker` (Daemon mit `--privileged`) | kein Warn-Eintrag: die Stufe ist OK, die Meldung nennt `--privileged`; die Lockerung weist die Befehlsausgabe aus ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)) | keines (Stufe unverändert OK) |
| `devcontainer.features.allowlist`, `devcontainer.features.drift` | Warn-Einträge zu Feature-Quellen: Feature außerhalb des Katalogs ohne `source:` oder ohne `enabled`-Schlüssel (Allowlist), Abweichung zwischen `u-boot.yaml` und `devcontainer.json` (Drift) | nein: sie betreffen die Konfiguration und die Herkunft der Features (Drittanbieter-Herkunft, Entscheidung 5 des CR), keine Eigenschaft des laufenden Sandbox-Profils | keines (entfällt unter `--quiet`) |

  Die Fehler-Einträge (Politik `fail`, AppArmor-Sysctl, Allowlist-Verstoß) zeigt `--quiet` ohnehin. Der Eintrag zu `nestedRuntime: docker` hat die Stufe OK und bleibt es: eine Stufenänderung ändert `--strict` und die JSON-Ausgabe und gehört nicht in diesen Slice (§6). Die übrigen `doctor`-Warn-Einträge außerhalb der Sandbox-Prüfungen (Docker-Version, Pflichtdateien, `compose.yaml`, `forwardPorts`) sind keine Sicherheitseigenschaft des Profils und bleiben unter `--quiet` unterdrückt.
- Die sshd-Warnung gehört zur Befehlsausgabe (`init`, `generate`, `up`) und nutzt das vorhandene `SecurityRelaxation` am `WarningEntry`; sie hängt nicht an diesem Slice (siehe §4).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Drittanbieter-Code-Warnung beim Image-Build** ([`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md)) — der Projektinhaber hat sie nicht gewählt; sie entfällt unter `--quiet` weiter (`thirdPartyFeatureWarnings` setzt kein Kennzeichen).
- **Neue `doctor`-Check-Codes** — sie wären an `coderegistry.go` und [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) gebunden (Drift-Test); der Slice kennzeichnet vorhandene Einträge und führt keine Kennung ein.
- **Die Warnung zum Login-Dienst `sshd` und die Prüfung der Agenten-Variable** — beide gehören zu den Feature-Slices ([`slice-v1-devcontainer-feature-sshd`](slice-v1-devcontainer-feature-sshd.md), [`slice-v1-devcontainer-feature-claude-code`](slice-v1-devcontainer-feature-claude-code.md)); sie setzen das Kennzeichen, das dieser Slice einführt.
- **Stufen oder Wortlaut der Meldungen, JSON-Ausgabe, Exit-Codes** — unverändert; das Kennzeichen steuert nur die Sichtbarkeit unter `--quiet` im Textbericht (`--quiet --json` ändert nichts).
- **Die Spezifikation** ([`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)) — sie folgt im letzten Slice der Welle ([`slice-v2-spec-devcontainer-sshd-claude-code`](slice-v2-spec-devcontainer-sshd-claude-code.md)), nachdem Lastenheft und Code es tragen.

## 2. Definition of Done

- [ ] **Kennzeichen und Filter:** `domain.Diagnostic` trägt ein Kennzeichen `SecurityRelaxation` (Eigenschaft des Eintrags, nicht in der JSON-Wire-Form); die vier Warn-Einträge der Tabelle in §1 (credentials ×2, egress ×2: mit `docker` und „has no effect“) setzen es; `writeDoctorReport` zeigt unter `--quiet` Fehler-Einträge, Einträge mit `DegradationFallback` und Warn-Einträge mit `SecurityRelaxation`, sonst nichts (kein Kopf, keine Summary-Zeile); Hilfetext und Kommentare nennen die geltende Zusage (`AGENTS.md §3.9`).
- [ ] **Tests und Mutationen:** Happy (Klartext-Schlüssel, fehlende Token-Quelle, Egress mit `docker` und Egress „has no effect“ bleiben unter `--quiet` sichtbar), Boundary (ohne `--quiet` unverändert; dieselbe Check-Kennung `devcontainer.sandbox.runtime` ist mit „has no effect“ nicht sichtbar, mit dem Fallback `vfs` als Degradations-Fallback sichtbar), Negative (nur OK-Einträge, `nestedRuntime` „has no effect“ oder Feature-Allowlist- und Drift-Warnungen: `--quiet` schreibt nichts; `--quiet --json` unverändert, das Kennzeichen erscheint nicht im Umschlag; `--strict` mit nur diesen Warn-Einträgen endet weiter mit `11`); der Test `TestDoctor_SandboxCredentialsWarnIsNoFallback` und `TestDoctor_SandboxFallbackFlagIsOnTheItem` nennen das neue Kennzeichen (kein Fallback, aber eingestuft). Mutationen fallen rot: Filter ohne den Zweig `SecurityRelaxation`, Kennzeichen an der falschen Stelle (`nestedRuntime` „has no effect“ gekennzeichnet; Egress „has no effect“ ohne Kennzeichen), Kennzeichen am `Severity`-OK-Eintrag, Kennzeichen an der Check-Kennung statt am Eintrag.
- [ ] **Doku:** Benutzerhandbuch (beide `--quiet`-Zeilen der Optionstabelle und der `doctor`-Abschnitt: welche Warn-Einträge unter `--quiet` bleiben) und `CHANGELOG.md` (`[Unreleased]`) nennen das Verhalten; der Satz „follow in a later release“ im CHANGELOG-Eintrag des CR (Lastenheft 0.11.0) ist für den Teil dieses Slice angepasst (der letzte landende Code-Slice der Welle entfernt ihn); `README.md` und `README.de.md` nur, wenn sie `doctor --quiet` nennen (heute nicht).
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
| `internal/hexagon/domain/diagnostic.go` | update | Feld `SecurityRelaxation` am `Diagnostic`, Doc-Kommentar mit Zusage und Abgrenzung zu `DegradationFallback` |
| `internal/hexagon/application/doctor_sandbox.go` | update | Kennzeichen an den vier Warn-Einträgen (credentials ×2, egress mit `docker`, egress „has no effect“); `nestedRuntime` „has no effect“ und `vfs`-Fallback unverändert |
| `internal/adapter/driving/cli/doctor.go` | update | Filter in `writeDoctorReport`, Hilfetext des Befehls, Doc-Kommentare |
| `internal/hexagon/application/doctor_sandbox_test.go`, `internal/adapter/driving/cli/quiet_test.go`, `internal/adapter/driving/cli/doctor_json_test.go` | update | Happy/Boundary/Negative und Pins; JSON unverändert — nach [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | update | öffentlicher Vertrag |

Die Pfade stammen aus dem Quelltext (Lesestand 2026-10-17); der Implementer prüft jede Zeile vor der Änderung gegen den Code und trägt Abweichungen in die Evidence ein. <!-- d-check:ignore (Pfad-Zitate in der Tabelle darüber) -->

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der CR [`slice-cr-devcontainer-agent-features-sandbox`](../in-progress/slice-cr-devcontainer-agent-features-sandbox.md) liegt mit Lastenheft 0.11.0 in `done/`; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Reihenfolge in der Welle (aus dem Code gelesen):** Nach dem CR, weil die Ausnahme 1 erst dort im Lastenheft steht. Vor [`slice-v1-devcontainer-feature-claude-code`](slice-v1-devcontainer-feature-claude-code.md) **zwingend**: dessen neue `doctor`-Warnung zur Agenten-Variable trägt das Kennzeichen `SecurityRelaxation`, das dieser Slice am `Diagnostic` einführt. Vor [`slice-v1-devcontainer-feature-sshd`](slice-v1-devcontainer-feature-sshd.md) aus **Wahl**, nicht aus Abhängigkeit: Die sshd-Warnung ist ein `WarningEntry` der Befehlsausgabe, dessen `SecurityRelaxation` seit der Ausnahme 1 existiert; der sshd-Slice könnte auch davor stehen. Die gewählte Reihenfolge schließt die Ausnahme 1 im Code (Befehlsausgabe und `doctor`), bevor ein Feature eine weitere Warnung unter sie stellt.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice überschreitet 3 Liefer-Punkte, etwa weil die Stufe des `nestedRuntime: docker`-Eintrags mitgeändert werden soll.
- `in-progress` → `open` (blockiert): die Einstufung der Warn-Einträge in §1 trägt den Wortlaut der Ausnahme 1 nicht (der Projektinhaber stuft anders ein) — dann zuerst der CR.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Der Wortlaut der Ausnahme 1 nennt „Sicherheitslockerungen des Docker-Daemons nach [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)“, doch der `doctor`-Eintrag zu `nestedRuntime: docker` hat die Stufe OK**; die einzige Warnung zur Docker-Daemon-Lockerung ist der Egress-Eintrag mit `docker` (nach [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)). — **Ausgang:** entfallen: der Slice kennzeichnet nur vorhandene Warn-Einträge; die Menge der Daemon-Warn-Einträge ist leer und die Einstufung steht in der Tabelle in §1 (Verification-Evidence [`2026-10-17-verification-slice-cr-devcontainer-agent-features-sandbox`](../../../reviews/2026-10-17-verification-slice-cr-devcontainer-agent-features-sandbox.md), Abschnitt (d)); eine Warn-Stufe für den Daemon-Eintrag wäre eine Vertragsänderung (`--strict`, JSON) und ein eigener Vorgang.
- **Neue Warn-Einträge späterer Slices vergessen das Kennzeichen** (Prüfung der Agenten-Variable, weitere `doctor`-Warnungen im Sandbox-Profil). — **Ausgang:** eingetreten: [`slice-v1-devcontainer-feature-claude-code`](slice-v1-devcontainer-feature-claude-code.md) trägt es in der DoD (Warn-Einträge zur Agenten-Variable mit `SecurityRelaxation`); für künftige Einträge gilt der Test je Eintrag (Boundary mit derselben Check-Kennung).
- **Plan-Aussage über den Code stimmt nicht** (Klasse „Plan-Aussage über Code stimmt nicht“, siehe §8): die Fundstellen in §1 und §3 sind ein Lesestand. — **Ausgang:** weiter offen: → [BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) im Register (verkörpert in `AGENTS.md §3.10`); der Implementer prüft jede Aussage vor der Änderung im Code, Abweichungen stehen in der Evidence.

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `internal/hexagon/domain` (DOM), `internal/hexagon/application` (APP), `internal/adapter/driving/cli` (CLI) und, für die Doku, `docs/user/` (DOK) — alle Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006) (Abschnitt Modus-Deklaration), alle erfüllen die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen — für CLI die Beobachtung `BEO-CLI/kommentar-beschreibt-abwesendes` (3×, verkörpert in `AGENTS.md §3.9`): der Kommentar-Selbstcheck gilt für Hilfetext und Doc-Kommentare des Slice (die Zusage nennen, nicht die frühere Fassung „errors and degradation fallbacks only“); sonst keine Treffer für DOM, APP und DOK. Die Beobachtung [BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (verkörpert in `AGENTS.md §3.10`) gilt für jede Plan-Aussage über Code; sie steht in §6.

Alle berührten Sub-Areas GF.
