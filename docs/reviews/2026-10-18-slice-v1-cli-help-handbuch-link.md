# Review-Report: slice-v1-cli-help-handbuch-link — 2026-10-18 <!-- d-check:ignore (Slice-Kennung im Titel, Slice liegt bei Closure in done/) -->

**Review-Art:** Code — Diff gegen Plan, Entscheidungen und Hard Rules (Maintainability); DoD-Prüfung ist Sache des Verifiers.

**Gegenstand:** Commit-Range `2ab2ee6..HEAD` (HEAD `3213066`), Slice `slice-v1-cli-help-handbuch-link`. <!-- d-check:ignore (Slice-Kennung, Fundstelle des Gegenstands) -->

**Skill:** `.harness/skills/reviewer.md` @ HEAD `3213066` · <!-- d-check:ignore (Skill-Pfad als Zitat) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-18

**Eingangs-Kontext:**

- Slice-Plan `slice-v1-cli-help-handbuch-link` §1–§6 (in `docs/plan/planning/in-progress/`) <!-- d-check:ignore (Slice-Kennung) -->
- Entscheidung des Projektinhabers 2026-10-18 (zwei Links) und Festlegung des Leads (Tag-Link nur bei `X.Y.Z` oder `X.Y.Z-rc.N`)
- [`LH-FA-CLI-002`](../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe), [`LH-FA-CLI-003`](../../spec/lastenheft.md#lh-fa-cli-003--versionsausgabe), [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache), [`LH-NFA-PERF-001`](../../spec/lastenheft.md#lh-nfa-perf-001--schnelle-cli-antwort)
- `.github/workflows/publish.yml` (Schritt „Validate SemVer tag“), `docs/maintainer/releasing.md` §1/§3/§5
- `AGENTS.md` (§3.2, §3.6, §3.7, §3.9, §3.10, §5)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der Kommentar nennt `releaseVersionPattern` „exactly the versions the publish workflow tags“; der Workflow lässt mehr Vorab-Formen zu (`-rc1`, `-beta.1`), die Regex ist eine echte Teilmenge. Der Kommentar von `rootLongText` begründet den Fortfall des Tag-Links mit „because no tag exists for it“; für `0.12.0-rc1` kann der Tag existieren. | `AGENTS.md` §3.7 (Zusage/Grenze beschreibt, was da ist) | `internal/adapter/driving/cli/root.go` · „matches exactly the versions the publish workflow tags as“ und „because no tag exists for it“ | Review-only (Abgleich mit `publish.yml` Zeile „Validate SemVer tag“) | Kommentar-Zusage weicht von gelesener Quelle ab |
| F-2 | LOW | Plan §2 (DoD, Test-Punkt) verlangt weiter einen dritten Test, der für `0.12.0-rc1` den Tag-Link `v0.12.0-rc1` prüft, und §6 (erstes Risiko) sagt „nur `-dev` ist ausgenommen“ sowie `0.0.0-test` bekäme einen Tag-Link. Seit der Festlegung des Leads (§1 Entwurf 2) und im Code gilt das Gegenteil (`-rc1` und `0.0.0-test`: nur `main`). | `AGENTS.md` §3.10 (Aussagen im Plan gegen Code prüfen) | `docs/plan/planning/in-progress/slice-v1-cli-help-handbuch-link.md` · „(`0.12.0-rc1`) prüft den Tag-Link `v0.12.0-rc1`“ | Review-only | Plan-Aussage nach Entscheidungswechsel nicht nachgezogen | <!-- d-check:ignore (Slice-Kennung im Pfad des geprüften Gegenstands) -->
| F-3 | INFO | `docs/maintainer/releasing.md` §5 räumt Vorab-Tags nach dem Cut auf (Release und Tag löschen). Ein `X.Y.Z-rc.N`-Binary, das danach noch läuft, zeigt einen Tag-Link, dessen Tag nicht mehr existiert; die Aussage „Tag existiert für die Release-Form“ gilt nur, solange der Vorab-Tag lebt. | Maintainability (Plan §6, Risiko „Link zum Handbuch einer Version ohne Tag“) | `docs/maintainer/releasing.md` · „Vorab-Tags aufräumen: Release und Tag löschen“ | Review-only | Randfall bewusst klein (rc-Binaries sind kurzlebig); §6-Ausgang bei Closure benennen |
| F-4 | INFO | `regexp.MustCompile` läuft bei jedem Aufruf von `rootLongText`, das Muster ist eine Konstante: ein Panic ist nur bei einer Änderung der Konstante möglich, der Test `TestExecute_RootHelpManualLinks` kompiliert sie in jedem Fall. `Execute` baut das Wurzelkommando einmal je Prozess; der Aufwand liegt im Mikrosekundenbereich, weit unter der Grenze aus [`LH-NFA-PERF-001`](../../spec/lastenheft.md#lh-nfa-perf-001--schnelle-cli-antwort). | Maintainability | `internal/adapter/driving/cli/root.go` · „regexp.MustCompile(releaseVersionPattern)“ | `make test` | Beobachtung ohne Änderungspflicht |

## Prüfung der Aufträge

**(1) Regex gegen `publish.yml` und Praxis.** Die Regex `^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-rc\.(0|[1-9][0-9]*))?$` ist eine strikte Teilmenge der Tag-Prüfung in `publish.yml`: jede ihrer Versionen besteht dort (`rc.N` mit numerischem Identifier ohne führende Null; `rc.0` in beiden zulässig). Umgekehrt lässt der Workflow `-rc1`, `-beta.1`, `-alpha`, `-rc.0a` zu; für sie zeigt der Help-Text nur `main` — die sichere Richtung (kein toter Tag-Link; Preis: ein Vorab-Binary mit existierendem Tag verlinkt das Handbuch seiner Version nicht, das ist eine Funktions-, keine Korrektheitslücke und im CHANGELOG ehrlich benannt). Die Praxis-Schreibweise `-rc.N` stimmt mit `releasing.md` §1/§5 (`vX.Y.Z-rc.1`, `vX.Y.Z-rc.N`) und den Commit-Meldungen der Cut-Slices (rc.1, rc.2) überein. Beim Handbuch (Zeile „`vX.Y.Z-rc.1` landen nicht im Tap“) gleich. `MustCompile`-Panic: nur bei geänderter Konstante, vom Test gedeckt (F-4). `gochecknoglobals`: nur `const`, keine Variable; kein `nolint`. `make lint` danach grün (Formatkorrektur im letzten Commit wirksam).

**(2) Verhalten (Runtime-Image aus dem Arbeitsbaum, `--build-arg UBOOT_VERSION=0.11.0` und Default `0.11.0-dev`).** `u-boot --help` mit `0.11.0`: README-Zeile unverändert, danach „User manual (German), this version: …/blob/v0.11.0/docs/user/benutzerhandbuch.md“ und „… latest: …/blob/main/…“. Mit `0.11.0-dev`: nur die `latest`-Zeile. `u-boot help` zeigt denselben Wurzeltext; `u-boot init --help` und `completion bash` enthalten keinen `blob`-Link (nur die Wurzel-`Long` ist geändert). `--version` unverändert (`u-boot version 0.11.0` bzw. `0.11.0-dev`). Kein `LH-` im Text (bestehender Test plus neuer Fall); Text Englisch, „(German)“ benannt, Sprachvertrag ([`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)) gewahrt. Längste Zeile des Wurzeltextes 111 Zeichen (die URL-Zeile; die Flag-Beschreibungen sind ähnlich lang): kein Konventions-Anker, kein Finding. Windows/macOS: reine String-Logik, unberührt. `Long` wird bei jedem `Execute` neu gebaut, `a.version` liegt vor. Laufzeit: ein Docker-Lauf misst 1,5 s einschließlich Container-Start und ist als Messung der 200-ms-Grenze ([`LH-NFA-PERF-001`](../../spec/lastenheft.md#lh-nfa-perf-001--schnelle-cli-antwort)) ungeeignet; der Mehraufwand ist ein einmaliger Regex-Abgleich je Prozess (F-4).

**(3) Test und Mutationen.** Die 13 Fälle nutzen hartkodierte erwartete Strings (nicht aus den Konstanten des Codes abgeleitet): nicht tautologisch. Eigene Mutationen (nur `make test`, danach `git checkout -- <Datei>`, Arbeitsbaum sauber):

| Mutation | Ergebnis |
|---|---|
| Regex ohne `^` und `$` | gekillt (`0.11.0-dev`, `0.0.0-test`, `0.12.0-rc1` u. a.) |
| `-rc\.` zu `-rc` | gekillt (`0.12.0-rc.1`, `0.12.0-rc1`) |
| `main`-Link entfernt | gekillt (alle Fälle) |
| `repoURL` geändert (`uboot`) | gekillt, auch der bestehende README-Test |
| `manualPath` geändert | gekillt |
| führende Nullen zugelassen (`([0-9]+)`) | gekillt (`0.12.0-rc.01`, `01.2.3`) |
| Beschriftung `latest:` zu `latest` | überlebt — der Wortlaut der Beschriftung ist laut Plan §6 bewusst nicht gepinnt (kein Finding) |

**(4) Praxis.** `gh api repos/pt9912/u-boot/contents/docs/user/benutzerhandbuch.md?ref=v0.11.0`: Exit 0, Pfad vorhanden; Tag `v0.11.0` und Release `v0.11.0` existieren; die Links für `0.11.0` sind nicht tot. Die Tags `v0.10.0-rc.1` und `v0.11.0-rc.1` existieren nicht mehr (aufgeräumt, siehe F-3). Images `uboot-review-rel` und `uboot-review-dev` entfernt.

**(5) Doku.** Handbuch-FAQ: ohne Kennungen, nennt den Stand („die Fassung Ihrer Version, bei einer selbst gebauten oder Entwicklungsversion nur die aktuelle“) und trägt keine Chronik; gegen den Code wahr für Release-Versionen und `-dev`/Selbstbau (sie lässt die Randfälle `-rc1`/`-beta.N` offen, ohne zu irren). CHANGELOG `[Unreleased]` / Added: ehrlich (Nennung der Release-Form und des main-only-Falls), Anker auf [`LH-FA-CLI-002`](../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe). `releasing.md` unverändert korrekt: §3 Schritt 2 setzt die drei Versionsstrings, die Tag-Version kommt zur Laufzeit aus `publish.yml` (`VERSION="${TAG_REF#v}"`) — auch für `-rc.N` entsteht kein Handarbeitsschritt. Alter Wortlaut: „Quickstart and user docs“ steht nur noch in `root.go` (README-Zeile, unverändert); in Tests wird nur der README-Link geprüft; kein Treffer in Doku.

**(6) Plan-Tabelle gegen `git diff --stat`.** Diff berührt sechs Dateien: `CHANGELOG.md`, `roadmap.md`, der Slice-Plan, `benutzerhandbuch.md`, `cli_test.go`, `root.go`. Die Tabelle §3 führt alle außer dem Slice-Plan selbst (der Plan als Gegenstand, kein Eintrag nötig) und `export_test.go` als „prüfen“ (Aussage gelesen: `WalkRootCommandPathsForTest` ruft `buildRootCommand(a)`, Signatur unverändert, Datei im Diff nicht berührt: stimmt). §3.9 (a): kein „früher/former“-Kommentar im Diff; (b): keine Kommentare auf geänderte Symbole, die alte Zeichenkette „Quickstart and user docs“ bleibt gültig; (c): `rootLongText` steht vor `buildRootCommand` und trennt keinen Doc-Kommentar von seiner Funktion. Zu §3.7 siehe F-1.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapter/driving/cli/root.go` (Logik, Konstanten, `gochecknoglobals`, kein `nolint`) | geprüft, ohne Befund (außer F-1, F-4) |
| `internal/adapter/driving/cli/cli_test.go` (13 Fälle, Mutationsfestigkeit) | geprüft, ohne Befund |
| `u-boot --help`, `help`, `<cmd> --help`, `--version`, `completion` im Runtime-Image | geprüft, ohne Befund |
| `CHANGELOG.md`, `docs/user/benutzerhandbuch.md`, `docs/maintainer/releasing.md` | geprüft, ohne Befund |
| `roadmap.md` (Zustandszelle: Zustand und Beleg, keine Chronik) | geprüft, ohne Befund |
| Referenzmodell, Architekturgrenze, Dual-Classifier, Dateisicherheit | nicht berührt (reine Darstellungs-Änderung im Adapter) |
| `make test-docker`, `make verify-depguard` | nicht gelaufen: kein Docker-Pfad und keine Importregel berührt |
| `make docs-check` | siehe Verifikation des Reports unten |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Kommentar-Zusage weicht von gelesener Quelle ab · Plan-Aussage nach Entscheidungswechsel nicht nachgezogen · Vorab-Tag-Aufräumen macht Tag-Link tot

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM; die beiden LOW sind Wartbarkeit (Kommentar und Plan-Text) ohne Vertragsbruch.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Die Verifikation der DoD bleibt beim Verifier.
