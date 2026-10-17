# Slice CR: Degradations-Fallbacks werden dort ausgewiesen, wo sie ermittelt werden (`LH-FA-DEV-007`)

> **CR-Vehikel.** u-boot führt kein separates CR-Artefakt — dieser Slice
> **ist** der Change Request am Vertrags-Stratum (Precedent:
> [`slice-cr-cli-quiet-vorschau`](../done/slice-cr-cli-quiet-vorschau.md),
> [`slice-cr-sandbox-docker-runtime`](../done/slice-cr-sandbox-docker-runtime.md)).
> Anlass: Der zweite Review des Code-Slices
> [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md)
> ([Report](../../../reviews/2026-10-16-slice-v1-cli-quiet-nur-fehler-2.md), F-4) und die Verifier-Auflage 3
> ([Verification Evidence](../../../reviews/2026-10-16-verification-slice-v1-cli-quiet-nur-fehler.md)) fanden eine Lücke zwischen
> [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)
> („jeder Fallback wird in der Befehlsausgabe und in `u-boot doctor` ausgewiesen“) und dem, was der Host ermitteln kann.
> Der Projektinhaber hat am 2026-10-17 entschieden: **B — Lastenheft präzisieren**. Der Code-Zweig (Lösung A, Kanal Container → Host) entfällt.

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — Vertragspflege ohne Release-Bezug, keine Closure-Bedingung, die von der DoD dieses Slice verschieden wäre.

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Änderung),
[`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Ausnahme 2, Änderung),
[`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Wegfall der Egress-Restriktion als Fallback, Klarstellung).
Übernimmt den Gegenstand von
[`slice-v1-sandbox-fallback-ausweisung`](../done/slice-v1-sandbox-fallback-ausweisung.md).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
[`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)), §16 (Historie), Versionsfeld;
`spezifikation.md` ([`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten),
[`LH-FA-DEV-008.a`](../../../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability),
[`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), §7 Historie).
`architecture.md` berührt der CR nicht.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** Lastenheft-Änderung **0.12.0 → 0.13.0**: [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) verspricht die
Ausweisung eines Degradations-Fallbacks nur noch dort, wo ihn der Host ermittelt; was erst im Container entsteht, gibt das erzeugte Startscript im Container aus.
Spezifikation, Handbuch und CHANGELOG folgen derselben Aussage.

**Warum Präzisierung statt Code:** Die Zusage nennt drei Fallbacks, aber nur einer ist vom Host aus ermittelbar, und auch der nur auf einem Linux-Host. Ein Kanal Container → Host
(Statusdatei, Volume, Lesen des Container-Zustands durch `doctor`) wäre eine neue Schnittstelle am Sandbox-Profil
([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): kein Host-Socket) und eine neue Zusage; der Projektinhaber hat sie nicht gewollt (Entscheidung 2026-10-17).
Die Präzisierung nimmt eine Zusage zurück, die der Code nie erfüllte (siehe Ist-Stand) — sie ändert kein beobachtbares Verhalten.

### Gefundener Stand im Code (Ist-Stand, gelesen am 2026-10-17)

| Fallback | Wo er entsteht | Was ausgegeben wird | Fundstelle |
|---|---|---|---|
| `vfs`-Storage (`/dev/fuse` fehlt) | host-seitig ermittelbar nur auf Linux-Hosts: `doctor`-Item `devcontainer.sandbox.runtime` mit `DegradationFallback` (nur bei Policy `warn`; bei `fail` ist es ein Fehler ohne Kennzeichen); auf anderen Hosts (macOS/Colima: die Engine läuft in einer VM) meldet das Item OK und verweist auf das Startscript | `doctor`-Zeile | `internal/hexagon/application/doctor_sandbox.go`, `checkSandboxRuntime` (Linux-Prüfung `/proc/sys/kernel/osrelease`, `/dev/fuse`-Prüfung, Kennzeichen `DegradationFallback`) |
| `vfs`-Storage (im Container festgestellt) | `sandbox-init.sh` prüft `[ ! -c /dev/fuse ]` im Container, schreibt `storage.conf` mit `driver = "vfs"` | Zeile auf Standardfehlerausgabe `u-boot sandbox: /dev/fuse is not available; falling back to vfs storage`; Exit-Code des Skripts bei `warn` unverändert (kein `exit`), bei `fail` `11` | `internal/hexagon/application/templates/devcontainer/sandbox-init.sh.tmpl`, Funktion `unavailable`, Zeilen 7–12 und 20–25; Einhängung als `postCreateCommand` in `devcontainer_sandbox.go` (`composePostCreate`, Aufruf `sh /usr/local/bin/u-boot-sandbox-init`) |
| Container ohne Docker-Daemon | `docker-init.sh` startet `dockerd`; kommt der Daemon nicht hoch (Prozess endet oder nach 60 s nicht bereit) | Zeile auf Standardfehlerausgabe `u-boot sandbox: … no daemon available, see /var/log/u-boot-dockerd.log`; Exit-Code `0` bei `warn`, `11` bei `fail` | `docker-init.sh.tmpl`, Funktion `unavailable` (Zeilen 11–17), Aufrufe Zeilen 41 und 45; `postStartCommand` |
| Wegfall der Egress-Restriktion | `egress-init.sh` prüft `nft list tables`; ohne `NET_ADMIN` entfällt die Restriktion | Zeile auf Standardfehlerausgabe `u-boot sandbox: egress restriction unavailable (NET_ADMIN not granted); consider restricting egress at network or DNS level`; Exit-Code `0` bei `warn`, `11` bei `fail` | `egress-init.sh.tmpl`, Zeilen 14–20; `postStartCommand`, vor dem Daemon-Start eingereiht (`devcontainer_sandbox.go`, `postStart`) |

Weitere Befunde, gelesen an der Stelle:

- `doctor` hat kein Item für den Daemon- oder Egress-Fall; das Egress-Item `devcontainer.sandbox.egress` nennt ausdrücklich, dass die Capability beim Containerstart geprüft wird
  (`doctor_sandbox.go`, `checkSandboxEgress`, letzter Rückgabewert und Kommentar zu `checkIDSandboxEgress`).
- Für `init`, `generate` und `up` erzeugt keine produktive Quelle einen Eintrag mit `DegradationFallback`: Der Kommentar an `sandboxRelaxationWarnings`
  (`internal/hexagon/application/devcontainer_sandbox.go`) hält fest, dass die Fallbacks dort nicht entstehen. Die Kennzeichnung samt Filter unter `--quiet` ist vorhanden
  (`internal/adapter/driving/cli/warnings.go`, `visibleWarnings`; `internal/adapter/driving/cli/doctor.go`, `writeDoctorReport`) und über Fake-Einträge gepinnt (`internal/adapter/driving/cli/quiet_test.go`, `quietFallbacks`).
- Die Spezifikation formuliert die Zusage bereits enger als das Lastenheft: [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) Ausnahme 2 nennt „soweit sie vom Host aus ermittelbar sind“,
  [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) sagt „(soweit vom Host aus ermittelbar)“ bei der Prüfung und
  [`LH-FA-DEV-008.a`](../../../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability) hält fest, dass die Capability vom Host aus nicht zuverlässig bestimmbar ist.
  Das Lastenheft ist die einzige Stelle ohne diese Einschränkung; der Satz „jeder Fallback wird in der Befehlsausgabe und in `u-boot doctor` ausgewiesen“ steht im Lastenheft und, wortgleich, in der Spezifikation.
- Die Lücke besteht unabhängig von `--quiet`: Ausnahme 2 hält nur sichtbar, was `u-boot` ausgibt; was nicht ausgegeben wird, kann `--quiet` nicht verbergen.

### Prüfung gegen die älteren Zusagen (Steering-Loop-Lehre aus den CRs zu `LH-FA-CLI-005`)

Vor der Änderung geprüft per Suche im Lastenheft (`Fallback`, `ausgewiesen`, `ausweisen`, `Befehlsausgabe`, `Degradation`, `vfs`, `ohne Daemon`, `Wegfall`):

| Stelle | Berührung | Ergebnis |
|---|---|---|
| [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Absatz **Degradation**, letzter Satz | verspricht Ausweisung jedes Fallbacks in Befehlsausgabe und `doctor` | **geändert** (Stelle 1) |
| [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Ausnahme 2 | wiederholt dasselbe Versprechen für alle drei Fallbacks | **geändert** (Stelle 2) |
| [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), letzter Punkt | verweist auf die Degradation, „geprüft beim Containerstart; `u-boot doctor` prüft die Konfiguration“ | **klargestellt** (Stelle 3, optional); trägt die Änderung schon, ohne sie zu widersprechen |
| [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), Punkte „jede Sicherheitslockerung … in der Befehlsausgabe ausweisen“ | betrifft Lockerungen (host-seitig bekannt, `generate` kennt die `runArgs`), nicht Fallbacks | geprüft, unverändert |
| [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Punkt `docker` („die Befehlsausgabe weist das als Lockerung aus“) | Lockerung, kein Fallback | geprüft, unverändert |
| [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Punkte `podman` (Exit-Code `11` des Startscripts) und Degradations-Exit-Codes | beschreiben das Startscript bereits als Quelle des Exit-Codes | geprüft, unverändert; die neue Fassung folgt demselben Muster |
| [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Ausnahme 1 (Lockerungen, `doctor`-Sicherheits-Einträge) | host-seitig ausgewiesen | geprüft, unverändert |
| [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Code `11`) | Umgebungsproblem; das Startscript endet mit `11` im Container, `u-boot` selbst mit `11` bei `doctor`-Fehler | geprüft, unverändert |

Die Suche fand keine weitere Zusage, die „jeder Fallback in Befehlsausgabe und `doctor`“ wiederholt. Der Verifier wiederholt sie vor dem Schließen.

### CR-Kern (Wortlaut-Vorschlag; im Lastenheft ohne Markierung)

Die Fassungen stehen als Vorschlag. **Der Projektinhaber gibt den Wortlaut vor dem Lastenheft-Commit frei**; bis dahin bleibt das Lastenheft auf 0.12.0.

**Stelle 1 — [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Absatz Degradation, vorletzter und letzter Satz der Folge „Eine ausdrücklich angeforderte Runtime …“.**

Alt:

> Eine ausdrücklich angeforderte Runtime wird nie stillschweigend ersetzt; jeder Fallback wird in der Befehlsausgabe und in `u-boot doctor` ausgewiesen.

Neu:

> Eine ausdrücklich angeforderte Runtime wird nie stillschweigend ersetzt; jeder Fallback wird dort ausgewiesen, wo er ermittelt wird. Was der Host ermittelt (ein fehlendes `/dev/fuse`, das zum Fallback auf `vfs`-Storage führt), weist `u-boot doctor` aus, soweit der Host es bestimmen kann, und die Befehlsausgabe, wenn ein Befehl es feststellt. Ein Fallback, den der Host nicht ermittelt (Container ohne Docker-Daemon, Wegfall der Egress-Restriktion und `vfs`-Storage, sofern erst der Container das Fehlen von `/dev/fuse` feststellt), entsteht beim Containerstart und wird vom erzeugten Startscript im Container ausgegeben (Ausgabe des Containerstarts; Exit-Code des Startscripts: `0` bei `warn`, `11` bei `fail`); die Befehlsausgabe von `u-boot` und `u-boot doctor` weisen ihn nicht aus.

**Stelle 2 — [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Ausnahme 2.**

Alt:

> 2. die Ausweisung der Degradations-Fallbacks (Container ohne Docker-Daemon, `vfs`-Storage, Wegfall der Egress-Restriktion) in der Befehlsausgabe und in `u-boot doctor` ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer));

Neu:

> 2. die Ausweisung der Degradations-Fallbacks, soweit der Host sie ermittelt (`vfs`-Storage bei fehlendem `/dev/fuse`), in der Befehlsausgabe und in `u-boot doctor`; die Fallbacks, die erst im Container entstehen (Container ohne Docker-Daemon, Wegfall der Egress-Restriktion), gibt das Startscript im Container aus, sie liegen außerhalb der Ausgabe von `u-boot` und damit außerhalb dieser Ausnahme ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer));

Die Zählung „Vier Ausnahmen“ bleibt: Ausnahme 2 hat weiterhin einen Gegenstand (den `doctor`-Eintrag zu `vfs` unter `doctor --quiet`).

**Stelle 3 — [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), letzter Punkt (Nebenstelle, Klarstellung; entfällt, wenn der Projektinhaber sie nicht will).**

Alt:

> - Ist die nötige Capability nicht gewährbar, greift die Degradation aus [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (geprüft beim Containerstart; `u-boot doctor` prüft die Konfiguration).

Neu:

> - Ist die nötige Capability nicht gewährbar, greift die Degradation aus [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (geprüft beim Containerstart; das Startscript gibt den Wegfall der Restriktion im Container aus; `u-boot doctor` prüft die Konfiguration, die Capability kennt es nicht).

**Stelle 4 — Kopf und [§16 Historie](../../../../spec/lastenheft.md#16-historie).** Versionsfeld `0.12.0` → `0.13.0`; neue Zeile:

```text
| 0.13.0 | 2026-10-17 | [`LH-FA-DEV-007`](#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer): Die Ausweisung der Degradations-Fallbacks auf das beschränkt, was der Host ermittelt (`vfs`-Storage bei fehlendem `/dev/fuse`, in `u-boot doctor` und, wo ein Befehl es feststellt, in der Befehlsausgabe); Container ohne Docker-Daemon und Wegfall der Egress-Restriktion entstehen beim Containerstart und werden vom Startscript im Container ausgegeben (Exit-Code `0` bei `warn`, `11` bei `fail`), der Host kennt sie nicht. [`LH-FA-CLI-005`](#lh-fa-cli-005--verbosity-und-logging): Ausnahme 2 unter `--quiet` entsprechend; [`LH-FA-DEV-008`](#lh-fa-dev-008--egress-restriktion): Klarstellung zum Ort der Prüfung. | Vereinbarung mit dem Projektinhaber |
```

**Warum 0.13.0 und nicht 0.12.1:** Die Historie (§16) führt `0.3.x`-Patches für Klarstellungen, die einen Fall nachtragen, ohne eine Zusage zu ändern, und Minor-Wechsel, wenn eine Zusage
ihren Umfang ändert (`0.7.0`: Zusage für Docker Desktop zurückgenommen; `0.8.0`–`0.12.0`). Hier nimmt der CR einen Teil einer Zusage zurück („jeder Fallback in der Befehlsausgabe und in `doctor`“); das ist eine Minor-Änderung nach dem Muster von `0.7.0`.
Die Wahl ist vom Projektinhaber zu bestätigen.

**Unverändert, mit Begründung:**

- Die Konfigurationsschlüssel, die Degradationstabelle (`warn`/`fail`) und die Exit-Codes `10`/`11`: der CR ändert nur, *wer* den Fallback ausgibt, nicht *ob* er eintritt.
- Die Startscripte (`sandbox-init.sh`, `docker-init.sh`, `egress-init.sh`): sie geben aus, was der CR beschreibt; Änderungen am Code gibt es nicht.
- [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): Die Ausweisung der Lockerungen ist host-seitig und bleibt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Kanal Container → Host (Lösung A)** — vom Projektinhaber abgelehnt; er wäre eine neue Schnittstelle am Sandbox-Profil und eine neue Zusage.
- **Eine Änderung an `--quiet` oder am Code von `doctor`, `generate`, `init`, `up`** — das Verhalten ändert sich nicht; die Kennzeichnung `DegradationFallback` und ihre Pins bleiben bestehen, weil `vfs` unter `doctor --quiet` sie nutzt.
- **Eine Änderung der Startscript-Ausgabe** — der CR beschreibt sie nur; Wortlaut und Exit-Codes der Skripte bleiben.
- **Neue Fallback-Klassen oder Sicherheitslockerungen** — eine neue Zusage; [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) zählt drei Fallbacks.
- **Der `--device=/dev/fuse`-Befund (siehe §6)** — ein anderer Vorgang; er verlangt eine Messung, keine Präzisierung.

## 2. Definition of Done

- [ ] **Lastenheft 0.13.0**: Stellen 1, 2 (und, falls freigegeben, 3) in der vom Projektinhaber freigegebenen Fassung, Versionsfeld, Historie-Zeile (Verweis „Vereinbarung mit dem Projektinhaber“)
      ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)).
- [ ] **Spezifikation nachgezogen**: [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten),
      [`LH-FA-DEV-008.a`](../../../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability) und
      [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) Ausnahme 2 tragen dieselbe Aussage; Historie-Zeile in `spezifikation.md` §7.
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` (`--quiet`-Ausnahme 2, Degradationstabelle, Daemon- und Egress-Abschnitt) und CHANGELOG nennen den geltenden Stand der Ausweisung.
- [ ] `make gates` grün; `make doc-trace` ohne neue Waise.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Reihenfolge: Wortlaut-Freigabe des Projektinhabers → Lastenheft → Spezifikation → Handbuch und CHANGELOG → Gates → Review → Closure.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Stellen 1–3 ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)), Versionsfeld 0.13.0, §16-Zeile; erst nach Freigabe |
| `spec/spezifikation.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) (Satz „Jeder Fallback wird …“ und Prüfsatz), [`LH-FA-DEV-008.a`](../../../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability), [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) Ausnahme 2, §7 Historie |
| `docs/user/benutzerhandbuch.md` | update | `--quiet`-Ausnahme 2 (Tabellenzeile Allgemeine Optionen), Degradationstabelle (Zeilen um 1020), Docker-Daemon-Absatz (um 1073), Egress-Absatz (um 1131); `doctor --quiet`-Absatz (um 617) bleibt, weil `vfs` unter `doctor --quiet` sichtbar ist |
| `CHANGELOG.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Eintrag unter `[Unreleased]` / `Changed`; der Abschnitt `[0.10.0]` bleibt unverändert (veröffentlichter Stand) |
| `docs/plan/planning/open/README.md` | update | Zeile des übernommenen Slice durch diesen CR ersetzen |
| `docs/plan/planning/done/slice-v1-sandbox-fallback-ausweisung.md` | update (Inhalt), Move | Übergang mit `Gegenstand:` auf diesen CR; Move per `make slice-mv` nach dem Inhalt |
| Code und Tests | — | keine Änderung (siehe §1, „Ausdrücklich NICHT“); die Pins in `internal/adapter/driving/cli/quiet_test.go` bleiben gültig |

Nicht berührt, geprüft: `spec/architecture.md` (enthält weder Fallback-Zusage noch `--quiet`-Ausnahme), `README.md`/`README.de.md` (nennen „vier Ausweisungen“, nicht die drei Fallbacks), `docs/maintainer/cli-json-contract.md` §2.2 (beschreibt den `doctor`-Filter, der gilt weiter).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Projektinhaber gibt den Wortlaut der Stellen 1–4 frei (Entscheidung zu B liegt seit 2026-10-17 vor, der Wortlaut noch nicht).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Freigabe verlangt, weitere Zusagen mitzuregeln (etwa das `--device=/dev/fuse`-Verhalten, §6) — dann je Zusage ein Slice.
- `in-progress` → `open` (blockiert): Der Projektinhaber ändert die Entscheidung von B auf A (Code-Zweig); dann entsteht ein neuer Code-Slice, dieser CR entfällt.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7). Lastenheft 0.13.0 ist angenommen.

## 6. Risiken und offene Punkte

- **Der CR nimmt eine Zusage zurück:** „jeder Fallback in der Befehlsausgabe und in `doctor`“ gilt künftig nur für das Host-Ermittelbare. — **Ausgang:** offen bis zur Freigabe; bei Freigabe: entfallen: vom Projektinhaber ausdrücklich als CR entschieden (2026-10-17, „B: Lastenheft präzisieren“).
- **`vfs` ist nur auf Linux-Hosts host-seitig ermittelbar:** Auf macOS/Colima, den gemessenen Hosts, meldet `checkSandboxRuntime` OK und verweist auf das Startscript; der Wortlaut „soweit der Host es bestimmen kann“ trägt das, ohne die Plattformen zu nennen. — **Ausgang:** offen bis zur Freigabe; Frage an den Projektinhaber (siehe Fragen unten).
- **Spec ist dem Lastenheft voraus:** Die Spezifikation schränkt schon vor dem CR ein. Bleibt der CR aus, widerspricht ihr Satz dem Lastenheft (Lastenheft gewinnt). — **Ausgang:** entfallen: der CR zieht das Lastenheft nach, die Spezifikation wird im selben Slice präzisiert.
- **Der Befund `--device=/dev/fuse` (ungemessen):** `podmanRelaxations` (`devcontainer_sandbox.go`) setzt `--device=/dev/fuse`; fehlt das Device auf dem Container-Host, lehnt die Engine nach dem Muster von `/dev/net/tun` womöglich den Start ab, sodass der `vfs`-Zweig im Startscript nicht erreicht würde. Der Plan behauptet das nicht; er stützt sich nur auf den Wortlaut der Skripte. — **Ausgang:** weiter offen: → Frage an den Projektinhaber; bei Bestätigung ein eigener Slice mit Messung.
- **CHANGELOG-Release `[0.10.0]` trägt das alte Versprechen:** Der veröffentlichte Abschnitt nennt „the disclosure of degradation fallbacks (container without daemon, `vfs` storage, dropped egress restriction)“ als sichtbar bleibende Ausweisung. — **Ausgang:** entfallen: Veröffentlichte Abschnitte bleiben; der Eintrag unter `[Unreleased]` stellt die Aussage richtig.

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `spec/`, `harness/`, `docs/plan/`
(eine Sub-Area, Kürzel DOK, Modus Greenfield laut
[`harness/conventions.md`](../../../../harness/conventions.md#mr-006), Abschnitt Modus-Deklaration). `docs/user/` und `CHANGELOG.md` sind öffentliche Verträge, keine Sub-Area; der Code
(`internal/hexagon/application`, `internal/adapter/driving/cli`) wird nicht geändert.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register
([`../observations/`](../observations/README.md)) durchgegangen — Verzeichnisse `BEO-CLI`, `BEO-DOK`, `BEO-E2E`, `BEO-WZG` liegen vor; keiner betrifft die Ausweisung von Fallbacks
(`BEO-DOK/plan-aussage-ueber-code` betrifft Plan-Aussagen über Code und wurde hier beachtet: jede Aussage in §1 steht mit Fundstelle).

Alle berührten Sub-Areas GF.
