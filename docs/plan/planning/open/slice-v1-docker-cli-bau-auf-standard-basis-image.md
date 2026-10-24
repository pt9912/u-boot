# Slice V1: `docker-cli` baut mit seiner Standard-Option nicht auf dem Standard-Basis-Image — Feature-Optionen in `u-boot.yaml`, Pin des Basis-Images oder ein Hinweis (Nebenbefund N4 von `slice-v1-sandbox-features-socket-und-lockerungen`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Nebenbefund N4 von [`slice-v1-sandbox-features-socket-und-lockerungen`](../done/slice-v1-sandbox-features-socket-und-lockerungen.md) (dort §3 *Offene Punkte* 1 und §7); der Weg ist eine Entscheidung des Projektinhabers. Es gibt keine Closure-Bedingung, die von der DoD dieses Slice verschieden wäre. Die Entscheidung zum Weg liegt vor dem Text des Advisory von `v0.13.1` (Verification-Report des Vorgänger-Slice, V-3 und Prüfliste des Release-Cuts); der Slice selbst ist nicht Teil des Release. *Vermerk 2026-10-22 (Entscheidung des Projektinhabers, [`slice-v1-release-cut-v0.13.1`](../done/slice-v1-release-cut-v0.13.1.md) §7 Entscheidung 4):* N4 bleibt unentschieden; Security-Eintrag und Advisory von `v0.13.1` nennen das Gemessene mit Datum und Umgebung und sagen, dass die Betroffenheit vom Basis-Image abhängt, und warten nicht auf die Entscheidung. *Vermerk 2026-10-22, nach dem Release:* `v0.13.1` ist veröffentlicht ([`slice-v1-release-cut-v0.13.1`](../done/slice-v1-release-cut-v0.13.1.md) §7); der Eintrag *Security* im CHANGELOG und das Advisory nennen die Abhängigkeit vom Basis-Image („depends on the base image“) und zwei Läufe (der Verifier maß einen dritten, siehe unten). Die Wahl des Weges A, B oder C steht weiter aus.

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog, „Docker CLI“ als Beispiel), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (im Sandbox-Profil ist der Eintrag abgewiesen; der Bau betrifft das Standardprofil); [`ADR-0032`](../../adr/0032-katalog-features-socket-und-lockerungen.md) Entscheidung 5 (die Aussage, dass `docker-cli` den Socket einbindet, hängt am Docker-Test, der heute mit einer Option von Hand läuft).

**Berührte Spec-Stellen:** [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (Basis-Image, Zeile 1347; Katalog, Zeile 1348), je nach Weg [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) (Drift-Prüfung der Schlüssel). Der Verweis zeigt **aufwärts**: Lastenheft und Spezifikation nennen diesen Slice nie.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-22.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Projekt, das `u-boot` mit aktivem `docker-cli` im Standardprofil erzeugt, baut auf dem erzeugten Dockerfile ohne Handarbeit — oder `u-boot` und das Handbuch sagen vor dem Bau, was dem Nutzer fehlt.

**Befund** (gelesen am 2026-10-22; was nicht gelesen und nicht gemessen ist, steht je Punkt; [`AGENTS.md`](../../../../AGENTS.md) §3.10):

- **Beobachtet, drei Läufe.** Beim ersten Lauf des Docker-Tests `TestDevcontainerUp_Docker_StandardProfileDockerCLIMountsTheSocket` des Vorgänger-Slice (Implementer, 2026-10-22, Docker 29.9.2, `linux/amd64`, lokale Engine) brach der Bau des erzeugten Projekts mit `docker-cli` ab, bevor ein Container entstand: das Installationsskript von `docker-outside-of-docker` endete mit „The 'moby' option is not supported on debian 'trixie'“ (Slice §3 *Umsetzung* 1). Mit `{"moby": false}` im gerenderten Schlüssel lief der Bau (zwei Läufe, Implementer 43,98 s und Verifier 47,46 s). **Wiederholt** (Implementer des Release-Cuts `v0.13.1`, 2026-10-22, Docker 29.9.2, `linux/amd64`, lokale Engine, Debian 13 `trixie` als Betriebssystem des Basis-Images gelesen): `u-boot devcontainer up` auf einem Wegwerf-Projekt im Standardprofil mit `docker-cli` und der Standard-Option brach nach 11 s mit derselben Meldung ab, es entstand kein Container ([`slice-v1-release-cut-v0.13.1`](../done/slice-v1-release-cut-v0.13.1.md) §3 *Umsetzung*). **Ein dritter Lauf** (Verifier des Release-Cuts, 2026-10-22, Linux `amd64`, Docker 29.9.2, lokale Engine; Host-Binary aus `make build-binaries`, Wegwerf-Projekt im Standardprofil mit `docker-cli` und der Standard-Option) endete mit Exit `12` nach 4,6 s mit derselben Meldung („The 'moby' option is not supported on debian 'trixie' because 'moby-cli' and related system packages are not available in that distribution.“), ohne Container; `mcr.microsoft.com/devcontainers/base:debian` war lokal Debian 13 `trixie` ([Verification-Report](../../../reviews/2026-10-22-slice-v1-release-cut-v0.13.1-verify.md) *Einzelnachweise* und V-2). Der Kommentar des Docker-Tests (Zeilen 454 bis 460) nennt den Abbruch „one run on 2026-10-22“; wer den Test für den Weg A oder B ändert, liest die Zahl neu ([`AGENTS.md`](../../../../AGENTS.md) §3.9 b, Verification V-3). Vom Verifier des Vorgänger-Slice ist nur gelesen, dass `mcr.microsoft.com/devcontainers/base:debian` lokal Debian 13 (`trixie`) ist ([Verification-Report](../../../reviews/2026-10-22-slice-v1-sandbox-features-socket-und-lockerungen-verify.md) *Nicht ausgeführt*). Dass die Standard-Option des Features `moby: true` ist, steht als Angabe des Implementers im Vorgänger-Slice (§7, N4) und im Testkommentar (`devcontainer_start_docker_test.go`, Zeilen 454 bis 460); das Installationsskript habe ich nicht gelesen.
- **Beweglicher Tag.** `Dockerfile.tmpl` (`internal/hexagon/application/templates/devcontainer/`, Zeile 3): `FROM mcr.microsoft.com/devcontainers/base:debian`; [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (Zeile 1347) nennt dasselbe Basis-Image ohne Version. Seit wann der Tag auf Debian 13 zeigt, ist nicht gemessen; damit ist offen, wie weit `v0.6.0` bis `v0.13.0` (die betroffenen Versionen des Vorgänger-Slice) mit `docker-cli` bauten (Verification-Report V-3).
- **Keine Feature-Optionen.** `ubootYAMLDevcontainerFeature` (`initproject.go`, Zeilen 129 bis 133) trägt `Enabled`, `Source` und `Version`; `featureCatalogueEntry` (`devcontainer_features.go`, Zeilen 155 bis 179) trägt keine Option; die Vorlage rendert jeden Schlüssel mit `{}` (`devcontainer.json.tmpl`, Zeile 47). Der Docker-Test setzt die Option darum von Hand in den gerenderten Schlüssel (`dcFeatureOptions`, `internal/e2e/devcontainer_start_docker_test.go`, Zeile 429, Aufruf Zeile 465) — er misst eine Variante, nicht das ausgelieferte Projekt ([`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md)).
- **Handbuch, Spezifikation, CHANGELOG sagen nichts dazu** (Verification-Report V-3: Suche nach `moby` und `trixie` in `CHANGELOG.md` leer; Handbuch und Spezifikation nach dem Slice-Plan ohne den Satz, weil er ohne Test über einen beweglichen Tag stünde).

**Wege** (Entscheidung des Projektinhabers vor dem Start; der Slice wählt nicht):

- **A, Feature-Optionen in `u-boot.yaml`** (oder ein Katalog-Default `moby: false` für `docker-cli`): neues Feld im Schema, im Rendern und in der Drift-Prüfung; das Lastenheft nennt nach dem Gelesenen keine Feature-Optionen ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) nicht auf Optionen durchgesehen — prüfen), also vermutlich ein Lastenheft-CR. Der Docker-Test lässt die Handarbeit fallen.
- **B, Pin des Basis-Images** (ein Tag mit Debian 12 oder eine Version): ändert `Dockerfile.tmpl` Zeile 3 und [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) Zeile 1347 und damit jedes erzeugte Dockerfile (Golden Cases, `harness/replay.md`, alle Docker-Tests, die auf dem Basis-Image bauen). **Nicht gelesen:** welche der übrigen Blöcke des Dockerfiles (`nestedRuntime: podman` und `docker`, `ollama`, Egress) an die Debian-Version gebunden sind.
- **C, ein Hinweis in Handbuch und `CHANGELOG.md`:** nach [`AGENTS.md`](../../../../AGENTS.md) §3.10 nennt er nur, was gemessen ist, mit Datum, Umgebung und Zahl der Läufe, und trägt keine Bedingung über einen beweglichen Tag ohne Test (Empfehlung (b) des Reviews, [Review-Report](../../../reviews/2026-10-22-slice-v1-sandbox-features-socket-und-lockerungen.md) F-5: „Known issue“ mit Datum, ohne Bezug auf den Test). Das Problem bleibt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Quellen der Katalog-Features, die Render-Schlüssel externer Quellen (N1 und N2)** — [`slice-v1-katalog-features-quellen-und-freigabe`](slice-v1-katalog-features-quellen-und-freigabe.md), und **die Freigabeprüfung in `generate` (N3)** — [`slice-v1-generate-prueft-freigabeliste`](slice-v1-generate-prueft-freigabeliste.md): andere Ursache und kein Entscheid des Projektinhabers über Optionen oder Basis-Image.
- **Der Mount von `docker-cli`, die Abweisung im Sandbox-Profil und die Ausweisungen** — geliefert in [`slice-v1-sandbox-features-socket-und-lockerungen`](../done/slice-v1-sandbox-features-socket-und-lockerungen.md); sie hängen nicht daran, dass der Bau gelingt (sie sagen, was das Feature einbindet, nicht, dass der Bau gelingt).
- **Eine Prüfung anderer Katalog-Features auf dieselbe Ursache** — nicht gemessen; der Slice misst `docker-cli`. Findet die Messung ein weiteres Feature, das mit dem Basis-Image nicht baut, nennt sie es und der Projektinhaber entscheidet über einen eigenen Vorgang.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Messung** ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)): `docker-cli` mit der Standard-Option und mit `moby: false` baut auf dem aktuellen Tag des Basis-Images und, falls die Registry einen führt, auf einem Tag mit Debian 12; mindestens zwei Läufe je Zelle, datiert, mit Umgebung (Engine, Architektur). *Beleg:* Messprotokoll unter `docs/messungen/<kennung>/` mit Rohdaten, Skript unter `tools/messung/<name>/` ([`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten)).
- [ ] **Der gewählte Weg** (A, B oder C): bei A oder B baut das ausgelieferte Projekt ohne Handarbeit, der Docker-Test misst den Mount am erzeugten Projekt und entfällt die Variante (`dcFeatureOptions`); bei C steht der Hinweis mit Datum und Umgebung in Handbuch und `CHANGELOG.md`. *Beleg:* Tests je Bedingung und Mutation; Spezifikation und Handbuch nach dem Code mit Beleg-Tabelle (Bedingung → Test oder Streichung); Handbuch und READMEs ohne Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6).
- [ ] `make gates` grün; `make docs-check` 0 Befunde; `make doc-trace` ohne Waise; `make doc-commits RANGE="origin/main..HEAD"` 0 Befunde; `make test-docker` (Nachweis ist der CI-Lauf `integration-docker` mit sichtbar bestandenem Test, [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md)).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8). Ändert eine Fix-Runde Vertragstext, liest ein Nachlauf-Review im frischen Kontext genau diesen Diff, bevor der Verifier startet.
- [ ] Doku-Update: Handbuch, `CHANGELOG.md`, bei A oder B auch Spezifikation und `harness/replay.md`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt, mit der Stelle, an der der Planner die Aussage gelesen hat (gelesen am 2026-10-22). Eine Zeile „prüfen“ nennt eine Datei, die der Planner nicht gelesen hat. Sie wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen (`AGENTS.md` §3.10); die Zeilennummern gelten am 2026-10-22 und verschieben sich, der Implementer liest sie am Ist-Stand nach. Welche Zeilen der Weg berührt, entscheidet der Projektinhaber; die Tabelle nennt alle drei.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/messung/<name>/`, `docs/messungen/<kennung>/` | neu | Skript und Protokoll der Messung (Standard-Option und `moby: false`, aktueller Tag und Debian-12-Tag) |
| `internal/e2e/devcontainer_start_docker_test.go` | update (Wege A und B) | `TestDevcontainerUp_Docker_StandardProfileDockerCLIMountsTheSocket` (Zeile 461) setzt die Option mit `dcFeatureOptions` (Zeile 429, Aufruf Zeile 465) von Hand; die Handarbeit und die Begründung im Doc-Kommentar (Zeilen 454 bis 460) entfallen, wenn das erzeugte Projekt baut |
| `internal/hexagon/application/initproject.go`, `devcontainer_features.go`, `templates/devcontainer/devcontainer.json.tmpl`, `config.go` | update (Weg A) | `ubootYAMLDevcontainerFeature` (Zeilen 129 bis 133) bekommt ein Feld für Optionen oder `featureCatalogueEntry` (Zeilen 155 bis 179) einen Default; die Vorlage (Zeile 47) rendert sie statt `{}`; der Pfad `config set` kennt die Option (nicht gelesen — prüfen); die Drift-Prüfung von `doctor` (`projectAllFeatureEntries`, `devcontainer_features.go` Zeile 697) vergleicht nur die Schlüssel (nicht auf Optionen gelesen — prüfen) |
| `internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl` | update (Weg B) | Zeile 3 `FROM mcr.microsoft.com/devcontainers/base:debian`; alle Golden Cases des Dockerfiles ändern sich (nicht gelesen — prüfen: Testdateien mit dem Wortlaut `FROM`) |
| `spec/spezifikation.md` | update (Wege A und B) | [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) Zeile 1347 (Basis-Image) und Zeile 1348 (Katalog); Historie; nach dem Code, mit Beleg-Tabelle |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | update (prüfen) | je nach Weg; ohne Kennungen im Nutzertext; die genauen Stellen sind nicht gelesen |
| `spec/lastenheft.md` | **unverändert**, prüfen | bei Weg A ein CR-Slice (§4), nicht Teil dieses Slice |
| `docs/plan/adr/`, `spec/architecture.md`, `Makefile`, `.golangci.yml`, `.d-check.yml` | **unverändert** | Accepted-ADRs sind nicht änderbar; ein Pin des Basis-Images oder ein neues Schema kann eine Folge-ADR verlangen; Gates und Schwellen bleiben ([`AGENTS.md`](../../../../AGENTS.md) §3.6) |
| diese Datei, `docs/plan/planning/open/README.md` | update | Zeile dieses Slice; Moves per `make slice-mv` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der Projektinhaber hat den Weg (A, B oder C) entschieden und den Slice einem Release zugeordnet; kein anderer Slice in `in-progress/` (WIP-Limit 1, gelesen beim Start). Die Messung ist der erste Schritt des Slice und kann den Weg ändern (Rückführung).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Weg A berührt das Lastenheft (Feature-Optionen als neue Zusage) — dann ist das ein Lastenheft-CR-Slice vor dem Code, mit Wortlaut-Freigabe des Projektinhabers; Weg B ändert jedes erzeugte Dockerfile und braucht eine Folge-ADR.
- `in-progress` → `open` (blockiert): (a) die Messung kann den Abbruch nicht wiederholen (der Tag des Basis-Images oder das Feature hat sich geändert) — dann entfällt der Gegenstand, §7 nennt `Gegenstand:` mit dem Grund und der Messung; (b) der Projektinhaber entscheidet den Weg nicht.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` und `make docs-check` sind grün; bei A oder B bestätigt der Lauf `integration-docker` des CI den Test am erzeugten Projekt ohne Handarbeit (oder die Closure benennt den Nachweis, der an dessen Stelle trat, mit Grund); Review-Report, ein Nachlauf-Review über jeden Diff von Vertragstext und der Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Die Ausgänge stehen bei der Closure; bis dahin nennt jede Zeile die Vorsorge.

1. **Der Abbruch ist nicht wiederholbar oder hängt an etwas anderem als der Option** (ein Lauf, bewegliche Tags, ein Installationsskript im Netz). — **Ausgang:** offen bis zur Closure. Vorsorge: Die Messung steht am Anfang; zwei Läufe je Zelle, datiert, mit Umgebung; keine Aussage über andere Hosts oder Architekturen ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
2. **Der gewählte Weg hat Nebenwirkungen auf andere Blöcke des Dockerfiles** (Weg B: Podman, Docker, Ollama, Egress hängen an der Debian-Version). — **Ausgang:** offen bis zur Closure. Vorsorge: Die Tests der Blöcke laufen gegen das gepinnte Basis-Image; das Ergebnis steht im Messprotokoll.
3. **Eine Bedingung in Handbuch oder Spezifikation hat keinen Test** (Weg C ist anfällig: eine Aussage über einen beweglichen Tag). — **Ausgang:** offen bis zur Closure. Vorsorge: Beleg-Tabelle, jede Aussage mit Datum und Umgebung oder gestrichen ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)).
4. **Der Docker-Test braucht Netz und hängt am Cache des Hosts.** — **Ausgang:** offen bis zur Closure. Vorsorge: Aufräumen wie die Nachbartests, Nachweis ist der CI-Lauf ([`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md), [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md)).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders als geplant · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6 · Drei Paarungen.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden je nach Weg `internal/hexagon/application` (APP), `internal/e2e` (E2E), `tools/` (WZG, Messskript) sowie `spec/`, `harness/`, `docs/plan/` und `docs/messungen/` (DOK); `docs/user/`, die READMEs und `CHANGELOG.md` zählen wie in den übrigen Slices zu DOK. Alle sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) als Greenfield auditiert; eine zu grobe Sub-Area ist nicht dabei.

**Vorgelagert — offene Beobachtungen sichten** (Register gelesen am 2026-10-22, nach der Closure des Vorgänger-Slice; Zähler = Dateien in `evidence/`): [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) (2): der Docker-Test des Vorgänger-Slice ist der zweite Beleg; dieser Slice schließt die Lücke zwischen Messvariante und ausgeliefertem Projekt (Bei A oder B); [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (12, offen über der Schwelle) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (16, verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10): Prüfpflicht für Weg C und für jeden Satz über den Tag des Basis-Images (Risiken 1 und 3); [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (37, verkörpert in §3.10): jede Aussage dieses Plans über Code steht mit der gelesenen Stelle da oder als „prüfen“; [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md) (1), [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (1) und [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (2): Docker-Tests mit Bau treffen alle drei (Risiko 4). Evidenz-/Diskrepanz-Risiko insgesamt: mittel (die Beobachtung ist ein Lauf).

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
