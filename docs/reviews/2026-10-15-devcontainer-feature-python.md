# Review-Report: [`slice-v1-devcontainer-feature-python`](../plan/planning/done/slice-v1-devcontainer-feature-python.md) — 2026-10-15

**Review-Art:** Code — Diff gegen Slice-Plan, [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md), Spec ([`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) und Replay Rules (`harness/replay.md` §Replay Rules).

**Gegenstand:** Commit `a0ab388` (Katalog-Feature `python`; Pfade `internal/hexagon/application`, `spec`, `docs/user`, READMEs, CHANGELOG, `harness/replay.md`) und der Slice-Plan.

**Skill:** `.harness/skills/reviewer.md` @ a0ab388 · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-15

**Eingangs-Kontext:**

- [`slice-v1-devcontainer-feature-python`](../plan/planning/done/slice-v1-devcontainer-feature-python.md) (DoD, Abgrenzung, Risiken)
- [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) (Katalog-Regel offizielles Präfix)
- [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)
- [SPEC-019](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [SPEC-023](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)
- `AGENTS.md` (Hard Rules), `harness/review.md`, `harness/replay.md`, Slice-Vorlage der vendorten Baseline
- Vorläufer: Report `2026-10-15-devcontainer-features-homebrew-sdkman`

Zusätzlich ausgeführt: `make docs-check` (205 Dateien, 0 Befunde); `gofmt -d internal/hexagon/application/devcontainer_egress.go` im Image `golang:1.27.2`, Quellen read-only gemountet, ohne Netz (Ausgabe leer: konform). `make gates` und `make test` wurden nicht ausgeführt.

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Aussage „das Feature installiert das Betriebssystem-Python“ steht als Tatsache im Handbuch und im Slice-Risiko (Ausgang „entfallen: gemessen, Python 3.13.5“). Im Repo liegt kein Messbeleg (anders als die Messtabelle in [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) beim Vorläufer-Slice). Der Katalog übergibt dem Feature keine Optionen; welche Python-Quelle das Feature mit seinen Standardoptionen wählt, ist im Repo nirgends festgehalten. Ob die PEP-668-Aussage für jede Option des Features gilt, ist damit nur über den Slice-Text nachvollziehbar. | Slice §6, Doku-Präzision | `docs/user/benutzerhandbuch.md:763`; [slice-v1-devcontainer-feature-python](../plan/planning/done/slice-v1-devcontainer-feature-python.md), Zeile 66 | nein — Review-only (Messung braucht Image-Build, außerhalb der Reviewer-Rolle) | Messaussage ohne Beleg im Repo |
| F-2 | LOW | Slice §1 sagt, der Nutzer überschreibe „die Version“ über `devcontainer.features.python.version`. Dieses Feld ist der OCI-Tag des Features (Render-Schlüssel `source:version`), nicht die Python-Version; `version: "3.12"` rendert `…/python:3.12`. Die Formulierung kann als Python-Versionswahl gelesen werden; Feature-Optionen schließt der Slice selbst aus. Das Handbuch führt dieselbe Mehrdeutigkeit schon beim `java`-Beispiel (`features.java.version 21`), nicht Gegenstand dieses Diffs. | Slice §1, [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | [slice-v1-devcontainer-feature-python](../plan/planning/done/slice-v1-devcontainer-feature-python.md), Zeile 27 | ja — Gegenlesen; `projectFeatureSourceVersion` in `internal/hexagon/application/devcontainer_features.go` | Feature-Tag als Python-Version gelesen |
| F-3 | LOW | Slice §3 führt die Testzeile (`*_test.go`) ohne Akzeptanzkriterien-ID in der Begründung; die Vorlage verlangt die ID dort („nach `<Akzeptanzkriterien-ID>`“). Der Vorläufer-Slice hat dieselbe Form, daher kein Einzelfall. | Slice-Vorlage §3 | [slice-v1-devcontainer-feature-python](../plan/planning/done/slice-v1-devcontainer-feature-python.md), Zeile 48 | ja — Gegenlesen | Testzeile ohne LH-Anker |
| F-4 | INFO | Replay-Abgleich: Der Fall `python` liegt als Tabellenfall in `TestGenerateDevcontainer_FeaturesHomebrewSdkman` und passt auf das Muster `FeaturesHomebrew*` der Golden-Zeile in `harness/replay.md`; Fresh-State (Render `python:1`) und Double-Run-NoOp sind gepinnt. Der Hosts-Fall steht in `TestGenerateDevcontainer_SandboxEgress_FeatureHosts`. Die Gegenrichtung („`pypi.org` fehlt ohne aktives `python`“) ist nur indirekt gepinnt: die `absent`-Listen der Fälle `homebrew` und `sdkman` enthalten `pypi.org` nicht. Die Warnungszählung des Falls `python` filtert auf `meaningful-ooo`; eine versehentliche Drittanbieter-Warnung für `python` würde sie nicht fangen (die Zeile `wantWarns` ist 0 für einen anderen Zähler). Stale-Block und User-Content sind über die allgemeinen Feature-Tests abgedeckt; der Katalogeintrag ändert die Planungsphase nicht. | `harness/replay.md` §Replay Rules | `internal/hexagon/application/generate_features_test.go:641`; `internal/hexagon/application/generate_sandbox_test.go:883-886` | ja — Gegenlesen | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Katalogeintrag | geprüft, ohne Befund: `python` → `ghcr.io/devcontainers/features/python`, Standardversion `1`, nicht per Herausgeber-Namensraum markiert; Render `python:1` im Golden Case belegt |
| Katalog-Invariantentest | geprüft, ohne Befund: `TestFeatureCatalogue_KeysCoverSpecExamples` führt `python` in sortierter Position; die Prüfung „offiziell per Präfix oder Herausgeber als Namensraum, nie beides“ gilt für den neuen Eintrag und braucht keine Änderung |
| Egress-Hosts (Code, Spezifikation, Handbuch) | geprüft, ohne Befund: `pypi.org`, `files.pythonhosted.org` sind in `devcontainer_egress.go:77`, [SPEC-019](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) und Handbuch-Tabelle (`docs/user/benutzerhandbuch.md:1042`) identisch und gleich platziert (nach `java`, vor `sdkman`) |
| gofmt-Ausrichtung der Map | geprüft, ohne Befund: `gofmt -d` im Docker-Image liefert keine Ausgabe; `"python":   {` ist auf die Breite von `"homebrew":` ausgerichtet |
| Doctor und Drift-Detektor | geprüft, ohne Befund: `projectAllFeatureEntries` und `knownProjectableKeys` entstehen aus den Konfigurationseinträgen über den Katalog; `python` erscheint dort nur, wenn der Nutzer den Schlüssel konfiguriert, und erzeugt dann `python:1` wie der Generator. Keine feste Schlüsselliste im Doctor, keine Codeänderung nötig; ein zuvor als verwaist behandelter Schlüssel `python` ist jetzt ein Katalogschlüssel (gewollt) |
| Weitere Schlüssel-Aufzählungen | geprüft, ohne Befund: im Code kommt die Aufzählung nur im Katalog, in `egressFeatureHosts` und in den Tests vor; das Lastenheft führt keine Schlüsselliste und bleibt unverändert; [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) nennt für seinen Entscheidungszeitpunkt acht Schlüssel (accepted, immutable, kein Befund) |
| Doku-Konsistenz | geprüft, ohne Befund: [SPEC-023](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (`python` → `python`), Handbuch-Katalogtabelle mit PEP-668-Hinweis und venv-Anweisung, READMEs „11 built-in features“ beziehungsweise „11 Built-in-Features“ mit `python` in der Liste, CHANGELOG `[Unreleased]` / `Added` mit Hosts; `make docs-check` grün |
| Decken-Regel und Hard Rule 3.4 | geprüft, ohne Befund: Der Diff in `spec/spezifikation.md` enthält weder ADR-, Slice- noch Wellenbezug; `spec/architecture.md` unverändert |
| Replay-Doku | geprüft, ohne Befund: Die Golden-Zeile in `harness/replay.md` nennt `python` samt Render `python:1`; Testnamen und Pfade stimmen mit dem Bestand überein |
| Slice gegen Vorlage | geprüft, ohne Befund außer F-3: Kopf (Welle, Bezug, Berührte Spec-Stellen aufwärts, Verantwortlich, Autor), §1 mit begründeten Ausschlüssen, ≤ 3 Liefer-Punkte in der DoD, §4 Rückführungen, §5, §6 mit Ausgang je Risiko, §7 als offene Vorlage bis zur Closure, §8 mit Sub-Area-Hinweis; Roadmap und Verzeichnis-README führen den Slice; Lastenheft bleibt unverändert und verweist nicht abwärts |
| Slice-Lifecycle | kein Befund: Der Slice liegt in `in-progress/` und wird nach Closure verschoben |
| Architektur, Suppression, Gates | geprüft, ohne Befund: kein neuer Import, kein `//nolint`, keine Gate- oder Coverage-Änderung |
| File Safety, Exit Codes | geprüft, ohne Befund: keine neue Schreib- oder Fehlerpfad-Logik; der Katalogeintrag läuft durch die bestehende Planungsphase |

Nicht geprüft:

- `make gates`, `make test`, `make lint` — nicht ausgeführt; die Findings stützen sich auf Lesen, die Docker-`gofmt`-Prüfung und `make docs-check`.
- Die Messung (Python 3.13.5, `pip` im venv mit genau zwei Hosts, PEP-668-Ablehnung außerhalb des venv) — Image-Build und Netzzugriff liegen außerhalb der Reviewer-Rolle; siehe F-1.
- Verhalten der Feature-Optionen und des Upstream-Features unter aktiver Egress-Firewall — ohne Netz und ohne privilegierten Container nicht prüfbar.
- DoD-Abgleich und Closure — Verifier-Aufgabe.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Messaussage ohne Beleg im Repo · Feature-Tag als Python-Version gelesen · Testzeile ohne LH-Anker

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. Die gofmt-Abweichung aus dem Vorläufer-Lauf kehrt nicht wieder. F-1 bis F-3 sind Doku-Präzision im Slice und Handbuch und können im Verifier-Lauf oder in der Closure-Notiz entschieden werden.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Der Report ersetzt keine Verifikation.
