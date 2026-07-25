# Slice Gate: `--print-mk`-Fragment einbinden, `docs-check` als Alias

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** `welle-gate-ausbau-v0.51` (s. [`roadmap.md`](../in-progress/roadmap.md)).

**Bezug:** `MR-005` ([`harness/conventions.md`](../../../../harness/conventions.md))
— die Ablehnung des Fragments wird ersetzt. Kein `LH`-Bezug: Werkzeug-Wiring.

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

Das von d-check erzeugte `d-check.mk`-Fragment einbinden statt den
`docker run`-Aufruf selbst zu pflegen — und `docs-check` als dünnen Alias auf
`doc-check` erhalten, damit kein einziger bestehender Doku-Verweis bricht.

Der Gewinn ist konkret: `--network none` an jedem Target (ein Doku-Gate
braucht kein Netz), fertige Targets für alle opt-in-Module samt der jeweils
rund achtzehn Glieder langen `--disable`-Ketten, und der Image-Pin an genau
einer Stelle. Die `--disable`-Ketten sind das eigentliche Argument: Sie wachsen
mit jedem neuen d-check-Modul, und von Hand gepflegt sind sie eine
Drift-Quelle, die kein Sensor fängt.

## 2. Definition of Done

- [ ] **Fragment erzeugt und committet:** `d-check.mk` im Repo-Wurzelverzeichnis,
  erzeugt per `--print-mk` aus dem gepinnten Image; als generiertes Artefakt
  im Kopf kenntlich.
- [ ] **Pin über `DCHECK_DIGEST`:** Der Digest steht im `Makefile` (nicht im
  generierten Fragment), damit ein Re-Generieren den Pin nicht überschreibt.
- [ ] **`docs-check` bleibt der u-boot-Name:** als `.PHONY`-Alias auf
  `doc-check`. Kein Doku-Sweep — `AGENTS.md`, `harness/verification.md`, die
  CI-Workflows, `docs/user/quality.md` und die Gate-Nennungen in bestehenden
  `done/`-Closures bleiben unangetastet.
- [ ] **Alt-Aufruf entfernt:** `D_CHECK_IMAGE` und das handgeschriebene
  `docker run`-Recipe entfallen; es gibt genau einen Weg.
- [ ] **CI unverändert grün:** Die Workflows rufen `make docs-check` — der
  Alias muss dort ohne Anpassung durchlaufen.
- [ ] **`MR-005` umgeschrieben:** Von „kein Fragment" auf „Fragment + Alias",
  mit dem Grund der Neubewertung (Stand `v0.51.1` statt `0.2.0`) und dem
  Namens-Konflikt als bewusst gelöstem Punkt.
- [ ] **Re-Generierungs-Weg dokumentiert:** Wie das Fragment bei einem
  künftigen Image-Bump neu erzeugt wird, und dass dabei der Digest im
  `Makefile` bleibt.
- [ ] `make docs-check` grün.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `d-check.mk` | neu (generiert) | Tool-gepflegte Targets statt Handkopie |
| [`Makefile`](../../../../Makefile) | update | `include`, `DCHECK_DIGEST`, `docs-check`-Alias; Alt-Recipe raus |
| [`harness/conventions.md`](../../../../harness/conventions.md) `MR-005` | update | Adaption dreht sich um |

## 4. Trigger

Gefeuert: Entscheidung des Projektinhabers nach dem Image-Bump
([`slice-harness-dcheck-image-bump`](../done/slice-harness-dcheck-image-bump.md)).

## 5. Closure-Trigger

Fragment eingebunden, Alias trägt, CI grün, `MR-005` umgeschrieben,
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **Generiertes Artefakt im Repo:** `d-check.mk` ist Tool-Output und wird
  beim nächsten Bump neu erzeugt. Handänderungen daran wären stille Drift —
  der Kopf muss das ausdrücklich sagen.
- **`.PHONY`-Kollision:** Das Fragment definiert eigene `.PHONY`-Ziele; die
  bestehende `help`-Konvention (`## `-Kommentare) muss weiter greifen.
- **CI-Netz-Annahme:** `--network none` ist neu. Falls ein Workflow-Schritt
  bisher implizit Netz nutzte (er sollte nicht), fällt es hier auf.
- **Kein Carveout erwartet:** Wiring-Wechsel ohne Gate-Lockerung; der
  Befundsatz bleibt identisch.

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas: *harness / Konventionen* und *Harness-Tooling* — beide
**GF** nach [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration.
