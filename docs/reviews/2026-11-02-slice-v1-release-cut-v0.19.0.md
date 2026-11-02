# Review-Report: slice-v1-release-cut-v0.19.0 — 2026-11-02

**Review-Art:** Code (Release-Cut-Nacharbeit: Cut-Commit gegen Plan, Hard Rules und
[`releasing.md`](../maintainer/releasing.md); Lead des Blocks `[0.19.0]` gegen
Lastenheft `0.30.0`, [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) und die ausgelieferten Texte)

**Gegenstand:** [`slice-v1-release-cut-v0.19.0`](../plan/planning/done/slice-v1-release-cut-v0.19.0.md),
Tag `v0.19.0` auf `3564ac31` (gepubliziert: Release `u-boot v0.19.0`, 12
Assets, 2026-11-02 — Probe `gh release view v0.19.0`). Der Cut selbst ist der
Diff des Commits `3564ac31` (7 Dateien); dazu die Cut-Vorbereitung im Range
`9639db21..v0.19.0` (Slice-Anlage `b4d56997`, reine Moves `eff895c2`/`3f223739`,
Verweis-Nachzug `bf0eb01f`) und der Nachzugs-Commit `ec556f96` (Roadmap, 2
Zeilen). Träger: der Umsetzungs-Slice
[`slice-v1-katalog-gepinnte-toolchain`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md)
(separat gereviewt), der CR
[`slice-cr-devcontainer-gepinnte-werkzeugkette`](../plan/planning/done/slice-cr-devcontainer-gepinnte-werkzeugkette.md)
(separat gereviewt) und [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md);
sie sind hier nur am Produkttext und als Eintrags-Belege gelesen.

**Skill:** `.harness/skills/reviewer.md` @ `v0.19.0`
**Modell:** glm-5.3-flash · **Datum:** 2026-11-02

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan [`slice-v1-release-cut-v0.19.0`](../plan/planning/done/slice-v1-release-cut-v0.19.0.md)
  (Plan-Tabelle §3, DoD §2, Risiken §6), [`releasing.md`](../maintainer/releasing.md) §§3–6
- Lastenheft, Stand
  [`0.30.0`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)
  ([`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)
  Absatz zu `go-toolchain`/`golangci-lint`,
  [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)
  Absatz `containerEnv`, [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)
  neu, Historie-Zeile `0.30.0`),
  [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)
  samt Beleg-Tabelle
- [`CHANGELOG.md`](../../CHANGELOG.md) Block `[0.19.0]` samt der drei
  Träger-Einträge und dem Changed-Eintrag; Benutzerhandbuch (Abschnitt
  *Werkzeugketten aus gepinntem Release optional*, Features-Tabelle),
  [`README.md`](../../README.md)/[`README.de.md`](../../README.de.md)
- [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md),
  Muster-Report der Vorgänger-Cuts
  [`2026-11-02-slice-v1-release-cut-v0.18.0.md`](2026-11-02-slice-v1-release-cut-v0.18.0.md)
  (Finding-Klassen als Suchbild),
  Beobachtung [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md)
  (Plan-Risiko 1)
- [`AGENTS.md`](../../AGENTS.md) §§3.7, 3.9, 3.10, 5 Regel 6;
  [`AGENTS.md`](../../AGENTS.md) §3.10 Abgleich Plan-Tabelle gegen
  `git diff --stat 9639db21..3564ac31` und gegen `git diff v0.19.0..HEAD`

**Nicht ausgeführte Sensors** (Auftrag Review-only — keine Gate-Läufe, keine
Commits): `make gates`, `make ci`, `make fullbuild` (beim Hauptagenten grün auf
dem Cut-Stand, Stempel-HEAD nicht nachgelesen — Auftragskontext; CI-Läufe
`ci` `41460180976`, `integration-docker` `41460180946`, `tap-check`
`41463190182`, `publish` `41463337361` je success — Auftragskontext, nicht
eigenständig verifiziert; Probe auf das Release selbst bestanden), 
`make test-docker` (Auftrag), `make verify-depguard` (keine Importregel
berührt — kein Produktcode im Cut), `make docs-check` (Reviewer schreibt den
Report vor dem Gate-Lauf des Verifiers; Link-Tiefen manuell gegen die
Fundstellen geprüft). Verifikation der DoD und der CI-Evidence: Verifier,
nicht dieser Lauf.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Lead-Aussage „…, so the caches survive `u-boot devcontainer up --recreate`“ lässt die Profil-Grenze weg: der Vertrag zugesagt das Überstehen des Neuanlegens **nur im Sandbox-Profil** (Lastenheft `0.30.0`, [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen): „das Überstehen des Neuanlegens nur im Sandbox-Profil auf einem Linux-Host (`amd64`) …, soweit ein Test mit Build-Tag `docker` es belegt“; Handbuch: „Im Sandbox-Profil liegt … Im Standardprofil schreibt `u-boot` keine Aussage zu, wo der Container sein Workspace-Verzeichnis hält“). Der eigene Träger-Eintrag trägt die Grenze („In the sandbox profile these caches … survive the recreation“; Not-promised: „the standard profile for the cache effect“) — der Lead verwirft sie. Die Not-promised-Liste des Lead nennt „other hosts and engines“, aber nicht das Standardprofil für den Cache-Effekt. Ein Nutzer im Standardprofil liest den Lead so, dass seine Caches das `--recreate` überstehen — die Zusage dazu steht nicht. | [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Beobachtung [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Plan-Risiko 1) | `CHANGELOG.md` · “so the caches survive `u-boot devcontainer up --recreate`” (Lead, Zeile 17) | nein — kein Gate liest Lead gegen Vertragsaussage; `make docs-check` prüft nur Links | Lead-Grenzliste kürzer als Vertragsliste (3. Befund — `v0.17.0` F-2 LOW, `v0.18.0` F-2 MEDIUM; Schwelle erreicht, Steering-Loop) |
| F-2 | MEDIUM | Die Lead-Formulierung „the build then reads the official checksum source of the chosen version“ legt den Lesevorgang in den Build. Der Vertrag sagt das Gegenteil: der Lesevorgang liegt bei der **Erzeugung** ([`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv): „Der Generator liest bei der Erzeugung die Prüfsummen-Quelle des Eintrags für die gewählte Version“; Handbuch: „Beim `generate devcontainer` liest `u-boot` die offizielle Prüfsummen-Quelle … der Lauf braucht für eine Übersteuerung dieselbe Internetverbindung wie der Image-Build“). Der Image-Build liest keine Quellen-Angabe — er verifiziert das Archiv gegen die eingetragene Prüfsumme (`sha256sum -c`). Der Träger-Eintrag sagt es richtig („u-boot then reads … and builds the override against the read checksums“); ein Leser des Lead legt den Netzgriff in den falschen Schritt des Ablaufs. | [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Beobachtung [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md) | `CHANGELOG.md` · “the build then reads the official checksum source” (Lead, Zeile 19) | nein — Review-only | Lead-Formulierung gegen die Vertragsaussage lesbar (2. Befund nach `v0.18.0` F-3 — Drift-Signal) |
| F-3 | LOW | Die Lead-Aufzählung der Abweisungen „Without a readable checksum source, an entry without a source, and the combination of the external `go` feature with `go-toolchain` … exit code `10`“ nennt den dritten Abweisungsfall der Vertragsliste nicht: eine **lesbare** Quelle, die die gewählte Version nicht überdeckt, weist der Vertrag ebenfalls mit Exit-Code `10` ab ([`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) Abweisungen: „Ist die Quelle nicht lesbar oder überdeckt sie die gewählte Version nicht …“; Träger-Eintrag: „whose source is unreadable or does not cover the version“). Der Zählwert „the three exit codes“ stimmt, die genannte Einteilung nicht: ein Leser des Lead hält eine lesbare, aber nicht überdeckende Quelle für angenommen. | [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv), Beobachtung [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md) | `CHANGELOG.md` · “Without a readable checksum source, an entry without a source” (Lead, Zeile 19) | nein — Review-only | Lead-Abweisungsliste enger als die Vertrags-Abweisung (1. Befund) |
| F-4 | MEDIUM | Slice §7 trägt „Noch nicht geschlossen; der Slice liegt in `open/` (Plan).“ — der Slice liegt seit `3f223739` in `in-progress/`; der Marker ist abgelaufen. Denselbe Wortlaut trug der `v0.18.0`-Cut (dort F-4, LOW); das Muster wiederholt sich unverändert, obwohl der `v0.18.0`-Report es bereits als Klasse mit der Schwelle geführt hat. | [`AGENTS.md`](../../AGENTS.md) §3.10 (Zustandsmarker mit Ablaufpunkt) | [`slice-v1-release-cut-v0.19.0.md`](../plan/planning/done/slice-v1-release-cut-v0.19.0.md) §7 · “der Slice liegt in `open/` (Plan)” | nein — Review-only | Zustandsmarker nach Ablauf nicht abgebaut (Wortlaut identisch mit `v0.18.0` F-4 — 2. Befund) |
| F-5 | LOW | Die Plan-Tabelle (§3) führt `docs/plan/planning/open/README.md` mit „Bestand-Zeile für diesen Slice“, und die Commit-Message von `bf0eb01f` behauptet „Bestand open/README“ — im ganzen Range `9639db21..v0.19.0` trägt **kein** Commit eine Änderung an `open/README.md`, und `git log -S` findet die Kennung des Cut-Slices dort nie: eine Bestand-Zeile für diesen Slice ist nie angelegt worden. Der DoD-Punkt „Bestand-Zeilen“ (§2) wird damit gegen einen nicht existierenden Stand geprüft, und die Message behauptet eine Änderung, die der Diff nicht trägt. Die Berichtigung (Plan-Zeile streichen oder Bestand-Zeile nachtragen) ist vor dem Abhaken fällig; eine korrigierte Message ist nicht möglich — der Abgleich gehört in die Closure-Notiz. | [`AGENTS.md`](../../AGENTS.md) §3.10 (Abgleich Plan-Tabelle gegen `git diff --stat`; „jeder Eintrag ohne Diff wird berichtigt“), [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md)-Klasse | [`slice-v1-release-cut-v0.19.0.md`](../plan/planning/done/slice-v1-release-cut-v0.19.0.md) §3 · “docs/plan/planning/open/README.md” (Plan-Zeile) | nein — Review-only | Plan-Eintrag ohne Diff; Commit-Message behauptet die Änderung (`BEO-DOK/plan-aussage-ueber-code`, Verwandter) |
| F-6 | LOW | Die Roadmap-Zeile unter *Abgeschlossene Wellen* trägt den Anker-Text “[`slice-v1-release-cut-v0.19.0`](../plan/planning/done/slice-v1-release-cut-v0.19.0.md)” — mit einem überzähligen Apostroph in der Kennung. Die Kennung existiert im ID-Schema nicht ([`MR-000`](../../harness/conventions.md#mr-000--baseline-aussage): `slice-<phase>-<slug>`), alle vorangegangenen Zeilen der Tabelle tragen die nackte Kennung; eine Suche nach dem Anker-Text geht leer. | [`MR-000`](../../harness/conventions.md#mr-000--baseline-aussage) (ID-Schema), [`AGENTS.md`](../../AGENTS.md) §3.7 (Beleg als auflösbarer Anker) | [`roadmap.md`](../plan/planning/in-progress/roadmap.md) · “[slice-v1-release-cut-v0.19.0](../plan/planning/done/slice-v1-release-cut-v0.19.0.md)" (Apostroph berichtigt) (Zeile 117) | nein — Review-only | Beleg-Kennung existiert nicht (Apostroph) |
| F-7 | LOW | Die Planungs-Zeilen mit den Vorschlags-Versionen sind an der vergebenen Version vorbei: die Roadmap führt unter *Nächste Wellen* „`v0.19.0` — Feature-Katalog: gepinnte Werkzeugkette, Quellen, Optionen und Basis-Image (geöffnet als CR-Slice; Umsetzungs-Slice nach dem CR)“ mit Trigger „`v0.18.0` veröffentlicht“ — der Tag `v0.19.0` ist jetzt gesetzt und trägt nur die Werkzeugkette; Quellen, Optionen und Basis-Image stehen in `open/` ohne Zuordnung. Dasselbe gilt für die Bestand-Zeile des JDK-Entwurfs „ohne Welle (Vorschlag `v0.19.0`)“ in [`open/README.md`](../plan/planning/open/README.md) (Zeile 86). Die Muster der `v0.18.0`-Berichtigungsrunde (dort F-8/F-9) wiederholen sich; die nächste Zuordnung dieser Zeilen ist die Folgerversion. | [`AGENTS.md`](../../AGENTS.md) §3.10 (Zustandsmarker), [`releasing.md`](../maintainer/releasing.md) §3 Punkt 4 | [`roadmap.md`](../plan/planning/in-progress/roadmap.md) · “`v0.19.0` — Feature-Katalog: gepinnte Werkzeugkette, Quellen, Optionen und Basis-Image” (Zeile 32); `open/README.md` · “Vorschlag `v0.19.0`” (Zeile 86) | nein — Review-only | Vorschlags-Version vergeben (2. Befund nach `v0.18.0` F-8/F-9) |
| F-8 | INFO | Die Nachzugs-Zustandsmarker sind heute korrekt und laufen erst mit der Closure ab (Urteil zu der Auftragfrage): die Zeile „liegt in `in-progress/` — … `publish` `41463337361` grün mit Job `tap`; Doku-Nachzug und Closure folgen“ nennt Zustand und Beleg als auflösbare Anker (Tag-Hash, Lauf-ID — keine Chronik, Form korrekt); die Zeile unter *Abgeschlossene Wellen* trägt „(liegt in `in-progress/`, nach der Closure `done/`)“ und den Zeiger auf `../in-progress/`. Beide veralten mit der Closure: dann müssen der Klammerzusatz weg, der Run-ID-Satz zum Zustand „geöffnet, Nachzug erledigt“ gesetzt und der Zeiger auf `../done/` berichtigt werden (die Message von `ec556f96` trägt das bereits als Auftrag „Ruheort-Berichtigung des Zeigers folgt mit der Closure“). Keine Regelverletzung zum Zeitpunkt dieses Laufs. | [`AGENTS.md`](../../AGENTS.md) §3.10 | [`roadmap.md`](../plan/planning/in-progress/roadmap.md) · “Doku-Nachzug und Closure folgen” (Zeile 24) | nein — Review-only | Zustandsmarker mit Ablaufpunkt (noch nicht abgelaufen) |
| F-9 | INFO | Die Handbuch-Version bleibt `2.1`, obwohl das Release den Abschnitt *Werkzeugketten* ausliefert (durch den Umsetzungs-Slice, mit eigenem Review). [`releasing.md`](../maintainer/releasing.md) §5 verlangt nur den Abgleich von `Stand` gegen das Tag-Datum (2026-11-02, stimmt) und keine Anhebung der Handbuch-Version; die Inhaltsänderung kam je im Träger-Review. Keine Regel verletzt. | [`releasing.md`](../maintainer/releasing.md) §5 | [`benutzerhandbuch.md`](../user/benutzerhandbuch.md) · “Handbuch-Version: 2.1” | nein — Review-only | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Lead des Blocks `[0.19.0]` gegen Lastenheft `0.30.0` (Satz zu S1–S11) und [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv): zwei Schlüssel `go-toolchain`/`golangci-lint`, Installation aus einem gepinnten Release-Archiv in beiden Profilen, `sha256sum -c` vor dem Entpacken (am Template `Dockerfile.tmpl` gelesen, Zeilen 161/190), Exit `10` ohne Schreiben, Cache-Werte unter `${containerWorkspaceFolder}`, Pin-Politik der vier Stellen mit Test | geprüft, ohne Befund abgesehen von F-1 (Profil-Grenze), F-2 (Lese-Moment), F-3 (Abweisungsliste) |
| Lead „Nineteenth release.“ (Zählung fortlaufend: `0.18.0` „Eighteenth"), Host-Eingrenzung „on a Linux host (`amd64`) with a local, non-nested engine, as far as the tests of this release show — unit tests, Golden Cases … and Docker tests with build tag `docker`" gegen die Vertrags-Bedingung (Erzeugung und Bau nur `amd64`, nur Docker-Test), Not-promised „the working behaviour of the tools themselves in an overridden version, … the lint run without a Go installation" gegen die Beleg-Tabelle der [SPEC-054](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) (der Fall `golangci-lint` allein trägt dort keine Zusage, `TestE2E_Toolchain_GolangciLintAloneBuilds`), „and a limit beyond the measured ones" | geprüft, ohne Befund |
| Drei `[0.19.0]`-Träger-Einträge gegen die Lastenheft-Stände (gelesen, nicht neu bewertet — eigenes Review im Umsetzungs-Slice): Werkzeugketten-Eintrag (beide Schlüssel, Prüfsummen je Architektur, Übersteuerung mit Quellen-Lesung, drei Abweisungen mit Exit `10`, Ausweisung ohne Signatur als `warn` unter `--quiet`), Cache-Eintrag (`GOCACHE`/`GOMODCACHE` je einmal, Sandbox-Grenze), Pin-Politik-Eintrag (eine Stelle, vier Stellen gleich, Test) | geprüft, ohne Befund (die Lead-Abweichungen betreffen nur den Lead) |
| Handbuch-Stelle der Träger gelesen (*Werkzeugketten aus gepinntem Release optional*): Installationsform, `/usr/local`-Entpacken, Version/Pin-Stand `1.27.3`/`v2.14.0`, Übersteuerung mit Quellen-Lesung beim Erzeugen (auch unter `--dry-run`), Fall „nicht überdeckende Quelle" ausdrücklich abgewiesen, `go` + `go-toolchain` abgewiesen, `golangci-lint` allein mit Grenze, Cache-Pfade mit der Sandbox-Grenze im Wortlaut — gegen den Träger-Eintrag und gegen F-1/F-2 konsistent | geprüft, ohne Befund |
| Handbuch-Kopf: `Software-Version` v0.19.0, `Stand: 2026-11-02` = Tag-Datum, drei Versionsbeispiele (`0.19.0`, `download/v0.19.0`, `VERSION=0.19.0`) und Gültigkeitsbereich; Versionsrest `0.18.0` in READMEs, Handbuch, `cmd/`, `repo.mk`, `Dockerfile`, `.github/`, `packaging/`: 0 Treffer (Suchbild); Handbuch-Version (F-9) | geprüft, ohne Befund |
| Versionszeilen `0.19.0-dev` an den drei Stellen (`cmd/uboot/main.go:45`, `repo.mk:66`, `Dockerfile:47`) samt Kommentar-Nennungen (`Dockerfile:44,46`, `repo.mk:60,64`): Kommentare beschreiben den Stand, keine Verweis-auf-Geändertes-Muster (§3.9 a/b); `grep -rIn '0.18.0-dev'`: 0 Treffer | geprüft, ohne Befund |
| Kennungen in READMEs und Handbuch = 0 (`AGENTS.md` §5 Regel 6, `grep`-Abgleich); READMEs tragen nur Status-Zeile und `VERSION=`-Beispiel, kein „released", keine Release-Tabelle | geprüft, ohne Befund |
| Compare-Links und leerer Anker: `## [Unreleased]` leer (nur Kopfzeile), `[Unreleased]: …/compare/v0.19.0...HEAD`, `[0.19.0]: …/compare/v0.18.0...v0.19.0`, `[0.18.0]` unverändert; Block-Datum 2026-11-02 = Tag-Datum | geprüft, ohne Befund |
| Plan-Tabelle (§3) gegen `git show --stat 3564ac31` (7 Dateien: `CHANGELOG.md`, `Dockerfile`, beide READMEs, `cmd/uboot/main.go`, `benutzerhandbuch.md`, `repo.mk`) — jede Datei mit Plan-Zeile; „kein Diff" für `spec/`, `docs/plan/adr/`, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` im Cut-Commit bestätigt; Plan-Zeile `open/README.md` ohne Diff im Range (F-5) | geprüft (F-5) |
| Plan-Tabelle gegen den Nachzugs-Diff `v0.19.0..HEAD` (2 Zeilen `roadmap.md`): Zeile unter *Abgeschlossene Wellen* mit Tag-Hash `3564ac31`, Datum 2026-11-02 und Zustandszeile im Verweis-Nachzug — von den Plan-Zeilen gedeckt | geprüft, ohne Befund (Ablauf-Marker: F-8) |
| Cut-Vorbereitung im Range `9639db21..v0.19.0`: Slice-Anlage `b4d56997` (nur Slice-Datei), `eff895c2`/`3f223739` reine Moves (0 Änderungen, Rename-Erkennung intakt), `bf0eb01f` Verweis-Nachzug (nur `roadmap.md`; Message behauptet mehr als der Diff — F-5), Git-mv-Disziplin (§3.3 Fall 1) eingehalten | geprüft, ohne Befund abgesehen von F-5 |
| Roadmap nach dem Tag: `Letzte Änderung: 2026-11-02` aktuell; Zeile „liegt in `in-progress/`" mit Beleg-Ankern (Tag-Hash, `publish`-Lauf-ID, Job `tap`) — Zustand und Anker, keine Chronik | geprüft, ohne Befund (Ablauf-Marker: F-8) |
| `releasing.md`-Residuen: §3 Punkt 5 nennt die Skelett-Vorlage (Rest des `v0.18.0`-Berichts F-7 ist behoben); „zuletzt" nur als PAT-Hinweis; keine veralteten Versionswerte im Dokument | geprüft, ohne Befund |
| §3.7/§3.9 über beide Diffs (Cut und Nachzug): nur Versionswerte und Kommentar-Nennungen auf den neuen Stand; kein Kommentar über Verworfenenes, kein Verweis auf den alten Wert in Kommentaren, kein Helfer zwischen Doc-Kommentar und Funktion | geprüft, ohne Befund |
| Reviewer-HIGH-Liste (Dual-Classifier, Referenzmodell, Sprachvertrag, Architekturgrenze, Gate-Lockerung, Kommentar-Klassen, Carveout, Immutabilität, Lifecycle): der Cut ändert keinen Produktcode, kein Gate, kein ADR, keine Spec; kein abwärtsgerichteter Verweis aus Rang 1–3; CHANGELOG-Links auf Lastenheft-Anker laufen aufwärts | geprüft, ohne Befund |
| CI-/Release-Evidence: Release `v0.19.0` gepubliziert (Probe `gh release view`: 12 Assets, 2026-11-02) — Rest der Lauf-IDs als Auftragskontext, nicht eigenständig verifiziert | nicht im Ganzen geprüft: Review-only; Nachweis beim Verifier |
| Produktcode und Beleg-Tabelle der Träger (Sentinel-Abdeckung, Docker-Tests, Pin-Stellen-Test) | nicht geprüft im Ganzen: Umsetzungs- und CR-Slices tragen je ein Review; hier nur die für die Lead-Linsen nötigen Stellen (Template, Kommentare) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 (F-1, F-2, F-4) |
| LOW | 3 (F-3, F-5, F-6, F-7) |
| INFO | 2 (F-8, F-9) |

**Finding-Klassen dieses Laufs:** Lead-Grenzliste kürzer als Vertragsliste
(3. Befund — Schwelle, v0.17.0 F-2 / v0.18.0 F-2 / dieser Lauf F-1) ·
Lead-Formulierung gegen die Vertragsaussage lesbar (2. Befund nach `v0.18.0`
F-3) · Lead-Abweisungsliste enger als die Vertrags-Abweisung (1. Befund) ·
Zustandsmarker nach Ablauf nicht abgebaut (Wortlaut identisch mit `v0.18.0`
F-4; 2. Befund dieses Wortlauts) · Plan-Eintrag ohne Diff / Commit-Message
behauptet die Änderung (`BEO-DOK/plan-aussage-ueber-code`, verwandt) ·
Beleg-Kennung existiert nicht (Apostroph) · Vorschlags-Version vergeben
(2. Befund) · Zustandsmarker mit Ablaufpunkt (noch nicht abgelaufen, F-8).

**Steering-Loop-Vermerk:** „Lead-Grenzliste kürzer als Vertragsliste“ erreicht
mit diesem Lauf die Schwelle von drei Vorgängen (`v0.17.0` Cut F-2, `v0.18.0`
Cut F-2, dieser Lauf F-1) — sie gehört in den Lerneintrag der Closure §7 und
zur Architect-Prüfung, ob ein computational Sensor möglich ist (etwa eine
`docs-check`-Heuristik oder ein Checklisten-Schritt „Lead-Grenzen gegen die
Historie-Zeile des Lastenheft-Stands abticken“ vor dem Tag nach
[`releasing.md`](../maintainer/releasing.md) §4). „Lead-Formulierung gegen die
Vertragsaussage lesbar“ steht bei zwei Vorgängen; „Zustandsmarker nach Ablauf
nicht abgebaut“ übersteht die Schwelle des v0.18.0-Laufs bereits als
ungelöste Klasse — der Cut-Wortlaut in §7 wiederholt sie. Das
Bestand-Zeilen-Muster (F-5) ist ein Verwandter der Klasse
`BEO-DOK/plan-aussage-ueber-code` mit einer neuen Facette (Message behauptet
mehr als der Diff); es steht bei 1. Vorgang.

## Verdikt

**Merge-blockierend:** nein — der Tag `v0.19.0` ist gesetzt und veröffentlicht
und wird nicht zurückgenommen ([`releasing.md`](../maintainer/releasing.md)
*Zweck*). F-1 bis F-3 betreffen den Lead des ausgelieferten
[`CHANGELOG.md`](../../CHANGELOG.md)-Blocks und führen zu Nachzugs-Commits auf
`main`; die LOW-Findings (F-5–F-7) betreffen Plan- und Planning-Dokumente und
gehen vor der Closure an den Implementer. Der Slice geht nicht nach `done/`,
bevor der Hauptagent F-1 bis F-3 entschieden hat und F-5/F-6/F-7 (und die
Ablauf-Marker des F-8) berichtigt sind. Die Findings werden nicht
herabgestuft, wenn der Implementer widerspricht; HIGH gibt es nicht, ein
Konfliktpfad über den Architect entfällt.

**DoD-Vermerk:** Nicht Bestandteil dieses Laufs (Verifier). Festgehalten: die
offenen DoD-Punkte (Review-Report — dieser Lauf —, Verification-Report,
Closure-Notiz mit Lerneintrag, `git mv` nach `done/`) sind im Slice ehrlich
offen; der DoD-Punkt „Bestand-Zeilen“ ist durch F-5 gegen den Ist-Stand
(Bestand-Zeile nie angelegt) zu berichtigen, bevor er abgehakt wird; §7 trägt
den abgelaufenen Stand-Marker (F-4) und wird mit der Closure-Notiz ohnehin
neu gefasst. Der Report wurde vor dem Gate-Lauf des Verifiers geschrieben;
seine Links sind manuell gegen die Fundstellen geprüft, nicht durch
`make docs-check`.