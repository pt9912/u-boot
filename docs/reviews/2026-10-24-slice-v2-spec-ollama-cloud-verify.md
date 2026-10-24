# Verifier-Report: slice-v2-spec-ollama-cloud — 2026-10-24

**Rolle:** Verifier (Modul 8). Frage: Bauen wir es richtig — gegen Plan, DoD, Lastenheft, ADR und Sensor-Belege.

**Gegenstand:** `git diff c3e8044e..HEAD` (Commits `c0b5a757`, `19508c19`, `4293a4cc`, `f8a6dd33`, Review-Report `3207da2e`, Fix-Runde `e942484d`), Slice
[`slice-v2-spec-ollama-cloud`](../plan/planning/done/slice-v2-spec-ollama-cloud.md)
(Welle [`welle-ollama-cloud-modelle`](../plan/planning/welle-ollama-cloud-modelle.md)).

**Maßstab:** [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Absätze `ollama`, Cloud-Modelle, `ollama-host`, Fassung 0.23.2),
[`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Ausnahme 1),
[`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) (Accepted: Z1 bis Z13, *Nicht gemessen, nicht zugesagt*),
Review-Report [`2026-10-24-slice-v2-spec-ollama-cloud.md`](2026-10-24-slice-v2-spec-ollama-cloud.md),
die Messprotokolle [`slice-v1-ollama-cloud-messung-mit-konto`](../messungen/slice-v1-ollama-cloud-messung-mit-konto/README.md) und
[`slice-v1-ollama-cloud-messung`](../messungen/slice-v1-ollama-cloud-messung/README.md),
`harness/verification.md`. Nicht Gegenstand: Closure-Notiz, Auflösung von Carveout 1 (Planner danach).

## Urteil

**Konform.** DoD-Liefer-Punkte 1 bis 3 und der Gate-Punkt sind durch den Diff, existierende pinnende Tests und selbst gefahrene Sensoren belegt. Kein Befund, der
die Closure blockiert. Die Review-Befunde F-1 bis F-8 sind im Text behoben (Tabelle unten, am Wortlaut gelesen). Vier Hinweise (V-1 bis V-4) betreffen
Zustandsmarker und Zähler im Slice-§7, die bei der Closure zu setzen sind; kein Hinweis betrifft eine Zusage. Nicht geprüft und beim Validator der Welle:
die Abhilfe `ollama.com` am Auslieferungspfad mit `u-boot devcontainer up --recreate` und die leere Variable `ANTHROPIC_API_KEY` im erzeugten Sandbox-Profil.

## Sensoren (selbst gefahren, Belege gelesen)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` (sauberer Arbeitsbaum) | Exit 0 | `d-check: 873 Datei(en) geprüft, 0 Befund(e)`; lint `0 issues.`; `coverage-gate: OK — coverage 92.10% meets threshold 90%`; `ablage-pruefung: OK`, `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`; Arbeitsbaum danach sauber |
| `make doc-trace` | Exit 0 | `84 Anforderung(en), 0 Waise(n)` |
| `make doc-commits RANGE="origin/main..HEAD"` | Exit 0 | `9 Commit(s) — OK`, `0 Befund(e)` |
| `make doc-immutable RANGE="d9dd171e..HEAD"` | Exit 0 | `46 Commit(s) — OK`, `0 Befund(e)` |
| `git diff d9dd171e..HEAD` über `docs/plan/adr/0036-…md` und `0024-…md` | leer | keine Zeile Diff |
| `git diff --stat c3e8044e..HEAD` über `spec/architecture.md`, `spec/lastenheft.md`, `docs/plan/adr/`, `internal/`, `cmd/` | leer | Slice berührt keinen Code, kein Template, keine ADR, nicht das Lastenheft, nicht die Sicht |
| `make test-docker` | nicht gefahren | Docs-only-Diff (`internal/`, `cmd/`, Dockerfile, Template unberührt); die E2E-Tests der Zeile B15 sind Bestand des Code-Slice (dort gefahren); ihre Rümpfe habe ich gelesen, nicht erneut ausgeführt |
| `make verify-depguard` | nicht gefahren | keine Importregel berührt |

Mutationen habe ich nicht gefahren: Der Slice ändert keinen Code; die Rot-Belege der Zusagen stehen im Mutationsprotokoll des Code-Slice (`done/`, §7) und sind dort
vom Verifier teilweise reproduziert.

## DoD, Punkt für Punkt

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| 1 Spezifikation und Sicht nach der Fundstellen-Tabelle | erfüllt | [`SPEC-052`](../../spec/spezifikation.md#spec-052--feature-ollama-host-im-devcontainer-eintragstyp-container-option-meldung-und-grenzen) neu (Aktiv, Option, Wirkung in beiden Profilen, `version`/`source`, Meldung als Tabelle mit 4 Zellen je Zeile, Hosts und Zugangsdaten, Auflösung des Namens, Nicht zugesagt); Nachträge in [`SPEC-047`](../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen) (Absatz *Cloud-Modelle*, Halbsatz aus *Nicht zugesagt* entfernt), [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [`SPEC-023`](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (dritter Eintragstyp), [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Ausnahme 1, *gekennzeichnet*, *ungekennzeichnet*), [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) (Sandbox-Regel, Warnungstabelle 4 Zellen, *Zugangsdaten*), [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) (Zeile und Schritt 6), Änderungshistorie mit einer Zeile; `architecture.md` unverändert (Diffstat leer; Suche `ollama`, `runArgs`, `host.docker.internal`, `Anschluss` 0 Treffer); `+`-Zeilen in `spec/` nach `ADR-`, `Slice`, `Welle`, `Roadmap`, `Carveout`, `docs/plan`: 0 Treffer |
| 2 Handbuch, READMEs, `CHANGELOG.md` | erfüllt | Punkt *Cloud-Modelle von Ollama* ersetzt „nicht enthalten“; neuer Abschnitt *Anschluss an einen Ollama-Dienst des Hosts (optional)*; Nachträge in `--quiet`-Zeile, Katalogzeile (Link), *Features im Sandbox-Profil*, Profil-Eigenschaftentabelle (3 Zellen), Egress-Tabelle (2 Zellen), *Claude Code* (Zugangsdaten), *Zugangsdaten (Git und Agent)*; READMEs verweisen auf das Handbuch; `+`-Zeilen in `docs/user`, `README.md`, `README.de.md` nach Anforderungs-, Spec-, ADR-, Slice- und Wellen-Kennungen, `roadmap`, `spec/`, `docs/plan`, Carveout, Register-Kennungen und Token-Präfixen: 0 Treffer (`AGENTS.md` §5 Regel 6 gewahrt); Tokenwerte: nur `<Platzhalter>` |
| 3 Beleg-Tabelle | erfüllt | Alle in B1 bis B20 genannten Testnamen existieren (Stichprobe: 20 Namen je mit Datei und Zeile gefunden, die 22 Funktionen aus `generate_ollama_host_test.go` gelistet); Rümpfe gelesen: `DefaultProfileRunArgs`, `DefaultProfileKeyOrder`, `SandboxRunArgs`, `InactiveIsByteIdentical`, `WritesNoVariableAndNoRelay`, `AddsNoEgressHost`, `AddsNoCredentialVariable`, `SandboxNoticeIsRelaxation`, `DefaultNoticeIsOrdinary`, `EgressRestrictionIsAWarning` (inkl. `egress-init.sh` gleich), `RendersNoKey`, `VersionWithoutSourceRejected`, `ConfigSet…Rejected`, `SourceOverrideIsExternalFeature`, `NotActive`, `ReplayAndLaterActivation`, `TogetherWithOtherFeaturesAndNestedRuntime`, `NoticeText`, `TestRun_Generate_OllamaHostNotices` (je Profil mit und ohne `--quiet`), `TestFeatureCatalogue_EntriesAreValid`/`SecurityProperties`, `TestFixedConfigPaths_SetIsPinned` und `…UsesFeaturePlaceholderKinds` (inkl. `.relay` ungültig), `TestGenerateDevcontainer_Ollama_Sandbox`, E2E D1/D2 (Gegenprobe ohne Feature, `hosts=`-Zeile). Jede gelesene Zusage der Spezifikation hat einen tragenden Rumpf; die Streichungen S1 bis S16 stehen in der Tabelle |
| 4 Gates | erfüllt | selbst gefahren (Tabelle oben) |
| 5 Review durchgeführt | erfüllt | Report liegt vor; Fix-Runde nachgelesen (Tabelle unten) |
| 6 bis 9 (Doku-Update, Closure-Notiz, Register, Risiken, Paarungen) | offen, nicht Gegenstand | gehören der Closure; Carveout 1 steht bewusst bis dahin; `roadmap.md` trägt die Arbeitszeile |

## Fix-Runde: F-1 bis F-12

| Befund | Stand im Text | Urteil |
|---|---|---|
| F-1 (MEDIUM) `--quiet`-Aufzählung | Handbuch-Zeile `--quiet`: „der Netzanschluss an den Host durch `ollama-host` im Sandbox-Profil“; Profil-Tabelle: Zeile „Netzanschluss an den Host“ (3 Zellen) | behoben |
| F-2 `Nicht zugesagt` unbedingt | [`SPEC-052`](../../spec/spezifikation.md#spec-052--feature-ollama-host-im-devcontainer-eintragstyp-container-option-meldung-und-grenzen): Erreichbarkeit „außer in dem Aufbau, den das Lastenheft benennt“ mit dessen Bedingungen; Test prüft nur die Auflösung | behoben |
| F-3 Allsatz „keine Umgebungsvariable“ | [`SPEC-047`](../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen): „Er setzt dafür keine Umgebungsvariable (weder …)“ mit den fünf gepinnten Namen | behoben |
| F-4 Qualifier Neuanlage | Handbuch: „(hinter der Egress-Restriktion, der Eintrag `ollama.com` von Hand gesetzt)“, „zugesagt ist nichts“ | behoben |
| F-5 §7 zu `[1m]` | Rohdatei `rohdaten-06-discover-1-114526.txt` gelesen: viermal `Error: invalid model name` (Fälle `nocloud` und `apikey`, `run` und `launch`); Handbuch nennt beide Pfade mit Einschränkung, stimmt | behoben |
| F-6 CHANGELOG „runs no sign-in“ | „you sign in yourself; `u-boot` sets no variable for it“ | behoben |
| F-7 Bestandssätze zu `ollama` | [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) und *Kein Start, nur Loopback*: `runArgs` nur „im Sandbox-Profil ohne Egress-Restriktion und ohne `nestedRuntime`“; Pin `TestGenerateDevcontainer_Ollama_Sandbox` (Rumpf gelesen) | behoben |
| F-8 Zusage-Ton | Einleitung des Host-Abschnitts nennt zuerst die Bedingungen; „In der Messung war die Zeile von Hand gesetzt“; Standardprofil-Punkt „wirkt auch im Standardprofil“, kein „bietet“ | behoben |
| F-9 Zustandsmarker | nicht geändert, siehe V-1 | offen bis Closure |
| F-10 Abhilfe `ollama.com` | Satz nennt, dass der Neuaufbau mit der Dev-Container-CLI gemessen ist | ausgewiesen; Validator |
| F-11 Begriffsdoppelung | nicht geändert | INFO, keine Zusage berührt |
| F-12 Commit-Historie | `make doc-commits` 0 Befunde | erledigt |

## Prüfpunkte des Auftrags

1. **Aussage → Beleg.** Spezifikation: jede neue Bedingung hat einen existierenden, tragenden Test (oben). Handbuch: Messaussagen gegen die Rohdaten gelesen — `Error: ollama cloud is disabled: remote model details are unavailable`, `Warning: ignoring --model … because cloud is disabled`, `Error: headless --yes launch for claude requires --model <model>` (`…-06-discover-1-114559.txt`, Exit 1), `Claude Code doesn't recognize the model …` (`…-114507.txt`), `dial tcp: lookup ollama.com on 127.0.0.1:53: server misbehaving` (Protokoll ohne Konto), `You need to be signed in to Ollama to run Cloud models.`, Exit 124 nach 120 s (`…-05-restricted-1-111639.txt`, `…-05-allow-1-111902.txt`), `40 von 40`, Ubuntu 24.04.5, Docker 29.9.2, Server `0.40.10`, Claude Code `2.1.310`, Node `v24.22.1`: alle stimmen. Qualifier gegen die Pfad-Spalte: „den Sie selbst starten“, „von Hand gesetzt“, „Handvariante“, „nicht zugesagt“ stehen im Satz. Der Browser-Schritt der Anmeldung ist im Protokoll belegt („Bestätigung im Browser des Projektinhabers“); der Start des Servers (`ollama serve > /tmp/ollama.log 2>&1 &`) entspricht `inner.sh`. Nichts steht als gemessen, was nur abgeleitet ist („abgeleitet“, „Folgerung“, „nicht geprüft“ stehen im Satz).
2. **`AGENTS.md` §5 Regel 6 und Sicht.** 0 Treffer (oben); `architecture.md` unverändert, ADR- und meilensteinfrei; Spezifikation ohne Verweis nach unten. Der Satz „… no longer say …“ steht nur im CHANGELOG, die Änderungshistorie der Spezifikation ist der zulässige Ort für Chronik.
3. **Widersprüche.** Suche nach „nicht enthalten“, „not included“, „gibt es nicht“, „ist nicht aufgenommen“ in `docs/user`, `README.md`, `README.de.md`, `spec/spezifikation.md`, `CHANGELOG.md` [Unreleased]: verbleibende Treffer sind die Release-Zeilen `v0.12.0` der READMEs (historisch, bleiben) und Sätze über andere Gegenstände (APT-Repository, Wildcard-Einträge, `version` ohne `source`). Kein Widerspruch Lastenheft ↔ Spezifikation ↔ Handbuch ↔ READMEs ↔ CHANGELOG gefunden; „runArgs stammt allein“ und „ein Anschluss … unter Nicht zugesagt“ haben 0 Treffer. Das Handbuch schreibt weniger als das Lastenheft (z. B. „ohne Relay oder geänderte Bindung durch Sie nicht erreichbar“ als notwendige Bedingung, danach „für andere Bindungen keine Zusage“).
4. **Z4, L-1, `ANTHROPIC_API_KEY`, Relay.** Z4 steht in der Spezifikation nicht und im Handbuch nur als „zugesagt ist nichts“ mit dem Qualifier des Handaufbaus. Der Schlüssel-Satz in [`SPEC-047`](../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen) (*Volume*) ist unverändert; das Handbuch dehnt ihn nicht auf einen Container ohne Egress-Skript aus und sagt, dass er für den Container, den `u-boot devcontainer up` ohne Egress-Restriktion startet, nicht gemessen ist. Die Egress-Sätze gelten für `ollama run` im Container der Dev-Container-CLI. Die leere Variable `ANTHROPIC_API_KEY` im Container steht als „nicht gemessen“ (Cloud-Abschnitt und *Claude Code*). Der Relay erscheint nur als Voraussetzung und Wirkung (jeder Container im Netz `bridge`, „nicht geprüft“, „Folgerung, nicht gemessen“); „Ein Verfahren für den Relay steht hier nicht“; die Bedingungen S4 und S13 sind nicht als Eigenschaft des Codes übernommen.
5. **Plan-Abgleich (§3.10, Stand nach der Fix-Runde).** `git diff --stat c3e8044e..HEAD`: `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `roadmap.md`, `carveouts.md`, die Slice-Datei, der Review-Report — jede Datei steht in der Plan-Tabelle; Zeilen „update“ ohne Diff sind nur die der Closure (Welle §4, `open/README.md`, Register) oder bereits in `c7da473b` erledigt. Carveout 1 steht unverändert bis auf den Ort des Plans (`in-progress/`), Carveout 2 unverändert.
6. **Dem Validator bleibt:** (a) die Abhilfe `ollama.com` in `egress.allow` mit `u-boot devcontainer up --recreate` am Auslieferungspfad, ohne Konto (Handbuch und Review F-10 weisen die Handvariante aus); (b) die Umgebung des Containers mit `claude-code` im erzeugten Sandbox-Profil bei auf dem Host nicht gesetztem `ANTHROPIC_API_KEY` (leer gesetzt oder ungesetzt); (c) `getent hosts host.docker.internal` und Abwesenheit in der Egress-Liste am ausgelieferten Binary (Entscheidung 11 (iii) der ADR).

## Hinweise für die Closure (keine Blocker)

| Nr. | Hinweis |
|---|---|
| V-1 | §7 des Slice trägt Zustandsmarker mit Ablaufpunkt (§3.10): „Review, Verification und Closure stehen aus“, „Abgleich nach der Fix-Runde folgt“, Gate-Läufe „Stand vor der Niederschrift … der Lauf auf dem Endstand steht im Übergabebericht“, „ein weiterer Commit verlinkt die Kennungen“, die Commit-Zeile nennt `e942484d` und `3207da2e` nicht. Mit der Closure auf den Ist-Stand setzen; den Plan-Abgleich (Punkt 5 oben) dort eintragen. |
| V-2 | Die Zähler der Vorab-Suche in §7 stehen auf dem Stand vor der Fix-Runde (zeilenweise gezählt jetzt Handbuch `host.docker.internal` 8 statt 7, `Sicherheitslockerung` 9 statt 8). Bei der Closure neu zählen oder als Stand vor der Fix-Runde kennzeichnen. |
| V-3 | Carveout 1 trägt den Satz, Spezifikation und Handbuch führten Z4 „als Handvariante ohne Zusage“; die Spezifikation führt Z4 nicht. Der Eintrag entfällt bei der Closure (Reihenfolge in *Ansatz*); die Abweichung ADR gegen Lastenheft gehört dann in Welle §7 *Offene Entscheidungen* (Planner). |
| V-4 | Der `Stand` des Handbuch-Kopfs (`2026-10-24`) und die `Software-Version` (`v0.15.1`) bleiben für den Cut; der CHANGELOG-Eintrag steht unter `[Unreleased]` / *Changed*. |

## Offene Beobachtungen (kein Eintrag im Register durch den Verifier)

- `BEO-DOK/spec-bedingung-ohne-code-entsprechung`: die Bestandssätze zu `sshd`, `claude-code` und `gh`/`glab` sind im Diff als ungepinnt gestrichen beziehungsweise begrenzt (S9, S14, S15); der Kandidat steht im Slice-§7.
- Vorfall der ersetzten Commit-Message (Review F-12) ist sauber: `make doc-commits` 0 Befunde.

## Nicht ausgeführt

`make test-docker`, `make verify-depguard`, `make ci` (`govulncheck`, `image-scan`), Mutationen — Docs-only-Diff, kein Docker-, Importregel-, Release- oder Security-Pfad berührt; Replay entfällt (kein Generator geändert, Golden-Zeile in `harness/replay.md` besteht seit dem Code-Slice).
