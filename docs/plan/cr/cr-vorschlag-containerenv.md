# CR-Vorschlag: Konfig-Schlüssel `devcontainer.containerEnv` (Projekt-Umgebungswerte) — Alt/Neu-Wortlaut je Stelle

**Status:** freigegeben und übernommen (2026-11-03; Freigabe des Projektinhabers „wie vorgelegt” im Agenten-Gespräch, Vertrags-Übernahme committet `d62aad5e`).

**Abweichungen der Übernahme je Stelle** (Ausweis nach dem Review-Report F-2/F-3; je Bestätigung des Projektinhabers, 2026-11-03):

- **Stelle (c):** der Schlüssel trägt eine **eigene Wertetabellen-Zeile**; die zwei Übersteuerungs-Zeilen (`devcontainer.features.<feature>.version`, Pipeline-Schritt 6) bleiben unberührt — die vorgelegte Formulierung „beide Zeilen tragen den Namen” bezog sie semantisch unzutreffend ein. Bestätigt („eigene Zeile”).
- **Stelle (c, Nachtrag der Fix-Runde):** die Wertetabellen-Zeile sagt Schreiben nur über die Unterpfade `devcontainer.containerEnv.<name>`; der Basis-Pfad ist nicht schreibbar (kein Ganz-Map-Set). Bestätigt („nur Unterpfade”).
- **Stelle (b):** die übernommene Historie-Zeile folgt der vollen Tabellenform (vier Zellen: Datum, Anforderung, Inhalt, Freigabe); der Vorschlag nannte den Gehalt nur in einer Zelle. Keine Down-Referenz (Decken-Regel) — die im Vorschlag genannten Links entfallen im Vertrag.
- **Stelle (d):** im Vertrag heißt die Namens-Norm nur *gültiger Dev-Container-Umgebungsname*; die konkrete Prüffunktion, die der Vorschlag als Kandidat nannte (`driven.ValidateDevcontainerEnvName`, Fundstelle beim Umsetzungs-Schnitt), ist Umsetzungs-Fundstelle und kein Vertrags-Inhalt.

**Gegenstand:** die Lieferung zu DoD-Zeile 1 des CR-Slices
[`slice-cr-devcontainer-containerenv`](../done/slice-cr-devcontainer-containerenv.md).
Diese Datei ist kein Slice und trägt keinen Lebenszyklus; sie ändert nichts an
den Verträgen. Der *Neu*-Wortlaut der Stellen (a) bis (d) ist freizugeben; die
vier Entscheidungen des Umsetzungs-Plans
([`§4`](../done/slice-v1-katalog-projekt-containerenv.md)) tragen ihn.

**Grundlage:** Vorbefund des Umsetzungs-Plans §1 (die ENV-Ebene ist heute die
einzige Trägerin für Projekt-Werte; die zwei Go-Schlüssel tragen die einzigen
`containerEnv`-Werte des Generators, [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)
Abschnitt *Cache-Werte*; `BUILD_MODE`/`GOFLAGS`/`GOTOOLCHAIN` stehen in der
Historie ausdrücklich als „nicht belegt"). Die vier Entscheidungen sind
gebilligt (2026-11-03, „Ok dein Vorschlag dazu"): **Schlüsselname** —
`devcontainer.containerEnv`, Map Name → Wert (String), gelesen aus
`u-boot.yaml`; **Merge-Verhalten** — Abweisung mit Exit-Code `10`, wenn ein
Projekt-Wert `GOCACHE` oder `GOMODCACHE` setzt und der Toolchain-Schlüssel
aktiv ist (zwei Schreiber einer Stelle), sonst darf das Projekt sie setzen;
**Ebene** — nur `containerEnv` (Container-Prozesse), Dockerfile-`ENV` bleibt
dem Generator und dem handgeschriebenen Rest vorbehalten, gemessen per
Docker-Test im Umsetzungs-Slice; **Verwaltet-Form** — die Werte stehen
idempotent im verwalteten Teil, Regeneration schreibt neu, ein entfernter
Wert verschwindet.

---

## Stelle (a) — Lastenheft, Absatz der Zusage unter `LH-FA-DEV-001`

**Alt:** *(kein Absatz; der `containerEnv`-Träger ist der Cache-Absatz der
Schlüssel, [`SPEC-054`](../../../../spec/spezifikation.md) verfeinert ihn.)*

**Neu (Vorschlag, Sätze S1–S6 in der Lastenheft-Satzform):**

> **Projekt-Umgebungswerte.** `u-boot generate devcontainer` liest
> `devcontainer.containerEnv` (Map von Namen zu String-Werten) und trägt die
> Werte in das verwaltete `containerEnv`-Objekt der erzeugten
> `.devcontainer/devcontainer.json` — neben den Generator-Werten der
> Schlüssel, ohne Duplikate; die Erzeugung ist idempotent (eine Regeneration
> schreibt die Werte neu, ein entfernter Wert verschwindet). Der Name jedes
> Wertes ist ein gültiger Dev-Container-Umgebungsname, der Wert ist ein
> String; gegen beide setzt `u-boot config set` je fachlichen Fehler ab
> (Exit-Code `10`, nichts geschrieben). Setzt das Projekt mit aktivem
> Toolchain-Schlüssel (`go-toolchain` oder `golangci-lint`) einen der
> Cache-Werte (`GOCACHE`, `GOMODCACHE`), die der Schlüssel selbst trägt,
> beendet `u-boot generate devcontainer` denselben Fachfehler mit Exit-Code
> `10`, ehe nichts geschrieben wird; ohne den Schlüssel darf das Projekt
> beide setzen. Die Werte stehen nur in der `containerEnv`-Ebene der
> erzeugten Datei — der Dockerfile-Block und der handgeschriebene Rest tragen
> sie nicht. `u-boot generate devcontainer` meldet jede erzeugte
> Projekt-Wert-Zeile in der Ausweisung der Erzeugung.

## Stelle (b) — Lastenheft, Versionshistorie

**Alt:** *(letzte Zeile: 0.31.0 — der java-toolchain-CR.)*

**Neu:** `| 0.32.0 | Projekt-Umgebungswerte: `devcontainer.containerEnv` als Map des erzeugten DevContainers (Umsetzungs-Slice nach diesem CR) |`

## Stelle (c) — `SPEC-040` (Konfigurationspfade)

**Alt:** die zwei Pfad-Zeilen tragen die Schlüssel bis `...features.<key>.version` (JDK-Stand).

**Neu (Vorschlag):** beide Zeilen tragen den Namen `devcontainer.containerEnv`
neben die bestehenden; die Schreib-Pipeline (`u-boot config set`) und die
Anzeige decken die Map ab.

## Stelle (d) — Technik-Unterabschnitt (innerhalb `SPEC-040`)

**Alt:** *(kein Unterabschnitt.)* — **Neu (Vorschlag):** Unterabschnitt
*Projekt-Umgebungswerte* mit den Zusage-Sätzen zur Render-Form (verwalteter
Teil, idempotent, ohne Duplikate), dem Merge-Verhalten (Abweisung bei
Kollision mit aktiven Toolchain-Schlüsseln, Vorschlags-Form von Stelle (a)),
der Namens-Validation (`driven.ValidateDevcontainerEnvName`, Fundstelle beim
Umsetzungs-Schnitt nachzulesen) und der Beleg-Tabelle in Auftrags-Form (je
Bedingung ein Test-Auftrag; die Test-Namen liefert der Umsetzungs-Diff, Muster
[`SPEC-054`](../../../../spec/spezifikation.md) unten). Spec-Historie-Zeile.

---

**Offene Punkte an die Freigabe:** der Schlüsselname
(`devcontainer.containerEnv`), die Kollisions-Abweisung der Cache-Werte
(Frage 2), die Verwaltet-Form (Frage 4) und die TECH-Stelle innerhalb
[`SPEC-040`](../../../../spec/spezifikation.md) statt als eigener `SPEC-<NNN>`
(F6-Muster der java-Vorschlags). Bei Freigabe mit ändern trägt §7 des
CR-Slices je Abweichung; der Umsetzungs-Slice folgt.