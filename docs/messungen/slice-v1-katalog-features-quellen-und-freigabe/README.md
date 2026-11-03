# Messung: Registry-Lage der Katalogquellen `cpp` und `postgres-client` (N1)

| Dokument    | Messprotokoll |
| ----------- | ------------- |
| Auslösender Vorgang | [`slice-v1-katalog-features-quellen-und-freigabe`](../../plan/planning/in-progress/slice-v1-katalog-features-quellen-und-freigabe.md) (Nebenbefunde N1/N2 aus [`ADR-0032`](../../plan/adr/0032-katalog-features-socket-und-lockerungen.md)), [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) |
| Gelesen am  | 2026-11-03 |
| Skripte     | `tools/messung/katalog-quellen/registry-quellen.sh` (Aufruf je Quelle, Umleitung in die Rohdaten) |

## Umgebung

Linux-Host `amd64`, `curl` 8.18.0 netzgebunden; keine Engine- und CLI-Abhängigkeit
(Die Registry- und GitHub-API-Lesung lief per `curl` (M1); die Schlüsselform-Probe (M2) lief mit der CLI 0.89.0 im Helper-Container mit gemountetem Docker-Socket — siehe Abschnitt M2). Die
Messung liest nur; sie schreibt nichts ins Repo außer diesem Protokoll samt
Rohdaten.

## M1 — Registry-Lage der zwei Katalogquellen (gemessen 2026-11-03)

- **Upstream-Bestand:** das Quell-Repository `devcontainers/features` (GitHub-API,
  `contents/src?per_page=100`) trägt **29 Features** (Liste in den Rohdaten) —
  **weder `cpp` noch `postgres*`** ist im Bestand.
- **Registry:** `ghcr.io/devcontainers/features/cpp` und
  `ghcr.io/devcontainers/features/postgresql-client`: der anonyme
  Pull-Token-Abruf endet mit `403 DENIED invalid token` und die Tag-Listung
  entsprechend — **ghcr stellt für beide Repositorys keinen Token aus**. Das
  Muster des Nicht-Existierens (Annahme der [`ADR-0032`](../../plan/adr/0032-katalog-features-socket-und-lockerungen.md)
  *Nicht Gegenstand* (a) „NAME_UNKNOWN“) ist damit am lebenden Registry-Verhalten
  als **Token-Verweigerung** bestätigt, nicht als `NAME_UNKNOWN`-Antwort des
  Manifest-Zugriffs gelesen.
- **Kontrollproben:** `devcontainers/features/go` und `devcontainers/features/java`
  enden je `200` mit Tag-Liste — die Lesung und die Token-Infrastruktur tragen
  die Messung; die beiden Fehlverläufe sind messbar den nicht existierenden
  beiden Katalogquellen zuzuordnen, nicht der Messmethode.
- **Fehlwahl der Erst-Controlle:** `typescript` existiert im Upstream ebenfalls
  nicht (404 in beiden Lesungen) — die Erst-Läufe mit dieser Controlle tragen
  das gleiche Muster und sind als Steuerwahl-Fehler notiert, nicht als
  Befund-Kontrast.

## M2 — Schlüsselform-Probe der Dev-Container-CLI (gemessen 2026-11-03)

Fixture-Workspace (`busybox:latest`, Features je Form), CLI 0.89.0 (Image
`u-boot-devcontainer-cli:0.89.0-f1a94cce795e`, docker.sock-Mount), Befehl
`devcontainer up`, Formen mit existierendem Tag `go:1.0.3`
(`rohdaten-m2.txt`; Erst-Läufe mit Nicht-Tag `9.9.9` diskrimierten nicht,
`rohdaten-m2-erstrun.txt`):

| Form | Ergebnis |
|---|---|
| `ghcr.io/devcontainers/features/go:1.0.3` (Controlle) | Feature-Metadaten und Dependencies gelesen („Feature : Go"); der weitere Bau scheiterte am `docker buildx build`-Zug der Probe-Umgebung (rohdaten-m2-f1-voll.log) — keine Form-Aussage, Auflösung trägt |
| `oci://ghcr.io/devcontainers/features/go:1.0.3` | scheitert mit eigener Fehlerklasse („Feature manifest for 'oci://…' … Ensure the Feature is published and accessible") — **unverarbeitet** |
| `https://ghcr.io/devcontainers/features/go:1.0.3` | scheitert am Identifier-Verarbeiten („could not be processed") — **unverarbeitet** |

Das Muster des [ADR-0032](../../plan/adr/0032-katalog-features-socket-und-lockerungen.md)-Vermerks (b) ist für die zwei Form-Schreibweisen
bestätigt; u-boot rendert `source:version` wörtlich — externe Quellen in
dieser Form bauen nicht.

## Grenzen

- Der `https://…/devcontainer-feature-<id>.tgz`-Tarball-Weg ([ADR-0032](../../plan/adr/0032-katalog-features-socket-und-lockerungen.md) (b)) wurde
  **nicht gemessen** — nur die bloßen `https://…:1`-Schlüsselformen; eine
  Tarball-servierte Quelle bleibt offene Probe.
- Die Manifest-Zugriffsform (`/v2/.../manifests/<tag>`) wurde nicht probeziert
 — der Token-Ausgang reicht fürs Urteil der Registry-Lage; die Antwortform
 `NAME_UNKNOWN` (Annahme der ADR) bleibt so nicht wörtlich gemessen, sondern
 als Token-Verweigerung belegt.
- Die Schlüsselform-Probe (N2: ob die Dev-Container-CLI `0.89.0` einen Schlüssel
 `https://…:1` oder `oci://…:1` auflöst) ist als M2 geführt (Abschnitt M2); die
 Tarball-Variante des `https://`-Weges blieb ungemessen und bleibt offene
 Grenze — deshalb trägt der [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)-Wortlaut
 „Registry-Tarball-Form offen“; keine Aussage über die Verarbeitbarkeit des
 Tarball-Weges steht in diesem Protokoll.
- Der Stand gilt am 2026-11-03; künftige Upstream-Änderungen alter das
 Protokoll, nicht die Mess-Basis. Die Messung liest nur; sie schreibt nichts.