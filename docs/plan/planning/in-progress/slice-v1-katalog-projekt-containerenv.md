# Slice V1: Projekt-Umgebungswerte des DevContainers — der Konfig-Schlüssel `devcontainer.containerEnv` (CR nach dem Muster der Werkzeugketten-Kette, Umsetzung nach der Freigabe)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle (Vorschlag des Projektinhabers, 2026-11-03: „im
`.devcontainer`-Dockerfile gibt es noch einen kleinen Block am Ende mit
ENV-Var. Können wir diese nicht in `u-boot.yaml` setzen" — Vorschlag
angenommen; CR-Slice zuerst, Muster der Werkzeugketten-Kette der zwei Go- und
des JDK-Schlüssels).

**Bezug:** [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)
(erzeugte `devcontainer.json`; der Schlüssel erweitert die Erzeugung),
[`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)
(Konfig-Schema der Dev-Container-Form), [`LH-FA-CONF-005`](../../../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern)
(`config set`/Anzeige berühren den neuen Schlüssel), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
(Abweisungen je Exit-Code `10`, Muster [`SPEC-054`](../../../../spec/spezifikation.md)).
Muster-ADR [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)
(nur als Form-Vorbild für die CR-→-Umsetzungs-Kette).

**Berührte Spec-Stellen** (gelesen am 2026-11-03): [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)
(Config-Pfade; der neue Schlüssel kommt in die zwei Pfade-Zeilen), [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)
(Kontext Erzeugung, gelesen am JDK-Stand), [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)
(Abschnitt *Cache-Werte*: die bestehenden `containerEnv`-Werte der zwei
Go-Schlüssel — der Merge-Verhaltens-Fall; Historie-Zeile 1608 nennt
`BUILD_MODE`/`GOFLAGS`/`GOTOOLCHAIN` als „nicht belegt und nicht zugesagt" —
genau die Lücke, die dieser Slice schließt), [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)
(`/devcontainer/devcontainer.json`-Namenslesung und `containerEnv`-Auflösung
zu `Config.Env` — bestehend, wird nicht geändert). Der Verweis zeigt
**aufwärts**: Lastenheft und Spezifikation nennen diesen Slice nie.

**Verantwortlich:** pt9912 (Implementer; Umsetzung nach der CR-Closure 2026-11-03, Reihenfolge-Beschluss des Projektinhabers)

**Autor:** pt9912 (Rolle Planner, im Auftrag des Hauptagenten). **Datum:** 2026-11-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten.

**Ziel:** Ein Projekt trägt seine Umgebungswerte für den erzeugten
DevContainer in [`u-boot.yaml`](../../../../spec/lastenheft.md) (Schlüssel
Vorschlag: `devcontainer.containerEnv`, Map Name → Wert); `u-boot generate
devcontainer` lädt sie verwaltet und idempotent in das `containerEnv`-Objekt
der erzeugten `.devcontainer/devcontainer.json` — neben bestehenden
Generator-Werten (die Cache-Verzeichnisse der zwei Go-Schlüssel bleiben), mit
Validation und Exit-Code `10` je fachlichem Fehler — und das u-boot-Projekt
selbst hat seine drei Rest-ENV-Werte (`BUILD_MODE`, `GOFLAGS`,
`GOTOOLCHAIN`) in seine eigene `u-boot.yaml` gezogen: der handgeschriebene
ENV-Block hinter dem Managed-Block von `.devcontainer/Dockerfile` entfällt
(Fundstelle: `.devcontainer/Dockerfile`, Abschnitt *Hand-written rest*).

**Vorbefund** (gelesen am 2026-11-03):

- **Die ENV-Ebene ist heute die einzige Trägerin für Projekt-Werte.** Der
  Generator schreibt `containerEnv` nur für die zwei Go-Schlüssel
  (`containerEnv`-Abschnitt der erzeugten `devcontainer.json`, Cache-Werte
  unter `${containerWorkspaceFolder}`, [`SPEC-054`](../../../../spec/spezifikation.md));
  ein Projekt hat keinen Weg, eigene Werte konfiguriert zu setzen.
- **Die Semantik der zwei Ebenen bleibt ein Messpunkt.** Dockerfile-`ENV`
  backt in das Image (gilt in `docker build`-RUN-Schritten und mit jedem
  Neubau), `containerEnv` gilt pro Container-Prozess der erzeugten
  `devcontainer.json` (Auflösung zu `Config.Env`, [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)
  Abschnitt *Grenzen*); für die drei Werte (Build-Umgebung des Go-Einsatzes
  in der Container-Shell) reicht `containerEnv` vermutlich — der Docker-Test
  misst es, es steht nicht aus dem Kopf im CR
  ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
- **Die Namens-Validation existiert schon als Werkzeug.**
  `driven.ValidateDevcontainerEnvName` prüft Namen der Umgebung der
  Dev-Container-Übergabe (`internal/hexagon/adapter/driven/`, Fundstelle beim
  Start nachzulesen) — der CR entscheidet, ob sie den Schlüssel-Namensraum
  trägt.
- **Der Managed-Block-Einsatz ist bestehender Vertrag.** u-boot generiert
  `devcontainer.json` (und die `containerEnv`-Werte der Schlüssel stehen im
  verwalteten Teil); das Merge-Verhalten der Projekt-Werte mit den
  Cache-Werten ist §4-Frage 2 — nicht still mergen
  ([`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz)).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Dockerfile-`ENV`-Zusätze als zweite Ebene** — keine zweite
  Installationsform unterhalb von `containerEnv`: der CR entscheidet die
  Ebene (Vorschlag: nur `containerEnv`); ein Wert, der ins Image gehöre,
  braucht einen eigenen Auftrag (gibt es heute nicht — Vorbefund liest sich
  leer).
- **`${localEnv:…}`-Interpolation in Werten** — die Render-Stelle `remoteEnv`
  liest die Form ([`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes));
  ob auch `containerEnv`-Werte die Form tragen (und der Hilfs-Container die
  Namen durchreicht), ist ein eigener Vorgang mit eigener Messung, nicht Teil
  dieses Slices.
- **`remoteEnv` als Projekt-Schlüssel** — die bestehende Render-Stelle bleibt
 ; eine Projekt-Konfiguration für sie wäre ein weiterer Schlüssel, kein Teil
  dieses Vorschlags.
- **Die Namen-Übergabe an den Hilfs-Container** (`up`/`shell` lesen Namen aus
  der erzeugten `devcontainer.json`) — Bestand; nur wenn der CR die
  Namensquellen erweitert (Vorschlag: tut er nicht), berührt er sie,
  sonst unverändert.

**Keine Mindestzahl.** Ein echter Ausschluss ist besser als vier erfundene.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; Gate-Läufe und Closure-Pflichten
zählen nicht mit.

- [ ] **CR-Wortlaut und Freigabe:** Lastenheft-Absatz für
  [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)
  (Schlüssel, Map-Form, Merge-Verhalten je §4-Frage 2, Abweisungen je
  Exit-Code `10`, Grenz-Sätze S10/S11-Muster), [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)-Zeilen;
  freigegeben vom Projektinhaber, **bevor** Verträge geändert werden.
- [ ] **Umsetzung:** Konfig-Schema samt Validation (`u-boot config set`
  prüft die Map — Name gegen den Namensraum, Wert vom Typ String),
  Generator-Render der Werte in das `containerEnv`-Objekt der erzeugten
  `devcontainer.json` (idempotent, Golden Cases: mit Werten, ohne Werte
  byte-identisch zum heutigen Stand, je Profil, Merge-Fall mit den
  Cache-Werten je §4-Frage 2), Abweisungs-Meldung je Fehlerfall, Docker-Test
  misst die Werte im Container-Prozess (`env`/`go env GOTOOLCHAIN` im
  gebauten Container — der ENV-vs-containerEnv-Satz des CR).
- [ ] **Eigenes Devcontainer umgestellt:** das u-boot-Projekt trägt die drei
  Rest-Werte in `.devcontainer/u-boot.yaml` (oder wie die Projektdatei am
  Umsetzungs-Stand heißt), der ENV-Block hinter dem Managed-Block entfällt;
  `.devcontainer` regeneriert; `make targets` der Sicht mit `make gates`
  grün am umgestellten Stand.

Zusätzlich (nicht gezählt): `make gates` grün, `make docs-check` 0 Befunde;
Review und Verification mit Reports unter `docs/reviews/`; Doku-Update
(Handbuch, READMEs falls betroffen, `CHANGELOG.md` ohne Kennungen im
Nutzertext, `harness/replay.md`); Closure-Notiz mit
Steering-Loop-Lerneintrag; Beobachtungs-Register fortgeschrieben oder „keine
Beobachtung"; jedes Risiko aus §6 mit Ausgang; die drei Paarungen getragen.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt;
Fundstellen liest der Implementer am Start neu
([`AGENTS.md`](../../../../AGENTS.md) §3.10) und gleicht nach jeder Fix-Runde
erneut ab.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` (CR-Slice zuerst) | update | Absatz der Zusage (Muster der Toolchain-Zeilen) samt Merge-Verhalten und Abweisungen; Versionshistorie |
| `spec/spezifikation.md` | update | [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) die zwei Pfade-Zeilen; Technik-Abschnitt (Render-Form, Merge, Validation) samt Beleg-Tabelle (Muster [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)) |
| Konfig-Schema und Validation (`internal/hexagon/`, Fundstellen beim Start) | update | Schlüssel, Map-Typ, `config set`/Anzeige, Validation je Exit-Code `10` |
| Generator (`templates.go`, `generate.go`, `templates/devcontainer/devcontainer.json.tmpl`) | update | Render der Projekt-Werte in `containerEnv` (verwaltet, idempotent), Merge-Fall mit den Cache-Werten je §4-Frage 2 |
| Tests (`internal/hexagon/application/`, `cmd/uboot/`, `internal/e2e/` je Fundstelle) | neu / update | Golden Cases, Validation-Fälle, Nutzerpfad-Test, Docker-Test der Werte im Container-Prozess |
| `repo.mk` (Kommentar des Docker-only-Absatzes) | update | Berichtigt nach §3.10 (erste Fix-Runde): der Kommentar nannte `.devcontainer/Dockerfile` als Träger von `BUILD_MODE=local` — der Träger ist jetzt die `containerEnv` des erzeugten `devcontainer.json`; Fundstelle gelesen `repo.mk:36` |
| `.devcontainer/` (dieses Projekt) | update | Die Projektdatei heißt am Umsetzungs-Stand `u-boot.yaml` im Repo-Wurzel (nicht `.devcontainer/u-boot.yaml`, DoD-Lesung „oder wie die Projektdatei am Umsetzungs-Stand heißt“): die drei Werte stehen dort unter `devcontainer.containerEnv` (geschrieben mit `config set`, Docker-Modus); der ENV-Block hinter dem Managed-Block von `.devcontainer/Dockerfile` entfällt; Regeneration aktualisierte das `containerEnv`-Objekt von `.devcontainer/devcontainer.json` (NoOp im zweiten Lauf) |
| `docs/user/benutzerhandbuch.md`, READMEs, `CHANGELOG.md`, `harness/replay.md` | update | Schlüssel-Dokumentation ohne Kennungen im Nutzertext; Golden-Ablage in `replay.md` |
| `spec/architecture.md`, ADRs, `Makefile`, `.golangci.yml`, `.d-check.yml` | **unverändert** | Gates bleiben ([`AGENTS.md`](../../../../AGENTS.md) §3.6); ein ADR-Bedarf ergibt sich nur, wenn der CR eine Architektur-Frage aufwirft (Vorschlag: nicht) |

## 4. Trigger und Entscheidungen des Projektinhabers

(CR-Auftrag mit Vorschlägen — der Implementer wählt nicht; die Freigabe des
Wortlauts liegt beim Projektinhaber, Muster der Werkzeugketten-Kette.)

1. **Schlüsselname und Ort** — Vorschlag: `devcontainer.containerEnv` als Map
   Name → Wert (String), gelesen aus der Projekt-Konfigurationsdatei
   (`u-boot.yaml`); Namens-Norm per `ValidateDevcontainerEnvName` (Vorschlag,
   zu lesen beim CR-Schnitt).
2. **Merge-Verhalten mit den Cache-Werten der zwei Go-Schlüssel** — Vorschlag:
   Abweisung mit Exit-Code `10`, wenn ein Projekt-Wert `GOCACHE` oder
   `GOMODCACHE` setzt und der Schlüssel aktiv ist (zwei Schreiber einer
   Stelle, Muster des Konflikt-Falls `java` + `java-toolchain`); ohne
   Schlüssel darf das Projekt sie setzen.
3. **Ebene der Werte** — Vorschlag: nur `containerEnv` (Container-Prozesse);
   Dockerfile-`ENV` bleibt dem Generator und dem handgeschriebenen Rest
   vorbehalten; der Docker-Test misst, dass die Container-Shell die Werte
   trägt (`go env GOTOOLCHAIN`/`echo $BUILD_MODE`).
4. **Idempotenz und Verwaltet-Form in `devcontainer.json`** — Vorschlag: die
   Werte stehen im verwalteten Teil neben den Generator-Werten; eine
   Regeneration schreibt sie neu, ein entfernter Wert verschwindet (Muster
   des Managed-Blocks; keine Duplikate).

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der
CR-Schnitt ist der erste Vorgang (eigener CR-Slice nach dem Muster
[`slice-cr-devcontainer-java-toolchain`](../done/slice-cr-devcontainer-java-toolchain.md)),
kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der CR verlangt
  mehr als den Schlüssel (etwa ein Render-Verhalten außerhalb der
  `containerEnv`-Ebene) — der Architect trennt.
- `in-progress` → `open` (blockiert): der Docker-Test messbar widerspricht
  der ENV-Annahme (Werte erreichen die Build-Umgebung nicht über
  `containerEnv`) — dann steht der Weg für die drei u-boot-Werte zur
  erneuten Entscheidung (ENV-Rest bleibt).

## 5. Closure-Trigger

DoD vollständig; `make gates` und `make docs-check` grün; der Docker-Test
läuft im CI-Lauf `integration-docker` sichtbar bestanden (oder die Closure
benennt den ersetzenden Nachweis mit Grund); Review- und Verification-Report
liegen vor; Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

1. **Eine Zusage reicht weiter als die Messung** (ENV-vs-`containerEnv`-Semantik
   aus der Annahme statt per Docker-Test). — **Ausgang:** bei Closure.
   Vorsorge: der CR-Wortlaut trägt nur Gemessenes; der Docker-Test misst den
   Wert-Weg, bevor der Satz in den Vertrag geht
   ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
2. **Still gemergte `containerEnv`-Konflikte** (Projekt-Wert und
   Cache-Werte-Kollision untergehen ohne Abweisung). — **Ausgang:** bei
   Closure. Vorsorge: Validation je §4-Frage 2 mit Exit-Code `10` und
   Wächter-Test; [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz)-Disziplin.
3. **Eine Bedingung der Technik-Stelle ohne Test** (Merge-Verhalten,
   Idempotenz, byte-identisch-ohne-Werte) — **Ausgang:** bei Closure.
   Vorsorge: Beleg-Tabelle nach dem Code (Muster [`SPEC-054`](../../../../spec/spezifikation.md)),
   Golden Cases
   ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)).
4. **Plan-Aussagen über Code, die am Fundort nicht gelesen sind** (Schemata,
   Render-Stellen) — **Ausgang:** bei Closure. Vorsorge: Fundstellen beim
   Start neu gelesen; §3.10-Abgleich nach jeder Fix-Runde
   ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md)).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders als
geplant · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices ·
Risiken aus §6 · Drei Paarungen.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `internal/hexagon/application` (APP),
`docs/user/`/READMEs/`CHANGELOG.md` (DOK), `docs/plan/` (DOK); laut
[`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)
alle Greenfield.

**Vorgelagert — offene Beobachtungen sichten:** [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)
(Risiko 1), [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)
(Risiko 3), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md)
(Risiko 4). Evidenz-/Diskrepanz-Risiko: mittel (die ENV-Ebene ist ein
Verhalten des Dev-Container-CLI-Einsatzes, nicht von `u-boot` — die
Annahme trägt kein Messprotokoll).

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield (kein
Modus-Begründungsblock nötig).