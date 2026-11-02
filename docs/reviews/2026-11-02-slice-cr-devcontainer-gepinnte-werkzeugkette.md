# Review-Report: slice-cr-devcontainer-gepinnte-werkzeugkette — 2026-11-02

**Review-Art:** Plan-Review des CR-Diffs — geprüft wird der Schreib-Lauf
(Lastenheft 0.30.0, [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv),
Planning-Nachzug, §7-Verifikation) gegen den freigegebenen
[CR-Vorschlag](../plan/planning/in-progress/cr-vorschlag-gepinnte-werkzeugkette.md),
die Hard Rules ([`AGENTS.md`](../../AGENTS.md) §3) und die Rollenentscheidungen.
Keine DoD-Closure — die prüft der Verifier
([`harness/roles.md`](../../harness/roles.md) §Kontext-Trennung).

**Gegenstand:** Commit-Range `cc028b69..e3490c65` (vier Commits: `0c214395`,
`56469d8b`, `8f704cfb`, `e3490c65`) gegen den Slice
[`slice-cr-devcontainer-gepinnte-werkzeugkette`](../plan/planning/done/slice-cr-devcontainer-gepinnte-werkzeugkette.md).

**Skill:** `.harness/skills/reviewer.md` @ Commit `e3490c65` ·
**Modell:** `glm-5.3-flash:cloud[1m]` · **Datum:** 2026-11-02

**Eingangs-Kontext:**

- [Slice-Plan](../plan/planning/done/slice-cr-devcontainer-gepinnte-werkzeugkette.md) inkl. §7-Verifikation des Implementers
- [CR-Vorschlag](../plan/planning/in-progress/cr-vorschlag-gepinnte-werkzeugkette.md) (Freigabe-Block, §§2 Stellen 1 bis 8, §4 F1 bis F6)
- [`spec/lastenheft.md`](../../spec/lastenheft.md) ([`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), neue [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), Grabstein [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik)), [`spec/spezifikation.md`](../../spec/spezifikation.md) ([`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)), [`spec/architecture.md`](../../spec/architecture.md)
- [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) (`Accepted`), [`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md), [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) (nur Referenzebene)
- Messprotokoll [`docs/messungen/slice-v1-katalog-gepinnte-toolchain-entwurf/README.md`](../messungen/slice-v1-katalog-gepinnte-toolchain-entwurf/README.md) (m1c/m2c/m3/m3z)
- Beobachtungen: [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md), [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md)
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules §3), [`harness/conventions.md`](../../harness/conventions.md) (ID-Schema, Sub-Area-Modi), [`harness/review.md`](../../harness/review.md)
- Umsetzungs-Slice [`slice-v1-katalog-gepinnte-toolchain`](../plan/planning/open/slice-v1-katalog-gepinnte-toolchain.md) (Konsument der offenen Punkte)

---

## Findings

Schema nach `.harness/skills/reviewer.md` §Output-Schema
([`harness/review.md`](../../harness/review.md) ist die kanonische Prosa).

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Planning-Nachzug `8f704cfb` upgrade-te die Kennungs-Links, ließ aber die `d-check:ignore`-Marker am [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)-Link an zwei Stellen im Umsetzungs-Slice stehen; dieselbe Klasse in der Titel-Zeile des CR-Slice („Anker entsteht mit diesem CR") und an der bare-Kennung in [`done/slice-v1-katalog-gepinnte-toolchain-entwurf.md`](../plan/planning/done/slice-v1-katalog-gepinnte-toolchain-entwurf.md) Zeile 114. Der Marker-Text „Anker entsteht mit dem Lastenheft-Commit" trägt nicht: der Anker besteht seit `0c214395`, und die Commit-Message von `8f704cfb` wie der CR-Slice-Absatz behaupten, die Marker an den Entwurfs-Kennungen seien entfallen. | [`AGENTS.md`](../../AGENTS.md) §3.10 (Zustandsmarker mit Ablaufpunkt) | [`docs/plan/planning/open/slice-v1-katalog-gepinnte-toolchain.md`](../plan/planning/open/slice-v1-katalog-gepinnte-toolchain.md) · „`<!-- d-check:ignore (Kennung des Entwurfs; Anker entsteht mit dem Lastenheft-Commit) -->`" — zwei Stellen, Zeilen 14 und 26 | ja — `make docs-check` nach Marker-Entfernung (die Links lösen bereits; der Gate-Lauf deckt dann die bislang unterdrückte Prüfung) | d-check:ignore-Marker übersteht seinen Ablaufpunkt |
| F-2 | MEDIUM | [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) Abschnitt *Eintragstyp* sagt, `u-boot doctor` melde für die neuen Schlüssel keinen Befund der Prüfungen `devcontainer.features.allowlist` und `devcontainer.features.drift` — eine testbare Verhaltensbedingung, die in der Beleg-Tabelle (zehn Zeilen) keinen pinnenden Test trägt; die als erledigt markierte DoD-Zeile 2 und die Spec-Historie-Zeile behaupten, je Bedingung stehe der Zeiger auf den pinnenden Test. | [`AGENTS.md`](../../AGENTS.md) §3.10 (letzte Aufzählung), [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) | [`spec/spezifikation.md`](../../spec/spezifikation.md) · „`u-boot doctor` meldet für sie keinen Befund der Prüfungen `devcontainer.features.allowlist` und `devcontainer.features.drift`" | ja — `make test` am Umsetzungs-Diff (Beleg-Tabelle-Zeile oder Streichung des Satzes), bis dahin Review-only | Spec-Bedingung ohne pinnenden Test-Zeiger (BEO-DOK/spec-bedingung-ohne-code-entsprechung) |
| F-3 | MEDIUM | Die §7-Sentence „ein Zeilenversatz ergab sich nicht" trägt nicht: der Commit verschiebt die Fundstellen — Grabstein [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) steht nun bei Zeilen 894 bis 896 (vorher 884 bis 886), Historie-Zeilen 0.27.0 bis 0.29.0 bei 1624 bis 1626 (vorher 1614 bis 1616), `## 4.4` bei Zeile 526 (vorher 516). Die Zeilenangaben im CR-Vorschlag (Zeilen 26 und 51), im CR-Slice §1 Zeile 26 und im Umsetzungs-Slice DoD-„Hebung-Abgleich" („Fundstelle Zeilen 884 bis 886 bleibt Grabstein") sind damit für jeden Leser nach `0c214395` veraltet; der Nachzug `8f704cfb` korrigierte keine Zeilenangabe, obwohl §3.9 (b) einen Suchlauf über den ganzen Baum verlangt. | [`AGENTS.md`](../../AGENTS.md) §3.9 (b), §3.10; [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) | [`docs/plan/planning/done/slice-cr-devcontainer-gepinnte-werkzeugkette.md`](../plan/planning/done/slice-cr-devcontainer-gepinnte-werkzeugkette.md) §7 · „ein Zeilenversatz ergab sich nicht" | Review-only (Zeilen-Abgleich per `grep -n`, hier durchgeführt: 884→894, 1614→1624, 516→526) | State-Claim über Fundstellen-Abgleich trägt nicht (BEO-DOK/plan-aussage-ueber-code) |
| F-4 | LOW | Die Plan-Tabelle §3 führt die Komponente „`docs/plan/planning/done/` (Referenz-Nachzug)" als `update`; die Range berührt keine Datei unter `done/`, und die §7-Abgleicherklärung „jede Datei trägt einen Plan-Eintrag, kein Eintrag ohne Diff" behandelt die Richtung „Eintrag ohne Diff" nicht — eine Berichtigung des Eintrags (auf den Ist-Stand oder „entfällt") unterblieb. | [`AGENTS.md`](../../AGENTS.md) §3.10 | [`docs/plan/planning/done/slice-cr-devcontainer-gepinnte-werkzeugkette.md`](../plan/planning/done/slice-cr-devcontainer-gepinnte-werkzeugkette.md) §3 · „`docs/plan/planning/done/` (Referenz-Nachzug)" | Review-only (`git diff --stat cc028b69..HEAD`, hier durchgeführt: sieben Dateien, keine unter `done/`) | Plan-Eintrag ohne Diff nicht berichtigt |

**Failure-Szenarien.** F-1: der nächste Änderer benennt die
[`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)-Überschrift um oder löscht sie; `docs-check`, der das mit
dem Marker unterdrückte Ende prüfen würde, bleibt grün, und der tote Link
fällt lautlos an — dieselbe Blende, die den Zustandsmarker-Regelbefund der
Beobachtung `BEO-DOK/plan-aussage-ueber-code` trug (Vorgang
[`slice-v2-spec-devcontainer-sshd-claude-code`](../plan/planning/done/slice-v2-spec-devcontainer-sshd-claude-code.md), dritter Auftritt dieser
Klasse). F-2: die Umsetzung oder ein Folge-Commit ändert das Doctor-Verhalten
für einen der Schlüssel; die Bedingung „kein Befund" ist durch keinen Test
vertraut, die Drift fällt erst durch manuelles Nachlesen auf — genau das
Versagen, das die Verkörperung von §3.10 an dieses Muster trug. F-3: die
Umsetzungs-Slice-DoD liest „Zeilen 884 bis 886" für den Grabstein-Kommentar
am Root-`Dockerfile`-Nachzug und greift an der um zehn Zeilen verschobenen
Stelle daneben; die Closure trägt die §7-Aussage „kein Zeilenversatz" als
geprüften Abgleich, obwohl der versprochene Baum-Suchlauf nicht belegt ist.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| **Linse 1 — Wortlaut-Treue, Stelle 1**: die zwei Beispiel-Zeilen und der Eintrags-Absatz (S1 bis S11) gegen den CR-Vorschlag | geprüft, ohne Befund — byte-identisch nach den drei angekündigten Abweichungen; Nachvollzug `cmp` über die extrahierten Passagen (`cmp`: Stelle 1 Absatz 2542 Bytes, Stelle 2 Absatz 744 Bytes, Stelle 3 Absatz 912 Bytes, Überschrift und Beispiel-Zeilen identisch) |
| **Linse 1 — Wortlaut-Treue, Stelle 2** (Überschrift + Absatz [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)) | geprüft, ohne Befund — Abweichung 3 belegt: der Blockquote-Titel des CR-Vorschlags trug vor `8f704cfb` keine Link-Klammer (nachgelesen an `cc028b69`), der Vertragstext ist unverändert |
| **Linse 1 — Wortlaut-Treue, Stelle 3** (`containerEnv`-Absatz in [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)) | geprüft, ohne Befund — byte-identisch |
| **Linse 1 — Positionen**: Stellen-Einsetzpunkte (nach `resources`/vor `Mindestdateien` Zeile 387; nach `gh`/`glab`/vor „Auf Wunsch"; nach [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)/vor `## 4.4`) | geprüft, ohne Befund |
| **Linse 2 — Kein Abgleich nach unten**: [`docs/user/benutzerhandbuch.md`](../../docs/user/benutzerhandbuch.md), `README.md`, `README.de.md`, [`CHANGELOG.md`](../../CHANGELOG.md) | geprüft, ohne Befund — `git diff --stat cc028b69..HEAD` auf die vier Pfade leer; Plan-Tabelle trägt „im CR kein Diff" je Zeile; Roadmap- und open/README-Zeilen tragen die Entscheidung |
| **Linse 3 — Decken-Regel im Lastenheft**: Suchlauf `SPEC-`, `ADR-`, `slice-`, `welle-`, `BEO-`, `CO-` | geprüft, ohne Befund — Treffer nur „ADR-Format" als Begriff in [`LH-FA-PROJDOCS-002`](../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format) (Bestand); [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) verweist nicht abwärts (nur `LH-*` und `SPEC-*` des gleichen Stratum); [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) nur in Planning-Dokumenten verlinkt |
| **Linse 3 — Referenz-Richtung der Verweis-Form**: [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) → [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) Absatz *Grenzen* (Auflösung `containerEnv` zu `Config.Env`) | geprüft, ohne Befund — Fundstelle `spec/spezifikation.md` Zeile 882 trägt die Aussage |
| **Linse 3 — Meldungs-Code [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)** | geprüft, ohne Befund — Fundstellen `internal/hexagon/application/devcontainer_features.go` Zeilen 635/655/666 (`forge`-CLI-Meldungen tragen den Code als `code`), Muster gelesen; der offene Punkt des CR-Slice an den Umsetzungs-Slice bleibt dokumentiert |
| **Linse 4 — ID-Schema**: [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) (Familie `LH-FA-DEV-*`, nächste freie Nummer nach [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) in §4.3), [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) (nächste freie nach [`SPEC-053`](../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs)) | geprüft, ohne Befund |
| **Linse 4 — Grabstein**: [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) bleibt Grabstein; keine Wiederbelebung in den vier Commits | geprüft, ohne Befund — die Kennung taucht im neuen Vertragstext nicht auf; Fundstelle unverändert |
| **Linse 4 — Historien**: Lastenheft-Zeile `0.30.0` (2026-11-02, Kennungen → Kennung, Grenz-Sätze, Pin-Zahlen als Chronik nach CR Stelle 7) und Spec-Historie-Zeile `2026-11-02` | geprüft, ohne Befund — konform zum Muster der Vorgänger-Zeilen |
| **Linse 5 — Zusage-Menge gegen Messung**: die Grenz-Sätze S10/S11 und der Stelle-3-Grenz-Halbsatz gegen das Messprotokoll (m1c löst `${containerWorkspaceFolder}`, m2c `containerEnv` übersteuert Block-`ENV`, m3/m3z Recreate übersteht, nur Sandbox-Profil gemessen) | geprüft, ohne Befund — die Zusage-Menge trägt die Protokoll-Grenze; die Lese-Prüfung des Protokolls an drei Prüfstellen nachvollzogen |
| **Gate-Beleg** | nachgelesen, ohne eigener Lauf — `.harness/state/gates-passed.head` = `e3490c65` (= HEAD), diffsha gesetzt; der Stempel trägt den Stand inklusive §7-Commit |
| **§3.7/§3.9 über den neuen Vertrags- und Spec-Wortlaut** | geprüft, ohne Befund — keine Kommentar-Klassen-Brüche im Doku-Diff; die Historie-Zeilen sind Vertrags-Historie nach dem Bestand, keine Chronik-Prosa über abwesende Texte |

**Nicht geprüft:**

- `make gates` / `make ci` / `make fullbuild` — Review-only, keine Gate-Läufe
  aus diesem Lauf; Beleg ist der nachgelesene Stempel (`gates-passed.head` =
  `e3490c65`) und der Implementer-Lauf (zweiter Versuch grün an `8f704cfb`,
  nach §7); der CI-Lauf auf dem Push-Tip ist Auftragskontext, nicht nachgelesen.
- `make test-docker` — kein Code- oder Docker-Pfad-Diff in der Range
  (Plan-Tabelle „ungeöffnet / unverändert", vom Diff belegt).
- `make verify-depguard` — keine Importregel berührt.
- Voller Satz-für-Satz-Abgleich der [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)-Beleg-Tabelle
  gegen die Rohdaten des Protokolls — Verifier- und Umsetzungs-Kontext; hier
  nur die Grenz-Aussagen gegen die Protokoll-Zusammenfassung gelesen.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** d-check:ignore-Marker übersteht seinen
Ablaufpunkt · Spec-Bedingung ohne pinnenden Test-Zeiger
(BEO-DOK/spec-bedingung-ohne-code-entsprechung) · State-Claim über
Fundstellen-Abgleich trägt nicht (BEO-DOK/plan-aussage-ueber-code) ·
Plan-Eintrag ohne Diff nicht berichtigt

## Verdikt

**Merge-blockierend:** ja — die drei MEDIUM-Findings blockieren die
Slice-Closure (nicht die drei Contract-Commits: der freigegebene Wortlaut ist
byte-identisch umgesetzt, die Straten- und Grabstein-Disziplin ist ohne
Befund). Die Behebung ist Prosa-Nachzug in Planning-Dokumenten und eine
Beleg-Tabelle-Zeile in [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv);
berührt die Fix-Runde dennoch eine Lastenheft-Zeile, greift die
Nachlauf-Review-Pflicht der DoD-Zeile 3.

**Übergabe:** Findings an den Implementer; die Finding-Klassen zusätzlich
in die Slice-Closure §7 und von dort in den Beobachtungs-Register-Zähler.
Für die Closure fehlen danach: Fix-Runde zu F-1 (Marker entfernen; Titel-Zeile
des CR-Slice und `done/entwurf` nach demselben Ablauf), F-2 (Beleg-Tabelle-Zeile
oder Streichung des Doctor-Satzes), F-3 (Zeilenangaben auf den Ist-Stand,
§7-Satz berichtigen), F-4 (Plan-Eintrag ohne Diff berichtigen); dann
Verification-Report, Closure-Notiz mit Steering-Loop-Lerneintrag,
Beobachtungs-Register-Fortschreibung, Risiko-Ausgänge und die drei Paarungen
der DoD-Zeile 3; ein `make gates`-Lauf am Fix-HEAD mit `record-gates`-Stempel.
Kein Self-Review dieses Reports — der Folge-Lauf bei einer Fix-Runde ist ein
neuer Report (`-r2`), keine Überschreibung.