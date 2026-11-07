# Maintainer-Dokumentation

Dokumente für Menschen, die an `u-boot` selbst arbeiten (Qualitätssicherung, CI, Release). Nutzer
von `u-boot` finden ihre Anleitung unter [`docs/user/`](../user/README.md).

- [`quality.md`](quality.md) — Quality-Gate- und Linter-Profil.
- [`devcontainer-arbeiten.md`](devcontainer-arbeiten.md) — Arbeit im Devcontainer des Repos: die Bedien-Targets `make dev-up`/`dev-shell`/`dev-down`/`dev-remove` (Instanzen), die Umgebung des Containers und die Gate-Lage.
- [`zusagen-und-messstaende.md`](zusagen-und-messstaende.md) — Zusagen, Nicht-Zusagen und Messstände zu den Aussagen des Benutzerhandbuchs, mit Fundstelle des Messprotokolls unter [`docs/messungen/`](../messungen/README.md).
- [`branch-protection.md`](branch-protection.md) — Einrichtung der Required Status Checks für `main`.
- [`releasing.md`](releasing.md) — Release-Ablauf (Tag, `publish`, GHCR, Binaries, Pakete, Homebrew-Tap, Nachziehen der Tap-Dateien, Fehlerbilder).
- [`cli-json-contract.md`](cli-json-contract.md) — Vertrag der maschinenlesbaren Ausgabe (`--json`, `--dry-run`, `--diff`): Schema, Code-Registry, Migrationsstand. Nutzersicht: [Benutzerhandbuch, Anhang „JSON-Ausgabe“](../user/benutzerhandbuch.md#json-ausgabe).
