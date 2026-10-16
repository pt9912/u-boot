# Slice CR: Lastenheft nur mit Produktanforderungen — Prozess-Anforderungen als Grabsteine

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — keine Closure-Bedingung, die von der DoD dieses Slice
verschieden ist (Change Request am Vertrags-Stratum).

**Bezug:** Dieser Slice **ist** der Change Request. Prinzip (Projektinhaber,
2026-10-16): Das Lastenheft beschreibt die Anforderungen **des Tools**, nicht die
des Entwicklungsprozesses; der steht im Regelwerk. Betroffen sind
[`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
und die Familien `LH-FA-BUILD-*`, `LH-FA-PROJDOCS-*`, `LH-FA-ARCH-*`, `LH-QA-*`;
[`LH-LESE-003`](../../../../spec/lastenheft.md#lh-lese-003--dokumentenordnung)
(Dokumentenordnung) ist die Grundlage. Vorgänger-CR:
[`slice-cr-adr-format-madr`](../done/slice-cr-adr-format-madr.md).
Löst die Blockade von
[`slice-harness-wellen-betrieb`](../open/slice-harness-wellen-betrieb.md) (kein
Vertragswortlaut für Welle-Dateien mehr nötig).

**Berührte Spec-Stellen:** Lastenheft (Vertrags-Stratum) §4.11–§4.13, §8, §12, §16;
Spezifikation (Verfeinerungen der entfallenden Anforderungen); Architektur-Sicht
(Verweise). Änderung des Lastenhefts nur nach Vereinbarung mit dem Projektinhaber
(Lastenheft §16) — die Einzelentscheidungen in §2 stehen zur Freigabe.

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Lastenheft enthält nur noch Zusagen, die das Produkt `u-boot` dem
Nutzer gibt. Die 22 Anforderungen, die den Bau, die Doku und die Qualitätssicherung
der u-boot-Codebase regeln, entfallen als **Grabstein** (Überschrift und ID bleiben,
ein Satz nennt Version und neuen Ort); die ID wird nie wieder vergeben.

**Warum Grabsteine:** Die IDs stehen in Accepted ADRs (unveränderlich), in
erledigten Slices, Reviews, CI, Skripten und Code-Kommentaren (zusammen über 600
Verweise in mehr als 200 Dateien). Ein entfernter Anker bräche sie; ein Grabstein
hält jeden dieser Links gültig.

| Familie | Anforderungen | Neuer Ort der Regelung |
|---|---|---|
| [`LH-FA-BUILD-001`](../../../../spec/lastenheft.md#lh-fa-build-001--multi-stage-dockerfile-u-boot-repo)..`009` (§4.11) | Dockerfile, Runtime-Stage, Pins, `.dockerignore`, Makefile, Aggregatoren, Docker-only, Coverage-Bootstrap, Repo-Layout | Regelwerk (Docker-Harness, Quality Gates); Festlegungen bleiben als `SPEC-<NNN>`; Sensor-Dateien und `harness/README.md` |
| [`LH-FA-PROJDOCS-001`](../../../../spec/lastenheft.md#lh-fa-projdocs-001--mindeststruktur)..`006` (§4.12) | Mindeststruktur, ADR-Format, Planning-Lifecycle, Carveouts, Referenzmodell, Archivierung | Regelwerk (Spec, ADRs, Planung, Carveouts); u-boot-Abweichungen als `MR-<NNN>` |
| [`LH-FA-ARCH-001`](../../../../spec/lastenheft.md#lh-fa-arch-001--hexagonales-pattern)..`003` (§4.13) | hexagonales Pattern, Schichten, Import-Regeln | Architektur-Sicht `spec/architecture.md`, [`SPEC-012`](../../../../spec/spezifikation.md#spec-012--schichten-und-verzeichnislayout-unter-internal)/[`SPEC-013`](../../../../spec/spezifikation.md#spec-013--import-regel-tabelle-der-schichten), [`ADR-0002`](../../adr/0002-hexagonale-architektur.md) |
| [`LH-QA-001`](../../../../spec/lastenheft.md#lh-qa-001--automatisierte-tests)..`004` (§8) | Tests, Akzeptanzkriterien, CI, Lint | Regelwerk (Quality Gates, Verification), [`ADR-0003`](../../adr/0003-solid-nahes-lint-profil.md), [`ADR-0004`](../../adr/0004-ci-system.md) |

**Prüfpunkte — Produktzusagen in den Prozess-Blöcken** (je Einzelentscheidung
zur Freigabe, Vorschlag in Klammern):

1. [`LH-FA-BUILD-002`](../../../../spec/lastenheft.md#lh-fa-build-002--runtime-stage-pflichten) (Runtime-Image: minimal, shell-los, Non-root, OCI-Labels): das
   ausgelieferte Image ist ein Produktmerkmal (Distribution, siehe [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung)).
   (Vorschlag: als Produktzusage in den Distributions-Abschnitt übernehmen, die
   Stage-Mechanik entfällt.)
2. [`LH-QA-002`](../../../../spec/lastenheft.md#lh-qa-002--testbare-akzeptanzkriterien) (jede funktionale Anforderung ist durch einen Akzeptanztest
   überprüfbar): Eigenschaft des Lastenhefts selbst. (Vorschlag: als Satz in den
   Lesehinweisen behalten.)
3. [`LH-QA-001`](../../../../spec/lastenheft.md#lh-qa-001--automatisierte-tests), [`LH-QA-003`](../../../../spec/lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions), [`LH-QA-004`](../../../../spec/lastenheft.md#lh-qa-004--linting-solid-nahes-lint-profil): reine Prozess-Anforderungen.
   (Vorschlag: entfallen.)

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Änderung produktbezogener Anforderungen (`LH-FA-CLI-*`, `-INIT-*`, `-DEV-*` usw.):
  der Slice verschiebt nur die Grenze zwischen Produkt und Prozess.
- Umschreiben von Accepted ADRs, erledigten Slices, Reviews und Code-Kommentaren:
  ihre Verweise bleiben über die Grabsteine gültig; ADRs sind unveränderlich.
- Der Wellen-Betrieb selbst: [`slice-harness-wellen-betrieb`](../open/slice-harness-wellen-betrieb.md).

## 2. Definition of Done

- [x] **Freigabe** der Einzelentscheidungen (§1 Prüfpunkte, Grabstein-Form, Version
  `0.6.0`) durch den Projektinhaber liegt vor (2026-10-16).
- [x] **Lastenheft:** 22 Anforderungen als Grabsteine (Überschrift und Anker
  unverändert, darunter ein Satz `Entfallen (0.6.0)` mit neuem Ort, ohne Abwärtsverweis
  auf Planung/ADR), Einleitungen der Abschnitte §4.11–§4.13 angepasst, §12 MVP-Umfang
  bereinigt, Version `0.6.0`, Zeile in §16 Historie; neu
  [`LH-NFA-PORT-004`](../../../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image)
  und [`LH-LESE-004`](../../../../spec/lastenheft.md#lh-lese-004--prüfbarkeit).
- [x] **Spezifikation:** die 13 Verfeinerungen `.a`/`.b` der entfallenden Anforderungen
  tragen neue Kennungen (Anhang), die alten Anker bleiben als HTML-Anker erhalten; Eltern-Verweise
  zeigen auf die neuen Stellen.
- [x] **Lebende Verweise nachgezogen:** `spec/architecture.md`, `AGENTS.md`, `harness/`
  (Sensor-Dateien, `conventions`, MR-Einträge), Skills, `docs/maintainer/`, Planning-READMEs,
  Carveout-Inventar, READMEs zeigen auf den neuen Ort. Historische Stellen bleiben über die
  Grabsteine gültig (Accepted ADRs, `done/`, Reviews, `MR-008`, Code- und Konfig-Kommentare).
- [x] `trace` und `doc-complete`: 83 Anforderungen, 0 Waisen — die Grabsteine erzeugen
  keine Requirement-Waisen.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      ([`Report`](../../../reviews/2026-10-16-slice-cr-lastenheft-produktanforderungen.md); 0 HIGH, 0 MEDIUM, 6 LOW, 2 INFO; F-1 bis F-3 und F-5 behoben; `.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Grabsteine, Abschnittseinleitungen, §9, §12, Version, §16 (nach Freigabe) |
| `spec/spezifikation.md` | update | `.a`/`.b`-Verfeinerungen → `SPEC-<NNN>`; Eltern-Verweise |
| `spec/architecture.md` | update | 23 Verweise auf `LH-FA-ARCH-*` |
| `harness/sensors/*.md`, `harness/README.md`, `harness/conventions.md`, `harness/conventions/MR-*.md` | update | Bindungen auf Anforderungen, die entfallen |
| `AGENTS.md` | update | Zeiger-Tabelle (§3.8) nennt `LH-FA-PROJDOCS-*` |
| `docs/user/`, `docs/maintainer/`, `README*.md` | update | lebende Verweise |
| `.d-check.yml` | prüfen | `trace`-Familien, `ids`-Muster |

Ansatz: erst die Zuordnungstabelle (alt → neu) und die Grabstein-Texte als Entwurf
im Anhang dieses Slice zur Freigabe, dann Lastenheft, dann Spezifikation, dann die
lebenden Verweise, je ein Commit.

## 4. Trigger

**Start** (`next` → `in-progress`): Priorisierung durch den Projektinhaber
(erfolgt 2026-10-16: Entscheidung „alle vier Familien“).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Familien einzeln als eigene CR-Slices führen
  (Reihenfolge `PROJDOCS`, `BUILD`, `QA`, `ARCH`).
- `in-progress` → `open` (blockiert): der Projektinhaber gibt eine Einzelentscheidung
  (§1 Prüfpunkte) nicht frei.

## 5. Closure-Trigger

DoD vollständig, Review- und Verification-Evidence liegen vor, Closure-Notiz
geschrieben, `make gates` grün.

## 6. Risiken und offene Punkte

- Verweise in lebenden Dokumenten bleiben auf eine entfallene Anforderung stehen und
  wirken wie eine geltende Zusage — **Ausgang:** offen; Suche nach den 22 IDs
  außerhalb von `done/`, `docs/reviews/` und Accepted ADRs vor der Closure.
- Eine Produktzusage in den Prozess-Blöcken geht verloren (Prüfpunkte §1) —
  **Ausgang:** offen bis zur Freigabe.
- `ARCH`/`QA` tragen die Bindungen der Gates (`depguard`, Coverage-Schwelle); die
  Sensor-Dateien und ADRs verweisen auf sie — **Ausgang:** offen; neue Bindung ist
  die ADR oder `SPEC-<NNN>`.

## 7. Closure-Notiz

Offen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/`, `harness/` und
`docs/plan/`
([`conventions.md`](../../../../harness/conventions.md) §Modus-Deklaration);
keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** das Register ist leer.

**Modus:** alle berührten Sub-Areas GF.

## Anhang: Zuordnung der Verfeinerungen (alt → neu)

| Alt | Neu |
|---|---|
| [`LH-FA-BUILD-001`](../../../../spec/lastenheft.md#lh-fa-build-001--multi-stage-dockerfile-u-boot-repo)`.a` | [`SPEC-025`](../../../../spec/spezifikation.md#spec-025--stages-des-multi-stage-dockerfiles) |
| [`LH-QA-003`](../../../../spec/lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions)`.a` | [`SPEC-026`](../../../../spec/spezifikation.md#spec-026--komposition-der-ci-pipeline-github-actions) |
| [`LH-QA-004`](../../../../spec/lastenheft.md#lh-qa-004--linting-solid-nahes-lint-profil)`.a` | [`SPEC-027`](../../../../spec/spezifikation.md#spec-027--zusammensetzung-des-lint-profils) |
| [`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)`.a` | [`SPEC-028`](../../../../spec/spezifikation.md#spec-028--bestandsregel-und-abgleich-zum-adr-format) |
| [`LH-FA-BUILD-001`](../../../../spec/lastenheft.md#lh-fa-build-001--multi-stage-dockerfile-u-boot-repo)`.b` | [`SPEC-029`](../../../../spec/spezifikation.md#spec-029--mindestumfang-des-multi-stage-dockerfiles) |
| [`LH-FA-BUILD-003`](../../../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik)`.a` | [`SPEC-030`](../../../../spec/spezifikation.md#spec-030--overrides-der-build-args) |
| [`LH-FA-BUILD-007`](../../../../spec/lastenheft.md#lh-fa-build-007--docker-only-workflow)`.a` | [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) |
| [`LH-FA-BUILD-008`](../../../../spec/lastenheft.md#lh-fa-build-008--coverage-bootstrap)`.a` | [`SPEC-032`](../../../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) |
| [`LH-FA-BUILD-009`](../../../../spec/lastenheft.md#lh-fa-build-009--repository-layout)`.a` | [`SPEC-033`](../../../../spec/spezifikation.md#spec-033--go-layout-des-repositories) |
| [`LH-FA-ARCH-003`](../../../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement)`.a` | [`SPEC-034`](../../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln) |
| [`LH-QA-004`](../../../../spec/lastenheft.md#lh-qa-004--linting-solid-nahes-lint-profil)`.b` | [`SPEC-035`](../../../../spec/spezifikation.md#spec-035--pflichten-des-lint-profils) |
| [`LH-FA-BUILD-002`](../../../../spec/lastenheft.md#lh-fa-build-002--runtime-stage-pflichten)`.a` / `.b` | [`LH-NFA-PORT-004`](../../../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image)`.a` / `.b` (Verfeinerungen der neuen Produktzusage) |

Eltern-Verweise auf entfallene Anforderungen: Bau (`BUILD-001`..`009`) → [`SPEC-025`](../../../../spec/spezifikation.md#spec-025--stages-des-multi-stage-dockerfiles)..[`SPEC-033`](../../../../spec/spezifikation.md#spec-033--go-layout-des-repositories),
[`SPEC-006`](../../../../spec/spezifikation.md#spec-006--build-args-des-dockerfiles-und-pin-politik)..[`SPEC-008`](../../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles); Architektur → [`SPEC-012`](../../../../spec/spezifikation.md#spec-012--schichten-und-verzeichnislayout-unter-internal), [`SPEC-034`](../../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln); [`LH-QA-002`](../../../../spec/lastenheft.md#lh-qa-002--testbare-akzeptanzkriterien) → [`LH-LESE-004`](../../../../spec/lastenheft.md#lh-lese-004--prüfbarkeit);
Planung, Carveouts, Referenzmodell → Baseline-Regelwerk.
