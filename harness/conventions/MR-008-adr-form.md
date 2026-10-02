# MR-008 — ADR-Form: Bestand lean+grandfathered, neue ADRs MADR (per CR)

- **Datum:** 2026-07-24
- **Geltungsbereich:** ADR-Form-Politik; verweist auf den CR-Slice
  [`slice-cr-adr-format-madr`](../../docs/plan/planning/done/slice-cr-adr-format-madr.md).
  Aendert `spec/lastenheft.md` NICHT von hier aus.
- **Ersetzt-Baseline-Regel:** [`modul-04-adrs.md` §Ziel-Form: ADR (MADR)](../../.harness/baseline/v6.13.0/regelwerk/modul-04-adrs.md#ziel-form-adr-madr)
- **Adaption:** Das vendored ADR-Template (MADR-/Nygard-Stil) kollidiert mit dem
  heutigen [`LH-FA-PROJDOCS-002`](../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)
  (Status/Datum als Inline-Felder statt `##`-Ueberschriften; Titel-/Superseded-
  Format; zusaetzliche Pflicht-Sections). Die Angleichung aendert das
  **Vertrags-Stratum** und ist damit ein **Change Request**, nicht per
  conventions-MR moeglich. Die Aenderung traegt der CR-Slice, nicht dieser Block.
  **CR ausgefuehrt (2026-07-24):** [`LH-FA-PROJDOCS-002`](../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format) traegt jetzt die
  MADR-Form; die zum CR-Zeitpunkt Accepted ADRs (`0001`-`0010`, `0013`) bleiben
  lean + immutabel (grandfathered, Hard Rule); Proposed (`0011`, `0012`) und alle
  neuen ADRs sind MADR-konform.
- **Begruendung:** conventions.md ist form-bringend, nicht vertrags-aendernd; eine
  Vertrags-Anforderung darf hier nur referenziert werden. Das Template-Feld
  "Schaerft" deckt sich mit u-boots Referenzmodell
  ([`ADR-0013`](../../docs/plan/adr/0013-dokumentationsreferenzmodell.md) /
  [`LH-FA-PROJDOCS-006`](../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell)):
  die `Schaerft`-Aufwaerts-Deklaration ist die Aenderungskopplung Spec-ADR.
- **Aufloesungs-Trigger:** erledigt (CR ausgefuehrt 2026-07-24; Delivery-Hash im
  CR-Slice). Die status-basierte Grandfather-Grenze gilt permanent; neue
  Accepted-ADRs entstehen bereits MADR-konform, kein neuer Grandfather noetig.
- **Nachtrag Baseline-`v3.5.2` (2026-07-25):** Die Baseline schaerft, dass
  "Change Request" **kein Harness-Konstrukt** ist (kein `CR-*`-Schema, keine
  eigene Datei, kein Gate), sondern der *externe* Vorgang der Vertrags-
  vereinbarung; im Repo hinterlaesst ein angenommener CR nur einen **Fussabdruck**
  (Version-Bump des Lastenhefts + Zeile in dessen `## Historie` + die geaenderten
  `LH-*`). Dazu die Hard Rule: **weder ADR noch Slice duerfen `LH-*` je aendern**
  - sie referenzieren nur. u-boots Praxis ist damit **inhaltlich** vereinbar
  (die Entscheidung fiel ausserhalb des Repos, der Slice war nur das
  Ausfuehrungs-Vehikel), der **Fussabdruck fehlt aber**: `spec/lastenheft.md`
  hat keine `## Historie` und steht unveraendert auf Version `0.1.0`, obwohl
  [`LH-FA-PROJDOCS-002`](../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format) am 2026-07-24 geaendert wurde. **Nachgezogen am
  2026-07-25** im Folge-Slice
  [`slice-harness-lastenheft-historie-cr-fussabdruck`](../../docs/plan/planning/done/slice-harness-lastenheft-historie-cr-fussabdruck.md):
  Das Lastenheft traegt jetzt Version `0.2.0`, Status `Accepted` und einen
  Abschnitt Historie mit der MADR-Zeile. Die Bezeichnung "CR-Slice" oben bleibt
  als historischer Name stehen, meint aber das **Vehikel**, nicht die
  Entscheidungs-Autoritaet - die lag beim Projektinhaber.
