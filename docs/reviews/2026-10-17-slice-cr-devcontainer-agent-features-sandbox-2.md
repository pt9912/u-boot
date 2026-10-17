# Review-Report: slice-cr-devcontainer-agent-features-sandbox (Nachlauf) — 2026-10-17

**Review-Art:** Plan/Design — Nachlauf-Review der Erweiterung des CR-Kerns gegen Lastenheft, Hard Rules und die älteren Zusagen; Maintainability des Vertragstextes, keine DoD-Prüfung.

**Gegenstand:** Diff `5de388c..4df2b54` (`spec/lastenheft.md`, `CHANGELOG.md`) für den Slice [`slice-cr-devcontainer-agent-features-sandbox`](../plan/planning/done/slice-cr-devcontainer-agent-features-sandbox.md). Der erste Review (`2026-10-17-slice-cr-devcontainer-agent-features-sandbox.md`) lief vor dieser Erweiterung. <!-- d-check:ignore (Report-Dateiname als Zitat) -->

**Skill:** `.harness/skills/reviewer.md` (aktueller Stand im Arbeitsbaum), HEAD `48b2537` · <!-- d-check:ignore (Skill-Pfad als Zitat) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

**Eingangs-Kontext:**

- Slice [`slice-cr-devcontainer-agent-features-sandbox`](../plan/planning/done/slice-cr-devcontainer-agent-features-sandbox.md) (§1 Entscheidung 5, §6 Risiken) und der offene Code-Slice [`slice-v1-cli-quiet-sicherheitswarnungen`](../plan/planning/in-progress/slice-v1-cli-quiet-sicherheitswarnungen.md) (§1 Einstufungstabelle)
- [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DIAG-003`](../../spec/lastenheft.md#lh-fa-diag-003--fehlerklassifikation)
- Code-Bestand zum Abgleich (gelesen, nicht geändert): `internal/hexagon/application/doctor_sandbox.go` (`checkSandboxRuntime`, `checkSandboxEgress`, `checkSandboxCredentials`), `internal/adapter/driving/cli/doctor.go` (`DegradationFallback`, `--strict`) <!-- d-check:ignore (Pfade als Zitat) -->
- `AGENTS.md` (Hard Rules §3.4, §3.7, §3.9)
- Entscheidung des Projektinhabers 2026-10-17 (aus dem Review-Auftrag): Sicherheitswarnungen bleiben unter `--quiet` sichtbar; Drittanbieter-Code-Warnung beim Image-Build nicht; Egress „has no effect“ sichtbar; `nestedRuntime` „has no effect“ nicht sichtbar

**Systematische Suche gegen ältere Zusagen** (Lastenheft; Begriffe `doctor`, `--quiet`, `sichtbar`, `Warn`, `strict`, `LH-FA-DEV-008`, `LH-FA-DEV-009`, `LH-FA-CLI-005`, `LH-FA-CLI-005A`): Treffer im Normtext in `LH-FA-CLI-005` (Absätze zu `doctor`, Ausnahmen 1–4, „Alle anderen Hinweise und Warnungen entfallen“), `LH-FA-DIAG-003` (Stufen, `--strict`), `LH-FA-DEV-006` bis `-009`, `LH-FA-CLI-006`, §16. Kein Widerspruch zu einer älteren Zusage; die einzige unscharfe Stelle ist F-1. <!-- d-check:ignore (Kennungen als Zitat) -->

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Wortlaut trägt die Unterscheidung der Entscheidung nicht. Das Hauptkriterium „die eine Sicherheitseigenschaft des Sandbox-Profils betreffen“ trifft die beiden „has no effect“-Warnungen (egress, nestedRuntime) gleich: beide entstehen genau dann, wenn das Profil nicht `sandbox` ist. Die Klammer ist ohne „z. B.“/„insbesondere“ gesetzt (liest sich abschließend), nennt „Egress-Restriktion“ als Gegenstand (deckt „has no effect“ des Egress-Schlüssels, wenn man den Gegenstand statt der Lockerung liest) und „Sicherheitslockerungen des Docker-Daemons“ (deckt „has no effect“ von `nestedRuntime` nicht, auch nicht bei `docker`, weil eine wirkungslose Einstellung nichts lockert). Die Unterscheidung „Egress sichtbar, Runtime nicht“ ergibt sich nur aus dieser Wortwahl, nicht aus einem Satz. Failure-Szenario: Der Code-Slice liest den Wortlaut und stuft „egress has no effect“ als nicht sichtbar ein; seine Einstufungstabelle tut das heute („nein: Konfigurationsmeldung, die Einstellung wirkt nicht“, `keines (entfällt unter --quiet)`). Dann liefert er eine Ausgabe, die der Entscheidung des Projektinhabers widerspricht, und der Verifier hat im Lastenheft keinen Satz, gegen den er das beanstanden kann; im umgekehrten Fall zeigt ein Implementer die `nestedRuntime`-Meldung zu `docker` („nested Docker daemon“-Bezug) und bricht „Alle anderen Hinweise und Warnungen entfallen“. | [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | `spec/lastenheft.md` · "die Warn-Einträge von `u-boot doctor`, die eine Sicherheitseigenschaft des Sandbox-Profils betreffen" | nein — Review-only (Vertragsprosa; erst die Tests des Code-Slice pinnen es: `has no effect` mit `egress` sichtbar, mit `nestedRuntime` nicht) | Ausnahme-Klammer ohne Abgrenzungssatz für wirkungslose Einstellungen |
| F-2 | INFO | Der offene Code-Slice führt `egress.enabled` ohne Sandbox-Profil („has no effect“) als „nein … entfällt unter `--quiet`“ und markiert die Abweichung vom Projektinhaber als „weiter offen bis zum Review“. Die Entscheidung vom 2026-10-17 (Egress „has no effect“ sichtbar) ist dort nicht nachgezogen; der Plan widerspricht dem Vertragsziel. Zuständigkeit: Implementer des Plans, vor dem Start des Code-Slice; kein Teil des Diffs. | [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | `docs/plan/planning/in-progress/slice-v1-cli-quiet-sicherheitswarnungen.md` · "`egress.enabled` ohne Sandbox-Profil („has no effect“)" | nein — Review-only | Plan-Einstufung widerspricht Projektinhaber-Entscheidung | <!-- d-check:ignore (Kennungen und Pfad als Zitat) -->
| F-3 | INFO | Die Klammer „Sicherheitslockerungen des Docker-Daemons nach `LH-FA-DEV-007`“ benennt heute keinen Warn-Eintrag von `doctor`: der Eintrag zu `nestedRuntime: docker` hat die Stufe OK (`doctor_sandbox.go`, „Nested Docker daemon configured; it needs `--privileged`“). Bestätigt: leere Menge, kein Vertragsloch, denn die Lockerung selbst ist durch den ersten Teil der Ausnahme 1 („in der Befehlsausgabe“, [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer): „weist das als Lockerung aus“) abgedeckt, und die einzige Warnung mit Docker-Daemon-Bezug (Egress zusammen mit `docker`) ist über die Egress-Nennung sichtbar. Restrisiko als Szenario: Ein Implementer, der die Klammer wörtlich nimmt, hebt den OK-Eintrag auf `warn`, damit er unter `--quiet` erscheint; das ändert `--strict` und die JSON-Ausgabe ([`LH-FA-DIAG-003`](../../spec/lastenheft.md#lh-fa-diag-003--fehlerklassifikation)). Der Wortlaut „Warn-Einträge“ verlangt das nicht, der Code-Slice nennt die Stufenänderung ausdrücklich als Nicht-Ziel. | [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | `spec/lastenheft.md` · "Sicherheitslockerungen des Docker-Daemons nach" | nein — Review-only (Verifier hat dieselbe leere Menge festgestellt) | Klammerglied ohne heutigen Referenten | <!-- d-check:ignore (Kennungen und Pfad als Zitat) -->
| F-4 | INFO | Wortdrift zwischen Normtext und §16-Zeile 0.11.0: der Normtext sagt „Warn-Einträge … die eine Sicherheitseigenschaft … betreffen“, die Historie „die Sicherheits-Einträge von `u-boot doctor`“. „Sicherheits-Einträge“ ließe auch OK-Einträge zu. Die Historie ist nicht normativ, der Normtext ist eindeutig in „Warn“; die Abweichung ist ohne Vertragsfolge. | Maintainability | `spec/lastenheft.md` · "nennt auch diesen Login-Dienst und die Sicherheits-Einträge von `u-boot doctor`" | nein — Review-only | Historie-Zeile weicht im Begriff vom Normtext ab |
| F-5 | INFO | Im selben Abschnitt des CHANGELOG stehen eine gelieferte Aussage („`u-boot doctor --quiet` no longer lists warn items other than a degradation fallback“) und die Vertragsaussage, dass Sicherheits-Warn-Einträge unter `--quiet` sichtbar bleiben. Beide stimmen mit Lastenheft und Stand überein, weil der Schlusssatz („states the contract only … their visibility under `--quiet` … follow in a later release“) die zweite ausdrücklich als nicht geliefert ausweist; der Slice führt den Anpassungspunkt (§Release-Hinweis, DoD des Code-Slice). Kein Widerspruch, nur Lesehinweis. | Maintainability | `CHANGELOG.md` · "This entry states the contract only" | nein — Review-only | Gelieferter Stand und Vertragsstand im selben Eintrag |

## Prüfpunkte aus dem Auftrag

1. **Satzstruktur und Eindeutigkeit.** Grammatik: „die Ausweisung von … sowie die Warn-Einträge von `u-boot doctor`, die …“ ist korrekt gebaut (zweites Objekt zu „bleiben … sichtbar“, der Relativsatz bezieht sich auf „Warn-Einträge“); die Nummerierung 1–4 und die Zahl „Vier Ausnahmen“ bleiben. Die Klammer ist nach Form abschließend (keine Beispielmarke), der Hauptsatz formuliert aber ein offenes Kriterium; beides zusammen trägt die Unterscheidung „Egress ‚no effect‘ sichtbar, Runtime ‚no effect‘ nicht“ nicht ausdrücklich: F-1 (MEDIUM). Gegenstück im Code-Bestand: [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) deckt beide Credentials-Warnungen (Klartext in Projektdatei, fehlende Quelle) als „Quelle der Zugangsdaten“, ohne Zweifel; Egress mit `docker` ist durch „Egress-Restriktion“ sichtbar. Der Degradations-`vfs`-Fall bleibt in Ausnahme 2, keine Doppelregelung.
2. **Docker-Daemon-Lockerung ohne Warn-Eintrag.** Bestätigt als leere Menge ohne Vertragsloch: F-3 (INFO).
3. **Kollision mit älteren Zusagen.** `--strict` ([`LH-FA-DIAG-003`](../../spec/lastenheft.md#lh-fa-diag-003--fehlerklassifikation)) bindet Warn-Einträge an einen Exit-Code ungleich `0`, unabhängig von der Sichtbarkeit; der Code liefert `11` ([`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes): Umgebungs-/Prüfungsfehler) und nennt laut CHANGELOG die Warn-Einträge in der Fehlermeldung, die `--quiet` zeigt. `doctor --quiet --strict` mit einem verborgenen „has no effect“ endet also mit `11` und einer Fehlermeldung, die den Eintrag benennt; die neue Ausnahme verkleinert die Menge der verborgenen Warnungen nur. Keine neue Wechselwirkung, kein Widerspruch zu [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes).
4. **Decken-Regel, Anker, Sprache, §3.7/§3.9.** Der Diff fügt im Lastenheft keinen Verweis auf Slice, Welle, Carveout, Roadmap oder ADR ein (§3.4); alle neuen Verweise sind Anker auf Lastenheft-Abschnitte und löschen weder Überschrift noch Anker (`#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer` existiert). Lastenheft deutsch, CHANGELOG englisch (Sprachvertrag). Die Sätze beschreiben den Zustand im Indikativ; kein „früher/former/no longer“ im Diff, kein Verweis auf umbenannte Symbole (§3.9 a/b), kein eingefügter Helfer zwischen Doc-Kommentar und Funktion (§3.9 c). Im CHANGELOG-Diff selbst steht kein „no longer“; das vorhandene stammt aus unverändertem Text.
5. **CHANGELOG.** „states the contract only“ und die Liste (Übergabe der Agenten-Zugangsdaten, neue `doctor`-Prüfungen und ihre Sichtbarkeit unter `--quiet`, Login-Dienst-Ausweisung, `sshd`- und Claude-Code-Features) folgt dem Lastenheft; `doctor`-Zeilen unter `--quiet` stehen als Vertrag, nicht als geliefert. Die Klammer „(credential source, egress restriction, relaxations of the Docker daemon)“ spiegelt die Klammer des Lastenhefts 1:1 und erbt F-1. Lesehinweis: F-5.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `LH-FA-CLI-005` Satzbau, Nummerierung, Zahl der Ausnahmen | geprüft, ohne Befund | <!-- d-check:ignore (Kennung als Zitat) -->
| Wechselwirkung mit `--strict` und Exit-Code-Vertrag (`LH-FA-DIAG-003`, `LH-FA-CLI-006`) | geprüft, ohne Befund | <!-- d-check:ignore (Kennungen als Zitat) -->
| Überschneidung Ausnahme 1 gegen Ausnahme 2 (Degradations-Fallbacks in `doctor`) | geprüft, ohne Befund |
| Decken-Regel (`AGENTS.md` §3.4) im Lastenheft-Diff | geprüft, ohne Befund |
| Link-Anker im Diff (Lastenheft und CHANGELOG) | geprüft, ohne Befund |
| Sprachvertrag Lastenheft/CHANGELOG | geprüft, ohne Befund |
| `AGENTS.md` §3.7/§3.9 im Diff | geprüft, ohne Befund |
| CHANGELOG gegen Lastenheft (Vertrag statt geliefert) | geprüft, ohne Befund (Lesehinweis F-5) |
| §16-Zeile 0.11.0 (Muster, Quelle, Verweise) | geprüft, ohne Befund (Begriffsdrift F-4) |

Nicht geprüft:

- Code, Tests und Handbuch (`SPEC-043`, `docs/user/`) — der Diff ändert nur Vertragstext; Umsetzung und Nachzug liegen im Code-Slice und im Spezifikations-Folge-Slice. <!-- d-check:ignore (Kennung als Zitat) -->
- `make gates`, `make test-docker`, `make verify-depguard` — kein Code- und kein Importregel-Pfad berührt; `make docs-check` läuft nach dem Schreiben dieses Reports (Ergebnis im Handoff).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Ausnahme-Klammer ohne Abgrenzungssatz für wirkungslose Einstellungen · Plan-Einstufung widerspricht Projektinhaber-Entscheidung · Klammerglied ohne heutigen Referenten · Historie-Zeile weicht im Begriff vom Normtext ab · Gelieferter Stand und Vertragsstand im selben Eintrag

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (MEDIUM): Der Wortlaut der Ausnahme 1 pinnt die Entscheidung zu den beiden „has no effect“-Meldungen nicht, und der offene Code-Slice folgt ihr heute nicht (F-2). Ein Abgrenzungssatz im Lastenheft oder die Nachführung der Entscheidung in Normtext und Plan klärt es; der Reviewer schlägt keinen Wortlaut vor.

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation.
