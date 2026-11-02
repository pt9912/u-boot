# Verification-Report: slice-cr-devcontainer-gepinnte-werkzeugkette — 2026-11-02

**Rolle:** Verifier (`harness/roles.md`) — Modell `glm-5.3-flash:cloud[1m]`, Datum 2026-11-02. Frage: *"Bauen wir es richtig?"* — DoD und Entscheidungs-Konformität gegen Plan und Review-Fix-Stand.

**Gegenstand (Scope):** die sechs Commits `cc028b69..fd65d6f7` — `0c214395` (Lastenheft 0.30.0, Stellen 1 bis 3), `56469d8b` ([`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)), `8f704cfb` (Planning-Nachzug), `e3490c65` (Slice-§7-Verifikation), `cd4f7c11` (Review-Fix-Runde zu F-1 bis F-4), `fd65d6f7` (d-check-Berichtigung am Review-Report) — gegen den Slice [`slice-cr-devcontainer-gepinnte-werkzeugkette`](../plan/planning/done/slice-cr-devcontainer-gepinnte-werkzeugkette.md), den [CR-Vorschlag](../plan/planning/done/cr-vorschlag-gepinnte-werkzeugkette.md) und den [Review-Report](2026-11-02-slice-cr-devcontainer-gepinnte-werkzeugkette.md) vom 2026-11-02.

**Beleg-Quellen:** eigene Nachlesung am Stand `fd65d6f7` (Datei-für-Datei), Byte-Vergleich der Wortlaut-Stellen, `git diff cc028b69..fd65d6f7`, Stempel `cat .harness/state/gates-passed.head`, CI-Nachlesung per `gh run list`. Kein Blindvertrauen in die Implementer-Behauptungen: Wortlaut, Fundstellen, Marker-Restbestand und Doku-Diff sind von diesem Lauf selbst geprüft. Der Report liegt uncommitted im Arbeitsbaum (Auftrag: keine Commits von der Verifier-Rolle).

---

## DoD-Abgleich (Slice §2)

| DoD-Zeile | Stand | Beleg |
|---|---|---|
| 1 — Lastenheft `0.30.0`, Stellen 1 bis 8 wortgleich | **erfüllt** | Byte-Vergleich dieses Laufs (unten Abschnitt Wortlaut-Verifikation): Eintrags-Absatz in [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) und der `containerEnv`-Absatz in [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) sind nach den drei angekündigten Abweichungen (Satz-Marker, Link-Form, Leerzeichen-Normalisierung) **byte-identisch** gegen den CR-Vorschlag; die zwei Beispiel-Zeilen stehen (Position: nach `glab`, vor `ollama-host`); Versionsfeld Zeile 9 = `0.30.0`; Historie-Zeile an Zeile 1627. Decken-Regel: Suchlauf des Laufs über die hinzugefügten Lastenheft-Zeilen nach `SPEC-`, `ADR-`, `slice-`, `welle-`, `BEO-`, `MR-`, `CO-` — 0 Treffer. Grabstein [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) unangetastet (Fundstelle Zeilen 894 bis 896, gelesen). Commit `0c214395`. |
| 2 — Technik [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv), kein Abgleich nach unten | **erfüllt, mit Rest (Befund V-2)** | [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) liest am Stand: Prüfsummen-Quellen je Schlüssel (Tabelle Zeilen 983 bis 986), URL-Form als Umsetzungs-Stelle (988), drei Abweisungen je Exit-Code `10` (990, F3-Stand), Cache-Werte unter `${containerWorkspaceFolder}` (992), Meldungen samt [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)-Kennung und kein `--quiet`-Ausschluss (994), Beleg-Tabelle mit Beleg-Aufträgen (996 bis 1009), Nicht-zugesagt (1011), Spec-Historie-Zeile (1593). Je Bedingung trägt der Zeiger auf den pinnenden Test — **Rest:** die Doctor-Bedingung `allowlist`/`drift` (977) hat zwar den Beleg-Auftrag inline, aber keine eigene Zeile in der Beleg-Tabelle, obwohl der Satz auf sie verweist (V-2 unten). **Kein Abgleich nach unten bestätigt:** `git diff cc028b69..fd65d6f7 -- docs/user/ README.md README.de.md CHANGELOG.md` ist leer (Projektinhaber-Entscheid 2026-11-02 getragen); die Plan-Tabelle Zeilen 54 bis 56 trägt die Berichtigung je Zeile. Commits `56469d8b`, `0c214395`. |
| 3 — Gates grün, Review, Verifier-Report, Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | **teilweise** | `make gates`: Stempel `.harness/state/gates-passed.head` = `fd65d6f7` = HEAD; CI-Lauf `ci` am Push-Tip `fd65d6f7` **success** (Run 41429191226). Review-Report liegt vor, Fix-Runde `cd4f7c11`/`fd65d6f7` abgeschlossen (F-1 bis F-4 unten geprüft). Verifier-Report: diese Datei. **Offen** (das ist das Closure-Paket selbst, kein Vermessungs-Gegenstand): Closure-Notiz mit Steering-Loop-Lerneintrag (§7 markiert „nicht geschlossen"), Fortschreibung des Beobachtungs-Registers oder §7-Vermerk, Ausgang je §6-Risiko im abgeschlossenen Stand, die drei Paarungen. |

## Wortlaut-Verifikation (Probe gegen den CR-Vorschlag, §§2)

Sonden dieses Laufs mit sed-Normalisierung (CR-Satz-Marker `(**Sn**)` entfernt, Leerzeichen zusammengeführt, Kennungs-Links von der Datei-Form auf die Anker-Form normalisiert) und `cmp`:

- **Stelle 1** — Eintrags-Absatz in [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features): 2511 Zeichen je Seite, **byte-identisch**. Zwei Beispiel-Zeilen je einmal in CR und Lastenheft: `Go-Toolchain aus gepinntem Release`, `golangci-lint mit gepinntem Release` (je eine Fundstelle je Datei).
- **Stelle 2** — Überschrift und Absatz [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette): 659 Zeichen je Seite (Link-Form normalisiert), **byte-identisch**; die Überschrift steht exakt so im Lastenheft (Zeile 520).
- **Stelle 3** — `containerEnv`-Absatz in [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen): 913 Zeichen je Seite, **byte-identisch** ohne Normalisierung des Links (kein Link im Absatz).
- **Einsetzpunkte gelesen:** Stelle 3 nach dem `resources`-Absatz, vor `Mindestdateien:`; Stelle 1 nach dem `gh`/`glab`-Absatz, vor `Auf Wunsch des Nutzers …`; Stelle 2 nach [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), vor `## 4.4` (Zeile 526). Ohne Befund.
- Die drei angekündigten Abweichungen des §7 (Satz-Marker fortgelassen, Anker statt Datei-Link, Link-Fassung der Stelle-2-Überschrift im CR-Vorschlag) sind zur Kenntnis genommen und durch die Sonden als genau diese Klassen bestätigt; die Exit-Code-Zuordnung im Lastenheft trägt durchgehend `10` je Fehlerfall (S5, S6, S8) — F3-Stand, keine 13/10-Reste im Vertragswortlaut.

## Review-Findings am Fix-Stand

| Finding | Prüfung dieses Laufs | Ergebnis |
|---|---|---|
| F-1 (Marker über Ablaufpunkt) | `grep -rn "Anker entsteht" docs/plan/planning/` — **0 Treffer**. An den drei F-1-Stellen: [`slice-v1-katalog-gepinnte-toolchain`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md) Zeilen 14 und 26 (Diff `cd4f7c11`) tragen jetzt echte Links, der CR-Slice-Titel ebenso, [`slice-v1-katalog-gepinnte-toolchain-entwurf`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain-entwurf.md) (Diff `cd4f7c11`) trägt den Anker-Link. Die verbleibenden `d-check:ignore`-Marker im Planning-Baum sind eine andere Klasse (Pfad-Zitate, Zitat-Formen) und nicht abgelaufen. | **behoben** |
| F-2 (Doctor-Bedingung ohne pinnenden Test) | [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) Zeile 977 trägt den Beleg-Auftrag inline (` gepinnt mit einem Test im Umsetzungs-Slice; Beleg-Tabelle Zusage → Test unten `). **Rest:** die Beleg-Tabelle (Zeilen 998 bis 1009, zehn Zeilen) hat keine Zeile für die Doctor-Bedingung `allowlist`/`drift` — der Verweis „unten“ löst auf nichts. | **überwiegend behoben, Rest → V-2** |
| F-3 (Fundstellen-Spanne) | Selbst nachgelesen: Grabstein [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) bei Zeilen 894/896, Historie `0.27.0` bis `0.29.0` bei Zeilen 1624/1626, `## 4.4` bei 526, Versionsfeld 0.30.0 bei Zeile 9 — deckungsgleich mit der berichtigten §7-Aussage und dem CR-Vorschlag (Zeilen 26 und 51). | **behoben** |
| F-4 (Plan-Eintrag ohne Berichtigung) | Die Plan-Tabelle §3 trägt die Berichts-Zeilen (`docs/reviews/` als `neu`, Nachtrag-Vermerk) und die „done/ (Referenz-Nachzug)“-Zeile ist entfernt. **Rest:** die Fix-Runde `cd4f7c11` berührte selbst eine Datei ohne Plan-Zeile — siehe V-1. | **behoben, Rest → V-1** |

## Sensors

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass | `.harness/state/gates-passed.head` = `fd65d6f7` — der Stempel deckt den HEAD ab (Nachlese-Beleg, kein eigener Lauf nötig); dazu Implementer-Lauf grün an `8f704cfb` (§7, zweiter Versuch nach Registry-502) und CI-Lauf `ci` **success** an `fd65d6f7`. |
| CI am Push-Tip | pass | `gh run list` am Commit `fd65d6f7`: `ci` success (Run 41429191226); `integration-docker` zum Report-Moment noch in_progress (kein Code-Diff in der Range — Rest-Risiko gering, siehe unten). |
| Vorgeschichteter Befund | — | Der CI-Lauf am Vorgänger-Commit `cd4f7c11` fiel mit `docs-check` (9 Befunde — die bare-Backtick-Kennungen des Review-Reports) aus; `fd65d6f7` berichtigt genau das, und der Tip-Lauf ist grün. Die Nachlauf-Review-Pflicht (Lastenheft-Diff in Fix-Runde) greift nicht: kein Lastenheft-Diff in den Fix-Commits. |
| `make docs-check` | pass (über `gates`/CI) | Implementer-Lauf: 1023 Dateien, 0 Befunde (§7); CI am Tip grün. |
| `make test-docker` | not run | kein Code-, Docker-Pfad- oder Template-Diff in der Range (`git diff --stat cc028b69..fd65d6f7`: neun Dateien, alle Doku/Planning/Spec-Text) |
| `make verify-depguard` | not run | keine Importregel berührt (kein Code-Diff) |
| `make gates` (eigener Lauf) | nicht nötig | Stempel deckt HEAD; Auftrag sah eigenen Lauf nur bei Stempel-Lücke vor — keine Lücke |

## Plan-vs-Code (`git diff --stat cc028b69..fd65d6f7` gegen Plan-Tabelle §3)

Neun Dateien im Range; je Zeile ein Plan-Eintrag: `spec/lastenheft.md`, `spec/spezifikation.md`, `roadmap.md`, `open/README.md`, [Umsetzungs-Slice](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md), [CR-Vorschlag](../plan/planning/done/cr-vorschlag-gepinnte-werkzeugkette.md), diese Slice-Datei (`diese Datei`), Review-Report unter `docs/reviews/` (Zeile `docs/reviews/`, Nachtrag-Vermerk). **Eine Datei ohne Plan-Zeile:** [`slice-v1-katalog-gepinnte-toolchain-entwurf`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain-entwurf.md) (2 Zeilen, F-1-Fix) — Befund V-1 unten. „Unverändert“ geführt und durch denselben Diff belegt: `spec/architecture.md`, `internal/**`, `cmd/**`, Root-`Dockerfile`, `repo.mk`, `.devcontainer/Dockerfile`, `docs/user/`, beide READMEs, `CHANGELOG.md`, `docs/plan/adr/` (keine ADR-Datei im Range — [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) ist `Accepted` und unverändert, gelesen). Diese Datei (Verifier-Report) geht nicht in den Abgleich ein — sie liegt uncommitted außerhalb des Ranges.

## Befunde dieses Verifier-Laufs

| ID | Kategorie | Befund | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|
| V-1 | LOW | Die F-1-Fix-Datei [`slice-v1-katalog-gepinnte-toolchain-entwurf`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain-entwurf.md) (2 Zeilen in `cd4f7c11`) trägt keine Plan-Tabelle-Zeile; der F-4-Fix entfernte die „done/“-Zeile, und der anschließende F-1-Fix schuf genau einen „done/“-Diff, der nicht nachgetragen wurde. Der Fundstellen-Abgleich in §7 erklärt „jede Datei trägt einen Plan-Eintrag" für den Stand `e3490c65` und wiederholt sich nicht nach der Fix-Runde. | [Slice §3/§7](../plan/planning/done/slice-cr-devcontainer-gepinnte-werkzeugkette.md) | ja — `git diff --stat cc028b69..fd65d6f7` gegen die Plan-Tabelle (dieser Lauf) | Plan-Abgleich nach letzter Fix-Runde nicht wiederholt — Vorstufe von [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md), Berichtigung im Closure-Commit |
| V-2 | LOW | Die Doctor-Bedingung der Eintragstyp-Prüfungen `devcontainer.features.allowlist` und `devcontainer.features.drift` trägt in [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) den F-2-Beleg-Auftrag inline, steht aber in keiner Zeile der Beleg-Tabelle (998 bis 1009) — der Auftrags-Verweis „Beleg-Tabelle unten“ löst für sie auf nichts; ein Leser sucht und findet die Zusage→Test-Zeile nicht. | `spec/spezifikation.md` Zeile 977 gegen Zeilen 998 bis 1009 | ja — Lesung des Abschnitts am Fix-Stand (dieser Lauf) | Beleg-Tabelle-Zeile fehlt zum Bedingungs-Satz — Zuordnung an den [Umsetzungs-Slice](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md): die Zeile kommt mit dem pinnenden Doctor-Test, oder der inline-Verweis wird auf „beim Umsettzen ergänzt“ berichtigt |
| V-3 | INFO | Vor der Schwelle, außerhalb des Scopes: in [`slice-v1-devcontainer-start-ohne-cli`](../plan/planning/done/slice-v1-devcontainer-start-ohne-cli.md) Zeile 211 trägt die [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)-Überschrift den Marker „im Lastenheft erst nach dem CR“ — der Anker besteht seit dem CR der Version 0.25.0; dieselbe Ablaufpunkt-Klasse wie F-1 in archivierter Historie. Kein Bestandteil dieses Slices; Kandidat für ein späteres Register-Ereignis. | [`slice-v1-devcontainer-start-ohne-cli`](../plan/planning/done/slice-v1-devcontainer-start-ohne-cli.md) | ja — `grep -rn "d-check:ignore (Kennung des Vorschlags" docs/plan/planning/done/` | Marker-Derivat über Ablaufpunkt in Archiv-Prosa |

## Traceability

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | Lastenheft 0.30.0, Stelle 1 (Commit `0c214395`, byte-identisch zu CR); Technikaufträge in [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) |
| [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) | Stelle 3 (`containerEnv`-Zusätze), Commit `0c214395`; Cache-Bedingungen in [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) |
| [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) | Neue Anforderung in §4.3 (Stelle 2, Commit `0c214395`); Pinnungs-Auftrag in [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) Zeile 979 und Beleg-Tabelle Zeile 1002 — Hebung der Pin-Werte bleibt beim [Umsetzungs-Slice](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md) |
| [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | Exit-Code `10` je Abweisungsfall (S5, S6, S8) — Lastenheft und Spec-Zeile 990 |
| [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) | Grabstein bleibt Grabstein (Fundstelle gelesen, keine Wiederbelebung im Diff) |
| [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) | `Accepted`, unverändert im Range; in Planning verlinkt, nicht im Spec-Stratum (Decken-Regel, Sondierung 0 Treffer) |
| [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | Meldungs-Kennung als `code` in [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) Zeile 994; Code-Gegenprobe mit der Umsetzung (offener Punkt §7) |

## Carveouts

Neu: keine. Gelöst: keine. Unverändert: [`carveouts.md`](../plan/carveouts/README.md) liegt nicht im Diff; kein Gate gelockert, keine Schwelle gesenkt (der Lastenheft-Diff trägt ausschließlich freigegebenen Wortlaut).

## Nicht ausgeführt

- `make gates` (eigener Lauf) — Stempel `fd65d6f7` deckt HEAD; der Auftrag schrieb einen eigenen Lauf nur bei Stempel-Lücke vor.
- `make test-docker` — kein Code- oder Docker-Pfad-Diff (Plan-Tabelle „ungeöffnet / unverändert“, durch `git diff --stat` belegt).
- `make verify-depguard` — keine Importregel berührt.
- Satz-für-Satz-Abgleich der vollen Beleg-Tabelle gegen die Rohdaten des Messprotokolls — Umsetzungs-Kontext; hier die Grenz-Sätze gegen die von Review- und Implementer-Lauf gelesenen Protokoll-Aussagen geprüft.

## Commit / Artefakte

- Range `cc028b69..fd65d6f7`: `0c214395`, `56469d8b`, `8f704cfb`, `e3490c65`, `cd4f7c11`, `fd65d6f7`.
- Gate-Stempel `.harness/state/gates-passed.head` = `fd65d6f7` · CI-Run `ci` 41429191226 success am Tip.
- Diese Datei: Verifier-Report (uncommitted im Arbeitsbaum, keine Commits von dieser Rolle).

## Rest-Risiken (an den Umsetzungs-Slice übernommen, aus §7 und der Beleg-Tabelle)

1. **Doctor-Rest (V-2):** die Zusage→Test-Zeile zu den Doctor-Prüfungen `allowlist`/`drift` fehlt in der Beleg-Tabelle — der [Umsetzungs-Slice](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain.md) ergänzt sie mit dem pinnenden Doctor-Test, oder der inline-Verweis wird berichtigt, bevor er ein stiller Tot-Verweis wird.
2. **F6-ENV-Rest:** Träger der Umgebungswerte (Hand-written-Rest oder eigener `containerEnv`-Schritt) entscheidet der Umsetzungs-Slice (Freigabe-Punkt F6, keine Zusage im CR).
3. **URL-Form des Prüfsummen-Abrufs je Version** — Umsetzungs-Stelle, je Schlüssel ein Docker-Test; für andere Umgebungen unbelegt und nicht zugesagt.
4. **`containerEnv`-Unterpfade** (Muster `go-build`/`go-mod`) — Umsetzungs-Stelle, Docker-Test-Pflicht.
5. **Meldungs-Code-Kennung:** der CR sagt „Kennung als `code`"; die Zuordnung zu [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) ist im Umsetzungs-Diff gegen den Code zu bestätigen (Meldungs-Muster `internal/hexagon/application/devcontainer_features.go`).
6. **Pin-Politik-Kommentar-Nachzug:** Root-`Dockerfile` und `repo.mk` bekommen [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) statt der Grabstein-Kennung — bis dahin bleibt die tote Kennung im Code-Kommentar (Doku-Zitat-Formen im Planning verweisen bewusst darauf).
7. **Nutzer-Doku mit Implementierung:** Handbuch, READMEs und `CHANGELOG.md` tragen den Ist-Zustand erst mit dem Umsetzungs-Slice (Wortlaut-Vorlage CR-Vorschlag §8); bis dahin beschreiben sie die zwei Schlüssel nicht — der geprüft-leere Diff ist der Stand.
8. **CI `integration-docker` auf dem Tip** lag zum Report-Moment in progress; die Range trägt keinen Code und keinen Docker-Pfad, ein Einfluss auf den CR-Nachweis ist nicht erkennbar (Rest-Risiko: erst nach Abschluss formell belegt).

## Gesamturteil

**Closure-fähig — ja, mit zwei kleinen Berichtigungen im Closure-Paket.** DoD-Zeile 1 und 2 sind erfüllt und von diesem Lauf belegt: der freigegebene Wortlaut ist byte-identisch in Lastenheft und Spec, die Straten-Disziplin (Decken-Regel, Referenz-Richtung, Grabstein) ist ohne Befund, der Projektinhaber-Entscheid (kein Abgleich nach unten) ist im Plan und im Diff getragen, Gates und CI sind am HEAD grün. DoD-Zeile 3 ist zu dem durch den Verifier prüfbaren Teil erfüllt (Gates, Review inklusive Fix-Runde, Verifier-Report); offen ist ausschließlich das Closure-Paket selbst: **V-1-Berichtigung** (Plan-Zeile für die done-Entwurf-Datei des F-1-Fixes), Closure-Notiz mit Steering-Loop-Lerneintrag, Beobachtungs-Register-Fortschreibung, Ausgang je §6-Risiko und die drei Paarungen. **V-2** geht als Auftrag mit an den Umsetzungs-Slice (Beleg-Tabelle-Zeile mit dem Doctor-Test) — er blockiert die CR-Closure nicht, solange der inline-Beleg-Auftrag in [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) stehen bleibt und nicht als erledigt gelesen wird.