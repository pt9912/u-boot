# Messprotokoll: Forge-CLIs `gh` und `glab` im Devcontainer — Quelle, Lieferkette, Hosts, Token-Umgang

Rohdaten und Ablauf der Messung des Slice
[`slice-v1-forge-cli-quellen-messung`](../../plan/planning/in-progress/slice-v1-forge-cli-quellen-messung.md).
Skripte und Hilfsdateien liegen unter `tools/messung/forge-cli/` (Ablage: [`docs/messungen/README.md`](../README.md)).
Die Entscheidungsvorlage zur Quelle ist [`ADR-0035`](../../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) (`Proposed`), der Wortlaut-Entwurf für den
Lastenheft-CR steht in [`wortlaut-entwurf.md`](wortlaut-entwurf.md).

**Messstand** (`rohdaten-01.txt`): Linux `amd64`, Kernel 6.8.0-147-generic, 13th Gen Intel Core i9-13900H (20 Threads), 31,1 GiB RAM;
Docker 29.9.2 (Client und Server, lokale Engine, `overlayfs`); Dev-Container-CLI 0.89.0 (Hilfs-Image, `Dockerfile.cli`); Basis-Image
`mcr.microsoft.com/devcontainers/base@sha256:6ee2b4097c75d33e9a94d941577a3286e5d629434fba384688b06daeeaaeac99` (Debian 13, `amd64`);
`u-boot` aus dem Arbeitsverzeichnis (`make build`, `u-boot version 0.14.0-dev`, HEAD `2e986785`, Arbeitsverzeichnis ohne Änderung an `internal/`, `cmd/`,
Vorlagen; `rohdaten-04.txt`). Gemessen am 2026-10-23 (Zeiten in den Rohdaten: UTC), Netzzugriff auf `ghcr.io`, `github.com`, `api.github.com`,
`cli.github.com`, `gitlab.com`, `containers.dev` und die Debian-Spiegel; **ohne jede Anmeldung** (siehe *Token-Sicherheit*).

## Reihenfolge und Aufruf

`WORK=<messordner ausserhalb des Repos> tools/messung/forge-cli/00-alles.sh` führt 01 bis 07 aus; `08-aufraeumen.sh` läuft getrennt. Jedes Skript
schreibt seine Rohdaten (geschwärzt, siehe unten) selbst in dieses Verzeichnis (Vorgabe `OUT`). Voraussetzungen auf dem Host: Docker und GNU `make`;
`jq`, `curl`, `sed`, `awk`, `gpg` laufen nur in Containern. Der Messordner `WORK` (Downloads, erzeugte Projekte) liegt außerhalb des Repos.

| Schritt | Skript | Rohdaten | Messfrage |
|---|---|---|---|
| 1 | `01-messstand.sh` | `rohdaten-01.txt`, `rohdaten-08-vorher.txt` | Messstand, Basis-Image, Erreichbarkeit der Hosts, Docker-Objekte vorher |
| 2 | `02-oci-quellen.sh` mit `inner-02*.sh` | `rohdaten-02.txt` | Quellen als OCI-Feature: Feature-Index, Tags, Manifest, Metadaten, `install.sh` |
| 3 | `03-lieferkette.sh` mit `inner-03*.sh` | `rohdaten-03.txt` | Lieferkette der Upstream-Releases und der Debian-Pakete |
| 4 | `04-projekte.sh` | `rohdaten-04.txt`, `rohdaten-04-erzeugt-<variante>.txt` | `u-boot generate devcontainer` für `gh`/`glab` über die externe Quelle, Katalog-Form (`kat`), Egress-Liste |
| 5 | `05-bau.sh` mit `Dockerfile.cli` | `rohdaten-05.txt`, `rohdaten-05-<variante>.log` | Bau der erzeugten Projekte mit der Dev-Container-CLI, kalt; Inhalt der Images |
| 6 | `06-dockerfile-bloecke.sh` mit `block-*.dockerfile` | `rohdaten-06.txt`, `rohdaten-06-<variante>.log` | Dockerfile-Block-Varianten (Debian-Paket, gepinntes Archiv), Negativprobe des Pins |
| 7 | `07-laufzeit.sh` mit `inner-07.sh` | `rohdaten-07-<variante>-<modus>.txt` | Laufzeit: Hosts, Anmeldung per Umgebungsvariable, Dateien, Git; ohne und hinter der Egress-Restriktion |
| 8 | `08-aufraeumen.sh` | `rohdaten-08.txt`, `rohdaten-08-nachher.txt` | Aufräumen der Docker-Objekte mit dem Präfix `u-boot-messung-forge` |

**Varianten.** `bas`/`basn` = erzeugtes Projekt ohne Forge-Werkzeug (mit/ohne Sandbox-Profil, mit Egress-Restriktion bei `bas`); `gh`, `ghx`, `glabs`, `glabp`,
`glabz`, `ghn` = erzeugtes Projekt mit einem Feature über `source:`-Überschreibung (`ghn` ohne Sandbox-Profil); Endung `k` = Handvariante (Schlüssel ohne
Schema, siehe B1); `kat` = erzeugtes Projekt mit dem Katalog-Feature `git`; `ega` = `bas` mit `egress.allow gitlab.com`; `basd`, `ghdeb`, `glabdeb`, `ghtgz`,
`glabtgz` = Handvarianten mit Dockerfile-Block (`docker build`). Feature-Kürzel: `gh` = `ghcr.io/devcontainers/features/github-cli`, `ghx` =
`ghcr.io/devcontainers-extra/features/gh-cli`, `glabs` = `ghcr.io/skriptfabrik/devcontainer-features/gitlab-cli`, `glabp` =
`ghcr.io/postfinance/devcontainer-features/gitlab-cli`, `glabz` = `ghcr.io/srzstephen/devcontainer-features/gitlab-cli`.

## Was Messvariante ist und was Auslieferung (Risiko 2 des Slice)

| Gegenstand | Von `u-boot` erzeugt | Handvariante, nicht Auslieferung |
|---|---|---|
| Projekt, `devcontainer.json`, `Dockerfile`, `egress-init.sh` | alle Varianten ohne Endung `k` und ohne `d`/`deb`/`tgz` (`rohdaten-04-erzeugt-*.txt`) | — |
| Feature-Schlüssel in `devcontainer.json` | `https://ghcr.io/…:1` (externe Quelle) und `ghcr.io/…:1` (Katalogeintrag, nur `kat`) | Endung `k`: `sed` ersetzt `"https://ghcr.io/` durch `"ghcr.io/` (`05-bau.sh`); der Katalog hat keinen Eintrag für `gh`/`glab`, die Handvariante nimmt die Schlüssel-Form eines Katalogeintrags vorweg |
| Installation im Dockerfile | — (`u-boot` hat keinen Dockerfile-Block für `gh`/`glab`) | `ghdeb`, `glabdeb`, `ghtgz`, `glabtgz`: Block vor der letzten Zeile `USER vscode` des erzeugten Dockerfiles, Form des Ollama-Blocks (`block-*.dockerfile`); `docker build` statt Dev-Container-CLI |
| Optionen eines Features | — (`generate` rendert `"<Schlüssel>": {}`, keine Optionen) | keine Variante setzt Optionen |

Das erzeugte Ergebnis eines Katalog-Eintrags für `gh`/`glab` ist **nicht gemessen**, weil es ihn nicht gibt; gemessen ist seine Vorwegnahme durch die Handvarianten.

## Aussagen

Umgebung für alle Aussagen: der Messstand oben. *Läufe*: Zahl der Läufe, deren Rohdaten im Verzeichnis liegen; **gewertet ist der letzte Lauf**. Vor dem
letzten Lauf liefen die Skripte während der Entwicklung mehrfach (früherer Läufe: 02: 1, 03: 4, 04: 2, 05: 2, 07: 2, für `ghk` und `glabsk` 3); die Exit-Codes
und Versionen wichen dabei nicht ab, soweit sie beim Entwickeln verglichen wurden (nicht Zeile für Zeile; die Rohdaten der früheren Läufe sind überschrieben).
Die Erzeugnisse von 04 sind im letzten Lauf gegen den vorletzten byte-gleich (`cmp` über neun `rohdaten-04-erzeugt-*.txt` im Verlauf der Messung, kein Rohdatum).

### A. Quellen und ihre Lieferkette

| Nr. | Aussage | Rohdaten | Läufe |
|---|---|---|---|
| A1 | Der Feature-Index (`containers.dev/static/devcontainer-index.json`, abgerufen 2026-10-23) führt für `gh` das offizielle Feature `ghcr.io/devcontainers/features/github-cli` (1.1.3, Sammlung *Reference Implementation Features*) und Drittanbieter-Features (`devcontainers-extra` `gh-cli` 1.0.13, `sliekens`, `postfinance`, `joshuanianji`, `cirolosapio`). Für `glab` führt er **kein** Feature der Referenzsammlung; Einträge mit `gitlab-cli`: `skriptfabrik` 1.1.1, `postfinance` 1.0.1, `srzstephen` 1.0.1. Sechs Namen unter `devcontainers`, `devcontainers-extra`, `devcontainers-contrib` (`gitlab-cli`, `glab`, `gh-cli`) liefern HTTP 404 (nur diese Namen, keine Aussage über andere) | `rohdaten-02.txt` | 2 |
| A2 | Manifest-Digest = SHA-256 der Antwort und Layer-Digest = SHA-256 des Layers, für alle fünf Features. Digests von `:latest`: `github-cli` `sha256:bd7ab48a…`, `gh-cli` `sha256:bd667bbc…`, `skriptfabrik` `sha256:f87aea79…`, `postfinance` `sha256:e2b0903c…`, `srzstephen` `sha256:b72b7173…` (vollständig in der Rohdatei) | `rohdaten-02.txt` | 2 |
| A3 | Am Feature-Artefakt: Referrers-Schnittstelle HTTP 404 und keine Tags der Form `sha256-*`, `*.sig`, `*.att` — für alle fünf. Das heißt: kein Signatur- oder Attestierungs-Artefakt in diesen beiden Mechanismen; andere Mechanismen sind nicht geprüft | `rohdaten-02.txt` | 2 |
| A4 | Metadaten (Annotation `dev.containers.metadata`): `mounts`, `privileged`, `init`, `capAdd`, `securityOpt`, `entrypoint`, `dependsOn` sind bei allen fünf Features leer. Optionen: `github-cli` `extensions`, `installDirectlyFromGitHubRelease` (Standard `true`), `version` (Standard `latest`); `gh-cli` `version`; `skriptfabrik` `version`; `postfinance` `version`, `downloadUrl`; `srzstephen` `version` | `rohdaten-02.txt` | 2 |
| A5 | `github-cli` 1.1.3, Standardpfad (`installDirectlyFromGitHubRelease=true`): `wget` lädt `gh_<version>_linux_<arch>.deb` von `github.com/cli/cli/releases`, `dpkg -i`; **keine Prüfsumme, keine Signaturprüfung im Skript**; die Version `latest` wird aus den Git-Tags von `github.com/cli/cli` bestimmt. Mit `…=false`: apt-Repository `cli.github.com/packages`, Schlüssel mit dem fest eingetragenen Fingerabdruck `7F38BBB5…` über Keyserver geholt (`keyserver.ubuntu.com`, `keys.openpgp.org`, …); dieser Pfad ist **nicht gebaut** | `rohdaten-02.txt` (`grep -n`-Zeilen des Skripts, Nummern der Zeilen von `install.sh`), Bau `rohdaten-05-ghk.log` (Zeilen „Downloading github CLI … gh_2.103.0_linux_amd64.deb“) | 2 / 1 |
| A6 | `gh-cli` (`devcontainers-extra`): 23 Zeilen; lädt `nanolayer` v0.5.4 und installiert mit dem verschachtelten Feature `gh-release` 1.0.25 (`repo=cli/cli`, `version=$VERSION`, im Bau „Soft-dependency … gh-release is not required“); **welche Prüfung `nanolayer`/`gh-release` machen, ist nicht gelesen** | `rohdaten-02.txt`, `rohdaten-05-ghxk.log` | 2 / 1 |
| A7 | `skriptfabrik` `gitlab-cli` 1.1.1: installiert per `apt` `curl ca-certificates jq git`, bestimmt `latest` über die GitLab-API, lädt `glab_<version>_linux_<arch>.deb` von `gitlab.com/gitlab-org/cli/-/releases/…`, `dpkg -i`; **keine Prüfsumme, keine Signaturprüfung** | `rohdaten-02.txt` | 2 |
| A8 | `postfinance` `gitlab-cli` 1.0.1: `install.sh` (5 Zeilen) führt ein mitgeliefertes Binary aus (`installer_amd64` 6 230 200 Byte, `installer_arm64` 5 767 352 Byte, im 12-MB-Layer); **was es lädt und prüft, ist nicht gelesen** (Binary nicht untersucht) | `rohdaten-02.txt` | 2 |
| A9 | `srzstephen` `gitlab-cli` 1.0.1: lädt Archiv und `checksums.txt` aus dem GitLab-Paketregister (`…/packages/generic/glab/<version>/`), prüft SHA-256 mit `verify_checksum`; Prüfsumme und Archiv stammen aus derselben Quelle | `rohdaten-02.txt` | 2 |
| A10 | `gh` Upstream `cli/cli` v2.103.0 (veröffentlicht 2026-10-21): Release `immutable: true`; `gh_2.103.0_checksums.txt` stimmt für das `linux_amd64`-Archiv mit Datei und GitHub-Digest überein (`fdd768eb…`); **keine** Assets mit Signatur-, Provenance- oder SBOM-Namen; die Attestierungs-Schnittstelle liefert HTTP 200 mit zwei Attestierungen (Initiator `github` und `user`), Inhalt nicht dekodiert (Bundle-Datei ist nicht als JSON lesbar), Prüflogik nicht gemessen; Release-Workflow enthält `attestations: write`, `actions/attest` (Pin auf Commit) und einen GPG-Schritt für das apt-Repository | `rohdaten-03.txt` | 1 |
| A11 | `gh` apt-Repository `cli.github.com`: `InRelease` trägt `GOODSIG`/`VALIDSIG` des Schlüssels `7F38BBB5…`, derselbe Fingerabdruck wie im offiziellen Feature; das Repository führt gh 2.103.0 | `rohdaten-03.txt` | 1 |
| A12 | `glab` Upstream `gitlab-org/cli` v1.124.0 (veröffentlicht 2026-10-20): `checksums.txt` stimmt für das `linux_amd64`-Archiv mit der Datei überein (`63df6848…`); dieselbe Datei im Paketregister ist byte-gleich; **keine** Assets mit Signatur-, Provenance- oder SBOM-Namen; `.goreleaser.yml`/`.gitlab-ci.yml` (Suchmuster `sign`, `cosign`, `attest`, …) enthalten Signierschritte nur für macOS- und Windows-Binaries; keine Attestierung gefunden (Suchmuster, nicht abschließend) | `rohdaten-03.txt` | 1 |
| A13 | Rhythmus: `gh` 8 Releases von 2026-07-31 bis 2026-10-21 (v2.97.0 bis v2.103.0 und v2.102.1), alle `immutable`; `glab` 8 Releases von 2026-09-08 bis 2026-10-20, etwa wöchentlich | `rohdaten-03.txt` | 1 |
| A14 | Debian 13: `gh` 2.46.0-3, `glab` 1.53.0-1+b3 (Abhängigkeit nur `libc6`), Quelle `deb.debian.org`, signiert vom Debian-Archivschlüssel (`Signed-By` der Paketquellen) | `rohdaten-03.txt` | 1 |

### B. Form der Auslieferung (`u-boot generate devcontainer`)

| Nr. | Aussage | Rohdaten | Läufe |
|---|---|---|---|
| B1 | Für ein Feature über `source:` (mit Freigabeliste) schreibt `generate` den Schlüssel **mit Schema**: `"https://ghcr.io/devcontainers/features/github-cli:1": {}` (so für `gh`, `ghx`, `glabs`, `glabp`, `glabz` und `ghn` ohne Sandbox-Profil). Die Dev-Container-CLI 0.89.0 bricht **alle sechs** Bauten nach 0,8 bis 0,9 s mit `Feature '<Schlüssel>' could not be processed. You may not have permission to access this Feature, or may not be logged in.` ab (Exit 1). Damit ist der Weg über `source:` in der heute erzeugten Form **nicht baufähig**; die Form des Schlüssels ist die Ursache (B2: derselbe Aufbau mit Schlüssel ohne Schema baut). Das misst die Annahme zu N2 von [`ADR-0032`](../../plan/adr/0032-katalog-features-socket-und-lockerungen.md) ([`slice-v1-katalog-features-quellen-und-freigabe`](../../plan/planning/open/slice-v1-katalog-features-quellen-und-freigabe.md), dort als Schlussfolgerung geführt) für diese fünf Features | `rohdaten-04-erzeugt-gh.txt` u. a., `rohdaten-05.txt` | 3 / 3 |
| B2 | Dasselbe Projekt mit dem Schlüssel **ohne Schema** (Handvariante `k`) baut für alle fünf Features und `ghnk` | `rohdaten-05.txt` | 3 |
| B3 | Ein Katalog-Feature (`git`) erzeugt den Schlüssel ohne Schema: `"ghcr.io/devcontainers/features/git:1": {}` — die Form der Handvarianten ist die Form eines Katalogeintrags | `rohdaten-04-erzeugt-kat.txt` | 2 |
| B4 | Im Sandbox-Profil weist `generate` ein Feature aus einer externen Quelle mit der Warnung aus, `u-boot` lese seine Metadaten nicht und das Profil mache für das Feature keine Zusage (Code [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)); eine Meldung nach [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Drittanbieter) erscheint für die externe Quelle nicht | `rohdaten-04.txt` | 2 |
| B5 | Das Sandbox-Projekt trägt `remoteEnv` mit `GIT_TOKEN` (`${localEnv:GIT_TOKEN}`), aber keine Variable für `gh` oder `glab`; die Datei enthält keinen Wert | `rohdaten-04-erzeugt-gh.txt` | 3 |
| B6 | Die Basis-Liste der Egress-Allowlist trägt `github.com`, `api.github.com`, `codeload.github.com`, `raw.githubusercontent.com`, `objects.githubusercontent.com`, `deb.debian.org`, `security.debian.org`; `gitlab.com` fehlt (`egress list` jedes Sandbox-Projekts; `ega` mit `egress.allow gitlab.com`: Eintrag Quelle `user`) | `rohdaten-04.txt` | 3 |

### C. Bau auf dem Standard-Basis-Image

Alle Bauten kalt (`--no-cache`), Image-Größe in Byte (`docker image inspect`), Dauer in Sekunden.

| Nr. | Variante | Ergebnis | Rohdaten |
|---|---|---|---|
| C1 | `bas` / `basn` (Referenz Dev-Container-CLI) | Exit 0, 25,2 s / 4,8 s, 1 419 218 247 / 1 412 193 634 Byte; weder `gh` noch `glab` im Image | `rohdaten-05.txt` |
| C2 | `ghk` (offizielles Feature) | Exit 0, 59,2 s, 1 479 255 931 Byte (+60 037 684 gegen `bas`); `/usr/bin/gh`, `gh version 2.103.0 (2026-10-21)`; ein Paket mehr (`gh 2.103.0`) | `rohdaten-05.txt`, `-ghk.log` |
| C3 | `ghxk` (`devcontainers-extra`) | Exit 0, 71,0 s, 1 476 974 056 Byte; `/usr/local/bin/gh` (`stat` ohne Folgen eines Links: Mode 777, 47 Byte, also nach Größe und Modus ein symbolischer Link; Ziel nicht ausgegeben), `gh version 2.103.0`; kein Debian-Paket mehr | `rohdaten-05.txt`, `-ghxk.log` |
| C4 | `glabsk` (`skriptfabrik`) | Exit 0, 70,3 s, 1 501 865 746 Byte; `/usr/bin/glab` 1.124.0; Paket `glab 1.124.0`, dazu Aktualisierungen von `curl`, `libcurl*`, `ca-certificates` (11 Zeilen Paketunterschied) | `rohdaten-05.txt`, `-glabsk.log` |
| C5 | `glabpk` (`postfinance`) | Exit 0, 59,9 s, 1 504 119 000 Byte; `/usr/local/bin/glab` 1.124.0 (49 426 592 Byte); kein Debian-Paket | `rohdaten-05.txt`, `-glabpk.log` |
| C6 | `glabzk` (`srzstephen`) | Exit 0, 69,6 s, 1 500 216 385 Byte; `/usr/local/bin/glab` 1.124.0 (49 426 592 Byte); Aktualisierungen von `curl`, `libcurl*`, `ca-certificates` (10 Zeilen) | `rohdaten-05.txt`, `-glabzk.log` |
| C7 | `ghnk` (ohne Sandbox-Profil) | Exit 0, 36,2 s, 1 472 224 607 Byte (+60 030 973 gegen `basn`); `gh` 2.103.0 | `rohdaten-05.txt` |
| C8 | `basd` (Referenz `docker build`) | Exit 0, 10,7 s, 1 419 125 578 Byte | `rohdaten-06.txt` |
| C9 | `ghdeb` / `glabdeb` (Debian-Paket) | Exit 0, 19,3 s / 19,0 s; +48 809 142 / +40 798 568 Byte gegen `basd`; `gh version 2.46.0 (2025-01-13 Debian 2.46.0-3)` / `glab` meldet `Current glab version: 1.53.0` | `rohdaten-06.txt` |
| C10 | `ghtgz` / `glabtgz` (gepinntes Archiv, SHA-256 geprüft) | Exit 0, 16,0 s / 17,2 s; +56 100 573 / +67 439 177 Byte gegen `basd`; `gh version 2.103.0` / `glab 1.124.0 (d43ccacec)`; Build-Zeile `gh_2.103.0_linux_amd64.tar.gz: OK` / `glab_1.124.0_linux_amd64.tar.gz: OK`; Binary Eigentümer `root:root` Mode 755 | `rohdaten-06.txt` |
| C11 | Negativprobe: erste Stelle der SHA-256 in `ghtgz` verfälscht | Build Exit 1 (`computed checksum did NOT match`, `FAILED`), kein Image | `rohdaten-06.txt` |

Versionen, die der Bau zog (`latest` am 2026-10-23): `gh` 2.103.0, `glab` 1.124.0. Die Features pinnen die Werkzeug-Version **nicht**: ohne Option `version` gilt `latest` zur Bauzeit
(A5 bis A9). Die Archiv-Pins der Varianten `*tgz` stehen in `block-ghtgz.dockerfile` und `block-glabtgz.dockerfile`; der `arm64`-Pin (je Datei) ist aus den Prüfsummen-Dateien
der Releases gelesen, **nicht** gebaut.

### D. Hosts im laufenden Container (Laufzeit, Docker, `--cap-add=NET_ADMIN`)

Verfahren: ein protokollierender `dnsmasq` ersetzt den Resolver (Modus `discover`); `restricted` führt das **erzeugte** Startscript `egress-init.sh` aus
(Liste B6) und lässt danach denselben `dnsmasq` protokollieren; `restricted-allow` (nur `glab`) ersetzt dieses Skript durch das von `u-boot` für `ega`
erzeugte. Jeder Aufruf steht mit Exit-Code, DNS-Namen und geänderten Dateien in den Rohdaten (`rohdaten-07-<variante>-<modus>.txt`). Läufe: 1 gewertet je Variante und Modus.

| Nr. | Aussage | Rohdaten |
|---|---|---|
| D1 | `gh` 2.103.0 spricht in allen gemessenen Aufrufen nur `api.github.com` (mit Platzhalter-Token: `auth status`, `api rate_limit`, `repo view`, `release download`, `repo clone`) und `github.com` (`git ls-remote`) an; beide stehen in der Basis-Liste. `restricted` liefert dieselben Exit-Codes wie `discover`. Gleich für `ghxk`, `ghtgz`; auch `ghdeb` (2.46.0) | `rohdaten-07-ghk-*`, `-ghxk-*`, `-ghtgz-*`, `-ghdeb-*` |
| D2 | `glab` 1.124.0 spricht in allen gemessenen Aufrufen nur `gitlab.com` an (API, `release list`, `release download`, `git ls-remote`); der Download von `checksums.txt` ging über `gitlab.com/api/v4/projects/…/packages/generic/…` ohne weiteren Host. `gitlab.com` steht **nicht** in der Basis-Liste | `rohdaten-07-glabsk-discover.txt` u. a. |
| D3 | Hinter der erzeugten Standard-Allowlist scheitert jeder `glab`-Aufruf, der das Netz braucht (`api`, `repo view`, `release list`, `release download`, `auth status`), mit `dial tcp: lookup gitlab.com on 127.0.0.1:53: server misbehaving`, und `git ls-remote https://gitlab.com/…` mit Exit 128; der Client wiederholt die Anfragen: Gesamtdauer des Containerlaufs 127 bis 131 s statt 15 bis 21 s. Mit `egress.allow gitlab.com` (Variante `ega`, `restricted-allow`) gelingen `api`, `repo view`, `release list`, `release download` und `git ls-remote` wie im Modus `discover` | `rohdaten-07-glabsk-restricted.txt`, `-restricted-allow.txt` |
| D4 | Ergebnis D1 bis D3 gilt für `glabsk`, `glabpk`, `glabzk`, `glabtgz` (alle 1.124.0) und, mit gleichen Hosts, `glabdeb` (1.53.0) | `rohdaten-07-glab*-*` |
| D5 | In keiner Sitzung erschien ein weiterer Name als `api.github.com`, `github.com` oder `gitlab.com` — auch nicht für Telemetrie oder Aktualisierungsprüfung, obwohl `glab` 1.124.0 `telemetry=true` und `check_update=true` führt und `gh` 2.103.0 `GH_TELEMETRY` kennt (Dokumentation in `gh help environment`). Die Aufrufe liefen ohne Terminal; ob Telemetrie oder Update-Prüfung in einer interaktiven Sitzung andere Hosts braucht, ist **nicht gemessen** | `rohdaten-07-ghk-discover.txt`, `-glabsk-discover.txt` |
| D6 | Hosts des **Baus** (nicht durch die Egress-Restriktion begrenzt, aus den Skripten gelesen, nicht per DNS gemessen): offizielles Feature: `github.com` (Git-Tags, Release-Download; Weiterleitung nicht gemessen), `deb.debian.org`; mit `…=false` zusätzlich `cli.github.com` und Keyserver; `skriptfabrik`: `gitlab.com`, `deb.debian.org`; `srzstephen`: `gitlab.com`, `deb.debian.org`; Dockerfile-Block: `github.com` beziehungsweise `gitlab.com` (nur im Archiv-Download) | `rohdaten-02.txt`, `block-*.dockerfile` |

### E. Umgang mit Tokens

Alle Aufrufe mit dem Platzhalter `MESSUNG-KEIN-ECHTES-TOKEN` (kein gültiges Token) oder ohne Variable. Es wurde **kein Token angelegt oder verwendet**; ein
Aufruf mit gültigem Token ist **nicht gemessen** (siehe *Nicht gemessen*).

| Nr. | Aussage | Rohdaten |
|---|---|---|
| E1 | Variablen laut Hilfe des Werkzeugs: `gh`: `GH_TOKEN`, `GITHUB_TOKEN` (in dieser Rangfolge; für `github.com` und Unterdomänen von `ghe.com`), `GH_ENTERPRISE_TOKEN`, `GITHUB_ENTERPRISE_TOKEN`; `glab`: `GITLAB_TOKEN`, `GITLAB_ACCESS_TOKEN`, `OAUTH_TOKEN` (haben Vorrang vor gespeicherten Zugangsdaten), `GITLAB_HOST` | `rohdaten-07-ghk-discover.txt` (`help-environment`), `rohdaten-07-glabsk-discover.txt` (`glab auth login --help`) |
| E2 | Mit der Variable und **ohne** `auth login` versucht das Werkzeug die Anfrage: `gh` mit `GH_TOKEN` und mit `GITHUB_TOKEN` antwortet `HTTP 401: Bad credentials`; `glab` mit `GITLAB_TOKEN`, `GITLAB_ACCESS_TOKEN` und `OAUTH_TOKEN` `401 Unauthorized`; `glab auth status` meldet `Token found in environment variable GITLAB_TOKEN` mit maskiertem Wert. Dass dieselbe Anfrage mit einem **gültigen** Token gelingt, ist nicht gemessen | `rohdaten-07-*-discover.txt` |
| E3 | Ohne Variable: `gh api rate_limit` endet mit Exit 4 und dem Hinweis, `gh auth login` auszuführen oder `GH_TOKEN` zu setzen; `gh auth status` Exit 1 (`not logged into any GitHub hosts`). `glab api`, `glab repo view`, `glab release list`, `glab release download` gelingen für das öffentliche Projekt `gitlab-org/cli` **ohne** Token (Exit 0); `glab auth status` Exit 1 (`No token found`); `glab repo clone` scheitert ohne gültiges Token (`Get https://gitlab.com/api/v4/user: 401`) | `rohdaten-07-*-discover.txt` |
| E4 | Dateien: `gh` legt beim ersten Aufruf (auch `--version`) `~/.local/state/gh/device-id` (36 Byte) an, sonst nichts unter `~/.config/gh/`. `glab` legt beim ersten Aufruf (auch `--version`) `~/.config/glab-cli/config.yml` (7 534 Byte), `config.lock` (0 Byte) und `aliases.yml` (32 Byte) an; die Konfiguration führt unter anderem `git_protocol`, `telemetry`, `check_update`. Der Platzhalter steht in **keiner** dieser Dateien und in keiner Datei unter `~/.config`, `~/.cache`, `~/.local` (Zähler `platzhalter-treffer=0`). Die Dateien liegen im Dateisystem des Containers, nicht im Workspace-Volume | `rohdaten-07-*-discover.txt` (`geaenderte Dateien`, `Dateien unter HOME`) |
| E5 | Git-Zugriff: Das Dockerfile des Sandbox-Profils trägt den Credential-Helper `!f() { test -n "$GIT_TOKEN" && printf "username=x-access-token\npassword=%s\n" "$GIT_TOKEN"; }; f` (System-Konfiguration). `git ls-remote` gegen `github.com/cli/cli` und `gitlab.com/gitlab-org/cli` gelingt (Exit 0) mit und ohne `GIT_TOKEN`, weil beide Projekte öffentlich sind; ob der Helper bei einem privaten Repository mit `GITLAB_TOKEN`/`GIT_TOKEN` gelingt, ist nicht gemessen | `rohdaten-07-ghk-discover.txt` (Kopf), `rohdaten-07-glabsk-discover.txt` |
| E6 | `gh auth setup-git --hostname github.com` mit gesetztem `GH_TOKEN` (Platzhalter) schreibt `~/.gitconfig` (181 Byte) mit je einem leeren und einem `!/usr/bin/gh auth git-credential`-Eintrag für `github.com` und `gist.github.com`; kein Token im Inhalt. Der Befehl ersetzt damit für diese Hosts die Helper-Kette des Images | `rohdaten-07-ghk-discover.txt` |
| E7 | Standard für Git über `https`: `gh config list` → `git_protocol=https`; `glab config get git_protocol` → `ssh`; `glab auth status` meldet „Git operations for gitlab.com configured to use ssh protocol“. Der Clone-Weg selbst ist nicht gemessen (`glab repo clone` scheitert vorher am Token, E3). `glab config set git_protocol https` scheitert außerhalb eines Repositorys ohne `--global` (`not a Git repository`), mit `--global` Exit 0 | `rohdaten-07-ghk-discover.txt`, `-glabsk-discover.txt` |
| E8 | Version-Unterschied: `gh auth status` mit ungültigem Token endet in 2.103.0 mit Exit 1, in der Debian-Version 2.46.0 mit Exit 0 und derselben Meldung `The token in GH_TOKEN is invalid`; in 2.103.0 sind `GH_TELEMETRY`/`DO_NOT_TRACK` dokumentiert, in 2.46.0 nicht; `glab` 1.53.0 führt `check_update=false` und kennt `telemetry` nicht | `rohdaten-07-ghdeb-discover.txt`, `-glabdeb-discover.txt` |
| E9 | `glab auth status` mit gesetzter Variable gibt den Hinweis aus, den Wert mit `env | grep -E 'GITLAB_TOKEN|…'` anzuzeigen; der Wert selbst steht nur maskiert in der Ausgabe | `rohdaten-07-glabsk-discover.txt` |

### F. Nicht gemessen

- **Jeder Aufruf mit einem gültigen Token**: Erfolgspfad von `gh api`, `gh release download` (Weiterleitungs-Hosts), `gh repo clone`, `glab repo clone`, Push; Hosts der
  Anmeldung (`auth login`, Device-Flow); Telemetrie- und Update-Hosts in einer interaktiven Sitzung; Schreiben/Lesen eines Tokens aus `hosts.yml` beziehungsweise Schlüsselbund nach `auth login`.
  Grund: kein Zugangsdatum durfte angelegt oder verwendet werden (Auftrag), das Messen mit einem echten Token wäre eine Entscheidung des Projektinhabers.
- Der Credential-Helper des Sandbox-Dockerfiles gegen ein **privates** Repository und gegen GitLab (Benutzername `x-access-token`).
- `arm64`, macOS, Windows, Docker Desktop, Podman; Docker-Daemon in Colima; Host mit anderer Engine.
- Die Features mit gesetzten Optionen (`version`, `installDirectlyFromGitHubRelease=false`, `downloadUrl`); der apt-Pfad des offiziellen Features (A5) ist gelesen, nicht gebaut.
- Der Inhalt von `installer_amd64` (`postfinance`) und die Prüfung von `nanolayer`/`gh-release` (A6, A8).
- `gh attestation verify` und die Dekodierung der Attestierungen (A10); die Signaturen der Debian-Pakete über die Paketquellen hinaus.
- Hosts des **Baus** per DNS (D6 ist aus den Skripten gelesen); Verhalten bei fehlendem Netz während des Baus.
- Der Bau eines vom Generator erzeugten Dockerfiles mit einem `gh`-/`glab`-Block (es gibt keinen) und der Bau mit einem Katalogeintrag für `gh`/`glab` (es gibt keinen).
- Selbst gehostete Instanzen (GitHub Enterprise, GitLab Self-Managed) und weitere Hosts; das Sandbox-Profil mit `nestedRuntime`.
- Die Stabilität der Aussagen über Zeit: `latest` und Release-Stand sind die vom 2026-10-23.

## Token-Sicherheit der Rohdaten

Die Skripte verwenden nur den Platzhalter `MESSUNG-KEIN-ECHTES-TOKEN` und rufen die Forge-APIs anonym aus Containern ab (keine Anmeldung des Hosts, kein `gh auth`). `redact` in
`lib.sh` ersetzt vor dem Schreiben jeder Rohdatei vorsignierte Query-Strings, `Authorization`/`Bearer`-Werte, bekannte Token-Präfixe (`ghp_`, `gho_`, `ghu_`, `ghs_`,
`ghr_`, `github_pat_`, `glpat-`) und den Home-Pfad. Die signierte Adresse eines Attestierungs-Bundles (A10) wird abgerufen und nicht ausgegeben. Vor dem Commit durchsucht
(`grep -rIn` über dieses Verzeichnis und `tools/messung/forge-cli` nach den Präfixen, `Bearer `, `Authorization`, `X-Amz`, `Signature=`, `sig=`,
dem Namen des Home-Verzeichnisses): Treffer **nur** in den Skripten, und zwar die Regeln von `redact` und die Header-Zeile `Authorization: Bearer $tok` der anonymen
`ghcr.io`-Abfrage (Variable, kein Wert); in den Rohdateien kein Treffer (die Beschreibung in diesem Absatz nennt die Muster und trifft sie selbst). Die Suche nach Token-Werten (Präfixe `ghp_`, `gho_`, `ghu_`, `ghs_`, `ghr_`, `github_pat_`, `glpat-`, `sk-ant-`, je mindestens 20 Zeichen danach) über `docs/messungen/` und `tools/messung/` ist leer. Die Regel für `Authorization:`-Zeilen in `redact` ließ den Wert hinter `Bearer` stehen und wurde nach den Läufen berichtigt (Gegenprobe mit einer Beispielzeile: danach entfällt der Wert); die Rohdaten sind mit der früheren Regel entstanden und tragen keine solche Zeile (Suche oben).

## Lesehilfen

- `AUSSAGE <schlüssel>=<wert>`-Zeilen sind die stabilen Werte, die ein zweiter Lauf wiederholen muss; Zeilen `---- <Zeit> <Etikett>` trennen die Aufrufe von `inner-07.sh`.
- Im Modus `restricted` listet *Antworten ausserhalb der Liste* nur Zeilen des Resolver-Protokolls; `dnsmasq[…]: config error` erscheint, wenn ein Name nicht freigegeben ist, sein genauer Wortlaut ist nicht weiter untersucht.
- Die Zahl der DNS-Zeilen je Name ist eine Zählung von `query[A]`/`query[AAAA]` einschließlich Wiederholungen des Clients, kein Maß für die Zahl der Verbindungen.
- `rohdaten-08-vorher.txt` und `-nachher.txt` zählen alle Docker-Objekte des Hosts; die Zunahme von Volumes und Containern über die Messung stammt von Arbeiten außerhalb dieser Messung (die Skripte starten alle Container mit `--rm`; Volumes legen sie nicht an), alle Objekte mit dem Präfix sind nachher entfernt (`objekte-mit-praefix-nachher=0`).
- Handvarianten sind im Namen kenntlich: Endung `k`, Namen `*deb`, `*tgz`, `basd`.
