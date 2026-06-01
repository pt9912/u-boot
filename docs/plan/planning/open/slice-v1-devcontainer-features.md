# Slice V1: Devcontainer-Features / Toolchains (`LH-FA-DEV-003`)

> **Status:** scheduled für v0.4.0 — Spec ✅, ADR-Hook ✅,
> Code-Anker ✅, Implementation ausstehend. Plan ist eine erste
> Version; Tranchen-Schnitt wird beim Übergang nach `next/`
> verfeinert.

## Auslöser

Der heutige Add-on-Katalog kennt nur **Laufzeit-Services**
(Postgres, Keycloak, OpenTelemetry) — alles Compose-Services mit
Image-Pin, Healthcheck und Port. Eine reproduzierbare
Entwicklungsumgebung braucht aber auch **Toolchains** (Java,
Gradle, Go, Node.js, C++, …). Ohne sie ist `u-boot init` zwar
ein lauffähiger Compose-Stack, aber kein einsatzbereites
Dev-Setup — der User muss die Sprach-Toolchain selbst im
Devcontainer-`Dockerfile` oder per `features:`-Section
nachziehen.

Die Spec hat das vorgesehen, aber die Anforderung ist bisher
nicht implementiert:

- **`LH-FA-DEV-003`** (Priorität V1) listet explizit Git, Docker
  CLI, Node.js, Java, SDKMAN, PostgreSQL Client, Kubernetes Tools
  als Beispiele und legt das Security-Modell fest
  (`--allow-external-feature-sources`,
  `devcontainer.featureSources.allow`, `--yes` reicht nicht).
- **ADR-0008** §78 nennt die Allowlist ausdrücklich als einzigen
  ausnahme-pflichtigen externer-Code-Pfad ausserhalb des
  Compose-Stacks.
- **Code-Anker:** `internal/hexagon/application/initproject.go:66-69`
  markiert die Future-Field-Stelle im YAML-Schema
  (`Future fields:`-Kommentarblock mit
  `FeatureSources ubootYAMLFeatureSources`).
- **Schema-Skizze:** `spec/lastenheft.md:1340-1343` zeigt die
  YAML-Form (`devcontainer.featureSources.allow: [URL, …]`).

## Aufhebungsbedingung

`u-boot init --devcontainer` plus ein neuer Feature-Flow erzeugt
einen `.devcontainer/devcontainer.json`, der eine vom User
gewählte Toolchain als `features:`-Section enthält; die Allowlist
in `u-boot.yaml` ist gepflegt; `u-boot doctor` läuft ohne Errors;
`u-boot generate devcontainer` ist idempotent (managed-block-
Stabilität analog zu M7).

## Akzeptanzkriterien

- ✅ **Schema:** `ubootYAMLDevcontainer.FeatureSources` ist
  implementiert (heute als Future-Field-Kommentar in
  `initproject.go:66-69`). `devcontainer.featureSources.allow` ist
  Liste von URL-Strings, dedupliziert, mit Validation gegen
  leere/ungültige Quellen (`spec/lastenheft.md:1350-1353`).
- ✅ **Feature-Auswahl in u-boot.yaml:** Ein neuer Schema-Pfad
  hält die *aktivierten* Features (nicht nur die *erlaubten
  Quellen*). Design-Entscheidung in T0: Subkey
  `devcontainer.features` (Liste) vs. `services`-analog
  (`devcontainer.features.<name>.enabled`).
- ✅ **CLI:** Mindestens einer der folgenden Pfade ist verdrahtet
  — Design-Entscheidung in T0:
  - `u-boot add feature <name>` (analog zu `u-boot add postgres`,
    aber unter neuem `feature`-Namespace), oder
  - `u-boot devcontainer feature add <name>`, oder
  - `u-boot config set devcontainer.features.<name>.enabled true`
    + nachgelagertes `u-boot generate devcontainer`.
- ✅ **Allowlist-Enforcement:** Ohne explizit erlaubte Quelle
  bricht der Versuch, eine externe Feature-Quelle zu aktivieren,
  mit `code LH-FA-DEV-003` / Exit-Code `10` ab. `--yes` reicht
  nicht (`LH-NFA-SEC-004`). `--allow-external-feature-sources`-
  Flag-Argumente werden in `u-boot.yaml`
  (`devcontainer.featureSources.allow`) persistiert — Flag ist
  Erweiterung der Liste, nicht Ein-Aufruf-Override
  (`spec/lastenheft.md:719`: „als explizit freigegebene Liste in
  der Projektkonfiguration gespeichert").
- ✅ **Devcontainer-Template:** Der bestehende Managed-Block in
  `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl`
  wird um eine `features:`-Section erweitert (oder ein zweiter
  Managed-Block); JSON-Stabilität (sortierte Keys, deterministische
  Ausgabe) ist gewahrt.
- ✅ **Doctor-Integration:** `u-boot doctor` validiert (a)
  Allowlist-Konformität: `devcontainer.features` referenziert nur
  Quellen aus `devcontainer.featureSources.allow` *oder* aus dem
  eingebauten Katalog (Spec-mandatiert via `LH-FA-DEV-003` +
  Doctor-Pin [`spec/lastenheft.md:2394`](../../../../spec/lastenheft.md)).
  (b) Drift-Erkennung: `devcontainer.json` enthält die aktivierten
  Features tatsächlich — **über Spec hinaus, in Anlehnung an M5/M7
  Managed-Block-Disziplin** (kann zur v0.4.0-Last-Reduktion in
  Folge-Slice verschoben werden, wenn Tranchen-Budget knapp wird).
- ✅ **Statischer Katalog** (ADR-0008-konform): u-boot bringt
  einen kuratierten Default-Satz Features mit (mindestens die
  Spec-Beispiele: Git, Docker CLI, Node, Java/SDKMAN, Go, C++,
  K8s-Tools), die ohne Allowlist-Eintrag aktivierbar sind. Externe
  Quellen jenseits davon brauchen die Allowlist. **Lesart des
  Spec-Begriffs „lokal hinterlegt"** (`spec/lastenheft.md:711`):
  hier verstanden als *in u-boot eingebauter Katalog* (statisches
  Go-Mapping). Die Alternativ-Lesart „Features im Repo-Pfad" ist
  bewusst ausgeschlossen und Trigger für ein eigenes Folge-Slice
  (siehe Out of Scope).
- ✅ **Spec-Pin:** `acceptance_test.go` deckt `LH-FA-DEV-003` ab
  (Test-Skelett `TC-DEV-003`); Allowlist-Verstoß-Pfad gepinnt.
- ✅ **Doku:** README (EN + DE) erweitert um Feature-Beispiel;
  `docs/user/devcontainer-features.md` (oder vergleichbar) listet
  Katalog + Allowlist-Mechanik.

## Tranchen (Skizze, wird beim Übergang nach `next/` verfeinert)

| T   | Inhalt (Skizze) |
| --- | --------------- |
| T0  | **Discovery / Design.** Drei Design-Entscheidungen festzurren: (a) CLI-Pfad (`add feature` vs. `devcontainer feature add` vs. nur `config set`); (b) u-boot.yaml-Schema für aktivierte Features; (c) Katalog-Format (Go-Struct analog `addservice_execute.go:serviceCatalogue()` + `serviceCatalogueEntry`-Typ?). Ergebnis: kurzer Design-Memo oder Mini-ADR. |
| T1  | **Schema-Erweiterung.** `ubootYAMLDevcontainer.FeatureSources` aus dem `Future fields`-Kommentar in `initproject.go:66-69` ziehen; YAML-Codec + Schema-Validierung; Edge-Cases (leere Liste, ungültige URLs). Tests in `application/*_test.go`. |
| T2  | **Feature-Katalog.** Statischer Go-Katalog (analog `addservice_execute.go:serviceCatalogue()`) mit Spec-Beispielen (Git, Docker CLI, Node, Java/SDKMAN, Go, C++, K8s-Tools). Quell-URLs gepinnt (`ghcr.io/devcontainers/features/<name>:<version>`). |
| T3  | **Generator-Patch.** `devcontainer.json.tmpl` um `features:`-Section erweitern (oder zweiter Managed-Block); JSON-Determinismus (sortierte Keys, stabile Reihenfolge); `generate devcontainer` idempotent. |
| T4  | **CLI-Subkommando** gemäß T0-Entscheidung. Allowlist-Enforcement-Pfad (Exit-Code 10, Spec-konforme Fehlermeldung). `--allow-external-feature-sources <quelle>[,<quelle>...]` parsen. |
| T5  | **Doctor-Checks.** Mindestens zwei: `devcontainer.features.in-allowlist` + `devcontainer.json.features-drift`. Severity-Klassifikation analog M5. |
| T6  | **E2E + Spec-Pin.** `acceptance_test.go` `TC-DEV-003`: positiver Pfad (Feature aktivieren + `devcontainer.json` enthält Eintrag) + negativer Pfad (externe Quelle ohne Allowlist → Exit-Code 10). |
| T7  | **Doku-Closure.** READMEs (EN + DE), `docs/user/devcontainer-features.md`, CHANGELOG `## [Unreleased]`-Eintrag, Slice `open/` → `done/`, Roadmap-Status-Update. |

LOC-Schätzung wird in T0 nachgereicht (analog zum
v0.3.0-Add-on-Slice-Stil).

## Out of Scope

- **Eigene/lokale Features** (Custom-Feature im Repo-Pfad statt
  externer Quelle): geht über `LH-FA-DEV-003` hinaus; eigenes
  Folge-Slice mit eigenem Trigger.
- **Sprach-spezifische Build-Aktionen** (Gradle-Wrapper anlegen,
  Go-Module-Init, npm-init, …): das ist Template-Job
  (`LH-FA-TPL-*`), nicht Devcontainer-Feature-Job. Klare
  Trennung: u-boot pinnt die Toolchain, das Template (oder der
  User) ruft sie auf.
- **Feature-Version-Updates** (Renovate/Dependabot-Hook für
  Katalog-Pins): nice-to-have; eigener Folge-Slice nach
  Erstauslieferung.
- **Migration bestehender Devcontainer-Configs ohne managed-Block:**
  M3-Re-Init-Disziplin gilt; User mit selbst-gepflegtem
  `devcontainer.json` müssen `--force`/`--backup` benutzen oder
  `features:` manuell migrieren.

## Bezug

- Spec: `LH-FA-DEV-003` ([`spec/lastenheft.md:692`](../../../../spec/lastenheft.md))
  + Schema-Skizze [`spec/lastenheft.md:1340-1353`](../../../../spec/lastenheft.md)
  + Doctor-Pin [`spec/lastenheft.md:2394`](../../../../spec/lastenheft.md).
- ADR: [ADR-0008 §78](../../adr/0008-plugin-system-statisch.md)
  — Allowlist als einziger ausnahme-pflichtiger externer-Code-Pfad.
- Code-Anker: `internal/hexagon/application/initproject.go:66-69`
  (`Future fields:`-Kommentarblock mit `FeatureSources …`).
- Generator-Pfad: `internal/hexagon/application/generate.go:705`
  (devcontainer.json-Mapping) +
  `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl`
  (heute nur `init`-Managed-Block).
- Vorbild-Slices:
  [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md)
  (Add-on-Pattern), [`slice-v1-keycloak`](../done/slice-v1-keycloak.md)
  (Service-Catalogue-Erweiterung),
  [`slice-m7-generate`](../done/slice-m7-generate.md)
  (Managed-Block-Disziplin in `devcontainer.json`).
- Roadmap:
  [`roadmap.md`](../in-progress/roadmap.md) §v0.4.0.
- Phase: V1, geplant für v0.4.0-Bündelung mit weiteren
  V1-Generators (`u-boot logs`, `--json`/`--dry-run`); deren
  Slice-Pläne existieren noch nicht (Roadmap-Stichworte) und
  folgen separat.
