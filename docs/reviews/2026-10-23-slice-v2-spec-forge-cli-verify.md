# Verification-Report: slice-v2-spec-forge-cli — 2026-10-23

**Rolle:** Verifier (Modul 11, [`harness/verification.md`](../../harness/verification.md)), frischer Kontext, kein Self-Review, kein Reparieren. Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD, [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md), Lastenheft `0.22.0` und Code?

**Gegenstand:** [`slice-v2-spec-forge-cli`](../plan/planning/done/slice-v2-spec-forge-cli.md) (Welle [`welle-forge-cli-gh-glab`](../plan/planning/done/welle-forge-cli-gh-glab.md)), Diff `d22b7d61..b52a531d`, Vertrag [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte). Review: [`2026-10-23-slice-v2-spec-forge-cli.md`](2026-10-23-slice-v2-spec-forge-cli.md) (kein HIGH/MEDIUM, drei LOW), Fix-Runde `b52a531d`.

## Ergebnis

**Schließbar** (nach den Closure-Pflichten des Planners: Haken, Lerneintrag, Risiko-Ausgänge, Register, Moves). Kein Blocker, kein MEDIUM, kein neues LOW. Zwei INFO. Alle Sensoren, deren Ausgabe ich nicht aus dem Bericht des Implementers übernehmen wollte, habe ich selbst gefahren (Host `linux/amd64`, Docker-Daemon lokal, Docker 29.9.2).

## Sensoren — selbst gefahren

| Sensor | Ergebnis | Beleg |
| --- | --- | --- |
| `make gates` (auf `b52a531d`) | **Exit 0** (sichtbar als `GATES_EXIT=0`) | Coverage 92,00 % gegen 90 %; `ablage-pruefung` 26 Fälle, 0 Fehlschläge; `docs-check` und `baseline-verify` grün |
| `make doc-trace` | Exit 0 | 84 Anforderungen, 0 Waisen |
| `make doc-planning` | Exit 0 | 766 Dateien, 0 Befunde |
| `make doc-immutable RANGE="4655bcff..HEAD"` | Exit 0 | 46 Commits, 0 Befunde; `git diff 4655bcff..HEAD --stat -- docs/plan/adr` leer |
| Image `u-boot:test-docker-tools` neu gebaut, `go test -tags docker -run 'TestE2E_ForgeCLI_(BuildsInBothProfiles\|EgressHosts\|PinMismatchBreaksBuild\|TokenReferencesCarryNoValue)' ./internal/e2e/` | **Exit 0**, `ok … 145.280s` | `BuildsInBothProfiles` PASS (default_profile, sandbox_profile), `PinMismatchBreaksBuild` PASS (GH, GLAB), `EgressHosts` PASS, `TokenReferencesCarryNoValue` PASS |
| Container (`docker ps -a -q`) vorher / nachher | 4 / 4 | kein Container bleibt zurück |
| `make verify-depguard` | nicht gelaufen | Grund: keine Importregel und kein Go-Import berührt; der Diff enthält keinen Code |
| `make test-docker` (voll) | nicht gelaufen | Grund: der Slice ändert weder Code noch Test; die vier Tests, an denen Spezifikationssätze allein hängen, sind gezielt gelaufen |

Arbeitsbaum vor und nach allen Läufen und Mutationen: `git status --short` leer.

## Nachlauf der Fix-Runde (L-1 bis L-3)

- **L-1 Dauer.** Nachgerechnet an [`rohdaten-07-glabtgz-restricted.txt`](../messungen/slice-v1-forge-cli-quellen-messung/rohdaten-07-glabtgz-restricted.txt): Kopf 10:38:29,440, letzter Eintrag 10:40:36,964, `gesamtdauer=127.8s`; `…-restricted-allow.txt`: Kopf 10:40:40,296, `gesamtdauer=20.0s`. Das Handbuch (Abschnitt *Forge-CLIs (optional)*, Punkt *Egress*) sagt „Der Containerlauf der Messung dauerte 127,8 statt 20,0 Sekunden … ganzer Lauf samt Start, ein Host, ein Lauf, keine Dauer eines einzelnen Aufrufs“. Stimmt mit den Rohdaten, die Auslieferungs-Variante (`glabtgz`) ist die Quelle. Die rund 124 s gegen rund 17 s der Plan-Zeile H14 stimmen mit den Zeitstempeln der Aufrufe (10:38:32 bis 10:40:36 gegen 10:40:43 bis 10:40:59).
- **L-2 Platzhalter-Satz.** `rohdaten-07-glabtgz-discover.txt`: ohne Variable `api`, `repo-view`, `release-list`, `release-download` Exit 0; mit `GITLAB_TOKEN`-Platzhalter `api` und `repo-view` Exit 1, `auth-status` Exit 1 (Token gefunden). `release list` und `release download` mit Variable: keine Zeile `AUSSAGE` — das Handbuch nennt sie „nicht gemessen“. Stimmt.
- **L-3 Allaussagen.** `Wirkung in beiden Profilen`, `Installation` in der neuen Stelle der Spezifikation und die Katalogstelle in der Spezifikation sind auf den gepinnten Umfang eingeschränkt (`USER root` genau einmal nur ohne `ollama` und Laufzeit-Feature; „Block steht vor dem letzten `USER vscode`“ und „Image läuft als `vscode`“; „kein Schlüssel unter `features`“ nur im Standardprofil). Jede Einschränkung gegen den Test gelesen: `assertForgeBlock` (`generate_forge_cli_test.go` Zeilen 98 bis 118) prüft `LastIndex(df, "USER vscode")`, E2E-Skript `USER:vscode`, `TestGenerateDevcontainer_ForgeCLI_DefaultProfileJSONUnchanged` das Standardprofil. Deckt.

## Eigene Mutationen (Rücknahme per `git apply -R`, Baum danach sauber)

| Nr. | Mutation | Erwartung der Beleg-Tabelle | Ergebnis (`make test`, Exit 2) |
| --- | --- | --- | --- |
| V1 | `agentCredentialEnv`: `GITLAB_TOKEN` vor `GH_TOKEN` (Reihenfolge der Credential-Variablen) | B13 | rot: `TestGenerateDevcontainer_ForgeCLI_RemoteEnvReferences` (both, with_claude); zusätzlich `TestDoctor_SandboxCredentials_ForgeCLI` (both) und `…RemedyRestoresReference` |
| V2 | `doctor_sandbox.go`: Text „has no reference for: %s“ → „lacks a reference: %s“ | B15 (Meldung `no reference for: GH_TOKEN, GITLAB_TOKEN`) | rot: `TestDoctor_SandboxCredentials_ForgeCLI` (6 Fälle), `…_Agent` (3 Fälle), `…RemedyRestoresReference` |
| V3 | `egressFeatureHosts`: `gh` ohne `api.github.com` | B11 (Egress-Zeilen) | rot: `TestRun_EgressList_ForgeCLIHostOrigins` (Herkunft `base,feature:gh` für `api.github.com`). `TestGenerateDevcontainer_ForgeCLI_EgressHosts` bleibt grün: der Host steht auch in der Basisliste, die sortierte Vereinigung ändert sich nicht; der Satz „führt für `gh` die Hosts `api.github.com` und `github.com`“ ist an der Anzeige gepinnt, nicht an der Liste. Die Spezifikation sagt das so (Herkunft `base` und `feature:gh`). |
| V4 | Template: `--no-same-owner` im `glab`-Block entfernt | B4 | rot: `TestGenerateDevcontainer_ForgeCLI_BlocksInBothProfiles` (glab und both in allen drei Profilfällen) |
| V5 | Template: `glab`-Block vor den `gh`-Block verschoben (21 Zeilen) | B6 (Reihenfolge `ollama`, `gh`, `glab`) | rot: `TestGenerateDevcontainer_ForgeCLI_TogetherWithOllamaAndRuntime` (podman, docker) |

V1 bis V4 in einem Lauf, V5 allein. Jede Zusage der Stichprobe wird von dem Test getragen, den die Beleg-Tabelle nennt.

## DoD-Abgleich

| DoD-Punkt | Befund |
| --- | --- |
| Spezifikation und Sicht | erfüllt. Die neue Stelle [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) mit der geplanten Gliederung; [`SPEC-023`](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) (Sandbox-Regel, Zugangsdaten, Warnungstabelle, Prüfung Punkt 3 und Schluss), [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), [`SPEC-017`](../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry), [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) und Änderungshistorie nachgezogen; der falsche Satz „Andere Variablen reicht er nicht durch“ ist berichtigt (Reihenfolge `ANTHROPIC_API_KEY`, `GH_TOKEN`, `GITLAB_TOKEN`, belegt durch V1). `architecture.md` ohne Diff (Diff-Stat leer), Vorab-Suche mit Befund in §7. Verweise nach unten im Diff der Spezifikation: Suche nach `ADR-`, `slice-`, `welle-`, Roadmap, Carveout, `docs/plan` — 0 Treffer. |
| Handbuch | erfüllt. Abschnitt *Forge-CLIs (optional)* (Zeilen 2036 bis 2147), Links aus den vier Tabellenzeilen, Nachträge in *Zugangsdaten*, Sandbox-Anleitung, *Umgebung*, *Umgebungsvariablen*, `doctor`-Aufzählung und Prüftabelle. Suche im Diff nach `LH-`, `SPEC-`, `ADR-`, `ARC-`, `slice-`, `welle-`, `roadmap`, `spec/`, `docs/plan`, `ghp_`, `glpat`, `github_pat`: **0 Treffer** (§5 Regel 6 gehalten). Kein Tokenwert. Handbuch und Spezifikation stimmen überein: Aktivierung, `version`/`source`, Meldung je Feature und Profil, `--quiet`, Hosts, Reihenfolge der Variablen, Standardprofil ohne `remoteEnv`, Nicht-Zusagen. |
| Beleg-Tabelle | erfüllt. Alle 24 genannten Testnamen: `grep -rn "^func <Name>(" internal cmd` je genau eine Definition (selbst gefahren). Streichungen S1 bis S19 vorhanden; die geschärften Zeilen B3, B5, B6, B7, B10 stehen in Spezifikation und Handbuch in der geschärften Form. |
| Gates | erfüllt (siehe Sensoren). `doc-immutable` 0 Befunde, ADR-Diff leer. |
| Review und Fix-Runde | Report liegt vor; L-1 bis L-3 nachgelesen und gegen Rohdaten und Tests belegt. |
| Doku-Update | `roadmap.md` (Arbeitszeile), `open/README.md`, Welle-Datei durch den Move; ADR-Index, READMEs, `CHANGELOG.md` unverändert (Plan: kein Verhalten). |
| Closure-Notiz, Lerneintrag, Register, Risiko-Ausgänge, drei Paarungen, Haken in §2 | **offen**: §7 trägt „(Wird bei der Closure gefüllt.)“, alle Kästchen in §2 sind leer, §6 hat überall „(bei der Closure)“. Das ist Auftrag des Planners vor dem Move nach `done/`, keine Lücke der Umsetzung. |

## §3.10, §3.4, Zellenzahl, Plan-Tabelle

- **§3.10 Bedingung ohne Test:** Stichprobe der neuen Stelle und des Handbuchs gegen die Tests. Jede Bedingung hat einen Test oder steht in der Streichungstabelle: Block-Eigenschaften (`assertForgeBlock`, `assertNoForgeSecrets`), Reihenfolge (V5), `version` ohne `source` Exit `10` (`…VersionOverrideWithoutSourceRejected`, MAIN-Teiltest), Meldungen (`mustNot` für Version und Attestierung bei `glab`, Zeile 266 und 267 des Tests), Token-Referenzen (V1), `doctor` (V2), Hosts (V3, E2E), Kein-Token-Speicher und `Config.Env` (E2E, gelaufen).
- **§3.4:** keine ADR-, Slice-, Wellen-, Roadmap-Kennung in `spec/`; `architecture.md` unberührt.
- **Zellenzahl:** alle neuen und berührten Tabellenzeilen gezählt: Katalog 2, Egress-Tabelle des Handbuchs 2, Umgebungsvariablen 2, [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) 3, Warnungstabelle von [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) 4, Prüfliste von [`SPEC-017`](../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) 2. Keine Zelle zu viel oder zu wenig.
- **Plan-Tabelle gegen `git diff --stat d22b7d61..HEAD`:** sieben Dateien — `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`, `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/open/README.md`, `docs/plan/planning/welle-forge-cli-gh-glab.md`, diese Slice-Datei, der Review-Report. Jede hat eine Zeile der Plan-Tabelle (§3); jeder Eintrag „unverändert“ (`architecture.md`, `benutzerhandbuch-standard.md`, `cli-json-contract.md`, READMEs, `CHANGELOG.md`, `internal/`, `cmd/`, Lastenheft, ADRs) ist im Diff tatsächlich leer (selbst geprüft für `architecture.md`, ADR-Verzeichnis, Code). Mit diesem Report kommt `docs/reviews/…-verify.md` unter die Plan-Zeile `docs/reviews/`.

## Ungemessen oder nicht gepinnt — trägt es den Haken?

| Gegenstand | Stand im Text | Trägt den Haken? |
| --- | --- | --- |
| Nutzerpfad Standardprofil (Token in den Container bringen, danach `gh api`) | Handbuch: „weder zugesagt noch gemessen“; Streichung S7; Abwesenheit der Referenz gepinnt (B19) | ja — als Nicht-Zusage geschrieben, keine Abhilfe ohne Probe. Fahrt am Stück ist Prüfgegenstand des Validators der Welle. |
| `arm64`-Bau | Handbuch und Spezifikation unter „Nicht zugesagt“; die Meldung „has no archive for <arch>“ ist nur als Text im Golden gepinnt, nicht als Lauf | ja |
| Aufruf mit gültigem Token, Telemetrie- und Update-Hosts, leere Variable im Container, `release list`/`release download` mit Variable | je als „nicht gemessen“ oder „nicht zugesagt“ benannt | ja |
| Kombination mit `ollama` und `nestedRuntime` | nur Reihenfolge, Volumes, Hosts (V5); Streichung S6 | ja |

## Befunde

- **INFO-1.** Die Aussage „eine andere Architektur als `amd64` und `arm64` bricht den Build ab“ (neue Stelle, *Installation*) ist als Text des `case`-Zweigs im Golden gepinnt, nicht als Lauf. Das ist ein Shell-Verhalten mit einer Zeile Schutz; als Zusage zulässig, weil der Text des Zweigs gepinnt ist. Kein Handlungsbedarf.
- **INFO-2.** V3 zeigt: Die Liste der Hosts von `gh` ist gegen das Entfernen von `api.github.com` nur an der Anzeige (`TestRun_EgressList_ForgeCLIHostOrigins`) abgesichert, weil der Host auch in der Basis steht. Der Spezifikationssatz benennt genau diese Anzeige. Kein Handlungsbedarf.

Keine Findings der Klassen HIGH, MEDIUM oder LOW.

## Sensoren, Carveouts, Commit

- Carveouts: keine neuen, keine gelösten, Inventar unverändert (Diff berührt `carveouts.md` nicht).
- Commit des geprüften Stands: `b52a531d`. Dieser Report ist ein weiterer Commit darauf; nichts gepusht, kein Tag.
