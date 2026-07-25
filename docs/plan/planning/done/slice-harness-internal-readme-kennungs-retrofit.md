# Slice Harness: `internal/`-READMEs in den Referenz-Scan holen

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** noch keiner Welle zugeordnet (Wartungs-Kandidat in [`roadmap.md`](../in-progress/roadmap.md) §Nächste Wellen) — Folge-Befund aus
[`slice-harness-sub-area-modus-audit`](slice-harness-sub-area-modus-audit.md)
§9.2.

**Bezug:** [`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell)
(Kennungs-/Link-Disziplin) und
[`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005--carveout-disziplin)
(die `scan.ignore`-Ausnahme ist ein Carveout und braucht einen Plan — dieser
hier). Berührt
[`ADR-0013`](../../adr/0013-dokumentationsreferenzmodell.md) lesend.

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

Die Code-Paket-READMEs unter `internal/**` aus dem toten Winkel holen:
`LH-*`-Kennungen verlinken, den zeitlichen Ballast (Meilenstein-/Tranchen-Tags
in Status-Abschnitten) entfernen und den `scan.ignore`-Glob `internal/**` in
[`.d-check.yml`](../../../../.d-check.yml) auflösen. Damit bekommt die
Inventur-Linie „README-Aussage vs. Code-Bestand" erstmals einen Sensor.

Der Befund stammt aus dem Sub-Area-Audit: Diese READMEs sind die einzige
**Brownfield**-Sub-Area des Repos (Doku beschreibt Bestand nachträglich), und
ihre Ausnahme vom Referenz-Gate ist bisher nur im Konfigurations-Kommentar
begründet — ohne Inventar-Eintrag, ohne Plan.

## 2. Definition of Done

- [x] **Kennungs-Inventar:** alle ungelinkten `LH-*`-/`ADR-*`-Kennungen in
  `internal/**/README.md` erfasst (Ist-Zahl im Slice festgehalten, nicht nur
  „einige").
- [x] **Links gesetzt** — relative Pfade auf
  [`spec/lastenheft.md`](../../../../spec/lastenheft.md) bzw. den ADR-Index,
  Anker-Form wie im übrigen Repo.
- [x] **Zeitliche Schicht entfernt:** Status-Abschnitte tragen keine
  Meilenstein-/Tranchen-Tags (`M3-T2`, `MVP-Closure-T1`, „Stand M8") mehr,
  sondern beschreiben den Ist-Zustand. Die zeitliche Schicht lebt in Roadmap
  und Closure-Notizen, nicht im Code-README.
- [x] **`scan.ignore` verengt:** `internal/**` entfällt; falls einzelne
  Restfälle bleiben müssen, wird der Glob so eng wie möglich gefasst und der
  Rest als Carveout mit Plan geführt — kein pauschaler Verzeichnis-Ausschluss.
- [x] **Carveout-Eintrag aufgelöst:** der Eintrag zum `internal/**`-Ausschluss
  im Inventar wird geschlossen (Audit-Trail bleibt).
- [x] `make docs-check` grün mit dem erweiterten Scope (Datei-Zahl steigt
  sichtbar — das ist die Evidence).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/**/README.md` (7 Dateien) | update | Kennungen verlinken, Status entzeitlichen |
| [`.d-check.yml`](../../../../.d-check.yml) `scan.ignore` | update | `internal/**` entfernen bzw. minimal fassen |
| [`carveouts.md`](../in-progress/carveouts.md) | update | Eintrag auflösen mit Audit-Trail |

## 4. Trigger

Kein externer Trigger; direkt einplanbar. Sinnvoll gebündelt mit der nächsten
Doku-/Harness-Welle, weil der Diff breit (7 Dateien), aber flach ist.

## 5. Closure-Trigger

`make docs-check` grün ohne `internal/**`-Ausschluss, Carveout-Eintrag
aufgelöst, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **Befund-Lawine:** Der erweiterte Scan kann mehr Befunde melden als erwartet
  (Anker-Drift in älteren READMEs). Gegenmaßnahme: erst Trockenlauf mit
  entferntem Ignore, Ist-Zahl notieren, dann entscheiden, ob der Slice in
  Tranchen zerfällt.
- **Entzeitlichen vs. Informationsverlust:** Die Meilenstein-Tags tragen
  historischen Kontext. Sie werden nicht ersatzlos gelöscht, sondern durch
  Ist-Aussagen ersetzt; die Historie steht ohnehin in den `done/`-Slices.
- **Carveout bleibt möglich:** Falls ein Restbestand nicht sinnvoll linkbar ist
  (z. B. Kennungen in Code-Blöcken), bleibt ein enger, begründeter Ausschluss —
  dann aber inventarisiert und mit Plan-Anker.

## 7. Closure-Notiz (nach `done/`)

### Ist-Zahlen

Trockenlauf mit entferntem Ignore, wie im Plan-Risiko „Befund-Lawine"
vorgesehen: **26 `id-unlinked`-Befunde** über sechs der sieben READMEs —
deutlich weniger als befürchtet, und keine Anker-Drift. Der Scan wuchs von
132 auf 139 Dateien.

Ein 27. Befund war ein Fehlalarm mit Signalwert: Die Wendung
„slice-mandatorisch" traf das `slice-`-Kennungsmuster, ohne eine Slice-Referenz
zu sein. Umformuliert statt unterdrückt — ein Gate, das man mit einer
Ausnahme beruhigt, prüft danach weniger.

**Wichtiger Nebenbefund zur Reichweite des Gates:** Es meldet nur *bare*
Kennungen. Steht eine ID in Backticks (`` `LH-FA-UP-004` ``), gilt sie als
Code-Beispiel und bleibt ungeprüft — so lag der Großteil der Kennungen in
diesen Dateien. Der Sensor deckt die Inventur-Linie also **nicht
vollständig** ab; er fängt die prosaischen Referenzen, nicht die
code-formatierten. Das ist eine Eigenschaft des Werkzeugs, keine Entscheidung
dieses Slice — festgehalten, damit die Abdeckung nicht überschätzt wird.

### Verification Evidence

Scope:
- Slice: `slice-harness-internal-readme-kennungs-retrofit`
- IDs: **keine** Anforderung geändert.
  [`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell)
  (Kennungs-/Link-Disziplin) und
  [`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005--carveout-disziplin)
  (Carveout aufgelöst) sind die tragenden Anker.
- Artefakte: sieben `internal/**/README.md`,
  [`.d-check.yml`](../../../../.d-check.yml) (`scan.ignore` verengt),
  [`carveouts.md`](../in-progress/carveouts.md) (Eintrag aufgelöst, Audit-Trail
  ergänzt).

DoD-Abgleich: alle Punkte erfüllt. Der `scan.ignore`-Glob ist **vollständig**
aufgelöst, nicht nur verengt — es blieb kein Restfall, der eine Ausnahme
gebraucht hätte. Über den Plan hinaus wurden die Inventar-Listen der READMEs
mit dem Bestand abgeglichen (fehlende Services und Subkommandos ergänzt): Ein
entzeitlichter Status-Abschnitt, der weiter „sieben Use-Cases" bei zwölf
behauptet, hätte die eine Alterungsform gegen die andere getauscht.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make docs-check` (Trockenlauf ohne Ignore) | 26 Befunde | Ist-Zahl als Ausgangspunkt, alle `id-unlinked` |
| `make docs-check` (nach Retrofit) | pass | **139 Dateien / 0 Befunde** — die Dateizahl ist die Evidence für den erweiterten Scope |

Traceability:
| ID / Pflicht | Beleg |
| --- | --- |
| [`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell) | 26 Kennungen verlinkt; `ids` grün bei erweitertem Scope |
| [`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005--carveout-disziplin) | Carveout aufgelöst statt verlängert; Audit-Trail im Inventar |
| `MR-006` Graduation der BF-Sub-Area | Inventur-Linie „README vs. Code-Bestand" hat jetzt einen Sensor → Modus-Wechsel BF → GF |

Carveouts: Neu: none. **Gelöst:** `internal/**` im `scan.ignore`.
Unverändert: none.

Nicht ausgeführt:
- `make gates` / `make test` — kein Go-Delta; ausschließlich README-Prosa und
  Gate-Konfiguration.

Independent Review: nicht durchgeführt. Der Sensor ist hier der Reviewer: Die
Aussage „alle Kennungen verlinkt" ist maschinell geprüft, und die inhaltlichen
Ergänzungen sind gegen dieselbe Code-Inventur belegt, die zwei Slices zuvor
`spec/architecture.md` §2 getragen hat.

Commit / Artefakt: `8cab0e2` (7 READMEs, `.d-check.yml`, Carveout-Aufloesung, `MR-006`-Graduation).

### Steering-Loop-Lerneintrag

- **Automatisches Entzeitlichen hat Prosa zerschossen — zweimal.** Der erste
  Durchlauf entfernte Meilenstein-Tags per Regex und ließ Satzfragmente zurück
  („Der `--devcontainer`-Flag: zwei Templates durchlaufen…"); ein
  Whitespace-Aufräumer zerlegte zusätzlich die ASCII-Baumstruktur in
  `internal/README.md`. Korrektur war `git checkout` und ein zweiter, engerer
  Durchlauf. Lehre: Ein Tag ist Teil eines Satzes, kein Token — Massenersetzung
  liefert die Kandidaten, das Umformulieren bleibt Handarbeit. Und generische
  Whitespace-Regeln haben in Dateien mit Layout (Bäume, Tabellen, Code-Blöcke)
  nichts verloren.
- **Der Rückzieher war billiger als die Reparatur.** Weil alle Änderungen
  scriptgestützt und unkommittet waren, kostete das Zurücksetzen nichts. Das
  spricht dafür, mechanische Massenläufe *vor* dem ersten Commit zu prüfen —
  danach wäre dieselbe Korrektur ein Revert-Commit mit Erklärungsbedarf gewesen.
- **Ein Gate mit Blindfleck ist kein Freibrief.** Dass code-formatierte
  Kennungen ungeprüft bleiben, wurde beim Verlinken sichtbar. Die ehrliche
  Konsequenz ist, die Abdeckung zu benennen, statt „`ids` grün" als
  Vollständigkeit zu lesen.
- **Folge-Slices:** keine. Sollte die Backtick-Lücke einmal stören, ist das ein
  d-check-Thema (Werkzeug), kein u-boot-Slice.

## 8. Sub-Area-Modus-Begründung

### Sub-Area: `internal/**/README.md` (Code-Paket-READMEs)

- **Modus:** BF (Doku beschreibt den Code-Bestand nachträglich).
- **Konventionen-Dichte:** niedrig — keine Form-Regel für diese READMEs
  vorhanden; dieser Slice zieht die erste ein (Kennungs-Linkpflicht,
  keine zeitliche Schicht).
- **Phase-Reife:** 3 (partiell) → Ziel 4 (kohärent, mit Sensor).
- **Evidenz-/Diskrepanz-Risiko:** hoch — genau das macht dieser Slice
  sichtbar; die Ist-Zahl der Befunde ist das erwartete Diskrepanz-Signal.
- **Reconciliation-Aufwand:** ein Slice, ggf. zwei Tranchen (Kennungen /
  Entzeitlichung). Graduation nach GF, sobald `internal/**` im Scan liegt.
