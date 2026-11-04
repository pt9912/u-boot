# Review-Report: slice-v1-release-cut-v0.23.0 — 2026-11-04

**Review-Art:** Release-Cut-Review (Text und Versionsstellen) gegen Plan, Spec, Hard Rules und `releasing.md` §3/§4 (Modul 10). Keine DoD-Prüfung — das ist die Frage des Verifiers.

**Gegenstand:** Cut-Commit `84a1ff3e` (8 Dateien, 23 Einfügungen, 18 Löschungen) und `git diff v0.22.0..v0.23.0 --stat` (82 Dateien; Produkt-Diff `-- cmd internal ':!*_test.go' go.mod Dockerfile repo.mk`: 15 Dateien, davon 12 unter `internal/`, `go.mod` unberührt). Slice [`slice-v1-release-cut-v0.23.0`](../plan/planning/done/slice-v1-release-cut-v0.23.0.md); der Tag `v0.23.0` steht auf `84a1ff3e` (gelesen: `git rev-parse v0.23.0^{commit}`).

**Skill:** `.harness/skills/reviewer.md`
**Modell:** claude (Agentenlauf) · **Datum:** 2026-11-04

**Eingangs-Kontext:**

- Slice-Plan (§1 Ziel, §2 DoD, §3 Plan-Tabelle, §6 Risiken 1, 6, 7); der Slice ist bewusst nicht geschlossen (§7 leer, DoD-Haken offen, Roadmap-Zeile „Push und Tag stehen aus“) — kein Finding.
- [`docs/maintainer/releasing.md`](../maintainer/releasing.md) §3 (Vorbereitung) und §4 (Vorab-Prüfung).
- [`LH-FA-DEV-012`](../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`ADR-0044`](../plan/adr/0044-devcontainer-umgebungsdatei-env-file.md), [`ADR-0045`](../plan/adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md); Lastenheft-Versionstabelle Zeile 0.35.0.
- `AGENTS.md` §3.7, §3.9, §3.10, §5 Regel 6; [`harness/review.md`](../../harness/review.md).
- Vorgänger-Report derselben Klasse: [`verification-2026-11-04-slice-v1-release-cut-v0.22.0`](verification-2026-11-04-slice-v1-release-cut-v0.22.0.md) (Lead-Aussage vor dem Cut-Commit gemessen).

**Sensors:**

- `make docs-check`: nach dem Schreiben dieses Reports gefahren (Ergebnis im Handback).
- `make gates`, `make test-docker`, `make verify-depguard`: nicht gefahren — Review-Rolle liest, wiederholt keine Gate-Läufe; keine Importregel berührt (Diff im Produktcode nur `internal/hexagon/application`, Templates, ein Test-Helfer-Registry-Eintrag).
- CI-Lauf `41837587498`: **nicht lesbar** (`gh` hat in dieser Sitzung keine Anmeldung); die Aussage „`--- PASS`“ stammt aus dem Auftrag und ist nicht von mir geprüft (siehe F-2).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der Lead sagt „In a fresh clone the file is missing until `u-boot generate devcontainer` runs once, and the start fails at the engine until then.“ ohne Bedingung. Das gilt nur, wenn die Datei nicht im Repository liegt, also in einem Projekt, dessen `.gitignore` aus `u-boot init` stammt; der *Added*-Eintrag führt die Bedingung nur als eigene Grenze („Not promised: … the ignoring of the file in a project whose `.gitignore` does not come from `u-boot init`“), nicht an dieser Aussage. „Fails at the engine“ ist mit `docker run` (Exit-Code 125) gemessen, der Lead nennt das Messmittel nicht (der Eintrag nennt es). Der Tag ist veröffentlicht; der Block ist nicht mehr in place zu ändern. | `AGENTS.md` §3.10 (Aussage über Verhalten nennt nur, was gelesen oder gemessen wurde) | `CHANGELOG.md` · „In a fresh clone the file is missing until `u-boot generate devcontainer` runs once“ (Lead und *Added*) | Review-only | Lead-Aussage ohne Geltungsbedingung |
| F-2 | INFO | Zu Frage (7): Der Stand „nur ein Lauf in der verschachtelten Engine“ steht im *Changed*-Eintrag („as far as the Docker test shows it (one run in the nested engine of the development devcontainer)“) und im Lead als „measured in one container in one run“. Er ist **kein Befund** — siehe Bewertung unten. | `AGENTS.md` §3.10 | `CHANGELOG.md` · „one run in the nested engine of the development devcontainer“ | Review-only (Lauf `41837587498` lesen) | Messstand-Aussage hinter neuerem Beleg |

### Bewertung zu (7)

1. **Die Aussage ist wahr und deckt sich mit dem Plan.** Sie nennt, was beim Schreiben gemessen war (ein Lauf, verschachtelte Engine); der Lauf auf dem Runner entstand am Push-Tip `84a1ff3e`, also *nach* dem Cut-Commit — der Text konnte ihn nicht tragen, ohne ihn vor seiner Messung zu behaupten. Das ist die Umkehrung der Klasse von `v0.22.0` (dort stand eine Aussage, die nach dem Cut-Commit nicht mehr stimmte); hier steht eine, die weiter stimmt, aber hinter dem Beleg zurückbleibt.
2. **„Zu wenig behauptet“ ist kein Fehler im Sinn von §3.10.** Die Regel verbietet die Aussage ohne Lesung, nicht die vorsichtige; ein Handbuch- oder CHANGELOG-Satz, der enger ist als der Beleg, lässt keinen Änderer einer falschen Aussage folgen (kein Failure-Szenario nach dem Skill, daher kein MEDIUM). Er bleibt im Rahmen von *Not promised* („that limits reach nested containers on any other host or engine“).
3. **Die Asymmetrie ist sichtbar:** Der *Added*-Eintrag der Umgebungsdatei nennt beide Läufe (verschachtelte Engine und `integration-docker` auf dem Runner), der *Changed*-Eintrag des cgroup-Umbaus nur einen. Die Läufe messen aber nicht dasselbe: der Runner-Lauf ist ein nicht verschachtelter Host-Daemon mit einem privilegierten Wegwerf-Container, in dem der verschachtelte `dockerd` läuft — eine zweite Umgebung für denselben Test, keine Wiederholung.
4. **Zwei Vorbehalte für Closure und Verifier, nicht für den Cut:** (a) der Test besitzt einen `t.Skipf`-Pfad („`memory` is already in cgroup.subtree_control of the throwaway container … the test would pass without it“); ein `--- PASS` ist nur dann ein Beleg für den Umbau, wenn der Lauf-Log kein `SKIP` für diesen Test trägt — das lese ich nicht (kein Zugriff auf Lauf `41837587498`); (b) ein Runner-`PASS` stützt die Aussage „Neuanlage-Messung ausstehend“ nicht ab: die Neuanlage des Entwicklungs-Devcontainers bleibt ungemessen ([`slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung`](../plan/planning/open/slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung.md)).
5. **Folge:** Der Runner-Lauf gehört in die Closure-Notiz des Slice (Lauf-ID, Job, Ergebnis des Tests) und in einen künftigen CHANGELOG-Block erst, wenn die Aussage dort gebraucht wird; ein Eingriff in den veröffentlichten Block `[0.23.0]` ist nicht Gegenstand dieses Reports.

---

## Zu den Prüflinsen

**(1) Lead gegen die zwei Einträge und den Diff `v0.22.0..v0.23.0`.** Jede Aussage einzeln gelesen:

| Aussage des Leads | Beleg am Tag-Stand | Urteil |
|---|---|---|
| „Twenty-third release.“ | `[0.22.0]` trägt „Twenty-second release.“ | trägt |
| `generate devcontainer` und `init --devcontainer` schreiben `--env-file=${localWorkspaceFolder}/.devcontainer/devcontainer.env` in `runArgs` | `devcontainer_sandbox.go` · `envFileRunArg`, `devcontainerTemplateData` (`RunArgs` beginnt in beiden Profilen damit); `initproject.go` · `devcontainerFileTemplates` | trägt |
| „the output of every project changes by that line“ | beide Zweige von `devcontainerTemplateData` (`!sandbox` und Sandbox) setzen die Zeile | trägt (für jedes Projekt mit erzeugtem Devcontainer) |
| legt `.devcontainer/devcontainer.env` „with comment lines only“ an, „only when it is missing“ | `templates/devcontainer/devcontainer.env.tmpl` (nur `#`-Zeilen, 10 Zeilen); `generate.go` · `planDevcontainerEnvFile` (`exists` → NoOp), `initproject.go` · `actionKeepExisting` | trägt |
| neue `doctor`-Prüfung warnt bei fehlender Datei | `doctor_envfile.go` · `checkDevcontainerEnvFile` (Warnung nur, wenn `devcontainer.json` die Zeile trägt — im Lead ohne diese Bedingung, im *Added*-Eintrag mit) | trägt, Bedingung nur im Eintrag (kein Befund: der Lead ist Kurzform) |
| „fresh clone … start fails at the engine until then“ | Eintrag: Exit-Code 125 mit `docker run` gemessen | trägt eingeschränkt — F-1 |
| `docker-init.sh` „of every sandbox project with `nestedRuntime: docker`“ ändert sich | `generate.go` · `if data.NestedDocker { … docker-init.sh … }`; `docker-init.sh.tmpl` +45 Zeilen | trägt |
| baut die Root-cgroup um, „before the nested `dockerd` starts, best-effort“ | `docker-init.sh.tmpl` · `cgroup_prepare` steht im Zweig vor `setsid nohup dockerd`; jeder Schritt endet mit `cgroup_prepare_failed` und `return 0` | trägt |
| „so that `--memory` and `--cpus` can reach nested containers“ | Zweckangabe; die Wirkung nennt der Lead im Folgesatz ausdrücklich als „measured in one container in one run“, und *Not promised* nimmt andere Hosts und Engines aus | trägt (Zweck, nicht Zusage) |
| „its effect after a recreation of the development devcontainer is not measured“ | Slice §1 und §6 Risiko 6; Folge-Slice [`slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung`](../plan/planning/open/slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung.md) liegt in `open/` | trägt |
| „Lastenheft 0.35.0 (new requirement for the environment file)“ | `spec/lastenheft.md` Versionstabelle Zeile 0.35.0 (2026-11-04, [`LH-FA-DEV-012`](../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers) neu, Zusage-Ausnahme in [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) und [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)); keine Lastenheft-Zeile zu [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) im Cut | trägt |

Gilt jede Aussage am Tag-Stand? Der Cut-Commit ändert in `cmd/` und `internal/` allein `var version`; der Produkt-Diff `v0.22.0..v0.23.0` ist der Diff, gegen den der Lead gelesen wurde, und der Tag steht auf dem Cut-Commit: keine Aussage hängt an einer Messung, die zwischen Cut-Commit und Tag liegt. Der Lead trägt **keine** Zahl (kein Zählwert, kein Coverage-Wert, kein Testzählwert), die nach dem Commit abweichen könnte — das ist der Fehlerpfad von `v0.22.0`, hier vermieden.

**(2) Versionsstellen.** `grep` auf `0.22.0` im Baum außer CHANGELOG-Chronik, Compare-Link, `docs/plan/`, `docs/reviews/`, `docs/messungen/`: Treffer nur im Lastenheft (Versionszeile 0.22.0, anderes Thema), in `spec/spezifikation.md` (Historie, anderes Thema) und in den CHANGELOG-Zeilen `[0.22.0]` / `Lastenheft 0.22.0`. `0.22` ohne `.0` ohne Treffer. `cmd/uboot/main.go` (`var version`), `repo.mk` (`VERSION ?=` samt zwei Kommentarstellen), `Dockerfile` (`ARG UBOOT_VERSION=` samt zwei Kommentarstellen: „in-source fallback“ und „coherent“) alle `0.23.0-dev`; READMEs (Status-Zeile, `VERSION=`), Handbuch (`Software-Version`, `u-boot version 0.23.0`, `download/v0.23.0`, `VERSION=0.23.0`, Gültigkeitsbereich) auf `v0.23.0`. `Stand: 2026-11-04` im Handbuchkopf deckt sich mit dem Cut-Datum. Ohne Befund.

**(3) Plan-Tabelle §3 gegen `git show --stat 84a1ff3e`.** Acht Dateien im Commit: `CHANGELOG.md`, `Dockerfile`, `README.de.md`, `README.md`, `cmd/uboot/main.go`, `docs/plan/planning/in-progress/roadmap.md`, `docs/user/benutzerhandbuch.md`, `repo.mk` — jede hat eine Plan-Zeile; die Zeilen „`open/README.md`“ (Plan-Commit, nicht Cut-Commit), „dieser Slice“ und „kein Diff durch den Cut“ (`spec/`, ADRs, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod`) stimmen: `go.mod` und `Makefile` fehlen im Diff `v0.22.0..v0.23.0`, `.github/` und `scripts/` im Cut-Commit. Die Plan-Aussage „12 Nicht-Test-Dateien unter `internal/`“ stimmt mit dem Diff überein (gezählt: 12). Ohne Befund.

**(4) Nutzer-Doku ohne Kennungen.** `grep -E "LH-|ADR-|SPEC-"` über `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`: kein Treffer. Der Handbuch-Diff `v0.22.0..v0.23.0` führt keine Anforderungs-, ADR-, Spec-, Slice- oder Roadmap-Kennung (die Prüf-ID `devcontainer.envfile` ist ein Befehlsausgabe-Name, keine Anforderungs-Kennung). Die beiden READMEs ändern dieselben zwei Zeilen (Status-Zeile, `VERSION=`) mit identischem Versionswert: deckungsgleich. Ohne Befund.

**(5) Compare-Links und Datum.** `[Unreleased]: …/compare/v0.23.0...HEAD`, `[0.23.0]: …/compare/v0.22.0...v0.23.0`, `[0.22.0]` unverändert darunter; `## [Unreleased]` leer als Anker. Block-Datum `2026-11-04` = Datum des Commits `84a1ff3e` (`2026-11-04T13:20:52Z`) = Tag-Datum (der Tag zeigt auf diesen Commit). Ohne Befund.

**(6) Lastenheft-Aussage und Ausgabe-Änderung.** Der Lead nennt „Lastenheft 0.35.0“, die neue `runArgs`-Zeile („the output of every project changes by that line“) und die geänderte `docker-init.sh` („of every sandbox project with `nestedRuntime: docker`“). Alle drei stehen im Lead; die Ausgabe-Änderung steht zusätzlich in *Added* und *Changed*. Der Lead führt die Vorlage `devcontainer.env.example` und die `.gitignore`-Zeile von `init` nicht — beide stehen im *Added*-Eintrag und im Handbuch; kein Befund (Kurzform). Ohne Befund.

**§3.7 / §3.9 über den Cut-Commit.** (a) Kein „früher“, „no longer“, „former“ in den eingefügten Zeilen; (b) die Kommentare, die Symbole oder Zahlenwerte nennen (`repo.mk`, `Dockerfile`, `main.go` samt `0.23.0-dev`), tragen den neuen Wert an allen Stellen; (c) keine neu eingefügte Funktion. Die Zustandsangabe der Roadmap-Zeile („Push und Tag stehen aus“) ist laut Auftrag bewusst bis zum Tag offen — kein Finding.

## Negativbefunde

- Geprüft, ohne Befund: `cmd/uboot/main.go`, `repo.mk`, `Dockerfile` (Versionswert und Kommentare), `README.md`/`README.de.md`, `docs/user/benutzerhandbuch.md` (Kopf, `u-boot version`, `download/v…`, `VERSION=`, Gültigkeitsbereich), CHANGELOG-Compare-Links und Datum, Plan-Tabelle §3.
- Geprüft, ohne Befund: Produktcode `internal/hexagon/application/devcontainer_sandbox.go`, `doctor.go`, `doctor_envfile.go`, `generate.go`, `initproject.go`, `templates.go`, Templates (`devcontainer.env.tmpl`, `devcontainer.env.example.tmpl`, `docker-init.sh.tmpl`, `gitignore.tmpl`, `.gitignore.tmpl` von `basic`) — nur als Beleg gegen den Lead gelesen, nicht als Code-Review (Gegenstand ist der Cut). Sentinel/Klassifikatoren: kein neuer Driving-Sentinel (`doctor`-Prüfung ist eine Warnung; Code in der `jsontestutil`-Registry nachgezogen).
- Nicht geprüft: Lauf `41837587498` und der Log des cgroup-Tests (kein `gh`-Zugang; Auftrag nennt `--- PASS`, ein `SKIP` wäre nicht von mir ausgeschlossen); `make gates`, `make ci`, `make test-docker`, Tap-Trockenlauf, `tap-check`, `publish`/`tap`-Läufe (Verifier bzw. Maintainer-Aktionen M-T/M2/M4/M6); Mac-Probe M7 (Projektinhaber); die Neuanlage-Messung (Folge-Slice, bewusst nach dem Release).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Lead-Aussage ohne Geltungsbedingung · Messstand-Aussage hinter neuerem Beleg.

## Verdikt

**Merge-blockierend:** nein (der Tag ist veröffentlicht; F-1 ist LOW, F-2 INFO). Der Lead sagt, was die zwei Einträge und der Diff belegen; keine Aussage hängt an einem Stand, der zwischen Cut-Commit und Tag liegt.

**Übergabe:** F-1 an den Implementer (Entscheidung, ob die Bedingung bei der Closure im Slice vermerkt oder in einem künftigen Block geführt wird; der Block `[0.23.0]` selbst ist veröffentlicht). F-2 an den Verifier: Lauf `41837587498` auf `SKIP` im cgroup-Test lesen und das Ergebnis in die Closure-Notiz (Lauf-ID, Job, Ergebnis) nehmen; keine Änderung am CHANGELOG-Block nötig.
