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
(Abweisungen je Exit-Code `10`); Umsetzungs-Plan [`slice-v1-katalog-projekt-containerenv`](../open/slice-v1-katalog-projekt-containerenv.md)
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
([`cr-vorschlag-containerenv.md`](cr-vorschlag-containerenv.md)), vom
Projektinhaber freigegeben und wortgleich in die Verträge übernommen. **Kein
Produktcode und keine Implementierung in diesem Slice** (Muster der CR-Slices).

**Ausdrücklich NICHT in diesem Slice:** der Umsetzungs-Diff (Tests, Generator,
Validation — der Umsetzungs-Slice [`slice-v1-katalog-projekt-containerenv`](../open/slice-v1-katalog-projekt-containerenv.md)
trägt ihn, nachdem der CR geschlossen hat) · die Messung ENV-vs-`containerEnv`
(per Docker-Test, Umsetzungs-Slice — der CR trägt nur die Ebene-Entscheidung,
nicht den Mess-Weg) · die Umstellung des eigenen DevContainers (Umsetzungs-Slice,
Liefer-Punkt 3) · Handbuch/READMEs/CHANGELOG-Nutzertext (der Umsetzungs-Slice —
Muster der Werkzeugketten-Kette).

## 2. Definition of Done

- [ ] **Vorschlags-Datei:** [`cr-vorschlag-containerenv.md`](cr-vorschlag-containerenv.md)
  mit Alt/Neu-Wortlaut je Stelle (Muster [`cr-vorschlag-java-toolchain`](../done/cr-vorschlag-java-toolchain.md))
  — Stellen (a) Lastenheft-Absatz, (b) Historie-Zeile `0.32.0`, (c)
  [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)
  Pfad-Zeilen, (d) TECH-Unterabschnitt (Render/Merge/Validation, Ebene der
  Werte) — samt den vier Entscheidungen des Umsetzungs-Plans.
- [ ] **Freigabe:** Projektinhaber erteilt den Neu-Wortlaut (je Stelle oder
  mit benannten Änderungen); Abweichungen je ausgewiesen.
- [ ] **Vertrags-Übernahme wortgleich:** `spec/lastenheft.md` (Absatz,
  Historie `0.32.0`, Versionsfeld) und `spec/spezifikation.md`
  ([`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)-Zeilen,
  TECH-Unterabschnitt mit Zusage-Sätzen) — Abgleich Zeichen für Zeichen im
  §7; **erst nach der Freigabe** (`AGENTS.md` §Source-Precedence des
  Umsetzungs-Plans §4).
- [ ] `make gates` grün; Review (kein Self-Review) und Verification mit
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
| dieser Slice, [`roadmap.md`](../in-progress/roadmap.md) (die Zeile beim Start), `open/README.md` (Bestand-Zeile) | update | Lifecycle-Moves per `make slice-mv`; Zustand je Verzeichnis |
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

*(Wird bei der Closure gefüllt.)*

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