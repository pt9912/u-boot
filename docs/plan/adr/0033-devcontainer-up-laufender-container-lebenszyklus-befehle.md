# ADR 0033: Dev-Container-CLI bei `up` gegen einen laufenden Container — das `postAttachCommand` läuft bei jedem `up`, das `postCreateCommand` und das `postStartCommand` nicht erneut; die Prämisse „kein Lebenszyklus-Befehl läuft“ aus ADR-0030 und ADR-0031 wird berichtigt, die Entscheidungen von ADR-0030 gelten

**Status:** Accepted

**Datum:** 2026-10-22

**Autor:** pt9912 (Rolle Architect)

**Bezug:** [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [ADR-0028](0028-devcontainer-start-ohne-editor-hilfs-image.md), [ADR-0029](0029-devcontainer-start-exit-11-und-socket-haltung.md), [ADR-0030](0030-devcontainer-up-failed-nach-zustand.md), [ADR-0031](0031-devcontainer-shell-ohne-editor.md)

**Schärft:** [`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (die Aussagen darüber, was die Dev-Container-CLI bei `up` gegen einen laufenden Container ausführt: *Ausgabe von `up`* und die Startphase von `shell`; die Spezifikation führt sie nach dem Code nach). Das Lastenheft sagt dazu nichts ([`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) gelesen); diese ADR ändert es nicht.

**Supersedes (teilweise):** vier Stellen, und nur diese. Aus [ADR-0030](0030-devcontainer-up-failed-nach-zustand.md): (1) im Kontext, Spiegelstrich „Verhalten der Dev-Container-CLI `0.89.0` bei `up` gegen einen laufenden Container“, der Satz „Kein Lebenszyklus-Befehl läuft.“; (2) in Entscheidung 2 der Satzteil „dass die CLI bei einem laufenden Container keine Lebenszyklus-Befehle ausführt“; (3) der erste Re-Evaluierungs-Trigger. Aus [ADR-0031](0031-devcontainer-shell-ohne-editor.md): (4) im Kontext, Absatz „Was die Frage entscheidet“, der Satzteil „und führt keine Lebenszyklus-Befehle aus ([ADR-0030](0030-devcontainer-up-failed-nach-zustand.md), Kontext)“. Für diese vier Stellen ist diese ADR die geltende Fassung (Wortlaut in Entscheidung 1; dieselbe Form wie [ADR-0030](0030-devcontainer-up-failed-nach-zustand.md) gegenüber ADR-0029 und [ADR-0032](0032-katalog-features-socket-und-lockerungen.md) gegenüber ADR-0029 und ADR-0031). Alles Übrige beider ADRs gilt unverändert, auch die Entscheidungen 1 bis 3 von ADR-0030 (Entscheidung 2 dieser ADR); beide bleiben `Accepted` und werden nicht geändert.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0030](0030-devcontainer-up-failed-nach-zustand.md) Kontext sagt für die Dev-Container-CLI `0.89.0` bei `up` gegen einen laufenden Container: „Kein Lebenszyklus-Befehl läuft“; Entscheidung 2 verlangt, dass das Handbuch es sagt; [ADR-0031](0031-devcontainer-shell-ohne-editor.md) übernimmt den Satz als Grund dafür, `up` bei jedem Aufruf von `shell` als Startphase zu nehmen. Beim Umsetzen von ADR-0030 hat der Docker-Test zu Entscheidung 2 (Fitness Function) neben dem `postStartCommand` auch das `postCreateCommand` und das `postAttachCommand` gezählt und gezeigt, dass das `postAttachCommand` bei jedem `up` läuft. Das ist der Wortlaut des ersten Re-Evaluierungs-Triggers von ADR-0030. Das Review des Diffs ([Review-Report](../../reviews/2026-10-22-slice-v1-devcontainer-up-failed-nach-zustand.md), F-2) hält fest, dass die Feststellung zu diesem Trigger beim Architect liegt und nicht bei der Implementation. **Eine neue Produktentscheidung steht nicht an;** die ADR prüft, was an dem widerlegten Satz hängt, und berichtigt ihn dort, wo zwei `Accepted`-ADRs ihn tragen.

**Gelesen am 2026-10-22** (jede Aussage über Code, Spezifikation oder Verhalten, `AGENTS.md` §3.10):

- `internal/e2e/devcontainer_start_docker_test.go`, `TestDevcontainerUp_Docker_AlreadyRunningLifecycleCounters`: ein von Hand geschriebenes Projekt (`dcHandWritten`), dessen `devcontainer.json` das `postCreateCommand`, das `postStartCommand` und das `postAttachCommand` in der Stringform trägt, je `echo run >> <Datei>` mit einer Datei je Befehl im Container; der Test liest die Zeilen der drei Dateien mit `wc -l` unmittelbar nach der Rückkehr von `up`. Der Test fordert (`counts`) und der Lauf ergab (Fix-Runde, PASS in 17,75 s, laut Bericht der Umsetzung; von mir nicht wiederholt):

  | Schritt (derselbe Container) | `postCreateCommand` | `postStartCommand` | `postAttachCommand` |
  |---|---|---|---|
  | erstes `up` (Container neu) | 1 | 1 | 1 |
  | zweites `up` (läuft; `--json`, `alreadyRunning: true`) | 1 | 1 | 2 |
  | drittes `up` (läuft; Textmodus) | 1 | 1 | 3 |
  | `down`, danach `up` (derselbe Container, `alreadyRunning: false`) | 1 | 2 | 4 |

  **Messgrenzen:** Dev-Container-CLI `0.89.0` (der Pin im Hilfs-Image), Linux `amd64`, lokale Engine, Docker 29.9.2, ein Lauf je Aufruf; die drei Befehle in der Stringform und gelingend.
- [`rohdaten-05-m4.txt`](../../messungen/slice-v1-sandbox-egress-neuladen-entwurf/rohdaten-05-m4.txt) (Colima, `arm64`, Messvariante, ein Lauf): Die Zeile „Running the postStartCommand“ steht in Schritt 0, 2 und 4 (Z. 4, 30, 55), in Schritt 3 (`up` gegen den laufenden Container, Z. 42–53) steht keine Zeile zu einem Lebenszyklus-Befehl. Die Messung hatte also nur ein `postStartCommand` zu zählen (Schluss aus der Ausgabe; die Konfiguration liegt den Rohdaten nicht bei). Die Läufe des Validators, die ADR-0030 nennt, betreffen ein `postStartCommand` mit Exit 3 und ein `postCreateCommand` (Strg-C).
- `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl`, Zeilen 27 bis 32: `postCreateCommand` (unter `if .PostCreate`) und `postStartCommand` (unter `if .PostStart`), beide im Zweig des Sandbox-Profils; `internal/hexagon/application/devcontainer_sandbox.go`, `devcontainerTemplateData`: das `postCreateCommand` trägt `u-boot-sandbox-init` (Nested Podman, schreibt die Datei `vfs`) und den Clone, das `postStartCommand` `u-boot-egress-init` und `u-boot-docker-init`. Die Suche nach `postAttach` in `internal` und `cmd` findet nur den Docker-Test; `initializeCommand` kommt außerhalb von Tests nur im Text der Ausweisung des Sockets vor (Kommentar, Meldung und Hilfetext), `onCreateCommand` und `updateContentCommand` gar nicht. **`u-boot` erzeugt kein `postAttachCommand`.**
- `internal/adapter/driven/docker/devcontainer_launcher.go`, `runArgs`: Der Hilfs-Container ruft `devcontainer up --workspace-folder <Ordner>` auf, bei `--recreate` mit `--remove-existing-container`; keine Option, die Lebenszyklus-Befehle auslässt.
- `internal/hexagon/application/devcontainer_start.go`, `failedAfterStartText`: die Meldung nach einem gelungenen Lauf der CLI sagt „the container is running and this call did not run the start script“ und nennt `u-boot devcontainer up --recreate`.
- [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): keine Aussage zu Lebenszyklus-Befehlen bei einem laufenden Container. [`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) enthält den Satz „keine Lebenszyklus-Befehle“ nicht (Suche nach „Lebenszyklus“ in `spec/spezifikation.md`: zwei Treffer, beide in anderem Zusammenhang). Handbuch (4.9 *Erneutes `up` nach einem Fehler oder Abbruch*) und `CHANGELOG.md` nennen die drei Befehle einzeln und mit den Messgrenzen.

**Befund.** (1) Die Verallgemeinerung „kein Lebenszyklus-Befehl“ stammt aus einer Messung, die nur ein `postStartCommand` kannte, und aus Läufen mit einem `postStartCommand` und einem `postCreateCommand`; keine der genannten Messungen nennt ein `postAttachCommand`. Sie ist für diese zwei Befehle belegt und für das `postAttachCommand` widerlegt. (2) Die Entscheidungen 1 bis 3 von ADR-0030 hängen nicht am `postAttachCommand` (Entscheidung 2 dieser ADR). (3) Die Folge, die der Trigger in Klammern nennt („Entscheidung 2 und der Schluss im Kontext fallen“), tritt nicht ein: Sie hängt daran, dass die Startscripte erneut laufen, und das zeigt der Test für das `postCreateCommand` und das `postStartCommand`, die sie tragen, nicht.

**Nicht gemessen, nicht zugesagt:** `initializeCommand`, `onCreateCommand` und `updateContentCommand` gegen einen laufenden Container; die Array- und die Objektform eines Befehls beim zweiten `up`; ein scheiternder Befehl beim zweiten `up` (die Läufe des Validators nennen ein `postStartCommand` mit Exit 3 und ein abgebrochenes `postCreateCommand`, je ein Lauf mit `u-boot` 0.13.0; ob ein gescheitertes oder abgebrochenes `postCreateCommand` nach `down` und `up` erneut läuft, bleibt unbelegt, wie in ADR-0030 Entscheidung 2, *Grenze*); andere Hosts als Linux `amd64`.

## Entscheidung

Jede Entscheidung nennt die Gegenoption, die verworfen ist.

1. **Die Prämisse lautet, wie gemessen.** Gegen einen laufenden Container führt die Dev-Container-CLI `0.89.0` bei `up` das `postCreateCommand` und das `postStartCommand` nicht erneut aus, das `postAttachCommand` dagegen bei jedem `up`; nach `down` und `up` im selben Container läuft das `postStartCommand` wieder, das `postCreateCommand` nicht (Messgrenzen im Kontext, nur gelingende Befehle). Die vier Stellen lauten künftig:

   | Stelle | geltende Fassung |
   |---|---|
   | ADR-0030, Kontext, Spiegelstrich zum Verhalten der CLI bei `up` gegen einen laufenden Container | „Das `postCreateCommand` und das `postStartCommand` laufen nicht erneut, das `postAttachCommand` läuft bei jedem `up`.“ Die Belege dort (`rohdaten-05-m4.txt` Schritt 3, der Validator) zeigen das `postStartCommand` und das `postCreateCommand`. |
   | ADR-0030, Entscheidung 2, Satz zum Handbuch | „Das Handbuch sagt dazu, was der Docker-Test zeigt: dass die CLI gegen einen laufenden Container das `postCreateCommand` und das `postStartCommand` nicht erneut ausführt, das `postAttachCommand` aber bei jedem `up`, und dass nach einem Fehler oder Abbruch `--recreate` hilft; der Satz nennt die Messgrenzen.“ |
   | ADR-0030, Re-Evaluierungs-Trigger, erster | „Der Docker-Test zeigt, dass die CLI bei `up` gegen einen laufenden Container das `postCreateCommand` oder das `postStartCommand` erneut ausführt (der Schluss im Kontext, die Meldung nach einem gelungenen Lauf und ihr Rat `--recreate` fallen).“ |
   | ADR-0031, Kontext, Absatz „Was die Frage entscheidet“ | „… und `up` gegen einen laufenden Container ist kurz (1,1 s und 1,2 s, je ein Lauf des Validators; ein `postAttachCommand` der Konfiguration kommt dazu) und führt das `postCreateCommand` und das `postStartCommand` nicht erneut aus; das `postAttachCommand` läuft bei jedem Aufruf (ADR-0033).“ |

   Unverändert bleibt in ADR-0030 der Satz im ersten Absatz des Kontexts („ohne die Lebenszyklus-Befehle zu wiederholen“): Er berichtet die Läufe des Validators und gilt für die Befehle, die sie nennen. **Gegenoption:** die Sätze stehen lassen und den Befund nur im Plan vermerken — der Plan hat im Rang weniger Gewicht als die ADR (Quellenrang), und ADR-0031 ist die Quelle, aus der die Pläne der Shell abgeleitet werden; ein Vermerk neben einem `Accepted`-Satz ließe den Satz gelten.

2. **Die Entscheidungen 1 bis 3 von ADR-0030 gelten unverändert; der Trigger ist im Wortlaut ausgelöst, seine Folge nicht.** Geprüft wurde jede Entscheidung gegen ihren Anker:
   - *Entscheidung 1* (ein `failed` in einer konfigurierten Statusdatei ergibt Exit 11, auch nach einem gelungenen Lauf der CLI) hängt an den Statusdateien. Sie schreiben `u-boot-sandbox-init` (aus dem `postCreateCommand`), `u-boot-egress-init` und `u-boot-docker-init` (aus dem `postStartCommand`); für genau diese zwei Befehle zeigt der Test (gelingende Befehle) und beobachtete der Validator (scheiternde Befehle, je ein Lauf), dass `up` gegen einen laufenden Container sie nicht erneut ausführt. Der Schluss im Kontext von ADR-0030 trägt damit; das `postAttachCommand` schreibt keine Statusdatei, `u-boot` erzeugt keines.
   - *Entscheidung 2* (Erfolgszeile `already running`, Feld `alreadyRunning`) hängt am Zustand und an der Kennung des Containers, nicht an dem, was die CLI ausführt; die Zeile sagt, dass der Container lief, und behauptet nicht, dass nichts lief. Nur der Handbuch-Satz hing an der Verallgemeinerung (Entscheidung 1 dieser ADR).
   - *Entscheidung 3* (kein Zustandsspeicher, keine Selbstheilung) beschreibt, was `u-boot` nicht tut; das `postAttachCommand` läuft, weil die CLI es bei jedem `up` ausführt, nicht als Heilung.

   **Gegenoption:** ADR-0030 ablösen und neu entscheiden — keine Entscheidung hängt an dem widerlegten Satz, der Projektinhaber hat sie am 2026-10-21 getroffen, und eine neue Entscheidung ohne Anlass wäre Nachbessern unter anderem Namen.

3. **`up` und die Startphase von `shell` ändern sich nicht; ein handgeschriebenes `postAttachCommand` läuft bei jedem Start.** `u-boot` erzeugt keines und übergibt keine Option, die es auslässt (`runArgs`); die Startphase von `shell` ist ein `up` ohne `--recreate` ([ADR-0031](0031-devcontainer-shell-ohne-editor.md) Entscheidung 1), ein `postAttachCommand` der Konfiguration läuft also bei jedem Aufruf von `up` und von `shell`, innerhalb des Aufrufs und vor dem Beginn der Sitzung (der Test liest die Zähler unmittelbar nach der Rückkehr von `up`). ADR-0031 Entscheidung 1 gilt unverändert. **Gegenoption:** der Launcher unterdrückt das `postAttachCommand` — ob die CLI `0.89.0` dafür eine Option anbietet, ist nicht gelesen; die Änderung der Argumentliste ist eine Produktentscheidung ohne Bedarfsmeldung und ändert, was ein handgeschriebenes `postAttachCommand` im Alltag tut; nicht gewählt.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: Vermerk im Plan, die zwei ADRs bleiben | kein neues Dokument | zwei `Accepted`-ADRs tragen einen widerlegten Satz mit höherem Rang als der Vermerk; die Pläne der Spezifikation und der Shell leiten aus ADR-0031 ab |
| **B — diese ADR: Prämisse berichtigt (`Supersedes (teilweise)`), Entscheidungen von ADR-0030 bestätigt** | die ADR-Dateien bleiben unverändert; die Berichtigung steht dort, wo man den Satz findet; der Trigger wird eng und wieder prüfbar; Form wie ADR-0030 und ADR-0032 | ein Dokument mehr; wer ADR-0030 liest, liest auch diese |
| C — ADR-0030 ablösen, neu entscheiden | eine geschlossene Fassung | keine Entscheidung hängt am widerlegten Satz; Neuentscheidung ohne Anlass |
| D — Launcher unterdrückt das `postAttachCommand` | `up` führte nur Start und Anlage aus | Option der CLI nicht gelesen; Produktentscheidung ohne Bedarf; ändert das Verhalten handgeschriebener Konfigurationen |

## Konsequenzen

- Positiv: Wer ADR-0030 oder ADR-0031 liest, findet die berichtigte Prämisse; die Entscheidungen von ADR-0030 sind gegen den widerlegten Satz geprüft und bestätigt; der Trigger fragt nach dem, woran die Entscheidungen hängen (`postCreateCommand`, `postStartCommand`); keine Änderung an Code, Tests, Handbuch oder `CHANGELOG.md` (gelesen: sie sagen schon, was der Test zeigt).
- Negativ: Ein `postAttachCommand` einer handgeschriebenen Konfiguration läuft bei jedem `up` und jedem `shell`; seine Laufzeit gehört zur Startphase (die Dauer „rund eine Sekunde“ aus ADR-0031 trägt keine Aussage über eine Konfiguration mit `postAttachCommand`; die Laufzeit kommt hinzu, nicht gemessen). Die Zusage reicht nur für die drei gemessenen Befehle.
- Folgepflicht: Spezifikation und Sicht nach dem Code ([`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)): Sie übernehmen den Satz „keine Lebenszyklus-Befehle“ nicht, sondern nennen die Befehle einzeln mit den Messgrenzen; Aussagen über das Verhalten der CLI nennen den Befehl `up`, den der Docker-Test pinnt; für `shell` gilt, dass seine Startphase ein `up` ist ([ADR-0031](0031-devcontainer-shell-ohne-editor.md) Entscheidung 1), und eine eigene Aussage über Lebenszyklus-Befehle bei `shell` braucht einen Test am Befehl oder entfällt (`AGENTS.md` §3.10). Keine eigene Folge-ADR.
- **Akzeptierte Negative, mit Grund:** (i) Die Befehle `initializeCommand`, `onCreateCommand`, `updateContentCommand`, die Array- und Objektform und scheiternde Befehle gegen einen laufenden Container sind nicht gemessen; keine Entscheidung hängt daran, und `u-boot` erzeugt keinen dieser Befehle, eine Messung hätte keinen Leser. (ii) Ein handgeschriebenes `postAttachCommand`, das ein Startscript erneut ausführt und dessen Fehlschlag verschluckt, ließe ein `failed` aus diesem Aufruf entstehen, und die Meldung „this call did not run the start script“ wäre dort falsch; das Startscript erzeugt nur `u-boot`, ein solcher Aufbau ist ein Eingriff in die erzeugte Konfiguration, und kein Bericht kennt ihn.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Integrationstest (`//go:build docker`) | `TestDevcontainerUp_Docker_AlreadyRunningLifecycleCounters` zählt das `postCreateCommand`, das `postStartCommand` und das `postAttachCommand` (Tabelle im Kontext); er färbt sich rot, wenn die CLI das `postCreateCommand` oder das `postStartCommand` gegen einen laufenden Container erneut ausführt oder das `postAttachCommand` nicht mehr bei jedem `up` | `make test-docker` |
| Unit | die Meldung nach einem gelungenen Lauf der CLI (`TestDevcontainerUp_FailedStatusAfterSuccess`) bleibt, wie sie ist | `make test` |

Die Wortwahl der Dokumente („die Befehle einzeln nennen“) ist nicht maschinell prüfbar; `AGENTS.md` §3.10 (Bedingung mit Test) und die Prüfung im Review decken sie.

## Re-Evaluierungs-Trigger

Eine Hebung des Pins der Dev-Container-CLI: der Docker-Test mit den drei Zählern wird mitgelesen, bevor der Pin gehoben wird; der Docker-Test zeigt, dass die CLI gegen einen laufenden Container das `postCreateCommand` oder das `postStartCommand` erneut ausführt (der Schluss im Kontext von ADR-0030 fällt) oder dass das `postAttachCommand` nicht mehr bei jedem `up` läuft (Entscheidung 3 dieser ADR ändert sich); ein Nutzerbericht, dass ein `postAttachCommand` die Startphase von `up` oder `shell` aufhält oder stört (Option D neu bewerten); eine Messung auf einem anderen Host als Linux `amd64`, die abweicht; `u-boot` erzeugt künftig ein `postAttachCommand`.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-22 | Entwurf und `Accepted` durch den Architect im Auftrag des Projektinhabers, im Gespräch („ohne nachzufragen“); die ADR berichtigt eine Tatsachen-Prämisse und trifft keine neue Produktentscheidung | [Review-Report](../../reviews/2026-10-22-slice-v1-devcontainer-up-failed-nach-zustand.md) F-2 |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-0033`.
