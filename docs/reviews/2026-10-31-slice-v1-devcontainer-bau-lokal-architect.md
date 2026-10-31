# Review-Report: [`slice-v1-devcontainer-bau-lokal`](../plan/planning/in-progress/slice-v1-devcontainer-bau-lokal.md) — 2026-10-31

**Review-Art:** Plan — Plan-Review der Rolle Architect gegen Spec und ADR (`harness/roles.md`: „Passt die Lösung zu ADRs und Architektur?“). Kein Code-Review und keine DoD-Prüfung; beides gehört Reviewer und Verifier.

**Gegenstand:** Slice-Plan [`slice-v1-devcontainer-bau-lokal`](../plan/planning/in-progress/slice-v1-devcontainer-bau-lokal.md) im Stand des Commits `b371e262` (vor der Einarbeitung dieser Befunde) und die zugehörige [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) (`Accepted`, unverändert gelassen).

**Skill:** keiner; die Rolle Architect (Agent `architect`) führt kein Skill-File unter `.harness/skills/`. Der Auftrag stand im Aufruf. **Modell:** nicht angegeben (Subagent; die Modell-ID liegt dem Auftraggeber nicht vor) · **Datum:** 2026-10-31

**Form:** Die Befunde stammen aus dem Lauf des Architect-Agents. Der Planner hat sie am 2026-10-31 in die Form der Vorlage `review-report.template.md` gebracht (Tabelle, Negativbefunde, Summary, Verdikt); Befunde, Fundstellen und Schweregrade sind unverändert. Lösungswege stehen unter *Übergabe*, nicht in der Befund-Spalte.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- [`slice-v1-devcontainer-bau-lokal`](../plan/planning/in-progress/slice-v1-devcontainer-bau-lokal.md) (§1 bis §8)
- [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) (`Accepted`), [`ADR-0004`](../plan/adr/0004-ci-system.md) (CI-System, Docker-only), [`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md) (Referenzmodell)
- [`LH-NFA-PORT-002`](../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten); [`LH-FA-BUILD-007`](../../spec/lastenheft.md#lh-fa-build-007--docker-only-workflow) (im Lastenheft entfallen, Zeilen 896 bis 898)
- [`SPEC-008`](../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles), [`SPEC-025`](../../spec/spezifikation.md#spec-025--stages-des-multi-stage-dockerfiles), [`SPEC-026`](../../spec/spezifikation.md#spec-026--komposition-der-ci-pipeline-github-actions), [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows), [`SPEC-032`](../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage); `spec/architecture.md`
- `AGENTS.md` (Hard Rules §3.1, §3.4, §3.6, §3.7, §3.10; §5 Regel 3)

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.1` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | INFO | Die Docker-only-Zusage steht nach dem Entfall von [`LH-FA-BUILD-007`](../../spec/lastenheft.md#lh-fa-build-007--docker-only-workflow) nur noch in [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows), die [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) schärft; [`ADR-0004`](../plan/adr/0004-ci-system.md) bleibt wahr, weil die CI im Docker-Modus läuft (`publish.yml` ruft `make build-binaries` ohne `BUILD_MODE`). Die Aussage „identische Build-/Lint-/Test-Pfade lokal und in CI“ aus [`ADR-0004`](../plan/adr/0004-ci-system.md) §Konsequenzen gilt für den Devcontainer nicht mehr uneingeschränkt, was [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) selbst als Negativ benennt. | [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md), [`ADR-0004`](../plan/adr/0004-ci-system.md), [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows), `AGENTS.md` §3.6 | `spec/spezifikation.md` · „ausschließlich `docker build`, `docker run`“ (Zeilen 368 bis 374); `.github/workflows/publish.yml` Zeile 177 | nein — kein Gate prüft den Anker | — |
| F-2 | LOW | Das Feld `Schärft:` der [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) nennt [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) und [`SPEC-008`](../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles); der Slice führt im Kopf zusätzlich [`SPEC-025`](../../spec/spezifikation.md#spec-025--stages-des-multi-stage-dockerfiles) und [`SPEC-032`](../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) als berührt, die Plan-Tabelle §3 aber nur [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) und [`SPEC-008`](../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles). Ob [`SPEC-025`](../../spec/spezifikation.md#spec-025--stages-des-multi-stage-dockerfiles) und [`SPEC-032`](../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) Text ändern müssen, ist offen. | [`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md) Regel 2, [`SPEC-025`](../../spec/spezifikation.md#spec-025--stages-des-multi-stage-dockerfiles), [`SPEC-032`](../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) | [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) Kopf · `Schärft:`; Slice-Kopf · „Berührte Spec-Stellen“ | nein — `make docs-check` prüft die Richtung der Verweise, nicht die Vollständigkeit von `Schärft:` | Schärft-Feld deckt berührte Spec-Stelle nicht |
| F-3 | MEDIUM | Die Plan-Zeile zu `scripts/bau.sh` setzt `CGO_ENABLED=0` und `GOFLAGS` „wie die Stage `deps`“ und lässt die Stage `lint` dasselbe Skript rufen. Die Stage `lint` baut aber auf `golangci/golangci-lint`, nicht auf `deps`, und ruft `golangci-lint run ./...` ohne diese Variablen; mit dem Skript ändert sich der Docker-Modus der Stage, was [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) Entscheidung 4 ausschließt. Die Go-Version im Lint-Image ist nicht geprüft. | [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) Entscheidung 4, `AGENTS.md` §3.10 | `Dockerfile` · „FROM golangci/golangci-lint:${GOLANGCI_LINT_VERSION} AS lint“ (Zeilen 71 bis 76); `ENV GOFLAGS` in `deps` (Zeilen 42 bis 46) | ja — `make lint` im Docker-Modus vor und nach der Änderung | Gemeinsames Skript ändert eine Stage im Docker-Modus |
| F-4 | MEDIUM | Go und `golangci-lint` stehen künftig an drei Stellen (`Dockerfile`, `repo.mk`, `.devcontainer/Dockerfile`); der Plan nennt weder einen Kopplungskommentar „Hebung betrifft drei Stellen“ noch, dass der handgeschriebene Teil hinter dem `END`-Marker nach `USER vscode` liegt und `USER root` braucht. Dass handgeschriebene Zeilen `u-boot generate devcontainer` überstehen, ist in ADR und Plan behauptet, hier nicht nachgemessen. | [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) Entscheidung 5, `AGENTS.md` §3.7 | `Dockerfile` Zeilen 22 bis 24 (Pin-Policy); `repo.mk` Zeilen 43 bis 45; `.devcontainer/Dockerfile` Zeilen 28 bis 55 (Muster `ARG …_SHA256`) und Zeile 75 (`END`-Marker) | nein — der Abgleich der drei Stellen fehlt bis zum d-check-Modul (Issue 14 in `pt9912/d-check`) | Pin an mehreren Stellen ohne Kopplungskommentar |
| F-5 | LOW | Im Devcontainer (`BUILD_MODE=local`) baut `make packages` die Binaries lokal und die Pakete per `docker run`; ein solches Paket ist kein Release-Artefakt. `make test PKG=… RUN=…` ist ein Teillauf, und der Plan sagt nicht, dass `make gates` `PKG` und `RUN` ablehnt oder leert. | [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) Entscheidungen 3 und 6, [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) | `repo.mk` Zeile 167 (`packages: build-binaries`) und Zeilen 170 bis 177; `scripts/build-packages.sh` Zeilen 3 bis 4 und 29 | ja — `make modus-pruefung` (geplant) | Teillauf oder Release-Artefakt im falschen Modus ohne Sperre |
| F-6 | LOW | Die Plan-Tabelle §3 führt `docs/plan/adr/README.md` als „update: Indexzeile für ADR 0037“, die Zeile steht aber schon in Zeile 65 (Commit `b371e262`). | `AGENTS.md` §3.10 | `docs/plan/adr/README.md` · Zeile 65 | ja — `git diff --stat` gegen §3 nach der Umsetzung | Plan-Zeile ohne Diff |
| F-7 | LOW | `spec/architecture.md` Zeilen 714 und 716 beschreiben `make test` als Stage im `Dockerfile` und ein getrenntes Make-Target in der CI; die Datei fehlt in der Plan-Tabelle. Für `.dockerignore` fehlt der Pfadname des Workspace-Caches. | `AGENTS.md` §3.10, [`SPEC-025`](../../spec/spezifikation.md#spec-025--stages-des-multi-stage-dockerfiles) | `spec/architecture.md` · „`make test` (Stage `test` im Dockerfile …)“ (Zeile 714); `.dockerignore` | ja — `git diff --stat` gegen §3 nach der Umsetzung | Berührte Datei fehlt in der Plan-Tabelle |
| F-8 | INFO | Die Kopf-Verweise der [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) laufen aufwärts (Lastenheft, Spec) und zu einer aktiven ADR; die ADR verlinkt keinen Slice, der Slice verlinkt ADR und Spec. Die Decken-Regel gilt für den Nachzug in [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) und [`SPEC-008`](../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles): Er nennt weder ADR noch Slice. | [`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md), [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md), `AGENTS.md` §3.4 | [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) Kopf · `Bezug:`, `Schärft:` | ja — `make docs-check` | — |
| F-9 | INFO | Die Stichproben tragen: `repo.mk` Zeilen 77 bis 79, 100, 103, 142; `Dockerfile` (sechs `COPY . .`, kein `--mount=type=cache`, Stages `compile`, `lint`, `test`, `coverage`); [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) ab Zeile 368, [`SPEC-008`](../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) ab Zeile 1146; `tools/mk/devcontainer.mk`. `test-docker`, `govulncheck` und `build-binaries` stehen in keiner Tabelle von [`SPEC-008`](../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles); [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) Entscheidung 3 nennt sie. | `AGENTS.md` §3.10, [`SPEC-008`](../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles), [`SPEC-026`](../../spec/spezifikation.md#spec-026--komposition-der-ci-pipeline-github-actions) | `repo.mk`; `Dockerfile`; `spec/spezifikation.md`; `tools/mk/devcontainer.mk` | ja — Lesen am Ist-Stand (kein Gate) | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Widerspruch zu [`ADR-0004`](../plan/adr/0004-ci-system.md) und [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) | geprüft, ohne Befund: kein Widerspruch, [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) reicht als Plan-Anker nach `AGENTS.md` §3.6 (F-1) |
| Lastenheft-Änderung | geprüft, ohne Befund: nicht nötig, weil [`LH-FA-BUILD-007`](../../spec/lastenheft.md#lh-fa-build-007--docker-only-workflow) entfallen ist und [`LH-NFA-PORT-002`](../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten) nur „möglichst wenige Systemabhängigkeiten“ verlangt |
| Referenzmodell ([`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md)) | geprüft, ohne Befund (F-8) |
| Schichten und Importregeln (`internal/hexagon`, depguard, `spec/architecture.md` Zeile 490) | geprüft, ohne Befund: der Slice ändert keinen Go-Quelltext |
| Folge-ADR | geprüft, keine heute nötig; Auslöser einer künftigen stehen unter *Übergabe* |
| Messwerte der [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) (Zeiten, Kontextgröße 24 MB) | **nicht geprüft** |
| Verhalten von `u-boot generate devcontainer` (handgeschriebene Zeilen bleiben stehen) | **nicht geprüft** |
| `docs/maintainer/quality.md` Zeilen 72 und 300 bis 314 | **nicht geprüft** (nur Zeilen 17 bis 32 gelesen) |
| Go-Version und `bash` im Image `golangci/golangci-lint:v2.14.0` | **nicht geprüft** |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 4 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Gemeinsames Skript ändert eine Stage im Docker-Modus · Pin an mehreren Stellen ohne Kopplungskommentar · Teillauf oder Release-Artefakt im falschen Modus ohne Sperre · Schärft-Feld deckt berührte Spec-Stelle nicht · Plan-Zeile ohne Diff · Berührte Datei fehlt in der Plan-Tabelle

## Verdikt

**Merge-blockierend:** nein — kein HIGH. Die zwei MEDIUM (F-3, F-4) sind Zusätze zum Plan ohne neue Entscheidung. Voraussetzung für `open` → `next` nach dem Architect: F-3 und F-5 im Plan vermerkt, F-6 und F-7 berichtigt; F-2 und F-4 sind Notizen für die Umsetzung. Alles davon ist am 2026-10-31 in §3 und §6 des Slice eingearbeitet.

**Übergabe:** Findings gehen an den Implementer (Rückkante
Review → Plan bei Plan-Defekt); die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).

Lösungswege des Architects zu den Befunden:

- F-2: Der Plan führt [`SPEC-025`](../../spec/spezifikation.md#spec-025--stages-des-multi-stage-dockerfiles) und [`SPEC-032`](../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) in §3 als „unverändert, gelesen“ mit Fundstelle. Werden Test und Coverage ein Lauf, schärft eine kurze Folge-ADR [`SPEC-032`](../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) (nicht jetzt).
- F-3: §3 benennt, welche Umgebungsvariablen das Skript je Stage selbst setzt und welche die Umgebung liefert; `lint` bleibt unverändert oder wird bewusst angeglichen und in §7 vermerkt. M0 prüft die Go-Version des Lint-Images gegen `GO_VERSION` und ob `bash` dort vorhanden ist.
- F-4: Ein Kopplungskommentar „Hebung betrifft drei Stellen“ in `repo.mk` und `Dockerfile` (kein Prüfskript, Entscheidung des Projektinhabers); der Plan nennt `USER root` und `USER vscode` nach dem `END`-Marker. Eine Folge-ADR ist nicht nötig.
- F-5: Release-Artefakte entstehen nur im Docker-Modus, ein Satz in [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows); `make gates` lehnt `PKG` und `RUN` ab oder leert sie, und `modus-pruefung` prüft das.
- F-6 und F-7: Plan-Zeile auf „unverändert, Zeile 65 gelesen“ setzen; `spec/architecture.md` Zeile 714 erhält „im Docker-Modus“ oder wird als unverändert mit gelesener Stelle ausgewiesen.
- Künftige Folge-ADR nur bei Eintritt: (a) Test und Coverage werden ein Lauf, dann schärft eine ADR „Coverage-Lauf teilt den Test-Lauf im lokalen Modus“ [`SPEC-032`](../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage); (b) die CI wechselt auf das Toolchain-Image (Option E), dann ersetzt eine neue ADR mit `Supersedes` auf [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md), wie dessen Re-Evaluierungs-Trigger vorsieht.
