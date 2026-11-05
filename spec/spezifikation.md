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
- `u-boot generate devcontainer` erzeugt die `source:`-Überschreibung eines aktivierten Features nur, wenn die Quelle in `devcontainer.featureSources.allow` oder in `--allow-external-feature-sources` desselben Aufrufs steht; sonst Exit-Code `10`, Code [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features), nichts geschrieben ([`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen), Absatz *Freigabe der Quelle*).

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

- Bei `podman`: rootless Podman, ein `docker`-Kompatibilitäts-Alias, Zugriff auf `/dev/fuse` und `/dev/net/tun` (Netzwerk der verschachtelten Container über `pasta`; ohne das Device scheitert jedes `RUN`; beide Devices müssen auf dem Container-Host vorhanden sein, fehlt eines, lehnt die Container-Engine den Start ab (`error gathering device information while adding custom device`), das Startscript läuft dann nicht; `doctor` prüft `/dev/net/tun` nicht und weist ein fehlendes `/dev/fuse` nur aus, soweit der Host es bestimmen kann, nicht aus einem Container), subuid-/subgid-Bereiche für den Container-Benutzer und ein Volume für den Storage. `runArgs` enthält neben der Zeile der Umgebungsdatei ([`SPEC-055`](#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung)) sechs Optionen: `--cap-add=SYS_ADMIN`, `--security-opt=seccomp=unconfined`, `--security-opt=apparmor=unconfined`, `--security-opt=systempaths=unconfined`, `--device=/dev/fuse`, `--device=/dev/net/tun`; jede weist die Befehlsausgabe einzeln aus. Podman kommt nicht aus der Distribution, sondern aus dem statisch gelinkten Release-Archiv `podman-linux-<arch>.tar.gz` des Projekts `podman-static` (Version 6.1.4; `amd64` und `arm64`; enthält `podman`, `crun`, `runc`, `conmon`, `netavark`, `aardvark-dns`, `passt`/`pasta`, `fuse-overlayfs` und `catatonit` unter `/usr/local`). Version, SHA-256 je Architektur und der Primär-Fingerprint des Signaturschlüssels (`0CCF102C4F95D89E583FF1D4F8B5AF50344BB503`) stehen als `ARG` an einer Stelle im erzeugten Dockerfile; der Image-Build bricht vor dem Entpacken ab, wenn die Prüfsumme nicht stimmt, der Schlüssel (von `keys.openpgp.org`, nach Fingerprint) nicht importierbar ist oder die Signatur des Archivs nicht vom gepinnten Primärschlüssel stammt (`VALIDSIG`). Das Archiv wird entpackt, nicht ausgeführt; Archiv und Schlüssel lädt der Image-Build von `github.com` und `objects.githubusercontent.com` (Archiv, Signatur) und `keys.openpgp.org` (Schlüssel); der Build läuft auf dem Container-Host ohne die Egress-Restriktion, braucht diese Hosts also offen, und im laufenden Container ist keiner davon nötig ([`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion) bleibt unverändert). Aus Debian installiert das Dockerfile nur `uidmap` und `ca-certificates` ohne empfohlene Pakete (`podman`, `fuse-overlayfs` und `passt` kommen nicht aus Debian); `nftables` kommt nur mit der Egress-Restriktion ins Image. Der Alias ist ein Symlink `/usr/local/bin/docker` → `/usr/local/bin/podman`. Das Archiv liefert unter `/etc/containers/` `containers.conf`, `policy.json`, `registries.conf`, `seccomp.json`, `storage.conf` und `storage.conf.d/00-podman-static.conf`; sie gelten mit einer Änderung, und das `README.md` des Archivs wird nach dem Entpacken entfernt. `registries.conf` setzt `unqualified-search-registries = ['docker.io']`. Daraus folgt für Builds: `podman build` und `docker build` sind derselbe Aufruf mit der in Podman eingebauten Buildah-Bibliothek; das Kommando `buildah`, `docker buildx`, `docker compose` und BuildKit-Frontends (`# syntax=`) sind nicht Teil der Ausstattung. Ein Dockerfile mit `FROM`, `RUN` und `COPY` baut kalt, aus dem Layer-Cache und mit `--no-cache`, mit und ohne Egress-Restriktion (Default-Allowlist). Der Kurzname `alpine` löst über Docker Hub auf (`docker.io/library/alpine`), `mycompany/img` wird nur dort gesucht; Dockerfiles verwenden voll qualifizierte Image-Namen. Fehlt `/dev/fuse` im Container (der Container wurde ohne das Device gestartet), schreibt das Startscript ein Benutzer-Drop-in (`~/.config/containers/storage.conf.d/99-u-boot-vfs.conf` mit `driver = "vfs"` und leerem `mount_program`, das die Vorgabe der mitgelieferten Konfiguration überstimmt) und setzt `image_copy_tmp_dir` (`~/.config/containers/containers.conf`) auf ein Verzeichnis im Storage-Volume; `podman run`, `podman build` und `docker build` laufen dann mit `vfs`-Storage. Der Beleg sind Integrationstests; sie laufen nicht auf Hosts, die verschachtelte User-Namespaces blockieren (`kernel.apparmor_restrict_unprivileged_userns=1`, etwa Ubuntu 24.04 mit Standardeinstellung), dort endet das Startscript mit Exit-Code `11`.
- Bei `docker`: Docker-Daemon, `docker`-CLI sowie Buildx- und Compose-Plugin aus dem Docker-APT-Repository (Schlüssel per Fingerprint geprüft); der Container-Benutzer gehört der Gruppe `docker` an, der Daemon startet beim Containerstart (`postStartCommand`) und gilt nach spätestens 60 Sekunden als nicht betriebsbereit; `runArgs` enthält neben der Zeile der Umgebungsdatei ([`SPEC-055`](#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung)) `--privileged` und, nur mit Egress-Restriktion, `--cap-add=NET_ADMIN`, sonst keine Fähigkeitsoptionen; der Daemon-Speicher (`/var/lib/docker`) liegt in einem Volume `<projekt>-docker-${devcontainerId}`; kein Socket-Mount.

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

- **Der Bau-Pfad.** Das Repo trägt einen Bau-Pfad: die Rezepte von `scripts/bau.sh` laufen in der Devcontainer-Umgebung (`BUILD_MODE=local` ist der einzige Modus-Wert, eine zweite Modus-Entscheidung entfällt). Die Gatter tragen ihre Ebenen je Messbasis: `lint`/`test`/`coverage`/`docs-check`/`baseline-verify`/`ablage-pruefung` laufen im Devcontainer-Container — dessen verschachtelter Daemon trägt die Container-Züge der Dokumentation (Messung M2/M3) — und `test-docker` läuft an dem Docker, der den Devcontainer-Container erzeugt (Sockelform; der verschachtelte Daemon trägt die docker-getaggten Tests nicht — Messung M4). `modus-pruefung` entfällt; ein Aufruf mit dem entfallenen Schalter (`BUILD_MODE=docker`) ist ein unbekanntes Verhalten und endet mit dem `$(error)`-Wächter der geöffneten Guard-Region (Probe Exit 2).
- **Voraussetzung am Host:** Docker Engine und `make` (Bestand). `make` ist ein bewusster Carveout zu [`LH-NFA-PORT-002`](lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten) (weit verbreitet, einzige zusätzliche Host-Abhängigkeit neben Docker); der Host-Vertrag senkt sich mit dem Bau-Pfad nicht. Eine Go-Toolchain am Host darf für Standard-Targets nicht vorausgesetzt werden; die Bau-Nutzung läuft über den Devcontainer des Repos, dessen Image (`.devcontainer/Dockerfile`) die gepinnte Toolchain trägt, und die CI baut dieses Image je `gates`-Job-Beginn (`docker build -f .devcontainer/Dockerfile .devcontainer`) und ruft `make gates` in einem privilegierten `docker run` des Images auf (Muster je Messbasis M3/M5).
- Die Befehle von `compile`, `lint`, `test` und `coverage` stehen in `scripts/bau.sh`, das die Targets des Makefiles (`compile`, `lint`, `test`, `coverage-gate`) und die gleichnamigen Stages des Multi-Stage-Dockerfiles rufen; die Stages bleiben je per `docker build --target <stage>` einzeln baubar ([`SPEC-029`](#spec-029--mindestumfang-des-multi-stage-dockerfiles)). Wer Container startet, braucht einen Docker-Daemon: im Devcontainer trägt ihn der verschachtelte Daemon, an Host und CI der Docker der Umgebung. Die Träger-Ebene je Gatter mit Beleg:

  | Gatter | Ebene (je Messbasis) | Beleg |
  |---|---|---|
  | `lint`, `test`, `coverage-gate`, `docs-check`, `baseline-verify`, `ablage-pruefung`, `tool-pruefung` (die `make gates`-Kette) | Devcontainer-Container; der verschachtelte Daemon trägt den Container-Zug von `docs-check` | `make gates` im Container, Messung M3 ([`roh-gates.txt`](../docs/messungen/slice-v1-bau-nur-devcontainer/roh-gates.txt)); der CI-Lauf je Job (`ci.yml::gates` baut das Image und ruft `make gates` im Container), M5-Analogie der [`roh-analogie.txt`](../docs/messungen/slice-v1-bau-nur-devcontainer/roh-analogie.txt) |
  | `test-docker` | Sockelform an dem Docker, der den Devcontainer-Container erzeugt | der CI-Lauf je Job (`integration.yml::integration-docker` am Runner-Docker, Socket-Form), M5-Analogie der [`roh-analogie.txt`](../docs/messungen/slice-v1-bau-nur-devcontainer/roh-analogie.txt); Grenze: Messung M4 ([`roh-test-docker.txt`](../docs/messungen/slice-v1-bau-nur-devcontainer/roh-test-docker.txt)) — im verschachtelten Daemon 12 Fehlschläge von 1771, die Tests dort nicht als Nachweis zugesagt |
  | `govulncheck`, `image-scan`, `build`, `run`, `build-binaries`, `packages` am Docker der Umgebung | Docker-Container am aufrufenden Docker | `docker run`- und `docker build`-Form der Targets; die CI-Jobs (`ci.yml::security-gates`/`ci.yml::image-scan`, [`publish.yml`](../.github/workflows/publish.yml)) |
  | `modus-pruefung` | entfallen | Messung M2 ([`roh-dind.txt`](../docs/messungen/slice-v1-bau-nur-devcontainer/roh-dind.txt)): die Prüfung lief im verbleibenden Stand OK; sie entfällt, weil die Frage, die sie trug (die zwei Modi), entfällt |

- Release-Artefakte (Binaries, Pakete, Images) entstehen in der Release-Pipeline am Runtime-Image-Pfad (`build`, `build-binaries`, `packages`, je Docker-Container).
- Nachweis ist `make gates` im Devcontainer-Container. Ein Teillauf (`PKG` oder `RUN` bei `make test`, nur auf der Kommandozeile gültig) ist kein Nachweis; `gates`, `record-gates`, `ci` und `fullbuild` lehnen ihn ab.
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

**`devcontainer`** arbeitet in zwei Phasen: Zuerst werden alle Dateien klassifiziert, ohne zu schreiben; nur wenn keine Datei den Abbruch-Zustand hat, wird geschrieben. Ein Konflikt hinterlässt also keinen Teilzustand. `forwardPorts` stammt aus den Container-Ports der aktiven Services in `compose.yaml` ([`LH-FA-DEV-005`](lastenheft.md#lh-fa-dev-005--ports)) und teilt sich die Quelle mit der `doctor`-Prüfung `devcontainer.forwardPorts.consistency`. Ist `compose.yaml` nicht parsebar, gilt Exit-Code `10`. Die Freigabe externer Feature-Quellen ([`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features)) wird erst nach erfolgreicher Erzeugung in `u-boot.yaml` übernommen. Die Umgebungsdatei des Devcontainers gehört zum Plan, ist aber keine verwaltete Datei: Eine vorhandene wird weder gelesen noch geändert ([`SPEC-055`](#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung)).

### SPEC-040 — Konfigurationspfade und Schreibpipeline von `u-boot config`

Verfeinert [`LH-FA-CONF-005`](lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern) und [`LH-FA-CONF-002`](lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration).

**Erlaubte Pfade** (jeder andere Pfad: fachlicher Validierungsfehler, Exit-Code `10`):

| Pfad | Lesen | Schreiben | Werte |
|---|---|---|---|
| `project.name` | ja | ja | gültiger Projektname ([`SPEC-004`](#spec-004--projektnamen-muster)) |
| `devcontainer.enabled` | ja | ja | Wahrheitswert |
| `devcontainer.containerEnv` | nein (der Basis-Pfad ist keine Pfadform — auch `config get` liest ihn nicht; die Map bleibt via `config` und Anzeige lesbar) | ja (nur Unterpfade `devcontainer.containerEnv.<name>`, je Name ein Wert) | Map von Namen zu String-Werten; jeder Name ein gültiger Dev-Container-Umgebungsname, jeder Wert ein String; der Basis-Pfad selbst ist nicht schreibbar (kein Ganz-Map-Set) |
| `devcontainer.featureSources.allow` | ja | ja (Liste) | kommagetrennte Quellen |
| `devcontainer.features.<feature>.enabled` | ja | ja | Wahrheitswert |
| `devcontainer.features.<feature>.source` | ja | ja | gültige Quelle, die wörtlich in `devcontainer.featureSources.allow` steht |
| `devcontainer.features.<feature>.packages` | ja | ja (Liste) | kommagetrennte Formeln, `<formel>` oder `<user>/<tap>/<formel>`; nur am Schlüssel `brew` ohne `source` (sonst `ErrConfigValueInvalid`, Exit-Code `10`); bei aktivem `brew` nur Formeln freigegebener Taps ([`SPEC-057`](#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen)) |
| `devcontainer.brewTaps.allow` | ja | ja (Liste) | kommagetrennte Taps `<user>/<tap>`, unabhängig von `devcontainer.featureSources.allow` ([`SPEC-057`](#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen)) |
| `devcontainer.features.<feature>.version` | ja | ja | nicht-leerer String; bei einem Feature ohne OCI-Quelle nur, wenn der Eintrag eine `source` trägt oder einer der Schlüssel der gepinnten Werkzeugkette ist (`go-toolchain`, `golangci-lint`, `java-toolchain` — ihre Übersteuerung prüft `generate devcontainer` gegen die Quellen-Angabe, [`SPEC-054`](#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)) (sonst `ErrConfigValueInvalid`, Exit-Code `10`) |
| `devcontainer.user.uid` | ja | ja | Ganzzahl `1`..`65535` |
| `devcontainer.profile` | ja | ja | `default` \| `sandbox` |
| `devcontainer.sandbox.nestedRuntime` | ja | ja | `none` \| `podman` \| `docker` |
| `devcontainer.sandbox.onUnavailable` | ja | ja | `warn` \| `fail` |
| `devcontainer.sandbox.repository` | ja | ja | Clone-URL aus den Zeichen `A–Z a–z 0–9 . _ ~ : / @ % + = -`, nicht mit `-` beginnend, ohne eingebettete Zugangsdaten ([`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)) |
| `devcontainer.sandbox.egress.enabled` | ja | ja | Wahrheitswert |
| `devcontainer.sandbox.egress.allow` | ja | ja (Liste) | kommagetrennte Hostnamen |
| `devcontainer.sandbox.hostGitToken` | ja | ja | Wahrheitswert; ohne Sandbox-Profil ohne Wirkung ([`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)) |
| `devcontainer.resources.memory` | ja | ja | positive Ganzzahl mit optionaler Einheit `b`, `k`, `m` oder `g` (klein), Zeichenmenge `0-9bkmg`, nicht `0` ([`SPEC-053`](#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs)) |
| `devcontainer.resources.cpus` | ja | ja | positive Dezimalzahl mit `.`, Zeichenmenge `0-9` und `.`, nicht `0` ([`SPEC-053`](#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs)) |
| `services.<svc>.enabled` | ja | **nein** | — |

`services.<svc>.enabled` ist nur lesbar: Ein Schreibversuch endet mit Exit-Code `10` und dem Hinweis auf `u-boot add <svc>`, weil das bloße Umschalten die Zustandsregeln aus [`LH-FA-ADD-005.a`](#lh-fa-add-005a--zustandsregeln-für-registrierte-und-aktive-services) umginge. `<svc>` und `<feature>` werden auf Namensformat geprüft. Wahrheitswerte nehmen die Schreibweisen `true`, `false`, `1`, `0`, `T`, `F` u. ä. an und werden kanonisch als `true`/`false` geführt.

**Schreibpipeline `config set`** — jede Stufe läuft vor dem Schreiben; schlägt eine fehl, bleibt `u-boot.yaml` byte-identisch:

1. Pfad nicht schreibbar → Abbruch (vor jedem Dateizugriff).
2. Wert in den erwarteten Typ überführen und prüfen (Wertetabelle oben).
3. Entspricht der kanonische neue Wert dem bestehenden, wird nichts geschrieben (No-Op, alter und neuer Wert gleich).
4. Wert im Dokument ersetzen, Kommentare und Formatierung bleiben erhalten.
5. Das Ergebnis muss sich wieder als Projektkonfiguration einlesen lassen (sonst Schemafehler, Exit-Code `10`).
6. Fachliche Prüfung des Ergebnisses: Der geschriebene Wert muss gebunden und gültig sein; bei `.source` muss die Quelle wörtlich (Groß-/Kleinschreibung und abschließender Schrägstrich zählen) in `devcontainer.featureSources.allow` stehen, eine Form, die die Dev-Container-CLI nicht verarbeitet (gemessen 2026-11-03: `oci://…:<tag>` und `https://…:<tag>` unverarbeitbar) weist den Aufruf mit Exit-Code `10` ab, die Abweisungs-Meldung nennt die Abhilfe der Registry-Form (die Registry-Tarball-Form ist offen und bleibt beobachtet), `--yes` ersetzt diese Freigabe nicht; bei `.version` eines Features ohne OCI-Quelle (Dockerfile-Block oder Container-Option) muss der Eintrag eine `source` tragen oder einer der Schlüssel der gepinnten Werkzeugkette sein (`go-toolchain`, `golangci-lint`, `java-toolchain`; [`SPEC-054`](#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)) (dieselbe Regel wie beim Lesen durch `generate devcontainer`; bei mehreren Paaren in einem Aufruf steht das Paar der `source` vor dem Paar der `version`); bei `.enabled` von `brew` und `homebrew` darf der Schlüssel nicht aktiv werden, wenn der andere aktiv ist oder die Formel-Liste von `brew` einen nicht freigegebenen Tap trägt (Exit-Code `10`, [`SPEC-057`](#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen)).
7. Schreiben.

Die vier Listenpfade (`devcontainer.featureSources.allow`, `devcontainer.sandbox.egress.allow`, `devcontainer.features.brew.packages`, `devcontainer.brewTaps.allow`) hängen die neuen Einträge an die bestehende Liste an und entdoppeln sie; die Datei wird dabei neu serialisiert, Kommentare gehen verloren. `config set devcontainer.features.<feature>.enabled true` für ein Feature, das weder im eingebauten Katalog ([`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen)) steht noch eine `source` trägt, ist kein Fehler, erzeugt aber eine `warn`-Diagnose ([`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features)): Der Generator überspringt es, `doctor` meldet es. `config` ohne Argumente gibt `u-boot.yaml` unverändert (mit Kommentaren) aus.

**Projekt-Umgebungswerte (`devcontainer.containerEnv`).** Der Pfad trägt eine Map von Namen zu String-Werten; schreiben nur die Unterpfade `devcontainer.containerEnv.<name>` (je Name ein Wert), der Basis-Pfad selbst ist nicht schreibbar (Ganz-Map-Set ist kein erlaubter Pfad); `u-boot config set devcontainer.containerEnv.<name> <wert>` validiert den Namen als Dev-Container-Umgebungsname und den Wert als String und weist je fachlichen Fehler mit Exit-Code `10` ab, ehe nichts geschrieben wird (Wertetabelle oben). `u-boot generate devcontainer` trägt die Werte der Map in das verwaltete `containerEnv`-Objekt der erzeugten `.devcontainer/devcontainer.json` — neben den Generator-Werten der Schlüssel, ohne Duplikate, idempotent (eine Regeneration schreibt die Werte neu, ein entfernter Wert verschwindet); die Erzeugung meldet jede erzeugte Projekt-Wert-Zeile in der Ausweisung der Erzeugung. Setzt das Projekt mit aktivem Toolchain-Schlüssel (`go-toolchain` oder `golangci-lint`) einen der Cache-Werte (`GOCACHE`, `GOMODCACHE`), die der Schlüssel selbst trägt, beendet die Erzeugung mit demselben fachlichen Fehler (Exit-Code `10`, nichts geschrieben; Muster der Abweisungen von [`SPEC-054`](#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)); ohne den Schlüssel darf das Projekt beide setzen. Die Werte stehen nur in der `containerEnv`-Ebene der erzeugten Datei; der Dockerfile-Block und der handgeschriebene Rest tragen sie nicht ([`LH-FA-DEV-001`](lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), Absatz *Projekt-Umgebungswerte*).

Beleg-Tabelle (Muster [`SPEC-054`](#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv); die Test-Namen sind die des Umsetzungs-Diffs):

| Bedingung | Pinnender Test (Umsetzung) |
|---|---|
| Werte im verwalteten `containerEnv`, je Profil, ohne Duplikate, idempotent; ohne die Map byte-identische erzeugte Datei | `TestGenerateDevcontainer_ContainerEnv_InBothProfiles` (je Profil die drei Einträge sortiert, NoOp des zweiten Laufs) samt `…/golden default profile` (Byte-Golden), Mischfall `TestGenerateDevcontainer_ContainerEnv_MixedWithTheCacheValues` (Byte-Golden mit `go-toolchain`, `golangci-lint` allein unter dem `--sandbox`-Flag je genau ein Cache-Eintrag, `java-toolchain` ohne Cache-Raum), `TestGenerateDevcontainer_ContainerEnv_WithoutTheMapByteIdentical` (die zwei Bestands-Golden ohne Map unverändert) und `TestGenerateDevcontainer_ContainerEnv_RemovedEntryDisappears` (Regeneration schreibt neu, der entfernte Eintrag verschwindet, dritter Lauf NoOp) (`internal/hexagon/application/generate_containerenv_test.go`) |
| Validation Name/Wert je Exit-Code `10`, nichts geschrieben (auch `config set`) | `TestConfigSet_ContainerEnv_InvalidNameExits10` (je Fehlername der Umgebungs-Namensregel, `u-boot.yaml` byte-identisch), `TestConfigSet_ContainerEnv_HappyPathAndNoOp` (Wert eine Zeichenkette mit leerem Wert, NoOp bei gleichem Wert, Kommentare bleiben), `TestConfigSetMany_ContainerEnv_Entry`, `TestConfigGet_ContainerEnv_Extraction` (Eintrag verbatim, fehlender Eintrag NotSet) und am Nutzerpfad `TestRun_ConfigSet_ContainerEnvValueAndNames` (`cmd/uboot/containerenv_test.go`) |
| Kollision mit `GOCACHE`/`GOMODCACHE` bei aktivem Toolchain-Schlüssel → Exit `10`, nichts geschrieben; ohne den Schlüssel erlaubt | `TestGenerateDevcontainer_ContainerEnv_CollisionExits10` (je Schlüsselkombination: einzelner Name je Schlüssel, beide Namen bei beiden Schlüsseln, `--sandbox`-Flag, ohne die Schlüssel erlaubt, neben `java-toolchain` erlaubt, `java`-Feature neben dem kollidierenden Go-Schlüssel hebt die Abweisung nicht) und am Nutzerpfad mit `--dry-run --json` `TestRun_Generate_ContainerEnvCollisionExits10` (Code [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features), nichts geschrieben) |
| `config set devcontainer.containerEnv` (Basis-Pfad) endet mit Exit-Code `10` (kein Ganz-Map-Set, Wertetabelle oben) | `TestNewConfigPath_ContainerEnv_EntryParsesTheName` (`internal/hexagon/domain/configpath_test.go`: der Basis-Pfad und der leere Name sind keine Pfad-Formen) und `TestRun_ConfigSet_ContainerEnvValueAndNames` (Basis-Pfad je Lesen und Schreiben mit Exit `10` am Nutzerpfad) |
| Die Ausweisung der Erzeugung nennt jede erzeugte Projekt-Wert-Zeile | `TestGenerateDevcontainer_ContainerEnv_Notices` (je Eintrag eine gewöhnliche Warnung mit Code [`LH-FA-DEV-001`](lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), ohne `SecurityRelaxation`/`DegradationFallback`, ohne die Map keine) und `TestRun_Generate_ContainerEnvNotices` (am Nutzerpfad, unter `--quiet` keine) |
| Werte im Container-Prozess sichtbar (`env`/`go env`), nicht im Dockerfile-Block oder im handgeschriebenen Rest | Docker-Test `TestE2E_ContainerEnv_ValuesReachTheContainerProcess` (`internal/e2e/containerenv_docker_test.go`, `make test-docker`): `printenv` im gestarteten Container trägt die Werte, die `docker run`-Gegenprobe des gebauten Images sieht keinen der Werte, das erzeugte Dockerfile trägt keinen der Namen; die Konvertierung des eigenen Devcontainers pinnt `TestOwnDevcontainer_ProjectValuesUnderContainerEnv` |
| Schritt 6: eine `.source`-Form, die die Dev-Container-CLI nicht verarbeitet (gemessen 2026-11-03: `oci://…:<tag>`, `https://…:<tag>` unverarbeitbar), weist den Aufruf mit Exit-Code `10` ab, ehe nichts geschrieben wird — auch wenn der Wert byte-gleich in `devcontainer.featureSources.allow` steht; die Abweisungs-Meldung nennt die Abhilfe (die verarbeitbare Registry-Form), `--yes` ersetzt sie nicht | `TestConfigSet_FeatureSource_UnprocessableFormRejected` (beide gemessenen Formen, `u-boot.yaml` byte-identisch, Abhilfe-Satz, Form-Aussage statt des Freigabe-Hinweises; die Anwendungsschicht der Schreib-Pipeline) |
| Schritt 6: die Abhilfe ist ausführbar — die Registry-Referenz (host/pfad ohne URL-Schema) als `.source` mit Liste-Eintrag wird angenommen; der Schlüssel rendert als `<source>:<version>` | `TestConfigSet_FeatureSource_RegistryReferenceAccepted`, `TestGenerateDevcontainer_SourceFormRegistryReferenceRenders` |

Bezieht sich auf: [`LH-FA-DEV-001`](lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes).

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
- Vier Ausnahmen bleiben sichtbar: (1) die Ausweisung von Sicherheitslockerungen des Sandbox-Profils und eines im Container aktivierten Login-Dienstes ([`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion)) in der Warnungsausgabe von `generate devcontainer` (die Lockerungen des Sandbox-Profils und den Login-Dienst weist nur dieser Befehl aus; zu ihnen gehören die Lockerungen eines aktiven eingebauten Features, der auf Wunsch des Nutzers erzeugte Netzanschluss an den Host (`ollama-host`, [`SPEC-052`](#spec-052--feature-ollama-host-im-devcontainer-eintragstyp-container-option-meldung-und-grenzen)) und die Warnung je Feature aus einer externen Quelle, [`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen)), die Ausweisung der Einbindung des Container-Runtime-Sockets in den Hilfs-Container in der Ausgabe von `devcontainer up` und `devcontainer shell` (bei `shell` vor dem Beginn der Sitzung auf der Standardfehlerausgabe; [`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)) sowie die Warn-Einträge von `doctor`, die eine Sicherheitseigenschaft des Sandbox-Profils betreffen: die Quelle der Zugangsdaten ([`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)) und die Egress-Restriktion ([`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion)), darunter der Eintrag, dass sie ohne Sandbox-Profil nicht greift; (2) die Ausweisung der Degradations-Fallbacks (Container ohne Daemon, `vfs`-Storage, Wegfall der Egress-Restriktion) in `doctor` (Einträge `devcontainer.sandbox.fallback.*`) und in der Warnungsausgabe von `generate devcontainer`, `init`, `devcontainer up` und `devcontainer shell` (bei `shell` vor dem Beginn der Sitzung auf der Standardfehlerausgabe) nach [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung); ein Zustand, den der Host nicht lesen kann, ist „nicht ermittelbar“ und kein Fallback, und sein Eintrag entfällt unter `--quiet`; (3) der Compose-Fortschritt von `up` ([`LH-NFA-PERF-002`](lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker)); (4) die Pfad-Zusammenfassung, die `init --force` vor dem Überschreiben ausgibt ([`LH-FA-CLI-005A`](lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-INIT-005`](lastenheft.md#lh-fa-init-005--überschreibschutz)). Jede andere Warnung und jeder andere Hinweis entfällt.
- Die Einstufung hängt am Kennzeichen des einzelnen Eintrags (`SecurityRelaxation`, `DegradationFallback`), nicht an der Kennung der Prüfung oder Anforderung; dieselbe Kennung trägt gekennzeichnete und ungekennzeichnete Einträge. Im Einzelnen:
  - Gekennzeichnet sind in der Warnungsausgabe von `generate devcontainer` die Lockerungen von `nestedRuntime` und Egress-Restriktion (je Option eine Meldung), die Meldung, dass `egress.enabled` ohne Sandbox-Profil keine Wirkung hat (es entsteht keine Restriktion), der Login-Dienst, die Lockerungen eines aktiven eingebauten Features (`go`: je Option eine Meldung), die eine Meldung zum Netzanschluss an den Host im Sandbox-Profil (`ollama-host`, [`SPEC-052`](#spec-052--feature-ollama-host-im-devcontainer-eintragstyp-container-option-meldung-und-grenzen)) und die Warnung je Feature aus einer externen Quelle ([`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen)); in `doctor` die Warn-Einträge von `devcontainer.sandbox.credentials` (Klartext-Zugangsdaten, https-Clone ohne Token-Quelle, fehlende Referenz der Token-Variable eines aktiven Features: `claude-code`, `gh` oder `glab`) und von `devcontainer.sandbox.egress` (ohne Sandbox-Profil, zusammen mit `nestedRuntime: docker`); in der Ausgabe von `devcontainer up` und `devcontainer shell` die Ausweisung der Einbindung des Container-Runtime-Sockets in den Hilfs-Container ([`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)).
  - Mit `DegradationFallback` gekennzeichnet sind der `warn`-Eintrag je `devcontainer.sandbox.fallback.*` bei `result=fallback` und die entsprechende Warnung von `generate devcontainer`, `init`, `devcontainer up` und `devcontainer shell`; ungekennzeichnet bleiben die Einträge dieser Prüfungen bei `ok`, „nicht ermittelbar“ (`info`) und ungültigem Dokument (`warn`), der Eintrag bei `failed` ist ein Fehler und bleibt unter `--quiet` als Fehler sichtbar.
  - `nestedRuntime: docker` hat in `doctor` die Stufe OK (kein Warn-Eintrag; `--privileged` weist die Befehlsausgabe aus): Die Sicherheitslockerung des Docker-Daemons erscheint nur in der Befehlsausgabe.
  - Ungekennzeichnet sind die Meldungen „has no effect“ zu `nestedRuntime` in der Warnungsausgabe von `generate devcontainer` (Konfigurationsmeldung), der `doctor`-Eintrag, dass `nestedRuntime` ohne Sandbox-Profil keine Wirkung hat, die Drittanbieter-Warnung und der Hinweis auf das Feature `node` ([`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)), die beiden Meldungen zum Feature `ollama` (Lieferkette; keine Wirkung ohne Sandbox-Profil, [`SPEC-047`](#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen)), die Meldung zur Lieferkette je aktivem Feature `gh` und `glab` ([`SPEC-051`](#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten)), die Meldung zur Lieferkette je aktivem Schlüssel `go-toolchain`, `golangci-lint` und `java-toolchain` ([`SPEC-054`](#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)), die Meldung zu `ollama-host` im Standardprofil ([`SPEC-052`](#spec-052--feature-ollama-host-im-devcontainer-eintragstyp-container-option-meldung-und-grenzen)), der Hinweis zum Neuaufbau bei geändertem `egress-init.sh` ([`SPEC-046`](#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung)) sowie die Warn-Einträge von `devcontainer.features.allowlist`, `devcontainer.features.drift` und `devcontainer.envfile` ([`SPEC-055`](#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung)); sie entfallen. Dieselbe Aussage zur Egress-Restriktion ohne Sandbox-Profil ist in `doctor` wie in der Befehlsausgabe gekennzeichnet, also sichtbar.
- `down --quiet` hält die Compose-Ausgabe zurück und gibt sie nur aus, wenn `down` fehlschlägt, weil sie die Fehlerdetails trägt, die der Fehler selbst nicht nennt.
- `devcontainer shell --quiet` hält den Fortschritt des Starts zurück und gibt ihn nur aus, wenn der Befehl fehlschlägt; `u-boot` gibt selbst nur die Ausnahmen 1 und 2 aus, auf der Standardfehlerausgabe und vor dem Beginn der Sitzung. Die Ausgabe der Sitzung bleibt unverändert, denn die Shell ist nicht Ausgabe von `u-boot` ([`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)).
- `devcontainer up --quiet` hält den Fortschritt zurück und gibt ihn nur aus, wenn der Befehl fehlschlägt; `devcontainer down --quiet` gibt nichts aus ([`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)).
- Mit `--json` ändert `--quiet` nichts: Der Umschlag und seine `diagnostics[]` bleiben vollständig, und die Compose-Ausgabe von `down` wird nicht zurückgehalten.

### SPEC-044 — Features `sshd` und `claude-code` im Devcontainer: Warnungen und Zugangsdaten

Verfeinert [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) (Login-Dienst), [`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten des Agenten und eines Forge-Werkzeugs), [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Drittanbieter-Feature) und berührt [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts: [`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion)). Der Katalog steht in [`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen), die Wirkung von `--quiet` in [`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet). Die Regeln zu `docker-cli`, `go` und zu Features aus externen Quellen im Sandbox-Profil stehen in [`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen). Die Features `gh` und `glab` (Forge-Werkzeuge, keine Agenten) stehen in [`SPEC-051`](#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten); diese Stelle führt für sie die Token-Variablen, den Verweis auf ihre Meldung und die Prüfung `devcontainer.sandbox.credentials`.

**Aktiv** heißt: Ein Feature-Eintrag mit `enabled: true` löst auf die Katalogquelle von `sshd` bzw. `claude-code` auf; verglichen wird die aufgelöste Quelle über alle aktivierten Einträge, nicht der Schlüssel: Auch ein Eintrag mit anderem Schlüssel und `source` gleich der Katalogquelle ist aktiv. Eine `version`-Überschreibung ändert das nicht. Ein Eintrag, dessen aufgelöste Quelle eine andere ist, ist kein aktives `sshd` oder `claude-code`: Ein anderes SSH-Feature oder ein anderer Agent über `source` bleibt unausgewiesen.

**Sandbox-Regel.** Beide Features sind mit und ohne Sandbox-Profil wählbar; das Profil gibt keine externe Quelle frei, und der Katalogschlüssel braucht keinen Eintrag in `devcontainer.featureSources.allow`. Der Generator erzeugt für beide Features den Schlüssel unter `features`; der Katalog führt für beide weder die Einbindung des Container-Runtime-Sockets des Hosts noch Container-Optionen ([`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen)). Das Feature `ollama` ist ein Dockerfile-Block, ergänzt mit Sandbox-Profil ein Volume unter `mounts`, ohne Egress-Restriktion und ohne `nestedRuntime` weder Optionen in `runArgs` noch einen Startbefehl ([`SPEC-047`](#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen)). Die Features `gh` und `glab` sind ebenfalls Dockerfile-Blöcke und wirken in beiden Profilen ([`SPEC-051`](#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten)). Das Feature `ollama-host` ergänzt in beiden Profilen eine Zeile in `runArgs`, im Sandbox-Profil hinter den Optionen von `nestedRuntime` und Egress-Restriktion ([`SPEC-052`](#spec-052--feature-ollama-host-im-devcontainer-eintragstyp-container-option-meldung-und-grenzen)).

**Zugangsdaten des Agenten und der Forge-Werkzeuge.** Mit Sandbox-Profil trägt `remoteEnv` neben `GIT_TOKEN` je aktivem Feature den Eintrag mit der Referenz auf die Umgebungsvariable des Hosts, nie ihren Wert, in dieser festen Reihenfolge: `"ANTHROPIC_API_KEY": "${localEnv:ANTHROPIC_API_KEY}"` und `"CLAUDE_CODE_OAUTH_TOKEN": "${localEnv:CLAUDE_CODE_OAUTH_TOKEN}"` bei aktivem `claude-code` (der API-Schlüssel und das OAuth-Token eines Abonnements sind Anmeldewege, die der Nutzer wechselt; der Generator rendert immer beide Referenzen), `"GH_TOKEN": "${localEnv:GH_TOKEN}"` bei aktivem `gh` und `"GITLAB_TOKEN": "${localEnv:GITLAB_TOKEN}"` bei aktivem `glab`. „Aktiv“ heißt für `gh` und `glab`, was [`SPEC-051`](#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) für sie festlegt. Ohne Sandbox-Profil erzeugt der Generator kein `remoteEnv`. Andere Variablen reicht er nicht durch, auch nicht für die Features `ollama` und `ollama-host`; ein Anmeldeweg über Login im Container wird weder eingerichtet noch geprüft, und die Hosts der Anmeldung stehen nicht in der Default-Allowlist ([`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion)).

**Warnungen der Befehlsausgabe** (Erzeuger: `generate devcontainer`; Level `warn`, Subjekt `.devcontainer/devcontainer.json`):

| Auslöser | Code | Sicherheitslockerung (`--quiet`) | Inhalt |
|---|---|---|---|
| `sshd` aktiv, mit Sandbox-Profil | [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) | ja, bleibt sichtbar | Login-Dienst `sshd` ist im Container aktiv (Port 2222; Passwort- und Schlüssel-Login, Root-Login erlaubt; die Konten bleiben gesperrt, bis ein Passwort gesetzt oder ein Schlüssel eingetragen ist); das Profil begrenzt den Schaden und ist keine Isolationsgrenze. Ohne Sandbox-Profil entsteht keine Warnung. |
| `homebrew` oder `claude-code` aktiv, ohne eigene `source` | [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | nein, entfällt | Eine Meldung je Feature: Feature, Herausgeber, Quelle mit Version; sein Installationsskript läuft beim Bau des Images. Beide Profile. |
| `gh` oder `glab` aktiv | [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | nein, entfällt | Eine Meldung je Feature zur Lieferkette des Release-Archivs; Inhalt und Grenzen in [`SPEC-051`](#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten). Beide Profile. |
| `claude-code` aktiv, `node` nicht aktiv | [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) | nein, entfällt | `claude-code` installiert über npm, das Basis-Image enthält kein Node.js; `node` ist zu aktivieren, sonst schlägt der Bau des Images fehl. `u-boot` aktiviert `node` nicht selbst. Beide Profile. |
| `ollama-host` aktiv | [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) im Sandbox-Profil, [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) im Standardprofil | im Sandbox-Profil ja, bleibt sichtbar; im Standardprofil nein, entfällt | Eine Meldung zum Netzanschluss an den Host; Inhalt und Grenzen in [`SPEC-052`](#spec-052--feature-ollama-host-im-devcontainer-eintragstyp-container-option-meldung-und-grenzen). Beide Profile. |

**Prüfung `devcontainer.sandbox.credentials`** (`u-boot doctor`; kein eigener Check-Code für Agenten, die Kennung steht in [`SPEC-017`](#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry)). Ohne Sandbox-Profil ist der Eintrag OK und die Prüfung übersprungen. Mit Sandbox-Profil meldet ein Lauf höchstens den ersten Befund in dieser Reihenfolge, jeweils `warn` mit dem Kennzeichen der Sicherheitseigenschaft nach [`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (es ändert weder Severity noch Exit-Code):

1. **Klartext in einer Projektdatei.** Gescannt werden `u-boot.yaml`, `compose.yaml`, `.env.example`, `.devcontainer/devcontainer.json`, `.devcontainer/Dockerfile`, `.devcontainer/sandbox-init.sh` und `.devcontainer/docker-init.sh`. Treffer sind: Zugangsdaten in einer http(s)-URL (Benutzer, optional Passwort, vor `@`); Git-Token-Präfixe `ghp_`, `gho_`, `ghu_`, `ghs_`, `ghr_` (mindestens 20 Zeichen danach), `github_pat_` und `glpat-` (je mindestens 20 Zeichen danach); die Zuweisung (`=` oder `:`) eines Werts von mindestens 12 Zeichen an `git_token`, `gh_token`, `github_token` oder `gitlab_token` (Schreibung und Unterstrich beliebig); ein Anthropic-Schlüssel `sk-ant-` mit mindestens 40 Zeichen aus `A–Z a–z 0–9 _ -` danach, dem kein solches Zeichen vorangeht; die Zuweisung (`=` oder `:`) eines Werts von mindestens 12 Zeichen an `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN` oder `CLAUDE_CODE_OAUTH_TOKEN`. Referenzen wie `${localEnv:GIT_TOKEN}` treffen nicht. Die Meldung nennt die Dateien, nie den Wert. **Grenze:** Der Zuweisungs-Ausdruck unterscheidet einen Platzhalter nicht von einem Wert; `ANTHROPIC_API_KEY=your-api-key-here` in `.env.example` löst die Warnung aus.
2. **https-Clone ohne Token-Quelle.** Die Clone-Quelle beginnt mit `https://` oder `http://`, und `devcontainer.json` enthält weder `${localEnv:GIT_TOKEN}` noch die Zeichenfolgen `/run/secrets` und `readonly` zusammen (Secret-Mount; die Kurzform `ro` der Mount-Option zählt nicht).
3. **Token-Referenz fehlt.** Ein Feature mit Token-Variable ist aktiv, und `devcontainer.json` enthält die Referenz `${localEnv:<NAME>}` seiner Variablen nicht (`ANTHROPIC_API_KEY` und `CLAUDE_CODE_OAUTH_TOKEN` bei `claude-code`, `GH_TOKEN` bei `gh`, `GITLAB_TOKEN` bei `glab`); eine nicht lesbare Datei gilt als ohne Referenz. Fehlen mehrere Referenzen, nennt die Meldung alle fehlenden Variablen. Bei `claude-code` verlangt die Prüfung beide Referenzen, auch wenn nur einer der beiden Anmeldewege genutzt wird; fehlt eine, warnt sie und nennt nur die fehlende. Ein abgeschaltetes Feature verlangt keine Referenz, ebenso kein Eintrag von `gh` oder `glab` mit `source` ([`SPEC-051`](#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) zu „aktiv“).

Ohne Befund ist der Eintrag OK; bei aktivem Feature mit Token-Variable nennt die Meldung die Variablen und dass nur die Referenz durchgereicht wird.

### SPEC-050 — Katalog-Features im Sandbox-Profil: Abweisung von `docker-cli`, Lockerungen von `go`, Features aus externen Quellen

Verfeinert [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) (kein Container-Runtime-Socket des Hosts, keine zusätzlichen Capabilities, Ausweisung jeder Sicherheitslockerung) und [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) (Features, externe Quellen) und berührt [`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (`--quiet`, Ausnahme 1) sowie [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes `10` und `11`). Die Eigenschaften je Katalogeintrag stehen in [`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen), die Wirkung von `--quiet` in [`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), die Warnungen zu `sshd` und `claude-code` in [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten).

**Aktiv** heißt wie in [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten): Ein Feature-Eintrag mit `enabled: true` löst auf die Katalogquelle von `docker-cli` bzw. `go` auf; verglichen wird die aufgelöste Quelle über alle aktivierten Einträge. Eine `version`-Überschreibung ändert das nicht. Ein Eintrag mit `source:`-Überschreibung ist ein Feature aus einer externen Quelle, auch wenn sein Schlüssel ein Katalogschlüssel ist (`go` mit `source:` ist kein aktives `go`). Ein Eintrag mit `enabled: false` oder ohne `enabled` ist nicht aktiv.

**`docker-cli` im Sandbox-Profil.** Ist das Profil `sandbox` (aus `devcontainer.profile: sandbox` oder aus `--sandbox`) und `docker-cli` aktiv, endet `u-boot generate devcontainer` mit einem fachlichen Fehler: Exit-Code `10`, Code [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) im Umschlag von `--json`. Der Aufruf schreibt nichts: keine Datei und kein Verzeichnis, auch unter `--dry-run` nicht, und `--sandbox` wird nicht in `u-boot.yaml` übernommen. Die Meldung nennt das Feature und die Abhilfe in dieser Reihenfolge: das Abschalten (`u-boot config set devcontainer.features.docker-cli.enabled false`); wer Docker im Container braucht, setzt zusätzlich `u-boot config set devcontainer.sandbox.nestedRuntime docker` für einen Docker-Daemon im Container ([`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)); danach `generate devcontainer` erneut und einen vorhandenen Container neu anlegen (`u-boot devcontainer up --recreate`). `nestedRuntime: docker` allein hebt den Fehler nicht auf; mit abgeschaltetem `docker-cli` erzeugt `generate devcontainer` das Projekt mit `nestedRuntime: docker` (`--privileged` in `runArgs`, kein Socket-Mount). `config set` bleibt möglich: die Prüfung gehört nicht zur Schemaprüfung, die `config` beim Laden von `u-boot.yaml` ausführt. Im Standardprofil bleibt `docker-cli` wählbar: Exit-Code `0`, der Schlüssel `ghcr.io/devcontainers/features/docker-outside-of-docker:1` steht unter `features`, und das Feature löst keine Warnung aus.

**`doctor`.** Die Prüfung `devcontainer.features.allowlist` ([`SPEC-017`](#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry)) meldet denselben Zustand als Fehler: Ist `devcontainer.profile` in `u-boot.yaml` `sandbox` und `docker-cli` aktiv, ist der Eintrag ein Fehler (`doctor` endet bei einem Fehler-Eintrag mit Exit-Code `11`), und die Meldung nennt dieselbe Abhilfe. Verstößt zugleich eine Quelle gegen die Freigabeliste, steht beides in derselben Meldung. Im Standardprofil, bei abgeschaltetem `docker-cli` und für `go` meldet die Prüfung dazu nichts. Es gibt keine eigene Prüfungskennung.

**Lockerungen von `go`.** Im Sandbox-Profil gibt `generate devcontainer` für ein aktives `go` zwei Warnungen aus, je eine zu `--cap-add=SYS_PTRACE` und zu `--security-opt=seccomp=unconfined` (Level `warn`, Code [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil), Subjekt `.devcontainer/devcontainer.json`, Kennzeichen `SecurityRelaxation`); der Text beginnt mit „security relaxation for the devcontainer feature go:“. `init: true` des Features ist keine Sicherheitslockerung und wird nicht ausgewiesen. Der Generator schreibt die Optionen nicht in `devcontainer.json`: `runArgs` und `mounts` sind mit `go` dieselben wie ohne. Eine `version`-Überschreibung ändert die zwei Warnungen nicht; ohne Sandbox-Profil, bei abgeschaltetem oder ungesetztem `go` und bei `go` mit `source:` entstehen sie nicht.

**Freigabe der Quelle (beide Profile).** Hat ein Eintrag `enabled: true` und eine `source:`-Überschreibung, die nicht Zeichen für Zeichen in der wirksamen Freigabe steht, endet `u-boot generate devcontainer` mit einem fachlichen Fehler: Exit-Code `10`, Code [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) im Umschlag von `--json` ([`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)). Die wirksame Freigabe ist `devcontainer.featureSources.allow` zusammen mit den Quellen von `--allow-external-feature-sources` desselben Aufrufs; der Anhang der Flag-Quellen an die Liste in `u-boot.yaml` geschieht wie bisher erst nach einem erfolgreichen Lauf und unter `--dry-run` nicht. Der Vergleich ist derselbe wie bei `u-boot config set devcontainer.features.<Name>.source` und bei `doctor`: Zeichen für Zeichen, ein abschließender Schrägstrich und die Groß-/Kleinschreibung des Hosts zählen. Der Aufruf schreibt nichts: keine Datei und kein Verzeichnis, auch unter `--dry-run` nicht, und `u-boot.yaml` bleibt unverändert. Die Prüfung steht vor der Abweisung von `docker-cli` im Sandbox-Profil: Verstößt ein Aufruf gegen beides, gilt der Code [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features). Die Meldung nennt das Feature und die Quelle, dann die Abhilfe in dieser Reihenfolge: das Freigeben (`u-boot config set devcontainer.featureSources.allow <Quelle>` oder `--allow-external-feature-sources <Quelle>`), dann das Abschalten (`u-boot config set devcontainer.features.<Name>.enabled false`), und den Satz, dass der Vergleich Zeichen für Zeichen geht. Ein Eintrag mit `enabled: false` oder ohne `enabled` wird nicht erzeugt und nicht geprüft; ein Katalogschlüssel ohne `source:` braucht keine Freigabe.

**Freigabe der Taps (zweite, getrennte Freigabe).** Die Freigabe der Quelle gilt Features mit `source:`. Eine Formel aus einem Tap ([`LH-FA-DEV-013`](lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer)) braucht die eigene Freigabe `devcontainer.brewTaps.allow`; ein Eintrag der einen Liste gibt in der anderen nichts frei ([`SPEC-057`](#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen)).

**Features aus externen Quellen.** Im Sandbox-Profil gibt `generate devcontainer` für jeden Eintrag mit `enabled: true` und `source:`-Überschreibung genau eine Warnung aus, in der Reihenfolge der Namen (Level `warn`, Code [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil), Subjekt `.devcontainer/devcontainer.json`, Kennzeichen `SecurityRelaxation`). Der Text nennt den Namen und die Quelle („devcontainer feature <Name> comes from the external source <Quelle>“) und sagt, dass `u-boot` die Metadaten des Features nicht liest und das Profil für dieses Feature keine seiner Eigenschaften zusagt. Der Eintrag wird erzeugt, wenn seine Quelle freigegeben ist (Absatz davor). Ohne Sandbox-Profil, bei `enabled: false`, ohne `enabled` und für einen Katalogschlüssel ohne `source:` entsteht diese Warnung nicht.

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

- Der Container ist der, der das Label `devcontainer.local_folder` mit dem absoluten Projektordner (dem Basisverzeichnis des Befehls) trägt und kein Label `u-boot.instance` (der Container ohne Instanz; die Container von Instanzen kennt nur `u-boot devcontainer` mit `--instance`). Bei mehreren Treffern gilt der zuletzt erzeugte, auch ein gestoppter. Hat der Projektordner nur Container von Instanzen, ist der Zustand nicht ermittelbar. `u-boot` setzt ein eigenes Label (`u-boot.instance`) nur an Container von Instanzen und fügt keine `runArgs` hinzu. Bei einem Container, den `devcontainer up` oder die Startphase von `devcontainer shell` gestartet hat, setzt die Dev-Container-CLI das Label mit dem Projektordner des Hosts als Wert ([`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)).
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

**Kein Start, nur Loopback.** Der Generator ergänzt für `ollama` keinen Startbefehl und kein `forwardPorts`, im Sandbox-Profil ohne Egress-Restriktion, ohne `nestedRuntime` und ohne `devcontainer.resources.*` auch keine Option in `runArgs`, und setzt `OLLAMA_HOST` weder im Dockerfile noch in `devcontainer.json`. Der Nutzer startet den Server im Container selbst (`ollama serve`). Bei ungesetztem `OLLAMA_HOST` lauscht der Server der gepinnten Version nur auf `127.0.0.1:11434`; ein Aufruf über die Adresse des Containers scheitert. Setzt der Nutzer `OLLAMA_HOST`, ist das nicht zugesagt. Der Server läuft auf der CPU (`library=cpu` im Start-Protokoll); ein Pull über die beiden Hosts von [`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion) hinter der Egress-Restriktion gelingt, ein nicht gelisteter Host bleibt gesperrt.

**Cloud-Modelle.** Für Cloud-Modelle von Ollama (Modellname mit der Endung `-cloud`) erzeugt der Generator nichts, was über dieses Feature hinausgeht: Er setzt dafür keine Umgebungsvariable (weder `OLLAMA_HOST` noch `OLLAMA_API_KEY` noch `OLLAMA_NO_CLOUD` noch `ANTHROPIC_BASE_URL` oder `ANTHROPIC_AUTH_TOKEN`), ergänzt keinen Startbefehl, und die Zeile `ollama` der Egress-Liste führt weiter genau zwei Hosts, ohne `ollama.com` ([`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion)). Anmeldung, Start des Servers und Aufruf liegen beim Nutzer; den Weg im Container und die Grenzen der Zusage nennt das Lastenheft ([`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features)). Einen Anschluss an einen Ollama-Dienst des Hosts erzeugt das Feature `ollama-host`, unabhängig von `ollama` ([`SPEC-052`](#spec-052--feature-ollama-host-im-devcontainer-eintragstyp-container-option-meldung-und-grenzen)).

**Meldungen der Befehlsausgabe** (Erzeuger: `generate devcontainer`; Level `warn`; bei aktivem Feature genau eine je Lauf, und keine Drittanbieter-Meldung nach [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)):

| Sandbox-Profil | Code | Inhalt |
|---|---|---|
| ja | [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | Quelle ist das Release-Archiv von `github.com/ollama/ollama`; die Version nennt die Meldung nicht, sie verweist auf `OLLAMA_VERSION` im Dockerfile. Geprüft wird allein die SHA-256-Prüfsumme; der Upstream veröffentlicht keine Signatur. Subjekt `.devcontainer/devcontainer.json`. |
| nein | [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) | Das Feature hat ohne Sandbox-Profil keine Wirkung (`devcontainer.profile=sandbox` fehlt); es wurde keine Ollama-Installation erzeugt. |

Beide Meldungen tragen weder `SecurityRelaxation` noch `DegradationFallback`: Die Ausweisung der Lieferkette ist die Meldung der Befehlsausgabe, sie ist keine Sicherheitslockerung des Sandbox-Profils ([`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)). Unter `--quiet` entfallen beide.

**Nicht zugesagt.** GPU-Beschleunigung, Auswahl oder Verwaltung von Modellen durch `u-boot`, andere Architekturen und Container-Engines als die in [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) genannten und das Überstehen des Neubaus des Images.

### SPEC-051 — Features `gh` und `glab` im Devcontainer: Installation, Meldungen, Hosts und Zugangsdaten

Verfeinert [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog-Features `gh` und `glab`), [`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Token-Variablen, `doctor`), [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Lieferkette) und [`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Wirkung von `--quiet`) und berührt [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts: [`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion)). Der Katalog steht in [`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen), die Wirkung von `--quiet` in [`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), die Prüfung der Zugangsdaten in [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten). Die Grenzen der Zusagen stehen im Lastenheft ([`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features)). `gh` und `glab` sind Forge-Werkzeuge, keine Agenten.

**Aktiv** heißt: Der Eintrag `devcontainer.features.gh` beziehungsweise `devcontainer.features.glab` hat `enabled: true` und keine `source`. Mit `enabled: false`, ohne Eintrag oder mit `source` (externes Feature, [`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen)) erzeugt der Generator für das Feature keinen Dockerfile-Block, keinen Host und keine Token-Referenz; ohne Eintrag oder mit `source` entsteht auch die Meldung zur Lieferkette nicht.

**Wirkung in beiden Profilen.** Die Features wirken im Standardprofil und im Sandbox-Profil (`devcontainer.profile: sandbox` oder `u-boot generate devcontainer --sandbox`). Im Standardprofil ändert ein aktives Feature das Dockerfile und sonst nichts: `devcontainer.json` ist byte-gleich zum Lauf ohne die Features, ohne `remoteEnv`, ohne `mounts` und ohne Schlüssel unter `features`. In keinem Profil ergänzen `gh` und `glab` ein Volume unter `mounts`. Mit `gh` und `glab` als einzigen Features, die Dockerfile-Blöcke tragen (ohne `ollama` und ohne Laufzeit-Feature), steht `USER root` in beiden Profilen genau einmal; im Standardprofil steht es nur mit einem Grund, einem aktiven Feature oder der Überschreibung `devcontainer.user.uid`, und bei beiden einmal.

**Installation.** Das erzeugte Dockerfile trägt hinter `USER root` je aktivem Feature einen Block, der das Programm aus dem Release-Archiv des Upstream-Projekts installiert: `gh` aus dem Archiv von `cli/cli` auf `github.com`, `glab` aus dem Archiv von `gitlab-org/cli` auf `gitlab.com`. Version und SHA-256-Prüfsumme je Architektur (`amd64`, `arm64`) stehen als `GH_VERSION`, `GH_SHA256_AMD64`, `GH_SHA256_ARM64` beziehungsweise `GLAB_VERSION`, `GLAB_SHA256_AMD64`, `GLAB_SHA256_ARM64` an einer Stelle im Dockerfile; die Spezifikation führt sie nicht ein zweites Mal. Der Block lädt das Archiv, prüft die Prüfsumme mit `sha256sum -c` vor dem Entpacken (eine Abweichung bricht den Image-Build ab, ebenso eine andere Architektur als `amd64` und `arm64`) und entpackt nur das Programm nach `/usr/local`, ohne den Eigentümer des Archivs zu übernehmen (`--no-same-owner`). Der Block setzt weder `ENV` noch `EXPOSE` noch `ENTRYPOINT` noch `CMD`, führt kein `auth login` aus, nennt weder `GH_TOKEN` noch `GITLAB_TOKEN` und enthält keinen Schritt, der eine Signatur oder eine Attestierung prüft. Der Block steht vor dem letzten `USER vscode` des Dockerfiles, und das gebaute Image läuft als `vscode`; es enthält weder `/home/vscode/.config/gh` noch `/home/vscode/.config/glab-cli`. Ein zweiter Lauf von `generate devcontainer` ändert nichts. Wird ein Feature nach dem ersten Lauf aktiviert, ersetzt der Generator den verwalteten Block des Dockerfiles und lässt Zeilen des Nutzers dahinter stehen; nach dem Abschalten entfällt der Block des Features, Zeilen des Nutzers bleiben.

**Reihenfolge der Blöcke.** Mit `nestedRuntime: podman` oder `docker`, dem Feature `ollama` und der Egress-Restriktion stehen die Blöcke in dieser Reihenfolge im Dockerfile: Block der Container-Runtime, `ollama`, `gh`, `glab`, Egress-Restriktion. Die Egress-Liste ist die sortierte Vereinigung der Hosts ohne Duplikat, und ein zweiter Lauf ändert nichts.

**`version` und `source`.** Eine `version` ohne `source` ist bei `gh` und `glab` in beiden Profilen ein Validierungsfehler von `generate devcontainer` (Exit-Code `10`, es entsteht kein Dockerfile); die Meldung nennt `features.<Schlüssel>.version`. Dieselbe Regel gilt beim Schreiben: `config set devcontainer.features.<Schlüssel>.version <Wert>` ohne `source` endet mit Exit-Code `10` und lässt die Datei unverändert ([`SPEC-040`](#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)). Die Prüfsummen hängen an der gepinnten Version. Mit `source` ist der Schlüssel ein externes Feature ([`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen)): Er steht unter `features`, das Dockerfile trägt keinen Block und im Standardprofil kein `USER root`, `remoteEnv` keine Token-Referenz, und im Sandbox-Profil entsteht die eine Warnung zur externen Quelle ([`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen)).

**Meldungen der Befehlsausgabe** (Erzeuger: `generate devcontainer`; Level `warn`, Code [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), Subjekt `.devcontainer/devcontainer.json`). Bei einem aktiven Feature entsteht in jedem Profil genau eine Meldung je Feature und Lauf; das andere Feature fügt keine hinzu. Der Text ist Englisch und sagt:

- die Quelle des Archivs (`cli/cli` beziehungsweise `gitlab-org/cli`) und dass die Version im Dockerfile als `GH_VERSION` beziehungsweise `GLAB_VERSION` steht; die Meldung nennt die Version selbst nicht;
- dass die SHA-256-Prüfsumme des Archivs geprüft wird;
- dass für das Archiv keine Signatur gefunden wurde, nach Suchmuster und nicht abschließend;
- bei `gh`, dass die GitHub-Attestierung des Projekts von `u-boot` nicht geprüft wird; die Meldung zu `glab` nennt keine Attestierung.

Die Meldungen tragen weder `SecurityRelaxation` noch `DegradationFallback`: Sie sind die Ausweisung der Lieferkette und keine Sicherheitslockerung des Sandbox-Profils ([`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)). Unter `--quiet` entfallen sie, die Blöcke im Dockerfile bleiben.

**Hosts.** Die Egress-Liste führt für `gh` die Hosts `api.github.com` und `github.com` (beide stehen auch in der Basisliste) und für `glab` den Host `gitlab.com`, sonst nichts; die Zeilen stehen in [`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion), die Anzeige führt `github.com` und `api.github.com` mit den Herkünften `base` und `feature:gh`, `gitlab.com` mit `feature:glab`. Hinter der Egress-Restriktion führt das Startscript die drei Hosts; ein Aufruf von `glab` gegen ein öffentliches Projekt auf `gitlab.com` gelingt ohne Token, wenn `gitlab.com` in der Liste steht, und scheitert ohne diesen Host mit einer Meldung zur Namensauflösung (`lookup gitlab.com`).

**Zugangsdaten.** Im Sandbox-Profil trägt `remoteEnv` neben `GIT_TOKEN` je aktivem Feature die Referenz `"GH_TOKEN": "${localEnv:GH_TOKEN}"` beziehungsweise `"GITLAB_TOKEN": "${localEnv:GITLAB_TOKEN}"`, nur als Referenz, in der festen Reihenfolge `ANTHROPIC_API_KEY`, `CLAUDE_CODE_OAUTH_TOKEN`, `GH_TOKEN`, `GITLAB_TOKEN` ([`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)); im Standardprofil entsteht kein `remoteEnv` und damit keine Token-Referenz. Der Wert eines Tokens steht in keiner erzeugten Datei, keinem Image, keinem Volume und nicht in der Umgebung des Containers (`Config.Env`). Im Sandbox-Profil gibt `u-boot devcontainer shell` jede Variable, auf die `remoteEnv` verweist (bis zu vier: `ANTHROPIC_API_KEY`, `CLAUDE_CODE_OAUTH_TOKEN`, `GH_TOKEN`, `GITLAB_TOKEN`), an die Sitzung weiter, soweit sie auf dem Host gesetzt sind oder der Befehl sie mit der Auflösung der Token-Quelle ermittelt hat ([`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)), und nur als Name ([`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)). `u-boot` führt weder `gh auth login` noch `glab auth login` aus und speichert kein Token. `u-boot doctor` prüft im Sandbox-Profil die Referenz und den Klartext nach [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten); im Standardprofil ist der Eintrag übersprungen.

**Nicht zugesagt.** Aufrufe mit gültigem Token und deren Hosts; Hosts der Anmeldung sowie von Telemetrie und Aktualisierung; selbst gehostete Instanzen (GitHub Enterprise, GitLab Self-Managed); andere Versionen als die gepinnten; `glab duo`, `glab orbit` und das Nachladen weiterer Programme durch `glab`; ein Bau auf `arm64`; die Kombination mit `ollama` und `nestedRuntime` über die Reihenfolge der Blöcke, die Volumes und die Hosts hinaus.

### SPEC-052 — Feature `ollama-host` im Devcontainer: Eintragstyp, Container-Option, Meldung und Grenzen

Verfeinert [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog-Feature `ollama-host`, Anschluss an einen Ollama-Dienst des Hosts), [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) (Ausweisung als Sicherheitslockerung) und [`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Wirkung von `--quiet`, Ausnahme 1) und berührt [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts: [`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion)) und [`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten). Der Katalog steht in [`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen), die Wirkung von `--quiet` in [`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), die Zugangsdaten der Agenten in [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten), das Feature `ollama` in [`SPEC-047`](#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen). Die Grenzen der Zusagen stehen im Lastenheft ([`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features)).

**Aktiv** heißt: Der Eintrag `devcontainer.features.ollama-host` hat `enabled: true` und keine `source`. Mit `enabled: false`, ohne `enabled` oder ohne Eintrag erzeugt der Generator weder die Zeile noch die Meldung, in beiden Profilen, und die erzeugten Dateien des Projekts außer `u-boot.yaml` sind gleich dem Lauf ohne den Eintrag; im Standardprofil trägt `runArgs` dann nur die Zeile der Umgebungsdatei ([`SPEC-055`](#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung)). Mit `source` ist der Schlüssel ein externes Feature ([`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen)): Er steht unter `features`, `runArgs` trägt keine Zeile des Anschlusses, und im Sandbox-Profil entsteht die eine Warnung zur externen Quelle ([`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen)) und keine Meldung des Anschlusses. Das Feature ist unabhängig von `ollama`: Es braucht weder dieses Feature noch ein anderes.

**Option.** Der Eintrag trägt genau eine Container-Option, die Zeile `--add-host=host.docker.internal:host-gateway`; sie macht den Host aus dem Container unter dem Namen `host.docker.internal` adressierbar. Zu diesem Anschluss erzeugt der Generator sonst nichts: keinen Relay, kein `forwardPorts`, keinen Eintrag der Egress-Liste und keine Umgebungsvariable. Auch zusammen mit den Features `ollama`, `claude-code` und `node` steht keiner der Namen `OLLAMA_HOST`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_AUTH_TOKEN`, `OLLAMA_API_KEY` und `OLLAMA_NO_CLOUD` und nicht der Port `11434` in einer erzeugten Datei, in beiden Profilen, mit und ohne Egress-Restriktion und mit `nestedRuntime: podman` und `docker`.

**Wirkung in beiden Profilen.** Mit aktivem Eintrag trägt `devcontainer.json` die Zeile in `runArgs`.

- Standardprofil: `runArgs` hat ohne `devcontainer.resources.*` genau zwei Einträge, die Zeile der Umgebungsdatei ([`SPEC-055`](#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung)) und diese Zeile (mit den Schlüsseln stehen `--memory` und `--cpus` zwischen beiden, [`SPEC-053`](#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs)) und steht nach `build` und vor `forwardPorts`, `features` und `remoteUser`; die Datei ist sonst gleich dem Lauf ohne den Eintrag, ebenso das Dockerfile.
- Sandbox-Profil (`devcontainer.profile: sandbox` oder `u-boot generate devcontainer --sandbox`): Die Zeile ist der letzte Eintrag von `runArgs`, hinter den Optionen von `nestedRuntime`, der Egress-Restriktion und von `devcontainer.resources.*`, und steht genau einmal. Die übrigen Schlüssel von `devcontainer.json`, darunter `mounts`, `postStartCommand` und `remoteEnv`, das Dockerfile und mit Egress-Restriktion `egress-init.sh` sind gleich dem Lauf ohne den Eintrag.

Mit den Features `ollama`, `gh` und `glab` und mit `nestedRuntime: podman` oder `docker` und Egress-Restriktion sind `runArgs` (in dieser Reihenfolge: Zeile der Umgebungsdatei, Optionen der Container-Runtime, Egress-Restriktion, Optionen von `devcontainer.resources.*` falls gesetzt, Zeile des Anschlusses), `mounts` und Dockerfile gleich dem Lauf ohne den Eintrag bis auf die Zeile; die Hosts der Egress-Liste sind die sortierte Vereinigung wie ohne ihn. Ein zweiter Lauf ändert nichts. Wird der Eintrag nach dem ersten Lauf aktiviert, ersetzt der Generator den verwalteten Teil der `devcontainer.json` und lässt Zeilen des Nutzers dahinter stehen; nach dem Abschalten entfällt die Zeile, Zeilen des Nutzers bleiben.

**`version` und `source`.** Der Eintrag hat keine Version. Eine `version` ohne `source` ist in beiden Profilen ein Validierungsfehler von `generate devcontainer` (Exit-Code `10`, es wird nichts geschrieben); die Meldung nennt `features.ollama-host.version` und sagt, dass der Eintrag kein OCI-Feature ist. Dieselbe Regel gilt beim Schreiben: `config set devcontainer.features.ollama-host.version <Wert>` ohne `source` endet mit Exit-Code `10` und lässt die Datei unverändert ([`SPEC-040`](#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)). Die Pfade `devcontainer.features.ollama-host.enabled`, `.source` und `.version` sind Platzhalter-Pfade der Konfiguration; die Menge der festen Pfade ändert das nicht, und `devcontainer.features.ollama-host.relay` ist ungültig.

**Meldung der Befehlsausgabe** (Erzeuger: `generate devcontainer`; Level `warn`, Subjekt `.devcontainer/devcontainer.json`). Bei aktivem Eintrag entsteht in jedem Profil genau eine Meldung je Lauf, auch im zweiten Lauf ohne Änderung; sie ist getrennt von den Meldungen des Features `ollama`.

| Profil | Code | Sicherheitslockerung (`--quiet`) | Inhalt |
|---|---|---|---|
| Sandbox-Profil | [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil) | ja, bleibt sichtbar | Der Text unten, dazu die Klausel, dass das Profil den Schaden begrenzt und keine Isolationsgrenze ist („limits damage and is not an isolation boundary“). |
| Standardprofil | [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) | nein, entfällt | Der Text unten, ohne die Klausel zum Profil. |

Der Text ist Englisch. Er nennt die Option, den Namen `host.docker.internal`, dass die Erreichbarkeit von der Adresse abhängt, auf der der Dienst lauscht („depends on the address it listens on“), dass ein Dienst, der nur auf der Loopback-Adresse des Hosts lauscht, nur über einen Relay oder eine andere Bindung erreichbar ist, die der Nutzer einrichtet („only through a relay or another binding that you set up“), und dass `u-boot` keinen Relay startet. Bei Sandbox-Profil **und** aktiver Egress-Restriktion steht im selben Eintrag zusätzlich, dass der Anschluss dahinter nicht funktioniert („does not work behind it“) und dass `u-boot` keine Adresse in die Allowlist einträgt; das ist eine Warnung und kein Fehler (Exit-Code `0`, dieselbe Aktion wie ohne den Eintrag, die übrigen Warnungen unverändert). Im Standardprofil mit gesetztem `egress.enabled` und im Sandbox-Profil ohne Egress-Restriktion steht dieser Zusatz nicht. Der Text nennt weder den Cloud-Dienst noch Modelle noch Kontingent, keine Kennung und keine Zahl.

Unter `--quiet` bleibt die Meldung im Sandbox-Profil sichtbar, mit dem Zusatz zur Egress-Restriktion, wenn er gilt (Kennzeichen `SecurityRelaxation`, [`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) Ausnahme 1); im Standardprofil entfällt sie. Sie trägt in keinem Profil `DegradationFallback`.

**Hosts und Zugangsdaten.** Der Eintrag ergänzt keinen Host: Die Egress-Liste mit dem Eintrag ist gleich der ohne ihn, `host.docker.internal` und `ollama.com` stehen weder in der Liste noch im Startscript, die Liste führt keine IP-Adresse, und `u-boot egress list` führt keinen Host mit der Herkunft `feature:ollama-host` ([`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion)). Er führt keine Variable und keine Token-Referenz: `remoteEnv` des Sandbox-Profils trägt dieselben Einträge wie ohne ihn ([`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)), und `u-boot doctor` meldet für ihn nichts Neues; die Prüfungen `devcontainer.features.allowlist` und `devcontainer.features.drift` melden für ihn keinen Befund, auch wenn `devcontainer.json` keinen Schlüssel unter `features` enthält.

**Auflösung des Namens im Container.** Im Container, den `u-boot devcontainer up` aus dem erzeugten Projekt startet, steht die Zeile unter den zusätzlichen Hosts der Container-Konfiguration, und `getent hosts host.docker.internal` im Container ergibt eine Adresse; geprüft für das Standardprofil und für das Sandbox-Profil mit Egress-Restriktion (Integrationstest, Linux-Host mit lokaler Docker-Engine). Ohne das Feature löst der Name im Standardprofil nicht auf.

**Nicht zugesagt.** Die Erreichbarkeit eines Dienstes des Hosts aus dem Container außer in dem Aufbau, den das Lastenheft benennt (Linux-Host mit lokaler Docker-Engine, Dienst auf der Loopback-Adresse mit einem Relay des Nutzers, Features `claude-code` und `node`; der Integrationstest prüft nur die Auflösung des Namens, nicht den Dienst); die Wirkung der erzeugten Zeile gegenüber einer gleichen Zeile, die der Nutzer von Hand gesetzt hat; der Anschluss hinter der Egress-Restriktion; das Sandbox-Profil, wie `u-boot` es ohne Egress-Restriktion erzeugt, als Aufbau des Anschlusses; andere Hosts, Engines und Bindungen als die im Lastenheft genannten ([`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features)).

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

**Aufruf.** `u-boot devcontainer up [--recreate] [--instance <name>]`, `u-boot devcontainer down [--instance <name>]` und `u-boot devcontainer shell [--instance <name>]`. Der Projektordner ist das Arbeitsverzeichnis des Aufrufs. Ohne Unterbefehl, mit einem unbekannten Unterbefehl, mit zusätzlichen Argumenten sowie mit `--dry-run` oder `--diff` endet der Befehl mit Exit-Code `2`, denn die Unterbefehle schreiben keine Projektdatei. Mit `--json` entsteht bei `--dry-run`, `--diff` und zusätzlichen Argumenten ein Umschlag mit Exit-Code `2` (bei den ersten beiden mit Diagnose [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes)); der Aufruf ohne Unterbefehl endet auch mit `--json` ohne Umschlag. `shell` weist außerdem `--json` mit Exit-Code `2` ab, weil der Befehl eine interaktive Sitzung öffnet und die Standardausgabe ihr gehört; die Abweisung ist ein Umschlag mit `subcommand: "shell"` und Diagnose [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes); keine dieser Abweisungen startet etwas.

**Voraussetzungen von `up` und `shell` (Exit-Code `10`).** `up` und `shell` prüfen vor der Ausweisung des Sockets und vor dem Aufruf des Launchers: (1) der Projektordner ist ein absoluter Pfad ohne Komma, Anführungszeichen, Zeilenwechsel, Wagenrücklauf und NUL, weil ein `--mount`-Wert diese Zeichen nicht trägt; die Meldung nennt den Pfad nicht; (2) `.devcontainer/devcontainer.json` existiert; (3) jeder Name, auf den diese Datei mit `${localEnv:<NAME>}` verweist, ist ein einfacher Variablenname (`[A-Za-z_][A-Za-z0-9_]*`). Ein Verstoß endet mit Exit-Code `10`, ein Lesefehler der Datei mit Exit-Code `14`; in beiden Fällen startet kein Hilfs-Container, und die Ausweisung des Sockets entfällt. `down` liest keine Projektdatei und braucht keine `devcontainer.json`.

**Instanzen** (`--instance <name>`). Ein Projektordner kann neben dem Container ohne Instanz, den `up` ohne Flag anlegt, mehrere Container von Instanzen haben; `u-boot` führt keine Liste der Instanzen und räumt keine auf.

- **Name.** Ein Name besteht aus Kleinbuchstaben, Ziffern und Bindestrichen, beginnt und endet mit einem Buchstaben oder einer Ziffer und ist 1 bis 32 Zeichen lang. Ungültig sind die leere Zeichenkette (auch `--instance=`), weil die Dev-Container-CLI sie ablehnt, und ein Name mit `=`, bei dem `devcontainer up` nicht zurückkehrt. Ein ungültiger Name endet bei `up`, `down` und `shell` mit Exit-Code `10`, bevor `u-boot` die Engine befragt, die Ausweisung des Sockets ausgibt oder einen Hilfs-Container startet; die Regel steht einmal im Port `driven`, und der Launcher und der Inspector wenden sie vor jedem Aufruf von `docker` noch einmal an.
- **Identität.** Mit Instanz hängt `u-boot` dem Aufruf der Dev-Container-CLI zwei `--id-label` an (Hilfs-Container unten): `devcontainer.local_folder=<Ordner>` und `u-boot.instance=<Name>`. Die CLI bildet die `devcontainerId` aus ihren Id-Labels, deshalb hat jede Instanz einen eigenen Container und, im Sandbox-Profil, eigene Volumes (ihre Namen tragen `${devcontainerId}`); die Messwerte stehen im Messprotokoll der Instanzen, die Zusage der Menge in der Anforderung. Der Container ohne Instanz bekommt kein Label von `u-boot` und keine `--id-label`: Ein Label änderte die `devcontainerId` und damit die Volume-Namen bestehender Projekte.
- **Suche.** Jede Suche eines Containers (`down`, die Abfrage von `up` und `shell`, ob der Container schon läuft, die Lesung der Statusdateien nach dem Start) wählt mit Instanz den jüngsten Container, der das Label des Projektordners und `u-boot.instance=<Name>` trägt (`docker ps -a` mit zwei `--filter label=…`). Ohne Instanz wählt sie den jüngsten Container mit dem Label des Projektordners, der kein Label `u-boot.instance` trägt; `u-boot` filtert dafür nicht auf das fehlende Label, sondern liest das Label jeder Zeile der Ausgabe von `docker ps` (`{{.Label "u-boot.instance"}}`). `doctor`, `generate devcontainer` und `init` suchen ohne Instanz; für sie ist ein Projekt mit nur Containern von Instanzen „nicht ermittelbar“ ([`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)), kein Befund.
- **Ohne `--instance`, nur Instanzen.** Hat der Projektordner Container, aber keinen ohne Instanz, endet `up` (auch mit `--recreate`), `down` und `shell` ohne `--instance` mit Exit-Code `10` und einem Hinweis auf `--instance`; es läuft weder ein Hilfs-Container noch `docker stop` noch eine Sitzung. Bei `up` und `shell` steht die Prüfung nach der Ausweisung des Sockets (die Ausweisung entsteht, bevor die Engine befragt wird) und vor der Vorab-Prüfung des Launchers. Hat der Projektordner gar keinen Container, startet `up` und `shell` ohne `--instance` den Container ohne Instanz, und `down` endet mit Exit-Code `0`. Ein Container, der zwischen der Prüfung und dem Aufruf der CLI entsteht oder verschwindet, ist nicht abgesichert.
- **Mit `--instance`.** `up` startet den Container der Instanz oder findet ihn; ein zweites `up --instance <name>` ohne `--recreate` trifft denselben Container (Zeile und Feld `alreadyRunning` wie sonst). `up --instance <name> --recreate` entfernt nur den Container dieser Instanz und legt ihn neu an, benannte Volumes bleiben; die übrigen Container des Projekts laufen weiter. `down --instance <name>` stoppt nur den Container der Instanz; gibt es keinen, endet es mit Exit-Code `0`, ohne `docker stop`.
- **Container neben dem Container ohne Instanz.** Gibt es den Container ohne Instanz und eine neuere Instanz, treffen `up`, `down` und `shell` ohne `--instance` den Container ohne Instanz, und mit `--instance` den der Instanz.
- **Containername** ([`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)). Der Name des gestarteten Containers lautet `u-boot-<Projektname>`, mit `--instance <name>` `u-boot-<Projektname>-<name>`; er setzt sich aus dem Präfix `u-boot`, dem Projektnamen und der Instanz zusammen. Der Projektname ist der Wert `project.name` der `u-boot.yaml` des Projektordners, oder, wenn die Datei fehlt, nicht lesbar oder keinen Namen trägt, der normalisierte Basisname des Projektordners (Regeln [`LH-FA-INIT-002`](lastenheft.md#lh-fa-init-002--projektname)); er ist keine Prüfung des Projektnamens selbst, sondern die Quelle des Namenssteils. Die Zeichengrenzen der Engine gelten für den ganzen Namen: Erstes Zeichen Buchstabe oder Ziffer, danach Buchstaben, Ziffern, `.`, `_` und `-` (gemessen an Docker Engine 29.9.2, Linux `amd64`; die Messung fand bis zu 257 Zeichen keine Grenze; die Zusammensetzung bleibt unter den langen gemessenen Namen). Eine Zusammenstellung, die diese Regel verletzt, endet mit Exit-Code `10` und einer Meldung zum Namen, bevor die Ausweisung des Sockets entsteht und die Engine befragt wird; dasselbe gilt für einen Namen der Instanz. Die Regel steht einmal im Port `driven`; der Launcher wendet sie vor dem ersten `docker`-Aufruf noch einmal an.
- **So kommt der Name in die Konfiguration.** `u-boot` rendert für den Aufruf eine Kopie der gelesenen Konfiguration und setzt dort in `runArgs` den Wert `--name=<Name>` — an eine bestehende `runArgs`-Liste angehängt, an eine fehlende ergänzt. Das geschieht im Text der für diesen Aufruf gerenderten Kopie: Die erzeugte `devcontainer.json` des Projekts ändert sich nicht (kein Golden-Case-Unterschied), und Kommentare oder Formatierung der Quelldatei sind nicht Teil der Kopie. `u-boot` übergibt die Kopie dem Hilfs-Aufruf per `--override-config`, als Datei außerhalb des Projektordners, die während des Aufrufs gemountet und nach dem Aufruf entfernt wird; auch ohne das `runArgs`-Overlay wäre `--name` kein Argument des `docker run` des Hilfs-Containers. Trägt die Konfiguration keine mergebare Struktur (etwa weil sie kein lesbares JSON mit Kommentaren ist), wird kein Name übergeben: Der Start behält sein Verhalten, der Container trägt dann keinen Namen, und die Zeile der Ausgabe nennt keinen ([`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Zusage nur für den docker-getesteten Weg). Die CLI löst die relativen Pfade der Konfiguration unabhängig vom Ablageort der Kopie auf; der Hilfs-Aufruf ändert die Id-Labels nicht — mit den Labels wechselt weder die `devcontainerId` des Containers noch die Namen der Volumes (Messung M3 und M6 im Messprotokoll [`slice-v1-devcontainer-container-name`](../docs/messungen/slice-v1-devcontainer-container-name/README.md)).
- **Namenskollision vor dem Start.** Der Launcher fragt nach der Vorab-Prüfung der Engine und vor dem Bau des Hilfs-Images, ob ein Container des Namens besteht (Suche `docker ps -a` mit Namensfilter) und dieser zu einem anderen Projektordner gehört; ein Container dieses Projektordners ist der zu übernehmende, ein Container ohne `devcontainer.local_folder`-Label gehört zu keinem Projektordner von `u-boot`. Kollision: Der Aufruf endet mit Exit-Code `10` und einer klaren Meldung, die Namen und Projektordner des haltenden Containers nennt, bevor ein Hilfs-Container entsteht; der haltende Container wird nicht entfernt oder umbenannt, auch nicht mit `--recreate` — die CLI entfernt nach Label, nicht nach Name (gemessen). Ein Container, der zwischen der Prüfung und dem Start den Namen belegt, ist nicht abgesichert; der Start endet dann mit dem Fehler der CLI (Exit-Code `12`). Ein gescheiterte Befragung ist kein Befund: der Start läuft weiter. Bei `down` gibt es diese Prüfung nicht: `down` stoppt über die Labels und ohne Namen.

**Vorab-Prüfung der Engine** (bei `up` und in der Startphase von `shell` vor dem Bau des Hilfs-Images und dem Start des Hilfs-Containers, bei `shell` zusätzlich vor der Probe der Sitzung, bei `down` vor `docker stop`): (1) das `docker`-Binary ist auf dem `PATH`, (2) der Daemon antwortet, (3) die Engine ist lokal (`DOCKER_HOST` oder der aktive Docker-Kontext zeigt auf einen lokalen Unix-Socket, wie in [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)). Scheitert Schritt 1 oder 2, lässt sich der Engine-Zugang nicht bestimmen (etwa weil der Kontext nicht lesbar ist) oder ist die Engine nicht lokal (`DOCKER_HOST` mit `tcp://…`, ein Kontext mit `ssh://…`), endet der Befehl mit Exit-Code `11`. Nach einem gescheiterten Schritt laufen weder Bau noch Start noch `docker stop` noch `docker exec`.

**Hilfs-Image.** `u-boot` baut das Image lokal aus einem eingebetteten Kontext aus genau drei Dateien (`Dockerfile`, `package.json`, `package-lock.json`). Der Tag hat die Form `u-boot-devcontainer-cli:<exakte Version der Dev-Container-CLI>-<zwölf Hexzeichen>`. Gebaut wird nur, wenn kein Image mit diesem Tag besteht: `docker build --force-rm -t <Tag> -` mit dem Kontext als Tar-Strom auf der Standardeingabe. Pin-Form: beide Basis-Images stehen mit einem Digest (`@sha256:…`) im `FROM`; `package.json` pinnt die CLI auf eine exakte Version, jeder Eintrag der Lockdatei trägt eine `integrity`, die Lockdatei nennt dieselbe Version der CLI, und die Installation ist `npm ci --ignore-scripts`. Eine verfälschte `integrity` bricht den Bau. Scheitert der Bau, endet der Befehl mit Exit-Code `12`; die Meldung nennt das Ende der Standardfehlerausgabe des Baus.

**Hilfs-Container.** Der Aufruf lautet:

```
docker run --rm --name <Name> \
  --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock \
  --mount type=bind,source=<Ordner>,target=<Ordner> \
  [-e <NAME> …] <Hilfs-Image> \
  devcontainer up --workspace-folder <Ordner> [--remove-existing-container] \
  [--id-label devcontainer.local_folder=<Ordner> --id-label u-boot.instance=<Name>]
```

Der Name beginnt mit einem festen Präfix. Die Argumentliste enthält weder `-v` noch `--volume`, `--privileged`, `--network` oder `--cap-add`, und sie bindet kein Home-Verzeichnis ein. Der Projektordner ist schreibend eingehängt; ein `initializeCommand`, der in ihn schreibt, gelingt. `--recreate` ergänzt `--remove-existing-container`. Die beiden `--id-label` stehen nur mit `--instance` im Aufruf (Absatz *Instanzen*).

- **Spiegelpfad.** Der Projektordner liegt im Hilfs-Container unter demselben absoluten Pfad wie auf dem Host, und derselbe Pfad geht als `--workspace-folder` an die CLI. Das Label `devcontainer.local_folder` des gestarteten Containers trägt damit den Pfad des Hosts, und die Suche über dieses Label ([`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)) findet den Container; ein Pfad unter `/w/<Name>` im Hilfs-Container ergäbe ein anderes Label. Das Label setzt die Dev-Container-CLI; `u-boot` setzt ein eigenes Label (`u-boot.instance`) nur an den Container einer Instanz, nie an den Container ohne Instanz. Ein Projektordner, den die Engine nicht sieht, ergibt Exit-Code `12`: `--mount type=bind` bricht ab und legt den Quellpfad nicht an.
- **Lebensdauer.** Nach `up` bleibt kein Hilfs-Container zurück. Endet `docker run` mit einem Fehler oder bricht der Aufrufer ab, entfernt `u-boot` den Hilfs-Container mit `docker rm -f <Name>`; gelingt das nicht, nennt die Fehlermeldung es. Ein Abbruch beendet den Prozess mit einem Exit-Code ungleich `0`; der Dev-Container, den die CLI bis dahin angelegt hat, bleibt bestehen.
- **Fortschritt.** Die Standardfehlerausgabe der CLI geht als Fortschritt an die Konsole. Die Ergebniszeile der CLI steht auf deren Standardausgabe, nicht auf der Standardfehlerausgabe.

**Umgebung und Zugangsdaten.** Die Namen liest `u-boot` aus dem Text von `.devcontainer/devcontainer.json`: jede Form `${localEnv:<NAME>}` und `${localEnv:<NAME>:<Vorgabe>}`, sortiert und ohne Duplikate; ein Verweis in einem Kommentar zählt mit, andere Formen (`${env:…}`, `${containerEnv:…}`) werden nicht gelesen. Der Hilfs-Container bekommt nur die Namen, die auf dem Host gesetzt sind, als `-e <NAME>` ohne Wert. Variablen, auf die die Datei nicht verweist (etwa Proxy-Variablen), reicht der Befehl nicht durch.

**Auflösung der Token-Quelle** (`devcontainer.sandbox.hostGitToken`, [`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)). Steht im Sandbox-Profil der Schalter auf `true` und verweist `devcontainer.json` auf die Variable, und ist die Variable `GIT_TOKEN` auf dem Host **nicht** gesetzt, ermittelt `u-boot` den Token auf dem Host in fester Reihenfolge: Quelle `gh auth token` (das `gh` des Hosts ist installiert und angemeldet; ein Ende mit Fehler — etwa „nicht angemeldet“ — gilt wie fehlend). Bei einer Quelle übergibt `u-boot` den Wert dem Hilfs-Container wie eine Host-Variable: der Name als `-e GIT_TOKEN`, der Wert in der Umgebung des `docker`-Client-Prozesses, nie in einer Kommandozeile, in einer Datei, in einem Image, in einem Volume oder in einer Ausgabe; der Wert der Quelle steht in keiner Datei des Projekts (die erzeugte Konfiguration trägt weiter die Referenz). Der Befehl nennt in einer Warnung die benutzte Quelle — „the git token of this container comes from the environment variable GIT_TOKEN of the host“ beziehungsweise „… from the gh session of the host (gh auth token)“ — und nie den Wert; findet er keine Quelle, weist er ohne Abbruch darauf hin. Ohne Sandbox-Profil, ohne `true` oder ohne Verweis löst der Befehl nichts auf; mit Schalter und ohne Verweis weist er darauf hin. Der Wert, den `u-boot` in einem Aufruf ermittelt, ist derselbe für den Hilfs-Container des Starts und für die Sitzung von `shell`; die Sitzung bekommt ihn wie die durchgereichten Host-Variablen als `-e GIT_TOKEN` mit dem Wert in der Umgebung des `docker`-Client-Prozesses des `docker exec`.

- **Wozu `u-boot` das „gh“ und „git“ des Hosts befragt.** Nur die Befehle `up` und `shell` rufen mit der Auflösung `gh auth token` und `git config` des Hosts auf; `doctor` liest nur Projektdateien und sieht weder die Umgebung des Hosts noch diese Quellen (Grenze der Prüfung, [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)). Befragt werden: die Variable `GIT_TOKEN` der Umgebung, die Ausgabe von `gh auth token` auf dem Host und die Werte von `git -C <Projektordner> config user.name` und `… user.email`; keine anderen Schlüssel, und keine Konfigurationsdatei des Hosts öffnet `u-boot` selbst.

- **Identität** (beide Profile, [`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)). Nach einem gelungenen Start setzt der Befehl mit `git config --global` im Container-Heim des Remote-Benutzers die Werte von `user.name` und `user.email`, die er auf dem Host mit `git -C <Projektordner> config <Schlüssel>` liest (der wirksame Wert, mit den bedingten Einfügungen; ein leerer Wert gilt wie fehlend); ein im Container bereits gesetzter, nichtleerer Wert bleibt; ein leerer Wert des Containers wird überschrieben. Jeder Wert geht als `-e <VARIABLENNAME>` in den `docker exec`, der Wert in der Umgebung des `docker`-Client-Prozesses, nie in einer Kommandozeile; die Aufrufe sind best effort: scheitern sie (kein `git` im Container), steht eine Warnung, und der Exit-Code des Starts bleibt. Liefert der Host keine Werte oder kein `git`, weist der Befehl ohne Abbruch darauf hin.

- **Wo der Wert nicht steht.** Das `docker run` und der `docker exec` des `docker`-Client-Prozesses von `u-boot` tragen nur den Namen, in keiner Zeile der Ausgabe steht ein Wert. Mit der Form `remoteEnv` steht der Wert (aus der Host-Umgebung oder aus der Auflösung) zudem in keinem Label und keinen Metadaten des gestarteten Containers (`Config.Env` trägt ihn nicht; das Label `devcontainer.metadata` trägt die Referenz, nicht den Wert), nicht im Hilfs-Image, in keiner Datei unter dem Projektordner und in keinem Volume des Sandbox-Containers.
- **Grenzen.** Die Dev-Container-CLI schreibt den aufgelösten Wert für die Dauer eines Lebenszyklus-Befehls in die Kommandozeilen der `docker`-Prozesse im Hilfs-Container (auch bei `remoteEnv`; die Messung des Slices belegt dies für die Container-Befehle der CLI). Bei `containerEnv` steht der aufgelöste Wert in den Metadaten (`Config.Env`) des Dev-Containers, solange dieser besteht.
- **Maskierung.** In dem, was `u-boot` als Fortschritt weitergibt, und in seinen Fehlermeldungen ersetzt es den Wert einer durchgereichten Variable durch `[masked]`, wenn er unverändert vorkommt und mindestens vier Byte lang ist; zerlegte oder kodierte Formen und kürzere Werte bleiben stehen. Der aufgelöste Wert ist wie ein Host-Wert maskiert. Die Felder der Ergebniszeile sind Daten und werden nicht maskiert. Die letzten Bytes des Fortschritts (bis zur Länge des längsten maskierbaren Werts minus eins) erscheinen mit dem nächsten Schreiben oder am Ende des Aufrufs.

**Start und Ergebniszeile.** Die letzte nicht leere Zeile der Standardausgabe der CLI ist die Ergebniszeile (JSON); sie bleibt auch nach umfangreicher Ausgabe davor lesbar. Der Start gilt als gelungen, wenn `outcome` den Wert `success` hat und `containerId` aus Hexzeichen besteht; `remoteUser` und `remoteWorkspaceFolder` übernimmt `u-boot` aus der Zeile. Ein anderes `outcome` (die Meldung der CLI steht im Fehlertext), eine fehlende oder keine JSON-Ergebniszeile und ein Exit-Code des `docker run` ungleich `0` enden mit Exit-Code `12`, vorbehaltlich des Urteils nach dem Zustand des Containers (Abbildung auf `11`, Absatz *Fallbacks und Abbruch eines Startscripts*). `up` nach `down` startet denselben Container; `up --recreate` entfernt ihn und legt einen neuen an, benannte Volumes bleiben samt Inhalt (der Workspace des Sandbox-Profils). Die Images `vsc-<Ordnername>-…`, die die CLI baut, bleiben nach `down` und nach `--recreate` bestehen.

**Lebenszyklus-Befehle bei `up` gegen einen laufenden Container.** Gegen einen laufenden Container führt die Dev-Container-CLI bei `up` das `postCreateCommand` und das `postStartCommand` nicht erneut aus, das `postAttachCommand` bei jedem `up`; nach `down` und `up` im selben Container läuft das `postStartCommand` wieder, das `postCreateCommand` nicht. Das gilt für die gepinnte Version der CLI, für Befehle in der Form einer Zeichenkette, die gelingen, unter Linux `amd64` mit lokaler Engine. Nicht zugesagt sind das Verhalten des `initializeCommand`, des `onCreateCommand` und des `updateContentCommand`, die Array- und die Objektform der Befehle, scheiternde Befehle gegen einen laufenden Container und andere Hosts und Engines. Für `shell` ist nur zugesagt, dass seine Startphase ein `up` ohne `--recreate` ist; wann ein `postAttachCommand` bei `shell` läuft, ist nicht zugesagt.

**Fallbacks und Abbruch eines Startscripts.** Nach einem Lauf der Dev-Container-CLI, der nicht an der Vorab-Prüfung scheitert und nicht abgebrochen wird, liest `up` die Statusdateien auf dem Weg von [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), und zwar je Fähigkeit, die die `u-boot.yaml` im Projektordner mit Sandbox-Profil konfiguriert. Bei `result=fallback` nach einem gelungenen Start entsteht eine Warnung (Level `warn`, Subjekt `.devcontainer/devcontainer.json`, Code [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) für `vfs` und Daemon, [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) für Egress) mit dem Kennzeichen `DegradationFallback` und ohne `SecurityRelaxation`; der Text beginnt mit „state of the devcontainer container at its last start (<Zeitpunkt>): “ und nennt die Folge des Fallbacks. Er enthält weder „existing“ noch die Anweisung zum Neubau, weil `up` keine Datei neu erzeugt. Jeder andere Zustand außer `failed` bleibt still. Ohne `u-boot.yaml`, ohne Sandbox-Profil und ohne die Fähigkeit entsteht keine Ausweisung, und ein gelungener Start endet mit Exit-Code `0`.

`up` urteilt nach dem Zustand des Containers, nicht nach dem Verlauf des Aufrufs: Meldet die Statusdatei einer konfigurierten Fähigkeit `result=failed`, endet der Befehl mit Exit-Code `11`, wie `doctor` ([`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), [`SPEC-048`](#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf), [`LH-FA-DEV-007.a`](#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten)). Das gilt nach einem Lauf der CLI, der mit `success` endete, wie nach einem Fehler der Klasse `12` (dann `11` statt `12`), und auch dann, wenn die Datei aus einem früheren Start desselben Containers stammt, denn die Datei beschreibt den letzten Lauf ihres Scripts, nicht notwendig den letzten Start des Containers ([`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), *Stand*). Ein erneutes `up` ohne `--recreate` gegen einen Container mit gescheitertem Startscript endet damit wieder mit `11`; `u-boot devcontainer up --recreate` legt den Container neu an. Melden mehrere konfigurierte Fähigkeiten `failed`, nennt die Meldung die erste in der Reihenfolge der Prüfung (`vfs` vor `egress`). Nach einem Fehler der CLI behält der Fehlertext deren Ausgabe; nach einem gelungenen Lauf nennt er den Satz der Fähigkeit (derselbe wie bei `doctor`), den Zeitpunkt der Datei als `(reported at <Zeitpunkt nach RFC 3339>)`, dass der Container läuft und dieser Aufruf das Startscript nicht ausgeführt hat, und `u-boot devcontainer up --recreate`. Fallback-Warnungen anderer Fähigkeiten entfallen bei `11`. Die Klassen `11` und `12` überschneiden sich nicht.

Das Urteil nach dem Zustand setzt eine lesbare Statusdatei einer konfigurierten Fähigkeit voraus. Ohne `u-boot.yaml`, ohne Sandbox-Profil, ohne die Fähigkeit und ohne lesbare Statusdatei (etwa bei blockierten verschachtelten User-Namespaces, bei denen keine Datei entsteht) bleibt der Exit-Code der Startphase: `0` nach einem gelungenen Lauf der CLI, `12` nach einem Fehler dieser Klasse, und ein früherer Abbruch eines Startscripts bleibt unbemerkt. Auch ein ungültiges Dokument, ein Lesefehler, ein nicht gefundener Container sowie `ok` und `fallback` einer Fähigkeit ändern den Exit-Code nicht. Fehler der Vorab-Prüfung und ein Abbruch lesen keine Statusdatei und behalten ihre Klasse. Die Lesung findet den Container über das Label des Projektordners ([`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)), mit `--instance` zusätzlich über das Label der Instanz; tragen mehrere Container dieselben Labels, ist nicht zugesagt, dass der gelesene der gestartete ist.

**Ausweisung des Sockets.** `up` und `shell` geben eine Warnung aus (Level `warn`, Code [`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Kennzeichen `SecurityRelaxation`, nicht `DegradationFallback`). Sie entsteht vor dem Aufruf des Launchers (bei `up` vor der Abfrage, ob der Container schon läuft); dessen Vorab-Prüfung der Engine, der Bau und der Start folgen danach. Bei einem späteren Fehler steht sie weiter in der Ausgabe. Sie bleibt unter `--quiet` ([`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), Ausnahme 1). Der Text ist Englisch, in Konsole und JSON gleich, und sagt fünf Dinge: (1) „this command mounts the host container-engine socket into a short-lived helper container“, (2) „root-equivalent access“ zur Engine für die Dauer des Aufrufs, bei `shell` für die Dauer des Starts („for the duration of the call“ beziehungsweise „for the duration of the start“; die Sitzung von `shell` bindet den Socket nicht ein), (3) „the devcontainer configuration that u-boot generates mounts it if a feature that does so is enabled, such as the built-in feature docker-cli (the sandbox profile rejects that feature)“, (4) „a hand-written .devcontainer/devcontainer.json can mount it“, (5) „an initializeCommand declared there runs with this access“. Der Text sagt bei (3) nicht „only“: `u-boot` liest die Metadaten eines Features aus einer externen Quelle nicht ([`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen)). Er sagt nicht, dass der gestartete Devcontainer das Socket nicht erhalte. `down` startet keinen Hilfs-Container und weist nichts aus.

**Ausgabe von `up`.**

- **Text.** Auf der Standardausgabe stehen in dieser Reihenfolge die Ausweisung des Sockets (sofort, vor dem Start), die Zeile `Devcontainer started: container <erste zwölf Zeichen der Kennung>[, name <Containername>], remote user <Benutzer>, workspace <Pfad>` und die Fallback-Warnungen und die Warnungen der Auflösung (Quelle, fehlende Quelle, Identität; sie stehen nach den Fallback-Warnungen), jede als `Warning: <Text>`. Der Name steht in der Zeile, wenn der Launcher dem Aufruf einen übergeben hat (*Containername*); der Abschnitt `--json` ändert daran nichts. War der Container schon gestartet (nächster Punkt), lautet die Zeile `Devcontainer already running: container <erste zwölf Zeichen der Kennung>[, name <Containername>], remote user <Benutzer>, workspace <Pfad>`. Der Fortschritt geht an die Standardfehlerausgabe.
- **Schon laufender Container.** Ohne `--recreate` fragt `up` nach der Ausweisung des Sockets und vor dem Aufruf des Launchers über das Label des Projektordners ([`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)), ob der Container des Projekts läuft. Lief er und meldet die CLI dieselbe Kennung, gilt er als schon laufend; das ändert nur die Zeile und das Feld `alreadyRunning`, weder den Exit-Code noch `status` noch die Warnungen. Ein gestoppter oder fehlender Container, eine andere Kennung und eine Abfrage, die scheitert, ergeben „nicht schon laufend“; ein Fehler der Abfrage ist kein Fehler des Befehls. Mit `--recreate` meldet `up` den Container nie als schon laufend: die Zeile lautet immer `started`, und `alreadyRunning` ist `false` (ohne `--instance` fragt `up --recreate` den Container dennoch ab, für die Prüfung des Absatzes *Instanzen*, und wertet den Zustand nicht aus). Die Abfrage und die Lesung der Statusdateien wählen den Container der Instanz, wenn `--instance` gegeben ist (Absatz *Instanzen*).
- **`--quiet`.** Die Erfolgszeile, auch die Zeile `already running`, und der Fortschritt entfallen; die Ausweisung und die Fallback-Warnungen bleiben ([`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Ausnahmen 1 und 2); jede andere Warnung entfällt. Schlägt der Aufruf fehl, gibt der Befehl den zurückgehaltenen Fortschritt nachträglich auf der Standardfehlerausgabe aus.
- **`--json`.** Genau ein Umschlag auf der Standardausgabe mit `command: "devcontainer"` und `subcommand: "up"`; es gibt keine Fortschrittsausgabe. `data` trägt `containerId` (die volle Kennung), `remoteUser`, `remoteWorkspaceFolder` und `alreadyRunning` (immer vorhanden, `true` nur für einen schon laufenden Container) sowie `gitTokenSource` (der Name der Quelle der Auflösung; leer ohne Token oder ohne Auflösung). `diagnostics[]` führt die Ausweisung zuerst (mit dem Text der Konsole, ohne `file`), danach die Fallback-Warnungen (mit `file`). `status` folgt [`LH-FA-CLI-007.a`](#lh-fa-cli-007a--konvention-für-diagnose-codes-und-kopplung-des-status): Der Umschlag eines gelungenen `up` trägt wegen der Ausweisung `status: "warn"` bei `exitCode: 0`. `--quiet` ändert mit `--json` nichts. Bei einem Fehler enthält `diagnostics[]` die Ausweisung (wenn sie entstanden ist), danach den Fehler (Level `error`); `data` trägt dann keine `containerId`. Der Code der Fehler-Diagnose ist [`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) für die Exit-Codes `10`, `11` und `12`, [`LH-NFA-REL-003`](lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern) für `14` und [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes) sonst; sie nennt den Pfad des Projektordners nicht.

**`down`.** Der Befehl findet den Container über das Label des Projektordners ([`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)), mit `--instance` zusätzlich über das Label der Instanz (Absatz *Instanzen*). Läuft er, ruft `u-boot` `docker stop <Kennung>` auf; der Befehl entfernt weder Container noch Volume noch Image und braucht weder Hilfs-Container noch Socket. Ohne Container und bei einem schon gestoppten Container endet er mit Exit-Code `0` und ohne `docker stop`. Die Textausgabe nennt den Fall: `Devcontainer stopped: container <zwölf Zeichen> (not removed)`, `Devcontainer container <zwölf Zeichen> is not running; nothing to stop` oder `No devcontainer container found for this project; nothing to stop`; unter `--quiet` entfällt sie. Mit `--json` trägt `data` `containerId` (leer ohne Container) und `stopped` (`false`, wenn nichts zu stoppen war), `status` ist `ok` und `diagnostics` leer; `--quiet` ändert daran nichts. Kann die Engine nicht befragt werden (kein `docker`, der Daemon antwortet nicht, die Engine ist nicht lokal), endet `down` mit Exit-Code `11`, nicht mit „kein Container“ und Exit-Code `0`; ein Fehler von `docker stop` endet mit Exit-Code `12`.

**Exit-Codes** ([`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes)); unter `--json` ist `exitCode` des Umschlags gleich dem Exit-Code des Prozesses:

| Code | `up` | `down` | `shell` |
|---|---|---|---|
| `0` | gestartet oder schon laufend | gestoppt oder nichts zu stoppen | die Sitzung hat begonnen und ist zu Ende, unabhängig vom Exit-Status der Shell (Ausnahme: die Grenzen der Erkennung im Unterabschnitt *`shell`*) |
| `2` | Aufruf (siehe *Aufruf*) | Aufruf (siehe *Aufruf*) | Aufruf (siehe *Aufruf*), auch `--json` |
| `10` | Voraussetzungen von `up` verletzt; ungültiger Name der Instanz; ein Namesteil des Containernamens verletzt die Engine-Regel; ein Container eines anderen Projektordners hält den Containernamen (*Namenskollision*); ohne `--instance` nur Container von Instanzen | ungültiger Name der Instanz; ohne `--instance` nur Container von Instanzen | Voraussetzungen verletzt, wie bei `up`; ungültiger Name der Instanz; Namesteil oder Kollision wie bei `up`; ohne `--instance` nur Container von Instanzen |
| `11` | Vorab-Prüfung der Engine gescheitert; Abbruch eines Startscripts mit `failed` in der Statusdatei (nach dem Zustand des Containers) | Vorab-Prüfung oder Befragung der Engine gescheitert | wie bei `up`; es öffnet sich keine Sitzung |
| `12` | Bau des Hilfs-Images, Start des Hilfs-Containers und Fehler der CLI (namentlich `outcome` ungleich `success`, kein Ergebnis, Fehler eines Lebenszyklus-Befehls, Bau des Dev-Containers, ein Ordner, den die Engine nicht sieht) | `docker stop` gescheitert | Fehler des Starts wie bei `up`; die Probe der Sitzung scheitert; der Client scheitert vor dem Beginn der Sitzung (Unterabschnitt *`shell`*) |
| `14` | Lesefehler von `.devcontainer/devcontainer.json` | — | Lesefehler von `.devcontainer/devcontainer.json` |
| `1` | jeder andere Fehler | jeder andere Fehler | jeder andere Fehler |

**Reichweite.** Ohne `docker` im `PATH` enden alle drei Unterbefehle mit Exit-Code `11`. Die Hosts und Engines, auf denen er zugesagt ist, und seine Grenze im ausgelieferten Container-Image nennt die Anforderung ([`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)).

**`shell`: Startphase und Sitzung.** `u-boot devcontainer shell` hat zwei Phasen. Die *Startphase* ist `up` ohne `--recreate`: Voraussetzungen (Exit-Code `10`), Ausweisung des Sockets, Vorab-Prüfung der Engine, Hilfs-Image, Hilfs-Container samt seiner Umgebung, Ergebniszeile der CLI, Fallback-Warnungen und das Urteil nach dem Zustand des Containers sind die von `up`. Es entfällt nur die Auswertung, ob der Container schon lief, weil `shell` kein `alreadyRunning` ausgibt; die Abfrage für die Prüfung des Absatzes *Instanzen* bleibt. Meldet die Statusdatei einer konfigurierten Fähigkeit `failed`, endet der Befehl mit Exit-Code `11` und öffnet keine Sitzung, auch gegen einen Container, der schon lief und dessen Startscript dieser Aufruf nicht ausführte. Die *Sitzungsphase* beginnt nach einem gelungenen Start und öffnet im gestarteten Container eine interaktive Shell mit der Kennung, dem Benutzer (`remoteUser`) und dem Arbeitsverzeichnis (`remoteWorkspaceFolder`) der Ergebniszeile.

- **Sitzung.** Der `docker`-Client des Hosts führt zwei Aufrufe aus und sonst nichts im Container: die Probe `docker exec -u <Benutzer> -w <Arbeitsverzeichnis> <Kennung> bash -c true` (ohne `-i`, `-t` und Variablen; endet sie nicht mit Status `0`, dieselbe Probe mit `sh`) und die Sitzung `docker exec -i [-t] -u <Benutzer> -w <Arbeitsverzeichnis> [-e <NAME> …] <Kennung> <Shell>`. Die Shell der Sitzung ist die erste, deren Probe mit Status `0` endet (`bash`, sonst `sh`); sie startet ohne Argument. Die Argumentliste enthält weder `-v` noch `--volume`, `--mount`, `--privileged`, `--network`, `--cap-add` oder `--detach`. Die Sitzung läuft ohne Hilfs-Container: Während der Sitzung besteht keiner, und für die Sitzung bindet der Befehl keinen Container-Runtime-Socket des Hosts ein. Zugesagt ist, was der Befehl für die Sitzung einbindet, nicht der Zustand der Shell: ein eigener Socket im Container (`nestedRuntime: docker`) und ein Socket, den die Konfiguration selbst einbindet (etwa über das Katalog-Feature `docker-cli`), sind nicht Gegenstand dieser Aussage.
- **Umgebung der Sitzung.** `u-boot` liest den Text von `.devcontainer/devcontainer.json` ohne seine Kommentare (`//` bis zum Zeilenende und `/* … */`; Zeichenketten bleiben, ein auskommentierter Eintrag gibt also nichts weiter) und sammelt jeden Eintrag der Form `"<NAME>": "${localEnv:<NAME>}"` mit demselben einfachen Variablennamen im Schlüssel und im Verweis, sortiert und ohne Duplikate, an jeder Stelle der Datei. Das ist die Form, in der die erzeugte Konfiguration unter `remoteEnv` die Zugangsdaten, die Variablen eines Agenten und die Token-Variablen eines Forge-Werkzeugs führt ([`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)). Ein Verweis anderer Form (ein anderer Name, ein Vorgabewert, Teil eines längeren Werts, ein Pfad in `mounts`) ist kein Eintrag; der Hilfs-Container der Startphase bekommt ihn dennoch wie bei `up`. Die Sitzung bekommt die Namen als `-e <NAME>` ohne Wert, und nur die, die auf dem Host gesetzt sind oder die der Start desselben Aufrufs mit der Auflösung der Token-Quelle ermittelt hat (Absatz *Auflösung der Token-Quelle*); der Client liest den Wert aus seiner Umgebung, und der Wert steht in keiner Kommandozeile. Eine auf dem Host ungesetzte Variable fehlt in der Sitzung. `u-boot` maskiert in der Sitzung nichts: Was die Shell ausgibt, erreicht den Aufrufer unverändert; die Maskierung gilt für den Fortschritt und die Fehlermeldungen der Startphase.
- **Terminal.** Die Sitzung bekommt ein Terminal (`-t`), wenn die Standardeingabe des Befehls eine Datei, ein Zeichengerät und nicht das Nullgerät ist. Mit Terminal gibt der Launcher Standardeingabe und Standardausgabe unverändert an den Client, damit der Client das Terminal dort findet. Ob die Ausgabe ebenfalls ein Terminal ist, prüft die Erkennung nicht (Eingabe Terminal, Ausgabe Pipe: nicht zugesagt). Ein anderes Zeichengerät als ein Terminal (etwa `/dev/zero`) gälte als Terminal; lehnt der Client es ab („stdin is not a terminal“), endet der Befehl mit `12`. Der Terminal-Pfad ist in drei Teilen zugesagt: die Weitergabe der Erkennung an die Sitzung, die Sitzung am gebauten Binary und die Sitzung am Launcher, beide unter einem Pseudo-Terminal; die Strecke vom Terminal des Nutzers bis zur Sitzung ist nicht als Ganzes zugesagt.
- **Ende der Sitzung und Signale.** Die Sitzung endet mit dem Ende der Shell (`exit`), der Befehl dann mit `0`. Strg-C erreicht mit Terminal die Shell im Container. Eine Shell, die einen Befehl ausführt, bleibt samt Befehl bestehen, wenn die Sitzung von außen endet: ohne Terminal bei Abbruch des Kontexts, mit Terminal bei `SIGKILL` am Client; `u-boot` räumt sie nicht auf. Eine untätige Shell ohne Terminal endet bei Abbruch des Kontexts; eine untätige Shell mit Terminal bleibt bei `SIGKILL` am Client bestehen. Mit Terminal beobachtet der Launcher den Kontext nach dem Beginn der Sitzung nicht; ohne Terminal sendet er bei einem Abbruch `SIGINT` an den Client und meldet den Abbruch als Fehler des Kontexts. Nicht zugesagt sind: dass die Shell mit `u-boot` endet, dass ein Signal an `u-boot` die Sitzung beendet (mit Terminal beendet es sie nicht), das Schließen des Terminals und Strg-C ohne Terminal.
- **Art der Shell.** Mit Terminal ist die Shell keine Login-Shell (`shopt -q login_shell` ergibt nein), `$0` ist `bash`, sie ist interaktiv, liest `~/.bashrc` und liest `~/.profile` nicht; ohne `bash` im Container ist `sh` die Shell. Nicht zugesagt sind `/etc/profile`, die Art der Shell ohne Terminal und Einstellungen, die nur Login-Shells lesen.
- **Reichweite** der Aussagen zum Ende der Sitzung, zu den Signalen und zur Art der Shell: Sie sind zugesagt für Container mit dem Basis-Image `mcr.microsoft.com/devcontainers/base:debian`, deren Hauptprozess ein Dauerlauf ist; für einen von `devcontainer up` erzeugten Container sind sie nicht zugesagt, für andere Basis-Images ebenfalls nicht. Andere Engines, Hosts und Terminals sind nicht zugesagt.
- **Ausgabe.** Alles, was `u-boot` selbst ausgibt, steht auf der Standardfehlerausgabe, in dieser Reihenfolge: die Ausweisung des Sockets, der Fortschritt der Startphase, jede Fallback-Warnung (als `Warning: <Text>`, im Wortlaut von `up`) und, danach, vor dem Beginn der Sitzung, die Zeile `Devcontainer started: container <erste zwölf Zeichen der Kennung>, name <Containername>` beziehungsweise, wenn die Startphase den Container als schon laufend meldete, `Devcontainer already running: …` in derselben Form — als Ergebnis eines ändernden Befehls also nur bei einem Container, dessen Name `u-boot` übergeben hat ([`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)); bei einem Fehler folgt dieser. Auf der Standardausgabe steht nur, was die Sitzung ausgibt; eine Erfolgszeile gibt es daneben nicht. Standardeingabe, Standardausgabe und Standardfehlerausgabe des Befehls gehen unverändert an die Sitzung.
- **`--quiet`** ([`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)). Der Fortschritt der Startphase wird zurückgehalten und nur ausgegeben, wenn der Befehl fehlschlägt; die Ausweisung des Sockets und die Fallback-Warnungen bleiben (Ausnahmen 1 und 2 von [`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging)); jede andere Warnung entfällt. Die Zeile mit der Kennung und dem Namen des Containers steht in demselben zurückgehaltenen Fortschritt: bei einem gelungenen Befehl entfällt sie, bei einem fehlgeschlagenen erscheint sie mit ihm. Die Ausgabe der Sitzung bleibt unverändert: `--quiet` wirkt nur auf die Ausgabe des Starts.
- **`--yes` und `--no-interactive`** wirken bei `shell` nicht; der Befehl fragt nichts.
- **Exit-Code `12`.** Der Befehl endet mit `12`, wenn die Probe scheitert (Benutzer, Arbeitsverzeichnis oder Shell fehlen im Container, oder der Container läuft nicht; die Meldung nennt das Ende der Ausgabe jedes Versuchs) und keine Sitzung öffnet, und wenn der Client `docker exec` vor dem Beginn der Sitzung mit einem eigenen Fehler beendet (*Grenzen der Erkennung*). Der Exit-Status der Shell nach dem Beginn der Sitzung ist kein Fehler: Er erreicht den Aufrufer nicht, auch nicht nach dem Ende des Containers während der Sitzung; der Befehl endet mit `0` (Ausnahme: die *Grenzen der Erkennung*).
- **Grenzen der Erkennung** (Fall „die Probe gelingt, die Sitzung beginnt nicht“, etwa weil der Container zwischen Probe und Sitzung endet oder Benutzer, Arbeitsverzeichnis oder Shell dort verschwinden; [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes) bindet ihn an `10`, `11` oder `12`). Der Exit-Status des Clients trennt den Fehler des Clients nicht von dem einer Shell, die mit einem fehlgeschlagenen Befehl endet (Status 1, 126 und 127 hat auch eine Shell); der Launcher erkennt den Fehler an der Meldung des Clients und meldet ihn als `12`, auch „Engine nicht erreichbar“ (`11` bleibt der Vorab-Prüfung vorbehalten). Erkannt wird:
  - eine Sitzung, deren erste nicht leere Zeile der Fehlerausgabe mit `Error response from daemon: ` beginnt und die Kennung des Containers oder den Benutzer der Sitzung nennt (Container läuft nicht oder fehlt, Benutzer fehlt), bei einem Status ungleich `0`; ohne Terminal muss die Standardausgabe leer sein, mit Terminal gehört sie dem Client;
  - die Zeile `Cannot connect to the Docker daemon at `, wenn die Engine auch bei einer erneuten Frage nicht antwortet;
  - die Zeile `cannot attach stdin to a TTY-enabled container because stdin is not a terminal`, wenn ein Terminal angefordert war;
  - ohne Terminal die Zeile `OCI runtime exec failed: ` als erste nicht leere Zeile der Standardausgabe bei leerer Fehlerausgabe und Status `127` (Arbeitsverzeichnis oder Shell nach der Probe entfernt); mit Status `1` gilt dieselbe Zeile als Sitzung;
  - ohne Terminal die Zeile `unable to find user <Benutzer der Sitzung>: …` als einzige Zeile der Standardausgabe (keine zweite nicht leere Zeile, nichts abgeschnitten) bei leerer Fehlerausgabe und Status `126`; der Doppelpunkt ist die Namensgrenze, `vscode2` ist nicht `vscode`. Das ist die Form, die der Launcher für die Engine 28.0.4 erkennt (nicht zugesagt, nicht bestätigt); er prüft sie unabhängig von der Version der Engine, die Engine 28.0.4 selbst ist nicht zugesagt (*Engine-Abhängigkeit*).

  Nicht zu trennen ist der Fall „die Sitzung schreibt die Zeile selbst“. Ein `docker exec -u <Benutzer der Sitzung> <anderer Container>` in der Sitzung schreibt auf der Engine 29.9.2 mit dem Client 29.9.2 (gemessen am 2026-11-01, Linux `amd64`) dieselbe Zeile in derselben Form auf demselben Strom wie der Client bei fehlendem Benutzer: `Error response from daemon: unable to find user U: …` auf der Fehlerausgabe. Ist dieser Befehl der letzte der Sitzung ohne Terminal, endet die Shell mit Status 1 und ist die Standardausgabe leer, erfüllt die Zeile die Bedingungen des ersten Punkts, und der Befehl endet mit `12` statt `0`; der Launcher hat keinen Anhaltspunkt, der beides trennt. Das ist eine Lesung des Codes gegen die gemessene Form, kein gefahrener Lauf des Launchers. Endet die Shell danach mit Status `0` (etwa `exit 0`), gilt die Sitzung. Das Lastenheft nennt diese Ausnahme zu `0` nicht; bei Konflikt gilt es.

  Alles andere gilt als Sitzung: eine Meldung des Clients nach anderer Ausgabe der Sitzung, auch hinter einer Warnzeile (nur die erste nicht leere Zeile zählt), die Meldung der Engine, solange sie antwortet, und jede Meldung, die nicht zu den genannten gehört. Mit Terminal ist die Erkennung auf die Meldungen der Fehlerausgabe beschränkt, und von ihnen sind „Container läuft nicht“ und „Standardeingabe ist kein Terminal“ zugesagt; die OCI-Zeile und die Zeile des fehlenden Benutzers ohne Präfix stehen im Terminal-Strom des Clients, den der Launcher nicht liest, und der Befehl endet dann mit `0`. **Engine-Abhängigkeit:** Die Meldung kommt je Engine in anderer Form. Zugesagt ist die Engine 29.9.2 mit dem Client 29.9.2 unter Linux `amd64`: Der fehlende Benutzer endet dort mit Status 1 und `Error response from daemon: unable to find user U: …` auf der Fehlerausgabe. Die Engine 28.0.4 und andere Engine-Versionen sind nicht zugesagt. Für die Engine 28.0.4 erkennt der Launcher die Form Status 126 und Zeile `unable to find user U: …` ohne Präfix auf der Standardausgabe (mit `\r\n`, letzter Punkt der Liste oben), ohne dass diese Form zugesagt ist; die Meldungen gelten für den Client 29.9.2.
- **Ohne `u-boot.yaml`.** Das Urteil nach dem Zustand setzt die `u-boot.yaml` des Projektordners voraus ([`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): Der Exit-Code `11` für den Startscript-Abbruch setzt die Statusdatei des Containers voraus). Ohne diese Datei, ohne Sandbox-Profil oder ohne die Fähigkeit liest `shell` keine Statusdatei und öffnet die Sitzung auch in einem Container, dessen Startscript abgebrochen ist; scheitert der Start selbst, endet der Befehl dann mit `12`.

### SPEC-053 — Grenzen für Arbeitsspeicher und CPU des Devcontainers: Schlüssel, Werte und `runArgs`

Verfeinert [`LH-FA-CONF-002`](lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (Schlüssel und Exit-Code `10`) und [`LH-FA-DEV-001`](lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (Erzeugung der `devcontainer.json`).

**Schlüssel.** `devcontainer.resources.memory` und `devcontainer.resources.cpus`, beide als Zeichenkette in `u-boot.yaml`, beide les- und schreibbar über `u-boot config` ([`SPEC-040`](#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)). Leer oder ohne Eintrag heißt nicht gesetzt.

**Werte.** `memory`: Ziffern, danach höchstens ein Buchstabe aus `b`, `k`, `m`, `g` (klein); der Zahlenwert ist nicht null. `cpus`: Ziffern, optional `.` und weitere Ziffern; der Wert ist nicht null. Jede andere Zeichenfolge ist ungültig, darunter Großbuchstaben, `t`, Vorzeichen, Leer- und Anführungszeichen, Rückstrich, Strichpunkt, Zeilenwechsel und ein Komma als Dezimaltrenner. Die Zeichenmenge ist geschlossen, weil der Wert wörtlich in eine JSON-Zeichenkette der `devcontainer.json` gelangt. Die Prüfung kennt die Grenzen der Engine nicht (Mindestwert von `--memory`, Obergrenze von `--cpus` nach der Zahl der CPUs des Hosts); eine Engine kann einen Wert abweisen, den `u-boot` annimmt.

**Fehler.** (Leerer String und fehlender Eintrag gelten gleich als nicht gesetzt.) Ein ungültiger Wert endet bei `config set` mit Exit-Code `10`; `generate devcontainer` bricht bei einem von Hand eingetragenen ungültigen Wert mit Exit-Code `10` ab, schreibt nichts und nennt den Schlüssel; `doctor` meldet ihn im Befund zu `u-boot.yaml` (Schema ungültig), ohne neue Prüfung. Der Konstruktor der Werte liefert den Sentinel `ErrInvalidResourceLimit`; `config set` wickelt ihn in `ErrConfigValueInvalid`, `generate` und die Ladepfade verflachen ihn zu Text innerhalb ihrer eigenen Sentinels, die den Exit-Code `10` tragen. Ein eigener Diagnose-Code besteht nicht.

**`runArgs`.** Sind Schlüssel gesetzt, trägt `devcontainer.json` `--memory=<Wert>` vor `--cpus=<Wert>`, in beiden Profilen. Im Sandbox-Profil stehen sie hinter den Optionen von `nestedRuntime` und der Egress-Restriktion, im Standardprofil trägt `runArgs` sonst nur die Zeile der Umgebungsdatei; in beiden Profilen stehen sie hinter dieser Zeile ([`SPEC-055`](#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung)) und vor der Zeile des Features `ollama-host`. Ohne Schlüssel (auch mit leerem Block `resources: {}`) ist die erzeugte Datei byte-gleich zum Lauf ohne den Block. Die Grenzen sind keine Sicherheitslockerung: Die Befehlsausgabe führt sie nicht auf, und die Meldungen sind die des Laufs ohne die Schlüssel.

**Belege.**

| Bedingung | Test |
|---|---|
| Zeichenmenge von `memory` und `cpus` | `TestNewMemoryLimit`, `TestNewCPULimit` (`internal/hexagon/domain/resources_test.go`) |
| Pfade les- und schreibbar, ungültig endet mit `ErrConfigValueInvalid` | `TestConfigSetGet_SandboxKeys`, `TestConfigSet_SandboxKeys_InvalidValue` (`config_sandbox_test.go`) |
| ungültiger Wert in `u-boot.yaml` beim Lesen: `ErrConfigSchemaInvalid` | `TestConfigGet_SandboxKeys_InvalidOnLoad` |
| `doctor` meldet ihn im Befund `uboot.yaml.valid` | `TestDoctor_UbootYaml_ErrorOnBadResourceLimit` |
| `runArgs` im Standardprofil, Reihenfolge | `TestGenerateDevcontainer_Resources_DefaultProfile` |
| `runArgs` im Sandbox-Profil, Reihenfolge zu `nestedRuntime`, Egress, `ollama-host` | `TestGenerateDevcontainer_Resources_SandboxProfile` |
| ohne Schlüssel byte-gleich | `TestGenerateDevcontainer_Resources_AbsentIsByteIdentical` |
| keine Sicherheitslockerung | `TestGenerateDevcontainer_Resources_NoRelaxationWarning` |
| ungültiger Wert bricht `generate` ab, nichts geschrieben | `TestGenerateDevcontainer_Resources_InvalidValue` |
| Grenze kommt im gestarteten Container an (Linux-Host `amd64`, lokale, nicht verschachtelte Engine; Standardprofil und Sandbox-Profil mit `nestedRuntime: docker`) | `TestE2E_Resources_DefaultProfileLimitsReachTheContainer`, `TestE2E_Resources_SandboxNestedDockerLimitsReachTheContainer` (Build-Tag `docker`) |

**Nicht zugesagt:** die Wirkung auf den verschachtelten Docker und seine Container (der Umbau davor: [`SPEC-056`](#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag)), das Verhalten bei einem zu kleinen Limit (Abbruch eines Startscripts), die Summe mehrerer Instanzen, andere Hosts und Engines, andere Kombinationen als die der beiden Docker-Tests, ob eine geänderte Grenze erst für einen neu erzeugten Container gilt (angenommen, nicht gemessen).

### SPEC-054 — Katalog-Eintragstyp „gepinnte Werkzeugkette“: `go-toolchain`, `golangci-lint` und `java-toolchain`, Prüfsummen-Quellen, `version:`-Übersteuerung und Cache-Werte unter `containerEnv`

Verfeinert [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) (Eintragstyp mit Übersteuerung und Prüfsummen-Quelle) und [`LH-FA-DEV-001`](lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (Cache-Werte unter `containerEnv`) und berührt [`LH-FA-DEV-011`](lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) (Pin-Politik der vier Stellen), [`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Wirkung von `--quiet`) und [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Code `10`). Der Katalog und die anderen Eintragstypen stehen in [`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen), die Wirkung von `--quiet` in [`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), das Feature `go` in [`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen). Die Grenzen der Zusagen stehen im Lastenheft ([`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-001`](lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)).

**Eintragstyp.** Die drei Katalog-Schlüssel `go-toolchain` (Go-Toolchain), `golangci-lint` und `java-toolchain` (JDK der Eclipse-Temurin-Distribution der LTS-Version 21) sind Dockerfile-Blocks wie `ollama`, `gh` und `glab` ([`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen)): Sie haben keine OCI-Quelle und keinen Herausgeber, erzeugen keinen Schlüssel unter `features` in `devcontainer.json`, und `u-boot doctor` meldet für sie keinen Befund der Prüfungen `devcontainer.features.allowlist` und `devcontainer.features.drift` (gepinnt mit einem Test im Umsetzungs-Slice; Beleg-Tabelle Zusage → Test unten). Sie wirken in beiden Profilen. Der aktive Block installiert sein Werkzeug im erzeugten Dockerfile aus einem gepinnten Release-Archiv — bei den zwei Go-Schlüsseln eines Upstream-Werkzeugs, bei `java-toolchain` aus der Veröffentlichung der Temurin-Binaries — und prüft die SHA-256-Prüfsumme des Archivs, bevor es entpackt wird (eine Abweichung bricht den Image-Build ab); die Version und die Prüfsummen stehen an einer Stelle im Dockerfile. Die Archive des Typs sind nicht signiert, und `u-boot` prüft keine Signatur; die Grenze der Prüfung weist der Generator aus (Meldungen unten).

**Default-Pin und Übersteuerung.** Die Version jedes Schlüssels steht als Build-Argument im Katalog-Template (Default-Pin); sie wird nur mit einem u-boot-Release geändert. Die Build-Argumente sind nach [`LH-FA-DEV-011`](lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) die eine Quellen-Stelle der Pin-Werte; das Root-`Dockerfile` des u-boot-Projekts, `repo.mk` und die erzeugte `.devcontainer/Dockerfile` folgen ihr, die vier Stellen tragen denselben Wert, und ein Test des u-boot-Projekts liest alle vier Stellen und prüft die Gleichheit. Eine `version:`-Angabe der Projektkonfiguration (`devcontainer.features.<Schlüssel>.version`, les- und schreibbar per `u-boot config`, [`SPEC-040`](#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)) ist eine Übersteuerung des Default-Pins: Der Generator liest bei der Erzeugung die Prüfsummen-Quelle des Eintrags für die gewählte Version, und der erzeugte Dockerfile-Block trägt die gelesenen Prüfsummen anstelle der des Default-Pins. Die erzeugte `.devcontainer/Dockerfile` trägt allein den übersteuerten Wert des Projekts, das die Übersteuerung setzt; die vier Stellen der Kopplung ändert die Übersteuerung nicht. Die Übersteuerung normalisiert auf die ARG-Form des Schlüssels (Go ohne, `golangci-lint` mit dem führenden `v`; die ARG-Form von `java-toolchain` ist eine Umsetzungs-Stelle der Messung und wird von einem Docker-Test gepinnt), so dass beide Schreibweisen einer Version dieselbe Quelle lesen. Die Vier-Stellen-Kopplung nach [`LH-FA-DEV-011`](lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) trägt die zwei Go-Schlüssel; `java-toolchain` ist nicht Teil der Kopplung ([`LH-FA-DEV-011`](lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) bleibt bei den zwei Go-Schlüsseln) — seine Pin-Stellen sind die Template-`ARG` (Default-Pin) und die Übersteuerung des Projekts.

**Prüfsummen-Quellen.** Jeder Schlüssel des Typs trägt seine offizielle Prüfsummen-Quelle als Angabe des Katalog-Eintrags; dieselbe Quelle liefert das Release-Archiv (Same-Origin am Archiv, Muster der Blöcke `gh`, `glab` und `ollama`, [`SPEC-047`](#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen), [`SPEC-051`](#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten)):

| Schlüssel | Offizielle Prüfsummen-Quelle |
|---|---|
| `go-toolchain` | die Download-Listung von `go.dev/dl`, je gewählter Version |
| `golangci-lint` | die `CHECKSUMS`-Datei des Releases im Release-Verzeichnis von `golangci/golangci-lint` |
| `java-toolchain` | die Prüfsummen-Datei derselben Veröffentlichung der Temurin-Binaries, aus der das Release-Archiv stammt, je gewählter Version |

Die JSON-Schnittstelle `api.adoptium.net` ist für `java-toolchain` ausdrücklich nicht die Quelle: sie liefert die Prüfsummen nicht aus derselben Veröffentlichung wie das Archiv (kein Same-Origin am Archiv); der Abgleich der Übersteuerung liest nur die Quelle der Tabelle. Ob die Prüfsummen-Datei einer Veröffentlichung eine beliebige nutzergewählte Version abdeckt, liest die Messung der Umsetzung; deckt sie eine gewählte Version nicht, bleibt die Abweisung des Abschnitts *Abweisungen*. Die URL-Form des Abrufs je Version und der Meldungs-Wortlaut der Quelle (der Katalog-Eintrag verdoppelt die Angabe als Meldungstext der Erzeugung) sind Umsetzungs-Stellen dieses Eintragstyps; die URL-Form wird je Schlüssel von einem Docker-Test auf die drei Schlüssel gepinnt, für andere Umgebungen als die der Messung (CLI, Engine, Linux `amd64`) ist der Abruf nicht zugesagt.

**Abweisungen (Exit-Code `10`, nichts geschrieben).** Ist die Quelle nicht lesbar oder überdeckt sie die gewählte Version nicht, schlägt die Erzeugung der Übersteuerung mit einem fachlichen Validierungsfehler fehl ([`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes)); die nicht lesbare Quelle ist kein Feature-Source-Fehler ([`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen), Absatz *Freigabe der Quelle*), die Eingabe der Übersteuerung selbst ist unbrauchbar — ein Wert mit einem Zeichen außerhalb `A-Z a-z 0-9 . _ + -` ist dieselbe unbrauchbare Eingabe und wird vor dem Lesen abgewiesen. Ein Eintrag des Typs ohne Angabe seiner offiziellen Prüfsummen-Quelle weist eine `version:`-Angabe mit einem fachlichen Validierungsfehler ab (Muster der bestehenden `version:`-Abweisung, [`SPEC-047`](#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen) und [`SPEC-051`](#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten)); die drei Schlüssel des Katalogs tragen ihre Quelle nach der Tabelle oben. Die gleichzeitige Aktivierung des externen Features `go` und des Schlüssels `go-toolchain` ist ein fachlicher Validierungsfehler (das externe `go` installiert eine zweite Go-Installation, welche die Befehle finden, hängt an der Pfad-Reihenfolge); das externe `go` allein bleibt unverändert wählbar, und die Katalog-Eigenschaften von `go` ([`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen)) ändern sich nicht. Die gleichzeitige Aktivierung des externen Features `java` und des Schlüssels `java-toolchain` ist ein fachlicher Validierungsfehler ([`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes): Exit-Code `10`; das externe `java` installiert ein zweites JDK, welche die Befehle finden, hängt wie im Go-Fall an der `PATH`-Reihenfolge); das externe `java` allein bleibt unverändert wählbar, und die Katalog-Eigenschaften von `java` ([`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen)) ändern sich nicht.

**Cache-Werte unter `containerEnv`.** Sind `devcontainer.features.go-toolchain.enabled` oder `devcontainer.features.golangci-lint.enabled` gesetzt, trägt die erzeugte `devcontainer.json` unter `containerEnv` die Cache-Verzeichnisse `GOCACHE` und `GOMODCACHE`, jeweils unterhalb von `${containerWorkspaceFolder}`; sind beide Schlüssel gesetzt, trägt die Datei die Werte genau einmal — der Generator ergänzt die Einträge des zweiten Blocks nicht ein zweites Mal. Die konkreten Unterpfade sind eine Umsetzungs-Stelle und werden von einem Docker-Test gepinnt ([`LH-FA-DEV-001`](lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)). Die bestehende Render-Stelle `remoteEnv` mit den `${localEnv:…}`-Referenzen bleibt unverändert und trägt diese Werte nicht; die Auflösung von `containerEnv` zu den Metadaten (`Config.Env`) des Containers steht in [`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), Absatz *Grenzen*. Der Schlüssel `java-toolchain` trägt keinen Wert unter `containerEnv`; `JAVA_HOME` und `PATH` stehen als `ENV`-Zeilen in seinem Dockerfile-Block ([`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features)) — die Installations-Umgebung übersteht den Image-Neubau mit dem Bild, nicht das Neuanlegen des Containers —, und der konkrete `JAVA_HOME`-Wert (Installationspfad der entpackten Distribution) und die Form des Entpackens sind Umsetzungs-Mess-Stellen, die der Docker-Test pinnt (`java -version` und `JAVA_HOME` im gebauten Container).

**Meldungen des Erzeugens** (Erzeuger: `generate devcontainer`; Level `warn`; je aktiviertem Schlüssel des Typs eine Meldung; Regeln der maschinenlesbaren Ausgabe nach [`LH-NFA-USE-004`](lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), Muster der Ollama- und Forge-CLI-Meldungen, [`SPEC-047`](#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen), [`SPEC-051`](#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten)). Der Eintrag trägt die Anforderung der Lieferkette als `code` ([`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)) und nennt: den aktivierten Block des Schlüssels, ob die Version der Default-Pin oder eine Übersteuerung ist und, bei einer Übersteuerung, die Prüfsummen-Quelle; dazu die Grenze der Prüfung — die Archive des Typs sind nicht signiert, und `u-boot` prüft keine Signatur. Die Meldung des Schlüssels `java-toolchain` trägt daneben die Vertriebs-Form (die Eclipse-Temurin-Distribution der LTS-Version 21); sie konkretisiert die gemeldete Quelle, sie ersetzt sie nicht. Die Meldungen tragen weder `SecurityRelaxation` noch `DegradationFallback`; sie entfallen unter `--quiet` ([`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging)), und [`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) bekommt keine fünfte Ausnahme.

**Beleg-Tabelle.** Jede Bedingung der Abschnitte oben wird von einem Test der Umsetzung gepinnt; die Spalte nennt den Test als Beleg-Auftrag, den Namen trägt der Umsetzungs-Diff. Docker-Bedingungen mit Build-Tag `docker` ([`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features)).

| Bedingung | Pinnender Test (Umsetzung: die Namen der Beleg-Aufträge) |
|---|---|
| Die zwei Schlüssel sind als eingebaute Features in beiden Profilen wählbar | `TestFeatureCatalogue_ToolchainChecksumSources` und `TestFeatureCatalogue_KeysCoverSpecExamples`/`TestFeatureCatalogue_EntriesAreValid` (Katalog-Tabelle samt Prüfsummen-Quelle), `TestGenerateDevcontainer_Toolchain_BlocksInBothProfiles` (`internal/hexagon/application/generate_toolchain_test.go`; je Profil erzeugtes `devcontainer.json` ohne Schlüssel unter `features`, Dockerfile mit dem Block, NoOp des zweiten Laufs) |
| Block installiert aus gepinntem Archiv, Prüfsumme vor dem Entpacken | Docker-Tests `TestE2E_Toolchain_BuildsAndReportsPins` und `TestE2E_Toolchain_PinMismatchBreaksBuild` (`internal/e2e/toolchain_docker_test.go`, `make test-docker`): der Bau läuft, die Werkzeuge melden die gepinnten Versionen; je Architektur-Schlüssel `GO`/`GOLANGCI_LINT` eine Negativprobe: eine verfälschte Prüfsumme bricht den Bau vor dem Entpacken ab, der Zwischencontainer ist geräumt |
| Default-Pin und Vier-Stellen-Gleichheit ([`LH-FA-DEV-011`](lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)) | `TestToolchainPinPlacesEqual` liest die vier Stellen (Root-`Dockerfile`, `repo.mk`, Katalog-Template, erzeugte `.devcontainer/Dockerfile` des eigenen Projekts) und prüft die Gleichheit |
| Übersteuerung: Quelle wird gelesen, der Block trägt die gelesenen Prüfsummen | `TestGenerateDevcontainer_Toolchain_OverrideRendersReadChecksums` (Render-Seite mit Fixture-Quellen, Double-Run re-liest) und Docker-Test `TestE2E_Toolchain_OverrideBuildsAgainstReadChecksums` (`make test-docker`): der Bau mit Übersteuerung gegen die für die gewählte Version gelesene Prüfsumme läuft, die Werkzeuge melden die übersteuerten Versionen |
| Quelle nicht lesbar oder überdeckt die Version nicht: Abbruch Exit `10` | `TestGenerateDevcontainer_Toolchain_UnreadableChecksumSourceExit10` und `TestGenerateDevcontainer_Toolchain_VersionNotCoveredExit10`: Abbruch mit der Klasse von Exit-Code `10` (Wrap von `ErrConfigValueInvalid`, [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes)), nichts geschrieben |
| Eintrag ohne Quellen-Angabe weist `version:` ab, Exit `10` | `TestFeatureCatalogue_ToolchainChecksumSources` (Regel-Brücke `ToolchainVersionOverrideAcceptedForTest`: ein Eintrag des Typs ohne Quellen-Angabe nimmt keine `version:` an; ein `version:` ohne `source:` an `gh`/`glab`/`ollama` bleibt durch die bestehende Regel abgewiesen, `TestRun_Generate_ForgeCLINotices`) |
| `go` (extern) plus `go-toolchain`: Abbruch Exit `10`; externes `go` allein wählbar | `TestGenerateDevcontainer_Toolchain_GoFeatureConflictExit10` (die Kombination endet mit der Klasse von Exit-Code `10` und nichts geschrieben; allein aktiviertes externes `go` erzeugt unverändert) und `TestRun_Generate_ToolchainConflictExits10` (`cmd/uboot/toolchain_test.go`; am Nutzerpfad, mit `--dry-run --json` und Code [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features)) |
| `containerEnv` mit `GOCACHE`/`GOMODCACHE` unter `${containerWorkspaceFolder}`, beide Schlüssel → Werte genau einmal | `TestGenerateDevcontainer_Toolchain_JSONGolden` (Byte-Golden der erzeugten `devcontainer.json` mit `go-toolchain`, beide Schlüssel → beide Werte genau einmal, ohne die Schlüssel byte-identisch) |
| Cache liegt unter dem Workspace-Volume und übersteht das Neuanlegen (nur Sandbox-Profil) | Docker-Test `TestE2E_Toolchain_CacheSurvivesRecreate` (Marker + 64-KiB-Deck im Cache, Recreate über `u-boot devcontainer up --recreate`, aufgelöste Umgebung, Eigentümer `vscode:vscode`; Linux-Host `amd64`, lokale, nicht verschachtelte Engine) |
| Warn-Meldungen je aktiviertem Schlüssel samt Grenz-Hinweis; keine unter `--quiet` | `TestGenerateDevcontainer_Toolchain_Notices` (je Schlüssel ein `warn`-Eintrag mit Code [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), ohne `SecurityRelaxation`/`DegradationFallback`; bei Übersteuerung Quelle und Grenze) und `TestRun_Generate_ToolchainNotices` (`cmd/uboot/toolchain_test.go`: beide Profile, unter `--quiet` keine) |
| `u-boot doctor` meldet für die drei Schlüssel keinen Befund der Prüfungen `devcontainer.features.allowlist` und `devcontainer.features.drift` | `TestDoctor_ToolchainKeys_NoFindings` ((mit den drei aktivierten Schlüsseln und ohne) `features`-Schlüssel meldet `doctor` keine dieser Prüfstellen) |
| Der Schlüssel `java-toolchain` ist als eingebautes Feature in beiden Profilen wählbar | Tabellentest des Katalogs samt der Quellen-Angabe: `TestFeatureCatalogue_KeysCoverSpecExamples`, `TestFeatureCatalogue_EntriesAreValid/java-toolchain` (`devcontainer_features_test.go`); je Profil erzeugtes `devcontainer.json` ohne Schlüssel unter `features` und erzeugtes Dockerfile mit dem Block: `TestGenerateDevcontainer_Toolchain_BlocksInBothProfiles` (Profil-Kombinationen) |
| Der JDK-Block baut aus dem gepinnten Temurin-Archiv; eine verfälschte Prüfsumme bricht den Bau vor dem Entpacken ab | Docker-Tests: `TestE2E_Toolchain_BuildsAndReportsPins` (der Bau läuft, `java -version` meldet das gepinnte Release, `JAVA_HOME` zeigt auf die entpackte Distribution) und `TestE2E_Toolchain_PinMismatchBreaksBuild/jdk` (eine verfälschte Prüfsumme bricht den Bau vor dem Entpacken ab; `--force-rm` — der Zwischencontainer ist geräumt) |
| `version:`-Übersteuerung des Schlüssels liest die Prüfsummen-Datei derselben Veröffentlichung und trägt die gelesenen Prüfsummen | Render-Test mit Fixture-Quellen: `TestGenerateDevcontainer_Toolchain_OverrideRendersReadChecksums` (drei Schlüssel, Double-Run NoOp) samt `TestGenerateDevcontainer_Toolchain_OverrideSpellings/java spellings`; Docker-Test: `TestE2E_Toolchain_OverrideBuildsAgainstReadChecksums` (der Bau mit `21.0.9+10` gegen die gelesenen Prüfsummen) |
| Quelle nicht lesbar oder überdeckt die Version nicht: Abbruch Exit `10` (`java-toolchain`) | `TestGenerateDevcontainer_Toolchain_VersionNotCoveredExit10` (zweiter Lauf, `java-toolchain` ohne Archiv-Zeile) und `TestGenerateDevcontainer_Toolchain_UnreadableChecksumSourceExit10` (dritter Lauf, `java-toolchain` mit lesender Fehlerquelle); je Fall der exit-10-Wrapper, nichts geschrieben |
| `java` (extern) plus `java-toolchain`: Abbruch Exit `10`; externes `java` allein wählbar | `TestGenerateDevcontainer_Toolchain_JavaFeatureConflictExit10` (die Kombination: Klasse Exit `10`, nichts geschrieben; `java` allein erzeugt unverändert; `java` neben `go-toolchain` erzeugt; kein `sdkman`-Fall) und `TestRun_Generate_JavaToolchainConflictExits10` (`cmd/uboot/toolchain_test.go`, am Nutzerpfad mit `--dry-run --json` und Code [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features)) |
| `java-toolchain` trägt keinen Wert unter `containerEnv`; `JAVA_HOME` und `PATH` stehen im Block | Golden Case: `TestGenerateDevcontainer_Toolchain_JSONGolden/golden with java-toolchain carries no containerEnv values` (Byte-Golden ohne `containerEnv`) und `assertNoContainerEnv` in den Lauf-Tests; Docker-Test mit `java -version` und `JAVA_HOME` im gebauten Container (`TestE2E_Toolchain_BuildsAndReportsPins`) |
| Warn-Meldung des Schlüssels `java-toolchain` samt Vertriebs-Form und Grenz-Hinweis; keine unter `--quiet` | Test der Warnungsausgabe: `TestGenerateDevcontainer_Toolchain_Notices` (drei Schlüssel, Code [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), ohne `SecurityRelaxation`/`DegradationFallback`) und `TestRun_Generate_ToolchainNotices` (beide Profile, unter `--quiet` keine) |
| `u-boot doctor` meldet für `java-toolchain` keinen Befund der Prüfungen `devcontainer.features.allowlist` und `devcontainer.features.drift` | `TestDoctor_ToolchainKeys_NoFindings` (drei Schlüssel, beide Prüfungen OK) |

Die URL-Form der Lese-Stelle wird je Schlüssel gepinnt (`TestToolchainChecksumParsing`: Abschnitt **url forms**; der Docker-Test der Übersteuerung liest dasselbe Formular am echten Netz), und die Normalisierung der Schreibweisen trägt `TestGenerateDevcontainer_Toolchain_OverrideSpellings` (je Schlüssel die abweichende Schreibweise, der erzeugte Block die ARG-Form) neben dem Parser-Pin derselben Normalisierung in beiden Schreibweisen. Der Fall `golangci-lint` allein — ohne `go-toolchain` und ohne das externe Feature `go` — liest der Umsetzungs-Diff als Beobachtung: der Bau läuft, `golangci-lint version` meldet die Version; für den Lauf des Lints ohne Go-Installation steht das in dieser Tabelle mit keiner Zusage (`TestE2E_Toolchain_GolangciLintAloneBuilds`). Die URL-Form des Lese-Abrufs von `java-toolchain` trägt der Umsetzungs-Diff mit seiner Lese-Stelle neben `TestToolchainChecksumParsing`; die Java-Zeilen der Beleg-Tabelle füllt der Umsetzungs-Diff mit den gelieferten Test-Namen.

**Nicht zugesagt:** das Arbeitsverhalten der Werkzeuge in einer übersteuerten Version, das Überstehen des Neubaus des Images durch die Installationen, soweit kein Test es belegt, die Umgebungswerte `BUILD_MODE`, `GOFLAGS` und `GOTOOLCHAIN` (die Schlüssel setzen sie nicht) und, für die Cache-Werte, das Standardprofil sowie andere Hosts, Engines und Basis-Images als die der Docker-Tests.

### SPEC-055 — Umgebungsdatei des Devcontainers: Pfad, Inhalt, Reihenfolge in `runArgs`, Anlage und Prüfung

Verfeinert [`LH-FA-DEV-012`](lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers) (Umgebungsdatei, Anlage, Prüfung) und berührt [`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) und [`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (die Ausnahme der Zusage „keine Zugangsdaten in einer erzeugten Datei“ für die benutzergeführte Datei) sowie [`LH-FA-INIT-005`](lastenheft.md#lh-fa-init-005--überschreibschutz) (Überschreibschutz). Die Reihenfolge der übrigen `runArgs` steht in [`SPEC-053`](#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) und [`SPEC-052`](#spec-052--feature-ollama-host-im-devcontainer-eintragstyp-container-option-meldung-und-grenzen), die Registry der Prüfungen im Abschnitt *Externe Verträge*.

**Pfade.** Die Umgebungsdatei ist `.devcontainer/devcontainer.env`, die Beispieldatei `.devcontainer/devcontainer.env.example`; beide liegen relativ zum Projektordner. `generate devcontainer` und `init --devcontainer` behandeln beide Dateien gleich.

**Zeile in `runArgs`.** Der Generator trägt in beiden Profilen immer — auch ohne Features und ohne `devcontainer.resources.*` — den Eintrag `--env-file=${localWorkspaceFolder}/.devcontainer/devcontainer.env` als **ersten** Eintrag von `runArgs` ein. `${localWorkspaceFolder}` steht wörtlich in der erzeugten Datei; `u-boot` löst die Variable nicht auf. Die Reihenfolge von `runArgs` ist damit:

- Standardprofil: Zeile der Umgebungsdatei, `--memory=<Wert>`, `--cpus=<Wert>` (wenn gesetzt), Zeile von `ollama-host` (wenn aktiv). Ohne Schlüssel und ohne `ollama-host` ist `runArgs` die Liste mit genau diesem einen Eintrag.
- Sandbox-Profil: Zeile der Umgebungsdatei, dann die Optionen von `nestedRuntime` (`podman`: die Optionen von [`LH-FA-DEV-007.b`](#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman); `docker`: `--privileged`), `--cap-add=NET_ADMIN` mit Egress-Restriktion, `--memory`, `--cpus`, Zeile von `ollama-host`.

Alle Stellen dieser Spezifikation, die den Inhalt von `runArgs` beschreiben, nennen die Optionen ihres Gegenstands und zählen die Zeile der Umgebungsdatei nicht mit, wo sie es nicht ausdrücklich sagen. Ein zweiter Lauf ändert die Datei nicht.

**Verhalten der Engine (gemessen 2026-11-04, Docker Engine 29.9.3, `linux/amd64`, verschachtelt, ein Lauf je Fall, nur `docker run --env-file`).** Eine fehlende Datei und ein Präfix `export` beenden den Start mit Exit-Code `125`; eine Datei nur aus Kommentarzeilen startet mit Exit-Code `0`; `NAME=WERT` setzt den Wert samt Leerzeichen, Anführungszeichen bleiben Teil des Werts, `NAME=` setzt einen leeren Wert; ein mit `-e` gesetzter Wert gewinnt gegen denselben Namen in der Datei; ein relativer Pfad wird gegen das Arbeitsverzeichnis des Clients aufgelöst. Daraus folgen die absolute Form der Zeile und die Anlage der Datei.

**Anlage der Umgebungsdatei.**

- Fehlt `.devcontainer/devcontainer.env`, legt der Befehl sie an. Der Inhalt besteht nur aus Leerzeilen und Kommentarzeilen (`#`), ohne Zuweisung und ohne Markierung eines verwalteten Blocks; die Kommentare sind Englisch und nennen das Format (`NAME=WERT`, Anführungszeichen sind Teil des Werts, kein Präfix `export`) und den Namen der Beispieldatei.
- Eine vorhandene Datei — gleich welchen Inhalts, auch leer — bleibt byte-gleich: der Befehl liest sie nicht, schreibt sie nicht und legt keine Sicherung an. Das gilt auch für `init --devcontainer` mit `--force` und `--backup`; die Datei bricht dort weder den Aufruf ab noch wird sie gesichert.
- Die Datei ist Teil des Plans: Schlägt ein Schritt der Erzeugung fehl, nimmt der Befehl auch die neu angelegte Datei zurück (dasselbe Verfahren wie für die übrigen Dateien des Befehls). Die Aktion der Antwort folgt der Regel der übrigen Dateien: Ist die neu angelegte Datei die einzige Änderung, ist sie `created`.
- `--dry-run` schreibt nichts. Die Liste der geplanten Änderungen nennt die Umgebungsdatei, wenn sie angelegt würde, und nennt sie nicht, wenn sie vorhanden ist.
- Der Befehl erzeugt wegen der Umgebungsdatei keine Warnung und gibt keinen Inhalt aus.

**Beispieldatei.** `.devcontainer/devcontainer.env.example` ist eine verwaltete Datei mit dem Block `init` (Kommentarform `#`, wie das `Dockerfile`): Der Befehl legt sie an, wenn sie fehlt, und ersetzt bei einem Lauf nur den Block, Zeilen außerhalb des Blocks bleiben. Der Block enthält mindestens eine Zuweisung; jeder Wert einer Zuweisung beginnt mit `CHANGEME_`. Kein Wert eines Projekts oder einer Umgebung steht darin.

**`.gitignore`.** `u-boot init` trägt die Zeile `.devcontainer/devcontainer.env` in die erzeugte `.gitignore` ein, mit und ohne `--devcontainer`; die Zeile der Beispieldatei steht dort nicht. `generate devcontainer` schreibt die `.gitignore` nie.

**Prüfung `devcontainer.envfile`** (`u-boot doctor`; Registry im Abschnitt *Externe Verträge*). Der Befund ist `warn`, wenn `.devcontainer/devcontainer.json` die Zeile der Umgebungsdatei (`--env-file=${localWorkspaceFolder}/.devcontainer/devcontainer.env`) trägt und die Datei nicht existiert; die Meldung nennt den Pfad, der Hinweis `u-boot generate devcontainer`. Ist die Datei vorhanden, ist der Befund `ok`. Fehlt `.devcontainer/devcontainer.json`, ist sie nicht lesbar oder trägt sie die Zeile nicht, ist der Befund `ok` mit der Meldung, dass die Prüfung übersprungen ist. Die Prüfung fragt nur die Existenz: Sie liest den Inhalt nicht, fragt Git nicht und prüft die Ignorierung nicht. Sie ist nie ein Fehler und trägt weder `SecurityRelaxation` noch `DegradationFallback`; unter `--quiet` entfällt der Warn-Eintrag ([`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)). Die Prüfung `devcontainer.sandbox.credentials` liest die Umgebungsdatei nicht.

**Belege.** Die Test-Namen sind die des Umsetzungs-Diffs (Pakete: `internal/hexagon/application` ohne Angabe, sonst genannt).

| Bedingung | Test |
|---|---|
| erster Eintrag von `runArgs`, wörtlich, Standardprofil (nur diese Zeile; mit `resources` und `ollama-host` in der Reihenfolge der Liste oben) | `TestGenerateDevcontainer_EnvFile_RunArgsDefaultProfile` |
| Reihenfolge im Sandbox-Profil mit `nestedRuntime: podman`, `nestedRuntime: docker` plus Egress-Restriktion, `resources` und `ollama-host` | `TestGenerateDevcontainer_EnvFile_RunArgsSandboxProfileOrder` |
| zweiter Lauf ändert nichts; gelöschte Datei wird neu angelegt, Aktion `created` | `TestGenerateDevcontainer_EnvFile_SecondRunNoOp`, `TestGenerateDevcontainer_EnvFile_DeletedIsRecreated` |
| Anlage: nur Leer- und Kommentarzeilen, keine Zuweisung, keine Block-Markierung | `TestGenerateDevcontainer_EnvFile_CreatedWithCommentsOnly` |
| vorhandene Datei (mit Inhalt und leer) bleibt byte-gleich, wird nie gelesen, geschrieben oder gesichert | `TestGenerateDevcontainer_EnvFile_ExistingNeverTouched` (lesezählender Fake am Pfad) |
| Rücknahme der neu angelegten Datei bei einem Fehler des Befehls | `TestGenerateDevcontainer_EnvFile_RollbackRemovesCreatedFile` |
| `--dry-run`: nichts geschrieben, Liste nennt die Datei nur, wenn sie fehlt | `TestGenerateDevcontainer_EnvFile_DryRun` |
| keine Warnung wegen der Datei | `TestGenerateDevcontainer_EnvFile_NoWarning` |
| Beispieldatei: verwalteter Block, `CHANGEME_`-Werte, Block-Ersatz lässt Zeilen außerhalb stehen | `TestGenerateDevcontainer_EnvFile_ExampleManagedBlock` |
| `generate devcontainer` schreibt die `.gitignore` nicht | `TestGenerateDevcontainer_EnvFile_DoesNotWriteGitignore` |
| `init --devcontainer`: legt Datei und Beispieldatei an, `runArgs` wie `generate` (Anti-Drift der zwei Pfade) | `TestInit_Devcontainer_EnvFileCreatedAndRunArgsMatchGenerate` (`initproject_envfile_test.go`) |
| `init --devcontainer --force --backup` mit vorhandener Datei: kein Abbruch, byte-gleich, keine Sicherung | `TestInit_Devcontainer_ExistingEnvFileKept` |
| `.gitignore` von `init` trägt die Zeile, mit und ohne `--devcontainer`, nicht die der Beispieldatei | `TestInit_Gitignore_IgnoresEnvFileNotExample` |
| `doctor`: `warn` bei fehlender Datei mit Zeile; `ok` bei vorhandener; übersprungen ohne Zeile und ohne `devcontainer.json` | `TestDoctor_EnvFile_WarnsWhenMissing`, `TestDoctor_EnvFile_OKWhenPresent`, `TestDoctor_EnvFile_SkippedWithoutReference` |
| `doctor` liest die Datei nicht (lesezählender Fake) | `TestDoctor_EnvFile_NeverReadsContent` |
| `doctor`-Credential-Prüfung liest sie nicht (Sandbox-Profil, Datei mit einem Muster einer Klartext-Zugangsdate bleibt ohne Fund) | `TestDoctor_SandboxCredentials_IgnoresEnvFile` |
| `--quiet` unterdrückt den Warn-Eintrag | `TestDoctor_EnvFile_QuietDropsWarning` (`internal/adapter/driving/cli`) |
| Registry: Code in `DefaultAllowedCodes` und in der Tabelle | der bestehende symmetrische Abgleichstest (Zeile in `coderegistry.go`) |
| Variable aus der Datei ist im Container gesetzt (Start mit `docker run` und den erzeugten `runArgs`; Linux-Host `amd64`) | `TestE2E_EnvFile_VariableReachesContainer` (Build-Tag `docker`) |
| Variable aus der Datei ist im Container gesetzt, den `u-boot devcontainer up` startet (`${localWorkspaceFolder}` löst im Hilfs-Container zum Host-Pfad auf) | `TestE2E_EnvFile_VariableReachesContainerViaDevcontainerUp` (Build-Tag `docker`) |

**Messgrenze der zwei Docker-Tests.** Beide liefen in der verschachtelten Engine des Entwicklungs-Devcontainers (Docker Engine 29.9.3, ein Lauf) und auf dem GitHub-Runner `ubuntu-latest` im CI-Job `integration-docker` (ein Lauf, beide `--- PASS`; die Engine-Version des Runners ist nicht erhoben). Die Messgrenze verengt nichts und erweitert nichts: Die Zusage steht in [`LH-FA-DEV-012`](lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers), sie gilt nur, soweit diese Tests sie belegen.

**Zugesagt über die Tabelle hinaus:** der Start über `u-boot devcontainer up` gehört zur Zusage-Menge von [`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), soweit der Docker-Test der Tabelle ihn belegt. **Nicht zugesagt:** der Start über den Editor (ein Editor-Lauf ist nicht automatisierbar, ohne Test); andere Engines und Hosts als die Zusage-Menge; die Reihenfolge gegen `remoteEnv` und `containerEnv`; Dateien mit `CRLF` und mehrzeilige Werte; die Ignorierung der Datei in einem Projekt, dessen `.gitignore` nicht von `u-boot init` stammt; der Git-Status der Datei; der Dateimodus der angelegten Datei.

### SPEC-056 — cgroup-Umbau im Startscript `docker-init.sh`: Bedingung, Schritte und Fehlschlag

Verfeinert [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Container-Runtime im Sandbox-Devcontainer) und [`LH-FA-DEV-001`](lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (erzeugte Dateien) und berührt [`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Statusdatei `daemon`) und [`SPEC-053`](#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) (Grenzen).

**Anlass (gemessen 2026-11-04 an einem Container, dem Entwicklungs-Devcontainer dieses Repos).** Dort war die Wurzel-cgroup `domain threaded` (Prozesse in der Wurzel und `cpuset cpu pids` in `cgroup.subtree_control`); ein verschachtelter Start mit `--memory` oder `--cpus` endete mit Exit-Code `125` (`cannot enter cgroupv2 … in threaded mode`). Nach dem Umbau unten zeigte `docker run --memory=64m --cpus=0.5` `memory.max = 67108864` und `cpu.max = 50000 100000`. Andere Hosts und Engines sind nicht gemessen.

**Ort.** Das erzeugte `docker-init.sh` (nur bei `nestedRuntime: docker`) führt den Umbau bei jedem Lauf aus, der `dockerd` startet: nach der Bereitschaftsprüfung und dem Entfernen verwaister PID- und Socket-Dateien, unmittelbar vor dem Start von `dockerd`. Der Zweig wird nur betreten, wenn kein laufender `dockerd` erkannt wird (Bereitschaftsprüfung und PID-Datei); ein bereiter oder startender Daemon löst keinen Umbau aus. Ein `dockerd`, den das Skript nicht erkennt (etwa ohne PID-Datei), ist die Ausnahme: der Umbau liefe dann unter ihm. Die Skripte `sandbox-init.sh` und `egress-init.sh` und das Standardprofil tragen keinen Umbau. Der Ort ist `postStartCommand`; der `ENTRYPOINT` des Images wird nicht verwendet.

**Bedingung.** Der Umbau läuft nur, wenn `/sys/fs/cgroup/cgroup.controllers` lesbar ist (cgroup v2), `memory` in `cgroup.controllers` steht und `memory` nicht in `/sys/fs/cgroup/cgroup.subtree_control` der Wurzel steht. Sonst geschieht nichts, ohne Ausgabe.

**Schritte** (POSIX-`sh`, Reihenfolge fest):

1. `-<Controller>` für jeden Eintrag von `cgroup.subtree_control` der Wurzel in `cgroup.subtree_control` schreiben (im Normalfall `-cpuset -cpu -pids`); der Schritt entfällt, wenn der Inhalt leer ist.
2. Unter-cgroup `/sys/fs/cgroup/init` anlegen und jede Prozess-ID aus `cgroup.procs` der Wurzel in `init/cgroup.procs` schreiben; ein Fehler beim Verschieben eines einzelnen Prozesses (ein beendeter oder ein nicht verschiebbarer Prozess) wird ignoriert. Ein nicht verschiebbarer Prozess zeigt sich erst als Fehlschlag von Schritt 3 (`enable controllers`).
3. Alle Controller aus `cgroup.controllers` als `+<Controller>` in `cgroup.subtree_control` der Wurzel schreiben.

**Fehlschlag.** Scheitert Schritt 1 oder 3 oder das Anlegen von `init`, schreibt das Skript genau eine Zeile `u-boot sandbox: cgroup preparation failed (<Schritt>); limits (--memory, --cpus) may not reach nested containers` auf stderr, bricht den Umbau ab und startet `dockerd` wie ohne Umbau. Der Exit-Code des Skripts und die Statusdatei `daemon` ([`SPEC-045`](#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)) sind die des Daemon-Starts: der Umbau schreibt keine Statusdatei, hat keinen eigenen Wert und folgt nicht der Richtlinie `onUnavailable`. Ein Fehlschlag ist keine Degradationsstufe der Container-Runtime.

**Grenze.** Zugesagt ist die erzeugte Datei und ihr Verhalten nach der Tabelle unten. Dass die Grenzen in jedem Container ankommen, ist nicht zugesagt (gemessen: ein Container). Ungemessen ist, ob die Wurzel nach dem Umbau neue Prozesse aufnimmt (`docker exec`, neue Terminals der Erweiterung; nach Kenntnis des Kernels kann eine Domain-cgroup mit gesetztem `subtree_control` das verweigern); der Docker-Test unten führt dazu ein `docker exec` nach dem Umbau aus. Die Prozesse der Erweiterung werden mit in `init` verschoben. Schritt 2 verschiebt eine Momentaufnahme von `cgroup.procs`: ein Prozess, der danach in der Wurzel entsteht (von einem noch nicht verschobenen Elternprozess oder per `docker exec`), bleibt dort und lässt Schritt 3 mit `EBUSY` scheitern; es gibt keinen Wiederholungsversuch, die Folge ist die Zeile des Fehlschlags.

**Belege** (Dateien in `internal/hexagon/application` und `internal/e2e`):

| Bedingung | Test |
|---|---|
| Umbau steht im Zweig des `dockerd`-Starts, nach der Bereitschaftsprüfung und vor `dockerd` | `TestGenerateDevcontainer_SandboxDocker_CgroupPrepBeforeDaemon` (`generate_sandbox_test.go`) |
| Bedingung: cgroup v2, `memory` in `cgroup.controllers`, `memory` nicht in `cgroup.subtree_control`; Schrittfolge `-cpuset -cpu -pids`, `init`, `+Controller` | `TestGenerateDevcontainer_SandboxDocker_CgroupPrepConditionsAndSteps` |
| Fehlschlag: Zeile auf stderr, kein `write_status`, kein neuer `exit`, kein Aufruf von `unavailable` aus dem Umbau | `TestGenerateDevcontainer_SandboxDocker_CgroupPrepBestEffort` |
| keine Wirkung ohne `nestedRuntime: docker`: weder in `sandbox-init.sh`, `egress-init.sh` noch im Standardprofil | `TestGenerateDevcontainer_CgroupPrepOnlyWithNestedDocker` |
| zweiter Lauf, veralteter Block und vorhandenes Skript ohne Marker wie bisher (NoOp, `UpdatedBlock`, Konflikt) | bestehende `TestGenerateDevcontainer_SandboxDocker`, `…_ScriptConflict`, `…_StaleBlockKeepsUserContent`, angepasst auf den neuen Wortlaut |
| Zeilen und Zählungen von Exit-Code und `write_status` unverändert; neue stderr-Zeile gepinnt | bestehender `TestGenerateDevcontainer_SandboxStartScriptsWriteStatusFiles`, Liste der erwarteten Fragmente um die Zeile erweitert |
| Grenzen kommen im verschachtelten Container an: im Wegwerf-Container mit `--privileged` aus dem erzeugten Image läuft `u-boot-docker-init`, dann `docker run --memory=64m --cpus=0.5` mit einem per `docker load` geladenen Image (kein Netz) liest `memory.max = 67108864` und `cpu.max = 50000 100000`; danach gelingt ein `docker exec` in den Wegwerf-Container (Risiko `EBUSY`); Stopp des Daemons und erneuter Lauf: Exit-Code `0` ohne die Zeile des Fehlschlags (Bedingung greift) | `TestE2E_SandboxDevcontainer_NestedDocker_CgroupLimitsReachNestedContainer` (`sandbox_devcontainer_docker_test.go`, Build-Tag `docker`) |

Der Docker-Test liest vor dem Skript `cgroup.subtree_control` der Wurzel des Wegwerf-Containers und überspringt sich, wenn `memory` dort schon steht (der Umbau liefe nicht, der Test bestünde leer); nach dem ersten Lauf verlangt er `memory` dort und wertet die Fehlschlagzeile als Fehler. Er unterscheidet den Umbau nur im Normalzustand „Domain-Wurzel mit Prozessen, `memory` nicht delegiert“; `domain threaded` stellt er nicht her, das ist nur im Devcontainer dieses Repos gemessen. Die zwei Tests `TestE2E_Resources_DefaultProfileLimitsReachTheContainer` und `TestE2E_Resources_SandboxNestedDockerLimitsReachTheContainer` (`resources_docker_test.go`) prüfen `HostConfig` des äußeren Containers; in der verschachtelten Engine des Devcontainers dieses Repos laufen sie erst nach der Neuanlage grün und belegen den Umbau nicht.

### SPEC-057 — Homebrew-Formeln im Devcontainer: Katalog-Schlüssel `brew`, Pfade, Freigabe der Taps, Konflikt mit `homebrew`, Pin, Meldungen, Hosts und Grenzen

Verfeinert [`LH-FA-DEV-013`](lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer) (Formel-Liste in der Projektkonfiguration, Tap als externe Quelle nur nach Freigabe, Formeln ungepinnt und ausgewiesen) und berührt [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog: [`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen)), [`LH-FA-CONF-002`](lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (die zwei optionalen Schlüssel: [`SPEC-040`](#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)), [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts: [`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion)), [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Installer-Skript und Tap als Quellen), [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Code `10`) und [`LH-FA-CLI-005`](lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Wirkung von `--quiet`: [`SPEC-043`](#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)). Die Zusage des Lastenhefts erweitert diese Stelle nicht; sie legt Schlüssel, Pfade und Verfahren fest. Das Messprotokoll steht in [`docs/messungen/slice-v1-devcontainer-homebrew-pakete`](../docs/messungen/slice-v1-devcontainer-homebrew-pakete/README.md).

**Aktiv** heißt: Der Eintrag `devcontainer.features.brew` hat `enabled: true` und keine `source`. Mit `enabled: false`, ohne `enabled` oder mit `source` (externes Feature, [`SPEC-023`](#spec-023--devcontainer-basis-image-und-feature-quellen)) erzeugt der Generator keinen Dockerfile-Block, keine Meldung, keinen Host und prüft die Tap-Freigabe nicht; die Liste wird dann nur auf ihre Form geprüft, und ein Abschalten per `enabled: false` lässt sie gültig stehen.

**Schlüssel und Pfade.** Der Katalog-Schlüssel `brew` ist ein Dockerfile-Block ohne OCI-Quelle, ohne Standardversion und ohne Herausgeber (Eintragstyp wie `gh`, [`SPEC-051`](#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten)); er erzeugt keinen Schlüssel unter `features` in `devcontainer.json`, `u-boot doctor` meldet für ihn keinen Befund der Prüfung `devcontainer.features.drift`, und er wirkt in beiden Profilen. `brew` und `homebrew` sind zwei Schlüssel zu zwei Wegen nach Homebrew: `homebrew` ist das OCI-Feature des Drittanbieters (`meaningful-ooo`, unverändert), `brew` der Block von `u-boot`.

| Pfad | Lesen | Schreiben | Wert |
|---|---|---|---|
| `devcontainer.features.brew.packages` | ja | ja (Liste) | Formeln, kommagetrennt; jeder Eintrag `<formel>` oder `<user>/<tap>/<formel>` |
| `devcontainer.brewTaps.allow` | ja | ja (Liste) | Taps, kommagetrennt; jeder Eintrag `<user>/<tap>` |

Beide Listen sind optional; fehlen sie, sind sie leer. Die Schreibung hängt Einträge an und entdoppelt in der Reihenfolge des ersten Auftretens (Listenpfad, [`SPEC-040`](#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)); ein Eintrag wird in `u-boot.yaml` entfernt, nicht per `config set`. `devcontainer.features.<schluessel>.packages` ist für jeden Schlüsselnamen eine Pfadform; angenommen wird sie nur am Schlüssel `brew` ohne `source`.

**Form der Einträge** (Exit-Code `10`, nichts geschrieben, beim Lesen von `u-boot.yaml` und beim Schreiben mit `config set`, auch mit mehreren Paaren in einem Aufruf): Jedes Segment passt auf `[A-Za-z0-9][A-Za-z0-9@._+-]*` — kein Leerzeichen, kein Shell-Zeichen, kein führendes `-`; `@` ist Teil von Formelnamen wie `node@22` (`u-boot` bietet damit keinen Versions-Pin an). Ein Eintrag der Formel-Liste hat ein oder drei Segmente, ein Eintrag der Freigabeliste zwei; zwei oder mehr als drei Segmente, URL-Formen, Optionen von `brew install`, Casks und `Brewfile` sind keine Einträge. In der Datei wird der Eintrag Byte für Byte genommen (ein Leerzeichen am Rand ist ein Formfehler); `config set` kürzt Leerzeichen um die Kommas wie die anderen Listenpfade. Eine Formel-Liste an einem anderen Schlüssel (auch an `homebrew`) oder an `brew` mit `source` endet mit Exit-Code `10`; die Meldung nennt `brew`. `enabled: true` mit fehlender oder leerer Liste ist gültig und installiert Homebrew allein. Eine `version` an `brew` ist ein Validierungsfehler wie bei `gh` ([`SPEC-051`](#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten)): der Pin von Homebrew steht im Template und ist nicht übersteuerbar.

**Freigabe der Taps.** Ein Tap ist eine zusätzliche, vom Benutzer genannte Quelle: Seine Formel-Dateien lädt Homebrew aus dem Repository des Tap-Betreibers und führt sie beim Installieren aus. Eine Formel der Form `<user>/<tap>/<formel>` eines aktiven `brew` wird nur erzeugt, wenn `<user>/<tap>` Byte für Byte (Groß-/Kleinschreibung zählt) in `devcontainer.brewTaps.allow` steht; eine Formel ohne Tap braucht keine Freigabe. Die Liste ist unabhängig von `devcontainer.featureSources.allow` und von `devcontainer.sandbox.egress.allow`: Ein Eintrag einer Liste gibt in keiner anderen etwas frei, und das Flag `--allow-external-feature-sources` gibt keinen Tap frei. `--yes` und `--no-interactive` ersetzen die Freigabe nicht; ein Flag für Taps gibt es nicht. Der Weg ist `u-boot config set devcontainer.brewTaps.allow <user>/<tap>` oder die Datei selbst. Geprüft wird an vier Stellen mit demselben Text:

- `u-boot generate devcontainer` vor jedem Plan und Schreiben, auch unter `--dry-run`: Exit-Code `10`, im Umschlag von `--json` der Code [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) (die Abbildung der Fehlerklasse, wie bei der Freigabe einer Quelle, [`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen)), keine Datei und kein Verzeichnis, `u-boot.yaml` unverändert.
- `config set devcontainer.features.brew.packages`, wenn `brew` aktiv ist: Exit-Code `10`, die Datei bleibt byte-gleich. Bei mehreren Paaren in einem Aufruf steht das Paar der Tap-Liste vor dem Paar der Formeln.
- `config set devcontainer.features.brew.enabled true`: Der Schreibweg lässt das Aktivieren nicht zu, wenn die Liste einen nicht freigegebenen Tap trägt; das Abschalten geht immer, an `brew` wie an `homebrew`, auch in einer Datei, die `generate devcontainer` abweist.
- `u-boot doctor` als Fehler der vorhandenen Prüfung `devcontainer.features.allowlist` (keine eigene Prüfungskennung, [`SPEC-017`](#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry)); `doctor` endet bei einem Fehler-Eintrag mit Exit-Code `11`.

Die Meldung nennt je Tap den Tap und seine Formeln (nach Tap sortiert), dann die Auswege in dieser Reihenfolge — `u-boot config set devcontainer.brewTaps.allow <user>/<tap>`, den Eintrag aus der Liste in `u-boot.yaml` entfernen, das Abschalten mit `u-boot config set devcontainer.features.brew.enabled false` — und den Satz, dass die Freigabe von `devcontainer.featureSources.allow` getrennt ist, der Vergleich Byte für Byte geht und `--yes` sie nicht ersetzt.

**Konflikt mit dem OCI-Feature `homebrew`.** Ist `brew` aktiv und `homebrew` aktiviert (mit oder ohne `source`, mit oder ohne `version`), endet `u-boot generate devcontainer` mit Exit-Code `10`, ehe etwas geschrieben wird, auch unter `--dry-run`: Welche der zwei Installationen die Befehle finden, hängt an der Reihenfolge des `PATH`. Die Meldung nennt beide Schlüssel und beide Auswege (`u-boot config set devcontainer.features.brew.enabled false` oder `u-boot config set devcontainer.features.homebrew.enabled false`). Dieselbe Regel gilt beim Schreiben: `config set` weigert sich, `brew` oder `homebrew` zu aktivieren, wenn der andere Schlüssel aktiv ist (Exit-Code `10`, Datei byte-gleich; bei mehreren Paaren in einem Aufruf steht das Abschalten vor dem Aktivieren). Ein abgeschalteter Schlüssel und `brew` mit `source` (externes Feature ohne Block) lösen den Konflikt nicht aus. `u-boot doctor` meldet den Konflikt nicht.

**Dockerfile-Block.** Das erzeugte Dockerfile trägt hinter dem ersten `USER root` und hinter den Blöcken der Container-Runtime, von `ollama`, `gh`, `glab` und der drei Werkzeugketten, vor der Egress-Restriktion einen Block, der in dieser Reihenfolge

1. die Build-Abhängigkeiten von Homebrew mit `apt-get` installiert (`build-essential procps curl file git ca-certificates`), das Verzeichnis `/home/linuxbrew/.linuxbrew` dem Benutzer `vscode` übergibt und die Datei `/etc/profile.d/u-boot-brew.sh` schreibt (Eigentümer `root`, eine feste Zeile `export PATH=…` mit denselben zwei Verzeichnissen wie Schritt 2; ohne `brew shellenv` und ohne Installer-Skript: eine Login-Shell setzt den `PATH` aus `/etc/profile` neu, die Datei stellt die Verzeichnisse von Homebrew wieder her),
2. `PATH` per `ENV` um `/home/linuxbrew/.linuxbrew/bin` und `/home/linuxbrew/.linuxbrew/sbin` ergänzt,
3. als Benutzer `vscode` `Homebrew/brew` von `https://github.com/Homebrew/brew` auf den Commit `HOMEBREW_REVISION` holt (`git fetch --depth 1` des Commits, `git checkout --detach`) und den Bau abbricht, wenn `git rev-parse HEAD` ihn nicht liefert,
4. als Benutzer `vscode` die Formeln der Liste mit `brew install --formula '<eintrag>' …` in der Reihenfolge der Liste installiert, mit `HOMEBREW_NO_AUTO_UPDATE=1` und `HOMEBREW_NO_ANALYTICS=1` nur für diesen Schritt (kein `ENV`; die Zeile entfällt bei leerer Liste),
5. mit `USER root` zurückkehrt.

**Der Pin** ist das `ARG HOMEBREW_REVISION` im Template: der volle Commit-Hash von `Homebrew/brew`, aus dem Release-Tag `7.0.12` gelesen (`git ls-remote --tags https://github.com/Homebrew/brew`; der Tag ist leichtgewichtig, sein Hash der Commit). Der Hash adressiert den Inhalt; eine Prüfsumme wie bei den Release-Archiven gibt es für ein Git-Repository nicht. Er steht an einer Stelle, wird nur mit einem Release von `u-boot` gehoben und ist nicht übersteuerbar (kein `version:`); er altert ohne rotes Gate. Er deckt nur Homebrew selbst: nicht die Formeln, nicht die Metadaten, die `brew install` lädt, und nicht einen Tap. Die Vier-Stellen-Kopplung der gepinnten Werkzeugkette ([`LH-FA-DEV-011`](lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)) gilt für `brew` nicht. Der Block führt kein Installer-Skript aus, hat kein `curl … | sh`, kein `brew tap` und kein `brew trust`, keine Anmeldung und keinen Token und setzt keine `HOMEBREW_*`-Variable per `ENV`. Die Trust-Stufe von Homebrew (ab 6.0 brauchen Formeln aus Taps ein ausdrückliches Vertrauen) ist die der voll qualifizierten Formel: `brew install --formula <user>/<tap>/<formel>` vertraut nur dieser Formel, nicht dem Tap (gemessen).

**Meldungen der Befehlsausgabe** (Erzeuger: `generate devcontainer`; Level `warn`, Code [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), Subjekt `.devcontainer/devcontainer.json`; die Codes der Meldungen folgen den bestehenden Klassen der Ausweisung der Lieferkette, nicht der Anforderung, die sie verlangt). Bei einem aktiven `brew` entsteht in jedem Profil eine Meldung für den Block und, nach Tap sortiert, eine je Tap, den die Formel-Liste benutzt. Der Text ist Englisch:

- Block: Homebrew kommt aus `github.com/Homebrew/brew`, der Commit steht als `HOMEBREW_REVISION` im Dockerfile (die Meldung nennt den Hash nicht), ohne Installer-Skript; die Formeln der Liste in der Reihenfolge der Liste, ungepinnt und nicht geprüft, ein späterer Bau kann anderes installieren als ein früherer — oder, bei leerer Liste, dass keine Formel eingetragen ist.
- Tap: Der Tap (`<user>/<tap>`, freigegeben in `devcontainer.brewTaps.allow`) ist eine Quelle eigenen Rechts, seine Formeln sind Code, den Homebrew aus dem Repository des Tap lädt und ausführt, und er ist ungepinnt und ungeprüft.

Die Meldungen tragen weder `SecurityRelaxation` noch `DegradationFallback`; unter `--quiet` entfallen sie, der Block bleibt. Der Fehler der Tap-Freigabe trägt im Umschlag den Code [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) (Fehlerklasse „Konfigurationswert ungültig“), nicht [`LH-FA-DEV-013`](lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer).

**Hosts.** Die Egress-Liste führt für `brew` die drei Hosts `formulae.brew.sh`, `ghcr.io` und `pkg-containers.githubusercontent.com` (Herkunft `feature:brew`, Zeile in [`SPEC-019`](#spec-019--default-allowlist-der-egress-restriktion)). Gemessen im laufenden Container hinter der Restriktion, in einem Image mit leerer Formel-Liste (Homebrew allein): Ohne die Hosts scheitert `brew install` schon beim ersten Start von Homebrew (Portable Ruby von `ghcr.io`); mit ihnen gelingt `brew install` einer Formel von Homebrew selbst. Für ein Image mit Formel-Liste ist der erste Fehler ohne die Hosts nicht gemessen. Die Hosts, von denen eine Formel lädt, folgen der Formel und stehen nicht in der Zeile: Für die Formel `pt9912/u-boot/u-boot` scheiterte der Download am Host `release-assets.githubusercontent.com`, der nicht in der Basisliste steht; wer ihn braucht, trägt ihn in `devcontainer.sandbox.egress.allow` ein. Der Bau liegt nicht hinter der Restriktion: Sie greift beim Containerstart.

**Gemessen** ([Messprotokoll](../docs/messungen/slice-v1-devcontainer-homebrew-pakete/README.md); Linux-Host `amd64`, Docker Engine 29.9.3 verschachtelt, Basis-Image `mcr.microsoft.com/devcontainers/base:debian`, Debian 13, Homebrew `7.0.12`): Der Bau mit `ripgrep` und der Formel eines Taps gelingt im Standardprofil und im Sandbox-Profil mit eingeschalteter Egress-Restriktion (kalt 35 s und 41 s, ein Lauf; das Image wächst gegenüber dem Basis-Image um rund 292 MB, mit Homebrew allein um rund 60 MB); der Bau mit einem nicht vorhandenen Pin bricht beim Holen ab. Das Homebrew-Verzeichnis ist ein flacher Klon (`git rev-parse --is-shallow-repository` ergibt `true`); `brew --version` meldet deshalb `Homebrew >=4.3.0 (shallow or no git repository)`. `brew` steht im `PATH` von `bash`, `sh` und `zsh` als Kommando (`-c`), als Login-Shell (`-lc`), von interaktiven Shells (`bash -ic`, `bash -lic`), von `docker exec` (`bash -c`, `sh -c`, `bash -lc`, `sh -lc`) und der Login-Sitzung von `su -l vscode` (mit `bash` und mit `sh`), in beiden Profilen; Debian setzt in `/etc/profile` den `PATH` neu, die Datei `/etc/profile.d/u-boot-brew.sh` stellt die Verzeichnisse von Homebrew vorn und einmal wieder her. Entfernt man die Datei im laufenden Container, fehlt `brew` in `bash -lc` und `su -l vscode` und steht weiter in `bash -c` (über `ENV PATH`). Die Shell von `u-boot devcontainer shell` ist keine Login-Shell ([`SPEC-049`](#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)). Ein `brew install` mit veraltetem `FETCH_HEAD` (Alter über der Aktualisierungs-Automatik) verschiebt den Klon auf das neueste Release, wenn der Schalter fehlt (Probe auf dem Commit des Releases `7.0.10`: nach der Installation `7.0.12`); mit `HOMEBREW_NO_AUTO_UPDATE=1` bleibt der Commit. Zeigt die Abhängigkeit einer Tap-Formel auf einen anderen Tap, richtet Homebrew diesen Tap nicht selbst ein: Der Bau mit `shivammathur/extensions/xdebug@8.2` (Abhängigkeit `shivammathur/php/php@8.2`) und nur dem Tap `shivammathur/extensions` in `devcontainer.brewTaps.allow` bricht im Schritt `brew install` mit Exit-Code 1 ab (`This command requires the tap shivammathur/php`), kein Image; die Formel-Dateien des anderen Taps werden nicht geladen. Der Plan (`brew install --dry-run`) mit beiden Formeln in der Liste richtet beide Taps ein, in beiden Reihenfolgen; gebaut ist diese Liste nicht.

**Nicht zugesagt.** Eine andere Architektur als `amd64` (`arm64` ist nicht gemessen, die Plattform stand in der Messumgebung nicht zur Verfügung), andere Basis-Images, Engines und Hosts als die gemessenen; der Inhalt und die Version jeder Formel und jedes Tap (ungepinnt und ungeprüft); dass Homebrew eine bestimmte Formel auf der Plattform installiert (Bottles gibt es nicht für jede); Abhängigkeiten einer Tap-Formel aus anderen Taps (nur an einem Formelpaar gemessen, siehe „Gemessen“); die Hosts einer Formel und die Vollständigkeit der Host-Zeile für andere Formeln als `jq` (im laufenden Container gemessen); `brew` im `PATH` anderer Shells als `bash`, `sh` und `zsh` und in Sitzungen des Features `sshd` (nicht gemessen; `su -l` ist gemessen); Image-Größe und Bauzeit (nur gemessen, nicht begrenzt); die Kombination mit `ollama`, `gh`, `glab`, den Werkzeugketten und `nestedRuntime` über die Reihenfolge der Blöcke hinaus (gebaut ist sie nicht); Casks, `brew services`, `Brewfile` und Versions-Pins je Formel.

**Belege** (Dateien in `internal/hexagon/domain`, `internal/hexagon/application`, `cmd/uboot` und `internal/e2e`):

| Bedingung | Pinnender Test |
|---|---|
| Form der Einträge: ein oder drei Segmente der Formel-Liste, zwei der Freigabeliste, Zeichenmenge, kein führendes `-` | `TestNewBrewPackage_Accepts`, `TestNewBrewPackage_Rejects`, `TestNewBrewTap` (`internal/hexagon/domain/brewpackage_test.go`) |
| Form und Platzierung beim Lesen: Exit-Code `10`, nichts geschrieben; Liste an anderem Schlüssel, an `brew` mit `source`; `version` an `brew`; schlechter Tap-Eintrag; abgeschalteter Schlüssel hält eine gültige Liste gültig | `TestGenerateDevcontainer_Brew_InvalidEntriesExit10`, `TestGenerateDevcontainer_Brew_InactiveEntriesLeaveNoTrace`, `TestGenerateDevcontainer_Brew_FirstErrorIsDeterministic`, `TestConfig_BrewFileOutsideTheFormIsRejectedOnLoad`, `TestDoctor_UbootYaml_ErrorOnBadBrewEntries` |
| Form und Platzierung beim Schreiben (`config set`, auch mit mehreren Paaren), Listenpfad hängt an, NoOp, `config get`, `config list` | `TestConfigSet_BrewPackages_AppendsLikeTheOtherLists`, `TestConfigSet_BrewLists_FormIsCheckedOnTheWritePath`, `TestConfigSet_BrewPackages_PlacementIsCheckedOnTheWritePath`, `TestConfigGet_BrewLists_NotSet`, `TestConfigList_ShowsTheBrewLists`, `TestRun_ConfigSet_BrewPackagesFormAndPlacementExit10` (`config_brew_test.go`, `cmd/uboot/brew_test.go`) |
| Pfadform `devcontainer.features.<schluessel>.packages` und `devcontainer.brewTaps.allow` | `TestNewConfigPath_*`, `TestFixedConfigPaths_SetIsPinned`, `TestConfigPath_String_RoundTrip` (`internal/hexagon/domain/configpath_test.go`), `TestConfigPath_PackagesParsesForEveryFeatureName` |
| Block: Pin, Schritte und Reihenfolge der Zeilen (Golden, mit der Zeile der Datei `/etc/profile.d/u-boot-brew.sh`), kein Installer-Skript, kein `curl … | sh`, kein `brew tap`/`brew trust`, keine Anmeldung, kein `ENV HOMEBREW`; beide Profile und `--sandbox`; zweiter Lauf NoOp; kein Schlüssel unter `features` | `TestGenerateDevcontainer_Brew_BlockInEveryProfile`, `TestGenerateDevcontainer_Brew_ListForms` (`generate_brew_test.go`) |
| Reihenfolge zu den übrigen Blöcken, `USER root` nach dem Block, Volumes unverändert | `TestGenerateDevcontainer_Brew_OrderAgainstTheOtherBlocks` |
| Standardprofil: `devcontainer.json` byte-gleich zum Lauf ohne den Schlüssel; Aktivieren und Abschalten später ersetzt den Block und lässt Zeilen des Nutzers stehen | `TestGenerateDevcontainer_Brew_DefaultProfileJSONUnchanged`, `TestGenerateDevcontainer_Brew_ActivateAndDeactivateLater` |
| Katalog: `brew` ist ein Dockerfile-Block ohne Quelle, Version und Herausgeber, mit dem Pin `HOMEBREW_REVISION` im Template; `doctor` meldet für ihn keinen Befund der Prüfungen `devcontainer.features.allowlist` und `devcontainer.features.drift` | `TestFeatureCatalogue_KeysCoverSpecExamples`, `TestFeatureCatalogue_EntriesAreValid` (`devcontainer_features_test.go`), `TestDoctor_BrewKey_NoDriftFinding` |
| Tap nur nach Freigabe, vor dem Schreiben, auch `--dry-run`: Exit-Code `10`, nichts geschrieben, Meldung nennt Tap, Formel, Auswege; Byte-Gleichheit; eigene Liste (weder `featureSources.allow` noch das Flag geben frei); inaktiver Eintrag und Formel ohne Tap ungefragt | `TestGenerateDevcontainer_Brew_UnapprovedTapIsRejectedBeforeWriting`, `TestGenerateDevcontainer_Brew_ApprovalIsByteEqualAndItsOwnList` |
| `--yes` und `--no-interactive` ersetzen die Freigabe nicht; Umschlag mit Code [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) und dem Befehl der Freigabe; das Flag hängt nichts an die Liste | `TestRun_Generate_BrewTapNeedsApprovalAndYesDoesNotSubstitute` (`cmd/uboot/brew_test.go`) |
| Freigabe beim Schreiben: `packages` bei aktivem `brew`, das Aktivieren, mehrere Paare in einem Aufruf; Abschalten immer möglich, an beiden Schlüsseln und in einer Datei, die `generate devcontainer` abweist | `TestConfigSet_BrewPackages_TapNeedsApprovalOnTheWritePath`, `TestConfigSet_BrewEnabling_AsksTheApprovalAndTheConflict`, `TestConfigSetMany_BrewTapListBeforePackages`, `TestRun_ConfigSet_BrewTapApprovalOnTheWritePath` |
| Freigabe in `doctor`: Fehler der Prüfung `devcontainer.features.allowlist` neben den anderen Befunden; kein Befund bei Freigabe, Kernformel, abgeschaltetem Schlüssel | `TestDoctor_FeaturesAllowlist_BrewTapApproval`, `TestDoctor_FeaturesAllowlist_BrewFindingNextToTheOthers` (`doctor_brew_test.go`) |
| Konflikt mit `homebrew`: beide Richtungen, mit und ohne `source`, Exit-Code `10`, nichts geschrieben; Auswege; abgeschaltet und `brew` mit `source` ohne Konflikt | `TestGenerateDevcontainer_Brew_ConflictWithTheHomebrewFeature`, `TestRun_Brew_ConflictWithHomebrewFeatureExits10`, `TestConfigSet_BrewEnabling_AsksTheApprovalAndTheConflict`, `TestConfigSetMany_BrewConflictOrder` |
| Meldungen: eine für den Block, eine je Tap, Code, Level, Kennzeichen, Wortlaut, kein Hash; `--quiet` | `TestGenerateDevcontainer_Brew_Notices`, `TestRun_Generate_BrewNotices` |
| Hosts: drei Hosts, Herkunft `feature:brew`, keine Hosts einer Formel, abgeschaltet keine | `TestGenerateDevcontainer_Brew_EgressHosts`, `TestEgressList_BrewOrigins` |
| Bau in beiden Profilen, Pin im Image, Benutzer, `rg` und Tap-Formel laufen, `brew` im `PATH` von `bash -c`, `sh -c`, `bash -lc`, `sh -lc` und `su -l vscode`; die Kontrolle ohne die Datei `/etc/profile.d/u-boot-brew.sh` findet es in der Login-Sitzung nicht | `TestE2E_Brew_BuildsInBothProfiles` (Build-Tag `docker`, `internal/e2e/brew_docker_test.go`) |
| Nicht vorhandener Pin bricht den Bau; der Schalter hält die Revision, die Kontrolle ohne Schalter verschiebt sie | `TestE2E_Brew_PinThatDoesNotExistBreaksTheBuild`, `TestE2E_Brew_UpdateSwitchKeepsTheRevision` |
| Hinter der Egress-Restriktion gelingt `brew install` einer Formel von Homebrew mit den Hosts der erzeugten Liste und scheitert ohne zwei von ihnen | `TestE2E_Brew_EgressHostsOfHomebrew` |

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
#   resources:
#     memory: 3g                   # --memory, z. B. 512m, 3g
#     cpus: '1.5'                  # --cpus, positive Dezimalzahl
```

`u-boot init --devcontainer` schreibt `devcontainer.enabled: true` in die neu angelegte `u-boot.yaml` ([`LH-AK-005`](lastenheft.md#lh-ak-005--devcontainer-flow)).

### SPEC-008 — Eigenschaften und Pflicht-Targets des Makefiles

Pflicht-Eigenschaften:

- `.DEFAULT_GOAL := help`
- `.PHONY` für alle Targets gesetzt
- `help`-Target mit Übersicht über alle verfügbaren Targets
- Variablen mit `?=`-Defaults für Overridability (`IMAGE`, `GO_VERSION`, `GOLANGCI_LINT_VERSION`, `THRESHOLD`, `BUILD_MODE`)

MVP-Pflicht-Targets (die Spalte *Zweck* nennt den Lauf je [`SPEC-031`](#spec-031--regeln-des-docker-only-workflows); die Targets `deps`, `compile`, `lint`, `test` und `coverage-gate` führen die in der Klammer genannten Befehle mit der Toolchain der Devcontainer-Umgebung aus):

| Target          | Zweck                                                           |
| --------------- | --------------------------------------------------------------- |
| `help`          | Übersicht aller Targets                                         |
| `deps`          | `go mod download` (das Devcontainer-Umfeld; die Stage `deps` bleibt ein einzelner `docker build`-Zielort)         |
| `compile`       | `scripts/bau.sh compile` (das Devcontainer-Umfeld; die Stage `compile` ruft dasselbe) |
| `lint`          | `scripts/bau.sh lint` (das Devcontainer-Umfeld; die Stage `lint` ruft dasselbe)     |
| `test`          | `scripts/bau.sh test` (das Devcontainer-Umfeld; die Stage `test` ruft dasselbe; `PKG=…`/`RUN=…` = Teillauf, kein Nachweis)     |
| `coverage`      | Alias auf `coverage-gate`                                       |
| `coverage-gate` | `scripts/bau.sh coverage` mit `COVERAGE_THRESHOLD` (das Devcontainer-Umfeld; die Stage `coverage` ruft dasselbe) |
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

**Pin-Inventar** (Pin-Politik der gepinnten Werkzeugkette, [`LH-FA-DEV-011`](lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)). Die Version der zwei Go-Werkzeugketten-Schlüssel steht an den vier Kopplungsstellen — `ARG GO_VERSION` und `ARG GOLANGCI_LINT_VERSION` in diesem Root-`Dockerfile` und im Katalog-Template, die `?=`-Variablen in `repo.mk`, die erzeugte `.devcontainer/Dockerfile` des eigenen Projekts — und zusätzlich an der `go`-Direktive der `go.mod` (fünfter Ort, [`SPEC-016`](#spec-016--go-toolchain-mindestversion-und-dockerfile-pin)). Shell-Werkzeuge des Bau- und Entwicklungs-Verlaufs (`bash`, `sed` und weitere) tragen keine Einzel-Pins: ihre Versionen kommen mit dem Basis-Image des jeweiligen Pins. GNU make ist der zweite Host-Bedarf neben Docker und trägt eine Mindest-Version am Kopf der `Makefile` (`MAKE_MIN_VERSION`, gegenwärtig `4.3`): der Lauf bricht vor jeder anderen Logik ab, wenn die laufende make-Version älter als die Mindest-Version ist (make-eigener Guard mit `$(error)`, kein Stellenvergleich); der gemessene Stand ist 4.4.1 am Host und im Devcontainer-Image, und die CI-Runner (ubuntu-24.04) tragen die make-Version `4.3` der Distribution.

**Pin-Gleichheit im Gatter (`pinsync`).** Ein `pinsync`-Wächter (d-check-Modul, Mindest-Stand v0.91.0; Gruppen-Konfiguration in `.d-check.yml`) prüft beim Gatter-Lauf die Gleichheit der Quellen-Zeilen je Pin-Gruppe, byte-genau über die Fanggruppe der Quellen-Zeile. Die Gruppe `go` nimmt die vier `GO_VERSION`-Stellen und die `go`-Direktive auf — ihre Fanggruppe liest den Major.Minor-Teil (die Direktive steht auf dem Pin-Minor, die Patch-Stelle der Toolchain bleibt in den `GO_VERSION`-Stellen), der Vergleich ist der Minor. Die Gruppe `golangci-lint` nimmt die vier `GOLANGCI_LINT_VERSION`-Stellen mit vollem Wert-Feld (mit dem `v`-Präfix), die Gruppe `devcontainer-cli` die zwei Pin-Stellen der gepinnten Dev-Container-CLI — die Abhängigkeits-Angabe in `package.json` und die Resolved-URL in `package-lock.json`; das Tag des Hilfs-Images leitet sich aus der CLI-Version und dem Hash der drei Dateien ab und ist keine eigene Text-Stelle. Das JDK ist nicht Teil einer Pin-Gruppe: seine einzige Einstellen-Stelle ist der Template-`ARG`. Der Wächter prüft Gleichheit, keine Aktualität; je abweichender, nicht treffender oder mehrdeutiger Quelle ein Befund (`pinsync-mismatch`, `pinsync-no-match`, `pinsync-ambiguous`), der den `docs-check`-Lauf rot macht. Der d-check-Stand ist an einen Digest gepinnt (`d-check.mk`, `DCHECK_DIGEST`), gegenwärtig v0.91.0, Digest aus den Release-Notes verifiziert und gegen den gezogenen Lauf gegengeprüft.

Beleg-Tabelle (Muster [`SPEC-054`](#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv); die Test-Namen sind die des Umsetzungs-Diffs, die Gate-Zeilen den Gatter-Läufen):

| Bedingung | Pinnender Test beziehungsweise Gate-Beleg |
|---|---|
| Die Pin-Werte der vier Kopplungsstellen der zwei Go-Schlüssel sind gleich (`1.27.3`, `v2.14.0`) und folgen den Template-Args | `TestToolchainPinPlacesEqual` (`internal/hexagon/application/generate_toolchain_test.go`; liest die vier Stellen gegen die Code-Konstanten) |
| Die `go`-Direktive der `go.mod` trägt den Pin-Minor (`1.27.0` neben `1.27.3`) und zieht je Pin-Hebung nach | Pin-Gruppe `go` im `docs-check`-Gatter (`pinsync`, Fanggruppe Major.Minor); Negativprobe mit `pinsync-mismatch` in der Mismatch-Gegenprobe des Laufes (Messprotokoll [`slice-harness-dcheck-pinsync`](../docs/messungen/slice-harness-dcheck-pinsync/README.md)) |
| Die vier `GOLANGCI_LINT_VERSION`-Stellen tragen denselben Wert (`v2.14.0`) | Pin-Gruppe `golangci-lint` im `docs-check`-Gatter (derselbe Lauf), neben `TestToolchainPinPlacesEqual` für die Vier-Stellen-Kopplung |
| Die zwei `devcontainer-cli`-Stellen tragen dieselbe Version (`0.89.0`) | Pin-Gruppe `devcontainer-cli` im `docs-check`-Gatter (derselbe Lauf; Mismatch-Gegenprobe im Messprotokoll) |
| Ein Lauf unter GNU make unterhalb der Mindest-Version bricht vor jeder Logik ab | make-eigener Guard im `Makefile`: `$(error)` in der `ifeq`-Klammer zum `$(shell if …)`, Negativprobe mit `make MAKE_VERSION=4.2` (Messprotokoll) |

### SPEC-007 — Mindest-Ausschlüsse der `.dockerignore`

Mindestens auszuschließen:

- `.git`, `.gitignore`, `.github`
- IDE-Verzeichnisse: `.idea`, `.vscode`
- Agent-Verzeichnisse: `.claude`, `.codex`, `.agents`
- lokale Build-Artefakte und Caches (z. B. `dist/`, `coverage*`, `*.log`)

### SPEC-016 — Go-Toolchain: Mindestversion und Dockerfile-Pin

Mindest-Toolchain: Go 1.27 oder neuer — die `go`-Direktive der `go.mod` trägt `1.27.0`, den Pin-Minor der `GO_VERSION`-Stellen (Default-Pin `ARG GO_VERSION` im Dockerfile, [`SPEC-006`](#spec-006--build-args-des-dockerfiles-und-pin-politik), Abschnitt *Pin-Inventar*: die Direktive ist der fünfte Ort der Kopplung, sie zieht je Pin-Hebung nach; die Patch-Stelle der Toolchain bleibt in den `GO_VERSION`-Stellen).

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
| Feature `brew` ([`SPEC-057`](#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen); gemessen in einem Image mit leerer Formel-Liste: ohne die drei Hosts scheitert der erste Start von Homebrew) | `formulae.brew.sh`, `ghcr.io`, `pkg-containers.githubusercontent.com` | `feature:brew` |
| Feature `claude-code` | `registry.npmjs.org`, `api.anthropic.com` | `feature:claude-code` |
| Feature `ollama` | `registry.ollama.ai`, `dd20bb891979d25aebc8bec07b2b3bbc.r2.cloudflarestorage.com` | `feature:ollama` |
| Feature `gh` | `api.github.com`, `github.com` (beide stehen auch in der Zeile „immer“) | `feature:gh` |
| Feature `glab` | `gitlab.com` | `feature:glab` |
| Feature `java-toolchain` | `github.com` (steht auch in der Zeile „immer“), `release-assets.githubusercontent.com` | `feature:java-toolchain` |
| Clone-Quelle | Host von `devcontainer.sandbox.repository` beziehungsweise `origin` ([`SPEC-046`](#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung)) | `clone` |
| Nutzer | Einträge von `devcontainer.sandbox.egress.allow` (kleingeschriebene Hostnamen mit mindestens einem Punkt, ohne Schema, Port und Wildcard) | `user` |

Die Features `sshd` und `ollama-host` ergänzen keinen Host ([`SPEC-052`](#spec-052--feature-ollama-host-im-devcontainer-eintragstyp-container-option-meldung-und-grenzen)). Die Zeile `claude-code` führt den Host der Anthropic-API und die npm-Registry (die ohnehin zum Feature `node` gehört, das `claude-code` braucht); die Hosts der Anmeldung sowie von Telemetrie, Updates, MCP und Plugins führt die Liste nicht, der Nutzer trägt sie über `devcontainer.sandbox.egress.allow` ein. Ein Feature mit `enabled: false` liefert keinen Host.

**Herkunft.** Die Anzeige führt sechs feste Namen: `base`, `runtime`, `docker-apt`, `feature:<Schlüssel>`, `clone` und `user`. Das Lastenheft nennt fünf Klassen (Basis, Container-Runtime, Feature, Clone-Quelle, Nutzerliste); `docker-apt` ist eine Verfeinerung der Klasse Container-Runtime und trägt allein den Host `download.docker.com`, nur bei `nestedRuntime: docker`. Ein Host, den mehrere Quellen liefern (etwa `ghcr.io` bei `podman` zusammen mit `homebrew`, `registry.npmjs.org` bei `node` zusammen mit `claude-code`, oder `github.com` in der Basis zusammen mit `gh`), trägt jede seiner Herkünfte einmal, in der Reihenfolge `base`, `runtime`, `docker-apt`, Features nach Schlüssel, `clone`, `user`.

**Zeile `ollama`.** Die Liste führt den Host der Manifeste und den Host der Modell-Dateien (Blob-Speicher) des Modell-Pull; für den Namen des Blob-Hosts führt sie keinen breiteren Eintrag. `ollama.com` steht nicht in der Liste, auch nicht für Cloud-Modelle; für lokale Modelle funktionieren Start, Pull und Lauf des Servers ohne ihn ([`SPEC-047`](#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen)). Wer den Namen für ein Cloud-Modell braucht, trägt ihn über `devcontainer.sandbox.egress.allow` ein; der Eintrag schließt Subdomains ein (Absatz *Namen und Subdomains*).

**Zeilen `gh` und `glab`.** Die Liste führt für `gh` die Hosts der GitHub-API und von `github.com` und für `glab` den Host `gitlab.com`, sonst nichts. Die beiden Hosts von `gh` gehören schon zur Basis: Die Anzeige führt sie mit `base` und `feature:gh`. Hosts, die auf eine Anmeldung folgen, Hosts von Telemetrie und Aktualisierung sowie die Hosts von `glab duo` und `glab orbit` führt die Liste nicht; der Nutzer trägt sie über `devcontainer.sandbox.egress.allow` ein ([`SPEC-051`](#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten)).

**Zeile `java-toolchain`.** Die Liste führt die Hosts der Veröffentlichung der Temurin-Binaries: das Release-Verzeichnis von `adoptium/temurin21-binaries` auf `github.com` — es steht auch in der Zeile „immer“, die Anzeige führt es mit `base` und `feature:java-toolchain` — und den Asset-Host `release-assets.githubusercontent.com`; die Release-Archive und die Prüfsummen-Dateien derselben Veröffentlichung stammen aus demselben Release ([`SPEC-054`](#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv), Abschnitt *Prüfsummen-Quellen*). Die JSON-Schnittstelle `api.adoptium.net` gehört nicht dazu. Die Form des Asset-Hosts pinnt der Docker-Test der Umsetzung; weicht ihre Messung ab, ist die Zeile berichtigt. Kein Host von Telemetrie, Aktualisierung oder Brew-Angebot (Cask `temurin`/`temurin@21`) gehört zur Zeile.

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
| `devcontainer.features.allowlist` | `devcontainer` Features auf Allowlist; ein Tap der Formel-Liste von `brew` ohne Eintrag in `devcontainer.brewTaps.allow` ist ein Fehler, [`LH-FA-DEV-013`](lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer), Einzelheiten [`SPEC-057`](#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen); im Sandbox-Profil kein Feature, das den Socket des Hosts einbindet (Fehler bei `docker-cli`), [`LH-FA-DEV-006`](lastenheft.md#lh-fa-dev-006--sandbox-profil), Einzelheiten [`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) |
| `devcontainer.features.drift` | `devcontainer` Features ohne Drift |
| `devcontainer.envfile` | `devcontainer`: die Umgebungsdatei `.devcontainer/devcontainer.env`, auf die `runArgs` verweist, existiert (nur Existenz), [`LH-FA-DEV-012`](lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers), Einzelheiten [`SPEC-055`](#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung) |
| `devcontainer.sandbox.runtime` | Sandbox: nested Podman (`/dev/fuse` auf dem Linux-Host, nicht im Container: Fehler bei Fehlen, Profil-Konsistenz), [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) |
| `devcontainer.sandbox.egress` | Sandbox: Egress-Restriktion konsistent konfiguriert, [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) |
| `devcontainer.sandbox.credentials` | Sandbox: keine Klartext-Zugangsdaten (Git, Agent und Forge-Werkzeug), Token-Quelle bei https-Clone, Referenz `${localEnv:<NAME>}` bei aktivem Feature mit Token-Variable (`claude-code`, `gh`, `glab`), [`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), Einzelheiten [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) |
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
- Eingebaute Features (Schlüssel → Quelle unter `ghcr.io/devcontainers/features/`, jeweils mit der Standardversion `1`): `git` → `git`, `docker-cli` → `docker-outside-of-docker`, `node` → `node`, `java` → `java`, `go` → `go`, `kubectl-helm` → `kubectl-helm-minikube`, `python` → `python`, `sshd` → `sshd` (installiert einen SSH-Server, also einen Login-Dienst; Warnung nach [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)) und `sdkman` → `java` (das Feature installiert SDKMAN; teilt sich den Render-Schlüssel mit `java`, der einmal erscheint). Zwei Schlüssel stammen von einem Drittanbieter (Quelle außerhalb von `ghcr.io/devcontainers/features/`): `homebrew` → `ghcr.io/meaningful-ooo/devcontainer-features/homebrew` mit der Standardversion `2` (Herausgeber `meaningful-ooo`) und `claude-code` → `ghcr.io/anthropics/devcontainer-features/claude-code` mit der Standardversion `1` (Herausgeber `anthropics`; das Feature installiert über npm und braucht das Feature `node`). Für beide weist die Befehlsausgabe den Herausgeber aus, solange der Eintrag ohne eigene `source` aktiviert ist ([`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`SPEC-044`](#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)); alle Katalogschlüssel sind ohne Eintrag in `devcontainer.featureSources.allow` aktivierbar. Jede andere Quelle gilt als extern und braucht eine Freigabe in `devcontainer.featureSources.allow`.
- Der Katalog führt je Eintrag zwei Eigenschaften des Features in der Standardversion: ob es den Container-Runtime-Socket des Hosts einbindet (nur `docker-cli`) und welche Container-Optionen es beim Start mitbringt, soweit der Katalog sie führt (nur `go`: `--cap-add=SYS_PTRACE`, `--security-opt=seccomp=unconfined`; bei `docker-cli` führt er keine, denn das Sandbox-Profil weist den Eintrag ab). Bei keinem anderen Schlüssel trägt der Katalog eine von beiden, auch bei `kubectl-helm`, `ollama`, `ollama-host`, `gh` und `glab` nicht; die Container-Option von `ollama-host` ist die Wirkung des Eintrags selbst und keine Eigenschaft eines Features (Punkt zu `ollama-host` weiter unten). Die Metadaten sind für `git`, `docker-cli`, `node`, `java`, `sdkman`, `go`, `kubectl-helm`, `python`, `sshd`, `homebrew` und `claude-code` gelesen; die Schlüssel `cpp` und `postgres-client` entfallen (Lastenheft 0.33.0 — ihre Quellen existieren nicht, gemessen 2026-11-03). Was das Sandbox-Profil daraus macht, steht in [`SPEC-050`](#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen).
- Der Schlüssel `ollama` gehört zum Katalog, ist aber kein OCI-Feature, sondern ein Dockerfile-Block (zweiter Eintragstyp): Er hat weder Quelle noch Standardversion, erzeugt keinen Schlüssel unter `features` in `devcontainer.json`, und `u-boot doctor` meldet für ihn keinen Befund der Prüfungen `devcontainer.features.allowlist` und `devcontainer.features.drift`, auch wenn `devcontainer.json` keinen Schlüssel dafür enthält. Mit Sandbox-Profil trägt das erzeugte Dockerfile die Installation ([`SPEC-047`](#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen)); ohne Sandbox-Profil entsteht nichts. Eine `version` ohne `source` ist bei `ollama` ein Validierungsfehler (Exit-Code `10`); mit `source` ist der Schlüssel ein externes Feature wie jeder Katalogschlüssel mit eigener Quelle, und der Dockerfile-Block entfällt.
- Die Schlüssel `gh` und `glab` gehören ebenfalls zum Katalog und sind Dockerfile-Blöcke wie `ollama`: Sie haben weder Quelle noch Standardversion noch Herausgeber, erzeugen im Standardprofil keinen Schlüssel unter `features` in `devcontainer.json`, und `u-boot doctor` meldet für sie keinen Befund der Prüfungen `devcontainer.features.allowlist` und `devcontainer.features.drift`, auch wenn `devcontainer.json` keinen Schlüssel dafür enthält. Anders als `ollama` wirken sie in beiden Profilen ([`SPEC-051`](#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten)). Eine `version` ohne `source` ist ein Validierungsfehler (Exit-Code `10`); mit `source` ist der Schlüssel ein externes Feature, und der Dockerfile-Block entfällt.
- Der Schlüssel `brew` gehört ebenfalls zum Katalog und ist ein Dockerfile-Block wie `gh`: Er hat weder Quelle noch Standardversion noch Herausgeber, erzeugt keinen Schlüssel unter `features` in `devcontainer.json`, und `u-boot doctor` meldet für ihn keinen Befund der Prüfung `devcontainer.features.drift`. Er installiert Homebrew aus einem im Template gepinnten Commit und die Formeln seiner Liste; er wirkt in beiden Profilen und schließt sich mit dem Schlüssel `homebrew` (dem OCI-Feature des Drittanbieters) aus. Pfade, Freigabe der Taps, Pin, Meldungen und Grenzen stehen in [`SPEC-057`](#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen). Eine `version` ohne `source` ist bei `brew` ein Validierungsfehler (Exit-Code `10`).
- Der Schlüssel `ollama-host` gehört ebenfalls zum Katalog und ist ein dritter Eintragstyp: weder OCI-Feature noch Dockerfile-Block, sondern ein Eintrag, der eine Container-Option trägt (die Zeile `--add-host=host.docker.internal:host-gateway`). Er hat weder Quelle noch Standardversion noch Herausgeber, erzeugt keinen Schlüssel unter `features` in `devcontainer.json`, und `u-boot doctor` meldet für ihn keinen Befund der Prüfungen `devcontainer.features.allowlist` und `devcontainer.features.drift`, auch wenn `devcontainer.json` keinen Schlüssel dafür enthält. Er wirkt in beiden Profilen und unabhängig von `ollama` ([`SPEC-052`](#spec-052--feature-ollama-host-im-devcontainer-eintragstyp-container-option-meldung-und-grenzen)). Eine `version` ohne `source` ist ein Validierungsfehler (Exit-Code `10`); mit `source` ist der Schlüssel ein externes Feature, und die Option entfällt.
- Externe Feature-Quellen werden als Quell-Angabe getragen, die die Dev-Container-CLI verarbeitet — gemessen (2026-11-03) ist die Registry-Referenz (`host/pfad`, ohne Schema-Suffix) die verarbeitbare Form; `oci://…:<tag>` und `https://…:<tag>` sind unverarbeitbar, und die Registry-Tarball-Form ist offen. `u-boot generate devcontainer` weist eine unverarbeitbare Form je Exit-Code `10` ab, ehe nichts geschrieben wird, mit Abhilfe-Meldung der verarbeitbaren Registry-Form (Lastenheft 0.33.0), und `config set …source` wendet dieselbe Prüfung an (Schreib-Pipeline von [`SPEC-040`](#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)).

Beleg-Tabelle (Muster [`SPEC-054`](#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv); die Test-Namen sind die des Umsetzungs-Diffs; Messbasis 2026-11-03, Registry- und CLI-Lesung):

| Bedingung | Pinnender Test (Umsetzungs-Diff) |
|---|---|
| Katalog-Aufzählung: die Aufzählung oben trägt die Schlüssel ohne `cpp` und `postgres-client`, je Eintrag Quelle, Standardversion, Socket-/Options-Eigenschaften | `TestFeatureCatalogue_KeysCoverSpecExamples` (aufzählter Key-Stand) und `TestFeatureCatalogue_EntriesAreValid` (je Eintrag Quelle, Version, Kurzbeschreibung, Typ) (`internal/hexagon/application/devcontainer_features_test.go`) |
| Die Registry-Referenz (host/pfad ohne URL-Schema) ist die verarbeitete Quell-Form: `u-boot generate devcontainer` rendert sie als Schlüssel `<source>:<version>` | `TestGenerateDevcontainer_SourceFormRegistryReferenceRenders` (`internal/hexagon/application/generate_sourceform_test.go`) |
| Eine unverarbeitbare Form (`oci://…:<tag>`, `https://…:<tag>`) wird je Exit-Code `10` (Klasse `ErrConfigValueInvalid`) abgewiesen, ehe nichts geschrieben wird — auch unter `--dry-run`, auch wenn sie byte-gleich in `devcontainer.featureSources.allow` steht | `TestGenerateDevcontainer_SourceFormRejectsMeasuredForms` (beide gemessenen Formen, je Profil-Modus: nichts geschrieben, `u-boot.yaml` byte-gleich, nichts unter `.devcontainer`), `TestGenerateDevcontainer_SourceFormBeforeFreigabe` (die Form-Abweisung steht vor der Freigabepflicht) und am Nutzerpfad `TestRun_Generate_UnprocessableFeatureSourceFormIsRejected` (Exit `10`, `--json` Code [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features), Abhilfe läuft zu Exit `0`) (`internal/hexagon/application/`, `cmd/uboot/main_test.go`) |
| Die Abweisungs-Meldung nennt Feature, Quelle, die Form-Aussage und die Abhilfe (die verarbeitbare Registry-Form), in dieser Reihenfolge, ohne eine Freigabe-Möglichkeit der unverarbeitbaren Form | `TestGenerateDevcontainer_SourceFormMessageParts` (je Teilsatz ein Pin, Reihenfolge geprüft, kein `featureSources.allow`-Weg) |
| Deaktivierte Einträge und Einträge ohne Quelle werden nicht geprüft | `TestGenerateDevcontainer_SourceFormIgnoresInactiveEntries` |
| `config set …source` wendet dieselbe Prüfung an: unverarbeitbare Form je Exit-Code `10`, nichts geschrieben, Abhilfe-Zeiger; die Registry-Referenz mit Liste-Eintrag wird angenommen | `TestConfigSet_FeatureSource_UnprocessableFormRejected` (beide gemessenen Formen, `u-boot.yaml` byte-identisch, Abhilfe-Satz, kein Freigabe-Hinweis statt der Form-Aussage) und `TestConfigSet_FeatureSource_RegistryReferenceAccepted` (`internal/hexagon/application/config_features_test.go`) |

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
| 2026-10-19 | Egress-Anzeige festgelegt, deckungsgleich mit Lastenheft 0.16.0 und Code (`SPEC-046`): Berechnung (dieselbe wie die Zeile `hosts=` des Startscripts), Eingabe, Aktivierungszustand, Auflösung der Clone-Quelle samt Fehlerfall, Ausgabeformen, Wirkung von `--quiet`, Exit-Codes und Abgrenzung zum laufenden Container (ein Neustart behält die alte Liste). `SPEC-043` führt `egress list` unter den Abfragebefehlen, `SPEC-001` den Wert `egress` im Enum von `command`, `SPEC-024` einen Beispielaufruf. |
| 2026-10-19 | `SPEC-019` führt die Herkunft je Quelle (sechs Namen, `docker-apt` als Verfeinerung der Klasse Container-Runtime) und die Subdomain-Aussage mit Messstand: gepinnte und gemessene Fälle getrennt, der übergeordnete Name bei geteilter IP-Adresse als nicht geprüft ausgewiesen, der offene Upstream-Resolver ohne `sudo` als gemessen; `*`-Einträge werden mit Exit-Code `10` abgewiesen; zu `CNAME`-Zielen keine Aussage. |
| 2026-10-19 | Feature `ollama` festgelegt, deckungsgleich mit Lastenheft 0.15.0 und Code (`SPEC-047`): Dockerfile-Block mit gepinntem Archiv und Prüfsumme, Volume, kein Start, Bindung an Loopback bei ungesetztem `OLLAMA_HOST`, zwei Meldungen, Grenzen. `SPEC-023` führt den zweiten Eintragstyp, `SPEC-019` die Zeile `ollama` (zwei Hosts), `SPEC-043` die beiden Meldungen unter den entfallenden, `SPEC-044` begrenzt die Sandbox-Regel auf `sshd` und `claude-code`. |
| 2026-10-19 | `SPEC-045` führt im Absatz *Grenze* nach, dass ein Prozess ohne `sudo` die Statusdatei ersetzen kann (Verzeichnis gehört `vscode`). |
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
| 2026-10-23 | Features `gh` und `glab` festgelegt, deckungsgleich mit Lastenheft 0.22.0 und Code (`SPEC-051`): Dockerfile-Block mit gepinntem Archiv und Prüfsumme in beiden Profilen, Reihenfolge der Blöcke, `version` und `source`, Meldung zur Lieferkette je Feature, Hosts, Token-Referenzen in `remoteEnv`, Grenzen. `SPEC-023` führt die beiden Schlüssel als Dockerfile-Blöcke, `SPEC-019` die Zeilen `gh` und `glab`, `SPEC-044` die Token-Variablen und die Prüfung `devcontainer.sandbox.credentials` für beide Features (die Aussage „Andere Variablen reicht er nicht durch“ nennt jetzt die drei Variablen), `SPEC-043` die Einstufung unter `--quiet`, `SPEC-017` und `SPEC-049` nennen die Token-Variablen eines Forge-Werkzeugs. Nicht durch einen Test belegt und deshalb nicht zugesagt sind unter anderem `runArgs` und Startbefehl mit den Features, der Eigentümer der Binaries und der Inhalt der Meldung über die Wirkungsgrenze der Prüfsumme. |
| 2026-10-24 | Freigabe der Quelle bei `generate devcontainer` festgelegt, deckungsgleich mit Lastenheft 0.22.0 und Code (`SPEC-050`, Absatz *Freigabe der Quelle*; Aufruf-Abschnitt [`LH-FA-DEV-003.a`](#lh-fa-dev-003a--mechanik-der-freigabe-externer-feature-quellen)): Exit-Code `10` und Code [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features) für die Quelle eines aktivierten Features außerhalb der wirksamen Freigabe, in beiden Profilen, vor der Abweisung von `docker-cli`, ohne Schreiben, Aussage der Meldung. Der Satz „Der Eintrag wird wie bisher erzeugt“ gilt für freigegebene Quellen. |
| 2026-10-24 | Feature `ollama-host` festgelegt und Cloud-Modelle abgegrenzt, deckungsgleich mit Lastenheft 0.23.2 und Code (`SPEC-052`): Eintragstyp mit Container-Option, eine Zeile in `runArgs` in beiden Profilen, Meldung samt Einstufung unter `--quiet` und Zusatz hinter der Egress-Restriktion, Hosts und Zugangsdaten, Auflösung des Namens im Container. `SPEC-047` führt den Anschluss nicht mehr unter „Nicht zugesagt“ und grenzt Cloud-Modelle ab, `SPEC-023` den dritten Eintragstyp, `SPEC-019` den Fall des Cloud-Modells und `ollama-host` ohne Host, `SPEC-043` die Einstufung der Meldung, `SPEC-040` die Container-Option bei `version` ohne `source`; `SPEC-044` berichtigt die Sandbox-Regel: Die Aussage, `runArgs` stamme allein aus `nestedRuntime` und Egress-Restriktion, und die Aussage zu `runArgs`, `mounts` und Startbefehl bei `sshd` und `claude-code` entfallen. Nicht durch einen Test belegt und deshalb nicht zugesagt sind unter anderem `runArgs`, `mounts` und Startbefehl mit `sshd` und `claude-code`, der Inhalt von `sandbox-init.sh` und `docker-init.sh` mit dem Eintrag `ollama-host` und der Ablauf mit Cloud-Modellen (Verhalten von Ollama und Claude Code). |
| 2026-11-01 | `SPEC-049`, Unterabschnitt `shell`, *Grenzen der Erkennung*: der Fall „die Sitzung schreibt die Zeile des fehlenden Benutzers selbst“ steht als Grenze mit der gemessenen Form der Engine 29.9.2 (Fehlerausgabe, Status 1); er ist eine Lesung des Codes, kein gefahrener Lauf des Launchers. Die Form der Engine 28.0.4 ist nicht neu gemessen. |
| 2026-11-01 | `SPEC-049`, Unterabschnitt `shell`, *Grenzen der Erkennung*: zugesagt ist nur die Engine 29.9.2 (Entscheidung des Projektinhabers); die Engine 28.0.4 ist nicht zugesagt. Die Erkennung der Zeile `unable to find user …` auf der Standardausgabe mit Status 126 bleibt im Code und ist als nicht zugesagte Engine-Form eingeordnet. |
| 2026-11-01 | Grenzen für Arbeitsspeicher und CPU festgelegt, deckungsgleich mit Lastenheft 0.27.0 und Code (`SPEC-053`): Schlüssel `devcontainer.resources.memory` und `.cpus`, Zeichenmenge der Werte, Exit-Code `10`, `--memory` vor `--cpus` in `runArgs` beider Profile, Beleg-Tabelle je Bedingung; `SPEC-005` (Schema) und `SPEC-040` (erlaubte Pfade) nennen die zwei Schlüssel. Nicht belegt und nicht zugesagt: Wirkung auf den verschachtelten Docker, zu kleines Limit, mehrere Instanzen. |
| 2026-11-01 | Auflösung der Git-Zugangsdaten und Identität auf dem Host festgelegt, deckungsgleich mit Lastenheft 0.28.0 und Code (`SPEC-049`): der Schalter `devcontainer.sandbox.hostGitToken` (Sandbox-Profil; [`SPEC-040`](#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)), die Quelle `gh auth token`, wenn die Host-Variable fehlt, die Quelle als Warnung genannt, nie der Wert, die fehlende Quelle ohne Abbruch, der Wert nur im Prozess-Umfeld der `docker`-Client-Prozesse, die Grenze der Metadaten, die Identität `user.name`/`user.email` per `git config --global` im Container (beide Profile, best effort) und das JSON-Feld `gitTokenSource`. Belege: Unit-Tests der Adapter und des Start-Service, Docker-Tests mit FAKE-Quellen hinter dem Build-Tag `docker`. |
| 2026-11-02 | Katalog-Eintragstyp „gepinnte Werkzeugkette“ als Zielbild zum Lastenheft 0.30.0 festgelegt (`SPEC-054`): die zwei Schlüssel `go-toolchain` und `golangci-lint` als Dockerfile-Blocks in beiden Profilen, Prüfsummen-Quelle je Schlüssel (Same-Origin am Archiv: `go.dev/dl` je Version beziehungsweise die `CHECKSUMS`-Datei des Releases von `golangci/golangci-lint`), `version:`-Übersteuerung mit Lese-Prüfung der Quelle bei der Erzeugung, die drei Abweisungen je Exit-Code `10` (Quelle nicht lesbar oder überdeckt die Version nicht, Eintrag ohne Quellen-Angabe, externes `go` plus `go-toolchain`), Cache-Werte `GOCACHE`/`GOMODCACHE` unter `containerEnv` je unterhalb von `${containerWorkspaceFolder}`, Meldungen je aktiviertem Schlüssel samt Einstufung unter `--quiet`, Pin-Politik der vier Stellen. Jede Bedingung trägt den Zeiger auf den pinnenden Test; die Beleg-Tabelle Zusage → Test entsteht mit dem Umsetzungs-Diff. Nicht belegt und nicht zugesagt: das Arbeitsverhalten in übersteuerten Versionen, das Überstehen eines Neubaus des Images, das Standardprofil und andere Hosts, Engines und Basis-Images für die Cache-Werte, `BUILD_MODE`/`GOFLAGS`/`GOTOOLCHAIN`. |
| 2026-11-02 | `SPEC-054` um den dritten Schlüssel `java-toolchain` erweitert (Zielbild zum Lastenheft 0.31.0; die drei Schlüssel im Titel): JDK der Eclipse-Temurin-Distribution der LTS-Version 21 aus der Veröffentlichung der Temurin-Binaries, Prüfsummen-Quelle die Prüfsummen-Datei derselben Veröffentlichung je gewählter Version (`api.adoptium.net` ausdrücklich nicht Quelle), die drei Abweisungen des Typs angewandt auf den dritten Schlüssel samt Konflikt externes `java` plus `java-toolchain` (Exit-Code `10`), der Schlüssel trägt keinen Wert unter `containerEnv` (`JAVA_HOME`/`PATH` im Block), die Vier-Stellen-Kopplung bleibt bei den zwei Go-Schlüsseln, Meldung des Schlüssels mit Vertriebs-Form; `SPEC-019` die Hosts-Zeile `feature:java-toolchain`, `SPEC-040` der dritte Schlüssel-Namen in den zwei Übersteuerungs-Zeilen. Die Java-Zeilen der Beleg-Tabelle tragen die Bedingung als Beleg-Auftrag; ihre Test-Namen entstehen mit dem Umsetzungs-Diff. Nicht belegt und nicht zugesagt für den neuen Schlüssel: die ARG-Form des JDK-Schlüssels (Umsetzungs-Messung per Docker-Test), das Arbeitsverhalten des JDK in übersteuerten Versionen, das Überstehen eines Neubaus des Images. |
| 2026-11-02 | Umsetzung des dritten Schlüssels: die Test-Namen der Java-Zeilen der Beleg-Tabelle sind geliefert (Katalog-Tabellentest samt Quellen-Angabe, Render-Tests mit Fixture-Quellen je Fehlerfall, Golden Case der erzeugten `devcontainer.json` ohne `containerEnv`-Werte, Docker-Tests des Baus samt Negativprobe der Prüfsumme und des Baus gegen die für `21.0.9+10` gelesene Prüfsumme). Die ARG-Form des JDK-Schlüssels ist die bare Version mit `+` als Trenner des Builds (`21.0.13+9`), beide Schreibweisen (`jdk-`-Präfix oder bare) normalisieren auf sie, gepinnt durch die Offline- und Docker-Tests; der Default-Pin des Templates trägt die Revision `21.0.13+9` der Veröffentlichung vom 2026-10-22 (Prüfsummen und Archiv-Form gelesen, Protokoll unter `docs/messungen/`). |
| 2026-11-03 | Konfig-Schlüssel `devcontainer.containerEnv` als Zielbild zum Lastenheft 0.32.0 festgelegt (`SPEC-040`, Unterabschnitt *Projekt-Umgebungswerte*): der Pfad trägt eine Map von Namen zu String-Werten, die Erzeugung trägt die Werte in das verwaltete `containerEnv`-Objekt der erzeugten `.devcontainer/devcontainer.json` neben den Generator-Werten der Schlüssel (ohne Duplikate, idempotent), Validation von Name und Wert mit Exit-Code `10` (auch `config set`), die Kollision eines Projekt-Wertes mit `GOCACHE`/`GOMODCACHE` eines aktiven Toolchain-Schlüssels beendet die Erzeugung mit Exit-Code `10` (ohne den Schlüssel erlaubt), die Werte stehen nur in der `containerEnv`-Ebene; die Ausweisung der Erzeugung nennt jede erzeugte Projekt-Wert-Zeile. Die Beleg-Tabelle trägt die Bedingungen als Aufträge; ihre Test-Namen entstehen mit dem Umsetzungs-Diff, die Ebene-Messung (Werte im Container-Prozess; nicht im Dockerfile-Block) ist Pflicht der Umsetzung. Nicht zugesagt: `${localEnv:…}`-Interpolation in Werten, eine `remoteEnv`-Konfiguration, Dockerfile-`ENV`-Zusätze aus dem Pfad. |
| 2026-11-03 | Umsetzung des Schlüssels (Lastenheft 0.32.0): die Test-Namen der ContainerEnv-Zeilen der Beleg-Tabelle sind geliefert, und der Schlüssel `devcontainer.containerEnv` trägt die Map in das verwaltete `containerEnv`-Objekt (Sortierung, JSON-Escaping), die Basis-Pfad-Auslegung (Exit `10`, auch bei `config get`) ist über Wertetabelle und Tests gepinnt; die Cache-Goldens bleiben byte-identisch. Eigenes Devcontainer umgestellt (die drei Rest-Werte unter `devcontainer.containerEnv`, ENV-Rest entfallen). |
| 2026-11-03 | Lastenheft 0.33.0 als Zielbild im Umsetzungs-Muster (`SPEC-023`/`SPEC-040`): die zwei Katalog-Einträge `cpp`/`postgres-client` entfallen (Quelle-Fehlen gemessen, Kontrollproben `go`/`java` je 200; Registrier-Zugriff mit Token-Verweigerung statt `NAME_UNKNOWN` gelesen), die Form-Tabelle der externen Quellen trägt das Gemessene (Registry-Referenz verarbeitet; `oci://…:<tag>`, `https://…:<tag>` unverarbeitbar; Tarball-Form offen), die Abweisung unverarbeitbarer Formen trägt der Schreib-Pipeline-Schritt 6 und die Erzeugung (Beleg-Aufträge; Test-Namen mit dem Umsetzungs-Diff). Nicht belegt und nicht zugesagt: der Tarball-Weg der CLI und seine Verarbeitbarkeit. |
| 2026-11-03 | Umsetzung von Lastenheft 0.33.0: die Test-Namen der Beleg-Zeilen von `SPEC-023` und `SPEC-040` sind geliefert — die Katalogtabelle des Codes trägt die zwei Einträge nicht mehr, die Erzeugung rendert die Registry-Referenz (host/pfad ohne URL-Schema) als Schlüssel `<source>:<version>` und weist die unverarbeitbaren Formen je Exit-Code `10` ab (vor der Freigabepflicht, auch unter `--dry-run`, auch bei byte-gleichem Liste-Eintrag), `config set …source` wendet dieselbe Prüfung in Schritt 6 an, beide Abhilfe-Zeiger sind je am Nutzerpfad ausgeführt. Die Stage-1-Grammatik der Feature-Quellen erkennt die Registry-Referenz daneben an; die Freigabeliste selbst bleibt form-permissiv. Nicht belegt und nicht zugesagt: der Tarball-Weg der CLI und seine Verarbeitbarkeit. Die Abweisung an den zwei Stellen trägt die Regel, dass nur die Registry-Referenz die verarbeitete Form ist; die zwei gemessenen Formen sind ihre gepinnten Fälle, URL-Schema-Formen außerhalb der zwei tragen keine eigene Messaussage. |
| 2026-11-03 | Pin-Inventar als Zielbild und Umsetzung zum Lastenheft 0.34.0 (`SPEC-006`, Abschnitte *Pin-Inventar* und *Pin-Gleichheit im Gatter*, `SPEC-016`): die Direktive der `go.mod` als fünfter Ort der Pin-Kopplung auf den Pin-Minor gehoben (`1.27.0` neben dem Pin `1.27.3`, `go mod tidy`-Gegenprobe am Docker-Stage mit byte-stabilem Nachlauf — der Lauf normalisiert das Require-Layout), Shell-Werkzeuge ohne Einzel-Pins (Basis-Image-Gruppe), GNU make mit Mindest-Version `4.3` am Kopf der Makefile und ab-brechendem Guard (gemessener Stand 4.4.1 am Host und im Devcontainer-Image, `4.3` auf den CI-Runnern Ubuntu 24.04), d-check auf v0.91.0 am verifizierten Digest gehoben und das Modul `pinsync` ins Gatter genommen; die Pin-Gruppen `go`, `golangci-lint` und `devcontainer-cli` tragen je Quelle einen RE2-Ausdruck mit einer Fanggruppe, Beleg-Tabelle mit den Gate-Läufen und der Mismatch-Gegenprobe (Protokoll unter `docs/messungen/`). |
| 2026-11-04 | Umgebungsdatei des Devcontainers als Zielbild zum Lastenheft 0.35.0 festgelegt (`SPEC-055`): Pfade, die Zeile `--env-file=${localWorkspaceFolder}/.devcontainer/devcontainer.env` als erster Eintrag von `runArgs` beider Profile samt Reihenfolge gegen `resources`, Egress, `nestedRuntime` und `ollama-host`, Anlage nur wenn fehlend (nie gelesen, geändert, gesichert), Beispieldatei als verwalteter Block, `.gitignore`-Zeile nur über `init`, Prüfung `devcontainer.envfile` (nur Existenz, `warn`, ungekennzeichnet unter `--quiet`); `SPEC-052`, `SPEC-053`, `SPEC-047`, `SPEC-044` und die Zeilen zu `nestedRuntime` ([`LH-FA-DEV-007.b`](#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman)) zählen die neue Zeile mit. Die Beleg-Tabelle trägt die Bedingungen als Aufträge; ihre Test-Namen entstehen mit dem Umsetzungs-Diff. Der Start über `devcontainer up` gehört zur Zusage-Menge von [`LH-FA-DEV-010`](lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), soweit der Docker-Test ihn belegt; nicht zugesagt: der Start über den Editor (nicht automatisierbar). |
| 2026-11-04 | cgroup-Umbau im Startscript `docker-init.sh` festgelegt (`SPEC-056`): Ort vor `dockerd` im `postStartCommand`-Skript, Bedingung (cgroup v2, `memory` noch nicht in `subtree_control`), drei Schritte, Fehlschlag best-effort ohne Statusänderung, Beleg-Tabelle; `SPEC-053` verweist darauf. Lastenheft unverändert. |
| 2026-11-05 | Homebrew-Formeln im Devcontainer festgelegt, deckungsgleich mit Lastenheft 0.36.1 und Code (`SPEC-057`): Katalog-Schlüssel `brew` als Dockerfile-Block mit dem Pin `HOMEBREW_REVISION` (Commit von `Homebrew/brew` im Template), die Pfade `devcontainer.features.brew.packages` und `devcontainer.brewTaps.allow` samt Form der Einträge, Freigabe der Taps an vier Stellen (Erzeugen, Schreiben, Aktivieren, `doctor`), Konflikt mit dem OCI-Feature `homebrew`, Meldungen (Code [`LH-NFA-SEC-004`](lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte); der Fehler der Freigabe trägt im Umschlag [`LH-FA-DEV-003`](lastenheft.md#lh-fa-dev-003--devcontainer-features)), Hosts-Zeile `brew` in `SPEC-019`, Grenzen der Zusage mit den Messwerten (nur `amd64`; `brew` im `PATH` von Login-Shells durch die Datei `/etc/profile.d/u-boot-brew.sh`), Beleg-Tabelle. |
