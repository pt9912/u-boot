# ADR 0017: Linux-Pakete (.deb/.rpm) per nfpm als Release-Assets

**Status:** Accepted

**Datum:** 2026-10-01

**Autor:** pt9912

**Bezug:** [`LH-OPEN-002`](../../../spec/lastenheft.md#lh-open-002--paketierung), [ADR-0007](0007-distributionswege-ghcr.md), [ADR-0016](0016-homebrew-distribution-per-tap.md)

**Schärft:** [`LH-OPEN-002`](../../../spec/lastenheft.md#lh-open-002--paketierung) — Debian/RPM ist von „vertagt mit Trigger“ auf „gewählt“ gesetzt; Werkzeug und Verteilungsform.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0007](0007-distributionswege-ghcr.md) vertagte Debian/RPM wegen des Tooling-Overheads (`debhelper`, `rpmbuild`, Repository-Hosting). Die statisch gelinkten Linux-Binaries (amd64, arm64) hängen seit v0.1.1 an jedem Tag. Der Projektinhaber will installierbare Pakete, ohne eine eigene Paket-Infrastruktur zu betreiben. Der Accepted-ADR-0007 bleibt unverändert; diese Folgeentscheidung ersetzt die Zeile „Debian/RPM: vertagt“.

## Entscheidung

1. **Ein Werkzeug für beide Formate:** `nfpm` (`goreleaser/nfpm`, Image-Pin in `scripts/build-packages.sh`) erzeugt `.deb` und `.rpm` aus einer YAML-Spec (`packaging/nfpm.yaml`). Das Paket enthält nur das fertige Binary (`/usr/bin/u-boot`) und die Lizenz; keine Maintainer-Skripte.
2. **Vier Pakete pro Tag:** amd64 und arm64, je `.deb` und `.rpm`, gebaut aus denselben Binaries wie die Binary-Assets (`make packages` lokal, `publish.yml` pro `v*`-Tag).
3. **Nur Release-Assets, kein Repository:** Die Pakete hängen am GitHub-Release und stehen in `SHA256SUMS`; Installation per `apt install ./u-boot_<v>_amd64.deb` bzw. `dnf install ./u-boot-<v>-1.x86_64.rpm`. Kein APT-/DNF-Repository, keine Signierung im ersten Schritt.
4. **Smoke nach dem Publish:** Zwei Jobs installieren die veröffentlichten Pakete auf Ubuntu (`.deb`) und in einem Fedora-Container (`.rpm`) und prüfen `--version` sowie `init`.
5. **Nicht Teil dieser Entscheidung:** Aufnahme in offizielle Distro-Repositories, gehostetes APT-/DNF-Repository, Paket-Signierung, Snap/Flatpak/AppImage.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — `debhelper` + `rpmbuild` nativ | volle Kontrolle, Distro-Konformität | zwei Toolchains, hoher Pflegeaufwand für ein statisches Einzel-Binary |
| **B — `nfpm` als Single Source** | eine YAML-Spec, Docker-only lauffähig, passt zum Einzel-Binary | weniger Distro-Policy-Prüfung (lintian/rpmlint nicht Teil des Gates) |
| C — gehostetes APT-/DNF-Repository (PPA, OBS, Cloudsmith) | `apt install u-boot` ohne Download | externe Infrastruktur, Schlüsselverwaltung, laufende Kosten |
| D — nichts tun | kein Aufwand | Linux-Nutzer ohne Paketweg |

## Konsequenzen

- Positiv: Ein Tag erzeugt Binaries, Pakete und Prüfsummen in einem Vorgang; lokal reproduzierbar per `make packages`.
- Negativ: Kein `apt upgrade`; Updates erfolgen durch erneute Installation des neuen Pakets. Pakete sind unsigniert (Integrität über `SHA256SUMS`).
- Folgepflicht: Pakete auf dem ersten Tag nach Einführung auf echten Systemen verifizieren; ein gehostetes Repository bleibt ein möglicher Folgeschritt.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Release-Workflow | Smoke-Jobs installieren `.deb` (Ubuntu) und `.rpm` (Fedora) aus dem Release und prüfen `--version` | `publish.yml` (nur auf Tag) |

## Re-Evaluierungs-Trigger

Nachfrage nach `apt upgrade`/`dnf upgrade` (gehostetes Repository), Signierungsanforderung oder Aufnahme in eine Distribution.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-01 | Entschieden und `Accepted` | Vereinbarung mit dem Projektinhaber |
