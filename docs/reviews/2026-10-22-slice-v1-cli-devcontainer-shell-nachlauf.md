# Review (Nachlauf): slice-v1-cli-devcontainer-shell, Fix-Runde

* Rolle: Reviewer (frischer Kontext, nicht der Autor)
* Gegenstand: Diff `3b2b5984..2ac89eb4` (drei Commits: Fix-Runde, Verlinkung des Reports, Plan-Nachtrag §3), Plan
  [`slice-v1-cli-devcontainer-shell`](../plan/planning/in-progress/slice-v1-cli-devcontainer-shell.md),
  Vorbericht [`2026-10-22-slice-v1-cli-devcontainer-shell`](2026-10-22-slice-v1-cli-devcontainer-shell.md)
* Geprüft gegen: die fünf Findings des Vorberichts, [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md),
  [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten),
  [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), AGENTS.md §3.7, §3.9, §3.10, §5 Regel 6.
  Nicht geprüft gegen die DoD (Verifier).
* Sensoren: `make gates` Exit 0 auf `2ac89eb4` (Stempel `gates-passed.head` = HEAD, Arbeitsbaum sauber).
  `make test-docker` und `make verify-depguard` nicht gefahren: Der Diff berührt keinen Docker-Test
  und keine Importregel (nur Kommentare, ein Hilfetext, ein Testfall, Handbuch).

## Findings

### HIGH

Keine.

### MEDIUM

Keine.

### LOW

**L-1 · Plan-Nachtrag verweist auf einen Bericht ohne Artefakt**
* quelle: AGENTS.md §3.10 (Abgleich der Plan-Tabelle nach der letzten Fix-Runde), `harness/roles.md` (Übergaben als benennbares Artefakt)
* pfad: [`docs/plan/planning/in-progress/slice-v1-cli-devcontainer-shell.md`](../plan/planning/in-progress/slice-v1-cli-devcontainer-shell.md) · „Sensoren und Abgleich der Plan-Tabelle: siehe den Bericht an den Hauptagenten“
* klasse: uebergabe-ohne-artefakt
* befund: Die Sensorläufe und der Tabellenabgleich der Fix-Runde stehen nicht im Plan, sondern verweisen auf eine
  Nachricht, die nach dem Lauf nicht auflösbar ist. Die Dateien der Runde (`devcontainer_shell_test.go` in
  `internal/adapter/driving/cli` und `internal/hexagon/application`) tragen keine eigene Tabellenzeile; sie sind durch
  „Ist gegen Plan“ gedeckt, der Nachtrag nennt sie nicht.
* risiko: Der Verifier findet die Evidenz der Runde nicht an einem Ort; §3.10 verlangt den Abgleich nach der letzten Fix-Runde im Text.
* verifizierbar: Review-only

### INFO

**I-1 · Spezifikation führt `devcontainer shell` in den `--quiet`-Aufzählungen nicht**
`spec/spezifikation.md` (Abschnitt zu `--quiet`) nennt `devcontainer up`, aber nicht `shell`, in der Liste der ändernden
Subkommandos, in Ausnahme (1) und (2) sowie bei den Kennzeichnungen. Das Handbuch und das Lastenheft nennen `shell`.
Der Plan führt `spec/` als „unverändert“ und verweist auf den Folge-Slice
[`slice-v2-spec-devcontainer-shell-und-zweites-up`](../plan/planning/open/slice-v2-spec-devcontainer-shell-und-zweites-up.md);
die Behauptung „keine weitere Aufzählung“ in der M-1-Zeile des Nachtrags gilt damit nur für `docs/user/`, READMEs, Hilfe und
CHANGELOG. Zuständigkeit: Folge-Slice; Eingabe für den Verifier.

**I-2 · Hilfetext „this needs the u-boot.yaml“ ist gröber als die Exit-Zeile**
Die Zeile `11` der Exit-Liste sagt „if u-boot.yaml configures the capability“, der Absatz „needs the u-boot.yaml of the
project“. Ein `u-boot.yaml` ohne Fähigkeit liefert ebenfalls keine Statusdatei. Das Handbuch ist hier genau („über die
Fähigkeiten, die `u-boot.yaml` anfordert“). Kein Widerspruch im Regelfall, kein Finding ohne Konventions-Anker.

**I-3 · Verifier-Hinweis**: Das Handbuch-Wortlaut der `--quiet`-Zelle hat keinen Test (Review-only, der Plan sagt es selbst).

## Prüfung der fünf Findings (Frage 1)

* **M-1 behoben.** Die Zelle `--quiet` nennt `devcontainer shell` in Ausnahme (1) (das „nur“ entfällt) und in Ausnahme (2);
  der Absatz „Fallbacks nach dem Start“ nennt `shell` („vor dem Beginn der Sitzung“). Eigene Suche in `docs/user/`,
  `README.md`, `README.de.md`, `CHANGELOG.md`, `root.go` (`--quiet`-Usage allgemein), `devcontainer.go`/`up.go`-Hilfe: keine weitere
  Aufzählung, die `shell` auslässt. Zeilen 62 und 2549 des Handbuchs („Nur `up`, … `devcontainer shell`“) nennen `shell`.
  Aufzählungen in `spec/` fehlen `shell`: siehe I-1 (geplanter Folge-Slice).
* **L-1 (Vorbericht) behoben.** Der Langtext nennt die Bedingung; `TestExecute_DevcontainerShell_Help` pinnt den vollen Satz; Handbuch
  und Exit-Zeile stimmen überein (siehe I-2 für die Feinheit).
* **L-2 behoben.** Das Handbuch nennt keinen Ort mehr („als Einträge … weitergibt“); der neue Fall in `TestSessionEnvNames`
  (`containerEnv`) pinnt, dass die Form überall zählt.
* **L-3 behoben.** Der Port-Kommentar sagt, `Stdin` sei die Terminal-Datei, `Stdout` werde unverändert durchgereicht und nicht
  geprüft, der gemischte Fall sei nicht getestet. Das stimmt mit `isTerminal` (nur Standardeingabe) überein.
* **L-4 behoben.** Der Testkommentar beschreibt die Gruppe samt Umhüllung, ohne Zahl.

## Neue Aussagen und Kommentare (Frage 2)

* §3.10: Neu ist nur die Bedingung im Hilfetext (durch den Help-Test gepinnt) und die Entfernung des Ortes im Handbuch (durch den
  neuen Testfall gedeckt). Die Ergänzung „`shell`“ in Ausnahme (1)/(2) und im Fallbacks-Absatz stützt sich auf bestehende Tests
  (`TestExecute_DevcontainerShell_Quiet`, `…FallbackWarningsComeBeforeTheSession`), laut Plan Review-only im Wortlaut.
* §5 Regel 6: keine neue Kennung in `docs/user/`; der Diff dort enthält weder `LH-`, `ADR-`, `SPEC-`, `ARC-` noch `slice-`.
* §3.7/§3.9 (a bis c): keine Chronik-Wörter in den hinzugefügten Zeilen unter `internal/`; der Port-Kommentar nennt die geltende
  Grenze; kein Helfer zwischen Doc-Kommentar und Funktion; keine umbenannten Symbole.

## Hilfe-Stand nach Zurücksetzen der Mutation (Frage 3)

Der Commit-Stand von `newDevcontainerShellCommand` (Langtext, Exit-Zeile 11, Beispiele) und der Erwartungstext in
`TestExecute_DevcontainerShell_Help` stimmen wörtlich überein (nach `strings.Fields`-Normalisierung); der Arbeitsbaum ist sauber,
`make gates` (inklusive `make test`) auf diesem Stand grün. Die Hilfe ist im Commit vollständig.

## Report-Verlinkung (Frage 4)

Der Diff von [`docs/reviews/2026-10-22-slice-v1-cli-devcontainer-shell.md`](2026-10-22-slice-v1-cli-devcontainer-shell.md) hat fünf Hunks; in allen ersetzen die `+`-Zeilen nur
nackte `ADR-`/`LH-`-Kennungen durch Links und brechen Zeilen um. Die Wortlaute von Befunden, Klassen und Empfehlung sind
unverändert.

## Plan-Nachtrag (Frage 5)

Der Ausgang je Finding stimmt mit dem Diff überein (fünf Zeilen, M-1 und L-1 bis L-4, je mit passender Datei). Der Abgleich
gegen `git diff --stat 3b2b5984..HEAD`: `devcontainer.go` (CLI), `devcontainer.go` (Port), Handbuch haben Tabellenzeilen;
die beiden `devcontainer_shell_test.go` sind unter „Ist gegen Plan“ geführt; Report und Plan selbst ebenfalls. Offen: L-1 oben.

## Negativbefund (geprüft, ohne Befund)

* `internal/adapter/driving/cli/devcontainer.go`, `devcontainer_shell_test.go`: Hilfetext Englisch, ohne Kennung (Test).
* `internal/hexagon/port/driving/devcontainer.go`: nur Kommentar; keine Schichtänderung, Dual-Classifier unberührt.
* `internal/hexagon/application/devcontainer_shell_test.go`: ein zusätzlicher Fall, kein Produktionscode.
* `docs/user/benutzerhandbuch.md`: `--quiet`-Zelle, Umgebungs-Absatz, Fallbacks-Absatz; READMEs, CHANGELOG ohne Drift.
* Gates: keine Lockerung, kein `//nolint`.

## Nicht geprüft

* `make test-docker`, `make verify-depguard` (Begründung oben); die Mutationen der Fix-Runde nicht reproduziert.
* DoD, Closure, Evidence: Verifier.

## Empfehlung

**GO.** Kein HIGH, kein MEDIUM; die fünf Findings des Vorberichts sind behoben. L-1 (Plan-Nachtrag ohne auflösbare
Evidenz) vor der Closure im Plan nachtragen; I-1 geht an den Spezifikations-Slice.
