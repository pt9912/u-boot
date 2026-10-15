# Review-Report: [`slice-v1-devcontainer-features-homebrew-sdkman`](../plan/planning/done/slice-v1-devcontainer-features-homebrew-sdkman.md) — 2026-10-15

**Review-Art:** Code — Diff gegen Slice-Plan, [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md), Spec ([`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) und Replay Rules (`harness/replay.md` §Replay Rules).

**Gegenstand:** Diff-Range `e6b003c..HEAD` (Commits `95609ce`, `98a6828`, Plan-Commits davor), Pfade `internal`, `spec`, `docs/plan/adr`, `docs/user`, READMEs, CHANGELOG, `harness`.

**Skill:** `.harness/skills/reviewer.md` @ 98a6828 · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-15

**Eingangs-Kontext:**

- [`slice-v1-devcontainer-features-homebrew-sdkman`](../plan/planning/done/slice-v1-devcontainer-features-homebrew-sdkman.md) (DoD, Abgrenzung, Risiken)
- [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md), [ADR-0012](../plan/adr/0012-devcontainer-egress-firewall.md)
- [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)
- [SPEC-019](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [SPEC-023](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)
- `AGENTS.md` (Hard Rules), `harness/review.md`, `harness/roles.md`, `harness/replay.md`

Zusätzlich ausgeführt: `make docs-check` (203 Dateien, 0 Befunde); `gofmt -l` / `gofmt -d` im Image `golang:1.27.2` ohne Netz, Quellen read-only gemountet (kein Host-Go-Nachweis, nur Beleg für F-3). `make gates` und `make test` wurden nicht ausgeführt.

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Ein Drittanbieter-Feature (`ghcr.io/meaningful-ooo/...`) steht ohne Eintrag in `devcontainer.featureSources.allow` im Katalog und wird allein durch `features.homebrew.enabled: true` gerendert. Die Kennzeichnung ist eine Warnung, die `generate devcontainer` nach dem Schreiben ausgibt; es gibt keine Bestätigung vor dem Schreiben und keinen Hinweis bei `config set`. Der Lastenheft-Text bleibt unverändert und verlangt weiter „lokal hinterlegt oder ausdrücklich freigegeben“ beziehungsweise „ausdrücklich in `u-boot.yaml` freigegeben (siehe `devcontainer.featureSources.allow`)“; [SPEC-023](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) schließt weiter mit „Jede andere Quelle gilt als extern und braucht eine Freigabe“. Ob „aktivieren plus Warnung“ die Zustimmung ist, entscheidet nur [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md), das im Rang unter dem Lastenheft steht. | [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) | `internal/hexagon/application/devcontainer_features.go:217-222`, `:321-349`; `internal/hexagon/application/generate.go:782`; `spec/spezifikation.md:890` | nein — Review-only (Auslegung des Lastenhefts); der Pfad `config set` ohne Hinweis ist über `maybeWarnOrphanFeatureActivation` in `internal/hexagon/application/config.go:320-345` lesbar | Zustimmung nur nachträglich per Warnung |
| F-2 | MEDIUM | `thirdPartyFeatureWarnings` nennt immer `cat.defaultVersion`. Setzt der Nutzer `features.homebrew.version: "3"` ohne `source`, rendert `projectFeatureSourceVersion` `…/homebrew:3`, die Warnung nennt `…/homebrew:2`. [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) sagt, die Warnung nenne „Herausgeber, Quelle und Version“; kein Test setzt eine Versionsüberschreibung. | [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) Entscheidung 2, [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | `internal/hexagon/application/devcontainer_features.go:346` (gegen `:413`) | ja — `make test` mit Fixture `homebrew` plus `version: "3"` und Abgleich Warnung gegen `devcontainer.json` | Warnung weicht vom gerenderten Wert ab |
| F-3 | LOW | `egressFeatureHosts` ist nicht gofmt-konform: Die neuen Schlüssel `homebrew` und `sdkman` sind länger als `node`, `go`, `java`; `gofmt -d` richtet die drei bestehenden Zeilen neu aus. Die Lint-Konfiguration enthält kein gofmt/gofumpt, `make lint` fängt es nicht; `gofmt -l` meldet 20 Dateien im Paket, darunter solche außerhalb des Diffs. | Maintainability | `internal/hexagon/application/devcontainer_egress.go:74-79` | ja — `gofmt -d` im Docker-Image (siehe Kopf); kein Gate | Format-Drift ohne Gate |
| F-4 | LOW | Die Zusage „die Doku nennt das“ ([ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) Entscheidung 5, Slice §6 Ausgang des Egress-Risikos) steht nur im CHANGELOG. Das Handbuch führt bei der Egress-Tabelle und bei den Katalog-Zeilen keinen Hinweis, dass Download-Hosts weiterer SDKMAN-Kandidaten (etwa `services.gradle.org`) oder Hersteller-Hosts von `brew`-Formeln in `devcontainer.sandbox.egress.allow` gehören; `grep` nach `gradle`/`Kandidat` in `docs/user` ist leer. | Doku-Drift, Slice §6 | `docs/user/benutzerhandbuch.md:763-764`, `:1041-1042` | ja — Gegenlesen | Nutzerführung für Egress fehlt im Handbuch |
| F-5 | LOW | Die Katalog-Überschrift im Handbuch sagt „Ohne weitere Freigabe aktivierbar“, die Tabelle führt darunter jetzt `homebrew` als Drittanbieter-Feature; die Zeile selbst nennt die Besonderheit, die Überschrift bleibt allgemein. | Doku-Präzision | `docs/user/benutzerhandbuch.md:750` | ja — Gegenlesen | Überschrift und Tabelle uneinheitlich |
| F-6 | LOW | Das Doctor-Drift-Verhalten bei `java` plus `sdkman` ist korrekt (Mengen-Semantik in `projectAllFeatureEntries`, `expectedKeys` und `knownProjectableKeys` führen den gemeinsamen Schlüssel `java:1` einmal), aber nicht gepinnt: Kein Test in `doctor_features_drift_test.go` deckt `sdkman`, weder „nur `sdkman` aktiv“ noch „`java` aus, `sdkman` an, JSON trägt `java:1`“. Ein späterer Umbau der Projektion auf Listen oder Zähler könnte falschen Drift melden, ohne dass ein Test rot wird. | [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Replay Rules (Idempotenz) | `internal/hexagon/application/doctor.go:1413`; Lücke in `internal/hexagon/application/doctor_features_drift_test.go` | ja — `make test` mit ergänztem Drift-Case | Gemeinsamer Render-Schlüssel nicht im Drift-Test |
| F-7 | LOW | `dedupeFeatureData` entfernt nur identische Schlüssel. Aktiviert der Nutzer `java` mit `version: "21"` und `sdkman` (Default `1`), erscheinen `java:1` und `java:21` beide in `devcontainer.json`. [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) Entscheidung 4 deckt nur „gleiche Render-Schlüssel“; Verhalten und Test für diesen Fall fehlen, Doku erwähnt ihn nicht. | Maintainability, [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) | `internal/hexagon/application/devcontainer_features.go:307-317` | ja — `make test` mit Fixture (Verhalten ist derzeit unspezifiziert) | Zwei Versionen desselben Features |
| F-8 | INFO | Replay-Abgleich: Die neue Golden-Zeile in `harness/replay.md` stimmt mit den Testnamen überein; Fresh-State und Double-Run-NoOp sind gepinnt. Stale-Block und User-Content sowie Konflikt ohne Write für den neuen Fall fehlen als eigene Cases; sie sind über die allgemeinen Feature-Tests (`TestGenerateDevcontainer_Features*`) abgedeckt, der Katalogeintrag ändert die Planungsphase nicht. Die Tests prüfen Code und Substring der Warnung (`meaningful-ooo`), nicht `Level` und `Subject`. | `harness/replay.md` §Replay Rules | `internal/hexagon/application/generate_features_test.go:630-676` | ja — Gegenlesen | — |
| F-9 | INFO | Plan und Diff weichen ab: Slice §3 nennt `devcontainer_sandbox.go` (Warnung) und `internal/hexagon/domain/featurename_test.go` als Änderungsziel, der Diff ändert beide nicht (Warnung in `generate.go`). DoD-Punkte Review, `make gates`, Closure-Notiz, Beobachtungs-Register sind offen; das ist Sache des Verifiers. [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) trägt eine Messreihe und benennt „Nicht verifiziert“ ausdrücklich; die Egress-Allowlist selbst wurde dort nicht unter aktiver Firewall gemessen (die Messung ist DNS-basiert) und das ADR behauptet das auch nicht. | `AGENTS.md` §6, `harness/verification.md` | Slice-Datei, §3 und §2 | Review-only | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `dedupeFeatureData` (Korrektheit) | geprüft, ohne Befund: Die Sortierung läuft über (Source, Version), die Struktur hat genau diese zwei Felder; gleiche Schlüssel liegen damit immer benachbart. Das In-Place-Filtern (`sorted[:0]`) ist sicher, weil der Aufrufer den Slice nicht weiterverwendet |
| Doctor-Drift bei `java`/`sdkman` | geprüft, ohne Befund im Code: nur `java` aktiv, nur `sdkman` aktiv, beide aktiv, `java` aus und `sdkman` an liefern jeweils denselben erwarteten Schlüssel `java:1`; kein falsches „fehlender Key“ und kein falsches „Orphan“ (Test-Lücke siehe F-6) |
| `thirdPartyFeatureWarnings` (Sortierung, Override, Dubletten) | geprüft, ohne Befund: Namen werden sortiert, ein `source:`-Override unterdrückt die Warnung (Test `…_ExternalOverrideNoWarning`), je Katalogname höchstens eine Warnung, `sdkman` erzeugt keine; Abweichung bei Versionsüberschreibung siehe F-2 |
| `init --devcontainer` (auch `--sandbox`) | geprüft, ohne Befund: `initproject.go` ruft `devcontainerTemplateData(name, nil, …)` ohne Features; `collectDevcontainerFeatures` hat genau einen Aufrufer (`generate.go:771`). Der Init-Pfad rendert kein Katalog-Feature und braucht die Warnung nicht |
| `generate devcontainer --sandbox` und `--dry-run` | geprüft, ohne Befund: Beide laufen durch `generateDevcontainer`; die Warnung wird vor Planung angehängt und in der Antwort auch bei NoOp/Dry-Run geliefert; die CLI druckt `Warning (code): …` und mappt sie in das JSON |
| Egress-Hosts (`devcontainer_egress.go`, Spezifikation, Handbuch) | geprüft, ohne Befund bis auf F-3 (Format) und F-4 (Hinweis): Hosts in Code, [SPEC-019](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), Handbuch und [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) sind identisch; Ausgabe wird als Menge gesammelt und sortiert, `java` plus `sdkman` erzeugen keine Dubletten; die gemessenen GitHub-Hosts der Brew-Läufe sind bereits in der Basisliste |
| Katalog-Invarianten-Test | geprüft, ohne Befund: Die Regel „offiziell per Präfix oder Herausgeber als Registry-Namensraum, nie beides“ pinnt beide Richtungen; `sdkman` zeigt auf den offiziellen Namensraum ohne Herausgeber |
| Decken-Regel und Hard Rule 3.4 | geprüft, ohne Befund: `spec/spezifikation.md` verweist in der Diff-Zeile nicht auf ADR, Slice oder Welle; `spec/architecture.md` unverändert |
| Drittanbieter-ADR (Format, Index, Messaussagen) | geprüft, ohne Befund: MADR-Felder inklusive `Schärft:` vorhanden, ADR-Index aktualisiert, Messtabelle ist im ADR belegt, offene Punkte stehen unter „Nicht verifiziert“; nichts als gemessen behauptet, was nur Annahme ist. Die Auslegungsfrage zum Lastenheft steht in F-1 |
| Doku-Konsistenz (CHANGELOG, READMEs, Handbuch-Tabelle) | geprüft, ohne Befund außer F-4/F-5: CHANGELOG `[Unreleased]`/`Added` mit Verweisen, beide READMEs „10 built-in features“ beziehungsweise „10 Built-in-Features“ mit Drittanbieter-Kennzeichnung, Handbuch-Tabelle vollständig; `make docs-check` grün |
| Architektur und Importregeln | geprüft, ohne Befund: neuer Import `port/driving` in `application` ist die bestehende Richtung (Hexagon-intern); kein Adapter-Zugriff |
| Suppression-Verbot, Docker-only, Gates | geprüft, ohne Befund: kein `//nolint`, keine Gate-Änderung |
| Slice-Lifecycle | kein Befund: Der Slice liegt in `in-progress/` und wird nach Closure verschoben |

Nicht geprüft:

- `make gates`, `make test`, `make lint` — nicht ausgeführt; die Findings stützen sich auf Lesen, die Docker-`gofmt`-Prüfung und `make docs-check`.
- Tatsächlicher Image-Build mit `homebrew:2` und `java:1` sowie `brew install`/`sdk install` unter aktiver Egress-Firewall — außerhalb der Reviewer-Rolle und ohne privilegierten Container; das ADR benennt den Rest selbst als nicht verifiziert.
- Aktuelle Verfügbarkeit des Community-Features bei ghcr.io — kein Netzzugriff vorgesehen; Re-Evaluierungs-Trigger im ADR.
- DoD-Abgleich und Closure — Verifier-Aufgabe.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 5 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zustimmung nur nachträglich per Warnung · Warnung weicht vom gerenderten Wert ab · Format-Drift ohne Gate · Nutzerführung für Egress fehlt im Handbuch · Gemeinsamer Render-Schlüssel nicht im Drift-Test

## Verdikt

**Merge-blockierend:** nein — es gibt kein HIGH. F-1 und F-2 sind vor dem Merge zu klären: F-1 ist eine Auslegungs- und Spec-Frage, die der Projektinhaber über [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) oder eine Lastenheft-Klarstellung entscheidet; F-2 ist ein kleiner, testbarer Fehler in der Sicherheitsausgabe.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Der Report ersetzt keine Verifikation.
