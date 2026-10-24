# Review — slice-v1-generate-prueft-freigabeliste

* Rolle: Reviewer (`.harness/skills/reviewer.md`), frischer Kontext
* Gegenstand: `git diff 750139e5..0cdec3b2`, Slice `docs/plan/planning/in-progress/slice-v1-generate-prueft-freigabeliste.md`
* Bezug: LH-FA-DEV-003, LH-NFA-SEC-004, LH-FA-CLI-006
* Sensor-Lauf: `make gates` selbst gefahren, Exit-Code 0. Gezielte Tests (`Freigabe|ThreeCallers|SandboxHostSocket|SandboxFeatureWarn|Homebrew|ClaudeCodeNode` in `internal/hexagon/application`, `UnlistedFeature|SandboxFeatureWarnings` in `cmd/uboot`) im Image `u-boot:test-docker-tools`: grün. Nicht gelaufen: `make test-docker` (Anweisung: Speicher; kein Docker-Pfad berührt), `make verify-depguard` (keine Importregel berührt, kein neuer Import im Diff).

## Findings

Kein HIGH, kein MEDIUM.

### LOW-1
* kategorie: LOW
* quelle: AGENTS.md §3.10 (Bedingung in Handbuch/Spec braucht einen Test)
* pfad: `docs/user/benutzerhandbuch.md` · „Features mit `enabled: false` oder ohne `enabled` werden weder erzeugt noch geprüft; `doctor` meldet sie weiter.“
* klasse: aussage-ueber-code-ohne-pin
* befund: Der Teilsatz „doctor meldet sie weiter“ stützt sich auf `classifyFeatureEntries` (`doctor.go`). Die Doctor-Tests (`doctor_features_test.go`, Fälle mit `uninvited.test`) pinnen nur `enabled: true`; `T5` (`TestFeatureSourceAllow_ThreeCallersAgree`) ebenso. Ein Fall `enabled: false` + nicht freigegebene Quelle am Doctor fehlt; der Slice trägt selbst „nicht neu gepinnt“ ein. Für „ohne `enabled`“ meldet Doctor außerdem nicht die Freigabe-Verletzung, sondern `enabledKeyMissing` – der Satz ist dort ungenau.
* risiko: Ändert jemand `classifyFeatureEntries` (z. B. deaktivierte Einträge überspringen), widerspricht das Handbuch dem Code, kein Gate fängt es.
* verifizierbar: `make test` (mit neuem Pin) – bzw. Satz auf `enabled: false` einengen.

### LOW-2
* kategorie: LOW
* quelle: AGENTS.md §3.10; Spec SPEC-050 „Freigabe der Quelle (beide Profile)“
* pfad: `internal/hexagon/application/generate_freigabe_test.go` · `TestGenerateDevcontainer_FreigabeBeforeSandboxRejection`
* klasse: bedingung-ohne-eigenen-pin
* befund: „Beide Profile“ ist nur indirekt gepinnt. Mutation `len(unlisted) > 0 && !sandbox` (Prüfung im Sandbox-Profil ausgesetzt) wird allein von T6 erkannt, und T6 trägt zugleich `docker-cli`. Ein Test „Sandbox-Profil, nur unfreigegebene Quelle, kein `docker-cli`“ fehlt; der Binary-Test läuft nur im Standardprofil. Die umgestellten Sandbox-Fixtures belegen nur den Positivfall (Quelle gelistet).
* risiko: Eine Mutation, die das Sandbox-Profil ausnimmt und T6 nicht berührt, überlebt; das wäre genau die Lücke, die der Slice schließt, im sichereren Profil.
* verifizierbar: `make test`.

### LOW-3
* kategorie: LOW
* quelle: Maintainability (Test deckt Zusage der Meldung)
* pfad: `internal/hexagon/application/devcontainer_features.go` · `sortedFeatureNames(cfg.Devcontainer.Features)` in `unlistedFeatureSources`
* klasse: zusage-ohne-pin
* befund: Die Doc-Zusage „sorted by feature name“ und das Verbinden mehrerer Verstöße mit `"; "` haben keinen Test (keiner der Tests T1–T8 hat zwei unfreigegebene Einträge).
* risiko: Reihenfolge der Meldung (Determinismus, Map-Iteration) kann unbemerkt kippen.
* verifizierbar: `make test`.

### INFO-1 (Sicherheitseinschätzung, Advisory)
Kein Advisory-Kandidat im Sinne einer Vertrauensgrenzen-Verletzung. Begründung aus dem Code: `devcontainer.featureSources.allow` und `devcontainer.features.<n>.source` stehen in derselben Datei `u-boot.yaml` (`cfg.Devcontainer.FeatureSources.Allow` neben `Features[...]`), und `generate` liest sie ohne weitere Signatur (`readProjectConfig`). Wer `u-boot.yaml` kontrollieren kann (Schreibrecht im Projekt, PR/Fremd-Repository vor dem Klonen), kann Quelle und Freigabe zugleich eintragen – oder `devcontainer.json` direkt ändern. Die Prüfung schützt also vor Versehen und vor dem Umweg über u-boot-Befehle (`config set` verweigert schon, `doctor` meldet), nicht vor einem Angreifer mit Schreibrecht auf `u-boot.yaml`. Es ist ein Vertragsfehler (LH-FA-DEV-003 / LH-NFA-SEC-004: Freigabe als expliziter Schritt wurde von `generate` still umgangen, Exit 0) und als Verhaltensänderung im CHANGELOG richtig unter Fixed verbucht. Einzuschätzen bleibt die Wirkung für Nutzer, die einem fremden `u-boot.yaml` mit vorbereiteter Liste vertrauen: unverändert. Die Einstufung als Advisory liegt beim Projektinhaber; aus dem Code folgt eher „Hardening/Bugfix“ (niedrige Schwere), nicht „Vulnerability“.

### INFO-2
Sicherheits-Randfälle am Vergleich (`featureSourceInAllow`, `allowed == src`), selbst gelesen und durch Mutation geprüft:
* Kein Präfix-Bypass: exakter String-Vergleich. Mutation M5 (`strings.HasPrefix(src, allowed)`) wird von `TestFeatureSourceAllow_ThreeCallersAgree` (Fall „trailing slash in source“) erkannt.
* Trailing Slash, Host-Groß-/Kleinschreibung: strikt ungleich, in T5 gepinnt (drei Aufrufer gleich).
* Führende/folgende Leerzeichen: Allow-Einträge werden von `normaliseFeatureSources` getrimmt, `entry.Source` nicht; ein Source mit Leerzeichen passt daher nicht zu einer getrimmten Allow – fail-closed, kein Bypass. Kein Test.
* Query/Userinfo (`https://u@h/p?x`): `validateFeatureSource` prüft nur Schema und Host; der Vergleich ist exakt, ein abweichender Query-/Userinfo-Anteil wird abgewiesen. Kein Test.
* `enabled: false` und fehlendes `enabled`: nicht geprüft, nicht gerendert (`collectDevcontainerFeatures` filtert gleich); Test T4. Katalogschlüssel ohne `source`: frei, gepinnt.
* Flag-Quelle: wird in `unlistedFeatureSources` zur Liste addiert und durchläuft `normaliseFeatureSources` (Formatprüfung; ungültiger Eintrag -> `ErrConfigValueInvalid`), nichts wird in der Prüfung geschrieben; Persistierung wie bisher nach Erfolg; `--dry-run` schreibt nichts (Unit T3-DryRun plus Binary-Test mit Byte-Vergleich).

## Mutationen (selbst nachgefahren, Rücknahme per `git checkout`, Baum danach sauber)

| Mutation | Erkannt durch |
|---|---|
| M1a: `if len(unlisted) > 0` -> `if false && …` | T1, T6, T8, Unrelated-Flag, ThreeCallers, Binary-Test |
| M3: Flag nicht in die Freigabe (`nil` statt `req.AllowExternalFeatureSources`) | `FreigabeByFlagOfSameCall`, `FreigabeByFlagDryRun`, Binary-Test |
| M4: `enabled`-Filter in `unlistedFeatureSources` entfernt | `FreigabeIgnoresInactiveEntries`, Binary-Test (Abhilfe 3) |
| M5: Vergleich als Präfix | `ThreeCallersAgree` |
| M7b: Prüfung im Sandbox-Profil ausgesetzt | nur `FreigabeBeforeSandboxRejection` (siehe LOW-2) |

M7b wurde als Abwandlung gefahren (Reihenfolge-Verschiebung per Skript nicht zuverlässig). Die Fixture-Berichtigungen (Homebrew, ClaudeCodeNodeNotice, SandboxHostSocketFeature, `RunGenerateSandboxFeatureWarnings`) ergänzen jeweils die Liste um die tatsächliche Quelle (byte-gleich); der Vergleich ist nicht gelockert (`config.go` nur Doc-Kommentar).

## Negativbefund (geprüft, ohne Befund)

* Dual-Classifier: `isConfigValidationError` (`cli.go`, `ErrConfigValueInvalid` -> Exit 10) und `mapGenerateErrorToDiagnostic` (`cli/generate.go`, -> `LH-FA-DEV-003`) selbst gelesen; Exit-Code-Pin und Envelope-Pin im Binary-Test (`TestRun_Generate_UnlistedFeatureSourceIsRejected`: Exit 10, genau ein Diagnostic `error` / `LH-FA-DEV-003`). Kein neuer Sentinel.
* Architekturgrenze: Änderungen nur in `internal/hexagon/application`; keine neuen Imports, kein I/O in der Prüfung, kein `//nolint`.
* Eine Vergleichsstelle: `featureSourceInAllow` mit drei Aufrufern (`config.go`, `doctor.go`, `devcontainer_features.go`); `unlistedFeatureSources` baut nur die gemergte Liste, keine zweite Vergleichslogik.
* Reihenfolge in `generateDevcontainer`: nach `collectDevcontainerFeatures`, vor Sandbox-Abweisung und Plan/Schreiben; Freigabe-Code gewinnt bei Doppelverstoß (T6).
* Sprachvertrag: Fehlermeldung Englisch.
* Kommentare §3.7/§3.9: neue Kommentare im Indikativ, ohne Chronik; die ersetzten Kommentare („T3's renderer does not enforce…“, „T4 enforces…“) sind auf den Ist-Ort umgestellt. Verweise auf alte Namen geprüft (kein Rest in der Berührungszone).
* Spec: SPEC-050 „Freigabe der Quelle“ und Bullet `LH-FA-DEV-003.a` zeigen aufwärts/innerhalb der Spec, keine Slice-/ADR-Verweise; Historie-Zeile ergänzt, ältere Einträge unverändert. Jede Bedingung (Liste, Flag desselben Aufrufs, `enabled`, Katalogschlüssel, Reihenfolge vor `docker-cli`, kein Schreiben auch `--dry-run`, Meldungsteile) hat einen Test (T1–T8, Binary-Test); Ausnahmen siehe LOW-1/2.
* Handbuch: kennungsfrei (§5 Regel 6); CHANGELOG `[Unreleased]` / Fixed mit Kennungen (zulässig), Abhilfe genannt. `harness/replay.md`: neue Zeile, bestehende Zeilen unverändert.
* Plan-Tabelle gegen `git diff --stat` (13 Dateien): alle Dateien im Slice-Ist-Stand aufgeführt, kein Eintrag ohne Diff; `README*`, `architecture.md`, `doctor.go`, `cli*.go` unverändert (Diff-Stat bestätigt).
* Lastenheft und ADR im Diff unverändert; Tag `v0.15.0` (`3377aa80`) unverändert; nichts gepusht, kein Tag.

## Nicht geprüft

* `make test-docker`, `make verify-depguard`, `make ci` (Begründung oben). Das Verhalten am gebauten Binary (Slice-Aussage „Wirkung am Binary“) nicht nachgefahren; Ersatz ist der Binary-nahe Test mit echten Adaptern in `cmd/uboot`.
* DoD-Closure: Verifier-Aufgabe.
