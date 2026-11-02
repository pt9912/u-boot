# Verifier-Report: [`slice-v1-katalog-gepinnte-toolchain`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md) — 2026-11-02

**Review-Art:** Verification — DoD- und Entscheidungs-Konformität samt
Plan-vs-Code; Schema und Standardblock nach
[`../../harness/verification.md`](../../harness/verification.md). Kein Self-Review
des Reviewer-Reports und keine Repair-Arbeit: Findings werden geprüft und
berichtet, nicht behoben. Keine Commits in diesem Lauf.

**Gegenstand:** Commit-Range `41667865..b32e9c38` (sieben Commits: `7ae8faea`
Code plus Tests, `37062fdc` Spec- und Vertrag-Nachzug, `58141a2c`
Hebung-Abgleich und eigener Devcontainer, `5872a1fa` Testpaket-Umzug,
`2ed2d956` Plan-Fundstellen-Abgleich, `8ab58b95` Fix-Runde zu F-1 bis F-8,
`b32e9c38` Berichtigungsrunde r2 zu F-N1 bis F-N3); 32 Dateien,
+2729/−91 nach `git diff --stat`.

**Modell:** `glm-5.3-flash:cloud[1m]` · **Datum:** 2026-11-02 ·
**Rang:** `main` = `b32e9c38` (`origin/main` gleich; Push-Tip geprüft).

---

## Verification Evidence

### Scope

- Slice: [`slice-v1-katalog-gepinnte-toolchain`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md)
- IDs: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (S1 bis S11) · [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) · [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) · [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) · [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) · [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) (version-Pfad) · [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Aufzählung) · [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)
- Artefakte: Katalog-Eintragstyp mit den zwei Schlüsseln
  `go-toolchain` und `golangci-lint`, Port `driven.Fetcher` plus Adapter
  `fetch/`, Template-Block und `containerEnv`-Render-Stelle, Docker-Tests,
  Golden-Set, Pin-Stellen-Abgleichstest, Spec-, Handbuch-, README- und
  CHANGELOG-Nachzug, Review-Reports r1 und r2.

### DoD-Abgleich (die drei Liefer-Punkte des Slices, §2)

| DoD-Punkt | Stand | Evidence |
| --- | --- | --- |
| Zeile 1 — Katalog-Regeln und Renderer, S1 bis S9 | erfüllt (Beleg-Tiefe gelesen, siehe Tabelle unten) | [Beleg-Tabelle](../../spec/spezifikation.md) [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) je Zeile mit geliefertem Namen; alle 23 genannten Testfunktionen im Baum (je `grep -rn`), die Kern-Tests im Umfang gelesen |
| Zeile 2 — Erzeugnis und Pin-Stellen | erfüllt, mit Auftrags-Residuum beim Docker-Lauf — siehe Sensors | `TestToolchainPinPlacesEqual` liest die vier Stellen (root Dockerfile, repo.mk, Katalog-Template, eigene `.devcontainer/Dockerfile`) samt Übersteuer-Subtest; Golden-Byte-Test und ohne-Schlüssel-Fall gelesen; der Docker-Lauf ist der von `2ed2d956` (Auftragskontext, siehe Sensors) plus der CI-Lauf auf dem Push-Tip |
| Zeile 3 — Hebung-Abgleich und Spezifikation | erfüllt bis auf Fundstellen-Rest (Finding V-1) und die noch offenen Closure-Pflichten | Pin-Kommentare an beiden Stellen gelesen ([`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), Hebung löst auf [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) Abschnitt Konsequenzen auf); Grabstein [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) in Build-Dateien abwesend; [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) mit Beleg-Tabelle und [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) angepasst; Closure-Notiz, Beobachtungs-Register-Fortschreibung und Risiko-Ausgänge stehen noch aus (Closure-Schritt, nicht Liefer-Punkt) |

Erzeugnis-Stand gelesen: die eigene `.devcontainer/devcontainer.json`
trägt `containerEnv` mit den Pfaden `.cache/go-build` und `.cache/go-mod`
unter `${containerWorkspaceFolder}`; `.devcontainer/Dockerfile` trägt den
verwalteten Block mit den Installationen (115 Zeilen) und den
Hand-written-Rest mit nur `BUILD_MODE`, `GOFLAGS`, `GOTOOLCHAIN` (F6-Entscheidung
aus §7, Zeilen 117 bis 129); `u-boot.yaml` trägt die zwei `enabled`-Einträge.

### Beleg-Tabelle von [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) — Existenz und Tiefe

Existenz jeder der 23 genannten Testfunktionen per `grep -rn "func <Name>("`
über den Baum (hier durchgeführt): alle 23 stehen, je in der vom
Beleg-Auftrag erwarteten Datei — Unit-Tests in
`internal/hexagon/application/generate_toolchain_test.go`,
`devcontainer_features_test.go` und `cmd/uboot/toolchain_test.go`
(plus `cmd/uboot/main_test.go` für `TestRun_Generate_ForgeCLINotices`),
Docker-Tests in `internal/e2e/toolchain_docker_test.go` (fünf
`func TestE2E_Toolchain*`, Zeilen 61, 107, 186, 226, 323),
Adapter-Negativtest in `internal/adapter/driven/fetch/fetch_test.go`.

Tiefe gelesen (fixe Fundstellen am Stand `b32e9c38`):

| Test | Gelesene Pinntiefe |
| --- | --- |
| `TestToolchainPinPlacesEqual` (Zeile 694) | liest die vier Stellen per Datei, vergleicht je Pin auf Gleichheit und gegen die Katalog-Defaults; Subtest „an override changes none of the four places“: nach einem Generate-Lauf mit Übersteuerung werden root Dockerfile, `repo.mk`, Katalog-Template und `u-boot.yaml` erneut gelesen und auf Byte-Gleichheit geprüft, und die erzeugte Datei trägt allein die übersteuerten Werte — die [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)-Zusage samt Übersteuer-Wechselwirkung vollständig |
| `TestFeatureCatalogue_ToolchainChecksumSources` (Zeile 824) | Typgrenze über die Brücke `ToolchainVersionOverrideAcceptedForTest` in vier Synthetic-Fällen plus Positiv-Fall; Last-Pfad über `ValidateDevcontainerFeaturesForTest`: Übersteuerung beider Schlüssel angenommen, `gh` ohne Quelle abgewiesen — die F-1- und F-2-Behebung ist am Stand sichtbar |
| `TestToolchainChecksumParsing` (Zeile 865, Abschnitt url forms 911) | Parser beider Quellen gegen Festwerte, nicht-überdeckt- und malformed-Fälle, und die URL-Formen: `https://go.dev/dl/?mode=json&include=all` und die Release-Checksums-URL `releases/download/v2.13.2/golangci-lint-2.13.2-checksums.txt` — dieselbe Form, die der Docker-Test der Übersteuerung am echten Netz liest |
| `TestGenerateDevcontainer_Toolchain_OverrideSpellings` (Zeile 927) | trägt die abweichenden Schreibweisen (`go`-Präfix für Go, tag-los für golangci-lint) in der Konfiguration, die erzeugten Blöcke in der normalisierten ARG-Form, Assert `ARG GOLANGCI_LINT_VERSION=v2.13.2` — die F-3-Behebung |
| `TestGenerateDevcontainer_Toolchain_JSONGolden` (Zeile 322) | Byte-Golden der erzeugten `devcontainer.json` mit `go-toolchain`; ohne die Schlüssel bleibt die Datei ohne `containerEnv` und die Dockerfile ohne die ARGs; mit beiden Schlüsseln erscheinen `GOCACHE` und `GOMODCACHE` je genau einmal (Zähl-Assert) |
| `TestGenerateDevcontainer_Toolchain_OverrideRendersReadChecksums` (Zeile 391) | gelesene Prüfsummen im Block, Double-Run re-liest dieselbe Quelle und bleibt NoOp |
| `TestGenerateDevcontainer_Toolchain_Notices` (Zeile 540) und `TestRun_Generate_ToolchainNotices` (`cmd/uboot/toolchain_test.go` Zeile 18) | je Schlüssel eine Meldung, Code [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), Level `warn`, ohne `SecurityRelaxation` und `DegradationFallback`; die CLI-Variante prüft beide Profile |
| `TestGet_ResponseOverTheLimitIsAnError` (`internal/adapter/driven/fetch/fetch_test.go` Zeile 91, F-6) | chunked Server ohne Content-Length streamt `64<<20+1` Bytes von einem Zero-Reader; Assert auf den Fehler samt Kappen-Wortlaut; der Helfer `unitReader` steht nach der Funktion (§3.9 (c)) |
| `TestDoctor_ToolchainKeys_NoFindings` (Zeile 793) | mit beiden aktivierten Schlüsseln und in beiden Zuständen (ohne `features`-Schlüssel, leeres `features`) tritt keine `allowlist`- und keine `drift`-Prüfstelle aus |

Die fünf Docker-Case-Namen (`TestE2E_Toolchain_BuildsAndReportsPins`,
`…PinMismatchBreaksBuild`, `…OverrideBuildsAgainstReadChecksums`,
`…CacheSurvivesRecreate`, `…GolangciLintAloneBuilds`) sind im Baum
existenzgeprüft; ihr Lauf ist der Docker-Nachweis (siehe Sensors).

### Sensors

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| Stempel gegen HEAD | pass — deckend | `.harness/state/gates-passed.head` = `b32e9c383acd892b…` = `git rev-parse HEAD`; `gates-passed.diffsha` = `0449717295248915…` = frisch berechneter Arbeitsbaum-Hash via `tools/harness/working-tree-hash.sh` (Baum clean) — der Stempel deckt den Push-Tip samt Arbeitsbaum, ein wiederholter `make gates` ist entbehrlich |
| CI auf dem Push-Tip (`ci`-Workflow) | pass | `gh run list`: Lauf `41452705877` auf `b32e9c38` `conclusion: success` (2026-11-02T13:24:28Z, 2 min 6 s) |
| `integration-docker` auf dem Push-Tip | pass | Lauf `41452705759` auf `b32e9c38`: `conclusion: success`, von 2026-11-02T13:24:28Z bis 2026-11-02T13:43:30Z — 19 min 2 s, mit klarem Abstand unter der 30-min-Jobgrenze (erster CI-Lauf überhaupt mit den fünf Toolchain-E2E-Tests — alle sieben Commits kamen mit diesem einen Push) |
| `make test-docker` | not run — gedeckt, Begründung unten | siehe Zeile unten |
| `make verify-depguard` | not run, mit Grund | siehe Zeile unten |

`make test-docker` nicht wiederholt — Begründung (Auftrag 3): der Abgleich `git diff --stat 2ed2d956..HEAD` über `internal/e2e` und
`internal/hexagon/application/templates` liefert nur einen Treffer:
`internal/adapter/driven/fetch/fetch_test.go` (+27) — ein Unit-Test ohne
Build-Tag `docker` (gelesen: Paket `fetch_test`, `httptest` auf
`127.0.0.1`, kein `//go:build docker`); `internal/e2e` und
`internal/hexagon/application/templates` sind leer (keine Treffer) — die
fünf E2E-Tests und die Templates sind seit dem berichteten grünen Lauf
(`2ed2d956`, 25 Pakete ok, 1697,8 s) byte-unverändert; die Änderungen an
`Dockerfile` und `repo.mk` seit `2ed2d956` sind Kommentar-only (gelesen:
Pin-Politik-Kommentar und [ADR-0041-Verweis](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)s, keine `ARG`-Zeile berührt).
Der Lauf von `2ed2d956` trägt damit für den Codepfad; der eigentliche
Nachweis für den Tip ist der CI-Lauf `integration-docker` (Zeile oben).

`make verify-depguard` nicht ausgeführt — die Importregeln sind in der Range
unverändert (`git diff --stat`: kein Treffer für `.golangci.yml`, `Makefile`,
`.d-check.yml`); der neue Adapter `fetch/` liegt im Driven-Stratum, wo
`net/http` nach der Sicht (neue Zeile in [`../../spec/architecture.md`](../../spec/architecture.md))
der einzige verzeichnete Ort ist; die CI des Tips führt `make lint` mit
grünem Ergebnis. Der on-demand-Sensor prüft Regeln, die dieser Diff nicht
berührt.

### Plan-vs-Code (Auftrag 4 — `git diff --stat 41667865..b32e9c38` gegen die Plan-Tabelle §3)

Alle 32 Dateien des Diffs tragen eine Zeile in der Plan-Tabelle (gelesen am
Stand nach `b32e9c38`):

- Verdrahtung und Träger [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features): `internal/hexagon/application/devcontainer_features.go` (Zeile Katalog-Einträge, update), `devcontainer_toolchain.go` (neu, Zeile „der Träger des Typs“), `generate.go`, `templates.go`, `devcontainer_sandbox.go` (Zeile „Fundstellen beim Umsetzen gelesen“), Port `internal/hexagon/port/driven/fetch.go` und Adapter `internal/adapter/driven/fetch/fetch.go` (Zeile Driven-Port), `cmd/uboot/main.go` (Zeile Verdrahtung) — geprüft, deckungsgleich.
- Tests: `generate_toolchain_test.go`, `fakes_test.go`, `export_test.go` (Zeile Unit-/Integrationstest), `devcontainer_features_test.go` (Zeile Katalog-Invarianten), `cmd/uboot/toolchain_test.go` (Zeile cmd-Datei), `internal/e2e/toolchain_docker_test.go` (Zeile Docker-Test mit der Testzahl fünf, seit der Fix-Runde berichtigt), Goldens und Pin-Stellen-Abgleich nach den Plan-Zeilen.
- Öffentliche Verträge: `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` (Zeile der öffentlichen Verträge); `harness/replay.md` (Zeile Coverage-Tabelle).
- Spec und Planung: `spec/spezifikation.md`, `spec/architecture.md` (Zeile Spec), die Slice-Datei selbst, `docs/reviews/` mit den drei Reports (Erst-Report, r2, dieser Verifier-Report) nach der berichtigten Zeile.
- Erzeugnis: `.devcontainer/Dockerfile`, `.devcontainer/devcontainer.json`, `u-boot.yaml` (Zeile Erzeugnis).
- Build-Dateien: `Dockerfile`, `repo.mk` (Zeile Pin-Politik-Kommentar, Pins unverändert).
- Gegenprobe der drei Zeilen „unverändert“: `spec/lastenheft.md`, alle ADRs inkl. [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md), `docs/plan/carveouts/README.md`, `.golangci.yml`, `Makefile`, `.d-check.yml`, `internal/adapter/driving/**` — je leerer Treffer in `git diff --stat 41667865..b32e9c38` (hier durchgeführt).

### Traceability

| ID | Beleg |
| --- | --- |
| [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | Beleg-Tabelle von [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) samt Docker-Tests; Konflikt- und Lesefälle über die exit-10-Klasse [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) |
| [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) | `containerEnv`-Golden (`GOCACHE`/`GOMODCACHE` je einmal, Unterpfade `.cache/go-build` und `.cache/go-mod`) und `TestE2E_Toolchain_CacheSurvivesRecreate` |
| [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) | `TestToolchainPinPlacesEqual` (vier Stellen plus Übersteuer-Wechselwirkung); Pin-Politik-Kommentare lösen die Hebung-Routine auf [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) §Konsequenzen auf |
| [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | die Lieferketten-Meldungen (Code im `WarningEntry`, Tests `TestGenerateDevcontainer_Toolchain_Notices` und `TestRun_Generate_ToolchainNotices`) |
| [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) | die version-Pfad-Zeile nennt die zwei Schlüssel (gelesen, Spezifikation) |
| [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) | die Aufzählung der ungekennzeichneten Meldungen trägt die zwei Schlüssel (F-8, gelesen am Fix-Stand) |

### Carveouts

- Neu: keine — `carveouts.md` unberührt in der Range.
- Geloest: keine.
- Unveraendert: keine Berührung.

### Nicht ausgefuehrt

- `make gates` — nicht wiederholt: der Stempel deckt den Push-Tip samt
  Arbeitsbaum (Zeile Stempel oben); ein zweiter Lauf im selben Modus ist
  nach dem Verfahren kein zusätzlicher Nachweis.
- `make test-docker` — nicht wiederholt, Begründung im Sensors-Block; die
  Docker-Nachweise trägt der Lauf von `2ed2d956` (Codepfad unverändert)
  plus der CI-Lauf auf dem Tip.
- `make verify-depguard` — Importregeln unberührt, Begründung im
  Sensors-Block.

### Commit / Artefakt

- Push-Tip `b32e9c383acd892bfbc717d58cc4e0e1250f8256` (`origin/main` gleich).
- Die erzeugten Projektdateien (`.devcontainer/Dockerfile`,
  `.devcontainer/devcontainer.json`) sind als Erzeugnis committet; der
  zweite `generate`-Lauf ist NoOp (gelesen am Test, nicht hier ausgeführt).

---

## Review-Findings — Stand am Fix-Stand `b32e9c38` (Auftrag 2)

Geprüft am Baum, nicht an den Berichts-Zeilen allein.

| Finding | Verifiziert | Evidence |
| --- | --- | --- |
| r1 F-1 (Beleg-Namen ohne Test) | behoben | `TestFeatureCatalogue_ToolchainChecksumSources` (Zeile 824) und `TestToolchainChecksumParsing` mit Abschnitt url forms (Zeilen 865, 911) existieren; die Beleg-Tabelle zitiert die gelieferten Namen (Existenzprüfung über alle 23 Namen, siehe oben) |
| r1 F-2 (Typgrenze, Brücken ohne Konsumenten) | behoben | alle drei Brücken werden gerufen: `GoChecksumsForTest` und `GolangciChecksumsForTest` im Parser-Test, `ToolchainVersionOverrideAcceptedForTest` in vier Fällen plus Last-Pfad über `ValidateDevcontainerFeaturesForTest` (gelesen, Zeilen 824 bis 858) |
| r1 F-3 (Schreibweisen-Zusage ohne Test) | behoben | `TestGenerateDevcontainer_Toolchain_OverrideSpellings` (Zeile 927) trägt `go`-Präfix und tag-lose Schreibweise und asserted die ARG-Form; der Parser-Test pinnt dieselbe Normalisierung in beiden Schreibweisen |
| r1 F-4 (Plan-Tabelle unausgeglichen) | behoben bis auf den verbliebenen Fundstellen-Rest — siehe V-1 | die zwei fehlenden Plan-Zeilen (`devcontainer_toolchain.go`, `devcontainer_features_test.go`) stehen; die driving-Zeile steht auf unverändert mit leerem `git diff --stat`-Treffer (hier nachgeprüft); die e2e-Zeile nennt fünf Tests (fünf `func TestE2E_Toolchain*` gezählt) |
| r1 F-5 (Coverage-Tabelle ohne Bestands-Zeile) | behoben | [`../../harness/replay.md`](../../harness/replay.md) trägt die Zeile der zwei Schlüssel mit Golden-Vertrag und Evidence-Pfaden; alle dort genannten Testnamen existieren (geprüft) |
| r1 F-6 (Kappe ohne Negativtest) | behoben | `TestGet_ResponseOverTheLimitIsAnError` (Zeile 91) am chunked Server, gelesen |
| r1 F-7 (Hebung-Aussage ohne auflösbare Herkunft) | behoben | der Pin-Politik-Kommentar des root Dockerfile (Zeilen 23 bis 30 gelesen) und von `repo.mk` (Zeilen 47 bis 54 gelesen) nennen die Hebung-Routine mit [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) Abschnitt Konsequenzen |
| r1 F-8 (Aufzählung in [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)) | behoben | die Aufzählung nennt die Meldung je aktivem Schlüssel mit [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)-Verweis (Spezifikation, gelesen am Fix-Stand) |
| r1 F-9 (Laufzeit an der Jobgrenze) | entkräftet | der erste CI-Lauf mit den fünf E2E-Tests (Lauf `41452705759`) endete nach 19 min 2 s bei einer Jobgrenze von 30 min ([`../../integration.yml`](../../.github/workflows/integration.yml) Zeile 37) — der berichtete lokale Lauf von 1697,8 s war host-slow; die gemessene CI-Dauer trägt die Grenze mit klarem Abstand |
| r1 F-10 (Handbuch-Aussage außerhalb der Code-Belege) | offen — Übergabe an den Validator | der Satz „Lint-Läufe brauchen für viele Prüfungen eine Go-Installation“ (Handbuch, gelesen) ist eine Fremdwerkzeug-Eigenschaft ohne Test; [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) sagt zu diesem Fall bewusst keine Zusage (gelesen, Absatz unten der Beleg-Tabelle) |
| r2 F-N1 | mit Rest behoben — siehe V-2 | die unverifizierbare [`SPEC-034`](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)-Behauptung steht an der Plan-Zeile nicht mehr; die Ersatzbehauptung (Bericht-Zeile 62 trägt die berichtigte Fundstelle der Architektur-Linse) ist am Blatt prüfbar: Zeile 60 bis 62 des Erst-Reports nennt die Fundstellen `internal/hexagon/port/driven/fetch.go`, `internal/adapter/driven/fetch/` und `.golangci.yml` — nachgelesen, Fundstelle passt |
| r2 F-N2 | zum Teil behoben — Rest siehe V-1 | die drei berichtigten Fundstellen habe ich am Original nachgelesen: root Dockerfile — Pin-Politik-Kommentar Zeilen 23 bis 32, `ARG GO_VERSION=` Zeile 41, `ARG GOLANGCI_LINT_VERSION=` Zeile 42; `repo.mk` — Hebung-Kommentar Zeilen 47 bis 54, `GO_VERSION            ?=` Zeile 55, `GOLANGCI_LINT_VERSION ?=` Zeile 56. Die §3-Preamble (23 bis 32, 47 bis 54, ?= bei 55/56) und der Plan-Tabellen-Adapter (23 bis 42, 47 bis 56) und die DoD-Zeile 2 (41/42, 55/56) passen am Original |
| r2 F-N3 (doppelter Checkbox-Marker) | behoben | die drei DoD-Zeilen tragen je einen Marker (gelesen) |
| r2 F-N4 (Beleg-Name in Zeile mit anderem Gegenstand) | unverändert, INFO ohne Blockierwirkung | die Zeile der Profil-Wahl listet drei Namen; die Profil-Wahl tragen `TestFeatureCatalogue_KeysCoverSpecExamples`, `TestFeatureCatalogue_EntriesAreValid` und `TestGenerateDevcontainer_Toolchain_BlocksInBothProfiles` (gelesen: je Profil und Schlüssel-Gruppe ein Subtest mit Block-Assert, kein features-Schlüssel, `USER root` genau einmal, Double-Run NoOp) |

## Neue Verifier-Findings

| ID | Kategorie | Befund | Path |
| --- | --- | --- | --- |
| V-1 | MEDIUM (nur Planning-Baum) | Der §3.10-Abgleich der Plan-Fundstellen trägt am Abgangs-Stand `b32e9c38` noch an zwei Stellen nicht — dieselbe Klasse wie F-4/F-N2 (Beobachtung [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md), weiterer Auftritt): (a) DoD-Zeile 1 sagt, die bestehende Abweisungsregel stehe an „Fundstellen Zeilen 798 bis 814“ — an diesen Zeilen steht heute die Prüfung der Quellen-Freigabe (`featureSourceViolation`); die erwähnten Symbole `featureVersionNeedsSource` (Aufruf 762 bis 764, Regel-Kommentar 829, Funktion 840) und `featureWithoutVersionReason` (857) liegen weiter unten in `internal/hexagon/application/devcontainer_features.go` (je nachgelesen). (b) DoD-Zeile 3 nennt den Pin-Politik-Kommentar „Root-Dockerfile (Zeilen 23 bis 26)“ und „repo.mk (Zeilen 47 bis 48)“ — der Kommentar umfasst 23 bis 32 beziehungsweise 47 bis 54 (nachgelesen, siehe F-N2 oben). Ferner lösen die Vorher-Fundstellen `.devcontainer/Dockerfile` „Zeilen 75 bis 120“, „109 bis 112“, „113 bis 118“ und „84 bis 85“ (§3-Preamble, DoD-Zeile 3) an der regenerierten Datei (129 Zeilen, verwalteter Block 1 bis 116, Hand-Rest 117 bis 129) nicht mehr — sie beschreiben den Vorher-Stand und sind als Lesungs-Anker ausgewiesen, nicht als Ist-Säulen; der Rest unter (a) und (b) ist dagegen als Ist-Stand formuliert. | [`../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md) §2 |
| V-2 | LOW (nur Planning-Baum) | Die F-N1-Berichtigung der Plan-Zeile ist inhaltlich prüfbar gestellt (siehe F-N1 oben), aber der Ersatzsatz ist syntaktisch gebrochen: „(Anchor-Berichtigung in der Fix-Runde am Fundstellen-Verweis der Linse 3; Bericht-Zeile 62 trägt die berichtigte Fundstelle der Architektur-Linse) statt auf die Sicht)“ — der Klammer-Schweif „statt auf die Sicht)" ist ein Waisenfragment der alten Fassung. Kein Gate fängt das ( [`make docs-check`](../../harness/sensors/coverage-gate.md)-Familie — `docs-check` war im CI-Lauf des Tips grün ). | [`../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md) §3, Zeile der Tabelle |
| V-F9 | INFO — entkräftet | die Laufzeit-Schwellen-Frage von r1 F-9 ist durch die Messung des CI-Laufs beantwortet: 19 min 2 s bei 30 min Grenze; kein verbleibender Befund — das Risiko bleibt künftigen E2E-Erweiterungen gewarnt | [`.github/workflows/integration.yml`](../../.github/workflows/integration.yml) |

## Replay- und Golden-Abgleich (Auftrag 5)

Die Coverage-Zeile der zwei Schlüssel in
[`../../harness/replay.md`](../../harness/replay.md) (Abschnitt Current
Generator Golden Coverage, letzte Tabellen-Blockzeile) habe ich gegen den
Baum gelesen: der Golden-Vertrag nennt Byte-Golden der `devcontainer.json`
mit den Unterpfaden `.cache/go-build` und `.cache/go-mod` (im Test als
Fest-Strings gepinnt), ohne-Schlüssel-Byte-Gleichheit, Abweisungen je Exit
10 ohne Write, Übersteuer- und Schreibweisen-Fälle, Pin-Stellen-Gleichheit,
Egress-Gleichheit, doctor-Freiheit und Double-Run NoOp — und die genannten
Evidence-Namen (`generate_toolchain_test.go`, `devcontainer_features_test.go`,
`cmd/uboot/toolchain_test.go`, `toolchain_docker_test.go`) stehen alle im
Baum mit den zitierten Funktionen. Der Fresh-State-/NoOp-Vertrag ist
gelesen, der Name-zuerst-Kontrast des Replay-Blocks erfüllt.

## Rest-Risiken nach dem Lauf

- **F-9 (Laufzeit):** durch die Messung entkräftet — der CI-Lauf mit den
  fünf Toolchain-E2E-Tests endete nach 19 min 2 s (Grenze 30 min). Die
  Warnung bleibt für künftige E2E-Erweiterungen stehen: die Marge entsteht
  heute aus dem Vergleich lokaler Host- und CI-Geschwindigkeit, nicht aus
  einer Grenze im Test.
- **F-10 an den Validator:** der Handbuch-Satz zur Go-Installation bleibt
  eine Aussage über ein Fremdwerkzeug ohne Test; falls er zur Nutzerzusage
  aufgewertet wird, braucht er einen Beleg-Auftrag nach dem Muster der
  Beleg-Tabelle.
- Der Docker-Test der Lesung `TestE2E_Toolchain_GolangciLintAloneBuilds`
  loggt die Fremdwerkzeug-Eigenschaft nur (nicht zugesagt) — Grenze bleibt
  dokumentiert, keine Handlung hier.
- Die noch offene Closure (Closure-Notiz §7 über die F6-Vorentscheidung
  hinaus, Steering-Loop-Lerneintrag, Beobachtungs-Register-Fortschreibung
  mit den Finding-Klassen der beiden Review-Läufe, Risiko-Ausgänge §6 und
  die drei Paarungen) sind Closure-Schritte, keine offenen Liefer-Punkte
  der DoD-Liste.

## Verdikt

**Closure-fähig nach einer kleinen Berichtigungsrunde im Planning-Baum.**
Der Code- und Spec-Stand erfüllt die drei Liefer-Punkte: die Beleg-Tabelle
steht auf 23 existierenden Tests mit gelesener Pinntiefe, alle 8 MEDIUM /
LOW-Findings der r1- und die r2-Findings sind am Fix-Stand behoben oder
als Rest (V-1, V-2, F-9, F-10) sauber ausgewiesen, der Gate-Stempel deckt
den Push-Tip, und alle 32 Diffs-Dateien tragen Plan-Zeilen mit leerer
Trefferliste für die „unverändert“-Zeilen. Vor der Closure:

1. V-1: die zwei stale Fundstellen in DoD-Zeile 1 und DoD-Zeile 3 auf den
   Nachweis-Stand berichtigen (kleine Schreib-Änderung, §3.10-Abgleich
   gegen `git diff --stat` danach wiederholen); V-2 (Satz-Waisenfragment)
   mit derselben Runde mitziehen.
2. Nach der Berichtigung den Stop-Hook-Stempel neu schreiben (neuer HEAD).
3. Closure nach §5 des Slices: Closure-Notiz, Lerneintrag,
   Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen.

## Nachtrag — gemessene CI-Laufzeit integration-docker (F-9)

Gemessen (Messung der Koordination, Lauf-IDs und Zeitsperren über
`gh run list` belegt): Lauf `41452705759` auf dem Tip `b32e9c38` ist
`completed success`, von 2026-11-02T13:24:28Z bis 2026-11-02T13:43:30Z =
19 min 2 s. Grenze: `timeout-minutes: 30` in `integration.yml`. Ausgang:
das F-9-Zeit-Risiko ist entkräftet — Restmarge rund 11 min; der berichtete
lokale Lauf von 1697,8 s war host-slow, die CI-Jobgrenze trägt. Der
zweite Lauf des Tips, `ci` (`41452705877`), ist ebenfalls success.