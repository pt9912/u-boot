# Slice V1: Release-Cut `v0.18.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** keine Welle-Datei — Minor-Release `v0.18.0` von `main` nach `v0.17.0` (Release-Zuordnung und Einstufung entschieden: Minor, *Changed* — Projektinhaber 2026-11-01 im Slice [`slice-v1-devcontainer-container-name`](../done/slice-v1-devcontainer-container-name.md), §6 Frage 8; wie die Cuts ohne Welle-Datei, geführt in der [Roadmap](../in-progress/roadmap.md), bis der Tag gesetzt ist). Der Release trägt, was seit dem Tag `v0.17.0` (`f89424e1`) unter `## [Unreleased]` steht: die Auflösung von Git-Zugangsdaten und Identität des Hosts (opt-in, Lastenheft `0.28.0`, [`slice-v1-devcontainer-host-git-uebernehmen`](../done/slice-v1-devcontainer-host-git-uebernehmen.md), closed mit CI-Nachweis 41328395477), die Schlüssel `devcontainer.resources.memory`/`.cpus` (Lastenheft `0.27.0`, [`slice-v1-devcontainer-ressourcenlimits`](../done/slice-v1-devcontainer-ressourcenlimits.md)) und der Containername `u-boot-<project>[-<instance>]` (Lastenheft `0.29.0`, [`slice-v1-devcontainer-container-name`](../done/slice-v1-devcontainer-container-name.md), geschlossen mit CI-Nachweis 41340432215). Planungs-Commits ohne Release-Wirkung: der CR-Entwurf [`slice-v1-katalog-gepinnte-toolchain-entwurf`](../done/slice-v1-katalog-gepinnte-toolchain-entwurf.md) und die Messungs-Commits des Container-Name-Slices (Messprotokoll, kein Produktcode), die Messungs-Commits des Host-git-Slices (`tools/messung/` und `docs/messungen/`) und zwei `tools/mk`-Commits (`82682508`, `7bbb8185`: GIT_TOKEN/GH_TOKEN-Brücke für `make dev-up`/`dev-shell`, `INSTANCE`-Weiterleitung).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Containername und Instanz-Aufrufe), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten/Identität), [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (Ressourcen-Keys, HostGitToken-Schalter), [`LH-FA-BUILD-003`](../../../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) (Pins unverändert), [`ADR-0004`](../../adr/0004-ci-system.md), [`ADR-0007`](../../adr/0007-distributionswege-ghcr.md), [`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md), [`ADR-0017`](../../adr/0017-linux-pakete-nfpm.md), [ADR-0040](../../adr/0040-host-git-quellen-fuer-token-und-identitaet.md). Ablauf (einzig normativ): [`docs/maintainer/releasing.md`](../../../maintainer/releasing.md) §§3 bis 8 — Form des Artefakts aus der vendored Vorlage [`slice.template.md`](../../../../.harness/baseline/v6.18.1/templates/docs/plan/planning/slice.template.md).

**Berührte Spec-Stellen:** — (Spec-Version ≠ Release-Version; dieses Lastenheft `0.29.0` lieferte der CR des Container-Name-Slices; dieser Slice ändert `spec/` nicht).

**Verantwortlich:** pt9912 (Implementer; Umsetzung und Gate-Läufe durch den Hauptagenten im Auftrag)

**Autor:** pt9912 (Rolle Planner, im Auftrag des Hauptagenten). **Datum:** 2026-11-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** `v0.18.0` liefert den Block `[0.18.0]` des CHANGELOG als GHCR-Image, Binaries, Pakete und Homebrew-Formel. Der Slice bereitet vor; **Push und Tag führt nur der Hauptagent nach ausdrücklicher Freigabe des Projektinhabers im Gespräch**.

**Nummer:** Minor (Zusage-Menge wächst: zwei neue Konfigurationsbereiche und ein opt-in Schalter), mit *Changed* für den Container-Verhalten-Wechsel (`up` startet künftig benannte Container); entschieden 2026-11-01. Kein Vorab-Tag außer auf Wunsch.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Push, Tag und alles mit Außenwirkung durch den Implementer** ([`releasing.md`](../../../maintainer/releasing.md) *Zweck*); `tap-check` startet der Auftraggeber mit dem PAT.
- **`make test-docker` lokal** — Entscheidung des Projektinhabers (2026-11-01): der CI-Lauf `integration-docker` auf dem Push-Tip ist der Nachweis; der Monitor der beiden Läufe `41335734851`/`41335734895` trägt den Stand des Diffs.
- **Produktcode, `spec/`, ADRs** — der Cut ändert nur Release-Dokumente und drei Versionszeilen.
- **Offene Entwürfe** ([`slice-v1-katalog-gepinnte-toolchain-entwurf`](../done/slice-v1-katalog-gepinnte-toolchain-entwurf.md), [`slice-v1-katalog-feature-jdk-gepinnt`](../next/slice-v1-katalog-feature-jdk-gepinnt.md), [`slice-v1-docker-cli-bau-auf-standard-basis-image`](../open/slice-v1-docker-cli-bau-auf-standard-basis-image.md), die beiden Test-docker-Slices) — sind nicht Inhalt des Release und geben ihm keine Bedingung.
- **Hebung von `GO_VERSION`/`GOLANGCI_LINT_VERSION`** — nur Anlass bei rotem `image-scan` ([`releasing.md`](../../../maintainer/releasing.md) §4, §8).

**Vor dem Start zu klären** (Implementer, M-Schritt): jeder Eintrag unter `## [Unreleased]` trägt einen Beleg (`LH-*`, Slice) und einen Commit im Diff `v0.17.0..HEAD`; ein Eintrag ohne beides ist ein Befund, der vor dem Cut geklärt wird (gelesen: drei Einträge mit Trägern; die Hygiene-Einträge der Hilfe-Texte gehören zu einem älteren Commit und brauchen diesen Abgleich, nicht Annahme).

## 2. Definition of Done

- [ ] **Release-Text und Versionsstellen:** `CHANGELOG.md` trägt `## [0.18.0] - <Tag-Datum>` mit Lead (Stil `[0.17.0]`: „Eighteenth release.“), leerem Anker `[Unreleased]` und Compare-Links; Versionsstrings `0.18.0-dev` an den drei Stellen (`cmd/uboot/main.go`, `repo.mk`, `Dockerfile`); beide READMEs und das Handbuch (Kopf, Versionsbeispiele, Gültigkeitsbereich) tragen `v0.18.0`; Handbuch-`Stand` gleich dem Tag-Datum.
- [ ] **Vorab-Prüfungen:** `make gates`, `make ci`, `make fullbuild` lokal grün; Trockenlauf des Tap-Pfads nach [`releasing.md`](../../../maintainer/releasing.md) §4 (mit Gegenprobe auf den Vorgänger-Tag); `ci` und `integration-docker` auf dem Push-Tip grün; `tap-check` grün; **danach und nur nach Freigabe**: Tag `v0.18.0` gesetzt und gepusht; Kontrolle nach [`releasing.md`](../../../maintainer/releasing.md) §6.
- [ ] **Doku nachgezogen:** Roadmap (Zeiger → nach dem Tag *Abgeschlossene Wellen* mit Tag-Hash), Zeilen der Bestands-Readmes — der `releasing.md`-Zeiger „zuletzt" ist entfallen (gestrichen in `8cd0976b`, Plan-Tabelle „im Cut kein Diff").
- [ ] `make gates` grün nach dem letzten Commit vor dem Handoff und nach der Closure.
- [ ] Review (kein Self-Review) und Verification mit Reports unter `docs/reviews/`; Nachlauf-Review über jeden Diff von Vertragstext.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register fortgeschrieben oder „keine Beobachtung" notiert; jedes Risiko mit Ausgang; die drei Paarungen getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung (Fundstellen) |
|---|---|---|
| `harness/conventions/MR-003-roadmap-wellen.md`, `harness/conventions.md` (Index-Zeile) | update | MR-Anker der welle-losen Cut-Praxis — Beschluss 2026-11-02 (§6); MR-003 neu gefasst: Release-Cuts sind keine Wellen, Grundlagen vier welle-lose Cuts seit `v0.13.1` |
| `docs/plan/planning/in-progress/roadmap.md` (Format-Regel), `docs/plan/planning/README.md` (Wellen-Absatz), `.claude/commands/plan-welle.md`, `.claude/commands/close-welle.md`, `.claude/agents/planner.md` | update | Stellen, die „Wellen = Release-Versionen" wiederholten, folgen der Neufassung von MR-003 (§6 Beschluss); keine neuen Regeln, nur die wiederholte Aussage berichtigt |
| `CHANGELOG.md` | update | `[Unreleased]`-Inhalt (gelesen: drei Zusage-Einträge, [`AGENTS.md`](../../../../AGENTS.md) §3.10) unter `## [0.18.0] - <Tag-Datum>` mit Lead — Lead nennt die drei Nutzer-Wirkungen samt Grenzen der Lastenheft-Stände `0.27.0`/`0.28.0`/`0.29.0`; Fuß: `[Unreleased]: …/compare/v0.18.0...HEAD`, neue Zeile `[0.18.0]: …/compare/v0.17.0...v0.18.0` |
| `cmd/uboot/main.go`, `repo.mk`, `Dockerfile` | update | `0.17.0-dev` → `0.18.0-dev` je Stelle samt Kommentar (Fundstellen beim Start mit `grep -rIn '0\.17\.0-dev'` erheben) |
| `README.md`, `README.de.md` | update | Status-Zeile und `VERSION=`-Beispiel (README.md Zeilen 10 und 98, README.de.md Zeilen 11 und 101); keine Kennungen, kein „released" |
| `docs/user/benutzerhandbuch.md` | update | Kopf `Software-Version` und `Stand` (heut `2026-11-01`), Versionsbeispiele, Gültigkeitsbereich |
| `docs/plan/planning/in-progress/roadmap.md` | update | Zeiger auf die Umsetzung (Container-Name) weg nach Closure; Zeiger auf diesen Slice unter *Offene Wellen*; nach dem Tag *Abgeschlossene Wellen* mit Tag-Hash und Datum |
| `docs/plan/planning/open/README.md` | update | Bestand-Zeile für diesen Slice |
| dieser Slice | neu/Moves | `make slice-mv` je Lifecycle-Schritt; nach dem Tag nach `done/` |
| `docs/maintainer/releasing.md` | **im Cut kein Diff** | der Schnitt von "Muster der Release-Cut-Slices" zu "Form aus der Skelett-Vorlage, Ablauf dieses Dokuments" ist vor dem Cut gezogen (Commits dieser Korrektur); der Cut zieht `releasing.md` danach nicht nach |
| `spec/`, `docs/plan/adr/`, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` | kein Diff | Kein Pin-Bump, keine Pipeline-Änderung ([`AGENTS.md`](../../../../AGENTS.md) §3.6) |

### Reihenfolge der Maintainer-Aktionen (aus [`releasing.md`](../../../maintainer/releasing.md) §§4 bis 6 abgeleitet)

M0 Datum-Abgleich vor jedem Tag (`date +%F`, CHANGELOG-Überschrift und Handbuch-`Stand`) · M-T Trockenlauf des Tap-Pfads · M1 Push von `main` (nach Freigabe) · M2 CI-Läufe des Push-Tips lesen · M3 Vorab-Tag entfällt außer auf Wunsch · M4 `gh workflow run tap-check` (Auftraggeber, unmittelbar vor M5) · M5 Tag `v0.18.0` gesetzt und gepusht (nur nach Freigabe des Projektinhabers) · M6 Kontrolle nach [`releasing.md`](../../../maintainer/releasing.md) §6 · M7 Mac-Probe des Tap.

**Sensoren und Nachweise** aus [`releasing.md`](../../../maintainer/releasing.md) §4 und dem Gate-Index [`harness/README.md`](../../../../harness/README.md) §Sensors: `make gates` (inklusive `docs-check`), `make ci`, `make fullbuild`, `make build VERSION=0.18.0` mit `--version`- und OCI-Label-Gegenprobe (die Abweisung des falschen Versionswerts ist Step des `publish`-Workflows), Trockenlauf des Tap-Pfads (§4), `gh workflow run tap-check` (§2) und die grünen Läufe `ci`/`integration-docker` auf dem Push-Tip (§4); `verify-depguard` nur bei berührten Importregeln; `make doc-immutable RANGE="v0.17.0..HEAD"`, weil der Cut die Accepted-ADRs nicht ändert und der Nachweis das belegt (on-demand-Sensor).

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): die getragenen Slices liegen in `done/` — **der Container-Name-Slice liegt in `done/`** (Verification-Report 2026-11-01; `make gates` 2026-11-02 grün auf dem Tip `71cf2f23`); `Verantwortlich:` gesetzt; `in-progress/` trägt sonst nichts (WIP-Limit 1); die CI-Läufe des Diffs `33e9cd5d` sind grün oder der Cut wartet.

**Rückführungen — vorab benannt:**

- `in-progress` → `open` (blockiert): ein `--- FAIL` rot in `integration-docker` auf dem Push-Tip (dann ein eigener Fix-Slice vor dem Cut); ein roter `image-scan` (Go-Pin heben, eigener Vorgang); ein roter `tap-check`; ein Träger-Eintrag im CHANGELOG ohne Beleg (§1 Vor dem Start).
- `in-progress` → `next`: nicht vorgesehen (der Cut ist ein Slice mit drei Textstellen).

## 5. Closure-Trigger

DoD vollständig; Tag `v0.18.0` gesetzt und gepusht (M5, nur mit Freigabe); `publish` und Job `tap` kontrolliert (M6); Trockenlauf-Serge festgehalten; Closure-Notiz mit Tag-Hash, Lauf-IDs und Lerneintrag; der Slice geht **erst nach dem Tag** nach `done/` ([`AGENTS.md`](../../../../AGENTS.md) §3.3 Fall 2). Ohne Freigabe bleibt er in `in-progress/`.

## 6. Risiken und offene Punkte

**Entscheidung (2026-11-02, Projektinhaber):** Die Praxis bleibt welle-los; [`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md) ist neu gefasst (2026-11-02) und trägt sie als Adaption — Release-Cuts sind keine Wellen und bekommen keine Welle-Datei; die Closure lebt im `done/`-Cut-Slice (§7 mit Tag-Hash und Lauf-IDs) und der Roadmap-Zeile. Grundlage die dokumentierte Praxis der vier welle-losen Cuts (`v0.13.1`, `v0.15.1`, `v0.16.1`, `v0.17.0`). Der Drift-Audit-Kandidat ist damit aufgelöst; die Stellen, die die alte Aussage „Wellen = Release-Versionen" wiederholten, sind mitgezogen (Roadmap-Format-Regel, Planning-README, `plan-welle`/`close-welle`, `planner`).

1. **Der Lead sagt mehr als das Lastenheft** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)). Vorsorge: der Lead trägt je Eintrag die Bedingungen des Lastenheft-Stands und die Liste des Nicht-Zugesagten im Wortlaut des CHANGELOG-Eintrags; der Reviewer liest ihn gegen die drei Stände (`0.27.0`, `0.28.0`, `0.29.0`). — **Ausgang:** bei Closure.
2. **Ein Unreleased-Eintrag ohne Träger** (etwa der Hygiene-Eintrag der Hilfetexte). Vorsorge: §1-Abgleich mit `git diff v0.17.0..HEAD`. — **Ausgang:** bei Closure.
3. **Das Datum weicht vom Tag-Datum ab** (Mitternacht). Vorsorge: M0 vor M1 und M5. — **Ausgang:** bei Closure.
4. **Ein neuer Advisory macht `image-scan` rot.** Vorsorge: `make ci` lokal, Läufe des Push-Tips. — **Ausgang:** bei Closure.
5. **Job `tap` scheitert im stabilen Lauf.** Vorsorge: `tap-check` unmittelbar vor M5; Trockenlauf M-T. — **Ausgang:** bei Closure.
6. **`integration-docker` am Push-Tip rot oder am Zeitlimit.** Vorsorge: Monitor läuft; ein Fix-Runde gehört dem trägenden Slice. — **Ausgang:** bei Closure.
7. **Der `up`-Verhalten-Wechsel mit den benannten Containern überrascht Skripte, die den Container-Namen parsen** (`docker ps --format` mit Zufallsnamen-Konvention). Vorsorge: der Lead nennt die neue Form als *Changed* mit dem Hinweis auf die Such-Identität (Labels unverändert). — **Ausgang:** bei Closure.

## 7. Closure-Notiz

**Evidence (Berichte und Läufe):**

- Review-Report [`2026-11-02-slice-v1-release-cut-v0.18.0`](../../../reviews/2026-11-02-slice-v1-release-cut-v0.18.0.md): 0 HIGH, 2 MEDIUM (F-1, F-2), 8 LOW (F-3 bis F-10), 1 INFO (F-11); Plan-Tabelle gegen `git show --stat e6c87379` ohne Befund, Negativbefunde über die Klassen des Vorgänger-Cuts ohne neuen Fund auf den gelesenen Linien.
- Verification-Report [`2026-11-02-slice-v1-release-cut-v0.18.0-verify`](../../../reviews/2026-11-02-slice-v1-release-cut-v0.18.0-verify.md): DoD-Konform verifiziert, closure-fähig; alle Lauf-IDs und Digests vom Verifier selbst nachgelesen (`gh run view`, `gh release view`, Manifest-Inspekte), Gate-Stempel gegen Inhalts-Hash geprüft, Tap-Trockenlauf im Verifier-Lauf neu gefahren.
- Läufe: lokal `make gates` Exit 0 (Stempel = HEAD `724c5e5e`, Docker-Modus), `make ci` Exit 0, `make fullbuild` Exit 0; `make build VERSION=0.18.0` melde `u-boot version 0.18.0` und OCI-Label `0.18.0`; Tap-Trockenlauf v0.18.0 Exit 0 (Protokoll clone/add/`commit -m u-boot 0.18.0`/push), Gegenprobe v0.17.0 Exit 1 vor `clone`; `make doc-immutable RANGE="v0.17.0..HEAD"` grün (98 Commits nach den Nachzügen). CI auf dem Tag-Commit `e6c87379`: [`ci`](https://github.com/pt9912/u-boot/actions/runs/41403124364) `41403124364`, [`integration-docker`](https://github.com/pt9912/u-boot/actions/runs/41403124438) `41403124438`, [`tap-check`](https://github.com/pt9912/u-boot/actions/runs/41403131805) `41403131805` — je success. Nach dem Tag: [`publish`](https://github.com/pt9912/u-boot/actions/runs/41405100611) `41405100611` success (Jobs publish, package smoke deb+rpm, `tap` je success).
- Release: `u-boot v0.18.0`, 12 Assets, kein Draft, kein Prerelease, publiziert `2026-11-02T07:17:01Z`; Tag lightweight (`git cat-file -t` = `commit`) auf `e6c87379`. GHCR-Digests `:0.18.0` und `:latest` gleich (`sha256:bef5b03dff71467cc439ed015e98f0b8552e827ddef79789f8c19ca82ad50383`), `docker run --rm ghcr.io/pt9912/u-boot:0.18.0 --version` = `u-boot version 0.18.0`. Tap-Commit `b52a0bc` (`u-boot 0.18.0`).
- M7 (Mac): gelaufen, Aktion des Projektinhabers — Meldung im Hauptagenten-Gespräch 2026-11-02: `u-boot v0.18.0` läuft auf dem macOS. Der Job `tap` ersetzt die Bestätigung nicht; hier ist sie der Meldung nach erbracht.

**DoD je Liefer-Punkt (alle erbracht):**

1. **Release-Text und Versionsstellen** — `CHANGELOG.md` trägt `[0.18.0] - 2026-11-02` mit Lead (Tag-Datum = Datum), leerem Anker `[Unreleased]` und Compare-Links; `0.18.0-dev` an den drei Stellen samt Kommentaren; READMEs und Handbuch auf `v0.18.0`, Handbuch-`Stand` = Tag-Datum (Review-Negativbefunde zu Kopf, Fuß, Versionen und Kennungsfreiheit).
2. **Vorab-Prüfungen** — `make gates`, `make ci`, `make fullbuild` lokal grün; Trockenlauf des Tap-Pfads mit Gegenprobe auf den Vorgänger-Tag; `ci` und `integration-docker` auf dem Push-Tip grün; `tap-check` grün; danach Tag `v0.18.0` gesetzt und gepusht (Freigabe erteilt 2026-11-02); Kontrolle nach `releasing.md` §6 (M6) erbracht.
3. **Doku nachgezogen** — Roadmap-Zeile unter *Abgeschlossene Wellen* mit Tag-Hash (diese Closure), Bestand-Zeile in `open/README.md`; der [`releasing.md`](../../../maintainer/releasing.md)-Zeiger „zuletzt" ist entfallen (`8cd0976b`, Review F-6) und die DoD-Zeile ist berichtigt; `releasing.md` Punkt 5 trägt die Skelett-Vorlage statt „Muster der früheren Cuts" (Review F-7, `9a7db2f1`).
4. `make gates` grün nach dem letzten Commit vor dem Handoff (Abschnitt Verifikation unten) und mit dieser Closure.
5. **Review und Verification** — Reports unter `docs/reviews/` (oben); Nachlauf-Review über den Diff von Vertragstext (`9a7db2f1`).
6. Closure-Notiz mit Steering-Loop-Lerneintrag; Register fortgeschrieben (unten); jedes Risiko aus §6 mit Ausgang; die drei Paarungen getragen.

**Befunde von Review und Verification, Ausgang je Befund:**

| Befund | Ausgang |
|---|---|
| F-1 MEDIUM (Lead-Satz „`down`, `shell`, `doctor`, `init` and `generate` work as before" widerspricht der Lieferung) | behoben `9a7db2f1`: der Lead nennt jetzt die Fehlerklassen von `up` und `shell` (startet wie `up` ohne `--recreate`), die Label-Suche von `down` und die geänderten Dateien von `init`/`generate`; Verifier hat den Fix gegen [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) gelesen |
| F-2 MEDIUM (Grenzliste ohne „the behaviour of a host `gh` with several accounts") | behoben `9a7db2f1`; Verifier liest die Grenzliste deckungsgleich mit den drei Einträgen |
| F-3 LOW („the value travels by name on the command line" lesbar als Wert auf der Kommandozeile) | behoben `9a7db2f1`: „the handover travels by name on the command line, the value in the environment of the `docker` client process" |
| F-4 LOW (Slice-§7-Marker „liegt in `open/`" abgelaufen) | entfällt mit dieser Closure (die Notiz ersetzt den Marker) |
| F-5 LOW (Roadmap „Letzte Änderung" einen Tag zurück) | behoben `9a7db2f1` |
| F-6 LOW (DoD-Punkt nennt den gestrichenen `releasing.md`-Zeiger) | behoben `9a7db2f1` |
| F-7 LOW (`releasing.md` Punkt 5 „nach dem Muster der früheren Cuts") | behoben `9a7db2f1` (Skelett-Vorlage) |
| F-8 LOW (Roadmap-Zeile „`v0.18.0` — Feature-Katalog (geplant)", Tag vergeben) | behoben `9a7db2f1` (`v0.19.0`, Trigger „`v0.18.0` veröffentlicht") |
| F-9 LOW (Vorschlag `v0.18.0` in zwei CR-Entwürfen und Bestand-Zeilen) | behoben `9a7db2f1` (`v0.19.0`) |
| F-10 LOW (Kopfliste „Planungs-Commits ohne Release-Wirkung" unvollständig) | behoben `9a7db2f1` (Host-git-Mess-Commits, zwei `tools/mk`-Commits) |
| F-11 INFO (Handbuch-Version `2.1`) | aktionslos (`releasing.md` verlangt keine Anhebung) |
| Verification-Rest-Risiken | M7 erbracht (Meldung des Projektinhabers, oben); die zwei vom Verifier benannten mess-getragenen Bedingungen des Container-Name-Slices sind dort zur Closure geführt und nicht Gegenstand dieses Slices |
| Review-Rest (nicht geprüft: CI-/Gate-Evidence, Träger-Produktcode im Ganzen) | vom Verifier nachgeprüft (IDs und Digests, Report oben); die Träger tragen je ein eigenes Review |

**Was ging anders als geplant:** Der Release-Text trug zwei MEDIUM-Fehler (F-1 Widerspruch zur eigenen Lieferung, F-2 Grenzliste kürzer als der Vertrag), gefangen erst vom Review **nach dem Tag** — derselbe Hergang wie im Cut `v0.17.0` (dort F-1/N-1). Die Nachzugs-Commits berichtigen nur `main`, den Tag-Stand nicht ([`releasing.md`](../../../maintainer/releasing.md) *Zweck*). Dazu die Klasse „Zustandsmarker nach Ablauf": vier weitere Funde (F-4 bis F-6, F-8, F-9) trotz der Berichtigungsrunde `827601be`, die vor dem Tag lief.

**Steering-Loop-Eintrag** (gezählt, nicht verkörpert; kein Eintrag trägt das Feld `liegt in`):

- *Lehre 1 (Muster, 3× überschritten innerhalb der verkörperten Regel):* **Zustandsmarker mit Beleg- oder Versionsbezug stehen nach ihrem Ablauf** — im Cut `v0.17.0` zwei Funde (dort F-3/F-4), hier trotz Berichtigungsrunde vier (dieser Review F-4 bis F-6, F-8, F-9). Rückfälle nach der Verkörperung in [`AGENTS.md`](../../../../AGENTS.md) §3.10; Register [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md). Kandidat eines Sensors (zur Entscheidung des Projektinhabers): eine `d-check`-Regel, die „Vorschlag Release `v<X.Y.0>`"-Zusätze in `open/` und Bestand nach einem gesetzten Tag `v<X.Y.0>` rot färbt; ein solcher Sensor ist mechanisch prüfbar (Tag-Liste gegen Textstellen). Es wird nichts verkörpert, ohne die Entscheidung.
- *Lehre 2 (Muster, 2. Vorgang):* **Der Lead des Release-Textes wird nur vom Review gegen die ausgelieferten Einträge gelesen** — kein Gate liest Lead gegen Einträge oder Code [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Plan-Risiko 1); F-1 und F-3 gehören zur Klasse der ausgelieferten Vertragsstellen. Wie beim Vorgänger-Cut zählt sie zur Closure; ein Pflichtschritt „Lead gegen die drei Einträge lesen" in [`releasing.md`](../../../maintainer/releasing.md) §3 bliebe eine Setzung des Projektinhabers.
- *Neuer Sensor:* keiner. *Benannte Spec-Lücke:* keine. *Folge-Slices:* keiner; die Katalog-Entwürfe liegen als Dateien in `open/` und die Roadmap führt `v0.19.0` unter *Nächste Wellen*.

**Beobachtungs-Register (`../observations/`):**

- *Evidence neu (ein Vorgang je Kürzel):* [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (`evidence/slice-v1-release-cut-v0.18.0.md`: F-5, F-6, F-8, F-9, F-10 — Zustandsmarker nach Ablauf und unvollständige Kopfliste, behoben) · [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (`evidence/slice-v1-release-cut-v0.18.0.md`: F-1 und F-3 — Lead-Aussagen über Verhalten des Codes ohne Beleg, behoben).
- *Nicht gezählt:* [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) — F-2 ist die engere Grenzliste (der Lead sagt weniger als das Lastenheft, nicht mehr), wie beim Vorgänger-Cut F-2; zweiter Vorgang der Klasse, kein Register-Dach. [`BEO-DOK/ci-log-falsch-gelesen`](../observations/BEO-DOK/ci-log-falsch-gelesen/observation.md) — Prüfweg über Lauf-IDs und `gh run view` angewandt, kein Fund. [`BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle`](../observations/BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle/observation.md) — `tap-check` lief vor dem Tag, der Trockenlauf mit Gegenprobe davor, kein Fund.

**Carveouts:** keiner neu, keiner gelöst, keiner verändert.

**Weiter offen (nicht abgehakt):** nichts aus der DoD; der Pflichtschritt „Lead gegen die drei Einträge lesen" und der Sensor-Kandidat für die Marker-Klasse liegen als Entscheidungen beim Projektinhaber (Lehre 1 und 2, keine Dateien versprochen).

**Abgleich der Plan-Tabelle** (`git show --stat e6c87379` und die Nachzugs-Commits `827601be`, `9a7db2f1`, `724c5e5e`): der Cut trägt die sieben Dateien der Plan-Tabelle; die Nachzüge sind Review-Fix (Review-Report, CHANGELOG-Lead, Planning-Marker) und Link-Berichtigung, jeder mit Zeile in §3 oder als Review-Artefakt ausgewiesen; „kein Diff"-Zeilen (`spec/`, ADRs, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod`) durch die Lesung des Reviewers belegt.

**Drei Paarungen** (geprüft auf den Inhalt dieser Notiz): *Anker* — kein Steering-Loop-Eintrag trägt ein Zielort-Feld; nichts verkörpert, nichts zu paaren. *Folge-Slice* — kein neuer Slice versprochen; die zwei Katalog-Entwürfe und der Basis-Image-Entwurf liegen als Dateien in `open/` (`ls` bei der Closure). *Register* — jede genannte Kennung `BEO-<KUERZEL>/<slug>` ist ein Verzeichnis mit der Evidence-Datei dieser Closure (`ls` bei der Closure).

**Verifikation** (Läufe nach Closure-Inhalt, Move, Verweise-Nachzug und den zwei d-check-Berichtigungen am Verification-Report; Docker-Modus): `make gates` Exit 0 am 2026-11-02 (`record-gates`-Stempel = HEAD `d4008b00`). Der Lauf gilt für den Stand vor dieser Zeile; die Zeile ändert nur Prosa.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `CHANGELOG.md`, READMEs, Handbuch und Planung (alle DOK laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)); die drei Versionszeilen liegen in `cmd/` und dem Build-Werkzeug — `cmd/` führt keine Sub-Area (Modus-Deklaration, Zeile Hexagon-Schichtungs-Linie), die Builds der Tool-Dateien sind keine Sub-Area-Änderung, weil keine Konvention berührt wird.

**Vorgelagert — offene Beobachtungen sichten** (Register 2026-11-01 durchgesehen, nur zitierte Einträge gelesen): [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Risiko 1), [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (CI-Nachweis statt test-docker lokal). Evidenz-/Diskrepanz-Risiko: gering (der Cut trägt geschlossene Zusagen; der Reviewer liest den Lead gegen die Lastenheft-Stände). Alle berührten Sub-Areas Greenfield — kein Modus-Begründungsblock nötig.