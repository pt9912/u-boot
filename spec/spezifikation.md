# Spezifikation — u-boot

**Bezug zum Lastenheft:** Diese Spezifikation präzisiert die in
[`lastenheft.md`](lastenheft.md) formulierten Anforderungen (`LH-*`-Kennungen). Bei
Konflikt gewinnt das Lastenheft — präzisieren ja, erweitern nie.

**Rolle:** Technik-Stratum — verbindlich, aber ohne Änderung des Lastenhefts
fortschreibbar. Sie verweist nur aufwärts auf das Lastenheft, nie auf ADRs, Slices,
Carveouts oder die Roadmap.

**Kennungen:** Eine Verfeinerung einer einzelnen Anforderung trägt deren Kennung mit
Buchstabensuffix (`<Anforderungs-Kennung>.a`, `.b`, …). Alles, was keine einzelne
Anforderung verfeinert (Schemata, Defaults, Fehler-Codes, Metrik-Felder, externe
Verträge), trägt `SPEC-<NNN>`: dreistellig, fortlaufend je Datei, nicht je Abschnitt.
Eine `SPEC-*` ist eine Adresse, keine Anforderung.

---

## 1. Algorithmen und Datenflüsse

Verfeinerungen einzelner Anforderungen: Ablauf, Zustandsfolgen, Entscheidungsregeln.
Tatsächlicher Code gehört nicht hierher.

### LH-FA-INIT-002.a — Normalisierung des abgeleiteten Projektnamens

Verfeinert [`LH-FA-INIT-002`](lastenheft.md#lh-fa-init-002--projektname).

1. Der Basisname des Arbeitsverzeichnisses wird auf Kleinbuchstaben gesetzt.
2. Alle Zeichen außer `a-z`, `0-9` und `-` werden auf `-` abgebildet.
3. aufeinanderfolgende `-` werden zu einem einzelnen `-` zusammengeführt.
4. führende und nachgestellte `-` sowie Leerzeichen werden entfernt.
5. Die Länge wird auf 1 bis 63 Zeichen begrenzt.
6. Nach Kürzung auf 63 Zeichen wird erneut auf führende/nachgestellte `-` geprüft und diese notfalls entfernt.
7. Anschließend wird der Name gegen die Validierung in [`LH-FA-INIT-006`](lastenheft.md#lh-fa-init-006--projektnamen-validierung) geprüft.

### LH-FA-DEV-004.a — Übergabe der Benutzer-ID an den Image-Build

Verfeinert [`LH-FA-DEV-004`](lastenheft.md#lh-fa-dev-004--benutzerrechte).

- Der Wert wird als Build-Argument `USER_UID` an den Image-Build übergeben; der Container-Benutzer wird mit dieser UID angelegt.

### LH-FA-DIAG-002.a — Prüfung von Docker und Docker-kompatiblen Drop-ins

Verfeinert [`LH-FA-DIAG-002`](lastenheft.md#lh-fa-diag-002--lokale-voraussetzungen-prüfen).

- Docker installiert (Mindestversion: 24.0.0 oder neuer) **oder**
  ein Docker-API-kompatibler Drop-in (z. B. Podman ≥ 4.0 mit
  aktivem `podman.socket` und `DOCKER_HOST` darauf gezeigt). Die
  Checks `docker.installed` / `docker.reachable` /
  `docker.compose.installed` shellen aus zum `docker`-Binary;
  Drop-ins, die nicht den Docker-Version-Format-Pin
  (`<major>.<minor>.<patch>`-Bereich 24.0 / 2.20) treffen,
  werden als `Severity: warn` („unrecognized version") emittiert,
  ohne den Exit-Code zu eskalieren. Eine formal getestete
  Podman-Variante folgt erst bei konkretem Bedarf — die heutige
  MVP-Pflicht ist Docker.

### LH-FA-BUILD-001.a — Stages des Multi-Stage-Dockerfiles

Verfeinert [`LH-FA-BUILD-001`](lastenheft.md#lh-fa-build-001--multi-stage-dockerfile-u-boot-repo).

- BuildKit-Direktive in der ersten Zeile: `# syntax=docker/dockerfile:1.7`.
- Pflicht-Stages im MVP:
  - `deps` – Modulauflösung (`go mod download`) als Cache-Layer.
  - `compile` – schnelles Compile-Feedback (`go build`) ohne Tests/Lint.
  - `test` – `go test ./...`.
  - `lint` – `golangci-lint run ./...`.
  - `coverage` – `go test -coverprofile` + Coverage-Gate gegen `COVERAGE_THRESHOLD`; Bootstrap-Verhalten und Scope sind in [`LH-FA-BUILD-008`](lastenheft.md#lh-fa-build-008--coverage-bootstrap) und [`LH-FA-BUILD-009`](lastenheft.md#lh-fa-build-009--repository-layout) definiert.
  - `build` – statisch gelinktes Binary (`CGO_ENABLED=0`, `-ldflags="-s -w"`).
  - `runtime` – minimales Endimage ([`LH-FA-BUILD-002`](lastenheft.md#lh-fa-build-002--runtime-stage-pflichten)).
- Jeder Stage ist ein eigenständiges Build-Ziel und wird per `docker build --target <stage>` einzeln baubar.

### LH-FA-BUILD-002.a — Eigenschaften der Runtime-Stage

Verfeinert [`LH-FA-BUILD-002`](lastenheft.md#lh-fa-build-002--runtime-stage-pflichten).

- Base-Image: `gcr.io/distroless/static-debian12:nonroot` (oder gleichwertig minimal und ohne Shell).
- Non-root-Benutzer (Distroless-`nonroot`-User, `USER 65532:65532`).
- `ENTRYPOINT` zeigt auf das im `build`-Stage erzeugte Binary; der konkrete Pfad (Empfehlung: `/usr/local/bin/u-boot`) ist im Dockerfile dokumentiert.
- OCI Image Labels gesetzt:
  - `org.opencontainers.image.source`
  - `org.opencontainers.image.description`
  - `org.opencontainers.image.licenses`
  - `org.opencontainers.image.title`
- Keine Build-Toolchain im Endimage; alle Build-Artefakte werden aus dem `build`-Stage per `COPY --from=build` übernommen.

### LH-QA-003.a — Komposition der CI-Pipeline (GitHub Actions)

Verfeinert [`LH-QA-003`](lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions).

- Workflow-Datei: [`.github/workflows/ci.yml`](../.github/workflows/ci.yml).
- Trigger: `pull_request` und `push` auf den Branch `main`.
- Drei Jobs, parallel, alle PR-blockierend (Required-Status-Checks im GitHub-UI nach dem ersten grünen Lauf zu konfigurieren; die Required-Status-Check-Liste muss die tatsächlichen Workflow-`name:`-Felder verwenden, nicht die Kurz-Keys). Die drei delegierten Make-Targets sind in [`LH-FA-BUILD-005`](lastenheft.md#lh-fa-build-005--makefile-mit-standard-targets) (MVP) und [`LH-FA-BUILD-006`](lastenheft.md#lh-fa-build-006--aggregator-targets) (V1) definiert; die PR-Blocking-**Pflicht** für alle drei kommt aus diesem [`LH-QA-003`](lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions)-Eintrag (MVP). Damit ist `make govulncheck` und `make image-scan` über [`LH-QA-003`](lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions) MVP-bindend, auch wenn die jeweiligen Make-Target-Definitionen unter [`LH-FA-BUILD-006`](lastenheft.md#lh-fa-build-006--aggregator-targets) (V1) liegen:
  - `gates (lint + test + coverage-gate)` — führt `make gates` aus (lint + test + coverage-gate, [`LH-FA-BUILD-005`](lastenheft.md#lh-fa-build-005--makefile-mit-standard-targets)/[`LH-FA-BUILD-006`](lastenheft.md#lh-fa-build-006--aggregator-targets)).
  - `security-gates (govulncheck)` — führt `make govulncheck` aus ([`LH-FA-BUILD-006`](lastenheft.md#lh-fa-build-006--aggregator-targets), MVP-Pflicht via diesem [`LH-QA-003`](lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions)).
  - `image-scan (trivy HIGH+CRITICAL)` — führt `make image-scan` aus (Trivy gegen das Runtime-Image, severity `HIGH,CRITICAL`, exit-code `1`; [`LH-FA-BUILD-006`](lastenheft.md#lh-fa-build-006--aggregator-targets), MVP-Pflicht via diesem [`LH-QA-003`](lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions)).
- Runner: `ubuntu-latest`. Keine Host-Go-Toolchain ([`LH-FA-BUILD-007`](lastenheft.md#lh-fa-build-007--docker-only-workflow)); der Runner braucht nur das vorinstallierte Docker + BuildKit.
- Actions sind **SHA-gepinnt** mit Tag-Kommentar (Supply-Chain-Härtung gegen Tag-Move). Pin-Hebung ist Routine; neuer Commit-SHA via `gh api repos/<owner>/<repo>/git/refs/tags/<tag>`.
- Top-Level `permissions: {}` (alle Tokens entzogen); jeder Job lockert auf das Minimum (Defense-in-Depth).
- Jeder Job mit `timeout-minutes` versehen (Empfehlung: 20).
- Die konkrete Workflow-Ausprägung muss die hier genannten
  Anforderungen und die Build-Target-Verträge erfüllen.

### LH-QA-004.a — Zusammensetzung des Lint-Profils

Verfeinert [`LH-QA-004`](lastenheft.md#lh-qa-004--linting-solid-nahes-lint-profil).

Profil-Komposition (29 Linter insgesamt):

- 5 Default-Linter (`govet`, `errcheck`, `staticcheck`, `unused`, `ineffassign`).
- 24 SOLID-nahe Zusatz-Linter (Komplexitäts-, Funktionslänge-, Interface-, Kopplungs- und Boundary-Signale). **`depguard`** für die Schicht-Regeln aus [`LH-FA-ARCH-003`](lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement) ist Teil dieser 24.

### LH-FA-CLI-005A.a — Detailregeln des Flags `--assume-existing`

Verfeinert [`LH-FA-CLI-005A`](lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung).

- Für `u-boot init` ist zusätzlich das Flag `--assume-existing` definiert (nicht global, nur für diesen Befehl):
  - Ohne `--assume-existing` wird eine implizite Erkennung als bestehendes Projekt im nicht-interaktiven Modus nicht automatisch akzeptiert.
  - Mit `--assume-existing` wird die implizite Erkennung als bestehendes Projekt in nicht-interaktiven Läufen akzeptiert.
  - `--yes` ist für diesen Sonderfall **nicht** ausreichend; die implizite Erkennung bleibt abgelehnt, wenn keine `--assume-existing` gesetzt ist.
  - Ohne `--assume-existing` und bei nicht-interaktivem Lauf ist die implizite Erkennung zwingend ablehnend und erzeugt einen fachlichen Fehler.
  - Der Fehlercode für diese Abweisung ist `10`.

### LH-FA-CLI-005A.b — Verhalten bei aktivierter Nicht-Interaktivität

Verfeinert [`LH-FA-CLI-005A`](lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung).

- Bei aktivierter Nicht-Interaktivität darf keine neue Rückfrage erzeugt werden:
  - mit `--no-interactive` bricht der Aufruf bei jeder offenen Bestätigungsfrage mit Exit-Code `2` ab,
  - mit `--yes` wird die vorgesehene Standardentscheidung deterministisch ausgeführt.
- Für bereits deterministische Ausführungspfade (keine relevante Rückfrage) ist das Verhalten in beiden Modi unverändert.

### LH-FA-CLI-005A.c — Auswertungsreihenfolge von `init` im nicht-interaktiven Modus

Verfeinert [`LH-FA-CLI-005A`](lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung).

Bei `u-boot init` gilt zusätzlich die feste Auswertungsreihenfolge im nicht-interaktiven Modus:

- ohne `--assume-existing`: keine implizite Annahme einer bestehenden Projekterkennung, deterministisch abbrechen (Exit-Code `10` bei bestehendem Projekt),
- mit `--assume-existing`: implizite Annahme als bestehendes Projekt (soweit kompatibel mit den übrigen Validierungen).

### LH-FA-CLI-005A.d — Auswertungslogik der Bestätigungsmodi

Verfeinert [`LH-FA-CLI-005A`](lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung).

Deterministische Auswertungslogik für bestätigungsrelevante Modi:

- `--yes` und `--no-interactive` sind exklusiv.
- `--no-interactive` erlaubt keinerlei Rückfragen. Alle Entscheidungswege müssen deterministisch sein oder mit [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes)-Code `2` abbrechen, wenn eine notwendige Bestätigung fehlt.
- `--yes` erlaubt deterministische Standardpfade ohne Nutzerinteraktion.
- `--force` und/oder `--backup` sind in nicht-interaktiven Läufen explizit zulässig, weil beide Modi deterministisch arbeiten.
- `--no-interactive` + `--force` erlaubt das Überschreiben ohne Rückfrage; dabei ist immer eine vollständige Zusammenfassung der betroffenen Pfade auszugeben.
- `--force` darf keine zusätzlichen Rückfragen erzeugen; die Sicherheitslogik beschränkt sich auf die Validierung der Eingabedaten.
- `--backup` ist optional. Wenn `--backup` gesetzt ist, dürfen Dateischutz-Szenarien mit automatischer Sicherung deterministisch abgearbeitet werden.
- Bei fehlender Möglichkeit zur sicheren automatischen Abarbeitung (z. B. fehlender verwalteter Block ohne `--backup` bei vollständig kontrolliertem Überschreiben) muss der Befehl mit Fehlercode `10` abbrechen.

### LH-FA-INIT-005.a — Schutzregeln für strukturierte Konfigurationsdateien

Verfeinert [`LH-FA-INIT-005`](lastenheft.md#lh-fa-init-005--überschreibschutz).

Zusätzliche Schutzregeln für strukturierte Konfigurationsdateien (`compose.yaml`, `.env.example`, `README.md`, `CHANGELOG.md`, `.devcontainer/devcontainer.json`):

- bestehende, nicht verwaltete Inhalte bleiben in `--force`-Ausführung erhalten.
- wird ein `U-BOOT MANAGED BLOCK` erkannt, darf bei `--force` nur dieser Block verändert werden. Das Markierungsformat pro Dateityp ist in [`LH-SA-FILE-002`](lastenheft.md#lh-sa-file-002--markierte-verwaltete-bereiche) definiert.
- für `.devcontainer/devcontainer.json` gilt der JSONC-Markerstil (`// BEGIN U-BOOT MANAGED BLOCK: <name>` / `// END U-BOOT MANAGED BLOCK: <name>`); für strikte JSON-Dateien ohne Kommentar-Support wird die gesamte Datei als verwaltet behandelt und in `u-boot.yaml` referenziert.
- fehlt ein verwalteter Block in einer vorhandenen Datei:
  - ist `--backup` gesetzt, wird vor jedem vollständigen Überschreiben der komplette Dateiinhalt gesichert und danach ersetzt.
  - ist `--backup` nicht gesetzt, wird der Vorgang mit einem fachlichen Fehler (Code `10`) abgebrochen; es erfolgt ein klarer Hinweis auf die nötige Option `--backup`.
- bei vollständiger Überschreibung ohne verwalteten Block gilt ein vollständiges Backup vor dem Schreiben als Pflicht.

### LH-FA-DEV-003.a — Mechanik der Freigabe externer Feature-Quellen

Verfeinert [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features).

- Die Freigabe erfolgt als klarer, protokollierter Schritt im interaktiven Modus oder im Skriptmodus nur über die explizite Option:
  - `--allow-external-feature-sources <quelle>[,<quelle>...]` (`interaktiv`: Quelle bei Nachfrage bestätigen, `nicht-interaktiv`: alle Quellen als Flag-Argumente übergeben).
  - Die Option ist nur für diese Befehle gültig:
    - `u-boot init --devcontainer`
    - `u-boot generate devcontainer`
    - `u-boot config set devcontainer.featureSources.allow`
- Ein einzelnes `--allow-external-feature-sources` kann mehrere explizit erlaubte Quellen über Komma trennen.

### LH-FA-DEV-006.a — Workspace-Volume und Clone-Quelle im Sandbox-Profil

Verfeinert [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil).

- das Host-Arbeitsverzeichnis nicht per Bind-Mount einbinden; der Workspace liegt in einem benannten Volume, das Repository wird im Container geklont (Quelle: `devcontainer.sandbox.repository` (URL, optional), sonst die URL des Remotes `origin` des Projekt-Repositories; so kann ein anderes Repository als das Projekt-Repository geklont werden, aus dem im Container gepullt und gepusht wird; enthält die URL Zugangsdaten oder unzulässige Zeichen, ein fachlicher Fehler, Exit-Code `10`; ohne Quelle wird kein Clone-Schritt erzeugt und die Befehlsausgabe weist mit einer Warnung darauf hin, `u-boot generate devcontainer` ergänzt den Schritt, sobald eine Quelle existiert);

### LH-FA-DEV-007.a — Degradation und Strenge der Sandbox-Fähigkeiten

Verfeinert [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer).

**Degradation und Strenge** (gilt für dieses und das folgende Egress-Feature): `devcontainer.sandbox.onUnavailable` (`warn` | `fail`, Default `warn`).

| Zustand | `warn` (Default) | `fail` |
| ------- | ---------------- | ------ |
| `/dev/fuse` nicht verfügbar | Fallback auf `vfs`-Storage, Warnung mit Hinweis | Umgebungsproblem, Exit-Code `11` |
| Nested User-Namespaces durch Seccomp/AppArmor blockiert und `nestedRuntime: podman` | Umgebungsproblem, Exit-Code `11` | Umgebungsproblem, Exit-Code `11` |
| Egress-Capability (`NET_ADMIN`) nicht gewährbar ([`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion)) | Egress-Restriktion entfällt, Warnung mit Hinweis auf Restriktion auf Netz-/DNS-Ebene | Umgebungsproblem, Exit-Code `11` |

Ausdrücklich angeforderte Runtime (`nestedRuntime: podman`) wird nie stillschweigend durch etwas anderes ersetzt. Jeder Fallback wird in der Befehlsausgabe und in `u-boot doctor` ausgewiesen. Die Prüfung erfolgt durch `u-boot doctor` (soweit vom Host aus ermittelbar) und beim Containerstart durch das erzeugte Startscript, das im Fehlerfall nicht-null endet. Ungültige Werte der Schlüssel führen zu einem fachlichen Validierungsfehler (Exit-Code `10`).

### LH-FA-ADD-005.a — Zustandsregeln für registrierte und aktive Services

Verfeinert [`LH-FA-ADD-005`](lastenheft.md#lh-fa-add-005--mehrfaches-hinzufügen-verhindern).

- Ein bereits vorhandener Service darf nicht doppelt eingefügt werden.
- Ein Service gilt als registriert, sobald `services.<name>` in `u-boot.yaml` existiert.
- Er gilt als aktiv vorhanden, wenn `services.<name>.enabled` explizit auf `true` steht **und** ein verwalteter Eintrag in `compose.yaml` existiert.
- `services.<name>.enabled` ist immer explizit zu setzen. Ein registrierter Service ohne expliziten `enabled`-Schlüssel gilt als deaktiviert (`false`) und führt bei `u-boot doctor` zu einer `warn`-Diagnose, die das explizite Setzen empfiehlt.
- Liegt `services.<name>.enabled: false` vor, gilt der Service als deaktiviert (weiterhin registriert), und `u-boot add <service>` darf ihn idempotent reaktivieren.
- Besteht `services.<name>` nicht in `u-boot.yaml`, aber ein verwalteter Block in `compose.yaml`, darf die Inkonsistenz nicht stillschweigend ignoriert werden. Der Befehl muss mit klarer Diagnose abbrechen und auf manuelle Bereinigung oder Re-Konfiguration verweisen.
- Besteht `services.<name>` in `u-boot.yaml` mit `enabled: true`, aber der verwaltete Compose-Eintrag fehlt, muss das Verhalten deterministisch sein: `u-boot add <service>` erzeugt den fehlenden Compose-Block wieder.

### LH-FA-ADD-006.a — Verhalten bei erkannter Add-on-Abhängigkeit

Verfeinert [`LH-FA-ADD-006`](lastenheft.md#lh-fa-add-006--add-on-abhängigkeiten).

Verhalten bei erkannter abhängiger Konfiguration:

- Ist `services.keycloak.persistence: external-postgres` gesetzt und PostgreSQL nicht vorhanden, darf der Aufruf nicht stillschweigend fortfahren.
- Ist die optionale Abhängigkeit nicht aktiv, darf Keycloak ohne PostgreSQL angelegt werden.

- Im interaktiven Modus (Standardmodus) muss das Produkt nachfragen, ob das fehlende Add-on automatisch hinzugefügt werden soll.
- Im nicht-interaktiven Modus (`--no-interactive`) ohne `--with-deps` muss das Produkt mit Exit-Code `10` abbrechen und auf die fehlende Abhängigkeit hinweisen.
- Über die Option `--with-deps` muss das Produkt fehlende Abhängigkeiten automatisch hinzufügen. `--with-deps` ist mit `--no-interactive` kombinierbar; in dem Fall werden Abhängigkeiten deterministisch und ohne Rückfrage installiert.
- Mit `--yes` (ohne `--with-deps`) wird die Standardentscheidung "Abhängigkeit hinzufügen" deterministisch ausgeführt, ohne dass eine Rückfrage gestellt wird.
- Mit `--yes` oder `--no-interactive` (jeweils exklusiv) muss das Verhalten in Skript-/CI-Umgebungen deterministisch und nicht-blockierend sein.

### LH-FA-DIAG-002.b — Schweregrade der Devcontainer-Prüfungen

Verfeinert [`LH-FA-DIAG-002`](lastenheft.md#lh-fa-diag-002--lokale-voraussetzungen-prüfen).

- falls Devcontainer-Dateien vorhanden sind:
  - Ist `u-boot.yaml` vorhanden und `devcontainer.enabled == true`, müssen diese Prüfungen mit `error` bewertet werden:
    - syntaktische Gültigkeit von `.devcontainer/devcontainer.json`
    - Mindestkompatibilität mit VS Code Dev Containers (`name` gesetzt; mindestens `image` oder `build` vorhanden)
    - `forwardPorts`-Konsistenz zu aktivierten Services, falls Portangaben existieren
  - Ist `u-boot.yaml` vorhanden und `devcontainer.enabled == false`, sind die obigen Prüfungen optional (`warn`, keine harte Validierungspflicht).
  - Ist keine `u-boot.yaml` vorhanden, werden die obigen Prüfungen als ergänzende Qualitätsdiagnosen mit `warn` ausgegeben.

### LH-FA-DIAG-002.c — Konsistenzregeln für `forwardPorts`

Verfeinert [`LH-FA-DIAG-002`](lastenheft.md#lh-fa-diag-002--lokale-voraussetzungen-prüfen).

- `forwardPorts`-Konsistenzregeln:
  - Für jeden aktivierten Service mit expliziter `ports`-Zuordnung (TCP) ist der Host-Port in `forwardPorts` enthalten.
  - bei mehreren TCP-Ports werden eindeutige Portzahlen eingetragen (Duplikate dedupliziert).
  - UDP- oder nicht eindeutig auflösbare Portangaben dürfen in `forwardPorts` fehlen; dafür ist ein `warn`-Diagnoseeintrag zulässig.

### LH-FA-PROJDOCS-002.a — Bestandsregel und Abgleich zum ADR-Format

Verfeinert [`LH-FA-PROJDOCS-002`](lastenheft.md#lh-fa-projdocs-002--adr-format).

**Grandfathering (Bestand).** Die zum Zeitpunkt der Format-Umstellung (Regelwerk-v3.5.1-Adoption) bereits `Accepted` ADRs (`0001`–`0010`, `0013`) bleiben in der vorherigen leanen Form (`## Status`/`## Datum` als Überschriften; Abschnitte Kontext/Entscheidung/Konsequenzen) und sind als `Accepted` **unveränderlich**; sie werden **nicht** migriert. Das MADR-Format gilt für alle **neu** angelegten ADRs sowie für noch mutable `Proposed`-ADRs (`0011`, `0012`) beim nächsten inhaltlichen Anfassen.

**Reconciliation.** Titel-Form (`# ADR <Nr>: <Titel>`) und Superseded-Referenz (`<NNNN>-<slug>`, klickbar) folgen der bestehenden u-boot-Konvention, nicht dem Template-Wortlaut (`# ADR-NNNN:` bzw. `Superseded by ADR-NNNN`); die MADR-**Substanz** (Inline-Kopf-Felder, `Schärft`-Aufwärtskopplung, Alternativen/Fitness-Function/Re-Eval/Geschichte) wird übernommen. Die Umstellung ist ein Change Request am Vertrags-Stratum (Trigger: v3.5.1-Adoption).

### LH-FA-PROJDOCS-005.a — Pflichten der Carveout-Disziplin

Verfeinert [`LH-FA-PROJDOCS-005`](lastenheft.md#lh-fa-projdocs-005--carveout-disziplin).

Pflichten:

- Der Slice-Plan folgt der Dateiname-Konvention aus [`LH-FA-PROJDOCS-003`](lastenheft.md#lh-fa-projdocs-003--planning-lifecycle) (`slice-<phase>-<slug>.md`).
- Der Plan benennt mindestens: Auslöser (was wurde wo bewusst weggelassen), Aufhebungsbedingung (was muss passieren), Akzeptanzkriterien.
- Wo der Carveout in einer Spec-Anforderung dokumentiert ist (z. B. [`LH-FA-BUILD-008`](lastenheft.md#lh-fa-build-008--coverage-bootstrap) für Coverage-Bootstrap), bleibt die Spec-Anforderung die normative Quelle; der Plan-Verweis lebt im Carveout-Inventar und in der Roadmap.
- **Doppelte Verankerung:** jeder temporäre Carveout ist sowohl im Carveout-Inventar als auch in der Roadmap als Slice-Zeile sichtbar. Carveouts ohne Roadmap-Eintrag oder Slice-Pläne ohne Carveout-Inventar-Verweis sind Verstoß gegen diese Anforderung.
- Auch Spec-Open-Punkte (`LH-OPEN-*`) und ADR-Folgepunkte gelten als temporäre Carveouts und brauchen einen Slice-Plan — kein „bleibt offen bis MVP-Closure" als Inventar-Eintrag.
- Ein Master-Inventar in `carveouts.md` listet alle aktuellen Carveouts mit Status (`temporär` + Plan-Verweis vs. `permanent` + Begründung). Diese Datei lebt analog zur `roadmap.md` dauerhaft in `in-progress/`.

### LH-FA-PROJDOCS-006.a — Pflichten des Dokumentationsreferenzmodells

Verfeinert [`LH-FA-PROJDOCS-006`](lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell).

Pflichten:

- Das Lastenheft darf nur `LH-*`-Anforderungen intern normativ
  querverweisen. ADRs, Slices, Carveouts und Roadmap/Wellen dürfen im
  Lastenheft keine Quelle der Normativität sein.
- Technische oder Sicht-Specs dürfen auf das Lastenheft und innerhalb
  ihres Stratums referenzieren, aber keine ADRs, Slices, Carveouts oder
  Roadmap/Wellen als bindenden Text verlinken.
- ADRs dürfen `LH-*`, betroffene Spec-Stellen und aktive ADRs normativ
  referenzieren. Superseded ADRs dürfen nur innerhalb der ADR-Lineage als
  Historie referenziert werden.
- Slices dürfen `LH-*` und aktive ADRs normativ referenzieren.
  Slice-zu-Slice-, Slice-zu-Carveout- und Slice-zu-Roadmap-Referenzen
  sind ausschließlich Kontext.
- Carveouts dürfen `LH-*` und aktive ADRs normativ referenzieren.
  Carveout-zu-Slice- und Slice-zu-Carveout-Referenzen sind nur Owner-,
  Trigger- oder Closure-Buchführung.
- Roadmap/Wellen orchestrieren Arbeit, tragen aber keine normative
  Ableitungskraft.
- `docs-check` muss Markdown-Link-Pfade, Heading-Anker, verlinkte
  `ADR-*`-Kennungen und die Referenzmatrix für Lastenheft-, ADR-,
  Spec-, Slice-, Carveout- und Roadmap/Wellen-Artefakte in `docs/`,
  `spec/`, `harness/` und Root-Markdown prüfen.
- Die Kennungs-Linkpflicht wird stufenweise aktiviert: `ADR-*` gilt
  global, `LH-*` gilt in Spec-Straten außerhalb des Lastenhefts, in
  README-Dateien und in `docs/user/`. Eindeutig auflösbare Slice- und
  Tranche-IDs gelten in README-Dateien und in `docs/user/`; Markdown-
  Überschriften sind ausgenommen, damit bestehende Section-Anker stabil
  bleiben. Konkrete `PH-*`- und `TC-*`-Kennungen in der Traceability-
  Matrix sind bis zu getrennten Pflichtenheft-/Testfall-Artefakten
  Traceability-Aliase und verlinken auf die zugehörige `LH-*`-
  Anforderung derselben Matrixzeile. Konkrete `CO-*`-Kennungen sind

### LH-NFA-USE-004.a — Regeln für `--json`-Antworten

Verfeinert [`LH-NFA-USE-004`](lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe).

Für `--json`-Antworten gilt zusätzlich:

- `diagnostics`, wenn leer, darf als `[]` ausgegeben werden.
- `diagnostics.level` darf nur `warn` oder `error` enthalten.
- `diagnostics.code` folgt der Konvention: LH-Kennung der verursachenden Anforderung (z. B. [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features)). Tool-interne Codes ohne LH-Bezug dürfen nur dann verwendet werden, wenn ihre Bedeutung in der Dokumentation festgehalten ist (Verweis: [`LH-FA-CLI-007`](lastenheft.md#lh-fa-cli-007--dry-run)).
- `diagnostics.file` ist optional.
- `status` ist an den höchsten in `diagnostics` enthaltenen `level` gekoppelt: `error` → `status == "error"`; `warn` ohne `error` → `status == "warn"`; sonst `status == "ok"`.
- Bei `command == "template"` oder `command == "config"` ist `subcommand` verpflichtend.
- Die Felder `status`, `command`, `diagnostics` und `exitCode` sind minimal verpflichtend und sollten mit anderen Feldern ergänzt werden.

## 2. Datenstrukturen und Schemas

Formate und Schemata (`u-boot.yaml`, JSON-Ausgabe, CLI-Tabellen). Jede Struktur trägt
eine `SPEC-<NNN>`.

### SPEC-001 — JSON-Schema der Vorschau-Ausgabe (`--dry-run --json`)

Gilt für `--dry-run --json` und `--diff --json`; die Pflichtfelder nennt die Anforderung.

Für `--dry-run --json` ist die Ausgabe mindestens wie folgt als maschinenlesbares JSON zu liefern:

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["status", "command", "dryRun", "diff", "plannedFiles", "changes", "diagnostics", "exitCode"],
  "properties": {
    "subcommand": {
      "type": "string",
      "description": "Unterkommando bei gruppierten Hauptkommandos wie `template` oder `config`"
    },
    "status": {
      "type": "string",
      "enum": ["ok", "warn", "error"]
    },
    "command": {
      "type": "string",
      "enum": ["init", "add", "remove", "up", "down", "doctor", "logs", "generate", "config", "template"]
    },
    "dryRun": {
      "type": "boolean"
    },
    "diff": {
      "type": "boolean"
    },
    "plannedFiles": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["path", "action"],
        "properties": {
          "path": { "type": "string" },
          "action": {
            "type": "string",
            "enum": ["create", "modify", "delete"]
          }
        },
        "additionalProperties": true
      }
    },
    "changes": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["path", "count"],
        "properties": {
          "path": { "type": "string" },
          "count": { "type": "integer", "minimum": 0 }
        },
        "additionalProperties": true
      }
    },
    "diagnostics": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["level", "code", "message"],
        "properties": {
          "level": { "type": "string", "enum": ["warn", "error"] },
          "code": { "type": "string" },
          "message": { "type": "string" },
          "file": { "type": "string" }
        },
        "additionalProperties": true
      }
    },
    "exitCode": {
      "type": "integer",
      "minimum": 0
    }
  },
  "allOf": [
    {
      "if": {
        "properties": {
          "command": { "const": "template" }
        },
        "required": ["command"]
      },
      "then": {
        "required": ["subcommand"]
      }
    },
    {
      "if": {
        "properties": {
          "command": { "const": "config" }
        },
        "required": ["command"]
      },
      "then": {
        "required": ["subcommand"]
      }
    }
  ],
  "additionalProperties": true
}
```

### SPEC-002 — Beispielinstanz einer Vorschau-Ausgabe (`add --dry-run --json`)

Beispielinstanz:

```json
{
  "status": "warn",
  "command": "add",
  "dryRun": true,
  "diff": false,
  "plannedFiles": [
    { "path": "compose.yaml", "action": "create" }
  ],
  "changes": [
    { "path": "compose.yaml", "count": 12 }
  ],
  "diagnostics": [
    { "level": "warn", "code": "LH-FA-CLI-007", "message": "Geplante Datei fehlt bereits" }
  ],
  "exitCode": 0
}
```

### SPEC-003 — Beispielinstanz einer Diff-Ausgabe (`add --diff --json`, ohne `--dry-run`)

Beispiel für `--diff --json` ohne `--dry-run` (Vorschau mit anschließendem Schreiben):

```json
{
  "status": "ok",
  "command": "add",
  "dryRun": false,
  "diff": true,
  "plannedFiles": [
    { "path": "compose.yaml", "action": "modify" }
  ],
  "changes": [
    { "path": "compose.yaml", "count": 6 }
  ],
  "diagnostics": [],
  "exitCode": 0
}
```

### SPEC-004 — Projektnamen-Muster

Gilt für den Projektnamen (`project.name`) bei `init` und `config set`.

- regulärer Ausdruck: `^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$`

### SPEC-005 — Schema der Projektkonfiguration `u-boot.yaml`

Gilt als Mindestinhalt der Konfigurationsdatei; die Hinweise zu Schlüsseln und Werten nennt die Anforderung.

```yaml
schemaVersion: 1
project:
  name: my-service

services:
  postgres:
    enabled: false

devcontainer:
  enabled: false

# Optionale, V1-relevante Felder:
# services:
#   keycloak:
#     enabled: false
#     persistence: embedded   # embedded | external-postgres
#   otel:
#     enabled: false
#
# devcontainer:
#   featureSources:
#     allow:
#       - https://ghcr.io/devcontainers/features/node
#   user:
#     uid: 1000                    # 1..65535
#   profile: default               # default | sandbox
#   sandbox:
#     nestedRuntime: none          # none | podman
#     onUnavailable: warn          # warn | fail
#     repository: ''               # Clone-Quelle statt origin (URL)
#     egress:
#       enabled: false
#       allow: []
```

### SPEC-008 — Eigenschaften und Pflicht-Targets des Makefiles

Pflicht-Eigenschaften:

- `.DEFAULT_GOAL := help`
- `.PHONY` für alle Targets gesetzt
- `help`-Target mit Übersicht über alle verfügbaren Targets
- Variablen mit `?=`-Defaults für Overridability (`IMAGE`, `GO_VERSION`, `GOLANGCI_LINT_VERSION`, `THRESHOLD`)

MVP-Pflicht-Targets:

| Target          | Zweck                                                           |
| --------------- | --------------------------------------------------------------- |
| `help`          | Übersicht aller Targets                                         |
| `deps`          | `docker build --target deps`                                    |
| `compile`       | `docker build --target compile`                                 |
| `lint`          | `docker build --target lint`                                    |
| `test`          | `docker build --target test`                                    |
| `coverage`      | Alias auf `coverage-gate`                                       |
| `coverage-gate` | `docker build --target coverage --build-arg COVERAGE_THRESHOLD` |
| `build`         | `docker build --target runtime`                                 |
| `run`           | `docker run --rm <image> --help` (Smoketest); Dependency: `build` |
| `clean`         | lokale Artefakte und gebaute Images entfernen                   |

### SPEC-009 — Mindestlayout des u-boot-Repositories

Mindestlayout:

```text
.
├── cmd/
│   └── uboot/
│       └── main.go              # Entry point der CLI (Wiring-Schicht)
├── internal/                    # nicht-exportierbare Implementierung
│   ├── hexagon/                 # innere Schichten (LH-FA-ARCH-002)
│   │   ├── domain/
│   │   ├── application/
│   │   └── port/{driving,driven}/
│   └── adapter/                 # äußere Schichten (LH-FA-ARCH-002)
│       ├── driving/
│       └── driven/
├── spec/                        # Lastenheft, weitere Spezifikationen
│   ├── lastenheft.md
│   └── <weitere-spezifikation>.md
├── docs/                        # Doku-Struktur (LH-FA-PROJDOCS-001)
├── go.mod
├── go.sum
├── Dockerfile
├── Makefile
├── .dockerignore
├── .gitignore
├── LICENSE
└── README.md
```

### SPEC-010 — Mindest-Verzeichnisstruktur unter `docs/`

```text
docs/
├── archive/                  # abgelöste oder ersetzte Inhalte
├── plan/
│   ├── adr/                  # Architecture Decision Records
│   └── planning/
│       ├── open/             # Backlog
│       ├── next/             # priorisiert für nächsten Schritt
│       ├── in-progress/      # aktiv bearbeitet
│       └── done/             # abgeschlossen
└── user/                     # User-facing Dokumentation
```

### SPEC-011 — Format der Architecture Decision Records

- Dateiname beginnt mit vierstelliger Nummer, beginnend bei `0001` und monoton steigend: `0001-<slug>.md`, `0002-<slug>.md`; Slug in Kebab-Case (z. B. `0001-implementierungssprache-go.md`).
- Dokumenttitel als `#`-Überschrift: `# ADR <Nr>: <Titel>`.
- Direkt darunter die Kopf-Felder als **fette Inline-Felder** (nicht als `##`-Überschriften):
  - `**Status:**` – einer aus `Proposed`, `Accepted`, `Deprecated`, `Superseded by <NNNN>-<slug>`.
  - `**Datum:**` – Entscheidungsdatum im Format `YYYY-MM-DD`.
  - `**Autor:**` – verantwortliche Rolle oder Person.
  - `**Bezug:**` – betroffene `LH-*`- und ggf. Vorgänger-ADR-IDs als Markdown-Links (optional, wenn zutreffend).
  - `**Schärft:**` – welche Spec-Stelle (Abschnitt `§N` der technischen oder der Sicht-Spezifikation) diese ADR verbindlich macht, als Aufwärts-Deklaration der Änderungskopplung (wer die ADR ändert, zieht von hier die Spec-Stellen nach); `—`, wenn Prozess-ADR ohne Spec-Bezug.
- Danach die Abschnitte, jeweils als `##`-Überschrift, in dieser Reihenfolge:
  1. `## Kontext` – Ausgangslage, auslösende Anforderung, tragende Annahmen.
  2. `## Entscheidung` – die Wahl, eindeutig.
  3. `## Verglichene Alternativen` – Optionen mit Pro/Contra (auch „nichts tun").
  4. `## Konsequenzen` – kurz- und langfristige Folgen, positiv und negativ, inkl. Folgepflichten.
  5. `## Fitness Function` – die maschinell prüfbare Regel, falls die Entscheidung sich in einer Code-Eigenschaft niederschlägt (sonst entfällt der Abschnitt).
  6. `## Re-Evaluierungs-Trigger` – wann die Entscheidung erneut zu prüfen ist.
  7. `## Geschichte` – Tabelle Datum/Ereignis/Verweis (`Proposed`, `Accepted`, …).
- ADR-Nummern werden nie wiederverwendet; abgelöste ADRs bleiben mit Status `Superseded by <NNNN>-<slug>` erhalten und verweisen auf den Nachfolger über den vollen Dateinamen-Stamm (ohne `.md`), als klickbaren Link.

### SPEC-012 — Schichten und Verzeichnislayout unter `internal/`

```text
internal/
├── hexagon/
│   ├── domain/         # reine Datentypen + invariantes Verhalten, keine I/O
│   ├── application/    # Use-Cases; ruft ausschließlich Ports auf
│   └── port/
│       ├── driving/    # Interfaces, die von außen (CLI/HTTP) konsumiert werden
│       └── driven/     # Interfaces, die das Application nach außen ruft
└── adapter/
    ├── driving/        # konkrete Driver (z. B. cli/ mit Cobra-Commands)
    └── driven/         # konkrete Adapter (z. B. docker/, fs/, yaml/)
```

### SPEC-013 — Import-Regel-Tabelle der Schichten

| Schicht | darf importieren | darf nicht importieren |
| --- | --- | --- |
| `hexagon/domain` | Go-Standard-Library | alle anderen `internal/`-Pakete, I/O-Libraries |
| `hexagon/application` | `hexagon/domain`, `hexagon/port/driving`, `hexagon/port/driven` | `adapter/*`, externe I/O-Libraries |
| `hexagon/port/driving` | `hexagon/domain` | `hexagon/application`, `hexagon/port/driven`, `adapter/*` |
| `hexagon/port/driven` | `hexagon/domain` | `hexagon/application`, `hexagon/port/driving`, `adapter/*` |
| `adapter/driving` | `hexagon/domain`, `hexagon/port/driving`, externe Libraries | `hexagon/application`, `adapter/driven` |
| `adapter/driven` | `hexagon/domain`, `hexagon/port/driven`, externe Libraries | `hexagon/application`, `adapter/driving` |
| `cmd/uboot` | `internal/...`, Standardbibliothek, externe Libraries | keine Einschränkung; Wiring-Schicht |

### SPEC-014 — Markierungsformat verwalteter Bereiche je Dateityp

- YAML, `.env`, `Dockerfile`, Shell-Skripte (`#`-Kommentare):

  ```yaml
  # BEGIN U-BOOT MANAGED BLOCK: postgres
  # ...
  # END U-BOOT MANAGED BLOCK: postgres
  ```

- Markdown (`README.md`, `CHANGELOG.md`) als HTML-Kommentar:

  ```markdown
  <!-- BEGIN U-BOOT MANAGED BLOCK: postgres -->
  ...
  <!-- END U-BOOT MANAGED BLOCK: postgres -->
  ```

- JSONC (z. B. `.devcontainer/devcontainer.json`):

  ```jsonc
  // BEGIN U-BOOT MANAGED BLOCK: postgres
  // ...
  // END U-BOOT MANAGED BLOCK: postgres
  ```

### SPEC-015 — Beispielinstanz einer Minimalkontrakt-Ausgabe (`doctor --json`)

Beispiel:

```json
{
  "status": "ok",
  "command": "doctor",
  "diagnostics": [],
  "exitCode": 0
}
```

## 3. Defaults und Konstanten

Werte, die im Produkt fest sind (Standardwerte, Grenzwerte, Versionsuntergrenzen).

### SPEC-006 — Build-Args des Dockerfiles und Pin-Politik

Gilt für das Dockerfile im Repo-Root; Overrides per `docker build --build-arg <NAME>=<value>` bzw. Makefile-Variable.

- `ARG GO_VERSION` – mit Default-Pin (Stand 2026-10-02: `1.27.1`); Hebung ist Routine ohne separaten Spec-Eintrag.
- `ARG GOLANGCI_LINT_VERSION` – mit Default-Pin; gleiche Pin-Politik.
- `ARG COVERAGE_THRESHOLD` – mit Default `90` (Prozent) und Override-Pfad `make coverage-gate THRESHOLD=…`; das Bootstrap-Verhalten bei leerer Coverage-Eingabe steht in der Anforderung.

### SPEC-007 — Mindest-Ausschlüsse der `.dockerignore`

Mindestens auszuschließen:

- `.git`, `.gitignore`, `.github`
- IDE-Verzeichnisse: `.idea`, `.vscode`
- Agent-Verzeichnisse: `.claude`, `.codex`, `.agents`
- lokale Build-Artefakte und Caches (z. B. `dist/`, `coverage*`, `*.log`)

### SPEC-016 — Go-Toolchain: Mindestversion und Dockerfile-Pin

Mindest-Toolchain: Go 1.26 oder neuer (`go 1.26.0` in `go.mod`, analog Referenzprojekt `k-deskflight`); Default-Pin im Dockerfile als `ARG GO_VERSION` (am Entscheidungsdatum `1.26.3`, Stand 2026-10-02 `1.27.1`). Pin-Hebung ist Routine ohne separaten Spec-Eintrag.

### SPEC-018 — Standardwerte der Build- und Laufzeitumgebung

| Name | Wert | Quelle |
|---|---|---|
| Coverage-Schwellwert (`COVERAGE_THRESHOLD` / `THRESHOLD`) | 90 Prozent | Dockerfile-Build-Arg, Makefile-Variable |
| Polling-Intervall von `u-boot up` | 500 ms | Stabilisierungsschleife (`docker compose ps`) |
| Standard-Wartezeit von `u-boot up` | 60 s | `--timeout` |
| UID des Container-Benutzers | 1000 | `devcontainer.user.uid` |

### SPEC-019 — Default-Allowlist der Egress-Restriktion

Die wirksame Liste ist die Vereinigung aus Basis, Ergänzungen und Nutzerliste, sortiert und ohne Duplikate.

| Quelle | Hosts |
|---|---|
| immer | `github.com`, `api.github.com`, `codeload.github.com`, `raw.githubusercontent.com`, `objects.githubusercontent.com`, `deb.debian.org`, `security.debian.org` |
| `nestedRuntime: podman` | `registry-1.docker.io`, `auth.docker.io`, `production.cloudflare.docker.com`, `ghcr.io` |
| Feature `node` | `registry.npmjs.org` |
| Feature `go` | `proxy.golang.org`, `sum.golang.org`, `storage.googleapis.com` |
| Feature `java` | `repo.maven.apache.org`, `repo1.maven.org` |
| Clone-Quelle | Host von `devcontainer.sandbox.repository` beziehungsweise `origin` |
| Nutzer | Einträge von `devcontainer.sandbox.egress.allow` (kleingeschriebene Hostnamen, ohne Schema, Port und Wildcard) |

## 4. Fehler-Codes und Logging-Felder

Verbindliche Diagnose- und Fehler-Codes sowie Logging-Felder.

### SPEC-017 — Diagnose-Codes der `doctor`-Prüfungen (Code-Registry)

Jede Prüfung von `u-boot doctor` trägt einen stabilen Code in `diagnostics[].code`. Die Tabelle ist die kanonische Registry; der Quelltext (`DefaultAllowedCodes`) und diese Tabelle werden von einem Test symmetrisch abgeglichen. Andere Befehle tragen als Code die Kennung der verursachenden Anforderung.

<!-- code-registry:start -->

| Code | Bedeutung |
| --- | --- |
| `fs.write-permissions` | Schreib-Permission im Working Directory |
| `git.installed` | Git-Binary verfügbar |
| `docker.installed` | Docker-Binary verfügbar |
| `docker.reachable` | Docker-Daemon erreichbar |
| `docker.compose.installed` | Compose-Plugin verfügbar |
| `uboot.yaml.valid` | `u-boot.yaml` syntaktisch valide |
| `compose.yaml.valid` | `compose.yaml` syntaktisch valide |
| `devcontainer.json.valid` | `.devcontainer/devcontainer.json` syntaktisch valide |
| `devcontainer.dockerfile.valid` | `.devcontainer/Dockerfile` parsebar |
| `services.enabled-key` | `u-boot.yaml` `services`-Block konsistent |
| `devcontainer.forwardPorts.consistency` | `devcontainer.json` `forwardPorts` konsistent |
| `devcontainer.features.allowlist` | `devcontainer` Features auf Allowlist |
| `devcontainer.features.drift` | `devcontainer` Features ohne Drift |
| `devcontainer.sandbox.runtime` | Sandbox: nested Podman (`/dev/fuse`, Profil-Konsistenz), [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) |
| `devcontainer.sandbox.egress` | Sandbox: Egress-Restriktion konsistent konfiguriert, [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) |
| `devcontainer.sandbox.credentials` | Sandbox: keine Klartext-Git-Zugangsdaten, Token-Quelle, [`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) |

<!-- code-registry:end -->

## 5. Metriken und Tracing-Felder

Verbindliche Telemetrie-Felder pro Span.

### SPEC-020 — Telemetrie des Produkts

`u-boot` erzeugt selbst keine Metriken oder Traces und hat keine Tracing-Felder. Der Add-on `otel` richtet lediglich einen OpenTelemetry-Collector für das Zielprojekt ein (siehe Add-on-Katalog); dessen Standardports stehen dort.

## 6. Externe Verträge

Schnittstellen zu Drittsystemen mit Versionsannahme (Docker, Compose, Devcontainer).

### SPEC-021 — Docker, Docker Compose und Podman

| System | Mindestversion | Hinweis |
|---|---|---|
| Docker Engine | 24.0.0 | Prüfung durch `u-boot doctor` |
| Docker Compose (Plugin) | 2.20.0 | Prüfung durch `u-boot doctor` |
| Podman (als Docker-Ersatz) | 4.0 | mit aktivem `podman.socket` und `DOCKER_HOST` auf den Socket; nicht erkannte Versionen werden als Warnung gemeldet |

### SPEC-022 — Add-on-Katalog: Images und Ports

| Add-on | Image | Host-Ports |
|---|---|---|
| `postgres` | `postgres:16-alpine` | 5432 |
| `keycloak` | `quay.io/keycloak/keycloak:26.0` | 8080 |
| `otel` | `otel/opentelemetry-collector:0.108.0` | 4317 (OTLP/gRPC), 4318 (OTLP/HTTP) |

### SPEC-023 — Devcontainer: Basis-Image und Feature-Quellen

- Basis-Image des erzeugten Dockerfiles: `mcr.microsoft.com/devcontainers/base:debian`.
- Eingebaute Features (Schlüssel → Quelle unter `ghcr.io/devcontainers/features/`, jeweils mit der Standardversion `1`): `git` → `git`, `docker-cli` → `docker-outside-of-docker`, `node` → `node`, `java` → `java`, `go` → `go`, `cpp` → `cpp`, `kubectl-helm` → `kubectl-helm-minikube`, `postgres-client` → `postgresql-client`. Jede andere Quelle gilt als extern und braucht eine Freigabe in `devcontainer.featureSources.allow`.
- Externe Feature-Quellen müssen als URL mit Schema `http://`, `https://` oder `oci://` angegeben werden.

## 7. Historie

| Datum | Änderung |
|---|---|
| 2026-10-02 | Angelegt als Technik-Stratum; Schemata, Beispielinstanzen, Algorithmen, Build-/CI-, Doku- und Architektur-Details sowie Markierungsformate wörtlich aus dem Lastenheft übernommen (`SPEC-001`..`SPEC-016`, Verfeinerungen); neu aus dem Bestand von Code und Konfiguration: Standardwerte, Egress-Default-Allowlist, Code-Registry der `doctor`-Prüfungen, Telemetrie-Aussage, externe Verträge (`SPEC-017`..`SPEC-023`). |
