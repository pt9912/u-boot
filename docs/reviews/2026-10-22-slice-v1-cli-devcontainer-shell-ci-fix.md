# Review (Nachlauf): slice-v1-cli-devcontainer-shell, CI-Fix-Runde

* Rolle: Reviewer (frischer Kontext, nicht der Autor)
* Gegenstand: Diff `d537ae24..7d63e04f` (Move, Roadmap, Test-Fix `36b30fcd`, Plan-Nachtrag), Plan
  [`slice-v1-cli-devcontainer-shell`](../plan/planning/in-progress/slice-v1-cli-devcontainer-shell.md) §3 *CI-Fix-Runde* und §7
* Geprüft gegen: [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md),
  [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), AGENTS.md §3.7, §3.9, §3.10.
  Nicht geprüft gegen die DoD (Verifier).
* Sensoren: `make gates` Exit 0 auf `7d63e04f`. `make test-docker` nicht gefahren: Der Befund ist ein Wettlauf zweier
  Pakete auf einer Engine, lokal nicht erzwingbar; die CI-Läufe (`integration-docker` grün auf dem Tip) sind die Evidenz des
  Implementers und wurden von mir nicht eingesehen. `make verify-depguard` nicht gefahren: keine Importregel berührt.

## Findings

### HIGH

Keine.

### MEDIUM

Keine.

### LOW

**L-1 · Kein Tabellenabgleich nach der CI-Fix-Runde**
* quelle: AGENTS.md §3.10 (Abgleich gegen `git diff --stat` nach jeder späteren Änderung)
* pfad: [`docs/plan/planning/in-progress/slice-v1-cli-devcontainer-shell.md`](../plan/planning/in-progress/slice-v1-cli-devcontainer-shell.md) · „CI-Fix-Runde (Implementer, 2026-10-22)“
* klasse: tabellenabgleich-nach-fix-runde-fehlt
* befund: Der Nachtrag nennt Befund, Ursache, Fix, Mutationen und Engine-Matrix, aber keinen Abgleich der Plan-Tabelle gegen den
  Diff der Runde. Der Abgleich weiter oben („66 Namen“) gilt für den Stand vor Move und Fix.
* risiko: Der Diff der Runde berührt `internal/e2e/devcontainer_shell_docker_test.go` (Tabellenzeile vorhanden) und
  Verweis-Nachzüge; ohne Zeile ist nicht belegt, dass die Runde nichts außerhalb der Tabelle berührt.
* verifizierbar: Review-only

### INFO

**I-1 · Ursache: Namensform und Hash belegt, Gleichzeitigkeit geschlossen**
Selbst nachgerechnet: `sha256("TestLauncher_Docker_S6_GuardHandWrittenSocketMountIsDetected")` beginnt mit `b4c391e5`.
`helperPrefix(t)` (`devcontainer_launcher_docker_test.go`) ist `HelperContainerPrefixForTest + "t" + hex(sha256(t.Name())[:4]) + "-"`,
also `u-boot-devcontainer-up-tb4c391e5-`; das Binary vergibt `helperContainerPrefix + randomSuffix()` mit genau acht Hex-Zeichen
(`devcontainer_launcher.go`). Der Name `…-tb4c391e5-eb6f2965` ist damit von der Namensform her nur aus diesem Adapter-Test
erklärbar (das `t` ist kein Hex-Zeichen, kein Binary-Name kann so lauten); belegt. Nur plausibel, nicht beobachtet, ist die
Gleichzeitigkeit der beiden Pakete (Plan: lokal nicht reproduziert; Lauf-Log der CI nicht von mir eingesehen). Eine zweite Ursache
ist für den gemeldeten Fehler nicht nötig: Ein verspätet entfernter Hilfs-Container des Binary trüge die Form `…-up-<8 hex>` und
würde vom Regex weiter gezählt.

**I-2 · Der engere Filter schwächt den Wächter nicht (soweit lesbar)**
Das Produkt kennt genau eine Namensform für Hilfs-Container (`helperContainerPrefix` + `randomSuffix`); im Produktcode
(`internal/`, `cmd/`) gibt es keinen Namen `u-boot-devcontainer-down-` oder ein anderes Präfix für Hilfs-Container. Driftet die
Form des Binary (Länge, Zeichensatz, Präfix), zählt der Regex nichts, und die bestehende Gegenprobe „no helper container was seen
during the start“ wird rot; die Kopplung an `randomSuffix` ist damit abgesichert. Das Fenster zwischen Start und Entfernung bleibt
durch `dcWatchHelpers` (100 ms) und die Prüfung nach Sitzungsbeginn abgedeckt.

**I-3 · Wächterfall und Mutationen**
Der Fall ist logisch tragfähig: `foreign` (`…-tb4c391e5-<8 Zeichen>`) liegt neben `hand` (Form des Binary), erwartet wird genau
`hand`; mit Präfix-only-Filter wären es zwei Treffer, also rot. Die Meldung der Mutation 1 im Plan passt zum neuen Code (`hand`
ist jetzt acht Hex-Zeichen). Mutation 2 (`--rm` entfernt, `"run", "--rm"` in `runArgs` vorhanden) ist plausibel rot. Die Läufe habe
ich nicht wiederholt (Review-only). Der Plan benennt ehrlich, dass der Wettlauf selbst lokal nicht reproduziert wurde.

**I-4 · Weitere globale Suchen im Paket**
Ein Treffer für `--filter name=<Präfix>` gibt es nur in `devcontainer_shell_docker_test.go`. Andere e2e-Tests filtern nach Label
des Projektordners (`dcFolderLbl`), der von Adapter-Tests nicht geteilt wird. `newContainersWithCommand` in
`sandbox_devcontainer_docker_test.go` listet alle Container und sucht einen Pin in Path/Args/Env; der Kommentar benennt die
Grenze („Two runs against the same daemon at the same time are not told apart“). Ein Container eines parallelen Pakets mit
demselben Pin ist denkbar, aber nicht belegt; kein Finding, Hinweis für den Verifier.

**I-5 · Zustandsmarker nach grüner CI**
§3 („die Closure bleibt bis dahin bedingt“) und der Vermerk in §7 („bis zum grünen Lauf … nicht gültig“) stehen auf dem Stand
vor dem grünen Lauf auf `7d63e04f`. Aufgabe der erneuten Closure, nur angemerkt.

**I-6 · Kommentare**
§3.7/§3.9 (a bis c): keine Chronik-Wörter in den neuen Kommentaren (`dcHelperName`, `dcHelperNames`, Wächterfall); der Kommentar
nennt die geltende Abgrenzung. Symbole (`randomSuffix`) existieren; `dcHelperName` steht nach der Konstante mit eigenem
Doc-Kommentar, kein Helfer zwischen Doc-Kommentar und Funktion.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/e2e/devcontainer_shell_docker_test.go` (Regex, Filter, Wächterfall) | geprüft, ohne Befund (I-2, I-3) |
| `internal/adapter/driven/docker/devcontainer_launcher.go` (Namensform, `--rm`, Entfernung) | geprüft, ohne Befund |
| Paket `internal/e2e` und Adapter-Tests, Suche nach Präfix-Filtern | geprüft, ohne Befund (I-4) |
| Produktcode des Diffs | keiner geändert (nur Test) |
| Commits `7c7afceb`, `eeb62cba`, `d0fa0747` (Move, Verweise, Roadmap) | geprüft, ohne Befund; Move rein, Roadmap-Zeile nennt Zustand |
| Wiederholung der Mutationen und der Engine-Matrix | nicht geprüft: Review-only, Läufe des Implementers nicht reproduziert |
| `make test-docker`, `make verify-depguard` | nicht gefahren, Gründe oben |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 6 |

**Finding-Klassen dieses Laufs:** tabellenabgleich-nach-fix-runde-fehlt

## Verdikt

**GO.** Merge-blockierend: nein. Die Ursache ist über Hash und Namensform belegt, der Fix ist auf den Test beschränkt und
schwächt den Wächter nicht erkennbar. L-1 geht an den Implementer (Abgleich nachtragen), I-5 an die erneute Closure.
