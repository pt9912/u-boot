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
([`slice-harness-internal-readme-kennungs-retrofit`](../done/slice-harness-internal-readme-kennungs-retrofit.md))
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

- [ ] **Delta gesichtet, nicht nur gezogen:** Modul- und
  Konfigurations-Schema-Änderungen zwischen 0.2.0 und v0.51.1 sind erfasst
  (welche Module neu, welche Defaults geändert, welche Config-Schlüssel
  hinzugekommen). Keine Vollständigkeit über 50 Changelog-Einträge — die
  Frage ist: *ändert sich etwas an unserem Vertrag?*
- [ ] **Regressionsfreiheit belegt:** Trockenlauf des neuen Images gegen die
  **unveränderte** `.d-check.yml`. Abweichungen gegenüber dem alten Stand
  wären Befunde, keine Nebenwirkung.
- [ ] **Pin gebumpt:** `D_CHECK_IMAGE` in [`Makefile`](../../../../Makefile) auf
  den **Digest** von `v0.51.1` (nicht auf `:latest` — reproduzierbare Läufe,
  bestehende u-boot-Konvention).
- [ ] **`MR-005` nachgezogen:** Der Adaptions-Block nennt den Stand und die
  Bump-Prozedur für dieses Gate; heute steht dort nur „Modul-Auswahl bei
  d-check-Upgrade re-evaluieren" ohne Anhaltspunkt, wann ein Upgrade fällig
  ist.
- [ ] **Modul-Kandidaten benannt:** Die für u-boot einschlägigen neuen Module
  sind mit Nutzen und geschätzter Befundlast dokumentiert — als
  Entscheidungsgrundlage, nicht als Umsetzung.
- [ ] **Freshness-Lücke benannt:** Für die vendorte Regelwerk-Baseline gibt es
  seit heute einen Aktualitäts-Sensor, für den Gate-Image-Pin nicht. Der Slice
  hält fest, ob das ein Folge-Slice wird oder bewusst offen bleibt.
- [ ] `make docs-check` grün mit dem neuen Pin.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

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

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas: *harness / Konventionen* und *Harness-Tooling* — beide
**GF** nach [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration.
