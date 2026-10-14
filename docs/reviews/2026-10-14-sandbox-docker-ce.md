# Review-Report: slice-v1-sandbox-docker-ce — 2026-10-14

**Review-Art:** Code — Diff gegen Slice-Plan, ADR und Spezifikation.

**Gegenstand:** Commit `eea88cc` (Docker-Daemon aus dem Docker-APT-Repo, Fingerprint-Pin, `download.docker.com` in der Egress-Allowlist).

**Skill:** `.harness/skills/reviewer.md` @ `eea88cc`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-14

**Eingangs-Kontext:**

- [`slice-v1-sandbox-docker-ce`](../plan/planning/in-progress/slice-v1-sandbox-docker-ce.md) (DoD, Abgrenzung, Risiken)
- [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md), [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md), [ADR-0012](../plan/adr/0012-devcontainer-egress-firewall.md)
- [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)
- Vorgänger-Reports `2026-10-14-sandbox-docker-runtime.md`, `2026-10-14-sandbox-docker-runtime-r2.md`
- `AGENTS.md` (Hard Rules), `harness/review.md`, `harness/replay.md`

Nicht ausgeführt: Docker-Builds und Integrationstests (Reviewer-Rolle, keine privilegierten Container). Der RUN-Block wurde per `sh -n` auf Syntax geprüft.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Golden Case prüft den Fingerprint nur als Substring. Er prüft weder, dass die Fingerprint-Prüfung als `test … =` in der `&&`-Kette vor dem Eintrag von `docker.list` und vor `apt-get install` steht, noch `signed-by=`. Eine Änderung, die den Fingerprint nur noch im Kommentar oder in einem `echo` stehen lässt, bleibt grün. | [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) (Fitness Function), `harness/replay.md` §Replay Rules | internal/hexagon/application/generate_sandbox_test.go:667 | ja — `make test` mit einer mutierten Vorlage | Sicherheitsrelevante Eigenschaft nur als Substring gepinnt |
| F-2 | MEDIUM | Weder Golden Case noch E2E belegen die Version oder Herkunft des Daemons. Die DoD verlangt „belegen die Version“. Der E2E akzeptiert `overlay2` oder `overlayfs` und prüft nicht, dass `docker-ce` läuft (zum Beispiel `docker version` oder `docker info` mit Hauptversion ≥ 29). Der Test bliebe mit dem Debian-Daemon grün, wenn nur das Dockerfile zurückgesetzt würde. | Slice §2 DoD (Generator), [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) §Fitness Function | internal/e2e/sandbox_devcontainer_docker_test.go:437 | ja — `make test-docker` | Integrationstest pinnt die neue Eigenschaft (Herkunft/Version) nicht |
| F-3 | MEDIUM | Der Egress-Grund der Aufnahme (`apt-get update` im eingeschränkten Container erreicht `download.docker.com`, ein CloudFront-Host) ist nirgends getestet. Der Egress-E2E läuft ohne `nestedRuntime: docker`. Der Golden Case prüft nur, dass der Hostname im Skript steht. | [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) Entscheidung 3 | internal/e2e/sandbox_devcontainer_docker_test.go:123 | ja — `make test-docker` (Egress plus Docker, `apt-get update` im Container) | Begründung einer Allowlist-Aufnahme ohne Laufzeitbeleg |
| F-4 | LOW | Die Slice-Planung (§3) führt `harness/replay.md` als zu aktualisieren („Golden-Zeile“). Der Commit ändert die Datei nicht; die Zeile für `nestedRuntime: docker` nennt weder Paketquelle/Pin noch den Host-Fall `download.docker.com`. | Slice §3, `harness/replay.md` §Current Generator Golden Coverage | harness/replay.md:52 | ja — Review-only | Plan-/Replay-Tabelle driftet gegen Diff |
| F-5 | LOW | Die Prüfung nimmt nur den ersten `fpr`-Eintrag (`awk … exit`). `signed-by=` vertraut allen Schlüsseln der Datei. Enthielte die geladene `gpg`-Datei weitere Schlüssel, würden sie nicht geprüft. Die Voraussetzung ist ein kompromittierter TLS-Pfad oder eine kompromittierte Quelle; Pin und TLS bleiben die einzige Barriere. | [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl:31 | nein — Review-only | Pin prüft nur den ersten Schlüssel der Datei |
| F-6 | LOW | Der Doc-Kommentar von `egressHosts` nennt weiter „default base + registries (podman, docker)“ und erwähnt den neuen Host nicht, obwohl er Teil der Zusammensetzung ist. | Maintainability, `AGENTS.md` §3.7 (Kommentar beschreibt, was da ist) | internal/hexagon/application/devcontainer_egress.go:98 | nein — Review-only | Kommentar zur Zusammensetzung unvollständig |
| F-7 | INFO | Zwei Builds zu verschiedenen Zeiten liefern verschiedene Docker-Versionen; ein Versions-Pin ist in [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) als Re-Evaluierungs-Trigger vermerkt. Docker Desktop ist weiter ungemessen (Folge-Slice [slice-v1-sandbox-docker-nachholmessung-desktop](../plan/planning/open/slice-v1-sandbox-docker-nachholmessung-desktop.md)). Zuständigkeit: Validator beziehungsweise Folge-Slice. | [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) §Konsequenzen | docs/plan/adr/0019-docker-daemon-aus-docker-apt-repo.md | nein — Review-only | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Dockerfile.tmpl, Shell-Syntax und Quoting | geprüft, ohne Befund: `sh -n` des extrahierten RUN-Blocks ok; `$(…)` und `"$VERSION_CODENAME"` sind für `text/template` und Dockerfile-RUN unkritisch (keine `{{ }}`) |
| Fingerprint-Abbruch und `&&`-Semantik | geprüft, ohne Befund: schlägt `gpg` fehl oder weicht der Fingerprint ab, ist `test` ungleich 0 und die Kette bricht vor `docker.list`/`apt-get` ab; leere `gpg`-Ausgabe ist ebenfalls ungleich dem Pin |
| `curl`/`gpg`, Architektur, Suite, Gruppe | geprüft, ohne Befund: `curl`/`gpg` im Base-Image laut Messung in [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md), `dpkg --print-architecture`, `VERSION_CODENAME`, `.asc` mit `signed-by`, `usermod -aG docker vscode` nach dem Paketinstall (Gruppe aus dem `docker-ce`-Paket) |
| Kein Fremdskript | geprüft, ohne Befund: es wird nur ein APT-Repo mit geprüftem Schlüssel eingetragen, kein Installationsskript ausgeführt ([`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)); kein Host-Socket-Mount |
| devcontainer_egress.go | geprüft, ohne Befund: Host nur bei `nestedDocker`, Menge plus `sort.Strings` (sortiert, dedupliziert), Nutzer-`egress.allow` unverändert hinzugefügt; Golden Case `DockerAptHostOnlyWithDocker` pinnt docker/podman und den Nutzereintrag |
| [SPEC-019](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) und Handbuch | geprüft, ohne Befund: Tabellenzeile `nestedRuntime: docker` stimmt mit Code, Handbuch und CHANGELOG überein |
| Golden Cases NotGenerated | geprüft, ohne Befund: bei `none`/`podman` kein `docker-ce` und kein `download.docker.com` im Dockerfile |
| E2E `overlayDriver` | geprüft, ohne Befund: `HasPrefix` auf `overlay2`/`overlayfs` entspricht Entscheidung 4 in [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md); Prüfung an allen drei Stellen ersetzt; beide Treiber sind im Test dokumentiert |
| Doku-Konsistenz, Altaussagen | geprüft, ohne Befund: Grep nach „Debian-Pakete“, `docker.io`, `docker-cli`, `docker-buildx`, Altversionen außerhalb von Historie ([0.8.0], [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md), Reviews, done/) ohne Treffer; CHANGELOG [Unreleased]/Changed vorhanden; ADR-Index trägt ADR 0019 |
| Decken-Regel | geprüft, ohne Befund: Änderungen in `spec/spezifikation.md` verweisen nicht auf ADR, Slice oder Roadmap |
| [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) Format und Immutabilität | geprüft, ohne Befund: `Schärft:`-Feld vorhanden; [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) im Commit unverändert (letzte Änderung `6d726e0`); Ersatz/Ergänzung ist in [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) deklariert |
| Sprachvertrag, Exit-Codes, Architektur | geprüft, ohne Befund: keine neue CLI-Ausgabe, keine neuen Sentinels, keine Importänderung; Dateisicherheit unberührt (Managed Block unverändert) |
| Suppression/Gates | geprüft, ohne Befund: kein `//nolint`, keine Gate-Lockerung |

Nicht geprüft:

- Tatsächlicher Image-Build, Daemon-Start und `make test-docker` — Reviewer führt keine privilegierten Container aus; Evidence ist Aufgabe des Verifiers.
- `make gates` — nicht ausgeführt (nur `make docs-check` im Rahmen dieses Reports).
- Live-Gültigkeit des Fingerprints und Paketstände (29.9/0.38/5.6) — kein Netzzugriff auf `download.docker.com`; Angaben aus [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) übernommen.
- Docker Desktop, Podman als Host-Engine, VS Code/Codespaces — laut [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) ungemessen, Folge-Slice.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Sicherheitsrelevante Eigenschaft nur als Substring gepinnt · Integrationstest pinnt die neue Eigenschaft nicht · Plan-/Replay-Tabelle driftet gegen Diff

## Verdikt

**Merge-blockierend:** nein — keine HIGH-Findings; die drei MEDIUM-Findings betreffen Testtiefe (kein Produktfehler gefunden) und sind vor Slice-Closure zu klären oder in der Closure-Notiz mit Ausgang zu führen.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation.
