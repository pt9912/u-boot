# ADR 0010: Kein HTTP-Driving-Adapter — u-boot bleibt CLI-only

## Status

Accepted

## Datum

2026-05-31

## Kontext

`spec/architecture.md` §7 nennt seit dem Bootstrap-Stand
(2026-05-22) als „geplante Erweiterung":

> *„HTTP-Driving-Adapter, falls u-boot perspektivisch eine
> Daemon-Variante bekommen soll."*

Keine Spec-Anforderung (`LH-FA-*`, `LH-AK-*`, `LH-SA-*`) fordert
einen HTTP-Adapter; auch `LH-OPEN-001..004` listen ihn nicht.
Der Slice-Plan
[`slice-later-http-driving-adapter`](../planning/done/slice-later-http-driving-adapter.md)
verankert diesen prospektiven Architektur-Hinweis als temporären
Carveout (`LH-FA-PROJDOCS-005`), ohne einen konkreten Use-Case
heute.

Hexagonale Architektur (`LH-FA-ARCH-001..003`, ADR-0002) erlaubt
weitere Driving-Adapter ohne Anpassung der Domain/Application-
Schicht — der HTTP-Adapter wäre eine reine Adapter-Ergänzung
unter `internal/adapter/driving/http/`. Trotzdem gilt: jeder
Adapter erhöht die Code-/Test-Surface dauerhaft.

`LH-NFA-USE-004` (Maschinen-lesbare Ausgabe) deckt bereits den
naheliegenden Maschinen-Schnittstellen-Bedarf ab: jeder relevante
Subkommando-Output ist via `--json`/`--dry-run`-Flags
maschinenlesbar (siehe `LH-FA-CLI-007/008` und das
`status`-Schema in §3.6). Für die heute absehbaren Use-Cases —
CI/CD-Pipelines, Editor-Integrationen, Skript-Orchestrierung —
ist die JSON-CLI-Schnittstelle ausreichend; ein HTTP-Server
würde keine zusätzlichen Fähigkeiten bringen.

Vergleichbare Tools:

- `kubectl`: CLI-only; HTTP wird durch das Kubernetes-API-Server-
  Ökosystem abgedeckt, nicht durch eigene CLI-Daemon-Variante.
- `helm`: CLI-only; ehemalige Tiller-Daemon-Architektur (Helm 2)
  wurde mit Helm 3 explizit zurückgebaut wegen Sicherheits- und
  Komplexitäts-Last.
- `gh`: CLI-only; nutzt die GitHub-API direkt statt eigenen
  HTTP-Server.
- `docker`: Daemon-getrieben, aber das ist der originale
  Use-Case (Container-Engine), nicht der u-boot-Use-Case
  (Compose-Orchestrierung als Build-Tooling).

Helm 3's Tiller-Rollback ist der direkteste Präzedenzfall: ein
prospektiver Daemon ohne klaren Use-Case hat zusätzliche Auth-,
Sandbox- und Versionierungs-Probleme produziert und wurde
zurückgebaut.

## Entscheidung

**Kein HTTP-Driving-Adapter.** u-boot bleibt CLI-only. Maschinen-
lesbare Schnittstellen werden ausschließlich über die bestehenden
`--json`/`--dry-run`-Flags (`LH-NFA-USE-004`, `LH-FA-CLI-007/008`)
ausgeliefert.

Konkrete Setzungen:

- **Kein `internal/adapter/driving/http/`-Verzeichnis.** Die
  hexagonale Driving-Schicht bleibt CLI-only (Cobra,
  [ADR-0005](0005-cli-framework-cobra.md)).
- **`spec/architecture.md` §7** wird auf den Stand „kein HTTP-
  Adapter geplant" umgeschrieben (analog ADR-0008-Pattern); der
  prospektive Bullet wird durch einen Verweis auf dieses ADR
  ersetzt, nicht stumm gelöscht.
- **JSON-CLI als kanonische Maschinen-Schnittstelle.** Wer
  u-boot programmatisch ansprechen möchte, nutzt `subprocess.run`
  / `os/exec`-Aufrufe mit `--json`-Flag und parst die
  strukturierte Ausgabe nach dem `LH-FA-CLI-007`-Schema.
- **Re-Evaluation-Trigger explizit dokumentiert** (siehe
  §Folgepunkte). Sobald einer der genannten Trigger eintritt,
  wird ein neuer `slice-v2-http-driving-adapter`-Plan in `open/`
  angelegt, der dieses ADR superseded.

## Konsequenzen

Positiv:

- **Minimale Adapter-Surface.** Nur eine Driving-Adapter-Implementierung
  (CLI/Cobra); keine zusätzliche Test-Strecke für HTTP-Server-
  Lifecycle, Routing, Authentifizierung, Concurrent-Request-Handling.
- **Keine Sicherheits-Diskussion.** Kein HTTP-Endpoint bedeutet keine
  TLS-/Auth-/CORS-/Rate-Limit-Pflichtsetzungen, kein Bind-Address-
  Default-Risiko, keine `LH-NFA-SEC-*`-Erweiterung für eine
  Daemon-Variante.
- **`LH-NFA-USE-004`-konsistent.** Maschinen-Lesbarkeit liegt
  bereits auf dem CLI-Layer; jedes Subkommando ist ohnehin
  `--json`-fähig, ein HTTP-Layer würde dieselbe Information nur
  über ein zweites Protokoll exponieren.
- **Konsistent mit ADR-0008** (statisches Add-on-System, keine
  prospektive Plugin-Architektur): prospektive Erweiterungen ohne
  Use-Case-Trigger werden nicht vorgebaut.

Negativ / Trade-offs:

- **Editor-/IDE-Integrationen müssen Subprocess-Aufrufe machen.**
  Wer u-boot aus einer Editor-Extension nutzen will, startet
  jedes Mal einen kurzen CLI-Prozess statt einen langlaufenden
  Daemon zu kontaktieren. Bei der typischen u-boot-Operationsfrequenz
  (Init, Add, Up, Down: mehrere Sekunden, nicht Sub-Millisekunden)
  ist der Subprocess-Overhead vernachlässigbar.
- **Keine Multi-Projekt-Orchestrierung aus einem Prozess heraus.**
  Wer mehrere u-boot-Projekte gleichzeitig steuern will, ruft
  mehrfache CLI-Prozesse auf — kein gemeinsamer State-Cache. Heute
  ist kein konkreter Use-Case dafür dokumentiert; falls einer
  auftaucht, greift Re-Eval-Trigger 1.
- **`spec/architecture.md` §7 Plugin-Bullet wurde mit ADR-0008
  schon geschlossen; HTTP-Bullet schließt mit diesem ADR.** Die
  §7-Sektion „Evolution" wird damit zu einem reinen ADR-Verweis-
  Block — was OK ist, weil die Architektur-Evolution per ADR
  läuft (`LH-FA-PROJDOCS-002`).

Alternativen (verworfen):

- **HTTP-Adapter ab `v0.1.0` bauen:** prospektive Architektur ohne
  Use-Case-Trigger. Helm 2/Tiller-Präzedenzfall (Daemon ohne
  klares Bedürfnis wurde später schmerzhaft zurückgebaut)
  spricht direkt dagegen.
- **HTTP-Adapter als optionaler Build-Tag (`//go:build http`):**
  optionaler Adapter halbiert die Test-Surface scheinbar, hält
  aber die Wartungslast voll (Build-Tag-Pfade müssen separat
  gepflegt werden). Kein Vorteil ohne konkreten Use-Case.
- **gRPC-Adapter stattdessen:** dieselben Trade-offs wie HTTP,
  zusätzlich Protocol-Buffer-Toolchain-Abhängigkeit. Heute
  ohne Use-Case ebenfalls verworfen.

## Folgepunkte (Re-Evaluation-Trigger)

Dieses ADR ist **revertierbar**: sobald einer der folgenden Trigger
eintritt, wird ein neuer Slice-Plan in `open/`
(`slice-v2-http-driving-adapter.md`, ADR-Nummer noch zu vergeben)
angelegt, der diesen ADR superseded.

1. **Konkreter Daemon-Use-Case.** Beispiele: Multi-Projekt-
   Orchestrierung mit gemeinsamem State, Web-Dashboard für
   Compose-Stack-Monitoring, langlaufender Health-Watcher.
   Trigger: ein dokumentierter Use-Case, der mit
   Subprocess-Aufrufen praktisch nicht abdeckbar ist.
2. **Maschinen-Schnittstelle über `LH-NFA-USE-004` hinaus.**
   Beispiele: Streaming-Output für lange Operationen (Compose-
   Logs in Echtzeit), bidirektionale Kommunikation (Inputs an
   Healthcheck-Probes durchreichen), Push-Notifications an
   Editor-Extensions. Trigger: ein dokumentiertes Anwendungs-
   szenario, das die `--json`-Ausgabe nicht effizient bedient.

Solange keiner dieser Trigger eintritt, bleibt u-boot CLI-only.

Der carveouts-Eintrag „HTTP-Driving-Adapter" wird mit der
[`slice-later-http-driving-adapter`](../planning/done/slice-later-http-driving-adapter.md)-
Closure entfernt.
