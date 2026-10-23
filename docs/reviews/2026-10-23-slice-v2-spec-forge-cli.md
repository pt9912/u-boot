# Review-Report: slice-v2-spec-forge-cli (Spezifikation und Handbuch zu `gh` und `glab`) — 2026-10-23

**Rolle:** Reviewer (Modul 10, [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md), [`harness/review.md`](../../harness/review.md)), frischer Kontext. Prüfung gegen Plan, Entscheidungen und Hard Rules, nicht gegen die DoD (Verifier).

**Gegenstand:** `git diff d22b7d61..HEAD` (HEAD `877d943f`, 6 Dateien): Slice [`slice-v2-spec-forge-cli`](../plan/planning/done/slice-v2-spec-forge-cli.md), Welle [`welle-forge-cli-gh-glab`](../plan/planning/welle-forge-cli-gh-glab.md), [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) (Accepted), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer).

**Eingang gelesen:** vollständiger Diff von `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md` und des Slice (§3 Beleg-, Handbuch- und Streichungs-Tabelle, §7 Umsetzung); `generate_forge_cli_test.go` vollständig; `TestDoctor_SandboxCredentials_ForgeCLI` und `…RemedyRestoresReference`, `TestDoctor_FeaturesDrift_ForgeCLI`, `TestRun_Generate_ForgeCLINotices`, `TestRun_EgressList_ForgeCLIHostOrigins`, E2E `forge_cli_docker_test.go` (Bau, Pin, Egress, Token); `Dockerfile.tmpl` (Blöcke `gh`/`glab`, `USER`-Struktur, kein `ENV`); `forgeCLIWarnings` (Meldungstext), `agentCredentialEnv`, `doctor_sandbox.go` (Meldung); Rohdaten `rohdaten-07-{gh,glab}tgz-{discover,restricted,restricted-allow}.txt` und Lastenheft 0.22.0 (Zeile zu `gh`/`glab`).

## Ergebnis

Kein HIGH, kein MEDIUM. Drei LOW (zwei Messaussagen im Handbuch einen Schritt weiter als die Rohdaten, Spezifikationssätze mit Allaussage, die ein Teilpin trägt), vier INFO. Die Teilpins B3/B5/B6/B7/B10 sind im Slice ehrlich ausgewiesen; die Sätze zu `setup-git` (H12), Telemetrie (H13) und Dauer (H14) sind bis auf die genannten Punkte durch die Nicht-gemessen-Spalte gedeckt.

## Findings

### L-1 — LOW — „Aufrufreihe dauerte 127,8 statt 20,0 Sekunden“ ist die Gesamtdauer des Containerlaufs

- `kategorie`: LOW · `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.10 (eine Aussage nennt nur, was gemessen wurde) · `klasse`: zusage-weiter-als-messung
- `pfad`: `docs/user/benutzerhandbuch.md` · „Die Aufrufreihe der Messung dauerte 127,8 statt 20,0 Sekunden“ (auch Slice §7, Zeile H14: „Aufrufreihe 127,8 s“)
- `befund`: Die Rohdateien führen `gesamtdauer=127.8s` (`rohdaten-07-glabtgz-restricted.txt`) und `gesamtdauer=20.0s` (`…-restricted-allow.txt`) als Dauer des ganzen Containerlaufs. Die Zeitstempel der Aufrufe (10:38:32 bis rund 10:40:36 gegen 10:40:43 bis rund 10:41:00) ergeben rund 124 s gegen rund 17 s; Start, Egress-Init, Schlüsselliste und DNS-Auszug liegen in der Gesamtdauer mit drin. Der Satz nennt die Zahlen als Dauer der „Aufrufreihe“.
- `risiko`: Die Größenordnung stimmt; das Handbuch schreibt die Zahl einer engeren Größe zu, als die Messung trägt (die Variante `glabtgz` ist die Auslieferung, die Zahl gehört zum ganzen Lauf).
- `verifizierbar`: Review-only (Rohdaten)

### L-2 — LOW — „bei `glab` scheitern damit auch die Aufrufe gegen ein öffentliches Projekt, die ohne Variable gelingen“

- `kategorie`: LOW · `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.10 · `klasse`: zusage-weiter-als-messung
- `pfad`: `docs/user/benutzerhandbuch.md` · „Mit Variable, ohne Anmeldung durch `u-boot`“, Satz „bei `glab` scheitern damit auch die Aufrufe gegen ein öffentliches Projekt, die ohne Variable gelingen“
- `befund`: Mit Platzhalter-Token gelaufen und mit `401` gescheitert sind in `rohdaten-07-glabtgz-discover.txt` nur `api-gitlab_token`, `api-gitlab_access_token`, `api-oauth_token` und `repo-view-gitlab_token`. `release list` und `release download` hat die Messung mit Variable nicht gefahren (kein Eintrag `release-*-gitlab_token`); der Satz steht neben der Liste „ohne Variable“ (`api`, `repo view`, `release list`, `release download`) und liest sich als Aussage über alle vier.
- `risiko`: Der Leser nimmt an, auch `release download` scheitere mit einem (ungültigen) Token; gemessen ist es für `api` und `repo view`.
- `verifizierbar`: Review-only (Rohdaten)

### L-3 — LOW — Allaussagen in der Spezifikation, die nur ein Teilpin trägt

- `kategorie`: LOW · `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.10 (Bedingung in Spec ohne pinnenden Test oder gestrichen) · `klasse`: spec-bedingung-ohne-code-entsprechung
- `pfad`: `spec/spezifikation.md` · (a) „`USER root` steht im Dockerfile höchstens einmal“ und „Das Dockerfile endet auf `USER vscode`“ ([`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten), *Wirkung* und *Installation*); (b) „erzeugen keinen Schlüssel unter `features` in `devcontainer.json`“ ([`SPEC-023`](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), neuer Spiegelstrich)
- `befund`: (a) Die Tests zählen `USER root` nur im Standardprofil (`…UserRootOnlyWithReason`) und im Sandbox-Profil ohne Laufzeit (`…BlocksInBothProfiles`); `…TogetherWithOllamaAndRuntime` prüft nur den ersten Marker, nicht die Anzahl. „Endet auf `USER vscode`“ prüft `assertForgeBlock` nur als „der Block steht vor dem letzten `USER vscode`“, das Image-`USER` der E2E-Test. (b) Für die erzeugte Sandbox-`devcontainer.json` gibt es keinen Test auf leeres `features` (im Slice als Teilpin B3 ausgewiesen). Beide Aussagen gelten nach `Dockerfile.tmpl` (ein `USER root`-Block, `USER vscode` am Ende) beziehungsweise nach der Katalog-Invariante (keine Quelle), sind aber nicht auf der Breite der Aussage gepinnt.
- `risiko`: Ein Eingriff in die Vorlage (zweiter `USER root`-Block, Zeile nach `USER vscode`) oder in die Feature-Auflösung würde die Spezifikation widerlegen, ohne dass ein Gate rot wird; der Slice weist die Lücke aus, der Spezifikationstext trägt den Vorbehalt nicht.
- `verifizierbar`: `make test` (kein Test auf der vollen Breite)

## INFO

- **I-1 (§3.10, Docker-Tests).** Mehrere Sätze der Spezifikation und des Handbuchs hängen allein an Tests mit Build-Tag `docker` (Bau in beiden Profilen, Version, Pin-Abbruch, kein Token-Speicher im Image, `lookup gitlab.com`, `Config.Env`, Pfad `/usr/local`, Eigentümer). `make gates` fährt sie nicht; `make test-docker` ist in diesem Lauf nicht gefahren (Diff ohne Code und Test; die Läufe liegen im Code-Slice und im Verify-Report dazu). „Andere Architektur als `amd64`/`arm64` bricht ab“ ist nur als Text der `case`-Zeile im Golden-Test gepinnt, nicht als Lauf. Aufgabe des Verifiers.
- **I-2 (Schwerpunkt 3, §3.4).** Kein Verweis nach unten: Der Diff von `spec/spezifikation.md` enthält 0 Treffer für ADR, Slice, Welle, Roadmap, Carveout, `BEO-`, Messprotokoll; `spec/architecture.md` und `spec/lastenheft.md` haben keinen Diff; die Sicht nennt weder `gh`, `glab`, `claude-code` noch `ollama` (eigene Suche, 0 Treffer) und bleibt korrekt unverändert.
- **I-3 (Schwerpunkt 6).** Die Handbuchorte Zeile ~639 (`doctor --quiet`: „aktives Feature (`claude-code`, `gh`, `glab`)“) und ~2811 (Tabelle der Prüfungen) sind stimmig mit der Code-Registry ([`SPEC-017`](../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry)) und der Prüfung ([`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)). Die „feste Reihenfolge `ANTHROPIC_API_KEY`, `GH_TOKEN`, `GITLAB_TOKEN`“ in dieser Stelle und in [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) deckt sich mit `agentCredentialEnv` (Reihenfolge der Anhänge `claude-code`, `gh`, `glab`) und dem Test `…RemoteEnvReferences` (Fall „with claude“). Die Doctor-Meldung „no reference for: GH_TOKEN, GITLAB_TOKEN“ stützt den Satz „nennt alle fehlenden Variablen“ (`…ForgeCLIRemedyRestoresReference`).
- **I-4 (Schwerpunkt 7 und 8).** Plan-Tabelle gegen `git diff --stat`: `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`, der Slice, `roadmap.md`, `open/README.md` und die Welle-Datei (die beiden letzten durch den Move des Slice) sind je durch eine Zeile gedeckt; die als „unverändert“ geführten Dateien haben keinen Diff. Commit-Nachrichten: 0 Treffer für `SPEC-`/`ARC-` über `d22b7d61..HEAD`; jede nennt `LH-`/`ADR-`/Slice-Kennung.

## Prüfung nach Schwerpunkten

1. **§3.10, Testnamen einzeln gelesen.** `TestGenerateDevcontainer_ForgeCLI_BlocksInBothProfiles` (B4/B5/B19), `…UserRootOnlyWithReason`, `…Notices` (Wortteile `cli/cli`, `GH_VERSION`, `SHA-256`, „no signature was found“, „not exhaustive“, `attestation`; `mustNot` Version), `…EgressHosts`, `…RemoteEnvReferences`, `…DefaultProfileJSONUnchanged`, `…VersionOverrideWithoutSourceRejected`, `…SourceOverrideIsExternalFeature`, `…TogetherWithOllamaAndRuntime`, `…ActivateAndDeactivateLater`, `TestDoctor_SandboxCredentials_ForgeCLI` (16 Fälle, Platzhalter-Werte erfunden), `TestRun_Generate_ForgeCLINotices` (einmal je Feature, `--quiet`, Exit 10), `TestRun_EgressList_ForgeCLIHostOrigins` und die E2E-Tests tragen die Sätze, die die Spezifikation darauf stützt. Teilpins B3/B5/B6/B7/B10: im Slice §7 so ausgewiesen und im Text der Spezifikation auf den gepinnten Umfang eingeschränkt (kein „allein“ zu SHA-256, „es entsteht kein Dockerfile“, Reihenfolge nur der fünf Blöcke), Ausnahme L-3.
2. **Sätze zu `setup-git` (H12), Telemetrie (H13), Dauer (H14).** H12: `~/.gitconfig` 193 Byte mit je einem leeren und einem `gh auth git-credential`-Eintrag für `github.com` und `gist.github.com` in `rohdaten-07-ghtgz-discover.txt`, Platzhalter-Treffer 0; „Wirkung auf einen Git-Aufruf nicht gemessen“ im Text; der Pfad `/usr/bin/gh` der Feature-Variante steht nicht im Satz. H13: `GH_TELEMETRY`/`DO_NOT_TRACK` aus `gh help environment`, `telemetry=true`/`check_update=true` aus `glab config get`; „kein `ENV`“ stimmt (`Dockerfile.tmpl` hat keine `ENV`-Zeile), Wirkung der Schalter und Hosts stehen unter „Nicht gemessen“. H14: siehe L-1.
3. **Messvariante gegen Auslieferung.** Die Aussagen H7/H8/H11 sind an `ghtgz`/`glabtgz` belegt (dasselbe Archiv wie der erzeugte Block): `gh api` ohne Variable Exit 4 mit Hinweis auf `GH_TOKEN`, `glab`-Aufrufe ohne Variable Exit 0, `repo clone` Exit 1 mit `401`, `git_protocol` `ssh`, `config set … https` Exit 0, `gh config list` `git_protocol=https`. „Der Clone-Weg nach der Umstellung ist nicht gemessen“ ist vorsichtiger als die Rohdaten (`repo-clone-https` lief, Exit 1); kein Befund. Abweichungen: L-1, L-2.
4. **§5 Regel 6.** `docs/user/benutzerhandbuch.md`, README und `docs/user/` tragen im Diff keine `LH-`, `SPEC-`, `ADR-`, `ARC-`, Slice-, Wellen- oder Roadmap-Kennung und keine Tokenform (`ghp_`, `glpat`, `github_pat`: 0 Treffer, eigene Suche).
5. **Tabellenzeilen, d-check, Anker.** Die neuen Zeilen in [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (3 Zellen), [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) (4 Zellen) und die Registerzeile in [`SPEC-017`](../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) passen zur Spaltenzahl; `make docs-check` im `make gates` ohne Befund, der Anker auf die neue Stelle löst auf.
6. **Kommentare und Zustandsfelder (§3.7, §3.9).** Kein Konjunktiv über Verworfenes, keine „früher/ohne X wäre“-Muster in den neuen Texten; die Historie-Zeile in Abschnitt 7 der Spezifikation ist Chronik am dafür vorgesehenen Ort. Die Arbeitszeile der Roadmap nennt Zustand und Slice, keine Chronik.

## Nicht geprüft

- `make test-docker`: nicht gelaufen (Docker-Tests, reine Doku-Änderung; die Läufe liegen im Code-Slice). Siehe I-1.
- `make verify-depguard`: nicht gelaufen, keine Importregel berührt.
- `docs/user/benutzerhandbuch-standard.md`, `docs/maintainer/cli-json-contract.md`: unverändert, nur auf Diff-Freiheit geprüft.

## Sensors

`make gates`: Exit-Code `0` (selbst gefahren; lint, Tests, Coverage, `docs-check`, `baseline-verify`, `ablage-pruefung` mit „26 Faelle, 0 Fehlschlaege“, Ausgabe zuletzt `EXIT=0`).
