# Messung: CI-Fähigkeit der Devcontainer-Basis (Messung M1)

| Dokument    | Messprotokoll |
| ----------- | ------------- |
| Auslösender Vorgang | [`slice-v1-bau-nur-devcontainer`](../../plan/planning/done/slice-v1-bau-nur-devcontainer.md), Plan-Risiko 1 — die Messung, die vor dem Vertrags-Wortlaut (Umformulierung von [`SPEC-031`](../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows), Folge-ADR zu [`ADR-0037`](../../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md)) die CI-Fähigkeit der Devcontainer-Umgebung belegt |
| Gelesen am  | 2026-11-03 |
| Skripte     | `tools/messung/devcontainer-ci/` — Fahrer [`messen.sh`](../../../tools/messung/devcontainer-ci/messen.sh) (Modi `bau`, `dind`, `gates`, `test-docker`, `analogie`); im Container [`im-container.sh`](../../../tools/messung/devcontainer-ci/im-container.sh) (M2), [`im-gates.sh`](../../../tools/messung/devcontainer-ci/im-gates.sh) (M3), [`im-test-docker.sh`](../../../tools/messung/devcontainer-ci/im-test-docker.sh) (M4), [`im-test-docker-ergaenzung.sh`](../../../tools/messung/devcontainer-ci/im-test-docker-ergaenzung.sh) (M4-Ergänzung); Analogie [`analogie-ebene-a.sh`](../../../tools/messung/devcontainer-ci/analogie-ebene-a.sh) mit den Ebene-2-Skripten (M5) |
| Rohdaten    | [`roh-bau.txt`](roh-bau.txt), [`roh-dind.txt`](roh-dind.txt), [`roh-gates.txt`](roh-gates.txt), [`roh-test-docker.txt`](roh-test-docker.txt), [`roh-test-docker-ergaenzung.txt`](roh-test-docker-ergaenzung.txt), [`roh-analogie.txt`](roh-analogie.txt) |

## Umgebung

- 2026-11-03, Linux-Host `x86_64` (4 CPUs, 7 GB RAM, 23 GB frei vor der Messung,
  11 GB danach — das Messbild selbst belegt 2,47 GB), Docker Engine am Host
  29.9.2 (containerd v2.4.2, runc 1.5.2). Kein Host-Go; alle Läufe über Docker.
- Arbeitsbaum der Messung: der Stand nach Commit `d8104049`; die Messdateien
  waren zum Laufzeitpunkt noch uncommittet (der Arbeitsbaum wurde je Lauf
  per Bind-Mount eingehängt).
- Die Werte von `containerEnv` der erzeugten
  [devcontainer.json](../../../.devcontainer/devcontainer.json) — `BUILD_MODE=local`,
  `GOTOOLCHAIN=local`, `GOFLAGS=-mod=readonly -buildvcs=false`,
  `GOCACHE`/`GOMODCACHE` unter dem Workspace — hat der bare
  `docker run`-Aufruf je Messlauf per Optionen nachgebaut; die Bare-Form
  übernimmt die `containerEnv` nicht — die devcontainer-CLI tut es.
- Die Messung liest und protokolliert; sie baut nichts zur Behebung.

## M1 — Das Devcontainer-Image bauen ([`roh-bau.txt`](roh-bau.txt))

- **Kontext-Befund:** die notierte Form `docker build -f .devcontainer/Dockerfile .`
  (Kontext = Repo-Wurzel) **scheitert** — Exit 1 an `COPY docker-init.sh`
  („/docker-init.sh: not found“), denn die Datei liegt neben dem Dockerfile.
  Die arbeitsfähige Form ist `docker build -f .devcontainer/Dockerfile .devcontainer`
  (Kontext = `.devcontainer`-Verzeichnis); mit ihr ist der Bau Exit 0.
- Image `uboot-devc-m1`: **Exit 0, wall 40 s** (alle Ebenen kamen aus dem
  lokalen Bau-Cache der bisherigen Devcontainer-Bauten; ein kalter Netzbau
  ist hier **nicht** gemessen). Image-Id (= manifest list, Buildx mit
  Attestation-Manifest): `sha256:fc2aa5e219c8e2b7955e478388f461f446bf6acf5431fe9bd2c5011f8560d0b8`,
  Konfigurations-Digest `sha256:8d346eb514f7bf04f05f763c83f657efe92b4326f6d886aca61eef7ce72242a6`,
  2.474.989.058 Byte, `amd64`, erzeugt 2026-11-03T14:18:42Z.
- **Was das Image trägt** (gelesen, Benutzer `vscode` (uid 1000),
  `HOME=/home/vscode`):
  - `make` 4.4.1, `bash` 5.2.37, `git` 2.56.0,
  - `go` 1.27.3, `golangci-lint` 2.14.0 (gebaut mit go1.27.0),
  - `docker` 29.9.2 — Client **und** `dockerd` im Image (Paket `docker-ce`/
    `docker-ce-cli`/`containerd.io`/`buildx`/`compose-plugin`/`iptables`),
    der Daemon startet erst zur Container-Laufzeit (Startscript
    `u-boot-docker-init`, Policy `warn` → fallback),
  - `ollama` 0.40.6 (Client-Archiv, kein Server im Image), `gh` 2.103.0,
  - **`node` fehlt im Bild** — es ist ein Feature-Eintrag der
    `devcontainer.json` (`node:2.1.1`), die die devcontainer-CLI erst beim
    echten Start zusätzlich schreibt; der bare `docker build` trägt ihn nicht.
- Ohne gestarteten Daemon ist `docker info` im Image ein Fehler (`dial unix
  /var/run/docker.sock: connect: no such file or directory`).

## M2 — Der verschachtelte Daemon ([`roh-dind.txt`](roh-dind.txt))

Privilegiert gestarteter Container des Messbilds, Arbeitsbaum per Bind-Mount
unter `/workspaces/u-boot`, `/var/lib/docker` als frisches Volume:

- Startscript `u-boot-docker-init`: Exit 0; `docker version` meldet danach im
  laufenden Container den **Server-Anteil: Engine 29.9.2** (containerd
  v2.4.2, runc 1.5.2), `overlayfs`, `cgroupfs`, cgroup v2.
- Test-Container im Container: `docker pull busybox` Exit 0 (2 s), `docker
  run --rm busybox true` Exit 0 (1 s).
- `make modus-pruefung` im Container: **OK** (alle 8 Target-Regeln beider
  Modi), Exit 0 (0 s) — die Prüfung liest die `BUILD_MODE=local`-Umgebung
  nach wie vor als Devcontainer-Umgebung.

## M3 — `make gates` im Container ([`roh-gates.txt`](roh-gates.txt))

Derselbe Container-Aufbau wie M2; die Go-Caches sitzen (wie im echten
Devcontainer je `containerEnv`) unter dem Workspace — bei der Bind-Mount-Form
im Baum `./.cache`, vor Lauf 1 geräumt (kalt):

| Lauf | Exit | Wanduhr |
|---|---|---|
| `make gates`, Lauf 1 (`.cache` vor dem Lauf geräumt) | 0 | 62,0 s |
| `make gates`, Lauf 2 (warm am selben Baum) | 0 | 33,0 s |
| `docs-check` (warm, einzeln) | 0 | 6,8 s |
| `lint` / `test` / `coverage-gate` (warm, einzeln) | 0 / 0 / 0 | 5,4 s / 0,4 s / 0,7 s |

- **Coverage 92,20 %** (Schwelle 90 %) in beiden Läufen.
- **`docs-check` läuft im verschachtelten Daemon:** das d-check-Image
  (`ghcr.io/pt9912/d-check@sha256:5d2f650d102c1d8224bb26930ec96d237506c07612fb8e0dbcdc42291702c57f`)
  baut/pullt der Lauf im Nested-Daemon, 1086 Dateien geprüft, 0 Befund(e),
  inklusive pinsync (3 Pins, 11 Quellen gelesen).

## M4 — `make test-docker` im Container ([`roh-test-docker.txt`](roh-test-docker.txt), Ergänzung [`roh-test-docker-ergaenzung.txt`](roh-test-docker-ergaenzung.txt))

**Exit 2, wall 869,0 s**; Zählung der Testfunktionen der obersten Ebene:
**1758 PASS, 12 FAIL, 1 SKIP** (Skip: `TestDevcontainerUp_Docker_S3_ObserveUpdateRemoteUserUID`).
Paketergebnisse: `internal/adapter/driven/docker` FAIL (219,9 s), `internal/e2e` FAIL (864,9 s).

Die 12 Fehlschläge, je mit Meldungsklasse (die sechs ohne Meldung im
Hauptauschnitt im gezielten Ergänzungs-Lauf nachgelesen — ein Filter-Lauf
über alle sechs, Exit 1, dieselben Fehlschläge, und ein Zweitlauf nur für die
ForgeCLI-Meldung):

| Klasse (Meldung) | Tests |
|---|---|
| „restricted host: init exit 0, want 11 … kernel.apparmor_restrict_unprivileged_userns=1; nested Podman cannot run here“ — der Startscript fällt im Nested-Daemon auf Exit 0 zurück, der Test erwartet den abgewiesenen Start (Exit 11) | `TestE2E_SandboxDevcontainer_NestedPodman`, `…_NestedPodman_Build`, `…_NestedPodman_BuildEgress`, `TestE2E_SandboxDevcontainer_StatusFiles_VFS` |
| „build output lacks the checksum failure:“ — der Bau bricht im Nested-Daemon ab, aber nicht mit der erwarteten Prüfsummen-Meldung (Ursache nicht geklärt) | `TestE2E_ForgeCLI_PinMismatchBreaksBuild` (gh und glab), `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild`, `TestE2E_Toolchain_PinMismatchBreaksBuild` (go, golangci-lint, jdk) |
| „the build failed for another reason than the integrity check: exit status 1: Dockerfile:17“ | `TestLauncher_Docker_S5_GuardFalsifiedIntegrityBreaksTheBuild` |
| „the value was on a command line: [docker exec -i -u vscode -e HOSTNAME=… -e HOME=…]“ — der S6-Wert-Guard sieht den Wert auf einer exec-Kommandozeile | `TestLauncher_Docker_S6_ValueIsNeverInImageLabelFileOrCommandLine` |
| „Up: devcontainer up failed (exit status 1: … #11 exporting layers …)“ — ein `devcontainer up` im Nested-Daemon bricht im Image-Bau/Export ab | `TestE2E_Resources_DefaultProfileLimitsReachTheContainer`, `TestE2E_Resources_SandboxNestedDockerLimitsReachTheContainer` |
| „build output does not show the failed step's container removed, so the leftover check proves nothing (builder is not the classic one, or --force-rm is miss…)“ | `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild` |

Zum Vergleich (Messung von 2026-10-31, Protokoll
[`slice-v1-devcontainer-bau-lokal`](../slice-v1-devcontainer-bau-lokal/README.md),
Abschnitt G/Gegenprobe): `make test-docker` im **Docker-Modus am Host** —
Exit 0, 1619 s, 1627 PASS, **0 FAIL**, 4 SKIP; im Nested-Daemon damals 9 FAIL,
heute 12 (drei Tests sind seitdem neu: die zwei Ressourcen-Tests und der
Toolchain-Pin-Test; die neun damaligen Namen sind je wieder in dieser Liste).

Die Klassen-Diagnose ist die des Meldungstexts; die Ursachen der Klassen
„Prüfsummen-Meldung fehlt“ und „build failed for another reason“ sind
**nicht geklärt** — die Messung liest, sie behebt nicht.

## M5 — CI-Runner-Analogie ([`roh-analogie.txt`](roh-analogie.txt))

Die Kette des CI-Falls (Job-Container über `docker.sock` statt privilegiertem
DIND auf Ebene 1), lokal nachgemessen: Container **A** (das Messbild,
**unprivilegiert** gestartet — `CapEff=0000000000000000`, nur der
Host-Docker-Socket eingehängt; die Analogie zum CI-Job-Container):

- Über den Socket meldet sich A an der Ebene-1-Engine (29.9.2). `docker run`
  aus A erzeugt den Ebene-2-Container — **die `--privileged`-Flag honorert
  die Engine, die den Container erzeugt, nicht die Privilegien von A**:
  - **(a) mit `--privileged`:** Startscript Exit 0, Server-Anteil 29.9.2
    (`overlayfs`) auf Ebene 2, `busybox` pull Exit 0 (2 s) und
    `docker run --rm busybox true` auf **Ebene 3** Exit 0. Die DIND-Kette
    Host-Docker → Job-Container → privilegiertes Devcontainer-Container →
    verschachtelter Daemon → Test-Container läuft.
  - **(b) ohne `--privileged`:** Startscript Exit 0 (Policy `warn`),
    Status-Datei `result=fallback`, kein Daemon („Cannot connect“);
    `docker run` auf Ebene 3 scheitert. Unprivilegiert trägt kein DIND —
    die `--privileged`-Flag (das `runArgs`-Bild des erzeugten Containers)
    ist die Bedingung.

## Was daraus folgt

- **M1–M3, M5: die CI-Migration ist beobachtbar belegt** soweit die Gatter
  der `make gates`-Kette reichen: das Image baut aus dem `.devcontainer`-Kontext,
  der verschachtelte Daemon startet, `make gates` läuft im Container grün
  (62,0 s kalt, 33,0 s warm, Coverage 92,20 %), und seine Docker-Benutzung
  (`docs-check` mit dem d-check-Image im Nested-Daemon) läuft grün. Die
  Docker-Verhältnisse des CI-Falls (Job-Container am Socket, darin das
  privilegierte Devcontainer-Container mit DIND) sind von der Ebene-1-Engine
  aus analog reproduziert.
- **M4 blockiert den ganzen Zug:** `test-docker` im verschachtelten Daemon ist
  **nicht gleichwertig** (12 FAIL gegen 0 im Docker-Modus am Host) — das
  Bestätigt die heutige Zusage-Lage von [`SPEC-031`](../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows)
  (die Tests sind im Nested-Daemon „nicht als Nachweis zugesagt“) und
  begrenzt die Migration: die CI-Migration kann die Gatter-Kette
  (`make gates`) in den Devcontainer-Container verlegen, **`test-docker` trägt
  sie nicht** — er bleibt am Runner-Docker (Docker-Modus, Socket-Run wie im
  heutigen `repo.mk`-Target) oder braucht Folge-Arbeit an den
  Test-Umgebungsannahmen (sechs Klassen oben), die nicht Bestandteil dieser
  Messung ist. Der Vertrags-Wortlaut kann diese Messbasis je Punkt anführen.

## Lesehilfen

- In [`roh-analogie.txt`](roh-analogie.txt) ist die Zeile „(info gelang
  unerwartet)“ ein Artefakt der Probe-Skript-Logik (`||`/`&&`-Vorrang), kein
  Ergebnis — der Meldungstext „Cannot connect“ der Zeile davor ist das
  Ergebnis von Variante (b).
- Die Go-Caches von M3 sind als frische **named Volumes** an den
  Workspace-Unterpfaden nicht gemessen: frisch angelegte Volumes an Pfaden,
  die das Image nicht vorhält, sind root-owned (`mkdir .cache/go-build/00:
  permission denied` im Zwischenlauf, ersetzt); die gemessene Form ist die
  vom `containerEnv` vorgeschriebene, Cache unter dem Workspace-Bind-Mount.
- Ein Lauf je Zelle, ein Host. Die Wanduhr des Image-Baus (40 s) ist
  **cache-warm**; der Bau hat in dieser Messung nichts aus dem Netz gezogen.
- Die M4-Ergänzung ist ein Filter-Lauf über die sechs Fehlschläge (Exit 1)
  plus der ForgeCLI-Zweitlauf; die Zählung des ganzen Zugs (1758/12/1) kommt
  aus dem Hauptlauf, die Zählung der Ergänzungsläufe beschreibt nur ihren
  Umfang.
- Die Lauf-Volumes (`uboot-m1-docker`, `uboot-m1-gocache`, `uboot-m1-mod`,
  `uboot-m1-docker-analogie-priv/nop`, Probe-Volume und -Image) werden nach
  dem Commit entfernt (Platz am Host: 11 GB frei nach der Messung); die
  Rohdaten tragen die Belege.

## Grenzen

- Ein Host (`amd64`), je Messpunkt ein Lauf; keine Streuung. Die CI-Runner
  (`ubuntu-24.04` bis `ubuntu-latest`) sind nicht selbst gemessen — M5 ist
  die Analogie über `docker.sock`, nicht ein GitHub-Lauf; die Migration
  selbst braucht den echten CI-Lauf als Beleg (DoD des Slices).
- Das Messbild ist der bare `docker build`-Stand ohne die Features der
  `devcontainer.json` (`node`, `claude-code`, `homebrew`); die `containerEnv`
  kommt als `-e`-Optionen statt durch die devcontainer-CLI. Der echte
  Devcontainer-Start mit Features ist in der bau-lokal-Messung (Abschnitt H)
  gemessen; hier geht es um die CI-Form (bare `docker run` auf dem Image).
- Der Image-Bau lief cache-warm; ein Netzbau von Null ist hier nicht belegt —
  die Registry-Züge der Ebene (Basis-Image, `docker-ce`-Pakete, die
  Archive von `ollama`, `gh`, `go`, `golangci-lint`) tragen die Host-Lage
  der bisherigen Bauten, und die CI-Migration braucht für sie einen
  Netzbau-Beleg der Runnerversion.
- Die Ursachen der Klassen „Prüfsummen-Meldung fehlt“ (4 Fehlschläge) und
  „build failed for another reason“ (1) sind aus den Meldungstexten nicht
  bestimmt; die zugehörigen Codes stecken in den Tests
  (`forge_cli_docker_test.go:107`, `sandbox_devcontainer_docker_test.go:1362`,
  `toolchain_docker_test.go:153`, `devcontainer_launcher_docker_test.go:825`).
- `node`/`claude`/`brew` im bare-Bild ungetestet (Feature-Stand — siehe M1);
  ein CI-Job, der sie braucht, ist hier nicht abgedeckt.