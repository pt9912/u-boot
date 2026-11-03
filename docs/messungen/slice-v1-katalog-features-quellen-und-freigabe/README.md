# Messung: Registry-Lage der Katalogquellen `cpp` und `postgres-client` (N1)

| Dokument    | Messprotokoll |
| ----------- | ------------- |
| Auslösender Vorgang | [`slice-v1-katalog-features-quellen-und-freigabe`](../../plan/planning/in-progress/slice-v1-katalog-features-quellen-und-freigabe.md) (Nebenbefunde N1/N2 aus [`ADR-0032`](../../plan/adr/0032-katalog-features-socket-und-lockerungen.md)), [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) |
| Gelesen am  | 2026-11-03 |
| Skripte     | `tools/messung/katalog-quellen/registry-quellen.sh` (Aufruf je Quelle, Umleitung in die Rohdaten) |

## Umgebung

Linux-Host `amd64`, `curl` 8.18.0 netzgebunden; keine Engine- und CLI-Abhängigkeit
(Die Registry- und GitHub-API-Lesung lief per `curl`; die
Dev-Container-CLI trat nicht auf — die Schlüsselform-Probe (N2, Form `https://…:1`
/ `oci://…:1`) ist als M2 geplant und in diesem Protokoll nicht geführt). Die
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

## Grenzen

- Die Manifest-Zugriffsform (`/v2/.../manifests/<tag>`) wurde nicht probeziert
 — der Token-Ausgang reicht fürs Urteil der Registry-Lage; die Antwortform
 `NAME_UNKNOWN` (Annahme der ADR) bleibt so nicht wörtlich gemessen, sondern
 als Token-Verweigerung belegt.
- Die Schlüsselform-Probe (N2: ob die Dev-Container-CLI `0.89.0` einen Schlüssel
 `https://…:1` oder `oci://…:1` auflöst) ist **M2, hier nicht geführt** — der
 Render-Schlüssel-Aspekt von [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) bleibt offen
 (keine Aussage über Form-Schwäche der externen Quellen steht in diesem Protokoll).
- Der Stand gilt am 2026-11-03; künftige Upstream-Änderungen alter das
 Protokoll, nicht die Mess-Basis. Die Messung liest nur; sie schreibt nichts.