# CO-001: CRLF-Edge-Case in `u-boot logs --json`

**Status:** Aktiv

**Datum angelegt:** 2026-06-07 (Pre-T8-Bestätigungsrunde LOW-2 des Slices
[`slice-v1-cli-json-dry-run-logs`](../planning/done/slice-v1-cli-json-dry-run-logs.md),
geschlossen 2026-06-07). **Letzte Prüfung:** 2026-11-02 (Lesung `splitLogLines`
am Fundort und des Vertrags-Absatzes, Migrationsslice [`slice-harness-carveouts-baseline-form`](../planning/done/slice-harness-carveouts-baseline-form.md)).

**Betroffenes Gate:** `make test` — Deckung der CLI-Vertrags-Tests; kein
Schwellen-Verweis: der Carveout nimmt keine Schwelle aus, er dokumentiert die
Grenze der getesteten Zusage.

**Geltungsbereich:** `internal/adapter/driving/cli/logs.go` — `splitLogLines`
(gelesen 2026-11-02, Zeilen 309 bis 320: TrimSuffix `"\n"`, Split auf `\n`).

**Folge-Slice:** [`slice-v1-logs-crlf`](../planning/open/slice-v1-logs-crlf.md)
(Anlage 2026-11-02; Start an seinem Trigger).

Regeln: Baseline-Regelwerk `modul-07-carveouts.md` §Ziel-Form: Carveout — ein
Carveout braucht immer einen Auflösungs-Trigger **und** einen Folge-Slice.

---

## Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-07-carveouts.md`
§Ziel-Form: Carveout — technische Begründung, keine
„noch nicht geschafft“-Aussagen.

`splitLogLines` splittet den gepufferten Compose-Output nur auf `\n`. Bei
einem CRLF-terminierten Compose-Output trüge jede `lines[i]` ein trailing
`\r`. Heute kein bekannter Konsument-Pfad für so einen Output — die
Compose-CLI emittiert auf allen gemessenen Plattformen `\n`; der Fall ist
hypothetisch, und ein Fix ohne Befund würde Verhalten in den Contract
schreiben, das kein Konsument verlangt (Entscheidung der Bestätigungsrunde:
dokumentierte Grenze statt speculative Zusage).

## Auflösungs-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-07-carveouts.md`
§Ziel-Form: Carveout — konkret und prüfbar. „Wenn Zeit ist“ ist kein Trigger.

**Erster bestätigter Befund** über ein trailing `\r` in der `--json`-Antwort
von `u-boot logs` — eine Konsumer-Meldung (Issue), ein Review-Finding oder
ein Testbefund mit dem Fall; der Start-Trigger des Folge-Slices (§4, Absatz
*Start*). Ohne Befund bleibt der Carveout die dokumentierte Grenze; er wird
bei jeder Welle-Closure geprüft (Carveout-Audit, Modul 7).

## Geltungs-Konfiguration

Keine Gate-Konfiguration — der Carveout ist keine Werkzeug-Ausnahme. Sein
Träger ist die Vertrags-Doku:

| Datei | Zeile/Section | Wert |
|---|---|---|
| [`docs/maintainer/cli-json-contract.md`](../../maintainer/cli-json-contract.md) | §6.8, Absatz „CRLF-Edge-Case“ | dokumentierte Grenze („Falls je relevant: Folge-Slice ergänzt explizites `\r\n`-Handling“) |

## Verifikation (nach Auflösung)

- [ ] `splitLogLines` normalisiert `\r` unmittelbar vor `\n` (Fix des
      Folge-Slices; Pin-Tests je Fall in `logs_acceptance_test.go`).
- [ ] `make gates` grün ohne Ausnahme (die Gate-Konfiguration kennt keine
      CO-Ausnahme; nichts zurückzunehmen).
- [ ] Die Datei wandert nach `docs/plan/carveouts/done/` (reiner `git mv`).
- [ ] Der Folge-Slice [`slice-v1-logs-crlf`](../planning/open/slice-v1-logs-crlf.md)
      ist geschlossen oder mit Kennung dokumentiert.
- [ ] Der Vertrags-Absatz (§6.8) trägt den Ist-Stand statt der Grenze
      (Zustandsmarker mit Ablaufpunkt, `AGENTS.md` §3.10).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-06-07 | Angelegt (Pre-T8-Bestätigungsrunde LOW-2) | [`slice-v1-cli-json-dry-run-logs`](../planning/done/slice-v1-cli-json-dry-run-logs.md) |
| 2026-11-02 | Geprüft, weiterhin gültig; aus dem Master-Inventar in die Baseline-Form überführt (Migrationsslice [`slice-harness-carveouts-baseline-form`](../planning/done/slice-harness-carveouts-baseline-form.md), [`MR-002`](../../../harness/conventions/done/MR-002-carveout-inventar.md) resolved) | [`slice-v1-logs-crlf`](../planning/open/slice-v1-logs-crlf.md) (Anlage) |