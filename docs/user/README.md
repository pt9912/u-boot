# docs/user

User-facing Dokumentation für u-boot:

- Installationsanleitungen
- Quickstart-Guides
- Troubleshooting
- Beispiel-Workflows

Stand M8: nur zwei meta-Themen sind aktuell hier publiziert
([`quality.md`](quality.md) und
[`branch-protection.md`](branch-protection.md)); pro-Command-
Guides für die sieben verdrahteten Subcommands (`init`, `doctor`,
`add`, `up`, `down`, `generate`, `config`) folgen, sobald sich
Quickstart-Beispiele gegen reale Compose-Workloads validieren
lassen — siehe das Carveout-Slice
[`slice-m6-docker-integrationstests`](../plan/planning/done/slice-m6-docker-integrationstests.md)
(Done) für die End-to-End-Pins, die diesen Guides als
Quelle dienen werden.

Bis dahin ist die kanonische User-Dokumentation `u-boot --help`
und `u-boot <command> --help` im built Binary, die das Lastenheft
(`spec/lastenheft.md`) und die Slice-Pläne in
[`docs/plan/planning/done/`](../plan/planning/done/) referenzieren.
