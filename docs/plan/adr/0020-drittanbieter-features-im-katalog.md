# ADR 0020: Drittanbieter-Feature `homebrew` und Alias `sdkman` im eingebauten Feature-Katalog

**Status:** Accepted

**Datum:** 2026-10-15

**Autor:** pt9912

**Bezug:** [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [ADR-0012](0012-devcontainer-egress-firewall.md)

**Schärft:** [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) — Herkunft und Kennzeichnung der eingebauten Features; der Katalog umfasst bisher nur das offizielle Namensraum-Präfix `ghcr.io/devcontainers/features/`.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der eingebaute Katalog kennt acht Schlüssel (`git`, `docker-cli`, `node`, `java`, `go`, `cpp`, `kubectl-helm`, `postgres-client`), alle aus dem offiziellen Namensraum `ghcr.io/devcontainers/features/`. [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) nennt SDKMAN als Beispiel für ein Feature; für Homebrew gibt es keine Zusage. Gewünscht sind beide als eingebaute Schlüssel.

**Messung (2026-10-15).** Dev-Container-CLI 0.89.0, Basis `mcr.microsoft.com/devcontainers/base:debian` (Colima 0.10.3, arm64):

| Prüfung | Ergebnis |
|---|---|
| Offizielles Feature `homebrew` oder `sdkman` unter `ghcr.io/devcontainers/features/` | gibt es nicht (Index containers.dev/features) |
| Community-Feature `ghcr.io/meaningful-ooo/devcontainer-features/homebrew:2` (2.0.5) | baut; `brew --version` → Homebrew 7.0.9; `brew install jq` läuft |
| Community-Feature `ghcr.io/ebaskoro/devcontainer-features/sdkman:1` (1.0.0) ohne Optionen | Build bricht ab („Stop! latest is not a valid candidate“): das Feature installiert einen SDK-Kandidaten und braucht `candidate` und `version`, nicht SDKMAN selbst |
| Offizielles Feature `ghcr.io/devcontainers/features/java:1` | installiert SDKMAN mit Java; `sdk list java`, `sdk install maven`, `sdk install gradle` laufen |
| DNS-Abfragen von `brew install jq` | `formulae.brew.sh`, `ghcr.io`, `pkg-containers.githubusercontent.com`, `github.com`, `api.github.com` |
| DNS-Abfragen von `sdk list java`, `sdk install maven`, `sdk install gradle` | `api.sdkman.io`, `broker.sdkman.io`; je Kandidat die Hersteller-Hosts (`repo.maven.apache.org`, `services.gradle.org`, `github.com`, `release-assets.githubusercontent.com`) |

## Entscheidung

1. **`homebrew` kommt als eingebauter Schlüssel in den Katalog** mit dem Community-Feature `ghcr.io/meaningful-ooo/devcontainer-features/homebrew`, Standardversion `2` (Major-Pin). Der Katalogeintrag trägt den Herausgeber (`meaningful-ooo`); Katalogeinträge außerhalb des offiziellen Namensraums brauchen diesen Herausgeber als Registry-Namensraum.
2. **Kennzeichnung ([`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)):** `u-boot generate devcontainer` weist jedes aktivierte Drittanbieter-Feature des Katalogs als Warnung aus (Code [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)): Herausgeber, Quelle und Version, mit dem Hinweis, dass dessen Installationsskript beim Image-Build läuft. Eine ausdrückliche `source:`-Überschreibung (mit Freigabe nach [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)) bleibt der externe Weg ohne diese Warnung, weil die Freigabe die Zustimmung ist.
3. **`sdkman` ist ein Alias auf das offizielle Feature `ghcr.io/devcontainers/features/java:1`** (Standardversion `1`), das SDKMAN installiert. Keine Drittanbieter-Quelle; keine Feature-Optionen. Das Community-Feature `ebaskoro/sdkman` ist nicht im Katalog, weil es keinen SDKMAN-Installer liefert, sondern einen Kandidaten und Optionen verlangt, die der Katalog nicht kennt.
4. **Gleiche Render-Schlüssel erscheinen einmal:** `java` und `sdkman` zusammen aktiviert rendern `ghcr.io/devcontainers/features/java:1` einmal in `devcontainer.json`.
5. **Egress ([ADR-0012](0012-devcontainer-egress-firewall.md)):** Die Default-Allowlist enthält je Feature: `homebrew` → `formulae.brew.sh`, `ghcr.io`, `pkg-containers.githubusercontent.com`; `sdkman` → `api.sdkman.io`, `broker.sdkman.io`, `repo.maven.apache.org`, `repo1.maven.org`. Die Download-Hosts weiterer Hersteller-Kandidaten (etwa `services.gradle.org`) ergänzt der Nutzer über `devcontainer.sandbox.egress.allow`; die Doku nennt das.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nur offizielle Features im Katalog (Status quo) | keine Drittanbieter-Skripte im Katalog | Homebrew und ein SDKMAN-Schlüssel fehlen |
| **B — `homebrew` als Community-Feature im Katalog (gekennzeichnet), `sdkman` als Alias auf `java`** | beide Schlüssel nutzbar; Drittanbieter sichtbar ausgewiesen; SDKMAN aus offizieller Quelle | Drittanbieter-Feature im Katalog; Pflege der Egress-Hosts |
| C — `sdkman` aus `ebaskoro/sdkman` mit Optionen | SDKMAN ohne Java-Zwang | Optionsmodell für den ganzen Katalog nötig; Drittanbieter-Feature; Default-Kandidat willkürlich |
| D — Nur Doku (`featureSources.allow` plus Quelle) | keine Spec-Änderung | beide Schlüssel nicht eingebaut |

## Konsequenzen

- Positiv: `u-boot config set devcontainer.features.homebrew.enabled true` und `...sdkman.enabled true` genügen; die Herkunft von `homebrew` ist in der Befehlsausgabe sichtbar.
- Negativ: Ein Drittanbieter-Feature im Katalog hängt an dessen Pflege; ändert oder entfernt der Herausgeber `homebrew:2`, bricht der Image-Build. `sdkman` installiert auch Java (Verhalten des `java`-Features).
- Folgepflicht: Golden Cases für Render, Warnung, Dedupe und Egress-Hosts; die Dev-Container-CLI wird im Repo nicht ausgeführt, die Lauffähigkeit der Features ist durch die Messung dieser ADR belegt, nicht durch einen Test.
- **Nicht verifiziert (Stand 2026-10-15):** Lauf unter VS Code und Codespaces; `brew install` mit anderen Formeln als `jq` und weitere SDKMAN-Kandidaten (Hosts je Hersteller).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Katalog-Test | jeder Eintrag ist offiziell (Präfix) oder nennt seinen Herausgeber als Namensraum | `make test` |
| Golden Case | `homebrew` rendert `…/homebrew:2` mit Warnung; `sdkman` rendert `java:1`; `java` und `sdkman` zusammen einmal | `make test` |

## Re-Evaluierungs-Trigger

Ein offizielles `homebrew`- oder `sdkman`-Feature unter `ghcr.io/devcontainers/features/` (dann dorthin umstellen), ein Optionsmodell für Katalog-Features, Änderung oder Wegfall von `meaningful-ooo/homebrew:2`, oder ein Nutzerbericht, dass die Default-Hosts für `brew install`/`sdk install` nicht reichen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-15 | Entwurf mit Messreihe (Dev-Container-CLI 0.89.0, Colima, arm64); `Accepted` durch den Projektinhaber | Vereinbarung mit dem Projektinhaber |
