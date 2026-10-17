# Review-Report: slice-v1-devcontainer-podman-build — 2026-10-17

**Review-Art:** Code — gegen Plan, Spec, ADR und Hard Rules (Maintainability; keine DoD-Prüfung, die ist Verifier-Aufgabe).

**Gegenstand:** Commit-Range `5f372bc..HEAD` (13 Commits, Spitze `e017907`), Slice `slice-v1-devcontainer-podman-build` <!-- d-check:ignore (Slice-Kennung, Slice liegt in in-progress/) --> (Welle `welle-devcontainer-sshd-claude-code`, siebter Slice).

**Skill:** `.harness/skills/reviewer.md` @ `e017907` · <!-- d-check:ignore (Skill-Pfad des Repos) -->
**Modell:** Sonnet 5.5 (`claude-sonnet-5-5`) · **Datum:** 2026-10-17

**Eingangs-Kontext:**

- Slice-Plan `slice-v1-devcontainer-podman-build` (§1 Entscheidungsmatrix, §2 DoD, §3 Plan-Tabelle) <!-- d-check:ignore (Slice-Kennung, liegt in in-progress/ und wandert) -->
- [ADR-0021](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md), [ADR-0014](../plan/adr/0014-nested-podman-sandbox-devcontainer.md), [ADR-0012](../plan/adr/0012-devcontainer-egress-firewall.md)
- [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-NFA-SEC-*`](../../spec/lastenheft.md) (Sicherheitsfamilie, gelesen)
- [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [`LH-FA-DEV-007.b`](../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman)
- `AGENTS.md` §3.2 bis §3.10, `harness/review.md`, `harness/conventions.md` (`MR-008`)

Eigene Sensor-Läufe: `make doc-immutable RANGE="5f372bc..HEAD"` ([ADR-0014](../plan/adr/0014-nested-podman-sandbox-devcontainer.md) unverändert, 0 Befunde), `make docs-check` (0 Befunde), `make lint` (grün), `make test` (grün auf dem Stand `e017907`) und drei Mutationen (unten), jede per `git checkout -- <Datei>` zurückgesetzt; Arbeitsbaum danach sauber.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die sechste Lockerung macht `/dev/net/tun` zur Host-Vorbedingung: fehlt das Device auf dem Container-Host, startet der Container nicht (Docker-Fehler, keine u-boot-Meldung, keine Degradationsstufe wie bei `/dev/fuse`). Das Lastenheft nennt diese Vorbedingung nicht: Die Host-Zusage („startet unter Docker auf Hosts, die verschachtelte User-Namespaces erlauben“) und der Degradations-Absatz erwähnen nur User-Namespaces und `/dev/fuse`; auch die Spezifikation sagt nur „ohne das Device scheitert jedes `RUN`“. Handbuch und [ADR-0021](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) nennen die Vorbedingung, `u-boot doctor` prüft `/dev/net/tun` nicht (der Slice schließt das ausdrücklich aus). Failure-Szenario: Ein Nutzer mit bestehender `nestedRuntime: podman`-Konfiguration auf einem Linux-Host ohne geladenes `tun` aktualisiert u-boot, `generate devcontainer` schreibt das neue `runArgs`, der Devcontainer startet nicht mehr; vorher lief dort `podman run --network=host`. | [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Host-Zusage, Degradation „nie stillschweigend“), [ADR-0021](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) §Konsequenzen | `spec/lastenheft.md` · „Bei `podman` startet das Ergebnis unter Docker auf Hosts, die verschachtelte User-Namespaces erlauben“; `internal/hexagon/application/doctor_sandbox.go` (prüft nur `/dev/fuse`) | nein — Review-only (`make test-docker` auf einem Host ohne `tun` würde es zeigen; Colima hat das Device) | Neue Host-Vorbedingung ohne Vertrags- und doctor-Abdeckung |
| F-2 | LOW | Das Lastenheft sagt `podman build` „auch hinter der Egress-Restriktion mit ihrer Default-Allowlist“ zu; gemessen sind `docker.io/library/alpine:3` und `debian:stable-slim`, und der Pull hängt am Blob-Host `production.cloudfront.docker.com`, den Docker Hub ohne Änderung an u-boot wechseln kann. Es gibt keinen laufenden Sensor (siehe F-4); der Re-Evaluierungs-Trigger der ADR nennt den Fall, fängt ihn aber nicht. | [ADR-0021](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) §Re-Evaluierungs-Trigger | `spec/lastenheft.md` · „auch hinter der Egress-Restriktion ([`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) mit ihrer Default-Allowlist“ | nein — Review-only | Zusage hängt an externer Host-Topologie ohne Dauer-Sensor |
| F-3 | LOW | Der Doc-Kommentar von `podmanRelaxations()` verweist mit „measured, see the nested-Podman ADR“ auf eine ADR; die sechste Option ist in [ADR-0021](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) gemessen, nicht in [ADR-0014](../plan/adr/0014-nested-podman-sandbox-devcontainer.md). Der Diff ändert die Funktion, nicht diesen Zeiger. | `AGENTS.md` §3.9 (b), §3.7 Rang-Zeiger | `internal/hexagon/application/devcontainer_sandbox.go` · „(measured, see the\n// nested-Podman ADR)“ | nein — Review-only | Kommentar-Zeiger nennt bei erweitertem Gegenstand nur eine von zwei Quellen |
| F-4 | INFO | `.github/workflows/integration.yml` läuft auf `ubuntu-latest`; dort steht `kernel.apparmor_restrict_unprivileged_userns` nach Plan und ADR vermutlich auf 1, womit beide neuen Tests (und der Nachbartest) nach der Degradationsprüfung überspringen. Die Zusage hat dann nur manuelle Messevidenz; die ADR benennt das als Folgepflicht, der Slice schließt eine Workflow-Änderung als Gate-Änderung aus. Der Skip ist ehrlich (`Errorf` vor `Skip` markiert Fehlschlag der Degradationsprüfung, Skip erscheint als SKIP, nicht als PASS). Zuständigkeit: Verifier/Projektinhaber. | [ADR-0021](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) §Konsequenzen (Folgepflicht) | `internal/e2e/sandbox_devcontainer_docker_test.go` · `skipIfNestedUsernsRestricted` | nein — `make test-docker` auf CI-Runner | — |
| F-5 | INFO | In `nestedPodmanEgressBuildScript` stehen die Marker `ALLOWED_RESOLVES` und `UNLISTED_BLOCKED` im `RUN`-Text der Dockerfile-Zeile, die `podman build` als `STEP …` ausgibt; sie sind als Ausgabe-Marker allein nicht beweiskräftig. Den Beleg trägt der Exit-Code der `RUN`-Zeilen (`getent … &&`, `! getent … &&`) und der Marker `EGRESS_BUILD_OK`; der Test ist nicht tautologisch. | Maintainability | `internal/e2e/sandbox_devcontainer_docker_test.go` · „RUN getent hosts deb.debian.org > /dev/null && echo ALLOWED_RESOLVES“ | nein — Review-only | Marker im Echo des eigenen Befehls |
| F-6 | INFO | `production.cloudfront.docker.com` steht bei `nestedRuntime: docker` mit in der Allowlist (Code, [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), Handbuch, Tests, CHANGELOG, Entscheidung 3 der ADR benennen das konsistent); gemessen ist der Pull nur im Podman-Pfad. Für den Docker-Daemon ist es dieselbe Registry, aber ungemessen. | [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [ADR-0021](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) Entscheidung 3 | `CHANGELOG.md` · „without it a pull behind the restriction failed with a DNS error“ | nein — Review-only | Messaussage auf zweiten Pfad übertragen |
| F-7 | INFO | Die Entscheidungsmatrix des Slice sieht bei „fehlende Lockerung“ die Rückführung nach `next/` vor (Sicherheitsentscheidung); der Slice verfährt stattdessen nach einer Vorab-Freigabe des Projektinhabers. Die Freigabe steht nur als Prosa im Slice, in [ADR-0021](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) (`Accepted`) und in der Lastenheft-Historie („Vereinbarung mit dem Projektinhaber“); ich kann sie nicht prüfen. Die Messung belegt die Notwendigkeit der Option (`pasta … Failed to open() /dev/net/tun`). Zuständigkeit: Verifier. | `AGENTS.md` §3.6 | Slice-Plan (Datei unter `docs/plan/planning/in-progress/`) · „Nach Vorab-Entscheidung des Projektinhabers“ | nein — Review-only | Matrix-Abweichung belegt nur durch Prosa |
| F-8 | INFO | Spezifikation und Handbuch sagen „Kurznamen aus der Alias-Liste der Distribution lösen auf, andere nicht“; gemessen ist ein Kurzname (`alpine`) und ein Nicht-Kurzname (`mycompany/img`). Die Aussage ist verallgemeinert; die ADR nennt „Kurznamen außer `alpine`“ ausdrücklich als nicht gemessen. | [ADR-0021](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) §Kontext (Nicht gemessen) | `docs/user/benutzerhandbuch.md` · „löst auf, `mycompany/img` nicht“ | nein — Review-only | Verallgemeinerung aus einer Einzelmessung |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Sechste Lockerung in `podmanRelaxations()` | geprüft, ohne Befund: ein Eintrag je Option, gleiche Warnungsform wie die fünf anderen (`relax(...)` mit `SecurityRelaxation`, jede Option einzeln); `nestedRuntime: docker` nutzt `dockerRelaxation()` und bleibt unverändert; `--privileged`/`NET_ADMIN` nicht ergänzt |
| Zählung fünf → sechs im Baum | geprüft, ohne Befund: Suche nach „five/fünf/fuenf“ im Zusammenhang mit Podman in Code, Tests, Spec, Lastenheft, Handbuch, `harness/replay.md`, CHANGELOG ohne Treffer außerhalb Slice-Plan (Lesestand-Texte) und [ADR-0021](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) (Messtabelle, dort Bezeichnung des Vor-Zustands); Tests zählen 6 bzw. 7 mit Egress |
| `runArgs` im generierten `devcontainer.json` | geprüft, ohne Befund: Test pinnt sechs Optionen in fester Reihenfolge, `NET_ADMIN` als siebte mit Egress |
| Mutation 1: `--device=/dev/net/tun` aus `podmanRelaxations()` entfernt | `make test` rot (`TestSandboxRelaxationWarnings_FlagIsOnTheEntryNotTheCode`: flagged=5, want 6), zurückgesetzt |
| Mutation 2: `production.cloudfront.docker.com` aus `egressRegistryHosts()` entfernt | `make test` rot (`TestGenerateDevcontainer_SandboxEgressWithPodman`, `…WithDocker`), zurückgesetzt |
| Mutation 3: Lockerungswarnung für Podman auf nicht gekennzeichnet gestellt | `make test` rot (`…SandboxPodman`, `TestSandboxRelaxationWarnings_*`), zurückgesetzt |
| Egress-Konsistenz Code / [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) / Handbuch / Tests / `harness/replay.md` | geprüft, ohne Befund: Host überall genannt; Wirkung auf `nestedRuntime: docker` korrekt benannt (gleiche Zeile) |
| Lastenheft 0.12.0: Version in Kopf, §16-Historie, CHANGELOG | geprüft, ohne Befund: Kopf 0.12.0, Historienzeile mit Messstand und Quelle; keine weitere Datei im Baum führt einen veralteten Lastenheft-Stand für diesen Gegenstand (README, Handbuch ohne Versionsangabe; Spec-Historie nennt 0.11.0 für den vorigen Eintrag korrekt) |
| Lastenheft-Wortlaut gegen Messung | geprüft, ohne Befund außer F-2: Messstand (Docker 29.5.2, Colima 0.10.3, arm64, Schalter 0) steht im Text; amd64, andere Engines, `buildx`/`compose`/`# syntax=` ausdrücklich ausgenommen; Podman als Host-Engine für `podman` nicht mehr zugesagt (Rücknahme der früheren „engine-neutral“-Zusage steht in der Historie); Korrektur der Colima-Aussage (Standardwert 1 → Exit-Code 11) stimmt mit dem unveränderten Degradations-Absatz überein |
| Lastenheft ↔ Spezifikation ↔ Handbuch ↔ Code | geprüft, ohne Befund außer F-1/F-8: [`LH-FA-DEV-007.b`](../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman) sagt dasselbe wie Code und Dockerfile-Pin (sechs `runArgs`, Pakete, Symlink, kein `buildah`) |
| [ADR-0021](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md): Form, Schärft-Feld, Index | geprüft, ohne Befund: MADR-Form nach `MR-008`, `Schärft:` nach oben auf [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), ADR-Index ergänzt; Verweise nur aufwärts oder auf ADRs; `Accepted` ist durch das Messprotokoll in der ADR gedeckt (Freigabe selbst: F-7) |
| [ADR-0014](../plan/adr/0014-nested-podman-sandbox-devcontainer.md) immutable | geprüft, ohne Befund: `make doc-immutable RANGE="5f372bc..HEAD"` meldet 0 Befunde; `git diff --stat` nennt die Datei nicht |
| Integrationstests `…_NestedPodman_Build`, `…_BuildEgress`, Helfer | geprüft, ohne Befund außer F-4/F-5: Build-Tag `docker`, Zeitlimit 25 min, `--rm` und `rmi` im Cleanup, Marker plus Exit-Code, Inhalt des gebauten Images wird geprüft (`COPY`/`RUN`), Cache- und `--no-cache`-Zweig unterscheiden sich über `Using cache` |
| Plan-Tabelle gegen `git diff --stat 5f372bc..HEAD` (`AGENTS.md` §3.10) | geprüft, ohne Befund: alle 14 berührten Dateien stehen in der Tabelle; die als „nicht berührt“ geführten Dateien (`Dockerfile.tmpl`, `sandbox-init.sh.tmpl`, `devcontainer.json.tmpl`, READMEs, Welle-Datei, Workflow, `doctor_sandbox.go`) tauchen im Diff nicht auf; Roadmap-Ruhe-Marker „Nichts in Arbeit.“ entfernt, solange der Slice in `in-progress/` liegt |
| Sprachvertrag ([`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)), Architekturgrenzen, Suppression | geprüft, ohne Befund: neue Warnungen und Strings Englisch, keine Importänderung, kein `//nolint` |
| Kommentare im Diff (`AGENTS.md` §3.7, §3.9) | geprüft, ohne Befund außer F-3 |
| `make lint`, `make docs-check`, `make test` | grün, 0 Befunde |
| `make test-docker`-Mutation (Tun-Device aus den `runArgs` entfernt, nur `…NestedPodman_Build$` gefahren) | nicht ausgewertet: in dieser Sandbox ist das Netz gesperrt (`apt-get` im erzeugten Dockerfile erreicht `deb.debian.org` nicht), der Lauf scheiterte vor der Messung am Image-Build; kein Beleg, dass der Test die Mutation fängt; zurückgesetzt |
| `make verify-depguard` | nicht gelaufen: keine Importregel und kein `.golangci.yml` berührt |
| amd64, Docker Desktop, Podman-Host, große Images | nicht geprüft: nicht gemessen und im Lastenheft nicht zugesagt |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Neue Host-Vorbedingung ohne Vertrags- und doctor-Abdeckung · Zusage hängt an externer Host-Topologie ohne Dauer-Sensor · Kommentar-Zeiger nennt bei erweitertem Gegenstand nur eine von zwei Quellen · Marker im Echo des eigenen Befehls · Messaussage auf zweiten Pfad übertragen · Matrix-Abweichung belegt nur durch Prosa · Verallgemeinerung aus einer Einzelmessung

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (MEDIUM, Failure-Szenario benannt). Kein HIGH: Die Lockerung ist einzeln ausgewiesen, im Baum konsistent gezählt und gemessen begründet; der Container-Startfehler bei fehlendem Device ist laut, nicht still.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation (DoD-Abgleich, F-4 und F-7 sind Verifier-Fragen).
