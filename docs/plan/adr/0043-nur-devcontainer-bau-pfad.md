# ADR 0043: Nur Devcontainer als Bau-Pfad — `BUILD_MODE` entfällt als Modus-Entscheidung, die Gatter-Kette läuft im Devcontainer-Container, `test-docker` in der Sockelform

**Status:** Accepted

**Datum:** 2026-11-03

**Autor:** pt9912 (Rolle Implementer — Umsetzung des freigegebenen CR [`cr-vorschlag-spec-031-devcontainer-bau`](../cr/cr-vorschlag-spec-031-devcontainer-bau.md); die drei Entscheidungen und die test-docker-Platzierung trägt der Projektinhaber am 2026-11-03, Wortlaut wie vorgelegt)

**Bezug:** [`LH-FA-DEV-001`](../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-NFA-PORT-002`](../../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten), [ADR-0004](0004-ci-system.md); Messbasis [`docs/messungen/slice-v1-bau-nur-devcontainer/README.md`](../../messungen/slice-v1-bau-nur-devcontainer/README.md) (M1–M5 samt Grenzen)

**Schärft:** [`SPEC-031`](../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) (der einzig-gültige Bau-Pfad und die Träger-Ebenen je Gatter samt Beleg-Tabelle) und [`SPEC-008`](../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) (die zweimodale Lesart der Pflicht-Targets entfällt) — die Kennung nennen, wo das Zielelement eine trägt.

**Supersedes:** [ADR-0037](0037-bauen-im-devcontainer-lokaler-modus.md)

**Verhältnis:** ADR-0037 bleibt `Accepted` und inhaltlich unverändert; seine Stellen des Zwei-Modus-Modells sind hier abgelöst: (a.) zwei Modi mit `BUILD_MODE`-Schalter (Entscheidung 1 und 2), (b.) „Am Host und in der CI gilt der Docker-Modus“ (Entscheidung 1), (c.) „die CI zunächst unverändert“ (Entscheidung, Kopfannahme), (d.) Nachweis „im Modus der Umgebung“ (Entscheidung 6), (e.) `modus-pruefung` als Fitness-Function (Modul 4, Abschnitt Fitness Function). Nicht abgelöst und unverändert übernommen: die Pin-Plätze der Toolchain im Devcontainer-Dockerfile (Katalog-Schlüssel nach [ADR-0041](0041-gepinnte-werkzeugkette-hybrid-versionierung.md)), der Go-Cache im Workspace-Verzeichnis (Entscheidung 8 von ADR-0037), `PKG`/`RUN` über die Umgebung statt die Kommandozeile, und die Ein-Quelle der Stage-Befehle (`scripts/bau.sh`).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0037](0037-bauen-im-devcontainer-lokaler-modus.md) führte die zwei Modi ein, um das Bauen im Devcontainer zu beschleunigen, ohne Host und CI zu ändern. Der Vorgang des Mess-Slices hat vor dem Vertrags-Wortlaut die CI-Fähigkeit der Devcontainer-Basis gemessen (M1–M5, [`docs/messungen/slice-v1-bau-nur-devcontainer/README.md`](../../messungen/slice-v1-bau-nur-devcontainer/README.md)):

- **M1:** das Devcontainer-Image baut mit dem Kontext `.devcontainer` (`docker build -f .devcontainer/Dockerfile .devcontainer`); der Repo-Wurzel-Kontext scheitert an `COPY docker-init.sh` — die Kontext-Konvention ist eine Grenze der Messung. Der Netzbau is nicht belegt (cache-warm).
- **M2:** der verschachtelte Daemon startet im privilegierten Container und trägt die Container-Züge (`docker run busybox`, `make modus-pruefung` OK).
- **M3:** `make gates` läuft im Container Exit 0 (62,0 s kalt, 33,0 s warm, Coverage 92,20 %), `docs-check` mit dem d-check-Image im verschachtelten Daemon grün.
- **M4:** `make test-docker` im verschachtelten Daemon ist **nicht gleichwertig**: 12 Fehlschläge von 1771 Testfunktionen (1758 PASS), gegen null Fehlschläge im Docker-Modus am Host (Messung 2026-10-31). Die Ursachen zweier Meldungs-Klassen sind bis zur Meldungs-Stelle bestimmt, nicht geklärt.
- **M5:** die CI-Topologie ist analog reproduziert: der unprivilegierte Job-Container am Runner-Socket erzeugt mit `docker run --privileged` den Devcontainer-Container, dessen verschachtelter Daemon Container-Züge der Ebene 3 trägt; ohne `--privileged` fällt der Start auf fallback zurück.

**Stand im Code (gelesen am 2026-11-03, mit Fundstelle).**

- `repo.mk` Zeilen 80 bis 102 und 130 bis 216: der `BUILD_MODE`-Zweig (`?=`-Default `docker`, `$(error)` für unbekannte Werte, der PKG/RUN-Lokal-Modus-Wächter) und die zwei Rezept-Zweige je Target; der lokale Zweig ruft `scripts/bau.sh`.
- `.github/workflows/ci.yml` Jobs `gates`/`security-gates`/`image-scan` und `integration.yml` Job `integration-docker`: die Runner rufen `make gates`, `make govulncheck`, `make build` + Trivy und `make test-docker` direkt im Docker-Modus (`BUILD_MODE`-Default `docker`).
- `modus-pruefung` (`repo.mk`, Ziel `modus-pruefung` mit den Kreuzproben über `make -n`): trägt die Zwei-Modus-Bedingungen je acht Targets samt PKG/RUN-Wächtern.

Damit stand der Docker-Modus als Standard im Workflow-Vertrag, während der Devcontainer nur der schnelleren Entwickler-Lauf diente; der Zug des Slices zielt darauf, den Devcontainer zur Bau-Umgebung zu erheben, statt die zwei Modi weiter zu pflegen.

## Entscheidung

Vom Projektinhaber am 2026-11-03 entschieden (drei Entscheidungen bejaht, test-docker-Platzierung Hybrid je Messbasis, CR freigegeben wie vorgelegt):

1. **CI im Devcontainer-Container.** Der `gates`-Job von `ci.yml` baut das Devcontainer-Image je Job-Beginn (`docker build -f .devcontainer/Dockerfile .devcontainer`, der Kontext je Messbasis M1) und ruft `make gates` in einem privilegierten `docker run` des Images auf (Muster je Messbasis M3; M5 belegt die Topologie: unprivilegierte Ebene 1 am Runner-Socket erzeugt die privilegierte Ebene 2). Die Sockelform-Züge bleiben am Runner-Docker: `test-docker` in `integration.yml` und die Container-Züge der `ci.yml`-Jobs `security-gates` und `image-scan`.
2. **Der Host-Vertrag bleibt.** Der Host braucht Docker und GNU `make` ([`LH-NFA-PORT-002`](../../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten)); die Nutzungs-Sicht verschiebt sich auf den Devcontainer — die Rezepte von `scripts/bau.sh` sind der einzige Pfad, die übrigen Gatter laufen im Container ([`SPEC-031`](../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows), [`AGENTS.md`](../../../AGENTS.md) §3.1). Kein Host-Go-Weg entsteht.
3. **Der `BUILD_MODE`-Zweig und `modus-pruefung` entfallen.** `repo.mk` trägt einen Weg (die `bau.sh`-Rezepte für die Container-freien Targets, die Docker-Formen für die Container-Züge); der Wert `local`, den die `containerEnv` des Devcontainers setzt, bleibt der erlaubte Modus-Wert, und der entfallene Schalter (`BUILD_MODE=docker`) endet mit dem make-Fehler der ungeöffneten Variable. `modus-pruefung` entfällt mit seiner Messbasis M2 (die Prüfung lief im verbleibenden Stand OK; die Frage, die sie trug — die zwei Modi — entfällt); seine Index-Zeilen (`harness/README.md` Sensors-Tabelle, `docs/maintainer/quality.md`) entfallen mit ihm. Die `PKG`/`RUN`-Wächter bleiben: der Teillauf-Schutz trägt ohne Modus weiter, nur die Modus-Sense der Wächter entfällt.
4. **Die Träger-Ebenen je Gatter** stehen als Beleg-Tabelle in [`SPEC-031`](../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) — der CI-Lauf je Job trägt die Läufe der Form nach, je Zeile ein Beleg der Messbasis M3/M4/M5 oder der CI-Jobs.

Die Zusage des verschachtelten Daemon für die docker-getaggten Tests entfällt damit endgültig aus dem Vertrag: [ADR-0037](0037-bauen-im-devcontainer-lokaler-modus.md) hatte sie schon im Docker-Modus nicht zugesagt; die Sockelform ist die bewusste Hybrid-Platzierung je Messbasis M4.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Status quo: zwei Modi ([ADR-0037](0037-bauen-im-devcontainer-lokaler-modus.md)), CI unverändert | kleinster Schritt; CI-Nachweis unverändert | der Docker-Modus-Zweig trägt ein zweites Nachweis-Umfeld ohne Zweit-Nachweis; `modus-pruefung` und seine Kreuzproben bleiben dauerhaft; die Messbasis M3 zeigt die Gatter-Kette im Container ohne Vorteil der Docker-Stage-Ebene (kalt 62,0 s warm 33,0 s, kein Kontext-Transfer) |
| B — der volle Zug auch mit `test-docker` im Container | ein Umfeld für alles inklusive der Integrationstests | die Messbasis M4 blockiert: 12 Fehlschläge im verschachtelten Daemon gegen null im Docker-Modus am Host — der Migration-Zug verliert die einzige Form der Nachweis-Gleichheit; die Ursachen von zwei Meldungs-Klassen sind ungeklärt |
| C — Toolchain-Image ohne Devcontainer (ADR-0037 Option E, `docker run` mit Bind-Mount und Cache-Volumes) | ein Modus für Host, Devcontainer und CI ohne den Devcontainer-Vertrag | größerer Umbau; die Devcontainer-Basis der Messung trägt den Zug nicht: M5 belegt die Topologie am echten Devcontainer-Image; Entscheidung des Projektinhabers: später |
| D — nur `repo.mk` einfädern, CI bleibt im Docker-Modus | kein CI-Risiko | die CI ruft dann einen Pfad auf, den `repo.mk` nicht mehr trägt (der Docker-Mode-Zweig ist entfernt); ein toter Modus-Schalter in CI bleibt, und der Nachweis der CI träfe einen nicht mehr existenten Vertrag |
| **E — die drei Entscheidungen (gewählt)** | ein Nachweis statt zwei (ADR-0037 Entscheidung 6, übernommen); der Modus-Gewirr entfällt; die Messbasis M1–M5 trägt jede Ebene je Punkt | der Netzbau des Images ist nicht belegt (Grenze M1); die CI-Runner sind nicht selbst gemessen (Grenze der Messbasis) — der erste echte CI-Lauf am Push trägt die Belege; `test-docker` bleibt hybrid platziert |

## Konsequenzen

- **Positiv:** ein Bau-Pfad, ein Nachweis (`make gates` im Devcontainer-Container); die Modus-Mechanik (`BUILD_MODE`-Zweig, `modus-pruefung` samt Kreuzproben) entfällt als Lauf- und Doku-Last; die Go-Caches bleiben im Workspace-Verzeichnis (ADR-0037 Entscheidung 8 übernommen); die CI-Basis ist dieselbe wie die Entwickler-Basis.
- **Negativ / Trade-offs:** der Devcontainer ist für die `gates`-Kette verpflichtend — ein bare-Host-`make gates` ohne Container endet rot (das ist die Zusage, nicht ein Defekt); der Image-Netzbau der CI ist nicht belegt (Grenze M1) und beim ersten echten Lauf zu belegen; zwei der sechs M4-Meldungs-Klassen sind nur bis zur Meldungs-Stelle bestimmt — der Hybrid-Weg umgeht die Frage, eine Berichtigung ist nicht zugesagt (Grenze der Messbasis).
- **Folgepflicht:** [`SPEC-031`](../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) samt Beleg-Tabelle, [`SPEC-008`](../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) (Zweck-Spalte), [`AGENTS.md`](../../../AGENTS.md) §3.1, `harness/README.md` (Sensors-Zeile und Lauf-Ebene), `docs/maintainer/quality.md`, READMEs und `CHANGELOG.md` nachziehen; die Workflow-Runs am Push tragen die Run-IDs im Slice-Bericht.
- **Die Übergangsform der Ablösung.** Die Status-Zeile des Vorgänger-ADRs trägt weiter `Accepted`; die Ablösung tragen die `Supersedes`-Kante dieses ADRs, die `Geschichte`-Zeile des Vorgängers und seine Index-Zeile. Grund: der gemessene Zustand der Gate-Konfiguration — das Flipping auf `Superseded by` in der Status-Zeile macht `[`matrix`]` jeden verlinkten Verweis auf den Vorgänger zum Befund (gemessen 2026-11-03 gegen d-check 0.91.0: 28 `matrix-inactive`-Befunde über Closing-Archive, Beobachtungs-Evidence und unveränderliche Bestands-ADRs), während die Entfernung der unverlinkten Kennung je `ids` zum Befund führt — das `matrix`-Modul und das `ids`-Modul bilden zusammen eine Referenz-Falle für Archiv-Quellen. Diese Stelle ist ein ehrlicher Messbefund, keine Lockerung: kein Gate (Statuslisten, Schwellen, Richtung) wurde verändert. Eine Hebung der d-check-Pin-Stelle mit einer Lineage-Ausnahme über den Bestand hinaus wäre eine eigene Entscheidung.

## Fitness Function

| Tooling | Regel | Make-Target |
|---|---|---|
| `make` | ein Aufruf mit dem entfallenen Schalter (`BUILD_MODE=docker …`) endet mit dem make-Fehler der ungeöffneten Variable; ohne den Schalter steht das `bau.sh`-Rezept | `BUILD_MODE=docker make -n compile` (rot) gegen `make -n compile` (Rezept ohne `docker build`-Stage-Aufruf) |
| `make docs-check` | ein dokumentiertes Target ohne Makefile-Regel meldet `gate-phantom`: die entfernte `modus-pruefung`-Zeile in `harness/README.md` bricht den Lauf rot, bis beide Seiten entfernt sind; die Pin-Gruppen stehen unverändert | `make docs-check` |
| `ci.yml` (Form) | der `gates`-Job baut das Devcontainer-Image und ruft `make gates` im Container; der erste CI-Lauf am Push trägt den Netzbau-/Lauf-Beleg | der CI-Lauf je Job (`gh run watch`-Form), M5-Analogie der [`roh-analogie.txt`](../../messungen/slice-v1-bau-nur-devcontainer/roh-analogie.txt) |

## Re-Evaluierungs-Trigger

- Steht die M4-Grenze nicht mehr (die sechs Meldungs-Klassen der verschachtelten-Daemon-Läufe sind bearbeitet oder die Test-Umgebungsannahme geändert), rückt `test-docker` in die Nachweis-Kette des Devcontainer-Containers — dann eine Folge-ADR mit `Supersedes ADR-0043`.
- Scheitert der Netzbau des Images in der CI an der Grenze M1 (Cache-Warm-Beleg nicht tragend), oder scheiterte ein CI-Job am Punkt, der die Docker-Stage-Träger trug, ist die Rückführung des Slices je Trigger (`in-progress` → `open`, Slice §4) und die Neuentcheidung fällig.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-11-03 | Entwurf, `Proposed` (Implementer), als Folge-ADR des CR [`cr-vorschlag-spec-031-devcontainer-bau`](../cr/cr-vorschlag-spec-031-devcontainer-bau.md) | CR, Messbasis M1–M5 |
| 2026-11-03 | `Accepted` durch den Projektinhaber — CR freigegeben wie vorgelegt, die drei Entscheidungen bejaht, test-docker-Platzierung Hybrid je Messbasis | CR (Status-Zeile), Messbasis (Bezug-Zeile) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).