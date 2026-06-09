# docs/user

User-facing Dokumentation für u-boot:

- Installationsanleitungen
- Quickstart-Guides
- Troubleshooting
- Beispiel-Workflows

Aktuell publiziert:

- [`examples.md`](examples.md) — Beispiel-Workflows als Kommando-Rezepte
  (Postgres-Stack, Keycloak+OTel, Devcontainer, Templates, CI/JSON,
  Cleanup, Config).
- [`cli-json-output.md`](cli-json-output.md) — `--json`/`--dry-run`/
  `--diff`-Envelope-Schema und Exit-Code-Matrix pro Subcommand.
- [`devcontainer-features.md`](devcontainer-features.md) — Devcontainer-
  Features + Drift-Doctor-Check.
- [`quality.md`](quality.md) — Quality-Gate-/Linter-Profil.
- [`branch-protection.md`](branch-protection.md) — Required-Checks-Setup.

`examples.md` führt **Befehle**, keinen committeten Output: die
byte-genaue Ausgabe ist in den Acceptance-/e2e-Tests gegen reale
Compose-Workloads gepinnt (Single-Source-of-Truth, siehe
[`slice-m6-docker-integrationstests`](../plan/planning/done/slice-m6-docker-integrationstests.md)).
Eine generierte+gegatete `examples/`-Variante (Doku-Variante A) ist als
Backlog-Idee denkbar, aber bewusst nicht umgesetzt — statische
Beispiel-Outputs würden ohne Gate driften.

Die kanonische pro-Command-Referenz bleibt `u-boot --help` /
`u-boot <command> --help` im gebauten Binary, die das Lastenheft
(`spec/lastenheft.md`) referenziert.
