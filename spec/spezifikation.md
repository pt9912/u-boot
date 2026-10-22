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
| Startscript der Egress-Restriktion bricht aus einem anderen Grund als der fehlenden Capability ab (Benutzer des Resolvers fehlt, nftables-Regeln nicht anwendbar, Resolver startet nicht oder läuft unter anderer UID, [`SPEC-048`](#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf)) | Umgebungsproblem, Exit-Code `11` | Umgebungsproblem, Exit-Code `11` |

Ein fehlendes `/dev/fuse` oder `/dev/net/tun` auf dem Container-Host ist keine Degradationsstufe und steht nicht in der Tabelle: die Container-Engine lehnt den Start ab, das Startscript läuft nie. `u-boot doctor` prüft `/dev/net/tun` nicht; ein fehlendes `/dev/fuse` weist es auf einem Linux-Host, soweit der Host es bestimmen kann, als Fehler aus (Umgebungsproblem, Exit-Code `11`, unabhängig von `onUnavailable` und ohne Kennzeichen eines Fallbacks). Läuft `u-boot` in einem Container, prüft `doctor` das Device nicht (Info, kein Beitrag zum Exit-Code), weil das `/dev/fuse` des Containers nicht das des Container-Hosts ist; ein ohne das Device gestarteter Devcontainer erscheint im Eintrag `devcontainer.sandbox.fallback.vfs`, soweit dessen Statusdatei lesbar ist (sonst „nicht ermittelbar“, [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)).

Ausdrücklich angeforderte Runtime (`nestedRuntime: podman` oder `docker`) wird nie stillschweigend durch etwas anderes ersetzt. Jeder Fallback wird in der Befehlsausgabe und in `u-boot doctor` ausgewiesen: Das erzeugte Startscript schreibt bei jedem Lauf eine Statusdatei in den Container (`vfs`, `daemon`, `egress`) und endet im Fehlerfall nicht-null; `u-boot doctor` liest die Datei vom Host über die Container-Engine und weist je Fallback einen Eintrag aus (`devcontainer.sandbox.fallback.vfs`, `….daemon`, `….egress`), `u-boot generate devcontainer` weist nach dem Schreiben den Fallback eines bestehenden Containers aus; `u-boot init` fragt nach dem Schreiben auf demselben Weg, die von ihm geschriebene Konfiguration enthält jedoch weder `nestedRuntime` noch Egress-Restriktion, sodass dort keine Fähigkeit abzufragen ist; `u-boot devcontainer up` und `u-boot devcontainer shell` weisen den Fallback des gestarteten Containers nach dem Start aus und bilden den Abbruch eines Startscripts (`result=failed`) auf Exit-Code `11` ab, wenn eine `u-boot.yaml` mit Sandbox-Profil die Fähigkeit konfiguriert und die Statusdatei lesbar ist (nach einem gelungenen Lauf der CLI wie nach einem Fehler, auch bei einem `failed` aus einem früheren Start des Containers), sonst bleibt es beim Exit-Code des Starts (`0` oder `12`) ([`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)). Format, Ermittlung und Grenzen: [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung). Der Host meldet einen Fallback nur, wenn ein Container zum Projektordner gefunden und seine Statusdatei gelesen wird; sonst ist der Zustand „nicht ermittelbar“, was kein Fallback ist und auch nicht ausschließt, dass einer vorliegt. Ungültige Werte der Schlüssel führen zu einem fachlichen Validierungsfehler (Exit-Code `10`).

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
- Bei `command == "template"`, `command == "config"`, `command == "egress"` oder `command == "devcontainer"` ist `subcommand` verpflichtend.
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

- Bei `podman`: rootless Podman, ein `docker`-Kompatibilitäts-Alias, Zugriff auf `/dev/fuse` und `/dev/net/tun` (Netzwerk der verschachtelten Container über `pasta`; ohne das Device scheitert jedes `RUN`; beide Devices müssen auf dem Container-Host vorhanden sein, fehlt eines, lehnt die Container-Engine den Start ab (`error gathering device information while adding custom device`), das Startscript läuft dann nicht; `doctor` prüft `/dev/net/tun` nicht und weist ein fehlendes `/dev/fuse` nur aus, soweit der Host es bestimmen kann, nicht aus einem Container), subuid-/subgid-Bereiche für den Container-Benutzer und ein Volume für den Storage. `runArgs` enthält sechs Optionen: `--cap-add=SYS_ADMIN`, `--security-opt=seccomp=unconfined`, `--security-opt=apparmor=unconfined`, `--security-opt=systempaths=unconfined`, `--device=/dev/fuse`, `--device=/dev/net/tun`; jede weist die Befehlsausgabe einzeln aus. Podman kommt nicht aus der Distribution, sondern aus dem statisch gelinkten Release-Archiv `podman-linux-<arch>.tar.gz` des Projekts `podman-static` (Version 6.1.4; `amd64` und `arm64`; enthält `podman`, `crun`, `runc`, `conmon`, `netavark`, `aardvark-dns`, `passt`/`pasta`, `fuse-overlayfs` und `catatonit` unter `/usr/local`). Version, SHA-256 je Architektur und der Primär-Fingerprint des Signaturschlüssels (`0CCF102C4F95D89E583FF1D4F8B5AF50344BB503`) stehen als `ARG` an einer Stelle im erzeugten Dockerfile; der Image-Build bricht vor dem Entpacken ab, wenn die Prüfsumme nicht stimmt, der Schlüssel (von `keys.openpgp.org`, nach Fingerprint) nicht importierbar ist oder die Signatur des Archivs nicht vom gepinnten Primärschlüssel stammt (`VALIDSIG`). Das Archiv wird entpackt, nicht ausgeführt; Archiv und Schlüssel lädt der Image-Build von `github.com` und `objects.githubusercontent.com` (Archiv, Signatur) und `keys.openpgp.org` (Schlüssel); der Build läuft auf dem Container-Host ohne die Egress-Restriktion, braucht diese Hosts also offen, und im laufenden Container ist keiner davon nötig ([`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion) bleibt unverändert). Aus Debian installiert das Dockerfile nur `uidmap` und `ca-certificates` ohne empfohlene Pakete (`podman`, `fuse-overlayfs` und `passt` kommen nicht aus Debian); `nftables` kommt nur mit der Egress-Restriktion ins Image. Der Alias ist ein Symlink `/usr/local/bin/docker` → `/usr/local/bin/podman`. Das Archiv liefert unter `/etc/containers/` `containers.conf`, `policy.json`, `registries.conf`, `seccomp.json`, `storage.conf` und `storage.conf.d/00-podman-static.conf`; sie gelten mit einer Änderung, und das `README.md` des Archivs wird nach dem Entpacken entfernt. `registries.conf` setzt `unqualified-search-registries = ['docker.io']`. Daraus folgt für Builds: `podman build` und `docker build` sind derselbe Aufruf mit der in Podman eingebauten Buildah-Bibliothek; das Kommando `buildah`, `docker buildx`, `docker compose` und BuildKit-Frontends (`# syntax=`) sind nicht Teil der Ausstattung. Ein Dockerfile mit `FROM`, `RUN` und `COPY` baut kalt, aus dem Layer-Cache und mit `--no-cache`, mit und ohne Egress-Restriktion (Default-Allowlist). Der Kurzname `alpine` löst über Docker Hub auf (`docker.io/library/alpine`), `mycompany/img` wird nur dort gesucht; Dockerfiles verwenden voll qualifizierte Image-Namen. Fehlt `/dev/fuse` im Container (der Container wurde ohne das Device gestartet), schreibt das Startscript ein Benutzer-Drop-in (`~/.config/containers/storage.conf.d/99-u-boot-vfs.conf` mit `driver = "vfs"` und leerem `mount_program`, das die Vorgabe der mitgelieferten Konfiguration überstimmt) und setzt `image_copy_tmp_dir` (`~/.config/containers/containers.conf`) auf ein Verzeichnis im Storage-Volume; `podman run`, `podman build` und `docker build` laufen dann mit `vfs`-Storage. Der Beleg sind Integrationstests; sie laufen nicht auf Hosts, die verschachtelte User-Namespaces blockieren (`kernel.apparmor_restrict_unprivileged_userns=1`, etwa Ubuntu 24.04 mit Standardeinstellung), dort endet das Startscript mit Exit-Code `11`.
- Bei `docker`: Docker-Daemon, `docker`-CLI sowie Buildx- und Compose-Plugin aus dem Docker-APT-Repository (Schlüssel per Fingerprint geprüft); der Container-Benutzer gehört der Gruppe `docker` an, der Daemon startet beim Containerstart (`postStartCommand`) und gilt nach spätestens 60 Sekunden als nicht betriebsbereit; `runArgs` enthält `--privileged` und, nur mit Egress-Restriktion, `--cap-add=NET_ADMIN`, sonst keine Fähigkeitsoptionen; der Daemon-Speicher (`/var/lib/docker`) liegt in einem Volume `<projekt>-docker-${devcontainerId}`; kein Socket-Mount.

### LH-FA-DEV-008.a — Verhalten bei nicht gewährbarer Capability

Verfeinert [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion).

- Ist die nötige Capability nicht gewährbar, greift die Degradationstabelle der Verfeinerung [`LH-FA-DEV-007.a`](#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten); die Capability wird beim Containerstart geprüft (Warnung und Wegfall der Restriktion, bei `onUnavailable: fail` Exit-Code `11`), und das Startscript schreibt das Ergebnis in die Statusdatei `egress`. Für die Abbrüche des Startscripts aus anderem Grund als der fehlenden Capability gilt nicht die Zeile der Degradationstabelle, sondern [`SPEC-048`](#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf) (Exit-Code `11` unter jeder Richtlinie). `u-boot doctor` weist den Wegfall der Restriktion vom Host aus (Eintrag `devcontainer.sandbox.fallback.egress`, [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)) und prüft zusätzlich die Konfiguration (Schlüssel ohne Sandbox-Profil: `warn`). Die Capability selbst bestimmt der Host nicht, er liest nur das gemeldete Ergebnis; ohne lesbare Statusdatei ist der Zustand nicht ermittelbar.

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

Verfeinert [`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Stufen, Wirkung von `--quiet`, Mehrfachangabe) und berührt [`LH-FA-CLI-007`](lastenheft.md#lh-fa-cli-007--dry-run) (Vorschau), [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) und [`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Ausweisungen), [`LH-NFA-PERF-002`](lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker) (Fortschritt von `up`) sowie [`LH-FA-INIT-005`](lastenheft.md#lh-fa-init-005--überschreibschutz) und [`LH-FA-CLI-005A`](lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (Pfad-Zusammenfassung bei `--force`).

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

- Die ändernden Subkommandos (`init`, `generate`, `add`, `remove`, `config set`, `up`, `down`, `devcontainer up`, `devcontainer down`, `devcontainer shell`) geben nur Fehler aus; ihre Erfolgs- und Zusammenfassungsmeldungen entfallen (`devcontainer shell` gibt keine Erfolgszeile aus), bei `remove` auch die Purge-Warnungen.
- `doctor` schreibt nur Fehler-Einträge, Einträge eines Degradations-Fallbacks und die Warn-Einträge, die eine Sicherheitseigenschaft des Sandbox-Profils betreffen (Ausnahme 1), ohne Kopf- und Summary-Zeile; der Exit-Code trägt das Ergebnis.
- Die Abfragebefehle `config get`, `config list`, `egress list`, `logs` und `template list` behalten ihre Ausgabe; bei `egress list` bleibt die gesamte Ausgabe, auch der Hinweis, dass die Restriktion nicht aktiv ist ([`SPEC-046`](#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung)).
- Die Ausgabe von `--dry-run` und `--diff` (Änderungsplan, Diff) bleibt ([`LH-FA-CLI-007`](lastenheft.md#lh-fa-cli-007--dry-run)); die Erfolgsmeldung eines tatsächlichen Schreibens entfällt auch bei `--diff` ohne `--dry-run`.
- Vier Ausnahmen bleiben sichtbar: (1) die Ausweisung von Sicherheitslockerungen des Sandbox-Profils und eines im Container aktivierten Login-Dienstes ([`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion)) in der Warnungsausgabe von `generate devcontainer` (die Lockerungen des Sandbox-Profils und den Login-Dienst weist nur dieser Befehl aus; zu ihnen gehören die Lockerungen eines aktiven eingebauten Features und die Warnung je Feature aus einer externen Quelle, [`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen)), die Ausweisung der Einbindung des Container-Runtime-Sockets in den Hilfs-Container in der Ausgabe von `devcontainer up` und `devcontainer shell` (bei `shell` vor dem Beginn der Sitzung auf der Standardfehlerausgabe; [`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)) sowie die Warn-Einträge von `doctor`, die eine Sicherheitseigenschaft des Sandbox-Profils betreffen: die Quelle der Zugangsdaten ([`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)) und die Egress-Restriktion ([`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion)), darunter der Eintrag, dass sie ohne Sandbox-Profil nicht greift; (2) die Ausweisung der Degradations-Fallbacks (Container ohne Daemon, `vfs`-Storage, Wegfall der Egress-Restriktion) in `doctor` (Einträge `devcontainer.sandbox.fallback.*`) und in der Warnungsausgabe von `generate devcontainer`, `init`, `devcontainer up` und `devcontainer shell` (bei `shell` vor dem Beginn der Sitzung auf der Standardfehlerausgabe) nach [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung); ein Zustand, den der Host nicht lesen kann, ist „nicht ermittelbar“ und kein Fallback, und sein Eintrag entfällt unter `--quiet`; (3) der Compose-Fortschritt von `up` ([`LH-NFA-PERF-002`](lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker)); (4) die Pfad-Zusammenfassung, die `init --force` vor dem Überschreiben ausgibt ([`LH-FA-CLI-005A`](lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-INIT-005`](lastenheft.md#lh-fa-init-005--überschreibschutz)). Jede andere Warnung und jeder andere Hinweis entfällt.
- Die Einstufung hängt am Kennzeichen des einzelnen Eintrags (`SecurityRelaxation`, `DegradationFallback`), nicht an der Kennung der Prüfung oder Anforderung; dieselbe Kennung trägt gekennzeichnete und ungekennzeichnete Einträge. Im Einzelnen:
  - Gekennzeichnet sind in der Warnungsausgabe von `generate devcontainer` die Lockerungen von `nestedRuntime` und Egress-Restriktion (je Option eine Meldung), die Meldung, dass `egress.enabled` ohne Sandbox-Profil keine Wirkung hat (es entsteht keine Restriktion), der Login-Dienst, die Lockerungen eines aktiven eingebauten Features (`go`: je Option eine Meldung) und die Warnung je Feature aus einer externen Quelle ([`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen)); in `doctor` die Warn-Einträge von `devcontainer.sandbox.credentials` (Klartext-Zugangsdaten, https-Clone ohne Token-Quelle, fehlende Referenz der Agenten-Variable) und von `devcontainer.sandbox.egress` (ohne Sandbox-Profil, zusammen mit `nestedRuntime: docker`); in der Ausgabe von `devcontainer up` und `devcontainer shell` die Ausweisung der Einbindung des Container-Runtime-Sockets in den Hilfs-Container ([`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)).
  - Mit `DegradationFallback` gekennzeichnet sind der `warn`-Eintrag je `devcontainer.sandbox.fallback.*` bei `result=fallback` und die entsprechende Warnung von `generate devcontainer`, `init`, `devcontainer up` und `devcontainer shell`; ungekennzeichnet bleiben die Einträge dieser Prüfungen bei `ok`, „nicht ermittelbar“ (`info`) und ungültigem Dokument (`warn`), der Eintrag bei `failed` ist ein Fehler und bleibt unter `--quiet` als Fehler sichtbar.
  - `nestedRuntime: docker` hat in `doctor` die Stufe OK (kein Warn-Eintrag; `--privileged` weist die Befehlsausgabe aus): Die Sicherheitslockerung des Docker-Daemons erscheint nur in der Befehlsausgabe.
  - Ungekennzeichnet sind die Meldungen „has no effect“ zu `nestedRuntime` in der Warnungsausgabe von `generate devcontainer` (Konfigurationsmeldung), der `doctor`-Eintrag, dass `nestedRuntime` ohne Sandbox-Profil keine Wirkung hat, die Drittanbieter-Warnung und der Hinweis auf das Feature `node` ([`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)), die beiden Meldungen zum Feature `ollama` (Lieferkette; keine Wirkung ohne Sandbox-Profil, [`SPEC-047`](#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen)), der Hinweis zum Neuaufbau bei geändertem `egress-init.sh` ([`SPEC-046`](#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung)) sowie die Warn-Einträge von `devcontainer.features.allowlist` und `devcontainer.features.drift`; sie entfallen. Dieselbe Aussage zur Egress-Restriktion ohne Sandbox-Profil ist in `doctor` wie in der Befehlsausgabe gekennzeichnet, also sichtbar.
- `down --quiet` hält die Compose-Ausgabe zurück und gibt sie nur aus, wenn `down` fehlschlägt, weil sie die Fehlerdetails trägt, die der Fehler selbst nicht nennt.
- `devcontainer shell --quiet` hält den Fortschritt des Starts zurück und gibt ihn nur aus, wenn der Befehl fehlschlägt; `u-boot` gibt selbst nur die Ausnahmen 1 und 2 aus, auf der Standardfehlerausgabe und vor dem Beginn der Sitzung. Die Ausgabe der Sitzung bleibt unverändert, denn die Shell ist nicht Ausgabe von `u-boot` ([`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)).
- `devcontainer up --quiet` hält den Fortschritt zurück und gibt ihn nur aus, wenn der Befehl fehlschlägt; `devcontainer down --quiet` gibt nichts aus ([`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)).
- Mit `--json` ändert `--quiet` nichts: Der Umschlag und seine `diagnostics[]` bleiben vollständig, und die Compose-Ausgabe von `down` wird nicht zurückgehalten.

### SPEC-044 — Features `sshd` und `claude-code` im Devcontainer: Warnungen und Zugangsdaten

Verfeinert [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) (Login-Dienst), [`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten des Agenten), [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Drittanbieter-Feature) und berührt [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts: [`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion)). Der Katalog steht in [`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen), die Wirkung von `--quiet` in [`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet). Die Regeln zu `docker-cli`, `go` und zu Features aus externen Quellen im Sandbox-Profil stehen in [`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen).

**Aktiv** heißt: Ein Feature-Eintrag mit `enabled: true` löst auf die Katalogquelle von `sshd` bzw. `claude-code` auf; verglichen wird die aufgelöste Quelle über alle aktivierten Einträge, nicht der Schlüssel: Auch ein Eintrag mit anderem Schlüssel und `source` gleich der Katalogquelle ist aktiv. Eine `version`-Überschreibung ändert das nicht. Ein Eintrag, dessen aufgelöste Quelle eine andere ist, ist kein aktives `sshd` oder `claude-code`: Ein anderes SSH-Feature oder ein anderer Agent über `source` bleibt unausgewiesen.

**Sandbox-Regel.** Beide Features sind mit und ohne Sandbox-Profil wählbar; das Profil gibt keine externe Quelle frei, und der Katalogschlüssel braucht keinen Eintrag in `devcontainer.featureSources.allow`. Der Generator erzeugt für beide Features nur den Schlüssel unter `features`; keines von beiden ergänzt `runArgs`, `mounts` oder einen Startbefehl (`runArgs` stammt allein aus `nestedRuntime` und Egress-Restriktion). Der Satz gilt für `sshd` und `claude-code`; das Feature `ollama` ist ein Dockerfile-Block, ergänzt mit Sandbox-Profil ein Volume unter `mounts` und gehört nicht dazu ([`SPEC-047`](#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen)); auch `ollama` ergänzt weder `runArgs` noch einen Startbefehl.

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

### SPEC-050 — Katalog-Features im Sandbox-Profil: Abweisung von `docker-cli`, Lockerungen von `go`, Features aus externen Quellen

Verfeinert [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) (kein Container-Runtime-Socket des Hosts, keine zusätzlichen Capabilities, Ausweisung jeder Sicherheitslockerung) und [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) (Features, externe Quellen) und berührt [`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (`--quiet`, Ausnahme 1) sowie [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes `10` und `11`). Die Eigenschaften je Katalogeintrag stehen in [`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen), die Wirkung von `--quiet` in [`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), die Warnungen zu `sshd` und `claude-code` in [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten).

**Aktiv** heißt wie in [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten): Ein Feature-Eintrag mit `enabled: true` löst auf die Katalogquelle von `docker-cli` bzw. `go` auf; verglichen wird die aufgelöste Quelle über alle aktivierten Einträge. Eine `version`-Überschreibung ändert das nicht. Ein Eintrag mit `source:`-Überschreibung ist ein Feature aus einer externen Quelle, auch wenn sein Schlüssel ein Katalogschlüssel ist (`go` mit `source:` ist kein aktives `go`). Ein Eintrag mit `enabled: false` oder ohne `enabled` ist nicht aktiv.

**`docker-cli` im Sandbox-Profil.** Ist das Profil `sandbox` (aus `devcontainer.profile: sandbox` oder aus `--sandbox`) und `docker-cli` aktiv, endet `u-boot generate devcontainer` mit einem fachlichen Fehler: Exit-Code `10`, Code [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) im Umschlag von `--json`. Der Aufruf schreibt nichts: keine Datei und kein Verzeichnis, auch unter `--dry-run` nicht, und `--sandbox` wird nicht in `u-boot.yaml` übernommen. Die Meldung nennt das Feature und die Abhilfe in dieser Reihenfolge: das Abschalten (`u-boot config set devcontainer.features.docker-cli.enabled false`); wer Docker im Container braucht, setzt zusätzlich `u-boot config set devcontainer.sandbox.nestedRuntime docker` für einen Docker-Daemon im Container ([`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)); danach `generate devcontainer` erneut und einen vorhandenen Container neu anlegen (`u-boot devcontainer up --recreate`). `nestedRuntime: docker` allein hebt den Fehler nicht auf; mit abgeschaltetem `docker-cli` erzeugt `generate devcontainer` das Projekt mit `nestedRuntime: docker` (`--privileged` in `runArgs`, kein Socket-Mount). `config set` bleibt möglich: die Prüfung gehört nicht zur Schemaprüfung, die `config` beim Laden von `u-boot.yaml` ausführt. Im Standardprofil bleibt `docker-cli` wählbar: Exit-Code `0`, der Schlüssel `ghcr.io/devcontainers/features/docker-outside-of-docker:1` steht unter `features`, und das Feature löst keine Warnung aus.

**`doctor`.** Die Prüfung `devcontainer.features.allowlist` ([`SPEC-017`](#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry)) meldet denselben Zustand als Fehler: Ist `devcontainer.profile` in `u-boot.yaml` `sandbox` und `docker-cli` aktiv, ist der Eintrag ein Fehler (`doctor` endet bei einem Fehler-Eintrag mit Exit-Code `11`), und die Meldung nennt dieselbe Abhilfe. Verstößt zugleich eine Quelle gegen die Freigabeliste, steht beides in derselben Meldung. Im Standardprofil, bei abgeschaltetem `docker-cli` und für `go` meldet die Prüfung dazu nichts. Es gibt keine eigene Prüfungskennung.

**Lockerungen von `go`.** Im Sandbox-Profil gibt `generate devcontainer` für ein aktives `go` zwei Warnungen aus, je eine zu `--cap-add=SYS_PTRACE` und zu `--security-opt=seccomp=unconfined` (Level `warn`, Code [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil), Subjekt `.devcontainer/devcontainer.json`, Kennzeichen `SecurityRelaxation`); der Text beginnt mit „security relaxation for the devcontainer feature go:“. `init: true` des Features ist keine Sicherheitslockerung und wird nicht ausgewiesen. Der Generator schreibt die Optionen nicht in `devcontainer.json`: `runArgs` und `mounts` sind mit `go` dieselben wie ohne. Eine `version`-Überschreibung ändert die zwei Warnungen nicht; ohne Sandbox-Profil, bei abgeschaltetem oder ungesetztem `go` und bei `go` mit `source:` entstehen sie nicht.

**Features aus externen Quellen.** Im Sandbox-Profil gibt `generate devcontainer` für jeden Eintrag mit `enabled: true` und `source:`-Überschreibung genau eine Warnung aus, in der Reihenfolge der Namen (Level `warn`, Code [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil), Subjekt `.devcontainer/devcontainer.json`, Kennzeichen `SecurityRelaxation`). Der Text nennt den Namen und die Quelle („devcontainer feature <Name> comes from the external source <Quelle>“) und sagt, dass `u-boot` die Metadaten des Features nicht liest und das Profil für dieses Feature keine seiner Eigenschaften zusagt. Der Eintrag wird wie bisher erzeugt. Ohne Sandbox-Profil, bei `enabled: false`, ohne `enabled` und für einen Katalogschlüssel ohne `source:` entsteht diese Warnung nicht.

**`--quiet`.** Die Warnungen zu den Lockerungen von `go` und zu den Features aus externen Quellen tragen das Kennzeichen `SecurityRelaxation` und bleiben unter `--quiet` sichtbar ([`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), Ausnahme 1).

### SPEC-045 — Statusdateien der Startscripte und ihre Ausweisung

Verfeinert [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Degradation, Ausweisung jedes Fallbacks), [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) (Wegfall der Egress-Restriktion) und [`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Ausnahme 2) und berührt [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) (kein Pfad vom Container zum Host) und [`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Ausweisung und Abbildung von `failed` bei `devcontainer up` und `devcontainer shell`, [`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)).

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

`result` ist `fallback` bei der Richtlinie `warn` und eingetretenem Fallback (`vfs`-Storage, Container ohne Daemon, Wegfall der Egress-Restriktion), `failed`, wenn das Startscript mit dem Exit-Code `11` endet, weil die Richtlinie `fail` gilt, und beim Startscript der Egress-Restriktion zusätzlich unter jeder Richtlinie bei den Abbrüchen aus [`SPEC-048`](#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf) (das Startscript schreibt die Datei vor dem Exit-Code `11`), sonst `ok`. `at` ist der Zeitpunkt des Scriptlaufs in UTC, Sekundengenauigkeit, Form `YYYY-MM-DDTHH:MM:SSZ`. Die Datei trägt keinen Freitext und keine Container-Kennung. Das Startscript schreibt in eine Temporärdatei und benennt sie um, sodass kein halbes Dokument lesbar ist. Blockierte verschachtelte User-Namespaces (Exit-Code `11` in beiden Richtlinien) sind ein Umgebungsproblem und kein Fallback; es entsteht keine Datei. Exit-Codes und Standardfehlerausgabe der Startscripte bleiben unverändert.

**Strenge beim Lesen.** Ein Dokument ist gültig, wenn es höchstens 512 Byte umfasst, nur druckbare ASCII-Zeichen und Zeilenwechsel enthält, genau die drei Zeilen oben trägt (ein abschließender Zeilenwechsel ist zulässig), `result` einen der drei Werte hat und `at` genau die Form oben hat (keine Bruchteile, kein Offset). Jede Abweichung, auch Freitext, macht es ungültig. Ausgegeben werden nur die aufgelösten Felder, nie der Inhalt der Datei.

**Ermittlung durch den Host.** `doctor`, `generate devcontainer`, `init` und `devcontainer up` sowie die Startphase von `devcontainer shell` lesen die Dateien über die Container-Engine, ohne ein Programm im Container auszuführen (die Probe und die Sitzung von `shell` führen Programme im Container aus und gehören nicht zum Lesen der Dateien):

- Der Container ist der, der das Label `devcontainer.local_folder` mit dem absoluten Projektordner (dem Basisverzeichnis des Befehls) trägt. Bei mehreren Treffern gilt der zuletzt erzeugte, auch ein gestoppter. `u-boot` setzt kein eigenes Label und fügt keine `runArgs` hinzu. Bei einem Container, den `devcontainer up` oder die Startphase von `devcontainer shell` gestartet hat, setzt die Dev-Container-CLI das Label mit dem Projektordner des Hosts als Wert ([`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)).
- Gelesen wird mit `docker cp <Container>:<Pfad> -`, ohne Link-Auflösung und ohne `docker exec`. Der Tar-Kopf entscheidet vor dem Lesen: nur eine reguläre Datei von höchstens 512 Byte; Symlink, FIFO, Verzeichnis, Gerät und Übergröße machen das Dokument ungültig. Die Container-Kennung wird vor jeder Verwendung gegen `[0-9a-f]{12,64}` geprüft, der Projektordner als ein Argument ohne Shell übergeben, jeder Aufruf der Engine ist auf zehn Sekunden begrenzt.
- Ein entfernter Engine-Zugang (`DOCKER_HOST` oder aktiver Docker-Kontext zeigt nicht auf einen lokalen Unix-Socket) gilt als nicht ermittelbar: Ein Container dort, der denselben Ordnerwert trägt, gehört zu einer anderen Maschine.
- Läuft `u-boot` selbst im Devcontainer (Container-Erkennung), liest `doctor` die Dateien direkt aus dem Dateisystem, nach derselben Typ- und Größenprüfung; fehlt das Verzeichnis (etwa `u-boot` aus dem Runtime-Image), ist der Zustand nicht ermittelbar. `generate devcontainer` und `init` lesen nicht aus dem Dateisystem, sondern fragen immer über die Engine; im Devcontainer bleibt die Abfrage still (mit verschachteltem Docker erreicht sie den verschachtelten Daemon). `devcontainer up` und die Startphase von `devcontainer shell` kennen die Container-Erkennung nicht und fragen ebenfalls immer über die Engine.

**Nicht ermittelbar** heißt: kein Container zum Projektordner, keine oder entfernte Engine, keine Datei (ein Container, der vor den Statusdateien erzeugt wurde oder noch nie gestartet ist, hat keine; er braucht `u-boot generate devcontainer` und einen Neubau; ebenso fehlt die Datei, wenn das Startscript vor dem Schreiben endete, etwa bei blockierten User-Namespaces), ein Lesefehler. Das ist nie ein Fallback.

**Stand.** Jeder Lauf von `egress-init.sh` und `docker-init.sh` überschreibt die Datei `egress` beziehungsweise `daemon`; die Datei beschreibt den letzten Lauf ihres Scripts, nicht notwendig den letzten Start des Containers. Das erzeugte `postStartCommand` verkettet beide Scripte mit `&&`, `egress-init.sh` zuerst; endet `egress-init.sh` mit dem Exit-Code `11` (auch unter `warn`, [`SPEC-048`](#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf)), folgt `docker-init.sh` nach der Semantik von `&&` nicht, und `daemon` bleibt die Datei des letzten Laufs von `docker-init.sh`. `vfs` gilt bis zum Entfernen des Containers wie die `storage.conf`, die das Script schreibt. Ein gestoppter Container wird gelesen; der Eintrag nennt den Zeitstempel und den Vermerk „container stopped".

**Abbildung in `u-boot doctor`.** Je konfigurierter Fähigkeit ein Eintrag, und nur dann: `devcontainer.sandbox.fallback.vfs` (`nestedRuntime: podman`), `….daemon` (`nestedRuntime: docker`), `….egress` (Egress-Restriktion), jeweils mit Sandbox-Profil. Ohne Sandbox-Profil oder ohne die Fähigkeit entsteht kein Eintrag.

| Zustand | Stufe | Kennzeichen `DegradationFallback` (sichtbar unter `--quiet`) |
|---|---|---|
| `result=ok` | OK | nein |
| `result=fallback` | `warn` | ja |
| `result=failed` | `error` (Exit-Code `11`) | nein |
| nicht ermittelbar | `info` mit Grund | nein |
| ungültiges Dokument | `warn` | nein |

`doctor --strict` endet bei `warn` mit Exit-Code `11`, wie bei jedem `warn`-Eintrag. Ein Fehler der Engine verschiebt den Exit-Code nicht: der Bericht entsteht immer, der Fehler ist „nicht ermittelbar". `diagnostics[]` der JSON-Ausgabe trägt wie sonst nur `warn` und `error`; ein `info`-Eintrag erscheint dort nicht.

**Abbildung in `generate devcontainer` und `init`.** Nach dem Schreiben fragen beide für jede Fähigkeit, die `u-boot.yaml` (bei `init` die geschriebene) konfiguriert, denselben Weg ab (ein Container-Lookup je Lauf). `init` schreibt ein Sandbox-Profil ohne `nestedRuntime` und ohne Egress-Restriktion; die geschriebene Konfiguration trägt daher keine abfragbare Fähigkeit, und `init` stellt bei einem bestehenden Projekt keine Abfrage nach den Fähigkeiten der überschriebenen `u-boot.yaml` an. Meldet ein bestehender Container `result=fallback`, steht in der Warnungsausgabe eine Warnung (Level `warn`, Subjekt `.devcontainer/devcontainer.json`, Code [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) für `vfs` und Daemon, [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) für Egress) mit dem Kennzeichen `DegradationFallback`, ohne `SecurityRelaxation`; sie bleibt unter `--quiet` sichtbar ([`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 2). Der Text nennt, dass es der Stand des letzten Starts des bestehenden Containers ist und die neu erzeugten Dateien erst nach einem Neubau wirken. Die Abfrage ist best effort: ohne Engine, ohne Container, ohne Datei, bei ungültigem Dokument und bei `ok` oder `failed` bleibt sie still (keine Meldung, keine Änderung des Exit-Codes); mit `--dry-run` und ohne konfigurierte Fähigkeit findet sie nicht statt. `devcontainer up` und die Startphase von `devcontainer shell` fragen auf demselben Weg, nach einem gelungenen Start für `result=fallback` (Warnung) und nach einem gelungenen wie nach einem fehlgeschlagenen für `result=failed` (Exit-Code `11`; [`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)); das Compose-`up` und jeder andere Befehl fragen nicht.

**Grenze.** Die Datei stammt aus einem Container, in dem ein Agent läuft, und ist dort mit Root-Rechten schreibbar. Auch ein Prozess ohne `sudo` kann sie ersetzen: das Verzeichnis `/var/lib/u-boot/status` gehört dem Container-Benutzer `vscode` (das erzeugte Dockerfile legt es mit `chown vscode:vscode` so an), die Datei dem Benutzer `root`. Sie ist ein Hinweis, kein Beleg: „kein Fallback gemeldet" heißt nicht „kein Fallback aufgetreten". Das Sandbox-Profil bleibt Schadensbegrenzung, keine Grenze ([`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil)).

### SPEC-046 — Egress-Anzeige (`u-boot egress list`): Berechnung, Ausgabe und Abgrenzung

Verfeinert [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) (die wirksame Allowlist samt Herkunft jedes Hosts ist abfragbar) und berührt [`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Abfragebefehl unter `--quiet`), [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes), [`LH-FA-CLI-007`](lastenheft.md#lh-fa-cli-007--dry-run) und [`LH-NFA-USE-004`](lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) (`--json`, Wert `egress` des Feldes `command`) sowie [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) (Clone-Quelle). Die Hosts und die Herkunftsnamen stehen in [`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion).

**Berechnung und Eingabe.** Der Befehl zeigt die Liste, die `u-boot generate devcontainer` im Sandbox-Profil als Zeile `hosts=` in das Startscript `egress-init.sh` schreibt, mit der Herkunft je Host: Beide Listen entstehen aus derselben Berechnung und sind für dieselbe Konfiguration gleich. Der Befehl liest im Arbeitsverzeichnis nur `u-boot.yaml` und die Git-Konfiguration `.git/config`, beides lesend; er braucht weder Netz noch Docker noch einen laufenden Container und schreibt nichts.

**Aktivierungszustand.** Die Restriktion ist aktiv (`active`), wenn beides gilt: das Sandbox-Profil (`devcontainer.profile: sandbox`) und `devcontainer.sandbox.egress.enabled: true`. Ist sie es nicht, zeigt der Befehl dieselben Hosts wie im aktiven Zustand, also die Liste, die bei Aktivierung gälte, und nennt den Grund: das Profil, `egress.enabled` oder beides. Der Zustand ist ein Datum der Ausgabe, keine Warnung und kein Fehler (Exit-Code `0`).

**Clone-Quelle.** Die Clone-Quelle wird aufgelöst wie bei `u-boot generate devcontainer --sandbox`, auch wenn das Profil nicht gesetzt oder die Restriktion nicht aktiv ist: zuerst `devcontainer.sandbox.repository`, sonst die URL des Remotes `origin` in `.git/config`; ihr Host ist der Host der URL (Formen `https://host/pfad` und `benutzer@host:pfad`). Ist `.git` eine Datei (Worktree, Submodul), gibt es kein `.git/config` und damit keine Quelle aus `origin`. Ohne Quelle trägt die Liste keinen Clone-Host: Das ist kein Fehler (der Befehl zeigt `cloneSource: false` und einen Hinweis; `generate devcontainer` warnt dagegen). Für eine URL, die keine dieser Formen hat (etwa ein lokaler Pfad), ist der abgeleitete Host nicht zugesagt. Eine `origin`-URL, die `generate devcontainer --sandbox` abweist (etwa mit Zugangsdaten), weist auch der Befehl ab: derselbe fachliche Fehler, Exit-Code `10`, Diagnose [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil), auch ohne Sandbox-Profil und ohne aktive Restriktion und ohne Liste. Ein `devcontainer.sandbox.repository` mit Zugangsdaten weist schon die Prüfung der Konfiguration ab (Exit-Code `10`, Diagnose [`LH-FA-CONF-002`](lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration)).

**Ausgabe.** Die Textform beginnt mit der Zustandszeile (`Egress restriction: active.` beziehungsweise `Egress restriction: not active (<Grund>). …`), dann folgt eine Tabelle `HOST` und `ORIGIN` (je Host alle Herkünfte, durch Komma getrennt, nach Host sortiert) und ein Block `Notes:` mit den Hinweisen: Subdomains eines Namens sind eingeschlossen (eine und zwei Ebenen darunter; ein Eintrag für eine Subdomain öffnet den übergeordneten Namen bei der Verbindung per Name nicht), ein fehlender Clone-Host (nur dann), und die Abgrenzung zum Zustand eines Containers (unten). Mit `--json` entsteht der Umschlag nach [`LH-NFA-USE-004`](lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) mit `command: "egress"`, `subcommand: "list"`, `status: "ok"`, `diagnostics: []`, `exitCode: 0` und einem Objekt `data` mit den Wahrheitswerten `active`, `sandboxProfile`, `egressEnabled` und `cloneSource` sowie `hosts`, einer nach Host sortierten Liste von Objekten `host` und `sources` (die Herkünfte); ohne Host ist `hosts` die leere Liste `[]`. Die Sprache der Ausgabe ist Englisch.

**Wirkung von `--quiet`.** Die Ausgabe bleibt vollständig und ist byte-gleich zur Ausgabe ohne `--quiet`, die Zustandszeile und die Hinweise eingeschlossen: Sie sind Teil des Ergebnisses der Abfrage, kein Hinweis neben ihm. Mit `--json` ändert `--quiet` nichts ([`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)).

**Optionen und Exit-Codes.** `--dry-run` und `--diff` weist der Befehl vor jeder Berechnung mit Exit-Code `2` ab (mit `--json` als Umschlag mit Diagnose [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes)); zusätzliche Argumente sind ein Aufruffehler (Exit-Code `2`). `u-boot egress` ohne Unterbefehl gibt die Hilfe aus, mit `--json` endet es mit Exit-Code `2` und ohne Umschlag. Fehler: kein `u-boot.yaml` (Exit-Code `10`, Diagnose [`LH-FA-INIT-001`](lastenheft.md#lh-fa-init-001--neues-projekt-initialisieren)), ungültige Konfiguration einschließlich eines ungültigen Eintrags in `devcontainer.sandbox.egress.allow` (Exit-Code `10`, [`LH-FA-CONF-002`](lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration)), unzulässige `origin`-URL (Exit-Code `10`, [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil)), Lesefehler von `u-boot.yaml` (Exit-Code `14`, [`LH-NFA-REL-003`](lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern)), jeder übrige Fehler (Exit-Code `1`, [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes)). Im Textmodus bleibt bei einem Fehler die Standardausgabe leer; mit `--json` entsteht der Fehler-Umschlag mit Exit-Code und Diagnose. Der Befehl führt keinen neuen Diagnose-Code ein ([`SPEC-017`](#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) bleibt unverändert).

**Abgrenzung zum laufenden Container.** Der Befehl zeigt die Liste der Projektdateien, nicht den Zustand eines laufenden Containers; ob die Restriktion beim Start des Containers griff, weist `u-boot doctor` aus ([`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)). Das Startscript mit der Liste wird beim Bau ins Image kopiert. Die Ausgabe des Befehls nennt, dass ein Container die Liste des Images behält, aus dem er gebaut wurde, dass eine geänderte Liste einen Neuaufbau des Containers braucht, und den Aufruf `devcontainer up --remove-existing-container`.

**Hinweis bei `u-boot generate devcontainer`.** Ersetzt der Befehl den verwalteten Block eines bestehenden `.devcontainer/egress-init.sh`, gibt er eine Warnung mit Code [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) und Subjekt `.devcontainer/egress-init.sh` aus: Ein bestehender Container behält die Egress-Liste des Images, aus dem er gebaut wurde, und der Container ist neu aufzubauen. Der Hinweis erscheint nur, wenn der Block eines bestehenden `egress-init.sh` ersetzt wird: nicht beim ersten Schreiben der Datei (auch nicht bei der ersten Aktivierung der Restriktion), nicht beim Abschalten der Restriktion (die Datei wird nicht mehr geplant und bleibt unverändert stehen), nicht ohne Änderung der Datei und nicht unter `--dry-run`; er trägt weder `SecurityRelaxation` noch `DegradationFallback` und entfällt unter `--quiet`, bleibt aber mit `--json` in `diagnostics[]` ([`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)).

### SPEC-047 — Feature `ollama` im Sandbox-Devcontainer: Installation, Volume, Meldungen und Grenzen

Verfeinert [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog-Feature `ollama`), [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Lieferkette) und [`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Wirkung von `--quiet`) und berührt [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil) sowie [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts: [`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion)). Der Katalog steht in [`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen), die Wirkung von `--quiet` in [`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet). Die Grenzen der Zusagen stehen im Lastenheft ([`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features)).

**Aktiv** heißt: Der Eintrag `devcontainer.features.ollama` hat `enabled: true` und keine `source`. Mit `enabled: false`, ohne Eintrag oder mit `source` (externes Feature, [`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen)) erzeugt der Generator keinen Dockerfile-Block, kein Volume und keine Meldung zu Ollama.

**Nur im Sandbox-Profil.** Das Feature wirkt, wenn `devcontainer.profile: sandbox` gilt oder `u-boot generate devcontainer --sandbox` aufgerufen wird. Ohne Sandbox-Profil bleiben `Dockerfile` und `devcontainer.json` byte-gleich zum Lauf ohne den Eintrag (kein Block, kein Volume), und die Befehlsausgabe sagt es (Meldungen unten).

**Installation.** Mit Sandbox-Profil trägt das erzeugte Dockerfile hinter `USER root` einen Block, der Ollama ohne Grafikbeschleunigung (nur CPU) aus dem Release-Archiv `ollama-linux-<arch>.tar.zst` von `github.com/ollama/ollama` installiert. Version und die SHA-256-Prüfsumme je Architektur (`amd64`, `arm64`) stehen als `OLLAMA_VERSION`, `OLLAMA_SHA256_AMD64` und `OLLAMA_SHA256_ARM64` an einer Stelle im Dockerfile; die Spezifikation führt sie nicht ein zweites Mal. Der Block prüft die Prüfsumme mit `sha256sum -c` vor dem Entpacken (eine Abweichung bricht den Image-Build ab, ebenso eine andere Architektur), entpackt nach `/usr/local` ohne die Verzeichnisse `lib/ollama/cuda_*` und `lib/ollama/vulkan`, installiert `zstd` nur dafür und entfernt es danach wieder (das Image enthält kein `zstd`), legt `/home/vscode/.ollama` an und übergibt es dem Benutzer `vscode`. Das Dockerfile enthält außerhalb von Kommentaren weder `ENV` noch `EXPOSE` noch `CMD` noch `ENTRYPOINT` noch `OLLAMA_HOST` noch einen Aufruf von `ollama serve`. Eine `version` ohne `source` ist ein Validierungsfehler von `generate devcontainer` (Exit-Code `10`, nichts wird geschrieben): Die Prüfsummen hängen an der gepinnten Version.

**Volume.** Mit aktivem Feature und Sandbox-Profil trägt `mounts` ein weiteres benanntes Volume `<Projekt>-ollama-${devcontainerId}` mit dem Ziel `/home/vscode/.ollama`, nach dem Workspace-Volume und dem Volume der Container-Runtime (`nestedRuntime`), in dieser Reihenfolge; die `mounts` bestehender Konfigurationen bleiben byte-gleich. Die Modelle in diesem Volume überstehen das Löschen des Containers: Ein neuer Container mit demselben Volume listet das Modell. Dass sie den Neubau des Images überstehen, ist nicht zugesagt. Das Volume enthält außerdem das Schlüsselpaar `id_ed25519`, das der Server beim ersten Start anlegt; es ist nicht weiterzugeben.

**Kein Start, nur Loopback.** Der Generator ergänzt für `ollama` keinen Startbefehl, kein `forwardPorts`, kein `runArgs` und setzt `OLLAMA_HOST` weder im Dockerfile noch in `devcontainer.json`. Der Nutzer startet den Server im Container selbst (`ollama serve`). Bei ungesetztem `OLLAMA_HOST` lauscht der Server der gepinnten Version nur auf `127.0.0.1:11434`; ein Aufruf über die Adresse des Containers scheitert. Setzt der Nutzer `OLLAMA_HOST`, ist das nicht zugesagt. Der Server läuft auf der CPU (`library=cpu` im Start-Protokoll); ein Pull über die beiden Hosts von [`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion) hinter der Egress-Restriktion gelingt, ein nicht gelisteter Host bleibt gesperrt.

**Meldungen der Befehlsausgabe** (Erzeuger: `generate devcontainer`; Level `warn`; bei aktivem Feature genau eine je Lauf, und keine Drittanbieter-Meldung nach [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)):

| Sandbox-Profil | Code | Inhalt |
|---|---|---|
| ja | [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | Quelle ist das Release-Archiv von `github.com/ollama/ollama`; die Version nennt die Meldung nicht, sie verweist auf `OLLAMA_VERSION` im Dockerfile. Geprüft wird allein die SHA-256-Prüfsumme; der Upstream veröffentlicht keine Signatur. Subjekt `.devcontainer/devcontainer.json`. |
| nein | [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) | Das Feature hat ohne Sandbox-Profil keine Wirkung (`devcontainer.profile=sandbox` fehlt); es wurde keine Ollama-Installation erzeugt. |

Beide Meldungen tragen weder `SecurityRelaxation` noch `DegradationFallback`: Die Ausweisung der Lieferkette ist die Meldung der Befehlsausgabe, sie ist keine Sicherheitslockerung des Sandbox-Profils ([`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)). Unter `--quiet` entfallen beide.

**Nicht zugesagt.** GPU-Beschleunigung, ein Anschluss an einen Ollama-Dienst des Hosts, Auswahl oder Verwaltung von Modellen durch `u-boot`, andere Architekturen und Container-Engines als die in [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) genannten und das Überstehen des Neubaus des Images.

### SPEC-048 — Startscript der Egress-Restriktion: Upstream-DNS, Abbruch und Wiederholungslauf

Verfeinert [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) (Mechanismus des Startscripts `egress-init.sh`) und [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Abbruch als Umgebungsproblem) und berührt [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) (Guardrail, keine Sandbox-Grenze). Die Hosts der Allowlist stehen in [`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion), die Statusdatei in [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung). Belege sind Golden Cases des erzeugten Scripts und Integrationstests (Build-Tag `docker`).

**Upstream-DNS.** Das Startscript legt den Benutzer des lokalen Resolvers einmal fest (`nobody`): Die Konfiguration des Resolvers trägt `user=<Benutzer>`, und die UID der Regeln wird zur Laufzeit aus demselben Wert gelesen. Port 53 (UDP und TCP) zu den Upstream-Resolvern ist nur für diese UID offen (`meta skuid`); jede Regel für Port 53 zu den Upstream-Resolvern trägt die Benutzerbedingung. Ein Prozess, der nicht als Benutzer des Resolvers fragt, bekommt von den Upstream-Resolvern keine Antwort (UDP und TCP); erlaubte Namen löst der lokale Resolver auf, und die Verbindung zu einem erlaubten Host kommt zustande. Die `ip6`-Tabelle lässt nur Loopback zu.

**Grenze der Zusage.** Zugesagt ist nur dieser Pfad. Wer als Benutzer des Resolvers fragt (etwa nach einem Wechsel auf diesen Benutzer, der root-Rechte voraussetzt), bekommt von den Upstream-Resolvern eine Antwort; die Restriktion bleibt ein Guardrail ([`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion)). Nicht zugesagt sind der Verkehr zu den Adressen erlaubter Namen und die Anfragen nach Unternamen erlaubter Namen (zur Erreichbarkeit der Unternamen: [`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion)).

**Abbruch.** Das Startscript endet mit dem Exit-Code `11` und schreibt `failed` in die Statusdatei `egress`, unabhängig von `devcontainer.sandbox.onUnavailable`, in vier Abbruchfällen; die Meldung auf der Standardfehlerausgabe nennt den Grund:

1. Der Benutzer des Resolvers existiert im Image nicht. Das Script bricht ab, bevor es Regeln anlegt; im ersten Lauf besteht keine Tabelle, und `/etc/resolv.conf` zeigt nicht auf den lokalen Resolver.
2. Die `nft`-Eingabe mit allen Regeln (siehe Wiederholungslauf) wird nicht angenommen.
3. Der Resolver startet nicht.
4. Der Resolver läuft nicht unter der UID der Regeln (andere UID).

Bei 3 und 4 stehen beide Tabellen, die Kette mit der Richtlinie `drop` besteht, und `/etc/resolv.conf` wird nicht auf den lokalen Resolver gesetzt. Dass die Capability `NET_ADMIN` fehlt, ist kein fünfter Fall dieser Liste, sondern die Zeile der Degradationstabelle ([`LH-FA-DEV-007.a`](#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten)): Exit-Code `11` nur bei `fail`. `u-boot doctor` weist `failed` des Eintrags `devcontainer.sandbox.fallback.egress` als Fehler aus, und `u-boot devcontainer up` und `u-boot devcontainer shell` bilden es ebenfalls auf Exit-Code `11` ab, wenn eine `u-boot.yaml` mit Sandbox-Profil die Egress-Restriktion konfiguriert und die Statusdatei lesbar ist (nach einem gelungenen Lauf der CLI wie nach einem Fehler, auch bei einem `failed` aus einem früheren Start des Containers), sonst bleibt es beim Exit-Code des Starts (`0` oder `12`) ([`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)); Die Meldung nennt vier Ursachen (die Capability, den Benutzer, die Regeln, den Resolver; sie fasst die Abbruchfälle 3 und 4 zum Resolver zusammen, die Zählung unterscheidet sich also von den vier Abbruchfällen), der Hinweis nennt die Capability, den Benutzer und die Regeln und verweist auf die Meldung des Startscripts im Container-Log, weil die Statusdatei keinen Freitext trägt.

**Wiederholungslauf.** Das Script ist wiederholbar, bei jedem Start durch das `postStartCommand` und von Hand im laufenden Container (mit `sudo`). Es schreibt zuerst die Konfiguration des Resolvers und ersetzt dann beide Tabellen (`ip u_boot_egress`, `ip6 u_boot_egress6`) in **einer** `nft -f`-Eingabe: je Tabelle `add table`, `delete table`, `add table`, danach der Aufbau; kein anderer `nft`-Aufruf ändert Tabellen oder Regeln. Danach beendet es den alten Resolver mit begrenztem Warten, startet den neuen, prüft Start und UID und schreibt `/etc/resolv.conf` zuletzt. Ein Lauf, der keinen der Abbrüche trifft (mit der Capability, einem startenden Resolver und einem erreichbaren Upstream-Resolver), endet mit dem Exit-Code `0` und `ok`, mit genau einem Resolver unter der UID der Regeln und mit beiden Tabellen. Nach einem Wiederholungslauf ist das Set der erlaubten Adressen leer und füllt sich neu. Nimmt `nft` die Eingabe nicht an (Fall 2; gepinnt für eine ungültige Regelzeile der `ip`-Tabelle), bleiben in einem wiederholten Lauf die Kette `output` der `ip`-Tabelle, der Resolver und der Eintrag des lokalen Resolvers in `/etc/resolv.conf` des vorigen Laufs unverändert; die Konfigurationsdatei des Resolvers ist zu diesem Zeitpunkt bereits neu geschrieben. Zugesagt ist diese Struktur; nicht zugesagt sind ein Fenster von null Länge ohne Restriktion und eine unterbrechungsfreie Namensauflösung.

### SPEC-049 — Devcontainer ohne Editor starten und stoppen (`u-boot devcontainer up` und `down`): Aufruf, Vorab-Prüfung, Umgebung, Ausgabe und Exit-Codes

Verfeinert [`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) und berührt [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) (der Hilfs-Container ist nicht das erzeugte Ergebnis), [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Ausweisung der Fallbacks), [`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten), [`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (`--quiet`), [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes), [`LH-FA-CLI-007`](lastenheft.md#lh-fa-cli-007--dry-run) (`--dry-run`) und [`LH-NFA-USE-004`](lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) (`--json`, Wert `devcontainer` des Feldes `command`). Die Statusdateien stehen in [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), die Vorab-Prüfung vor Compose-Aufrufen in [`SPEC-042`](#spec-042--aufruf-der-compose-engine-und-preflight). Belege sind Unit-Tests und Integrationstests (Build-Tag `docker`). Digests, Versionen und Prüfsummen des Hilfs-Images stehen in dessen eingebetteten Dateien, nicht in dieser Stelle. Die Überschrift nennt `up` und `down`; `u-boot devcontainer shell` steht im Unterabschnitt *`shell`* nach der Tabelle der Exit-Codes. Wo ein Absatz davor „`up`“ sagt, gilt er für die Startphase von `shell`, soweit der Unterabschnitt es sagt.

**Aufruf.** `u-boot devcontainer up [--recreate]`, `u-boot devcontainer down` und `u-boot devcontainer shell`. Der Projektordner ist das Arbeitsverzeichnis des Aufrufs. Ohne Unterbefehl, mit einem unbekannten Unterbefehl, mit zusätzlichen Argumenten sowie mit `--dry-run` oder `--diff` endet der Befehl mit Exit-Code `2`, denn die Unterbefehle schreiben keine Projektdatei. Mit `--json` entsteht bei `--dry-run`, `--diff` und zusätzlichen Argumenten ein Umschlag mit Exit-Code `2` (bei den ersten beiden mit Diagnose [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes)); der Aufruf ohne Unterbefehl endet auch mit `--json` ohne Umschlag. `shell` weist außerdem `--json` mit Exit-Code `2` ab, weil der Befehl eine interaktive Sitzung öffnet und die Standardausgabe ihr gehört; die Abweisung ist ein Umschlag mit `subcommand: "shell"` und Diagnose [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes); keine dieser Abweisungen startet etwas.

**Voraussetzungen von `up` und `shell` (Exit-Code `10`).** `up` und `shell` prüfen vor der Ausweisung des Sockets und vor dem Aufruf des Launchers: (1) der Projektordner ist ein absoluter Pfad ohne Komma, Anführungszeichen, Zeilenwechsel, Wagenrücklauf und NUL, weil ein `--mount`-Wert diese Zeichen nicht trägt; die Meldung nennt den Pfad nicht; (2) `.devcontainer/devcontainer.json` existiert; (3) jeder Name, auf den diese Datei mit `${localEnv:<NAME>}` verweist, ist ein einfacher Variablenname (`[A-Za-z_][A-Za-z0-9_]*`). Ein Verstoß endet mit Exit-Code `10`, ein Lesefehler der Datei mit Exit-Code `14`; in beiden Fällen startet kein Hilfs-Container, und die Ausweisung des Sockets entfällt. `down` liest keine Projektdatei und braucht keine `devcontainer.json`.

**Vorab-Prüfung der Engine** (bei `up` und in der Startphase von `shell` vor dem Bau des Hilfs-Images und dem Start des Hilfs-Containers, bei `shell` zusätzlich vor der Probe der Sitzung, bei `down` vor `docker stop`): (1) das `docker`-Binary ist auf dem `PATH`, (2) der Daemon antwortet, (3) die Engine ist lokal (`DOCKER_HOST` oder der aktive Docker-Kontext zeigt auf einen lokalen Unix-Socket, wie in [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)). Scheitert Schritt 1 oder 2, lässt sich der Engine-Zugang nicht bestimmen (etwa weil der Kontext nicht lesbar ist) oder ist die Engine nicht lokal (`DOCKER_HOST` mit `tcp://…`, ein Kontext mit `ssh://…`), endet der Befehl mit Exit-Code `11`. Nach einem gescheiterten Schritt laufen weder Bau noch Start noch `docker stop` noch `docker exec`.

**Hilfs-Image.** `u-boot` baut das Image lokal aus einem eingebetteten Kontext aus genau drei Dateien (`Dockerfile`, `package.json`, `package-lock.json`). Der Tag hat die Form `u-boot-devcontainer-cli:<exakte Version der Dev-Container-CLI>-<zwölf Hexzeichen>`. Gebaut wird nur, wenn kein Image mit diesem Tag besteht: `docker build --force-rm -t <Tag> -` mit dem Kontext als Tar-Strom auf der Standardeingabe. Pin-Form: beide Basis-Images stehen mit einem Digest (`@sha256:…`) im `FROM`; `package.json` pinnt die CLI auf eine exakte Version, jeder Eintrag der Lockdatei trägt eine `integrity`, die Lockdatei nennt dieselbe Version der CLI, und die Installation ist `npm ci --ignore-scripts`. Eine verfälschte `integrity` bricht den Bau. Scheitert der Bau, endet der Befehl mit Exit-Code `12`; die Meldung nennt das Ende der Standardfehlerausgabe des Baus.

**Hilfs-Container.** Der Aufruf lautet:

```
docker run --rm --name <Name> \
  --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock \
  --mount type=bind,source=<Ordner>,target=<Ordner> \
  [-e <NAME> …] <Hilfs-Image> \
  devcontainer up --workspace-folder <Ordner> [--remove-existing-container]
```

Der Name beginnt mit einem festen Präfix. Die Argumentliste enthält weder `-v` noch `--volume`, `--privileged`, `--network` oder `--cap-add`, und sie bindet kein Home-Verzeichnis ein. Der Projektordner ist schreibend eingehängt; ein `initializeCommand`, der in ihn schreibt, gelingt. `--recreate` ergänzt `--remove-existing-container`.

- **Spiegelpfad.** Der Projektordner liegt im Hilfs-Container unter demselben absoluten Pfad wie auf dem Host, und derselbe Pfad geht als `--workspace-folder` an die CLI. Das Label `devcontainer.local_folder` des gestarteten Containers trägt damit den Pfad des Hosts, und die Suche über dieses Label ([`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)) findet den Container; ein Pfad unter `/w/<Name>` im Hilfs-Container ergäbe ein anderes Label. Das Label setzt die Dev-Container-CLI, `u-boot` setzt kein eigenes. Ein Projektordner, den die Engine nicht sieht, ergibt Exit-Code `12`: `--mount type=bind` bricht ab und legt den Quellpfad nicht an.
- **Lebensdauer.** Nach `up` bleibt kein Hilfs-Container zurück. Endet `docker run` mit einem Fehler oder bricht der Aufrufer ab, entfernt `u-boot` den Hilfs-Container mit `docker rm -f <Name>`; gelingt das nicht, nennt die Fehlermeldung es. Ein Abbruch beendet den Prozess mit einem Exit-Code ungleich `0`; der Dev-Container, den die CLI bis dahin angelegt hat, bleibt bestehen.
- **Fortschritt.** Die Standardfehlerausgabe der CLI geht als Fortschritt an die Konsole. Die Ergebniszeile der CLI steht auf deren Standardausgabe, nicht auf der Standardfehlerausgabe.

**Umgebung und Zugangsdaten.** Die Namen liest `u-boot` aus dem Text von `.devcontainer/devcontainer.json`: jede Form `${localEnv:<NAME>}` und `${localEnv:<NAME>:<Vorgabe>}`, sortiert und ohne Duplikate; ein Verweis in einem Kommentar zählt mit, andere Formen (`${env:…}`, `${containerEnv:…}`) werden nicht gelesen. Der Hilfs-Container bekommt nur die Namen, die auf dem Host gesetzt sind, als `-e <NAME>` ohne Wert. Variablen, auf die die Datei nicht verweist (etwa Proxy-Variablen), reicht der Befehl nicht durch.

- **Wo der Wert nicht steht.** Das `docker run` von `u-boot` trägt nur den Namen. Mit der Form `remoteEnv` steht der Wert zudem in keinem Label und keinen Metadaten des gestarteten Containers, nicht im Hilfs-Image, in keiner Datei unter dem Projektordner und in keinem Volume des Sandbox-Containers.
- **Grenzen.** Die Dev-Container-CLI schreibt den aufgelösten Wert für die Dauer eines Lebenszyklus-Befehls als `docker exec -e <NAME>=<Wert>` in die Kommandozeile eines Prozesses im Hilfs-Container (auch bei `remoteEnv`). Bei `containerEnv` steht der aufgelöste Wert in den Metadaten (`Config.Env`) des Dev-Containers, solange dieser besteht.
- **Maskierung.** In dem, was `u-boot` als Fortschritt weitergibt, und in seinen Fehlermeldungen ersetzt es den Wert einer durchgereichten Variable durch `[masked]`, wenn er unverändert vorkommt und mindestens vier Byte lang ist; zerlegte oder kodierte Formen und kürzere Werte bleiben stehen. Die Felder der Ergebniszeile sind Daten und werden nicht maskiert. Die letzten Bytes des Fortschritts (bis zur Länge des längsten maskierbaren Werts minus eins) erscheinen mit dem nächsten Schreiben oder am Ende des Aufrufs.

**Start und Ergebniszeile.** Die letzte nicht leere Zeile der Standardausgabe der CLI ist die Ergebniszeile (JSON); sie bleibt auch nach umfangreicher Ausgabe davor lesbar. Der Start gilt als gelungen, wenn `outcome` den Wert `success` hat und `containerId` aus Hexzeichen besteht; `remoteUser` und `remoteWorkspaceFolder` übernimmt `u-boot` aus der Zeile. Ein anderes `outcome` (die Meldung der CLI steht im Fehlertext), eine fehlende oder keine JSON-Ergebniszeile und ein Exit-Code des `docker run` ungleich `0` enden mit Exit-Code `12`, vorbehaltlich des Urteils nach dem Zustand des Containers (Abbildung auf `11`, Absatz *Fallbacks und Abbruch eines Startscripts*). `up` nach `down` startet denselben Container; `up --recreate` entfernt ihn und legt einen neuen an, benannte Volumes bleiben samt Inhalt (der Workspace des Sandbox-Profils). Die Images `vsc-<Ordnername>-…`, die die CLI baut, bleiben nach `down` und nach `--recreate` bestehen.

**Lebenszyklus-Befehle bei `up` gegen einen laufenden Container.** Gegen einen laufenden Container führt die Dev-Container-CLI bei `up` das `postCreateCommand` und das `postStartCommand` nicht erneut aus, das `postAttachCommand` bei jedem `up`; nach `down` und `up` im selben Container läuft das `postStartCommand` wieder, das `postCreateCommand` nicht. Gemessen ist das mit der gepinnten Version der CLI, mit Befehlen in der Form einer Zeichenkette, die gelingen, unter Linux `amd64` mit lokaler Engine (Docker 29.9.2), ein Lauf je Aufruf. Nicht zugesagt sind das Verhalten des `initializeCommand`, des `onCreateCommand` und des `updateContentCommand`, die Array- und die Objektform der Befehle, scheiternde Befehle gegen einen laufenden Container und andere Hosts und Engines. Für `shell` ist nur zugesagt, dass seine Startphase ein `up` ohne `--recreate` ist; wann ein `postAttachCommand` bei `shell` läuft, ist nicht zugesagt.

**Fallbacks und Abbruch eines Startscripts.** Nach einem Lauf der Dev-Container-CLI, der nicht an der Vorab-Prüfung scheitert und nicht abgebrochen wird, liest `up` die Statusdateien auf dem Weg von [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), und zwar je Fähigkeit, die die `u-boot.yaml` im Projektordner mit Sandbox-Profil konfiguriert. Bei `result=fallback` nach einem gelungenen Start entsteht eine Warnung (Level `warn`, Subjekt `.devcontainer/devcontainer.json`, Code [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) für `vfs` und Daemon, [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) für Egress) mit dem Kennzeichen `DegradationFallback` und ohne `SecurityRelaxation`; der Text beginnt mit „state of the devcontainer container at its last start (<Zeitpunkt>): “ und nennt die Folge des Fallbacks. Er enthält weder „existing“ noch die Anweisung zum Neubau, weil `up` keine Datei neu erzeugt. Jeder andere Zustand außer `failed` bleibt still. Ohne `u-boot.yaml`, ohne Sandbox-Profil und ohne die Fähigkeit entsteht keine Ausweisung, und ein gelungener Start endet mit Exit-Code `0`.

`up` urteilt nach dem Zustand des Containers, nicht nach dem Verlauf des Aufrufs: Meldet die Statusdatei einer konfigurierten Fähigkeit `result=failed`, endet der Befehl mit Exit-Code `11`, wie `doctor` ([`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), [`SPEC-048`](#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf), [`LH-FA-DEV-007.a`](#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten)). Das gilt nach einem Lauf der CLI, der mit `success` endete, wie nach einem Fehler der Klasse `12` (dann `11` statt `12`), und auch dann, wenn die Datei aus einem früheren Start desselben Containers stammt, denn sie beschreibt den Zustand des Containers ([`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), *Stand*). Ein erneutes `up` ohne `--recreate` gegen einen Container mit gescheitertem Startscript endet damit wieder mit `11`; `u-boot devcontainer up --recreate` legt den Container neu an. Melden mehrere konfigurierte Fähigkeiten `failed`, nennt die Meldung die erste in der Reihenfolge der Prüfung (`vfs` vor `egress`). Nach einem Fehler der CLI behält der Fehlertext deren Ausgabe; nach einem gelungenen Lauf nennt er den Satz der Fähigkeit (derselbe wie bei `doctor`), den Zeitpunkt der Datei als `(reported at <Zeitpunkt nach RFC 3339>)`, dass der Container läuft und dieser Aufruf das Startscript nicht ausgeführt hat, und `u-boot devcontainer up --recreate`. Fallback-Warnungen anderer Fähigkeiten entfallen bei `11`. Die Klassen `11` und `12` überschneiden sich nicht.

Das Urteil nach dem Zustand setzt eine lesbare Statusdatei einer konfigurierten Fähigkeit voraus. Ohne `u-boot.yaml`, ohne Sandbox-Profil, ohne die Fähigkeit und ohne lesbare Statusdatei (etwa bei blockierten verschachtelten User-Namespaces, bei denen keine Datei entsteht) bleibt der Exit-Code der Startphase: `0` nach einem gelungenen Lauf der CLI, `12` nach einem Fehler dieser Klasse, und ein früherer Abbruch eines Startscripts bleibt unbemerkt. Auch ein ungültiges Dokument, ein Lesefehler, ein nicht gefundener Container sowie `ok` und `fallback` einer Fähigkeit ändern den Exit-Code nicht. Fehler der Vorab-Prüfung und ein Abbruch lesen keine Statusdatei und behalten ihre Klasse. Die Lesung findet den Container über das Label des Projektordners ([`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)); tragen mehrere Container dasselbe Label, ist nicht zugesagt, dass der gelesene der gestartete ist.

**Ausweisung des Sockets.** `up` und `shell` geben eine Warnung aus (Level `warn`, Code [`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Kennzeichen `SecurityRelaxation`, nicht `DegradationFallback`). Sie entsteht vor dem Aufruf des Launchers (bei `up` vor der Abfrage, ob der Container schon läuft); dessen Vorab-Prüfung der Engine, der Bau und der Start folgen danach. Bei einem späteren Fehler steht sie weiter in der Ausgabe. Sie bleibt unter `--quiet` ([`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), Ausnahme 1). Der Text ist Englisch, in Konsole und JSON gleich, und sagt fünf Dinge: (1) „this command mounts the host container-engine socket into a short-lived helper container“, (2) „root-equivalent access“ zur Engine für die Dauer des Aufrufs, bei `shell` für die Dauer des Starts („for the duration of the call“ beziehungsweise „for the duration of the start“; die Sitzung von `shell` bindet den Socket nicht ein), (3) „the devcontainer configuration that u-boot generates mounts it if a feature that does so is enabled, such as the built-in feature docker-cli (the sandbox profile rejects that feature)“, (4) „a hand-written .devcontainer/devcontainer.json can mount it“, (5) „an initializeCommand declared there runs with this access“. Der Text sagt bei (3) nicht „only“: `u-boot` liest die Metadaten eines Features aus einer externen Quelle nicht ([`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen)). Er sagt nicht, dass der gestartete Devcontainer das Socket nicht erhalte. `down` startet keinen Hilfs-Container und weist nichts aus.

**Ausgabe von `up`.**

- **Text.** Auf der Standardausgabe stehen in dieser Reihenfolge die Ausweisung des Sockets (sofort, vor dem Start), die Zeile `Devcontainer started: container <erste zwölf Zeichen der Kennung>, remote user <Benutzer>, workspace <Pfad>` und die Fallback-Warnungen, jede als `Warning (<Code>): <Text>`. War der Container schon gestartet (nächster Punkt), lautet die Zeile `Devcontainer already running: container <erste zwölf Zeichen der Kennung>, remote user <Benutzer>, workspace <Pfad>`. Der Fortschritt geht an die Standardfehlerausgabe.
- **Schon laufender Container.** Ohne `--recreate` fragt `up` nach der Ausweisung des Sockets und vor dem Aufruf des Launchers über das Label des Projektordners ([`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)), ob der Container des Projekts läuft. Lief er und meldet die CLI dieselbe Kennung, gilt er als schon laufend; das ändert nur die Zeile und das Feld `alreadyRunning`, weder den Exit-Code noch `status` noch die Warnungen. Ein gestoppter oder fehlender Container, eine andere Kennung und eine Abfrage, die scheitert, ergeben „nicht schon laufend“; ein Fehler der Abfrage ist kein Fehler des Befehls. Mit `--recreate` findet die Abfrage nicht statt: die Zeile lautet immer `started`, und `alreadyRunning` ist `false`.
- **`--quiet`.** Die Erfolgszeile, auch die Zeile `already running`, und der Fortschritt entfallen; die Ausweisung und die Fallback-Warnungen bleiben ([`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Ausnahmen 1 und 2); jede andere Warnung entfällt. Schlägt der Aufruf fehl, gibt der Befehl den zurückgehaltenen Fortschritt nachträglich auf der Standardfehlerausgabe aus.
- **`--json`.** Genau ein Umschlag auf der Standardausgabe mit `command: "devcontainer"` und `subcommand: "up"`; es gibt keine Fortschrittsausgabe. `data` trägt `containerId` (die volle Kennung), `remoteUser`, `remoteWorkspaceFolder` und `alreadyRunning` (immer vorhanden, `true` nur für einen schon laufenden Container). `diagnostics[]` führt die Ausweisung zuerst (mit dem Text der Konsole, ohne `file`), danach die Fallback-Warnungen (mit `file`). `status` folgt [`LH-FA-CLI-007.a`](#lh-fa-cli-007a--konvention-für-diagnose-codes-und-kopplung-des-status): Der Umschlag eines gelungenen `up` trägt wegen der Ausweisung `status: "warn"` bei `exitCode: 0`. `--quiet` ändert mit `--json` nichts. Bei einem Fehler enthält `diagnostics[]` die Ausweisung (wenn sie entstanden ist), danach den Fehler (Level `error`); `data` trägt dann keine `containerId`. Der Code der Fehler-Diagnose ist [`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) für die Exit-Codes `10`, `11` und `12`, [`LH-NFA-REL-003`](lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern) für `14` und [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes) sonst; sie nennt den Pfad des Projektordners nicht.

**`down`.** Der Befehl findet den Container über das Label des Projektordners ([`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)). Läuft er, ruft `u-boot` `docker stop <Kennung>` auf; der Befehl entfernt weder Container noch Volume noch Image und braucht weder Hilfs-Container noch Socket. Ohne Container und bei einem schon gestoppten Container endet er mit Exit-Code `0` und ohne `docker stop`. Die Textausgabe nennt den Fall: `Devcontainer stopped: container <zwölf Zeichen> (not removed)`, `Devcontainer container <zwölf Zeichen> is not running; nothing to stop` oder `No devcontainer container found for this project; nothing to stop`; unter `--quiet` entfällt sie. Mit `--json` trägt `data` `containerId` (leer ohne Container) und `stopped` (`false`, wenn nichts zu stoppen war), `status` ist `ok` und `diagnostics` leer; `--quiet` ändert daran nichts. Kann die Engine nicht befragt werden (kein `docker`, der Daemon antwortet nicht, die Engine ist nicht lokal), endet `down` mit Exit-Code `11`, nicht mit „kein Container“ und Exit-Code `0`; ein Fehler von `docker stop` endet mit Exit-Code `12`.

**Exit-Codes** ([`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes)); unter `--json` ist `exitCode` des Umschlags gleich dem Exit-Code des Prozesses:

| Code | `up` | `down` | `shell` |
|---|---|---|---|
| `0` | gestartet oder schon laufend | gestoppt oder nichts zu stoppen | die Sitzung hat begonnen und ist zu Ende, unabhängig vom Exit-Status der Shell |
| `2` | Aufruf (siehe *Aufruf*) | Aufruf (siehe *Aufruf*) | Aufruf (siehe *Aufruf*), auch `--json` |
| `10` | Voraussetzungen von `up` verletzt | — | Voraussetzungen verletzt, wie bei `up` |
| `11` | Vorab-Prüfung der Engine gescheitert; Abbruch eines Startscripts mit `failed` in der Statusdatei (nach dem Zustand des Containers) | Vorab-Prüfung oder Befragung der Engine gescheitert | wie bei `up`; es öffnet sich keine Sitzung |
| `12` | Bau des Hilfs-Images, Start des Hilfs-Containers und Fehler der CLI (namentlich `outcome` ungleich `success`, kein Ergebnis, Fehler eines Lebenszyklus-Befehls, Bau des Dev-Containers, ein Ordner, den die Engine nicht sieht) | `docker stop` gescheitert | Fehler des Starts wie bei `up`; die Probe der Sitzung scheitert; der Client scheitert vor dem Beginn der Sitzung (Unterabschnitt *`shell`*) |
| `14` | Lesefehler von `.devcontainer/devcontainer.json` | — | Lesefehler von `.devcontainer/devcontainer.json` |
| `1` | jeder andere Fehler | jeder andere Fehler | jeder andere Fehler |

**Reichweite.** Ohne `docker` im `PATH` enden alle drei Unterbefehle mit Exit-Code `11`. Die Hosts und Engines, auf denen er zugesagt ist, und seine Grenze im ausgelieferten Container-Image nennt die Anforderung ([`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)).

**`shell`: Startphase und Sitzung.** `u-boot devcontainer shell` hat zwei Phasen. Die *Startphase* ist `up` ohne `--recreate`: Voraussetzungen (Exit-Code `10`), Ausweisung des Sockets, Vorab-Prüfung der Engine, Hilfs-Image, Hilfs-Container samt seiner Umgebung, Ergebniszeile der CLI, Fallback-Warnungen und das Urteil nach dem Zustand des Containers sind die von `up`. Es entfällt nur die Abfrage, ob der Container schon lief, weil `shell` kein `alreadyRunning` ausgibt. Meldet die Statusdatei einer konfigurierten Fähigkeit `failed`, endet der Befehl mit Exit-Code `11` und öffnet keine Sitzung, auch gegen einen Container, der schon lief und dessen Startscript dieser Aufruf nicht ausführte. Die *Sitzungsphase* beginnt nach einem gelungenen Start und öffnet im gestarteten Container eine interaktive Shell mit der Kennung, dem Benutzer (`remoteUser`) und dem Arbeitsverzeichnis (`remoteWorkspaceFolder`) der Ergebniszeile.

- **Sitzung.** Der `docker`-Client des Hosts führt zwei Aufrufe aus und sonst nichts im Container: die Probe `docker exec -u <Benutzer> -w <Arbeitsverzeichnis> <Kennung> bash -c true` (ohne `-i`, `-t` und Variablen; endet sie nicht mit Status `0`, dieselbe Probe mit `sh`) und die Sitzung `docker exec -i [-t] -u <Benutzer> -w <Arbeitsverzeichnis> [-e <NAME> …] <Kennung> <Shell>`. Die Shell der Sitzung ist die erste, deren Probe mit Status `0` endet (`bash`, sonst `sh`); sie startet ohne Argument. Die Argumentliste enthält weder `-v` noch `--volume`, `--mount`, `--privileged`, `--network`, `--cap-add` oder `--detach`. Die Sitzung läuft ohne Hilfs-Container: Während der Sitzung besteht keiner, und für die Sitzung bindet der Befehl keinen Container-Runtime-Socket des Hosts ein. Zugesagt ist, was der Befehl für die Sitzung einbindet, nicht der Zustand der Shell: ein eigener Socket im Container (`nestedRuntime: docker`) und ein Socket, den die Konfiguration selbst einbindet (etwa über das Katalog-Feature `docker-cli`), sind nicht Gegenstand dieser Aussage.
- **Umgebung der Sitzung.** `u-boot` liest den Text von `.devcontainer/devcontainer.json` ohne seine Kommentare (`//` bis zum Zeilenende und `/* … */`; Zeichenketten bleiben, ein auskommentierter Eintrag gibt also nichts weiter) und sammelt jeden Eintrag der Form `"<NAME>": "${localEnv:<NAME>}"` mit demselben einfachen Variablennamen im Schlüssel und im Verweis, sortiert und ohne Duplikate, an jeder Stelle der Datei. Das ist die Form, in der die erzeugte Konfiguration unter `remoteEnv` die Zugangsdaten und die Variablen eines Agenten führt ([`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)). Ein Verweis anderer Form (ein anderer Name, ein Vorgabewert, Teil eines längeren Werts, ein Pfad in `mounts`) ist kein Eintrag; der Hilfs-Container der Startphase bekommt ihn dennoch wie bei `up`. Die Sitzung bekommt die Namen als `-e <NAME>` ohne Wert, und nur die, die auf dem Host gesetzt sind; der Client liest den Wert aus seiner Umgebung, und der Wert steht in keiner Kommandozeile. Eine auf dem Host ungesetzte Variable fehlt in der Sitzung. `u-boot` maskiert in der Sitzung nichts: Was die Shell ausgibt, erreicht den Aufrufer unverändert; die Maskierung gilt für den Fortschritt und die Fehlermeldungen der Startphase.
- **Terminal.** Die Sitzung bekommt ein Terminal (`-t`), wenn die Standardeingabe des Befehls eine Datei, ein Zeichengerät und nicht das Nullgerät ist; die Erkennung benutzt kein `ioctl`. Mit Terminal gibt der Launcher Standardeingabe und Standardausgabe unverändert an den Client, damit der Client das Terminal dort findet. Ob die Ausgabe ebenfalls ein Terminal ist, prüft die Erkennung nicht (Eingabe Terminal, Ausgabe Pipe: nicht zugesagt). Ein anderes Zeichengerät als ein Terminal (etwa `/dev/zero`) gälte als Terminal; lehnt der Client es ab („stdin is not a terminal“), endet der Befehl mit `12`. Der Terminal-Pfad ist in drei Teilen zugesagt: die Weitergabe der Erkennung an die Sitzung, die Sitzung am gebauten Binary und die Sitzung am Launcher, beide unter einem Pseudo-Terminal; die Strecke vom Terminal des Nutzers bis zur Sitzung ist nicht als Ganzes zugesagt.
- **Ende der Sitzung und Signale.** Die Sitzung endet mit dem Ende der Shell (`exit`), der Befehl dann mit `0`. Strg-C erreicht mit Terminal die Shell im Container. Eine Shell, die einen Befehl ausführt, bleibt samt Befehl bestehen, wenn die Sitzung von außen endet: ohne Terminal bei Abbruch des Kontexts, mit Terminal bei `SIGKILL` am Client; `u-boot` räumt sie nicht auf. Eine untätige Shell ohne Terminal endet bei Abbruch des Kontexts; eine untätige Shell mit Terminal bleibt bei `SIGKILL` am Client bestehen. Mit Terminal beobachtet der Launcher den Kontext nach dem Beginn der Sitzung nicht; ohne Terminal sendet er bei einem Abbruch `SIGINT` an den Client und meldet den Abbruch als Fehler des Kontexts. Nicht zugesagt sind: dass die Shell mit `u-boot` endet, dass ein Signal an `u-boot` die Sitzung beendet (mit Terminal beendet es sie nicht), das Schließen des Terminals und Strg-C ohne Terminal.
- **Art der Shell.** Mit Terminal ist die Shell keine Login-Shell (`shopt -q login_shell` ergibt nein), `$0` ist `bash`, sie ist interaktiv, liest `~/.bashrc` und liest `~/.profile` nicht; ohne `bash` im Container ist `sh` die Shell. Nicht zugesagt sind `/etc/profile`, die Art der Shell ohne Terminal und Einstellungen, die nur Login-Shells lesen.
- **Messbedingungen** der Aussagen zum Ende der Sitzung, zu den Signalen und zur Art der Shell: Docker 29.9.2 als Engine und Client, Linux `amd64`, lokale Engine; der Container ist ein `sleep 3600` des Basis-Images `mcr.microsoft.com/devcontainers/base:debian`, nicht ein von `devcontainer up` erzeugter; das Terminal ist ein Pseudo-Terminal von `script`; ein Lauf je Fall. Andere Engines, Hosts und Terminals sind nicht zugesagt.
- **Ausgabe.** Alles, was `u-boot` selbst ausgibt, steht auf der Standardfehlerausgabe, in dieser Reihenfolge: die Ausweisung des Sockets, der Fortschritt der Startphase und jede Fallback-Warnung (als `Warning (<Code>): <Text>`, im Wortlaut von `up`); bei einem Fehler folgt dieser. Die Fallback-Warnungen erscheinen vor dem Beginn der Sitzung. Auf der Standardausgabe steht nur, was die Sitzung ausgibt; eine Erfolgszeile gibt es nicht. Standardeingabe, Standardausgabe und Standardfehlerausgabe des Befehls gehen unverändert an die Sitzung.
- **`--quiet`** ([`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)). Der Fortschritt der Startphase wird zurückgehalten und nur ausgegeben, wenn der Befehl fehlschlägt; die Ausweisung des Sockets und die Fallback-Warnungen bleiben (Ausnahmen 1 und 2 von [`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging)); jede andere Warnung entfällt. Die Ausgabe der Sitzung bleibt unverändert: `--quiet` wirkt nur auf die Ausgabe des Starts.
- **`--yes` und `--no-interactive`** wirken bei `shell` nicht; der Befehl fragt nichts.
- **Exit-Code `12`.** Der Befehl endet mit `12`, wenn die Probe scheitert (Benutzer, Arbeitsverzeichnis oder Shell fehlen im Container, oder der Container läuft nicht; die Meldung nennt das Ende der Ausgabe jedes Versuchs) und keine Sitzung öffnet, und wenn der Client `docker exec` vor dem Beginn der Sitzung mit einem eigenen Fehler beendet (*Grenzen der Erkennung*). Der Exit-Status der Shell nach dem Beginn der Sitzung ist kein Fehler: Er erreicht den Aufrufer nicht, auch nicht nach dem Ende des Containers während der Sitzung; der Befehl endet mit `0`.
- **Grenzen der Erkennung** (Fall „die Probe gelingt, die Sitzung beginnt nicht“, etwa weil der Container zwischen Probe und Sitzung endet oder Benutzer, Arbeitsverzeichnis oder Shell dort verschwinden; [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes) bindet ihn an `10`, `11` oder `12`). Der Exit-Status des Clients trennt den Fehler des Clients nicht von dem einer Shell, die mit einem fehlgeschlagenen Befehl endet (Status 1, 126 und 127 hat auch eine Shell); der Launcher erkennt den Fehler an der Meldung des Clients und meldet ihn als `12`, auch „Engine nicht erreichbar“ (`11` bleibt der Vorab-Prüfung vorbehalten). Erkannt wird:
  - eine Sitzung, deren erste nicht leere Zeile der Fehlerausgabe mit `Error response from daemon: ` beginnt und die Kennung des Containers oder den Benutzer der Sitzung nennt (Container läuft nicht oder fehlt, Benutzer fehlt), bei einem Status ungleich `0`; ohne Terminal muss die Standardausgabe leer sein, mit Terminal gehört sie dem Client;
  - die Zeile `Cannot connect to the Docker daemon at `, wenn die Engine auch bei einer erneuten Frage nicht antwortet;
  - die Zeile `cannot attach stdin to a TTY-enabled container because stdin is not a terminal`, wenn ein Terminal angefordert war;
  - ohne Terminal die Zeile `OCI runtime exec failed: ` als erste nicht leere Zeile der Standardausgabe bei leerer Fehlerausgabe und Status `127` (Arbeitsverzeichnis oder Shell nach der Probe entfernt); mit Status `1` gilt dieselbe Zeile als Sitzung;
  - ohne Terminal auf der Engine 28.0.4 die Zeile `unable to find user <Benutzer der Sitzung>: …` als einzige Zeile der Standardausgabe (keine zweite nicht leere Zeile, nichts abgeschnitten) bei leerer Fehlerausgabe und Status `126`; der Doppelpunkt ist die Namensgrenze, `vscode2` ist nicht `vscode`.

  Alles andere gilt als Sitzung: eine Meldung des Clients nach anderer Ausgabe der Sitzung, auch hinter einer Warnzeile (nur die erste nicht leere Zeile zählt), die Meldung der Engine, solange sie antwortet, und jede Meldung, die nicht zu den genannten gehört. Mit Terminal ist die Erkennung auf die Meldungen der Fehlerausgabe beschränkt, und von ihnen sind „Container läuft nicht“ und „Standardeingabe ist kein Terminal“ zugesagt; die OCI-Zeile und die Zeile des fehlenden Benutzers der Engine 28.0.4 stehen im Terminal-Strom des Clients, den der Launcher nicht liest, und der Befehl endet dann mit `0`. **Engine-Abhängigkeit:** Die Meldung kommt je Engine in anderer Form. Der fehlende Benutzer endet auf der Engine 29.9.2 mit Status 1 und `Error response from daemon: unable to find user U: …` auf der Fehlerausgabe, auf der Engine 28.0.4 mit Status 126 und der Zeile `unable to find user U: …` ohne Präfix auf der Standardausgabe (mit `\r\n`). Messbedingungen: zwei Engines (29.9.2 und 28.0.4), Client 29.9.2, Linux `amd64`, ein Lauf je Fall; andere Engine-Versionen sind nicht zugesagt, und die Meldungen sind für den Client 29.9.2 gemessen.
- **Ohne `u-boot.yaml`.** Das Urteil nach dem Zustand setzt die `u-boot.yaml` des Projektordners voraus ([`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): Der Exit-Code `11` für den Startscript-Abbruch setzt die Statusdatei des Containers voraus). Ohne diese Datei, ohne Sandbox-Profil oder ohne die Fähigkeit liest `shell` keine Statusdatei und öffnet die Sitzung auch in einem Container, dessen Startscript abgebrochen ist; scheitert der Start selbst, endet der Befehl dann mit `12`.

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
      "description": "Unterkommando der gruppierten Hauptkommandos `template`, `config`, `egress` und `devcontainer`"
    },
    "status": {
      "type": "string",
      "enum": ["ok", "warn", "error"]
    },
    "command": {
      "type": "string",
      "enum": ["init", "add", "remove", "up", "down", "doctor", "logs", "generate", "config", "egress", "devcontainer", "template"]
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
    },
    {
      "if": {
        "properties": {
          "command": { "const": "egress" }
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
          "command": { "const": "devcontainer" }
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

Der Wert `egress` gehört zum Feld `command` für die `--json`-Ausgabe von `u-boot egress list` ([`LH-NFA-USE-004`](lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) bindet alle `--json`-Ausgaben an dieses Enum); der Befehl setzt `subcommand: "list"` und weist `--dry-run` und `--diff` ab ([`SPEC-046`](#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung)).

Der Wert `devcontainer` gehört zum Feld `command` für die `--json`-Ausgabe von `u-boot devcontainer up` und `u-boot devcontainer down`; der Befehl setzt `subcommand` auf `up` beziehungsweise `down` und weist `--dry-run` und `--diff` ab ([`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)). `u-boot devcontainer shell` weist `--json`, `--dry-run`, `--diff` und zusätzliche Argumente mit Exit-Code `2` ab; der Umschlag dieser Abweisung trägt `command: "devcontainer"` und `subcommand: "shell"`, und eine Erfolgsausgabe mit `--json` gibt es bei `shell` nicht.

Die Pflicht von `subcommand` für `template`, `config`, `egress` und `devcontainer` gilt für jede `--json`-Ausgabe ([`LH-NFA-USE-004.a`](#lh-nfa-use-004a--regeln-für---json-antworten)); `allOf` gibt sie für das Schema der Vorschau wieder, auch für `egress` und `devcontainer`, die keine Vorschau-Ausgabe erzeugen.

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
- [LH-FA-DEV-008](lastenheft.md#lh-fa-dev-008--egress-restriktion): `u-boot egress list`, `u-boot egress list --json`
- [LH-FA-DEV-010](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): `u-boot devcontainer up`, `u-boot devcontainer up --recreate`, `u-boot devcontainer down`, `u-boot devcontainer shell`
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

- `ARG GO_VERSION` – mit Default-Pin; Hebung ist Routine ohne separaten Spec-Eintrag.
- `ARG GOLANGCI_LINT_VERSION` – mit Default-Pin; gleiche Pin-Politik.
- `ARG COVERAGE_THRESHOLD` – mit Default `90` (Prozent) und Override-Pfad `make coverage-gate THRESHOLD=…`; das Bootstrap-Verhalten bei leerer Coverage-Eingabe steht in der Anforderung.

### SPEC-007 — Mindest-Ausschlüsse der `.dockerignore`

Mindestens auszuschließen:

- `.git`, `.gitignore`, `.github`
- IDE-Verzeichnisse: `.idea`, `.vscode`
- Agent-Verzeichnisse: `.claude`, `.codex`, `.agents`
- lokale Build-Artefakte und Caches (z. B. `dist/`, `coverage*`, `*.log`)

### SPEC-016 — Go-Toolchain: Mindestversion und Dockerfile-Pin

Mindest-Toolchain: Go 1.26 oder neuer (`go 1.26.0` in `go.mod`, analog Referenzprojekt `k-deskflight`); Default-Pin im Dockerfile als `ARG GO_VERSION`. Pin-Hebung ist Routine ohne separaten Spec-Eintrag.

### SPEC-018 — Standardwerte der Build- und Laufzeitumgebung

| Name | Wert | Quelle |
|---|---|---|
| Coverage-Schwellwert (`COVERAGE_THRESHOLD` / `THRESHOLD`) | 90 Prozent | Dockerfile-Build-Arg, Makefile-Variable |
| Polling-Intervall von `u-boot up` | 500 ms | Stabilisierungsschleife (`docker compose ps`) |
| Standard-Wartezeit von `u-boot up` | 60 s | `--timeout` |
| UID des Container-Benutzers | 1000 | `devcontainer.user.uid` |

### SPEC-019 — Default-Allowlist der Egress-Restriktion

Die wirksame Liste ist die Vereinigung aus Basis, Ergänzungen und Nutzerliste, sortiert und ohne Duplikate. Dieselbe Berechnung schreibt die Zeile `hosts=` des Startscripts `egress-init.sh` und liefert die Anzeige von `u-boot egress list` ([`SPEC-046`](#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung)); die dritte Spalte nennt den Herkunftsnamen, den die Anzeige je Host führt.

| Quelle | Hosts | Herkunft |
|---|---|---|
| immer | `github.com`, `api.github.com`, `codeload.github.com`, `raw.githubusercontent.com`, `objects.githubusercontent.com`, `deb.debian.org`, `security.debian.org` | `base` |
| `nestedRuntime: podman` oder `docker` | `registry-1.docker.io`, `auth.docker.io`, `production.cloudflare.docker.com`, `production.cloudfront.docker.com`, `ghcr.io` | `runtime` |
| `nestedRuntime: docker` | zusätzlich `download.docker.com` (Paketquelle bleibt im Image eingetragen) | `docker-apt` |
| Feature `node` | `registry.npmjs.org` | `feature:node` |
| Feature `go` | `proxy.golang.org`, `sum.golang.org`, `storage.googleapis.com` | `feature:go` |
| Feature `java` | `repo.maven.apache.org`, `repo1.maven.org` | `feature:java` |
| Feature `python` | `pypi.org`, `files.pythonhosted.org` | `feature:python` |
| Feature `sdkman` | `api.sdkman.io`, `broker.sdkman.io`, `repo.maven.apache.org`, `repo1.maven.org` | `feature:sdkman` |
| Feature `homebrew` | `formulae.brew.sh`, `ghcr.io`, `pkg-containers.githubusercontent.com` | `feature:homebrew` |
| Feature `claude-code` | `registry.npmjs.org`, `api.anthropic.com` | `feature:claude-code` |
| Feature `ollama` | `registry.ollama.ai`, `dd20bb891979d25aebc8bec07b2b3bbc.r2.cloudflarestorage.com` | `feature:ollama` |
| Clone-Quelle | Host von `devcontainer.sandbox.repository` beziehungsweise `origin` ([`SPEC-046`](#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung)) | `clone` |
| Nutzer | Einträge von `devcontainer.sandbox.egress.allow` (kleingeschriebene Hostnamen mit mindestens einem Punkt, ohne Schema, Port und Wildcard) | `user` |

Das Feature `sshd` ergänzt keinen Host. Die Zeile `claude-code` führt den Host der Anthropic-API und die npm-Registry (die ohnehin zum Feature `node` gehört, das `claude-code` braucht); die Hosts der Anmeldung sowie von Telemetrie, Updates, MCP und Plugins führt die Liste nicht, der Nutzer trägt sie über `devcontainer.sandbox.egress.allow` ein. Ein Feature mit `enabled: false` liefert keinen Host.

**Herkunft.** Die Anzeige führt sechs feste Namen: `base`, `runtime`, `docker-apt`, `feature:<Schlüssel>`, `clone` und `user`. Das Lastenheft nennt fünf Klassen (Basis, Container-Runtime, Feature, Clone-Quelle, Nutzerliste); `docker-apt` ist eine Verfeinerung der Klasse Container-Runtime und trägt allein den Host `download.docker.com`, nur bei `nestedRuntime: docker`. Ein Host, den mehrere Quellen liefern (etwa `ghcr.io` bei `podman` zusammen mit `homebrew`, oder `registry.npmjs.org` bei `node` zusammen mit `claude-code`), trägt jede seiner Herkünfte einmal, in der Reihenfolge `base`, `runtime`, `docker-apt`, Features nach Schlüssel, `clone`, `user`.

**Zeile `ollama`.** Die Liste führt den Host der Manifeste und den Host der Modell-Dateien (Blob-Speicher) des Modell-Pull; für den Namen des Blob-Hosts führt sie keinen breiteren Eintrag. `ollama.com` steht nicht in der Liste; Start, Pull und Lauf des Servers funktionieren ohne ihn ([`SPEC-047`](#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen)).

**Namen und Subdomains.** Mit dem Eintrag `example.org` sind `example.org` und `www.example.org` erreichbar, mit dem Eintrag `amazonaws.com` auch `s3.us-east-1.amazonaws.com` (zwei Ebenen darunter), und `example.com` ist es nicht; mit dem Eintrag `www.example.org` ist `www.example.org` erreichbar und `example.org` bei der Verbindung per Name nicht. Einen Eintrag der Form `*.example.org` gibt es nicht: `u-boot` weist jeden Eintrag mit `*` mit dem Exit-Code `10` ab, und die Meldung nennt den Namen ohne `*.`, der stattdessen einzutragen ist.

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
| `devcontainer.features.allowlist` | `devcontainer` Features auf Allowlist; im Sandbox-Profil kein Feature, das den Socket des Hosts einbindet (Fehler bei `docker-cli`), [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil), Einzelheiten [`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) |
| `devcontainer.features.drift` | `devcontainer` Features ohne Drift |
| `devcontainer.sandbox.runtime` | Sandbox: nested Podman (`/dev/fuse` auf dem Linux-Host, nicht im Container: Fehler bei Fehlen, Profil-Konsistenz), [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) |
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
- Der Katalog führt je Eintrag zwei Eigenschaften des Features in der Standardversion: ob es den Container-Runtime-Socket des Hosts einbindet (nur `docker-cli`) und welche Container-Optionen es beim Start mitbringt, soweit der Katalog sie führt (nur `go`: `--cap-add=SYS_PTRACE`, `--security-opt=seccomp=unconfined`; bei `docker-cli` führt er keine, denn das Sandbox-Profil weist den Eintrag ab). Bei keinem anderen Schlüssel trägt der Katalog eine von beiden, auch bei `kubectl-helm` nicht. Die Metadaten sind für `git`, `docker-cli`, `node`, `java`, `sdkman`, `go`, `kubectl-helm`, `python`, `sshd`, `homebrew` und `claude-code` gelesen; für `cpp` und `postgres-client` sind sie nicht gelesen, ihre Felder sind leer. Was das Sandbox-Profil daraus macht, steht in [`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen).
- Der Schlüssel `ollama` gehört zum Katalog, ist aber kein OCI-Feature, sondern ein Dockerfile-Block (zweiter Eintragstyp): Er hat weder Quelle noch Standardversion, erzeugt keinen Schlüssel unter `features` in `devcontainer.json`, und `u-boot doctor` meldet für ihn keinen Befund der Prüfungen `devcontainer.features.allowlist` und `devcontainer.features.drift`, auch wenn `devcontainer.json` keinen Schlüssel dafür enthält. Mit Sandbox-Profil trägt das erzeugte Dockerfile die Installation ([`SPEC-047`](#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen)); ohne Sandbox-Profil entsteht nichts. Eine `version` ohne `source` ist bei `ollama` ein Validierungsfehler (Exit-Code `10`); mit `source` ist der Schlüssel ein externes Feature wie jeder Katalogschlüssel mit eigener Quelle, und der Dockerfile-Block entfällt.
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
| 2026-10-18 | [`LH-FA-DEV-007.b`](#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman) beschreibt die Herkunft von Podman bei `nestedRuntime: podman`, gelesen gegen Dockerfile-Vorlage und Startscript: Release-Archiv von `podman-static` in Version 6.1.4 mit gepinnter Prüfsumme und gepinntem Signatur-Fingerprint statt der Debian-Pakete, aus Debian nur `uidmap` und `ca-certificates`, `registries.conf` mit Docker Hub als einziger Kurznamen-Registry, `vfs`-Fallback über ein Benutzer-Drop-in und `image_copy_tmp_dir` im Storage-Volume; keine Änderung der Default-Allowlist. |
| 2026-10-18 | [`LH-FA-DEV-007.a`](#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) und [`.b`](#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman) nennen die Container-Bedingung der Prüfung `devcontainer.sandbox.runtime`, gelesen gegen den Code: Im Container prüft `doctor` `/dev/fuse` nicht (Info), weil das Device das des Containers ist; die Konfigurationszweige und die Prüfung des AppArmor-Werts laufen unverändert. |
| 2026-10-19 | Egress-Anzeige festgelegt, deckungsgleich mit Lastenheft 0.16.0 und Code (`SPEC-046`): Berechnung (dieselbe wie die Zeile `hosts=` des Startscripts), Eingabe, Aktivierungszustand, Auflösung der Clone-Quelle samt Fehlerfall, Ausgabeformen, Wirkung von `--quiet`, Exit-Codes und Abgrenzung zum laufenden Container (ein Neustart behält die alte Liste, gemessen in einem Lauf). `SPEC-043` führt `egress list` unter den Abfragebefehlen, `SPEC-001` den Wert `egress` im Enum von `command`, `SPEC-024` einen Beispielaufruf. |
| 2026-10-19 | `SPEC-019` führt die Herkunft je Quelle (sechs Namen, `docker-apt` als Verfeinerung der Klasse Container-Runtime) und die Subdomain-Aussage mit Messstand: gepinnte und gemessene Fälle getrennt, der übergeordnete Name bei geteilter IP-Adresse als nicht geprüft ausgewiesen, der offene Upstream-Resolver ohne `sudo` als gemessen; `*`-Einträge werden mit Exit-Code `10` abgewiesen; zu `CNAME`-Zielen keine Aussage. |
| 2026-10-19 | Feature `ollama` festgelegt, deckungsgleich mit Lastenheft 0.15.0 und Code (`SPEC-047`): Dockerfile-Block mit gepinntem Archiv und Prüfsumme, Volume, kein Start, Bindung an Loopback bei ungesetztem `OLLAMA_HOST`, zwei Meldungen, Grenzen. `SPEC-023` führt den zweiten Eintragstyp, `SPEC-019` die Zeile `ollama` (zwei Hosts, Stabilität des Blob-Hosts ungemessen), `SPEC-043` die beiden Meldungen unter den entfallenden, `SPEC-044` begrenzt die Sandbox-Regel auf `sshd` und `claude-code`. |
| 2026-10-19 | `SPEC-045` führt im Absatz *Grenze* nach, dass ein Prozess ohne `sudo` die Statusdatei ersetzen kann (Verzeichnis gehört `vscode`, gemessen in einem Lauf). |
| 2026-10-19 | Spezifikation auf Zielbild: `SPEC-019`, `SPEC-045`, `SPEC-046`, `SPEC-047` und [`LH-FA-DEV-007.b`](#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman) tragen keine Messstände, Messumgebungen und Nicht-Messungen mehr; `SPEC-006` und `SPEC-016` keinen datierten Pin-Stand; der offene Upstream-Resolver und die `CNAME`-Ziele sind nicht mehr Gegenstand von `SPEC-019`. Zusagen bleiben, wo ein Test sie pinnt. |
| 2026-10-19 | `SPEC-046` führt den Hinweis von `generate devcontainer` nach, wenn ein bestehendes `egress-init.sh` ersetzt wird (Warnung, entfällt unter `--quiet` und `--dry-run`, `SPEC-043` führt ihn unter den entfallenden), und die Zusage, dass eine geänderte Liste einen Neuaufbau braucht, wieder in der Abgrenzung zum laufenden Container; die Subdomain-Aussage der Ausgabe nennt nur die gepinnten Fälle. |
| 2026-10-20 | Egress-Startscript festgelegt, gelesen gegen Template, Golden Cases und Integrationstests (`SPEC-048`): Upstream-DNS nur für den Benutzer des lokalen Resolvers samt Grenze, vier Abbruchfälle mit Exit-Code `11` und `failed` unter jeder Richtlinie, Wiederholungslauf in einer `nft -f`-Eingabe. |
| 2026-10-20 | `SPEC-045` nennt `failed` auch unter `warn` und nennt die Verkettung `egress && docker` nach der Semantik von `&&`; *Nicht ermittelbar* nennt keine Ursache aus der Verkettung. |
| 2026-10-20 | Die Degradationstabelle von [`LH-FA-DEV-007.a`](#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) führt den Abbruch aus anderem Grund als der fehlenden Capability. |
| 2026-10-20 | [`LH-FA-DEV-008.a`](#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability) grenzt den Abbruch von der Capability-Prüfung ab. |
| 2026-10-20 | `SPEC-048` auf die gepinnte Stärke begrenzt: Aussagen zum Wiederholungslauf, zur Ablehnung durch `nft`, zu den Abbruchfällen 3 und 4, zur Doctor-Meldung und zur IPv4-Beschränkung der Upstream-Regeln. |
| 2026-10-20 | `SPEC-048` und `SPEC-045` führen nur noch Bedingungen mit pinnendem Test: gestrichen sind der Weiterlauf bestehender Verbindungen, die IPv4-Beschränkung samt Set-Typ, die Pid-Datei als Ursache des Falls 4, das „erste“ der Abbruchfälle 3 und 4 und der Satz zum Teststand der Verkettung. |
| 2026-10-21 | Devcontainer ohne Editor festgelegt, deckungsgleich mit Lastenheft 0.18.0 und Code (`SPEC-049`): Aufruf, Voraussetzungen und Vorab-Prüfung der Engine, Hilfs-Image mit Pin-Form, Argumente des Hilfs-Containers, Spiegelpfad, Umgebung und Maskierung samt Grenzen, Abbildung eines Startscript-Abbruchs auf Exit-Code `11`, Ausweisung des Sockets, Ausgabe unter `--quiet` und `--json`, Exit-Codes. `SPEC-043`, `SPEC-045`, `SPEC-001`, `SPEC-024`, `SPEC-048` und [`LH-FA-DEV-007.a`](#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) nennen den Befehl in ihren Aufzählungen. Nicht durch einen Test belegt und deshalb nicht zugesagt sind der Exit-Code bei einem `failed` aus einem früheren Lauf und das Verhalten der Vorab-Prüfung ohne Compose-Plugin; nicht jede Bedingung ist durch eine Mutation des Codes als wirksam gesehen. |
| 2026-10-21 | [`LH-NFA-USE-004.a`](#lh-nfa-use-004a--regeln-für---json-antworten) und `SPEC-001` folgen dem Lastenheft 0.19.0: `subcommand` ist auch bei `egress` verpflichtend (Regel, Beschreibung des Feldes, Block in `allOf`); ein Absatz nennt die Herkunft der Pflicht aus dem Minimalkontrakt. |
| 2026-10-22 | Katalog-Features im Sandbox-Profil festgelegt, deckungsgleich mit Lastenheft 0.21.0 und Code (`SPEC-050`): Abweisung von `docker-cli` (Exit-Code `10`, nichts geschrieben, Abhilfe in der Meldung, `config set` bleibt möglich), Fehler in `doctor`, Warnungen zu den Lockerungen von `go` und je zu einem Feature aus einer externen Quelle samt Einstufung unter `--quiet` (`SPEC-043`); Eigenschaften je Katalogeintrag (`SPEC-023`); Teilaussage (3) der Ausweisung von `up` (`SPEC-049`); Beschreibung der Prüfung `devcontainer.features.allowlist` (`SPEC-017`). Nicht durch einen Test belegt und deshalb nicht zugesagt sind der Zustand eines vorhandenen Containers nach der Änderung und die Eigenschaften eines Features aus einer externen Quelle. |
| 2026-10-22 | Devcontainer-Shell und erneutes `up` festgelegt, deckungsgleich mit Lastenheft 0.21.0 und Code (`SPEC-049`): Unterabschnitt `shell` (Startphase wie `up` ohne `--recreate`, Probe und Sitzung, Umgebung aus den Einträgen `"<NAME>": "${localEnv:<NAME>}"`, Terminal-Erkennung, Ende der Sitzung und Signale, Art der Shell, Ausgabe auf der Standardfehlerausgabe, `--quiet`, Exit-Code `12`, Grenzen der Erkennung von Fehlern des Clients, `shell` ohne `u-boot.yaml`); Urteil von `up` nach dem Zustand des Containers (Exit-Code `11` bei `failed` auch nach einem gelungenen Lauf und aus einem früheren Start), Zeile `already running` und Feld `alreadyRunning`, Lebenszyklus-Befehle bei `up` gegen einen laufenden Container. `SPEC-043`, `SPEC-045`, `SPEC-048`, `SPEC-001`, `SPEC-024` und [`LH-FA-DEV-007.a`](#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) nennen `shell` in ihren Aufzählungen. Nicht durch einen Test belegt und deshalb nicht zugesagt sind unter anderem der Exit-Code eines Abbruchs, Strg-C ohne Terminal, das Schließen des Terminals, die Erkennung von Fehlern des Clients mit Terminal außer „Container läuft nicht“ und „Standardeingabe ist kein Terminal“, das `postAttachCommand` bei `shell` und die Gleichheit des gelesenen mit dem gestarteten Container bei mehreren Containern am selben Label. |
