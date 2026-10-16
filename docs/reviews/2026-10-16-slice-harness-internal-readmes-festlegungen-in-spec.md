# Review-Report: slice-harness-internal-readmes-festlegungen-in-spec — 2026-10-16

**Review-Art:** Code (Doku-/Spec-Diff) — geprüft gegen Plan, Spezifikation, Code-Ist-Stand und Hard Rules.

**Gegenstand:** Commit-Range `145fc42..db27b3a` (zehn Commits): Inventur im Slice, neue `SPEC-036`..`SPEC-042` in `spec/spezifikation.md`, Korrektur `spec/architecture.md` §2.1, sieben gekürzte `internal/**/README.md`, Modus-Zeile in `harness/conventions.md`. <!-- d-check:ignore (Kennungen als Zitat) -->

**Skill:** `.harness/skills/reviewer.md` @ 3ec4535 · <!-- d-check:ignore (Skill-Pfad als Zitat) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingangs-Kontext:**

- `slice-harness-internal-readmes-festlegungen-in-spec` (Plan, DoD, §9 Inventur) <!-- d-check:ignore (Kennungen als Zitat) -->
- `LH-LESE-003`, `LH-FA-UP-001`, `LH-FA-UP-003`, `LH-FA-UP-004`, `LH-FA-CLI-005`, `LH-FA-CLI-005A`, `LH-FA-GEN-001`..`LH-FA-GEN-005`, `LH-AK-007`, `LH-FA-CONF-002`, `LH-FA-CONF-005`, `LH-FA-DIAG-002`, `LH-AK-001` <!-- d-check:ignore (Kennungen als Zitat) -->
- `AGENTS.md` (Hard Rules §3.4, §3.7), `harness/conventions.md` (Modus-Deklaration, `MR-006`)
- Spezifikations-Kopf (Abschnitt „Kennungen") und `v6.18.1` · `regelwerk/grundlagen-source-precedence.md` §ID-Schema
- Code-Ist-Stand: `internal/hexagon/application/{upservice,upservice_portparse,downservice,generate,config,doctor,devcontainer_sandbox,devcontainer_sandbox_config,devcontainer_egress}.go`, `internal/hexagon/domain/{serviceup,configpath,sandbox,diagnostic}.go`, `internal/adapter/driven/{docker/engine.go,netprobe}`, `internal/adapter/driving/cli/{root,cli,up}.go`

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Spec sagt, ein leerer `## [Unreleased]`-Abschnitt werde eingefügt, wenn er „außerhalb des Bereichs" fehlt. Der Code prüft die ganze Datei (außer Codeblöcke), also auch den verwalteten Bereich; die Standardvorlage trägt `## [Unreleased]` im Bereich (Vorlage Zeile 8). Ein nutzerbearbeiteter Bereich, der `## [Unreleased]` behält, wird nie repariert (Test `TestGenerateChangelog_UserEditedBlock_WithUnreleased_NoOp`). Ebenso wird nur nach dem END-Marker nach einem Release-Abschnitt gesucht, die Spec sagt nur „gibt es einen Release-Abschnitt". | `LH-AK-007`, Spec-Kopf („präzisieren, nie erweitern"), Slice §6 („Werte kommen aus dem Code") | `spec/spezifikation.md` · "Fehlt außerhalb des Bereichs der Abschnitt `## [Unreleased]`" | ja — `make test` (der genannte Test pinnt das Code-Verhalten); Spec-Seite Review-only | Spec-Bedingung weicht vom Code ab | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-2 | LOW | `SPEC-037` verfeinert genau eine Anforderung (`LH-FA-UP-001`) und trägt dennoch `SPEC-<NNN>`. Der Spec-Kopf und das Baseline-ID-Schema reservieren `SPEC-*` für Festlegungen ohne einzelne Anforderung; eine Verfeinerung genau einer ID trägt `<LH-ID>.<Buchstabe>`. Die übrigen sechs Einträge nennen je zwei oder mehr Anforderungen. | Spec-Kopf „Kennungen"; `v6.18.1` · `regelwerk/grundlagen-source-precedence.md` §ID-Schema | `spec/spezifikation.md` · "### SPEC-037 — Erkennung probebarer Ports in `compose.yaml`" | nein — Review-only (kein Gate prüft die Form) | ID-Form folgt nicht der Kopf-Regel | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-3 | INFO | Die Kopf-Regel regelt nur die Verfeinerung einer einzelnen Anforderung; eine Festlegung, die zwei bis sechs Anforderungen verfeinert (`SPEC-036`, `SPEC-038`..`SPEC-042`), ist nicht beschrieben. Die Wahl `SPEC-<NNN>` ist die naheliegende Lesart, aber nirgends deklariert. Ein Verzeichnis oder Index der `SPEC-*` existiert nicht; es musste nichts nachgezogen werden. | Spec-Kopf „Kennungen" | `spec/spezifikation.md` · "Eine Verfeinerung einer einzelnen Anforderung trägt deren Kennung mit" | nein — Review-only | Kopf-Regel deckt Mehr-Anforderungs-Verfeinerung nicht ab | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-4 | LOW | Zwei Inventur-Zeilen nennen ein Ziel, das die Aussage nicht trägt. (a) `Clock.Sleep` ohne Kontext ist als Festlegung mit Ziel `SPEC-042` und Architektur §7.1 geführt; beide enthalten sie nicht, der Code-Kommentar von `Clock.Sleep` nennt den Grund auch nicht, das alte README nannte ihn. (b) `init --devcontainer` setzt `devcontainer.enabled: true` ist mit Ziel `LH-AK-005` als „A" geführt; `LH-AK-005` nennt den Wert nicht. | Slice-DoD „Inventur … Ziel" | `internal/hexagon/port/driven/clock.go` · "Sleep pauses the calling goroutine for at least"; Slice §9 · "`Clock.Sleep` ohne Kontext" | nein — Review-only | Inventur-Ziel trägt Aussage nicht | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-5 | INFO | Die Abweichung `LH-FA-CLI-005` (letzte Option gewinnt, `--quiet` nur Fehler) gegen `root.go` (feste Rangfolge `--debug` > `--verbose` > `--quiet`, `--quiet` setzt `slog.LevelWarn`) ist richtig beschrieben und zu Recht nicht in die Spec übernommen. Sie steht nur in Slice §9, nicht in §6 mit Ausgang und in keinem Register; die Formulierung „hebt nur den Pegel auf Warn" lässt aus, dass `--quiet` zusätzlich Statustabelle (`up`), Erfolgsmeldung (`down`) und OK-Items (`doctor`) unterdrückt. Zuständig für die Entscheidung: Projektinhaber/Architect. | `LH-FA-CLI-005` | `docs/plan/planning/in-progress/slice-harness-internal-readmes-festlegungen-in-spec.md` · "offener Befund für den Projektinhaber" | nein — Review-only | Lastenheft-Code-Abweichung ohne dauerhaften Ort | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-6 | INFO | `SPEC-042` „Aufrufformen" nennt `up`, `down`, `ps --format json`, `logs`; der Adapter ruft zusätzlich `config --format json`, `volume ls`, `volume rm` und `--dry-run up -d`. Der Satz „Scheitert erst der Compose-Aufruf selbst, gilt Exit-Code `12`" gilt für `logs` bei Abbruch des Kontexts nicht (Kontextfehler wird unverpackt zurückgegeben, Exit `0`). | `LH-FA-CLI-006` | `spec/spezifikation.md` · "**Aufrufformen:** `docker compose -f <Projektverzeichnis>/compose.yaml`" | nein — Review-only | Aufzählung unvollständig | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-7 | INFO | `SPEC-036` sagt, die `warn`-Diagnosen (`up.state.<service>.unknown`, `up.port.<service>.<index>`, `up.port.<service>.unreachable`) würden „ausgegeben". Im Code tragen nur die Antwort bei erfolgreicher Stabilisierung (`buildResult(..., diagnostics)`) die Diagnosen; Fehlerpfad (`ErrComposeRuntime`) und Timeout liefern nur `PartialServices`. | `LH-FA-UP-003` | `spec/spezifikation.md` · "Je Service und Aufruf wird einmalig eine `warn`-Diagnose `up.state.<service>.unknown` ausgegeben" | nein — Review-only | Diagnose-Reichweite nicht eingegrenzt | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-8 | INFO | Die Zelle der Sub-Area `internal/**/README.md` trägt „Seit … beschreiben sie den Code-Bestand nicht mehr". Zustand (Zweck und Verweis) und Beleg (Slice-Verweis) stehen darin; der Modus „Greenfield (graduiert 2026-07-25)" ist unverändert. Keine Änderungspflicht. | `AGENTS.md` §3.7 | `harness/conventions.md` · "beschreiben sie den Code-Bestand nicht mehr" | nein — Review-only | Zustandsfeld mit Zeitform (unkritisch) |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `SPEC-036` Takt `500 ms`/`300 ms` (`pollInterval`, `dialTimeout`), Schwelle `3` (`RestartLoopThreshold`), Zustandsklassen-Tabelle (`ParseContainerState`), Entscheidungsmatrix, Healthcheck-Dominanz, alphabetische Reihenfolge, Exit `12` bei Fehlschlag und Timeout, `--timeout=0` | geprüft, ohne Befund (außer F-7) | <!-- d-check:ignore (Kennungen als Zitat) -->
| `SPEC-037` acht Portformen, Protokoll-Suffix, Bereich, Langform, IPv6 → nicht probebar, Bereich `1`..`65535`, Normalisierung `""`/`0.0.0.0`/`::`, einzelne Portzahl als Host-Port | geprüft, ohne Befund (`upservice_portparse.go`) | <!-- d-check:ignore (Kennungen als Zitat) -->
| `SPEC-038` Wahrheitstabelle `down --volumes`, Exit `10`, `--json` → Ablehnung (`SilenceConfirmer`) | geprüft, ohne Befund (`downservice.go`) | <!-- d-check:ignore (Kennungen als Zitat) -->
| `SPEC-039` Dateien und Markierungsstile, Zustandstabelle, LF-Normalisierung, Zwei-Phasen-Devcontainer, `forwardPorts`-Quelle, Exit `10` bei nicht parsebarer `compose.yaml` | geprüft, ohne Befund (außer F-1) | <!-- d-check:ignore (Kennungen als Zitat) -->
| `SPEC-040` 14 Konfigurationspfade (Zählung gegen `ConfigPathKind`), 7-Stufen-Pipeline (gegen `Set`), Repository-Zeichenmenge (`cloneURLSafeRE`), UID-Bereich, Listenpfade mit Neuserialisierung, Orphan-Feature-Warnung | geprüft, ohne Befund | <!-- d-check:ignore (Kennungen als Zitat) -->
| `SPEC-041` Schweregrad-Tabelle `compose.yaml.valid` | geprüft, ohne Befund (`checkComposeYaml`, `loadProjectYAML`) | <!-- d-check:ignore (Kennungen als Zitat) -->
| `SPEC-042` Preflight-Reihenfolge, Exit `11`/`12`, Fire-and-forget, Netprobe-Kontext | geprüft, ohne Befund (außer F-6) | <!-- d-check:ignore (Kennungen als Zitat) -->
| Referenzrichtung (Decken-Regel) | geprüft, ohne Befund: keine neue Zeile in `spec/` verweist auf Slice, ADR, Carveout oder Roadmap; `spec/architecture.md` trägt im Diff nur die Zeile §2.1 ohne ADR-, Slice- oder Meilenstein-Bezug (`AGENTS.md` §3.4) |
| Architektur §2.1 (`SeverityInfo`, streng monoton) | geprüft, ohne Befund (`diagnostic.go`) |
| Vollständigkeit der Verlagerung: alle sieben alten README-Fassungen gegen Slice §9 und aktuellen Stand | geprüft; verlorene Festlegungen nur in F-4; übrige Streichungen sind in Spec, Architektur (§2.2, §2.4, §2.6, §7.1), `carveouts.md`, `.golangci.yml` oder am Code belegt |
| Inventur-Abweichungen README ↔ Code (16 Doctor-Kennungen, 14 Pfad-Arten, 8 `DockerEngine`-Methoden, 14 Adapter-Verzeichnisse, fünf Devcontainer-Dateien) | geprüft, ohne Befund |
| Kommentar-Klassen und Zustandsfelder in neuen Spec-Einträgen und READMEs (`AGENTS.md` §3.7) | geprüft, ohne Befund: Indikativ, keine verworfene Alternative, „Grenze"/„Folge" tragen zulässige Klassen |
| Sprache/ASCII-Konvention | geprüft, ohne Befund: Spec deutsch mit Umlauten (Konvention der Datei), `harness/conventions.md`-Zeile ASCII |
| `harness/conventions.md` Zeile `internal/**/README.md` | geprüft: Modus unverändert, neue Begründung trägt (F-8 nur INFO) |
| `make docs-check` | gelaufen: 246 Dateien, 0 Befunde (vor dem Report; nach dem Report siehe Verdikt) |

Nicht geprüft:

- `make gates`, `make test`, `make test-docker`, `make verify-depguard` — kein Code- und kein Importregel-Diff; Implementer meldet `make gates` Exit 0. Das F-1-Testverhalten wurde durch Lesen von `generate_test.go` belegt, nicht ausgeführt.
- DoD-Haken und Closure-Notiz des Slice — Verifier-Aufgabe.
- `SPEC-039` Sandbox-Dateien im Detail (Bedingungen `nestedRuntime`/Egress) — nur gegen `planDevcontainerFiles` gelesen, keine Vorlageninhalte. <!-- d-check:ignore (Kennungen als Zitat) -->

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Spec-Bedingung weicht vom Code ab · ID-Form folgt nicht der Kopf-Regel · Kopf-Regel deckt Mehr-Anforderungs-Verfeinerung nicht ab · Inventur-Ziel trägt Aussage nicht · Lastenheft-Code-Abweichung ohne dauerhaften Ort · Aufzählung unvollständig · Diagnose-Reichweite nicht eingegrenzt · Zustandsfeld mit Zeitform (unkritisch)

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM) ist ein Spec-Fehler mit Verhaltenswirkung: Wer die Spec umsetzt oder gegen sie testet, erwartet bei einem nutzerbearbeiteten Bereich mit `## [Unreleased]` eine Reparatur, die der Code (und sein Test) nicht ausführt; eine Angleichung des Codes an die Spec fügte einen doppelten Abschnitt in Nutzer-Changelogs ein. LOW und INFO blockieren nicht.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. F-5 braucht eine Entscheidung des Projektinhabers über Architect (Lastenheft ändern oder Code ändern); F-3 ist eine Frage an den Architect (Kopf-Regel). Dieser Report ersetzt keine Verifikation.
