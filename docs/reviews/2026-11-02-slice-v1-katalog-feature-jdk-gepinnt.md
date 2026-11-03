# Review-Report: slice-v1-katalog-feature-jdk-gepinnt — 2026-11-02

**Review-Art:** Code-Review — geprüft gegen den Slice-Plan (§1 Befund und Nicht-Ziele,
§2 DoD, §3 Plan-Tabelle, §4 Trigger/Entscheidungen, §6 Risiken), die Hard Rules
(`AGENTS.md` §3 mit §3.10-Abgleich), die Beleg-Tabelle der Spezifikation
([`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)),
den Freigabe-Wortlaut des Lastenhefts 0.31.0 und die Vorlagen des Vorgangs
[`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) /
[`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md).

**Gegenstand:** die Slice-Commits `0ba3c6fa` / `d55002a8` (Lifecycle-Move und
Verweis-Nachzug, reine Moves), `767c1122` (Roadmap-State-Berichtigung),
`9faa3137` (Umsetzung: dritter Katalog-Schlüssel `java-toolchain`, Dockerfile-Block,
Tests, Vertrag und Doku, Messung) und `d213e44e` = `HEAD` (Lint-Fix-Runde:
`featureCatalogue`-Split, gofmt).

**Skill:** `v6.18.1` · `.harness/skills/reviewer.md` @ `d213e44e` · <!-- d-check:ignore (Adopter-spezifischer Skill-Pfad, existiert im Ziel-Repo ggf. nicht) -->
**Modell:** claude (glm-5.3-flash) · **Datum:** 2026-11-02

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- [Slice-Plan](../plan/planning/done/slice-v1-katalog-feature-jdk-gepinnt.md) (Stand `d213e44e`; §3-Tabelle gegen `git show --name-stat` der Slice-Commits abgeglichen)
- [ADR-0041](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)
  (`Accepted`, der Typ), [ADR-0024](../plan/adr/0024-ollama-im-sandbox-devcontainer.md)
  (Release-Archiv-/Prüfsummen-Muster), [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) (Namensraum nur bei OCI-Quelle)
- [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)
  (0.31.0), [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette),
  [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte),
  [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes),
  [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv),
  [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)
  (Zeile `java-toolchain`), [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config),
  [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)
- Der CR-Vorgang [cr-vorschlag-java-toolchain](../plan/planning/done/cr-vorschlag-java-toolchain.md) (Freigabe, Fragen F4 bis F8 mit Antworten)
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules, §3.10), [`harness/conventions.md`](../../harness/conventions.md)
  (Ablage-Regel, `MR-*`), [`harness/replay.md`](../../harness/replay.md) (Werkzeugketten-Zeile),
  Messprotokoll [slice-v1-katalog-feature-jdk-gepinnt](../messungen/slice-v1-katalog-feature-jdk-gepinnt/README.md)
- Vorherige Findings am selben Bereich: die drei Reports dieses Tages zum
  Vorgänger-Vorgang [`slice-v1-katalog-gepinnte-toolchain`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md)
  in [`docs/reviews/`](README.md) (Klassen-Tradition; keine offenen Findings übernommen)

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.1` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Beleg-Zeile „Der Schlüssel `java-toolchain` ist als eingebautes Feature in beiden Profilen wählbar" nennt für die Profil-Hälfte `TestGenerateDevcontainer_Toolchain_BlocksInBothProfiles` — der Umsetzungs-Diff hat die Schlüssel-Tabelle dieses Tests nicht um `java-toolchain` erweitert (sie trägt `go-toolchain`, `golangci-lint` und `both`); für den dritten Schlüssel pinnt die erzeugte `devcontainer.json` ohne Schlüssel unter `features` nur des Standardprofils (Byte-Golden), und der Block in beiden Profilen trägt nur der cmd-Lauf (`TestRun_Generate_ToolchainNotices`, `ARG`-Assert). | [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) · `AGENTS.md` §3.10 | [spec/spezifikation.md](../../spec/spezifikation.md) · "`TestFeatureCatalogue_EntriesAreValid/java-toolchain`" | nein — Review-only; ein erweiterter Lauf der Schlüssel-Tabelle um ein `java-toolchain`-Paar träfe es | Beleg-Tabelle nennt einen Test, der die Bedingung nicht trägt |
| F-2 | MEDIUM | Die Beleg-Zeile „Warn-Meldung des Schlüssels `java-toolchain` samt Vertriebs-Form und Grenz-Hinweis" nennt `TestGenerateDevcontainer_Toolchain_Notices` und `TestRun_Generate_ToolchainNotices` — kein Assert dieser Tests liest die Vertriebs-Form („the JDK of the Eclipse Temurin distribution of the LTS version 21"): der erste prüft Schlüssel-Präsenz, Code, Level und die zwei Flag-Felder, der zweite das `Warning:`-Präfix der Ausgabe. | [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) · `AGENTS.md` §3.10 | [internal/hexagon/application/devcontainer_toolchain.go](../../internal/hexagon/application/devcontainer_toolchain.go) · "the JDK of the Eclipse Temurin distribution of the LTS version 21" | nein — Review-only; ein Assert auf die Vertriebs-Form in der Notices-Tabelle träfe es | Beleg-Tabelle nennt einen Test, der die Bedingung nicht trägt |
| F-3 | MEDIUM | Der §3.10-Abgleich der Plan-Tabelle gegen die Slice-Commits ist unvollständig: drei berührte Dateien tragen keine §3-Zeile — `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` (die Abgrenzung „kein containerEnv-Wert"), `cmd/uboot/toolchain_test.go` (der Nutzerpfad-Test zum Konflikt, 51 Zeilen) und `docs/plan/planning/in-progress/roadmap.md` (State-Berichtigung `767c1122`); die Tabelle nennt für Templates nur `Dockerfile.tmpl` und für Tests nur die application-, e2e- und export-Dateien. | `AGENTS.md` §3.10 · [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) | [Slice-Plan](../plan/planning/done/slice-v1-katalog-feature-jdk-gepinnt.md) · "Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt" | nein — Review-only; der Abgleich der §3-Tabelle gegen die `--name-only`-Listen der Commits trägt es | Plan-Aussage über Code/Berührung trägt dem Stand nicht (BEO-DOK/plan-aussage-ueber-code) |
| F-4 | LOW | Der Bestand des Messprotokoll-Index führt die neue Messung nicht: `docs/messungen/README.md` trägt 19 Einträge, alle sonstigen Protokoll-Verzeichnisse stehen in der Tabelle — [`slice-v1-katalog-feature-jdk-gepinnt`](../plan/planning/done/slice-v1-katalog-feature-jdk-gepinnt.md) fehlt trotz neu angelegten Verzeichnisses samt Skript-Zeile im eigenen README. | [`harness/conventions.md`](../../harness/conventions.md) §Ablage von Messprotokollen und Skripten | [docs/messungen/README.md](../messungen/README.md) · "## Bestand" | nein — Review-only; kein Sensor prüft Bestand-Vollständigkeit | Neu angelegtes Artefakt fehlt im Bestand-Index |
| F-5 | LOW | Vier Kommentar-Stellen tragen die Zwei-Schlüssel-Aufzählung des Eintragstyps weiter, obwohl der Umsetzungs-Diff sie auf drei erweitert: der Typ-Kommentar von `devcontainer_toolchain.go` („the keys ..., carry the toolchains of the u-boot build itself", „lastenheft 0.30.0"), die Zeichen-Grenze desselben Files („the character set of a ... value of the two keys" — sie gilt für alle drei Schlüssel, auch `21.0.13+9`), die Test-Kommentare `devcontainer_features_test.go` („the two pinned-toolchain keys carry theirs") und `generate_toolchain_test.go` („accepts the overrides of the two keys" — das Fixture trägt drei), dazu die Spec-Zeile der doctor-Prüfung („mit beiden aktivierten Schlüsseln" — das Test-Fixture trägt drei). | `AGENTS.md` §3.7, §3.9(b) · `v6.18.1` · `regelwerk/grundlagen-harness-dateien.md` §Was ein Kommentar trägt | [internal/hexagon/application/devcontainer_toolchain.go](../../internal/hexagon/application/devcontainer_toolchain.go) · "the keys `go-toolchain` and `golangci-lint` carry the" | nein — Review-only; kein Gate liest Kommentar-Aufzählungen gegen den Code | Kommentar-Aufzählung trägt der durch den Diff geänderten Menge nicht |
| F-6 | LOW | Die Roadmap-State-Berichtigung (`767c1122`) hat die Leerzeile vor dem Abschnitt gestrichen, so dass unter dem Bestand die Zeile `)).## Nächste Wellen` in einem Absatz endet: Der ATX-Head leitet nur am Zeilenanfang, die Abschnitts-Überschrift ist im Rendering Teil des Absatzes, und der fünf-Abschnitte-Aufbau (Modul 6) zeigt die Sektion „Nächste Wellen" ohne eigene Head-Zeile; `make docs-check` fällt das nicht auf, weil kein Link den Anker `#nächste-wellen` trägt — eine eingehende Referenz auf die Sektion würde erst damit rot. | [MR-003](../../harness/conventions.md#mr-003) · `v6.18.1` · `regelwerk/modul-06-roadmap.md` §Roadmap-Struktur (Modul 6) | [docs/plan/planning/in-progress/roadmap.md](../plan/planning/in-progress/roadmap.md) · ").## Nächste Wellen" | nein — Review-only; d-check prüft Link-Ziele und Anker, nicht die Head-Struktur der Datei selbst | Abschnitts-Überschrift in die Absatzzeile gezogen |
| F-7 | LOW | Die Berichtigung des Handbuch-Beispiels trägt einen Genus-Fehler: „ein gepinnte JDK-Installation steuert der Katalog-Schlüssel ...". | Maintainability (Doku-Präzision, Kategorie LOW des Reviewer-Skills) | [docs/user/benutzerhandbuch.md](../../docs/user/benutzerhandbuch.md) · "ein gepinnte JDK-Installation" | nein — Review-only | Genus-Fehler in der Nutzer-Doku |
| F-8 | LOW | Das Messprotokoll trägt die Werte der drei Messungen nur im README, ohne Rohdaten-Dateien (`.txt`/`.log`) neben sich — die Ablage-Regel nennt die Rohdaten als Bestandteil des Protokoll-Verzeichnisses, und frühere reine Lese-Messungen im selben Baums tragen sie je Lauf; die Umgebung-Zeile deklariert den Zustand, aber die Begründung trägt nicht, warum die gelesenen Ausgaben nicht archiviert sind — die gelesenen Prüfsummen sind so nur gegen das Repo-Statement, nicht gegen den Lauf selbst abgleichbar. | [`harness/conventions.md`](../../harness/conventions.md) §Ablage von Messprotokollen und Skripten | [docs/messungen/slice-v1-katalog-feature-jdk-gepinnt/README.md](../messungen/slice-v1-katalog-feature-jdk-gepinnt/README.md) · "Die Messung liest nur; sie schreibt nichts ins Repo außer diesem Protokoll." | nein — Review-only; `make ablage-pruefung` prüft Zulässigkeit, nicht Vollständigkeit der Rohdaten | Messprotokoll ohne Rohdaten |
| F-9 | INFO | Der DoD-Punkt der Messung verlangt „Datum und Dev-Container-CLI-Version im Protokoll" — das Protokoll trägt das Datum, aber statt einer CLI-Version die Angabe, dass keine CLI- und Engine-Abhängigkeit bestand (die Tag- und Cask-Lesung lief über Registry- und GitHub-API per `curl`); die Abweichung ist offengelegt, ihre Bewertung gegen die DoD-Formulierung ist Verifier-Aufgabe. | Maintainability (Übergabe an den Verifier) | [Slice-Plan](../plan/planning/done/slice-v1-katalog-feature-jdk-gepinnt.md) · "Datum und Dev-Container-CLI-Version im Protokoll" | nein — Review-only | DoD-Wortlaut trägt der Protokoll-Stand nicht |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `spec/lastenheft.md` — Nicht-Berührung und Freigabe-Wortlaut | geprüft, ohne Befund — kein Slice-Commit (`9faa3137`, `d213e44e`) berührt die Datei (`grep` über beide `--name-only`-Listen: 0 Treffer); Lastenheft 0.31.0 stammt aus dem CR-Vorgang [`slice-cr-devcontainer-java-toolchain`](../plan/planning/done/slice-cr-devcontainer-java-toolchain.md) (`1ddada3e`, davor); der java-Absatz (Zeile 427) nennt keine ADR-, keine Slice- und keine Zahlen-Angabe jenseits der Freigabe — die Decken-Regel der Zusagen hält |
| Decken-Regel der Zusagen (Linse [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md)) | geprüft, ohne Befund — CHANGELOG trägt Promised/Not promised treu zum CR-Muster, das Handbuch begrenzt die Übersteuerung auf die gemessene Umgebung, der Brew-Befund steht als Vertriebs-Vergleich ohne zweite Installationsform |
| Dual-Classifier-Linse | geprüft, ohne Befund — der Konflikt ist kein neuer Sentinel: `javaFeatureActive`/`javaToolchainConflictMessage` erzeugen den Fehler am selben Ort wie der Go-Fall, und der Wrap in `driving.ErrConfigValueInvalid` steht zentral am Aufrufer (`generate.go`, der `resolveToolchain`-Aufruf); die Diagnostic-Abbildung des `generate`-Subkommandos kennt die Klasse [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) schon, der cmd-Test pinnt den Umschlag; die Warn-Einträge tragen je Schlüssel [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) |
| Exit-Code-10-Fälle des dritten Schlüssels | geprüft, ohne Befund — alle drei Fälle (Quelle unlesbar, Version nicht überdeckt, Konflikt externes `java` plus `java-toolchain`) enden im Wrap von `driving.ErrConfigValueInvalid` (Exit 10, nichts geschrieben), je Fall ein application-Test samt `assertNothingWritten` und ein cmd-Lauf am Nutzerpfad mit `--dry-run --json` und dem Code der Klasse (siehe oberste Zeile) |
| Dockerfile-Block (`Dockerfile.tmpl`, dritter Block) | geprüft, ohne Befund — Reihenfolge Download → `sha256sum -c -` → `tar` nach `/opt`, Same-Origin an derselben Release-URL (Archiv und Prüfsummen-Datei aus derselben Veröffentlichung, frisch am Vertrieb nachgelesen, s. Sensors), POSIX-sh-Festigkeit (`printf`/`tr`/`sed`, keine bash-Substitution im `RUN`), Case-Zeile je Architektur mit Abbruchtext, `ENV JAVA_HOME`/`PATH`, die `or`-Zeile trägt `.JdkToolchain`, kein `--privileged`, kein Schlüssel unter `features` |
| `devcontainer.json.tmpl` — kein `containerEnv`-Wert | geprüft, ohne Befund — der containerEnv-Block bleibt an die zwei Go-Schlüssel gebunden; der neue Template-Kommentar steht als `{{- /* ... */ -}}` und liest nicht in die erzeugte Datei (Byte-Golden ohne `containerEnv`) |
| Gate-Dateien — Lockerung (AGENTS.md §3.6) | geprüft, ohne Befund — kein Slice-Commit berührt `Makefile`, `Dockerfile`, `.golangci.yml`, `.d-check.yml`, `spec/architecture.md` oder `docs/plan/adr/` (`--name-only`-Listen der zwei Umsetzungs-Commits) |
| Catalogue-Split (`d213e44e`) | geprüft, ohne Befund — `featureCatalogue()` fasst `ociFeatureCatalogue` und `blockAndOptionCatalogue` zusammen, kein Eintrag doppelt oder verloren (`TestFeatureCatalogue_KeysCoverSpecExamples` trägt 20 Schlüssel, `TestFeatureCatalogue_EntriesAreValid` läuft je Eintrag und prüft die `dockerfileBlock`-Invarianten; die Einrückung der Konstante `featureJavaToolchain` ist gofmt-Ausrichtung, ohne Umzug) |
| Test-Namen gegen die Beleg-Tabelle ([`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)) — die übrigen Java-Zeilen | geprüft, ohne Befund — alle genannten Test-Namen existieren am Diff und tragen die Bedingung: Katalog-Tabellen, `OverrideRendersReadChecksums` samt `OverrideSpellings/java spellings`, `VersionNotCoveredExit10`, `UnreadableChecksumSourceExit10`, `JavaFeatureConflictExit10` samt `TestRun_Generate_JavaToolchainConflictExits10` (Nutzerpfad, Code wie in der Dual-Classifier-Zeile), `JSONGolden` samt `assertNoContainerEnv`, `Doctor_ToolchainKeys_NoFindings`, `ToolchainChecksumParsing` (url forms, beide Architekturen), Docker-Rows in `TestE2E_Toolchain_BuildsAndReportsPins` / `PinMismatchBreaksBuild/jdk` / `OverrideBuildsAgainstReadChecksums` |
| Handbuch — Kennungs-Freiheit und Ist-Zustand (AGENTS.md §5 Regel 6) | geprüft, ohne Befund — die neuen Absätze und Tabellen-Zeilen tragen keine `LH-*`-, `SPEC-*`- oder slice-Kennung; die Revision `21.0.13+9` und die Abwesenheit des Tags `21` sind gemessen (Protokoll M2) und in diesem Lauf am Registry-Gegenstand frisch nachgelesen (41 Tags, `1.8.3` vorhanden, kein `21`) |
| `README.md`/`README.de.md` — Schlüsselzählung | geprüft, ohne Befund — 20 Built-in-Features passend zum Katalog (20 Einträge, `TestFeatureCatalogue_KeysCoverSpecExamples`) |
| `harness/replay.md` — Werkzeugketten-Zeile | geprüft, ohne Befund — die Zeile nennt den dritten Schlüssel mit beiden Schreibweisen, die JDK-Abgrenzung am `containerEnv`, die Abweisungen samt Konflikt („kein `sdkman`-Fall") und die Grenze der Vier-Stellen-Kopplung auf die zwei Go-Schlüssel; Fresh-State, Double-Run, Stale und User-Content bleiben getragen |
| Messskripte am [Protokoll](../messungen/slice-v1-katalog-feature-jdk-gepinnt/README.md) (Skript-Zeile) | geprüft, ohne Befund — drei Skripte am Konventions-Ort: Same-Origin-Gegenprobe (M1), Tag- und Metadaten-Lesung (M2), Cask-Lesung (M3) |
| Prüfsummen gegen die Vertriebs-Quelle (Risiko 1 des Plans) | geprüft, ohne Befund — die Werte der Template-`ARG` (Default-Pin) und der Test-Fixtures (Override `21.0.9+10`) stimmen mit dem Protokoll überein und wurden in diesem Lauf live nachgelesen (`jdk-21.0.13+9` x64/aarch64 und `jdk-21.0.9+10` x64/aarch64, alle Werte deckungsgleich, Tag- und Asset-Form wie im Code) |
| Lifecycle- und Move-Disziplin | geprüft, ohne Befund — `0ba3c6fa` reiner Move, Verweis-Nachzug als eigener Commit (`d55002a8`, 15 eingehende), State-Berichtigung nach dem Move (`767c1122`); der Plan liegt nur in `in-progress/` |
| §3-Abgleich — die mit Zeilen genannten Dateien | geprüft, ohne Befund (Zuordnung) — jede andere Änderungs-Art der Plan-Tabelle entspricht der Berührung des Diffes; die drei Zeilen-Lücken stehen als F-3 |

## Sensors

| Sensor | Stand |
|---|---|
| `make docs-check` | gelaufen am Handoff (0 Befunde, 1049 Dateien) und in diesem Lauf vor dem Report; nach dem Schreiben des Reports erneut — Ausgang in der Schluss-Antwort (die Report-Datei trägt neue Links) |
| `make lint` | gelaufen in diesem Lauf — grün (Docker-Modus, der Standard) |
| `make test` | gelaufen in diesem Lauf — grün (Unit-Tests im Docker-Test-Stage, inklusive application und cmd) |
| `make ablage-pruefung` | gelaufen in diesem Lauf — grün (26 Fälle, 0 Fehlschläge) |
| `make gates` | gelaufen in diesem Lauf — grün (Exit 0; inklusive `coverage-gate`: 92,20 % gegen die Schwelle 90 %) |
| `make test-docker` | nicht gelaufen — die Docker-Tests der JDK-Rows sind Implementer-Evidence (gezielte Läufe gemeldet: Pins, Negativprobe `jdk`, Override am echten Engine-Host); der Nachweis des CI-Laufs `integration-docker` ist Verifier-Closure |
| `make verify-depguard` | nicht gelaufen — keine Import-Regel berührt: der Diff trägt keinen neuen Import, `devcontainer_toolchain.go` nutzt den bestehenden `driven.Fetcher`-Port |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 5 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Beleg-Tabelle nennt einen Test, der die Bedingung nicht trägt · Plan-Aussage über Code/Berührung trägt dem Stand nicht (BEO-DOK/plan-aussage-ueber-code) · Neu angelegtes Artefakt fehlt im Bestand-Index · Kommentar-Aufzählung trägt der durch den Diff geänderten Menge nicht · Abschnitts-Überschrift in die Absatzzeile gezogen · Genus-Fehler in der Nutzer-Doku · Messprotokoll ohne Rohdaten · DoD-Wortlaut trägt der Protokoll-Stand nicht

## Verdikt

**Merge-blockierend:** nein — keine HIGH-Findings. Die drei MEDIUM-Findings sind
Berichtigt-Läufe vor der Closure: F-1 und F-2 tragen die Beleg-Tabelle der
Spezifikation nach (Erweiterung der Schlüssel-Tabelle der Profil-Kombinationen um das
`java-toolchain`-Paar und ein Assert auf die Vertriebs-Form — oder eine schärfere
Beleg-Formulierung auf den cmd-Lauf, der die Profil-Hälfte heute de facto trägt),
F-3 trägt drei §3-Zeilen nach. Alle drei bleiben am Slice: Die Java-Zeilen der
Beleg-Tabelle sind mit dem Umsetzungs-Diff geliefert — ihre Richtigkeit ist damit
Slice-Gegenstand, nicht Rand.

**Übergabe:** Findings gehen an den Implementer (Rückkante
Review → Plan bei Plan-Defekt: F-3); die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext). Die frische
Gegenprobe der Prüfsummen und der java-Tags ist in diesem Report als
Lauf-Evidence festgehalten, nicht als Closure-Beweis der DoD.