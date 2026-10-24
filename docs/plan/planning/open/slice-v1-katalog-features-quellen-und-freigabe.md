# Slice V1: Quellen der Features — die Katalogquellen von `cpp` und `postgres-client` und der Render-Schlüssel externer Quellen (Nebenbefunde N1 und N2 von `ADR-0032`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Folge-Punkt von [`slice-v1-sandbox-features-socket-und-lockerungen`](../done/slice-v1-sandbox-features-socket-und-lockerungen.md) (Nebenbefunde N1 und N2, [`ADR-0032`](../../adr/0032-katalog-features-socket-und-lockerungen.md) *Nicht Gegenstand* (a) und (b)); nicht Teil von `v0.13.1`, Release-Zuordnung beim Projektinhaber. Es gibt keine Closure-Bedingung, die von der DoD dieses Slice verschieden wäre.

**Teilung (2026-10-24):** Der Punkt N3 (die Prüfung der Freigabeliste in `generate devcontainer`) ist ausgegliedert in [`slice-v1-generate-prueft-freigabeliste`](../done/slice-v1-generate-prueft-freigabeliste.md): Er ist entscheidungsfrei und brauchte keine Messung (die Lücke ist am Binary gezeigt, dort §1). Dieser Slice behält N1 und N2 mit ihrem Start-Trigger (§4): die Messung steht am Anfang, die Entscheidung des Projektinhabers zu N1 folgt der Messung. Die Reihenfolge „N3 vor N2“ bleibt: der Schlüssel aus N2 wird erst committet, wenn der ausgegliederte Slice geliefert ist (§4, §6). **Stand 2026-10-24:** N3 liegt in `done/` ([`slice-v1-generate-prueft-freigabeliste`](../done/slice-v1-generate-prueft-freigabeliste.md)); der Schlüssel aus N2 darf damit committet werden. Der Start-Trigger von N1 und N2 (§4: Priorisierung und Release-Zuordnung durch den Projektinhaber, Messung zuerst) ist unverändert.

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog; ohne explizit erlaubte Quelle ist der Versuch, externe Quellen zu nutzen, ein fachlicher Fehler mit Exit-Code `10`), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (letzter Absatz: Features aus externen Quellen), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes); [`ADR-0032`](../../adr/0032-katalog-features-socket-und-lockerungen.md) (*Nicht Gegenstand* (a) bis (c), *Re-Evaluierungs-Trigger*: der Render-Schlüssel externer Quellen wird berichtigt), [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md).

**Berührte Spec-Stellen:** [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (Katalogtabelle und Form der externen Quelle), [`SPEC-050`](../../../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) (Warnung je Feature aus einer externen Quelle), [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) (Beschreibung von `devcontainer.features.allowlist`, voraussichtlich). Der Verweis zeigt **aufwärts**: Lastenheft und Spezifikation nennen diesen Slice nie.

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

**Ziel:** Jeder Schlüssel, den `u-boot generate devcontainer` unter `features` schreibt, nennt eine Quelle, die es gibt und die die Dev-Container-CLI in dieser Form auflöst, und eine Quelle außerhalb von `devcontainer.featureSources.allow` endet bei `generate` mit Exit-Code `10` — wie bei `config set …source` und bei `doctor`.

**Befund** (gelesen am 2026-10-22 im Code und in den Quellen der ADR; was nicht gelesen und was nicht gemessen ist, steht je Punkt; [`AGENTS.md`](../../../../AGENTS.md) §3.10):

- **N1, Quellen von `cpp` und `postgres-client`.** `featureCatalogue` (`internal/hexagon/application/devcontainer_features.go`) führt `cpp` → `ghcr.io/devcontainers/features/cpp` (Zeilen 232 bis 236) und `postgres-client` → `ghcr.io/devcontainers/features/postgresql-client` (Zeilen 242 bis 246), beide mit der Standardversion `1`; der Doc-Kommentar von `featureCatalogue` (Zeilen 194 bis 202) sagt, `ghcr.io` habe unter diesen beiden Quellen kein Feature gekannt (Lesung vom 2026-10-21) und beide Einträge stünden ungelesen. [`ADR-0032`](../../adr/0032-katalog-features-socket-und-lockerungen.md) *Nicht Gegenstand* (a): `ghcr.io` meldet beide Quellen als unbekannt, und `devcontainers/features/src` führt sie nicht; der Hauptagent nennt als Antwort der Registry `NAME_UNKNOWN` (Mitteilung an den Planner vom 2026-10-22). **Nicht gelesen und nicht gemessen in dieser Rolle:** weder die Registry noch `devcontainers/features/src` (kein Netz); die Aussage „existiert nicht“ stammt aus der ADR und der Mitteilung. Die Einträge stehen in [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (Zeile 1348), in der Katalogtabelle des Handbuchs (Zeile 981 für `postgres-client`) und in den Namenslisten von `README.md` (Zeile 155) und `README.de.md` (Zeile 157); [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) nennt „PostgreSQL Client“ als Beispiel (Zeile 405), `C++` nicht.
- **N2, Render-Schlüssel einer externen Quelle.** `projectFeatureSourceVersion` (`devcontainer_features.go`, Zeile 477) gibt eine nicht leere `source:` wörtlich zurück (Version `1`, wenn keine gesetzt ist); die Vorlage rendert `"{{$f.Source}}:{{$f.Version}}": {}` (`templates/devcontainer/devcontainer.json.tmpl`, Zeile 47); `validateFeatureSource` (Zeile 60) lässt die Schemata `http`, `https` und `oci` zu ([`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), Zeile 1351). Ein Schlüssel hat damit die Form `https://host/pfad:1` oder `oci://host/pfad:1`. Nach [`ADR-0032`](../../adr/0032-katalog-features-socket-und-lockerungen.md) *Nicht Gegenstand* (b) verlangt die Dev-Container-CLI `0.89.0` für einen `https://`-Bezeichner einen Tarball `devcontainer-feature-<id>.tgz` (`getFeatureIdType`, `processFeatureIdentifier` in `src/spec-configuration/containerFeaturesConfiguration.ts`, vom Architect am 2026-10-21 im Quelltext des Tags gelesen). **Nicht gelesen in dieser Rolle:** der Quelltext der CLI; **nicht gemessen:** ob die CLI einen Schlüssel der Form `https://…:1` oder `oci://…:1` auflöst. Der Review des Vorgänger-Slice führte dieselbe Vermutung für `oci://…:1` (Probe P2, [Review-Report](../../../reviews/2026-10-22-slice-v1-sandbox-features-socket-und-lockerungen.md) F-6: „vermutlich nicht“, nicht neu gemessen).
- **N3, `generate` prüft die Freigabeliste nicht — ausgegliedert** in [`slice-v1-generate-prueft-freigabeliste`](../done/slice-v1-generate-prueft-freigabeliste.md) (Befund, Aufrufgraph und Reproduktion am Binary dort §1). Die Aussage dieses Punktes „an keinem Lauf gezeigt“ ist überholt: Am 2026-10-24 lief das Release-Image `0.15.0`, und `generate devcontainer` rendert eine nicht freigegebene Quelle mit Exit-Code `0`.
- **Zusammenhang.** Solange der Schlüssel aus N2 nicht baut, bleibt die fehlende Prüfung ohne Folge für den Bau; wird er berichtigt, schreibt `generate` eine Quelle ohne Freigabe in einen Schlüssel, den die CLI auflöst und dessen Skripte beim Bau laufen ([`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)). Das ist eine Schlussfolgerung aus dem Gelesenen, nicht gemessen. Sie bestimmt die Reihenfolge: [`slice-v1-generate-prueft-freigabeliste`](../done/slice-v1-generate-prueft-freigabeliste.md) vor N2 oder zusammen mit N2, nie danach.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **N4, der Bau von `docker-cli` mit der Standard-Option auf dem Standard-Basis-Image** — [`slice-v1-docker-cli-bau-auf-standard-basis-image`](slice-v1-docker-cli-bau-auf-standard-basis-image.md): andere Ursache (Installationsskript und Debian-Version, nicht die Auflösung der Quelle), und der Weg braucht eine Entscheidung des Projektinhabers (Feature-Optionen, Pin des Basis-Images oder Hinweis).
- **Der Satz zum Volume `minikube-config` von `kubectl-helm` im Handbuch** — Entscheidung 9 des Vorgänger-Slice (Herkunft: Upstream-Metadaten, nicht gemessen, kein Test); ein eigener Vorgang, dem der Projektinhaber noch keinen Plan gegeben hat.
- **Metadaten der Features online lesen, ein Sensor für spätere Minor-Versionen** — braucht Netz und einen neuen Port; `generate` bliebe nicht offline ([`ADR-0032`](../../adr/0032-katalog-features-socket-und-lockerungen.md) *Verglichene Alternativen*).
- **N3, die Prüfung der Freigabeliste in `generate devcontainer`** — [`slice-v1-generate-prueft-freigabeliste`](../done/slice-v1-generate-prueft-freigabeliste.md) (Teilung 2026-10-24): entscheidungsfrei, braucht weder Netz noch Messung der CLI; die Liefer-Reihenfolge vor dem Schlüssel aus N2 steht in §4 und §6.
- **Eine Änderung der Form der Freigabeliste oder des Flags `--allow-external-feature-sources`** — [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) bleibt unverändert; die Slices stellen die Zusage her, die dort steht. Ein Eingriff in das Lastenheft (etwa der Entfall des Beispiels „PostgreSQL Client“ bei einer Entfernung des Eintrags) wäre ein eigener CR-Slice, nicht Teil dieses Slice (§4, Rückführung).
- **Die Warnungen zu `go` und zu externen Features im Sandbox-Profil** — geliefert in [`slice-v1-sandbox-features-socket-und-lockerungen`](../done/slice-v1-sandbox-features-socket-und-lockerungen.md); sie bleiben, auch wenn der Schlüssel berichtigt wird.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Render-Schlüssel externer Quellen** (N2): erst die Messung, dann der Schlüssel. Gemessen wird mit der Dev-Container-CLI aus dem Hilfs-Image (Version aus `internal/adapter/driven/docker/devcontainercli/package.json`), welche Schlüsselformen für `https://`- und `oci://`-Quellen sie auflöst; `u-boot` rendert die Form, die sie auflöst, und weist ab, was keine Form hat (Exit-Code `10`, Abhilfe in der Meldung). *Beleg:* Messprotokoll unter `docs/messungen/<kennung>/` ([`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten)), Golden Case je Form, ein Docker-Test, in dem ein Feature aus einer externen Quelle baut; `devcontainer.features.drift` in `doctor` vergleicht dieselbe Form.
- [ ] **Katalogquellen `cpp` und `postgres-client`** (N1): Nach der Messung (Registry und `devcontainers/features/src`, am Tag der Messung datiert) entscheidet der Projektinhaber über den Weg — Eintrag entfernen, Quelle eines Drittanbieters (mit `publisher` und Ausweisung nach [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)) oder Dockerfile-Block wie `ollama`; der Implementer wählt nicht. *Beleg:* Katalogtabelle und Tabellentest (`TestFeatureCatalogue_*`), [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), Handbuch, beide READMEs und `CHANGELOG.md` sagen dasselbe; bleibt `cpp` oder `postgres-client` im Katalog, baut je ein Docker-Test das Feature, sonst entfällt der Eintrag mit seinen Texten.
- [ ] `make gates` grün; `make docs-check` 0 Befunde; `make doc-trace` ohne Waise; `make doc-commits RANGE="origin/main..HEAD"` 0 Befunde; `make test-docker` (der Docker-Pfad ist berührt; Nachweis ist der CI-Lauf `integration-docker` mit sichtbar bestandenem neuen Test, [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md)).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8). Ändert eine Fix-Runde Vertragstext (Spezifikation), liest ein Nachlauf-Review im frischen Kontext genau diesen Diff, bevor der Verifier startet ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)).
- [ ] Doku-Update: Spezifikation nach dem Code mit Beleg-Tabelle (Bedingung → Test oder Streichung, [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)), Handbuch, READMEs und `CHANGELOG.md` ohne Kennungen im Nutzertext ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6), `harness/replay.md`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt, mit der Stelle, an der der Planner die Aussage gelesen hat (gelesen am 2026-10-22). Eine Zeile „prüfen“ nennt eine Datei, die der Planner nicht gelesen hat. Sie wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen (`AGENTS.md` §3.10); die Zeilennummern gelten am 2026-10-22 und verschieben sich, der Implementer liest sie am Ist-Stand nach.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/application/generate.go` | update (bedingt) | nur, wenn die Messung zu N2 eine Quelle ohne baubare Schlüsselform ergibt: die Abweisung (Exit-Code `10`) steht in `generateDevcontainer` neben der Freigabeprüfung aus [`slice-v1-generate-prueft-freigabeliste`](../done/slice-v1-generate-prueft-freigabeliste.md) (die Freigabeprüfung selbst gehört nicht mehr zu diesem Slice; Stelle und Sentinel `ErrConfigValueInvalid`, beide Klassifikatoren gelesen, stehen dort in §3) |
| `internal/hexagon/application/devcontainer_features.go` | update | Katalogeinträge `cpp` und `postgres-client` (Zeilen 232 bis 246) je nach Entscheidung; `projectFeatureSourceVersion` (Zeile 477) und `projectFeatureEntry` (Zeile 447) für die gemessene Schlüsselform; Doc-Kommentar von `featureCatalogue` (Zeilen 194 bis 202: „stand unread“) wird zur gelesenen oder entfallenen Aussage; `projectAllFeatureEntries` (Zeile 697) trägt dieselbe Form in den Drift-Vergleich |
| `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` | update (bedingt) | Zeile 47 `"{{$f.Source}}:{{$f.Version}}": {}`, falls die Messung eine andere Form verlangt; sonst bleibt die Vorlage |
| `internal/hexagon/application/doctor.go` | update (prüfen) | `classifyFeatureEntries` (Zeile 1333) prüft die Freigabe schon; ob die Meldung und die Drift-Prüfung die gemessene Schlüsselform lesen, ist nicht gelesen |
| `internal/hexagon/application/export_test.go`, `generate_features_test.go`, `devcontainer_features_test.go`, `doctor_features_test.go`, | Fälle zur Freigabe in `generate` entfallen hier (liegen im ausgegliederten Slice) | update / neu | Golden je gemessener Schlüsselform; Tabellentest des Katalogs; die Namen wählt der Implementer |
| `cmd/uboot/main_test.go` | update (bedingt) | nur für die Abweisung einer Quelle ohne baubare Schlüsselform: ein Fall mit den echten Adaptern (Exit `10`, `--json`, nichts angelegt); Muster `TestRun_Generate_SandboxHostSocketFeatureIsRejected` |
| `internal/e2e/devcontainer_start_docker_test.go` oder ein neuer Docker-Test | update / neu | ein Feature aus einer externen Quelle baut (braucht eine erreichbare Quelle; Wahl des Implementers); Hilfsfunktionen `dcGenerateProject` und `dcFeatureOptions` sind vorhanden (gelesen), Aufräumen wie die Nachbartests |
| `tools/messung/<name>/`, `docs/messungen/<kennung>/` | neu | Skript und Messprotokoll zur Auflösung der Schlüsselformen und zur Registry-Abfrage ([`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten)) |
| `spec/spezifikation.md` | update | [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (Katalogtabelle Zeile 1348, Form der Quelle Zeile 1351), [`SPEC-050`](../../../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) (Abschnitt *Features aus externen Quellen*), [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry), Historie; nach dem Code, mit Beleg-Tabelle |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | update (prüfen) | Katalogtabelle (Handbuch Zeile 981 für `postgres-client`), Namenslisten der READMEs (Zeilen 155 und 157), Abschnitt *Externe Feature-Quellen freigeben*, `[Unreleased]`; ohne Kennungen im Nutzertext; die genauen Stellen sind nicht gelesen |
| `spec/lastenheft.md` | **unverändert**, prüfen | die Zusage steht; entfällt `postgres-client` aus dem Katalog, steht „PostgreSQL Client“ weiter als Beispiel in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Zeile 405) — dann ist das ein CR-Slice (§4) |
| `docs/plan/adr/`, `spec/architecture.md`, `Makefile`, `.golangci.yml`, `.d-check.yml` | **unverändert** | Accepted-ADRs sind nicht änderbar; ändert sich der Weg des Katalogs grundlegend, entsteht eine Folge-ADR mit `Supersedes`; Gates und Schwellen bleiben ([`AGENTS.md`](../../../../AGENTS.md) §3.6) |
| diese Datei, `docs/plan/planning/open/README.md` | update | Zeile dieses Slice; Moves per `make slice-mv` |

Die Zeilennummern dieser Tabelle gelten am 2026-10-22; seit der Closure von [`slice-v1-config-set-validiert-feature-version`](../done/slice-v1-config-set-validiert-feature-version.md) verschoben sich `config.go` und `devcontainer_features.go` (Stand 2026-10-24 siehe [`slice-v1-generate-prueft-freigabeliste`](../done/slice-v1-generate-prueft-freigabeliste.md) §3).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der Projektinhaber hat den Slice priorisiert und einem Release zugeordnet; kein anderer Slice in `in-progress/` (WIP-Limit 1, gelesen beim Start); der Release-Cut `v0.13.1` ist geschnitten, damit der Slice das Patch-Release nicht hält. Die Messung zu N1 und N2 ist der erste Schritt des Slice, nicht seine Voraussetzung. **Vor dem Commit des Schlüssels aus N2** liegt [`slice-v1-generate-prueft-freigabeliste`](../done/slice-v1-generate-prueft-freigabeliste.md) in `done/` (beobachtbar: die Datei liegt dort); die Messung selbst darf vorher laufen.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Messung zu N2 zeigt, dass eine neue Form der Quelle nötig ist (etwa ein anderes Schema in `devcontainer.featureSources.allow`) — dann ist das ein Lastenheft-CR und eine Entscheidung des Projektinhabers; ebenso, wenn der Weg zu N1 das Lastenheft berührt. Die Teilung entscheidet der Projektinhaber; der Implementer teilt nicht selbst.
- `in-progress` → `open` (blockiert): (a) die Messung zeigt, dass die Dev-Container-CLI keine externe Quelle in keiner Form auflöst — dann ist Entscheidung 3 von [`ADR-0032`](../../adr/0032-katalog-features-socket-und-lockerungen.md) ohne Gegenstand, und eine Folge-ADR kommt vor weiterem Code; (b) der Projektinhaber entscheidet den Weg zu N1 nicht.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` und `make docs-check` sind grün; die Docker-Tests laufen im Lauf `integration-docker` des CI sichtbar bestanden (oder die Closure benennt den Nachweis, der an dessen Stelle trat, mit Grund); jede Bedingung der Beleg-Tabelle (§3 nach dem Code) ist durch einen Test belegt oder gestrichen; Review-Report, ein Nachlauf-Review über jeden Diff von Vertragstext und der Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Die Ausgänge stehen bei der Closure; bis dahin nennt jede Zeile die Vorsorge.

1. **Die Messung widerlegt die Annahmen aus der ADR** (die Registry kennt die Quellen von `cpp` und `postgres-client`; die CLI löst einen Schlüssel der Form `https://…:1` oder `oci://…:1` auf). — **Ausgang:** offen bis zur Closure. Vorsorge: Die Annahmen stammen aus einer Lesung des Architect und einer Mitteilung (§1); die Messung steht am Anfang des Slice, ein anderes Ergebnis führt zu §4 (Rückführung), nicht zu einem Umbau im Slice.
2. **N3 wird nach N2 geliefert** und `generate` schreibt für einen Moment eine Quelle ohne Freigabe in einen baubaren Schlüssel. — **Ausgang:** offen bis zur Closure. Vorsorge: N3 ist ausgegliedert in [`slice-v1-generate-prueft-freigabeliste`](../done/slice-v1-generate-prueft-freigabeliste.md); der Implementer committet den Schlüssel erst, wenn dieser Slice in `done/` liegt (§4, Start).
3. **Eine Bedingung in Spezifikation oder Handbuch hat keinen Test** (Form der Quelle, Verhalten bei einer nicht freigegebenen Quelle, Katalogtabelle). — **Ausgang:** offen bis zur Closure. Vorsorge: Beleg-Tabelle nach dem Code, jeder Teilsatz einer Meldung mit eigenem Pin ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), im Vorgänger-Slice viermal in einem Slice).
4. **Eine Zusage reicht weiter als die Messung** (eine Schlüsselform „funktioniert“, gemessen mit einer Quelle und einer CLI-Version). — **Ausgang:** offen bis zur Closure. Vorsorge: Messsatz mit Umgebung, Datum und Zahl der Läufe; keine Aussage über andere Quellen ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
5. **Der Docker-Test braucht Netz** (eine erreichbare externe Quelle, die Registry) und hängt am Cache des Hosts. — **Ausgang:** offen bis zur Closure. Vorsorge: Aufräumen wie die Nachbartests, Nachweis ist der CI-Lauf ([`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md), [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md)).

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/application` (APP), voraussichtlich `internal/adapter/driving/cli` (CLI, Klassifikatoren), `internal/e2e` (E2E), `tools/` (WZG, Messskript) sowie `spec/`, `harness/`, `docs/plan/` und `docs/messungen/` (DOK); `docs/user/`, die READMEs und `CHANGELOG.md` zählen wie in den übrigen Slices zu DOK. Alle sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) als Greenfield auditiert; eine zu grobe Sub-Area ist nicht dabei.

**Vorgelagert — offene Beobachtungen sichten** (Register gelesen am 2026-10-22, nach der Closure des Vorgänger-Slice; Zähler = Dateien in `evidence/`): [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (16, verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10) und [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (37, verkörpert in §3.10): Prüfpflicht für jede Aussage dieses Plans über Code und jede Bedingung im Vertragstext (Risiken 3); [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (12, offen über der Schwelle) und [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) (2): die Messung zu N2 läuft mit der Form, die `u-boot` ausliefert, nicht mit einer Hilfsvariante (Risiko 4); [`BEO-DOK/abhilfe-ohne-wirkungsprobe`](../observations/BEO-DOK/abhilfe-ohne-wirkungsprobe/observation.md) (1): gilt für jede Abhilfe, die eine Meldung dieses Slice nennt (die Freigabeprüfung selbst liegt in [`slice-v1-generate-prueft-freigabeliste`](../done/slice-v1-generate-prueft-freigabeliste.md)); [`BEO-APP/pruefung-nur-beim-lesen`](../observations/BEO-APP/pruefung-nur-beim-lesen/observation.md) (1) betrifft die Teilung nur dort; [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md) (1), [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (1) und [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (2): ein neuer Docker-Test trifft alle drei (Risiko 5). Evidenz-/Diskrepanz-Risiko insgesamt: mittel (zwei Annahmen der Planung sind nicht gemessen).

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
