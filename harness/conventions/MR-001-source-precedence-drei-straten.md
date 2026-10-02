# MR-001 — Source Precedence mit drei Spec-Straten (Vollform)

- **Datum:** 2026-10-02 (Fassung 2; Fassung 1 vom 2026-07-24 fuehrte nur zwei Straten und liess das Technik-Stratum bewusst weg)
- **Geltungsbereich:** [`AGENTS.md`](../../AGENTS.md) §2,
  [`harness/README.md`](../README.md) §Source precedence,
  `.d-check.yml` (Referenzmatrix).
- **Ersetzt-Baseline-Regel:** [`grundlagen-source-precedence.md` §Source Precedence](../../.harness/baseline/v6.13.0/regelwerk/grundlagen-source-precedence.md#source-precedence)
- **Adaption:** u-boot fuehrt eine Source-Precedence mit **drei**
  Spec-Straten an den Raengen 1 bis 3 (Baseline-Default) und zwei repo-lokalen
  Zusatz-Raengen: Rang 5 fasst Roadmap und aktive Slices zusammen, Rang 6 fuehrt die
  ausfuehrbaren Vertraege (`Makefile`, `Dockerfile`, `.golangci.yml`, Workflows) als eigenen
  Rang; `AGENTS.md` und `harness/README.md` stehen als Raenge 9 und 10 nach Baseline-Regelwerk `modul-03-spec.md`:
  `contract_spec` ([`spec/lastenheft.md`](../../spec/lastenheft.md), Vertrag, das
  *Was*), `tech_spec` ([`spec/spezifikation.md`](../../spec/spezifikation.md),
  Technik, das *Wie genau*) und `view_spec`
  ([`spec/architecture.md`](../../spec/architecture.md), Sicht, das *Wodurch*).
  Konfliktregel: Lastenheft vor Spezifikation vor Architektur - die untere
  Schicht praezisiert, erweitert nie. Die **Decken-Regel** gilt fuer alle drei:
  kein Spec-Stratum nennt eine ADR, einen Slice, einen Carveout oder die
  Roadmap; die Spezifikation verweist nur aufwaerts auf das Lastenheft.
  Repo-Klasse: **Tooling/Referenz**.
- **ID-Schemata (Technik-Stratum):** Eine **Verfeinerung** einer einzelnen
  Anforderung traegt deren Kennung mit Buchstabensuffix (`<Anforderungs-ID>.a`,
  `.b`, ...); alles, was keine einzelne Anforderung verfeinert (Datenschemata,
  Defaults, Fehler-Codes, Metrik-Felder, externe Vertraege), traegt
  `SPEC-<NNN>` (dreistellig, fortlaufend je Datei). Eine `SPEC-*` ist eine
  Adresse, keine Anforderung. Die Sicht fuehrt `ARC-<NNN>` fuer Komponenten.
- **Begruendung:** Fassung 1 stuetzte sich auf die Optionalitaet des
  Technik-Stratums; das Lastenheft wuchs dadurch auf rund 3000 Zeilen und
  vermischte Vertrag und Technik (Schemata, Defaults, Build-Details). Das
  Regelwerk (v6.13.0) macht alle drei Straten obligatorisch: Technik im
  Vertrag zu falten verschiebt den Aenderungs-Prozess (nur per Change Request
  aenderbar, keine schaerfende ADR moeglich). Der Projektinhaber hat die
  Vollform am 2026-10-02 beschlossen.
- **Aufloesungs-Trigger:** permanent. Die Befuellung der Spezifikation ist
  Gegenstand der Welle `welle-spec-technik-stratum` (Roadmap).
