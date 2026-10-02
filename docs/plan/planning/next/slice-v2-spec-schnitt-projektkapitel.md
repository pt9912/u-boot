# Slice V2: Lastenheft-Schnitt 3: Projektkapitel als Randbedingungen, Schnittstellen und Daten

> **Status:** **geplant** (`next/`) — startet nach Schnitt 2.

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

**Schnittregeln:** siehe [`slice-v2-spec-technik-stratum-grundlage`](../in-progress/slice-v2-spec-technik-stratum-grundlage.md) §Schnittregeln (bleibt im Lastenheft / zieht in die Spezifikation, Decken-Regel, Verbleib-Tabelle, Zeilenverweise).

**Ausdrücklich NICHT in diesem Slice:**

- **Neue Anforderungs-Kennungen** — alle bestehenden `LH-*`-Kennungen bleiben; ein Umbau zu `-RB-`-Kennungen wäre eine Vertragsänderung mit Link-Bruch in ADRs, Slices und Tests und ist nicht beschlossen.
- **Inhaltliche Änderung der Importregeln** — sie bleiben in `spec/architecture.md`; die Angleichung der Sicht an die Vorlage ist ein eigener Slice.


## Definition of Done

- [ ] §4.11–§4.13 auf Vorgabe plus Nachweis gekürzt; §6.2, §6.3, §7 geschnitten; Verbleib-Tabelle vollständig; Kennungen und Anker unverändert.
- [ ] Spezifikation ergänzt (Build-/CI-Festlegungen, Doku-Struktur, Dateischnittstellen, Datenformate) bzw. Verweis auf die Architektur-Sicht an den passenden Stellen, ohne ADR- und Slice-Verweis.
- [ ] `make gates` grün (inklusive `make docs-check`), Exit-Code separat geprüft.
- [ ] Review durch eine andere Rolle als die Umsetzung (`harness/review.md`), Report unter `docs/reviews/`.

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

Offen. Beim Schließen: Verifikation-Evidence nach `harness/verification.md` (DoD, Sensoren, nicht ausgeführte
Sensoren, Carveouts), Steering-Loop-Lerneintrag, Delivery-Hash im Kopf.
