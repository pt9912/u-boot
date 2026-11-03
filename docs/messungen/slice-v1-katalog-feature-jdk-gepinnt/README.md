# Messung: Vertriebs-Lage des Katalog-Schlüssels `java-toolchain`

| Dokument    | Messprotokoll |
| ----------- | ------------- |
| Auslösender Vorgang | [`slice-v1-katalog-feature-jdk-gepinnt`](../../plan/planning/in-progress/slice-v1-katalog-feature-jdk-gepinnt.md) (Umsetzung zu Lastenheft 0.31.0, [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`SPEC-054`](../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)) |
| Gelesen am  | 2026-11-02 |
| Skripte     | `tools/messung/katalog-feature-jdk/` (M1 `temurin-checksums.sh` — Aufruf je Version, Modus `download`/`verifyx64`; M2 `java-feature-tags.sh`; M3 `brew-casks.sh`) |

## Umgebung

Linux-Host `amd64`, `curl` 8.18.0 (x86_64-pc-linux-gnu) netzgebunden; keine
Engine- und CLI-Abhängigkeit (die Tag- und Cask-Lesung lief über Registry- und
GitHub-API per `curl`; der Devcontainer-CLI war nicht beteiligt). Die Messung
liest nur; sie schreibt nichts ins Repo außer diesem Protokoll samt Rohdaten
(`rohdaten-m1.txt`, `rohdaten-m2.txt`, `rohdaten-m3.txt` — die gelesenen
Läufe, archiviert).

## M1 — Temurin-Release-Archiv und Prüfsummen-Quelle

- Die Veröffentlichung `adoptium/temurin21-binaries`, Tag `jdk-21.0.13+9`,
  veröffentlicht 2026-10-22 (GitHub-API, `releases/latest`): die JDK-Archive
  der zwei Linux-Architekturen heißen
  `OpenJDK21U-jdk_x64_linux_hotspot_21.0.13_9.tar.gz` (207.858.352 Byte) und
  `OpenJDK21U-jdk_aarch64_linux_hotspot_21.0.13_9.tar.gz` (205.995.218 Byte);
  daneben je ein Paar `…tar.gz.json` und `…tar.gz.sha256.txt` (Prüfsummen-Datei
  des Archivs) und eine `.sig`-Datei (kein Gegenstand der Prüfung von `u-boot`).
- **Die Prüfsummen-Dateien je Architektur sind Assets derselben
  Veröffentlichung** (Same-Origin am Archiv):
  - x64: `d7f5d5cf865a2902b70ff549b97586a7d3f234e016f041ab447ea740c7d66d1c`
  - aarch64: `e3522e707fc7dbe2a72c30746b1e7068301a0213ecd213d4491f35a14eb55578`
  (Zeilenform `<sha256>  <Dateiname>`.)
- **Gegenprobe (Stream):** der Digest des heruntergestreamten x64-Archivs
  entspricht der gelesenen Prüfsumme (`d7f5…d1c`, Modus `verifyx64`).
- **Vertriebs-Hosts:** der Download eines Assets heißt
  `github.com/adoptium/temurin21-binaries/releases/download/<Tag>/<Asset>` und
  ein verfolgter Redirect endet (je Asset, gemessen an den zwei
  Prüfsummen-Dateien und dem Archiv) auf
  `release-assets.githubusercontent.com`. Keine weitere Host-Lage gemessen.
- **Die Form generalisiert auf eine nutzergewählte Version:** die
  Prüfsummen-Dateien der Veröffentlichung `jdk-21.0.9+10` (2025-10-23) tragen
  `810d3773df7e0d6c4394e4e244b264c8b30e0b05a0acf542d065fd78a6b65c2f` (x64) und
  `edf0da4debe7cf475dbe320d174d6eed81479eb363f41e38a2efb740428c603a`
  (aarch64) — dieselbe Asset-Konvention (`OpenJDK21U-jdk_<slug>_linux_hotspot_`
  mit `_` für `+`), damit deckt die Quelle auch eine ältere nutzergewählte
  Version; die Docker-Tests bauen genau diesen Weg (Übersteuerung `21.0.9+10`).

## M2 — Tags und Metadaten des OCI-Features `java`

- Die Tags von `ghcr.io/devcontainers/features/java` (Registry-Listung,
  Bearer-Token): `latest`, `1` (der Render-Default), `1.0`–`1.8.3` — **kein
  Tag `21`**. Der Render-Suffix eines Feature-Eintrags heißt das OCI-Release
  des Features, nicht die Java-Version.
- Die Metadaten des Features (main, Feature-Version `1.8.3`): „Installs Java,
  SDKMAN! (if not installed), and needed dependencies."; die Option `version`
  des Features nimmt die Java-Versionen (`latest`, `none`, `21`, `17`, `11`,
  `8`) — sie wählt das JDK, das das Feature selbst installiert, und hat mit
  dem Tag-Suffix des Render-Schlüssels nichts zu tun. Die Fassung vor diesem
  Slice sagte zum Beispiel `java.version 21`, der Eintrag hieße dann
  `ghcr.io/devcontainers/features/java:21` — der Eintrag hieße so, aber ein
  solcher Tag existiert nicht (der Feature-Build bricht ab); berichtigt auf
  `1.8.3` samt Auslegungs-Zeile zum Schlüssel `java-toolchain`. **Kein
  separater Abweisungsfall `sdkman`** — der Konflikt des Typs gilt der zwei
  JDKs (externes `java` plus Schlüssel).

## M3 — Brew-Form des Vertriebs (Auflage)

- **Cask `temurin`** (tap `homebrew/cask`): „Eclipse Temurin Java Development
  Kit", Version `27.0.1.0.2,1`; die Installation lädt die `.pkg`-Form der
  mac-Veröffentlichung `temurin27-binaries`.
- **Cask `temurin@21`** (tap `homebrew/cask`): „Eclipse Temurin 21", Version
  `21.0.13,9`,Installation über die `.pkg`-Form **derselben Veröffentlichung**
  (`jdk-21.0.13+9`), `aarch64_mac`.
- **Befund:** brew verteilt die macOS-Installer-Form (`pkg`) derselben
  Veröffentlichungen, keine Linux-Archivform. Im Devcontainer-Kontext (Linux)
  installiert `brew` Casks nicht; der Katalog-Weg (Linux-`tar.gz`-Archiv)
  bleibt die einzige Installationsform im Devcontainer. **Der Befund begründet
  keine zweite Zusage** — er steht als Vertriebs-Vergleich (Chronik,
  `CHANGELOG`-Eintrag); die Mess-Pflicht der Umsetzungs-Form bleibt.

## Grenzen

- Der Stream-Vergleich lief für das x64-Archiv; das aarch64-Archiv wurde nicht
  gestreamt (Prüfsumme gelesen, nicht verifiziert).
- Der Stand der Tags und Metadaten ist der vom 2026-11-02; Änderungen des
  Vertriebs (neue Tags, neue Casks) altert das Protokoll, nicht die Zusage.
- Der Default-Pin des Templates trägt die Revision `21.0.13+9`; eine spätere
  Hebung folgt einem u-boot-Release (Muster der Pin-Politik).