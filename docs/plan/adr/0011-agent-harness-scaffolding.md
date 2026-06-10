# ADR 0011: u-boot scaffoldet Agent-Harness-Artefakte (GF/BF)

## Status

Proposed

> **Entwurf — noch nicht ratifiziert.** Festgehalten, damit die Idee
> nicht verloren geht. Die §Entscheidung unten ist ein **Vorschlag**;
> die §Offenen Fragen müssen vor `Accepted` beantwortet werden. Kein
> Code, bis ratifiziert + Spec-Erweiterung + Planning-Artefakt stehen.

## Datum

2026-06-09

## Kontext

`u-boot` ist laut [`LH-ZB-002`](../../../spec/lastenheft.md#lh-zb-002-produktvision) ein „Bootloader für
Entwicklungsumgebungen". Bisher scaffoldet es die **Laufzeit-/Projekt-
Schicht** (Docker, Compose, Devcontainer, Services, `u-boot.yaml`,
README/CHANGELOG, optional lokale/Katalog-Templates via
[ADR-0009](0009-template-format-yaml-files.md)).

Eine Schicht *darüber* ist die **Agent-Harness-Schicht**: die
Dokument-Artefakte, mit denen ein Repo reproduzierbar entlang einer
Spec von KI-Agenten bearbeitet wird — `AGENTS.md`, `harness/README.md`,
`harness/conventions.md`, `spec/lastenheft.md`, `spec/architecture.md`,
`docs/plan/` (Roadmap, ADRs, Planning-Lifecycle), Quality-Gate-`Makefile`.

Der externe Kurs **`ai-harness-course`** (Modul 2 — Harness-Bootstrap;
nicht Teil dieses Repos, daher hier nur als Klartext referenziert)
beschreibt diesen Einstiegsprozess als **GF/BF-Modus pro Sub-Area**:

- **Greenfield (GF):** leeres Repo, Doku führt → Code. Skelette werden
  zuerst angelegt (Templates → `conventions.md` → Lastenheft-Outline →
  Roadmap → Sensors → Architektur → erste ADR), dann Code.
- **Brownfield (BF):** bestehender Code ohne Harness, Code führt → Doku.
  Code-Inventur (Discovery) → Reverse-Engineering von
  Spec/Architektur/retroaktiven ADRs → Diskrepanz-Klassifikation
  (orphan code → Carveout, orphan requirement → Reconciliation-Artefakt,
  implicit decision → retro-ADR) → Roadmap als Reconciliation-Plan mit
  Graduation BF→GF.

**Warum das zu u-boot passt:**

- u-boot *ist* ein Bootstrapper — Harness-Bootstrap ist dieselbe
  Verb-Klasse, eine Schicht höher.
- **GF/BF ist in u-boot bereits erstklassig:** [`LH-FA-INIT-004`](../../../spec/lastenheft.md#lh-fa-init-004-bestehendes-projekt-erkennen)
  (Bestehendes-Projekt-Erkennung) + `--assume-existing`
  ([`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a-interaktivität-und-automatisierung)) ist exakt die GF-(frisch)- vs.
  BF-(bestehend)-Unterscheidung. Der Discovery-Trigger existiert.
- Die **Template-Engine existiert** ([ADR-0009](0009-template-format-yaml-files.md)): ein Harness-Skelett ist
  ein weiteres `text/template`-Template-Set.
- u-boot **dogfoodet das Layout selbst** (`AGENTS.md`, `harness/`,
  `spec/`, `docs/plan/`) → glaubwürdige Referenz-Ausgabe.

**Die zentrale Scope-Grenze:** Der BF-Kern (Code-Inventur +
Reconciliation-Backlog aus beliebigem Code ableiten) ist
**inferenz-lastige Agenten-Arbeit, kein deterministischer
Template-Fill**. u-boot kann nur das **Skelett** liefern; die
Diskrepanz-Analyse selbst nicht. Das ist die Linie, an der jeder
Vorschlag scheitern oder gelingen wird.

## Entscheidung

**(Vorschlag — Status `Proposed`, noch nicht ratifiziert.)**

u-boot scaffoldet Agent-Harness-Artefakte als **opt-in**, GF/BF-aware:

1. **Opt-in, nicht Core-`init`.** Auslieferung als eigenes Artefakt
   (`u-boot generate harness`) **oder** als Katalog-Template
   (`u-boot init --template harness`) — nicht als Default-`init`-
   Verhalten. Begründung: die Harness-Methodik ist eine *spezifische*
   Schule; sie darf nicht jedem u-boot-Projekt aufgezwungen werden.
2. **u-boot scaffoldet, der Agent reconciled.** u-boot emittiert das
   **GF-Skelett** (Doku-Outlines in Phase-1/Outline-Reife) und im
   erkannten **BF-Fall** zusätzlich einen **BF-Modus-Block** in
   `conventions.md` + einen **Reconciliation-Backlog-Stub** + leere
   Diskrepanz-Tabelle. u-boot führt **keine** Code-Inventur/Reverse-
   Engineering durch — das bleibt Agenten-Arbeit (Scope-Grenze oben).
3. **GF/BF wird automatisch gewählt** über die bestehende
   [`LH-FA-INIT-004`](../../../spec/lastenheft.md#lh-fa-init-004-bestehendes-projekt-erkennen)-Detection (kein neues Detektions-Konzept):
   leeres/neues Verzeichnis → GF-Skelett; erkanntes bestehendes Projekt
   → BF-Skelett + Modus-Block. `--assume-existing`/`--no-interactive`-
   Semantik aus [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a-interaktivität-und-automatisierung) gilt unverändert.
4. **Engine + Format wie [ADR-0009](0009-template-format-yaml-files.md)** (`text/template` + `template.yaml`),
   damit kein zweiter Stack entsteht.

## Konsequenzen

Positiv:

- u-boots Vision ([`LH-ZB-002`](../../../spec/lastenheft.md#lh-zb-002-produktvision)) wächst konsistent um eine Schicht, ohne
  neuen technischen Stack ([ADR-0009](0009-template-format-yaml-files.md)-Engine wiederverwendet).
- GF/BF kostet kein neues Konzept — [`LH-FA-INIT-004`](../../../spec/lastenheft.md#lh-fa-init-004-bestehendes-projekt-erkennen) trägt es schon.
- Dogfooding: u-boots eigenes Repo ist die Referenz-Ausgabe.

Negativ / Risiken:

- **Methodik-Kopplung.** u-boot-Output bindet sich an *eine*
  Harness-Schule (MR-*/LH-*/Sub-Areas/Trigger-Klassen/ADR). Opt-in
  mildert das, beseitigt es nicht.
- **BF-Erwartungslücke.** Nutzer könnten erwarten, dass u-boot die
  Reconciliation *macht*, nicht nur das Skelett legt. Muss in der Doku
  scharf abgegrenzt werden.
- **Lizenz/Attribution.** Falls Skelette von den
  `ai-harness-course`-Templates (`lab/templates/`, CC-BY für Markdown /
  MIT für Code) abgeleitet werden, ist Namensnennung zu klären, bevor
  abgeleitete Vorlagen mit u-boot ausgeliefert werden.
- **Spec-Wachstum.** Neue `LH-FA-*`-IDs + [`LH-ZB-002`](../../../spec/lastenheft.md#lh-zb-002-produktvision)-Visions-
  Erweiterung nötig — kein reines Adapter-Inkrement.

## Offene Fragen (vor `Accepted` zu beantworten)

1. **Subcommand-Form:** `generate harness` (Artefakt-Generator,
   [`LH-FA-GEN-001`](../../../spec/lastenheft.md#lh-fa-gen-001-generate-befehl)-Familie) vs. `init --template harness`
   (Template-Pfad, [ADR-0009](0009-template-format-yaml-files.md)). Tendenz: `generate harness`, weil es ein
   *Zusatz-Artefakt* in ein bestehendes oder frisches Projekt legt, kein
   vollständiges Projekt-Skelett ersetzt.
2. **BF-Skelett-Umfang:** nur Modus-Block + leere Stubs, oder
   zusätzlich ein Diskrepanz-Backlog-Gerüst mit vorbefüllten
   Sektions-Headern aus erkannten Projektdateien (deterministisch, ohne
   Inferenz)?
3. **Methodik-Verbindlichkeit:** fest verdrahtete Kurs-Konvention vs.
   konfigurierbares/austauschbares Harness-Profil.
4. **Lizenz/Attribution** der abgeleiteten Skelette (siehe
   §Konsequenzen).

## Folgepunkte

Dieses ADR liefert nur die Entscheidungs-Rahmung. Vor Implementierung:

- Ratifizierung (`Proposed` → `Accepted`) nach Klärung der §Offenen
  Fragen.
- Spec-Erweiterung: neue `LH-FA-*`-Anforderungen +
  [`LH-ZB-002`](../../../spec/lastenheft.md#lh-zb-002-produktvision)-Visions-Satz.
- Planning-Artefakt fuer die Umsetzung.
- Lizenz-/Attribution-Check gegen `ai-harness-course`.

Re-Evaluation-Trigger: konkrete Nutzer-/Team-Nachfrage nach
Harness-Scaffolding, oder ein zweites u-boot-Projekt, das die
Harness-Schicht manuell nachzieht (Dogfooding-Schmerz).
