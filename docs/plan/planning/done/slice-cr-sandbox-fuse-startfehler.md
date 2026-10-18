# Slice CR: Ein fehlendes `/dev/fuse` auf dem Container-Host ist ein Startfehler, keine Degradationsstufe (`LH-FA-DEV-007`)

> **CR-Vehikel.** u-boot führt kein separates CR-Artefakt — dieser Slice
> **ist** der Change Request am Vertrags-Stratum (Precedent:
> [`slice-cr-sandbox-docker-runtime`](../done/slice-cr-sandbox-docker-runtime.md),
> [`slice-cr-devcontainer-agent-features-sandbox`](../done/slice-cr-devcontainer-agent-features-sandbox.md)).
> Anlass: Die Messung der [`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) (M13, M14) zeigt, dass das Lastenheft 0.12.0 an genau einer
> Stelle ein Ziel beschreibt, das unter Docker nicht erreichbar ist: „ein fehlendes `/dev/fuse` führt zu einem Fallback auf `vfs`-Storage“. Der erzeugte Container
> trägt `--device=/dev/fuse`; fehlt das Device auf dem Container-Host, lehnt die Engine den Start ab (`docker run` Exit-Code `127`, `docker start` Exit-Code `1`), das Startscript läuft nie.
> Der Projektinhaber hat am 2026-10-17 entschieden: **„A: Lastenheft-CR“** — ein fehlendes `/dev/fuse` wird wie `/dev/net/tun` als Startfehler geführt.
> Grundsatz (Memory `lastenheft-ist-zielbild`): Lastenheft und Specs sind Zielbild; der CR ist hier zulässig, weil die Messung das Ziel an dieser Stelle als unerreichbar zeigt und der Projektinhaber es entschieden hat.

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — Vertragspflege ohne Release-Bezug, keine Closure-Bedingung, die von der DoD dieses Slice verschieden wäre.

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Änderung),
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Devices als ausgewiesene Lockerung, unverändert),
[`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Code `11`),
[`ADR-0014`](../../adr/0014-nested-podman-sandbox-devcontainer.md), [`ADR-0021`](../../adr/0021-podman-build-im-sandbox-devcontainer.md),
[`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) (Entscheidung 11).

**Berührte Spec-Stellen:** `lastenheft.md` §4.3 ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer): Punkt „Bei `podman` startet …“ und Absatz Degradation), §16 (Historie), Versionsfeld;
`spezifikation.md` ([`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) Tabelle,
[`LH-FA-DEV-007.b`](../../../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman), Beschreibung der Prüfung `devcontainer.sandbox.runtime` in der Code-Registry, §7 Historie).
`architecture.md` berührt der CR nicht (gelesen: kein Treffer für `fuse`, `vfs`).
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Lastenheft-Änderung **0.12.0 → 0.13.0**: Ein fehlendes `/dev/fuse` auf dem Container-Host ist keine Degradationsstufe, sondern ein Startfehler der Container-Engine (wie `/dev/net/tun`).
Der Fallback auf `vfs`-Storage gilt nur noch für einen Container, der **ohne** das Device gestartet wurde. Spezifikation, Host-Prüfung in `u-boot doctor`, Handbuch und CHANGELOG folgen.

**Reihenfolge und Abhängigkeit (vom Projektinhaber gewünscht): dieser CR → [`slice-v1-sandbox-fallback-host-ausweisung`](../in-progress/slice-v1-sandbox-fallback-host-ausweisung.md).**
Der CR legt den Wortlaut der Host-Prüfung in `checkSandboxRuntime` und der Zeilen in [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) und [`.b`](../../../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman) fest; der Umsetzungs-Slice zieht danach [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), die übrigen Sätze von [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten), die Item-Beschreibung
und die neuen `doctor`-Items nach, ohne die hier gesetzten Zeilen erneut zu ändern. Der Umsetzungs-Slice startet erst, wenn dieser CR in `done/` liegt.

### Gefundener Stand (Ist-Stand, gelesen am 2026-10-18)

| Stelle | Aussage heute | Fundstelle |
|---|---|---|
| Lastenheft, Degradation | „Bei `warn` führen ein fehlendes `/dev/fuse` zu einem Fallback auf `vfs`-Storage …; bei `fail` sind alle drei ein Umgebungsproblem (Exit-Code `11`)“; für `/dev/net/tun` steht ausdrücklich: „keine Degradationsstufe: die Container-Engine lehnt den Start ab, und `u-boot doctor` prüft das Device nicht“ | `spec/lastenheft.md` Zeile 464 (Abschnitt [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Absatz **Degradation**) |
| Lastenheft, Start unter Docker | nennt als Voraussetzung nur `/dev/net/tun` auf dem Container-Host | `spec/lastenheft.md` Zeile 462 |
| Lastenheft, `--quiet` Ausnahme 2 | zählt `vfs`-Storage unter den Fallbacks; bleibt unverändert, weil `vfs` ein Fallback bleibt (nur seine Ursache wird enger) | `spec/lastenheft.md` Zeile 192 |
| Spezifikation, Degradationstabelle | Zeile „`/dev/fuse` nicht verfügbar“ → `vfs`, Warnung / Exit-Code `11` | `spec/spezifikation.md` Zeile 199 ([`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten)) |
| Spezifikation, Ausstattung `podman` | Absatz zu `/dev/net/tun` (Engine lehnt Start ab) ohne `/dev/fuse`; Satz „Ohne `/dev/fuse` baut es mit `vfs`-Storage“ | `spec/spezifikation.md` Zeile 327 ([`LH-FA-DEV-007.b`](../../../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman)) |
| Spezifikation, Code-Registry | `devcontainer.sandbox.runtime`: „nested Podman (`/dev/fuse`, Profil-Konsistenz)“ | `spec/spezifikation.md` Zeile 1071 |
| Code, erzeugte `runArgs` | `podmanRelaxations()` enthält `--device=/dev/fuse` und `--device=/dev/net/tun` | `internal/hexagon/application/devcontainer_sandbox.go`, `podmanRelaxations` |
| Code, Host-Prüfung | Linux-Host: fehlt `/dev/fuse`, Warnung „nested Podman will fall back to vfs storage (slower)“ mit `DegradationFallback`, bei Policy `fail` Fehler ohne Kennzeichen; Hint „Load the fuse module … or set `onUnavailable` to `warn`“. Nicht-Linux-Host: OK mit „the container start script checks /dev/fuse and user namespaces“ | `internal/hexagon/application/doctor_sandbox.go`, `checkSandboxRuntime` (Zweige ab `apparmorUsernsSysctl`, `s.fs.Exists("/dev/fuse")`) |
| Code, Startscript | `[ ! -c /dev/fuse ]` → `unavailable` (Exit-Code `11` bei `fail`), danach `storage.conf` mit `driver = "vfs"`; wird als `postCreateCommand` im Container ausgeführt, also nur, wenn der Container gestartet ist | `internal/hexagon/application/templates/devcontainer/sandbox-init.sh.tmpl`, Zeilen 20–25 |
| Test, Host-Prüfung | Fälle „fuse missing, warn policy“ (Warnung, „vfs“) und „fuse missing, fail policy“ (Fehler, „`fail`“); Flag-Test erwartet Warnung mit `DegradationFallback` bzw. Fehler ohne | `internal/hexagon/application/doctor_sandbox_test.go`, `TestDoctor_SandboxRuntime`, `TestDoctor_SandboxFallbackFlagIsOnTheItem` |
| Test, Startscript | Container wird mit allen `runArgs` **außer** `--device=/dev/fuse` gestartet; Skript schreibt `vfs`; `fail`-Policy: Exit-Code `11` | `internal/e2e/sandbox_devcontainer_docker_test.go` (Abschnitte „Degradation: without /dev/fuse …“ und `TestE2E_SandboxDevcontainer_FailPolicyWithoutFuse`) |
| Handbuch | Tabelle „`/dev/fuse` fehlt → Rückfall auf `vfs`-Speicher, Warnung / Abbruch mit Exit-Code 11“; „Ohne `/dev/fuse` baut Podman mit `vfs`-Speicher“; „`u-boot doctor` prüft `/dev/fuse` nur auf Linux-Hosts“; Fehlerbild „… oder fehlendes `/dev/fuse`“ unter Exit-Code 11 des Startscripts | `docs/user/benutzerhandbuch.md` Zeilen 1005, 1022, 1039, 1672; das Muster für `/dev/net/tun` steht in den Zeilen 987–989 |
| Messung | `docker run --device=/dev/fuse-nonexistent` Exit-Code `127`, `docker start` des erzeugten Containers Exit-Code `1` (`error gathering device information while adding custom device … no such file or directory`); ohne `--device=/dev/fuse` in `runArgs` startet der Container und das Skript wählt `vfs` | [`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) M13, M14 (Colima 0.10.3, Docker 29.5.2, arm64; Linux ohne VM nicht gemessen) |

**Entscheidung zum Startscript-Zweig (Vorschlag, gelesen gegen [`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) Entscheidung 11):** Der `vfs`-Zweig in `sandbox-init.sh.tmpl` **bleibt** als Rückfall für einen ohne das Device gestarteten Container (nach geänderten `runArgs`, M14);
er ist dann die einzige Quelle des `vfs`-Fallbacks und liefert die Statusdatei `vfs` des Umsetzungs-Slice. Er entfällt nicht, weil (a) die `runArgs` im Projekt des Nutzers liegen und änderbar sind, (b) eine Engine das Device
nicht prüfen könnte (Podman als Host-Engine ist nicht gemessen) und (c) [`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) ihn als ausgewiesenen Rückfall führt. Das Skript ändert dieser CR nicht (keine Zeile; die Status-Schreibung kommt im Umsetzungs-Slice).

### CR-Kern (Wortlaut-Vorschlag; im Lastenheft ohne Markierung)

Die Fassungen stehen als Vorschlag. **Der Projektinhaber gibt den Wortlaut vor dem Lastenheft-Commit frei**; bis dahin bleibt das Lastenheft auf 0.12.0.

**Stelle 1 — [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Punkt „Bei `podman` startet das Ergebnis …“** (nur der geänderte Satzteil; der Rest des Punktes bleibt wörtlich).

Alt:

> … und wenn das Device `/dev/net/tun` auf dem Container-Host vorhanden ist; mit dem Wert `1`, dem Standard einer Ubuntu-24.04-VM, endet das Startscript mit Exit-Code `11`, ohne das Device startet der Container nicht (Fehler der Container-Engine). …

Neu:

> … und wenn die Devices `/dev/fuse` und `/dev/net/tun` auf dem Container-Host vorhanden sind; mit dem Wert `1`, dem Standard einer Ubuntu-24.04-VM, endet das Startscript mit Exit-Code `11`, ohne eines der beiden Devices startet der Container nicht (Fehler der Container-Engine). …

**Stelle 2 — [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Absatz Degradation.**

Alt (vollständig):

> **Degradation:** `devcontainer.sandbox.onUnavailable` (`warn` | `fail`, Default `warn`) steuert den Umgang mit fehlenden Fähigkeiten: Bei `warn` führen ein fehlendes `/dev/fuse` zu einem Fallback auf `vfs`-Storage, eine nicht gewährbare Egress-Capability zum Wegfall der Restriktion und ein Docker-Daemon, der beim Containerstart nicht betriebsbereit wird, zu einem Container ohne Daemon, jeweils mit Warnung; bei `fail` sind alle drei ein Umgebungsproblem (Exit-Code `11`). Blockierte verschachtelte User-Namespaces (Seccomp/AppArmor) bei `nestedRuntime: podman` sind in beiden Modi ein Umgebungsproblem (Exit-Code `11`). Ein fehlendes `/dev/net/tun` auf dem Container-Host ist keine Degradationsstufe: die Container-Engine lehnt den Start ab, und `u-boot doctor` prüft das Device nicht. Eine ausdrücklich angeforderte Runtime wird nie stillschweigend ersetzt; jeder Fallback wird in der Befehlsausgabe und in `u-boot doctor` ausgewiesen. Ungültige Werte der Schlüssel führen zu einem fachlichen Validierungsfehler (Exit-Code `10`).

Neu (vollständig):

> **Degradation:** `devcontainer.sandbox.onUnavailable` (`warn` | `fail`, Default `warn`) steuert den Umgang mit fehlenden Fähigkeiten: Bei `warn` führen ein im Container fehlendes `/dev/fuse` (der Container wurde ohne das Device gestartet, etwa weil die erzeugten `runArgs` geändert wurden) zu einem Fallback auf `vfs`-Storage, eine nicht gewährbare Egress-Capability zum Wegfall der Restriktion und ein Docker-Daemon, der beim Containerstart nicht betriebsbereit wird, zu einem Container ohne Daemon, jeweils mit Warnung; bei `fail` sind alle drei ein Umgebungsproblem (Exit-Code `11`). Blockierte verschachtelte User-Namespaces (Seccomp/AppArmor) bei `nestedRuntime: podman` sind in beiden Modi ein Umgebungsproblem (Exit-Code `11`). Ein fehlendes `/dev/fuse` oder `/dev/net/tun` auf dem Container-Host ist keine Degradationsstufe: die Container-Engine lehnt den Start ab; `u-boot doctor` prüft `/dev/net/tun` nicht und weist ein fehlendes `/dev/fuse` nur aus, soweit der Host es bestimmen kann, als Umgebungsproblem (Exit-Code `11`). Eine ausdrücklich angeforderte Runtime wird nie stillschweigend ersetzt; jeder Fallback wird in der Befehlsausgabe und in `u-boot doctor` ausgewiesen. Ungültige Werte der Schlüssel führen zu einem fachlichen Validierungsfehler (Exit-Code `10`).

**Stelle 3 — Kopf und [§16 Historie](../../../../spec/lastenheft.md#16-historie).** Versionsfeld `0.12.0` → `0.13.0`; neue Zeile:

```text
| 0.13.0 | (Annahmedatum) | [`LH-FA-DEV-007`](#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer): Ein fehlendes `/dev/fuse` auf dem Container-Host ist wie `/dev/net/tun` keine Degradationsstufe, sondern ein Startfehler der Container-Engine (gemessen: Docker 29.5.2 auf Colima 0.10.3, `docker run` Exit-Code `127`, `docker start` Exit-Code `1`); der Fallback auf `vfs`-Storage gilt nur für einen Container, der ohne das Device gestartet wurde (etwa nach geänderten `runArgs`); `u-boot doctor` weist ein fehlendes `/dev/fuse` auf dem Container-Host als Umgebungsproblem aus, nicht als Fallback. | Vereinbarung mit dem Projektinhaber |
```

**Warum 0.13.0 und nicht 0.12.1:** Der CR nimmt eine Zusage zurück („ein fehlendes `/dev/fuse` führt zu einem Fallback auf `vfs`“); das ist nach dem Muster von `0.7.0` und `0.12.0` (Zusage auf gemessene Hosts begrenzt) eine Minor-Änderung.
Die Versionskonvention ergibt sich aus der Historie in §16 (kein eigener Regelabschnitt, gelesen: Einträge `0.3.x` für Klarstellungen, Minor bei geänderter Zusage). Die Wahl ist vom Projektinhaber zu bestätigen.

**Prüfung gegen die älteren Zusagen** (Suche im Lastenheft nach `fuse`, `vfs`, `Fallback`, `Degradation`, `ausgewiesen`):

| Stelle | Berührung | Ergebnis |
|---|---|---|
| [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Punkt `podman` startet … (Zeile 462) | nennt nur `/dev/net/tun` als Start-Voraussetzung | **geändert** (Stelle 1) |
| [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Absatz Degradation (Zeile 464) | `vfs` aus fehlendem `/dev/fuse`; Aussage zu `/dev/net/tun` | **geändert** (Stelle 2) |
| [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Ausnahme 2 (Zeile 192) | zählt `vfs`-Storage unter den Fallbacks | geprüft, **unverändert**: `vfs` bleibt Fallback, nur mit engerer Ursache |
| [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), Punkt „jede Sicherheitslockerung (… zusätzliche Devices) … ausweisen“ | `--device=/dev/fuse` ist weiterhin eine ausgewiesene Lockerung | geprüft, unverändert |
| [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), letzter Punkt | verweist auf die Degradation | geprüft, unverändert |
| `grep fuse` im Lastenheft | nur die Zeilen 462 und 464 | keine weitere Stelle |

**Unverändert, mit Begründung:**

- Die Konfigurationsschlüssel, `onUnavailable` und die Exit-Codes `10`/`11`: der CR ändert nur die **Ursache** des `vfs`-Fallbacks und die Einordnung des Host-Befunds.
- Die erzeugten `runArgs` und die Startscripte: Der CR verändert kein erzeugtes Artefakt; **Golden Cases bleiben unverändert** (Beleg: `generate_sandbox_test.go` pinnt die sechs Optionen, Zeilen 311–313).
- Die Ausweisung der Fallbacks in `doctor` und den Befehlen: [`slice-v1-sandbox-fallback-host-ausweisung`](../in-progress/slice-v1-sandbox-fallback-host-ausweisung.md) (nach diesem CR).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Kanal Container → Host und die neuen `doctor`-Items** — gehören dem Umsetzungs-Slice [`slice-v1-sandbox-fallback-host-ausweisung`](../in-progress/slice-v1-sandbox-fallback-host-ausweisung.md); der CR ändert nur Zusage und Host-Prüfung.
- **Eine Host-Prüfung für `/dev/net/tun`** — das Lastenheft sagt ausdrücklich „`u-boot doctor` prüft das Device nicht“; eine neue Prüfung wäre eine neue Zusage und ein anderer Vorgang.
- **Eine Änderung der erzeugten `runArgs` oder des Startscripts** — kein beobachtbares Verhalten der Generierung ändert sich; der `vfs`-Zweig bleibt (siehe Entscheidung zum Startscript-Zweig).
- **Neue Fallback-Klassen** — [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) zählt drei Fallbacks; der CR streicht keinen davon, er verengt die Ursache von `vfs`.
- **Messung auf Linux ohne VM, amd64, Docker Desktop** — der Startfehler ist auf Colima gemessen (M13); die Übertragung ist angenommen und steht im Wortlaut der CHANGELOG-Zeile als Messstand.

## 2. Definition of Done

- [x] **Lastenheft 0.13.0**: Stellen 1 und 2 in der vom Projektinhaber freigegebenen Fassung, Versionsfeld, Historie-Zeile (Verweis „Vereinbarung mit dem Projektinhaber“)
      ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)).
- [x] **Nachzug in Spezifikation, Code und Test:** [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) (Tabelle),
      [`LH-FA-DEV-007.b`](../../../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman), Beschreibung von `devcontainer.sandbox.runtime`, Historie-Zeile in `spezifikation.md` §7;
      `checkSandboxRuntime` weist ein fehlendes `/dev/fuse` auf dem Linux-Host als Fehler ohne `DegradationFallback` aus und sagt nicht mehr „fall back to vfs“; Tests dazu und ein Integrationstest für den Startfehler der Engine.
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` und CHANGELOG nennen den geltenden Stand (öffentlicher Vertrag; Änderung der `doctor`-Einstufung benannt).
- [x] `make gates` grün; `make test-docker` auf dem Entwicklungs-Host (Colima) grün oder mit benanntem Grund nicht gelaufen; `make doc-trace` ohne neue Waise.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Reihenfolge: Wortlaut-Freigabe des Projektinhabers → Lastenheft → Spezifikation → Code und Tests → Handbuch und CHANGELOG → Gates → Review → Closure. Danach startet
[`slice-v1-sandbox-fallback-host-ausweisung`](../in-progress/slice-v1-sandbox-fallback-host-ausweisung.md).
Die Tabelle nennt jede Datei, die der Diff voraussichtlich berührt (§3.10 von [`AGENTS.md`](../../../../AGENTS.md)); der Implementer gleicht sie vor dem Abschluss gegen `git diff --stat` ab.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Stellen 1 und 2, Versionsfeld (Zeile 9) `0.13.0`, §16-Zeile; erst nach Freigabe des Wortlauts |
| `spec/spezifikation.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten): Tabellenzeile `/dev/fuse` (Zeile 199) aufteilen — Zeile „`/dev/fuse` im Container nicht verfügbar (Container ohne `--device=/dev/fuse` gestartet)“ mit `vfs`/Exit-Code `11`, danach ein Satz: ein fehlendes `/dev/fuse` oder `/dev/net/tun` auf dem Container-Host ist keine Degradationsstufe, die Engine lehnt den Start ab, `doctor` weist `/dev/fuse` nur aus, soweit der Host es bestimmt. [`LH-FA-DEV-007.b`](../../../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman) (Zeile 327): Satz zu `/dev/net/tun` auf beide Devices erweitern; „Ohne `/dev/fuse` baut es mit `vfs`-Storage“ → „Fehlt `/dev/fuse` im Container (ohne das Device gestartet), baut es mit `vfs`-Storage“. Code-Registry (Zeile 1071): `devcontainer.sandbox.runtime` nennt `/dev/fuse` auf dem Linux-Host als Fehler bei Fehlen. §7 Historie. Die übrigen Sätze von `.a` („Jeder Fallback wird …“, Prüfsatz) und [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) Ausnahme 2 fasst erst der Umsetzungs-Slice an |
| `internal/hexagon/application/doctor_sandbox.go` | update | `checkSandboxRuntime`: der Zweig `/dev/fuse` fehlt (Linux-Host) liefert **immer** `SeverityError` (unabhängig von `onUnavailable`), kein `DegradationFallback`, Text sagt, dass die Engine den Start ablehnt (nicht „fall back to vfs“); der Nicht-Linux-Zweig nennt „the container engine refuses to start without /dev/fuse and /dev/net/tun; the start script checks user namespaces“ statt „checks /dev/fuse“; Hint zu `modprobe fuse` bzw. `nestedRuntime: none` (der Hinweis auf `onUnavailable: warn` entfällt); der Funktionskommentar (Zeilen vor `checkSandboxRuntime`) nennt den neuen Befund |
| `internal/hexagon/application/doctor_sandbox_test.go` | update | `TestDoctor_SandboxRuntime`: die beiden Fälle „fuse missing, warn policy“/„fail policy“ erwarten beide `SeverityError` mit „refuses“ bzw. „engine“; neuer Fall: `fail`-Policy gleicher Text. `TestDoctor_SandboxFallbackFlagIsOnTheItem`: Zeile „runtime: vfs fallback (warn policy)“ wird zu „fuse missing: error, unflagged“ (`fallback=false`); der Kommentar („Host-side, only the `vfs` fallback is determinable“) wird berichtigt — Happy/Boundary/Negative nach [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) |
| `internal/e2e/sandbox_devcontainer_docker_test.go` (Build-Tag `docker`) | update | neuer Test: der erzeugte Container startet mit einem nicht vorhandenen Device (`--device=/dev/fuse-nonexistent` statt `--device=/dev/fuse`) nicht (Exit-Code ≠ 0, Meldung `error gathering device information`) — pinnt M13 auf der lokalen Engine; die beiden Bestandstests (`vfs` ohne Device, `FailPolicyWithoutFuse`) bleiben, ihr Kommentar „without /dev/fuse“ sagt „container started without the device“ |
| `internal/adapter/driving/cli/quiet_test.go`, `cli_test.go`, `doctor_json_test.go`, `quiet_modifying_test.go` | prüfen | enthalten Fake-Einträge mit „falls back to vfs storage“ als Vertreter eines `DegradationFallback`-Items (Fixtures, nicht der Text der Prüfung); voraussichtlich unverändert, weil der Umsetzungs-Slice ein echtes `vfs`-Item einführt |
| `internal/hexagon/application/templates/devcontainer/sandbox-init.sh.tmpl` | prüfen | unverändert (Zweig bleibt, Kommentar „without it fuse-overlayfs fails, fall back to vfs storage“ gilt für den ohne Device gestarteten Container) |
| `docs/user/benutzerhandbuch.md` | update | Zeilen 987–989 (`/dev/net/tun`-Absatz) um `/dev/fuse` erweitern; Zeile 1005 („Ohne `/dev/fuse` baut Podman mit `vfs`-Speicher“) auf den ohne Device gestarteten Container begrenzen; Tabelle Zeile 1022: `/dev/fuse` fehlt im Container → `vfs`/Exit-Code 11, plus Satz, dass ein fehlendes Device auf dem Container-Host den Start verhindert; Zeile 1039 („`u-boot doctor` prüft `/dev/fuse` nur auf Linux-Hosts“) nennt den Fehler-Eintrag; Zeile 1672 (Fehlerbild) verschiebt „fehlendes `/dev/fuse`“ vom Exit-Code 11 des Startscripts zum Startfehler der Engine |
| `CHANGELOG.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | `[Unreleased]` / `Changed` (Muster: `/dev/net/tun`-Eintrag im Abschnitt `[0.10.0]`, unverändert): ein fehlendes `/dev/fuse` auf dem Container-Host ist ein Startfehler der Engine, kein `vfs`-Fallback; `u-boot doctor` weist es auf einem Linux-Host als **Fehler** aus (Exit-Code `11`, bisher Warnung mit Exit-Code `0`, mit `onUnavailable: fail` schon Fehler); Lastenheft 0.13.0; Messstand (Colima, Docker 29.5.2, arm64); das Erzeugnis von `generate devcontainer` ändert sich nicht |
| `docs/plan/planning/open/README.md` | update | Bestandszeile dieses Slice: Freigabe-Hinweis ersetzt durch den Stand (Lastenheft 0.13.0 angenommen); der Link auf `done/` kommt von `make slice-mv` |
| `docs/plan/planning/in-progress/roadmap.md` | update | Ruhe-Marker „Nichts in Arbeit.“ wieder gesetzt (kein Slice in `in-progress/`; d-check `planning`), Datum „Letzte Änderung“; der Diff gegen den Stand vor dem Slice ist nur das Datum |
| `docs/plan/planning/next/README.md`, `docs/plan/planning/in-progress/slice-v1-sandbox-fallback-host-ausweisung.md` | update | Bestandszeile und Verweise auf diesen Slice zeigen nach `done/` (`make slice-mv`); die Zeile nennt den CR als erledigt und trägt die Zeichenkodierung von „unverändert“ wieder korrekt |
| `docs/plan/planning/done/…` (Closure) | update | dieser Slice liegt nach dem reinen Move in `done/`; DoD, §6 und §7 gefüllt |
| `docs/plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/evidence/slice-cr-sandbox-fuse-startfehler.md` | neu | neunter Beleg der Beobachtung (Review F-2) |
| `docs/reviews/2026-10-18-verification-slice-cr-sandbox-fuse-startfehler.md` | neu | Verification Evidence |
| `README.md`, `README.de.md` | update | Statuszeile nennt als Host-Voraussetzung `/dev/fuse` und `/dev/net/tun` (EN/DE gespiegelt); die Suche nach `fuse` fand sie nicht, weil dort nur `/dev/net/tun` stand |
| `docs/reviews/2026-10-18-slice-cr-sandbox-fuse-startfehler.md` | neu | Review-Report |
| `architecture.md`, `docs/maintainer/cli-json-contract.md`, `.golangci.yml`, `.d-check.yml`, Carveout-Inventar | geprüft, unverändert | `fuse`/`vfs` kommt dort nicht vor (Suche im Baum); kein neuer Carveout |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Projektinhaber hat den Wortlaut der Stellen 1 und 2 und die Versionswahl `0.13.0` freigegeben;
[`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) ist `Accepted`. Der Slice geht von `open/` nach `next/` mit der Freigabe.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Freigabe verlangt, auch eine Host-Prüfung für `/dev/net/tun` oder das Entfernen des `vfs`-Zweigs aus dem Startscript mitzuregeln — dann je Zusage ein Slice.
- `in-progress` → `open` (blockiert): Der Projektinhaber nimmt „A“ zurück (Zielbild halten, `vfs` wird für Docker gebaut) — dann entfällt der CR mit `Gegenstand:` nach dem Muster von [`slice-cr-sandbox-fallback-ausweisung`](../done/slice-cr-sandbox-fallback-ausweisung.md).

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7). Lastenheft 0.13.0 ist angenommen; der Umsetzungs-Slice
[`slice-v1-sandbox-fallback-host-ausweisung`](../in-progress/slice-v1-sandbox-fallback-host-ausweisung.md) kann starten.

## 6. Risiken und offene Punkte

- **Die Einstufung in `doctor` ändert sich:** Auf einem Linux-Host mit fehlendem `/dev/fuse` war die Prüfung eine Warnung (Exit-Code `0`, mit `--strict` `11`), sie wird ein Fehler (Exit-Code `11`). Ein Skript, das `doctor` ohne `--strict` auswertet, scheitert dort künftig. — **Ausgang:** eingetreten: der Projektinhaber hat „Immer Fehler, Exit 11“ entschieden (die Alternative „Warnung ohne Kennzeichen“ entfällt); der CHANGELOG nennt die Regression für Skripte ohne `--strict` (`[Unreleased]` · Changed).
- **Falsch-positiv bei entfernter Engine:** Die Host-Prüfung geht davon aus, dass die Engine auf dem Host läuft (so schon die AppArmor-Prüfung in `checkSandboxRuntime`); bei einem entfernten `DOCKER_HOST` kann `/dev/fuse` lokal fehlen, während die Engine es hat; ein Fehler statt einer Warnung verschärft das. — **Ausgang:** weiter offen: Grenze, nicht zugesagt; der Vertragstext bindet die Einstufung nicht an eine lokale Engine (Review F-1, behoben in `f839d28`), `checkSandboxRuntime` liest keinen `DOCKER_HOST`; der Umsetzungs-Slice [`slice-v1-sandbox-fallback-host-ausweisung`](../in-progress/slice-v1-sandbox-fallback-host-ausweisung.md) oder ein eigener Vorgang entscheidet, ob die Prüfung bei entferntem `DOCKER_HOST` entfällt.
- **Messung nur auf Colima:** Der Startfehler (M13) ist auf Colima/Docker 29.5.2 gemessen; Linux ohne VM, amd64, Docker Desktop und Podman als Host-Engine sind angenommen bzw. nicht zugesagt. — **Ausgang:** weiter offen: der Messstand steht im Lastenheft-Eintrag (§16, 0.13.0) und in der CHANGELOG-Zeile; die Verification hat `TestE2E_SandboxDevcontainer_MissingHostFuseDeviceRefusesStart` erneut auf Colima gemessen, andere Engines bleiben ungemessen.
- **Der `vfs`-Zweig bleibt ein Fallback mit enger Ursache:** Das Lastenheft führt ihn weiter in der Degradationstabelle und in `--quiet` Ausnahme 2, obwohl er unter den erzeugten `runArgs` nicht eintritt. — **Ausgang:** weiter offen: der Projektinhaber hat „vfs-Fallback bleibt“ entschieden (die Alternative „Zweig und Fallback entfallen“ entfällt); getragen vom Umsetzungs-Slice [`slice-v1-sandbox-fallback-host-ausweisung`](../in-progress/slice-v1-sandbox-fallback-host-ausweisung.md), der die Statusdatei `vfs` und den Nachzug von [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) Ausnahme 2 und [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 2 liefert (Review F-5).
- **Überschneidung mit dem Umsetzungs-Slice:** beide ändern `spezifikation.md` im Absatz [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten). — **Ausgang:** entfallen: Reihenfolge CR → Umsetzung (§1) und Aufteilung der Sätze in §3 gehalten; [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) und die übrigen Sätze von [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) blieben unberührt (Verification, Plan-Tabelle gegen `git diff --stat`).
- **Spec-Nachzug ohne Verweise nach unten:** Spezifikation darf weder ADR noch diesen Slice nennen. — **Ausgang:** entfallen: `make docs-check` ohne Befund, der Review fand keinen Verweis nach unten in Lastenheft und Spezifikation.

## 7. Closure-Notiz

**Belege:** Review-Report [`2026-10-18-slice-cr-sandbox-fuse-startfehler`](../../../reviews/2026-10-18-slice-cr-sandbox-fuse-startfehler.md) (0 HIGH, 0 MEDIUM, 3 LOW, 4 INFO; sechs Mutationen getötet, eine überlebte; F-1 bis F-4 behoben in `f839d28`, `2ff6db8`; F-5 bewusst offen, Umsetzungs-Slice; F-6, F-7 ohne Handlung). Verification Evidence [`2026-10-18-verification-slice-cr-sandbox-fuse-startfehler`](../../../reviews/2026-10-18-verification-slice-cr-sandbox-fuse-startfehler.md) (keine Diskrepanz; Auflagen 1 bis 5 in dieser Closure erfüllt). `make gates` grün (Coverage 91,20 %), `make test-docker` grün (Colima), `make doc-trace`: 83 Anforderungen, 0 Waisen, `make doc-immutable`: 0 Befunde. Commits `c1ef60e` (Lastenheft), `4c16901` (Spezifikation), `a28a9f5` (Code, Tests), `e4e5d0d` (Handbuch, CHANGELOG), `b0919ad` (Roadmap), `702335a`, `2ff6db8`, `f839d28`, `1204c1e` (Nachzüge), `175d02d` (Review-Report).

**Freigaben des Projektinhabers (2026-10-17/18):** Der Wortlaut wurde vorgelegt, die Antworten gelten als Freigabe: „ADR-0022 Accepted“ · „doctor, plus generate“ · „A: Lastenheft-CR“ · Version „0.13.0“ · „Immer Fehler, Exit 11“ · „vfs-Fallback bleibt“. Weitere Freigaben liegen nicht vor. <!-- d-check:ignore (Wortlaut der Freigabe als Zitat) -->

**Fundstellen je Aussage (`AGENTS.md` §3.10):**

| Aussage | Fundstelle |
|---|---|
| Ein fehlendes `/dev/fuse` oder `/dev/net/tun` auf dem Container-Host ist ein Startfehler der Engine; `vfs` nur für den ohne das Device gestarteten Container; Version 0.13.0 | `spec/lastenheft.md` · Abschnitt [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Punkt „Bei `podman` startet …“, Absatz **Degradation**), Versionsfeld, §16-Zeile 0.13.0 |
| Tabellenzeile „`/dev/fuse` im Container nicht verfügbar“, Host-Fall kein Degradationsschritt, Ausstattung mit beiden Devices, Registry-Beschreibung | `spec/spezifikation.md` · [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten), [`LH-FA-DEV-007.b`](../../../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman), Code-Registry-Zeile `devcontainer.sandbox.runtime`, §7 Historie |
| Fehlendes `/dev/fuse` auf dem Linux-Host ist immer `SeverityError`, ohne `DegradationFallback`, unabhängig von `onUnavailable`; Nicht-Linux-Zweig OK mit Hinweis auf die Engine-Ablehnung | `checkSandboxRuntime` in `internal/hexagon/application/doctor_sandbox.go` |
| Beide Policies, Nicht-Linux-Fälle, Hint, Flag `fallback=false` gepinnt | `TestDoctor_SandboxRuntime`, `TestDoctor_SandboxRuntime_FuseMissingHint`, `TestDoctor_SandboxFallbackFlagIsOnTheItem` in `internal/hexagon/application/doctor_sandbox_test.go` |
| Die Engine lehnt den Start mit einem nicht vorhandenen Device ab (M13), Kontrolllauf mit echten `runArgs` | `TestE2E_SandboxDevcontainer_MissingHostFuseDeviceRefusesStart` in `internal/e2e/sandbox_devcontainer_docker_test.go` |
| Handbuch, CHANGELOG, READMEs nennen beide Devices und die geänderte `doctor`-Einstufung | `docs/user/benutzerhandbuch.md` (Absätze zu `/dev/net/tun` und `/dev/fuse`, Tabelle, Fehlerbild); `CHANGELOG.md` · `[Unreleased]` · Changed; Statuszeile in `README.md` und `README.de.md` |
| `runArgs` und Startscript unverändert (Golden Cases unberührt) | `podmanRelaxations` in `internal/hexagon/application/devcontainer_sandbox.go`; `internal/hexagon/application/templates/devcontainer/sandbox-init.sh.tmpl` (kein Diff) |

- **Was hat funktioniert:** Die Messung (M13, M14) vor dem CR machte aus einer Annahme einen Vertragsbefund; der Verifier maß den Startfehler auf Colima erneut (Test und Praxislauf von `doctor` mit und ohne `--device /dev/fuse`). Die Mutationsläufe des Reviews trafen sechs von sieben Zweigen; die siebte (Hint) führte zu `TestDoctor_SandboxRuntime_FuseMissingHint`. Die Prüfung gegen ältere Zusagen im Lastenheft (Suche nach `fuse`, `vfs`, `Fallback`) fand die beiden geänderten Stellen und keine weitere.
- **Was ging anders als geplant:** Die Plan-Tabelle führte `README.md` und `README.de.md` zunächst als „unverändert“ (Review F-2); die Suche nach `fuse` fand sie nicht, weil dort nur `/dev/net/tun` stand. Die Spezifikation band die Einstufung an eine „lokale Engine“, die der Code nicht prüft (F-1); der Hinweistext der Prüfung war ungepinnt (F-3); eine Testzeile war nicht `gofmt`-konform (F-4). Alle vier wurden vor der Verification behoben. §6 stand nach der Lieferung noch auf „offen“ und wurde erst in der Closure gesetzt (Verifier-Auflage 2), obwohl die Entscheidungen früher vorlagen.
- **Finding-Klassen des Reviews:** Spec-Bedingung ohne Code-Entsprechung (F-1) · Öffentliche Voraussetzungs-Aussage nicht auf beide Devices nachgezogen (F-2) · Hinweistext einer Prüfung ungepinnt (F-3) · Formatierung ohne Gate (F-4) · Bewusst offener Nachzug (F-5) · Zeilenverweis im Auftrag (F-6) · Sensor nicht gelaufen (F-7).
- **Steering-Loop-Eintrag:** **Gezählt, nicht verkörpert.** Kein `liegt in`, nichts verkörpert, keine Entscheidung über Verkörperung. **Lehren:** (1) Eine Plan-Tabelle prüft eine Datei gegen die *Aussage*, nicht gegen das *Wort*: Für die READMEs galt die Suche nach `fuse`, die Aussage (Voraussetzung auf dem Container-Host) nannte nur ein Device (F-2). (2) Ein Hinweistext einer Prüfung ist Teil ihres Vertrags und braucht einen Pin (F-3). F-2 gehört zu [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (die Plan-Aussage „README unverändert“ stimmte nicht): neunter Beleg als `evidence/slice-cr-sandbox-fuse-startfehler.md`. Die Beobachtung steht `verkörpert` in `AGENTS.md §3.10` (acht Belege bei der Verkörperung); dieser Vorgang tritt **nach** der Verkörperung auf, die Regel hat ihn nicht verhindert, und der Plan nannte die Fundstellen für Code, nicht für die READMEs. **Nur genannt, nicht entschieden:** ob die Regel dafür geschärft wird (Aussage statt Wort suchen) oder ob ein Rückfall nur gezählt wird, entscheidet der Projektinhaber. F-3 und F-4 sind einmalig und ohne bestehenden Pfad; sie stehen hier benannt, nicht im Register.
- **Beobachtungs-Register (`../observations/`):** ein Beleg fortgeschrieben (`BEO-DOK/plan-aussage-ueber-code`, siehe oben); keine neue Beobachtung angelegt.
- **Hinweis an den Projektinhaber (Nebenbefund des Implementers):** Viele Bestandsdateien sind nicht `gofmt`-konform (z. B. `internal/adapter/driving/cli/cli.go`, `addservice.go`, `git_test.go`); `.golangci.yml` aktiviert keinen Formatter, daher fängt kein Gate die Abweichung (Review F-4). Zu entscheiden: ein Formatter-Gate (Verschärfung eines Gates, `AGENTS.md` §3.6: Plan- oder ADR-Anker nötig) und ein Format-Slice für den Bestand. Dieser Slice ändert beides nicht; kein Register-Eintrag (ein Beleg, Zuordnung zu einer Sub-Area `WZG` offen).
- **Folge-Slices:** [`slice-v1-sandbox-fallback-host-ausweisung`](../in-progress/slice-v1-sandbox-fallback-host-ausweisung.md) (liegt in `next/`, startet mit diesem Slice in `done/`); er trägt den offenen Nachzug von [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) Ausnahme 2 und [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 2 (Review F-5).
- **Risiken aus §6:** alle sechs tragen einen Ausgang: eingetreten (Einstufung `doctor`), weiter offen (entfernter `DOCKER_HOST`, Messung nur auf Colima, `vfs`-Zweig mit enger Ursache), entfallen (Überschneidung mit dem Umsetzungs-Slice, Spec ohne Verweise nach unten).
- **Roadmap:** Der Ruhe-Marker „Nichts in Arbeit.“ steht nach dem Move wieder unter „Offene Wellen“ (`in-progress/` trägt keinen Slice; d-check `planning`).
- **Grenzen:** (1) Der Startfehler ist auf Colima (Docker 29.5.2, arm64) gemessen; Linux ohne VM, Docker Desktop, amd64 und Podman als Host-Engine sind angenommen. (2) Der Nicht-Linux-Zweig von `checkSandboxRuntime` ist nur im Unit-Test, nicht an einem macOS-Binary geprüft. (3) Host-seitig ist `vfs` nicht mehr ermittelbar; die Ausnahme 2 in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) und [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) nennt es bis zum Umsetzungs-Slice weiter.
- **Drei Paarungen:** Anker (kein `liegt in`, nichts verkörpert) · Folge-Slice (`slice-v1-sandbox-fallback-host-ausweisung` in `next/`) · Register (ein Beleg unter `BEO-DOK/plan-aussage-ueber-code`) — nach dem `git mv` nach `done/` geprüft.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/`, `harness/`, `docs/plan/` (DOK), `internal/hexagon/application` (APP) und `internal/e2e` (E2E) — alle Greenfield laut
[`harness/conventions.md`](../../../../harness/conventions.md#mr-006) (Abschnitt Modus-Deklaration), jede erfüllt die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig.
`docs/user/` und `CHANGELOG.md` sind öffentliche Verträge, keine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register
([`../observations/`](../observations/README.md)) durchgegangen — Verzeichnisse `BEO-CLI`, `BEO-DOK`, `BEO-E2E`, `BEO-WZG` liegen vor; keiner betrifft `/dev/fuse` oder die Degradation.
`BEO-DOK/plan-aussage-ueber-code` ist beachtet: jede Aussage über Code in §1 und §3 steht mit Fundstelle, an der Stelle gelesen am 2026-10-18.

Alle berührten Sub-Areas GF.
