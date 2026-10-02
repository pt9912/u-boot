# Slice V2: Lastenheft-Schnitt 4 und Abschluss: Qualität, Akzeptanz, Traceability, Version 0.4.0

> **Status:** **Done** (2026-10-02, **Delivery-Hash: `3b98b30`**).

**Welle:** `welle-spec-technik-stratum` (siehe [`roadmap.md`](../in-progress/roadmap.md)).
**Bezug:** Lastenheft §5 Nichtfunktionale Anforderungen, §8 Qualitätsanforderungen, §9 Akzeptanzkriterien, §10–§16 (Abgrenzung, Risiken, MVP, Traceability, Offene Punkte, Glossar, Historie) (Scope).
**Berührte Spec-Stellen:** `spec/lastenheft.md` §5, §8–§16 (Version, Historie); `spec/spezifikation.md` §1–§7.
**Autor:** pt9912. **Datum:** 2026-10-02.

## Auslöser

Die Qualitätsanforderungen tragen Messmethoden, die teils technische Festlegungen enthalten (Schwellen, Messaufbau). Zum Abschluss
bekommt das Lastenheft den neuen Vertragsstand: Version 0.4.0 mit Historie-Zeile, Lesehinweis auf die drei Straten und
konsistente Traceability-Matrix.

## Ziel und Abgrenzung

**Ziel:** §5 und §8–§16 sind geschnitten; das Lastenheft trägt Version 0.4.0 mit Historie-Zeile; die Spezifikation ist vollständig befüllt und frei von Platzhaltern.

**Schnittregeln:** siehe [`slice-v2-spec-technik-stratum-grundlage`](../done/slice-v2-spec-technik-stratum-grundlage.md) §Schnittregeln (bleibt im Lastenheft / zieht in die Spezifikation, Decken-Regel, Verbleib-Tabelle, Zeilenverweise).

**Ausdrücklich NICHT in diesem Slice:**

- **Änderung von Anforderungen oder Schwellenwerten** — der Schnitt verschiebt nur; Coverage-Schwelle, Exit-Codes usw. bleiben wertgleich.
- **Architektur-Sicht** — eigener Slice.


## Definition of Done

- [x] §5, §8–§16 geschnitten, Verbleib-Tabelle der gesamten Welle vollständig (jede Quelle genau ein Ziel, nichts verwaist); Version 0.4.0, Historie-Zeile, Traceability-Matrix konsistent.
- [x] Spezifikation vollständig (§1–§7, Historie), `SPEC-<NNN>` fortlaufend, keine Platzhalter; `LH-*`-Anker unverändert.
- [x] `make gates` grün (inklusive `make docs-check`), Exit-Code separat geprüft.
- [x] Review durch eine andere Rolle als die Umsetzung (`harness/review.md`), Report unter `docs/reviews/`.

## Tranchen

| T | Inhalt |
| - | ------ |
| T1 | **Inventur und Verbleib-Tabelle** für §5, §8–§16. |
| T2 | **Übernahme** der technischen Messaufbauten und Schwellen-Festlegungen; Wertgleichheit der Schwellen gegen die Gates prüfen. |
| T3 | **Vertragsstand:** Version 0.4.0, Historie-Zeile, Lesehinweis auf die drei Straten, Traceability-Matrix. |
| T4 | **Gegenlesen der gesamten Welle,** `make gates`, Review, Welle-Closure-Notiz. |

## Risiken

- Wertdrift bei Schwellen (zum Beispiel Coverage 90 %): T2 vergleicht die Zahlen mit `Makefile` und `.golangci.yml`.

## Closure

- **Geliefert (`3b98b30`, Nachbesserung aus dem Review `76988a5`):** Lastenheft §5, §8 und §14 geschnitten (Beispielinstanz der Minimalkontrakt-Ausgabe, CI-Komposition und Lint-Profil-Zusammensetzung als Verfeinerungen, Go-Toolchain-Pin); Version **0.4.0** mit Historie-Zeile, neue Lesehinweis-Anforderung [`LH-LESE-003`](../../../../spec/lastenheft.md#lh-lese-003--dokumentenordnung) (Dokumentenordnung) samt Traceability-Zeile; die Spezifikation ist in §1–§7 befüllt: aus dem Bestand neu Standardwerte, Egress-Default-Allowlist, Code-Registry der `doctor`-Prüfungen (der Drift-Test liest jetzt die Spezifikation), Telemetrie-Aussage und externe Verträge (Docker/Compose/Podman, Add-on-Katalog, Devcontainer-Quellen).
- **Review (unabhängig, Rolle getrennt):** `docs/reviews/2026-10-02-welle-spec-technik-stratum.md` — 0 HIGH, 4 MEDIUM, 6 LOW, 4 INFO. Alle MEDIUM und LOW behoben (Enum in CLI-007, Feature-Quellen, Coverage-/Go-Pins, Verweis-Artefakte aus der Zeilenumstellung, Lesbarkeit, READMEs, Wortlaut); verbleibend bewusst: die Coverage-Bootstrap-Aussage im Lastenheft (Vertragsänderung, eigener Slice `slice-v2-spec-lastenheft-coverage-default`), `§611` in den erzeugten `.gitignore`-Vorlagen (Produktausgabe) und die Fehlertext-Änderungen (INFO).
- **Vollständigkeit:** Review-Abgleich aller entfernten Lastenheft-Zeilen gegen Lastenheft plus Spezifikation: nichts verloren (9 bewusste Verkürzungen, keine Zusage abgeschwächt); alle `LH-*`-Anker unverändert.
- **Sensoren:** `make gates` grün (lint, test, coverage-gate, docs-check). Nicht ausgeführt: `make ci`, `make test-docker` (keine Produktlogik berührt; Hilfe- und Fehlertexte einzelner Meldungen tragen nun Kennungen statt Zeilennummern).
- **Lerneintrag:** Verweise der Form „Zeile N“ sind ein Drift-Risiko ohne Sensor; Auflösung über den Stand des Kommentar-Commits funktioniert, erzeugt aber Artefakte, wenn der Verweis auf einen noch älteren Stand zielt. Ein Review der Zeilen-Umstellung gehört in jeden solchen Schnitt.

### Verbleib-Tabelle

| Quelle (Lastenheft) | Ziel (Spezifikation) | Ort |
|---|---|---|
| §5.1 [LH-NFA-USE-004](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) (Beispielinstanz, 1702–1711) | [SPEC-015](../../../../spec/spezifikation.md#spec-015--beispielinstanz-einer-minimalkontrakt-ausgabe-doctor---json) | spezifikation.md §2 |
| §8 [LH-QA-003](../../../../spec/lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions) (Pflicht-Komposition, 2091–2102) | [LH-QA-003.a](../../../../spec/spezifikation.md#lh-qa-003a--komposition-der-ci-pipeline-github-actions) | spezifikation.md §1 |
| §8 [LH-QA-004](../../../../spec/lastenheft.md#lh-qa-004--linting-solid-nahes-lint-profil) (Profil-Komposition, 2112–2115) | [LH-QA-004.a](../../../../spec/spezifikation.md#lh-qa-004a--zusammensetzung-des-lint-profils) | spezifikation.md §1 |
| §14 [LH-OPEN-001](../../../../spec/lastenheft.md#lh-open-001--implementierungssprache-entschieden) (Toolchain-Pin, 2554) | [SPEC-016](../../../../spec/spezifikation.md#spec-016--go-toolchain-mindestversion-und-dockerfile-pin) | spezifikation.md §3 |
