# Slice CR: CR-Schnitt für Projekt-Umgebungswerte — der Schlüssel `devcontainer.containerEnv`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv` (Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine).

**Welle:** ohne Welle (Restlinie der Feature-Katalog-Zeile nach dem Release
`v0.20.0`; Anlass: billigung des Vorschlags des Projektinhabers, 2026-11-03 —
die drei ENV-Werte des eigenen DevContainers in `u-boot.yaml`). Eintragstyp-
Adaption des CR-Musters [`slice-cr-devcontainer-java-toolchain`](../done/slice-cr-devcontainer-java-toolchain.md)
/ [`slice-cr-devcontainer-gepinnte-werkzeugkette`](../done/slice-cr-devcontainer-gepinnte-werkzeugkette.md).

**Bezug:** [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)
(erzeugte `devcontainer.json` — der Absatz der Zusage kommt dort), [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)
(Konfig-Schema der Dev-Container-Form — Pfad-Zeilen), [`LH-FA-CONF-005`](../../../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern)
(`config set`/Anzeige berühren den Schlüssel), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
(Abweisungen je Exit-Code `10`); Umsetzungs-Plan [`slice-v1-katalog-projekt-containerenv`](../done/slice-v1-katalog-projekt-containerenv.md)
(liegt in `open/`, vier Entscheidungen je §4 — der CR trägt sie als Freigabe-Form).

**Berührte Spec-Stellen** (gelesen am 2026-11-03 im Vorbefund des
Umsetzungs-Plans): [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)
(die zwei Pfad-Zeilen — Fundstellen beim Schnitt nachlesen; bei der
Übersteuerungs-Zeile 522 und 544 des JDK-Stands), [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)
Abschnitt *Cache-Werte* (bestehende `containerEnv`-Werte; der CR berührt die
Stelle nicht, der Merge-Fall kommt in die Umsetzung), [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)
(Namenslesung, `Config.Env`-Auflösung — unverändert). Der Verweis zeigt
**aufwärts**.

**Verantwortlich:** pt9912 (Implementer; Umsetzung und Freigabe-Vorlage durch den Hauptagenten im Auftrag)

**Autor:** pt9912 (Rolle Planner, im Auftrag des Hauptagenten). **Datum:** 2026-11-03.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der freigegebene CR-Wortlaut für den Konfig-Schlüssel
`devcontainer.containerEnv` (Projekt-Umgebungswerte des erzeugten
DevContainers, Map Name → Wert in [`u-boot.yaml`](../../../../spec/lastenheft.md)) — Lastenheft-Absatz je
[`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen),
[`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)-Zeilen
und die Beleg-Aufträge der Technik-Stelle — wird vorgelegt
([`cr-vorschlag-containerenv.md`](../../cr/cr-vorschlag-containerenv.md)), vom
Projektinhaber freigegeben und wortgleich in die Verträge übernommen. **Kein
Produktcode und keine Implementierung in diesem Slice** (Muster der CR-Slices).

**Ausdrücklich NICHT in diesem Slice:** der Umsetzungs-Diff (Tests, Generator,
Validation — der Umsetzungs-Slice [`slice-v1-katalog-projekt-containerenv`](../done/slice-v1-katalog-projekt-containerenv.md)
trägt ihn, nachdem der CR geschlossen hat) · die Messung ENV-vs-`containerEnv`
(per Docker-Test, Umsetzungs-Slice — der CR trägt nur die Ebene-Entscheidung,
nicht den Mess-Weg) · die Umstellung des eigenen DevContainers (Umsetzungs-Slice,
Liefer-Punkt 3) · Handbuch/READMEs/CHANGELOG-Nutzertext (der Umsetzungs-Slice —
Muster der Werkzeugketten-Kette).

## 2. Definition of Done

- [x] **Vorschlags-Datei:** [`cr-vorschlag-containerenv.md`](../../cr/cr-vorschlag-containerenv.md)
  mit Alt/Neu-Wortlaut je Stelle (Muster [`cr-vorschlag-java-toolchain`](../../cr/cr-vorschlag-java-toolchain.md))
  — Stellen (a) Lastenheft-Absatz, (b) Historie-Zeile `0.32.0`, (c)
  [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)
  Pfad-Zeilen, (d) TECH-Unterabschnitt (Render/Merge/Validation, Ebene der
  Werte) — samt den vier Entscheidungen des Umsetzungs-Plans.
- [x] **Freigabe:** Projektinhaber erteilt den Neu-Wortlaut (je Stelle oder
  mit benannten Änderungen); Abweichungen je ausgewiesen.
- [x] **Vertrags-Übernahme wortgleich:** `spec/lastenheft.md` (Absatz,
  Historie `0.32.0`, Versionsfeld) und `spec/spezifikation.md`
  ([`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)-Zeilen,
  TECH-Unterabschnitt mit Zusage-Sätzen) — Abgleich Zeichen für Zeichen im
  §7; **erst nach der Freigabe** (`AGENTS.md` §Source-Precedence des
  Umsetzungs-Plans §4).
- [x] `make gates` grün; Review (kein Self-Review) und Verification mit
  Reports unter `docs/reviews/`; Doku-Nutzertext bleibt ohne Kennungen (der
  CR-Diff berührt nur die Vertrags-Straten — Muster [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features));
  Closure-Notiz mit Steering-Loop-Lerneintrag; Register fortgeschrieben oder
  „keine Beobachtung"; jedes Risiko aus §6 mit Ausgang; die drei Paarungen
  getragen.

## 3. Plan (vor Code)

| Datei | Änderungs-Art | Begründung |
|---|---|---|
| `cr-vorschlag-containerenv.md` (neu) | neu | DoD-Zeile 1 — Alt/Neu-Wortlaut je Stelle, Entscheidungen, Freigabe-Form |
| `spec/lastenheft.md` | update (nach Freigabe) | Absatz der Zusage [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), Versionsfeld `0.32.0`, Historie-Zeile |
| `spec/spezifikation.md` | update (nach Freigabe) | [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)-Zeilen samt TECH-Unterabschnitt (Zusage-Sätze, Abweisungen; Beleg-Aufträge statt Test-Namen — die Namen liefert der Umsetzungs-Diff, Muster [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)); Spec-Historie-Zeile |
| dieser Slice, [`roadmap.md`](../in-progress/roadmap.md) (die Zeile beim Start und die Freigabe-Berichtigung nach der Übernahme), [`cr-vorschlag-containerenv.md`](../../cr/cr-vorschlag-containerenv.md) (Status- und Abweichungs-Zeilen je Review; **Move nach `done/` bei der Closure**, Muster [`cr-vorschlag-java-toolchain`](../../cr/cr-vorschlag-java-toolchain.md) — Review F-7), `open/README.md` (Bestand-Zeile — **ohne Diff in diesem CR**: die Zeile trägt den Start-Stand, die „geliefert"-Form folgt der Closure, Review F-4) | update | Lifecycle-Moves per `make slice-mv`; Zustand je Verzeichnis; Verweis-Nachzug je Move |
| `docs/reviews/` (Review-Report, ggf. Nachlauf-Review, Verify-Report) | neu | nach Vertrags-Diff: Nachlauf-Review über genau diesen Diff (Muster [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)); Reports der zwei Rollen |
| Produktcode, Templates, Tests, Handbuch, READMEs, `CHANGELOG.md` | **kein Diff in diesem CR** | der Umsetzungs-Slice trägt sie (Muster der Werkzeugketten-Kette) |

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): `Verantwortlich:` gesetzt; kein
anderer Slice in `in-progress/` (WIP-Limit 1, gelesen beim Start).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): der CR verlangt mehr als die Absatz-, Zeilen-
  und TECH-Stellen (etwa ein zweiter Schlüssel oder ein Render-Verhalten in
  weiteren Dateien) — Rück zur Zerlegung, der Architect trennt.
- `in-progress` → `open` (blockiert): der Projektinhaber zieht die Ebene-Frage
  (ENV-vs-`containerEnv`) zurück und verlangt eine Dockerfile-`ENV`-Form — dann
  braucht der CR eine Messung vor dem Wortlaut (Umsetzungs-Messung vorziehen).

## 5. Closure-Trigger

DoD vollständig; der Vertrags-Diff ist wortgleich zum freigegebenen Wortlaut
(je Stelle, Abgleich in §7); `make gates` und `make docs-check` grün; Reports
liegen vor; Closure-Notiz mit Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

1. **Der freigegebene Wortlaut geht nicht wortgleich ins Lastenheft** (Anker
   statt Datei-Links, Historie-Datum, Zeilenversatz). — **Ausgang:** bei
   Closure. Vorsorge: Abgleich Zeichen für Zeichen, Ergebnis in §7 (Muster
   der CR-Slices).
2. **Eine Zusage geht weiter als die Ebene-Frage trägt** („containerEnv
   reicht" steht aus dem Kopf im CR, nicht gemessen). — **Ausgang:** bei
   Closure. Vorsorge: der CR-Wortlaut trägt die Ebene nur als Entscheidung
   (Vorschlag, gebilligt), nicht als Verhaltens-Zusage; die Wirkungs-Messung
   bleibt Umsetzungs-Pflicht
   ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
3. **Eine Abweisung ohne Träger im Code** (Merge-Konflikt, Namens-Norm) —
   **Ausgang:** bei Closure. Vorsorge: der CR trägt die Abweisung als
   Contract-Text, die Umsetzung als Test (Beleg-Auftrag); das Muster ist
   [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)-Abweisungen
   ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)).
4. **Plan-Aussagen über Code, die am Fundort nicht gelesen sind** — **Ausgang:**
   bei Closure. Vorsorge: Fundstellen am Schnitt gelesen; §3.10-Abgleich nach
   jeder Fix-Runde
   ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md)).

## 7. Closure-Notiz

### Closure (Planner, 2026-11-03, nach Review-Fix `df6c40ce`, Verifier-Berichtigung `45833353` und Vertrags-Übernahme `d62aad5e`)

**Evidence:** Review-Report [`2026-11-03-slice-cr-devcontainer-containerenv`](../../../reviews/2026-11-03-slice-cr-devcontainer-containerenv.md) (2 MEDIUM/3 LOW/2 INFO — F-1 Zustandsmarker bei der Übernahme, F-2 Übertragungs-Lesart der Stelle (c), F-3 Historie-Zellen, F-4/F-5 §3- und Basis-Pfad-Klärung, F-6/F-7 INFO; je berichtigt, F-2/F-5 je Projektinhaber-Bestätigung) · Verify-Report [`…-verify`](../../../reviews/2026-11-03-slice-cr-devcontainer-containerenv-verify.md) (1 INFO/3 LOW — V-1 doppelter Artikel, V-2 Ausweis unterzählt Zellen, V-3 Stelle-(d)-Ausweis ohne Abgrenzung, V-4 zwei Verhaltens-Zusagen ohne Beleg-Auftrag; je berichtigt; DoD-Zeilen 1 bis 3 belegt erfüllt).

**Abgleich je Stelle §§2 der Freigabe:** (a) zeichenweise wortgleich (difflib-Abgleich des Verifiers) · (b) vollen vier Tabellen-Zellen statt der vorgeschlagenen Zelle — Ausweis · (c) eigene Wertetabellen-Zeile samt Nur-Unterpfade-Präzisierung statt Anhängen an die zwei Übersteuerungs-Zeilen — Ausweis samt zwei Bestätigungen · (d) Namens-Norm im Vertrag ohne die Prüffunktion — Ausweis.

**Risiken aus §6 — je Ausgang:**
1. *Der freigegebene Wortlaut geht nicht wortgleich ins Lastenheft* — **eingetreten (abgemildert):** die Stellen (b)–(d) tragen je bestätigte Abweichungen (Ausweis in der Vorschlags-Datei und hier); die Stelle (a) ist zeichenweise wortgleich.
2. *Eine Zusage geht weiter als die Ebene-Frage trägt* — **entfallen:** der CR trägt die Ebene als Entscheidung (nur `containerEnv`), die Wirkungs-Messung bleibt Pflicht des Umsetzungs-Slices; der Absatz formuliert keine Verhaltens-Zusage über die Erzeugung hinaus.
3. *Eine Abweisung ohne Träger im Code* — **entfallen für den CR:** die Abweisungen stehen als Vertrags-Sätze mit Beleg-Aufträgen in der Beleg-Tabelle von [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) (seit der Verifier-Runde auch die zwei Nachgetragenen); die Code-Träger liefert der Umsetzungs-Diff (V-4-Klasse).
4. *Plan-Aussagen über Code, die am Fundort nicht gelesen sind* — **eingetreten (abgemildert):** die Ausweis-Familie trat (V-1 bis V-3) und ist berichtigt; der §3.10-Abgleich des Ranges ist voll deckend (die zwei `open/`-Diffs trägt der eigene Slices-Commit `60f9b6c5`).

**Was ging anders als geplant:** (a) Die Übernahme brauchte zwei Bestätigungen des Projektinhabers, weil der Vorschlag an Stelle (c) die zwei Übersteuerungs-Zeilen semantisch unzutreffend nannte und die Basis-Pfad-Schreibform unbestimmt ließ (F-2, F-5) — die Freigabe „wie vorgelegt" galt je Stelle, und die Abweichungen sind je ausgewiesen. (b) Die Zustandsmarker-Klasse trat bei der Übernahme erneut auf (F-1: Roadmap-Zeile und Vorschlags-Status trugen „Freigabe steht aus" über den Vertrags-Commit hinaus) — dieselbe Träger-Aussage wie bei den Release-Cuts. (c) Die Historie-Zelle-Zahl (F-3, V-2) und zwei fehlende Beleg-Aufträge (V-4) kamen erst vom Verifier.

**Steering-Loop-Eintrag** (gezählt, nicht verkörpert; kein Eintrag trägt das Feld `liegt in`): *Lehre 1:* die Zustandsmarker-Variante „Abgelaufen bei der Übernahme eines Vertrags" ist von [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) getragen (dritter Träger-Kontext in Folge: Release-Cuts, CR-Übernahme). *Lehre 2:* die Abweichungs-Klasse (Freigabe-Wortlaut gegenüber dem Entwurf verdichtet, Abweichung nicht je Stelle ausgewiesen) ist von keiner bestehenden Klasse exakt getragen — **Klassenaufstellung (eigener `BEO-DOK/<slug>` oder Zitat-Gewichtung) liegt zur Entscheidung beim Projektinhaber**; die Findings F-2/F-3/V-2/V-3 sind in der Evidence-Datei dieses Vorgangs getragen. *Neuer Sensor:* keiner. *Benannte Spec-Lücke:* keine.

**Beobachtungs-Register:** *Evidence neu:* [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) ([`evidence/slice-cr-devcontainer-containerenv.md`](../observations/BEO-DOK/plan-aussage-ueber-code/evidence/slice-cr-devcontainer-containerenv.md): F-1, F-2/F-3, V-1 bis V-3; berichtigt `df6c40ce` und `45833353`) · [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) ([`evidence/slice-cr-devcontainer-containerenv.md`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/evidence/slice-cr-devcontainer-containerenv.md): V-4, zwei Beleg-Aufträge nachgetragen). *Klassenfrage:* Lehre 2 oben (Projektinhaber).

**Carveouts:** keiner — kein Gate, keine Schwelle und kein Code im Range. **Folge-Slice:** die Umsetzung [`slice-v1-katalog-projekt-containerenv`](../done/slice-v1-katalog-projekt-containerenv.md) trägt die Produkt-Träger (Konfig-Schema, Generator, Validation, Golden Cases, Docker-Test der Ebene, Umstellung des eigenen Devcontainers) und die Test-Namen der Beleg-Aufträge; ihr Start folgt der Closure (Reihenfolge-Beschluss 2026-11-03). **Weiter offen:** nichts aus der DoD.

**Drei Paarungen:** *Anker* — kein Steering-Loop-Eintrag trägt ein Zielort-Feld; nichts verkörpert. *Folge-Slice* — der Umsetzungs-Slice ist eine Datei in `open/` (Verzeichnis-Lesung); er übernimmt die Beleg-Aufträge als Test-Aufträge. *Register* — jede `BEO-*`-Kennung in §7 löst in ein Verzeichnis mit der Evidence-Datei dieser Closure (`ls`).

**Verifikation** (Läufe nach Closure-Inhalt, Move und Vorschlags-Nachzug; Docker-Modus): `make gates` Exit 0 am 2026-11-03 (`record-gates`-Stempel = HEAD nach dem Nachzug). Der Lauf gilt für den Stand vor dieser Zeile; die Zeile ändert nur Prosa.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `spec/` (DOK) und `docs/plan/` (DOK);
laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)
Greenfield. Keine Code-Sub-Area im CR-Diff.

**Vorgelagert — offene Beobachtungen sichten:** [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)
(Risiko 2), [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)
(Risiko 3), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md)
(Risiko 4). Evidenz-/Diskrepanz-Risiko: mittel (Wortlaut-Zusagen ohne Code im
Slice; die Ebene-Frage bleibt Messung des Umsetzungs-Slices).

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield.

**Belege (Closure 2026-11-03, je DoD-Zeile):** (1) [`cr-vorschlag-containerenv.md`](../../cr/cr-vorschlag-containerenv.md) trägt die vier Stellen samt den vier Entscheidungen und dem je-Stelle-Abweichungs-Ausweis. (2) Freigabe „wie vorgelegt" 2026-11-03 samt den zwei Bestätigungen der Übertragungs-Lesarten (Stelle (c) eigene Wertetabellen-Zeile; Nur-Unterpfade-Schreiben) — die vier Träger-Spuren stimmen überein (Vorschlags-Status, Roadmap, Fix-Commit-Message, Plan; Verifier-Zeile 2). (3) Stelle (a) ist zeichenweise wortgleich (difflib-Abgleich des Verifiers); Stellen (b)/(c)/(d) tragen die bestätigten Abweichungen; kein Produktcode-, kein Nutzertext- und kein Gate-Diff im Range. (4) [`make gates`] grün je Berichtigungs-Runde (Stempel `df6c40ce` und `45833353` = HEAD, Coverage 92,20 %, `docs-check` 0 Befunde); [Review-Report](../../../reviews/2026-11-03-slice-cr-devcontainer-containerenv.md) (2 MEDIUM/3 LOW/2 INFO, je berichtigt) und [Verify-Report](../../../reviews/2026-11-03-slice-cr-devcontainer-containerenv-verify.md) (1 INFO/3 LOW, je berichtigt) committet; Register, Risiko-Ausgänge und Paarungen in §7.
