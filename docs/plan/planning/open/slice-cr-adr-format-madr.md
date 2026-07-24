# Slice CR: ADR-Format auf MADR-Template ([`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)-Change-Request)

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Phase:** CR (Change Request am Vertrags-Stratum). Kein Produkt-Meilenstein.

**Bezug:** Dieser Slice **ist** der Change Request, mit dem
[`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)
(Vertrags-Stratum `spec/lastenheft.md`, „vertraglich abnahmebindend") auf die
vendored MADR-ADR-Template-Form gehoben wird. Trigger ist die v3.5.1-Adoption
([`slice-harness-regelwerk-adoption-v3.5.1`](slice-harness-regelwerk-adoption-v3.5.1.md),
dortiger `MR-008` verweist auf diesen CR, ändert aber `LH-FA-PROJDOCS-002`
nicht). Kohärent mit dem Referenzmodell
([`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell),
[`ADR-0013`](../../adr/0013-dokumentationsreferenzmodell.md)): das
Template-Feld `Schärft:` ist die Aufwärts-Kopplung Spec↔ADR.

**Autor:** pt9912. **Datum:** 2026-07-24.

---

## 1. Warum ein CR (und kein conventions-MR)

Das vendored ADR-Template
(`.harness/baseline/v3.5.1/templates/docs/plan/adr/NNNN-titel.template.md`,
MADR-/Nygard-Stil) **kollidiert** mit dem heutigen
[`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)
— es ist nicht additiv:

| Punkt | `LH-FA-PROJDOCS-002` (heute) | MADR-Template |
|---|---|---|
| Status / Datum | `##`-Überschriften | fette Inline-Felder (`**Status:**`) |
| Kopf-Felder | keine | `**Autor:**`, `**Bezug:**`, `**Schärft:**` |
| Pflicht-Sections | 5 (Status·Datum·Kontext·Entscheidung·Konsequenzen) | Kontext·Entscheidung·**Verglichene Alternativen**·Konsequenzen·**Fitness Function**·**Re-Evaluierungs-Trigger**·**Geschichte** |
| Titel | `# ADR <Nr>: <Titel>` | `# ADR-NNNN: <Titel>` |
| Superseded-Ref | `Superseded by <NNNN>-<slug>` | `Superseded by ADR-NNNN` |

Weil das Ziel-Stratum `spec/lastenheft.md` ist, dessen Änderungsmechanismus das
Regelwerk als **Change Request** führt (`grundlagen-konventionen` §Straten:
„Vertrag → Change Request"), ist ein `harness/conventions.md`-MR **nicht**
zulässig (der Adaptions-Block ist form-bringend, nicht vertrags-ändernd). Die
Anforderung selbst wird geändert — kontrolliert, mit Trigger, Begründung und
Review. Dieser Slice ist das CR-Vehikel (u-boot führt kein separates
CR-Artefakt).

## 2. Definition of Done

<!-- Prüfbare Kriterien für den Übergang nach done/. -->

- [ ] **CR-Charakter explizit dokumentiert** (dieser Slice §1): Trigger
  (v3.5.1-Adoption), betroffenes Vertrags-Stratum, Reversibilität,
  Grandfathering-Entscheidung. Kein stiller Spec-Edit.
- [ ] **[`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)
  umgeschrieben** auf die MADR-Template-Struktur als neue Pflicht-Form:
  Inline-Kopf-Felder (`Status`, `Datum`, `Autor`, `Bezug`, `Schärft`),
  Sections `Kontext` · `Entscheidung` · `Verglichene Alternativen` ·
  `Konsequenzen` · `Fitness Function` (falls maschinell prüfbar) ·
  `Re-Evaluierungs-Trigger` · `Geschichte`.
- [ ] **Reconciliation der Format-Divergenzen** explizit entschieden und in
  `LH-FA-PROJDOCS-002` festgehalten (nicht offen lassen):
  - **Titel:** `# ADR-NNNN: <Titel>` (Template) **oder** u-boots bisheriges
    `# ADR <Nr>: <Titel>` — eine Form wählen, begründen.
  - **Superseded-Ref:** klickbarer Datei-Link statt roher Tokens; Abgleich mit
    der d-check-`ids`-Regel (`\bADR-\d{4}\b`, Ziel `docs/plan/adr/`) und dem
    Referenzmodell — u-boots Link-Form gewinnt (Anker-/Datei-Auflösung).
  - **Status-Werte:** deutscher/englischer Satz vereinheitlicht (Template:
    `Proposed | Accepted | Deprecated | Superseded by …`).
- [ ] **Grandfathering (status-basiert, nicht Nummern-Präfix):** die zum
  CR-Zeitpunkt **Accepted** ADRs (`0001`–`0010`, `0013` = 11 Stück) bleiben in
  der leanen 5-Abschnitt-Form und **immutabel** (Hard Rule „Accepted =
  nicht inhaltlich überschreiben"); sie werden **nicht** auf MADR migriert.
  `LH-FA-PROJDOCS-002` trägt die Grandfather-Klausel explizit.
- [ ] **Proposed-ADRs (`0011`, `0012`) + alle künftigen ADRs** MADR-konform:
  0011/0012 sind noch mutabel → beim nächsten Anfassen (ihre offenen
  Ratifizierungs-Slices) auf MADR heben; neu angelegte ADRs kopieren die
  vendored Vorlage. (Kein Zwang, 0011/0012 *jetzt* zu migrieren, solange sie
  Proposed bleiben — als Klausel festhalten.)
- [ ] **[`docs/plan/adr/README.md`](../../adr/README.md) nachgezogen:**
  `## Konventionen`-Abschnitt auf die MADR-Form + die `Schärft`-Aufwärts-
  Deklaration; Grandfather-Hinweis für den Accepted-Bestand.
- [ ] **`MR-008` in [`harness/conventions.md`](../../../../harness/conventions.md)
  fortgeschrieben:** von „pending CR" auf „CR ausgeführt" (Bestand grandfathered,
  neue ADRs MADR); Verweis auf diesen Slice.
- [ ] **`make gates` grün** (`links`/`anchors`/`ids`/`matrix` — die neuen
  `Bezug`/`Schärft`-Links in künftigen ADRs sind linkpflichtig und müssen
  auflösen; die Grandfather-ADRs bleiben unverändert und damit gate-neutral).
- [ ] **Review + Verification-Evidence**; Closure-Notiz mit Commit-Hash.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` (`LH-FA-PROJDOCS-002`) | update (CR) | Pflicht-ADR-Form auf MADR; Grandfather-Klausel; Format-Reconciliation |
| `docs/plan/adr/README.md` | update | `## Konventionen` auf MADR + `Schärft`-Deklaration + Grandfather-Hinweis |
| `harness/conventions.md` (`MR-008`) | update | „pending CR" → „ausgeführt"; Grandfather-Record |
| `docs/plan/adr/0011`, `0012` | prüfen | Proposed → optional MADR beim nächsten Anfassen (nicht erzwungen) |
| Accepted-ADRs `0001`–`0010`, `0013` | **unverändert** | immutabel, grandfathered — kein Retrofit |

## 4. Trigger

Manuell, nach Review dieses Plans **und** nach Schließung des
[`slice-harness-regelwerk-adoption-v3.5.1`](slice-harness-regelwerk-adoption-v3.5.1.md)
(die vendored Baseline + `MR-008`-Verweis müssen stehen).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Review nach
[`harness/review.md`](../../../../harness/review.md) + Verification-Evidence nach
[`harness/verification.md`](../../../../harness/verification.md) + Closure-Notiz.

## 6. Risiken und offene Punkte

- **Vertrags-Änderung ist abnahmebindend:** `LH-FA-PROJDOCS-002` hat Priorität
  MVP. Die Umschreibung muss RTM-/`ids`-kompatibel bleiben (keine gebrochenen
  `LH-*`-Anker durch Heading-Umbau der Anforderung selbst).
- **Immutabilitäts-Kollision vermeiden:** Kein Wort an den 11 Accepted-ADRs
  ändern (auch keine additiven `Schärft`-Felder) — das würde die Immutabilität
  verletzen (vgl. belief-agent-Erfahrung: schon ein additiver Link erzeugt
  Core-Drift). Deshalb reine Grandfather-Klausel, kein Bestands-Retrofit.
- **Link-Dichte:** MADR-`Bezug`/`Schärft` machen jede ADR link-reicher; unter
  der `ids`-Linkpflicht müssen diese Anker real auflösen. Betrifft nur neue ADRs.
- **Titel-/Superseded-Form:** Divergenz Template ↔ Bestand bewusst zugunsten der
  u-boot-Link-Konventionen auflösen (nicht blind Template übernehmen).
- **Kein Carveout erwartet:** additive Spec-Präzisierung; keine Gate-Lockerung.

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen: Commit-Hash, Steering-Loop-Lerneintrag. -->

## 8. Sub-Area-Modus-Begründung

Berührt schreibend nur `spec/` und `docs/plan/` (Doku-/Spec-führt-Sub-Areas).

- **Modus:** GF (spec-führt; die Anforderung wird präzisiert, nicht aus
  Bestandscode rekonstruiert).
- **Evidenz-/Diskrepanz-Risiko:** niedrig — reine Doku/Spec, Absicherung durch
  `docs-check` (Anker/Links/Referenzmodell).
- **Graduation-Bedingung:** n/a (GF).
