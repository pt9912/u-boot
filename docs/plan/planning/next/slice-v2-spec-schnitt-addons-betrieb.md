# Slice V2: Lastenheft-Schnitt 2: Docker/Compose, Add-ons, Betrieb, Diagnose, Generatoren, Templates, Konfiguration

> **Status:** **geplant** (`next/`) — startet nach Schnitt 1.

**Welle:** `welle-spec-technik-stratum` (siehe [`roadmap.md`](../in-progress/roadmap.md)).
**Bezug:** Lastenheft §4.4 Docker/Compose, §4.5 Service-Add-ons, §4.6 Starten und Stoppen, §4.7 Diagnose, §4.8 Generatoren, §4.9 Template-System, §4.10 Konfigurationsdatei (Scope).
**Berührte Spec-Stellen:** `spec/spezifikation.md` §1–§6; `spec/lastenheft.md` §4.4–§4.10.
**Autor:** pt9912. **Datum:** 2026-10-02.

## Auslöser

Add-on-Kataloge (Images, Ports, Healthchecks), das YAML-Schema von `u-boot.yaml`, Doctor-Prüfliste und Generator-Formate sind
technische Festlegungen. Sie folgen den im Pilot bestimmten Schnittregeln.

## Ziel und Abgrenzung

**Ziel:** §4.4–§4.10 enthalten nur noch Vertrag; Kataloge, Schemata, Prüflisten und Defaults stehen wörtlich in der Spezifikation mit Kennung.

**Schnittregeln:** siehe [`slice-v2-spec-technik-stratum-grundlage`](../in-progress/slice-v2-spec-technik-stratum-grundlage.md) §Schnittregeln (bleibt im Lastenheft / zieht in die Spezifikation, Decken-Regel, Verbleib-Tabelle, Zeilenverweise).

**Ausdrücklich NICHT in diesem Slice:**

- **§4.1–§4.3** — Pilot-Slice, bereits geschnitten.
- **Projektkapitel §4.11–§4.13** — eigener Slice, weil sie als Randbedingungen gekürzt werden (andere Schnittart).


## Definition of Done

- [ ] §4.4–§4.10 geschnitten; Verbleib-Tabelle vollständig; `LH-*`-Überschriften und Anker unverändert.
- [ ] Spezifikation ergänzt (Katalog, `u-boot.yaml`-Schema, Doctor-Prüfungen, Generator-Formate, Defaults), ohne ADR- und Slice-Verweis; Zeilenverweise umgestellt.
- [ ] `make gates` grün (inklusive `make docs-check`), Exit-Code separat geprüft.
- [ ] Review durch eine andere Rolle als die Umsetzung (`harness/review.md`), Report unter `docs/reviews/`.

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

Offen. Beim Schließen: Verifikation-Evidence nach `harness/verification.md` (DoD, Sensoren, nicht ausgeführte
Sensoren, Carveouts), Steering-Loop-Lerneintrag, Delivery-Hash im Kopf.
