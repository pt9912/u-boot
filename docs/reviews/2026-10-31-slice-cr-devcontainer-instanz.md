# Review-Report: slice-cr-devcontainer-instanz — 2026-10-31

**Review-Art:** Plan — Wortlaut-Review gegen den Plan (Alt/Neu je Stelle), die Entscheidungen des Projektinhabers (2026-10-31) und die Hard Rules.

**Gegenstand:** Commit `772685aa` (Basis `f7d89606`), Änderung an `spec/lastenheft.md` (Version 0.25.0) und ein Häkchen im Slice [`slice-cr-devcontainer-instanz`](../plan/planning/in-progress/slice-cr-devcontainer-instanz.md). Der Commit war beim Review bereits auf origin/main.

**Skill:** `.harness/skills/reviewer.md` @ 3ec45359 ·
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-31

**Eingangs-Kontext:**

- Slice-Plan [`slice-cr-devcontainer-instanz`](../plan/planning/in-progress/slice-cr-devcontainer-instanz.md) (§1: Entscheidungen, Stellen 1 bis 6 Alt/Neu, Tabelle Nicht-Zusagen; §2 DoD)
- [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Nachbarn [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle)
- `AGENTS.md` (Hard Rules 3.4, 3.5, 3.7, 3.9, 3.10)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der Slice [`slice-cr-devcontainer-instanz`](../plan/planning/in-progress/slice-cr-devcontainer-instanz.md) trägt nach dem Commit überholte Zustandsmarker: „`spec/lastenheft.md` ist nicht geändert und wird erst nach der Freigabe geschrieben“, „ENTWURF, nicht freigegeben“, „voraussichtlich `0.25.0`“, „Änderung an `spec/lastenheft.md` — erst nach der Freigabe“; nur ein DoD-Häkchen wurde gesetzt. | `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt) | [`slice-cr-devcontainer-instanz`](../plan/planning/in-progress/slice-cr-devcontainer-instanz.md) (Datei) · „Dieser Wortlaut ist ein ENTWURF zur Freigabe“ | nein — Review-only | Zustandsmarker im Plan nach Ablauf nicht nachgezogen |
| F-2 | LOW | Der Neutext in [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) formuliert den Prüfzustand: „Gemessen ist das mit zwei, drei und vier …“ und „andere Hosts und Engines als gemessen“. Es ist das einzige „gemessen“ im Hauptteil des Lastenhefts; die Historie 0.24.0 hat Anforderungstexte von „Messumgebung und Prüfzustand“ auf die Zusage-Menge im Wortlaut umgestellt. Der Inhalt (Stückzahlen, Grenzen, Profil, Host-Klasse) ist Zusage-Menge; es stehen keine Versionen, Läufe oder konkrete Hosts darin. | Historie 0.24.0 in [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)-Umfeld (§16), Maintainability | `spec/lastenheft.md` · „Gemessen ist das mit zwei, drei und vier gleichzeitig laufenden Instanzen“ | nein — Review-only | Prüfzustands-Verb in Zusage-Menge |
| F-3 | INFO | Die Zusage „Jede Instanz hat ein eigenes Workspace-Volume …“ nennt nur benannte Instanzen. Ob der Container ohne Instanz neben Instanzen mitgezählt wird („bis zu vier“) und ob er von ihnen isoliert ist, sagt der Text nicht; die Aussage steht außerhalb der Definition. | Maintainability | `spec/lastenheft.md` · „Bis zu vier gleichzeitig laufende Instanzen desselben Projekts teilen“ | nein — Review-only | Zusage-Rand zum Container ohne Instanz offen |
| F-4 | INFO | Die Definition „Instanz“ steht in [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten); [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) sagt unverändert „Volume je Instanz“ im Sinne des einzelnen Devcontainers (Entscheidung 2, vom Projektinhaber so gesetzt). Ein Leser von [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) mit dem neuen Begriff liest „je Instanz“ enger als gemeint. | Entscheidung 2 des Projektinhabers | `spec/lastenheft.md` · „Der Daemon-Speicher liegt in einem benannten Volume je Instanz“ | nein — Review-only | Doppelter Wortgebrauch „Instanz“ (bekannt, entschieden) |
| F-5 | INFO | Der Lastenheft-Text sagt `--instance` an `up`, `down`, `shell` zu, ohne dass der Code es kann (Zielbild, Entscheidung 5). Nutzer-Doku, Befehlstabelle [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) (ohne Flags) und der Code stehen bis zum Code-Slice dahinter; das ist Zuständigkeit des Code-Slice und des Verifiers. | Entscheidung 5 | `spec/lastenheft.md` · „Mit `--instance <name>` stoppt er nur den Container dieser Instanz“ | nein — Review-only | Zielbild vor Code |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Wörtlichkeit Stellen 1 bis 6 (maschineller Abgleich Slice-„Neu“ gegen Lastenheft, `**` entfernt, Links `#…`, Datum der Historie-Zeile 2026-10-31): `up`, `down`, `shell` (erster Satz), neuer Spiegelstrich, Fehlerklassen (Anfang bis `gibt`, Rest wie Alt), Nicht-Ziele, Historie-Zeile | geprüft, ohne Befund: jeder Neutext stimmt wörtlich; nichts Zusätzliches (kein Beleg, Begleittext, Slice-/ADR-/Protokoll-Verweis); nichts Gefordertes ausgelassen |
| Entscheidungen 1 bis 5 im Text (Namensregel `[a-z0-9-]`, Rand, 1 bis 32, leer und `=`, Exit-Code `10`; Definition einmal in [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)/`-007` unverändert; Flag an `up`/`down`/`shell`; Fassung (i) mit Hinweis auf `--instance`; Isolation bis vier, Sandbox, nicht Podman, Standardprofil nicht isoliert, Messbedingung 1 CPU und 3 GiB nur direkt nach dem Start) | geprüft, ohne Befund |
| Konsistenz im Lastenheft: Singular „Devcontainer des Projekts“ (durch die Definition in [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) für den Abschnitt gedeckt; [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), -003, -006, -007 nicht berührt); Fehlerklassen und [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Namensvalidierung gehört zu `10`); `--quiet`/`--json` (kein neuer Text, Wortlaut sagt nichts zum Schema); Obergrenze „bis zu vier“ (ausdrücklich „Nicht zugesagt sind mehr als vier“, einleitend „mehrere Instanzen … können gleichzeitig laufen“ ohne Grenze) | geprüft, ohne Befund außer F-3, F-4 |
| Decken-Regel `AGENTS.md` §3.4: im Neutext keine Verweise auf ADR, Slice, Welle, Roadmap, Carveout, Commit; Links nur auf Lastenheft-Anker; `make docs-check` | geprüft, ohne Befund |
| Version und Historie: Kopf `0.25.0`; Historie-Zeile `0.25.0 \| 2026-10-31 \| … \| Vereinbarung mit dem Projektinhaber` in der Form der Zeile 0.24.0; Tabelle intakt; `git show --stat` nennt nur `spec/lastenheft.md` und den Slice; keine Änderung an [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung), Zeile 416, älteren Historie-Zeilen, `spec/spezifikation.md`, `spec/architecture.md` | geprüft, ohne Befund |
| Kommentar-Muster `AGENTS.md` §3.7/§3.9 im Neutext (Chronik, abwesender Text, Konjunktiv) | geprüft, ohne Befund |
| Revert-Frage (Commit vor Review gepusht) | geprüft: kein Befund verlangt Revert; F-1 und F-2 sind Korrekturen im Folge-Commit |
| `make gates`, `make test-docker`, `make verify-depguard`, Code | nicht geprüft: Doku-Commit, kein Code; `make gates` ist Sache des Verifiers, es wurden keine Container gestartet |

Sensors: `make docs-check` — 943 Dateien, 0 Befunde. `make ablage-pruefung` — 26 Fälle, 0 Fehlschläge.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Zustandsmarker im Plan nach Ablauf nicht nachgezogen · Prüfzustands-Verb in Zusage-Menge · Zusage-Rand zum Container ohne Instanz offen · Doppelter Wortgebrauch „Instanz“ · Zielbild vor Code

## Verdikt

**Merge-blockierend:** nein — kein HIGH oder MEDIUM; kein Revert nötig. F-1 und F-2 gehören in einen Folge-Commit bzw. an den Projektinhaber (F-2 betrifft freigegebenen Wortlaut).

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Die DoD-Prüfung bleibt beim Verifier.
