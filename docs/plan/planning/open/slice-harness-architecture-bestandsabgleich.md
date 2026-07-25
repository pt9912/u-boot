# Slice Harness: `spec/architecture.md` §2 gegen den Ist-Bestand abgleichen

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** ohne Welle (Spec-Wartung) — Befund aus
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

- [ ] **Driving-Ports vollständig:** §2.3 beschreibt alle vorhandenen
  Use-Case-Schnittstellen. Ist-Stand beim Befund: elf Port-Dateien, drei
  dokumentierte Use-Cases.
- [ ] **Driven-Ports vollständig:** §2.4 beschreibt alle vorhandenen
  Driven-Ports. Nicht dokumentiert waren unter anderem der Runtime-Engine-Port,
  der Netz-Probe-Port und die Template-Katalog-/Datei-Ports; der
  Runtime-Engine-Port stand fälschlich unter „Vorgesehene Erweiterungen".
- [ ] **„Vorgesehene Erweiterungen" bereinigt:** In §2.2, §2.3, §2.5 und §2.6
  enthält der Abschnitt nur noch tatsächlich Offenes. Was existiert, wandert in
  den Inhalt.
- [ ] **Application-Services vollständig:** §2.2 nennt die vorhandenen
  Use-Case-Services statt nur zwei.
- [ ] **Zielbild-Regel gewahrt:** keine Meilenstein-/Tranchen-Tags, keine
  Slice-Referenzen, keine Commit-Hashes — die Kopf-Hard-Rule der Datei gilt auch
  für §2.
- [ ] **Abgleich ist belegt:** Jede neue Aussage ist am Code prüfbar; für jede
  entfernte Aussage ist klar, warum sie nicht mehr gilt.
- [ ] `make docs-check` grün.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

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

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Area: *spec / Architektur-Sicht* — **GF** nach
[`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration. Der Abgleich läuft hier ausnahmsweise rückwärts
(Code → Doku); das ist die einmalige Schließung einer Alterungslücke, keine
Änderung der Sub-Area-Richtung.
