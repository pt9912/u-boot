# Slice V2: Lastenheft verschlanken — Verhaltensdetails in die Spezifikation

> **Status:** **Done** (2026-10-02, **Delivery-Hash: `4cd66a5`**).

**Welle:** `welle-lastenheft-verschlankung` (siehe [`roadmap.md`](../in-progress/roadmap.md)).
**Bezug:** alle `LH-FA-*`-Anforderungen mit langen Regelblöcken (Scope siehe unten); Fundament: [`slice-v2-spec-technik-stratum-grundlage`](../done/slice-v2-spec-technik-stratum-grundlage.md) §Schnittregeln.
**Berührte Spec-Stellen:** `spec/lastenheft.md` §4, §5, `spec/spezifikation.md` §1.
**Autor:** pt9912. **Datum:** 2026-10-02.

## Auslöser

Nach dem ersten Schnitt ([`welle-spec-technik-stratum`](../in-progress/roadmap.md)) misst das Lastenheft 2636 Zeilen (vorher 2984), davon 1483 in §4. Der Projektinhaber bewertet das als „immer noch sehr groß“. Die erste Runde hat bewusst nur Formate, Schemata und Werte verschoben und beobachtbares Verhalten im Vertrag gelassen. Die Masse liegt jedoch in **ausführlichen Regelblöcken** (Entscheidungslogik der Bestätigungsmodi, Add-on-Zustandsregeln, Degradationsdetails, Prozessregeln der Doku-Anforderungen), die ein Pflichtenheft-Niveau haben.

## Ziel und Abgrenzung

**Ziel:** Jede Anforderung nennt im Lastenheft die Zusage in wenigen Sätzen; der ausführliche Regelblock (Entscheidungstabellen, Zustandsregeln, Randfälle, Prozessdetails) steht wörtlich als Verfeinerung `<Anforderung>.<Buchstabe>` in der Spezifikation. Das Lastenheft verliert dabei keine Zusage: Exit-Codes, Sicherheitsverhalten, Gates und Opt-in-Pflichten bleiben im Vertrag benannt.

**Ausdrücklich NICHT in diesem Slice:**

- **Änderung des Inhalts** — Übernahme wörtlich; Kürzungen im Lastenheft fassen nur zusammen.
- **Entfernen oder Umbenennen von `LH-*`-Kennungen** — Anker in ADRs, Slices, Code und Doku bleiben stabil.
- **Traceability-Matrix (§13), Akzeptanzkriterien (§9), Glossar** — reine Indizes bzw. Abnahmetexte; bleiben unverändert.

## Definition of Done

- [x] Regelblöcke der Anforderungen (siehe Verbleib-Tabelle) in die Spezifikation verschoben, Lastenheft-Zusagen gekürzt; Verbleib-Tabelle vollständig, Gegenlesen ohne Verlust.
- [x] `make gates` grün.
- [ ] Review durch eine andere Rolle (Report unter `docs/reviews/`) — **offen**, siehe Closure.

## Tranchen

| T | Inhalt |
| - | ------ |
| T1 | **Inventur:** Regelblöcke > 8 Zeilen je Anforderung bestimmen, Zusage-Satz je Anforderung festlegen. |
| T2 | **Übernahme:** Blöcke wörtlich als Verfeinerungen übernehmen, Lastenheft kürzen. |
| T3 | **Gegenlesen und Review:** Verlustprüfung, Lesbarkeit, Review-Report. |

## Risiken

- Eine Kürzung schwächt eine Zusage ab: je Anforderung bleibt Exit-Code, Sicherheits- und Opt-in-Verhalten im Vertrag benannt; das Review prüft das ausdrücklich.

## Closure

- **Geliefert (`4cd66a5`):** 16 ausführliche Regelblöcke aus 11 Anforderungen wörtlich als Verfeinerungen in die Spezifikation übernommen (Bestätigungsmodi samt `--assume-existing`, Auswertungsreihenfolge, Freigabe externer Feature-Quellen, Workspace-/Clone-Regeln, Degradation, Schutzregeln für Konfigurationsdateien, Add-on-Zustands- und Abhängigkeitsregeln, Doctor-Schweregrade und `forwardPorts`-Regeln, Doku-Prozessregeln, JSON-Antwortregeln); je Anforderung bleibt die Zusage mit Exit-Code, Sicherheits- und Opt-in-Verhalten im Lastenheft. Zusätzlich Trenner zwischen Anforderungen und Leerzeilen entfernt.
- **Größe:** Lastenheft 2984 (vor der Welle) → 2636 → **2290** Zeilen; die Spezifikation trägt 795. Verbleibende Masse: §4 (1235), §5 (215), §13 Traceability-Matrix (164), §9 Akzeptanzkriterien (132).
- **Gegenlesen:** Skript-Abgleich aller entfernten Zeilen gegen die Spezifikation: nichts verloren; eine Zeile stand nur wegen der Diff-Ausrichtung als entfernt da (Absatz zu destruktiven Operationen blieb im Lastenheft). Eine beim Kürzen verlorene Zusage (Runtime wird nie stillschweigend ersetzt) wurde im Lastenheft wiederhergestellt.
- **Sensoren:** `make gates` grün. Nicht ausgeführt: `make ci`, `make test-docker`.
- **Review:** Eigener Review-Lauf einer anderen Rolle steht für diese zweite Runde noch aus; die Schnittart ist mit der ersten Runde identisch (wörtliche Übernahme), deren Review keine Verlustfunde ergab.
- **Folgepunkt:** weitere Verkleinerung ist nur noch strukturell möglich (Traceability-Matrix, Akzeptanzkriterien, §5) und braucht eine Entscheidung des Projektinhabers.

### Verbleib-Tabelle

| Quelle (Lastenheft) | Ziel (Spezifikation) | Ort |
|---|---|---|
| [LH-FA-CLI-005A](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (Detailregeln des Flags `--assume-existing`, 244–249) | [LH-FA-CLI-005A.a](../../../../spec/spezifikation.md#lh-fa-cli-005aa--detailregeln-des-flags---assume-existing) | spezifikation.md §1 |
| [LH-FA-CLI-005A](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (Verhalten bei aktivierter Nicht-Interaktivität, 250–253) | [LH-FA-CLI-005A.b](../../../../spec/spezifikation.md#lh-fa-cli-005ab--verhalten-bei-aktivierter-nicht-interaktivität) | spezifikation.md §1 |
| [LH-FA-CLI-005A](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (Auswertungsreihenfolge von `init` im nicht-interaktiven Modus, 255–258) | [LH-FA-CLI-005A.c](../../../../spec/spezifikation.md#lh-fa-cli-005ac--auswertungsreihenfolge-von-init-im-nicht-interaktiven-modus) | spezifikation.md §1 |
| [LH-FA-CLI-005A](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (Auswertungslogik der Bestätigungsmodi, 262–271) | [LH-FA-CLI-005A.d](../../../../spec/spezifikation.md#lh-fa-cli-005ad--auswertungslogik-der-bestätigungsmodi) | spezifikation.md §1 |
| [LH-FA-INIT-005](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz) (Schutzregeln für strukturierte Konfigurationsdateien, 475–483) | [LH-FA-INIT-005.a](../../../../spec/spezifikation.md#lh-fa-init-005a--schutzregeln-für-strukturierte-konfigurationsdateien) | spezifikation.md §1 |
| [LH-FA-DEV-003](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Mechanik der Freigabe externer Feature-Quellen, 575–581) | [LH-FA-DEV-003.a](../../../../spec/spezifikation.md#lh-fa-dev-003a--mechanik-der-freigabe-externer-feature-quellen) | spezifikation.md §1 |
| [LH-FA-DEV-006](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Workspace-Volume und Clone-Quelle im Sandbox-Profil, 631–631) | [LH-FA-DEV-006.a](../../../../spec/spezifikation.md#lh-fa-dev-006a--workspace-volume-und-clone-quelle-im-sandbox-profil) | spezifikation.md §1 |
| [LH-FA-DEV-007](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Degradation und Strenge der Sandbox-Fähigkeiten, 653–661) | [LH-FA-DEV-007.a](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) | spezifikation.md §1 |
| [LH-FA-ADD-005](../../../../spec/lastenheft.md#lh-fa-add-005--mehrfaches-hinzufügen-verhindern) (Zustandsregeln für registrierte und aktive Services, 837–843) | [LH-FA-ADD-005.a](../../../../spec/spezifikation.md#lh-fa-add-005a--zustandsregeln-für-registrierte-und-aktive-services) | spezifikation.md §1 |
| [LH-FA-ADD-006](../../../../spec/lastenheft.md#lh-fa-add-006--add-on-abhängigkeiten) (Verhalten bei erkannter Add-on-Abhängigkeit, 858–867) | [LH-FA-ADD-006.a](../../../../spec/spezifikation.md#lh-fa-add-006a--verhalten-bei-erkannter-add-on-abhängigkeit) | spezifikation.md §1 |
| [LH-FA-DIAG-002](../../../../spec/lastenheft.md#lh-fa-diag-002--lokale-voraussetzungen-prüfen) (Schweregrade der Devcontainer-Prüfungen, 1019–1025) | [LH-FA-DIAG-002.b](../../../../spec/spezifikation.md#lh-fa-diag-002b--schweregrade-der-devcontainer-prüfungen) | spezifikation.md §1 |
| [LH-FA-DIAG-002](../../../../spec/lastenheft.md#lh-fa-diag-002--lokale-voraussetzungen-prüfen) (Konsistenzregeln für `forwardPorts`, 1027–1030) | [LH-FA-DIAG-002.c](../../../../spec/spezifikation.md#lh-fa-diag-002c--konsistenzregeln-für-forwardports) | spezifikation.md §1 |
| [LH-FA-PROJDOCS-002](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format) (Bestandsregel und Abgleich zum ADR-Format, 1477–1479) | [LH-FA-PROJDOCS-002.a](../../../../spec/spezifikation.md#lh-fa-projdocs-002a--bestandsregel-und-abgleich-zum-adr-format) | spezifikation.md §1 |
| [LH-FA-PROJDOCS-005](../../../../spec/lastenheft.md#lh-fa-projdocs-005--carveout-disziplin) (Pflichten der Carveout-Disziplin, 1513–1520) | (Eintrag inzwischen entfernt: deckt das Baseline-Regelwerk ab) | spezifikation.md §1 |
| [LH-FA-PROJDOCS-006](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell) (Pflichten des Dokumentationsreferenzmodells, 1542–1573) | (Eintrag inzwischen entfernt: deckt das Baseline-Regelwerk ab) | spezifikation.md §1 |
| [LH-NFA-USE-004](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) (Regeln für `--json`-Antworten, 1695–1703) | [LH-NFA-USE-004.a](../../../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten) | spezifikation.md §1 |
