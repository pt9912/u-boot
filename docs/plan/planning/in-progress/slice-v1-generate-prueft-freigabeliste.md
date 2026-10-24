# Slice V1: `generate devcontainer` prüft die Freigabeliste der Feature-Quellen (Nebenbefund N3 von `ADR-0032`, ausgegliedert aus `slice-v1-katalog-features-quellen-und-freigabe`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Teilung des Slice [`slice-v1-katalog-features-quellen-und-freigabe`](../open/slice-v1-katalog-features-quellen-und-freigabe.md): dessen Punkt N3 ist entscheidungsfrei (die Zusage steht im Lastenheft, der Weg ist im Code vorgezeichnet), N1 und N2 brauchen eine Messung und eine Entscheidung des Projektinhabers. Release-Zuordnung: offen, beim Projektinhaber; kein Release ist festgelegt. Es gibt keine Closure-Bedingung, die von der DoD dieses Slice verschieden wäre.

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Aufzählung zu optionalen externen Feature-Quellen, Zeilen 415 bis 422: „Standardmäßig sind nur lokal hinterlegte oder ausdrücklich freigegebene Features erlaubt“, `--allow-external-feature-sources` gültig auch für `generate devcontainer`, „Ohne explizit erlaubte Quelle […] fachlicher Fehler (Code [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Exit-Code `10`)“), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Zeilen 1076 bis 1080: kein externer Code aus nicht freigegebenen Quellen, Devcontainer-Features „nicht ausdrücklich in `u-boot.yaml` freigegeben“), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes); [`ADR-0032`](../../adr/0032-katalog-features-socket-und-lockerungen.md) (*Nicht Gegenstand* (c)). [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) trägt die Zusage nicht: der Abschnitt handelt von Zugangsdaten und nennt keine Freigabeliste (Zeilen 484 bis 500 auf „Quelle“ und „extern“ gelesen: Treffer sind die Quelle der Zugangsdaten, Zeile 490, und die gepinnten Quellen des Hilfs-Images, Zeile 496); er bleibt unberührt.

**Berührte Spec-Stellen:** [`SPEC-050`](../../../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) (Abschnitt *Features aus externen Quellen*, Zeile 638: „Der Eintrag wird wie bisher erzeugt“ gilt weiter für freigegebene Quellen und wird um die Abweisung ergänzt), [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) (Zeile 1368, Beschreibung von `devcontainer.features.allowlist`, voraussichtlich unverändert), und der Aufruf-Abschnitt von `generate devcontainer` in der Spezifikation, soweit er das Flag `--allow-external-feature-sources` beschreibt (Zeile nicht gelesen, der Implementer sucht ihn). Der Verweis zeigt **aufwärts**: Lastenheft und Spezifikation nennen diesen Slice nie.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-24.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `u-boot generate devcontainer` mit einem aktivierten Eintrag, dessen `source:` nicht in `devcontainer.featureSources.allow` steht (und den der Aufruf nicht selbst über `--allow-external-feature-sources` freigibt), endet mit Exit-Code `10` und Code [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), ohne etwas zu schreiben — wie `config set …source` und `doctor`, über dieselbe Prüffunktion.

**Befund** (gelesen am 2026-10-24; [`AGENTS.md`](../../../../AGENTS.md) §3.10):

- **Aufrufgraph.** `featureSourceInAllow` (`internal/hexagon/application/config.go`, Zeile 633; Vergleich byte-gleich) hat im Produktcode zwei Aufrufer: `revalidateFeatureEntry` (`config.go`, Zeile 557, Prüfung Zeile 575, Sentinel `driving.ErrConfigValueInvalid`) und `classifyFeatureEntries` (`doctor.go`, Zeile 1333, Prüfung Zeile 1341; Fehler-Eintrag von `devcontainer.features.allowlist`, [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry), Zeile 1368). `generateDevcontainer` (`generate.go`, Zeile 764) liest über `readProjectConfig` (Zeile 268; `validateDevcontainer` prüft nur das Format, `validateDevcontainerFeatures` in `devcontainer_features.go`, Zeile 613, kommentiert „`Source set but not in Allow` is also deferred (T4 enforcement)“), ruft `collectDevcontainerFeatures` (Zeile 789; Doc-Kommentar in `devcontainer_features.go`: „T4 enforces allowlist membership — T3 trusts the value here“) und rendert `entry.Source` wörtlich (`projectFeatureSourceVersion`). Ein T4-Durchsetzer in `generate` existiert nicht.
- **Reproduziert am Binary, nicht nur gelesen.** Am 2026-10-24 lief das Release-Image `ghcr.io/pt9912/u-boot:0.15.0` (Linux `amd64`, Aufruf als Host-Benutzer, Verzeichnis mit `.git/config` und `u-boot.yaml`: `devcontainer.enabled: true`, Eintrag `evil` mit `enabled: true` und `source: https://evil.example.org/feat`, **keine** `featureSources.allow`): `generate devcontainer` endete mit **Exit-Code `0`** und schrieb `"https://evil.example.org/feat:1": {}` nach `.devcontainer/devcontainer.json`; im selben Verzeichnis meldete `doctor` die Prüfung `devcontainer.features.allowlist` als Fehler („not in devcontainer.featureSources.allow“), und `config set devcontainer.features.evil.source https://other.example.org/x` endete mit Exit-Code `10`. Die Lücke ist damit an einem Lauf gezeigt (ein Lauf, ein Image, ein Host). Der Quelltext von `main` weicht seit `v0.15.0` in `config.go` und `devcontainer_features.go` nur um die Regel zur `version` ab (`git diff v0.15.0 HEAD`, gelesen); `generate.go` und `doctor.go` sind unverändert. **Nicht gemessen:** das Verhalten von `main` am Binary, `--json`, `--dry-run`, ob die Dev-Container-CLI den Schlüssel baut (das ist N2 des Ursprungs-Slice, nicht dieser Slice).
- **Fehlerklasse.** `driving.ErrConfigValueInvalid` steht in beiden Klassifikatoren: `isConfigValidationError` (`internal/adapter/driving/cli/cli.go`, Zeile 455) ergibt Exit-Code `10`, `mapGenerateErrorToDiagnostic` (`internal/adapter/driving/cli/generate.go`, Zeile 293) den Code [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features); `validateAllowExternalFeatureSourcesEntries` (`generate.go` der Application, Zeile 1180) nutzt denselben Sentinel für dieselbe Zusage. `ErrInvalidFeatureSource` bleibt unbenutzt: `cli.go` Zeile 475 gibt ihm Exit-Code `10`, aber `mapGenerateErrorToDiagnostic` kennt ihn nicht und fiele auf [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (gelesen) — die Dual-Classifier-Regel wäre verletzt. Es entsteht **kein neuer Sentinel und kein neuer Vertrag**.
- **Zusammenhang mit N2.** Ob die Dev-Container-CLI den Schlüssel auflöst, ist offen (Ursprungs-Slice, N2). Dieser Slice braucht die Antwort nicht: Die Zusage von [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) gilt für den Versuch, eine externe Quelle zu nutzen, nicht erst für den gelungenen Bau. Wird N2 später berichtigt, liegt die Prüfung dann schon vor ([`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **N1 (Katalogquellen von `cpp` und `postgres-client`) und N2 (Render-Schlüssel externer Quellen, Messung der Dev-Container-CLI)** — bleiben in [`slice-v1-katalog-features-quellen-und-freigabe`](../open/slice-v1-katalog-features-quellen-und-freigabe.md): sie brauchen Netz, eine Messung und eine Entscheidung des Projektinhabers.
- **Eine Änderung der Form der Freigabeliste, des Vergleichs (byte-gleich) oder des Flags** — [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) bleibt unverändert; der Slice stellt die Zusage her, die dort steht. Eine Normalisierung (abschließender Schrägstrich, Groß-/Kleinschreibung des Hosts) wäre ein eigener Vorgang (Kommentar zu `featureSourceInAllow`, Zeilen 610 bis 632).
- **Die Prüfung in `readProjectConfig` (Lesepfad für alle Befehle)** — würde `doctor` und `config set` als Abhilfe sperren und kennt das Flag des Aufrufs nicht; die Prüfung steht in `generateDevcontainer`, die Prädikatsfunktion ist die gemeinsame.
- **Deaktivierte Einträge (`enabled: false`) und Einträge ohne `enabled`** — `generate` rendert sie nicht (`collectDevcontainerFeatures`), also gibt es nichts zu verweigern; `doctor` meldet sie weiter (`classifyFeatureEntries` prüft jeden Eintrag mit `enabled`-Schlüssel). Das Abschalten bleibt die Abhilfe.
- **Die Warnungen zu externen Features im Sandbox-Profil** ([`SPEC-050`](../../../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen)) — bleiben für freigegebene Quellen.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Freigabeprüfung in `generate devcontainer`** ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)): ein aktivierter Eintrag mit `source:` außerhalb der wirksamen Freigabe endet mit Exit-Code `10` und Code [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) im Fehlerpfad (`ExitCode`) **und** im `--json`-Umschlag (`mapGenerateErrorToDiagnostic`), ohne etwas zu schreiben (keine Datei, kein Verzeichnis, `u-boot.yaml` unverändert), auch unter `--dry-run`. Die **wirksame Freigabe** ist `devcontainer.featureSources.allow` plus die Quellen von `--allow-external-feature-sources` desselben Aufrufs (in der Form von `normaliseFeatureSources`); die Datei wird dabei nicht vorab geändert, der Anhang nach erfolgreichem Lauf (`applyAllowExternalFeatureSources`) bleibt wie er ist. Die Meldung nennt Feature und Quelle und die Abhilfe in dieser Reihenfolge: das Freigeben (`u-boot config set devcontainer.featureSources.allow <quelle>` oder `--allow-external-feature-sources <quelle>`), dann das Abschalten (`u-boot config set devcontainer.features.<feature>.enabled false`), und den Satz zum byte-gleichen Vergleich wie bei `config set`; der Implementer **führt beide Abhilfen aus**, bevor sie in den Text gehen ([`BEO-DOK/abhilfe-ohne-wirkungsprobe`](../observations/BEO-DOK/abhilfe-ohne-wirkungsprobe/observation.md)). Eine **einzige** Prädikatsfunktion (`featureSourceInAllow`, oder ihr Nachfolger mit derselben Rolle) trägt `revalidateFeatureEntry`, `classifyFeatureEntries` und die neue Prüfung. *Beleg:* Tests und Mutationen in §3; Beleg-Tabelle (Bedingung → Test) für jeden neuen Satz in Spezifikation und Handbuch.
- [x] **Vertragstext nach dem Code:** [`SPEC-050`](../../../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) und der Aufruf-Abschnitt von `generate devcontainer` sagen die Abweisung (unabhängig vom Profil: die Freigabepflicht gilt in beiden Profilen), Handbuch (Abschnitt *Externe Feature-Quellen freigeben*, Zeilen 1064 bis 1122: der Hinweis „Ohne Freigabe endet `config set …source` mit Exit-Code 10“ nennt `generate devcontainer` mit, die Tabelle *Prüfungen durch `u-boot doctor`* bleibt) und `CHANGELOG.md` (`[Unreleased]`, Eintrag *Fixed* oder *Security*, englisch, ohne Kennungen im Nutzertext der Handbuch-Zeilen, [`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6) folgen dem Code; [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) bleiben unverändert (Lastenheft nicht berührt). Muss der Wortlaut des Lastenhefts doch geändert werden, ist das eine Entscheidung des Projektinhabers: der Implementer stoppt, der Slice geht nach `open/` zurück (§4).
- [x] `make gates` grün; `make docs-check` 0 Befunde; `make doc-trace` ohne Waise; `make doc-commits RANGE="origin/main..HEAD"` 0 Befunde. Der Docker-Pfad ist nicht berührt (`make test-docker` entfällt; die Prüfung wirkt vor jedem Docker-Aufruf).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8). Ändert eine Fix-Runde Vertragstext (Spezifikation), liest ein Nachlauf-Review im frischen Kontext genau diesen Diff, bevor der Verifier startet ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt, mit der Stelle, an der der Planner die Aussage gelesen hat (gelesen am 2026-10-24). Eine Zeile „prüfen“ nennt eine Datei, die der Planner nicht gelesen hat. Sie wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen (`AGENTS.md` §3.10); Zeilennummern verschieben sich, der Implementer liest sie am Ist-Stand nach.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/application/generate.go` | update | `generateDevcontainer` (Zeile 764): die Prüfung steht **nach** `collectDevcontainerFeatures` (Zeile 789) und **vor** der Sandbox-Abweisung (`sandboxHostSocketFeatureMessage`, Zeile 795), `resolveSandboxClone` und `planDevcontainerFiles`; Reihenfolge festgelegt: Freigabe vor Profil, damit bei zwei Verstößen der Code [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) gilt (festgehalten durch einen Test). Die wirksame Freigabe bildet `normaliseFeatureSources(existing + req.AllowExternalFeatureSources)` ohne Schreiben (die Flag-Liste ist durch `validateAllowExternalFeatureSourcesEntries`, Zeile 1180, schon geprüft). Fehler: `fmt.Errorf("%w: …", driving.ErrConfigValueInvalid, …)` wie dort. Die Kommentare an Zeile 773 („T3's renderer does not enforce the allowlist“) und an `collectDevcontainerFeatures` werden zur gelesenen Aussage (AGENTS.md §3.7) |
| `internal/hexagon/application/devcontainer_features.go` | update | neuer Helfer, der die aktivierten Einträge mit `source:` außerhalb der wirksamen Freigabe liefert (sortiert nach Namen, deterministische Meldung), gebaut auf `featureSourceInAllow`; Doc-Kommentar von `collectDevcontainerFeatures` („T4 enforces allowlist membership“) nennt den Ort der Durchsetzung; der Kommentar in `validateDevcontainerFeatures` („also deferred (T4 enforcement)“) ebenso. Muster der einen Regel an einer Stelle: `featureVersionNeedsSource` (Zeile 653) |
| `internal/hexagon/application/config.go` | update (Kommentar, bedingt) | `featureSourceInAllow` (Zeile 633): bleibt das eine Prädikat; Doc-Kommentar nennt den dritten Aufrufer. Falls der Helfer ein `ubootYAMLConfig` mit gemergter Liste braucht, entsteht keine zweite Vergleichslogik |
| `internal/hexagon/application/doctor.go` | **unverändert**, gelesen | `classifyFeatureEntries` (Zeile 1333) ruft `featureSourceInAllow` (Zeile 1341) und bleibt; der Pin-Test (unten) hält die Gleichheit der drei Aufrufer fest |
| `internal/adapter/driving/cli/cli.go`, `internal/adapter/driving/cli/generate.go` | **unverändert**, gelesen | beide Klassifikatoren tragen `ErrConfigValueInvalid` (`cli.go` Zeile 455; `mapGenerateErrorToDiagnostic`, Zeile 293); Dual-Classifier-Regel ist erfüllt, ohne Eingriff — der Beleg ist der Test am Binary-Pfad (`main_test.go`), nicht die Wortsuche |
| neue Testdatei `internal/hexagon/application/generate_freigabe_test.go`, `export_test.go` (bedingt) | neu / update | Unit-Fälle mit `newGenerateService`/`seedUBootYAMLWithFeatures` (Muster `generate_features_test.go`, Zeilen 42 und 237) und Fake-Dateisystem; die Namen wählt der Implementer |
| `internal/hexagon/application/generate_features_test.go`, `generate_host_socket_test.go`, `generate_ollama_test.go`, `generate_forge_cli_test.go`, `devcontainer_sandbox_relaxation_test.go`, weitere mit `source:` | update (Fixtures) | Tests, die einen aktivierten Eintrag mit `source:` rendern, brauchen eine passende Freigabe. **Bekannt rot nach dem Code:** `TestGenerateDevcontainer_FeaturesHomebrew_ExternalOverrideNoWarning` (`generate_features_test.go`, Zeile 689): die Liste führt `https://ghcr.io/stu-bell/devcontainer-features`, die Quelle lautet `…/devcontainer-features/homebrew` — nicht byte-gleich. Welche weiteren Tests rot werden, ist **nicht gelesen**: der Implementer führt `make test` nach dem Code aus und trägt jede berichtigte Fixture hier ein; er berichtigt die Fixture (Freigabe ergänzen), nie den Vergleich |
| `cmd/uboot/main_test.go` | update | ein Fall mit den echten Adaptern nach dem Muster `TestRun_Generate_SandboxHostSocketFeatureIsRejected` (Zeile 384, Hilfsfunktion `writeFeatureProject`, Zeile 316): Exit-Code `10` für `generate devcontainer` und `--dry-run`, `--json` mit genau einem Eintrag Level `error`, Code [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features); `.devcontainer` nicht angelegt, `u-boot.yaml` byte-gleich; Abhilfe 1 (`config set …featureSources.allow`) → Exit-Code `0`; Abhilfe 2 (`…enabled false`) → Exit-Code `0`; Flag `--allow-external-feature-sources` mit der Quelle → Exit-Code `0` und die Quelle danach in der Liste. Bestehende Fälle mit `source:` im Aufruf (`TestRun_Generate_SandboxFeatureWarnings`, Zeile 337, und weitere) auf Freigabe prüfen |
| `spec/spezifikation.md` | update | [`SPEC-050`](../../../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) Zeile 638 („Der Eintrag wird wie bisher erzeugt“) und der Aufruf-Abschnitt von `generate devcontainer`; nach dem Code, mit Beleg-Tabelle (jeder Teilsatz → Test oder Streichung); Historie-Zeile |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | update | Handbuch Zeilen 1064 bis 1122 (Hinweise zu Exit-Code `10`, Flag), `[Unreleased]` (Zeile 10 ff.); ohne Kennungen im Handbuchtext. `README.md`, `README.de.md` und `harness/replay.md`: **nicht gelesen**; der Implementer prüft per gelesener Aussage, ob `generate devcontainer` dort als „rendert jede Quelle“ beschrieben ist, und trägt das Ergebnis ein |
| `spec/lastenheft.md`, `docs/plan/adr/`, `spec/architecture.md`, `Makefile`, `.golangci.yml`, `.d-check.yml` | **unverändert** | die Zusage steht (Zeilen 415 bis 422); Accepted-ADRs sind nicht änderbar; Gates und Schwellen bleiben ([`AGENTS.md`](../../../../AGENTS.md) §3.6); `architecture.md` nennt die Prüfung nicht (nicht gelesen — der Implementer liest den Abschnitt zu `generate` und trägt die Fundstelle ein) |
| diese Datei, `slice-v1-katalog-features-quellen-und-freigabe.md`, `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` | update | Zeile dieses Slice; Moves per `make slice-mv` |

**Pinnende Tests mit Mutationen** (jede Mutation wird gefahren und färbt genau die genannten Tests rot; die Zeile im Verification-Report nennt Mutation, Datei und roten Test):

| # | Pin | Mutation |
|---|---|---|
| T1 | aktivierter Eintrag, `source` nicht in der Liste: Fehler mit `driving.ErrConfigValueInvalid`, kein `WriteFile`, `u-boot.yaml` byte-gleich, kein `.devcontainer`; je einmal mit und ohne `--dry-run` | die Prüfung entfällt; sie läuft erst nach `writeDevcontainerPhase` |
| T2 | Quelle steht byte-gleich in der Liste: Lauf gelingt, Schlüssel gerendert | die Prüfung weist jede externe Quelle ab |
| T3 | Quelle nur durch das Flag desselben Aufrufs frei: gelingt; danach steht sie in der Liste (bestehender Pfad); ohne das Flag im nächsten Aufruf zählt die Liste | die Flag-Liste fließt nicht in die wirksame Freigabe ein |
| T4 | `enabled: false` und fehlendes `enabled` mit nicht freigegebener Quelle: Lauf gelingt, Schlüssel fehlt (Abhilfe „abschalten“) | die Prüfung bezieht deaktivierte Einträge ein |
| T5 | Vergleichstabelle über **alle drei Aufrufer** mit denselben Fällen (gleiche Quelle, abschließender Schrägstrich, Groß-/Kleinschreibung des Hosts, Leerliste, fehlender Block): `generateDevcontainer`, `ConfigService` (`config set …source`) und `doctor` urteilen gleich — das Gegenstück zur Klasse [`BEO-APP/pruefung-nur-beim-lesen`](../observations/BEO-APP/pruefung-nur-beim-lesen/observation.md) (Lese- und Schreibpfad dieselbe Frage; hier umgekehrt: der Rendern-Pfad hatte die Regel nie) | in `generate` ein eigener Vergleich (z. B. `strings.EqualFold` oder `TrimRight "/"`); der Test färbt sich |
| T6 | zwei Verstöße zugleich (Freigabe und `docker-cli` im Sandbox-Profil): Code [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | Reihenfolge vertauscht |
| T7 | `main_test.go` (echte Adapter): Exit-Code `10`, `--json` Code [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) Level `error`, Abhilfen führen zu Exit-Code `0`, Dateien wie oben | `mapGenerateErrorToDiagnostic` ohne den Zweig für `ErrConfigValueInvalid` (Code fällt auf [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)); `ExitCode` ohne `ErrConfigValueInvalid` (Exit-Code `1`) |
| T8 | Meldung: jeder Teilsatz (Feature, Quelle, Abhilfe 1, Abhilfe 2, byte-gleich) mit eigenem Pin | je Teilsatz gestrichen |

**Ansatz in Sätzen.**

- Die Prüfung ist eine reine Funktion über `cfg`, die Flag-Liste und die aktivierten Einträge; sie liest nichts neu und hängt an keinem Port.
- Der Implementer committet Code und Vertragstext getrennt; Vertragstext erst nach dem grünen Code ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)).
- Die Meldung wird vor dem Text am Binary gefahren (beide Abhilfen), mit dem Image aus `make`, nicht mit einem veröffentlichten Image.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der Projektinhaber hat den Slice priorisiert und einem Release zugeordnet (offen); kein anderer Slice in `in-progress/` (WIP-Limit 1, beim Start zu lesen). Keine Messung und keine Entscheidung des Projektinhabers ist Voraussetzung.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Zahl der berichtigten Fixtures (§3) übersteigt, was ein Review in einer Runde liest (Richtwert: mehr als zehn Testdateien); dann teilt der Planner nach Testdatei, der Implementer teilt nicht selbst.
- `in-progress` → `open` (blockiert): der Wortlaut von [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) oder [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) trägt die Abweisung doch nicht (etwa weil ein Pfad „lokal hinterlegte“ Features über `source:` meint) — **Entscheidung des Projektinhabers**, Lastenheft-CR als eigener Slice; der Slice startet dann nicht und schreibt keinen Code.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` und `make docs-check` sind grün; jede Mutation T1 bis T8 ist gefahren und im Verification-Report belegt; jede Bedingung der Beleg-Tabelle ist durch einen Test belegt oder gestrichen; Review-Report, Nachlauf-Review über jeden Diff von Vertragstext und Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben. Die Wirkung am Binary (Exit-Code `10` mit dem Image aus `make`, im Ursprungsfall dieses Plans) steht im Verification-Report.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Die Ausgänge stehen bei der Closure; bis dahin nennt jede Zeile die Vorsorge.

1. **Bestehende Projekte brechen:** wer eine externe Quelle ohne Freigabe in `u-boot.yaml` trägt und `generate devcontainer` laufen lässt, bekommt ab dem Release Exit-Code `10` statt eines gerenderten Schlüssels (Verhaltensänderung, aber die Zusage des Lastenhefts). — **Ausgang:** offen bis zur Closure. Vorsorge: die Meldung nennt beide Abhilfen; `CHANGELOG.md` nennt die Änderung unter *Fixed* oder *Security* (Wahl des Projektinhabers beim Release-Schnitt, nicht des Implementers) mit der Abhilfe; `doctor` meldet den Zustand schon heute.
2. **Zweite Vergleichslogik schleicht sich ein** (die neue Prüfung normalisiert „gutmütig“ und weicht von `config set` und `doctor` ab). — **Ausgang:** offen bis zur Closure. Vorsorge: ein Prädikat, Pin T5 über drei Aufrufer, Mutation in T5.
3. **Fixtures werden gelockert statt berichtigt** (Tests ohne Freigabe laufen weiter, weil die Prüfung für sie abgeschaltet wird). — **Ausgang:** offen bis zur Closure. Vorsorge: §3 verbietet es; der Reviewer prüft den Diff der Testdateien auf eine Option „ohne Prüfung“ (es gibt keine).
4. **Eine Aussage über Code im Plan trägt nicht** ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md)): die nicht gelesenen Dateien (README, `architecture.md`, `harness/replay.md`, Aufruf-Abschnitt der Spezifikation) sind in §3 als solche geführt. — **Ausgang:** offen bis zur Closure. Vorsorge: Abgleich der Tabelle gegen `git diff --stat` nach der letzten Fix-Runde.
5. **Die Messung am Binary gilt für `v0.15.0`, nicht für `main`.** — **Ausgang:** offen bis zur Closure. Vorsorge: der Implementer wiederholt den Ursprungsfall (Eintrag `evil`) vor dem Code am Image aus `make` und danach, beide Läufe im Verification-Report.

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

**Umsetzungsstand (Implementer, 2026-10-24; vor Review).**

*Ist-Stand der Plan-Tabelle (§3) gegen `git diff --stat` von `750139e5..HEAD`:* berührt sind `internal/hexagon/application/generate.go` (Prüfung nach `collectDevcontainerFeatures`, vor `sandboxHostSocketFeatureMessage`; Kommentar zur Allowlist), `devcontainer_features.go` (`unlistedFeatureSources`, `featureSourceViolationMessage`, zwei Kommentare), `config.go` (nur Doc-Kommentar von `featureSourceInAllow`), die neue `generate_freigabe_test.go`, `cmd/uboot/main_test.go` (neuer Fall plus eine Fixture), `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, `harness/replay.md` (neue Zeile; die Tabelle war als „nicht gelesen“ geführt, gelesen: sie führt Generator-Fälle) und `roadmap.md` (Zeiger unter *Offene Wellen*). `doctor.go`, `cli.go`, `cli/generate.go` unverändert (gelesen: `cli.go` Zeile 455 und `mapGenerateErrorToDiagnostic` Zeile 293 tragen `ErrConfigValueInvalid`; Beleg ist der Test am Binary-Pfad). `README.md`, `README.de.md` (Zeile 168 bzw. 171, nur „Freigabe nötig“) und `spec/architecture.md` (Wortsuche nach Freigabe/allowlist/featureSources: keine Aussage zu `generate`) sind unverändert; `export_test.go` nicht berührt.

*Berichtigte Fixtures (Liste um die tatsächliche Quelle, kein Vergleich gelockert):* `TestGenerateDevcontainer_FeaturesHomebrew_ExternalOverrideNoWarning` und `TestGenerateDevcontainer_ClaudeCodeNodeNotice` (Fall „own source“) in `generate_features_test.go`, `TestGenerateDevcontainer_SandboxHostSocketFeature_OnlyActiveEntryIsRejected` (Fall „source override“) in `generate_host_socket_test.go`, `TestRun_Generate_SandboxFeatureWarnings` in `cmd/uboot/main_test.go`. Die übrigen Tests mit `source:` (`generate_ollama_test.go`, `generate_forge_cli_test.go`, `devcontainer_sandbox_relaxation_test.go`, `generate_features_test.go` sonst) blieben grün, ihre Listen standen schon.

*Beleg-Tabelle (neuer Satz in Spezifikation/Handbuch → Test):*

| Bedingung | Test |
|---|---|
| Quelle außerhalb der Freigabe: Exit-Code `10`, nichts geschrieben, auch `--dry-run` | `TestGenerateDevcontainer_FreigabeRejectsUnlistedSource`, `TestRun_Generate_UnlistedFeatureSourceIsRejected` |
| Code [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) im `--json`-Umschlag | `TestRun_Generate_UnlistedFeatureSourceIsRejected` |
| wirksame Freigabe = Liste plus Flag desselben Aufrufs; Flag-Quelle danach in der Liste, unter `--dry-run` nicht | `…_FreigabeByFlagOfSameCall`, `…_FreigabeByFlagDryRun`, `…_FreigabeRejectsAlsoWithUnrelatedFlag`, `TestRun_Generate_UnlistedFeatureSourceIsRejected` |
| Vergleich wie `config set` und `doctor` (Schrägstrich, Host-Schreibung) | `TestFeatureSourceAllow_ThreeCallersAgree` |
| Prüfung vor der Abweisung von `docker-cli`; beide Profile | `…_FreigabeBeforeSandboxRejection` (Sandbox), `…_FreigabeRejectsUnlistedSource` (Standardprofil) |
| Meldung: Feature, Quelle, Freigeben, Abschalten (Reihenfolge), Vergleichssatz | `…_FreigabeMessageParts` |
| `enabled: false` / ohne `enabled` nicht geprüft; Katalogschlüssel ohne `source:` frei | `…_FreigabeIgnoresInactiveEntries`, `…_FreigabeCatalogueKeyNeedsNoList` |
| beide Abhilfen führen zu Exit-Code `0` | `TestRun_Generate_UnlistedFeatureSourceIsRejected` |
| Handbuchsatz „doctor meldet sie weiter“ | bestehend: `classifyFeatureEntries` prüft jeden Eintrag mit `enabled`-Schlüssel (`TestDoctor_FeaturesAllowlist_*`, nicht neu gepinnt) |

*Mutationsprotokoll (Mutation per Diff gefahren, mit `git apply -R` zurückgenommen; die Tests waren vorher committet):* M1a Prüfung abgeschaltet → T1, Unrelated-Flag, Sandbox-Reihenfolge, Meldung, T5, `TestRun_…UnlistedFeatureSourceIsRejected` rot. M1b Prüfung nach `writeDevcontainerPhase` → `…FreigabeRejectsUnlistedSource`, `…BeforeSandboxRejection`, Binary-Test rot. M2 jede Quelle abgewiesen → T2, T3, mehrere Bestandstests rot. M3 Flag fließt nicht ein → `…ByFlagOfSameCall`, `…ByFlagDryRun`, Binary-Test rot. M4 deaktivierte Einträge einbezogen → `…IgnoresInactiveEntries`, Binary-Test rot. M5 eigener Vergleich in `generate` (`TrimRight "/"`) → `TestFeatureSourceAllow_ThreeCallersAgree` rot. M6 Reihenfolge vertauscht → `…BeforeSandboxRejection` rot. M7a `mapGenerateErrorToDiagnostic` ohne den Zweig → Binary-Test rot; M7b `ExitCode` ohne `ErrConfigValueInvalid` → Binary-Test rot. M8a bis M8g je ein Teilsatz der Meldung (Feature, Quelle, Allowlist-Satz, `config set`, Flag, Abschalten, Vergleichssatz) verändert → `…FreigabeMessageParts` rot, je mit eigener Fehlzeile. M9 Katalogschlüssel einbezogen → `…CatalogueKeyNeedsNoList` und viele Bestandstests rot. Nicht gefahren: Vertauschung der Reihenfolge der Abhilfen 1 und 2 in der Meldung (der Reihenfolge-Pin steht im Test, ist nicht rot gesehen).

*Wirkung am Binary* (Eintrag `evil`, Image `u-boot:latest` aus `make build`, Linux `amd64`): vor dem Code dasselbe Verhalten am Release-Image `0.15.0` (Exit-Code `0`, Schlüssel geschrieben; das Image aus `make` vor dem Code wurde nicht gebaut — Abweichung von Risiko 5, `generate.go` ist seit `v0.15.0` unverändert); danach Exit-Code `10`, `--json` Code [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), `--dry-run` Exit-Code `10`, `u-boot.yaml` byte-gleich, kein `.devcontainer`. Abhilfe 1 (`config set …allow`) → `generate` Exit-Code `0`; Flag → Exit-Code `0`, Quelle in der Liste; Abhilfe 2 (`…enabled false`) → Exit-Code `0`, kein Schlüssel.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/application` (APP; Prüfung, Test-Fixtures), `cmd/uboot` (Test, keine Sub-Area, `main_test.go`), nur lesend `internal/adapter/driving/cli` (CLI), sowie `spec/`, `docs/plan/` und `docs/user/` samt `CHANGELOG.md` (DOK). Alle sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) als Greenfield auditiert; eine zu grobe Sub-Area ist nicht dabei.

**Vorgelagert — offene Beobachtungen sichten** (Register `../observations/` gelesen am 2026-10-24; Zähler = Dateien in `evidence/`; die Treffer nach Stichwort „Freigabe“, „allowlist“, „featureSources“ und nach Sub-Area APP/DOK/CLI): [`BEO-APP/pruefung-nur-beim-lesen`](../observations/BEO-APP/pruefung-nur-beim-lesen/observation.md) (1; Klasse „eine Regel prüft nur auf einem der Pfade“ — hier umgekehrt: der Render-Pfad prüft sie nicht; ein Prädikat, beide Pfade fragen, Pin T5 über drei Aufrufer; das Auftreten dieses Slice ist ein Evidence-Kandidat und wird bei der Closure eingetragen, falls die Reviewerin es als zweites Vorkommen der Klasse liest); [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (16, verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10) und [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (37, verkörpert in §3.10): Beleg-Tabelle und Fundstellen (Risiko 4); [`BEO-DOK/abhilfe-ohne-wirkungsprobe`](../observations/BEO-DOK/abhilfe-ohne-wirkungsprobe/observation.md) (1): beide Abhilfen der Meldung werden ausgeführt (DoD 1); [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md): Nachlauf-Review (DoD 3). Die Zähler sind aus dem Ursprungs-Slice vom 2026-10-22 übernommen bzw. für `pruefung-nur-beim-lesen` am 2026-10-24 gelesen (eine Datei in `evidence/`); nicht neu gezählt, außer dem einen. Evidenz-/Diskrepanz-Risiko insgesamt: niedrig (Aufrufgraph gelesen, Lücke am Binary gezeigt; offen bleibt die Zahl der Fixtures).

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
