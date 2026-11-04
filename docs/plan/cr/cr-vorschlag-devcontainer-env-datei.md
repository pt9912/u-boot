# CR-Vorschlag: Umgebungsdatei des Devcontainers — Lastenheft 0.34.0 → 0.35.0, Alt/Neu je Stelle

**Status:** freigegeben wie vorgelegt (Projektinhaber 2026-11-04, Auftrag „bitte alle … umsetzen“ nach Vorlage der fünf Punkte mit Empfehlung; die Lesart als Freigabe hat der Hauptagent im Gespräch ausgewiesen) und übernommen (Lastenheft 0.35.0).

**Gegenstand:** die Lieferung zum ersten Liefer-Punkt der DoD des Slice [`slice-v1-devcontainer-env-datei`](../planning/in-progress/slice-v1-devcontainer-env-datei.md) (Vertrag). Diese Datei ist kein Slice und hat keinen eigenen Lebenszyklus; sie bleibt dauerhaft im CR-Ordner.

**Grundlage:** Entscheidungen des Projektinhabers 2026-11-04 (Sitzung): Zweck ist Docker `--env-file` ohne Umweg über `${localEnv:<NAME>}`; `generate devcontainer` legt die Datei an, wenn sie fehlt; `generate devcontainer` erzeugt zusätzlich eine Vorlage mit Platzhaltern. Messbasis: [`docs/messungen/slice-v1-devcontainer-env-datei`](../../messungen/slice-v1-devcontainer-env-datei/README.md) — ein Lauf je Fall, Docker Engine 29.9.3 (`linux/amd64`, verschachtelt), nur `docker run --env-file`; **nicht gemessen** ist der Start über `u-boot devcontainer up` und den Editor.

---

## Stelle (a) — Lastenheft, [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), erster Absatz (Zeile 504)

**Alt:** „… dürfen Git-Zugangsdaten und Zugangsdaten eines Agenten oder eines Forge-Werkzeugs (…) weder im Image, im Workspace-Volume noch in `u-boot.yaml` oder einer anderen erzeugten Datei stehen.“

**Neu (Vorschlag):**
> „… weder im Image, im Workspace-Volume noch in `u-boot.yaml` oder einer anderen von `u-boot` erzeugten Datei stehen. Ausgenommen ist die Datei `.devcontainer/devcontainer.env`: Sie führt der Benutzer; `u-boot` legt sie ohne Werte an, liest ihren Inhalt nicht und gibt keinen ihrer Werte aus (LH-FA-DEV-012).“ <!-- d-check:ignore (Kennung der neuen Anforderung; der Anker entsteht mit der Übernahme) -->

## Stelle (b) — Lastenheft, [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), vierter Absatz (Zeile 519)

**Alt:** „… sie stehen nie in einer Datei oder einem Image.“

**Neu (Vorschlag):**
> „… sie stehen nie in einer von `u-boot` erzeugten Datei oder einem Image. Ausgenommen ist die vom Benutzer geführte Datei `.devcontainer/devcontainer.env` (LH-FA-DEV-012); der Befehl liest sie nicht.“ <!-- d-check:ignore (Kennung der neuen Anforderung; der Anker entsteht mit der Übernahme) -->

## Stelle (c) — Lastenheft, neue Anforderung nach [`LH-FA-DEV-011`](../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)

**Neu (Vorschlag):**
> ### LH-FA-DEV-012 – Umgebungsdatei des Devcontainers <!-- d-check:ignore (Kennung der neuen Anforderung; der Anker entsteht mit der Übernahme) -->
>
> *Priorität: V2.* Der erzeugte Devcontainer bindet die Datei `.devcontainer/devcontainer.env` als Umgebungsdatei der Container-Engine ein; der Benutzer führt darin die Umgebungsvariablen seines Containers.
>
> - `u-boot generate devcontainer` trägt in beiden Profilen `--env-file=${localWorkspaceFolder}/.devcontainer/devcontainer.env` in `runArgs` der erzeugten Konfiguration ein.
> - Fehlt die Datei, legt der Befehl sie an: ein Kommentar-Kopf, keine Variablen. Eine vorhandene Datei ändert er nie. Er liest ihren Inhalt nicht und gibt keinen Wert aus. Er erzeugt außerdem `.devcontainer/devcontainer.env.example` als verwaltete Datei mit Platzhaltern (nie reale Werte).
> - Die Datei gehört in die `.gitignore` des Projekts; `u-boot init` trägt sie dort ein, die Beispieldatei bleibt getrackt. `u-boot doctor` weist mit `warn` darauf hin, wenn die Datei fehlt.
> - Format, gemessen am 2026-11-04 mit Docker Engine 29.9.3 (`linux/amd64`, verschachtelt, ein Lauf je Fall): Zeilen `NAME=WERT` und Kommentarzeilen mit `#`; Anführungszeichen sind Teil des Werts; ein Präfix `export` und eine fehlende Datei lassen den Start der Engine mit Exit-Code `125` scheitern; ein mit `-e` gesetzter Wert gewinnt gegen die Datei; ein relativer Pfad wird gegen das Arbeitsverzeichnis des Clients aufgelöst, deshalb steht der absolute Pfad über `${localWorkspaceFolder}` in der Konfiguration.
> - **Zugesagt** ist die erzeugte Konfiguration und die Anlage der Dateien, soweit Tests es belegen. **Nicht zugesagt:** der Start über `u-boot devcontainer up` und über den Editor (die Auflösung von `${localWorkspaceFolder}` dort ist nicht gemessen), andere Engines und Hosts als die gemessene, die Reihenfolge gegen `remoteEnv` und `containerEnv`, Dateien mit `CRLF` und mehrzeilige Werte.

## Stelle (d) — Lastenheft, Kopf und Historie

**Alt:** Version `0.34.0`; letzte Historie-Zeile `0.34.0` (2026-11-03).
**Neu (Vorschlag):** Version `0.35.0`; neue Zeile `0.35.0` | 2026-11-04 | Neue Anforderung LH-FA-DEV-012 (Umgebungsdatei des Devcontainers) samt Ausnahme der Zusage „Zugangsdaten stehen in keiner erzeugten Datei“ in LH-FA-DEV-009 und LH-FA-DEV-010 für die benutzergeführte Datei | Freigabe des Projektinhabers (Datum der Freigabe). <!-- d-check:ignore (Kennungen der neuen Anforderung; Anker entstehen mit der Übernahme) -->

---

**Offen zur Freigabe:**

1. **Die Ausnahme der Sicherheitszusage** (a) und (b): Zugangsdaten dürfen in der benutzergeführten, ignorierten Datei stehen. Gegenposition: die Zusage unverändert lassen und die Datei nur für Nicht-Zugangsdaten vorsehen — das ließe sich nicht durchsetzen, `u-boot` liest die Datei nicht.
2. **Priorität V2** statt V1.
3. **Umfang von `u-boot doctor`:** nur „fehlt“ (weil Git-Status und Ignorierung `u-boot` heute nicht abfragt — ungeprüft, ob der Git-Port das trägt). Der Entwurf des Slice nennt zusätzlich „getrackt oder nicht ignoriert“; das bräuchte eine Git-Abfrage.
4. **`.gitignore` nur über `u-boot init`:** für bestehende Projekte schreibt `generate devcontainer` die `.gitignore` nicht. Gegenposition: der Befehl ergänzt die Zeile selbst.
5. **Der Start über `u-boot devcontainer up`** ist nicht zugesagt, bis eine Messung ihn belegt. Gegenposition: die Zusage erst nach der Messung formulieren und den CR bis dahin zurückstellen.

**Regeln dieser Form:** Die Kennungen sind verlinkt, soweit ihr Anker besteht; der Neu-Wortlaut trägt nur Gemessenes und benennt die Messgrenzen.
