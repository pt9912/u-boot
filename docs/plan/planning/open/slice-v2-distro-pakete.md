# Slice V2: Debian-/RPM-Pakete für u-boot ([`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung)-Restweg)

> **Status:** on hold — Trigger noch nicht gefeuert. Plan-Stub
> existiert, damit [`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005--carveout-disziplin)-Disziplin den Carveout-
> Anker erfüllt (siehe
> [`carveouts.md`](../in-progress/carveouts.md) §Temporäre
> Carveouts, [ADR-0007](../../adr/0007-distributionswege-ghcr.md)
> §Entscheidung Tabelle „Debian/RPM").

## Auslöser

`spec/lastenheft.md` §14 listet Debian- und RPM-Pakete als
mögliche Distributionsoptionen für [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung). [ADR-0007](../../adr/0007-distributionswege-ghcr.md) hat
beide gemeinsam vertagt — der Tooling-Overhead (`debhelper`,
`rpmbuild`, Repository-Hosting) ist hoch und ohne konkrete
Distro-Anfrage nicht zu rechtfertigen.

Voraussetzung ist erfüllt: das Linux-amd64- und Linux-arm64-Binary
existiert seit
[`slice-v2-binary-distribution`](../done/slice-v2-binary-distribution.md)
T1 (`dc9a336`/`f3f1731`) und wird pro `v*`-Tag als GitHub-Release-
Asset gepublished. Debian- und RPM-Pakete würden diese Binaries
in `.deb`- bzw. `.rpm`-Container wrappen, plus Postinst-/
Postuninst-Scripts für die `/usr/local/bin/u-boot`-Platzierung
und Bash-Completion-Installation.

## Trigger

**Konkrete Distro-Anfrage** für Debian/Ubuntu (`.deb` via
`debhelper`) oder Fedora/RHEL/openSUSE (`.rpm` via `rpmbuild`).
Ein Trigger reicht — die Implementierung kann optional von Anfang
an beide Pakete-Formate liefern, wenn der Overhead-Schmerz pro
Format gleich groß ist.

## Aufhebungsbedingung

`sudo apt install u-boot` (nach `apt repository add`) bzw.
`sudo dnf install u-boot` (nach `dnf config-manager add`)
installiert das aktuelle Tag-Binary auf einem frisch aufgesetzten
Debian/Ubuntu- bzw. Fedora/RHEL-System; `u-boot --version` zeigt
die korrekte Version; `u-boot doctor` läuft ohne Errors.

## Akzeptanzkriterien

- ✅ `.deb`-Paket-Spec via `debhelper` (oder `nfpm` als
  Light-Wrapper) — produziert ein installierbares Paket pro
  Linux-amd64 und Linux-arm64.
- ✅ `.rpm`-Paket-Spec via `rpmbuild` (oder `nfpm`) — analog
  zwei Pakete.
- ✅ Repository-Hosting-Entscheidung: entweder PPA / OBS /
  packagecloud / Cloudsmith **oder** GitHub-Release-Assets
  (`.deb` / `.rpm` direkt zum Tag) ohne dedizierten Repo-Aufbau.
- ✅ `publish.yml` (oder ein zusätzlicher Workflow) baut beide
  Paket-Formate pro Tag-Push automatisch.
- ✅ README-Install-Block (EN + DE) listet `apt install`-bzw.
  `dnf install`-Snippet.
- ✅ Smoke-Workflow auf `ubuntu-latest` + `fedora`-Container, der
  nach Paket-Build `apt install ./u-boot_*.deb` bzw.
  `rpm -i u-boot-*.rpm` ausführt + drei `LH-AK-*`-Pre-Checks
  durchläuft.

## Tranchen (vorgeschlagen, wird beim Trigger ausgearbeitet)

| T | Inhalt (Skizze) |
| - | --------------- |
| T1 | **Tooling-Decision:** `nfpm` als Single-Source-of-Truth für beide Formate (einfacher YAML-Spec, gut für Solo-Projekte) **oder** native `debhelper`/`rpmbuild`-Specs (mehr Kontrolle, mehr Lernkurve). Decision dokumentiert im Slice-Plan. Erste `.deb` + `.rpm` lokal gebaut, `dpkg -i` / `rpm -i` validiert. |
| T2 | CI-Automatisierung: `publish.yml` baut beide Pakete pro Tag-Push, hängt sie als zusätzliche GitHub-Release-Assets an (oder pusht in Repository-Hosting-Variante). |
| T3 | Smoke-Workflow für beide Pakete (ubuntu-latest + fedora-container), READMEs (EN + DE) ergänzt um `apt install` / `dnf install`-Snippets. |
| T4 | Closure: CHANGELOG, carveouts.md [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung)-Zeile entweder gelöscht (alle Restwege beliefert) oder reduziert, [ADR-0007](../../adr/0007-distributionswege-ghcr.md) §Entscheidung „Vertagt → Gewählt" für Debian/RPM. Slice-Plan `open/` → `done/`. |

## Out of Scope

- **Aufnahme in offizielle Distro-Repositories** (Debian Sid,
  Fedora Rawhide): das ist ein eigener Antragsprozess pro Distro
  mit Paket-Maintainer-Rolle und Qualitätsregeln. Erst sinnvoll
  wenn die eigene Repository-Hosting-Variante stabil läuft und es
  eine nicht-triviale Nutzerbasis gibt.
- **Snap / Flatpak / AppImage**: nicht in [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung) §14
  gelistet. Separater Slice falls Nachfrage entsteht.

## Bezug

- Spec: [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung) §14 (offene Distributionswege).
- ADR: [ADR-0007 §Entscheidung Tabelle „Debian/RPM"](../../adr/0007-distributionswege-ghcr.md)
  — verbindlicher Plan-Anker bis Trigger.
- Voraussetzungs-Slice:
  [`slice-v2-binary-distribution`](../done/slice-v2-binary-distribution.md)
  — Linux-Binaries existieren seit T1 `dc9a336`/`f3f1731` als
  GitHub-Release-Asset.
- Carveout:
  [`carveouts.md`](../in-progress/carveouts.md) §Temporäre
  Carveouts, [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung)-Zeile.
- Roadmap:
  [`roadmap.md`](../in-progress/roadmap.md) §v0.4.0+ Backlog.
- Phase: V2 (nach v0.3.0-Milestone, Trigger-getrieben).
- Geschwister-Slice:
  [`slice-v2-homebrew-formula.md`](slice-v2-homebrew-formula.md) (parallel auf hold; macOS-
  Pendant zu diesem Linux-Paket-Slice).
