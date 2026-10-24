# Review-Report: slice-v1-devcontainer-ollama-cloud — 2026-10-24

**Review-Art:** Code — gegen Plan, [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) (Accepted), Lastenheft 0.23.2 und Hard Rules.

**Gegenstand:** `git diff 8ad49398..HEAD` (Commits `7f082a8a`, `5b983434`, `24711dd3`), Slice
[`slice-v1-devcontainer-ollama-cloud`](../plan/planning/done/slice-v1-devcontainer-ollama-cloud.md)
(Welle [`welle-ollama-cloud-modelle`](../plan/planning/welle-ollama-cloud-modelle.md)).

**Skill:** `.harness/skills/reviewer.md` @ `24711dd3` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-24

**Eingangs-Kontext:**

- Slice-Plan (Ist-Stand, H1 bis H15, D1/D2, F1 bis F4, §7 Implementer-Beleg mit Mutationstabelle)
- [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) (Entscheidungen 3, 4, 6, 8, O1 bis O3, Fitness Function)
- [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)
- `AGENTS.md` §3.1 bis §3.10, `harness/review.md`, `harness/replay.md`

---

## Urteil

**Kein HIGH, kein MEDIUM.** Verhalten, Meldung, Eintragstyp und Tests tragen Lastenheft und ADR. Zwei LOW und
vier INFO unten; keiner blockiert. Die Rollenpflicht des Reviewers (DoD-Punkt „Review durchgeführt“) ist mit diesem
Report erfüllt; die Closure-Pflichten gehören nicht dem Reviewer.

## Sensoren (selbst gefahren)

| Sensor | Ergebnis |
|---|---|
| `make gates` | Exit 0; `coverage-gate: OK — coverage 92.10% meets threshold 90%`; Baum danach sauber |
| Mutationen M07, M15, M22, M23 (Patch per `git apply --recount`, Rücknahme per `git apply -R`) | je rot: M07 vier Tests (`…OllamaHost_InactiveIsByteIdentical`, `…DefaultProfileRunArgs`, bestehend `…SandboxEgress_NotGenerated`, `…SandboxPodman_NotGenerated`); M22 drei Tests im Lese- und Schreibpfad (`…VersionWithoutSourceRejected`, `TestConfigSet_OllamaHostVersionWithoutSource_Rejected`, `TestRun_OllamaHostVersionWithoutSourceExits10`); M15 drei Tests (`…SandboxNoticeIsRelaxation`, `…EgressRestrictionIsAWarning`, `TestRun_Generate_OllamaHostNotices`); M23 `…SourceOverrideIsExternalFeature`. `git diff` nach jeder Rücknahme leer (0 Zeilen). Deckt sich mit der Mutationstabelle in §7 des Slice |
| D1 `TestE2E_OllamaHost_DefaultProfileResolvesHostDockerInternal` (gezielt, `-tags docker -p 1`) | PASS, 10,6 s; `host.docker.internal -> 172.17.0.1`; 0 Container danach |
| D2 `TestE2E_OllamaHost_SandboxWithEgressResolvesAndAddsNoHost` (gezielt) | PASS, 7,1 s; `hosts=`-Zeile ohne `ollama.com`, `host.docker.internal`, Adresse; 0 Container danach |
| `make test-docker` (volles Paket) | nicht gefahren: Auftrag, Speicher, Paket am Zeitlimit; gezielte Läufe ersetzen es für die neuen Tests, nicht für den Rest des Pakets |
| `make verify-depguard` | nicht gefahren: keine Importregel berührt (Diff nur in `application`, Tests, Template) |
| Mutationen M01 bis M06, M08 bis M14, M16 bis M21, M24 bis M32 | nicht selbst gefahren (Stichprobe von vier); Angaben des Implementers gelesen und gegen die Tests plausibilisiert, nicht reproduziert |
| Hinweis | Die Docker-Tests liefen mit `-v /tmp:/tmp`; ohne diesen Mount scheitert `devcontainer up` im Testimage am Bind-Pfad (Umgebung, kein Befund am Code) |

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Plan-Tabelle führt `config.go` als „unverändert“ (bedingt update), der Diff ändert dort Doc-Kommentar und die Argumente der Fehlermeldung von `revalidateFeatureEntry`; §7 *Abweichungen* nennt `config.go` nicht. Der Abgleich gegen `git diff --stat` nach der letzten Fix-Runde steht noch aus. | `AGENTS.md` §3.10 (Plan-Tabelle gegen `git diff --stat`) | `internal/hexagon/application/config.go` · `featureWithoutVersionReason(path.Feature.String())` | Review-only (Abgleich bei der Closure) | Plan-Tabelle ohne Abgleich mit dem Diff |
| F-2 | LOW | `ollamaNotices` filtert auf die Teilzeichenkette `feature ollama`, die auch die Meldung `devcontainer feature ollama-host adds …` trifft. Heute aktiviert kein Test beide Features und zählt über diesen Filter; ein künftiger Test mit `ollama` und `ollama-host` zählt zwei statt einer Meldung. | `AGENTS.md` §3.9 (b) (der neue Name ist Präfix-Treffer eines bestehenden Filters) | `internal/hexagon/application/generate_ollama_test.go` · `strings.Contains(w.Message, "feature ollama")` | `make test` (grün heute; latent) | Präfix-Überdeckung eines Test-Filters durch neuen Namen |
| F-3 | INFO | Die Meldung sagt, ein Dienst mit Loopback-Bindung sei „only through a relay or another binding that you set up“ erreichbar. Der Wortlaut folgt [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) Entscheidung 8; das Lastenheft sagt für eine vom Nutzer geänderte Bindung „keine Zusage“. Der Text formuliert eine notwendige Bedingung, keine Zusage; kein Widerspruch, aber die Nähe ist ein Lesepunkt für den Validator und den Spezifikations-Slice. | [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) Absatz `ollama-host` | `internal/hexagon/application/devcontainer_features.go` · `"is reachable only through a relay or another binding that you set up; "` | Review-only | Meldungswortlaut nahe an Zusage-Grenze |
| F-4 | INFO | Der Test `…InactiveIsByteIdentical` schaltet den Literal-Pin „kein `runArgs`-Schlüssel“ über Teilzeichenketten des Fallnamens (`egress`, `podman`) ab; ein neuer Fall mit Option, dessen Name keines der Wörter trägt, schlägt fälschlich an, einer mit dem Wort ohne Option bleibt ungeprüft. | Maintainability | `internal/hexagon/application/generate_ollama_host_test.go` · `!strings.Contains(hname, "egress") && !strings.Contains(hname, "podman")` | Review-only | Testbedingung an Fallnamen gekoppelt |
| F-5 | INFO | `hostGatewayFeatureActive` prüft `containerOption == ""` für einen Katalogeintrag, der die Option als Konstante trägt; der Zweig ist im Produktcode nicht erreichbar. | Maintainability | `internal/hexagon/application/devcontainer_features.go` · `if featureCatalogue()[featureOllamaHost].containerOption == ""` | `make coverage-gate` (Schwelle trotzdem erfüllt) | nicht erreichbarer Schutzzweig |
| F-6 | INFO | Zwischenzustand wie geplant (Carveout 1): Die neue Katalogzeile im Handbuch (`ollama-host … schreibt genau eine Zeile in runArgs`) steht neben dem Satz „ein Anschluss an einen Ollama-Dienst des Hosts … gibt es nicht“ und der Aufzählung „Cloud-Modelle sind nicht enthalten“; `README.de.md` hat in der Namensliste zudem ein fehlendes Komma vor `ollama`. Der Diff schafft den Widerspruch bewusst; die Auflösung liegt beim Spezifikations-Slice der Welle. Zuständigkeit: Planner, Validator. | `AGENTS.md` §3.10; Carveout 1 | `docs/user/benutzerhandbuch.md` · „gibt es nicht.“; `README.de.md` · „(GitHub- und GitLab-CLI) `ollama` (nur CPU“ | `make docs-check` (grün) | bewusster Zwischenwiderspruch, Kommaversehen |

Nicht als Finding geführt, aber vermerkt (Zuständigkeit Verifier): Die Mutationstabelle zeigt Befehl und
Wiederherstellung einmal in Prosa statt je Zeile (Form des Plans in §3 *Mutationsprotokoll*); die vier gefahrenen
Mutationen belegen den Inhalt.

## Prüfschwerpunkte, Ergebnis

1. **Verhalten gegen Lastenheft und ADR.** Genau eine Zeile `--add-host=host.docker.internal:host-gateway` aus `hostGatewayRunArgs`, im Standardprofil vor dem frühen `return`, im Sandbox-Profil als letztes Element nach `nestedRuntime` und Egress (Golden in `…SandboxRunArgs`, `…TogetherWithOtherFeaturesAndNestedRuntime`). Keine Variable, kein `forwardPorts`, kein Relay, kein Egress-Eintrag (`egressFeatureHosts` unverändert, nur Kommentar), `agentCredentialEnv` und `doctor` unverändert; H4, H5 und H6 prüfen die Abwesenheit in allen erzeugten Dateien. Inaktiv byte-identisch: Template `{{- else if .RunArgs}}`, bestehende Goldens unverändert grün, M07 färbt vier Tests rot. Ohne Befund.
2. **Meldung.** Eine `WarningEntry` je Lauf (auch beim NoOp). Sandbox: Code der Sandbox-Anforderung ([`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)), `SecurityRelaxation: true`, unter `--quiet` sichtbar (`TestRun_Generate_OllamaHostNotices`, am Nutzerpfad). Standard: Code der Feature-Anforderung ([`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)), ungeflaggt, unter `--quiet` entfällt sie. Egress-Zusatz nur bei Sandbox und `egressEnabled`, derselbe Eintrag, Exit 0, keine Adresse in der Menge. Text ohne Kennung, ohne `ollama.com`, ohne „cloud“ und „model“, ohne Ziffern, ohne unbedingte Erreichbarkeit; auf Linux, lokale Engine, Loopback begrenzt wie in der ADR. Ohne Befund (F-3 als Lesepunkt).
3. **Eintragstyp.** `hasNoSource()` an `projectFeatureSourceVersion` und `featureVersionNeedsSource`; die eine Regel wird von `validateDevcontainerFeatures` (Lesen) und `revalidateFeatureEntry` (`config set`) gefragt; Fehlertext je Typ über `featureWithoutVersionReason`. Keine neue Fehlerklasse: `ErrInvalidFeatureSource` und `ErrConfigValueInvalid` sind vorhanden und auf Exit 10 abgebildet; der Exit-Pin liegt am Nutzerpfad (`TestRun_OllamaHostVersionWithoutSourceExits10`, drei Vorschaumodi). Dual-Classifier-Regel: entfällt, kein neuer Sentinel. `source:` macht ein externes OCI-Feature ohne Zeile und ohne Meldung des Anschlusses (`…SourceOverrideIsExternalFeature`). Ohne Befund.
4. **Tests.** Vier Mutationen selbst rot, Rücknahme sauber; H3-Literalpin vorhanden und wirksam (M07). D1 und D2 nutzen den Weg von `u-boot devcontainer up`, kein Aufruf von `ollama.com`, kein Konto, Laufzeit unter 11 s; Budget `dcTestBudget` 20 Minuten aus dem Bestand; Aufräumen belegt (0 Container). Ohne Befund (F-4 als Lesepunkt).
5. **Kommentare, Plan, Doku.** Kein Treffer auf Muster §3.9 (a) in hinzugefügten Zeilen; Kommentare zu `featureCatalogueEntry`, `RunArgs`, `egressFeatureHosts`, `devcontainerTemplateData`, `validateDevcontainerFeatures` sind angepasst und tragen die geltende Zusage; Verweise auf `dockerfileBlock`-Symbole im Baum nachgesucht, keine veraltete Aussage. Nutzer-Doku ohne Kennungen (`AGENTS.md` §5 Regel 6), die CHANGELOG-Aussagen sind durch Tests gedeckt (auch „that a service of the host is reachable is not tested“). F-1, F-6.
6. **`ollamaNotices`-Filter und Lint.** `make lint` (Teil von `make gates`) meldet `0 issues`; der Filter überdeckt den neuen Namen (F-2). Doctor und Drift: `TestDoctor_FeaturesDrift_OllamaHost` und `TestDoctor_OllamaHost_NoNewFinding` pinnen Stille für Orphan, Drift, Allowlist und Credentials.

## Negativbefunde

Geprüft, ohne Befund:

- `internal/hexagon/application/devcontainer_sandbox.go` · `devcontainerTemplateData` (Reihenfolge, nil-sicherer Pfad für `init --devcontainer`)
- `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` (gültiges JSONC, Schlüsselreihenfolge, keine Leerzeilen-Drift)
- `internal/hexagon/application/devcontainer_egress.go` (nur Kommentar, auf Gelesenes begrenzt)
- `internal/hexagon/application/generate.go` (Aufruf neben `ollamaWarnings`, ein Eintrag je Lauf)
- Sentinel- und Exit-Code-Abdeckung (kein neuer Sentinel; Exit-10-Pins am Nutzerpfad)
- Architekturgrenzen (kein Import, keine I/O im Hexagon; `internal/adapter/driving/cli/warnings.go` unverändert, Doc-Kommentar weiter wahr)
- Sprachvertrag (Meldung, Fehlertexte, generierte Dateien Englisch)
- File Safety (Managed Block, Nutzerzeile bleibt; `…ReplayAndLaterActivation`)
- Referenzmodell (keine ADR- oder Slice-Verweise in `spec/`, Nutzer-Doku ohne Kennungen)
- `harness/replay.md` (Golden-Zeile vorhanden)
- Gate-Lockerung (`.golangci.yml`, `Makefile`, Schwellen unverändert; kein `//nolint`)
- Docker-only (alle Läufe über Docker)

Nicht geprüft:

- Wirkung der erzeugten Zeile gegen einen Dienst des Hosts, Relay, Egress-Verhalten im Lauf — nicht Gegenstand des Slice (ADR *Nicht gemessen*); Validator-Aufgabe
- Mutationen außerhalb der vier selbst gefahrenen — Aufwand und Speicher; Angaben des Implementers plausibilisiert
- volles `make test-docker`, `make verify-depguard`, `make ci` — siehe Sensoren
- Plan-Abgleich nach der letzten Fix-Runde, Closure-Notiz, Register, Risiko-Ausgänge — Closure-Pflichten, Verifier und Implementer

## Summary

0 HIGH · 0 MEDIUM · 2 LOW (F-1, F-2) · 4 INFO (F-3 bis F-6). Handoff an den Verifier möglich; F-1 gehört in den Plan-Abgleich der Closure.
