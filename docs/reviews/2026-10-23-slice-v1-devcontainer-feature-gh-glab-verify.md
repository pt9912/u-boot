# Verify-Report: slice-v1-devcontainer-feature-gh-glab — 2026-10-23

**Rolle:** Verifier (Modul 11, [`harness/verification.md`](../../harness/verification.md)), frischer Kontext. Frage: Wurde der Slice richtig gebaut, gemessen an Plan, [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) (`Accepted`), Lastenheft `0.22.0` und DoD?

**Gegenstand:** Commit-Range `db989760..HEAD` (HEAD `604dea7f`), Slice [`slice-v1-devcontainer-feature-gh-glab`](../plan/planning/done/slice-v1-devcontainer-feature-gh-glab.md), Welle [`welle-forge-cli-gh-glab`](../plan/planning/done/welle-forge-cli-gh-glab.md), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), Review [`2026-10-23-slice-v1-devcontainer-feature-gh-glab.md`](2026-10-23-slice-v1-devcontainer-feature-gh-glab.md).

## Urteil

**Schließbar**, sobald der Planner die offenen Closure-Pflichten der DoD erledigt (siehe DoD-Abgleich, Zeilen „offen“). Kein Befund, der den Code- oder Test-Haken verhindert. Zwei LOW-/INFO-Befunde (V-1, V-2), keiner blockierend.

## Selbst gefahrene Sensoren (Belege, nicht Behauptungen)

Host: Linux 6.8.0-147 `x86_64`, Docker-Image `u-boot:test-docker-tools`, Netz zu `github.com` und `gitlab.com` vorhanden. Arbeitsbaum vor dem ersten Lauf: `git status --short` leer.

| Sensor | Ergebnis |
|---|---|
| `make gates` | **Exit 0** (lint, test, `coverage-gate: OK — coverage 92.00% meets threshold 90%`, `docs-check`, `baseline-verify: v6.18.1 OK — 54 Dateien`, `ablage-pruefung: 26 Faelle, 0 Fehlschlaege`) |
| `TestE2E_ForgeCLI_BuildsInBothProfiles` (`-tags docker`, Optionen wie `make test-docker`) | PASS, 63,4 s (Standardprofil 26,6 s, Sandbox-Profil 36,8 s) — **vom Reviewer nicht gefahren** |
| `TestE2E_ForgeCLI_EgressHosts` | PASS, 52,4 s — **vom Reviewer nicht gefahren** |
| `TestE2E_ForgeCLI_PinMismatchBreaksBuild` (GH, GLAB) | PASS, 15,7 s; Container-Zahl vorher/nachher 4/4 |
| `TestE2E_ForgeCLI_TokenReferencesCarryNoValue` | PASS, 15,3 s |
| `make verify-depguard` | nicht gelaufen: keine Importregel berührt (`git diff --stat` zeigt keine neuen Importe, `.golangci.yml` unverändert) |
| vollständiges `make test-docker` | nicht gelaufen: Paket-Zeitlimit bekannt ([`BEO-E2E/e2e-paket-am-zeitlimit`](../plan/planning/observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md)); gezielt gefahren wie im Plan zugesagt |

## Mutationen (selbst gewählt, per `git apply -R` zurückgenommen)

Vor jeder Mutation Baum sauber, nach jeder Rücknahme `git status --short` leer (Rücknahme Exit 0). Unit-Läufe: `go test -count=1 ./internal/hexagon/application/ ./cmd/uboot/` im Image; Docker-Läufe gezielt. Alle roten Läufe Exit 1. Die Mutation wurde mit `sed` gesetzt und als Patch gesichert (Abweichung wie in Plan §7 F-3 des Reviews).

| # | Bedingung | Mutation | Rot (Unit) | Rot (Docker) |
|---|---|---|---|---|
| M1 | B1 | `Dockerfile.tmpl`: `or .UID .Sandbox .Gh .Glab` → `or .UID .Sandbox` | `…BlocksInBothProfiles`, `…UserRootOnlyWithReason`, `…ActivateAndDeactivateLater`, `TestRun_Generate_ForgeCLINotices` | `BuildsInBothProfiles/default_profile` (`gh: command not found`), Sandbox-Unterfall grün |
| M2 | B5 | Attestierungs-Passage aus der `gh`-Meldung gestrichen | `TestGenerateDevcontainer_ForgeCLI_Notices` | — |
| M2c | B5 | `SecurityRelaxation: true` an der ersten Meldung von `forgeCLIWarnings` | `…ForgeCLI_Notices`, `TestRun_Generate_ForgeCLINotices` (`--quiet`) | — |
| M3 | B9 | `GH_TOKEN` aus `agentCredentialEnv` | `…RemoteEnvReferences`, `TestDoctor_SandboxCredentials_ForgeCLI`, `…ForgeCLIRemedyRestoresReference` | `TokenReferencesCarryNoValue` (Sitzung sieht nur den `glpat-`-Wert, Referenz wirkt nicht) |
| M4 | B6, Host `gitlab.com` | Zeile `glab` aus `egressFeatureHosts()` | `…ForgeCLI_EgressHosts`, `…TogetherWithOllamaAndRuntime`, `TestRun_EgressList_ForgeCLIHostOrigins` | `EgressHosts` (Liste ohne `gitlab.com`, `glab` erreicht den Host nicht) |
| M5 | Eigentümer im Block | `--no-same-owner` aus dem `gh`-`tar` | vier `ForgeCLI`-Golden Cases | `BuildsInBothProfiles` beide Profile (`OWNER:UNKNOWN root`) |
| M6 | Eigentümer im Block | `chown vscode:vscode /usr/local/bin/gh` im `gh`-Block | **grün (keine Wirkung)** | `BuildsInBothProfiles` beide Profile (`OWNER:vscode root`) |
| M7 | B3 | `sha256sum -c -` aus dem `gh`-Block | vier `ForgeCLI`-Golden Cases | `PinMismatchBreaksBuild/GH` (Bau gelingt trotz Nullprüfsumme), `GLAB` grün |

Eine erste Fassung von M6 traf versehentlich einen fremden Block (Zeile 43) und blieb grün; sie wurde verworfen und gezielt auf den `gh`-Block (Zeile 114) wiederholt. Eine Zwischenfassung von M3 trug zusätzlich eine übrig gebliebene Mutation (M2b); der Lauf wurde mit sauber gesetzter Mutation wiederholt (Tabelle zeigt nur den sauberen).

## Befunde

### V-1 — LOW — Eigentümer `root` des Binaries ist nur vom Docker-Test gepinnt, nicht vom Golden Case

Mutation M6 (`chown vscode:vscode` im `gh`-Block) färbt keinen Unit-/Golden-Test, nur `BuildsInBothProfiles`. Die Zusage „Binary-Eigentümer `root`“ (DoD Punkt 2) hängt damit allein an einem Test, der nicht in `make gates` läuft. Das Dokument des Plans verlangt für Punkt 2 den Docker-Test, nicht ein Golden; der DoD-Haken ist getragen. Als Schärfung wäre eine Prüfung „kein `chown` im Block“ im Golden Case billig. Nicht blockierend.

### V-2 — INFO — Pin `glab` `arm64` im Messprotokoll nicht als eigene Zeile belegt

Die `amd64`-Werte beider Werkzeuge und der `arm64`-Wert von `gh` stehen in [`rohdaten-03.txt`](../messungen/slice-v1-forge-cli-quellen-messung/rohdaten-03.txt). Der `arm64`-Wert von `glab` (`0716395a…8d16`) steht dort nicht als eigene Zeile, sondern in der ADR, den Handvarianten und im Bau-Log `rohdaten-06-glabtgz.log`. Ich habe ihn gegen die live abgerufenen Prüfsummen-Dateien der Releases nachgelesen (siehe unten): gleich. Keine Wirkung.

## Prüfung gegen Plan, ADR und Lastenheft

- **Pin-Werte** (`Dockerfile.tmpl` `GH_VERSION=2.103.0`, `GLAB_VERSION=1.124.0`, vier SHA-256): zeichengleich mit `tools/messung/forge-cli/block-ghtgz.dockerfile`, `block-glabtgz.dockerfile`, [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) Entscheidung 1 und — selbst abgerufen — mit `…/cli/releases/download/v2.103.0/gh_2.103.0_checksums.txt` sowie `https://gitlab.com/gitlab-org/cli/-/releases/v1.124.0/downloads/checksums.txt` (alle vier `linux_amd64`/`linux_arm64`-Archive). Die Rohdaten-Zeilen 8, 10 (gh) und 78–81 (glab amd64) von `rohdaten-03.txt` stimmen ebenfalls.
- **Byte-Identität bestehender Konfigurationen:** Im Template ist genau eine Zeile geändert (die `if`-Bedingung, nur um `.Gh .Glab` erweitert); die Blöcke liegen innerhalb ihrer eigenen `{{- if}}`. Kein Golden-/Testdata-Bestand im Diff (`git diff --stat -- '*testdata*'` leer, keine gelöschte Datei); in den vorhandenen Tests ist nur die Katalog-Schlüsselliste, ein Kommentar und ein Helfer (`ollamaSandboxDockerfile` → `blockSandboxDockerfile(t, key)`, gleicher Inhalt für `ollama`) geändert — keine Erwartung eines vorhandenen Golden Case wurde angefasst. `make gates` grün.
- **Entscheidungen 1–8:** (2) Einträge ohne `source`/`defaultVersion`/`publisher`, `dockerfileBlock: true`; (3) `Gh`/`Glab` vor dem frühen `return` des Standardprofils gesetzt, `AgentEnv` danach (nur Sandbox); (4) `egressFeatureHosts` genau `gh` → `api.github.com`, `github.com`, `glab` → `gitlab.com`; (5) `GH_TOKEN`/`GITLAB_TOKEN` als `${localEnv:…}` über `agentCredentialEnv`, feste Reihenfolge; (6) `forgeCLIWarnings` eine Meldung je aktivem Feature, [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), `warn`, ohne `SecurityRelaxation`/`DegradationFallback` (`--quiet` unterdrückt sie, M2c); (7) keine Telemetrie-Voreinstellung (kein `ENV`/`containerEnv`); (8) kein Pin-Hebungs-Code. Keine Abweichung.
- **Tokenwerte:** Im Diff erscheinen nur erkennbar erfundene Werte (`ghp_forgeCANARYcanary…`, `glpat-forgeCANARYcanary…`, `ghp_ABCDEF…`, `glpat-ABCDEF…`); kein echtes Token, kein Wert in Template, Block oder erzeugten Dateien (Test `TokenReferencesCarryNoValue` grün, Wächter bewiesen durch M3/Wächter-Zeilen im Test).
- **Kommentare (§3.7/§3.9):** Suche nach „sandbox Dockerfile“ im Baum ohne Treffer außer Zitat im Review-Report und Plan; die berichtigten Kommentare (`dockerfileBlock`, `dockerfileBlockFeatureActive`, `agentCredentialEnv`, `checkSandboxCredentials`, `AgentEnv`, `export_test.go`) sagen den Ist-Zustand; Block-Kommentare im Template nennen nur Gelesenes („nach Suchmuster, nicht abschließend“). Kein „früher/former/no longer“-Muster im Diff gefunden. Handbuchzeilen tragen den Zusatz „nach Suchmuster, nicht abschließend“ (F-2 des Reviews behoben).
- **Plan-Tabelle gegen `git diff --stat`:** alle 25 Dateien des Diffs haben einen Plan-Eintrag (`CHANGELOG.md`, beide READMEs, `main_test.go`, Roadmap, Slice selbst, `open/README.md`, Welle-Datei, Review, Handbuch, `replay.md`, `forge_cli_docker_test.go`, sechs Quelltext-Dateien samt Template, fünf Testdateien); jede „update/neu“-Zeile hat einen Diff. Die Zeile `export_test.go` trägt nach der Fix-Runde „geändert“ mit dem Zusatz; `generate_ollama_test.go` ist als „update (bedingt)“ geführt und wurde geändert (Helfer verallgemeinert). Kein Eintrag ohne Diff, keine Diff-Zeile ohne Eintrag.

## DoD-Abgleich

| DoD-Punkt | Haken | Urteil |
|---|---|---|
| 1 Katalog/Generator mit B1–B11, X1–X6, je Mutation | [x] | **getragen.** Eigene Mutationen M1–M7 rot; Mutationstabelle in §7 des Slice vorhanden; B4/B8/B10 begründet ohne Verhaltenstest. |
| 2 Lauf belegt (`docker`): Bau beide Profile, Negativprobe, Egress, Kanarienwerte | [x] | **getragen**, von mir selbst gefahren (vier Tests PASS). Standardprofil-Bau, im Plan als ungemessen geführt, gelingt. |
| 3 Kleine Doku (kennungsfrei in `docs/user`/READMEs; CHANGELOG; replay) | [x] | **getragen.** READMEs 16 Features = 16 Einträge in `featureCatalogue()`; Aussagen im Handbuch haben Tests (Hosts, Profile, Pin, Ausweisung). |
| `make gates` grün, `make test-docker` gezielt | [x] | **getragen**, selbst gefahren (Exit 0; Docker gezielt, Grund genannt). |
| Review durchgeführt, Report liegt vor | [ ] | Report liegt vor, Fix-Runde (F-1, F-2) nachgelesen und im Code belegt; der Haken fehlt noch — vom Planner zu setzen. |
| Doku-Update öffentlicher Vertrag | [x] | getragen (Umfang Punkt 3). |
| Closure-Notiz mit Steering-Loop-Eintrag; Register fortgeschrieben; Risiken §6 mit Ausgang; drei Paarungen | [ ] | **offen — Planner-Pflicht bei der Closure** (§7 enthält Umsetzung, aber keine Closure-Notiz; §6 trägt noch „Ausgang: bei der Closure“). Keine Verifier-Beanstandung der Haken; erwartete Register-Einträge siehe Slice §7. |

Die gesetzten Haken stimmen mit den tatsächlichen Läufen überein; kein Sensor wurde als gelaufen ausgewiesen, der nicht lief. Nicht gelaufene Sensoren (`make verify-depguard`, vollständiges `make test-docker`) sind im Slice nicht als gelaufen behauptet und hier begründet.

## Ungemessen — und ob es den Haken trägt

- **`arm64`:** Der Block trägt die `arm64`-Prüfsummen; sie sind gegen die veröffentlichten Prüfsummen-Dateien nachgelesen, aber **kein `arm64`-Bau ist gelaufen** (Host `amd64`). Der Slice nennt das (Plan §1, ADR-Konsequenzen, Zusage-Menge Linux `amd64`). Trägt den Haken, weil `arm64` nicht zugesagt ist; die Negativprobe überschreibt beide Architekturwerte, gilt also auf beiden Hosts. Eine falsche `arm64`-Prüfsumme würde dort mit Bau-Abbruch auffallen, nicht still.
- **Gültiges Token:** Kein Aufruf mit gültigem Token, keine Hosts nach der Anmeldung, keine Telemetrie-/Update-Hosts, kein `glab duo`/`glab orbit`, keine selbst gehostete Instanz. Vom Projektinhaber aus der Zusage genommen (B8, [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) Punkte f/h). Der Egress-Test zeigt `glab` → `gitlab.com` im Sinn der Messung (DNS-Auflösung, kein Erfolgspfad mit Anmeldung); `gh` hat im Test nur `gh --version` und die Prüfung, dass die Hosts in der Liste stehen. Trägt den Haken, weil es nicht zugesagt ist.
- **Kombination mit Ollama (Lauf):** Nur das Golden X4 (Reihenfolge, Hosts, `mounts`, NoOp); die Kombination `gh` + `glab` + `ollama` + `nestedRuntime` wurde nicht gebaut oder gestartet. Der Trigger der Beobachtung [`BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf`](../plan/planning/observations/BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf/observation.md) ist eingetreten; der Slice führt es als Risiko (§6) und nicht als Zusage. Trägt den Haken, ist aber bei der Closure als „weiter offen“ im Register zu führen.
- **Standardprofil ohne Token-Referenz:** Nutzerpfad nur als Abwesenheit gepinnt (X1); Prüfgegenstand des Validators, Handbuch im Spec-Slice.
- **Netz:** `EgressHosts` und der Bau rufen `github.com`/`gitlab.com` live an; ein Netzfehler färbt rot, ohne dass sich der Code ändert (benannt im Plan §3 Ansatz und im Test-Kopf). Meine Läufe waren grün, die Pins stimmen zum Prüfzeitpunkt (Austausch eines Release-Assets ist nicht ausgeschlossen, [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) Entscheidung 8).

## Carveouts

Neu: keine. Gelöst: keine. Unverändert: keine berührt (`docs/plan/planning/in-progress/carveouts.md` nicht im Diff).
