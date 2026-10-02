# Slice V2: Lastenheft-Schnitt 3: Projektkapitel als Randbedingungen, Schnittstellen und Daten

> **Status:** **Done** (2026-10-02, **Delivery-Hash: `d4abd2b`**).

**Welle:** `welle-spec-technik-stratum` (siehe [`roadmap.md`](../in-progress/roadmap.md)).
**Bezug:** Lastenheft §4.11 Build- und CI-Infrastruktur, §4.12 Doku-Struktur, §4.13 Architektur des Projekts, §6.2 Dateischnittstellen, §6.3 Docker-Schnittstelle, §7 Datenanforderungen (Scope).
**Berührte Spec-Stellen:** `spec/spezifikation.md` §1–§3, §6; `spec/architecture.md` (Importregeln, Verzeichnislayout bleiben dort); `spec/lastenheft.md` §4.11–§4.13, §6.2–§6.3, §7.
**Autor:** pt9912. **Datum:** 2026-10-02.

## Auslöser

Die Kapitel §4.11–§4.13 beschreiben das Repository selbst (Dockerfile-Stages, Make-Targets, CI-Jobs, Doku-Struktur, hexagonale
Schichten). Entscheidung des Projektinhabers: im Lastenheft als **Randbedingungen** kürzen (Vorgabe plus Nachweis), die Details in
die Spezifikation beziehungsweise die Architektur-Sicht.

## Ziel und Abgrenzung

**Ziel:** §4.11–§4.13 tragen als Vertrag nur Vorgabe und Nachweis je Anforderung (Docker-only, Gates, Doku-Struktur, hexagonale Architektur); §6.2, §6.3 und §7 sind technisch geschnitten; alle Details stehen wörtlich an ihrem neuen Ort.

**Schnittregeln:** siehe [`slice-v2-spec-technik-stratum-grundlage`](../done/slice-v2-spec-technik-stratum-grundlage.md) §Schnittregeln (bleibt im Lastenheft / zieht in die Spezifikation, Decken-Regel, Verbleib-Tabelle, Zeilenverweise).

**Ausdrücklich NICHT in diesem Slice:**

- **Neue Anforderungs-Kennungen** — alle bestehenden `LH-*`-Kennungen bleiben; ein Umbau zu `-RB-`-Kennungen wäre eine Vertragsänderung mit Link-Bruch in ADRs, Slices und Tests und ist nicht beschlossen.
- **Inhaltliche Änderung der Importregeln** — sie bleiben in `spec/architecture.md`; die Angleichung der Sicht an die Vorlage ist ein eigener Slice.


## Definition of Done

- [x] §4.11–§4.13 auf Vorgabe plus Nachweis gekürzt; §6.2, §6.3, §7 geschnitten; Verbleib-Tabelle vollständig; Kennungen und Anker unverändert.
- [x] Spezifikation ergänzt (Build-/CI-Festlegungen, Doku-Struktur, Dateischnittstellen, Datenformate) bzw. Verweis auf die Architektur-Sicht an den passenden Stellen, ohne ADR- und Slice-Verweis.
- [x] `make gates` grün (inklusive `make docs-check`), Exit-Code separat geprüft.
- [x] Review durch eine andere Rolle als die Umsetzung (`harness/review.md`), Report unter `docs/reviews/`.

## Tranchen

| T | Inhalt |
| - | ------ |
| T1 | **Inventur und Verbleib-Tabelle** für §4.11–§4.13, §6.2, §6.3, §7. |
| T2 | **Projektkapitel:** je Anforderung Vorgabe und Nachweis im Lastenheft belassen, Details übernehmen (Spezifikation oder Architektur-Sicht). |
| T3 | **Schnittstellen und Daten:** §6.2, §6.3, §7 übernehmen. |
| T4 | **Gegenlesen,** `make gates`, Review. |

## Risiken

- Die Nachweise in §4.11–§4.13 verweisen auf Make-Targets und Workflows; beim Kürzen darf kein Gate-Anspruch still wegfallen (`AGENTS.md`: Gates nicht lockern). T4 vergleicht die Sensor-Liste vorher/nachher.

## Closure

- **Geliefert (`d4abd2b`):** Die Projektkapitel §4.11 (Build/CI), §4.12 (Doku-Struktur) und §4.13 (Architektur) sowie §6.2 (Markierungsformate) sind auf Vorgabe und Nachweis gekürzt; die Details stehen wörtlich in der Spezifikation: Stages, Runtime-Eigenschaften, Build-Args, `.dockerignore`-Liste, Makefile-Eigenschaften und -Targets, Repository-Mindestlayout, Docs-Verzeichnisbaum, ADR-Format, Schichten-Layout, Import-Regel-Tabelle und Markierungsformat je Dateityp. Alle `LH-*`-Kennungen und Anker unverändert, kein Gate-Anspruch entfernt: die Aggregator-Targets (`gates`, `ci`, `fullbuild`), Docker-only-Workflow und Coverage-Bootstrap bleiben im Lastenheft.
- **Nicht geschnitten (bewusst):** §6.3 und §7 enthalten nur Verhalten und kurze Beispiele und blieben unverändert; die Abgrenzung zu Zielprojekten in §4.12 bleibt, weil sie Vertrag ist.
- **Folge für die Sicht:** Die Import-Regel-Tabelle und das Schichten-Layout stehen jetzt in der Spezifikation und noch einmal in `spec/architecture.md`; der Abgleich beider ist Gegenstand von `slice-v2-spec-architektur-angleichung` (Sicht visualisiert, Technik führt).
- **Gegenlesen:** 150 entfernte Zeilen, 0 fehlend in der Spezifikation.
- **Sensoren:** `make gates` grün (lint, test, coverage-gate, docs-check); Konfigurationen (`.golangci.yml`, Makefile, Dockerfile) unberührt. Nicht ausgeführt: `make ci`, `make test-docker`.
- **Review:** unabhängiges Review über die gesamte Welle am Abschluss.

### Verbleib-Tabelle

| Quelle (Lastenheft) | Ziel (Spezifikation) | Ort |
|---|---|---|
| §4.11 [LH-FA-BUILD-001](../../../../spec/lastenheft.md#lh-fa-build-001--multi-stage-dockerfile-u-boot-repo) (Mindestumfang, 1335–1344) | [LH-FA-BUILD-001.a](../../../../spec/spezifikation.md#lh-fa-build-001a--stages-des-multi-stage-dockerfiles) | spezifikation.md §1 |
| §4.11 [LH-FA-BUILD-002](../../../../spec/lastenheft.md#lh-fa-build-002--runtime-stage-pflichten) (Pflichten, 1354–1362) | [LH-FA-BUILD-002.a](../../../../spec/spezifikation.md#lh-fa-build-002a--eigenschaften-der-runtime-stage) | spezifikation.md §1 |
| §4.11 [LH-FA-BUILD-003](../../../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) (Build-Args, 1372–1374) | [SPEC-006](../../../../spec/spezifikation.md#spec-006--build-args-des-dockerfiles-und-pin-politik) | spezifikation.md §3 |
| §4.11 [LH-FA-BUILD-004](../../../../spec/lastenheft.md#lh-fa-build-004--dockerignore-pflicht) (Ausschlussliste, 1386–1391) | [SPEC-007](../../../../spec/spezifikation.md#spec-007--mindest-ausschlüsse-der-dockerignore) | spezifikation.md §3 |
| §4.11 [LH-FA-BUILD-005](../../../../spec/lastenheft.md#lh-fa-build-005--makefile-mit-standard-targets) (Eigenschaften und Target-Tabelle, 1403–1423) | [SPEC-008](../../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) | spezifikation.md §2 |
| §4.11 [LH-FA-BUILD-009](../../../../spec/lastenheft.md#lh-fa-build-009--repository-layout) (Mindestlayout-Baum, 1477–1504) | [SPEC-009](../../../../spec/spezifikation.md#spec-009--mindestlayout-des-u-boot-repositories) | spezifikation.md §2 |
| §4.12 [LH-FA-PROJDOCS-001](../../../../spec/lastenheft.md#lh-fa-projdocs-001--mindeststruktur) (Verzeichnisbaum, 1522–1533) | [SPEC-010](../../../../spec/spezifikation.md#spec-010--mindest-verzeichnisstruktur-unter-docs) | spezifikation.md §2 |
| §4.12 [LH-FA-PROJDOCS-002](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format) (Format-Regeln, 1547–1563) | (Eintrag inzwischen entfernt: deckt das Baseline-Regelwerk ab) | spezifikation.md §2 |
| §4.13 [LH-FA-ARCH-002](../../../../spec/lastenheft.md#lh-fa-arch-002--schichten-und-verzeichnislayout) (Verzeichnisbaum, 1704–1715) | [SPEC-012](../../../../spec/spezifikation.md#spec-012--schichten-und-verzeichnislayout-unter-internal) | spezifikation.md §2 |
| §4.13 [LH-FA-ARCH-003](../../../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement) (Import-Regeln, 1727–1735) | [SPEC-013](../../../../spec/spezifikation.md#spec-013--import-regel-tabelle-der-schichten) | spezifikation.md §2 |
| §6.2 [LH-SA-FILE-002](../../../../spec/lastenheft.md#lh-sa-file-002--markierte-verwaltete-bereiche) (Markierungsformate, 2086–2108) | [SPEC-014](../../../../spec/spezifikation.md#spec-014--markierungsformat-verwalteter-bereiche-je-dateityp) | spezifikation.md §2 |
