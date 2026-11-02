# Review-Report (Nachlauf): [`slice-v1-katalog-gepinnte-toolchain`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md) — 2026-11-02 (r2)

**Review-Art:** Nachlauf der Review-Fix-Runde (2. Lauf, frischer Kontext) gegen den
Erst-Report [`2026-11-02-slice-v1-katalog-gepinnte-toolchain.md`](2026-11-02-slice-v1-katalog-gepinnte-toolchain.md).

**Gegenstand:** `git show 8ab58b95` (Fix-Runde, „Berichtigungsrunde"; 8 Dateien,
+312/−14). Alles vor `8ab58b95` ist der gelesene Stand des Erst-Laufs
(Range `41667865..2ed2d956`, 29 Dateien). Der Erst-Report selbst kommt mit
diesem Commit erstmals in den Baum (+127, komplett neu) — sein Bestand ist hier
geprüft worden, nicht ein Diff gegen eine frühere Fassung.

**Skill:** `.harness/skills/reviewer.md` @ `8ab58b95` ·
**Modell:** `glm-5.3-flash:cloud[1m]` · **Datum:** 2026-11-02

---

## Stand der Vorgänger-Findings

| Finding | Stand | Beleg |
|---|---|---|
| F-1 (MEDIUM, Beleg-Name ohne Test) | behoben | [`TestFeatureCatalogue_ToolchainChecksumSources`](../../spec/spezifikation.md) im Baum: `internal/hexagon/application/generate_toolchain_test.go` Zeile 819 (Typgrenze über die Brücke `ToolchainVersionOverrideAcceptedForTest` je Fall, Last-Pfad mit `ValidateDevcontainerFeaturesForTest` für die zwei Schlüssel und Gegenprobe `gh` ohne Quelle); [`TestToolchainChecksumParsing`](../../spec/spezifikation.md) Zeile 876 mit Abschnitt **url forms** (Zeile 955); die Beleg-Tabelle der [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) steht auf den gelieferten Namen (Zeilen 1000, 1005, 1012) |
| F-2 (MEDIUM, Typgrenze/Brücken ohne Konsumenten) | behoben | alle drei Brücken werden von Tests gerufen (`GoChecksumsForTest` und `GolangciChecksumsForTest` je Subtest gegen Festwerte, `ToolchainVersionOverrideAcceptedForTest` in vier Synthetic-Fällen); der Last-Pfad läuft über `ValidateDevcontainerFeaturesForTest` (Accept beider Schlüssel, Abweisung von `gh` ohne Quelle). Die Zufassung der Nicht-Quelle ist im Test als nicht-offizielle Quelle geführt („no real key of the catalogue carries none; the case is synthetic", Zeile 843) — die leere Form fällt unter dasselbe Prädikat (`toolchainChecksumSourceSupported`, `devcontainer_toolchain.go` Zeilen 159 bis 167, Vergleich gegen die zwei Konstanten `checksumSourceGo`/`checksumSourceGolangci`); gelesen, die Pinntiefe ist größer als die Wortlaut-Fassung |
| F-3 (MEDIUM, Schreibweisen-Zusage ohne Test) | behoben | `TestGenerateDevcontainer_Toolchain_OverrideSpellings` (Zeile 1000): `go` + Tag-lose golangci-Schreibweise in der Konfiguration, erzeugter Block mit `assertToolchainBlock` je Schlüssel und `ARG GOLANGCI_LINT_VERSION=v2.13.2`; daneben Parser-Asserts beider Schreibweisen in beiden Subtests des Parser-Pins. Der [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)-Normalsatz (Zeile 980) blieb unverändert, die Berichtung steht in Zeile 1012; der Handbuch-Satz ist unberührt gelesen (`docs/user/benutzerhandbuch.md` Zeile 2323) |
| F-4 (MEDIUM, Plan-Tabelle unausgeglichen) | überwiegend behoben, Rest in F-N2 | `devcontainer_toolchain.go` (neu) und `devcontainer_features_test.go` (Update) tragen Plan-Zeilen; die driving-Zeile steht auf `unverändert` mit gelesener Fundstelle — der Abgleich bestätigt: kein `internal/adapter/driving`-Treffer in `git diff --stat 41667865..HEAD`; die e2e-Zeile nennt „fünf Tests" — gelesen, fünf `func TestE2E_Toolchain*` (Zeilen 61, 107, 186, 226, 323). Rest: die berichtigten Pin-Fundstellen der DoD-Zeile 2 tragen nicht (F-N2) und die DoD-Zeile 2 trägt einen doppelten Checkbox-Marker (F-N3) |
| F-5 (MEDIUM, Coverage-Tabelle ohne Bestands-Zeile) | behoben | [`harness/replay.md`](../../harness/replay.md) Zeile 61: Golden-Vertrag der zwei Schlüssel nach dem Muster der Vorgänger-Zeilen; die genannten Beleg-Namen liegen im Baum (`TestGenerateDevcontainer_Toolchain_*`, `TestToolchainPinPlacesEqual`, `TestDoctor_ToolchainKeys_NoFindings`, `TestRun_Generate_ToolchainNotices` und `TestRun_Generate_ToolchainConflictExits10` in `cmd/uboot/toolchain_test.go` Zeilen 18/70, fünf E2E-Tests); die Aufzählung der Abweisungen trägt auch die Zeichengrenze (`TestGenerateDevcontainer_Toolchain_InvalidVersionValueExit10`, Zeile 461) und den Nutzerpfad mit `--dry-run --json` (Zeilen 89/90) |
| F-6 (LOW, Kappe ohne Negativtest) | behoben | `TestGet_ResponseOverTheLimitIsAnError` (`internal/adapter/driven/fetch/fetch_test.go`): chunked Server ohne Content-Length streamt `64<<20+1` Bytes von `unitReader`, Assert auf den Fehler und den Wortlaut — der Code trägt genau diesen Wortlaut (`fetch.go` Zeile 81 „response over the %d-byte limit"); der Helfer `unitReader` steht nach der Funktion (§3.9 (c) erfüllt) |
| F-7 (LOW, Hebung-Aussage ohne auflösbare Herkunft) | behoben | Root-`Dockerfile` Zeilen 23 bis 31 und `repo.mk` Zeilen 47 bis 52 lösen die Hebung-Routine auf [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) §Konsequenzen auf — gelesen: Positiv trägt die Werkzeug-Hebung für ein Projekt ohne u-boot-Release per `version:` und die release-gekoppelte Rolle des Default-Pins. Der Teil „ohne eigenen ADR" bleibt laut der freigegebenen CR-Setzung bewusst Prozess-Satz im Code-Kommentar (Fundstelle [`cr-vorschlag-gepinnte-werkzeugkette.md`](../plan/planning/in-progress/cr-vorschlag-gepinnte-werkzeugkette.md) §Beleg: „Default-Hebung als Routine (kein Folge-ADR) … trägt der Code-Kommentar weiter … im Vertragstext steht sie bewusst nicht") — dokumentierte Setzung, das F-7-Ziel (auflösbare Routine-Aussage) ist erfüllt |
| F-8 (LOW, §-Aufzählung ohne neuen Typ) | behoben | [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) Aufzählung ungekennzeichneter Meldungen nennt „die Meldung zur Lieferkette je aktivem Schlüssel `go-toolchain` und `golangci-lint`" mit [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)-Verweis (Zeile 600); der Wortlaut der übrigen Aufzählung unverändert |
| F-9 (INFO, Laufzeit an der Jobgrenze) | offen | bleibt Verifier (CI-Laufzeit nach dem Push) |
| F-10 (INFO, Handbuch-Aussage außerhalb der Code-Belege) | offen | bleibt Validator |

## Neue Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-N1 | MEDIUM | Die behauptete „Anchor-Berichtigung" des Erst-Reports löst im committeten Blatt nicht: Plan-Zeile `docs/reviews/` („der Fundstellen-Verweis der Linse 3 zeigt auf [`SPEC-034`](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln) statt auf die Sicht") und Commit-Message („plus Anchor-Berichtigung des Review-Reports ([SPEC-034](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)-Verweis)") — das in diesem Commit erstmals committete Report trägt keinen [`SPEC-034`](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)-Verweis an der Linse 3 (Zeile 62) und im gesamten Dokument (Suchlauf `034`: kein Treffer). Die inhaltliche Absicht ist richtig — [`SPEC-034`](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln) ist die Fundstelle der Import-Regel-*Enforcement* in der Technik-Spec — aber die Aussage über den Zustand des Dokuments steht im Plan, und das Blatt trägt sie nicht. | [`AGENTS.md`](../../AGENTS.md) §3.10, Klasse [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) | [`slice-v1-katalog-gepinnte-toolchain.md`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md) §3 · „Anchor-Berichtigung in der Fix-Runde" | Review-only | Berichtigungs-Behauptung ohne Beleg im Baum |
| F-N2 | MEDIUM | Der §3.10-Abgleich trägt am Abgangs-Stand der Fix-Runde nicht: die in derselben Runde berichtigten Pin-Fundstellen der DoD-Zeile 2 sind durch das Kommentarwachstum desselben Commits schon wieder gewandert — „Root-`Dockerfile` Zeilen 39 bis 40" (Ist 41/42, `grep -n "ARG GO_VERSION="`), „`repo.mk` Zeilen 53 bis 54" (Ist 55/56); dazu weitere Fundstellen am Stand vor der Runde: DoD-Zeile 3 „Dockerfile Zeilen 23 bis 26 / `repo.mk` Zeilen 47 bis 48" (Kommentare umfassen 23–31 bzw. 47–52), §3-Preamble „`repo.mk` … die `?=`-Zeilen stehen bei 49 bis 50" (Ist 55/56) und die Plan-Zeile „Root-`Dockerfile` (Zeilen 23 bis 36), `repo.mk` (Zeilen 44 bis 50)". Der Abgleich „nach jeder späteren Änderung" schloss genau die Änderungen dieses Commits ein. | [`AGENTS.md`](../../AGENTS.md) §3.10, Klasse [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (F-4 derselben Klasse, Wiederholung) | [`slice-v1-katalog-gepinnte-toolchain.md`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md) §2 · „Root-`Dockerfile` Zeilen 39 bis 40" | Review-only (Fundstellen je Zeile nachgelesen) | Fundstellen-Berichtigung im eigenen Abgangs-Stand veraltet |
| F-N3 | LOW | Die DoD-Zeile 2 trägt einen doppelten Checkbox-Marker „- [ ] - [ ]" — durch das Nachzugs-Editing der Fix-Runde entstanden; sie bricht die DoD-Listenform des Vorlagen-Musters (`slice.template.md`: je DoD-Zeile ein „`- [ ]`"). Kein Gate fängt das: `docs-check` läuft grün über die Zeile (gemeldet, 0 Befunde, Stempel gedeckt). Failure-Szenario: Closure und Verifier lesen die DoD-Liste checkbox-getrieben; die zweite Liefer-Punkt-Zeile ist als malformed Marker nicht checkbox-getrieben. | DoD-Listenform des Musters [`slice.template.md`](../../.harness/baseline/v6.18.1/templates/docs/plan/planning/slice.template.md) | [`slice-v1-katalog-gepinnte-toolchain.md`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md) §2 · „- [ ] - [ ] **Erzeugnis und Pin-Stellen**" | Review-only | Doppelten Checkbox-Marker durch Nachzugs-Editing |
| F-N4 | INFO | Die Beleg-Tabelle der [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) führt den gelieferten Test `TestFeatureCatalogue_ToolchainChecksumSources` zusätzlich in der Zeile „Die zwei Schlüssel sind als eingebaute Features in beiden Profilen wählbar" auf (Zeile 1000), dessen Inhalt die Typgrenze S5/S6 prüft, nicht die Profil-Wahl — die Profil-Wahl tragen die beiden Katalog- und Render-Tests derselben Zeile. Name-zuerst-Ergebnis der F-1-Füllung; als Beleg für die Profil-Zusage bleibt die Zeile durch die beiden anderen Tests verständlich. | [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (Grenze des Musters) | [`spec/spezifikation.md`](../../spec/spezifikation.md) Beleg-Tabelle · „`TestFeatureCatalogue_ToolchainChecksumSources` und `TestFeatureCatalogue_KeysCoverSpecExamples`" | Review-only | Beleg-Name in Zeile mit anderem Gegenstand |

## Negativbefunde

Geprüft, ohne Befund:

- **§3.7/§3.9 über den Fix-Diff:** die geänderten Kommentare (Root-`Dockerfile`, `repo.mk`, Test-Kommentare) sind Indikativ über den geltenden Zustand, keine Chronik, keine „früher"/„no longer"-Formulierung, kein Verweis auf einen im Diff geänderten Namen oder Wert; Helfer stehen nach ihren Funktionen (`unitReader` in `fetch_test.go`); alle Einfüge-Zeilen des Diffs tragen gerade Backtick-Parität (Zählung pro Zeile, kein ungerader Treffer)
- **Beleg-Tabelle der [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) und Handbuch-Satz inhaltlich unverändert:** die Fix-Runde ändert nur die Berichtung in Zeile 1012 (Testname ergänzt, Beobachtungs-Satz umgestellt); der Normalisierungs-Satz (Zeile 980) und der Handbuch-Satz (Zeile 2323) blieben wortgleich — keine Zusage erweitert
- **F-2-Prädikat:** Pinntiefe größer als die Wortlaut-Fassung (siehe F-2 in der Stand-Tabelle) — die Zusage ist nicht schmaler geworden
- **Architektur:** die Fix-Runde ändert keine Importe, keine Schicht, kein Wiring; „driving unverändert" bestätigt (kein Treffer in `git diff --stat`); der einzige `net/http`-Ort bleibt der `fetch`-Adapter (Bestand des Erst-Laufs, ungeändert)
- **Exit-Klassen:** keine `internal/adapter/driving`-Änderung in der Fix-Runde; zentrale Klassifikation und `generate`-Diagnostic-Abbildung unverändert — kein Dual-Classifier-Bruch
- **Lastenheft und [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md):** beide in der Fix-Runde unberührt (kein Treffer in `git show --stat`); der Grabstein [`spec/lastenheft.md`](../../spec/lastenheft.md) bleibt unangetastet
- **Handbuch/README/CHANGELOG:** in der Fix-Runde nicht berührt — „nur wo belegt" geprüft, ohne Befund
- **Gate-Änderungen:** keine (Makefile, `repo.mk`-Rezepte, `.golangci.yml`, `.d-check.yml` in der Fix-Runde außer den Kommentaren unberührt)
- **Coverage-Zeile der [`harness/replay.md`](../../harness/replay.md):** Beleg-Namen je Vorgänger-Zeilen-Muster gegen den Baum gelesen (alle genannten Tests existieren), Zeilensitzung nach dem Muster (Bedingung → Vertrag → Beleg-Pfade)

Nicht geprüft:

- `make gates` — nicht ausgeführt (Review-only); der Stempel deckt: `.harness/state/gates-passed.head` = `8ab58b95` = `git rev-parse HEAD` (mündlich am 2026-11-02, 13:11, nach dem Commit 13:09), und `gates-passed.diffsha` entspricht dem frisch berechneten Arbeitsbaum-Hash `working-tree-hash.sh` (`33f21f45…` = Wert im Stempel, Baum clean)
- `make test` / `make docs-check` — berichtet (alle Pakete ok, 0 Befunde) und durch denselben Stempel gedeckt; nicht wiederholt
- `make test-docker` — kein eigener Lauf (Review-only): die Fix-Runde berührt keinen Build-Tag-`docker`-Pfad (`git show --stat` listet keine `internal/e2e`-Datei); die fünf E2E-Werkzeugkette-Tests sind seit ihrem Umsetzungs-Commit unverändert; der berichtete Gesamtlauf (1697,8 s, F-9) bleibt Verifier-Kontext
- `make verify-depguard` — die Import-Aussagen sind hier per `git show --stat` und `grep` gelesen (keine Import-Änderung in der Fix-Runde); der Bestands-Nachweis (9/9, Auftragskontext des Erst-Laufs) steht nicht zur Wiederholung an
- CI — kein Lauf auf dem Umsetzung-/Fix-Stand: `gh run list` nennt als neuesten `main`-Stand das Push von `41667865` (Vor-Umsetzung, beide Workflows `ci` und `integration-docker` grün, 2026-11-02T11:03); `HEAD` (`8ab58b95`) ist nicht gepusht — der Push nach diesem Review zieht den CI-Lauf nach; Verifier-Auftrag: CI-Ausgang auf `8ab58b95`, Laufzeit-Analyse nach F-9
- DoD-Closure, Verification-Evidence, Golden-/Replay-Evidence-Block, Risiko-Ausgänge §6 und die drei Paarungen — Verifier; hier nur die Review-Prüflinsen

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Berichtigungs-Behauptung ohne Beleg im Baum
(BEO-DOK/plan-aussage-ueber-code) · Fundstellen-Berichtigung veraltet
(= F-4-Klasse, wiederholter Auftritt) · Doppelten Checkbox-Marker (Beleg:
slice.template.md DoD-Form) · Beleg-Name in Zeile mit anderem Gegenstand.

## Verdikt

**Merge-blockierend: ja, für die Slice-Closure.** Alle acht Befunde des
Erst-Laufs sind inhaltlich behoben (F-1 bis F-8, siehe Stand-Tabelle; F-7 mit
der gelesenen CR-Setzung erfüllt), und der Gate-Stempel deckt den
Fix-Commit-Kopf samt Arbeitsbaum — aber die Fix-Runde trägt zwei eigene
MEDIUM-Befunde: die Plan-Zeilen behaupten eine Report-Berichtigung (F-N1) und
Fundstellen (F-N2), die der committete Stand nicht trägt; beide sind schreibbar
in einer Berichtigungsrunde ohne Codepfad (Fundstellen der DoD-Zeilen 2/3 und
der §3-Runde berichtigen, F-N1 an Plan-Zeile oder Blatt klären,
Checkbox-Marker F-N3 mitziehen). F-N4 ist INFO ohne Blockierwirkung.

**Übergabe an den Verifier:** Fix-Runde r2 für F-N1/F-N2 (mit F-N3), danach
§3.10-Abgleich noch einmal gegen `git diff --stat` auf dem Abgangs-Commit; der
Stop-Hook-Stempel muss nach der Berichtigung neu geschrieben werden (neuer
HEAD). Verifier-Aufträge unverändert: CI-Ausgang auf dem gepushten Stand,
Laufzeit-Analyse nach F-9, Verification-Evidence und Closure; F-10 an den
Validator.