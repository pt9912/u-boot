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

- ✅ **Schema (`featureSources.allow`, Spec-mandatiert):**
  `ubootYAMLDevcontainer.FeatureSources` ist implementiert (heute
  als Future-Field-Kommentar in `initproject.go:66-69`).
  `devcontainer.featureSources.allow` ist Liste von Quell-Strings,
  dedupliziert (`spec/lastenheft.md:1352`). Konkrete
  Failure-Cases der Validation: leerer Quell-String, fehlendes
  URL-Scheme, fehlende Host-Komponente — siehe T1 für die
  vollständige Failure-Tabelle. Diagnose-Code: `LH-FA-DEV-003`,
  Exit-Code `10` (`spec/lastenheft.md:1353`).
- ✅ **Feature-Aktivierung in u-boot.yaml (Plan-eigene Schema-
  Erweiterung, *nicht* Spec-mandatiert):** Spec definiert nur
  `featureSources.allow` (Quellen-Whitelist, §719). Für die
  deterministische Regeneration des `devcontainer.json`-Blocks
  braucht der Plan zusätzlich eine *Aktivierungs-Liste*, analog
  zur `services.<name>.enabled`-Konvention aus M5
  ([`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md)).
  Diese Erweiterung geht bewusst über den Spec-Wortlaut hinaus
  und wird hier als Plan-Hinzufügung explizit markiert.
  Design-Entscheidung in T0: Subkey `devcontainer.features`
  (Liste) vs. `services`-analog
  (`devcontainer.features.<name>.enabled` mit Pointer-`*bool`-
  Semantik wie `ubootYAMLService.Enabled` aus
  `initproject.go:54-56`).
- ✅ **CLI:** Mindestens einer der folgenden Pfade ist verdrahtet
  — Design-Entscheidung in T0; muss mit dem **Spec-Gültigkeits-
  bereich des `--allow-external-feature-sources`-Flags** abgeglichen
  sein. `spec/lastenheft.md:714-717` listet exakt drei Befehle,
  in denen das Flag gültig ist: `u-boot init --devcontainer`,
  `u-boot generate devcontainer`, `u-boot config set
  devcontainer.featureSources.allow`:
  - `u-boot add feature <name>` (analog zu `u-boot add postgres`,
    aber unter neuem `feature`-Namespace), oder
  - `u-boot devcontainer feature add <name>`, oder
  - `u-boot config set devcontainer.features.<name>.enabled true`
    + nachgelagertes `u-boot generate devcontainer`.

  T0-Sub-Entscheidung: Akzeptiert der neue Aktivierungs-Pfad
  `--allow-external-feature-sources` ebenfalls (vierte Stelle —
  Spec-Stretch, dafür UX-konsistent), oder läuft Aktivierung
  externer Quellen **ausschließlich** über die drei Spec-Pfade
  und der Aktivierungs-Pfad aktiviert nur Katalog-Features?
  Default-Annahme bis T0: zweitere Variante (Spec-treu).
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
- ✅ **Doctor-Integration Teil A (Spec-mandatiert):** `u-boot
  doctor` validiert Allowlist-Konformität: `devcontainer.features`
  referenziert nur Quellen aus `devcontainer.featureSources.allow`
  *oder* aus dem eingebauten Katalog (Spec-Pin
  [`spec/lastenheft.md:2394`](../../../../spec/lastenheft.md):
  „`u-boot doctor` enthält keinen `error` zu `devcontainer`-
  Konfiguration oder Feature-Quellen"). Check-ID:
  `devcontainer.features.allowlist` (Punktnotation analog zum
  bestehenden `doctorCheckID`-Enum in `doctor.go:69 ff.`).
- 🟡 **Doctor-Integration Teil B (über Spec hinaus,
  konditional):** Drift-Erkennung — `devcontainer.json` enthält
  die aktivierten Features tatsächlich (Managed-Block-Disziplin
  analog M5/M7). Check-ID: `devcontainer.features.drift`.
  **Carveout-Trigger:** Wenn nach Abschluss T1-T4 die kumulierte
  LOC > 800 (Add-on-Vergleichswert aus
  [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md) +
  [`slice-v1-keycloak`](../done/slice-v1-keycloak.md)), wird Teil
  B in einen Folge-Slice
  `slice-followup-devcontainer-features-drift-doctor` ausgelagert
  (Plan in `open/` anlegen *vor* T5-Start). Andernfalls bleibt
  Teil B im aktuellen Slice und wird in T5 implementiert.
- ✅ **Statischer Katalog** (ADR-0008-konform): u-boot bringt
  einen kuratierten Default-Satz Features mit (mindestens die
  Spec-Beispiele: Git, Docker CLI, Node, Java/SDKMAN, Go, C++,
  K8s-Tools), die ohne Allowlist-Eintrag aktivierbar sind. Externe
  Quellen jenseits davon brauchen die Allowlist. **Lesart des
  Spec-Begriffs „lokal hinterlegt"** (`spec/lastenheft.md:711`):
  hier verstanden als *in u-boot eingebauter Katalog* (statisches
  Go-Mapping analog `serviceCatalogue()` in
  `addservice_execute.go:234 ff.`, mit eigener
  `featureCatalogueEntry`-Struct). Die Alternativ-Lesart
  „Features im Repo-Pfad" ist bewusst ausgeschlossen und Trigger
  für ein eigenes Folge-Slice (siehe Out of Scope). Katalog-Key
  ist eine **neue Domain-Type `domain.FeatureName`** (T0-
  Entscheidung: nur `name` oder Tuple `(name, version)`, weil
  Spec-Beispiel `ghcr.io/devcontainers/features/<name>:<version>`
  Versionspinning vorsieht).
- ✅ **Spec-Pin:** `internal/hexagon/application/acceptance_test.go`
  deckt `LH-FA-DEV-003` ab. Test-Naming-Konvention analog zum
  bestehenden `TestLHAK###_<Title>`-Muster (siehe
  `acceptance_test.go:39, 98`):
  `TestLHFADEV003_AllowlistEnforcement` (negativer Pfad — externe
  Quelle ohne Allowlist → Exit-Code `10`) +
  `TestLHFADEV003_CatalogueActivation` (positiver Pfad — Katalog-
  Feature aktivieren + `devcontainer.json` enthält Eintrag).
  Docker-e2e in `internal/e2e/` ist **nicht** erforderlich:
  Devcontainer-Features werden nicht zur Laufzeit gegen einen
  echten Docker-Stack geprüft, sondern strukturell im
  generierten JSON.
- ✅ **Doku:** README (EN + DE) erweitert um Feature-Beispiel;
  `docs/user/devcontainer-features.md` (oder vergleichbar) listet
  Katalog + Allowlist-Mechanik.

## Tranchen (Skizze, wird beim Übergang nach `next/` verfeinert)

| T   | Inhalt (Skizze) |
| --- | --------------- |
| T0  | **Discovery / Design.** Vier Design-Entscheidungen festzurren: (a) CLI-Pfad (`add feature` vs. `devcontainer feature add` vs. nur `config set`) **inkl. Spec-§714-717-Verträglichkeit**: akzeptiert der neue Aktivierungs-Pfad `--allow-external-feature-sources` (vierte Stelle) oder ausschließlich die drei Spec-Pfade? (b) u-boot.yaml-Schema für aktivierte Features (Liste vs. Map-analog `services.<name>.enabled` mit `*bool`-Pointer-Semantik); (c) Katalog-Format (Go-Struct analog `addservice_execute.go:serviceCatalogue()` + `serviceCatalogueEntry`-Typ) inkl. Domain-Type-Frage **`domain.FeatureName` mit oder ohne Version-Slot**; (d) Managed-Block-Design: Erweiterung des bestehenden `init`-Blocks in `devcontainer.json.tmpl` oder zweiter, eigenständiger Managed-Block (Idempotenz- und Re-Init-Implikationen). Ergebnis: kurzer Design-Memo oder Mini-ADR; LOC-Schätzung pro Tranche wird hier nachgereicht. |
| T1  | **Schema-Erweiterung.** `ubootYAMLDevcontainer.FeatureSources` aus dem `Future fields`-Kommentar in `initproject.go:66-69` ziehen; YAML-Codec + Schema-Validierung. **Failure-Tabelle:** leerer Quell-String, fehlendes URL-Scheme (`http://`/`https://`/`oci://`), fehlende Host-Komponente, doppelter Eintrag (silent-dedupe gemäß `spec/lastenheft.md:1352`). Tests in `application/*_test.go`. Optional in T1 (sonst T2): Domain-Type `domain.FeatureName` gemäß T0-Entscheidung. |
| T2  | **Feature-Katalog.** Statischer Go-Katalog (analog `addservice_execute.go:234 ff. serviceCatalogue()`) mit Spec-Beispielen (Git, Docker CLI, Node, Java/SDKMAN, Go, C++, K8s-Tools). Quell-URLs gepinnt (`ghcr.io/devcontainers/features/<name>:<version>`). |
| T3  | **Generator-Patch.** `devcontainer.json.tmpl` gemäß T0-(d)-Entscheidung: bestehender `init`-Managed-Block erweitert *oder* zweiter Managed-Block `features`. JSON-Determinismus (sortierte Keys, stabile Reihenfolge); `generate devcontainer` idempotent. |
| T4  | **CLI-Subkommando** gemäß T0-(a)-Entscheidung. Allowlist-Enforcement-Pfad (Exit-Code `10`, Spec-konforme Fehlermeldung). `--allow-external-feature-sources <quelle>[,<quelle>...]`-Parser: Komma-Trennung, Whitespace-Toleranz (trim pro Element), Multi-Flag-Vorkommen kumulieren (kein last-wins, analog `--with-deps`-Pattern aus M5), Duplikate gegen bestehende Allowlist silent-dedupen (`spec/lastenheft.md:1352`). |
| T5  | **Doctor-Checks (Teil A obligatorisch; Teil B konditional).** Check-ID `devcontainer.features.allowlist` (Spec-mandatiert) ist Pflicht. Check-ID `devcontainer.features.drift` (Drift-Erkennung) gemäß AK-Carveout-Trigger: in diesem Slice implementieren, wenn T1-T4-Summe ≤ 800 LOC; sonst Folge-Slice `slice-followup-devcontainer-features-drift-doctor` (in `open/` anlegen *vor* T5-Start, damit kein Plan-Loch entsteht). Severity-Klassifikation analog M5. |
| T6  | **E2E + Spec-Pin.** `internal/hexagon/application/acceptance_test.go`: `TestLHFADEV003_CatalogueActivation` (positiver Pfad — Feature aktivieren + `devcontainer.json` enthält Eintrag) + `TestLHFADEV003_AllowlistEnforcement` (negativer Pfad — externe Quelle ohne Allowlist → Exit-Code `10`). Docker-e2e in `internal/e2e/` nicht erforderlich (kein Laufzeit-Stack zu prüfen). |
| T7  | **Doku-Closure.** READMEs (EN + DE), `docs/user/devcontainer-features.md`, CHANGELOG `## [Unreleased]`-Eintrag (Conventional-Commit-Stil: `feat(devcontainer): Devcontainer-Features-Allowlist und Katalog (LH-FA-DEV-003)`), Slice `open/` → `done/` **mit DoD-Hash-Line direkt** (kein `git log --grep`-Platzhalter, vgl. `feedback-done-slice-dod-hash`), Roadmap-Status-Update. |

LOC-Schätzung pro Tranche wird in T0 nachgereicht (analog zum
v0.3.0-Add-on-Slice-Stil); die kumulierte T1-T4-Summe ist
gleichzeitig Trigger-Wert für den Carveout-Entscheid an T5
(Schwelle 800 LOC, Vergleichswert aus
[`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md) +
[`slice-v1-keycloak`](../done/slice-v1-keycloak.md)).

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
  (Add-on-Pattern + `enabled`-State-Machine als Vorbild für die
  `devcontainer.features.<name>.enabled`-Disziplin — sofern T0
  diesen Schema-Pfad wählt; abweichend: keine Managed-Block-
  Drift-Zustände, weil `devcontainer.json` einfacher zu klassi-
  fizieren ist als ein Compose-Service-Block),
  [`slice-v1-keycloak`](../done/slice-v1-keycloak.md)
  (Service-Catalogue-Erweiterung um `serviceCatalogueEntry`-Typ
  als Vorlage für `featureCatalogueEntry`),
  [`slice-m7-generate`](../done/slice-m7-generate.md)
  (Managed-Block-Disziplin in `devcontainer.json`, inkl. der
  T0-(d)-Frage Erweiterung vs. zweiter Block).
- Roadmap:
  [`roadmap.md`](../in-progress/roadmap.md) §v0.4.0.
- Phase: V1, geplant für v0.4.0-Bündelung mit weiteren
  V1-Generators (`u-boot logs`, `--json`/`--dry-run`); deren
  Slice-Pläne existieren noch nicht (Roadmap-Stichworte) und
  folgen separat.
