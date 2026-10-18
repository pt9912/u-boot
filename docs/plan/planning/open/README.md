# docs/plan/planning/open

Backlog: Slice- und Tranchen-Pläne, die noch nicht für die nächste
Iteration eingeplant sind.

Übergang nach `next/` per `git mv`, sobald ein Artefakt priorisiert
wird. Siehe `docs/plan/planning/README.md` für Lifecycle und
Dateiname-Konventionen.

## Bestand

| Datei | Gegenstand |
| ----- | ---------- |
| [`slice-cr-sandbox-fuse-startfehler`](../done/slice-cr-sandbox-fuse-startfehler.md) | Lastenheft-CR 0.12.0 → 0.13.0: ein fehlendes `/dev/fuse` auf dem Container-Host ist ein Startfehler der Engine, keine Degradationsstufe ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)); Lastenheft 0.13.0 angenommen (2026-10-18); Umsetzung in [`slice-v1-sandbox-fallback-host-ausweisung`](../done/slice-v1-sandbox-fallback-host-ausweisung.md) |
| [`slice-v1-handbuch-sshd-anleitung`](../done/slice-v1-handbuch-sshd-anleitung.md) | Handbuch: Schritt-für-Schritt-Anleitung für den `sshd` im Sandbox-Container (aktivieren, Zugang einrichten, verbinden, Sicherheitshinweis); keine neue Zusage; Teil von `v0.11.0` |
| [`slice-v1-e2e-status-ci-helfer`](../done/slice-v1-e2e-status-ci-helfer.md) | e2e-Test `StatusFiles_InspectorAgainstEngine` ist auf einem Runner ohne Image-Cache rot (Container-ID aus gemischter Ausgabe von `docker run -d`); Fix im Helfer `startLabeled`; geliefert (liegt in `done/`), Teil von `v0.11.0` |
| [`slice-v1-podman-aktuelle-version`](../done/slice-v1-podman-aktuelle-version.md) | Sandbox-Devcontainer mit `nestedRuntime: podman`: neueste Podman-Version aus gepinnter, geprüfter Quelle statt des ungepinnten Debian-Pakets (Quellen-Alternativen, Messung, Folge-ADR [`ADR-0023`](../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md)); geliefert (liegt in `done/`), Podman 6.1.4 aus `podman-static`; Teil von `v0.11.0` |
| [`slice-v1-release-cut-v0.11.0`](slice-v1-release-cut-v0.11.0.md) | Release-Cut `v0.11.0` nach dem Muster von [`slice-v1-release-cut-v0.10.0`](../done/slice-v1-release-cut-v0.10.0.md), mit Trockenlauf des Tap-Pfads vor dem stabilen Tag; startet erst, wenn `slice-v1-podman-aktuelle-version` in `done/` liegt und `integration-docker` auf dem Commit nach dem Push grün ist; Push und Tags nur nach Freigabe des Projektinhabers |
| [`slice-v1-cli-dry-run-log-would`](slice-v1-cli-dry-run-log-would.md) | Dry-Run: keine Info-Logzeilen in Vergangenheitsform ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run)); nach `v0.11.0`, nicht Teil des Release |
| [`slice-v1-harness-gofmt-und-depguard`](slice-v1-harness-gofmt-und-depguard.md) | Harness-Hygiene: Bestand formatieren und Formatter im Lint-Gate aktivieren; `verify-depguard` Fall `domain-isoliert` ohne Importzyklus; nach `v0.11.0`, nicht Teil des Release |
