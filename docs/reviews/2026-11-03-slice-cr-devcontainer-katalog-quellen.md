# Review-Report: slice-cr-devcontainer-katalog-quellen — 2026-11-03

**Review-Art:** Code-Review (Vertrags-Diff) — geprüft wird der freigegebene CR-Wortlaut
gegen den übernommenen Vertrags-Text und die Hard Rules.

**Gegenstand:** CR-Schnitt [slice-cr-devcontainer-katalog-quellen](../plan/planning/in-progress/slice-cr-devcontainer-katalog-quellen.md),
Range `d3c24eb6..HEAD` (Ab Mess-Commit `d3c24eb6`): Schnitt `4c5cd730` (Vorschlag + CR-Plan),
reine Moves `dd9db460`/`d6877b7f` (R100) mit Verweis-Nachzug je eigenem Commit `621dba91`/`362dab38`,
**Vertrags-Übernahme `151a7976`** ([`spec/lastenheft.md`](../../spec/lastenheft.md) Beispiele-Zeile,
Abweisungs-Bullet, Versionsfeld `0.33.0`, Historie-Zeile; [`spec/spezifikation.md`](../../spec/spezifikation.md)
[SPEC-023](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) Katalog-Aufzählung
und [SPEC-040](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)
Schreib-Pipeline-Schritt 6), Link-Berichtigung `137f2f92`.
Wortlaut-Norm: [`cr-vorschlag-katalog-quellen.md`](../plan/planning/open/cr-vorschlag-katalog-quellen.md)
(Status *freigegeben wie vorgelegt*, Projektinhaber 2026-11-03). Mess-Basis:
[`docs/messungen/slice-v1-katalog-features-quellen-und-freigabe/README.md`](../messungen/slice-v1-katalog-features-quellen-und-freigabe/README.md) (M1/M2).

**Skill:** `.harness/skills/reviewer.md` (Modul 10, Baseline v6.18.1) ·
**Modell:** `glm-5.3-flash:cloud[1m]` · **Datum:** 2026-11-03 · frischer Kontext, kein Self-Review.

**Eingangs-Kontext:** der CR-Slice-Plan und die WIP-Slot-Bewegung des quellen-Slices
[slice-v1-katalog-features-quellen-und-freigabe](../plan/planning/open/slice-v1-katalog-features-quellen-und-freigabe.md),
[`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
[`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes),
[ADR-0032](../plan/adr/0032-katalog-features-socket-und-lockerungen.md),
`AGENTS.md` (Hard Rules §3.4/§3.7/§3.9/§3.10), `harness/conventions.md` (Modus-Deklaration),
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Roadmap-Absatz des quellen-Slices behauptet „liegt in `in-progress/`", während die Slice-Datei nach `open/` zurückgemoved wurde; zugleich trägt die Roadmap den aktiven CR-Schnitt [slice-cr-devcontainer-katalog-quellen](../plan/planning/in-progress/slice-cr-devcontainer-katalog-quellen.md)` gar nicht (Muster der drei Vorgänger-CRs: je Absatz bzw. Start-Zeile, zuletzt [slice-cr-devcontainer-containerenv](../plan/planning/done/slice-cr-devcontainer-containerenv.md)). | AGENTS.md §3.7/§3.10 (Zustandsmarker auf Ist-Stand) | `docs/plan/planning/in-progress/roadmap.md` · „liegt in `in-progress/`" (Absatz zum quellen-Slice) | review-only — d-check prüft Link-Ziele, nicht Präsens-Zustandsangaben | Zustandsmarker widerspricht Position |
| F-2 | MEDIUM | Zustandsmarker über den quellen-Slice sind am Ist-Stand abgelaufen und wurden beim Rück-Move nicht berichtigt: open/README-Zeile „die Katalogquellen von `cpp` und `postgres-client` (laut ADR bei `ghcr.io` unbekannt, **nicht gemessen**)" (gemessen am 2026-11-03, Wege entschieden) und Slice-Plan-Welle-Zeile „(CR-Schnitt liegt in `in-progress/`)" ohne Kennung — der gemeinte CR ([slice-cr-devcontainer-containerenv](../plan/planning/done/slice-cr-devcontainer-containerenv.md)) ist geliefert, der aktive CR ist ein anderer. Der Vorgänger-CR ließ diesen Refresh als eigenen Marker-Refresh-Commit laufen; hier entstand er nicht. | AGENTS.md §3.10 (Zustandsmarker mit Ablaufpunkt) | `docs/plan/planning/open/README.md` · „(laut ADR bei `ghcr.io` unbekannt, nicht gemessen)" | review-only | Abgelaufener Zustandsmarker |
| F-3 | MEDIUM | [SPEC-023](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) trägt weiterhin „für `cpp` und `postgres-client` sind sie nicht gelesen, ihre Felder sind leer" — derselbe Abschnitt, dessen Katalog-Aufzählung die zwei Einträge in diesem Commit verlor. Das widerspricht Lastenheft `0.33.0` („die Katalog-Einträge … entfallen") und der eigenen Tabelle; die Zeile liest sich fälschbar als Mess-Stand über existierende Einträge. | Source Precedence (Spec driftet gegen Lastenheft) | `spec/spezifikation.md` · „für `cpp` und `postgres-client` sind sie nicht gelesen, ihre Felder sind leer" | review-only (d-check prüft keine Inhaltskohärenz) | Vertriebszeile zu entfernten Katalogeinträgen |
| F-4 | MEDIUM | Vorschlags-Stelle (c) ist nur zur Hälfte übernommen: die Quelle-Form-Aufzählung („Externe Feature-Quellen müssen als URL mit Schema `http://`, `https://` oder `oci://` angegeben werden.") trägt **keine** gemessene Verarbeitungstabelle; die Form-Abweisung steht in der Spec nur im `config set`-Pipeline-Schritt 6. Lastenheft `0.33.0` trägt die Abweisung zuerst für `u-boot generate devcontainer` — für diesen Generator-Pfad ist in der Spec keine Form-Prüfung verankert. Szenario: der Umsetzungs-Diff liest [SPEC-023](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) als vollständige Formregel und setzt die Abweisung nur in der Schreibpipeline um. | Freigegebener Wortlaut (Vorschlags-Stelle c) + CR-Plan §3 („Form-Tabelle") | `spec/spezifikation.md` · „Externe Feature-Quellen müssen als URL mit Schema `http://`, `https://` oder `oci://` angegeben werden." | review-only | Freigegebene Stelle nicht übernommen |
| F-5 | LOW | Die Historie-Zeile `0.33.0` wurde in der Lastenheft-Historie **vor** der `0.32.0`-Zeile eingefügt; alle übrigen Zeilen laufen aufsteigend. Lesereihenfolge der Chronik bricht an der Naht. | AGENTS.md §3.7 (Historie als Chronik — Zustand und Beleg lesbar) | `spec/lastenheft.md` · „\| 0.33.0 \| 2026-11-03 \| …" | review-only | Chronik-Einfügeposition |
| F-6 | LOW | Übernahme nicht wortgleich je Stelle, ohne benannte Änderung (Status „freigegeben wie vorgelegt"): Stelle (b) blendet den zweiten Satz des Vorschlags aus („Der Tarball-Weg der CLI bleibt Offenpunkt; sein Messergebnis kommt in die Historie."); Stelle (d) ist umformuliert („Quellen existieren nicht, gemessen" → „ihre Katalogquellen existieren nicht — Messung …"; „gemessen mit der CLI des Repos" → „gemessen:"). Inhaltlich gedeckt (Tarball-Grenze steht in Bullet-Klammer und Historie), aber der §7-Abgleich des CR muss sie je Stelle als benannte Abweichung führen. | CR-Plan §2 DoD („wortgleich … Abgleich je Stelle in §7") | `spec/lastenheft.md` · Abweisungs-Bullet und Historie-Zeile `0.33.0` | review-only | Unbenannte Wortlaut-Abweichung |
| F-7 | LOW | Die Historie-Zeile `0.33.0` verweist auf die Mess-Basis als „Messung `docs/messungen/`" — ein unauflösender generischer Pfad (welche Messung?), als Backtick-Text ohne Link und ohne Kennung; kein Gate prüft Backtick-Pfade, das Lastenheft darf die Planungs-Kennung nicht nennen. Verbal mit Datum wäre die Herkunft auflösbar genug. | AGENTS.md §3.7 (Herkunft in **ein** auflösbarem Feld) | `spec/lastenheft.md` · „(ihre Katalogquellen existieren nicht — Messung `docs/messungen/`)" | review-only — `make docs-check` bestätigt die Prüf-Lücke indirekt | Unauflösender Beleg-Verweis |
| F-8 | LOW | Im Schreib-Pipeline-Schritt 6 steht das Literal „[Beleg-Auftrag]" als Klammer-Token im Fließtext. Das etablierte Muster des Dokuments ([SPEC-054](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) Beleg-Tabelle) trägt den Beleg-Auftrag als Bedingungsspalte mit der Zusage, dass der Testname im Umsetzungs-Diff entsteht — nicht als einmaliges Einschub-Token im Vertrags-Text, das kein Gate liest und dessen Verbleib/Form im Umsetzungs-Diff offen bleibt. | [SPEC-054](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) Beleg-Tabelle (konventionsbildende Form im selben Stratum) | `spec/spezifikation.md` · „(die Registry-Tarball-Form ist offen und bleibt beobachtet) [Beleg-Auftrag]" | review-only | Beleg-Auftrag als Literal-Token |
| F-9 | LOW | Der CR-Plan §3 listet die berührten Verweis-Träger-Dateien der zwei Moves nicht: 22 eingehende Verweise auf den quellen-Slice (14 `done/`-Slices und Welle-Dateien, 3 ältere Review-/Verify-Berichte, 2 Messungs-READMEs, ein weiterer `open/`-Slice) und die quellen-Slice-Datei selbst fehlen je Zeile; nur „dieser Slice, `roadmap.md`, `open/README.md` (Bestand)" sind gelistet. Der §3.10-Abgleich gegen `git diff --name-status` ist je Fix-Runde fällig. | AGENTS.md §3.10 (Plan-Tabelle listet jede berührte Datei) | der CR-Plan ([slice-cr-devcontainer-katalog-quellen](../plan/planning/in-progress/slice-cr-devcontainer-katalog-quellen.md)) · Plan-Tabelle §3 | review-only | Plan-Tabelle unvollständig gegen Move-Nachzüge |
| F-10 | LOW | Die Spezifikations-Historie (§7) trägt keinen Eintrag zur `0.33.0`-Vertrags-Änderung ([SPEC-023](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) Katalog-Aufzählung und Schritt 6); der CR-Plan §3 nennt „Historie" für `spec/spezifikation.md`, und der Vorgänger-Vertrags-Eintrag `2026-11-02` ([SPEC-054](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)) etabliert die Zeile je Vertrags-Änderung. | CR-Plan §3 („Stelle (c) — …; Historie") + Spec-§7-Muster | `spec/spezifikation.md` §7 Historie · letzte Zeile `2026-11-01` | review-only | Fehlende Spec-Historien-Zeile |
| F-11 | INFO | Das neue Abweisungs-Bullet trägt wieder einen Messstand im Anforderungstext („(gemessen am 2026-11-03 mit der CLI des Repos: …)") — die Lastenheft-Richtung `0.24.0` hatte Messbezüge aus den Anforderungstexten entfernt; der Projektinhaber hat den Wortlaut so freigegeben, die Spanne (Beleg-Klammer statt Zusage-Bedingung) steht zur Beobachtung, keine Änderungspflicht. | Lastenheft-Historie `0.24.0` (Messbezüge aus Anforderungstexten) | `spec/lastenheft.md` · „(gemessen am 2026-11-03 mit der CLI des Repos: …)" | review-only | Messbezug im Anforderungstext |
| F-12 | INFO | Der Schritt-6-Satz trägt "(gemessen 2026-11-03: …)" und "(die Registry-Tarball-Form ist offen und bleibt beobachtet)" in den Vertrags-Text der Spec — die Spec-Historie-Zeile `2026-10-19` hat die Spec „auf Zielbild" gebracht („tragen keine Messstände, Messumgebungen und Nicht-Messungen mehr"). Spannung zur etablierten Richtung; zuständig Architect/Planner bei der Umsetzung, der freigegebene Wortlaut überwiegt. | Spec-Historie `2026-10-19` (Zielbild-Richtung) | `spec/spezifikation.md` · „(gemessen 2026-11-03: `oci://…:<tag>` und `https://…:<tag>` unverarbeitbar)" | review-only | Messstand im Spec-Vertrags-Text |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Stellung (a): Beispiele-Zeile — Restsuche „PostgreSQL Client" über `spec/lastenheft.md` | geprüft, ohne Befund: einzig verbleibenden Treffer trägt die Historie-Chronik (Zeile `0.33.0`), kein Rest im Zusage-Text; Zusage-Bullets sonst unverändert |
| Stellung (b): Formen in Wortlaut, Vertrags-Text und Messung (`oci://…:<tag>`, `https://…:<tag>`) | geprüft, ohne Befund: deckungsgleich der Form nach; die Messung trägt beide Formen mit existierendem Tag (M2, `go:1.0.3`, CLI `0.89.0` „CLI des Repos"), je eigener Fehlerklasse |
| Stellung (b): Tarball-Grenz-Formulierung („offen") | geprüft, ohne Befund: Lastenheft-Bullet-Klammer, [SPEC-040](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)-Klammer und Historie-Zeile tragen die Offenheit; das Protokoll führt den Tarball-Weg als „nicht gemessen" — keine erfundene Messung, keine Namen-Form-Aussage über den Tarball im Vertrags-Text |
| Stellung (b): Zusage-Spanne „Form verallgemeinert über die Probe-Quelle hinaus" | geprüft, ohne Befund: genau dies benennt CR-Risiko 2 („Zusage geht weiter als die Messung") mit Vorsorge; dem Review bleibt die Spanne am Umsetzungs-Slice beobachtbar |
| Decken-Regel (AGENTS.md §3.4): Down-Referenzen in Lastenheft/Spezifikation | geprüft, ohne Befund: Historien-Spalten „Freigabe des Projektinhabers … (Wortlaut wie vorgelegt, CR der Katalogquellen)" und die `0.32.0`-Zeile („CR der Projekt-Umgebungswerte") nennen die Herkunft **verbal** ohne Kennung oder Link; keine `SPEC-*`-Kennung im Lastenheft (`grep` 0 Treffer) — nur F-7 bleibt am generischen Pfad |
| Lastenheft-Struktur: Listenform, Ankerform [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | geprüft, ohne Befund: das Bullet endet die Liste „Für optionale externe Feature-Quellen gilt:" (nach der java-toolchain-Zusage), Anker-Link-Wortlaut konform |
| `grep -n 'cpp\|postgres' spec/spezifikation.md` | geprüft, mit Befund nur F-3: die Treffer Z.228 (Keycloak-„external-postgres"), Z.1238/1375–1393 (postgres-Service, Managed-Blocks) und Z.1415–1442/1573 (`postgres`-Befehle/Service-Katalog) sind Services, nicht Katalog-Features — Negativbefund |
| Pipeline-Schritt 6: grammatikalischer Einbau und Beleg-Auftrag-Form ohne erfundenen Testnamen | geprüft, ohne Befund für den Satzbau (eingebettet zwischen Freigabe- und `--yes`-Teilsatz); zur Token-Form „[Beleg-Auftrag]" siehe F-8 |
| Kein Produktcode, kein Nutzertext, kein Gate im Range | geprüft, ohne Befund: `git diff --name-status d3c24eb6~1..HEAD` berührt nur `spec/`, `docs/plan/`, `docs/messungen/` (Verweis-Nachzug), `docs/reviews/` (Verweis-Nachzug); kein `internal/`, `cmd/`, `templates/`, `docs/user/`, README, `CHANGELOG.md`, Makefile, Workflow — Negativbefund auch für §3.9 (keine Code-Kommentar-Berührung im Vertrags-Range) |
| Lifecycle der WIP-Slot-Bewegung (Baseline `modul-05-planning-harness.md` §Lifecycle als State Machine) | geprüft, ohne Befund für die Bewegung selbst: `in_progress → open` ist ein erlaubter Übergang; `dd9db460`/`d6877b7f` je reiner Move (R100), Verweis-Nachzug je eigenem Commit (`621dba91`/`362dab38`), quellen-Slice vor dem CR aus `in-progress/` hinaus (WIP-Limit 1), CR-Start mit `Verantwortlich:`-Feld — die **Lese-Wahrheit** der Roadmap trägt F-1 |
| Plan §3: `open/README.md` (Bestand des quellen-Slices) | geprüft, mit Befund im Umfeld siehe F-9 (die Zeile ist im Plan gelistet; die übrigen Nachzugs-Dateien fehlen) |
| Status-Zeile des Vorschlags („freigegeben wie vorgelegt", 2026-11-03) | geprüft, ohne Befund: Zustandsmarker mit Zustand und Beleg, kein Ablaufpunkt; Änderung im Vertrags-Commit je Commit-Message |
| Gates/Carveouts: keine Lockerung, kein neuer Carveout | geprüft, ohne Befund: keine Gate-Datei im Range, keine CO-Datei fällig |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 6 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zustandsmarker widerspricht Position · Abgelaufener Zustandsmarker · Vertriebszeile zu entfernten Katalogeinträgen · Freigegebene Stelle nicht übernommen · Chronik-Einfügeposition · Unbenannte Wortlaut-Abweichung · Unauflösender Beleg-Verweis · Beleg-Auftrag als Literal-Token · Plan-Tabelle unvollständig gegen Move-Nachzüge · Fehlende Spec-Historien-Zeile

## Sensors

| Sensor | Ergebnis |
|---|---|
| `make docs-check` | gelaufen nach dem Schreiben dieses Reports; Ergebnis am Ende des Laufs im Handoff genannt (siehe Zusammenfassung der Antwort) |
| `make lint`, `make test`, `make coverage-gate` | nicht gelaufen — Grund: keine Produktcode-/Gate-Datei im Range (`--name-status` geprüft); kein Nutzertext-Diff |
| `make test-docker` | nicht gelaufen — Grund: kein Docker-Pfad berührt |
| `make verify-depguard` | nicht gelaufen — Grund: keine Import-Regel berührt |

## Verdikt

**Merge-blockierend:** ja — die vier MEDIUM-Findings blockieren die CR-Closure normalerweise: F-3/F-4 sind Vertrags-Text-Befunde in gerade geändertem Vertrags-Bestand (Spec-Drift gegen Lastenheft `0.33.0` bzw. unvollständige Übernahme der freigegebenen Stelle (c)), F-1/F-2 lassen die Roadmap die Lifecycle-Wahrheit (in-progress = nur der CR) nicht tragen. Kein HIGH: der Vertrags-Satz des Lastenhefts trägt die freigegebene Zusage je Stelle vollständig und mit der Tarball-Grenze; die Befund-Familien haben konkrete Szenarien am nächsten Änderer (Umsetzungs-Diff, nächster Lifecycle-Move), nicht am laufenden Gate.

**Übergabe:** Findings an den Implementer (Fix-Runde des CR vor dem §7-Abgleich); die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den Zähler. DoD-Closure und Verification bleiben Verifier-Aufgabe — dieser Report prüft den Vertrags-Diff gegen den freigegebenen Wortlaut und die Hard Rules, nicht die DoD.