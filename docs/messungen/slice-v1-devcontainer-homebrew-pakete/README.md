# Messung: Homebrew-Formeln im Devcontainer (Katalog-Schlüssel `brew`)

| Dokument    | Messprotokoll |
| ----------- | ------------- |
| Auslösender Vorgang | [`slice-v1-devcontainer-homebrew-pakete`](../../plan/planning/in-progress/slice-v1-devcontainer-homebrew-pakete.md) (Umsetzung zu Lastenheft 0.36.1, [`LH-FA-DEV-013`](../../../spec/lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer), [`SPEC-057`](../../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen)); Messpunkte aus [`ADR-0046`](../../plan/adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md) |
| Gelesen am  | 2026-11-05 |
| Skripte     | `tools/messung/devcontainer-homebrew-pakete/` (`01-bau.sh`, `02-pfad.sh`, `03-aktualisierung.sh`, `04-egress.sh` mit `inner-04.sh`, `06-arm64.sh`, `05-aufraeumen.sh`, gemeinsam `lib.sh`; `00-alles.sh` ruft 01 bis 04 und 06) |

## Umgebung

Linux-Host `amd64` (`x86_64`), Docker Engine 29.9.3 (`linux/amd64`, die verschachtelte Engine des Entwicklungs-Devcontainers dieses Repos),
Basis-Image `mcr.microsoft.com/devcontainers/base:debian` (Debian 13, 1.412.272.019 Byte lokal), Netz offen (kein Proxy). Das Binary `u-boot`
(`u-boot version 0.23.0-dev`, der Stand des Slice) erzeugt die Projekte: `init --devcontainer`, die Einstellungen per `config set` in der
Reihenfolge der Skripte, `generate devcontainer`. Ein Lauf je Messpunkt; die Rohdaten sind die archivierten Läufe
(`rohdaten-m1.txt` bis `rohdaten-m5.txt`, M4 je Modus eine Datei). Die Messung liest keine Zugangsdaten und setzt keine; sie schreibt nichts ins
Repo außer diesem Protokoll samt Rohdaten.

## Pin-Revision (Quelle)

`git ls-remote --tags https://github.com/Homebrew/brew` am 2026-11-05: der höchste Release-Tag ist `7.0.12`, sein Hash
`bf464f57275cae0be106b43d190c98564ac59b0e`. Die Tags sind leichtgewichtig (`git ls-remote … refs/tags/7.0.12 'refs/tags/7.0.12^{}'` liefert
nur die eine Zeile), der Hash ist also der Commit. Dieser Hash steht als `ARG HOMEBREW_REVISION` im Template. Der Commit des Releases davor
(`7.0.10`) ist `8f10429b56be28a25032b87463fffc95e556b5bc`; er dient M3 als älterer Stand.

## M1 — Bau (`rohdaten-m1.txt`)

Formeln `ripgrep` und `pt9912/u-boot/u-boot` (die Formel des Taps `pt9912/homebrew-u-boot`), Tap freigegeben, alle Baue mit `--no-cache`:

- **Standardprofil:** Bau gelingt, 35 s (kalt, ein Lauf); Image 1.704.035.546 Byte, das sind 291.763.527 Byte (rund 292 MB) mehr als das Basis-Image.
- **Sandbox-Profil mit eingeschalteter Egress-Restriktion** (ohne Eintrag für Homebrew in `egress.allow`; die Zeile `brew` der Liste steht im
  Startscript, der Bau liegt nicht hinter der Restriktion): Bau gelingt, 41 s; Image 1.710.992.029 Byte.
- **Homebrew allein** (`enabled: true`, leere Liste): Bau gelingt, 11 s; Image 1.472.764.145 Byte, das sind 60.492.126 Byte (rund 60 MB) mehr
  als das Basis-Image.
- Im Image (Standard- und Sandbox-Profil gleich): Benutzer `vscode` (UID 1000), `brew` unter `/home/linuxbrew/.linuxbrew/bin/brew`, die Revision
  des Homebrew-Verzeichnisses ist der Pin, Eigentümer `vscode:vscode`, installiert sind `ripgrep` (15.2.0), `u-boot` (0.23.0, aus dem Tap) und
  die Abhängigkeiten `bzip2`, `pcre2`, `zlib-ng-compat`; das Verzeichnis ist ein flacher Klon (`git rev-parse --is-shallow-repository`: `true`),
  `brew --version` meldet deshalb `Homebrew >=4.3.0 (shallow or no git repository)`; es gibt keine Umgebungsvariable `HOMEBREW*` im Image.
- **Negativprobe:** ein Pin, den es nicht gibt (`--build-arg HOMEBREW_REVISION=000…001`), bricht den Bau beim `git fetch` ab (Exit-Code 128 des
  Schritts, kein Image).
- Ein Bau braucht keinen Eintrag in einer Egress-Liste: Die Restriktion greift beim Containerstart, nicht beim Bau (die Sandbox-Variante ist mit
  eingeschalteter Restriktion gebaut).
- Die Versionen der Formeln und die Größen sind Momentaufnahmen des 2026-11-05; die Formeln sind nicht gepinnt.

## M2 — `PATH` und Sitzungsarten (`rohdaten-m2.txt`)

Image der Standardprofil-Variante, Benutzer `vscode`. Das Image setzt `ENV PATH` mit `/home/linuxbrew/.linuxbrew/bin` und `…/sbin` vor dem
bisherigen Wert; `brew` wird gefunden bei `bash -c`, `bash -ic`, `sh -c`, `zsh -c`, `zsh -lc`, `zsh -ic`, `zsh -lic` und in der Sitzung von
`docker exec` (mit und ohne `-u vscode`). `bash -lc`, `bash -lic` und `sh -lc` finden `brew` **nicht**: `/etc/profile` setzt den `PATH`
neu (Zeilen 5 und 7), `/etc/profile.d/00-restore-env.sh` stellt ihn nicht her. Die Shell von `u-boot devcontainer shell` ist keine Login-Shell
([`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)),
der Block setzt keine Datei unter `/etc/profile.d/`.

## M3 — Selbst-Aktualisierung (`rohdaten-m3.txt`)

Image auf den Commit des Releases `7.0.10` gebaut (`--build-arg`). Das Alter von `.git/FETCH_HEAD` wird auf drei Tage gesetzt (die Automatik
prüft es), dann `brew install --formula jq`:

- **mit `HOMEBREW_NO_AUTO_UPDATE=1 HOMEBREW_NO_ANALYTICS=1`** (die Umgebung des Blocks): die Revision bleibt `8f10429b…`.
- **ohne den Schalter** (nur `HOMEBREW_NO_ANALYTICS=1`): `==> Auto-updating Homebrew…`, danach steht der Klon auf `bf464f57…` (dem Release `7.0.12`),
  Branch `stable`.

Folgerung aus dem Befund, kein gefahrener Lauf: Ein Bau, dessen Checkout-Schicht aus dem Zwischenspeicher stammt und älter ist als der
Zeitraum der Automatik (zum Beispiel nach einer Änderung der Formel-Liste), würde den Pin verschieben; der Schalter im Befehl der Installation
verhindert es. Bei jungem `FETCH_HEAD` (ein frischer Bau) greift die Automatik nicht; ein Lauf ohne Schalter bei jungem `FETCH_HEAD` ist nicht
gemessen.

## M4 — Laufzeit hinter der Egress-Restriktion (`rohdaten-m4-egr-*.txt`)

Image der Variante `egr` (Sandbox-Profil, Restriktion an, Homebrew allein), Container mit `--cap-add=NET_ADMIN`, ohne Volume; im Container
`brew install --formula jq` (eine Formel von Homebrew selbst) und `brew install --formula pt9912/u-boot/u-boot` (Formel des Taps), jeweils mit der
Umgebung eines Nutzers (Automatik der Aktualisierung und Analytik von Homebrew an):

| Modus | Liste im Startscript | `jq` | Tap-Formel |
|---|---|---|---|
| `discover` (ohne Restriktion) | — | gelingt; Namen: `formulae.brew.sh`, `ghcr.io`, `pkg-containers.githubusercontent.com` | gelingt; Namen: `api.github.com`, `github.com`, `release-assets.githubusercontent.com` |
| `restricted-ohne` (die drei Hosts der Zeile `brew` aus der Liste entfernt) | Basisliste | **scheitert** schon beim ersten Start von Homebrew: `Failed to install Homebrew Portable Ruby` (`ghcr.io` wird nicht aufgelöst) | scheitert gleich |
| `restricted` (das erzeugte Startscript) | Basisliste plus `formulae.brew.sh`, `ghcr.io`, `pkg-containers.githubusercontent.com` | gelingt | **scheitert**: `Could not resolve host: release-assets.githubusercontent.com` |
| `restricted-formelhost` (erzeugte Liste plus `release-assets.githubusercontent.com`) | wie `restricted` plus der Host der Formel | gelingt | gelingt |

Befund: Die drei Hosts sind die Hosts von Homebrew selbst (Portable Ruby, Metadaten-API, Bottles); ohne sie scheitert jedes `brew install`. Der Host,
von dem die gemessene Tap-Formel lädt, ist Sache der Formel (`release-assets.githubusercontent.com`, nicht in der Basisliste) und steht nicht in
der Zeile; der Nutzer trägt ihn in `devcontainer.sandbox.egress.allow` ein. Die Tap-Formel `pt9912/u-boot/u-boot` wird bei der Installation
getappt und ihr Vertrauen gesetzt (`==> Tapping pt9912/u-boot`, `==> Trusted formula pt9912/u-boot/u-boot`): Eine voll qualifizierte Formel
braucht weder `brew tap` noch `brew trust` im Block (Homebrew: `docs/Tap-Trust.md`, „Installing a fully qualified formula … trusts only that
item“); das Vertrauen gilt der Formel, nicht dem Tap.

## M5 — Plattform (`rohdaten-m5.txt`)

`docker run --platform linux/arm64 mcr.microsoft.com/devcontainers/base:debian uname -m` endet mit `exec format error` (Exit-Code 255): Auf diesem
Host steht keine Emulation für `arm64` bereit. **`arm64` ist nicht gemessen.**

## Aufruf

```bash
# Binary des Slice-Stands (im Devcontainer-Container):
make compile BAU_COMPILE_OUT=<pfad>/u-boot
UBOOT=<pfad>/u-boot WORK=<messordner-ausserhalb-des-repos> tools/messung/devcontainer-homebrew-pakete/00-alles.sh
UBOOT=<pfad>/u-boot WORK=<messordner> tools/messung/devcontainer-homebrew-pakete/05-aufraeumen.sh
```

Lesehilfe: Zeilen `AUSSAGE …` sind die stabilen Aussagen, die ein zweiter Lauf wiederholen muss; Dauern und Größen streuen.

## Grenzen

- Ein Lauf je Messpunkt, ein Host, die verschachtelte Engine; Dauern streuen (frühere, nicht archivierte Läufe lagen im selben Bereich).
- `arm64` ist nicht gemessen (M5); die Zusage gilt nur für `amd64`.
- Gemessen sind `ripgrep`, `jq`, `tree` (im Docker-Test) und eine Formel eines Taps; andere Formeln können andere Hosts, Build-Abhängigkeiten
  oder Bottles verlangen. Der Stand der Formeln und des Tap ist der vom 2026-11-05.
- Abhängigkeiten einer Tap-Formel aus anderen Taps, die Kombination mit den übrigen Dockerfile-Blöcken (`ollama`, `gh`, `glab`,
  Werkzeugketten, `nestedRuntime`) und das Verhalten nach mehreren Tagen Betrieb (Aktualisierung der Formeln im laufenden Container) sind nicht
  gemessen.
- Der Pin von Homebrew altert ohne rotes Gate; eine Hebung folgt einem Release von `u-boot`.
