# Slice Harness: `d-check`-Gate-Image von 0.2.0 auf v0.51.1

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** noch keiner Welle zugeordnet (Wartungs-Kandidat in
[`roadmap.md`](../in-progress/roadmap.md) §Nächste Wellen).

**Bezug:** `MR-005` ([`harness/conventions.md`](../../../../harness/conventions.md),
Gate-Haltung `docs-check`) und
[`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell)
(das Referenzmodell, dessen Durchsetzung an diesem Gate hängt). Keine
Anforderungsänderung — Werkzeug-Pin.

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

Den digest-gepinnten `D_CHECK_IMAGE`-Stand von **0.2.0** auf **v0.51.1**
heben. Der Pin stammt aus der Migration am 2026-06-12 und ist seither nie
nachgezogen worden — 50 Releases Abstand.

Auslöser ist ein konkreter Befund: Der Retrofit der `internal/`-READMEs
([`slice-harness-internal-readme-kennungs-retrofit`](slice-harness-internal-readme-kennungs-retrofit.md))
hat festgestellt, dass `ids` nur *nackte* Kennungen prüft und
Inline-Code-Spans auslässt. Das wurde dort als Werkzeug-Grenze notiert — es
ist aber seit d-check **0.8.0** eine Konfigurations-Option
(`ids.patterns[].link-policy: always`), die das alte Image schlicht nicht
kennt (`field link-policy not found`).

Zweiter, größerer Punkt: Zwischen 0.2.0 und v0.51.1 sind Regelmodule
entstanden, die Disziplinen mechanisieren, die u-boot heute **nur per Hand**
einhält (Lifecycle-Konsistenz, Doku-↔-Makefile-Abgleich, Immutabilität von
`done/`-Slices und Accepted-ADRs). Dieser Slice **bumpt nur**; die Auswahl
neuer Module ist Folge-Arbeit — sonst wird aus einem Pin-Bump ein
Gate-Umbau.

## 2. Definition of Done

- [x] **Delta gesichtet, nicht nur gezogen:** Modul- und
  Konfigurations-Schema-Änderungen zwischen 0.2.0 und v0.51.1 sind erfasst
  (welche Module neu, welche Defaults geändert, welche Config-Schlüssel
  hinzugekommen). Keine Vollständigkeit über 50 Changelog-Einträge — die
  Frage ist: *ändert sich etwas an unserem Vertrag?*
- [x] **Regressionsfreiheit belegt:** Trockenlauf des neuen Images gegen die
  **unveränderte** `.d-check.yml`. Abweichungen gegenüber dem alten Stand
  wären Befunde, keine Nebenwirkung.
- [x] **Pin gebumpt:** `D_CHECK_IMAGE` in [`Makefile`](../../../../Makefile) auf
  den **Digest** von `v0.51.1` (nicht auf `:latest` — reproduzierbare Läufe,
  bestehende u-boot-Konvention).
- [x] **`MR-005` nachgezogen:** Der Adaptions-Block nennt den Stand und die
  Bump-Prozedur für dieses Gate; heute steht dort nur „Modul-Auswahl bei
  d-check-Upgrade re-evaluieren" ohne Anhaltspunkt, wann ein Upgrade fällig
  ist.
- [x] **Modul-Kandidaten benannt:** Die für u-boot einschlägigen neuen Module
  sind mit Nutzen und geschätzter Befundlast dokumentiert — als
  Entscheidungsgrundlage, nicht als Umsetzung.
- [x] **Freshness-Lücke benannt:** Für die vendorte Regelwerk-Baseline gibt es
  seit heute einen Aktualitäts-Sensor, für den Gate-Image-Pin nicht. Der Slice
  hält fest, ob das ein Folge-Slice wird oder bewusst offen bleibt.
- [x] `make docs-check` grün mit dem neuen Pin.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`Makefile`](../../../../Makefile) `D_CHECK_IMAGE` | update | Digest-Pin auf `v0.51.1` |
| [`harness/conventions.md`](../../../../harness/conventions.md) `MR-005` | update | Stand + Bump-Prozedur für das Gate-Image |
| dieser Slice §Delta | neu | Sichtung als Protokoll, nicht in `conventions.md` |

Kein Code-Delta, keine Änderung an `.d-check.yml` — die Modul-Auswahl bleibt
wie sie ist.

## 4. Trigger

Bereits gefeuert: Befund aus dem README-Retrofit plus der Versions-Abstand
selbst. Kein externer Trigger nötig.

## 5. Closure-Trigger

Pin gebumpt, Trockenlauf regressionsfrei, `MR-005` nachgezogen,
Modul-Kandidaten dokumentiert, `make docs-check` grün, Closure-Notiz
geschrieben.

## 6. Risiken und offene Punkte

- **Verschärfte Defaults:** Ein neueres Gate kann bei gleicher Konfiguration
  mehr finden (strengere Slug-Regeln, neue Grund-Codes). Deshalb ist der
  Trockenlauf gegen die unveränderte Config ein eigener DoD-Punkt und **vor**
  dem Pin-Wechsel auszuführen.
- **Config-Schema-Bruch:** Umgekehrter Fall — ein Schlüssel, den das neue
  Image nicht mehr kennt. Fällt im Trockenlauf sofort auf (Exit 2).
- **Scope-Drift:** Die neuen Module sind verlockend. Sie gehören **nicht** in
  diesen Slice; jedes eingeschaltete Modul bringt eigene Befundlast und eigene
  Ausnahme-Entscheidungen.
- **Kein Carveout erwartet:** Werkzeug-Pin ohne Gate-Lockerung. Sollte der
  Bump Befunde erzeugen, die nicht sofort behebbar sind, gilt
  Diagnose-vor-Carveout.

## 7. Closure-Notiz (nach `done/`)

### Delta 0.2.0 → v0.51.1

50 Releases. Die Sichtung folgt der Frage „ändert sich etwas an unserem
Vertrag?", nicht dem Anspruch auf Vollständigkeit.

**Für unsere Konfiguration: nichts.** Der Trockenlauf des neuen Images gegen
die **unveränderte** `.d-check.yml` liefert dasselbe Ergebnis wie der alte Pin
(140 Dateien, 0 Befunde). Kein Schema-Bruch, keine verschärften Defaults, die
uns treffen.

**Was hinzugekommen ist:**

| Ebene | 0.2.0 | v0.51.1 |
|---|---|---|
| Regelmodule | 4 (`links`, `anchors`, `ids`, `matrix`) | 20 |
| `ids`-Schema | Muster + Target | zusätzlich `link-policy` (`prose`/`always`) und `exempt-paths` je Muster |
| Betriebs-Modi | Befund-Zeilen | zusätzlich `--doctor` (gruppierte Diagnose mit Fix-Kandidaten), `--repair` (git-apply-fähiger Patch), `--json`/`--yaml`, `--trace` (Traceability-Matrix), `--suggest-config` |

Der `link-policy`-Schlüssel (seit 0.8.0) ist der unmittelbare Auslöser: Er
schließt genau die Lücke, die der README-Retrofit als „Werkzeug-Grenze"
notiert hatte. Diese Einordnung war falsch — es war eine
Konfigurations-Entscheidung, die wir nie treffen konnten, weil das Image sie
nicht kannte.

### Modul-Kandidaten für u-boot (Entscheidungsgrundlage, nicht Umsetzung)

Vier der sechzehn neuen Module mechanisieren Regeln, die u-boot heute **nur
per Hand** einhält:

| Modul | Was es prüft | Warum es hier passt |
|---|---|---|
| `planning` | Der Ruhe-Marker steht im `## Aktuelle Welle`-Block genau dann, wenn kein `slice-*` im Lifecycle-Verzeichnis liegt | Genau die Buchführung, die in dieser Session mehrfach von Hand nachgezogen wurde. Konfiguration: eine Zeile (`roadmap:`) |
| `targets` | Jedes in einer Doku-Tabelle behauptete `make X` ist eine Makefile-Regel — und jede Regel steht in der Autoritäts-Doku | [`AGENTS.md`](../../../../AGENTS.md) §Quality Gates ist genau so eine Tabelle; heute prüft sie niemand gegen das [`Makefile`](../../../../Makefile) |
| `immutable` / `vcs` | Der normalisierte Kern einer Datei ist seit dem Pinnen unverändert (hermetisch bzw. über eine Commit-Range) | u-boot hat zwei Immutabilitäts-Hard-Rules — Accepted-ADRs und `done/`-Slices — die bisher **kein** Sensor stützt |
| `codepaths` | Explizite Pfade in Inline-Code existieren | Die Slice- und Spec-Prosa ist voll von Datei-Pfaden; ein umbenanntes Paket bricht sie lautlos |

Bewusst **nicht** in diesem Slice eingeschaltet: Jedes Modul bringt eigene
Befundlast und eigene Ausnahme-Entscheidungen mit. Gemessen ist bisher nur
`link-policy: always` — **100 Befunde**, davon 78 in den `internal/`-READMEs,
13 im Review-Report und 15 in `done/`-Slices. Die letzten beiden Gruppen sind
der eigentliche Knackpunkt: Beide Artefaktklassen sind per Konvention
unveränderlich, ihre Kennungen nachträglich zu verlinken widerspräche der
eigenen Regel. `exempt-paths` ist dafür da — aber das ist eine Entscheidung,
kein Handgriff.

### Befund: `--print-mk` ist neu zu bewerten (`MR-005`)

`MR-005` hält fest, dass u-boot `docs-check` **direkt** per `docker run`
fährt und bewusst **kein** `--print-mk`-Fragment einbindet. Diese Entscheidung
fiel gegen den 0.2.0-Stand. Gegen v0.51.1 sieht die Rechnung anders aus — das
Fragment liefert heute:

- `--network none` an jedem Target (unser direkter Aufruf setzt das nicht:
  ein Doku-Gate braucht kein Netz, und die Härtung ist geschenkt),
- fertige Targets für die opt-in-Module samt der jeweils nötigen, sehr langen
  `--disable`-Ketten (die will niemand von Hand pflegen),
- den Image-Pin an **einer** Stelle, mit `DCHECK_DIGEST`-Override für
  reproduzierbare Läufe.

Dagegen steht ein Namens-Bruch: Das Fragment heißt `doc-check`, u-boot heißt
`docs-check` — und dieser Name steht in [`AGENTS.md`](../../../../AGENTS.md),
[`harness/verification.md`](../../../../harness/verification.md), in
Slice-Closures und in der CI. Ein Wechsel braucht entweder einen dünnen Alias
oder einen breiten Doku-Sweep. **Nicht in diesem Slice entschieden** —
festgehalten als der nächste anstehende Punkt.

### Befund: keine Aktualitäts-Routine für den Gate-Pin

Für die vendorte Regelwerk-Baseline gibt es seit heute `--check-freshness`.
Für den `D_CHECK_IMAGE`-Pin gibt es nichts — und genau dieser Pin ist
unbemerkt 50 Releases alt geworden. Dasselbe Alterungsmuster, zweimal am
selben Tag aufgetreten. Bleibt bewusst offen: Ob das eine eigene Routine
verdient oder mit der `--print-mk`-Entscheidung zusammenfällt (dort lebt der
Pin an einer Stelle), ist Teil derselben Folge-Entscheidung.

### Verification Evidence

Scope:
- Slice: `slice-harness-dcheck-image-bump`
- IDs: **keine** Anforderung geändert. `MR-005` (Gate-Haltung) nachgezogen.
- Artefakte: [`Makefile`](../../../../Makefile) (`D_CHECK_IMAGE`-Digest),
  [`harness/conventions.md`](../../../../harness/conventions.md) (`MR-005`).

DoD-Abgleich: alle Punkte erfüllt. Die beiden „benennen"-Punkte
(Modul-Kandidaten, Freshness-Lücke) sind oben ausgeführt; die
`--print-mk`-Neubewertung kam als zusätzlicher Befund hinzu.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| Trockenlauf neues Image, Config unverändert | pass | 140 Dateien / 0 Befunde — identisch zum alten Pin |
| Trockenlauf neues Image + `link-policy: always` | 100 Befunde | Messwert als Entscheidungsgrundlage, **nicht** übernommen |
| `make docs-check` (neuer Pin) | pass | 140 Dateien / 0 Befunde |

Traceability:
| ID / Pflicht | Beleg |
| --- | --- |
| `MR-005` Auflösungs-Trigger („Modul-Auswahl bei d-check-Upgrade re-evaluieren") | eingelöst: Delta gesichtet, Kandidaten benannt, Auswahl bewusst vertagt |
| Reproduzierbarkeit | Digest-Pin statt `:latest`; Digest von `v0.51.1` verifiziert (`:latest` zeigt derzeit auf denselben) |

Carveouts: Neu: none. Gelöst: none. Unverändert: none.

Nicht ausgeführt:
- `make gates` / `make ci` — kein Go-Delta; der Gate-Pin betrifft nur
  `docs-check`.
- Vollständige Changelog-Lektüre über 50 Releases — bewusst nicht: Die
  Regressionsfrage ist empirisch beantwortet (Trockenlauf), die
  Vertragsfrage über Modul- und Schema-Delta.

Independent Review: nicht durchgeführt. Der Diff ist eine Zeile; die
inhaltliche Leistung ist die Sichtung, und ihr Ergebnis ist mit zwei
Trockenläufen belegt.

Commit / Artefakt: `225627c` (Digest-Pin im `Makefile`, `MR-005` Stand + Neubewertungs-Auftrag).

### Steering-Loop-Lerneintrag

- **Ein Werkzeug-Pin ist eine Baseline wie jede andere.** Wir haben heute für
  die Regelwerk-Baseline einen Aktualitäts-Sensor gebaut — und im selben
  Atemzug übersehen, dass der Gate-Pin dieselbe Eigenschaft hat: Er altert
  ohne Ausschlag. 50 Releases sind kein Zufall, sondern die Abwesenheit eines
  Sensors.
- **„Werkzeug-Grenze" ist eine Diagnose, die man belegen muss.** Die
  `ids`-Lücke habe ich in der vorherigen Closure als Eigenschaft von d-check
  notiert, ohne in Konfiguration oder Quelle nachzusehen. Sie war eine
  abgeschaltete Option. Lehre: Bevor eine Einschränkung als „geht nicht"
  dokumentiert wird, gehört ein Blick in das Werkzeug — hier lag die Quelle
  sogar im Schwester-Repo bereit.
- **Eine Adaptions-Entscheidung altert mit ihrer Grundlage.** `MR-005` hat
  `--print-mk` gegen den 0.2.0-Stand abgelehnt. Gegen v0.51.1 fällt dieselbe
  Abwägung womöglich anders aus. Adaptions-Blöcke sollten deshalb nicht nur
  einen Auflösungs-Trigger tragen, sondern auch den *Stand*, gegen den sie
  entschieden wurden.
- **Folge-Entscheidungen:** (a) `--print-mk` einbinden oder direkt bleiben,
  (b) welche der vier Kandidaten-Module scharfgeschaltet werden, (c)
  `link-policy: always` samt `exempt-paths`-Politik, (d) Aktualitäts-Routine
  für den Gate-Pin.

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas: *harness / Konventionen* und *Harness-Tooling* — beide
**GF** nach [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration.
