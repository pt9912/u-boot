# Release-Ablauf — u-boot

| Dokument    | Release-Anleitung für Maintainer |
| ----------- | -------------------------------- |
| Projektname | `u-boot` |
| Bezug       | [`.github/workflows/publish.yml`](../../.github/workflows/publish.yml), [ADR-0007](../plan/adr/0007-distributionswege-ghcr.md), [ADR-0016](../plan/adr/0016-homebrew-distribution-per-tap.md) |
| Zielgruppe  | Maintainer, die einen Release schneiden |

## Zweck

Beschreibt, wie ein Release von `u-boot` entsteht: was vorab zu prüfen ist, was der
Tag-Push auslöst, was danach zu kontrollieren ist und was bei einem Fehler zu tun ist.
Ein Release ist **öffentlich** (GHCR-Image, Binaries, Homebrew-Tap) und lässt sich nicht
sauber zurücknehmen — der Tag-Push ist deshalb bewusst ein von Hand ausgeführter Schritt.

## 1. Was ein Tag auslöst

Der Push eines Tags `vMAJOR.MINOR.PATCH[-PRERELEASE]` startet den Workflow `publish`:

| Schritt | Ergebnis |
| --- | --- |
| SemVer-Prüfung | bricht vor Login/Build/Push ab bei ungültigen Tags (`v1.2`, `vNext`, Build-Metadaten mit `+`) |
| Runtime-Image bauen | `make build VERSION=<x.y.z>` (Docker-only) |
| OCI-Labels und `--version` prüfen | das Image meldet genau die Tag-Version, sonst Abbruch |
| GHCR-Push | `ghcr.io/pt9912/u-boot:<x.y.z>`; `:latest` **nur** für stabile Tags |
| Binaries | sechs Plattformen (Linux/macOS/Windows × amd64/arm64) als Release-Assets |
| Linux-Pakete | vier Pakete (`.deb`/`.rpm` × amd64/arm64) per `nfpm` ([ADR-0017](../plan/adr/0017-linux-pakete-nfpm.md)) als Release-Assets |
| `SHA256SUMS` | Prüfsummen der Binaries und Pakete als weiteres Asset |
| Paket-Smoke | zwei Jobs installieren `.deb` (Ubuntu) und `.rpm` (Fedora) aus dem Release, prüfen `--version` und `init`; der Asset-Name kommt aus der Asset-Liste des Release (`scripts/package-asset-name.sh`), nicht aus dem Tag |
| Homebrew-Formel | nur stabile Tags: `u-boot.rb` aus den Digests von `SHA256SUMS` |
| Job `tap` | nur stabile Tags: zieht die Formel nach `pt9912/homebrew-u-boot` nach |

Ein **Vorab-Tag** (`v0.6.0-rc.1`) durchläuft dieselbe Kette, wird aber als Prerelease
markiert, bekommt kein `:latest`, keine Formel und keinen Tap-Commit. Das ist die
Probe für den Release-Pfad (und hinterher löschbar). Bei Vorabversionen heißen die Pakete
in der Datei `0.6.0~rc.1` (nfpm), das Release-Asset `0.6.0.rc.1` (GitHub ersetzt `~`); `SHA256SUMS`
führt die Dateinamen mit `~`. Stabile Tags haben nur eine Schreibweise.

## 2. Voraussetzungen (einmalig)

- **Secret `HOMEBREW_TAP_GITHUB_TOKEN`** im Repo `pt9912/u-boot` (Settings → Secrets and
  variables → Actions): ein PAT mit Schreibrecht (Contents: Read & Write) auf
  `pt9912/homebrew-u-boot`. Ohne Secret **überspringt** der Job `tap` den Nachzug mit einer
  Warnung; der Release selbst hängt nicht am Tap.
- **Token prüfen:** `gh workflow run tap-check` (Workflow
  [`tap-check.yml`](../../.github/workflows/tap-check.yml)) meldet, ob das Secret auf das Tap
  **schreiben** darf — per echtem Klon und `git push --dry-run` (die API-Sicht allein täuscht:
  bei klassischen PATs zeigt `permissions.push` die Rechte des Benutzers, nicht die Scopes des
  Tokens). Ein Rerun des Jobs `tap` zeigt das nicht, wenn die Formel schon im Tap liegt.
- **Tap-Repo** `pt9912/homebrew-u-boot` (README und Smoke-Workflow aus
  [`packaging/homebrew-tap/`](../../packaging/homebrew-tap/); eine Änderung dieser zwei Dateien kommt nach §9 in den Tap).
- Schreibrecht auf Tags und auf `ghcr.io/pt9912/u-boot` (über `GITHUB_TOKEN` des Workflows).

## 3. Vorbereitung (ein Commit auf `main`)

Der Cut ist ein Release-Cut-Slice als eigenes Arbeitsartefakt, geführt wie jeder
Slice (`open/` → `next/` → `in-progress/`, nach dem Tag nach `done/`); seine
**Form** kommt aus der Skelett-Vorlage
[`.harness/baseline/v6.18.1/templates/docs/plan/planning/slice.template.md`](../../.harness/baseline/v6.18.1/templates/docs/plan/planning/slice.template.md),
der **Ablauf** steht vollständig in diesem Dokument:

1. **CHANGELOG** (`CHANGELOG.md`): `## [Unreleased]` bleibt als leerer Anker; der bisherige
   Inhalt wandert unter `## [x.y.z] - <Tag-Datum>` mit kurzem Lead-Absatz. Compare-Links am
   Fuß umstellen (`[Unreleased]: …/compare/vx.y.z...HEAD`, neue Zeile `[x.y.z]: …`).
2. **Versionsstrings** `<alt>-dev` → `<neu>-dev` an genau drei Stellen:
   `cmd/uboot/main.go` (`var version`), `repo.mk` (`VERSION ?=`), `Dockerfile`
   (`ARG UBOOT_VERSION=`), jeweils samt Kommentar-Erwähnungen. Der Tag-Build überschreibt
   den Wert per `UBOOT_VERSION`; der `-dev`-Stand gilt für lokale Builds.
3. **READMEs** (`README.md`, `README.de.md`): die Versionsangabe in der Status-Zeile (eine
   Zeile mit der aktuellen Version und den Zeigern auf Handbuch und CHANGELOG), ohne
   Anforderungs-, ADR- oder Spec-Kennungen ([`AGENTS.md`](../../AGENTS.md) §5). Die Zusagen,
   Bedingungen und Grenzen eines Releases stehen im Handbuch und im CHANGELOG, nicht im
   README; das README trägt weder eine Release-Tabelle noch eine Chronik und kein „released".
   Das Datum steht im CHANGELOG und wird vor dem Tag gegen das Tag-Datum gehalten (§5); §7
   prüft es nachher noch einmal.
4. **Roadmap** (`docs/plan/planning/in-progress/roadmap.md`): den Release unter „Offene Wellen"
   führen, bis der Tag gesetzt ist; danach nach „Abgeschlossene Wellen".
5. **Release-Cut-Slice** aus der Skelett-Vorlage `slice.template.md` anlegen und nach dem Tag nach
   `done/` verschieben (Planning-Sensor: der Ruhe-Marker der Roadmap gilt nur, wenn kein
   Slice in `in-progress/` liegt).

## 4. Vorab-Prüfung

```bash
make gates            # lint + test + coverage-gate + docs-check + baseline-verify + ablage-pruefung + tool-pruefung
make ci               # gates + govulncheck + image-scan (Trivy HIGH/CRITICAL)
make test-docker      # Integrationstests gegen die echte Docker Engine
make fullbuild        # ci + Runtime-Image-Build
make build VERSION=X.Y.Z   # Gegenprobe zu den Versionszeilen (schritt 2): das Image meldet
                           # X.Y.Z mit `docker run --rm u-boot:latest --version` und im
                           # OCI-Label org.opencontainers.image.version; der Standardbau
                           # meldet X.Y.Z-dev. lokal, nicht gepusht
```

**Trockenlauf des Tap-Pfads** (vor dem stabilen Tag, ohne Netzzugriff — Stellvertreter für `curl` und `git` im `PATH`, Container mit
`--network none` — und ohne Schreiben ins Tap; im Docker, zum Beispiel im Image `u-boot:test` aus
`make test`, mit dem Arbeitsbaum als schreibgeschütztem Mount).
Der Job `tap` läuft nur bei stabilen Tags; ein Vorab-Tag probt ihn nicht (§1, §7). Der Trockenlauf
führt deshalb die Skripte des Jobs gegen eine Fixture aus:

1. `scripts/homebrew-formula-fill.sh vX.Y.Z <SHA256SUMS-Fixture> scripts/homebrew-formula.rb.tmpl <Ausgabe>`
   erzeugt die Formel (Fixture: vier Plattform-Digests, Muster `sumsFull` in
   `cmd/uboot/homebrew_scripts_test.go`).
2. `scripts/tap-formula-check.sh vX.Y.Z <Ausgabe>` endet mit Exit `0`; mit dem Tag des Vorgängers
   endet er mit Exit `1`.
3. `scripts/tap-nachzug.sh` läuft mit `TAP_TOKEN` (Dummy) und `TAG=vX.Y.Z` und mit Stellvertretern für
   `curl` (kopiert die Fixture-Formel an das Ziel von `-o`) und `git` (protokolliert die Aufrufe,
   schreibt nichts) im `PATH`. Erwartet: das Protokoll zeigt `clone`, `add`, `commit -m u-boot X.Y.Z`
   und `push`; mit dem Tag des Vorgängers bricht das Skript vor `clone` ab.

Das Schreibrecht des Tokens belegt der Trockenlauf nicht; das tut nur `gh workflow run tap-check` (§2).
Der Maintainer startet es vor dem stabilen Tag, auch wenn es zuletzt grün war (der PAT läuft ab).
Der Lauf steht im Release-Cut-Slice mit Befehl und Ausgabe. **Er beweist nicht:** die Bedingung
`needs.publish.outputs.is_stable == 'true'` und den Schritt im Workflow, den echten `curl` gegen das
Release-Asset, den echten `git push` und das Token. Das bleibt dem ersten stabilen Lauf (§6)
vorbehalten; der Trockenlauf ersetzt ihn nicht.

Zusätzlich: `main` ist auf `origin/main` gepusht und der **letzte CI-Lauf auf `main` ist
grün** (`gh run list --limit 3`). Ein roter `image-scan` blockiert den Release — typischer
Anlass ist ein neuer Go-stdlib-Advisory; dann den Go-Pin heben (`Dockerfile` `ARG
GO_VERSION`, `repo.mk` `GO_VERSION`), siehe CHANGELOG-Einträge zu Toolchain-Bumps.

Sichtschritt vor einem stabilen Tag, kein Gate: Quelle und Tap von `README.md` und `.github/workflows/smoke.yml`
vergleichen (Schritt 1 aus §9); der Smoke läuft mit der Datei, die im Tap liegt.

## 5. Tag setzen (Maintainer-Aktion)

Vor jedem Tag (auch dem Vorab-Tag) `date +%F` prüfen und das Datum in CHANGELOG und im
Handbuch-`Stand` gegen das Tag-Datum halten; weicht es ab, ist das ein weiterer Cut-Commit
vor dem Tag.

```bash
git switch main && git pull --ff-only
git tag vX.Y.Z
git push origin vX.Y.Z      # löst publish aus
```

Optional zuerst ein Vorab-Tag (`vX.Y.Z-rc.1`), um die Kette ohne `:latest`, Formel und
Tap-Commit zu proben; er lässt sich danach löschen (`gh release delete`, `git push --delete
origin <tag>`, GHCR-Tag entfernen). Der Vorab-Tag läuft den Zweig `is_stable == true` und den Job
`tap` **nicht**; wer ihn als Probe führt, benennt in der Closure, welche Zweige er nicht lief.

## 6. Nach dem Push kontrollieren

```bash
gh run watch                                   # publish-Lauf (Job publish, dann tap)
gh release view vX.Y.Z --json assets -q '.assets[].name'
docker run --rm ghcr.io/pt9912/u-boot:X.Y.Z --version
```

Erwartet: sechs Binaries, vier Pakete (`.deb`/`.rpm`), `SHA256SUMS`, bei stabilen Tags `u-boot.rb`; `--version` gleich
`X.Y.Z`; `:latest` zeigt auf dieselbe Version (nur stabil). Der Job `tap` ist grün **und
wirklich gelaufen**: Bei fehlendem Secret endet er mit einer Warnung und Exit `0`, grün ist dann nicht
„nachgezogen“. Das Log des Schritts „Pull the formula into the tap“ zeigt `tap-nachzug: … updated to vX.Y.Z.`
oder `tap-nachzug: tap already at vX.Y.Z; nothing to do.`; die Warnung `HOMEBREW_TAP_GITHUB_TOKEN is not
configured` heißt: nicht gelaufen (§2, §8). Danach auf einem Mac:

```bash
brew tap pt9912/u-boot
brew trust pt9912/u-boot      # Homebrew 5 verlangt das für Drittanbieter-Taps
brew install u-boot
u-boot --version              # X.Y.Z
```

Im Tap-Repo läuft nach dem Formel-Commit der Smoke-Workflow (macOS: Install mit
`--build-from-source`, `brew test`, `--version`, `init`, `add postgres`; `doctor` ohne Docker nur
auf „kein Absturz" geprüft). `brew audit --strict --online` gatet den Lauf; ein Befund färbt ihn rot.

## 7. Nacharbeit

- Roadmap: Welle nach „Abgeschlossene Wellen", Tag-Commit-Hash und Datum eintragen.
- Release-Cut-Slice nach `done/` (Closure-Notiz mit Hash und Sensoren).
- CHANGELOG-Datum auf das tatsächliche Tag-Datum prüfen (`git log -1 --format=%cs vX.Y.Z`); dasselbe
  Datum steht im Handbuch-`Stand`.
- Vorab-Tags aufräumen: Release und Tag löschen (§5) und die GHCR-Tags der Vorab-Versionen entfernen;
  das braucht für `gh` die Scopes `read:packages` und `delete:packages`, sonst über die Web-Oberfläche.
- Job `tap` prüfen (§6): Er läuft nur bei stabilen Tags, ein Vorab-Tag probt ihn nicht. Schlägt er fehl,
  steht die Behebung auf `main`, der Tag bleibt; den Tap zieht der Maintainer manuell nach (§8).

## 8. Wenn etwas schiefgeht

| Symptom | Ursache / Maßnahme |
| --- | --- |
| `publish` bricht in „Validate SemVer tag" | Tag ungültig → Tag löschen, korrekt neu setzen (noch nichts veröffentlicht) |
| `--version` des Images ≠ Tag | Build-Arg-Drift → Abbruch vor dem Push; Versionsstrings/`make build VERSION` prüfen |
| Lauf bricht nach dem GHCR-Push | `:x.y.z` ist schon veröffentlicht; Workflow nach Behebung erneut ausführen (`--clobber` ersetzt Assets), nicht den Tag verschieben |
| Job `tap` warnt „secret not configured" | Secret setzen (§2), Job erneut ausführen (`gh run rerun --job`) |
| Job `tap` schlägt fehl (z. B. `403 Permission denied`) | `gh workflow run tap-check` ausführen; PAT-Rechte/Ablauf prüfen (fine-grained: Repo ergänzen, Contents = Read and write); manuell nachziehen: `TAP_TOKEN=… TAG=vX.Y.Z scripts/tap-nachzug.sh` |
| Job `tap` meldet „the u-boot.rb asset of vX.Y.Z does not carry four download URLs“ | Die Formel trägt keine `version`-Zeile; `scripts/tap-formula-check.sh` liest die vier `url`-Zeilen und weist einen falschen Tag oder eine fehlende URL ab. Das Asset ist dann tatsächlich falsch befüllt (`scripts/homebrew-formula-fill.sh`, `SHA256SUMS` prüfen); Tap nach der Behebung manuell nachziehen (`TAP_TOKEN=… TAG=vX.Y.Z scripts/tap-nachzug.sh`) |
| Formel fehlt im Tap | Asset `u-boot.rb` des Releases fehlt (Prerelease?) oder Tap-Job nicht gelaufen |
| Smoke im Tap rot, ein Job ohne Runner, Annotation „The configuration '…' is not supported“ | Das Runner-Label der Matrix in `smoke.yml` wird nicht unterstützt → Label in der Vorlage ändern und die Datei nach §9 in den Tap bringen |
| Paket-Smoke endet mit `curl: (22) 404` | Asset-Name falsch gewählt (Vorab-Tag: `~` im Dateinamen, `.` im Asset); `gh release view vX.Y.Z-rc.N --json assets` zeigt die Namen, `scripts/package-asset-name.sh` wählt sie; Workflow-Fix auf `main`, neuer Vorab-Tag (der alte Tag wird nicht verschoben) |
| Trivy-`image-scan` rot | Go-Pin heben (siehe §4), erst dann taggen |

Ein bereits veröffentlichter Tag wird **nicht verschoben**; ein Fehler wird mit dem nächsten
Patch-Release behoben.

## 9. Dateien des Tap-Repos nachziehen (`packaging/homebrew-tap/`)

Der Abschnitt gilt für `README.md` und `.github/workflows/smoke.yml` aus
[`packaging/homebrew-tap/`](../../packaging/homebrew-tap/). Die Formel `Formula/u-boot.rb` zieht der Job `tap` nach
(§1, §6) und gehört nicht hierher. Die Quelle in `packaging/homebrew-tap/` führt, der Tap folgt: Der Smoke läuft
mit der Datei, die im Tap-Repo `pt9912/homebrew-u-boot` liegt, nicht mit der Vorlage (`smoke.yml` Zeilen 10 bis 14:
Trigger `push` auf `Formula/u-boot.rb` und `workflow_dispatch`).

**Rechte, als Hinweis.** Die Seite „Create or update file contents“ der GitHub-REST-Dokumentation nennt für klassische
PATs den Scope `repo` und zusätzlich `workflow`, um Dateien unter `.github/workflows` zu ändern, und für fine-grained
PATs Contents (write) sowie Contents (write) mit Workflows (write). Welche Rechte das Token trägt, mit dem der Maintainer
pusht, zeigt nur der Push; über `HOMEBREW_TAP_GITHUB_TOKEN` (§2) sagt dieser Abschnitt nichts.

Die Befehle setzen `P` (der Pfad der Datei) und lesen die Quelle aus dem Arbeitsbaum:

```bash
P=.github/workflows/smoke.yml      # oder README.md
```

1. **Ausgangslage vergleichen.** Beide Ausgaben sind gleich, wenn Quelle und Tap übereinstimmen. Weichen sie ab, ist
   zuerst zu klären, was nur im Tap steht, bevor die Quelle geändert wird:

   ```bash
   git hash-object packaging/homebrew-tap/$P
   gh api repos/pt9912/homebrew-u-boot/contents/$P --jq .sha
   ```

2. **Quelle ändern.** `make gates` im Devcontainer-Container, auf `main` committen, Review und Verification der Quelle
   abwarten. Erst danach folgt der Push in den Tap.
3. **Push in den Tap** per `gh api` mit `PUT` auf `contents/<Pfad>`. Die Befehle lesen die Quelle aus dem Arbeitsbaum;
   vorher muss dieser dem geprüften Commit auf `main` entsprechen: Arbeitsbaum sauber und die Quelle gleich `origin/main`
   (beide Ausgaben leer):

   ```bash
   git fetch origin
   git status --short
   git diff origin/main -- packaging/homebrew-tap/
   ```

   `sha` ist der Blob-sha der Datei, die ersetzt wird (Schritt 1):

   ```bash
   B64=$(base64 -w0 < packaging/homebrew-tap/$P)   # Base64 der Quelle ohne Zeilenumbruch
   SHA=$(gh api repos/pt9912/homebrew-u-boot/contents/$P --jq .sha)
   gh api -X PUT repos/pt9912/homebrew-u-boot/contents/$P \
     -f message='<Commit-Text>' -f content="$B64" -f sha="$SHA" -q '.commit.sha'
   ```

   Der Push startet keinen Lauf: Der Trigger nennt nur `Formula/u-boot.rb` und `workflow_dispatch`.
4. **Nachvergleich.** Der Inhalt im Tap ist byte-gleich mit der Quelle (`cmp` endet mit Exit `0`), und der sha stimmt
   mit dem aus Schritt 1 für die Quelle überein:

   ```bash
   gh api repos/pt9912/homebrew-u-boot/contents/$P --jq .content | base64 -d | cmp - packaging/homebrew-tap/$P
   gh api repos/pt9912/homebrew-u-boot/contents/$P --jq .sha
   ```

5. **Lauf starten und lesen.** Der Dispatch ist ein Schreibzugriff auf den Tap; die Lesebefehle danach ändern nichts:

   ```bash
   gh workflow run smoke --repo pt9912/homebrew-u-boot
   gh api 'repos/pt9912/homebrew-u-boot/actions/workflows/smoke.yml/runs?per_page=5' \
     --jq '.workflow_runs[] | [.id,.created_at,.head_sha,.event,.status,.conclusion] | @tsv'
   ```

   Den neuen Lauf erkennt `created_at` (der jüngste `workflow_dispatch`-Lauf); `head_sha` ist der Commit aus Schritt 3.
   Zwei Dispatch-Läufe auf demselben Commit tragen denselben `head_sha`, und ein älterer Dispatch-Lauf ohne Job kann
   daneben in der Liste stehen. Dann Jobs, Annotationen und Log des neuen Laufs lesen (`<id>` und `<job-id>` aus den
   vorigen Ausgaben):

   ```bash
   gh api repos/pt9912/homebrew-u-boot/actions/runs/<id>/jobs \
     --jq '.jobs[] | [.id,.name,(.labels|join(",")),.runner_name,.conclusion] | @tsv'
   gh api repos/pt9912/homebrew-u-boot/check-runs/<job-id>/annotations --jq '.[] | [.annotation_level,.message] | @tsv'
   gh api --allow-escape-sequences repos/pt9912/homebrew-u-boot/actions/jobs/<job-id>/logs \
     | grep -E 'Image:|Cellar/u-boot/[0-9.]+:|u-boot version'
   ```

   Pro Matrix-Eintrag (`labels`) trägt ein Job `runner_name` und `conclusion`. Das Log (nicht allein `gh run view
   --log`) zeigt die Image-Zeile (`Image: …`) und die installierte Formel-Version (`Cellar/u-boot/<Version>`, `u-boot
   version <Version>`); „grün“ heißt diese Version und diese Image-Zeile, nicht „macOS grün“. Ein Job ohne Runner trägt
   die Annotation „The configuration '…' is not supported“ (§8).
6. **Rückweg.** Ein weiterer Commit mit dem vorigen Inhalt: derselbe Aufruf wie in Schritt 3, mit `B64` aus der
   früheren Revision der Quelle (`git show <Revision>:packaging/homebrew-tap/$P | base64 -w0`) und dem dann aktuellen
   `sha`.

**Gelaufen und nicht gelaufen.** Die lesenden Befehle der Schritte 1, 4 und 5 sind gegen den Tap ausgeführt, ebenso die
Befehle vor Schritt 3. Der Aufruf von Schritt 3 ist durch den Hauptagenten in der Sitzung vom 2026-11-06 für
`.github/workflows/smoke.yml` gelaufen (ohne `branch`, `B64` und `SHA` vorab gesetzt); sein Wortlaut stammt aus dem
Sitzungsprotokoll. Die Zeile für `B64` im Block oben ist nicht als gelaufen ausgewiesen. Für `README.md` und für den
Rückweg (Schritt 6) ist der Aufruf nicht gelaufen. Eine Probe, die nichts schreibt, gibt es für den Schreibschritt nicht;
der Dispatch aus Schritt 5 ist ebenfalls nur berichtet.

**Grenze.** Kein Skript, kein Gate und kein Test hält Quelle und Tap gleich oder schreibt diese zwei Dateien in den Tap:
`scripts/tap-nachzug.sh` klont den Tap, kopiert die Formel nach `Formula/u-boot.rb` und fügt diese Datei hinzu (Zeilen
25 bis 27 und 34); `scripts/tap-check.sh` führt `git push --dry-run` auf den geklonten Stand aus (Zeile 34) und belegt
kein Recht für `.github/workflows/`. Der Vergleich der Schritte 1 und 4 ist die einzige Prüfung; dazwischen kann jede
Seite allein geändert werden, und ein stabiler Tag startet den Smoke mit der Datei, die im Tap liegt.
`base64 -w0` (Schritt 3 und 6) ist eine GNU-Option; ausgeführt wurde sie unter Linux, über andere Plattformen sagt der
Abschnitt nichts.
