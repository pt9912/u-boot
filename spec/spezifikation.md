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
Verträge) und jede Festlegung, die zwei oder mehr Anforderungen zugleich verfeinert, trägt `SPEC-<NNN>`: dreistellig, fortlaufend je Datei, nicht je Abschnitt.
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

### SPEC-025 — Stages des Multi-Stage-Dockerfiles

<a id="lh-fa-build-001a--stages-des-multi-stage-dockerfiles"></a>

- BuildKit-Direktive in der ersten Zeile: `# syntax=docker/dockerfile:1.7`.
- Pflicht-Stages im MVP:
  - `deps` – Modulauflösung (`go mod download`) als Cache-Layer.
  - `compile` – schnelles Compile-Feedback (`go build`) ohne Tests/Lint.
  - `test` – `go test ./...`.
  - `lint` – `golangci-lint run ./...`.
  - `coverage` – `go test -coverprofile` + Coverage-Gate gegen `COVERAGE_THRESHOLD`; Bootstrap-Verhalten und Scope sind in [`SPEC-032`](#spec-032--bootstrap-verhalten-der-coverage-stage) und [`SPEC-033`](#spec-033--go-layout-des-repositories) definiert.
  - `build` – statisch gelinktes Binary (`CGO_ENABLED=0`, `-ldflags="-s -w"`).
  - `runtime` – minimales Endimage ([`LH-NFA-PORT-004`](lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image)).
- Jeder Stage ist ein eigenständiges Build-Ziel und wird per `docker build --target <stage>` einzeln baubar.

### LH-NFA-PORT-004.a — Eigenschaften der Runtime-Stage

<a id="lh-fa-build-002a--eigenschaften-der-runtime-stage"></a>

Verfeinert [`LH-NFA-PORT-004`](lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image).

- Base-Image: `gcr.io/distroless/static-debian12:nonroot` (oder gleichwertig minimal und ohne Shell).
- Non-root-Benutzer (Distroless-`nonroot`-User, `USER 65532:65532`).
- `ENTRYPOINT` zeigt auf das im `build`-Stage erzeugte Binary; der konkrete Pfad (Empfehlung: `/usr/local/bin/u-boot`) ist im Dockerfile dokumentiert.
- OCI Image Labels gesetzt:
  - `org.opencontainers.image.source`
  - `org.opencontainers.image.description`
  - `org.opencontainers.image.licenses`
  - `org.opencontainers.image.title`
- Keine Build-Toolchain im Endimage; alle Build-Artefakte werden aus dem `build`-Stage per `COPY --from=build` übernommen.

### SPEC-026 — Komposition der CI-Pipeline (GitHub Actions)

<a id="lh-qa-003a--komposition-der-ci-pipeline-github-actions"></a>

- Workflow-Datei: [`.github/workflows/ci.yml`](../.github/workflows/ci.yml).
- Trigger: `pull_request` und `push` auf den Branch `main`.
- Drei Jobs, parallel, alle PR-blockierend (Required-Status-Checks im GitHub-UI nach dem ersten grünen Lauf zu konfigurieren; die Required-Status-Check-Liste muss die tatsächlichen Workflow-`name:`-Felder verwenden, nicht die Kurz-Keys). Die drei delegierten Make-Targets sind in [`SPEC-008`](#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) definiert; die PR-Blocking-**Pflicht** für alle drei kommt aus diesem [`SPEC-026`](#spec-026--komposition-der-ci-pipeline-github-actions)-Eintrag. Damit ist `make govulncheck` und `make image-scan` über [`SPEC-026`](#spec-026--komposition-der-ci-pipeline-github-actions) verbindlich, die Make-Target-Definitionen stehen in [`SPEC-008`](#spec-008--eigenschaften-und-pflicht-targets-des-makefiles):
  - `gates (lint + test + coverage-gate)` — führt `make gates` aus (lint + test + coverage-gate, [`SPEC-008`](#spec-008--eigenschaften-und-pflicht-targets-des-makefiles)).
  - `security-gates (govulncheck)` — führt `make govulncheck` aus ([`SPEC-008`](#spec-008--eigenschaften-und-pflicht-targets-des-makefiles), MVP-Pflicht via diesem [`SPEC-026`](#spec-026--komposition-der-ci-pipeline-github-actions)).
  - `image-scan (trivy HIGH+CRITICAL)` — führt `make image-scan` aus (Trivy gegen das Runtime-Image, severity `HIGH,CRITICAL`, exit-code `1`; [`SPEC-008`](#spec-008--eigenschaften-und-pflicht-targets-des-makefiles), MVP-Pflicht via diesem [`SPEC-026`](#spec-026--komposition-der-ci-pipeline-github-actions)).
- Runner: `ubuntu-latest`. Keine Host-Go-Toolchain ([`SPEC-031`](#spec-031--regeln-des-docker-only-workflows)); der Runner braucht nur das vorinstallierte Docker + BuildKit.
- Actions sind **SHA-gepinnt** mit Tag-Kommentar (Supply-Chain-Härtung gegen Tag-Move). Pin-Hebung ist Routine; neuer Commit-SHA via `gh api repos/<owner>/<repo>/git/refs/tags/<tag>`.
- Top-Level `permissions: {}` (alle Tokens entzogen); jeder Job lockert auf das Minimum (Defense-in-Depth).
- Jeder Job mit `timeout-minutes` versehen (Empfehlung: 20).
- Die konkrete Workflow-Ausprägung muss die hier genannten
  Anforderungen und die Build-Target-Verträge erfüllen.

### SPEC-027 — Zusammensetzung des Lint-Profils

<a id="lh-qa-004a--zusammensetzung-des-lint-profils"></a>

Profil-Komposition (29 Linter insgesamt):

- 5 Default-Linter (`govet`, `errcheck`, `staticcheck`, `unused`, `ineffassign`).
- 24 SOLID-nahe Zusatz-Linter (Komplexitäts-, Funktionslänge-, Interface-, Kopplungs- und Boundary-Signale). **`depguard`** für die Schicht-Regeln aus [`SPEC-034`](#spec-034--enforcement-der-import-regeln) ist Teil dieser 24.

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
| `/dev/fuse` im Container nicht verfügbar (Container ohne `--device=/dev/fuse` gestartet) | Fallback auf `vfs`-Storage, Warnung mit Hinweis | Umgebungsproblem, Exit-Code `11` |
| Nested User-Namespaces durch Seccomp/AppArmor blockiert und `nestedRuntime: podman` | Umgebungsproblem, Exit-Code `11` | Umgebungsproblem, Exit-Code `11` |
| Egress-Capability (`NET_ADMIN`) nicht gewährbar ([`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion)) | Egress-Restriktion entfällt, Warnung mit Hinweis auf Restriktion auf Netz-/DNS-Ebene | Umgebungsproblem, Exit-Code `11` |
| Docker-Daemon wird beim Containerstart nicht betriebsbereit (`nestedRuntime: docker`, etwa Start ohne `--privileged`) | Container ohne Daemon, Warnung mit Hinweis; keine andere Runtime als Ersatz | Umgebungsproblem, Exit-Code `11` |

Ein fehlendes `/dev/fuse` oder `/dev/net/tun` auf dem Container-Host ist keine Degradationsstufe und steht nicht in der Tabelle: die Container-Engine lehnt den Start ab, das Startscript läuft nie. `u-boot doctor` prüft `/dev/net/tun` nicht; ein fehlendes `/dev/fuse` weist es auf einem Linux-Host, soweit der Host es bestimmen kann, als Fehler aus (Umgebungsproblem, Exit-Code `11`, unabhängig von `onUnavailable` und ohne Kennzeichen eines Fallbacks).

Ausdrücklich angeforderte Runtime (`nestedRuntime: podman` oder `docker`) wird nie stillschweigend durch etwas anderes ersetzt. Jeder Fallback wird in der Befehlsausgabe und in `u-boot doctor` ausgewiesen: Das erzeugte Startscript schreibt bei jedem Lauf eine Statusdatei in den Container (`vfs`, `daemon`, `egress`) und endet im Fehlerfall nicht-null; `u-boot doctor` liest die Datei vom Host über die Container-Engine und weist je Fallback einen Eintrag aus (`devcontainer.sandbox.fallback.vfs`, `….daemon`, `….egress`), `u-boot generate devcontainer` und `u-boot init` weisen nach dem Schreiben den Fallback eines bestehenden Containers aus. Format, Ermittlung und Grenzen: [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung). Der Host meldet einen Fallback nur, wenn ein Container zum Projektordner gefunden und seine Statusdatei gelesen wird; sonst ist der Zustand „nicht ermittelbar“, was kein Fallback ist und auch nicht ausschließt, dass einer vorliegt. Ungültige Werte der Schlüssel führen zu einem fachlichen Validierungsfehler (Exit-Code `10`).

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

### SPEC-028 — Bestandsregel und Abgleich zum ADR-Format

<a id="lh-fa-projdocs-002a--bestandsregel-und-abgleich-zum-adr-format"></a>

**Grandfathering (Bestand).** Die zum Zeitpunkt der Format-Umstellung (Regelwerk-v3.5.1-Adoption) bereits `Accepted` ADRs (`0001`–`0010`, `0013`) bleiben in der vorherigen leanen Form (`## Status`/`## Datum` als Überschriften; Abschnitte Kontext/Entscheidung/Konsequenzen) und sind als `Accepted` **unveränderlich**; sie werden **nicht** migriert. Das MADR-Format gilt für alle **neu** angelegten ADRs sowie für noch mutable `Proposed`-ADRs (`0011`, `0012`) beim nächsten inhaltlichen Anfassen.

**Reconciliation.** Titel-Form (`# ADR <Nr>: <Titel>`) und Superseded-Referenz (`<NNNN>-<slug>`, klickbar) folgen der bestehenden u-boot-Konvention, nicht dem Template-Wortlaut (`# ADR-NNNN:` bzw. `Superseded by ADR-NNNN`); die MADR-**Substanz** (Inline-Kopf-Felder, `Schärft`-Aufwärtskopplung, Alternativen/Fitness-Function/Re-Eval/Geschichte) wird übernommen. Die Umstellung ist ein Change Request am Vertrags-Stratum (Trigger: v3.5.1-Adoption).

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

### LH-FA-CLI-006.a — Zuordnung von Fehlern zu Exit-Codes

Verfeinert [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes).

Empfehlung für typische Fehlerzuordnung:

- `10` bei Struktur-, Namens- oder Konfigurationsvalidierungsfehlern
- `11` bei Umgebungsproblemen (z. B. fehlende Tools, Versionsinkompatibilität)
- `12` bei Laufzeitfehlern beim Ausführen von Docker/Compose-Operationen

Für alle fachlichen Fehler ist die Verwendung von `10`, `11` oder `12` bindend.
Nicht-fachliche Fehler dürfen standardmäßig mit `1` codiert werden, wenn eine feinere technische Klassifikation nicht sinnvoll ist.
`13` bis `15` dürfen zusätzlich verwendet werden, wenn deren Bedeutung für den aufrufenden Kontext explizit dokumentiert ist.
`16` bis `19` sind in der aktuellen Spezifikation nicht zu verwenden.

### LH-FA-CLI-007.a — Konvention für Diagnose-Codes und Kopplung des Status

Verfeinert [`LH-FA-CLI-007`](lastenheft.md#lh-fa-cli-007--dry-run).

Konvention für `diagnostics[*].code`: LH-Kennung der verursachenden Anforderung (z. B. [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-CLI-007`](lastenheft.md#lh-fa-cli-007--dry-run)). Tool-interne Codes ohne LH-Bezug dürfen nur dann verwendet werden, wenn ihre Bedeutung in der Dokumentation festgehalten ist.

Das Feld `status` ist an den höchsten in `diagnostics` enthaltenen `level` gekoppelt: enthält `diagnostics` mindestens einen `error`-Eintrag, ist `status == "error"`; enthält es mindestens einen `warn`-Eintrag (und keinen `error`), ist `status == "warn"`; andernfalls `status == "ok"`. Diese Regel gilt für alle `--json`-Ausgaben ([`LH-FA-CLI-007`](lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](lastenheft.md#lh-fa-cli-008--diff-ausgabe), [`LH-NFA-USE-004`](lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)).

### LH-FA-INIT-004.a — Relevante Dateien der Projekterkennung

Verfeinert [`LH-FA-INIT-004`](lastenheft.md#lh-fa-init-004--bestehendes-projekt-erkennen).

Relevante Dateien sind mindestens:

- `u-boot.yaml`
- `compose.yaml`
- `.env.example`
- `README.md`
- `CHANGELOG.md`
- `.gitignore`
- `docs/`
- `scripts/`
- `docker/`
- `.devcontainer/devcontainer.json`

### LH-FA-INIT-005.b — Backup-Verfahren bei `--backup`

Verfeinert [`LH-FA-INIT-005`](lastenheft.md#lh-fa-init-005--überschreibschutz).

- `--backup` – bestehende Datei als `<name>.bak` sichern und ersetzen; ist `<name>.bak` bereits vorhanden, wird automatisch `<name>.bak.1`, `<name>.bak.2`, ... verwendet (kleinster freier numerischer Suffix), ohne vorhandene Backups zu überschreiben.
- Für bestehende Verzeichnisse (z. B. `docs/`, `scripts/`, `docker/`, `.devcontainer/`) wird der komplette Verzeichnisbaum rekursiv als `<name>.bak*` gesichert und innerhalb derselben Operation ersetzt; bei Fehlern während des Ersetzens muss ein Rollback auf den ursprünglichen Zustand durchgeführt werden (POSIX-Atomarität für rekursive Bäume wird nicht garantiert).

### LH-FA-DEV-007.b — Ausstattung bei `nestedRuntime: podman`

Verfeinert [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer).

- Bei `podman`: rootless Podman, ein `docker`-Kompatibilitäts-Alias, Zugriff auf `/dev/fuse` und `/dev/net/tun` (Netzwerk der verschachtelten Container über `pasta`; ohne das Device scheitert jedes `RUN`; beide Devices müssen auf dem Container-Host vorhanden sein, fehlt eines, lehnt die Container-Engine den Start ab (`error gathering device information while adding custom device`), das Startscript läuft dann nicht; `doctor` prüft `/dev/net/tun` nicht und weist ein fehlendes `/dev/fuse` nur aus, soweit der Host es bestimmen kann), subuid-/subgid-Bereiche für den Container-Benutzer und ein Volume für den Storage. `runArgs` enthält sechs Optionen: `--cap-add=SYS_ADMIN`, `--security-opt=seccomp=unconfined`, `--security-opt=apparmor=unconfined`, `--security-opt=systempaths=unconfined`, `--device=/dev/fuse`, `--device=/dev/net/tun`; jede weist die Befehlsausgabe einzeln aus. Das Dockerfile installiert `podman`, `uidmap`, `fuse-overlayfs`, `passt` und `ca-certificates` ohne empfohlene Pakete; der Alias ist ein Symlink `docker` → `podman`. Daraus folgt für Builds: `podman build` und `docker build` sind derselbe Aufruf mit der in Podman eingebauten Buildah-Bibliothek; das Kommando `buildah`, `docker buildx`, `docker compose` und BuildKit-Frontends (`# syntax=`) sind nicht Teil der Ausstattung. Ein Dockerfile mit `FROM`, `RUN` und `COPY` baut kalt, aus dem Layer-Cache und mit `--no-cache`, mit und ohne Egress-Restriktion (Default-Allowlist; gemessen mit `docker.io/library/alpine:3` und `debian:stable-slim`); Der Kurzname `alpine` löst auf, `mycompany/img` nicht (gemessen; `unqualified-search-registries` ist nicht gesetzt, andere Kurznamen sind nicht gemessen), Dockerfiles verwenden voll qualifizierte Image-Namen. Fehlt `/dev/fuse` im Container (der Container wurde ohne das Device gestartet), baut es mit `vfs`-Storage. Der Beleg sind Integrationstests; sie laufen nicht auf Hosts, die verschachtelte User-Namespaces blockieren (`kernel.apparmor_restrict_unprivileged_userns=1`, etwa Ubuntu 24.04 mit Standardeinstellung), dort endet das Startscript mit Exit-Code `11`.
- Bei `docker`: Docker-Daemon, `docker`-CLI sowie Buildx- und Compose-Plugin aus dem Docker-APT-Repository (Schlüssel per Fingerprint geprüft); der Container-Benutzer gehört der Gruppe `docker` an, der Daemon startet beim Containerstart (`postStartCommand`) und gilt nach spätestens 60 Sekunden als nicht betriebsbereit; `runArgs` enthält `--privileged` und, nur mit Egress-Restriktion, `--cap-add=NET_ADMIN`, sonst keine Fähigkeitsoptionen; der Daemon-Speicher (`/var/lib/docker`) liegt in einem Volume `<projekt>-docker-${devcontainerId}`; kein Socket-Mount.

### LH-FA-DEV-008.a — Verhalten bei nicht gewährbarer Capability

Verfeinert [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion).

- Ist die nötige Capability nicht gewährbar, greift die Degradationstabelle der Verfeinerung [`LH-FA-DEV-007.a`](#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten); die Capability wird beim Containerstart geprüft (Warnung und Wegfall der Restriktion, bei `onUnavailable: fail` Exit-Code `11`), und das Startscript schreibt das Ergebnis in die Statusdatei `egress`. `u-boot doctor` weist den Wegfall der Restriktion vom Host aus (Eintrag `devcontainer.sandbox.fallback.egress`, [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)) und prüft zusätzlich die Konfiguration (Schlüssel ohne Sandbox-Profil: `warn`). Die Capability selbst bestimmt der Host nicht, er liest nur das gemeldete Ergebnis; ohne lesbare Statusdatei ist der Zustand nicht ermittelbar.

### SPEC-029 — Mindestumfang des Multi-Stage-Dockerfiles

<a id="lh-fa-build-001b--mindestumfang-des-multi-stage-dockerfiles"></a>

Mindestumfang:

- BuildKit-Direktive in der ersten Zeile des Dockerfiles.
- Pflicht-Stages: `deps`, `compile`, `test`, `lint`, `coverage`, `build` und `runtime`; die Coverage-Stage erzwingt den Coverage-Schwellwert, die Runtime-Stage ist das minimale Endimage ([`LH-NFA-PORT-004`](lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image)).
- Jede Stage ist ein eigenständiges Build-Ziel und per `docker build --target <stage>` einzeln baubar.

### LH-NFA-PORT-004.b — Pflichten der Runtime-Stage

<a id="lh-fa-build-002b--pflichten-der-runtime-stage"></a>

Verfeinert [`LH-NFA-PORT-004`](lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image).

- Minimales Base-Image ohne Shell.
- Ausführung als Non-root-Benutzer.
- `ENTRYPOINT` zeigt auf das im `build`-Stage erzeugte Binary.
- OCI Image Labels für Quelle, Beschreibung, Lizenz und Titel sind gesetzt.
- Keine Build-Toolchain im Endimage; alle Build-Artefakte stammen aus dem `build`-Stage.

### SPEC-030 — Overrides der Build-Args

<a id="lh-fa-build-003a--overrides-der-build-args"></a>

Overrides erfolgen über `docker build --build-arg <NAME>=<value>` bzw. die korrespondierende Makefile-Variable.

### SPEC-031 — Regeln des Docker-only-Workflows

<a id="lh-fa-build-007a--regeln-des-docker-only-workflows"></a>

- Alle MVP- und V1-Pflicht-Targets aus [`SPEC-008`](#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) müssen ausschließlich `docker build`, `docker run` oder die Aggregation anderer solcher Targets aufrufen.
- Voraussetzung am Host: Docker Engine und `make`. `make` ist ein bewusster Carveout zu [`LH-NFA-PORT-002`](lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten) (weit verbreitet, einzige zusätzliche Host-Abhängigkeit neben Docker). Eine Go-Toolchain am Host darf für Standard-Targets nicht vorausgesetzt werden.
- Carveouts (z. B. ein Bash-Skript, das nicht containerisiert wird) sind im `Makefile`-Header explizit zu dokumentieren.

### SPEC-032 — Bootstrap-Verhalten der Coverage-Stage

<a id="lh-fa-build-008a--bootstrap-verhalten-der-coverage-stage"></a>

- Default-Schwellwert `90` Prozent (`ARG COVERAGE_THRESHOLD=90`, Makefile `THRESHOLD ?= 90`); in der Bootstrap-Phase war er `0`.
- Der Override-Pfad `make coverage-gate THRESHOLD=…` muss funktionieren.
- Leere Coverage darf in der Bootstrap-Phase nicht zu einem falschen Grün führen, das echte Test-Failures maskiert; der `go test`-Exit-Code wird über `set -o pipefail` o. ä. an die Gate-Logik durchgereicht.

### SPEC-033 — Go-Layout des Repositories

<a id="lh-fa-build-009a--go-layout-des-repositories"></a>

- Modul-Pfad in `go.mod`: `github.com/pt9912/u-boot`.
- Implementierungspakete leben unter `./internal/...`; öffentlich konsumierbare Pakete unter `./pkg/...` werden im MVP nicht erzeugt.
- CLI-Entry-Points unter `./cmd/<binary>/`; das primäre Binary heißt `uboot` (Verzeichnis `./cmd/uboot/`, Go-konform ohne Bindestrich), wird beim Build aber als `u-boot` ausgeliefert (`-o /out/u-boot`).
- Unit-Tests stehen als `*_test.go` neben dem produktiven Code im selben Paket.
- Coverage-Messung ([`SPEC-025`](#spec-025--stages-des-multi-stage-dockerfiles), [`SPEC-032`](#spec-032--bootstrap-verhalten-der-coverage-stage)) bezieht sich auf `./internal/...`; `./cmd/...` ist bewusst ausgeschlossen, weil dort nur dünne Wireup-Logik liegt.

### SPEC-034 — Enforcement der Import-Regeln

<a id="lh-fa-arch-003a--enforcement-der-import-regeln"></a>

- Die Regeln werden im `lint`-Stage ([`SPEC-025`](#spec-025--stages-des-multi-stage-dockerfiles)) per `golangci-lint` mit `depguard` durchgesetzt; Verstöße sind PR-blockierend.
- Die `depguard`-Konfiguration in `.golangci.yml` ist deckungsgleich mit den Import-Regeln der Anforderung zu halten; Drift wird im Review zurückgewiesen.
- `//nolint:depguard`-Pragmas sind verboten. Carveouts werden zentral in `.golangci.yml` mit `Why:`-Kommentar dokumentiert.
- `depguard`-Regeln gelten production-only; `*_test.go`-Dateien sind ausgenommen, damit Tests Fakes und Test-Libraries (`testify`, …) frei nutzen können.
- Solange einzelne Schichten noch keine produktiven Pakete enthalten, dürfen `depguard`-Regelblöcke aktiv sein und nichts treffen — die Schicht-Regeln greifen automatisch, sobald das erste produktive Paket angelegt wird.

### SPEC-035 — Pflichten des Lint-Profils

<a id="lh-qa-004b--pflichten-des-lint-profils"></a>

- Die Konfiguration lebt in `.golangci.yml` (v2-Schema).
- Schwellen und Linter-Settings sind in `.golangci.yml` konfiguriert und in abgeleiteter Quality-Doku zu erklären; bei Drift gewinnt die Anforderung, Config und Doku sind anzupassen.
- `//nolint`-Pragmas sind verboten. Pro-Pfad-Carveouts (z. B. Tests, `cmd/uboot`) werden zentral in `.golangci.yml` unter `issues.exclude-rules` mit `Why:`-Kommentar dokumentiert.
- Verstöße brechen den `lint`-Stage ([`SPEC-025`](#spec-025--stages-des-multi-stage-dockerfiles)) und damit `make gates`/`make ci`/`make fullbuild`.
- Die konkrete Linter-Auswahl muss die hier genannten Anforderungen
  und die Architekturgrenzen aus [`SPEC-034`](#spec-034--enforcement-der-import-regeln) abdecken.

### SPEC-036 — Stabilisierung von `u-boot up`: Takt, Zustandsklassen und Entscheidungsmatrix

Verfeinert [`LH-FA-UP-001`](lastenheft.md#lh-fa-up-001--umgebung-starten) und [`LH-FA-UP-003`](lastenheft.md#lh-fa-up-003--startstatus-anzeigen).

- **Takt:** Der Compose-Status wird alle `500 ms` abgefragt; jede einzelne TCP-Erreichbarkeitsprobe ist auf `300 ms` begrenzt. Beide Werte sind feste Konstanten ohne Flag- oder Projekt-Override.
- **Zustandsklassen:** Der vom Compose gemeldete Container-Zustand wird ohne Beachtung von Groß-/Kleinschreibung und umgebendem Leerraum klassifiziert:

  | Compose-Zustand | Klasse |
  |---|---|
  | `running` | running |
  | `restarting` | restarting |
  | `created`, `starting`, `paused` | starting |
  | `exited`, `dead`, `removing`, `removed` | dead |
  | jeder andere Wert | unknown |

  `unknown` ist weich: Der Service gilt als noch nicht stabilisiert, `up` bricht nicht ab. Je Service und Aufruf wird einmalig eine `warn`-Diagnose `up.state.<service>.unknown` erzeugt.
- **Entscheidung je Service** (in alphabetischer Reihenfolge der Service-Namen aus `compose.yaml`, damit bei mehreren Fehlern stets derselbe zuerst gemeldet wird):
  - `dead` → fehlgeschlagen, das Warten endet sofort.
  - `restarting` → Der Zähler der aufeinanderfolgenden `restarting`-Beobachtungen steigt; bei `3` gilt der Service als fehlgeschlagen, darunter wird weiter gewartet. Jede andere Beobachtung setzt den Zähler auf `0` zurück. Grenze: Ein Restart-Loop, der zwischendurch `running` zeigt, wird auf dem Pfad ohne Healthcheck und ohne probebaren Port nicht erkannt.
  - `starting` und `unknown` → weiter warten.
  - `running` mit deklariertem Healthcheck (`healthcheck.disable` nicht `true`): solange der Health-Status nicht `healthy` ist (ohne Beachtung der Groß-/Kleinschreibung), wird weiter gewartet und nicht geprobt; bei `healthy` gilt der Service als stabilisiert. Die deklarierten TCP-Ports werden dann dennoch geprobt; ein Fehlschlag verhindert die Stabilisierung nicht, sondern erzeugt einmalig je Service die `warn`-Diagnose `up.port.<service>.unreachable` (der Healthcheck dominiert).
  - `running` ohne Healthcheck: Der Service gilt als stabilisiert, wenn alle probebaren Ports antworten; ohne probebaren Port genügt `running`.
  - In `compose.yaml` deklarierte, in der Compose-Statusabfrage (noch) nicht sichtbare Services gelten als noch nicht stabilisiert; zusätzliche, nur in der Statusabfrage sichtbare Services werden ignoriert.
- **Reichweite der `warn`-Diagnosen** (`up.state.<service>.unknown`, `up.port.<service>.<index>`, `up.port.<service>.unreachable`): Sie erreichen den Aufrufer nur in der Antwort einer erfolgreich stabilisierten Ausführung. Bei Fehlschlag (Exit-Code `12`) und bei Zeitüberschreitung liefert `up` nur den zuletzt beobachteten Service-Status, keine Diagnosen.
- **Abschluss:** Sind alle deklarierten Services stabilisiert, endet `up` erfolgreich. Ein fehlgeschlagener Service beendet `up` mit Exit-Code `12` (Compose-Laufzeitfehler). Ist die Wartezeit nach einer Statusabfrage abgelaufen, endet `up` mit Exit-Code `12` und nennt die noch ausstehenden Services; es findet also mindestens eine Statusabfrage statt.
- **Nicht probebare Portangaben:** Für jeden Eintrag in `ports`, der nach [`LH-FA-UP-001.a`](#lh-fa-up-001a--erkennung-probebarer-ports-in-composeyaml) nicht probebar ist, wird vor dem Warten eine `warn`-Diagnose `up.port.<service>.<index>` ausgegeben (`<index>` ist die Position im `ports`-Array); das Warten stützt sich dann auf Healthcheck oder `running`.
- **`--timeout=0`:** Es findet keine Statusabfrage und keine Probe statt; die Antwort trägt `stabilized=false` und eine `info`-Diagnose `up.fire-and-forget`.

### LH-FA-UP-001.a — Erkennung probebarer Ports in `compose.yaml`

Verfeinert [`LH-FA-UP-001`](lastenheft.md#lh-fa-up-001--umgebung-starten) (Erreichbarkeitsprüfung auf `localhost` für TCP-Ports).

Aus einem Eintrag des `ports`-Arrays eines Services wird ein Probe-Ziel (Host und Host-Port) abgeleitet, oder der Eintrag gilt als nicht probebar (Warnung nach [`SPEC-036`](#spec-036--stabilisierung-von-u-boot-up-takt-zustandsklassen-und-entscheidungsmatrix), kein Abbruch):

| Angabe | Ergebnis |
|---|---|
| Zahl (`5432`) oder Kurzform ohne Doppelpunkt (`"5432"`) | probe `localhost:5432` (die Zahl gilt als Host-Port) |
| `"HOST:CONTAINER"` (`"5432:5432"`) | probe `localhost:HOST` |
| `"IP:HOST:CONTAINER"` (`"127.0.0.1:5432:5432"`) | probe `IP:HOST` |
| Kurzform mit Protokoll-Suffix außer `tcp` (`"5432:5432/udp"`) | nicht probebar |
| Bereich im Host-Port (`"5000-5010:5000-5010"`) | nicht probebar |
| Langform (Mapping) mit `published`, `protocol` fehlend oder `tcp`, optional `host_ip` | probe `host_ip:published` (ohne `host_ip`: `localhost`); `published` darf Zahl oder Zahl-String sein |
| Langform mit `protocol` außer `tcp`, ohne `published` oder mit Bereich | nicht probebar |
| jede andere Form (leerer String, nicht numerischer Host-Port, mehr als drei Segmente einschließlich IPv6-Adressen in Klammern, gebrochene Zahl, Wahrheitswert, Liste) | nicht probebar |

Der Protokoll-Vergleich ist unabhängig von der Groß-/Kleinschreibung. Ein Host-Port außerhalb `1`..`65535` ist nicht probebar. Die Host-Angaben `""`, `0.0.0.0` und `::` (an alle Schnittstellen gebunden) werden zu `localhost` normalisiert.

### SPEC-038 — Bestätigungspfad von `u-boot down --volumes`

Verfeinert [`LH-FA-UP-004`](lastenheft.md#lh-fa-up-004--umgebung-stoppen) und wendet [`LH-FA-CLI-005A`](lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) auf den destruktiven Pfad an.

| `--volumes` | `--yes` | `--no-interactive` | Verhalten |
|---|---|---|---|
| nein | beliebig | beliebig | `down` läuft ohne Rückfrage |
| ja | ja | beliebig | `down` läuft ohne Rückfrage |
| ja | nein | ja | Abbruch mit Exit-Code `10`; weder Rückfrage noch Compose-Aufruf |
| ja | nein | nein | interaktive Rückfrage: Zustimmung → `down` läuft; Ablehnung → Abbruch mit Exit-Code `10`; Fehler beim Einholen der Antwort → Abbruch |

Bei `--json` wird nicht gefragt: Die Antwort gilt als abgelehnt, destruktives `--volumes` verlangt dort also ausdrücklich `--yes` ([`LH-NFA-USE-004`](lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)).

### SPEC-039 — Zustandsmaschinen der Generatoren (`u-boot generate`)

Verfeinert [`LH-FA-GEN-001`](lastenheft.md#lh-fa-gen-001--generate-befehl) bis [`LH-FA-GEN-005`](lastenheft.md#lh-fa-gen-005--idempotenz) und [`LH-AK-007`](lastenheft.md#lh-ak-007--changelog-generator).

Alle Generatoren verlangen ein initialisiertes Projekt (`u-boot.yaml` vorhanden) und schreiben in einen verwalteten Bereich mit dem Namen `init` (Markierungen nach [`SPEC-014`](#spec-014--markierungsformat-verwalteter-bereiche-je-dateityp)):

| Artefakt | Datei | Markierungsstil |
|---|---|---|
| `changelog` | `CHANGELOG.md` | HTML-Kommentar |
| `readme` | `README.md` | HTML-Kommentar |
| `env-example` | `.env.example` | `#`-Kommentar |
| `devcontainer` | `.devcontainer/devcontainer.json` | `//`-Kommentar |
| `devcontainer` | `.devcontainer/Dockerfile` | `#`-Kommentar |
| `devcontainer` (nur bei Sandbox-Profil mit der jeweiligen Fähigkeit) | `.devcontainer/sandbox-init.sh` (`nestedRuntime: podman`), `.devcontainer/docker-init.sh` (`nestedRuntime: docker`), `.devcontainer/egress-init.sh` (Egress-Restriktion) | `#`-Kommentar |

**Zustände je Datei** (`readme`, `env-example`, je Devcontainer-Datei):

| Vorgefundener Zustand | Ergebnis |
|---|---|
| Datei fehlt | vollständig aus der Vorlage anlegen (`created`) |
| Datei mit verwaltetem Bereich, Inhalt entspricht der Vorlage | nichts schreiben (`no-op`) |
| Datei mit verwaltetem Bereich, Inhalt weicht ab | nur den Bereich ersetzen (`updated-block`); Inhalt außerhalb bleibt byte-identisch |
| Datei ohne verwalteten Bereich | Abbruch, Exit-Code `10`, nichts wird geschrieben |
| Datei mit fehlerhaft markiertem Bereich | Abbruch, Exit-Code `10`, nichts wird geschrieben |

Zeilenenden werden für den Vergleich auf LF normalisiert, damit eine CRLF-Datei nicht als geändert gilt.

**`changelog`** ersetzt einen abweichenden Bereich nie, weil der Nutzer dort Einträge pflegt ([`LH-AK-007`](lastenheft.md#lh-ak-007--changelog-generator)): Fehlt die Datei, wird sie angelegt; fehlt der Bereich oder ist er fehlerhaft markiert, gilt Exit-Code `10`; entspricht der Bereich der Vorlage oder weicht er ab, bleibt er unverändert. Weicht der Bereich ab, wird nur dann repariert, wenn in der **gesamten Datei** — der verwaltete Bereich eingeschlossen — kein Abschnitt `## [Unreleased]` steht (Überschriften in Codeblöcken zählen nicht) und **nach der END-Markierung** ein Release-Abschnitt `## [<Version>]` (ebenfalls außerhalb von Codeblöcken) vorhanden ist: Dann wird vor dem ersten solchen Release-Abschnitt ein leerer Abschnitt `## [Unreleased]` mit den Unterabschnitten `### Added`, `### Changed` und `### Fixed` eingefügt (`repaired-manual`). Steht `## [Unreleased]` irgendwo in der Datei, oder gibt es nach der END-Markierung keinen Release-Abschnitt, bleibt die Datei unverändert (`no-op`). Folge: Ein nutzerbearbeiteter Bereich, der `## [Unreleased]` behält, wird nicht repariert; ein Release-Abschnitt innerhalb des Bereichs verankert keine Reparatur. Folge: Ändert sich die Vorlage oder der Projektname, wird ein bestehender Bereich nicht nachgezogen.

**`devcontainer`** arbeitet in zwei Phasen: Zuerst werden alle Dateien klassifiziert, ohne zu schreiben; nur wenn keine Datei den Abbruch-Zustand hat, wird geschrieben. Ein Konflikt hinterlässt also keinen Teilzustand. `forwardPorts` stammt aus den Container-Ports der aktiven Services in `compose.yaml` ([`LH-FA-DEV-005`](lastenheft.md#lh-fa-dev-005--ports)) und teilt sich die Quelle mit der `doctor`-Prüfung `devcontainer.forwardPorts.consistency`. Ist `compose.yaml` nicht parsebar, gilt Exit-Code `10`. Die Freigabe externer Feature-Quellen ([`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features)) wird erst nach erfolgreicher Erzeugung in `u-boot.yaml` übernommen.

### SPEC-040 — Konfigurationspfade und Schreibpipeline von `u-boot config`

Verfeinert [`LH-FA-CONF-005`](lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern) und [`LH-FA-CONF-002`](lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration).

**Erlaubte Pfade** (jeder andere Pfad: fachlicher Validierungsfehler, Exit-Code `10`):

| Pfad | Lesen | Schreiben | Werte |
|---|---|---|---|
| `project.name` | ja | ja | gültiger Projektname ([`SPEC-004`](#spec-004--projektnamen-muster)) |
| `devcontainer.enabled` | ja | ja | Wahrheitswert |
| `devcontainer.featureSources.allow` | ja | ja (Liste) | kommagetrennte Quellen |
| `devcontainer.features.<feature>.enabled` | ja | ja | Wahrheitswert |
| `devcontainer.features.<feature>.source` | ja | ja | gültige Quelle, die wörtlich in `devcontainer.featureSources.allow` steht |
| `devcontainer.features.<feature>.version` | ja | ja | nicht-leerer String |
| `devcontainer.user.uid` | ja | ja | Ganzzahl `1`..`65535` |
| `devcontainer.profile` | ja | ja | `default` \| `sandbox` |
| `devcontainer.sandbox.nestedRuntime` | ja | ja | `none` \| `podman` \| `docker` |
| `devcontainer.sandbox.onUnavailable` | ja | ja | `warn` \| `fail` |
| `devcontainer.sandbox.repository` | ja | ja | Clone-URL aus den Zeichen `A–Z a–z 0–9 . _ ~ : / @ % + = -`, nicht mit `-` beginnend, ohne eingebettete Zugangsdaten ([`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)) |
| `devcontainer.sandbox.egress.enabled` | ja | ja | Wahrheitswert |
| `devcontainer.sandbox.egress.allow` | ja | ja (Liste) | kommagetrennte Hostnamen |
| `services.<svc>.enabled` | ja | **nein** | — |

`services.<svc>.enabled` ist nur lesbar: Ein Schreibversuch endet mit Exit-Code `10` und dem Hinweis auf `u-boot add <svc>`, weil das bloße Umschalten die Zustandsregeln aus [`LH-FA-ADD-005.a`](#lh-fa-add-005a--zustandsregeln-für-registrierte-und-aktive-services) umginge. `<svc>` und `<feature>` werden auf Namensformat geprüft. Wahrheitswerte nehmen die Schreibweisen `true`, `false`, `1`, `0`, `T`, `F` u. ä. an und werden kanonisch als `true`/`false` geführt.

**Schreibpipeline `config set`** — jede Stufe läuft vor dem Schreiben; schlägt eine fehl, bleibt `u-boot.yaml` byte-identisch:

1. Pfad nicht schreibbar → Abbruch (vor jedem Dateizugriff).
2. Wert in den erwarteten Typ überführen und prüfen (Wertetabelle oben).
3. Entspricht der kanonische neue Wert dem bestehenden, wird nichts geschrieben (No-Op, alter und neuer Wert gleich).
4. Wert im Dokument ersetzen, Kommentare und Formatierung bleiben erhalten.
5. Das Ergebnis muss sich wieder als Projektkonfiguration einlesen lassen (sonst Schemafehler, Exit-Code `10`).
6. Fachliche Prüfung des Ergebnisses: Der geschriebene Wert muss gebunden und gültig sein; bei `.source` muss die Quelle wörtlich (Groß-/Kleinschreibung und abschließender Schrägstrich zählen) in `devcontainer.featureSources.allow` stehen, `--yes` ersetzt diese Freigabe nicht.
7. Schreiben.

Die beiden Listenpfade (`devcontainer.featureSources.allow`, `devcontainer.sandbox.egress.allow`) hängen die neuen Einträge an die bestehende Liste an und entdoppeln sie; die Datei wird dabei neu serialisiert, Kommentare gehen verloren. `config set devcontainer.features.<feature>.enabled true` für ein Feature, das weder im eingebauten Katalog ([`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen)) steht noch eine `source` trägt, ist kein Fehler, erzeugt aber eine `warn`-Diagnose ([`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features)): Der Generator überspringt es, `doctor` meldet es. `config` ohne Argumente gibt `u-boot.yaml` unverändert (mit Kommentaren) aus.

### SPEC-041 — Schweregrade der Prüfung `compose.yaml.valid`

Verfeinert [`LH-FA-DIAG-002`](lastenheft.md#lh-fa-diag-002--lokale-voraussetzungen-prüfen) und [`LH-AK-001`](lastenheft.md#lh-ak-001--minimaler-init-flow).

| Befund | Schweregrad |
|---|---|
| `compose.yaml` fehlt | `warn` |
| `compose.yaml` nicht lesbar | `error` |
| `compose.yaml` kein gültiges YAML | `error` |
| gültig, aber ohne Einträge unter `services` | `warn` |
| gültig mit mindestens einem Service | `ok` (Meldung nennt die Anzahl) |

Der leere `services`-Block ist `warn`, weil ihn ein frisches `u-boot init` erzeugt und der Ablauf `init`, dann `doctor` nach [`LH-AK-001`](lastenheft.md#lh-ak-001--minimaler-init-flow) keinen `error` enthalten darf.

### SPEC-042 — Aufruf der Compose-Engine und Preflight

Verfeinert [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes) (Trennung der Codes `11` und `12`) und [`LH-FA-UP-001`](lastenheft.md#lh-fa-up-001--umgebung-starten) (`--timeout=0` als Fire-and-forget).

- **Preflight** vor jedem Compose-Aufruf, in dieser Reihenfolge: (1) das `docker`-Binary ist auffindbar, (2) der Daemon antwortet, (3) das Compose-Plugin ist installiert. Scheitert ein Schritt, endet der Befehl mit Exit-Code `11` (Umgebungsproblem), ohne dass Compose aufgerufen wurde. Scheitert erst der Compose-Aufruf selbst, gilt Exit-Code `12`. Ausnahme: Ist der Kontext bei `logs`, `config`, `volume ls`, `volume rm` und `--dry-run up -d` abgebrochen oder abgelaufen, wird der Kontextfehler unverpackt durchgereicht statt als Exit-Code `12` gemeldet; bei `logs` behandelt der Aufrufer `context.Canceled` und `context.DeadlineExceeded` als Erfolg (Exit-Code `0`). Die Vorschau `--dry-run up -d` ist best-effort: Scheitert der Trockenlauf selbst, bleibt der Plan leer, ohne Fehler.
- **Aufrufformen:** `docker compose -f <Projektverzeichnis>/compose.yaml` mit `up` (mit `-d` im Detach-Modus), `down` (mit `-v` bei `--volumes`), `ps --format json`, `logs` und `config --format json`; ferner `docker volume ls` und `docker volume rm <Name>` sowie für die Vorschau `--dry-run up -d`. Der Fortschritt von `up` und `down` (Compose-Ausgabe auf stderr) wird an die Fortschrittsausgabe weitergereicht.
- **Fire-and-forget:** Nach `up` im Detach-Modus fragt der Adapter den Status nicht ab; das Warten und jede Statusabfrage gehören allein dem Aufrufer ([`SPEC-036`](#spec-036--stabilisierung-von-u-boot-up-takt-zustandsklassen-und-entscheidungsmatrix)).
- **Netzprobe:** Die TCP-Probe läuft mit der übergebenen Frist und dem Kontext; Abbruch oder Fristablauf des Kontexts sind am Fehler erkennbar (`context.Canceled`, `context.DeadlineExceeded`).

### SPEC-043 — Verbosity: Pegelabbildung, Rangfolge und Wirkung von `--quiet`

Verfeinert [`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Stufen, Wirkung von `--quiet`, Mehrfachangabe) und berührt [`LH-FA-CLI-007`](lastenheft.md#lh-fa-cli-007--dry-run) (Vorschau), [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) und [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) (Ausweisungen), [`LH-NFA-PERF-002`](lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker) (Fortschritt von `up`) sowie [`LH-FA-INIT-005`](lastenheft.md#lh-fa-init-005--überschreibschutz) und [`LH-FA-CLI-005A`](lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (Pfad-Zusammenfassung bei `--force`).

**Pegelabbildung** (Log-Pegel des Loggers, der nach dem Parsen der Optionen und vor der Ausführung des Subkommandos gesetzt wird):

| Gewählte Verbosity | Log-Pegel |
|---|---|
| `--quiet` | Fehler (`slog.LevelError`) |
| keine Option | Info (`slog.LevelInfo`) |
| `--verbose` | Debug (`slog.LevelDebug`) |
| `--debug` | Debug (`slog.LevelDebug`) |

`--verbose` und `--debug` wirken auf den Log-Pegel gleich; eine eigene Stufe zwischen Info und Debug gibt es nicht. Eine Quelle für den Pegel außerhalb dieser drei Optionen (Umgebungsvariable, Konfiguration) gibt es nicht.

**Rangfolge:** Bei mehreren Optionen gewinnt die zuletzt auf der Kommandozeile angegebene, unabhängig davon, ob sie vor oder nach dem Subkommando steht: `--quiet --verbose` ist `--verbose`, `--verbose --quiet` ist `--quiet`, `--quiet --verbose --quiet` ist `--quiet`. Eine Mehrfachangabe wird weder abgewiesen noch verändert sie den Exit-Code ([`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes)). Die drei Optionen teilen sich einen einzigen Zustand, den jede Angabe überschreibt; es gibt keine feste Rangfolge zwischen ihnen. `--quiet` gilt für die Ausgabe-Filter genau dann, wenn es die zuletzt angegebene Option ist: Ein nachfolgendes `--verbose` oder `--debug` hebt es auf, auch für Zusammenfassungen, Warnungen, den `doctor`-Bericht und die Pufferung der Compose-Ausgabe, nicht nur für den Log-Pegel.

**Wirkung von `--quiet` auf die Ausgabe:**

- Die ändernden Subkommandos (`init`, `generate`, `add`, `remove`, `config set`, `up`, `down`) geben nur Fehler aus; ihre Erfolgs- und Zusammenfassungsmeldungen entfallen, bei `remove` auch die Purge-Warnungen.
- `doctor` schreibt nur Fehler-Einträge, Einträge eines Degradations-Fallbacks und die Warn-Einträge, die eine Sicherheitseigenschaft des Sandbox-Profils betreffen (Ausnahme 1), ohne Kopf- und Summary-Zeile; der Exit-Code trägt das Ergebnis.
- Die Abfragebefehle `config get`, `config list`, `logs` und `template list` behalten ihre Ausgabe.
- Die Ausgabe von `--dry-run` und `--diff` (Änderungsplan, Diff) bleibt ([`LH-FA-CLI-007`](lastenheft.md#lh-fa-cli-007--dry-run)); die Erfolgsmeldung eines tatsächlichen Schreibens entfällt auch bei `--diff` ohne `--dry-run`.
- Vier Ausnahmen bleiben sichtbar: (1) die Ausweisung von Sicherheitslockerungen des Sandbox-Profils und eines im Container aktivierten Login-Dienstes ([`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion)) in der Warnungsausgabe von `generate devcontainer` (nur dieser Befehl erzeugt sie) sowie die Warn-Einträge von `doctor`, die eine Sicherheitseigenschaft des Sandbox-Profils betreffen: die Quelle der Zugangsdaten ([`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)) und die Egress-Restriktion ([`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion)), darunter der Eintrag, dass sie ohne Sandbox-Profil nicht greift; (2) die Ausweisung der Degradations-Fallbacks (Container ohne Daemon, `vfs`-Storage, Wegfall der Egress-Restriktion) in `doctor` (Einträge `devcontainer.sandbox.fallback.*`) und in der Warnungsausgabe von `generate devcontainer` und `init` nach [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung); ein Zustand, den der Host nicht lesen kann, ist „nicht ermittelbar“ und kein Fallback, und sein Eintrag entfällt unter `--quiet`; (3) der Compose-Fortschritt von `up` ([`LH-NFA-PERF-002`](lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker)); (4) die Pfad-Zusammenfassung, die `init --force` vor dem Überschreiben ausgibt ([`LH-FA-CLI-005A`](lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-INIT-005`](lastenheft.md#lh-fa-init-005--überschreibschutz)). Jede andere Warnung und jeder andere Hinweis entfällt.
- Die Einstufung hängt am Kennzeichen des einzelnen Eintrags (`SecurityRelaxation`, `DegradationFallback`), nicht an der Kennung der Prüfung oder Anforderung; dieselbe Kennung trägt gekennzeichnete und ungekennzeichnete Einträge. Im Einzelnen:
  - Gekennzeichnet sind in der Warnungsausgabe von `generate devcontainer` die Lockerungen von `nestedRuntime` und Egress-Restriktion (je Option eine Meldung), die Meldung, dass `egress.enabled` ohne Sandbox-Profil keine Wirkung hat (es entsteht keine Restriktion), und der Login-Dienst; in `doctor` die Warn-Einträge von `devcontainer.sandbox.credentials` (Klartext-Zugangsdaten, https-Clone ohne Token-Quelle, fehlende Referenz der Agenten-Variable) und von `devcontainer.sandbox.egress` (ohne Sandbox-Profil, zusammen mit `nestedRuntime: docker`).
  - Mit `DegradationFallback` gekennzeichnet sind der `warn`-Eintrag je `devcontainer.sandbox.fallback.*` bei `result=fallback` und die entsprechende Warnung von `generate devcontainer` und `init`; ungekennzeichnet bleiben die Einträge dieser Prüfungen bei `ok`, „nicht ermittelbar“ (`info`) und ungültigem Dokument (`warn`), der Eintrag bei `failed` ist ein Fehler und bleibt unter `--quiet` als Fehler sichtbar.
  - `nestedRuntime: docker` hat in `doctor` die Stufe OK (kein Warn-Eintrag; `--privileged` weist die Befehlsausgabe aus): Die Sicherheitslockerung des Docker-Daemons erscheint nur in der Befehlsausgabe.
  - Ungekennzeichnet sind die Meldungen „has no effect“ zu `nestedRuntime` in der Warnungsausgabe von `generate devcontainer` (Konfigurationsmeldung), der `doctor`-Eintrag, dass `nestedRuntime` ohne Sandbox-Profil keine Wirkung hat, die Drittanbieter-Warnung und der Hinweis auf das Feature `node` ([`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)) sowie die Warn-Einträge von `devcontainer.features.allowlist` und `devcontainer.features.drift`; sie entfallen. Dieselbe Aussage zur Egress-Restriktion ohne Sandbox-Profil ist in `doctor` wie in der Befehlsausgabe gekennzeichnet, also sichtbar.
- `down --quiet` hält die Compose-Ausgabe zurück und gibt sie nur aus, wenn `down` fehlschlägt, weil sie die Fehlerdetails trägt, die der Fehler selbst nicht nennt.
- Mit `--json` ändert `--quiet` nichts: Der Umschlag und seine `diagnostics[]` bleiben vollständig, und die Compose-Ausgabe von `down` wird nicht zurückgehalten.

### SPEC-044 — Features `sshd` und `claude-code` im Devcontainer: Warnungen und Zugangsdaten

Verfeinert [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) (Login-Dienst), [`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten des Agenten), [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Drittanbieter-Feature) und berührt [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts: [`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion)). Der Katalog steht in [`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen), die Wirkung von `--quiet` in [`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet).

**Aktiv** heißt: Ein Feature-Eintrag mit `enabled: true` löst auf die Katalogquelle von `sshd` bzw. `claude-code` auf; verglichen wird die aufgelöste Quelle über alle aktivierten Einträge, nicht der Schlüssel: Auch ein Eintrag mit anderem Schlüssel und `source` gleich der Katalogquelle ist aktiv. Eine `version`-Überschreibung ändert das nicht. Ein Eintrag, dessen aufgelöste Quelle eine andere ist, ist kein aktives `sshd` oder `claude-code`: Ein anderes SSH-Feature oder ein anderer Agent über `source` bleibt unausgewiesen.

**Sandbox-Regel.** Beide Features sind mit und ohne Sandbox-Profil wählbar; das Profil gibt keine externe Quelle frei, und der Katalogschlüssel braucht keinen Eintrag in `devcontainer.featureSources.allow`. Der Generator erzeugt für beide Features nur den Schlüssel unter `features`; keines von beiden ergänzt `runArgs`, `mounts` oder einen Startbefehl (`runArgs` stammt allein aus `nestedRuntime` und Egress-Restriktion).

**Zugangsdaten des Agenten.** Mit Sandbox-Profil und aktivem `claude-code` trägt `remoteEnv` neben `GIT_TOKEN` den Eintrag `"ANTHROPIC_API_KEY": "${localEnv:ANTHROPIC_API_KEY}"`, die Referenz auf die Umgebungsvariable des Hosts, nie ihren Wert. Ohne Sandbox-Profil erzeugt der Generator kein `remoteEnv`. Andere Variablen reicht er nicht durch; ein Anmeldeweg über Login im Container oder `CLAUDE_CODE_OAUTH_TOKEN` wird weder eingerichtet noch geprüft, und die Hosts der Anmeldung stehen nicht in der Default-Allowlist ([`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion)).

**Warnungen der Befehlsausgabe** (Erzeuger: `generate devcontainer`; Level `warn`, Subjekt `.devcontainer/devcontainer.json`):

| Auslöser | Code | Sicherheitslockerung (`--quiet`) | Inhalt |
|---|---|---|---|
| `sshd` aktiv, mit Sandbox-Profil | [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) | ja, bleibt sichtbar | Login-Dienst `sshd` ist im Container aktiv (Port 2222; Passwort- und Schlüssel-Login, Root-Login erlaubt; die Konten bleiben gesperrt, bis ein Passwort gesetzt oder ein Schlüssel eingetragen ist); das Profil begrenzt den Schaden und ist keine Isolationsgrenze. Ohne Sandbox-Profil entsteht keine Warnung. |
| `homebrew` oder `claude-code` aktiv, ohne eigene `source` | [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | nein, entfällt | Eine Meldung je Feature: Feature, Herausgeber, Quelle mit Version; sein Installationsskript läuft beim Bau des Images. Beide Profile. |
| `claude-code` aktiv, `node` nicht aktiv | [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) | nein, entfällt | `claude-code` installiert über npm, das Basis-Image enthält kein Node.js; `node` ist zu aktivieren, sonst schlägt der Bau des Images fehl. `u-boot` aktiviert `node` nicht selbst. Beide Profile. |

**Prüfung `devcontainer.sandbox.credentials`** (`u-boot doctor`; kein eigener Check-Code für Agenten, die Kennung steht in [`SPEC-017`](#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry)). Ohne Sandbox-Profil ist der Eintrag OK und die Prüfung übersprungen. Mit Sandbox-Profil meldet ein Lauf höchstens den ersten Befund in dieser Reihenfolge, jeweils `warn` mit dem Kennzeichen der Sicherheitseigenschaft nach [`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (es ändert weder Severity noch Exit-Code):

1. **Klartext in einer Projektdatei.** Gescannt werden `u-boot.yaml`, `compose.yaml`, `.env.example`, `.devcontainer/devcontainer.json`, `.devcontainer/Dockerfile`, `.devcontainer/sandbox-init.sh` und `.devcontainer/docker-init.sh`. Treffer sind: Zugangsdaten in einer http(s)-URL (Benutzer, optional Passwort, vor `@`); Git-Token-Präfixe `ghp_`, `gho_`, `ghu_`, `ghs_`, `ghr_` (mindestens 20 Zeichen danach), `github_pat_` und `glpat-` (je mindestens 20 Zeichen danach); die Zuweisung (`=` oder `:`) eines Werts von mindestens 12 Zeichen an `git_token`, `gh_token`, `github_token` oder `gitlab_token` (Schreibung und Unterstrich beliebig); ein Anthropic-Schlüssel `sk-ant-` mit mindestens 40 Zeichen aus `A–Z a–z 0–9 _ -` danach, dem kein solches Zeichen vorangeht; die Zuweisung (`=` oder `:`) eines Werts von mindestens 12 Zeichen an `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN` oder `CLAUDE_CODE_OAUTH_TOKEN`. Referenzen wie `${localEnv:GIT_TOKEN}` treffen nicht. Die Meldung nennt die Dateien, nie den Wert. **Grenze:** Der Zuweisungs-Ausdruck unterscheidet einen Platzhalter nicht von einem Wert; `ANTHROPIC_API_KEY=your-api-key-here` in `.env.example` löst die Warnung aus.
2. **https-Clone ohne Token-Quelle.** Die Clone-Quelle beginnt mit `https://` oder `http://`, und `devcontainer.json` enthält weder `${localEnv:GIT_TOKEN}` noch die Zeichenfolgen `/run/secrets` und `readonly` zusammen (Secret-Mount; die Kurzform `ro` der Mount-Option zählt nicht).
3. **Agenten-Referenz fehlt.** Ein Agenten-Feature ist aktiv, und `devcontainer.json` enthält die Referenz `${localEnv:<NAME>}` seiner Variablen nicht (`ANTHROPIC_API_KEY` bei `claude-code`); eine nicht lesbare Datei gilt als ohne Referenz.

Ohne Befund ist der Eintrag OK; bei aktivem Agenten-Feature nennt die Meldung die Variable und dass nur die Referenz durchgereicht wird.

### SPEC-045 — Statusdateien der Startscripte und ihre Ausweisung

Verfeinert [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Degradation, Ausweisung jedes Fallbacks), [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) (Wegfall der Egress-Restriktion) und [`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Ausnahme 2) und berührt [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) (kein Pfad vom Container zum Host).

**Statusdateien.** Jedes Startscript schreibt bei jedem Lauf, auch wenn seine Fähigkeit verfügbar ist, eine Statusdatei im Container unter `/var/lib/u-boot/status/`. Das erzeugte Dockerfile legt das Verzeichnis für den Container-Benutzer an, nur wenn mindestens ein Startscript installiert ist.

| Datei | Startscript | Fähigkeit | Läuft bei |
|---|---|---|---|
| `vfs` | `sandbox-init.sh` | `nestedRuntime: podman` | `postCreateCommand`, einmal je Container |
| `daemon` | `docker-init.sh` | `nestedRuntime: docker` | `postStartCommand`, bei jedem Start |
| `egress` | `egress-init.sh` | Egress-Restriktion | `postStartCommand`, bei jedem Start |

Das Dokument hat genau drei Zeilen:

```
u-boot-status 1
result=ok|fallback|failed
at=2026-10-17T18:52:54Z
```

`result` ist `fallback` bei der Richtlinie `warn` und eingetretenem Fallback (`vfs`-Storage, Container ohne Daemon, Wegfall der Egress-Restriktion), `failed` bei der Richtlinie `fail` (das Startscript schreibt die Datei vor dem Exit-Code `11`), sonst `ok`. `at` ist der Zeitpunkt des Scriptlaufs in UTC, Sekundengenauigkeit, Form `YYYY-MM-DDTHH:MM:SSZ`. Die Datei trägt keinen Freitext und keine Container-Kennung. Das Startscript schreibt in eine Temporärdatei und benennt sie um, sodass kein halbes Dokument lesbar ist. Blockierte verschachtelte User-Namespaces (Exit-Code `11` in beiden Richtlinien) sind ein Umgebungsproblem und kein Fallback; es entsteht keine Datei. Exit-Codes und Standardfehlerausgabe der Startscripte bleiben unverändert.

**Strenge beim Lesen.** Ein Dokument ist gültig, wenn es höchstens 512 Byte umfasst, nur druckbare ASCII-Zeichen und Zeilenwechsel enthält, genau die drei Zeilen oben trägt (ein abschließender Zeilenwechsel ist zulässig), `result` einen der drei Werte hat und `at` genau die Form oben hat (keine Bruchteile, kein Offset). Jede Abweichung, auch Freitext, macht es ungültig. Ausgegeben werden nur die aufgelösten Felder, nie der Inhalt der Datei.

**Ermittlung durch den Host.** `doctor`, `generate devcontainer` und `init` lesen die Dateien über die Container-Engine, ohne ein Programm im Container auszuführen:

- Der Container ist der, der das Label `devcontainer.local_folder` mit dem absoluten Projektordner (dem Basisverzeichnis des Befehls) trägt; das setzt das Devcontainer-Werkzeug (gemessen mit der Devcontainer-CLI). Bei mehreren Treffern gilt der zuletzt erzeugte, auch ein gestoppter. `u-boot` setzt kein eigenes Label und fügt keine `runArgs` hinzu.
- Gelesen wird mit `docker cp <Container>:<Pfad> -`, ohne Link-Auflösung und ohne `docker exec`. Der Tar-Kopf entscheidet vor dem Lesen: nur eine reguläre Datei von höchstens 512 Byte; Symlink, FIFO, Verzeichnis, Gerät und Übergröße machen das Dokument ungültig. Die Container-Kennung wird vor jeder Verwendung gegen `[0-9a-f]{12,64}` geprüft, der Projektordner als ein Argument ohne Shell übergeben, jeder Aufruf der Engine ist auf zehn Sekunden begrenzt.
- Ein entfernter Engine-Zugang (`DOCKER_HOST` oder aktiver Docker-Kontext zeigt nicht auf einen lokalen Unix-Socket) gilt als nicht ermittelbar: Ein Container dort, der denselben Ordnerwert trägt, gehört zu einer anderen Maschine.
- Läuft `u-boot` selbst im Devcontainer (Container-Erkennung), liest `doctor` die Dateien direkt aus dem Dateisystem, nach derselben Typ- und Größenprüfung; fehlt das Verzeichnis (etwa `u-boot` aus dem Runtime-Image), ist der Zustand nicht ermittelbar. `generate devcontainer` und `init` fragen dort nicht.

**Nicht ermittelbar** heißt: kein Container zum Projektordner, keine oder entfernte Engine, keine Datei (ein Container, der vor den Statusdateien erzeugt wurde oder noch nie gestartet ist, hat keine; er braucht `u-boot generate devcontainer` und einen Neubau), ein Lesefehler. Das ist nie ein Fallback.

**Stand.** `postStartCommand` überschreibt `daemon` und `egress` bei jedem Start; `vfs` gilt bis zum Entfernen des Containers wie die `storage.conf`, die das Script schreibt. Ein gestoppter Container wird gelesen; der Eintrag nennt den Zeitstempel und den Vermerk „container stopped".

**Abbildung in `u-boot doctor`.** Je konfigurierter Fähigkeit ein Eintrag, und nur dann: `devcontainer.sandbox.fallback.vfs` (`nestedRuntime: podman`), `….daemon` (`nestedRuntime: docker`), `….egress` (Egress-Restriktion), jeweils mit Sandbox-Profil. Ohne Sandbox-Profil oder ohne die Fähigkeit entsteht kein Eintrag.

| Zustand | Stufe | Kennzeichen `DegradationFallback` (sichtbar unter `--quiet`) |
|---|---|---|
| `result=ok` | OK | nein |
| `result=fallback` | `warn` | ja |
| `result=failed` | `error` (Exit-Code `11`) | nein |
| nicht ermittelbar | `info` mit Grund | nein |
| ungültiges Dokument | `warn` | nein |

`doctor --strict` endet bei `warn` mit Exit-Code `11`, wie bei jedem `warn`-Eintrag. Ein Fehler der Engine verschiebt den Exit-Code nicht: der Bericht entsteht immer, der Fehler ist „nicht ermittelbar". `diagnostics[]` der JSON-Ausgabe trägt wie sonst nur `warn` und `error`; ein `info`-Eintrag erscheint dort nicht.

**Abbildung in `generate devcontainer` und `init`.** Nach dem Schreiben fragen beide für jede Fähigkeit, die `u-boot.yaml` (bei `init` die geschriebene) konfiguriert, denselben Weg ab (ein Container-Lookup je Lauf). Meldet ein bestehender Container `result=fallback`, steht in der Warnungsausgabe eine Warnung (Level `warn`, Subjekt `.devcontainer/devcontainer.json`, Code [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) für `vfs` und Daemon, [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) für Egress) mit dem Kennzeichen `DegradationFallback`, ohne `SecurityRelaxation`; sie bleibt unter `--quiet` sichtbar ([`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 2). Der Text nennt, dass es der Stand des letzten Starts des bestehenden Containers ist und die neu erzeugten Dateien erst nach einem Neubau wirken. Die Abfrage ist best effort: ohne Engine, ohne Container, ohne Datei, bei ungültigem Dokument und bei `ok` oder `failed` bleibt sie still (keine Meldung, keine Änderung des Exit-Codes); mit `--dry-run` und ohne konfigurierte Fähigkeit findet sie nicht statt. `up` und jeder andere Befehl fragen nicht.

**Grenze.** Die Datei stammt aus einem Container, in dem ein Agent läuft, und ist dort mit Root-Rechten frei schreibbar. Sie ist ein Hinweis, kein Beleg: „kein Fallback gemeldet" heißt nicht „kein Fallback aufgetreten". Das Sandbox-Profil bleibt Schadensbegrenzung, keine Grenze ([`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil)).

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
#     nestedRuntime: none          # none | podman | docker
#     onUnavailable: warn          # warn | fail
#     repository: ''               # Clone-Quelle statt origin (URL)
#     egress:
#       enabled: false
#       allow: []
```

`u-boot init --devcontainer` schreibt `devcontainer.enabled: true` in die neu angelegte `u-boot.yaml` ([`LH-AK-005`](lastenheft.md#lh-ak-005--devcontainer-flow)).

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
│   ├── hexagon/                 # innere Schichten (SPEC-012)
│   │   ├── domain/
│   │   ├── application/
│   │   └── port/{driving,driven}/
│   └── adapter/                 # äußere Schichten (SPEC-012)
│       ├── driving/
│       └── driven/
├── spec/                        # Lastenheft, weitere Spezifikationen
│   ├── lastenheft.md
│   └── <weitere-spezifikation>.md
├── docs/                        # Doku-Struktur (SPEC-010)
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

### SPEC-024 — Beispielaufrufe zu den Anforderungen

Beispielaufrufe zu den Anforderungen (zur Veranschaulichung; verbindlich ist die jeweilige Anforderung).

- [LH-FA-CLI-001](lastenheft.md#lh-fa-cli-001--cli-aufruf): `u-boot --help`
- [LH-FA-CLI-003](lastenheft.md#lh-fa-cli-003--versionsausgabe): `u-boot --version`
- [LH-FA-CLI-007](lastenheft.md#lh-fa-cli-007--dry-run): `u-boot add postgres --dry-run`
- [LH-FA-CLI-008](lastenheft.md#lh-fa-cli-008--diff-ausgabe): `u-boot add postgres --diff`
- [LH-FA-INIT-001](lastenheft.md#lh-fa-init-001--neues-projekt-initialisieren): `u-boot init`, `u-boot init --assume-existing`
- [LH-FA-INIT-002](lastenheft.md#lh-fa-init-002--projektname): `u-boot init`, `u-boot init my-service`
- [LH-FA-DEV-001](lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen): `u-boot init --devcontainer`, `u-boot generate devcontainer`
- [LH-FA-DEV-006](lastenheft.md#lh-fa-dev-006--sandbox-profil): `u-boot init --devcontainer --sandbox`, `u-boot generate devcontainer --sandbox`
- [LH-FA-ADD-001](lastenheft.md#lh-fa-add-001--add-on-befehl): `u-boot add <service>`
- [LH-FA-ADD-002](lastenheft.md#lh-fa-add-002--postgresql-hinzufügen): `u-boot add postgres`
- [LH-FA-ADD-003](lastenheft.md#lh-fa-add-003--keycloak-hinzufügen): `u-boot add keycloak`
- [LH-FA-ADD-004](lastenheft.md#lh-fa-add-004--opentelemetry-hinzufügen): `u-boot add otel`
- [LH-FA-ADD-007](lastenheft.md#lh-fa-add-007--service-entfernen): `u-boot remove postgres`
- [LH-FA-UP-001](lastenheft.md#lh-fa-up-001--umgebung-starten): `u-boot up`
- [LH-FA-UP-004](lastenheft.md#lh-fa-up-004--umgebung-stoppen): `u-boot down`; `u-boot down --volumes`
- [LH-FA-UP-005](lastenheft.md#lh-fa-up-005--logs-anzeigen): `u-boot logs`, `u-boot logs postgres`
- [LH-FA-DIAG-001](lastenheft.md#lh-fa-diag-001--doctor-befehl): `u-boot doctor`
- [LH-FA-GEN-001](lastenheft.md#lh-fa-gen-001--generate-befehl): `u-boot generate <artifact>`
- [LH-FA-GEN-002](lastenheft.md#lh-fa-gen-002--changelog-erzeugen): `u-boot generate changelog`
- [LH-FA-GEN-003](lastenheft.md#lh-fa-gen-003--readme-erzeugen): `u-boot generate readme`
- [LH-FA-GEN-004](lastenheft.md#lh-fa-gen-004--beispiel-env-erzeugen): `u-boot generate env-example`
- [LH-FA-TPL-001](lastenheft.md#lh-fa-tpl-001--projektvorlagen): `u-boot init --template basic`, `u-boot init --template micronaut`, `u-boot init --template sveltekit`, `u-boot init --template micronaut-sveltekit`
- [LH-FA-TPL-003](lastenheft.md#lh-fa-tpl-003--eigene-templates): `u-boot init --template ./my-template`
- [LH-FA-TPL-004](lastenheft.md#lh-fa-tpl-004--templates-auflisten): `u-boot template list`
- [LH-FA-CONF-001](lastenheft.md#lh-fa-conf-001--projektkonfiguration): `u-boot config get project.name`, `u-boot config set project.name my-service`
- [LH-FA-CONF-005](lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern): `u-boot config`, `u-boot config get project.name`, `u-boot config set project.name foo`
- [LH-FA-CONF-006](lastenheft.md#lh-fa-conf-006--konfiguration-migrieren): `u-boot config migrate`
- [LH-NFA-USE-002](lastenheft.md#lh-nfa-use-002--klare-befehle): `u-boot init`, `u-boot add postgres`, `u-boot doctor`, `u-boot up`
- [LH-NFA-USE-004](lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe): `u-boot doctor --json`
- [LH-SA-CLI-001](lastenheft.md#lh-sa-cli-001--befehlsstruktur): `u-boot <command> [subcommand|args...] [options]`

## 3. Defaults und Konstanten

Werte, die im Produkt fest sind (Standardwerte, Grenzwerte, Versionsuntergrenzen).

### SPEC-006 — Build-Args des Dockerfiles und Pin-Politik

Gilt für das Dockerfile im Repo-Root; Overrides per `docker build --build-arg <NAME>=<value>` bzw. Makefile-Variable.

- `ARG GO_VERSION` – mit Default-Pin (Stand 2026-10-14: `1.27.2`); Hebung ist Routine ohne separaten Spec-Eintrag.
- `ARG GOLANGCI_LINT_VERSION` – mit Default-Pin; gleiche Pin-Politik.
- `ARG COVERAGE_THRESHOLD` – mit Default `90` (Prozent) und Override-Pfad `make coverage-gate THRESHOLD=…`; das Bootstrap-Verhalten bei leerer Coverage-Eingabe steht in der Anforderung.

### SPEC-007 — Mindest-Ausschlüsse der `.dockerignore`

Mindestens auszuschließen:

- `.git`, `.gitignore`, `.github`
- IDE-Verzeichnisse: `.idea`, `.vscode`
- Agent-Verzeichnisse: `.claude`, `.codex`, `.agents`
- lokale Build-Artefakte und Caches (z. B. `dist/`, `coverage*`, `*.log`)

### SPEC-016 — Go-Toolchain: Mindestversion und Dockerfile-Pin

Mindest-Toolchain: Go 1.26 oder neuer (`go 1.26.0` in `go.mod`, analog Referenzprojekt `k-deskflight`); Default-Pin im Dockerfile als `ARG GO_VERSION` (am Entscheidungsdatum `1.26.3`, Stand 2026-10-14 `1.27.2`). Pin-Hebung ist Routine ohne separaten Spec-Eintrag.

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
| `nestedRuntime: podman` oder `docker` | `registry-1.docker.io`, `auth.docker.io`, `production.cloudflare.docker.com`, `production.cloudfront.docker.com`, `ghcr.io` |
| `nestedRuntime: docker` | zusätzlich `download.docker.com` (Paketquelle bleibt im Image eingetragen) |
| Feature `node` | `registry.npmjs.org` |
| Feature `go` | `proxy.golang.org`, `sum.golang.org`, `storage.googleapis.com` |
| Feature `java` | `repo.maven.apache.org`, `repo1.maven.org` |
| Feature `python` | `pypi.org`, `files.pythonhosted.org` |
| Feature `sdkman` | `api.sdkman.io`, `broker.sdkman.io`, `repo.maven.apache.org`, `repo1.maven.org` |
| Feature `homebrew` | `formulae.brew.sh`, `ghcr.io`, `pkg-containers.githubusercontent.com` |
| Feature `claude-code` | `registry.npmjs.org`, `api.anthropic.com` |
| Clone-Quelle | Host von `devcontainer.sandbox.repository` beziehungsweise `origin` |
| Nutzer | Einträge von `devcontainer.sandbox.egress.allow` (kleingeschriebene Hostnamen, ohne Schema, Port und Wildcard) |

Das Feature `sshd` ergänzt keinen Host. Die Zeile `claude-code` führt den Host der Anthropic-API und die npm-Registry (die ohnehin zum Feature `node` gehört, das `claude-code` braucht); die Hosts der Anmeldung sowie von Telemetrie, Updates, MCP und Plugins führt die Liste nicht, der Nutzer trägt sie über `devcontainer.sandbox.egress.allow` ein.

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
| `devcontainer.sandbox.runtime` | Sandbox: nested Podman (`/dev/fuse` auf dem Linux-Host: Fehler bei Fehlen, Profil-Konsistenz), [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) |
| `devcontainer.sandbox.egress` | Sandbox: Egress-Restriktion konsistent konfiguriert, [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) |
| `devcontainer.sandbox.credentials` | Sandbox: keine Klartext-Zugangsdaten (Git und Agent), Token-Quelle bei https-Clone, Referenz `${localEnv:<NAME>}` bei aktivem Agenten-Feature, [`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), Einzelheiten [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) |
| `devcontainer.sandbox.fallback.vfs` | Sandbox: Zustand des `vfs`-Storage-Fallbacks beim Containerstart laut Statusdatei, nur bei `nestedRuntime: podman`, [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Einzelheiten [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) |
| `devcontainer.sandbox.fallback.daemon` | Sandbox: Zustand des Docker-Daemons beim Containerstart laut Statusdatei, nur bei `nestedRuntime: docker`, [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Einzelheiten [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) |
| `devcontainer.sandbox.fallback.egress` | Sandbox: Zustand der Egress-Restriktion beim Containerstart laut Statusdatei, nur mit Egress-Restriktion, [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion), Einzelheiten [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) |

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
- Eingebaute Features (Schlüssel → Quelle unter `ghcr.io/devcontainers/features/`, jeweils mit der Standardversion `1`): `git` → `git`, `docker-cli` → `docker-outside-of-docker`, `node` → `node`, `java` → `java`, `go` → `go`, `cpp` → `cpp`, `kubectl-helm` → `kubectl-helm-minikube`, `postgres-client` → `postgresql-client`, `python` → `python`, `sshd` → `sshd` (installiert einen SSH-Server, also einen Login-Dienst; Warnung nach [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)) und `sdkman` → `java` (das Feature installiert SDKMAN; teilt sich den Render-Schlüssel mit `java`, der einmal erscheint). Zwei Schlüssel stammen von einem Drittanbieter (Quelle außerhalb von `ghcr.io/devcontainers/features/`): `homebrew` → `ghcr.io/meaningful-ooo/devcontainer-features/homebrew` mit der Standardversion `2` (Herausgeber `meaningful-ooo`) und `claude-code` → `ghcr.io/anthropics/devcontainer-features/claude-code` mit der Standardversion `1` (Herausgeber `anthropics`; das Feature installiert über npm und braucht das Feature `node`). Für beide weist die Befehlsausgabe den Herausgeber aus, solange der Eintrag ohne eigene `source` aktiviert ist ([`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)); alle Katalogschlüssel sind ohne Eintrag in `devcontainer.featureSources.allow` aktivierbar. Jede andere Quelle gilt als extern und braucht eine Freigabe in `devcontainer.featureSources.allow`.
- Externe Feature-Quellen müssen als URL mit Schema `http://`, `https://` oder `oci://` angegeben werden.

## 7. Historie

| Datum | Änderung |
|---|---|
| 2026-10-02 | Angelegt als Technik-Stratum; Schemata, Beispielinstanzen, Algorithmen, Build-/CI-, Doku- und Architektur-Details sowie Markierungsformate wörtlich aus dem Lastenheft übernommen (`SPEC-001`..`SPEC-016`, Verfeinerungen); neu aus dem Bestand von Code und Konfiguration: Standardwerte, Egress-Default-Allowlist, Code-Registry der `doctor`-Prüfungen, Telemetrie-Aussage, externe Verträge (`SPEC-017`..`SPEC-023`). |
| 2026-10-16 | Festlegungen aus den Paket-READMEs unter `internal/` übernommen (`SPEC-036`, [`LH-FA-UP-001.a`](#lh-fa-up-001a--erkennung-probebarer-ports-in-composeyaml), `SPEC-038`..`SPEC-042`): Stabilisierung und Port-Erkennung von `up`, Bestätigungspfad von `down --volumes`, Zustandsmaschinen von `generate`, Pfade und Schreibpipeline von `config`, Schweregrade von `compose.yaml.valid`, Preflight der Compose-Engine. Zahlenwerte aus dem Quelltext. |
| 2026-10-16 | Verbosity festgelegt (`SPEC-043`): Pegelabbildung, Rangfolge „zuletzt angegeben gewinnt“ und Wirkung von `--quiet`, deckungsgleich mit Lastenheft 0.10.0 und Code. |
| 2026-10-17 | Features `sshd` und `claude-code` festgelegt, deckungsgleich mit Lastenheft 0.11.0 und Code: Katalogschlüssel (`SPEC-023`), Egress-Hosts von `claude-code` (`SPEC-019`), Beschreibung der Prüfung `devcontainer.sandbox.credentials` (`SPEC-017`), Warnungen und Zugangsdaten-Übergabe (`SPEC-044`); `SPEC-043` nennt in Ausnahme 1 den Login-Dienst und die `doctor`-Sicherheits-Einträge und hält die Einstufung je Eintrag fest. |
| 2026-10-17 | `SPEC-043` führt die Meldung, dass `egress.enabled` ohne Sandbox-Profil keine Wirkung hat, in der Befehlsausgabe von `generate devcontainer` unter den gekennzeichneten Einträgen; ungekennzeichnet bleibt nur die Meldung „has no effect“ zu `nestedRuntime`. |
| 2026-10-17 | [`LH-FA-DEV-007.b`](#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman) beschreibt bei `nestedRuntime: podman` die Ausstattung für Builds, gelesen gegen Code und Dockerfile: sechs `runArgs` (neu `--device=/dev/net/tun`), Pakete, `docker`-Symlink, kein `buildah`-Kommando, Kurznamen, `vfs`-Fallback, Belege per Integrationstest. [`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion) führt bei `podman` und `docker` zusätzlich `production.cloudfront.docker.com` (Blob-Host von Docker Hub). |
| 2026-10-18 | [`LH-FA-DEV-007.a`](#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) und [`.b`](#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman) folgen dem Lastenheft 0.13.0: ein fehlendes `/dev/fuse` auf dem Container-Host ist ein Startfehler der Container-Engine und keine Degradationsstufe; der `vfs`-Fallback gilt für einen ohne das Device gestarteten Container; `u-boot doctor` weist das fehlende Device auf dem Linux-Host als Fehler aus (Beschreibung von `devcontainer.sandbox.runtime`). |
| 2026-10-18 | Ausweisung der Degradations-Fallbacks festgelegt, deckungsgleich mit Lastenheft 0.13.0 und Code: Statusdateien der Startscripte, Ermittlung durch den Host und Abbildung in `doctor`, `generate devcontainer` und `init` (`SPEC-045`); [`LH-FA-DEV-007.a`](#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten), [`LH-FA-DEV-008.a`](#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability) und `SPEC-043` Ausnahme 2 tragen die Aussage des Lastenhefts ohne die Einschränkung „soweit vom Host aus ermittelbar“, mit der Grenze „nicht ermittelbar“; drei neue Prüfungen in der Code-Registry (`SPEC-017`). |
