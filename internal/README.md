# internal/

Nicht-exportierbare Go-Pakete für u-boot, nach dem hexagonalen
Architektur-Pattern geschnitten:
[`SPEC-012`](../spec/spezifikation.md#spec-012--schichten-und-verzeichnislayout-unter-internal)
(Layout),
[`SPEC-013`](../spec/spezifikation.md#spec-013--import-regel-tabelle-der-schichten)
und
[`SPEC-034`](../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)
(Importregeln und Enforcement). Schichten, Komponenten und
Sequenzen: [`spec/architecture.md`](../spec/architecture.md).

Coverage-Scope ist `./internal/...`
([`SPEC-032`](../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage),
[`SPEC-033`](../spec/spezifikation.md#spec-033--go-layout-des-repositories)).
Jede Schicht hat ein eigenes README mit Zweck und Verweis.
