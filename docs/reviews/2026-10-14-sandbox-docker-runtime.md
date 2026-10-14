# Review-Report: [`slice-cr-sandbox-docker-runtime`](../../docs/plan/planning/done/slice-cr-sandbox-docker-runtime.md) (+ drei Umsetzungs-Slices) — 2026-10-14

**Review-Art:** Code — Diff gegen Plan (Abnahmekriterien 1–7 des CR-Slice), Spec ([`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)), [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) und Hard Rules.

**Gegenstand:** Diff-Range `fa6aa1b..6d726e0` (Commits `effcf99`, `8db447c`, `6d726e0`): `devcontainer.sandbox.nestedRuntime: docker`.

**Skill:** `.harness/skills/reviewer.md` @ 6d726e0 · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-14

**Eingangs-Kontext:**

- [`slice-cr-sandbox-docker-runtime`](../../docs/plan/planning/done/slice-cr-sandbox-docker-runtime.md) (§1 Abnahmekriterien 1–7, §2 DoD), [`slice-v1-sandbox-docker-generator`](../../docs/plan/planning/done/slice-v1-sandbox-docker-generator.md), [`slice-v1-sandbox-docker-degradation-doctor`](../../docs/plan/planning/done/slice-v1-sandbox-docker-degradation-doctor.md), [`slice-v1-sandbox-docker-integrationstest`](../../docs/plan/planning/done/slice-v1-sandbox-docker-integrationstest.md)
- [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) (Schärft [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)), [ADR-0012](../plan/adr/0012-devcontainer-egress-firewall.md), [ADR-0014](../plan/adr/0014-nested-podman-sandbox-devcontainer.md), [ADR-0015](../plan/adr/0015-sandbox-volumes-pro-instanz.md)
- [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)/-007/-008`, [`[`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`[`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)`](../../spec/lastenheft.md#lh-lese-002--sprache), [`[`LH-FA-PROJDOCS-006`](../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell)`](../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell)
- `AGENTS.md` (Hard Rules), `harness/review.md`, `harness/roles.md`, `harness/replay.md`

Zusätzlich lokal ohne `--privileged` geprüft: Paketlage im Basis-Image `mcr.microsoft.com/devcontainers/base:debian` (Debian 13): `docker.io`, `docker-cli`, `docker-buildx`, `docker-compose` (2.26.1, installiert das Plugin unter `/usr/libexec/docker/cli-plugins/`), Gruppe `docker` entsteht, `setsid`/`nohup` vorhanden.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Das Startskript räumt vor dem Hintergrundstart keine Reste eines früheren `dockerd` auf. Der Devcontainer-Init-Prozess beendet beim Stopp den Container, der Kernel beendet `dockerd` per SIGKILL, die PID-Datei (`/var/run/docker.pid`) bleibt im Container-Dateisystem. Nach Neustart prüft `dockerd` die alte PID; belegt in der frischen PID-Namespace schon ein anderer Prozess diese Nummer, bricht er mit „process with PID N is still running“ ab. Das Skript meldet dann „daemon exited at start (without --privileged?)“ und degradiert (`warn`) oder endet mit Exit 11 (`fail`), obwohl `--privileged` gesetzt ist. Es gibt keinen Test für den Neustart-Pfad (Stopp, Start, erneutes `u-boot-docker-init`). Der Befund ist aus dem Code und dem bekannten Verhalten von `dockerd` abgeleitet, nicht per Lauf bestätigt (kein privilegierter Lauf). | [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Degradation nur bei fehlender Fähigkeit), [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) §Entscheidung 2 | `internal/hexagon/application/templates/devcontainer/docker-init.sh.tmpl:23` | nein — nur ein privilegierter Stopp/Start-Lauf (`make test-docker`-Erweiterung) bestätigt es | Neustart-Pfad ohne Test, Stale-State des Daemons |
| F-2 | MEDIUM | `harness/replay.md` führt für `generate devcontainer` je einen Golden-Case-Eintrag für `podman` und `egress`, aber keinen für `nestedRuntime: docker`; die Tests `TestGenerateDevcontainer_SandboxDocker*` sind darin nicht verankert. | Hard Rule §3.8 (Generator-Änderungen: Replay/Golden, `harness/replay.md`); AK 4 | `harness/replay.md:51` | ja — Review-only, `make docs-check` prüft die Tabelle nicht | Replay-Tabelle driftet gegen neue Generator-Pfade |
| F-3 | LOW | Die Spezifikation sagt für `docker`: „`runArgs` enthält `--privileged` und sonst keine Fähigkeitsoptionen“. Mit `egress.enabled` erzeugt der Generator zusätzlich `--cap-add=NET_ADMIN` (pinnt `TestGenerateDevcontainer_SandboxEgressWithDocker`: zwei Einträge); die Spezifikation nennt diese Kombination nicht. | [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) | `spec/spezifikation.md:326`, `internal/hexagon/application/devcontainer_sandbox.go:118` | nein — Review-only | Spec-Satz gilt nur ohne Egress, nicht ausgewiesen |
| F-4 | LOW | Das Lastenheft sagt zu, dass `docker` „unter Docker (inkl. Colima und Docker Desktop)“ startet. [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) und Benutzerhandbuch führen Docker Desktop als nicht geprüft (Nachholmessung offen, AK 6). | [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) §Konsequenzen | `spec/lastenheft.md:447`, `docs/user/benutzerhandbuch.md:1086` | nein — Review-only; Verifier-Aufgabe bei AK 6 | Zusage vor Messung |
| F-5 | LOW | Die Wartezeit ist nicht hart auf 60 s begrenzt: Jede Schleifenrunde ruft `docker info` ohne Timeout auf und schläft 1 s; ein hängender Aufruf verlängert die Wartezeit unbegrenzt. Die Spezifikation sagt „nach spätestens 60 Sekunden“. | [`[`[`[`[`[`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Spezifikation §[`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer).a) | `internal/hexagon/application/templates/devcontainer/docker-init.sh.tmpl:27` | nein — Review-only | Zeitzusage „spätestens“ nicht durchgesetzt |
| F-6 | LOW | Nachziehstellen für `docker-init.sh` fehlen: `docs/maintainer/cli-json-contract.md` nennt in der Rollback-Beschreibung nur `sandbox-init.sh`/`egress-init.sh`; der Abschnitt „Grenzen“ im Benutzerhandbuch nennt das Zurückbleiben von `sandbox-init.sh` nach Wechsel auf `none`, aber nicht das von `docker-init.sh`. | `AGENTS.md` §5 Regel 5 | `docs/maintainer/cli-json-contract.md:616`, `docs/user/benutzerhandbuch.md:1088` | nein — Review-only | Doku-Drift bei neuer generierter Datei |
| F-7 | LOW | DoD-Zustand im CR-Slice: „Folge-ADR“ ist unabgehakt, obwohl [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) (Accepted) und der ADR-Index-Eintrag im Diff vorliegen; „Doku-Update“ (CHANGELOG, ADR-Index, Roadmap) ist ebenfalls offen, die Inhalte stehen im Diff. | Hard Rule §3.7 (Zustandsfelder), `harness/verification.md` | [`slice-cr-sandbox-docker-runtime`](../plan/planning/done/slice-cr-sandbox-docker-runtime.md) (§2, Zeile 157 zum Reviewstand) | nein — Review-only; Verifier schließt | Zustandsfeld hinter Beleg |
| F-8 | INFO | `docker-compose` im Debian-Archiv ist Compose v2 und installiert das CLI-Plugin; `docker compose` funktioniert damit. Das Basis-Image `devcontainers/base:debian` ist unversioniert; ein Wechsel der Debian-Basis (Paketnamen/-stand) trifft den Dockerfile-Pfad. [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) nennt Debian 13 als Messgrundlage. | [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) | `internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl:27` | ja — `make test-docker` | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| AK 1 (`--privileged` in `runArgs`, Warnung, kein Socket-/Bind-Mount, Volume `<projekt>-docker-${devcontainerId}`) | geprüft, ohne Befund (`TestGenerateDevcontainer_SandboxDocker`, E2E-Mount-Prüfung) |
| AK 2/3 (Integrationstest, zwei Instanzen, getrennte Daemon-IDs) | geprüft, ohne Befund im Code des Tests; Lauf nicht wiederholt (siehe Nicht geprüft) |
| AK 4 (`none`/`podman` unverändert, `--privileged` nur bei `docker`) | geprüft, ohne Befund (`..._NotGenerated`) |
| AK 5 (Degradation: `warn` ohne Daemon, `fail` Exit 11, nie Podman) | geprüft, ohne Befund im Skript (Exit 11 an einer Stelle, Fehlerpfade Early-Exit und Timeout), E2E `..._Degradation` pinnt beide Policies |
| AK 7 (Doctor: ohne Sandbox warnt; mit Egress warnt) | geprüft, ohne Befund (`doctor_sandbox.go`, Tests) |
| Exit-Code-Klassifikation / Dual-Classifier | kein neuer Sentinel; ungültiger Wert nutzt `ErrInvalidSandboxSetting` (Exit 10); geprüft, ohne Befund |
| Decken-Regel (`spec/` verweist nicht abwärts auf ADR/Slice/Roadmap) | geprüft, ohne Befund (Diff von `lastenheft.md`/`spezifikation.md` enthält keine ADR-/Slice-Bezüge) |
| Architektur / Importregeln | nur `application`/`domain` berührt, keine neuen Importe; geprüft, ohne Befund |
| Kommentar-Regel §3.7 | Code-, Template- und Test-Kommentare sind Zusage/Kopplung/Rang-Zeiger im Indikativ; geprüft, ohne Befund |
| Sprachvertrag ([`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)) | Skript-Meldungen, Warnungen, Doctor-Texte Englisch; geprüft, ohne Befund |
| Dockerfile-/devcontainer.json-Ausgabe | gültig, `postStartCommand` bei Egress+Docker in Reihenfolge `egress-init && docker-init`, Mounts-Zweig korrekt (`else if`); geprüft, ohne Befund |
| Ordner-Lifecycle (`git mv` vs. Inhalt) | Verschiebung in `effcf99`, Link-/Inhaltsänderung in `8db447c`; geprüft, ohne Befund |
| [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) / ADR-Index / Immutabilität | neue ADR, keine Umschreibung einer Accepted-ADR; geprüft, ohne Befund |
| CHANGELOG, Benutzerhandbuch (Inhalt) | konsistent mit Lastenheft/Spezifikation außer F-4/F-6; geprüft |

Nicht geprüft:

- Echter Lauf von `make test-docker` / der `//go:build docker`-Tests — der Auftrag untersagt `--privileged`-Container; Begründung der Erwartung stützt sich auf die Messung in [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md).
- Neustart-/Stale-PID-Verhalten von `dockerd` (F-1) — kein privilegierter Lauf möglich.
- Docker Desktop, Podman als Host-Engine, Codespaces — laut [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) ungemessen.
- `make test`/`make lint`/`make gates` — laut Auftrag bereits grün, nicht erneut ausgeführt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 5 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Neustart-Pfad ohne Test, Stale-State des Daemons · Replay-Tabelle driftet gegen neue Generator-Pfade · Zusage vor Messung · Doku-Drift bei neuer generierter Datei · Zustandsfeld hinter Beleg

## Verdikt

**Merge-blockierend:** nein — keine HIGH-Findings; F-1 und F-2 sind vor Closure des Slice zu klären (F-1: Neustart-Pfad belegen oder absichern; F-2: Replay-Eintrag).

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation.
