# Reviewer-Skill — u-boot

* Status: Accepted
* Bezug: [`harness/review.md`](../../harness/review.md) (kanonische Review-Prosa),
  [`AGENTS.md`](../../AGENTS.md) §Hard Rules,
  [`harness/roles.md`](../../harness/roles.md) §Reviewer
* Gilt für: jeden Review-Lauf gegen einen u-boot-Diff (Plan-, Design- oder
  Code-Review). u-boot hat **kein** `agent-review`-Make-Target; der Lauf ist
  agentisch und wird über seinen Report belegt.

Diese Datei ist die *scharfe Kurzform* für den Reviewer-Lauf, nicht die
Vollkopie der Review-Prosa. Bei Abweichung gilt
[`harness/review.md`](../../harness/review.md) — dort stehen die vollständigen
Linsen-, Kategorie- und Steering-Loop-Tabellen. Zweck dieser Datei: Ein
Reviewer ohne Skill-Datei driftet zwischen Sessions (gleiche Eingabe → andere
Findings); hier steht das repo-spezifische „worauf achtest du".

Für die engere Closure-Prüfung gibt es den Schwester-Skill
[`closure-note-reviewer.md`](closure-note-reviewer.md). Report-Gerüst pro Lauf:
`.harness/baseline/v3.5.1/templates/docs/reviews/review-report.template.md`
(Kopiervorlage, siehe [`docs/reviews/README.md`](../../docs/reviews/README.md)).

## Kontext-Eingang (Pflicht)

Was der Reviewer *immer* mitbringt, bevor er den Diff liest:

- Diff bzw. Commit-Range des Vorhabens
- der aktive Slice-Plan (Scope, DoD, explizite Nicht-Ziele)
- [`spec/lastenheft.md`](../../spec/lastenheft.md) für die referenzierten
  `LH-*`-IDs, [`spec/architecture.md`](../../spec/architecture.md) für
  Schicht-/Portgrenzen
- ADRs, deren Kennung im Diff, im Slice oder in der Commit-Message vorkommt
- [`AGENTS.md`](../../AGENTS.md) §Hard Rules und
  [`harness/conventions.md`](../../harness/conventions.md) (Adaptions-Ledger
  `MR-*`, Modus-Deklaration pro Sub-Area)
- vorherige Findings am selben Bereich (letzte Reports in
  [`docs/reviews/`](../../docs/reviews/))

Ohne Plan- und Spec-Kontext ist der Lauf Codekritik, kein Harness-Review — dann
wird er als „nicht durchgeführt" berichtet, nicht als „ohne Befund".

## Klassifikation

Vollständige Kategorie-Definitionen: [`harness/review.md`](../../harness/review.md)
§Finding Categories. Hier die Anker, gegen die *dieser Repo-Diff* geprüft wird.

**HIGH** — eines der folgenden:

- **Dual-Classifier-Bruch (repo-spezifisch #1).** Ein neuer oder aufgeteilter
  Driving-Sentinel ist nicht in *allen* zuständigen Klassifikatoren des
  Driving-Adapters eingetragen. Zuständig sind: die **eine, zentrale**
  Exit-Code-Klassifikation **und** die Diagnostic-Abbildung **jedes**
  Subkommandos, über dessen Pfad der Sentinel auftreten kann — von diesen
  Abbildungen gibt es **eine pro Subkommando**, nicht eine gemeinsame. Ein
  querliegender Sentinel (z. B. „Projekt nicht initialisiert") steht damit an
  der zentralen Stelle plus in *mehreren* Kommando-Abbildungen; prüfe beim
  Review ausdrücklich, ob alle betroffenen Kommandos nachgezogen wurden.
  Andernfalls driften Exit-Code
  ([`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)) und
  maschinenlesbare Ausgabe
  ([`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe))
  je nach Subkommando auseinander. Ebenfalls HIGH: alle Einträge vorhanden, aber
  ohne Exit-Code-Pin-Test — dann fällt die Drift beim nächsten Split lautlos
  aus. Kein Gate prüft diese Abdeckung; sie ist Review-Aufgabe.
- **Referenzmodell-Verstoß (repo-spezifisch #2).** Ein Dokument der Ränge 1–3
  (Lastenheft, Architektur-Sicht, ADR) verweist **abwärts** auf Slice, Carveout
  oder Roadmap. Normative Kraft existiert nur aufwärts
  ([`LH-FA-PROJDOCS-006`](../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell));
  die Änderungskopplung deklariert die ADR im `Schärft`-Feld nach oben, nie die
  Spec nach unten.
- **Sprachvertrag verletzt.** CLI-Ausgabe, Fehlermeldung oder generierte
  Datei ist nicht Englisch
  ([`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)) — auch wenn
  Plan- und Spec-Dokumente deutsch sind. Der Fehler tritt typisch im
  Grenzbereich auf: deutscher Code-Kommentar ist zulässig, deutscher
  `error`-String nicht.
- **Architekturgrenze gebrochen.** Hexagonale Import-Regel, kreuz-blinde
  Port-Trennung oder Wiring-Grenze verletzt
  ([`LH-FA-ARCH-003`](../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement)),
  oder `//nolint:depguard` als Umgehung.
- **Gate- oder Harness-Lockerung.** Docker-only-Workflow umgangen
  ([`LH-FA-BUILD-007`](../../spec/lastenheft.md#lh-fa-build-007--docker-only-workflow)),
  Coverage-Schwelle oder Lint-Regel still gesenkt, Host-Toolchain-Lauf als
  alleiniger Handoff-Nachweis.
- **Dateisicherheit.** User-Dateien werden ohne managed-block-Schutz, Backup
  oder Spec-konforme Bestätigung überschrieben
  ([`LH-NFA-REL-001`](../../spec/lastenheft.md#lh-nfa-rel-001--kein-stilles-überschreiben),
  [`LH-SA-FILE-002`](../../spec/lastenheft.md#lh-sa-file-002--markierte-verwaltete-bereiche)).
- **Carveout ohne Buchführung.** Neue temporäre Ausnahme ohne Eintrag im
  Carveout-Inventar *und* ohne Plan-Anker
  ([`LH-FA-PROJDOCS-005`](../../spec/lastenheft.md#lh-fa-projdocs-005--carveout-disziplin)).
- **Immutabilitäts-Bruch.** Ein `Accepted`-ADR wird inhaltlich umgeschrieben
  statt durch eine Folge-Entscheidung abgelöst; oder ein grandfathered lean-ADR
  wird ohne Change Request auf die MADR-Form migriert
  ([`LH-FA-PROJDOCS-002`](../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)).
- **Lifecycle-Bruch.** Planning-Artefakt liegt in zwei Lifecycle-Verzeichnissen
  oder wurde per Copy-Delete statt `git mv` bewegt
  ([`LH-FA-PROJDOCS-003`](../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)).
- **Baseline-Bump nur teilweise.** Ein Versions-Bump der vendorten Baseline
  fasst nicht alle vier Stellen der `MR-004`-Bump-Prozedur an
  ([`harness/conventions.md`](../../harness/conventions.md)).

**MEDIUM** — eines der folgenden:

- Neuer öffentlicher CLI-Pfad ohne Negativtest oder ohne Sentinel-/
  Exit-Code-Abdeckung.
- Fehlerbehandlung plausibel, aber keiner Exit-Code-Kategorie `2/10/11/12/14`
  eindeutig zugeordnet.
- Doku (README, `docs/user/`, CHANGELOG) driftet gegen Spec oder gegen das
  implementierte Verhalten.
- Test belegt Verhalten, aber nicht die tragende `LH-*`-ID.
- Eine Sub-Area wird berührt, deren Modus-Aussage in
  [`harness/conventions.md`](../../harness/conventions.md) den Diff nicht mehr
  trägt (Modus-Drift).
- Ein LOW-Muster wiederholt sich und wird zum Drift-Signal.

**LOW** — Wartbarkeit, Doku-Präzision, einmalige Inkonsistenz ohne
Vertragsbruch.

**INFO** — Beobachtung ohne Änderungspflicht, inklusive Hinweisen auf
Zuständigkeiten anderer Rollen (Verifier, Validator).

## Was dieser Skill NICHT macht

- Keine Lösungsvorschläge als Ersatz für Findings — der Reviewer kategorisiert,
  der Implementer entscheidet.
- Kein Refactoring-Vorschlag über den Diff-Scope hinaus.
- Keine DoD-Closure und keine Verification-Evidence — das ist Verifier-Aufgabe
  ([`harness/verification.md`](../../harness/verification.md)).
- Keine Validation gegen realen Nutzerbedarf — Validator-Aufgabe.
- Keine Abwertung eines Findings, weil seine Behebung unbequem ist.

Fällt etwas auf, das in diese Kategorien gehört: INFO-Finding mit Verweis auf
die zuständige Rolle.

## Output-Schema

Jedes Finding:

- `kategorie`: HIGH | MEDIUM | LOW | INFO
- `quelle`: `LH-*`-ID, ADR-Kennung, Hard-Rule-Name, `MR-*` oder
  „Maintainability"
- `pfad`: Datei:Zeile
- `befund`: 1–2 beobachtbare Sätze, ohne Lösungsvorschlag
- `risiko`: warum das relevant ist
- `verifizierbar`: welcher Sensor bestätigt es (`make lint`, `make test`,
  `make docs-check`, `make verify-depguard`, `make image-scan`) — oder
  „Review-only"

Am Ende jedes Laufs eine Negativbefund-Zeile pro betrachtetem Bereich
(„geprüft, ohne Befund: `<Pfad/Linse>`") und eine Zeile pro *nicht* geprüftem
Bereich mit Grund. Ohne diesen Block ist „keine Findings" nicht von „nicht
geprüft" unterscheidbar.

Ein Report pro Lauf unter [`docs/reviews/`](../../docs/reviews/); Folgeläufe
bekommen eine neue Datei, keine Überschreibung.

## Pflege (Steering-Loop)

Bei dreimaligem Auftreten desselben Findings:

1. Klassifikation hier **und** in [`harness/review.md`](../../harness/review.md)
   schärfen — die beiden dürfen nicht auseinanderlaufen.
2. Prüfen, ob [`AGENTS.md`](../../AGENTS.md), eine ADR oder die Spec eine Hard
   Rule braucht.
3. Prüfen, ob ein computational Sensor möglich ist (`make lint`,
   `make docs-check`, `make verify-depguard`, ein Test) — ein Gate ist billiger
   als wiederholtes inferentielles Nachlesen.
4. Wird das Muster nur temporär toleriert: Carveout-Inventar und Plan-Anker
   nachziehen.

Diese Datei wird nicht überschrieben, sondern versioniert fortgeschrieben.
