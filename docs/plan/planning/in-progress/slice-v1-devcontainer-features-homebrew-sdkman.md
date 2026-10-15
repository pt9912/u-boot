# Slice V1: Devcontainer-Features `homebrew` und `sdkman` im eingebauten Katalog

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** v0.9.0 (Devcontainer-Features Homebrew und SDKMAN).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
[`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte),
[`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion).
Folge-ADR: Drittanbieter-Feature im eingebauten Katalog.

**Berührte Spec-Stellen:** `spezifikation.md` [SPEC-023](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (eingebaute Features) und [SPEC-019](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Egress-Hosts je Feature); das Lastenheft bleibt unverändert ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) nennt SDKMAN als Beispiel und erlaubt „lokal hinterlegte“ Features).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-15.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Schlüssel `homebrew` und `sdkman` sind eingebaute Devcontainer-Features: `homebrew` rendert das Community-Feature `ghcr.io/meaningful-ooo/devcontainer-features/homebrew:2`, `sdkman` das offizielle Feature `ghcr.io/devcontainers/features/java:1` (installiert SDKMAN mit Java); das Drittanbieter-Feature wird in der Befehlsausgabe ausgewiesen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Feature-Optionen (etwa `candidate` und `version` des Community-Features `ebaskoro/sdkman`) — der Katalog kennt keine Optionen, und das Community-Feature installiert einen SDK-Kandidaten, nicht SDKMAN selbst (Messung in der Folge-ADR); `sdkman` zeigt deshalb auf das offizielle `java`-Feature.
- Weitere Community-Features (`kubectl`-Varianten, `...-homebrew`-Features anderer Autoren) — jede Quelle braucht eine eigene Entscheidung; bis dahin gelten sie als extern mit Freigabe nach [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features).
- Homebrew als Installationsweg von `u-boot` — der Tap besteht bereits ([`slice-v2-homebrew-formula`](../in-progress/slice-v2-homebrew-formula.md)).

## 2. Definition of Done

- [x] **Folge-ADR** (nächste freie Nummer) entscheidet Drittanbieter-Features im Katalog (Quelle, Versionspin, Kennzeichnung, Messung); ADR-Index aktualisiert.
- [x] **Katalog und Generator:** `homebrew` und `sdkman` im Katalog; identische Render-Schlüssel (`java` und `sdkman` zusammen) erscheinen einmal in `devcontainer.json`; die Befehlsausgabe weist das Drittanbieter-Feature aus; Egress-Hosts je Feature; Golden Cases belegen Render, Warnung, Dedupe und Hosts.
- [x] **Doku:** Spezifikation ([SPEC-023](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [SPEC-019](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)), Benutzerhandbuch, READMEs und CHANGELOG nennen die neuen Schlüssel.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — kein Self-Review.
- [x] Doku-Update, falls ein öffentlicher Vertrag berührt ist.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/` (Folge-ADR, Index) | neu / update | Drittanbieter-Feature im Katalog |
| `internal/hexagon/application/devcontainer_features.go` | update | Katalogeinträge, Herausgeber-Kennzeichnung, Dedupe gleicher Render-Schlüssel |
| `internal/hexagon/application/generate.go`, `devcontainer_sandbox.go` | update | Warnung für Drittanbieter-Features |
| `internal/hexagon/application/devcontainer_egress.go` | update | Hosts für `homebrew` und `sdkman` |
| `internal/hexagon/application/*_test.go`, `internal/hexagon/domain/featurename_test.go` | update | Katalog, Golden Cases |
| `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | update | öffentlicher Vertrag und Golden-Zeile |

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; kein anderer Slice des Implementers in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice überschreitet 3 Liefer-Punkte.
- `in-progress` → `open` (blockiert): das Community-Feature `homebrew` ist nicht mehr verfügbar oder baut nicht mehr.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Das Community-Feature `homebrew` ändert sich oder verschwindet (Drittanbieter, Major-Pin `2`) — **Ausgang:** weiter offen: Re-Evaluierungs-Trigger der Folge-ADR ([ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md)); ein Build-Fehler fällt beim Image-Build auf.
- Die Egress-Hosts reichen für `brew install` oder `sdk install` nicht aus (Download-Hosts der Hersteller wechseln) — **Ausgang:** entfallen: gemessen für `brew install jq` sowie `sdk list java`, `sdk install maven` und `sdk install gradle` ([ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md)); weitere Hosts ergänzt der Nutzer über `devcontainer.sandbox.egress.allow`, das Handbuch nennt es.

## 7. Closure-Notiz

- **Was hat funktioniert:** `homebrew` (Community-Feature, gekennzeichnet) und `sdkman` (Alias auf `java`) sind im Katalog; Render, Dedupe, Warnung und Egress-Hosts sind durch Golden Cases gepinnt; die Hosts sind gemessen (DNS-Abfragen von `brew install jq` und `sdk install maven/gradle`).
- **Was ging anders als geplant:** Die Messung mit der Dev-Container-CLI zeigte, dass das Community-Feature `sdkman` einen SDK-Kandidaten statt SDKMAN installiert; `sdkman` zeigt deshalb auf das offizielle `java`-Feature (Entscheidung des Projektinhabers). Der Review [`2026-10-15-devcontainer-features-homebrew-sdkman`](../../../reviews/2026-10-15-devcontainer-features-homebrew-sdkman.md): 0 HIGH, 2 MEDIUM, 5 LOW; behoben bis auf F-1 (Zustimmungsmodell, Entscheidung des Projektinhabers: Kennzeichnung durch Warnung bei Aktivierung im Katalog).
- **Steering-Loop-Eintrag:** Eine Annahme über ein externes Artefakt (Funktion eines Community-Features) vor der Planung messen, nicht aus dem Namen ableiten — die Messung änderte hier den Zuschnitt.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angefallen.
- **Folge-Slices:** keine.
- **Risiken aus §6:** siehe §6
- **Drei Paarungen:** Anker (Zustimmungsmodell Drittanbieter-Feature, Review F-1) → Entscheidung des Projektinhabers im [ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md) · Register: keine Beobachtung angefallen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `internal/hexagon/application`, `internal/hexagon/domain` (Greenfield laut `harness/conventions.md`), `docs/plan/`, `docs/user/`, `spec/`.

**Vorgelagert — offene Beobachtungen sichten:** kein Beobachtungs-Register
mit Treffern für die berührten Sub-Areas.

Alle berührten Sub-Areas GF.

## Verification Evidence

Scope:
- Slice: [slice-v1-devcontainer-features-homebrew-sdkman](slice-v1-devcontainer-features-homebrew-sdkman.md)
- IDs: [LH-FA-DEV-003](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [LH-NFA-SEC-004](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [LH-FA-DEV-008](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md), [SPEC-023](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [SPEC-019](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)
- Artefakte: `internal/hexagon/application/devcontainer_features.go`, `devcontainer_egress.go`, `generate.go`, `generate_features_test.go`, `generate_sandbox_test.go`, `doctor_features_drift_test.go`; `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md`, ADR-Index; Review-Report [2026-10-15-devcontainer-features-homebrew-sdkman](../../../reviews/2026-10-15-devcontainer-features-homebrew-sdkman.md)

DoD-Abgleich (gegen Diff `7f1927e^..ac59b83` gelesen):
- [x] Folge-ADR — [ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md) liegt vor, Status Accepted, Eintrag im ADR-Index (`docs/plan/adr/README.md`).
- [x] Katalog und Generator — Katalogeinträge, Dedupe und Warnung in `devcontainer_features.go`/`generate.go`, Hosts in `devcontainer_egress.go`; Tests `TestGenerateDevcontainer_FeaturesHomebrewSdkman`, `…_FeaturesHomebrew_WarningNamesConfiguredVersion`, `…_SandboxEgress_FeatureHosts`. Abweichung von Plan §3: `devcontainer_sandbox.go` und `domain/featurename_test.go` sind unverändert (Warnung liegt in `generate.go`); ohne Wirkung auf den Vertrag.
- [x] Doku — [SPEC-023](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) und [SPEC-019](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) nennen Schlüssel und Hosts; Handbuch (Katalogtabelle, Egress-Tabelle), beide READMEs, CHANGELOG `[Unreleased]` nennen sie. Teilabweichung: siehe Abweichungen (§6-Zusage Handbuch).
- [x] `make gates` — grün (siehe Sensors).
- [x] Review — Report liegt vor (0 HIGH, 2 MEDIUM, 5 LOW, 2 INFO); F-2 bis F-7 durch `ac59b83` behoben (Tests `…_WarningNamesConfiguredVersion`, `…_JavaVersionAndSdkman_RenderTwoKeys`, `TestDoctor_FeaturesDrift_OKWithJavaAndSdkmanSharingOneKey`, gofmt in `devcontainer_egress.go`); F-1 trägt §7 als Entscheidung des Projektinhabers (offene Planner-Frage unten).
- [x] Doku-Update bei berührtem öffentlichem Vertrag — siehe Doku-Punkt.
- [x] Closure-Notiz mit Steering-Loop-Eintrag — §7 vorhanden.
- [x] Beobachtungs-Register — §7: „keine Beobachtung angefallen“.
- [x] §6-Risiken — beide tragen einen Ausgang (weiter offen mit Re-Evaluierungs-Trigger der ADR; entfallen durch DNS-Messung laut ADR).
- [x] Drei Paarungen — in §7 benannt (Anker, Register); Folge-Slices: keine.

Abweichungen und offene Punkte:
- [Review F-4](../../../reviews/2026-10-15-devcontainer-features-homebrew-sdkman.md) nicht umgesetzt: §6 sagt „das Handbuch nennt es“ (weitere Download-Hosts über `devcontainer.sandbox.egress.allow`); `grep -i "gradle\|Kandidat" docs/user/benutzerhandbuch.md` ist leer, `ac59b83` änderte im Handbuch nur die Katalog-Einleitung (Zeile 751). Der Hinweis steht im CHANGELOG und allgemein beim Egress-Abschnitt („Wichtig“ zu `egress.allow`), nicht für SDKMAN-Kandidaten oder `brew`-Formeln.
- Review F-5 (Überschrift „Ohne weitere Freigabe aktivierbar“) ist durch den Klammerzusatz im selben Satz abgedeckt, die Überschrift selbst bleibt unverändert.
- Offene Planner-Frage (Review F-1): Gilt „Aktivierung plus Warnung nach dem Schreiben“ als Zustimmung nach [LH-NFA-SEC-004](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)? [ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md) Entscheidung 2 trägt die Auslegung; [SPEC-023](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) schließt weiter mit „Jede andere Quelle gilt als extern und braucht eine Freigabe“, das Lastenheft ist unverändert. Entscheidung liegt beim Projektinhaber.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | Exit 0, „[gates] lint + test + coverage-gate + docs-check green“ |
| `make docs-check` | pass | d-check: 204 Datei(en) geprüft, 0 Befund(e) |
| `make test-docker` | pass | Exit 0, Paket `internal/adapter/driven/docker` ok (21.5 s), übrige Pakete ok, lokale Colima-Engine |

Traceability:
| ID | Beleg |
| --- | --- |
| [LH-FA-DEV-003](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | `TestGenerateDevcontainer_FeaturesHomebrewSdkman`, `…_FeaturesJavaVersionAndSdkman_RenderTwoKeys`, `TestDoctor_FeaturesDrift_OKWithJavaAndSdkmanSharingOneKey`; [SPEC-023](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)-Absatz (`spec/spezifikation.md:890`) |
| [LH-NFA-SEC-004](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | `TestGenerateDevcontainer_FeaturesHomebrew_WarningNamesConfiguredVersion`, `…_ExternalOverrideNoWarning`; Auslegung offen (F-1) |
| [LH-FA-DEV-008](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) | `TestGenerateDevcontainer_SandboxEgress_FeatureHosts`; Hosttabelle in [SPEC-019](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) und Handbuch |
| [ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md) | ADR Accepted, Index-Eintrag, Review-Report; die DNS-Messreihe der ADR wurde nicht erneut ausgeführt |

Carveouts:
- Neu: none
- Geloest: none
- Unveraendert: none (`docs/plan/planning/in-progress/carveouts.md` nicht berührt)

Nicht ausgefuehrt:
- `make ci` / `make fullbuild` — vom Verifier nicht selbst belegbar; der Implementer meldet einen lokalen Lauf, der hier nicht reproduziert wurde.
- Erneute Messung der Egress-Hosts (`brew install jq`, `sdk install maven/gradle`) — nicht Teil der Verifier-Sensors; Beleg bleibt die Messreihe in [ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md).
- Tests mit `--privileged` außer `make test-docker` — nicht ausgeführt.

Commit / Artefakt:
- Delivery-Hash: `ac59b83` (Implementation `95609ce`, Review-Fixes `ac59b83`)

Replay / Golden:
- Affected generators: `generate devcontainer` (Feature-Katalog, `devcontainer.json`-Features, Sandbox-Egress-Hosts, Warnungsausgabe)
- Golden cases added: `TestGenerateDevcontainer_FeaturesHomebrewSdkman`, `TestGenerateDevcontainer_FeaturesHomebrew_WarningNamesConfiguredVersion`, `TestGenerateDevcontainer_FeaturesJavaVersionAndSdkman_RenderTwoKeys`, `TestDoctor_FeaturesDrift_OKWithJavaAndSdkmanSharingOneKey`; Zeile in `harness/replay.md:54`
- Golden cases updated: `devcontainer_features_test.go` (Katalogumfang), `TestGenerateDevcontainer_SandboxEgress_FeatureHosts` (Hosts je Feature)
- Golden cases replayed: komplette Suite über `make gates` (Exit 0), Double-Run-NoOp und Fresh-State im Fall `homebrew`/`sdkman` gepinnt; Stale-Block und Konflikt ohne Write für den neuen Fall nur über die allgemeinen `TestGenerateDevcontainer_Features*`-Tests (Review F-8)
- Intentional output changes: `devcontainer.json` ändert sich nur für die neuen Schlüssel `homebrew` (`ghcr.io/meaningful-ooo/devcontainer-features/homebrew:2`) und `sdkman` (`ghcr.io/devcontainers/features/java:1`, einmal bei gemeinsamer Aktivierung mit `java`); neue Warnung für `homebrew`; neue Egress-Hosts nur bei aktivem Feature
