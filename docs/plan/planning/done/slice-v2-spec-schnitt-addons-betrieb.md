# Slice V2: Lastenheft-Schnitt 2: Docker/Compose, Add-ons, Betrieb, Diagnose, Generatoren, Templates, Konfiguration

> **Status:** **Done** (2026-10-02, **Delivery-Hash: `5d4716e`**).

**Welle:** `welle-spec-technik-stratum` (siehe [`roadmap.md`](../in-progress/roadmap.md)).
**Bezug:** Lastenheft §4.4 Docker/Compose, §4.5 Service-Add-ons, §4.6 Starten und Stoppen, §4.7 Diagnose, §4.8 Generatoren, §4.9 Template-System, §4.10 Konfigurationsdatei (Scope).
**Berührte Spec-Stellen:** `spec/spezifikation.md` §1–§6; `spec/lastenheft.md` §4.4–§4.10.
**Autor:** pt9912. **Datum:** 2026-10-02.

## Auslöser

Add-on-Kataloge (Images, Ports, Healthchecks), das YAML-Schema von `u-boot.yaml`, Doctor-Prüfliste und Generator-Formate sind
technische Festlegungen. Sie folgen den im Pilot bestimmten Schnittregeln.

## Ziel und Abgrenzung

**Ziel:** §4.4–§4.10 enthalten nur noch Vertrag; Kataloge, Schemata, Prüflisten und Defaults stehen wörtlich in der Spezifikation mit Kennung.

**Schnittregeln:** siehe [`slice-v2-spec-technik-stratum-grundlage`](../done/slice-v2-spec-technik-stratum-grundlage.md) §Schnittregeln (bleibt im Lastenheft / zieht in die Spezifikation, Decken-Regel, Verbleib-Tabelle, Zeilenverweise).

**Ausdrücklich NICHT in diesem Slice:**

- **§4.1–§4.3** — Pilot-Slice, bereits geschnitten.
- **Projektkapitel §4.11–§4.13** — eigener Slice, weil sie als Randbedingungen gekürzt werden (andere Schnittart).


## Definition of Done

- [x] §4.4–§4.10 geschnitten; Verbleib-Tabelle vollständig; `LH-*`-Überschriften und Anker unverändert.
- [x] Spezifikation ergänzt (Katalog, `u-boot.yaml`-Schema, Doctor-Prüfungen, Generator-Formate, Defaults), ohne ADR- und Slice-Verweis; Zeilenverweise umgestellt.
- [x] `make gates` grün (inklusive `make docs-check`), Exit-Code separat geprüft.
- [x] Review durch eine andere Rolle als die Umsetzung (`harness/review.md`), Report unter `docs/reviews/`.

## Tranchen

| T | Inhalt |
| - | ------ |
| T1 | **Inventur und Verbleib-Tabelle** für §4.4–§4.10. |
| T2 | **Übernahme Betrieb:** §4.4–§4.7 (Compose-Annahmen, Add-on-Katalog, Up/Down-Stabilisierung, Doctor-Prüfungen). |
| T3 | **Übernahme Artefakte:** §4.8–§4.10 (Generator-Formate, Template-Schema, `u-boot.yaml`-Schema). |
| T4 | **Gegenlesen,** `make gates`, Review. |

## Risiken

- Das Schema von `u-boot.yaml` in §4.10 wird von Code und Tests zitiert: T4 sucht Zitate (`grep -rn "lastenheft.md" internal/`) und stellt sie um.

## Closure

- **Geliefert (`5d4716e`):** Lastenheft §4.4–§4.10 geschnitten. In die Spezifikation gezogen: Mechanismus der Docker-Prüfung samt Drop-in-Regel (Verfeinerung der Diagnose-Anforderung) und das YAML-Schema der Projektkonfiguration. Das Lastenheft nennt an der Stelle weiter die Pflichtfelder und optionalen Schlüssel in Prosa. Alle `LH-*`-Überschriften und Anker unverändert.
- **Befund:** §4.4–§4.10 sind fast durchgehend beobachtbares Verhalten (Add-on-Zustandsregeln, Bestätigungsmodi, Doctor-Prüfliste); nach der Schnittregel bleibt es Vertrag. Der Technik-Anteil dieser Kapitel ist gering.
- **Gegenlesen:** 48 entfernte Zeilen, alle in der Spezifikation wiederzufinden; die einzige Abweichung ist die zusammengeführte Einleitungszeile des Schemas (bewusst, im Lastenheft ersetzt).
- **Zeilenverweise:** keine weiteren (alle Verweise sind seit Schnitt 1 auf Kennungen umgestellt).
- **Sensoren:** `make gates` grün. Nicht ausgeführt: `make ci`, `make test-docker`.
- **Review:** unabhängiges Review über die gesamte Welle am Abschluss.

### Verbleib-Tabelle

| Quelle (Lastenheft) | Ziel (Spezifikation) | Ort |
|---|---|---|
| §4.7 [LH-FA-DIAG-002](../../../../spec/lastenheft.md#lh-fa-diag-002--lokale-voraussetzungen-prüfen) (Docker-Prüfung, Mechanismus, 1006–1016) | [LH-FA-DIAG-002.a](../../../../spec/spezifikation.md#lh-fa-diag-002a--prüfung-von-docker-und-docker-kompatiblen-drop-ins) | spezifikation.md §1 |
| §4.10 [LH-FA-CONF-002](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (YAML-Schema, 1261–1295) | [SPEC-005](../../../../spec/spezifikation.md#spec-005--schema-der-projektkonfiguration-u-bootyaml) | spezifikation.md §2 |
