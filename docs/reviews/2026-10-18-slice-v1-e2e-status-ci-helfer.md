# Review-Report: slice-v1-e2e-status-ci-helfer — 2026-10-18

**Review-Art:** Code — geprüft gegen Plan, Hard Rules und Konventionen (nicht gegen die DoD; das ist Verifier-Aufgabe).

**Gegenstand:** Slice [`slice-v1-e2e-status-ci-helfer`](../plan/planning/done/slice-v1-e2e-status-ci-helfer.md), Commit-Range `575281f..HEAD` (`814072a`, `db4eb13`, `b440a1c`); Code-Änderung: `internal/e2e/sandbox_status_docker_test.go`, neu `internal/e2e/container_id_test.go`.

**Skill:** `.harness/skills/reviewer.md` @ `3ec4535` · <!-- d-check:ignore (Skill-Pfad, existiert im Repo) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-18

**Eingangs-Kontext:**

- Slice-Plan [`slice-v1-e2e-status-ci-helfer`](../plan/planning/done/slice-v1-e2e-status-ci-helfer.md) (§1 Befund, §2 DoD, §3 Plan-Tabelle, §6 Risiken)
- [`ADR-0022`](../plan/adr/0022-fallback-ausweisung-statusdateien-im-container.md), [`ADR-0004`](../plan/adr/0004-ci-system.md)
- [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)
- `spec/architecture.md` Abschnitt Tests (Build-Tag-Konvention)
- `AGENTS.md` §3.2, §3.6, §3.7, §3.9, §3.10

**Eigene Läufe dieses Reviews** (Docker-only):

- Vorher-Stand (`git archive 575281f`), `alpine:3` und `alpine:latest` entfernt, Einzeltest im Image `u-boot:test-docker-tools` mit der Rezeptur aus `repo.mk` (Ziel `test-docker`): **rot**, `sandbox_status_docker_test.go:362: seed files: Error response from daemon: page not found` (3,28 s). Das ist der Fehler des CI-Laufs 39103939800, auf Kommandoebene und im Test reproduziert.
- Kommandoebene mit der CLI des Test-Images (Docker 29.9.1), `alpine:3` nicht gecacht: `docker run -d alpine:3 sleep 5` schreibt **stdout** nur die 64-stellige ID; **stderr** `Unable to find image 'alpine:3' locally`, `3: Pulling from library/alpine`, Layer-Zeilen, `Digest:`, `Status: Downloaded newer image`. Mit `-d -q`: stdout nur die ID, stderr leer (ungecacht und gecacht); `-q` ist mit `-d` wirksam, Exit 0.
- Aktueller Stand, Einzeltest `-run StatusFiles_InspectorAgainstEngine -v -count=1`: ungecachtes `alpine:3` **grün, 7,51 s**; gecacht **grün, 5,48 s**. Danach: keine Container mit Vorfahr `alpine:3`, keine `devcontainer.local_folder`-Container übrig.
- `make test` (Baseline): `ok github.com/pt9912/u-boot/internal/e2e 0.006s`. `make lint`: Exit 0.
- Mutationen an `container_id_test.go` (nur `make test`, danach `git checkout -- <Datei>`, Arbeitsbaum am Ende sauber):

| Mutation | Ergebnis |
|---|---|
| Hex-Länge 64 → 63 | erkannt (4 Fälle rot) |
| letzte → erste nichtleere Zeile | erkannt (3 Fälle rot) |
| Großbuchstaben erlaubt (`[0-9a-fA-F]`) | erkannt (`upper-case hex` rot) |
| `strings.TrimSpace` der Zeile entfernt | **überlebt** (F-2) |
| CRLF-`ReplaceAll` entfernt | **überlebt** (F-2) |

- Zustand des Hosts am Ende: `alpine:3` und `alpine:latest` (Image `294b683cb724`) wieder gecacht wie zu Beginn; Zwischenverzeichnisse gelöscht.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Plan-Tabelle §3 führt `docs/plan/planning/in-progress/roadmap.md` nicht, obwohl `b440a1c` sie ändert (Ruhe-Marker durch Verweis auf den Slice ersetzt); die Commit-Message nennt den Abgleich gegen `git diff --stat`. `docs/plan/planning/open/README.md` und [`slice-v1-podman-aktuelle-version`](../plan/planning/in-progress/slice-v1-podman-aktuelle-version.md) stehen ebenfalls im Range (anderer Vorgang, Commit `db4eb13`) und nicht in der Tabelle. | `AGENTS.md` §3.10 (Tabelle listet jede berührte Datei, Abgleich gegen `git diff --stat`) | `docs/plan/planning/done/slice-v1-e2e-status-ci-helfer.md` · "<!-- d-check:ignore (Zitat der Plan-Tabellenzeile) -->`docs/plan/planning/done/slice-v1-e2e-status-ci-helfer.md` &#124; update" | nein — Review-only (`git diff --stat 575281f..HEAD` gegen §3) | Plan-Tabelle gegen diff --stat unvollständig |
| F-2 | INFO | Zwei Mutationen von `containerIDFromRunOutput` überleben `make test`: ohne `TrimSpace` und ohne den CRLF-`ReplaceAll`. Jede der beiden Stellen deckt die andere beim CRLF-Fall; ein Fall mit Leerzeichen oder Tab um die ID auf derselben Zeile fehlt, und ein Fall mit mehreren IDs (letzte gewinnt) ebenso. | Maintainability | `internal/e2e/container_id_test.go` · "CRLF and trailing blank lines" | ja — `make test` mit der Mutation | Mutant überlebt (redundante Normalisierung) |
| F-3 | INFO | `docker pull -q` vor dem Test macht den Lauf von der Registry abhängig, auch wenn `alpine:3` lokal gecacht ist (Offline-Lauf, Rate-Limit); ein Pull-Fehler endet mit `t.Fatalf`, nicht mit `t.Skip` — die Haltung ist ehrlich und entspricht §3 Schritt 3 und §6 des Plans. Der Plan nennt den Helfer allein als ausreichend für die Korrektheit. | Maintainability | `internal/e2e/sandbox_status_docker_test.go` · `docker pull alpine:3: %s` | nein — Review-only (nur mit Netz prüfbar) | Netzabhängigkeit durch Vorab-Pull |
| F-4 | INFO | Schlägt die ID-Prüfung nach einem erfolgreichen `docker run -d` fehl (z. B. eine weitere Zeile nach der ID), endet `startLabeled` mit `t.Fatalf`, bevor `t.Cleanup` registriert ist; der Container bleibt dann liegen. Der Erfolgspfad räumt mit `rm -f -v <ID>` auf. | Maintainability | `internal/e2e/sandbox_status_docker_test.go` · `t.Fatalf("%v\nstderr:\n%s", err, stderr.String())` | nein — Review-only | Cleanup erst nach ID-Prüfung registriert |
| F-5 | INFO | Der Doc-Kommentar von `startLabeled` sagt, stderr trage den Pull-Fortschritt eines nicht gecachten Images; mit `-q` bleibt stderr leer (am Test-Image-CLI gemessen). Die Aussage gilt für `docker run` allgemein und begründet das Lesen nur von stdout; sie ist kein Verstoß gegen §3.7, aber ungenauer als der Zustand. | `AGENTS.md` §3.7 | `internal/e2e/sandbox_status_docker_test.go` · "stderr carries the pull" | nein — Review-only | Kommentar ungenau zum Ist-Zustand |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `containerIDFromRunOutput`: letzte nichtleere Zeile, `^[0-9a-f]{64}$`, CRLF, Leerzeilen, Zeile mit Präfix, Kurz-ID, ID gefolgt von Fehlerzeile (Fehler), leer | geprüft, ohne Befund; Verhalten gelesen und über drei erkannte Mutationen belegt |
| `startLabeled`: stdout/stderr in getrennten `strings.Builder`, `cmd.Run` wartet auf beide Kopier-Goroutinen (kein Race, kein Pufferproblem); Exit-Code über `cmd.Run`-Fehler; `t.Fatalf` mit stdout und stderr; ID-Fehlerpfad mit stderr | geprüft, ohne Befund |
| `docker run -d -q`: stdout nur die ID, mit und ohne gecachtes Image, `-q` wirksam | geprüft, ohne Befund (am CLI des Test-Images gemessen) |
| Cleanup mit ID statt Name: `rm -f -v <ID>` im `t.Cleanup`, Aufrufer nutzen den Rückgabewert für `docker exec`, `stop`, `inspect`, `removeContainer`; nach dem Lauf bleibt nichts liegen (gemessen) | geprüft, ohne Befund (Ausnahme F-4) |
| Reproduktion laut Slice: alpine entfernt, Vorher-Stand rot mit dem CI-Fehlertext, Nachher-Stand grün (ungecacht und gecacht) | geprüft, ohne Befund; Slice-Ursachenkette bestätigt |
| Weitere Stellen im Paket `internal/e2e` mit Container-ID oder -Name aus `docker run -d`: `sandbox_devcontainer_docker_test.go` (`run -d --name`, Name selbst gesetzt, gebautes Image), `docker pull -q alpine:3` im Container, Skripte mit `FROM alpine:3`, `down_volumes_docker_test.go` (`.Output()`, getrennt), `sandbox_devcontainer_docker_test.go` `docker info --format '{{.ID}}'` (kein Pull, nur Vergleich) | geprüft, ohne Befund |
| `internal/adapter/driven/docker` Tests: `CombinedOutput` nicht vorhanden, `.Output()` (nur stdout) in `engine_psjsonschema_docker_test.go`; kein Registry-Image-Start mit gelesener ID | geprüft, ohne Befund |
| `container_id_test.go`: untagged, nur Testcode, läuft in `make test` ohne Engine (gemessen: `internal/e2e` ok); Build-Tag-Konvention (Abschnitt Tests der Architektur-Sicht) verlangt den Tag für Engine-Tests, diese Datei greift nicht auf die Engine zu; Tests nicht tautologisch (Mutationen erkannt) | geprüft, ohne Befund (Hinweis F-2) |
| Sprachvertrag ([`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)): Fehlermeldungen und Kommentare im Diff englisch, Plan deutsch | geprüft, ohne Befund |
| `AGENTS.md` §3.2 (kein `//nolint`), §3.6 (kein Gate, keine Schwelle, kein Workflow berührt; `repo.mk`, `Dockerfile`, `integration.yml` nicht im Diff) | geprüft, ohne Befund |
| `AGENTS.md` §3.9: (a) keine Chronik-Kommentare, (b) keine Symbole mit geändertem Namen in Kommentaren (Suche nach `container name` / altem Rückgabewert-Wort: der Doc-Kommentar wurde auf „container ID“ umgestellt), (c) `containerIDFromRunOutput` steht vor seinem Test, kein Helfer zwischen Doc-Kommentar und Funktion | geprüft, ohne Befund |
| §3.10 „unverändert“-Zeilen: `sandbox_devcontainer_docker_test.go`, `repo.mk`, `Dockerfile`, `integration.yml` nicht im Diff; `harness/replay.md` und `docs/maintainer/quality.md` enthalten keinen der Suchbegriffe (je 0 Treffer, selbst gesucht) | geprüft, ohne Befund |
| Roadmap-Zeile: nennt Zustand und auflösbaren Anker (Welle `v0.11.0`, Slice-Link), keine Chronik | geprüft, ohne Befund |
| Carveouts, ADR-Immutabilität, Referenzmodell, Exit-Codes, Dateisicherheit | nicht berührt (reiner Testcode und Planungsdokumente) |

Nicht geprüft:

- Vollständiger `make test-docker`, `make gates`, `make coverage-gate`, `make docs-check` über den Slice, `make doc-trace` — Verifier-Aufgabe (DoD); Review-Läufe: `make test`, `make lint`, Einzeltest. `make verify-depguard` — keine Importregel berührt.
- CI-Beleg `integration-docker` auf einem Runner — erst nach dem Push mit Freigabe des Projektinhabers (Slice §6, offen).
- [`slice-v1-podman-aktuelle-version`](../plan/planning/in-progress/slice-v1-podman-aktuelle-version.md) inhaltlich — anderer Vorgang, nicht Gegenstand dieses Slice.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Plan-Tabelle gegen diff --stat unvollständig · Mutant überlebt (redundante Normalisierung) · Netzabhängigkeit durch Vorab-Pull · Cleanup erst nach ID-Prüfung registriert · Kommentar ungenau zum Ist-Zustand

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM; das LOW (F-1) ist eine Tabellenberichtigung, die vier INFO ohne Änderungspflicht.

**Übergabe:** F-1 an den Implementer (Plan-Tabelle §3 berichtigen, dann Closure §7); die Finding-Klassen gehen in die Slice-Closure §7 und den Zähler. Verifier: DoD-Punkte `make test-docker`, `make gates`, `make doc-trace` und der CI-Beleg bleiben offen.
