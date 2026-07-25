# Slice Harness: `spec/architecture.md` §2 gegen den Ist-Bestand abgleichen

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** noch keiner Welle zugeordnet (Wartungs-Kandidat in [`roadmap.md`](../in-progress/roadmap.md) §Nächste Wellen) — Befund aus
[`slice-harness-architecture-template-konformitaet`](slice-harness-architecture-template-konformitaet.md).

**Bezug:** [`LH-FA-ARCH-002`](../../../../spec/lastenheft.md#lh-fa-arch-002--schichten-und-verzeichnislayout)
(Schichten und Verzeichnislayout) — die Sicht-Spec beschreibt den Bestand
unvollständig. Kein Vertrags-Delta: Das Lastenheft ist korrekt, die **Sicht**
ist veraltet.

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

`spec/architecture.md` §2 (Schichten und Verzeichnisse) auf den Ist-Bestand
heben. Der Abschnitt beschreibt in weiten Teilen den Stand der frühen
Meilensteine: Er führt unter „Vorgesehene Erweiterungen" auf, was längst
implementiert ist, und listet nur einen Teil der existierenden Ports.

## 2. Definition of Done

- [x] **Driving-Ports vollständig:** §2.3 beschreibt alle vorhandenen
  Use-Case-Schnittstellen. Ist-Stand beim Befund: elf Port-Dateien, drei
  dokumentierte Use-Cases.
- [x] **Driven-Ports vollständig:** §2.4 beschreibt alle vorhandenen
  Driven-Ports. Nicht dokumentiert waren unter anderem der Runtime-Engine-Port,
  der Netz-Probe-Port und die Template-Katalog-/Datei-Ports; der
  Runtime-Engine-Port stand fälschlich unter „Vorgesehene Erweiterungen".
- [x] **„Vorgesehene Erweiterungen" bereinigt:** In §2.2, §2.3, §2.5 und §2.6
  enthält der Abschnitt nur noch tatsächlich Offenes. Was existiert, wandert in
  den Inhalt.
- [x] **Application-Services vollständig:** §2.2 nennt die vorhandenen
  Use-Case-Services statt nur zwei.
- [x] **Zielbild-Regel gewahrt:** keine Meilenstein-/Tranchen-Tags, keine
  Slice-Referenzen, keine Commit-Hashes — die Kopf-Hard-Rule der Datei gilt auch
  für §2.
- [x] **Abgleich ist belegt:** Jede neue Aussage ist am Code prüfbar; für jede
  entfernte Aussage ist klar, warum sie nicht mehr gilt.
- [x] `make docs-check` grün.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`spec/architecture.md`](../../../../spec/architecture.md) §2.2–§2.6 | update | Ist-Bestand statt Frühstand |

Kein Code-Delta.

## 4. Trigger

Kein externer Trigger. Sinnvoll vor der nächsten größeren Architektur-Änderung,
weil eine veraltete Sicht-Spec jede Design-Review-Runde verteuert.

## 5. Closure-Trigger

§2 vollständig gegen den Bestand abgeglichen, `make docs-check` grün,
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **Umfang:** §2 ist der längste Abschnitt der Datei; der Abgleich betrifft
  fünf Unterabschnitte. Falls der Diff unübersichtlich wird, zerfällt der Slice
  in zwei Tranchen (Ports / Adapter).
- **Detailtiefe:** §2 trägt heute sehr viel Implementierungsdetail (Feldnamen,
  Funktionsnamen). Beim Abgleich ist zu entscheiden, ob diese Tiefe in einer
  *Sicht*-Spec richtig ist oder ob Teile in die Paket-Dokumentation gehören —
  die Entscheidung wird im Slice begründet, nicht nebenbei getroffen.
- **Doppelpflege:** Je detaillierter §2, desto schneller altert es wieder. Der
  Slice soll deshalb auch benennen, welche Detailtiefe künftig gepflegt wird.

## 7. Closure-Notiz (nach `done/`)

### Ist-Zahlen des Befunds

Was §2 vor diesem Slice beschrieb, gemessen am Bestand:

| Sektion | dokumentiert | vorhanden | Lücke |
|---|---|---|---|
| §2.2 Application-Services | 2 | 12 | 10 fehlten; drei davon standen als „Vorgesehene Erweiterungen", obwohl längst geliefert |
| §2.3 Driving-Ports | 3 | 11 | 8 fehlten; vier standen als „vorgesehen" |
| §2.4 Driven-Ports | 9 (+1 als „vorgesehen") | 14 | 4 fehlten ganz, 1 war als geplant geführt, obwohl implementiert |
| §2.5 CLI-Subkommandos | 1 als implementiert | 11 | 10 standen als „Vorgesehene Erweiterungen" |
| §2.6 Driven-Adapter-Pakete | 8 | 14 | 6 fehlten |

Der Abschnitt beschrieb damit im Kern den Stand nach den ersten beiden
Meilensteinen — vier Release-Wellen bevor er zuletzt angefasst wurde.

### Entscheidung zur Detailtiefe

Die im Plan offen gelassene Frage („ist diese Tiefe in einer *Sicht*-Spec
richtig?") ist so beantwortet:

**Die Sicht nennt Existenz, Verantwortung und Vertragsanker je Element sowie die
Regeln, die seine Änderung binden. Sie nennt keine Feldlisten, Funktionsnamen
oder Algorithmus-Schritte.**

Begründung: Feld- und Funktionsdetails altern mit jedem Refactoring, ohne dass
ein Gate ausschlägt — genau das Muster, das diesen Slice ausgelöst hat. Die
Struktur-Ebene dagegen ist per `depguard` abgesichert und driftet nicht
lautlos. Konkrete Folgen im Diff: Die Request-/Response-Feldlisten in §2.3 und
die Methodenlisten der Backup-/ManagedBlock-Helfer sind zu
Verantwortungs-Aussagen verdichtet; die Sentinel-Aufzählung ist durch die drei
**Familien** nach Exit-Klasse ersetzt, mit ausdrücklichem Verweis, dass die
vollständige Liste im Code steht und mit jedem Use-Case wächst.

Was **bleibt**, obwohl es detailliert ist: Aussagen, die eine Regel tragen —
der Port-Pin im Produktivcode, die Context-Konvention, die Trennung von
Probe und Engine, der `contextcheck`-Carveout. Sie sind keine Beschreibung des
Codes, sondern Bedingungen an ihn.

### Verification Evidence

Scope:
- Slice: `slice-harness-architecture-bestandsabgleich`
- IDs: **keine** Anforderung geändert.
  [`LH-FA-ARCH-002`](../../../../spec/lastenheft.md#lh-fa-arch-002--schichten-und-verzeichnislayout)
  lesend; die neu verlinkten `LH-*`-Anker (Add-on-, Lifecycle-, Generator-,
  Konfigurations- und Vorlagen-Familien) sind sämtlich Bestand.
- Artefakte: [`spec/architecture.md`](../../../../spec/architecture.md)
  §2.1–§2.7.

DoD-Abgleich: alle Punkte erfüllt. Ergänzend zum Plan wurden **§2.1** (Domäne)
und **§2.7** (Wiring) mit abgeglichen — §2.1 führte sechs Domänentypen als
„vorgesehen", die es gibt, und §2.7 nannte zwei von zwölf Services namentlich.
Beide waren im Plan nicht genannt, gehören aber zum selben Befund; die
Alternative wäre ein zweiter Slice für denselben Abschnitt gewesen.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make docs-check` | pass | 132 Dateien / 0 Befunde; `anchors` bestätigt jeden der neu gesetzten `LH-*`-Anker, `matrix` die unveränderte Aufwärts-Richtung |
| Code-Inventar als Gegenprobe | pass | Interface-, Service-, Adapter- und Subkommando-Listen direkt aus `internal/` erhoben (Zahlen oben) |

Traceability:
| ID / Pflicht | Beleg |
| --- | --- |
| [`LH-FA-ARCH-002`](../../../../spec/lastenheft.md#lh-fa-arch-002--schichten-und-verzeichnislayout) | §2 beschreibt jetzt das tatsächliche Schichten- und Verzeichnislayout |
| Kopf-Hard-Rule (keine zeitliche Schicht) | keine Meilenstein-Tags, keine Slice-Referenzen, keine Commit-Angaben im neuen Text |
| [`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell) | `matrix` grün — die Sicht verweist ausschließlich aufwärts |

Carveouts: Neu: none. Gelöst: none. Unverändert: der `contextcheck`-Carveout
in §2.5 bleibt beschrieben (permanent, im Inventar geführt).

Nicht ausgeführt:
- `make gates` / `make lint` / `make test` — kein Go-Delta; der Code wurde
  ausschließlich lesend inventarisiert.

Independent Review: nicht durchgeführt. Der Diff ist eine Inventur mit
maschinell prüfbarer Gegenprobe (jede genannte Struktur existiert, jede
Kennung ist verlinkt und von `docs-check` bestätigt). Das Risiko liegt hier
nicht bei falschen Aussagen, sondern bei *ausgelassenen* — und dagegen hilft
die Zähl-Tabelle oben mehr als ein zweiter Leser.

Commit / Artefakt: `b49fe9f` (§2.1-§2.7 abgeglichen, Detailtiefe-Regel gesetzt).

### Steering-Loop-Lerneintrag

- **„Vorgesehene Erweiterungen" ist die gefährlichste Rubrik einer Spec.** Sie
  altert ohne Widerspruch: Wenn das Vorgesehene geliefert wird, sagt niemand der
  Spec Bescheid, und der Eintrag wird von „Plan" stillschweigend zu „Lüge". In
  §2 standen vier Rubriken dieser Art, alle veraltet. Konsequenz: Solche
  Einträge brauchen entweder einen Trigger (wer liefert, streicht) oder sie
  gehören in die Roadmap, wo Veralten sichtbar ist.
- **Der Befund war ein Nebenprodukt.** Aufgefallen ist die Lücke nicht bei einer
  Architektur-Arbeit, sondern beim Schreiben eines *anderen* Abschnitts
  derselben Datei. Eine Spec wird dort geprüft, wo jemand ohnehin schreibt —
  das spricht dafür, Nachbarabschnitte beim Anfassen mitzulesen, statt sie
  bewusst zu meiden.
- **Detailtiefe ist eine Regel, keine Geschmacksfrage.** Ohne explizite Regel
  wächst jede Sektion auf die Tiefe, die der jeweilige Autor gerade nützlich
  fand — hier von einer Zeile bis zu Feldlisten. Die jetzt formulierte Regel
  („Existenz, Verantwortung, Vertragsanker; keine Feld- oder Funktionsnamen")
  ist beim nächsten Anfassen prüfbar.
- **Folge-Slices:** keine.

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Area: *spec / Architektur-Sicht* — **GF** nach
[`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration. Der Abgleich läuft hier ausnahmsweise rückwärts
(Code → Doku); das ist die einmalige Schließung einer Alterungslücke, keine
Änderung der Sub-Area-Richtung.
