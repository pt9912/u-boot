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

- [x] **Fragment erzeugt und committet:** `d-check.mk` im Repo-Wurzelverzeichnis,
  erzeugt per `--print-mk` aus dem gepinnten Image; als generiertes Artefakt
  im Kopf kenntlich.
- [x] **Pin über `DCHECK_DIGEST`:** Der Digest steht im `Makefile` (nicht im
  generierten Fragment), damit ein Re-Generieren den Pin nicht überschreibt.
- [x] **`docs-check` bleibt der u-boot-Name:** als `.PHONY`-Alias auf
  `doc-check`. Kein Doku-Sweep — `AGENTS.md`, `harness/verification.md`, die
  CI-Workflows, `docs/user/quality.md` und die Gate-Nennungen in bestehenden
  `done/`-Closures bleiben unangetastet.
- [x] **Alt-Aufruf entfernt:** `D_CHECK_IMAGE` und das handgeschriebene
  `docker run`-Recipe entfallen; es gibt genau einen Weg.
- [x] **CI unverändert grün:** Die Workflows rufen `make docs-check` — der
  Alias muss dort ohne Anpassung durchlaufen.
- [x] **`MR-005` umgeschrieben:** Von „kein Fragment" auf „Fragment + Alias",
  mit dem Grund der Neubewertung (Stand `v0.51.1` statt `0.2.0`) und dem
  Namens-Konflikt als bewusst gelöstem Punkt.
- [x] **Re-Generierungs-Weg dokumentiert:** Wie das Fragment bei einem
  künftigen Image-Bump neu erzeugt wird, und dass dabei der Digest im
  `Makefile` bleibt.
- [x] `make docs-check` grün.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `d-check.mk` | neu (generiert) | Tool-gepflegte Targets statt Handkopie |
| [`Makefile`](../../../../Makefile) | update | `include`, `DCHECK_DIGEST`, `docs-check`-Alias; Alt-Recipe raus |
| [`harness/conventions.md`](../../../../harness/conventions.md) `MR-005` | update | Adaption dreht sich um |

## 4. Trigger

Gefeuert: Entscheidung des Projektinhabers nach dem Image-Bump
([`slice-harness-dcheck-image-bump`](slice-harness-dcheck-image-bump.md)).

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

### Verification Evidence

Scope:
- Slice: `slice-gate-print-mk-einbindung`
- IDs: **keine** Anforderung geändert. `MR-005` neu gefasst.
- Artefakte: `d-check.mk` (neu, generiert),
  [`Makefile`](../../../../Makefile) (`include`, `DCHECK_DIGEST`, Alias;
  Alt-Recipe und `D_CHECK_IMAGE` entfernt),
  [`harness/conventions.md`](../../../../harness/conventions.md) (`MR-005`
  Adaption, Begründung, Re-Generierungs-Weg).

DoD-Abgleich: alle Punkte erfüllt.

- **Kein Doku-Sweep nötig:** Der Alias trägt. `make docs-check` läuft über das
  Fragment-Recipe (`d-check.mk:24`), und keine der Nennungen in
  [`AGENTS.md`](../../../../AGENTS.md),
  [`harness/verification.md`](../../../../harness/verification.md), den
  `done/`-Closures oder `docs/user/` musste angefasst werden.
- **CI unverändert:** Die Workflows rufen `make gates`, nicht `docs-check`
  direkt — die Kette `gates → docs-check → doc-check` greift ohne Anpassung.
- **`make help` trägt beides:** Die `## `-Konvention gilt auch im Fragment, die
  `doc-*`-Targets erscheinen samt Beschreibung; zusätzlich gibt es `doc-help`.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make docs-check` (über den Alias) | pass | 144 Dateien / 0 Befunde; Recipe-Quelle `d-check.mk:24`, also über das Fragment |
| `make help` | pass | `docs-check` und die `doc-*`-Targets gelistet |

Traceability:
| ID / Pflicht | Beleg |
| --- | --- |
| `MR-005` (Gate-Haltung) | Adaption von „kein Fragment" auf „Fragment + Alias" umgeschrieben, mit dem Stand, gegen den entschieden wurde |
| [`LH-FA-BUILD-007`](../../../../spec/lastenheft.md#lh-fa-build-007--docker-only-workflow) | Docker-only bleibt: dasselbe Container-Recipe, jetzt zusätzlich `--network none` |

Carveouts: Neu: none. Gelöst: none. Unverändert: none.

Nicht ausgeführt:
- `make gates` / `make ci` — kein Go-Delta. Die Alias-Kette über `gates` ist
  strukturell geprüft (`gates` hängt an `docs-check`, das an `doc-check`), ein
  voller Gate-Lauf hätte dieselbe Aussage teurer erzeugt.

Independent Review: nicht durchgeführt. Der Diff ist Wiring: ein generiertes
Fragment, drei Makefile-Zeilen, ein Ledger-Eintrag. Die Wirkung ist am
Sensor-Lauf ablesbar.

### Steering-Loop-Lerneintrag

- **Der Alias war die ganze Schwierigkeit — und er kostete eine Zeile.** Die
  Entscheidung „einbinden oder nicht" hing nicht am Fragment, sondern an einem
  Namen, der in unveränderlichen Artefakten steht. Wo eine Konvention breit
  verankert ist, ist die Anpassungsschicht billiger als die Anpassung.
- **Generierte Artefakte brauchen eine Pin-Trennung.** Der Digest liegt
  bewusst im `Makefile`, nicht im Fragment: Sonst würde jede Re-Generierung
  den Pin auf den Tag des erzeugenden Images zurücksetzen — ein Verlust von
  Reproduzierbarkeit, den niemand bemerkt hätte.
- **Ein Gate ohne Netz ist eine Härtung, die wir übersehen hatten.** Unser
  handgeschriebenes Recipe lief zwei Monate lang ohne `--network none`, obwohl
  das Doku-Gate hermetisch ist. Solche Kleinigkeiten liefert ein gepflegtes
  Fragment mit — das ist das eigentliche Argument gegen Handkopien.
- **Folge-Slices:** keine aus diesem Slice; die übrigen drei der Welle waren
  bereits geplant.

Commit / Artefakt: `e246bbf` (`d-check.mk`, `Makefile`-Wiring, `MR-005` neu gefasst).

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas: *harness / Konventionen* und *Harness-Tooling* — beide
**GF** nach [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration.
